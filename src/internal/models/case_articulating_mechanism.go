package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ArticulatingMechanismScopeMunicipal     = "municipal"
	ArticulatingMechanismScopeDepartamental = "departamental"
	ArticulatingMechanismScopeNacional      = "nacional"
	ArticulatingMechanismRoleCode           = "mc"

	// CaseArticulatingMechanismStatusActive es el estado inicial (abierto).
	CaseArticulatingMechanismStatusActive = "active"
)

// CaseArticulatingMechanism registra un escalamiento del caso ante un
// mecanismo articulador. El usuario asignado tiene rol mc; el alcance del
// modal debe coincidir con el tipo registrado en articulating_mechanism.
type CaseArticulatingMechanism struct {
	ID     string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID string `gorm:"type:varchar(36);not null;index;column:case_id" json:"caseId"`
	Scope  string `gorm:"type:varchar(20);not null;column:scope" json:"scope"`
	Status string `gorm:"type:varchar(20);not null;default:active;index;column:status" json:"status"`

	AssignedUserID   string `gorm:"type:varchar(36);not null;index;column:assigned_user_id" json:"assignedUserId"`
	AssignedUserName string `gorm:"type:varchar(256);not null;column:assigned_user_name" json:"assignedUserName"`

	DepartmentID   *string `gorm:"type:varchar(36);column:department_id" json:"departmentId,omitempty"`
	CityID         *string `gorm:"type:varchar(36);column:city_id" json:"cityId,omitempty"`
	TownCode       *string `gorm:"type:varchar(16);column:town_code" json:"townCode,omitempty"`
	DepartmentName *string `gorm:"type:varchar(128);column:department_name" json:"departmentName,omitempty"`
	CityName       *string `gorm:"type:varchar(128);column:city_name" json:"cityName,omitempty"`
	TownName       *string `gorm:"type:varchar(128);column:town_name" json:"townName,omitempty"`

	SalviaCommitment *string `gorm:"type:text;column:salvia_commitment" json:"salviaCommitment,omitempty"`

	CreatedByID string `gorm:"type:varchar(36);column:created_by_id" json:"createdById"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Comments []CaseArticulatingMechanismComment `gorm:"foreignKey:MechanismID;references:ID" json:"comments"`
}

func (CaseArticulatingMechanism) TableName() string {
	return "salvia.case_articulating_mechanism"
}

// CaseArticulatingMechanismComment es un comentario append-only sobre un
// registro de mecanismo (autor + fecha para trazabilidad).
type CaseArticulatingMechanismComment struct {
	ID          string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	MechanismID string    `gorm:"type:varchar(36);not null;index;column:mechanism_id" json:"mechanismId"`
	AuthorID    string    `gorm:"type:varchar(36);not null;column:author_id" json:"authorId"`
	AuthorName  string    `gorm:"type:varchar(256);not null;column:author_name" json:"authorName"`
	Body        string    `gorm:"type:text;not null;column:body" json:"body"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (CaseArticulatingMechanismComment) TableName() string {
	return "salvia.case_articulating_mechanism_comment"
}

// MechanismUserOption es una fila del select de usuarios mc.
type MechanismUserOption struct {
	ID             string `gorm:"-" json:"id"`
	ICode          string `gorm:"column:icode" json:"icode"`
	FullName       string `gorm:"column:full_name" json:"fullName"`
	TownCode       string `gorm:"column:town_code" json:"townCode"`
	TownName       string `gorm:"column:town_name" json:"townName"`
	CityName       string `gorm:"column:city_name" json:"cityName"`
	DepartmentName string `gorm:"column:department_name" json:"departmentName"`
}
