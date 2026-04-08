package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrEconomicStabilizationNotFound = errors.New("economic_stabilization: registro no encontrado")

type EconomicStabilizationService interface {
	GetByID(ctx context.Context, id string) (*models.EconomicStabilization, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.EconomicStabilization], error)
	Create(ctx context.Context, es *models.EconomicStabilization) error
	Update(ctx context.Context, es *models.EconomicStabilization) error
	Delete(ctx context.Context, id string) error
}

type economicStabilizationService struct {
	repo repository.EconomicStabilizationRepository
}

func NewEconomicStabilizationService(repo repository.EconomicStabilizationRepository) EconomicStabilizationService {
	return &economicStabilizationService{repo: repo}
}

func (s *economicStabilizationService) GetByID(ctx context.Context, id string) (*models.EconomicStabilization, error) {
	es, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEconomicStabilizationNotFound
		}
		return nil, err
	}
	return es, nil
}

func (s *economicStabilizationService) List(ctx context.Context, page, limit int) (repository.PageResult[models.EconomicStabilization], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *economicStabilizationService) Create(ctx context.Context, es *models.EconomicStabilization) error {
	if es.Status == "" {
		es.Status = "ACTIVE"
	}
	return s.repo.Create(ctx, es)
}

func (s *economicStabilizationService) Update(ctx context.Context, es *models.EconomicStabilization) error {
	return s.repo.Update(ctx, es)
}

func (s *economicStabilizationService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrEconomicStabilizationNotFound
	}
	return err
}
