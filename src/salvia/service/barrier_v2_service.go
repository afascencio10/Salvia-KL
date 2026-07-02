package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var ErrBarrierV2NotFound = errors.New("barrier_v2: registro no encontrado")

// BarrierFollowUpItem es el DTO de respuesta para el endpoint de seguimientos de barrera.
type BarrierFollowUpItem struct {
	models.BarrierFollowUp
	ActorName string `json:"actorName"`
}

type BarrierV2Service interface {
	GetByID(ctx context.Context, id string) (*models.BarrierV2, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.BarrierV2], error)
	Create(ctx context.Context, b *models.BarrierV2) error
	Update(ctx context.Context, b *models.BarrierV2) error
	Delete(ctx context.Context, id string) error
	// ListByCreatedByIDWithRelations devuelve las barreras creadas por el agente
	// enriquecidas con datos de victim_case (nombres, doc, case_code).
	ListByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error)
	// ListFollowUps devuelve los seguimientos de una barrera en orden cronológico
	// con el nombre del autor resuelto.
	ListFollowUps(ctx context.Context, barrierID string) ([]BarrierFollowUpItem, error)
}

type barrierV2Service struct {
	repo           repository.BarrierV2Repository
	followUpRepo   repository.BarrierFollowUpRepository
	agentLightRepo repository.AgentLightRepository
}

func NewBarrierV2Service(repo repository.BarrierV2Repository, followUpRepo repository.BarrierFollowUpRepository, agentLightRepo repository.AgentLightRepository) BarrierV2Service {
	return &barrierV2Service{repo: repo, followUpRepo: followUpRepo, agentLightRepo: agentLightRepo}
}

func (s *barrierV2Service) GetByID(ctx context.Context, id string) (*models.BarrierV2, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBarrierV2NotFound
		}
		return nil, err
	}
	return b, nil
}

func (s *barrierV2Service) List(ctx context.Context, page, limit int) (repository.PageResult[models.BarrierV2], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *barrierV2Service) Create(ctx context.Context, b *models.BarrierV2) error {
	if b.Status == "" {
		b.Status = "OPEN"
	}
	return s.repo.Create(ctx, b)
}

func (s *barrierV2Service) Update(ctx context.Context, b *models.BarrierV2) error {
	return s.repo.Update(ctx, b)
}

func (s *barrierV2Service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrBarrierV2NotFound
	}
	return err
}

func (s *barrierV2Service) ListByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error) {
	return s.repo.FindByCreatedByIDWithRelations(ctx, createdByID)
}

func (s *barrierV2Service) ListFollowUps(ctx context.Context, barrierID string) ([]BarrierFollowUpItem, error) {
	records, err := s.followUpRepo.FindByBarrierID(ctx, barrierID)
	if err != nil {
		return nil, err
	}
	items := make([]BarrierFollowUpItem, len(records))
	for i, r := range records {
		item := BarrierFollowUpItem{BarrierFollowUp: r}
		if s.agentLightRepo != nil && r.CreatedByID != "" {
			if agent, err := s.agentLightRepo.FindByICode(ctx, r.CreatedByID); err == nil && agent != nil {
				item.ActorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
			}
		}
		items[i] = item
	}
	return items, nil
}
