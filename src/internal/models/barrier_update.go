package models

import "time"

// BarrierUpdate representa la tabla salvia.barrier_update.
// Es append-only (log de cambios), sin DeletedAt.
type BarrierUpdate struct {
	ID               string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	BarrierID        string    `gorm:"type:varchar(36);not null;index"`
	FormSubmissionID *string   `gorm:"type:varchar(36);index"` // nullable
	Type             string    `gorm:"type:varchar(50);not null"`
	TargetEntity     string    `gorm:"type:varchar(255)"`
	AttachmentURL    string    `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (BarrierUpdate) TableName() string { return "salvia.barrier_update" }
