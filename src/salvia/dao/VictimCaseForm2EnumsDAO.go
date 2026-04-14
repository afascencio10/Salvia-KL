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
	"errors"
	"fmt"
)

var (
	// Entity / JSON / DB names
	VictimCaseForm2EnumsEntityName string = "VictimCaseForm2Enums"
	VictimCaseForm2EnumsJSONName   string = "enums"
	VictimCaseForm2EnumsDBName     string = "victim_case_form2_enums"
	VictimCaseForm2EnumsDBScheme   string = "salvia"

	// Field definitions – used by the generic CRUD helpers
	VictimCaseForm2EnumsFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimCaseForm2EnumsId":       {Name: "VictimCaseForm2EnumsId", DBName: "victim_case_form2_enums_id", Alias: "", ModelType: "uint64", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseForm2EnumsICode":    {Name: "VictimCaseForm2EnumsICode", DBName: "victim_case_form2_enums_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimCaseForm2EnumsName":     {Name: "VictimCaseForm2EnumsName", DBName: "victim_case_form2_enums_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 96, Required: true},
		"VictimCaseForm2EnumsCode":     {Name: "VictimCaseForm2EnumsCode", DBName: "victim_case_form2_enums_code", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: true},
		"VictimCaseForm2EnumsCategory": {Name: "VictimCaseForm2EnumsCategory", DBName: "victim_case_form2_enums_category", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: true},
	}
)

// ---------------------------------------------------------------------------
// DTO – JSON representation (used by the API)
// ---------------------------------------------------------------------------

type VictimCaseForm2EnumsDTO struct {
	VictimCaseForm2EnumsId       uint64 `json:"-"`     // PK – not exposed
	VictimCaseForm2EnumsICode    string `json:"icode"` // i‑code
	VictimCaseForm2EnumsName     string `json:"name"`  // name
	VictimCaseForm2EnumsCode     string `json:"code"`  // code
	VictimCaseForm2EnumsCategory string `json:"-"`     // category
}

// ---------------------------------------------------------------------------
// PgDB – database representation (used by the repository)
// ---------------------------------------------------------------------------

type VictimCaseForm2EnumsPgDB struct {
	VictimCaseForm2EnumsId       sql.NullString
	VictimCaseForm2EnumsICode    sql.NullString
	VictimCaseForm2EnumsName     sql.NullString
	VictimCaseForm2EnumsCode     sql.NullString
	VictimCaseForm2EnumsCategory sql.NullString
}

// Estructura enums[categoria][icode]=code
var VictimCaseForm2Enums map[string][]VictimCaseForm2EnumsDTO = make(map[string][]VictimCaseForm2EnumsDTO)

func SetVictimCaseForm2Enums(victimCaseForm2Enums *VictimCaseForm2EnumsDTO,
	connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Campos a insertar
	fieldsSlice := []string{
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{"VictimCaseForm2EnumsId"}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		victimCaseForm2Enums.VictimCaseForm2EnumsICode,
		victimCaseForm2Enums.VictimCaseForm2EnumsName,
		victimCaseForm2Enums.VictimCaseForm2EnumsCode,
		victimCaseForm2Enums.VictimCaseForm2EnumsCategory)

	persistenceCtrl.Scan(&victimCaseForm2Enums.VictimCaseForm2EnumsId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	return nil
}

func GetVictimCaseForm2Enums(by common_controllers.By, victimCaseForm2Enums *VictimCaseForm2EnumsDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm2Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsId",
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm2Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm2EnumsDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var victimCaseForm2Pg VictimCaseForm2EnumsPgDB

	persistenceCtrl.Scan(&victimCaseForm2Pg.VictimCaseForm2EnumsId,
		&victimCaseForm2Pg.VictimCaseForm2EnumsICode,
		&victimCaseForm2Pg.VictimCaseForm2EnumsName,
		&victimCaseForm2Pg.VictimCaseForm2EnumsCode,
		&victimCaseForm2Pg.VictimCaseForm2EnumsCategory)

	*victimCaseForm2Enums = victimCaseForm2Pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	return nil
}

func GetVictimCasesForm2Enums(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2EnumsDTO, error) {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm2Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm2Path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimCaseForm2EnumsDBName, by.AttrsName, []string{}, []string{}, by.Operator, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	var victimCaseForm2s []VictimCaseForm2EnumsDTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2EnumsPgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2EnumsICode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsName,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCategory)

		victimCaseForm2s = append(victimCaseForm2s, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimCaseForm2s, nil
}

func GetVictimCasesForm2EnumsByVictimCaseForm2Id(form2Id uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2EnumsDTO, error) {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm2Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	relPath := RelVictimCaseForm2EnumsVictimCaseForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsId",
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm2Path +
		` RIGHT JOIN ` + relPath + ` ON(` + victimCaseForm2EnumsForm2Path + `.` + VictimCaseForm2EnumsFieldDefinitions["VictimCaseForm2EnumsId"].DBName + ` = ` + relPath + `.` + RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimCaseForm2EnumsId"].DBName + `)` +
		` WHERE ` + RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimCaseForm2FormId"].DBName + ` = $1`

	persistenceCtrl.Query(context.Background(), query, form2Id)

	var victimCaseForm2s []VictimCaseForm2EnumsDTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2EnumsPgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2EnumsId,
			&victimCaseForm2Pg.VictimCaseForm2EnumsICode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsName,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCategory)

		victimCaseForm2s = append(victimCaseForm2s, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimCaseForm2s, nil
}

func GetVictimCasesForm2EnumsByVictimContactForm2Id(form2Id uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2EnumsDTO, error) {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm2Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	relPath := RelVictimCaseForm2EnumsVictimContactForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimContactForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsId",
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm2Path +
		` RIGHT JOIN ` + relPath + ` ON(` + victimCaseForm2EnumsForm2Path + `.` + VictimCaseForm2EnumsFieldDefinitions["VictimCaseForm2EnumsId"].DBName + ` = ` + relPath + `.` + RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimCaseForm2EnumsId"].DBName + `)` +
		` WHERE ` + RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimContactForm2FormId"].DBName + ` = $1`

	persistenceCtrl.Query(context.Background(), query, form2Id)

	var victimCaseForm2s []VictimCaseForm2EnumsDTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2EnumsPgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2EnumsId,
			&victimCaseForm2Pg.VictimCaseForm2EnumsICode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsName,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCategory)

		victimCaseForm2s = append(victimCaseForm2s, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimCaseForm2s, nil
}

func GetVictimCasesForm2EnumsByFeminicideForm1Id(form1Id uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2EnumsDTO, error) {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm1Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	relPath := RelVictimCaseForm2EnumsFeminicideForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsId",
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm1Path +
		` RIGHT JOIN ` + relPath + ` ON(` + victimCaseForm2EnumsForm1Path + `.` + VictimCaseForm2EnumsFieldDefinitions["VictimCaseForm2EnumsId"].DBName + ` = ` + relPath + `.` + RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideForm1EnumsId"].DBName + `)` +
		` WHERE ` + RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideForm1FormId"].DBName + ` = $1`

	persistenceCtrl.Query(context.Background(), query, form1Id)

	var victimCaseForm2s []VictimCaseForm2EnumsDTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2EnumsPgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2EnumsId,
			&victimCaseForm2Pg.VictimCaseForm2EnumsICode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsName,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCategory)

		victimCaseForm2s = append(victimCaseForm2s, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimCaseForm2s, nil
}

func GetVictimCasesForm2EnumsByFeminicideRiskForm1Id(form1Id uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2EnumsDTO, error) {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm1Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	relPath := RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideRiskForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsId",
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm1Path +
		` RIGHT JOIN ` + relPath + ` ON(` + victimCaseForm2EnumsForm1Path + `.` + VictimCaseForm2EnumsFieldDefinitions["VictimCaseForm2EnumsId"].DBName + ` = ` + relPath + `.` + RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId"].DBName + `)` +
		` WHERE ` + RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideRiskForm1FormId"].DBName + ` = $1`

	persistenceCtrl.Query(context.Background(), query, form1Id)

	var victimCaseForm2s []VictimCaseForm2EnumsDTO
	for persistenceCtrl.Next() {
		var victimCaseForm2Pg VictimCaseForm2EnumsPgDB
		persistenceCtrl.ScanRow(&victimCaseForm2Pg.VictimCaseForm2EnumsId,
			&victimCaseForm2Pg.VictimCaseForm2EnumsICode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsName,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCode,
			&victimCaseForm2Pg.VictimCaseForm2EnumsCategory)

		victimCaseForm2s = append(victimCaseForm2s, victimCaseForm2Pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimCaseForm2s, nil
}

func GetAllVictimCaseForm2Enums(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]VictimCaseForm2EnumsDTO, error) {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var victimCaseForm2EnumsForm2Path string = VictimCaseForm2EnumsDBScheme + "." + VictimCaseForm2EnumsDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsId",
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + victimCaseForm2EnumsForm2Path

	persistenceCtrl.Query(context.Background(), query)

	var victimCaseForm2Enums []VictimCaseForm2EnumsDTO
	for persistenceCtrl.Next() {
		var victimCaseForm2EnumsPg VictimCaseForm2EnumsPgDB
		persistenceCtrl.ScanRow(&victimCaseForm2EnumsPg.VictimCaseForm2EnumsId,
			&victimCaseForm2EnumsPg.VictimCaseForm2EnumsICode,
			&victimCaseForm2EnumsPg.VictimCaseForm2EnumsName,
			&victimCaseForm2EnumsPg.VictimCaseForm2EnumsCode,
			&victimCaseForm2EnumsPg.VictimCaseForm2EnumsCategory)

		victimCaseForm2Enums = append(victimCaseForm2Enums, victimCaseForm2EnumsPg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return victimCaseForm2Enums, nil
}

func UpdateVictimCaseForm2Enums(victimCaseForm2Enums *VictimCaseForm2EnumsDTO,
	connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimCaseForm2EnumsICode",
		"VictimCaseForm2EnumsName",
		"VictimCaseForm2EnumsCode",
		"VictimCaseForm2EnumsCategory",
	}
	fieldsAliasSlice := []string{}

	query := common_dao.GetSQL(common_dao.SQL_UPDATE, fieldsSlice, fieldsAliasSlice, VictimCaseForm2EnumsDBName, []string{"VictimCaseForm2EnumsId"}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseForm2EnumsDBScheme, VictimCaseForm2EnumsFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query,
		victimCaseForm2Enums.VictimCaseForm2EnumsId,
		victimCaseForm2Enums.VictimCaseForm2EnumsICode,
		victimCaseForm2Enums.VictimCaseForm2EnumsName,
		victimCaseForm2Enums.VictimCaseForm2EnumsCode,
		victimCaseForm2Enums.VictimCaseForm2EnumsCategory)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// Algunas utilidades
func SetVictimCaseForm2EnumsDefaults(victimCaseForm2 *VictimCaseForm2DTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		victimCaseForm2.VictimCaseForm2ICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// No hay campos que se actualicen automáticamente
	}
}

func (obj *VictimCaseForm2EnumsPgDB) ToDTO() VictimCaseForm2EnumsDTO {
	var dto VictimCaseForm2EnumsDTO

	if obj.VictimCaseForm2EnumsId.Valid {
		dto.VictimCaseForm2EnumsId, _ = strconv.ParseUint(obj.VictimCaseForm2EnumsId.String, 10, 64)
	}
	if obj.VictimCaseForm2EnumsICode.Valid {
		dto.VictimCaseForm2EnumsICode = obj.VictimCaseForm2EnumsICode.String
	}
	if obj.VictimCaseForm2EnumsName.Valid {
		dto.VictimCaseForm2EnumsName = obj.VictimCaseForm2EnumsName.String
	}
	if obj.VictimCaseForm2EnumsCode.Valid {
		dto.VictimCaseForm2EnumsCode = obj.VictimCaseForm2EnumsCode.String
	}
	if obj.VictimCaseForm2EnumsCategory.Valid {
		dto.VictimCaseForm2EnumsCategory = obj.VictimCaseForm2EnumsCategory.String
	}

	return dto
}

func GetLocalVictimCaseForm2EnumsByICode(victimCaseForm2Enums *VictimCaseForm2EnumsDTO) error {
	for _, slice := range VictimCaseForm2Enums {
		for _, e := range slice {
			if e.VictimCaseForm2EnumsICode == victimCaseForm2Enums.VictimCaseForm2EnumsICode {
				*victimCaseForm2Enums = e
				return nil
			}
		}
	}
	return errors.New("Not found")
}

func GetLocalVictimCaseForm2EnumsById(victimCaseForm2Enums *VictimCaseForm2EnumsDTO) error {
	for _, slice := range VictimCaseForm2Enums {
		for _, e := range slice {
			if e.VictimCaseForm2EnumsId == victimCaseForm2Enums.VictimCaseForm2EnumsId {
				*victimCaseForm2Enums = e
				return nil
			}
		}
	}
	return errors.New("Not found")
}

func GetLocalVictimCaseForm2EnumsByCategory(victimCaseForm2Enums VictimCaseForm2EnumsDTO) []VictimCaseForm2EnumsDTO {
	var enums []VictimCaseForm2EnumsDTO = []VictimCaseForm2EnumsDTO{}
	for _, slice := range VictimCaseForm2Enums {
		for _, e := range slice {
			if e.VictimCaseForm2EnumsCategory == victimCaseForm2Enums.VictimCaseForm2EnumsCategory {
				enums = append(enums, e)
			}
		}
	}
	return enums
}
