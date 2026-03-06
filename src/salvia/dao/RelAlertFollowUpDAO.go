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
	// RelAlertFollowUpEntityName es el nombre lógico de la entidad de relación entre Alerta y Seguimiento.
	RelAlertFollowUpEntityName string = "RelAlertFollowUp"
	// RelAlertFollowUpJSONName es el nombre JSON para la entidad RelAlertFollowUp.
	RelAlertFollowUpJSONName string = "relAlertFollowUp"
	// RelAlertFollowUpDBName es el nombre de la tabla en la base de datos para RelAlertFollowUp.
	RelAlertFollowUpDBName string = "rel_alert_follow_up"
	// RelAlertFollowUpDBScheme es el esquema de la base de datos para la tabla RelAlertFollowUp.
	RelAlertFollowUpDBScheme string = "salvia"

	// RelAlertFollowUpFieldDefinitions define las propiedades de los campos de la entidad RelAlertFollowUp,
	// incluyendo su nombre, nombre en la DB, alias, tipo de modelo, tamaño máximo y si es requerido.
	RelAlertFollowUpFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelAlertFollowUpId":           {Name: "RelAlertFollowUpId", DBName: "rel_alert_follow_up_id", Alias: "", ModelType: "uint", Required: false},
		"RelAlertFollowUpData":         {Name: "RelAlertFollowUpData", DBName: "rel_alert_follow_up_data", Alias: "", ModelType: "string", MaxSize: 128, Required: true},
		"RelAlertFollowUpCreationDate": {Name: "RelAlertFollowUpCreationDate", DBName: "rel_alert_follow_up_creation_date", Alias: "", ModelType: "datetime", Required: false},
		"RelAlertFollowUpFollowUp":     {Name: "RelAlertFollowUpFollowUp", DBName: "rel_alert_follow_up_follow_up", Alias: "", ModelType: "uint", Required: true},
		"RelAlertFollowUpAlert":        {Name: "RelAlertFollowUpAlert", DBName: "rel_alert_follow_up_alert", Alias: "", ModelType: "uint", Required: true},
	}
)

// RelAlertFollowUpDTO representa el Data Transfer Object para la relación entre una Alerta y un Seguimiento.
// Se utiliza para transferir datos entre las capas de la aplicación y la presentación.
type RelAlertFollowUpDTO struct {
	RelAlertFollowUpId           uint64      `json:"-"`            // ID de la relación (no serializado a JSON).
	RelAlertFollowUpData         string      `json:"data"`         // Datos adicionales asociados a la relación.
	RelAlertFollowUpCreationDate time.Time   `json:"creationDate"` // Fecha de creación de la relación.
	RelAlertFollowUpFollowUp     FollowUpDTO `json:"followUp"`     // Objeto DTO del seguimiento asociado.
	RelAlertFollowUpAlert        AlertDTO    `json:"alert"`        // Objeto DTO de la alerta asociada.
}

// RelAlertFollowUpPgDB representa la estructura de la relación entre Alerta y Seguimiento tal como se almacena en PostgreSQL.
// Utiliza tipos sql.Null para manejar valores nulos de la base de datos de forma segura.
type RelAlertFollowUpPgDB struct {
	RelAlertFollowUpId           sql.NullInt64  // ID de la relación.
	RelAlertFollowUpData         sql.NullString // Datos adicionales de la relación.
	RelAlertFollowUpCreationDate sql.NullTime   // Fecha de creación de la relación.
	RelAlertFollowUpFollowUp     sql.NullInt64  // ID del seguimiento asociado.
	RelAlertFollowUpAlert        sql.NullInt64  // ID de la alerta asociada.
}

// SetRelAlertFollowUp inserta una nueva relación entre una alerta y un seguimiento en la base de datos.
// Parámetros:
//   - rel: puntero al DTO que contiene los datos de la relación a insertar.
//   - connData: datos de la conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de producirse algún fallo.
func SetRelAlertFollowUp(rel *RelAlertFollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelAlertFollowUpData", "RelAlertFollowUpCreationDate", "RelAlertFollowUpFollowUp", "RelAlertFollowUpAlert"}
	var relFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, relFieldsSlice, relFieldsAliasSlice, RelAlertFollowUpDBName, []string{}, []string{}, []string{"RelAlertFollowUpId"}, common_dao.SQL_AND, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query, rel.RelAlertFollowUpData, rel.RelAlertFollowUpCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), rel.RelAlertFollowUpFollowUp.FollowUpId, rel.RelAlertFollowUpAlert.AlertId)

	persistenceCtrl.Scan(&rel.RelAlertFollowUpId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// RemoveRelAlertFollowUp elimina una relación específica entre una alerta y un seguimiento en la base de datos.
// Parámetros:
//   - rel: puntero al DTO que contiene el ID de la relación a eliminar.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func RemoveRelAlertFollowUp(rel *RelAlertFollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{}
	var relFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, relFieldsSlice, relFieldsAliasSlice, RelAlertFollowUpDBName, []string{}, []string{"RelAlertFollowUpId"}, []string{}, common_dao.SQL_AND, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, rel.RelAlertFollowUpId)

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

