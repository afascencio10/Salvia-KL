package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type RepeaterGroupRepository interface {
	Repository[models.RepeaterGroup]
	FindBySectionID(ctx context.Context, sectionID string) ([]models.RepeaterGroup, error)
}

type repeaterGroupRepository struct {
	repository[models.RepeaterGroup]
	db *gorm.DB
}

func NewRepeaterGroupRepository(db *gorm.DB) RepeaterGroupRepository {
	return &repeaterGroupRepository{
		repository: repository[models.RepeaterGroup]{db: db},
		db:         db,
	}
}

func (r *repeaterGroupRepository) FindBySectionID(ctx context.Context, sectionID string) ([]models.RepeaterGroup, error) {
	var items []models.RepeaterGroup
	return items, r.db.WithContext(ctx).Where("form_section_id = ?", sectionID).Find(&items).Error
}
