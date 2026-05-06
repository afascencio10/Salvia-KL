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
	FeminicideRiskForm1EntityName string = "FeminicideRiskForm1"
	FeminicideRiskForm1JSONName   string = "form"
	FeminicideRiskForm1DBName     string = "feminicide_risk_form1"
	FeminicideRiskForm1DBScheme   string = "salvia"

	FeminicideRiskForm1FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FeminicideRiskForm1Id":                                       {Name: "FeminicideRiskForm1Id", DBName: "feminicide_risk_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1ICode":                                    {Name: "FeminicideRiskForm1ICode", DBName: "feminicide_risk_form1_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"FeminicideRiskForm1CreationDate":                             {Name: "FeminicideRiskForm1CreationDate", DBName: "feminicide_risk_form1_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1UpdateDate":                               {Name: "FeminicideRiskForm1UpdateDate", DBName: "feminicide_risk_form1_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimIdentityName":                       {Name: "FeminicideRiskForm1VictimIdentityName", DBName: "feminicide_risk_form1_victim_identity_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"FeminicideRiskForm1BirthDate":                                {Name: "FeminicideRiskForm1BirthDate", DBName: "feminicide_risk_form1_birth_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimAddress":                            {Name: "FeminicideRiskForm1VictimAddress", DBName: "feminicide_risk_form1_victim_address", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 32, Required: true},
		"FeminicideRiskForm1VictimLivingZone":                         {Name: "FeminicideRiskForm1VictimLivingZone", DBName: "feminicide_risk_form1_victim_living_zone", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimLivingTownCode":                     {Name: "FeminicideRiskForm1VictimLivingTownCode", DBName: "feminicide_risk_form1_victim_living_town", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 8, Required: true},
		"FeminicideRiskForm1VictimSGSSSAffiliation":                   {Name: "FeminicideRiskForm1VictimSGSSSAffiliation", DBName: "feminicide_risk_form1_victim_s_g_s_s_s_affiliation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1EpsName":                                  {Name: "FeminicideForm1EpsName", DBName: "feminicide_risk_form1_eps_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: false},
		"FeminicideRiskForm1ContactPhone":                             {Name: "FeminicideRiskForm1ContactPhone", DBName: "feminicide_risk_form1_contact_phone", Alias: "", ModelType: "string", MinSize: 10, MaxSize: 10, Required: true},
		"FeminicideRiskForm1ContactEmergencyContactNames":             {Name: "FeminicideRiskForm1ContactEmergencyContactNames", DBName: "feminicide_risk_form1_contact_emergency_contact_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: true},
		"FeminicideRiskForm1ContactEmergencyContactNumber":            {Name: "FeminicideRiskForm1ContactEmergencyContactNumber", DBName: "feminicide_risk_form1_contact_emergency_contact_number", Alias: "", ModelType: "string", MinSize: 10, MaxSize: 10, Required: true},
		"FeminicideRiskForm1VictimMaritalStatus":                      {Name: "FeminicideRiskForm1VictimMaritalStatus", DBName: "feminicide_risk_form1_victim_marital_status", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimSex":                                {Name: "FeminicideRiskForm1VictimSex", DBName: "feminicide_risk_form1_victim_sex", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimGenderIdentity":                     {Name: "FeminicideRiskForm1VictimGenderIdentity", DBName: "feminicide_risk_form1_victim_gender_identity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimSexualOrientation":                  {Name: "FeminicideRiskForm1VictimSexualOrientation", DBName: "feminicide_risk_form1_victim_sexual_orientation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimEthnicAffiliation":                  {Name: "FeminicideRiskForm1VictimEthnicAffiliation", DBName: "feminicide_risk_form1_victim_ethnic_affiliation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimIndigenousPeople":                   {Name: "FeminicideRiskForm1VictimIndigenousPeople", DBName: "feminicide_risk_form1_victim_indigenous_people", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1VictimIsMigrant":                          {Name: "FeminicideRiskForm1VictimIsMigrant", DBName: "feminicide_risk_form1_victim_is_migrant", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimMigrationStatus":                    {Name: "FeminicideRiskForm1VictimMigrationStatus", DBName: "feminicide_risk_form1_victim_migration_status", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1VictimMaxEducationLevel":                  {Name: "FeminicideRiskForm1VictimMaxEducationLevel", DBName: "feminicide_risk_form1_victim_max_education_level", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimIsSpecialPopulation":                {Name: "FeminicideRiskForm1VictimIsSpecialPopulation", DBName: "feminicide_risk_form1_victim_is_special_population", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimCurrentlyHasJob":                    {Name: "FeminicideRiskForm1VictimCurrentlyHasJob", DBName: "feminicide_risk_form1_victim_currently_has_job", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimJobExplanation":                     {Name: "FeminicideRiskForm1VictimJobExplanation", DBName: "feminicide_risk_form1_victim_job_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimAbandonedJobDueToRisk":              {Name: "FeminicideRiskForm1VictimAbandonedJobDueToRisk", DBName: "feminicide_risk_form1_victim_abandoned_job_due_to_risk", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimMainOccupation":                     {Name: "FeminicideRiskForm1VictimMainOccupation", DBName: "feminicide_risk_form1_victim_main_occupation", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimHasDisability":                      {Name: "FeminicideRiskForm1VictimHasDisability", DBName: "feminicide_risk_form1_victim_has_disability", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimDisabilityType":                     {Name: "FeminicideRiskForm1VictimDisabilityType", DBName: "feminicide_risk_form1_victim_disability_type", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1VictimRiskDescription":                    {Name: "FeminicideRiskForm1VictimRiskDescription", DBName: "feminicide_risk_form1_victim_risk_description", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: true},
		"FeminicideRiskForm1VictimAdditionalInfo":                     {Name: "FeminicideRiskForm1VictimAdditionalInfo", DBName: "feminicide_risk_form1_victim_additional_info", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimIsEconomicProvider":                 {Name: "FeminicideRiskForm1VictimIsEconomicProvider", DBName: "feminicide_risk_form1_victim_is_economic_provider", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimEconomicProviderExplanation":        {Name: "FeminicideRiskForm1VictimEconomicProviderExplanation", DBName: "feminicide_risk_form1_victim_economic_provider_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimHasFamiliarSupport":                 {Name: "FeminicideRiskForm1VictimHasFamiliarSupport", DBName: "feminicide_risk_form1_victim_has_familiar_support", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimFamiliarSupportType":                {Name: "FeminicideRiskForm1VictimFamiliarSupportType", DBName: "feminicide_risk_form1_victim_familiar_support_type", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimPublicTransportAccess":              {Name: "FeminicideRiskForm1VictimPublicTransportAccess", DBName: "feminicide_risk_form1_victim_public_transport_access", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimCommonTransportMode":                {Name: "FeminicideRiskForm1VictimCommonTransportMode", DBName: "feminicide_risk_form1_victim_common_transport_mode", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimTransportModeExplanation":           {Name: "FeminicideRiskForm1VictimTransportModeExplanation", DBName: "feminicide_risk_form1_victim_transport_mode_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimEstimatedTravelCost":                {Name: "FeminicideRiskForm1VictimEstimatedTravelCost", DBName: "feminicide_risk_form1_victim_estimated_travel_cost", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimDifficultiesWithTransport":          {Name: "FeminicideRiskForm1VictimDifficultiesWithTransport", DBName: "feminicide_risk_form1_victim_difficulties_with_transport", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimDifficultiesExplanation":            {Name: "FeminicideRiskForm1VictimDifficultiesExplanation", DBName: "feminicide_risk_form1_victim_difficulties_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimEconomicResourcesForTransport":      {Name: "FeminicideRiskForm1VictimEconomicResourcesForTransport", DBName: "feminicide_risk_form1_victim_economic_resources_for_transport", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimEconomicResourcesExplanation":       {Name: "FeminicideRiskForm1VictimEconomicResourcesExplanation", DBName: "feminicide_risk_form1_victim_economic_resources_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport":        {Name: "FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport", DBName: "feminicide_risk_form1_victim_has_debt_or_help_due_to_transport", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimDebtImpactDetails":                  {Name: "FeminicideRiskForm1VictimDebtImpactDetails", DBName: "feminicide_risk_form1_victim_debt_impact_details", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimReceivedTransportSubsidy":           {Name: "FeminicideRiskForm1VictimReceivedTransportSubsidy", DBName: "feminicide_risk_form1_victim_received_transport_subsidy", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimTransportSubsidyExplanation":        {Name: "FeminicideRiskForm1VictimTransportSubsidyExplanation", DBName: "feminicide_risk_form1_victim_transport_subsidy_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimSafetyAvoidedTransport":             {Name: "FeminicideRiskForm1VictimSafetyAvoidedTransport", DBName: "feminicide_risk_form1_victim_safety_avoided_transport", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimSafetyAvoidedExplanation":           {Name: "FeminicideRiskForm1VictimSafetyAvoidedExplanation", DBName: "feminicide_risk_form1_victim_safety_avoided_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimFoodAccessFrequency":                {Name: "FeminicideRiskForm1VictimFoodAccessFrequency", DBName: "feminicide_risk_form1_victim_food_access_frequency", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimFoodAccessExplanation":              {Name: "FeminicideRiskForm1VictimFoodAccessExplanation", DBName: "feminicide_risk_form1_victim_food_access_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimAgressorFoodRestriction":            {Name: "FeminicideRiskForm1VictimAgressorFoodRestriction", DBName: "feminicide_risk_form1_victim_agressor_food_restriction", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation": {Name: "FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation", DBName: "feminicide_risk_form1_victim_agressor_food_restrict_explanat", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimFamilyFixedIncome":                  {Name: "FeminicideRiskForm1VictimFamilyFixedIncome", DBName: "feminicide_risk_form1_victim_family_fixed_income", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimFamilyFixedIncomeExplanation":       {Name: "FeminicideRiskForm1VictimFamilyFixedIncomeExplanation", DBName: "feminicide_risk_form1_victim_family_fixed_income_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimIsOnlyProviderForFood":              {Name: "FeminicideRiskForm1VictimIsOnlyProviderForFood", DBName: "feminicide_risk_form1_victim_is_only_provider_for_food", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimIsOnlyProviderExplanation":          {Name: "FeminicideRiskForm1VictimIsOnlyProviderExplanation", DBName: "feminicide_risk_form1_victim_is_only_provider_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimJuridicalAssistanceReceived":        {Name: "FeminicideRiskForm1VictimJuridicalAssistanceReceived", DBName: "feminicide_risk_form1_victim_juridical_assistance_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimJuridicalAssistanceExplain":         {Name: "FeminicideRiskForm1VictimJuridicalAssistanceExplain", DBName: "feminicide_risk_form1_victim_juridical_assistance_explain", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimWantsJuridicalAssistance":           {Name: "FeminicideRiskForm1VictimWantsJuridicalAssistance", DBName: "feminicide_risk_form1_victim_wants_juridical_assistance", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimRepresentationsOfVictims":           {Name: "FeminicideRiskForm1VictimRepresentationsOfVictims", DBName: "feminicide_risk_form1_victim_representations_of_victims", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimRepresentationsExplanation":         {Name: "FeminicideRiskForm1VictimRepresentationsExplanation", DBName: "feminicide_risk_form1_victim_representations_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimPsychosocialSupportReceived":        {Name: "FeminicideRiskForm1VictimPsychosocialSupportReceived", DBName: "feminicide_risk_form1_victim_psychosocial_support_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain": {Name: "FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain", DBName: "feminicide_risk_form1_victim_psychosocial_supp_received_explain", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimUrgentEmotionalCrisis":              {Name: "FeminicideRiskForm1VictimUrgentEmotionalCrisis", DBName: "feminicide_risk_form1_victim_urgent_emotional_crisis", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimUrgentCrisisExplanation":            {Name: "FeminicideRiskForm1VictimUrgentCrisisExplanation", DBName: "feminicide_risk_form1_victim_urgent_crisis_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1AggressorSameResidence":                   {Name: "FeminicideRiskForm1AggressorSameResidence", DBName: "feminicide_risk_form1_aggressor_same_residence", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1AggressorSameResidenceExplanation":        {Name: "FeminicideRiskForm1AggressorSameResidenceExplanation", DBName: "feminicide_risk_form1_aggressor_same_residence_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1AggressorKnowsVictimLocation":             {Name: "FeminicideRiskForm1AggressorKnowsVictimLocation", DBName: "feminicide_risk_form1_aggressor_knows_victim_location", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1AggressorKnowsVictimLocationExplanation":  {Name: "FeminicideRiskForm1AggressorKnowsVictimLocationExplanation", DBName: "feminicide_risk_form1_aggressor_knows_victim_loc_explanation", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimHousingHelpReceived":                {Name: "FeminicideRiskForm1VictimHousingHelpReceived", DBName: "feminicide_risk_form1_victim_housing_help_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimHousingHelpExplain":                 {Name: "FeminicideRiskForm1VictimHousingHelpExplain", DBName: "feminicide_risk_form1_victim_housing_help_explain", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimAbandonClothing":                    {Name: "FeminicideRiskForm1VictimAbandonClothing", DBName: "feminicide_risk_form1_victim_abandon_clothing", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimAbandonClothingExplain":             {Name: "FeminicideRiskForm1VictimAbandonClothingExplain", DBName: "feminicide_risk_form1_victim_abandon_clothing_explain", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimClothingHelpReceived":               {Name: "FeminicideRiskForm1VictimClothingHelpReceived", DBName: "feminicide_risk_form1_victim_clothing_help_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1VictimClothingHelpExplain":                {Name: "FeminicideRiskForm1VictimClothingHelpExplain", DBName: "feminicide_risk_form1_victim_clothing_help_explain", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1VictimOtherNeeds":                         {Name: "FeminicideRiskForm1VictimOtherNeeds", DBName: "feminicide_risk_form1_victim_other_needs", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 10000, Required: false},
		"FeminicideRiskForm1InterviewDate":                            {Name: "FeminicideRiskForm1InterviewDate", DBName: "feminicide_risk_form1_interview_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1Summary":                                  {Name: "FeminicideRiskForm1Summary", DBName: "feminicide_risk_form1_summary", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 20000, Required: false},
		"FeminicideRiskForm1FeminicideRisk":                           {Name: "FeminicideRiskForm1FeminicideRisk", DBName: "feminicide_form1_feminicide_risk", Alias: "", ModelType: "uint", Required: true},

		"FeminicideRiskForm1AnyAssistanceReceived":             {Name: "FeminicideForm1AnyAssistanceReceived", DBName: "feminicide_form1_any_assistance_received", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskForm1AnyAssistanceReceivedCityHall":     {Name: "FeminicideRiskForm1AnyAssistanceReceivedCityHall", DBName: "feminicide_risk_form1_any_assistance_received_city_hall", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1AnyAssistanceReceivedWomensOffice": {Name: "FeminicideRiskForm1AnyAssistanceReceivedWomensOffice", DBName: "feminicide_risk_form1_any_assistance_received_womens_office", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1AnyAssistanceReceivedOtherEntity":  {Name: "FeminicideRiskForm1AnyAssistanceReceivedOtherEntity", DBName: "feminicide_risk_form1_any_assistance_received_other_entity", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1AnyAssistanceReceivedOther":        {Name: "FeminicideRiskForm1AnyAssistanceReceivedOther", DBName: "feminicide_risk_form1_any_assistance_received_other", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"FeminicideRiskForm1AssistanceReceived":                {Name: "FeminicideRiskForm1AssistanceReceived", DBName: "feminicide_risk_form1_assistance_received", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideRiskForm1AssistanceReceivedCityHall":        {Name: "FeminicideRiskForm1AssistanceReceivedCityHall", DBName: "feminicide_risk_form1_assistance_received_city_hall", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideRiskForm1AssistanceReceivedWomensOffice":    {Name: "FeminicideRiskForm1AssistanceReceivedWomensOffice", DBName: "feminicide_risk_form1_assistance_received_womens_office", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideRiskForm1AssistanceReceivedOtherEntity":     {Name: "FeminicideRiskForm1AssistanceReceivedOtherEntity", DBName: "feminicide_risk_form1_assistance_received_other_entity", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"FeminicideRiskForm1AssistanceReceivedOther":           {Name: "FeminicideRiskForm1AssistanceReceivedOther", DBName: "feminicide_risk_form1_assistance_received_other", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},

		"FeminicideRiskForm1FamilyMother":       {Name: "FeminicideRiskForm1FamilyMother", DBName: "feminicide_risk_form1_family_mother", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyFather":       {Name: "FeminicideRiskForm1FamilyFather", DBName: "feminicide_risk_form1_family_father", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyStepfather":   {Name: "FeminicideRiskForm1FamilyStepfather", DBName: "feminicide_risk_form1_family_stepfather", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyStepmother":   {Name: "FeminicideRiskForm1FamilyStepmother", DBName: "feminicide_risk_form1_family_stepmother", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyPartner":      {Name: "FeminicideRiskForm1FamilyPartner", DBName: "feminicide_risk_form1_family_partner", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySibling1":     {Name: "FeminicideRiskForm1FamilySibling1", DBName: "feminicide_risk_form1_family_sibling_1", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySibling2":     {Name: "FeminicideRiskForm1FamilySibling2", DBName: "feminicide_risk_form1_family_sibling_2", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySibling3":     {Name: "FeminicideRiskForm1FamilySibling3", DBName: "feminicide_risk_form1_family_sibling_3", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySibling4":     {Name: "FeminicideRiskForm1FamilySibling4", DBName: "feminicide_risk_form1_family_sibling_4", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySibling5":     {Name: "FeminicideRiskForm1FamilySibling5", DBName: "feminicide_risk_form1_family_sibling_5", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySonDaughter1": {Name: "FeminicideRiskForm1FamilySonDaughter1", DBName: "feminicide_risk_form1_family_son_daughter_1", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySonDaughter2": {Name: "FeminicideRiskForm1FamilySonDaughter2", DBName: "feminicide_risk_form1_family_son_daughter_2", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySonDaughter3": {Name: "FeminicideRiskForm1FamilySonDaughter3", DBName: "feminicide_risk_form1_family_son_daughter_3", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySonDaughter4": {Name: "FeminicideRiskForm1FamilySonDaughter4", DBName: "feminicide_risk_form1_family_son_daughter_4", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilySonDaughter5": {Name: "FeminicideRiskForm1FamilySonDaughter5", DBName: "feminicide_risk_form1_family_son_daughter_5", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyGrandmother":  {Name: "FeminicideRiskForm1FamilyGrandmother", DBName: "feminicide_risk_form1_family_grandmother", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyGrandfather":  {Name: "FeminicideRiskForm1FamilyGrandfather", DBName: "feminicide_risk_form1_family_grandfather", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 120, Required: true},
		"FeminicideRiskForm1FamilyOtherMember":  {Name: "FeminicideRiskForm1FamilyOtherMember", DBName: "feminicide_risk_form1_family_other_member", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
	}
)

type FeminicideRiskForm1DTO struct {
	FeminicideRiskForm1Id                                       uint64                  `json:"-"`
	FeminicideRiskForm1ICode                                    string                  `json:"icode"`
	FeminicideRiskForm1CreationDate                             time.Time               `json:"creationDate"`
	FeminicideRiskForm1UpdateDate                               time.Time               `json:"updateDate"`
	FeminicideRiskForm1VictimIdentityName                       string                  `json:"victimIdentityName"`
	FeminicideRiskForm1BirthDate                                time.Time               `json:"birthDate"`
	FeminicideRiskForm1VictimAddress                            string                  `json:"victimAddress"`
	FeminicideRiskForm1VictimLivingZone                         VictimCaseForm2EnumsDTO `json:"victimLivingZone"`
	FeminicideRiskForm1VictimLivingTownCode                     string                  `json:"victimLivingTownCode"`
	FeminicideRiskForm1VictimSGSSSAffiliation                   VictimCaseForm2EnumsDTO `json:"victimSGSSSAffiliation"`
	FeminicideRiskForm1EpsName                                  string                  `json:"epsName"`
	FeminicideRiskForm1ContactPhone                             string                  `json:"contactPhone"`
	FeminicideRiskForm1ContactEmergencyContactNames             string                  `json:"contactEmergencyContactNames"`
	FeminicideRiskForm1ContactEmergencyContactNumber            string                  `json:"contactEmergencyContactNumber"`
	FeminicideRiskForm1VictimMaritalStatus                      VictimCaseForm2EnumsDTO `json:"victimMaritalStatus"`
	FeminicideRiskForm1VictimSex                                VictimCaseForm2EnumsDTO `json:"victimSex"`
	FeminicideRiskForm1VictimGenderIdentity                     VictimCaseForm2EnumsDTO `json:"victimGenderIdentity"`
	FeminicideRiskForm1VictimSexualOrientation                  VictimCaseForm2EnumsDTO `json:"victimSexualOrientation"`
	FeminicideRiskForm1VictimEthnicAffiliation                  VictimCaseForm2EnumsDTO `json:"victimEthnicAffiliation"`
	FeminicideRiskForm1VictimIndigenousPeople                   VictimCaseForm2EnumsDTO `json:"victimIndigenousPeople"`
	FeminicideRiskForm1VictimIsMigrant                          VictimCaseForm2EnumsDTO `json:"victimIsMigrant"`
	FeminicideRiskForm1VictimMigrationStatus                    VictimCaseForm2EnumsDTO `json:"victimMigrationStatus"`
	FeminicideRiskForm1VictimMaxEducationLevel                  VictimCaseForm2EnumsDTO `json:"victimMaxEducationLevel"`
	FeminicideRiskForm1VictimIsSpecialPopulation                VictimCaseForm2EnumsDTO `json:"victimIsSpecialPopulation"`
	FeminicideRiskForm1VictimCurrentlyHasJob                    VictimCaseForm2EnumsDTO `json:"victimCurrentlyHasJob"`
	FeminicideRiskForm1VictimJobExplanation                     string                  `json:"victimJobExplanation"`
	FeminicideRiskForm1VictimAbandonedJobDueToRisk              VictimCaseForm2EnumsDTO `json:"victimAbandonedJobDueToRisk"`
	FeminicideRiskForm1VictimMainOccupation                     VictimCaseForm2EnumsDTO `json:"victimMainOccupation"`
	FeminicideRiskForm1VictimHasDisability                      VictimCaseForm2EnumsDTO `json:"victimHasDisability"`
	FeminicideRiskForm1VictimDisabilityType                     VictimCaseForm2EnumsDTO `json:"victimDisabilityType"`
	FeminicideRiskForm1VictimRiskDescription                    string                  `json:"victimRiskDescription"`
	FeminicideRiskForm1VictimAdditionalInfo                     string                  `json:"victimAdditionalInfo"`
	FeminicideRiskForm1VictimIsEconomicProvider                 VictimCaseForm2EnumsDTO `json:"victimIsEconomicProvider"`
	FeminicideRiskForm1VictimEconomicProviderExplanation        string                  `json:"victimEconomicProviderExplanation"`
	FeminicideRiskForm1VictimHasFamiliarSupport                 VictimCaseForm2EnumsDTO `json:"victimHasFamiliarSupport"`
	FeminicideRiskForm1VictimFamiliarSupportType                VictimCaseForm2EnumsDTO `json:"victimFamiliarSupportType"`
	FeminicideRiskForm1VictimPublicTransportAccess              VictimCaseForm2EnumsDTO `json:"victimPublicTransportAccess"`
	FeminicideRiskForm1VictimCommonTransportMode                VictimCaseForm2EnumsDTO `json:"victimCommonTransportMode"`
	FeminicideRiskForm1VictimTransportModeExplanation           string                  `json:"victimTransportModeExplanation"`
	FeminicideRiskForm1VictimEstimatedTravelCost                int32                   `json:"victimEstimatedTravelCost"`
	FeminicideRiskForm1VictimDifficultiesWithTransport          VictimCaseForm2EnumsDTO `json:"victimDifficultiesWithTransport"`
	FeminicideRiskForm1VictimDifficultiesExplanation            string                  `json:"victimDifficultiesExplanation"`
	FeminicideRiskForm1VictimEconomicResourcesForTransport      VictimCaseForm2EnumsDTO `json:"victimEconomicResourcesForTransport"`
	FeminicideRiskForm1VictimEconomicResourcesExplanation       string                  `json:"victimEconomicResourcesExplanation"`
	FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport        VictimCaseForm2EnumsDTO `json:"victimHasDebtOrHelpDueToTransport"`
	FeminicideRiskForm1VictimDebtImpactDetails                  string                  `json:"victimDebtImpactDetails"`
	FeminicideRiskForm1VictimReceivedTransportSubsidy           VictimCaseForm2EnumsDTO `json:"victimReceivedTransportSubsidy"`
	FeminicideRiskForm1VictimTransportSubsidyExplanation        string                  `json:"victimTransportSubsidyExplanation"`
	FeminicideRiskForm1VictimSafetyAvoidedTransport             VictimCaseForm2EnumsDTO `json:"victimSafetyAvoidedTransport"`
	FeminicideRiskForm1VictimSafetyAvoidedExplanation           string                  `json:"victimSafetyAvoidedExplanation"`
	FeminicideRiskForm1VictimFoodAccessFrequency                VictimCaseForm2EnumsDTO `json:"victimFoodAccessFrequency"`
	FeminicideRiskForm1VictimFoodAccessExplanation              string                  `json:"victimFoodAccessExplanation"`
	FeminicideRiskForm1VictimAgressorFoodRestriction            VictimCaseForm2EnumsDTO `json:"victimAgressorFoodRestriction"`
	FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation string                  `json:"victimAgressorFoodRestrictionExplanation"`
	FeminicideRiskForm1VictimFamilyFixedIncome                  VictimCaseForm2EnumsDTO `json:"victimFamilyFixedIncome"`
	FeminicideRiskForm1VictimFamilyFixedIncomeExplanation       string                  `json:"victimFamilyFixedIncomeExplanation"`
	FeminicideRiskForm1VictimIsOnlyProviderForFood              VictimCaseForm2EnumsDTO `json:"victimIsOnlyProviderForFood"`
	FeminicideRiskForm1VictimIsOnlyProviderExplanation          string                  `json:"victimIsOnlyProviderExplanation"`
	FeminicideRiskForm1VictimJuridicalAssistanceReceived        VictimCaseForm2EnumsDTO `json:"victimJuridicalAssistanceReceived"`
	FeminicideRiskForm1VictimJuridicalAssistanceExplain         string                  `json:"victimJuridicalAssistanceExplain"`
	FeminicideRiskForm1VictimWantsJuridicalAssistance           VictimCaseForm2EnumsDTO `json:"victimWantsJuridicalAssistance"`
	FeminicideRiskForm1VictimRepresentationsOfVictims           VictimCaseForm2EnumsDTO `json:"victimRepresentationsOfVictims"`
	FeminicideRiskForm1VictimRepresentationsExplanation         string                  `json:"victimRepresentationsExplanation"`
	FeminicideRiskForm1VictimPsychosocialSupportReceived        VictimCaseForm2EnumsDTO `json:"victimPsychosocialSupportReceived"`
	FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain string                  `json:"victimPsychosocialSupportReceivedExplain"`
	FeminicideRiskForm1VictimUrgentEmotionalCrisis              VictimCaseForm2EnumsDTO `json:"victimUrgentEmotionalCrisis"`
	FeminicideRiskForm1VictimUrgentCrisisExplanation            string                  `json:"victimUrgentCrisisExplanation"`
	FeminicideRiskForm1AggressorSameResidence                   VictimCaseForm2EnumsDTO `json:"aggressorSameResidence"`
	FeminicideRiskForm1AggressorSameResidenceExplanation        string                  `json:"aggressorSameResidenceExplanation"`
	FeminicideRiskForm1AggressorKnowsVictimLocation             VictimCaseForm2EnumsDTO `json:"aggressorKnowsVictimLocation"`
	FeminicideRiskForm1AggressorKnowsVictimLocationExplanation  string                  `json:"aggressorKnowsVictimLocationExplanation"`
	FeminicideRiskForm1VictimHousingHelpReceived                VictimCaseForm2EnumsDTO `json:"victimHousingHelpReceived"`
	FeminicideRiskForm1VictimHousingHelpExplain                 string                  `json:"victimHousingHelpExplain"`
	FeminicideRiskForm1VictimAbandonClothing                    VictimCaseForm2EnumsDTO `json:"victimAbandonClothing"`
	FeminicideRiskForm1VictimAbandonClothingExplain             string                  `json:"victimAbandonClothingExplain"`
	FeminicideRiskForm1VictimClothingHelpReceived               VictimCaseForm2EnumsDTO `json:"victimClothingHelpReceived"`
	FeminicideRiskForm1VictimClothingHelpExplain                string                  `json:"victimClothingHelpExplain"`
	FeminicideRiskForm1VictimOtherNeeds                         string                  `json:"victimOtherNeeds"`
	FeminicideRiskForm1InterviewDate                            time.Time               `json:"interviewDate"`
	FeminicideRiskForm1Summary                                  string                  `json:"summary"`
	FeminicideRiskForm1FeminicideRisk                           interface{}             `json:"-"`

	FeminicideRiskForm1AnyAssistanceReceived             VictimCaseForm2EnumsDTO `json:"anyAssistanceReceived"`
	FeminicideRiskForm1AssistanceReceived                string                  `json:"assistanceReceived"`
	FeminicideRiskForm1AnyAssistanceReceivedCityHall     VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedCityHall"`
	FeminicideRiskForm1AssistanceReceivedCityHall        string                  `json:"assistanceReceivedCityHall"`
	FeminicideRiskForm1AnyAssistanceReceivedWomensOffice VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedWomensOffice"`
	FeminicideRiskForm1AssistanceReceivedWomensOffice    string                  `json:"assistanceReceivedWomensOffice"`
	FeminicideRiskForm1AnyAssistanceReceivedOtherEntity  VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedOtherEntity"`
	FeminicideRiskForm1AssistanceReceivedOtherEntity     string                  `json:"assistanceReceivedOtherEntity"`
	FeminicideRiskForm1AnyAssistanceReceivedOther        VictimCaseForm2EnumsDTO `json:"anyAssistanceReceivedOther"`
	FeminicideRiskForm1AssistanceReceivedOther           string                  `json:"assistanceReceivedOther"`

	FeminicideRiskForm1FamilyFather       int64  `json:"familyFather"`
	FeminicideRiskForm1FamilyMother       int64  `json:"familyMother"`
	FeminicideRiskForm1FamilyStepfather   int64  `json:"familyStepfather"`
	FeminicideRiskForm1FamilyStepmother   int64  `json:"familyStepmother"`
	FeminicideRiskForm1FamilyPartner      int64  `json:"familyPartner"`
	FeminicideRiskForm1FamilySibling1     int64  `json:"familySibling1"`
	FeminicideRiskForm1FamilySibling2     int64  `json:"familySibling2"`
	FeminicideRiskForm1FamilySibling3     int64  `json:"familySibling3"`
	FeminicideRiskForm1FamilySibling4     int64  `json:"familySibling4"`
	FeminicideRiskForm1FamilySibling5     int64  `json:"familySibling5"`
	FeminicideRiskForm1FamilySonDaughter1 int64  `json:"familySonDaughter1"`
	FeminicideRiskForm1FamilySonDaughter2 int64  `json:"familySonDaughter2"`
	FeminicideRiskForm1FamilySonDaughter3 int64  `json:"familySonDaughter3"`
	FeminicideRiskForm1FamilySonDaughter4 int64  `json:"familySonDaughter4"`
	FeminicideRiskForm1FamilySonDaughter5 int64  `json:"familySonDaughter5"`
	FeminicideRiskForm1FamilyGrandmother  int64  `json:"familyGrandmother"`
	FeminicideRiskForm1FamilyGrandfather  int64  `json:"familyGrandfather"`
	FeminicideRiskForm1FamilyOtherMember  string `json:"familyOtherMember"`

	//Campos múltiples
	FeminicideRiskForm1FinanciallyDependentPeople []VictimCaseForm2EnumsDTO `json:"financiallyDependentPeople"`
	FeminicideRiskForm1PlacesVisitRegularly       []VictimCaseForm2EnumsDTO `json:"placesVisitRegularly"`

	//Campos de formulario

	FeminicideRiskForm1VictimLivingDepartment security_daos.DepartmentDTO `json:"livingDepartment"`
	FeminicideRiskForm1VictimLivingCity       security_daos.CityDTO       `json:"livingCity"`
	FeminicideRiskForm1VictimLivingTown       security_daos.TownDTO       `json:"livingTown"`
}

type FeminicideRiskForm1PgDB struct {
	FeminicideRiskForm1Id                                       sql.NullInt64
	FeminicideRiskForm1ICode                                    sql.NullString
	FeminicideRiskForm1CreationDate                             sql.NullTime
	FeminicideRiskForm1UpdateDate                               sql.NullTime
	FeminicideRiskForm1VictimIdentityName                       sql.NullString
	FeminicideRiskForm1BirthDate                                sql.NullTime
	FeminicideRiskForm1VictimAddress                            sql.NullString
	FeminicideRiskForm1VictimLivingZone                         sql.NullInt64
	FeminicideRiskForm1VictimLivingTownCode                     sql.NullString
	FeminicideRiskForm1VictimSGSSSAffiliation                   sql.NullInt64
	FeminicideRiskForm1EpsName                                  sql.NullString
	FeminicideRiskForm1ContactPhone                             sql.NullString
	FeminicideRiskForm1ContactEmergencyContactNames             sql.NullString
	FeminicideRiskForm1ContactEmergencyContactNumber            sql.NullString
	FeminicideRiskForm1VictimMaritalStatus                      sql.NullInt64
	FeminicideRiskForm1VictimSex                                sql.NullInt64
	FeminicideRiskForm1VictimGenderIdentity                     sql.NullInt64
	FeminicideRiskForm1VictimSexualOrientation                  sql.NullInt64
	FeminicideRiskForm1VictimEthnicAffiliation                  sql.NullInt64
	FeminicideRiskForm1VictimIndigenousPeople                   sql.NullInt64
	FeminicideRiskForm1VictimIsMigrant                          sql.NullInt64
	FeminicideRiskForm1VictimMigrationStatus                    sql.NullInt64
	FeminicideRiskForm1VictimMaxEducationLevel                  sql.NullInt64
	FeminicideRiskForm1VictimIsSpecialPopulation                sql.NullInt64
	FeminicideRiskForm1VictimCurrentlyHasJob                    sql.NullInt64
	FeminicideRiskForm1VictimJobExplanation                     sql.NullString
	FeminicideRiskForm1VictimAbandonedJobDueToRisk              sql.NullInt64
	FeminicideRiskForm1VictimMainOccupation                     sql.NullInt64
	FeminicideRiskForm1VictimHasDisability                      sql.NullInt64
	FeminicideRiskForm1VictimDisabilityType                     sql.NullInt64
	FeminicideRiskForm1VictimRiskDescription                    sql.NullString
	FeminicideRiskForm1VictimAdditionalInfo                     sql.NullString
	FeminicideRiskForm1VictimIsEconomicProvider                 sql.NullInt64
	FeminicideRiskForm1VictimEconomicProviderExplanation        sql.NullString
	FeminicideRiskForm1VictimHasFamiliarSupport                 sql.NullInt64
	FeminicideRiskForm1VictimFamiliarSupportType                sql.NullInt64
	FeminicideRiskForm1VictimPublicTransportAccess              sql.NullInt64
	FeminicideRiskForm1VictimCommonTransportMode                sql.NullInt64
	FeminicideRiskForm1VictimTransportModeExplanation           sql.NullString
	FeminicideRiskForm1VictimEstimatedTravelCost                sql.NullInt32
	FeminicideRiskForm1VictimDifficultiesWithTransport          sql.NullInt64
	FeminicideRiskForm1VictimDifficultiesExplanation            sql.NullString
	FeminicideRiskForm1VictimEconomicResourcesForTransport      sql.NullInt64
	FeminicideRiskForm1VictimEconomicResourcesExplanation       sql.NullString
	FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport        sql.NullInt64
	FeminicideRiskForm1VictimDebtImpactDetails                  sql.NullString
	FeminicideRiskForm1VictimReceivedTransportSubsidy           sql.NullInt64
	FeminicideRiskForm1VictimTransportSubsidyExplanation        sql.NullString
	FeminicideRiskForm1VictimSafetyAvoidedTransport             sql.NullInt64
	FeminicideRiskForm1VictimSafetyAvoidedExplanation           sql.NullString
	FeminicideRiskForm1VictimFoodAccessFrequency                sql.NullInt64
	FeminicideRiskForm1VictimFoodAccessExplanation              sql.NullString
	FeminicideRiskForm1VictimAgressorFoodRestriction            sql.NullInt64
	FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation sql.NullString
	FeminicideRiskForm1VictimFamilyFixedIncome                  sql.NullInt64
	FeminicideRiskForm1VictimFamilyFixedIncomeExplanation       sql.NullString
	FeminicideRiskForm1VictimIsOnlyProviderForFood              sql.NullInt64
	FeminicideRiskForm1VictimIsOnlyProviderExplanation          sql.NullString
	FeminicideRiskForm1VictimJuridicalAssistanceReceived        sql.NullInt64
	FeminicideRiskForm1VictimJuridicalAssistanceExplain         sql.NullString
	FeminicideRiskForm1VictimWantsJuridicalAssistance           sql.NullInt64
	FeminicideRiskForm1VictimRepresentationsOfVictims           sql.NullInt64
	FeminicideRiskForm1VictimRepresentationsExplanation         sql.NullString
	FeminicideRiskForm1VictimPsychosocialSupportReceived        sql.NullInt64
	FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain sql.NullString
	FeminicideRiskForm1VictimUrgentEmotionalCrisis              sql.NullInt64
	FeminicideRiskForm1VictimUrgentCrisisExplanation            sql.NullString
	FeminicideRiskForm1AggressorSameResidence                   sql.NullInt64
	FeminicideRiskForm1AggressorSameResidenceExplanation        sql.NullString
	FeminicideRiskForm1AggressorKnowsVictimLocation             sql.NullInt64
	FeminicideRiskForm1AggressorKnowsVictimLocationExplanation  sql.NullString
	FeminicideRiskForm1VictimHousingHelpReceived                sql.NullInt64
	FeminicideRiskForm1VictimHousingHelpExplain                 sql.NullString
	FeminicideRiskForm1VictimAbandonClothing                    sql.NullInt64
	FeminicideRiskForm1VictimAbandonClothingExplain             sql.NullString
	FeminicideRiskForm1VictimClothingHelpReceived               sql.NullInt64
	FeminicideRiskForm1VictimClothingHelpExplain                sql.NullString
	FeminicideRiskForm1VictimOtherNeeds                         sql.NullString
	FeminicideRiskForm1InterviewDate                            sql.NullTime
	FeminicideRiskForm1Summary                                  sql.NullString
	FeminicideRiskForm1FeminicideRisk                           sql.NullInt64

	FeminicideRiskForm1AnyAssistanceReceived             sql.NullInt64
	FeminicideRiskForm1AnyAssistanceReceivedCityHall     sql.NullInt64
	FeminicideRiskForm1AnyAssistanceReceivedWomensOffice sql.NullInt64
	FeminicideRiskForm1AnyAssistanceReceivedOtherEntity  sql.NullInt64
	FeminicideRiskForm1AnyAssistanceReceivedOther        sql.NullInt64

	FeminicideRiskForm1AssistanceReceived             sql.NullString
	FeminicideRiskForm1AssistanceReceivedCityHall     sql.NullString
	FeminicideRiskForm1AssistanceReceivedWomensOffice sql.NullString
	FeminicideRiskForm1AssistanceReceivedOtherEntity  sql.NullString
	FeminicideRiskForm1AssistanceReceivedOther        sql.NullString
	FeminicideRiskForm1FamilyFather                   sql.NullInt64
	FeminicideRiskForm1FamilyMother                   sql.NullInt64
	FeminicideRiskForm1FamilyStepfather               sql.NullInt64
	FeminicideRiskForm1FamilyStepmother               sql.NullInt64
	FeminicideRiskForm1FamilyPartner                  sql.NullInt64
	FeminicideRiskForm1FamilySibling1                 sql.NullInt64
	FeminicideRiskForm1FamilySibling2                 sql.NullInt64
	FeminicideRiskForm1FamilySibling3                 sql.NullInt64
	FeminicideRiskForm1FamilySibling4                 sql.NullInt64
	FeminicideRiskForm1FamilySibling5                 sql.NullInt64
	FeminicideRiskForm1FamilySonDaughter1             sql.NullInt64
	FeminicideRiskForm1FamilySonDaughter2             sql.NullInt64
	FeminicideRiskForm1FamilySonDaughter3             sql.NullInt64
	FeminicideRiskForm1FamilySonDaughter4             sql.NullInt64
	FeminicideRiskForm1FamilySonDaughter5             sql.NullInt64
	FeminicideRiskForm1FamilyGrandmother              sql.NullInt64
	FeminicideRiskForm1FamilyGrandfather              sql.NullInt64
	FeminicideRiskForm1FamilyOtherMember              sql.NullString
}

func (f *FeminicideRiskForm1DTO) MarshalJSON() ([]byte, error) {
	type Alias FeminicideRiskForm1DTO

	return json.Marshal(&struct {
		*Alias
		FeminicideRiskForm1BirthDate     string `json:"birthDate"`
		FeminicideRiskForm1CreationDate  string `json:"creationDate"`
		FeminicideRiskForm1UpdateDate    string `json:"updateDate"`
		FeminicideRiskForm1InterviewDate string `json:"interviewDate"`
	}{
		Alias:                            (*Alias)(f),
		FeminicideRiskForm1BirthDate:     f.FeminicideRiskForm1BirthDate.Format(common_config.DateTime.DATE_FORMAT),
		FeminicideRiskForm1CreationDate:  f.FeminicideRiskForm1CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FeminicideRiskForm1UpdateDate:    f.FeminicideRiskForm1UpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FeminicideRiskForm1InterviewDate: time.Now().Format(common_config.DateTime.DATE_FORMAT),
	})
}

func (f *FeminicideRiskForm1DTO) UnmarshalJSON(data []byte) error {
	type Alias FeminicideRiskForm1DTO

	aux := &struct {
		*Alias
		FeminicideRiskForm1BirthDate     string `json:"birthDate"`
		FeminicideRiskForm1CreationDate  string `json:"creationDate"`
		FeminicideRiskForm1UpdateDate    string `json:"updateDate"`
		FeminicideRiskForm1InterviewDate string `json:"interviewDate"`
	}{
		Alias: (*Alias)(f),
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

	f.FeminicideRiskForm1BirthDate = parse(aux.FeminicideRiskForm1BirthDate, common_config.DateTime.DATE_FORMAT)
	f.FeminicideRiskForm1CreationDate = parse(aux.FeminicideRiskForm1CreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	f.FeminicideRiskForm1UpdateDate = parse(aux.FeminicideRiskForm1UpdateDate, common_config.DateTime.DATE_TIME_FORMAT)
	f.FeminicideRiskForm1InterviewDate = parse(aux.FeminicideRiskForm1InterviewDate, common_config.DateTime.DATE_FORMAT)

	return nil
}

func SetFeminicideRiskForm1(f *FeminicideRiskForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		return persistenceCtrl.Error
	}

	var fields []string = []string{
		"FeminicideRiskForm1ICode",
		"FeminicideRiskForm1CreationDate",
		"FeminicideRiskForm1UpdateDate",
		"FeminicideRiskForm1VictimIdentityName",
		"FeminicideRiskForm1BirthDate",
		"FeminicideRiskForm1VictimAddress",
		"FeminicideRiskForm1VictimLivingZone",
		"FeminicideRiskForm1VictimLivingTownCode",
		"FeminicideRiskForm1VictimSGSSSAffiliation",
		"FeminicideRiskForm1EpsName",
		"FeminicideRiskForm1ContactPhone",
		"FeminicideRiskForm1ContactEmergencyContactNames",
		"FeminicideRiskForm1ContactEmergencyContactNumber",
		"FeminicideRiskForm1VictimMaritalStatus",
		"FeminicideRiskForm1VictimSex",
		"FeminicideRiskForm1VictimGenderIdentity",
		"FeminicideRiskForm1VictimSexualOrientation",
		"FeminicideRiskForm1VictimEthnicAffiliation",
		"FeminicideRiskForm1VictimIndigenousPeople",
		"FeminicideRiskForm1VictimIsMigrant",
		"FeminicideRiskForm1VictimMigrationStatus",
		"FeminicideRiskForm1VictimMaxEducationLevel",
		"FeminicideRiskForm1VictimIsSpecialPopulation",
		"FeminicideRiskForm1VictimCurrentlyHasJob",
		"FeminicideRiskForm1VictimJobExplanation",
		"FeminicideRiskForm1VictimAbandonedJobDueToRisk",
		"FeminicideRiskForm1VictimMainOccupation",
		"FeminicideRiskForm1VictimHasDisability",
		"FeminicideRiskForm1VictimDisabilityType",
		"FeminicideRiskForm1VictimRiskDescription",
		"FeminicideRiskForm1VictimAdditionalInfo",
		"FeminicideRiskForm1VictimIsEconomicProvider",
		"FeminicideRiskForm1VictimEconomicProviderExplanation",
		"FeminicideRiskForm1VictimHasFamiliarSupport",
		"FeminicideRiskForm1VictimFamiliarSupportType",
		"FeminicideRiskForm1VictimPublicTransportAccess",
		"FeminicideRiskForm1VictimCommonTransportMode",
		"FeminicideRiskForm1VictimTransportModeExplanation",
		"FeminicideRiskForm1VictimEstimatedTravelCost",
		"FeminicideRiskForm1VictimDifficultiesWithTransport",
		"FeminicideRiskForm1VictimDifficultiesExplanation",
		"FeminicideRiskForm1VictimEconomicResourcesForTransport",
		"FeminicideRiskForm1VictimEconomicResourcesExplanation",
		"FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport",
		"FeminicideRiskForm1VictimDebtImpactDetails",
		"FeminicideRiskForm1VictimReceivedTransportSubsidy",
		"FeminicideRiskForm1VictimTransportSubsidyExplanation",
		"FeminicideRiskForm1VictimSafetyAvoidedTransport",
		"FeminicideRiskForm1VictimSafetyAvoidedExplanation",
		"FeminicideRiskForm1VictimFoodAccessFrequency",
		"FeminicideRiskForm1VictimFoodAccessExplanation",
		"FeminicideRiskForm1VictimAgressorFoodRestriction",
		"FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation",
		"FeminicideRiskForm1VictimFamilyFixedIncome",
		"FeminicideRiskForm1VictimFamilyFixedIncomeExplanation",
		"FeminicideRiskForm1VictimIsOnlyProviderForFood",
		"FeminicideRiskForm1VictimIsOnlyProviderExplanation",
		"FeminicideRiskForm1VictimJuridicalAssistanceReceived",
		"FeminicideRiskForm1VictimJuridicalAssistanceExplain",
		"FeminicideRiskForm1VictimWantsJuridicalAssistance",
		"FeminicideRiskForm1VictimRepresentationsOfVictims",
		"FeminicideRiskForm1VictimRepresentationsExplanation",
		"FeminicideRiskForm1VictimPsychosocialSupportReceived",
		"FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain",
		"FeminicideRiskForm1VictimUrgentEmotionalCrisis",
		"FeminicideRiskForm1VictimUrgentCrisisExplanation",
		"FeminicideRiskForm1AggressorSameResidence",
		"FeminicideRiskForm1AggressorSameResidenceExplanation",
		"FeminicideRiskForm1AggressorKnowsVictimLocation",
		"FeminicideRiskForm1AggressorKnowsVictimLocationExplanation",
		"FeminicideRiskForm1VictimHousingHelpReceived",
		"FeminicideRiskForm1VictimHousingHelpExplain",
		"FeminicideRiskForm1VictimAbandonClothing",
		"FeminicideRiskForm1VictimAbandonClothingExplain",
		"FeminicideRiskForm1VictimClothingHelpReceived",
		"FeminicideRiskForm1VictimClothingHelpExplain",
		"FeminicideRiskForm1VictimOtherNeeds",
		"FeminicideRiskForm1AnyAssistanceReceived",
		"FeminicideRiskForm1AnyAssistanceReceivedCityHall",
		"FeminicideRiskForm1AnyAssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AnyAssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AnyAssistanceReceivedOther",
		"FeminicideRiskForm1AssistanceReceived",
		"FeminicideRiskForm1AssistanceReceivedCityHall",
		"FeminicideRiskForm1AssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AssistanceReceivedOther",
		"FeminicideRiskForm1FamilyFather",
		"FeminicideRiskForm1FamilyMother",
		"FeminicideRiskForm1FamilyStepfather",
		"FeminicideRiskForm1FamilyStepmother",
		"FeminicideRiskForm1FamilyPartner",
		"FeminicideRiskForm1FamilySibling1",
		"FeminicideRiskForm1FamilySibling2",
		"FeminicideRiskForm1FamilySibling3",
		"FeminicideRiskForm1FamilySibling4",
		"FeminicideRiskForm1FamilySibling5",
		"FeminicideRiskForm1FamilySonDaughter1",
		"FeminicideRiskForm1FamilySonDaughter2",
		"FeminicideRiskForm1FamilySonDaughter3",
		"FeminicideRiskForm1FamilySonDaughter4",
		"FeminicideRiskForm1FamilySonDaughter5",
		"FeminicideRiskForm1FamilyGrandmother",
		"FeminicideRiskForm1FamilyGrandfather",
		"FeminicideRiskForm1FamilyOtherMember",
		"FeminicideRiskForm1InterviewDate",
		"FeminicideRiskForm1Summary",
		"FeminicideRiskForm1FeminicideRisk",
	}

	var values []interface{} = []interface{}{
		f.FeminicideRiskForm1ICode,
		f.FeminicideRiskForm1CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		f.FeminicideRiskForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		f.FeminicideRiskForm1VictimIdentityName,
		f.FeminicideRiskForm1BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		f.FeminicideRiskForm1VictimAddress,
		f.FeminicideRiskForm1VictimLivingZone.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimLivingTownCode,
		f.FeminicideRiskForm1VictimSGSSSAffiliation.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1EpsName,
		f.FeminicideRiskForm1ContactPhone,
		f.FeminicideRiskForm1ContactEmergencyContactNames,
		f.FeminicideRiskForm1ContactEmergencyContactNumber,
		f.FeminicideRiskForm1VictimMaritalStatus.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimSex.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimGenderIdentity.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimSexualOrientation.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimEthnicAffiliation.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimIndigenousPeople.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimIsMigrant.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimMigrationStatus.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimMaxEducationLevel.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimIsSpecialPopulation.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimCurrentlyHasJob.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimJobExplanation,
		f.FeminicideRiskForm1VictimAbandonedJobDueToRisk.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimMainOccupation.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimHasDisability.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimDisabilityType.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimRiskDescription,
		f.FeminicideRiskForm1VictimAdditionalInfo,
		f.FeminicideRiskForm1VictimIsEconomicProvider.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimEconomicProviderExplanation,
		f.FeminicideRiskForm1VictimHasFamiliarSupport.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimFamiliarSupportType.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimPublicTransportAccess.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimCommonTransportMode.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimTransportModeExplanation,
		f.FeminicideRiskForm1VictimEstimatedTravelCost,
		f.FeminicideRiskForm1VictimDifficultiesWithTransport.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimDifficultiesExplanation,
		f.FeminicideRiskForm1VictimEconomicResourcesForTransport.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimEconomicResourcesExplanation,
		f.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimDebtImpactDetails,
		f.FeminicideRiskForm1VictimReceivedTransportSubsidy.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimTransportSubsidyExplanation,
		f.FeminicideRiskForm1VictimSafetyAvoidedTransport.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimSafetyAvoidedExplanation,
		f.FeminicideRiskForm1VictimFoodAccessFrequency.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimFoodAccessExplanation,
		f.FeminicideRiskForm1VictimAgressorFoodRestriction.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation,
		f.FeminicideRiskForm1VictimFamilyFixedIncome.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation,
		f.FeminicideRiskForm1VictimIsOnlyProviderForFood.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimIsOnlyProviderExplanation,
		f.FeminicideRiskForm1VictimJuridicalAssistanceReceived.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimJuridicalAssistanceExplain,
		f.FeminicideRiskForm1VictimWantsJuridicalAssistance.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimRepresentationsOfVictims.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimRepresentationsExplanation,
		f.FeminicideRiskForm1VictimPsychosocialSupportReceived.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain,
		f.FeminicideRiskForm1VictimUrgentEmotionalCrisis.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimUrgentCrisisExplanation,
		f.FeminicideRiskForm1AggressorSameResidence.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AggressorSameResidenceExplanation,
		f.FeminicideRiskForm1AggressorKnowsVictimLocation.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation,
		f.FeminicideRiskForm1VictimHousingHelpReceived.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimHousingHelpExplain,
		f.FeminicideRiskForm1VictimAbandonClothing.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimAbandonClothingExplain,
		f.FeminicideRiskForm1VictimClothingHelpReceived.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1VictimClothingHelpExplain,
		f.FeminicideRiskForm1VictimOtherNeeds,
		f.FeminicideRiskForm1AnyAssistanceReceived.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AnyAssistanceReceivedCityHall.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AnyAssistanceReceivedOther.VictimCaseForm2EnumsId,
		f.FeminicideRiskForm1AssistanceReceived,
		f.FeminicideRiskForm1AssistanceReceivedCityHall,
		f.FeminicideRiskForm1AssistanceReceivedWomensOffice,
		f.FeminicideRiskForm1AssistanceReceivedOtherEntity,
		f.FeminicideRiskForm1AssistanceReceivedOther,
		f.FeminicideRiskForm1FamilyFather,
		f.FeminicideRiskForm1FamilyMother,
		f.FeminicideRiskForm1FamilyStepfather,
		f.FeminicideRiskForm1FamilyStepmother,
		f.FeminicideRiskForm1FamilyPartner,
		f.FeminicideRiskForm1FamilySibling1,
		f.FeminicideRiskForm1FamilySibling2,
		f.FeminicideRiskForm1FamilySibling3,
		f.FeminicideRiskForm1FamilySibling4,
		f.FeminicideRiskForm1FamilySibling5,
		f.FeminicideRiskForm1FamilySonDaughter1,
		f.FeminicideRiskForm1FamilySonDaughter2,
		f.FeminicideRiskForm1FamilySonDaughter3,
		f.FeminicideRiskForm1FamilySonDaughter4,
		f.FeminicideRiskForm1FamilySonDaughter5,
		f.FeminicideRiskForm1FamilyGrandmother,
		f.FeminicideRiskForm1FamilyGrandfather,
		f.FeminicideRiskForm1FamilyOtherMember,
		f.FeminicideRiskForm1InterviewDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		f.FeminicideRiskForm1Summary,
		f.FeminicideRiskForm1FeminicideRisk.(FeminicideRiskDTO).FeminicideRiskId,
	}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, fields, []string{}, FeminicideRiskForm1DBName, []string{}, []string{}, []string{"FeminicideRiskForm1Id"}, common_dao.SQL_AND, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query, values...)

	persistenceCtrl.Scan(&f.FeminicideRiskForm1Id)

	if persistenceCtrl.Error != nil {
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicideRiskForm1(by common_controllers.By, feminicideRiskForm1 *FeminicideRiskForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var path string = FeminicideRiskForm1DBScheme + "." + FeminicideRiskForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideRiskForm1Id",
		"FeminicideRiskForm1ICode",
		"FeminicideRiskForm1CreationDate",
		"FeminicideRiskForm1UpdateDate",
		"FeminicideRiskForm1VictimIdentityName",
		"FeminicideRiskForm1BirthDate",
		"FeminicideRiskForm1VictimAddress",
		"FeminicideRiskForm1VictimLivingZone",
		"FeminicideRiskForm1VictimLivingTownCode",
		"FeminicideRiskForm1VictimSGSSSAffiliation",
		"FeminicideRiskForm1EpsName",
		"FeminicideRiskForm1ContactPhone",
		"FeminicideRiskForm1ContactEmergencyContactNames",
		"FeminicideRiskForm1ContactEmergencyContactNumber",
		"FeminicideRiskForm1VictimMaritalStatus",
		"FeminicideRiskForm1VictimSex",
		"FeminicideRiskForm1VictimGenderIdentity",
		"FeminicideRiskForm1VictimSexualOrientation",
		"FeminicideRiskForm1VictimEthnicAffiliation",
		"FeminicideRiskForm1VictimIndigenousPeople",
		"FeminicideRiskForm1VictimIsMigrant",
		"FeminicideRiskForm1VictimMigrationStatus",
		"FeminicideRiskForm1VictimMaxEducationLevel",
		"FeminicideRiskForm1VictimIsSpecialPopulation",
		"FeminicideRiskForm1VictimCurrentlyHasJob",
		"FeminicideRiskForm1VictimJobExplanation",
		"FeminicideRiskForm1VictimAbandonedJobDueToRisk",
		"FeminicideRiskForm1VictimMainOccupation",
		"FeminicideRiskForm1VictimHasDisability",
		"FeminicideRiskForm1VictimDisabilityType",
		"FeminicideRiskForm1VictimRiskDescription",
		"FeminicideRiskForm1VictimAdditionalInfo",
		"FeminicideRiskForm1VictimIsEconomicProvider",
		"FeminicideRiskForm1VictimEconomicProviderExplanation",
		"FeminicideRiskForm1VictimHasFamiliarSupport",
		"FeminicideRiskForm1VictimFamiliarSupportType",
		"FeminicideRiskForm1VictimPublicTransportAccess",
		"FeminicideRiskForm1VictimCommonTransportMode",
		"FeminicideRiskForm1VictimTransportModeExplanation",
		"FeminicideRiskForm1VictimEstimatedTravelCost",
		"FeminicideRiskForm1VictimDifficultiesWithTransport",
		"FeminicideRiskForm1VictimDifficultiesExplanation",
		"FeminicideRiskForm1VictimEconomicResourcesForTransport",
		"FeminicideRiskForm1VictimEconomicResourcesExplanation",
		"FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport",
		"FeminicideRiskForm1VictimDebtImpactDetails",
		"FeminicideRiskForm1VictimReceivedTransportSubsidy",
		"FeminicideRiskForm1VictimTransportSubsidyExplanation",
		"FeminicideRiskForm1VictimSafetyAvoidedTransport",
		"FeminicideRiskForm1VictimSafetyAvoidedExplanation",
		"FeminicideRiskForm1VictimFoodAccessFrequency",
		"FeminicideRiskForm1VictimFoodAccessExplanation",
		"FeminicideRiskForm1VictimAgressorFoodRestriction",
		"FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation",
		"FeminicideRiskForm1VictimFamilyFixedIncome",
		"FeminicideRiskForm1VictimFamilyFixedIncomeExplanation",
		"FeminicideRiskForm1VictimIsOnlyProviderForFood",
		"FeminicideRiskForm1VictimIsOnlyProviderExplanation",
		"FeminicideRiskForm1VictimJuridicalAssistanceReceived",
		"FeminicideRiskForm1VictimJuridicalAssistanceExplain",
		"FeminicideRiskForm1VictimWantsJuridicalAssistance",
		"FeminicideRiskForm1VictimRepresentationsOfVictims",
		"FeminicideRiskForm1VictimRepresentationsExplanation",
		"FeminicideRiskForm1VictimPsychosocialSupportReceived",
		"FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain",
		"FeminicideRiskForm1VictimUrgentEmotionalCrisis",
		"FeminicideRiskForm1VictimUrgentCrisisExplanation",
		"FeminicideRiskForm1AggressorSameResidence",
		"FeminicideRiskForm1AggressorSameResidenceExplanation",
		"FeminicideRiskForm1AggressorKnowsVictimLocation",
		"FeminicideRiskForm1AggressorKnowsVictimLocationExplanation",
		"FeminicideRiskForm1VictimHousingHelpReceived",
		"FeminicideRiskForm1VictimHousingHelpExplain",
		"FeminicideRiskForm1VictimAbandonClothing",
		"FeminicideRiskForm1VictimAbandonClothingExplain",
		"FeminicideRiskForm1VictimClothingHelpReceived",
		"FeminicideRiskForm1VictimClothingHelpExplain",
		"FeminicideRiskForm1VictimOtherNeeds",
		"FeminicideRiskForm1AnyAssistanceReceived",
		"FeminicideRiskForm1AnyAssistanceReceivedCityHall",
		"FeminicideRiskForm1AnyAssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AnyAssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AnyAssistanceReceivedOther",
		"FeminicideRiskForm1AssistanceReceived",
		"FeminicideRiskForm1AssistanceReceivedCityHall",
		"FeminicideRiskForm1AssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AssistanceReceivedOther",
		"FeminicideRiskForm1FamilyFather",
		"FeminicideRiskForm1FamilyMother",
		"FeminicideRiskForm1FamilyStepfather",
		"FeminicideRiskForm1FamilyStepmother",
		"FeminicideRiskForm1FamilyPartner",
		"FeminicideRiskForm1FamilySibling1",
		"FeminicideRiskForm1FamilySibling2",
		"FeminicideRiskForm1FamilySibling3",
		"FeminicideRiskForm1FamilySibling4",
		"FeminicideRiskForm1FamilySibling5",
		"FeminicideRiskForm1FamilySonDaughter1",
		"FeminicideRiskForm1FamilySonDaughter2",
		"FeminicideRiskForm1FamilySonDaughter3",
		"FeminicideRiskForm1FamilySonDaughter4",
		"FeminicideRiskForm1FamilySonDaughter5",
		"FeminicideRiskForm1FamilyGrandmother",
		"FeminicideRiskForm1FamilyGrandfather",
		"FeminicideRiskForm1FamilyOtherMember",
		"FeminicideRiskForm1InterviewDate",
		"FeminicideRiskForm1Summary",
		"FeminicideRiskForm1FeminicideRisk",
	}

	var fieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, []string{}, FeminicideRiskForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideRiskForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var pgDB FeminicideRiskForm1PgDB
	persistenceCtrl.Scan(
		&pgDB.FeminicideRiskForm1Id,
		&pgDB.FeminicideRiskForm1ICode,
		&pgDB.FeminicideRiskForm1CreationDate,
		&pgDB.FeminicideRiskForm1UpdateDate,
		&pgDB.FeminicideRiskForm1VictimIdentityName,
		&pgDB.FeminicideRiskForm1BirthDate,
		&pgDB.FeminicideRiskForm1VictimAddress,
		&pgDB.FeminicideRiskForm1VictimLivingZone,
		&pgDB.FeminicideRiskForm1VictimLivingTownCode,
		&pgDB.FeminicideRiskForm1VictimSGSSSAffiliation,
		&pgDB.FeminicideRiskForm1EpsName,
		&pgDB.FeminicideRiskForm1ContactPhone,
		&pgDB.FeminicideRiskForm1ContactEmergencyContactNames,
		&pgDB.FeminicideRiskForm1ContactEmergencyContactNumber,
		&pgDB.FeminicideRiskForm1VictimMaritalStatus,
		&pgDB.FeminicideRiskForm1VictimSex,
		&pgDB.FeminicideRiskForm1VictimGenderIdentity,
		&pgDB.FeminicideRiskForm1VictimSexualOrientation,
		&pgDB.FeminicideRiskForm1VictimEthnicAffiliation,
		&pgDB.FeminicideRiskForm1VictimIndigenousPeople,
		&pgDB.FeminicideRiskForm1VictimIsMigrant,
		&pgDB.FeminicideRiskForm1VictimMigrationStatus,
		&pgDB.FeminicideRiskForm1VictimMaxEducationLevel,
		&pgDB.FeminicideRiskForm1VictimIsSpecialPopulation,
		&pgDB.FeminicideRiskForm1VictimCurrentlyHasJob,
		&pgDB.FeminicideRiskForm1VictimJobExplanation,
		&pgDB.FeminicideRiskForm1VictimAbandonedJobDueToRisk,
		&pgDB.FeminicideRiskForm1VictimMainOccupation,
		&pgDB.FeminicideRiskForm1VictimHasDisability,
		&pgDB.FeminicideRiskForm1VictimDisabilityType,
		&pgDB.FeminicideRiskForm1VictimRiskDescription,
		&pgDB.FeminicideRiskForm1VictimAdditionalInfo,
		&pgDB.FeminicideRiskForm1VictimIsEconomicProvider,
		&pgDB.FeminicideRiskForm1VictimEconomicProviderExplanation,
		&pgDB.FeminicideRiskForm1VictimHasFamiliarSupport,
		&pgDB.FeminicideRiskForm1VictimFamiliarSupportType,
		&pgDB.FeminicideRiskForm1VictimPublicTransportAccess,
		&pgDB.FeminicideRiskForm1VictimCommonTransportMode,
		&pgDB.FeminicideRiskForm1VictimTransportModeExplanation,
		&pgDB.FeminicideRiskForm1VictimEstimatedTravelCost,
		&pgDB.FeminicideRiskForm1VictimDifficultiesWithTransport,
		&pgDB.FeminicideRiskForm1VictimDifficultiesExplanation,
		&pgDB.FeminicideRiskForm1VictimEconomicResourcesForTransport,
		&pgDB.FeminicideRiskForm1VictimEconomicResourcesExplanation,
		&pgDB.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport,
		&pgDB.FeminicideRiskForm1VictimDebtImpactDetails,
		&pgDB.FeminicideRiskForm1VictimReceivedTransportSubsidy,
		&pgDB.FeminicideRiskForm1VictimTransportSubsidyExplanation,
		&pgDB.FeminicideRiskForm1VictimSafetyAvoidedTransport,
		&pgDB.FeminicideRiskForm1VictimSafetyAvoidedExplanation,
		&pgDB.FeminicideRiskForm1VictimFoodAccessFrequency,
		&pgDB.FeminicideRiskForm1VictimFoodAccessExplanation,
		&pgDB.FeminicideRiskForm1VictimAgressorFoodRestriction,
		&pgDB.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation,
		&pgDB.FeminicideRiskForm1VictimFamilyFixedIncome,
		&pgDB.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation,
		&pgDB.FeminicideRiskForm1VictimIsOnlyProviderForFood,
		&pgDB.FeminicideRiskForm1VictimIsOnlyProviderExplanation,
		&pgDB.FeminicideRiskForm1VictimJuridicalAssistanceReceived,
		&pgDB.FeminicideRiskForm1VictimJuridicalAssistanceExplain,
		&pgDB.FeminicideRiskForm1VictimWantsJuridicalAssistance,
		&pgDB.FeminicideRiskForm1VictimRepresentationsOfVictims,
		&pgDB.FeminicideRiskForm1VictimRepresentationsExplanation,
		&pgDB.FeminicideRiskForm1VictimPsychosocialSupportReceived,
		&pgDB.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain,
		&pgDB.FeminicideRiskForm1VictimUrgentEmotionalCrisis,
		&pgDB.FeminicideRiskForm1VictimUrgentCrisisExplanation,
		&pgDB.FeminicideRiskForm1AggressorSameResidence,
		&pgDB.FeminicideRiskForm1AggressorSameResidenceExplanation,
		&pgDB.FeminicideRiskForm1AggressorKnowsVictimLocation,
		&pgDB.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation,
		&pgDB.FeminicideRiskForm1VictimHousingHelpReceived,
		&pgDB.FeminicideRiskForm1VictimHousingHelpExplain,
		&pgDB.FeminicideRiskForm1VictimAbandonClothing,
		&pgDB.FeminicideRiskForm1VictimAbandonClothingExplain,
		&pgDB.FeminicideRiskForm1VictimClothingHelpReceived,
		&pgDB.FeminicideRiskForm1VictimClothingHelpExplain,
		&pgDB.FeminicideRiskForm1VictimOtherNeeds,
		&pgDB.FeminicideRiskForm1AnyAssistanceReceived,
		&pgDB.FeminicideRiskForm1AnyAssistanceReceivedCityHall,
		&pgDB.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice,
		&pgDB.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity,
		&pgDB.FeminicideRiskForm1AnyAssistanceReceivedOther,
		&pgDB.FeminicideRiskForm1AssistanceReceived,
		&pgDB.FeminicideRiskForm1AssistanceReceivedCityHall,
		&pgDB.FeminicideRiskForm1AssistanceReceivedWomensOffice,
		&pgDB.FeminicideRiskForm1AssistanceReceivedOtherEntity,
		&pgDB.FeminicideRiskForm1AssistanceReceivedOther,
		&pgDB.FeminicideRiskForm1FamilyFather,
		&pgDB.FeminicideRiskForm1FamilyMother,
		&pgDB.FeminicideRiskForm1FamilyStepfather,
		&pgDB.FeminicideRiskForm1FamilyStepmother,
		&pgDB.FeminicideRiskForm1FamilyPartner,
		&pgDB.FeminicideRiskForm1FamilySibling1,
		&pgDB.FeminicideRiskForm1FamilySibling2,
		&pgDB.FeminicideRiskForm1FamilySibling3,
		&pgDB.FeminicideRiskForm1FamilySibling4,
		&pgDB.FeminicideRiskForm1FamilySibling5,
		&pgDB.FeminicideRiskForm1FamilySonDaughter1,
		&pgDB.FeminicideRiskForm1FamilySonDaughter2,
		&pgDB.FeminicideRiskForm1FamilySonDaughter3,
		&pgDB.FeminicideRiskForm1FamilySonDaughter4,
		&pgDB.FeminicideRiskForm1FamilySonDaughter5,
		&pgDB.FeminicideRiskForm1FamilyGrandmother,
		&pgDB.FeminicideRiskForm1FamilyGrandfather,
		&pgDB.FeminicideRiskForm1FamilyOtherMember,
		&pgDB.FeminicideRiskForm1InterviewDate,
		&pgDB.FeminicideRiskForm1Summary,
		&pgDB.FeminicideRiskForm1FeminicideRisk)

	*feminicideRiskForm1 = pgDB.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicideRisksForm1(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideRiskForm1DTO, int, error) {
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var path string = FeminicideRiskForm1DBScheme + "." + FeminicideRiskForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideRiskForm1Id",
		"FeminicideRiskForm1ICode",
		"FeminicideRiskForm1CreationDate",
		"FeminicideRiskForm1UpdateDate",
		"FeminicideRiskForm1VictimIdentityName",
		"FeminicideRiskForm1BirthDate",
		"FeminicideRiskForm1VictimAddress",
		"FeminicideRiskForm1VictimLivingZone",
		"FeminicideRiskForm1VictimLivingTownCode",
		"FeminicideRiskForm1VictimSGSSSAffiliation",
		"FeminicideRiskForm1EpsName",
		"FeminicideRiskForm1ContactPhone",
		"FeminicideRiskForm1ContactEmergencyContactNames",
		"FeminicideRiskForm1ContactEmergencyContactNumber",
		"FeminicideRiskForm1VictimMaritalStatus",
		"FeminicideRiskForm1VictimSex",
		"FeminicideRiskForm1VictimGenderIdentity",
		"FeminicideRiskForm1VictimSexualOrientation",
		"FeminicideRiskForm1VictimEthnicAffiliation",
		"FeminicideRiskForm1VictimIndigenousPeople",
		"FeminicideRiskForm1VictimIsMigrant",
		"FeminicideRiskForm1VictimMigrationStatus",
		"FeminicideRiskForm1VictimMaxEducationLevel",
		"FeminicideRiskForm1VictimIsSpecialPopulation",
		"FeminicideRiskForm1VictimCurrentlyHasJob",
		"FeminicideRiskForm1VictimJobExplanation",
		"FeminicideRiskForm1VictimAbandonedJobDueToRisk",
		"FeminicideRiskForm1VictimMainOccupation",
		"FeminicideRiskForm1VictimHasDisability",
		"FeminicideRiskForm1VictimDisabilityType",
		"FeminicideRiskForm1VictimRiskDescription",
		"FeminicideRiskForm1VictimAdditionalInfo",
		"FeminicideRiskForm1VictimIsEconomicProvider",
		"FeminicideRiskForm1VictimEconomicProviderExplanation",
		"FeminicideRiskForm1VictimHasFamiliarSupport",
		"FeminicideRiskForm1VictimFamiliarSupportType",
		"FeminicideRiskForm1VictimPublicTransportAccess",
		"FeminicideRiskForm1VictimCommonTransportMode",
		"FeminicideRiskForm1VictimTransportModeExplanation",
		"FeminicideRiskForm1VictimEstimatedTravelCost",
		"FeminicideRiskForm1VictimDifficultiesWithTransport",
		"FeminicideRiskForm1VictimDifficultiesExplanation",
		"FeminicideRiskForm1VictimEconomicResourcesForTransport",
		"FeminicideRiskForm1VictimEconomicResourcesExplanation",
		"FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport",
		"FeminicideRiskForm1VictimDebtImpactDetails",
		"FeminicideRiskForm1VictimReceivedTransportSubsidy",
		"FeminicideRiskForm1VictimTransportSubsidyExplanation",
		"FeminicideRiskForm1VictimSafetyAvoidedTransport",
		"FeminicideRiskForm1VictimSafetyAvoidedExplanation",
		"FeminicideRiskForm1VictimFoodAccessFrequency",
		"FeminicideRiskForm1VictimFoodAccessExplanation",
		"FeminicideRiskForm1VictimAgressorFoodRestriction",
		"FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation",
		"FeminicideRiskForm1VictimFamilyFixedIncome",
		"FeminicideRiskForm1VictimFamilyFixedIncomeExplanation",
		"FeminicideRiskForm1VictimIsOnlyProviderForFood",
		"FeminicideRiskForm1VictimIsOnlyProviderExplanation",
		"FeminicideRiskForm1VictimJuridicalAssistanceReceived",
		"FeminicideRiskForm1VictimJuridicalAssistanceExplain",
		"FeminicideRiskForm1VictimWantsJuridicalAssistance",
		"FeminicideRiskForm1VictimRepresentationsOfVictims",
		"FeminicideRiskForm1VictimRepresentationsExplanation",
		"FeminicideRiskForm1VictimPsychosocialSupportReceived",
		"FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain",
		"FeminicideRiskForm1VictimUrgentEmotionalCrisis",
		"FeminicideRiskForm1VictimUrgentCrisisExplanation",
		"FeminicideRiskForm1AggressorSameResidence",
		"FeminicideRiskForm1AggressorSameResidenceExplanation",
		"FeminicideRiskForm1AggressorKnowsVictimLocation",
		"FeminicideRiskForm1AggressorKnowsVictimLocationExplanation",
		"FeminicideRiskForm1VictimHousingHelpReceived",
		"FeminicideRiskForm1VictimHousingHelpExplain",
		"FeminicideRiskForm1VictimAbandonClothing",
		"FeminicideRiskForm1VictimAbandonClothingExplain",
		"FeminicideRiskForm1VictimClothingHelpReceived",
		"FeminicideRiskForm1VictimClothingHelpExplain",
		"FeminicideRiskForm1VictimOtherNeeds",
		"FeminicideRiskForm1AnyAssistanceReceived",
		"FeminicideRiskForm1AnyAssistanceReceivedCityHall",
		"FeminicideRiskForm1AnyAssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AnyAssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AnyAssistanceReceivedOther",
		"FeminicideRiskForm1AssistanceReceived",
		"FeminicideRiskForm1AssistanceReceivedCityHall",
		"FeminicideRiskForm1AssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AssistanceReceivedOther",
		"FeminicideRiskForm1FamilyFather",
		"FeminicideRiskForm1FamilyMother",
		"FeminicideRiskForm1FamilyStepfather",
		"FeminicideRiskForm1FamilyStepmother",
		"FeminicideRiskForm1FamilyPartner",
		"FeminicideRiskForm1FamilySibling1",
		"FeminicideRiskForm1FamilySibling2",
		"FeminicideRiskForm1FamilySibling3",
		"FeminicideRiskForm1FamilySibling4",
		"FeminicideRiskForm1FamilySibling5",
		"FeminicideRiskForm1FamilySonDaughter1",
		"FeminicideRiskForm1FamilySonDaughter2",
		"FeminicideRiskForm1FamilySonDaughter3",
		"FeminicideRiskForm1FamilySonDaughter4",
		"FeminicideRiskForm1FamilySonDaughter5",
		"FeminicideRiskForm1FamilyGrandmother",
		"FeminicideRiskForm1FamilyGrandfather",
		"FeminicideRiskForm1FamilyOtherMember",
		"FeminicideRiskForm1InterviewDate",
		"FeminicideRiskForm1Summary",
		"FeminicideRiskForm1FeminicideRisk",
	}

	var fieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, []string{}, FeminicideRiskForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideRiskForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, true) +
		` ORDER BY ` + path + `.` + FeminicideRiskForm1FieldDefinitions["FeminicideRiskForm1CreationDate"].DBName + ` ASC ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	var list []FeminicideRiskForm1DTO
	for persistenceCtrl.Next() {
		var pgDB FeminicideRiskForm1PgDB
		persistenceCtrl.ScanRow(
			&pgDB.FeminicideRiskForm1Id,
			&pgDB.FeminicideRiskForm1ICode,
			&pgDB.FeminicideRiskForm1CreationDate,
			&pgDB.FeminicideRiskForm1UpdateDate,
			&pgDB.FeminicideRiskForm1VictimIdentityName,
			&pgDB.FeminicideRiskForm1BirthDate,
			&pgDB.FeminicideRiskForm1VictimAddress,
			&pgDB.FeminicideRiskForm1VictimLivingZone,
			&pgDB.FeminicideRiskForm1VictimLivingTownCode,
			&pgDB.FeminicideRiskForm1VictimSGSSSAffiliation,
			&pgDB.FeminicideRiskForm1EpsName,
			&pgDB.FeminicideRiskForm1ContactPhone,
			&pgDB.FeminicideRiskForm1ContactEmergencyContactNames,
			&pgDB.FeminicideRiskForm1ContactEmergencyContactNumber,
			&pgDB.FeminicideRiskForm1VictimMaritalStatus,
			&pgDB.FeminicideRiskForm1VictimSex,
			&pgDB.FeminicideRiskForm1VictimGenderIdentity,
			&pgDB.FeminicideRiskForm1VictimSexualOrientation,
			&pgDB.FeminicideRiskForm1VictimEthnicAffiliation,
			&pgDB.FeminicideRiskForm1VictimIndigenousPeople,
			&pgDB.FeminicideRiskForm1VictimIsMigrant,
			&pgDB.FeminicideRiskForm1VictimMigrationStatus,
			&pgDB.FeminicideRiskForm1VictimMaxEducationLevel,
			&pgDB.FeminicideRiskForm1VictimIsSpecialPopulation,
			&pgDB.FeminicideRiskForm1VictimCurrentlyHasJob,
			&pgDB.FeminicideRiskForm1VictimJobExplanation,
			&pgDB.FeminicideRiskForm1VictimAbandonedJobDueToRisk,
			&pgDB.FeminicideRiskForm1VictimMainOccupation,
			&pgDB.FeminicideRiskForm1VictimHasDisability,
			&pgDB.FeminicideRiskForm1VictimDisabilityType,
			&pgDB.FeminicideRiskForm1VictimRiskDescription,
			&pgDB.FeminicideRiskForm1VictimAdditionalInfo,
			&pgDB.FeminicideRiskForm1VictimIsEconomicProvider,
			&pgDB.FeminicideRiskForm1VictimEconomicProviderExplanation,
			&pgDB.FeminicideRiskForm1VictimHasFamiliarSupport,
			&pgDB.FeminicideRiskForm1VictimFamiliarSupportType,
			&pgDB.FeminicideRiskForm1VictimPublicTransportAccess,
			&pgDB.FeminicideRiskForm1VictimCommonTransportMode,
			&pgDB.FeminicideRiskForm1VictimTransportModeExplanation,
			&pgDB.FeminicideRiskForm1VictimEstimatedTravelCost,
			&pgDB.FeminicideRiskForm1VictimDifficultiesWithTransport,
			&pgDB.FeminicideRiskForm1VictimDifficultiesExplanation,
			&pgDB.FeminicideRiskForm1VictimEconomicResourcesForTransport,
			&pgDB.FeminicideRiskForm1VictimEconomicResourcesExplanation,
			&pgDB.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport,
			&pgDB.FeminicideRiskForm1VictimDebtImpactDetails,
			&pgDB.FeminicideRiskForm1VictimReceivedTransportSubsidy,
			&pgDB.FeminicideRiskForm1VictimTransportSubsidyExplanation,
			&pgDB.FeminicideRiskForm1VictimSafetyAvoidedTransport,
			&pgDB.FeminicideRiskForm1VictimSafetyAvoidedExplanation,
			&pgDB.FeminicideRiskForm1VictimFoodAccessFrequency,
			&pgDB.FeminicideRiskForm1VictimFoodAccessExplanation,
			&pgDB.FeminicideRiskForm1VictimAgressorFoodRestriction,
			&pgDB.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation,
			&pgDB.FeminicideRiskForm1VictimFamilyFixedIncome,
			&pgDB.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation,
			&pgDB.FeminicideRiskForm1VictimIsOnlyProviderForFood,
			&pgDB.FeminicideRiskForm1VictimIsOnlyProviderExplanation,
			&pgDB.FeminicideRiskForm1VictimJuridicalAssistanceReceived,
			&pgDB.FeminicideRiskForm1VictimJuridicalAssistanceExplain,
			&pgDB.FeminicideRiskForm1VictimWantsJuridicalAssistance,
			&pgDB.FeminicideRiskForm1VictimRepresentationsOfVictims,
			&pgDB.FeminicideRiskForm1VictimRepresentationsExplanation,
			&pgDB.FeminicideRiskForm1VictimPsychosocialSupportReceived,
			&pgDB.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain,
			&pgDB.FeminicideRiskForm1VictimUrgentEmotionalCrisis,
			&pgDB.FeminicideRiskForm1VictimUrgentCrisisExplanation,
			&pgDB.FeminicideRiskForm1AggressorSameResidence,
			&pgDB.FeminicideRiskForm1AggressorSameResidenceExplanation,
			&pgDB.FeminicideRiskForm1AggressorKnowsVictimLocation,
			&pgDB.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation,
			&pgDB.FeminicideRiskForm1VictimHousingHelpReceived,
			&pgDB.FeminicideRiskForm1VictimHousingHelpExplain,
			&pgDB.FeminicideRiskForm1VictimAbandonClothing,
			&pgDB.FeminicideRiskForm1VictimAbandonClothingExplain,
			&pgDB.FeminicideRiskForm1VictimClothingHelpReceived,
			&pgDB.FeminicideRiskForm1VictimClothingHelpExplain,
			&pgDB.FeminicideRiskForm1VictimOtherNeeds,
			&pgDB.FeminicideRiskForm1AnyAssistanceReceived,
			&pgDB.FeminicideRiskForm1AnyAssistanceReceivedCityHall,
			&pgDB.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice,
			&pgDB.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity,
			&pgDB.FeminicideRiskForm1AnyAssistanceReceivedOther,
			&pgDB.FeminicideRiskForm1AssistanceReceived,
			&pgDB.FeminicideRiskForm1AssistanceReceivedCityHall,
			&pgDB.FeminicideRiskForm1AssistanceReceivedWomensOffice,
			&pgDB.FeminicideRiskForm1AssistanceReceivedOtherEntity,
			&pgDB.FeminicideRiskForm1AssistanceReceivedOther,
			&pgDB.FeminicideRiskForm1FamilyFather,
			&pgDB.FeminicideRiskForm1FamilyMother,
			&pgDB.FeminicideRiskForm1FamilyStepfather,
			&pgDB.FeminicideRiskForm1FamilyStepmother,
			&pgDB.FeminicideRiskForm1FamilyPartner,
			&pgDB.FeminicideRiskForm1FamilySibling1,
			&pgDB.FeminicideRiskForm1FamilySibling2,
			&pgDB.FeminicideRiskForm1FamilySibling3,
			&pgDB.FeminicideRiskForm1FamilySibling4,
			&pgDB.FeminicideRiskForm1FamilySibling5,
			&pgDB.FeminicideRiskForm1FamilySonDaughter1,
			&pgDB.FeminicideRiskForm1FamilySonDaughter2,
			&pgDB.FeminicideRiskForm1FamilySonDaughter3,
			&pgDB.FeminicideRiskForm1FamilySonDaughter4,
			&pgDB.FeminicideRiskForm1FamilySonDaughter5,
			&pgDB.FeminicideRiskForm1FamilyGrandmother,
			&pgDB.FeminicideRiskForm1FamilyGrandfather,
			&pgDB.FeminicideRiskForm1FamilyOtherMember,
			&pgDB.FeminicideRiskForm1InterviewDate,
			&pgDB.FeminicideRiskForm1Summary,
			&pgDB.FeminicideRiskForm1FeminicideRisk)

		list = append(list, pgDB.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	if page == 0 {
		var countQuery string = `SELECT COUNT(*) FROM ` + path +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideRiskForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, true)
		persistenceCtrl.QueryRow(context.Background(), countQuery, by.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return list, count, nil
}

func UpdateFeminicideRiskForm1(feminicideRiskForm1 *FeminicideRiskForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideRiskForm1UpdateDate",
		"FeminicideRiskForm1VictimIdentityName",
		"FeminicideRiskForm1BirthDate",
		"FeminicideRiskForm1VictimAddress",
		"FeminicideRiskForm1VictimLivingZone",
		"FeminicideRiskForm1VictimLivingTownCode",
		"FeminicideRiskForm1VictimSGSSSAffiliation",
		"FeminicideRiskForm1EpsName",
		"FeminicideRiskForm1ContactPhone",
		"FeminicideRiskForm1ContactEmergencyContactNames",
		"FeminicideRiskForm1ContactEmergencyContactNumber",
		"FeminicideRiskForm1VictimMaritalStatus",
		"FeminicideRiskForm1VictimSex",
		"FeminicideRiskForm1VictimGenderIdentity",
		"FeminicideRiskForm1VictimSexualOrientation",
		"FeminicideRiskForm1VictimEthnicAffiliation",
		"FeminicideRiskForm1VictimIndigenousPeople",
		"FeminicideRiskForm1VictimIsMigrant",
		"FeminicideRiskForm1VictimMigrationStatus",
		"FeminicideRiskForm1VictimMaxEducationLevel",
		"FeminicideRiskForm1VictimIsSpecialPopulation",
		"FeminicideRiskForm1VictimCurrentlyHasJob",
		"FeminicideRiskForm1VictimJobExplanation",
		"FeminicideRiskForm1VictimAbandonedJobDueToRisk",
		"FeminicideRiskForm1VictimMainOccupation",
		"FeminicideRiskForm1VictimHasDisability",
		"FeminicideRiskForm1VictimDisabilityType",
		"FeminicideRiskForm1VictimRiskDescription",
		"FeminicideRiskForm1VictimAdditionalInfo",
		"FeminicideRiskForm1VictimIsEconomicProvider",
		"FeminicideRiskForm1VictimEconomicProviderExplanation",
		"FeminicideRiskForm1VictimHasFamiliarSupport",
		"FeminicideRiskForm1VictimFamiliarSupportType",
		"FeminicideRiskForm1VictimPublicTransportAccess",
		"FeminicideRiskForm1VictimCommonTransportMode",
		"FeminicideRiskForm1VictimTransportModeExplanation",
		"FeminicideRiskForm1VictimEstimatedTravelCost",
		"FeminicideRiskForm1VictimDifficultiesWithTransport",
		"FeminicideRiskForm1VictimDifficultiesExplanation",
		"FeminicideRiskForm1VictimEconomicResourcesForTransport",
		"FeminicideRiskForm1VictimEconomicResourcesExplanation",
		"FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport",
		"FeminicideRiskForm1VictimDebtImpactDetails",
		"FeminicideRiskForm1VictimReceivedTransportSubsidy",
		"FeminicideRiskForm1VictimTransportSubsidyExplanation",
		"FeminicideRiskForm1VictimSafetyAvoidedTransport",
		"FeminicideRiskForm1VictimSafetyAvoidedExplanation",
		"FeminicideRiskForm1VictimFoodAccessFrequency",
		"FeminicideRiskForm1VictimFoodAccessExplanation",
		"FeminicideRiskForm1VictimAgressorFoodRestriction",
		"FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation",
		"FeminicideRiskForm1VictimFamilyFixedIncome",
		"FeminicideRiskForm1VictimFamilyFixedIncomeExplanation",
		"FeminicideRiskForm1VictimIsOnlyProviderForFood",
		"FeminicideRiskForm1VictimIsOnlyProviderExplanation",
		"FeminicideRiskForm1VictimJuridicalAssistanceReceived",
		"FeminicideRiskForm1VictimJuridicalAssistanceExplain",
		"FeminicideRiskForm1VictimWantsJuridicalAssistance",
		"FeminicideRiskForm1VictimRepresentationsOfVictims",
		"FeminicideRiskForm1VictimRepresentationsExplanation",
		"FeminicideRiskForm1VictimPsychosocialSupportReceived",
		"FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain",
		"FeminicideRiskForm1VictimUrgentEmotionalCrisis",
		"FeminicideRiskForm1VictimUrgentCrisisExplanation",
		"FeminicideRiskForm1AggressorSameResidence",
		"FeminicideRiskForm1AggressorSameResidenceExplanation",
		"FeminicideRiskForm1AggressorKnowsVictimLocation",
		"FeminicideRiskForm1AggressorKnowsVictimLocationExplanation",
		"FeminicideRiskForm1VictimHousingHelpReceived",
		"FeminicideRiskForm1VictimHousingHelpExplain",
		"FeminicideRiskForm1VictimAbandonClothing",
		"FeminicideRiskForm1VictimAbandonClothingExplain",
		"FeminicideRiskForm1VictimClothingHelpReceived",
		"FeminicideRiskForm1VictimClothingHelpExplain",
		"FeminicideRiskForm1VictimOtherNeeds",
		"FeminicideRiskForm1AnyAssistanceReceived",
		"FeminicideRiskForm1AnyAssistanceReceivedCityHall",
		"FeminicideRiskForm1AnyAssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AnyAssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AnyAssistanceReceivedOther",
		"FeminicideRiskForm1AssistanceReceived",
		"FeminicideRiskForm1AssistanceReceivedCityHall",
		"FeminicideRiskForm1AssistanceReceivedWomensOffice",
		"FeminicideRiskForm1AssistanceReceivedOtherEntity",
		"FeminicideRiskForm1AssistanceReceivedOther",
		"FeminicideRiskForm1FamilyFather",
		"FeminicideRiskForm1FamilyMother",
		"FeminicideRiskForm1FamilyStepfather",
		"FeminicideRiskForm1FamilyStepmother",
		"FeminicideRiskForm1FamilyPartner",
		"FeminicideRiskForm1FamilySibling1",
		"FeminicideRiskForm1FamilySibling2",
		"FeminicideRiskForm1FamilySibling3",
		"FeminicideRiskForm1FamilySibling4",
		"FeminicideRiskForm1FamilySibling5",
		"FeminicideRiskForm1FamilySonDaughter1",
		"FeminicideRiskForm1FamilySonDaughter2",
		"FeminicideRiskForm1FamilySonDaughter3",
		"FeminicideRiskForm1FamilySonDaughter4",
		"FeminicideRiskForm1FamilySonDaughter5",
		"FeminicideRiskForm1FamilyGrandmother",
		"FeminicideRiskForm1FamilyGrandfather",
		"FeminicideRiskForm1FamilyOtherMember",
		"FeminicideRiskForm1InterviewDate",
		"FeminicideRiskForm1Summary",
	}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, fieldsSlice, []string{}, FeminicideRiskForm1DBName, []string{"FeminicideRiskForm1Id"}, []string{}, []string{}, common_dao.SQL_AND, FeminicideRiskForm1DBScheme, FeminicideRiskForm1FieldDefinitions, false)
	persistenceCtrl.Exec(context.Background(), query,
		feminicideRiskForm1.FeminicideRiskForm1Id,
		feminicideRiskForm1.FeminicideRiskForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideRiskForm1.FeminicideRiskForm1VictimIdentityName,
		feminicideRiskForm1.FeminicideRiskForm1BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideRiskForm1.FeminicideRiskForm1VictimAddress,
		feminicideRiskForm1.FeminicideRiskForm1VictimLivingZone.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimLivingTownCode,
		feminicideRiskForm1.FeminicideRiskForm1VictimSGSSSAffiliation.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1EpsName,
		feminicideRiskForm1.FeminicideRiskForm1ContactPhone,
		feminicideRiskForm1.FeminicideRiskForm1ContactEmergencyContactNames,
		feminicideRiskForm1.FeminicideRiskForm1ContactEmergencyContactNumber,
		feminicideRiskForm1.FeminicideRiskForm1VictimMaritalStatus.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimSex.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimGenderIdentity.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimSexualOrientation.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimEthnicAffiliation.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimIndigenousPeople.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimIsMigrant.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimMigrationStatus.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimMaxEducationLevel.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimIsSpecialPopulation.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimCurrentlyHasJob.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimJobExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimAbandonedJobDueToRisk.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimMainOccupation.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimHasDisability.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimDisabilityType.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimRiskDescription,
		feminicideRiskForm1.FeminicideRiskForm1VictimAdditionalInfo,
		feminicideRiskForm1.FeminicideRiskForm1VictimIsEconomicProvider.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimEconomicProviderExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimHasFamiliarSupport.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimFamiliarSupportType.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimPublicTransportAccess.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimCommonTransportMode.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimTransportModeExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimEstimatedTravelCost,
		feminicideRiskForm1.FeminicideRiskForm1VictimDifficultiesWithTransport.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimDifficultiesExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimEconomicResourcesForTransport.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimEconomicResourcesExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimDebtImpactDetails,
		feminicideRiskForm1.FeminicideRiskForm1VictimReceivedTransportSubsidy.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimTransportSubsidyExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimSafetyAvoidedTransport.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimSafetyAvoidedExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimFoodAccessFrequency.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimFoodAccessExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimAgressorFoodRestriction.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimFamilyFixedIncome.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimIsOnlyProviderForFood.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimIsOnlyProviderExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimJuridicalAssistanceReceived.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimJuridicalAssistanceExplain,
		feminicideRiskForm1.FeminicideRiskForm1VictimWantsJuridicalAssistance.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimRepresentationsOfVictims.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimRepresentationsExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimPsychosocialSupportReceived.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain,
		feminicideRiskForm1.FeminicideRiskForm1VictimUrgentEmotionalCrisis.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimUrgentCrisisExplanation,
		feminicideRiskForm1.FeminicideRiskForm1AggressorSameResidence.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1AggressorSameResidenceExplanation,
		feminicideRiskForm1.FeminicideRiskForm1AggressorKnowsVictimLocation.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation,
		feminicideRiskForm1.FeminicideRiskForm1VictimHousingHelpReceived.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimHousingHelpExplain,
		feminicideRiskForm1.FeminicideRiskForm1VictimAbandonClothing.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimAbandonClothingExplain,
		feminicideRiskForm1.FeminicideRiskForm1VictimClothingHelpReceived.VictimCaseForm2EnumsId,
		feminicideRiskForm1.FeminicideRiskForm1VictimClothingHelpExplain,
		feminicideRiskForm1.FeminicideRiskForm1VictimOtherNeeds,
		feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceived,
		feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedCityHall,
		feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice,
		feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity,
		feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedOther,
		feminicideRiskForm1.FeminicideRiskForm1AssistanceReceived,
		feminicideRiskForm1.FeminicideRiskForm1AssistanceReceivedCityHall,
		feminicideRiskForm1.FeminicideRiskForm1AssistanceReceivedWomensOffice,
		feminicideRiskForm1.FeminicideRiskForm1AssistanceReceivedOtherEntity,
		feminicideRiskForm1.FeminicideRiskForm1AssistanceReceivedOther,
		feminicideRiskForm1.FeminicideRiskForm1FamilyFather,
		feminicideRiskForm1.FeminicideRiskForm1FamilyMother,
		feminicideRiskForm1.FeminicideRiskForm1FamilyStepfather,
		feminicideRiskForm1.FeminicideRiskForm1FamilyStepmother,
		feminicideRiskForm1.FeminicideRiskForm1FamilyPartner,
		feminicideRiskForm1.FeminicideRiskForm1FamilySibling1,
		feminicideRiskForm1.FeminicideRiskForm1FamilySibling2,
		feminicideRiskForm1.FeminicideRiskForm1FamilySibling3,
		feminicideRiskForm1.FeminicideRiskForm1FamilySibling4,
		feminicideRiskForm1.FeminicideRiskForm1FamilySibling5,
		feminicideRiskForm1.FeminicideRiskForm1FamilySonDaughter1,
		feminicideRiskForm1.FeminicideRiskForm1FamilySonDaughter2,
		feminicideRiskForm1.FeminicideRiskForm1FamilySonDaughter3,
		feminicideRiskForm1.FeminicideRiskForm1FamilySonDaughter4,
		feminicideRiskForm1.FeminicideRiskForm1FamilySonDaughter5,
		feminicideRiskForm1.FeminicideRiskForm1FamilyGrandmother,
		feminicideRiskForm1.FeminicideRiskForm1FamilyGrandfather,
		feminicideRiskForm1.FeminicideRiskForm1FamilyOtherMember,
		feminicideRiskForm1.FeminicideRiskForm1InterviewDate.Format(common_config.DateTime.DB_DATE_FORMAT),
		feminicideRiskForm1.FeminicideRiskForm1Summary)

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

func SetFeminicideRiskForm1Defaults(f *FeminicideRiskForm1DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		f.FeminicideRiskForm1CreationDate = time.Now()
		f.FeminicideRiskForm1UpdateDate = time.Now()
		f.FeminicideRiskForm1ICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		f.FeminicideRiskForm1UpdateDate = time.Now()
	}
}

func (obj *FeminicideRiskForm1PgDB) ToDTO() FeminicideRiskForm1DTO {
	var dto FeminicideRiskForm1DTO

	if obj.FeminicideRiskForm1Id.Valid {
		dto.FeminicideRiskForm1Id = uint64(obj.FeminicideRiskForm1Id.Int64)
	}

	if obj.FeminicideRiskForm1ICode.Valid {
		dto.FeminicideRiskForm1ICode = obj.FeminicideRiskForm1ICode.String
	}

	if obj.FeminicideRiskForm1CreationDate.Valid {
		dto.FeminicideRiskForm1CreationDate = obj.FeminicideRiskForm1CreationDate.Time
	}

	if obj.FeminicideRiskForm1UpdateDate.Valid {
		dto.FeminicideRiskForm1UpdateDate = obj.FeminicideRiskForm1UpdateDate.Time
	}

	if obj.FeminicideRiskForm1VictimIdentityName.Valid {
		dto.FeminicideRiskForm1VictimIdentityName = obj.FeminicideRiskForm1VictimIdentityName.String
	}

	if obj.FeminicideRiskForm1BirthDate.Valid {
		dto.FeminicideRiskForm1BirthDate = obj.FeminicideRiskForm1BirthDate.Time
	}

	if obj.FeminicideRiskForm1VictimAddress.Valid {
		dto.FeminicideRiskForm1VictimAddress = obj.FeminicideRiskForm1VictimAddress.String
	}

	if obj.FeminicideRiskForm1VictimLivingZone.Valid {
		dto.FeminicideRiskForm1VictimLivingZone = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimLivingZone.Int64)}
	}

	if obj.FeminicideRiskForm1VictimLivingTownCode.Valid {
		dto.FeminicideRiskForm1VictimLivingTownCode = obj.FeminicideRiskForm1VictimLivingTownCode.String
	}

	if obj.FeminicideRiskForm1VictimSGSSSAffiliation.Valid {
		dto.FeminicideRiskForm1VictimSGSSSAffiliation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimSGSSSAffiliation.Int64)}
	}

	if obj.FeminicideRiskForm1EpsName.Valid {
		dto.FeminicideRiskForm1EpsName = obj.FeminicideRiskForm1EpsName.String
	}

	if obj.FeminicideRiskForm1ContactPhone.Valid {
		dto.FeminicideRiskForm1ContactPhone = obj.FeminicideRiskForm1ContactPhone.String
	}

	if obj.FeminicideRiskForm1ContactEmergencyContactNames.Valid {
		dto.FeminicideRiskForm1ContactEmergencyContactNames = obj.FeminicideRiskForm1ContactEmergencyContactNames.String
	}

	if obj.FeminicideRiskForm1ContactEmergencyContactNumber.Valid {
		dto.FeminicideRiskForm1ContactEmergencyContactNumber = obj.FeminicideRiskForm1ContactEmergencyContactNumber.String
	}

	if obj.FeminicideRiskForm1VictimMaritalStatus.Valid {
		dto.FeminicideRiskForm1VictimMaritalStatus = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimMaritalStatus.Int64)}
	}

	if obj.FeminicideRiskForm1VictimSex.Valid {
		dto.FeminicideRiskForm1VictimSex = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimSex.Int64)}
	}

	if obj.FeminicideRiskForm1VictimGenderIdentity.Valid {
		dto.FeminicideRiskForm1VictimGenderIdentity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimGenderIdentity.Int64)}
	}

	if obj.FeminicideRiskForm1VictimSexualOrientation.Valid {
		dto.FeminicideRiskForm1VictimSexualOrientation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimSexualOrientation.Int64)}
	}

	if obj.FeminicideRiskForm1VictimEthnicAffiliation.Valid {
		dto.FeminicideRiskForm1VictimEthnicAffiliation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimEthnicAffiliation.Int64)}
	}

	if obj.FeminicideRiskForm1VictimIndigenousPeople.Valid {
		dto.FeminicideRiskForm1VictimIndigenousPeople = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimIndigenousPeople.Int64)}
	}

	if obj.FeminicideRiskForm1VictimIsMigrant.Valid {
		dto.FeminicideRiskForm1VictimIsMigrant = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimIsMigrant.Int64)}
	}

	if obj.FeminicideRiskForm1VictimMigrationStatus.Valid {
		dto.FeminicideRiskForm1VictimMigrationStatus = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimMigrationStatus.Int64)}
	}

	if obj.FeminicideRiskForm1VictimMaxEducationLevel.Valid {
		dto.FeminicideRiskForm1VictimMaxEducationLevel = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimMaxEducationLevel.Int64)}
	}

	if obj.FeminicideRiskForm1VictimIsSpecialPopulation.Valid {
		dto.FeminicideRiskForm1VictimIsSpecialPopulation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimIsSpecialPopulation.Int64)}
	}

	if obj.FeminicideRiskForm1VictimCurrentlyHasJob.Valid {
		dto.FeminicideRiskForm1VictimCurrentlyHasJob = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimCurrentlyHasJob.Int64)}
	}

	if obj.FeminicideRiskForm1VictimJobExplanation.Valid {
		dto.FeminicideRiskForm1VictimJobExplanation = obj.FeminicideRiskForm1VictimJobExplanation.String
	}

	if obj.FeminicideRiskForm1VictimAbandonedJobDueToRisk.Valid {
		dto.FeminicideRiskForm1VictimAbandonedJobDueToRisk = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimAbandonedJobDueToRisk.Int64)}
	}

	if obj.FeminicideRiskForm1VictimMainOccupation.Valid {
		dto.FeminicideRiskForm1VictimMainOccupation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimMainOccupation.Int64)}
	}

	if obj.FeminicideRiskForm1VictimHasDisability.Valid {
		dto.FeminicideRiskForm1VictimHasDisability = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimHasDisability.Int64)}
	}

	if obj.FeminicideRiskForm1VictimDisabilityType.Valid {
		dto.FeminicideRiskForm1VictimDisabilityType = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimDisabilityType.Int64)}
	}

	if obj.FeminicideRiskForm1VictimRiskDescription.Valid {
		dto.FeminicideRiskForm1VictimRiskDescription = obj.FeminicideRiskForm1VictimRiskDescription.String
	}

	if obj.FeminicideRiskForm1VictimAdditionalInfo.Valid {
		dto.FeminicideRiskForm1VictimAdditionalInfo = obj.FeminicideRiskForm1VictimAdditionalInfo.String
	}

	if obj.FeminicideRiskForm1VictimIsEconomicProvider.Valid {
		dto.FeminicideRiskForm1VictimIsEconomicProvider = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimIsEconomicProvider.Int64)}
	}

	if obj.FeminicideRiskForm1VictimEconomicProviderExplanation.Valid {
		dto.FeminicideRiskForm1VictimEconomicProviderExplanation = obj.FeminicideRiskForm1VictimEconomicProviderExplanation.String
	}

	if obj.FeminicideRiskForm1VictimHasFamiliarSupport.Valid {
		dto.FeminicideRiskForm1VictimHasFamiliarSupport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimHasFamiliarSupport.Int64)}
	}

	if obj.FeminicideRiskForm1VictimFamiliarSupportType.Valid {
		dto.FeminicideRiskForm1VictimFamiliarSupportType = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimFamiliarSupportType.Int64)}
	}

	if obj.FeminicideRiskForm1VictimPublicTransportAccess.Valid {
		dto.FeminicideRiskForm1VictimPublicTransportAccess = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimPublicTransportAccess.Int64)}
	}

	if obj.FeminicideRiskForm1VictimCommonTransportMode.Valid {
		dto.FeminicideRiskForm1VictimCommonTransportMode = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimCommonTransportMode.Int64)}
	}

	if obj.FeminicideRiskForm1VictimTransportModeExplanation.Valid {
		dto.FeminicideRiskForm1VictimTransportModeExplanation = obj.FeminicideRiskForm1VictimTransportModeExplanation.String
	}

	if obj.FeminicideRiskForm1VictimEstimatedTravelCost.Valid {
		dto.FeminicideRiskForm1VictimEstimatedTravelCost = obj.FeminicideRiskForm1VictimEstimatedTravelCost.Int32
	}

	if obj.FeminicideRiskForm1VictimDifficultiesWithTransport.Valid {
		dto.FeminicideRiskForm1VictimDifficultiesWithTransport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimDifficultiesWithTransport.Int64)}
	}

	if obj.FeminicideRiskForm1VictimDifficultiesExplanation.Valid {
		dto.FeminicideRiskForm1VictimDifficultiesExplanation = obj.FeminicideRiskForm1VictimDifficultiesExplanation.String
	}

	if obj.FeminicideRiskForm1VictimEconomicResourcesForTransport.Valid {
		dto.FeminicideRiskForm1VictimEconomicResourcesForTransport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimEconomicResourcesForTransport.Int64)}
	}

	if obj.FeminicideRiskForm1VictimEconomicResourcesExplanation.Valid {
		dto.FeminicideRiskForm1VictimEconomicResourcesExplanation = obj.FeminicideRiskForm1VictimEconomicResourcesExplanation.String
	}

	if obj.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport.Valid {
		dto.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport.Int64)}
	}

	if obj.FeminicideRiskForm1VictimDebtImpactDetails.Valid {
		dto.FeminicideRiskForm1VictimDebtImpactDetails = obj.FeminicideRiskForm1VictimDebtImpactDetails.String
	}

	if obj.FeminicideRiskForm1VictimReceivedTransportSubsidy.Valid {
		dto.FeminicideRiskForm1VictimReceivedTransportSubsidy = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimReceivedTransportSubsidy.Int64)}
	}

	if obj.FeminicideRiskForm1VictimTransportSubsidyExplanation.Valid {
		dto.FeminicideRiskForm1VictimTransportSubsidyExplanation = obj.FeminicideRiskForm1VictimTransportSubsidyExplanation.String
	}

	if obj.FeminicideRiskForm1VictimSafetyAvoidedTransport.Valid {
		dto.FeminicideRiskForm1VictimSafetyAvoidedTransport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimSafetyAvoidedTransport.Int64)}
	}

	if obj.FeminicideRiskForm1VictimSafetyAvoidedExplanation.Valid {
		dto.FeminicideRiskForm1VictimSafetyAvoidedExplanation = obj.FeminicideRiskForm1VictimSafetyAvoidedExplanation.String
	}

	if obj.FeminicideRiskForm1VictimFoodAccessFrequency.Valid {
		dto.FeminicideRiskForm1VictimFoodAccessFrequency = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimFoodAccessFrequency.Int64)}
	}

	if obj.FeminicideRiskForm1VictimFoodAccessExplanation.Valid {
		dto.FeminicideRiskForm1VictimFoodAccessExplanation = obj.FeminicideRiskForm1VictimFoodAccessExplanation.String
	}

	if obj.FeminicideRiskForm1VictimAgressorFoodRestriction.Valid {
		dto.FeminicideRiskForm1VictimAgressorFoodRestriction = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimAgressorFoodRestriction.Int64)}
	}

	if obj.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation.Valid {
		dto.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation = obj.FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation.String
	}

	if obj.FeminicideRiskForm1VictimFamilyFixedIncome.Valid {
		dto.FeminicideRiskForm1VictimFamilyFixedIncome = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimFamilyFixedIncome.Int64)}
	}

	if obj.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation.Valid {
		dto.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation = obj.FeminicideRiskForm1VictimFamilyFixedIncomeExplanation.String
	}

	if obj.FeminicideRiskForm1VictimIsOnlyProviderForFood.Valid {
		dto.FeminicideRiskForm1VictimIsOnlyProviderForFood = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimIsOnlyProviderForFood.Int64)}
	}

	if obj.FeminicideRiskForm1VictimIsOnlyProviderExplanation.Valid {
		dto.FeminicideRiskForm1VictimIsOnlyProviderExplanation = obj.FeminicideRiskForm1VictimIsOnlyProviderExplanation.String
	}

	if obj.FeminicideRiskForm1VictimJuridicalAssistanceReceived.Valid {
		dto.FeminicideRiskForm1VictimJuridicalAssistanceReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimJuridicalAssistanceReceived.Int64)}
	}

	if obj.FeminicideRiskForm1VictimJuridicalAssistanceExplain.Valid {
		dto.FeminicideRiskForm1VictimJuridicalAssistanceExplain = obj.FeminicideRiskForm1VictimJuridicalAssistanceExplain.String
	}

	if obj.FeminicideRiskForm1VictimWantsJuridicalAssistance.Valid {
		dto.FeminicideRiskForm1VictimWantsJuridicalAssistance = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimWantsJuridicalAssistance.Int64)}
	}

	if obj.FeminicideRiskForm1VictimRepresentationsOfVictims.Valid {
		dto.FeminicideRiskForm1VictimRepresentationsOfVictims = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimRepresentationsOfVictims.Int64)}
	}

	if obj.FeminicideRiskForm1VictimRepresentationsExplanation.Valid {
		dto.FeminicideRiskForm1VictimRepresentationsExplanation = obj.FeminicideRiskForm1VictimRepresentationsExplanation.String
	}

	if obj.FeminicideRiskForm1VictimPsychosocialSupportReceived.Valid {
		dto.FeminicideRiskForm1VictimPsychosocialSupportReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimPsychosocialSupportReceived.Int64)}
	}

	if obj.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain.Valid {
		dto.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain = obj.FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain.String
	}

	if obj.FeminicideRiskForm1VictimUrgentEmotionalCrisis.Valid {
		dto.FeminicideRiskForm1VictimUrgentEmotionalCrisis = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimUrgentEmotionalCrisis.Int64)}
	}

	if obj.FeminicideRiskForm1VictimUrgentCrisisExplanation.Valid {
		dto.FeminicideRiskForm1VictimUrgentCrisisExplanation = obj.FeminicideRiskForm1VictimUrgentCrisisExplanation.String
	}

	if obj.FeminicideRiskForm1AggressorSameResidence.Valid {
		dto.FeminicideRiskForm1AggressorSameResidence = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AggressorSameResidence.Int64)}
	}

	if obj.FeminicideRiskForm1AggressorSameResidenceExplanation.Valid {
		dto.FeminicideRiskForm1AggressorSameResidenceExplanation = obj.FeminicideRiskForm1AggressorSameResidenceExplanation.String
	}

	if obj.FeminicideRiskForm1AggressorKnowsVictimLocation.Valid {
		dto.FeminicideRiskForm1AggressorKnowsVictimLocation = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AggressorKnowsVictimLocation.Int64)}
	}

	if obj.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation.Valid {
		dto.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation = obj.FeminicideRiskForm1AggressorKnowsVictimLocationExplanation.String
	}

	if obj.FeminicideRiskForm1VictimHousingHelpReceived.Valid {
		dto.FeminicideRiskForm1VictimHousingHelpReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimHousingHelpReceived.Int64)}
	}

	if obj.FeminicideRiskForm1VictimHousingHelpExplain.Valid {
		dto.FeminicideRiskForm1VictimHousingHelpExplain = obj.FeminicideRiskForm1VictimHousingHelpExplain.String
	}

	if obj.FeminicideRiskForm1VictimAbandonClothing.Valid {
		dto.FeminicideRiskForm1VictimAbandonClothing = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimAbandonClothing.Int64)}
	}

	if obj.FeminicideRiskForm1VictimAbandonClothingExplain.Valid {
		dto.FeminicideRiskForm1VictimAbandonClothingExplain = obj.FeminicideRiskForm1VictimAbandonClothingExplain.String
	}

	if obj.FeminicideRiskForm1VictimClothingHelpReceived.Valid {
		dto.FeminicideRiskForm1VictimClothingHelpReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1VictimClothingHelpReceived.Int64)}
	}

	if obj.FeminicideRiskForm1VictimClothingHelpExplain.Valid {
		dto.FeminicideRiskForm1VictimClothingHelpExplain = obj.FeminicideRiskForm1VictimClothingHelpExplain.String
	}

	if obj.FeminicideRiskForm1VictimOtherNeeds.Valid {
		dto.FeminicideRiskForm1VictimOtherNeeds = obj.FeminicideRiskForm1VictimOtherNeeds.String
	}

	if obj.FeminicideRiskForm1AnyAssistanceReceived.Valid {
		dto.FeminicideRiskForm1AnyAssistanceReceived = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AnyAssistanceReceived.Int64)}
	}

	if obj.FeminicideRiskForm1AnyAssistanceReceivedCityHall.Valid {
		dto.FeminicideRiskForm1AnyAssistanceReceivedCityHall = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AnyAssistanceReceivedCityHall.Int64)}
	}

	if obj.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice.Valid {
		dto.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice.Int64)}
	}

	if obj.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity.Valid {
		dto.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity.Int64)}
	}

	if obj.FeminicideRiskForm1AnyAssistanceReceivedOther.Valid {
		dto.FeminicideRiskForm1AnyAssistanceReceivedOther = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1AnyAssistanceReceivedOther.Int64)}
	}

	if obj.FeminicideRiskForm1AssistanceReceived.Valid {
		dto.FeminicideRiskForm1AssistanceReceived = obj.FeminicideRiskForm1AssistanceReceived.String
	}

	if obj.FeminicideRiskForm1AssistanceReceivedCityHall.Valid {
		dto.FeminicideRiskForm1AssistanceReceivedCityHall = obj.FeminicideRiskForm1AssistanceReceivedCityHall.String
	}

	if obj.FeminicideRiskForm1AssistanceReceivedWomensOffice.Valid {
		dto.FeminicideRiskForm1AssistanceReceivedWomensOffice = obj.FeminicideRiskForm1AssistanceReceivedWomensOffice.String
	}

	if obj.FeminicideRiskForm1AssistanceReceivedOtherEntity.Valid {
		dto.FeminicideRiskForm1AssistanceReceivedOtherEntity = obj.FeminicideRiskForm1AssistanceReceivedOtherEntity.String
	}

	if obj.FeminicideRiskForm1AssistanceReceivedOther.Valid {
		dto.FeminicideRiskForm1AssistanceReceivedOther = obj.FeminicideRiskForm1AssistanceReceivedOther.String
	}

	if obj.FeminicideRiskForm1FamilyFather.Valid {
		dto.FeminicideRiskForm1FamilyFather = obj.FeminicideRiskForm1FamilyFather.Int64
	}

	if obj.FeminicideRiskForm1FamilyMother.Valid {
		dto.FeminicideRiskForm1FamilyMother = obj.FeminicideRiskForm1FamilyMother.Int64
	}

	if obj.FeminicideRiskForm1FamilyStepfather.Valid {
		dto.FeminicideRiskForm1FamilyStepfather = obj.FeminicideRiskForm1FamilyStepfather.Int64
	}

	if obj.FeminicideRiskForm1FamilyStepmother.Valid {
		dto.FeminicideRiskForm1FamilyStepmother = obj.FeminicideRiskForm1FamilyStepmother.Int64
	}

	if obj.FeminicideRiskForm1FamilyPartner.Valid {
		dto.FeminicideRiskForm1FamilyPartner = obj.FeminicideRiskForm1FamilyPartner.Int64
	}

	if obj.FeminicideRiskForm1FamilySibling1.Valid {
		dto.FeminicideRiskForm1FamilySibling1 = obj.FeminicideRiskForm1FamilySibling1.Int64
	}

	if obj.FeminicideRiskForm1FamilySibling2.Valid {
		dto.FeminicideRiskForm1FamilySibling2 = obj.FeminicideRiskForm1FamilySibling2.Int64
	}

	if obj.FeminicideRiskForm1FamilySibling3.Valid {
		dto.FeminicideRiskForm1FamilySibling3 = obj.FeminicideRiskForm1FamilySibling3.Int64
	}

	if obj.FeminicideRiskForm1FamilySibling4.Valid {
		dto.FeminicideRiskForm1FamilySibling4 = obj.FeminicideRiskForm1FamilySibling4.Int64
	}

	if obj.FeminicideRiskForm1FamilySibling5.Valid {
		dto.FeminicideRiskForm1FamilySibling5 = obj.FeminicideRiskForm1FamilySibling5.Int64
	}

	if obj.FeminicideRiskForm1FamilySonDaughter1.Valid {
		dto.FeminicideRiskForm1FamilySonDaughter1 = obj.FeminicideRiskForm1FamilySonDaughter1.Int64
	}

	if obj.FeminicideRiskForm1FamilySonDaughter2.Valid {
		dto.FeminicideRiskForm1FamilySonDaughter2 = obj.FeminicideRiskForm1FamilySonDaughter2.Int64
	}

	if obj.FeminicideRiskForm1FamilySonDaughter3.Valid {
		dto.FeminicideRiskForm1FamilySonDaughter3 = obj.FeminicideRiskForm1FamilySonDaughter3.Int64
	}

	if obj.FeminicideRiskForm1FamilySonDaughter4.Valid {
		dto.FeminicideRiskForm1FamilySonDaughter4 = obj.FeminicideRiskForm1FamilySonDaughter4.Int64
	}

	if obj.FeminicideRiskForm1FamilySonDaughter5.Valid {
		dto.FeminicideRiskForm1FamilySonDaughter5 = obj.FeminicideRiskForm1FamilySonDaughter5.Int64
	}

	if obj.FeminicideRiskForm1FamilyGrandmother.Valid {
		dto.FeminicideRiskForm1FamilyGrandmother = obj.FeminicideRiskForm1FamilyGrandmother.Int64
	}

	if obj.FeminicideRiskForm1FamilyGrandfather.Valid {
		dto.FeminicideRiskForm1FamilyGrandfather = obj.FeminicideRiskForm1FamilyGrandfather.Int64
	}

	if obj.FeminicideRiskForm1FamilyOtherMember.Valid {
		dto.FeminicideRiskForm1FamilyOtherMember = obj.FeminicideRiskForm1FamilyOtherMember.String
	}

	if obj.FeminicideRiskForm1InterviewDate.Valid {
		dto.FeminicideRiskForm1InterviewDate = obj.FeminicideRiskForm1InterviewDate.Time
	}

	if obj.FeminicideRiskForm1Summary.Valid {
		dto.FeminicideRiskForm1Summary = obj.FeminicideRiskForm1Summary.String
	}

	if obj.FeminicideRiskForm1FeminicideRisk.Valid {
		dto.FeminicideRiskForm1FeminicideRisk = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(obj.FeminicideRiskForm1FeminicideRisk.Int64)}
	}

	return dto
}
