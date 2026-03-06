// Package security_daos proporciona las funciones de acceso a datos para gestionar la relación entre roles y usuarios generales
// dentro del módulo de seguridad. Incluye operaciones para insertar, eliminar y obtener relaciones mediante consultas SQL,
// utilizando una capa de persistencia que abstrae la conexión y la ejecución de consultas en la base de datos.
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

var (
	// RelRoleRelRoleGeneralUserName representa el nombre de la relación de roles y usuarios generales en el modelo.
	RelRoleRelRoleGeneralUserName string = "RelRoleGeneralUser"
	// RelRoleGeneralUserJSONName es el nombre que se utiliza en el JSON para representar la relación.
	RelRoleGeneralUserJSONName string = "rel"
	// RelRoleGeneralUserDBName es el nombre de la tabla en la base de datos que almacena la relación.
	RelRoleGeneralUserDBName string = "rel_role_general_user"
	// RelRoleGeneralUserDBScheme indica el esquema de la base de datos en el que se encuentra la tabla.
	RelRoleGeneralUserDBScheme string = "security"

	// RelRoleGeneralUserFieldDefinitions contiene las definiciones de los campos asociados a la relación,
	// especificando las validaciones, nombres en el modelo y en la base de datos, tipos y requisitos.
	RelRoleGeneralUserFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelRoleGeneralUserId":           {Name: "RelRoleGeneralUserId", DBName: "rel_role_general_user_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelRoleGeneralUserCreationDate": {Name: "RelRoleGeneralUserCreationDate", DBName: "rel_role_general_user_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"RelRoleGeneralUserRole":         {Name: "RelRoleGeneralUserRole", DBName: "role_id", Alias: "", ModelType: "uint", Required: false},
		"RelRoleGeneralUserGeneralUser":  {Name: "RelRoleGeneralUserGeneralUser", DBName: "general_user_id", Alias: "", ModelType: "uint", Required: false},
	}
)

// RelRoleGeneralUserDTO es el objeto de transferencia de datos que representa la relación entre un rol y un usuario general.
// Incluye la fecha de creación, y los detalles del rol y del usuario asociados.
type RelRoleGeneralUserDTO struct {
	RelRoleGeneralUserId           uint64         `json:"-"`
	RelRoleGeneralUserCreationDate time.Time      `json:"creation_date"`
	RelRoleGeneralUserRole         RoleDTO        `json:"role"`
	RelRoleGeneralUserGeneralUser  GeneralUserDTO `json:"user"`
}

// RelRoleGeneralUserPgDB representa el modelo de datos en PostgreSQL para la relación entre rol y usuario general,
// utilizando tipos sql.Null* para gestionar valores nulos.
type RelRoleGeneralUserPgDB struct {
	RelRoleGeneralUserId           sql.NullInt64
	RelRoleGeneralUserCreationDate sql.NullTime
	RelRoleGeneralUserRole         sql.NullInt64
	RelRoleGeneralUserGeneralUser  sql.NullInt64
}

