package models

import (
	"time"

	"gorm.io/gorm"
)

// EmergencyMeasure representa la tabla salvia.emergency_measure.
type EmergencyMeasure struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	CaseID           string         `gorm:"type:varchar(36);not null;index"`
	FollowUpID       string         `gorm:"type:varchar(36);not null;index"`
	Type             string         `gorm:"type:varchar(50);not null"`
	IssuingAuthority string         `gorm:"type:varchar(255)"`
	Status           string         `gorm:"type:varchar(20);default:'ACTIVE'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (EmergencyMeasure) TableName() string { return "salvia.emergency_measure" }
