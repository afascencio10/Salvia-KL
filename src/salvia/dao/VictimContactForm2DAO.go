// Package salvia_daos implementa los Data Access Objects (DAO) relacionados con el contacto de la víctima.
// Contiene funciones para insertar, obtener, actualizar y eliminar registros de contacto de la víctima en la base de datos.
package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"strconv"

	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

var (
	VictimContactForm2EntityName string = "VictimContactForm2"
	VictimContactForm2JSONName   string = "form2"
	VictimContactForm2DBName     string = "victim_contact_form2"
	VictimContactForm2DBScheme   string = "salvia"

	// Field definitions – used by the ORM / validation layer
	VictimContactForm2FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimContactForm2Id":                  {Name: "VictimContactForm2Id", DBName: "victim_contact_form2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2ICode":               {Name: "VictimContactForm2ICode", DBName: "victim_case_form2_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimContactForm2CreationDate":        {Name: "VictimContactForm2CreationDate", DBName: "victim_case_form2_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2UpdateDate":          {Name: "VictimContactForm2UpdateDate", DBName: "victim_case_form2_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2WillReceiveCall":     {Name: "VictimContactForm2WillReceiveCall", DBName: "victim_contact_form2_will_receive_call", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm2HasCareRole":         {Name: "VictimContactForm2HasCareRole", DBName: "victim_contact_form2_has_care_role", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm2VictimAwareOfReport": {Name: "VictimContactForm2VictimAwareOfReport", DBName: "victim_contact_form2_victim_aware_of_report", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm2ReporterNames":       {Name: "VictimContactForm2ReporterNames", DBName: "victim_contact_form2_reporter_names", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 64, Required: false},
		"VictimContactForm2ReporterPhone":       {Name: "VictimContactForm2ReporterPhone", DBName: "victim_contact_form2_reporter_phone", Alias: "", ModelType: "uint", MinSize: 1000000000, MaxSize: 9999999999, Required: false},
		"VictimContactForm2VictimColPhone":      {Name: "VictimContactForm2VictimColPhone", DBName: "victim_contact_form2_victim_col_phone", Alias: "", ModelType: "uint", MinSize: 1000000000, MaxSize: 9999999999, Required: true},
		"VictimContactForm2FactsDescription":    {Name: "VictimContactForm2FactsDescription", DBName: "victim_contact_form2_facts_description", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2BestContactTime":     {Name: "VictimContactForm2BestContactTime", DBName: "victim_contact_form2_best_contact_time", Alias: "", ModelType: "time", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2ReportType":          {Name: "VictimContactForm2ReportType", DBName: "victim_contact_form2_report_type", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2VictimContact":       {Name: "VictimContactForm2VictimContact", DBName: "victim_contact_form2_victim_contact", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// ---------------------------------------------------------------------------
// DTO – the structure that is sent/received via JSON
// ---------------------------------------------------------------------------
type VictimContactForm2DTO struct {
	VictimContactForm2Id                  uint64                  `json:"-"`
	VictimContactForm2ICode               string                  `json:"icode"`
	VictimContactForm2CreationDate        time.Time               `json:"creationDate"`
	VictimContactForm2UpdateDate          time.Time               `json:"updateDate"`
	VictimContactForm2WillReceiveCall     VictimCaseForm2EnumsDTO `json:"willReceiveCall"`
	VictimContactForm2HasCareRole         VictimCaseForm2EnumsDTO `json:"hasCareRole"`
	VictimContactForm2VictimAwareOfReport VictimCaseForm2EnumsDTO `json:"victimAwareOfReport"`
	VictimContactForm2ReporterNames       string                  `json:"reporterNames"`
	VictimContactForm2ReporterPhone       uint64                  `json:"reporterPhone"`
	VictimContactForm2VictimColPhone      uint64                  `json:"victimColPhone"`
	VictimContactForm2FactsDescription    string                  `json:"factsDescription"`
	VictimContactForm2BestContactTime     time.Time               `json:"bestContactTime"`
	VictimContactForm2ReportType          VictimCaseForm2EnumsDTO `json:"reportType"`
	VictimContactForm2VictimContact       interface{}             `json:"-"`

	// Parámetros auxiliares
	VictimContactForm2AuthorizationAnswer VictimCaseForm2EnumsDTO `json:"authorizationAnswer"`
}

// ---------------------------------------------------------------------------
// PgDB – the structure that maps directly to PostgreSQL rows
// ---------------------------------------------------------------------------
type VictimContactForm2PgDB struct {
	VictimContactForm2Id                  sql.NullString
	VictimContactForm2ICode               sql.NullString
	VictimContactForm2CreationDate        sql.NullString
	VictimContactForm2UpdateDate          sql.NullString
	VictimContactForm2WillReceiveCall     sql.NullString
	VictimContactForm2HasCareRole         sql.NullString
	VictimContactForm2VictimAwareOfReport sql.NullString
	VictimContactForm2ReporterNames       sql.NullString
	VictimContactForm2ReporterPhone       sql.NullInt64
	VictimContactForm2VictimColPhone      sql.NullInt64
	VictimContactForm2FactsDescription    sql.NullString
	VictimContactForm2BestContactTime     sql.NullString
	VictimContactForm2ReportType          sql.NullString
	VictimContactForm2VictimContact       sql.NullString
}

func (vcd VictimContactForm2DTO) MarshalJSON() ([]byte, error) {
	type Alias VictimContactForm2DTO

	return json.Marshal(&struct {
		*Alias
		VictimContactForm2CreationDate    string `json:"creationDate"`
		VictimContactForm2UpdateDate      string `json:"updateDate"`
		VictimContactForm2BestContactTime string `json:"bestContactTime"`
	}{
		Alias:                             (*Alias)(&vcd),
		VictimContactForm2CreationDate:    vcd.VictimContactForm2CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimContactForm2UpdateDate:      vcd.VictimContactForm2UpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimContactForm2BestContactTime: vcd.VictimContactForm2BestContactTime.Format(common_config.DateTime.TIME_FORMAT),
	})
}

// ---------------------------------------------------------------------------
// JSON unmarshalling – parse dates/times safely
// ---------------------------------------------------------------------------
func (vcd *VictimContactForm2DTO) UnmarshalJSON(data []byte) error {
	type Alias VictimContactForm2DTO

	aux := &struct {
		*Alias
		VictimContactForm2CreationDate    string          `json:"creationDate"`
		VictimContactForm2UpdateDate      string          `json:"updateDate"`
		VictimContactForm2BestContactTime string          `json:"bestContactTime"`
		VictimContactForm2ReporterPhone   json.RawMessage `json:"reporterPhone"`
		VictimContactForm2VictimColPhone  json.RawMessage `json:"victimColPhone"`
	}{
		Alias: (*Alias)(vcd),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{}
		}
		return t
	}

	parseUint := func(raw json.RawMessage) uint64 {
		if len(raw) == 0 {
			return 0
		}

		// Intentar como número
		var num uint64
		if err := json.Unmarshal(raw, &num); err == nil {
			return num
		}

		// Intentar como string
		var str string
		if err := json.Unmarshal(raw, &str); err == nil {
			n, err := strconv.ParseUint(str, 10, 64)
			if err == nil {
				return n
			}
		}

		return 0
	}

	vcd.VictimContactForm2CreationDate = parse(aux.VictimContactForm2CreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimContactForm2UpdateDate = parse(aux.VictimContactForm2UpdateDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimContactForm2BestContactTime = parse(aux.VictimContactForm2BestContactTime, common_config.DateTime.TIME_FORMAT)
	vcd.VictimContactForm2ReporterPhone = parseUint(aux.VictimContactForm2ReporterPhone)
	vcd.VictimContactForm2VictimColPhone = parseUint(aux.VictimContactForm2VictimColPhone)

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
func SetVictimContactForm2(victimContactForm2 *VictimContactForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	victimContactForm2Fields := []string{
		"VictimContactForm2ICode",
		"VictimContactForm2CreationDate",
		"VictimContactForm2UpdateDate",
		"VictimContactForm2WillReceiveCall",
		"VictimContactForm2HasCareRole",
		"VictimContactForm2VictimAwareOfReport",
		"VictimContactForm2ReporterNames",
		"VictimContactForm2ReporterPhone",
		"VictimContactForm2VictimColPhone",
		"VictimContactForm2FactsDescription",
		"VictimContactForm2BestContactTime",
		"VictimContactForm2ReportType",
		"VictimContactForm2VictimContact",
	}

	//Hacemos null los campos que no vengan con datos

	var victimContactForm2AliasSlice []string = []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, victimContactForm2Fields, victimContactForm2AliasSlice, VictimContactForm2DBName, []string{}, []string{}, []string{"VictimContactForm2Id"}, common_dao.SQL_AND, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		victimContactForm2.VictimContactForm2ICode,
		victimContactForm2.VictimContactForm2CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimContactForm2.VictimContactForm2UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		utils.NilIfZero(victimContactForm2.VictimContactForm2WillReceiveCall.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimContactForm2.VictimContactForm2HasCareRole.VictimCaseForm2EnumsId),
		utils.NilIfZero(victimContactForm2.VictimContactForm2VictimAwareOfReport.VictimCaseForm2EnumsId),
		victimContactForm2.VictimContactForm2ReporterNames,
		victimContactForm2.VictimContactForm2ReporterPhone,
		victimContactForm2.VictimContactForm2VictimColPhone,
		victimContactForm2.VictimContactForm2FactsDescription,
		victimContactForm2.VictimContactForm2BestContactTime.Format(common_config.DateTime.TIME_FORMAT),
		utils.NilIfZero(victimContactForm2.VictimContactForm2ReportType.VictimCaseForm2EnumsId),
		victimContactForm2.VictimContactForm2VictimContact.(VictimContactDTO).VictimContactId)

	persistenceCtrl.Scan(&victimContactForm2.VictimContactForm2Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetVictimContact obtiene un registro de contacto de víctima de la base de datos utilizando criterios de búsqueda.
// Parámetros:
//   - by: Criterio de búsqueda (estructura con nombres y valores de atributos).
//   - victimContact: Puntero al objeto VictimContactForm2DTO donde se almacenarán los datos obtenidos.

//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un posible error durante la operación.
func GetVictimContactForm2(by common_controllers.By, victimContactForm2 *VictimContactForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController

	victimContactForm2Path := VictimContactForm2DBScheme + "." + VictimContactForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	victimContactForm2Fields := []string{
		"VictimContactForm2Id",
		"VictimContactForm2ICode",
		"VictimContactForm2CreationDate",
		"VictimContactForm2UpdateDate",
		"VictimContactForm2WillReceiveCall",
		"VictimContactForm2HasCareRole",
		"VictimContactForm2VictimAwareOfReport",
		"VictimContactForm2ReporterNames",
		"VictimContactForm2ReporterPhone",
		"VictimContactForm2VictimColPhone",
		"VictimContactForm2FactsDescription",
		"VictimContactForm2BestContactTime",
		"VictimContactForm2ReportType",
		"VictimContactForm2VictimContact",
	}

	var victimContactForm2AliasSlice []string = []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactForm2Fields, victimContactForm2AliasSlice, VictimContactForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimContactForm2Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactForm2DBName, by.AttrsName, nil, nil, by.Operator, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var victimContactForm2Pg VictimContactForm2PgDB
	persistenceCtrl.Scan(
		&victimContactForm2Pg.VictimContactForm2Id,
		&victimContactForm2Pg.VictimContactForm2ICode,
		&victimContactForm2Pg.VictimContactForm2CreationDate,
		&victimContactForm2Pg.VictimContactForm2UpdateDate,
		&victimContactForm2Pg.VictimContactForm2WillReceiveCall,
		&victimContactForm2Pg.VictimContactForm2HasCareRole,
		&victimContactForm2Pg.VictimContactForm2VictimAwareOfReport,
		&victimContactForm2Pg.VictimContactForm2ReporterNames,
		&victimContactForm2Pg.VictimContactForm2ReporterPhone,
		&victimContactForm2Pg.VictimContactForm2VictimColPhone,
		&victimContactForm2Pg.VictimContactForm2FactsDescription,
		&victimContactForm2Pg.VictimContactForm2BestContactTime,
		&victimContactForm2Pg.VictimContactForm2ReportType,
		&victimContactForm2Pg.VictimContactForm2VictimContact)

	*victimContactForm2 = victimContactForm2Pg.ToDTO()

	//Agregamos los enums desde memoria
	GetLocalVictimCaseForm2EnumsById(&victimContactForm2.VictimContactForm2WillReceiveCall)
	GetLocalVictimCaseForm2EnumsById(&victimContactForm2.VictimContactForm2HasCareRole)
	GetLocalVictimCaseForm2EnumsById(&victimContactForm2.VictimContactForm2VictimAwareOfReport)
	GetLocalVictimCaseForm2EnumsById(&victimContactForm2.VictimContactForm2ReportType)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
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
//   - un slice de VictimContactForm2DTO y un posible error durante la operación.
func GetVictimContactsForm2(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimContactForm2DTO, int, error) {
	var count int

	persistenceCtrl := common_controllers.PersistenceController{}
	victimContactForm2Path := VictimContactForm2DBScheme + "." + VictimContactForm2DBName
	victimContactsForm2 := []VictimContactForm2DTO{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	victimContactForm2Fields := []string{
		"VictimContactForm2UpdateDate",
		"VictimContactForm2WillReceiveCall",
		"VictimContactForm2HasCareRole",
		"VictimContactForm2VictimAwareOfReport",
		"VictimContactForm2ReporterNames",
		"VictimContactForm2ReporterPhone",
		"VictimContactForm2VictimColPhone",
		"VictimContactForm2FactsDescription",
		"VictimContactForm2BestContactTime",
		"VictimContactForm2ReportType",
		"VictimContactForm2VictimContact",
	}

	var victimContactForm2AliasSlice []string = []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimContactForm2Fields, victimContactForm2AliasSlice, VictimContactForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimContactForm2Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, true) +
		` ORDER BY ` + victimContactForm2Path + `.` + VictimContactForm2FieldDefinitions["VictimContactForm2CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	for persistenceCtrl.Next() {
		var victimContactForm2Pg VictimContactForm2PgDB
		persistenceCtrl.ScanRow(
			&victimContactForm2Pg.VictimContactForm2ICode,
			&victimContactForm2Pg.VictimContactForm2CreationDate,
			&victimContactForm2Pg.VictimContactForm2UpdateDate,
			&victimContactForm2Pg.VictimContactForm2WillReceiveCall,
			&victimContactForm2Pg.VictimContactForm2HasCareRole,
			&victimContactForm2Pg.VictimContactForm2VictimAwareOfReport,
			&victimContactForm2Pg.VictimContactForm2ReporterNames,
			&victimContactForm2Pg.VictimContactForm2ReporterPhone,
			&victimContactForm2Pg.VictimContactForm2VictimColPhone,
			&victimContactForm2Pg.VictimContactForm2FactsDescription,
			&victimContactForm2Pg.VictimContactForm2BestContactTime,
			&victimContactForm2Pg.VictimContactForm2ReportType,
			&victimContactForm2Pg.VictimContactForm2VictimContact)

		victimContactsForm2 = append(victimContactsForm2, victimContactForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Count total rows for the first page
	if page == 0 {
		countQuery := `SELECT COUNT(*) FROM ` + victimContactForm2Path +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return victimContactsForm2, count, nil
}

// GetAllVictimContactsForm2 obtiene todos los registros de contacto de víctima de la base de datos.
// Parámetros:
//   - connData: Datos de conexión a la base de datos.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - un slice de VictimContactForm2DTO y un posible error durante la operación.

func GetAllVictimContactsForm2(page int,
	connData *db.ConnData,
	clientConfig *db.DBClientConfig,
	serverConfig *db.DBServerConfig) ([]VictimContactForm2DTO, int, error) {

	var count int

	persistenceCtrl := common_controllers.PersistenceController{}
	victimContactForm2Path := VictimContactForm2DBScheme + "." + VictimContactForm2DBName
	victimContactsForm2 := []VictimContactForm2DTO{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fields := []string{
		"VictimContactForm2Id",
		"VictimContactForm2ICode",
		"VictimContactForm2CreationDate",
		"VictimContactForm2UpdateDate",
		"VictimContactForm2WillReceiveCall",
		"VictimContactForm2HasCareRole",
		"VictimContactForm2VictimAwareOfReport",
		"VictimContactForm2ReporterNames",
		"VictimContactForm2ReporterPhone",
		"VictimContactForm2VictimColPhone",
		"VictimContactForm2FactsDescription",
		"VictimContactForm2BestContactTime",
		"VictimContactForm2ReportType",
		"VictimContactForm2VictimContact",
	}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fields, nil, VictimContactForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm2DBScheme, VictimContactForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimContactForm2Path +
		` ORDER BY ` + victimContactForm2Path + `.` + VictimContactForm2FieldDefinitions["VictimContactForm2CreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query)
	for persistenceCtrl.Next() {
		var victimContactForm2Pg VictimContactForm2PgDB
		persistenceCtrl.ScanRow(
			&victimContactForm2Pg.VictimContactForm2ICode,
			&victimContactForm2Pg.VictimContactForm2CreationDate,
			&victimContactForm2Pg.VictimContactForm2UpdateDate,
			&victimContactForm2Pg.VictimContactForm2WillReceiveCall,
			&victimContactForm2Pg.VictimContactForm2HasCareRole,
			&victimContactForm2Pg.VictimContactForm2VictimAwareOfReport,
			&victimContactForm2Pg.VictimContactForm2ReporterNames,
			&victimContactForm2Pg.VictimContactForm2ReporterPhone,
			&victimContactForm2Pg.VictimContactForm2VictimColPhone,
			&victimContactForm2Pg.VictimContactForm2FactsDescription,
			&victimContactForm2Pg.VictimContactForm2BestContactTime,
			&victimContactForm2Pg.VictimContactForm2ReportType,
			&victimContactForm2Pg.VictimContactForm2VictimContact)

		victimContactsForm2 = append(victimContactsForm2, victimContactForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Count total rows for the first page
	if page == 0 {
		countQuery := `SELECT COUNT(*) FROM ` + victimContactForm2Path

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return victimContactsForm2, count, nil
}

// SetVictimContactDefaults asigna valores predeterminados a un objeto VictimContactForm2DTO según la acción a realizar.
// Parámetros:
//   - victimContact: Puntero al objeto VictimContactForm2DTO al que se asignarán los valores predeterminados.
//   - action: Acción que se va a realizar (por ejemplo, common_dao.SQL_INSERT o common_dao.SQL_UPDATE).
func SetVictimContactForm2Defaults(victimContactForm2 *VictimContactForm2DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		victimContactForm2.VictimContactForm2CreationDate = time.Now()
		victimContactForm2.VictimContactForm2UpdateDate = time.Now()
		victimContactForm2.VictimContactForm2ICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		victimContactForm2.VictimContactForm2UpdateDate = time.Now()
	}
}

// ToDTO convierte una instancia de VictimContactPgDB a su correspondiente VictimContactDTO.
// Este método verifica la validez de cada campo (usando sql.Null*) y asigna el valor correspondiente.
func (obj *VictimContactForm2PgDB) ToDTO() VictimContactForm2DTO {
	dto := VictimContactForm2DTO{}

	if obj.VictimContactForm2Id.Valid {
		dto.VictimContactForm2Id, _ = strconv.ParseUint(obj.VictimContactForm2Id.String, 10, 64)
	}
	if obj.VictimContactForm2ICode.Valid {
		dto.VictimContactForm2ICode = obj.VictimContactForm2ICode.String
	}
	if obj.VictimContactForm2CreationDate.Valid {
		dto.VictimContactForm2CreationDate, _ = time.Parse(common_config.DateTime.DATE_TIME_FORMAT, obj.VictimContactForm2CreationDate.String)
	}
	if obj.VictimContactForm2UpdateDate.Valid {
		dto.VictimContactForm2UpdateDate, _ = time.Parse(common_config.DateTime.DATE_TIME_FORMAT, obj.VictimContactForm2UpdateDate.String)
	}
	if obj.VictimContactForm2WillReceiveCall.Valid {
		dto.VictimContactForm2WillReceiveCall = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: utils.ParseUint64(obj.VictimContactForm2WillReceiveCall.String)}
	}
	if obj.VictimContactForm2HasCareRole.Valid {
		dto.VictimContactForm2HasCareRole = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: utils.ParseUint64(obj.VictimContactForm2HasCareRole.String)}
	}
	if obj.VictimContactForm2VictimAwareOfReport.Valid {
		dto.VictimContactForm2VictimAwareOfReport = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: utils.ParseUint64(obj.VictimContactForm2VictimAwareOfReport.String)}
	}

	if obj.VictimContactForm2ReporterNames.Valid {
		dto.VictimContactForm2ReporterNames = obj.VictimContactForm2ReporterNames.String
	}
	if obj.VictimContactForm2ReporterPhone.Valid {
		dto.VictimContactForm2ReporterPhone = uint64(obj.VictimContactForm2ReporterPhone.Int64)
	}
	if obj.VictimContactForm2VictimColPhone.Valid {
		dto.VictimContactForm2VictimColPhone = uint64(obj.VictimContactForm2VictimColPhone.Int64)
	}

	if obj.VictimContactForm2FactsDescription.Valid {
		dto.VictimContactForm2FactsDescription = obj.VictimContactForm2FactsDescription.String
	}
	if obj.VictimContactForm2BestContactTime.Valid {
		dto.VictimContactForm2BestContactTime, _ = time.Parse(common_config.DateTime.TIME_WITH_MILLISECONDS_FORMAT, obj.VictimContactForm2BestContactTime.String)
	}
	if obj.VictimContactForm2ReportType.Valid {
		dto.VictimContactForm2ReportType = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: utils.ParseUint64(obj.VictimContactForm2ReportType.String)}
	}
	if obj.VictimContactForm2VictimContact.Valid {
		dto.VictimContactForm2VictimContact = VictimCaseDTO{VictimCaseId: utils.ParseUint64(obj.VictimContactForm2VictimContact.String)}
	}

	return dto
}
