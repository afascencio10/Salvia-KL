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
	"fmt"
	"time"
)

var (
	// Nombres y esquemas de la entidad VictimContactForm1
	VictimContactForm1EntityName string = "VictimContactForm1"   // Nombre de la entidad
	VictimContactForm1JSONName   string = "victimContactForm1"   // Nombre utilizado en formato JSON
	VictimContactForm1DBName     string = "victim_contact_form1" // Nombre de la tabla en la base de datos
	VictimContactForm1DBScheme   string = "salvia"               // Esquema de la base de datos

	// Atributos relacionados con las validaciones ------------------------------

	// Definiciones de campos para VictimContactForm1. Cada entrada mapea el nombre del campo a su definición,
	// incluyendo detalles como el nombre en la base de datos, tipo de dato, tamaño mínimo y máximo, y si es requerido.
	VictimContactForm1FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimContactForm1Id":                {Name: "VictimContactForm1Id", DBName: "victim_contact_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm1ICode":             {Name: "VictimContactForm1ICode", DBName: "victim_contact_form1_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimContactForm1CreationDate":      {Name: "VictimContactForm1CreationDate", DBName: "victim_contact_form1_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm1UpdateDate":        {Name: "VictimContactForm1UpdateDate", DBName: "victim_contact_form1_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm1Nick":              {Name: "VictimContactForm1Nick", DBName: "victim_contact_form1_nick", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: false},
		"VictimContactForm1DocType":           {Name: "VictimContactForm1DocType", DBName: "victim_contact_form1_doc_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimContactForm1DocNumber":         {Name: "VictimContactForm1DocNumber", DBName: "victim_contact_form1_doc_number", Alias: "", ModelType: "uint", MinSize: 4, MaxSize: 32, Required: false},
		"VictimContactForm1BirthDate":         {Name: "VictimContactForm1BirthDate", DBName: "victim_contact_form1_birth_date", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm1TownCode":          {Name: "VictimContactForm1TownCode", DBName: "victim_contact_form1_town_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 8, Required: false},
		"VictimContactForm1Address":           {Name: "VictimContactForm1Address", DBName: "victim_contact_form1_address", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 120, Required: false},
		"VictimContactForm1Phone":             {Name: "VictimContactForm1Phone", DBName: "victim_contact_form1_phone", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 10, Required: true},
		"VictimContactForm1GenderIdentity":    {Name: "VictimContactForm1GenderIdentity", DBName: "victim_contact_form1_gender_identity", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimContactForm1SexualOrientation": {Name: "VictimContactForm1SexualOrientation", DBName: "victim_contact_form1_sexual_orientation", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimContactForm1Origin":            {Name: "VictimContactForm1Origin", DBName: "victim_contact_form1_origin", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"VictimContactForm1Occupation":        {Name: "VictimContactForm1Occupation", DBName: "victim_contact_form1_occupation", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"VictimContactForm1OccupationOther":   {Name: "VictimContactForm1OccupationOther", DBName: "victim_contact_form1_occupation_other", Alias: "", ModelType: "string", MinSize: 4, MaxSize: 32, Required: false},
		"VictimContactForm1FactsDescription":  {Name: "VictimContactForm1FactsDescription", DBName: "victim_contact_form1_facts_description", Alias: "", ModelType: "string", MinSize: 14, MaxSize: 10000, Required: true},
		"VictimContactForm1VictimContact":     {Name: "VictimContactForm1VictimContact", DBName: "victim_contact_form1_victim_contact", Alias: "", ModelType: "uint", Required: true},
	}
)

// VictimContactForm1DTO representa el Data Transfer Object para el contacto de la víctima.
// Este DTO se utiliza para transferir datos entre capas de la aplicación y para serialización JSON.
type VictimContactForm1DTO struct {
	VictimContactForm1Id                uint64    `json:"-"`
	VictimContactForm1ICode             string    `json:"icode"`
	VictimContactForm1CreationDate      time.Time `json:"creationDate"`
	VictimContactForm1UpdateDate        time.Time `json:"updateDate"`
	VictimContactForm1Nick              string    `json:"nick"`
	VictimContactForm1DocType           string    `json:"docType"`
	VictimContactForm1DocNumber         string    `json:"docNumber"`
	VictimContactForm1BirthDate         time.Time `json:"birthDate"`
	VictimContactForm1TownCode          string    `json:"townCode"`
	VictimContactForm1Address           string    `json:"address"`
	VictimContactForm1Phone             string    `json:"phone"`
	VictimContactForm1GenderIdentity    string    `json:"genderIdentity"`
	VictimContactForm1SexualOrientation string    `json:"sexualOrientation"`
	VictimContactForm1Origin            string    `json:"origin"`
	VictimContactForm1Occupation        string    `json:"occupation"`
	VictimContactForm1OccupationOther   string    `json:"occupationOther"`
	VictimContactForm1FactsDescription  string    `json:"factsDescription"`

	VictimContactForm1VictimContact interface{} `json:"victimContact"`
}

// VictimContactForm1PgDB representa el mapeo de la tabla de VictimContactForm1 en la base de datos PostgreSQL.
// Utiliza sql.Null* para permitir valores nulos en la base de datos.
type VictimContactForm1PgDB struct {
	VictimContactForm1Id                sql.NullInt64
	VictimContactForm1ICode             sql.NullString
	VictimContactForm1CreationDate      sql.NullTime
	VictimContactForm1UpdateDate        sql.NullTime
	VictimContactForm1Nick              sql.NullString
	VictimContactForm1DocType           sql.NullString
	VictimContactForm1DocNumber         sql.NullString
	VictimContactForm1BirthDate         sql.NullTime
	VictimContactForm1TownCode          sql.NullString
	VictimContactForm1Address           sql.NullString
	VictimContactForm1Phone             sql.NullString
	VictimContactForm1GenderIdentity    sql.NullString
	VictimContactForm1SexualOrientation sql.NullString
	VictimContactForm1Origin            sql.NullString
	VictimContactForm1Occupation        sql.NullString
	VictimContactForm1OccupationOther   sql.NullString
	VictimContactForm1FactsDescription  sql.NullString
	VictimContactForm1VictimContact     sql.NullInt64
}

// MarshalJSON implementa la interfaz json.Marshaler para formatear las fechas con un formato específico.
// Se crea un alias de VictimContactForm1DTO para evitar recursión infinita durante la serialización.
func (vcd VictimContactForm1DTO) MarshalJSON() ([]byte, error) {
	type Alias VictimContactForm1DTO // Alias para evitar recursión infinita

	// Se formatean las fechas de nacimiento, creación y actualización utilizando el formato definido en common_config.
	return json.Marshal(&struct {
		*Alias
		VictimContactForm1BirthDate    string `json:"birthDate"`
		VictimContactForm1CreationDate string `json:"creationDate"`
		VictimContactForm1UpdateDate   string `json:"updateDate"`
	}{
		Alias:                          (*Alias)(&vcd),
		VictimContactForm1BirthDate:    vcd.VictimContactForm1BirthDate.Format(common_config.DateTime.DATE_FORMAT),
		VictimContactForm1CreationDate: vcd.VictimContactForm1CreationDate.Format(common_config.DateTime.DATE_FORMAT),
		VictimContactForm1UpdateDate:   vcd.VictimContactForm1UpdateDate.Format(common_config.DateTime.DATE_FORMAT),
	})
}

func (vcd *VictimContactForm1DTO) UnmarshalJSON(data []byte) error {
	type Alias VictimContactForm1DTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		VictimContactForm1BirthDate    string `json:"birthDate"`
		VictimContactForm1CreationDate string `json:"creationDate"`
		VictimContactForm1UpdateDate   string `json:"updateDate"`
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

	// Parseo seguro según formato
	vcd.VictimContactForm1BirthDate = parse(aux.VictimContactForm1BirthDate, common_config.DateTime.DATE_FORMAT)

	vcd.VictimContactForm1CreationDate = parse(aux.VictimContactForm1CreationDate, common_config.DateTime.DATE_FORMAT)

	vcd.VictimContactForm1UpdateDate = parse(aux.VictimContactForm1UpdateDate, common_config.DateTime.DATE_FORMAT)

	return nil
}

// SetVictimContactForm1 inserta un nuevo registro de contacto de víctima en la base de datos.
// Parámetros:
//   - victimContactForm1: Puntero al objeto VictimContactForm1DTO que contiene los datos a insertar.

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un posible error durante la operación.
func SetVictimContactForm1(victimContactForm1 *VictimContactForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se define el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Definición de campos que se insertarán en la base de datos.
	var victimContactForm1FieldsSlice []string = []string{"VictimContactForm1ICode", "VictimContactForm1CreationDate", "VictimContactForm1UpdateDate",
		"VictimContactForm1Nick", "VictimContactForm1DocType", "VictimContactForm1DocNumber", "VictimContactForm1BirthDate", "VictimContactForm1TownCode", "VictimContactForm1Address",
		"VictimContactForm1Phone", "VictimContactForm1GenderIdentity", "VictimContactForm1SexualOrientation", "VictimContactForm1Origin",
		"VictimContactForm1Occupation", "VictimContactForm1OccupationOther", "VictimContactForm1FactsDescription"}
	var victimContactForm1FieldsAliasSlice []string = []string{}

	// Se construye la consulta SQL de inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, victimContactForm1FieldsSlice, victimContactForm1FieldsAliasSlice, VictimContactForm1DBName, []string{}, []string{}, []string{"VictimContactForm1Id"}, common_dao.SQL_AND, VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, false)

	// Se ejecuta la consulta con los valores correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query,
		victimContactForm1.VictimContactForm1ICode, victimContactForm1.VictimContactForm1CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), victimContactForm1.VictimContactForm1UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimContactForm1.VictimContactForm1Nick, victimContactForm1.VictimContactForm1DocType,
		victimContactForm1.VictimContactForm1DocNumber, victimContactForm1.VictimContactForm1BirthDate.Format(common_config.DateTime.DB_DATE_FORMAT), victimContactForm1.VictimContactForm1TownCode, victimContactForm1.VictimContactForm1Address,
		victimContactForm1.VictimContactForm1Phone, victimContactForm1.VictimContactForm1GenderIdentity, victimContactForm1.VictimContactForm1SexualOrientation,
		victimContactForm1.VictimContactForm1Origin, victimContactForm1.VictimContactForm1Occupation, victimContactForm1.VictimContactForm1OccupationOther, victimContactForm1.VictimContactForm1FactsDescription)

	// Se escanea el resultado para obtener el ID generado.
	persistenceCtrl.Scan(&victimContactForm1.VictimContactForm1Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetVictimContactForm1 obtiene un registro de contacto de víctima de la base de datos utilizando criterios de búsqueda.
// Parámetros:
//   - by: Criterio de búsqueda (estructura con nombres y valores de atributos).
//   - victimContactForm1: Puntero al objeto VictimContactForm1DTO donde se almacenarán los datos obtenidos.

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un posible error durante la operación.
func GetVictimContactForm1(by common_controllers.By, victimContactForm1 *VictimContactForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se define el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactForm1Path string = VictimContactForm1DBScheme + "." + VictimContactForm1DBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var victimContactForm1FieldsSlice []string = []string{"VictimContactForm1Id", "VictimContactForm1ICode", "VictimContactForm1CreationDate", "VictimContactForm1UpdateDate",
		"VictimContactForm1Nick", "VictimContactForm1DocType", "VictimContactForm1DocNumber",
		"VictimContactForm1BirthDate", "VictimContactForm1TownCode", "VictimContactForm1Address", "VictimContactForm1Phone",
		"VictimContactForm1GenderIdentity", "VictimContactForm1SexualOrientation", "VictimContactForm1Origin", "VictimContactForm1Occupation", "VictimContactForm1OccupationOther",
		"VictimContactForm1FactsDescription"}
	var victimContactForm1FieldsAliasSlice []string = []string{}

	// Se construye la parte de la consulta que define los campos a seleccionar.
	var victimContactForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactForm1FieldsSlice, victimContactForm1FieldsAliasSlice, VictimContactForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, true)

	// Se construye la consulta SQL completa, incluyendo la cláusula WHERE basada en el criterio de búsqueda.
	var query string = `SELECT ` + victimContactForm1FieldsStr +
		` FROM ` + victimContactForm1Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, true)

	// Se imprime la consulta para fines de depuración.
	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se define una variable temporal para almacenar los resultados de la consulta.
	var victimContactForm1PG VictimContactForm1PgDB = VictimContactForm1PgDB{}

	// Se escanean los resultados en la estructura correspondiente.
	persistenceCtrl.Scan(&victimContactForm1PG.VictimContactForm1Id, &victimContactForm1PG.VictimContactForm1ICode, &victimContactForm1PG.VictimContactForm1CreationDate, &victimContactForm1PG.VictimContactForm1UpdateDate,
		&victimContactForm1PG.VictimContactForm1Nick,
		&victimContactForm1PG.VictimContactForm1DocType, &victimContactForm1PG.VictimContactForm1DocNumber, &victimContactForm1PG.VictimContactForm1BirthDate, &victimContactForm1PG.VictimContactForm1TownCode,
		&victimContactForm1PG.VictimContactForm1Address, &victimContactForm1PG.VictimContactForm1Phone,
		&victimContactForm1PG.VictimContactForm1GenderIdentity, &victimContactForm1PG.VictimContactForm1SexualOrientation, &victimContactForm1PG.VictimContactForm1Origin, &victimContactForm1PG.VictimContactForm1Occupation,
		&victimContactForm1PG.VictimContactForm1OccupationOther, &victimContactForm1PG.VictimContactForm1FactsDescription)

	// Se convierte la estructura obtenida a DTO y se asigna al puntero recibido.
	var tmp VictimContactForm1DTO = victimContactForm1PG.ToDTO()
	*victimContactForm1 = tmp
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetVictimContactForm1s obtiene múltiples registros de contacto de víctima basándose en un criterio de búsqueda.
// Parámetros:
//   - by: Criterio de búsqueda (estructura con nombres y valores de atributos).

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactForm1DTO y un posible error durante la operación.
func GetVictimContactForm1s(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactForm1DTO, error) {
	// Se define el controlador de persistencia y las variables necesarias.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactForm1Path string = VictimContactForm1DBScheme + "." + VictimContactForm1DBName

	var victimContactForm1 = VictimContactForm1DTO{}
	var victimContactForm1s []VictimContactForm1DTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var victimContactForm1FieldsSlice []string = []string{"VictimContactForm1ICode", "VictimContactForm1CreationDate",
		"VictimContactForm1Nick", "VictimContactForm1DocType", "VictimContactForm1DocNumber", "VictimContactForm1TownCode", "VictimContactForm1Address",
		"VictimContactForm1Phone"}
	var victimContactForm1FieldsAliasSlice []string = []string{}

	// Se construye la parte de la consulta que define los campos a seleccionar.
	var victimContactForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactForm1FieldsSlice, victimContactForm1FieldsAliasSlice, VictimContactForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, true)

	// Se construye la consulta SQL completa, incluyendo cláusula WHERE y ordenamiento.
	var query string = `SELECT ` + victimContactForm1FieldsStr +
		` FROM ` + victimContactForm1Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactForm1DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, true) +
		` ORDER BY ` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1CreationDate"].DBName + ` ASC  `

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Se itera sobre los resultados y se almacenan en el slice.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&victimContactForm1.VictimContactForm1ICode, &victimContactForm1.VictimContactForm1CreationDate,
			&victimContactForm1.VictimContactForm1Nick,
			&victimContactForm1.VictimContactForm1DocType, &victimContactForm1.VictimContactForm1DocNumber, &victimContactForm1.VictimContactForm1TownCode,
			&victimContactForm1.VictimContactForm1Address, &victimContactForm1.VictimContactForm1Phone)
		victimContactForm1s = append(victimContactForm1s, victimContactForm1)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return victimContactForm1s, nil
}

