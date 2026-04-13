package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type RepeaterEntryRepository interface {
	Repository[models.RepeaterEntry]
	FindBySubmissionID(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error)
	FindBySubmissionIDOrdered(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error)
	FindBySubmissionIDAndGroupIDs(ctx context.Context, submissionID string, groupIDs []string) ([]models.RepeaterEntry, error)
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

func (r *repeaterEntryRepository) FindBySubmissionIDOrdered(ctx context.Context, submissionID string) ([]models.RepeaterEntry, error) {
	var items []models.RepeaterEntry
	return items, r.db.WithContext(ctx).
		Where("form_submission_id = ?", submissionID).
		Order("repeater_group_id ASC, iteration ASC").
		Find(&items).Error
}

func (r *repeaterEntryRepository) FindBySubmissionIDAndGroupIDs(ctx context.Context, submissionID string, groupIDs []string) ([]models.RepeaterEntry, error) {
	var items []models.RepeaterEntry
	if len(groupIDs) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).
		Where("form_submission_id = ? AND repeater_group_id IN ?", submissionID, groupIDs).
		Order("iteration ASC").
		Find(&items).Error
}
