// Package security_ctrl contiene funciones para controlar la seguridad
// y la gestión de departamentos a través de operaciones en la base de datos.
package security_ctrl

import (
	// Importación de configuraciones, controladores, acceso a datos y utilidades comunes.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"
	"net/http"
)

// GetDepartmentByICode obtiene un departamento a partir de su código identificador (ICode).
//
// Parámetros:
//   - id: Código del departamento que se desea buscar.
//   - module: Nombre del módulo (se utiliza para contextos o logs, aunque en este fragmento no se usa directamente).
//   - connData: Datos de conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código HTTP representado como entero.
//   - Mensaje en formato string, que puede ser un mensaje de éxito o de error.
//   - Un objeto DepartmentDTO que contiene la información del departamento.
func GetDepartmentByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.DepartmentDTO) {
	var err error = nil

	// Si la conexión no tiene identificador, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el objeto DTO que contendrá los datos del departamento.
	var department security_daos.DepartmentDTO = security_daos.DepartmentDTO{}

	// Verifica si el identificador proporcionado es vacío.
	if id == "" {
		// Obtiene un mensaje de error localizado en español y retorna un Bad Request.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, security_daos.DepartmentDTO{}
		}
	} else {
		// Define los criterios de búsqueda para el departamento basado en su código.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentICode"},
			AttrsValue: []interface{}{id},
		}

		// Llama a la función que obtiene el departamento de la base de datos.
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, retorna un Internal Server Error junto con el mensaje de error.
			return http.StatusInternalServerError, err.Error(), security_daos.DepartmentDTO{}
		}

		// Retorna un estado OK, un mensaje de éxito en formato JSON y el objeto DepartmentDTO.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(department), department
	}
	// Retorno por defecto en caso de no cumplir las condiciones anteriores.
	return http.StatusInternalServerError, "", security_daos.DepartmentDTO{}
}

func GetDepartmentByTownCode(townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.DepartmentDTO) {
	var err error = nil

	// Si la conexión no tiene identificador, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el objeto DTO que contendrá los datos del departamento.
	var department security_daos.DepartmentDTO = security_daos.DepartmentDTO{}
	var town security_daos.TownDTO = security_daos.TownDTO{}
	var city security_daos.CityDTO = security_daos.CityDTO{}

	// Verifica si el identificador proporcionado es vacío.
	if townCode == "" {
		// Obtiene un mensaje de error localizado en español y retorna un Bad Request.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, security_daos.DepartmentDTO{}
		}
	} else {
		//Se ubica el municipio
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownCode"},
			AttrsValue: []interface{}{townCode},
		}
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, retorna un Internal Server Error junto con el mensaje de error.
			return http.StatusInternalServerError, err.Error(), security_daos.DepartmentDTO{}
		}

		//Se ubica la ciudad
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"CityId"},
			AttrsValue: []interface{}{town.TownCity},
		}
		err = security_daos.GetCity(by, &city, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, retorna un Internal Server Error junto con el mensaje de error.
			return http.StatusInternalServerError, err.Error(), security_daos.DepartmentDTO{}
		}

		//Se ubica el departamento

		// Define los criterios de búsqueda para el departamento basado en su código.
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{city.CityDepartment},
		}

		// Llama a la función que obtiene el departamento de la base de datos.
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, retorna un Internal Server Error junto con el mensaje de error.
			return http.StatusInternalServerError, err.Error(), security_daos.DepartmentDTO{}
		}

		// Retorna un estado OK, un mensaje de éxito en formato JSON y el objeto DepartmentDTO.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(department), department
	}
	// Retorno por defecto en caso de no cumplir las condiciones anteriores.
	return http.StatusInternalServerError, "", security_daos.DepartmentDTO{}
}

// GetDepartmentById obtiene un departamento a partir de su identificador numérico.
//
// Parámetros:
//   - id: Identificador numérico (uint64) del departamento.
//   - module: Nombre del módulo (utilizado para contexto o logging).
//   - connData: Datos de conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código HTTP representado como entero.
//   - Mensaje en formato string, que puede ser un mensaje de éxito o de error.
//   - Un objeto DepartmentDTO con la información del departamento.
func GetDepartmentById(id uint64, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.DepartmentDTO) {
	var err error = nil

	// Si la conexión no tiene identificador, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa el objeto DTO para almacenar la información del departamento.
	var department security_daos.DepartmentDTO = security_daos.DepartmentDTO{}

	// Valida que el identificador no sea cero.
	if id == 0 {
		// Retorna un mensaje de error localizado en español con un código de Bad Request.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, security_daos.DepartmentDTO{}
		}
	} else {
		// Define los criterios de búsqueda utilizando el identificador numérico.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"DepartmentId"},
			AttrsValue: []interface{}{id},
		}

		// Realiza la consulta a la base de datos para obtener el departamento.
		err = security_daos.GetDepartment(by, &department, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, retorna un estado de error interno junto al mensaje del error.
			return http.StatusInternalServerError, err.Error(), security_daos.DepartmentDTO{}
		}

		// Si la consulta es exitosa, retorna el estado OK, un mensaje de éxito y el objeto DepartmentDTO.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(department), department
	}
	// Retorno por defecto en caso de error no manejado.
	return http.StatusInternalServerError, "", security_daos.DepartmentDTO{}
}

// GetDepartmentByAll obtiene todos los departamentos disponibles.
//
// Parámetros:
//   - module: Nombre del módulo utilizado para la consulta.
//   - connData: Datos de conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - Código HTTP representado como entero.
//   - Mensaje en formato string, que puede ser un mensaje de éxito (incluyendo los datos en formato JSON)
//     o un mensaje de error.
func GetDepartmentByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no tiene identificador, se libera la conexión una vez finalizada la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se invoca la función para obtener todos los departamentos.
	users, err := security_daos.GetAllDepartment(connData, &dbClientConfig, &dbServerConfig)

	// Verifica si ocurrió un error durante la consulta.
	if err != nil {
		// Asigna el mensaje de error y establece el código de error.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, asigna el mensaje de éxito con los datos obtenidos.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(users)
	}
	// Retorna el código y el mensaje correspondiente.
	return resCode, resData
}
