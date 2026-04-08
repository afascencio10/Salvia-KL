package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// ─── FormSubmission ───────────────────────────────────────────────────────────

type FormSubmissionRepository interface {
	Repository[models.FormSubmission]
	FindByFormID(ctx context.Context, formID string) ([]models.FormSubmission, error)
}

type formSubmissionRepository struct {
	repository[models.FormSubmission]
	db *gorm.DB
}

func NewFormSubmissionRepository(db *gorm.DB) FormSubmissionRepository {
	return &formSubmissionRepository{repository: repository[models.FormSubmission]{db: db}, db: db}
}

func (r *formSubmissionRepository) FindByFormID(ctx context.Context, formID string) ([]models.FormSubmission, error) {
	var items []models.FormSubmission
	err := r.db.WithContext(ctx).Where("form_id = ?", formID).Find(&items).Error
	return items, err
}

// ─── RepeaterEntry ────────────────────────────────────────────────────────────

type RepeaterEntryRepository interface {
	Repository[models.RepeaterEntry]
	FindBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error)
}

type repeaterEntryRepository struct {
	repository[models.RepeaterEntry]
	db *gorm.DB
}

func NewRepeaterEntryRepository(db *gorm.DB) RepeaterEntryRepository {
	return &repeaterEntryRepository{repository: repository[models.RepeaterEntry]{db: db}, db: db}
}

func (r *repeaterEntryRepository) FindBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error) {
	var items []models.RepeaterEntry
	err := r.db.WithContext(ctx).Where("form_submission_id = ?", submissionID).Order("iteration ASC").Find(&items).Error
	return items, err
}

// ─── Answer ───────────────────────────────────────────────────────────────────

type AnswerRepository interface {
	Repository[models.Answer]
	FindBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error)
}

type answerRepository struct {
	repository[models.Answer]
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{repository: repository[models.Answer]{db: db}, db: db}
}

func (r *answerRepository) FindBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error) {
	var items []models.Answer
	err := r.db.WithContext(ctx).Where("form_submission_id = ?", submissionID).Find(&items).Error
	return items, err
}
