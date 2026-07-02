package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type BarrierFollowUpRepository interface {
	Repository[models.BarrierFollowUp]
	FindByBarrierID(ctx context.Context, barrierID string) ([]models.BarrierFollowUp, error)
}

type barrierFollowUpRepository struct {
	repository[models.BarrierFollowUp]
	db *gorm.DB
}

func NewBarrierFollowUpRepository(db *gorm.DB) BarrierFollowUpRepository {
	return &barrierFollowUpRepository{
		repository: repository[models.BarrierFollowUp]{db: db},
		db:         db,
	}
}

func (r *barrierFollowUpRepository) FindByBarrierID(ctx context.Context, barrierID string) ([]models.BarrierFollowUp, error) {
	var items []models.BarrierFollowUp
	return items, r.db.WithContext(ctx).Where("barrier_id = ?", barrierID).Order("created_at ASC").Find(&items).Error
}
