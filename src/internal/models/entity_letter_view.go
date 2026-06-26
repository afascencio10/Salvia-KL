package models

import "time"

// EntityLetterWithRelations es el DTO de lectura enriquecido que combina:
//   - Todos los campos de entity_letter
//   - barrier_sector y barrier_description de salvia.barrier_v2
//   - victim_name, victim_last_name, victim_doc_number y case_code de salvia.victim_case
//
// El join se realiza así:
//
//	entity_letter.barrier_id  → barrier_v2.id     (UUID almacenado como text)
//	entity_letter.case_id     → victim_case.victim_case_i_code
type EntityLetterWithRelations struct {
	// ── Campos de entity_letter ───────────────────────────────────────────────
	ID                 string    `gorm:"column:id"                   json:"id"`
	BarrierID          string    `gorm:"column:barrier_id"           json:"barrierId"`
	CaseID             string    `gorm:"column:case_id"              json:"caseId"`
	State              string    `gorm:"column:state"                json:"state"`
	Priority           string    `gorm:"column:priority"             json:"priority"`
	AgentID            *string   `gorm:"column:agent_id"             json:"agentId,omitempty"`
	NotificationUserID *string   `gorm:"column:notification_user_id" json:"notificationUserId,omitempty"`
	ReviewBy           *string   `gorm:"column:review_by"            json:"reviewBy,omitempty"`
	RadicadoBy         *string   `gorm:"column:radicado_by"          json:"radicadoBy,omitempty"`
	RegisterBy         *string   `gorm:"column:register_by"          json:"registerBy,omitempty"`
	Entidad            *string    `gorm:"column:entidad"              json:"entidad,omitempty"`
	Nivel              *string    `gorm:"column:nivel"                json:"nivel,omitempty"`
	UrlKofax           *string    `gorm:"column:url_kofax"            json:"urlKofax,omitempty"`
	EntityBranchID     *int64     `gorm:"column:entity_branch_id"     json:"entityBranchId,omitempty"`
	DepartmentID       *string    `gorm:"column:department_id"        json:"departmentId,omitempty"`
	CityID             *string    `gorm:"column:city_id"              json:"cityId,omitempty"`
	TownID             *string    `gorm:"column:town_id"              json:"townId,omitempty"`
	OfficialDependency *string    `gorm:"column:official_dependency"  json:"officialDependency,omitempty"`
	Subject            *string    `gorm:"column:subject"              json:"subject,omitempty"`
	TownName           string     `gorm:"column:town_name"            json:"townName"`
	AsuntoRadicado     *string    `gorm:"column:asunto_radicado"      json:"asuntoRadicado,omitempty"`
	CorreoEntidad      *string    `gorm:"column:correo_entidad"       json:"correoEntidad,omitempty"`
	NumeroRadicado     *string    `gorm:"column:numero_radicado"      json:"numeroRadicado,omitempty"`
	ResponseDate       *time.Time `gorm:"column:response_date"        json:"responseDate,omitempty"`
	CorreoRemitente    *string    `gorm:"column:correo_remitente"     json:"correoRemitente,omitempty"`
	AsuntoRespuesta    *string    `gorm:"column:asunto_respuesta"     json:"asuntoRespuesta,omitempty"`
	ResponseReviewBy   *string    `gorm:"column:response_review_by"   json:"responseReviewBy,omitempty"`
	ReasonCorrection   *string    `gorm:"column:reason_correction"    json:"reasonCorrection,omitempty"`
	CreatedAt          time.Time  `gorm:"column:created_at"           json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"           json:"updatedAt"`

	// ── Campos de barrier_v2 ─────────────────────────────────────────────────
	BarrierSector      string `gorm:"column:barrier_sector"       json:"barrierSector"`
	BarrierDescription string `gorm:"column:barrier_description"  json:"barrierDescription"`

	// ── Campos de victim_case ─────────────────────────────────────────────────
	VictimName      string `gorm:"column:victim_name"       json:"victimName"`
	VictimLastName  string `gorm:"column:victim_last_name"   json:"victimLastName"`
	VictimDocNumber string `gorm:"column:victim_doc_number"  json:"victimDocNumber"`
	CaseCode        string `gorm:"column:case_code"          json:"caseCode"`
}
