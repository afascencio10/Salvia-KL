// Package security_daos contiene las funciones DAO (Data Access Object) para la gestión de emails en el sistema de seguridad.
// Este paquete permite insertar, actualizar, eliminar y consultar registros de emails en la base de datos.
package security_daos

import (
	// Configuración común del sistema.
	// Controladores comunes para la persistencia.
	// Funciones DAO comunes.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"    // Módulo para la gestión de conexiones a la base de datos.
	"bitsflow/common/utils" // Utilidades varias, como generación de UUID y definiciones de campos.
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	// EMailntityName es el nombre lógico de la entidad Email en el sistema.
	EMailntityName string = "EMail"
	// EMailJSONName es el nombre que se utiliza para la representación JSON de la entidad Email.
	EMailJSONName string = "mail"
	// EMailDBName es el nombre de la tabla en la base de datos donde se almacenan los emails.
	EMailDBName string = "email"
	// EMailDBScheme es el esquema de la base de datos donde se encuentra la tabla de emails.
	EMailDBScheme string = "security"

	// EMailFieldDefinitions define los atributos y validaciones de cada campo del email.
	// Cada entrada mapea el nombre del campo a su definición, que incluye nombres en la base de datos,
	// el tipo en el modelo, restricciones de tamaño y obligatoriedad.
	EMailFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"EMailId":                 {Name: "EMailId", DBName: "e_mail_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"EMailICode":              {Name: "EMailICode", DBName: "e_mail_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"EMailCreationDate":       {Name: "EMailCreationDate", DBName: "e_mail_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"EMailUpdateDate":         {Name: "EMailUpdateDate", DBName: "e_mail_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"EMailData":               {Name: "EMailData", DBName: "e_mail_data", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 128, Required: true},
		"EMailGeneralUserProfile": {Name: "EMailGeneralUserProfile", DBName: "e_mail_general_user_profile", Alias: "", ModelType: "uint", Required: false},
	}
)

// EMailDTO es el objeto de transferencia de datos (DTO) para la entidad Email.
// Representa la estructura de datos que se utiliza para la comunicación entre capas de la aplicación.
type EMailDTO struct {
	EMailId                 uint64                `json:"-"`
	EMailICode              string                `json:"icode"`
	EMailCreationDate       time.Time             `json:"creation_date"`
	EMailUpdateDate         time.Time             `json:"update_date"`
	EMailData               string                `json:"data"`
	EMailGeneralUserProfile GeneralUserProfileDTO `json:"-"`
}

// EMailPgDB representa la estructura que mapea los campos de la tabla email en PostgreSQL.
// Utiliza tipos nulos de SQL para gestionar valores que pueden ser NULL en la base de datos.
type EMailPgDB struct {
	EMailId                 sql.NullInt64
	EMailICode              sql.NullString
	EMailCreationDate       sql.NullTime
	EMailUpdateDate         sql.NullTime
	EMailData               sql.NullString
	EMailGeneralUserProfile sql.NullInt64
}

// SetEMail inserta un nuevo registro de email en la base de datos.
// Recibe un puntero a EMailDTO con los datos del email, y parámetros de conexión y configuración.
// Retorna la conexión actualizada y un error en caso de que ocurra.
func SetEMail(email *EMailDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Si el perfil general de usuario está vacío (ID == 0), se asigna nil.
	var profileId interface{} = email.EMailGeneralUserProfile.GeneralUserProfileId
	if email.EMailGeneralUserProfile.GeneralUserProfileId == 0 {
		profileId = nil
	}

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Se definen los campos que se insertarán en la tabla.
	var emailFieldsSlice []string = []string{"EMailICode", "EMailCreationDate", "EMailUpdateDate", "EMailData", "EMailGeneralUserProfile"}
	var emailFieldsAliasSlice []string = []string{}

	// Se genera la consulta SQL para la inserción del email.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{}, []string{}, []string{"EMailId"}, common_dao.SQL_AND, EMailDBScheme, EMailFieldDefinitions, false)

	// Se ejecuta la consulta y se obtiene el ID generado.
	persistenceCtrl.QueryRow(context.Background(), query,
		email.EMailICode,
		email.EMailCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		email.EMailUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		email.EMailData,
		profileId)

	// Se escanea el resultado para obtener el ID del email.
	persistenceCtrl.Scan(&email.EMailId)

	// En caso de error, se imprime la consulta y el error, y se retorna.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateEMailByICode actualiza un registro de email identificado por su código único (ICode).
// Recibe el email a actualizar y parámetros de conexión.
// Retorna la conexión actualizada y un error en caso de fallo.
func UpdateEMailByICode(email *EMailDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se actualizarán.
	var emailFieldsSlice []string = []string{"EMailUpdateDate", "EMailData"}
	var emailFieldsAliasSlice []string = []string{}

	// Se genera la consulta SQL para actualizar el registro.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{}, []string{"EMailICode"}, []string{}, common_dao.SQL_AND, EMailDBScheme, EMailFieldDefinitions, false)

	// Se ejecuta la consulta con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		email.EMailICode,
		email.EMailUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		email.EMailData)

	// Se verifica si ocurrió algún error o si no se afectó ninguna fila.
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

// RemoveEMailByICode elimina un registro de email identificado por su código único (ICode).
// Recibe el email a eliminar y parámetros de conexión.
// Retorna la conexión actualizada y un error en caso de fallo.
func RemoveEMailByICode(email *EMailDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se generan los slices vacíos ya que no se requieren campos adicionales para la eliminación.
	var emailFieldsSlice []string = []string{}
	var emailFieldsAliasSlice []string = []string{}

	// Se genera la consulta SQL para eliminar el registro basado en el ICode.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{"EMailICode"}, []string{}, []string{}, common_dao.SQL_AND, EMailDBScheme, EMailFieldDefinitions, false)

	// Se ejecuta la consulta.
	persistenceCtrl.Exec(context.Background(), query, email.EMailICode)
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

// RemoveEMails elimina registros de email basándose en atributos dinámicos definidos en el parámetro 'by'.
// Permite eliminar múltiples registros que cumplan con las condiciones indicadas.
// Retorna la conexión actualizada y un error en caso de fallo.
func RemoveEMails(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se genera la consulta SQL utilizando los atributos y operadores proporcionados.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, EMailDBName, by.AttrsName, []string{}, []string{}, by.Operator, EMailDBScheme, EMailFieldDefinitions, false)

	// Se ejecuta la consulta con los valores correspondientes.
	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetEMail obtiene un registro de email junto con el perfil de usuario asociado.
// La consulta se realiza usando condiciones definidas en el parámetro 'by'.
// Retorna la conexión actualizada y un error en caso de fallo.
func GetEMail(by common_controllers.By, email *EMailDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construyen las rutas de la tabla email y la tabla de perfiles.
	var emailPath string = EMailDBScheme + "." + EMailDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName

	// Se inicializa el perfil del email.
	email.EMailGeneralUserProfile = GeneralUserProfileDTO{}

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a seleccionar tanto del email como del perfil.
	var emailFieldsSlice []string = []string{"EMailId", "EMailICode", "EMailCreationDate", "EMailUpdateDate", "EMailData", "EMailGeneralUserProfile"}
	var emailFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Se generan los strings de campos para la consulta SQL.
	var emailFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EMailDBScheme, EMailFieldDefinitions, true)
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Se construye la consulta SQL con LEFT JOIN para incluir datos del perfil.
	var query string = `SELECT ` + emailFieldsStr + `, ` + profileFieldsStr +
		` FROM ` + emailPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + emailPath + `.` + EMailFieldDefinitions["EMailGeneralUserProfile"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, EMailDBName, by.AttrsName, []string{}, []string{}, by.Operator, EMailDBScheme, EMailFieldDefinitions, true)

	// Se imprime la consulta para depuración (puede ser removido en producción).
	fmt.Printf(query, by.AttrsValue...)
	// Se ejecuta la consulta para obtener el registro.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se crea una variable para mapear los datos del perfil desde la base de datos.
	var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
	// Se escanean los resultados en las variables correspondientes.
	persistenceCtrl.Scan(&email.EMailId, &email.EMailICode, &email.EMailCreationDate, &email.EMailUpdateDate, &email.EMailData,
		&profile.GeneralUserProfileId,
		&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

	// Se convierte el perfil obtenido al DTO correspondiente.
	email.EMailGeneralUserProfile = profile.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetEMails obtiene una lista de emails que cumplen con las condiciones definidas en 'by'.
// Solo se consultan los datos del email sin incluir información de perfil.
// Retorna la conexión actualizada, la lista de emails y un error en caso de fallo.
func GetEMails(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EMailDTO, error) {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye la ruta de la tabla email.
	var emailPath string = EMailDBScheme + "." + EMailDBName

	var email = EMailDTO{}
	var emails []EMailDTO

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a seleccionar.
	var emailFieldsSlice []string = []string{"EMailId", "EMailICode", "EMailCreationDate", "EMailUpdateDate", "EMailData", "EMailGeneralUserProfile"}
	var emailFieldsAliasSlice []string = []string{}

	// Se genera el string de campos para la consulta SQL.
	var emailFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EMailDBScheme, EMailFieldDefinitions, true)

	// Se construye la consulta SQL para seleccionar los emails.
	var query string = `SELECT ` + emailFieldsStr +
		` FROM ` + emailPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, EMailDBName, by.AttrsName, []string{}, []string{}, by.Operator, EMailDBScheme, EMailFieldDefinitions, true)

	// Se imprime la consulta para depuración.
	fmt.Printf(query, by.AttrsValue...)
	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Se itera sobre los resultados y se construye la lista de emails.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&email.EMailId, &email.EMailICode, &email.EMailCreationDate, &email.EMailUpdateDate, &email.EMailData, &email.EMailGeneralUserProfile.GeneralUserProfileId)
		emails = append(emails, email)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return emails, nil
}

// GetAllEMail obtiene todos los registros de email existentes en la base de datos.
// No aplica filtros, por lo que retorna todos los registros.
// Retorna la conexión actualizada, la lista de emails y un error en caso de fallo.
func GetAllEMail(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EMailDTO, error) {
	// Se define una variable temporal para almacenar cada email.
	var email EMailDTO
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a seleccionar.
	var emailFieldsSlice []string = []string{"EMailId", "EMailICode", "EMailCreationDate", "EMailUpdateDate", "EMailData", "EMailGeneralUserProfile"}
	var emailFieldsAliasSlice []string = []string{}

	// Se genera la consulta SQL para seleccionar todos los emails.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{}, []string{}, []string{}, "", EMailDBScheme, EMailFieldDefinitions, true)

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var emails []EMailDTO
	// Se itera sobre los registros obtenidos.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		email = EMailDTO{}
		persistenceCtrl.ScanRow(&email.EMailId, &email.EMailICode, &email.EMailCreationDate, &email.EMailUpdateDate, &email.EMailData,
			&profile.GeneralUserProfileId)

		// Se convierte el perfil obtenido a DTO.
		email.EMailGeneralUserProfile = profile.ToDTO()
		emails = append(emails, email)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return emails, nil
}

// GetAllEMailWithProfile obtiene todos los registros de email junto con la información de perfil asociada.
// Utiliza LEFT JOIN para combinar la información de ambas tablas.
// Retorna la conexión actualizada, la lista de emails con perfiles y un error en caso de fallo.
func GetAllEMailWithProfile(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EMailDTO, error) {
	// Se crea una instancia del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se definen las rutas para las tablas email y perfil.
	var emailPath string = EMailDBScheme + "." + EMailDBName
	var profilePath string = GeneralUserProfileDBScheme + "." + GeneralUserProfileDBName

	// Se establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a seleccionar para el email.
	var emailFieldsSlice []string = []string{"EMailId", "EMailICode", "EMailCreationDate", "EMailUpdateDate", "EMailData", "EMailGeneralUserProfile"}
	var emailFieldsAliasSlice []string = []string{}

	// Se definen los campos a seleccionar para el perfil.
	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Se generan los strings de campos para ambas tablas.
	var emailFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, emailFieldsSlice, emailFieldsAliasSlice, EMailDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EMailDBScheme, EMailFieldDefinitions, true)
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, GeneralUserProfileDBScheme, GeneralUserProfileFieldDefinitions, true)

	// Se construye la consulta SQL con LEFT JOIN para obtener la información del email y del perfil.
	var query string = `SELECT ` + emailFieldsStr + ", " + profileFieldsStr +
		` FROM ` + emailPath +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + EMailFieldDefinitions["EMailGeneralUserProfile"].DBName + `)`

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var emails []EMailDTO
	// Se itera sobre los resultados y se construyen los DTOs correspondientes.
	for persistenceCtrl.Next() {
		var profile GeneralUserProfilePgDB = GeneralUserProfilePgDB{}
		var email = EMailDTO{}

		persistenceCtrl.ScanRow(&email.EMailId, &email.EMailICode, &email.EMailCreationDate, &email.EMailUpdateDate, &email.EMailData,
			&profile.GeneralUserProfileId,
			&profile.GeneralUserProfileId, &profile.GeneralUserProfileICode, &profile.GeneralUserProfileCreationDate, &profile.GeneralUserProfileUpdateDate, &profile.GeneralUserProfileGender, &profile.GeneralUserProfileNick, &profile.GeneralUserProfileDescription, &profile.GeneralUserProfileNames, &profile.GeneralUserProfileLastNames, &profile.GeneralUserProfileDocType, &profile.GeneralUserProfileDocNumber)

		// Se asigna el perfil convertido al DTO del email.
		email.EMailGeneralUserProfile = profile.ToDTO()

		emails = append(emails, email)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return emails, nil
}

