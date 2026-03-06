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

type FeminicideRiskRequest struct {
	FeminicideRisk salvia_daos.FeminicideRiskDTO `json:"feminicideRisk"`
}

type FeminicideRiskForm1Request struct {
	Form salvia_daos.FeminicideRiskForm1DTO `json:"form"`
}

// SetFeminicideRisk crea un nuevo caso de víctima a partir de la entrada JSON.
// Realiza la validación de los datos, crea o actualiza usuarios y relaciones,
// y gestiona la transacción en la base de datos. Devuelve un código HTTP y
// un mensaje en formato JSON (éxito o error).
func SetFeminicideRisk(dataInput string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de variables y obtención de conexión a la BD.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	// Se libera la conexión al finalizar la función.
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores durante la validación del DTO.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos para el caso, contacto y dueño.
	var feminicideRiskRequest FeminicideRiskRequest = FeminicideRiskRequest{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	// Variable para almacenar el DTO en formato mapa.
	var dtoMap map[string]interface{} = nil

	// Se parsea el JSON de entrada a un mapa.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FeminicideRiskJSONName, common_config.Locale, collectedErrors)
	utils.JSONToStruct(dataInput, &feminicideRiskRequest)

	// Se asignan valores por defecto al objeto FeminicideRisk.
	salvia_daos.SetFeminicideRiskDefaults(&feminicideRiskRequest.FeminicideRisk, common_dao.SQL_INSERT)
	salvia_daos.SetFeminicideRiskForm1Defaults(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, common_dao.SQL_INSERT)

	// Verificamos que los datos tengan la estructura esperada.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar (obligatorios y opcionales).
		var checkFields map[string]bool = map[string]bool{
			"FeminicideRiskNames":     true,
			"FeminicideRiskLastNames": true,
			"FeminicideRiskDocType":   true,
			"FeminicideRiskDocNumber": true,
		}
		// Validación de los campos del JSON contra la definición del DTO.
		utils.ValidateJSONInput(&feminicideRiskRequest.FeminicideRisk, dtoMap, salvia_daos.FeminicideRiskJSONName, salvia_daos.FeminicideRiskFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"FeminicideRiskForm1VictimIdentityName":                       true,
			"FeminicideRiskForm1BirthDate":                                true,
			"FeminicideRiskForm1VictimAddress":                            true,
			"FeminicideRiskForm1VictimZone":                               true,
			"FeminicideRiskForm1VictimLivingTown":                         true,
			"FeminicideRiskForm1VictimSGSSSAffiliation":                   true,
			"FeminicideRiskForm1ContactPhone":                             true,
			"FeminicideRiskForm1ContactEmergencyContactNames":             true,
			"FeminicideRiskForm1ContactEmergencyContactNumber":            true,
			"FeminicideRiskForm1VictimMaritalStatus":                      true,
			"FeminicideRiskForm1VictimSex":                                true,
			"FeminicideRiskForm1VictimGenderIdentity":                     true,
			"FeminicideRiskForm1VictimSexualOrientation":                  true,
			"FeminicideRiskForm1VictimEthnicAffiliation":                  true,
			"FeminicideRiskForm1VictimIsMigrant":                          true,
			"FeminicideRiskForm1VictimMigrationStatus":                    true,
			"FeminicideRiskForm1VictimMaxEducationLevel":                  true,
			"FeminicideRiskForm1VictimIsSpecialPopulation":                true,
			"FeminicideRiskForm1VictimCurrentlyHasJob":                    true,
			"FeminicideRiskForm1VictimJobExplanation":                     true,
			"FeminicideRiskForm1VictimAbandonedJobDueToRisk":              true,
			"FeminicideRiskForm1VictimMainOccupation":                     true,
			"FeminicideRiskForm1VictimHasDisability":                      true,
			"FeminicideRiskForm1VictimDisabilityType":                     true,
			"FeminicideRiskForm1VictimRiskDescription":                    true,
			"FeminicideRiskForm1VictimAdditionalInfo":                     true,
			"FeminicideRiskForm1VictimReceivedHelpOrSubsidy":              true,
			"FeminicideRiskForm1VictimIsEconomicProvider":                 true,
			"FeminicideRiskForm1VictimEconomicProviderExplanation":        true,
			"FeminicideRiskForm1VictimHasFamiliarSupport":                 true,
			"FeminicideRiskForm1VictimFamiliarSupportType":                true,
			"FeminicideRiskForm1VictimPublicTransportAccess":              true,
			"FeminicideRiskForm1VictimCommonTransportMode":                true,
			"FeminicideRiskForm1VictimTransportModeExplanation":           true,
			"FeminicideRiskForm1VictimEstimatedTravelCost":                true,
			"FeminicideRiskForm1VictimDifficultiesWithTransport":          true,
			"FeminicideRiskForm1VictimDifficultiesExplanation":            true,
			"FeminicideRiskForm1VictimEconomicResourcesForTransport":      true,
			"FeminicideRiskForm1VictimEconomicResourcesExplanation":       true,
			"FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport":        true,
			"FeminicideRiskForm1VictimDebtImpactDetails":                  true,
			"FeminicideRiskForm1VictimReceivedTransportSubsidy":           true,
			"FeminicideRiskForm1VictimTransportSubsidyExplanation":        true,
			"FeminicideRiskForm1VictimSafetyAvoidedTransport":             true,
			"FeminicideRiskForm1VictimSafetyAvoidedExplanation":           true,
			"FeminicideRiskForm1VictimFoodAccessFrequency":                true,
			"FeminicideRiskForm1VictimFoodAccessExplanation":              true,
			"FeminicideRiskForm1VictimAgressorFoodRestriction":            true,
			"FeminicideRiskForm1VictimAgressorFoodRestrictionExplanation": true,
			"FeminicideRiskForm1VictimFamilyFixedIncome":                  true,
			"FeminicideRiskForm1VictimFamilyFixedIncomeExplanation":       true,
			"FeminicideRiskForm1VictimIsOnlyProviderForFood":              true,
			"FeminicideRiskForm1VictimIsOnlyProviderExplanation":          true,
			"FeminicideRiskForm1VictimJuridicalAssistanceReceived":        true,
			"FeminicideRiskForm1VictimJuridicalAssistanceExplain":         true,
			"FeminicideRiskForm1VictimWantsJuridicalAssistance":           true,
			"FeminicideRiskForm1VictimRepresentationsOfVictims":           true,
			"FeminicideRiskForm1VictimRepresentationsExplanation":         true,
			"FeminicideRiskForm1VictimPsychosocialSupportReceived":        true,
			"FeminicideRiskForm1VictimPsychosocialSupportReceivedExplain": true,
			"FeminicideRiskForm1VictimUrgentEmotionalCrisis":              true,
			"FeminicideRiskForm1VictimUrgentCrisisExplanation":            true,
			"FeminicideRiskForm1AggressorSameResidence":                   true,
			"FeminicideRiskForm1AggressorSameResidenceExplanation":        true,
			"FeminicideRiskForm1AggressorKnowsVictimLocation":             true,
			"FeminicideRiskForm1AggressorKnowsVictimLocationExplanation":  true,
			"FeminicideRiskForm1VictimHousingHelpReceived":                true,
			"FeminicideRiskForm1VictimHousingHelpExplain":                 true,
			"FeminicideRiskForm1VictimAbandonClothing":                    true,
			"FeminicideRiskForm1VictimAbandonClothingExplain":             true,
			"FeminicideRiskForm1VictimClothingHelpReceived":               true,
			"FeminicideRiskForm1VictimClothingHelpExplain":                true,
			"FeminicideRiskForm1VictimOtherNeeds":                         true,
			"FeminicideRiskForm1InterviewDate":                            true,
			"FeminicideRiskForm1Summary":                                  true,
		}

		utils.ValidateJSONInput(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, dtoMap, salvia_daos.FeminicideRiskJSONName+"."+salvia_daos.FeminicideRiskForm1JSONName, salvia_daos.FeminicideRiskForm1FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.FeminicideRiskJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

	default:
		// Si el JSON no tiene la estructura esperada, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.FeminicideRiskJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	//Ahora se verifican los enums únicos
	getAndVerifyFeminicideRiskEnums(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, collectedErrors, true)

	//Ahora se verifican los enums múltiples que se asociarán al form2
	getAndVerifyFeminicideRiskEnumsMultiple(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, collectedErrors)

	//Ahora algunos campos opcionales
	/*if feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1IncomeGenerationMethod.FeminicideRiskForm1EnumsCode == "pr" {
		if feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsStartTime.IsZero() {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, "FeminicideRiskForm1FactsStartTime", "json"), "common_validation_field_date_error", "common_global_error", common_config.Locale)
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
	if err = salvia_daos.SetFeminicideRisk(&feminicideRiskRequest.FeminicideRisk, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FeminicideRisk = feminicideRiskRequest.FeminicideRisk

	if err = salvia_daos.SetFeminicideRiskForm1(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Ahora se insertan todos los campos múltiples
	//Ahora se verifican los enums múltiples que se asociarán al form2
	if err = setFeminicideRiskEnumsMultiple(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, connData, dbClientConfig, dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se asigna la autoría del caso al usuario en sesión.
	var relCaseOwnerFeminicideRisk salvia_daos.RelCaseOwnerFeminicideRiskDTO = salvia_daos.RelCaseOwnerFeminicideRiskDTO{}
	salvia_daos.SetRelCaseOwnerFeminicideRiskDefaults(&relCaseOwnerFeminicideRisk, common_dao.SQL_INSERT)
	relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_CaseOwner = owner.CaseOwnerId
	relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_FeminicideRisk = feminicideRiskRequest.FeminicideRisk.FeminicideRiskId

	if err = salvia_daos.SetRelCaseOwnerFeminicideRisk(&relCaseOwnerFeminicideRisk, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
	if err = salvia_daos.UpdateFeminicideRiskOwnersAndRolesByFeminicideRiskId(feminicideRiskRequest.FeminicideRisk.FeminicideRiskId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se realiza el commit de la transacción y se retorna el resultado.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}

// UpdateFeminicideRisk actualiza los datos de un caso de víctima existente.
// Recibe la entrada JSON, el identificador del caso (feminicideRiskICode) y realiza la actualización
// validando el DTO, actualizando datos y gestionando la transacción.
/*
func UpdateFeminicideRisk(dataInput string, feminicideRiskICode string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicializa variables y establece la conexión.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores de validación.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos: caso, dueño y el caso a actualizar.
	var feminicideRiskRequest FeminicideRiskRequest = FeminicideRiskRequest{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	var feminicideRiskToUpdate salvia_daos.FeminicideRiskDTO = salvia_daos.FeminicideRiskDTO{}
	var feminicideRiskFormToUpdate salvia_daos.FeminicideRiskForm1DTO = salvia_daos.FeminicideRiskForm1DTO{}
	var code int

	// Slices para gestionar los momentos asociados al caso.
	var momentsToDelete []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}
	var moments []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}

	// Variable para almacenar el DTO en forma de mapa.
	var dtoMap map[string]interface{} = nil

	// Se obtiene el mapa DTO a partir de la entrada JSON.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FeminicideRiskJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &feminicideRiskRequest)

	// Se establecen valores por defecto para el caso (modo actualización).
	salvia_daos.SetFeminicideRiskDefaults(&feminicideRiskRequest.FeminicideRisk, common_dao.SQL_UPDATE, s)
	salvia_daos.SetFeminicideRiskForm1Defaults(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, common_dao.SQL_INSERT)

	// Se verifica la estructura del DTO.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar.
		var checkFields map[string]bool = map[string]bool{
			"FeminicideRiskNames":     true,
			"FeminicideRiskLastNames": true,
			"FeminicideRiskDocType":   true,
			"FeminicideRiskDocNumber": true,
			"FeminicideRiskTownCode":  true,
		}
		// Validación de los datos del caso.
		utils.ValidateJSONInput(&feminicideRiskRequest.FeminicideRisk, dtoMap, salvia_daos.FeminicideRiskJSONName, salvia_daos.FeminicideRiskFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"FeminicideRiskForm1Nick":                              true,
			"FeminicideRiskForm1BirthDate":                         true,
			"FeminicideRiskForm1ViolenceTownCode":                  true,
			"FeminicideRiskForm1Address":                           true,
			"FeminicideRiskForm1LivingLatitude":                    false,
			"FeminicideRiskForm1LivingLongitude":                   false,
			"FeminicideRiskForm1Phone":                             true,
			"FeminicideRiskForm1Email":                             true,
			"FeminicideRiskForm1GenderIdentity":                    true,
			"FeminicideRiskForm1SexualOrientation":                 true,
			"FeminicideRiskForm1Origin":                            true,
			"FeminicideRiskForm1Occupation":                        true,
			"FeminicideRiskForm1OccupationOther":                   false,
			"FeminicideRiskForm1VictimEthnicGroup":                 true,
			"FeminicideRiskForm1VictimEthnicGroupOther":            false,
			"FeminicideRiskForm1VictimContactNames":                true,
			"FeminicideRiskForm1VictimContactPhone":                true,
			"FeminicideRiskForm1VictimContactKinship":              true,
			"FeminicideRiskForm1VictimNumChildren":                 false,
			"FeminicideRiskForm1VictimMaritalStatus":               true,
			"FeminicideRiskForm1VictimMaritalStatusOther":          false,
			"FeminicideRiskForm1VictimChildrenAge":                 false,
			"FeminicideRiskForm1VictimDisability":                  true,
			"FeminicideRiskForm1VictimSpecialSupport":              true,
			"FeminicideRiskForm1FactsOccurrence":                   true,
			"FeminicideRiskForm1FactsStartTime":                    true,
			"FeminicideRiskForm1FactsEndTime":                      true,
			"FeminicideRiskForm1FactsWeekday":                      true,
			"FeminicideRiskForm1FactsDate":                         true,
			"FeminicideRiskForm1FactsDescription":                  true,
			"FeminicideRiskForm1VictimViolenceExperienced":         true,
			"FeminicideRiskForm1VictimViolenceExperiencedOther":    false,
			"FeminicideRiskForm1VictimViolenceScope":               true,
			"FeminicideRiskForm1VictimFemicideRisk":                true,
			"FeminicideRiskForm1VictimAggressor":                   true,
			"FeminicideRiskForm1VictimRelationshipWithAggressor":   true,
			"FeminicideRiskForm1VictimAggressorName":               false,
			"FeminicideRiskForm1VictimAggressorDocType":            false,
			"FeminicideRiskForm1VictimAggressorDocNumber":          false,
			"FeminicideRiskForm1VictimAggressorAddress":            false,
			"FeminicideRiskForm1VictimAggressorPhone":              false,
			"FeminicideRiskForm1Age":                               true,
			"FeminicideRiskForm1VictimNationality":                 true,
			"FeminicideRiskForm1VictimNationalityOther":            false,
			"FeminicideRiskForm1VictimForeignerImmigrationStatus":  false,
			"FeminicideRiskForm1VictimGender":                      true,
			"FeminicideRiskForm1VictimGenderIdentityOther":         false,
			"FeminicideRiskForm1VictimSexualOrientationOther":      false,
			"FeminicideRiskForm1VictimDependents":                  true,
			"FeminicideRiskForm1VictimPreviouslyReportedSituation": true,
			"FeminicideRiskForm1VictimIfPreviouslyReported":        false,
			"FeminicideRiskForm1VictimIfAfro":                      false,
			"FeminicideRiskForm1VictimIfIndigenous":                false,
			"FeminicideRiskForm1VictimIfIndigenousTongue":          false,
			"FeminicideRiskForm1VictimIfPeasant":                   true,
			"FeminicideRiskForm1VictimIfArmedConflict":             true,
			"FeminicideRiskForm1VictimViolenceScene":               true,
			"FeminicideRiskForm1PhysicalViolenceIncreased":         true,
			"FeminicideRiskForm1SeparatedFromPartnerLastYear":      true,
			"FeminicideRiskForm1ThreatenedWithWeapon":              true,
			"FeminicideRiskForm1ThreatenedToKillOrHarmChildren":    true,
			"FeminicideRiskForm1JealousAndViolent":                 true,
			"FeminicideRiskForm1BelievesCapableOfKilling":          true,
		}

		utils.ValidateJSONInput(&feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1, dtoMap, salvia_daos.FeminicideRiskJSONName+"."+salvia_daos.FeminicideRiskForm1JSONName, salvia_daos.FeminicideRiskForm1FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.FeminicideRiskJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}


	default:
		// Si el JSON no cumple la estructura, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.FeminicideRiskJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si existen errores, se retorna un BadRequest.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se obtiene el caso a actualizar a partir del identificador.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideRiskICode"},
		AttrsValue: []interface{}{feminicideRiskICode},
	}

	if err = salvia_daos.GetFeminicideRisk(by, &feminicideRiskToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideRiskForm1FeminicideRisk"},
		AttrsValue: []interface{}{feminicideRiskToUpdate.FeminicideRiskId},
	}

	if err = salvia_daos.GetFeminicideRiskForm1(by, &feminicideRiskFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Verifica si se produjo un cambio en el documento de identidad.
	if feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocType != feminicideRiskToUpdate.FeminicideRiskDocType || feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocNumber != feminicideRiskToUpdate.FeminicideRiskDocNumber {
		var oldProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}
		var newProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
			AttrsValue: []interface{}{feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocType, feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocNumber},
		}

		err = security_daos.GetGeneralUserProfile(by, &newProfile, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// No existe un usuario con el nuevo documento; se actualiza el perfil del usuario.
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
				AttrsValue: []interface{}{feminicideRiskToUpdate.FeminicideRiskDocType, feminicideRiskToUpdate.FeminicideRiskDocNumber},
			}
			err = security_daos.GetGeneralUserProfile(by, &oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_feminicide_risk_find_profile_fail"]
			}
			// Se actualiza el perfil con el nuevo documento.
			security_daos.SetGeneralUserProfileDefaults(&oldProfile, common_dao.SQL_UPDATE)
			oldProfile.GeneralUserProfileDocType = feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocType
			oldProfile.GeneralUserProfileDocNumber = feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocNumber

			err = security_daos.UpdateGeneralUserProfile(&oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_feminicide_risk_user_not_found"]
			}
		} else {
			// Si ya existe un perfil con el nuevo documento, se retorna un error.
			utils.SetError(collectedErrors, salvia_daos.FeminicideRiskJSONName, utils.GetTag(&feminicideRiskRequest.FeminicideRisk, "FeminicideRiskDocNumber", "json"),
				"update_feminicide_risk_find_profile_already_exist", "salvia_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	// Se inicia la transacción para la actualización.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualizan los campos del caso con los nuevos valores.
	feminicideRiskToUpdate.FeminicideRiskNames = feminicideRiskRequest.FeminicideRisk.FeminicideRiskNames
	feminicideRiskToUpdate.FeminicideRiskLastNames = feminicideRiskRequest.FeminicideRisk.FeminicideRiskLastNames
	feminicideRiskToUpdate.FeminicideRiskDocType = feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocType
	feminicideRiskToUpdate.FeminicideRiskDocNumber = feminicideRiskRequest.FeminicideRisk.FeminicideRiskDocNumber
	feminicideRiskToUpdate.FeminicideRiskTownCode = feminicideRiskRequest.FeminicideRisk.FeminicideRiskTownCode

	feminicideRiskFormToUpdate.FeminicideRiskForm1Nick = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Nick
	feminicideRiskFormToUpdate.FeminicideRiskForm1BirthDate = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1BirthDate
	feminicideRiskFormToUpdate.FeminicideRiskForm1Address = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Address
	feminicideRiskFormToUpdate.FeminicideRiskForm1LivingLatitude = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1LivingLatitude
	feminicideRiskFormToUpdate.FeminicideRiskForm1LivingLongitude = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1LivingLongitude
	feminicideRiskFormToUpdate.FeminicideRiskForm1Phone = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Phone
	feminicideRiskFormToUpdate.FeminicideRiskForm1Email = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Email
	feminicideRiskFormToUpdate.FeminicideRiskForm1GenderIdentity = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1GenderIdentity
	feminicideRiskFormToUpdate.FeminicideRiskForm1SexualOrientation = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1SexualOrientation
	feminicideRiskFormToUpdate.FeminicideRiskForm1Origin = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Origin
	feminicideRiskFormToUpdate.FeminicideRiskForm1Occupation = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Occupation
	feminicideRiskFormToUpdate.FeminicideRiskForm1OccupationOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1OccupationOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimEthnicGroup = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimEthnicGroup
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimEthnicGroupOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimEthnicGroupOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimContactNames = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimContactNames
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimContactPhone = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimContactPhone
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimContactKinship = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimContactKinship
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimNumChildren = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimNumChildren
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimMaritalStatus = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimMaritalStatus
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimMaritalStatusOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimMaritalStatusOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimChildrenAge = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimChildrenAge
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimDisability = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimDisability
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimSpecialSupport = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimSpecialSupport
	feminicideRiskFormToUpdate.FeminicideRiskForm1FactsOccurrence = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsOccurrence
	feminicideRiskFormToUpdate.FeminicideRiskForm1FactsStartTime = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsStartTime
	feminicideRiskFormToUpdate.FeminicideRiskForm1FactsEndTime = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsEndTime
	feminicideRiskFormToUpdate.FeminicideRiskForm1FactsWeekday = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsWeekday
	feminicideRiskFormToUpdate.FeminicideRiskForm1FactsDate = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsDate
	feminicideRiskFormToUpdate.FeminicideRiskForm1FactsDescription = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1FactsDescription
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimViolenceExperienced = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimViolenceExperienced
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimViolenceExperiencedOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimViolenceExperiencedOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimViolenceScope = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimViolenceScope
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimFemicideRisk = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimFemicideRisk
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimAggressor = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAggressor
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimRelationshipWithAggressor = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimRelationshipWithAggressor
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimAggressorName = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAggressorName
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimAggressorDocType = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAggressorDocType
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimAggressorDocNumber = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAggressorDocNumber
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimAggressorAddress = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAggressorAddress
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimAggressorPhone = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAggressorPhone

	feminicideRiskFormToUpdate.FeminicideRiskForm1Age = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Age
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimNationality = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimNationality
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimNationalityOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimNationalityOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimForeignerImmigrationStatus = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimForeignerImmigrationStatus
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimGender = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimGender
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimGenderIdentityOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimGenderIdentityOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimSexualOrientationOther = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimSexualOrientationOther
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimDependents = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimDependents
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimPreviouslyReportedSituation = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimPreviouslyReportedSituation
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimIfPreviouslyReported = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIfPreviouslyReported

	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimIfAfro = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIfAfro
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimIfIndigenous = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIfIndigenous
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimIfIndigenousTongue = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIfIndigenousTongue
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimIfPeasant = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIfPeasant
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimIfArmedConflict = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIfArmedConflict
	feminicideRiskFormToUpdate.FeminicideRiskForm1VictimViolenceScene = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimViolenceScene

	feminicideRiskFormToUpdate.FeminicideRiskForm1PhysicalViolenceIncreased = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1PhysicalViolenceIncreased
	feminicideRiskFormToUpdate.FeminicideRiskForm1SeparatedFromPartnerLastYear = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1SeparatedFromPartnerLastYear
	feminicideRiskFormToUpdate.FeminicideRiskForm1ThreatenedWithWeapon = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1ThreatenedWithWeapon
	feminicideRiskFormToUpdate.FeminicideRiskForm1ThreatenedToKillOrHarmChildren = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1ThreatenedToKillOrHarmChildren
	feminicideRiskFormToUpdate.FeminicideRiskForm1JealousAndViolent = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1JealousAndViolent
	feminicideRiskFormToUpdate.FeminicideRiskForm1BelievesCapableOfKilling = feminicideRiskRequest.FeminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1BelievesCapableOfKilling

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
	if err = salvia_daos.UpdateFeminicideRisk(&feminicideRiskToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if err = salvia_daos.UpdateFeminicideRiskForm1(&feminicideRiskFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen los momentos (acciones o eventos) asociados al caso.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentFeminicideRisk"},
		AttrsValue: []interface{}{feminicideRiskToUpdate.FeminicideRiskId},
	}

	momentsToDelete, err = salvia_daos.GetMoments(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se procesan las ramas para determinar cuáles momentos eliminar y cuáles crear.
	for mcode, sectors := range feminicideRiskRequest.FeminicideRisk.FeminicideRiskEntityBranches {
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
					moment.MomentFeminicideRisk = feminicideRiskToUpdate

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
	var relCaseOwners []salvia_daos.RelCaseOwnerFeminicideRiskDTO = []salvia_daos.RelCaseOwnerFeminicideRiskDTO{}
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelCaseOwnerFeminicideRisk_CaseOwner", "RelCaseOwnerFeminicideRisk_FeminicideRisk"},
		AttrsValue: []interface{}{owner.CaseOwnerId, feminicideRiskToUpdate.FeminicideRiskId},
	}
	relCaseOwners, err = salvia_daos.GetRelCaseOwnerFeminicideRisks(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Si no existe la relación, se crea y se actualiza el contador del dueño.
	if len(relCaseOwners) == 0 {
		var relCaseOwnerFeminicideRisk salvia_daos.RelCaseOwnerFeminicideRiskDTO = salvia_daos.RelCaseOwnerFeminicideRiskDTO{}
		salvia_daos.SetRelCaseOwnerFeminicideRiskDefaults(&relCaseOwnerFeminicideRisk, common_dao.SQL_INSERT)
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_CaseOwner = owner.CaseOwnerId
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_FeminicideRisk = feminicideRiskToUpdate.FeminicideRiskId

		if err = salvia_daos.SetRelCaseOwnerFeminicideRisk(&relCaseOwnerFeminicideRisk, connData, &dbClientConfig, &dbServerConfig); err != nil {
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
	if err = salvia_daos.UpdateFeminicideRiskOwnersAndRolesByFeminicideRiskId(feminicideRiskToUpdate.FeminicideRiskId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se finaliza la transacción.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideRiskRequest.FeminicideRisk.FeminicideRiskNewUser)
}*/

// GetFeminicideRiskByICode retorna un caso de víctima a partir de su identificador único (ICode).
// Además, recupera información geográfica relacionada (municipio, ciudad, departamento)
// y los momentos (logs) asociados al caso.
func GetFeminicideRiskByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.FeminicideRiskDTO) {
	var err error = nil

	// Se libera la conexión si ésta no es válida.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el caso.
	var feminicideRisk salvia_daos.FeminicideRiskDTO = salvia_daos.FeminicideRiskDTO{}
	var feminicideRiskForm1 salvia_daos.FeminicideRiskForm1DTO = salvia_daos.FeminicideRiskForm1DTO{}

	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, feminicideRisk
		}

	} else {
		var town security_daos.TownDTO

		var city security_daos.CityDTO

		var department security_daos.DepartmentDTO

		// Se consulta el caso por su identificador.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"FeminicideRiskICode"},
			AttrsValue: []interface{}{id},
		}

		err = salvia_daos.GetFeminicideRisk(by, &feminicideRisk, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), feminicideRisk
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"FeminicideRiskForm1FeminicideRisk"},
			AttrsValue: []interface{}{feminicideRisk.FeminicideRiskId},
		}

		if err = salvia_daos.GetFeminicideRiskForm1(by, &feminicideRiskForm1, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideRiskDTO{}
		}

		feminicideRisk.FeminicideRiskForm1 = feminicideRiskForm1
		// Se recupera la información geográfica de la víctima y familiar
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownCode"},
			AttrsValue: []interface{}{feminicideRiskForm1.FeminicideRiskForm1VictimLivingTownCode},
		}
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideRiskDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town.TownCity},
		}
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideRiskDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city.CityDepartment},
		}
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideRiskDTO{}
		}

		//Ahora se cargan los enums sencillos
		loadFeminicideRiskEnums(&feminicideRisk)

		//Ahora se cargan los enums múltiples
		var enums []salvia_daos.VictimCaseForm2EnumsDTO

		enums, err = salvia_daos.GetVictimCasesForm2EnumsByFeminicideRiskForm1Id(feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1Id, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.FeminicideRiskDTO{}
		}

		loadFeminicideRiskEnumsMultiple(&feminicideRisk.FeminicideRiskForm1, enums)

		// Se asignan los datos geográficos, los momentos y el seguimiento al caso.
		feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimLivingDepartment = department
		feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimLivingCity = city
		feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimLivingTown = town

		return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideRisk), feminicideRisk

	}
	return http.StatusInternalServerError, "", salvia_daos.FeminicideRiskDTO{}
}

