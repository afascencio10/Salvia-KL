package models

import (
	"time"

	"gorm.io/gorm"
)

// FormSubmission representa la tabla salvia.form_submission.
type FormSubmission struct {
	ID         string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormID     string         `gorm:"type:varchar(36);not null;index"`
	FollowUpID *string        `gorm:"type:varchar(36);index"` // nullable
	AgentID    string         `gorm:"type:varchar(36);not null"`
	ScoreTotal *float64       `gorm:"type:numeric(5,2)"` // nullable
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (FormSubmission) TableName() string { return "salvia.form_submission" }
