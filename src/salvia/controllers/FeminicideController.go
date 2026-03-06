// Package salvia_ctrl contiene las funciones para el manejo de casos de víctimas,
// incluyendo su creación, actualización, aprobación y asignación de responsables.
package salvia_ctrl

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
	security_config "bitsflow/security/config"
	security_daos "bitsflow/security/dao"
	"math"
	"net/http"
)

type FeminicideRequest struct {
	Feminicide salvia_daos.FeminicideDTO `json:"feminicide"`
}

type FeminicideForm1Request struct {
	Form salvia_daos.FeminicideForm1DTO `json:"form"`
}

// SetFeminicide crea un nuevo caso de víctima a partir de la entrada JSON.
// Realiza la validación de los datos, crea o actualiza usuarios y relaciones,
// y gestiona la transacción en la base de datos. Devuelve un código HTTP y
// un mensaje en formato JSON (éxito o error).
func SetFeminicide(dataInput string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de variables y obtención de conexión a la BD.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	// Se libera la conexión al finalizar la función.
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores durante la validación del DTO.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos para el caso, contacto y dueño.
	var feminicideRequest FeminicideRequest = FeminicideRequest{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	// Variable para almacenar el DTO en formato mapa.
	var dtoMap map[string]interface{} = nil

	// Se parsea el JSON de entrada a un mapa.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FeminicideJSONName, common_config.Locale, collectedErrors)
	utils.JSONToStruct(dataInput, &feminicideRequest)

	// Se asignan valores por defecto al objeto Feminicide.
	salvia_daos.SetFeminicideDefaults(&feminicideRequest.Feminicide, common_dao.SQL_INSERT)
	salvia_daos.SetFeminicideForm1Defaults(&feminicideRequest.Feminicide.FeminicideForm1, common_dao.SQL_INSERT)

	// Verificamos que los datos tengan la estructura esperada.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar (obligatorios y opcionales).
		var checkFields map[string]bool = map[string]bool{
			"FeminicideNames":     true,
			"FeminicideLastNames": true,
			"FeminicideDocType":   true,
			"FeminicideDocNumber": true,
		}
		// Validación de los campos del JSON contra la definición del DTO.
		utils.ValidateJSONInput(&feminicideRequest.Feminicide, dtoMap, salvia_daos.FeminicideJSONName, salvia_daos.FeminicideFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"FeminicideForm1VictimIdentityName":                           true,
			"FeminicideForm1BirthDate":                                    true,
			"FeminicideForm1DeathDate":                                    true,
			"FeminicideForm1VictimAddress":                                true,
			"FeminicideForm1VictimZone":                                   true,
			"FeminicideForm1VictimLivingTownCode":                         true,
			"FeminicideForm1VictimMaritalStatus":                          true,
			"FeminicideForm1VictimSex":                                    true,
			"FeminicideForm1VictimGenderIdentity":                         true,
			"FeminicideForm1VictimSexualOrientation":                      true,
			"FeminicideForm1VictimEthnicity":                              true,
			"FeminicideForm1VictimSpecialPopulation":                      true,
			"FeminicideForm1VictimDisability":                             true,
			"FeminicideForm1VictimDisabilityType":                         true,
			"FeminicideForm1PresumedAggressorNames":                       true,
			"FeminicideForm1PresumedAggressorRelation":                    true,
			"FeminicideForm1PresumedAggressorKnownVGB":                    true,
			"FeminicideForm1InformantNames":                               true,
			"FeminicideForm1InformantIdentityName":                        true,
			"FeminicideForm1InformantDocType":                             true,
			"FeminicideForm1InformantDocNumber":                           true,
			"FeminicideForm1InformantBirthDate":                           true,
			"FeminicideForm1SGSSSAffiliation":                             true,
			"FeminicideForm1EpsName":                                      true,
			"FeminicideForm1InformantAddress":                             true,
			"FeminicideForm1InformantZone":                                true,
			"FeminicideForm1InformantLivingTownCode":                      true,
			"FeminicideForm1InformantPhone":                               true,
			"FeminicideForm1EmergencyContactNames":                        true,
			"FeminicideForm1EmergencyContactNumber":                       true,
			"FeminicideForm1InformantSex":                                 true,
			"FeminicideForm1InformantGenderIdentity":                      true,
			"FeminicideForm1InformantSexualOrientation":                   true,
			"FeminicideForm1InformantEthnicity":                           true,
			"FeminicideForm1InformantMigratoryStatus":                     true,
			"FeminicideForm1InformantMigrationSituation":                  true,
			"FeminicideForm1InformantHighestEducationLevel":               true,
			"FeminicideForm1InformantSpecialPopulation":                   true,
			"FeminicideForm1InformantCurrentEmployment":                   true,
			"FeminicideForm1InformantEmploymentAccess":                    true,
			"FeminicideForm1InformantEmploymentImpactDescription":         true,
			"FeminicideForm1InformantPrimaryOccupation":                   true,
			"FeminicideForm1InformantDisability":                          true,
			"FeminicideForm1InformantDisabilityType":                      true,
			"FeminicideForm1SituationAfterFeminicide":                     true,
			"FeminicideForm1AdditionalInformation":                        true,
			"FeminicideForm1AnyAssistanceReceived":                        true,
			"FeminicideForm1HouseholdExpenseResponsibility":               true,
			"FeminicideForm1PostDeathEconomicAssumption":                  true,
			"FeminicideForm1EconomicAssumptionExplanation":                true,
			"FeminicideForm1AnyDependentPeople":                           true,
			"FeminicideForm1AnyPublicOrPrivateEntity":                     true,
			"FeminicideForm1AnyPublicOrPrivateEntityExplanation":          true,
			"FeminicideForm1PublicTransportAccess":                        true,
			"FeminicideForm1PreferredTransportationMode":                  true,
			"FeminicideForm1PreferredTransportationModeExplanation":       true,
			"FeminicideForm1TransportationCostEstimate":                   true,
			"FeminicideForm1TransportDifficulty":                          true,
			"FeminicideForm1TransportDifficultyExplanation":               true,
			"FeminicideForm1EconomicResourcesForTransport":                true,
			"FeminicideForm1DebtOrHelpDueToTransport":                     true,
			"FeminicideForm1DebtImpactExplanation":                        true,
			"FeminicideForm1TransportSubsidyReceived":                     true,
			"FeminicideForm1TransportSubsidyExplanation":                  true,
			"FeminicideForm1SafetyTransportationConcern":                  true,
			"FeminicideForm1SafetyTransportationExplanation":              true,
			"FeminicideForm1FoodAccessFrequency":                          true,
			"FeminicideForm1FoodAccessExplanation":                        true,
			"FeminicideForm1AggressorFoodRestriction":                     true,
			"FeminicideForm1AggressorFoodRestrictionExplanation":          true,
			"FeminicideForm1FoodIncomeSupport":                            true,
			"FeminicideForm1FoodIncomeSupportExplanation":                 true,
			"FeminicideForm1JuridicalAssistanceReceived":                  true,
			"FeminicideForm1JuridicalAssistanceExplanation":               true,
			"FeminicideForm1VictimRepresentation":                         true,
			"FeminicideForm1VictimRepresentationExplanation":              true,
			"FeminicideForm1PsychosocialSupportReceived":                  true,
			"FeminicideForm1PsychosocialSupportExplanation":               true,
			"FeminicideForm1EmergencyEmotionalCrisis":                     true,
			"FeminicideForm1EmergencyEmotionalCrisisExplanation":          true,
			"FeminicideForm1AggressorSameResidence":                       true,
			"FeminicideForm1AggressorSameResidenceExplanation":            true,
			"FeminicideForm1AggressorLocationKnown":                       true,
			"FeminicideForm1AggressorLocationKnownExplanation":            true,
			"FeminicideForm1AnyTypeOfAssistanceReceived":                  true,
			"FeminicideForm1AnyTypeOfAssistanceReceivedExplanation":       true,
			"FeminicideForm1CompensationFundInsuranceCoverage":            true,
			"FeminicideForm1CompensationFundInsuranceCoverageExplanation": true,
			"FeminicideForm1FuneralSubsidyReceived":                       true,
			"FeminicideForm1FuneralSubsidyExplanation":                    true,
			"FeminicideForm1FuneralFundsAvailable":                        true,
			"FeminicideForm1FuneralFundsAvailableExplanation":             true,
			"FeminicideForm1FuneralCostValue":                             true,
			"FeminicideForm1RenameReputationImpact":                       true,
			"FeminicideForm1RenameReputationExplanation":                  true,
			"FeminicideForm1AdditionalNeedsDescription":                   true,
			"FeminicideForm1ActionPlan":                                   true,
			"FeminicideForm1Summary":                                      true,

			"FeminicideForm1AssistanceReceived":                false,
			"FeminicideForm1AnyAssistanceReceivedCityHall":     false,
			"FeminicideForm1AssistanceReceivedCityHall":        false,
			"FeminicideForm1AnyAssistanceReceivedWomensOffice": false,
			"FeminicideForm1AssistanceReceivedWomensOffice":    false,
			"FeminicideForm1AnyAssistanceReceivedOtherEntity":  false,
			"FeminicideForm1AssistanceReceivedOtherEntity":     false,
			"FeminicideForm1AnyAssistanceReceivedOther":        false,
			"FeminicideForm1AssistanceReceivedOther":           false,
		}

		utils.ValidateJSONInput(&feminicideRequest.Feminicide.FeminicideForm1, dtoMap, salvia_daos.FeminicideJSONName+"."+salvia_daos.FeminicideForm1JSONName, salvia_daos.FeminicideForm1FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.FeminicideJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

	default:
		// Si el JSON no tiene la estructura esperada, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.FeminicideJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	//Ahora se verifican los enums únicos
	getAndVerifyFeminicideEnums(&feminicideRequest.Feminicide.FeminicideForm1, collectedErrors, true)

	//Ahora se verifican los enums múltiples que se asociarán al form2
	getAndVerifyFeminicideEnumsMultiple(&feminicideRequest.Feminicide.FeminicideForm1, collectedErrors)

	//Ahora algunos campos opcionales
	/*if feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1IncomeGenerationMethod.FeminicideForm1EnumsCode == "pr" {
		if feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsStartTime.IsZero() {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRequest.Feminicide.FeminicideForm1, "FeminicideForm1FactsStartTime", "json"), "common_validation_field_date_error", "common_global_error", common_config.Locale)
		}
	}*/

	// Si se detectaron errores en la validación, se retorna BadRequest con los errores.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se determina el dueño del caso
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{s.UserICode},
	}

	err = salvia_daos.GetCaseOwner(by, &owner, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se crea el caso en la BD.
	if err = salvia_daos.SetFeminicide(&feminicideRequest.Feminicide, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Feminicide = feminicideRequest.Feminicide

	if err = salvia_daos.SetFeminicideForm1(&feminicideRequest.Feminicide.FeminicideForm1, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Ahora se insertan todos los campos múltiples
	//Ahora se verifican los enums múltiples que se asociarán al form2
	if err = setFeminicideEnumsMultiple(&feminicideRequest.Feminicide.FeminicideForm1, connData, dbClientConfig, dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se asigna la autoría del caso al usuario en sesión.
	var relCaseOwnerFeminicide salvia_daos.RelCaseOwnerFeminicideDTO = salvia_daos.RelCaseOwnerFeminicideDTO{}
	salvia_daos.SetRelCaseOwnerFeminicideDefaults(&relCaseOwnerFeminicide, common_dao.SQL_INSERT)
	relCaseOwnerFeminicide.RelCaseOwnerFeminicide_CaseOwner = owner.CaseOwnerId
	relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Feminicide = feminicideRequest.Feminicide.FeminicideId

	if err = salvia_daos.SetRelCaseOwnerFeminicide(&relCaseOwnerFeminicide, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
	if err = salvia_daos.UpdateFeminicideOwnersAndRolesByFeminicideId(feminicideRequest.Feminicide.FeminicideId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se realiza el commit de la transacción y se retorna el resultado.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}

// UpdateFeminicide actualiza los datos de un caso de víctima existente.
// Recibe la entrada JSON, el identificador del caso (feminicideICode) y realiza la actualización
// validando el DTO, actualizando datos y gestionando la transacción.
/*
func UpdateFeminicide(dataInput string, feminicideICode string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicializa variables y establece la conexión.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores de validación.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos: caso, dueño y el caso a actualizar.
	var feminicideRequest FeminicideRequest = FeminicideRequest{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	var feminicideToUpdate salvia_daos.FeminicideDTO = salvia_daos.FeminicideDTO{}
	var feminicideFormToUpdate salvia_daos.FeminicideForm1DTO = salvia_daos.FeminicideForm1DTO{}
	var code int

	// Slices para gestionar los momentos asociados al caso.
	var momentsToDelete []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}
	var moments []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}

	// Variable para almacenar el DTO en forma de mapa.
	var dtoMap map[string]interface{} = nil

	// Se obtiene el mapa DTO a partir de la entrada JSON.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FeminicideJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &feminicideRequest)

	// Se establecen valores por defecto para el caso (modo actualización).
	salvia_daos.SetFeminicideDefaults(&feminicideRequest.Feminicide, common_dao.SQL_UPDATE, s)
	salvia_daos.SetFeminicideForm1Defaults(&feminicideRequest.Feminicide.FeminicideForm1, common_dao.SQL_INSERT)

	// Se verifica la estructura del DTO.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar.
		var checkFields map[string]bool = map[string]bool{
			"FeminicideNames":     true,
			"FeminicideLastNames": true,
			"FeminicideDocType":   true,
			"FeminicideDocNumber": true,
			"FeminicideTownCode":  true,
		}
		// Validación de los datos del caso.
		utils.ValidateJSONInput(&feminicideRequest.Feminicide, dtoMap, salvia_daos.FeminicideJSONName, salvia_daos.FeminicideFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"FeminicideForm1Nick":                              true,
			"FeminicideForm1BirthDate":                         true,
			"FeminicideForm1ViolenceTownCode":                  true,
			"FeminicideForm1Address":                           true,
			"FeminicideForm1LivingLatitude":                    false,
			"FeminicideForm1LivingLongitude":                   false,
			"FeminicideForm1Phone":                             true,
			"FeminicideForm1Email":                             true,
			"FeminicideForm1GenderIdentity":                    true,
			"FeminicideForm1SexualOrientation":                 true,
			"FeminicideForm1Origin":                            true,
			"FeminicideForm1Occupation":                        true,
			"FeminicideForm1OccupationOther":                   false,
			"FeminicideForm1VictimEthnicGroup":                 true,
			"FeminicideForm1VictimEthnicGroupOther":            false,
			"FeminicideForm1VictimContactNames":                true,
			"FeminicideForm1VictimContactPhone":                true,
			"FeminicideForm1VictimContactKinship":              true,
			"FeminicideForm1VictimNumChildren":                 false,
			"FeminicideForm1VictimMaritalStatus":               true,
			"FeminicideForm1VictimMaritalStatusOther":          false,
			"FeminicideForm1VictimChildrenAge":                 false,
			"FeminicideForm1VictimDisability":                  true,
			"FeminicideForm1VictimSpecialSupport":              true,
			"FeminicideForm1FactsOccurrence":                   true,
			"FeminicideForm1FactsStartTime":                    true,
			"FeminicideForm1FactsEndTime":                      true,
			"FeminicideForm1FactsWeekday":                      true,
			"FeminicideForm1FactsDate":                         true,
			"FeminicideForm1FactsDescription":                  true,
			"FeminicideForm1VictimViolenceExperienced":         true,
			"FeminicideForm1VictimViolenceExperiencedOther":    false,
			"FeminicideForm1VictimViolenceScope":               true,
			"FeminicideForm1VictimFemicideRisk":                true,
			"FeminicideForm1VictimAggressor":                   true,
			"FeminicideForm1VictimRelationshipWithAggressor":   true,
			"FeminicideForm1VictimAggressorName":               false,
			"FeminicideForm1VictimAggressorDocType":            false,
			"FeminicideForm1VictimAggressorDocNumber":          false,
			"FeminicideForm1VictimAggressorAddress":            false,
			"FeminicideForm1VictimAggressorPhone":              false,
			"FeminicideForm1Age":                               true,
			"FeminicideForm1VictimNationality":                 true,
			"FeminicideForm1VictimNationalityOther":            false,
			"FeminicideForm1VictimForeignerImmigrationStatus":  false,
			"FeminicideForm1VictimGender":                      true,
			"FeminicideForm1VictimGenderIdentityOther":         false,
			"FeminicideForm1VictimSexualOrientationOther":      false,
			"FeminicideForm1VictimDependents":                  true,
			"FeminicideForm1VictimPreviouslyReportedSituation": true,
			"FeminicideForm1VictimIfPreviouslyReported":        false,
			"FeminicideForm1VictimIfAfro":                      false,
			"FeminicideForm1VictimIfIndigenous":                false,
			"FeminicideForm1VictimIfIndigenousTongue":          false,
			"FeminicideForm1VictimIfPeasant":                   true,
			"FeminicideForm1VictimIfArmedConflict":             true,
			"FeminicideForm1VictimViolenceScene":               true,
			"FeminicideForm1PhysicalViolenceIncreased":         true,
			"FeminicideForm1SeparatedFromPartnerLastYear":      true,
			"FeminicideForm1ThreatenedWithWeapon":              true,
			"FeminicideForm1ThreatenedToKillOrHarmChildren":    true,
			"FeminicideForm1JealousAndViolent":                 true,
			"FeminicideForm1BelievesCapableOfKilling":          true,
		}

		utils.ValidateJSONInput(&feminicideRequest.Feminicide.FeminicideForm1, dtoMap, salvia_daos.FeminicideJSONName+"."+salvia_daos.FeminicideForm1JSONName, salvia_daos.FeminicideForm1FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.FeminicideJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}


	default:
		// Si el JSON no cumple la estructura, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.FeminicideJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si existen errores, se retorna un BadRequest.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se obtiene el caso a actualizar a partir del identificador.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideICode"},
		AttrsValue: []interface{}{feminicideICode},
	}

	if err = salvia_daos.GetFeminicide(by, &feminicideToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideForm1Feminicide"},
		AttrsValue: []interface{}{feminicideToUpdate.FeminicideId},
	}

	if err = salvia_daos.GetFeminicideForm1(by, &feminicideFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Verifica si se produjo un cambio en el documento de identidad.
	if feminicideRequest.Feminicide.FeminicideDocType != feminicideToUpdate.FeminicideDocType || feminicideRequest.Feminicide.FeminicideDocNumber != feminicideToUpdate.FeminicideDocNumber {
		var oldProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}
		var newProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
			AttrsValue: []interface{}{feminicideRequest.Feminicide.FeminicideDocType, feminicideRequest.Feminicide.FeminicideDocNumber},
		}

		err = security_daos.GetGeneralUserProfile(by, &newProfile, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// No existe un usuario con el nuevo documento; se actualiza el perfil del usuario.
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
				AttrsValue: []interface{}{feminicideToUpdate.FeminicideDocType, feminicideToUpdate.FeminicideDocNumber},
			}
			err = security_daos.GetGeneralUserProfile(by, &oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_feminicide_find_profile_fail"]
			}
			// Se actualiza el perfil con el nuevo documento.
			security_daos.SetGeneralUserProfileDefaults(&oldProfile, common_dao.SQL_UPDATE)
			oldProfile.GeneralUserProfileDocType = feminicideRequest.Feminicide.FeminicideDocType
			oldProfile.GeneralUserProfileDocNumber = feminicideRequest.Feminicide.FeminicideDocNumber

			err = security_daos.UpdateGeneralUserProfile(&oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_feminicide_user_not_found"]
			}
		} else {
			// Si ya existe un perfil con el nuevo documento, se retorna un error.
			utils.SetError(collectedErrors, salvia_daos.FeminicideJSONName, utils.GetTag(&feminicideRequest.Feminicide, "FeminicideDocNumber", "json"),
				"update_feminicide_find_profile_already_exist", "salvia_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	// Se inicia la transacción para la actualización.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualizan los campos del caso con los nuevos valores.
	feminicideToUpdate.FeminicideNames = feminicideRequest.Feminicide.FeminicideNames
	feminicideToUpdate.FeminicideLastNames = feminicideRequest.Feminicide.FeminicideLastNames
	feminicideToUpdate.FeminicideDocType = feminicideRequest.Feminicide.FeminicideDocType
	feminicideToUpdate.FeminicideDocNumber = feminicideRequest.Feminicide.FeminicideDocNumber
	feminicideToUpdate.FeminicideTownCode = feminicideRequest.Feminicide.FeminicideTownCode

	feminicideFormToUpdate.FeminicideForm1Nick = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Nick
	feminicideFormToUpdate.FeminicideForm1BirthDate = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1BirthDate
	feminicideFormToUpdate.FeminicideForm1Address = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Address
	feminicideFormToUpdate.FeminicideForm1LivingLatitude = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1LivingLatitude
	feminicideFormToUpdate.FeminicideForm1LivingLongitude = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1LivingLongitude
	feminicideFormToUpdate.FeminicideForm1Phone = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Phone
	feminicideFormToUpdate.FeminicideForm1Email = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Email
	feminicideFormToUpdate.FeminicideForm1GenderIdentity = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1GenderIdentity
	feminicideFormToUpdate.FeminicideForm1SexualOrientation = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1SexualOrientation
	feminicideFormToUpdate.FeminicideForm1Origin = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Origin
	feminicideFormToUpdate.FeminicideForm1Occupation = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Occupation
	feminicideFormToUpdate.FeminicideForm1OccupationOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1OccupationOther
	feminicideFormToUpdate.FeminicideForm1VictimEthnicGroup = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimEthnicGroup
	feminicideFormToUpdate.FeminicideForm1VictimEthnicGroupOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimEthnicGroupOther
	feminicideFormToUpdate.FeminicideForm1VictimContactNames = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimContactNames
	feminicideFormToUpdate.FeminicideForm1VictimContactPhone = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimContactPhone
	feminicideFormToUpdate.FeminicideForm1VictimContactKinship = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimContactKinship
	feminicideFormToUpdate.FeminicideForm1VictimNumChildren = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimNumChildren
	feminicideFormToUpdate.FeminicideForm1VictimMaritalStatus = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimMaritalStatus
	feminicideFormToUpdate.FeminicideForm1VictimMaritalStatusOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimMaritalStatusOther
	feminicideFormToUpdate.FeminicideForm1VictimChildrenAge = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimChildrenAge
	feminicideFormToUpdate.FeminicideForm1VictimDisability = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimDisability
	feminicideFormToUpdate.FeminicideForm1VictimSpecialSupport = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimSpecialSupport
	feminicideFormToUpdate.FeminicideForm1FactsOccurrence = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsOccurrence
	feminicideFormToUpdate.FeminicideForm1FactsStartTime = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsStartTime
	feminicideFormToUpdate.FeminicideForm1FactsEndTime = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsEndTime
	feminicideFormToUpdate.FeminicideForm1FactsWeekday = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsWeekday
	feminicideFormToUpdate.FeminicideForm1FactsDate = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsDate
	feminicideFormToUpdate.FeminicideForm1FactsDescription = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1FactsDescription
	feminicideFormToUpdate.FeminicideForm1VictimViolenceExperienced = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimViolenceExperienced
	feminicideFormToUpdate.FeminicideForm1VictimViolenceExperiencedOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimViolenceExperiencedOther
	feminicideFormToUpdate.FeminicideForm1VictimViolenceScope = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimViolenceScope
	feminicideFormToUpdate.FeminicideForm1VictimFemicideRisk = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimFemicideRisk
	feminicideFormToUpdate.FeminicideForm1VictimAggressor = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimAggressor
	feminicideFormToUpdate.FeminicideForm1VictimRelationshipWithAggressor = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimRelationshipWithAggressor
	feminicideFormToUpdate.FeminicideForm1VictimAggressorName = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimAggressorName
	feminicideFormToUpdate.FeminicideForm1VictimAggressorDocType = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimAggressorDocType
	feminicideFormToUpdate.FeminicideForm1VictimAggressorDocNumber = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimAggressorDocNumber
	feminicideFormToUpdate.FeminicideForm1VictimAggressorAddress = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimAggressorAddress
	feminicideFormToUpdate.FeminicideForm1VictimAggressorPhone = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimAggressorPhone

	feminicideFormToUpdate.FeminicideForm1Age = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1Age
	feminicideFormToUpdate.FeminicideForm1VictimNationality = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimNationality
	feminicideFormToUpdate.FeminicideForm1VictimNationalityOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimNationalityOther
	feminicideFormToUpdate.FeminicideForm1VictimForeignerImmigrationStatus = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimForeignerImmigrationStatus
	feminicideFormToUpdate.FeminicideForm1VictimGender = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimGender
	feminicideFormToUpdate.FeminicideForm1VictimGenderIdentityOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimGenderIdentityOther
	feminicideFormToUpdate.FeminicideForm1VictimSexualOrientationOther = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimSexualOrientationOther
	feminicideFormToUpdate.FeminicideForm1VictimDependents = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimDependents
	feminicideFormToUpdate.FeminicideForm1VictimPreviouslyReportedSituation = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimPreviouslyReportedSituation
	feminicideFormToUpdate.FeminicideForm1VictimIfPreviouslyReported = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimIfPreviouslyReported

	feminicideFormToUpdate.FeminicideForm1VictimIfAfro = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimIfAfro
	feminicideFormToUpdate.FeminicideForm1VictimIfIndigenous = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimIfIndigenous
	feminicideFormToUpdate.FeminicideForm1VictimIfIndigenousTongue = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimIfIndigenousTongue
	feminicideFormToUpdate.FeminicideForm1VictimIfPeasant = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimIfPeasant
	feminicideFormToUpdate.FeminicideForm1VictimIfArmedConflict = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimIfArmedConflict
	feminicideFormToUpdate.FeminicideForm1VictimViolenceScene = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1VictimViolenceScene

	feminicideFormToUpdate.FeminicideForm1PhysicalViolenceIncreased = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1PhysicalViolenceIncreased
	feminicideFormToUpdate.FeminicideForm1SeparatedFromPartnerLastYear = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1SeparatedFromPartnerLastYear
	feminicideFormToUpdate.FeminicideForm1ThreatenedWithWeapon = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1ThreatenedWithWeapon
	feminicideFormToUpdate.FeminicideForm1ThreatenedToKillOrHarmChildren = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1ThreatenedToKillOrHarmChildren
	feminicideFormToUpdate.FeminicideForm1JealousAndViolent = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1JealousAndViolent
	feminicideFormToUpdate.FeminicideForm1BelievesCapableOfKilling = feminicideRequest.Feminicide.FeminicideForm1.FeminicideForm1BelievesCapableOfKilling

	// Se obtiene el dueño del caso basado en el usuario en sesión.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{s.UserICode},
	}

	err = salvia_daos.GetCaseOwner(by, &owner, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se actualiza el objeto en la BD.
	if err = salvia_daos.UpdateFeminicide(&feminicideToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if err = salvia_daos.UpdateFeminicideForm1(&feminicideFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen los momentos (acciones o eventos) asociados al caso.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentFeminicide"},
		AttrsValue: []interface{}{feminicideToUpdate.FeminicideId},
	}

	momentsToDelete, err = salvia_daos.GetMoments(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se procesan las ramas para determinar cuáles momentos eliminar y cuáles crear.
	for mcode, sectors := range feminicideRequest.Feminicide.FeminicideEntityBranches {
		for _, entities := range sectors {
			for _, branchICode := range entities {
				if branchICode != "" {
					var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
					code, _, branch = GetEntityBranchByICode(branchICode, connData, dbClientConfig, dbServerConfig)
					if code != http.StatusOK {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return code, ""
					}

					var moment salvia_daos.MomentDTO = salvia_daos.MomentDTO{}
					salvia_daos.SetMomentDefaults(&moment, common_dao.SQL_INSERT)
					moment.MomentCode = mcode
					moment.MomentEntityBranch = branch
					moment.MomentApprovalSource = salvia_config.MOMENT_APPROVAL_SOURCE["MANUAL"]
					moment.MomentFeminicide = feminicideToUpdate

					//
					moments = append(moments, moment)
				}
			}
		}
	}

	// Se eliminan de momentsToDelete aquellos momentos que ya existen en moments.
	momentsToDelete, moments = cleanMoments(momentsToDelete, moments)

	// Se insertan los nuevos momentos.
	for _, moment := range moments {
		if err = salvia_daos.SetMoment(&moment, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	// Se eliminan los momentos que quedaron pendientes de borrar.
	for _, moment := range momentsToDelete {
		if err = salvia_daos.RemoveMomentById(&moment, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	// Se asegura la relación de autoría entre el dueño y el caso.
	var relCaseOwners []salvia_daos.RelCaseOwnerFeminicideDTO = []salvia_daos.RelCaseOwnerFeminicideDTO{}
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelCaseOwnerFeminicide_CaseOwner", "RelCaseOwnerFeminicide_Feminicide"},
		AttrsValue: []interface{}{owner.CaseOwnerId, feminicideToUpdate.FeminicideId},
	}
	relCaseOwners, err = salvia_daos.GetRelCaseOwnerFeminicides(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Si no existe la relación, se crea y se actualiza el contador del dueño.
	if len(relCaseOwners) == 0 {
		var relCaseOwnerFeminicide salvia_daos.RelCaseOwnerFeminicideDTO = salvia_daos.RelCaseOwnerFeminicideDTO{}
		salvia_daos.SetRelCaseOwnerFeminicideDefaults(&relCaseOwnerFeminicide, common_dao.SQL_INSERT)
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_CaseOwner = owner.CaseOwnerId
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Feminicide = feminicideToUpdate.FeminicideId

		if err = salvia_daos.SetRelCaseOwnerFeminicide(&relCaseOwnerFeminicide, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		owner.CaseOwnerNumCases = owner.CaseOwnerNumCases + 1
		if err = salvia_daos.UpdateCaseOwner(&owner, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
	if err = salvia_daos.UpdateFeminicideOwnersAndRolesByFeminicideId(feminicideToUpdate.FeminicideId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se finaliza la transacción.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideRequest.Feminicide.FeminicideNewUser)
}*/

// GetFeminicideByICode retorna un caso de víctima a partir de su identificador único (ICode).
// Además, recupera información geográfica relacionada (municipio, ciudad, departamento)
// y los momentos (logs) asociados al caso.
func GetFeminicideByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.FeminicideDTO) {
	var err error = nil

	// Se libera la conexión si ésta no es válida.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el caso.
	var feminicide salvia_daos.FeminicideDTO = salvia_daos.FeminicideDTO{}
	var feminicideForm1 salvia_daos.FeminicideForm1DTO = salvia_daos.FeminicideForm1DTO{}

	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, feminicide
		}

	} else {
		var town security_daos.TownDTO
		var town2 security_daos.TownDTO

		var city security_daos.CityDTO
		var city2 security_daos.CityDTO

		var department security_daos.DepartmentDTO
		var department2 security_daos.DepartmentDTO

		// Se consulta el caso por su identificador.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"FeminicideICode"},
			AttrsValue: []interface{}{id},
		}

		err = salvia_daos.GetFeminicide(by, &feminicide, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), feminicide
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"FeminicideForm1Feminicide"},
			AttrsValue: []interface{}{feminicide.FeminicideId},
		}

		if err = salvia_daos.GetFeminicideForm1(by, &feminicideForm1, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		feminicide.FeminicideForm1 = feminicideForm1
		// Se recupera la información geográfica de la víctima y familiar
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownCode"},
			AttrsValue: []interface{}{feminicideForm1.FeminicideForm1VictimLivingTownCode},
		}
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town.TownCity},
		}
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city.CityDepartment},
		}
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownCode"},
			AttrsValue: []interface{}{feminicideForm1.FeminicideForm1VictimLivingTownCode},
		}
		err = security_daos.GetTown(by, &town2, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town2.TownCity},
		}
		err = security_daos.GetCity(by, &city2, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city2.CityDepartment},
		}
		err = security_daos.GetDepartment(by, &department2, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		//Ahora se cargan los enums sencillos
		loadFeminicideEnums(&feminicide)

		//Ahora se cargan los enums múltiples
		var enums []salvia_daos.VictimCaseForm2EnumsDTO

		enums, err = salvia_daos.GetVictimCasesForm2EnumsByFeminicideForm1Id(feminicide.FeminicideForm1.FeminicideForm1Id, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideDTO{}
		}

		loadFeminicideEnumsMultiple(&feminicide.FeminicideForm1, enums)

		// Se asignan los datos geográficos, los momentos y el seguimiento al caso.
		feminicide.FeminicideForm1.FeminicideForm1VictimLivingDepartment = department
		feminicide.FeminicideForm1.FeminicideForm1VictimLivingCity = city
		feminicide.FeminicideForm1.FeminicideForm1VictimLivingTown = town

		feminicide.FeminicideForm1.FeminicideForm1InformantLivingDepartment = department2
		feminicide.FeminicideForm1.FeminicideForm1InformantLivingCity = city2
		feminicide.FeminicideForm1.FeminicideForm1InformantLivingTown = town2

		return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicide), feminicide

	}
	return http.StatusInternalServerError, "", salvia_daos.FeminicideDTO{}
}

