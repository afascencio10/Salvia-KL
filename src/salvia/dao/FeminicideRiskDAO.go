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
	FeminicideRiskEntityName string = "FeminicideRisk"
	FeminicideRiskJSONName   string = "feminicideRisk"
	FeminicideRiskDBName     string = "feminicide_risk"
	FeminicideRiskDBScheme   string = "public" // Cambiar si está en otro esquema

	FeminicideRiskFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FeminicideRiskId":              {Name: "FeminicideRiskId", DBName: "feminicide_risk_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskICode":           {Name: "FeminicideRiskICode", DBName: "feminicide_risk_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"FeminicideRiskCreationDate":    {Name: "FeminicideRiskCreationDate", DBName: "feminicide_risk_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskUpdateDate":      {Name: "FeminicideRiskUpdateDate", DBName: "feminicide_risk_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"FeminicideRiskStatus":          {Name: "FeminicideRiskStatus", DBName: "feminicide_risk_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FeminicideRiskGeneralUser":     {Name: "FeminicideRiskGeneralUser", DBName: "feminicide_risk_general_user", Alias: "", ModelType: "string", MinSize: 36, MaxSize: 36, Required: true},
		"FeminicideRiskNames":           {Name: "FeminicideRiskNames", DBName: "feminicide_risk_names", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: true},
		"FeminicideRiskLastNames":       {Name: "FeminicideRiskLastNames", DBName: "feminicide_risk_last_names", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: true},
		"FeminicideRiskVictimDocType":   {Name: "FeminicideRiskVictimDocType", DBName: "feminicide_risk_victim_doc_type", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: true},
		"FeminicideRiskVictimDocNumber": {Name: "FeminicideRiskVictimDocNumber", DBName: "feminicide_risk_victim_doc_number", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 32, Required: true},
	}
)

type FeminicideRiskDTO struct {
	FeminicideRiskId              uint64    `json:"-"`
	FeminicideRiskICode           string    `json:"icode"`
	FeminicideRiskCreationDate    time.Time `json:"creationDate"`
	FeminicideRiskUpdateDate      time.Time `json:"updateDate"`
	FeminicideRiskStatus          string    `json:"status"`
	FeminicideRiskGeneralUser     string    `json:"generalUser"`
	FeminicideRiskNames           string    `json:"names"`
	FeminicideRiskLastNames       string    `json:"lastNames"`
	FeminicideRiskVictimDocType   string    `json:"victimDocType"`
	FeminicideRiskVictimDocNumber string    `json:"victimDocNumber"`

	FeminicideRiskForm1      FeminicideRiskForm1DTO `json:"form"`
	FeminicideRiskVictimCase VictimCaseDTO          `json:"victimCase"`
}

type FeminicideRiskPgDB struct {
	FeminicideRiskId              sql.NullInt64
	FeminicideRiskICode           sql.NullString
	FeminicideRiskCreationDate    sql.NullTime
	FeminicideRiskUpdateDate      sql.NullTime
	FeminicideRiskStatus          sql.NullString
	FeminicideRiskGeneralUser     sql.NullString
	FeminicideRiskNames           sql.NullString
	FeminicideRiskLastNames       sql.NullString
	FeminicideRiskVictimDocType   sql.NullString
	FeminicideRiskVictimDocNumber sql.NullString
}

