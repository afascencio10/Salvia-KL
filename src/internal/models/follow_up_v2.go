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
	FollowUpStatusCerrado      = "CERRADO"
)

// FollowUpV2 representa la tabla salvia.follow_up_v2 (Fase 2).
type FollowUpV2 struct {
	// ── Identificadores ─────────────────────────────────────────────────────
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID           string         `gorm:"type:varchar(36);index;not null"                       json:"case_id"`
	FormSubmissionID *string        `gorm:"type:varchar(36)"                                      json:"form_submission_id,omitempty"`
	FormID           *string        `gorm:"type:varchar(36)"                                      json:"form_id,omitempty"`
	AgentID          *string        `gorm:"type:varchar(36)"                                      json:"agent_id,omitempty"`
	SupervisorID     string         `gorm:"type:varchar(36)"                                      json:"supervisor_id"`

	// ── Estado y equipo ──────────────────────────────────────────────────────
	Status           string         `gorm:"type:varchar(20);not null"  json:"status"`
	Team             string         `gorm:"type:varchar(50)"           json:"team"`
	RiskStatus       *string        `gorm:"type:varchar(20)"           json:"risk_status,omitempty"`

	// ── Calendario ───────────────────────────────────────────────────────────
	ScheduledDate      time.Time  `gorm:"not null"           json:"scheduled_date"`
	ScheduledTime      string     `gorm:"type:varchar(8)"    json:"scheduled_time"`
	CompletedAt        *time.Time `                          json:"completed_at,omitempty"`
	SequenceNumber     int        `gorm:"default:0"          json:"sequence_number"`

	// ── Intentos y seguimiento ───────────────────────────────────────────────
	Attempts           int        `gorm:"default:0"          json:"attempts"`
	LastAttemptAt      *time.Time `                          json:"last_attempt_at,omitempty"`
	ReassignmentReason *string    `gorm:"type:text"          json:"reassignment_reason,omitempty"`

	// ── Flags y resumen ──────────────────────────────────────────────────────
	IsPriority       bool    `gorm:"default:false"      json:"is_priority"`
	Summary          *string `gorm:"type:text"          json:"summary,omitempty"`
	// IDs de barreras activas al momento de cargar el seguimiento por primera vez.
	// Comma-separated UUIDs: "id1,id2,id3". Se fija en E-01 y no cambia después.
	ActiveBarrierIDs *string `gorm:"type:text"          json:"active_barrier_ids,omitempty"`

	// ── Auditoría ────────────────────────────────────────────────────────────
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// ── Campos virtuales (no persistidos) ────────────────────────────────────
	AgentNames       string           `gorm:"-" json:"agent_names,omitempty"`
	AgentLastNames   string           `gorm:"-" json:"agent_last_names,omitempty"`
	FollowUpAttempts []FollowUpAttempt `gorm:"-" json:"follow_up_attempts,omitempty"`
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
	ReassignmentReason *string         `json:"reassignment_reason,omitempty"`
	Team             string            `json:"team"`
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
