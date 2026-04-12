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
	LastAttemptAt      *time.Time `                                 json:"last_attempt_at,omitempty"`

	// ── Campos enriquecidos (no persistidos) ────────────────────────────────
	AgentNames     string `gorm:"-" json:"agent_names,omitempty"`
	AgentLastNames string `gorm:"-" json:"agent_last_names,omitempty"`
}

func (FollowUpV2) TableName() string { return "salvia.follow_up_v2" }

// MyDayFollowUpResponse representa un seguimiento enriquecido con datos del caso para la vista "Mis Seguimientos"
type MyDayFollowUpResponse struct {
	ID               string            `json:"id"`
	CaseID           string            `json:"case_id"`
	Case             *VictimCaseLight  `json:"case,omitempty"`
	RiskStatus       string            `json:"risk_status"`
	ScheduledTime    string            `json:"scheduled_time"`
	Attempts         int               `json:"attempts"`
	IsPriority       bool              `json:"is_priority"`
	Status           string            `json:"status"`
	SequenceNumber   int               `json:"sequence_number"`
	LastAttemptAt    *time.Time        `json:"last_attempt_at,omitempty"`
	FollowUpAttempts []FollowUpAttempt `json:"follow_up_attempts,omitempty"`
}

// MyDayResponse representa la respuesta completa del endpoint /follow-ups/my-day
type MyDayResponse struct {
	Success string            `json:"success"`
	Data    MyDayDataResponse `json:"data"`
}

type MyDayDataResponse struct {
	FollowUpsPendingCount         int                     `json:"followUpsPendingCount"`
	FollowUpsPendingPriorityCount int                     `json:"followUpsPendingPriorityCount"`
	FollowUpsCompletedCount       int                     `json:"followUpsCompletedCount"`
	FollowUpsPending              []MyDayFollowUpResponse `json:"followUpsPending"`
	FollowUpsPendingPriority      []MyDayFollowUpResponse `json:"followUpsPendingPriority"`
	FollowUpsCompleted            []MyDayFollowUpResponse `json:"followUpsCompleted"`
}
