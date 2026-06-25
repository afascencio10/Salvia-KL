package models

import (
	"time"

	"gorm.io/gorm"
)

// Directory representa una entidad de contacto institucional listada en la app Flutter.
// city_id almacena security.city.city_i_code (no el city_id numérico).
//
// SQL equivalente:
//
//	CREATE TABLE salvia.directories (
//	    id              VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    city_id         VARCHAR(36)  NOT NULL,
//	    creation_date   TIMESTAMPTZ  NOT NULL,
//	    update_date     TIMESTAMPTZ  NOT NULL,
//	    name            VARCHAR(128) NOT NULL,
//	    address         VARCHAR(256) NOT NULL,
//	    phone           VARCHAR(32),
//	    email           VARCHAR(128),
//	    opening_hours   VARCHAR(128),
//	    type            VARCHAR(2)   NOT NULL
//	);
type Directory struct {
	ID           string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CityID       string    `gorm:"type:varchar(36);not null;index;column:city_id" json:"cityId"`
	CreationDate time.Time `gorm:"column:creation_date;not null" json:"creationDate"`
	UpdateDate   time.Time `gorm:"column:update_date;not null" json:"updateDate"`
	Name         string    `gorm:"type:varchar(128);not null" json:"name"`
	Address      string    `gorm:"type:varchar(256);not null" json:"address"`
	Phone        *string   `gorm:"type:varchar(32)" json:"phone,omitempty"`
	Email        *string   `gorm:"type:varchar(128)" json:"email,omitempty"`
	OpeningHours *string   `gorm:"type:varchar(128);column:opening_hours" json:"openingHours,omitempty"`
	Type         string    `gorm:"type:varchar(2);not null;index" json:"type"`
}

func (Directory) TableName() string { return "salvia.directories" }

func (d *Directory) BeforeCreate(_ *gorm.DB) error {
	now := time.Now()
	if d.CreationDate.IsZero() {
		d.CreationDate = now
	}
	if d.UpdateDate.IsZero() {
		d.UpdateDate = now
	}
	return nil
}

func (d *Directory) BeforeUpdate(_ *gorm.DB) error {
	d.UpdateDate = time.Now()
	return nil
}

// Tipos válidos de directorio (código de 2 caracteres).
const (
	DirectoryTypeFiscalia           = "FI"
	DirectoryTypeComisariaFamilia   = "CF"
	DirectoryTypeRedUrgencias       = "RU"
	DirectoryTypeLineasEmergencia   = "LE"
)

// DirectoryTypeLabels mapea código → etiqueta para UI y respuestas API.
var DirectoryTypeLabels = map[string]string{
	DirectoryTypeFiscalia:         "Fiscalías",
	DirectoryTypeComisariaFamilia: "Comisarías de familia",
	DirectoryTypeRedUrgencias:     "Red de Urgencias",
	DirectoryTypeLineasEmergencia: "Líneas de Emergencia",
}

// ValidDirectoryTypes lista los códigos permitidos en el campo type.
var ValidDirectoryTypes = []string{
	DirectoryTypeFiscalia,
	DirectoryTypeComisariaFamilia,
	DirectoryTypeRedUrgencias,
	DirectoryTypeLineasEmergencia,
}

// IsValidDirectoryType indica si el código de tipo es uno de los permitidos.
func IsValidDirectoryType(code string) bool {
	_, ok := DirectoryTypeLabels[code]
	return ok
}
