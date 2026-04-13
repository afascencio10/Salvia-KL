package models

import (
	"time"

	"gorm.io/gorm"
)

// RepeaterGroup representa la tabla salvia.repeater_group.
type RepeaterGroup struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormSectionID   string         `gorm:"type:varchar(36);not null;index"                       json:"formSectionId"`
	Name            string         `gorm:"type:varchar(255);not null"                            json:"name"`
	ItemName        *string        `gorm:"type:varchar(100)"                                     json:"itemName,omitempty"`
	Order           int            `gorm:"column:order;default:0"                                json:"order"`
	MinRepetitions  int            `gorm:"default:0"                                             json:"minRepetitions"`
	MaxRepetitions  *int           `                                                             json:"maxRepetitions,omitempty"`
	CreatedAt       time.Time      `                                                             json:"createdAt"`
	UpdatedAt       time.Time      `                                                             json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (RepeaterGroup) TableName() string { return "salvia.repeater_group" }
