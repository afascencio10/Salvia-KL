// Package salvia_ctrl contiene controladores relacionados con la gestión de FollowUp
// y su asociación con casos de víctimas en el sistema.
package salvia_ctrl

import (
	"net/http"
	"time"

	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
)

// SetFollowUp crea un FollowUp.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func SetFollowUp(dataInput string, victimCaseICode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var followUp salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}
	var followUpEntries []salvia_daos.FollowUpEntryDTO = []salvia_daos.FollowUpEntryDTO{}
	var dtoMap map[string]interface{} = nil

	// Se establecen valores por defecto para la creación de followUp
	salvia_daos.SetFollowUpDefaults(&followUp, common_dao.SQL_INSERT)

	// Se especifican los campos obligatorios exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"FollowUpPhysicalViolenceWitnessedByFamily":  true,
		"FollowUpViolenceEscalation":                 true,
		"FollowUpAssaultWithWeapon":                  true,
		"FollowUpRecentControllingOrJealousBehavior": true,
		"FollowUpViolenceHistoryWithExPartner":       true,
		"FollowUpViolenceHistoryWithOthers":          true,
		"FollowUpSubstanceAbuse":                     true,
		"FollowUpViolenceJustification":              true,
		"FollowUpVictimVulnerability":                true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FollowUpJSONName, common_config.Locale, collectedErrors)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en followUp a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&followUp, dtoMap, salvia_daos.FollowUpJSONName, salvia_daos.FollowUpFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.FollowUpJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	var vcase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{victimCaseICode},
	}

	if err = salvia_daos.GetVictimCase(by, &vcase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se calcula el riesgo y creamos el seguimiento primero para obtener el id
	followUp.FollowUpOwnerGeneralUser = s.UserICode
	followUp.FollowUpRiskLevel = calculateRiskLevel(followUp)

	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if err = salvia_daos.SetFollowUp(&followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	//Ahora definimos y creamos las entradas de seguimiento con sus fechas de vencimiento
	followUpEntries = calculateFollowUpEntries(followUp)

	//Almacenamos las entradas de seguimiento
	for _, fue := range followUpEntries {
		if err = salvia_daos.SetFollowUpEntry(&fue, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}

	//Finalmente lo asociamos al caso
	vcase.VictimCaseFollowUp = followUp
	if err = salvia_daos.UpdateVictimCaseFollowUp(&vcase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// UpdateFollowUpByICode actualiza un FollowUp.
// Recibe como parámetro principal el JSON que viene del cliente HTML
// Datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (vacío en caso de éxito o con el error correspondiente).
func UpdateFollowUpByICode(dataInput string, icode string, connData *db.ConnData, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable error.
	var err error = nil

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Mapa para almacenar errores de validación
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean las estructuras y DTO requeridas para la operación.
	var followUp salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}
	var followUpToUpdate salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}
	var dtoMap map[string]interface{} = nil

	// Se especifican los campos obligatorios para actualizar exigidos al cliente HTML
	var checkFields map[string]bool = map[string]bool{
		"FollowUpPhysicalViolenceWitnessedByFamily":  true,
		"FollowUpViolenceEscalation":                 true,
		"FollowUpAssaultWithWeapon":                  true,
		"FollowUpRecentControllingOrJealousBehavior": true,
		"FollowUpViolenceHistoryWithExPartner":       true,
		"FollowUpViolenceHistoryWithOthers":          true,
		"FollowUpSubstanceAbuse":                     true,
		"FollowUpViolenceJustification":              true,
		"FollowUpVictimVulnerability":                true,
		"FollowUpRiskLevel":                          true,
	}

	// Se obtiene el DTO representado como map a partir de la entrada
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.FollowUpJSONName, common_config.Locale, collectedErrors)

	// Se obtiene el followUp actual a actualizar basado en su icode
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetFollowUp(by, &followUpToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen valores por defecto para la actualización de followUp
	salvia_daos.SetFollowUpDefaults(&followUpToUpdate, common_dao.SQL_UPDATE)

	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		//ValidateJSONInput valida los campos y si son válidos, entonces construye el DTO en followUp a partir del mapa.
		// En caso contrario adiciona los errores a collectedErrors
		utils.ValidateJSONInput(&followUp, dtoMap, salvia_daos.FollowUpJSONName, salvia_daos.FollowUpFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
		//--------------------------------------------------------------------

	default:
		// Error si la estructura del DTO es incorrecta
		utils.SetError(collectedErrors, salvia_daos.FollowUpJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si hay errores en la validación, se retorna el mensaje de error
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	//Actualización de los campos
	followUpToUpdate.FollowUpPhysicalViolenceWitnessedByFamily = followUp.FollowUpPhysicalViolenceWitnessedByFamily
	followUpToUpdate.FollowUpViolenceEscalation = followUp.FollowUpViolenceEscalation
	followUpToUpdate.FollowUpAssaultWithWeapon = followUp.FollowUpAssaultWithWeapon
	followUpToUpdate.FollowUpRecentControllingOrJealousBehavior = followUp.FollowUpRecentControllingOrJealousBehavior
	followUpToUpdate.FollowUpViolenceHistoryWithExPartner = followUp.FollowUpViolenceHistoryWithExPartner
	followUpToUpdate.FollowUpViolenceHistoryWithOthers = followUp.FollowUpViolenceHistoryWithOthers
	followUpToUpdate.FollowUpSubstanceAbuse = followUp.FollowUpSubstanceAbuse
	followUpToUpdate.FollowUpViolenceJustification = followUp.FollowUpViolenceJustification
	followUpToUpdate.FollowUpVictimVulnerability = followUp.FollowUpVictimVulnerability
	followUpToUpdate.FollowUpRiskLevel = followUp.FollowUpRiskLevel
	followUpToUpdate.FollowUpOwnerGeneralUser = followUp.FollowUpOwnerGeneralUser
	followUpToUpdate.FollowUpStatus = followUp.FollowUpStatus

	// Se crea el usuario en la base de datos
	if err = salvia_daos.UpdateFollowUpByICode(&followUpToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna un estado OK en caso de éxito.
	return http.StatusOK, ""
}

// GetFollowUpByICode obtiene el FollowUp dado su código interno.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetFollowUpByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.FollowUpDTO) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las FollowUps y la localidad.
	var followUp salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}
	var followUpEntries []salvia_daos.FollowUpEntryDTO = []salvia_daos.FollowUpEntryDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpICode"},
		AttrsValue: []interface{}{icode},
	}

	// Se consulta el seguimiento.
	if err = salvia_daos.GetFollowUp(by, &followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), followUp
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryFollowUp"},
		AttrsValue: []interface{}{followUp.FollowUpId},
	}
	// Se traen las entradas y complementa el seguimiento
	if followUpEntries, err = salvia_daos.GetFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), followUp
	}

	followUp.FollowUpEntries = followUpEntries

	// Se retorna el estado OK y el mensaje JSON con el followUp.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(followUp), followUp
}

