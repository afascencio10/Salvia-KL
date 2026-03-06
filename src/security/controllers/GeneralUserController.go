// Package security_ctrl contiene las funciones de control relacionadas con la seguridad de usuarios generales.
package security_ctrl

import (
	// Importación de paquetes comunes y específicos para la aplicación
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
	"strings"

	"github.com/dchest/captcha"
	"golang.org/x/crypto/bcrypt"
)

type GeneralUserRequest struct {
	User security_daos.GeneralUserDTO `json:"user"`
}

type GeneralUserProfileRequest struct {
	Profile security_daos.GeneralUserProfileDTO `json:"profile"`
}

////////////////////////////////////////////////////////////////////////////////
// LoginGeneralUser
////////////////////////////////////////////////////////////////////////////////

// LoginGeneralUser autentica a un usuario general a partir de la información de entrada.
// Valida el DTO, el captcha, verifica la existencia del usuario, la contraseña y sus roles.
// Retorna un código HTTP, un mensaje en formato JSON y el DTO del usuario.
func LoginGeneralUser(dataInput string, checkCaptcha bool, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.GeneralUserDTO) {

	// Variable para capturar errores
	var err error = nil

	// Si no se especifica una conexión previamente, se libera la conexión al salir
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crea el DTO del usuario general y se inicializan variables adicionales
	var userRequest GeneralUserRequest

	var currPass string

	var dtoMap map[string]interface{} = nil

	var roles []security_daos.RoleDTO

	// Se establecen valores por defecto en el DTO
	security_daos.SetGeneralUserDefaults(&userRequest.User, common_dao.SQL_INSERT)

	// Se especifican los campos que se deben validar obligatoriamente
	var checkFields map[string]bool = map[string]bool{"GeneralUserPassword": true, "GeneralUserLogin": true}

	// Se obtiene el mapa de datos (DTO) a partir de la cadena de entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, security_daos.GeneralUserJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &userRequest)

	// Verificamos que los datos tengan la estructura esperada y validamos el JSON
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Se valida la entrada JSON y se mapean los valores al DTO del usuario
		utils.ValidateJSONInput(&userRequest.User, dtoMap, security_daos.GeneralUserJSONName, security_daos.GeneralUserFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
	default:
		// Si la estructura del JSON es incorrecta, se agrega un error global
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Variable para almacenar el resultado de la validación del captcha
	var captchaPassed bool = false

	//Validamos que no hayan habido errores en el procesamiento de los datos de entrada
	if len(collectedErrors) == 0 {
		// Se verifica el captcha usando la librería correspondiente
		captchaPassed = !checkCaptcha || captcha.VerifyString(userRequest.User.GeneralUserCaptchaID, userRequest.User.GeneralUserCaptchaSolution)
		// Si el captcha no pasó la validación, se registra el error
		if !captchaPassed {
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserCaptchaSolution", "json"), "security_general_user_captcha_fail", "security_general_user_captcha_fail", security_config.Locale)
		}
	}

	if len(collectedErrors) > 0 {
		// Se crea una nueva imagen de captcha
		resetCaptcha(collectedErrors)

		// Se retorna un error indicando que el DTO no fue validado correctamente
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.GeneralUserDTO{}
	}

	// Se procede a realizar las verificaciones de existencia del usuario y contraseña

	// Se arma el criterio de búsqueda para validar los campos únicos (login y estado activo "e")
	var by []common_controllers.By = []common_controllers.By{
		{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserLogin", "GeneralUserStatus"},
			AttrsValue: []interface{}{userRequest.User.GeneralUserLogin, "e"},
		},
	}
	currPass = userRequest.User.GeneralUserPassword

	// Se consulta una versión simple del usuario para evitar ataques
	if err = security_daos.GetGeneralUserLogin(by[0], &userRequest.User, connData, &dbClientConfig, &dbServerConfig); err != nil {
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", "login_fail", "login_fail", security_config.Locale)
		resetCaptcha(collectedErrors)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.GeneralUserDTO{}
	}

	// Se verifica la contraseña comparando el hash almacenado con el password ingresado
	err = bcrypt.CompareHashAndPassword([]byte(userRequest.User.GeneralUserPassword), []byte(currPass))
	if err != nil {
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", "login_fail", "login_fail", security_config.Locale)
		resetCaptcha(collectedErrors)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.GeneralUserDTO{}
	}

	// Se consulta la versión completa del usuario, usando su ID
	by = []common_controllers.By{
		{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserId"},
			AttrsValue: []interface{}{userRequest.User.GeneralUserId},
		},
	}
	if err = security_daos.GetGeneralUser(by[0], &userRequest.User, connData, &dbClientConfig, &dbServerConfig); err != nil {
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", "login_fail", "login_fail", security_config.Locale)
		resetCaptcha(collectedErrors)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.GeneralUserDTO{}
	}

	// Se obtienen los roles asociados al usuario
	roles, err = security_daos.GetRolesByGeneralUser(&userRequest.User, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se retorna un error si no se pudieron obtener los roles
		resetCaptcha(collectedErrors)
		return http.StatusInternalServerError, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.GeneralUserDTO{}
	}
	userRequest.User.GeneralUserRoles = roles
	// Se llenan los códigos de los roles en el DTO del usuario
	for i := 0; i < len(roles); i++ {
		userRequest.User.GeneralUserRoleCodes = append(userRequest.User.GeneralUserRoleCodes, roles[i].RoleCode)
	}

	// Se retorna la respuesta exitosa con el usuario autenticado
	return http.StatusOK, utils.CommMsgGetJSONSuccess(security_daos.GeneralUserDTO{GeneralUserICode: userRequest.User.GeneralUserICode, GeneralUserLanguage: userRequest.User.GeneralUserLanguage}), userRequest.User
}

////////////////////////////////////////////////////////////////////////////////
// UpdateGeneralUserByReset
////////////////////////////////////////////////////////////////////////////////

// UpdateGeneralUserByReset actualiza la contraseña de un usuario mediante un token de reinicio.
// Valida el DTO de reinicio, compara las contraseñas, genera el hash y actualiza el usuario.
// Retorna un código HTTP y un mensaje en formato JSON.
func UpdateGeneralUserByReset(dataInput string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	var err error = nil

	// Se libera la conexión si no se ha especificado una previamente
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Inicialización del mapa de errores
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Creación de DTOs requeridos: reset de contraseña y usuario
	var resetPassword security_daos.ResetPasswordDTO = security_daos.ResetPasswordDTO{}
	var userRequest GeneralUserRequest
	var usrToUpdate security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}

	var dtoMap map[string]interface{} = nil

	// Se especifican los campos obligatorios para la validación
	var checkFields map[string]bool = map[string]bool{"GeneralUserPassword": true, "GeneralUserPasswordRepeat": true, "GeneralUserResetPasswordICode": true}

	// Se obtiene el DTO a partir de la cadena de entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, security_daos.GeneralUserJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &userRequest)

	// Se valida la estructura del JSON y se asignan los valores al DTO
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		utils.ValidateJSONInput(&userRequest.User, dtoMap, security_daos.GeneralUserJSONName, security_daos.GeneralUserFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		// Error si no hay token de reinicio
		if userRequest.User.GeneralUserResetPasswordICode == "" {
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserResetPasswordICode", "json"), "security_general_user_reset_token_not_found", "common_global_error", security_config.Locale)
		}

		// Se valida que ambas contraseñas coincidan y se verifica la fortaleza de la misma
		if userRequest.User.GeneralUserPassword != userRequest.User.GeneralUserPasswordRepeat {
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPassword", "json"), "security_general_user_password_do_not_match", "common_global_error", security_config.Locale)
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPasswordRepeat", "json"), "security_general_user_password_do_not_match", "common_global_error", security_config.Locale)
		} else if passErrs := utils.CheckPasswordStrength(userRequest.User.GeneralUserPassword); len(passErrs) > 0 {
			if len(passErrs) > 0 {
				utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPassword", "json"), passErrs[0], "common_global_error", common_config.Locale)
				utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPasswordRepeat", "json"), passErrs[0], "common_global_error", common_config.Locale)
			}
		}
		// Se genera un hash de la contraseña con costo 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.User.GeneralUserPassword), 10)
		if err != nil {
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPassword", "json"), "security_general_user_password_internal_error", "common_global_error", security_config.Locale)
		} else {
			userRequest.User.GeneralUserPassword = string(hashedPassword)
		}

		//--------------------------------------------------------------------

	default:
		// Error genérico si la estructura no coincide
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si hay errores de validación, se retorna el error sin proceder con la actualización
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se arma el criterio de búsqueda para obtener el reset de contraseña
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"ResetPasswordICode"},
		AttrsValue: []interface{}{userRequest.User.GeneralUserResetPasswordICode},
	}

	// Se consulta el token de reinicio de contraseña
	err = security_daos.GetResetPassword(by, &resetPassword, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtiene el usuario a actualizar a partir del token de reinicio
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserId"},
		AttrsValue: []interface{}{resetPassword.ResetPasswordGeneralUser.GeneralUserId},
	}

	// Se trae el usuario que se actualizará
	if err = security_daos.GetGeneralUser(by, &usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusBadRequest, err.Error()
	}

	// Se inicia una transacción para actualizar de manera atómica
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el token de reinicio ya utilizado
	if err = security_daos.RemoveResetPasswordByICode(resetPassword.ResetPasswordICode, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se actualiza la contraseña del usuario con el nuevo hash generado
	usrToUpdate.GeneralUserPassword = userRequest.User.GeneralUserPassword

	if err = security_daos.UpdateGeneralUserByICode(&usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}
	// Se confirma la transacción
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	// Se retorna la respuesta exitosa
	return http.StatusOK, utils.CommMsgGetJSONSuccess(resetPassword)
}

