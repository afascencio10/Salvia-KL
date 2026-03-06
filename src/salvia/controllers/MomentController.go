// Package salvia_ctrl proporciona funciones de control para gestionar "moments" (momentos)
// en el módulo Salvia, incluyendo la obtención y actualización de momentos asociados a casos de víctimas y sucursales de entidades.
package salvia_ctrl

import (
	// Importación de configuraciones, controladores, acceso a datos y utilidades comunes.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
	"mime/multipart"
	"net/http"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// GetMomentByICode
///////////////////////////////////////////////////////////////////////////////

// GetMomentByICode obtiene un "moment" (momento) a partir de su código identificador (ICode).
//
// Parámetros:
//   - icode: Código identificador del momento.
//   - module: Nombre del módulo desde donde se invoca la función.
//   - connData: Datos de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de base de datos.
//   - dbServerConfig: Configuración del servidor de base de datos.
//
// Retorna:
//   - Código HTTP resultante.
//   - Mensaje de respuesta (éxito o error).
//   - Estructura MomentDTO con los detalles del momento obtenido.
func GetMomentByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.MomentDTO) {
	var err error = nil

	// Si no se especifica un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO requerido para almacenar la información del momento.
	var branch salvia_daos.MomentDTO = salvia_daos.MomentDTO{}
	if icode == "" {
		// Si el id está vacío, se retorna un error utilizando el mensaje localizado.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, branch
		}

	} else {
		// Se crea una condición de búsqueda (By) para localizar el momento por "MomentICode".
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"MomentICode"},
			AttrsValue: []interface{}{icode},
		}

		// Se obtiene el momento desde la base de datos.
		err = salvia_daos.GetMoment(by, &branch, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// Se retorna el error ocurrido durante la obtención.
			return http.StatusInternalServerError, err.Error(), branch
		}

		// Se retorna el éxito con un mensaje formateado y el DTO obtenido.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(branch), branch

	}
	return http.StatusInternalServerError, "", branch
}

///////////////////////////////////////////////////////////////////////////////
// GetMomentByAll
///////////////////////////////////////////////////////////////////////////////

// GetMomentByAll obtiene todos los momentos disponibles para un módulo determinado.
//
// Parámetros:
//   - module: Nombre del módulo desde donde se invoca la función.
//   - connData: Datos de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de base de datos.
//   - dbServerConfig: Configuración del servidor de base de datos.
//
// Retorna:
//   - Código HTTP resultante.
//   - Mensaje de respuesta (éxito o error).
//   - Slice de MomentDTO con todos los momentos obtenidos.
func GetMomentByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.MomentDTO) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se especifica un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se obtienen todos los momentos.
	ebranches, err := salvia_daos.GetAllMoment(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		// Si ocurre un error, se asigna el mensaje de error.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Si la operación es exitosa, se retorna el mensaje de éxito junto con los datos.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(ebranches)
	}
	return resCode, resData, ebranches
}

///////////////////////////////////////////////////////////////////////////////
// UpdateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch
///////////////////////////////////////////////////////////////////////////////

// UpdateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch actualiza un momento en la base de datos
// utilizando el código del caso de la víctima, el código del momento y el código de la entidad (sucursal).
//
// Parámetros:
//   - dataInput: Cadena JSON de entrada con la información a actualizar.
//   - sessionUser: Información de la sesión del usuario que realiza la actualización.
//   - vCaseIcode: Código identificador del caso de la víctima.
//   - momentCode: Código del momento a actualizar.
//   - eBranchIcode: Código identificador de la sucursal de la entidad.
//   - module: Nombre del módulo desde donde se invoca la función.
//   - connData: Datos de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de base de datos.
//   - dbServerConfig: Configuración del servidor de base de datos.
//
// Retorna:
//   - Código HTTP resultante.
//   - Mensaje de respuesta (éxito o error).
func UpdateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch(dataInput string, sessionUser utils.CommonSession, vCaseIcode string, momentCode string, eBranchIcode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Validación de que los campos requeridos no estén vacíos.
	if vCaseIcode == "" || momentCode == "" || eBranchIcode == "" {
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v
		}
	}

	// Mapa para recolectar errores durante la validación del DTO.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crea un objeto MomentDTO y se obtiene el DTO a partir del JSON de entrada.
	var moment salvia_daos.MomentDTO = salvia_daos.MomentDTO{}
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.MomentJSONName, common_config.Locale, collectedErrors)

	// Se define un mapa de campos que se deben verificar, en este caso la descripción de aprobación.
	var checkFields map[string]bool = map[string]bool{"MomentApprovalDescription": true}

	// Se valida el JSON de entrada según las definiciones de campos del DTO.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		var dtoMap map[string]interface{} = rawDto
		utils.ValidateJSONInput(&moment, dtoMap, salvia_daos.MomentJSONName, salvia_daos.MomentFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)
	default:
		// Si el DTO no es del tipo esperado, se registra un error.
		utils.SetError(collectedErrors, salvia_daos.MomentJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", salvia_config.Locale)
	}

	// Si se encontraron errores durante la validación, se retorna el error.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se llama a la función interna para realizar la actualización del momento.
	code, res := updateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch(moment, sessionUser, vCaseIcode, momentCode, eBranchIcode, connData, dbClientConfig, dbServerConfig)

	return code, res
}