// UpdateFeminicideStatusByICode actualiza el estado de un caso de víctima basado en su ICode.
/*func UpdateFeminicideStatusByICode(newStatus string, id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicideToUpdate salvia_daos.FeminicideDTO = salvia_daos.FeminicideDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideICode"},
		AttrsValue: []interface{}{id},
	}

	if err = salvia_daos.GetFeminicide(by, &feminicideToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se asigna el nuevo estado.
	feminicideToUpdate.FeminicideStatus = newStatus

	if err = salvia_daos.UpdateFeminicide(&feminicideToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideToUpdate)
}*/

// GetFeminicidesByTownCode obtiene todos los casos de víctimas asociados a un código de municipio.
func GetFeminicidesByTownCode(townCode string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var count = 0

	if townCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, count
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicides []salvia_daos.FeminicideDTO = []salvia_daos.FeminicideDTO{}

	// Se consulta la BD utilizando el código del municipio.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideTownCode"},
		AttrsValue: []interface{}{townCode},
	}

	if feminicides, count, err = salvia_daos.GetFeminicides(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicides), count
}

// GetDepartmentFeminicidesByTownCode obtiene todos los casos de víctimas asociados al departamento relacionado con el código de municipio.
/*func GetDepartmentFeminicidesByTownCode(townCode string, feminicideStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var count = 0

	if townCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, count
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	resCode, errStr, department := security_ctrl.GetDepartmentByTownCode(townCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

	if resCode != http.StatusOK {
		return http.StatusInternalServerError, errStr, count
	}

	var feminicides []salvia_daos.FeminicideDTO = []salvia_daos.FeminicideDTO{}

	if feminicides, count, err = salvia_daos.GetFeminicidesByDepartmentICode(department.DepartmentICode, feminicideStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicides), count
}*/

