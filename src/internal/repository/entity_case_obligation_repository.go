package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EntityCaseObligationRepository interface {
	Repository[models.EntityCaseObligation]
	FindByEntityCaseID(ctx context.Context, entityCaseID string) ([]models.EntityCaseObligation, error)
}

type entityCaseObligationRepository struct {
	repository[models.EntityCaseObligation]
	db *gorm.DB
}

func NewEntityCaseObligationRepository(db *gorm.DB) EntityCaseObligationRepository {
	return &entityCaseObligationRepository{
		repository: repository[models.EntityCaseObligation]{db: db},
		db:         db,
	}
}

func (r *entityCaseObligationRepository) FindByEntityCaseID(ctx context.Context, entityCaseID string) ([]models.EntityCaseObligation, error) {
	var items []models.EntityCaseObligation
	return items, r.db.WithContext(ctx).Where("entity_case_id = ?", entityCaseID).Order("created_at ASC").Find(&items).Error
}
