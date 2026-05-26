package models

import (
	"time"

	"gorm.io/gorm"
)

// MenTeamRemision modelo GORM que mapea la tabla salvia.men_team_remision.
type MenTeamRemision struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID       string         `gorm:"type:varchar(36);not null"                      json:"caseId"`
	FollowUpID   string         `gorm:"type:uuid;not null"                             json:"followUpId"`
	Type         string         `gorm:"type:varchar(50);not null"                      json:"type"`
	Provider     string         `gorm:"type:varchar(255)"                              json:"provider"`
	ScheduledAt  *time.Time     `                                                      json:"scheduledAt"`
	SessionCount int            `gorm:"default:0"                                      json:"sessionCount"`
	Status       string         `gorm:"type:varchar(20);default:'ACTIVE'"               json:"status"`
	Notes        *string        `gorm:"type:text"                                      json:"notes"`
	CreatedAt    time.Time      `                                                      json:"createdAt"`
	UpdatedAt    time.Time      `                                                      json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

func (MenTeamRemision) TableName() string {
	return "salvia.men_team_remision"
}
