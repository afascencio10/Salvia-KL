package models

import "time"

// VictimCaseLight es un modelo GORM de solo lectura parcial de salvia.victim_case
type VictimCaseLight struct {
	ID              int       `gorm:"column:victim_case_id;primaryKey" json:"id"`
	VictimICode     string    `gorm:"column:victim_case_i_code" json:"i_code"`
	CreationDate    time.Time `gorm:"column:victim_case_creation_date" json:"creationDate"`
	Status          string    `gorm:"column:victim_case_status" json:"status"`
	VictimNames     string    `gorm:"column:victim_case_victim_names" json:"names"`
	VictimLastNames string    `gorm:"column:victim_case_victim_last_names" json:"lastNames"`
	DocType         string    `gorm:"column:victim_case_victim_doc_type" json:"docType"`
	DocNumber       string    `gorm:"column:victim_case_victim_doc_number" json:"docNumber"`
	TownCode        string    `gorm:"column:victim_case_victim_town_code" json:"town_code"`
	Municipality    string    `gorm:"-" json:"municipality,omitempty"` // Campo virtual, se llena en runtime
}

func (VictimCaseLight) TableName() string {
	return "salvia.victim_case"
}

// TownLight es un modelo GORM de solo lectura parcial de salvia.town
type TownLight struct {
	ID        int    `gorm:"column:town_id;primaryKey" json:"id"`
	TownICode string `gorm:"column:town_i_code" json:"i_code"`
	TownCode  string `gorm:"column:town_code" json:"code"`
	TownName  string `gorm:"column:town_name" json:"name"`
	TownType  string `gorm:"column:town_type" json:"type"`
}

func (TownLight) TableName() string {
	return "security.town"
}
