package models

import (
	"time"

	"gorm.io/gorm"
)

// ─── FormSubmission ───────────────────────────────────────────────────────────

type FormSubmission struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormID    string         `gorm:"type:varchar(36);not null;index"                       json:"formId"`
	CreatedAt time.Time      `                                                             json:"createdAt"`
	UpdatedAt time.Time      `                                                             json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (FormSubmission) TableName() string { return "salvia.form_submission" }

// RepeaterEntry y Answer fueron movidos a archivos individuales:
// repeater_entry.go, answer.go
