package models

import (
	"time"

	"gorm.io/gorm"
)

// FormSection representa la tabla salvia.form_section.
type FormSection struct {
	ID         string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormID     string         `gorm:"type:varchar(36);not null;index"`
	Name       string         `gorm:"type:varchar(255);not null"`
	OrderIndex int            `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (FormSection) TableName() string { return "salvia.form_section" }
