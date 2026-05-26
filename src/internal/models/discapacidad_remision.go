package models

import (
	"time"

	"gorm.io/gorm"
)

// DiscapacidadRemision modelo GORM que mapea la tabla salvia.discapacidad_remision.
// Registra la remisión al equipo de Discapacidad (Salvia Dignidad).
// Se crea un registro por cada servicio seleccionado:
//   - apoyo_lsc           → Apoyo comunicativo en Lengua de Señas Colombiana (LSC)
//   - enfoque_discapacidad → Apoyo con enfoque de discapacidad (ley de capacidad legal)
type DiscapacidadRemision struct {
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID     string         `gorm:"type:varchar(36);not null"                      json:"caseId"`
	FollowUpID string         `gorm:"type:uuid;not null"                             json:"followUpId"`
	Service    string         `gorm:"type:varchar(50);not null"                      json:"service"`
	Status     string         `gorm:"type:varchar(20);default:'ACTIVE'"               json:"status"`
	Notes      *string        `gorm:"type:text"                                      json:"notes"`
	CreatedAt  time.Time      `                                                      json:"createdAt"`
	UpdatedAt  time.Time      `                                                      json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

func (DiscapacidadRemision) TableName() string {
	return "salvia.discapacidad_remision"
}
