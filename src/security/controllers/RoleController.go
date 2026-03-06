// Package security_ctrl implementa funciones relacionadas con el manejo
// de roles de seguridad, interactuando con los DAOs correspondientes para
// obtener la información necesaria.
package security_ctrl

import (
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"
)

// GetRoleByAll obtiene todos los roles disponibles para un módulo específico.
// Esta función se conecta a la base de datos a través de los parámetros provistos
// y utiliza el DAO de seguridad para recuperar la información de los roles.
//
// Parámetros:
//   - module: Nombre del módulo para el cual se requieren los roles.
//   - connData: Puntero a la estructura de conexión de la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - int: Código de estado (200 para éxito, 500 para error).
//   - string: Mensaje en formato JSON con el resultado o descripción del error.
//   - []security_daos.RoleDTO: Lista de roles obtenidos.
func GetRoleByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []security_daos.RoleDTO) {
	var resData string = ""
	var resCode int = 500

	// Si no se dispone de un identificador de conexión, se libera la conexión
	// al finalizar la ejecución de la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se realiza la llamada al DAO para obtener todos los roles.
	// El primer valor retornado no se utiliza, el segundo es la lista de roles,
	// y el tercero es el error en caso de que ocurra.
	roles, err := security_daos.GetAllRoles(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// En caso de error, se captura el mensaje de error y se asigna el código 500.
		resData = err.Error()
		resCode = 500
	} else {
		// Si la operación es exitosa, se asigna el código 200 y se formatea la respuesta
		// en un mensaje JSON indicando éxito.
		resCode = 200
		resData = utils.CommMsgGetJSONSuccess(roles)
	}
	return resCode, resData, roles
}

// GetRoles obtiene roles específicos para un módulo, aplicando un filtro basado en la estructura "By".
// Este filtro excluye aquellos roles cuyo "RoleCode" sea igual a "us".
// La función se conecta a la base de datos y utiliza el DAO de seguridad para realizar la consulta.
//
// Parámetros:
//   - module: Nombre del módulo para el cual se requieren los roles.
//   - connData: Puntero a la estructura de conexión de la base de datos.
//   - dbClientConfig: Configuración del cliente de la base de datos.
//   - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - int: Código de estado (200 para éxito, 500 para error).
//   - string: Mensaje en formato JSON con el resultado o descripción del error.
//   - []security_daos.RoleDTO: Lista de roles filtrados obtenidos.
func GetRoles(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []security_daos.RoleDTO) {
	var resData string = ""
	var resCode int = 500

	// Se define el filtro de búsqueda para excluir roles con "RoleCode" igual a "us".
	// La estructura "By" especifica que se usará el operador AND y la condición de desigualdad.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RoleCode"},
		AttrsValue: []interface{}{"us"},
		EqualSign:  []string{"<>"},
	}

	// Si no se dispone de un identificador de conexión, se libera la conexión
	// al finalizar la ejecución de la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se realiza la llamada al DAO para obtener los roles aplicando el filtro definido.
	// El primer valor retornado no se utiliza, el segundo es la lista de roles,
	// y el tercero es el error en caso de que ocurra.
	roles, err := security_daos.GetRoles(by, connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// En caso de error, se captura el mensaje y se asigna el código 500.
		resData = err.Error()
		resCode = 500
	} else {
		// Si la operación es exitosa, se asigna el código 200 y se formatea la respuesta
		// en un mensaje JSON indicando éxito.
		resCode = 200
		resData = utils.CommMsgGetJSONSuccess(roles)
	}
	return resCode, resData, roles
}
