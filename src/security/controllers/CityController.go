// Package security_ctrl contiene funciones de control relacionadas con la seguridad,
// en particular para la gestión y consulta de información de ciudades.
package security_ctrl

import (
	// Importa configuraciones y utilidades comunes de la aplicación.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////
// GetCityByCode
///////////////////////////////////////////////////////////////////////////////

// GetCityByCode consulta la información de una ciudad utilizando el código de ciudad.
// Recibe el código de la ciudad (id), el módulo que realiza la solicitud,
// y la configuración de conexión y bases de datos necesarias.
// Devuelve un código de estado HTTP, un mensaje (usualmente en formato JSON) y el DTO de la ciudad.
//
// Parámetros:
//   - id: código de la ciudad a buscar. Si es una cadena vacía se retorna un error.
//   - module: módulo o contexto de la solicitud.
//   - connData: datos de conexión utilizados para la consulta.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - int: código HTTP representando el estado de la operación.
//   - string: mensaje informativo o de error.
//   - security_daos.CityDTO: objeto que contiene la información de la ciudad.
func GetCityByCode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.CityDTO) {
	var err error = nil

	// Si no se ha asignado una conexión, se asegura de liberar la conexión al final.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se instancia el DTO para almacenar los datos de la ciudad.
	var city security_daos.CityDTO = security_daos.CityDTO{}

	// Validación del parámetro 'id'. Si está vacío se retorna un error.
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, security_daos.CityDTO{}
		}
	} else {
		// Se construye el criterio de búsqueda para filtrar por "CityCode".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityCode"},
			AttrsValue: []interface{}{id},
		}

		// Se invoca la función para obtener la ciudad desde la base de datos.
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error durante la consulta, se retorna el error.
			return http.StatusInternalServerError, err.Error(), security_daos.CityDTO{}
		}

		// Retorna el resultado exitoso junto con un mensaje de éxito.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(city), city
	}
	// Retorno por defecto en caso de error no especificado.
	return http.StatusInternalServerError, "", security_daos.CityDTO{}
}

///////////////////////////////////////////////////////////////////////////////
// GetCityById
///////////////////////////////////////////////////////////////////////////////

// GetCityById consulta la información de una ciudad utilizando su identificador único (ID).
// Recibe el ID de la ciudad, el módulo que realiza la solicitud,
// y la configuración de conexión y bases de datos necesarias.
// Devuelve un código de estado HTTP, un mensaje y el DTO de la ciudad.
//
// Parámetros:
//   - id: identificador numérico único de la ciudad. Un valor de 0 se considera inválido.
//   - module: módulo o contexto de la solicitud.
//   - connData: datos de conexión utilizados para la consulta.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - int: código HTTP representando el estado de la operación.
//   - string: mensaje informativo o de error.
//   - security_daos.CityDTO: objeto que contiene la información de la ciudad.
func GetCityById(id uint64, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.CityDTO) {
	var err error = nil

	// Asegura la liberación de la conexión si no está previamente establecida.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Instancia el DTO para la ciudad.
	var city security_daos.CityDTO = security_daos.CityDTO{}

	// Valida el parámetro 'id'. Un id igual a 0 se considera inválido.
	if id == 0 {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, security_daos.CityDTO{}
		}
	} else {
		// Construye el criterio de búsqueda utilizando "CityId".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{id},
		}

		// Consulta la ciudad en la base de datos.
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna un error en caso de fallo en la consulta.
			return http.StatusInternalServerError, err.Error(), security_daos.CityDTO{}
		}

		// Retorna el resultado exitoso junto con un mensaje en formato JSON.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(city), city
	}
	// Retorno por defecto en caso de error.
	return http.StatusInternalServerError, "", security_daos.CityDTO{}
}

///////////////////////////////////////////////////////////////////////////////
// GetCitiesByDeparment
///////////////////////////////////////////////////////////////////////////////

// GetCitiesByDeparment consulta las ciudades asociadas a un departamento específico.
// Recibe el identificador del departamento, el módulo y las configuraciones de conexión y base de datos.
// Devuelve un código de estado HTTP y un mensaje en formato JSON.
//
// Parámetros:
//   - deparmentId: identificador numérico del departamento. Un valor de 0 se considera inválido.
//   - module: módulo o contexto de la solicitud.
//   - connData: datos de conexión para la consulta.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - int: código HTTP representando el estado de la operación.
//   - string: mensaje informativo o de error en formato JSON.
func GetCitiesByDeparment(deparmentId uint64, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []security_daos.CityDTO) {
	var err error = nil

	// Libera la conexión si ésta no está establecida.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Crea el slice que almacenará los DTOs de las ciudades.
	var cities []security_daos.CityDTO = []security_daos.CityDTO{}

	// Valida que el identificador del departamento sea válido.
	if deparmentId == 0 {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, []security_daos.CityDTO{}
		}
	} else {
		// Define el criterio de búsqueda utilizando el atributo "CityDepartment".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_OR,
			EqualSign:  []string{},
			AttrsName:  []string{"CityDepartment"},
			AttrsValue: []interface{}{deparmentId},
		}

		// Consulta las ciudades asociadas al departamento.
		cities, err = security_daos.GetCities(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna el error en caso de fallo durante la consulta.
			return http.StatusInternalServerError, err.Error(), []security_daos.CityDTO{}
		}

		// Retorna un mensaje de éxito junto con los datos en formato JSON.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(cities), cities
	}
	// Retorno por defecto en caso de error.
	return http.StatusInternalServerError, "", []security_daos.CityDTO{}
}

