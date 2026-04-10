package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type PsychosocialSupportRepository interface {
	FindByFollowUpID(ctx context.Context, followUpId string) ([]models.PsychosocialSupport, error)
}

type psychosocialSupportRepository struct {
	db *gorm.DB
}

func NewPsychosocialSupportRepository(db *gorm.DB) PsychosocialSupportRepository {
	return &psychosocialSupportRepository{db: db}
}

func (r *psychosocialSupportRepository) FindByFollowUpID(ctx context.Context, followUpId string) ([]models.PsychosocialSupport, error) {
	var results []models.PsychosocialSupport
	err := r.db.WithContext(ctx).Where("follow_up_id = ?", followUpId).Find(&results).Error
	return results, err
}