// UpdateFeminicideRiskStatusByICode actualiza el estado de un caso de víctima basado en su ICode.
/*func UpdateFeminicideRiskStatusByICode(newStatus string, id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicideRiskToUpdate salvia_daos.FeminicideRiskDTO = salvia_daos.FeminicideRiskDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideRiskICode"},
		AttrsValue: []interface{}{id},
	}

	if err = salvia_daos.GetFeminicideRisk(by, &feminicideRiskToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se asigna el nuevo estado.
	feminicideRiskToUpdate.FeminicideRiskStatus = newStatus

	if err = salvia_daos.UpdateFeminicideRisk(&feminicideRiskToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideRiskToUpdate)
}*/

// GetFeminicideRisksByTownCode obtiene todos los casos de víctimas asociados a un código de municipio.
func GetFeminicideRisksByTownCode(townCode string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
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

	var feminicideRisks []salvia_daos.FeminicideRiskDTO = []salvia_daos.FeminicideRiskDTO{}

	// Se consulta la BD utilizando el código del municipio.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideRiskTownCode"},
		AttrsValue: []interface{}{townCode},
	}

	if feminicideRisks, count, err = salvia_daos.GetFeminicideRisks(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideRisks), count
}

// GetDepartmentFeminicideRisksByTownCode obtiene todos los casos de víctimas asociados al departamento relacionado con el código de municipio.
/*func GetDepartmentFeminicideRisksByTownCode(townCode string, feminicideRiskStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
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

	var feminicideRisks []salvia_daos.FeminicideRiskDTO = []salvia_daos.FeminicideRiskDTO{}

	if feminicideRisks, count, err = salvia_daos.GetFeminicideRisksByDepartmentICode(department.DepartmentICode, feminicideRiskStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(feminicideRisks), count
}*/

