package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type BarrierUpdateRepository interface {
	Repository[models.BarrierUpdate]
	FindByBarrierID(ctx context.Context, barrierID string) ([]models.BarrierUpdate, error)
}

type barrierUpdateRepository struct {
	repository[models.BarrierUpdate]
	db *gorm.DB
}

func NewBarrierUpdateRepository(db *gorm.DB) BarrierUpdateRepository {
	return &barrierUpdateRepository{
		repository: repository[models.BarrierUpdate]{db: db},
		db:         db,
	}
}

func (r *barrierUpdateRepository) FindByBarrierID(ctx context.Context, barrierID string) ([]models.BarrierUpdate, error) {
	var items []models.BarrierUpdate
	return items, r.db.WithContext(ctx).Where("barrier_id = ?", barrierID).Order("created_at ASC").Find(&items).Error
}