// RemoveRelAlertFollowUps elimina una o más relaciones entre alertas y seguimientos basándose en criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (nombres de atributos, alias, operadores y valores).
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func RemoveRelAlertFollowUps(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, RelAlertFollowUpDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelAlertFollowUp obtiene una relación específica entre una alerta y un seguimiento según los criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (atributos, alias, operador y valores).
//   - rel: puntero al DTO donde se almacenarán los datos obtenidos.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func GetRelAlertFollowUp(by common_controllers.By, rel *RelAlertFollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var relPath string = RelAlertFollowUpDBScheme + "." + RelAlertFollowUpDBName
	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName
	var alertPath string = AlertDBScheme + "." + AlertDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelAlertFollowUpId", "RelAlertFollowUpData", "RelAlertFollowUpCreationDate", "RelAlertFollowUpFollowUp", "RelAlertFollowUpAlert"}
	var relFieldsAliasSlice []string = []string{}

	var followUpFieldsSlice []string = []string{"FollowUpICode", "FollowUpCreationDate", "FollowUpUpdateDate", "FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation", "FollowUpAssaultWithWeapon", "FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner", "FollowUpViolenceHistoryWithOthers", "FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification", "FollowUpVictimVulnerability", "FollowUpRiskLevel", "FollowUpOwnerGeneralUser", "FollowUpStatus"}
	var followUpFieldsAliasSlice []string = []string{}

	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}

	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelAlertFollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, true)
	var followUpFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, followUpFieldsSlice, followUpFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, true)
	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)

	var query string = `SELECT ` + relFieldsStr + `, ` + followUpFieldsStr + `, ` + alertFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + followUpPath + ` ON (` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + ` = ` + relPath + `.` + RelAlertFollowUpFieldDefinitions["RelAlertFollowUpFollowUp"].DBName + `)` +
		` LEFT JOIN ` + alertPath + ` ON (` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + ` = ` + relPath + `.` + RelAlertFollowUpFieldDefinitions["RelAlertFollowUpAlert"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelAlertFollowUpDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var followUp FollowUpPgDB = FollowUpPgDB{}
	var alert AlertPgDB = AlertPgDB{}

	persistenceCtrl.Scan(&rel.RelAlertFollowUpId, &rel.RelAlertFollowUpData, &rel.RelAlertFollowUpCreationDate,
		&followUp.FollowUpId, &alert.AlertId,
		&followUp.FollowUpICode, &followUp.FollowUpCreationDate, &followUp.FollowUpUpdateDate, &followUp.FollowUpPhysicalViolenceWitnessedByFamily, &followUp.FollowUpViolenceEscalation,
		&followUp.FollowUpAssaultWithWeapon, &followUp.FollowUpRecentControllingOrJealousBehavior, &followUp.FollowUpViolenceHistoryWithExPartner, &followUp.FollowUpViolenceHistoryWithOthers,
		&followUp.FollowUpSubstanceAbuse, &followUp.FollowUpViolenceJustification, &followUp.FollowUpVictimVulnerability, &followUp.FollowUpRiskLevel, &followUp.FollowUpOwnerGeneralUser, &followUp.FollowUpStatus,
		&alert.AlertId, &alert.AlertICode, &alert.AlertCreationDate, &alert.AlertType, &alert.AlertPriority, &alert.AlertCode)

	rel.RelAlertFollowUpFollowUp = followUp.ToDTO()
	rel.RelAlertFollowUpAlert = alert.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelAlertFollowUps obtiene una lista de relaciones entre alertas y seguimientos basándose en criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna una lista de DTOs con los datos de las relaciones y un error en caso de fallo.
func GetRelAlertFollowUps(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelAlertFollowUpDTO, error) {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var relPath string = RelAlertFollowUpDBScheme + "." + RelAlertFollowUpDBName
	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName
	var alertPath string = AlertDBScheme + "." + AlertDBName

	var rels []RelAlertFollowUpDTO

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelAlertFollowUpId", "RelAlertFollowUpData", "RelAlertFollowUpCreationDate", "RelAlertFollowUpFollowUp", "RelAlertFollowUpAlert"}
	var relFieldsAliasSlice []string = []string{}

	var followUpFieldsSlice []string = []string{"FollowUpICode", "FollowUpCreationDate", "FollowUpUpdateDate", "FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation", "FollowUpAssaultWithWeapon", "FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner", "FollowUpViolenceHistoryWithOthers", "FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification", "FollowUpVictimVulnerability", "FollowUpRiskLevel", "FollowUpOwnerGeneralUser", "FollowUpStatus"}
	var followUpFieldsAliasSlice []string = []string{}

	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}

	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelAlertFollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, true)
	var followUpFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, followUpFieldsSlice, followUpFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, true)
	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)

	var query string = `SELECT ` + relFieldsStr + `, ` + followUpFieldsStr + `, ` + alertFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + followUpPath + ` ON (` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + ` = ` + relPath + `.` + RelAlertFollowUpFieldDefinitions["RelAlertFollowUpFollowUp"].DBName + `)` +
		` LEFT JOIN ` + alertPath + ` ON (` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + ` = ` + relPath + `.` + RelAlertFollowUpFieldDefinitions["RelAlertFollowUpAlert"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelAlertFollowUpDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	for persistenceCtrl.Next() {
		var rel RelAlertFollowUpDTO
		var followUp FollowUpPgDB = FollowUpPgDB{}
		var alert AlertPgDB = AlertPgDB{}

		persistenceCtrl.ScanRow(&rel.RelAlertFollowUpId, &rel.RelAlertFollowUpData, &rel.RelAlertFollowUpCreationDate,
			&followUp.FollowUpId, &alert.AlertId,
			&followUp.FollowUpICode, &followUp.FollowUpCreationDate, &followUp.FollowUpUpdateDate, &followUp.FollowUpPhysicalViolenceWitnessedByFamily, &followUp.FollowUpViolenceEscalation,
			&followUp.FollowUpAssaultWithWeapon, &followUp.FollowUpRecentControllingOrJealousBehavior, &followUp.FollowUpViolenceHistoryWithExPartner, &followUp.FollowUpViolenceHistoryWithOthers,
			&followUp.FollowUpSubstanceAbuse, &followUp.FollowUpViolenceJustification, &followUp.FollowUpVictimVulnerability, &followUp.FollowUpRiskLevel, &followUp.FollowUpOwnerGeneralUser, &followUp.FollowUpStatus,
			&alert.AlertId, &alert.AlertICode, &alert.AlertCreationDate, &alert.AlertType, &alert.AlertPriority, &alert.AlertCode)

		rel.RelAlertFollowUpFollowUp = followUp.ToDTO()
		rel.RelAlertFollowUpAlert = alert.ToDTO()
		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// GetAllRelAlertFollowUp obtiene todas las relaciones entre alertas y seguimientos sin aplicar filtros.
// Parámetros:
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna una lista de DTOs con todas las relaciones y un error en caso de fallo.
func GetAllRelAlertFollowUp(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelAlertFollowUpDTO, error) {
	var rel RelAlertFollowUpDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelAlertFollowUpId", "RelAlertFollowUpData", "RelAlertFollowUpCreationDate", "RelAlertFollowUpFollowUp", "RelAlertFollowUpAlert"}
	var relFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, relFieldsSlice, relFieldsAliasSlice, RelAlertFollowUpDBName, []string{}, []string{}, []string{}, "", RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query)
	var rels []RelAlertFollowUpDTO
	for persistenceCtrl.Next() {
		var followUp FollowUpPgDB = FollowUpPgDB{}
		var alert AlertPgDB = AlertPgDB{}
		rel = RelAlertFollowUpDTO{}
		persistenceCtrl.ScanRow(&rel.RelAlertFollowUpId, &rel.RelAlertFollowUpData, &rel.RelAlertFollowUpCreationDate,
			&followUp.FollowUpId, &alert.AlertId)

		rel.RelAlertFollowUpFollowUp = followUp.ToDTO()
		rel.RelAlertFollowUpAlert = alert.ToDTO()
		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// GetAllRelAlertFollowUpWithFollowUpAndAlert obtiene todas las relaciones entre alertas y seguimientos
// junto con los detalles asociados del seguimiento y de la alerta.
// Parámetros:
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna una lista de DTOs con las relaciones y sus datos asociados, y un error en caso de fallo.
func GetAllRelAlertFollowUpWithFollowUpAndAlert(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelAlertFollowUpDTO, error) {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var relPath string = RelAlertFollowUpDBScheme + "." + RelAlertFollowUpDBName
	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName
	var alertPath string = AlertDBScheme + "." + AlertDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelAlertFollowUpId", "RelAlertFollowUpData", "RelAlertFollowUpCreationDate", "RelAlertFollowUpFollowUp", "RelAlertFollowUpAlert"}
	var relFieldsAliasSlice []string = []string{}

	var followUpFieldsSlice []string = []string{"FollowUpICode", "FollowUpCreationDate", "FollowUpUpdateDate", "FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation", "FollowUpAssaultWithWeapon", "FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner", "FollowUpViolenceHistoryWithOthers", "FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification", "FollowUpVictimVulnerability", "FollowUpRiskLevel", "FollowUpOwnerGeneralUser", "FollowUpStatus"}
	var followUpFieldsAliasSlice []string = []string{}

	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}

	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelAlertFollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelAlertFollowUpDBScheme, RelAlertFollowUpFieldDefinitions, true)
	var followUpFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, followUpFieldsSlice, followUpFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, true)
	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)

	var query string = `SELECT ` + relFieldsStr + ", " + followUpFieldsStr + ", " + alertFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + followUpPath + ` ON (` + followUpPath + `.` + FollowUpFieldDefinitions["FollowUpId"].DBName + ` = ` + relPath + `.` + RelAlertFollowUpFieldDefinitions["RelAlertFollowUpFollowUp"].DBName + `)` +
		` LEFT JOIN ` + alertPath + ` ON (` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + ` = ` + relPath + `.` + RelAlertFollowUpFieldDefinitions["RelAlertFollowUpAlert"].DBName + `)`

	persistenceCtrl.Query(context.Background(), query)
	var rels []RelAlertFollowUpDTO
	for persistenceCtrl.Next() {
		var followUp FollowUpPgDB = FollowUpPgDB{}
		var alert AlertPgDB = AlertPgDB{}
		var rel = RelAlertFollowUpDTO{}

		persistenceCtrl.ScanRow(&rel.RelAlertFollowUpId, &rel.RelAlertFollowUpData, &rel.RelAlertFollowUpCreationDate,
			&followUp.FollowUpId, &alert.AlertId,
			&followUp.FollowUpICode, &followUp.FollowUpCreationDate, &followUp.FollowUpUpdateDate, &followUp.FollowUpPhysicalViolenceWitnessedByFamily, &followUp.FollowUpViolenceEscalation,
			&followUp.FollowUpAssaultWithWeapon, &followUp.FollowUpRecentControllingOrJealousBehavior, &followUp.FollowUpViolenceHistoryWithExPartner, &followUp.FollowUpViolenceHistoryWithOthers,
			&followUp.FollowUpSubstanceAbuse, &followUp.FollowUpViolenceJustification, &followUp.FollowUpVictimVulnerability, &followUp.FollowUpRiskLevel, &followUp.FollowUpOwnerGeneralUser, &followUp.FollowUpStatus,
			&alert.AlertId, &alert.AlertICode, &alert.AlertCreationDate, &alert.AlertType, &alert.AlertPriority, &alert.AlertCode)

		rel.RelAlertFollowUpFollowUp = followUp.ToDTO()
		rel.RelAlertFollowUpAlert = alert.ToDTO()

		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// SetRelAlertFollowUserDefaults asigna valores por defecto a ciertos campos de la relación en función de la acción que se realizará.
// Por ejemplo, en una inserción se establece la fecha de creación.
// Parámetros:
//   - rel: puntero al DTO de la relación a modificar.
//   - action: acción que se va a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetRelAlertFollowUserDefaults(rel *RelAlertFollowUpDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		rel.RelAlertFollowUpCreationDate = time.Now()

	case common_dao.SQL_UPDATE:
		// Para una actualización, se podrían asignar otros valores por defecto si fuera necesario.
	}
}

