package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// LocationRepository provee acceso de solo lectura a departamentos, ciudades y municipios.
type LocationRepository interface {
	GetDepartments(ctx context.Context) ([]models.LocationOption, error)
	GetCities(ctx context.Context) ([]models.CityLocationOption, error)
	GetTownsByCityID(ctx context.Context, cityID uint64) ([]models.LocationOption, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetDepartments(ctx context.Context) ([]models.LocationOption, error) {
	var rows []models.DepartmentLight
	if err := r.db.WithContext(ctx).
		Select("department_id, department_name").
		Order("department_name ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]models.LocationOption, len(rows))
	for i, d := range rows {
		out[i] = models.LocationOption{
			Label: d.DepartmentName,
			Value: fmt.Sprintf("%d", d.DepartmentId),
		}
	}
	return out, nil
}

func (r *locationRepository) GetCities(ctx context.Context) ([]models.CityLocationOption, error) {
	var rows []models.CityLight
	if err := r.db.WithContext(ctx).
		Select("city_id, city_name, department_id").
		Order("city_name ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]models.CityLocationOption, len(rows))
	for i, c := range rows {
		out[i] = models.CityLocationOption{
			Label:        c.CityName,
			Value:        fmt.Sprintf("%d", c.CityId),
			DepartmentId: c.DepartmentId,
		}
	}
	return out, nil
}

func (r *locationRepository) GetTownsByCityID(ctx context.Context, cityID uint64) ([]models.LocationOption, error) {
	var rows []models.TownLight
	if err := r.db.WithContext(ctx).
		Select("town_id, town_name").
		Where("city_id = ?", cityID).
		Order("town_name ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]models.LocationOption, len(rows))
	for i, t := range rows {
		out[i] = models.LocationOption{
			Label: t.TownName,
			Value: fmt.Sprintf("%d", t.ID),
		}
	}
	return out, nil
}
