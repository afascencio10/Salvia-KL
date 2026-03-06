// Package security_daos contiene las funciones de acceso a datos (DAO) para la entidad PhoneNumber,
// permitiendo operaciones de inserción, actualización, eliminación y consulta en la base de datos.
package security_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Variables globales para la entidad PhoneNumber.
// Estas variables definen nombres de la entidad en distintos contextos (JSON, base de datos y esquema)
// y las definiciones de los campos que se utilizarán en las validaciones y generación de queries.
var (
	// Nombre de la entidad para referencia interna.
	PhoneNumberntityName string = "PhoneNumber"
	// Nombre utilizado en el JSON.
	PhoneNumberJSONName string = "phoneNumber"
	// Nombre de la tabla en la base de datos.
	PhoneNumberDBName string = "phone_number"
	// Esquema de la base de datos en el que se encuentra la tabla.
	PhoneNumberDBScheme string = "security"

	// Atributos relacionados con las validaciones.
	// Map que relaciona el nombre del campo en el modelo con su definición (nombre en DB, alias, tipo, tamaños, requerido).
	PhoneNumberFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"PhoneNumberId":                 {Name: "PhoneNumberId", DBName: "phone_number_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"PhoneNumberICode":              {Name: "PhoneNumberICode", DBName: "phone_number_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"PhoneNumberCreationDate":       {Name: "PhoneNumberCreationDate", DBName: "phone_number_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"PhoneNumberUpdateDate":         {Name: "PhoneNumberUpdateDate", DBName: "phone_number_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"PhoneNumberData":               {Name: "PhoneNumberData", DBName: "phone_number_data", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 128, Required: true},
		"PhoneNumberGeneralUserProfile": {Name: "PhoneNumberGeneralUserProfile", DBName: "phone_number_general_user_profile", Alias: "", ModelType: "uint", Required: false},
	}
)

// PhoneNumberDTO representa el Data Transfer Object para la entidad PhoneNumber.
// Incluye la información del teléfono y el perfil de usuario asociado.
type PhoneNumberDTO struct {
	PhoneNumberId                 uint64                `json:"-"`
	PhoneNumberICode              string                `json:"icode"`
	PhoneNumberCreationDate       time.Time             `json:"creation_date"`
	PhoneNumberUpdateDate         time.Time             `json:"update_date"`
	PhoneNumberData               string                `json:"data"`
	PhoneNumberGeneralUserProfile GeneralUserProfileDTO `json:"profile"`
}

// PhoneNumberPgDB representa la estructura de la entidad PhoneNumber tal y como se almacena en la base de datos PostgreSQL.
// Se utiliza para el mapeo de datos desde y hacia la base de datos, utilizando tipos sql.Null*.
type PhoneNumberPgDB struct {
	PhoneNumberId                 sql.NullInt64
	PhoneNumberICode              sql.NullString
	PhoneNumberCreationDate       sql.NullTime
	PhoneNumberUpdateDate         sql.NullTime
	PhoneNumberData               sql.NullString
	PhoneNumberGeneralUserProfile sql.NullInt64
}

// SetPhoneNumber inserta un nuevo registro de PhoneNumber en la base de datos.
//
// Parámetros:
//   - phoneNumber: puntero al DTO con la información del teléfono a insertar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un error en caso de fallo durante la inserción.
func SetPhoneNumber(phoneNumber *PhoneNumberDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se prepara el valor del campo del perfil asociado.
	var profileId interface{} = phoneNumber.PhoneNumberGeneralUserProfile.GeneralUserProfileId
	if phoneNumber.PhoneNumberGeneralUserProfile.GeneralUserProfileId == 0 {
		profileId = nil
	}

	// Se obtiene o configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Preparación de la lista de campos que se insertarán en la base de datos.
	var phoneNumberFieldsSlice []string = []string{"PhoneNumberICode", "PhoneNumberCreationDate", "PhoneNumberUpdateDate", "PhoneNumberData", "PhoneNumberGeneralUserProfile"}
	var phoneNumberFieldsAliasSlice []string = []string{}

	// Generación dinámica de la query SQL de inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{}, []string{"PhoneNumberId"}, common_dao.SQL_AND, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, false)

	// Ejecución de la query utilizando el controlador de persistencia.
	persistenceCtrl.QueryRow(context.Background(), query,
		phoneNumber.PhoneNumberICode, phoneNumber.PhoneNumberCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), phoneNumber.PhoneNumberUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), phoneNumber.PhoneNumberData, profileId)

	// Escaneo del resultado para obtener el ID asignado.
	persistenceCtrl.Scan(&phoneNumber.PhoneNumberId)

	// Manejo de errores en caso de fallo en la query.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdatePhoneNumberByICode actualiza un registro de PhoneNumber identificado por su iCode.