// GetFeminicidesByDocument obtiene casos de víctimas filtrados por tipo y número de documento.
func GetFeminicidesByDocument(docType string, docNumber string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicides []salvia_daos.FeminicideDTO = []salvia_daos.FeminicideDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideDocType", "FeminicideDocNumber"},
		AttrsValue: []interface{}{docType, docNumber},
	}

	if feminicides, count, err = salvia_daos.GetFeminicides(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{feminicides, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}

/*
// GetFeminicidesByOwnerUserICode obtiene casos de víctimas asociados al ICode del usuario dueño.
func GetFeminicidesByOwnerUserICode(userIcode string, feminicideStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicides []salvia_daos.FeminicideDTO = []salvia_daos.FeminicideDTO{}

	if feminicides, count, err = salvia_daos.GetFeminicidesByOwnerUserICode(userIcode, feminicideStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{feminicides, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}*/

// GetFeminicideByAll retorna todos los casos de víctimas almacenados en la base de datos.

func GetFeminicideByAll(page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	cases, count, err := salvia_daos.GetAllFeminicides(page, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{cases, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"})
	}
	return resCode, resData, count
}

func getAndVerifyFeminicideEnums(feminicideForm1 *salvia_daos.FeminicideForm1DTO, collectedErrors map[string]map[string]string, verify bool) bool {
	var opRes bool = true

	err := salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimSex)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimSex", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimMaritalStatus)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimMaritalStatus", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimGenderIdentity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimGenderIdentity", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimSexualOrientation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimSexualOrientation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimEthnicity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimEthnicity", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if feminicideForm1.FeminicideForm1VictimEthnicity.VictimCaseForm2EnumsCode == "in" {
		// VictimCaseForm2IndigenousPeople
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantIndigenousPeople)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantIndigenousPeople", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimSpecialPopulation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimSpecialPopulation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimDisability)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimDisability", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimDisabilityType)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimDisabilityType", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1PresumedAggressorRelation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1PresumedAggressorRelation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1PresumedAggressorKnownVGB)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1PresumedAggressorKnownVGB", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1SGSSSAffiliation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1SGSSSAffiliation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantZone)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantZone", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantSex)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantSex", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantGenderIdentity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantGenderIdentity", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantSexualOrientation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantSexualOrientation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantEthnicity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantEthnicity", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if feminicideForm1.FeminicideForm1InformantEthnicity.VictimCaseForm2EnumsCode == "in" {
		// VictimCaseForm2IndigenousPeople
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantIndigenousPeople)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantIndigenousPeople", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantMigratoryStatus)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantMigratoryStatus", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantMigrationSituation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantMigrationSituation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantHighestEducationLevel)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantHighestEducationLevel", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantSpecialPopulation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantSpecialPopulation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantCurrentEmployment)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantCurrentEmployment", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantEmploymentAccess)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantEmploymentAccess", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantPrimaryOccupation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantPrimaryOccupation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantDisability)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantDisability", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1InformantDisabilityType)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1InformantDisabilityType", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyAssistanceReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyAssistanceReceived", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if feminicideForm1.FeminicideForm1AnyAssistanceReceived.VictimCaseForm2EnumsCode == "y" {

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyAssistanceReceivedCityHall)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyAssistanceReceivedCityHall", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyAssistanceReceivedWomensOffice)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyAssistanceReceivedWomensOffice", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyAssistanceReceivedOtherEntity)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyAssistanceReceivedOtherEntity", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyAssistanceReceivedOther)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyAssistanceReceivedOther", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1HouseholdExpenseResponsibility)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1HouseholdExpenseResponsibility", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1PostDeathEconomicAssumption)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1PostDeathEconomicAssumption", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyDependentPeople)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyDependentPeople", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyPublicOrPrivateEntity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyPublicOrPrivateEntity", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1PublicTransportAccess)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1PublicTransportAccess", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1PreferredTransportationMode)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1PreferredTransportationMode", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1TransportDifficulty)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1TransportDifficulty", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1EconomicResourcesForTransport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1EconomicResourcesForTransport", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1DebtOrHelpDueToTransport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1DebtOrHelpDueToTransport", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1TransportSubsidyReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1TransportSubsidyReceived", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1SafetyTransportationConcern)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1SafetyTransportationConcern", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1FoodAccessFrequency)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1FoodAccessFrequency", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AggressorFoodRestriction)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AggressorFoodRestriction", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1FoodIncomeSupport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1FoodIncomeSupport", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1JuridicalAssistanceReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1JuridicalAssistanceReceived", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1VictimRepresentation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1VictimRepresentation", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1PsychosocialSupportReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1PsychosocialSupportReceived", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1EmergencyEmotionalCrisis)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1EmergencyEmotionalCrisis", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AggressorSameResidence)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AggressorSameResidence", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AggressorLocationKnown)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AggressorLocationKnown", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1AnyTypeOfAssistanceReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1AnyTypeOfAssistanceReceived", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1CompensationFundInsuranceCoverage)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1CompensationFundInsuranceCoverage", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1FuneralSubsidyReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1FuneralSubsidyReceived", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1FuneralFundsAvailable)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1FuneralFundsAvailable", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1FuneralCostValue)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1FuneralCostValue", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1RenameReputationImpact)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1RenameReputationImpact", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1ActionPlan)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1ActionPlan", "json"), "feminicide_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

