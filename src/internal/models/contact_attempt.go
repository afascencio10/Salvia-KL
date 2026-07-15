package models

import (
	"time"

	"gorm.io/gorm"
)

// ContactAttempt representa un intento de contacto (exitoso o fallido) del flujo
// 3x3 de Atención Psicosocial. Guarda TODOS los intentos de un proceso
// psychosocial_support; es independiente de follow_up_attempts (Seguimientos).
type ContactAttempt struct {
	ID             string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PsicosocialID  string         `gorm:"type:varchar(36);not null;index;column:psicosocial_id" json:"psicosocialId"`
	CaseID         string         `gorm:"type:varchar(36);not null;column:case_id" json:"caseId"`
	WasAnswered    bool           `gorm:"default:false;column:was_answered" json:"wasAnswered"`
	Note           *string        `gorm:"type:text" json:"note,omitempty"`
	ConsentGiven   *bool          `gorm:"column:consent_given" json:"consentGiven,omitempty"`
	ProfessionalID *string        `gorm:"type:varchar(36);column:professional_id" json:"professionalId,omitempty"`
	Team           *string        `gorm:"type:varchar(50)" json:"team,omitempty"`
	AttemptAt      time.Time      `gorm:"column:attempt_at;default:now()" json:"attemptAt"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ContactAttempt) TableName() string { return "salvia.contact_attempts" }
