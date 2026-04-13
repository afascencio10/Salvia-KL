package models

import (
	"time"

	"gorm.io/gorm"
)

// Question representa la tabla salvia.question.
type Question struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormID          string         `gorm:"type:varchar(36);not null;index"                       json:"formId"`
	FormSectionID   string         `gorm:"type:varchar(36);not null;index"                       json:"formSectionId"`
	RepeaterGroupID *string        `gorm:"type:varchar(36);index"                                json:"repeaterGroupId,omitempty"`
	QuestionTypeID  string         `gorm:"column:question_type;type:varchar(50);not null"        json:"questionTypeId"`
	Description     string         `gorm:"type:text;not null"                                    json:"description"`
	Required        bool           `gorm:"column:required;default:false"                         json:"required"`
	Order           int            `gorm:"column:order;default:0"                                json:"order"`
	CreatedAt       time.Time      `                                                             json:"createdAt"`
	UpdatedAt       time.Time      `                                                             json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Question) TableName() string { return "salvia.question" }
