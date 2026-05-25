package models

import "time"

// CaseTaskWithRelations es el DTO de lectura enriquecido para la pantalla
// "Mis Barreras". Combina los campos principales de case_task con datos
// de barrier_v2 y victim_case obtenidos mediante LEFT JOIN.
//
// Joins:
//
//	case_task.barrier_id → barrier_v2.id::text
//	case_task.case_id    → victim_case.victim_case_i_code
type CaseTaskWithRelations struct {
	// ── Campos de case_task ───────────────────────────────────────────────────
	ID             string     `gorm:"column:id"               json:"id"`
	Category       string     `gorm:"column:category"         json:"category"`
	Type           string     `gorm:"column:type"             json:"type"`
	Description    string     `gorm:"column:description"      json:"description"`
	Result         *string    `gorm:"column:result"           json:"result,omitempty"`
	CompletedAt    *time.Time `gorm:"column:completed_at"     json:"completedAt,omitempty"`
	AssignedUserID string     `gorm:"column:assigned_user_id" json:"assignedUserId"`
	Status         string     `gorm:"column:status"           json:"status"`
	CaseID         string     `gorm:"column:case_id"          json:"caseId"`
	BarrierID      *string    `gorm:"column:barrier_id"       json:"barrierId,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at"       json:"createdAt"`

	// ── Campos de barrier_v2 ─────────────────────────────────────────────────
	BarrierSector      string `gorm:"column:barrier_sector"      json:"barrierSector"`
	BarrierDescription string `gorm:"column:barrier_description" json:"barrierDescription"`
	BarrierStatus      string `gorm:"column:barrier_status"      json:"barrierStatus"`

	// ── Campos de victim_case ─────────────────────────────────────────────────
	VictimName      string `gorm:"column:victim_name"       json:"victimName"`
	VictimLastName  string `gorm:"column:victim_last_name"  json:"victimLastName"`
	VictimDocNumber string `gorm:"column:victim_doc_number" json:"victimDocNumber"`
	CaseCode        string `gorm:"column:case_code"         json:"caseCode"`
}
