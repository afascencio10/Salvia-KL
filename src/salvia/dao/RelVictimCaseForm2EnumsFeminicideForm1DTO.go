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
	RelVictimCaseForm2EnumsFeminicideForm1EntityName string = "RelVictimCaseForm2EnumsFeminicideForm1"
	RelVictimCaseForm2EnumsFeminicideForm1JSONName   string = "relVictimVaseForm2EnumsFeminicideForm1"
	RelVictimCaseForm2EnumsFeminicideForm1DBName     string = "rel_victim_case_form2_enums_feminicide_form1"
	RelVictimCaseForm2EnumsFeminicideForm1DBScheme   string = "salvia"

	// Field definitions – used by the ORM / validation layer
	RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelVictimCaseForm2EnumsFeminicideForm1Id":      {Name: "RelVictimCaseForm2EnumsFeminicideForm1Id", DBName: "rel_victim_case_form2_enums_feminicide_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsFeminicideForm1EnumsId": {Name: "RelVictimCaseForm2EnumsFeminicideForm1EnumsId", DBName: "victim_case_form2_enums_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsFeminicideForm1FormId":  {Name: "RelVictimCaseForm2EnumsFeminicideForm1FormId", DBName: "feminicide_form1_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// ---------------------------------------------------------------------------
// DTO – JSON representation (used by the API)
// ---------------------------------------------------------------------------

type RelVictimCaseForm2EnumsFeminicideForm1DTO struct {
	RelVictimCaseForm2EnumsFeminicideForm1Id    uint64                  `json:"-"`
	RelVictimCaseForm2EnumsFeminicideForm1Enums VictimCaseForm2EnumsDTO `json:"enumsId"`
	RelVictimCaseForm2EnumsFeminicideForm1Form  FeminicideForm1DTO      `json:"formId"`
}

// PgDB – the database representation (sql.Null* types) ---------------------

type RelVictimCaseForm2EnumsFeminicideForm1PgDB struct {
	RelVictimCaseForm2EnumsFeminicideForm1Id    sql.NullInt64
	RelVictimCaseForm2EnumsFeminicideForm1Enums sql.NullInt64
	RelVictimCaseForm2EnumsFeminicideForm1Form  sql.NullInt64
}

func SetRelVictimCaseForm2EnumsFeminicideForm1(r *RelVictimCaseForm2EnumsFeminicideForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	// Prepare persistence controller
	persistenceCtrl := common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Build the INSERT query
	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideForm1FormId",
	}
	aliasSlice := []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideForm1DBName, []string{}, []string{}, []string{"RelVictimCaseForm2EnumsFeminicideForm1Id"}, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, false)

	// Execute the query
	persistenceCtrl.QueryRow(
		context.Background(),
		query,
		r.RelVictimCaseForm2EnumsFeminicideForm1Enums.VictimCaseForm2EnumsId,
		r.RelVictimCaseForm2EnumsFeminicideForm1Form.FeminicideForm1Id)

	persistenceCtrl.Scan(&r.RelVictimCaseForm2EnumsFeminicideForm1Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCaseForm2EnumsFeminicideForm1(b common_controllers.By, r *RelVictimCaseForm2EnumsFeminicideForm1DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsFeminicideForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideForm1Id",
		"RelVictimCaseForm2EnumsFeminicideForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideForm1FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideForm1DBName, []string{}, []string{}, nil, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsFeminicideForm1DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, true)

	fmt.Printf(query, b.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, b.AttrsValue...)

	var pg RelVictimCaseForm2EnumsFeminicideForm1PgDB
	persistenceCtrl.Scan(
		&pg.RelVictimCaseForm2EnumsFeminicideForm1Id,
		&pg.RelVictimCaseForm2EnumsFeminicideForm1Enums,
		&pg.RelVictimCaseForm2EnumsFeminicideForm1Form)

	*r = pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCasesForm2EnumsFeminicideForm1(b common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsFeminicideForm1DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsFeminicideForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideForm1Id",
		"RelVictimCaseForm2EnumsFeminicideForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideForm1FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsFeminicideForm1DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, true) +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideForm1Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, b.AttrsValue...)

	var list []RelVictimCaseForm2EnumsFeminicideForm1DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsFeminicideForm1PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsFeminicideForm1Id,
			&pg.RelVictimCaseForm2EnumsFeminicideForm1Enums,
			&pg.RelVictimCaseForm2EnumsFeminicideForm1Form)

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
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsFeminicideForm1DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, b.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return list, count, nil
}

// ---------------------------------------------------------------------------
// GetAllRelVictimCaseForm2EnumsFeminicideForm1 – SELECT all (no filter)
// ---------------------------------------------------------------------------
func GetAllRelVictimCaseForm2EnumsFeminicideForm1(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsFeminicideForm1DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsFeminicideForm1DBScheme + "." + RelVictimCaseForm2EnumsFeminicideForm1DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsFeminicideForm1Id",
		"RelVictimCaseForm2EnumsFeminicideForm1EnumsId",
		"RelVictimCaseForm2EnumsFeminicideForm1FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsFeminicideForm1DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsFeminicideForm1DBScheme, RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsFeminicideForm1FieldDefinitions["RelVictimCaseForm2EnumsFeminicideForm1Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query)

	var list []RelVictimCaseForm2EnumsFeminicideForm1DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsFeminicideForm1PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsFeminicideForm1Id,
			&pg.RelVictimCaseForm2EnumsFeminicideForm1Enums,
			&pg.RelVictimCaseForm2EnumsFeminicideForm1Form)

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

func SetRelVictimCaseForm2EnumsFeminicideForm1Defaults(r *RelVictimCaseForm2EnumsFeminicideForm1DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:

	case common_dao.SQL_UPDATE:

	}
}

func (pg *RelVictimCaseForm2EnumsFeminicideForm1PgDB) ToDTO() RelVictimCaseForm2EnumsFeminicideForm1DTO {
	var dto RelVictimCaseForm2EnumsFeminicideForm1DTO

	if pg.RelVictimCaseForm2EnumsFeminicideForm1Id.Valid {
		dto.RelVictimCaseForm2EnumsFeminicideForm1Id = uint64(pg.RelVictimCaseForm2EnumsFeminicideForm1Id.Int64)
	}
	if pg.RelVictimCaseForm2EnumsFeminicideForm1Enums.Valid {
		dto.RelVictimCaseForm2EnumsFeminicideForm1Enums = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(pg.RelVictimCaseForm2EnumsFeminicideForm1Enums.Int64)}
	}
	if pg.RelVictimCaseForm2EnumsFeminicideForm1Form.Valid {
		dto.RelVictimCaseForm2EnumsFeminicideForm1Form = FeminicideForm1DTO{FeminicideForm1Id: uint64(pg.RelVictimCaseForm2EnumsFeminicideForm1Form.Int64)}
	}

	return dto
}
