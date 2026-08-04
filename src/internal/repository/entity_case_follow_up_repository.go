package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EntityCaseFollowUpRepository interface {
	Repository[models.EntityCaseFollowUp]
	FindByEntityCaseID(ctx context.Context, entityCaseID string) ([]models.EntityCaseFollowUp, error)
}

type entityCaseFollowUpRepository struct {
	repository[models.EntityCaseFollowUp]
	db *gorm.DB
}

func NewEntityCaseFollowUpRepository(db *gorm.DB) EntityCaseFollowUpRepository {
	return &entityCaseFollowUpRepository{
		repository: repository[models.EntityCaseFollowUp]{db: db},
		db:         db,
	}
}

func (r *entityCaseFollowUpRepository) FindByEntityCaseID(ctx context.Context, entityCaseID string) ([]models.EntityCaseFollowUp, error) {
	var items []models.EntityCaseFollowUp
	return items, r.db.WithContext(ctx).Where("entity_case_id = ?", entityCaseID).Order("created_at ASC").Find(&items).Error
}
