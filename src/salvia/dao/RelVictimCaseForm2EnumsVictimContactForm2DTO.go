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
	RelVictimCaseForm2EnumsVictimContactForm2EntityName string = "RelVictimCaseForm2EnumsVictimContactForm2"
	RelVictimCaseForm2EnumsVictimContactForm2JSONName   string = "rel_victim_case_form2_enums_victim_contact_form2"
	RelVictimCaseForm2EnumsVictimContactForm2DBName     string = "rel_victim_case_form2_enums_victim_contact_form2"
	RelVictimCaseForm2EnumsVictimContactForm2DBScheme   string = "salvia"

	// Field definitions – used by the ORM / validation layer
	RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelVictimCaseForm2EnumsVictimContactForm2Id":     {Name: "RelVictimCaseForm2EnumsVictimContactForm2Id", DBName: "rel_victim_case_form2_enums_victim_contact_form2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId":   {Name: "RelVictimCaseForm2EnumsVictimCaseForm2EnumsId", DBName: "victim_case_form2_enums_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelVictimCaseForm2EnumsVictimContactForm2FormId": {Name: "RelVictimCaseForm2EnumsVictimContactForm2FormId", DBName: "victim_contact_form2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// ---------------------------------------------------------------------------
// DTO – JSON representation (used by the API)
// ---------------------------------------------------------------------------

type RelVictimCaseForm2EnumsVictimContactForm2DTO struct {
	RelVictimCaseForm2EnumsVictimContactForm2Id   uint64                  `json:"-"`
	RelVictimCaseForm2EnumsVictimCaseForm2Enums   VictimCaseForm2EnumsDTO `json:"enumsId"`
	RelVictimCaseForm2EnumsVictimContactForm2Form VictimContactForm2DTO   `json:"formId"`
}

// PgDB – the database representation (sql.Null* types) ---------------------

type RelVictimCaseForm2EnumsVictimContactForm2PgDB struct {
	RelVictimCaseForm2EnumsVictimContactForm2Id   sql.NullInt64
	RelVictimCaseForm2EnumsVictimCaseForm2Enums   sql.NullInt64
	RelVictimCaseForm2EnumsVictimContactForm2Form sql.NullInt64
}

func SetRelVictimCaseForm2EnumsVictimContactForm2(r *RelVictimCaseForm2EnumsVictimContactForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

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
		"RelVictimCaseForm2EnumsVictimContactForm2FormId",
	}
	aliasSlice := []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimContactForm2DBName, []string{}, []string{}, []string{"RelVictimCaseForm2EnumsVictimContactForm2Id"}, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, false)

	// Execute the query
	persistenceCtrl.QueryRow(
		context.Background(),
		query,
		r.RelVictimCaseForm2EnumsVictimCaseForm2Enums.VictimCaseForm2EnumsId,
		r.RelVictimCaseForm2EnumsVictimContactForm2Form.VictimContactForm2Id)

	persistenceCtrl.Scan(&r.RelVictimCaseForm2EnumsVictimContactForm2Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimCaseForm2EnumsVictimContactForm2(b common_controllers.By, r *RelVictimCaseForm2EnumsVictimContactForm2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsVictimContactForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimContactForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimContactForm2Id",
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimContactForm2FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimContactForm2DBName, []string{}, []string{}, nil, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsVictimContactForm2DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, true)

	fmt.Printf(query, b.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, b.AttrsValue...)

	var pg RelVictimCaseForm2EnumsVictimContactForm2PgDB
	persistenceCtrl.Scan(
		&pg.RelVictimCaseForm2EnumsVictimContactForm2Id,
		&pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums,
		&pg.RelVictimCaseForm2EnumsVictimContactForm2Form)

	*r = pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetRelVictimContactsForm2EnumsVictimContactForm2(b common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsVictimContactForm2DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsVictimContactForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimContactForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimContactForm2Id",
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimContactForm2FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimContactForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsVictimContactForm2DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, true) +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimContactForm2Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, b.AttrsValue...)

	var list []RelVictimCaseForm2EnumsVictimContactForm2DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsVictimContactForm2PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsVictimContactForm2Id,
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums,
			&pg.RelVictimCaseForm2EnumsVictimContactForm2Form)

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
			common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, b.AttrsName, b.AttrsAliasName, RelVictimCaseForm2EnumsVictimContactForm2DBName, b.AttrsName, []string{}, []string{}, b.Operator, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, true)

		persistenceCtrl.QueryRow(context.Background(), countQuery, b.AttrsValue...)
		persistenceCtrl.Scan(&count)
	}

	return list, count, nil
}

