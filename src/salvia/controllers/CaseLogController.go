package salvia_ctrl

import (
	// Importaciones de configuración, controladores, acceso a datos, utilidades y manejo de base de datos.
	common_config "bitsflow/common/config"           // Configuración común de la aplicación.
	common_controllers "bitsflow/common/controllers" // Funciones y estructuras comunes para controladores.
	common_dao "bitsflow/common/dao"                 // Acceso a datos y estructuras comunes de DAO.
	"bitsflow/common/db"                             // Funciones y estructuras para la gestión de conexiones a la base de datos.
	"bitsflow/common/utils"                          // Utilidades generales, incluyendo validación de JSON y manejo de errores.
	security_config "bitsflow/salvia/config"         // Configuración de seguridad específica para Salvia.
	salvia_daos "bitsflow/salvia/dao"                // Acceso a datos y estructuras específicas de Salvia.
	"net/http"                                       // Constantes y utilidades para el manejo de HTTP.
)

// SetCaseLog crea un registro de log de caso basado en los datos de entrada proporcionados.
// Realiza la validación del JSON de entrada, establece valores por defecto, asocia el log a un "moment"
// (si se proporciona) y genera una alerta relacionada.
//
// Parámetros:
// - dataInput: Cadena JSON que contiene los datos del log del caso.
// - momentICode: Código identificador del momento asociado (opcional).
// - module: Identificador del módulo que invoca la función.
// - s: Sesión del usuario que incluye el código de usuario y otros datos de sesión.
// - dbClientConfig: Configuración del cliente de la base de datos.
// - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
// - Código de estado HTTP que indica el resultado de la operación.
// - Cadena JSON con los datos del log creado en caso de éxito o mensajes de error.
func SetCaseLog(dataInput string, momentICode string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable de error.
	var err error = nil

	// Se obtiene la conexión a la base de datos.
	var connData *db.ConnData = &db.ConnData{}
	// Se garantiza la liberación de la conexión al finalizar la función.
	defer db.ReleaseConnection(connData)

	// Mapa para recolectar errores durante el procesamiento.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO (Data Transfer Object) requeridos para el log de caso y el momento.
	var caseLog salvia_daos.CaseLogDTO = salvia_daos.CaseLogDTO{}
	var moment salvia_daos.MomentDTO = salvia_daos.MomentDTO{}

	var dtoMap map[string]interface{} = nil

	// Se establecen los valores por defecto para el DTO del log de caso (operación de inserción SQL).
	salvia_daos.SetCaseLogDefaults(&caseLog, common_dao.SQL_INSERT)

	// Definición de campos requeridos para la validación del JSON (por ejemplo, "CaseLogDescription").
	var checkFields map[string]bool = map[string]bool{"CaseLogDescription": true}

	// Se parsea la cadena JSON de entrada en una estructura de mapa.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.CaseLogJSONName, common_config.Locale, collectedErrors)

	// Validación de la estructura del DTO obtenido.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Se valida la entrada JSON contra las definiciones esperadas de los campos.
		utils.ValidateJSONInput(&caseLog, dtoMap, salvia_daos.CaseLogJSONName, salvia_daos.CaseLogFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------
	default:
		// Si la estructura no es la esperada, se registra un error global.
		utils.SetError(collectedErrors, salvia_daos.CaseLogJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si se han detectado errores durante la validación, se retorna un error 400 con los detalles en formato JSON.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Si se proporcionó un momentICode, se procede a buscar el momento asociado.
	var by common_controllers.By = common_controllers.By{}
	if momentICode != "" {
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"MomentICode"},
			AttrsValue: []interface{}{momentICode},
		}

		// Se obtiene el registro del momento desde la base de datos.
		err = salvia_daos.GetMoment(by, &moment, connData, &dbClientConfig, &dbServerConfig)
		if err == nil {
			// Se asigna el identificador del momento al DTO del log.
			caseLog.CaseLogMoment.MomentId = moment.MomentId
		}
	}

	// Se asocia el objeto completo del momento y el usuario al DTO del log.
	caseLog.CaseLogMoment = moment
	caseLog.CaseLogGeneralUser = s.UserICode

	// Se procede a crear el registro del log en la base de datos.
	if err = salvia_daos.SetCaseLog(&caseLog, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// En caso de error durante la creación, se retorna un error 500 con el mensaje correspondiente.
		return http.StatusInternalServerError, err.Error()
	}

	// Se adiciona la alerta respectiva relacionada con el log del momento.
	code, error := SetAlert(moment.MomentVictimCase.VictimCaseId, security_config.ALERT_TYPE["MOMENT_LOG"], "", connData, s, dbClientConfig, dbServerConfig)

	// Si la creación de la alerta falla, se registra el error internamente sin exponerlo al usuario.
	if code != 200 {
		// TODO: Registrar el error en el log interno sin escalarlo al usuario.
		println("Error: CaseLogController.go: 110 - ", error)
	}

	// Se retorna el estado 200 OK con el DTO del log de caso en formato JSON.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(caseLog)
}

