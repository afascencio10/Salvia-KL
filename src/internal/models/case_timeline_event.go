// Package models — case_timeline_event.go
// Modelo para registrar eventos/trazabilidad de un caso en el timeline.
package models

import "time"

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
	EventType   string    `gorm:"type:varchar(30);not null"                             json:"event_type"`
	Description string    `gorm:"type:text"                                             json:"description"`
	ActorID     string    `gorm:"type:varchar(36)"                                      json:"actor_id"`
	ActorName   string    `gorm:"type:varchar(128)"                                     json:"actor_name"`
	CreatedAt   time.Time `                                                             json:"created_at"`
}

func (CaseTimelineEvent) TableName() string { return "salvia.case_timeline_event" }