func GetCitiesByTownCode(townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []security_daos.CityDTO) {
	var err error = nil

	// Libera la conexión si ésta no está establecida.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Crea el slice que almacenará los DTOs de las ciudades.
	var cities []security_daos.CityDTO = []security_daos.CityDTO{}
	var town security_daos.TownDTO = security_daos.TownDTO{}
	var city security_daos.CityDTO = security_daos.CityDTO{}

	// Valida que el identificador del departamento sea válido.
	if townCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, []security_daos.CityDTO{}
		}
		return http.StatusBadRequest, "", []security_daos.CityDTO{}
	}
	// Traemos el CP
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_OR,
		EqualSign:  []string{},
		AttrsName:  []string{"TownCode"},
		AttrsValue: []interface{}{townCode},
	}

	err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Retorna el error en caso de fallo durante la consulta.
		return http.StatusInternalServerError, err.Error(), []security_daos.CityDTO{}
	}

	// Traemos la ciudad
	by = common_controllers.By{
		Operator:   common_dao.SQL_OR,
		EqualSign:  []string{},
		AttrsName:  []string{"CityId"},
		AttrsValue: []interface{}{town.TownCity},
	}

	err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Retorna el error en caso de fallo durante la consulta.
		return http.StatusInternalServerError, err.Error(), []security_daos.CityDTO{}
	}

	// Consulta las ciudades asociadas al departamento.
	by = common_controllers.By{
		Operator:   common_dao.SQL_OR,
		EqualSign:  []string{},
		AttrsName:  []string{"CityDepartment"},
		AttrsValue: []interface{}{city.CityDepartment},
	}
	cities, err = security_daos.GetCities(by, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Retorna el error en caso de fallo durante la consulta.
		return http.StatusInternalServerError, err.Error(), []security_daos.CityDTO{}
	}

	// Retorna un mensaje de éxito junto con los datos en formato JSON.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(cities), cities

}

///////////////////////////////////////////////////////////////////////////////
// GetCitiesByCityCode
///////////////////////////////////////////////////////////////////////////////

// GetCitiesByCityCode consulta las ciudades utilizando el código de ciudad.
// Recibe el código de ciudad, el módulo y las configuraciones necesarias para la consulta.
// Devuelve un código de estado HTTP y un mensaje en formato JSON.
//
// Parámetros:
//   - cityCode: código de la ciudad a buscar. Una cadena vacía se considera inválida.
//   - module: módulo o contexto de la solicitud.
//   - connData: datos de conexión utilizados para la consulta.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - int: código HTTP representando el estado de la operación.
//   - string: mensaje informativo o de error en formato JSON.
func GetCitiesByCityCode(cityCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil

	// Libera la conexión si aún no se ha establecido.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Crea el slice para almacenar los DTOs de las ciudades.
	var cities []security_daos.CityDTO = []security_daos.CityDTO{}

	// Valida el parámetro 'cityCode'. Una cadena vacía retorna un error.
	if cityCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	} else {
		// Define el criterio de búsqueda para filtrar por "CityCode".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_OR,
			EqualSign:  []string{},
			AttrsName:  []string{"CityCode"},
			AttrsValue: []interface{}{cityCode},
		}

		// Consulta las ciudades que coinciden con el código.
		cities, err = security_daos.GetCities(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna el error en caso de fallo durante la consulta.
			return http.StatusInternalServerError, err.Error()
		}

		// Retorna un mensaje de éxito junto con los datos en formato JSON.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(cities)
	}
	// Retorno por defecto en caso de error.
	return http.StatusInternalServerError, ""
}

///////////////////////////////////////////////////////////////////////////////
// GetCityByAll
///////////////////////////////////////////////////////////////////////////////

// GetCityByAll obtiene todas las ciudades disponibles en el sistema.
// Recibe el módulo y las configuraciones necesarias para la consulta.
// Devuelve un código de estado HTTP y un mensaje en formato JSON que contiene la información de todas las ciudades.
//
// Parámetros:
//   - module: módulo o contexto de la solicitud.
//   - connData: datos de conexión utilizados para la consulta.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - int: código HTTP representando el estado de la operación.
//   - string: mensaje informativo o de error en formato JSON.
func GetCityByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Libera la conexión si aún no se ha establecido.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se consulta la información de todas las ciudades disponibles.
	users, err := security_daos.GetAllCities(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// Si ocurre un error, se prepara el mensaje de error.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la consulta es exitosa, se retorna un mensaje con el JSON de los datos.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(users)
	}
	return resCode, resData
}
