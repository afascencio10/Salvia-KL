package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EconomicStabilizationRepository interface {
	Repository[models.EconomicStabilization]
	FindByCaseID(ctx context.Context, caseID string) ([]models.EconomicStabilization, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.EconomicStabilization, error)
}

type economicStabilizationRepository struct {
	repository[models.EconomicStabilization]
	db *gorm.DB
}

func NewEconomicStabilizationRepository(db *gorm.DB) EconomicStabilizationRepository {
	return &economicStabilizationRepository{
		repository: repository[models.EconomicStabilization]{db: db},
		db:         db,
	}
}

func (r *economicStabilizationRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.EconomicStabilization, error) {
	var items []models.EconomicStabilization
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *economicStabilizationRepository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.EconomicStabilization, error) {
	var items []models.EconomicStabilization
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
