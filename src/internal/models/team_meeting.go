package models

import (
	"time"

	"gorm.io/gorm"
)

// TeamMeeting reunión interna del equipo psicosocial (sin víctima).
// Convocada por supervisor. Se asocia a N agentes via team_meeting_agent.
type TeamMeeting struct {
	ID            string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	Title         string         `gorm:"type:varchar(200);not null" json:"title"`
	Description   *string        `gorm:"type:text" json:"description"`
	ScheduledDate time.Time      `gorm:"column:scheduled_date;type:date;not null" json:"scheduledDate"`
	ScheduledTime string         `gorm:"type:varchar(8);column:scheduled_time;not null" json:"scheduledTime"`
	DurationMin   int            `gorm:"column:duration_min;default:60" json:"durationMin"`
	CreatedBy     string         `gorm:"type:varchar(36);column:created_by;not null" json:"createdBy"`
	CreatedByName *string        `gorm:"type:varchar(255);column:created_by_name" json:"createdByName"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TeamMeeting) TableName() string { return "salvia.team_meeting" }
