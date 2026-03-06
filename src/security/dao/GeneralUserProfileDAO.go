// Package security_daos contiene las definiciones y funciones para el manejo de
// datos relacionados con el perfil de usuario general en el sistema de seguridad.
package security_daos

import (
	// Configuraciones comunes
	// Controladores comunes para persistencia y lógica de negocio
	// Funciones y utilidades para acceso a datos
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"    // Conexión y configuración de base de datos
	"bitsflow/common/utils" // Funciones utilitarias generales (por ejemplo, para generar UUID)
	"encoding/json"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	// Nombres y esquemas de la entidad GeneralUserProfile en diferentes contextos.
	GeneralUserProfileEntityName string = "GeneralUserProfile"   // Nombre de la entidad
	GeneralUserProfileJSONName   string = "profile"              // Nombre en JSON
	GeneralUserProfileDBName     string = "general_user_profile" // Nombre en la base de datos
	GeneralUserProfileDBScheme   string = "security"             // Esquema en la base de datos

	// Definición de campos para validaciones, mapeando los nombres de campo a sus propiedades.
	// Cada campo especifica nombre, nombre en BD, tipo de modelo, tamaños mínimos/máximos y si es requerido.
	GeneralUserProfileFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"GeneralUserProfileId":           {Name: "GeneralUserProfileId", DBName: "general_user_profile_id", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 0, Required: false},
		"GeneralUserProfileICode":        {Name: "GeneralUserProfileICode", DBName: "general_user_profile_i_code", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 36, Required: false},
		"GeneralUserProfileCreationDate": {Name: "GeneralUserProfileCreationDate", DBName: "general_user_profile_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"GeneralUserProfileUpdateDate":   {Name: "GeneralUserProfileUpdateDate", DBName: "general_user_profile_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"GeneralUserProfileGender":       {Name: "GeneralUserProfileGender", DBName: "general_user_profile_gender", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"GeneralUserProfileNick":         {Name: "GeneralUserProfileNick", DBName: "general_user_profile_nick", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 120, Required: false},
		"GeneralUserProfileDescription":  {Name: "GeneralUserProfileDescription", DBName: "general_user_profile_description", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 1024, Required: false},
		"GeneralUserProfileNames":        {Name: "GeneralUserProfileNames", DBName: "general_user_profile_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 128, Required: true},
		"GeneralUserProfileLastNames":    {Name: "GeneralUserProfileLastNames", DBName: "general_user_profile_last_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 128, Required: false},
		"GeneralUserProfileDocType":      {Name: "GeneralUserProfileDocType", DBName: "general_user_profile_doc_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"GeneralUserProfileDocNumber":    {Name: "GeneralUserProfileDocNumber", DBName: "general_user_profile_doc_number", Alias: "", ModelType: "string", MinSize: 5, MaxSize: 32, Required: false},
		"GeneralUserProfileTown":         {Name: "GeneralUserProfileTown", DBName: "general_user_profile_town", Alias: "", ModelType: "string", MinSize: 5, MaxSize: 8, Required: true},
	}
)

// GeneralUserProfileDTO representa el Data Transfer Object (DTO) para el perfil de usuario general.
// Se utiliza para el intercambio de datos entre la capa de persistencia y la lógica de negocio o presentación.
type GeneralUserProfileDTO struct {
	GeneralUserProfileId           uint64           `json:"-"`
	GeneralUserProfileICode        string           `json:"icode"`
	GeneralUserProfileCreationDate time.Time        `json:"creationDate"`
	GeneralUserProfileUpdateDate   time.Time        `json:"updateDate"`
	GeneralUserProfileGender       string           `json:"gender"`
	GeneralUserProfileNick         string           `json:"nick"`
	GeneralUserProfileDescription  string           `json:"description"`
	GeneralUserProfileNames        string           `json:"names"`
	GeneralUserProfileLastNames    string           `json:"lastNames"`
	GeneralUserProfileDocType      string           `json:"docType"`
	GeneralUserProfileDocNumber    string           `json:"docNumber"`
	GeneralUserProfileEmails       []EMailDTO       `json:"emails"`
	GeneralUserProfilePhoneNumbers []PhoneNumberDTO `json:"phoneNumbers"`
	GeneralUserProfileTown         TownDTO          `json:"town"`
	//Campos de formulario que no hacen parte del modelo o no direcxtamente en la BD -----------------------------------
	GeneralUserProfileTownCode             string        `json:"townSelected"`
	GeneralUserProfileEntityBranchSelected string        `json:"branchSelected"`
	GeneralUserProfileEmailStr             string        `json:"emailStr"`
	GeneralUserProfilePhoneNumberStr       string        `json:"phoneNumberStr"`
	GeneralUserProfileDepartment           DepartmentDTO `json:"department"`
	GeneralUserProfileCity                 CityDTO       `json:"city"`
	GeneralUserProfileSelectedCities       []CityDTO     `json:"selectedCities"`
	GeneralUserProfileSelectedTowns        []TownDTO     `json:"selectedTowns"`
	GeneralUserProfileSelectedBranches     []interface{} `json:"selectedBranches"` //Está ripo interface para que no haya un ciclo concurrente entre paquetes
	GeneralUserProfileEntity               interface{}   `json:"entity"`
}

// GeneralUserProfilePgDB representa la estructura que mapea los datos del perfil de usuario
// provenientes de la base de datos PostgreSQL. Utiliza tipos sql.Null* para manejar valores nulos.
type GeneralUserProfilePgDB struct {
	GeneralUserProfileId                   sql.NullInt64
	GeneralUserProfileICode                sql.NullString
	GeneralUserProfileCreationDate         sql.NullTime
	GeneralUserProfileUpdateDate           sql.NullTime
	GeneralUserProfileGender               sql.NullString
	GeneralUserProfileNick                 sql.NullString
	GeneralUserProfileDescription          sql.NullString
	GeneralUserProfileNames                sql.NullString
	GeneralUserProfileLastNames            sql.NullString
	GeneralUserProfileDocType              sql.NullString
	GeneralUserProfileDocNumber            sql.NullString
	GeneralUserProfileTown                 sql.NullString
	GeneralUserProfileEntityBranchSelected sql.NullString
}

func (up GeneralUserProfileDTO) MarshalJSON() ([]byte, error) {
	type Alias GeneralUserProfileDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		GeneralUserProfileCreationDate string `json:"creationDate"`
		GeneralUserProfileUpdateDate   string `json:"updateDate"`
	}{
		Alias:                          (*Alias)(&up),
		GeneralUserProfileCreationDate: up.GeneralUserProfileCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		GeneralUserProfileUpdateDate:   up.GeneralUserProfileUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (up *GeneralUserProfileDTO) UnmarshalJSON(data []byte) error {
	type Alias GeneralUserProfileDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		GeneralUserProfileCreationDate string `json:"creationDate"`
		GeneralUserProfileUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(up),
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
	up.GeneralUserProfileCreationDate = parse(aux.GeneralUserProfileCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	up.GeneralUserProfileUpdateDate = parse(aux.GeneralUserProfileUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetGeneralUserProfile inserta un nuevo perfil de usuario en la base de datos.
//
// Parámetros:
//   - profile: Puntero al DTO con la información del perfil a insertar.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo que invoca la función.
//   - connData: Datos de conexión a la BD.
//   - clientConfig: Configuración del cliente de BD.
//   - serverConfig: Configuración del servidor de BD.
//
// Retorna:
//   - La conexión actualizada y un error en caso de fallo.
func SetGeneralUserProfile(profile *GeneralUserProfileDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia para gestionar la conexión y la ejecución de la query.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de campos para la inserción en la tabla.
	var usrProfileFieldsSlice []string = []string{"GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber", "GeneralUserProfileTown"}
	var usrProfileFieldsAliasSlice []string = []string{}

	// Construcción dinámica de la query SQL para inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, usrProfileFieldsSlice, usrProfileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{"GeneralUserProfileId"}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, false)

	// Se establecen valores predeterminados para el perfil.
	profile.GeneralUserProfileICode = utils.GetUUID()
	profile.GeneralUserProfileCreationDate = time.Now()
	profile.GeneralUserProfileUpdateDate = time.Now()

	// Ejecución de la query de inserción y asignación del valor generado (ID).
	persistenceCtrl.QueryRow(context.Background(), query,
		profile.GeneralUserProfileICode,
		profile.GeneralUserProfileCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		profile.GeneralUserProfileUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		profile.GeneralUserProfileGender, profile.GeneralUserProfileNick, profile.GeneralUserProfileDescription,
		profile.GeneralUserProfileNames, profile.GeneralUserProfileLastNames, profile.GeneralUserProfileDocType,
		profile.GeneralUserProfileDocNumber, profile.GeneralUserProfileTownCode)

	persistenceCtrl.Scan(&profile.GeneralUserProfileId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateGeneralUserProfile actualiza un perfil de usuario existente en la base de datos.
//
// Parámetros:
//   - profile: Puntero al DTO que contiene la información actualizada.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo que invoca la función.
//   - connData: Datos de conexión a la BD.
//   - clientConfig: Configuración del cliente de BD.
//   - serverConfig: Configuración del servidor de BD.
//
// Retorna:
//   - La conexión actualizada y un error en caso de fallo.
func UpdateGeneralUserProfile(profile *GeneralUserProfileDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración de la conexión.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de campos a actualizar.
	var usrProfileFieldsSlice []string = []string{"GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber", "GeneralUserProfileTown"}
	var usrProfileFieldsAliasSlice []string = []string{}

	// Construcción de la query de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, usrProfileFieldsSlice, usrProfileFieldsAliasSlice, GeneralUserProfileDBName, []string{"GeneralUserProfileId"}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, false)

	// Ejecución de la actualización utilizando los valores del perfil.
	persistenceCtrl.Exec(context.Background(), query,
		profile.GeneralUserProfileId,
		profile.GeneralUserProfileUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		profile.GeneralUserProfileGender, profile.GeneralUserProfileNick, profile.GeneralUserProfileDescription,
		profile.GeneralUserProfileNames, profile.GeneralUserProfileLastNames, profile.GeneralUserProfileDocType,
		profile.GeneralUserProfileDocNumber, profile.GeneralUserProfileTownCode)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	// Verifica que se haya afectado al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveGeneralUserProfile elimina un perfil de usuario de la base de datos.
//
// Parámetros:
//   - profile: Puntero al DTO que identifica el perfil a eliminar.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo que invoca la función.
//   - connData: Datos de conexión a la BD.
//   - clientConfig: Configuración del cliente de BD.
//   - serverConfig: Configuración del servidor de BD.
//
// Retorna:
//   - La conexión actualizada y un error en caso de fallo.
func RemoveGeneralUserProfile(profile *GeneralUserProfileDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración de la conexión.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// No se definen campos para eliminación, se utiliza la clave primaria.
	var usrProfileFieldsSlice []string = []string{}
	var usrProfileFieldsAliasSlice []string = []string{}

	// Construcción de la query de eliminación.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, usrProfileFieldsSlice, usrProfileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{"GeneralUserProfileId"}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, false)

	// Ejecución de la query de eliminación.
	persistenceCtrl.Exec(context.Background(), query, profile.GeneralUserProfileId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica que se haya afectado al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetGeneralUserProfile obtiene la información de un perfil de usuario según los criterios indicados.
//
// Parámetros:
//   - by: Estructura que define los atributos y operadores de filtrado.
//   - profile: Puntero al DTO donde se almacenará la información obtenida.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo que invoca la función.
//   - connData: Datos de conexión a la BD.
//   - clientConfig: Configuración del cliente de BD.
//   - serverConfig: Configuración del servidor de BD.
//
// Retorna:
//   - La conexión actualizada y un error en caso de fallo.
func GetGeneralUserProfile(by common_controllers.By, profile *GeneralUserProfileDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se forma el path completo (esquema.tabla) para la consulta.
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName

	// Configuración de la conexión.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate",
		"GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames",
		"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Genera la parte de la query que selecciona solo los campos necesarios.
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Construcción completa de la query con cláusula WHERE.
	var query string = `SELECT ` + profileFieldsStr +
		` FROM ` + profilePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, GeneralUserProfileDBName, by.AttrsName, []string{}, []string{}, by.Operator, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Ejecución de la consulta y escaneo de los resultados en el DTO.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)
	persistenceCtrl.Scan(&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate,
		&profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames,
		&profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetGeneralUserProfileByEmail obtiene un perfil de usuario basado en el correo electrónico asociado.
//
// Parámetros:
//   - email: Dirección de correo electrónico a buscar.
//   - profile: Puntero al DTO donde se almacenará la información obtenida.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo que invoca la función.
//   - connData: Datos de conexión a la BD.
//   - clientConfig: Configuración del cliente de BD.
//   - serverConfig: Configuración del servidor de BD.
//
// Retorna:
//   - La conexión actualizada y un error en caso de fallo.
func GetGeneralUserProfileByEmail(email string, profile *GeneralUserProfileDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Definición de los paths completos para la tabla de perfil y la tabla de email.
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName
	var emailPath string = EMailDBScheme + "." + EMailDBName

	// Configuración de la conexión.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de campos a seleccionar para el perfil.
	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate",
		"GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames",
		"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}
	// Genera la parte de la query para la selección de campos.
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Construcción de la query con un RIGHT JOIN para relacionar el perfil con el correo electrónico.
	var query string = `SELECT ` + profileFieldsStr +
		` FROM ` + profilePath +
		` RIGHT JOIN ` + emailPath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + emailPath + `.` + EMailFieldDefinitions["EMailGeneralUserProfile"].DBName + `)` +
		` WHERE ` + emailPath + `.` + EMailFieldDefinitions["EMailData"].DBName + ` = $1`

	// Ejecución de la consulta y escaneo de los resultados en el DTO.
	persistenceCtrl.QueryRow(context.Background(), query, email)
	persistenceCtrl.Scan(&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate,
		&profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames,
		&profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// SetGeneralUserProfileDefaults establece valores predeterminados para un perfil de usuario,
// dependiendo de la acción (inserción o actualización).
//
// Parámetros:
//   - profile: Puntero al DTO del perfil.
//   - action: Acción a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetGeneralUserProfileDefaults(profile *GeneralUserProfileDTO, action string) {
	switch action {
	// En caso de inserción, se establecen la fecha de creación, de actualización y un nuevo UUID.
	case common_dao.SQL_INSERT:
		profile.GeneralUserProfileCreationDate = time.Now()
		profile.GeneralUserProfileUpdateDate = time.Now()
		profile.GeneralUserProfileICode = utils.GetUUID()
	// En caso de actualización, solo se actualiza la fecha de última modificación.
	case common_dao.SQL_UPDATE:
		profile.GeneralUserProfileUpdateDate = time.Now()
	}
}

// ToDTO convierte una estructura GeneralUserProfilePgDB (obtenida desde la BD) a un DTO.
//
// Retorna:
//   - Un objeto GeneralUserProfileDTO con los datos convertidos.
func (obj *GeneralUserProfilePgDB) ToDTO() GeneralUserProfileDTO {
	var dto GeneralUserProfileDTO

	if obj.GeneralUserProfileId.Valid {
		dto.GeneralUserProfileId = uint64(obj.GeneralUserProfileId.Int64)
	}

	if obj.GeneralUserProfileICode.Valid {
		dto.GeneralUserProfileICode = obj.GeneralUserProfileICode.String
	}

	if obj.GeneralUserProfileCreationDate.Valid {
		dto.GeneralUserProfileCreationDate = obj.GeneralUserProfileCreationDate.Time
	}

	if obj.GeneralUserProfileUpdateDate.Valid {
		dto.GeneralUserProfileUpdateDate = obj.GeneralUserProfileUpdateDate.Time
	}

	if obj.GeneralUserProfileGender.Valid {
		dto.GeneralUserProfileGender = obj.GeneralUserProfileGender.String
	}

	if obj.GeneralUserProfileNick.Valid {
		dto.GeneralUserProfileNick = obj.GeneralUserProfileNick.String
	}

	if obj.GeneralUserProfileDescription.Valid {
		dto.GeneralUserProfileDescription = obj.GeneralUserProfileDescription.String
	}

	if obj.GeneralUserProfileNames.Valid {
		dto.GeneralUserProfileNames = obj.GeneralUserProfileNames.String
	}

	if obj.GeneralUserProfileLastNames.Valid {
		dto.GeneralUserProfileLastNames = obj.GeneralUserProfileLastNames.String
	}

	if obj.GeneralUserProfileDocType.Valid {
		dto.GeneralUserProfileDocType = obj.GeneralUserProfileDocType.String
	}

	if obj.GeneralUserProfileDocNumber.Valid {
		dto.GeneralUserProfileDocNumber = obj.GeneralUserProfileDocNumber.String
	}

	if obj.GeneralUserProfileEntityBranchSelected.Valid {
		dto.GeneralUserProfileEntityBranchSelected = obj.GeneralUserProfileEntityBranchSelected.String
	}

	if obj.GeneralUserProfileTown.Valid {
		dto.GeneralUserProfileTown = TownDTO{TownCode: obj.GeneralUserProfileTown.String}
	}

	return dto
}
