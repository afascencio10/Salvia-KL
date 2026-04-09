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
	"bitsflow/salvia/service"
	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"
	security_daos "bitsflow/security/dao"
	"context"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// FollowUpSvc es inyectado desde main.go para generar el calendario automáticamente
// al crear un caso. Si es nil, la generación automática se omite silenciosamente.
var FollowUpSvc service.FollowUpV2Service

type VictimCaseRequest struct {
	VCase salvia_daos.VictimCaseDTO `json:"victimCase"`
}

type VictimCaseForm1Request struct {
	Form salvia_daos.VictimCaseForm1DTO `json:"form"`
}

// SetVictimCase crea un nuevo caso de víctima a partir de la entrada JSON.
// Realiza la validación de los datos, crea o actualiza usuarios y relaciones,
// y gestiona la transacción en la base de datos. Devuelve un código HTTP y
// un mensaje en formato JSON (éxito o error).
func SetVictimCase(dataInput string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de variables y obtención de conexión a la BD.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	// Se libera la conexión al finalizar la función.
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores durante la validación del DTO.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos para el caso, contacto y dueño.
	var vCaseRequest VictimCaseRequest = VictimCaseRequest{}
	var victimContact salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	// Variable para almacenar el DTO en formato mapa.
	var dtoMap map[string]interface{} = nil

	// Se parsea el JSON de entrada a un mapa.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.VictimCaseJSONName, common_config.Locale, collectedErrors)
	utils.JSONToStruct(dataInput, &vCaseRequest)

	// Se asignan valores por defecto al objeto VictimCase.
	salvia_daos.SetVictimCaseDefaults(&vCaseRequest.VCase, common_dao.SQL_INSERT, s)
	salvia_daos.SetVictimCaseForm2Defaults(&vCaseRequest.VCase.VictimCaseForm2, common_dao.SQL_INSERT)

	// Verificamos que los datos tengan la estructura esperada.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar (obligatorios y opcionales).
		var checkFields map[string]bool = map[string]bool{
			"VictimCaseNames":     true,
			"VictimCaseLastNames": true,
			"VictimCaseDocType":   true,
			"VictimCaseDocNumber": true,
		}
		// Validación de los campos del JSON contra la definición del DTO.
		utils.ValidateJSONInput(&vCaseRequest.VCase, dtoMap, salvia_daos.VictimCaseJSONName, salvia_daos.VictimCaseFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"VictimCaseForm2IdentityName":                          false,
			"VictimCaseForm2VictimPhone":                           true,
			"VictimCaseForm2FactsDescription":                      true,
			"VictimCaseForm2FactsDate":                             true,
			"VictimCaseForm2FactsStartTime":                        true,
			"VictimCaseForm2FactsTownCode":                         true,
			"VictimCaseForm2FactsZone":                             true,
			"VictimCaseForm2FactsAddress":                          true,
			"VictimCaseForm2ScenarioViolence":                      true,
			"VictimCaseForm2ReportedPreviously":                    true,
			"VictimCaseForm2RecurrenceAggression":                  true,
			"VictimCaseForm2NumAgressors":                          true,
			"VictimCaseForm2ProximityPrincipalAggressor":           true,
			"VictimCaseForm2RelationshipWithPresumedAggressor":     true,
			"VictimCaseForm2EconomicallyDependent":                 true,
			"VictimCaseForm2AggressorGenderIdentity":               false,
			"VictimCaseForm2AggressorNames":                        false,
			"VictimCaseForm2AggressorDocType":                      false,
			"VictimCaseForm2AggressorDocNumber":                    false,
			"VictimCaseForm2AggressorAddress":                      false,
			"VictimCaseForm2AggressorPhone":                        false,
			"VictimCaseForm2AggressorViolencePhysicalIncrease":     true,
			"VictimCaseForm2AggressorWeaponUsed":                   true,
			"VictimCaseForm2AggressorThreatKill":                   true,
			"VictimCaseForm2AggressorPursuesSpiesDestroys":         true,
			"VictimCaseForm2AggressorCapableOfKilling":             true,
			"VictimCaseForm2AggressorHasAccessToWeapons":           true,
			"VictimCaseForm2PartnerUnemployed":                     false,
			"VictimCaseForm2PartnerOtherDenunciations":             false,
			"VictimCaseForm2AggressorHasPenalBackground":           false,
			"VictimCaseForm2AggressorForcedSex":                    false,
			"VictimCaseForm2AggressorAttemptedStrangulation":       false,
			"VictimCaseForm2AggressorConsumesDrugs":                false,
			"VictimCaseForm2AggressorIsAlcoholic":                  false,
			"VictimCaseForm2PartnerControls":                       false,
			"VictimCaseForm2AggressorHadHitInVulnerability":        false,
			"VictimCaseForm2PartnerThreatenedSuicide":              false,
			"VictimCaseForm2PartnerThreatenedDamageMembers":        false,
			"VictimCaseForm2ThoughtsOfSelfHarm":                    false,
			"VictimCaseForm2AggressorLimitsContactSupportNetworks": false,
			"VictimCaseForm2StillLivesWithAggressor":               false,
			"VictimCaseForm2AggressorViolentlyJealous":             false,
			"VictimCaseForm2AggressorUnemployed":                   false,
			"VictimCaseForm2AggressorHasPenalBackground2":          false,
			"VictimCaseForm2AggressorSexuallyHarassment":           false,
			"VictimCaseForm2AggressorUseDrugs":                     false,
			"VictimCaseForm2AggressorIsAlcoholic2":                 false,
			"VictimCaseForm2AggressorControls":                     false,
			"VictimCaseForm2AggressorThreatenedDamageMembers":      false,
			"VictimCaseForm2ThoughtsOfSelfHarm2":                   false,
			"VictimCaseForm2AggressorCommonSpaces":                 false,
			"VictimCaseForm2AggressorHierarchy":                    false,
			"VictimCaseForm2BirthDate":                             true,
			"VictimCaseForm2PhysicalMentalSensoryDifficulties":     true,
			"VictimCaseForm2Nationality":                           true,
			"VictimCaseForm2SpecifiedNationality":                  false,
			"VictimCaseForm2MigrationCondition":                    true,
			"VictimCaseForm2GenderIdentity":                        true,
			"VictimCaseForm2SexualOrientation":                     true,
			"VictimCaseForm2AssignedSexAtBirth":                    true,
			"VictimCaseForm2EthnicAffiliation":                     true,
			"VictimCaseForm2IndigenousPeople":                      false,
			"VictimCaseForm2CampesinoRecognition":                  true,
			"VictimCaseForm2MaritalStatus":                         true,
			"VictimCaseForm2LastEducationLevel":                    true,
			"VictimCaseForm2Occupation":                            true,
			"VictimCaseForm2IncomeGenerationMethod":                true,
			"VictimCaseForm2EmploymentRelationship":                false,
			"VictimCaseForm2ApproxStartAsp":                        false,
			"VictimCaseForm2HousingTenancyForm":                    true,
			"VictimCaseForm2HousingStratum":                        true,
			"VictimCaseForm2CurrentlyPregnant":                     true,
			"VictimCaseForm2ResidenceTown":                         true,
			"VictimCaseForm2ResidenceAddress":                      true,
			"VictimCaseForm2ResidenceZone":                         true,
			"VictimCaseForm2SupportContactNames":                   false,
			"VictimCaseForm2SupportContactPhone":                   false,
			"VictimCaseForm2SupportContactEmail":                   false,
			"VictimCaseForm2SupportContactKinship":                 false,
			"VictimCaseForm2SalivaManagementExplanation":           true,
			"VictimCaseForm2ActivitiesUnableToHear":                false,
			"VictimCaseForm2ActivitiesUnableToTalk":                false,
			"VictimCaseForm2ActivitiesUnableToSee":                 false,
			"VictimCaseForm2ActivitiesUnableToMove":                false,
			"VictimCaseForm2ActivitiesUnableToTake":                false,
			"VictimCaseForm2ActivitiesUnableToUnderstand":          false,
			"VictimCaseForm2ActivitiesUnableToEat":                 false,
			"VictimCaseForm2ActivitiesUnableToInteract":            false,
			"VictimCaseForm2ActivitiesUnableToDoEveryday":          false,
			"VictimCaseForm2PersonWithDisability":                  true,
			"VictimCaseForm2RequireLanguageInterpreter":            true,
			"VictimCaseForm2WorkplaceSectorOccurrence":             false,
			"VictimCaseForm2ViolenceMotivatedByGender":             true,
			"VictimCaseForm2AttentionWasAppropriate":               true,
			"VictimCaseForm2AggressorOccupation":                   true,
		}

		utils.ValidateJSONInput(&vCaseRequest.VCase.VictimCaseForm2, dtoMap, salvia_daos.VictimCaseJSONName+"."+salvia_daos.VictimCaseForm2JSONName, salvia_daos.VictimCaseForm2FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

	default:
		// Si el JSON no tiene la estructura esperada, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	//Ahora se verifican los enums únicos
	getAndVerifyVictimCaseEnums(&vCaseRequest.VCase.VictimCaseForm2, collectedErrors, true)

	//Ahora se verifican los enums múltiples que se asociarán al form2
	getAndVerifyVictimCaseEnumsMultiple(&vCaseRequest.VCase.VictimCaseForm2, collectedErrors)

	//Ahora algunos campos opcionales
	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsCode == "pr" {
		if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ApproxStartAsp.IsZero() {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2ApproxStartAsp", "json"), "common_validation_field_date_error", "common_global_error", common_config.Locale)
		}
	}

	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RequireLanguageInterpreter.VictimCaseForm2EnumsCode == "y" {
		var l = int64(len(vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2LanguageAssistance))
		if l == 0 {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2LanguageAssistance", "json"), "common_validation_field_required_error", "common_global_error", common_config.Locale)
		} else if l < salvia_daos.VictimCaseForm2FieldDefinitions["VictimCaseForm2LanguageAssistance"].MinSize {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2LanguageAssistance", "json"), "common_validation_field_string_min_size_error", "common_global_error", common_config.Locale)
		} else if l >= salvia_daos.VictimCaseForm2FieldDefinitions["VictimCaseForm2LanguageAssistance"].MaxSize {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2LanguageAssistance", "json"), "common_validation_field_string_max_size_error", "common_global_error", common_config.Locale)
		}
	}

	//Verificamos que ninguna fecha supere la actual

	today, _ := time.Parse(common_config.DateTime.DATE_FORMAT, time.Now().Format(common_config.DateTime.DATE_FORMAT))

	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsDate.After(today) {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2FactsDate", "json"), "common_validation_field_date_future_error", "common_global_error", common_config.Locale)
	}

	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2BirthDate.After(today) {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2BirthDate", "json"), "common_validation_field_date_future_error", "common_global_error", common_config.Locale)
	}

	// Si se detectaron errores en la validación, se retorna BadRequest con los errores.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Calculamos el nivel de riesgo
	riskScore, riskLevel := getRiskScore(vCaseRequest.VCase.VictimCaseForm2)

	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RiskScore != riskScore {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2RiskScore", "json"), "victim_case_form2_risk_score_client_error", "common_global_error", common_config.Locale)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Si la comparación anterior coincide entonces sólo se actualiza el nivel de riesgo
	vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RiskLevel = riskLevel

	// Si existe un contacto definido, se busca en la base de datos.
	var by common_controllers.By = common_controllers.By{}
	if vCaseRequest.VCase.VictimCaseVictimContact.VictimContactICode != "" {
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"VictimContactICode"},
			AttrsValue: []interface{}{vCaseRequest.VCase.VictimCaseVictimContact.VictimContactICode},
		}

		victimContact.VictimContactICode = vCaseRequest.VCase.VictimCaseVictimContact.VictimContactICode

		// Se intenta obtener el contacto de la víctima.
		err = salvia_daos.GetVictimContact(by, &victimContact, connData, &dbClientConfig, &dbServerConfig)
		if err == nil {
			vCaseRequest.VCase.VictimCaseVictimContact.VictimContactId = victimContact.VictimContactId
		}
	}

	// Se obtiene el perfil de usuario basado en el documento (tipo y número).
	var profile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}
	var user security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
	var role security_daos.RoleDTO = security_daos.RoleDTO{}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
		AttrsValue: []interface{}{vCaseRequest.VCase.VictimCaseDocType, vCaseRequest.VCase.VictimCaseDocNumber},
	}

	err = security_daos.GetGeneralUserProfile(by, &profile, connData, &dbClientConfig, &dbServerConfig)
	if err == nil {
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserGeneralUserProfile"},
			AttrsValue: []interface{}{profile.GeneralUserProfileId},
		}

		_ = security_daos.GetGeneralUser(by, &user, connData, &dbClientConfig, &dbServerConfig)
	}
	// Si el usuario no existe, se procede a crearlo.
	// Se inicia la transacción.
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se genera una contraseña aleatoria para el usuario.
	var newPassword string = getRandomPassword(4)
	if user.GeneralUserId == 0 {
		// Configuración de valores por defecto para crear el perfil y el usuario.
		security_daos.SetGeneralUserProfileDefaults(&profile, common_dao.SQL_INSERT)
		security_daos.SetGeneralUserDefaults(&user, common_dao.SQL_INSERT)
		profile.GeneralUserProfileDocNumber = vCaseRequest.VCase.VictimCaseDocNumber
		profile.GeneralUserProfileDocType = vCaseRequest.VCase.VictimCaseDocType
		profile.GeneralUserProfileGender = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2GenderIdentity.VictimCaseForm2EnumsCode
		profile.GeneralUserProfileNames = vCaseRequest.VCase.VictimCaseNames
		profile.GeneralUserProfileLastNames = vCaseRequest.VCase.VictimCaseLastNames
		profile.GeneralUserProfileNick = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2IdentityName

		user.GeneralUserLanguage = s.Lang
		user.GeneralUserStatus = "e"
		user.GeneralUserLogin = getRandomUserName(6)
		// Se encripta la contraseña generada.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
		if err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		} else {
			user.GeneralUserPassword = string(hashedPassword)
		}

		// Se crea el perfil de usuario.
		if err = security_daos.SetGeneralUserProfile(&profile, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		// Se asigna el perfil al usuario.
		user.GeneralUserGeneralUserProfile = profile

		// Se crea el usuario.
		if err = security_daos.SetGeneralUser(&user, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		// Se obtiene el rol 'us' para el usuario.
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"RoleCode"},
			AttrsValue: []interface{}{"us"},
		}

		if err = security_daos.GetRole(by, &role, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
		var rel security_daos.RelRoleGeneralUserDTO
		security_daos.SetRelRoleGeneralUserDefaults(&rel, common_dao.SQL_INSERT)
		rel.RelRoleGeneralUserGeneralUser = user
		rel.RelRoleGeneralUserRole = role
		// Se establece la relación entre el usuario y el rol.
		if err = security_daos.SetRelRoleGeneralUser(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	} else {
		// Si el usuario ya existe y posee rol de 'us', se resetea la contraseña.
		if hasRole("us", user.GeneralUserRoles) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
			if err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, err.Error()
			} else {
				user.GeneralUserPassword = string(hashedPassword)
			}
			if err = security_daos.UpdateGeneralUserByICode(&user, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, err.Error()
			}
		} else {
			newPassword = salvia_config.Locale[s.Lang]["victim_case_set_user_same_password"]
		}
	}

	// Se asigna la contraseña nueva al usuario.
	user.GeneralUserPassword = newPassword
	vCaseRequest.VCase.VictimCaseGeneralUser = user.GeneralUserICode
	vCaseRequest.VCase.VictimCaseNewUser = user

	// Se determina el dueño del caso (para usuarios con rol "op").
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

	// Si el usuario en sesión tiene rol "op", se asigna la aprobación del caso.
	if s.CurrentRole == "op" {
		vCaseRequest.VCase.VictimCaseApprovedBy = owner
	}

	//Si no se eligió lugar de atención se utiliza erl lugar de vivienda
	if vCaseRequest.VCase.VictimCaseTownCode == "" {
		vCaseRequest.VCase.VictimCaseTownCode = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ResidenceTownCode
	}

	// Se crea el caso en la BD.
	if err = salvia_daos.SetVictimCase(&vCaseRequest.VCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2VictimCase = vCaseRequest.VCase

	if err = salvia_daos.SetVictimCaseForm2(&vCaseRequest.VCase.VictimCaseForm2, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Ahora se insertan todos los campos múltiples
	//Ahora se verifican los enums múltiples que se asociarán al form2
	if err = setVictimCaseEnumsMultiple(&vCaseRequest.VCase.VictimCaseForm2, connData, dbClientConfig, dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se procesan las ramas (entity branches) asociadas al caso.
	var code int
	for mcode, sectors := range vCaseRequest.VCase.VictimCaseEntityBranches {
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
					moment.MomentVictimCase = vCaseRequest.VCase

					if err = salvia_daos.SetMoment(&moment, connData, &dbClientConfig, &dbServerConfig); err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusInternalServerError, err.Error()
					}
				}
			}
		}
	}

	// Se asigna la autoría del caso al usuario en sesión.
	var relCaseOwnerVictimCase salvia_daos.RelCaseOwnerVictimCaseDTO = salvia_daos.RelCaseOwnerVictimCaseDTO{}
	salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&relCaseOwnerVictimCase, common_dao.SQL_INSERT)
	relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner = owner.CaseOwnerId
	relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase = vCaseRequest.VCase.VictimCaseId

	if err = salvia_daos.SetRelCaseOwnerVictimCase(&relCaseOwnerVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualiza el contador de casos del dueño.
	owner.CaseOwnerNumCases = owner.CaseOwnerNumCases + 1
	if err = salvia_daos.UpdateCaseOwner(&owner, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Si el usuario tiene rol "et", se busca un usuario "op" para compartir la autoría.
	if s.CurrentRole == "et" {
		var opOwners []salvia_daos.CaseOwnerDTO = []salvia_daos.CaseOwnerDTO{}
		opOwners, err = salvia_daos.GetCaseOwnersByActiveUsersAndRoleCode("op", connData, &dbClientConfig, &dbServerConfig)
		if err != nil || len(opOwners) == 0 {
			return http.StatusInternalServerError, err.Error()
		}
		relCaseOwnerVictimCase = salvia_daos.RelCaseOwnerVictimCaseDTO{}
		salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&relCaseOwnerVictimCase, common_dao.SQL_INSERT)
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner = opOwners[0].CaseOwnerId
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase = vCaseRequest.VCase.VictimCaseId
		if err = salvia_daos.SetRelCaseOwnerVictimCase(&relCaseOwnerVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
	if err = salvia_daos.UpdateVictimCaseOwnersAndRolesByVictimCaseId(vCaseRequest.VCase.VictimCaseId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se realiza el commit de la transacción y se retorna el resultado.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// ── HU-027: Generar calendario de seguimientos automáticamente ───────────
	// Se ejecuta DESPUÉS del commit para no afectar la transacción del caso.
	// Si falla, se loggea pero NO se deshace la creación del caso.
	if FollowUpSvc != nil {
		go func() {
			calendarInput := service.GenerateCalendarInput{
				RiskLevel: int(riskLevel), // 1=Bajo, 2=Moderado, 3=Alto, 4=Extremo
				AgentID:   s.UserICode,    // Operador que creó el caso
				Team:      "",             // Se asignará "SIN_EQUIPO" por defecto en el servicio
			}
			_, calErr := FollowUpSvc.GenerateOrRecalculate(context.Background(), vCaseRequest.VCase.VictimCaseICode, calendarInput)
			if calErr != nil {
				log.Printf("[WARN] HU-027: Error generando calendario para caso %s: %v", vCaseRequest.VCase.VictimCaseICode, calErr)
			} else {
				log.Printf("[INFO] HU-027: Calendario generado para caso %s (risk_level=%d)", vCaseRequest.VCase.VictimCaseICode, riskLevel)
			}
		}()
	}
	// ─────────────────────────────────────────────────────────────────────────

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCaseRequest.VCase.VictimCaseNewUser)
}

// UpdateVictimCase actualiza los datos de un caso de víctima existente.
// Recibe la entrada JSON, el identificador del caso (vCaseICode) y realiza la actualización
// validando el DTO, actualizando datos y gestionando la transacción.
func UpdateVictimCaseForm1(dataInput string, vCaseICode string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicializa variables y establece la conexión.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores de validación.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos: caso, dueño y el caso a actualizar.
	var vCaseRequest VictimCaseRequest = VictimCaseRequest{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	var vCaseToUpdate salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var vCaseFormToUpdate salvia_daos.VictimCaseForm1DTO = salvia_daos.VictimCaseForm1DTO{}
	var code int

	// Slices para gestionar los momentos asociados al caso.
	var momentsToDelete []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}
	var moments []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}

	// Variable para almacenar el DTO en forma de mapa.
	var dtoMap map[string]interface{} = nil

	// Se obtiene el mapa DTO a partir de la entrada JSON.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.VictimCaseJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &vCaseRequest)

	// Se establecen valores por defecto para el caso (modo actualización).
	salvia_daos.SetVictimCaseDefaults(&vCaseRequest.VCase, common_dao.SQL_UPDATE, s)
	salvia_daos.SetVictimCaseForm1Defaults(&vCaseRequest.VCase.VictimCaseForm1, common_dao.SQL_INSERT)

	// Se verifica la estructura del DTO.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar.
		var checkFields map[string]bool = map[string]bool{
			"VictimCaseNames":     true,
			"VictimCaseLastNames": true,
			"VictimCaseDocType":   true,
			"VictimCaseDocNumber": true,
		}
		// Validación de los datos del caso.
		utils.ValidateJSONInput(&vCaseRequest.VCase, dtoMap, salvia_daos.VictimCaseJSONName, salvia_daos.VictimCaseFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"VictimCaseForm1Nick":                              true,
			"VictimCaseForm1BirthDate":                         true,
			"VictimCaseForm1ViolenceTownCode":                  true,
			"VictimCaseForm1Address":                           true,
			"VictimCaseForm1LivingLatitude":                    false,
			"VictimCaseForm1LivingLongitude":                   false,
			"VictimCaseForm1Phone":                             true,
			"VictimCaseForm1Email":                             true,
			"VictimCaseForm1GenderIdentity":                    true,
			"VictimCaseForm1SexualOrientation":                 true,
			"VictimCaseForm1Origin":                            true,
			"VictimCaseForm1Occupation":                        true,
			"VictimCaseForm1OccupationOther":                   false,
			"VictimCaseForm1VictimEthnicGroup":                 true,
			"VictimCaseForm1VictimEthnicGroupOther":            false,
			"VictimCaseForm1VictimContactNames":                true,
			"VictimCaseForm1VictimContactPhone":                true,
			"VictimCaseForm1VictimContactKinship":              true,
			"VictimCaseForm1VictimNumChildren":                 false,
			"VictimCaseForm1VictimMaritalStatus":               true,
			"VictimCaseForm1VictimMaritalStatusOther":          false,
			"VictimCaseForm1VictimChildrenAge":                 false,
			"VictimCaseForm1VictimDisability":                  true,
			"VictimCaseForm1VictimSpecialSupport":              true,
			"VictimCaseForm1FactsOccurrence":                   true,
			"VictimCaseForm1FactsStartTime":                    true,
			"VictimCaseForm1FactsEndTime":                      true,
			"VictimCaseForm1FactsWeekday":                      true,
			"VictimCaseForm1FactsDate":                         true,
			"VictimCaseForm1FactsDescription":                  true,
			"VictimCaseForm1VictimViolenceExperienced":         true,
			"VictimCaseForm1VictimViolenceExperiencedOther":    false,
			"VictimCaseForm1VictimViolenceScope":               true,
			"VictimCaseForm1VictimFemicideRisk":                true,
			"VictimCaseForm1VictimAggressor":                   true,
			"VictimCaseForm1VictimRelationshipWithAggressor":   true,
			"VictimCaseForm1VictimAggressorName":               false,
			"VictimCaseForm1VictimAggressorDocType":            false,
			"VictimCaseForm1VictimAggressorDocNumber":          false,
			"VictimCaseForm1VictimAggressorAddress":            false,
			"VictimCaseForm1VictimAggressorPhone":              false,
			"VictimCaseForm1Age":                               true,
			"VictimCaseForm1VictimNationality":                 true,
			"VictimCaseForm1VictimNationalityOther":            false,
			"VictimCaseForm1VictimForeignerImmigrationStatus":  false,
			"VictimCaseForm1VictimGender":                      true,
			"VictimCaseForm1VictimGenderIdentityOther":         false,
			"VictimCaseForm1VictimSexualOrientationOther":      false,
			"VictimCaseForm1VictimDependents":                  true,
			"VictimCaseForm1VictimPreviouslyReportedSituation": true,
			"VictimCaseForm1VictimIfPreviouslyReported":        false,
			"VictimCaseForm1VictimIfAfro":                      false,
			"VictimCaseForm1VictimIfIndigenous":                false,
			"VictimCaseForm1VictimIfIndigenousTongue":          false,
			"VictimCaseForm1VictimIfPeasant":                   true,
			"VictimCaseForm1VictimIfArmedConflict":             true,
			"VictimCaseForm1VictimViolenceScene":               true,
			"VictimCaseForm1PhysicalViolenceIncreased":         true,
			"VictimCaseForm1SeparatedFromPartnerLastYear":      true,
			"VictimCaseForm1ThreatenedWithWeapon":              true,
			"VictimCaseForm1ThreatenedToKillOrHarmChildren":    true,
			"VictimCaseForm1JealousAndViolent":                 true,
			"VictimCaseForm1BelievesCapableOfKilling":          true,
		}

		utils.ValidateJSONInput(&vCaseRequest.VCase.VictimCaseForm1, dtoMap, salvia_daos.VictimCaseJSONName+"."+salvia_daos.VictimCaseForm1JSONName, salvia_daos.VictimCaseForm1FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

	default:
		// Si el JSON no cumple la estructura, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si existen errores, se retorna un BadRequest.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se obtiene el caso a actualizar a partir del identificador.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{vCaseICode},
	}

	if err = salvia_daos.GetVictimCase(by, &vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseForm1VictimCase"},
		AttrsValue: []interface{}{vCaseToUpdate.VictimCaseId},
	}

	if err = salvia_daos.GetVictimCaseForm1(by, &vCaseFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Verifica si se produjo un cambio en el documento de identidad.
	if vCaseRequest.VCase.VictimCaseDocType != vCaseToUpdate.VictimCaseDocType || vCaseRequest.VCase.VictimCaseDocNumber != vCaseToUpdate.VictimCaseDocNumber {
		var oldProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}
		var newProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
			AttrsValue: []interface{}{vCaseRequest.VCase.VictimCaseDocType, vCaseRequest.VCase.VictimCaseDocNumber},
		}

		err = security_daos.GetGeneralUserProfile(by, &newProfile, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// No existe un usuario con el nuevo documento; se actualiza el perfil del usuario.
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
				AttrsValue: []interface{}{vCaseToUpdate.VictimCaseDocType, vCaseToUpdate.VictimCaseDocNumber},
			}
			err = security_daos.GetGeneralUserProfile(by, &oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_victim_case_find_profile_fail"]
			}
			// Se actualiza el perfil con el nuevo documento.
			security_daos.SetGeneralUserProfileDefaults(&oldProfile, common_dao.SQL_UPDATE)
			oldProfile.GeneralUserProfileDocType = vCaseRequest.VCase.VictimCaseDocType
			oldProfile.GeneralUserProfileDocNumber = vCaseRequest.VCase.VictimCaseDocNumber

			err = security_daos.UpdateGeneralUserProfile(&oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_victim_case_user_not_found"]
			}
		} else {
			// Si ya existe un perfil con el nuevo documento, se retorna un error.
			utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, utils.GetTag(&vCaseRequest.VCase, "VictimCaseDocNumber", "json"),
				"update_victim_case_find_profile_already_exist", "salvia_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	// Se inicia la transacción para la actualización.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualizan los campos del caso con los nuevos valores.
	vCaseToUpdate.VictimCaseNames = vCaseRequest.VCase.VictimCaseNames
	vCaseToUpdate.VictimCaseLastNames = vCaseRequest.VCase.VictimCaseLastNames
	vCaseToUpdate.VictimCaseDocType = vCaseRequest.VCase.VictimCaseDocType
	vCaseToUpdate.VictimCaseDocNumber = vCaseRequest.VCase.VictimCaseDocNumber
	vCaseToUpdate.VictimCaseTownCode = vCaseRequest.VCase.VictimCaseTownCode

	vCaseFormToUpdate.VictimCaseForm1Nick = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Nick
	vCaseFormToUpdate.VictimCaseForm1BirthDate = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1BirthDate
	vCaseFormToUpdate.VictimCaseForm1Address = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Address
	vCaseFormToUpdate.VictimCaseForm1LivingLatitude = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1LivingLatitude
	vCaseFormToUpdate.VictimCaseForm1LivingLongitude = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1LivingLongitude
	vCaseFormToUpdate.VictimCaseForm1Phone = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Phone
	vCaseFormToUpdate.VictimCaseForm1Email = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Email
	vCaseFormToUpdate.VictimCaseForm1GenderIdentity = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1GenderIdentity
	vCaseFormToUpdate.VictimCaseForm1SexualOrientation = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1SexualOrientation
	vCaseFormToUpdate.VictimCaseForm1Origin = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Origin
	vCaseFormToUpdate.VictimCaseForm1Occupation = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Occupation
	vCaseFormToUpdate.VictimCaseForm1OccupationOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1OccupationOther
	vCaseFormToUpdate.VictimCaseForm1VictimEthnicGroup = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimEthnicGroup
	vCaseFormToUpdate.VictimCaseForm1VictimEthnicGroupOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimEthnicGroupOther
	vCaseFormToUpdate.VictimCaseForm1VictimContactNames = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimContactNames
	vCaseFormToUpdate.VictimCaseForm1VictimContactPhone = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimContactPhone
	vCaseFormToUpdate.VictimCaseForm1VictimContactKinship = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimContactKinship
	vCaseFormToUpdate.VictimCaseForm1VictimNumChildren = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimNumChildren
	vCaseFormToUpdate.VictimCaseForm1VictimMaritalStatus = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimMaritalStatus
	vCaseFormToUpdate.VictimCaseForm1VictimMaritalStatusOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimMaritalStatusOther
	vCaseFormToUpdate.VictimCaseForm1VictimChildrenAge = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimChildrenAge
	vCaseFormToUpdate.VictimCaseForm1VictimDisability = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimDisability
	vCaseFormToUpdate.VictimCaseForm1VictimSpecialSupport = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimSpecialSupport
	vCaseFormToUpdate.VictimCaseForm1FactsOccurrence = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1FactsOccurrence
	vCaseFormToUpdate.VictimCaseForm1FactsStartTime = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1FactsStartTime
	vCaseFormToUpdate.VictimCaseForm1FactsEndTime = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1FactsEndTime
	vCaseFormToUpdate.VictimCaseForm1FactsWeekday = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1FactsWeekday
	vCaseFormToUpdate.VictimCaseForm1FactsDate = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1FactsDate
	vCaseFormToUpdate.VictimCaseForm1FactsDescription = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1FactsDescription
	vCaseFormToUpdate.VictimCaseForm1VictimViolenceExperienced = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimViolenceExperienced
	vCaseFormToUpdate.VictimCaseForm1VictimViolenceExperiencedOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimViolenceExperiencedOther
	vCaseFormToUpdate.VictimCaseForm1VictimViolenceScope = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimViolenceScope
	vCaseFormToUpdate.VictimCaseForm1ViolenceTownCode = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1ViolenceTown.TownCode
	vCaseFormToUpdate.VictimCaseForm1VictimFemicideRisk = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimFemicideRisk
	vCaseFormToUpdate.VictimCaseForm1VictimAggressor = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimAggressor
	vCaseFormToUpdate.VictimCaseForm1VictimRelationshipWithAggressor = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimRelationshipWithAggressor
	vCaseFormToUpdate.VictimCaseForm1VictimAggressorName = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimAggressorName
	vCaseFormToUpdate.VictimCaseForm1VictimAggressorDocType = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimAggressorDocType
	vCaseFormToUpdate.VictimCaseForm1VictimAggressorDocNumber = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimAggressorDocNumber
	vCaseFormToUpdate.VictimCaseForm1VictimAggressorAddress = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimAggressorAddress
	vCaseFormToUpdate.VictimCaseForm1VictimAggressorPhone = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimAggressorPhone

	vCaseFormToUpdate.VictimCaseForm1Age = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1Age
	vCaseFormToUpdate.VictimCaseForm1VictimNationality = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimNationality
	vCaseFormToUpdate.VictimCaseForm1VictimNationalityOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimNationalityOther
	vCaseFormToUpdate.VictimCaseForm1VictimForeignerImmigrationStatus = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimForeignerImmigrationStatus
	vCaseFormToUpdate.VictimCaseForm1VictimGender = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimGender
	vCaseFormToUpdate.VictimCaseForm1VictimGenderIdentityOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimGenderIdentityOther
	vCaseFormToUpdate.VictimCaseForm1VictimSexualOrientationOther = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimSexualOrientationOther
	vCaseFormToUpdate.VictimCaseForm1VictimDependents = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimDependents
	vCaseFormToUpdate.VictimCaseForm1VictimPreviouslyReportedSituation = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimPreviouslyReportedSituation
	vCaseFormToUpdate.VictimCaseForm1VictimIfPreviouslyReported = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimIfPreviouslyReported

	vCaseFormToUpdate.VictimCaseForm1VictimIfAfro = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimIfAfro
	vCaseFormToUpdate.VictimCaseForm1VictimIfIndigenous = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimIfIndigenous
	vCaseFormToUpdate.VictimCaseForm1VictimIfIndigenousTongue = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimIfIndigenousTongue
	vCaseFormToUpdate.VictimCaseForm1VictimIfPeasant = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimIfPeasant
	vCaseFormToUpdate.VictimCaseForm1VictimIfArmedConflict = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimIfArmedConflict
	vCaseFormToUpdate.VictimCaseForm1VictimViolenceScene = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1VictimViolenceScene

	vCaseFormToUpdate.VictimCaseForm1PhysicalViolenceIncreased = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1PhysicalViolenceIncreased
	vCaseFormToUpdate.VictimCaseForm1SeparatedFromPartnerLastYear = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1SeparatedFromPartnerLastYear
	vCaseFormToUpdate.VictimCaseForm1ThreatenedWithWeapon = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1ThreatenedWithWeapon
	vCaseFormToUpdate.VictimCaseForm1ThreatenedToKillOrHarmChildren = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1ThreatenedToKillOrHarmChildren
	vCaseFormToUpdate.VictimCaseForm1JealousAndViolent = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1JealousAndViolent
	vCaseFormToUpdate.VictimCaseForm1BelievesCapableOfKilling = vCaseRequest.VCase.VictimCaseForm1.VictimCaseForm1BelievesCapableOfKilling

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
	if err = salvia_daos.UpdateVictimCase(&vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if err = salvia_daos.UpdateVictimCaseForm1(&vCaseFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen los momentos (acciones o eventos) asociados al caso.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentVictimCase"},
		AttrsValue: []interface{}{vCaseToUpdate.VictimCaseId},
	}

	momentsToDelete, err = salvia_daos.GetMoments(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se procesan las ramas para determinar cuáles momentos eliminar y cuáles crear.
	for mcode, sectors := range vCaseRequest.VCase.VictimCaseEntityBranches {
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
					moment.MomentVictimCase = vCaseToUpdate

					/* Se agregan los nuevos momentos al slice para crear. */
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
	var relCaseOwners []salvia_daos.RelCaseOwnerVictimCaseDTO = []salvia_daos.RelCaseOwnerVictimCaseDTO{}
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelCaseOwnerVictimCase_CaseOwner", "RelCaseOwnerVictimCase_VictimCase"},
		AttrsValue: []interface{}{owner.CaseOwnerId, vCaseToUpdate.VictimCaseId},
	}
	relCaseOwners, err = salvia_daos.GetRelCaseOwnerVictimCases(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Si no existe la relación, se crea y se actualiza el contador del dueño.
	if len(relCaseOwners) == 0 {
		var relCaseOwnerVictimCase salvia_daos.RelCaseOwnerVictimCaseDTO = salvia_daos.RelCaseOwnerVictimCaseDTO{}
		salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&relCaseOwnerVictimCase, common_dao.SQL_INSERT)
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner = owner.CaseOwnerId
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase = vCaseToUpdate.VictimCaseId

		if err = salvia_daos.SetRelCaseOwnerVictimCase(&relCaseOwnerVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
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
	if err = salvia_daos.UpdateVictimCaseOwnersAndRolesByVictimCaseId(vCaseToUpdate.VictimCaseId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se finaliza la transacción.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCaseRequest.VCase.VictimCaseNewUser)
}

func UpdateVictimCaseForm2(dataInput string, vCaseICode string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicializa variables y establece la conexión.
	var err error = nil
	var connData *db.ConnData = &db.ConnData{}
	defer db.ReleaseConnection(connData)

	// Mapa para almacenar errores de validación.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos: caso, dueño y el caso a actualizar.
	var vCaseRequest VictimCaseRequest = VictimCaseRequest{}
	var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	var vCaseToUpdate salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var vCaseFormToUpdate salvia_daos.VictimCaseForm2DTO = salvia_daos.VictimCaseForm2DTO{}
	var code int
	var enums []salvia_daos.VictimCaseForm2EnumsDTO

	// Slices para gestionar los momentos asociados al caso.
	var momentsToDelete []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}
	var moments []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}

	// Variable para almacenar el DTO en forma de mapa.
	var dtoMap map[string]interface{} = nil

	// Se obtiene el mapa DTO a partir de la entrada JSON.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.VictimCaseJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &vCaseRequest)

	// Se establecen valores por defecto para el caso (modo actualización).
	salvia_daos.SetVictimCaseDefaults(&vCaseRequest.VCase, common_dao.SQL_UPDATE, s)
	salvia_daos.SetVictimCaseForm2Defaults(&vCaseRequest.VCase.VictimCaseForm2, common_dao.SQL_INSERT)

	// Se verifica la estructura del DTO.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Definición de campos a validar.
		var checkFields map[string]bool = map[string]bool{
			"VictimCaseNames":     true,
			"VictimCaseLastNames": true,
			"VictimCaseDocType":   true,
			"VictimCaseDocNumber": true,
		}
		// Validación de los datos del caso.
		utils.ValidateJSONInput(&vCaseRequest.VCase, dtoMap, salvia_daos.VictimCaseJSONName, salvia_daos.VictimCaseFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"VictimCaseForm2IdentityName":                          false,
			"VictimCaseForm2VictimPhone":                           true,
			"VictimCaseForm2FactsDescription":                      true,
			"VictimCaseForm2FactsDate":                             true,
			"VictimCaseForm2FactsStartTime":                        true,
			"VictimCaseForm2FactsTownCode":                         true,
			"VictimCaseForm2FactsZone":                             true,
			"VictimCaseForm2FactsAddress":                          true,
			"VictimCaseForm2ScenarioViolence":                      true,
			"VictimCaseForm2ReportedPreviously":                    true,
			"VictimCaseForm2RecurrenceAggression":                  true,
			"VictimCaseForm2NumAgressors":                          true,
			"VictimCaseForm2ProximityPrincipalAggressor":           true,
			"VictimCaseForm2RelationshipWithPresumedAggressor":     true,
			"VictimCaseForm2EconomicallyDependent":                 true,
			"VictimCaseForm2AggressorGenderIdentity":               false,
			"VictimCaseForm2AggressorNames":                        false,
			"VictimCaseForm2AggressorDocType":                      false,
			"VictimCaseForm2AggressorDocNumber":                    false,
			"VictimCaseForm2AggressorAddress":                      false,
			"VictimCaseForm2AggressorPhone":                        false,
			"VictimCaseForm2AggressorViolencePhysicalIncrease":     true,
			"VictimCaseForm2AggressorWeaponUsed":                   true,
			"VictimCaseForm2AggressorThreatKill":                   true,
			"VictimCaseForm2AggressorPursuesSpiesDestroys":         true,
			"VictimCaseForm2AggressorCapableOfKilling":             true,
			"VictimCaseForm2AggressorHasAccessToWeapons":           true,
			"VictimCaseForm2PartnerUnemployed":                     false,
			"VictimCaseForm2PartnerOtherDenunciations":             false,
			"VictimCaseForm2AggressorHasPenalBackground":           false,
			"VictimCaseForm2AggressorForcedSex":                    false,
			"VictimCaseForm2AggressorAttemptedStrangulation":       false,
			"VictimCaseForm2AggressorConsumesDrugs":                false,
			"VictimCaseForm2AggressorIsAlcoholic":                  false,
			"VictimCaseForm2PartnerControls":                       false,
			"VictimCaseForm2AggressorHadHitInVulnerability":        false,
			"VictimCaseForm2PartnerThreatenedSuicide":              false,
			"VictimCaseForm2PartnerThreatenedDamageMembers":        false,
			"VictimCaseForm2ThoughtsOfSelfHarm":                    false,
			"VictimCaseForm2AggressorLimitsContactSupportNetworks": false,
			"VictimCaseForm2StillLivesWithAggressor":               false,
			"VictimCaseForm2AggressorViolentlyJealous":             false,
			"VictimCaseForm2AggressorUnemployed":                   false,
			"VictimCaseForm2AggressorHasPenalBackground2":          false,
			"VictimCaseForm2AggressorSexuallyHarassment":           false,
			"VictimCaseForm2AggressorUseDrugs":                     false,
			"VictimCaseForm2AggressorIsAlcoholic2":                 false,
			"VictimCaseForm2AggressorControls":                     false,
			"VictimCaseForm2AggressorThreatenedDamageMembers":      false,
			"VictimCaseForm2ThoughtsOfSelfHarm2":                   false,
			"VictimCaseForm2AggressorCommonSpaces":                 false,
			"VictimCaseForm2AggressorHierarchy":                    false,
			"VictimCaseForm2BirthDate":                             true,
			"VictimCaseForm2PhysicalMentalSensoryDifficulties":     true,
			"VictimCaseForm2Nationality":                           true,
			"VictimCaseForm2SpecifiedNationality":                  false,
			"VictimCaseForm2MigrationCondition":                    true,
			"VictimCaseForm2GenderIdentity":                        true,
			"VictimCaseForm2SexualOrientation":                     true,
			"VictimCaseForm2AssignedSexAtBirth":                    true,
			"VictimCaseForm2EthnicAffiliation":                     true,
			"VictimCaseForm2IndigenousPeople":                      false,
			"VictimCaseForm2CampesinoRecognition":                  true,
			"VictimCaseForm2MaritalStatus":                         true,
			"VictimCaseForm2LastEducationLevel":                    true,
			"VictimCaseForm2Occupation":                            true,
			"VictimCaseForm2IncomeGenerationMethod":                true,
			"VictimCaseForm2EmploymentRelationship":                false,
			"VictimCaseForm2ApproxStartAsp":                        false,
			"VictimCaseForm2HousingTenancyForm":                    true,
			"VictimCaseForm2HousingStratum":                        true,
			"VictimCaseForm2CurrentlyPregnant":                     true,
			"VictimCaseForm2ResidenceTown":                         true,
			"VictimCaseForm2ResidenceAddress":                      true,
			"VictimCaseForm2ResidenceZone":                         true,
			"VictimCaseForm2SupportContactNames":                   false,
			"VictimCaseForm2SupportContactPhone":                   false,
			"VictimCaseForm2SupportContactEmail":                   false,
			"VictimCaseForm2SupportContactKinship":                 false,
			"VictimCaseForm2LanguageAssistance":                    false,
			"VictimCaseForm2SalivaManagementExplanation":           true,
			"VictimCaseForm2ActivitiesUnableToHear":                false,
			"VictimCaseForm2ActivitiesUnableToTalk":                false,
			"VictimCaseForm2ActivitiesUnableToSee":                 false,
			"VictimCaseForm2ActivitiesUnableToMove":                false,
			"VictimCaseForm2ActivitiesUnableToTake":                false,
			"VictimCaseForm2ActivitiesUnableToUnderstand":          false,
			"VictimCaseForm2ActivitiesUnableToEat":                 false,
			"VictimCaseForm2ActivitiesUnableToInteract":            false,
			"VictimCaseForm2ActivitiesUnableToDoEveryday":          false,
			"VictimCaseForm2PersonWithDisability":                  true,
			"VictimCaseForm2RequireLanguageInterpreter":            true,
			"VictimCaseForm2WorkplaceSectorOccurrence":             false,
			"VictimCaseForm2ViolenceMotivatedByGender":             true,
			"VictimCaseForm2AttentionWasAppropriate":               true,
			"VictimCaseForm2AggressorOccupation":                   true,
		}

		utils.ValidateJSONInput(&vCaseRequest.VCase.VictimCaseForm2, dtoMap, salvia_daos.VictimCaseJSONName+"."+salvia_daos.VictimCaseForm2JSONName, salvia_daos.VictimCaseForm2FieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

	default:
		// Si el JSON no cumple la estructura, se establece un error global.
		utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	//Ahora se verifican los enums únicos
	getAndVerifyVictimCaseEnums(&vCaseRequest.VCase.VictimCaseForm2, collectedErrors, true)

	//Ahora se verifican los enums múltiples que se asociarán al form2
	getAndVerifyVictimCaseEnumsMultiple(&vCaseRequest.VCase.VictimCaseForm2, collectedErrors)

	//Ahora algunos campos opcionales
	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsCode == "pr" {
		if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ApproxStartAsp.IsZero() {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2ApproxStartAsp", "json"), "common_validation_field_date_error", "common_global_error", common_config.Locale)
		}
	}

	//Calculamos el nivel de riesgo
	riskScore, riskLevel := getRiskScore(vCaseRequest.VCase.VictimCaseForm2)

	if vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RiskScore != riskScore {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vCaseRequest.VCase.VictimCaseForm2, "VictimCaseForm2RiskScore", "json"), "victim_case_form2_risk_score_client_error", "common_global_error", common_config.Locale)
	}

	//Si la comparación anterior coincide entonces sólo se actualiza el nivel de riesgo
	vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RiskLevel = riskLevel

	// Si existen errores, se retorna un BadRequest.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se obtiene el caso a actualizar a partir del identificador.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{vCaseICode},
	}

	if err = salvia_daos.GetVictimCase(by, &vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseForm2VictimCase"},
		AttrsValue: []interface{}{vCaseToUpdate.VictimCaseId},
	}

	if err = salvia_daos.GetVictimCaseForm2(by, &vCaseFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	//Traemos todos los enums asociados a este formulario
	enums, err = salvia_daos.GetVictimCasesForm2EnumsByVictimcaseForm2Id(vCaseFormToUpdate.VictimCaseForm2Id, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	loadVictimCaseEnumsMultiple(&vCaseFormToUpdate, enums)

	// Verifica si se produjo un cambio en el documento de identidad.
	if vCaseRequest.VCase.VictimCaseDocType != vCaseToUpdate.VictimCaseDocType || vCaseRequest.VCase.VictimCaseDocNumber != vCaseToUpdate.VictimCaseDocNumber {
		var oldProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}
		var newProfile security_daos.GeneralUserProfileDTO = security_daos.GeneralUserProfileDTO{}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
			AttrsValue: []interface{}{vCaseRequest.VCase.VictimCaseDocType, vCaseRequest.VCase.VictimCaseDocNumber},
		}

		err = security_daos.GetGeneralUserProfile(by, &newProfile, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// No existe un usuario con el nuevo documento; se actualiza el perfil del usuario.
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
				AttrsValue: []interface{}{vCaseToUpdate.VictimCaseDocType, vCaseToUpdate.VictimCaseDocNumber},
			}
			err = security_daos.GetGeneralUserProfile(by, &oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_victim_case_find_profile_fail"]
			}
			// Se actualiza el perfil con el nuevo documento.
			security_daos.SetGeneralUserProfileDefaults(&oldProfile, common_dao.SQL_UPDATE)
			oldProfile.GeneralUserProfileDocType = vCaseRequest.VCase.VictimCaseDocType
			oldProfile.GeneralUserProfileDocNumber = vCaseRequest.VCase.VictimCaseDocNumber

			err = security_daos.UpdateGeneralUserProfile(&oldProfile, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["update_victim_case_user_not_found"]
			}
		} else {
			// Si ya existe un perfil con el nuevo documento, se retorna un error.
			utils.SetError(collectedErrors, salvia_daos.VictimCaseJSONName, utils.GetTag(&vCaseRequest.VCase, "VictimCaseDocNumber", "json"),
				"update_victim_case_find_profile_already_exist", "salvia_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	// Se inicia la transacción para la actualización.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualizan los campos del caso con los nuevos valores.
	vCaseToUpdate.VictimCaseNames = vCaseRequest.VCase.VictimCaseNames
	vCaseToUpdate.VictimCaseLastNames = vCaseRequest.VCase.VictimCaseLastNames
	vCaseToUpdate.VictimCaseDocType = vCaseRequest.VCase.VictimCaseDocType
	vCaseToUpdate.VictimCaseDocNumber = vCaseRequest.VCase.VictimCaseDocNumber
	vCaseToUpdate.VictimCaseTownCode = vCaseRequest.VCase.VictimCaseTownCode

	vCaseFormToUpdate.VictimCaseForm2IdentityName = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2IdentityName
	vCaseFormToUpdate.VictimCaseForm2VictimPhone = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2VictimPhone
	vCaseFormToUpdate.VictimCaseForm2FactsDescription = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsDescription
	vCaseFormToUpdate.VictimCaseForm2FactsDate = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsDate
	vCaseFormToUpdate.VictimCaseForm2FactsStartTime = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsStartTime
	vCaseFormToUpdate.VictimCaseForm2FactsTownCode = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsTownCode
	vCaseFormToUpdate.VictimCaseForm2FactsZone = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsZone
	vCaseFormToUpdate.VictimCaseForm2FactsAddress = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2FactsAddress
	vCaseFormToUpdate.VictimCaseForm2ScenarioViolence = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ScenarioViolence
	vCaseFormToUpdate.VictimCaseForm2ReportedPreviously = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ReportedPreviously
	vCaseFormToUpdate.VictimCaseForm2RecurrenceAggression = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RecurrenceAggression
	vCaseFormToUpdate.VictimCaseForm2NumAgressors = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2NumAgressors
	vCaseFormToUpdate.VictimCaseForm2ProximityPrincipalAggressor = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ProximityPrincipalAggressor
	vCaseFormToUpdate.VictimCaseForm2RelationshipWithPresumedAggressor = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor
	vCaseFormToUpdate.VictimCaseForm2EconomicallyDependent = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2EconomicallyDependent
	vCaseFormToUpdate.VictimCaseForm2AggressorGenderIdentity = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorGenderIdentity
	vCaseFormToUpdate.VictimCaseForm2AggressorNames = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorNames
	vCaseFormToUpdate.VictimCaseForm2AggressorDocType = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorDocType
	vCaseFormToUpdate.VictimCaseForm2AggressorDocNumber = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorDocNumber
	vCaseFormToUpdate.VictimCaseForm2AggressorAddress = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorAddress
	vCaseFormToUpdate.VictimCaseForm2AggressorPhone = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorPhone
	vCaseFormToUpdate.VictimCaseForm2AggressorViolencePhysicalIncrease = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease
	vCaseFormToUpdate.VictimCaseForm2AggressorWeaponUsed = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorWeaponUsed
	vCaseFormToUpdate.VictimCaseForm2AggressorThreatKill = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorThreatKill
	vCaseFormToUpdate.VictimCaseForm2AggressorPursuesSpiesDestroys = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys
	vCaseFormToUpdate.VictimCaseForm2AggressorCapableOfKilling = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorCapableOfKilling
	vCaseFormToUpdate.VictimCaseForm2AggressorHasAccessToWeapons = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons
	vCaseFormToUpdate.VictimCaseForm2PartnerUnemployed = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PartnerUnemployed
	vCaseFormToUpdate.VictimCaseForm2PartnerOtherDenunciations = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PartnerOtherDenunciations
	vCaseFormToUpdate.VictimCaseForm2AggressorHasPenalBackground = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorHasPenalBackground
	vCaseFormToUpdate.VictimCaseForm2AggressorForcedSex = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorForcedSex
	vCaseFormToUpdate.VictimCaseForm2AggressorAttemptedStrangulation = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorAttemptedStrangulation
	vCaseFormToUpdate.VictimCaseForm2AggressorConsumesDrugs = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorConsumesDrugs
	vCaseFormToUpdate.VictimCaseForm2AggressorIsAlcoholic = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorIsAlcoholic
	vCaseFormToUpdate.VictimCaseForm2PartnerControls = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PartnerControls
	vCaseFormToUpdate.VictimCaseForm2AggressorHadHitInVulnerability = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorHadHitInVulnerability
	vCaseFormToUpdate.VictimCaseForm2PartnerThreatenedSuicide = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PartnerThreatenedSuicide
	vCaseFormToUpdate.VictimCaseForm2PartnerThreatenedDamageMembers = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PartnerThreatenedDamageMembers
	vCaseFormToUpdate.VictimCaseForm2ThoughtsOfSelfHarm = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm
	vCaseFormToUpdate.VictimCaseForm2AggressorLimitsContactSupportNetworks = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorLimitsContactSupportNetworks
	vCaseFormToUpdate.VictimCaseForm2StillLivesWithAggressor = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2StillLivesWithAggressor
	vCaseFormToUpdate.VictimCaseForm2AggressorViolentlyJealous = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorViolentlyJealous
	vCaseFormToUpdate.VictimCaseForm2AggressorUnemployed = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorUnemployed
	vCaseFormToUpdate.VictimCaseForm2AggressorHasPenalBackground2 = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorHasPenalBackground2
	vCaseFormToUpdate.VictimCaseForm2AggressorSexuallyHarassment = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment
	vCaseFormToUpdate.VictimCaseForm2AggressorUseDrugs = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorUseDrugs
	vCaseFormToUpdate.VictimCaseForm2AggressorIsAlcoholic2 = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorIsAlcoholic2
	vCaseFormToUpdate.VictimCaseForm2AggressorControls = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorControls
	vCaseFormToUpdate.VictimCaseForm2AggressorThreatenedDamageMembers = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorThreatenedDamageMembers
	vCaseFormToUpdate.VictimCaseForm2ThoughtsOfSelfHarm2 = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm2
	vCaseFormToUpdate.VictimCaseForm2AggressorCommonSpaces = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorCommonSpaces
	vCaseFormToUpdate.VictimCaseForm2AggressorHierarchy = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorHierarchy
	vCaseFormToUpdate.VictimCaseForm2BirthDate = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2BirthDate
	vCaseFormToUpdate.VictimCaseForm2PhysicalMentalSensoryDifficulties = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PhysicalMentalSensoryDifficulties
	vCaseFormToUpdate.VictimCaseForm2Nationality = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2Nationality
	vCaseFormToUpdate.VictimCaseForm2SpecifiedNationality = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SpecifiedNationality
	vCaseFormToUpdate.VictimCaseForm2MigrationCondition = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2MigrationCondition
	vCaseFormToUpdate.VictimCaseForm2GenderIdentity = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2GenderIdentity
	vCaseFormToUpdate.VictimCaseForm2SexualOrientation = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SexualOrientation
	vCaseFormToUpdate.VictimCaseForm2AssignedSexAtBirth = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AssignedSexAtBirth
	vCaseFormToUpdate.VictimCaseForm2EthnicAffiliation = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2EthnicAffiliation
	vCaseFormToUpdate.VictimCaseForm2IndigenousPeople = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2IndigenousPeople
	vCaseFormToUpdate.VictimCaseForm2CampesinoRecognition = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2CampesinoRecognition
	vCaseFormToUpdate.VictimCaseForm2MaritalStatus = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2MaritalStatus
	vCaseFormToUpdate.VictimCaseForm2LastEducationLevel = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2LastEducationLevel
	vCaseFormToUpdate.VictimCaseForm2Occupation = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2Occupation
	vCaseFormToUpdate.VictimCaseForm2IncomeGenerationMethod = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2IncomeGenerationMethod
	vCaseFormToUpdate.VictimCaseForm2ApproxStartAsp = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ApproxStartAsp
	vCaseFormToUpdate.VictimCaseForm2HousingTenancyForm = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2HousingTenancyForm
	vCaseFormToUpdate.VictimCaseForm2HousingStratum = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2HousingStratum
	vCaseFormToUpdate.VictimCaseForm2CurrentlyPregnant = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2CurrentlyPregnant
	vCaseFormToUpdate.VictimCaseForm2ResidenceTownCode = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ResidenceTownCode
	vCaseFormToUpdate.VictimCaseForm2ResidenceAddress = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ResidenceAddress
	vCaseFormToUpdate.VictimCaseForm2ResidenceZone = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ResidenceZone
	vCaseFormToUpdate.VictimCaseForm2SupportContactNames = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SupportContactNames
	vCaseFormToUpdate.VictimCaseForm2SupportContactPhone = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SupportContactPhone
	vCaseFormToUpdate.VictimCaseForm2SupportContactEmail = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SupportContactEmail
	vCaseFormToUpdate.VictimCaseForm2SupportContactKinship = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SupportContactKinship
	vCaseFormToUpdate.VictimCaseForm2LanguageAssistance = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2LanguageAssistance
	vCaseFormToUpdate.VictimCaseForm2SalivaManagementExplanation = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2SalivaManagementExplanation
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToHear = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToHear
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToTalk = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToTalk
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToSee = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToSee
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToMove = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToMove
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToTake = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToTake
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToUnderstand = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToUnderstand
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToEat = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToEat
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToInteract = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToInteract
	vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToDoEveryday = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ActivitiesUnableToDoEveryday
	vCaseFormToUpdate.VictimCaseForm2PersonWithDisability = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2PersonWithDisability
	vCaseFormToUpdate.VictimCaseForm2RequireLanguageInterpreter = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2RequireLanguageInterpreter
	vCaseFormToUpdate.VictimCaseForm2WorkplaceSectorOccurrence = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2WorkplaceSectorOccurrence
	vCaseFormToUpdate.VictimCaseForm2ViolenceMotivatedByGender = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender
	vCaseFormToUpdate.VictimCaseForm2AttentionWasAppropriate = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AttentionWasAppropriate
	vCaseFormToUpdate.VictimCaseForm2AggressorOccupation = vCaseRequest.VCase.VictimCaseForm2.VictimCaseForm2AggressorOccupation

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
	if err = salvia_daos.UpdateVictimCase(&vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if err = salvia_daos.UpdateVictimCaseForm2(&vCaseFormToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen los momentos (acciones o eventos) asociados al caso.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentVictimCase"},
		AttrsValue: []interface{}{vCaseToUpdate.VictimCaseId},
	}

	momentsToDelete, err = salvia_daos.GetMoments(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se procesan las ramas para determinar cuáles momentos eliminar y cuáles crear.
	for mcode, sectors := range vCaseRequest.VCase.VictimCaseEntityBranches {
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
					moment.MomentVictimCase = vCaseToUpdate

					/* Se agregan los nuevos momentos al slice para crear. */
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
	var relCaseOwners []salvia_daos.RelCaseOwnerVictimCaseDTO = []salvia_daos.RelCaseOwnerVictimCaseDTO{}
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelCaseOwnerVictimCase_CaseOwner", "RelCaseOwnerVictimCase_VictimCase"},
		AttrsValue: []interface{}{owner.CaseOwnerId, vCaseToUpdate.VictimCaseId},
	}
	relCaseOwners, err = salvia_daos.GetRelCaseOwnerVictimCases(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Si no existe la relación, se crea y se actualiza el contador del dueño.
	if len(relCaseOwners) == 0 {
		var relCaseOwnerVictimCase salvia_daos.RelCaseOwnerVictimCaseDTO = salvia_daos.RelCaseOwnerVictimCaseDTO{}
		salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&relCaseOwnerVictimCase, common_dao.SQL_INSERT)
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner = owner.CaseOwnerId
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase = vCaseToUpdate.VictimCaseId

		if err = salvia_daos.SetRelCaseOwnerVictimCase(&relCaseOwnerVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		owner.CaseOwnerNumCases = owner.CaseOwnerNumCases + 1
		if err = salvia_daos.UpdateCaseOwner(&owner, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	//Se actualizan los enums mútiples
	if err = updateVictimCaseEnumsMultiple(vCaseFormToUpdate, vCaseRequest.VCase.VictimCaseForm2, connData, dbClientConfig, dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
	if err = salvia_daos.UpdateVictimCaseOwnersAndRolesByVictimCaseId(vCaseToUpdate.VictimCaseId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se finaliza la transacción.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCaseRequest.VCase.VictimCaseNewUser)
}

// GetVictimCaseByICode retorna un caso de víctima a partir de su identificador único (ICode).
// Además, recupera información geográfica relacionada (municipio, ciudad, departamento)
// y los momentos (logs) asociados al caso.
func GetVictimCaseByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.VictimCaseDTO) {
	var err error = nil

	// Se libera la conexión si ésta no es válida.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el caso.
	var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var vCaseForm1 salvia_daos.VictimCaseForm1DTO = salvia_daos.VictimCaseForm1DTO{}
	var vCaseForm2 salvia_daos.VictimCaseForm2DTO = salvia_daos.VictimCaseForm2DTO{}

	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, vCase
		}

	} else {
		var town security_daos.TownDTO
		var town2 security_daos.TownDTO
		var town3 security_daos.TownDTO
		var city2 security_daos.CityDTO
		var city3 security_daos.CityDTO
		var department security_daos.DepartmentDTO
		var department2 security_daos.DepartmentDTO
		var department3 security_daos.DepartmentDTO
		// Se consulta el caso por su identificador.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"VictimCaseICode"},
			AttrsValue: []interface{}{id},
		}

		err = salvia_daos.GetVictimCase(by, &vCase, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), vCase
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"VictimCaseForm1VictimCase"},
			AttrsValue: []interface{}{vCase.VictimCaseId},
		}

		if err = salvia_daos.GetVictimCaseForm1(by, &vCaseForm1, connData, &dbClientConfig, &dbServerConfig); err != nil {
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"VictimCaseForm2VictimCase"},
				AttrsValue: []interface{}{vCase.VictimCaseId},
			}

			if err = salvia_daos.GetVictimCaseForm2(by, &vCaseForm2, connData, &dbClientConfig, &dbServerConfig); err != nil {
				return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
			} else {
				vCase.VictimCaseForm2 = vCaseForm2
				// Se recupera la información geográfica del lugar de hechos (violencia) y residencia
				by = common_controllers.By{
					Operator:   common_dao.SQL_AND,
					AttrsName:  []string{"TownCode"},
					AttrsValue: []interface{}{vCaseForm2.VictimCaseForm2FactsTownCode},
				}
				err = security_daos.GetTown(by, &town2, connData, &dbClientConfig, &dbServerConfig)
				if err != nil {
					return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
				}

				by = common_controllers.By{
					Operator:   common_dao.SQL_AND,
					AttrsName:  []string{"TownCode"},
					AttrsValue: []interface{}{vCaseForm2.VictimCaseForm2ResidenceTownCode},
				}
				err = security_daos.GetTown(by, &town3, connData, &dbClientConfig, &dbServerConfig)
				if err != nil {
					return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
				}

				by = common_controllers.By{
					Operator:   common_dao.SQL_AND,
					AttrsName:  []string{"CityId"},
					AttrsValue: []interface{}{town3.TownCity},
				}
				err = security_daos.GetCity(by, &city3, connData, &dbClientConfig, &dbServerConfig)
				if err != nil {
					return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
				}

				by = common_controllers.By{
					Operator:   common_dao.SQL_AND,
					AttrsName:  []string{"DepartmentId"},
					AttrsValue: []interface{}{city3.CityDepartment},
				}
				err = security_daos.GetDepartment(by, &department3, connData, &dbClientConfig, &dbServerConfig)
				if err != nil {
					return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
				}

				//Ahora se cargan los enums sencillos
				loadVictimCaseEnums(&vCase)

				//Ahora se cargan los enums múltiples
				var enums []salvia_daos.VictimCaseForm2EnumsDTO

				//Cargamos todos los enums múltiples asociados al formulario
				enums, err = salvia_daos.GetVictimCasesForm2EnumsByVictimcaseForm2Id(vCase.VictimCaseForm2.VictimCaseForm2Id, connData, &dbClientConfig, &dbServerConfig)
				if err != nil {
					return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
				}

				loadVictimCaseEnumsMultiple(&vCase.VictimCaseForm2, enums)
			}
		} else {
			vCase.VictimCaseForm1 = vCaseForm1
			// Se recupera la información geográfica del lugar de hechos (violencia)
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"TownCode"},
				AttrsValue: []interface{}{vCaseForm1.VictimCaseForm1ViolenceTownCode},
			}
			err = security_daos.GetTown(by, &town2, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
			}
		}

		//Town2, department2 u Town y department son comúnes entonces los dejamos aquí afuera]
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town2.TownCity},
		}
		err = security_daos.GetCity(by, &city2, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city2.CityDepartment},
		}
		err = security_daos.GetDepartment(by, &department2, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
		}

		// Se recupera la información geográfica del lugar de atención.

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownCode"},
			AttrsValue: []interface{}{vCase.VictimCaseTownCode},
		}
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
		}

		var city security_daos.CityDTO
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town.TownCity},
		}
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city.CityDepartment},
		}
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
		}

		// Se recuperan los momentos (acciones/eventos) asociados al caso.
		var moments []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"MomentVictimCase"},
			AttrsValue: []interface{}{vCase.VictimCaseId},
		}
		moments, err = salvia_daos.GetMoments(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
		}

		// Para cada momento, se recuperan sus logs.
		var logs []salvia_daos.CaseLogDTO = []salvia_daos.CaseLogDTO{}
		for i := 0; i < len(moments); i++ {
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"CaseLogMoment"},
				AttrsValue: []interface{}{moments[i].MomentId},
			}
			logs, err = salvia_daos.GetCaseLogs(by, connData, &dbClientConfig, &dbServerConfig)
			if err == nil {
				moments[i].MomentCaseLogs = logs
			}
		}

		// Se recupera el seguimiento asociado al caso.
		var followUp salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}
		var resCode int
		if vCase.VictimCaseFollowUp.FollowUpId != 0 {
			resCode, _, followUp = GetFollowUpById(vCase.VictimCaseFollowUp.FollowUpId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			if resCode != http.StatusOK {
				return http.StatusInternalServerError, err.Error(), salvia_daos.VictimCaseDTO{}
			}
		}

		// Se asignan los datos geográficos, los momentos y el seguimiento al caso.
		vCase.VictimCaseTown = town
		vCase.VictimCaseDepartment = department
		vCase.VictimCaseCity = city
		vCase.VictimCaseMoments = moments

		if vCase.VictimCaseForm1.VictimCaseForm1Id > 0 {
			vCase.VictimCaseForm1.VictimCaseForm1ViolenceTown = town2
			vCase.VictimCaseForm1.VictimCaseForm1ViolenceDepartment = department2
			vCase.VictimCaseForm1.VictimCaseForm1ViolenceCity = city2
		} else if vCase.VictimCaseForm2.VictimCaseForm2Id > 0 {
			vCase.VictimCaseForm2.VictimCaseForm2FactsDepartment = department2
			vCase.VictimCaseForm2.VictimCaseForm2FactsCity = city2
			vCase.VictimCaseForm2.VictimCaseForm2FactsTown = town2
			vCase.VictimCaseForm2.VictimCaseForm2ResidenceDepartment = department3
			vCase.VictimCaseForm2.VictimCaseForm2ResidenceCity = city3
			vCase.VictimCaseForm2.VictimCaseForm2ResidenceTown = town3
		}

		vCase.VictimCaseFollowUp = followUp

		return http.StatusOK, utils.CommMsgGetJSONSuccess(vCase), vCase

	}
	return http.StatusInternalServerError, "", salvia_daos.VictimCaseDTO{}
}

