package models

import (
	"time"

	"gorm.io/gorm"
)

// Option representa una opción de respuesta para una pregunta.
// Aplica a preguntas de tipo: single, multiple, dropdown.
type Option struct {
	ID         string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	QuestionID string         `gorm:"type:varchar(36);not null;index"                       json:"questionId"`
	Label      string         `gorm:"type:varchar(255);not null"                            json:"label"`
	Value      string         `gorm:"type:varchar(255);not null"                            json:"value"`
	Order      int            `gorm:"column:order;default:0"                                json:"order"`
	CreatedAt  time.Time      `                                                             json:"createdAt"`
	UpdatedAt  time.Time      `                                                             json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Option) TableName() string { return "salvia.option" }
