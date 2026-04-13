package models

import (
	"time"

	"gorm.io/gorm"
)

// Answer representa la tabla salvia.answer.
type Answer struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormSubmissionID string         `gorm:"type:varchar(36);not null;index"                       json:"formSubmissionId"`
	QuestionID       string         `gorm:"type:varchar(36);not null;index"                       json:"questionId"`
	RepeaterEntryID  *string        `gorm:"type:varchar(36);index"                                json:"repeaterEntryId,omitempty"`
	Value            string         `gorm:"type:text"                                             json:"value"`
	QuestionSnapshot *string        `gorm:"type:text"                                             json:"questionSnapshot,omitempty"`
	CreatedAt        time.Time      `                                                             json:"createdAt"`
	UpdatedAt        time.Time      `                                                             json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Answer) TableName() string { return "salvia.answer" }
