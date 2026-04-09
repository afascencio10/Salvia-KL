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
// Los campos originales (ID, CaseID, AgentID, Status, CreatedAt, UpdatedAt, DeletedAt)
// se conservan intactos. Se agregan los campos del calendario (HU-027).
type FollowUpV2 struct {
	// ── Campos originales (NO modificar) ─────────────────────────────────────
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID    string         `gorm:"type:varchar(36);index;not null"                       json:"case_id"`
	AgentID   string         `gorm:"type:varchar(36);not null"                             json:"agent_id"`
	Status    string         `gorm:"type:varchar(20);not null"                             json:"status"`
	CreatedAt time.Time      `                                                             json:"created_at"`
	UpdatedAt time.Time      `                                                             json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                                 json:"-"`

	// ── Campos nuevos HU-027 (calendario) ────────────────────────────────────
	FormSubmissionID *string   `gorm:"type:varchar(36)"          json:"form_submission_id,omitempty"`
	Team             string    `gorm:"type:varchar(50)"          json:"team"`
	RiskStatus       string    `gorm:"type:varchar(20)"          json:"risk_status"`
	ScheduledDate    time.Time `gorm:"not null"                  json:"scheduled_date"`
	IsCompleted      bool      `gorm:"default:false"             json:"is_completed"`
	SequenceNumber   int       `gorm:"default:0"                 json:"sequence_number"`
}

func (FollowUpV2) TableName() string { return "salvia.follow_up_v2" }
