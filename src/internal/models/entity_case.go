package models

import (
	"time"

	"gorm.io/gorm"
)

// EntityCase relaciona una sede de entidad (entity_branch) con un caso (victim_case).
// Representa "esta víctima/caso está gestionando algo con esta entidad" — se crea
// manualmente por un agente desde el componente case-entities, no se infiere
// automáticamente de barreras u oficios.
//
// Ver DocsMD/Componentes/case-entities/related-tables.md.
type EntityCase struct {
	ID             string  `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID         string  `gorm:"type:varchar(36);not null;index;column:case_id" json:"caseId"`
	EntityBranchID int64   `gorm:"not null;index;column:entity_branch_id" json:"entityBranchId"`
	Objetivo       *string `gorm:"type:text;column:objetivo" json:"objetivo,omitempty"`

	// LastAction es un campo de texto libre editado directamente — no se calcula.
	// Pendiente definir el mecanismo exacto para poblarlo (ver GAP en related-tables.md).
	LastAction *string `gorm:"type:text;column:last_action" json:"lastAction,omitempty"`

	CreatedByID string `gorm:"type:varchar(36);column:created_by_id" json:"createdById"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityCase) TableName() string {
	return "salvia.entity_case"
}

// EntityCaseWithRelations es el shape enriquecido que devuelve
// GET /api/v1/casos/:caseId/entidades — una fila de entity_case ya resuelta
// contra entity_branch, entity, la cadena town/city/department, y los conteos
// de entity_letter (oficios) y barrier_v2 (barreras activas).
type EntityCaseWithRelations struct {
	RelID              string  `gorm:"column:rel_id"                json:"relId"`
	EntityBranchID     int64   `gorm:"column:entity_branch_id"      json:"entityBranchId"`
	EntityBranchICode  string  `gorm:"column:entity_branch_icode"   json:"entityBranchICode"`
	EntityBranchName   string  `gorm:"column:entity_branch_name"    json:"entityBranchName"`
	Sector             string  `gorm:"column:sector"                json:"sector"`
	Address            string  `gorm:"column:address"                json:"address"`
	DepartmentName     string  `gorm:"column:department_name"       json:"departmentName"`
	CityName           string  `gorm:"column:city_name"             json:"cityName"`
	TownName           string  `gorm:"column:town_name"             json:"townName"`
	Objetivo           *string `gorm:"column:objetivo"              json:"objetivo,omitempty"`
	LastAction         *string `gorm:"column:last_action"           json:"lastAction,omitempty"`
	OficiosCount       int64   `gorm:"column:oficios_count"         json:"oficiosCount"`
	BarrerasActivasCount int64 `gorm:"column:barreras_activas_count" json:"barrerasActivasCount"`
	CreatedByID        string  `gorm:"column:created_by_id"         json:"createdById"`
	CreatedAt          time.Time `gorm:"column:created_at"          json:"createdAt"`
}
