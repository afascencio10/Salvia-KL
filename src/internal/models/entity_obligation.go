package models

import (
	"time"

	"gorm.io/gorm"
)

// EntityObligation es el catálogo de obligaciones institucionales que puede tener
// una entidad (organización) en un caso. Se define a nivel de entity, no de
// entity_branch — todas las sedes de una misma organización comparten catálogo.
type EntityObligation struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	EntityID  int64          `gorm:"not null;index;column:entity_id"                        json:"entityId"`
	Label     string         `gorm:"type:varchar(255);not null"                             json:"label"`
	Order     int            `gorm:"column:order;default:0"                                 json:"order"`
	CreatedAt time.Time      `                                                               json:"createdAt"`
	UpdatedAt time.Time      `                                                               json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                                   json:"deletedAt,omitempty"`
}

func (EntityObligation) TableName() string { return "salvia.entity_obligation" }
