// Package salvia_ctrl contiene los controladores para la gestión de sedes de entidad,
// incluyendo funciones para obtener, actualizar y crear sedes a partir de diferentes criterios.
// Se apoya en las capas DAO y de configuración para interactuar con la base de datos y manejar mensajes de respuesta.
package salvia_ctrl

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type EntityBranchRequest struct {
	Branch salvia_daos.EntityBranchDTO `json:"entityBranch"`
}

// GetEntityBranchByICode recupera una sede de entidad (EntityBranch) basándose en su ICode.
// Parámetros:
//   - id: código identificador de la sede.
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la estructura de conexión a la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta (error o mensaje de éxito en formato JSON).
//   - Estructura EntityBranchDTO con la información de la sede.
func GetEntityBranchByICode(icode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, salvia_daos.EntityBranchDTO) {
	var err error = nil

	// Si no se tiene un identificador de conexión, se libera la conexión al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea la estructura DTO para la sede.
	var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
	if icode == "" {
		// Si no se proporciona un id, se retorna un error de petición incorrecta.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, branch
		}

	} else {
		// Se define el criterio de búsqueda para la sede utilizando el ICode.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityBranchICode"},
			AttrsValue: []interface{}{icode},
		}

		// Se invoca el DAO para obtener la sede.
		err = salvia_daos.GetEntityBranch(by, &branch, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, se retorna el error obtenido.
			return http.StatusInternalServerError, err.Error(), branch
		}

		// Se retorna la sede encontrada junto con un mensaje de éxito.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(branch), branch

	}
	return http.StatusInternalServerError, "", branch
}

// GetEntityBranchesBySector recupera una sede de entidad basándose en el sector proporcionado.
// Parámetros:
//   - sector: código del sector a buscar.
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la conexión con la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta.
//   - Slice de EntityBranchDTO; aunque en este caso se retorna un slice vacío.
func GetEntityBranchesBySector(sector string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityBranchDTO) {
	var err error = nil

	// Si la conexión no tiene un ID, se libera al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Se crea el DTO requerido.
	var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
	if sector == "" {
		// Si no se proporciona el sector, se retorna un error.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, []salvia_daos.EntityBranchDTO{}
		}

	} else {
		// Se construye el criterio de búsqueda usando el sector como EntityBranchICode.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityBranchICode"},
			AttrsValue: []interface{}{sector},
		}

		// Se intenta obtener la sede a través del DAO.
		err = salvia_daos.GetEntityBranch(by, &branch, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			// En caso de error, se retorna el mensaje correspondiente.
			return http.StatusInternalServerError, err.Error(), []salvia_daos.EntityBranchDTO{}
		}

		// Se retorna un mensaje de éxito y, en este caso, se retorna un slice vacío.
		return http.StatusOK, utils.CommMsgGetJSONSuccess(branch), []salvia_daos.EntityBranchDTO{}
	}
	return http.StatusInternalServerError, "", []salvia_daos.EntityBranchDTO{}
}

// GetEntityBranchesByTownCodeWithMoments recupera las sedes de entidad asociadas a un código de ciudad
// incluyendo información de "moments".
// Parámetros:
//   - townCode: código de la ciudad.
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la conexión con la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta.
//   - Slice de EntityBranchDTO con la información de las sedes.
func GetEntityBranchesByTownCodeWithMoments(townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityBranchDTO) {
	var err error
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no tiene un ID, se libera al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Inicializa el slice para las sedes.
	var entityBranches []salvia_daos.EntityBranchDTO = []salvia_daos.EntityBranchDTO{}
	if townCode == "" {
		// Retorna error si no se proporciona el código de la ciudad.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, entityBranches
		}
	} else {
		// Llama al DAO para obtener las sedes con información de moments.
		entityBranches, err = salvia_daos.GetEntityBranchesByTownCodeWithMoments(townCode, connData, &dbClientConfig, &dbServerConfig)

		if err != nil {
			// En caso de error, se prepara la respuesta de error.
			resData = err.Error()
			resCode = http.StatusInternalServerError
		} else {
			// Si la operación fue exitosa, se retorna un mensaje de éxito.
			resCode = http.StatusOK
			resData = utils.CommMsgGetJSONSuccess(entityBranches)
		}
	}
	return resCode, resData, entityBranches
}

