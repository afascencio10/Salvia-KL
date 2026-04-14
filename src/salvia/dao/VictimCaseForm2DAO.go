package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"

	"encoding/json"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	VictimCaseForm2EntityName string = "VictimCaseForm2"
	VictimCaseForm2JSONName   string = "form2"
	VictimCaseForm2DBName     string = "victim_case_form2"
	VictimCaseForm2DBScheme   string = "salvia"

	VictimCaseForm2FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimCaseForm2Id":                                    {Name: "VictimCaseForm2Id", DBName: "victim_case_form2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2ICode":                                 {Name: "VictimCaseForm2ICode", DBName: "victim_case_form2_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimCaseForm2CreationDate":                          {Name: "VictimCaseForm2CreationDate", DBName: "victim_case_form2_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2UpdateDate":                            {Name: "VictimCaseForm2UpdateDate", DBName: "victim_case_form2_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2IdentityName":                          {Name: "VictimCaseForm2IdentityName", DBName: "victim_case_form2_identity_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm2VictimPhone":                           {Name: "VictimCaseForm2VictimPhone", DBName: "victim_case_form2_victim_phone", Alias: "", ModelType: "uint", MinSize: 1000000000, MaxSize: 9999999999, Required: false},
		"VictimCaseForm2FactsDescription":                      {Name: "VictimCaseForm2FactsDescription", DBName: "victim_case_form2_facts_description", Alias: "", ModelType: "string", MinSize: 5, MaxSize: 20000, Required: true},
		"VictimCaseForm2FactsDate":                             {Name: "VictimCaseForm2FactsDate", DBName: "victim_case_form2_facts_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2FactsStartTime":                        {Name: "VictimCaseForm2FactsStartTime", DBName: "victim_case_form2_facts_start_time", Alias: "", ModelType: "time", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2FactsTownCode":                         {Name: "VictimCaseForm2FactsTownCode", DBName: "victim_case_form2_facts_town_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 8, Required: true},
		"VictimCaseForm2FactsZone":                             {Name: "VictimCaseForm2FactsZone", DBName: "victim_case_form2_facts_zone", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2FactsAddress":                          {Name: "VictimCaseForm2FactsAddress", DBName: "victim_case_form2_facts_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 150, Required: true},
		"VictimCaseForm2ScenarioViolence":                      {Name: "VictimCaseForm2ScenarioViolence", DBName: "victim_case_form2_scenario_violence", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2ReportedPreviously":                    {Name: "VictimCaseForm2ReportedPreviously", DBName: "victim_case_form2_reported_previously", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2RecurrenceAggression":                  {Name: "VictimCaseForm2RecurrenceAggression", DBName: "victim_case_form2_recurrence_aggression", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2NumAgressors":                          {Name: "VictimCaseForm2NumAgressors", DBName: "victim_case_form2_num_agressors", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2ProximityPrincipalAggressor":           {Name: "VictimCaseForm2ProximityPrincipalAggressor", DBName: "victim_case_form2_proximity_principal_aggressor", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2RelationshipWithPresumedAggressor":     {Name: "VictimCaseForm2RelationshipWithPresumedAggressor", DBName: "victim_case_form2_relationship_with_presumed_aggressor", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2EconomicallyDependent":                 {Name: "VictimCaseForm2EconomicallyDependent", DBName: "victim_case_form2_economically_dependent", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorGenderIdentity":               {Name: "VictimCaseForm2AggressorGenderIdentity", DBName: "victim_case_form2_aggressor_gender_identity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorNames":                        {Name: "VictimCaseForm2AggressorNames", DBName: "victim_case_form2_aggressor_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: false},
		"VictimCaseForm2AggressorDocType":                      {Name: "VictimCaseForm2AggressorDocType", DBName: "victim_case_form2_aggressor_doc_type", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorDocNumber":                    {Name: "VictimCaseForm2AggressorDocNumber", DBName: "victim_case_form2_aggressor_doc_number", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm2AggressorAddress":                      {Name: "VictimCaseForm2AggressorAddress", DBName: "victim_case_form2_aggressor_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 150, Required: false},
		"VictimCaseForm2AggressorPhone":                        {Name: "VictimCaseForm2AggressorPhone", DBName: "victim_case_form2_aggressor_phone", Alias: "", ModelType: "uint", MinSize: 1000000000, MaxSize: 9999999999, Required: false},
		"VictimCaseForm2AggressorViolencePhysicalIncrease":     {Name: "VictimCaseForm2AggressorViolencePhysicalIncrease", DBName: "victim_case_form2_aggressor_violence_physical_increase", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AggressorWeaponUsed":                   {Name: "VictimCaseForm2AggressorWeaponUsed", DBName: "victim_case_form2_aggressor_weapon_used", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AggressorThreatKill":                   {Name: "VictimCaseForm2AggressorThreatKill", DBName: "victim_case_form2_aggressor_threat_kill", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AggressorPursuesSpiesDestroys":         {Name: "VictimCaseForm2AggressorPursuesSpiesDestroys", DBName: "victim_case_form2_aggressor_pursues_spies_destroys", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AggressorCapableOfKilling":             {Name: "VictimCaseForm2AggressorCapableOfKilling", DBName: "victim_case_form2_aggressor_capable_of_killing", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AggressorHasAccessToWeapons":           {Name: "VictimCaseForm2AggressorHasAccessToWeapons", DBName: "victim_case_form2_aggressor_has_access_to_weapons", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2PartnerUnemployed":                     {Name: "VictimCaseForm2PartnerUnemployed", DBName: "victim_case_form2_partner_unemployed", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2PartnerOtherDenunciations":             {Name: "VictimCaseForm2PartnerOtherDenunciations", DBName: "victim_case_form2_partner_other_denunciations", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorHasPenalBackground":           {Name: "VictimCaseForm2AggressorHasPenalBackground", DBName: "victim_case_form2_aggressor_has_penal_background", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorForcedSex":                    {Name: "VictimCaseForm2AggressorForcedSex", DBName: "victim_case_form2_aggressor_forced_sex", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorAttemptedStrangulation":       {Name: "VictimCaseForm2AggressorAttemptedStrangulation", DBName: "victim_case_form2_aggressor_attempted_strangulation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorConsumesDrugs":                {Name: "VictimCaseForm2AggressorConsumesDrugs", DBName: "victim_case_form2_aggressor_consumes_drugs", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorIsAlcoholic":                  {Name: "VictimCaseForm2AggressorIsAlcoholic", DBName: "victim_case_form2_aggressor_is_alcoholic", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2PartnerControls":                       {Name: "VictimCaseForm2PartnerControls", DBName: "victim_case_form2_partner_controls", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorHadHitInVulnerability":        {Name: "VictimCaseForm2AggressorHadHitInVulnerability", DBName: "victim_case_form2_aggressor_had_hit_in_vulnerability", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2PartnerThreatenedSuicide":              {Name: "VictimCaseForm2PartnerThreatenedSuicide", DBName: "victim_case_form2_partner_threatened_suicide", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2PartnerThreatenedDamageMembers":        {Name: "VictimCaseForm2PartnerThreatenedDamageMembers", DBName: "victim_case_form2_partner_threatened_damage_members", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2ThoughtsOfSelfHarm":                    {Name: "VictimCaseForm2ThoughtsOfSelfHarm", DBName: "victim_case_form2_thoughts_of_self_harm", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorLimitsContactSupportNetworks": {Name: "VictimCaseForm2AggressorLimitsContactSupportNetworks", DBName: "victim_case_form2_aggressor_limits_contact_support_networks", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2StillLivesWithAggressor":               {Name: "VictimCaseForm2StillLivesWithAggressor", DBName: "victim_case_form2_still_lives_with_aggressor", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorViolentlyJealous":             {Name: "VictimCaseForm2AggressorViolentlyJealous", DBName: "victim_case_form2_aggressor_violently_jealous", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorUnemployed":                   {Name: "VictimCaseForm2AggressorUnemployed", DBName: "victim_case_form2_aggressor_unemployed", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorHasPenalBackground2":          {Name: "VictimCaseForm2AggressorHasPenalBackground2", DBName: "victim_case_form2_aggressor_has_penal_background_2", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorSexuallyHarassment":           {Name: "VictimCaseForm2AggressorSexuallyHarassment", DBName: "victim_case_form2_aggressor_sexually_harassment", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorUseDrugs":                     {Name: "VictimCaseForm2AggressorUseDrugs", DBName: "victim_case_form2_aggressor_use_drugs", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorIsAlcoholic2":                 {Name: "VictimCaseForm2AggressorIsAlcoholic2", DBName: "victim_case_form2_aggressor_is_alcoholic_2", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorControls":                     {Name: "VictimCaseForm2AggressorControls", DBName: "victim_case_form2_aggressor_controls", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorThreatenedDamageMembers":      {Name: "VictimCaseForm2AggressorThreatenedDamageMembers", DBName: "victim_case_form2_aggressor_threatened_damage_members", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2ThoughtsOfSelfHarm2":                   {Name: "VictimCaseForm2ThoughtsOfSelfHarm2", DBName: "victim_case_form2_thoughts_of_self_harm_2", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorCommonSpaces":                 {Name: "VictimCaseForm2AggressorCommonSpaces", DBName: "victim_case_form2_aggressor_common_spaces", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorHierarchy":                    {Name: "VictimCaseForm2AggressorHierarchy", DBName: "victim_case_form2_aggressor_hierarchy", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2BirthDate":                             {Name: "VictimCaseForm2BirthDate", DBName: "victim_case_form2_birth_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2PhysicalMentalSensoryDifficulties":     {Name: "VictimCaseForm2PhysicalMentalSensoryDifficulties", DBName: "victim_case_form2_physical_mental_sensory_difficulties", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2Nationality":                           {Name: "VictimCaseForm2Nationality", DBName: "victim_case_form2_nationality", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2SpecifiedNationality":                  {Name: "VictimCaseForm2SpecifiedNationality", DBName: "victim_case_form2_specified_nationality", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2MigrationCondition":                    {Name: "VictimCaseForm2MigrationCondition", DBName: "victim_case_form2_migration_condition", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2GenderIdentity":                        {Name: "VictimCaseForm2GenderIdentity", DBName: "victim_case_form2_gender_identity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2SexualOrientation":                     {Name: "VictimCaseForm2SexualOrientation", DBName: "victim_case_form2_sexual_orientation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AssignedSexAtBirth":                    {Name: "VictimCaseForm2AssignedSexAtBirth", DBName: "victim_case_form2_assigned_sex_at_birth", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2EthnicAffiliation":                     {Name: "VictimCaseForm2EthnicAffiliation", DBName: "victim_case_form2_ethnic_affiliation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2IndigenousPeople":                      {Name: "VictimCaseForm2IndigenousPeople", DBName: "victim_case_form2_indigenous_people", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2CampesinoRecognition":                  {Name: "VictimCaseForm2CampesinoRecognition", DBName: "victim_case_form2_campesino_recognition", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2MaritalStatus":                         {Name: "VictimCaseForm2MaritalStatus", DBName: "victim_case_form2_marital_status", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2LastEducationLevel":                    {Name: "VictimCaseForm2LastEducationLevel", DBName: "victim_case_form2_last_education_level", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2Occupation":                            {Name: "VictimCaseForm2Occupation", DBName: "victim_case_form2_occupation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2IncomeGenerationMethod":                {Name: "VictimCaseForm2IncomeGenerationMethod", DBName: "victim_case_form2_income_generation_method", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2EmploymentRelationship":                {Name: "VictimCaseForm2EmploymentRelationship", DBName: "victim_case_form2_employment_relationship", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2ApproxStartAsp":                        {Name: "VictimCaseForm2ApproxStartAsp", DBName: "victim_case_form2_approx_start_asp", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2HousingTenancyForm":                    {Name: "VictimCaseForm2HousingTenancyForm", DBName: "victim_case_form2_housing_tenancy_form", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2HousingStratum":                        {Name: "VictimCaseForm2HousingStratum", DBName: "victim_case_form2_housing_stratum", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2CurrentlyPregnant":                     {Name: "VictimCaseForm2CurrentlyPregnant", DBName: "victim_case_form2_currently_pregnant", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2ResidenceTownCode":                     {Name: "VictimCaseForm2ResidenceTownCode", DBName: "victim_case_form2_residence_town", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 8, Required: true},
		"VictimCaseForm2ResidenceAddress":                      {Name: "VictimCaseForm2ResidenceAddress", DBName: "victim_case_form2_residence_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 150, Required: true},
		"VictimCaseForm2ResidenceZone":                         {Name: "VictimCaseForm2ResidenceZone", DBName: "victim_case_form2_residence_zone", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2SupportContactNames":                   {Name: "VictimCaseForm2SupportContactNames", DBName: "victim_case_form2_support_contact_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: false},
		"VictimCaseForm2SupportContactPhone":                   {Name: "VictimCaseForm2SupportContactPhone", DBName: "victim_case_form2_support_contact_phone", Alias: "", ModelType: "uint", MinSize: 1000000000, MaxSize: 9999999999, Required: false},
		"VictimCaseForm2SupportContactEmail":                   {Name: "VictimCaseForm2SupportContactEmail", DBName: "victim_case_form2_support_contact_email", Alias: "", ModelType: "email", MinSize: 3, MaxSize: 128, Required: false},
		"VictimCaseForm2SupportContactKinship":                 {Name: "VictimCaseForm2SupportContactKinship", DBName: "victim_case_form2_support_contact_kinship", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2SalivaManagementExplanation":           {Name: "VictimCaseForm2SalivaManagementExplanation", DBName: "victim_case_form2_saliva_management_explanation", Alias: "", ModelType: "string", MinSize: 5, MaxSize: 20000, Required: true},
		"VictimCaseForm2ActivitiesUnableToHear":                {Name: "VictimCaseForm2ActivitiesUnableToHear", DBName: "victim_case_form2_activities_unable_to_hear", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToTalk":                {Name: "VictimCaseForm2ActivitiesUnableToTalk", DBName: "victim_case_form2_activities_unable_to_talk", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToSee":                 {Name: "VictimCaseForm2ActivitiesUnableToSee", DBName: "victim_case_form2_activities_unable_to_see", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToMove":                {Name: "VictimCaseForm2ActivitiesUnableToMove", DBName: "victim_case_form2_activities_unable_to_move", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToTake":                {Name: "VictimCaseForm2ActivitiesUnableToTake", DBName: "victim_case_form2_activities_unable_to_take", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToUnderstand":          {Name: "VictimCaseForm2ActivitiesUnableToUnderstand", DBName: "victim_case_form2_activities_unable_to_understand", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToEat":                 {Name: "VictimCaseForm2ActivitiesUnableToEat", DBName: "victim_case_form2_activities_unable_to_eat", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToInteract":            {Name: "VictimCaseForm2ActivitiesUnableToInteract", DBName: "victim_case_form2_activities_unable_to_interact", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2ActivitiesUnableToDoEveryday":          {Name: "VictimCaseForm2ActivitiesUnableToDoEveryday", DBName: "victim_case_form2_activities_unable_to_do_everyday", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 4, Required: true},
		"VictimCaseForm2PersonWithDisability":                  {Name: "VictimCaseForm2PersonWithDisability", DBName: "victim_case_form2_person_with_disability", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2RequireLanguageInterpreter":            {Name: "VictimCaseForm2RequireLanguageInterpreter", DBName: "victim_case_form2_require_language_interpreter", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2LanguageAssistance":                    {Name: "VictimCaseForm2LanguageAssistance", DBName: "victim_case_form2_language_assistance", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimCaseForm2WorkplaceSectorOccurrence":             {Name: "VictimCaseForm2WorkplaceSectorOccurrence", DBName: "victim_case_form2_workplace_sector_occurrence", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2ViolenceMotivatedByGender":             {Name: "VictimCaseForm2ViolenceMotivatedByGender", DBName: "victim_case_form2_violence_motivated_by_gender", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2AttentionWasAppropriate":               {Name: "VictimCaseForm2AttentionWasAppropriate", DBName: "victim_case_form2_attention_was_appropriate", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorOccupation":                   {Name: "VictimCaseForm2AggressorOccupation", DBName: "victim_case_form2_aggressors_occupation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},

		"VictimCaseForm2StoppedSeekingHelp":                           {Name: "VictimCaseForm2StoppedSeekingHelp", DBName: "victim_case_form2_stopped_seeking_help", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2VictimHealthToBlackmail":                      {Name: "VictimCaseForm2VictimHealthToBlackmail", DBName: "victim_case_form2_victim_health_to_blackmail", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2ThreatenedRevealSexualOrientation":            {Name: "VictimCaseForm2ThreatenedRevealSexualOrientation", DBName: "victim_case_form2_threatened_reveal_sexual_orientation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2ViolenceMotivatedByGender2":                   {Name: "VictimCaseForm2ViolenceMotivatedByGender2", DBName: "victim_case_form2_violence_motivated_by_gender_2", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability": {Name: "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", DBName: "victim_case_form2_aggressor_taken_advantage_physical_vulnerabil", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorSexuallyHarassment2":                 {Name: "VictimCaseForm2AggressorSexuallyHarassment2", DBName: "victim_case_form2_aggressor_sexually_harassment_2", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AggressorUsedPositionAuthority":               {Name: "VictimCaseForm2AggressorUsedPositionAuthority", DBName: "victim_case_form2_aggressor_used_position_authority", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimCaseForm2AllowsEasyReport":                             {Name: "VictimCaseForm2AllowsEasyReport", DBName: "victim_case_form2_allows_easy_report", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},

		"VictimCaseForm2RiskScore": {Name: "VictimCaseForm2RiskScore", DBName: "victim_case_form2_risk_score", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2RiskLevel": {Name: "VictimCaseForm2RiskLevel", DBName: "victim_case_form2_risk_level", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 0, Required: true},

		"VictimCaseForm2VictimCase": {Name: "VictimCaseForm2VictimCase", DBName: "victim_case_form2_victim_case", Alias: "", ModelType: "uint", Required: true},
	}
)

// ---------------------------------------------------------------------------
//  DTO – JSON representation
// ---------------------------------------------------------------------------

type VictimCaseForm2DTO struct {
	VictimCaseForm2Id                                    uint64                  `json:"-"`
	VictimCaseForm2ICode                                 string                  `json:"icode"`
	VictimCaseForm2CreationDate                          time.Time               `json:"creationDate"`
	VictimCaseForm2UpdateDate                            time.Time               `json:"updateDate"`
	VictimCaseForm2IdentityName                          string                  `json:"identityName"`
	VictimCaseForm2VictimPhone                           uint64                  `json:"phone"`
	VictimCaseForm2FactsDescription                      string                  `json:"factsDescription"`
	VictimCaseForm2FactsDate                             time.Time               `json:"factsDate"`
	VictimCaseForm2FactsStartTime                        time.Time               `json:"factsStartTime"`
	VictimCaseForm2FactsTownCode                         string                  `json:"factsTownCode"`
	VictimCaseForm2FactsZone                             VictimCaseForm2EnumsDTO `json:"factsZone"`
	VictimCaseForm2FactsAddress                          string                  `json:"factsAddress"`
	VictimCaseForm2ScenarioViolence                      VictimCaseForm2EnumsDTO `json:"scenarioViolence"`
	VictimCaseForm2ReportedPreviously                    VictimCaseForm2EnumsDTO `json:"reportedPreviously"`
	VictimCaseForm2RecurrenceAggression                  VictimCaseForm2EnumsDTO `json:"recurrenceAggression"`
	VictimCaseForm2NumAgressors                          VictimCaseForm2EnumsDTO `json:"numAgressors"`
	VictimCaseForm2ProximityPrincipalAggressor           VictimCaseForm2EnumsDTO `json:"proximityPrincipalAggressor"`
	VictimCaseForm2RelationshipWithPresumedAggressor     VictimCaseForm2EnumsDTO `json:"relationshipWithPresumedAggressor"`
	VictimCaseForm2EconomicallyDependent                 VictimCaseForm2EnumsDTO `json:"economicallyDependent"`
	VictimCaseForm2AggressorGenderIdentity               VictimCaseForm2EnumsDTO `json:"aggressorGenderIdentity"`
	VictimCaseForm2AggressorNames                        string                  `json:"aggressorNames"`
	VictimCaseForm2AggressorDocType                      VictimCaseForm2EnumsDTO `json:"aggressorDocType"`
	VictimCaseForm2AggressorDocNumber                    string                  `json:"aggressorDocNumber"`
	VictimCaseForm2AggressorAddress                      string                  `json:"aggressorAddress"`
	VictimCaseForm2AggressorPhone                        uint64                  `json:"aggressorPhone"`
	VictimCaseForm2AggressorViolencePhysicalIncrease     VictimCaseForm2EnumsDTO `json:"aggressorViolencePhysicalIncrease"`
	VictimCaseForm2AggressorWeaponUsed                   VictimCaseForm2EnumsDTO `json:"aggressorWeaponUsed"`
	VictimCaseForm2AggressorThreatKill                   VictimCaseForm2EnumsDTO `json:"aggressorThreatKill"`
	VictimCaseForm2AggressorPursuesSpiesDestroys         VictimCaseForm2EnumsDTO `json:"aggressorPursuesSpiesDestroys"`
	VictimCaseForm2AggressorCapableOfKilling             VictimCaseForm2EnumsDTO `json:"aggressorCapableOfKilling"`
	VictimCaseForm2AggressorHasAccessToWeapons           VictimCaseForm2EnumsDTO `json:"aggressorHasAccessToWeapons"`
	VictimCaseForm2PartnerUnemployed                     VictimCaseForm2EnumsDTO `json:"partnerUnemployed"`
	VictimCaseForm2PartnerOtherDenunciations             VictimCaseForm2EnumsDTO `json:"partnerOtherDenunciations"`
	VictimCaseForm2AggressorHasPenalBackground           VictimCaseForm2EnumsDTO `json:"aggressorHasPenalBackground"`
	VictimCaseForm2AggressorForcedSex                    VictimCaseForm2EnumsDTO `json:"aggressorForcedSex"`
	VictimCaseForm2AggressorAttemptedStrangulation       VictimCaseForm2EnumsDTO `json:"aggressorAttemptedStrangulation"`
	VictimCaseForm2AggressorConsumesDrugs                VictimCaseForm2EnumsDTO `json:"aggressorConsumesDrugs"`
	VictimCaseForm2AggressorIsAlcoholic                  VictimCaseForm2EnumsDTO `json:"aggressorIsAlcoholic"`
	VictimCaseForm2PartnerControls                       VictimCaseForm2EnumsDTO `json:"partnerControls"`
	VictimCaseForm2AggressorHadHitInVulnerability        VictimCaseForm2EnumsDTO `json:"aggressorHadHitInVulnerability"`
	VictimCaseForm2PartnerThreatenedSuicide              VictimCaseForm2EnumsDTO `json:"partnerThreatenedSuicide"`
	VictimCaseForm2PartnerThreatenedDamageMembers        VictimCaseForm2EnumsDTO `json:"partnerThreatenedDamageMembers"`
	VictimCaseForm2ThoughtsOfSelfHarm                    VictimCaseForm2EnumsDTO `json:"thoughtsOfSelfHarm"`
	VictimCaseForm2AggressorLimitsContactSupportNetworks VictimCaseForm2EnumsDTO `json:"aggressorLimitsContactSupportNetworks"`
	VictimCaseForm2StillLivesWithAggressor               VictimCaseForm2EnumsDTO `json:"stillLivesWithAggressor"`
	VictimCaseForm2AggressorViolentlyJealous             VictimCaseForm2EnumsDTO `json:"aggressorViolentlyJealous"`
	VictimCaseForm2AggressorUnemployed                   VictimCaseForm2EnumsDTO `json:"aggressorUnemployed"`
	VictimCaseForm2AggressorHasPenalBackground2          VictimCaseForm2EnumsDTO `json:"aggressorHasPenalBackground2"`
	VictimCaseForm2AggressorSexuallyHarassment           VictimCaseForm2EnumsDTO `json:"aggressorSexuallyHarassment"`
	VictimCaseForm2AggressorUseDrugs                     VictimCaseForm2EnumsDTO `json:"aggressorUseDrugs"`
	VictimCaseForm2AggressorIsAlcoholic2                 VictimCaseForm2EnumsDTO `json:"aggressorIsAlcoholic2"`
	VictimCaseForm2AggressorControls                     VictimCaseForm2EnumsDTO `json:"aggressorControls"`
	VictimCaseForm2AggressorThreatenedDamageMembers      VictimCaseForm2EnumsDTO `json:"aggressorThreatenedDamageMembers"`
	VictimCaseForm2ThoughtsOfSelfHarm2                   VictimCaseForm2EnumsDTO `json:"thoughtsOfSelfHarm2"`
	VictimCaseForm2AggressorCommonSpaces                 VictimCaseForm2EnumsDTO `json:"aggressorCommonSpaces"`
	VictimCaseForm2AggressorHierarchy                    VictimCaseForm2EnumsDTO `json:"aggressorHierarchy"`
	VictimCaseForm2BirthDate                             time.Time               `json:"birthDate"`
	VictimCaseForm2PhysicalMentalSensoryDifficulties     VictimCaseForm2EnumsDTO `json:"physicalMentalSensoryDifficulties"`
	VictimCaseForm2Nationality                           VictimCaseForm2EnumsDTO `json:"nationality"`
	VictimCaseForm2SpecifiedNationality                  VictimCaseForm2EnumsDTO `json:"specifiedNationality"`
	VictimCaseForm2MigrationCondition                    VictimCaseForm2EnumsDTO `json:"migrationCondition"`
	VictimCaseForm2GenderIdentity                        VictimCaseForm2EnumsDTO `json:"genderIdentity"`
	VictimCaseForm2SexualOrientation                     VictimCaseForm2EnumsDTO `json:"sexualOrientation"`
	VictimCaseForm2AssignedSexAtBirth                    VictimCaseForm2EnumsDTO `json:"assignedSexAtBirth"`
	VictimCaseForm2EthnicAffiliation                     VictimCaseForm2EnumsDTO `json:"ethnicAffiliation"`
	VictimCaseForm2IndigenousPeople                      VictimCaseForm2EnumsDTO `json:"indigenousPeople"`
	VictimCaseForm2CampesinoRecognition                  VictimCaseForm2EnumsDTO `json:"campesinoRecognition"`
	VictimCaseForm2MaritalStatus                         VictimCaseForm2EnumsDTO `json:"maritalStatus"`
	VictimCaseForm2LastEducationLevel                    VictimCaseForm2EnumsDTO `json:"lastEducationLevel"`
	VictimCaseForm2Occupation                            VictimCaseForm2EnumsDTO `json:"occupation"`
	VictimCaseForm2IncomeGenerationMethod                VictimCaseForm2EnumsDTO `json:"incomeGenerationMethod"`
	VictimCaseForm2EmploymentRelationship                VictimCaseForm2EnumsDTO `json:"employmentRelationship"`
	VictimCaseForm2ApproxStartAsp                        time.Time               `json:"approxStartAsp"`
	VictimCaseForm2HousingTenancyForm                    VictimCaseForm2EnumsDTO `json:"housingTenancyForm"`
	VictimCaseForm2HousingStratum                        VictimCaseForm2EnumsDTO `json:"housingStratum"`
	VictimCaseForm2CurrentlyPregnant                     VictimCaseForm2EnumsDTO `json:"currentlyPregnant"`
	VictimCaseForm2ResidenceTownCode                     string                  `json:"residenceTownCode"`
	VictimCaseForm2ResidenceAddress                      string                  `json:"residenceAddress"`
	VictimCaseForm2ResidenceZone                         VictimCaseForm2EnumsDTO `json:"residenceZone"`
	VictimCaseForm2SupportContactNames                   string                  `json:"supportContactNames"`
	VictimCaseForm2SupportContactPhone                   uint64                  `json:"supportContactPhone"`
	VictimCaseForm2SupportContactEmail                   string                  `json:"supportContactEmail"`
	VictimCaseForm2SupportContactKinship                 VictimCaseForm2EnumsDTO `json:"supportContactKinship"`
	VictimCaseForm2LanguageAssistance                    string                  `json:"languageAssistance"`

	VictimCaseForm2PersonWithDisability       VictimCaseForm2EnumsDTO `json:"personWithDisability"`
	VictimCaseForm2RequireLanguageInterpreter VictimCaseForm2EnumsDTO `json:"requireLanguageInterpreter"`
	VictimCaseForm2WorkplaceSectorOccurrence  VictimCaseForm2EnumsDTO `json:"workplaceSectorOccurrence"`
	VictimCaseForm2ViolenceMotivatedByGender  VictimCaseForm2EnumsDTO `json:"violenceMotivatedByGender"`
	VictimCaseForm2AttentionWasAppropriate    VictimCaseForm2EnumsDTO `json:"attentionWasAppropriate"`
	VictimCaseForm2AggressorOccupation        VictimCaseForm2EnumsDTO `json:"aggressorOccupation"`

	VictimCaseForm2StoppedSeekingHelp                           VictimCaseForm2EnumsDTO `json:"stoppedSeekingHelp"`
	VictimCaseForm2VictimHealthToBlackmail                      VictimCaseForm2EnumsDTO `json:"victimHealthToBlackmail"`
	VictimCaseForm2ThreatenedRevealSexualOrientation            VictimCaseForm2EnumsDTO `json:"threatenedRevealSexualOrientation"`
	VictimCaseForm2ViolenceMotivatedByGender2                   VictimCaseForm2EnumsDTO `json:"violenceMotivatedByGender2"`
	VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability VictimCaseForm2EnumsDTO `json:"aggressorTakenAdvantagePhysicalVulnerability"`
	VictimCaseForm2AggressorSexuallyHarassment2                 VictimCaseForm2EnumsDTO `json:"aggressorSexuallyHarassment2"`
	VictimCaseForm2AggressorUsedPositionAuthority               VictimCaseForm2EnumsDTO `json:"aggressorUsedPositionAuthority"`
	VictimCaseForm2AllowsEasyReport                             VictimCaseForm2EnumsDTO `json:"allowsEasyReport"`

	//Campos de calificación
	VictimCaseForm2ActivitiesUnableToHear       int64 `json:"activitiesUnableToHear"`
	VictimCaseForm2ActivitiesUnableToTalk       int64 `json:"activitiesUnableToTalk"`
	VictimCaseForm2ActivitiesUnableToSee        int64 `json:"activitiesUnableToSee"`
	VictimCaseForm2ActivitiesUnableToMove       int64 `json:"activitiesUnableToMove"`
	VictimCaseForm2ActivitiesUnableToTake       int64 `json:"activitiesUnableToTake"`
	VictimCaseForm2ActivitiesUnableToUnderstand int64 `json:"activitiesUnableToUnderstand"`
	VictimCaseForm2ActivitiesUnableToEat        int64 `json:"activitiesUnableToEat"`
	VictimCaseForm2ActivitiesUnableToInteract   int64 `json:"activitiesUnableToInteract"`
	VictimCaseForm2ActivitiesUnableToDoEveryday int64 `json:"activitiesUnableToDoEveryday"`

	//Campos de selección múltiple
	VictimCaseForm2TypeViolenceExperienced      []VictimCaseForm2EnumsDTO `json:"typeViolenceExperienced"`
	VictimCaseForm2SubtypeViolenceExperienced   []VictimCaseForm2EnumsDTO `json:"subtypeViolenceExperienced"`
	VictimCaseForm2ScopeOfViolence              []VictimCaseForm2EnumsDTO `json:"scopeOfViolence"`
	VictimCaseForm2WhoReportTo                  []VictimCaseForm2EnumsDTO `json:"whoReportTo"`
	VictimCaseForm2ActivitiesUnableToPerform    []VictimCaseForm2EnumsDTO `json:"activitiesUnableToPerform"`
	VictimCaseForm2AdjustmentsGBV               []VictimCaseForm2EnumsDTO `json:"adjustmentsGBV"`
	VictimCaseForm2Law1996                      []VictimCaseForm2EnumsDTO `json:"law1996"`
	VictimCaseForm2SpeciallyProtectedPopulation []VictimCaseForm2EnumsDTO `json:"speciallyProtectedPopulation"`
	VictimCaseForm2ASPMode                      []VictimCaseForm2EnumsDTO `json:"aspMode"`
	VictimCaseForm2ReasonASP                    []VictimCaseForm2EnumsDTO `json:"reasonASP"`
	VictimCaseForm2HasDependents                []VictimCaseForm2EnumsDTO `json:"hasDependents"`
	VictimCaseForm2ActionPlan                   []VictimCaseForm2EnumsDTO `json:"actionPlan"`

	VictimCaseForm2SalivaManagementExplanation string `json:"salivaManagementExplanation"`

	VictimCaseForm2RiskScore int64 `json:"riskScore"`
	VictimCaseForm2RiskLevel int64 `json:"riskLevel"`

	// Reference to the parent case – stored as ID in DB
	VictimCaseForm2VictimCase interface{} `json:"case"` // use VictimCaseDTO internally

	//Campos de formulario que no hacen parte del modelo o no directamente en la BD -----------------------------------

	VictimCaseForm2FactsDepartment security_daos.DepartmentDTO `json:"factsDepartment"`
	VictimCaseForm2FactsTown       security_daos.TownDTO       `json:"factsTown"`
	VictimCaseForm2FactsCity       security_daos.CityDTO       `json:"factsCity"`

	VictimCaseForm2ResidenceDepartment security_daos.DepartmentDTO `json:"residenceDepartment"`
	VictimCaseForm2ResidenceTown       security_daos.TownDTO       `json:"residenceTown"`
	VictimCaseForm2ResidenceCity       security_daos.CityDTO       `json:"residenceCity"`
}

type VictimCaseForm2PgDB struct {
	VictimCaseForm2Id               sql.NullInt64
	VictimCaseForm2ICode            sql.NullString
	VictimCaseForm2CreationDate     sql.NullTime
	VictimCaseForm2UpdateDate       sql.NullTime
	VictimCaseForm2IdentityName     sql.NullString
	VictimCaseForm2VictimPhone      sql.NullInt64
	VictimCaseForm2FactsDescription sql.NullString
	VictimCaseForm2FactsDate        sql.NullTime
	VictimCaseForm2FactsStartTime   sql.NullString
	VictimCaseForm2FactsTownCode    sql.NullString
	VictimCaseForm2FactsZone        sql.NullInt64
	VictimCaseForm2FactsAddress     sql.NullString

	//Campos de calificación
	VictimCaseForm2ActivitiesUnableToHear       sql.NullInt64
	VictimCaseForm2ActivitiesUnableToTalk       sql.NullInt64
	VictimCaseForm2ActivitiesUnableToSee        sql.NullInt64
	VictimCaseForm2ActivitiesUnableToMove       sql.NullInt64
	VictimCaseForm2ActivitiesUnableToTake       sql.NullInt64
	VictimCaseForm2ActivitiesUnableToUnderstand sql.NullInt64
	VictimCaseForm2ActivitiesUnableToEat        sql.NullInt64
	VictimCaseForm2ActivitiesUnableToInteract   sql.NullInt64
	VictimCaseForm2ActivitiesUnableToDoEveryday sql.NullInt64

	// Enums – stored as ID in DB, but exposed as a struct
	VictimCaseForm2ScenarioViolence                      sql.NullInt64
	VictimCaseForm2ReportedPreviously                    sql.NullInt64
	VictimCaseForm2RecurrenceAggression                  sql.NullInt64
	VictimCaseForm2NumAgressors                          sql.NullInt64
	VictimCaseForm2ProximityPrincipalAggressor           sql.NullInt64
	VictimCaseForm2RelationshipWithPresumedAggressor     sql.NullInt64
	VictimCaseForm2EconomicallyDependent                 sql.NullInt64
	VictimCaseForm2AggressorGenderIdentity               sql.NullInt64
	VictimCaseForm2AggressorNames                        sql.NullString
	VictimCaseForm2AggressorDocType                      sql.NullInt64
	VictimCaseForm2AggressorDocNumber                    sql.NullString
	VictimCaseForm2AggressorAddress                      sql.NullString
	VictimCaseForm2AggressorPhone                        sql.NullInt64
	VictimCaseForm2AggressorViolencePhysicalIncrease     sql.NullInt64
	VictimCaseForm2AggressorWeaponUsed                   sql.NullInt64
	VictimCaseForm2AggressorThreatKill                   sql.NullInt64
	VictimCaseForm2AggressorPursuesSpiesDestroys         sql.NullInt64
	VictimCaseForm2AggressorCapableOfKilling             sql.NullInt64
	VictimCaseForm2AggressorHasAccessToWeapons           sql.NullInt64
	VictimCaseForm2PartnerUnemployed                     sql.NullInt64
	VictimCaseForm2PartnerOtherDenunciations             sql.NullInt64
	VictimCaseForm2AggressorHasPenalBackground           sql.NullInt64
	VictimCaseForm2AggressorForcedSex                    sql.NullInt64
	VictimCaseForm2AggressorAttemptedStrangulation       sql.NullInt64
	VictimCaseForm2AggressorConsumesDrugs                sql.NullInt64
	VictimCaseForm2AggressorIsAlcoholic                  sql.NullInt64
	VictimCaseForm2PartnerControls                       sql.NullInt64
	VictimCaseForm2AggressorHadHitInVulnerability        sql.NullInt64
	VictimCaseForm2PartnerThreatenedSuicide              sql.NullInt64
	VictimCaseForm2PartnerThreatenedDamageMembers        sql.NullInt64
	VictimCaseForm2ThoughtsOfSelfHarm                    sql.NullInt64
	VictimCaseForm2AggressorLimitsContactSupportNetworks sql.NullInt64
	VictimCaseForm2StillLivesWithAggressor               sql.NullInt64
	VictimCaseForm2AggressorViolentlyJealous             sql.NullInt64
	VictimCaseForm2AggressorUnemployed                   sql.NullInt64
	VictimCaseForm2AggressorHasPenalBackground2          sql.NullInt64
	VictimCaseForm2AggressorSexuallyHarassment           sql.NullInt64
	VictimCaseForm2AggressorUseDrugs                     sql.NullInt64
	VictimCaseForm2AggressorIsAlcoholic2                 sql.NullInt64
	VictimCaseForm2AggressorControls                     sql.NullInt64
	VictimCaseForm2AggressorThreatenedDamageMembers      sql.NullInt64
	VictimCaseForm2ThoughtsOfSelfHarm2                   sql.NullInt64
	VictimCaseForm2AggressorCommonSpaces                 sql.NullInt64
	VictimCaseForm2AggressorHierarchy                    sql.NullInt64
	VictimCaseForm2BirthDate                             sql.NullTime
	VictimCaseForm2PhysicalMentalSensoryDifficulties     sql.NullInt64
	VictimCaseForm2Nationality                           sql.NullInt64
	VictimCaseForm2SpecifiedNationality                  sql.NullInt64
	VictimCaseForm2MigrationCondition                    sql.NullInt64
	VictimCaseForm2GenderIdentity                        sql.NullInt64
	VictimCaseForm2SexualOrientation                     sql.NullInt64
	VictimCaseForm2AssignedSexAtBirth                    sql.NullInt64
	VictimCaseForm2EthnicAffiliation                     sql.NullInt64
	VictimCaseForm2IndigenousPeople                      sql.NullInt64
	VictimCaseForm2CampesinoRecognition                  sql.NullInt64
	VictimCaseForm2MaritalStatus                         sql.NullInt64
	VictimCaseForm2LastEducationLevel                    sql.NullInt64
	VictimCaseForm2Occupation                            sql.NullInt64
	VictimCaseForm2IncomeGenerationMethod                sql.NullInt64
	VictimCaseForm2EmploymentRelationship                sql.NullInt64
	VictimCaseForm2ApproxStartAsp                        sql.NullTime
	VictimCaseForm2HousingTenancyForm                    sql.NullInt64
	VictimCaseForm2HousingStratum                        sql.NullInt64
	VictimCaseForm2CurrentlyPregnant                     sql.NullInt64
	VictimCaseForm2ResidenceTownCode                     sql.NullString
	VictimCaseForm2ResidenceAddress                      sql.NullString
	VictimCaseForm2ResidenceZone                         sql.NullInt64
	VictimCaseForm2SupportContactNames                   sql.NullString
	VictimCaseForm2SupportContactPhone                   sql.NullInt64
	VictimCaseForm2SupportContactEmail                   sql.NullString
	VictimCaseForm2SupportContactKinship                 sql.NullInt64
	VictimCaseForm2LanguageAssistance                    sql.NullString

	VictimCaseForm2StoppedSeekingHelp                           sql.NullInt64
	VictimCaseForm2VictimHealthToBlackmail                      sql.NullInt64
	VictimCaseForm2ThreatenedRevealSexualOrientation            sql.NullInt64
	VictimCaseForm2ViolenceMotivatedByGender2                   sql.NullInt64
	VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability sql.NullInt64
	VictimCaseForm2AggressorSexuallyHarassment2                 sql.NullInt64
	VictimCaseForm2AggressorUsedPositionAuthority               sql.NullInt64
	VictimCaseForm2AllowsEasyReport                             sql.NullInt64

	VictimCaseForm2RiskScore sql.NullInt64
	VictimCaseForm2RiskLevel sql.NullInt64

	VictimCaseForm2PersonWithDisability       sql.NullInt64
	VictimCaseForm2RequireLanguageInterpreter sql.NullInt64
	VictimCaseForm2WorkplaceSectorOccurrence  sql.NullInt64
	VictimCaseForm2ViolenceMotivatedByGender  sql.NullInt64
	VictimCaseForm2AttentionWasAppropriate    sql.NullInt64
	VictimCaseForm2AggressorOccupation        sql.NullInt64

	VictimCaseForm2SalivaManagementExplanation sql.NullString

	// Reference to the parent case – stored as ID in DB
	VictimCaseForm2VictimCase sql.NullInt64
}

func (vcd VictimCaseForm2DTO) MarshalJSON() ([]byte, error) {
	type Alias VictimCaseForm2DTO

	return json.Marshal(&struct {
		*Alias
		VictimCaseForm2CreationDate   string `json:"creationDate"`
		VictimCaseForm2UpdateDate     string `json:"updateDate"`
		VictimCaseForm2BirthDate      string `json:"birthDate"`
		VictimCaseForm2FactsDate      string `json:"factsDate"`
		VictimCaseForm2FactsStartTime string `json:"factsStartTime"`
		VictimCaseForm2ApproxStartAsp string `json:"approxStartAsp"`
	}{
		Alias:                       (*Alias)(&vcd),
		VictimCaseForm2CreationDate: vcd.VictimCaseForm2CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimCaseForm2UpdateDate:   vcd.VictimCaseForm2UpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimCaseForm2BirthDate:    vcd.VictimCaseForm2BirthDate.Format(common_config.DateTime.DATE_FORMAT),

		VictimCaseForm2FactsDate:      vcd.VictimCaseForm2FactsDate.Format(common_config.DateTime.DATE_FORMAT),
		VictimCaseForm2FactsStartTime: vcd.VictimCaseForm2FactsStartTime.Format(common_config.DateTime.TIME_FORMAT),
		VictimCaseForm2ApproxStartAsp: vcd.VictimCaseForm2ApproxStartAsp.Format(common_config.DateTime.DATE_FORMAT),
	})

}

func (vcd *VictimCaseForm2DTO) UnmarshalJSON(data []byte) error {
	type Alias VictimCaseForm2DTO

	aux := &struct {
		*Alias
		VictimCaseForm2CreationDate   string `json:"creationDate"`
		VictimCaseForm2UpdateDate     string `json:"updateDate"`
		VictimCaseForm2BirthDate      string `json:"birthDate"`
		VictimCaseForm2FactsDate      string `json:"factsDate"`
		VictimCaseForm2FactsStartTime string `json:"factsStartTime"`
		VictimCaseForm2ApproxStartAsp string `json:"approxStartAsp"`
	}{
		Alias: (*Alias)(vcd),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{} // invalid → zero value
		}
		return t
	}

	vcd.VictimCaseForm2CreationDate = parse(aux.VictimCaseForm2CreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimCaseForm2UpdateDate = parse(aux.VictimCaseForm2UpdateDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimCaseForm2BirthDate = parse(aux.VictimCaseForm2BirthDate, common_config.DateTime.DATE_FORMAT)
	vcd.VictimCaseForm2FactsDate = parse(aux.VictimCaseForm2FactsDate, common_config.DateTime.DATE_FORMAT)
	vcd.VictimCaseForm2FactsStartTime = parse(aux.VictimCaseForm2FactsStartTime, common_config.DateTime.TIME_FORMAT)
	vcd.VictimCaseForm2ApproxStartAsp = parse(aux.VictimCaseForm2ApproxStartAsp, common_config.DateTime.DATE_FORMAT)

	return nil
}

func SetVictimCaseForm2(victimCaseForm2 *VictimCaseForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fields := []string{
		"VictimCaseForm2ICode", "VictimCaseForm2CreationDate", "VictimCaseForm2UpdateDate", "VictimCaseForm2IdentityName", "VictimCaseForm2VictimPhone",
		"VictimCaseForm2FactsDescription", "VictimCaseForm2FactsDate", "VictimCaseForm2FactsStartTime", "VictimCaseForm2FactsTownCode", "VictimCaseForm2FactsZone",
		"VictimCaseForm2FactsAddress", "VictimCaseForm2ScenarioViolence", "VictimCaseForm2ReportedPreviously", "VictimCaseForm2RecurrenceAggression", "VictimCaseForm2NumAgressors",
		"VictimCaseForm2ProximityPrincipalAggressor", "VictimCaseForm2RelationshipWithPresumedAggressor", "VictimCaseForm2EconomicallyDependent", "VictimCaseForm2AggressorGenderIdentity", "VictimCaseForm2AggressorNames",
		"VictimCaseForm2AggressorDocType", "VictimCaseForm2AggressorDocNumber", "VictimCaseForm2AggressorAddress", "VictimCaseForm2AggressorPhone", "VictimCaseForm2AggressorViolencePhysicalIncrease",
		"VictimCaseForm2AggressorWeaponUsed", "VictimCaseForm2AggressorThreatKill", "VictimCaseForm2AggressorPursuesSpiesDestroys", "VictimCaseForm2AggressorCapableOfKilling", "VictimCaseForm2AggressorHasAccessToWeapons",
		"VictimCaseForm2PartnerUnemployed", "VictimCaseForm2PartnerOtherDenunciations", "VictimCaseForm2AggressorHasPenalBackground", "VictimCaseForm2AggressorForcedSex", "VictimCaseForm2AggressorAttemptedStrangulation",
		"VictimCaseForm2AggressorConsumesDrugs", "VictimCaseForm2AggressorIsAlcoholic", "VictimCaseForm2PartnerControls", "VictimCaseForm2AggressorHadHitInVulnerability", "VictimCaseForm2PartnerThreatenedSuicide",
		"VictimCaseForm2PartnerThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm", "VictimCaseForm2AggressorLimitsContactSupportNetworks", "VictimCaseForm2StillLivesWithAggressor", "VictimCaseForm2AggressorViolentlyJealous",
		"VictimCaseForm2AggressorUnemployed", "VictimCaseForm2AggressorHasPenalBackground2", "VictimCaseForm2AggressorSexuallyHarassment", "VictimCaseForm2AggressorUseDrugs", "VictimCaseForm2AggressorIsAlcoholic2",
		"VictimCaseForm2AggressorControls", "VictimCaseForm2AggressorThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm2", "VictimCaseForm2AggressorCommonSpaces", "VictimCaseForm2AggressorHierarchy",
		"VictimCaseForm2BirthDate", "VictimCaseForm2PhysicalMentalSensoryDifficulties", "VictimCaseForm2Nationality", "VictimCaseForm2SpecifiedNationality", "VictimCaseForm2MigrationCondition",
		"VictimCaseForm2GenderIdentity", "VictimCaseForm2SexualOrientation", "VictimCaseForm2AssignedSexAtBirth", "VictimCaseForm2EthnicAffiliation", "VictimCaseForm2IndigenousPeople",
		"VictimCaseForm2CampesinoRecognition", "VictimCaseForm2MaritalStatus", "VictimCaseForm2LastEducationLevel", "VictimCaseForm2Occupation", "VictimCaseForm2IncomeGenerationMethod", "VictimCaseForm2EmploymentRelationship",
		"VictimCaseForm2ApproxStartAsp", "VictimCaseForm2HousingTenancyForm", "VictimCaseForm2HousingStratum", "VictimCaseForm2CurrentlyPregnant", "VictimCaseForm2ResidenceTownCode",
		"VictimCaseForm2ResidenceAddress", "VictimCaseForm2ResidenceZone", "VictimCaseForm2SupportContactNames", "VictimCaseForm2SupportContactPhone", "VictimCaseForm2SupportContactEmail",
		"VictimCaseForm2SupportContactKinship", "VictimCaseForm2LanguageAssistance",
		"VictimCaseForm2PersonWithDisability", "VictimCaseForm2RequireLanguageInterpreter", "VictimCaseForm2WorkplaceSectorOccurrence", "VictimCaseForm2ViolenceMotivatedByGender", "VictimCaseForm2AttentionWasAppropriate", "VictimCaseForm2AggressorOccupation",
		"VictimCaseForm2SalivaManagementExplanation",
		"VictimCaseForm2ActivitiesUnableToHear", "VictimCaseForm2ActivitiesUnableToTalk", "VictimCaseForm2ActivitiesUnableToSee", "VictimCaseForm2ActivitiesUnableToMove", "VictimCaseForm2ActivitiesUnableToTake",
		"VictimCaseForm2ActivitiesUnableToUnderstand", "VictimCaseForm2ActivitiesUnableToEat", "VictimCaseForm2ActivitiesUnableToInteract", "VictimCaseForm2ActivitiesUnableToDoEveryday",
		"VictimCaseForm2StoppedSeekingHelp", "VictimCaseForm2VictimHealthToBlackmail", "VictimCaseForm2ThreatenedRevealSexualOrientation", "VictimCaseForm2ViolenceMotivatedByGender2", "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", "VictimCaseForm2AggressorSexuallyHarassment2", "VictimCaseForm2AggressorUsedPositionAuthority",
		"VictimCaseForm2AllowsEasyReport", "VictimCaseForm2RiskScore", "VictimCaseForm2RiskLevel",
		"VictimCaseForm2VictimCase",
	}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fields, []string{}, VictimCaseForm2DBName, []string{}, []string{}, []string{"VictimCaseForm2Id"}, common_dao.SQL_AND, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		victimCaseForm2.VictimCaseForm2ICode,
		victimCaseForm2.VictimCaseForm2CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCaseForm2.VictimCaseForm2UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCaseForm2.VictimCaseForm2IdentityName,
		victimCaseForm2.VictimCaseForm2VictimPhone,
		victimCaseForm2.VictimCaseForm2FactsDescription,
		victimCaseForm2.VictimCaseForm2FactsDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		victimCaseForm2.VictimCaseForm2FactsStartTime.Format(common_config.DateTime.TIME_FORMAT),
		victimCaseForm2.VictimCaseForm2FactsTownCode,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2FactsZone.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2FactsAddress),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ScenarioViolence.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ReportedPreviously.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2RecurrenceAggression.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2NumAgressors.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ProximityPrincipalAggressor.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2EconomicallyDependent.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorGenderIdentity.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2AggressorNames,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorDocType.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2AggressorDocNumber,
		victimCaseForm2.VictimCaseForm2AggressorAddress,
		victimCaseForm2.VictimCaseForm2AggressorPhone,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorWeaponUsed.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorThreatKill.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorCapableOfKilling.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerUnemployed.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerOtherDenunciations.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHasPenalBackground.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorForcedSex.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorAttemptedStrangulation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorConsumesDrugs.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorIsAlcoholic.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerControls.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHadHitInVulnerability.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerThreatenedSuicide.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerThreatenedDamageMembers.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorLimitsContactSupportNetworks.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2StillLivesWithAggressor.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorViolentlyJealous.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorUnemployed.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHasPenalBackground2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorUseDrugs.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorIsAlcoholic2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorControls.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorThreatenedDamageMembers.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorCommonSpaces.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHierarchy.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PhysicalMentalSensoryDifficulties.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2Nationality.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2SpecifiedNationality.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2MigrationCondition.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2GenderIdentity.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2SexualOrientation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AssignedSexAtBirth.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2EthnicAffiliation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2IndigenousPeople.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2CampesinoRecognition.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2MaritalStatus.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2LastEducationLevel.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2Occupation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2EmploymentRelationship.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2ApproxStartAsp.Format(common_config.DateTime.DB_DATE_FORMAT),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2HousingTenancyForm.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2HousingStratum.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2CurrentlyPregnant.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2ResidenceTownCode,
		victimCaseForm2.VictimCaseForm2ResidenceAddress,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ResidenceZone.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2SupportContactNames,
		victimCaseForm2.VictimCaseForm2SupportContactPhone,
		victimCaseForm2.VictimCaseForm2SupportContactEmail,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2SupportContactKinship.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2LanguageAssistance,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PersonWithDisability.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2RequireLanguageInterpreter.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2WorkplaceSectorOccurrence.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AttentionWasAppropriate.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorOccupation.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2SalivaManagementExplanation,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToHear,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToTalk,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToSee,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToMove,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToTake,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToUnderstand,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToEat,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToInteract,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToDoEveryday,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2StoppedSeekingHelp.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2VictimHealthToBlackmail.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ThreatenedRevealSexualOrientation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorUsedPositionAuthority.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AllowsEasyReport.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2RiskScore,
		victimCaseForm2.VictimCaseForm2RiskLevel,
		victimCaseForm2.VictimCaseForm2VictimCase.(VictimCaseDTO).VictimCaseId)

	persistenceCtrl.Scan(&victimCaseForm2.VictimCaseForm2Id)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	return nil
}

func GetVictimCaseForm2(by common_controllers.By, victimCaseForm2 *VictimCaseForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2Path string = VictimCaseForm2DBScheme + "." + VictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fields := []string{
		"VictimCaseForm2Id", "VictimCaseForm2ICode", "VictimCaseForm2CreationDate", "VictimCaseForm2UpdateDate", "VictimCaseForm2IdentityName",
		"VictimCaseForm2VictimPhone", "VictimCaseForm2FactsDescription", "VictimCaseForm2FactsDate", "VictimCaseForm2FactsStartTime", "VictimCaseForm2FactsTownCode",
		"VictimCaseForm2FactsZone", "VictimCaseForm2FactsAddress", "VictimCaseForm2ScenarioViolence", "VictimCaseForm2ReportedPreviously", "VictimCaseForm2RecurrenceAggression",
		"VictimCaseForm2NumAgressors", "VictimCaseForm2ProximityPrincipalAggressor", "VictimCaseForm2RelationshipWithPresumedAggressor", "VictimCaseForm2EconomicallyDependent", "VictimCaseForm2AggressorGenderIdentity",
		"VictimCaseForm2AggressorNames", "VictimCaseForm2AggressorDocType", "VictimCaseForm2AggressorDocNumber", "VictimCaseForm2AggressorAddress", "VictimCaseForm2AggressorPhone",
		"VictimCaseForm2AggressorViolencePhysicalIncrease", "VictimCaseForm2AggressorWeaponUsed", "VictimCaseForm2AggressorThreatKill", "VictimCaseForm2AggressorPursuesSpiesDestroys", "VictimCaseForm2AggressorCapableOfKilling",
		"VictimCaseForm2AggressorHasAccessToWeapons", "VictimCaseForm2PartnerUnemployed", "VictimCaseForm2PartnerOtherDenunciations", "VictimCaseForm2AggressorHasPenalBackground", "VictimCaseForm2AggressorForcedSex",
		"VictimCaseForm2AggressorAttemptedStrangulation", "VictimCaseForm2AggressorConsumesDrugs", "VictimCaseForm2AggressorIsAlcoholic", "VictimCaseForm2PartnerControls", "VictimCaseForm2AggressorHadHitInVulnerability",
		"VictimCaseForm2PartnerThreatenedSuicide", "VictimCaseForm2PartnerThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm", "VictimCaseForm2AggressorLimitsContactSupportNetworks", "VictimCaseForm2StillLivesWithAggressor",
		"VictimCaseForm2AggressorViolentlyJealous", "VictimCaseForm2AggressorUnemployed", "VictimCaseForm2AggressorHasPenalBackground2", "VictimCaseForm2AggressorSexuallyHarassment", "VictimCaseForm2AggressorUseDrugs",
		"VictimCaseForm2AggressorIsAlcoholic2", "VictimCaseForm2AggressorControls", "VictimCaseForm2AggressorThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm2", "VictimCaseForm2AggressorCommonSpaces",
		"VictimCaseForm2AggressorHierarchy", "VictimCaseForm2BirthDate", "VictimCaseForm2PhysicalMentalSensoryDifficulties", "VictimCaseForm2Nationality", "VictimCaseForm2SpecifiedNationality",
		"VictimCaseForm2MigrationCondition", "VictimCaseForm2GenderIdentity", "VictimCaseForm2SexualOrientation", "VictimCaseForm2AssignedSexAtBirth", "VictimCaseForm2EthnicAffiliation",
		"VictimCaseForm2IndigenousPeople", "VictimCaseForm2CampesinoRecognition", "VictimCaseForm2MaritalStatus", "VictimCaseForm2LastEducationLevel", "VictimCaseForm2Occupation",
		"VictimCaseForm2IncomeGenerationMethod", "VictimCaseForm2EmploymentRelationship", "VictimCaseForm2ApproxStartAsp", "VictimCaseForm2HousingTenancyForm", "VictimCaseForm2HousingStratum", "VictimCaseForm2CurrentlyPregnant",
		"VictimCaseForm2ResidenceTownCode", "VictimCaseForm2ResidenceAddress", "VictimCaseForm2ResidenceZone", "VictimCaseForm2SupportContactNames", "VictimCaseForm2SupportContactPhone",
		"VictimCaseForm2SupportContactEmail", "VictimCaseForm2SupportContactKinship", "VictimCaseForm2LanguageAssistance",
		"VictimCaseForm2PersonWithDisability", "VictimCaseForm2RequireLanguageInterpreter", "VictimCaseForm2WorkplaceSectorOccurrence", "VictimCaseForm2ViolenceMotivatedByGender", "VictimCaseForm2AttentionWasAppropriate", "VictimCaseForm2AggressorOccupation", "VictimCaseForm2SalivaManagementExplanation",
		"VictimCaseForm2ActivitiesUnableToHear", "VictimCaseForm2ActivitiesUnableToTalk", "VictimCaseForm2ActivitiesUnableToSee", "VictimCaseForm2ActivitiesUnableToMove", "VictimCaseForm2ActivitiesUnableToTake",
		"VictimCaseForm2ActivitiesUnableToUnderstand", "VictimCaseForm2ActivitiesUnableToEat", "VictimCaseForm2ActivitiesUnableToInteract", "VictimCaseForm2ActivitiesUnableToDoEveryday",
		"VictimCaseForm2StoppedSeekingHelp", "VictimCaseForm2VictimHealthToBlackmail", "VictimCaseForm2ThreatenedRevealSexualOrientation", "VictimCaseForm2ViolenceMotivatedByGender2", "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", "VictimCaseForm2AggressorSexuallyHarassment2", "VictimCaseForm2AggressorUsedPositionAuthority",
		"VictimCaseForm2AllowsEasyReport", "VictimCaseForm2RiskScore", "VictimCaseForm2RiskLevel", "VictimCaseForm2VictimCase",
	}

	var fieldsStr string = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fields, []string{}, VictimCaseForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var victimCaseForm2Pg VictimCaseForm2PgDB
	persistenceCtrl.Scan(&victimCaseForm2Pg.VictimCaseForm2Id,
		&victimCaseForm2Pg.VictimCaseForm2ICode,
		&victimCaseForm2Pg.VictimCaseForm2CreationDate,
		&victimCaseForm2Pg.VictimCaseForm2UpdateDate,
		&victimCaseForm2Pg.VictimCaseForm2IdentityName,
		&victimCaseForm2Pg.VictimCaseForm2VictimPhone,
		&victimCaseForm2Pg.VictimCaseForm2FactsDescription,
		&victimCaseForm2Pg.VictimCaseForm2FactsDate,
		&victimCaseForm2Pg.VictimCaseForm2FactsStartTime,
		&victimCaseForm2Pg.VictimCaseForm2FactsTownCode,
		&victimCaseForm2Pg.VictimCaseForm2FactsZone,
		&victimCaseForm2Pg.VictimCaseForm2FactsAddress,
		&victimCaseForm2Pg.VictimCaseForm2ScenarioViolence,
		&victimCaseForm2Pg.VictimCaseForm2ReportedPreviously,
		&victimCaseForm2Pg.VictimCaseForm2RecurrenceAggression,
		&victimCaseForm2Pg.VictimCaseForm2NumAgressors,
		&victimCaseForm2Pg.VictimCaseForm2ProximityPrincipalAggressor,
		&victimCaseForm2Pg.VictimCaseForm2RelationshipWithPresumedAggressor,
		&victimCaseForm2Pg.VictimCaseForm2EconomicallyDependent,
		&victimCaseForm2Pg.VictimCaseForm2AggressorGenderIdentity,
		&victimCaseForm2Pg.VictimCaseForm2AggressorNames,
		&victimCaseForm2Pg.VictimCaseForm2AggressorDocType,
		&victimCaseForm2Pg.VictimCaseForm2AggressorDocNumber,
		&victimCaseForm2Pg.VictimCaseForm2AggressorAddress,
		&victimCaseForm2Pg.VictimCaseForm2AggressorPhone,
		&victimCaseForm2Pg.VictimCaseForm2AggressorViolencePhysicalIncrease,
		&victimCaseForm2Pg.VictimCaseForm2AggressorWeaponUsed,
		&victimCaseForm2Pg.VictimCaseForm2AggressorThreatKill,
		&victimCaseForm2Pg.VictimCaseForm2AggressorPursuesSpiesDestroys,
		&victimCaseForm2Pg.VictimCaseForm2AggressorCapableOfKilling,
		&victimCaseForm2Pg.VictimCaseForm2AggressorHasAccessToWeapons,
		&victimCaseForm2Pg.VictimCaseForm2PartnerUnemployed,
		&victimCaseForm2Pg.VictimCaseForm2PartnerOtherDenunciations,
		&victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground,
		&victimCaseForm2Pg.VictimCaseForm2AggressorForcedSex,
		&victimCaseForm2Pg.VictimCaseForm2AggressorAttemptedStrangulation,
		&victimCaseForm2Pg.VictimCaseForm2AggressorConsumesDrugs,
		&victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic,
		&victimCaseForm2Pg.VictimCaseForm2PartnerControls,
		&victimCaseForm2Pg.VictimCaseForm2AggressorHadHitInVulnerability,
		&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedSuicide,
		&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedDamageMembers,
		&victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm,
		&victimCaseForm2Pg.VictimCaseForm2AggressorLimitsContactSupportNetworks,
		&victimCaseForm2Pg.VictimCaseForm2StillLivesWithAggressor,
		&victimCaseForm2Pg.VictimCaseForm2AggressorViolentlyJealous,
		&victimCaseForm2Pg.VictimCaseForm2AggressorUnemployed,
		&victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground2,
		&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment,
		&victimCaseForm2Pg.VictimCaseForm2AggressorUseDrugs,
		&victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic2,
		&victimCaseForm2Pg.VictimCaseForm2AggressorControls,
		&victimCaseForm2Pg.VictimCaseForm2AggressorThreatenedDamageMembers,
		&victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm2,
		&victimCaseForm2Pg.VictimCaseForm2AggressorCommonSpaces,
		&victimCaseForm2Pg.VictimCaseForm2AggressorHierarchy,
		&victimCaseForm2Pg.VictimCaseForm2BirthDate,
		&victimCaseForm2Pg.VictimCaseForm2PhysicalMentalSensoryDifficulties,
		&victimCaseForm2Pg.VictimCaseForm2Nationality,
		&victimCaseForm2Pg.VictimCaseForm2SpecifiedNationality,
		&victimCaseForm2Pg.VictimCaseForm2MigrationCondition,
		&victimCaseForm2Pg.VictimCaseForm2GenderIdentity,
		&victimCaseForm2Pg.VictimCaseForm2SexualOrientation,
		&victimCaseForm2Pg.VictimCaseForm2AssignedSexAtBirth,
		&victimCaseForm2Pg.VictimCaseForm2EthnicAffiliation,
		&victimCaseForm2Pg.VictimCaseForm2IndigenousPeople,
		&victimCaseForm2Pg.VictimCaseForm2CampesinoRecognition,
		&victimCaseForm2Pg.VictimCaseForm2MaritalStatus,
		&victimCaseForm2Pg.VictimCaseForm2LastEducationLevel,
		&victimCaseForm2Pg.VictimCaseForm2Occupation,
		&victimCaseForm2Pg.VictimCaseForm2IncomeGenerationMethod,
		&victimCaseForm2Pg.VictimCaseForm2EmploymentRelationship,
		&victimCaseForm2Pg.VictimCaseForm2ApproxStartAsp,
		&victimCaseForm2Pg.VictimCaseForm2HousingTenancyForm,
		&victimCaseForm2Pg.VictimCaseForm2HousingStratum,
		&victimCaseForm2Pg.VictimCaseForm2CurrentlyPregnant,
		&victimCaseForm2Pg.VictimCaseForm2ResidenceTownCode,
		&victimCaseForm2Pg.VictimCaseForm2ResidenceAddress,
		&victimCaseForm2Pg.VictimCaseForm2ResidenceZone,
		&victimCaseForm2Pg.VictimCaseForm2SupportContactNames,
		&victimCaseForm2Pg.VictimCaseForm2SupportContactPhone,
		&victimCaseForm2Pg.VictimCaseForm2SupportContactEmail,
		&victimCaseForm2Pg.VictimCaseForm2SupportContactKinship,
		&victimCaseForm2Pg.VictimCaseForm2LanguageAssistance,
		&victimCaseForm2Pg.VictimCaseForm2PersonWithDisability,
		&victimCaseForm2Pg.VictimCaseForm2RequireLanguageInterpreter,
		&victimCaseForm2Pg.VictimCaseForm2WorkplaceSectorOccurrence,
		&victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender,
		&victimCaseForm2Pg.VictimCaseForm2AttentionWasAppropriate,
		&victimCaseForm2Pg.VictimCaseForm2AggressorOccupation,
		&victimCaseForm2Pg.VictimCaseForm2SalivaManagementExplanation,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToHear,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTalk,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToSee,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToMove,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTake,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToUnderstand,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToEat,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToInteract,
		&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToDoEveryday,
		&victimCaseForm2Pg.VictimCaseForm2StoppedSeekingHelp,
		&victimCaseForm2Pg.VictimCaseForm2VictimHealthToBlackmail,
		&victimCaseForm2Pg.VictimCaseForm2ThreatenedRevealSexualOrientation,
		&victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender2,
		&victimCaseForm2Pg.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability,
		&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment2,
		&victimCaseForm2Pg.VictimCaseForm2AggressorUsedPositionAuthority,
		&victimCaseForm2Pg.VictimCaseForm2AllowsEasyReport,
		&victimCaseForm2Pg.VictimCaseForm2RiskScore,
		&victimCaseForm2Pg.VictimCaseForm2RiskLevel,
		&victimCaseForm2Pg.VictimCaseForm2VictimCase)

	*victimCaseForm2 = victimCaseForm2Pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	return nil
}

func GetVictimCasesForm2(by common_controllers.By, page int,
	connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2DTO, int, error) {

	var count int
	persistenceCtrl := common_controllers.PersistenceController{}
	victimCaseForm2Path := VictimCaseForm2DBScheme + "." + VictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fields := []string{
		"VictimCaseForm2Id", "VictimCaseForm2ICode", "VictimCaseForm2CreationDate", "VictimCaseForm2UpdateDate", "VictimCaseForm2IdentityName",
		"VictimCaseForm2VictimPhone", "VictimCaseForm2FactsDescription", "VictimCaseForm2FactsDate", "VictimCaseForm2FactsStartTime", "VictimCaseForm2FactsTownCode",
		"VictimCaseForm2FactsZone", "VictimCaseForm2FactsAddress", "VictimCaseForm2ScenarioViolence", "VictimCaseForm2ReportedPreviously", "VictimCaseForm2RecurrenceAggression",
		"VictimCaseForm2NumAgressors", "VictimCaseForm2ProximityPrincipalAggressor", "VictimCaseForm2RelationshipWithPresumedAggressor", "VictimCaseForm2EconomicallyDependent", "VictimCaseForm2AggressorGenderIdentity",
		"VictimCaseForm2AggressorNames", "VictimCaseForm2AggressorDocType", "VictimCaseForm2AggressorDocNumber", "VictimCaseForm2AggressorAddress", "VictimCaseForm2AggressorPhone",
		"VictimCaseForm2AggressorViolencePhysicalIncrease", "VictimCaseForm2AggressorWeaponUsed", "VictimCaseForm2AggressorThreatKill", "VictimCaseForm2AggressorPursuesSpiesDestroys", "VictimCaseForm2AggressorCapableOfKilling",
		"VictimCaseForm2AggressorHasAccessToWeapons", "VictimCaseForm2PartnerUnemployed", "VictimCaseForm2PartnerOtherDenunciations", "VictimCaseForm2AggressorHasPenalBackground", "VictimCaseForm2AggressorForcedSex",
		"VictimCaseForm2AggressorAttemptedStrangulation", "VictimCaseForm2AggressorConsumesDrugs", "VictimCaseForm2AggressorIsAlcoholic", "VictimCaseForm2PartnerControls", "VictimCaseForm2AggressorHadHitInVulnerability",
		"VictimCaseForm2PartnerThreatenedSuicide", "VictimCaseForm2PartnerThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm", "VictimCaseForm2AggressorLimitsContactSupportNetworks", "VictimCaseForm2StillLivesWithAggressor",
		"VictimCaseForm2AggressorViolentlyJealous", "VictimCaseForm2AggressorUnemployed", "VictimCaseForm2AggressorHasPenalBackground2", "VictimCaseForm2AggressorSexuallyHarassment", "VictimCaseForm2AggressorUseDrugs",
		"VictimCaseForm2AggressorIsAlcoholic2", "VictimCaseForm2AggressorControls", "VictimCaseForm2AggressorThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm2", "VictimCaseForm2AggressorCommonSpaces",
		"VictimCaseForm2AggressorHierarchy", "VictimCaseForm2BirthDate", "VictimCaseForm2PhysicalMentalSensoryDifficulties", "VictimCaseForm2Nationality", "VictimCaseForm2SpecifiedNationality",
		"VictimCaseForm2MigrationCondition", "VictimCaseForm2GenderIdentity", "VictimCaseForm2SexualOrientation", "VictimCaseForm2AssignedSexAtBirth", "VictimCaseForm2EthnicAffiliation",
		"VictimCaseForm2IndigenousPeople", "VictimCaseForm2CampesinoRecognition", "VictimCaseForm2MaritalStatus", "VictimCaseForm2LastEducationLevel", "VictimCaseForm2Occupation",
		"VictimCaseForm2IncomeGenerationMethod", "VictimCaseForm2EmploymentRelationship", "VictimCaseForm2ApproxStartAsp", "VictimCaseForm2HousingTenancyForm", "VictimCaseForm2HousingStratum", "VictimCaseForm2CurrentlyPregnant",
		"VictimCaseForm2ResidenceTownCode", "VictimCaseForm2ResidenceAddress", "VictimCaseForm2ResidenceZone", "VictimCaseForm2SupportContactNames", "VictimCaseForm2SupportContactPhone",
		"VictimCaseForm2SupportContactEmail", "VictimCaseForm2SupportContactKinship", "VictimCaseForm2LanguageAssistance",
		"VictimCaseForm2PersonWithDisability", "VictimCaseForm2RequireLanguageInterpreter", "VictimCaseForm2WorkplaceSectorOccurrence", "VictimCaseForm2ViolenceMotivatedByGender", "VictimCaseForm2AttentionWasAppropriate", "VictimCaseForm2AggressorOccupation", "VictimCaseForm2SalivaManagementExplanation",
		"VictimCaseForm2ActivitiesUnableToHear", "VictimCaseForm2ActivitiesUnableToTalk", "VictimCaseForm2ActivitiesUnableToSee", "VictimCaseForm2ActivitiesUnableToMove", "VictimCaseForm2ActivitiesUnableToTake",
		"VictimCaseForm2ActivitiesUnableToUnderstand", "VictimCaseForm2ActivitiesUnableToEat", "VictimCaseForm2ActivitiesUnableToInteract", "VictimCaseForm2ActivitiesUnableToDoEveryday",
		"VictimCaseForm2StoppedSeekingHelp", "VictimCaseForm2VictimHealthToBlackmail", "VictimCaseForm2ThreatenedRevealSexualOrientation", "VictimCaseForm2ViolenceMotivatedByGender2", "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", "VictimCaseForm2AggressorSexuallyHarassment2", "VictimCaseForm2AggressorUsedPositionAuthority",
		"VictimCaseForm2AllowsEasyReport", "VictimCaseForm2RiskScore", "VictimCaseForm2RiskLevel", "VictimCaseForm2VictimCase",
	}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fields, []string{}, VictimCaseForm2DBName,
		[]string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2DBScheme,
		VictimCaseForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, true) +
		` ORDER BY ` + victimCaseForm2Path + `.` + VictimCaseForm2FieldDefinitions["VictimCaseForm2CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	var victimCasesForm2 []VictimCaseForm2DTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2PgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2Id,
			&victimCaseForm2Pg.VictimCaseForm2ICode,
			&victimCaseForm2Pg.VictimCaseForm2CreationDate,
			&victimCaseForm2Pg.VictimCaseForm2UpdateDate,
			&victimCaseForm2Pg.VictimCaseForm2IdentityName,
			&victimCaseForm2Pg.VictimCaseForm2VictimPhone,
			&victimCaseForm2Pg.VictimCaseForm2FactsDescription,
			&victimCaseForm2Pg.VictimCaseForm2FactsDate,
			&victimCaseForm2Pg.VictimCaseForm2FactsStartTime,
			&victimCaseForm2Pg.VictimCaseForm2FactsTownCode,
			&victimCaseForm2Pg.VictimCaseForm2FactsZone,
			&victimCaseForm2Pg.VictimCaseForm2FactsAddress,
			&victimCaseForm2Pg.VictimCaseForm2ScenarioViolence,
			&victimCaseForm2Pg.VictimCaseForm2ReportedPreviously,
			&victimCaseForm2Pg.VictimCaseForm2RecurrenceAggression,
			&victimCaseForm2Pg.VictimCaseForm2NumAgressors,
			&victimCaseForm2Pg.VictimCaseForm2ProximityPrincipalAggressor,
			&victimCaseForm2Pg.VictimCaseForm2RelationshipWithPresumedAggressor,
			&victimCaseForm2Pg.VictimCaseForm2EconomicallyDependent,
			&victimCaseForm2Pg.VictimCaseForm2AggressorGenderIdentity,
			&victimCaseForm2Pg.VictimCaseForm2AggressorNames,
			&victimCaseForm2Pg.VictimCaseForm2AggressorDocType,
			&victimCaseForm2Pg.VictimCaseForm2AggressorDocNumber,
			&victimCaseForm2Pg.VictimCaseForm2AggressorAddress,
			&victimCaseForm2Pg.VictimCaseForm2AggressorPhone,
			&victimCaseForm2Pg.VictimCaseForm2AggressorViolencePhysicalIncrease,
			&victimCaseForm2Pg.VictimCaseForm2AggressorWeaponUsed,
			&victimCaseForm2Pg.VictimCaseForm2AggressorThreatKill,
			&victimCaseForm2Pg.VictimCaseForm2AggressorPursuesSpiesDestroys,
			&victimCaseForm2Pg.VictimCaseForm2AggressorCapableOfKilling,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHasAccessToWeapons,
			&victimCaseForm2Pg.VictimCaseForm2PartnerUnemployed,
			&victimCaseForm2Pg.VictimCaseForm2PartnerOtherDenunciations,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground,
			&victimCaseForm2Pg.VictimCaseForm2AggressorForcedSex,
			&victimCaseForm2Pg.VictimCaseForm2AggressorAttemptedStrangulation,
			&victimCaseForm2Pg.VictimCaseForm2AggressorConsumesDrugs,
			&victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic,
			&victimCaseForm2Pg.VictimCaseForm2PartnerControls,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHadHitInVulnerability,
			&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedSuicide,
			&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedDamageMembers,
			&victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm,
			&victimCaseForm2Pg.VictimCaseForm2AggressorLimitsContactSupportNetworks,
			&victimCaseForm2Pg.VictimCaseForm2StillLivesWithAggressor,
			&victimCaseForm2Pg.VictimCaseForm2AggressorViolentlyJealous,
			&victimCaseForm2Pg.VictimCaseForm2AggressorUnemployed,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment,
			&victimCaseForm2Pg.VictimCaseForm2AggressorUseDrugs,
			&victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorControls,
			&victimCaseForm2Pg.VictimCaseForm2AggressorThreatenedDamageMembers,
			&victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorCommonSpaces,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHierarchy,
			&victimCaseForm2Pg.VictimCaseForm2BirthDate,
			&victimCaseForm2Pg.VictimCaseForm2PhysicalMentalSensoryDifficulties,
			&victimCaseForm2Pg.VictimCaseForm2Nationality,
			&victimCaseForm2Pg.VictimCaseForm2SpecifiedNationality,
			&victimCaseForm2Pg.VictimCaseForm2MigrationCondition,
			&victimCaseForm2Pg.VictimCaseForm2GenderIdentity,
			&victimCaseForm2Pg.VictimCaseForm2SexualOrientation,
			&victimCaseForm2Pg.VictimCaseForm2AssignedSexAtBirth,
			&victimCaseForm2Pg.VictimCaseForm2EthnicAffiliation,
			&victimCaseForm2Pg.VictimCaseForm2IndigenousPeople,
			&victimCaseForm2Pg.VictimCaseForm2CampesinoRecognition,
			&victimCaseForm2Pg.VictimCaseForm2MaritalStatus,
			&victimCaseForm2Pg.VictimCaseForm2LastEducationLevel,
			&victimCaseForm2Pg.VictimCaseForm2Occupation,
			&victimCaseForm2Pg.VictimCaseForm2IncomeGenerationMethod,
			&victimCaseForm2Pg.VictimCaseForm2EmploymentRelationship,
			&victimCaseForm2Pg.VictimCaseForm2ApproxStartAsp,
			&victimCaseForm2Pg.VictimCaseForm2HousingTenancyForm,
			&victimCaseForm2Pg.VictimCaseForm2HousingStratum,
			&victimCaseForm2Pg.VictimCaseForm2CurrentlyPregnant,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceTownCode,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceAddress,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceZone,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactNames,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactPhone,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactEmail,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactKinship,
			&victimCaseForm2Pg.VictimCaseForm2LanguageAssistance,
			&victimCaseForm2Pg.VictimCaseForm2PersonWithDisability,
			&victimCaseForm2Pg.VictimCaseForm2RequireLanguageInterpreter,
			&victimCaseForm2Pg.VictimCaseForm2WorkplaceSectorOccurrence,
			&victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender,
			&victimCaseForm2Pg.VictimCaseForm2AttentionWasAppropriate,
			&victimCaseForm2Pg.VictimCaseForm2AggressorOccupation,
			&victimCaseForm2Pg.VictimCaseForm2SalivaManagementExplanation,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToHear,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTalk,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToSee,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToMove,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTake,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToUnderstand,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToEat,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToInteract,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToDoEveryday,
			&victimCaseForm2Pg.VictimCaseForm2StoppedSeekingHelp,
			&victimCaseForm2Pg.VictimCaseForm2VictimHealthToBlackmail,
			&victimCaseForm2Pg.VictimCaseForm2ThreatenedRevealSexualOrientation,
			&victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability,
			&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorUsedPositionAuthority,
			&victimCaseForm2Pg.VictimCaseForm2AllowsEasyReport,
			&victimCaseForm2Pg.VictimCaseForm2RiskScore,
			&victimCaseForm2Pg.VictimCaseForm2RiskLevel,
			&victimCaseForm2Pg.VictimCaseForm2VictimCase)

		victimCasesForm2 = append(victimCasesForm2, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Count total rows for first page
	if page == 0 {
		countQuery := `SELECT COUNT(*) FROM ` + victimCaseForm2Path +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, by.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return victimCasesForm2, count, nil
}

func GetAllVictimCasesForm2(victimCaseForm2Status string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2DTO, int, error) {

	var count int
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2Path string = VictimCaseForm2DBScheme + "." + VictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fields := []string{
		"VictimCaseForm2Id", "VictimCaseForm2ICode", "VictimCaseForm2CreationDate", "VictimCaseForm2UpdateDate", "VictimCaseForm2IdentityName",
		"VictimCaseForm2VictimPhone", "VictimCaseForm2FactsDescription", "VictimCaseForm2FactsDate", "VictimCaseForm2FactsStartTime", "VictimCaseForm2FactsTownCode",
		"VictimCaseForm2FactsZone", "VictimCaseForm2FactsAddress", "VictimCaseForm2ScenarioViolence", "VictimCaseForm2ReportedPreviously", "VictimCaseForm2RecurrenceAggression",
		"VictimCaseForm2NumAgressors", "VictimCaseForm2ProximityPrincipalAggressor", "VictimCaseForm2RelationshipWithPresumedAggressor", "VictimCaseForm2EconomicallyDependent", "VictimCaseForm2AggressorGenderIdentity",
		"VictimCaseForm2AggressorNames", "VictimCaseForm2AggressorDocType", "VictimCaseForm2AggressorDocNumber", "VictimCaseForm2AggressorAddress", "VictimCaseForm2AggressorPhone",
		"VictimCaseForm2AggressorViolencePhysicalIncrease", "VictimCaseForm2AggressorWeaponUsed", "VictimCaseForm2AggressorThreatKill", "VictimCaseForm2AggressorPursuesSpiesDestroys", "VictimCaseForm2AggressorCapableOfKilling",
		"VictimCaseForm2AggressorHasAccessToWeapons", "VictimCaseForm2PartnerUnemployed", "VictimCaseForm2PartnerOtherDenunciations", "VictimCaseForm2AggressorHasPenalBackground", "VictimCaseForm2AggressorForcedSex",
		"VictimCaseForm2AggressorAttemptedStrangulation", "VictimCaseForm2AggressorConsumesDrugs", "VictimCaseForm2AggressorIsAlcoholic", "VictimCaseForm2PartnerControls", "VictimCaseForm2AggressorHadHitInVulnerability",
		"VictimCaseForm2PartnerThreatenedSuicide", "VictimCaseForm2PartnerThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm", "VictimCaseForm2AggressorLimitsContactSupportNetworks", "VictimCaseForm2StillLivesWithAggressor",
		"VictimCaseForm2AggressorViolentlyJealous", "VictimCaseForm2AggressorUnemployed", "VictimCaseForm2AggressorHasPenalBackground2", "VictimCaseForm2AggressorSexuallyHarassment", "VictimCaseForm2AggressorUseDrugs",
		"VictimCaseForm2AggressorIsAlcoholic2", "VictimCaseForm2AggressorControls", "VictimCaseForm2AggressorThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm2", "VictimCaseForm2AggressorCommonSpaces",
		"VictimCaseForm2AggressorHierarchy", "VictimCaseForm2BirthDate", "VictimCaseForm2PhysicalMentalSensoryDifficulties", "VictimCaseForm2Nationality", "VictimCaseForm2SpecifiedNationality",
		"VictimCaseForm2MigrationCondition", "VictimCaseForm2GenderIdentity", "VictimCaseForm2SexualOrientation", "VictimCaseForm2AssignedSexAtBirth", "VictimCaseForm2EthnicAffiliation",
		"VictimCaseForm2IndigenousPeople", "VictimCaseForm2CampesinoRecognition", "VictimCaseForm2MaritalStatus", "VictimCaseForm2LastEducationLevel", "VictimCaseForm2Occupation",
		"VictimCaseForm2IncomeGenerationMethod", "VictimCaseForm2EmploymentRelationship", "VictimCaseForm2ApproxStartAsp", "VictimCaseForm2HousingTenancyForm", "VictimCaseForm2HousingStratum", "VictimCaseForm2CurrentlyPregnant",
		"VictimCaseForm2ResidenceTownCode", "VictimCaseForm2ResidenceAddress", "VictimCaseForm2ResidenceZone", "VictimCaseForm2SupportContactNames", "VictimCaseForm2SupportContactPhone",
		"VictimCaseForm2SupportContactEmail", "VictimCaseForm2SupportContactKinship", "VictimCaseForm2LanguageAssistance",
		"VictimCaseForm2PersonWithDisability", "VictimCaseForm2RequireLanguageInterpreter", "VictimCaseForm2WorkplaceSectorOccurrence", "VictimCaseForm2ViolenceMotivatedByGender", "VictimCaseForm2AttentionWasAppropriate", "VictimCaseForm2AggressorOccupation",
		"VictimCaseForm2SalivaManagementExplanation",
		"VictimCaseForm2ActivitiesUnableToHear", "VictimCaseForm2ActivitiesUnableToTalk", "VictimCaseForm2ActivitiesUnableToSee", "VictimCaseForm2ActivitiesUnableToMove", "VictimCaseForm2ActivitiesUnableToTake",
		"VictimCaseForm2ActivitiesUnableToUnderstand", "VictimCaseForm2ActivitiesUnableToEat", "VictimCaseForm2ActivitiesUnableToInteract", "VictimCaseForm2ActivitiesUnableToDoEveryday",
		"VictimCaseForm2StoppedSeekingHelp", "VictimCaseForm2VictimHealthToBlackmail", "VictimCaseForm2ThreatenedRevealSexualOrientation", "VictimCaseForm2ViolenceMotivatedByGender2", "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", "VictimCaseForm2AggressorSexuallyHarassment2", "VictimCaseForm2AggressorUsedPositionAuthority",
		"VictimCaseForm2AllowsEasyReport", "VictimCaseForm2RiskScore", "VictimCaseForm2RiskLevel", "VictimCaseForm2VictimCase",
	}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fields, []string{}, VictimCaseForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2Path +
		` WHERE ` + victimCaseForm2Path + `.` + VictimCaseForm2FieldDefinitions["VictimCaseForm2Status"].DBName + ` = $1 ` +
		` ORDER BY ` + victimCaseForm2Path + `.` + VictimCaseForm2FieldDefinitions["VictimCaseForm2CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, victimCaseForm2Status)
	var results []VictimCaseForm2DTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2PgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2Id,
			&victimCaseForm2Pg.VictimCaseForm2ICode,
			&victimCaseForm2Pg.VictimCaseForm2CreationDate,
			&victimCaseForm2Pg.VictimCaseForm2UpdateDate,
			&victimCaseForm2Pg.VictimCaseForm2IdentityName,
			&victimCaseForm2Pg.VictimCaseForm2VictimPhone,
			&victimCaseForm2Pg.VictimCaseForm2FactsDescription,
			&victimCaseForm2Pg.VictimCaseForm2FactsDate,
			&victimCaseForm2Pg.VictimCaseForm2FactsStartTime,
			&victimCaseForm2Pg.VictimCaseForm2FactsTownCode,
			&victimCaseForm2Pg.VictimCaseForm2FactsZone,
			&victimCaseForm2Pg.VictimCaseForm2FactsAddress,
			&victimCaseForm2Pg.VictimCaseForm2ScenarioViolence,
			&victimCaseForm2Pg.VictimCaseForm2ReportedPreviously,
			&victimCaseForm2Pg.VictimCaseForm2RecurrenceAggression,
			&victimCaseForm2Pg.VictimCaseForm2NumAgressors,
			&victimCaseForm2Pg.VictimCaseForm2ProximityPrincipalAggressor,
			&victimCaseForm2Pg.VictimCaseForm2RelationshipWithPresumedAggressor,
			&victimCaseForm2Pg.VictimCaseForm2EconomicallyDependent,
			&victimCaseForm2Pg.VictimCaseForm2AggressorGenderIdentity,
			&victimCaseForm2Pg.VictimCaseForm2AggressorNames,
			&victimCaseForm2Pg.VictimCaseForm2AggressorDocType,
			&victimCaseForm2Pg.VictimCaseForm2AggressorDocNumber,
			&victimCaseForm2Pg.VictimCaseForm2AggressorAddress,
			&victimCaseForm2Pg.VictimCaseForm2AggressorPhone,
			&victimCaseForm2Pg.VictimCaseForm2AggressorViolencePhysicalIncrease,
			&victimCaseForm2Pg.VictimCaseForm2AggressorWeaponUsed,
			&victimCaseForm2Pg.VictimCaseForm2AggressorThreatKill,
			&victimCaseForm2Pg.VictimCaseForm2AggressorPursuesSpiesDestroys,
			&victimCaseForm2Pg.VictimCaseForm2AggressorCapableOfKilling,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHasAccessToWeapons,
			&victimCaseForm2Pg.VictimCaseForm2PartnerUnemployed,
			&victimCaseForm2Pg.VictimCaseForm2PartnerOtherDenunciations,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground,
			&victimCaseForm2Pg.VictimCaseForm2AggressorForcedSex,
			&victimCaseForm2Pg.VictimCaseForm2AggressorAttemptedStrangulation,
			&victimCaseForm2Pg.VictimCaseForm2AggressorConsumesDrugs,
			&victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic,
			&victimCaseForm2Pg.VictimCaseForm2PartnerControls,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHadHitInVulnerability,
			&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedSuicide,
			&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedDamageMembers,
			&victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm,
			&victimCaseForm2Pg.VictimCaseForm2AggressorLimitsContactSupportNetworks,
			&victimCaseForm2Pg.VictimCaseForm2StillLivesWithAggressor,
			&victimCaseForm2Pg.VictimCaseForm2AggressorViolentlyJealous,
			&victimCaseForm2Pg.VictimCaseForm2AggressorUnemployed,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment,
			&victimCaseForm2Pg.VictimCaseForm2AggressorUseDrugs,
			&victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorControls,
			&victimCaseForm2Pg.VictimCaseForm2AggressorThreatenedDamageMembers,
			&victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorCommonSpaces,
			&victimCaseForm2Pg.VictimCaseForm2AggressorHierarchy,
			&victimCaseForm2Pg.VictimCaseForm2BirthDate,
			&victimCaseForm2Pg.VictimCaseForm2PhysicalMentalSensoryDifficulties,
			&victimCaseForm2Pg.VictimCaseForm2Nationality,
			&victimCaseForm2Pg.VictimCaseForm2SpecifiedNationality,
			&victimCaseForm2Pg.VictimCaseForm2MigrationCondition,
			&victimCaseForm2Pg.VictimCaseForm2GenderIdentity,
			&victimCaseForm2Pg.VictimCaseForm2SexualOrientation,
			&victimCaseForm2Pg.VictimCaseForm2AssignedSexAtBirth,
			&victimCaseForm2Pg.VictimCaseForm2EthnicAffiliation,
			&victimCaseForm2Pg.VictimCaseForm2IndigenousPeople,
			&victimCaseForm2Pg.VictimCaseForm2CampesinoRecognition,
			&victimCaseForm2Pg.VictimCaseForm2MaritalStatus,
			&victimCaseForm2Pg.VictimCaseForm2LastEducationLevel,
			&victimCaseForm2Pg.VictimCaseForm2Occupation,
			&victimCaseForm2Pg.VictimCaseForm2IncomeGenerationMethod,
			&victimCaseForm2Pg.VictimCaseForm2EmploymentRelationship,
			&victimCaseForm2Pg.VictimCaseForm2ApproxStartAsp,
			&victimCaseForm2Pg.VictimCaseForm2HousingTenancyForm,
			&victimCaseForm2Pg.VictimCaseForm2HousingStratum,
			&victimCaseForm2Pg.VictimCaseForm2CurrentlyPregnant,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceTownCode,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceAddress,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceZone,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactNames,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactPhone,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactEmail,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactKinship,
			&victimCaseForm2Pg.VictimCaseForm2LanguageAssistance,
			&victimCaseForm2Pg.VictimCaseForm2PersonWithDisability,
			&victimCaseForm2Pg.VictimCaseForm2RequireLanguageInterpreter,
			&victimCaseForm2Pg.VictimCaseForm2WorkplaceSectorOccurrence,
			&victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender,
			&victimCaseForm2Pg.VictimCaseForm2AttentionWasAppropriate,
			&victimCaseForm2Pg.VictimCaseForm2AggressorOccupation,
			&victimCaseForm2Pg.VictimCaseForm2SalivaManagementExplanation,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToHear,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTalk,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToSee,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToMove,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTake,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToUnderstand,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToEat,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToInteract,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToDoEveryday,
			&victimCaseForm2Pg.VictimCaseForm2StoppedSeekingHelp,
			&victimCaseForm2Pg.VictimCaseForm2VictimHealthToBlackmail,
			&victimCaseForm2Pg.VictimCaseForm2ThreatenedRevealSexualOrientation,
			&victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability,
			&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorUsedPositionAuthority,
			&victimCaseForm2Pg.VictimCaseForm2AllowsEasyReport,
			&victimCaseForm2Pg.VictimCaseForm2RiskScore,
			&victimCaseForm2Pg.VictimCaseForm2RiskLevel,
			&victimCaseForm2Pg.VictimCaseForm2VictimCase)

		results = append(results, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Count total rows for first page
	if page == 0 {
		countQuery := `SELECT COUNT(*) FROM ` + victimCaseForm2Path +
			` WHERE ` + victimCaseForm2Path + `.` + VictimCaseForm2FieldDefinitions["VictimCaseForm2Status"].DBName + ` = $1 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, victimCaseForm2Status)
		persistenceCtrl.Scan(&count)
	}

	return results, count, nil
}

func UpdateVictimCaseForm2(victimCaseForm2 *VictimCaseForm2DTO,
	connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fields := []string{
		"VictimCaseForm2UpdateDate", "VictimCaseForm2IdentityName",
		"VictimCaseForm2VictimPhone", "VictimCaseForm2FactsDescription", "VictimCaseForm2FactsDate", "VictimCaseForm2FactsStartTime", "VictimCaseForm2FactsTownCode",
		"VictimCaseForm2FactsZone", "VictimCaseForm2FactsAddress", "VictimCaseForm2ScenarioViolence", "VictimCaseForm2ReportedPreviously", "VictimCaseForm2RecurrenceAggression",
		"VictimCaseForm2NumAgressors", "VictimCaseForm2ProximityPrincipalAggressor", "VictimCaseForm2RelationshipWithPresumedAggressor", "VictimCaseForm2EconomicallyDependent", "VictimCaseForm2AggressorGenderIdentity",
		"VictimCaseForm2AggressorNames", "VictimCaseForm2AggressorDocType", "VictimCaseForm2AggressorDocNumber", "VictimCaseForm2AggressorAddress", "VictimCaseForm2AggressorPhone",
		"VictimCaseForm2AggressorViolencePhysicalIncrease", "VictimCaseForm2AggressorWeaponUsed", "VictimCaseForm2AggressorThreatKill", "VictimCaseForm2AggressorPursuesSpiesDestroys", "VictimCaseForm2AggressorCapableOfKilling",
		"VictimCaseForm2AggressorHasAccessToWeapons", "VictimCaseForm2PartnerUnemployed", "VictimCaseForm2PartnerOtherDenunciations", "VictimCaseForm2AggressorHasPenalBackground", "VictimCaseForm2AggressorForcedSex",
		"VictimCaseForm2AggressorAttemptedStrangulation", "VictimCaseForm2AggressorConsumesDrugs", "VictimCaseForm2AggressorIsAlcoholic", "VictimCaseForm2PartnerControls", "VictimCaseForm2AggressorHadHitInVulnerability",
		"VictimCaseForm2PartnerThreatenedSuicide", "VictimCaseForm2PartnerThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm", "VictimCaseForm2AggressorLimitsContactSupportNetworks", "VictimCaseForm2StillLivesWithAggressor",
		"VictimCaseForm2AggressorViolentlyJealous", "VictimCaseForm2AggressorUnemployed", "VictimCaseForm2AggressorHasPenalBackground2", "VictimCaseForm2AggressorSexuallyHarassment", "VictimCaseForm2AggressorUseDrugs",
		"VictimCaseForm2AggressorIsAlcoholic2", "VictimCaseForm2AggressorControls", "VictimCaseForm2AggressorThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm2", "VictimCaseForm2AggressorCommonSpaces",
		"VictimCaseForm2AggressorHierarchy", "VictimCaseForm2BirthDate", "VictimCaseForm2PhysicalMentalSensoryDifficulties", "VictimCaseForm2Nationality", "VictimCaseForm2SpecifiedNationality",
		"VictimCaseForm2MigrationCondition", "VictimCaseForm2GenderIdentity", "VictimCaseForm2SexualOrientation", "VictimCaseForm2AssignedSexAtBirth", "VictimCaseForm2EthnicAffiliation",
		"VictimCaseForm2IndigenousPeople", "VictimCaseForm2CampesinoRecognition", "VictimCaseForm2MaritalStatus", "VictimCaseForm2LastEducationLevel", "VictimCaseForm2Occupation",
		"VictimCaseForm2IncomeGenerationMethod", "VictimCaseForm2EmploymentRelationship", "VictimCaseForm2ApproxStartAsp", "VictimCaseForm2HousingTenancyForm", "VictimCaseForm2HousingStratum", "VictimCaseForm2CurrentlyPregnant",
		"VictimCaseForm2ResidenceTownCode", "VictimCaseForm2ResidenceAddress", "VictimCaseForm2ResidenceZone", "VictimCaseForm2SupportContactNames", "VictimCaseForm2SupportContactPhone",
		"VictimCaseForm2SupportContactEmail", "VictimCaseForm2SupportContactKinship", "VictimCaseForm2LanguageAssistance",
		"VictimCaseForm2PersonWithDisability", "VictimCaseForm2RequireLanguageInterpreter", "VictimCaseForm2WorkplaceSectorOccurrence", "VictimCaseForm2ViolenceMotivatedByGender", "VictimCaseForm2AttentionWasAppropriate", "VictimCaseForm2AggressorOccupation",
		"VictimCaseForm2SalivaManagementExplanation",
		"VictimCaseForm2ActivitiesUnableToHear", "VictimCaseForm2ActivitiesUnableToTalk", "VictimCaseForm2ActivitiesUnableToSee", "VictimCaseForm2ActivitiesUnableToMove", "VictimCaseForm2ActivitiesUnableToTake",
		"VictimCaseForm2ActivitiesUnableToUnderstand", "VictimCaseForm2ActivitiesUnableToEat", "VictimCaseForm2ActivitiesUnableToInteract", "VictimCaseForm2ActivitiesUnableToDoEveryday",
		"VictimCaseForm2StoppedSeekingHelp", "VictimCaseForm2VictimHealthToBlackmail", "VictimCaseForm2ThreatenedRevealSexualOrientation", "VictimCaseForm2ViolenceMotivatedByGender2", "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", "VictimCaseForm2AggressorSexuallyHarassment2", "VictimCaseForm2AggressorUsedPositionAuthority",
		"VictimCaseForm2AllowsEasyReport", "VictimCaseForm2RiskScore", "VictimCaseForm2RiskLevel",
	}

	query := common_dao.GetSQL(common_dao.SQL_UPDATE, fields, []string{}, VictimCaseForm2DBName, []string{"VictimCaseForm2Id"}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query,
		victimCaseForm2.VictimCaseForm2Id,

		victimCaseForm2.VictimCaseForm2UpdateDate,
		victimCaseForm2.VictimCaseForm2IdentityName,
		victimCaseForm2.VictimCaseForm2VictimPhone,
		victimCaseForm2.VictimCaseForm2FactsDescription,
		victimCaseForm2.VictimCaseForm2FactsDate,
		victimCaseForm2.VictimCaseForm2FactsStartTime,
		victimCaseForm2.VictimCaseForm2FactsTownCode,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2FactsZone.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2FactsAddress,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ScenarioViolence.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ReportedPreviously.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2RecurrenceAggression.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2NumAgressors.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ProximityPrincipalAggressor.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2EconomicallyDependent.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorGenderIdentity.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2AggressorNames,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorDocType.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2AggressorDocNumber,
		victimCaseForm2.VictimCaseForm2AggressorAddress,
		victimCaseForm2.VictimCaseForm2AggressorPhone,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorWeaponUsed.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorThreatKill.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorCapableOfKilling.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerUnemployed.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerOtherDenunciations.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHasPenalBackground.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorForcedSex.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorAttemptedStrangulation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorConsumesDrugs.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorIsAlcoholic.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerControls.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHadHitInVulnerability.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerThreatenedSuicide.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PartnerThreatenedDamageMembers.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorLimitsContactSupportNetworks.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2StillLivesWithAggressor.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorViolentlyJealous.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorUnemployed.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHasPenalBackground2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorUseDrugs.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorIsAlcoholic2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorControls.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorThreatenedDamageMembers.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorCommonSpaces.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorHierarchy.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2BirthDate,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PhysicalMentalSensoryDifficulties.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2Nationality.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2SpecifiedNationality.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2MigrationCondition.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2GenderIdentity.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2SexualOrientation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AssignedSexAtBirth.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2EthnicAffiliation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2IndigenousPeople.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2CampesinoRecognition.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2MaritalStatus.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2LastEducationLevel.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2Occupation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2EmploymentRelationship.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2ApproxStartAsp,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2HousingTenancyForm.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2HousingStratum.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2CurrentlyPregnant.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2ResidenceTownCode,
		victimCaseForm2.VictimCaseForm2ResidenceAddress,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ResidenceZone.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2SupportContactNames,
		victimCaseForm2.VictimCaseForm2SupportContactPhone,
		victimCaseForm2.VictimCaseForm2SupportContactEmail,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2SupportContactKinship.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2LanguageAssistance,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2PersonWithDisability.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2RequireLanguageInterpreter.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2WorkplaceSectorOccurrence.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AttentionWasAppropriate.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorOccupation.VictimCaseForm2EnumsId),
		victimCaseForm2.VictimCaseForm2SalivaManagementExplanation,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToHear,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToTalk,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToSee,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToMove,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToTake,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToUnderstand,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToEat,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToInteract,
		victimCaseForm2.VictimCaseForm2ActivitiesUnableToDoEveryday,
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2StoppedSeekingHelp.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2VictimHealthToBlackmail.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ThreatenedRevealSexualOrientation.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment2.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AggressorUsedPositionAuthority.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimCaseForm2.VictimCaseForm2AllowsEasyReport.VictimCaseForm2EnumsId),
		&victimCaseForm2.VictimCaseForm2RiskScore,
		&victimCaseForm2.VictimCaseForm2RiskLevel,
	)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// Algunas utilidades
func SetVictimCaseForm2Defaults(victimCaseForm2 *VictimCaseForm2DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		victimCaseForm2.VictimCaseForm2CreationDate = time.Now()
		victimCaseForm2.VictimCaseForm2UpdateDate = time.Now()
		victimCaseForm2.VictimCaseForm2ICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		victimCaseForm2.VictimCaseForm2UpdateDate = time.Now()
	}
}

func (obj *VictimCaseForm2PgDB) ToDTO() VictimCaseForm2DTO {
	var dto VictimCaseForm2DTO

	if obj.VictimCaseForm2Id.Valid {
		dto.VictimCaseForm2Id = uint64(obj.VictimCaseForm2Id.Int64)
	}
	if obj.VictimCaseForm2ICode.Valid {
		dto.VictimCaseForm2ICode = obj.VictimCaseForm2ICode.String
	}
	if obj.VictimCaseForm2CreationDate.Valid {
		dto.VictimCaseForm2CreationDate = obj.VictimCaseForm2CreationDate.Time
	}
	if obj.VictimCaseForm2UpdateDate.Valid {
		dto.VictimCaseForm2UpdateDate = obj.VictimCaseForm2UpdateDate.Time
	}
	if obj.VictimCaseForm2IdentityName.Valid {
		dto.VictimCaseForm2IdentityName = obj.VictimCaseForm2IdentityName.String
	}
	if obj.VictimCaseForm2VictimPhone.Valid {
		dto.VictimCaseForm2VictimPhone = uint64(obj.VictimCaseForm2VictimPhone.Int64)
	}
	if obj.VictimCaseForm2FactsDescription.Valid {
		dto.VictimCaseForm2FactsDescription = obj.VictimCaseForm2FactsDescription.String
	}
	if obj.VictimCaseForm2FactsDate.Valid {
		dto.VictimCaseForm2FactsDate = obj.VictimCaseForm2FactsDate.Time
	}
	if obj.VictimCaseForm2FactsStartTime.Valid {
		dto.VictimCaseForm2FactsStartTime, _ = time.Parse(common_config.DateTime.TIME_WITH_MILLISECONDS_FORMAT, obj.VictimCaseForm2FactsStartTime.String)
	}
	if obj.VictimCaseForm2FactsTownCode.Valid {
		dto.VictimCaseForm2FactsTownCode = obj.VictimCaseForm2FactsTownCode.String
	}
	if obj.VictimCaseForm2FactsAddress.Valid {
		dto.VictimCaseForm2FactsAddress = obj.VictimCaseForm2FactsAddress.String
	}

	if obj.VictimCaseForm2AggressorNames.Valid {
		dto.VictimCaseForm2AggressorNames = obj.VictimCaseForm2AggressorNames.String
	}
	if obj.VictimCaseForm2AggressorDocNumber.Valid {
		dto.VictimCaseForm2AggressorDocNumber = obj.VictimCaseForm2AggressorDocNumber.String
	}
	if obj.VictimCaseForm2AggressorAddress.Valid {
		dto.VictimCaseForm2AggressorAddress = obj.VictimCaseForm2AggressorAddress.String
	}
	if obj.VictimCaseForm2AggressorPhone.Valid {
		dto.VictimCaseForm2AggressorPhone = uint64(obj.VictimCaseForm2AggressorPhone.Int64)
	}

	//Campos de calificación

	if obj.VictimCaseForm2ActivitiesUnableToHear.Valid {
		dto.VictimCaseForm2ActivitiesUnableToHear = obj.VictimCaseForm2ActivitiesUnableToHear.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToTalk.Valid {
		dto.VictimCaseForm2ActivitiesUnableToTalk = obj.VictimCaseForm2ActivitiesUnableToTalk.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToSee.Valid {
		dto.VictimCaseForm2ActivitiesUnableToSee = obj.VictimCaseForm2ActivitiesUnableToSee.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToMove.Valid {
		dto.VictimCaseForm2ActivitiesUnableToMove = obj.VictimCaseForm2ActivitiesUnableToMove.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToTake.Valid {
		dto.VictimCaseForm2ActivitiesUnableToTake = obj.VictimCaseForm2ActivitiesUnableToTake.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToUnderstand.Valid {
		dto.VictimCaseForm2ActivitiesUnableToUnderstand = obj.VictimCaseForm2ActivitiesUnableToUnderstand.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToEat.Valid {
		dto.VictimCaseForm2ActivitiesUnableToEat = obj.VictimCaseForm2ActivitiesUnableToEat.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToInteract.Valid {
		dto.VictimCaseForm2ActivitiesUnableToInteract = obj.VictimCaseForm2ActivitiesUnableToInteract.Int64
	}

	if obj.VictimCaseForm2ActivitiesUnableToDoEveryday.Valid {
		dto.VictimCaseForm2ActivitiesUnableToDoEveryday = obj.VictimCaseForm2ActivitiesUnableToDoEveryday.Int64
	}

	// Enums – convert to DTO struct
	if obj.VictimCaseForm2AggressorDocType.Valid {
		dto.VictimCaseForm2AggressorDocType = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorDocType.Int64)}
	}
	if obj.VictimCaseForm2FactsZone.Valid {
		dto.VictimCaseForm2FactsZone = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2FactsZone.Int64)}
	}
	if obj.VictimCaseForm2ScenarioViolence.Valid {
		dto.VictimCaseForm2ScenarioViolence = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ScenarioViolence.Int64)}
	}
	if obj.VictimCaseForm2ReportedPreviously.Valid {
		dto.VictimCaseForm2ReportedPreviously = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ReportedPreviously.Int64)}
	}
	if obj.VictimCaseForm2RecurrenceAggression.Valid {
		dto.VictimCaseForm2RecurrenceAggression = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2RecurrenceAggression.Int64)}
	}
	if obj.VictimCaseForm2NumAgressors.Valid {
		dto.VictimCaseForm2NumAgressors = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2NumAgressors.Int64)}
	}
	if obj.VictimCaseForm2ProximityPrincipalAggressor.Valid {
		dto.VictimCaseForm2ProximityPrincipalAggressor = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ProximityPrincipalAggressor.Int64)}
	}
	if obj.VictimCaseForm2RelationshipWithPresumedAggressor.Valid {
		dto.VictimCaseForm2RelationshipWithPresumedAggressor = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2RelationshipWithPresumedAggressor.Int64)}
	}
	if obj.VictimCaseForm2EconomicallyDependent.Valid {
		dto.VictimCaseForm2EconomicallyDependent = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2EconomicallyDependent.Int64)}
	}
	if obj.VictimCaseForm2AggressorGenderIdentity.Valid {
		dto.VictimCaseForm2AggressorGenderIdentity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorGenderIdentity.Int64)}
	}
	if obj.VictimCaseForm2AggressorViolencePhysicalIncrease.Valid {
		dto.VictimCaseForm2AggressorViolencePhysicalIncrease = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorViolencePhysicalIncrease.Int64)}
	}
	if obj.VictimCaseForm2AggressorWeaponUsed.Valid {
		dto.VictimCaseForm2AggressorWeaponUsed = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorWeaponUsed.Int64)}
	}
	if obj.VictimCaseForm2AggressorThreatKill.Valid {
		dto.VictimCaseForm2AggressorThreatKill = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorThreatKill.Int64)}
	}
	if obj.VictimCaseForm2AggressorPursuesSpiesDestroys.Valid {
		dto.VictimCaseForm2AggressorPursuesSpiesDestroys = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorPursuesSpiesDestroys.Int64)}
	}
	if obj.VictimCaseForm2AggressorCapableOfKilling.Valid {
		dto.VictimCaseForm2AggressorCapableOfKilling = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorCapableOfKilling.Int64)}
	}
	if obj.VictimCaseForm2AggressorHasAccessToWeapons.Valid {
		dto.VictimCaseForm2AggressorHasAccessToWeapons = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorHasAccessToWeapons.Int64)}
	}
	if obj.VictimCaseForm2PartnerUnemployed.Valid {
		dto.VictimCaseForm2PartnerUnemployed = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PartnerUnemployed.Int64)}
	}
	if obj.VictimCaseForm2PartnerOtherDenunciations.Valid {
		dto.VictimCaseForm2PartnerOtherDenunciations = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PartnerOtherDenunciations.Int64)}
	}
	if obj.VictimCaseForm2AggressorHasPenalBackground.Valid {
		dto.VictimCaseForm2AggressorHasPenalBackground = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorHasPenalBackground.Int64)}
	}
	if obj.VictimCaseForm2AggressorForcedSex.Valid {
		dto.VictimCaseForm2AggressorForcedSex = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorForcedSex.Int64)}
	}
	if obj.VictimCaseForm2AggressorAttemptedStrangulation.Valid {
		dto.VictimCaseForm2AggressorAttemptedStrangulation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorAttemptedStrangulation.Int64)}
	}
	if obj.VictimCaseForm2AggressorConsumesDrugs.Valid {
		dto.VictimCaseForm2AggressorConsumesDrugs = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorConsumesDrugs.Int64)}
	}
	if obj.VictimCaseForm2AggressorIsAlcoholic.Valid {
		dto.VictimCaseForm2AggressorIsAlcoholic = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorIsAlcoholic.Int64)}
	}
	if obj.VictimCaseForm2PartnerControls.Valid {
		dto.VictimCaseForm2PartnerControls = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PartnerControls.Int64)}
	}
	if obj.VictimCaseForm2AggressorHadHitInVulnerability.Valid {
		dto.VictimCaseForm2AggressorHadHitInVulnerability = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorHadHitInVulnerability.Int64)}
	}
	if obj.VictimCaseForm2PartnerThreatenedSuicide.Valid {
		dto.VictimCaseForm2PartnerThreatenedSuicide = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PartnerThreatenedSuicide.Int64)}
	}
	if obj.VictimCaseForm2PartnerThreatenedDamageMembers.Valid {
		dto.VictimCaseForm2PartnerThreatenedDamageMembers = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PartnerThreatenedDamageMembers.Int64)}
	}
	if obj.VictimCaseForm2ThoughtsOfSelfHarm.Valid {
		dto.VictimCaseForm2ThoughtsOfSelfHarm = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ThoughtsOfSelfHarm.Int64)}
	}
	if obj.VictimCaseForm2AggressorLimitsContactSupportNetworks.Valid {
		dto.VictimCaseForm2AggressorLimitsContactSupportNetworks = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorLimitsContactSupportNetworks.Int64)}
	}
	if obj.VictimCaseForm2StillLivesWithAggressor.Valid {
		dto.VictimCaseForm2StillLivesWithAggressor = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2StillLivesWithAggressor.Int64)}
	}
	if obj.VictimCaseForm2AggressorViolentlyJealous.Valid {
		dto.VictimCaseForm2AggressorViolentlyJealous = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorViolentlyJealous.Int64)}
	}
	if obj.VictimCaseForm2AggressorUnemployed.Valid {
		dto.VictimCaseForm2AggressorUnemployed = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorUnemployed.Int64)}
	}
	if obj.VictimCaseForm2AggressorHasPenalBackground2.Valid {
		dto.VictimCaseForm2AggressorHasPenalBackground2 = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorHasPenalBackground2.Int64)}
	}
	if obj.VictimCaseForm2AggressorSexuallyHarassment.Valid {
		dto.VictimCaseForm2AggressorSexuallyHarassment = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorSexuallyHarassment.Int64)}
	}
	if obj.VictimCaseForm2AggressorUseDrugs.Valid {
		dto.VictimCaseForm2AggressorUseDrugs = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorUseDrugs.Int64)}
	}
	if obj.VictimCaseForm2AggressorIsAlcoholic2.Valid {
		dto.VictimCaseForm2AggressorIsAlcoholic2 = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorIsAlcoholic2.Int64)}
	}
	if obj.VictimCaseForm2AggressorControls.Valid {
		dto.VictimCaseForm2AggressorControls = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorControls.Int64)}
	}
	if obj.VictimCaseForm2AggressorThreatenedDamageMembers.Valid {
		dto.VictimCaseForm2AggressorThreatenedDamageMembers = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorThreatenedDamageMembers.Int64)}
	}
	if obj.VictimCaseForm2ThoughtsOfSelfHarm2.Valid {
		dto.VictimCaseForm2ThoughtsOfSelfHarm2 = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ThoughtsOfSelfHarm2.Int64)}
	}
	if obj.VictimCaseForm2AggressorCommonSpaces.Valid {
		dto.VictimCaseForm2AggressorCommonSpaces = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorCommonSpaces.Int64)}
	}
	if obj.VictimCaseForm2AggressorHierarchy.Valid {
		dto.VictimCaseForm2AggressorHierarchy = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorHierarchy.Int64)}
	}

	if obj.VictimCaseForm2BirthDate.Valid {
		dto.VictimCaseForm2BirthDate = obj.VictimCaseForm2BirthDate.Time
	}
	if obj.VictimCaseForm2PhysicalMentalSensoryDifficulties.Valid {
		dto.VictimCaseForm2PhysicalMentalSensoryDifficulties = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PhysicalMentalSensoryDifficulties.Int64)}
	}
	if obj.VictimCaseForm2Nationality.Valid {
		dto.VictimCaseForm2Nationality = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2Nationality.Int64)}
	}
	if obj.VictimCaseForm2SpecifiedNationality.Valid {
		dto.VictimCaseForm2SpecifiedNationality = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2SpecifiedNationality.Int64)}
	}
	if obj.VictimCaseForm2MigrationCondition.Valid {
		dto.VictimCaseForm2MigrationCondition = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2MigrationCondition.Int64)}
	}
	if obj.VictimCaseForm2GenderIdentity.Valid {
		dto.VictimCaseForm2GenderIdentity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2GenderIdentity.Int64)}
	}
	if obj.VictimCaseForm2SexualOrientation.Valid {
		dto.VictimCaseForm2SexualOrientation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2SexualOrientation.Int64)}
	}
	if obj.VictimCaseForm2AssignedSexAtBirth.Valid {
		dto.VictimCaseForm2AssignedSexAtBirth = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AssignedSexAtBirth.Int64)}
	}
	if obj.VictimCaseForm2EthnicAffiliation.Valid {
		dto.VictimCaseForm2EthnicAffiliation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2EthnicAffiliation.Int64)}
	}
	if obj.VictimCaseForm2IndigenousPeople.Valid {
		dto.VictimCaseForm2IndigenousPeople = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2IndigenousPeople.Int64)}
	}
	if obj.VictimCaseForm2CampesinoRecognition.Valid {
		dto.VictimCaseForm2CampesinoRecognition = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2CampesinoRecognition.Int64)}
	}
	if obj.VictimCaseForm2MaritalStatus.Valid {
		dto.VictimCaseForm2MaritalStatus = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2MaritalStatus.Int64)}
	}
	if obj.VictimCaseForm2LastEducationLevel.Valid {
		dto.VictimCaseForm2LastEducationLevel = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2LastEducationLevel.Int64)}
	}
	if obj.VictimCaseForm2Occupation.Valid {
		dto.VictimCaseForm2Occupation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2Occupation.Int64)}
	}
	if obj.VictimCaseForm2IncomeGenerationMethod.Valid {
		dto.VictimCaseForm2IncomeGenerationMethod = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2IncomeGenerationMethod.Int64)}
	}

	if obj.VictimCaseForm2EmploymentRelationship.Valid {
		dto.VictimCaseForm2EmploymentRelationship = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2EmploymentRelationship.Int64)}
	}

	if obj.VictimCaseForm2ApproxStartAsp.Valid {
		dto.VictimCaseForm2ApproxStartAsp = obj.VictimCaseForm2ApproxStartAsp.Time
	}
	if obj.VictimCaseForm2HousingTenancyForm.Valid {
		dto.VictimCaseForm2HousingTenancyForm = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2HousingTenancyForm.Int64)}
	}
	if obj.VictimCaseForm2HousingStratum.Valid {
		dto.VictimCaseForm2HousingStratum = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2HousingStratum.Int64)}
	}
	if obj.VictimCaseForm2CurrentlyPregnant.Valid {
		dto.VictimCaseForm2CurrentlyPregnant = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2CurrentlyPregnant.Int64)}
	}
	if obj.VictimCaseForm2ResidenceTownCode.Valid {
		dto.VictimCaseForm2ResidenceTownCode = obj.VictimCaseForm2ResidenceTownCode.String
	}
	if obj.VictimCaseForm2ResidenceAddress.Valid {
		dto.VictimCaseForm2ResidenceAddress = obj.VictimCaseForm2ResidenceAddress.String
	}
	if obj.VictimCaseForm2ResidenceZone.Valid {
		dto.VictimCaseForm2ResidenceZone = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ResidenceZone.Int64)}
	}
	if obj.VictimCaseForm2SupportContactNames.Valid {
		dto.VictimCaseForm2SupportContactNames = obj.VictimCaseForm2SupportContactNames.String
	}
	if obj.VictimCaseForm2SupportContactPhone.Valid {
		dto.VictimCaseForm2SupportContactPhone = uint64(obj.VictimCaseForm2SupportContactPhone.Int64)
	}
	if obj.VictimCaseForm2SupportContactEmail.Valid {
		dto.VictimCaseForm2SupportContactEmail = obj.VictimCaseForm2SupportContactEmail.String
	}
	if obj.VictimCaseForm2SupportContactKinship.Valid {
		dto.VictimCaseForm2SupportContactKinship = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2SupportContactKinship.Int64)}
	}

	if obj.VictimCaseForm2LanguageAssistance.Valid {
		dto.VictimCaseForm2LanguageAssistance = obj.VictimCaseForm2LanguageAssistance.String
	}

	if obj.VictimCaseForm2PersonWithDisability.Valid {
		dto.VictimCaseForm2PersonWithDisability = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2PersonWithDisability.Int64)}
	}

	if obj.VictimCaseForm2RequireLanguageInterpreter.Valid {
		dto.VictimCaseForm2RequireLanguageInterpreter = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2RequireLanguageInterpreter.Int64)}
	}

	if obj.VictimCaseForm2WorkplaceSectorOccurrence.Valid {
		dto.VictimCaseForm2WorkplaceSectorOccurrence = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2WorkplaceSectorOccurrence.Int64)}
	}

	if obj.VictimCaseForm2ViolenceMotivatedByGender.Valid {
		dto.VictimCaseForm2ViolenceMotivatedByGender = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ViolenceMotivatedByGender.Int64)}
	}

	if obj.VictimCaseForm2AttentionWasAppropriate.Valid {
		dto.VictimCaseForm2AttentionWasAppropriate = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AttentionWasAppropriate.Int64)}
	}

	if obj.VictimCaseForm2AggressorOccupation.Valid {
		dto.VictimCaseForm2AggressorOccupation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorOccupation.Int64)}
	}

	if obj.VictimCaseForm2StoppedSeekingHelp.Valid {
		dto.VictimCaseForm2StoppedSeekingHelp = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2StoppedSeekingHelp.Int64)}
	}

	if obj.VictimCaseForm2VictimHealthToBlackmail.Valid {
		dto.VictimCaseForm2VictimHealthToBlackmail = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2VictimHealthToBlackmail.Int64)}
	}

	if obj.VictimCaseForm2ThreatenedRevealSexualOrientation.Valid {
		dto.VictimCaseForm2ThreatenedRevealSexualOrientation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ThreatenedRevealSexualOrientation.Int64)}
	}

	if obj.VictimCaseForm2ViolenceMotivatedByGender2.Valid {
		dto.VictimCaseForm2ViolenceMotivatedByGender2 = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2ViolenceMotivatedByGender2.Int64)}
	}

	if obj.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability.Valid {
		dto.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability.Int64)}
	}

	if obj.VictimCaseForm2AggressorSexuallyHarassment2.Valid {
		dto.VictimCaseForm2AggressorSexuallyHarassment2 = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorSexuallyHarassment2.Int64)}
	}

	if obj.VictimCaseForm2AggressorUsedPositionAuthority.Valid {
		dto.VictimCaseForm2AggressorUsedPositionAuthority = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AggressorUsedPositionAuthority.Int64)}
	}

	if obj.VictimCaseForm2AllowsEasyReport.Valid {
		dto.VictimCaseForm2AllowsEasyReport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.VictimCaseForm2AllowsEasyReport.Int64)}
	}

	if obj.VictimCaseForm2RiskScore.Valid {
		dto.VictimCaseForm2RiskScore = obj.VictimCaseForm2RiskScore.Int64
	}

	if obj.VictimCaseForm2RiskLevel.Valid {
		dto.VictimCaseForm2RiskLevel = obj.VictimCaseForm2RiskLevel.Int64
	}

	if obj.VictimCaseForm2SalivaManagementExplanation.Valid {
		dto.VictimCaseForm2SalivaManagementExplanation = obj.VictimCaseForm2SalivaManagementExplanation.String
	}
	if obj.VictimCaseForm2VictimCase.Valid {
		dto.VictimCaseForm2VictimCase = VictimCaseDTO{VictimCaseId: uint64(obj.VictimCaseForm2VictimCase.Int64)}
	}

	return dto
}