func loadFeminicideEnums(feminicide *salvia_daos.FeminicideDTO) {

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimSex)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimMaritalStatus)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimGenderIdentity)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimSexualOrientation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimEthnicity)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimSpecialPopulation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimDisability)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimDisabilityType)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1PresumedAggressorRelation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1PresumedAggressorKnownVGB)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1SGSSSAffiliation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantZone)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantSex)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantGenderIdentity)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantSexualOrientation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantEthnicity)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantMigratoryStatus)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantMigrationSituation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantHighestEducationLevel)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantSpecialPopulation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantCurrentEmployment)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantEmploymentAccess)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantPrimaryOccupation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantDisability)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1InformantDisabilityType)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AnyAssistanceReceived)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1HouseholdExpenseResponsibility)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1PostDeathEconomicAssumption)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AnyDependentPeople)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AnyPublicOrPrivateEntity)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1PublicTransportAccess)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1PreferredTransportationMode)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1TransportDifficulty)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1EconomicResourcesForTransport)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1DebtOrHelpDueToTransport)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1TransportSubsidyReceived)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1SafetyTransportationConcern)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1FoodAccessFrequency)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AggressorFoodRestriction)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1FoodIncomeSupport)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1JuridicalAssistanceReceived)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1VictimRepresentation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1PsychosocialSupportReceived)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1EmergencyEmotionalCrisis)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AggressorSameResidence)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AggressorLocationKnown)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1AnyTypeOfAssistanceReceived)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1CompensationFundInsuranceCoverage)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1FuneralSubsidyReceived)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1FuneralFundsAvailable)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1FuneralCostValue)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1RenameReputationImpact)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicide.FeminicideForm1.FeminicideForm1ActionPlan)
}

