// Package security_daos contiene las funciones y estructuras necesarias para gestionar
// los datos de usuarios generales en el contexto de seguridad. Se encarga de las operaciones
// CRUD (Crear, Actualizar, Eliminar y Consultar) sobre la entidad GeneralUser.
package security_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"encoding/json"
	"fmt"
	"strings"

	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	// GeneralUserEntityName es el nombre de la entidad para el usuario general.
	GeneralUserEntityName string = "GeneralUser"
	// GeneralUserJSONName es el nombre del usuario en formato JSON.
	GeneralUserJSONName string = "user"
	// GeneralUserDBName es el nombre de la tabla en la base de datos.
	GeneralUserDBName string = "general_user"
	// GeneralUserDBScheme es el esquema de la base de datos para la seguridad.
	GeneralUserDBScheme string = "security"

	// GeneralUserFieldDefinitions define los atributos y validaciones para cada campo
	// de la entidad GeneralUser, relacionando el nombre del campo, el nombre en la BD,
	// el tipo de dato, tamaño mínimo/máximo y si es requerido.
	// Los booleanos en cada definición indican si el campo se trata como string en el modelo.
	GeneralUserFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"GeneralUserId":                 {Name: "GeneralUserId", DBName: "general_user_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"GeneralUserICode":              {Name: "GeneralUserICode", DBName: "general_user_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"GeneralUserCreationDate":       {Name: "GeneralUserCreationDate", DBName: "general_user_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"GeneralUserUpdateDate":         {Name: "GeneralUserUpdateDate", DBName: "general_user_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"GeneralUserLogin":              {Name: "GeneralUserLogin", DBName: "general_user_login", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 128, Required: true},
		"GeneralUserPassword":           {Name: "GeneralUserPassword", DBName: "general_user_password", Alias: "", ModelType: "string", MinSize: 4, MaxSize: 512, Required: true},
		"GeneralUserStatus":             {Name: "GeneralUserStatus", DBName: "general_user_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: false},
		"GeneralUserLanguage":           {Name: "GeneralUserLanguage", DBName: "general_user_language", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"GeneralUserGeneralUserProfile": {Name: "GeneralUserGeneralUserProfile", DBName: "general_user_general_user_profile", Alias: "", ModelType: "uint", Required: false},
		"GeneralUserTeam":               {Name: "GeneralUserTeam", DBName: "general_user_team", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 50, Required: false},
		"GeneralUserAssignedDepartment": {Name: "GeneralUserAssignedDepartment", DBName: "general_user_assigned_department", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 20, Required: false},
	}
)

// GeneralUserDTO es la estructura de datos que representa la entidad GeneralUser en las operaciones
// de entrada y salida de la aplicación, y se utiliza para el manejo de datos en formato JSON.
type GeneralUserDTO struct {
	GeneralUserId                 uint64                `json:"-"`
	GeneralUserICode              string                `json:"icode"`
	GeneralUserCreationDate       time.Time             `json:"creationDate"`
	GeneralUserUpdateDate         time.Time             `json:"updateDate"`
	GeneralUserLogin              string                `json:"login"`
	GeneralUserPassword           string                `json:"pass"`
	GeneralUserStatus             string                `json:"status"`
	GeneralUserLanguage           string                `json:"lang"`
	GeneralUserGeneralUserProfile GeneralUserProfileDTO `json:"profile"`
	// Campos de formulario que no hacen parte del modelo o no directamente en la BD
	GeneralUserPasswordRepeat     string    `json:"pass2"`
	GeneralUserRoles              []RoleDTO `json:"roles"`
	GeneralUserRoleIds            []string  `json:"roleIds"`
	GeneralUserRoleCodes          []string  `json:"roleCodes"`
	GeneralUserFullName           string    `json:"fullName"`
	GeneralUserResetPasswordICode string    `json:"resetToken"`
	GeneralUserCaptchaID          string    `json:"captchaID"`
	GeneralUserCaptchaSolution    string    `json:"captchaSolution"`
	GeneralUserTeam               string    `json:"team"`
	GeneralUserAssignedDepartment string    `json:"assignedDepartment"`
}

// GeneralUserPgDB es la estructura que representa el modelo de datos de GeneralUser tal
// como se almacena en la base de datos, utilizando tipos sql.Null* para gestionar valores nulos.
type GeneralUserPgDB struct {
	GeneralUserId                 sql.NullInt64
	GeneralUserICode              sql.NullString
	GeneralUserCreationDate       sql.NullTime
	GeneralUserUpdateDate         sql.NullTime
	GeneralUserLogin              sql.NullString
	GeneralUserPassword           sql.NullString
	GeneralUserStatus             sql.NullString
	GeneralUserLanguage           sql.NullString
	GeneralUserGeneralUserProfile sql.NullInt64
	GeneralUserTeam              sql.NullString
	GeneralUserAssignedDepartment sql.NullString
	// Campos de formulario que no hacen parte del modelo en la BD
	GeneralUserPasswordRepeat string           `json:"pass2"`
	GeneralUserRoles          []sql.NullInt64  `json:"roles"`
	GeneralUserRoleIds        []sql.NullString `json:"roleIds"`
	GeneralUserRoleCodes      []sql.NullString `json:"roleCodes"`
}

func (u GeneralUserDTO) MarshalJSON() ([]byte, error) {
	type Alias GeneralUserDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		GeneralUserCreationDate string `json:"creationDate"`
		GeneralUserUpdateDate   string `json:"updateDate"`
	}{
		Alias:                   (*Alias)(&u),
		GeneralUserCreationDate: u.GeneralUserCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		GeneralUserUpdateDate:   u.GeneralUserUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
	})
}

