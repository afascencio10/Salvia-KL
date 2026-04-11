package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type FormSectionRepository interface {
	Repository[models.FormSection]
	FindByFormID(ctx context.Context, formID string) ([]models.FormSection, error)
}

type formSectionRepository struct {
	repository[models.FormSection]
	db *gorm.DB
}

func NewFormSectionRepository(db *gorm.DB) FormSectionRepository {
	return &formSectionRepository{
		repository: repository[models.FormSection]{db: db},
		db:         db,
	}
}

func (r *formSectionRepository) FindByFormID(ctx context.Context, formID string) ([]models.FormSection, error) {
	var items []models.FormSection
	return items, r.db.WithContext(ctx).Where("form_id = ?", formID).Order(`"order" ASC`).Find(&items).Error
}
