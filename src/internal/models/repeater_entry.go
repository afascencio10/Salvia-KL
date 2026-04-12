package models

import "time"

// RepeaterEntry representa la tabla salvia.repeater_entry.
// No tiene DeletedAt: es un registro de iteración inmutable.
type RepeaterEntry struct {
	ID               string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormSubmissionID string    `gorm:"type:varchar(36);not null;index"                       json:"formSubmissionId"`
	RepeaterGroupID  string    `gorm:"type:varchar(36);not null;index"                       json:"repeaterGroupId"`
	Iteration        int       `gorm:"not null;default:1"                                    json:"iteration"`
	CreatedAt        time.Time `                                                             json:"createdAt"`
	UpdatedAt        time.Time `                                                             json:"updatedAt"`
}

func (RepeaterEntry) TableName() string { return "salvia.repeater_entry" }