// GetFeminicideRisksByDocument obtiene casos de víctimas filtrados por tipo y número de documento.
func GetFeminicideRisksByDocument(docType string, docNumber string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicideRisks []salvia_daos.FeminicideRiskDTO = []salvia_daos.FeminicideRiskDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FeminicideRiskDocType", "FeminicideRiskDocNumber"},
		AttrsValue: []interface{}{docType, docNumber},
	}

	if feminicideRisks, count, err = salvia_daos.GetFeminicideRisks(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{feminicideRisks, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}

/*
// GetFeminicideRisksByOwnerUserICode obtiene casos de víctimas asociados al ICode del usuario dueño.
func GetFeminicideRisksByOwnerUserICode(userIcode string, feminicideRiskStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var feminicideRisks []salvia_daos.FeminicideRiskDTO = []salvia_daos.FeminicideRiskDTO{}

	if feminicideRisks, count, err = salvia_daos.GetFeminicideRisksByOwnerUserICode(userIcode, feminicideRiskStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{feminicideRisks, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}*/

// GetFeminicideRiskByAll retorna todos los casos de víctimas almacenados en la base de datos.

func GetFeminicideRiskByAll(page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	cases, count, err := salvia_daos.GetAllFeminicideRisks(page, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{cases, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"})
	}
	return resCode, resData, count
}

func getAndVerifyFeminicideRiskEnums(feminicideRiskForm1 *salvia_daos.FeminicideRiskForm1DTO, collectedErrors map[string]map[string]string, verify bool) bool {
	var opRes bool = true

	err := salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimZone)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimZone", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimSGSSSAffiliation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimSGSSSAffiliation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimMaritalStatus)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimMaritalStatus", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimSex)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimSex", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimGenderIdentity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimGenderIdentity", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimSexualOrientation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimSexualOrientation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimEthnicAffiliation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimEthnicAffiliation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimEthnicAffiliation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimEthnicAffiliation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if feminicideRiskForm1.FeminicideRiskForm1VictimEthnicAffiliation.VictimCaseForm2EnumsCode == "in" {
		// VictimCaseForm2IndigenousPeople
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimIndigenousPeople)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimIndigenousPeople", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimIsMigrant)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimIsMigrant", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimMigrationStatus)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimMigrationStatus", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimMaxEducationLevel)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimMaxEducationLevel", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimIsSpecialPopulation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimIsSpecialPopulation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimCurrentlyHasJob)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimCurrentlyHasJob", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimAbandonedJobDueToRisk)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimAbandonedJobDueToRisk", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimMainOccupation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimMainOccupation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimHasDisability)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimHasDisability", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimDisabilityType)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimDisabilityType", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimReceivedHelpOrSubsidy)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimReceivedHelpOrSubsidy", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimIsEconomicProvider)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimIsEconomicProvider", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimHasFamiliarSupport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimHasFamiliarSupport", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimFamiliarSupportType)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimFamiliarSupportType", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimPublicTransportAccess)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimPublicTransportAccess", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimCommonTransportMode)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimCommonTransportMode", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimDifficultiesWithTransport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimDifficultiesWithTransport", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimEconomicResourcesForTransport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimEconomicResourcesForTransport", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimReceivedTransportSubsidy)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimReceivedTransportSubsidy", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimSafetyAvoidedTransport)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimSafetyAvoidedTransport", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimFoodAccessFrequency)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimFoodAccessFrequency", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimAgressorFoodRestriction)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimAgressorFoodRestriction", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimFamilyFixedIncome)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimFamilyFixedIncome", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimIsOnlyProviderForFood)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimIsOnlyProviderForFood", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimJuridicalAssistanceReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimJuridicalAssistanceReceived", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimWantsJuridicalAssistance)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimWantsJuridicalAssistance", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimRepresentationsOfVictims)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimRepresentationsOfVictims", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimPsychosocialSupportReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimPsychosocialSupportReceived", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimUrgentEmotionalCrisis)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimUrgentEmotionalCrisis", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AggressorSameResidence)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AggressorSameResidence", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AggressorKnowsVictimLocation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AggressorKnowsVictimLocation", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimHousingHelpReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimHousingHelpReceived", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimAbandonClothing)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimAbandonClothing", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1VictimClothingHelpReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1VictimClothingHelpReceived", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceived)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AnyAssistanceReceived", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedCityHall)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AnyAssistanceReceivedCityHall", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AnyAssistanceReceivedWomensOffice", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AnyAssistanceReceivedOtherEntity", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedOther)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1AnyAssistanceReceivedOther", "json"), "feminicide_risk_form1_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