// ---------------------------------------------------------------------------
// GetAllRelVictimCaseForm2EnumsVictimContactForm2 – SELECT all (no filter)
// ---------------------------------------------------------------------------
func GetAllRelVictimCaseForm2EnumsVictimContactForm2(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelVictimCaseForm2EnumsVictimContactForm2DTO, int, error) {

	count := 0
	persistenceCtrl := common_controllers.PersistenceController{}
	path := RelVictimCaseForm2EnumsVictimContactForm2DBScheme + "." + RelVictimCaseForm2EnumsVictimContactForm2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"RelVictimCaseForm2EnumsVictimContactForm2Id",
		"RelVictimCaseForm2EnumsVictimCaseForm2EnumsId",
		"RelVictimCaseForm2EnumsVictimContactForm2FormId",
	}
	aliasSlice := []string{}
	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, RelVictimCaseForm2EnumsVictimContactForm2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		` ORDER BY ` + path + `.` + RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions["RelVictimCaseForm2EnumsVictimContactForm2Id"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query)

	var list []RelVictimCaseForm2EnumsVictimContactForm2DTO
	for persistenceCtrl.Next() {
		var pg RelVictimCaseForm2EnumsVictimContactForm2PgDB
		persistenceCtrl.ScanRow(
			&pg.RelVictimCaseForm2EnumsVictimContactForm2Id,
			&pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums,
			&pg.RelVictimCaseForm2EnumsVictimContactForm2Form)

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

func RemoveRelVictimCaseForm2EnumsVictimContactForm2(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se genera la consulta SQL para eliminar relaciones basándose en los atributos indicados.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, RelVictimCaseForm2EnumsVictimContactForm2DBName, by.AttrsName, []string{}, []string{}, by.Operator, RelVictimCaseForm2EnumsVictimContactForm2DBScheme, RelVictimCaseForm2EnumsVictimContactForm2FieldDefinitions, false)

	// Se ejecuta la consulta con los valores de los atributos.
	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func SetRelVictimCaseForm2EnumsVictimContactForm2Defaults(r *RelVictimCaseForm2EnumsVictimContactForm2DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:

	case common_dao.SQL_UPDATE:

	}
}

func (pg *RelVictimCaseForm2EnumsVictimContactForm2PgDB) ToDTO() RelVictimCaseForm2EnumsVictimContactForm2DTO {
	var dto RelVictimCaseForm2EnumsVictimContactForm2DTO

	if pg.RelVictimCaseForm2EnumsVictimContactForm2Id.Valid {
		dto.RelVictimCaseForm2EnumsVictimContactForm2Id = uint64(pg.RelVictimCaseForm2EnumsVictimContactForm2Id.Int64)
	}
	if pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums.Valid {
		dto.RelVictimCaseForm2EnumsVictimCaseForm2Enums = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: uint64(pg.RelVictimCaseForm2EnumsVictimCaseForm2Enums.Int64)}
	}
	if pg.RelVictimCaseForm2EnumsVictimContactForm2Form.Valid {
		dto.RelVictimCaseForm2EnumsVictimContactForm2Form = VictimContactForm2DTO{VictimContactForm2Id: uint64(pg.RelVictimCaseForm2EnumsVictimContactForm2Form.Int64)}
	}

	return dto
}