// GetFollowUpByICode obtiene el FollowUp dada su llave principal.
// Recibe el código y los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetFollowUpById(id uint64, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.FollowUpDTO) {
	var err error = nil
	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se inicializan las estructuras DTO requeridas para las FollowUps y la localidad.
	var followUp salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}
	var followUpEntries []salvia_daos.FollowUpEntryDTO = []salvia_daos.FollowUpEntryDTO{}

	// Se construye el criterio de búsqueda para obtener la localidad mediante su código.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpId"},
		AttrsValue: []interface{}{id},
	}

	// Se consulta el seguimiento.
	if err = salvia_daos.GetFollowUp(by, &followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), followUp
	}

	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpEntryFollowUp"},
		AttrsValue: []interface{}{followUp.FollowUpId},
	}
	// Se traen las entradas y complementa el seguimiento
	if followUpEntries, err = salvia_daos.GetFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error(), followUp
	}

	//Se traen las barreras de cada entrada

	for i := 0; i < len(followUpEntries); i++ {
		var barriers []salvia_daos.BarrierDTO = []salvia_daos.BarrierDTO{}
		var rel []salvia_daos.RelBarrierFollowUpEntryDTO = []salvia_daos.RelBarrierFollowUpEntryDTO{}
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"RelBarrierFollowUpEntryFollowUpEntry"},
			AttrsValue: []interface{}{followUpEntries[i].FollowUpEntryId},
		}
		if rel, err = salvia_daos.GetRelBarrierFollowUpEntries(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusInternalServerError, err.Error(), followUp
		}
		for i := 0; i < len(rel); i++ {
			barriers = append(barriers, rel[i].RelBarrierFollowUpEntryBarrier)
			//Ahora se traen las actuaciones

			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"FollowUpEntryActingRelBarrierFollowUpEntry"},
				AttrsValue: []interface{}{rel[i].RelBarrierFollowUpEntryId},
			}
			if barriers[i].BarrierFollowUpEntriesActing, err = salvia_daos.GetFollowUpEntriesActing(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
				return http.StatusInternalServerError, err.Error(), followUp
			}

		}
		followUpEntries[i].FollowUpEntryIdentifiedBarriers = barriers
		//Ahora se traen las actuaciones

	}

	followUp.FollowUpEntries = followUpEntries

	// Se retorna el estado OK y el mensaje JSON con el followUp.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(followUp), followUp
}

