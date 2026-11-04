package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// EmergencyMeasure modelo GORM que mapea la tabla salvia.emergency_measure.
type EmergencyMeasure struct {
	ID     string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID string `gorm:"type:varchar(36);not null" json:"caseId"`
	// FollowUpID es nil cuando la medida se crea manualmente (fuera del flujo
	// de seguimiento), p.ej. desde el modal de remisión a equipos.
	FollowUpID       *string    `gorm:"type:uuid" json:"followUpId,omitempty"`
	Type             string     `gorm:"type:varchar(50);not null" json:"type"`
	IssuingAuthority string     `gorm:"type:varchar(255)" json:"issuingAuthority"`
	IssuedAt         *time.Time `json:"issuedAt"`
	ExpiresAt        *time.Time `json:"expiresAt"`
	Status           string     `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	Notes            *string    `gorm:"type:text" json:"notes"`
	// ReferralID agrupa las medidas de una misma remisión de la bandeja.
	// Las filas sin este dato no entran al listado.
	ReferralID *string `gorm:"type:uuid;column:referral_id" json:"referralId,omitempty"`
	// RiesgoInminente marca las medidas de una remisión que el seguimiento
	// declaró de riesgo inminente. Las filas anteriores quedan en false.
	RiesgoInminente bool           `gorm:"column:riesgo_inminente;not null;default:false" json:"riesgoInminente"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmergencyMeasure) TableName() string {
	return "salvia.emergency_measure"
}

const (
	EmergencyMeasureStatusActive    = "ACTIVE"
	EmergencyMeasureStatusEntregada = "ENTREGADA"
	EmergencyMeasureStatusClosed    = "CLOSED"
)

// Tipos de medida creados desde la entrevista de caracterización (G-9).
const (
	EmergencyMeasureTypeAseo            = "aseo"
	EmergencyMeasureTypeExequial        = "exequial"
	EmergencyMeasureTypeAccionSimbolica = "accion_simbolica"
)

// EmergencyMeasureTypeCodes lista, en orden estable, los tipos de medida que
// la caracterización puede solicitar (chips del modal y mapeo del formulario).
var EmergencyMeasureTypeCodes = []string{
	"alojamiento",
	"transporte",
	"alimentacion",
	"vestuario",
	EmergencyMeasureTypeAseo,
	"apoyo_psico",
	EmergencyMeasureTypeExequial,
	EmergencyMeasureTypeAccionSimbolica,
}

// EmergencyMeasureTypeFromLabel devuelve el código del tipo cuyo label
// (EmergencyMeasureTypeLabel) coincide con el texto dado, o "" si no hay
// coincidencia. Se usa con los chips "Medidas requeridas" del modal.
func EmergencyMeasureTypeFromLabel(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return ""
	}
	// Los chips del modal también ofrecen "Otras ME" y "Otras medidas de
	// emergencia"; se reconocen para no perder esa selección.
	for _, code := range append(append([]string{}, EmergencyMeasureTypeCodes...), "otras_me", "otras") {
		if strings.EqualFold(EmergencyMeasureTypeLabel(code), label) {
			return code
		}
	}
	return ""
}

// EmergencyMeasureTypeLabel es el nombre que ve Cruz Roja y el timeline.
func EmergencyMeasureTypeLabel(code string) string {
	switch code {
	case EmergencyMeasureTypeAseo:
		return "Aseo"
	case EmergencyMeasureTypeExequial:
		return "Apoyo exequial"
	case EmergencyMeasureTypeAccionSimbolica:
		return "Acción simbólica"
	case "alojamiento":
		return "Alojamiento"
	case "transporte":
		return "Transporte"
	case "alimentacion":
		return "Alimentación"
	case "vestuario":
		return "Vestuario"
	case "apoyo_psico":
		return "Apoyo psicosocial"
	case "otras_me":
		return "Otras ME"
	case "otras":
		return "Otras medidas de emergencia"
	default:
		return code
	}
}
