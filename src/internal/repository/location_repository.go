package repository

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrCityNotFound = errors.New("location: ciudad no encontrada")

// LocationRepository provee acceso de solo lectura a departamentos, ciudades y municipios.
type LocationRepository interface {
	GetDepartments(ctx context.Context) ([]models.LocationOption, error)
	GetCities(ctx context.Context, departmentID *uint64) ([]models.CityLocationOption, error)
	GetTownsByCityID(ctx context.Context, cityID uint64) ([]models.LocationOption, error)
	FindBestCityMatchByName(ctx context.Context, cityName string) (*models.CityICodeLight, error)
	FindCityByICode(ctx context.Context, cityICode string) (*models.CityICodeLight, error)
	FindCityNamesByICodes(ctx context.Context, icodes []string) (map[string]string, error)
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

func (r *locationRepository) GetCities(ctx context.Context, departmentID *uint64) ([]models.CityLocationOption, error) {
	var rows []models.CityLight
	query := r.db.WithContext(ctx).
		Select("city_id, city_name, department_id").
		Order("city_name ASC")

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	if err := query.Find(&rows).Error; err != nil {
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
		Select("town_id, town_name, town_code").
		Where("city_id = ?", cityID).
		Order("town_name ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]models.LocationOption, len(rows))
	for i, t := range rows {
		out[i] = models.LocationOption{
			Label: t.TownName,
			Value: t.TownCode,
		}
	}
	return out, nil
}

func (r *locationRepository) FindBestCityMatchByName(ctx context.Context, cityName string) (*models.CityICodeLight, error) {
	normalizedInput := utils.NormalizeLocationName(cityName)
	if normalizedInput == "" {
		return nil, ErrCityNotFound
	}

	var rows []models.CityICodeLight
	query := r.db.WithContext(ctx).Select("city_i_code, city_name")

	if len(normalizedInput) >= 3 {
		query = query.Where("city_name ILIKE ?", normalizedInput[:3]+"%")
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		if err := r.db.WithContext(ctx).
			Select("city_i_code, city_name").
			Find(&rows).Error; err != nil {
			return nil, err
		}
	}

	var best *models.CityICodeLight
	bestScore := 0
	bestNormLen := 0

	for i := range rows {
		normalizedCity := utils.NormalizeLocationName(rows[i].CityName)
		score := utils.ScoreCityNameMatch(normalizedInput, normalizedCity)
		if score == 0 {
			continue
		}

		normLen := len(normalizedCity)
		if score > bestScore || (score == bestScore && (best == nil || normLen < bestNormLen)) {
			copyRow := rows[i]
			best = &copyRow
			bestScore = score
			bestNormLen = normLen
		}
	}

	if best == nil {
		return nil, fmt.Errorf("%w: %s", ErrCityNotFound, cityName)
	}

	return best, nil
}

func (r *locationRepository) FindCityByICode(ctx context.Context, cityICode string) (*models.CityICodeLight, error) {
	var city models.CityICodeLight
	err := r.db.WithContext(ctx).
		Select("city_i_code, city_name").
		Where("city_i_code = ?", cityICode).
		First(&city).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCityNotFound
		}
		return nil, err
	}
	return &city, nil
}

func (r *locationRepository) FindCityNamesByICodes(ctx context.Context, icodes []string) (map[string]string, error) {
	out := make(map[string]string, len(icodes))
	if len(icodes) == 0 {
		return out, nil
	}

	var rows []models.CityICodeLight
	if err := r.db.WithContext(ctx).
		Select("city_i_code, city_name").
		Where("city_i_code IN ?", icodes).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		out[row.CityICode] = row.CityName
	}
	return out, nil
}