// SetRelRoleGeneralUser inserta una nueva relación entre un rol y un usuario general en la base de datos.
// Parámetros:
//   - rel: puntero al DTO que contiene los datos de la relación a insertar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación (útil para logs).
//   - connData: datos de la conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada y un error en caso de producirse algún fallo.
func SetRelRoleGeneralUser(rel *RelRoleGeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se van a insertar.
	var relFieldsSlice []string = []string{"RelRoleGeneralUserCreationDate", "RelRoleGeneralUserRole", "RelRoleGeneralUserGeneralUser"}
	var relFieldsAliasSlice []string = []string{}

	// Se genera la consulta SQL para la inserción utilizando las definiciones de campo.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, relFieldsSlice, relFieldsAliasSlice, RelRoleGeneralUserDBName, []string{}, []string{}, []string{"RelRoleGeneralUserId"}, common_dao.SQL_AND, RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, false)

	// Se ejecuta la consulta pasando los valores correspondientes, formateando la fecha de creación.
	persistenceCtrl.QueryRow(context.Background(), query, rel.RelRoleGeneralUserCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), rel.RelRoleGeneralUserRole.RoleId, rel.RelRoleGeneralUserGeneralUser.GeneralUserId)

	// Se escanea el resultado para obtener el ID generado de la nueva relación.
	persistenceCtrl.Scan(&rel.RelRoleGeneralUserId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// RemoveRelRoleGeneralUser elimina una relación específica entre un rol y un usuario general en la base de datos.
// Parámetros:
//   - rel: puntero al DTO que contiene el ID de la relación a eliminar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que invoca la operación.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada y un error en caso de fallo.
func RemoveRelRoleGeneralUser(rel *RelRoleGeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// No se definen campos adicionales para la eliminación.
	var relFieldsSlice []string = []string{}
	var relFieldsAliasSlice []string = []string{}

	// Se genera la consulta SQL para eliminar la relación basándose en el ID.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, relFieldsSlice, relFieldsAliasSlice, RelRoleGeneralUserDBName, []string{}, []string{"RelRoleGeneralUserId"}, []string{}, common_dao.SQL_AND, RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, false)

	// Se ejecuta la consulta pasando el ID de la relación a eliminar.
	persistenceCtrl.Exec(context.Background(), query, rel.RelRoleGeneralUserId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Si no se afectaron filas, se retorna un error indicando que no se encontró la relación.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveRelRoleGeneralUsers elimina una o más relaciones entre roles y usuarios generales basándose en criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (nombres de atributos, alias, operadores y valores).
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que invoca la operación.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada y un error en caso de fallo.
func RemoveRelRoleGeneralUsers(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se genera la consulta SQL para eliminar relaciones basándose en los atributos indicados.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, RelRoleGeneralUserDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, false)

	// Se ejecuta la consulta con los valores de los atributos.
	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelRoleGeneralUser obtiene una relación específica entre un rol y un usuario general según los criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (atributos, alias, operador y valores).
//   - rel: puntero al DTO donde se almacenarán los datos obtenidos.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que invoca la operación.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada y un error en caso de fallo.
func GetRelRoleGeneralUser(by common_controllers.By, rel *RelRoleGeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia y definición de rutas de tabla.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var relPath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de campos a seleccionar para la relación, el rol y el usuario.
	var relFieldsSlice []string = []string{"RelRoleGeneralUserId", "RelRoleGeneralUserCreationDate", "RelRoleGeneralUserRole", "RelRoleGeneralUserGeneralUser"}
	var relFieldsAliasSlice []string = []string{}

	var roleFieldsSlice []string = []string{"RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	var userFieldsSlice []string = []string{"GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserStatus", "GeneralUserLanguage", "GeneralUserGeneralUserProfile"}
	var userFieldsAliasSlice []string = []string{}

	// Generación de cadenas SQL para la selección de campos.
	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelRoleGeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, true)
	var roleFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, true)
	var userFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, userFieldsSlice, userFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)

	// Construcción de la consulta SQL con LEFT JOIN para incluir los datos del rol y del usuario.
	var query string = `SELECT ` + relFieldsStr + `, ` + roleFieldsStr + `, ` + userFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + rolePath + ` ON (` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + ` = ` + relPath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + ` = ` + relPath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelRoleGeneralUserDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, true)

	// Ejecución de la consulta con los valores de búsqueda especificados.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Variables temporales para almacenar los datos del rol y del usuario.
	var role RolePgDB = RolePgDB{}
	var user GeneralUserPgDB = GeneralUserPgDB{}

	// Escaneo de los resultados en el DTO de la relación y en los modelos de rol y usuario.
	persistenceCtrl.Scan(&rel.RelRoleGeneralUserId, &rel.RelRoleGeneralUserCreationDate,
		&role.RoleId, &user.GeneralUserId,
		&role.RoleICode, &role.RoleCreationDate, &role.RoleUpdateDate, &role.RoleCode, &role.RoleName, &role.RoleDescription,
		&user.GeneralUserICode, &user.GeneralUserCreationDate, &user.GeneralUserUpdateDate, &user.GeneralUserLogin, &user.GeneralUserStatus, &user.GeneralUserLanguage, &user.GeneralUserGeneralUserProfile)

	// Conversión de los datos de rol y usuario a sus respectivos DTO.
	rel.RelRoleGeneralUserRole = role.ToDTO()
	rel.RelRoleGeneralUserGeneralUser = user.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelRoleGeneralUsers obtiene una lista de relaciones entre roles y usuarios generales basándose en criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que invoca la operación.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada, una lista de DTOs con los datos de las relaciones y un error en caso de fallo.
func GetRelRoleGeneralUsers(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelRoleGeneralUserDTO, error) {
	// Inicialización del controlador de persistencia y definición de la ruta de la tabla de rol.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var rolePath string = RoleDBScheme + "." + RoleDBName

	var role = RelRoleGeneralUserDTO{}
	var roles []RelRoleGeneralUserDTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos que se van a seleccionar.
	var roleFieldsSlice []string = []string{"RelRoleGeneralUserCreationDate", "RelRoleGeneralUserGeneralUser", "RelRoleGeneralUserRole"}
	var roleFieldsAliasSlice []string = []string{}

	// Generación de la cadena SQL que define los campos a seleccionar.
	var roleFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, true)

	// Construcción de la consulta SQL para obtener las relaciones.
	var query string = `SELECT ` + roleFieldsStr +
		` FROM ` + rolePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RoleDBName, by.AttrsName, []string{}, []string{}, by.Operator, RoleDBScheme, RoleFieldDefinitions, true)

	// Impresión del query (posiblemente para depuración) y ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Iteración sobre los resultados obtenidos y almacenamiento en el slice de roles.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&role.RelRoleGeneralUserId, &role.RelRoleGeneralUserCreationDate, &role.RelRoleGeneralUserGeneralUser.GeneralUserId, &role.RelRoleGeneralUserRole.RoleId)
		roles = append(roles, role)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un 'by' para todos los DAOS
	return roles, nil
}

