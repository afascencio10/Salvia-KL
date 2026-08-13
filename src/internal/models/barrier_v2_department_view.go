package models

import "time"

// BarrierV2DepartmentItem es el DTO de lectura enriquecido para la pantalla
// "Barreras Departamento". A diferencia de CaseTaskWithRelations (Mis
// Barreras), no depende de case_task — trae directamente todas las barreras
// activas de un departamento, sin filtrar por asignación de tareas.
//
// Joins:
//
//	barrier_v2.case_id          → victim_case.victim_case_i_code
//	victim_case.victim_case_id  → victim_case_form2.victim_case_form2_victim_case
//	barrier_v2.entity_branch_id → entity_branch.entity_branch_id
//	barrier_v2.city_id          → security.city.city_id
type BarrierV2DepartmentItem struct {
	// ── Campos de barrier_v2 ─────────────────────────────────────────────────
	ID           string    `gorm:"column:id"            json:"id"`
	CaseID       string    `gorm:"column:case_id"       json:"caseId"`
	Sector       string    `gorm:"column:sector"        json:"sector"`
	Description  string    `gorm:"column:description"   json:"description"`
	Status       string    `gorm:"column:status"        json:"status"`
	BarrierDate  string    `gorm:"column:barrier_date"  json:"barrierDate"`
	CreatedAt    time.Time `gorm:"column:created_at"    json:"createdAt"`

	// ── Entidad involucrada (entity_branch) ─────────────────────────────────
	EntityBranchName string `gorm:"column:entity_branch_name" json:"entityBranchName"`

	// ── Ubicación ────────────────────────────────────────────────────────────
	CityName string `gorm:"column:city_name" json:"cityName"`

	// ── Campos de victim_case ────────────────────────────────────────────────
	VictimName      string `gorm:"column:victim_name"       json:"victimName"`
	VictimLastName  string `gorm:"column:victim_last_name"  json:"victimLastName"`
	VictimDocNumber string `gorm:"column:victim_doc_number" json:"victimDocNumber"`
	CaseCode        string `gorm:"column:case_code"         json:"caseCode"`

	// ── Nivel de prioridad del caso (victim_case_form2_risk_level 1-4) ──────
	RiskLevel *int `gorm:"column:risk_level" json:"riskLevel"`
}