// GetAllFollowUps obtiene todas los FollowUps registrados en el sistema.
// Recibe los datos de conexión y configuraciones de base de datos.
// Retorna un código HTTP y un mensaje (JSON con la respuesta o el error).
func GetAllFollowUps(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no ha sido establecida, se libera la conexión al finalizar la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todas las FollowUps registradas.
	followUps, err := salvia_daos.GetAllFollowUps(connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		// Se asigna el mensaje de error en caso de fallo en la consulta.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// En caso de éxito, se asigna el mensaje con las FollowUps en formato JSON.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(followUps)
	}
	return resCode, resData
}

////////////////////////////////////////////////////////////////////////////////
// RemoveFollowUpByICode
////////////////////////////////////////////////////////////////////////////////

// RemoveFollowUpByICode elimina completamente un followUp
// Recibe el icode del followUp
// Retorna un código HTTP y un mensaje JSON.
func RemoveFollowUpByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error = nil
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para el followUp y se establecen valores por defecto
	var followUp salvia_daos.FollowUpDTO = salvia_daos.FollowUpDTO{}

	// Se consulta el usuario a eliminar
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"FollowUpICode"},
		AttrsValue: []interface{}{icode},
	}

	if err = salvia_daos.GetFollowUp(by, &followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se elimina el followUp
	if err = salvia_daos.RemoveFollowUpByICode(&followUp, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, utils.CommMsgGetJSONSuccess("")
}

func calculateRiskLevel(fu salvia_daos.FollowUpDTO) string {
	var riskLevelScore int = 0
	var riskLevel string = ""

	if fu.FollowUpPhysicalViolenceWitnessedByFamily == "y" {
		riskLevelScore += 1
	}
	if fu.FollowUpViolenceEscalation == "y" {
		riskLevelScore += 2
	}
	if fu.FollowUpAssaultWithWeapon == "y" {
		riskLevelScore += 1
	}
	if fu.FollowUpRecentControllingOrJealousBehavior == "y" {
		riskLevelScore += 3
	}
	if fu.FollowUpViolenceHistoryWithExPartner == "y" {
		riskLevelScore += 1
	}
	if fu.FollowUpViolenceHistoryWithOthers == "y" {
		riskLevelScore += 2
	}
	if fu.FollowUpSubstanceAbuse == "y" {
		riskLevelScore += 3
	}
	if fu.FollowUpViolenceJustification == "y" {
		riskLevelScore += 3
	}
	if fu.FollowUpVictimVulnerability == "y" {
		riskLevelScore += 2
	}

	if riskLevelScore > 5 {
		riskLevel = "h"
	} else if riskLevelScore > 3 {
		riskLevel = "m"
	} else {
		riskLevel = "l"
	}

	return riskLevel
}

func calculateFollowUpEntries(fu salvia_daos.FollowUpDTO) []salvia_daos.FollowUpEntryDTO {
	var followUpEntries []salvia_daos.FollowUpEntryDTO = []salvia_daos.FollowUpEntryDTO{}
	var now time.Time = time.Now()

	switch fu.FollowUpRiskLevel {
	case "h":
		//Creamos 7 seguimientos
		var fue salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}

		//A las 4 horas
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.Add(4 * time.Hour)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 1 día
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 1)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 2 día
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 2)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 3 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 3)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 30 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 30)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 2 meses
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 2, 0)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 3 meses
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 3, 0)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

	case "m":

		//Creamos 6 seguimientos
		var fue salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}

		//A las 8 horas
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.Add(8 * time.Hour)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 4 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 4)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 15 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 15)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 30 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 30)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 2 meses
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 2, 0)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 3 meses
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 3, 0)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

	case "l":

		//Creamos 5 seguimientos
		var fue salvia_daos.FollowUpEntryDTO = salvia_daos.FollowUpEntryDTO{}

		//En 1 día
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 1)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 15 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 15)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 30 días
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 0, 30)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 2 meses
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 2, 0)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)

		//En 3 meses
		fue = salvia_daos.FollowUpEntryDTO{}
		salvia_daos.SetFollowUpEntryDefaults(&fue, common_dao.SQL_INSERT)
		fue.FollowUpEntryCompletionDate = now.AddDate(0, 3, 0)
		fue.FollowUpEntryFollowUp = fu
		followUpEntries = append(followUpEntries, fue)
	}

	return followUpEntries
}
