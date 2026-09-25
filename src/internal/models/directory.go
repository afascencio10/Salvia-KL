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
//	    name            TEXT         NOT NULL,
//	    address         TEXT         NOT NULL,
//	    phone           TEXT,
//	    email           TEXT,
//	    opening_hours   TEXT,
//	    type            VARCHAR(2)   NOT NULL
//	);
//
// name/address/phone/email/opening_hours son TEXT y no VARCHAR(n): el "Mapeo de actores"
// de SALVIA trae teléfonos, horarios y nombres de varios cientos de caracteres.
type Directory struct {
	ID           string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CityID       string    `gorm:"type:varchar(36);not null;index;column:city_id" json:"cityId"`
	DepartmentID *uint64   `gorm:"column:department_id;index" json:"departmentId,omitempty"`
	TownICode    *string   `gorm:"type:varchar(36);column:town_i_code" json:"townICode,omitempty"`
	CreationDate time.Time `gorm:"column:creation_date;not null" json:"creationDate"`
	UpdateDate   time.Time `gorm:"column:update_date;not null" json:"updateDate"`
	Name         string    `gorm:"type:text;not null" json:"name"`
	Address      string    `gorm:"type:text;not null" json:"address"`
	Phone        *string   `gorm:"type:text" json:"phone,omitempty"`
	Email        *string   `gorm:"type:text" json:"email,omitempty"`
	OpeningHours *string   `gorm:"type:text;column:opening_hours" json:"openingHours,omitempty"`
	Type         string    `gorm:"type:varchar(2);not null;index" json:"type"`
	Sector       *string   `gorm:"type:varchar(2);column:sector;index" json:"sector,omitempty"`
	Latitude     *float64  `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude    *float64  `gorm:"column:longitude" json:"longitude,omitempty"`
	Active       bool      `gorm:"column:active;not null;default:true;index" json:"active"`
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
//
// FI/CF/RU/LE son los 4 tipos originales que consume la app Flutter (una pantalla
// por tipo, siempre consultando con ?type=). El resto replica el desplegable
// "Nombre de la institución/entidad" del "Mapeo de actores" de SALVIA, usado por
// la carga masiva (directory_bulk_upload_service.go).
const (
	DirectoryTypeFiscalia         = "FI"
	DirectoryTypeComisariaFamilia = "CF"
	DirectoryTypeRedUrgencias     = "RU"
	DirectoryTypeLineasEmergencia = "LE"

	DirectoryTypeGobernacion             = "GO"
	DirectoryTypeAlcaldia                = "AL"
	DirectoryTypeSecretariaSalud         = "SS"
	DirectoryTypeSecretariaGenero        = "SG"
	DirectoryTypeSecretariaEducacion     = "SE"
	DirectoryTypeUARIV                   = "UA"
	DirectoryTypeUNP                     = "UN"
	DirectoryTypeDefensoriaPueblo        = "DP"
	DirectoryTypeProcuraduria            = "PC"
	DirectoryTypePersoneria              = "PE"
	DirectoryTypeHospital                = "HO"
	DirectoryTypeCentroSalud             = "CS"
	DirectoryTypeIPSEPS                  = "IE"
	DirectoryTypeClinicaPrivada          = "CP"
	DirectoryTypePoliciaNacional         = "PN"
	DirectoryTypeICBF                    = "IC"
	DirectoryTypeCentroZonalICBF         = "CZ"
	DirectoryTypeCAIVAS                  = "CI"
	DirectoryTypeCAVIF                   = "CV"
	DirectoryTypeJuzgados                = "JZ"
	DirectoryTypeInspeccionesPolicia     = "IP"
	DirectoryTypeInstitucionesEducativas = "IN"
	DirectoryTypeCentroIntegrate         = "CT"
	DirectoryTypeCasaLGBTIQ              = "CL"
	DirectoryTypeJusticiaPropia          = "JP"
	DirectoryTypeCasaJusticia            = "CJ"
	DirectoryTypeCabildoIndigena         = "CB"
	DirectoryTypeOtra                    = "OT"
)

// DirectoryTypeLabels mapea código → etiqueta para UI y respuestas API.
// Es la fuente de verdad del catálogo (IsValidDirectoryType).
var DirectoryTypeLabels = map[string]string{
	DirectoryTypeFiscalia:         "Fiscalías",
	DirectoryTypeComisariaFamilia: "Comisarías de familia",
	DirectoryTypeRedUrgencias:     "Red de Urgencias",
	DirectoryTypeLineasEmergencia: "Líneas de Emergencia",

	DirectoryTypeGobernacion:             "Gobernación",
	DirectoryTypeAlcaldia:                "Alcaldía",
	DirectoryTypeSecretariaSalud:         "Secretaría de Salud",
	DirectoryTypeSecretariaGenero:        "Secretaría de Género / de la Mujer",
	DirectoryTypeSecretariaEducacion:     "Secretaría de Educación",
	DirectoryTypeUARIV:                   "UARIV",
	DirectoryTypeUNP:                     "UNP",
	DirectoryTypeDefensoriaPueblo:        "Defensoría del Pueblo",
	DirectoryTypeProcuraduria:            "Procuraduría",
	DirectoryTypePersoneria:              "Personería",
	DirectoryTypeHospital:                "Hospital",
	DirectoryTypeCentroSalud:             "Centro / Puesto de Salud",
	DirectoryTypeIPSEPS:                  "IPS/EPS",
	DirectoryTypeClinicaPrivada:          "Clínica privada",
	DirectoryTypePoliciaNacional:         "Policía Nacional",
	DirectoryTypeICBF:                    "ICBF",
	DirectoryTypeCentroZonalICBF:         "Centro zonal ICBF",
	DirectoryTypeCAIVAS:                  "CAIVAS",
	DirectoryTypeCAVIF:                   "CAVIF",
	DirectoryTypeJuzgados:                "Juzgados",
	DirectoryTypeInspeccionesPolicia:     "Inspecciones de Policía",
	DirectoryTypeInstitucionesEducativas: "Instituciones educativas",
	DirectoryTypeCentroIntegrate:         "Centro Intégrate",
	DirectoryTypeCasaLGBTIQ:              "Casa LGBTIQ+",
	DirectoryTypeJusticiaPropia:          "Justicia propia",
	DirectoryTypeCasaJusticia:            "Casa de Justicia",
	DirectoryTypeCabildoIndigena:         "Cabildos indígenas",
	DirectoryTypeOtra:                    "Otra",
}

// IsValidDirectoryType indica si el código de tipo es uno de los permitidos.
func IsValidDirectoryType(code string) bool {
	_, ok := DirectoryTypeLabels[code]
	return ok
}

// Sectores agrupadores para el Directorio de Entidades (código de 2 caracteres).
// Los 5 últimos vienen de la columna de sector del "Mapeo de actores" de SALVIA.
const (
	DirectorySectorSalud                   = "SA"
	DirectorySectorJusticia                = "JU"
	DirectorySectorProteccion              = "PR"
	DirectorySectorJusticiaPropia          = "JP"
	DirectorySectorMinisterioPublico       = "MP"
	DirectorySectorAtencionHumanitaria     = "AH"
	DirectorySectorEducacion               = "ED"
	DirectorySectorServiciosEspecializados = "ES"
)

// DirectorySectorLabels mapea código de sector → etiqueta para UI y respuestas API.
// Es la fuente de verdad del catálogo (IsValidDirectorySector).
var DirectorySectorLabels = map[string]string{
	DirectorySectorSalud:                   "Salud",
	DirectorySectorJusticia:                "Justicia",
	DirectorySectorProteccion:              "Protección",
	DirectorySectorJusticiaPropia:          "Justicia Propia",
	DirectorySectorMinisterioPublico:       "Ministerio Público",
	DirectorySectorAtencionHumanitaria:     "Atención Humanitaria",
	DirectorySectorEducacion:               "Educación",
	DirectorySectorServiciosEspecializados: "Servicios especializados",
}

// IsValidDirectorySector indica si el código de sector es uno de los permitidos.
func IsValidDirectorySector(code string) bool {
	_, ok := DirectorySectorLabels[code]
	return ok
}
