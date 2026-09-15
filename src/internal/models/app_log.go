// Package models — app_log.go
// Modelo genérico de logs estructurados para funciones del backend (ej. auto-asignación).
package models

import (
	"time"

	"gorm.io/datatypes"
)

// AppLog es un registro de log estructurado, append-only.
//
// Key es el criterio de filtrado (ej. la cédula/i_code del caso) para poder
// recuperar rápido todos los logs relacionados con un caso puntual, por
// ejemplo cuando se reporta un caso mal asignado.
type AppLog struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	Key       string         `gorm:"type:varchar(150);index"                                json:"key"`
	Message   string         `gorm:"type:text"                                              json:"message"`
	Data      datatypes.JSON `gorm:"type:jsonb"                                             json:"data"`
	CreatedAt time.Time      `gorm:"index"                                                  json:"createdAt"`
}

func (AppLog) TableName() string { return "salvia.app_log" }
