package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type BarrierV2Repository interface {
	Repository[models.BarrierV2]
	FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.BarrierV2, error)
}

type barrierV2Repository struct {
	repository[models.BarrierV2]
	db *gorm.DB
}

func NewBarrierV2Repository(db *gorm.DB) BarrierV2Repository {
	return &barrierV2Repository{
		repository: repository[models.BarrierV2]{db: db},
		db:         db,
	}
}

func (r *barrierV2Repository) FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *barrierV2Repository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
