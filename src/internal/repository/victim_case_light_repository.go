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
	UpdateStatus(ctx context.Context, caseID string, status string) error
	FindRiskLevelByICode(ctx context.Context, iCode string) (int, error)
	UpdateRiskLevelByICode(ctx context.Context, iCode string, newLevel int) error
	UpdateTeamAndAgent(ctx context.Context, iCode string, team string, agentID string) error
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

func (r *victimCaseLightRepository) UpdateStatus(ctx context.Context, caseID string, status string) error {
	return r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Where("victim_case_i_code = ?", caseID).
		Update("victim_case_status", status).Error
}

func (r *victimCaseLightRepository) FindRiskLevelByICode(ctx context.Context, iCode string) (int, error) {
	var level int
	err := r.db.WithContext(ctx).
		Table("salvia.victim_case_form2").
		Select("victim_case_form2_risk_level").
		Joins("JOIN salvia.victim_case ON victim_case.victim_case_id = victim_case_form2.victim_case_form2_victim_case").
		Where("victim_case.victim_case_i_code = ?", iCode).
		Limit(1).
		Scan(&level).Error
	return level, err
}

func (r *victimCaseLightRepository) UpdateRiskLevelByICode(ctx context.Context, iCode string, newLevel int) error {
	return r.db.WithContext(ctx).
		Table("salvia.victim_case_form2").
		Where("victim_case_form2_victim_case = (SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?)", iCode).
		Update("victim_case_form2_risk_level", newLevel).Error
}

func (r *victimCaseLightRepository) UpdateTeamAndAgent(ctx context.Context, iCode string, team string, agentID string) error {
	return r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Where("victim_case_i_code = ?", iCode).
		Updates(map[string]interface{}{
			"victim_case_team": team,
			"agent_id":         agentID,
		}).Error
}
