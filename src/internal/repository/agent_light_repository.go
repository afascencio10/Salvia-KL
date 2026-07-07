package repository

import (
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

type AgentLightRepository interface {
	FindByICode(ctx context.Context, iCode string) (*models.AgentLight, error)
	FindByLogin(ctx context.Context, login string) (*models.AgentLight, error)
	// FindAllByRoleAndTeam retorna los agentes activos (general_user_status = 'e')
	FindAllByRoleAndTeam(ctx context.Context, roleCode, team string) ([]models.AgentLight, error)
	// FindAllByRole retorna agentes activos con el rol indicado (sin filtro de equipo).
	FindAllByRole(ctx context.Context, roleCode string) ([]models.AgentLight, error)
	// SearchByName busca agentes activos por nombre/apellido (E-07 casos-component).
	SearchByName(ctx context.Context, query string, limit int) ([]models.AgentLight, error)
	// SearchPsicosocialByName busca profesionales activos del equipo psicosocial (E-07 remisiones).
	SearchPsicosocialByName(ctx context.Context, query string, limit int) ([]models.AgentPsicosocialLight, error)
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
		Select("security.general_user.general_user_i_code, security.general_user.general_user_login, security.general_user_profile.general_user_profile_names, security.general_user_profile.general_user_profile_last_names, security.general_user.general_user_team").
		Joins("JOIN security.general_user_profile ON security.general_user.general_user_general_user_profile = security.general_user_profile.general_user_profile_id").
		Where("security.general_user.general_user_i_code = ?", iCode).
		Take(&agent).Error

	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *agentLightRepository) FindByLogin(ctx context.Context, login string) (*models.AgentLight, error) {
	var agent models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user").
		Select("security.general_user.general_user_i_code, security.general_user.general_user_login, security.general_user_profile.general_user_profile_names, security.general_user_profile.general_user_profile_last_names, security.general_user.general_user_team").
		Joins("JOIN security.general_user_profile ON security.general_user.general_user_general_user_profile = security.general_user_profile.general_user_profile_id").
		Where("LOWER(security.general_user.general_user_login) = LOWER(?)", login).
		Take(&agent).Error

	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// FindAllByRoleAndTeam retorna los agentes activos que tienen asignado roleCode como rol
// y cuyo campo general_user_team coincide con team.
//
// Usa LEFT JOIN con general_user_profile para incluir agentes que aún no tienen
// perfil creado (general_user_general_user_profile = NULL). El JOIN con las tablas
// de roles sigue siendo INNER JOIN para garantizar el filtro por rol.
func (r *agentLightRepository) FindAllByRoleAndTeam(ctx context.Context, roleCode, team string) ([]models.AgentLight, error) {
	var agents []models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user").
		Select(
			"security.general_user.general_user_i_code, "+
				"COALESCE(security.general_user_profile.general_user_profile_names, '') AS general_user_profile_names, "+
				"COALESCE(security.general_user_profile.general_user_profile_last_names, '') AS general_user_profile_last_names, "+
				"security.general_user.general_user_team",
		).
		Joins("LEFT JOIN security.general_user_profile ON security.general_user_profile.general_user_profile_id = security.general_user.general_user_general_user_profile").
		Joins("JOIN security.rel_role_general_user ON security.rel_role_general_user.general_user_id = security.general_user.general_user_id").
		Joins("JOIN security.role ON security.role.role_id = security.rel_role_general_user.role_id").
		Where("security.role.role_code = ? AND security.general_user.general_user_team = ? AND security.general_user.general_user_status = ?", roleCode, team, "e").
		Find(&agents).Error
	return agents, err
}

// FindAllByRole retorna agentes activos con roleCode asignado, ordenados por apellido y nombre.
func (r *agentLightRepository) FindAllByRole(ctx context.Context, roleCode string) ([]models.AgentLight, error) {
	var agents []models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user gu").
		Select(
			"gu.general_user_i_code, "+
				"COALESCE(gup.general_user_profile_names, '') AS general_user_profile_names, "+
				"COALESCE(gup.general_user_profile_last_names, '') AS general_user_profile_last_names, "+
				"gu.general_user_team",
		).
		Joins("LEFT JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile").
		Joins("JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id").
		Joins("JOIN security.role r ON r.role_id = rr.role_id").
		Where("r.role_code = ? AND gu.general_user_status = ?", roleCode, "e").
		Order("gup.general_user_profile_last_names ASC, gup.general_user_profile_names ASC").
		Find(&agents).Error
	return agents, err
}

// allowedAgentSearchTeams equipos válidos para el autocomplete de persona asignada (E-07).
var allowedAgentSearchTeams = []string{"Riesgo bajo", "Riesgo alto", "Riesgo Alto", "Hombres"}

const agentSearchMinQueryLen = 3

// SearchByName retorna agentes activos cuyo nombre o apellido contiene query (ILIKE).
// Filtra general_user_status = 'e' y general_user_team en los equipos Salvia permitidos.
func (r *agentLightRepository) SearchByName(ctx context.Context, query string, limit int) ([]models.AgentLight, error) {
	query = strings.TrimSpace(query)
	if len(query) < agentSearchMinQueryLen {
		return []models.AgentLight{}, nil
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	pattern := "%" + query + "%"
	var agents []models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user gu").
		Select(
			"gu.general_user_i_code, "+
				"COALESCE(gup.general_user_profile_names, '') AS general_user_profile_names, "+
				"COALESCE(gup.general_user_profile_last_names, '') AS general_user_profile_last_names, "+
				"gu.general_user_team",
		).
		Joins("JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile").
		Where("gu.general_user_status = ?", "e").
		Where("gu.general_user_team IN ?", allowedAgentSearchTeams).
		Where(
			"(gup.general_user_profile_names ILIKE ? OR gup.general_user_profile_last_names ILIKE ? OR (gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names) ILIKE ?)",
			pattern, pattern, pattern,
		).
		Order("gup.general_user_profile_last_names ASC, gup.general_user_profile_names ASC").
		Limit(limit).
		Find(&agents).Error
	return agents, err
}

const agentPsicosocialSpecialtySelect = `(
	SELECT CASE
		WHEN COUNT(*) FILTER (
			WHERE LOWER(r.role_name) LIKE '%psicolog%'
			   OR LOWER(r.role_code) LIKE '%psi%'
			   OR LOWER(r.role_code) = 'psicologia'
		) > 0 THEN 'psicologia'
		WHEN COUNT(*) FILTER (
			WHERE LOWER(r.role_name) LIKE '%trab%social%'
			   OR LOWER(r.role_code) LIKE '%ts%'
			   OR LOWER(r.role_code) IN ('trab. social', 'trab_social', 'tss')
		) > 0 THEN 'trab. social'
		ELSE ''
	END
	FROM security.rel_role_general_user rr
	JOIN security.role r ON r.role_id = rr.role_id
	WHERE rr.general_user_id = gu.general_user_id
)`

// SearchPsicosocialByName retorna profesionales activos del equipo psicosocial (E-07).
func (r *agentLightRepository) SearchPsicosocialByName(ctx context.Context, query string, limit int) ([]models.AgentPsicosocialLight, error) {
	query = strings.TrimSpace(query)
	if len(query) < agentSearchMinQueryLen {
		return []models.AgentPsicosocialLight{}, nil
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	pattern := "%" + query + "%"
	sql := `
SELECT
	gu.general_user_i_code AS icode,
	COALESCE(gup.general_user_profile_names, '') AS names,
	COALESCE(gup.general_user_profile_last_names, '') AS last_names,
	gu.general_user_team AS team,
	` + agentPsicosocialSpecialtySelect + ` AS specialty
FROM security.general_user gu
JOIN security.general_user_profile gup
     ON gup.general_user_profile_id = gu.general_user_general_user_profile
WHERE gu.general_user_status = 'e'
  AND LOWER(BTRIM(gu.general_user_team)) = 'psicosocial'
  AND (
        gup.general_user_profile_names ILIKE ?
     OR gup.general_user_profile_last_names ILIKE ?
     OR (gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names) ILIKE ?
      )
ORDER BY gup.general_user_profile_last_names ASC, gup.general_user_profile_names ASC
LIMIT ?`

	var agents []models.AgentPsicosocialLight
	err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(sql, pattern, pattern, pattern, limit).Scan(&agents).Error
	})
	if err != nil {
		return nil, err
	}
	if agents == nil {
		agents = []models.AgentPsicosocialLight{}
	}
	return agents, nil
}
