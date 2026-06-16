package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type DirectoryRepository interface {
	Repository[models.Directory]
	FindByFilters(ctx context.Context, cityICode string, dirType string, page, limit int) ([]models.Directory, int64, error)
}

type directoryRepository struct {
	repository[models.Directory]
	db *gorm.DB
}

func NewDirectoryRepository(db *gorm.DB) DirectoryRepository {
	return &directoryRepository{
		repository: repository[models.Directory]{db: db},
		db:         db,
	}
}

func (r *directoryRepository) FindByFilters(ctx context.Context, cityICode string, dirType string, page, limit int) ([]models.Directory, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page < 0 {
		page = 0
	}

	query := r.db.WithContext(ctx).Model(&models.Directory{})
	if cityICode != "" {
		query = query.Where("city_id = ?", cityICode)
	}
	if dirType != "" {
		query = query.Where("type = ?", dirType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Directory
	err := query.
		Order("name ASC").
		Offset(page * limit).
		Limit(limit).
		Find(&items).Error
	return items, total, err
}