// SetEMailDefaults establece los valores por defecto para un objeto EMailDTO dependiendo de la acción que se realice.
// Para inserción, se inicializan las fechas de creación y actualización y se genera un código único (UUID).
// Para actualización, solo se actualiza la fecha de modificación.
func SetEMailDefaults(email *EMailDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		email.EMailCreationDate = time.Now()
		email.EMailUpdateDate = time.Now()
		email.EMailICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		email.EMailUpdateDate = time.Now()
	}
}

// PgDBToDTO convierte una estructura EMailPgDB, que representa los datos de la base de datos, en un objeto EMailDTO.
// Realiza la conversión de tipos nulos a sus correspondientes valores en el DTO.
func (obj *EMailPgDB) PgDBToDTO() EMailDTO {
	var dto EMailDTO

	if obj.EMailId.Valid {
		dto.EMailId = uint64(obj.EMailId.Int64)
	}

	if obj.EMailICode.Valid {
		dto.EMailICode = obj.EMailICode.String
	}

	if obj.EMailCreationDate.Valid {
		dto.EMailCreationDate = obj.EMailCreationDate.Time
	}

	if obj.EMailUpdateDate.Valid {
		dto.EMailUpdateDate = obj.EMailUpdateDate.Time
	}

	if obj.EMailData.Valid {
		dto.EMailData = obj.EMailData.String
	}

	if obj.EMailGeneralUserProfile.Valid {
		dto.EMailGeneralUserProfile = GeneralUserProfileDTO{GeneralUserProfileId: uint64(obj.EMailGeneralUserProfile.Int64)}
	}

	return dto
}
