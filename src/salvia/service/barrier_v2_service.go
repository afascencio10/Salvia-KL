package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrBarrierV2NotFound = errors.New("barrier_v2: registro no encontrado")

type BarrierV2Service interface {
	GetByID(ctx context.Context, id string) (*models.BarrierV2, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.BarrierV2], error)
	Create(ctx context.Context, b *models.BarrierV2) error
	Update(ctx context.Context, b *models.BarrierV2) error
	Delete(ctx context.Context, id string) error
}

type barrierV2Service struct {
	repo repository.BarrierV2Repository
}

func NewBarrierV2Service(repo repository.BarrierV2Repository) BarrierV2Service {
	return &barrierV2Service{repo: repo}
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
