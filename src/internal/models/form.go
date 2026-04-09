package models

import (
	"time"

	"gorm.io/gorm"
)

// Form representa un formulario del sistema.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.form (
//	    id          VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    name        VARCHAR(255) NOT NULL,
//	    description TEXT NOT NULL,
//	    status      VARCHAR(20) NOT NULL DEFAULT 'active',
//	    created_at  TIMESTAMPTZ,
//	    updated_at  TIMESTAMPTZ,
//	    deleted_at  TIMESTAMPTZ
//	);
type Form struct {
	ID          string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null"                            json:"name"`
	Description string         `gorm:"type:text;not null"                                    json:"description"`
	Status      string         `gorm:"type:varchar(20);not null;default:'active'"            json:"status"`
	CreatedAt   time.Time      `                                                             json:"createdAt"`
	UpdatedAt   time.Time      `                                                             json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Form) TableName() string { return "salvia.form" }

