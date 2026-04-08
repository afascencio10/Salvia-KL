package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type FormSubmissionRepository interface {
	Repository[models.FormSubmission]
	FindByFormID(ctx context.Context, formID string) ([]models.FormSubmission, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.FormSubmission, error)
}

type formSubmissionRepository struct {
	repository[models.FormSubmission]
	db *gorm.DB
}

func NewFormSubmissionRepository(db *gorm.DB) FormSubmissionRepository {
	return &formSubmissionRepository{
		repository: repository[models.FormSubmission]{db: db},
		db:         db,
	}
}

func (r *formSubmissionRepository) FindByFormID(ctx context.Context, formID string) ([]models.FormSubmission, error) {
	var items []models.FormSubmission
	return items, r.db.WithContext(ctx).Where("form_id = ?", formID).Find(&items).Error
}

func (r *formSubmissionRepository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.FormSubmission, error) {
	var items []models.FormSubmission
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
