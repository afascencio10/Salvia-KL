// Package models ÔÇö victim_case.go
// Modelos GORM para las tablas de seguimiento de casos (esquema salvia).
package models

import "time"

// VictimCase mapea salvia.victim_case.
// Solo se mapean los campos que necesita la pantalla de detalle.
type VictimCase struct {
	VictimCaseId           int64   `gorm:"column:victim_case_id;primaryKey"      json:"victimCaseId"`
	VictimCaseICode        string  `gorm:"column:victim_case_i_code"             json:"victimCaseICode"`
	VictimCaseFollowUpId   *int64  `gorm:"column:victim_case_follow_up_id"       json:"victimCaseFollowUpId"`
	VictimCaseStatus       string  `gorm:"column:victim_case_status"             json:"victimCaseStatus"`
	VictimCaseVictimName   string  `gorm:"column:victim_case_victim_names"       json:"victimCaseVictimName"`
	VictimCaseVictimLastName   string  `gorm:"column:victim_case_victim_last_names"       json:"victimCaseVictimLastName"`
	VictimCaseVictimAge    *int    `gorm:"column:victim_case_victim_age"         json:"victimCaseVictimAge"`
	VictimCaseVictimGender string  `gorm:"column:victim_case_victim_gender"      json:"victimCaseVictimGender"`
	VictimCaseVictimDocNumber string `gorm:"column:victim_case_victim_doc_number" json:"victimCaseVictimDocNumber"`
	VictimCaseVictimDocType   string `gorm:"column:victim_case_victim_doc_type"   json:"victimCaseVictimDocType"`
	VictimCaseTownCode     string  `gorm:"column:victim_case_victim_town_code"   json:"victimCaseTownCode"`
	VictimCaseOwnerDesc    string  `gorm:"column:victim_case_owner_description"  json:"victimCaseOwnerDescription"`
	VictimCaseTeam         string  `gorm:"column:victim_case_team"               json:"victimCaseTeam"`
	VictimCaseCreatedAt    time.Time `gorm:"column:victim_case_creation_date"       json:"victimCaseCreatedAt"`
}

func (VictimCase) TableName() string { return "salvia.victim_case" }

