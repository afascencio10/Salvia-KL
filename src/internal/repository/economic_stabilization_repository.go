package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EconomicStabilizationRepository interface {
	FindByFollowUpID(ctx context.Context, followUpId string) ([]models.EconomicStabilization, error)
}

type economic_stabilizationRepository struct {
	db *gorm.DB
}

func NewEconomicStabilizationRepository(db *gorm.DB) EconomicStabilizationRepository {
	return &economic_stabilizationRepository{db: db}
}

func (r *economic_stabilizationRepository) FindByFollowUpID(ctx context.Context, followUpId string) ([]models.EconomicStabilization, error) {
	var results []models.EconomicStabilization
	err := r.db.WithContext(ctx).Where("follow_up_id = ?", followUpId).Find(&results).Error
	return results, err
}
