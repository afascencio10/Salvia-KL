// Package salvia_daos contiene las funciones y estructuras necesarias para interactuar
// con la entidad CaseOwner en la base de datos, incluyendo operaciones de inserción,
// actualización, eliminación y consulta.
package salvia_daos

import (
	// Importación de paquetes internos y externos necesarios para la configuración,
	// control de persistencia, conexión a la base de datos y utilidades generales.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	// Constantes y variables que definen los nombres utilizados para la entidad CaseOwner
	CaseOwnerCaseOwnerName string = "CaseOwner"
	CaseOwnerJSONName      string = "caseOwner"
	CaseOwnerDBName        string = "case_owner"
	CaseOwnerDBScheme      string = "salvia"

	// Definición de los campos de la entidad CaseOwner junto con sus propiedades para validación.
	// Cada entrada en el mapa indica el nombre del campo, el nombre en la base de datos, tipo de dato,
	// tamaño mínimo y máximo, y si es requerido o no.
	CaseOwnerFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"CaseOwnerId":           {Name: "CaseOwnerId", DBName: "case_owner_id", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 0, Required: true},
		"CaseOwnerICode":        {Name: "CaseOwnerICode", DBName: "case_owner_i_code", Alias: "", ModelType: "string", MinSize: 32, MaxSize: 36, Required: true},
		"CaseOwnerCreationDate": {Name: "CaseOwnerCreationDate", DBName: "case_owner_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"CaseOwnerUpdateDate":   {Name: "CaseOwnerUpdateDate", DBName: "case_owner_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"CaseOwnerGeneralUser":  {Name: "CaseOwnerGeneralUser", DBName: "case_owner_general_user", Alias: "", ModelType: "string", MinSize: 32, MaxSize: 36, Required: true},
		"AttentionLine":         {Name: "AttentionLine", DBName: "attention_line_id", Alias: "", ModelType: "uint", Required: false},
		"EntityBranch":          {Name: "EntityBranch", DBName: "entity_branch_id", Alias: "", ModelType: "uint", Required: false},
		"CaseOwnerNumCases":     {Name: "CaseOwnerNumCases", DBName: "case_owner_num_cases", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// CaseOwnerDTO representa el Data Transfer Object de la entidad CaseOwner.
// Se utiliza para transferir los datos entre capas de la aplicación.
type CaseOwnerDTO struct {
	CaseOwnerId           uint64           `json:"-"`
	CaseOwnerICode        string           `json:"icode"`
	CaseOwnerCreationDate time.Time        `json:"creationDate"`
	CaseOwnerUpdateDate   time.Time        `json:"updateDate"`
	CaseOwnerGeneralUser  string           `json:"user"`
	CaseOwnerNumCases     int              `json:"-"`
	AttentionLine         AttentionLineDTO `json:"-"`
	EntityBranch          EntityBranchDTO  `json:"-"`
}

// CaseOwnerPgDB representa la estructura que se mapea directamente a los
// registros de la base de datos para la entidad CaseOwner.
// Utiliza tipos sql.Null* para manejar valores nulos en la BD.
type CaseOwnerPgDB struct {
	CaseOwnerId           sql.NullInt64
	CaseOwnerICode        sql.NullString
	CaseOwnerCreationDate sql.NullTime
	CaseOwnerUpdateDate   sql.NullTime
	CaseOwnerGeneralUser  sql.NullString
	VictimCase            sql.NullInt64
	AttentionLine         sql.NullInt64
	EntityBranch          sql.NullInt64
	CaseOwnerNumCases     sql.NullInt64
}

func (cw CaseOwnerDTO) MarshalJSON() ([]byte, error) {
	type Alias CaseOwnerDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		CaseOwnerCreationDate string `json:"creationDate"`
		CaseOwnerUpdateDate   string `json:"updateDate"`
	}{
		Alias:                 (*Alias)(&cw),
		CaseOwnerCreationDate: cw.CaseOwnerCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		CaseOwnerUpdateDate:   cw.CaseOwnerUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (cw *CaseOwnerDTO) UnmarshalJSON(data []byte) error {
	type Alias CaseOwnerDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		CaseOwnerCreationDate string `json:"creationDate"`
		CaseOwnerUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(cw),
	}

	// Unmarshal estructural (NO falla por fechas)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Helper de parseo tolerante
	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{} // fecha inválida → zero value
		}
		return t
	}

	// Parseo seguro
	cw.CaseOwnerCreationDate = parse(aux.CaseOwnerCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	cw.CaseOwnerUpdateDate = parse(aux.CaseOwnerUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetCaseOwner inserta un nuevo registro de CaseOwner en la base de datos.
// Realiza la configuración de conexión, genera el SQL correspondiente y ejecuta la inserción.
// Parámetros:
//   - caseOwner: puntero al DTO con los datos a insertar.

//   - connData: datos de conexión actual.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna  un error en caso de fallo.
func SetCaseOwner(caseOwner *CaseOwnerDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se crea un controlador de persistencia para ejecutar operaciones SQL
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Preparación del campo para EntityBranch (si se especifica)
	var branch sql.NullInt64 = sql.NullInt64{}
	if caseOwner.EntityBranch.EntityBranchId > 0 {
		branch.Int64 = int64(caseOwner.EntityBranch.EntityBranchId)
		branch.Valid = true
	}

	// Lista de campos a insertar y sus alias (si aplica)
	var usrProfileFieldsSlice []string = []string{"CaseOwnerICode", "CaseOwnerCreationDate", "CaseOwnerUpdateDate", "CaseOwnerGeneralUser", "CaseOwnerNumCases", "EntityBranch"}
	var usrProfileFieldsAliasSlice []string = []string{}

	// Genera dinámicamente la consulta SQL de inserción utilizando los parámetros y definiciones de campo
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, usrProfileFieldsSlice, usrProfileFieldsAliasSlice, CaseOwnerDBName, []string{}, []string{}, []string{"CaseOwnerId"}, common_dao.SQL_AND, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, false)

	// Asignación de valores por defecto antes de la inserción
	caseOwner.CaseOwnerICode = utils.GetUUID()
	caseOwner.CaseOwnerCreationDate = time.Now()
	caseOwner.CaseOwnerUpdateDate = time.Now()

	// Ejecuta la consulta y captura el ID generado
	persistenceCtrl.QueryRow(context.Background(), query,
		caseOwner.CaseOwnerICode,
		caseOwner.CaseOwnerCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		caseOwner.CaseOwnerUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		caseOwner.CaseOwnerGeneralUser, caseOwner.CaseOwnerNumCases, branch)

	persistenceCtrl.Scan(&caseOwner.CaseOwnerId)

	// Manejo de error en la ejecución de la consulta
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateCaseOwner actualiza los datos de un CaseOwner existente en la base de datos.
// Parámetros:
//   - caseOwner: puntero al DTO con los datos actualizados.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error si ocurre algún fallo.
func UpdateCaseOwner(caseOwner *CaseOwnerDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Prepara el valor de EntityBranch si se especifica
	var branch sql.NullInt64 = sql.NullInt64{}
	if caseOwner.EntityBranch.EntityBranchId > 0 {
		branch.Int64 = int64(caseOwner.EntityBranch.EntityBranchId)
		branch.Valid = true
	}

	// Campos que se actualizarán en la consulta
	var caseOwnerFieldsSlice []string = []string{"CaseOwnerUpdateDate", "CaseOwnerGeneralUser", "CaseOwnerNumCases", "EntityBranch"}
	var caseOwnerFieldsAliasSlice []string = []string{}

	// Genera dinámicamente la consulta SQL de actualización
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, caseOwnerFieldsSlice, caseOwnerFieldsAliasSlice, CaseOwnerDBName, []string{"CaseOwnerId"}, []string{}, []string{}, common_dao.SQL_AND, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, false)

	// Se imprime la consulta para debug (opcional)
	fmt.Printf(query, caseOwner.CaseOwnerUpdateDate, caseOwner.CaseOwnerGeneralUser)

	// Ejecuta la actualización con los parámetros correspondientes
	persistenceCtrl.Exec(context.Background(), query, caseOwner.CaseOwnerId,
		caseOwner.CaseOwnerUpdateDate, caseOwner.CaseOwnerGeneralUser, caseOwner.CaseOwnerNumCases, branch)

	// Manejo de errores y verificación de filas afectadas
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveCaseOwner elimina un registro de CaseOwner de la base de datos.
// Parámetros:
//   - caseOwner: puntero al DTO que identifica el registro a eliminar.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func RemoveCaseOwner(caseOwner *CaseOwnerDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Genera la consulta SQL de eliminación de forma dinámica
	var usrProfileFieldsSlice []string = []string{}
	var usrProfileFieldsAliasSlice []string = []string{}
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, usrProfileFieldsSlice, usrProfileFieldsAliasSlice, CaseOwnerDBName, []string{}, []string{"CaseOwnerId"}, []string{}, common_dao.SQL_AND, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, false)

	// Ejecuta la consulta para eliminar el registro
	persistenceCtrl.Exec(context.Background(), query, caseOwner.CaseOwnerId)

	// Manejo de errores y validación de filas afectadas
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetCaseOwner consulta un registro de CaseOwner a partir de un criterio específico.
// Parámetros:
//   - by: estructura que define el criterio de búsqueda (atributos y operador).
//   - caseOwner: puntero al DTO donde se almacenarán los datos obtenidos.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func GetCaseOwner(by common_controllers.By, caseOwner *CaseOwnerDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construye el path completo de la tabla en la base de datos
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName

	// Configuración de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Selección de los campos que se desean consultar
	var caseOwnerFieldsSlice []string = []string{"CaseOwnerId", "CaseOwnerICode", "CaseOwnerCreationDate", "CaseOwnerUpdateDate", "CaseOwnerGeneralUser", "EntityBranch"}
	var caseOwnerFieldsAliasSlice []string = []string{}
	// Genera la cadena de campos para la consulta
	var caseOwnerFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, caseOwnerFieldsSlice, caseOwnerFieldsAliasSlice, CaseOwnerDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, true)

	// Genera la consulta SQL completa incluyendo la cláusula WHERE basada en el criterio recibido
	var query string = `SELECT ` + caseOwnerFieldsStr +
		` FROM ` + caseOwnerPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, CaseOwnerDBName, by.AttrsName, []string{}, []string{}, by.Operator, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, true)

	// Ejecuta la consulta pasando los valores correspondientes
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Variable para capturar el ID de EntityBranch (puede ser nulo)
	var branchId sql.NullInt64
	// Escanea los resultados en los campos del DTO
	persistenceCtrl.Scan(&caseOwner.CaseOwnerId, &caseOwner.CaseOwnerICode, &caseOwner.CaseOwnerCreationDate, &caseOwner.CaseOwnerUpdateDate,
		&caseOwner.CaseOwnerGeneralUser, &branchId)

	// Asigna el valor de EntityBranch si se obtuvo un valor válido
	if branchId.Valid {
		caseOwner.EntityBranch.EntityBranchId = uint64(branchId.Int64)
	}
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetCaseOwnersByActiveUsers obtiene todos los registros de CaseOwner asociados a usuarios activos.
// La consulta realiza un RIGHT JOIN entre la tabla de CaseOwner y la de GeneralUser.
// Parámetros:
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un slice de CaseOwnerDTO y un error en caso de fallo.
func GetCaseOwnersByActiveUsers(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CaseOwnerDTO, error) {
	// Variable para almacenar temporalmente los datos de la base de datos
	var entity CaseOwnerPgDB
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Construcción de los path completos de las tablas a consultar
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName

	// Configuración de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Selección de los campos de CaseOwner a consultar
	var caseOwnerFieldsSlice []string = []string{"CaseOwnerId", "CaseOwnerICode", "CaseOwnerCreationDate", "CaseOwnerUpdateDate", "CaseOwnerGeneralUser"}
	var caseOwnerFieldsAliasSlice []string = []string{}
	var caseOwnerFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, caseOwnerFieldsSlice, caseOwnerFieldsAliasSlice, CaseOwnerDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, true)

	// Genera la consulta SQL que une CaseOwner y GeneralUser, filtrando solo usuarios activos ('e')
	var query string = `SELECT ` + caseOwnerFieldsStr +
		` FROM ` + caseOwnerPath +
		` RIGHT JOIN ` + userPath + ` ON (` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + ` = ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + `)` +
		` WHERE ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserStatus"].DBName + ` = 'e' ` +
		` ORDER BY ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerNumCases"].DBName + ` ASC `

	// Ejecuta la consulta
	persistenceCtrl.Query(context.Background(), query)

	var entities []CaseOwnerDTO
	// Itera sobre los resultados y convierte cada registro al DTO correspondiente
	for persistenceCtrl.Next() {
		entity = CaseOwnerPgDB{}
		persistenceCtrl.ScanRow(&entity.CaseOwnerICode, &entity.CaseOwnerCreationDate, &entity.CaseOwnerUpdateDate, &entity.CaseOwnerGeneralUser)
		entities = append(entities, entity.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// GetCaseOwnersByActiveUsersAndRoleCode obtiene los registros de CaseOwner asociados a usuarios activos
// que tienen asignado un rol específico identificado por roleCode.
// Realiza JOINs entre las tablas CaseOwner, GeneralUser, RelRoleGeneralUser y Role.
// Parámetros:
//   - roleCode: código del rol a filtrar.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un slice de CaseOwnerDTO y un error en caso de fallo.
func GetCaseOwnersByActiveUsersAndRoleCode(roleCode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CaseOwnerDTO, error) {
	var entity CaseOwnerPgDB
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Construcción de los path completos para las tablas implicadas en la consulta
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName
	var relRolePath string = security_daos.RelRoleGeneralUserDBScheme + "." + security_daos.RelRoleGeneralUserDBName
	var rolePath string = security_daos.RoleDBScheme + "." + security_daos.RoleDBName

	// Configuración de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Selección de los campos de CaseOwner a consultar
	var caseOwnerFieldsSlice []string = []string{"CaseOwnerId", "CaseOwnerICode", "CaseOwnerCreationDate", "CaseOwnerUpdateDate", "CaseOwnerGeneralUser"}
	var caseOwnerFieldsAliasSlice []string = []string{}
	var caseOwnerFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, caseOwnerFieldsSlice, caseOwnerFieldsAliasSlice, CaseOwnerDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CaseOwnerDBScheme, CaseOwnerFieldDefinitions, true)

	// Genera la consulta SQL con JOINs y filtro por roleCode y estado de usuario activo ('e')
	var query string = `SELECT ` + caseOwnerFieldsStr +
		` FROM ` + caseOwnerPath +
		` RIGHT JOIN ` + userPath + ` ON (` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + ` = ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + `)` +
		` LEFT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` LEFT JOIN ` + rolePath + ` ON (` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + ` = ` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + `)` +
		` WHERE ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserStatus"].DBName + ` = 'e' ` + ` AND ` +
		rolePath + `.` + security_daos.RoleFieldDefinitions["RoleCode"].DBName + ` = $1 ` +
		` ORDER BY ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerNumCases"].DBName + ` ASC `

	// Ejecuta la consulta con el parámetro roleCode
	persistenceCtrl.Query(context.Background(), query, roleCode)

	var entities []CaseOwnerDTO
	// Itera sobre los resultados y convierte cada registro a DTO
	for persistenceCtrl.Next() {
		entity = CaseOwnerPgDB{}
		persistenceCtrl.ScanRow(&entity.CaseOwnerId, &entity.CaseOwnerICode, &entity.CaseOwnerCreationDate, &entity.CaseOwnerUpdateDate, &entity.CaseOwnerGeneralUser)
		entities = append(entities, entity.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// GetAllCaseOwners recupera todos los registros de la entidad CaseOwner.
// Parámetros:
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un slice de CaseOwnerDTO y un error en caso de fallo.
func GetAllCaseOwners(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CaseOwnerDTO, error) {
	var entity CaseOwnerPgDB
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración de la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Selección de campos para la consulta
	var entityFieldsSlice []string = []string{"CaseOwnerICode", "CaseOwnerCreationDate", "CaseOwnerUpdateDate", "CaseOwnerGeneralUser"}
	var entityFieldsAliasSlice []string = []string{}

	// Genera la consulta SQL para seleccionar todos los registros
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, entityFieldsSlice, entityFieldsAliasSlice, CaseOwnerDBName, []string{}, []string{}, []string{}, "", CaseOwnerDBScheme, CaseOwnerFieldDefinitions, true)

	// Ejecuta la consulta
	persistenceCtrl.Query(context.Background(), query)
	var entities []CaseOwnerDTO
	// Itera sobre los registros obtenidos y los convierte a DTO
	for persistenceCtrl.Next() {
		entity = CaseOwnerPgDB{}
		persistenceCtrl.ScanRow(&entity.CaseOwnerICode, &entity.CaseOwnerCreationDate, &entity.CaseOwnerUpdateDate, &entity.CaseOwnerGeneralUser)
		entities = append(entities, entity.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// SetCaseOwnerDefaults asigna valores por defecto al DTO de CaseOwner según la acción a realizar.
// Parámetros:
//   - caseOwner: puntero al DTO al que se asignarán los valores por defecto.
//   - action: acción que se está realizando (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetCaseOwnerDefaults(caseOwner *CaseOwnerDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción: se asignan la fecha actual, se genera un UUID y se inicializa el contador de casos.
		caseOwner.CaseOwnerCreationDate = time.Now()
		caseOwner.CaseOwnerUpdateDate = time.Now()
		caseOwner.CaseOwnerICode = utils.GetUUID()
		caseOwner.CaseOwnerNumCases = 0
	case common_dao.SQL_UPDATE:
		// Para actualización: solo se actualiza la fecha de modificación.
		caseOwner.CaseOwnerUpdateDate = time.Now()
	}
}

// ToDTO convierte una instancia de CaseOwnerPgDB (representación de la BD)
// a su equivalente DTO (CaseOwnerDTO) para ser utilizada en la aplicación.
func (obj *CaseOwnerPgDB) ToDTO() CaseOwnerDTO {
	var dto CaseOwnerDTO

	if obj.CaseOwnerId.Valid {
		dto.CaseOwnerId = uint64(obj.CaseOwnerId.Int64)
	}
	if obj.CaseOwnerICode.Valid {
		dto.CaseOwnerICode = obj.CaseOwnerICode.String
	}
	if obj.CaseOwnerCreationDate.Valid {
		dto.CaseOwnerCreationDate = obj.CaseOwnerCreationDate.Time
	}
	if obj.CaseOwnerUpdateDate.Valid {
		dto.CaseOwnerUpdateDate = obj.CaseOwnerUpdateDate.Time
	}
	if obj.CaseOwnerGeneralUser.Valid {
		dto.CaseOwnerGeneralUser = obj.CaseOwnerGeneralUser.String
	}
	if obj.CaseOwnerNumCases.Valid {
		dto.CaseOwnerNumCases = int(obj.CaseOwnerNumCases.Int64)
	}
	if obj.AttentionLine.Valid {
		dto.AttentionLine = AttentionLineDTO{AttentionLineId: uint64(obj.AttentionLine.Int64)}
	}
	if obj.EntityBranch.Valid {
		dto.EntityBranch = EntityBranchDTO{EntityBranchId: uint64(obj.EntityBranch.Int64)}
	}

	return dto
}
