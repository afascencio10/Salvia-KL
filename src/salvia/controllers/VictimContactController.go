// Package salvia_ctrl contiene los controladores para manejar operaciones relacionadas con
// "VictimContact" en el módulo Salvia, incluyendo la creación, consulta e invalidación de registros.
package salvia_ctrl

import (
	"math"
	"net/http"

	"github.com/dchest/captcha"

	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"

	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"

	security_config "bitsflow/security/config"
)

type VictimContactRequest struct {
	VContact salvia_daos.VictimContactDTO `json:"victimContact"`
}

// SetVictimContact crea o actualiza un registro de VictimContact en la base de datos.
// Recibe como parámetros la entrada de datos en formato JSON (dataInput), el módulo actual,
// la información de conexión (connData) y las configuraciones de cliente y servidor de la base de datos.
// Retorna un código de estado HTTP y un mensaje en formato JSON (con información de éxito o error).
func SetVictimContact(dataInput string, checkCaptcha bool, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	var err error = nil

	// Si no se provee un ID de conexión, se libera la conexión al final de la función.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Mapa para recopilar errores durante el procesamiento del DTO.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crea el DTO (Data Transfer Object) para VictimContact.
	var vContactRequest VictimContactRequest = VictimContactRequest{}

	var dtoMap map[string]interface{} = nil

	// Se obtiene el mapa DTO a partir del JSON de entrada.
	// La función GetDTOMap convierte el string JSON en un mapa de datos.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.VictimContactJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &vContactRequest)

	// Se establecen los valores por defecto en el DTO.
	salvia_daos.SetVictimContactDefaults(&vContactRequest.VContact, common_dao.SQL_INSERT)
	salvia_daos.SetVictimContactForm2Defaults(&vContactRequest.VContact.VictimContactForm2, common_dao.SQL_INSERT)

	var captchaPassed bool = false

	// Se verifica que los datos tengan la estructura esperada.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto
		// Definición de campos requeridos y opcionales.
		var checkFields map[string]bool = map[string]bool{
			"VictimContactNames":     true,
			"VictimContactLastNames": true,
			"VictimContactLatitude":  false,
			"VictimContactLongitude": false,
		}

		// Valida la estructura del JSON y asigna los valores al DTO.
		utils.ValidateJSONInput(&vContactRequest.VContact, dtoMap, salvia_daos.VictimContactJSONName, salvia_daos.VictimContactFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		checkFields = map[string]bool{
			"VictimContactForm2WillReceiveCall":     false,
			"VictimContactForm2HasCareRole":         false,
			"VictimContactForm2VictimAwareOfReport": false,
			"VictimContactForm2ReporterNames":       false,
			"VictimContactForm2ReporterPhone":       false,
			"VictimContactForm2VictimColPhone":      true,
			"VictimContactForm2FactsDescription":    true,
			"VictimContactForm2BestContactTime":     true,
			"VictimContactForm2ReportType":          true,
			"VictimContactForm2ReportTypeDetails":   false,
		}

		// Valida la estructura del JSON y asigna los valores al DTO.
		utils.ValidateJSONInput(&vContactRequest.VContact.VictimContactForm2, dtoMap, salvia_daos.VictimContactJSONName+"."+salvia_daos.VictimContactForm2JSONName, salvia_daos.VictimContactForm2FieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		//--------------------------------------------------------------------
		// Extrae y verifica el captcha.

		var captchaID string = vContactRequest.VContact.VictimContactCaptchaID
		var captchaSolution string = vContactRequest.VContact.VictimContactCaptchaSolution
		// Se extrae el ID del captcha.

		// Verifica la solución del captcha.
		captchaPassed = !checkCaptcha || captcha.VerifyString(captchaID, captchaSolution)

		if len(collectedErrors) > 0 {
			// Error si la estructura del DTO es incorrecta
			utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
		}
		/*
			container, found := dtoMap[salvia_daos.VictimContactJSONName]
			if found {
				switch container := container.(type) {
				case map[string]interface{}:
					// Se obtienen las etiquetas de captcha definidas en el DTO.
					var captchaIDTag string = utils.GetTag(&victimContact, "VictimContactCaptchaID", "json")
					var captchaSolutionTag string = utils.GetTag(&victimContact, "VictimContactCaptchaSolution", "json")
					var captchaID string
					var captchaSolution string
					// Se extrae el ID del captcha.
					if value, found := container[captchaIDTag]; found {
						captchaID = value.(string)
					}
					// Se extrae la solución del captcha.
					if value, found := container[captchaSolutionTag]; found {
						captchaSolution = value.(string)
					}

					// Verifica la solución del captcha.
					captchaPassed = !checkCaptcha || captcha.VerifyString(captchaID, captchaSolution)
				}
			}*/
	default:
		// Si la estructura no es la esperada, se registra un error.
		utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si el captcha no es válido, se registra el error correspondiente.
	if !captchaPassed {
		utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, utils.GetTag(&vContactRequest.VContact, "VictimContactCaptchaSolution", "json"), "security_general_user_captcha_fail", "common_global_error", security_config.Locale)
	}

	getAndVerifyVictimContactEnumsMultiple(&vContactRequest.VContact.VictimContactForm2, collectedErrors)

	// Si se han recopilado errores, se genera un nuevo captcha y se retorna el error.
	if len(collectedErrors) > 0 {
		// Se genera un nuevo ID de captcha.
		collectedErrors["default"]["captchaID"] = captcha.New()

		// Retorna un error de validación con código 400.
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Procede a crear el objeto VictimContact en la base de datos.
	if err = salvia_daos.SetVictimContact(&vContactRequest.VContact, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// Retorna error interno si ocurre algún fallo durante la creación.
		return http.StatusInternalServerError, err.Error()
	}

	vContactRequest.VContact.VictimContactForm2.VictimContactForm2VictimContact = vContactRequest.VContact

	//Ahora se verifican los enums

	//Tipo de reporte
	err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportType)
	if err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		// Retorna error interno en caso de fallo.
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReportType", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	if vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportType.VictimCaseForm2EnumsCode == "r" {
		//Significa que es un tercero reportando, entonces se verifican los campos correspondientes
		//Sólo con consultarlos se validan que existan. Y otros se trabaja con su valor

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vContactRequest.VContact.VictimContactForm2.VictimContactForm2VictimAwareOfReport)
		if err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			// Retorna error interno en caso de fallo.
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2VictimAwareOfReport", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vContactRequest.VContact.VictimContactForm2.VictimContactForm2HasCareRole)
		if err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			// Retorna error interno en caso de fallo.
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2HasCareRole", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}

		if vContactRequest.VContact.VictimContactForm2.VictimContactForm2HasCareRole.VictimCaseForm2EnumsCode == "y" || vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportTypeDetails.VictimCaseForm2EnumsCode == "ae" || vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportTypeDetails.VictimCaseForm2EnumsCode == "et" || vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportTypeDetails.VictimCaseForm2EnumsCode == "ex" || vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportTypeDetails.VictimCaseForm2EnumsCode == "ot" {
			//Significa que la víctima tiene rol de cuidado, entonces se verifican 2 preguntas más
			size := len(vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReporterNames)
			if size == 0 {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReporterNames", "json"), "common_validation_field_empty_error", "common_global_error", common_config.Locale)
			} else if int64(size) < salvia_daos.VictimContactForm2FieldDefinitions["VictimContactForm2ReporterNames"].MinSize {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReporterNames", "json"), "common_validation_field_number_min_size_error", "common_global_error", common_config.Locale)
			} else if int64(size) > salvia_daos.VictimContactForm2FieldDefinitions["VictimContactForm2ReporterNames"].MaxSize {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReporterNames", "json"), "common_validation_field_number_max_size_error", "common_global_error", common_config.Locale)
			}

			sizePhone := int64(vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReporterPhone)
			if sizePhone == 0 {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReporterPhone", "json"), "common_validation_field_number_error", "common_global_error", common_config.Locale)
			} else if sizePhone < salvia_daos.VictimContactForm2FieldDefinitions["VictimContactForm2ReporterPhone"].MinSize {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReporterPhone", "json"), "common_validation_field_number_min_size_error", "common_global_error", common_config.Locale)

			} else if sizePhone > salvia_daos.VictimContactForm2FieldDefinitions["VictimContactForm2ReporterPhone"].MaxSize {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2ReporterPhone", "json"), "common_validation_field_number_max_size_error", "common_global_error", common_config.Locale)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
			if len(collectedErrors) > 0 {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
			}
		}
	} else if vContactRequest.VContact.VictimContactForm2.VictimContactForm2ReportType.VictimCaseForm2EnumsCode == "v" {
		//Revisamos los campos adicionales para cuando es víctima

		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vContactRequest.VContact.VictimContactForm2.VictimContactForm2WillReceiveCall)
		if err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			// Retorna error interno en caso de fallo.
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(&vContactRequest.VContact.VictimContactForm2, "VictimContactForm2WillReceiveCall", "json"), "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
		}
	}

	if err = salvia_daos.SetVictimContactForm2(&vContactRequest.VContact.VictimContactForm2, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		// Retorna error interno si ocurre algún fallo durante la creación.
		return http.StatusInternalServerError, err.Error()
	}

	//Ahora se insertan todos los campos múltiples
	//Ahora se verifican los enums múltiples que se asociarán al form2
	if err = setVictimContactEnumsMultiple(&vContactRequest.VContact.VictimContactForm2, connData, dbClientConfig, dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Retorna éxito con código 200 y el objeto VictimContact en formato JSON.
	return http.StatusOK, ""
}

