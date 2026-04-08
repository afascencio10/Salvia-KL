// Package models contiene los structs GORM para el nuevo patrón de repositorios.
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

// FormSection representa una sección dentro de un Form.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.form_section (
//	    id          VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    form_id     VARCHAR(36) NOT NULL REFERENCES salvia.form(id),
//	    name        VARCHAR(255) NOT NULL,
//	    description TEXT,
//	    "order"     INT NOT NULL DEFAULT 0,
//	    created_at  TIMESTAMPTZ,
//	    updated_at  TIMESTAMPTZ,
//	    deleted_at  TIMESTAMPTZ
//	);
type FormSection struct {
	ID          string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	FormID      string         `gorm:"type:varchar(36);not null;index"                       json:"formId"`
	Name        string         `gorm:"type:varchar(255);not null"                            json:"name"`
	Description *string        `gorm:"type:text"                                             json:"description"`
	Order       int            `gorm:"column:order;not null;default:0"                       json:"order"`
	CreatedAt   time.Time      `                                                             json:"createdAt"`
	UpdatedAt   time.Time      `                                                             json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (FormSection) TableName() string { return "salvia.form_section" }