// ApproveVictimCase aprueba un caso de víctima actualizando su estado a "ra" (aprobado).
func ApproveVictimCase(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crean los DTO para el caso.
	var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var vCaseToUpdate salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}

	// Se obtiene el caso por su identificador.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{id},
	}

	if err = salvia_daos.GetVictimCase(by, &vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualiza el estado del caso a "ra" (aprobado).
	vCaseToUpdate.VictimCaseStatus = "ra"

	// Se actualiza el estado del caso en la BD.
	if err = salvia_daos.UpdateVictimCaseStatus(&vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Commit de la transacción.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCase)
}

// UpdateVictimCaseStatusByICode actualiza el estado de un caso de víctima basado en su ICode.
func UpdateVictimCaseStatusByICode(newStatus string, id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCaseToUpdate salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{id},
	}

	if err = salvia_daos.GetVictimCase(by, &vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se asigna el nuevo estado.
	vCaseToUpdate.VictimCaseStatus = newStatus

	if err = salvia_daos.UpdateVictimCase(&vCaseToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCaseToUpdate)
}

// GetVictimCasesByTownCode obtiene todos los casos de víctimas asociados a un código de municipio.
func GetVictimCasesByTownCode(townCode string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
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

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	// Se consulta la BD utilizando el código del municipio.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseTownCode"},
		AttrsValue: []interface{}{townCode},
	}

	if vCases, count, err = salvia_daos.GetVictimCases(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

// GetDepartmentVictimCasesByTownCode obtiene todos los casos de víctimas asociados al departamento relacionado con el código de municipio.
func GetDepartmentVictimCasesByTownCode(townCode string, victimCaseStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
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

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	if vCases, count, err = salvia_daos.GetVictimCasesByDepartmentICode(department.DepartmentICode, victimCaseStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

func GetVictimCasesByTownCodeAndFollowUpStatusExcluded(townCode string, followUpStatusExcluded string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
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

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	if vCases, count, err = salvia_daos.GetVictimCasesByDepartmentICodeAndFollowUpStatusExcluded(department.DepartmentICode, followUpStatusExcluded, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

// GetVictimCasesByTownCodeWithAttend obtiene casos de víctimas por código de municipio incluyendo la información de atención.
func GetVictimCasesByTownCodeWithAttend(townCode string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var count int

	if townCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, count
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO

	if vCases, count, err = salvia_daos.GetVictimCasesByTownCodeWithAttend(townCode, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

// GetVictimCasesByTownCodeAndEntityBranchWithAttend obtiene casos de víctimas filtrados por código de municipio y código de rama de entidad.
func GetVictimCasesByTownCodeAndEntityBranchWithAttend(townCode string, entityBranchICode string, victimCaseStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var count int

	if townCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, count
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO

	if vCases, count, err = salvia_daos.GetVictimCasesByTownCodeAndEntityBranchWithAttend(townCode, entityBranchICode, victimCaseStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{vCases, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}

// GetVictimCasesByDocument obtiene casos de víctimas filtrados por tipo y número de documento.
func GetVictimCasesByDocument(docType string, docNumber string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseDocType", "VictimCaseDocNumber"},
		AttrsValue: []interface{}{docType, docNumber},
	}

	if vCases, count, err = salvia_daos.GetVictimCases(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{vCases, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}

// GetVictimCaseByDocumentAndCreationDate obtiene un caso de víctima utilizando documento y fecha de creación.
func GetVictimCaseByDocumentAndCreationDate(docType string, docNumber string, creationDate string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.VictimCaseDTO) {
	var err error = nil
	var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseCreationDate"},
		AttrsValue: []interface{}{docType, docNumber, creationDate},
	}

	if err = salvia_daos.GetVictimCase(by, &vCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, "", vCase
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCase), vCase
}

// GetVictimCasesByDocumentAndTownCode obtiene casos de víctimas filtrados por documento y código de municipio.
func GetVictimCasesByDocumentAndTownCode(docType string, docNumber string, townCode string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseTownCode"},
		AttrsValue: []interface{}{docType, docNumber, townCode},
	}

	if vCases, count, err = salvia_daos.GetVictimCases(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

// GetVictimCasesByDocumentAndTownCodeWithAttend obtiene casos de víctimas filtrados por documento y código de municipio incluyendo información de atención.
func GetVictimCasesByDocumentAndTownCodeWithAttend(docType string, docNumber string, townCode string, victimCaseStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseTownCode", "VictimCaseStatus"},
		AttrsValue: []interface{}{docType, docNumber, townCode, victimCaseStatus},
	}

	if vCases, count, err = salvia_daos.GetVictimCases(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

// GetVictimCasesByOwnerUserICode obtiene casos de víctimas asociados al ICode del usuario dueño.
func GetVictimCasesByOwnerUserICode(userIcode string, victimCaseStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var err error = nil
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	if vCases, count, err = salvia_daos.GetVictimCasesByOwnerUserICode(userIcode, victimCaseStatus, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{vCases, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"}), count
}

// GetVictimCasesByVictimUser obtiene casos de víctimas asociados a un usuario específico.
func GetVictimCasesByVictimUser(userICode string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var count int

	if userICode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, count
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseGeneralUser"},
		AttrsValue: []interface{}{userICode},
	}

	if vCases, count, err = salvia_daos.GetVictimCases(by, page, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), count
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases), count
}

// ReportVictimCases genera un reporte de casos de víctimas a partir de la entrada JSON
// que contiene los filtros de fecha, departamento, ciudad, etc.
func ReportVictimCases(dataInput string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	var dtoMap map[string]interface{} = nil
	var report salvia_daos.VictimCaseReportDTO

	// Definición de campos requeridos para el reporte.
	var checkFields map[string]bool = map[string]bool{
		"VictimCaseReportStartDate":    true,
		"VictimCaseReportEndDate":      true,
		"VictimCaseReportDepartment":   false,
		"VictimCaseReportCity":         false,
		"VictimCaseReportTown":         false,
		"VictimCaseReportViolenceType": false,
		"VictimCaseReportCaseStatus":   false,
	}

	// Se obtiene el DTO a partir del JSON.
	var rawDto interface{} = utils.GetDTOMap(dataInput, security_daos.GeneralUserJSONName, common_config.Locale, collectedErrors)

	var err error = nil

	var vCases []salvia_daos.VictimCaseDTO = []salvia_daos.VictimCaseDTO{}

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto
		utils.ValidateJSONInput(&report, dtoMap, salvia_daos.VictimCaseReportJSONName, salvia_daos.VictimCaseReportFieldDefinitions,
			checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT,
			common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
	default:
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default",
			common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	if vCases, err = salvia_daos.GetVictimCasesReport(report, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	//Ahora cargamos los enums de los registros que sean de formulario v2
	for idx := 0; idx < len(vCases); idx++ {

		if vCases[idx].VictimCaseForm2.VictimCaseForm2Id > 0 {
			loadVictimCaseEnums(&vCases[idx])

			var enums []salvia_daos.VictimCaseForm2EnumsDTO

			//Cargamos todos los enums múltiples asociados al formulario
			enums, err = salvia_daos.GetVictimCasesForm2EnumsByVictimcaseForm2Id(vCases[idx].VictimCaseForm2.VictimCaseForm2Id, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, err.Error()
			}
			loadVictimCaseEnumsMultiple(&vCases[idx].VictimCaseForm2, enums)
		}

	}
	return http.StatusOK, utils.CommMsgGetJSONSuccess(vCases)
}

// GetVictimCaseByAll retorna todos los casos de víctimas almacenados en la base de datos.
func GetVictimCaseByAll(victimCaseStatus string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	var count int

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	cases, count, err := salvia_daos.GetAllVictimCases(victimCaseStatus, page, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{cases, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"})
	}
	return resCode, resData, count
}

// AssignOperators asigna un nuevo operador a todos los casos de víctimas que
// pertenecen al operador anterior. Se actualizan las relaciones y se inactiva la antigua.
func AssignOperators(oldOperatorICode string, newOperatorICode string, s utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen los DTO de los dueños (casos) asociados a cada operador.
	var oldOwner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	var newOwner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	var relCaseOwners []salvia_daos.RelCaseOwnerVictimCaseDTO = []salvia_daos.RelCaseOwnerVictimCaseDTO{}

	// Obtiene el dueño asociado al operador antiguo.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{oldOperatorICode},
	}
	err = salvia_daos.GetCaseOwner(by, &oldOwner, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Obtiene el dueño asociado al nuevo operador.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{newOperatorICode},
	}
	err = salvia_daos.GetCaseOwner(by, &newOwner, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se obtienen las relaciones del dueño antiguo con sus casos.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelCaseOwnerVictimCase_CaseOwner"},
		AttrsValue: []interface{}{oldOwner.CaseOwnerId},
	}
	relCaseOwners, err = salvia_daos.GetRelCaseOwnerVictimCases(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_contact_error_loading_owner"]
	}

	// Se inicia la transacción para actualizar las relaciones.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualizan cada una de las relaciones: se asigna el nuevo dueño y se inactiva la anterior.
	for _, rel := range relCaseOwners {
		var relCaseOwnerVictimCase salvia_daos.RelCaseOwnerVictimCaseDTO = salvia_daos.RelCaseOwnerVictimCaseDTO{}
		salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&relCaseOwnerVictimCase, common_dao.SQL_INSERT)
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner = newOwner.CaseOwnerId
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase = rel.RelCaseOwnerVictimCase_VictimCase

		if err = salvia_daos.SetRelCaseOwnerVictimCase(&relCaseOwnerVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&rel, common_dao.SQL_UPDATE)
		if err = salvia_daos.UpdateRelCaseOwnerVictimCase(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
		if err = salvia_daos.UpdateVictimCaseOwnersAndRolesByVictimCaseId(rel.RelCaseOwnerVictimCase_VictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, ""
}

// AssignOperator asigna un nuevo operador a un caso de víctima específico.
// Actualiza las relaciones de autoría según el nuevo responsable.
func AssignOperator(caseICode string, newOperatorICode string, s utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen los DTO del caso y del nuevo dueño.
	var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var newOwner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	var relCaseOwners []salvia_daos.RelCaseOwnerVictimCaseDTO = []salvia_daos.RelCaseOwnerVictimCaseDTO{}

	// Se carga el caso a partir de su ICode.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{caseICode},
	}

	err = salvia_daos.GetVictimCase(by, &vCase, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_case_error_loading_vcase"]
	}

	// Se obtiene el nuevo dueño basado en el nuevo operador.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{newOperatorICode},
	}
	err = salvia_daos.GetCaseOwner(by, &newOwner, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_case_error_loading_new_owner"]
	}

	// Se obtienen las relaciones activas del caso.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelCaseOwnerVictimCase_VictimCase", "RelCaseOwnerVictimCase_Status"},
		AttrsValue: []interface{}{vCase.VictimCaseId, "a"},
	}

	relCaseOwners, err = salvia_daos.GetRelCaseOwnerVictimCases(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["rel_case_owner_error_loading_victim_case"]
	}

	if len(relCaseOwners) == 0 {
		return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["rel_case_owner_error_not_found"]
	}
	// Se recorre cada relación para actualizar el dueño.
	for _, rel := range relCaseOwners {
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CaseOwnerId"},
			AttrsValue: []interface{}{rel.RelCaseOwnerVictimCase_CaseOwner},
		}
		var owner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
		err = salvia_daos.GetCaseOwner(by, &owner, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_case_error_loading_owner"]
		}

		// Se verifica que el usuario del dueño no sea el mismo que el nuevo.
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserICode"},
			AttrsValue: []interface{}{owner.CaseOwnerGeneralUser},
		}
		var user security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
		err = security_daos.GetGeneralUser(by, &user, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, salvia_config.Locale[s.Lang]["victim_case_error_loading_user"]
		}
		// Si el usuario actual es distinto al nuevo, se procede a la actualización.
		if user.GeneralUserICode != newOperatorICode {
			var found bool = false
			for _, r := range user.GeneralUserRoles {
				if r.RoleCode == "op" {
					found = true
					break
				}
			}

			if found {
				if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
					return http.StatusInternalServerError, err.Error()
				}

				var relCaseOwnerVictimCase salvia_daos.RelCaseOwnerVictimCaseDTO = salvia_daos.RelCaseOwnerVictimCaseDTO{}
				salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&relCaseOwnerVictimCase, common_dao.SQL_INSERT)
				relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner = newOwner.CaseOwnerId
				relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase = rel.RelCaseOwnerVictimCase_VictimCase

				if err = salvia_daos.SetRelCaseOwnerVictimCase(&relCaseOwnerVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
					db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
					return http.StatusInternalServerError, err.Error()
				}

				salvia_daos.SetRelCaseOwnerVictimCaseDefaults(&rel, common_dao.SQL_UPDATE)
				if err = salvia_daos.UpdateRelCaseOwnerVictimCase(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
					db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
					return http.StatusInternalServerError, err.Error()
				}

				if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
					return http.StatusInternalServerError, err.Error()
				}
			}
		}
	}

	//Finalmente se actualiza la columna con los datos informativos sobre los funcionarios intervinientes
	if err = salvia_daos.UpdateVictimCaseOwnersAndRolesByVictimCaseId(vCase.VictimCaseId, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, ""
}

// Esta función no retorna nada porque está diseñada para ser ejecutada mediante TaskScheduler y sólo muestra errores en consola.
// Actualiza automaticamente la columna de operarios y roles para que quede quemada en la tabla y no haya que hacer la misma consulta cada vez que listan casos
func UpdateVictimCasesOwnersAndRoles(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) {

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se ejecuta a bajo nivel
	salvia_daos.UpdateVictimCasesOwnersAndRoles(connData, &dbClientConfig, &dbServerConfig)

}

// getRandomUserName genera un nombre de usuario aleatorio con la longitud especificada.
func getRandomUserName(size int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// getRandomPassword genera una contraseña aleatoria con la longitud especificada.
func getRandomPassword(size int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// hasRole verifica si el usuario posee un rol específico.
func hasRole(roleCode string, roles []security_daos.RoleDTO) bool {
	for _, r := range roles {
		if r.RoleCode == roleCode {
			return true
		}
	}
	return false
}

// cleanMoments compara dos slices de momentos y retorna dos slices: uno con los momentos a eliminar
// y otro con los nuevos momentos a crear.
func cleanMoments(momentsToDelete []salvia_daos.MomentDTO, moments []salvia_daos.MomentDTO) ([]salvia_daos.MomentDTO, []salvia_daos.MomentDTO) {

	for i := len(momentsToDelete) - 1; i >= 0; i-- {
		var idx int = getMomentIn(momentsToDelete[i], moments)
		// Si el momento ya existe en ambos slices, se elimina de ambos.
		if idx != -1 {
			momentsToDelete = removeMomentFromSlice(momentsToDelete, i)
			moments = removeMomentFromSlice(moments, idx)
		}
	}
	return momentsToDelete, moments
}

// getMomentIn busca un momento dentro de un slice y retorna su índice; si no se encuentra, retorna -1.
func getMomentIn(momentToFind salvia_daos.MomentDTO, moments []salvia_daos.MomentDTO) int {
	for i := 0; i < len(moments); i++ {
		if momentToFind.MomentCode == moments[i].MomentCode && momentToFind.MomentEntityBranch.EntityBranchId == moments[i].MomentEntityBranch.EntityBranchId {
			return i
		}
	}
	return -1
}

// removeMomentFromSlice elimina el elemento en el índice dado de un slice de momentos.
func removeMomentFromSlice(slice []salvia_daos.MomentDTO, idx int) []salvia_daos.MomentDTO {
	return append(slice[:idx], slice[idx+1:]...)
}

// cleanEnums compara dos slices de enumeraciones y retorna dos slices: uno con los enums a eliminar
// y otro con los nuevos enums a crear.
func cleanEnums(enumsToDelete []salvia_daos.VictimCaseForm2EnumsDTO, enums []salvia_daos.VictimCaseForm2EnumsDTO) ([]salvia_daos.VictimCaseForm2EnumsDTO, []salvia_daos.VictimCaseForm2EnumsDTO) {

	for i := len(enumsToDelete) - 1; i >= 0; i-- {
		var idx int = getEnumIn(enumsToDelete[i], enums)
		// Si el enumo ya existe en ambos slices, se elimina de ambos.
		if idx != -1 {
			enumsToDelete = removeEnumFromSlice(enumsToDelete, i)
			enums = removeEnumFromSlice(enums, idx)
		}
	}
	return enumsToDelete, enums
}

// getEnumIn busca un enumo dentro de un slice y retorna su índice; si no se encuentra, retorna -1.
func getEnumIn(enumToFind salvia_daos.VictimCaseForm2EnumsDTO, enums []salvia_daos.VictimCaseForm2EnumsDTO) int {
	for i := 0; i < len(enums); i++ {
		if enumToFind.VictimCaseForm2EnumsICode == enums[i].VictimCaseForm2EnumsICode {
			return i
		}
	}
	return -1
}

// removeEnumFromSlice elimina el elemento en el índice dado de un slice de enumos.
func removeEnumFromSlice(slice []salvia_daos.VictimCaseForm2EnumsDTO, idx int) []salvia_daos.VictimCaseForm2EnumsDTO {
	return append(slice[:idx], slice[idx+1:]...)
}

func getAndVerifyVictimCaseEnums(vCaseForm2 *salvia_daos.VictimCaseForm2DTO, collectedErrors map[string]map[string]string, verify bool) bool {
	var opRes bool = true

	//VictimCaseForm2PersonWithDisability
	err := salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PersonWithDisability)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PersonWithDisability", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2RequireLanguageInterpreter
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2RequireLanguageInterpreter)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2RequireLanguageInterpreter", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2ViolenceMotivatedByGender
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ViolenceMotivatedByGender)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ViolenceMotivatedByGender", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorOccupation
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorOccupation)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorOccupation", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2SupportContactKinship
	if vCaseForm2.VictimCaseForm2SupportContactKinship.VictimCaseForm2EnumsICode != "" {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2SupportContactKinship)
		if err != nil && verify {
			// Retorna error interno en caso de fallo.
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2SupportContactKinship", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	//VictimCaseForm2FactsZone
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2FactsZone)
	if err != nil && verify {
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2FactsZone", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2ScenarioViolence
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ScenarioViolence)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ScenarioViolence", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2ReportedPreviously
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ReportedPreviously)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ReportedPreviously", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if vCaseForm2.VictimCaseForm2ReportedPreviously.VictimCaseForm2EnumsCode == "y" {
		//VictimCaseForm2AttentionWasAppropriate
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AttentionWasAppropriate)
		if err != nil && verify {
			// Retorna error interno en caso de fallo.
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AttentionWasAppropriate", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	// VictimCaseForm2RecurrenceAggression
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2RecurrenceAggression)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2RecurrenceAggression", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2NumAgressors
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2NumAgressors)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2NumAgressors", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2ProximityPrincipalAggressor
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ProximityPrincipalAggressor)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ProximityPrincipalAggressor", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2RelationshipWithPresumedAggressor
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2RelationshipWithPresumedAggressor", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	var partnerKnown bool = (vCaseForm2.VictimCaseForm2ProximityPrincipalAggressor.VictimCaseForm2EnumsCode == "pc" || vCaseForm2.VictimCaseForm2ProximityPrincipalAggressor.VictimCaseForm2EnumsCode == "pn")

	var partnerValidated bool = vCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor.VictimCaseForm2EnumsCode == "pi" || vCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor.VictimCaseForm2EnumsCode == "ex"

	// VictimCaseForm2EconomicallyDependent
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2EconomicallyDependent)
	//Depende de VictimCaseForm2ProximityPrincipalAggressor
	if err != nil && verify && partnerKnown {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2EconomicallyDependent", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2AggressorGenderIdentity
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorGenderIdentity)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorGenderIdentity", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorDocType
	if vCaseForm2.VictimCaseForm2AggressorDocType.VictimCaseForm2EnumsICode != "" {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorDocType)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorDocType", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	//VictimCaseForm2AggressorViolencePhysicalIncrease
	var violentValidation bool = false
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorViolencePhysicalIncrease", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else {
		violentValidation = violentValidation || vCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease.VictimCaseForm2EnumsCode == "y"
	}

	//VictimCaseForm2AggressorWeaponUsed
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorWeaponUsed)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorWeaponUsed", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else {
		violentValidation = violentValidation || vCaseForm2.VictimCaseForm2AggressorWeaponUsed.VictimCaseForm2EnumsCode == "y"
	}

	//VictimCaseForm2AggressorThreatKill
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorThreatKill)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorThreatKill", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else {
		violentValidation = violentValidation || vCaseForm2.VictimCaseForm2AggressorThreatKill.VictimCaseForm2EnumsCode == "y"
	}

	//VictimCaseForm2AggressorPursuesSpiesDestroys
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorPursuesSpiesDestroys", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else {
		violentValidation = violentValidation || vCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys.VictimCaseForm2EnumsCode == "y"
	}

	//VictimCaseForm2AggressorCapableOfKilling
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorCapableOfKilling)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorCapableOfKilling", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else {
		violentValidation = violentValidation || vCaseForm2.VictimCaseForm2AggressorCapableOfKilling.VictimCaseForm2EnumsCode == "y"
	}

	//VictimCaseForm2AggressorHasAccessToWeapons
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorHasAccessToWeapons", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else {
		violentValidation = violentValidation || vCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons.VictimCaseForm2EnumsCode == "y"
	}

	//VictimCaseForm2PartnerUnemployed
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PartnerUnemployed)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PartnerUnemployed", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2PartnerOtherDenunciations
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PartnerOtherDenunciations)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PartnerOtherDenunciations", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorHasPenalBackground
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorHasPenalBackground)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorHasPenalBackground", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorForcedSex
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorForcedSex)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorForcedSex", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorAttemptedStrangulation
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorAttemptedStrangulation)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorAttemptedStrangulation", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorConsumesDrugs
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorConsumesDrugs)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorConsumesDrugs", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorIsAlcoholic
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorIsAlcoholic)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorIsAlcoholic", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2PartnerControls
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PartnerControls)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PartnerControls", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorHadHitInVulnerability
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorHadHitInVulnerability)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorHadHitInVulnerability", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2PartnerThreatenedSuicide
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PartnerThreatenedSuicide)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PartnerThreatenedSuicide", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2PartnerThreatenedDamageMembers
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PartnerThreatenedDamageMembers)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PartnerThreatenedDamageMembers", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2ThoughtsOfSelfHarm
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ThoughtsOfSelfHarm", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorLimitsContactSupportNetworks
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorLimitsContactSupportNetworks)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorLimitsContactSupportNetworks", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2StillLivesWithAggressor
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2StillLivesWithAggressor)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2StillLivesWithAggressor", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorViolentlyJealous
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorViolentlyJealous)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorViolentlyJealous", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2StoppedSeekingHelp
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2StoppedSeekingHelp)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2StoppedSeekingHelp", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2VictimHealthToBlackmail
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2VictimHealthToBlackmail)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2VictimHealthToBlackmail", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2ThreatenedRevealSexualOrientation
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ThreatenedRevealSexualOrientation)
	if err != nil && verify && violentValidation && partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ThreatenedRevealSexualOrientation", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorUnemployed
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorUnemployed)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorUnemployed", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorHasPenalBackground2
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorHasPenalBackground2)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorHasPenalBackground2", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorSexuallyHarassment
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorSexuallyHarassment)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorSexuallyHarassment", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorUseDrugs
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorUseDrugs)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorUseDrugs", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorIsAlcoholic2
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorIsAlcoholic2)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorIsAlcoholic2", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorControls
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorControls)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorControls", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorThreatenedDamageMembers
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorThreatenedDamageMembers)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorThreatenedDamageMembers", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2ThoughtsOfSelfHarm2
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm2)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ThoughtsOfSelfHarm2", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorCommonSpaces
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorCommonSpaces)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorCommonSpaces", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorHierarchy
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorHierarchy)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorHierarchy", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2ViolenceMotivatedByGender2
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ViolenceMotivatedByGender2)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ViolenceMotivatedByGender2", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorSexuallyHarassment2
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorSexuallyHarassment2)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorSexuallyHarassment2", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	//VictimCaseForm2AggressorUsedPositionAuthority
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AggressorUsedPositionAuthority)
	if err != nil && verify && violentValidation && !partnerValidated {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AggressorUsedPositionAuthority", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2PhysicalMentalSensoryDifficulties
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2PhysicalMentalSensoryDifficulties)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2PhysicalMentalSensoryDifficulties", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2Nationality
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2Nationality)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2Nationality", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if vCaseForm2.VictimCaseForm2Nationality.VictimCaseForm2EnumsCode == "ex" {
		// VictimCaseForm2SpecifiedNationality
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2SpecifiedNationality)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2SpecifiedNationality", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}

		// VictimCaseForm2MigrationCondition
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2MigrationCondition)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2MigrationCondition", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	// VictimCaseForm2GenderIdentity
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2GenderIdentity)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2GenderIdentity", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2SexualOrientation
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2SexualOrientation)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2SexualOrientation", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2AssignedSexAtBirth
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AssignedSexAtBirth)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AssignedSexAtBirth", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2EthnicAffiliation
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2EthnicAffiliation)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2EthnicAffiliation", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if vCaseForm2.VictimCaseForm2EthnicAffiliation.VictimCaseForm2EnumsCode == "in" {
		// VictimCaseForm2IndigenousPeople
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2IndigenousPeople)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2IndigenousPeople", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	// VictimCaseForm2CampesinoRecognition
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2CampesinoRecognition)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2CampesinoRecognition", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2MaritalStatus
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2MaritalStatus)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2MaritalStatus", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2LastEducationLevel
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2LastEducationLevel)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2LastEducationLevel", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2Occupation
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2Occupation)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2Occupation", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2IncomeGenerationMethod
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2IncomeGenerationMethod)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2IncomeGenerationMethod", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	} else if vCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsCode == "em" {
		// VictimCaseForm2EmploymentRelationship
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2EmploymentRelationship)
		if err != nil && verify {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2EmploymentRelationship", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			opRes = false
		}
	}

	// VictimCaseForm2HousingTenancyForm
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2HousingTenancyForm)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2HousingTenancyForm", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2HousingStratum
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2HousingStratum)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2HousingStratum", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2CurrentlyPregnant
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2CurrentlyPregnant)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2CurrentlyPregnant", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	// VictimCaseForm2ResidenceZone
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ResidenceZone)
	if err != nil && verify {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ResidenceZone", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

func loadVictimCaseEnums(vCase *salvia_daos.VictimCaseDTO) {
	//VictimCaseForm2FactsZone
	//salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseDocType)

	//VictimCaseForm2FactsZone
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2FactsZone)

	// VictimCaseForm2ScenarioViolence
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ScenarioViolence)

	// VictimCaseForm2ReportedPreviously
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ReportedPreviously)

	// VictimCaseForm2RecurrenceAggression
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2RecurrenceAggression)

	// VictimCaseForm2NumAgressors
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2NumAgressors)

	// VictimCaseForm2ProximityPrincipalAggressor
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ProximityPrincipalAggressor)

	// VictimCaseForm2RelationshipWithPresumedAggressor
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor)

	// VictimCaseForm2EconomicallyDependent
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2EconomicallyDependent)

	// VictimCaseForm2AggressorGenderIdentity
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorGenderIdentity)

	//VictimCaseForm2AggressorDocType
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorDocType)

	//VictimCaseForm2AggressorViolencePhysicalIncrease
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease)

	//VictimCaseForm2AggressorWeaponUsed
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorWeaponUsed)

	//VictimCaseForm2AggressorThreatKill
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorThreatKill)

	//VictimCaseForm2AggressorPursuesSpiesDestroys
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys)

	//VictimCaseForm2AggressorCapableOfKilling
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorCapableOfKilling)

	//VictimCaseForm2AggressorHasAccessToWeapons
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons)

	//VictimCaseForm2PartnerUnemployed
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PartnerUnemployed)

	//VictimCaseForm2PartnerOtherDenunciations
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PartnerOtherDenunciations)

	//VictimCaseForm2AggressorHasPenalBackground
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorHasPenalBackground)

	//VictimCaseForm2AggressorForcedSex
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorForcedSex)

	//VictimCaseForm2AggressorAttemptedStrangulation
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorAttemptedStrangulation)

	//VictimCaseForm2AggressorConsumesDrugs
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorConsumesDrugs)

	//VictimCaseForm2AggressorIsAlcoholic
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorIsAlcoholic)

	//VictimCaseForm2PartnerControls
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PartnerControls)

	//VictimCaseForm2AggressorHadHitInVulnerability
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorHadHitInVulnerability)

	//VictimCaseForm2PartnerThreatenedSuicide
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PartnerThreatenedSuicide)

	//VictimCaseForm2PartnerThreatenedDamageMembers
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PartnerThreatenedDamageMembers)

	//VictimCaseForm2ThoughtsOfSelfHarm
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm)

	//VictimCaseForm2AggressorLimitsContactSupportNetworks
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorLimitsContactSupportNetworks)

	//VictimCaseForm2StillLivesWithAggressor
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2StillLivesWithAggressor)

	//VictimCaseForm2AggressorViolentlyJealous
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorViolentlyJealous)

	//VictimCaseForm2AggressorUnemployed
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorUnemployed)

	//VictimCaseForm2AggressorHasPenalBackground2
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorHasPenalBackground2)

	//VictimCaseForm2AggressorSexuallyHarassment
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment)

	//VictimCaseForm2AggressorUseDrugs
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorUseDrugs)

	//VictimCaseForm2AggressorIsAlcoholic2
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorIsAlcoholic2)

	//VictimCaseForm2AggressorControls
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorControls)

	//VictimCaseForm2AggressorThreatenedDamageMembers
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorThreatenedDamageMembers)

	//VictimCaseForm2ThoughtsOfSelfHarm2
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm2)

	//VictimCaseForm2AggressorCommonSpaces
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorCommonSpaces)

	//VictimCaseForm2AggressorHierarchy
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorHierarchy)

	// VictimCaseForm2PhysicalMentalSensoryDifficulties
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PhysicalMentalSensoryDifficulties)

	// VictimCaseForm2Nationality
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2Nationality)

	// VictimCaseForm2SpecifiedNationality
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2SpecifiedNationality)

	// VictimCaseForm2MigrationCondition
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2MigrationCondition)

	// VictimCaseForm2GenderIdentity
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2GenderIdentity)

	// VictimCaseForm2SexualOrientation
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2SexualOrientation)

	// VictimCaseForm2AssignedSexAtBirth
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AssignedSexAtBirth)

	// VictimCaseForm2EthnicAffiliation
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2EthnicAffiliation)

	// VictimCaseForm2IndigenousPeople
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2IndigenousPeople)

	// VictimCaseForm2CampesinoRecognition
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2CampesinoRecognition)

	// VictimCaseForm2MaritalStatus
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2MaritalStatus)

	// VictimCaseForm2LastEducationLevel
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2LastEducationLevel)

	// VictimCaseForm2Occupation
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2Occupation)

	// VictimCaseForm2IncomeGenerationMethod
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2IncomeGenerationMethod)

	// VictimCaseForm2EmploymentRelationship
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2EmploymentRelationship)

	// VictimCaseForm2HousingTenancyForm
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2HousingTenancyForm)

	// VictimCaseForm2HousingStratum
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2HousingStratum)

	// VictimCaseForm2CurrentlyPregnant
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2CurrentlyPregnant)

	// VictimCaseForm2ResidenceZone
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ResidenceZone)

	// VictimCaseForm2SupportContactKinship
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2SupportContactKinship)

	//VictimCaseForm2PersonWithDisability
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2PersonWithDisability)

	//VictimCaseForm2RequireLanguageInterpreter
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2RequireLanguageInterpreter)

	//VictimCaseForm2WorkplaceSectorOccurrence
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2WorkplaceSectorOccurrence)

	//VictimCaseForm2ViolenceMotivatedByGender
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender)

	//VictimCaseForm2AttentionWasAppropriate
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AttentionWasAppropriate)

	//VictimCaseForm2AggressorOccupation
	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorOccupation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2StoppedSeekingHelp)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2VictimHealthToBlackmail)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ThreatenedRevealSexualOrientation)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2ViolenceMotivatedByGender2)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorSexuallyHarassment2)

	salvia_daos.GetLocalVictimCaseForm2EnumsById(&vCase.VictimCaseForm2.VictimCaseForm2AggressorUsedPositionAuthority)

}