// GetEntityBranchesByTownCodeAndEntityBranchIcodeWithMoments recupera las sedes de entidad filtradas por
// el código de ciudad y el ICode de la sede, incluyendo información de "moments".
// Parámetros:
//   - townCode: código de la ciudad.
//   - entityBranchICode: ICode de la sede.
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la conexión con la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta.
//   - Slice de EntityBranchDTO con la información de las sedes.
func GetEntityBranchesByTownCodeAndEntityBranchIcodeWithMoments(townCode string, entityBranchICode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityBranchDTO) {
	var err error
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no existe un identificador de conexión, se libera al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Inicializa el slice para almacenar las sedes.
	var entityBranches []salvia_daos.EntityBranchDTO = []salvia_daos.EntityBranchDTO{}
	if townCode == "" || entityBranchICode == "" {
		// Retorna error si faltan parámetros requeridos.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, entityBranches
		}
	} else {
		// Llama al DAO para obtener las sedes con la combinación de townCode e entityBranchICode.
		entityBranches, err = salvia_daos.GetEntityBranchesByTownCodeAndEntityBranchIcodeWithMoments(townCode, entityBranchICode, connData, &dbClientConfig, &dbServerConfig)

		if err != nil {
			// En caso de error, se prepara la respuesta con el mensaje de error.
			resData = err.Error()
			resCode = http.StatusInternalServerError
		} else {
			// Operación exitosa.
			resCode = http.StatusOK
			resData = utils.CommMsgGetJSONSuccess(entityBranches)
		}
	}
	return resCode, resData, entityBranches
}

// GetEntityBranchesByTownCodeAndEntityICode recupera las sedes de entidad basándose en el código
// de ciudad y el ICode de la entidad asociada.
// Parámetros:
//   - entityICode: código identificador de la entidad.
//   - townCode: código de la ciudad.
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la conexión con la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta.
//   - Slice de EntityBranchDTO con la información de las sedes encontradas.
func GetEntityBranchesByTownCodeAndEntityICode(entityICode string, townCode string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityBranchDTO) {
	var err error
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si no se tiene un ID de conexión, se libera al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}
	// Inicializa los slices y estructuras necesarios.
	var entityBranches []salvia_daos.EntityBranchDTO = []salvia_daos.EntityBranchDTO{}
	var entity salvia_daos.EntityDTO = salvia_daos.EntityDTO{}

	if townCode == "" {
		// Retorna error si falta el código de ciudad.
		if v, found := common_config.Locale["sp"][common_config.Enums.GLOBAL_ERROR]; found {
			return http.StatusBadRequest, v, entityBranches
		}
	} else {
		// Define el criterio para buscar la entidad a partir del entityICode.
		var by common_controllers.By = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityICode"},
			AttrsValue: []interface{}{entityICode},
		}

		// Obtiene la entidad a través del DAO.
		err = salvia_daos.GetEntity(by, &entity, connData, &dbClientConfig, &dbServerConfig)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), entityBranches
		}

		// Construye el criterio para obtener las sedes asociadas a la entidad y ciudad.
		by = common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityBranchTownCode", "EntityBranchEntity"},
			AttrsValue: []interface{}{townCode, entity.EntityId},
		}
		entityBranches, err = salvia_daos.GetEntityBranches(by, connData, &dbClientConfig, &dbServerConfig)

		if err != nil {
			return http.StatusInternalServerError, err.Error(), entityBranches
		} else {
			resCode = http.StatusOK
			resData = utils.CommMsgGetJSONSuccess(entityBranches)
		}
	}
	return resCode, resData, entityBranches
}

