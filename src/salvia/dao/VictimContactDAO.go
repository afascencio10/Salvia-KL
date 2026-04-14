// Package salvia_daos implementa los Data Access Objects (DAO) relacionados con el contacto de la víctima.
// Contiene funciones para insertar, obtener, actualizar y eliminar registros de contacto de la víctima en la base de datos.
package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	// Nombres y esquemas de la entidad VictimContact
	VictimContactEntityName string = "VictimContact"  // Nombre de la entidad
	VictimContactJSONName   string = "victimContact"  // Nombre utilizado en formato JSON
	VictimContactDBName     string = "victim_contact" // Nombre de la tabla en la base de datos
	VictimContactDBScheme   string = "salvia"         // Esquema de la base de datos

	// Atributos relacionados con las validaciones ------------------------------

	// Definiciones de campos para VictimContact. Cada entrada mapea el nombre del campo a su definición,
	// incluyendo detalles como el nombre en la base de datos, tipo de dato, tamaño mínimo y máximo, y si es requerido.
	VictimContactFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimContactId":                {Name: "VictimContactId", DBName: "victim_contact_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactICode":             {Name: "VictimContactICode", DBName: "victim_contact_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimContactCreationDate":      {Name: "VictimContactCreationDate", DBName: "victim_contact_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactUpdateDate":        {Name: "VictimContactUpdateDate", DBName: "victim_contact_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactGeneralUser":       {Name: "VictimContactGeneralUser", DBName: "victim_contact_general_user", Alias: "", ModelType: "string", MinSize: 16, MaxSize: 36, Required: false},
		"VictimContactStatus":            {Name: "VictimContactStatus", DBName: "victim_contact_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimContactStatusDescription": {Name: "VictimContactStatusDescription", DBName: "victim_contact_status_description", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 256, Required: false},
		"VictimContactLatitude":          {Name: "VictimContactLatitude", DBName: "victim_contact_latitude", Alias: "", ModelType: "float", MinSize: -90, MaxSize: 90, Required: false},
		"VictimContactLongitude":         {Name: "VictimContactLongitude", DBName: "victim_contact_longitude", Alias: "", ModelType: "float", MinSize: -180, MaxSize: 180, Required: false},
		"VictimContactNames":             {Name: "VictimContactNames", DBName: "victim_contact_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimContactLastNames":         {Name: "VictimContactLastNames", DBName: "victim_contact_last_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
	}
)

// VictimContactDTO representa el Data Transfer Object para el contacto de la víctima.
// Este DTO se utiliza para transferir datos entre capas de la aplicación y para serialización JSON.
type VictimContactDTO struct {
	VictimContactId                uint64    `json:"-"`
	VictimContactICode             string    `json:"icode"`
	VictimContactCreationDate      time.Time `json:"creationDate"`
	VictimContactUpdateDate        time.Time `json:"updateDate"`
	VictimContactStatus            string    `json:"status"`
	VictimContactStatusDescription string    `json:"statusDescription"`
	VictimContactGeneralUser       string    `json:"generalUser"`
	VictimContactNames             string    `json:"names"`
	VictimContactLastNames         string    `json:"lastNames"`
	VictimContactLatitude          float64   `json:"latitude"`
	VictimContactLongitude         float64   `json:"longitude"`

	// Parámetros auxiliares
	VictimContactCaptchaID       string `json:"captchaID"`
	VictimContactCaptchaSolution string `json:"captchaSolution"`

	VictimContactForm1 VictimContactForm1DTO `json:"form"`
	VictimContactForm2 VictimContactForm2DTO `json:"form2"`
}

// VictimContactPgDB representa el mapeo de la tabla de VictimContact en la base de datos PostgreSQL.
// Utiliza sql.Null* para permitir valores nulos en la base de datos.
type VictimContactPgDB struct {
	VictimContactId                sql.NullInt64
	VictimContactICode             sql.NullString
	VictimContactCreationDate      sql.NullTime
	VictimContactUpdateDate        sql.NullTime
	VictimContactStatus            sql.NullString
	VictimContactStatusDescription sql.NullString
	VictimContactGeneralUser       sql.NullString
	VictimContactNames             sql.NullString
	VictimContactLastNames         sql.NullString
	VictimContactLatitude          sql.NullFloat64
	VictimContactLongitude         sql.NullFloat64
}

// MarshalJSON implementa la interfaz json.Marshaler para formatear las fechas con un formato específico.
// Se crea un alias de VictimContactDTO para evitar recursión infinita durante la serialización.
func (vcd VictimContactDTO) MarshalJSON() ([]byte, error) {
	type Alias VictimContactDTO // Alias para evitar recursión infinita

	// Se formatean las fechas de nacimiento, creación y actualización utilizando el formato definido en common_config.
	return json.Marshal(&struct {
		*Alias
		VictimContactCreationDate string `json:"creationDate"`
		VictimContactUpdateDate   string `json:"updateDate"`
	}{
		Alias:                     (*Alias)(&vcd),
		VictimContactCreationDate: vcd.VictimContactCreationDate.Format(common_config.DateTime.DATE_FORMAT),
		VictimContactUpdateDate:   vcd.VictimContactUpdateDate.Format(common_config.DateTime.DATE_FORMAT),
	})
}

func (vcd *VictimContactDTO) UnmarshalJSON(data []byte) error {
	type Alias VictimContactDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		VictimContactCreationDate string `json:"creationDate"`
		VictimContactUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(vcd),
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

	vcd.VictimContactCreationDate = parse(aux.VictimContactCreationDate, common_config.DateTime.DATE_FORMAT)

	vcd.VictimContactUpdateDate = parse(aux.VictimContactUpdateDate, common_config.DateTime.DATE_FORMAT)

	return nil
}

// SetVictimContact inserta un nuevo registro de contacto de víctima en la base de datos.
// Parámetros:
//   - victimContact: Puntero al objeto VictimContactDTO que contiene los datos a insertar.

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un posible error durante la operación.
func SetVictimContact(victimContact *VictimContactDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se define el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Definición de campos que se insertarán en la base de datos.
	var victimContactFieldsSlice []string = []string{"VictimContactICode", "VictimContactCreationDate", "VictimContactUpdateDate", "VictimContactGeneralUser", "VictimContactStatus", "VictimContactNames", "VictimContactLastNames", "VictimContactLatitude", "VictimContactLongitude"}
	var victimContactFieldsAliasSlice []string = []string{}

	// Se construye la consulta SQL de inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, victimContactFieldsSlice, victimContactFieldsAliasSlice, VictimContactDBName, []string{}, []string{}, []string{"VictimContactId"}, common_dao.SQL_AND, VictimContactDBScheme, VictimContactFieldDefinitions, false)

	// Se ejecuta la consulta con los valores correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query,
		victimContact.VictimContactICode, victimContact.VictimContactCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), victimContact.VictimContactUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimContact.VictimContactGeneralUser, victimContact.VictimContactStatus, victimContact.VictimContactNames, victimContact.VictimContactLastNames, victimContact.VictimContactLatitude, victimContact.VictimContactLongitude)

	// Se escanea el resultado para obtener el ID generado.
	persistenceCtrl.Scan(&victimContact.VictimContactId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetVictimContact obtiene un registro de contacto de víctima de la base de datos utilizando criterios de búsqueda.
// Parámetros:
//   - by: Criterio de búsqueda (estructura con nombres y valores de atributos).
//   - victimContact: Puntero al objeto VictimContactDTO donde se almacenarán los datos obtenidos.

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un posible error durante la operación.
func GetVictimContact(by common_controllers.By, victimContact *VictimContactDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se define el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactPath string = VictimContactDBScheme + "." + VictimContactDBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var victimContactFieldsSlice []string = []string{"VictimContactId", "VictimContactICode", "VictimContactCreationDate", "VictimContactUpdateDate",
		"VictimContactGeneralUser", "VictimContactStatus", "VictimContactNames", "VictimContactLastNames", "VictimContactLatitude", "VictimContactLongitude", "VictimContactStatusDescription"}
	var victimContactFieldsAliasSlice []string = []string{}

	// Se construye la parte de la consulta que define los campos a seleccionar.
	var victimContactFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactFieldsSlice, victimContactFieldsAliasSlice, VictimContactDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactDBScheme, VictimContactFieldDefinitions, true)

	// Se construye la consulta SQL completa, incluyendo la cláusula WHERE basada en el criterio de búsqueda.
	var query string = `SELECT ` + victimContactFieldsStr +
		` FROM ` + victimContactPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimContactDBScheme, VictimContactFieldDefinitions, true)

	// Se imprime la consulta para fines de depuración.
	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se define una variable temporal para almacenar los resultados de la consulta.
	var victimContactPG VictimContactPgDB = VictimContactPgDB{}

	// Se escanean los resultados en la estructura correspondiente.
	persistenceCtrl.Scan(&victimContactPG.VictimContactId, &victimContactPG.VictimContactICode, &victimContactPG.VictimContactCreationDate, &victimContactPG.VictimContactUpdateDate,
		&victimContactPG.VictimContactGeneralUser, &victimContactPG.VictimContactStatus, &victimContactPG.VictimContactNames, &victimContactPG.VictimContactLastNames, &victimContactPG.VictimContactLatitude, &victimContactPG.VictimContactLongitude, &victimContactPG.VictimContactStatusDescription)

	// Se convierte la estructura obtenida a DTO y se asigna al puntero recibido.
	var tmp VictimContactDTO = victimContactPG.ToDTO()
	*victimContact = tmp
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetVictimContacts obtiene múltiples registros de contacto de víctima basándose en un criterio de búsqueda.
// Parámetros:
//   - by: Criterio de búsqueda (estructura con nombres y valores de atributos).

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactDTO y un posible error durante la operación.
func GetVictimContacts(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactDTO, error) {
	// Se define el controlador de persistencia y las variables necesarias.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactPath string = VictimContactDBScheme + "." + VictimContactDBName

	var victimContacts []VictimContactDTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var victimContactFieldsSlice []string = []string{"VictimContactICode", "VictimContactCreationDate", "VictimContactStatus", "VictimContactNames", "VictimContactLastNames", "VictimContactLatitude", "VictimContactLongitude"}
	var victimContactFieldsAliasSlice []string = []string{}

	// Se construye la parte de la consulta que define los campos a seleccionar.
	var victimContactFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactFieldsSlice, victimContactFieldsAliasSlice, VictimContactDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactDBScheme, VictimContactFieldDefinitions, true)

	// Se construye la consulta SQL completa, incluyendo cláusula WHERE y ordenamiento.
	var query string = `SELECT ` + victimContactFieldsStr +
		` FROM ` + victimContactPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimContactDBScheme, VictimContactFieldDefinitions, true) +
		` ORDER BY ` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactCreationDate"].DBName + ` ASC  `

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Se itera sobre los resultados y se almacenan en el slice.
	for persistenceCtrl.Next() {
		var victimContactPG VictimContactPgDB = VictimContactPgDB{}
		persistenceCtrl.ScanRow(&victimContactPG.VictimContactICode, &victimContactPG.VictimContactCreationDate, &victimContactPG.VictimContactStatus, &victimContactPG.VictimContactNames, &victimContactPG.VictimContactLastNames, &victimContactPG.VictimContactLatitude, &victimContactPG.VictimContactLongitude)
		victimContacts = append(victimContacts, victimContactPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return victimContacts, nil
}

// GetVictimContactsWithoutVictimCase obtiene los registros de contacto de víctima que no tienen asociados un caso de víctima.
// Parámetros:
//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactDTO y un posible error durante la operación.
func GetVictimContactsWithoutVictimCase(page int, status string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactDTO, int, error) {
	// Se define el controlador de persistencia y las variables necesarias.
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactPath string = VictimContactDBScheme + "." + VictimContactDBName
	// Se asume que VictimCaseDBScheme y VictimCaseDBName y VictimCaseFieldDefinitions están definidos en otro lugar del proyecto.
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName

	var victimContacts []VictimContactDTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var victimContactFieldsSlice []string = []string{"VictimContactICode", "VictimContactCreationDate", "VictimContactStatus", "VictimContactNames", "VictimContactLastNames", "VictimContactLatitude", "VictimContactLongitude"}
	var victimContactFieldsAliasSlice []string = []string{}

	// Se construye la parte de la consulta que define los campos a seleccionar.
	var victimContactFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactFieldsSlice, victimContactFieldsAliasSlice, VictimContactDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactDBScheme, VictimContactFieldDefinitions, true)

	// Se construye la consulta SQL con LEFT JOIN para identificar registros sin caso asociado.
	var query string = `SELECT ` + victimContactFieldsStr +
		` FROM ` + victimContactPath +
		` LEFT JOIN ` + victimCasePath + ` ON (` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactId"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContact"].DBName + `)` +
		` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContact"].DBName + ` IS NULL  AND ` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactStatus"].DBName + ` = $1` +
		` ORDER BY ` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, status)

	// Se itera sobre los resultados y se almacena cada registro en el slice.
	for persistenceCtrl.Next() {
		var victimContactPG VictimContactPgDB = VictimContactPgDB{}
		persistenceCtrl.ScanRow(&victimContactPG.VictimContactICode, &victimContactPG.VictimContactCreationDate, &victimContactPG.VictimContactStatus, &victimContactPG.VictimContactNames, &victimContactPG.VictimContactLastNames, &victimContactPG.VictimContactLatitude, &victimContactPG.VictimContactLongitude)
		victimContacts = append(victimContacts, victimContactPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimContactPath +
			` LEFT JOIN ` + victimCasePath + ` ON (` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactId"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContact"].DBName + `)` +
			` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContact"].DBName + ` IS NULL AND ` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactStatus"].DBName + ` = $1`

		persistenceCtrl.QueryRow(context.Background(), countQuery, status)
		persistenceCtrl.Scan(&count)
	}

	return victimContacts, count, nil
}

