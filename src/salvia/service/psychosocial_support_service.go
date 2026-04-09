package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrPsychosocialSupportNotFound = errors.New("psychosocial_support: registro no encontrado")

type PsychosocialSupportService interface {
	GetByID(ctx context.Context, id string) (*models.PsychosocialSupport, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.PsychosocialSupport], error)
	Create(ctx context.Context, ps *models.PsychosocialSupport) error
	Update(ctx context.Context, ps *models.PsychosocialSupport) error
	Delete(ctx context.Context, id string) error
}

type psychosocialSupportService struct {
	repo repository.PsychosocialSupportRepository
}

func NewPsychosocialSupportService(repo repository.PsychosocialSupportRepository) PsychosocialSupportService {
	return &psychosocialSupportService{repo: repo}
}

func (s *psychosocialSupportService) GetByID(ctx context.Context, id string) (*models.PsychosocialSupport, error) {
	ps, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPsychosocialSupportNotFound
		}
		return nil, err
	}
	return ps, nil
}

func (s *psychosocialSupportService) List(ctx context.Context, page, limit int) (repository.PageResult[models.PsychosocialSupport], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *psychosocialSupportService) Create(ctx context.Context, ps *models.PsychosocialSupport) error {
	if ps.Status == "" {
		ps.Status = "ACTIVE"
	}
	return s.repo.Create(ctx, ps)
}

func (s *psychosocialSupportService) Update(ctx context.Context, ps *models.PsychosocialSupport) error {
	return s.repo.Update(ctx, ps)
}

func (s *psychosocialSupportService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrPsychosocialSupportNotFound
	}
	return err
}
