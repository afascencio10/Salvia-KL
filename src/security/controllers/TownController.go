// Package security_ctrl contiene los controladores relacionados con la seguridad,
// específicamente las operaciones para obtener información de "towns" (pueblos)
// a partir de diferentes criterios de búsqueda.
package security_ctrl

import (
	// Importa la configuración común
	common_config "bitsflow/common/config"
	// Importa controladores comunes para operaciones compartidas
	common_controllers "bitsflow/common/controllers"
	// Importa funciones y constantes DAO comunes
	common_dao "bitsflow/common/dao"
	// Importa funcionalidades para manejo de conexiones a la base de datos
	"bitsflow/common/db"
	// Importa utilidades comunes, por ejemplo para formatear mensajes
	"bitsflow/common/utils"
	// Importa los DAO específicos de seguridad
	security_daos "bitsflow/security/dao"
	"net/http"
)

// GetTownByICode obtiene la información de un pueblo utilizando su ICode.
// Recibe como parámetros el identificador del pueblo (id), el módulo, datos de conexión,
// configuración del cliente y del servidor de la base de datos.
// Retorna un código HTTP y un mensaje en formato JSON con el resultado de la operación.
func GetTownByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no posee un identificador, se asegura liberar la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO (Data Transfer Object) para representar el pueblo.
	var town security_daos.TownDTO = security_daos.TownDTO{}

	// Verifica si el identificador está vacío y retorna un error de solicitud en ese caso.
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	} else {
		// Define los criterios de búsqueda utilizando el ICode del pueblo.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"TownICode"},
			AttrsValue: []interface{}{id},
		}

		// Realiza la búsqueda del pueblo en la base de datos mediante el DAO de seguridad.
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna un error interno del servidor en caso de falla en la consulta.
			return http.StatusInternalServerError, err.Error()
		}

		// Retorna el mensaje de éxito con el pueblo encontrado, formateado en JSON.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(town)
	}
	// En caso de condiciones no satisfechas, retorna un error interno.
	return http.StatusInternalServerError, ""
}

// GetTownsByCity obtiene la lista de pueblos asociados a una ciudad específica.
// Recibe el identificador de la ciudad (cityId), el módulo, datos de conexión,
// y configuraciones del cliente y del servidor de la base de datos.
// Retorna un código HTTP y un mensaje en formato JSON que contiene la lista de pueblos o un mensaje de error.
func GetTownsByCity(cityId uint64, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Libera la conexión si el identificador de la conexión está vacío.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se declara un slice de DTOs para almacenar los pueblos obtenidos.
	var towns []security_daos.TownDTO = []security_daos.TownDTO{}

	// Verifica que el identificador de la ciudad no sea cero, retornando un error si lo es.
	if cityId == 0 {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	} else {
		// Define los criterios de búsqueda para obtener los pueblos asociados a la ciudad.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_OR,
			EqualSign:  []string{},
			AttrsName:  []string{"TownCity"},
			AttrsValue: []interface{}{cityId},
		}

		// Llama a la función GetTowns del DAO de seguridad para realizar la búsqueda.
		towns, err = security_daos.GetTowns(by, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna un error interno en caso de que la consulta falle.
			return http.StatusInternalServerError, err.Error()
		}

		// Retorna el mensaje de éxito con la lista de pueblos en formato JSON.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(towns)
	}
	// Retorna un error interno si no se cumple alguna condición.
	return http.StatusInternalServerError, ""
}

// GetTownByTownCode obtiene la información de un pueblo utilizando su código único (townCode).
// Recibe como parámetros el código del pueblo, el módulo, datos de conexión, y configuraciones
// del cliente y del servidor de la base de datos.
// Retorna un código HTTP, un mensaje en formato JSON y el DTO del pueblo encontrado.
func GetTownByTownCode(townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, security_daos.TownDTO) {
	var err error = nil
	// Libera la conexión si no se dispone de un identificador de conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para representar el pueblo.
	var town security_daos.TownDTO = security_daos.TownDTO{}

	// Verifica si el código del pueblo está vacío; en ese caso, retorna un error.
	if townCode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, town
		}
	} else {
		// Define los criterios de búsqueda utilizando el código del pueblo.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_OR,
			EqualSign:  []string{},
			AttrsName:  []string{"TownCode"},
			AttrsValue: []interface{}{townCode},
		}

		// Realiza la búsqueda del pueblo mediante el DAO de seguridad.
		err = security_daos.GetTown(by, &town, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna un error interno del servidor en caso de falla.
			return http.StatusInternalServerError, err.Error(), town
		}

		// Retorna el pueblo encontrado junto con un mensaje de éxito formateado en JSON.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(town), town
	}
	// Retorna un error interno si no se cumplen las condiciones necesarias.
	return http.StatusInternalServerError, "", town
}

// GetTownByAll obtiene todos los pueblos disponibles en la base de datos.
// Recibe el módulo, datos de conexión, y configuraciones del cliente y del servidor de la base de datos.
// Retorna un código HTTP y un mensaje en formato JSON que contiene la lista de pueblos o un mensaje de error.
func GetTownByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Libera la conexión si el identificador de conexión no está definido.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Llama a la función GetAllTowns del DAO de seguridad para obtener la lista de todos los pueblos.
	users, err := security_daos.GetAllTowns(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// En caso de error, asigna el mensaje de error y establece el código HTTP de error.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la operación es exitosa, asigna el resultado en formato JSON y el código HTTP OK.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(users)
	}
	return resCode, resData
}
