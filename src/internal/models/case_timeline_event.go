// Package models — case_timeline_event.go
// Modelo para registrar eventos/trazabilidad de un caso en el timeline.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Tipos de evento para el timeline
const (
	TimelineEventRegistro       = "REGISTRO"
	TimelineEventSeguimiento    = "SEGUIMIENTO_CREADO"
	TimelineEventReasignacion   = "REASIGNACION"
	TimelineEventReasignSeg     = "REASIGNACION_SEGUIMIENTO"
	TimelineEventEstadoCambio   = "CAMBIO_ESTADO"
	TimelineEventBarrera        = "BARRERA_IDENTIFICADA"
	TimelineEventNota           = "NOTA"
)

// CaseTimelineEvent registra un evento en la vida de un caso.
// Es append-only — los eventos no se modifican ni eliminan.
type CaseTimelineEvent struct {
	ID          string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID      string    `gorm:"type:varchar(36);index;not null"                       json:"case_id"`
	EventType   string    `gorm:"type:varchar(30);not null"                             json:"event_type"` // Deprecated in favor of Type?
	Description string    `gorm:"type:text"                                             json:"description"`
	ActorID     string    `gorm:"type:varchar(36)"                                      json:"actor_id"` // Deprecated in favor of EventUserID?
	ActorName   string    `gorm:"type:varchar(128)"                                     json:"actor_name"`
	
	// Nuevos campos solicitados
	Category    string    `gorm:"type:varchar(50)"                                      json:"category"`
	Type        string    `gorm:"type:varchar(50)"                                      json:"type"`
	Icon        string    `gorm:"type:varchar(100)"                                     json:"icon"`
	Date        time.Time `gorm:"type:timestamp with time zone"                          json:"date"`
	EventUserID string    `gorm:"type:varchar(36)"                                      json:"event_user_id"`
	Color       string    `gorm:"type:varchar(20)"                                      json:"color"`

	// Relaciones
	TaskID                  string `gorm:"type:varchar(36)" json:"task_id"`
	FollowUpID              string `gorm:"type:varchar(36)" json:"follow_up_id"`
	EntityLetterID          string `gorm:"type:varchar(36)" json:"entity_letter_id"`
	BarrierID               string `gorm:"type:varchar(36)" json:"barrier_id"`
	EmergencyMeasureID      string `gorm:"type:varchar(36)" json:"emergency_measure_id"`
	PsychosocialSupportID   string `gorm:"type:varchar(36)" json:"psychosocial_support_id"`
	EconomicStabilizationID string `gorm:"type:varchar(36)" json:"economic_stabilization_id"`

	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CaseTimelineEvent) TableName() string { return "salvia.case_timeline_event" }
