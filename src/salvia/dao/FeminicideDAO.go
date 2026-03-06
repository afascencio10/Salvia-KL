package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_daos "bitsflow/security/dao"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	FeminicideEntityName string = "Feminicide"
	FeminicideJSONName   string = "feminicide"
	FeminicideDBName     string = "feminicide"
	FeminicideDBScheme   string = "public" // Ajusta según tu esquema

	FeminicideFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FeminicideId":              {Name: "FeminicideId", DBName: "feminicide_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideICode":           {Name: "FeminicideICode", DBName: "feminicide_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"FeminicideCreationDate":    {Name: "FeminicideCreationDate", DBName: "feminicide_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideUpdateDate":      {Name: "FeminicideUpdateDate", DBName: "feminicide_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideStatus":          {Name: "FeminicideStatus", DBName: "feminicide_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FeminicideGeneralUser":     {Name: "FeminicideGeneralUser", DBName: "feminicide_general_user", Alias: "", ModelType: "string", MinSize: 36, MaxSize: 36, Required: true},
		"FeminicideNames":           {Name: "FeminicideNames", DBName: "feminicide_names", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: true},
		"FeminicideLastNames":       {Name: "FeminicideLastNames", DBName: "feminicide_last_names", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: true},
		"FeminicideVictimDocType":   {Name: "FeminicideVictimDocType", DBName: "feminicide_victim_doc_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"FeminicideVictimDocNumber": {Name: "FeminicideVictimDocNumber", DBName: "feminicide_victim_doc_number", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: true},
	}
)

type FeminicideDTO struct {
	FeminicideId              uint64             `json:"id"`
	FeminicideICode           string             `json:"icode"`
	FeminicideCreationDate    time.Time          `json:"creationDate"`
	FeminicideUpdateDate      time.Time          `json:"updateDate"`
	FeminicideStatus          string             `json:"status"`
	FeminicideGeneralUser     string             `json:"generalUser"`
	FeminicideNames           string             `json:"names"`
	FeminicideLastNames       string             `json:"lastNames"`
	FeminicideVictimDocType   string             `json:"victimDocType"`
	FeminicideVictimDocNumber string             `json:"victimDocNumber"`
	FeminicideForm1           FeminicideForm1DTO `json:"form"`
}

type FeminicidePgDB struct {
	FeminicideId              sql.NullInt64
	FeminicideICode           sql.NullString
	FeminicideCreationDate    sql.NullTime
	FeminicideUpdateDate      sql.NullTime
	FeminicideStatus          sql.NullString
	FeminicideGeneralUser     sql.NullString
	FeminicideNames           sql.NullString
	FeminicideLastNames       sql.NullString
	FeminicideVictimDocType   sql.NullString
	FeminicideVictimDocNumber sql.NullString
}

// Manejo para las fechas y tiempos
func (f *FeminicideDTO) MarshalJSON() ([]byte, error) {
	type Alias FeminicideDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		FeminicideCreationDate string `json:"creationDate"`
		FeminicideUpdateDate   string `json:"updateDate"`
	}{
		Alias:                  (*Alias)(f),
		FeminicideCreationDate: f.FeminicideCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FeminicideUpdateDate:   f.FeminicideUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

// Manejo para las fechas y tiempos
func (f *FeminicideDTO) UnmarshalJSON(data []byte) error {
	type Alias FeminicideDTO

	// Estructura auxiliar donde las fechas vienen como string
	aux := &struct {
		*Alias
		FeminicideCreationDate string `json:"creationDate"`
		FeminicideUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(f),
	}

	// Primero unmarshalea todo el JSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Función helper para parsear fechas sin romper el flujo
	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{} // Fecha inválida → zero value
		}
		return t
	}

	// Parseo tolerante
	f.FeminicideCreationDate = parse(aux.FeminicideCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	f.FeminicideUpdateDate = parse(aux.FeminicideUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)
	return nil
}

func SetFeminicide(feminicide *FeminicideDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideICode", "FeminicideCreationDate", "FeminicideUpdateDate",
		"FeminicideStatus", "FeminicideGeneralUser", "FeminicideNames",
		"FeminicideLastNames", "FeminicideVictimDocType", "FeminicideVictimDocNumber",
	}

	var aliasesSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, aliasesSlice, FeminicideDBName, []string{}, []string{}, []string{"FeminicideId"}, common_dao.SQL_AND, FeminicideDBScheme, FeminicideFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		feminicide.FeminicideICode,
		feminicide.FeminicideCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicide.FeminicideUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicide.FeminicideStatus,
		feminicide.FeminicideGeneralUser,
		feminicide.FeminicideNames,
		feminicide.FeminicideLastNames,
		feminicide.FeminicideVictimDocType,
		feminicide.FeminicideVictimDocNumber,
	)

	persistenceCtrl.Scan(&feminicide.FeminicideId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicide(by common_controllers.By, feminicide *FeminicideDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicidePath string = FeminicideDBScheme + "." + FeminicideDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideId", "FeminicideICode", "FeminicideCreationDate", "FeminicideUpdateDate",
		"FeminicideStatus", "FeminicideGeneralUser", "FeminicideNames",
		"FeminicideLastNames", "FeminicideVictimDocType", "FeminicideVictimDocNumber",
	}
	var aliasesSlice []string = []string{}

	var fieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasesSlice, FeminicideDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideDBScheme, FeminicideFieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + feminicidePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideDBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideDBScheme, FeminicideFieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var pg FeminicidePgDB
	persistenceCtrl.Scan(
		&pg.FeminicideId,
		&pg.FeminicideICode,
		&pg.FeminicideCreationDate,
		&pg.FeminicideUpdateDate,
		&pg.FeminicideStatus,
		&pg.FeminicideGeneralUser,
		&pg.FeminicideNames,
		&pg.FeminicideLastNames,
		&pg.FeminicideVictimDocType,
		&pg.FeminicideVictimDocNumber,
	)

	*feminicide = pg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicides(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideDTO, int, error) {
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicidePath string = FeminicideDBScheme + "." + FeminicideDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{"FeminicideId",
		"FeminicideICode", "FeminicideCreationDate", "FeminicideUpdateDate",
		"FeminicideStatus", "FeminicideGeneralUser", "FeminicideNames",
		"FeminicideLastNames", "FeminicideVictimDocType", "FeminicideVictimDocNumber",
	}
	var aliasesSlice []string = []string{}

	var fieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasesSlice, FeminicideDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideDBScheme, FeminicideFieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + feminicidePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideDBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideDBScheme, FeminicideFieldDefinitions, true) +
		` ORDER BY ` + feminicidePath + `.` + FeminicideFieldDefinitions["FeminicideCreationDate"].DBName + ` ASC ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	var feminicides []FeminicideDTO
	for persistenceCtrl.Next() {
		var pg FeminicidePgDB
		persistenceCtrl.ScanRow(
			&pg.FeminicideId,
			&pg.FeminicideICode,
			&pg.FeminicideCreationDate,
			&pg.FeminicideUpdateDate,
			&pg.FeminicideStatus,
			&pg.FeminicideGeneralUser,
			&pg.FeminicideNames,
			&pg.FeminicideLastNames,
			&pg.FeminicideVictimDocType,
			&pg.FeminicideVictimDocNumber,
		)
		feminicides = append(feminicides, pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	if page == 0 {
		var countQuery string = `SELECT COUNT(*) FROM ` + feminicidePath + common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideDBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideDBScheme, FeminicideFieldDefinitions, true)
		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return feminicides, count, nil
}

func GetAllFeminicides(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideDTO, int, error) {
	// Declaración de variables locales.
	var count int
	var feminicide FeminicideDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicidePath string = FeminicideDBScheme + "." + FeminicideDBName

	// Configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var feminicideFieldsSlice []string = []string{"FeminicideICode", "FeminicideCreationDate", "FeminicideStatus", "FeminicideNames", "FeminicideLastNames", "FeminicideVictimDocType", "FeminicideVictimDocNumber"}
	var feminicideFieldsAliasSlice []string = []string{}

	// Construye el query SQL para seleccionar todos los registros de la tabla.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, feminicideFieldsSlice, feminicideFieldsAliasSlice, FeminicideDBName, []string{}, []string{}, []string{}, "", FeminicideDBScheme, FeminicideFieldDefinitions, true) +
		` ORDER BY ` + feminicidePath + `.` + FeminicideFieldDefinitions["FeminicideCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	// Ejecuta el query.
	persistenceCtrl.Query(context.Background(), query)
	var feminicides []FeminicideDTO
	// Itera sobre cada fila del resultado, escaneándola en la estructura FeminicideDTO.
	for persistenceCtrl.Next() {
		feminicide = FeminicideDTO{}
		persistenceCtrl.ScanRow(&feminicide.FeminicideICode, &feminicide.FeminicideCreationDate, &feminicide.FeminicideStatus,
			&feminicide.FeminicideNames, &feminicide.FeminicideLastNames, &feminicide.FeminicideVictimDocType, &feminicide.FeminicideVictimDocNumber)
		feminicides = append(feminicides, feminicide)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + feminicidePath

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return feminicides, count, nil
}

func UpdateFeminicide(feminicide *FeminicideDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideUpdateDate", "FeminicideStatus", "FeminicideGeneralUser",
		"FeminicideNames", "FeminicideLastNames", "FeminicideVictimDocType",
		"FeminicideVictimDocNumber",
	}
	var aliasesSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, fieldsSlice, aliasesSlice, FeminicideDBName, []string{"FeminicideId"}, []string{}, []string{}, common_dao.SQL_AND, FeminicideDBScheme, FeminicideFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query,
		feminicide.FeminicideId,
		feminicide.FeminicideUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicide.FeminicideStatus,
		feminicide.FeminicideGeneralUser,
		feminicide.FeminicideNames,
		feminicide.FeminicideLastNames,
		feminicide.FeminicideVictimDocType,
		feminicide.FeminicideVictimDocNumber,
	)

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

func UpdateFeminicideOwnersAndRolesByFeminicideId(feminicideId uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicidePath string = FeminicideDBScheme + "." + FeminicideDBName
	var relCaseOwnerFeminicidePath string = RelCaseOwnerFeminicideDBScheme + "." + RelCaseOwnerFeminicideDBName
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
		SELECT string_agg( '('||salvia.rel_case_owner_feminicide.rel_case_owner_feminicide_creation_date||') ' ||security.general_user_profile.general_user_profile_names || ' ' || security.general_user_profile.general_user_profile_last_names ||
			(
				SELECT string_agg( ' ['||security.role.role_name || '] ', ', ') as roles
				FROM security.role
				RIGHT JOIN security.rel_role_general_user ON (security.rel_role_general_user.role_id = security.role.role_id)
				LEFT JOIN security.general_user ON (security.rel_role_general_user.general_user_id = security.general_user.general_user_id)

				WHERE us.general_user_id = security.general_user.general_user_id

			), ', ' ORDER BY salvia.rel_case_owner_feminicide.rel_case_owner_feminicide_creation_date ASC) as case_owners
			FROM salvia.rel_case_owner_feminicide
			LEFT JOIN salvia.case_owner ON (salvia.case_owner.case_owner_id = salvia.rel_case_owner_feminicide.case_owner_id)
			LEFT JOIN security.general_user us ON (salvia.case_owner.case_owner_general_user = us.general_user_i_code)
			LEFT JOIN security.general_user_profile ON (security.general_user_profile.general_user_profile_id = us.general_user_general_user_profile)

			WHERE salvia.rel_case_owner_feminicide.feminicide_id = salvia.feminicide.feminicide_id AND salvia.feminicide.feminicide_id = $1
	*/

	var owners string = ` (SELECT string_agg( '('||` + relCaseOwnerFeminicidePath + `.` + RelCaseOwnerFeminicideFieldDefinitions["RelCaseOwnerFeminicide_CreationDate"].DBName + `||') ' ||` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileNames"].DBName + ` || ' ' || ` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileLastNames"].DBName + ` || ` +
		` (SELECT string_agg( ' ['||` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleName"].DBName + ` || '] ', ', ') as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` WHERE us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName +
		`), ', ' ORDER BY ` + relCaseOwnerFeminicidePath + `.` + RelCaseOwnerFeminicideFieldDefinitions["RelCaseOwnerFeminicide_CreationDate"].DBName + ` ASC) as case_owners ` +

		` FROM ` + relCaseOwnerFeminicidePath +
		` LEFT JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerFeminicidePath + `.` + RelCaseOwnerFeminicideFieldDefinitions["RelCaseOwnerFeminicide_CaseOwner"].DBName + `) ` +
		` LEFT JOIN ` + userPath + ` AS us ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `) ` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +

		` WHERE ` + relCaseOwnerFeminicidePath + `.` + RelCaseOwnerFeminicideFieldDefinitions["RelCaseOwnerFeminicide_Feminicide"].DBName + ` = ` + feminicidePath + `.` + FeminicideFieldDefinitions["FeminicideId"].DBName + ` AND ` +
		feminicidePath + `.` + FeminicideFieldDefinitions["FeminicideId"].DBName + ` = $1) WHERE ` + feminicidePath + `.` + FeminicideFieldDefinitions["FeminicideId"].DBName + ` = $1 `

	var query string = `UPDATE ` + feminicidePath +
		` SET ` + FeminicideFieldDefinitions["FeminicideOwnerDescription"].DBName + ` =  ` + owners + ` `

	persistenceCtrl.Exec(context.Background(), query, feminicideId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
	}
	return nil
}

func SetFeminicideDefaults(feminicide *FeminicideDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		feminicide.FeminicideCreationDate = time.Now()
		feminicide.FeminicideUpdateDate = time.Now()
		feminicide.FeminicideICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		feminicide.FeminicideUpdateDate = time.Now()
	}
}

func (obj *FeminicidePgDB) ToDTO() FeminicideDTO {
	var dto FeminicideDTO

	if obj.FeminicideId.Valid {
		dto.FeminicideId = uint64(obj.FeminicideId.Int64)
	}

	if obj.FeminicideICode.Valid {
		dto.FeminicideICode = obj.FeminicideICode.String
	}

	if obj.FeminicideCreationDate.Valid {
		dto.FeminicideCreationDate = obj.FeminicideCreationDate.Time
	}

	if obj.FeminicideUpdateDate.Valid {
		dto.FeminicideUpdateDate = obj.FeminicideUpdateDate.Time
	}

	if obj.FeminicideStatus.Valid {
		dto.FeminicideStatus = obj.FeminicideStatus.String
	}

	if obj.FeminicideGeneralUser.Valid {
		dto.FeminicideGeneralUser = obj.FeminicideGeneralUser.String
	}

	if obj.FeminicideNames.Valid {
		dto.FeminicideNames = obj.FeminicideNames.String
	}

	if obj.FeminicideLastNames.Valid {
		dto.FeminicideLastNames = obj.FeminicideLastNames.String
	}

	if obj.FeminicideVictimDocType.Valid {
		dto.FeminicideVictimDocType = obj.FeminicideVictimDocType.String
	}

	if obj.FeminicideVictimDocNumber.Valid {
		dto.FeminicideVictimDocNumber = obj.FeminicideVictimDocNumber.String
	}

	return dto
}
