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
//	    is_psico_session    BOOLEAN DEFAULT false
//	);
type TeamContact struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID           string         `gorm:"type:varchar(36);not null;column:case_id" json:"caseId"`
	FormSubmissionID *string        `gorm:"type:varchar(36);column:form_submission_id" json:"formSubmissionId"`
	DuplaID          *string        `gorm:"type:varchar(36);column:dupla_id" json:"duplaId"`
	PsicosocialID    *string        `gorm:"type:varchar(36);column:psicosocial_id" json:"psicosocialId"`
	ProfessionalID   *string        `gorm:"type:varchar(36);column:professional_id" json:"professionalId"`
	Team             *string        `gorm:"type:varchar(50)" json:"team"`
	ScheduledDate    *time.Time     `gorm:"column:scheduled_date" json:"scheduledDate"`
	IsCompleted      bool           `gorm:"column:is_completed;default:false" json:"isCompleted"`
	Status           *string        `gorm:"type:varchar(20)" json:"status"`
	ScheduledTime    *string        `gorm:"type:varchar(8);column:scheduled_time" json:"scheduledTime"`
	CompletedAt      *time.Time     `gorm:"column:completed_at" json:"completedAt"`
	Summary          *string        `gorm:"type:text" json:"summary"`
	IsPsicoSession   bool           `gorm:"column:is_psico_session;default:false" json:"isPsicoSession"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (TeamContact) TableName() string { return "salvia.team_contact" }
