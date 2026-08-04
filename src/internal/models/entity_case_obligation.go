package models

import "time"

// EntityCaseObligation registra una obligación marcada (completada o pendiente)
// para un entity_case concreto, durante un seguimiento. Append-only — sin
// DeletedAt, igual que BarrierFollowUp: cada seguimiento agrega su propia foto
// del estado del checklist, no se sobreescribe.
type EntityCaseObligation struct {
	ID                 string  `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	EntityCaseID       string  `gorm:"type:varchar(36);not null;index;column:entity_case_id" json:"entityCaseId"`
	EntityObligationID *string `gorm:"type:varchar(36);column:entity_obligation_id"          json:"entityObligationId,omitempty"` // nil si es "Otra"
	CustomLabel        *string `gorm:"type:text;column:custom_label"                         json:"customLabel,omitempty"`        // solo si es "Otra"
	Status             string  `gorm:"type:varchar(20);not null"                             json:"status"`                       // "completada" | "pendiente"
	// FollowUpID — nil cuando se crea desde Registro de Caso (el caso recién se
	// está creando, no existe ningún follow_up_v2 todavía).
	FollowUpID  *string `gorm:"type:varchar(36);column:follow_up_id" json:"followUpId,omitempty"`
	CreatedByID string  `gorm:"type:varchar(36);not null;column:created_by_id" json:"createdById"`
	// TeamContactID — sesión psicosocial (team_contact) que generó este registro.
	// nil cuando se crea desde Hacer Seguimiento o Registro de Caso.
	TeamContactID *string `gorm:"type:varchar(36);column:team_contact_id;index" json:"teamContactId,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}

func (EntityCaseObligation) TableName() string { return "salvia.entity_case_obligation" }

// Valores válidos de Status.
const (
	EntityCaseObligationStatusCompletada = "completada"
	EntityCaseObligationStatusPendiente  = "pendiente"
)
