package models

import (
	"time"

	"gorm.io/gorm"
)

// Estados válidos para FollowUpV2.Status
const (
	FollowUpStatusPendiente    = "PENDIENTE"
	FollowUpStatusRealizado    = "REALIZADO"
	FollowUpStatusVencido      = "VENCIDO"
	FollowUpStatusReprogramado = "REPROGRAMADO"
)

// FollowUpV2 representa la tabla salvia.follow_up_v2 (Fase 2).
type FollowUpV2 struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID           string         `gorm:"type:varchar(36);index;not null"                       json:"case_id"`
	FormSubmissionID *string        `gorm:"type:varchar(36)"                                      json:"form_submission_id,omitempty"`
	AgentID          string         `gorm:"type:varchar(36);not null"                             json:"agent_id"`
	SupervisorID     string         `gorm:"type:varchar(36)"                                      json:"supervisor_id"`
	Status           string         `gorm:"type:varchar(20);not null"                             json:"status"`
	Team             string         `gorm:"type:varchar(50)"                                      json:"team"`
	ScheduledDate    time.Time      `gorm:"not null"                                              json:"scheduled_date"`
	ScheduledTime    *string        `gorm:"type:varchar(8)"                                       json:"scheduled_time,omitempty"`
	CompletedAt      *time.Time     `                                                             json:"completed_at,omitempty"`
	IsPriority       bool           `gorm:"default:false"                                         json:"is_priority"`
	Summary          *string        `gorm:"type:text"                                             json:"summary,omitempty"`
	RiskStatus       *string        `gorm:"type:varchar(20)"                                      json:"risk_status,omitempty"`
	CreatedAt        time.Time      `                                                             json:"created_at"`
	UpdatedAt        time.Time      `                                                             json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index"                                                 json:"-"`
}

func (FollowUpV2) TableName() string { return "salvia.follow_up_v2" }
