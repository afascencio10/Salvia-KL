package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrRepeaterGroupNotFound        = errors.New("repeater_group: registro no encontrado")
var ErrQuestionNotFound             = errors.New("question: registro no encontrado")
var ErrVisibilityConditionNotFound  = errors.New("visibility_condition: registro no encontrado")

// ─── RepeaterGroup ────────────────────────────────────────────────────────────

type CreateRepeaterGroupInput struct {
	FormSectionID  string
	Name           string
	Order          int
	MinRepetitions *int
	MaxRepetitions *int
	AddButtonLabel *string
}

type UpdateRepeaterGroupInput struct {
	Name           *string
	Order          *int
	MinRepetitions *int
	MaxRepetitions *int
	AddButtonLabel *string
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
	return s.repo.FindByFormSectionID(ctx, formSectionID)
}

func (s *repeaterGroupService) Create(ctx context.Context, input CreateRepeaterGroupInput) (*models.RepeaterGroup, error) {
	rg := &models.RepeaterGroup{
		FormSectionID:  input.FormSectionID,
		Name:           input.Name,
		Order:          input.Order,
		MinRepetitions: input.MinRepetitions,
		MaxRepetitions: input.MaxRepetitions,
		AddButtonLabel: input.AddButtonLabel,
	}
	return rg, s.repo.Create(ctx, rg)
}

func (s *repeaterGroupService) Update(ctx context.Context, id string, input UpdateRepeaterGroupInput) (*models.RepeaterGroup, error) {
	fields := map[string]interface{}{}
	if input.Name != nil           { fields["name"] = *input.Name }
	if input.Order != nil          { fields["order"] = *input.Order }
	if input.MinRepetitions != nil { fields["min_repetitions"] = *input.MinRepetitions }
	if input.MaxRepetitions != nil { fields["max_repetitions"] = *input.MaxRepetitions }
	if input.AddButtonLabel != nil { fields["add_button_label"] = *input.AddButtonLabel }
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

// ─── Question ─────────────────────────────────────────────────────────────────

type CreateQuestionInput struct {
	FormID          string
	FormSectionID   string
	RepeaterGroupID *string
	QuestionType    string
	Description     string
	Metadata        *string
	Options         *string
	Order           int
}

type UpdateQuestionInput struct {
	QuestionType    *string
	Description     *string
	Metadata        *string
	Options         *string
	Order           *int
}

type QuestionService interface {
	GetByID(ctx context.Context, id string) (*models.Question, error)
	ListByFormID(ctx context.Context, formID string) ([]models.Question, error)
	ListByFormSectionID(ctx context.Context, formSectionID string) ([]models.Question, error)
	Create(ctx context.Context, input CreateQuestionInput) (*models.Question, error)
	Update(ctx context.Context, id string, input UpdateQuestionInput) (*models.Question, error)
	Delete(ctx context.Context, id string) error
}

type questionService struct{ repo repository.QuestionRepository }

func NewQuestionService(repo repository.QuestionRepository) QuestionService {
	return &questionService{repo: repo}
}

func (s *questionService) GetByID(ctx context.Context, id string) (*models.Question, error) {
	q, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrQuestionNotFound }
	return q, err
}

func (s *questionService) ListByFormID(ctx context.Context, formID string) ([]models.Question, error) {
	return s.repo.FindByFormID(ctx, formID)
}

func (s *questionService) ListByFormSectionID(ctx context.Context, formSectionID string) ([]models.Question, error) {
	return s.repo.FindByFormSectionID(ctx, formSectionID)
}

func (s *questionService) Create(ctx context.Context, input CreateQuestionInput) (*models.Question, error) {
	q := &models.Question{
		FormID:          input.FormID,
		FormSectionID:   input.FormSectionID,
		RepeaterGroupID: input.RepeaterGroupID,
		QuestionType:    input.QuestionType,
		Description:     input.Description,
		Metadata:        input.Metadata,
		Options:         input.Options,
		Order:           input.Order,
	}
	return q, s.repo.Create(ctx, q)
}

func (s *questionService) Update(ctx context.Context, id string, input UpdateQuestionInput) (*models.Question, error) {
	fields := map[string]interface{}{}
	if input.QuestionType != nil { fields["question_type"] = *input.QuestionType }
	if input.Description != nil  { fields["description"] = *input.Description }
	if input.Metadata != nil     { fields["metadata"] = *input.Metadata }
	if input.Options != nil      { fields["options"] = *input.Options }
	if input.Order != nil        { fields["order"] = *input.Order }
	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrQuestionNotFound }
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *questionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return ErrQuestionNotFound }
	return err
}

// ─── VisibilityCondition ──────────────────────────────────────────────────────

type CreateVisibilityConditionInput struct {
	TargetType        string
	TargetID          string
	TriggerQuestionID string
	TriggerOptionID   *string
	TriggerValue      *string
	Operator          string
	Logic             string
}

type UpdateVisibilityConditionInput struct {
	TargetType        *string
	TargetID          *string
	TriggerQuestionID *string
	TriggerOptionID   *string
	TriggerValue      *string
	Operator          *string
	Logic             *string
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
		TriggerOptionID:   input.TriggerOptionID,
		TriggerValue:      input.TriggerValue,
		Operator:          input.Operator,
		Logic:             input.Logic,
	}
	return vc, s.repo.Create(ctx, vc)
}

func (s *visibilityConditionService) Update(ctx context.Context, id string, input UpdateVisibilityConditionInput) (*models.VisibilityCondition, error) {
	fields := map[string]interface{}{}
	if input.TargetType != nil        { fields["target_type"] = *input.TargetType }
	if input.TargetID != nil          { fields["target_id"] = *input.TargetID }
	if input.TriggerQuestionID != nil { fields["trigger_question_id"] = *input.TriggerQuestionID }
	if input.TriggerOptionID != nil   { fields["trigger_option_id"] = *input.TriggerOptionID }
	if input.TriggerValue != nil      { fields["trigger_value"] = *input.TriggerValue }
	if input.Operator != nil          { fields["operator"] = *input.Operator }
	if input.Logic != nil             { fields["logic"] = *input.Logic }
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
