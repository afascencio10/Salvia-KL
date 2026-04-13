package service

// RepeaterGroup, Question y VisibilityCondition fueron movidos a archivos individuales:
// question_service.go (Question), y los demás pendientes de separar.

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrRepeaterGroupNotFound       = errors.New("repeater_group: registro no encontrado")
var ErrVisibilityConditionNotFound = errors.New("visibility_condition: registro no encontrado")

// ─── RepeaterGroup ────────────────────────────────────────────────────────────

type CreateRepeaterGroupInput struct {
	FormSectionID  string
	Name           string
	ItemName       *string
	Order          int
	MinRepetitions int
	MaxRepetitions *int
}

type UpdateRepeaterGroupInput struct {
	Name           *string
	ItemName       *string
	Order          *int
	MinRepetitions *int
	MaxRepetitions *int
}

type RepeaterGroupService interface {
	GetByID(ctx context.Context, id string) (*models.RepeaterGroup, error)
	ListByFormSectionID(ctx context.Context, formSectionID string) ([]models.RepeaterGroup, error)
	Create(ctx context.Context, input CreateRepeaterGroupInput) (*models.RepeaterGroup, error)
	Update(ctx context.Context, id string, input UpdateRepeaterGroupInput) (*models.RepeaterGroup, error)
	Delete(ctx context.Context, id string) error
}

type repeaterGroupService struct{ repo repository.RepeaterGroupRepository }

func NewRepeaterGroupService(repo repository.RepeaterGroupRepository) RepeaterGroupService {
	return &repeaterGroupService{repo: repo}
}

func (s *repeaterGroupService) GetByID(ctx context.Context, id string) (*models.RepeaterGroup, error) {
	rg, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrRepeaterGroupNotFound }
	return rg, err
}

func (s *repeaterGroupService) ListByFormSectionID(ctx context.Context, formSectionID string) ([]models.RepeaterGroup, error) {
	return s.repo.FindBySectionID(ctx, formSectionID)
}

func (s *repeaterGroupService) Create(ctx context.Context, input CreateRepeaterGroupInput) (*models.RepeaterGroup, error) {
	rg := &models.RepeaterGroup{
		FormSectionID:  input.FormSectionID,
		Name:           input.Name,
		ItemName:       input.ItemName,
		Order:          input.Order,
		MinRepetitions: input.MinRepetitions,
		MaxRepetitions: input.MaxRepetitions,
	}
	return rg, s.repo.Create(ctx, rg)
}

func (s *repeaterGroupService) Update(ctx context.Context, id string, input UpdateRepeaterGroupInput) (*models.RepeaterGroup, error) {
	fields := map[string]interface{}{}
	if input.Name != nil           { fields["name"] = *input.Name }
	if input.ItemName != nil       { fields["item_name"] = *input.ItemName }
	if input.Order != nil          { fields["order"] = *input.Order }
	if input.MinRepetitions != nil { fields["min_repetitions"] = *input.MinRepetitions }
	if input.MaxRepetitions != nil { fields["max_repetitions"] = *input.MaxRepetitions }
	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrRepeaterGroupNotFound }
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *repeaterGroupService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return ErrRepeaterGroupNotFound }
	return err
}

// ─── VisibilityCondition ──────────────────────────────────────────────────────

type CreateVisibilityConditionInput struct {
	TargetType        string
	TargetID          string
	TriggerQuestionID string
	TriggerValue      *string
	Operator          string
}

type UpdateVisibilityConditionInput struct {
	TargetType        *string
	TargetID          *string
	TriggerQuestionID *string
	TriggerValue      *string
	Operator          *string
}

type VisibilityConditionService interface {
	GetByID(ctx context.Context, id string) (*models.VisibilityCondition, error)
	ListByTargetID(ctx context.Context, targetID string) ([]models.VisibilityCondition, error)
	Create(ctx context.Context, input CreateVisibilityConditionInput) (*models.VisibilityCondition, error)
	Update(ctx context.Context, id string, input UpdateVisibilityConditionInput) (*models.VisibilityCondition, error)
	Delete(ctx context.Context, id string) error
}

type visibilityConditionService struct{ repo repository.VisibilityConditionRepository }

func NewVisibilityConditionService(repo repository.VisibilityConditionRepository) VisibilityConditionService {
	return &visibilityConditionService{repo: repo}
}

func (s *visibilityConditionService) GetByID(ctx context.Context, id string) (*models.VisibilityCondition, error) {
	vc, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrVisibilityConditionNotFound }
	return vc, err
}

func (s *visibilityConditionService) ListByTargetID(ctx context.Context, targetID string) ([]models.VisibilityCondition, error) {
	return s.repo.FindByTargetID(ctx, targetID)
}

func (s *visibilityConditionService) Create(ctx context.Context, input CreateVisibilityConditionInput) (*models.VisibilityCondition, error) {
	vc := &models.VisibilityCondition{
		TargetType:        input.TargetType,
		TargetID:          input.TargetID,
		TriggerQuestionID: input.TriggerQuestionID,
		TriggerValue:      input.TriggerValue,
		Operator:          input.Operator,
	}
	return vc, s.repo.Create(ctx, vc)
}

func (s *visibilityConditionService) Update(ctx context.Context, id string, input UpdateVisibilityConditionInput) (*models.VisibilityCondition, error) {
	fields := map[string]interface{}{}
	if input.TriggerValue != nil { fields["trigger_value"] = *input.TriggerValue }
	if input.Operator != nil    { fields["operator"] = *input.Operator }
	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrVisibilityConditionNotFound }
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *visibilityConditionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return ErrVisibilityConditionNotFound }
	return err
}
