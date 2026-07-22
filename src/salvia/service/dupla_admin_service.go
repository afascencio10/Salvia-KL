package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrDuplaNameInUse       = errors.New("nombre de la dupla en uso")
	ErrDuplaMembersConflict = errors.New("Uno o ambos profesionales ya pertenecen a otra dupla")
	ErrDuplaInvalidMember   = errors.New("Profesional inválido o inactivo")
	ErrDuplaValidation      = errors.New("datos de dupla inválidos")
	ErrDuplaNotFound        = errors.New("dupla no encontrada")
	ErrDuplaInUse           = errors.New("dupla_en_uso")
)

const DuplaInUseMessage = "Esta dupla no se puede eliminar porque se está utilizando en una sesión."

// DuplaAdminProfessionalItem profesional ps/ts para Administrar Duplas (E01).
type DuplaAdminProfessionalItem struct {
	ICode    string `json:"icode"`
	FullName string `json:"fullName"`
}

// DuplaAdminProfessionalsResponse GET /api/v1/duplas/profesionales.
type DuplaAdminProfessionalsResponse struct {
	Psychologists []DuplaAdminProfessionalItem `json:"psychologists"`
	SocialWorkers []DuplaAdminProfessionalItem `json:"socialWorkers"`
}

// DuplaAdminActivasResponse GET /api/v1/duplas/activas.
type DuplaAdminActivasResponse struct {
	Items []models.DuplaAdminItem `json:"items"`
}

// DuplaAdminSaveInput payload create/update (E04).
type DuplaAdminSaveInput struct {
	Name           string
	PsychologistID string
	SocialWorkerID string
}

// DuplaAdminService lecturas y escritura para pantalla Administrar Duplas.
type DuplaAdminService interface {
	ListProfessionals(ctx context.Context) (DuplaAdminProfessionalsResponse, error)
	ListActiveEnriched(ctx context.Context) (DuplaAdminActivasResponse, error)
	Create(ctx context.Context, input DuplaAdminSaveInput) (*models.DuplaAdminItem, error)
	Update(ctx context.Context, id string, input DuplaAdminSaveInput) (*models.DuplaAdminItem, error)
	SoftDelete(ctx context.Context, id string) error
}

type duplaAdminService struct {
	duplaRepo    repository.DuplaRepository
	reassignRepo repository.PsychosocialReassignRepository
}

func NewDuplaAdminService(
	duplaRepo repository.DuplaRepository,
	reassignRepo repository.PsychosocialReassignRepository,
) DuplaAdminService {
	return &duplaAdminService{
		duplaRepo:    duplaRepo,
		reassignRepo: reassignRepo,
	}
}

func (s *duplaAdminService) ListProfessionals(ctx context.Context) (DuplaAdminProfessionalsResponse, error) {
	rows, err := s.reassignRepo.ListActivePsTsProfessionals(ctx)
	if err != nil {
		return DuplaAdminProfessionalsResponse{}, err
	}

	psychologists := make([]DuplaAdminProfessionalItem, 0)
	socialWorkers := make([]DuplaAdminProfessionalItem, 0)

	for _, row := range rows {
		name := strings.TrimSpace(row.FullName)
		if name == "" || strings.TrimSpace(row.ICode) == "" {
			continue
		}
		item := DuplaAdminProfessionalItem{ICode: row.ICode, FullName: name}
		switch strings.ToLower(strings.TrimSpace(row.Role)) {
		case "ts":
			socialWorkers = append(socialWorkers, item)
		case "ps":
			psychologists = append(psychologists, item)
		}
	}

	return DuplaAdminProfessionalsResponse{
		Psychologists: psychologists,
		SocialWorkers: socialWorkers,
	}, nil
}

