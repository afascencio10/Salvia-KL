package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrOptionNotFound = errors.New("option: registro no encontrado")

// ─── Option ───────────────────────────────────────────────────────────────────

type CreateOptionInput struct {
	QuestionID string
	Label      string
	Value      string
	Order      int
}

type UpdateOptionInput struct {
	Label *string
	Value *string
	Order *int
}

type OptionService interface {
	GetByID(ctx context.Context, id string) (*models.Option, error)
	ListByQuestionID(ctx context.Context, questionID string) ([]models.Option, error)
	Create(ctx context.Context, input CreateOptionInput) (*models.Option, error)
	Update(ctx context.Context, id string, input UpdateOptionInput) (*models.Option, error)
	Delete(ctx context.Context, id string) error
}

type optionService struct {
	repo repository.OptionRepository
}

func NewOptionService(repo repository.OptionRepository) OptionService {
	return &optionService{repo: repo}
}

func (s *optionService) GetByID(ctx context.Context, id string) (*models.Option, error) {
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOptionNotFound
		}
		return nil, err
	}
	return o, nil
}

func (s *optionService) ListByQuestionID(ctx context.Context, questionID string) ([]models.Option, error) {
	return s.repo.FindByQuestionID(ctx, questionID)
}

func (s *optionService) Create(ctx context.Context, input CreateOptionInput) (*models.Option, error) {
	o := &models.Option{
		QuestionID: input.QuestionID,
		Label:      input.Label,
		Value:      input.Value,
		Order:      input.Order,
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *optionService) Update(ctx context.Context, id string, input UpdateOptionInput) (*models.Option, error) {
	fields := map[string]interface{}{}
	if input.Label != nil { fields["label"] = *input.Label }
	if input.Value != nil { fields["value"] = *input.Value }
	if input.Order != nil { fields["order"] = *input.Order }

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOptionNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *optionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrOptionNotFound
	}
	return err
}