// Manejo para las fechas y tiempos
func (frd FeminicideRiskDTO) MarshalJSON() ([]byte, error) {
	type Alias FeminicideRiskDTO

	return json.Marshal(&struct {
		*Alias
		FeminicideRiskCreationDate string `json:"creationDate"`
		FeminicideRiskUpdateDate   string `json:"updateDate"`
	}{
		Alias:                      (*Alias)(&frd),
		FeminicideRiskCreationDate: frd.FeminicideRiskCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FeminicideRiskUpdateDate:   frd.FeminicideRiskUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (frd *FeminicideRiskDTO) UnmarshalJSON(data []byte) error {
	type Alias FeminicideRiskDTO

	aux := &struct {
		*Alias
		FeminicideRiskCreationDate string `json:"creationDate"`
		FeminicideRiskUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(frd),
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

	frd.FeminicideRiskCreationDate = parse(aux.FeminicideRiskCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	frd.FeminicideRiskUpdateDate = parse(aux.FeminicideRiskUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

func SetFeminicideRisk(feminicideRisk *FeminicideRiskDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideRiskICode", "FeminicideRiskCreationDate", "FeminicideRiskUpdateDate",
		"FeminicideRiskStatus", "FeminicideRiskGeneralUser", "FeminicideRiskNames",
		"FeminicideRiskLastNames", "FeminicideRiskVictimDocType", "FeminicideRiskVictimDocNumber",
	}

	var query string = common_dao.GetSQL(
		common_dao.SQL_INSERT,
		fieldsSlice, []string{},
		FeminicideRiskDBName, []string{}, []string{}, []string{"FeminicideRiskId"},
		common_dao.SQL_AND, FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		feminicideRisk.FeminicideRiskICode,
		feminicideRisk.FeminicideRiskCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideRisk.FeminicideRiskUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideRisk.FeminicideRiskStatus,
		feminicideRisk.FeminicideRiskGeneralUser,
		feminicideRisk.FeminicideRiskNames,
		feminicideRisk.FeminicideRiskLastNames,
		feminicideRisk.FeminicideRiskVictimDocType,
		feminicideRisk.FeminicideRiskVictimDocNumber)

	persistenceCtrl.Scan(&feminicideRisk.FeminicideRiskId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicideRisk(by common_controllers.By, feminicideRisk *FeminicideRiskDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var path string = FeminicideRiskDBScheme + "." + FeminicideRiskDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideRiskId", "FeminicideRiskICode", "FeminicideRiskCreationDate", "FeminicideRiskUpdateDate",
		"FeminicideRiskStatus", "FeminicideRiskGeneralUser", "FeminicideRiskNames",
		"FeminicideRiskLastNames", "FeminicideRiskVictimDocType", "FeminicideRiskVictimDocNumber",
	}

	var fieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, []string{}, FeminicideRiskDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideRiskDBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, true)

	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var pgDB FeminicideRiskPgDB
	persistenceCtrl.Scan(
		&pgDB.FeminicideRiskId,
		&pgDB.FeminicideRiskICode,
		&pgDB.FeminicideRiskCreationDate,
		&pgDB.FeminicideRiskUpdateDate,
		&pgDB.FeminicideRiskStatus,
		&pgDB.FeminicideRiskGeneralUser,
		&pgDB.FeminicideRiskNames,
		&pgDB.FeminicideRiskLastNames,
		&pgDB.FeminicideRiskVictimDocType,
		&pgDB.FeminicideRiskVictimDocNumber)

	*feminicideRisk = pgDB.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

func GetFeminicideRisks(by common_controllers.By, page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideRiskDTO, int, error) {
	var count int

	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicideRiskPath string = FeminicideRiskDBScheme + "." + FeminicideRiskDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{"FeminicideRiskId",
		"FeminicideRiskICode", "FeminicideRiskCreationDate", "FeminicideRiskUpdateDate",
		"FeminicideRiskStatus", "FeminicideRiskGeneralUser", "FeminicideRiskNames",
		"FeminicideRiskLastNames", "FeminicideRiskVictimDocType", "FeminicideRiskVictimDocNumber",
	}
	var aliasesSlice []string = []string{}

	var fieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasesSlice, FeminicideRiskDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, true)

	var query string = `SELECT ` + fieldsStr +
		` FROM ` + feminicideRiskPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideRiskDBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, true) +
		` ORDER BY ` + feminicideRiskPath + `.` + FeminicideRiskFieldDefinitions["FeminicideRiskCreationDate"].DBName + ` ASC ` +
		common_dao.GetOffsetQuery(page)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	var feminicidesRisk []FeminicideRiskDTO
	for persistenceCtrl.Next() {
		var pg FeminicideRiskPgDB
		persistenceCtrl.ScanRow(
			&pg.FeminicideRiskId,
			&pg.FeminicideRiskICode,
			&pg.FeminicideRiskCreationDate,
			&pg.FeminicideRiskUpdateDate,
			&pg.FeminicideRiskStatus,
			&pg.FeminicideRiskGeneralUser,
			&pg.FeminicideRiskNames,
			&pg.FeminicideRiskLastNames,
			&pg.FeminicideRiskVictimDocType,
			&pg.FeminicideRiskVictimDocNumber,
		)
		feminicidesRisk = append(feminicidesRisk, pg.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	if page == 0 {
		var countQuery string = `SELECT COUNT(*) FROM ` + feminicideRiskPath + common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FeminicideDBName, by.AttrsName, []string{}, []string{}, by.Operator, FeminicideDBScheme, FeminicideFieldDefinitions, true)
		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return feminicidesRisk, count, nil
}

func GetAllFeminicideRisks(page int, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FeminicideRiskDTO, int, error) {
	// Declaración de variables locales.
	var count int
	var feminicideRisk FeminicideRiskDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicideRiskPath string = FeminicideRiskDBScheme + "." + FeminicideRiskDBName

	// Configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var feminicideRiskFieldsSlice []string = []string{"FeminicideRiskICode", "FeminicideRiskCreationDate", "FeminicideRiskStatus", "FeminicideRiskNames", "FeminicideRiskLastNames", "FeminicideRiskVictimDocType", "FeminicideRiskVictimDocNumber"}
	var feminicideRiskFieldsAliasSlice []string = []string{}

	// Construye el query SQL para seleccionar todos los registros de la tabla.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, feminicideRiskFieldsSlice, feminicideRiskFieldsAliasSlice, FeminicideRiskDBName, []string{}, []string{}, []string{}, "", FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, true) +
		` ORDER BY ` + feminicideRiskPath + `.` + FeminicideRiskFieldDefinitions["FeminicideRiskCreationDate"].DBName + ` ASC  ` +
		common_dao.GetOffsetQuery(page)

	// Ejecuta el query.
	persistenceCtrl.Query(context.Background(), query)
	var feminicideRisks []FeminicideRiskDTO
	// Itera sobre cada fila del resultado, escaneándola en la estructura FeminicideRiskDTO.
	for persistenceCtrl.Next() {
		feminicideRisk = FeminicideRiskDTO{}
		persistenceCtrl.ScanRow(&feminicideRisk.FeminicideRiskICode, &feminicideRisk.FeminicideRiskCreationDate, &feminicideRisk.FeminicideRiskStatus,
			&feminicideRisk.FeminicideRiskNames, &feminicideRisk.FeminicideRiskLastNames, &feminicideRisk.FeminicideRiskVictimDocType, &feminicideRisk.FeminicideRiskVictimDocNumber)
		feminicideRisks = append(feminicideRisks, feminicideRisk)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, count, persistenceCtrl.Error
	}

	// Se hace el conteo de elementos totales para la primera consulta
	if page == 0 {
		var countQuery string = `SELECT COUNT(*) ` +
			` FROM ` + feminicideRiskPath

		persistenceCtrl.QueryRow(context.Background(), countQuery)
		persistenceCtrl.Scan(&count)
	}

	return feminicideRisks, count, nil
}

func UpdateFeminicideRisk(feminicideRisk *FeminicideRiskDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var fieldsSlice []string = []string{
		"FeminicideRiskUpdateDate", "FeminicideRiskStatus", "FeminicideRiskGeneralUser",
		"FeminicideRiskNames", "FeminicideRiskLastNames", "FeminicideRiskVictimDocType",
		"FeminicideRiskVictimDocNumber",
	}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, fieldsSlice, []string{}, FeminicideRiskDBName, []string{"FeminicideRiskId"}, []string{}, []string{}, common_dao.SQL_AND, FeminicideRiskDBScheme, FeminicideRiskFieldDefinitions, false)
	persistenceCtrl.Exec(context.Background(), query,
		feminicideRisk.FeminicideRiskId,
		feminicideRisk.FeminicideRiskUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		feminicideRisk.FeminicideRiskStatus,
		feminicideRisk.FeminicideRiskGeneralUser,
		feminicideRisk.FeminicideRiskNames,
		feminicideRisk.FeminicideRiskLastNames,
		feminicideRisk.FeminicideRiskVictimDocType,
		feminicideRisk.FeminicideRiskVictimDocNumber)

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

func UpdateFeminicideRiskOwnersAndRolesByFeminicideRiskId(feminicideId uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Definición de variables
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var feminicideRiskPath string = FeminicideRiskDBScheme + "." + FeminicideRiskDBName
	var relCaseOwnerFeminicideRiskPath string = RelCaseOwnerFeminicideRiskDBScheme + "." + RelCaseOwnerFeminicideRiskDBName
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
		SELECT string_agg( '('||salvia.rel_case_owner_feminicide.rel_case_owner_feminicide_risk_creation_date||') ' ||security.general_user_profile.general_user_profile_names || ' ' || security.general_user_profile.general_user_profile_last_names ||
			(
				SELECT string_agg( ' ['||security.role.role_name || '] ', ', ') as roles
				FROM security.role
				RIGHT JOIN security.rel_role_general_user ON (security.rel_role_general_user.role_id = security.role.role_id)
				LEFT JOIN security.general_user ON (security.rel_role_general_user.general_user_id = security.general_user.general_user_id)

				WHERE us.general_user_id = security.general_user.general_user_id

			), ', ' ORDER BY salvia.rel_case_owner_feminicide.rel_case_owner_feminicide_risk_creation_date ASC) as case_owners
			FROM salvia.rel_case_owner_feminicide
			LEFT JOIN salvia.case_owner ON (salvia.case_owner.case_owner_id = salvia.rel_case_owner_feminicide.case_owner_id)
			LEFT JOIN security.general_user us ON (salvia.case_owner.case_owner_general_user = us.general_user_i_code)
			LEFT JOIN security.general_user_profile ON (security.general_user_profile.general_user_profile_id = us.general_user_general_user_profile)

			WHERE salvia.rel_case_owner_feminicide.feminicide_risk_id = salvia.feminicide.feminicide_risk_id AND salvia.feminicide.feminicide_risk_id = $1
	*/

	var owners string = ` (SELECT string_agg( '('||` + relCaseOwnerFeminicideRiskPath + `.` + RelCaseOwnerFeminicideRiskFieldDefinitions["RelCaseOwnerFeminicideRisk_CreationDate"].DBName + `||') ' ||` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileNames"].DBName + ` || ' ' || ` +
		profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileLastNames"].DBName + ` || ` +
		` (SELECT string_agg( ' ['||` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleName"].DBName + ` || '] ', ', ') as roles` +
		` FROM ` + rolePath +
		` RIGHT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + security_daos.RoleFieldDefinitions["RoleId"].DBName + `)` +
		` LEFT JOIN ` + userPath + ` ON (` + relRolePath + `.` + security_daos.RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		` WHERE us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserId"].DBName +
		`), ', ' ORDER BY ` + relCaseOwnerFeminicideRiskPath + `.` + RelCaseOwnerFeminicideRiskFieldDefinitions["RelCaseOwnerFeminicideRisk_CreationDate"].DBName + ` ASC) as case_owners ` +

		` FROM ` + relCaseOwnerFeminicideRiskPath +
		` LEFT JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerFeminicideRiskPath + `.` + RelCaseOwnerFeminicideRiskFieldDefinitions["RelCaseOwnerFeminicideRisk_CaseOwner"].DBName + `) ` +
		` LEFT JOIN ` + userPath + ` AS us ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `) ` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = us.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +

		` WHERE ` + relCaseOwnerFeminicideRiskPath + `.` + RelCaseOwnerFeminicideRiskFieldDefinitions["RelCaseOwnerFeminicideRisk_FeminicideRisk"].DBName + ` = ` + feminicideRiskPath + `.` + FeminicideRiskFieldDefinitions["FeminicideRiskId"].DBName + ` AND ` +
		feminicideRiskPath + `.` + FeminicideRiskFieldDefinitions["FeminicideRiskId"].DBName + ` = $1) WHERE ` + feminicideRiskPath + `.` + FeminicideRiskFieldDefinitions["FeminicideRiskId"].DBName + ` = $1 `

	var query string = `UPDATE ` + feminicideRiskPath +
		` SET ` + FeminicideRiskFieldDefinitions["FeminicideRiskOwnerDescription"].DBName + ` =  ` + owners + ` `

	persistenceCtrl.Exec(context.Background(), query, feminicideId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
	}
	return nil
}

func SetFeminicideRiskDefaults(feminicideRisk *FeminicideRiskDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		feminicideRisk.FeminicideRiskCreationDate = time.Now()
		feminicideRisk.FeminicideRiskUpdateDate = time.Now()
		feminicideRisk.FeminicideRiskICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		feminicideRisk.FeminicideRiskUpdateDate = time.Now()
	}
}

func (obj *FeminicideRiskPgDB) ToDTO() FeminicideRiskDTO {
	var dto FeminicideRiskDTO

	if obj.FeminicideRiskId.Valid {
		dto.FeminicideRiskId = uint64(obj.FeminicideRiskId.Int64)
	}

	if obj.FeminicideRiskICode.Valid {
		dto.FeminicideRiskICode = obj.FeminicideRiskICode.String
	}

	if obj.FeminicideRiskCreationDate.Valid {
		dto.FeminicideRiskCreationDate = obj.FeminicideRiskCreationDate.Time
	}

	if obj.FeminicideRiskUpdateDate.Valid {
		dto.FeminicideRiskUpdateDate = obj.FeminicideRiskUpdateDate.Time
	}

	if obj.FeminicideRiskStatus.Valid {
		dto.FeminicideRiskStatus = obj.FeminicideRiskStatus.String
	}

	if obj.FeminicideRiskGeneralUser.Valid {
		dto.FeminicideRiskGeneralUser = obj.FeminicideRiskGeneralUser.String
	}

	if obj.FeminicideRiskNames.Valid {
		dto.FeminicideRiskNames = obj.FeminicideRiskNames.String
	}

	if obj.FeminicideRiskLastNames.Valid {
		dto.FeminicideRiskLastNames = obj.FeminicideRiskLastNames.String
	}

	if obj.FeminicideRiskVictimDocType.Valid {
		dto.FeminicideRiskVictimDocType = obj.FeminicideRiskVictimDocType.String
	}

	if obj.FeminicideRiskVictimDocNumber.Valid {
		dto.FeminicideRiskVictimDocNumber = obj.FeminicideRiskVictimDocNumber.String
	}

	return dto
}
