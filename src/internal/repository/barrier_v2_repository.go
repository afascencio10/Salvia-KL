package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// BarrierV2Repository define el acceso a datos para barrier_v2
type BarrierV2Repository interface {
	FindByFollowUpID(ctx context.Context, followUpId string) ([]models.BarrierV2, error)
}

type barrierV2Repository struct {
	db *gorm.DB
}

// NewBarrierV2Repository crea la instancia del repositorio
func NewBarrierV2Repository(db *gorm.DB) BarrierV2Repository {
	return &barrierV2Repository{db: db}
}

func (r *barrierV2Repository) FindByFollowUpID(ctx context.Context, followUpId string) ([]models.BarrierV2, error) {
	var barriers []models.BarrierV2
	err := r.db.WithContext(ctx).Where("follow_up_id = ?", followUpId).Find(&barriers).Error
	return barriers, err
}
