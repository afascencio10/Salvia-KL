package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// ─── Form ────────────────────────────────────────────────────────────────────

type FormRepository interface {
	Repository[models.Form]
	FindByStatus(ctx context.Context, status string) ([]models.Form, error)
}

type formRepository struct {
	repository[models.Form]
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) FormRepository {
	return &formRepository{
		repository: repository[models.Form]{db: db},
		db:         db,
	}
}

func (r *formRepository) FindByStatus(ctx context.Context, status string) ([]models.Form, error) {
	var items []models.Form
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&items).Error
	return items, err
}

// FormSection fue movido a form_section_repository.go