func getAndVerifyFeminicideEnumsMultiple(feminicideForm1 *salvia_daos.FeminicideForm1DTO, collectedErrors map[string]map[string]string) bool {
	var opRes bool = true
	var err error
	var entityIntervened []salvia_daos.VictimCaseForm2EnumsDTO = []salvia_daos.VictimCaseForm2EnumsDTO{}
	for idx := range feminicideForm1.FeminicideForm1EntityIntervened {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideForm1.FeminicideForm1EntityIntervened[idx])
		if err != nil {
			opRes = false
			break
		} else {
			var enumTmp salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsCategory: "feminicide_form1_entities_intervend" + feminicideForm1.FeminicideForm1EntityIntervened[idx].VictimCaseForm2EnumsCode}
			var eTmp []salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.GetLocalVictimCaseForm2EnumsByCategory(enumTmp)
			entityIntervened = append(entityIntervened, eTmp...)
		}
	}
	if len(feminicideForm1.FeminicideForm1EntityIntervened) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideForm1, "FeminicideForm1EntityIntervened", "json"), "feminicide_form1_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

// En enums están todos los del caso actual, hay que extraer los necesarios para llenar el objeto
func loadFeminicideEnumsMultiple(feminicideForm1 *salvia_daos.FeminicideForm1DTO, enums []salvia_daos.VictimCaseForm2EnumsDTO) {
	var err error

	for _, e := range enums {
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsById(&e); err == nil {
			switch e.VictimCaseForm2EnumsCategory {
			case "feminicide_form1_entities_intervend":
				feminicideForm1.FeminicideForm1EntityIntervened = append(feminicideForm1.FeminicideForm1EntityIntervened, e)
			}
		}
	}
}

// Dado que ya están validados, entonces se relacionan los enums con contenido
func setFeminicideEnumsMultiple(feminicideForm1 *salvia_daos.FeminicideForm1DTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	var err error
	for _, e := range feminicideForm1.FeminicideForm1EntityIntervened {
		var rel salvia_daos.RelVictimCaseForm2EnumsFeminicideForm1DTO
		rel.RelVictimCaseForm2EnumsFeminicideForm1Form = *feminicideForm1
		rel.RelVictimCaseForm2EnumsFeminicideForm1Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsFeminicideForm1(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	return nil
}