////////////////////////////////////////////////////////////////////////////////
// SetGeneralUser
////////////////////////////////////////////////////////////////////////////////

// SetGeneralUser crea un nuevo usuario general, su perfil y asocia roles, correos y teléfonos.
// Realiza la validación de datos, verificación de campos únicos y maneja transacciones.
// Retorna un código HTTP y un mensaje JSON indicando el resultado.
func SetGeneralUser(dataInput string, s utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	var err error = nil
	var isProfileEmpty bool

	// Se libera la conexión si es necesario
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos para el usuario, teléfonos, roles y perfil
	var userRequest GeneralUserRequest
	var profileRequest GeneralUserProfileRequest

	var roles []security_daos.RoleDTO

	var mail security_daos.EMailDTO
	var phone security_daos.PhoneNumberDTO
	var dtoMap map[string]interface{} = nil

	var entityRoleFound bool = false
	var departmentOperatorRoleFound bool = false

	// Se especifican los campos obligatorios para el usuario
	var checkFields map[string]bool = map[string]bool{"GeneralUserPassword": true, "GeneralUserPasswordRepeat": true, "GeneralUserLanguage": true, "GeneralUserLogin": true}

	// Se obtiene el DTO a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, security_daos.GeneralUserJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &userRequest)
	utils.JSONToStruct(dataInput, &profileRequest)

	// Se establecen valores por defecto para el usuario y su perfil
	security_daos.SetGeneralUserDefaults(&userRequest.User, common_dao.SQL_INSERT)
	security_daos.SetGeneralUserProfileDefaults(&profileRequest.Profile, common_dao.SQL_INSERT)

	// Validación de la estructura del JSON y asignación de datos al DTO del usuario
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		utils.ValidateJSONInput(&userRequest.User, dtoMap, security_daos.GeneralUserJSONName, security_daos.GeneralUserFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		// Se valida el perfil del usuario (campos obligatorios y opcionales)
		checkFields = map[string]bool{"GeneralUserProfileGender": true, "GeneralUserProfileDescription": false, "GeneralUserProfileLastNames": true, "GeneralUserProfileNames": true, "GeneralUserProfileEmail": true, "GeneralUserProfileDocType": true, "GeneralUserProfileDocNumber": true}

		isProfileEmpty = utils.ValidateJSONInput(&profileRequest.Profile, dtoMap, security_daos.GeneralUserProfileJSONName, security_daos.GeneralUserProfileFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

		// Se valida que ambas contraseñas coincidan y se verifica la fortaleza de la misma
		if userRequest.User.GeneralUserPassword != userRequest.User.GeneralUserPasswordRepeat {
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPassword", "json"), "security_general_user_password_do_not_match", "common_global_error", security_config.Locale)
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPasswordRepeat", "json"), "security_general_user_password_do_not_match", "common_global_error", security_config.Locale)
		} else if passErrs := utils.CheckPasswordStrength(userRequest.User.GeneralUserPassword); len(passErrs) > 0 {
			if len(passErrs) > 0 {
				utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPassword", "json"), passErrs[0], "common_global_error", common_config.Locale)
				utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPasswordRepeat", "json"), passErrs[0], "common_global_error", common_config.Locale)
			}
		}
		// Se genera un hash de la contraseña con costo 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.User.GeneralUserPassword), 10)
		if err != nil {
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserPassword", "json"), "security_general_user_password_internal_error", "common_global_error", security_config.Locale)
		} else {
			userRequest.User.GeneralUserPassword = string(hashedPassword)
		}

		// Manejo de roles: se extraen y validan los roles asignados al usuario
		if s.CurrentRole == "do" {
			//Significa que sólo podrá crear usuarios entidad, entonces se agrega el respectivo rol e ignora los roles que vienen
			var role security_daos.RoleDTO
			var by common_controllers.By = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RoleCode"},
				AttrsValue: []interface{}{"et"},
			}

			if err = security_daos.GetRole(by, &role, connData, &dbClientConfig, &dbServerConfig); err != nil {
				// Se retorna error si no se encuentra el rol
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&userRequest.User, by.AttrsName[0], "json"), "security_role_not_found", "common_global_error", security_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
			roles = append(roles, role)
			userRequest.User.GeneralUserRoleIds = []string{role.RoleICode}
		} else {
			if len(userRequest.User.GeneralUserRoleIds) == 0 {
				utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, "GeneralUserRoleIds", "json"), "security_role_not_found", "common_global_error", security_config.Locale)
			}
			for _, r := range userRequest.User.GeneralUserRoleIds {
				var dto security_daos.RoleDTO
				security_daos.SetRoleDefaults(&dto, common_dao.SQL_INSERT)
				dto.RoleICode = r
				roles = append(roles, dto)
			}
		}

		// Se procesa el perfil del usuario para correos y otros campos específicos
		for _, r := range roles {
			var by common_controllers.By = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RoleICode"},
				AttrsValue: []interface{}{r.RoleICode},
			}

			if err = security_daos.GetRole(by, &r, connData, &dbClientConfig, &dbServerConfig); err != nil {
				// Se retorna error si no se encuentra el rol
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&profileRequest.Profile, by.AttrsName[0], "json"), "security_role_not_found", "common_global_error", security_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
			//Revisamos si es un rol que requiere de ubicación geográfica
			if r.RoleCode == "et" {
				entityRoleFound = true
				break
			}
			if r.RoleCode == "do" {
				departmentOperatorRoleFound = true
				break
			}
		}
		if entityRoleFound || departmentOperatorRoleFound {
			// Se obtienen valores específicos para entidad y localidad

			//Campos exclusivos para entidad
			//TODO: Poner mensaje global que diga que no hay sedes y no se puede hacer la operación
			if entityRoleFound && profileRequest.Profile.GeneralUserProfileEntityBranchSelected == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&profileRequest.Profile, "GeneralUserProfileEntityBranchSelected", "json"), "security_general_user_profile_branch", "common_global_error", security_config.Locale)
			}

			if profileRequest.Profile.GeneralUserProfileTownCode == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&profileRequest.Profile, "GeneralUserProfileTownCode", "json"), "security_general_user_profile_town", "common_global_error", security_config.Locale)
			}

		}

		// Manejo del correo: se procesa la cadena de correos y se separan en DTOs individuales
		var mails []string = strings.Split(profileRequest.Profile.GeneralUserProfileEmailStr, ",")
		for _, m := range mails {
			if m == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "emailStr", "common_validation_field_required_error", "", common_config.Locale)
				break
			}
			var dto security_daos.EMailDTO
			security_daos.SetEMailDefaults(&dto, common_dao.SQL_INSERT)
			dto.EMailData = m
			profileRequest.Profile.GeneralUserProfileEmails = append(profileRequest.Profile.GeneralUserProfileEmails, dto)
		}

		if len(profileRequest.Profile.GeneralUserProfileEmails) == 0 {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "emailStr", "common_validation_field_required_error", "", common_config.Locale)
		}

		var phones []string = strings.Split(profileRequest.Profile.GeneralUserProfilePhoneNumberStr, ",")
		for _, p := range phones {
			if p == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "phoneNumberStr", "common_validation_field_required_error", "", common_config.Locale)
				break
			}
			var dto security_daos.PhoneNumberDTO
			security_daos.SetPhoneNumberDefaults(&dto, common_dao.SQL_INSERT)
			dto.PhoneNumberData = p
			profileRequest.Profile.GeneralUserProfilePhoneNumbers = append(profileRequest.Profile.GeneralUserProfilePhoneNumbers, dto)
		}

		if len(profileRequest.Profile.GeneralUserProfilePhoneNumbers) == 0 {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "phoneNumberStr", "common_validation_field_required_error", "", common_config.Locale)
		}
		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Verificación de campos únicos para el perfil

	// Se arma el criterio de búsqueda para validar la unicidad de documento y nombre completo
	var by []common_controllers.By = []common_controllers.By{

		{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileDocNumber", "GeneralUserProfileDocType"},
			AttrsValue: []interface{}{userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDocType, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDocNumber},
		},
		{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileNames", "GeneralUserProfileLastNames"},
			AttrsValue: []interface{}{userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileNames, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames},
		},
	}
	for i := 0; i < len(by); i++ {
		if err = security_daos.GetGeneralUserProfile(by[i], &security_daos.GeneralUserProfileDTO{}, connData, &dbClientConfig, &dbServerConfig); err == nil {
			// Se detecta duplicidad y se retorna error
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&userRequest.User.GeneralUserGeneralUserProfile, by[i].AttrsName[0], "json"), "security_general_user_unique", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	// Verificación de campos únicos para el usuario
	by = []common_controllers.By{
		{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserLogin"},
			AttrsValue: []interface{}{userRequest.User.GeneralUserLogin},
		},
	}
	for i := 0; i < len(by); i++ {
		if err = security_daos.GetGeneralUser(by[i], &security_daos.GeneralUserDTO{}, connData, &dbClientConfig, &dbServerConfig); err == nil {
			// Se detecta duplicidad en el login y se retorna error
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, by[i].AttrsName[0], "json"), "security_general_user_unique", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	/*
		Validación de campos únicos para el correo:
		Se arma el criterio de búsqueda iterando sobre cada correo ingresado.
	*/
	var attrs []string
	var vals []interface{}

	for _, m := range profileRequest.Profile.GeneralUserProfileEmails {
		attrs = append(attrs, "EMailData")
		vals = append(vals, m.EMailData)
	}
	by = []common_controllers.By{}
	if len(profileRequest.Profile.GeneralUserProfileEmails) > 0 {
		by = []common_controllers.By{
			{
				Operator:   common_dao.SQL_OR,
				AttrsName:  attrs,
				AttrsValue: vals,
			},
		}
	}

	for i := 0; i < len(by); i++ {
		if err = security_daos.GetEMail(by[i], &mail, connData, &dbClientConfig, &dbServerConfig); err == nil {
			// Error de duplicidad de correo

			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&profileRequest.Profile, "GeneralUserProfileEmailStr", "json"), "security_email_unique", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	/*
		Validación similar para teléfonos, aunque en este caso el arreglo puede estar vacío.
	*/
	attrs = []string{}
	vals = []interface{}{}
	for _, p := range profileRequest.Profile.GeneralUserProfilePhoneNumbers {
		attrs = append(attrs, "PhoneNumberData")
		vals = append(vals, p.PhoneNumberData)
	}
	by = []common_controllers.By{}
	if len(profileRequest.Profile.GeneralUserProfilePhoneNumbers) > 0 {
		by = []common_controllers.By{
			{
				Operator:   common_dao.SQL_OR,
				AttrsName:  attrs,
				AttrsValue: vals,
			},
		}
	}

	for i := 0; i < len(by); i++ {
		if err = security_daos.GetPhoneNumber(by[i], &phone, connData, &dbClientConfig, &dbServerConfig); err == nil {
			// Error de duplicidad en teléfono
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&profileRequest.Profile, "GeneralUserProfilePhoneNumberStr", "json"), "security_phone_unique", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	// Se inicia la transacción para crear el usuario y sus relaciones
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se procede a crear el perfil del usuario si no está vacío
	if !isProfileEmpty {
		if err = security_daos.SetGeneralUserProfile(&profileRequest.Profile, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		userRequest.User.GeneralUserGeneralUserProfile = profileRequest.Profile

		// Se crean los correos asociados al perfil
		for _, m := range profileRequest.Profile.GeneralUserProfileEmails {
			m.EMailGeneralUserProfile = profileRequest.Profile
			if err = security_daos.SetEMail(&m, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, err.Error()
			}
		}

		// Se crean los teléfonos asociados al perfil
		for _, p := range profileRequest.Profile.GeneralUserProfilePhoneNumbers {
			p.PhoneNumberGeneralUserProfile = profileRequest.Profile
			if err = security_daos.SetPhoneNumber(&p, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, err.Error()
			}
		}

		// Se crea el CaseOwner asociado a la entidad, si tiene rol de entidad (entityRoleFound)
		var cowner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
		salvia_daos.SetCaseOwnerDefaults(&cowner, common_dao.SQL_INSERT)
		if entityRoleFound {
			var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
			var by common_controllers.By = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"EntityBranchICode"},
				AttrsValue: []interface{}{profileRequest.Profile.GeneralUserProfileEntityBranchSelected},
			}

			if err = salvia_daos.GetEntityBranch(by, &branch, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&branch, by.AttrsName[0], "json"), "security_general_user_branch_not_found", "common_global_error", security_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
			cowner.EntityBranch = branch
		}

		cowner.CaseOwnerGeneralUser = userRequest.User.GeneralUserICode

		if err = salvia_daos.SetCaseOwner(&cowner, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

	}

	// Se crea el usuario en la base de datos
	if err = security_daos.SetGeneralUser(&userRequest.User, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se asocian los roles al usuario
	for _, r := range roles {
		by = []common_controllers.By{
			{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RoleICode"},
				AttrsValue: []interface{}{r.RoleICode},
			},
		}
		if err = security_daos.GetRole(by[0], &r, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, by[0].AttrsName[0], "json"), "security_role_not_found", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
		var rel security_daos.RelRoleGeneralUserDTO
		security_daos.SetRelRoleGeneralUserDefaults(&rel, common_dao.SQL_INSERT)
		rel.RelRoleGeneralUserGeneralUser = userRequest.User
		rel.RelRoleGeneralUserRole = r
		if err = security_daos.SetRelRoleGeneralUser(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	// Se confirma la transacción y se retorna la respuesta exitosa
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, utils.CommMsgGetJSONSuccess(userRequest.User)
}

////////////////////////////////////////////////////////////////////////////////
// GetGeneralUserByICode
////////////////////////////////////////////////////////////////////////////////

// GetGeneralUserByICode obtiene la información completa de un usuario general a partir de su código único (ICode).
// Además, obtiene los roles, correos y datos geográficos asociados al usuario.
// Retorna un código HTTP y un mensaje JSON.
func GetGeneralUserByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.GeneralUserDTO) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el usuario y se inicializan variables auxiliares
	var usr security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
	var roles []security_daos.RoleDTO
	var mails []security_daos.EMailDTO
	var phones []security_daos.PhoneNumberDTO

	if id != "" {
		// Se arma el criterio de búsqueda para obtener el usuario
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserICode"},
			AttrsValue: []interface{}{id},
		}

		err = security_daos.GetGeneralUser(by, &usr, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), security_daos.GeneralUserDTO{}
		}

		// Se obtienen los roles asociados al usuario
		roles, err = security_daos.GetRolesByGeneralUser(&usr, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), security_daos.GeneralUserDTO{}
		}
		usr.GeneralUserRoles = roles

		// Se llena el arreglo de RoleIds
		for i := 0; i < len(roles); i++ {
			usr.GeneralUserRoleIds = append(usr.GeneralUserRoleIds, roles[i].RoleICode)
		}

		// Se obtienen los correos asociados al perfil del usuario
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EMailGeneralUserProfile"},
			AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileId},
		}
		mails, err = security_daos.GetEMails(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), security_daos.GeneralUserDTO{}
		}
		usr.GeneralUserGeneralUserProfile.GeneralUserProfileEmails = mails
		// Se construye una cadena con todos los correos
		for i := 0; i < len(mails); i++ {
			if i > 0 {
				usr.GeneralUserGeneralUserProfile.GeneralUserProfileEmailStr += ","
			}
			usr.GeneralUserGeneralUserProfile.GeneralUserProfileEmailStr += mails[i].EMailData
		}

		// Se obtienen los teléfonos asociados al perfil del usuario
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"PhoneNumberGeneralUserProfile"},
			AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileId},
		}
		phones, err = security_daos.GetPhoneNumbers(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), security_daos.GeneralUserDTO{}
		}
		usr.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers = phones
		// Se construye una cadena con todos los correos
		for i := 0; i < len(phones); i++ {
			if i > 0 {
				usr.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumberStr += ","
			}
			usr.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumberStr += phones[i].PhoneNumberData
		}

		// Se obtiene el CaseOwner para determinar la entidad y datos geográficos
		var cowner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CaseOwnerGeneralUser"},
			AttrsValue: []interface{}{usr.GeneralUserICode},
		}

		if err = salvia_daos.GetCaseOwner(by, &cowner, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
		}

		var town security_daos.TownDTO
		if usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownId > 0 {

			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"TownId"},
				AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownId},
			}
			err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
			}

			var city security_daos.CityDTO
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"CityId"},
				AttrsValue: []interface{}{town.TownCity},
			}
			err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
			}

			var department security_daos.DepartmentDTO
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"DepartmentId"},
				AttrsValue: []interface{}{city.CityDepartment},
			}
			err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}

			}

			var cities []security_daos.CityDTO
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"CityDepartment"},
				AttrsValue: []interface{}{department.DepartmentId},
			}
			cities, err = security_daos.GetCities(by, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}

			}

			var towns []security_daos.TownDTO
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"TownCity"},
				AttrsValue: []interface{}{city.CityId},
			}
			towns, err = security_daos.GetTowns(by, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}

			}

			usr.GeneralUserGeneralUserProfile.GeneralUserProfileSelectedTowns = towns
			usr.GeneralUserGeneralUserProfile.GeneralUserProfileSelectedCities = cities
			usr.GeneralUserGeneralUserProfile.GeneralUserProfileCity = city
			usr.GeneralUserGeneralUserProfile.GeneralUserProfileDepartment = department
			usr.GeneralUserGeneralUserProfile.GeneralUserProfileTownCode = town.TownCode
		}

		if cowner.EntityBranch.EntityBranchId > 0 {
			var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}

			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"EntityBranchId"},
				AttrsValue: []interface{}{cowner.EntityBranch.EntityBranchId},
			}

			if err = salvia_daos.GetEntityBranch(by, &branch, connData, &dbClientConfig, &dbServerConfig); err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
			}

			var entity salvia_daos.EntityDTO = salvia_daos.EntityDTO{}

			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"EntityId"},
				AttrsValue: []interface{}{branch.EntityBranchEntity.EntityId},
			}

			if err = salvia_daos.GetEntity(by, &entity, connData, &dbClientConfig, &dbServerConfig); err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
			}

			usr.GeneralUserGeneralUserProfile.GeneralUserProfileEntity = entity

			var branches []salvia_daos.EntityBranchDTO

			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"EntityBranchTownCode"},
				AttrsValue: []interface{}{town.TownCode},
			}
			branches, err = salvia_daos.GetEntityBranches(by, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
			}

			var branchesI []interface{} = make([]interface{}, len(branches))
			for i, b := range branches {
				branchesI[i] = b
			}

			usr.GeneralUserGeneralUserProfile.GeneralUserProfileSelectedBranches = branchesI

		}

		return http.StatusOK, utils.CommMsgGetJSONSuccess(usr), usr

	}

	if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
		return http.StatusBadRequest, v, security_daos.GeneralUserDTO{}
	}
	return http.StatusInternalServerError, "", security_daos.GeneralUserDTO{}
}