// GetEntityBranchByAll recupera todas las sedes de entidad existentes.
// Parámetros:
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la conexión con la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta.
//   - Slice de EntityBranchDTO con todas las sedes.
func GetEntityBranchByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string, []salvia_daos.EntityBranchDTO) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Si la conexión no tiene ID, se libera al finalizar.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Llama al DAO para obtener todas las sedes.
	ebranches, err := salvia_daos.GetAllEntityBranch(connData, &dbClientConfig, &dbServerConfig)

	if err != nil {
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(ebranches)
	}
	return resCode, resData, ebranches
}

// SetUpdateEntityBranchByPlainFile procesa un archivo de texto plano para crear o actualizar una sede.
// La función lee el archivo, extrae los campos requeridos y, en función de si la sede ya existe,
// la crea o actualiza en la base de datos.
// Parámetros:
//   - file: puntero a multipart.FileHeader que contiene el archivo.
//   - entityICode: código identificador de la entidad asociada.
//   - sessionUser: información de la sesión del usuario (no se utiliza explícitamente en el código mostrado).
//   - module: nombre del módulo para contexto.
//   - connData: puntero a la conexión con la base de datos.
//   - dbClientConfig: configuración del cliente de base de datos.
//   - dbServerConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - Código de estado HTTP.
//   - Mensaje de respuesta (formateado en JSON).
func SetUpdateEntityBranchByPlainFile(file *multipart.FileHeader, entityICode string, sessionUser utils.CommonSession, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	// Si la conexión no tiene ID, se libera al finalizar y se configura la conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)

		// Se obtiene una conexión válida utilizando el controlador de persistencia.
		var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
		persistenceCtrl.Setup(connData, &dbClientConfig, &dbServerConfig)
	}

	var code int

	// Se lee el contenido del archivo en forma de cadena.
	fileStr, err := utils.ReadStringFileHeader(file)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	/*
		Formato esperado del archivo (separado por líneas):
		[0] Código de la sede
		[1] Nombre de la sede
		[2] Descripción
		[3] Dirección de la sede
		[4] Latitud GPS
		[5] Longitud GPS
		[6] Código Ciudad
	*/

	// Normaliza los saltos de línea.
	fileStr = strings.ReplaceAll(fileStr, "\r\n", "\n")

	// Separa el contenido del archivo en líneas.
	lines := strings.Split(fileStr, "\n")

	// Se inicializa la entidad asociada.
	var entity salvia_daos.EntityDTO = salvia_daos.EntityDTO{}

	// Se obtiene la entidad a partir de su ICode.
	code, _, entity = GetEntityByICode(entityICode, connData, dbClientConfig, dbServerConfig)
	if code != http.StatusOK {
		return http.StatusInternalServerError, `{"loadPlainFiles":["Error interno: No se encontró la entidad"]}`
	}

	// Verifica que el archivo contenga al menos una línea de contenido.
	if len(lines) < 2 {
		return http.StatusBadRequest, `{"loadPlainFiles":["El archivo debe contener al menos una fila de contenido."]}`
	}

	// Inicia una transacción para asegurar la integridad de las operaciones.
	if connData, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Itera sobre cada línea del archivo (se omite la primera línea que suele ser encabezado).
	for idx, l := range lines {
		if idx > 0 {
			// Se separan los campos usando "|" como delimitador.
			fields := strings.Split(l, "|")
			if len(fields) == 7 {
				// Se asignan los campos a variables locales para mayor claridad.
				var (
					branchCode        string = fields[0]
					branchName        string = fields[1]
					branchDescription string = fields[2]
					branchAddress     string = fields[3]
					branchLatitude    string = fields[4]
					branchLongitude   string = fields[5]
					branchTownCode    string = fields[6]
					floatTmp          float64
				)

				// Se verifica si la sede ya existe.
				var branch salvia_daos.EntityBranchDTO = salvia_daos.EntityBranchDTO{}
				branch.EntityBranchSource = "s"

				code, _, branch = GetEntityBranchByICode(branchCode, connData, dbClientConfig, dbServerConfig)
				if code != http.StatusOK {
					// La sede no existe, por lo que se crea.
					branch.EntityBranchEntity = entity
					salvia_daos.SetEntityBranchDefaults(&branch, common_dao.SQL_INSERT)
					branch.EntityBranchICode = branchCode
					branch.EntityBranchName = branchName
					branch.EntityBranchDescription = branchDescription
					branch.EntityBranchAddress = branchAddress
					// Conversión de latitud a float64.
					if floatTmp, err = strconv.ParseFloat(branchLatitude, 64); err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - formato de latitud incorrecto"]}`
					}
					branch.EntityBranchLatitude = floatTmp
					// Conversión de longitud a float64.
					if floatTmp, err = strconv.ParseFloat(branchLongitude, 64); err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - formato de longitud incorrecto"]}`
					}
					branch.EntityBranchLongitude = floatTmp
					// TODO: validar que el town code sea válido.
					branch.EntityBranchTownCode = branchTownCode

					// Se intenta crear la sede.
					err = salvia_daos.SetEntityBranch(&branch, connData, &dbClientConfig, &dbServerConfig)
					if err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - no es posible crear la sede"]}`
					}

				} else {
					// La sede existe, se procede a actualizarla.
					branch.EntityBranchEntity = entity
					salvia_daos.SetEntityBranchDefaults(&branch, common_dao.SQL_UPDATE)
					branch.EntityBranchName = branchName
					branch.EntityBranchDescription = branchDescription
					branch.EntityBranchAddress = branchAddress
					// Conversión de latitud a float64.
					if floatTmp, err = strconv.ParseFloat(branchLatitude, 64); err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - formato de latitud incorrecto"]}`
					}
					branch.EntityBranchLatitude = floatTmp
					// Conversión de longitud a float64.
					if floatTmp, err = strconv.ParseFloat(branchLongitude, 64); err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - formato de longitud incorrecto"]}`
					}
					branch.EntityBranchLongitude = floatTmp
					// TODO: validar que el town code sea válido.
					branch.EntityBranchTownCode = branchTownCode

					// Se intenta actualizar la sede.
					err = salvia_daos.UpdateEntityBranch(&branch, connData, &dbClientConfig, &dbServerConfig)
					if err != nil {
						db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
						return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - no es posible actualizar la sede"]}`
					}
				}

			} else if idx < len(lines)-1 {
				// Si el número de columnas es incorrecto, se revierte la transacción.
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusBadRequest, `{"loadPlainFiles":["Línea: ` + strconv.Itoa(idx) + ` - la fila no contiene la cantidad de columnas requeridas"]}`
			}
		}
	}

	// Se confirma la transacción si todas las operaciones fueron exitosas.
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, ""
}

