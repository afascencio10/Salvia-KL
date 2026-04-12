package models

import (
	"time"

	"gorm.io/gorm"
)

// EmergencyMeasure modelo GORM que mapea la tabla salvia.emergency_measure.
type EmergencyMeasure struct {
	ID               string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID           string         `gorm:"type:varchar(36);not null" json:"caseId"`
	FollowUpID       string         `gorm:"type:uuid;not null" json:"followUpId"`
	Type             string         `gorm:"type:varchar(50);not null" json:"type"`
	IssuingAuthority string         `gorm:"type:varchar(255)" json:"issuingAuthority"`
	IssuedAt         *time.Time     `json:"issuedAt"`
	ExpiresAt        *time.Time     `json:"expiresAt"`
	Status           string         `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	Notes            *string        `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmergencyMeasure) TableName() string {
	return "salvia.emergency_measure"
}