// GetVictimContactByICode consulta un registro de VictimContact utilizando el ICode (identificador único).
// Recibe como parámetros el ICode, el módulo actual, y la información de conexión y configuraciones de la base de datos.
// Retorna un código HTTP, un mensaje (JSON) y el DTO de VictimContact encontrado.
func GetVictimContactByICode(id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.VictimContactDTO) {
	var err error = nil
	// Libera la conexión si no se proporcionó un ID de conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO para VictimContact.
	var victimContact salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}
	var victimContactForm1 salvia_daos.VictimContactForm1DTO = salvia_daos.VictimContactForm1DTO{}
	var victimContactForm2 salvia_daos.VictimContactForm2DTO = salvia_daos.VictimContactForm2DTO{}

	var enums []salvia_daos.VictimCaseForm2EnumsDTO
	// Valida que se haya proporcionado un ID.
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, victimContact
		}

	} else {
		// Construye el criterio de búsqueda utilizando el ICode.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"VictimContactICode"},
			AttrsValue: []interface{}{id},
		}

		// Consulta el registro en la base de datos.
		err = salvia_daos.GetVictimContact(by, &victimContact, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Retorna error interno en caso de fallo.
			return http.StatusInternalServerError, err.Error(), victimContact
		}

		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"VictimContactForm1VictimContact"},
			AttrsValue: []interface{}{victimContact.VictimContactId},
		}

		if err = salvia_daos.GetVictimContactForm1(by, &victimContactForm1, connData, &dbClientConfig, &dbServerConfig); err != nil {
			by = common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"VictimContactForm2VictimContact"},
				AttrsValue: []interface{}{victimContact.VictimContactId},
			}

			if err = salvia_daos.GetVictimContactForm2(by, &victimContactForm2, connData, &dbClientConfig, &dbServerConfig); err != nil {
				return http.StatusInternalServerError, err.Error(), salvia_daos.VictimContactDTO{}
			} else {
				victimContact.VictimContactForm2 = victimContactForm2
			}
		} else {
			victimContact.VictimContactForm1 = victimContactForm1
		}
		if victimContact.VictimContactForm2.VictimContactForm2Id > 0 {
			//Traemos todos los enums asociados a este formulario
			enums, err = salvia_daos.GetVictimCasesForm2EnumsByVictimContactForm2Id(victimContact.VictimContactForm2.VictimContactForm2Id, connData, &dbClientConfig, &dbServerConfig)
			if err != nil {
				return http.StatusInternalServerError, "", salvia_daos.VictimContactDTO{}
			}
			loadVictimContactEnumsMultiple(&victimContact.VictimContactForm2, enums)
		}

		// Retorna éxito con el DTO actualizado.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(victimContact), victimContact

	}
	// En caso de error no identificado, se retorna un error interno.
	return http.StatusInternalServerError, "", victimContact
}

