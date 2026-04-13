package models

import (
	"time"

	"gorm.io/gorm"
)

// BarrierV2 modelo GORM que mapea la tabla salvia.barrier_v2.
type BarrierV2 struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID      string         `gorm:"type:varchar(36);not null" json:"caseId"`
	FollowUpID  string         `gorm:"type:uuid;not null" json:"followUpId"`
	Sector      string         `gorm:"type:varchar(50);not null" json:"sector"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Status      string         `gorm:"type:varchar(20);default:'OPEN'" json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BarrierV2) TableName() string {
	return "salvia.barrier_v2"
}
