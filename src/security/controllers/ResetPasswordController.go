// Package security_ctrl contiene funciones relacionadas con el control de la seguridad,
// específicamente para la gestión del reinicio de contraseñas de los usuarios.
package security_ctrl

import (
	"net/http"

	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_config "bitsflow/security/config"
	security_daos "bitsflow/security/dao"
)

// SetResetPassword procesa la solicitud de reinicio de contraseña para un usuario.
// Recibe como entrada un string en formato JSON con los datos necesarios, la URL base de la
// aplicación para construir el enlace del correo (ver security_config.ResolveAppDomain),
// y la configuración de la conexión y de la base de datos. La función valida los datos de entrada,
// obtiene el usuario, verifica la existencia de solicitudes de reinicio previas, crea un nuevo registro
// de reinicio de contraseña y envía un email con el enlace para restablecer la contraseña.
// Retorna un código HTTP, un mensaje en formato JSON y el DTO de reinicio de contraseña generado.
func SetResetPassword(dataInput string, appDomain string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.ResetPasswordDTO) {

	// Inicialización de variable para capturar errores
	var err error = nil

	// Si la conexión no tiene un ID asignado, se libera la conexión al finalizar la función
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para recopilar errores encontrados durante la validación de datos
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO requeridos para el proceso
	var resetPassword security_daos.ResetPasswordDTO = security_daos.ResetPasswordDTO{}
	var usr security_daos.GeneralUserDTO = security_daos.GeneralUserDTO{}

	// Variable para almacenar el DTO en forma de mapa
	var dtoMap map[string]interface{} = nil

	// Definición de campos requeridos para la validación del DTO
	var checkFields map[string]bool = map[string]bool{"GeneralUserLogin": true}

	// Se obtiene el mapa DTO a partir del JSON de entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, security_daos.GeneralUserJSONName, common_config.Locale, collectedErrors)

	// Validación del DTO y verificación de la estructura esperada
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Valida que el JSON de entrada cumpla con las definiciones de campos y formato esperado
		utils.ValidateJSONInput(&usr, dtoMap, security_daos.GeneralUserJSONName, security_daos.GeneralUserFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------
	default:
		// Si la estructura del JSON no es la esperada, se registra un error global
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si se han encontrado errores en la validación, se retorna un error BadRequest con los mensajes correspondientes
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.ResetPasswordDTO{}
	}

	// Construcción del criterio de búsqueda para obtener el usuario basado en el campo "GeneralUserLogin"
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserLogin"},
		AttrsValue: []interface{}{usr.GeneralUserLogin},
	}

	// Se obtiene el usuario general utilizando el DAO correspondiente.
	// En caso de error, se retorna un BadRequest con el mensaje de error.
	if err = security_daos.GetGeneralUser(by, &usr, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusBadRequest, err.Error(), security_daos.ResetPasswordDTO{}
	}

	// Construcción del criterio de búsqueda para obtener los correos electrónicos asociados al perfil del usuario
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EMailGeneralUserProfile"},
		AttrsValue: []interface{}{usr.GeneralUserGeneralUserProfile.GeneralUserProfileId},
	}

	// Se obtienen los correos electrónicos del usuario
	var emailsDTO []security_daos.EMailDTO
	if emailsDTO, err = security_daos.GetEMails(by, connData, &dbClientConfig, &dbServerConfig); err != nil || len(emailsDTO) == 0 {
		// Si no se encuentran correos, se registra un error indicando que no se halló correo para el usuario
		utils.SetError(collectedErrors, security_daos.GeneralUserJSONName, utils.GetTag(&usr, "GeneralUserLogin", "json"), "security_general_user_reset_mail_not_found", "", security_config.Locale)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors), security_daos.ResetPasswordDTO{}
	}

	// Se extraen las direcciones de correo electrónico de los DTOs para formar una lista de emails
	var emailsStr []string
	for _, e := range emailsDTO {
		emailsStr = append(emailsStr, e.EMailData)
	}

	// Verifica si ya existe una solicitud de reinicio pendiente para el usuario
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"ResetPasswordGeneralUser"},
		AttrsValue: []interface{}{usr.GeneralUserId},
	}

	// Si se encuentra una solicitud de reinicio pendiente, se elimina para evitar duplicados
	if err = security_daos.GetResetPassword(by, &resetPassword, connData, &dbClientConfig, &dbServerConfig); err == nil {
		if err = security_daos.RemoveResetPasswordByICode(resetPassword.ResetPasswordICode, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusBadRequest, err.Error(), security_daos.ResetPasswordDTO{}
		}
	}

	// Se establecen los valores por defecto en el DTO de reinicio de contraseña y se asocia el usuario
	security_daos.SetResetPasswordDefaults(&resetPassword, common_dao.SQL_INSERT)
	resetPassword.ResetPasswordGeneralUser = usr

	// Se crea el registro de reinicio de contraseña en la base de datos.
	// Si ocurre algún error, se retorna el error correspondiente.
	if err = security_daos.SetResetPassword(&resetPassword, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusBadRequest, err.Error(), security_daos.ResetPasswordDTO{}
	}

	// Se configura el email que será enviado al usuario para reiniciar su contraseña, incluyendo los detalles del servidor y el mensaje.
	var mail utils.EMail = utils.EMail{
		FromEMail:    security_config.EMAIL_SERVER_FROM_USER,
		FromName:     "Salvia Bot",
		To:           emailsStr,
		ServerHost:   security_config.EMAIL_SERVER_HOST_PATH,
		ServerPort:   security_config.EMAIL_SERVER_HOST_PORT,
		Username:     security_config.EMAIL_SERVER_HOST_USERNAME,
		UserPassword: security_config.EMAIL_SERVER_HOST_PASS,
		Subject:      "Reinicio de contraseña",
		Body:         `<h1>Reinicio de contraseña</h1><p>Se ha solicitado un cambio de contraseña en su cuenta de salvia. SI usted ha solicitado el cambio, por favor haga click en el siguiente enlace para establecer una nueva contraseña: <a href="` + appDomain + `/seguridad/login/` + resetPassword.ResetPasswordICode + `">Nueva contraseña</a></p>`}

	// Se envía el correo electrónico. Si ocurre algún error durante el envío, se retorna el error correspondiente.
	if err = mail.SendEMail(); err != nil {
		return http.StatusBadRequest, err.Error(), security_daos.ResetPasswordDTO{}
	}

	// Como el proceso se completó con éxito, se retorna un status OK, junto con un mensaje JSON de éxito y el DTO de reinicio.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(resetPassword), resetPassword
}