////////////////////////////////////////////////////////////////////////////////
// GetDepartmentGeneralUserByICodeAndTownCode
////////////////////////////////////////////////////////////////////////////////

// GetDepartmentGeneralUserByICodeAndTownCode obtiene la información completa de un usuario general a partir de su código único (ICode) y verifica que pertenezca al departamento asociado al TownCode.
// Además, obtiene los roles, correos y datos geográficos asociados al usuario.
// Retorna un código HTTP y un mensaje JSON.
func GetDepartmentGeneralUserByICodeAndTownCode(id string, townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.GeneralUserDTO) {
	var errStr string
	var usr security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
	var code int
	//var err error = nil
	var department1 security_daos.DepartmentDTO
	var department2 security_daos.DepartmentDTO
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	//Se consulta el usuario a través del método GetGeneralUserByICode
	code, errStr, usr = GetGeneralUserByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

	if code != http.StatusOK {
		return code, errStr, usr
	}

	//Verificamos que haya código de centro poblado
	if townCode == "" || usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownCode == "" {
		return http.StatusBadRequest, security_config.Locale["sp"]["get_department_users_empty_town"], usr
	}
	//Ahora se verifica si pertenece al mismo departamento
	code, errStr, department1 = GetDepartmentByTownCode(townCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
	if code != http.StatusOK {
		return code, errStr, usr
	}

	code, errStr, department2 = GetDepartmentByTownCode(usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
	if code != http.StatusOK {
		return code, errStr, usr
	}

	if department1.DepartmentCode != department2.DepartmentCode {
		return http.StatusBadRequest, salvia_config.Locale["sp"]["permission_denied"], usr
	}
	/*
		var town security_daos.TownDTO
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownId"},
			AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownId},
		}
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, ""
		}

		var city security_daos.CityDTO
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town.TownCity},
		}
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, ""
		}

		var department security_daos.DepartmentDTO
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city.CityDepartment},
		}
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, ""

		}

		var cities []security_daos.CityDTO
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityDepartment"},
			AttrsValue: []interface{}{department.DepartmentId},
		}
		cities, err = security_daos.GetCities(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, ""

		}

		var towns []security_daos.TownDTO
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownCity"},
			AttrsValue: []interface{}{city.CityId},
		}
		towns, err = security_daos.GetTowns(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, ""

		}

		usr.GeneralUserGeneralUserProfile.GeneralUserProfileSelectedTowns = towns
		usr.GeneralUserGeneralUserProfile.GeneralUserProfileSelectedCities = cities
		usr.GeneralUserGeneralUserProfile.GeneralUserProfileCity = city
		usr.GeneralUserGeneralUserProfile.GeneralUserProfileDepartment = department
		usr.GeneralUserGeneralUserProfile.GeneralUserProfileTownCode = town.TownCode*/

	return http.StatusOK, utils.CommMsgGetJSONSuccess(usr), usr
}

////////////////////////////////////////////////////////////////////////////////
// UpdateGeneralUserByICode
////////////////////////////////////////////////////////////////////////////////

// UpdateGeneralUserByICode actualiza la información de un usuario general identificado por su código (ICode).
// Actualiza datos del usuario, su perfil, roles, correos y teléfonos, manejando eliminaciones y adiciones mediante transacción.
// Retorna un código HTTP y un mensaje JSON.
func UpdateGeneralUserByICode(dataInput string, id string, s utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Se verifica que el ID no sea vacío
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil

	// Se libera la conexión si es necesario
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Inicialización del mapa para errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Creación de DTOs para el usuario y sus relaciones a actualizar
	var userRequest GeneralUserRequest
	var usrToUpdate security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}

	// Variables para correos, teléfonos y roles actuales
	//var phones []security_daos.PhoneNumberDTO
	var roles []security_daos.RoleDTO

	// Variables para elementos que se eliminarán
	var mailsToDelete []security_daos.EMailDTO
	var phonesToDelete []security_daos.PhoneNumberDTO
	var rolesToDelete []security_daos.RoleDTO

	var dtoMap map[string]interface{} = nil

	var entityRoleFound bool = false
	var departmentOperatorRoleFound bool = false

	// Se obtiene el usuario actual a actualizar basado en su ICode
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserICode"},
		AttrsValue: []interface{}{id},
	}

	if err = security_daos.GetGeneralUser(by, &usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se traen los correos actuales asociados al perfil del usuario
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EMailGeneralUserProfile"},
		AttrsValue: []interface{}{usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileId},
	}
	if mailsToDelete, err = security_daos.GetEMails(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen los teléfonos actuales asociados al perfil del usuario
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"PhoneNumberGeneralUserProfile"},
		AttrsValue: []interface{}{usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileId},
	}

	if phonesToDelete, err = security_daos.GetPhoneNumbers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen los roles actuales asignados al usuario
	if rolesToDelete, err = security_daos.GetRolesByGeneralUser(&usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se cargan los valores por defecto para el usuario y su perfil en modo actualización
	security_daos.SetGeneralUserDefaults(&userRequest.User, common_dao.SQL_UPDATE)
	security_daos.SetGeneralUserProfileDefaults(&userRequest.User.GeneralUserGeneralUserProfile, common_dao.SQL_UPDATE)

	// Se valida el DTO de entrada; se ignora la validación completa para permitir datos parciales
	var checkFields map[string]bool = map[string]bool{"GeneralUserLanguage": true}

	var rawDto interface{} = utils.GetDTOMap(dataInput, security_daos.GeneralUserJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &userRequest)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		utils.ValidateJSONInput(&userRequest.User, dtoMap, security_daos.GeneralUserJSONName, security_daos.GeneralUserFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{"GeneralUserProfileGender": true, "GeneralUserProfileDescription": false, "GeneralUserProfileLastNames": true, "GeneralUserProfileNames": true, "GeneralUserProfileEmail": true, "GeneralUserProfileDocType": true, "GeneralUserProfileDocNumber": true}
		utils.ValidateJSONInput(&userRequest.User.GeneralUserGeneralUserProfile, dtoMap, security_daos.GeneralUserJSONName+"."+security_daos.GeneralUserProfileJSONName, security_daos.GeneralUserProfileFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, true)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}

		// Manejo de roles: se extraen y validan los roles asignados al usuario
		if s.CurrentRole == "do" {
			//Significa que sólo podrá crear usuarios entidad, entonces se agrega el respectivo rol e ignora los roles que vienen
			var role security_daos.RoleDTO
			var by common_controllers.By = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RoleCode"},
				AttrsValue: []interface{}{"et"},
			}

			if err = security_daos.GetRole(by, &role, connData, &dbClientConfig, &dbServerConfig); err != nil {
				// Se retorna error si no se encuentra el rol
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&userRequest.User, by.AttrsName[0], "json"), "security_role_not_found", "common_global_error", security_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
			roles = append(roles, role)
			userRequest.User.GeneralUserRoleIds = []string{role.RoleICode}
		} else {
			for _, r := range userRequest.User.GeneralUserRoleIds {
				var dto security_daos.RoleDTO
				security_daos.SetRoleDefaults(&dto, common_dao.SQL_INSERT)
				dto.RoleICode = r
				roles = append(roles, dto)
			}
		}

		// Se procesa el perfil del usuario para correos y otros campos específicos
		for _, r := range roles {
			var by common_controllers.By = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RoleICode"},
				AttrsValue: []interface{}{r.RoleICode},
			}

			if err = security_daos.GetRole(by, &r, connData, &dbClientConfig, &dbServerConfig); err != nil {
				// Se retorna error si no se encuentra el rol
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&userRequest.User.GeneralUserGeneralUserProfile, by.AttrsName[0], "json"), "security_role_not_found", "common_global_error", security_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
			//Revisamos si es un rol que requiere de ubicación geográfica
			if r.RoleCode == "et" {
				entityRoleFound = true
				break
			}
			if r.RoleCode == "do" {
				departmentOperatorRoleFound = true
				break
			}
		}
		if entityRoleFound || departmentOperatorRoleFound {
			// Se obtienen valores específicos para entidad y localidad

			//Campos exclusivos para entidad
			if entityRoleFound && userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEntityBranchSelected == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&userRequest.User.GeneralUserGeneralUserProfile, "GeneralUserProfileEntityBranchSelected", "json"), "security_general_user_profile_branch", "common_global_error", security_config.Locale)
			}

			if userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileTownCode == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&userRequest.User.GeneralUserGeneralUserProfile, "GeneralUserProfileTownCode", "json"), "security_general_user_profile_town", "common_global_error", security_config.Locale)
			}

		}

		// Manejo del correo: se procesa la cadena de correos y se separan en DTOs individuales
		var mails []string = strings.Split(userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmailStr, ",")
		userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails = []security_daos.EMailDTO{}
		for _, m := range mails {
			if m == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "emailStr", "common_validation_field_required_error", "", common_config.Locale)
				break
			}
			var dto security_daos.EMailDTO
			security_daos.SetEMailDefaults(&dto, common_dao.SQL_INSERT)
			dto.EMailData = m
			userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails = append(userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails, dto)
		}

		if len(userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails) == 0 {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "emailStr", "common_validation_field_required_error", "", common_config.Locale)
		}

		var phones []string = strings.Split(userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumberStr, ",")
		userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers = []security_daos.PhoneNumberDTO{}
		for _, p := range phones {
			if p == "" {
				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "phoneNumberStr", "common_validation_field_required_error", "", common_config.Locale)
				break
			}
			var dto security_daos.PhoneNumberDTO
			security_daos.SetPhoneNumberDefaults(&dto, common_dao.SQL_INSERT)
			dto.PhoneNumberData = p
			userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers = append(userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers, dto)
		}

		if len(userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers) == 0 {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "phoneNumberStr", "common_validation_field_required_error", "", common_config.Locale)
		}

	default:
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si existen errores, se retorna el mensaje de error sin proceder con la actualización
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Verificación de campos únicos para el perfil
	var verifBy []common_controllers.By = []common_controllers.By{

		{
			Operator:      common_dao.SQL_AND,
			AttrsName:     []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
			AttrsValue:    []interface{}{userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDocType, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDocNumber},
			AttrsOldValue: []interface{}{usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileDocType, usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileDocNumber},
		},
		{
			Operator:      common_dao.SQL_AND,
			AttrsName:     []string{"GeneralUserProfileNames", "GeneralUserProfileLastNames"},
			AttrsValue:    []interface{}{userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileNames, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames},
			AttrsOldValue: []interface{}{usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileNames, usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames},
		},
	}

	for i := 0; i < len(verifBy); i++ {
		if !utils.CompareInterfaceSlices(verifBy[i].AttrsValue, verifBy[i].AttrsOldValue) {
			if err = security_daos.GetGeneralUserProfile(verifBy[i], &usrToUpdate.GeneralUserGeneralUserProfile, connData, &dbClientConfig, &dbServerConfig); err == nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)

				utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&usrToUpdate.GeneralUserGeneralUserProfile, verifBy[i].AttrsName[0], "json"), "security_general_user_unique", "common_global_error", security_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
		}
	}

	// Actualización de los campos del usuario
	usrToUpdate.GeneralUserLanguage = userRequest.User.GeneralUserLanguage

	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileGender = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileGender
	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileDescription = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDescription
	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames
	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileNames = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileNames
	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileDocType = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDocType
	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileDocNumber = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileDocNumber
	usrToUpdate.GeneralUserGeneralUserProfile.GeneralUserProfileTownCode = userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileTownCode

	// Se determinan los correos y teléfonos a eliminar y a crear mediante funciones de limpieza
	mailsToDelete, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails = cleanEMails(mailsToDelete, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails)
	phonesToDelete, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers = cleanPhones(phonesToDelete, userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers)
	rolesToDelete, roles = cleanRoles(rolesToDelete, roles)

	// Se inicia la transacción para la actualización
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Manejo del CaseOwner y actualización de entidad si corresponde
	var cowner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{usrToUpdate.GeneralUserICode},
	}

	if err = salvia_daos.GetCaseOwner(by, &cowner, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	salvia_daos.SetCaseOwnerDefaults(&cowner, common_dao.SQL_UPDATE)
	var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
	if entityRoleFound {

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityBranchICode"},
			AttrsValue: []interface{}{userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEntityBranchSelected},
		}

		if err = salvia_daos.GetEntityBranch(by, &branch, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(&branch, by.AttrsName[0], "json"), "security_general_user_branch_not_found", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}

	}

	cowner.EntityBranch = branch

	if err = salvia_daos.UpdateCaseOwner(&cowner, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Eliminación de teléfonos, correos y roles que ya no corresponden
	for i := 0; i < len(phonesToDelete); i++ {
		if err = security_daos.RemovePhoneNumberByICode(&phonesToDelete[i], connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}
	for i := 0; i < len(mailsToDelete); i++ {
		if err = security_daos.RemoveEMailByICode(&mailsToDelete[i], connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	for i := 0; i < len(rolesToDelete); i++ {
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"RelRoleGeneralUserRole", "RelRoleGeneralUserGeneralUser"},
			AttrsValue: []interface{}{rolesToDelete[i].RoleId, usrToUpdate.GeneralUserId},
		}
		if err = security_daos.RemoveRelRoleGeneralUsers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	// Se agregan los nuevos teléfonos, correos y roles
	for _, p := range userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfilePhoneNumbers {
		p.PhoneNumberGeneralUserProfile = usrToUpdate.GeneralUserGeneralUserProfile
		if err = security_daos.SetPhoneNumber(&p, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	for _, m := range userRequest.User.GeneralUserGeneralUserProfile.GeneralUserProfileEmails {
		m.EMailGeneralUserProfile = usrToUpdate.GeneralUserGeneralUserProfile
		if err = security_daos.SetEMail(&m, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	for _, r := range roles {
		var by = []common_controllers.By{
			{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RoleICode"},
				AttrsValue: []interface{}{r.RoleICode},
			},
		}
		if err = security_daos.GetRole(by[0], &r, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&userRequest.User, by[0].AttrsName[0], "json"), "security_role_not_found", "common_global_error", security_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}

		var rel security_daos.RelRoleGeneralUserDTO
		security_daos.SetRelRoleGeneralUserDefaults(&rel, common_dao.SQL_INSERT)
		rel.RelRoleGeneralUserGeneralUser = usrToUpdate
		rel.RelRoleGeneralUserRole = r
		if err = security_daos.SetRelRoleGeneralUser(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

	}

	// Se actualiza el usuario y su perfil en la base de datos
	if err = security_daos.UpdateGeneralUserByICode(&usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if err = security_daos.UpdateGeneralUserProfile(&usrToUpdate.GeneralUserGeneralUserProfile, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess(userRequest.User)
}

////////////////////////////////////////////////////////////////////////////////
// DisableGeneralUserByICode
////////////////////////////////////////////////////////////////////////////////

// DisableGeneralUserByICode deshabilita un usuario general marcándolo con el estado "d" (deshabilitado).
// Retorna un código HTTP y, en caso de éxito, una cadena vacía.
func DisableGeneralUserByICode(dataInput string, id string, s utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtiene el usuario a actualizar
	var usrToUpdate security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
	var code int

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserICode"},
		AttrsValue: []interface{}{id},
	}

	if s.CurrentRole == "ad" {
		if err = security_daos.GetGeneralUser(by, &usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, err.Error()
		}
	} else {
		if code, _, usrToUpdate = GetDepartmentGeneralUserByICodeAndTownCode(id, s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig); code != http.StatusOK {
			return http.StatusInternalServerError, salvia_config.Locale["sp"]["permission_denied"]
		}
	}

	// Se actualiza el estado del usuario a "d" (deshabilitado)
	usrToUpdate.GeneralUserStatus = "d"

	if err = security_daos.UpdateGeneralUserByICode(&usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, ""
}

////////////////////////////////////////////////////////////////////////////////
// EnableGeneralUserByICode
////////////////////////////////////////////////////////////////////////////////

// EnableGeneralUserByICode habilita un usuario general marcándolo con el estado "e" (habilitado).
// Retorna un código HTTP y, en caso de éxito, una cadena vacía.
func EnableGeneralUserByICode(dataInput string, id string, s utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var usrToUpdate security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
	var code int

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserICode"},
		AttrsValue: []interface{}{id},
	}
	if s.CurrentRole == "ad" {
		if err = security_daos.GetGeneralUser(by, &usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, err.Error()
		}
	} else {
		if code, _, usrToUpdate = GetDepartmentGeneralUserByICodeAndTownCode(id, s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig); code != http.StatusOK {
			return http.StatusInternalServerError, salvia_config.Locale["sp"]["permission_denied"]
		}
	}

	// Se actualiza el estado del usuario a "e" (habilitado)
	usrToUpdate.GeneralUserStatus = "e"

	if err = security_daos.UpdateGeneralUserByICode(&usrToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, ""
}

////////////////////////////////////////////////////////////////////////////////
// RemoveGeneralUserByICode
////////////////////////////////////////////////////////////////////////////////

// RemoveGeneralUserByICode elimina completamente un usuario general, su perfil, correos, teléfonos y relaciones de roles.
// Se ejecuta dentro de una transacción para asegurar la integridad de la operación.
// Retorna un código HTTP y un mensaje JSON.
func RemoveGeneralUserByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el usuario y se establecen valores por defecto
	var usr security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}
	security_daos.SetGeneralUserDefaults(&usr, common_dao.SQL_INSERT)

	// Se inicia una transacción
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se consulta el usuario a eliminar
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserICode"},
		AttrsValue: []interface{}{id},
	}

	if err = security_daos.GetGeneralUser(by, &usr, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se eliminan correos asociados al perfil del usuario
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EMailGeneralUserProfile"},
		AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileId},
	}

	if err = security_daos.RemoveEMails(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se eliminan teléfonos asociados al perfil
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"PhoneNumberGeneralUserProfile"},
		AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileId},
	}
	if err = security_daos.RemovePhoneNumbers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se eliminan relaciones de roles del usuario
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelRoleGeneralUserGeneralUser"},
		AttrsValue: []interface{}{usr.GeneralUserId},
	}
	if err = security_daos.RemoveRelRoleGeneralUsers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el usuario
	if err = security_daos.RemoveGeneralUserByICode(&usr, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el perfil, si existe
	if usr.GeneralUserGeneralUserProfile.GeneralUserProfileId > 0 {
		if err = security_daos.RemoveGeneralUserProfile(&usr.GeneralUserGeneralUserProfile, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, err.Error()
		}
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	} else {
		return http.StatusOK, utils.CommMsgGetJSONSuccess("")
	}
}

////////////////////////////////////////////////////////////////////////////////
// GetGeneralUserByAll
////////////////////////////////////////////////////////////////////////////////

// GetGeneralUserByAll obtiene la lista de todos los usuarios generales junto con sus perfiles.
// Retorna un código HTTP y un mensaje JSON con los datos o el error correspondiente.
func GetGeneralUserByAll(page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	users, count, err := security_daos.GetAllGeneralUserWithProfile(page, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{users, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"})
	}
	return resCode, resData, count
}

func GetGeneralUserByStatus(status string, page int, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	users, count, err := security_daos.GetGeneralUserWithProfileByStatus(status, page, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{users, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"})
	}
	return resCode, resData, count
}

////////////////////////////////////////////////////////////////////////////////
// GetGeneralUsersByRole
////////////////////////////////////////////////////////////////////////////////

// GetGeneralUsersByRole obtiene la lista de usuarios generales asociados a un rol específico.
// Retorna un código HTTP y un mensaje JSON con la información de los usuarios o un error.
func GetGeneralUsersByRole(roleCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	users, err := security_daos.GetGeneralUsersByRoleWithProfile(roleCode, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(users)
	}
	return resCode, resData
}

////////////////////////////////////////////////////////////////////////////////
// GetGeneralUsersByRole
////////////////////////////////////////////////////////////////////////////////

// GetGeneralUsersByRole obtiene la lista de usuarios generales asociados a un rol específico.
// Retorna un código HTTP y un mensaje JSON con la información de los usuarios o un error.
func GetDepartmentGeneralUsersByTownCodeAndRoleCode(roleCode string, townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	code, errStr, department := GetDepartmentByTownCode(townCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

	if code != http.StatusOK {
		resData = errStr
		resCode = http.StatusInternalServerError
		return resCode, resData
	}

	users, err := security_daos.GetGeneralUsersByDepartmentICodeAndRoleCode(roleCode, department.DepartmentICode, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(users)
	}
	return resCode, resData
}

////////////////////////////////////////////////////////////////////////////////
// Funciones de limpieza de slices (EMails, Phones, Roles)
////////////////////////////////////////////////////////////////////////////////

// cleanEMails compara dos listas de correos y retorna dos slices: uno con los correos a eliminar y otro con los correos nuevos a insertar.
func cleanEMails(mailsToDelete []security_daos.EMailDTO, mails []security_daos.EMailDTO) ([]security_daos.EMailDTO, []security_daos.EMailDTO) {
	for i := 0; i < len(mailsToDelete); i++ {
		var idx int = getMailIn(mailsToDelete[i], mails)
		// Si el correo se encuentra en ambos slices, se elimina de ambos
		if idx != -1 {
			mailsToDelete = removeMailFromSlice(mailsToDelete, i)
			mails = removeMailFromSlice(mails, idx)
		}
	}
	return mailsToDelete, mails
}

// getMailIn busca un correo en una lista y retorna el índice donde se encuentra; si no se encuentra, retorna -1.
func getMailIn(mailToFind security_daos.EMailDTO, mails []security_daos.EMailDTO) int {
	for i := 0; i < len(mails); i++ {
		if mailToFind.EMailData == mails[i].EMailData {
			return i
		}
	}
	return -1
}

// cleanPhones compara dos listas de teléfonos y retorna dos slices: uno con los teléfonos a eliminar y otro con los nuevos.
func cleanPhones(phonesToDelete []security_daos.PhoneNumberDTO, phones []security_daos.PhoneNumberDTO) ([]security_daos.PhoneNumberDTO, []security_daos.PhoneNumberDTO) {
	for i := 0; i < len(phonesToDelete); i++ {
		var idx int = getPhoneIn(phonesToDelete[i], phones)
		if idx != -1 {
			phonesToDelete = removePhoneNumberFromSlice(phonesToDelete, i)
			phones = removePhoneNumberFromSlice(phones, idx)
		}
	}
	return phonesToDelete, phones
}

// getPhoneIn busca un teléfono en una lista y retorna su índice; si no se encuentra, retorna -1.
func getPhoneIn(phoneToFind security_daos.PhoneNumberDTO, phones []security_daos.PhoneNumberDTO) int {
	for i := 0; i < len(phones); i++ {
		if phoneToFind.PhoneNumberData == phones[i].PhoneNumberData {
			return i
		}
	}
	return -1
}

// cleanRoles compara dos listas de roles y retorna dos slices: uno con los roles a eliminar y otro con los nuevos.
func cleanRoles(rolesToDelete []security_daos.RoleDTO, roles []security_daos.RoleDTO) ([]security_daos.RoleDTO, []security_daos.RoleDTO) {
	for i := 0; i < len(rolesToDelete); i++ {
		var idx int = getRoleIn(rolesToDelete[i], roles)
		if idx != -1 {
			rolesToDelete = removeRoleFromSlice(rolesToDelete, i)
			roles = removeRoleFromSlice(roles, idx)
		}
	}
	return rolesToDelete, roles
}

// getRoleIn busca un rol en una lista y retorna el índice donde se encuentra; de lo contrario, retorna -1.
func getRoleIn(roleToFind security_daos.RoleDTO, roles []security_daos.RoleDTO) int {
	for i := 0; i < len(roles); i++ {
		if roleToFind.RoleICode == roles[i].RoleICode {
			return i
		}
	}
	return -1
}

// removeMailFromSlice elimina un elemento de un slice de correos en el índice especificado.
func removeMailFromSlice(slice []security_daos.EMailDTO, idx int) []security_daos.EMailDTO {
	return append(slice[:idx], slice[idx+1:]...)
}

// removePhoneNumberFromSlice elimina un elemento de un slice de teléfonos en el índice indicado.
func removePhoneNumberFromSlice(slice []security_daos.PhoneNumberDTO, idx int) []security_daos.PhoneNumberDTO {
	return append(slice[:idx], slice[idx+1:]...)
}

// removeRoleFromSlice elimina un elemento de un slice de roles en el índice especificado.
func removeRoleFromSlice(slice []security_daos.RoleDTO, idx int) []security_daos.RoleDTO {
	return append(slice[:idx], slice[idx+1:]...)
}

////////////////////////////////////////////////////////////////////////////////
// resetCaptcha
////////////////////////////////////////////////////////////////////////////////

// resetCaptcha asigna un nuevo captcha al mapa de errores en la clave "default" con la llave "captchaID".
// Esto se utiliza para reiniciar el captcha en caso de error.
func resetCaptcha(collectedErrors map[string]map[string]string) {
	collectedErrors["default"]["captchaID"] = captcha.New()
}
