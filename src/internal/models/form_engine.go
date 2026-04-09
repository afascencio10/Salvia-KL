package models

import (
	"time"

	"gorm.io/gorm"
)

// ─── RepeaterGroup ────────────────────────────────────────────────────────────

type RepeaterGroup struct {
	ID             string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormSectionID  string         `gorm:"type:varchar(36);not null;index"                       json:"formSectionId"`
	Name           string         `gorm:"type:varchar(255);not null"                            json:"name"`
	Order          int            `gorm:"column:order;not null;default:0"                       json:"order"`
	MinRepetitions *int           `gorm:"column:min_repetitions"                                json:"minRepetitions"`
	MaxRepetitions *int           `gorm:"column:max_repetitions"                                json:"maxRepetitions"`
	AddButtonLabel *string        `gorm:"type:varchar(255)"                                     json:"addButtonLabel"`
	CreatedAt      time.Time      `                                                             json:"createdAt"`
	UpdatedAt      time.Time      `                                                             json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (RepeaterGroup) TableName() string { return "salvia.repeater_group" }

// ─── Question ─────────────────────────────────────────────────────────────────

type Question struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormID          string         `gorm:"type:varchar(36);not null;index"                       json:"formId"`
	FormSectionID   string         `gorm:"type:varchar(36);not null;index"                       json:"formSectionId"`
	RepeaterGroupID *string        `gorm:"type:varchar(36);index"                                json:"repeaterGroupId"`
	QuestionType    string         `gorm:"type:varchar(50);not null"                             json:"questionType"`
	Description     string         `gorm:"type:text;not null"                                    json:"description"`
	Metadata        *string        `gorm:"type:text"                                             json:"metadata"`
	Options         *string        `gorm:"type:text"                                             json:"options"`
	Order           int            `gorm:"column:order;not null;default:0"                       json:"order"`
	CreatedAt       time.Time      `                                                             json:"createdAt"`
	UpdatedAt       time.Time      `                                                             json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Question) TableName() string { return "salvia.question" }

// ─── VisibilityCondition ──────────────────────────────────────────────────────

type VisibilityCondition struct {
	ID                string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	TargetType        string         `gorm:"type:varchar(50);not null"                             json:"targetType"`
	TargetID          string         `gorm:"type:varchar(36);not null;index"                       json:"targetId"`
	TriggerQuestionID string         `gorm:"type:varchar(36);not null;index"                       json:"triggerQuestionId"`
	TriggerOptionID   *string        `gorm:"type:varchar(36)"                                      json:"triggerOptionId"`
	TriggerValue      *string        `gorm:"type:varchar(255)"                                     json:"triggerValue"`
	Operator          string         `gorm:"type:varchar(50);not null"                             json:"operator"`
	Logic             string         `gorm:"type:varchar(10);not null"                             json:"logic"`
	CreatedAt         time.Time      `                                                             json:"createdAt"`
	UpdatedAt         time.Time      `                                                             json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (VisibilityCondition) TableName() string { return "salvia.visibility_condition" }
