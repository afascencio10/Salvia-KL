package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// DuplaRepository acceso a salvia.dupla.
type DuplaRepository interface {
	ListActive(ctx context.Context) ([]models.DuplaOption, error)
}

type duplaRepository struct {
	db *gorm.DB
}

func NewDuplaRepository(db *gorm.DB) DuplaRepository {
	return &duplaRepository{db: db}
}

func (r *duplaRepository) ListActive(ctx context.Context) ([]models.DuplaOption, error) {
	var items []models.DuplaOption
	err := r.db.WithContext(ctx).
		Table("salvia.dupla").
		Select("id, name").
		Where("deleted_at IS NULL").
		Order("name ASC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.DuplaOption{}
	}
	return items, nil
}
