// Package repository provee repositorios GORM para Form y FormSection.
package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// ─── Form ────────────────────────────────────────────────────────────────────

// FormRepository extiende el CRUD genérico con consultas propias de Form.
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

// ─── FormSection ─────────────────────────────────────────────────────────────

// FormSectionRepository extiende el CRUD genérico con consultas propias de FormSection.
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
	err := r.db.WithContext(ctx).
		Where("form_id = ?", formID).
		Order(`"order" ASC`).
		Find(&items).Error
	return items, err
}