func (s *duplaAdminService) ListActiveEnriched(ctx context.Context) (DuplaAdminActivasResponse, error) {
	items, err := s.duplaRepo.ListActiveEnriched(ctx)
	if err != nil {
		return DuplaAdminActivasResponse{}, err
	}
	if items == nil {
		items = []models.DuplaAdminItem{}
	}
	return DuplaAdminActivasResponse{Items: items}, nil
}

func (s *duplaAdminService) Create(ctx context.Context, input DuplaAdminSaveInput) (*models.DuplaAdminItem, error) {
	normalized, err := s.normalizeSaveInput(input)
	if err != nil {
		return nil, err
	}
	if err := s.validateSave(ctx, "", normalized); err != nil {
		return nil, err
	}

	created, err := s.duplaRepo.Create(ctx, normalized.Name, normalized.PsychologistID, normalized.SocialWorkerID)
	if err != nil {
		return nil, err
	}
	return s.duplaRepo.FindEnrichedByID(ctx, created.ID)
}

func (s *duplaAdminService) Update(ctx context.Context, id string, input DuplaAdminSaveInput) (*models.DuplaAdminItem, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrDuplaNotFound
	}

	normalized, err := s.normalizeSaveInput(input)
	if err != nil {
		return nil, err
	}
	if err := s.validateSave(ctx, id, normalized); err != nil {
		return nil, err
	}

	if err := s.duplaRepo.UpdateActive(ctx, id, normalized.Name, normalized.PsychologistID, normalized.SocialWorkerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDuplaNotFound
		}
		return nil, err
	}
	item, err := s.duplaRepo.FindEnrichedByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDuplaNotFound
		}
		return nil, err
	}
	return item, nil
}

func (s *duplaAdminService) SoftDelete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrDuplaNotFound
	}

	exists, err := s.duplaRepo.ExistsActiveByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrDuplaNotFound
	}

	inUse, err := s.duplaRepo.IsInUseNonClosed(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrDuplaInUse
	}

	if err := s.duplaRepo.SoftDelete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDuplaNotFound
		}
		return err
	}
	return nil
}

func (s *duplaAdminService) normalizeSaveInput(input DuplaAdminSaveInput) (DuplaAdminSaveInput, error) {
	name := strings.TrimSpace(input.Name)
	psychID := strings.TrimSpace(input.PsychologistID)
	tsID := strings.TrimSpace(input.SocialWorkerID)

	if name == "" || psychID == "" || tsID == "" {
		return DuplaAdminSaveInput{}, ErrDuplaValidation
	}
	if len(name) > 36 {
		return DuplaAdminSaveInput{}, ErrDuplaValidation
	}
	if psychID == tsID {
		return DuplaAdminSaveInput{}, ErrDuplaValidation
	}

	return DuplaAdminSaveInput{
		Name:           name,
		PsychologistID: psychID,
		SocialWorkerID: tsID,
	}, nil
}

func (s *duplaAdminService) validateSave(ctx context.Context, excludeID string, input DuplaAdminSaveInput) error {
	nameTaken, err := s.duplaRepo.ExistsActiveByName(ctx, input.Name, excludeID)
	if err != nil {
		return err
	}
	if nameTaken {
		return ErrDuplaNameInUse
	}

	ps, err := s.reassignRepo.FindProfessionalForReassign(ctx, input.PsychologistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDuplaInvalidMember
		}
		return err
	}
	if strings.ToLower(strings.TrimSpace(ps.Role)) != "ps" {
		return ErrDuplaInvalidMember
	}

	ts, err := s.reassignRepo.FindProfessionalForReassign(ctx, input.SocialWorkerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDuplaInvalidMember
		}
		return err
	}
	if strings.ToLower(strings.TrimSpace(ts.Role)) != "ts" {
		return ErrDuplaInvalidMember
	}

	conflict, err := s.duplaRepo.ExistsActiveMemberConflict(ctx, input.PsychologistID, input.SocialWorkerID, excludeID)
	if err != nil {
		return err
	}
	if conflict {
		return ErrDuplaMembersConflict
	}

	return nil
}
