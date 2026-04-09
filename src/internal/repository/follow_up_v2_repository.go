package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type FollowUpV2Repository interface {
	Repository[models.FollowUpV2]
	FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindPending(ctx context.Context) ([]models.FollowUpV2, error)
}

type followUpV2Repository struct {
	repository[models.FollowUpV2]
	db *gorm.DB
}

func NewFollowUpV2Repository(db *gorm.DB) FollowUpV2Repository {
	return &followUpV2Repository{
		repository: repository[models.FollowUpV2]{db: db},
		db:         db,
	}
}

func (r *followUpV2Repository) FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *followUpV2Repository) FindPending(ctx context.Context) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	return items, r.db.WithContext(ctx).Where("is_completed = false").Find(&items).Error
}
