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

// RepeaterEntry y Answer fueron movidos a archivos individuales:
// repeater_entry_repository.go, answer_repository.go