// GetVictimContactForm1sWithoutVictimCase obtiene los registros de contacto de víctima que no tienen asociados un caso de víctima.
// Parámetros:
//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactForm1DTO y un posible error durante la operación.
func GetVictimContactForm1sWithoutVictimCase(page int, status string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactForm1DTO, int, error) {
	// Se define el controlador de persistencia y las variables necesarias.
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactForm1Path string = VictimContactForm1DBScheme + "." + VictimContactForm1DBName
	// Se asume que VictimCaseDBScheme y VictimCaseDBName y VictimCaseFieldDefinitions están definidos en otro lugar del proyecto.
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName

	var victimContactForm1 = VictimContactForm1DTO{}
	var victimContactForm1s []VictimContactForm1DTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var victimContactForm1FieldsSlice []string = []string{"VictimContactForm1ICode", "VictimContactForm1CreationDate",
		"VictimContactForm1Nick", "VictimContactForm1DocType", "VictimContactForm1DocNumber", "VictimContactForm1TownCode", "VictimContactForm1Address",
		"VictimContactForm1Phone"}
	var victimContactForm1FieldsAliasSlice []string = []string{}

	// Se construye la parte de la consulta que define los campos a seleccionar.
	var victimContactForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactForm1FieldsSlice, victimContactForm1FieldsAliasSlice, VictimContactForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, true)

	// Se construye la consulta SQL con LEFT JOIN para identificar registros sin caso asociado.
	var query string = `SELECT ` + victimContactForm1FieldsStr +
		` FROM ` + victimContactForm1Path +
		` LEFT JOIN ` + victimCasePath + ` ON (` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1Id"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContactForm1"].DBName + `)` +
		` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContactForm1"].DBName + ` IS NULL  AND ` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1Status"].DBName + ` = $1` +
		` ORDER BY ` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, status)

	// Se itera sobre los resultados y se almacena cada registro en el slice.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&victimContactForm1.VictimContactForm1ICode, &victimContactForm1.VictimContactForm1CreationDate,
			&victimContactForm1.VictimContactForm1Nick,
			&victimContactForm1.VictimContactForm1DocType, &victimContactForm1.VictimContactForm1DocNumber, &victimContactForm1.VictimContactForm1TownCode,
			&victimContactForm1.VictimContactForm1Address, &victimContactForm1.VictimContactForm1Phone)
		victimContactForm1s = append(victimContactForm1s, victimContactForm1)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimContactForm1Path +
			` LEFT JOIN ` + victimCasePath + ` ON (` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1Id"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContactForm1"].DBName + `)` +
			` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseVictimContactForm1"].DBName + ` IS NULL AND ` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1Status"].DBName + ` = 'v'`

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return victimContactForm1s, count, nil
}

