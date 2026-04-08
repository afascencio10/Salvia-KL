package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// FormRepository extiende el CRUD genérico con búsquedas propias de Form.
type FormRepository interface {
	Repository[models.Form]
	FindByStatus(ctx context.Context, status string) ([]models.Form, error)
	FindByCampaignID(ctx context.Context, campaignID string) ([]models.Form, error)
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
	return items, r.db.WithContext(ctx).Where("status = ?", status).Find(&items).Error
}

func (r *formRepository) FindByCampaignID(ctx context.Context, campaignID string) ([]models.Form, error) {
	var items []models.Form
	return items, r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Find(&items).Error
}
