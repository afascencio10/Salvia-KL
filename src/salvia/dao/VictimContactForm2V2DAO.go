// Package salvia_daos — DAO de la tabla de extensión victim_contact_form2_v2.
//
// PSWYA-755: Registro de Contacto Víctima V2 (Flujo A). Esta tabla existe
// para no tocar VictimContactForm2DAO.go ni la tabla victim_contact_form2 —
// solo contiene los campos que V1 no tiene (ver related-tables.md en
// DocsMD/Screens/Registro de Contacto Víctima V2/).
package salvia_daos

import (
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/utils"

	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"bitsflow/common/db"
)

var (
	VictimContactForm2V2EntityName string = "VictimContactForm2V2"
	VictimContactForm2V2JSONName   string = "form2V2"
	VictimContactForm2V2DBName     string = "victim_contact_form2_v2"
	VictimContactForm2V2DBScheme   string = "salvia"

	VictimContactForm2V2FieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimContactForm2V2Id":                      {Name: "VictimContactForm2V2Id", DBName: "victim_contact_form2_v2_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2V2ICode":                   {Name: "VictimContactForm2V2ICode", DBName: "victim_contact_form2_v2_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"VictimContactForm2V2CreationDate":            {Name: "VictimContactForm2V2CreationDate", DBName: "victim_contact_form2_v2_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2V2UpdateDate":              {Name: "VictimContactForm2V2UpdateDate", DBName: "victim_contact_form2_v2_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2V2VictimContactForm2":      {Name: "VictimContactForm2V2VictimContactForm2", DBName: "victim_contact_form2_v2_victim_contact_form2", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2V2AuthorizationAnswer":     {Name: "VictimContactForm2V2AuthorizationAnswer", DBName: "victim_contact_form2_v2_authorization_answer", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"VictimContactForm2V2Anonimato":               {Name: "VictimContactForm2V2Anonimato", DBName: "victim_contact_form2_v2_anonimato", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm2V2Parentesco":              {Name: "VictimContactForm2V2Parentesco", DBName: "victim_contact_form2_v2_parentesco", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm2V2ReporterLastNames":       {Name: "VictimContactForm2V2ReporterLastNames", DBName: "victim_contact_form2_v2_reporter_last_names", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 64, Required: false},
		"VictimContactForm2V2ReporterEmail":           {Name: "VictimContactForm2V2ReporterEmail", DBName: "victim_contact_form2_v2_reporter_email", Alias: "", ModelType: "string", MinSize: 0, MaxSize: 120, Required: false},
		"VictimContactForm2V2ReporterWillReceiveCall": {Name: "VictimContactForm2V2ReporterWillReceiveCall", DBName: "victim_contact_form2_v2_reporter_will_receive_call", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"VictimContactForm2V2ReporterBestContactTime": {Name: "VictimContactForm2V2ReporterBestContactTime", DBName: "victim_contact_form2_v2_reporter_best_contact_time", Alias: "", ModelType: "time", MinSize: 0, MaxSize: 0, Required: false},
	}
)

// ---------------------------------------------------------------------------
// DTO – estructura enviada/recibida vía JSON
// ---------------------------------------------------------------------------
type VictimContactForm2V2DTO struct {
	VictimContactForm2V2Id                      uint64                  `json:"-"`
	VictimContactForm2V2ICode                   string                  `json:"icode"`
	VictimContactForm2V2CreationDate            time.Time               `json:"creationDate"`
	VictimContactForm2V2UpdateDate              time.Time               `json:"updateDate"`
	VictimContactForm2V2VictimContactForm2      interface{}             `json:"-"`
	VictimContactForm2V2AuthorizationAnswer     VictimCaseForm2EnumsDTO `json:"authorizationAnswer"`
	VictimContactForm2V2Anonimato               VictimCaseForm2EnumsDTO `json:"anonimato"`
	VictimContactForm2V2Parentesco              VictimCaseForm2EnumsDTO `json:"parentesco"`
	VictimContactForm2V2ReporterLastNames       string                  `json:"reporterLastNames"`
	VictimContactForm2V2ReporterEmail           string                  `json:"reporterEmail"`
	VictimContactForm2V2ReporterWillReceiveCall VictimCaseForm2EnumsDTO `json:"reporterWillReceiveCall"`
	VictimContactForm2V2ReporterBestContactTime time.Time               `json:"reporterBestContactTime"`
}

// ---------------------------------------------------------------------------
// PgDB – estructura que mapea directamente las filas de PostgreSQL
// ---------------------------------------------------------------------------
type VictimContactForm2V2PgDB struct {
	VictimContactForm2V2Id                      sql.NullString
	VictimContactForm2V2ICode                   sql.NullString
	VictimContactForm2V2CreationDate            sql.NullString
	VictimContactForm2V2UpdateDate              sql.NullString
	VictimContactForm2V2VictimContactForm2      sql.NullString
	VictimContactForm2V2AuthorizationAnswer     sql.NullString
	VictimContactForm2V2Anonimato               sql.NullString
	VictimContactForm2V2Parentesco              sql.NullString
	VictimContactForm2V2ReporterLastNames       sql.NullString
	VictimContactForm2V2ReporterEmail           sql.NullString
	VictimContactForm2V2ReporterWillReceiveCall sql.NullString
	VictimContactForm2V2ReporterBestContactTime sql.NullString
}

func (vcd VictimContactForm2V2DTO) MarshalJSON() ([]byte, error) {
	type Alias VictimContactForm2V2DTO

	return json.Marshal(&struct {
		*Alias
		VictimContactForm2V2CreationDate            string `json:"creationDate"`
		VictimContactForm2V2UpdateDate              string `json:"updateDate"`
		VictimContactForm2V2ReporterBestContactTime string `json:"reporterBestContactTime"`
	}{
		Alias:                            (*Alias)(&vcd),
		VictimContactForm2V2CreationDate: vcd.VictimContactForm2V2CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimContactForm2V2UpdateDate:   vcd.VictimContactForm2V2UpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		VictimContactForm2V2ReporterBestContactTime: func() string {
			if vcd.VictimContactForm2V2ReporterBestContactTime.IsZero() {
				return ""
			}
			return vcd.VictimContactForm2V2ReporterBestContactTime.Format(common_config.DateTime.TIME_FORMAT)
		}(),
	})
}

// SetVictimContactForm2Defaults ya existe para VictimContactForm2 — este es
// el equivalente para la fila de extensión V2 (icode + fechas al insertar).
func SetVictimContactForm2V2Defaults(v2 *VictimContactForm2V2DTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		v2.VictimContactForm2V2CreationDate = time.Now()
		v2.VictimContactForm2V2UpdateDate = time.Now()
		v2.VictimContactForm2V2ICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		v2.VictimContactForm2V2UpdateDate = time.Now()
	}
}

// SetVictimContactForm2V2 inserta la fila de extensión, ligada por FK al
// VictimContactForm2 ya creado (form2.VictimContactForm2Id debe estar
// resuelto — es decir, llamar después de salvia_daos.SetVictimContactForm2).
func SetVictimContactForm2V2(v2 *VictimContactForm2V2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimContactForm2V2ICode",
		"VictimContactForm2V2CreationDate",
		"VictimContactForm2V2UpdateDate",
		"VictimContactForm2V2VictimContactForm2",
		"VictimContactForm2V2AuthorizationAnswer",
		"VictimContactForm2V2Anonimato",
		"VictimContactForm2V2Parentesco",
		"VictimContactForm2V2ReporterLastNames",
		"VictimContactForm2V2ReporterEmail",
		"VictimContactForm2V2ReporterWillReceiveCall",
		"VictimContactForm2V2ReporterBestContactTime",
	}

	var aliasSlice []string = []string{}

	query := common_dao.GetSQL(common_dao.SQL_INSERT, fieldsSlice, aliasSlice, VictimContactForm2V2DBName, []string{}, []string{}, []string{"VictimContactForm2V2Id"}, common_dao.SQL_AND, VictimContactForm2V2DBScheme, VictimContactForm2V2FieldDefinitions, false)

	// TIME WITHOUT TIME ZONE nullable: solo se formatea si el reportante
	// diligenció hora de recontacto (Familiar/Hombre); si no, va NULL.
	var reporterBestContactTime interface{} = nil
	if !v2.VictimContactForm2V2ReporterBestContactTime.IsZero() {
		reporterBestContactTime = v2.VictimContactForm2V2ReporterBestContactTime.Format(common_config.DateTime.TIME_FORMAT)
	}

	persistenceCtrl.QueryRow(context.Background(), query,
		v2.VictimContactForm2V2ICode,
		v2.VictimContactForm2V2CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		v2.VictimContactForm2V2UpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		v2.VictimContactForm2V2VictimContactForm2.(VictimContactForm2DTO).VictimContactForm2Id,
		utils.NilIfZero(v2.VictimContactForm2V2AuthorizationAnswer.VictimCaseForm2EnumsId),
		utils.NilIfZero(v2.VictimContactForm2V2Anonimato.VictimCaseForm2EnumsId),
		utils.NilIfZero(v2.VictimContactForm2V2Parentesco.VictimCaseForm2EnumsId),
		utils.NilIfZero(v2.VictimContactForm2V2ReporterLastNames),
		utils.NilIfZero(v2.VictimContactForm2V2ReporterEmail),
		utils.NilIfZero(v2.VictimContactForm2V2ReporterWillReceiveCall.VictimCaseForm2EnumsId),
		reporterBestContactTime)

	persistenceCtrl.Scan(&v2.VictimContactForm2V2Id)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetVictimContactForm2V2 obtiene la fila de extensión asociada a un
// VictimContactForm2 (usada por el Expediente/futuras pantallas de agente).
func GetVictimContactForm2V2(by common_controllers.By, v2 *VictimContactForm2V2DTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {

	var persistenceCtrl common_controllers.PersistenceController
	path := VictimContactForm2V2DBScheme + "." + VictimContactForm2V2DBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	fieldsSlice := []string{
		"VictimContactForm2V2Id",
		"VictimContactForm2V2ICode",
		"VictimContactForm2V2CreationDate",
		"VictimContactForm2V2UpdateDate",
		"VictimContactForm2V2VictimContactForm2",
		"VictimContactForm2V2AuthorizationAnswer",
		"VictimContactForm2V2Anonimato",
		"VictimContactForm2V2Parentesco",
		"VictimContactForm2V2ReporterLastNames",
		"VictimContactForm2V2ReporterEmail",
		"VictimContactForm2V2ReporterWillReceiveCall",
		"VictimContactForm2V2ReporterBestContactTime",
	}

	var aliasSlice []string = []string{}

	fieldsStr := common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fieldsSlice, aliasSlice, VictimContactForm2V2DBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimContactForm2V2DBScheme, VictimContactForm2V2FieldDefinitions, true)

	query := `SELECT ` + fieldsStr +
		` FROM ` + path +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, VictimContactForm2V2DBName, by.AttrsName, nil, nil, by.Operator, VictimContactForm2V2DBScheme, VictimContactForm2V2FieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var pg VictimContactForm2V2PgDB
	persistenceCtrl.Scan(
		&pg.VictimContactForm2V2Id,
		&pg.VictimContactForm2V2ICode,
		&pg.VictimContactForm2V2CreationDate,
		&pg.VictimContactForm2V2UpdateDate,
		&pg.VictimContactForm2V2VictimContactForm2,
		&pg.VictimContactForm2V2AuthorizationAnswer,
		&pg.VictimContactForm2V2Anonimato,
		&pg.VictimContactForm2V2Parentesco,
		&pg.VictimContactForm2V2ReporterLastNames,
		&pg.VictimContactForm2V2ReporterEmail,
		&pg.VictimContactForm2V2ReporterWillReceiveCall,
		&pg.VictimContactForm2V2ReporterBestContactTime)

	if persistenceCtrl.Error != nil {
		return persistenceCtrl.Error
	}

	*v2 = pg.ToDTO()
	return nil
}

func (obj *VictimContactForm2V2PgDB) ToDTO() VictimContactForm2V2DTO {
	dto := VictimContactForm2V2DTO{}

	if obj.VictimContactForm2V2Id.Valid {
		dto.VictimContactForm2V2Id, _ = strconv.ParseUint(obj.VictimContactForm2V2Id.String, 10, 64)
	}
	if obj.VictimContactForm2V2ICode.Valid {
		dto.VictimContactForm2V2ICode = obj.VictimContactForm2V2ICode.String
	}
	if obj.VictimContactForm2V2ReporterLastNames.Valid {
		dto.VictimContactForm2V2ReporterLastNames = obj.VictimContactForm2V2ReporterLastNames.String
	}
	if obj.VictimContactForm2V2ReporterEmail.Valid {
		dto.VictimContactForm2V2ReporterEmail = obj.VictimContactForm2V2ReporterEmail.String
	}
	if obj.VictimContactForm2V2Anonimato.Valid {
		id, _ := strconv.ParseUint(obj.VictimContactForm2V2Anonimato.String, 10, 64)
		dto.VictimContactForm2V2Anonimato = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: id}
		GetLocalVictimCaseForm2EnumsById(&dto.VictimContactForm2V2Anonimato)
	}
	if obj.VictimContactForm2V2Parentesco.Valid {
		id, _ := strconv.ParseUint(obj.VictimContactForm2V2Parentesco.String, 10, 64)
		dto.VictimContactForm2V2Parentesco = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: id}
		GetLocalVictimCaseForm2EnumsById(&dto.VictimContactForm2V2Parentesco)
	}
	if obj.VictimContactForm2V2AuthorizationAnswer.Valid {
		id, _ := strconv.ParseUint(obj.VictimContactForm2V2AuthorizationAnswer.String, 10, 64)
		dto.VictimContactForm2V2AuthorizationAnswer = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: id}
		GetLocalVictimCaseForm2EnumsById(&dto.VictimContactForm2V2AuthorizationAnswer)
	}
	if obj.VictimContactForm2V2ReporterWillReceiveCall.Valid {
		id, _ := strconv.ParseUint(obj.VictimContactForm2V2ReporterWillReceiveCall.String, 10, 64)
		dto.VictimContactForm2V2ReporterWillReceiveCall = VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsId: id}
		GetLocalVictimCaseForm2EnumsById(&dto.VictimContactForm2V2ReporterWillReceiveCall)
	}

	return dto
}
