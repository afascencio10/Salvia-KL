package models

import (
	"time"

	"gorm.io/gorm"
)

// FollowUpV2 representa la tabla salvia.follow_up_v2 (versión modernizada).
type FollowUpV2 struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	CaseID           string         `gorm:"type:varchar(36);not null;index"`
	FormSubmissionID *string        `gorm:"type:varchar(36)"` // nullable, se actualiza al enviar Kobo
	AgentID          string         `gorm:"type:varchar(36);not null"`
	Team             string         `gorm:"type:varchar(50)"`
	RiskStatus       string         `gorm:"type:varchar(20)"`
	ScheduledDate    time.Time      `gorm:"type:date;not null"`
	IsCompleted      bool           `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (FollowUpV2) TableName() string { return "salvia.follow_up_v2" }