func loadFeminicideRiskEnums(feminicideRisk *salvia_daos.FeminicideRiskDTO) {

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimZone)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimSGSSSAffiliation)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimMaritalStatus)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimSex)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimGenderIdentity)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimSexualOrientation)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimEthnicAffiliation)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIndigenousPeople)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIsMigrant)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimMigrationStatus)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimMaxEducationLevel)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIsSpecialPopulation)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimCurrentlyHasJob)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAbandonedJobDueToRisk)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimMainOccupation)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimHasDisability)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimDisabilityType)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimReceivedHelpOrSubsidy)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIsEconomicProvider)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimHasFamiliarSupport)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimFamiliarSupportType)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimPublicTransportAccess)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimCommonTransportMode)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimDifficultiesWithTransport)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimEconomicResourcesForTransport)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimHasDebtOrHelpDueToTransport)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimReceivedTransportSubsidy)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimSafetyAvoidedTransport)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimFoodAccessFrequency)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAgressorFoodRestriction)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimFamilyFixedIncome)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimIsOnlyProviderForFood)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimJuridicalAssistanceReceived)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimWantsJuridicalAssistance)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimRepresentationsOfVictims)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimPsychosocialSupportReceived)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimUrgentEmotionalCrisis)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AggressorSameResidence)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AggressorKnowsVictimLocation)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimHousingHelpReceived)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimAbandonClothing)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1VictimClothingHelpReceived)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceived)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedCityHall)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedWomensOffice)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedOtherEntity)
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&feminicideRisk.FeminicideRiskForm1.FeminicideRiskForm1AnyAssistanceReceivedOther)

}

