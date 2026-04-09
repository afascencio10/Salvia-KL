package models

import (
	"time"

	"gorm.io/gorm"
)

// BarrierV2 representa la tabla salvia.barrier_v2.
type BarrierV2 struct {
	ID          string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	CaseID      string         `gorm:"type:varchar(36);not null;index"`
	FollowUpID  string         `gorm:"type:varchar(36);not null;index"`
	Sector      string         `gorm:"type:varchar(50);not null"`
	Description string         `gorm:"type:text;not null"`
	Status      string         `gorm:"type:varchar(20);default:'OPEN'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (BarrierV2) TableName() string { return "salvia.barrier_v2" }
