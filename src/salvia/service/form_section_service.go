package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrFormSectionNotFound = errors.New("form_section: registro no encontrado")

type FormSectionService interface {
	GetByID(ctx context.Context, id string) (*models.FormSection, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.FormSection], error)
	Create(ctx context.Context, s *models.FormSection) error
	Update(ctx context.Context, s *models.FormSection) error
	Delete(ctx context.Context, id string) error
}

type formSectionService struct {
	repo repository.FormSectionRepository
}

func NewFormSectionService(repo repository.FormSectionRepository) FormSectionService {
	return &formSectionService{repo: repo}
}

func (s *formSectionService) GetByID(ctx context.Context, id string) (*models.FormSection, error) {
	fs, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormSectionNotFound
		}
		return nil, err
	}
	return fs, nil
}

func (s *formSectionService) List(ctx context.Context, page, limit int) (repository.PageResult[models.FormSection], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *formSectionService) Create(ctx context.Context, fs *models.FormSection) error {
	return s.repo.Create(ctx, fs)
}

func (s *formSectionService) Update(ctx context.Context, fs *models.FormSection) error {
	return s.repo.Update(ctx, fs)
}

func (s *formSectionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFormSectionNotFound
	}
	return err
}
