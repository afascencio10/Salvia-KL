package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// TownLightRepository define el acceso a datos para salvia.town (light)
type TownLightRepository interface {
	FindByCode(ctx context.Context, townCode string) (*models.TownLight, error)
}

type townLightRepository struct {
	db *gorm.DB
}

// NewTownLightRepository crea la instancia del repositorio
func NewTownLightRepository(db *gorm.DB) TownLightRepository {
	return &townLightRepository{db: db}
}

func (r *townLightRepository) FindByCode(ctx context.Context, townCode string) (*models.TownLight, error) {
	var town models.TownLight
	err := r.db.WithContext(ctx).Where("town_code = ?", townCode).First(&town).Error
	if err != nil {
		return nil, err
	}
	return &town, nil
}
