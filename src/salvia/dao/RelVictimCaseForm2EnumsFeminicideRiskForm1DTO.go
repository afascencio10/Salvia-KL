package salvia_daos

import (
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"

	"context"
	"database/sql"
	"fmt"
)

var (
	RelVictimCaseForm2EnumsFeminicideRiskForm1EntityName string = "RelVictimCaseForm2EnumsFeminicideRiskForm1"
	RelVictimCaseForm2EnumsFeminicideRiskForm1JSONName   string = "rel_victim_case_form2_enums_feminicide_risk_form1"
	RelVictimCaseForm2EnumsFeminicideRiskForm1DBName     string = "rel_victim_case_form2_enums_feminicide_risk_form1"
	RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme   string = "salvia"

	// Field definitions – used by the ORM / validation layer
	RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelVictimCaseForm2EnumsFeminicideRiskForm1Id":      {Name: "RelVictimCaseForm2EnumsFeminicideRiskForm1Id", DBName: "rel_victim_case_form2_enums_feminicide_risk_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId": {Name: "RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId", DBName: "victim_case_form2_enums_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsFeminicideRiskForm1FormId":  {Name: "RelVictimCaseForm2EnumsFeminicideRiskForm1FormId", DBName: "feminicide_risk_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// ---------------------------------------------------------------------------
// DTO – JSON representation (used by the API)
// ---------------------------------------------------------------------------

type RelVictimCaseForm2EnumsFeminicideRiskForm1DTO struct {
	RelVictimCaseForm2EnumsFeminicideRiskForm1Id    uint64                  `json:"-"`
	RelVictimCaseForm2EnumsFeminicideRiskForm1Enums VictimCaseForm2EnumsDTO `json:"enumsId"`
	RelVictimCaseForm2EnumsFeminicideRiskForm1Form  FeminicideRiskForm1DTO  `json:"formId"`
}

// PgDB – the database representation (sql.Null* types) ---------------------

type RelVictimCaseForm2EnumsFeminicideRiskForm1PgDB struct {
	RelVictimCaseForm2EnumsFeminicideRiskForm1Id    sql.NullInt64
	RelVictimCaseForm2EnumsFeminicideRiskForm1Enums sql.NullInt64
	RelVictimCaseForm2EnumsFeminicideRiskForm1Form  sql.NullInt64
}

func SetRelVictimCaseForm2EnumsFeminicideRiskForm1(r *RelVictimCaseForm2EnumsFeminicideRiskForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	// Prepare persistence controller
	persistenceCtrl := common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Build the INSERT query
	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1FormId",
	}
	aliasSlice := []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, []string{}, []string{}, []string{"RelVictimCaseForm2EnumsFeminicideRiskForm1Id"}, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, false)

	// Execute the query
	persistenceCtrl.QueryRow(
		context.Background(),
		query,
		r.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums.VictimCaseForm2EnumsId,
		r.RelVictimCaseForm2EnumsFeminicideRiskForm1Form.FeminicideRiskForm1Id)

	persistenceCtrl.Scan(&r.RelVictimCaseForm2EnumsFeminicideRiskForm1Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCaseForm2EnumsFeminicideRiskForm1(b common_controllers.By, r *RelVictimCaseForm2EnumsFeminicideRiskForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideRiskForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideRiskForm1Id",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, []string{}, []string{}, nil, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, true)

	fmt.Printf(query, b.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, b.AttrsValue...)

	var pg RelVictimCaseForm2EnumsFeminicideRiskForm1PgDB
	persistenceCtrl.Scan(
		&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Id,
		&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums,
		&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Form)

	*r = pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCasesForm2EnumsFeminicideRiskForm1(b common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsFeminicideRiskForm1DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideRiskForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideRiskForm1Id",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, true) +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideRiskForm1Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, b.AttrsValue...)

	var list []RelVictimCaseForm2EnumsFeminicideRiskForm1DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsFeminicideRiskForm1PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Id,
			&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums,
			&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Form)

		list = append(list, pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Count total rows for the first page
	if page == 0 {
		countQuery := `SELECT COUNT(*) FROM ` + path +
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, b.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return list, count, nil
}

// ---------------------------------------------------------------------------
// GetAllRelVictimCaseForm2EnumsFeminicideRiskForm1 – SELECT all (no filter)
// ---------------------------------------------------------------------------
func GetAllRelVictimCaseForm2EnumsFeminicideRiskForm1(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsFeminicideRiskForm1DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideRiskForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideRiskForm1Id",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideRiskForm1FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideRiskForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideRiskForm1DBScheme, RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsFeminicideRiskForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideRiskForm1Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query)

	var list []RelVictimCaseForm2EnumsFeminicideRiskForm1DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsFeminicideRiskForm1PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Id,
			&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums,
			&pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Form)

		list = append(list, pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Count total rows for the first page
	if page == 0 {
		countQuery := `SELECT COUNT(*) FROM ` + path

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return list, count, nil
}

func SetRelVictimCaseForm2EnumsFeminicideRiskForm1Defaults(r *RelVictimCaseForm2EnumsFeminicideRiskForm1DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:

	case common_dao.SQL_UPDATE:

	}
}

func (pg *RelVictimCaseForm2EnumsFeminicideRiskForm1PgDB) ToDTO() RelVictimCaseForm2EnumsFeminicideRiskForm1DTO {
	var dto RelVictimCaseForm2EnumsFeminicideRiskForm1DTO

	if pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Id.Valid {
		dto.RelVictimCaseForm2EnumsFeminicideRiskForm1Id = uint64(pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Id.Int64)
	}
	if pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums.Valid {
		dto.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Enums.Int64)}
	}
	if pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Form.Valid {
		dto.RelVictimCaseForm2EnumsFeminicideRiskForm1Form = FeminicideRiskForm1DTO{FeminicideRiskForm1Id: uint64(pg.RelVictimCaseForm2EnumsFeminicideRiskForm1Form.Int64)}
	}

	return dto
}