// GetAllRelRoleGeneralUser obtiene todas las relaciones entre roles y usuarios generales sin aplicar filtros.
// Parámetros:
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que invoca la operación.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada, una lista de DTOs con todas las relaciones y un error en caso de fallo.
func GetAllRelRoleGeneralUser(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelRoleGeneralUserDTO, error) {
	// Declaración de variables locales.
	var rel RelRoleGeneralUserDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar para la relación.
	var relFieldsSlice []string = []string{"RelRoleGeneralUserId", "RelRoleGeneralUserCreationDate", "RelRoleGeneralUserRole", "RelRoleGeneralUserGeneralUser"}
	var relFieldsAliasSlice []string = []string{}

	// Generación de la consulta SQL para seleccionar todas las relaciones.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, relFieldsSlice, relFieldsAliasSlice, RelRoleGeneralUserDBName, []string{}, []string{}, []string{}, "", RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, true)

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var rels []RelRoleGeneralUserDTO
	for persistenceCtrl.Next() {
		var role RolePgDB = RolePgDB{}
		var user GeneralUserPgDB = GeneralUserPgDB{}
		rel = RelRoleGeneralUserDTO{}
		persistenceCtrl.ScanRow(&rel.RelRoleGeneralUserId, &rel.RelRoleGeneralUserCreationDate,
			&role.RoleId, &user.GeneralUserId)

		// Conversión de los datos de rol y usuario a sus respectivos DTO.
		rel.RelRoleGeneralUserRole = role.ToDTO()
		rel.RelRoleGeneralUserGeneralUser = user.ToDTO()
		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// GetAllRelRoleGeneralUserWithRoleAndGeneralUser obtiene todas las relaciones entre roles y usuarios generales
// junto con los detalles asociados del rol y del usuario.
// Parámetros:
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que invoca la operación.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna la conexión actualizada, una lista de DTOs con las relaciones y sus datos asociados, y un error en caso de fallo.
func GetAllRelRoleGeneralUserWithRoleAndGeneralUser(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelRoleGeneralUserDTO, error) {
	// Inicialización del controlador de persistencia y definición de rutas para las tablas.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var relPath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar para la relación, el rol y el usuario.
	var relFieldsSlice []string = []string{"RelRoleGeneralUserId", "RelRoleGeneralUserCreationDate", "RelRoleGeneralUserRole", "RelRoleGeneralUserGeneralUser"}
	var relFieldsAliasSlice []string = []string{}

	var roleFieldsSlice []string = []string{"RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	var userFieldsSlice []string = []string{"GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserStatus", "GeneralUserLanguage", "GeneralUserGeneralUserProfile"}
	var userFieldsAliasSlice []string = []string{}

	// Generación de las cadenas SQL que representan los campos a seleccionar.
	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelRoleGeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelRoleGeneralUserDBScheme, RelRoleGeneralUserFieldDefinitions, true)
	var roleFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, true)
	var userFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, userFieldsSlice, userFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)

	// Construcción de la consulta SQL con LEFT JOIN para incluir detalles del rol y del usuario.
	var query string = `SELECT ` + relFieldsStr + ", " + roleFieldsStr + ", " + userFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + rolePath + ` ON (` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + ` = ` + relPath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + ` = ` + relPath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + `)`

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var rels []RelRoleGeneralUserDTO
	for persistenceCtrl.Next() {
		var role RolePgDB = RolePgDB{}
		var user GeneralUserPgDB = GeneralUserPgDB{}
		var rel = RelRoleGeneralUserDTO{}

		persistenceCtrl.ScanRow(&rel.RelRoleGeneralUserId, &rel.RelRoleGeneralUserCreationDate,
			&role.RoleId, &user.GeneralUserId,
			&role.RoleICode, &role.RoleCreationDate, &role.RoleUpdateDate, &role.RoleCode, &role.RoleName, &role.RoleDescription,
			&user.GeneralUserICode, &user.GeneralUserCreationDate, &user.GeneralUserUpdateDate, &user.GeneralUserLogin, &user.GeneralUserStatus, &user.GeneralUserLanguage, &user.GeneralUserGeneralUserProfile)

		// Conversión de los modelos de rol y usuario a sus respectivos DTO.
		rel.RelRoleGeneralUserRole = role.ToDTO()
		rel.RelRoleGeneralUserGeneralUser = user.ToDTO()

		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// SetRelRoleGeneralUserDefaults asigna valores por defecto a ciertos campos de la relación en función de la acción que se realizará.
// Por ejemplo, en una inserción se establece la fecha de creación.
// Parámetros:
//   - rel: puntero al DTO de la relación a modificar.
//   - action: acción que se va a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetRelRoleGeneralUserDefaults(rel *RelRoleGeneralUserDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se asigna la fecha actual como fecha de creación.
		rel.RelRoleGeneralUserCreationDate = time.Now()

	case common_dao.SQL_UPDATE:
		// Para una actualización, se podrían asignar otros valores por defecto si fuera necesario.
	}
}

// PgDBToDTO convierte una instancia de RelRoleGeneralUserPgDB (modelo de datos de PostgreSQL)
// a su correspondiente objeto de transferencia de datos (DTO) RelRoleGeneralUserDTO.
// Retorna el DTO con los datos convertidos.
func (obj *RelRoleGeneralUserPgDB) PgDBToDTO() RelRoleGeneralUserDTO {
	var dto RelRoleGeneralUserDTO

	if obj.RelRoleGeneralUserId.Valid {
		dto.RelRoleGeneralUserId = uint64(obj.RelRoleGeneralUserId.Int64)
	}

	if obj.RelRoleGeneralUserCreationDate.Valid {
		dto.RelRoleGeneralUserCreationDate = obj.RelRoleGeneralUserCreationDate.Time
	}

	if obj.RelRoleGeneralUserRole.Valid {
		dto.RelRoleGeneralUserRole = RoleDTO{RoleId: uint64(obj.RelRoleGeneralUserRole.Int64)}
	}

	if obj.RelRoleGeneralUserGeneralUser.Valid {
		dto.RelRoleGeneralUserGeneralUser = GeneralUserDTO{GeneralUserId: uint64(obj.RelRoleGeneralUserGeneralUser.Int64)}
	}

	return dto
}
