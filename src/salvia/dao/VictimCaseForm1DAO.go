package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	security_daos "bitsflow/security/dao"

	"encoding/json"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	VictimCaseForm1EntityName string = "VictimCaseForm1"
	VictimCaseForm1JSONName   string = "form"
	VictimCaseForm1DBName     string = "victim_case_form1"
	VictimCaseForm1DBScheme   string = "salvia"

	//Atributos relacionados con las validaciones ------------------------------

	//Campos que vienen como string del JSON. El booleano indica si son strings en el modelo o no (como en el caso de una fecha)
	VictimCaseForm1FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimCaseForm1Id":                              {Name: "VictimCaseForm1Id", DBName: "victim_case_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm1ICode":                           {Name: "VictimCaseForm1ICode", DBName: "victim_case_form1_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimCaseForm1CreationDate":                    {Name: "VictimCaseForm1CreationDate", DBName: "victim_case_form1_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm1UpdateDate":                      {Name: "VictimCaseForm1UpdateDate", DBName: "victim_case_form1_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm1Nick":                            {Name: "VictimCaseForm1Nick", DBName: "victim_case_form1_victim_nick", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm1BirthDate":                       {Name: "VictimCaseForm1BirthDate", DBName: "victim_case_form1_victim_birth_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm1Address":                         {Name: "VictimCaseForm1Address", DBName: "victim_case_form1_victim_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 120, Required: true},
		"VictimCaseForm1LivingLatitude":                  {Name: "VictimCaseForm1LivingLatitude", DBName: "victim_case_form1_victim_living_latitude", Alias: "", ModelType: "float", MinSize: -90, MaxSize: 90, Required: false},
		"VictimCaseForm1LivingLongitude":                 {Name: "VictimCaseForm1LivingLongitude", DBName: "victim_case_form1_victim_living_longitude", Alias: "", ModelType: "float", MinSize: -180, MaxSize: 180, Required: false},
		"VictimCaseForm1Phone":                           {Name: "VictimCaseForm1Phone", DBName: "victim_case_form1_victim_phone", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 10, Required: true},
		"VictimCaseForm1Email":                           {Name: "VictimCaseForm1Email", DBName: "victim_case_form1_victim_e_mail", Alias: "", ModelType: "email", MinSize: 5, MaxSize: 128, Required: false},
		"VictimCaseForm1GenderIdentity":                  {Name: "VictimCaseForm1GenderIdentity", DBName: "victim_case_form1_victim_gender_identity", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1SexualOrientation":               {Name: "VictimCaseForm1SexualOrientation", DBName: "victim_case_form1_victim_sexual_orientation", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1Origin":                          {Name: "VictimCaseForm1Origin", DBName: "victim_case_form1_victim_origin", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1Occupation":                      {Name: "VictimCaseForm1Occupation", DBName: "victim_case_form1_victim_occupation", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1OccupationOther":                 {Name: "VictimCaseForm1OccupationOther", DBName: "victim_case_form1_victim_occupation_other", Alias: "", ModelType: "string", MinSize: 4, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimEthnicGroup":               {Name: "VictimCaseForm1VictimEthnicGroup", DBName: "victim_case_form1_victim_ethnic_group", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1VictimEthnicGroupOther":          {Name: "VictimCaseForm1VictimEthnicGroupOther", DBName: "victim_case_form1_victim_ethnic_group_other", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimContactNames":              {Name: "VictimCaseForm1VictimContactNames", DBName: "victim_case_form1_victim_contact_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"VictimCaseForm1VictimContactPhone":              {Name: "VictimCaseForm1VictimContactPhone", DBName: "victim_case_form1_victim_contact_phone", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 10, Required: true},
		"VictimCaseForm1VictimContactKinship":            {Name: "VictimCaseForm1VictimContactKinship", DBName: "victim_case_form1_victim_contact_kinship", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"VictimCaseForm1VictimNumChildren":               {Name: "VictimCaseForm1VictimNumChildren", DBName: "victim_case_form1_victim_children_number", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimMaritalStatus":             {Name: "VictimCaseForm1VictimMaritalStatus", DBName: "victim_case_form1_victim_marital_status", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1VictimMaritalStatusOther":        {Name: "VictimCaseForm1VictimMaritalStatusOther", DBName: "victim_case_form1_victim_marital_status_other", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimChildrenAge":               {Name: "VictimCaseForm1VictimChildrenAge", DBName: "victim_case_form1_victim_children_age", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimDisability":                {Name: "VictimCaseForm1VictimDisability", DBName: "victim_case_form1_victim_disability", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1VictimSpecialSupport":            {Name: "VictimCaseForm1VictimSpecialSupport", DBName: "victim_case_form1_victim_special_support", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 128, Required: false},
		"VictimCaseForm1FactsOccurrence":                 {Name: "VictimCaseForm1FactsOccurrence", DBName: "victim_case_form1_facts_occurrence", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1FactsStartTime":                  {Name: "VictimCaseForm1FactsStartTime", DBName: "victim_case_form1_facts_start_time", Alias: "", ModelType: "time", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm1FactsEndTime":                    {Name: "VictimCaseForm1FactsEndTime", DBName: "victim_case_form1_facts_end_time", Alias: "", ModelType: "time", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm1FactsWeekday":                    {Name: "VictimCaseForm1FactsWeekday", DBName: "victim_case_form1_facts_weekday", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 7, Required: false},
		"VictimCaseForm1FactsDate":                       {Name: "VictimCaseForm1FactsDate", DBName: "victim_case_form1_facts_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm1FactsDescription":                {Name: "VictimCaseForm1FactsDescription", DBName: "victim_case_form1_facts_description", Alias: "", ModelType: "string", MinSize: 16, MaxSize: 20000, Required: true},
		"VictimCaseForm1VictimViolenceExperienced":       {Name: "VictimCaseForm1VictimViolenceExperienced", DBName: "victim_case_form1_victim_violence_experienced", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1VictimViolenceExperiencedOther":  {Name: "VictimCaseForm1VictimViolenceExperiencedOther", DBName: "victim_case_form1_victim_violence_experienced_other", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimViolenceScope":             {Name: "VictimCaseVictimViolenceScope", DBName: "victim_case_form1_victim_violence_scope", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseForm1VictimFemicideRisk":              {Name: "VictimCaseForm1VictimFemicideRisk", DBName: "victim_case_form1_victim_femicide_risk", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimAggressor":                 {Name: "VictimCaseForm1VictimAggressor", DBName: "victim_case_form1_victim_aggressor", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimRelationshipWithAggressor": {Name: "VictimCaseForm1VictimRelationshipWithAggressor", DBName: "victim_case_form1_victim_relationship_with_aggressor", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimCaseForm1VictimAggressorName":             {Name: "VictimCaseForm1VictimAggressorName", DBName: "victim_case_form1_victim_aggressor_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: false},
		"VictimCaseForm1VictimAggressorDocType":          {Name: "VictimCaseForm1VictimAggressorDocType", DBName: "victim_case_form1_victim_aggressor_doc_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimCaseForm1VictimAggressorDocNumber":        {Name: "VictimCaseVictimAggressorDocNumber", DBName: "victim_case_form1_victim_aggressor_doc_number", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimAggressorAddress":          {Name: "VictimCaseForm1VictimAggressorAddress", DBName: "victim_case_form1_victim_aggressor_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 128, Required: false},
		"VictimCaseForm1VictimAggressorPhone":            {Name: "VictimCaseForm1VictimAggressorPhone", DBName: "victim_case_form1_victim_aggressor_phone", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 10, Required: false},

		"VictimCaseForm1Age":                                       {Name: "VictimCaseAge", DBName: "victim_case_form1_age", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm1VictimNationality":                         {Name: "VictimCaseVictimNationality", DBName: "victim_case_form1_victim_nationality", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimNationalityOther":                    {Name: "VictimCaseVictimNationalityOther", DBName: "victim_case_form1_victim_nationality_other", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimForeignerImmigrationStatus":          {Name: "VictimCaseVictimForeignerImmigrationStatus", DBName: "victim_case_form1_victim_foreigner_immigration_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimCaseForm1VictimGender":                              {Name: "VictimCaseVictimGender", DBName: "victim_case_form1_victim_gender", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimGenderIdentityOther":                 {Name: "VictimCaseVictimGenderIdentityOther", DBName: "victim_case_form1_victim_gender_identity_other", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimSexualOrientationOther":              {Name: "VictimCaseVictimSexualOrientationOther", DBName: "victim_case_form1_victim_sexual_orientation_other", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 32, Required: false},
		"VictimCaseForm1VictimDependents":                          {Name: "VictimCaseVictimDependents", DBName: "victim_case_form1_victim_dependents", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimDeathThreats":                        {Name: "VictimCaseVictimDeathThreats", DBName: "victim_case_form1_victim_death_threats", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimCaseForm1VictimAggressorHasWeapons":                 {Name: "VictimCaseVictimAggressorHasWeapons", DBName: "victim_case_form1_victim_aggressor_has_weapons", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore": {Name: "VictimCaseExperiencedPhysicalOrSexualViolenceBefore", DBName: "victim_case_form1_experienced_physical_or_sexual_violence_befor", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimCaseForm1VictimImminentRisk":                        {Name: "VictimCaseVictimImminentRisk", DBName: "victim_case_form1_victim_imminent_risk", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimCaseForm1VictimPreviouslyReportedSituation":         {Name: "VictimCaseVictimPreviouslyReportedSituation", DBName: "victim_case_form1_victim_previously_reported_situation", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimIfPreviouslyReported":                {Name: "VictimCaseVictimIfPreviouslyReported", DBName: "victim_case_form1_victim_if_previously_reported", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimCaseForm1ViolenceTownCode":                          {Name: "VictimCaseForm1ViolenceTownCode", DBName: "victim_case_form1_violence_town_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 8, Required: true},

		"VictimCaseForm1VictimIfAfro":             {Name: "VictimCaseForm1VictimIfAfro", DBName: "victim_case_form1_victim_if_afro", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimCaseForm1VictimIfIndigenous":       {Name: "VictimCaseForm1VictimIfIndigenous", DBName: "victim_case_form1_victim_if_indigenous", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimCaseForm1VictimIfIndigenousTongue": {Name: "VictimCaseForm1VictimIfIndigenousTongue", DBName: "victim_case_form1_victim_if_indigenous_tongue", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimCaseForm1VictimIfPeasant":          {Name: "VictimCaseForm1VictimIfPeasant", DBName: "victim_case_form1_victim_if_peasant", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimIfArmedConflict":    {Name: "VictimCaseForm1VictimIfArmedConflict", DBName: "victim_case_form1_victim_if_armed_conflict", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1VictimViolenceScene":      {Name: "VictimCaseForm1VictimViolenceScene", DBName: "victim_case_form1_victim_violence_scene", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},

		"VictimCaseForm1PhysicalViolenceIncreased":      {Name: "VictimCaseForm1PhysicalViolenceIncreased", DBName: "victim_case_form1_physical_violence_increased", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1SeparatedFromPartnerLastYear":   {Name: "VictimCaseForm1SeparatedFromPartnerLastYear", DBName: "victim_case_form1_separated_from_partner_last_year", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1ThreatenedWithWeapon":           {Name: "VictimCaseForm1ThreatenedWithWeapon", DBName: "victim_case_form1_threatened_with_weapon", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1ThreatenedToKillOrHarmChildren": {Name: "VictimCaseForm1ThreatenedToKillOrHarmChildren", DBName: "victim_case_form1_threatened_to_kill_or_harm_children", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1JealousAndViolent":              {Name: "VictimCaseForm1JealousAndViolent", DBName: "victim_case_form1_jealous_and_violent", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"VictimCaseForm1BelievesCapableOfKilling":       {Name: "VictimCaseForm1BelievesCapableOfKilling", DBName: "victim_case_form1_believes_capable_of_killing", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},

		"VictimCaseForm1VictimCase": {Name: "VictimCaseForm1VictimCase", DBName: "victim_case_form1_victim_case", Alias: "", ModelType: "uint", Required: true},
	}
)

type VictimCaseForm1DTO struct {
	VictimCaseForm1Id                              uint64    `json:"-"`
	VictimCaseForm1ICode                           string    `json:"icode"`
	VictimCaseForm1CreationDate                    time.Time `json:"creationDate"`
	VictimCaseForm1UpdateDate                      time.Time `json:"updateDate"`
	VictimCaseForm1Nick                            string    `json:"nick"`
	VictimCaseForm1BirthDate                       time.Time `json:"birthDate"`
	VictimCaseForm1ViolenceTownCode                string    `json:"violenceTownCode"`
	VictimCaseForm1Address                         string    `json:"address"`
	VictimCaseForm1LivingLatitude                  float64   `json:"livingLatitude"`
	VictimCaseForm1LivingLongitude                 float64   `json:"livingLongitude"`
	VictimCaseForm1Phone                           string    `json:"phone"`
	VictimCaseForm1Email                           string    `json:"email"`
	VictimCaseForm1GenderIdentity                  string    `json:"genderIdentity"`
	VictimCaseForm1SexualOrientation               string    `json:"sexualOrientation"`
	VictimCaseForm1Origin                          string    `json:"origin"`
	VictimCaseForm1Occupation                      string    `json:"occupation"`
	VictimCaseForm1OccupationOther                 string    `json:"occupationOther"`
	VictimCaseForm1VictimEthnicGroup               string    `json:"ethnicGroup"`
	VictimCaseForm1VictimEthnicGroupOther          string    `json:"ethnicGroupOther"`
	VictimCaseForm1VictimContactNames              string    `json:"contactNames"`
	VictimCaseForm1VictimContactPhone              string    `json:"contactPhone"`
	VictimCaseForm1VictimContactKinship            string    `json:"contactKinship"`
	VictimCaseForm1VictimNumChildren               int16     `json:"numChildren"`
	VictimCaseForm1VictimMaritalStatus             string    `json:"maritalStatus"`
	VictimCaseForm1VictimMaritalStatusOther        string    `json:"maritalStatusOther"`
	VictimCaseForm1VictimChildrenAge               string    `json:"childrenAge"`
	VictimCaseForm1VictimDisability                string    `json:"disability"`
	VictimCaseForm1VictimSpecialSupport            string    `json:"specialSupport"`
	VictimCaseForm1FactsOccurrence                 string    `json:"factsOccurrence"`
	VictimCaseForm1FactsStartTime                  time.Time `json:"factsStartTime"`
	VictimCaseForm1FactsEndTime                    time.Time `json:"factsEndTime"`
	VictimCaseForm1FactsWeekday                    uint16    `json:"factsWeekday"`
	VictimCaseForm1FactsDate                       time.Time `json:"factsDate"`
	VictimCaseForm1FactsDescription                string    `json:"factsDescription"`
	VictimCaseForm1VictimViolenceExperienced       string    `json:"violenceExperienced"`
	VictimCaseForm1VictimViolenceExperiencedOther  string    `json:"violenceExperiencedOther"`
	VictimCaseForm1VictimViolenceScope             string    `json:"violenceScope"`
	VictimCaseForm1VictimFemicideRisk              string    `json:"femicideRisk"`
	VictimCaseForm1VictimAggressor                 string    `json:"aggressor"`
	VictimCaseForm1VictimRelationshipWithAggressor string    `json:"relationshipWithAggressor"`
	VictimCaseForm1VictimAggressorName             string    `json:"aggressorName"`
	VictimCaseForm1VictimAggressorDocType          string    `json:"aggressorDocType"`
	VictimCaseForm1VictimAggressorDocNumber        string    `json:"aggressorDocNumber"`
	VictimCaseForm1VictimAggressorAddress          string    `json:"aggressorAddress"`
	VictimCaseForm1VictimAggressorPhone            string    `json:"aggressorPhone"`

	VictimCaseForm1Age                                       uint16 `json:"age"`
	VictimCaseForm1VictimNationality                         string `json:"nationality"`
	VictimCaseForm1VictimNationalityOther                    string `json:"nationalityOther"`
	VictimCaseForm1VictimForeignerImmigrationStatus          string `json:"foreignImmigrationStatus"`
	VictimCaseForm1VictimGender                              string `json:"gender"`
	VictimCaseForm1VictimGenderIdentityOther                 string `json:"genderIdentityOther"`
	VictimCaseForm1VictimSexualOrientationOther              string `json:"sexualOrientationOther"`
	VictimCaseForm1VictimDependents                          string `json:"dependents"`
	VictimCaseForm1VictimDeathThreats                        string `json:"deathThreats"`
	VictimCaseForm1VictimAggressorHasWeapons                 string `json:"aggressorHasWeapons"`
	VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore string `json:"experiencedViolenceBefore"`
	VictimCaseForm1VictimImminentRisk                        string `json:"imminentRisk"`
	VictimCaseForm1VictimPreviouslyReportedSituation         string `json:"previouslyReportedSituation"`
	VictimCaseForm1VictimIfPreviouslyReported                string `json:"ifPreviouslyReported"`

	VictimCaseForm1VictimIfAfro             string `json:"ifAfro"`
	VictimCaseForm1VictimIfIndigenous       string `json:"ifIndigenous"`
	VictimCaseForm1VictimIfIndigenousTongue string `json:"ifIndigenousTongue"`
	VictimCaseForm1VictimIfPeasant          string `json:"ifPeasant"`
	VictimCaseForm1VictimIfArmedConflict    string `json:"ifArmedConflict"`
	VictimCaseForm1VictimViolenceScene      string `json:"violenceScene"`

	VictimCaseForm1PhysicalViolenceIncreased      string `json:"physicalViolenceIncreased"`
	VictimCaseForm1SeparatedFromPartnerLastYear   string `json:"separatedFromPartnerLastYear"`
	VictimCaseForm1ThreatenedWithWeapon           string `json:"threatenedWithWeapon"`
	VictimCaseForm1ThreatenedToKillOrHarmChildren string `json:"threatenedToKillOrHarmChildren"`
	VictimCaseForm1JealousAndViolent              string `json:"jealousAndViolent"`
	VictimCaseForm1BelievesCapableOfKilling       string `json:"believesCapableOfKilling"`

	VictimCaseForm1VictimCase interface{} `json:"-"` //En lugar de VictimCaseDTO porque al hacerlo se crea un ciclo ya que VictimCaseDTO usa esta clase

	//Campos de formulario que no hacen parte del modelo o no directamente en la BD -----------------------------------

	VictimCaseForm1ViolenceDepartment security_daos.DepartmentDTO `json:"violenceDepartment"`
	VictimCaseForm1ViolenceCity       security_daos.CityDTO       `json:"violenceCity"`
	VictimCaseForm1ViolenceTown       security_daos.TownDTO       `json:"violenceTown"`
}
type VictimCaseForm1PgDB struct {
	VictimCaseForm1Id                sql.NullInt64
	VictimCaseForm1ICode             sql.NullString
	VictimCaseForm1CreationDate      sql.NullTime
	VictimCaseForm1UpdateDate        sql.NullTime
	VictimCaseForm1Nick              sql.NullString
	VictimCaseForm1BirthDate         sql.NullTime
	VictimCaseForm1ViolenceTownCode  sql.NullString
	VictimCaseForm1Address           sql.NullString
	VictimCaseForm1LivingLatitude    sql.NullFloat64
	VictimCaseForm1LivingLongitude   sql.NullFloat64
	VictimCaseForm1Phone             sql.NullString
	VictimCaseForm1Email             sql.NullString
	VictimCaseForm1GenderIdentity    sql.NullString
	VictimCaseForm1SexualOrientation sql.NullString
	VictimCaseForm1Origin            sql.NullString
	VictimCaseForm1Occupation        sql.NullString
	VictimCaseForm1OccupationOther   sql.NullString

	VictimCaseForm1VictimEthnicGroup               sql.NullString
	VictimCaseForm1VictimEthnicGroupOther          sql.NullString
	VictimCaseForm1VictimContactNames              sql.NullString
	VictimCaseForm1VictimContactPhone              sql.NullString
	VictimCaseForm1VictimContactKinship            sql.NullString
	VictimCaseForm1VictimNumChildren               sql.NullInt16
	VictimCaseForm1VictimMaritalStatus             sql.NullString
	VictimCaseForm1VictimMaritalStatusOther        sql.NullString
	VictimCaseForm1VictimChildrenAge               sql.NullString
	VictimCaseForm1VictimDisability                sql.NullString
	VictimCaseForm1VictimSpecialSupport            sql.NullString
	VictimCaseForm1FactsOccurrence                 sql.NullString
	VictimCaseForm1FactsStartTime                  sql.NullString
	VictimCaseForm1FactsEndTime                    sql.NullString
	VictimCaseForm1FactsWeekday                    sql.NullInt16
	VictimCaseForm1FactsDate                       sql.NullTime
	VictimCaseForm1FactsDescription                sql.NullString
	VictimCaseForm1VictimViolenceExperienced       sql.NullString
	VictimCaseForm1VictimViolenceExperiencedOther  sql.NullString
	VictimCaseForm1VictimViolenceScope             sql.NullString
	VictimCaseForm1VictimFemicideRisk              sql.NullString
	VictimCaseForm1VictimAggressor                 sql.NullString
	VictimCaseForm1VictimRelationshipWithAggressor sql.NullString
	VictimCaseForm1VictimAggressorName             sql.NullString
	VictimCaseForm1VictimAggressorDocType          sql.NullString
	VictimCaseForm1VictimAggressorDocNumber        sql.NullString
	VictimCaseForm1VictimAggressorAddress          sql.NullString
	VictimCaseForm1VictimAggressorPhone            sql.NullString

	VictimCaseForm1Age                                       sql.NullInt16
	VictimCaseForm1VictimNationality                         sql.NullString
	VictimCaseForm1VictimNationalityOther                    sql.NullString
	VictimCaseForm1VictimForeignerImmigrationStatus          sql.NullString
	VictimCaseForm1VictimGender                              sql.NullString
	VictimCaseForm1VictimGenderIdentityOther                 sql.NullString
	VictimCaseForm1VictimSexualOrientationOther              sql.NullString
	VictimCaseForm1VictimDependents                          sql.NullString
	VictimCaseForm1VictimDeathThreats                        sql.NullString
	VictimCaseForm1VictimAggressorHasWeapons                 sql.NullString
	VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore sql.NullString
	VictimCaseForm1VictimImminentRisk                        sql.NullString
	VictimCaseForm1VictimPreviouslyReportedSituation         sql.NullString
	VictimCaseForm1VictimIfPreviouslyReported                sql.NullString

	VictimCaseForm1VictimIfAfro             sql.NullString
	VictimCaseForm1VictimIfIndigenous       sql.NullString
	VictimCaseForm1VictimIfIndigenousTongue sql.NullString
	VictimCaseForm1VictimIfPeasant          sql.NullString
	VictimCaseForm1VictimIfArmedConflict    sql.NullString
	VictimCaseForm1VictimViolenceScene      sql.NullString

	VictimCaseForm1PhysicalViolenceIncreased      sql.NullString
	VictimCaseForm1SeparatedFromPartnerLastYear   sql.NullString
	VictimCaseForm1ThreatenedWithWeapon           sql.NullString
	VictimCaseForm1ThreatenedToKillOrHarmChildren sql.NullString
	VictimCaseForm1JealousAndViolent              sql.NullString
	VictimCaseForm1BelievesCapableOfKilling       sql.NullString

	VictimCaseForm1VictimCase sql.NullInt64
}

func (vcd VictimCaseForm1DTO) MarshalJSON() ([]byte, error) {
	type Alias VictimCaseForm1DTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		VictimCaseForm1BirthDate      string `json:"birthDate"`
		VictimCaseForm1CreationDate   string `json:"creationDate"`
		VictimCaseForm1UpdateDate     string `json:"updateDate"`
		VictimCaseForm1FactsStartTime string `json:"factsStartTime"`
		VictimCaseForm1FactsEndTime   string `json:"factsEndTime"`
		VictimCaseForm1FactsDate      string `json:"factsDate"`
	}{
		Alias:                         (*Alias)(&vcd),
		VictimCaseForm1BirthDate:      vcd.VictimCaseForm1BirthDate.Format(common_config.DateTime.DATE_FORMAT),
		VictimCaseForm1CreationDate:   vcd.VictimCaseForm1CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimCaseForm1UpdateDate:     vcd.VictimCaseForm1UpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimCaseForm1FactsStartTime: vcd.VictimCaseForm1FactsStartTime.Format(common_config.DateTime.TIME_FORMAT),
		VictimCaseForm1FactsEndTime:   vcd.VictimCaseForm1FactsEndTime.Format(common_config.DateTime.TIME_FORMAT),
		VictimCaseForm1FactsDate:      vcd.VictimCaseForm1FactsDate.Format(common_config.DateTime.DATE_FORMAT),
	})
}

func (vcd *VictimCaseForm1DTO) UnmarshalJSON(data []byte) error {
	type Alias VictimCaseForm1DTO

	// Estructura auxiliar donde las fechas vienen como string
	aux := &struct {
		*Alias
		VictimCaseForm1BirthDate      string `json:"birthDate"`
		VictimCaseForm1CreationDate   string `json:"creationDate"`
		VictimCaseForm1UpdateDate     string `json:"updateDate"`
		VictimCaseForm1FactsStartTime string `json:"factsStartTime"`
		VictimCaseForm1FactsEndTime   string `json:"factsEndTime"`
		VictimCaseForm1FactsDate      string `json:"factsDate"`
	}{
		Alias: (*Alias)(vcd),
	}

	// Primero unmarshalea todo el JSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Función helper para parsear fechas sin romper el flujo
	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{} // Fecha inválida → zero value
		}
		return t
	}

	// Parseo tolerante
	vcd.VictimCaseForm1BirthDate = parse(aux.VictimCaseForm1BirthDate, common_config.DateTime.DATE_FORMAT)
	vcd.VictimCaseForm1CreationDate = parse(aux.VictimCaseForm1CreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimCaseForm1UpdateDate = parse(aux.VictimCaseForm1UpdateDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimCaseForm1FactsStartTime = parse(aux.VictimCaseForm1FactsStartTime, common_config.DateTime.TIME_FORMAT)
	vcd.VictimCaseForm1FactsEndTime = parse(aux.VictimCaseForm1FactsEndTime, common_config.DateTime.TIME_FORMAT)
	vcd.VictimCaseForm1FactsDate = parse(aux.VictimCaseForm1FactsDate, common_config.DateTime.DATE_FORMAT)
	return nil
}

func SetVictimCaseForm1(victimCaseForm1 *VictimCaseForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query

	var victimCaseForm1FieldsSlice []string = []string{"VictimCaseForm1ICode", "VictimCaseForm1CreationDate", "VictimCaseForm1UpdateDate",
		"VictimCaseForm1Nick", "VictimCaseForm1BirthDate", "VictimCaseForm1ViolenceTownCode", "VictimCaseForm1Address",
		"VictimCaseForm1LivingLatitude", "VictimCaseForm1LivingLongitude", "VictimCaseForm1Phone", "VictimCaseForm1Email", "VictimCaseForm1GenderIdentity", "VictimCaseForm1SexualOrientation", "VictimCaseForm1Origin",
		"VictimCaseForm1Occupation", "VictimCaseForm1OccupationOther", "VictimCaseForm1VictimEthnicGroup", "VictimCaseForm1VictimEthnicGroupOther", "VictimCaseForm1VictimContactNames",
		"VictimCaseForm1VictimContactPhone", "VictimCaseForm1VictimContactKinship", "VictimCaseForm1VictimNumChildren", "VictimCaseForm1VictimMaritalStatus",
		"VictimCaseForm1VictimMaritalStatusOther", "VictimCaseForm1VictimChildrenAge", "VictimCaseForm1VictimDisability", "VictimCaseForm1VictimSpecialSupport", "VictimCaseForm1FactsOccurrence",
		"VictimCaseForm1FactsStartTime", "VictimCaseForm1FactsEndTime", "VictimCaseForm1FactsWeekday", "VictimCaseForm1FactsDate", "VictimCaseForm1FactsDescription",
		"VictimCaseForm1VictimViolenceExperienced", "VictimCaseForm1VictimViolenceExperiencedOther", "VictimCaseForm1VictimViolenceScope", "VictimCaseForm1VictimFemicideRisk",
		"VictimCaseForm1VictimAggressor", "VictimCaseForm1VictimRelationshipWithAggressor", "VictimCaseForm1VictimAggressorName", "VictimCaseForm1VictimAggressorDocType",
		"VictimCaseForm1VictimAggressorDocNumber", "VictimCaseForm1VictimAggressorAddress", "VictimCaseForm1VictimAggressorPhone",
		"VictimCaseForm1Age", "VictimCaseForm1VictimNationality", "VictimCaseForm1VictimNationalityOther", "VictimCaseForm1VictimForeignerImmigrationStatus",
		"VictimCaseForm1VictimGender", "VictimCaseForm1VictimGenderIdentityOther", "VictimCaseForm1VictimSexualOrientationOther", "VictimCaseForm1VictimDependents",
		"VictimCaseForm1VictimDeathThreats", "VictimCaseForm1VictimAggressorHasWeapons", "VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore",
		"VictimCaseForm1VictimImminentRisk", "VictimCaseForm1VictimPreviouslyReportedSituation", "VictimCaseForm1VictimIfPreviouslyReported",
		"VictimCaseForm1VictimIfAfro", "VictimCaseForm1VictimIfIndigenous", "VictimCaseForm1VictimIfIndigenousTongue", "VictimCaseForm1VictimIfPeasant",
		"VictimCaseForm1VictimIfArmedConflict", "VictimCaseForm1VictimViolenceScene",
		"VictimCaseForm1PhysicalViolenceIncreased", "VictimCaseForm1SeparatedFromPartnerLastYear", "VictimCaseForm1ThreatenedWithWeapon", "VictimCaseForm1ThreatenedToKillOrHarmChildren",
		"VictimCaseForm1JealousAndViolent", "VictimCaseForm1BelievesCapableOfKilling", "VictimCaseForm1VictimCase"}

	var victimCaseForm1FieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, victimCaseForm1FieldsSlice, victimCaseForm1FieldsAliasSlice, VictimCaseForm1DBName, []string{}, []string{}, []string{"VictimCaseForm1Id"}, common_dao.SQL_AND, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		victimCaseForm1.VictimCaseForm1ICode, victimCaseForm1.VictimCaseForm1CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), victimCaseForm1.VictimCaseForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCaseForm1.VictimCaseForm1Nick, victimCaseForm1.VictimCaseForm1BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		victimCaseForm1.VictimCaseForm1ViolenceTownCode, victimCaseForm1.VictimCaseForm1Address,
		victimCaseForm1.VictimCaseForm1LivingLatitude, victimCaseForm1.VictimCaseForm1LivingLongitude, victimCaseForm1.VictimCaseForm1Phone, victimCaseForm1.VictimCaseForm1Email, victimCaseForm1.VictimCaseForm1GenderIdentity, victimCaseForm1.VictimCaseForm1SexualOrientation,
		victimCaseForm1.VictimCaseForm1Origin, victimCaseForm1.VictimCaseForm1Occupation, victimCaseForm1.VictimCaseForm1OccupationOther, victimCaseForm1.VictimCaseForm1VictimEthnicGroup, victimCaseForm1.VictimCaseForm1VictimEthnicGroupOther,
		victimCaseForm1.VictimCaseForm1VictimContactNames, victimCaseForm1.VictimCaseForm1VictimContactPhone, victimCaseForm1.VictimCaseForm1VictimContactKinship, victimCaseForm1.VictimCaseForm1VictimNumChildren,
		victimCaseForm1.VictimCaseForm1VictimMaritalStatus, victimCaseForm1.VictimCaseForm1VictimMaritalStatusOther, victimCaseForm1.VictimCaseForm1VictimChildrenAge, victimCaseForm1.VictimCaseForm1VictimDisability,
		victimCaseForm1.VictimCaseForm1VictimSpecialSupport, victimCaseForm1.VictimCaseForm1FactsOccurrence, victimCaseForm1.VictimCaseForm1FactsStartTime.Format(common_config.DateTime.TIME_FORMAT),
		victimCaseForm1.VictimCaseForm1FactsEndTime.Format(common_config.DateTime.TIME_FORMAT), victimCaseForm1.VictimCaseForm1FactsWeekday,
		victimCaseForm1.VictimCaseForm1FactsDate.Format(common_config.DateTime.DB_DATE_FORMAT), victimCaseForm1.VictimCaseForm1FactsDescription, victimCaseForm1.VictimCaseForm1VictimViolenceExperienced,
		victimCaseForm1.VictimCaseForm1VictimViolenceExperiencedOther, victimCaseForm1.VictimCaseForm1VictimViolenceScope, victimCaseForm1.VictimCaseForm1VictimFemicideRisk,
		victimCaseForm1.VictimCaseForm1VictimAggressor, victimCaseForm1.VictimCaseForm1VictimRelationshipWithAggressor, victimCaseForm1.VictimCaseForm1VictimAggressorName,
		victimCaseForm1.VictimCaseForm1VictimAggressorDocType, victimCaseForm1.VictimCaseForm1VictimAggressorDocNumber, victimCaseForm1.VictimCaseForm1VictimAggressorAddress,
		victimCaseForm1.VictimCaseForm1VictimAggressorPhone, victimCaseForm1.VictimCaseForm1Age, victimCaseForm1.VictimCaseForm1VictimNationality, victimCaseForm1.VictimCaseForm1VictimNationalityOther,
		victimCaseForm1.VictimCaseForm1VictimForeignerImmigrationStatus, victimCaseForm1.VictimCaseForm1VictimGender, victimCaseForm1.VictimCaseForm1VictimGenderIdentityOther,
		victimCaseForm1.VictimCaseForm1VictimSexualOrientationOther, victimCaseForm1.VictimCaseForm1VictimDependents, victimCaseForm1.VictimCaseForm1VictimDeathThreats,
		victimCaseForm1.VictimCaseForm1VictimAggressorHasWeapons, victimCaseForm1.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore,
		victimCaseForm1.VictimCaseForm1VictimImminentRisk, victimCaseForm1.VictimCaseForm1VictimPreviouslyReportedSituation, victimCaseForm1.VictimCaseForm1VictimIfPreviouslyReported,
		victimCaseForm1.VictimCaseForm1VictimIfAfro, victimCaseForm1.VictimCaseForm1VictimIfIndigenous, victimCaseForm1.VictimCaseForm1VictimIfIndigenousTongue,
		victimCaseForm1.VictimCaseForm1VictimIfPeasant, victimCaseForm1.VictimCaseForm1VictimIfArmedConflict, victimCaseForm1.VictimCaseForm1VictimViolenceScene,
		victimCaseForm1.VictimCaseForm1PhysicalViolenceIncreased, victimCaseForm1.VictimCaseForm1SeparatedFromPartnerLastYear, victimCaseForm1.VictimCaseForm1ThreatenedWithWeapon,
		victimCaseForm1.VictimCaseForm1ThreatenedToKillOrHarmChildren, victimCaseForm1.VictimCaseForm1JealousAndViolent, victimCaseForm1.VictimCaseForm1BelievesCapableOfKilling, victimCaseForm1.VictimCaseForm1VictimCase.(VictimCaseDTO).VictimCaseId)

	persistenceCtrl.Scan(&victimCaseForm1.VictimCaseForm1Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetVictimCaseForm1(by common_controllers.By, victimCaseForm1 *VictimCaseForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var victimCaseForm1Path string = VictimCaseForm1DBScheme + "." + VictimCaseForm1DBName

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var victimCaseForm1FieldsSlice []string = []string{"VictimCaseForm1Id", "VictimCaseForm1ICode", "VictimCaseForm1CreationDate", "VictimCaseForm1UpdateDate",
		"VictimCaseForm1Nick", "VictimCaseForm1BirthDate", "VictimCaseForm1ViolenceTownCode", "VictimCaseForm1Address",
		"VictimCaseForm1LivingLatitude", "VictimCaseForm1LivingLongitude", "VictimCaseForm1Phone", "VictimCaseForm1Email", "VictimCaseForm1GenderIdentity", "VictimCaseForm1SexualOrientation", "VictimCaseForm1Origin",
		"VictimCaseForm1Occupation", "VictimCaseForm1OccupationOther", "VictimCaseForm1VictimEthnicGroup", "VictimCaseForm1VictimEthnicGroupOther", "VictimCaseForm1VictimContactNames",
		"VictimCaseForm1VictimContactPhone", "VictimCaseForm1VictimContactKinship", "VictimCaseForm1VictimNumChildren", "VictimCaseForm1VictimMaritalStatus",
		"VictimCaseForm1VictimMaritalStatusOther", "VictimCaseForm1VictimChildrenAge", "VictimCaseForm1VictimDisability", "VictimCaseForm1VictimSpecialSupport", "VictimCaseForm1FactsOccurrence",
		"VictimCaseForm1FactsStartTime", "VictimCaseForm1FactsEndTime", "VictimCaseForm1FactsWeekday", "VictimCaseForm1FactsDate", "VictimCaseForm1FactsDescription",
		"VictimCaseForm1VictimViolenceExperienced", "VictimCaseForm1VictimViolenceExperiencedOther", "VictimCaseForm1VictimViolenceScope", "VictimCaseForm1VictimFemicideRisk",
		"VictimCaseForm1VictimAggressor", "VictimCaseForm1VictimRelationshipWithAggressor", "VictimCaseForm1VictimAggressorName", "VictimCaseForm1VictimAggressorDocType",
		"VictimCaseForm1VictimAggressorDocNumber", "VictimCaseForm1VictimAggressorAddress", "VictimCaseForm1VictimAggressorPhone",
		"VictimCaseForm1Age", "VictimCaseForm1VictimNationality", "VictimCaseForm1VictimNationalityOther", "VictimCaseForm1VictimForeignerImmigrationStatus",
		"VictimCaseForm1VictimGender", "VictimCaseForm1VictimGenderIdentityOther", "VictimCaseForm1VictimSexualOrientationOther", "VictimCaseForm1VictimDependents",
		"VictimCaseForm1VictimDeathThreats", "VictimCaseForm1VictimAggressorHasWeapons", "VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore",
		"VictimCaseForm1VictimImminentRisk", "VictimCaseForm1VictimPreviouslyReportedSituation", "VictimCaseForm1VictimIfPreviouslyReported",
		"VictimCaseForm1VictimIfAfro", "VictimCaseForm1VictimIfIndigenous", "VictimCaseForm1VictimIfIndigenousTongue", "VictimCaseForm1VictimIfPeasant",
		"VictimCaseForm1VictimIfArmedConflict", "VictimCaseForm1VictimViolenceScene",
		"VictimCaseForm1PhysicalViolenceIncreased", "VictimCaseForm1SeparatedFromPartnerLastYear", "VictimCaseForm1ThreatenedWithWeapon", "VictimCaseForm1ThreatenedToKillOrHarmChildren",
		"VictimCaseForm1JealousAndViolent", "VictimCaseForm1BelievesCapableOfKilling", "VictimCaseForm1VictimCase"}

	var victimCaseForm1FieldsAliasSlice []string = []string{}

	var victimCaseForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseForm1FieldsSlice, victimCaseForm1FieldsAliasSlice, VictimCaseForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true)

	var query string = `SELECT ` + victimCaseForm1FieldsStr +
		` FROM ` + victimCaseForm1Path +

		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true)

	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var victimCaseForm1Pg VictimCaseForm1PgDB = VictimCaseForm1PgDB{}

	persistenceCtrl.Scan(&victimCaseForm1Pg.VictimCaseForm1Id,
		&victimCaseForm1Pg.VictimCaseForm1ICode, &victimCaseForm1Pg.VictimCaseForm1CreationDate, &victimCaseForm1Pg.VictimCaseForm1UpdateDate,
		&victimCaseForm1Pg.VictimCaseForm1Nick, &victimCaseForm1Pg.VictimCaseForm1BirthDate,
		&victimCaseForm1Pg.VictimCaseForm1ViolenceTownCode, &victimCaseForm1Pg.VictimCaseForm1Address,
		&victimCaseForm1Pg.VictimCaseForm1LivingLatitude, &victimCaseForm1Pg.VictimCaseForm1LivingLongitude, &victimCaseForm1Pg.VictimCaseForm1Phone, &victimCaseForm1Pg.VictimCaseForm1Email, &victimCaseForm1Pg.VictimCaseForm1GenderIdentity, &victimCaseForm1Pg.VictimCaseForm1SexualOrientation,
		&victimCaseForm1Pg.VictimCaseForm1Origin, &victimCaseForm1Pg.VictimCaseForm1Occupation, &victimCaseForm1Pg.VictimCaseForm1OccupationOther, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroup, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroupOther,
		&victimCaseForm1Pg.VictimCaseForm1VictimContactNames, &victimCaseForm1Pg.VictimCaseForm1VictimContactPhone, &victimCaseForm1Pg.VictimCaseForm1VictimContactKinship, &victimCaseForm1Pg.VictimCaseForm1VictimNumChildren,
		&victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatus, &victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatusOther, &victimCaseForm1Pg.VictimCaseForm1VictimChildrenAge, &victimCaseForm1Pg.VictimCaseForm1VictimDisability,
		&victimCaseForm1Pg.VictimCaseForm1VictimSpecialSupport, &victimCaseForm1Pg.VictimCaseForm1FactsOccurrence, &victimCaseForm1Pg.VictimCaseForm1FactsStartTime,
		&victimCaseForm1Pg.VictimCaseForm1FactsEndTime, &victimCaseForm1Pg.VictimCaseForm1FactsWeekday,
		&victimCaseForm1Pg.VictimCaseForm1FactsDate, &victimCaseForm1Pg.VictimCaseForm1FactsDescription, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperienced,
		&victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperiencedOther, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScope, &victimCaseForm1Pg.VictimCaseForm1VictimFemicideRisk,
		&victimCaseForm1Pg.VictimCaseForm1VictimAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimRelationshipWithAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorName,
		&victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocType, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocNumber, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorAddress,
		&victimCaseForm1Pg.VictimCaseForm1VictimAggressorPhone, &victimCaseForm1Pg.VictimCaseForm1Age, &victimCaseForm1Pg.VictimCaseForm1VictimNationality, &victimCaseForm1Pg.VictimCaseForm1VictimNationalityOther,
		&victimCaseForm1Pg.VictimCaseForm1VictimForeignerImmigrationStatus, &victimCaseForm1Pg.VictimCaseForm1VictimGender, &victimCaseForm1Pg.VictimCaseForm1VictimGenderIdentityOther,
		&victimCaseForm1Pg.VictimCaseForm1VictimSexualOrientationOther, &victimCaseForm1Pg.VictimCaseForm1VictimDependents, &victimCaseForm1Pg.VictimCaseForm1VictimDeathThreats,
		&victimCaseForm1Pg.VictimCaseForm1VictimAggressorHasWeapons, &victimCaseForm1Pg.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore,
		&victimCaseForm1Pg.VictimCaseForm1VictimImminentRisk, &victimCaseForm1Pg.VictimCaseForm1VictimPreviouslyReportedSituation, &victimCaseForm1Pg.VictimCaseForm1VictimIfPreviouslyReported,
		&victimCaseForm1Pg.VictimCaseForm1VictimIfAfro, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenous, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenousTongue,
		&victimCaseForm1Pg.VictimCaseForm1VictimIfPeasant, &victimCaseForm1Pg.VictimCaseForm1VictimIfArmedConflict, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScene,
		&victimCaseForm1Pg.VictimCaseForm1PhysicalViolenceIncreased, &victimCaseForm1Pg.VictimCaseForm1SeparatedFromPartnerLastYear, &victimCaseForm1Pg.VictimCaseForm1ThreatenedWithWeapon,
		&victimCaseForm1Pg.VictimCaseForm1ThreatenedToKillOrHarmChildren, &victimCaseForm1Pg.VictimCaseForm1JealousAndViolent, &victimCaseForm1Pg.VictimCaseForm1BelievesCapableOfKilling, &victimCaseForm1Pg.VictimCaseForm1VictimCase)

	*victimCaseForm1 = victimCaseForm1Pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetVictimCasesForm1(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm1DTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm1Path string = VictimCaseForm1DBScheme + "." + VictimCaseForm1DBName

	var victimCaseForm1s []VictimCaseForm1DTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseForm1FieldsSlice []string = []string{"VictimCaseForm1ICode", "VictimCaseForm1CreationDate", "VictimCaseForm1UpdateDate",
		"VictimCaseForm1Nick", "VictimCaseForm1BirthDate", "VictimCaseForm1ViolenceTownCode", "VictimCaseForm1Address",
		"VictimCaseForm1LivingLatitude", "VictimCaseForm1LivingLongitude", "VictimCaseForm1Phone", "VictimCaseForm1Email", "VictimCaseForm1GenderIdentity", "VictimCaseForm1SexualOrientation", "VictimCaseForm1Origin",
		"VictimCaseForm1Occupation", "VictimCaseForm1OccupationOther", "VictimCaseForm1VictimEthnicGroup", "VictimCaseForm1VictimEthnicGroupOther", "VictimCaseForm1VictimContactNames",
		"VictimCaseForm1VictimContactPhone", "VictimCaseForm1VictimContactKinship", "VictimCaseForm1VictimNumChildren", "VictimCaseForm1VictimMaritalStatus",
		"VictimCaseForm1VictimMaritalStatusOther", "VictimCaseForm1VictimChildrenAge", "VictimCaseForm1VictimDisability", "VictimCaseForm1VictimSpecialSupport", "VictimCaseForm1FactsOccurrence",
		"VictimCaseForm1FactsStartTime", "VictimCaseForm1FactsEndTime", "VictimCaseForm1FactsWeekday", "VictimCaseForm1FactsDate", "VictimCaseForm1FactsDescription",
		"VictimCaseForm1VictimViolenceExperienced", "VictimCaseForm1VictimViolenceExperiencedOther", "VictimCaseForm1VictimViolenceScope", "VictimCaseForm1VictimFemicideRisk",
		"VictimCaseForm1VictimAggressor", "VictimCaseForm1VictimRelationshipWithAggressor", "VictimCaseForm1VictimAggressorName", "VictimCaseForm1VictimAggressorDocType",
		"VictimCaseForm1VictimAggressorDocNumber", "VictimCaseForm1VictimAggressorAddress", "VictimCaseForm1VictimAggressorPhone",
		"VictimCaseForm1Age", "VictimCaseForm1VictimNationality", "VictimCaseForm1VictimNationalityOther", "VictimCaseForm1VictimForeignerImmigrationStatus",
		"VictimCaseForm1VictimGender", "VictimCaseForm1VictimGenderIdentityOther", "VictimCaseForm1VictimSexualOrientationOther", "VictimCaseForm1VictimDependents",
		"VictimCaseForm1VictimDeathThreats", "VictimCaseForm1VictimAggressorHasWeapons", "VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore",
		"VictimCaseForm1VictimImminentRisk", "VictimCaseForm1VictimPreviouslyReportedSituation", "VictimCaseForm1VictimIfPreviouslyReported",
		"VictimCaseForm1VictimIfAfro", "VictimCaseForm1VictimIfIndigenous", "VictimCaseForm1VictimIfIndigenousTongue", "VictimCaseForm1VictimIfPeasant",
		"VictimCaseForm1VictimIfArmedConflict", "VictimCaseForm1VictimViolenceScene",
		"VictimCaseForm1PhysicalViolenceIncreased", "VictimCaseForm1SeparatedFromPartnerLastYear", "VictimCaseForm1ThreatenedWithWeapon", "VictimCaseForm1ThreatenedToKillOrHarmChildren",
		"VictimCaseForm1JealousAndViolent", "VictimCaseForm1BelievesCapableOfKilling", "VictimCaseForm1VictimCase"}
	var victimCaseForm1FieldsAliasSlice []string = []string{}

	var victimCaseForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseForm1FieldsSlice, victimCaseForm1FieldsAliasSlice, VictimCaseForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true)

	var query string = `SELECT ` + victimCaseForm1FieldsStr +
		` FROM ` + victimCaseForm1Path +

		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true) +
		` ORDER BY ` + victimCaseForm1Path + `.` + VictimCaseForm1FieldDefinitions["VictimCaseForm1CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	for persistenceCtrl.Next() {
		var victimCaseForm1Pg VictimCaseForm1PgDB = VictimCaseForm1PgDB{}
		persistenceCtrl.ScanRow(&victimCaseForm1Pg.VictimCaseForm1ICode, &victimCaseForm1Pg.VictimCaseForm1CreationDate, &victimCaseForm1Pg.VictimCaseForm1UpdateDate,
			&victimCaseForm1Pg.VictimCaseForm1Nick, &victimCaseForm1Pg.VictimCaseForm1BirthDate,
			&victimCaseForm1Pg.VictimCaseForm1ViolenceTownCode, &victimCaseForm1Pg.VictimCaseForm1Address,
			&victimCaseForm1Pg.VictimCaseForm1LivingLatitude, &victimCaseForm1Pg.VictimCaseForm1LivingLongitude, &victimCaseForm1Pg.VictimCaseForm1Phone, &victimCaseForm1Pg.VictimCaseForm1Email, &victimCaseForm1Pg.VictimCaseForm1GenderIdentity, &victimCaseForm1Pg.VictimCaseForm1SexualOrientation,
			&victimCaseForm1Pg.VictimCaseForm1Origin, &victimCaseForm1Pg.VictimCaseForm1Occupation, &victimCaseForm1Pg.VictimCaseForm1OccupationOther, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroup, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroupOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimContactNames, &victimCaseForm1Pg.VictimCaseForm1VictimContactPhone, &victimCaseForm1Pg.VictimCaseForm1VictimContactKinship, &victimCaseForm1Pg.VictimCaseForm1VictimNumChildren,
			&victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatus, &victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatusOther, &victimCaseForm1Pg.VictimCaseForm1VictimChildrenAge, &victimCaseForm1Pg.VictimCaseForm1VictimDisability,
			&victimCaseForm1Pg.VictimCaseForm1VictimSpecialSupport, &victimCaseForm1Pg.VictimCaseForm1FactsOccurrence, &victimCaseForm1Pg.VictimCaseForm1FactsStartTime,
			&victimCaseForm1Pg.VictimCaseForm1FactsEndTime, &victimCaseForm1Pg.VictimCaseForm1FactsWeekday,
			&victimCaseForm1Pg.VictimCaseForm1FactsDate, &victimCaseForm1Pg.VictimCaseForm1FactsDescription, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperienced,
			&victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperiencedOther, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScope, &victimCaseForm1Pg.VictimCaseForm1VictimFemicideRisk,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimRelationshipWithAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorName,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocType, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocNumber, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorAddress,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorPhone, &victimCaseForm1Pg.VictimCaseForm1Age, &victimCaseForm1Pg.VictimCaseForm1VictimNationality, &victimCaseForm1Pg.VictimCaseForm1VictimNationalityOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimForeignerImmigrationStatus, &victimCaseForm1Pg.VictimCaseForm1VictimGender, &victimCaseForm1Pg.VictimCaseForm1VictimGenderIdentityOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimSexualOrientationOther, &victimCaseForm1Pg.VictimCaseForm1VictimDependents, &victimCaseForm1Pg.VictimCaseForm1VictimDeathThreats,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorHasWeapons, &victimCaseForm1Pg.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore,
			&victimCaseForm1Pg.VictimCaseForm1VictimImminentRisk, &victimCaseForm1Pg.VictimCaseForm1VictimPreviouslyReportedSituation, &victimCaseForm1Pg.VictimCaseForm1VictimIfPreviouslyReported,
			&victimCaseForm1Pg.VictimCaseForm1VictimIfAfro, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenous, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenousTongue,
			&victimCaseForm1Pg.VictimCaseForm1VictimIfPeasant, &victimCaseForm1Pg.VictimCaseForm1VictimIfArmedConflict, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScene,
			&victimCaseForm1Pg.VictimCaseForm1PhysicalViolenceIncreased, &victimCaseForm1Pg.VictimCaseForm1SeparatedFromPartnerLastYear, &victimCaseForm1Pg.VictimCaseForm1ThreatenedWithWeapon,
			&victimCaseForm1Pg.VictimCaseForm1ThreatenedToKillOrHarmChildren, &victimCaseForm1Pg.VictimCaseForm1JealousAndViolent, &victimCaseForm1Pg.VictimCaseForm1BelievesCapableOfKilling, &victimCaseForm1Pg.VictimCaseForm1VictimCase)

		victimCaseForm1s = append(victimCaseForm1s, victimCaseForm1Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCaseForm1Path +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return victimCaseForm1s, count, nil
}

func GetAllVictimCasesForm1(victimCaseForm1Status string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm1DTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm1Path string = VictimCaseForm1DBScheme + "." + VictimCaseForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseForm1FieldsSlice []string = []string{"VictimCaseForm1ICode", "VictimCaseForm1CreationDate", "VictimCaseForm1UpdateDate",
		"VictimCaseForm1Nick", "VictimCaseForm1BirthDate", "VictimCaseForm1ViolenceTownCode", "VictimCaseForm1Address",
		"VictimCaseForm1LivingLatitude", "VictimCaseForm1LivingLongitude", "VictimCaseForm1Phone", "VictimCaseForm1Email", "VictimCaseForm1GenderIdentity", "VictimCaseForm1SexualOrientation", "VictimCaseForm1Origin",
		"VictimCaseForm1Occupation", "VictimCaseForm1OccupationOther", "VictimCaseForm1VictimEthnicGroup", "VictimCaseForm1VictimEthnicGroupOther", "VictimCaseForm1VictimContactNames",
		"VictimCaseForm1VictimContactPhone", "VictimCaseForm1VictimContactKinship", "VictimCaseForm1VictimNumChildren", "VictimCaseForm1VictimMaritalStatus",
		"VictimCaseForm1VictimMaritalStatusOther", "VictimCaseForm1VictimChildrenAge", "VictimCaseForm1VictimDisability", "VictimCaseForm1VictimSpecialSupport", "VictimCaseForm1FactsOccurrence",
		"VictimCaseForm1FactsStartTime", "VictimCaseForm1FactsEndTime", "VictimCaseForm1FactsWeekday", "VictimCaseForm1FactsDate", "VictimCaseForm1FactsDescription",
		"VictimCaseForm1VictimViolenceExperienced", "VictimCaseForm1VictimViolenceExperiencedOther", "VictimCaseForm1VictimViolenceScope", "VictimCaseForm1VictimFemicideRisk",
		"VictimCaseForm1VictimAggressor", "VictimCaseForm1VictimRelationshipWithAggressor", "VictimCaseForm1VictimAggressorName", "VictimCaseForm1VictimAggressorDocType",
		"VictimCaseForm1VictimAggressorDocNumber", "VictimCaseForm1VictimAggressorAddress", "VictimCaseForm1VictimAggressorPhone",
		"VictimCaseForm1Age", "VictimCaseForm1VictimNationality", "VictimCaseForm1VictimNationalityOther", "VictimCaseForm1VictimForeignerImmigrationStatus",
		"VictimCaseForm1VictimGender", "VictimCaseForm1VictimGenderIdentityOther", "VictimCaseForm1VictimSexualOrientationOther", "VictimCaseForm1VictimDependents",
		"VictimCaseForm1VictimDeathThreats", "VictimCaseForm1VictimAggressorHasWeapons", "VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore",
		"VictimCaseForm1VictimImminentRisk", "VictimCaseForm1VictimPreviouslyReportedSituation", "VictimCaseForm1VictimIfPreviouslyReported",
		"VictimCaseForm1VictimIfAfro", "VictimCaseForm1VictimIfIndigenous", "VictimCaseForm1VictimIfIndigenousTongue", "VictimCaseForm1VictimIfPeasant",
		"VictimCaseForm1VictimIfArmedConflict", "VictimCaseForm1VictimViolenceScene",
		"VictimCaseForm1PhysicalViolenceIncreased", "VictimCaseForm1SeparatedFromPartnerLastYear", "VictimCaseForm1ThreatenedWithWeapon", "VictimCaseForm1ThreatenedToKillOrHarmChildren",
		"VictimCaseForm1JealousAndViolent", "VictimCaseForm1BelievesCapableOfKilling", "VictimCaseForm1VictimCase"}
	var victimCaseForm1FieldsAliasSlice []string = []string{}

	/*
		var profileFieldsSlice []string = []string{"GeneralUserProfileNames", "GeneralUserProfileLastNames"}
		var profileFieldsAliasSlice []string = []string{}

		var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, security_daos.GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.GeneralUserProfileDBScheme, security_daos.GeneralUserProfileFieldDefinitions, true)
	*/

	var victimCaseForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseForm1FieldsSlice, victimCaseForm1FieldsAliasSlice, VictimCaseForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true)

	var query string = `SELECT ` + victimCaseForm1FieldsStr +
		` FROM ` + victimCaseForm1Path +
		` WHERE ` + victimCaseForm1Path + `.` + VictimCaseForm1FieldDefinitions["VictimCaseForm1Status"].DBName + ` = $1 ` +
		` ORDER BY ` + victimCaseForm1Path + `.` + VictimCaseForm1FieldDefinitions["VictimCaseForm1CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	//var query string = `SELECT victimCaseForm1_id, victimCaseForm1_i_code, victimCaseForm1_creation_date,victimCaseForm1_updatdate,victimCaseForm1_data,victimCaseForm1_password,victimCaseForm1_status,victimCaseForm1_language FROM ` + VictimCaseForm1DBScheme + `.victimCaseForm1`

	persistenceCtrl.Query(context.Background(), query, victimCaseForm1Status)
	var victimCaseForm1s []VictimCaseForm1DTO
	for persistenceCtrl.Next() {
		var victimCaseForm1Pg VictimCaseForm1PgDB = VictimCaseForm1PgDB{}

		persistenceCtrl.ScanRow(&victimCaseForm1Pg.VictimCaseForm1ICode, &victimCaseForm1Pg.VictimCaseForm1CreationDate, &victimCaseForm1Pg.VictimCaseForm1UpdateDate,
			&victimCaseForm1Pg.VictimCaseForm1Nick, &victimCaseForm1Pg.VictimCaseForm1BirthDate,
			&victimCaseForm1Pg.VictimCaseForm1ViolenceTownCode, &victimCaseForm1Pg.VictimCaseForm1Address,
			&victimCaseForm1Pg.VictimCaseForm1LivingLatitude, &victimCaseForm1Pg.VictimCaseForm1LivingLongitude, &victimCaseForm1Pg.VictimCaseForm1Phone, &victimCaseForm1Pg.VictimCaseForm1Email, &victimCaseForm1Pg.VictimCaseForm1GenderIdentity, &victimCaseForm1Pg.VictimCaseForm1SexualOrientation,
			&victimCaseForm1Pg.VictimCaseForm1Origin, &victimCaseForm1Pg.VictimCaseForm1Occupation, &victimCaseForm1Pg.VictimCaseForm1OccupationOther, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroup, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroupOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimContactNames, &victimCaseForm1Pg.VictimCaseForm1VictimContactPhone, &victimCaseForm1Pg.VictimCaseForm1VictimContactKinship, &victimCaseForm1Pg.VictimCaseForm1VictimNumChildren,
			&victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatus, &victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatusOther, &victimCaseForm1Pg.VictimCaseForm1VictimChildrenAge, &victimCaseForm1Pg.VictimCaseForm1VictimDisability,
			&victimCaseForm1Pg.VictimCaseForm1VictimSpecialSupport, &victimCaseForm1Pg.VictimCaseForm1FactsOccurrence, &victimCaseForm1Pg.VictimCaseForm1FactsStartTime,
			&victimCaseForm1Pg.VictimCaseForm1FactsEndTime, &victimCaseForm1Pg.VictimCaseForm1FactsWeekday,
			&victimCaseForm1Pg.VictimCaseForm1FactsDate, &victimCaseForm1Pg.VictimCaseForm1FactsDescription, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperienced,
			&victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperiencedOther, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScope, &victimCaseForm1Pg.VictimCaseForm1VictimFemicideRisk,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimRelationshipWithAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorName,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocType, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocNumber, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorAddress,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorPhone, &victimCaseForm1Pg.VictimCaseForm1Age, &victimCaseForm1Pg.VictimCaseForm1VictimNationality, &victimCaseForm1Pg.VictimCaseForm1VictimNationalityOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimForeignerImmigrationStatus, &victimCaseForm1Pg.VictimCaseForm1VictimGender, &victimCaseForm1Pg.VictimCaseForm1VictimGenderIdentityOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimSexualOrientationOther, &victimCaseForm1Pg.VictimCaseForm1VictimDependents, &victimCaseForm1Pg.VictimCaseForm1VictimDeathThreats,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorHasWeapons, &victimCaseForm1Pg.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore,
			&victimCaseForm1Pg.VictimCaseForm1VictimImminentRisk, &victimCaseForm1Pg.VictimCaseForm1VictimPreviouslyReportedSituation, &victimCaseForm1Pg.VictimCaseForm1VictimIfPreviouslyReported,
			&victimCaseForm1Pg.VictimCaseForm1VictimIfAfro, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenous, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenousTongue,
			&victimCaseForm1Pg.VictimCaseForm1VictimIfPeasant, &victimCaseForm1Pg.VictimCaseForm1VictimIfArmedConflict, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScene,
			&victimCaseForm1Pg.VictimCaseForm1PhysicalViolenceIncreased, &victimCaseForm1Pg.VictimCaseForm1SeparatedFromPartnerLastYear, &victimCaseForm1Pg.VictimCaseForm1ThreatenedWithWeapon,
			&victimCaseForm1Pg.VictimCaseForm1ThreatenedToKillOrHarmChildren, &victimCaseForm1Pg.VictimCaseForm1JealousAndViolent, &victimCaseForm1Pg.VictimCaseForm1BelievesCapableOfKilling, &victimCaseForm1Pg.VictimCaseForm1VictimCase)
		victimCaseForm1s = append(victimCaseForm1s, victimCaseForm1Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCaseForm1Path +
			` WHERE ` + victimCaseForm1Path + `.` + VictimCaseForm1FieldDefinitions["VictimCaseForm1Status"].DBName + ` = $1 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, victimCaseForm1Status)
		persistenceCtrl.Scan(&count)
	}

	return victimCaseForm1s, count, nil
}

func UpdateVictimCaseForm1(victimCaseForm1 *VictimCaseForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var victimCaseForm1FieldsSlice []string = []string{"VictimCaseForm1UpdateDate",
		"VictimCaseForm1Nick", "VictimCaseForm1BirthDate", "VictimCaseForm1ViolenceTownCode", "VictimCaseForm1Address",
		"VictimCaseForm1LivingLatitude", "VictimCaseForm1LivingLongitude", "VictimCaseForm1Phone", "VictimCaseForm1Email", "VictimCaseForm1GenderIdentity", "VictimCaseForm1SexualOrientation", "VictimCaseForm1Origin",
		"VictimCaseForm1Occupation", "VictimCaseForm1OccupationOther", "VictimCaseForm1VictimEthnicGroup", "VictimCaseForm1VictimEthnicGroupOther", "VictimCaseForm1VictimContactNames",
		"VictimCaseForm1VictimContactPhone", "VictimCaseForm1VictimContactKinship", "VictimCaseForm1VictimNumChildren", "VictimCaseForm1VictimMaritalStatus",
		"VictimCaseForm1VictimMaritalStatusOther", "VictimCaseForm1VictimChildrenAge", "VictimCaseForm1VictimDisability", "VictimCaseForm1VictimSpecialSupport", "VictimCaseForm1FactsOccurrence",
		"VictimCaseForm1FactsStartTime", "VictimCaseForm1FactsEndTime", "VictimCaseForm1FactsWeekday", "VictimCaseForm1FactsDate", "VictimCaseForm1FactsDescription",
		"VictimCaseForm1VictimViolenceExperienced", "VictimCaseForm1VictimViolenceExperiencedOther", "VictimCaseForm1VictimViolenceScope", "VictimCaseForm1VictimFemicideRisk",
		"VictimCaseForm1VictimAggressor", "VictimCaseForm1VictimRelationshipWithAggressor", "VictimCaseForm1VictimAggressorName", "VictimCaseForm1VictimAggressorDocType",
		"VictimCaseForm1VictimAggressorDocNumber", "VictimCaseForm1VictimAggressorAddress", "VictimCaseForm1VictimAggressorPhone",
		"VictimCaseForm1Age", "VictimCaseForm1VictimNationality", "VictimCaseForm1VictimNationalityOther", "VictimCaseForm1VictimForeignerImmigrationStatus",
		"VictimCaseForm1VictimGender", "VictimCaseForm1VictimGenderIdentityOther", "VictimCaseForm1VictimSexualOrientationOther", "VictimCaseForm1VictimDependents",
		"VictimCaseForm1VictimDeathThreats", "VictimCaseForm1VictimAggressorHasWeapons", "VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore",
		"VictimCaseForm1VictimImminentRisk", "VictimCaseForm1VictimPreviouslyReportedSituation", "VictimCaseForm1VictimIfPreviouslyReported",
		"VictimCaseForm1VictimIfAfro", "VictimCaseForm1VictimIfIndigenous", "VictimCaseForm1VictimIfIndigenousTongue", "VictimCaseForm1VictimIfPeasant",
		"VictimCaseForm1VictimIfArmedConflict", "VictimCaseForm1VictimViolenceScene",
		"VictimCaseForm1PhysicalViolenceIncreased", "VictimCaseForm1SeparatedFromPartnerLastYear", "VictimCaseForm1ThreatenedWithWeapon", "VictimCaseForm1ThreatenedToKillOrHarmChildren",
		"VictimCaseForm1JealousAndViolent", "VictimCaseForm1BelievesCapableOfKilling", "VictimCaseForm1VictimCase"}
	var victimCaseForm1FieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, victimCaseForm1FieldsSlice, victimCaseForm1FieldsAliasSlice, VictimCaseForm1DBName, []string{"VictimCaseForm1Id"}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, victimCaseForm1.VictimCaseForm1Id,
		victimCaseForm1.VictimCaseForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCaseForm1.VictimCaseForm1Nick, victimCaseForm1.VictimCaseForm1BirthDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCaseForm1.VictimCaseForm1ViolenceTownCode, victimCaseForm1.VictimCaseForm1Address,
		victimCaseForm1.VictimCaseForm1LivingLatitude, victimCaseForm1.VictimCaseForm1LivingLongitude, victimCaseForm1.VictimCaseForm1Phone, victimCaseForm1.VictimCaseForm1Email, victimCaseForm1.VictimCaseForm1GenderIdentity, victimCaseForm1.VictimCaseForm1SexualOrientation,
		victimCaseForm1.VictimCaseForm1Origin, victimCaseForm1.VictimCaseForm1Occupation, victimCaseForm1.VictimCaseForm1OccupationOther, victimCaseForm1.VictimCaseForm1VictimEthnicGroup, victimCaseForm1.VictimCaseForm1VictimEthnicGroupOther,
		victimCaseForm1.VictimCaseForm1VictimContactNames, victimCaseForm1.VictimCaseForm1VictimContactPhone, victimCaseForm1.VictimCaseForm1VictimContactKinship, victimCaseForm1.VictimCaseForm1VictimNumChildren,
		victimCaseForm1.VictimCaseForm1VictimMaritalStatus, victimCaseForm1.VictimCaseForm1VictimMaritalStatusOther, victimCaseForm1.VictimCaseForm1VictimChildrenAge, victimCaseForm1.VictimCaseForm1VictimDisability,
		victimCaseForm1.VictimCaseForm1VictimSpecialSupport, victimCaseForm1.VictimCaseForm1FactsOccurrence, victimCaseForm1.VictimCaseForm1FactsStartTime.Format(common_config.DateTime.TIME_FORMAT),
		victimCaseForm1.VictimCaseForm1FactsEndTime.Format(common_config.DateTime.TIME_FORMAT), victimCaseForm1.VictimCaseForm1FactsWeekday,
		victimCaseForm1.VictimCaseForm1FactsDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), victimCaseForm1.VictimCaseForm1FactsDescription, victimCaseForm1.VictimCaseForm1VictimViolenceExperienced,
		victimCaseForm1.VictimCaseForm1VictimViolenceExperiencedOther, victimCaseForm1.VictimCaseForm1VictimViolenceScope, victimCaseForm1.VictimCaseForm1VictimFemicideRisk,
		victimCaseForm1.VictimCaseForm1VictimAggressor, victimCaseForm1.VictimCaseForm1VictimRelationshipWithAggressor, victimCaseForm1.VictimCaseForm1VictimAggressorName,
		victimCaseForm1.VictimCaseForm1VictimAggressorDocType, victimCaseForm1.VictimCaseForm1VictimAggressorDocNumber, victimCaseForm1.VictimCaseForm1VictimAggressorAddress,
		victimCaseForm1.VictimCaseForm1VictimAggressorPhone, victimCaseForm1.VictimCaseForm1Age, victimCaseForm1.VictimCaseForm1VictimNationality, victimCaseForm1.VictimCaseForm1VictimNationalityOther,
		victimCaseForm1.VictimCaseForm1VictimForeignerImmigrationStatus, victimCaseForm1.VictimCaseForm1VictimGender, victimCaseForm1.VictimCaseForm1VictimGenderIdentityOther,
		victimCaseForm1.VictimCaseForm1VictimSexualOrientationOther, victimCaseForm1.VictimCaseForm1VictimDependents, victimCaseForm1.VictimCaseForm1VictimDeathThreats,
		victimCaseForm1.VictimCaseForm1VictimAggressorHasWeapons, victimCaseForm1.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore,
		victimCaseForm1.VictimCaseForm1VictimImminentRisk, victimCaseForm1.VictimCaseForm1VictimPreviouslyReportedSituation, victimCaseForm1.VictimCaseForm1VictimIfPreviouslyReported,
		victimCaseForm1.VictimCaseForm1VictimIfAfro, victimCaseForm1.VictimCaseForm1VictimIfIndigenous, victimCaseForm1.VictimCaseForm1VictimIfIndigenousTongue,
		victimCaseForm1.VictimCaseForm1VictimIfPeasant, victimCaseForm1.VictimCaseForm1VictimIfArmedConflict, victimCaseForm1.VictimCaseForm1VictimViolenceScene,
		victimCaseForm1.VictimCaseForm1PhysicalViolenceIncreased, victimCaseForm1.VictimCaseForm1SeparatedFromPartnerLastYear, victimCaseForm1.VictimCaseForm1ThreatenedWithWeapon,
		victimCaseForm1.VictimCaseForm1ThreatenedToKillOrHarmChildren, victimCaseForm1.VictimCaseForm1JealousAndViolent, victimCaseForm1.VictimCaseForm1BelievesCapableOfKilling, victimCaseForm1.VictimCaseForm1VictimCase.(VictimCaseDTO).VictimCaseId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// Algunas utilidades
func SetVictimCaseForm1Defaults(victimCaseForm1 *VictimCaseForm1DTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		victimCaseForm1.VictimCaseForm1CreationDate = time.Now()
		victimCaseForm1.VictimCaseForm1UpdateDate = time.Now()
		victimCaseForm1.VictimCaseForm1ICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		victimCaseForm1.VictimCaseForm1UpdateDate = time.Now()
	}

}

func (obj *VictimCaseForm1PgDB) ToDTO() VictimCaseForm1DTO {
	var dto VictimCaseForm1DTO

	if obj.VictimCaseForm1Id.Valid {

		dto.VictimCaseForm1Id = uint64(obj.VictimCaseForm1Id.Int64)
	}

	if obj.VictimCaseForm1ICode.Valid {

		dto.VictimCaseForm1ICode = obj.VictimCaseForm1ICode.String
	}

	if obj.VictimCaseForm1CreationDate.Valid {

		dto.VictimCaseForm1CreationDate = obj.VictimCaseForm1CreationDate.Time
	}

	if obj.VictimCaseForm1UpdateDate.Valid {

		dto.VictimCaseForm1UpdateDate = obj.VictimCaseForm1UpdateDate.Time
	}

	if obj.VictimCaseForm1Nick.Valid {
		dto.VictimCaseForm1Nick = obj.VictimCaseForm1Nick.String
	}

	if obj.VictimCaseForm1BirthDate.Valid {
		dto.VictimCaseForm1BirthDate = obj.VictimCaseForm1BirthDate.Time
	}

	if obj.VictimCaseForm1ViolenceTownCode.Valid {
		dto.VictimCaseForm1ViolenceTownCode = obj.VictimCaseForm1ViolenceTownCode.String
	}

	if obj.VictimCaseForm1Address.Valid {
		dto.VictimCaseForm1Address = obj.VictimCaseForm1Address.String
	}

	if obj.VictimCaseForm1LivingLatitude.Valid {
		dto.VictimCaseForm1LivingLatitude = obj.VictimCaseForm1LivingLatitude.Float64
	}

	if obj.VictimCaseForm1LivingLongitude.Valid {
		dto.VictimCaseForm1LivingLongitude = obj.VictimCaseForm1LivingLongitude.Float64
	}

	if obj.VictimCaseForm1Phone.Valid {
		dto.VictimCaseForm1Phone = obj.VictimCaseForm1Phone.String
	}

	if obj.VictimCaseForm1Email.Valid {
		dto.VictimCaseForm1Email = obj.VictimCaseForm1Email.String
	}

	if obj.VictimCaseForm1GenderIdentity.Valid {
		dto.VictimCaseForm1GenderIdentity = obj.VictimCaseForm1GenderIdentity.String
	}

	if obj.VictimCaseForm1SexualOrientation.Valid {
		dto.VictimCaseForm1SexualOrientation = obj.VictimCaseForm1SexualOrientation.String
	}

	if obj.VictimCaseForm1Origin.Valid {
		dto.VictimCaseForm1Origin = obj.VictimCaseForm1Origin.String
	}

	if obj.VictimCaseForm1Occupation.Valid {
		dto.VictimCaseForm1Occupation = obj.VictimCaseForm1Occupation.String
	}

	if obj.VictimCaseForm1OccupationOther.Valid {
		dto.VictimCaseForm1OccupationOther = obj.VictimCaseForm1OccupationOther.String
	}

	if obj.VictimCaseForm1VictimEthnicGroup.Valid {
		dto.VictimCaseForm1VictimEthnicGroup = obj.VictimCaseForm1VictimEthnicGroup.String
	}

	if obj.VictimCaseForm1VictimEthnicGroupOther.Valid {
		dto.VictimCaseForm1VictimEthnicGroupOther = obj.VictimCaseForm1VictimEthnicGroupOther.String
	}

	if obj.VictimCaseForm1VictimContactNames.Valid {
		dto.VictimCaseForm1VictimContactNames = obj.VictimCaseForm1VictimContactNames.String
	}

	if obj.VictimCaseForm1VictimContactPhone.Valid {
		dto.VictimCaseForm1VictimContactPhone = obj.VictimCaseForm1VictimContactPhone.String
	}

	if obj.VictimCaseForm1VictimContactKinship.Valid {
		dto.VictimCaseForm1VictimContactKinship = obj.VictimCaseForm1VictimContactKinship.String
	}

	if obj.VictimCaseForm1VictimNumChildren.Valid {
		dto.VictimCaseForm1VictimNumChildren = obj.VictimCaseForm1VictimNumChildren.Int16
	}

	if obj.VictimCaseForm1VictimMaritalStatus.Valid {
		dto.VictimCaseForm1VictimMaritalStatus = obj.VictimCaseForm1VictimMaritalStatus.String
	}

	if obj.VictimCaseForm1VictimMaritalStatusOther.Valid {
		dto.VictimCaseForm1VictimMaritalStatusOther = obj.VictimCaseForm1VictimMaritalStatusOther.String
	}

	if obj.VictimCaseForm1VictimChildrenAge.Valid {
		dto.VictimCaseForm1VictimChildrenAge = obj.VictimCaseForm1VictimChildrenAge.String
	}

	if obj.VictimCaseForm1VictimDisability.Valid {
		dto.VictimCaseForm1VictimDisability = obj.VictimCaseForm1VictimDisability.String
	}

	if obj.VictimCaseForm1VictimSpecialSupport.Valid {
		dto.VictimCaseForm1VictimSpecialSupport = obj.VictimCaseForm1VictimSpecialSupport.String
	}

	if obj.VictimCaseForm1FactsOccurrence.Valid {
		dto.VictimCaseForm1FactsOccurrence = obj.VictimCaseForm1FactsOccurrence.String
	}

	if obj.VictimCaseForm1FactsStartTime.Valid {
		dto.VictimCaseForm1FactsStartTime, _ = time.Parse(common_config.DateTime.TIME_WITH_MILLISECONDS_FORMAT, obj.VictimCaseForm1FactsStartTime.String)
	}

	if obj.VictimCaseForm1FactsEndTime.Valid {
		dto.VictimCaseForm1FactsEndTime, _ = time.Parse(common_config.DateTime.TIME_WITH_MILLISECONDS_FORMAT, obj.VictimCaseForm1FactsEndTime.String)
	}

	if obj.VictimCaseForm1FactsWeekday.Valid {
		dto.VictimCaseForm1FactsWeekday = uint16(obj.VictimCaseForm1FactsWeekday.Int16)
	}

	if obj.VictimCaseForm1FactsDate.Valid {
		dto.VictimCaseForm1FactsDate = obj.VictimCaseForm1FactsDate.Time
	}

	if obj.VictimCaseForm1FactsDescription.Valid {
		dto.VictimCaseForm1FactsDescription = obj.VictimCaseForm1FactsDescription.String
	}

	if obj.VictimCaseForm1VictimViolenceExperienced.Valid {
		dto.VictimCaseForm1VictimViolenceExperienced = obj.VictimCaseForm1VictimViolenceExperienced.String
	}

	if obj.VictimCaseForm1VictimViolenceExperiencedOther.Valid {
		dto.VictimCaseForm1VictimViolenceExperiencedOther = obj.VictimCaseForm1VictimViolenceExperiencedOther.String
	}

	if obj.VictimCaseForm1VictimViolenceScope.Valid {
		dto.VictimCaseForm1VictimViolenceScope = obj.VictimCaseForm1VictimViolenceScope.String
	}

	if obj.VictimCaseForm1VictimFemicideRisk.Valid {
		dto.VictimCaseForm1VictimFemicideRisk = obj.VictimCaseForm1VictimFemicideRisk.String
	}

	if obj.VictimCaseForm1VictimAggressor.Valid {
		dto.VictimCaseForm1VictimAggressor = obj.VictimCaseForm1VictimAggressor.String
	}

	if obj.VictimCaseForm1VictimRelationshipWithAggressor.Valid {
		dto.VictimCaseForm1VictimRelationshipWithAggressor = obj.VictimCaseForm1VictimRelationshipWithAggressor.String
	}

	if obj.VictimCaseForm1VictimAggressorName.Valid {
		dto.VictimCaseForm1VictimAggressorName = obj.VictimCaseForm1VictimAggressorName.String
	}

	if obj.VictimCaseForm1VictimAggressorDocType.Valid {
		dto.VictimCaseForm1VictimAggressorDocType = obj.VictimCaseForm1VictimAggressorDocType.String
	}

	if obj.VictimCaseForm1VictimAggressorDocNumber.Valid {
		dto.VictimCaseForm1VictimAggressorDocNumber = obj.VictimCaseForm1VictimAggressorDocNumber.String
	}

	if obj.VictimCaseForm1VictimAggressorAddress.Valid {
		dto.VictimCaseForm1VictimAggressorAddress = obj.VictimCaseForm1VictimAggressorAddress.String
	}

	if obj.VictimCaseForm1VictimAggressorPhone.Valid {
		dto.VictimCaseForm1VictimAggressorPhone = obj.VictimCaseForm1VictimAggressorPhone.String
	}

	if obj.VictimCaseForm1Age.Valid {
		dto.VictimCaseForm1Age = uint16(obj.VictimCaseForm1Age.Int16)
	}

	if obj.VictimCaseForm1VictimNationality.Valid {
		dto.VictimCaseForm1VictimNationality = obj.VictimCaseForm1VictimNationality.String
	}

	if obj.VictimCaseForm1VictimNationalityOther.Valid {
		dto.VictimCaseForm1VictimNationalityOther = obj.VictimCaseForm1VictimNationalityOther.String
	}

	if obj.VictimCaseForm1VictimForeignerImmigrationStatus.Valid {
		dto.VictimCaseForm1VictimForeignerImmigrationStatus = obj.VictimCaseForm1VictimForeignerImmigrationStatus.String
	}

	if obj.VictimCaseForm1VictimGender.Valid {
		dto.VictimCaseForm1VictimGender = obj.VictimCaseForm1VictimGender.String
	}

	if obj.VictimCaseForm1VictimGenderIdentityOther.Valid {
		dto.VictimCaseForm1VictimGenderIdentityOther = obj.VictimCaseForm1VictimGenderIdentityOther.String
	}

	if obj.VictimCaseForm1VictimSexualOrientationOther.Valid {
		dto.VictimCaseForm1VictimSexualOrientationOther = obj.VictimCaseForm1VictimSexualOrientationOther.String
	}

	if obj.VictimCaseForm1VictimDependents.Valid {
		dto.VictimCaseForm1VictimDependents = obj.VictimCaseForm1VictimDependents.String
	}

	if obj.VictimCaseForm1VictimDeathThreats.Valid {
		dto.VictimCaseForm1VictimDeathThreats = obj.VictimCaseForm1VictimDeathThreats.String
	}

	if obj.VictimCaseForm1VictimAggressorHasWeapons.Valid {
		dto.VictimCaseForm1VictimAggressorHasWeapons = obj.VictimCaseForm1VictimAggressorHasWeapons.String
	}

	if obj.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore.Valid {
		dto.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore = obj.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore.String
	}

	if obj.VictimCaseForm1VictimImminentRisk.Valid {
		dto.VictimCaseForm1VictimImminentRisk = obj.VictimCaseForm1VictimImminentRisk.String
	}

	if obj.VictimCaseForm1VictimIfAfro.Valid {
		dto.VictimCaseForm1VictimIfAfro = obj.VictimCaseForm1VictimIfAfro.String
	}
	if obj.VictimCaseForm1VictimIfIndigenous.Valid {
		dto.VictimCaseForm1VictimIfIndigenous = obj.VictimCaseForm1VictimIfIndigenous.String
	}
	if obj.VictimCaseForm1VictimIfIndigenousTongue.Valid {
		dto.VictimCaseForm1VictimIfIndigenousTongue = obj.VictimCaseForm1VictimIfIndigenousTongue.String
	}
	if obj.VictimCaseForm1VictimIfPeasant.Valid {
		dto.VictimCaseForm1VictimIfPeasant = obj.VictimCaseForm1VictimIfPeasant.String
	}
	if obj.VictimCaseForm1VictimIfArmedConflict.Valid {
		dto.VictimCaseForm1VictimIfArmedConflict = obj.VictimCaseForm1VictimIfArmedConflict.String
	}
	if obj.VictimCaseForm1VictimViolenceScene.Valid {
		dto.VictimCaseForm1VictimViolenceScene = obj.VictimCaseForm1VictimViolenceScene.String
	}

	if obj.VictimCaseForm1VictimPreviouslyReportedSituation.Valid {
		dto.VictimCaseForm1VictimPreviouslyReportedSituation = obj.VictimCaseForm1VictimPreviouslyReportedSituation.String
	}

	if obj.VictimCaseForm1VictimIfPreviouslyReported.Valid {
		dto.VictimCaseForm1VictimIfPreviouslyReported = obj.VictimCaseForm1VictimIfPreviouslyReported.String
	}

	if obj.VictimCaseForm1PhysicalViolenceIncreased.Valid {
		dto.VictimCaseForm1PhysicalViolenceIncreased = obj.VictimCaseForm1PhysicalViolenceIncreased.String
	}

	if obj.VictimCaseForm1SeparatedFromPartnerLastYear.Valid {
		dto.VictimCaseForm1SeparatedFromPartnerLastYear = obj.VictimCaseForm1SeparatedFromPartnerLastYear.String
	}

	if obj.VictimCaseForm1ThreatenedWithWeapon.Valid {
		dto.VictimCaseForm1ThreatenedWithWeapon = obj.VictimCaseForm1ThreatenedWithWeapon.String
	}

	if obj.VictimCaseForm1ThreatenedToKillOrHarmChildren.Valid {
		dto.VictimCaseForm1ThreatenedToKillOrHarmChildren = obj.VictimCaseForm1ThreatenedToKillOrHarmChildren.String
	}

	if obj.VictimCaseForm1JealousAndViolent.Valid {
		dto.VictimCaseForm1JealousAndViolent = obj.VictimCaseForm1JealousAndViolent.String
	}

	if obj.VictimCaseForm1BelievesCapableOfKilling.Valid {
		dto.VictimCaseForm1BelievesCapableOfKilling = obj.VictimCaseForm1BelievesCapableOfKilling.String
	}

	if obj.VictimCaseForm1VictimCase.Valid {
		dto.VictimCaseForm1VictimCase = VictimCaseDTO{VictimCaseId: uint64(obj.VictimCaseForm1VictimCase.Int64)}
	}

	return dto
}

func (obj *VictimCaseForm1PgDB) ToDTOTranslated() VictimCaseForm1DTO {
	var dto VictimCaseForm1DTO

	if obj.VictimCaseForm1Id.Valid {

		dto.VictimCaseForm1Id = uint64(obj.VictimCaseForm1Id.Int64)
	}

	if obj.VictimCaseForm1ICode.Valid {

		dto.VictimCaseForm1ICode = obj.VictimCaseForm1ICode.String
	}

	if obj.VictimCaseForm1CreationDate.Valid {

		dto.VictimCaseForm1CreationDate = obj.VictimCaseForm1CreationDate.Time
	}

	if obj.VictimCaseForm1UpdateDate.Valid {

		dto.VictimCaseForm1UpdateDate = obj.VictimCaseForm1UpdateDate.Time
	}

	if obj.VictimCaseForm1Nick.Valid {
		dto.VictimCaseForm1Nick = obj.VictimCaseForm1Nick.String
	}

	if obj.VictimCaseForm1BirthDate.Valid {
		dto.VictimCaseForm1BirthDate = obj.VictimCaseForm1BirthDate.Time
	}

	if obj.VictimCaseForm1ViolenceTownCode.Valid {
		dto.VictimCaseForm1ViolenceTownCode = obj.VictimCaseForm1ViolenceTownCode.String
	}

	if obj.VictimCaseForm1Address.Valid {
		dto.VictimCaseForm1Address = obj.VictimCaseForm1Address.String
	}

	if obj.VictimCaseForm1LivingLatitude.Valid {
		dto.VictimCaseForm1LivingLatitude = obj.VictimCaseForm1LivingLatitude.Float64
	}

	if obj.VictimCaseForm1LivingLongitude.Valid {
		dto.VictimCaseForm1LivingLongitude = obj.VictimCaseForm1LivingLongitude.Float64
	}

	if obj.VictimCaseForm1Phone.Valid {
		dto.VictimCaseForm1Phone = obj.VictimCaseForm1Phone.String
	}

	if obj.VictimCaseForm1Email.Valid {
		dto.VictimCaseForm1Email = obj.VictimCaseForm1Email.String
	}

	if obj.VictimCaseForm1GenderIdentity.Valid {
		dto.VictimCaseForm1GenderIdentity = common_config.GENDER_IDENTITY[obj.VictimCaseForm1GenderIdentity.String]
	}

	if obj.VictimCaseForm1SexualOrientation.Valid {
		dto.VictimCaseForm1SexualOrientation = salvia_config.SEXUAL_ORIENTATION[obj.VictimCaseForm1SexualOrientation.String]
	}

	if obj.VictimCaseForm1Origin.Valid {
		dto.VictimCaseForm1Origin = salvia_config.ORIGIN_PLACE[obj.VictimCaseForm1Origin.String]
	}

	if obj.VictimCaseForm1Occupation.Valid {
		dto.VictimCaseForm1Occupation = salvia_config.OCCUPATION[obj.VictimCaseForm1Occupation.String]
	}

	if obj.VictimCaseForm1OccupationOther.Valid {
		dto.VictimCaseForm1OccupationOther = obj.VictimCaseForm1OccupationOther.String
	}

	if obj.VictimCaseForm1VictimEthnicGroup.Valid {
		dto.VictimCaseForm1VictimEthnicGroup = salvia_config.ETHNIC_GROUP[obj.VictimCaseForm1VictimEthnicGroup.String]
	}

	if obj.VictimCaseForm1VictimEthnicGroupOther.Valid {
		dto.VictimCaseForm1VictimEthnicGroupOther = obj.VictimCaseForm1VictimEthnicGroupOther.String
	}

	if obj.VictimCaseForm1VictimContactNames.Valid {
		dto.VictimCaseForm1VictimContactNames = obj.VictimCaseForm1VictimContactNames.String
	}

	if obj.VictimCaseForm1VictimContactPhone.Valid {
		dto.VictimCaseForm1VictimContactPhone = obj.VictimCaseForm1VictimContactPhone.String
	}

	if obj.VictimCaseForm1VictimContactKinship.Valid {
		dto.VictimCaseForm1VictimContactKinship = obj.VictimCaseForm1VictimContactKinship.String
	}

	if obj.VictimCaseForm1VictimNumChildren.Valid {
		dto.VictimCaseForm1VictimNumChildren = obj.VictimCaseForm1VictimNumChildren.Int16
	}

	if obj.VictimCaseForm1VictimMaritalStatus.Valid {
		dto.VictimCaseForm1VictimMaritalStatus = common_config.MARITAL_STATUS[obj.VictimCaseForm1VictimMaritalStatus.String]
	}

	if obj.VictimCaseForm1VictimMaritalStatusOther.Valid {
		dto.VictimCaseForm1VictimMaritalStatusOther = obj.VictimCaseForm1VictimMaritalStatusOther.String
	}

	if obj.VictimCaseForm1VictimChildrenAge.Valid {
		dto.VictimCaseForm1VictimChildrenAge = obj.VictimCaseForm1VictimChildrenAge.String
	}

	if obj.VictimCaseForm1VictimDisability.Valid {
		dto.VictimCaseForm1VictimDisability = salvia_config.DISABILITY[obj.VictimCaseForm1VictimDisability.String]
	}

	if obj.VictimCaseForm1VictimSpecialSupport.Valid {
		dto.VictimCaseForm1VictimSpecialSupport = obj.VictimCaseForm1VictimSpecialSupport.String
	}

	if obj.VictimCaseForm1FactsOccurrence.Valid {
		dto.VictimCaseForm1FactsOccurrence = salvia_config.OCCURRENCE[obj.VictimCaseForm1FactsOccurrence.String]
	}

	if obj.VictimCaseForm1FactsStartTime.Valid {
		dto.VictimCaseForm1FactsStartTime, _ = time.Parse(common_config.DateTime.TIME_WITH_MILLISECONDS_FORMAT, obj.VictimCaseForm1FactsStartTime.String)
	}

	if obj.VictimCaseForm1FactsEndTime.Valid {
		dto.VictimCaseForm1FactsEndTime, _ = time.Parse(common_config.DateTime.TIME_WITH_MILLISECONDS_FORMAT, obj.VictimCaseForm1FactsEndTime.String)
	}

	if obj.VictimCaseForm1FactsWeekday.Valid {
		dto.VictimCaseForm1FactsWeekday = uint16(obj.VictimCaseForm1FactsWeekday.Int16)
	}

	if obj.VictimCaseForm1FactsDate.Valid {
		dto.VictimCaseForm1FactsDate = obj.VictimCaseForm1FactsDate.Time
	}

	if obj.VictimCaseForm1FactsDescription.Valid {
		dto.VictimCaseForm1FactsDescription = obj.VictimCaseForm1FactsDescription.String
	}

	if obj.VictimCaseForm1VictimViolenceExperienced.Valid {
		dto.VictimCaseForm1VictimViolenceExperienced = salvia_config.VIOLENCE_EXPERIENCED[obj.VictimCaseForm1VictimViolenceExperienced.String]
	}

	if obj.VictimCaseForm1VictimViolenceExperiencedOther.Valid {
		dto.VictimCaseForm1VictimViolenceExperiencedOther = obj.VictimCaseForm1VictimViolenceExperiencedOther.String
	}

	if obj.VictimCaseForm1VictimViolenceScope.Valid {
		dto.VictimCaseForm1VictimViolenceScope = salvia_config.VIOLENCE_SCOPE[obj.VictimCaseForm1VictimViolenceScope.String]
	}

	if obj.VictimCaseForm1VictimFemicideRisk.Valid {
		dto.VictimCaseForm1VictimFemicideRisk = common_config.YES_NO[obj.VictimCaseForm1VictimFemicideRisk.String]
	}

	if obj.VictimCaseForm1VictimAggressor.Valid {
		dto.VictimCaseForm1VictimAggressor = salvia_config.AGGRESSOR[obj.VictimCaseForm1VictimAggressor.String]
	}

	if obj.VictimCaseForm1VictimRelationshipWithAggressor.Valid {
		dto.VictimCaseForm1VictimRelationshipWithAggressor = salvia_config.RELATIONSHIP_WITH_AGGRESSOR[obj.VictimCaseForm1VictimRelationshipWithAggressor.String]
	}

	if obj.VictimCaseForm1VictimAggressorName.Valid {
		dto.VictimCaseForm1VictimAggressorName = obj.VictimCaseForm1VictimAggressorName.String
	}

	if obj.VictimCaseForm1VictimAggressorDocType.Valid {
		dto.VictimCaseForm1VictimAggressorDocType = common_config.DOCUMENT_TYPE[obj.VictimCaseForm1VictimAggressorDocType.String]
	}

	if obj.VictimCaseForm1VictimAggressorDocNumber.Valid {
		dto.VictimCaseForm1VictimAggressorDocNumber = obj.VictimCaseForm1VictimAggressorDocNumber.String
	}

	if obj.VictimCaseForm1VictimAggressorAddress.Valid {
		dto.VictimCaseForm1VictimAggressorAddress = obj.VictimCaseForm1VictimAggressorAddress.String
	}

	if obj.VictimCaseForm1VictimAggressorPhone.Valid {
		dto.VictimCaseForm1VictimAggressorPhone = obj.VictimCaseForm1VictimAggressorPhone.String
	}

	if obj.VictimCaseForm1Age.Valid {
		dto.VictimCaseForm1Age = uint16(obj.VictimCaseForm1Age.Int16)
	}

	if obj.VictimCaseForm1VictimNationality.Valid {
		dto.VictimCaseForm1VictimNationality = salvia_config.VICTIM_CASE_VICTIM_NATIONALITY[obj.VictimCaseForm1VictimNationality.String]
	}

	if obj.VictimCaseForm1VictimNationalityOther.Valid {
		dto.VictimCaseForm1VictimNationalityOther = obj.VictimCaseForm1VictimNationalityOther.String
	}

	if obj.VictimCaseForm1VictimForeignerImmigrationStatus.Valid {
		dto.VictimCaseForm1VictimForeignerImmigrationStatus = salvia_config.VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS[obj.VictimCaseForm1VictimForeignerImmigrationStatus.String]
	}

	if obj.VictimCaseForm1VictimGender.Valid {
		dto.VictimCaseForm1VictimGender = common_config.GENDER[obj.VictimCaseForm1VictimGender.String]
	}

	if obj.VictimCaseForm1VictimGenderIdentityOther.Valid {
		dto.VictimCaseForm1VictimGenderIdentityOther = obj.VictimCaseForm1VictimGenderIdentityOther.String
	}

	if obj.VictimCaseForm1VictimSexualOrientationOther.Valid {
		dto.VictimCaseForm1VictimSexualOrientationOther = obj.VictimCaseForm1VictimSexualOrientationOther.String
	}

	if obj.VictimCaseForm1VictimDependents.Valid {
		dto.VictimCaseForm1VictimDependents = salvia_config.VICTIM_CASE_VICTIM_DEPENDENTS[obj.VictimCaseForm1VictimDependents.String]
	}

	if obj.VictimCaseForm1VictimDeathThreats.Valid {
		dto.VictimCaseForm1VictimDeathThreats = common_config.YES_NO[obj.VictimCaseForm1VictimDeathThreats.String]
	}

	if obj.VictimCaseForm1VictimAggressorHasWeapons.Valid {
		dto.VictimCaseForm1VictimAggressorHasWeapons = common_config.YES_NO[obj.VictimCaseForm1VictimAggressorHasWeapons.String]
	}

	if obj.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore.Valid {
		dto.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore = common_config.YES_NO[obj.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore.String]
	}

	if obj.VictimCaseForm1VictimImminentRisk.Valid {
		dto.VictimCaseForm1VictimImminentRisk = common_config.YES_NO[obj.VictimCaseForm1VictimImminentRisk.String]
	}

	if obj.VictimCaseForm1VictimIfAfro.Valid {
		dto.VictimCaseForm1VictimIfAfro = salvia_config.AFRO_COMMUNTITIES["sp"][obj.VictimCaseForm1VictimIfAfro.String]
	}
	if obj.VictimCaseForm1VictimIfIndigenous.Valid {
		dto.VictimCaseForm1VictimIfIndigenous = salvia_config.COLOMBIAN_INDIGENOUS["sp"][obj.VictimCaseForm1VictimIfIndigenous.String]
	}
	if obj.VictimCaseForm1VictimIfIndigenousTongue.Valid {
		dto.VictimCaseForm1VictimIfIndigenousTongue = salvia_config.INDIGENOUS_TONGUES["sp"][obj.VictimCaseForm1VictimIfIndigenousTongue.String]
	}
	if obj.VictimCaseForm1VictimIfPeasant.Valid {
		dto.VictimCaseForm1VictimIfPeasant = common_config.YES_NO[obj.VictimCaseForm1VictimIfPeasant.String]
	}
	if obj.VictimCaseForm1VictimIfArmedConflict.Valid {
		dto.VictimCaseForm1VictimIfArmedConflict = common_config.YES_NO[obj.VictimCaseForm1VictimIfArmedConflict.String]
	}
	if obj.VictimCaseForm1VictimViolenceScene.Valid {
		dto.VictimCaseForm1VictimViolenceScene = salvia_config.VIOLENCE_SCENES["sp"][obj.VictimCaseForm1VictimViolenceScene.String]
	}

	if obj.VictimCaseForm1VictimPreviouslyReportedSituation.Valid {
		dto.VictimCaseForm1VictimPreviouslyReportedSituation = common_config.YES_NO[obj.VictimCaseForm1VictimPreviouslyReportedSituation.String]
	}

	if obj.VictimCaseForm1VictimIfPreviouslyReported.Valid {
		dto.VictimCaseForm1VictimIfPreviouslyReported = salvia_config.VICTIM_CASE_VICTIM_IF_PREVIOUSLY_REPORTED[obj.VictimCaseForm1VictimIfPreviouslyReported.String]
	}

	if obj.VictimCaseForm1PhysicalViolenceIncreased.Valid {
		dto.VictimCaseForm1PhysicalViolenceIncreased = common_config.YES_NO_NA[obj.VictimCaseForm1PhysicalViolenceIncreased.String]
	}

	if obj.VictimCaseForm1SeparatedFromPartnerLastYear.Valid {
		dto.VictimCaseForm1SeparatedFromPartnerLastYear = common_config.YES_NO_NA[obj.VictimCaseForm1SeparatedFromPartnerLastYear.String]
	}

	if obj.VictimCaseForm1ThreatenedWithWeapon.Valid {
		dto.VictimCaseForm1ThreatenedWithWeapon = common_config.YES_NO_NA[obj.VictimCaseForm1ThreatenedWithWeapon.String]
	}

	if obj.VictimCaseForm1ThreatenedToKillOrHarmChildren.Valid {
		dto.VictimCaseForm1ThreatenedToKillOrHarmChildren = common_config.YES_NO_NA[obj.VictimCaseForm1ThreatenedToKillOrHarmChildren.String]
	}

	if obj.VictimCaseForm1JealousAndViolent.Valid {
		dto.VictimCaseForm1JealousAndViolent = common_config.YES_NO_NA[obj.VictimCaseForm1JealousAndViolent.String]
	}

	if obj.VictimCaseForm1BelievesCapableOfKilling.Valid {
		dto.VictimCaseForm1BelievesCapableOfKilling = common_config.YES_NO_NA[obj.VictimCaseForm1BelievesCapableOfKilling.String]
	}

	if obj.VictimCaseForm1VictimCase.Valid {
		dto.VictimCaseForm1VictimCase = VictimCaseDTO{VictimCaseId: uint64(obj.VictimCaseForm1VictimCase.Int64)}
	}

	return dto
}