func (u *GeneralUserDTO) UnmarshalJSON(data []byte) error {
	type Alias GeneralUserDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		GeneralUserCreationDate string `json:"creationDate"`
		GeneralUserUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(u),
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
	u.GeneralUserCreationDate = parse(aux.GeneralUserCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	u.GeneralUserUpdateDate = parse(aux.GeneralUserUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetGeneralUser inserta un nuevo usuario general en la base de datos.
// Recibe el objeto usuario, información de transacción, módulo, y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func SetGeneralUser(user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se asigna el identificador de perfil, validando si es 0 (se asigna nil en ese caso)
	var profileId interface{} = user.GeneralUserGeneralUserProfile.GeneralUserProfileId
	if user.GeneralUserGeneralUserProfile.GeneralUserProfileId == 0 {
		profileId = nil
	}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se define la lista de campos que se insertarán en la BD.
	var usrFieldsSlice []string = []string{"GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserPassword", "GeneralUserStatus", "GeneralUserLanguage", "GeneralUserGeneralUserProfile", "GeneralUserAssignedDepartment"}
	var usrFieldsAliasSlice []string = []string{}

	// Se genera la query de inserción utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{"GeneralUserId"}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query, user.GeneralUserICode, user.GeneralUserCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), user.GeneralUserUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), user.GeneralUserLogin, user.GeneralUserPassword, user.GeneralUserStatus, user.GeneralUserLanguage, profileId, user.GeneralUserAssignedDepartment)
	persistenceCtrl.Scan(&user.GeneralUserId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateGeneralUserByICode actualiza los datos de un usuario general en la BD, utilizando su ICode como referencia.
// Recibe el objeto usuario con los nuevos datos, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func UpdateGeneralUserByICode(user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se actualizarán.
	var usrFieldsSlice []string = []string{"GeneralUserUpdateDate", "GeneralUserPassword", "GeneralUserLanguage", "GeneralUserStatus", "GeneralUserAssignedDepartment"}
	var usrFieldsAliasSlice []string = []string{}

	// Se genera la query de actualización utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{"GeneralUserICode"}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		user.GeneralUserICode, user.GeneralUserUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), user.GeneralUserPassword, user.GeneralUserLanguage, user.GeneralUserStatus, user.GeneralUserAssignedDepartment)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se valida que se hayan afectado filas en la actualización.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveGeneralUserByICode elimina un usuario general de la base de datos utilizando su ICode.
// Recibe el objeto usuario, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func RemoveGeneralUserByICode(user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos y condiciones para la eliminación.
	var usrFieldsSlice []string = []string{}
	var usrFieldsAliasSlice []string = []string{}

	// Se genera la query de eliminación utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{"GeneralUserICode"}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, false)

	// Se ejecuta la query con el parámetro correspondiente.
	persistenceCtrl.Exec(context.Background(), query, user.GeneralUserICode)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se valida que se hayan afectado filas en la eliminación.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetGeneralUserLogin consulta el usuario general basado en atributos de búsqueda relacionados con el login.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto usuario a completar,
// información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func GetGeneralUserLogin(by common_controllers.By, user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla GeneralUser en la BD.
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var usrFieldsSlice []string = []string{"GeneralUserId", "GeneralUserPassword"}
	var usrFieldsAliasSlice []string = []string{}
	var usrFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + usrFieldsStr +
		` FROM ` + userPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, GeneralUserDBName, by.AttrsName, []string{}, []string{}, by.Operator, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)

	// Se ejecuta la query con los parámetros de filtrado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se preparan variables para capturar los datos relacionados con perfil, municipio y roles.
	var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
	var town TownPgDB = TownPgDB{}
	var rolesPg sql.NullString
	var userPg = GeneralUserPgDB{}

	// Se escanean los resultados de la query.
	persistenceCtrl.Scan(&userPg.GeneralUserId, &userPg.GeneralUserPassword)

	// Se procesan los roles (si existen) separándolos por comas.
	var rolesSlice []string
	if rolesPg.Valid {
		rolesSlice = strings.Split(rolesPg.String, ",")
	}

	// Se convierte el objeto de base de datos a DTO.
	*user = userPg.ToDTO()
	// Se asignan los roles al DTO.
	for _, r := range rolesSlice {
		var roleFields []string = strings.Split(r, "-")
		if len(roleFields) == 2 {
			user.GeneralUserRoles = append(user.GeneralUserRoles, RoleDTO{RoleName: roleFields[1], RoleCode: roleFields[0]})
		}
	}

	// Se asigna el perfil y el municipio al DTO.
	user.GeneralUserGeneralUserProfile = profile.ToDTO()
	user.GeneralUserGeneralUserProfile.GeneralUserProfileTown = town.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetGeneralUser consulta todos los datos asociados a un usuario general, incluyendo información de perfil,
// municipio, roles y datos de otras entidades relacionadas.
// Recibe un objeto 'By' para la búsqueda, el objeto usuario a completar, información de transacción,
// módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func GetGeneralUser(by common_controllers.By, user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construyen los paths para las diferentes tablas involucradas.
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName
	var townPath string = TownDBScheme + "." + TownDBName
	var caseOwnerPath string = "salvia.case_owner"
	var entityBranchPath string = "salvia.entity_branch"
	var relRolePath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a consultar para el usuario, el perfil y el municipio.
	var usrFieldsSlice []string = []string{"GeneralUserId", "GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserPassword", "GeneralUserStatus", "GeneralUserLanguage", "GeneralUserGeneralUserProfile", "GeneralUserTeam", "GeneralUserAssignedDepartment"}
	var usrFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate",
		"GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames",
		"GeneralUserProfileDocType", "GeneralUserProfileDocNumber", "GeneralUserProfileTown"}
	var profileFieldsAliasSlice []string = []string{}

	var townFieldsSlice []string = []string{"TownId", "TownICode", "TownCode", "TownName", "TownCity", "TownType"}
	var townFieldsAliasSlice []string = []string{}

	// Se generan los strings de campos para cada tabla utilizando GetSQL.
	var usrFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)
	var townFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, townFieldsSlice, townFieldsAliasSlice, TownDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, TownDBScheme, TownFieldDefinitions, true)

	// Se define una subquery para obtener los roles del usuario agregándolos en una cadena.
	var subquery string = ` ( SELECT string_agg( ` + RoleFieldDefinitions["RoleCode"].DBName + ` ||'-'||` + RoleFieldDefinitions["RoleName"].DBName + `, ', ' ) as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, GeneralUserDBName, by.AttrsName, []string{}, []string{}, by.Operator, GeneralUserDBScheme, GeneralUserFieldDefinitions, true) + `)`

	// Se genera la query completa con joins a las tablas de perfil, municipio, case_owner y entity_branch.
	var query string = `SELECT ` + usrFieldsStr + `, ` + profileFieldsStr + `, ` + townFieldsStr + `, ` + entityBranchPath + `.entity_branch_i_code` + ", " + subquery +
		` FROM ` + userPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +
		` LEFT JOIN ` + townPath + ` ON (` + townPath + `.` + TownFieldDefinitions["TownCode"].DBName + ` = ` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileTown"].DBName + `)` +
		` LEFT JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.case_owner_general_user` + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `)` +
		` LEFT JOIN ` + entityBranchPath + ` ON (` + entityBranchPath + `.entity_branch_id` + ` = ` + caseOwnerPath + `.entity_branch_id )` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, GeneralUserDBName, by.AttrsName, []string{}, []string{}, by.Operator, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)

	// Se ejecuta la query con el parámetro de filtrado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se instancian variables para almacenar los datos del perfil, municipio y roles.
	var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
	var town TownPgDB = TownPgDB{}
	var rolesPg sql.NullString
	var userPg = GeneralUserPgDB{}

	// Se escanean los resultados de la query.
	persistenceCtrl.Scan(&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
		&userPg.GeneralUserPassword, &userPg.GeneralUserStatus, &userPg.GeneralUserLanguage, &profile.GeneralUserProfileId, &userPg.GeneralUserTeam, &userPg.GeneralUserAssignedDepartment,
		&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate,
		&profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames,
		&profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber, &profile.GeneralUserProfileTown,
		&town.TownId, &town.TownICode, &town.TownCode, &town.TownName, &town.TownCity, &town.TownType, &profile.GeneralUserProfileEntityBranchSelected,
		&rolesPg)

	// Se procesan los roles obtenidos en la subquery.
	var rolesSlice []string
	if rolesPg.Valid {
		rolesSlice = strings.Split(rolesPg.String, ",")
	}

	// Se convierte el objeto de base de datos a DTO y se asignan los roles.
	*user = userPg.ToDTO()
	for _, r := range rolesSlice {
		var roleFields []string = strings.Split(r, "-")
		if len(roleFields) == 2 {
			user.GeneralUserRoles = append(user.GeneralUserRoles, RoleDTO{RoleName: roleFields[1], RoleCode: roleFields[0]})
		}
	}

	// Se asignan el perfil y el municipio al DTO.
	user.GeneralUserGeneralUserProfile = profile.ToDTO()
	user.GeneralUserGeneralUserProfile.GeneralUserProfileTown = town.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetAllGeneralUser consulta todos los usuarios generales existentes en la base de datos.
// Recibe información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada, un slice de GeneralUserDTO y un error en caso de producirse.
func GetAllGeneralUser(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]GeneralUserDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar para los usuarios.
	var usrFieldsSlice []string = []string{"GeneralUserId", "GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserStatus", "GeneralUserLanguage", "GeneralUserGeneralUserProfile"}
	var usrFieldsAliasSlice []string = []string{}
	var usrFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)

	// Se genera la query para obtener todos los usuarios.
	var query string = `SELECT ` + usrFieldsStr +
		` FROM ` + userPath +
		` WHERE TRUE`

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query)
	var users []GeneralUserDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var roles string
		var userPg = GeneralUserPgDB{}
		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
			&userPg.GeneralUserStatus, &userPg.GeneralUserLanguage, &profile.GeneralUserProfileId,
			&roles)
		// Se separan y asignan los roles al DTO.
		var rolesSlice []string = strings.Split(roles, ",")
		var user GeneralUserDTO = userPg.ToDTO()
		for _, r := range rolesSlice {
			user.GeneralUserRoles = append(user.GeneralUserRoles, RoleDTO{RoleName: r})
		}
		// Se asigna el perfil convertido a DTO.
		user.GeneralUserGeneralUserProfile = profile.ToDTO()
		users = append(users, user)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return users, nil
}

// GetAllGeneralUserWithProfile consulta todos los usuarios generales junto con sus perfiles asociados.
// Recibe información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada, un slice de GeneralUserDTO y un error en caso de producirse.
func GetAllGeneralUserWithProfile(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]GeneralUserDTO, int, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName
	var relRolePath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName
	var count int

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se generan los strings de campos para el usuario y el perfil.
	var usrFieldsStr = "u." + GeneralUserFieldDefinitions["GeneralUserId"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserICode"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserCreationDate"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserUpdateDate"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserLogin"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserStatus"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserLanguage"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName
	var profileFieldsStr = "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileICode"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileCreationDate"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileUpdateDate"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileGender"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileNick"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileDescription"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileNames"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileLastNames"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileDocType"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileDocNumber"].DBName
	// Se define una subquery para obtener los roles asociados al usuario.
	/*
		SELECT
			u.general_user_id,
			u.general_user_i_code,
			u.general_user_creation_date,
			u.general_user_update_date,
			u.general_user_login,
			u.general_user_status,
			u.general_user_language,
			u.general_user_general_user_profile,

			p.general_user_profile_id,
			p.general_user_profile_i_code,
			p.general_user_profile_creation_date,
			p.general_user_profile_update_date,
			p.general_user_profile_gender,
			p.general_user_profile_nick,
			p.general_user_profile_description,
			p.general_user_profile_names,
			p.general_user_profile_last_names,
			p.general_user_profile_doc_type,
			p.general_user_profile_doc_number,

			r.roles

			FROM security.general_user u

			LEFT JOIN security.general_user_profile p
			ON p.general_user_profile_id = u.general_user_general_user_profile

			LEFT JOIN (
				SELECT
					rgu.general_user_id,
					string_agg(role.role_code || '-' || role.role_name, ', ') AS roles
				FROM security.rel_role_general_user rgu
				JOIN security.role role
					ON role.role_id = rgu.role_id
				GROUP BY rgu.general_user_id
			) r
			ON r.general_user_id = u.general_user_id

			LIMIT 50 OFFSET 13500;
	*/
	var subquery string = ` LEFT JOIN ( SELECT rgu.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `,  string_agg( ` + RoleFieldDefinitions["RoleCode"].DBName + ` || '-' || ` + RoleFieldDefinitions["RoleName"].DBName + `, ', ' ) AS roles ` +
		` FROM ` + relRolePath + ` rgu ` +
		` JOIN ` + rolePath + ` role ON (rgu.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = role.` + RoleFieldDefinitions["RoleId"].DBName + `)` +
		` GROUP BY rgu.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + `) r ON (r.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = u.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)`

	// Se genera la query completa.
	var query string = `SELECT ` + usrFieldsStr + ", " + profileFieldsStr + ", r.roles" +
		` FROM ` + userPath + ` u ` +
		` LEFT JOIN ` + profilePath + ` p ON ( p.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = u.` + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `) ` +
		subquery +
		common_dao.GetOffsetQuery(page)

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query)
	var users []GeneralUserDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var userPg = GeneralUserPgDB{}
		var roles sql.NullString

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
			&userPg.GeneralUserStatus, &userPg.GeneralUserLanguage, &profile.GeneralUserProfileId,
			&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber,
			&roles)

		var rolesSlice []string
		var user GeneralUserDTO = userPg.ToDTO()

		// Se procesan los roles y se agregan al DTO.
		if roles.Valid {
			rolesSlice = strings.Split(roles.String, ",")
			for _, r := range rolesSlice {
				if r != "" {
					var roleFields []string = strings.Split(r, "-")
					if len(roleFields) == 2 {
						user.GeneralUserRoles = append(user.GeneralUserRoles, RoleDTO{RoleName: roleFields[1], RoleCode: roleFields[0]})
					}
				}
			}
		}

		// Se asigna el perfil convertido a DTO.
		user.GeneralUserGeneralUserProfile = profile.ToDTO()
		users = append(users, user)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + userPath + ` u ` +
			subquery

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return users, count, nil
}

func GetGeneralUserWithProfileByStatus(status string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]GeneralUserDTO, int, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName
	var relRolePath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName
	var count int

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se generan los strings de campos para el usuario y el perfil.
	var usrFieldsStr = "u." + GeneralUserFieldDefinitions["GeneralUserId"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserICode"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserCreationDate"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserUpdateDate"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserLogin"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserStatus"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserLanguage"].DBName + "," + "u." + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName
	var profileFieldsStr = "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileICode"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileCreationDate"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileUpdateDate"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileGender"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileNick"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileDescription"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileNames"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileLastNames"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileDocType"].DBName + "," + "p." + GeneralUserProfileFieldDefinitions["GeneralUserProfileDocNumber"].DBName
	// Se define una subquery para obtener los roles asociados al usuario.
	/*
		SELECT
			u.general_user_id,
			u.general_user_i_code,
			u.general_user_creation_date,
			u.general_user_update_date,
			u.general_user_login,
			u.general_user_status,
			u.general_user_language,
			u.general_user_general_user_profile,

			p.general_user_profile_id,
			p.general_user_profile_i_code,
			p.general_user_profile_creation_date,
			p.general_user_profile_update_date,
			p.general_user_profile_gender,
			p.general_user_profile_nick,
			p.general_user_profile_description,
			p.general_user_profile_names,
			p.general_user_profile_last_names,
			p.general_user_profile_doc_type,
			p.general_user_profile_doc_number,

			r.roles

			FROM security.general_user u

			LEFT JOIN security.general_user_profile p
			ON p.general_user_profile_id = u.general_user_general_user_profile

			LEFT JOIN (
				SELECT
					rgu.general_user_id,
					string_agg(role.role_code || '-' || role.role_name, ', ') AS roles
				FROM security.rel_role_general_user rgu
				JOIN security.role role
					ON role.role_id = rgu.role_id
				GROUP BY rgu.general_user_id
			) r
			ON r.general_user_id = u.general_user_id

			WHERE u.general_user_status = 'e'

			LIMIT 50 OFFSET 13500;
	*/
	var subquery string = ` LEFT JOIN ( SELECT rgu.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `,  string_agg( ` + RoleFieldDefinitions["RoleCode"].DBName + ` || '-' || ` + RoleFieldDefinitions["RoleName"].DBName + `, ', ' ) AS roles ` +
		` FROM ` + relRolePath + ` rgu ` +
		` JOIN ` + rolePath + ` role ON (rgu.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = role.` + RoleFieldDefinitions["RoleId"].DBName + `)` +
		` GROUP BY rgu.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + `) r ON (r.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = u.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)`

	// Se genera la query completa.
	var query string = `SELECT ` + usrFieldsStr + ", " + profileFieldsStr + ", r.roles" +
		` FROM ` + userPath + ` u ` +
		` LEFT JOIN ` + profilePath + ` p ON ( p.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = u.` + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `) ` +
		subquery +
		` WHERE u.` + GeneralUserFieldDefinitions["GeneralUserStatus"].DBName + ` = $1 ` +
		common_dao.GetOffsetQuery(page)

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query, status)
	var users []GeneralUserDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var userPg = GeneralUserPgDB{}
		var roles sql.NullString

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
			&userPg.GeneralUserStatus, &userPg.GeneralUserLanguage, &profile.GeneralUserProfileId,
			&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber,
			&roles)

		var rolesSlice []string
		var user GeneralUserDTO = userPg.ToDTO()

		// Se procesan los roles y se agregan al DTO.
		if roles.Valid {
			rolesSlice = strings.Split(roles.String, ",")
			for _, r := range rolesSlice {
				if r != "" {
					var roleFields []string = strings.Split(r, "-")
					if len(roleFields) == 2 {
						user.GeneralUserRoles = append(user.GeneralUserRoles, RoleDTO{RoleName: roleFields[1], RoleCode: roleFields[0]})
					}
				}
			}
		}

		// Se asigna el perfil convertido a DTO.
		user.GeneralUserGeneralUserProfile = profile.ToDTO()
		users = append(users, user)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + userPath + ` u ` +
			subquery +
			` WHERE u.` + GeneralUserFieldDefinitions["GeneralUserStatus"].DBName + ` = $1 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, status)
		persistenceCtrl.Scan(&count)
	}

	return users, count, nil
}

