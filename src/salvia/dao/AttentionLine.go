package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	AttentionLineEntityName string = "AttentionLine"
	AttentionLineJSONName   string = "attentionLine"
	AttentionLineDBName     string = "attention_line"
	AttentionLineDBScheme   string = "salvia"

	//Atributos relacionados con las validaciones ------------------------------

	//Campos que vienen como string del JSON. El booleano indica si son strings en el modelo o no (como en el caso de una fecha)
	AttentionLineFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"AttentionLineId":           {Name: "AttentionLineId", DBName: "attention_line_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"AttentionLineICode":        {Name: "AttentionLineICode", DBName: "attention_line_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"AttentionLineCreationDate": {Name: "AttentionLineCreationDate", DBName: "attention_line_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"AttentionLineUpdateDate":   {Name: "AttentionLineUpdateDate", DBName: "attention_line_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"AttentionLineName":         {Name: "AttentionLineName", DBName: "attention_line_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 254, Required: true},
		"AttentionLineDescription":  {Name: "AttentionLineDescription", DBName: "attention_line_description", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 5000, Required: true},
	}
)

type AttentionLineDTO struct {
	AttentionLineId           uint64         `json:"-"`
	AttentionLineICode        string         `json:"icode"`
	AttentionLineCreationDate time.Time      `json:"creationDate"`
	AttentionLineUpdateDate   time.Time      `json:"updateDate"`
	AttentionLineName         string         `json:"name"`
	AttentionLineDescription  string         `json:"description"`
	AttentionLineUsers        []CaseOwnerDTO `json:"users"`
}
type AttentionLinePgDB struct {
	AttentionLineId           sql.NullInt64
	AttentionLineICode        sql.NullString
	AttentionLineCreationDate sql.NullTime
	AttentionLineUpdateDate   sql.NullTime
	AttentionLineName         sql.NullString
	AttentionLineDescription  sql.NullString
	AttentionLineUsers        []sql.NullInt64
	//Campos de formulario que no hacen parte del modelo -----------------------------------

}

func SetAttentionLine(attentionLine *AttentionLineDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query

	var attentionLineFieldsSlice []string = []string{"AttentionLineICode", "AttentionLineCreationDate", "AttentionLineUpdateDate", "AttentionLineName", "AttentionLineDescription"}
	var attentionLineFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, attentionLineFieldsSlice, attentionLineFieldsAliasSlice, AttentionLineDBName, []string{}, []string{}, []string{"AttentionLineId"}, common_dao.SQL_AND, AttentionLineDBScheme, AttentionLineFieldDefinitions, false)

	//var query string = `INSERT INTO ` + AttentionLineDBScheme + `.attention_line (attention_line_creation_date,attention_line_update_date,attention_line_login,attention_line_password,attention_line_status,attention_line_language` + extraColumns + `)
	//	VALUES (NOW(),NOW(),$1,$2,$3,$4` + extraValues + `) RETURNING attention_line_id, attention_line_i_code`

	persistenceCtrl.QueryRow(context.Background(), query, attentionLine.AttentionLineICode, attentionLine.AttentionLineCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), attentionLine.AttentionLineUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), attentionLine.AttentionLineName, attentionLine.AttentionLineDescription)

	persistenceCtrl.Scan(&attentionLine.AttentionLineId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func UpdateAttentionLineByICode(attentionLine *AttentionLineDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var attentionLineFieldsSlice []string = []string{"AttentionLineUpdateDate", "AttentionLineName", "AttentionLineDescription"}
	var attentionLineFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, attentionLineFieldsSlice, attentionLineFieldsAliasSlice, AttentionLineDBName, []string{}, []string{"AttentionLineICode"}, []string{}, common_dao.SQL_AND, AttentionLineDBScheme, AttentionLineFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query,
		attentionLine.AttentionLineICode, attentionLine.AttentionLineUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), attentionLine.AttentionLineName, attentionLine.AttentionLineDescription)

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

func RemoveAttentionLineByICode(attentionLine *AttentionLineDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var attentionLineFieldsSlice []string = []string{}
	var attentionLineFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, attentionLineFieldsSlice, attentionLineFieldsAliasSlice, AttentionLineDBName, []string{}, []string{"AttentionLineICode"}, []string{}, common_dao.SQL_AND, AttentionLineDBScheme, AttentionLineFieldDefinitions, false)

	//var query string = `DELETE FROM ` + AttentionLineDBScheme + `.attention_line
	//	WHERE attention_line_id = $1`

	persistenceCtrl.Exec(context.Background(), query, attentionLine.AttentionLineICode)

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

