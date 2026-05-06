package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	FeminicideForm1EntityName string = "FeminicideForm1"
	FeminicideForm1JSONName   string = "form"
	FeminicideForm1DBName     string = "feminicide_form1"
	FeminicideForm1DBScheme   string = "salvia"

	FeminicideForm1FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FeminicideForm1Id":                                           {Name: "FeminicideForm1Id", DBName: "feminicide_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1ICode":                                        {Name: "FeminicideForm1ICode", DBName: "feminicide_form1_i_code", Alias: "", ModelType: "string", MinSize: 36, MaxSize: 36, Required: true},
		"FeminicideForm1CreationDate":                                 {Name: "FeminicideForm1CreationDate", DBName: "feminicide_form1_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1UpdateDate":                                   {Name: "FeminicideForm1UpdateDate", DBName: "feminicide_form1_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimIdentityName":                           {Name: "FeminicideForm1VictimIdentityName", DBName: "feminicide_form1_victim_identity_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"FeminicideForm1BirthDate":                                    {Name: "FeminicideForm1BirthDate", DBName: "feminicide_form1_birth_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1DeathDate":                                    {Name: "FeminicideForm1DeathDate", DBName: "feminicide_form1_death_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimAddress":                                {Name: "FeminicideForm1VictimAddress", DBName: "feminicide_form1_victim_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"FeminicideForm1VictimZone":                                   {Name: "FeminicideForm1VictimZone", DBName: "feminicide_form1_victim_zone", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimLivingTownCode":                         {Name: "FeminicideForm1VictimLivingTownCode", DBName: "feminicide_form1_victim_living_town_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 8, Required: true},
		"FeminicideForm1VictimMaritalStatus":                          {Name: "FeminicideForm1VictimMaritalStatus", DBName: "feminicide_form1_victim_marital_status", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimSex":                                    {Name: "FeminicideForm1VictimSex", DBName: "feminicide_form1_victim_sex", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimGenderIdentity":                         {Name: "FeminicideForm1VictimGenderIdentity", DBName: "feminicide_form1_victim_gender_identity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimSexualOrientation":                      {Name: "FeminicideForm1VictimSexualOrientation", DBName: "feminicide_form1_victim_sexual_orientation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimEthnicity":                              {Name: "FeminicideForm1VictimEthnicity", DBName: "feminicide_form1_victim_ethnicity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimIndigenousPeople":                       {Name: "FeminicideForm1VictimIndigenousPeople", DBName: "feminicide_form1_victim_indigenous_people", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideForm1VictimSpecialPopulation":                      {Name: "FeminicideForm1VictimSpecialPopulation", DBName: "feminicide_form1_victim_special_population", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimDisability":                             {Name: "FeminicideForm1VictimDisability", DBName: "feminicide_form1_victim_disability", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimDisabilityType":                         {Name: "FeminicideForm1VictimDisabilityType", DBName: "feminicide_form1_victim_disability_type", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1PresumedAggressorNames":                       {Name: "FeminicideForm1PresumedAggressorNames", DBName: "feminicide_form1_presumed_aggressor_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: true},
		"FeminicideForm1PresumedAggressorRelation":                    {Name: "FeminicideForm1PresumedAggressorRelation", DBName: "feminicide_form1_presumed_aggressor_relation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1PresumedAggressorKnownVGB":                    {Name: "FeminicideForm1PresumedAggressorKnownVGB", DBName: "feminicide_form1_presumed_aggressor_known_vgb", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantNames":                               {Name: "FeminicideForm1InformantNames", DBName: "feminicide_form1_informant_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: true},
		"FeminicideForm1InformantIdentityName":                        {Name: "FeminicideForm1InformantIdentityName", DBName: "feminicide_form1_informant_identity_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"FeminicideForm1InformantDocType":                             {Name: "FeminicideForm1InformantDocType", DBName: "feminicide_form1_informant_doc_type", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"FeminicideForm1InformantDocNumber":                           {Name: "FeminicideForm1InformantDocNumber", DBName: "feminicide_form1_informant_doc_number", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 10, Required: false},
		"FeminicideForm1InformantBirthDate":                           {Name: "FeminicideForm1InformantBirthDate", DBName: "feminicide_form1_informant_birth_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1SGSSSAffiliation":                             {Name: "FeminicideForm1SGSSSAffiliation", DBName: "feminicide_form1_s_g_s_s_s_affiliation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1EpsName":                                      {Name: "FeminicideForm1EpsName", DBName: "feminicide_form1_eps_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: false},
		"FeminicideForm1InformantAddress":                             {Name: "FeminicideForm1InformantAddress", DBName: "feminicide_form1_informant_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"FeminicideForm1InformantZone":                                {Name: "FeminicideForm1InformantZone", DBName: "feminicide_form1_informant_zone", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantLivingTownCode":                      {Name: "FeminicideForm1InformantLivingTownCode", DBName: "feminicide_form1_informant_living_town_code", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 8, Required: true},
		"FeminicideForm1InformantPhone":                               {Name: "FeminicideForm1InformantPhone", DBName: "feminicide_form1_informant_phone", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 10, Required: true},
		"FeminicideForm1EmergencyContactNames":                        {Name: "FeminicideForm1EmergencyContactNames", DBName: "feminicide_form1_emergency_contact_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: true},
		"FeminicideForm1EmergencyContactNumber":                       {Name: "FeminicideForm1EmergencyContactNumber", DBName: "feminicide_form1_emergency_contact_number", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 10, Required: true},
		"FeminicideForm1InformantSex":                                 {Name: "FeminicideForm1InformantSex", DBName: "feminicide_form1_informant_sex", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantGenderIdentity":                      {Name: "FeminicideForm1InformantGenderIdentity", DBName: "feminicide_form1_informant_gender_identity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantSexualOrientation":                   {Name: "FeminicideForm1InformantSexualOrientation", DBName: "feminicide_form1_informant_sexual_orientation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantEthnicity":                           {Name: "FeminicideForm1InformantEthnicity", DBName: "feminicide_form1_informant_ethnicity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantIndigenousPeople":                    {Name: "FeminicideForm1InformantIndigenousPeople", DBName: "feminicide_form1_informant_indigenous_people", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideForm1InformantMigratoryStatus":                     {Name: "FeminicideForm1InformantMigratoryStatus", DBName: "feminicide_form1_informant_migratory_status", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantMigrationSituation":                  {Name: "FeminicideForm1InformantMigrationSituation", DBName: "feminicide_form1_informant_migration_situation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantHighestEducationLevel":               {Name: "FeminicideForm1InformantHighestEducationLevel", DBName: "feminicide_form1_informant_highest_education_level", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantSpecialPopulation":                   {Name: "FeminicideForm1InformantSpecialPopulation", DBName: "feminicide_form1_informant_special_population", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantCurrentEmployment":                   {Name: "FeminicideForm1InformantCurrentEmployment", DBName: "feminicide_form1_informant_current_employment", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantEmploymentAccess":                    {Name: "FeminicideForm1InformantEmploymentAccess", DBName: "feminicide_form1_informant_employment_access", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantEmploymentImpactDescription":         {Name: "FeminicideForm1InformantEmploymentImpactDescription", DBName: "feminicide_form1_informant_employment_impact_description", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1InformantPrimaryOccupation":                   {Name: "FeminicideForm1InformantPrimaryOccupation", DBName: "feminicide_form1_informant_primary_occupation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantDisability":                          {Name: "FeminicideForm1InformantDisability", DBName: "feminicide_form1_informant_disability", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1InformantDisabilityType":                      {Name: "FeminicideForm1InformantDisabilityType", DBName: "feminicide_form1_informant_disability_type", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1SituationAfterFeminicide":                     {Name: "FeminicideForm1SituationAfterFeminicide", DBName: "feminicide_form1_situation_after_feminicide", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: true},
		"FeminicideForm1AdditionalInformation":                        {Name: "FeminicideForm1AdditionalInformation", DBName: "feminicide_form1_additional_information", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AnyAssistanceReceived":                        {Name: "FeminicideForm1AnyAssistanceReceived", DBName: "feminicide_form1_any_assistance_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1HouseholdExpenseResponsibility":               {Name: "FeminicideForm1HouseholdExpenseResponsibility", DBName: "feminicide_form1_household_expense_responsibility", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1PostDeathEconomicAssumption":                  {Name: "FeminicideForm1PostDeathEconomicAssumption", DBName: "feminicide_form1_post_death_economic_assumption", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1EconomicAssumptionExplanation":                {Name: "FeminicideForm1EconomicAssumptionExplanation", DBName: "feminicide_form1_economic_assumption_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AnyDependentPeople":                           {Name: "FeminicideForm1AnyDependentPeople", DBName: "feminicide_form1_any_dependent_people", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AnyPublicOrPrivateEntity":                     {Name: "FeminicideForm1AnyPublicOrPrivateEntity", DBName: "feminicide_form1_any_public_or_private_entity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AnyPublicOrPrivateEntityExplanation":          {Name: "FeminicideForm1AnyPublicOrPrivateEntityExplanation", DBName: "feminicide_form1_any_public_or_private_entity_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1PublicTransportAccess":                        {Name: "FeminicideForm1PublicTransportAccess", DBName: "feminicide_form1_public_transport_access", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1PreferredTransportationMode":                  {Name: "FeminicideForm1PreferredTransportationMode", DBName: "feminicide_form1_preferred_transportation_mode", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1PreferredTransportationModeExplanation":       {Name: "FeminicideForm1PreferredTransportationModeExplanation", DBName: "feminicide_form1_preferred_transportation_mode_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1TransportationCostEstimate":                   {Name: "FeminicideForm1TransportationCostEstimate", DBName: "feminicide_form1_transportation_cost_estimate", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1TransportDifficulty":                          {Name: "FeminicideForm1TransportDifficulty", DBName: "feminicide_form1_transport_difficulty", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1TransportDifficultyExplanation":               {Name: "FeminicideForm1TransportDifficultyExplanation", DBName: "feminicide_form1_transport_difficulty_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1EconomicResourcesForTransport":                {Name: "FeminicideForm1EconomicResourcesForTransport", DBName: "feminicide_form1_economic_resources_for_transport", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1DebtOrHelpDueToTransport":                     {Name: "FeminicideForm1DebtOrHelpDueToTransport", DBName: "feminicide_form1_debt_or_help_due_to_transport", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1DebtImpactExplanation":                        {Name: "FeminicideForm1DebtImpactExplanation", DBName: "feminicide_form1_debt_impact_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1TransportSubsidyReceived":                     {Name: "FeminicideForm1TransportSubsidyReceived", DBName: "feminicide_form1_transport_subsidy_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1TransportSubsidyExplanation":                  {Name: "FeminicideForm1TransportSubsidyExplanation", DBName: "feminicide_form1_transport_subsidy_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1SafetyTransportationConcern":                  {Name: "FeminicideForm1SafetyTransportationConcern", DBName: "feminicide_form1_safety_transportation_concern", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1SafetyTransportationExplanation":              {Name: "FeminicideForm1SafetyTransportationExplanation", DBName: "feminicide_form1_safety_transportation_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1FoodAccessFrequency":                          {Name: "FeminicideForm1FoodAccessFrequency", DBName: "feminicide_form1_food_access_frequency", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1FoodAccessExplanation":                        {Name: "FeminicideForm1FoodAccessExplanation", DBName: "feminicide_form1_food_access_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AggressorFoodRestriction":                     {Name: "FeminicideForm1AggressorFoodRestriction", DBName: "feminicide_form1_aggressor_food_restriction", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AggressorFoodRestrictionExplanation":          {Name: "FeminicideForm1AggressorFoodRestrictionExplanation", DBName: "feminicide_form1_aggressor_food_restriction_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1FoodIncomeSupport":                            {Name: "FeminicideForm1FoodIncomeSupport", DBName: "feminicide_form1_food_income_support", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1FoodIncomeSupportExplanation":                 {Name: "FeminicideForm1FoodIncomeSupportExplanation", DBName: "feminicide_form1_food_income_support_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1JuridicalAssistanceReceived":                  {Name: "FeminicideForm1JuridicalAssistanceReceived", DBName: "feminicide_form1_juridical_assistance_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1JuridicalAssistanceExplanation":               {Name: "FeminicideForm1JuridicalAssistanceExplanation", DBName: "feminicide_form1_juridical_assistance_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1VictimRepresentation":                         {Name: "FeminicideForm1VictimRepresentation", DBName: "feminicide_form1_victim_representation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1VictimRepresentationExplanation":              {Name: "FeminicideForm1VictimRepresentationExplanation", DBName: "feminicide_form1_victim_representation_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1PsychosocialSupportReceived":                  {Name: "FeminicideForm1PsychosocialSupportReceived", DBName: "feminicide_form1_psychosocial_support_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1PsychosocialSupportExplanation":               {Name: "FeminicideForm1PsychosocialSupportExplanation", DBName: "feminicide_form1_psychosocial_support_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1EmergencyEmotionalCrisis":                     {Name: "FeminicideForm1EmergencyEmotionalCrisis", DBName: "feminicide_form1_emergency_emotional_crisis", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1EmergencyEmotionalCrisisExplanation":          {Name: "FeminicideForm1EmergencyEmotionalCrisisExplanation", DBName: "feminicide_form1_emergency_emotional_crisis_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AggressorSameResidence":                       {Name: "FeminicideForm1AggressorSameResidence", DBName: "feminicide_form1_aggressor_same_residence", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AggressorSameResidenceExplanation":            {Name: "FeminicideForm1AggressorSameResidenceExplanation", DBName: "feminicide_form1_aggressor_same_residence_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AggressorLocationKnown":                       {Name: "FeminicideForm1AggressorLocationKnown", DBName: "feminicide_form1_aggressor_location_known", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AggressorLocationKnownExplanation":            {Name: "FeminicideForm1AggressorLocationKnownExplanation", DBName: "feminicide_form1_aggressor_location_known_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AnyTypeOfAssistanceReceived":                  {Name: "FeminicideForm1AnyTypeOfAssistanceReceived", DBName: "feminicide_form1_any_type_of_assistance_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AnyTypeOfAssistanceReceivedExplanation":       {Name: "FeminicideForm1AnyTypeOfAssistanceReceivedExplanation", DBName: "feminicide_form1_any_type_of_assistance_received_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1CompensationFundInsuranceCoverage":            {Name: "FeminicideForm1CompensationFundInsuranceCoverage", DBName: "feminicide_form1_compensation_fund_insurance_coverage", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1CompensationFundInsuranceCoverageExplanation": {Name: "FeminicideForm1CompensationFundInsuranceCoverageExplanation", DBName: "feminicide_form1_compensation_fund_insurance_cov_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1FuneralSubsidyReceived":                       {Name: "FeminicideForm1FuneralSubsidyReceived", DBName: "feminicide_form1_funeral_subsidy_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1FuneralSubsidyExplanation":                    {Name: "FeminicideForm1FuneralSubsidyExplanation", DBName: "feminicide_form1_funeral_subsidy_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1FuneralFundsAvailable":                        {Name: "FeminicideForm1FuneralFundsAvailable", DBName: "feminicide_form1_funeral_funds_available", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1FuneralFundsAvailableExplanation":             {Name: "FeminicideForm1FuneralFundsAvailableExplanation", DBName: "feminicide_form1_funeral_funds_available_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1FuneralCostValue":                             {Name: "FeminicideForm1FuneralCostValue", DBName: "feminicide_form1_funeral_cost_value", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1RenameReputationImpact":                       {Name: "FeminicideForm1RenameReputationImpact", DBName: "feminicide_form1_rename_reputation_impact", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1RenameReputationExplanation":                  {Name: "FeminicideForm1RenameReputationExplanation", DBName: "feminicide_form1_rename_reputation_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1AdditionalNeedsDescription":                   {Name: "FeminicideForm1AdditionalNeedsDescription", DBName: "feminicide_form1_additional_needs_description", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideForm1ActionPlan":                                   {Name: "FeminicideForm1ActionPlan", DBName: "feminicide_form1_action_plan", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideForm1AnyAssistanceReceivedCityHall":                {Name: "FeminicideForm1AnyAssistanceReceivedCityHall", DBName: "feminicide_form1_any_assistance_received_city_hall", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideForm1AnyAssistanceReceivedWomensOffice":            {Name: "FeminicideForm1AnyAssistanceReceivedWomensOffice", DBName: "feminicide_form1_any_assistance_received_womens_office", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideForm1AnyAssistanceReceivedOtherEntity":             {Name: "FeminicideForm1AnyAssistanceReceivedOtherEntity", DBName: "feminicide_form1_any_assistance_received_other_entity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideForm1AnyAssistanceReceivedOther":                   {Name: "FeminicideForm1AnyAssistanceReceivedOther", DBName: "feminicide_form1_any_assistance_received_other", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideForm1AssistanceReceived":                           {Name: "FeminicideForm1AssistanceReceived", DBName: "feminicide_form1_assistance_received", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideForm1AssistanceReceivedCityHall":                   {Name: "FeminicideForm1AssistanceReceivedCityHall", DBName: "feminicide_form1_assistance_received_city_hall", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideForm1AssistanceReceivedWomensOffice":               {Name: "FeminicideForm1AssistanceReceivedWomensOffice", DBName: "feminicide_form1_assistance_received_womens_office", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideForm1AssistanceReceivedOtherEntity":                {Name: "FeminicideForm1AssistanceReceivedOtherEntity", DBName: "feminicide_form1_assistance_received_other_entity", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideForm1AssistanceReceivedOther":                      {Name: "FeminicideForm1AssistanceReceivedOther", DBName: "feminicide_form1_assistance_received_other", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},

		"FeminicideForm1FamilyMother":       {Name: "FeminicideForm1FamilyMother", DBName: "feminicide_form1_family_mother", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyFather":       {Name: "FeminicideForm1FamilyFather", DBName: "feminicide_form1_family_father", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyStepfather":   {Name: "FeminicideForm1FamilyStepfather", DBName: "feminicide_form1_family_stepfather", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyStepmother":   {Name: "FeminicideForm1FamilyStepmother", DBName: "feminicide_form1_family_stepmother", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyPartner":      {Name: "FeminicideForm1FamilyPartner", DBName: "feminicide_form1_family_partner", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySibling1":     {Name: "FeminicideForm1FamilySibling1", DBName: "feminicide_form1_family_sibling_1", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySibling2":     {Name: "FeminicideForm1FamilySibling2", DBName: "feminicide_form1_family_sibling_2", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySibling3":     {Name: "FeminicideForm1FamilySibling3", DBName: "feminicide_form1_family_sibling_3", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySibling4":     {Name: "FeminicideForm1FamilySibling4", DBName: "feminicide_form1_family_sibling_4", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySibling5":     {Name: "FeminicideForm1FamilySibling5", DBName: "feminicide_form1_family_sibling_5", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySonDaughter1": {Name: "FeminicideForm1FamilySonDaughter1", DBName: "feminicide_form1_family_son_daughter_1", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySonDaughter2": {Name: "FeminicideForm1FamilySonDaughter2", DBName: "feminicide_form1_family_son_daughter_2", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySonDaughter3": {Name: "FeminicideForm1FamilySonDaughter3", DBName: "feminicide_form1_family_son_daughter_3", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySonDaughter4": {Name: "FeminicideForm1FamilySonDaughter4", DBName: "feminicide_form1_family_son_daughter_4", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilySonDaughter5": {Name: "FeminicideForm1FamilySonDaughter5", DBName: "feminicide_form1_family_son_daughter_5", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyGrandmother":  {Name: "FeminicideForm1FamilyGrandmother", DBName: "feminicide_form1_family_grandmother", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyGrandfather":  {Name: "FeminicideForm1FamilyGrandfather", DBName: "feminicide_form1_family_grandfather", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},
		"FeminicideForm1FamilyOtherMember":  {Name: "FeminicideForm1FamilyOtherMember", DBName: "feminicide_form1_family_other_member", Alias: "", ModelType: "int", MinSize: 1, MaxSize: 120, Required: true},

		"FeminicideForm1Summary":    {Name: "FeminicideForm1Summary", DBName: "feminicide_form1_summary", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 20000, Required: false},
		"FeminicideForm1Feminicide": {Name: "FeminicideForm1Feminicide", DBName: "feminicide_form1_feminicide", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

type FeminicideForm1DTO struct {
	FeminicideForm1Id                      uint64                  `json:"-"`
	FeminicideForm1ICode                   string                  `json:"icode"`
	FeminicideForm1CreationDate            time.Time               `json:"creationDate"`
	FeminicideForm1UpdateDate              time.Time               `json:"updateDate"`
	FeminicideForm1VictimIdentityName      string                  `json:"victimIdentityName"`
	FeminicideForm1BirthDate               time.Time               `json:"birthDate"`
	FeminicideForm1DeathDate               time.Time               `json:"deathDate"`
	FeminicideForm1VictimAddress           string                  `json:"victimAddress"`
	FeminicideForm1VictimZone              VictimCaseForm2EnumsDTO `json:"victimZone"`
	FeminicideForm1VictimLivingTownCode    string                  `json:"victimLivingTownCode"`
	FeminicideForm1VictimSex               VictimCaseForm2EnumsDTO `json:"victimSex"`
	FeminicideForm1VictimMaritalStatus     VictimCaseForm2EnumsDTO `json:"victimMaritalStatus"`
	FeminicideForm1VictimGenderIdentity    VictimCaseForm2EnumsDTO `json:"victimGenderIdentity"`
	FeminicideForm1VictimSexualOrientation VictimCaseForm2EnumsDTO `json:"victimSexualOrientation"`
	FeminicideForm1VictimEthnicity         VictimCaseForm2EnumsDTO `json:"victimEthnicity"`
	FeminicideForm1VictimIndigenousPeople  VictimCaseForm2EnumsDTO `json:"victimIndigenousPeople"`

	FeminicideForm1VictimSpecialPopulation                      VictimCaseForm2EnumsDTO `json:"victimSpecialPopulation"`
	FeminicideForm1VictimDisability                             VictimCaseForm2EnumsDTO `json:"victimDisability"`
	FeminicideForm1VictimDisabilityType                         VictimCaseForm2EnumsDTO `json:"victimDisabilityType"`
	FeminicideForm1PresumedAggressorNames                       string                  `json:"presumedAggressorNames"`
	FeminicideForm1PresumedAggressorRelation                    VictimCaseForm2EnumsDTO `json:"presumedAggressorRelation"`
	FeminicideForm1PresumedAggressorKnownVGB                    VictimCaseForm2EnumsDTO `json:"presumedAggressorKnownVGB"`
	FeminicideForm1InformantNames                               string                  `json:"informantNames"`
	FeminicideForm1InformantIdentityName                        string                  `json:"informantIdentityName"`
	FeminicideForm1InformantDocType                             string                  `json:"informantDocType"`
	FeminicideForm1InformantDocNumber                           string                  `json:"informantDocNumber"`
	FeminicideForm1InformantBirthDate                           time.Time               `json:"informantBirthDate"`
	FeminicideForm1SGSSSAffiliation                             VictimCaseForm2EnumsDTO `json:"sgssAffiliation"`
	FeminicideForm1EpsName                                      string                  `json:"epsName"`
	FeminicideForm1InformantAddress                             string                  `json:"informantAddress"`
	FeminicideForm1InformantZone                                VictimCaseForm2EnumsDTO `json:"informantZone"`
	FeminicideForm1InformantLivingTownCode                      string                  `json:"informantLivingTownCode"`
	FeminicideForm1InformantPhone                               string                  `json:"informantPhone"`
	FeminicideForm1EmergencyContactNames                        string                  `json:"emergencyContactNames"`
	FeminicideForm1EmergencyContactNumber                       string                  `json:"emergencyContactNumber"`
	FeminicideForm1InformantSex                                 VictimCaseForm2EnumsDTO `json:"informantSex"`
	FeminicideForm1InformantGenderIdentity                      VictimCaseForm2EnumsDTO `json:"informantGenderIdentity"`
	FeminicideForm1InformantSexualOrientation                   VictimCaseForm2EnumsDTO `json:"informantSexualOrientation"`
	FeminicideForm1InformantEthnicity                           VictimCaseForm2EnumsDTO `json:"informantEthnicity"`
	FeminicideForm1InformantIndigenousPeople                    VictimCaseForm2EnumsDTO `json:"informantIndigenousPeople"`
	FeminicideForm1InformantMigratoryStatus                     VictimCaseForm2EnumsDTO `json:"informantMigratoryStatus"`
	FeminicideForm1InformantMigrationSituation                  VictimCaseForm2EnumsDTO `json:"informantMigrationSituation"`
	FeminicideForm1InformantHighestEducationLevel               VictimCaseForm2EnumsDTO `json:"informantHighestEducationLevel"`
	FeminicideForm1InformantSpecialPopulation                   VictimCaseForm2EnumsDTO `json:"informantSpecialPopulation"`
	FeminicideForm1InformantCurrentEmployment                   VictimCaseForm2EnumsDTO `json:"informantCurrentEmployment"`
	FeminicideForm1InformantEmploymentAccess                    VictimCaseForm2EnumsDTO `json:"informantEmploymentAccess"`
	FeminicideForm1InformantEmploymentImpactDescription         string                  `json:"informantEmploymentImpactDescription"`
	FeminicideForm1InformantPrimaryOccupation                   VictimCaseForm2EnumsDTO `json:"informantPrimaryOccupation"`
	FeminicideForm1InformantDisability                          VictimCaseForm2EnumsDTO `json:"informantDisability"`
	FeminicideForm1InformantDisabilityType                      VictimCaseForm2EnumsDTO `json:"informantDisabilityType"`
	FeminicideForm1SituationAfterFeminicide                     string                  `json:"situationAfterFeminicide"`
	FeminicideForm1AdditionalInformation                        string                  `json:"additionalInformation"`
	FeminicideForm1AnyAssistanceReceived                        VictimCaseForm2EnumsDTO `json:"anyAssistanceReceived"`
	FeminicideForm1HouseholdExpenseResponsibility               VictimCaseForm2EnumsDTO `json:"householdExpenseResponsibility"`
	FeminicideForm1PostDeathEconomicAssumption                  VictimCaseForm2EnumsDTO `json:"postDeathEconomicAssumption"`
	FeminicideForm1EconomicAssumptionExplanation                string                  `json:"economicAssumptionExplanation"`
	FeminicideForm1AnyDependentPeople                           VictimCaseForm2EnumsDTO `json:"anyDependentPeople"`
	FeminicideForm1AnyPublicOrPrivateEntity                     VictimCaseForm2EnumsDTO `json:"anyPublicOrPrivateEntity"`
	FeminicideForm1AnyPublicOrPrivateEntityExplanation          string                  `json:"anyPublicOrPrivateEntityExplanation"`
	FeminicideForm1PublicTransportAccess                        VictimCaseForm2EnumsDTO `json:"publicTransportAccess"`
	FeminicideForm1PreferredTransportationMode                  VictimCaseForm2EnumsDTO `json:"preferredTransportationMode"`
	FeminicideForm1PreferredTransportationModeExplanation       string                  `json:"preferredTransportationModeExplanation"`
	FeminicideForm1TransportationCostEstimate                   string                  `json:"transportationCostEstimate"`
	FeminicideForm1TransportDifficulty                          VictimCaseForm2EnumsDTO `json:"transportDifficulty"`
	FeminicideForm1TransportDifficultyExplanation               string                  `json:"transportDifficultyExplanation"`
	FeminicideForm1EconomicResourcesForTransport                VictimCaseForm2EnumsDTO `json:"economicResourcesForTransport"`
	FeminicideForm1DebtOrHelpDueToTransport                     VictimCaseForm2EnumsDTO `json:"debtOrHelpDueToTransport"`
	FeminicideForm1DebtImpactExplanation                        string                  `json:"debtImpactExplanation"`
	FeminicideForm1TransportSubsidyReceived                     VictimCaseForm2EnumsDTO `json:"transportSubsidyReceived"`
	FeminicideForm1TransportSubsidyExplanation                  string                  `json:"transportSubsidyExplanation"`
	FeminicideForm1SafetyTransportationConcern                  VictimCaseForm2EnumsDTO `json:"safetyTransportationConcern"`
	FeminicideForm1SafetyTransportationExplanation              string                  `json:"safetyTransportationExplanation"`
	FeminicideForm1FoodAccessFrequency                          VictimCaseForm2EnumsDTO `json:"foodAccessFrequency"`
	FeminicideForm1FoodAccessExplanation                        string                  `json:"foodAccessExplanation"`
	FeminicideForm1AggressorFoodRestriction                     VictimCaseForm2EnumsDTO `json:"aggressorFoodRestriction"`
	FeminicideForm1AggressorFoodRestrictionExplanation          string                  `json:"aggressorFoodRestrictionExplanation"`
	FeminicideForm1FoodIncomeSupport                            VictimCaseForm2EnumsDTO `json:"foodIncomeSupport"`
	FeminicideForm1FoodIncomeSupportExplanation                 string                  `json:"foodIncomeSupportExplanation"`
	FeminicideForm1JuridicalAssistanceReceived                  VictimCaseForm2EnumsDTO `json:"juridicalAssistanceReceived"`
	FeminicideForm1JuridicalAssistanceExplanation               string                  `json:"juridicalAssistanceExplanation"`
	FeminicideForm1VictimRepresentation                         VictimCaseForm2EnumsDTO `json:"victimRepresentation"`
	FeminicideForm1VictimRepresentationExplanation              string                  `json:"victimRepresentationExplanation"`
	FeminicideForm1PsychosocialSupportReceived                  VictimCaseForm2EnumsDTO `json:"psychosocialSupportReceived"`
	FeminicideForm1PsychosocialSupportExplanation               string                  `json:"psychosocialSupportExplanation"`
	FeminicideForm1EmergencyEmotionalCrisis                     VictimCaseForm2EnumsDTO `json:"emergencyEmotionalCrisis"`
	FeminicideForm1EmergencyEmotionalCrisisExplanation          string                  `json:"emergencyEmotionalCrisisExplanation"`
	FeminicideForm1AggressorSameResidence                       VictimCaseForm2EnumsDTO `json:"aggressorSameResidence"`
	FeminicideForm1AggressorSameResidenceExplanation            string                  `json:"aggressorSameResidenceExplanation"`
	FeminicideForm1AggressorLocationKnown                       VictimCaseForm2EnumsDTO `json:"aggressorLocationKnown"`
	FeminicideForm1AggressorLocationKnownExplanation            string                  `json:"aggressorLocationKnownExplanation"`
	FeminicideForm1AnyTypeOfAssistanceReceived                  VictimCaseForm2EnumsDTO `json:"anyTypeOfAssistanceReceived"`
	FeminicideForm1AnyTypeOfAssistanceReceivedExplanation       string                  `json:"anyTypeOfAssistanceReceivedExplanation"`
	FeminicideForm1CompensationFundInsuranceCoverage            VictimCaseForm2EnumsDTO `json:"compensationFundInsuranceCoverage"`
	FeminicideForm1CompensationFundInsuranceCoverageExplanation string                  `json:"compensationFundInsuranceCoverageExplanation"`
	FeminicideForm1FuneralSubsidyReceived                       VictimCaseForm2EnumsDTO `json:"funeralSubsidyReceived"`
	FeminicideForm1FuneralSubsidyExplanation                    string                  `json:"funeralSubsidyExplanation"`
	FeminicideForm1FuneralFundsAvailable                        VictimCaseForm2EnumsDTO `json:"funeralFundsAvailable"`
	FeminicideForm1FuneralFundsAvailableExplanation             string                  `json:"funeralFundsAvailableExplanation"`
	FeminicideForm1FuneralCostValue                             VictimCaseForm2EnumsDTO `json:"funeralCostValue"`
	FeminicideForm1RenameReputationImpact                       VictimCaseForm2EnumsDTO `json:"renameReputationImpact"`
	FeminicideForm1RenameReputationExplanation                  string                  `json:"renameReputationExplanation"`
	FeminicideForm1AdditionalNeedsDescription                   string                  `json:"additionalNeedsDescription"`
	FeminicideForm1ActionPlan                                   VictimCaseForm2EnumsDTO `json:"actionPlan"`

	FeminicideForm1AssistanceReceived                string                  `json:"assistanceReceived"`
	FeminicideForm1AnyAssistanceReceivedCityHall     VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedCityHall"`
	FeminicideForm1AssistanceReceivedCityHall        string                  `json:"assistanceReceivedCityHall"`
	FeminicideForm1AnyAssistanceReceivedWomensOffice VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedWomensOffice"`
	FeminicideForm1AssistanceReceivedWomensOffice    string                  `json:"assistanceReceivedWomensOffice"`
	FeminicideForm1AnyAssistanceReceivedOtherEntity  VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedOtherEntity"`
	FeminicideForm1AssistanceReceivedOtherEntity     string                  `json:"assistanceReceivedOtherEntity"`
	FeminicideForm1AnyAssistanceReceivedOther        VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedOther"`
	FeminicideForm1AssistanceReceivedOther           string                  `json:"assistanceReceivedOther"`

	FeminicideForm1FamilyFather       int64  `json:"familyFather"`
	FeminicideForm1FamilyMother       int64  `json:"familyMother"`
	FeminicideForm1FamilyStepfather   int64  `json:"familyStepfather"`
	FeminicideForm1FamilyStepmother   int64  `json:"familyStepmother"`
	FeminicideForm1FamilyPartner      int64  `json:"familyPartner"`
	FeminicideForm1FamilySibling1     int64  `json:"familySibling1"`
	FeminicideForm1FamilySibling2     int64  `json:"familySibling2"`
	FeminicideForm1FamilySibling3     int64  `json:"familySibling3"`
	FeminicideForm1FamilySibling4     int64  `json:"familySibling4"`
	FeminicideForm1FamilySibling5     int64  `json:"familySibling5"`
	FeminicideForm1FamilySonDaughter1 int64  `json:"familySonDaughter1"`
	FeminicideForm1FamilySonDaughter2 int64  `json:"familySonDaughter2"`
	FeminicideForm1FamilySonDaughter3 int64  `json:"familySonDaughter3"`
	FeminicideForm1FamilySonDaughter4 int64  `json:"familySonDaughter4"`
	FeminicideForm1FamilySonDaughter5 int64  `json:"familySonDaughter5"`
	FeminicideForm1FamilyGrandmother  int64  `json:"familyGrandmother"`
	FeminicideForm1FamilyGrandfather  int64  `json:"familyGrandfather"`
	FeminicideForm1FamilyOtherMember  string `json:"familyOtherMember"`

	FeminicideForm1Summary string `json:"summary"`

	//Campos múltiples
	FeminicideForm1EntityIntervened []VictimCaseForm2EnumsDTO `json:"entityIntervened"`

	FeminicideForm1Feminicide interface{} `json:"-"`

	//Campos de formulario

	FeminicideForm1VictimLivingDepartment    security_daos.DepartmentDTO
	FeminicideForm1VictimLivingCity          security_daos.CityDTO
	FeminicideForm1VictimLivingTown          security_daos.TownDTO
	FeminicideForm1InformantLivingDepartment security_daos.DepartmentDTO
	FeminicideForm1InformantLivingCity       security_daos.CityDTO
	FeminicideForm1InformantLivingTown       security_daos.TownDTO
}

type FeminicideForm1PgDB struct {
	FeminicideForm1Id                                           sql.NullInt64
	FeminicideForm1ICode                                        sql.NullString
	FeminicideForm1CreationDate                                 sql.NullTime
	FeminicideForm1UpdateDate                                   sql.NullTime
	FeminicideForm1VictimIdentityName                           sql.NullString
	FeminicideForm1BirthDate                                    sql.NullTime
	FeminicideForm1DeathDate                                    sql.NullTime
	FeminicideForm1VictimAddress                                sql.NullString
	FeminicideForm1VictimZone                                   sql.NullInt64
	FeminicideForm1VictimLivingTownCode                         sql.NullString
	FeminicideForm1VictimMaritalStatus                          sql.NullInt64
	FeminicideForm1VictimSex                                    sql.NullInt64
	FeminicideForm1VictimGenderIdentity                         sql.NullInt64
	FeminicideForm1VictimSexualOrientation                      sql.NullInt64
	FeminicideForm1VictimEthnicity                              sql.NullInt64
	FeminicideForm1VictimIndigenousPeople                       sql.NullInt64
	FeminicideForm1VictimSpecialPopulation                      sql.NullInt64
	FeminicideForm1VictimDisability                             sql.NullInt64
	FeminicideForm1VictimDisabilityType                         sql.NullInt64
	FeminicideForm1PresumedAggressorNames                       sql.NullString
	FeminicideForm1PresumedAggressorRelation                    sql.NullInt64
	FeminicideForm1PresumedAggressorKnownVGB                    sql.NullInt64
	FeminicideForm1InformantNames                               sql.NullString
	FeminicideForm1InformantIdentityName                        sql.NullString
	FeminicideForm1InformantDocType                             sql.NullString
	FeminicideForm1InformantDocNumber                           sql.NullString
	FeminicideForm1InformantBirthDate                           sql.NullTime
	FeminicideForm1SGSSSAffiliation                             sql.NullInt64
	FeminicideForm1EpsName                                      sql.NullString
	FeminicideForm1InformantAddress                             sql.NullString
	FeminicideForm1InformantZone                                sql.NullInt64
	FeminicideForm1InformantLivingTownCode                      sql.NullString
	FeminicideForm1InformantPhone                               sql.NullString
	FeminicideForm1EmergencyContactNames                        sql.NullString
	FeminicideForm1EmergencyContactNumber                       sql.NullString
	FeminicideForm1InformantSex                                 sql.NullInt64
	FeminicideForm1InformantGenderIdentity                      sql.NullInt64
	FeminicideForm1InformantSexualOrientation                   sql.NullInt64
	FeminicideForm1InformantEthnicity                           sql.NullInt64
	FeminicideForm1InformantIndigenousPeople                    sql.NullInt64
	FeminicideForm1InformantMigratoryStatus                     sql.NullInt64
	FeminicideForm1InformantMigrationSituation                  sql.NullInt64
	FeminicideForm1InformantHighestEducationLevel               sql.NullInt64
	FeminicideForm1InformantSpecialPopulation                   sql.NullInt64
	FeminicideForm1InformantCurrentEmployment                   sql.NullInt64
	FeminicideForm1InformantEmploymentAccess                    sql.NullInt64
	FeminicideForm1InformantEmploymentImpactDescription         sql.NullString
	FeminicideForm1InformantPrimaryOccupation                   sql.NullInt64
	FeminicideForm1InformantDisability                          sql.NullInt64
	FeminicideForm1InformantDisabilityType                      sql.NullInt64
	FeminicideForm1SituationAfterFeminicide                     sql.NullString
	FeminicideForm1AdditionalInformation                        sql.NullString
	FeminicideForm1AnyAssistanceReceived                        sql.NullInt64
	FeminicideForm1HouseholdExpenseResponsibility               sql.NullInt64
	FeminicideForm1PostDeathEconomicAssumption                  sql.NullInt64
	FeminicideForm1EconomicAssumptionExplanation                sql.NullString
	FeminicideForm1AnyDependentPeople                           sql.NullInt64
	FeminicideForm1AnyPublicOrPrivateEntity                     sql.NullInt64
	FeminicideForm1AnyPublicOrPrivateEntityExplanation          sql.NullString
	FeminicideForm1PublicTransportAccess                        sql.NullInt64
	FeminicideForm1PreferredTransportationMode                  sql.NullInt64
	FeminicideForm1PreferredTransportationModeExplanation       sql.NullString
	FeminicideForm1TransportationCostEstimate                   sql.NullString
	FeminicideForm1TransportDifficulty                          sql.NullInt64
	FeminicideForm1TransportDifficultyExplanation               sql.NullString
	FeminicideForm1EconomicResourcesForTransport                sql.NullInt64
	FeminicideForm1DebtOrHelpDueToTransport                     sql.NullInt64
	FeminicideForm1DebtImpactExplanation                        sql.NullString
	FeminicideForm1TransportSubsidyReceived                     sql.NullInt64
	FeminicideForm1TransportSubsidyExplanation                  sql.NullString
	FeminicideForm1SafetyTransportationConcern                  sql.NullInt64
	FeminicideForm1SafetyTransportationExplanation              sql.NullString
	FeminicideForm1FoodAccessFrequency                          sql.NullInt64
	FeminicideForm1FoodAccessExplanation                        sql.NullString
	FeminicideForm1AggressorFoodRestriction                     sql.NullInt64
	FeminicideForm1AggressorFoodRestrictionExplanation          sql.NullString
	FeminicideForm1FoodIncomeSupport                            sql.NullInt64
	FeminicideForm1FoodIncomeSupportExplanation                 sql.NullString
	FeminicideForm1JuridicalAssistanceReceived                  sql.NullInt64
	FeminicideForm1JuridicalAssistanceExplanation               sql.NullString
	FeminicideForm1VictimRepresentation                         sql.NullInt64
	FeminicideForm1VictimRepresentationExplanation              sql.NullString
	FeminicideForm1PsychosocialSupportReceived                  sql.NullInt64
	FeminicideForm1PsychosocialSupportExplanation               sql.NullString
	FeminicideForm1EmergencyEmotionalCrisis                     sql.NullInt64
	FeminicideForm1EmergencyEmotionalCrisisExplanation          sql.NullString
	FeminicideForm1AggressorSameResidence                       sql.NullInt64
	FeminicideForm1AggressorSameResidenceExplanation            sql.NullString
	FeminicideForm1AggressorLocationKnown                       sql.NullInt64
	FeminicideForm1AggressorLocationKnownExplanation            sql.NullString
	FeminicideForm1AnyTypeOfAssistanceReceived                  sql.NullInt64
	FeminicideForm1AnyTypeOfAssistanceReceivedExplanation       sql.NullString
	FeminicideForm1CompensationFundInsuranceCoverage            sql.NullInt64
	FeminicideForm1CompensationFundInsuranceCoverageExplanation sql.NullString
	FeminicideForm1FuneralSubsidyReceived                       sql.NullInt64
	FeminicideForm1FuneralSubsidyExplanation                    sql.NullString
	FeminicideForm1FuneralFundsAvailable                        sql.NullInt64
	FeminicideForm1FuneralFundsAvailableExplanation             sql.NullString
	FeminicideForm1FuneralCostValue                             sql.NullInt64
	FeminicideForm1RenameReputationImpact                       sql.NullInt64
	FeminicideForm1RenameReputationExplanation                  sql.NullString
	FeminicideForm1AdditionalNeedsDescription                   sql.NullString
	FeminicideForm1ActionPlan                                   sql.NullInt64
	FeminicideForm1Summary                                      sql.NullString
	FeminicideForm1Feminicide                                   sql.NullInt64
	FeminicideForm1AnyAssistanceReceivedCityHall                sql.NullInt64
	FeminicideForm1AnyAssistanceReceivedWomensOffice            sql.NullInt64
	FeminicideForm1AnyAssistanceReceivedOtherEntity             sql.NullInt64
	FeminicideForm1AnyAssistanceReceivedOther                   sql.NullInt64
	FeminicideForm1AssistanceReceived                           sql.NullString
	FeminicideForm1AssistanceReceivedCityHall                   sql.NullString
	FeminicideForm1AssistanceReceivedWomensOffice               sql.NullString
	FeminicideForm1AssistanceReceivedOtherEntity                sql.NullString
	FeminicideForm1AssistanceReceivedOther                      sql.NullString

	FeminicideForm1FamilyMother       sql.NullInt64
	FeminicideForm1FamilyFather       sql.NullInt64
	FeminicideForm1FamilyStepfather   sql.NullInt64
	FeminicideForm1FamilyStepmother   sql.NullInt64
	FeminicideForm1FamilyPartner      sql.NullInt64
	FeminicideForm1FamilySibling1     sql.NullInt64
	FeminicideForm1FamilySibling2     sql.NullInt64
	FeminicideForm1FamilySibling3     sql.NullInt64
	FeminicideForm1FamilySibling4     sql.NullInt64
	FeminicideForm1FamilySibling5     sql.NullInt64
	FeminicideForm1FamilySonDaughter1 sql.NullInt64
	FeminicideForm1FamilySonDaughter2 sql.NullInt64
	FeminicideForm1FamilySonDaughter3 sql.NullInt64
	FeminicideForm1FamilySonDaughter4 sql.NullInt64
	FeminicideForm1FamilySonDaughter5 sql.NullInt64
	FeminicideForm1FamilyGrandmother  sql.NullInt64
	FeminicideForm1FamilyGrandfather  sql.NullInt64
	FeminicideForm1FamilyOtherMember  sql.NullString
}

func (fcd FeminicideForm1DTO) MarshalJSON() ([]byte, error) {
	type Alias FeminicideForm1DTO

	return json.Marshal(&struct {
		*Alias
		FeminicideForm1BirthDate          string `json:"birthDate"`
		FeminicideForm1DeathDate          string `json:"deathDate"`
		FeminicideForm1CreationDate       string `json:"creationDate"`
		FeminicideForm1UpdateDate         string `json:"updateDate"`
		FeminicideForm1InformantBirthDate string `json:"informantBirthDate"`
	}{
		Alias:                             (*Alias)(&fcd),
		FeminicideForm1BirthDate:          fcd.FeminicideForm1BirthDate.Format(common_config.DateTime.DATE_FORMAT),
		FeminicideForm1DeathDate:          fcd.FeminicideForm1DeathDate.Format(common_config.DateTime.DATE_FORMAT),
		FeminicideForm1CreationDate:       fcd.FeminicideForm1CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FeminicideForm1UpdateDate:         fcd.FeminicideForm1UpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FeminicideForm1InformantBirthDate: fcd.FeminicideForm1InformantBirthDate.Format(common_config.DateTime.DATE_FORMAT),
	})
}

func (fcd *FeminicideForm1DTO) UnmarshalJSON(data []byte) error {
	type Alias FeminicideForm1DTO

	aux := &struct {
		*Alias
		FeminicideForm1BirthDate          string `json:"birthDate"`
		FeminicideForm1DeathDate          string `json:"deathDate"`
		FeminicideForm1CreationDate       string `json:"creationDate"`
		FeminicideForm1UpdateDate         string `json:"updateDate"`
		FeminicideForm1InformantBirthDate string `json:"informantBirthDate"`
	}{
		Alias: (*Alias)(fcd),
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
			return time.Time{}
		}
		return t
	}

	fcd.FeminicideForm1BirthDate = parse(aux.FeminicideForm1BirthDate, common_config.DateTime.DATE_FORMAT)
	fcd.FeminicideForm1DeathDate = parse(aux.FeminicideForm1DeathDate, common_config.DateTime.DATE_FORMAT)
	fcd.FeminicideForm1CreationDate = parse(aux.FeminicideForm1CreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	fcd.FeminicideForm1UpdateDate = parse(aux.FeminicideForm1UpdateDate, common_config.DateTime.DATE_TIME_FORMAT)
	fcd.FeminicideForm1InformantBirthDate = parse(aux.FeminicideForm1InformantBirthDate, common_config.DateTime.DATE_FORMAT)

	return nil
}

func SetFeminicideForm1(feminicideForm1 *FeminicideForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var feminicideForm1FieldsSlice []string = []string{
		"FeminicideForm1ICode", "FeminicideForm1CreationDate", "FeminicideForm1UpdateDate",
		"FeminicideForm1VictimIdentityName", "FeminicideForm1BirthDate", "FeminicideForm1DeathDate",
		"FeminicideForm1VictimAddress", "FeminicideForm1VictimZone", "FeminicideForm1VictimLivingTownCode",
		"FeminicideForm1VictimMaritalStatus", "FeminicideForm1VictimSex", "FeminicideForm1VictimGenderIdentity",
		"FeminicideForm1VictimSexualOrientation", "FeminicideForm1VictimEthnicity", "FeminicideForm1VictimIndigenousPeople", "FeminicideForm1VictimSpecialPopulation",
		"FeminicideForm1VictimDisability", "FeminicideForm1VictimDisabilityType", "FeminicideForm1PresumedAggressorNames",
		"FeminicideForm1PresumedAggressorRelation", "FeminicideForm1PresumedAggressorKnownVGB",
		"FeminicideForm1InformantNames", "FeminicideForm1InformantIdentityName", "FeminicideForm1InformantDocType",
		"FeminicideForm1InformantDocNumber", "FeminicideForm1InformantBirthDate", "FeminicideForm1SGSSSAffiliation",
		"FeminicideForm1EpsName", "FeminicideForm1InformantAddress", "FeminicideForm1InformantZone",
		"FeminicideForm1InformantLivingTownCode", "FeminicideForm1InformantPhone", "FeminicideForm1EmergencyContactNames",
		"FeminicideForm1EmergencyContactNumber", "FeminicideForm1InformantSex", "FeminicideForm1InformantGenderIdentity",
		"FeminicideForm1InformantSexualOrientation", "FeminicideForm1InformantEthnicity", "FeminicideForm1InformantIndigenousPeople", "FeminicideForm1InformantMigratoryStatus",
		"FeminicideForm1InformantMigrationSituation", "FeminicideForm1InformantHighestEducationLevel", "FeminicideForm1InformantSpecialPopulation", "FeminicideForm1InformantCurrentEmployment",
		"FeminicideForm1InformantEmploymentAccess", "FeminicideForm1InformantEmploymentImpactDescription", "FeminicideForm1InformantPrimaryOccupation",
		"FeminicideForm1InformantDisability", "FeminicideForm1InformantDisabilityType", "FeminicideForm1SituationAfterFeminicide",
		"FeminicideForm1AdditionalInformation", "FeminicideForm1AnyAssistanceReceived", "FeminicideForm1HouseholdExpenseResponsibility",
		"FeminicideForm1PostDeathEconomicAssumption", "FeminicideForm1EconomicAssumptionExplanation", "FeminicideForm1AnyDependentPeople",
		"FeminicideForm1AnyPublicOrPrivateEntity", "FeminicideForm1AnyPublicOrPrivateEntityExplanation", "FeminicideForm1PublicTransportAccess",
		"FeminicideForm1PreferredTransportationMode", "FeminicideForm1PreferredTransportationModeExplanation", "FeminicideForm1TransportationCostEstimate",
		"FeminicideForm1TransportDifficulty", "FeminicideForm1TransportDifficultyExplanation", "FeminicideForm1EconomicResourcesForTransport",
		"FeminicideForm1DebtOrHelpDueToTransport", "FeminicideForm1DebtImpactExplanation", "FeminicideForm1TransportSubsidyReceived",
		"FeminicideForm1TransportSubsidyExplanation", "FeminicideForm1SafetyTransportationConcern", "FeminicideForm1SafetyTransportationExplanation",
		"FeminicideForm1FoodAccessFrequency", "FeminicideForm1FoodAccessExplanation", "FeminicideForm1AggressorFoodRestriction",
		"FeminicideForm1AggressorFoodRestrictionExplanation", "FeminicideForm1FoodIncomeSupport", "FeminicideForm1FoodIncomeSupportExplanation",
		"FeminicideForm1JuridicalAssistanceReceived", "FeminicideForm1JuridicalAssistanceExplanation", "FeminicideForm1VictimRepresentation",
		"FeminicideForm1VictimRepresentationExplanation", "FeminicideForm1PsychosocialSupportReceived", "FeminicideForm1PsychosocialSupportExplanation",
		"FeminicideForm1EmergencyEmotionalCrisis", "FeminicideForm1EmergencyEmotionalCrisisExplanation", "FeminicideForm1AggressorSameResidence",
		"FeminicideForm1AggressorSameResidenceExplanation", "FeminicideForm1AggressorLocationKnown", "FeminicideForm1AggressorLocationKnownExplanation",
		"FeminicideForm1AnyTypeOfAssistanceReceived", "FeminicideForm1AnyTypeOfAssistanceReceivedExplanation", "FeminicideForm1CompensationFundInsuranceCoverage",
		"FeminicideForm1CompensationFundInsuranceCoverageExplanation", "FeminicideForm1FuneralSubsidyReceived", "FeminicideForm1FuneralSubsidyExplanation",
		"FeminicideForm1FuneralFundsAvailable", "FeminicideForm1FuneralFundsAvailableExplanation", "FeminicideForm1FuneralCostValue",
		"FeminicideForm1RenameReputationImpact", "FeminicideForm1RenameReputationExplanation", "FeminicideForm1AdditionalNeedsDescription",
		"FeminicideForm1ActionPlan", "FeminicideForm1AnyAssistanceReceivedCityHall", "FeminicideForm1AnyAssistanceReceivedWomensOffice", "FeminicideForm1AnyAssistanceReceivedOtherEntity", "FeminicideForm1AnyAssistanceReceivedOther",
		"FeminicideForm1AssistanceReceived", "FeminicideForm1AssistanceReceivedCityHall", "FeminicideForm1AssistanceReceivedWomensOffice", "FeminicideForm1AssistanceReceivedOtherEntity", "FeminicideForm1AssistanceReceivedOther",
		"FeminicideForm1FamilyMother", "FeminicideForm1FamilyFather", "FeminicideForm1FamilyStepfather", "FeminicideForm1FamilyStepmother",
		"FeminicideForm1FamilyPartner", "FeminicideForm1FamilySibling1", "FeminicideForm1FamilySibling2", "FeminicideForm1FamilySibling3",
		"FeminicideForm1FamilySibling4", "FeminicideForm1FamilySibling5", "FeminicideForm1FamilySonDaughter1", "FeminicideForm1FamilySonDaughter2",
		"FeminicideForm1FamilySonDaughter3", "FeminicideForm1FamilySonDaughter4", "FeminicideForm1FamilySonDaughter5", "FeminicideForm1FamilyGrandmother",
		"FeminicideForm1FamilyGrandfather", "FeminicideForm1FamilyOtherMember",
		"FeminicideForm1Summary", "FeminicideForm1Feminicide",
	}

	var feminicideForm1FieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, feminicideForm1FieldsSlice, feminicideForm1FieldsAliasSlice, FeminicideForm1DBName, []string{}, []string{}, []string{"FeminicideForm1Id"}, common_dao.SQL_AND, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		feminicideForm1.FeminicideForm1ICode,
		feminicideForm1.FeminicideForm1CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideForm1.FeminicideForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideForm1.FeminicideForm1VictimIdentityName,
		feminicideForm1.FeminicideForm1BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideForm1.FeminicideForm1DeathDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideForm1.FeminicideForm1VictimAddress,
		feminicideForm1.FeminicideForm1VictimZone.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimLivingTownCode,
		feminicideForm1.FeminicideForm1VictimMaritalStatus.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimSex.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimGenderIdentity.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimSexualOrientation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimEthnicity.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimIndigenousPeople.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimSpecialPopulation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimDisability.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimDisabilityType.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1PresumedAggressorNames,
		feminicideForm1.FeminicideForm1PresumedAggressorRelation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1PresumedAggressorKnownVGB.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantNames,
		feminicideForm1.FeminicideForm1InformantIdentityName,
		feminicideForm1.FeminicideForm1InformantDocType,
		feminicideForm1.FeminicideForm1InformantDocNumber,
		feminicideForm1.FeminicideForm1InformantBirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideForm1.FeminicideForm1SGSSSAffiliation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1EpsName,
		feminicideForm1.FeminicideForm1InformantAddress,
		feminicideForm1.FeminicideForm1InformantZone.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantLivingTownCode,
		feminicideForm1.FeminicideForm1InformantPhone,
		feminicideForm1.FeminicideForm1EmergencyContactNames,
		feminicideForm1.FeminicideForm1EmergencyContactNumber,
		feminicideForm1.FeminicideForm1InformantSex.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantGenderIdentity.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantSexualOrientation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantEthnicity.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantIndigenousPeople.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantMigratoryStatus.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantMigrationSituation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantHighestEducationLevel.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantSpecialPopulation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantCurrentEmployment.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantEmploymentAccess.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantEmploymentImpactDescription,
		feminicideForm1.FeminicideForm1InformantPrimaryOccupation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantDisability.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1InformantDisabilityType.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1SituationAfterFeminicide,
		feminicideForm1.FeminicideForm1AdditionalInformation,
		feminicideForm1.FeminicideForm1AnyAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1HouseholdExpenseResponsibility.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1PostDeathEconomicAssumption.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1EconomicAssumptionExplanation,
		feminicideForm1.FeminicideForm1AnyDependentPeople.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyPublicOrPrivateEntity.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyPublicOrPrivateEntityExplanation,
		feminicideForm1.FeminicideForm1PublicTransportAccess.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1PreferredTransportationMode.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1PreferredTransportationModeExplanation,
		feminicideForm1.FeminicideForm1TransportationCostEstimate,
		feminicideForm1.FeminicideForm1TransportDifficulty.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1TransportDifficultyExplanation,
		feminicideForm1.FeminicideForm1EconomicResourcesForTransport.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1DebtOrHelpDueToTransport.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1DebtImpactExplanation,
		feminicideForm1.FeminicideForm1TransportSubsidyReceived.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1TransportSubsidyExplanation,
		feminicideForm1.FeminicideForm1SafetyTransportationConcern.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1SafetyTransportationExplanation,
		feminicideForm1.FeminicideForm1FoodAccessFrequency.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1FoodAccessExplanation,
		feminicideForm1.FeminicideForm1AggressorFoodRestriction.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AggressorFoodRestrictionExplanation,
		feminicideForm1.FeminicideForm1FoodIncomeSupport.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1FoodIncomeSupportExplanation,
		feminicideForm1.FeminicideForm1JuridicalAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1JuridicalAssistanceExplanation,
		feminicideForm1.FeminicideForm1VictimRepresentation.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1VictimRepresentationExplanation,
		feminicideForm1.FeminicideForm1PsychosocialSupportReceived.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1PsychosocialSupportExplanation,
		feminicideForm1.FeminicideForm1EmergencyEmotionalCrisis.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1EmergencyEmotionalCrisisExplanation,
		feminicideForm1.FeminicideForm1AggressorSameResidence.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AggressorSameResidenceExplanation,
		feminicideForm1.FeminicideForm1AggressorLocationKnown.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AggressorLocationKnownExplanation,
		feminicideForm1.FeminicideForm1AnyTypeOfAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation,
		feminicideForm1.FeminicideForm1CompensationFundInsuranceCoverage.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1CompensationFundInsuranceCoverageExplanation,
		feminicideForm1.FeminicideForm1FuneralSubsidyReceived.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1FuneralSubsidyExplanation,
		feminicideForm1.FeminicideForm1FuneralFundsAvailable.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1FuneralFundsAvailableExplanation,
		feminicideForm1.FeminicideForm1FuneralCostValue.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1RenameReputationImpact.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1RenameReputationExplanation,
		feminicideForm1.FeminicideForm1AdditionalNeedsDescription,
		feminicideForm1.FeminicideForm1ActionPlan.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyAssistanceReceivedCityHall.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyAssistanceReceivedWomensOffice.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyAssistanceReceivedOtherEntity.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AnyAssistanceReceivedOther.VictimCaseForm2EnumsId,
		feminicideForm1.FeminicideForm1AssistanceReceived,
		feminicideForm1.FeminicideForm1AssistanceReceivedCityHall,
		feminicideForm1.FeminicideForm1AssistanceReceivedWomensOffice,
		feminicideForm1.FeminicideForm1AssistanceReceivedOtherEntity,
		feminicideForm1.FeminicideForm1AssistanceReceivedOther,
		feminicideForm1.FeminicideForm1FamilyFather,
		feminicideForm1.FeminicideForm1FamilyMother,
		feminicideForm1.FeminicideForm1FamilyStepfather,
		feminicideForm1.FeminicideForm1FamilyStepmother,
		feminicideForm1.FeminicideForm1FamilyPartner,
		feminicideForm1.FeminicideForm1FamilySibling1,
		feminicideForm1.FeminicideForm1FamilySibling2,
		feminicideForm1.FeminicideForm1FamilySibling3,
		feminicideForm1.FeminicideForm1FamilySibling4,
		feminicideForm1.FeminicideForm1FamilySibling5,
		feminicideForm1.FeminicideForm1FamilySonDaughter1,
		feminicideForm1.FeminicideForm1FamilySonDaughter2,
		feminicideForm1.FeminicideForm1FamilySonDaughter3,
		feminicideForm1.FeminicideForm1FamilySonDaughter4,
		feminicideForm1.FeminicideForm1FamilySonDaughter5,
		feminicideForm1.FeminicideForm1FamilyGrandmother,
		feminicideForm1.FeminicideForm1FamilyGrandfather,
		feminicideForm1.FeminicideForm1FamilyOtherMember,
		feminicideForm1.FeminicideForm1Summary,
		feminicideForm1.FeminicideForm1Feminicide.(FeminicideDTO).FeminicideId)

	persistenceCtrl.Scan(&feminicideForm1.FeminicideForm1Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicideForm1(by common_controllers.By, feminicideForm1 *FeminicideForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicideForm1Path string = FeminicideForm1DBScheme + "." + FeminicideForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var feminicideForm1FieldsSlice []string = []string{
		"FeminicideForm1Id", "FeminicideForm1ICode", "FeminicideForm1CreationDate", "FeminicideForm1UpdateDate",
		"FeminicideForm1VictimIdentityName", "FeminicideForm1BirthDate", "FeminicideForm1DeathDate",
		"FeminicideForm1VictimAddress", "FeminicideForm1VictimZone", "FeminicideForm1VictimLivingTownCode",
		"FeminicideForm1VictimMaritalStatus", "FeminicideForm1VictimSex", "FeminicideForm1VictimGenderIdentity",
		"FeminicideForm1VictimSexualOrientation", "FeminicideForm1VictimEthnicity", "FeminicideForm1VictimIndigenousPeople", "FeminicideForm1VictimSpecialPopulation",
		"FeminicideForm1VictimDisability", "FeminicideForm1VictimDisabilityType", "FeminicideForm1PresumedAggressorNames",
		"FeminicideForm1PresumedAggressorRelation", "FeminicideForm1PresumedAggressorKnownVGB",
		"FeminicideForm1InformantNames", "FeminicideForm1InformantIdentityName", "FeminicideForm1InformantDocType",
		"FeminicideForm1InformantDocNumber", "FeminicideForm1InformantBirthDate", "FeminicideForm1SGSSSAffiliation",
		"FeminicideForm1EpsName", "FeminicideForm1InformantAddress", "FeminicideForm1InformantZone",
		"FeminicideForm1InformantLivingTownCode", "FeminicideForm1InformantPhone", "FeminicideForm1EmergencyContactNames",
		"FeminicideForm1EmergencyContactNumber", "FeminicideForm1InformantSex", "FeminicideForm1InformantGenderIdentity",
		"FeminicideForm1InformantSexualOrientation", "FeminicideForm1InformantEthnicity", "FeminicideForm1InformantIndigenousPeople", "FeminicideForm1InformantMigratoryStatus",
		"FeminicideForm1InformantMigrationSituation", "FeminicideForm1InformantHighestEducationLevel", "FeminicideForm1InformantSpecialPopulation", "FeminicideForm1InformantCurrentEmployment",
		"FeminicideForm1InformantEmploymentAccess", "FeminicideForm1InformantEmploymentImpactDescription", "FeminicideForm1InformantPrimaryOccupation",
		"FeminicideForm1InformantDisability", "FeminicideForm1InformantDisabilityType", "FeminicideForm1SituationAfterFeminicide",
		"FeminicideForm1AdditionalInformation", "FeminicideForm1AnyAssistanceReceived", "FeminicideForm1HouseholdExpenseResponsibility",
		"FeminicideForm1PostDeathEconomicAssumption", "FeminicideForm1EconomicAssumptionExplanation", "FeminicideForm1AnyDependentPeople",
		"FeminicideForm1AnyPublicOrPrivateEntity", "FeminicideForm1AnyPublicOrPrivateEntityExplanation", "FeminicideForm1PublicTransportAccess",
		"FeminicideForm1PreferredTransportationMode", "FeminicideForm1PreferredTransportationModeExplanation", "FeminicideForm1TransportationCostEstimate",
		"FeminicideForm1TransportDifficulty", "FeminicideForm1TransportDifficultyExplanation", "FeminicideForm1EconomicResourcesForTransport",
		"FeminicideForm1DebtOrHelpDueToTransport", "FeminicideForm1DebtImpactExplanation", "FeminicideForm1TransportSubsidyReceived",
		"FeminicideForm1TransportSubsidyExplanation", "FeminicideForm1SafetyTransportationConcern", "FeminicideForm1SafetyTransportationExplanation",
		"FeminicideForm1FoodAccessFrequency", "FeminicideForm1FoodAccessExplanation", "FeminicideForm1AggressorFoodRestriction",
		"FeminicideForm1AggressorFoodRestrictionExplanation", "FeminicideForm1FoodIncomeSupport", "FeminicideForm1FoodIncomeSupportExplanation",
		"FeminicideForm1JuridicalAssistanceReceived", "FeminicideForm1JuridicalAssistanceExplanation", "FeminicideForm1VictimRepresentation",
		"FeminicideForm1VictimRepresentationExplanation", "FeminicideForm1PsychosocialSupportReceived", "FeminicideForm1PsychosocialSupportExplanation",
		"FeminicideForm1EmergencyEmotionalCrisis", "FeminicideForm1EmergencyEmotionalCrisisExplanation", "FeminicideForm1AggressorSameResidence",
		"FeminicideForm1AggressorSameResidenceExplanation", "FeminicideForm1AggressorLocationKnown", "FeminicideForm1AggressorLocationKnownExplanation",
		"FeminicideForm1AnyTypeOfAssistanceReceived", "FeminicideForm1AnyTypeOfAssistanceReceivedExplanation", "FeminicideForm1CompensationFundInsuranceCoverage",
		"FeminicideForm1CompensationFundInsuranceCoverageExplanation", "FeminicideForm1FuneralSubsidyReceived", "FeminicideForm1FuneralSubsidyExplanation",
		"FeminicideForm1FuneralFundsAvailable", "FeminicideForm1FuneralFundsAvailableExplanation", "FeminicideForm1FuneralCostValue",
		"FeminicideForm1RenameReputationImpact", "FeminicideForm1RenameReputationExplanation", "FeminicideForm1AdditionalNeedsDescription",
		"FeminicideForm1ActionPlan", "FeminicideForm1AnyAssistanceReceivedCityHall", "FeminicideForm1AnyAssistanceReceivedWomensOffice", "FeminicideForm1AnyAssistanceReceivedOtherEntity", "FeminicideForm1AnyAssistanceReceivedOther",
		"FeminicideForm1AssistanceReceived", "FeminicideForm1AssistanceReceivedCityHall", "FeminicideForm1AssistanceReceivedWomensOffice", "FeminicideForm1AssistanceReceivedOtherEntity", "FeminicideForm1AssistanceReceivedOther",
		"FeminicideForm1FamilyMother", "FeminicideForm1FamilyFather", "FeminicideForm1FamilyStepfather", "FeminicideForm1FamilyStepmother",
		"FeminicideForm1FamilyPartner", "FeminicideForm1FamilySibling1", "FeminicideForm1FamilySibling2", "FeminicideForm1FamilySibling3",
		"FeminicideForm1FamilySibling4", "FeminicideForm1FamilySibling5", "FeminicideForm1FamilySonDaughter1", "FeminicideForm1FamilySonDaughter2",
		"FeminicideForm1FamilySonDaughter3", "FeminicideForm1FamilySonDaughter4", "FeminicideForm1FamilySonDaughter5", "FeminicideForm1FamilyGrandmother",
		"FeminicideForm1FamilyGrandfather", "FeminicideForm1FamilyOtherMember",
		"FeminicideForm1Summary", "FeminicideForm1Feminicide",
	}
	var feminicideForm1FieldsAliasSlice []string = []string{}

	var feminicideForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, feminicideForm1FieldsSlice, feminicideForm1FieldsAliasSlice, FeminicideForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, true)

	var query string = `SELECT ` + feminicideForm1FieldsStr +
		` FROM ` + feminicideForm1Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, true)

	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var feminicideForm1Pg FeminicideForm1PgDB = FeminicideForm1PgDB{}

	persistenceCtrl.Scan(&feminicideForm1Pg.FeminicideForm1Id,
		&feminicideForm1Pg.FeminicideForm1ICode,
		&feminicideForm1Pg.FeminicideForm1CreationDate,
		&feminicideForm1Pg.FeminicideForm1UpdateDate,
		&feminicideForm1Pg.FeminicideForm1VictimIdentityName,
		&feminicideForm1Pg.FeminicideForm1BirthDate,
		&feminicideForm1Pg.FeminicideForm1DeathDate,
		&feminicideForm1Pg.FeminicideForm1VictimAddress,
		&feminicideForm1Pg.FeminicideForm1VictimZone,
		&feminicideForm1Pg.FeminicideForm1VictimLivingTownCode,
		&feminicideForm1Pg.FeminicideForm1VictimMaritalStatus,
		&feminicideForm1Pg.FeminicideForm1VictimSex,
		&feminicideForm1Pg.FeminicideForm1VictimGenderIdentity,
		&feminicideForm1Pg.FeminicideForm1VictimSexualOrientation,
		&feminicideForm1Pg.FeminicideForm1VictimEthnicity,
		&feminicideForm1Pg.FeminicideForm1VictimIndigenousPeople,
		&feminicideForm1Pg.FeminicideForm1VictimSpecialPopulation,
		&feminicideForm1Pg.FeminicideForm1VictimDisability,
		&feminicideForm1Pg.FeminicideForm1VictimDisabilityType,
		&feminicideForm1Pg.FeminicideForm1PresumedAggressorNames,
		&feminicideForm1Pg.FeminicideForm1PresumedAggressorRelation,
		&feminicideForm1Pg.FeminicideForm1PresumedAggressorKnownVGB,
		&feminicideForm1Pg.FeminicideForm1InformantNames,
		&feminicideForm1Pg.FeminicideForm1InformantIdentityName,
		&feminicideForm1Pg.FeminicideForm1InformantDocType,
		&feminicideForm1Pg.FeminicideForm1InformantDocNumber,
		&feminicideForm1Pg.FeminicideForm1InformantBirthDate,
		&feminicideForm1Pg.FeminicideForm1SGSSSAffiliation,
		&feminicideForm1Pg.FeminicideForm1EpsName,
		&feminicideForm1Pg.FeminicideForm1InformantAddress,
		&feminicideForm1Pg.FeminicideForm1InformantZone,
		&feminicideForm1Pg.FeminicideForm1InformantLivingTownCode,
		&feminicideForm1Pg.FeminicideForm1InformantPhone,
		&feminicideForm1Pg.FeminicideForm1EmergencyContactNames,
		&feminicideForm1Pg.FeminicideForm1EmergencyContactNumber,
		&feminicideForm1Pg.FeminicideForm1InformantSex,
		&feminicideForm1Pg.FeminicideForm1InformantGenderIdentity,
		&feminicideForm1Pg.FeminicideForm1InformantSexualOrientation,
		&feminicideForm1Pg.FeminicideForm1InformantEthnicity,
		&feminicideForm1Pg.FeminicideForm1InformantIndigenousPeople,
		&feminicideForm1Pg.FeminicideForm1InformantMigratoryStatus,
		&feminicideForm1Pg.FeminicideForm1InformantMigrationSituation,
		&feminicideForm1Pg.FeminicideForm1InformantHighestEducationLevel,
		&feminicideForm1Pg.FeminicideForm1InformantSpecialPopulation,
		&feminicideForm1Pg.FeminicideForm1InformantCurrentEmployment,
		&feminicideForm1Pg.FeminicideForm1InformantEmploymentAccess,
		&feminicideForm1Pg.FeminicideForm1InformantEmploymentImpactDescription,
		&feminicideForm1Pg.FeminicideForm1InformantPrimaryOccupation,
		&feminicideForm1Pg.FeminicideForm1InformantDisability,
		&feminicideForm1Pg.FeminicideForm1InformantDisabilityType,
		&feminicideForm1Pg.FeminicideForm1SituationAfterFeminicide,
		&feminicideForm1Pg.FeminicideForm1AdditionalInformation,
		&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceived,
		&feminicideForm1Pg.FeminicideForm1HouseholdExpenseResponsibility,
		&feminicideForm1Pg.FeminicideForm1PostDeathEconomicAssumption,
		&feminicideForm1Pg.FeminicideForm1EconomicAssumptionExplanation,
		&feminicideForm1Pg.FeminicideForm1AnyDependentPeople,
		&feminicideForm1Pg.FeminicideForm1AnyPublicOrPrivateEntity,
		&feminicideForm1Pg.FeminicideForm1AnyPublicOrPrivateEntityExplanation,
		&feminicideForm1Pg.FeminicideForm1PublicTransportAccess,
		&feminicideForm1Pg.FeminicideForm1PreferredTransportationMode,
		&feminicideForm1Pg.FeminicideForm1PreferredTransportationModeExplanation,
		&feminicideForm1Pg.FeminicideForm1TransportationCostEstimate,
		&feminicideForm1Pg.FeminicideForm1TransportDifficulty,
		&feminicideForm1Pg.FeminicideForm1TransportDifficultyExplanation,
		&feminicideForm1Pg.FeminicideForm1EconomicResourcesForTransport,
		&feminicideForm1Pg.FeminicideForm1DebtOrHelpDueToTransport,
		&feminicideForm1Pg.FeminicideForm1DebtImpactExplanation,
		&feminicideForm1Pg.FeminicideForm1TransportSubsidyReceived,
		&feminicideForm1Pg.FeminicideForm1TransportSubsidyExplanation,
		&feminicideForm1Pg.FeminicideForm1SafetyTransportationConcern,
		&feminicideForm1Pg.FeminicideForm1SafetyTransportationExplanation,
		&feminicideForm1Pg.FeminicideForm1FoodAccessFrequency,
		&feminicideForm1Pg.FeminicideForm1FoodAccessExplanation,
		&feminicideForm1Pg.FeminicideForm1AggressorFoodRestriction,
		&feminicideForm1Pg.FeminicideForm1AggressorFoodRestrictionExplanation,
		&feminicideForm1Pg.FeminicideForm1FoodIncomeSupport,
		&feminicideForm1Pg.FeminicideForm1FoodIncomeSupportExplanation,
		&feminicideForm1Pg.FeminicideForm1JuridicalAssistanceReceived,
		&feminicideForm1Pg.FeminicideForm1JuridicalAssistanceExplanation,
		&feminicideForm1Pg.FeminicideForm1VictimRepresentation,
		&feminicideForm1Pg.FeminicideForm1VictimRepresentationExplanation,
		&feminicideForm1Pg.FeminicideForm1PsychosocialSupportReceived,
		&feminicideForm1Pg.FeminicideForm1PsychosocialSupportExplanation,
		&feminicideForm1Pg.FeminicideForm1EmergencyEmotionalCrisis,
		&feminicideForm1Pg.FeminicideForm1EmergencyEmotionalCrisisExplanation,
		&feminicideForm1Pg.FeminicideForm1AggressorSameResidence,
		&feminicideForm1Pg.FeminicideForm1AggressorSameResidenceExplanation,
		&feminicideForm1Pg.FeminicideForm1AggressorLocationKnown,
		&feminicideForm1Pg.FeminicideForm1AggressorLocationKnownExplanation,
		&feminicideForm1Pg.FeminicideForm1AnyTypeOfAssistanceReceived,
		&feminicideForm1Pg.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation,
		&feminicideForm1Pg.FeminicideForm1CompensationFundInsuranceCoverage,
		&feminicideForm1Pg.FeminicideForm1CompensationFundInsuranceCoverageExplanation,
		&feminicideForm1Pg.FeminicideForm1FuneralSubsidyReceived,
		&feminicideForm1Pg.FeminicideForm1FuneralSubsidyExplanation,
		&feminicideForm1Pg.FeminicideForm1FuneralFundsAvailable,
		&feminicideForm1Pg.FeminicideForm1FuneralFundsAvailableExplanation,
		&feminicideForm1Pg.FeminicideForm1FuneralCostValue,
		&feminicideForm1Pg.FeminicideForm1RenameReputationImpact,
		&feminicideForm1Pg.FeminicideForm1RenameReputationExplanation,
		&feminicideForm1Pg.FeminicideForm1AdditionalNeedsDescription,
		&feminicideForm1Pg.FeminicideForm1ActionPlan,
		&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedCityHall,
		&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedWomensOffice,
		&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedOtherEntity,
		&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedOther,
		&feminicideForm1Pg.FeminicideForm1AssistanceReceived,
		&feminicideForm1Pg.FeminicideForm1AssistanceReceivedCityHall,
		&feminicideForm1Pg.FeminicideForm1AssistanceReceivedWomensOffice,
		&feminicideForm1Pg.FeminicideForm1AssistanceReceivedOtherEntity,
		&feminicideForm1Pg.FeminicideForm1AssistanceReceivedOther,
		&feminicideForm1Pg.FeminicideForm1FamilyFather,
		&feminicideForm1Pg.FeminicideForm1FamilyMother,
		&feminicideForm1Pg.FeminicideForm1FamilyStepfather,
		&feminicideForm1Pg.FeminicideForm1FamilyStepmother,
		&feminicideForm1Pg.FeminicideForm1FamilyPartner,
		&feminicideForm1Pg.FeminicideForm1FamilySibling1,
		&feminicideForm1Pg.FeminicideForm1FamilySibling2,
		&feminicideForm1Pg.FeminicideForm1FamilySibling3,
		&feminicideForm1Pg.FeminicideForm1FamilySibling4,
		&feminicideForm1Pg.FeminicideForm1FamilySibling5,
		&feminicideForm1Pg.FeminicideForm1FamilySonDaughter1,
		&feminicideForm1Pg.FeminicideForm1FamilySonDaughter2,
		&feminicideForm1Pg.FeminicideForm1FamilySonDaughter3,
		&feminicideForm1Pg.FeminicideForm1FamilySonDaughter4,
		&feminicideForm1Pg.FeminicideForm1FamilySonDaughter5,
		&feminicideForm1Pg.FeminicideForm1FamilyGrandmother,
		&feminicideForm1Pg.FeminicideForm1FamilyGrandfather,
		&feminicideForm1Pg.FeminicideForm1FamilyOtherMember,
		&feminicideForm1Pg.FeminicideForm1Summary,
		&feminicideForm1Pg.FeminicideForm1Feminicide)

	*feminicideForm1 = feminicideForm1Pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicideCases(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideForm1DTO, int, error) {
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicideForm1Path string = FeminicideForm1DBScheme + "." + FeminicideForm1DBName

	var feminicideForm1s []FeminicideForm1DTO

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	var feminicideForm1FieldsSlice []string = []string{
		"FeminicideForm1Id", "FeminicideForm1ICode", "FeminicideForm1CreationDate", "FeminicideForm1UpdateDate",
		"FeminicideForm1VictimIdentityName", "FeminicideForm1BirthDate", "FeminicideForm1DeathDate",
		"FeminicideForm1VictimAddress", "FeminicideForm1VictimZone", "FeminicideForm1VictimLivingTownCode",
		"FeminicideForm1VictimMaritalStatus", "FeminicideForm1VictimSex", "FeminicideForm1VictimGenderIdentity",
		"FeminicideForm1VictimSexualOrientation", "FeminicideForm1VictimEthnicity", "FeminicideForm1VictimIndigenousPeople", "FeminicideForm1VictimSpecialPopulation",
		"FeminicideForm1VictimDisability", "FeminicideForm1VictimDisabilityType", "FeminicideForm1PresumedAggressorNames",
		"FeminicideForm1PresumedAggressorRelation", "FeminicideForm1PresumedAggressorKnownVGB",
		"FeminicideForm1InformantNames", "FeminicideForm1InformantIdentityName", "FeminicideForm1InformantDocType",
		"FeminicideForm1InformantDocNumber", "FeminicideForm1InformantBirthDate", "FeminicideForm1SGSSSAffiliation",
		"FeminicideForm1EpsName", "FeminicideForm1InformantAddress", "FeminicideForm1InformantZone",
		"FeminicideForm1InformantLivingTownCode", "FeminicideForm1InformantPhone", "FeminicideForm1EmergencyContactNames",
		"FeminicideForm1EmergencyContactNumber", "FeminicideForm1InformantSex", "FeminicideForm1InformantGenderIdentity",
		"FeminicideForm1InformantSexualOrientation", "FeminicideForm1InformantEthnicity", "FeminicideForm1InformantIndigenousPeople", "FeminicideForm1InformantMigratoryStatus",
		"FeminicideForm1InformantMigrationSituation", "FeminicideForm1InformantHighestEducationLevel", "FeminicideForm1InformantSpecialPopulation", "FeminicideForm1InformantCurrentEmployment",
		"FeminicideForm1InformantEmploymentAccess", "FeminicideForm1InformantEmploymentImpactDescription", "FeminicideForm1InformantPrimaryOccupation",
		"FeminicideForm1InformantDisability", "FeminicideForm1InformantDisabilityType", "FeminicideForm1SituationAfterFeminicide",
		"FeminicideForm1AdditionalInformation", "FeminicideForm1AnyAssistanceReceived", "FeminicideForm1HouseholdExpenseResponsibility",
		"FeminicideForm1PostDeathEconomicAssumption", "FeminicideForm1EconomicAssumptionExplanation", "FeminicideForm1AnyDependentPeople",
		"FeminicideForm1AnyPublicOrPrivateEntity", "FeminicideForm1AnyPublicOrPrivateEntityExplanation", "FeminicideForm1PublicTransportAccess",
		"FeminicideForm1PreferredTransportationMode", "FeminicideForm1PreferredTransportationModeExplanation", "FeminicideForm1TransportationCostEstimate",
		"FeminicideForm1TransportDifficulty", "FeminicideForm1TransportDifficultyExplanation", "FeminicideForm1EconomicResourcesForTransport",
		"FeminicideForm1DebtOrHelpDueToTransport", "FeminicideForm1DebtImpactExplanation", "FeminicideForm1TransportSubsidyReceived",
		"FeminicideForm1TransportSubsidyExplanation", "FeminicideForm1SafetyTransportationConcern", "FeminicideForm1SafetyTransportationExplanation",
		"FeminicideForm1FoodAccessFrequency", "FeminicideForm1FoodAccessExplanation", "FeminicideForm1AggressorFoodRestriction",
		"FeminicideForm1AggressorFoodRestrictionExplanation", "FeminicideForm1FoodIncomeSupport", "FeminicideForm1FoodIncomeSupportExplanation",
		"FeminicideForm1JuridicalAssistanceReceived", "FeminicideForm1JuridicalAssistanceExplanation", "FeminicideForm1VictimRepresentation",
		"FeminicideForm1VictimRepresentationExplanation", "FeminicideForm1PsychosocialSupportReceived", "FeminicideForm1PsychosocialSupportExplanation",
		"FeminicideForm1EmergencyEmotionalCrisis", "FeminicideForm1EmergencyEmotionalCrisisExplanation", "FeminicideForm1AggressorSameResidence",
		"FeminicideForm1AggressorSameResidenceExplanation", "FeminicideForm1AggressorLocationKnown", "FeminicideForm1AggressorLocationKnownExplanation",
		"FeminicideForm1AnyTypeOfAssistanceReceived", "FeminicideForm1AnyTypeOfAssistanceReceivedExplanation", "FeminicideForm1CompensationFundInsuranceCoverage",
		"FeminicideForm1CompensationFundInsuranceCoverageExplanation", "FeminicideForm1FuneralSubsidyReceived", "FeminicideForm1FuneralSubsidyExplanation",
		"FeminicideForm1FuneralFundsAvailable", "FeminicideForm1FuneralFundsAvailableExplanation", "FeminicideForm1FuneralCostValue",
		"FeminicideForm1RenameReputationImpact", "FeminicideForm1RenameReputationExplanation", "FeminicideForm1AdditionalNeedsDescription",
		"FeminicideForm1ActionPlan", "FeminicideForm1AnyAssistanceReceivedCityHall", "FeminicideForm1AnyAssistanceReceivedWomensOffice", "FeminicideForm1AnyAssistanceReceivedOtherEntity", "FeminicideForm1AnyAssistanceReceivedOther",
		"FeminicideForm1AssistanceReceived", "FeminicideForm1AssistanceReceivedCityHall", "FeminicideForm1AssistanceReceivedWomensOffice", "FeminicideForm1AssistanceReceivedOtherEntity", "FeminicideForm1AssistanceReceivedOther",
		"FeminicideForm1FamilyMother", "FeminicideForm1FamilyFather", "FeminicideForm1FamilyStepfather", "FeminicideForm1FamilyStepmother",
		"FeminicideForm1FamilyPartner", "FeminicideForm1FamilySibling1", "FeminicideForm1FamilySibling2", "FeminicideForm1FamilySibling3",
		"FeminicideForm1FamilySibling4", "FeminicideForm1FamilySibling5", "FeminicideForm1FamilySonDaughter1", "FeminicideForm1FamilySonDaughter2",
		"FeminicideForm1FamilySonDaughter3", "FeminicideForm1FamilySonDaughter4", "FeminicideForm1FamilySonDaughter5", "FeminicideForm1FamilyGrandmother",
		"FeminicideForm1FamilyGrandfather", "FeminicideForm1FamilyOtherMember",
		"FeminicideForm1Summary", "FeminicideForm1Feminicide",
	}
	var feminicideForm1FieldsAliasSlice []string = []string{}

	var feminicideForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, feminicideForm1FieldsSlice, feminicideForm1FieldsAliasSlice, FeminicideForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, true)

	var query string = `SELECT ` + feminicideForm1FieldsStr +
		` FROM ` + feminicideForm1Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, true) +
		` ORDER BY ` + feminicideForm1Path + `.` + FeminicideForm1FieldDefinitions["FeminicideForm1CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	for persistenceCtrl.Next() {
		var feminicideForm1Pg FeminicideForm1PgDB = FeminicideForm1PgDB{}
		persistenceCtrl.ScanRow(
			&feminicideForm1Pg.FeminicideForm1Id,
			&feminicideForm1Pg.FeminicideForm1ICode,
			&feminicideForm1Pg.FeminicideForm1CreationDate,
			&feminicideForm1Pg.FeminicideForm1UpdateDate,
			&feminicideForm1Pg.FeminicideForm1VictimIdentityName,
			&feminicideForm1Pg.FeminicideForm1BirthDate,
			&feminicideForm1Pg.FeminicideForm1DeathDate,
			&feminicideForm1Pg.FeminicideForm1VictimAddress,
			&feminicideForm1Pg.FeminicideForm1VictimZone,
			&feminicideForm1Pg.FeminicideForm1VictimLivingTownCode,
			&feminicideForm1Pg.FeminicideForm1VictimMaritalStatus,
			&feminicideForm1Pg.FeminicideForm1VictimSex,
			&feminicideForm1Pg.FeminicideForm1VictimGenderIdentity,
			&feminicideForm1Pg.FeminicideForm1VictimSexualOrientation,
			&feminicideForm1Pg.FeminicideForm1VictimEthnicity,
			&feminicideForm1Pg.FeminicideForm1VictimIndigenousPeople,
			&feminicideForm1Pg.FeminicideForm1VictimSpecialPopulation,
			&feminicideForm1Pg.FeminicideForm1VictimDisability,
			&feminicideForm1Pg.FeminicideForm1VictimDisabilityType,
			&feminicideForm1Pg.FeminicideForm1PresumedAggressorNames,
			&feminicideForm1Pg.FeminicideForm1PresumedAggressorRelation,
			&feminicideForm1Pg.FeminicideForm1PresumedAggressorKnownVGB,
			&feminicideForm1Pg.FeminicideForm1InformantNames,
			&feminicideForm1Pg.FeminicideForm1InformantIdentityName,
			&feminicideForm1Pg.FeminicideForm1InformantDocType,
			&feminicideForm1Pg.FeminicideForm1InformantDocNumber,
			&feminicideForm1Pg.FeminicideForm1InformantBirthDate,
			&feminicideForm1Pg.FeminicideForm1SGSSSAffiliation,
			&feminicideForm1Pg.FeminicideForm1EpsName,
			&feminicideForm1Pg.FeminicideForm1InformantAddress,
			&feminicideForm1Pg.FeminicideForm1InformantZone,
			&feminicideForm1Pg.FeminicideForm1InformantLivingTownCode,
			&feminicideForm1Pg.FeminicideForm1InformantPhone,
			&feminicideForm1Pg.FeminicideForm1EmergencyContactNames,
			&feminicideForm1Pg.FeminicideForm1EmergencyContactNumber,
			&feminicideForm1Pg.FeminicideForm1InformantSex,
			&feminicideForm1Pg.FeminicideForm1InformantGenderIdentity,
			&feminicideForm1Pg.FeminicideForm1InformantSexualOrientation,
			&feminicideForm1Pg.FeminicideForm1InformantEthnicity,
			&feminicideForm1Pg.FeminicideForm1InformantIndigenousPeople,
			&feminicideForm1Pg.FeminicideForm1InformantMigratoryStatus,
			&feminicideForm1Pg.FeminicideForm1InformantMigrationSituation,
			&feminicideForm1Pg.FeminicideForm1InformantHighestEducationLevel,
			&feminicideForm1Pg.FeminicideForm1InformantSpecialPopulation,
			&feminicideForm1Pg.FeminicideForm1InformantCurrentEmployment,
			&feminicideForm1Pg.FeminicideForm1InformantEmploymentAccess,
			&feminicideForm1Pg.FeminicideForm1InformantEmploymentImpactDescription,
			&feminicideForm1Pg.FeminicideForm1InformantPrimaryOccupation,
			&feminicideForm1Pg.FeminicideForm1InformantDisability,
			&feminicideForm1Pg.FeminicideForm1InformantDisabilityType,
			&feminicideForm1Pg.FeminicideForm1SituationAfterFeminicide,
			&feminicideForm1Pg.FeminicideForm1AdditionalInformation,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1HouseholdExpenseResponsibility,
			&feminicideForm1Pg.FeminicideForm1PostDeathEconomicAssumption,
			&feminicideForm1Pg.FeminicideForm1EconomicAssumptionExplanation,
			&feminicideForm1Pg.FeminicideForm1AnyDependentPeople,
			&feminicideForm1Pg.FeminicideForm1AnyPublicOrPrivateEntity,
			&feminicideForm1Pg.FeminicideForm1AnyPublicOrPrivateEntityExplanation,
			&feminicideForm1Pg.FeminicideForm1PublicTransportAccess,
			&feminicideForm1Pg.FeminicideForm1PreferredTransportationMode,
			&feminicideForm1Pg.FeminicideForm1PreferredTransportationModeExplanation,
			&feminicideForm1Pg.FeminicideForm1TransportationCostEstimate,
			&feminicideForm1Pg.FeminicideForm1TransportDifficulty,
			&feminicideForm1Pg.FeminicideForm1TransportDifficultyExplanation,
			&feminicideForm1Pg.FeminicideForm1EconomicResourcesForTransport,
			&feminicideForm1Pg.FeminicideForm1DebtOrHelpDueToTransport,
			&feminicideForm1Pg.FeminicideForm1DebtImpactExplanation,
			&feminicideForm1Pg.FeminicideForm1TransportSubsidyReceived,
			&feminicideForm1Pg.FeminicideForm1TransportSubsidyExplanation,
			&feminicideForm1Pg.FeminicideForm1SafetyTransportationConcern,
			&feminicideForm1Pg.FeminicideForm1SafetyTransportationExplanation,
			&feminicideForm1Pg.FeminicideForm1FoodAccessFrequency,
			&feminicideForm1Pg.FeminicideForm1FoodAccessExplanation,
			&feminicideForm1Pg.FeminicideForm1AggressorFoodRestriction,
			&feminicideForm1Pg.FeminicideForm1AggressorFoodRestrictionExplanation,
			&feminicideForm1Pg.FeminicideForm1FoodIncomeSupport,
			&feminicideForm1Pg.FeminicideForm1FoodIncomeSupportExplanation,
			&feminicideForm1Pg.FeminicideForm1JuridicalAssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1JuridicalAssistanceExplanation,
			&feminicideForm1Pg.FeminicideForm1VictimRepresentation,
			&feminicideForm1Pg.FeminicideForm1VictimRepresentationExplanation,
			&feminicideForm1Pg.FeminicideForm1PsychosocialSupportReceived,
			&feminicideForm1Pg.FeminicideForm1PsychosocialSupportExplanation,
			&feminicideForm1Pg.FeminicideForm1EmergencyEmotionalCrisis,
			&feminicideForm1Pg.FeminicideForm1EmergencyEmotionalCrisisExplanation,
			&feminicideForm1Pg.FeminicideForm1AggressorSameResidence,
			&feminicideForm1Pg.FeminicideForm1AggressorSameResidenceExplanation,
			&feminicideForm1Pg.FeminicideForm1AggressorLocationKnown,
			&feminicideForm1Pg.FeminicideForm1AggressorLocationKnownExplanation,
			&feminicideForm1Pg.FeminicideForm1AnyTypeOfAssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation,
			&feminicideForm1Pg.FeminicideForm1CompensationFundInsuranceCoverage,
			&feminicideForm1Pg.FeminicideForm1CompensationFundInsuranceCoverageExplanation,
			&feminicideForm1Pg.FeminicideForm1FuneralSubsidyReceived,
			&feminicideForm1Pg.FeminicideForm1FuneralSubsidyExplanation,
			&feminicideForm1Pg.FeminicideForm1FuneralFundsAvailable,
			&feminicideForm1Pg.FeminicideForm1FuneralFundsAvailableExplanation,
			&feminicideForm1Pg.FeminicideForm1FuneralCostValue,
			&feminicideForm1Pg.FeminicideForm1RenameReputationImpact,
			&feminicideForm1Pg.FeminicideForm1RenameReputationExplanation,
			&feminicideForm1Pg.FeminicideForm1AdditionalNeedsDescription,
			&feminicideForm1Pg.FeminicideForm1ActionPlan,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedCityHall,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedWomensOffice,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedOtherEntity,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedOther,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedCityHall,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedWomensOffice,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedOtherEntity,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedOther,
			&feminicideForm1Pg.FeminicideForm1FamilyFather,
			&feminicideForm1Pg.FeminicideForm1FamilyMother,
			&feminicideForm1Pg.FeminicideForm1FamilyStepfather,
			&feminicideForm1Pg.FeminicideForm1FamilyStepmother,
			&feminicideForm1Pg.FeminicideForm1FamilyPartner,
			&feminicideForm1Pg.FeminicideForm1FamilySibling1,
			&feminicideForm1Pg.FeminicideForm1FamilySibling2,
			&feminicideForm1Pg.FeminicideForm1FamilySibling3,
			&feminicideForm1Pg.FeminicideForm1FamilySibling4,
			&feminicideForm1Pg.FeminicideForm1FamilySibling5,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter1,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter2,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter3,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter4,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter5,
			&feminicideForm1Pg.FeminicideForm1FamilyGrandmother,
			&feminicideForm1Pg.FeminicideForm1FamilyGrandfather,
			&feminicideForm1Pg.FeminicideForm1FamilyOtherMember,
			&feminicideForm1Pg.FeminicideForm1Summary,
			&feminicideForm1Pg.FeminicideForm1Feminicide)

		feminicideForm1s = append(feminicideForm1s, feminicideForm1Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + feminicideForm1Path +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, by.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return feminicideForm1s, count, nil
}

func GetAllFeminicideCases(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideForm1DTO, int, error) {
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicideForm1Path string = FeminicideForm1DBScheme + "." + FeminicideForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	var feminicideForm1FieldsSlice []string = []string{
		"FeminicideForm1Id", "FeminicideForm1ICode", "FeminicideForm1CreationDate", "FeminicideForm1UpdateDate",
		"FeminicideForm1VictimIdentityName", "FeminicideForm1BirthDate", "FeminicideForm1DeathDate",
		"FeminicideForm1VictimAddress", "FeminicideForm1VictimZone", "FeminicideForm1VictimLivingTownCode",
		"FeminicideForm1VictimMaritalStatus", "FeminicideForm1VictimSex", "FeminicideForm1VictimGenderIdentity",
		"FeminicideForm1VictimSexualOrientation", "FeminicideForm1VictimEthnicity", "FeminicideForm1VictimIndigenousPeople", "FeminicideForm1VictimSpecialPopulation",
		"FeminicideForm1VictimDisability", "FeminicideForm1VictimDisabilityType", "FeminicideForm1PresumedAggressorNames",
		"FeminicideForm1PresumedAggressorRelation", "FeminicideForm1PresumedAggressorKnownVGB",
		"FeminicideForm1InformantNames", "FeminicideForm1InformantIdentityName", "FeminicideForm1InformantDocType",
		"FeminicideForm1InformantDocNumber", "FeminicideForm1InformantBirthDate", "FeminicideForm1SGSSSAffiliation",
		"FeminicideForm1EpsName", "FeminicideForm1InformantAddress", "FeminicideForm1InformantZone",
		"FeminicideForm1InformantLivingTownCode", "FeminicideForm1InformantPhone", "FeminicideForm1EmergencyContactNames",
		"FeminicideForm1EmergencyContactNumber", "FeminicideForm1InformantSex", "FeminicideForm1InformantGenderIdentity",
		"FeminicideForm1InformantSexualOrientation", "FeminicideForm1InformantEthnicity", "FeminicideForm1InformantIndigenousPeople", "FeminicideForm1InformantMigratoryStatus",
		"FeminicideForm1InformantMigrationSituation", "FeminicideForm1InformantHighestEducationLevel", "FeminicideForm1InformantSpecialPopulation", "FeminicideForm1InformantCurrentEmployment",
		"FeminicideForm1InformantEmploymentAccess", "FeminicideForm1InformantEmploymentImpactDescription", "FeminicideForm1InformantPrimaryOccupation",
		"FeminicideForm1InformantDisability", "FeminicideForm1InformantDisabilityType", "FeminicideForm1SituationAfterFeminicide",
		"FeminicideForm1AdditionalInformation", "FeminicideForm1AnyAssistanceReceived", "FeminicideForm1HouseholdExpenseResponsibility",
		"FeminicideForm1PostDeathEconomicAssumption", "FeminicideForm1EconomicAssumptionExplanation", "FeminicideForm1AnyDependentPeople",
		"FeminicideForm1AnyPublicOrPrivateEntity", "FeminicideForm1AnyPublicOrPrivateEntityExplanation", "FeminicideForm1PublicTransportAccess",
		"FeminicideForm1PreferredTransportationMode", "FeminicideForm1PreferredTransportationModeExplanation", "FeminicideForm1TransportationCostEstimate",
		"FeminicideForm1TransportDifficulty", "FeminicideForm1TransportDifficultyExplanation", "FeminicideForm1EconomicResourcesForTransport",
		"FeminicideForm1DebtOrHelpDueToTransport", "FeminicideForm1DebtImpactExplanation", "FeminicideForm1TransportSubsidyReceived",
		"FeminicideForm1TransportSubsidyExplanation", "FeminicideForm1SafetyTransportationConcern", "FeminicideForm1SafetyTransportationExplanation",
		"FeminicideForm1FoodAccessFrequency", "FeminicideForm1FoodAccessExplanation", "FeminicideForm1AggressorFoodRestriction",
		"FeminicideForm1AggressorFoodRestrictionExplanation", "FeminicideForm1FoodIncomeSupport", "FeminicideForm1FoodIncomeSupportExplanation",
		"FeminicideForm1JuridicalAssistanceReceived", "FeminicideForm1JuridicalAssistanceExplanation", "FeminicideForm1VictimRepresentation",
		"FeminicideForm1VictimRepresentationExplanation", "FeminicideForm1PsychosocialSupportReceived", "FeminicideForm1PsychosocialSupportExplanation",
		"FeminicideForm1EmergencyEmotionalCrisis", "FeminicideForm1EmergencyEmotionalCrisisExplanation", "FeminicideForm1AggressorSameResidence",
		"FeminicideForm1AggressorSameResidenceExplanation", "FeminicideForm1AggressorLocationKnown", "FeminicideForm1AggressorLocationKnownExplanation",
		"FeminicideForm1AnyTypeOfAssistanceReceived", "FeminicideForm1AnyTypeOfAssistanceReceivedExplanation", "FeminicideForm1CompensationFundInsuranceCoverage",
		"FeminicideForm1CompensationFundInsuranceCoverageExplanation", "FeminicideForm1FuneralSubsidyReceived", "FeminicideForm1FuneralSubsidyExplanation",
		"FeminicideForm1FuneralFundsAvailable", "FeminicideForm1FuneralFundsAvailableExplanation", "FeminicideForm1FuneralCostValue",
		"FeminicideForm1RenameReputationImpact", "FeminicideForm1RenameReputationExplanation", "FeminicideForm1AdditionalNeedsDescription",
		"FeminicideForm1ActionPlan", "FeminicideForm1AnyAssistanceReceivedCityHall", "FeminicideForm1AnyAssistanceReceivedWomensOffice", "FeminicideForm1AnyAssistanceReceivedOtherEntity", "FeminicideForm1AnyAssistanceReceivedOther",
		"FeminicideForm1AssistanceReceived", "FeminicideForm1AssistanceReceivedCityHall", "FeminicideForm1AssistanceReceivedWomensOffice", "FeminicideForm1AssistanceReceivedOtherEntity", "FeminicideForm1AssistanceReceivedOther",
		"FeminicideForm1FamilyMother", "FeminicideForm1FamilyFather", "FeminicideForm1FamilyStepfather", "FeminicideForm1FamilyStepmother",
		"FeminicideForm1FamilyPartner", "FeminicideForm1FamilySibling1", "FeminicideForm1FamilySibling2", "FeminicideForm1FamilySibling3",
		"FeminicideForm1FamilySibling4", "FeminicideForm1FamilySibling5", "FeminicideForm1FamilySonDaughter1", "FeminicideForm1FamilySonDaughter2",
		"FeminicideForm1FamilySonDaughter3", "FeminicideForm1FamilySonDaughter4", "FeminicideForm1FamilySonDaughter5", "FeminicideForm1FamilyGrandmother",
		"FeminicideForm1FamilyGrandfather", "FeminicideForm1FamilyOtherMember",
		"FeminicideForm1Summary", "FeminicideForm1Feminicide",
	}
	var feminicideForm1FieldsAliasSlice []string = []string{}

	var feminicideForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, feminicideForm1FieldsSlice, feminicideForm1FieldsAliasSlice, FeminicideForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, true)

	var query string = `SELECT ` + feminicideForm1FieldsStr +
		` FROM ` + feminicideForm1Path +
		` ORDER BY ` + feminicideForm1Path + `.` + FeminicideForm1FieldDefinitions["FeminicideForm1CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query)
	var feminicideForm1s []FeminicideForm1DTO
	for persistenceCtrl.Next() {
		var feminicideForm1Pg FeminicideForm1PgDB = FeminicideForm1PgDB{}

		persistenceCtrl.ScanRow(
			&feminicideForm1Pg.FeminicideForm1Id,
			&feminicideForm1Pg.FeminicideForm1ICode,
			&feminicideForm1Pg.FeminicideForm1CreationDate,
			&feminicideForm1Pg.FeminicideForm1UpdateDate,
			&feminicideForm1Pg.FeminicideForm1VictimIdentityName,
			&feminicideForm1Pg.FeminicideForm1BirthDate,
			&feminicideForm1Pg.FeminicideForm1DeathDate,
			&feminicideForm1Pg.FeminicideForm1VictimAddress,
			&feminicideForm1Pg.FeminicideForm1VictimZone,
			&feminicideForm1Pg.FeminicideForm1VictimLivingTownCode,
			&feminicideForm1Pg.FeminicideForm1VictimMaritalStatus,
			&feminicideForm1Pg.FeminicideForm1VictimSex,
			&feminicideForm1Pg.FeminicideForm1VictimGenderIdentity,
			&feminicideForm1Pg.FeminicideForm1VictimSexualOrientation,
			&feminicideForm1Pg.FeminicideForm1VictimEthnicity,
			&feminicideForm1Pg.FeminicideForm1VictimIndigenousPeople,
			&feminicideForm1Pg.FeminicideForm1VictimSpecialPopulation,
			&feminicideForm1Pg.FeminicideForm1VictimDisability,
			&feminicideForm1Pg.FeminicideForm1VictimDisabilityType,
			&feminicideForm1Pg.FeminicideForm1PresumedAggressorNames,
			&feminicideForm1Pg.FeminicideForm1PresumedAggressorRelation,
			&feminicideForm1Pg.FeminicideForm1PresumedAggressorKnownVGB,
			&feminicideForm1Pg.FeminicideForm1InformantNames,
			&feminicideForm1Pg.FeminicideForm1InformantIdentityName,
			&feminicideForm1Pg.FeminicideForm1InformantDocType,
			&feminicideForm1Pg.FeminicideForm1InformantDocNumber,
			&feminicideForm1Pg.FeminicideForm1InformantBirthDate,
			&feminicideForm1Pg.FeminicideForm1SGSSSAffiliation,
			&feminicideForm1Pg.FeminicideForm1EpsName,
			&feminicideForm1Pg.FeminicideForm1InformantAddress,
			&feminicideForm1Pg.FeminicideForm1InformantZone,
			&feminicideForm1Pg.FeminicideForm1InformantLivingTownCode,
			&feminicideForm1Pg.FeminicideForm1InformantPhone,
			&feminicideForm1Pg.FeminicideForm1EmergencyContactNames,
			&feminicideForm1Pg.FeminicideForm1EmergencyContactNumber,
			&feminicideForm1Pg.FeminicideForm1InformantSex,
			&feminicideForm1Pg.FeminicideForm1InformantGenderIdentity,
			&feminicideForm1Pg.FeminicideForm1InformantSexualOrientation,
			&feminicideForm1Pg.FeminicideForm1InformantEthnicity,
			&feminicideForm1Pg.FeminicideForm1InformantIndigenousPeople,
			&feminicideForm1Pg.FeminicideForm1InformantMigratoryStatus,
			&feminicideForm1Pg.FeminicideForm1InformantMigrationSituation,
			&feminicideForm1Pg.FeminicideForm1InformantHighestEducationLevel,
			&feminicideForm1Pg.FeminicideForm1InformantSpecialPopulation,
			&feminicideForm1Pg.FeminicideForm1InformantCurrentEmployment,
			&feminicideForm1Pg.FeminicideForm1InformantEmploymentAccess,
			&feminicideForm1Pg.FeminicideForm1InformantEmploymentImpactDescription,
			&feminicideForm1Pg.FeminicideForm1InformantPrimaryOccupation,
			&feminicideForm1Pg.FeminicideForm1InformantDisability,
			&feminicideForm1Pg.FeminicideForm1InformantDisabilityType,
			&feminicideForm1Pg.FeminicideForm1SituationAfterFeminicide,
			&feminicideForm1Pg.FeminicideForm1AdditionalInformation,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1HouseholdExpenseResponsibility,
			&feminicideForm1Pg.FeminicideForm1PostDeathEconomicAssumption,
			&feminicideForm1Pg.FeminicideForm1EconomicAssumptionExplanation,
			&feminicideForm1Pg.FeminicideForm1AnyDependentPeople,
			&feminicideForm1Pg.FeminicideForm1AnyPublicOrPrivateEntity,
			&feminicideForm1Pg.FeminicideForm1AnyPublicOrPrivateEntityExplanation,
			&feminicideForm1Pg.FeminicideForm1PublicTransportAccess,
			&feminicideForm1Pg.FeminicideForm1PreferredTransportationMode,
			&feminicideForm1Pg.FeminicideForm1PreferredTransportationModeExplanation,
			&feminicideForm1Pg.FeminicideForm1TransportationCostEstimate,
			&feminicideForm1Pg.FeminicideForm1TransportDifficulty,
			&feminicideForm1Pg.FeminicideForm1TransportDifficultyExplanation,
			&feminicideForm1Pg.FeminicideForm1EconomicResourcesForTransport,
			&feminicideForm1Pg.FeminicideForm1DebtOrHelpDueToTransport,
			&feminicideForm1Pg.FeminicideForm1DebtImpactExplanation,
			&feminicideForm1Pg.FeminicideForm1TransportSubsidyReceived,
			&feminicideForm1Pg.FeminicideForm1TransportSubsidyExplanation,
			&feminicideForm1Pg.FeminicideForm1SafetyTransportationConcern,
			&feminicideForm1Pg.FeminicideForm1SafetyTransportationExplanation,
			&feminicideForm1Pg.FeminicideForm1FoodAccessFrequency,
			&feminicideForm1Pg.FeminicideForm1FoodAccessExplanation,
			&feminicideForm1Pg.FeminicideForm1AggressorFoodRestriction,
			&feminicideForm1Pg.FeminicideForm1AggressorFoodRestrictionExplanation,
			&feminicideForm1Pg.FeminicideForm1FoodIncomeSupport,
			&feminicideForm1Pg.FeminicideForm1FoodIncomeSupportExplanation,
			&feminicideForm1Pg.FeminicideForm1JuridicalAssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1JuridicalAssistanceExplanation,
			&feminicideForm1Pg.FeminicideForm1VictimRepresentation,
			&feminicideForm1Pg.FeminicideForm1VictimRepresentationExplanation,
			&feminicideForm1Pg.FeminicideForm1PsychosocialSupportReceived,
			&feminicideForm1Pg.FeminicideForm1PsychosocialSupportExplanation,
			&feminicideForm1Pg.FeminicideForm1EmergencyEmotionalCrisis,
			&feminicideForm1Pg.FeminicideForm1EmergencyEmotionalCrisisExplanation,
			&feminicideForm1Pg.FeminicideForm1AggressorSameResidence,
			&feminicideForm1Pg.FeminicideForm1AggressorSameResidenceExplanation,
			&feminicideForm1Pg.FeminicideForm1AggressorLocationKnown,
			&feminicideForm1Pg.FeminicideForm1AggressorLocationKnownExplanation,
			&feminicideForm1Pg.FeminicideForm1AnyTypeOfAssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation,
			&feminicideForm1Pg.FeminicideForm1CompensationFundInsuranceCoverage,
			&feminicideForm1Pg.FeminicideForm1CompensationFundInsuranceCoverageExplanation,
			&feminicideForm1Pg.FeminicideForm1FuneralSubsidyReceived,
			&feminicideForm1Pg.FeminicideForm1FuneralSubsidyExplanation,
			&feminicideForm1Pg.FeminicideForm1FuneralFundsAvailable,
			&feminicideForm1Pg.FeminicideForm1FuneralFundsAvailableExplanation,
			&feminicideForm1Pg.FeminicideForm1FuneralCostValue,
			&feminicideForm1Pg.FeminicideForm1RenameReputationImpact,
			&feminicideForm1Pg.FeminicideForm1RenameReputationExplanation,
			&feminicideForm1Pg.FeminicideForm1AdditionalNeedsDescription,
			&feminicideForm1Pg.FeminicideForm1ActionPlan,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedCityHall,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedWomensOffice,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedOtherEntity,
			&feminicideForm1Pg.FeminicideForm1AnyAssistanceReceivedOther,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceived,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedCityHall,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedWomensOffice,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedOtherEntity,
			&feminicideForm1Pg.FeminicideForm1AssistanceReceivedOther,
			&feminicideForm1Pg.FeminicideForm1FamilyFather,
			&feminicideForm1Pg.FeminicideForm1FamilyMother,
			&feminicideForm1Pg.FeminicideForm1FamilyStepfather,
			&feminicideForm1Pg.FeminicideForm1FamilyStepmother,
			&feminicideForm1Pg.FeminicideForm1FamilyPartner,
			&feminicideForm1Pg.FeminicideForm1FamilySibling1,
			&feminicideForm1Pg.FeminicideForm1FamilySibling2,
			&feminicideForm1Pg.FeminicideForm1FamilySibling3,
			&feminicideForm1Pg.FeminicideForm1FamilySibling4,
			&feminicideForm1Pg.FeminicideForm1FamilySibling5,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter1,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter2,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter3,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter4,
			&feminicideForm1Pg.FeminicideForm1FamilySonDaughter5,
			&feminicideForm1Pg.FeminicideForm1FamilyGrandmother,
			&feminicideForm1Pg.FeminicideForm1FamilyGrandfather,
			&feminicideForm1Pg.FeminicideForm1FamilyOtherMember,
			&feminicideForm1Pg.FeminicideForm1Summary,
			&feminicideForm1Pg.FeminicideForm1Feminicide)

		feminicideForm1s = append(feminicideForm1s, feminicideForm1Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + feminicideForm1Path

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return feminicideForm1s, count, nil
}

func UpdateFeminicideForm1(feminicideForm1Form *FeminicideForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var feminicideForm1FormFieldsSlice []string = []string{
		"FeminicideForm1UpdateDate",
		"FeminicideForm1VictimIdentityName",
		"FeminicideForm1BirthDate",
		"FeminicideForm1DeathDate",
		"FeminicideForm1VictimAddress",
		"FeminicideForm1VictimZone",
		"FeminicideForm1VictimLivingTownCode",
		"FeminicideForm1VictimMaritalStatus",
		"FeminicideForm1VictimSex",
		"FeminicideForm1VictimGenderIdentity",
		"FeminicideForm1VictimSexualOrientation",
		"FeminicideForm1VictimEthnicity",
		"FeminicideForm1VictimIndigenousPeople",
		"FeminicideForm1VictimSpecialPopulation",
		"FeminicideForm1VictimDisability",
		"FeminicideForm1VictimDisabilityType",
		"FeminicideForm1PresumedAggressorNames",
		"FeminicideForm1PresumedAggressorRelation",
		"FeminicideForm1PresumedAggressorKnownVGB",
		"FeminicideForm1InformantNames",
		"FeminicideForm1InformantIdentityName",
		"FeminicideForm1InformantDocType",
		"FeminicideForm1InformantDocNumber",
		"FeminicideForm1InformantBirthDate",
		"FeminicideForm1SGSSSAffiliation",
		"FeminicideForm1EpsName",
		"FeminicideForm1InformantAddress",
		"FeminicideForm1InformantZone",
		"FeminicideForm1InformantLivingTownCode",
		"FeminicideForm1InformantPhone",
		"FeminicideForm1EmergencyContactNames",
		"FeminicideForm1EmergencyContactNumber",
		"FeminicideForm1InformantSex",
		"FeminicideForm1InformantGenderIdentity",
		"FeminicideForm1InformantSexualOrientation",
		"FeminicideForm1InformantEthnicity",
		"FeminicideForm1InformantIndigenousPeople",
		"FeminicideForm1InformantMigratoryStatus",
		"FeminicideForm1InformantMigrationSituation",
		"FeminicideForm1InformantHighestEducationLevel",
		"FeminicideForm1InformantSpecialPopulation",
		"FeminicideForm1InformantCurrentEmployment",
		"FeminicideForm1InformantEmploymentAccess",
		"FeminicideForm1InformantEmploymentImpactDescription",
		"FeminicideForm1InformantPrimaryOccupation",
		"FeminicideForm1InformantDisability",
		"FeminicideForm1InformantDisabilityType",
		"FeminicideForm1SituationAfterFeminicide",
		"FeminicideForm1AdditionalInformation",
		"FeminicideForm1AnyAssistanceReceived",
		"FeminicideForm1HouseholdExpenseResponsibility",
		"FeminicideForm1PostDeathEconomicAssumption",
		"FeminicideForm1EconomicAssumptionExplanation",
		"FeminicideForm1AnyDependentPeople",
		"FeminicideForm1AnyPublicOrPrivateEntity",
		"FeminicideForm1AnyPublicOrPrivateEntityExplanation",
		"FeminicideForm1PublicTransportAccess",
		"FeminicideForm1PreferredTransportationMode",
		"FeminicideForm1PreferredTransportationModeExplanation",
		"FeminicideForm1TransportationCostEstimate",
		"FeminicideForm1TransportDifficulty",
		"FeminicideForm1TransportDifficultyExplanation",
		"FeminicideForm1EconomicResourcesForTransport",
		"FeminicideForm1DebtOrHelpDueToTransport",
		"FeminicideForm1DebtImpactExplanation",
		"FeminicideForm1TransportSubsidyReceived",
		"FeminicideForm1TransportSubsidyExplanation",
		"FeminicideForm1SafetyTransportationConcern",
		"FeminicideForm1SafetyTransportationExplanation",
		"FeminicideForm1FoodAccessFrequency",
		"FeminicideForm1FoodAccessExplanation",
		"FeminicideForm1AggressorFoodRestriction",
		"FeminicideForm1AggressorFoodRestrictionExplanation",
		"FeminicideForm1FoodIncomeSupport",
		"FeminicideForm1FoodIncomeSupportExplanation",
		"FeminicideForm1JuridicalAssistanceReceived",
		"FeminicideForm1JuridicalAssistanceExplanation",
		"FeminicideForm1VictimRepresentation",
		"FeminicideForm1VictimRepresentationExplanation",
		"FeminicideForm1PsychosocialSupportReceived",
		"FeminicideForm1PsychosocialSupportExplanation",
		"FeminicideForm1EmergencyEmotionalCrisis",
		"FeminicideForm1EmergencyEmotionalCrisisExplanation",
		"FeminicideForm1AggressorSameResidence",
		"FeminicideForm1AggressorSameResidenceExplanation",
		"FeminicideForm1AggressorLocationKnown",
		"FeminicideForm1AggressorLocationKnownExplanation",
		"FeminicideForm1AnyTypeOfAssistanceReceived",
		"FeminicideForm1AnyTypeOfAssistanceReceivedExplanation",
		"FeminicideForm1CompensationFundInsuranceCoverage",
		"FeminicideForm1CompensationFundInsuranceCoverageExplanation",
		"FeminicideForm1FuneralSubsidyReceived",
		"FeminicideForm1FuneralSubsidyExplanation",
		"FeminicideForm1FuneralFundsAvailable",
		"FeminicideForm1FuneralFundsAvailableExplanation",
		"FeminicideForm1FuneralCostValue",
		"FeminicideForm1RenameReputationImpact",
		"FeminicideForm1RenameReputationExplanation",
		"FeminicideForm1AdditionalNeedsDescription",
		"FeminicideForm1ActionPlan",
		"FeminicideForm1AnyAssistanceReceivedCityHall",
		"FeminicideForm1AnyAssistanceReceivedWomensOffice",
		"FeminicideForm1AnyAssistanceReceivedOtherEntity",
		"FeminicideForm1AnyAssistanceReceivedOther",
		"FeminicideForm1AssistanceReceived",
		"FeminicideForm1AssistanceReceivedCityHall",
		"FeminicideForm1AssistanceReceivedWomensOffice",
		"FeminicideForm1AssistanceReceivedOtherEntity",

		"FeminicideForm1FamilyMother",
		"FeminicideForm1FamilyFather",
		"FeminicideForm1FamilyStepfather",
		"FeminicideForm1FamilyStepmother",
		"FeminicideForm1FamilyPartner",
		"FeminicideForm1FamilySibling1",
		"FeminicideForm1FamilySibling2",
		"FeminicideForm1FamilySibling3",
		"FeminicideForm1FamilySibling4",
		"FeminicideForm1FamilySibling5",
		"FeminicideForm1FamilySonDaughter1",
		"FeminicideForm1FamilySonDaughter2",
		"FeminicideForm1FamilySonDaughter3",
		"FeminicideForm1FamilySonDaughter4",
		"FeminicideForm1FamilySonDaughter5",
		"FeminicideForm1FamilyGrandmother",
		"FeminicideForm1FamilyGrandfather",
		"FeminicideForm1FamilyOtherMember",
		"FeminicideForm1AssistanceReceivedOther",

		"FeminicideForm1Summary",
	}
	var feminicideForm1FormFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, feminicideForm1FormFieldsSlice, feminicideForm1FormFieldsAliasSlice, FeminicideForm1DBName, []string{"FeminicideForm1Id"}, []string{}, []string{}, common_dao.SQL_AND, FeminicideForm1DBScheme, FeminicideForm1FieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, feminicideForm1Form.FeminicideForm1Id,
		feminicideForm1Form.FeminicideForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideForm1Form.FeminicideForm1VictimIdentityName,
		feminicideForm1Form.FeminicideForm1BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideForm1Form.FeminicideForm1DeathDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideForm1Form.FeminicideForm1VictimAddress,
		feminicideForm1Form.FeminicideForm1VictimZone.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimLivingTownCode,
		feminicideForm1Form.FeminicideForm1VictimMaritalStatus.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimSex.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimGenderIdentity.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimSexualOrientation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimEthnicity.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimIndigenousPeople.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimSpecialPopulation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimDisability.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimDisabilityType.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1PresumedAggressorNames,
		feminicideForm1Form.FeminicideForm1PresumedAggressorRelation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1PresumedAggressorKnownVGB.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantNames,
		feminicideForm1Form.FeminicideForm1InformantIdentityName,
		feminicideForm1Form.FeminicideForm1InformantDocType,
		feminicideForm1Form.FeminicideForm1InformantDocNumber,
		feminicideForm1Form.FeminicideForm1InformantBirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideForm1Form.FeminicideForm1SGSSSAffiliation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1EpsName,
		feminicideForm1Form.FeminicideForm1InformantAddress,
		feminicideForm1Form.FeminicideForm1InformantZone.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantLivingTownCode,
		feminicideForm1Form.FeminicideForm1InformantPhone,
		feminicideForm1Form.FeminicideForm1EmergencyContactNames,
		feminicideForm1Form.FeminicideForm1EmergencyContactNumber,
		feminicideForm1Form.FeminicideForm1InformantSex.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantGenderIdentity.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantSexualOrientation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantEthnicity.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantIndigenousPeople.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantMigratoryStatus.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantMigrationSituation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantHighestEducationLevel.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantSpecialPopulation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantCurrentEmployment.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantEmploymentAccess.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantEmploymentImpactDescription,
		feminicideForm1Form.FeminicideForm1InformantPrimaryOccupation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantDisability.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1InformantDisabilityType.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1SituationAfterFeminicide,
		feminicideForm1Form.FeminicideForm1AdditionalInformation,
		feminicideForm1Form.FeminicideForm1AnyAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1HouseholdExpenseResponsibility.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1PostDeathEconomicAssumption.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1EconomicAssumptionExplanation,
		feminicideForm1Form.FeminicideForm1AnyDependentPeople.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyPublicOrPrivateEntity.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyPublicOrPrivateEntityExplanation,
		feminicideForm1Form.FeminicideForm1PublicTransportAccess.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1PreferredTransportationMode.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1PreferredTransportationModeExplanation,
		feminicideForm1Form.FeminicideForm1TransportationCostEstimate,
		feminicideForm1Form.FeminicideForm1TransportDifficulty.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1TransportDifficultyExplanation,
		feminicideForm1Form.FeminicideForm1EconomicResourcesForTransport.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1DebtOrHelpDueToTransport.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1DebtImpactExplanation,
		feminicideForm1Form.FeminicideForm1TransportSubsidyReceived.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1TransportSubsidyExplanation,
		feminicideForm1Form.FeminicideForm1SafetyTransportationConcern.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1SafetyTransportationExplanation,
		feminicideForm1Form.FeminicideForm1FoodAccessFrequency.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1FoodAccessExplanation,
		feminicideForm1Form.FeminicideForm1AggressorFoodRestriction.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AggressorFoodRestrictionExplanation,
		feminicideForm1Form.FeminicideForm1FoodIncomeSupport.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1FoodIncomeSupportExplanation,
		feminicideForm1Form.FeminicideForm1JuridicalAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1JuridicalAssistanceExplanation,
		feminicideForm1Form.FeminicideForm1VictimRepresentation.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1VictimRepresentationExplanation,
		feminicideForm1Form.FeminicideForm1PsychosocialSupportReceived.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1PsychosocialSupportExplanation,
		feminicideForm1Form.FeminicideForm1EmergencyEmotionalCrisis.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1EmergencyEmotionalCrisisExplanation,
		feminicideForm1Form.FeminicideForm1AggressorSameResidence.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AggressorSameResidenceExplanation,
		feminicideForm1Form.FeminicideForm1AggressorLocationKnown.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AggressorLocationKnownExplanation,
		feminicideForm1Form.FeminicideForm1AnyTypeOfAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation,
		feminicideForm1Form.FeminicideForm1CompensationFundInsuranceCoverage.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1CompensationFundInsuranceCoverageExplanation,
		feminicideForm1Form.FeminicideForm1FuneralSubsidyReceived.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1FuneralSubsidyExplanation,
		feminicideForm1Form.FeminicideForm1FuneralFundsAvailable.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1FuneralFundsAvailableExplanation,
		feminicideForm1Form.FeminicideForm1FuneralCostValue.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1RenameReputationImpact.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1RenameReputationExplanation,
		feminicideForm1Form.FeminicideForm1AdditionalNeedsDescription,
		feminicideForm1Form.FeminicideForm1ActionPlan.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyAssistanceReceivedCityHall.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyAssistanceReceivedWomensOffice.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyAssistanceReceivedOtherEntity.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AnyAssistanceReceivedOther.VictimCaseForm2EnumsId,
		feminicideForm1Form.FeminicideForm1AssistanceReceived,
		feminicideForm1Form.FeminicideForm1AssistanceReceivedCityHall,
		feminicideForm1Form.FeminicideForm1AssistanceReceivedWomensOffice,
		feminicideForm1Form.FeminicideForm1AssistanceReceivedOtherEntity,
		feminicideForm1Form.FeminicideForm1AssistanceReceivedOther,
		feminicideForm1Form.FeminicideForm1FamilyFather,
		feminicideForm1Form.FeminicideForm1FamilyMother,
		feminicideForm1Form.FeminicideForm1FamilyStepfather,
		feminicideForm1Form.FeminicideForm1FamilyStepmother,
		feminicideForm1Form.FeminicideForm1FamilyPartner,
		feminicideForm1Form.FeminicideForm1FamilySibling1,
		feminicideForm1Form.FeminicideForm1FamilySibling2,
		feminicideForm1Form.FeminicideForm1FamilySibling3,
		feminicideForm1Form.FeminicideForm1FamilySibling4,
		feminicideForm1Form.FeminicideForm1FamilySibling5,
		feminicideForm1Form.FeminicideForm1FamilySonDaughter1,
		feminicideForm1Form.FeminicideForm1FamilySonDaughter2,
		feminicideForm1Form.FeminicideForm1FamilySonDaughter3,
		feminicideForm1Form.FeminicideForm1FamilySonDaughter4,
		feminicideForm1Form.FeminicideForm1FamilySonDaughter5,
		feminicideForm1Form.FeminicideForm1FamilyGrandmother,
		feminicideForm1Form.FeminicideForm1FamilyGrandfather,
		feminicideForm1Form.FeminicideForm1FamilyOtherMember,
		feminicideForm1Form.FeminicideForm1Summary)

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

func SetFeminicideForm1Defaults(feminicideForm1 *FeminicideForm1DTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		feminicideForm1.FeminicideForm1CreationDate = time.Now()
		feminicideForm1.FeminicideForm1UpdateDate = time.Now()
		feminicideForm1.FeminicideForm1ICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		feminicideForm1.FeminicideForm1UpdateDate = time.Now()
	}

}

func (obj *FeminicideForm1PgDB) ToDTO() FeminicideForm1DTO {
	var dto FeminicideForm1DTO

	if obj.FeminicideForm1Id.Valid {
		dto.FeminicideForm1Id = uint64(obj.FeminicideForm1Id.Int64)
	}

	if obj.FeminicideForm1ICode.Valid {
		dto.FeminicideForm1ICode = obj.FeminicideForm1ICode.String
	}

	if obj.FeminicideForm1CreationDate.Valid {
		dto.FeminicideForm1CreationDate = obj.FeminicideForm1CreationDate.Time
	}

	if obj.FeminicideForm1UpdateDate.Valid {
		dto.FeminicideForm1UpdateDate = obj.FeminicideForm1UpdateDate.Time
	}

	if obj.FeminicideForm1VictimIdentityName.Valid {
		dto.FeminicideForm1VictimIdentityName = obj.FeminicideForm1VictimIdentityName.String
	}

	if obj.FeminicideForm1BirthDate.Valid {
		dto.FeminicideForm1BirthDate = obj.FeminicideForm1BirthDate.Time
	}

	if obj.FeminicideForm1DeathDate.Valid {
		dto.FeminicideForm1DeathDate = obj.FeminicideForm1DeathDate.Time
	}

	if obj.FeminicideForm1VictimAddress.Valid {
		dto.FeminicideForm1VictimAddress = obj.FeminicideForm1VictimAddress.String
	}

	if obj.FeminicideForm1VictimZone.Valid {
		dto.FeminicideForm1VictimZone = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimZone.Int64)}
	}

	if obj.FeminicideForm1VictimLivingTownCode.Valid {
		dto.FeminicideForm1VictimLivingTownCode = obj.FeminicideForm1VictimLivingTownCode.String
	}

	if obj.FeminicideForm1VictimMaritalStatus.Valid {
		dto.FeminicideForm1VictimMaritalStatus = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimMaritalStatus.Int64)}
	}

	if obj.FeminicideForm1VictimSex.Valid {
		dto.FeminicideForm1VictimSex = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimSex.Int64)}
	}

	if obj.FeminicideForm1VictimGenderIdentity.Valid {
		dto.FeminicideForm1VictimGenderIdentity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimGenderIdentity.Int64)}
	}

	if obj.FeminicideForm1VictimSexualOrientation.Valid {
		dto.FeminicideForm1VictimSexualOrientation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimSexualOrientation.Int64)}
	}

	if obj.FeminicideForm1VictimEthnicity.Valid {
		dto.FeminicideForm1VictimEthnicity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimEthnicity.Int64)}
	}

	if obj.FeminicideForm1VictimIndigenousPeople.Valid {
		dto.FeminicideForm1VictimIndigenousPeople = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimIndigenousPeople.Int64)}
	}

	if obj.FeminicideForm1VictimSpecialPopulation.Valid {
		dto.FeminicideForm1VictimSpecialPopulation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimSpecialPopulation.Int64)}
	}

	if obj.FeminicideForm1VictimDisability.Valid {
		dto.FeminicideForm1VictimDisability = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimDisability.Int64)}
	}

	if obj.FeminicideForm1VictimDisabilityType.Valid {
		dto.FeminicideForm1VictimDisabilityType = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimDisabilityType.Int64)}
	}

	if obj.FeminicideForm1PresumedAggressorNames.Valid {
		dto.FeminicideForm1PresumedAggressorNames = obj.FeminicideForm1PresumedAggressorNames.String
	}

	if obj.FeminicideForm1PresumedAggressorRelation.Valid {
		dto.FeminicideForm1PresumedAggressorRelation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1PresumedAggressorRelation.Int64)}
	}

	if obj.FeminicideForm1PresumedAggressorKnownVGB.Valid {
		dto.FeminicideForm1PresumedAggressorKnownVGB = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1PresumedAggressorKnownVGB.Int64)}
	}

	if obj.FeminicideForm1InformantNames.Valid {
		dto.FeminicideForm1InformantNames = obj.FeminicideForm1InformantNames.String
	}

	if obj.FeminicideForm1InformantIdentityName.Valid {
		dto.FeminicideForm1InformantIdentityName = obj.FeminicideForm1InformantIdentityName.String
	}

	if obj.FeminicideForm1InformantDocType.Valid {
		dto.FeminicideForm1InformantDocType = obj.FeminicideForm1InformantDocType.String
	}

	if obj.FeminicideForm1InformantDocNumber.Valid {
		dto.FeminicideForm1InformantDocNumber = obj.FeminicideForm1InformantDocNumber.String
	}

	if obj.FeminicideForm1InformantBirthDate.Valid {
		dto.FeminicideForm1InformantBirthDate = obj.FeminicideForm1InformantBirthDate.Time
	}

	if obj.FeminicideForm1SGSSSAffiliation.Valid {
		dto.FeminicideForm1SGSSSAffiliation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1SGSSSAffiliation.Int64)}
	}

	if obj.FeminicideForm1EpsName.Valid {
		dto.FeminicideForm1EpsName = obj.FeminicideForm1EpsName.String
	}

	if obj.FeminicideForm1InformantAddress.Valid {
		dto.FeminicideForm1InformantAddress = obj.FeminicideForm1InformantAddress.String
	}

	if obj.FeminicideForm1InformantZone.Valid {
		dto.FeminicideForm1InformantZone = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantZone.Int64)}
	}

	if obj.FeminicideForm1InformantLivingTownCode.Valid {
		dto.FeminicideForm1InformantLivingTownCode = obj.FeminicideForm1InformantLivingTownCode.String
	}

	if obj.FeminicideForm1InformantPhone.Valid {
		dto.FeminicideForm1InformantPhone = obj.FeminicideForm1InformantPhone.String
	}

	if obj.FeminicideForm1EmergencyContactNames.Valid {
		dto.FeminicideForm1EmergencyContactNames = obj.FeminicideForm1EmergencyContactNames.String
	}

	if obj.FeminicideForm1EmergencyContactNumber.Valid {
		dto.FeminicideForm1EmergencyContactNumber = obj.FeminicideForm1EmergencyContactNumber.String
	}

	if obj.FeminicideForm1InformantSex.Valid {
		dto.FeminicideForm1InformantSex = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantSex.Int64)}
	}

	if obj.FeminicideForm1InformantGenderIdentity.Valid {
		dto.FeminicideForm1InformantGenderIdentity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantGenderIdentity.Int64)}
	}

	if obj.FeminicideForm1InformantSexualOrientation.Valid {
		dto.FeminicideForm1InformantSexualOrientation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantSexualOrientation.Int64)}
	}

	if obj.FeminicideForm1InformantEthnicity.Valid {
		dto.FeminicideForm1InformantEthnicity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantEthnicity.Int64)}
	}

	if obj.FeminicideForm1InformantIndigenousPeople.Valid {
		dto.FeminicideForm1InformantIndigenousPeople = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantIndigenousPeople.Int64)}
	}

	if obj.FeminicideForm1InformantMigratoryStatus.Valid {
		dto.FeminicideForm1InformantMigratoryStatus = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantMigratoryStatus.Int64)}
	}

	if obj.FeminicideForm1InformantMigrationSituation.Valid {
		dto.FeminicideForm1InformantMigrationSituation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantMigrationSituation.Int64)}
	}

	if obj.FeminicideForm1InformantHighestEducationLevel.Valid {
		dto.FeminicideForm1InformantHighestEducationLevel = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantHighestEducationLevel.Int64)}
	}

	if obj.FeminicideForm1InformantSpecialPopulation.Valid {
		dto.FeminicideForm1InformantSpecialPopulation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantSpecialPopulation.Int64)}
	}

	if obj.FeminicideForm1InformantCurrentEmployment.Valid {
		dto.FeminicideForm1InformantCurrentEmployment = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantCurrentEmployment.Int64)}
	}

	if obj.FeminicideForm1InformantEmploymentAccess.Valid {
		dto.FeminicideForm1InformantEmploymentAccess = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantEmploymentAccess.Int64)}
	}

	if obj.FeminicideForm1InformantEmploymentImpactDescription.Valid {
		dto.FeminicideForm1InformantEmploymentImpactDescription = obj.FeminicideForm1InformantEmploymentImpactDescription.String
	}

	if obj.FeminicideForm1InformantPrimaryOccupation.Valid {
		dto.FeminicideForm1InformantPrimaryOccupation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantPrimaryOccupation.Int64)}
	}

	if obj.FeminicideForm1InformantDisability.Valid {
		dto.FeminicideForm1InformantDisability = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantDisability.Int64)}
	}

	if obj.FeminicideForm1InformantDisabilityType.Valid {
		dto.FeminicideForm1InformantDisabilityType = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1InformantDisabilityType.Int64)}
	}

	if obj.FeminicideForm1SituationAfterFeminicide.Valid {
		dto.FeminicideForm1SituationAfterFeminicide = obj.FeminicideForm1SituationAfterFeminicide.String
	}

	if obj.FeminicideForm1AdditionalInformation.Valid {
		dto.FeminicideForm1AdditionalInformation = obj.FeminicideForm1AdditionalInformation.String
	}

	if obj.FeminicideForm1AnyAssistanceReceived.Valid {
		dto.FeminicideForm1AnyAssistanceReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyAssistanceReceived.Int64)}
	}

	if obj.FeminicideForm1HouseholdExpenseResponsibility.Valid {
		dto.FeminicideForm1HouseholdExpenseResponsibility = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1HouseholdExpenseResponsibility.Int64)}
	}

	if obj.FeminicideForm1PostDeathEconomicAssumption.Valid {
		dto.FeminicideForm1PostDeathEconomicAssumption = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1PostDeathEconomicAssumption.Int64)}
	}

	if obj.FeminicideForm1EconomicAssumptionExplanation.Valid {
		dto.FeminicideForm1EconomicAssumptionExplanation = obj.FeminicideForm1EconomicAssumptionExplanation.String
	}

	if obj.FeminicideForm1AnyDependentPeople.Valid {
		dto.FeminicideForm1AnyDependentPeople = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyDependentPeople.Int64)}
	}

	if obj.FeminicideForm1AnyPublicOrPrivateEntity.Valid {
		dto.FeminicideForm1AnyPublicOrPrivateEntity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyPublicOrPrivateEntity.Int64)}
	}

	if obj.FeminicideForm1AnyPublicOrPrivateEntityExplanation.Valid {
		dto.FeminicideForm1AnyPublicOrPrivateEntityExplanation = obj.FeminicideForm1AnyPublicOrPrivateEntityExplanation.String
	}

	if obj.FeminicideForm1PublicTransportAccess.Valid {
		dto.FeminicideForm1PublicTransportAccess = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1PublicTransportAccess.Int64)}
	}

	if obj.FeminicideForm1PreferredTransportationMode.Valid {
		dto.FeminicideForm1PreferredTransportationMode = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1PreferredTransportationMode.Int64)}
	}

	if obj.FeminicideForm1PreferredTransportationModeExplanation.Valid {
		dto.FeminicideForm1PreferredTransportationModeExplanation = obj.FeminicideForm1PreferredTransportationModeExplanation.String
	}

	if obj.FeminicideForm1TransportationCostEstimate.Valid {
		dto.FeminicideForm1TransportationCostEstimate = obj.FeminicideForm1TransportationCostEstimate.String
	}

	if obj.FeminicideForm1TransportDifficulty.Valid {
		dto.FeminicideForm1TransportDifficulty = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1TransportDifficulty.Int64)}
	}

	if obj.FeminicideForm1TransportDifficultyExplanation.Valid {
		dto.FeminicideForm1TransportDifficultyExplanation = obj.FeminicideForm1TransportDifficultyExplanation.String
	}

	if obj.FeminicideForm1EconomicResourcesForTransport.Valid {
		dto.FeminicideForm1EconomicResourcesForTransport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1EconomicResourcesForTransport.Int64)}
	}

	if obj.FeminicideForm1DebtOrHelpDueToTransport.Valid {
		dto.FeminicideForm1DebtOrHelpDueToTransport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1DebtOrHelpDueToTransport.Int64)}
	}

	if obj.FeminicideForm1DebtImpactExplanation.Valid {
		dto.FeminicideForm1DebtImpactExplanation = obj.FeminicideForm1DebtImpactExplanation.String
	}

	if obj.FeminicideForm1TransportSubsidyReceived.Valid {
		dto.FeminicideForm1TransportSubsidyReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1TransportSubsidyReceived.Int64)}
	}

	if obj.FeminicideForm1TransportSubsidyExplanation.Valid {
		dto.FeminicideForm1TransportSubsidyExplanation = obj.FeminicideForm1TransportSubsidyExplanation.String
	}

	if obj.FeminicideForm1SafetyTransportationConcern.Valid {
		dto.FeminicideForm1SafetyTransportationConcern = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1SafetyTransportationConcern.Int64)}
	}

	if obj.FeminicideForm1SafetyTransportationExplanation.Valid {
		dto.FeminicideForm1SafetyTransportationExplanation = obj.FeminicideForm1SafetyTransportationExplanation.String
	}

	if obj.FeminicideForm1FoodAccessFrequency.Valid {
		dto.FeminicideForm1FoodAccessFrequency = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1FoodAccessFrequency.Int64)}
	}

	if obj.FeminicideForm1FoodAccessExplanation.Valid {
		dto.FeminicideForm1FoodAccessExplanation = obj.FeminicideForm1FoodAccessExplanation.String
	}

	if obj.FeminicideForm1AggressorFoodRestriction.Valid {
		dto.FeminicideForm1AggressorFoodRestriction = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AggressorFoodRestriction.Int64)}
	}

	if obj.FeminicideForm1AggressorFoodRestrictionExplanation.Valid {
		dto.FeminicideForm1AggressorFoodRestrictionExplanation = obj.FeminicideForm1AggressorFoodRestrictionExplanation.String
	}

	if obj.FeminicideForm1FoodIncomeSupport.Valid {
		dto.FeminicideForm1FoodIncomeSupport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1FoodIncomeSupport.Int64)}
	}

	if obj.FeminicideForm1FoodIncomeSupportExplanation.Valid {
		dto.FeminicideForm1FoodIncomeSupportExplanation = obj.FeminicideForm1FoodIncomeSupportExplanation.String
	}

	if obj.FeminicideForm1JuridicalAssistanceReceived.Valid {
		dto.FeminicideForm1JuridicalAssistanceReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1JuridicalAssistanceReceived.Int64)}
	}

	if obj.FeminicideForm1JuridicalAssistanceExplanation.Valid {
		dto.FeminicideForm1JuridicalAssistanceExplanation = obj.FeminicideForm1JuridicalAssistanceExplanation.String
	}

	if obj.FeminicideForm1VictimRepresentation.Valid {
		dto.FeminicideForm1VictimRepresentation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1VictimRepresentation.Int64)}
	}

	if obj.FeminicideForm1VictimRepresentationExplanation.Valid {
		dto.FeminicideForm1VictimRepresentationExplanation = obj.FeminicideForm1VictimRepresentationExplanation.String
	}

	if obj.FeminicideForm1PsychosocialSupportReceived.Valid {
		dto.FeminicideForm1PsychosocialSupportReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1PsychosocialSupportReceived.Int64)}
	}

	if obj.FeminicideForm1PsychosocialSupportExplanation.Valid {
		dto.FeminicideForm1PsychosocialSupportExplanation = obj.FeminicideForm1PsychosocialSupportExplanation.String
	}

	if obj.FeminicideForm1EmergencyEmotionalCrisis.Valid {
		dto.FeminicideForm1EmergencyEmotionalCrisis = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1EmergencyEmotionalCrisis.Int64)}
	}

	if obj.FeminicideForm1EmergencyEmotionalCrisisExplanation.Valid {
		dto.FeminicideForm1EmergencyEmotionalCrisisExplanation = obj.FeminicideForm1EmergencyEmotionalCrisisExplanation.String
	}

	if obj.FeminicideForm1AggressorSameResidence.Valid {
		dto.FeminicideForm1AggressorSameResidence = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AggressorSameResidence.Int64)}
	}

	if obj.FeminicideForm1AggressorSameResidenceExplanation.Valid {
		dto.FeminicideForm1AggressorSameResidenceExplanation = obj.FeminicideForm1AggressorSameResidenceExplanation.String
	}

	if obj.FeminicideForm1AggressorLocationKnown.Valid {
		dto.FeminicideForm1AggressorLocationKnown = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AggressorLocationKnown.Int64)}
	}

	if obj.FeminicideForm1AggressorLocationKnownExplanation.Valid {
		dto.FeminicideForm1AggressorLocationKnownExplanation = obj.FeminicideForm1AggressorLocationKnownExplanation.String
	}

	if obj.FeminicideForm1AnyTypeOfAssistanceReceived.Valid {
		dto.FeminicideForm1AnyTypeOfAssistanceReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyTypeOfAssistanceReceived.Int64)}
	}

	if obj.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation.Valid {
		dto.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation = obj.FeminicideForm1AnyTypeOfAssistanceReceivedExplanation.String
	}

	if obj.FeminicideForm1CompensationFundInsuranceCoverage.Valid {
		dto.FeminicideForm1CompensationFundInsuranceCoverage = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1CompensationFundInsuranceCoverage.Int64)}
	}

	if obj.FeminicideForm1CompensationFundInsuranceCoverageExplanation.Valid {
		dto.FeminicideForm1CompensationFundInsuranceCoverageExplanation = obj.FeminicideForm1CompensationFundInsuranceCoverageExplanation.String
	}

	if obj.FeminicideForm1FuneralSubsidyReceived.Valid {
		dto.FeminicideForm1FuneralSubsidyReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1FuneralSubsidyReceived.Int64)}
	}

	if obj.FeminicideForm1FuneralSubsidyExplanation.Valid {
		dto.FeminicideForm1FuneralSubsidyExplanation = obj.FeminicideForm1FuneralSubsidyExplanation.String
	}

	if obj.FeminicideForm1FuneralFundsAvailable.Valid {
		dto.FeminicideForm1FuneralFundsAvailable = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1FuneralFundsAvailable.Int64)}
	}

	if obj.FeminicideForm1FuneralFundsAvailableExplanation.Valid {
		dto.FeminicideForm1FuneralFundsAvailableExplanation = obj.FeminicideForm1FuneralFundsAvailableExplanation.String
	}

	if obj.FeminicideForm1FuneralCostValue.Valid {
		dto.FeminicideForm1FuneralCostValue = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1FuneralCostValue.Int64)}
	}

	if obj.FeminicideForm1RenameReputationImpact.Valid {
		dto.FeminicideForm1RenameReputationImpact = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1RenameReputationImpact.Int64)}
	}

	if obj.FeminicideForm1RenameReputationExplanation.Valid {
		dto.FeminicideForm1RenameReputationExplanation = obj.FeminicideForm1RenameReputationExplanation.String
	}

	if obj.FeminicideForm1AdditionalNeedsDescription.Valid {
		dto.FeminicideForm1AdditionalNeedsDescription = obj.FeminicideForm1AdditionalNeedsDescription.String
	}

	if obj.FeminicideForm1ActionPlan.Valid {
		dto.FeminicideForm1ActionPlan = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1ActionPlan.Int64)}
	}

	if obj.FeminicideForm1AnyAssistanceReceivedCityHall.Valid {
		dto.FeminicideForm1AnyAssistanceReceivedCityHall = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyAssistanceReceivedCityHall.Int64)}
	}

	if obj.FeminicideForm1AnyAssistanceReceivedWomensOffice.Valid {
		dto.FeminicideForm1AnyAssistanceReceivedWomensOffice = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyAssistanceReceivedWomensOffice.Int64)}
	}

	if obj.FeminicideForm1AnyAssistanceReceivedOtherEntity.Valid {
		dto.FeminicideForm1AnyAssistanceReceivedOtherEntity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyAssistanceReceivedOtherEntity.Int64)}
	}

	if obj.FeminicideForm1AnyAssistanceReceivedOther.Valid {
		dto.FeminicideForm1AnyAssistanceReceivedOther = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideForm1AnyAssistanceReceivedOther.Int64)}
	}

	if obj.FeminicideForm1AssistanceReceived.Valid {
		dto.FeminicideForm1AssistanceReceived = obj.FeminicideForm1AssistanceReceived.String
	}

	if obj.FeminicideForm1AssistanceReceivedCityHall.Valid {
		dto.FeminicideForm1AssistanceReceivedCityHall = obj.FeminicideForm1AssistanceReceivedCityHall.String
	}

	if obj.FeminicideForm1AssistanceReceivedWomensOffice.Valid {
		dto.FeminicideForm1AssistanceReceivedWomensOffice = obj.FeminicideForm1AssistanceReceivedWomensOffice.String
	}

	if obj.FeminicideForm1AssistanceReceivedOtherEntity.Valid {
		dto.FeminicideForm1AssistanceReceivedOtherEntity = obj.FeminicideForm1AssistanceReceivedOtherEntity.String
	}

	if obj.FeminicideForm1AssistanceReceivedOther.Valid {
		dto.FeminicideForm1AssistanceReceivedOther = obj.FeminicideForm1AssistanceReceivedOther.String
	}

	if obj.FeminicideForm1FamilyMother.Valid {
		dto.FeminicideForm1FamilyMother = obj.FeminicideForm1FamilyMother.Int64
	}

	if obj.FeminicideForm1FamilyFather.Valid {
		dto.FeminicideForm1FamilyFather = obj.FeminicideForm1FamilyFather.Int64
	}

	if obj.FeminicideForm1FamilyStepfather.Valid {
		dto.FeminicideForm1FamilyStepfather = obj.FeminicideForm1FamilyStepfather.Int64
	}

	if obj.FeminicideForm1FamilyStepmother.Valid {
		dto.FeminicideForm1FamilyStepmother = obj.FeminicideForm1FamilyStepmother.Int64
	}

	if obj.FeminicideForm1FamilyPartner.Valid {
		dto.FeminicideForm1FamilyPartner = obj.FeminicideForm1FamilyPartner.Int64
	}

	if obj.FeminicideForm1FamilySibling1.Valid {
		dto.FeminicideForm1FamilySibling1 = obj.FeminicideForm1FamilySibling1.Int64
	}

	if obj.FeminicideForm1FamilySibling2.Valid {
		dto.FeminicideForm1FamilySibling2 = obj.FeminicideForm1FamilySibling2.Int64
	}

	if obj.FeminicideForm1FamilySibling3.Valid {
		dto.FeminicideForm1FamilySibling3 = obj.FeminicideForm1FamilySibling3.Int64
	}

	if obj.FeminicideForm1FamilySibling4.Valid {
		dto.FeminicideForm1FamilySibling4 = obj.FeminicideForm1FamilySibling4.Int64
	}

	if obj.FeminicideForm1FamilySibling5.Valid {
		dto.FeminicideForm1FamilySibling5 = obj.FeminicideForm1FamilySibling5.Int64
	}

	if obj.FeminicideForm1FamilySonDaughter1.Valid {
		dto.FeminicideForm1FamilySonDaughter1 = obj.FeminicideForm1FamilySonDaughter1.Int64
	}

	if obj.FeminicideForm1FamilySonDaughter2.Valid {
		dto.FeminicideForm1FamilySonDaughter2 = obj.FeminicideForm1FamilySonDaughter2.Int64
	}

	if obj.FeminicideForm1FamilySonDaughter3.Valid {
		dto.FeminicideForm1FamilySonDaughter3 = obj.FeminicideForm1FamilySonDaughter3.Int64
	}

	if obj.FeminicideForm1FamilySonDaughter4.Valid {
		dto.FeminicideForm1FamilySonDaughter4 = obj.FeminicideForm1FamilySonDaughter4.Int64
	}

	if obj.FeminicideForm1FamilySonDaughter5.Valid {
		dto.FeminicideForm1FamilySonDaughter5 = obj.FeminicideForm1FamilySonDaughter5.Int64
	}

	if obj.FeminicideForm1FamilyGrandmother.Valid {
		dto.FeminicideForm1FamilyGrandmother = obj.FeminicideForm1FamilyGrandmother.Int64
	}

	if obj.FeminicideForm1FamilyGrandfather.Valid {
		dto.FeminicideForm1FamilyGrandfather = obj.FeminicideForm1FamilyGrandfather.Int64
	}

	if obj.FeminicideForm1FamilyOtherMember.Valid {
		dto.FeminicideForm1FamilyOtherMember = obj.FeminicideForm1FamilyOtherMember.String
	}

	if obj.FeminicideForm1Summary.Valid {
		dto.FeminicideForm1Summary = obj.FeminicideForm1Summary.String
	}

	if obj.FeminicideForm1Feminicide.Valid {
		dto.FeminicideForm1Feminicide = FeminicideDTO{FeminicideId: uint64(obj.FeminicideForm1Feminicide.Int64)}
	}

	return dto
}