//
// Parámetros:
//   - phoneNumber: puntero al DTO con la información del teléfono a actualizar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un error en caso de fallo durante la actualización o si no se afecta ninguna fila.
func UpdatePhoneNumberByICode(phoneNumber *PhoneNumberDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Si ocurre un error al obtener la conexión, se retorna el error.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se define la lista de campos que se actualizarán.
	var phoneNumberFieldsSlice []string = []string{"PhoneNumberUpdateDate", "PhoneNumberData"}
	var phoneNumberFieldsAliasSlice []string = []string{}

	// Generación dinámica de la query SQL de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{"PhoneNumberICode"}, []string{}, common_dao.SQL_AND, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, false)

	// Ejecución de la query SQL con los nuevos valores.
	persistenceCtrl.Exec(context.Background(), query,
		phoneNumber.PhoneNumberICode, phoneNumber.PhoneNumberUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), phoneNumber.PhoneNumberData)

	// Manejo de errores en la ejecución de la query.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica si la actualización afectó al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemovePhoneNumberByICode elimina un registro de PhoneNumber identificado por su iCode.
//
// Parámetros:
//   - phoneNumber: puntero al DTO que contiene el iCode del teléfono a eliminar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un error en caso de fallo durante la eliminación o si no se afecta ninguna fila.
func RemovePhoneNumberByICode(phoneNumber *PhoneNumberDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Manejo de error en la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// No se requiere lista de campos para eliminación.
	var phoneNumberFieldsSlice []string = []string{}
	var phoneNumberFieldsAliasSlice []string = []string{}

	// Generación dinámica de la query SQL de eliminación.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{"PhoneNumberICode"}, []string{}, common_dao.SQL_AND, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, false)

	// Ejecución de la query SQL para eliminar el registro.
	persistenceCtrl.Exec(context.Background(), query, phoneNumber.PhoneNumberICode)

	// Manejo de errores en la ejecución.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica si la eliminación afectó al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemovePhoneNumbers elimina registros de PhoneNumber basándose en criterios dinámicos.
//
// Parámetros:
//   - by: estructura que contiene los atributos, alias y valores para filtrar la eliminación.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un error en caso de fallo durante la eliminación.
func RemovePhoneNumbers(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Manejo de error en la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Generación dinámica de la query SQL de eliminación basada en los criterios definidos en "by".
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, PhoneNumberDBName, by.AttrsName, []string{}, []string{}, by.Operator, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, false)

	// Ejecución de la query SQL con los valores de los atributos para filtrar.
	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)

	// Manejo de errores en la ejecución.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetPhoneNumber obtiene un registro de PhoneNumber y su perfil asociado a partir de criterios dinámicos.
