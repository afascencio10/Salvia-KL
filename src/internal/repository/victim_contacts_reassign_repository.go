package repository

import (
	"context"

	"gorm.io/gorm"
)

// VictimContactReassignRow datos mínimos del reporte para reasignación masiva.
type VictimContactReassignRow struct {
	ICode string `gorm:"column:victim_contact_i_code"`
}

// VictimContactsReassignRepository operaciones de persistencia para
// reasignación masiva de reportes (Listado Reportes E-06). Mismo patrón que
// CasesReassignRepository, con una diferencia: GAP-3 (related-tables.md) —
// sin bloqueo de equipo, el reporte SIEMPRE hereda el equipo del agente
// destino (no solo si está vacío, a diferencia de UpdateVictimCaseAgent).
type VictimContactsReassignRepository interface {
	FindAgentForReassign(ctx context.Context, agentICode string) (*AgentReassignRow, error)
	FindVictimContactForReassign(ctx context.Context, icode string) (*VictimContactReassignRow, error)
	UpdateVictimContactAgent(ctx context.Context, icode, agentICode, agentTeam string) error
}

type victimContactsReassignRepository struct {
	db *gorm.DB
}

func NewVictimContactsReassignRepository(db *gorm.DB) VictimContactsReassignRepository {
	return &victimContactsReassignRepository{db: db}
}

func (r *victimContactsReassignRepository) FindAgentForReassign(ctx context.Context, agentICode string) (*AgentReassignRow, error) {
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

func (r *victimContactsReassignRepository) FindVictimContactForReassign(ctx context.Context, icode string) (*VictimContactReassignRow, error) {
	var row VictimContactReassignRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT victim_contact_i_code
		FROM salvia.victim_contact
		WHERE victim_contact_i_code = ?
		LIMIT 1
	`, icode).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ICode == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

// UpdateVictimContactAgent SIEMPRE pisa victim_contact_team (GAP-3) — a
// diferencia de los casos, un reporte no tiene equipo previo que proteger.
func (r *victimContactsReassignRepository) UpdateVictimContactAgent(ctx context.Context, icode, agentICode, agentTeam string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE salvia.victim_contact
		SET victim_contact_agent_id = ?,
		    victim_contact_team = ?
		WHERE victim_contact_i_code = ?
	`, agentICode, agentTeam, icode).Error
}
