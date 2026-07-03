package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ReassignProfessionalItem profesional en select agrupado (RRM-03).
type ReassignProfessionalItem struct {
	ICode    string `json:"icode"`
	FullName string `json:"fullName"`
}

// ReassignProfessionalGroup grupo ps/ts para optgroup (RRM-03).
type ReassignProfessionalGroup struct {
	Role          string                     `json:"role"`
	Label         string                     `json:"label"`
	Professionals []ReassignProfessionalItem `json:"professionals"`
}

// ReassignProfessionalsResponse GET /api/v1/psychosocial-support/profesionales-reasignacion.
type ReassignProfessionalsResponse struct {
	Groups []ReassignProfessionalGroup `json:"groups"`
}

// DuplaReassignOption dupla con label enriquecido (RRM-03).
type DuplaReassignOption struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	PsychologistName string `json:"psychologistName"`
	SocialWorkerName string `json:"socialWorkerName"`
}

// DuplasReassignResponse GET /api/v1/duplas/reasignacion.
type DuplasReassignResponse struct {
	Duplas []DuplaReassignOption `json:"duplas"`
}

// PsychosocialReassignBulkInput parámetros POST reasignar-bulk (RRM-05).
type PsychosocialReassignBulkInput struct {
	RemisionIDs    []string
	AssignMode     string
	ProfessionalID string
	DuplaID        string
}

// PsychosocialReassignBulkResult respuesta POST reasignar-bulk.
type PsychosocialReassignBulkResult struct {
	OK                  bool  `json:"ok"`
	RemisionesUpdated   int   `json:"remisiones_updated"`
	TeamContactsUpdated int64 `json:"team_contacts_updated"`
}

// PsychosocialReassignService contrato RRM-03 y RRM-05.
type PsychosocialReassignService interface {
	ListProfessionalsForReassign(ctx context.Context) (ReassignProfessionalsResponse, error)
	ListDuplasForReassign(ctx context.Context) (DuplasReassignResponse, error)
	ReassignBulk(ctx context.Context, input PsychosocialReassignBulkInput) (PsychosocialReassignBulkResult, error)
}

type psychosocialReassignService struct {
	repo repository.PsychosocialReassignRepository
	db   *gorm.DB
}

func NewPsychosocialReassignService(repo repository.PsychosocialReassignRepository, db *gorm.DB) PsychosocialReassignService {
	return &psychosocialReassignService{repo: repo, db: db}
}

func (s *psychosocialReassignService) ListProfessionalsForReassign(ctx context.Context) (ReassignProfessionalsResponse, error) {
	rows, err := s.repo.ListActivePsTsProfessionals(ctx)
	if err != nil {
		return ReassignProfessionalsResponse{}, err
	}

	psGroup := ReassignProfessionalGroup{Role: "ps", Label: "Psicólogas", Professionals: []ReassignProfessionalItem{}}
	tsGroup := ReassignProfessionalGroup{Role: "ts", Label: "Trabajadoras Sociales", Professionals: []ReassignProfessionalItem{}}

	for _, row := range rows {
		name := strings.TrimSpace(row.FullName)
		if name == "" {
			continue
		}
		item := ReassignProfessionalItem{ICode: row.ICode, FullName: name}
		switch strings.ToLower(strings.TrimSpace(row.Role)) {
		case "ts":
			tsGroup.Professionals = append(tsGroup.Professionals, item)
		default:
			psGroup.Professionals = append(psGroup.Professionals, item)
		}
	}

	return ReassignProfessionalsResponse{
		Groups: []ReassignProfessionalGroup{psGroup, tsGroup},
	}, nil
}

func (s *psychosocialReassignService) ListDuplasForReassign(ctx context.Context) (DuplasReassignResponse, error) {
	rows, err := s.repo.ListActiveDuplasEnriched(ctx)
	if err != nil {
		return DuplasReassignResponse{}, err
	}

	duplas := make([]DuplaReassignOption, 0, len(rows))
	for _, row := range rows {
		psychName := strings.TrimSpace(row.PsychologistName)
		tsName := strings.TrimSpace(row.SocialWorkerName)
		name := strings.TrimSpace(row.Name)
		label := fmt.Sprintf("%s — %s + %s", name, psychName, tsName)
		duplas = append(duplas, DuplaReassignOption{
			ID:               row.ID,
			Name:             name,
			Label:            label,
			PsychologistName: psychName,
			SocialWorkerName: tsName,
		})
	}

	return DuplasReassignResponse{Duplas: duplas}, nil
}

func (s *psychosocialReassignService) ReassignBulk(ctx context.Context, input PsychosocialReassignBulkInput) (PsychosocialReassignBulkResult, error) {
	remisionIDs := uniqueNonEmptyStrings(input.RemisionIDs)
	if len(remisionIDs) == 0 {
		return PsychosocialReassignBulkResult{}, errors.New("debe indicar al menos una remisión")
	}

	assignMode := strings.ToLower(strings.TrimSpace(input.AssignMode))
	if assignMode != "professional" && assignMode != "dupla" {
		return PsychosocialReassignBulkResult{}, errors.New("modo de asignación inválido")
	}

	professionalID := strings.TrimSpace(input.ProfessionalID)
	duplaID := strings.TrimSpace(input.DuplaID)

	if assignMode == "professional" {
		if professionalID == "" {
			return PsychosocialReassignBulkResult{}, errors.New("debe indicar un profesional")
		}
		if _, err := s.repo.FindProfessionalForReassign(ctx, professionalID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PsychosocialReassignBulkResult{}, errors.New("profesional no encontrado o inactivo")
			}
			return PsychosocialReassignBulkResult{}, err
		}
	} else {
		if duplaID == "" {
			return PsychosocialReassignBulkResult{}, errors.New("debe indicar una dupla")
		}
		if _, err := s.repo.FindDuplaForReassign(ctx, duplaID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PsychosocialReassignBulkResult{}, errors.New("dupla no encontrada o inactiva")
			}
			return PsychosocialReassignBulkResult{}, err
		}
	}

	for _, remisionID := range remisionIDs {
		row, err := s.repo.FindPsychosocialForReassign(ctx, remisionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PsychosocialReassignBulkResult{}, fmt.Errorf("remisión no encontrada: %s", remisionID)
			}
			return PsychosocialReassignBulkResult{}, err
		}
		if strings.EqualFold(strings.TrimSpace(row.Status), models.PsychosocialSupportStatusCerrado) {
			return PsychosocialReassignBulkResult{}, errors.New("no se pueden reasignar remisiones en estado cerrado")
		}
	}

	var result PsychosocialReassignBulkResult

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewPsychosocialReassignRepository(tx)

		for _, remisionID := range remisionIDs {
			if assignMode == "professional" {
				if err := txRepo.UpdatePsychosocialProfessional(ctx, remisionID, professionalID); err != nil {
					return err
				}
				count, err := txRepo.UpdateTeamContactsProfessional(ctx, remisionID, professionalID)
				if err != nil {
					return err
				}
				result.TeamContactsUpdated += count
			} else {
				if err := txRepo.UpdatePsychosocialDupla(ctx, remisionID, duplaID); err != nil {
					return err
				}
				count, err := txRepo.UpdateTeamContactsDupla(ctx, remisionID, duplaID)
				if err != nil {
					return err
				}
				result.TeamContactsUpdated += count
			}
			result.RemisionesUpdated++
		}
		return nil
	})
	if err != nil {
		return PsychosocialReassignBulkResult{}, err
	}

	result.OK = true
	return result, nil
}
