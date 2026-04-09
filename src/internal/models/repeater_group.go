package models

import (
	"time"

	"gorm.io/gorm"
)

// RepeaterGroup representa la tabla salvia.repeater_group.
type RepeaterGroup struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormSectionID   string         `gorm:"type:varchar(36);not null;index"`
	Name            string         `gorm:"type:varchar(255);not null"`
	MinRepetitions  int            `gorm:"default:0"`
	MaxRepetitions  *int           // nullable
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (RepeaterGroup) TableName() string { return "salvia.repeater_group" }
