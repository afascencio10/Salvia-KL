// Package salvia_ctrl contiene controladores relacionados con la gestión de Barrier
// y su asociación con casos de víctimas en el sistema.
package salvia_ctrl

import (
	"net/http"

	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
)

// SetBarrier crea un Barrier.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func SetBarrier(dataInput string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var barrier salvia_daos.BarrierDTO = salvia_daos.BarrierDTO{}
	var dtoMap map[string]interface{} = nil

	// Se establecen valores por defecto para la creación de barrier
	salvia_daos.SetBarrierDefaults(&barrier, common_dao.SQL_INSERT)

	// Se especifican los campos obligatorios exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"BarrierICode":        true,
		"BarrierCreationDate": true,
		"BarrierUpdateDate":   true,
		"BarrierName":         true,
		"BarrierDescription":  false,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.BarrierJSONName, common_config.Locale, collectedErrors)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en barrier a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&barrier, dtoMap, salvia_daos.BarrierJSONName, salvia_daos.BarrierFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.BarrierJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se crea el usuario en la base de datos
	if err = salvia_daos.SetBarrier(&barrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// UpdateBarrierByICode actualiza un Barrier.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func UpdateBarrierByICode(dataInput string, icode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var barrier salvia_daos.BarrierDTO = salvia_daos.BarrierDTO{}
	var barrierToUpdate salvia_daos.BarrierDTO = salvia_daos.BarrierDTO{}
	var dtoMap map[string]interface{} = nil

	// Se especifican los campos obligatorios para actualizar exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"BarrierICode":        true,
		"BarrierCreationDate": true,
		"BarrierUpdateDate":   true,
		"BarrierName":         true,
		"BarrierDescription":  false,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.BarrierJSONName, common_config.Locale, collectedErrors)

	// Se obtiene el barrier actual a actualizar basado en su icode
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"BarrierICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetBarrier(by, &barrierToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen valores por defecto para la actualización de barrier
	salvia_daos.SetBarrierDefaults(&barrierToUpdate, common_dao.SQL_UPDATE)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en barrier a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&barrier, dtoMap, salvia_daos.BarrierJSONName, salvia_daos.BarrierFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.BarrierJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Actualización de los campos
	barrierToUpdate.BarrierName = barrier.BarrierName
	barrierToUpdate.BarrierDescription = barrier.BarrierDescription

	// Se crea el usuario en la base de datos
	if err = salvia_daos.UpdateBarrierByICode(&barrierToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// GetBarrierByICode obtiene el Barrier dado su código interno.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetBarrierByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las Barriers y la localidad.
	var barrier salvia_daos.BarrierDTO = salvia_daos.BarrierDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"BarrierICode"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta la localidad.
	if err = salvia_daos.GetBarrier(by, &barrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con el barrier.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(barrier)
}

// GetBarriersBySectorBarrier obtiene todas las Barriers dado el código de su padre SectorBarrier.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetBarriersBySectorBarrier(sectorICode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las Barriers y la localidad.
	var barriers []salvia_daos.BarrierDTO = []salvia_daos.BarrierDTO{}

	// Se construye el criterio de búsqueda para obtener el sector mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"SectorBarrierICode"},
		AttrsValue: []interface{}{sectorICode},
	}

	var sectorBarrier salvia_daos.SectorBarrierDTO

	// Se consulta el sector de barrera.
	if err = salvia_daos.GetSectorBarrier(by, &sectorBarrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	//Ahora se ejecuta la acción principal
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"BarrierSectorBarrier"},
		AttrsValue: []interface{}{sectorBarrier.SectorBarrierId},
	}

	// Se consultan las barreras.
	if barriers, err = salvia_daos.GetBarriers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con el barrier.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(barriers)
}

// GetAllBarriers obtiene todas los Barriers registrados en el sistema.
// Recibe los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAllBarriers(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todas las Barriers registradas.
	barriers, err := salvia_daos.GetAllBarriers(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se asigna el mensaje de error en caso de fallo en la consulta.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, se asigna el mensaje con las Barriers en formato JSON.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(barriers)
	}
	return resCode, resData
}

////////////////////////////////////////////////////////////////////////////////
// RemoveBarrierByICode
////////////////////////////////////////////////////////////////////////////////

// RemoveBarrierByICode elimina completamente un barrier
// Recibe el icode del barrier
// Retorna un código HTTP y un mensaje JSON.
func RemoveBarrierByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el barrier y se establecen valores por defecto
	var barrier salvia_daos.BarrierDTO = salvia_daos.BarrierDTO{}

	// Se consulta el usuario a eliminar
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"BarrierICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetBarrier(by, &barrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el barrier
	if err = salvia_daos.RemoveBarrierByICode(&barrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess("")
}
