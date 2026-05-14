package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type AgentLightRepository interface {
	FindByICode(ctx context.Context, iCode string) (*models.AgentLight, error)
	// FindAllByRoleAndTeam retorna los agentes cuyo rol coincide con roleCode
	// y cuyo equipo (general_user_team) coincide con team.
	// Usado por el algoritmo de auto-asignación para acotar el pool de candidatos
	// al mismo equipo que recibirá los nuevos seguimientos.
	FindAllByRoleAndTeam(ctx context.Context, roleCode, team string) ([]models.AgentLight, error)
}

type agentLightRepository struct {
	db *gorm.DB
}

func NewAgentLightRepository(db *gorm.DB) AgentLightRepository {
	return &agentLightRepository{db: db}
}

func (r *agentLightRepository) FindByICode(ctx context.Context, iCode string) (*models.AgentLight, error) {
	var agent models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user").
		Select("security.general_user.general_user_i_code, security.general_user_profile.general_user_profile_names, security.general_user_profile.general_user_profile_last_names").
		Joins("JOIN security.general_user_profile ON security.general_user.general_user_general_user_profile = security.general_user_profile.general_user_profile_id").
		Where("security.general_user.general_user_i_code = ?", iCode).
		Take(&agent).Error

	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// FindAllByRoleAndTeam retorna los agentes que tienen asignado roleCode como rol
// y cuyo campo general_user_team coincide con team.
// Replica el JOIN de GetGeneralUsersByRoleWithProfile del DAO legacy añadiendo
// el filtro de equipo para acotar el pool de candidatos a auto-asignación.
func (r *agentLightRepository) FindAllByRoleAndTeam(ctx context.Context, roleCode, team string) ([]models.AgentLight, error) {
	var agents []models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user").
		Select(
			"security.general_user.general_user_i_code, "+
				"security.general_user_profile.general_user_profile_names, "+
				"security.general_user_profile.general_user_profile_last_names, "+
				"security.general_user.general_user_team",
		).
		Joins("JOIN security.general_user_profile ON security.general_user_profile.general_user_profile_id = security.general_user.general_user_general_user_profile").
		Joins("JOIN security.rel_role_general_user ON security.rel_role_general_user.general_user_id = security.general_user.general_user_id").
		Joins("JOIN security.role ON security.role.role_id = security.rel_role_general_user.role_id").
		Where("security.role.role_code = ? AND security.general_user.general_user_team = ?", roleCode, team).
		Find(&agents).Error
	return agents, err
}
