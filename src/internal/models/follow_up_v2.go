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

type FollowUpV2 struct {
	// ── Campos originales ───────────────────────────────────────────────────
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID    string         `gorm:"type:varchar(36);index;not null"                       json:"case_id"`
	AgentID   string         `gorm:"type:varchar(36);not null"                             json:"agent_id"`
	Status    string         `gorm:"type:varchar(20);not null"                             json:"status"`
	CreatedAt time.Time      `                                                             json:"created_at"`
	UpdatedAt time.Time      `                                                             json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                                 json:"-"`

	// ── Campos de gestión y calendario ──────────────────────────────────────
	FormSubmissionID   *string    `gorm:"type:varchar(36)"          json:"form_submission_id,omitempty"`
	SupervisorID       string     `gorm:"type:varchar(36)"          json:"supervisor_id"`
	Team               string     `gorm:"type:varchar(50)"          json:"team"`
	ScheduledDate      time.Time  `gorm:"type:date;not null"        json:"scheduled_date"`
	ScheduledTime      string     `gorm:"type:time"                 json:"scheduled_time"`
	IsCompleted        bool       `gorm:"default:false"             json:"is_completed"`
	CompletedAt        *time.Time `                                 json:"completed_at,omitempty"`
	Attempts           int        `gorm:"default:0"                 json:"attempts"`
	IsPriority         bool       `gorm:"default:false"             json:"is_priority"`
	ReassignmentReason *string    `gorm:"type:text"                 json:"reassignment_reason,omitempty"`
	Summary            *string    `gorm:"type:text"                 json:"summary,omitempty"`
	RiskStatus         string     `gorm:"type:varchar(20)"          json:"risk_status"`
	SequenceNumber     int        `gorm:"default:0"                 json:"sequence_number"`
}

func (FollowUpV2) TableName() string { return "salvia.follow_up_v2" }
