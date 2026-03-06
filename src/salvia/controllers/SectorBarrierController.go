// Package salvia_ctrl contiene controladores relacionados con la gestión de SectorBarrier
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

// SetSectorBarrier crea un SectorBarrier.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func SetSectorBarrier(dataInput string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var sectorBarrier salvia_daos.SectorBarrierDTO = salvia_daos.SectorBarrierDTO{}
	var dtoMap map[string]interface{} = nil

	// Se establecen valores por defecto para la creación de sectorBarrier
	salvia_daos.SetSectorBarrierDefaults(&sectorBarrier, common_dao.SQL_INSERT)

	// Se especifican los campos obligatorios exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"SectorBarrierICode":        true,
		"SectorBarrierCreationDate": true,
		"SectorBarrierUpdateDate":   true,
		"SectorBarrierName":         true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.SectorBarrierJSONName, common_config.Locale, collectedErrors)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en sectorBarrier a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&sectorBarrier, dtoMap, salvia_daos.SectorBarrierJSONName, salvia_daos.SectorBarrierFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.SectorBarrierJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se crea el usuario en la base de datos
	if err = salvia_daos.SetSectorBarrier(&sectorBarrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// UpdateSectorBarrierByICode actualiza un SectorBarrier.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func UpdateSectorBarrierByICode(dataInput string, icode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var sectorBarrier salvia_daos.SectorBarrierDTO = salvia_daos.SectorBarrierDTO{}
	var sectorBarrierToUpdate salvia_daos.SectorBarrierDTO = salvia_daos.SectorBarrierDTO{}
	var dtoMap map[string]interface{} = nil

	// Se especifican los campos obligatorios para actualizar exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"SectorBarrierICode":        true,
		"SectorBarrierCreationDate": true,
		"SectorBarrierUpdateDate":   true,
		"SectorBarrierName":         true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.SectorBarrierJSONName, common_config.Locale, collectedErrors)

	// Se obtiene el sectorBarrier actual a actualizar basado en su icode
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"SectorBarrierICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetSectorBarrier(by, &sectorBarrierToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen valores por defecto para la actualización de sectorBarrier
	salvia_daos.SetSectorBarrierDefaults(&sectorBarrierToUpdate, common_dao.SQL_UPDATE)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en sectorBarrier a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&sectorBarrier, dtoMap, salvia_daos.SectorBarrierJSONName, salvia_daos.SectorBarrierFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.SectorBarrierJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Actualización de los campos
	sectorBarrierToUpdate.SectorBarrierName = sectorBarrier.SectorBarrierName
	sectorBarrierToUpdate.SectorBarrierDescription = sectorBarrier.SectorBarrierDescription
	sectorBarrierToUpdate.SectorBarrierSector = sectorBarrier.SectorBarrierSector

	// Se crea el usuario en la base de datos
	if err = salvia_daos.UpdateSectorBarrierByICode(&sectorBarrierToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// GetSectorBarrierByICode obtiene el SectorBarrier dado su código interno.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetSectorBarrierByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las SectorBarriers y la localidad.
	var sectorBarrier salvia_daos.SectorBarrierDTO = salvia_daos.SectorBarrierDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"SectorBarrierICode"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta la localidad.
	if err = salvia_daos.GetSectorBarrier(by, &sectorBarrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con el sectorBarrier.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(sectorBarrier)
}

// GetSectorBarriersBySector obtiene todas las SectorBarriers dado el código de su padre Sector.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetSectorBarriersBySector(sectorCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.SectorBarrierDTO) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las SectorBarriers y la localidad.
	var sectorBarriers []salvia_daos.SectorBarrierDTO = []salvia_daos.SectorBarrierDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"SectorBarrierSector"},
		AttrsValue: []interface{}{sectorCode},
	}

	// Se consulta la localidad.
	if sectorBarriers, err = salvia_daos.GetSectorBarriers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), sectorBarriers
	}

	// Se retorna el estado OK y el mensaje JSON con el sectorBarrier.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(sectorBarriers), sectorBarriers
}

// GetAllSectorBarriers obtiene todas los SectorBarriers registrados en el sistema.
// Recibe los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAllSectorBarriers(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.SectorBarrierDTO) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todas las SectorBarriers registradas.
	sectorBarriers, err := salvia_daos.GetAllSectorBarriers(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se asigna el mensaje de error en caso de fallo en la consulta.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, se asigna el mensaje con las SectorBarriers en formato JSON.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(sectorBarriers)
	}
	return resCode, resData, sectorBarriers
}

////////////////////////////////////////////////////////////////////////////////
// RemoveSectorBarrierByICode
////////////////////////////////////////////////////////////////////////////////

// RemoveSectorBarrierByICode elimina completamente un sectorBarrier
// Recibe el icode del sectorBarrier
// Retorna un código HTTP y un mensaje JSON.
func RemoveSectorBarrierByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el sectorBarrier y se establecen valores por defecto
	var sectorBarrier salvia_daos.SectorBarrierDTO = salvia_daos.SectorBarrierDTO{}

	// Se consulta el usuario a eliminar
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"SectorBarrierICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetSectorBarrier(by, &sectorBarrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el sectorBarrier
	if err = salvia_daos.RemoveSectorBarrierByICode(&sectorBarrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess("")
}
