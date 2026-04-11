package models

import "time"

// VisibilityCondition representa la tabla salvia.visibility_condition.
// No tiene DeletedAt: las condiciones se reemplazan, no se eliminan lógicamente.
type VisibilityCondition struct {
	ID                string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	TargetType        string    `gorm:"type:varchar(50);not null"` // QUESTION, SECTION, GROUP
	TargetID          string    `gorm:"type:varchar(36);not null"`
	TriggerQuestionID string    `gorm:"type:varchar(36);not null;index"`
	Operator          string    `gorm:"type:varchar(20);not null"` // EQUALS, CONTAINS, GT
	Logic             string    `gorm:"type:varchar(255);not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (VisibilityCondition) TableName() string { return "salvia.visibility_condition" }
