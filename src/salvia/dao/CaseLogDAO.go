// Package salvia_daos contiene los Data Access Objects (DAO) para manejar
// las operaciones relacionadas con el "CaseLog" en la base de datos.
// Este paquete permite insertar, obtener y listar registros de logs de casos.
package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"

	"context"
	"database/sql"
	"fmt"
	"time"
)

var (
	// Nombres y esquemas de la entidad CaseLog
	CaseLogEntityName string = "CaseLog"  // Nombre de la entidad en el sistema
	CaseLogJSONName   string = "caseLog"  // Nombre usado en el JSON
	CaseLogDBName     string = "case_log" // Nombre de la tabla en la base de datos
	CaseLogDBScheme   string = "salvia"   // Esquema de la base de datos en el que se encuentra la tabla

	// Definiciones de campo utilizadas para la validación y mapeo entre JSON, modelo y BD.
	// Cada definición indica el nombre del campo, nombre en la BD, tipo de dato, tamaño mínimo/máximo y si es requerido.
	// Los campos se reciben como string en el JSON, salvo aquellos que necesitan otro tipo en el modelo (como fechas).
	CaseLogFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"CaseLogId":           {Name: "CaseLogId", DBName: "case_log_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"CaseLogICode":        {Name: "CaseLogICode", DBName: "case_log_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"CaseLogCreationDate": {Name: "CaseLogCreationDate", DBName: "case_log_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"CaseLogDescription":  {Name: "CaseLogDescription", DBName: "case_log_description", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 5000, Required: true},
		"CaseLogMoment":       {Name: "CaseLogMoment", DBName: "moment_id", Alias: "", ModelType: "uint", Required: true},
		"CaseLogGeneralUser":  {Name: "CaseLogGeneralUser", DBName: "case_log_general_user", Alias: "", ModelType: "string", MinSize: 32, MaxSize: 36, Required: true},
	}
)

// CaseLogDTO es el Data Transfer Object que representa un registro de CaseLog.
// Se utiliza para el intercambio de datos entre la capa de negocio y la capa de presentación.
type CaseLogDTO struct {
	CaseLogId           uint64    `json:"-"`
	CaseLogICode        string    `json:"icode"`
	CaseLogCreationDate time.Time `json:"creationDate"`
	CaseLogDescription  string    `json:"description"`
	CaseLogMoment       MomentDTO `json:"moment"`
	CaseLogGeneralUser  string    `json:"user"`

	// Campos de formulario que no forman parte directamente del modelo o de la base de datos.
	CaseLogUserRoles   string                              `json:"roles"`
	CaseLogUserProfile security_daos.GeneralUserProfileDTO `json:"profile"`
}

// CaseLogPgDB representa la estructura de datos de CaseLog en la base de datos Postgres.
// Se utilizan tipos sql.Null* para manejar valores nulos.
type CaseLogPgDB struct {
	CaseLogId           sql.NullInt64
	CaseLogICode        sql.NullString
	CaseLogCreationDate sql.NullTime
	CaseLogDescription  sql.NullString
	CaseLogMoment       sql.NullInt64
	CaseLogGeneralUser  sql.NullString
	// Campos de formulario que no forman parte directamente del modelo.
	CaseLogUserRoles sql.NullString
}

// SetCaseLog inserta un nuevo registro de CaseLog en la base de datos.
//
// Parámetros:
//   - caseLog: puntero al DTO que contiene la información del CaseLog a insertar.
//   - connData: datos de conexión actuales a la base de datos.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:

// - error en caso de fallo, o nil si la operación fue exitosa.
func SetCaseLog(caseLog *CaseLogDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Obtiene la conexión a la base de datos utilizando la configuración proporcionada.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos a insertar en la base de datos.
	var caseLogFieldsSlice []string = []string{"CaseLogICode", "CaseLogCreationDate", "CaseLogDescription", "CaseLogGeneralUser", "CaseLogMoment"}
	var caseLogFieldsAliasSlice []string = []string{}

	// Construye la consulta SQL de inserción utilizando la función auxiliar GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, caseLogFieldsSlice, caseLogFieldsAliasSlice, CaseLogDBName, []string{}, []string{}, []string{"CaseLogId"}, common_dao.SQL_AND, CaseLogDBScheme, CaseLogFieldDefinitions, false)

	// Ejecuta la consulta, pasando los valores correspondientes del DTO.
	persistenceCtrl.QueryRow(context.Background(), query,
		caseLog.CaseLogICode,
		caseLog.CaseLogCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		caseLog.CaseLogDescription,
		caseLog.CaseLogGeneralUser,
		caseLog.CaseLogMoment.MomentId)

	// Escanea el valor generado para la clave primaria y lo asigna al DTO.
	persistenceCtrl.Scan(&caseLog.CaseLogId)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetCaseLog obtiene un único registro de CaseLog de la base de datos según el filtro especificado.
//
// Parámetros:
//   - by: estructura que contiene el filtro para la consulta (atributos, alias, operador y valores).
//   - caseLog: puntero al DTO donde se almacenará el resultado.
//   - connData: datos de conexión actuales a la base de datos.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:

// - error en caso de fallo, o nil si la operación fue exitosa.
func GetCaseLog(by common_controllers.By, caseLog *CaseLogDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Define las rutas completas (esquema.tabla) para las tablas involucradas en la consulta.
	var caseLogPath string = CaseLogDBScheme + "." + CaseLogDBName
	var momentPath string = MomentDBScheme + "." + MomentDBName // Se asume que MomentDBScheme y MomentDBName están definidos en otro lugar.
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName
	var profilePath string = security_daos.GeneralUserProfileDBScheme + "." + security_daos.GeneralUserProfileDBName
	var relRolePath string = security_daos.RelRoleGeneralUserDBScheme + "." + security_daos.RelRoleGeneralUserDBName
	var rolePath string = security_daos.RoleDBScheme + "." + security_daos.RoleDBName

	// Obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos a seleccionar para CaseLog.
	var caseLogFieldsSlice []string = []string{"CaseLogId", "CaseLogICode", "CaseLogGeneralUser", "CaseLogCreationDate", "CaseLogDescription", "CaseLogMoment"}
	var caseLogFieldsAliasSlice []string = []string{}
	// Construye la parte de la consulta que selecciona los campos del CaseLog.
	var caseLogFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, caseLogFieldsSlice, caseLogFieldsAliasSlice, CaseLogDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CaseLogDBScheme, CaseLogFieldDefinitions, true)

	// Define los campos a seleccionar para la tabla Moment.
	var momentFieldsSlice []string = []string{"MomentId", "MomentICode", "MomentCode", "MomentVictimCase", "MomentEntityBranch"}
	var momentFieldsAliasSlice []string = []string{}
	var momentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, momentFieldsSlice, momentFieldsAliasSlice, MomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, MomentDBScheme, MomentFieldDefinitions, true)

	// Subconsulta para obtener los roles del usuario asociado al CaseLog.
	var subquery string = `SELECT string_agg( ` + security_daos.RoleFieldDefinitions["RoleName"].DBName + `, ', ' ) as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` WHERE ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + ` = ` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogGeneralUser"].DBName + `)`

	// Construye la consulta SQL completa combinando las partes anteriores y aplicando el filtro recibido.
	var query string = `SELECT ` + caseLogFieldsStr + ", " + momentFieldsStr + ", " + subquery +
		` FROM ` + caseLogPath +
		` LEFT JOIN ` + momentPath + ` ON (` + momentPath + `.` + MomentFieldDefinitions["MomentId"].DBName + ` = ` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogMoment"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + ` = ` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogGeneralUser"].DBName + `)` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, CaseLogDBName, by.AttrsName, []string{}, []string{}, by.Operator, CaseLogDBScheme, CaseLogFieldDefinitions, true)

	// Ejecuta la consulta con los valores del filtro.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Variable temporal para almacenar los datos obtenidos en formato de base de datos.
	var caseLogPG = CaseLogPgDB{}
	persistenceCtrl.Scan(&caseLogPG.CaseLogId, &caseLogPG.CaseLogICode, &caseLogPG.CaseLogCreationDate,
		&caseLogPG.CaseLogDescription, &caseLogPG.CaseLogMoment)

	// Convierte el registro obtenido al DTO y lo asigna al parámetro de salida.
	var tmp CaseLogDTO = caseLogPG.ToDTO()
	*caseLog = tmp

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetCaseLogs obtiene una lista de registros de CaseLog de la base de datos según el filtro especificado.
//
// Parámetros:
//   - by: estructura que contiene el filtro para la consulta.
//   - connData: datos de conexión actuales a la base de datos.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:

