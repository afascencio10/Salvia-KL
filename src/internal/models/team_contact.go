package models

import (
	"time"

	"gorm.io/gorm"
)

// TeamContact registra contactos/sesiones del equipo de Atención Psicosocial.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.team_contact (
//	    id                  VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    case_id             VARCHAR(36) NOT NULL,
//	    form_submission_id  VARCHAR(36),
//	    dupla_id            VARCHAR(36),
//	    psicosocial_id      VARCHAR(36),
//	    professional_id     VARCHAR(36),
//	    team                VARCHAR(50),
//	    scheduled_date      TIMESTAMPTZ,
//	    is_completed        BOOLEAN DEFAULT false,
//	    created_at          TIMESTAMPTZ,
//	    updated_at          TIMESTAMPTZ,
//	    deleted_at          TIMESTAMPTZ,
//	    status              VARCHAR(20),
//	    scheduled_time      VARCHAR(8),
//	    completed_at        TIMESTAMP,
//	    summary             TEXT,
//	    is_psico_session    BOOLEAN DEFAULT false,
//	    form_id             VARCHAR(36),
//	    session_type        VARCHAR(30)
//	);
type TeamContact struct {
	ID               string  `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID           string  `gorm:"type:varchar(36);not null;column:case_id" json:"caseId"`
	FormSubmissionID *string `gorm:"type:varchar(36);column:form_submission_id" json:"formSubmissionId"`
	DuplaID          *string `gorm:"type:varchar(36);column:dupla_id" json:"duplaId"`
	PsicosocialID    *string `gorm:"type:varchar(36);column:psicosocial_id" json:"psicosocialId"`
	ProfessionalID   *string `gorm:"type:varchar(36);column:professional_id" json:"professionalId"`
	// FormID fija el formulario psicosocial (salvia.form.id) con el que se inició esta
	// sesión. Se resuelve una sola vez al cargar la pantalla (evento E-01) y se reutiliza
	// en cargas posteriores para que la sesión no "salte" de formulario si el estado del
	// proceso cambia mientras el team_contact sigue pendiente.
	FormID *string `gorm:"type:varchar(36);column:form_id" json:"formId,omitempty"`
	// SessionType distingue el tipo de sesión psicosocial registrada.
	// Valores: PRIMER_CONTACTO / PRIMERA_ATENCION / ATENCION_PSICOSOCIAL / CIERRE.
	SessionType    *string        `gorm:"type:varchar(30);column:session_type" json:"sessionType,omitempty"`
	Team           *string        `gorm:"type:varchar(50)" json:"team"`
	ScheduledDate  *time.Time     `gorm:"column:scheduled_date" json:"scheduledDate"`
	IsCompleted    bool           `gorm:"column:is_completed;default:false" json:"isCompleted"`
	Status         *string        `gorm:"type:varchar(20)" json:"status"`
	ScheduledTime  *string        `gorm:"type:varchar(8);column:scheduled_time" json:"scheduledTime"`
	CompletedAt    *time.Time     `gorm:"column:completed_at" json:"completedAt"`
	Summary        *string        `gorm:"type:text" json:"summary"`
	IsPsicoSession bool           `gorm:"column:is_psico_session;default:false" json:"isPsicoSession"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (TeamContact) TableName() string { return "salvia.team_contact" }

// Valores válidos de SessionType para sesiones psicosociales.
const (
	SessionTypePrimerContacto      = "PRIMER_CONTACTO"
	SessionTypePrimeraAtencion     = "PRIMERA_ATENCION"
	SessionTypeAtencionPsicosocial = "ATENCION_PSICOSOCIAL"
	SessionTypeCierre              = "CIERRE"

	// SessionTypePrimerContactoConAtencion: Form Primer Contacto con "Continuar Primera
	// Atención = Sí" y consentimiento = Sí en la Sección 4 (la Primera Atención se completó
	// en la misma sesión, sin pasar por el Form 2).
	SessionTypePrimerContactoConAtencion = "PRIMER_CONTACTO_CON_ATENCION"
	// SessionTypePrimerContactoSinConsentimiento: Form Primer Contacto con "Continuar = Sí"
	// pero Consentimiento Informado = No en la Sección 4. Se registra distinto de
	// PRIMER_CONTACTO solo para reporting; el efecto en psychosocial_support es el mismo
	// (avanza como un primer contacto regular — decisión del líder, Jul 2026).
	SessionTypePrimerContactoSinConsentimiento = "PRIMER_CONTACTO_SIN_CONSENTIMIENTO"
	// SessionTypeContactoSinAtencion: "¿Es atención o solo contacto?" = Solo Contacto
	// (Form Primera Atención / Atención Psicosocial / Cierre). No cuenta como sesión.
	SessionTypeContactoSinAtencion = "CONTACTO_SIN_ATENCION"
	// SessionTypeCierreNoConsentimiento: Consentimiento Informado = No en Form Primera
	// Atención (o en Primer Contacto con Continuar=Sí). status pasa a en_devolucion.
	SessionTypeCierreNoConsentimiento = "CIERRE_NO_CONSENTIMIENTO"
)
