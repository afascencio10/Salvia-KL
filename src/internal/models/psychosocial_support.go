package models

import (
	"time"

	"gorm.io/gorm"
)

// PsychosocialSupport modelo GORM que mapea la tabla salvia.psychosocial_support.
//
// Flujo de estados:
//
//	abierto       → al crear la remisión
//	en_gestion    → se logró primer contacto
//	en_devolucion → no cumplió criterios
//	cerrado       → por cualquier motivo
type PsychosocialSupport struct {
	ID              string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID          string     `gorm:"type:varchar(36);not null" json:"caseId"`
	FollowUpID      string     `gorm:"type:uuid;not null" json:"followUpId"`
	Type            string     `gorm:"type:varchar(50);not null" json:"type"`
	Provider        string     `gorm:"type:varchar(255)" json:"provider"`
	ScheduledAt     *time.Time `json:"scheduledAt"`
	SessionCount    int        `gorm:"default:0" json:"sessionCount"`
	Status          string     `gorm:"type:varchar(30);default:'abierto'" json:"status"`
    // NextContactAttemptAt: fecha/hora del próximo intento de contacto (acción 3a del flujo 3x3).
	NextContactAttemptAt *time.Time `gorm:"column:next_contact_attempt_at" json:"nextContactAttemptAt,omitempty"`
	Notes           *string    `gorm:"type:text" json:"notes"`
	SubmittedBy     *string    `gorm:"type:varchar(36);column:submitted_by" json:"submittedBy,omitempty"`
	SubmittedByTeam *string    `gorm:"type:varchar(50);column:submitted_by_team" json:"submittedByTeam,omitempty"`
	DuplaID         *string    `gorm:"type:varchar(36);column:dupla_id" json:"duplaId,omitempty"`
	ProfessionalID  *string    `gorm:"type:varchar(36);column:professional_id" json:"professionalId,omitempty"`
	// YaHizoPrimerContacto se activa al completar el formulario "Primer Contacto".
	// Determina si la pantalla de sesión carga el Form 1 o uno de los siguientes.
	YaHizoPrimerContacto bool `gorm:"column:ya_hizo_primer_contacto;default:false" json:"yaHizoPrimerContacto"`
	// YaHizoPrimeraAtencion se activa al completar el formulario "Primera Atención" con consentimiento.
	// Determina si la pantalla de sesión carga el Form 2 o Form 3/4 (Atención Psicosocial / Cierre).
	YaHizoPrimeraAtencion bool           `gorm:"column:ya_hizo_primera_atencion;default:false" json:"yaHizoPrimeraAtencion"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PsychosocialSupport) TableName() string {
	return "salvia.psychosocial_support"
}

// Estados válidos de una remisión de Atención Psicosocial.
const (
	PsychosocialSupportStatusAbierto      = "abierto"
	PsychosocialSupportStatusEnGestion    = "en_gestion"
	PsychosocialSupportStatusEnDevolucion = "en_devolucion"
	PsychosocialSupportStatusCerrado      = "cerrado"
)
