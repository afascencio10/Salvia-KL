package models

import (
	"time"

	"gorm.io/gorm"
)

<<<<<<< HEAD
// FormSubmission representa la tabla salvia.form_submission.
type FormSubmission struct {
	ID         string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	FormID     string         `gorm:"type:varchar(36);not null;index"`
	FollowUpID *string        `gorm:"type:varchar(36);index"` // nullable
	AgentID    string         `gorm:"type:varchar(36);not null"`
	ScoreTotal *float64       `gorm:"type:numeric(5,2)"` // nullable
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (FormSubmission) TableName() string { return "salvia.form_submission" }
=======
// ─── FormSubmission ───────────────────────────────────────────────────────────

type FormSubmission struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormID    string         `gorm:"type:varchar(36);not null;index"                       json:"formId"`
	CreatedAt time.Time      `                                                             json:"createdAt"`
	UpdatedAt time.Time      `                                                             json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (FormSubmission) TableName() string { return "salvia.form_submission" }

// ─── RepeaterEntry ────────────────────────────────────────────────────────────

type RepeaterEntry struct {
	ID               string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormSubmissionID string    `gorm:"type:varchar(36);not null;index"                       json:"formSubmissionId"`
	RepeaterGroupID  string    `gorm:"type:varchar(36);not null;index"                       json:"repeaterGroupId"`
	Iteration        int       `gorm:"not null;default:1"                                    json:"iteration"`
	CreatedAt        time.Time `                                                             json:"createdAt"`
	UpdatedAt        time.Time `                                                             json:"updatedAt"`
}

func (RepeaterEntry) TableName() string { return "salvia.repeater_entry" }

// ─── Answer ───────────────────────────────────────────────────────────────────

type Answer struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormSubmissionID string         `gorm:"type:varchar(36);not null;index"                       json:"formSubmissionId"`
	QuestionID       string         `gorm:"type:varchar(36);not null;index"                       json:"questionId"`
	RepeaterEntryID  *string        `gorm:"type:varchar(36);index"                                json:"repeaterEntryId"`
	Value            *string        `gorm:"type:text"                                             json:"value"`
	CreatedAt        time.Time      `                                                             json:"createdAt"`
	UpdatedAt        time.Time      `                                                             json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Answer) TableName() string { return "salvia.answer" }
>>>>>>> e10fe9a63702dc7f74e0c51178b6ec39a2464fcb
