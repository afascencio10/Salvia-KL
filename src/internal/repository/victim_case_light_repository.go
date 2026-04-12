package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// VictimCaseLightRepository define el acceso a datos para victim_case (light)
type VictimCaseLightRepository interface {
	FindByID(ctx context.Context, caseId string) (*models.VictimCaseLight, error)
}

type victimCaseLightRepository struct {
	db *gorm.DB
}

// NewVictimCaseLightRepository crea la instancia del repositorio
func NewVictimCaseLightRepository(db *gorm.DB) VictimCaseLightRepository {
	return &victimCaseLightRepository{db: db}
}

func (r *victimCaseLightRepository) FindByID(ctx context.Context, caseId string) (*models.VictimCaseLight, error) {
	var vcase models.VictimCaseLight
	err := r.db.WithContext(ctx).Where("victim_case_i_code = ?", caseId).First(&vcase).Error
	if err != nil {
		return nil, err
	}
	return &vcase, nil
}