// GetGeneralUsersByRoleWithProfile consulta los usuarios generales asociados a un determinado rol, incluyendo su perfil.
// Recibe el código del rol, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada, un slice de GeneralUserDTO y un error en caso de producirse.
func GetGeneralUsersByRoleWithProfile(roleCode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]GeneralUserDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName
	var relRolePath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar para el usuario y el perfil.
	var usrFieldsSlice []string = []string{"GeneralUserId", "GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserStatus", "GeneralUserLanguage"}
	var usrFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Se generan los strings de campos.
	var usrFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Se genera la query con joins a las tablas de perfil, rol y relación de roles.
	var query string = `SELECT ` + usrFieldsStr + ", " + profileFieldsStr +
		` FROM ` + userPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +
		` LEFT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` LEFT JOIN ` + rolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + `)` +
		` WHERE  ` + rolePath + `.` + RoleFieldDefinitions["RoleCode"].DBName + ` = $1`

	// Se ejecuta la query utilizando el código del rol.
	persistenceCtrl.Query(context.Background(), query, roleCode)
	var users []GeneralUserDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var userPg = GeneralUserPgDB{}

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
			&userPg.GeneralUserStatus, &userPg.GeneralUserLanguage,
			&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

		// Se convierte el objeto de base de datos a DTO y se asigna el perfil.
		var user GeneralUserDTO = userPg.ToDTO()
		user.GeneralUserGeneralUserProfile = profile.ToDTO()
		// Se construye el nombre completo a partir de los nombres y apellidos.
		user.GeneralUserFullName = user.GeneralUserGeneralUserProfile.GeneralUserProfileNames + " " + user.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames

		users = append(users, user)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return users, nil
}

// GetDepartmentGeneralUsersByTownCode consulta los usuarios generales asociados a un determinado departamento, incluyendo su perfil.
// Recibe el código del centro urbano, información de transacción, módulo y datos de conexión.
// Retorna un slice de GeneralUserDTO y un error en caso de producirse.
func GetGeneralUsersByDepartmentICodeAndRoleCode(roleCode string, departmentICode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]GeneralUserDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var userPath string = GeneralUserDBScheme + "." + GeneralUserDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName
	var relRolePath string = RelRoleGeneralUserDBScheme + "." + RelRoleGeneralUserDBName
	var rolePath string = RoleDBScheme + "." + RoleDBName
	var townPath string = TownDBScheme + "." + TownDBName
	var cityPath string = CityDBScheme + "." + CityDBName
	var departmentPath string = DepartmentDBScheme + "." + DepartmentDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	/* SQL
	SELECT general_user.*
	FROM security.general_user
	LEFT JOIN security.general_user_profile ON ( general_user.general_user_general_user_profile = general_user_profile.general_user_profile_id )
	LEFT JOIN security.town ON (town.town_code =  general_user_profile.general_user_profile_town)
	LEFT JOIN security.city ON (security.city.city_id = town.city_id)
	LEFT JOIN security.department ON (security.department.department_id = city.city_department_id)
	LEFT JOIN security.rel_role_general_user ON (rel_role_general_user.general_user_id = general_user.general_user_id)
	LEFT JOIN security.role ON (role.role_id = rel_role_general_user.role_id)
	WHERE town.town_code = '11001000' AND role.role_code = 'et'


	*/
	// Se definen los campos a consultar para el usuario y el perfil.
	var usrFieldsSlice []string = []string{"GeneralUserId", "GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserStatus", "GeneralUserLanguage"}
	var usrFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Se generan los strings de campos.
	var usrFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, usrFieldsSlice, usrFieldsAliasSlice, GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserDBScheme, GeneralUserFieldDefinitions, true)
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Se genera la query con joins a las tablas de perfil, rol y relación de roles.
	var query string = `SELECT ` + usrFieldsStr + ", " + profileFieldsStr +
		` FROM ` + userPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +
		` LEFT JOIN ` + townPath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileTown"].DBName + ` = ` + townPath + `.` + TownFieldDefinitions["TownCode"].DBName + `)` +
		` LEFT JOIN ` + cityPath + ` ON (` + cityPath + `.` + CityFieldDefinitions["CityId"].DBName + ` = ` + townPath + `.` + TownFieldDefinitions["TownCity"].DBName + `)` +
		` LEFT JOIN ` + departmentPath + ` ON (` + departmentPath + `.` + DepartmentFieldDefinitions["DepartmentId"].DBName + ` = ` + cityPath + `.` + CityFieldDefinitions["CityDepartment"].DBName + `)` +
		` LEFT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` LEFT JOIN ` + rolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + `)` +
		` WHERE  ` + rolePath + `.` + RoleFieldDefinitions["RoleCode"].DBName + ` = $1 AND ` + departmentPath + `.` + DepartmentFieldDefinitions["DepartmentICode"].DBName + ` = $2 `

	// Se ejecuta la query utilizando el código del rol.
	persistenceCtrl.Query(context.Background(), query, roleCode, departmentICode)
	var users []GeneralUserDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var userPg = GeneralUserPgDB{}

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
			&userPg.GeneralUserStatus, &userPg.GeneralUserLanguage,
			&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

		// Se convierte el objeto de base de datos a DTO y se asigna el perfil.
		var user GeneralUserDTO = userPg.ToDTO()
		user.GeneralUserGeneralUserProfile = profile.ToDTO()
		// Se construye el nombre completo a partir de los nombres y apellidos.
		user.GeneralUserFullName = user.GeneralUserGeneralUserProfile.GeneralUserProfileNames + " " + user.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames

		users = append(users, user)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return users, nil
}

