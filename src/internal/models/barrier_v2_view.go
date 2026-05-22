package models

import "time"

// BarrierV2WithRelations es el DTO de lectura enriquecido que combina:
//   - Todos los campos relevantes de salvia.barrier_v2
//   - victim_name, victim_last_name, victim_doc_number y case_code de salvia.victim_case
//
// El join se realiza:
//
//	barrier_v2.case_id → victim_case.victim_case_i_code
type BarrierV2WithRelations struct {
	// ── Campos de barrier_v2 ─────────────────────────────────────────────────
	ID          string    `gorm:"column:id"             json:"id"`
	CaseID      string    `gorm:"column:case_id"        json:"caseId"`
	FollowUpID  string    `gorm:"column:follow_up_id"   json:"followUpId"`
	Sector      string    `gorm:"column:sector"         json:"sector"`
	Description string    `gorm:"column:description"    json:"description"`
	Status      string    `gorm:"column:status"         json:"status"`
	CreatedByID string    `gorm:"column:created_by_id"  json:"createdById"`
	CreatedAt   time.Time `gorm:"column:created_at"     json:"createdAt"`

	// ── Campos de victim_case ─────────────────────────────────────────────────
	VictimName      string `gorm:"column:victim_name"       json:"victimName"`
	VictimLastName  string `gorm:"column:victim_last_name"  json:"victimLastName"`
	VictimDocNumber string `gorm:"column:victim_doc_number" json:"victimDocNumber"`
	CaseCode        string `gorm:"column:case_code"         json:"caseCode"`
}