// PgDBToDTO convierte una instancia de RelAlertFollowUpPgDB (modelo de datos de PostgreSQL)
// a su correspondiente objeto de transferencia de datos (DTO) RelAlertFollowUpDTO.
// Retorna el DTO con los datos convertidos.
func (obj *RelAlertFollowUpPgDB) PgDBToDTO() RelAlertFollowUpDTO {
	var dto RelAlertFollowUpDTO

	if obj.RelAlertFollowUpId.Valid {
		dto.RelAlertFollowUpId = uint64(obj.RelAlertFollowUpId.Int64)
	}

	if obj.RelAlertFollowUpData.Valid {
		dto.RelAlertFollowUpData = obj.RelAlertFollowUpData.String
	}

	if obj.RelAlertFollowUpCreationDate.Valid {
		dto.RelAlertFollowUpCreationDate = obj.RelAlertFollowUpCreationDate.Time
	}

	if obj.RelAlertFollowUpFollowUp.Valid {
		dto.RelAlertFollowUpFollowUp = FollowUpDTO{FollowUpId: uint64(obj.RelAlertFollowUpFollowUp.Int64)}
	}

	if obj.RelAlertFollowUpAlert.Valid {
		dto.RelAlertFollowUpAlert = AlertDTO{AlertId: uint64(obj.RelAlertFollowUpAlert.Int64)}
	}

	return dto
}
