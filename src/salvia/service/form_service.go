package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ─── Errores de dominio ───────────────────────────────────────────────────────

var ErrFormNotFound        = errors.New("form: registro no encontrado")
var ErrFormSectionNotFound = errors.New("form_section: registro no encontrado")

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

// ─── FormSection ──────────────────────────────────────────────────────────────

type CreateFormSectionInput struct {
	FormID      string
	Name        string
	Description *string
	Order       int
}

type UpdateFormSectionInput struct {
	Name        *string
	Description *string
	Order       *int
}

type FormSectionService interface {
	GetByID(ctx context.Context, id string) (*models.FormSection, error)
	ListByFormID(ctx context.Context, formID string) ([]models.FormSection, error)
	Create(ctx context.Context, input CreateFormSectionInput) (*models.FormSection, error)
	Update(ctx context.Context, id string, input UpdateFormSectionInput) (*models.FormSection, error)
	Delete(ctx context.Context, id string) error
}

type formSectionService struct {
	repo repository.FormSectionRepository
}

func NewFormSectionService(repo repository.FormSectionRepository) FormSectionService {
	return &formSectionService{repo: repo}
}

func (s *formSectionService) GetByID(ctx context.Context, id string) (*models.FormSection, error) {
	section, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormSectionNotFound
		}
		return nil, err
	}
	return section, nil
}

func (s *formSectionService) ListByFormID(ctx context.Context, formID string) ([]models.FormSection, error) {
	return s.repo.FindByFormID(ctx, formID)
}

func (s *formSectionService) Create(ctx context.Context, input CreateFormSectionInput) (*models.FormSection, error) {
	section := &models.FormSection{
		FormID:      input.FormID,
		Name:        input.Name,
		Description: input.Description,
		Order:       input.Order,
	}
	if err := s.repo.Create(ctx, section); err != nil {
		return nil, err
	}
	return section, nil
}

func (s *formSectionService) Update(ctx context.Context, id string, input UpdateFormSectionInput) (*models.FormSection, error) {
	fields := map[string]interface{}{}
	if input.Name != nil        { fields["name"] = *input.Name }
	if input.Description != nil { fields["description"] = *input.Description }
	if input.Order != nil       { fields["order"] = *input.Order }

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormSectionNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *formSectionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFormSectionNotFound
	}
	return err
}