//
// Parámetros:
//   - by: estructura que contiene los atributos, alias y valores para filtrar la consulta.
//   - phoneNumber: puntero al DTO donde se almacenarán los datos recuperados.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un error en caso de fallo durante la consulta.
func GetPhoneNumber(by common_controllers.By, phoneNumber *PhoneNumberDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Definición de las rutas completas de la tabla PhoneNumber y el perfil asociado.
	var phoneNumberPath string = PhoneNumberDBScheme + "." + PhoneNumberDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName

	// Inicializa el DTO del perfil asociado.
	phoneNumber.PhoneNumberGeneralUserProfile = GeneralUserProfileDTO{}

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Manejo de error en la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar para PhoneNumber y su perfil asociado.
	var phoneNumberFieldsSlice []string = []string{"PhoneNumberId", "PhoneNumberICode", "PhoneNumberCreationDate", "PhoneNumberUpdateDate", "PhoneNumberData", "PhoneNumberGeneralUserProfile"}
	var phoneNumberFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Generación de la parte de la query que selecciona los campos de PhoneNumber.
	var phoneNumberFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, true)
	// Generación de la parte de la query que selecciona los campos del perfil.
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Construcción completa de la query SQL con LEFT JOIN para obtener la información del PhoneNumber y su perfil.
	var query string = `SELECT ` + phoneNumberFieldsStr + `, ` + profileFieldsStr +
		` FROM ` + phoneNumberPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + phoneNumberPath + `.` + PhoneNumberFieldDefinitions["PhoneNumberGeneralUserProfile"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, PhoneNumberDBName, by.AttrsName, []string{}, []string{}, by.Operator, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, true)

	// Ejecución de la query.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se crea una variable auxiliar para almacenar los datos del perfil obtenido de la base de datos.
	var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
	// Escaneo de los resultados de la query en los campos del DTO y la estructura auxiliar del perfil.
	persistenceCtrl.Scan(&phoneNumber.PhoneNumberId, &phoneNumber.PhoneNumberICode, &phoneNumber.PhoneNumberCreationDate, &phoneNumber.PhoneNumberUpdateDate, &phoneNumber.PhoneNumberData,
		&profile.GeneralUserProfileId,
		&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

	// Conversión de la estructura auxiliar a DTO.
	phoneNumber.PhoneNumberGeneralUserProfile = profile.ToDTO()

	// Manejo de errores en la ejecución de la query.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetPhoneNumbers obtiene múltiples registros de PhoneNumber basándose en criterios dinámicos.
//
// Parámetros:
//   - by: estructura que contiene los atributos, alias y valores para filtrar la consulta.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un slice de DTO con los registros de PhoneNumber obtenidos.
//   - Un error en caso de fallo durante la consulta.
func GetPhoneNumbers(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]PhoneNumberDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Definición de la ruta completa de la tabla PhoneNumber.
	var phoneNumberPath string = PhoneNumberDBScheme + "." + PhoneNumberDBName

	var phoneNumber = PhoneNumberDTO{}
	var phoneNumbers []PhoneNumberDTO

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Manejo de error en la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar para PhoneNumber.
	var phoneNumberFieldsSlice []string = []string{"PhoneNumberId", "PhoneNumberICode", "PhoneNumberCreationDate", "PhoneNumberUpdateDate", "PhoneNumberData", "PhoneNumberGeneralUserProfile"}
	var phoneNumberFieldsAliasSlice []string = []string{}

	// Generación de la parte de la query que selecciona los campos.
	var phoneNumberFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, true)

	// Construcción de la query completa con cláusula WHERE basada en los criterios dinámicos.
	var query string = `SELECT ` + phoneNumberFieldsStr +
		` FROM ` + phoneNumberPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, PhoneNumberDBName, by.AttrsName, []string{}, []string{}, by.Operator, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, true)

	// Ejecución de la query.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Iteración sobre los resultados obtenidos.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&phoneNumber.PhoneNumberId, &phoneNumber.PhoneNumberICode, &phoneNumber.PhoneNumberCreationDate, &phoneNumber.PhoneNumberUpdateDate, &phoneNumber.PhoneNumberData, &phoneNumber.PhoneNumberGeneralUserProfile.GeneralUserProfileId)
		phoneNumbers = append(phoneNumbers, phoneNumber)
	}
	// Manejo de errores durante la iteración.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return phoneNumbers, nil
}

// GetAllPhoneNumber obtiene todos los registros de PhoneNumber de la base de datos.
//
// Parámetros:
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un slice de DTO con todos los registros de PhoneNumber.
//   - Un error en caso de fallo durante la consulta.
func GetAllPhoneNumber(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]PhoneNumberDTO, error) {
	// Inicialización del controlador de persistencia y variable temporal para PhoneNumber.
	var phoneNumber PhoneNumberDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Manejo de error en la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar para PhoneNumber.
	var phoneNumberFieldsSlice []string = []string{"PhoneNumberId", "PhoneNumberICode", "PhoneNumberCreationDate", "PhoneNumberUpdateDate", "PhoneNumberData", "PhoneNumberGeneralUserProfile"}
	var phoneNumberFieldsAliasSlice []string = []string{}

	// Generación dinámica de la query SQL de selección.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{}, []string{}, "", PhoneNumberDBScheme, PhoneNumberFieldDefinitions, true)

	// Ejecución de la query.
	persistenceCtrl.Query(context.Background(), query)
	var phoneNumbers []PhoneNumberDTO
	// Iteración sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		phoneNumber = PhoneNumberDTO{}
		persistenceCtrl.ScanRow(&phoneNumber.PhoneNumberId, &phoneNumber.PhoneNumberICode, &phoneNumber.PhoneNumberCreationDate, &phoneNumber.PhoneNumberUpdateDate, &phoneNumber.PhoneNumberData,
			&profile.GeneralUserProfileId)

		// Conversión del registro del perfil a DTO.
		phoneNumber.PhoneNumberGeneralUserProfile = profile.ToDTO()
		phoneNumbers = append(phoneNumbers, phoneNumber)
	}

	// Manejo de errores en la iteración.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return phoneNumbers, nil
}

