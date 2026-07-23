package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
)

// ErrEntityCaseDuplicate se retorna cuando ya existe una relación activa entre
// el caso y la sede de entidad indicados.
var ErrEntityCaseDuplicate = errors.New("entity_case: la entidad ya está asociada a este caso")

// CreateEntityCaseInput agrupa los datos necesarios para asociar una entidad a un caso.
type CreateEntityCaseInput struct {
	CaseID         string
	EntityBranchID int64
	Objetivo       *string
	CreatedByID    string
}

// EntityCaseService define las operaciones de negocio sobre EntityCase
// (componente case-entities — "Gestión institucional" de Detalle del Caso).
type EntityCaseService interface {
	// ListByCase devuelve las entidades relacionadas con un caso, ya enriquecidas
	// con datos de sede, ubicación, oficios y barreras activas.
	ListByCase(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error)

	// Create asocia una sede de entidad a un caso. Retorna ErrEntityCaseDuplicate
	// si ya existe una relación activa entre ambos.
	Create(ctx context.Context, input CreateEntityCaseInput) (*models.EntityCase, error)
}

type entityCaseService struct {
	repo repository.EntityCaseRepository
}

func NewEntityCaseService(repo repository.EntityCaseRepository) EntityCaseService {
	return &entityCaseService{repo: repo}
}

func (s *entityCaseService) ListByCase(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error) {
	return s.repo.FindByCaseIDWithRelations(ctx, caseID)
}

func (s *entityCaseService) Create(ctx context.Context, input CreateEntityCaseInput) (*models.EntityCase, error) {
	exists, err := s.repo.ExistsActive(ctx, input.CaseID, input.EntityBranchID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEntityCaseDuplicate
	}

	ec := &models.EntityCase{
		CaseID:         input.CaseID,
		EntityBranchID: input.EntityBranchID,
		Objetivo:       input.Objetivo,
		CreatedByID:    input.CreatedByID,
	}
	if err := s.repo.Create(ctx, ec); err != nil {
		return nil, err
	}
	return ec, nil
}