func getAndVerifyFeminicideRiskEnumsMultiple(feminicideRiskForm1 *salvia_daos.FeminicideRiskForm1DTO, collectedErrors map[string]map[string]string) bool {
	var opRes bool = true
	var err error
	var entityIntervened []salvia_daos.VictimCaseForm2EnumsDTO = []salvia_daos.VictimCaseForm2EnumsDTO{}
	for idx := range feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople[idx])
		if err != nil {
			opRes = false
			break
		} else {
			var enumTmp salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsCategory: "feminicide_risk_form1_financially_dependent_people" + feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople[idx].VictimCaseForm2EnumsCode}
			var eTmp []salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.GetLocalVictimCaseForm2EnumsByCategory(enumTmp)
			entityIntervened = append(entityIntervened, eTmp...)
		}
	}
	if len(feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1FinanciallyDependentPeople", "json"), "feminicide_risk_form1_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly[idx])
		if err != nil {
			opRes = false
			break
		} else {
			var enumTmp salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsCategory: "feminicide_risk_form1_places_visit_regularly" + feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly[idx].VictimCaseForm2EnumsCode}
			var eTmp []salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.GetLocalVictimCaseForm2EnumsByCategory(enumTmp)
			entityIntervened = append(entityIntervened, eTmp...)
		}
	}
	if len(feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(feminicideRiskForm1, "FeminicideRiskForm1PlacesVisitRegularly", "json"), "feminicide_risk_form1_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

// En enums están todos los del caso actual, hay que extraer los necesarios para llenar el objeto
func loadFeminicideRiskEnumsMultiple(feminicideRiskForm1 *salvia_daos.FeminicideRiskForm1DTO, enums []salvia_daos.VictimCaseForm2EnumsDTO) {
	var err error

	for _, e := range enums {
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsById(&e); err == nil {
			switch e.VictimCaseForm2EnumsCategory {
			case "feminicide_risk_form1_financially_dependent_people":
				feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople = append(feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople, e)
			}
		}
	}

	for _, e := range enums {
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsById(&e); err == nil {
			switch e.VictimCaseForm2EnumsCategory {
			case "feminicide_risk_form1_places_visit_regularly":
				feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly = append(feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly, e)
			}
		}
	}

}

// Dado que ya están validados, entonces se relacionan los enums con contenido
func setFeminicideRiskEnumsMultiple(feminicideRiskForm1 *salvia_daos.FeminicideRiskForm1DTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	var err error
	for _, e := range feminicideRiskForm1.FeminicideRiskForm1FinanciallyDependentPeople {
		var rel salvia_daos.RelVictimCaseForm2EnumsFeminicideRiskForm1DTO
		rel.RelVictimCaseForm2EnumsFeminicideRiskForm1Form = *feminicideRiskForm1
		rel.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsFeminicideRiskForm1(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range feminicideRiskForm1.FeminicideRiskForm1PlacesVisitRegularly {
		var rel salvia_daos.RelVictimCaseForm2EnumsFeminicideRiskForm1DTO
		rel.RelVictimCaseForm2EnumsFeminicideRiskForm1Form = *feminicideRiskForm1
		rel.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsFeminicideRiskForm1(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}
	return nil
}
