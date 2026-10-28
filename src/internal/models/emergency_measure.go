package models

import (
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
