package repository

import (
	"bitsflow/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

// DirectoryFilters agrupa los criterios opcionales del listado de entidades.
type DirectoryFilters struct {
	CityICode    string
	Type         string
	Sector       string
	DepartmentID *uint64
	NameQuery    string
	IncludeInactive bool
}

type DirectoryRepository interface {
	Repository[models.Directory]
	FindByFilters(ctx context.Context, filters DirectoryFilters, page, limit int) ([]models.Directory, int64, error)
	SetActive(ctx context.Context, id string, active bool) error
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

func (r *directoryRepository) FindByFilters(ctx context.Context, f DirectoryFilters, page, limit int) ([]models.Directory, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page < 0 {
		page = 0
	}

	query := r.db.WithContext(ctx).Model(&models.Directory{})
	if !f.IncludeInactive {
		query = query.Where("active = ?", true)
	}
	if f.CityICode != "" {
		query = query.Where("city_id = ?", f.CityICode)
	}
	if f.Type != "" {
		query = query.Where("type = ?", f.Type)
	}
	if f.Sector != "" {
		query = query.Where("sector = ?", f.Sector)
	}
	if f.DepartmentID != nil {
		query = query.Where("department_id = ?", *f.DepartmentID)
	}
	if q := strings.TrimSpace(f.NameQuery); q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
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

func (r *directoryRepository) SetActive(ctx context.Context, id string, active bool) error {
	return r.db.WithContext(ctx).
		Model(&models.Directory{}).
		Where("id = ?", id).
		Update("active", active).Error
}
