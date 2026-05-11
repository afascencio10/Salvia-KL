// Package service — case_detail_service.go
// Lógica de negocio para la pantalla de detalle de caso (rol sv).
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrCaseNotFound = errors.New("caso no encontrado")
var ErrInvalidICode = errors.New("icode inválido")

type CaseDetailService interface {
	GetDetail(ctx context.Context, caseICode string) (*repository.CaseDetailData, error)
	CreateFollowUp(ctx context.Context, caseICode, agentID, scheduledDate, notas string) (*models.FollowUpV2, error)
	AddTimelineEvent(ctx context.Context, caseICode, eventType, description, actorID, actorName string) error
	ReassignFollowUp(ctx context.Context, followUpID, newAgentID string) error
}

type caseDetailService struct {
	repo repository.CaseDetailRepository
}

func NewCaseDetailService(repo repository.CaseDetailRepository) CaseDetailService {
	return &caseDetailService{repo: repo}
}

func (s *caseDetailService) GetDetail(ctx context.Context, caseICode string) (*repository.CaseDetailData, error) {
	if caseICode == "" {
		return nil, ErrInvalidICode
	}
	detail, err := s.repo.GetByICode(ctx, caseICode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCaseNotFound
		}
		return nil, err
	}
	return detail, nil
}

func (s *caseDetailService) CreateFollowUp(ctx context.Context, caseICode, agentID, scheduledDate, notas string) (*models.FollowUpV2, error) {
	if caseICode == "" {
		return nil, ErrInvalidICode
	}

	fecha, err := time.Parse("2006-01-02T15:04", scheduledDate)
	if err != nil {
		fecha, err = time.Parse("2006-01-02", scheduledDate)
		if err != nil {
			return nil, errors.New("formato de fecha inválido, use YYYY-MM-DD o YYYY-MM-DDTHH:MM")
		}
	}

	// Validar que la fecha no sea del pasado
	hoy := time.Now().Truncate(24 * time.Hour)
	if fecha.Before(hoy) {
		return nil, errors.New("no se permiten seguimientos con fecha anterior a hoy")
	}

	// Extraer la hora como string para el campo scheduled_time
	hora := fecha.Format("15:04")

	// Calcular sequence_number: contar los existentes + 1
	count, err := s.repo.CountFollowUpsByCaseID(ctx, caseICode)
	if err != nil {
		return nil, err
	}

	// Validar máximo 8 seguimientos por caso
	if count >= 8 {
		return nil, errors.New("este caso ya tiene el máximo de 8 seguimientos permitidos")
	}

	nextSeq := count + 1

	sinEvaluar := "SIN_EVALUAR"
	followUp := &models.FollowUpV2{
		CaseID:         caseICode,
		AgentID:        &agentID,
		Status:         models.FollowUpStatusPendiente,
		ScheduledDate:  fecha,
		ScheduledTime:  hora,
		Summary:        &notas,
		Team:           "SIN_EQUIPO",
		RiskStatus:     &sinEvaluar,
		SequenceNumber: nextSeq,
	}

	if err := s.repo.CreateFollowUpV2(ctx, followUp); err != nil {
		return nil, err
	}

	// Registrar evento en el timeline
	s.repo.CreateTimelineEvent(ctx, &models.CaseTimelineEvent{
		CaseID:      caseICode,
		EventType:   models.TimelineEventSeguimiento,
		Description: "Seguimiento #" + fmt.Sprintf("%d", nextSeq) + " creado",
	})

	return followUp, nil
}

func (s *caseDetailService) AddTimelineEvent(ctx context.Context, caseICode, eventType, description, actorID, actorName string) error {
	return s.repo.CreateTimelineEvent(ctx, &models.CaseTimelineEvent{
		CaseID:      caseICode,
		EventType:   eventType,
		Description: description,
		ActorID:     actorID,
		ActorName:   actorName,
	})
}

func (s *caseDetailService) ReassignFollowUp(ctx context.Context, followUpID, newAgentID string) error {
	return s.repo.UpdateFollowUpAgent(ctx, followUpID, newAgentID)
}