///////////////////////////////////////////////////////////////////////////////
// updateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch (función interna)
///////////////////////////////////////////////////////////////////////////////

// updateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch realiza la actualización de un momento
// en la base de datos, validando y obteniendo los datos relacionados (caso de víctima, sucursal de entidad y propietario de aprobación),
// y actualizando el estado del momento y, de ser necesario, el estado del caso.
//
// Parámetros:
//   - moment: Objeto MomentDTO con la información actualizada.
//   - sessionUser: Información de la sesión del usuario que realiza la actualización.
//   - vCaseIcode: Código identificador del caso de la víctima.
//   - momentCode: Código del momento a actualizar.
//   - eBranchIcode: Código identificador de la sucursal de la entidad.
//   - module: Nombre del módulo desde donde se invoca la función.
//   - connData: Datos de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de base de datos.
//   - dbServerConfig: Configuración del servidor de base de datos.
//
// Retorna:
//   - Código HTTP resultante.
//   - Mensaje de respuesta (éxito o error).
func updateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch(moment salvia_daos.MomentDTO, sessionUser utils.CommonSession, vCaseIcode string, momentCode string, eBranchIcode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Si no se especifica un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Declaración de DTOs requeridos para la actualización.
	var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	var eBranch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
	var moments []salvia_daos.MomentDTO = []salvia_daos.MomentDTO{}
	var momentToUpdate salvia_daos.MomentDTO = salvia_daos.MomentDTO{}
	var approvalOwner salvia_daos.CaseOwnerDTO = salvia_daos.CaseOwnerDTO{}

	var err error = nil

	// Se obtiene el caso de víctima basado en el vCaseIcode.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"VictimCaseICode"},
		AttrsValue: []interface{}{vCaseIcode},
	}

	if err = salvia_daos.GetVictimCase(by, &vCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtiene la sucursal de la entidad basada en el eBranchIcode.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EntityBranchICode"},
		AttrsValue: []interface{}{eBranchIcode},
	}

	if err = salvia_daos.GetEntityBranch(by, &eBranch, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtiene el propietario de la aprobación del caso basado en el usuario de sesión.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"CaseOwnerGeneralUser"},
		AttrsValue: []interface{}{sessionUser.UserICode},
	}

	if err = salvia_daos.GetCaseOwner(by, &approvalOwner, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se obtiene el momento a actualizar utilizando múltiples criterios: caso, código del momento y sucursal.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentVictimCase", "MomentCode", "MomentEntityBranch"},
		AttrsValue: []interface{}{vCase.VictimCaseId, momentCode, eBranch.EntityBranchId},
	}

	if err = salvia_daos.GetMoment(by, &momentToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Se establecen los valores predeterminados para la actualización del momento.
	salvia_daos.SetMomentDefaults(&moment, common_dao.SQL_UPDATE)

	// Se actualizan los campos relevantes del momento.
	momentToUpdate.MomentApprovalDescription = moment.MomentApprovalDescription
	momentToUpdate.MomentApprovalSource = "m" // Temporal mientras se automatiza. Por ahora se asume que todas las aprobaciones son manuales.
	momentToUpdate.MomentUpdateDate = moment.MomentUpdateDate
	momentToUpdate.MomentApprovalOwner = approvalOwner
	momentToUpdate.MomentStatus = "a"

	// Se inicia una transacción en la base de datos.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	// Se actualiza el momento en la base de datos.
	if err = salvia_daos.UpdateMomentById(&momentToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// En caso de error, se revierte la transacción.
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se verifica si todos los momentos del caso han sido aprobados.
	by = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"MomentVictimCase"},
		AttrsValue: []interface{}{vCase.VictimCaseId},
	}

	if moments, err = salvia_daos.GetMoments(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	// Se recorre la lista de momentos para determinar si el caso está completo.
	var caseComplete bool = true
	for _, m := range moments {
		if m.MomentStatus == "u" {
			caseComplete = false
			break
		}
	}

	// Si el caso está completo, se actualiza el estado del caso a "c".
	if caseComplete {
		vCase.VictimCaseStatus = "c"
		if err = salvia_daos.UpdateVictimCaseStatus(&vCase, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, err.Error()
		}
	}
	// Se confirma la transacción.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}

///////////////////////////////////////////////////////////////////////////////
// UpdateMomentByPlainFile
///////////////////////////////////////////////////////////////////////////////

// UpdateMomentByPlainFile procesa un archivo de texto plano (CSV) para actualizar momentos en la base de datos.
// El archivo debe tener un formato específico en el que cada línea (excepto la cabecera) contiene 5 campos:
//
//	[0] Tipo de documento
//	[1] Número de documento
//	[2] Código de entidad
//	[3] Código del caso (fecha y hora de creación)
//	[4] Momento de la actuación
//
// Parámetros:
//   - file: Archivo multipart que contiene los datos en formato plano.
//   - sessionUser: Información de la sesión del usuario que realiza la operación.
//   - module: Nombre del módulo desde donde se invoca la función.
//   - connData: Datos de la conexión a la base de datos.
//   - dbClientConfig: Configuración del cliente de base de datos.
//   - dbServerConfig: Configuración del servidor de base de datos.
//
// Retorna:
//   - Código HTTP resultante.
//   - Mensaje de respuesta (éxito o error).
func UpdateMomentByPlainFile(file *multipart.FileHeader, sessionUser utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Si no se especifica un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se lee el contenido del archivo en una cadena.
	fileStr, err := utils.ReadStringFileHeader(file)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	/*
		[0] Tipo de documento
		[1] Número de documento
		[2] Código de entidad
		[3] Código del caso
		[4] Momento de la actuación
	*/
	// Se normalizan los saltos de línea.
	fileStr = strings.ReplaceAll(fileStr, "\r\n", "\n")

	// Se separa el contenido en líneas.
	lines := strings.Split(fileStr, "\n")

	// Se procesa cada línea del archivo (omitiendo la cabecera).
	for idx, l := range lines {
		if idx > 0 {
			// Se separan los campos por coma.
			fields := strings.Split(l, ",")
			if len(fields) == 5 {
				// Se asignan los valores correspondientes a cada campo.
				var (
					docType    string = fields[0]
					numDoc     string = fields[1]
					entityCode string = fields[2]
					caseCode   string = fields[3]
					momentCode string = fields[4]
				)

				// Se convierten los códigos de momentos globales a los internos del sistema.
				switch momentCode {
				case "03":
					momentCode = "01"
				case "04":
					momentCode = "02"
				case "05":
					momentCode = "03"
				}

				/*
					Se obtiene el caso de víctima utilizando el código externo (que contiene la fecha y hora de creación)
					y los datos del documento de la víctima.
				*/
				var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
				var code int
				var res string
				code, _, vCase = GetVictimCaseByDocumentAndCreationDate(docType, numDoc, caseCode, connData, dbClientConfig, dbServerConfig)
				if code != http.StatusOK {
					// TODO: Reportar error en la línea correspondiente.
					return http.StatusInternalServerError, ""
				}
				// Se obtiene la sucursal de la entidad a partir del código de entidad.
				var entityBranch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
				code, _, entityBranch = GetEntityBranchByICode(entityCode, connData, dbClientConfig, dbServerConfig)
				if code != http.StatusOK {
					// TODO: Reportar error en la línea correspondiente.
					return http.StatusInternalServerError, ""
				}
				// Se crea un objeto MomentDTO con una descripción por defecto, indicando que fue cargado sin comentarios.
				var moment salvia_daos.MomentDTO = salvia_daos.MomentDTO{MomentApprovalDescription: "Cargado desde archivo plano sin comentarios"}
				code, res = updateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch(moment, sessionUser, vCase.VictimCaseICode, momentCode, entityBranch.EntityBranchICode, connData, dbClientConfig, dbServerConfig)
				if code != http.StatusOK {
					// TODO: Reportar error en la línea correspondiente.
					return http.StatusInternalServerError, res
				}
			} else {
				// TODO: Indicar error en la línea correspondiente si el número de campos no es el esperado.
			}
		}
	}

	return http.StatusOK, ""
}
