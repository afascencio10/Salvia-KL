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

// SetFollowUpEntry crea un FollowUpEntry.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func SetFollowUpEntry(dataInput string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var followUpEntry salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}
	var dtoMap map[string]interface{} = nil

	// Se establecen valores por defecto para la creación de followUpEntry
	salvia_daos.SetFollowUpEntryDefaults(&followUpEntry, common_dao.SQL_INSERT)

	// Se especifican los campos obligatorios exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"FollowUpEntryOwnerGeneralUser": true,
		"FollowUpWasDone":               true,
		"FollowUpSector":                true,
		"FollowUpPersonVisitedEntity":   true,
		"FollowUpReceivedAttention":     true,
		"FollowUpComments":              true,
		"FollowUpCaseDocumentsPrepared": true,
		"FollowUpEntryStatus":           true,
		"FollowUpEntryFollowUp":         true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FollowUpEntryJSONName, common_config.Locale, collectedErrors)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en followUpEntry a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&followUpEntry, dtoMap, salvia_daos.FollowUpEntryJSONName, salvia_daos.FollowUpEntryFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.FollowUpEntryJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se crea el usuario en la base de datos
	if err = salvia_daos.SetFollowUpEntry(&followUpEntry, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// UpdateFollowUpEntryByICode actualiza un FollowUpEntry.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func UpdateFollowUpEntryByICode(dataInput string, icode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var followUpEntry salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}
	var followUpEntryToUpdate salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}
	var dtoMap map[string]interface{} = nil

	// Se especifican los campos obligatorios para actualizar exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"FollowUpEntryWasDone":               true,
		"FollowUpEntrySector":                true,
		"FollowUpEntryPersonVisitedEntity":   true,
		"FollowUpEntryReceivedAttention":     true,
		"FollowUpEntryComments":              true,
		"FollowUpEntryCaseDocumentsPrepared": true,
		"FollowUpEntryIdentifiedBarriers":    true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FollowUpEntryJSONName, common_config.Locale, collectedErrors)

	// Se obtiene el followUpEntry actual a actualizar basado en su icode
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetFollowUpEntry(by, &followUpEntryToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen valores por defecto para la actualización de followUpEntry
	salvia_daos.SetFollowUpEntryDefaults(&followUpEntryToUpdate, common_dao.SQL_UPDATE)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en followUpEntry a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&followUpEntry, dtoMap, salvia_daos.FollowUpEntryJSONName, salvia_daos.FollowUpEntryFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
		//--------------------------------------------------------------------

		//Traemos las barreras identificadas
		/*
			Reccorreremos la estructura:
			followUp{
				identifiedBarriers[
					{
						icode:"xxxxxxxx"
					}
				]
			}
		*/
		var container interface{}
		var found bool = false
		container, found = dtoMap[salvia_daos.FollowUpEntryJSONName]
		if found {
			switch container := container.(type) {
			case map[string]interface{}:
				var c map[string]interface{} = container
				var barriesTag string = utils.GetTag(&followUpEntry, "FollowUpEntryIdentifiedBarriers", "json")
				if value, found := c[barriesTag]; found {
					var barriers []interface{} = value.([]interface{})
					for _, barrier := range barriers {
						switch b := barrier.(type) {
						case map[string]interface{}:
							switch icode := b["icode"].(type) {
							case string:
								followUpEntry.FollowUpEntryIdentifiedBarriers = append(followUpEntry.FollowUpEntryIdentifiedBarriers, salvia_daos.BarrierDTO{BarrierICode: icode})
							}
						}
					}
				}
			}

		}
	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.FollowUpEntryJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Actualización de los campos
	followUpEntryToUpdate.FollowUpEntryStatus = "p"
	followUpEntryToUpdate.FollowUpEntryOwnerGeneralUser = followUpEntry.FollowUpEntryOwnerGeneralUser
	followUpEntryToUpdate.FollowUpEntryWasDone = followUpEntry.FollowUpEntryWasDone
	followUpEntryToUpdate.FollowUpEntrySector = followUpEntry.FollowUpEntrySector
	followUpEntryToUpdate.FollowUpEntryPersonVisitedEntity = followUpEntry.FollowUpEntryPersonVisitedEntity
	followUpEntryToUpdate.FollowUpEntryReceivedAttention = followUpEntry.FollowUpEntryReceivedAttention
	followUpEntryToUpdate.FollowUpEntryComments = followUpEntry.FollowUpEntryComments
	followUpEntryToUpdate.FollowUpEntryCaseDocumentsPrepared = followUpEntry.FollowUpEntryCaseDocumentsPrepared
	followUpEntryToUpdate.FollowUpEntryOwnerGeneralUser = s.UserICode

	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	// Se actualiza la entrada en la base de datos
	if err = salvia_daos.UpdateFollowUpEntryByICode(&followUpEntryToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	/*
		 Se relacionan las barreras para lo cual
		1. Se eliminan las que hayan anteriores
		2. Se adicionan las nuevas
	*/
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelBarrierFollowUpEntryFollowUpEntry"},
		AttrsValue: []interface{}{followUpEntryToUpdate.FollowUpEntryId},
	}
	if err = salvia_daos.RemoveRelBarrierFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	for _, b := range followUpEntry.FollowUpEntryIdentifiedBarriers {
		var rel salvia_daos.RelBarrierFollowUpEntryDTO = salvia_daos.RelBarrierFollowUpEntryDTO{}
		var barrier salvia_daos.BarrierDTO
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"BarrierICode"},
			AttrsValue: []interface{}{b.BarrierICode},
		}

		if err = salvia_daos.GetBarrier(by, &barrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		rel.RelBarrierFollowUpEntryFollowUpEntry = followUpEntryToUpdate
		rel.RelBarrierFollowUpEntryBarrier = barrier

		if err = salvia_daos.SetRelBarrierFollowUpEntry(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// GetFollowUpEntryByICode obtiene el FollowUpEntry dado su código interno.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetFollowUpEntryByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para la FollowUpEntry y la localidad.
	var followUpEntry salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryICode"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta la entrada de seguimiento.
	if err = salvia_daos.GetFollowUpEntry(by, &followUpEntry, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con la entrada de seguimiento.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(followUpEntry)
}

// GetFollowUpEntriesByFollowUpICode obtiene todas las FollowUpEntries dado el código de su padre FollowUp.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetFollowUpEntriesByFollowUpICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las FollowUpEntries.
	var followUpEntries []salvia_daos.FollowUpEntryDTO = []salvia_daos.FollowUpEntryDTO{}

	// Se construye el criterio de búsqueda.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryFollowUp"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta las entradas de seguimiento.
	if followUpEntries, err = salvia_daos.GetFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con las entradas de seguimiento.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(followUpEntries)
}

// GetAllFollowUpEntries obtiene todas las FollowUpEntries registradas en el sistema.
// Recibe los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAllFollowUpEntries(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todas las FollowUpEntries registradas.
	followUpEntries, err := salvia_daos.GetAllFollowUpEntries(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se asigna el mensaje de error en caso de fallo en la consulta.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, se asigna el mensaje con las FollowUpEntries en formato JSON.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(followUpEntries)
	}
	return resCode, resData
}

////////////////////////////////////////////////////////////////////////////////
// RemoveFollowUpEntryByICode
////////////////////////////////////////////////////////////////////////////////

// RemoveFollowUpEntryByICode elimina completamente un followUpEntry
// Recibe el icode del followUpEntry
// Retorna un código HTTP y un mensaje JSON.
func RemoveFollowUpEntryByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el followUpEntry y se establecen valores por defecto
	var followUpEntry salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}

	// Se consulta la entrada a eliminar
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetFollowUpEntry(by, &followUpEntry, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina la entrada de seguimiento
	if err = salvia_daos.RemoveFollowUpEntryByICode(&followUpEntry, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess("")
}
