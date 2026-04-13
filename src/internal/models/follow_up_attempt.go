package models

import (
	"time"

	"gorm.io/gorm"
)

// FollowUpAttempt representa un intento de contacto en un seguimiento
type FollowUpAttempt struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FollowUpID  string         `gorm:"type:uuid;not null;index" json:"follow_up_id"`
	Reason      string         `gorm:"type:text" json:"reason"`
	WasAnswered bool           `gorm:"default:false" json:"was_answered"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FollowUpAttempt) TableName() string {
	return "salvia.follow_up_attempts"
}
