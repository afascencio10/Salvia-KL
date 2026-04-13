package models

import (
	"time"

	"gorm.io/gorm"
)

// EconomicStabilization modelo GORM que mapea la tabla salvia.economic_stabilization.
type EconomicStabilization struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID       string         `gorm:"type:varchar(36);not null" json:"caseId"`
	FollowUpID   string         `gorm:"type:uuid;not null" json:"followUpId"`
	Type         string         `gorm:"type:varchar(50);not null" json:"type"`
	Institution  string         `gorm:"type:varchar(255)" json:"institution"`
	Benefit      *string        `gorm:"type:varchar(255)" json:"benefit"`
	StartDate    *time.Time     `json:"startDate"`
	EndDate      *time.Time     `json:"endDate"`
	Status       string         `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	Notes        *string        `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EconomicStabilization) TableName() string {
	return "salvia.economic_stabilization"
}
