package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// VictimCaseLightRepository define el acceso a datos para victim_case (light)
type VictimCaseLightRepository interface {
	FindByID(ctx context.Context, caseId string) (*models.VictimCaseLight, error)
	FindByICode(ctx context.Context, iCode string) (*models.VictimCaseLight, error)
	UpdateStatus(ctx context.Context, iCode string, status string) error
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
	err := r.db.WithContext(ctx).Where("victim_case_id::text = ?", caseId).First(&vcase).Error
	if err != nil {
		return nil, err
	}
	return &vcase, nil
}

func (r *victimCaseLightRepository) FindByICode(ctx context.Context, iCode string) (*models.VictimCaseLight, error) {
	var vcase models.VictimCaseLight
	err := r.db.WithContext(ctx).Where("victim_case_i_code = ?", iCode).First(&vcase).Error
	if err != nil {
		return nil, err
	}
	return &vcase, nil
}

func (r *victimCaseLightRepository) UpdateStatus(ctx context.Context, iCode string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.VictimCaseLight{}).
		Where("victim_case_i_code = ?", iCode).
		Update("victim_case_status", status).Error
}