func getAndVerifyVictimCaseEnumsMultiple(vCaseForm2 *salvia_daos.VictimCaseForm2DTO, collectedErrors map[string]map[string]string) bool {
	var opRes bool = true
	var err error
	var violenceSubTypes []salvia_daos.VictimCaseForm2EnumsDTO = []salvia_daos.VictimCaseForm2EnumsDTO{}
	for idx := range vCaseForm2.VictimCaseForm2TypeViolenceExperienced {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2TypeViolenceExperienced[idx])
		if err != nil {
			opRes = false
			break
		} else {
			var enumTmp salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsCategory: "victim_case_form2_subtype_violence_experienced_" + vCaseForm2.VictimCaseForm2TypeViolenceExperienced[idx].VictimCaseForm2EnumsCode}
			var typesTmp []salvia_daos.VictimCaseForm2EnumsDTO = salvia_daos.GetLocalVictimCaseForm2EnumsByCategory(enumTmp)
			violenceSubTypes = append(violenceSubTypes, typesTmp...)
		}
	}
	if len(vCaseForm2.VictimCaseForm2TypeViolenceExperienced) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2TypeViolenceExperienced", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx, t := range vCaseForm2.VictimCaseForm2SubtypeViolenceExperienced {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2SubtypeViolenceExperienced[idx])
		if err != nil {
			opRes = false
			break
		} else {
			//Verificamos las opciones
			var optVerified bool = false
			for _, e := range violenceSubTypes {
				if e.VictimCaseForm2EnumsICode == t.VictimCaseForm2EnumsICode {
					optVerified = true
					break
				}
			}
			if !optVerified {
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2SubtypeViolenceExperienced", "json"), "victim_case_form2_enums_incorrect_choise", "common_global_error", salvia_config.Locale)
				opRes = false
			}
		}
	}

	if len(vCaseForm2.VictimCaseForm2SubtypeViolenceExperienced) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2SubtypeViolenceExperienced", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2ScopeOfViolence {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ScopeOfViolence[idx])
		if err != nil {
			opRes = false
			break
		}

		if vCaseForm2.VictimCaseForm2ScopeOfViolence[idx].VictimCaseForm2EnumsCode == "al" {
			//Validamos uno de selección sencilla que depende de una opción aquí
			//VictimCaseForm2WorkplaceSectorOccurrence
			err2 := salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2WorkplaceSectorOccurrence)
			if err2 != nil {
				// Retorna error interno en caso de fallo.
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2WorkplaceSectorOccurrence", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
				opRes = false
			}
		}
	}

	if len(vCaseForm2.VictimCaseForm2ScopeOfViolence) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ScopeOfViolence", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2WhoReportTo {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2WhoReportTo[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if vCaseForm2.VictimCaseForm2ReportedPreviously.VictimCaseForm2EnumsCode == "y" && len(vCaseForm2.VictimCaseForm2WhoReportTo) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2WhoReportTo", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2AdjustmentsGBV {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2AdjustmentsGBV[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if len(vCaseForm2.VictimCaseForm2AdjustmentsGBV) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2AdjustmentsGBV", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2Law1996 {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2Law1996[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if vCaseForm2.VictimCaseForm2PhysicalMentalSensoryDifficulties.VictimCaseForm2EnumsCode == "y" && len(vCaseForm2.VictimCaseForm2Law1996) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2Law1996", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2SpeciallyProtectedPopulation {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2SpeciallyProtectedPopulation[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if len(vCaseForm2.VictimCaseForm2SpeciallyProtectedPopulation) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2SpeciallyProtectedPopulation", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2ASPMode {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ASPMode[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if vCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsCode == "pr" && len(vCaseForm2.VictimCaseForm2ASPMode) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ASPMode", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2ReasonASP {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ReasonASP[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if vCaseForm2.VictimCaseForm2IncomeGenerationMethod.VictimCaseForm2EnumsCode == "pr" && len(vCaseForm2.VictimCaseForm2ReasonASP) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ReasonASP", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2HasDependents {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2HasDependents[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if len(vCaseForm2.VictimCaseForm2HasDependents) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2HasDependents", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	for idx := range vCaseForm2.VictimCaseForm2ActionPlan {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vCaseForm2.VictimCaseForm2ActionPlan[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if len(vCaseForm2.VictimCaseForm2ActionPlan) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vCaseForm2, "VictimCaseForm2ActionPlan", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

// En enums están todos los del caso actual, hay que extraer los necesarios para llenar el objeto
func loadVictimCaseEnumsMultiple(vCaseForm2 *salvia_daos.VictimCaseForm2DTO, enums []salvia_daos.VictimCaseForm2EnumsDTO) {
	var err error

	for _, e := range enums {
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsById(&e); err == nil {
			switch e.VictimCaseForm2EnumsCategory {
			case "victim_case_form2_type_violence_experienced":
				vCaseForm2.VictimCaseForm2TypeViolenceExperienced = append(vCaseForm2.VictimCaseForm2TypeViolenceExperienced, e)

			case "victim_case_form2_subtype_violence_experienced_fi",
				"victim_case_form2_subtype_violence_experienced_ps",
				"victim_case_form2_subtype_violence_experienced_se",
				"victim_case_form2_subtype_violence_experienced_po",
				"victim_case_form2_subtype_violence_experienced_pl",
				"victim_case_form2_subtype_violence_experienced_re",
				"victim_case_form2_subtype_violence_experienced_vi":
				vCaseForm2.VictimCaseForm2SubtypeViolenceExperienced = append(vCaseForm2.VictimCaseForm2SubtypeViolenceExperienced, e)

			case "victim_case_form2_scope_of_violence":
				vCaseForm2.VictimCaseForm2ScopeOfViolence = append(vCaseForm2.VictimCaseForm2ScopeOfViolence, e)

			case "victim_case_form2_who_report_to":
				vCaseForm2.VictimCaseForm2WhoReportTo = append(vCaseForm2.VictimCaseForm2WhoReportTo, e)

			case "victim_case_form2_adjustments_gbv":
				vCaseForm2.VictimCaseForm2AdjustmentsGBV = append(vCaseForm2.VictimCaseForm2AdjustmentsGBV, e)

			case "victim_case_form2_law_1996":
				vCaseForm2.VictimCaseForm2Law1996 = append(vCaseForm2.VictimCaseForm2Law1996, e)

			case "victim_case_form2_specially_protected_population":
				vCaseForm2.VictimCaseForm2SpeciallyProtectedPopulation = append(vCaseForm2.VictimCaseForm2SpeciallyProtectedPopulation, e)

			case "victim_case_form2_asp_mode":
				vCaseForm2.VictimCaseForm2ASPMode = append(vCaseForm2.VictimCaseForm2ASPMode, e)

			case "victim_case_form2_reason_asp":
				vCaseForm2.VictimCaseForm2ReasonASP = append(vCaseForm2.VictimCaseForm2ReasonASP, e)

			case "victim_case_form2_has_dependents":
				vCaseForm2.VictimCaseForm2HasDependents = append(vCaseForm2.VictimCaseForm2HasDependents, e)

			case "victim_case_form2_action_plan":
				vCaseForm2.VictimCaseForm2ActionPlan = append(vCaseForm2.VictimCaseForm2ActionPlan, e)
			}

		}
	}
}

// Dado que ya están validados, entonces se relacionan los enums con contenido
func setVictimCaseEnumsMultiple(vCaseForm2 *salvia_daos.VictimCaseForm2DTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	var err error
	for _, e := range vCaseForm2.VictimCaseForm2TypeViolenceExperienced {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2SubtypeViolenceExperienced {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2ScopeOfViolence {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2WhoReportTo {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2ActivitiesUnableToPerform {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2AdjustmentsGBV {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2Law1996 {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2SpeciallyProtectedPopulation {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2ASPMode {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2ReasonASP {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2HasDependents {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	for _, e := range vCaseForm2.VictimCaseForm2ActionPlan {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = *vCaseForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	return nil
}

// Dado que ya están validados, entonces se actualizan los enums con contenido
func updateVictimCaseEnumsMultiple(vCaseFormToUpdate salvia_daos.VictimCaseForm2DTO, vCaseForm salvia_daos.VictimCaseForm2DTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	var err error

	var enumsToDelete []salvia_daos.VictimCaseForm2EnumsDTO = []salvia_daos.VictimCaseForm2EnumsDTO{}
	var currentEnums []salvia_daos.VictimCaseForm2EnumsDTO = []salvia_daos.VictimCaseForm2EnumsDTO{}

	//VictimCaseForm2TypeViolenceExperienced -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2TypeViolenceExperienced
	currentEnums = vCaseForm.VictimCaseForm2TypeViolenceExperienced

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2SubtypeViolenceExperienced -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2SubtypeViolenceExperienced
	currentEnums = vCaseForm.VictimCaseForm2SubtypeViolenceExperienced

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2ScopeOfViolence -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2ScopeOfViolence
	currentEnums = vCaseForm.VictimCaseForm2ScopeOfViolence

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2WhoReportTo -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2WhoReportTo
	currentEnums = vCaseForm.VictimCaseForm2WhoReportTo

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2ActivitiesUnableToPerform -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2ActivitiesUnableToPerform
	currentEnums = vCaseForm.VictimCaseForm2ActivitiesUnableToPerform

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2AdjustmentsGBV -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2AdjustmentsGBV
	currentEnums = vCaseForm.VictimCaseForm2AdjustmentsGBV

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2Law1996 -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2Law1996
	currentEnums = vCaseForm.VictimCaseForm2Law1996

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2SpeciallyProtectedPopulation -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2SpeciallyProtectedPopulation
	currentEnums = vCaseForm.VictimCaseForm2SpeciallyProtectedPopulation

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2ASPMode -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2ASPMode
	currentEnums = vCaseForm.VictimCaseForm2ASPMode

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2ReasonASP -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2ReasonASP
	currentEnums = vCaseForm.VictimCaseForm2ReasonASP

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2HasDependents -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2HasDependents
	currentEnums = vCaseForm.VictimCaseForm2HasDependents

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	//VictimCaseForm2ActionPlan -------------------------------------
	enumsToDelete = vCaseFormToUpdate.VictimCaseForm2ActionPlan
	currentEnums = vCaseForm.VictimCaseForm2ActionPlan

	enumsToDelete, currentEnums = cleanEnums(enumsToDelete, currentEnums)

	// Se insertan los nuevos enums.
	if err = storeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, currentEnums, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}

	// Se eliminan los enums que quedaron pendientes de borrar.
	if err = removeRelVictimCaseForm2EnumsVictimCaseForm2(vCaseFormToUpdate, enumsToDelete, connData, dbClientConfig, dbServerConfig); err != nil {
		return err
	}
	//-----------------------------------------------------------------------------

	return nil
}

func storeRelVictimCaseForm2EnumsVictimCaseForm2(form salvia_daos.VictimCaseForm2DTO, enums []salvia_daos.VictimCaseForm2EnumsDTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	// Se insertan los nuevos enums.
	for _, e := range enums {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = form
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err := salvia_daos.SetRelVictimCaseForm2EnumsVictimCaseForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}
	return nil
}
func removeRelVictimCaseForm2EnumsVictimCaseForm2(form salvia_daos.VictimCaseForm2DTO, enums []salvia_daos.VictimCaseForm2EnumsDTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	for _, e := range enums {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimCaseForm2DTO

		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId", "RelVictimCaseForm2EnumsVictimCaseForm2FormId"},
			AttrsValue: []interface{}{e.VictimCaseForm2EnumsId, form.VictimCaseForm2Id},
		}

		rel.RelVictimCaseForm2EnumsVictimCaseForm2Form = form
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err := salvia_daos.RemoveRelVictimCaseForm2EnumsVictimCaseForm2(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}
	return nil
}

// Retorna el cálculo de riesgo numéricamente y en nivel de riesgo. Analiza si es el caso de riesgo en pareja o riesgo cuando no es pareja
func getRiskScore(vCaseForm2 salvia_daos.VictimCaseForm2DTO) (int64, int64) {
	var partnerValidated bool = vCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor.VictimCaseForm2EnumsCode == "pi" || vCaseForm2.VictimCaseForm2RelationshipWithPresumedAggressor.VictimCaseForm2EnumsCode == "ex"
	var riskScore int64 = 0
	var riskLevel int64

	//Comunes:
	if vCaseForm2.VictimCaseForm2AggressorViolencePhysicalIncrease.VictimCaseForm2EnumsCode == "y" {
		riskScore++
	}
	if vCaseForm2.VictimCaseForm2AggressorWeaponUsed.VictimCaseForm2EnumsCode == "y" {
		riskScore++
	}
	if vCaseForm2.VictimCaseForm2AggressorThreatKill.VictimCaseForm2EnumsCode == "y" {
		riskScore++
	}
	if vCaseForm2.VictimCaseForm2AggressorPursuesSpiesDestroys.VictimCaseForm2EnumsCode == "y" {
		riskScore++
	}
	if vCaseForm2.VictimCaseForm2AggressorCapableOfKilling.VictimCaseForm2EnumsCode == "y" {
		riskScore++
	}
	if vCaseForm2.VictimCaseForm2AggressorHasAccessToWeapons.VictimCaseForm2EnumsCode == "y" {
		riskScore++
	}

	if partnerValidated {

		if vCaseForm2.VictimCaseForm2PartnerUnemployed.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2PartnerOtherDenunciations.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorHasPenalBackground.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2StoppedSeekingHelp.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorForcedSex.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorAttemptedStrangulation.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorConsumesDrugs.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorIsAlcoholic.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2PartnerControls.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorHadHitInVulnerability.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2VictimHealthToBlackmail.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2PartnerThreatenedSuicide.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2PartnerThreatenedDamageMembers.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorLimitsContactSupportNetworks.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2ThreatenedRevealSexualOrientation.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2StillLivesWithAggressor.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}
		if vCaseForm2.VictimCaseForm2AggressorViolentlyJealous.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		//Calculamos el nivel
		if riskScore <= 4 {
			riskLevel = 1
		} else if riskScore <= 8 {
			riskLevel = 2
		} else if riskScore <= 15 {
			riskLevel = 3
		} else if riskScore <= 24 {
			riskLevel = 4
		}

	} else {
		//Cuando no es pareja

		if vCaseForm2.VictimCaseForm2AggressorTakenAdvantagePhysicalVulnerability.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorUnemployed.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorHasPenalBackground2.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorSexuallyHarassment.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorSexuallyHarassment2.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2ViolenceMotivatedByGender2.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorUseDrugs.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorIsAlcoholic2.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorControls.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorThreatenedDamageMembers.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2ThoughtsOfSelfHarm2.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorCommonSpaces.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorHierarchy.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		if vCaseForm2.VictimCaseForm2AggressorUsedPositionAuthority.VictimCaseForm2EnumsCode == "y" {
			riskScore++
		}

		//Calculamos el nivel
		if riskScore <= 2 {
			riskLevel = 1
		} else if riskScore <= 5 {
			riskLevel = 2
		} else if riskScore <= 8 {
			riskLevel = 3
		} else if riskScore <= 20 {
			riskLevel = 4
		}

	}

	return riskScore, riskLevel
}