// - slice de CaseLogDTO con los registros obtenidos.
// - error en caso de fallo, o nil si la operación fue exitosa.
func GetCaseLogs(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CaseLogDTO, error) {
	// Inicializa el controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var caseLog CaseLogPgDB
	// Define las rutas completas para las tablas involucradas.
	var caseLogPath string = CaseLogDBScheme + "." + CaseLogDBName
	var momentPath string = MomentDBScheme + "." + MomentDBName
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName
	var profilePath string = security_daos.GeneralUserProfileDBScheme + "." + security_daos.GeneralUserProfileDBName
	var relRolePath string = security_daos.RelRoleGeneralUserDBScheme + "." + security_daos.RelRoleGeneralUserDBName
	var rolePath string = security_daos.RoleDBScheme + "." + security_daos.RoleDBName

	// Obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar para CaseLog.
	var caseLogFieldsSlice []string = []string{"CaseLogId", "CaseLogICode", "CaseLogGeneralUser", "CaseLogCreationDate", "CaseLogDescription", "CaseLogMoment"}
	var caseLogFieldsAliasSlice []string = []string{}
	var caseLogFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, caseLogFieldsSlice, caseLogFieldsAliasSlice, CaseLogDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CaseLogDBScheme, CaseLogFieldDefinitions, true)

	// Define los campos a seleccionar para la tabla Moment.
	var momentFieldsSlice []string = []string{"MomentId", "MomentICode", "MomentCode", "MomentVictimCase", "MomentEntityBranch"}
	var momentFieldsAliasSlice []string = []string{}
	var momentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, momentFieldsSlice, momentFieldsAliasSlice, MomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, MomentDBScheme, MomentFieldDefinitions, true)

	// Define los campos a seleccionar para el perfil del usuario.
	var profileFieldsSlice []string = []string{"GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber", "GeneralUserProfileTown"}
	var profileFieldsAliasSlice []string = []string{}
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, security_daos.GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.GeneralUserProfileDBScheme, security_daos.GeneralUserProfileFieldDefinitions, true)

	// Subconsulta para obtener los roles asociados al usuario.
	var subquery string = ` ( SELECT string_agg( ` + security_daos.RoleFieldDefinitions["RoleName"].DBName + `, ', ' ) as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` WHERE ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + ` = ` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogGeneralUser"].DBName + `)`

	// Construye la consulta SQL completa combinando las secciones anteriores y aplicando el filtro.
	var query string = `SELECT ` + caseLogFieldsStr + ", " + momentFieldsStr + ", " + profileFieldsStr + ", " + subquery +
		` FROM ` + caseLogPath +
		` LEFT JOIN ` + momentPath + ` ON (` + momentPath + `.` + MomentFieldDefinitions["MomentId"].DBName + ` = ` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogMoment"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `)` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, CaseLogDBName, by.AttrsName, []string{}, []string{}, by.Operator, CaseLogDBScheme, CaseLogFieldDefinitions, true) +
		` ORDER BY ` + caseLogPath + `.` + CaseLogFieldDefinitions["CaseLogCreationDate"].DBName + ` ASC `

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	var caseLogs []CaseLogDTO
	// Itera sobre los resultados obtenidos.
	for persistenceCtrl.Next() {
		var moment MomentPgDB = MomentPgDB{} // Se asume que MomentPgDB está definido en otro paquete.
		var profile security_daos.GeneralUserProfilePgDB
		caseLog = CaseLogPgDB{}

		// Escanea cada fila y asigna los valores a las variables correspondientes.
		persistenceCtrl.ScanRow(&caseLog.CaseLogId, &caseLog.CaseLogICode, &caseLog.CaseLogGeneralUser, &caseLog.CaseLogCreationDate, &caseLog.CaseLogDescription,
			&caseLog.CaseLogMoment, &moment.MomentId, &moment.MomentICode, &moment.MomentCode, &moment.MomentVictimCase, &moment.MomentEntityBranch,
			&profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate,
			&profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames,
			&profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber, &profile.GeneralUserProfileTown,
			&caseLog.CaseLogUserRoles)

		// Convierte el registro del formato BD a DTO.
		var m CaseLogDTO = caseLog.ToDTO()
		var p security_daos.GeneralUserProfileDTO = profile.ToDTO()
		m.CaseLogMoment = moment.ToDTO()
		m.CaseLogUserProfile = p
		caseLogs = append(caseLogs, m)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return caseLogs, nil
}

