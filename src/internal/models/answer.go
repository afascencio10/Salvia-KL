package models

import (
	"time"

	"gorm.io/gorm"
)

// Answer representa la tabla salvia.answer.
type Answer struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormSubmissionID string         `gorm:"type:varchar(36);not null;index"`
	QuestionID       string         `gorm:"type:varchar(36);not null;index"`
	RepeaterEntryID  *string        `gorm:"type:varchar(36);index"` // nullable
	Value            string         `gorm:"type:text"`
	QuestionSnapshot string         `gorm:"type:text;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (Answer) TableName() string { return "salvia.answer" }
