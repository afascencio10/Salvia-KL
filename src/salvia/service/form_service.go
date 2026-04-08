package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrFormNotFound = errors.New("form: registro no encontrado")

// FormService define el contrato de negocio para Form.
type FormService interface {
	GetByID(ctx context.Context, id string) (*models.Form, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error)
	Create(ctx context.Context, f *models.Form) error
	Update(ctx context.Context, f *models.Form) error
	Delete(ctx context.Context, id string) error
}

type formService struct {
	repo repository.FormRepository
}

func NewFormService(repo repository.FormRepository) FormService {
	return &formService{repo: repo}
}

func (s *formService) GetByID(ctx context.Context, id string) (*models.Form, error) {
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormNotFound
		}
		return nil, err
	}
	return f, nil
}

func (s *formService) List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *formService) Create(ctx context.Context, f *models.Form) error {
	if f.Status == "" {
		f.Status = "ACTIVE"
	}
	return s.repo.Create(ctx, f)
}

func (s *formService) Update(ctx context.Context, f *models.Form) error {
	return s.repo.Update(ctx, f)
}

func (s *formService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFormNotFound
	}
	return err
}
