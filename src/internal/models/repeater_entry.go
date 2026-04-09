package models

import "time"

// RepeaterEntry representa la tabla salvia.repeater_entry.
// No tiene DeletedAt: es un registro de iteración inmutable.
type RepeaterEntry struct {
	ID               string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormSubmissionID string    `gorm:"type:varchar(36);not null;index"`
	RepeaterGroupID  string    `gorm:"type:varchar(36);not null;index"`
	Iteration        int       `gorm:"not null;default:1"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (RepeaterEntry) TableName() string { return "salvia.repeater_entry" }