// GetAllVictimContactForm1 obtiene todos los registros de contacto de víctima de la base de datos.
// Parámetros:
//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactForm1DTO y un posible error durante la operación.
func GetAllVictimContactForm1(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactForm1DTO, error) {
	// Se define el controlador de persistencia y las variables necesarias.
	var victimContactForm1 VictimContactForm1DTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimContactForm1Path string = VictimContactForm1DBScheme + "." + VictimContactForm1DBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de campos a seleccionar.
	var victimContactForm1FieldsSlice []string = []string{"VictimContactForm1ICode", "VictimContactForm1CreationDate", "VictimContactForm1UpdateDate",
		"VictimContactForm1Nick", "VictimContactForm1DocType", "VictimContactForm1DocNumber", "VictimContactForm1TownCode", "VictimContactForm1Address",
		"VictimContactForm1Phone"}
	var victimContactForm1FieldsAliasSlice []string = []string{}

	// Se construye la consulta SQL completa.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, victimContactForm1FieldsSlice, victimContactForm1FieldsAliasSlice, VictimContactForm1DBName, []string{}, []string{}, []string{}, "", VictimContactForm1DBScheme, VictimContactForm1FieldDefinitions, true) +
		` ORDER BY ` + victimContactForm1Path + `.` + VictimContactForm1FieldDefinitions["VictimContactForm1CreationDate"].DBName + ` ASC  `

	// Se ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var victimContactForm1s []VictimContactForm1DTO
	// Se itera sobre los resultados y se almacena cada registro en el slice.
	for persistenceCtrl.Next() {
		victimContactForm1 = VictimContactForm1DTO{}
		persistenceCtrl.ScanRow(&victimContactForm1.VictimContactForm1ICode, &victimContactForm1.VictimContactForm1CreationDate, &victimContactForm1.VictimContactForm1UpdateDate,
			&victimContactForm1.VictimContactForm1Nick,
			&victimContactForm1.VictimContactForm1DocType, &victimContactForm1.VictimContactForm1DocNumber, &victimContactForm1.VictimContactForm1TownCode,
			&victimContactForm1.VictimContactForm1Address, &victimContactForm1.VictimContactForm1Phone)
		victimContactForm1s = append(victimContactForm1s, victimContactForm1)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimContactForm1s, nil
}

// SetVictimContactForm1Defaults asigna valores predeterminados a un objeto VictimContactForm1DTO según la acción a realizar.
// Parámetros:
//   - victimContactForm1: Puntero al objeto VictimContactForm1DTO al que se asignarán los valores predeterminados.
//   - action: Acción que se va a realizar (por ejemplo, common_dao.SQL_INSERT o common_dao.SQL_UPDATE).
func SetVictimContactForm1Defaults(victimContactForm1 *VictimContactForm1DTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se asignan la fecha de creación y actualización al momento actual,
		// se genera un código único y se asigna un estado predeterminado.
		victimContactForm1.VictimContactForm1CreationDate = time.Now()
		victimContactForm1.VictimContactForm1UpdateDate = time.Now()
		victimContactForm1.VictimContactForm1ICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para una actualización, solo se asigna la fecha de actualización.
		victimContactForm1.VictimContactForm1UpdateDate = time.Now()
	}

}