// GetAllVictimContact obtiene todos los registros de contacto de víctima de la base de datos.
// Parámetros:
//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactDTO y un posible error durante la operación.
func GetAllVictimContact(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactDTO, error) {
	// Se define el controlador de persistencia y las variables necesarias.

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactPath string = VictimContactDBScheme + "." + VictimContactDBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var victimContactFieldsSlice []string = []string{"VictimContactICode", "VictimContactCreationDate", "VictimContactUpdateDate", "VictimContactStatus", "VictimContactNames", "VictimContactLastNames", "VictimContactLatitude", "VictimContactLongitude"}
	var victimContactFieldsAliasSlice []string = []string{}

	// Se construye la consulta SQL completa.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, victimContactFieldsSlice, victimContactFieldsAliasSlice, VictimContactDBName, []string{}, []string{}, []string{}, "", VictimContactDBScheme, VictimContactFieldDefinitions, true) +
		` ORDER BY ` + victimContactPath + `.` + VictimContactFieldDefinitions["VictimContactCreationDate"].DBName + ` ASC  `

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var victimContacts []VictimContactDTO
	// Se itera sobre los resultados y se almacena cada registro en el slice.
	for persistenceCtrl.Next() {
		var victimContactPG VictimContactPgDB = VictimContactPgDB{}
		persistenceCtrl.ScanRow(&victimContactPG.VictimContactICode, &victimContactPG.VictimContactCreationDate, &victimContactPG.VictimContactUpdateDate, &victimContactPG.VictimContactStatus, &victimContactPG.VictimContactNames, &victimContactPG.VictimContactLastNames, &victimContactPG.VictimContactLatitude, &victimContactPG.VictimContactLongitude)
		victimContacts = append(victimContacts, victimContactPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimContacts, nil
}

// InvalidateVictimContactByICode actualiza el registro de contacto de víctima para invalidarlo (por ejemplo, cambiar su estado).
// Parámetros:
//   - victimContact: Puntero al objeto VictimContactDTO que se invalidará.
//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un posible error durante la operación.
func InvalidateVictimContactByICode(victimContact *VictimContactDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se define el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de campos a actualizar.
	var vcontactFieldsSlice []string = []string{"VictimContactUpdateDate", "VictimContactStatus", "VictimContactStatusDescription"}
	var vcontactFieldsAliasSlice []string = []string{}

	// Se construye la consulta SQL de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, vcontactFieldsSlice, vcontactFieldsAliasSlice, VictimContactDBName, []string{"VictimContactICode"}, []string{}, []string{}, common_dao.SQL_AND, VictimContactDBScheme, VictimContactFieldDefinitions, false)

	// Se ejecuta la consulta con los valores correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		victimContact.VictimContactICode, victimContact.VictimContactUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimContact.VictimContactStatus, victimContact.VictimContactStatusDescription)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se verifica que la actualización haya afectado alguna fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// SetVictimContactDefaults asigna valores predeterminados a un objeto VictimContactDTO según la acción a realizar.
// Parámetros:
//   - victimContact: Puntero al objeto VictimContactDTO al que se asignarán los valores predeterminados.
//   - action: Acción que se va a realizar (por ejemplo, common_dao.SQL_INSERT o common_dao.SQL_UPDATE).
func SetVictimContactDefaults(victimContact *VictimContactDTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se asignan la fecha de creación y actualización al momento actual,
		// se genera un código único y se asigna un estado predeterminado.
		victimContact.VictimContactCreationDate = time.Now()
		victimContact.VictimContactUpdateDate = time.Now()
		victimContact.VictimContactICode = utils.GetUUID()
		victimContact.VictimContactStatus = "v"
	case common_dao.SQL_UPDATE:
		// Para una actualización, solo se asigna la fecha de actualización.
		victimContact.VictimContactUpdateDate = time.Now()
	}

}

// ToDTO convierte una instancia de VictimContactPgDB a su correspondiente VictimContactDTO.
// Este método verifica la validez de cada campo (usando sql.Null*) y asigna el valor correspondiente.
func (obj *VictimContactPgDB) ToDTO() VictimContactDTO {
	var dto VictimContactDTO

	if obj.VictimContactId.Valid {
		dto.VictimContactId = uint64(obj.VictimContactId.Int64)
	}

	if obj.VictimContactICode.Valid {
		dto.VictimContactICode = obj.VictimContactICode.String
	}

	if obj.VictimContactCreationDate.Valid {
		dto.VictimContactCreationDate = obj.VictimContactCreationDate.Time
	}

	if obj.VictimContactUpdateDate.Valid {
		dto.VictimContactUpdateDate = obj.VictimContactUpdateDate.Time
	}

	if obj.VictimContactStatus.Valid {
		dto.VictimContactStatus = obj.VictimContactStatus.String
	}

	if obj.VictimContactStatusDescription.Valid {
		dto.VictimContactStatusDescription = obj.VictimContactStatusDescription.String
	}

	if obj.VictimContactGeneralUser.Valid {
		dto.VictimContactGeneralUser = obj.VictimContactGeneralUser.String
	}

	if obj.VictimContactNames.Valid {
		dto.VictimContactNames = obj.VictimContactNames.String
	}

	if obj.VictimContactLastNames.Valid {
		dto.VictimContactLastNames = obj.VictimContactLastNames.String
	}

	if obj.VictimContactLatitude.Valid {
		dto.VictimContactLatitude = obj.VictimContactLatitude.Float64
	}
	if obj.VictimContactLongitude.Valid {
		dto.VictimContactLongitude = obj.VictimContactLongitude.Float64
	}
	return dto
}