// VictimCaseForm1 mapea salvia.victim_case_form_1.
// Contiene los datos del formulario de registro del caso (datos de la v├¡ctima y del agresor).
// La FK victim_case_form1_victim_case apunta a VictimCase.VictimCaseICode.
type VictimCaseForm1 struct {
	VictimCaseForm1VictimCase                    int64   `gorm:"column:victim_case_form1_victim_case;primaryKey"                    json:"victimCaseForm1VictimCase"`
	VictimCaseForm1VictimAggressor               string  `gorm:"column:victim_case_form1_victim_aggressor"                          json:"victimAggressor"`
	VictimCaseForm1VictimRelationshipWithAggressor string `gorm:"column:victim_case_form1_victim_relationship_with_aggressor"        json:"relationshipWithAggressor"`
	VictimCaseForm1VictimAggressorName           string  `gorm:"column:victim_case_form1_victim_aggressor_name"                     json:"aggressorName"`
	VictimCaseForm1VictimAggressorDocType        string  `gorm:"column:victim_case_form1_victim_aggressor_doc_type"                 json:"aggressorDocType"`
	VictimCaseForm1VictimAggressorDocNumber      string  `gorm:"column:victim_case_form1_victim_aggressor_doc_number"               json:"aggressorDocNumber"`
	VictimCaseForm1VictimAggressorAddress        string  `gorm:"column:victim_case_form1_victim_aggressor_address"                  json:"aggressorAddress"`
	VictimCaseForm1VictimAggressorPhone          string  `gorm:"column:victim_case_form1_victim_aggressor_phone"                    json:"aggressorPhone"`
	VictimCaseForm1VictimEMail                   string  `gorm:"column:victim_case_form1_victim_e_mail"                             json:"victimEmail"`
	VictimCaseForm1Age                           *int    `gorm:"column:victim_case_form1_age"                                       json:"age"`
	VictimCaseForm1VictimChildrenNumber          *int    `gorm:"column:victim_case_form1_victim_children_number"                    json:"childrenNumber"`
	VictimCaseForm1VictimNationality             string  `gorm:"column:victim_case_form1_victim_nationality"                        json:"nationality"`
	VictimCaseForm1VictimNationalityOther        string  `gorm:"column:victim_case_form1_victim_nationality_other"                  json:"nationalityOther"`
	VictimCaseForm1VictimForeignerImmigrationStatus string `gorm:"column:victim_case_form1_victim_foreigner_immigration_status"     json:"foreignerImmigrationStatus"`
	VictimCaseForm1VictimGender                  string  `gorm:"column:victim_case_form1_victim_gender"                             json:"gender"`
	VictimCaseForm1VictimGenderIdentityOther     string  `gorm:"column:victim_case_form1_victim_gender_identity_other"              json:"genderIdentityOther"`
	VictimCaseForm1VictimSexualOrientationOther  string  `gorm:"column:victim_case_form1_victim_sexual_orientation_other"           json:"sexualOrientationOther"`
	VictimCaseForm1VictimDependents              string  `gorm:"column:victim_case_form1_victim_dependents"                         json:"dependents"`
	VictimCaseForm1VictimDeathThreats            string  `gorm:"column:victim_case_form1_victim_death_threats"                      json:"deathThreats"`
	VictimCaseForm1VictimAggressorHasWeapons     string  `gorm:"column:victim_case_form1_victim_aggressor_has_weapons"              json:"aggressorHasWeapons"`
	VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore string `gorm:"column:victim_case_form1_experienced_physical_or_sexual_violence_befor" json:"experiencedViolenceBefore"`
	VictimCaseForm1VictimPreviouslyReportedSituation string `gorm:"column:victim_case_form1_victim_previously_reported_situation"   json:"previouslyReportedSituation"`
	VictimCaseForm1VictimIfPreviouslyReported    string  `gorm:"column:victim_case_form1_victim_if_previously_reported"             json:"ifPreviouslyReported"`
	VictimCaseForm1VictimImminentRisk            string  `gorm:"column:victim_case_form1_victim_imminent_risk"                      json:"imminentRisk"`
	VictimCaseForm1VictimViolenceTownCode        string  `gorm:"column:victim_case_form1_victim_violence_town_code"                 json:"violenceTownCode"`
	VictimCaseForm1VictimIfAfro                  string  `gorm:"column:victim_case_form1_victim_if_afro"                            json:"ifAfro"`
	VictimCaseForm1VictimIfIndigenous            string  `gorm:"column:victim_case_form1_victim_if_indigenous"                      json:"ifIndigenous"`
	VictimCaseForm1VictimIfIndigenousTongue      string  `gorm:"column:victim_case_form1_victim_if_indigenous_tongue"               json:"ifIndigenousTongue"`
	VictimCaseForm1VictimIfPeasant               string  `gorm:"column:victim_case_form1_victim_if_peasant"                         json:"ifPeasant"`
	VictimCaseForm1VictimIfArmedConflict         string  `gorm:"column:victim_case_form1_victim_if_armed_conflict"                  json:"ifArmedConflict"`
	VictimCaseForm1VictimViolenceScene           string  `gorm:"column:victim_case_form1_victim_violence_scene"                     json:"violenceScene"`
	VictimCaseForm1PhysicalViolenceIncreased     string  `gorm:"column:victim_case_form1_physical_violence_increased"               json:"physicalViolenceIncreased"`
	VictimCaseForm1SeparatedFromPartnerLastYear  string  `gorm:"column:victim_case_form1_separated_from_partner_last_year"          json:"separatedFromPartnerLastYear"`
	VictimCaseForm1ThreatenedWithWeapon          string  `gorm:"column:victim_case_form1_threatened_with_weapon"                    json:"threatenedWithWeapon"`
	VictimCaseForm1ThreatenedToKillOrHarmChildren string `gorm:"column:victim_case_form1_threatened_to_kill_or_harm_children"      json:"threatenedToKillOrHarmChildren"`
	VictimCaseForm1JealousAndViolent             string  `gorm:"column:victim_case_form1_jealous_and_violent"                       json:"jealousAndViolent"`
	VictimCaseForm1BelievesCapableOfKilling      string  `gorm:"column:victim_case_form1_believes_capable_of_killing"               json:"believesCapableOfKilling"`
	VictimCaseForm1VictimFemicideRisk            string  `gorm:"column:victim_case_form1_victim_femicide_risk"                      json:"femicideRisk"`
}

func (VictimCaseForm1) TableName() string { return "salvia.victim_case_form1" }

// VictimCaseForm2 mapea salvia.victim_case_form_2 (solo campo de riesgo).
type VictimCaseForm2 struct {
	VictimCaseForm2VictimCase int64 `gorm:"column:victim_case_form2_victim_case;primaryKey" json:"victimCaseForm2VictimCase"`
	VictimCaseForm2RiskLevel  *int  `gorm:"column:victim_case_form2_risk_level"             json:"riskLevel"`
}

func (VictimCaseForm2) TableName() string { return "salvia.victim_case_form2" }

// FollowUp mapea salvia.follow_up.
type FollowUp struct {
	FollowUpId             int64   `gorm:"column:follow_up_id;primaryKey"        json:"followUpId"`
	FollowUpStatus         string  `gorm:"column:follow_up_status"               json:"followUpStatus"`
	FollowUpCreatedAt      time.Time `gorm:"column:follow_up_created_at"         json:"followUpCreatedAt"`
}

func (FollowUp) TableName() string { return "salvia.follow_up" }

// FollowUpEntry mapea salvia.follow_up_entry.
type FollowUpEntry struct {
	FollowUpEntryId               int64   `gorm:"column:follow_up_entry_id;primaryKey"          json:"followUpEntryId"`
	FollowUpId                    int64   `gorm:"column:follow_up_id"                           json:"followUpId"`
	FollowUpEntrySectorCode       string  `gorm:"column:follow_up_entry_sector_code"            json:"followUpEntrySectorCode"`
	FollowUpEntryStatus           string  `gorm:"column:follow_up_entry_status"                 json:"followUpEntryStatus"`
	FollowUpEntryCompletionDate   *time.Time `gorm:"column:follow_up_entry_completion_date"     json:"followUpEntryCompletionDate"`
}

func (FollowUpEntry) TableName() string { return "salvia.follow_up_entry" }
