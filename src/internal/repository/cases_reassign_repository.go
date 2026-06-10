package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// VictimCaseReassignRow datos mínimos del caso para reasignación masiva.
type VictimCaseReassignRow struct {
	ICode string `gorm:"column:victim_case_i_code"`
	Team  string `gorm:"column:victim_case_team"`
}

// AgentReassignRow agente destino con nombre y equipo.
type AgentReassignRow struct {
	ICode    string `gorm:"column:general_user_i_code"`
	FullName string `gorm:"column:full_name"`
	Team     string `gorm:"column:general_user_team"`
	Status   string `gorm:"column:general_user_status"`
}

// CasesReassignRepository operaciones de persistencia para reasignación masiva (M-05).
type CasesReassignRepository interface {
	FindAgentForReassign(ctx context.Context, agentICode string) (*AgentReassignRow, error)
	FindVictimCaseForReassign(ctx context.Context, caseICode string) (*VictimCaseReassignRow, error)
	UpdateVictimCaseAgent(ctx context.Context, caseICode, agentICode, agentTeam string) error
	UpdateFollowUpsForReassign(ctx context.Context, caseICode, agentICode, agentTeam string) (int64, error)
}

type casesReassignRepository struct {
	db *gorm.DB
}

func NewCasesReassignRepository(db *gorm.DB) CasesReassignRepository {
	return &casesReassignRepository{db: db}
}

func (r *casesReassignRepository) FindAgentForReassign(ctx context.Context, agentICode string) (*AgentReassignRow, error) {
	var row AgentReassignRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT gu.general_user_i_code,
		       TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')) AS full_name,
		       COALESCE(gu.general_user_team, '') AS general_user_team,
		       COALESCE(gu.general_user_status, '') AS general_user_status
		FROM security.general_user gu
		LEFT JOIN security.general_user_profile gup
		       ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE gu.general_user_i_code = ?
		LIMIT 1
	`, agentICode).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ICode == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *casesReassignRepository) FindVictimCaseForReassign(ctx context.Context, caseICode string) (*VictimCaseReassignRow, error) {
	var row VictimCaseReassignRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT victim_case_i_code,
		       COALESCE(victim_case_team, '') AS victim_case_team
		FROM salvia.victim_case
		WHERE victim_case_i_code = ?
		LIMIT 1
	`, caseICode).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ICode == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *casesReassignRepository) UpdateVictimCaseAgent(ctx context.Context, caseICode, agentICode, agentTeam string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE salvia.victim_case
		SET agent_id = ?,
		    victim_case_team = CASE
		        WHEN COALESCE(TRIM(victim_case_team), '') = '' THEN ?
		        ELSE victim_case_team
		    END
		WHERE victim_case_i_code = ?
	`, agentICode, agentTeam, caseICode).Error
}

func (r *casesReassignRepository) UpdateFollowUpsForReassign(ctx context.Context, caseICode, agentICode, agentTeam string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE salvia.follow_up_v2
		SET agent_id = ?,
		    team = ?,
		    updated_at = ?
		WHERE case_id = ?
		  AND status NOT IN ('REALIZADO', 'CERRADO')
		  AND deleted_at IS NULL
	`, agentICode, agentTeam, time.Now(), caseICode)
	return result.RowsAffected, result.Error
}