// ToDTO convierte una instancia de VictimContactForm1PgDB a su correspondiente VictimContactForm1DTO.
// Este método verifica la validez de cada campo (usando sql.Null*) y asigna el valor correspondiente.
func (obj *VictimContactForm1PgDB) ToDTO() VictimContactForm1DTO {
	var dto VictimContactForm1DTO

	if obj.VictimContactForm1Id.Valid {
		dto.VictimContactForm1Id = uint64(obj.VictimContactForm1Id.Int64)
	}

	if obj.VictimContactForm1ICode.Valid {
		dto.VictimContactForm1ICode = obj.VictimContactForm1ICode.String
	}

	if obj.VictimContactForm1CreationDate.Valid {
		dto.VictimContactForm1CreationDate = obj.VictimContactForm1CreationDate.Time
	}

	if obj.VictimContactForm1UpdateDate.Valid {
		dto.VictimContactForm1UpdateDate = obj.VictimContactForm1UpdateDate.Time
	}

	if obj.VictimContactForm1Nick.Valid {
		dto.VictimContactForm1Nick = obj.VictimContactForm1Nick.String
	}
	if obj.VictimContactForm1DocType.Valid {
		dto.VictimContactForm1DocType = obj.VictimContactForm1DocType.String
	}
	if obj.VictimContactForm1DocNumber.Valid {
		dto.VictimContactForm1DocNumber = obj.VictimContactForm1DocNumber.String
	}
	if obj.VictimContactForm1BirthDate.Valid {
		dto.VictimContactForm1BirthDate = obj.VictimContactForm1BirthDate.Time
	}
	if obj.VictimContactForm1TownCode.Valid {
		dto.VictimContactForm1TownCode = obj.VictimContactForm1TownCode.String
	}
	if obj.VictimContactForm1Address.Valid {
		dto.VictimContactForm1Address = obj.VictimContactForm1Address.String
	}

	if obj.VictimContactForm1Phone.Valid {
		dto.VictimContactForm1Phone = obj.VictimContactForm1Phone.String
	}
	if obj.VictimContactForm1GenderIdentity.Valid {
		dto.VictimContactForm1GenderIdentity = obj.VictimContactForm1GenderIdentity.String
	}
	if obj.VictimContactForm1SexualOrientation.Valid {
		dto.VictimContactForm1SexualOrientation = obj.VictimContactForm1SexualOrientation.String
	}
	if obj.VictimContactForm1Origin.Valid {
		dto.VictimContactForm1Origin = obj.VictimContactForm1Origin.String
	}
	if obj.VictimContactForm1Occupation.Valid {
		dto.VictimContactForm1Occupation = obj.VictimContactForm1Occupation.String
	}
	if obj.VictimContactForm1OccupationOther.Valid {
		dto.VictimContactForm1OccupationOther = obj.VictimContactForm1OccupationOther.String
	}
	if obj.VictimContactForm1FactsDescription.Valid {
		dto.VictimContactForm1FactsDescription = obj.VictimContactForm1FactsDescription.String
	}

	if obj.VictimContactForm1VictimContact.Valid {
		dto.VictimContactForm1VictimContact = VictimContactDTO{VictimContactId: uint64(obj.VictimContactForm1VictimContact.Int64)}
	}

	return dto
}
