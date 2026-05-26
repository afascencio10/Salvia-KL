package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type RenderModificationRepository interface {
	Repository[models.RenderModification]
	FindByTargetIDs(ctx context.Context, ids []string) ([]models.RenderModification, error)
}

type renderModificationRepository struct {
	repository[models.RenderModification]
	db *gorm.DB
}

func NewRenderModificationRepository(db *gorm.DB) RenderModificationRepository {
	return &renderModificationRepository{
		repository: repository[models.RenderModification]{db: db},
		db:         db,
	}
}

func (r *renderModificationRepository) FindByTargetIDs(ctx context.Context, ids []string) ([]models.RenderModification, error) {
	var items []models.RenderModification
	if len(ids) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).Where("target_id IN ?", ids).Find(&items).Error
}
