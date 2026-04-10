package models

import "time"

// VictimCaseLight es un modelo GORM de solo lectura parcial de salvia.victim_case
type VictimCaseLight struct {
	ID              string    `gorm:"column:victim_case_i_code;primaryKey" json:"id"`
	CreationDate    time.Time `gorm:"column:victim_case_creation_date" json:"creationDate"`
	Status          string    `gorm:"column:victim_case_status" json:"status"`
	VictimNames     string    `gorm:"column:victim_case_victim_names" json:"names"`
	VictimLastNames string    `gorm:"column:victim_case_victim_last_names" json:"lastNames"`
	DocType         string    `gorm:"column:victim_case_victim_doc_type" json:"docType"`
	DocNumber       string    `gorm:"column:victim_case_victim_doc_number" json:"docNumber"`
}

func (VictimCaseLight) TableName() string {
	return "salvia.victim_case"
}