// SetGeneralUserDefaults asigna valores por defecto a los campos de un usuario general,
// dependiendo de la acción que se esté realizando (insertar o actualizar).
func SetGeneralUserDefaults(user *GeneralUserDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción se asigna estado 'e', fechas actuales y se genera un UUID para el ICode.
		user.GeneralUserStatus = "e"
		user.GeneralUserCreationDate = time.Now()
		user.GeneralUserUpdateDate = time.Now()
		user.GeneralUserICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para actualización solo se actualiza la fecha de modificación.
		user.GeneralUserUpdateDate = time.Now()
	}
}

// ToDTO convierte un objeto GeneralUserPgDB obtenido de la base de datos en un objeto GeneralUserDTO.
// Se encarga de verificar la validez de los campos nulos y asignarlos correctamente.
func (obj *GeneralUserPgDB) ToDTO() GeneralUserDTO {
	var dto GeneralUserDTO

	if obj.GeneralUserId.Valid {
		dto.GeneralUserId = uint64(obj.GeneralUserId.Int64)
	}

	if obj.GeneralUserICode.Valid {
		dto.GeneralUserICode = obj.GeneralUserICode.String
	}

	if obj.GeneralUserCreationDate.Valid {
		dto.GeneralUserCreationDate = obj.GeneralUserCreationDate.Time
	}

	if obj.GeneralUserUpdateDate.Valid {
		dto.GeneralUserUpdateDate = obj.GeneralUserUpdateDate.Time
	}

	if obj.GeneralUserLogin.Valid {
		dto.GeneralUserLogin = obj.GeneralUserLogin.String
	}

	if obj.GeneralUserPassword.Valid {
		dto.GeneralUserPassword = obj.GeneralUserPassword.String
	}

	if obj.GeneralUserStatus.Valid {
		dto.GeneralUserStatus = obj.GeneralUserStatus.String
	}

	if obj.GeneralUserLanguage.Valid {
		dto.GeneralUserLanguage = obj.GeneralUserLanguage.String
	}

	if obj.GeneralUserGeneralUserProfile.Valid {
		dto.GeneralUserGeneralUserProfile = GeneralUserProfileDTO{GeneralUserProfileId: uint64(obj.GeneralUserGeneralUserProfile.Int64)}
	}

	if obj.GeneralUserTeam.Valid {
		dto.GeneralUserTeam = obj.GeneralUserTeam.String
	}

	if obj.GeneralUserAssignedDepartment.Valid {
		dto.GeneralUserAssignedDepartment = obj.GeneralUserAssignedDepartment.String
	}

	return dto
}
