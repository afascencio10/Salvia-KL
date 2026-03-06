// Package security_daos proporciona los Data Access Objects (DAOs) para la entidad Role,
// permitiendo realizar operaciones CRUD en la base de datos para el módulo de seguridad.
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
	// RoleRoleName es el nombre de la entidad Role en el contexto de la aplicación.
	RoleRoleName string = "Role"
	// RoleJSONName es el nombre utilizado para la serialización/deserialización JSON de Role.
	RoleJSONName string = "role"
	// RoleDBName es el nombre de la tabla en la base de datos correspondiente a Role.
	RoleDBName string = "role"
	// RoleDBScheme es el esquema de la base de datos donde se encuentra la tabla Role.
	RoleDBScheme string = "security"

	// RoleFieldDefinitions contiene la definición de campos y validaciones para la entidad Role.
	// Cada campo define su nombre, nombre en la base de datos, tipo de modelo, y restricciones de tamaño y obligatoriedad.
	RoleFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RoleId":           {Name: "RoleId", DBName: "role_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RoleICode":        {Name: "RoleICode", DBName: "role_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"RoleCreationDate": {Name: "RoleCreationDate", DBName: "role_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"RoleUpdateDate":   {Name: "RoleUpdateDate", DBName: "role_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"RoleCode":         {Name: "RoleCode", DBName: "role_code", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"RoleName":         {Name: "RoleName", DBName: "role_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 128, Required: true},
		"RoleDescription":  {Name: "RoleDescription", DBName: "role_description", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 256, Required: false},
	}
)

// RoleDTO representa el objeto de transferencia de datos (Data Transfer Object)
// para la entidad Role. Es utilizado para pasar la información del Role entre capas.
type RoleDTO struct {
	RoleId           uint64    `json:"-"`
	RoleICode        string    `json:"icode"`
	RoleCreationDate time.Time `json:"creation_date"`
	RoleUpdateDate   time.Time `json:"update_date"`
	RoleCode         string    `json:"code"`
	RoleName         string    `json:"name"`
	RoleDescription  string    `json:"description"`
}

// RolePgDB representa la estructura de la entidad Role tal como se almacena en la base de datos PostgreSQL.
// Utiliza tipos sql.Null para manejar valores nulos.
type RolePgDB struct {
	RoleId           sql.NullInt64
	RoleICode        sql.NullString
	RoleCreationDate sql.NullTime
	RoleUpdateDate   sql.NullTime
	RoleCode         sql.NullString
	RoleName         sql.NullString
	RoleDescription  sql.NullString
}

// SetRole inserta un nuevo registro de Role en la base de datos.
// Configura la conexión, construye la consulta SQL y realiza el INSERT.
//
// Parámetros:
//   - role: puntero a RoleDTO con la información del role a insertar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales, que pueden ser actualizados.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - error: en caso de producirse algún error durante la operación.
func SetRole(role *RoleDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia para ejecutar operaciones SQL.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura y obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a insertar.
	var roleFieldsSlice []string = []string{"RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de inserción utilizando utilidades comunes.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{"RoleId"}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, false)

	// Ejecución de la consulta con los parámetros correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query, role.RoleICode, role.RoleCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), role.RoleUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), role.RoleCode, role.RoleName, role.RoleDescription)

	// Se asigna el valor generado para RoleId al objeto role.
	persistenceCtrl.Scan(&role.RoleId)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateRoleByICode actualiza un registro de Role en la base de datos, identificándolo por su código único (ICode).
//
// Parámetros:
//   - role: puntero a RoleDTO que contiene los datos a actualizar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - error: en caso de producirse algún error durante la actualización o si no se afectó ninguna fila.
func UpdateRoleByICode(role *RoleDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia para ejecutar operaciones SQL.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura y obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a actualizar.
	var roleFieldsSlice []string = []string{"RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{"RoleICode"}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, false)

	// Ejecución de la consulta de actualización.
	persistenceCtrl.Exec(context.Background(), query,
		role.RoleICode, role.RoleUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), role.RoleCode, role.RoleName, role.RoleDescription)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verificación de que se haya afectado al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveRoleByICode elimina un registro de Role de la base de datos identificado por su código único (ICode).
//
// Parámetros:
//   - role: puntero a RoleDTO que contiene el código único del Role a eliminar.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - error: en caso de producirse algún error durante la eliminación o si no se afectó ninguna fila.
func RemoveRoleByICode(role *RoleDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia para ejecutar operaciones SQL.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura y obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// No se requieren campos adicionales para el DELETE.
	var roleFieldsSlice []string = []string{}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de eliminación.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{"RoleICode"}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, false)

	// Ejecución de la consulta de eliminación.
	persistenceCtrl.Exec(context.Background(), query, role.RoleICode)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verificación de que se haya afectado al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetRole recupera un registro de Role de la base de datos basándose en criterios de búsqueda proporcionados.
// Utiliza la estructura 'By' para determinar los atributos de búsqueda.
//
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (atributos, alias, valores, etc.).
//   - role: puntero a RoleDTO donde se almacenará la información recuperada.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - error: en caso de producirse algún error durante la consulta.
func GetRole(by common_controllers.By, role *RoleDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción del path completo de la tabla Role (esquema.nombre_tabla).
	var rolePath string = RoleDBScheme + "." + RoleDBName

	// Se configura y obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var roleFieldsSlice []string = []string{"RoleId", "RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la lista de campos a seleccionar.
	var roleFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, true)

	// Construcción de la consulta SQL de selección con cláusula WHERE basada en 'by'.
	var query string = `SELECT ` + roleFieldsStr +
		` FROM ` + rolePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RoleDBName, by.AttrsName, []string{}, []string{}, by.Operator, RoleDBScheme, RoleFieldDefinitions, true)

	// Ejecución de la consulta SQL.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Escaneo de los resultados en el objeto RoleDTO.
	persistenceCtrl.Scan(&role.RoleId, &role.RoleICode, &role.RoleCreationDate, &role.RoleUpdateDate,
		&role.RoleCode, &role.RoleName, &role.RoleDescription)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRolesByGeneralUser recupera una lista de Roles asociados a un usuario general.
// Se realiza una unión (JOIN) entre la tabla de Roles y la relación entre Roles y Usuarios Generales.
//
// Parámetros:
//   - user: puntero a GeneralUserDTO que identifica al usuario general.
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - []RoleDTO: slice con los roles asociados al usuario.
//   - error: en caso de producirse algún error durante la consulta.
func GetRolesByGeneralUser(user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RoleDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Construcción de los paths completos de las tablas involucradas.
	var relPath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName

	var role = RoleDTO{}
	var roles = []RoleDTO{}

	// Se configura y obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar para Role.
	var roleFieldsSlice []string = []string{"RoleId", "RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la lista de campos a seleccionar.
	var roleFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, true)

	// Construcción de la consulta SQL con RIGHT JOIN para obtener los roles asociados al usuario.
	var query string = `SELECT ` + roleFieldsStr +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relPath + ` ON (` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + ` = ` + relPath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + `) 
		WHERE ` + relPath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = $1`

	// Impresión de la consulta para depuración.
	fmt.Printf(query, user.GeneralUserId)
	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query, user.GeneralUserId)

	// Iteración sobre los resultados y conversión a objetos RoleDTO.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&role.RoleId, &role.RoleICode, &role.RoleCreationDate, &role.RoleUpdateDate, &role.RoleCode, &role.RoleName, &role.RoleDescription)
		roles = append(roles, role)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return roles, nil
}

// GetAllRoles recupera todos los registros de Role almacenados en la base de datos.
//
// Parámetros:
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - []RoleDTO: slice con todos los roles.
//   - error: en caso de producirse algún error durante la consulta.
func GetAllRoles(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RoleDTO, error) {
	var role RoleDTO
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura y obtiene la conexión.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var roleFieldsSlice []string = []string{"RoleId", "RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL para seleccionar todos los roles.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, "", RoleDBScheme, RoleFieldDefinitions, true)

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var roles []RoleDTO
	for persistenceCtrl.Next() {
		role = RoleDTO{}
		persistenceCtrl.ScanRow(&role.RoleId, &role.RoleICode, &role.RoleCreationDate, &role.RoleUpdateDate,
			&role.RoleCode, &role.RoleName, &role.RoleDescription)
		roles = append(roles, role)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return roles, nil
}

// GetRoles recupera una lista de registros de Role basándose en criterios de búsqueda específicos.
// Se utiliza la estructura 'By' para especificar los atributos y condiciones de búsqueda.
//
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (atributos, alias, operador, etc.).
//   - inTransaction: indica si la operación se realiza dentro de una transacción.
//   - module: identificador del módulo que realiza la operación.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna:
//   - connData: datos de conexión actualizados.
//   - []RoleDTO: slice con los roles que cumplen los criterios de búsqueda.
//   - error: en caso de producirse algún error durante la consulta.
func GetRoles(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RoleDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción del path completo de la tabla Role.
	var rolePath string = RoleDBScheme + "." + RoleDBName

	var role = RoleDTO{}
	var roles []RoleDTO

	// Se configura y obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var roleFieldsSlice []string = []string{"RoleId", "RoleICode", "RoleCreationDate", "RoleUpdateDate", "RoleCode", "RoleName", "RoleDescription"}
	var roleFieldsAliasSlice []string = []string{}

	// Construcción de la lista de campos a seleccionar.
	var roleFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, roleFieldsSlice, roleFieldsAliasSlice, RoleDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RoleDBScheme, RoleFieldDefinitions, true)

	// Construcción de la consulta SQL de selección con cláusula WHERE basada en 'by'.
	var query string = `SELECT ` + roleFieldsStr +
		` FROM ` + rolePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RoleDBName, by.AttrsName, by.EqualSign, []string{}, by.Operator, RoleDBScheme, RoleFieldDefinitions, true)

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Iteración sobre los resultados y conversión a objetos RoleDTO.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&role.RoleId, &role.RoleICode, &role.RoleCreationDate, &role.RoleUpdateDate,
			&role.RoleCode, &role.RoleName, &role.RoleDescription)
		roles = append(roles, role)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return roles, nil
}

// SetRoleDefaults establece valores por defecto para los campos de un RoleDTO, dependiendo de la acción que se
// va a realizar (inserción o actualización).
//
// Parámetros:
//   - role: puntero a RoleDTO sobre el que se establecerán los valores por defecto.
//   - action: cadena que indica la acción SQL (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetRoleDefaults(role *RoleDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se asignan la fecha actual y un UUID para el RoleICode.
		role.RoleCreationDate = time.Now()
		role.RoleUpdateDate = time.Now()
		role.RoleICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para una actualización, se actualiza únicamente la fecha de actualización.
		role.RoleUpdateDate = time.Now()
	}
}

// ToDTO convierte una estructura RolePgDB (representación de la base de datos) a RoleDTO (objeto de transferencia de datos).
//
// Retorna:
//   - RoleDTO: objeto con los valores convertidos y asignados desde RolePgDB.
func (obj *RolePgDB) ToDTO() RoleDTO {
	var dto RoleDTO

	// Conversión y asignación del RoleId si es válido.
	if obj.RoleId.Valid {
		dto.RoleId = uint64(obj.RoleId.Int64)
	}

	// Conversión y asignación del RoleICode si es válido.
	if obj.RoleICode.Valid {
		dto.RoleICode = obj.RoleICode.String
	}

	// Conversión y asignación de la fecha de creación si es válida.
	if obj.RoleCreationDate.Valid {
		dto.RoleCreationDate = obj.RoleCreationDate.Time
	}

	// Conversión y asignación de la fecha de actualización si es válida.
	if obj.RoleUpdateDate.Valid {
		dto.RoleUpdateDate = obj.RoleUpdateDate.Time
	}

	// Conversión y asignación del RoleCode si es válido.
	if obj.RoleCode.Valid {
		dto.RoleCode = obj.RoleCode.String
	}

	// Conversión y asignación del RoleName si es válido.
	if obj.RoleName.Valid {
		dto.RoleName = obj.RoleName.String
	}

	// Conversión y asignación del RoleDescription si es válido.
	if obj.RoleDescription.Valid {
		dto.RoleDescription = obj.RoleDescription.String
	}

	return dto
}
