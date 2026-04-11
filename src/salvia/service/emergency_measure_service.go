package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrEmergencyMeasureNotFound = errors.New("emergency_measure: registro no encontrado")

type EmergencyMeasureService interface {
	GetByID(ctx context.Context, id string) (*models.EmergencyMeasure, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.EmergencyMeasure], error)
	Create(ctx context.Context, em *models.EmergencyMeasure) error
	Update(ctx context.Context, em *models.EmergencyMeasure) error
	Delete(ctx context.Context, id string) error
}

type emergencyMeasureService struct {
	repo repository.EmergencyMeasureRepository
}

func NewEmergencyMeasureService(repo repository.EmergencyMeasureRepository) EmergencyMeasureService {
	return &emergencyMeasureService{repo: repo}
}

func (s *emergencyMeasureService) GetByID(ctx context.Context, id string) (*models.EmergencyMeasure, error) {
	em, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmergencyMeasureNotFound
		}
		return nil, err
	}
	return em, nil
}

func (s *emergencyMeasureService) List(ctx context.Context, page, limit int) (repository.PageResult[models.EmergencyMeasure], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *emergencyMeasureService) Create(ctx context.Context, em *models.EmergencyMeasure) error {
	if em.Status == "" {
		em.Status = "ACTIVE"
	}
	return s.repo.Create(ctx, em)
}

func (s *emergencyMeasureService) Update(ctx context.Context, em *models.EmergencyMeasure) error {
	return s.repo.Update(ctx, em)
}

func (s *emergencyMeasureService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrEmergencyMeasureNotFound
	}
	return err
}
