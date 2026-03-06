// Package salvia_ctrl proporciona controladores para la gestión de propietarios de casos (CaseOwner)
// en el sistema Salvia, permitiendo obtener la información de un CaseOwner a partir de distintos criterios.
package salvia_ctrl

import (
	// Configuración común del proyecto.
	common_config "bitsflow/common/config"
	// Controladores comunes para criterios de búsqueda.
	common_controllers "bitsflow/common/controllers"
	// Acceso a datos común en la aplicación.
	common_dao "bitsflow/common/dao"
	// Módulo para la conexión y gestión de la base de datos.
	"bitsflow/common/db"
	// Utilidades generales de la aplicación.
	"bitsflow/common/utils"
	// Acceso a datos específicos de Salvia.
	salvia_daos "bitsflow/salvia/dao"
	// Constantes para los códigos de estado HTTP.
	"net/http"
)

// GetCaseOwnerByICode obtiene la información de un CaseOwner a partir de su identificador "ICode".
//
// Parámetros:
//   - id: Identificador único (ICode) del CaseOwner.
//   - module: Nombre del módulo que solicita la operación.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje en formato JSON (éxito o error).
func GetCaseOwnerByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Inicialización de la variable error.
	var err error = nil

	// Si no se ha asignado un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO (Data Transfer Object) para almacenar la información del CaseOwner.
	var caseOwner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	if id == "" {
		// Si el ID es vacío, se retorna un error de solicitud incorrecta utilizando un mensaje localizado.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}

	} else {
		// Se define el criterio de búsqueda utilizando el operador SQL_AND.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CaseOwnerICode"},
			AttrsValue: []interface{}{id},
		}

		// Se ejecuta la búsqueda del CaseOwner en la base de datos.
		err = salvia_daos.GetCaseOwner(by, &caseOwner, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, se retorna un estado interno del servidor junto al mensaje de error.
			return http.StatusInternalServerError, err.Error()
		}

		// Se retorna el estado OK y el mensaje JSON de éxito que contiene el CaseOwner.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(caseOwner)
	}
	// Retorno por defecto en caso de que no se cumplan las condiciones anteriores.
	return http.StatusInternalServerError, ""
}

// GetCaseOwnerByGeneralUser obtiene la información de un CaseOwner basado en el identificador
// asociado a un usuario general.
//
// Parámetros:
//   - id: Identificador del usuario general asociado al CaseOwner.
//   - module: Nombre del módulo que invoca la función.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje en formato JSON (éxito o error).
//   - El objeto CaseOwnerDTO resultante (vacío en caso de error).
func GetCaseOwnerByGeneralUser(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.CaseOwnerDTO) {
	var err error = nil

	// Si no se ha asignado un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para almacenar la información del CaseOwner.
	var caseOwner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}
	if id == "" {
		// Retorna un error de solicitud incorrecta si el ID es vacío.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, caseOwner
		}

	} else {
		// Configura el criterio de búsqueda utilizando el atributo "CaseOwnerGeneralUser".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CaseOwnerGeneralUser"},
			AttrsValue: []interface{}{id},
		}

		// Se realiza la consulta para obtener el CaseOwner.
		err = salvia_daos.GetCaseOwner(by, &caseOwner, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, se retorna un estado interno del servidor junto con el mensaje de error y el objeto vacío.
			return http.StatusInternalServerError, err.Error(), caseOwner
		}

		// Retorna el estado OK, el mensaje de éxito y el objeto CaseOwner obtenido.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(caseOwner), caseOwner
	}
	// Retorno por defecto en caso de error.
	return http.StatusInternalServerError, "", caseOwner
}

// GetCaseOwnerByAll obtiene todos los CaseOwners registrados en el sistema.
//
// Parámetros:
//   - module: Nombre del módulo que invoca la función.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje en formato JSON (éxito o error).
//   - Una lista de objetos CaseOwnerDTO.
func GetCaseOwnerByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.CaseOwnerDTO) {
	// Variables para almacenar el mensaje de respuesta y el código de estado.
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se ha asignado un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se recuperan todos los CaseOwners de la base de datos.
	caseOwners, err := salvia_daos.GetAllCaseOwners(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// En caso de error, se asigna el mensaje de error y se retorna un estado interno del servidor.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la operación es exitosa, se retorna un estado OK junto con el mensaje de éxito.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(caseOwners)
	}
	return resCode, resData, caseOwners
}

// GetCaseOwnersByActiveUsers obtiene los CaseOwners asociados a usuarios activos.
//
// Parámetros:
//   - module: Nombre del módulo que invoca la función.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje en formato JSON (éxito o error).
//   - Una lista de objetos CaseOwnerDTO correspondientes a usuarios activos.
func GetCaseOwnersByActiveUsers(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.CaseOwnerDTO) {
	// Variables para almacenar el mensaje de respuesta y el código de estado.
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se ha asignado un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se recuperan los CaseOwners asociados a usuarios activos desde la base de datos.
	caseOwners, err := salvia_daos.GetCaseOwnersByActiveUsers(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// En caso de error, se asigna el mensaje de error y se retorna un estado interno del servidor.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la operación es exitosa, se retorna un estado OK junto con el mensaje de éxito.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(caseOwners)
	}
	return resCode, resData, caseOwners
}