// SetCaseLogDefaults asigna valores por defecto al DTO de CaseLog según la acción a realizar.
//
// Parámetros:
//   - caseLog: puntero al DTO de CaseLog que se va a modificar.
//   - action: tipo de acción (por ejemplo, SQL_INSERT o SQL_UPDATE).
//
// En el caso de inserción, se asigna la fecha actual y se genera un UUID para el código.
func SetCaseLogDefaults(caseLog *CaseLogDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Al insertar, se asigna la fecha actual y se genera un identificador único.
		caseLog.CaseLogCreationDate = time.Now()
		caseLog.CaseLogICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Aquí se podrían definir otros valores por defecto para la actualización.
	}
}

// ToDTO convierte una instancia de CaseLogPgDB a su correspondiente CaseLogDTO.
// Esta función permite transformar los datos obtenidos de la base de datos (con tipos sql.Null*)
// a un formato más manejable en la capa de negocio.
func (obj *CaseLogPgDB) ToDTO() CaseLogDTO {
	var dto CaseLogDTO

	if obj.CaseLogId.Valid {
		dto.CaseLogId = uint64(obj.CaseLogId.Int64)
	}

	if obj.CaseLogICode.Valid {
		dto.CaseLogICode = obj.CaseLogICode.String
	}

	if obj.CaseLogCreationDate.Valid {
		dto.CaseLogCreationDate = obj.CaseLogCreationDate.Time
	}

	if obj.CaseLogDescription.Valid {
		dto.CaseLogDescription = obj.CaseLogDescription.String
	}

	if obj.CaseLogGeneralUser.Valid {
		dto.CaseLogGeneralUser = obj.CaseLogGeneralUser.String
	}

	if obj.CaseLogUserRoles.Valid {
		dto.CaseLogUserRoles = obj.CaseLogUserRoles.String
	}

	if obj.CaseLogMoment.Valid {
		dto.CaseLogMoment = MomentDTO{MomentId: uint64(obj.CaseLogMoment.Int64)}
	}

	return dto
}
