package models

import (
	"time"

	"gorm.io/datatypes"
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

	Comments  []CaseArticulatingMechanismComment `gorm:"foreignKey:MechanismID;references:ID" json:"comments,omitempty"`
	Gestiones []MechanismGestionItem             `gorm:"-" json:"gestiones"`
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

// MechanismGestionItem es una case_task gestion_mecanismo Done, para la card
// del Detalle del Caso (comentario = result, Ver = formData).
type MechanismGestionItem struct {
	ID               string         `gorm:"column:id" json:"id"`
	MechanismID      string         `gorm:"column:mechanism_id" json:"mechanismId"`
	Result           *string        `gorm:"column:result" json:"result"`
	FormData         datatypes.JSON `gorm:"column:form_data" json:"formData"`
	CompletedAt      *time.Time     `gorm:"column:completed_at" json:"completedAt"`
	AssignedUserName string         `gorm:"column:assigned_user_name" json:"assignedUserName"`
}

// CaseMechanismMineItem es una fila del listado Casos Mecanismo (rol mc).
type CaseMechanismMineItem struct {
	ID               string     `gorm:"column:id" json:"id"`
	CaseID           string     `gorm:"column:case_id" json:"caseId"`
	CaseCode         string     `gorm:"column:case_code" json:"caseCode"`
	VictimFullName   string     `gorm:"column:victim_full_name" json:"victimFullName"`
	Document         string     `gorm:"column:document" json:"document"`
	Phone            string     `gorm:"column:phone" json:"phone"`
	City             string     `gorm:"column:city" json:"city"`
	RiskLevel        string     `gorm:"column:risk_level" json:"riskLevel"`
	CaseStatus       string     `gorm:"column:case_status" json:"caseStatus"`
	CaseStatusLabel  string     `gorm:"-" json:"caseStatusLabel"`
	Scope            string     `gorm:"column:scope" json:"scope"`
	SalviaCommitment *string    `gorm:"column:salvia_commitment" json:"salviaCommitment,omitempty"`
	LastComment      *string    `gorm:"column:last_comment" json:"lastComment,omitempty"`
	LastCommentAt    *time.Time `gorm:"column:last_comment_at" json:"lastCommentAt,omitempty"`
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
