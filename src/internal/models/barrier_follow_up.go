package models

import "time"

// BarrierFollowUp registra las respuestas del repeater de seguimiento a barreras
// (Sección 3 del formulario de seguimiento v2). Append-only — sin DeletedAt.
type BarrierFollowUp struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BarrierID string    `gorm:"type:uuid;not null;index" json:"barrierId"`
	FollowUpID string   `gorm:"type:uuid;not null;index" json:"followUpId"`
	CreatedByID string  `gorm:"type:varchar(36);not null" json:"createdById"`

	Persists              bool   `gorm:"not null;default:false" json:"persists"`                // Q2
	InstitutionalResponse string `gorm:"type:varchar(50)" json:"institutionalResponse"`         // Q3
	ManagementActions     string `gorm:"type:text" json:"managementActions"`                    // Q4 CSV
	Actions               string `gorm:"type:text" json:"actions"`                              // Q5
	ClosesBarrier         bool   `gorm:"not null;default:false" json:"closesBarrier"`           // Q6
	ClosureReason         string `gorm:"type:varchar(50)" json:"closureReason"`                 // Q7

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (BarrierFollowUp) TableName() string { return "salvia.barrier_follow_up" }