func GetAttentionLine(by common_controllers.By, attentionLine *AttentionLineDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var attentionLinePath string = AttentionLineDBScheme + "." + AttentionLineDBName

	//Obtenemos la conexión
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Query
	var attentionLineFieldsSlice []string = []string{"AttentionLineId", "AttentionLineICode", "AttentionLineCreationDate", "AttentionLineUpdateDate", "AttentionLineName", "AttentionLineDescription"}
	var attentionLineFieldsAliasSlice []string = []string{}

	var attentionLineFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, attentionLineFieldsSlice, attentionLineFieldsAliasSlice, AttentionLineDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AttentionLineDBScheme, AttentionLineFieldDefinitions, true)

	var query string = `SELECT ` + attentionLineFieldsStr +
		` FROM ` + attentionLinePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, AttentionLineDBName, by.AttrsName, []string{}, []string{}, by.Operator, AttentionLineDBScheme, AttentionLineFieldDefinitions, true)

	//var query string = `SELECT attention_line_id, attention_line_i_code, attention_line_creation_date,attention_line_update_date,attention_line_login,attention_line_password,attention_line_status,attention_line_language FROM ` + AttentionLineDBScheme + `.attention_line
	//	WHERE attention_line_i_code = $1`

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	persistenceCtrl.Scan(&attentionLine.AttentionLineId, &attentionLine.AttentionLineICode, &attentionLine.AttentionLineCreationDate, &attentionLine.AttentionLineUpdateDate,
		&attentionLine.AttentionLineName, &attentionLine.AttentionLineDescription)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetAllAttentionLine(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AttentionLineDTO, error) {
	// Definición de variables

	var attentionLine AttentionLineDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Query
	var attentionLineFieldsSlice []string = []string{"AttentionLineICode", "AttentionLineCreationDate", "AttentionLineUpdateDate", "AttentionLineName", "AttentionLineDescription"}
	var attentionLineFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, attentionLineFieldsSlice, attentionLineFieldsAliasSlice, AttentionLineDBName, []string{}, []string{}, []string{}, "", AttentionLineDBScheme, AttentionLineFieldDefinitions, true)

	//var query string = `SELECT attention_line_id, attention_line_i_code, attention_line_creation_date,attention_line_update_date,attention_line_login,attention_line_password,attention_line_status,attention_line_language FROM ` + AttentionLineDBScheme + `.attention_line`

	persistenceCtrl.Query(context.Background(), query)
	var attentionLines []AttentionLineDTO
	for persistenceCtrl.Next() {
		attentionLine = AttentionLineDTO{}
		persistenceCtrl.ScanRow(&attentionLine.AttentionLineICode, &attentionLine.AttentionLineCreationDate, &attentionLine.AttentionLineUpdateDate,
			&attentionLine.AttentionLineName, &attentionLine.AttentionLineDescription)

		attentionLines = append(attentionLines, attentionLine)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return attentionLines, nil
}

// Algunas utilidades
func SetAttentionLineDefaults(attentionLine *AttentionLineDTO, action string) {

	switch action {
	case common_dao.SQL_INSERT:
		attentionLine.AttentionLineCreationDate = time.Now()
		attentionLine.AttentionLineUpdateDate = time.Now()
		attentionLine.AttentionLineICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		attentionLine.AttentionLineUpdateDate = time.Now()
	}

}

func (obj *AttentionLinePgDB) ToDTO() AttentionLineDTO {
	var dto AttentionLineDTO

	if obj.AttentionLineId.Valid {

		dto.AttentionLineId = uint64(obj.AttentionLineId.Int64)
	}

	if obj.AttentionLineICode.Valid {

		dto.AttentionLineICode = obj.AttentionLineICode.String
	}

	if obj.AttentionLineCreationDate.Valid {

		dto.AttentionLineCreationDate = obj.AttentionLineCreationDate.Time
	}

	if obj.AttentionLineUpdateDate.Valid {

		dto.AttentionLineUpdateDate = obj.AttentionLineUpdateDate.Time
	}

	if obj.AttentionLineName.Valid {

		dto.AttentionLineName = obj.AttentionLineName.String
	}

	if obj.AttentionLineDescription.Valid {

		dto.AttentionLineDescription = obj.AttentionLineDescription.String
	}

	return dto
}
