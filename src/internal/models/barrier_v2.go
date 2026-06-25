package models

import (
	"time"

	"gorm.io/gorm"
)

// BarrierV2 modelo GORM que mapea la tabla salvia.barrier_v2.
// Cada registro representa una entrada completa del repeater de barreras.
type BarrierV2 struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID      string         `gorm:"type:varchar(36);not null" json:"caseId"`
	FollowUpID  string         `gorm:"type:uuid;not null" json:"followUpId"`
	CreatedByID string         `gorm:"type:varchar(36);index" json:"createdById"`
	Status      string         `gorm:"type:varchar(20);default:'OPEN'" json:"status"`

	// Sector y barreras específicas del bloque condicional
	Sector               string `gorm:"type:varchar(50)" json:"sector"`                 // Q1  dropdown
	SpecificBarriers     string `gorm:"type:text" json:"specificBarriers"`               // Q2/Q5/Q8  CSV de opciones del sector
	SpecificInstitutions string `gorm:"type:text" json:"specificInstitutions"`           // Q3/Q6/Q9  CSV de instituciones del sector
	OtherBarrierDesc     string `gorm:"type:text" json:"otherBarrierDesc"`               // Q4/Q7/Q10 texto "otra barrera"
	InstitutionName      string `gorm:"type:text" json:"institutionName"`                // Q11 nombre (solo otras_instituciones)

	// Ubicación geográfica
	DepartmentID string `gorm:"type:varchar(20)" json:"departmentId"` // Q12
	CityID       string `gorm:"type:varchar(20)" json:"cityId"`       // Q13
	TownID       string `gorm:"type:varchar(20)" json:"townId"`       // Q14

	// Barreras estructurales (comunes a toda entrada)
	StructuralInstitutional string `gorm:"type:text" json:"structuralInstitutional"` // Q15 CSV
	StructuralEconomic      string `gorm:"type:text" json:"structuralEconomic"`      // Q16 CSV
	StructuralTerritorial   string `gorm:"type:text" json:"structuralTerritorial"`   // Q17 CSV
	StructuralDifferential  string `gorm:"type:text" json:"structuralDifferential"`  // Q18 CSV

	// Detalles de la barrera
	BarrierDate        string `gorm:"type:varchar(20)" json:"barrierDate"`      // Q19 date
	OfficialDependency string `gorm:"type:text" json:"officialDependency"`      // Q20 text
	Description        string `gorm:"type:text" json:"description"`             // Q21 text
	ManagementActions  string `gorm:"type:text" json:"managementActions"`       // Q22 CSV

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BarrierV2) TableName() string {
	return "salvia.barrier_v2"
}

// Estados válidos de BarrierV2.
const (
	BarrierV2StatusOpen       = "OPEN"
	BarrierV2StatusEnGestion  = "En Gestion"
	BarrierV2StatusArticulada = "Articulada"
	BarrierV2StatusManaged    = "MANAGED"
)
