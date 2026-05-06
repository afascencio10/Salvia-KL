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
	RelVictimCaseForm2EnumsVictimCaseForm2EntityName string = "RelVictimCaseForm2EnumsVictimCaseForm2"
	RelVictimCaseForm2EnumsVictimCaseForm2JSONName   string = "relVictimCaseForm2EnumsVictimCaseForm2"
	RelVictimCaseForm2EnumsVictimCaseForm2DBName     string = "rel_victim_case_form2_enums_victim_case_form2"
	RelVictimCaseForm2EnumsVictimCaseForm2DBScheme   string = "salvia"

	// Field definitions – used by the ORM / validation layer
	RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelVictimCaseForm2EnumsVictimCaseForm2Id":      {Name: "RelVictimCaseForm2EnumsVictimCaseForm2Id", DBName: "rel_victim_case_form2_enums_victim_case_form2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId": {Name: "RelVictimCaseForm2EnumsVictimCaseForm2EnumsId", DBName: "victim_case_form2_enums_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsVictimCaseForm2FormId":  {Name: "RelVictimCaseForm2EnumsVictimCaseForm2FormId", DBName: "victim_case_form2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// ---------------------------------------------------------------------------
// DTO – JSON representation (used by the API)
// ---------------------------------------------------------------------------

type RelVictimCaseForm2EnumsVictimCaseForm2DTO struct {
	RelVictimCaseForm2EnumsVictimCaseForm2Id    uint64                  `json:"-"`
	RelVictimCaseForm2EnumsVictimCaseForm2Enums VictimCaseForm2EnumsDTO `json:"enumsId"`
	RelVictimCaseForm2EnumsVictimCaseForm2Form  VictimCaseForm2DTO      `json:"formId"`
}

// PgDB – the database representation (sql.Null* types) ---------------------

type RelVictimCaseForm2EnumsVictimCaseForm2PgDB struct {
	RelVictimCaseForm2EnumsVictimCaseForm2Id    sql.NullInt64
	RelVictimCaseForm2EnumsVictimCaseForm2Enums sql.NullInt64
	RelVictimCaseForm2EnumsVictimCaseForm2Form  sql.NullInt64
}

func SetRelVictimCaseForm2EnumsVictimCaseForm2(r *RelVictimCaseForm2EnumsVictimCaseForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	// Prepare persistence controller
	persistenceCtrl := common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Build the INSERT query
	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimCaseForm2FormId",
	}
	aliasSlice := []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimCaseForm2DBName, []string{}, []string{}, []string{"RelVictimCaseForm2EnumsVictimCaseForm2Id"}, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, false)

	// Execute the query
	persistenceCtrl.QueryRow(
		context.Background(),
		query,
		r.RelVictimCaseForm2EnumsVictimCaseForm2Enums.VictimCaseForm2EnumsId,
		r.RelVictimCaseForm2EnumsVictimCaseForm2Form.VictimCaseForm2Id)

	persistenceCtrl.Scan(&r.RelVictimCaseForm2EnumsVictimCaseForm2Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCaseForm2EnumsVictimCaseForm2(b common_controllers.By, r *RelVictimCaseForm2EnumsVictimCaseForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsVictimCaseForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimCaseForm2Id",
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimCaseForm2FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimCaseForm2DBName, []string{}, []string{}, nil, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsVictimCaseForm2DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, true)

	fmt.Printf(query, b.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, b.AttrsValue...)

	var pg RelVictimCaseForm2EnumsVictimCaseForm2PgDB
	persistenceCtrl.Scan(
		&pg.RelVictimCaseForm2EnumsVictimCaseForm2Id,
		&pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums,
		&pg.RelVictimCaseForm2EnumsVictimCaseForm2Form)

	*r = pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCasesForm2EnumsVictimCaseForm2(b common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsVictimCaseForm2DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsVictimCaseForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimCaseForm2Id",
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimCaseForm2FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimCaseForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsVictimCaseForm2DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, true) +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimCaseForm2Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, b.AttrsValue...)

	var list []RelVictimCaseForm2EnumsVictimCaseForm2DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsVictimCaseForm2PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Id,
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums,
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Form)

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
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsVictimCaseForm2DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, b.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return list, count, nil
}

// ---------------------------------------------------------------------------
// GetAllRelVictimCaseForm2EnumsVictimCaseForm2 – SELECT all (no filter)
// ---------------------------------------------------------------------------
func GetAllRelVictimCaseForm2EnumsVictimCaseForm2(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsVictimCaseForm2DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsVictimCaseForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimCaseForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimCaseForm2Id",
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimCaseForm2FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimCaseForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimCaseForm2Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query)

	var list []RelVictimCaseForm2EnumsVictimCaseForm2DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsVictimCaseForm2PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Id,
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums,
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Form)

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

func RemoveRelVictimCaseForm2EnumsVictimCaseForm2(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se genera la consulta SQL para eliminar relaciones basándose en los atributos indicados.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, RelVictimCaseForm2EnumsVictimCaseForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, RelVictimCaseForm2EnumsVictimCaseForm2DBScheme, RelVictimCaseForm2EnumsVictimCaseForm2FieldDefinitions, false)

	// Se ejecuta la consulta con los valores de los atributos.
	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func SetRelVictimCaseForm2EnumsVictimCaseForm2Defaults(r *RelVictimCaseForm2EnumsVictimCaseForm2DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:

	case common_dao.SQL_UPDATE:

	}
}

func (pg *RelVictimCaseForm2EnumsVictimCaseForm2PgDB) ToDTO() RelVictimCaseForm2EnumsVictimCaseForm2DTO {
	var dto RelVictimCaseForm2EnumsVictimCaseForm2DTO

	if pg.RelVictimCaseForm2EnumsVictimCaseForm2Id.Valid {
		dto.RelVictimCaseForm2EnumsVictimCaseForm2Id = uint64(pg.RelVictimCaseForm2EnumsVictimCaseForm2Id.Int64)
	}
	if pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums.Valid {
		dto.RelVictimCaseForm2EnumsVictimCaseForm2Enums = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums.Int64)}
	}
	if pg.RelVictimCaseForm2EnumsVictimCaseForm2Form.Valid {
		dto.RelVictimCaseForm2EnumsVictimCaseForm2Form = VictimCaseForm2DTO{VictimCaseForm2Id: uint64(pg.RelVictimCaseForm2EnumsVictimCaseForm2Form.Int64)}
	}

	return dto
}