// GetAllPhoneNumberWithProfile obtiene todos los registros de PhoneNumber junto con la información completa del perfil asociado.
//
// Parámetros:
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: nombre del módulo o contexto de la operación.
//   - connData: conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - La conexión de base de datos actualizada.
//   - Un slice de DTO con todos los registros de PhoneNumber y sus perfiles asociados.
//   - Un error en caso de fallo durante la consulta.
func GetAllPhoneNumberWithProfile(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]PhoneNumberDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Definición de las rutas completas de la tabla PhoneNumber y la tabla del perfil.
	var phoneNumberPath string = PhoneNumberDBScheme + "." + PhoneNumberDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName

	// Se configura o obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Manejo de error en la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar para PhoneNumber y el perfil asociado.
	var phoneNumberFieldsSlice []string = []string{"PhoneNumberId", "PhoneNumberICode", "PhoneNumberCreationDate", "PhoneNumberUpdateDate", "PhoneNumberData", "PhoneNumberGeneralUserProfile"}
	var phoneNumberFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Generación de la parte de la query que selecciona los campos de PhoneNumber.
	var phoneNumberFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, phoneNumberFieldsSlice, phoneNumberFieldsAliasSlice, PhoneNumberDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, PhoneNumberDBScheme, PhoneNumberFieldDefinitions, true)
	// Generación de la parte de la query que selecciona los campos del perfil.
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Construcción completa de la query SQL con LEFT JOIN para obtener la información completa.
	var query string = `SELECT ` + phoneNumberFieldsStr + ", " + profileFieldsStr +
		` FROM ` + phoneNumberPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + phoneNumberPath + `.` + PhoneNumberFieldDefinitions["PhoneNumberGeneralUserProfile"].DBName + `)`

	// Ejecución de la query.
	persistenceCtrl.Query(context.Background(), query)
	var phoneNumbers []PhoneNumberDTO
	// Iteración sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var phoneNumber = PhoneNumberDTO{}

		persistenceCtrl.ScanRow(&phoneNumber.PhoneNumberId, &phoneNumber.PhoneNumberICode, &phoneNumber.PhoneNumberCreationDate, &phoneNumber.PhoneNumberUpdateDate, &phoneNumber.PhoneNumberData,
			&profile.GeneralUserProfileId,
			&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

		// Conversión del registro del perfil a DTO.
		phoneNumber.PhoneNumberGeneralUserProfile = profile.ToDTO()

		phoneNumbers = append(phoneNumbers, phoneNumber)
	}

	// Manejo de errores durante la iteración.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return phoneNumbers, nil
}

// SetPhoneNumberDefaults asigna valores predeterminados a los campos de PhoneNumber en función de la acción a realizar.
//
// Parámetros:
//   - phoneNumber: puntero al DTO de PhoneNumber al que se asignarán los valores predeterminados.
//   - action: tipo de acción (por ejemplo, SQL_INSERT o SQL_UPDATE) que determina qué campos se deben inicializar.
func SetPhoneNumberDefaults(phoneNumber *PhoneNumberDTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción, se asignan las fechas actuales y se genera un UUID para el iCode.
		phoneNumber.PhoneNumberCreationDate = time.Now()
		phoneNumber.PhoneNumberUpdateDate = time.Now()
		phoneNumber.PhoneNumberICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		// Para actualización, se actualiza la fecha de modificación.
		phoneNumber.PhoneNumberUpdateDate = time.Now()
	}

}

// PgDBToDTO convierte la estructura PhoneNumberPgDB (formato de base de datos) a PhoneNumberDTO (formato de transferencia de datos).
//
// Retorna:
//   - Un objeto PhoneNumberDTO con los datos mapeados desde PhoneNumberPgDB.
func (obj *PhoneNumberPgDB) PgDBToDTO() PhoneNumberDTO {
	var dto PhoneNumberDTO

	if obj.PhoneNumberId.Valid {
		dto.PhoneNumberId = uint64(obj.PhoneNumberId.Int64)
	}

	if obj.PhoneNumberICode.Valid {
		dto.PhoneNumberICode = obj.PhoneNumberICode.String
	}

	if obj.PhoneNumberCreationDate.Valid {
		dto.PhoneNumberCreationDate = obj.PhoneNumberCreationDate.Time
	}

	if obj.PhoneNumberUpdateDate.Valid {
		dto.PhoneNumberUpdateDate = obj.PhoneNumberUpdateDate.Time
	}

	if obj.PhoneNumberData.Valid {
		dto.PhoneNumberData = obj.PhoneNumberData.String
	}

	if obj.PhoneNumberGeneralUserProfile.Valid {
		dto.PhoneNumberGeneralUserProfile = GeneralUserProfileDTO{GeneralUserProfileId: uint64(obj.PhoneNumberGeneralUserProfile.Int64)}
	}

	return dto
}
