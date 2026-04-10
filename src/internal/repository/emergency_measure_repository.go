package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EmergencyMeasureRepository interface {
	FindByFollowUpID(ctx context.Context, followUpId string) ([]models.EmergencyMeasure, error)
}

type emergencyMeasureRepository struct {
	db *gorm.DB
}

func NewEmergencyMeasureRepository(db *gorm.DB) EmergencyMeasureRepository {
	return &emergencyMeasureRepository{db: db}
}

func (r *emergencyMeasureRepository) FindByFollowUpID(ctx context.Context, followUpId string) ([]models.EmergencyMeasure, error) {
	var results []models.EmergencyMeasure
	err := r.db.WithContext(ctx).Where("follow_up_id = ?", followUpId).Find(&results).Error
	return results, err
}
