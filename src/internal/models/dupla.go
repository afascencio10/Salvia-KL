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

// DuplaAdminItem dupla activa enriquecida para pantalla Administrar Duplas (E01).
type DuplaAdminItem struct {
	ID               string `json:"id" gorm:"column:id"`
	Name             string `json:"name" gorm:"column:name"`
	PsychologistID   string `json:"psychologistId" gorm:"column:psychologist_id"`
	PsychologistName string `json:"psychologistName" gorm:"column:psychologist_name"`
	SocialWorkerID   string `json:"socialWorkerId" gorm:"column:social_worker_id"`
	SocialWorkerName string `json:"socialWorkerName" gorm:"column:social_worker_name"`
}
