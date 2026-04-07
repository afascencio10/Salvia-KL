// Package service contiene la lógica de negocio de la capa salvia (Fase 2).
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrFollowUpNotFound se retorna cuando el registro no existe o fue eliminado.
var ErrFollowUpNotFound = errors.New("followup: registro no encontrado")

// FollowUpV2Service define el contrato de negocio para FollowUpV2.
type FollowUpV2Service interface {
	GetFollowUpByID(ctx context.Context, id string) (*models.FollowUpV2, error)
	GetPaginatedFollowUps(ctx context.Context, page, limit int) (repository.PageResult[models.FollowUpV2], error)
}

type followUpV2Service struct {
	repo repository.FollowUpRepository
}

// NewFollowUpV2Service construye el servicio inyectando el repositorio.
func NewFollowUpV2Service(repo repository.FollowUpRepository) FollowUpV2Service {
	return &followUpV2Service{repo: repo}
}

// GetFollowUpByID busca un FollowUpV2 por su ID.
// Traduce gorm.ErrRecordNotFound al error de dominio ErrFollowUpNotFound.
func (s *followUpV2Service) GetFollowUpByID(ctx context.Context, id string) (*models.FollowUpV2, error) {
	fu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowUpNotFound
		}
		return nil, err
	}
	return fu, nil
}

// GetPaginatedFollowUps retorna una página de FollowUpV2.
// page es base-0; limit <= 0 usa el valor por defecto del repositorio (20).
func (s *followUpV2Service) GetPaginatedFollowUps(ctx context.Context, page, limit int) (repository.PageResult[models.FollowUpV2], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}
