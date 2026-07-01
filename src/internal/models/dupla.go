package models

import (
	"time"

	"gorm.io/gorm"
)

// Dupla representa un par psicóloga + trabajador social del equipo de Atención Psicosocial.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.dupla (
//	    id                 VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    name               VARCHAR(36) NOT NULL,
//	    psychologist_id    VARCHAR(36) NOT NULL,
//	    social_worker_id   VARCHAR(36) NOT NULL,
//	    created_at         TIMESTAMPTZ,
//	    updated_at         TIMESTAMPTZ,
//	    deleted_at         TIMESTAMPTZ
//	);
type Dupla struct {
	ID             string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string         `gorm:"type:varchar(36);not null"                         json:"name"`
	PsychologistID string         `gorm:"type:varchar(36);not null;column:psychologist_id"  json:"psychologistId"`
	SocialWorkerID string         `gorm:"type:varchar(36);not null;column:social_worker_id" json:"socialWorkerId"`
	CreatedAt      time.Time      `                                                          json:"createdAt"`
	UpdatedAt      time.Time      `                                                          json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index"                                              json:"deletedAt,omitempty"`
}

func (Dupla) TableName() string { return "salvia.dupla" }
