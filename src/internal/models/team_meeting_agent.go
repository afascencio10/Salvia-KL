package models

import (
	"time"

	"gorm.io/gorm"
)

// TeamMeetingAgent relación N:M entre reunión y agentes invitados.
type TeamMeetingAgent struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	MeetingID string         `gorm:"type:varchar(36);column:meeting_id;not null;index" json:"meetingId"`
	AgentID   string         `gorm:"type:varchar(36);column:agent_id;not null;index" json:"agentId"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TeamMeetingAgent) TableName() string { return "salvia.team_meeting_agent" }