func loadVictimContactEnumsMultiple(vContactForm2 *salvia_daos.VictimContactForm2DTO, enums []salvia_daos.VictimCaseForm2EnumsDTO) {
	var err error

	for _, e := range enums {
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsById(&e); err == nil {
			switch e.VictimCaseForm2EnumsCategory {
			case "victim_case_form2_adjustments_gbv":
				vContactForm2.VictimContactForm2AdjustmentsGBV = append(vContactForm2.VictimContactForm2AdjustmentsGBV, e)
			}
		}
	}
}

// InvalidateVictimContactByICode invalida un registro de VictimContact identificado por su ICode.
// Recibe como parámetros el JSON de entrada (dataInput) que contiene información adicional, el ICode, el módulo,
// y la información de conexión y configuraciones de la base de datos.
// Retorna un código HTTP y un mensaje en formato JSON.
func InvalidateVictimContactByICode(dataInput string, id string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Verifica que se haya proporcionado un ICode.
	if id == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	var err error = nil

	// Mapa para recopilar errores durante la validación del DTO.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Libera la conexión si no se proporcionó un ID de conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crean dos instancias del DTO para actualizar:
	// vContactToUpdate: registro existente a modificar.
	// vContact: registro con los nuevos datos.
	var vContactToUpdate salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}
	var vContact salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}

	// Definición de campos requeridos para la actualización.
	var checkFields map[string]bool = map[string]bool{"VictimContactStatusDescription": true}

	// Se obtiene el DTO a partir del JSON de entrada.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.VictimContactJSONName, common_config.Locale, collectedErrors)

	// Valida la estructura del DTO.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		var dtoMap map[string]interface{} = rawDto
		utils.ValidateJSONInput(&vContact, dtoMap, salvia_daos.VictimContactJSONName, salvia_daos.VictimContactFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
	default:
		// Registra un error si la estructura no es la esperada.
		utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si existen errores de validación, se retorna un error con el detalle.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se construye el criterio de búsqueda para identificar el registro a actualizar.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimContactICode"},
		AttrsValue: []interface{}{id},
	}

	// Consulta el registro existente a actualizar.
	if err = salvia_daos.GetVictimContact(by, &vContactToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Actualiza el estado del registro a "invalido" y asigna la descripción del estado.
	vContactToUpdate.VictimContactStatus = "i"
	vContactToUpdate.VictimContactStatusDescription = vContact.VictimContactStatusDescription

	// Se actualiza el registro en la base de datos.
	if err = salvia_daos.InvalidateVictimContactByICode(&vContactToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Retorna éxito con código 200.
	return http.StatusOK, ""
}

// GetVictimContactsWithoutVictimCase obtiene una lista de registros de VictimContact que no están asociados a ningún caso.
// Recibe como parámetros el módulo actual y la información de conexión y configuraciones de la base de datos.
// Retorna un código HTTP y un mensaje en formato JSON con la lista de registros o un mensaje de error.
func GetVictimContactsWithoutVictimCase(page int, status string, filters salvia_daos.VictimContactListFilters, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, int) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError
	var count int

	// Libera la conexión si no se proporcionó un ID de conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Consulta los registros sin caso asociado.
	contacts, count, err := salvia_daos.GetVictimContactsWithoutVictimCase(page, status, filters, connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// Retorna error interno en caso de fallo.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Retorna la lista de contactos en formato JSON con éxito.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccessMultiple([]string{"data", "numPages"}, []interface{}{contacts, int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE)))}, []string{"[]", "0"})
	}
	return resCode, resData, count
}

