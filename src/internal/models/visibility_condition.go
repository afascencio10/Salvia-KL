package models

import "time"

// VisibilityCondition representa la tabla salvia.visibility_condition.
// No tiene DeletedAt: las condiciones se reemplazan, no se eliminan lógicamente.
type VisibilityCondition struct {
	ID                string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	TargetType        string    `gorm:"type:varchar(50);not null"                             json:"targetType"`
	TargetID          string    `gorm:"type:varchar(36);not null"                             json:"targetId"`
	TriggerQuestionID string    `gorm:"type:varchar(36);not null;index"                       json:"triggerQuestionId"`
	TriggerStatePath  *string   `gorm:"type:varchar(255)"                                     json:"triggerStatePath,omitempty"`
	TriggerValue      *string   `gorm:"type:varchar(255)"                                     json:"triggerValue,omitempty"`
	Operator          string    `gorm:"type:varchar(20);not null"                             json:"operator"`
	CreatedAt         time.Time `                                                             json:"createdAt"`
	UpdatedAt         time.Time `                                                             json:"updatedAt"`
}

func (VisibilityCondition) TableName() string { return "salvia.visibility_condition" }