// GetResetPasswordByICode busca un registro de reinicio de contraseña utilizando su código único (ICode).
// Recibe el código del reinicio, el módulo de invocación y la configuración de conexión y de base de datos.
// Devuelve un código HTTP y un mensaje en formato JSON que contiene los datos del reinicio o un mensaje de error.
func GetResetPasswordByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Inicialización de variable para capturar errores
	var err error = nil

	// Libera la conexión al finalizar si no se ha asignado un ID de conexión
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para almacenar los datos del reinicio de contraseña
	var resetPassword security_daos.ResetPasswordDTO = security_daos.ResetPasswordDTO{}

	// Si se proporciona un ID válido, se procede a buscar el registro correspondiente
	if id != "" {
		// Construcción del criterio de búsqueda utilizando el código único del reinicio
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"ResetPasswordICode"},
			AttrsValue: []interface{}{id},
		}

		// Se obtiene el registro de reinicio de contraseña; en caso de error se retorna un InternalServerError
		err = security_daos.GetResetPassword(by, &resetPassword, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		// Se retorna el resultado exitoso junto con el DTO del reinicio
		return http.StatusOK, utils.CommMsgGetJSONSuccess(resetPassword)
	}

	// Si el ID no es válido, se retorna un error basado en la configuración regional (locale)
	if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
		return http.StatusBadRequest, v
	}
	return http.StatusInternalServerError, ""
}

// RemoveResetPasswordByICode elimina un registro de reinicio de contraseña utilizando su código único (ICode).
// Recibe el código, el módulo de invocación y la configuración de conexión y de base de datos.
// Retorna un código HTTP y un mensaje en formato JSON que indica el resultado de la operación.
func RemoveResetPasswordByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Inicialización de variable para capturar errores
	var err error = nil

	// Libera la conexión si no se ha asignado un ID de conexión
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el reinicio de contraseña
	var resetPassword security_daos.ResetPasswordDTO = security_daos.ResetPasswordDTO{}

	// Se establecen los valores por defecto en el DTO de reinicio
	security_daos.SetResetPasswordDefaults(&resetPassword, common_dao.SQL_INSERT)

	// Construcción del criterio de búsqueda para localizar el registro de reinicio utilizando el código único
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"ResetPasswordICode"},
		AttrsValue: []interface{}{id},
	}

	// Se consulta el registro de reinicio de contraseña; si ocurre algún error se retorna el mismo
	if err = security_daos.GetResetPassword(by, &resetPassword, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se procede a eliminar el registro de reinicio. Si ocurre un error durante la eliminación, se retorna el error.
	if err = security_daos.RemoveResetPasswordByICode(resetPassword.ResetPasswordICode, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un status OK junto con un mensaje de éxito en formato JSON.
	return http.StatusOK, utils.CommMsgGetJSONSuccess("")
}
