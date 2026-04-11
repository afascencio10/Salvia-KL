package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EmergencyMeasureRepository interface {
	Repository[models.EmergencyMeasure]
	FindByCaseID(ctx context.Context, caseID string) ([]models.EmergencyMeasure, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.EmergencyMeasure, error)
}

type emergencyMeasureRepository struct {
	repository[models.EmergencyMeasure]
	db *gorm.DB
}

func NewEmergencyMeasureRepository(db *gorm.DB) EmergencyMeasureRepository {
	return &emergencyMeasureRepository{
		repository: repository[models.EmergencyMeasure]{db: db},
		db:         db,
	}
}

func (r *emergencyMeasureRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.EmergencyMeasure, error) {
	var items []models.EmergencyMeasure
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *emergencyMeasureRepository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.EmergencyMeasure, error) {
	var items []models.EmergencyMeasure
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
