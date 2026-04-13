// Package models contiene los structs GORM para el Motor EAV y Auditoría Forense.
package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CaseTimeline es un registro de auditoría forense append-only.
// NO tiene DeletedAt: los eventos del timeline son inmutables.
type CaseTimeline struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	CaseID    string         `gorm:"type:varchar(36);not null;index"`
	ActorID   string         `gorm:"type:varchar(36);not null"`
	EventType string         `gorm:"type:varchar(50);not null"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time
}

// KoboSubmission representa el envío de un formulario Kobo (integración legacy).
// TimelineID es el ancla forense: vincula el submission al evento de auditoría.
type KoboSubmission struct {
	ID         string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormID     string         `gorm:"type:varchar(36);not null;index"`
	FollowUpID *string        `gorm:"type:varchar(36);index"` // opcional
	AgentID    string         `gorm:"type:varchar(36);not null"`
	TimelineID string         `gorm:"type:varchar(36);not null;index"`
	Timeline   *CaseTimeline  `gorm:"foreignKey:TimelineID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// KoboAnswer almacena la respuesta a una pregunta dentro de un KoboSubmission (integración legacy).
// QuestionSnapshot preserva el texto de la pregunta en el momento del envío.
type KoboAnswer struct {
	ID               string          `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormSubmissionID string          `gorm:"type:varchar(36);not null;index"`
	FormSubmission   *KoboSubmission `gorm:"foreignKey:FormSubmissionID"`
	QuestionID       string          `gorm:"type:varchar(36);not null"`
	Value            string          `gorm:"type:text"`
	QuestionSnapshot string          `gorm:"type:text"` // snapshot del enunciado al momento del envío
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
