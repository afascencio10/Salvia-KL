// Package models — case_timeline_event.go
// Modelo para registrar eventos/trazabilidad de un caso en el timeline.
package models

import (
	"time"

	"gorm.io/gorm"
)

// ─── Categorías del timeline ─────────────────────────────────────────────────
const (
	TimelineCategoryGeneral        = "General"
	TimelineCategoryBarreras       = "Barreras"
	TimelineCategorySeguimientos   = "Seguimientos"
	TimelineCategoryOficios        = "Oficios"
	TimelineCategoryMedidas        = "Medidas"
	TimelineCategoryPsicosocial    = "Psicosocial"
	TimelineCategoryEstabilizacion = "Estabilización"
)

// ─── Tipos de evento (Type) ──────────────────────────────────────────────────
const (
	TimelineTypeCreacionCaso          = "Creación de Caso"
	TimelineTypeIntentoSeguimiento    = "Intento de Seguimiento"
	TimelineTypeSeguimientoPospuesto  = "Seguimiento Pospuesto"
	TimelineTypeReasignacionCaso      = "Reasignación de Caso"
	TimelineTypeReasignacionSeg       = "Reasignación de Seguimiento"
	TimelineTypeSeguimientoProgramado = "Seguimiento Programado"
	TimelineTypeSeguimientoEjecutado  = "Seguimiento Ejecutado"
	TimelineTypeSeguimientoEditado    = "Seguimiento Editado"
	TimelineTypeCierreCaso            = "Cierre de Caso"
	TimelineTypeCambioEstado          = "Cambio de Estado"
	TimelineTypeBarreraIdentificada   = "Barrera Identificada"
	TimelineTypeBarreraArticulada     = "Barrera Articulada"
	TimelineTypeNota                  = "Nota"

	// Oficios (EntityLetter) — un Type por estado destino
	TimelineTypeOficioParaRevisar        = "Oficio para revisar"
	TimelineTypeOficioEnCorreccion       = "Oficio en corrección"
	TimelineTypeOficioAprobacionJuridica = "Oficio en aprobación jurídica"
	TimelineTypeOficioParaRadicar        = "Oficio para radicar"
	TimelineTypeOficioRadicado           = "Oficio radicado"
	TimelineTypeOficioRespondido         = "Oficio respondido"
)

// ─── Iconos ──────────────────────────────────────────────────────────────────
const (
	TimelineIconRegistro     = "clipboard-list"
	TimelineIconSeguimiento  = "calendar-check"
	TimelineIconReasignacion = "arrows-rotate"
	TimelineIconPospuesto    = "calendar-days"
	TimelineIconEstado       = "shuffle"
	TimelineIconBarrera          = "triangle-exclamation"
	TimelineIconBarreraArticulada = "circle-check"
	TimelineIconNota         = "note-sticky"
	TimelineIconCierre       = "circle-xmark"

	// Oficios
	TimelineIconOficioRevisar    = "eye"
	TimelineIconOficioCorreccion = "triangle-exclamation"
	TimelineIconOficioJuridica   = "scale-balanced"
	TimelineIconOficioRadicar    = "paper-plane"
	TimelineIconOficioRespondido = "envelope-open"
)

// ─── Colores ─────────────────────────────────────────────────────────────────
const (
	TimelineColorPurple = "#7c3aed"
	TimelineColorGreen  = "#22c55e"
	TimelineColorBlue   = "#3b82f6"
	TimelineColorOrange = "#f97316"
	TimelineColorGray   = "#6b7280"
	TimelineColorTeal   = "#63e6be"
	TimelineColorYellow = "#f8a625"
	TimelineColorRed    = "#dc2626"
)

// ─── Constantes legacy (compatibilidad con eventos existentes en BD) ─────────
const (
	TimelineEventRegistro        = "REGISTRO"
	TimelineEventSeguimiento     = "SEGUIMIENTO_CREADO"
	TimelineEventReasignacion    = "REASIGNACION"
	TimelineEventReasignSeg      = "REASIGNACION_SEGUIMIENTO"
	TimelineEventEstadoCambio    = "CAMBIO_ESTADO"
	TimelineEventBarrera           = "BARRERA_IDENTIFICADA"
	TimelineEventBarreraArticulada = "BARRERA_ARTICULADA"
	TimelineEventNota              = "NOTA"
	TimelineEventOficioActualizado = "OFICIO_ACTUALIZADO"
)

// CaseTimelineEvent registra un evento en la vida de un caso.
// Es append-only — los eventos no se modifican ni eliminan.
type CaseTimelineEvent struct {
	ID          string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID      string    `gorm:"type:varchar(36);index;not null"                       json:"case_id"`
	EventType   string    `gorm:"type:varchar(30)"                                      json:"event_type"` // Legacy — se mantiene para compatibilidad
	Description string    `gorm:"type:text"                                             json:"description"`
	ActorID     string    `gorm:"type:varchar(36)"                                      json:"actor_id"`
	ActorName   string    `gorm:"type:varchar(128)"                                     json:"actor_name"`

	// Campos del esquema nuevo
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