// SetEntityBranch crea un registro de entidad basado en los datos de entrada proporcionados.
// Realiza la validación del JSON de entrada, establece valores por defecto, asocia el log a una entidad y municipio

// Parámetros:
// - dataInput: Cadena JSON que contiene los datos de la sede
// - module: Identificador del módulo que invoca la función.
// - s: Sesión del usuario que incluye el código de usuario y otros datos de sesión.
// - dbClientConfig: Configuración del cliente de la base de datos.
// - dbServerConfig: Configuración del servidor de la base de datos.
//
// Retorna:
// - Código de estado HTTP que indica el resultado de la operación.
// - Cadena JSON con los datos del log creado en caso de éxito o mensajes de error.
func SetEntityBranchByUser(dataInput string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable de error.
	var err error = nil

	// Se obtiene la conexión a la base de datos.
	var connData *db.ConnData = &db.ConnData{}
	// Se garantiza la liberación de la conexión al finalizar la función.
	defer db.ReleaseConnection(connData)

	// Mapa para recolectar errores durante el procesamiento.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	var entityBranchRequest EntityBranchRequest

	var dtoMap map[string]interface{} = nil

	// Se establecen los valores por defecto para el DTO del log de caso (operación de inserción SQL).
	salvia_daos.SetEntityBranchDefaults(&entityBranchRequest.Branch, common_dao.SQL_INSERT)

	// Definición de campos requeridos para la validación del JSON (por ejemplo, "EntityBranchDescription").
	var checkFields map[string]bool = map[string]bool{"EntityBranchName": true, "EntityBranchDescription": false, "EntityBranchAddress": true, "EntityBranchLatitude": true, "EntityBranchLongitude": true,
		"EntityBranchEntityICode": true, "EntityBranchTownCode": true}

	// Se parsea la cadena JSON de entrada en una estructura de mapa.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.EntityBranchJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &entityBranchRequest)

	// Validación de la estructura del DTO obtenido.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Se valida la entrada JSON contra las definiciones esperadas de los campos.
		utils.ValidateJSONInput(&entityBranchRequest.Branch, dtoMap, salvia_daos.EntityBranchJSONName, salvia_daos.EntityBranchFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		entityBranchRequest.Branch.EntityBranchEntity = salvia_daos.EntityDTO{EntityICode: entityBranchRequest.Branch.EntityBranchEntityICode}

		//--------------------------------------------------------------------
	default:
		// Si la estructura no es la esperada, se registra un error global.
		utils.SetError(collectedErrors, salvia_daos.EntityBranchJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si se han detectado errores durante la validación, se retorna un error 400 con los detalles en formato JSON.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se verifica la entidad.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EntityICode"},
		AttrsValue: []interface{}{entityBranchRequest.Branch.EntityBranchEntity.EntityICode},
	}

	if err = salvia_daos.GetEntity(by, &entityBranchRequest.Branch.EntityBranchEntity, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	entityBranchRequest.Branch.EntityBranchSource = "u"

	// Se procede a crear el registro del log en la base de datos.
	if err = salvia_daos.SetEntityBranch(&entityBranchRequest.Branch, connData, &dbClientConfig, &dbServerConfig); err != nil {
		// En caso de error durante la creación, se retorna un error 500 con el mensaje correspondiente.
		return http.StatusInternalServerError, err.Error()
	}

	// Se retorna el estado 200 OK con el DTO del log de caso en formato JSON.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(entityBranchRequest.Branch)
}

func UpdateEntityBranchByUser(dataInput string, s utils.CommonSession, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Inicialización de la variable de error.
	var err error = nil

	// Se obtiene la conexión a la base de datos.
	var connData *db.ConnData = &db.ConnData{}
	// Se garantiza la liberación de la conexión al finalizar la función.
	defer db.ReleaseConnection(connData)

	// Mapa para recolectar errores durante el procesamiento.
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	// Se crean los DTO (Data Transfer Object) requeridos para el log de caso y el momento.
	var entityBranchRequest EntityBranchRequest
	var entityBranchToUpdate salvia_daos.EntityBranchDTO

	var dtoMap map[string]interface{} = nil

	// Se establecen los valores por defecto para el DTO del log de caso (operación de inserción SQL).
	salvia_daos.SetEntityBranchDefaults(&entityBranchRequest.Branch, common_dao.SQL_UPDATE)

	// Definición de campos requeridos para la validación del JSON (por ejemplo, "EntityBranchDescription").
	var checkFields map[string]bool = map[string]bool{"EntityBranchName": true, "EntityBranchDescription": false, "EntityBranchAddress": true, "EntityBranchLatitude": true, "EntityBranchLongitude": true}

	// Se parsea la cadena JSON de entrada en una estructura de mapa.
	var rawDto interface{} = utils.GetDTOMap(dataInput, salvia_daos.EntityBranchJSONName, common_config.Locale, collectedErrors)

	utils.JSONToStruct(dataInput, &entityBranchRequest)

	// Validación de la estructura del DTO obtenido.
	switch rawDto := rawDto.(type) {
	case map[string]interface{}:
		dtoMap = rawDto

		// Se valida la entrada JSON contra las definiciones esperadas de los campos.
		utils.ValidateJSONInput(&entityBranchRequest.Branch, dtoMap, salvia_daos.EntityBranchJSONName, salvia_daos.EntityBranchFieldDefinitions, checkFields, common_config.Locale, common_config.DateTime.DATE_TIME_FORMAT, common_config.DateTime.DATE_FORMAT, common_config.DateTime.TIME_FORMAT, collectedErrors, false)

		entityBranchRequest.Branch.EntityBranchEntity = salvia_daos.EntityDTO{EntityICode: entityBranchRequest.Branch.EntityBranchEntityICode}
		/*
			var container interface{}
			var found bool = false
			container, found = dtoMap[salvia_daos.EntityBranchJSONName]
			if found {
				switch container := container.(type) {
				case map[string]interface{}:
					var c map[string]interface{} = container
					var entityTag string = utils.GetTag(&entityBranch, "EntityBranchEntityICode", "json")
					if value, found := c[entityTag]; found {
						switch icode := value.(type) {
						case string:
							entityBranch.EntityBranchEntity = salvia_daos.EntityDTO{EntityICode: icode}
						}
					}
				}

			}*/
		//--------------------------------------------------------------------
	default:
		// Si la estructura no es la esperada, se registra un error global.
		utils.SetError(collectedErrors, salvia_daos.EntityBranchJSONName, "default", common_config.Enums.GLOBAL_ERROR, "", security_config.Locale)
	}

	// Si se han detectado errores durante la validación, se retorna un error 400 con los detalles en formato JSON.
	if len(collectedErrors) > 0 {
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// Se verifica la entidad.
	var by common_controllers.By = common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"EntityBranchICode"},
		AttrsValue: []interface{}{entityBranchRequest.Branch.EntityBranchICode},
	}

	if err = salvia_daos.GetEntityBranch(by, &entityBranchToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	entityBranchToUpdate.EntityBranchAddress = entityBranchRequest.Branch.EntityBranchAddress
	entityBranchToUpdate.EntityBranchDescription = entityBranchRequest.Branch.EntityBranchDescription
	entityBranchToUpdate.EntityBranchName = entityBranchRequest.Branch.EntityBranchName
	entityBranchToUpdate.EntityBranchLatitude = entityBranchRequest.Branch.EntityBranchLatitude
	entityBranchToUpdate.EntityBranchLongitude = entityBranchRequest.Branch.EntityBranchLongitude

	//Si es creado por el sistema, se duplica
	if entityBranchToUpdate.EntityBranchSource == "s" {
		salvia_daos.SetEntityBranchDefaults(&entityBranchToUpdate, common_dao.SQL_INSERT)
		entityBranchToUpdate.EntityBranchSource = "u"
		if err = salvia_daos.SetEntityBranch(&entityBranchToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
			// En caso de error durante la creación, se retorna un error 500 con el mensaje correspondiente.
			return http.StatusInternalServerError, err.Error()
		}
	} else {
		// Se procede a crear el registro del log en la base de datos.
		if err = salvia_daos.UpdateEntityBranch(&entityBranchToUpdate, connData, &dbClientConfig, &dbServerConfig); err != nil {
			// En caso de error durante la creación, se retorna un error 500 con el mensaje correspondiente.
			return http.StatusInternalServerError, err.Error()
		}
	}

	// Se retorna el estado 200 OK con el DTO del log de caso en formato JSON.
	return http.StatusOK, utils.CommMsgGetJSONSuccess(entityBranchToUpdate)
}
