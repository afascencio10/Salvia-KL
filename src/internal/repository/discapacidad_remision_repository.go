package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type DiscapacidadRemisionRepository interface {
	Repository[models.DiscapacidadRemision]
	FindByCaseID(ctx context.Context, caseID string) ([]models.DiscapacidadRemision, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.DiscapacidadRemision, error)
}

type discapacidadRemisionRepository struct {
	repository[models.DiscapacidadRemision]
	db *gorm.DB
}

func NewDiscapacidadRemisionRepository(db *gorm.DB) DiscapacidadRemisionRepository {
	return &discapacidadRemisionRepository{
		repository: repository[models.DiscapacidadRemision]{db: db},
		db:         db,
	}
}

func (r *discapacidadRemisionRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.DiscapacidadRemision, error) {
	var items []models.DiscapacidadRemision
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *discapacidadRemisionRepository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.DiscapacidadRemision, error) {
	var items []models.DiscapacidadRemision
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
