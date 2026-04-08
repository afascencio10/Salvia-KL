package models

import (
	"time"

	"gorm.io/gorm"
)

// Form representa la tabla salvia.form — maestro de formularios dinámicos.
type Form struct {
	ID         string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	Name       string         `gorm:"type:varchar(255);not null"`
	CampaignID string         `gorm:"type:varchar(50)"`
	Status     string         `gorm:"type:varchar(20);default:'ACTIVE'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Form) TableName() string { return "salvia.form" }
