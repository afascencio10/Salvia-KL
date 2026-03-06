// Package salvia_ctrl contiene funciones de control para gestionar entidades en el módulo Salvia.
package salvia_ctrl

import (
	// Importación de configuraciones, controladores, acceso a datos, utilidades y manejo de base de datos comunes
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_daos "bitsflow/salvia/dao"
	"net/http"
)

// GetEntityByICode busca y retorna una entidad a partir de su código identificador (ICode).
//
// Parámetros:
//   - id: Código identificador de la entidad.
//   - module: Nombre del módulo que realiza la solicitud.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente para la base de datos.
//   - dbServerConfig: Configuración del servidor para la base de datos.
//
// Retorna:
//   - Código HTTP que indica el resultado de la operación.
//   - Mensaje en formato JSON describiendo el resultado.
//   - Objeto EntityDTO con los datos de la entidad (vacío si ocurre un error).
func GetEntityByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.EntityDTO) {
	// Inicializa el error en nil.
	var err error = nil

	// Si no se ha establecido un identificador de conexión, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Se crea el DTO (Data Transfer Object) que almacenará la entidad consultada.
	var entity salvia_daos.EntityDTO = salvia_daos.EntityDTO{}
	if id == "" {
		// Si el parámetro id está vacío, se retorna un error de solicitud incorrecta (Bad Request)
		// utilizando un mensaje localizado.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, entity
		}
	} else {
		// Se configura el criterio de búsqueda para filtrar por "EntityICode".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityICode"},
			AttrsValue: []interface{}{id},
		}

		// Se realiza la consulta de la entidad utilizando el DAO específico.
		err = salvia_daos.GetEntity(by, &entity, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, se retorna un error interno del servidor junto con el mensaje de error.
			return http.StatusInternalServerError, err.Error(), entity
		}

		// Si la consulta es exitosa, se retorna el mensaje de éxito y la entidad encontrada.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(entity), entity
	}
	// Retorno por defecto en caso de no cumplirse las condiciones anteriores.
	return http.StatusInternalServerError, "", entity
}

// GetEntitiesByMoment obtiene y retorna una lista de entidades asociadas a un momento específico.
//
// Parámetros:
//   - moment: Identificador o marca temporal que se utiliza para filtrar las entidades.
//   - module: Nombre del módulo que realiza la solicitud.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente para la base de datos.
//   - dbServerConfig: Configuración del servidor para la base de datos.
//
// Retorna:
//   - Código HTTP que indica el resultado de la operación.
//   - Mensaje en formato JSON describiendo el resultado.
//   - Slice de EntityDTO con los datos de las entidades encontradas.
func GetEntitiesByMoment(moment string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityDTO) {
	var err error
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se ha establecido un identificador de conexión, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa el slice que almacenará los DTO de las entidades.
	var entities []salvia_daos.EntityDTO = []salvia_daos.EntityDTO{}
	if moment == "" {
		// Si el parámetro moment está vacío, se retorna un error de solicitud incorrecta (Bad Request)
		// utilizando un mensaje localizado.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, entities
		}
	} else {
		// Se realiza la consulta para obtener las entidades asociadas al momento especificado.
		entities, err = salvia_daos.GetEntitiesByMoment(moment, connData, &dbClientConfig, &dbServerConfig)

		if err != nil {
			// Si ocurre un error, se asigna el mensaje de error y se retorna un código de error interno del servidor.
			resData = err.Error()
			resCode = http.StatusInternalServerError
		} else {
			// Si la consulta es exitosa, se asigna el mensaje de éxito y el código HTTP correspondiente.
			resCode = http.StatusOK
			resData = utils.CommMsgGetJSONSuccess(entities)
		}
	}
	return resCode, resData, entities
}

// GetEntityByAll obtiene y retorna todas las entidades disponibles para un módulo determinado.
//
// Parámetros:
//   - module: Nombre del módulo desde el cual se realiza la consulta.
//   - connData: Información de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente para la base de datos.
//   - dbServerConfig: Configuración del servidor para la base de datos.
//
// Retorna:
//   - Código HTTP que indica el resultado de la operación.
//   - Mensaje en formato JSON describiendo el resultado.
//   - Slice de EntityDTO con los datos de todas las entidades encontradas.
func GetEntityByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityDTO) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se ha establecido un identificador de conexión, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se consulta la base de datos para obtener todas las entidades usando el DAO correspondiente.
	entities, err := salvia_daos.GetAllEntites(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Si ocurre un error, se asigna el mensaje de error y se retorna un código de error interno del servidor.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la consulta es exitosa, se asigna el mensaje de éxito y el código HTTP correspondiente.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(entities)
	}
	return resCode, resData, entities
}

func GetEntitiesBySector(sectorCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityDTO) {
	var err error
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se ha establecido un identificador de conexión, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializa el slice que almacenará los DTO de las entidades.
	var entities []salvia_daos.EntityDTO = []salvia_daos.EntityDTO{}

	// Se realiza la consulta para obtener las entidades asociadas al momento especificado.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EntitySector"},
		AttrsValue: []interface{}{sectorCode},
	}
	entities, err = salvia_daos.GetEntities(by, connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// Si ocurre un error, se asigna el mensaje de error y se retorna un código de error interno del servidor.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la consulta es exitosa, se asigna el mensaje de éxito y el código HTTP correspondiente.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(entities)
	}

	return resCode, resData, entities
}
