package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrFormNotFound = errors.New("form: registro no encontrado")

// ─── Form ─────────────────────────────────────────────────────────────────────

type CreateFormInput struct {
	Name        string
	Description string
	Status      string
}

type UpdateFormInput struct {
	Name        *string
	Description *string
	Status      *string
}

type FormService interface {
	GetByID(ctx context.Context, id string) (*models.Form, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error)
	Create(ctx context.Context, input CreateFormInput) (*models.Form, error)
	Update(ctx context.Context, id string, input UpdateFormInput) (*models.Form, error)
	Delete(ctx context.Context, id string) error
}

type formService struct {
	repo repository.FormRepository
}

func NewFormService(repo repository.FormRepository) FormService {
	return &formService{repo: repo}
}

func (s *formService) GetByID(ctx context.Context, id string) (*models.Form, error) {
	form, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormNotFound
		}
		return nil, err
	}
	return form, nil
}

func (s *formService) List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *formService) Create(ctx context.Context, input CreateFormInput) (*models.Form, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}
	form := &models.Form{
		Name:        input.Name,
		Description: input.Description,
		Status:      status,
	}
	if err := s.repo.Create(ctx, form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *formService) Update(ctx context.Context, id string, input UpdateFormInput) (*models.Form, error) {
	fields := map[string]interface{}{}
	if input.Name != nil        { fields["name"] = *input.Name }
	if input.Description != nil { fields["description"] = *input.Description }
	if input.Status != nil      { fields["status"] = *input.Status }

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *formService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFormNotFound
	}
	return err
}

// FormSection fue movido a form_section_service.go
