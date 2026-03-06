// Package salvia_ctrl contiene controladores relacionados con la gestión de alertas
// y su asociación con casos de víctimas en el sistema.
package salvia_ctrl

import (
	"net/http"

	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_daos "bitsflow/salvia/dao"
	security_daos "bitsflow/security/dao"
)

// SetAlert crea la relación entre un caso de víctima y una alerta.
// Recibe como parámetros el identificador del caso de víctima, el código de alerta, datos adicionales de la alerta,
// el módulo desde el cual se invoca, datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func SetAlert(victimCaseId uint64, alertCode string, alertData string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crean las estructuras DTO requeridas para la operación.
	var relAlertVictimCase salvia_daos.RelAlertVictimCaseDTO = salvia_daos.RelAlertVictimCaseDTO{}
	var victimCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var alert salvia_daos.AlertDTO = salvia_daos.AlertDTO{}

	// Se construye el criterio de búsqueda para obtener el caso de víctima.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseId"},
		AttrsValue: []interface{}{victimCaseId},
	}

	// Se consulta el caso de víctima utilizando el criterio definido.
	if err = salvia_daos.GetVictimCase(by, &victimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se redefine el criterio de búsqueda para obtener la alerta correspondiente al código proporcionado.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"AlertCode"},
		AttrsValue: []interface{}{alertCode},
	}

	// Se consulta la alerta utilizando el criterio definido.
	if err = salvia_daos.GetAlert(by, &alert, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen los valores por defecto para la relación entre la alerta y el caso de víctima.
	salvia_daos.SetRelAlertVictimCaseDefaults(&relAlertVictimCase, common_dao.SQL_INSERT)
	relAlertVictimCase.RelAlertVictimCase_Alert = alert.AlertId
	relAlertVictimCase.RelAlertVictimCase_VictimCase = victimCase.VictimCaseId
	relAlertVictimCase.RelAlertVictimCase_Data = alertData

	// Se inserta la relación en la base de datos.
	if err = salvia_daos.SetRelAlertVictimCase(&relAlertVictimCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// GetAlertsByVictimCaseICode obtiene las alertas asociadas a un caso de víctima a partir de su código único (ICode).
// Recibe el código del caso de víctima, el módulo, datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP, un mensaje (JSON con la respuesta o el error) y una lista de DTO de alerta.
func GetAlertsByVictimCaseICode(victimCaseIcode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.AlertDTO) {
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa la lista de DTOs para las alertas.
	var alerts []salvia_daos.AlertDTO = []salvia_daos.AlertDTO{}

	// Se obtienen las alertas asociadas al caso de víctima identificado por victimCaseIcode.
	alerts, err = salvia_daos.GetAlertsByVictimCaseICode(victimCaseIcode, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), alerts
	}
	// Se retorna el estado OK y el mensaje JSON con la lista de alertas.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(alerts), alerts
}

// GetAlertsByTownCode obtiene las alertas asociadas a una localidad a partir del código de la misma.
// Recibe el código de la localidad (townCode), el módulo, datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAlertsByTownCode(townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las alertas y la localidad.
	var alerts []salvia_daos.AlertDTO = []salvia_daos.AlertDTO{}
	var town security_daos.TownDTO = security_daos.TownDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"TownCode"},
		AttrsValue: []interface{}{townCode},
	}

	// Se consulta la localidad.
	if err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtienen las alertas asociadas al código de localidad.
	if alerts, err = salvia_daos.GetAlertsByTownCode(townCode, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// Se retorna el error en caso de fallo en la consulta.
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con la lista de alertas.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(alerts)
}

// GetAlertsByCaseOwner obtiene las alertas asociadas al propietario de un caso, identificado por su código único (ICode).
// Recibe el código del propietario, el módulo, datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAlertsByCaseOwner(ownerICode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa la lista de DTOs para las alertas.
	var alerts []salvia_daos.AlertDTO = []salvia_daos.AlertDTO{}

	// Se obtienen las alertas asociadas al propietario identificado por ownerICode.
	if alerts, err = salvia_daos.GetAlertsByUserICode(ownerICode, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// Se retorna el error en caso de fallo en la consulta.
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con la lista de alertas.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(alerts)
}

// GetAlertsInDanger obtiene las alertas que tienen un tipo de alerta definido como "d", indicativo de peligro.
// Recibe el módulo, datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAlertsInDanger(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa la lista de DTOs para las alertas.
	var alerts []salvia_daos.AlertDTO = []salvia_daos.AlertDTO{}

	// Se construye el criterio de búsqueda para filtrar alertas por tipo "d" (danger).
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"AlertType"},
		AttrsValue: []interface{}{"d"},
	}

	// Se obtienen las alertas que cumplen con el criterio definido.
	if alerts, err = salvia_daos.GetAlerts(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// Se retorna el error en caso de fallo en la consulta.
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con la lista de alertas.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(alerts)
}

// GetAlertByAll obtiene todas las alertas registradas en el sistema.
// Recibe el módulo, datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAlertByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todas las alertas registradas.
	alerts, err := salvia_daos.GetAllAlerts(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se asigna el mensaje de error en caso de fallo en la consulta.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, se asigna el mensaje con las alertas en formato JSON.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(alerts)
	}
	return resCode, resData
}
