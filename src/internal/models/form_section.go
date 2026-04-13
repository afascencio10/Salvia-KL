package models

import (
	"time"

	"gorm.io/gorm"
)

// FormSection representa la tabla salvia.form_section.
type FormSection struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FormID      string         `gorm:"type:uuid;not null;index"                       json:"formId"`
	Name        string         `gorm:"type:varchar(255);not null"                     json:"name"`
	Description *string        `gorm:"type:text"                                      json:"description,omitempty"`
	Order       int            `gorm:"column:order;default:0"                         json:"order"`
	CreatedAt   time.Time      `                                                      json:"createdAt"`
	UpdatedAt   time.Time      `                                                      json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                          json:"deletedAt,omitempty"`
}

func (FormSection) TableName() string { return "salvia.form_section" }
