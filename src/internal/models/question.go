package models

import (
	"time"

	"gorm.io/gorm"
)

// Question representa la tabla salvia.question.
type Question struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormID          string         `gorm:"type:varchar(36);not null;index"`
	FormSectionID   string         `gorm:"type:varchar(36);not null;index"`
	RepeaterGroupID *string        `gorm:"type:varchar(36);index"` // nullable
	QuestionTypeID  string         `gorm:"type:varchar(50);not null"`
	Description     string         `gorm:"type:text;not null"`
	OrderIndex      int            `gorm:"default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (Question) TableName() string { return "salvia.question" }
