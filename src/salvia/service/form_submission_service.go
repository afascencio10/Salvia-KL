package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrFormSubmissionNotFound = errors.New("form_submission: registro no encontrado")
var ErrRepeaterEntryNotFound  = errors.New("repeater_entry: registro no encontrado")
var ErrAnswerNotFound         = errors.New("answer: registro no encontrado")

// ─── FormSubmission ───────────────────────────────────────────────────────────

type FormSubmissionService interface {
	GetByID(ctx context.Context, id string) (*models.FormSubmission, error)
	ListByFormID(ctx context.Context, formID string) ([]models.FormSubmission, error)
	Create(ctx context.Context, formID string) (*models.FormSubmission, error)
	Delete(ctx context.Context, id string) error
}

type formSubmissionService struct{ repo repository.FormSubmissionRepository }

func NewFormSubmissionService(repo repository.FormSubmissionRepository) FormSubmissionService {
	return &formSubmissionService{repo: repo}
}

func (s *formSubmissionService) GetByID(ctx context.Context, id string) (*models.FormSubmission, error) {
	fs, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrFormSubmissionNotFound }
	return fs, err
}

func (s *formSubmissionService) ListByFormID(ctx context.Context, formID string) ([]models.FormSubmission, error) {
	return s.repo.FindByFormID(ctx, formID)
}

func (s *formSubmissionService) Create(ctx context.Context, formID string) (*models.FormSubmission, error) {
	fs := &models.FormSubmission{FormID: formID}
	return fs, s.repo.Create(ctx, fs)
}

func (s *formSubmissionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return ErrFormSubmissionNotFound }
	return err
}

// ─── RepeaterEntry ────────────────────────────────────────────────────────────

type CreateRepeaterEntryInput struct {
	FormSubmissionID string
	RepeaterGroupID  string
	Iteration        int
}

type RepeaterEntryService interface {
	GetByID(ctx context.Context, id string) (*models.RepeaterEntry, error)
	ListBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error)
	Create(ctx context.Context, input CreateRepeaterEntryInput) (*models.RepeaterEntry, error)
	Delete(ctx context.Context, id string) error
}

type repeaterEntryService struct{ repo repository.RepeaterEntryRepository }

func NewRepeaterEntryService(repo repository.RepeaterEntryRepository) RepeaterEntryService {
	return &repeaterEntryService{repo: repo}
}

func (s *repeaterEntryService) GetByID(ctx context.Context, id string) (*models.RepeaterEntry, error) {
	re, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrRepeaterEntryNotFound }
	return re, err
}

func (s *repeaterEntryService) ListBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error) {
	return s.repo.FindBySubmissionID(ctx, submissionID)
}

func (s *repeaterEntryService) Create(ctx context.Context, input CreateRepeaterEntryInput) (*models.RepeaterEntry, error) {
	re := &models.RepeaterEntry{
		FormSubmissionID: input.FormSubmissionID,
		RepeaterGroupID:  input.RepeaterGroupID,
		Iteration:        input.Iteration,
	}
	return re, s.repo.Create(ctx, re)
}

func (s *repeaterEntryService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return ErrRepeaterEntryNotFound }
	return err
}

// ─── Answer ───────────────────────────────────────────────────────────────────

type CreateAnswerInput struct {
	FormSubmissionID string
	QuestionID       string
	RepeaterEntryID  *string
	Value            *string
}

type UpdateAnswerInput struct {
	Value *string
}

type AnswerService interface {
	GetByID(ctx context.Context, id string) (*models.Answer, error)
	ListBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error)
	Create(ctx context.Context, input CreateAnswerInput) (*models.Answer, error)
	Update(ctx context.Context, id string, input UpdateAnswerInput) (*models.Answer, error)
	Delete(ctx context.Context, id string) error
}

type answerService struct{ repo repository.AnswerRepository }

func NewAnswerService(repo repository.AnswerRepository) AnswerService {
	return &answerService{repo: repo}
}

func (s *answerService) GetByID(ctx context.Context, id string) (*models.Answer, error) {
	a, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrAnswerNotFound }
	return a, err
}

func (s *answerService) ListBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error) {
	return s.repo.FindBySubmissionID(ctx, submissionID)
}

func (s *answerService) Create(ctx context.Context, input CreateAnswerInput) (*models.Answer, error) {
	a := &models.Answer{
		FormSubmissionID: input.FormSubmissionID,
		QuestionID:       input.QuestionID,
		RepeaterEntryID:  input.RepeaterEntryID,
		Value:            input.Value,
	}
	return a, s.repo.Create(ctx, a)
}

func (s *answerService) Update(ctx context.Context, id string, input UpdateAnswerInput) (*models.Answer, error) {
	fields := map[string]interface{}{}
	if input.Value != nil { fields["value"] = *input.Value }
	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrAnswerNotFound }
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *answerService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) { return ErrAnswerNotFound }
	return err
}
