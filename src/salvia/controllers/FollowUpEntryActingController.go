// Package salvia_ctrl contiene controladores relacionados con la gestión de FollowUpEntryActing
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

// SetFollowUpEntryActing crea un FollowUpEntryActing.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func SetFollowUpEntryActing(dataInput string, barrierICode string, entryICode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var followUpEntryActing salvia_daos.FollowUpEntryActingDTO = salvia_daos.FollowUpEntryActingDTO{}
	var dtoMap map[string]interface{} = nil

	// Se establecen valores por defecto para la creación de followUpEntryActing
	salvia_daos.SetFollowUpEntryActingDefaults(&followUpEntryActing, common_dao.SQL_INSERT)

	// Se especifican los campos obligatorios exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"FollowUpEntryActingDescription": true,
		"FollowUpEntryActionStatus":      true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FollowUpEntryActingJSONName, common_config.Locale, collectedErrors)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en followUpEntryActing a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&followUpEntryActing, dtoMap, salvia_daos.FollowUpEntryActingJSONName, salvia_daos.FollowUpEntryActingFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------
	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.FollowUpEntryActingJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	var barrier salvia_daos.BarrierDTO = salvia_daos.BarrierDTO{}
	var entry salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}
	var rel salvia_daos.RelBarrierFollowUpEntryDTO = salvia_daos.RelBarrierFollowUpEntryDTO{}
	var rels []salvia_daos.RelBarrierFollowUpEntryDTO = []salvia_daos.RelBarrierFollowUpEntryDTO{}

	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"BarrierICode"},
		AttrsValue: []interface{}{barrierICode},
	}

	if err = salvia_daos.GetBarrier(by, &barrier, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryICode"},
		AttrsValue: []interface{}{entryICode},
	}
	if err = salvia_daos.GetFollowUpEntry(by, &entry, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelBarrierFollowUpEntryFollowUpEntry", "RelBarrierFollowUpEntryBarrier"},
		AttrsValue: []interface{}{entry.FollowUpEntryId, barrier.BarrierId},
	}
	if err = salvia_daos.GetRelBarrierFollowUpEntry(by, &rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// No existe la relacion
		salvia_daos.SetRelBarrierFollowUpEntryDefaults(&rel, common_dao.SQL_INSERT)
		rel.RelBarrierFollowUpEntryFollowUpEntry = entry
		rel.RelBarrierFollowUpEntryBarrier = barrier

		if err = salvia_daos.SetRelBarrierFollowUpEntry(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	followUpEntryActing.FollowUpEntryActingRelBarrierFollowUpEntry = rel
	followUpEntryActing.FollowUpEntryActingOwnerGeneralUser = s.UserICode

	if err = salvia_daos.SetFollowUpEntryActing(&followUpEntryActing, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	//Si el estado es terminado, se verifica si todas las actuaciones han sido terminadas, caso en el cual se cerraría la actuación completa para la entrada de seguimiento.

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"RelBarrierFollowUpEntryFollowUpEntry"},
		AttrsValue: []interface{}{entry.FollowUpEntryId},
	}
	if rels, err = salvia_daos.GetRelBarrierFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	for _, rel := range rels {
		var barrierDone map[uint64]bool = make(map[uint64]bool)
		var barrierNotDone map[uint64]bool = make(map[uint64]bool)

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"FollowUpEntryActingRelBarrierFollowUpEntry"},
			AttrsValue: []interface{}{rel.RelBarrierFollowUpEntryId},
		}
		var acting []salvia_daos.FollowUpEntryActingDTO
		if acting, err = salvia_daos.GetFollowUpEntriesActing(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
		//Por cada actuación que haya puesto el estado en terminado (d), se marca esa relación con la barrera como solucionado en barrierDone, de lo contrario se marca la relación en barrierNotDone

		//Si no hay actuaciones, significa que las barreras están sin resolver
		if len(acting) == 0 {
			break
		}
		for _, a := range acting {
			if a.FollowUpEntryActionStatus == "d" {
				barrierDone[rel.RelBarrierFollowUpEntryId] = true
				break
			} else {
				barrierNotDone[rel.RelBarrierFollowUpEntryId] = true
			}
		}
		//Se eliminan las barreras duplicadas. Aunque se supone que no es necesario porque la última es la que cambia a "d", se hizo de esta manera sólo por precausión por si llega en desorden, poder detectar la que tenga "d".
		for k := range barrierNotDone {
			if _, found := barrierDone[k]; found {
				delete(barrierNotDone, k)
			}
		}
		//Si no hay elementos en barrierNotDone quiere decir que todas las barreras fueron resueltas
		if len(barrierNotDone) == 0 && entry.FollowUpEntryStatus == "p" {
			//Se cambia el estado a la entrada
			salvia_daos.SetFollowUpEntryDefaults(&entry, common_dao.SQL_UPDATE)
			entry.FollowUpEntryStatus = "d"

			if err = salvia_daos.UpdateFollowUpEntryByICode(&entry, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, err.Error()
			}
		}
	}

	//Si todas las entradas de seguimiento están completas, entonces se actualiza todo el caso
	var entries []salvia_daos.FollowUpEntryDTO
	var followUp salvia_daos.FollowUpDTO
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryFollowUp"},
		AttrsValue: []interface{}{entry.FollowUpEntryFollowUp.FollowUpId},
	}
	if entries, err = salvia_daos.GetFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}
	var isComplete bool = true
	for _, e := range entries {
		if e.FollowUpEntryStatus != "d" {
			isComplete = false
			break
		}
	}

	if isComplete {
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"FollowUpId"},
			AttrsValue: []interface{}{entry.FollowUpEntryFollowUp},
		}
		if err = salvia_daos.GetFollowUp(by, &followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}

		salvia_daos.SetFollowUpDefaults(&followUp, common_dao.SQL_UPDATE)
		followUp.FollowUpStatus = "d"

		if err = salvia_daos.UpdateFollowUpByICode(&followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
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

// UpdateFollowUpEntryActingByICode actualiza un FollowUpEntryActing.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func UpdateFollowUpEntryActingByICode(dataInput string, icode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var followUpEntryActing salvia_daos.FollowUpEntryActingDTO = salvia_daos.FollowUpEntryActingDTO{}
	var followUpEntryActingToUpdate salvia_daos.FollowUpEntryActingDTO = salvia_daos.FollowUpEntryActingDTO{}
	var dtoMap map[string]interface{} = nil

	// Se especifican los campos obligatorios para actualizar exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"FollowUpEntryActingOwnerGeneralUser": true,
		"FollowUpEntryActingDescription":      true,
		"FollowUpEntryActionStatus":           true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FollowUpEntryActingJSONName, common_config.Locale, collectedErrors)

	// Se obtiene el followUpEntryActing actual a actualizar basado en su icode
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryActingICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetFollowUpEntryActing(by, &followUpEntryActingToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen valores por defecto para la actualización de followUpEntryActing
	salvia_daos.SetFollowUpEntryActingDefaults(&followUpEntryActingToUpdate, common_dao.SQL_UPDATE)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en followUpEntryActing a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&followUpEntryActing, dtoMap, salvia_daos.FollowUpEntryActingJSONName, salvia_daos.FollowUpEntryActingFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.FollowUpEntryActingJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Actualización de los campos
	followUpEntryActingToUpdate.FollowUpEntryActingOwnerGeneralUser = followUpEntryActing.FollowUpEntryActingOwnerGeneralUser
	followUpEntryActingToUpdate.FollowUpEntryActingDescription = followUpEntryActing.FollowUpEntryActingDescription
	followUpEntryActingToUpdate.FollowUpEntryActionStatus = followUpEntryActing.FollowUpEntryActionStatus

	if err = salvia_daos.UpdateFollowUpEntryActingByICode(&followUpEntryActingToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// GetFollowUpEntryActingByICode obtiene el FollowUpEntryActing dado su código interno.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetFollowUpEntryActingByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las FollowUpEntryActings y la localidad.
	var followUpEntryActing salvia_daos.FollowUpEntryActingDTO = salvia_daos.FollowUpEntryActingDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryActingICode"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta la localidad.
	if err = salvia_daos.GetFollowUpEntryActing(by, &followUpEntryActing, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con el followUpEntryActing.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(followUpEntryActing)
}

// GetFollowUpEntryActingByICode obtiene todas las  FollowUpEntryActing dado el código de su padre FollowUpEntry.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetFollowUpEntriesActingByFollowUpEntryICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las FollowUpEntryActings y la localidad.
	var followUpEntriesActing []salvia_daos.FollowUpEntryActingDTO = []salvia_daos.FollowUpEntryActingDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryActingFollowUpEntry"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta la localidad.
	if followUpEntriesActing, err = salvia_daos.GetFollowUpEntriesActing(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado OK y el mensaje JSON con el followUpEntryActing.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(followUpEntriesActing)
}

// GetFollowUpEntryActingByAll obtiene todas los FollowUpEntryActings registrados en el sistema.
// Recibe los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAllFollowUpEntriesActing(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todas las FollowUpEntryActings registradas.
	followUpEntryActings, err := salvia_daos.GetAllFollowUpEntriesActing(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se asigna el mensaje de error en caso de fallo en la consulta.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, se asigna el mensaje con las FollowUpEntryActings en formato JSON.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(followUpEntryActings)
	}
	return resCode, resData
}

////////////////////////////////////////////////////////////////////////////////
// RemoveFollowUpEntryActingByICode
////////////////////////////////////////////////////////////////////////////////

// RemoveFollowUpEntryActingByICode elimina completamente un followUpEntryActing
// Recibe el icode del followUpEntryActing
// Retorna un código HTTP y un mensaje JSON.
func RemoveFollowUpEntryActingByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el followUpEntryActing y se establecen valores por defecto
	var followUpEntryActing salvia_daos.FollowUpEntryActingDTO = salvia_daos.FollowUpEntryActingDTO{}

	// Se consulta el usuario a eliminar
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryActingICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetFollowUpEntryActing(by, &followUpEntryActing, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el followUpEntryActing
	if err = salvia_daos.RemoveFollowUpEntryActingByICode(&followUpEntryActing, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess("")
}
