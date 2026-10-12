package models

import "time"

// EntityCaseFollowUp registra las respuestas del repeater "Seguimiento a
// Entidades" (Sección 5 del formulario de seguimiento v2) para una entidad ya
// vinculada al caso (entity_case). Append-only — sin DeletedAt, igual que
// BarrierFollowUp.
type EntityCaseFollowUp struct {
	ID                      string  `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	EntityCaseID            string  `gorm:"type:varchar(36);not null;index;column:entity_case_id" json:"entityCaseId"`
	FollowUpID              string  `gorm:"type:varchar(36);not null;column:follow_up_id"         json:"followUpId"`
	RutaActualizada         bool    `gorm:"not null;default:false;column:ruta_actualizada"        json:"rutaActualizada"`
	MotivoActualizacion     *string `gorm:"type:text;column:motivo_actualizacion"                 json:"motivoActualizacion,omitempty"`
	RequiereNuevaActivacion bool    `gorm:"not null;default:false;column:requiere_nueva_activacion" json:"requiereNuevaActivacion"`
	CanalActivacion         *string `gorm:"type:text;column:canal_activacion"                     json:"canalActivacion,omitempty"` // CSV
	CreatedByID             string  `gorm:"type:varchar(36);not null;column:created_by_id"        json:"createdById"`
	// TeamContactID — sesión psicosocial (team_contact) que generó este registro.
	// nil cuando se crea desde Hacer Seguimiento (que no tiene team_contact).
	TeamContactID *string `gorm:"type:varchar(36);column:team_contact_id;index" json:"teamContactId,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}

func (EntityCaseFollowUp) TableName() string { return "salvia.entity_case_follow_up" }