// GetVictimContactByAll obtiene todos los registros de VictimContact de la base de datos.
// Recibe como parámetros el módulo actual y la información de conexión y configuraciones de la base de datos.
// Retorna un código HTTP y un mensaje en formato JSON con la lista completa de registros o un mensaje de error.
func GetVictimContactByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Libera la conexión si no se proporcionó un ID de conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Consulta todos los registros de VictimContact.
	contacts, err := salvia_daos.GetAllVictimContact(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// Retorna error interno en caso de fallo.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Retorna la lista completa de contactos en formato JSON con éxito.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(contacts)
	}
	return resCode, resData
}

func getAndVerifyVictimContactEnumsMultiple(vContactForm2 *salvia_daos.VictimContactForm2DTO, collectedErrors map[string]map[string]string) bool {
	var opRes bool = true
	var err error

	for idx := range vContactForm2.VictimContactForm2AdjustmentsGBV {
		err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&vContactForm2.VictimContactForm2AdjustmentsGBV[idx])
		if err != nil {
			opRes = false
			break
		}
	}
	if len(vContactForm2.VictimContactForm2AdjustmentsGBV) == 0 || err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, utils.GetTag(vContactForm2, "VictimContactForm2AdjustmentsGBV", "json"), "victim_case_form2_enums_multiple_not_found", "common_global_error", salvia_config.Locale)
		opRes = false
	}

	return opRes
}

func setVictimContactEnumsMultiple(vContactForm2 *salvia_daos.VictimContactForm2DTO, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) error {
	var err error
	for _, e := range vContactForm2.VictimContactForm2AdjustmentsGBV {
		var rel salvia_daos.RelVictimCaseForm2EnumsVictimContactForm2DTO
		rel.RelVictimCaseForm2EnumsVictimContactForm2Form = *vContactForm2
		rel.RelVictimCaseForm2EnumsVictimCaseForm2Enums = e
		if err = salvia_daos.SetRelVictimCaseForm2EnumsVictimContactForm2(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return err
		}
	}

	return nil
}