// GetCaseLogByICode obtiene un registro de log de caso a partir de su identificador único (ICode).
//
// Parámetros:
// - id: Identificador único del log de caso.
// - module: Identificador del módulo que invoca la función.
// - connData: Datos de la conexión a la base de datos (se libera si no está previamente establecida).
// - dbClientConfig: Configuración del cliente de la base de datos.
// - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
// - Código de estado HTTP que indica el resultado de la operación.
// - Cadena JSON con los datos del log en caso de éxito o mensajes de error.
// - El objeto CaseLogDTO correspondiente al log de caso.
func GetCaseLogByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.CaseLogDTO) {
	var err error = nil

	// Si la conexión no está establecida, se libera al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa el DTO del log de caso.
	var caseLog salvia_daos.CaseLogDTO = salvia_daos.CaseLogDTO{}
	if id == "" {
		// Si no se proporciona el identificador, se retorna un error 400 con un mensaje localizado.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, caseLog
		}

	} else {
		// Se define la condición de búsqueda para el log de caso utilizando el identificador.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CaseLogICode"},
			AttrsValue: []interface{}{id},
		}

		// Se obtiene el log de caso desde la base de datos.
		err = salvia_daos.GetCaseLog(by, &caseLog, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, se retorna un error 500 con el mensaje correspondiente.
			return http.StatusInternalServerError, err.Error(), caseLog
		}

		// Se retorna el log de caso obtenido con un mensaje de éxito.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(caseLog), caseLog

	}
	// Retorno por defecto en caso de flujo inesperado.
	return http.StatusInternalServerError, "", caseLog
}

// GetCaseLogsByMomentICode obtiene todos los registros de log de caso asociados a un código de momento específico.
//
// Parámetros:
// - momentICode: Código identificador del momento cuyos logs se desean obtener.
// - module: Identificador del módulo que invoca la función.
// - connData: Datos de la conexión a la base de datos (se libera si no está previamente establecida).
// - dbClientConfig: Configuración del cliente de la base de datos.
// - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
// - Código de estado HTTP que indica el resultado de la operación.
// - Cadena JSON con la lista de logs en caso de éxito o mensajes de error.
func GetCaseLogsByMomentICode(momentICode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Se libera la conexión si no se ha establecido previamente.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Validación de que se haya proporcionado el código del momento.
	if momentICode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil

	// Se inicializa una lista vacía para almacenar los logs de caso.
	var caseLogs []salvia_daos.CaseLogDTO = []salvia_daos.CaseLogDTO{}

	// Se define la condición de búsqueda para obtener el momento utilizando su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentICode"},
		AttrsValue: []interface{}{momentICode},
	}

	// Se inicializa el DTO del momento.
	var moment salvia_daos.MomentDTO = salvia_daos.MomentDTO{}

	// Se obtiene el registro del momento desde la base de datos.
	err = salvia_daos.GetMoment(by, &moment, connData, &dbClientConfig, &dbServerConfig)
	if err == nil {
		// Se actualiza la condición de búsqueda para filtrar los logs asociados al identificador del momento.
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CaseLogMoment"},
			AttrsValue: []interface{}{moment.MomentId},
		}

		// Se obtienen los logs de caso asociados al momento.
		if caseLogs, err = salvia_daos.GetCaseLogs(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
			// En caso de error durante la consulta, se retorna un error 500.
			return http.StatusInternalServerError, err.Error()
		}

	}

	// Se retorna un mensaje de éxito con la lista de logs de caso en formato JSON.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(caseLogs)
}
