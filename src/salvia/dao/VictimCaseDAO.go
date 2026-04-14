package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"strconv"

	"encoding/json"

	security_daos "bitsflow/security/dao"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	VictimCaseEntityName string = "VictimCase"
	VictimCaseJSONName   string = "victimCase"
	VictimCaseDBName     string = "victim_case"
	VictimCaseDBScheme   string = "salvia"

	//Atributos relacionados con las validaciones ------------------------------

	//Campos que vienen como string del JSON. El booleano indica si son strings en el modelo o no (como en el caso de una fecha)
	VictimCaseFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimCaseId":               {Name: "VictimCaseId", DBName: "victim_case_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseICode":            {Name: "VictimCaseICode", DBName: "victim_case_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimCaseCreationDate":     {Name: "VictimCaseCreationDate", DBName: "victim_case_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseUpdateDate":       {Name: "VictimCaseUpdateDate", DBName: "victim_case_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseStatus":           {Name: "VictimCaseStatus", DBName: "victim_case_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: true},
		"VictimCaseGeneralUser":      {Name: "VictimCaseGeneralUser", DBName: "victim_case_general_user", Alias: "", ModelType: "string", MinSize: 32, MaxSize: 36, Required: true},
		"VictimCaseVictimContact":    {Name: "VictimCaseVictimContact", DBName: "victim_case_victim_contact", Alias: "", ModelType: "uint", Required: false},
		"VictimCaseApprovedBy":       {Name: "VictimCaseApprovedBy", DBName: "victim_case_approved_by", Alias: "", ModelType: "uint", Required: false},
		"VictimCaseFollowUp":         {Name: "VictimCaseFollowUp", DBName: "victim_case_follow_up", Alias: "", ModelType: "uint", Required: false},
		"VictimCaseOwnerDescription": {Name: "VictimCaseOwnerDescription", DBName: "victim_case_owner_description", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 512, Required: false},
		"VictimCaseTownCode":         {Name: "VictimCaseTownCode", DBName: "victim_case_victim_town_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 8, Required: true},

		"VictimCaseNames":     {Name: "VictimCaseNames", DBName: "victim_case_victim_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"VictimCaseLastNames": {Name: "VictimCaseLastNames", DBName: "victim_case_victim_last_names", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"VictimCaseDocType":   {Name: "VictimCaseDocType", DBName: "victim_case_victim_doc_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"VictimCaseDocNumber": {Name: "VictimCaseDocNumber", DBName: "victim_case_victim_doc_number", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
	}
)

type VictimCaseDTO struct {
	VictimCaseId           uint64    `json:"-"`
	VictimCaseICode        string    `json:"icode"`
	VictimCaseCreationDate time.Time `json:"creationDate"`
	VictimCaseUpdateDate   time.Time `json:"updateDate"`
	VictimCaseStatus       string    `json:"status"`
	VictimCaseGeneralUser  string    `json:"user"`
	VictimCaseTownCode     string    `json:"townCode"`

	VictimCaseVictimContact VictimContactDTO `json:"victimContact"`
	VictimCaseApprovedBy    CaseOwnerDTO     `json:"approvedBy"`
	VictimCaseFollowUp      FollowUpDTO      `json:"followUp"`

	VictimCaseOwnerDescription string `json:"owners"`

	VictimCaseNames     string `json:"names"`
	VictimCaseLastNames string `json:"lastNames"`
	VictimCaseDocType   string `json:"docType"`
	VictimCaseDocNumber string `json:"docNumber"`

	//Campos de formulario que no hacen parte del modelo o no directamente en la BD -----------------------------------
	VictimCaseEntityBranches     map[string]map[string]map[string]string `json:"entityBranchesByMomentIdx"`
	VictimCaseVictimContactICode string                                  `json:"victimContactIcode"`
	VictimCaseNewUser            security_daos.GeneralUserDTO            `json:"newUser"`
	VictimCaseDepartment         security_daos.DepartmentDTO             `json:"department"`
	VictimCaseCity               security_daos.CityDTO                   `json:"city"`
	VictimCaseTown               security_daos.TownDTO                   `json:"town"`
	VictimCaseMoments            []MomentDTO                             `json:"moments"`
	VictimCaseTownLatitude       float64                                 `json:"townLatitude"`
	VictimCaseTownLongitude      float64                                 `json:"townLongitude"`
	VictimCaseAttended           string                                  `json:"attended"`
	VictimCaseFormName           string                                  `json:"-"`
	VictimCaseForm1              VictimCaseForm1DTO                      `json:"form"`
	VictimCaseForm2              VictimCaseForm2DTO                      `json:"form2"`
}
type VictimCasePgDB struct {
	VictimCaseId           sql.NullInt64
	VictimCaseICode        sql.NullString
	VictimCaseCreationDate sql.NullTime
	VictimCaseUpdateDate   sql.NullTime
	VictimCaseStatus       sql.NullString
	VictimCaseGeneralUser  sql.NullString

	VictimCaseTownCode sql.NullString

	VictimCaseVictimContact sql.NullInt64
	VictimCaseApprovedBy    sql.NullInt64
	VictimCaseFollowUp      sql.NullInt64

	VictimCaseOwnerDescription sql.NullString

	VictimCaseNames     sql.NullString
	VictimCaseLastNames sql.NullString
	VictimCaseDocType   sql.NullString
	VictimCaseDocNumber sql.NullString

	VictimCaseAttended sql.NullString

	VictimCaseTownLatitude  sql.NullFloat64
	VictimCaseTownLongitude sql.NullFloat64
	TownName                sql.NullString
	CityName                sql.NullString
	DepartmentName          sql.NullString
}

func (vcd VictimCaseDTO) MarshalJSON() ([]byte, error) {
	type Alias VictimCaseDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		VictimCaseCreationDate string `json:"creationDate"`
		VictimCaseUpdateDate   string `json:"updateDate"`
	}{
		Alias:                  (*Alias)(&vcd),
		VictimCaseCreationDate: vcd.VictimCaseCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimCaseUpdateDate:   vcd.VictimCaseUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (vcd *VictimCaseDTO) UnmarshalJSON(data []byte) error {
	type Alias VictimCaseDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		VictimCaseCreationDate string `json:"creationDate"`
		VictimCaseUpdateDate   string `json:"updateDate"`
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

	// Parseo seguro
	vcd.VictimCaseCreationDate = parse(aux.VictimCaseCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	vcd.VictimCaseUpdateDate = parse(aux.VictimCaseUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

func SetVictimCase(victimCase *VictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var victimContactId interface{} = victimCase.VictimCaseVictimContact.VictimContactId
	if victimCase.VictimCaseVictimContact.VictimContactId == 0 {
		victimContactId = nil
	}

	var approvedBy interface{} = victimCase.VictimCaseApprovedBy.CaseOwnerId
	if victimCase.VictimCaseApprovedBy.CaseOwnerId == 0 {
		approvedBy = nil
	}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Query

	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy"}

	var victimCaseFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{"VictimCaseId"}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		victimCase.VictimCaseICode, victimCase.VictimCaseCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), victimCase.VictimCaseUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCase.VictimCaseStatus, victimCase.VictimCaseGeneralUser, victimCase.VictimCaseTownCode, victimCase.VictimCaseNames, victimCase.VictimCaseLastNames, victimCase.VictimCaseDocType,
		victimCase.VictimCaseDocNumber, victimContactId, approvedBy)

	persistenceCtrl.Scan(&victimCase.VictimCaseId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetVictimCase(by common_controllers.By, victimCase *VictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseId", "VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}

	var victimCaseFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	//		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var victimCasePg VictimCasePgDB = VictimCasePgDB{}

	persistenceCtrl.Scan(&victimCasePg.VictimCaseId, &victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
		&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
		&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp)

	*victimCase = victimCasePg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetVictimCases(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + victimCaseFieldsStr +
		` FROM ` + victimCasePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseDBScheme, VictimCaseFieldDefinitions, true) +
		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, by.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesByDepartmentICode(departmentICode string, victimCaseStatus string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var townPath string = security_daos.TownDBScheme + "." + security_daos.TownDBName
	var cityPath string = security_daos.CityDBScheme + "." + security_daos.CityDBName
	var departmentPath string = security_daos.DepartmentDBScheme + "." + security_daos.DepartmentDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + victimCaseFieldsStr +
		` FROM ` + victimCasePath +
		` LEFT JOIN ` + townPath + ` ON (` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + `)` +
		` LEFT JOIN ` + cityPath + ` ON (` + cityPath + `.` + security_daos.CityFieldDefinitions["CityId"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCity"].DBName + `)` +
		` LEFT JOIN ` + departmentPath + ` ON (` + departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentId"].DBName + ` = ` + cityPath + `.` + security_daos.CityFieldDefinitions["CityDepartment"].DBName + `)` +
		` WHERE  ` + departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentICode"].DBName + ` = $1 AND ` +
		victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $2 ` +
		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, departmentICode, victimCaseStatus)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			` LEFT JOIN ` + townPath + ` ON (` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + `)` +
			` LEFT JOIN ` + cityPath + ` ON (` + cityPath + `.` + security_daos.CityFieldDefinitions["CityId"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCity"].DBName + `)` +
			` LEFT JOIN ` + departmentPath + ` ON (` + departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentId"].DBName + ` = ` + cityPath + `.` + security_daos.CityFieldDefinitions["CityDepartment"].DBName + `)` +
			` WHERE  ` + departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentICode"].DBName + ` = $1  AND ` +
			victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $2 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, departmentICode, victimCaseStatus)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesByDepartmentICodeAndFollowUpStatusExcluded(departmentICode string, followUpStatusExcluded string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName
	var relPath string = RelBarrierFollowUpEntryDBScheme + "." + RelBarrierFollowUpEntryDBName
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var followUpEntryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName

	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var townPath string = security_daos.TownDBScheme + "." + security_daos.TownDBName
	var cityPath string = security_daos.CityDBScheme + "." + security_daos.CityDBName
	var departmentPath string = security_daos.DepartmentDBScheme + "." + security_daos.DepartmentDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}
	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + victimCaseFieldsStr +
		` FROM ` + victimCasePath +
		` RIGHT JOIN ` + followUpPath + ` ON (` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseFollowUp"].DBName + ` = ` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + `)` +
		` LEFT JOIN ` + followUpEntryPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryFollowUp"].DBName + ` = ` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + `)` +
		` LEFT JOIN ` + relPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryFollowUpEntry"].DBName + `)` +
		` LEFT JOIN ` + barrierPath + ` ON (` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryBarrier"].DBName + `)` +
		` LEFT JOIN ` + townPath + ` ON (` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + `)` +
		` LEFT JOIN ` + cityPath + ` ON (` + cityPath + `.` + security_daos.CityFieldDefinitions["CityId"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCity"].DBName + `)` +
		` LEFT JOIN ` + departmentPath + ` ON (` + departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentId"].DBName + ` = ` + cityPath + `.` + security_daos.CityFieldDefinitions["CityDepartment"].DBName + `)` +

		` WHERE  ` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` IS NOT NULL AND ` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryStatus"].DBName + ` != $1 AND ` +
		departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentICode"].DBName + ` = $2 ` +
		` GROUP BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName +
		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, followUpStatusExcluded, departmentICode)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			` RIGHT JOIN ` + followUpPath + ` ON (` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseFollowUp"].DBName + ` = ` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + `)` +
			` LEFT JOIN ` + followUpEntryPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryFollowUp"].DBName + ` = ` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + `)` +
			` LEFT JOIN ` + relPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryFollowUpEntry"].DBName + `)` +
			` LEFT JOIN ` + barrierPath + ` ON (` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryBarrier"].DBName + `)` +
			` LEFT JOIN ` + townPath + ` ON (` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + `)` +
			` LEFT JOIN ` + cityPath + ` ON (` + cityPath + `.` + security_daos.CityFieldDefinitions["CityId"].DBName + ` = ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCity"].DBName + `)` +
			` LEFT JOIN ` + departmentPath + ` ON (` + departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentId"].DBName + ` = ` + cityPath + `.` + security_daos.CityFieldDefinitions["CityDepartment"].DBName + `)` +

			` WHERE  ` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` IS NOT NULL AND ` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpStatus"].DBName + ` != $1 AND ` +
			departmentPath + `.` + security_daos.DepartmentFieldDefinitions["DepartmentICode"].DBName + ` = $2 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, followUpStatusExcluded, departmentICode)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesByOwnerUserICode(userICode string, victimCaseStatus string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName
	var relCaseOwnerVictimCasePath string = RelCaseOwnerVictimCaseDBScheme + "." + RelCaseOwnerVictimCaseDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + victimCaseFieldsStr +
		` FROM ` + victimCasePath +
		` LEFT JOIN ` + relCaseOwnerVictimCasePath + ` ON (` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
		` INNER JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CaseOwner"].DBName + `) ` +

		` WHERE ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = $1 AND ` +
		relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_Status"].DBName + ` = 'a'  AND ` +

		victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $2 ` +

		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, userICode, victimCaseStatus)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			` LEFT JOIN ` + relCaseOwnerVictimCasePath + ` ON (` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
			` INNER JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CaseOwner"].DBName + `) ` +

			` WHERE ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = $1 AND ` +
			relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_Status"].DBName + ` = 'a'  AND ` +

			victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $2 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, userICode, victimCaseStatus)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesByDocumentAndTownCodeWithAttend(docType string, docNumber string, townCode string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var momentPath string = MomentDBScheme + "." + MomentDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var commonWhere string = ` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseDocType"].DBName + ` = $1 AND ` +
		victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseDocNumber"].DBName + ` = $2 AND ` +
		victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = $3 `

	var subQuery string = `SELECT
	CASE
		WHEN COUNT(*) = SUM(CASE WHEN ` + momentPath + `.` + MomentFieldDefinitions["MomentStatus"].DBName + ` = 'a' THEN 1 ELSE 0 END)
			THEN CASE
				WHEN COUNT(*) = 0 THEN 'n'
				ELSE 'y'
			END
		ELSE 'p'
	END AS attended
	FROM ` + momentPath + `
		LEFT JOIN ` + victimCasePath + ` ON ( ` + momentPath + `.` + MomentFieldDefinitions["MomentVictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		commonWhere

	var query string = `SELECT ` + victimCaseFieldsStr + `, ` + subQuery + `, ` +
		` FROM ` + victimCasePath +

		commonWhere +

		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, docType, docNumber, townCode)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp, &victimCasePg.VictimCaseAttended)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			commonWhere

		persistenceCtrl.QueryRow(context.Background(), countQuery, docType, docNumber, townCode)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesByTownCodeWithAttend(townCode string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var momentPath string = MomentDBScheme + "." + MomentDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var commonWhere string = ` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = $1 `

	var subQuery string = `(SELECT CASE WHEN COUNT(*) = SUM(CASE WHEN ` + momentPath + `.` + MomentFieldDefinitions["MomentStatus"].DBName + ` = 'a' THEN 1 ELSE 0 END)` +
		`THEN CASE ` +
		`WHEN COUNT(*) = 0 THEN 'n' ` +
		`ELSE 'y' END ` +
		`ELSE 'p' END AS attended ` +
		`FROM ` + momentPath + `
		LEFT JOIN ` + victimCasePath + ` ON ( ` + momentPath + `.` + MomentFieldDefinitions["MomentVictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		commonWhere + `)`

	var query string = `SELECT ` + victimCaseFieldsStr + `, ` + subQuery + `, ` +
		` FROM ` + victimCasePath +

		commonWhere +

		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, townCode)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp, &victimCasePg.VictimCaseAttended)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			commonWhere

		persistenceCtrl.QueryRow(context.Background(), countQuery, townCode)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesByTownCodeAndEntityBranchWithAttend(townCode string, entityBranchICode string, victimCaseStatus string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var momentPath string = MomentDBScheme + "." + MomentDBName
	var entBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var commonWhere string = ` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + ` = $1 AND ` +
		entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchICode"].DBName + ` = $2 `

	var subQuery string = `(SELECT CASE WHEN COUNT(*) = SUM(CASE WHEN ` + momentPath + `.` + MomentFieldDefinitions["MomentStatus"].DBName + ` = 'a' THEN 1 ELSE 0 END)` +
		`THEN CASE ` +
		`WHEN COUNT(*) = 0 THEN 'n' ` +
		`ELSE 'y' END ` +
		`ELSE 'p' END AS attended ` +
		`FROM ` + momentPath + ` ` +
		`LEFT JOIN ` + victimCasePath + ` ON ( ` + momentPath + `.` + MomentFieldDefinitions["MomentVictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
		`LEFT JOIN ` + entBranchPath + ` ON ( ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchId"].DBName + ` = ` + momentPath + `.` + MomentFieldDefinitions["MomentEntityBranch"].DBName + `) ` +
		commonWhere + `)`

	var query string = `SELECT ` + victimCaseFieldsStr + `, ` + subQuery +
		` FROM ` + victimCasePath + ` ` +
		` LEFT JOIN ` + momentPath + ` ON ( ` + momentPath + `.` + MomentFieldDefinitions["MomentVictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
		` LEFT JOIN ` + entBranchPath + ` ON ( ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchId"].DBName + ` = ` + momentPath + `.` + MomentFieldDefinitions["MomentEntityBranch"].DBName + `) ` +

		commonWhere + ` AND ` +

		victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $3 ` +

		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, townCode, entityBranchICode, victimCaseStatus)

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp, &victimCasePg.VictimCaseAttended)

		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath + ` ` +
			` LEFT JOIN ` + momentPath + ` ON ( ` + momentPath + `.` + MomentFieldDefinitions["MomentVictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
			` LEFT JOIN ` + entBranchPath + ` ON ( ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchId"].DBName + ` = ` + momentPath + `.` + MomentFieldDefinitions["MomentEntityBranch"].DBName + `) ` +
			commonWhere + ` AND ` +

			victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $3 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, townCode, entityBranchICode, victimCaseStatus)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func GetVictimCasesReport(report VictimCaseReportDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, error) {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var departmentPath string = security_daos.DepartmentDBScheme + "." + security_daos.DepartmentDBName
	var victimCaseForm1Path string = VictimCaseForm1DBScheme + "." + VictimCaseForm1DBName
	var victimCaseForm2Path string = VictimCaseForm2DBScheme + "." + VictimCaseForm2DBName
	var cityPath string = security_daos.CityDBScheme + "." + security_daos.CityDBName
	var townPath string = security_daos.TownDBScheme + "." + security_daos.TownDBName

	var victimCases []VictimCaseDTO

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseForm1FieldsSlice []string = []string{"VictimCaseForm1Id", "VictimCaseForm1ICode", "VictimCaseForm1CreationDate", "VictimCaseForm1UpdateDate",
		"VictimCaseForm1Nick", "VictimCaseForm1BirthDate", "VictimCaseForm1ViolenceTownCode", "VictimCaseForm1Address",
		"VictimCaseForm1LivingLatitude", "VictimCaseForm1LivingLongitude", "VictimCaseForm1Phone", "VictimCaseForm1Email", "VictimCaseForm1GenderIdentity", "VictimCaseForm1SexualOrientation", "VictimCaseForm1Origin",
		"VictimCaseForm1Occupation", "VictimCaseForm1OccupationOther", "VictimCaseForm1VictimEthnicGroup", "VictimCaseForm1VictimEthnicGroupOther", "VictimCaseForm1VictimContactNames",
		"VictimCaseForm1VictimContactPhone", "VictimCaseForm1VictimContactKinship", "VictimCaseForm1VictimNumChildren", "VictimCaseForm1VictimMaritalStatus",
		"VictimCaseForm1VictimMaritalStatusOther", "VictimCaseForm1VictimChildrenAge", "VictimCaseForm1VictimDisability", "VictimCaseForm1VictimSpecialSupport", "VictimCaseForm1FactsOccurrence",
		"VictimCaseForm1FactsStartTime", "VictimCaseForm1FactsEndTime", "VictimCaseForm1FactsWeekday", "VictimCaseForm1FactsDate", "VictimCaseForm1FactsDescription",
		"VictimCaseForm1VictimViolenceExperienced", "VictimCaseForm1VictimViolenceExperiencedOther", "VictimCaseForm1VictimViolenceScope", "VictimCaseForm1VictimFemicideRisk",
		"VictimCaseForm1VictimAggressor", "VictimCaseForm1VictimRelationshipWithAggressor", "VictimCaseForm1VictimAggressorName", "VictimCaseForm1VictimAggressorDocType",
		"VictimCaseForm1VictimAggressorDocNumber", "VictimCaseForm1VictimAggressorAddress", "VictimCaseForm1VictimAggressorPhone",
		"VictimCaseForm1Age", "VictimCaseForm1VictimNationality", "VictimCaseForm1VictimNationalityOther", "VictimCaseForm1VictimForeignerImmigrationStatus",
		"VictimCaseForm1VictimGender", "VictimCaseForm1VictimGenderIdentityOther", "VictimCaseForm1VictimSexualOrientationOther", "VictimCaseForm1VictimDependents",
		"VictimCaseForm1VictimDeathThreats", "VictimCaseForm1VictimAggressorHasWeapons", "VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore",
		"VictimCaseForm1VictimImminentRisk", "VictimCaseForm1VictimPreviouslyReportedSituation", "VictimCaseForm1VictimIfPreviouslyReported",
		"VictimCaseForm1VictimIfAfro", "VictimCaseForm1VictimIfIndigenous", "VictimCaseForm1VictimIfIndigenousTongue", "VictimCaseForm1VictimIfPeasant",
		"VictimCaseForm1VictimIfArmedConflict", "VictimCaseForm1VictimViolenceScene",
		"VictimCaseForm1PhysicalViolenceIncreased", "VictimCaseForm1SeparatedFromPartnerLastYear", "VictimCaseForm1ThreatenedWithWeapon", "VictimCaseForm1ThreatenedToKillOrHarmChildren",
		"VictimCaseForm1JealousAndViolent", "VictimCaseForm1BelievesCapableOfKilling", "VictimCaseForm1VictimCase"}

	var victimCaseForm1FieldsAliasSlice []string = []string{}

	var victimCaseForm2FieldsSlice []string = []string{"VictimCaseForm2Id", "VictimCaseForm2ICode", "VictimCaseForm2CreationDate", "VictimCaseForm2UpdateDate", "VictimCaseForm2IdentityName",
		"VictimCaseForm2VictimPhone", "VictimCaseForm2FactsDescription", "VictimCaseForm2FactsDate", "VictimCaseForm2FactsStartTime", "VictimCaseForm2FactsTownCode",
		"VictimCaseForm2FactsZone", "VictimCaseForm2FactsAddress", "VictimCaseForm2ScenarioViolence", "VictimCaseForm2ReportedPreviously", "VictimCaseForm2RecurrenceAggression",
		"VictimCaseForm2NumAgressors", "VictimCaseForm2ProximityPrincipalAggressor", "VictimCaseForm2RelationshipWithPresumedAggressor", "VictimCaseForm2EconomicallyDependent", "VictimCaseForm2AggressorGenderIdentity",
		"VictimCaseForm2AggressorNames", "VictimCaseForm2AggressorDocType", "VictimCaseForm2AggressorDocNumber", "VictimCaseForm2AggressorAddress", "VictimCaseForm2AggressorPhone",
		"VictimCaseForm2AggressorViolencePhysicalIncrease", "VictimCaseForm2AggressorWeaponUsed", "VictimCaseForm2AggressorThreatKill", "VictimCaseForm2AggressorPursuesSpiesDestroys", "VictimCaseForm2AggressorCapableOfKilling",
		"VictimCaseForm2AggressorHasAccessToWeapons", "VictimCaseForm2PartnerUnemployed", "VictimCaseForm2PartnerOtherDenunciations", "VictimCaseForm2AggressorHasPenalBackground", "VictimCaseForm2AggressorForcedSex",
		"VictimCaseForm2AggressorAttemptedStrangulation", "VictimCaseForm2AggressorConsumesDrugs", "VictimCaseForm2AggressorIsAlcoholic", "VictimCaseForm2PartnerControls", "VictimCaseForm2AggressorHadHitInVulnerability",
		"VictimCaseForm2PartnerThreatenedSuicide", "VictimCaseForm2PartnerThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm", "VictimCaseForm2AggressorLimitsContactSupportNetworks", "VictimCaseForm2StillLivesWithAggressor",
		"VictimCaseForm2AggressorViolentlyJealous", "VictimCaseForm2AggressorUnemployed", "VictimCaseForm2AggressorHasPenalBackground2", "VictimCaseForm2AggressorSexuallyHarassment", "VictimCaseForm2AggressorUseDrugs",
		"VictimCaseForm2AggressorIsAlcoholic2", "VictimCaseForm2AggressorControls", "VictimCaseForm2AggressorThreatenedDamageMembers", "VictimCaseForm2ThoughtsOfSelfHarm2", "VictimCaseForm2AggressorCommonSpaces",
		"VictimCaseForm2AggressorHierarchy", "VictimCaseForm2RiskScore", "VictimCaseForm2RiskLevel", "VictimCaseForm2AggressorUsedPositionAuthority", "VictimCaseForm2BirthDate", "VictimCaseForm2PhysicalMentalSensoryDifficulties", "VictimCaseForm2Nationality", "VictimCaseForm2SpecifiedNationality",
		"VictimCaseForm2MigrationCondition", "VictimCaseForm2GenderIdentity", "VictimCaseForm2SexualOrientation", "VictimCaseForm2AssignedSexAtBirth", "VictimCaseForm2EthnicAffiliation",
		"VictimCaseForm2IndigenousPeople", "VictimCaseForm2CampesinoRecognition", "VictimCaseForm2MaritalStatus", "VictimCaseForm2LastEducationLevel", "VictimCaseForm2Occupation",
		"VictimCaseForm2IncomeGenerationMethod", "VictimCaseForm2ApproxStartAsp", "VictimCaseForm2HousingTenancyForm", "VictimCaseForm2HousingStratum", "VictimCaseForm2CurrentlyPregnant",
		"VictimCaseForm2ResidenceTownCode", "VictimCaseForm2ResidenceAddress", "VictimCaseForm2ResidenceZone", "VictimCaseForm2SupportContactNames", "VictimCaseForm2SupportContactPhone",
		"VictimCaseForm2SupportContactEmail", "VictimCaseForm2SupportContactKinship",
		"VictimCaseForm2PersonWithDisability", "VictimCaseForm2RequireLanguageInterpreter", "VictimCaseForm2WorkplaceSectorOccurrence", "VictimCaseForm2ViolenceMotivatedByGender", "VictimCaseForm2AttentionWasAppropriate", "VictimCaseForm2AggressorOccupation", "VictimCaseForm2SalivaManagementExplanation",
		"VictimCaseForm2ActivitiesUnableToHear", "VictimCaseForm2ActivitiesUnableToTalk", "VictimCaseForm2ActivitiesUnableToSee", "VictimCaseForm2ActivitiesUnableToMove", "VictimCaseForm2ActivitiesUnableToTake",
		"VictimCaseForm2ActivitiesUnableToUnderstand", "VictimCaseForm2ActivitiesUnableToEat", "VictimCaseForm2ActivitiesUnableToInteract", "VictimCaseForm2ActivitiesUnableToDoEveryday",
		"VictimCaseForm2AllowsEasyReport", "VictimCaseForm2VictimCase"}

	var victimCaseForm2FieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var townFieldsSlice []string = []string{"TownLatitude", "TownLongitude", "TownName"}
	var townFieldsAliasSlice []string = []string{}

	var cityFieldsSlice []string = []string{"CityName"}
	var cityFieldsAliasSlice []string = []string{}

	var departmentFieldsSlice []string = []string{"DepartmentName"}
	var departmentAliasSlice []string = []string{}

	var townFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, townFieldsSlice, townFieldsAliasSlice, security_daos.TownDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.TownDBScheme, security_daos.TownFieldDefinitions, true)
	var cityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, cityFieldsSlice, cityFieldsAliasSlice, security_daos.CityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.CityDBScheme, security_daos.CityFieldDefinitions, true)
	var departmentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, departmentFieldsSlice, departmentAliasSlice, security_daos.DepartmentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.DepartmentDBScheme, security_daos.DepartmentFieldDefinitions, true)
	var victimCaseForm1FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseForm1FieldsSlice, victimCaseForm1FieldsAliasSlice, VictimCaseForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm1DBScheme, VictimCaseForm1FieldDefinitions, true)
	var victimCaseForm2FieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseForm2FieldsSlice, victimCaseForm2FieldsAliasSlice, VictimCaseForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2DBScheme, VictimCaseForm2FieldDefinitions, true)

	var whereWithDates string = victimCasePath + "." + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` BETWEEN $1 AND $2 `
	var whereWithDepartment string = departmentPath + "." + security_daos.DepartmentFieldDefinitions["DepartmentICode"].DBName + " = '" + report.VictimCaseReportDepartment + "' "
	var whereWithCity string = cityPath + "." + security_daos.CityFieldDefinitions["CityCode"].DBName + " = '" + report.VictimCaseReportCity + "' "
	var whereWithTown string = townPath + "." + security_daos.TownFieldDefinitions["TownCode"].DBName + " = '" + report.VictimCaseReportTown + "' "
	var whereWithViolence string = victimCasePath + "." + VictimCaseFieldDefinitions["VictimCaseVictimViolenceExperienced"].DBName + " = '" + report.VictimCaseReportViolenceType + "' "
	var whereWithStatus string = victimCasePath + "." + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + " = '" + report.VictimCaseReportCaseStatus + "' "

	var finalWhere string = whereWithDates

	if report.VictimCaseReportDepartment != "" {
		finalWhere += " AND " + whereWithDepartment
	}

	if report.VictimCaseReportCity != "" {
		finalWhere += " AND " + whereWithCity
	}

	if report.VictimCaseReportTown != "" {
		finalWhere += " AND " + whereWithTown
	}

	if report.VictimCaseReportViolenceType != "" {
		finalWhere += " AND " + whereWithViolence
	}

	if report.VictimCaseReportCaseStatus != "" {
		finalWhere += " AND " + whereWithStatus
	}

	var query string = `SELECT ` + victimCaseFieldsStr + " , " + victimCaseForm1FieldsStr + " , " + victimCaseForm2FieldsStr + " , " + townFieldsStr + `, ` + cityFieldsStr + `, ` + departmentFieldsStr +
		` FROM ` + victimCasePath +
		` LEFT JOIN ` + victimCaseForm1Path + ` ON (` + victimCaseForm1Path + "." + VictimCaseForm1FieldDefinitions["VictimCaseForm1VictimCase"].DBName + ` = ` + victimCasePath + "." + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
		` LEFT JOIN ` + victimCaseForm2Path + ` ON (` + victimCaseForm2Path + "." + VictimCaseForm2FieldDefinitions["VictimCaseForm2VictimCase"].DBName + ` = ` + victimCasePath + "." + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
		` LEFT JOIN ` + townPath + ` ON (` + townPath + "." + security_daos.TownFieldDefinitions["TownCode"].DBName + ` = ` + victimCasePath + "." + VictimCaseFieldDefinitions["VictimCaseTownCode"].DBName + `) ` +
		` LEFT JOIN ` + cityPath + ` ON (` + cityPath + "." + security_daos.CityFieldDefinitions["CityId"].DBName + ` = ` + townPath + "." + security_daos.TownFieldDefinitions["TownCity"].DBName + `) ` +
		` LEFT JOIN ` + departmentPath + ` ON (` + cityPath + "." + security_daos.CityFieldDefinitions["CityDepartment"].DBName + ` = ` + departmentPath + "." + security_daos.DepartmentFieldDefinitions["DepartmentId"].DBName + `) ` +

		` WHERE ` + townPath + "." + security_daos.TownFieldDefinitions["TownLatitude"].DBName + ` IS NOT NULL AND ` +
		townPath + "." + security_daos.TownFieldDefinitions["TownLongitude"].DBName + ` IS NOT NULL AND ` +
		finalWhere +

		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  `

	var endDate time.Time = time.Date(report.VictimCaseReportEndDate.Year(), report.VictimCaseReportEndDate.Month(), report.VictimCaseReportEndDate.Day(), 23, 59, 59, 0, report.VictimCaseReportEndDate.Location())

	persistenceCtrl.Query(context.Background(), query, report.VictimCaseReportStartDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), endDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT))

	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}
		var victimCaseForm1Pg VictimCaseForm1PgDB = VictimCaseForm1PgDB{}
		var victimCaseForm2Pg VictimCaseForm2PgDB = VictimCaseForm2PgDB{}

		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp,

			&victimCaseForm1Pg.VictimCaseForm1Id,
			&victimCaseForm1Pg.VictimCaseForm1ICode, &victimCaseForm1Pg.VictimCaseForm1CreationDate, &victimCaseForm1Pg.VictimCaseForm1UpdateDate,
			&victimCaseForm1Pg.VictimCaseForm1Nick, &victimCaseForm1Pg.VictimCaseForm1BirthDate,
			&victimCaseForm1Pg.VictimCaseForm1ViolenceTownCode, &victimCaseForm1Pg.VictimCaseForm1Address,
			&victimCaseForm1Pg.VictimCaseForm1LivingLatitude, &victimCaseForm1Pg.VictimCaseForm1LivingLongitude, &victimCaseForm1Pg.VictimCaseForm1Phone, &victimCaseForm1Pg.VictimCaseForm1Email, &victimCaseForm1Pg.VictimCaseForm1GenderIdentity, &victimCaseForm1Pg.VictimCaseForm1SexualOrientation,
			&victimCaseForm1Pg.VictimCaseForm1Origin, &victimCaseForm1Pg.VictimCaseForm1Occupation, &victimCaseForm1Pg.VictimCaseForm1OccupationOther, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroup, &victimCaseForm1Pg.VictimCaseForm1VictimEthnicGroupOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimContactNames, &victimCaseForm1Pg.VictimCaseForm1VictimContactPhone, &victimCaseForm1Pg.VictimCaseForm1VictimContactKinship, &victimCaseForm1Pg.VictimCaseForm1VictimNumChildren,
			&victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatus, &victimCaseForm1Pg.VictimCaseForm1VictimMaritalStatusOther, &victimCaseForm1Pg.VictimCaseForm1VictimChildrenAge, &victimCaseForm1Pg.VictimCaseForm1VictimDisability,
			&victimCaseForm1Pg.VictimCaseForm1VictimSpecialSupport, &victimCaseForm1Pg.VictimCaseForm1FactsOccurrence, &victimCaseForm1Pg.VictimCaseForm1FactsStartTime,
			&victimCaseForm1Pg.VictimCaseForm1FactsEndTime, &victimCaseForm1Pg.VictimCaseForm1FactsWeekday,
			&victimCaseForm1Pg.VictimCaseForm1FactsDate, &victimCaseForm1Pg.VictimCaseForm1FactsDescription, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperienced,
			&victimCaseForm1Pg.VictimCaseForm1VictimViolenceExperiencedOther, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScope, &victimCaseForm1Pg.VictimCaseForm1VictimFemicideRisk,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimRelationshipWithAggressor, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorName,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocType, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorDocNumber, &victimCaseForm1Pg.VictimCaseForm1VictimAggressorAddress,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorPhone, &victimCaseForm1Pg.VictimCaseForm1Age, &victimCaseForm1Pg.VictimCaseForm1VictimNationality, &victimCaseForm1Pg.VictimCaseForm1VictimNationalityOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimForeignerImmigrationStatus, &victimCaseForm1Pg.VictimCaseForm1VictimGender, &victimCaseForm1Pg.VictimCaseForm1VictimGenderIdentityOther,
			&victimCaseForm1Pg.VictimCaseForm1VictimSexualOrientationOther, &victimCaseForm1Pg.VictimCaseForm1VictimDependents, &victimCaseForm1Pg.VictimCaseForm1VictimDeathThreats,
			&victimCaseForm1Pg.VictimCaseForm1VictimAggressorHasWeapons, &victimCaseForm1Pg.VictimCaseForm1ExperiencedPhysicalOrSexualViolenceBefore,
			&victimCaseForm1Pg.VictimCaseForm1VictimImminentRisk, &victimCaseForm1Pg.VictimCaseForm1VictimPreviouslyReportedSituation, &victimCaseForm1Pg.VictimCaseForm1VictimIfPreviouslyReported,
			&victimCaseForm1Pg.VictimCaseForm1VictimIfAfro, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenous, &victimCaseForm1Pg.VictimCaseForm1VictimIfIndigenousTongue,
			&victimCaseForm1Pg.VictimCaseForm1VictimIfPeasant, &victimCaseForm1Pg.VictimCaseForm1VictimIfArmedConflict, &victimCaseForm1Pg.VictimCaseForm1VictimViolenceScene,
			&victimCaseForm1Pg.VictimCaseForm1PhysicalViolenceIncreased, &victimCaseForm1Pg.VictimCaseForm1SeparatedFromPartnerLastYear, &victimCaseForm1Pg.VictimCaseForm1ThreatenedWithWeapon,
			&victimCaseForm1Pg.VictimCaseForm1ThreatenedToKillOrHarmChildren, &victimCaseForm1Pg.VictimCaseForm1JealousAndViolent, &victimCaseForm1Pg.VictimCaseForm1BelievesCapableOfKilling, &victimCaseForm1Pg.VictimCaseForm1VictimCase,

			&victimCaseForm2Pg.VictimCaseForm2Id, &victimCaseForm2Pg.VictimCaseForm2ICode, &victimCaseForm2Pg.VictimCaseForm2CreationDate, &victimCaseForm2Pg.VictimCaseForm2UpdateDate,
			&victimCaseForm2Pg.VictimCaseForm2IdentityName, &victimCaseForm2Pg.VictimCaseForm2VictimPhone, &victimCaseForm2Pg.VictimCaseForm2FactsDescription, &victimCaseForm2Pg.VictimCaseForm2FactsDate,
			&victimCaseForm2Pg.VictimCaseForm2FactsStartTime, &victimCaseForm2Pg.VictimCaseForm2FactsTownCode, &victimCaseForm2Pg.VictimCaseForm2FactsZone, &victimCaseForm2Pg.VictimCaseForm2FactsAddress,
			&victimCaseForm2Pg.VictimCaseForm2ScenarioViolence, &victimCaseForm2Pg.VictimCaseForm2ReportedPreviously, &victimCaseForm2Pg.VictimCaseForm2RecurrenceAggression, &victimCaseForm2Pg.VictimCaseForm2NumAgressors,
			&victimCaseForm2Pg.VictimCaseForm2ProximityPrincipalAggressor, &victimCaseForm2Pg.VictimCaseForm2RelationshipWithPresumedAggressor, &victimCaseForm2Pg.VictimCaseForm2EconomicallyDependent, &victimCaseForm2Pg.VictimCaseForm2AggressorGenderIdentity,
			&victimCaseForm2Pg.VictimCaseForm2AggressorNames, &victimCaseForm2Pg.VictimCaseForm2AggressorDocType, &victimCaseForm2Pg.VictimCaseForm2AggressorDocNumber, &victimCaseForm2Pg.VictimCaseForm2AggressorAddress,
			&victimCaseForm2Pg.VictimCaseForm2AggressorPhone, &victimCaseForm2Pg.VictimCaseForm2AggressorViolencePhysicalIncrease, &victimCaseForm2Pg.VictimCaseForm2AggressorWeaponUsed, &victimCaseForm2Pg.VictimCaseForm2AggressorThreatKill,
			&victimCaseForm2Pg.VictimCaseForm2AggressorPursuesSpiesDestroys, &victimCaseForm2Pg.VictimCaseForm2AggressorCapableOfKilling, &victimCaseForm2Pg.VictimCaseForm2AggressorHasAccessToWeapons, &victimCaseForm2Pg.VictimCaseForm2PartnerUnemployed,
			&victimCaseForm2Pg.VictimCaseForm2PartnerOtherDenunciations, &victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground, &victimCaseForm2Pg.VictimCaseForm2AggressorForcedSex, &victimCaseForm2Pg.VictimCaseForm2AggressorAttemptedStrangulation,
			&victimCaseForm2Pg.VictimCaseForm2AggressorConsumesDrugs, &victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic, &victimCaseForm2Pg.VictimCaseForm2PartnerControls, &victimCaseForm2Pg.VictimCaseForm2AggressorHadHitInVulnerability,
			&victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedSuicide, &victimCaseForm2Pg.VictimCaseForm2PartnerThreatenedDamageMembers, &victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm, &victimCaseForm2Pg.VictimCaseForm2AggressorLimitsContactSupportNetworks,
			&victimCaseForm2Pg.VictimCaseForm2StillLivesWithAggressor, &victimCaseForm2Pg.VictimCaseForm2AggressorViolentlyJealous, &victimCaseForm2Pg.VictimCaseForm2AggressorUnemployed, &victimCaseForm2Pg.VictimCaseForm2AggressorHasPenalBackground2,
			&victimCaseForm2Pg.VictimCaseForm2AggressorSexuallyHarassment, &victimCaseForm2Pg.VictimCaseForm2AggressorUseDrugs, &victimCaseForm2Pg.VictimCaseForm2AggressorIsAlcoholic2, &victimCaseForm2Pg.VictimCaseForm2AggressorControls,
			&victimCaseForm2Pg.VictimCaseForm2AggressorThreatenedDamageMembers, &victimCaseForm2Pg.VictimCaseForm2ThoughtsOfSelfHarm2, &victimCaseForm2Pg.VictimCaseForm2AggressorCommonSpaces, &victimCaseForm2Pg.VictimCaseForm2AggressorHierarchy,
			&victimCaseForm2Pg.VictimCaseForm2RiskScore, &victimCaseForm2Pg.VictimCaseForm2RiskLevel, &victimCaseForm2Pg.VictimCaseForm2AggressorUsedPositionAuthority,
			&victimCaseForm2Pg.VictimCaseForm2BirthDate, &victimCaseForm2Pg.VictimCaseForm2PhysicalMentalSensoryDifficulties, &victimCaseForm2Pg.VictimCaseForm2Nationality, &victimCaseForm2Pg.VictimCaseForm2SpecifiedNationality,
			&victimCaseForm2Pg.VictimCaseForm2MigrationCondition, &victimCaseForm2Pg.VictimCaseForm2GenderIdentity, &victimCaseForm2Pg.VictimCaseForm2SexualOrientation, &victimCaseForm2Pg.VictimCaseForm2AssignedSexAtBirth,
			&victimCaseForm2Pg.VictimCaseForm2EthnicAffiliation, &victimCaseForm2Pg.VictimCaseForm2IndigenousPeople, &victimCaseForm2Pg.VictimCaseForm2CampesinoRecognition, &victimCaseForm2Pg.VictimCaseForm2MaritalStatus,
			&victimCaseForm2Pg.VictimCaseForm2LastEducationLevel, &victimCaseForm2Pg.VictimCaseForm2Occupation, &victimCaseForm2Pg.VictimCaseForm2IncomeGenerationMethod, &victimCaseForm2Pg.VictimCaseForm2ApproxStartAsp,
			&victimCaseForm2Pg.VictimCaseForm2HousingTenancyForm, &victimCaseForm2Pg.VictimCaseForm2HousingStratum, &victimCaseForm2Pg.VictimCaseForm2CurrentlyPregnant, &victimCaseForm2Pg.VictimCaseForm2ResidenceTownCode,
			&victimCaseForm2Pg.VictimCaseForm2ResidenceAddress, &victimCaseForm2Pg.VictimCaseForm2ResidenceZone, &victimCaseForm2Pg.VictimCaseForm2SupportContactNames, &victimCaseForm2Pg.VictimCaseForm2SupportContactPhone,
			&victimCaseForm2Pg.VictimCaseForm2SupportContactEmail, &victimCaseForm2Pg.VictimCaseForm2SupportContactKinship, &victimCaseForm2Pg.VictimCaseForm2PersonWithDisability, &victimCaseForm2Pg.VictimCaseForm2RequireLanguageInterpreter,
			&victimCaseForm2Pg.VictimCaseForm2WorkplaceSectorOccurrence, &victimCaseForm2Pg.VictimCaseForm2ViolenceMotivatedByGender, &victimCaseForm2Pg.VictimCaseForm2AttentionWasAppropriate, &victimCaseForm2Pg.VictimCaseForm2AggressorOccupation,
			&victimCaseForm2Pg.VictimCaseForm2SalivaManagementExplanation, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToHear, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTalk, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToSee,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToMove, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToTake, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToUnderstand, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToEat,
			&victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToInteract, &victimCaseForm2Pg.VictimCaseForm2ActivitiesUnableToDoEveryday,
			&victimCaseForm2Pg.VictimCaseForm2AllowsEasyReport, &victimCaseForm2Pg.VictimCaseForm2VictimCase,

			&victimCasePg.VictimCaseTownLatitude, &victimCasePg.VictimCaseTownLongitude, &victimCasePg.TownName,
			&victimCasePg.CityName, &victimCasePg.DepartmentName)

		var victimCase VictimCaseDTO = victimCasePg.ToDTOTranslated()

		victimCase.VictimCaseForm1 = victimCaseForm1Pg.ToDTO()
		victimCase.VictimCaseForm2 = victimCaseForm2Pg.ToDTO()

		victimCases = append(victimCases, victimCase)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	//TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return victimCases, nil
}

func GetAllVictimCases(victimCaseStatus string, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseDTO, int, error) {
	// Definición de variables
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseICode", "VictimCaseCreationDate", "VictimCaseUpdateDate", "VictimCaseStatus", "VictimCaseGeneralUser", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber", "VictimCaseVictimContact", "VictimCaseApprovedBy", "VictimCaseOwnerDescription", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	/*
		var profileFieldsSlice []string = []string{"GeneralUserProfileNames", "GeneralUserProfileLastNames"}
		var profileFieldsAliasSlice []string = []string{}

		var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, security_daos.GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.GeneralUserProfileDBScheme, security_daos.GeneralUserProfileFieldDefinitions, true)
	*/

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + victimCaseFieldsStr +
		` FROM ` + victimCasePath +
		` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $1 ` +
		` ORDER BY ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	//var query string = `SELECT victimCase_id, victimCase_i_code, victimCase_creation_date,victimCase_updatdate,victimCase_data,victimCase_password,victimCase_status,victimCase_language FROM ` + VictimCaseDBScheme + `.victimCase`

	persistenceCtrl.Query(context.Background(), query, victimCaseStatus)
	var victimCases []VictimCaseDTO
	for persistenceCtrl.Next() {
		var victimCasePg VictimCasePgDB = VictimCasePgDB{}

		persistenceCtrl.ScanRow(&victimCasePg.VictimCaseICode, &victimCasePg.VictimCaseCreationDate, &victimCasePg.VictimCaseUpdateDate,
			&victimCasePg.VictimCaseStatus, &victimCasePg.VictimCaseGeneralUser, &victimCasePg.VictimCaseTownCode, &victimCasePg.VictimCaseNames, &victimCasePg.VictimCaseLastNames, &victimCasePg.VictimCaseDocType,
			&victimCasePg.VictimCaseDocNumber, &victimCasePg.VictimCaseVictimContact, &victimCasePg.VictimCaseApprovedBy, &victimCasePg.VictimCaseOwnerDescription, &victimCasePg.VictimCaseFollowUp)
		victimCases = append(victimCases, victimCasePg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + victimCasePath +
			` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseStatus"].DBName + ` = $1 `

		persistenceCtrl.QueryRow(context.Background(), countQuery, victimCaseStatus)
		persistenceCtrl.Scan(&count)
	}

	return victimCases, count, nil
}

func UpdateVictimCase(victimCase *VictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseUpdateDate", "VictimCaseTownCode", "VictimCaseNames", "VictimCaseLastNames",
		"VictimCaseDocType", "VictimCaseDocNumber"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{"VictimCaseId"}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, victimCase.VictimCaseId,
		victimCase.VictimCaseUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		victimCase.VictimCaseTownCode, victimCase.VictimCaseNames, victimCase.VictimCaseLastNames,
		victimCase.VictimCaseDocType, victimCase.VictimCaseDocNumber)

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

func UpdateVictimCaseFollowUp(victimCase *VictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseUpdateDate", "VictimCaseFollowUp"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{"VictimCaseId"}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, victimCase.VictimCaseId,
		victimCase.VictimCaseUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), victimCase.VictimCaseFollowUp.FollowUpId)

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

func UpdateVictimCaseStatus(victimCase *VictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var victimCaseFieldsSlice []string = []string{"VictimCaseStatus"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{"VictimCaseId"}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, victimCase.VictimCaseId, victimCase.VictimCaseStatus)

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

// Esta función no retorna nada porque está diseñada para ser ejecutada mediante TaskScheduler y sólo muestra errores en consola.
// Actualiza automaticamente la columna de operarios y roles para que quede quemada en la tabla y no haya que hacer la misma consulta cada vez que listan casos
func UpdateVictimCasesOwnersAndRoles(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var relCaseOwnerVictimCasePath string = RelCaseOwnerVictimCaseDBScheme + "." + RelCaseOwnerVictimCaseDBName
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName
	var profilePath string = security_daos.GeneralUserProfileDBScheme + "." + security_daos.GeneralUserProfileDBName

	var relRolePath string = security_daos.RelRoleGeneralUserDBScheme + "." + security_daos.RelRoleGeneralUserDBName
	var rolePath string = security_daos.RoleDBScheme + "." + security_daos.RoleDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return
	}

	// Query

	/*
		SELECT string_agg( '('||salvia.rel_case_owner_victim_case.rel_case_owner_victim_case_creation_date||') ' ||security.general_user_profile.general_user_profile_names || ' ' || security.general_user_profile.general_user_profile_last_names ||
			(
				SELECT string_agg( ' ['||security.role.role_name || '] ', ', ') as roles
				FROM security.role
				RIGHT JOIN security.rel_role_general_user ON (security.rel_role_general_user.role_id = security.role.role_id)
				LEFT JOIN security.general_user ON (security.rel_role_general_user.general_user_id = security.general_user.general_user_id)

				WHERE us.general_user_id = security.general_user.general_user_id

			), ', ' ORDER BY salvia.rel_case_owner_victim_case.rel_case_owner_victim_case_creation_date ASC) as case_owners
			FROM salvia.rel_case_owner_victim_case
			LEFT JOIN salvia.case_owner ON (salvia.case_owner.case_owner_id = salvia.rel_case_owner_victim_case.case_owner_id)
			LEFT JOIN security.general_user us ON (salvia.case_owner.case_owner_general_user = us.general_user_i_code)
			LEFT JOIN security.general_user_profile ON (security.general_user_profile.general_user_profile_id = us.general_user_general_user_profile)

			WHERE salvia.rel_case_owner_victim_case.victim_case_id = salvia.victim_case.victim_case_id
	*/

	var owners string = ` (SELECT string_agg( '('||` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CreationDate"].DBName + `||') ' ||` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileNames"].DBName + ` || ' ' || ` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileLastNames"].DBName + ` || ` +
		` (SELECT string_agg( ' ['||` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleName"].DBName + ` || '] ', ', ') as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` WHERE us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName +
		`), ', ' ORDER BY ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CreationDate"].DBName + ` ASC) as case_owners ` +

		` FROM ` + relCaseOwnerVictimCasePath +
		` LEFT JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CaseOwner"].DBName + `) ` +
		` LEFT JOIN ` + userPath + ` AS us ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `) ` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +

		` WHERE ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)`

	var query string = `UPDATE ` + victimCasePath +
		` SET ` + VictimCaseFieldDefinitions["VictimCaseOwnerDescription"].DBName + ` =  ` + owners + ` `

	persistenceCtrl.Exec(context.Background(), query)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
	}

}

// Igual que UpdateVictimCasesOwnersAndRoles pero se utiliza para actualizar un registro en particular. Actualiza automaticamente la columna de operarios y roles para que quede quemada en la tabla
func UpdateVictimCaseOwnersAndRolesByVictimCaseId(victimCaseId uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var relCaseOwnerVictimCasePath string = RelCaseOwnerVictimCaseDBScheme + "." + RelCaseOwnerVictimCaseDBName
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName
	var profilePath string = security_daos.GeneralUserProfileDBScheme + "." + security_daos.GeneralUserProfileDBName

	var relRolePath string = security_daos.RelRoleGeneralUserDBScheme + "." + security_daos.RelRoleGeneralUserDBName
	var rolePath string = security_daos.RoleDBScheme + "." + security_daos.RoleDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query

	/*
		SELECT string_agg( '('||salvia.rel_case_owner_victim_case.rel_case_owner_victim_case_creation_date||') ' ||security.general_user_profile.general_user_profile_names || ' ' || security.general_user_profile.general_user_profile_last_names ||
			(
				SELECT string_agg( ' ['||security.role.role_name || '] ', ', ') as roles
				FROM security.role
				RIGHT JOIN security.rel_role_general_user ON (security.rel_role_general_user.role_id = security.role.role_id)
				LEFT JOIN security.general_user ON (security.rel_role_general_user.general_user_id = security.general_user.general_user_id)

				WHERE us.general_user_id = security.general_user.general_user_id

			), ', ' ORDER BY salvia.rel_case_owner_victim_case.rel_case_owner_victim_case_creation_date ASC) as case_owners
			FROM salvia.rel_case_owner_victim_case
			LEFT JOIN salvia.case_owner ON (salvia.case_owner.case_owner_id = salvia.rel_case_owner_victim_case.case_owner_id)
			LEFT JOIN security.general_user us ON (salvia.case_owner.case_owner_general_user = us.general_user_i_code)
			LEFT JOIN security.general_user_profile ON (security.general_user_profile.general_user_profile_id = us.general_user_general_user_profile)

			WHERE salvia.rel_case_owner_victim_case.victim_case_id = salvia.victim_case.victim_case_id AND salvia.victim_case.victim_case_id = $1
	*/

	var owners string = ` (SELECT string_agg( '('||` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CreationDate"].DBName + `||') ' ||` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileNames"].DBName + ` || ' ' || ` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileLastNames"].DBName + ` || ` +
		` (SELECT string_agg( ' ['||` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleName"].DBName + ` || '] ', ', ') as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` WHERE us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName +
		`), ', ' ORDER BY ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CreationDate"].DBName + ` ASC) as case_owners ` +

		` FROM ` + relCaseOwnerVictimCasePath +
		` LEFT JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CaseOwner"].DBName + `) ` +
		` LEFT JOIN ` + userPath + ` AS us ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `) ` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +

		` WHERE ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + ` AND ` +
		victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + ` = $1) WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + ` = $1 `

	var query string = `UPDATE ` + victimCasePath +
		` SET ` + VictimCaseFieldDefinitions["VictimCaseOwnerDescription"].DBName + ` =  ` + owners + ` `

	persistenceCtrl.Exec(context.Background(), query, victimCaseId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
	}
	return nil
}

// Algunas utilidades
func SetVictimCaseDefaults(victimCase *VictimCaseDTO, action string, s utils.CommonSession) {

	switch action {
	case common_dao.SQL_INSERT:
		victimCase.VictimCaseCreationDate = time.Now()
		victimCase.VictimCaseUpdateDate = time.Now()
		victimCase.VictimCaseICode = utils.GetUUID()
		switch s.CurrentRole {
		case "op":
			victimCase.VictimCaseStatus = "ra"
		case "et":
			victimCase.VictimCaseStatus = "r"
		}

	case common_dao.SQL_UPDATE:
		victimCase.VictimCaseUpdateDate = time.Now()
	}

}

func (obj *VictimCasePgDB) ToDTO() VictimCaseDTO {
	var dto VictimCaseDTO

	if obj.VictimCaseId.Valid {

		dto.VictimCaseId = uint64(obj.VictimCaseId.Int64)
	}

	if obj.VictimCaseICode.Valid {

		dto.VictimCaseICode = obj.VictimCaseICode.String
	}

	if obj.VictimCaseCreationDate.Valid {

		dto.VictimCaseCreationDate = obj.VictimCaseCreationDate.Time
	}

	if obj.VictimCaseUpdateDate.Valid {

		dto.VictimCaseUpdateDate = obj.VictimCaseUpdateDate.Time
	}

	if obj.VictimCaseNames.Valid {
		dto.VictimCaseNames = obj.VictimCaseNames.String
	}

	if obj.VictimCaseStatus.Valid {
		dto.VictimCaseStatus = obj.VictimCaseStatus.String
	}

	if obj.VictimCaseGeneralUser.Valid {
		dto.VictimCaseGeneralUser = obj.VictimCaseGeneralUser.String
	}

	if obj.VictimCaseTownCode.Valid {
		dto.VictimCaseTownCode = obj.VictimCaseTownCode.String
	}

	if obj.VictimCaseLastNames.Valid {
		dto.VictimCaseLastNames = obj.VictimCaseLastNames.String
	}

	if obj.VictimCaseDocType.Valid {
		dto.VictimCaseDocType = obj.VictimCaseDocType.String
	}

	if obj.VictimCaseDocNumber.Valid {
		dto.VictimCaseDocNumber = obj.VictimCaseDocNumber.String
	}

	if obj.VictimCaseVictimContact.Valid {
		dto.VictimCaseVictimContact = VictimContactDTO{VictimContactId: uint64(obj.VictimCaseVictimContact.Int64)}
	}

	if obj.VictimCaseApprovedBy.Valid {
		dto.VictimCaseApprovedBy = CaseOwnerDTO{CaseOwnerId: uint64(obj.VictimCaseApprovedBy.Int64)}
	}

	if obj.VictimCaseFollowUp.Valid {
		dto.VictimCaseFollowUp = FollowUpDTO{FollowUpId: uint64(obj.VictimCaseFollowUp.Int64)}
	}

	if obj.VictimCaseOwnerDescription.Valid {
		dto.VictimCaseOwnerDescription = obj.VictimCaseOwnerDescription.String
	}

	if obj.VictimCaseAttended.Valid {
		dto.VictimCaseAttended = obj.VictimCaseAttended.String
	}

	if obj.VictimCaseTownLatitude.Valid {
		dto.VictimCaseTownLatitude = obj.VictimCaseTownLatitude.Float64
	}

	if obj.VictimCaseTownLongitude.Valid {
		dto.VictimCaseTownLongitude = obj.VictimCaseTownLongitude.Float64
	}

	return dto
}

func (obj *VictimCasePgDB) ToDTOTranslated() VictimCaseDTO {
	var dto VictimCaseDTO

	if obj.VictimCaseId.Valid {

		dto.VictimCaseId = uint64(obj.VictimCaseId.Int64)
	}

	if obj.VictimCaseICode.Valid {

		dto.VictimCaseICode = obj.VictimCaseICode.String
	}

	if obj.VictimCaseCreationDate.Valid {

		dto.VictimCaseCreationDate = obj.VictimCaseCreationDate.Time
	}

	if obj.VictimCaseUpdateDate.Valid {

		dto.VictimCaseUpdateDate = obj.VictimCaseUpdateDate.Time
	}

	if obj.VictimCaseNames.Valid {
		dto.VictimCaseNames = obj.VictimCaseNames.String
	}

	if obj.VictimCaseStatus.Valid {
		dto.VictimCaseStatus = obj.VictimCaseStatus.String
	}

	if obj.VictimCaseGeneralUser.Valid {
		dto.VictimCaseGeneralUser = obj.VictimCaseGeneralUser.String
	}

	if obj.VictimCaseLastNames.Valid {
		dto.VictimCaseLastNames = obj.VictimCaseLastNames.String
	}

	if obj.VictimCaseDocType.Valid {
		dto.VictimCaseDocType = common_config.DOCUMENT_TYPE[obj.VictimCaseDocType.String]
	}

	if obj.VictimCaseDocNumber.Valid {
		dto.VictimCaseDocNumber = obj.VictimCaseDocNumber.String
	}

	if obj.VictimCaseTownCode.Valid {
		dto.VictimCaseTownCode = obj.TownName.String
	}

	if obj.VictimCaseVictimContact.Valid {
		dto.VictimCaseVictimContact = VictimContactDTO{VictimContactId: uint64(obj.VictimCaseVictimContact.Int64)}
	}

	if obj.VictimCaseFollowUp.Valid {
		dto.VictimCaseFollowUp = FollowUpDTO{FollowUpId: uint64(obj.VictimCaseFollowUp.Int64)}
	}

	if obj.VictimCaseApprovedBy.Valid {
		dto.VictimCaseApprovedBy = CaseOwnerDTO{CaseOwnerId: uint64(obj.VictimCaseApprovedBy.Int64)}
	}

	if obj.VictimCaseOwnerDescription.Valid {
		dto.VictimCaseOwnerDescription = obj.VictimCaseOwnerDescription.String
	}

	if obj.VictimCaseTownLatitude.Valid {
		dto.VictimCaseTownLatitude = obj.VictimCaseTownLatitude.Float64
	}

	if obj.VictimCaseTownLongitude.Valid {
		dto.VictimCaseTownLongitude = obj.VictimCaseTownLongitude.Float64
	}

	return dto
}

func (obj *VictimCaseDTO) LoadFromVictimContactForm1(vc VictimContactDTO, moment []map[string]string, sector []map[string]string, entities []EntityDTO) {

	obj.VictimCaseVictimContactICode = vc.VictimContactICode
	obj.VictimCaseNames = vc.VictimContactNames
	obj.VictimCaseLastNames = vc.VictimContactLastNames
	obj.VictimCaseDocType = vc.VictimContactForm1.VictimContactForm1DocType
	obj.VictimCaseDocNumber = vc.VictimContactForm1.VictimContactForm1DocNumber
	obj.VictimCaseTownCode = vc.VictimContactForm1.VictimContactForm1TownCode
	obj.VictimCaseVictimContact = vc
	obj.VictimCaseForm2.VictimCaseForm2VictimPhone, _ = strconv.ParseUint(string(vc.VictimContactForm1.VictimContactForm1Phone), 10, 64)
	obj.VictimCaseForm2.VictimCaseForm2FactsDescription = vc.VictimContactForm1.VictimContactForm1FactsDescription
	obj.VictimCaseForm2.VictimCaseForm2IdentityName = vc.VictimContactForm1.VictimContactForm1Nick
	obj.VictimCaseForm2.VictimCaseForm2BirthDate = vc.VictimContactForm1.VictimContactForm1BirthDate
	obj.VictimCaseForm2.VictimCaseForm2ResidenceAddress = vc.VictimContactForm1.VictimContactForm1Address

	vc.VictimContactForm1 = VictimContactForm1DTO{}

	obj.VictimCaseEntityBranches = make(map[string]map[string]map[string]string)

	for _, m := range moment {
		obj.VictimCaseEntityBranches[m["code"]] = map[string]map[string]string{}
		for _, s := range sector {
			if obj.VictimCaseEntityBranches[m["code"]][s["code"]] == nil {
				obj.VictimCaseEntityBranches[m["code"]][s["code"]] = map[string]string{}
			}
			for _, e := range entities {
				if e.EntitySector == s["code"] {
					obj.VictimCaseEntityBranches[m["code"]][s["code"]][e.EntityICode] = ""
				}
			}
		}
	}
}

func (obj *VictimCaseDTO) LoadFromVictimContactForm2(vc VictimContactDTO, moment []map[string]string, sector []map[string]string, entities []EntityDTO) {

	obj.VictimCaseVictimContactICode = vc.VictimContactICode
	obj.VictimCaseNames = vc.VictimContactNames
	obj.VictimCaseLastNames = vc.VictimContactLastNames

	obj.VictimCaseVictimContact = vc
	obj.VictimCaseForm2.VictimCaseForm2FactsDescription = vc.VictimContactForm2.VictimContactForm2FactsDescription
	obj.VictimCaseForm2.VictimCaseForm2VictimPhone, _ = strconv.ParseUint(string(vc.VictimContactForm2.VictimContactForm2VictimColPhone), 10, 64)

	obj.VictimCaseEntityBranches = make(map[string]map[string]map[string]string)

	//Cargamos el ajuste VBG

	obj.VictimCaseForm2.VictimCaseForm2AdjustmentsGBV = vc.VictimContactForm2.VictimContactForm2AdjustmentsGBV

	for _, m := range moment {
		obj.VictimCaseEntityBranches[m["code"]] = map[string]map[string]string{}
		for _, s := range sector {
			if obj.VictimCaseEntityBranches[m["code"]][s["code"]] == nil {
				obj.VictimCaseEntityBranches[m["code"]][s["code"]] = map[string]string{}
			}
			for _, e := range entities {
				if e.EntitySector == s["code"] {
					obj.VictimCaseEntityBranches[m["code"]][s["code"]][e.EntityICode] = ""
				}
			}
		}
	}
}

func (obj *VictimCaseDTO) LoadFromWithoutVictimContact(moment []map[string]string, sector []map[string]string, entities []EntityDTO) {

	obj.VictimCaseEntityBranches = make(map[string]map[string]map[string]string)

	for _, m := range moment {
		obj.VictimCaseEntityBranches[m["code"]] = map[string]map[string]string{}
		for _, s := range sector {
			if obj.VictimCaseEntityBranches[m["code"]][s["code"]] == nil {
				obj.VictimCaseEntityBranches[m["code"]][s["code"]] = map[string]string{}
			}
			for _, e := range entities {
				if e.EntitySector == s["code"] {
					//A diferencia del anterior, aquí se referencia el ICode de la sede
					icode := obj.getBranchICodeByMomentCodeAndEntityId(m["code"], e.EntityId)
					obj.VictimCaseEntityBranches[m["code"]][s["code"]][e.EntityICode] = icode
				}
			}
		}
	}
}

func (obj *VictimCaseDTO) getBranchICodeByMomentCodeAndEntityId(momentCode string, entityId uint64) string {
	for _, m := range obj.VictimCaseMoments {
		if m.MomentCode == momentCode && m.MomentEntityBranch.EntityBranchEntity.EntityId == entityId {
			return m.MomentEntityBranch.EntityBranchICode
		}
	}
	return ""
}
