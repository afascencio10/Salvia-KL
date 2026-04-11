package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type RepeaterEntryRepository interface {
	Repository[models.RepeaterEntry]
	FindBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error)
}

type repeaterEntryRepository struct {
	repository[models.RepeaterEntry]
	db *gorm.DB
}

func NewRepeaterEntryRepository(db *gorm.DB) RepeaterEntryRepository {
	return &repeaterEntryRepository{
		repository: repository[models.RepeaterEntry]{db: db},
		db:         db,
	}
}

func (r *repeaterEntryRepository) FindBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error) {
	var items []models.RepeaterEntry
	return items, r.db.WithContext(ctx).Where("form_submission_id = ?", submissionID).Find(&items).Error
}
