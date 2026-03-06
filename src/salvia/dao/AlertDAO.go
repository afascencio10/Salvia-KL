// Package salvia_daos contiene los Data Access Objects (DAO) para el manejo de alertas
// en el módulo Salvia.
package salvia_daos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
)

var (
	// Nombres y esquemas asociados a la entidad Alert.
	AlertEntityName string = "Alert"
	AlertJSONName   string = "alert"
	AlertDBName     string = "alert"
	AlertDBScheme   string = "salvia"

	// AlertFieldDefinitions define las propiedades de cada campo de Alert en el JSON y en la base de datos.
	// La estructura utils.FieldDefinition contiene:
	// Name: nombre del campo en el modelo,
	// DBName: nombre del campo en la base de datos,
	// Alias: alias opcional,
	// ModelType: tipo de dato en el modelo,
	// MinSize, MaxSize: tamaño mínimo y máximo,
	// Required: indica si el campo es obligatorio.
	AlertFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"AlertId":           {Name: "AlertId", DBName: "alert_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"AlertICode":        {Name: "AlertICode", DBName: "alert_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"AlertCreationDate": {Name: "AlertCreationDate", DBName: "alert_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"AlertType":         {Name: "AlertType", DBName: "alert_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 16, Required: true},
		"AlertPriority":     {Name: "AlertPriority", DBName: "alert_priority", Alias: "", ModelType: "int", MinSize: 0, MaxSize: 254, Required: true},
		"AlertCode":         {Name: "AlertCode", DBName: "alert_code", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 64, Required: true},
	}
)

// AlertDTO representa el Data Transfer Object de una alerta, que se utiliza para el intercambio
// de información entre capas sin exponer directamente los campos de la base de datos.
type AlertDTO struct {
	AlertId           uint64    `json:"-"`        // Identificador único de la alerta (no se expone en JSON)
	AlertICode        string    `json:"icode"`    // Código identificador de la alerta
	AlertCreationDate time.Time `json:"-"`        // Fecha de creación de la alerta (no se expone en JSON)
	AlertType         string    `json:"type"`     // Tipo de alerta
	AlertPriority     int       `json:"priority"` // Prioridad de la alerta
	AlertCode         string    `json:"code"`     // Código de la alerta
	// Atributos sin persistencia en la base de datos
	AlertData             string    `json:"data"`            // Datos adicionales de la alerta
	AlertDataCreationDate time.Time `json:"creationDate"`    // Fecha de creación de los datos adicionales
	AlertVictimCaseICode  string    `json:"victimCaseICode"` // Código identificador del caso de víctima asociado
}

// AlertPgDB representa la estructura utilizada para mapear los resultados de las consultas
// a la base de datos, utilizando tipos sql.Null* para manejar valores nulos.
type AlertPgDB struct {
	AlertId               sql.NullInt64  // Identificador de la alerta
	AlertICode            sql.NullString // Código identificador de la alerta
	AlertCreationDate     sql.NullTime   // Fecha de creación de la alerta
	AlertType             sql.NullString // Tipo de alerta
	AlertPriority         sql.NullInt64  // Prioridad de la alerta
	AlertCode             sql.NullString // Código de la alerta
	AlertData             sql.NullString // Datos adicionales de la alerta
	AlertDataCreationDate sql.NullTime   // Fecha de creación de los datos adicionales
	AlertVictimCaseICode  sql.NullString // Código del caso de víctima asociado
}

// GetAlert obtiene una única alerta de la base de datos según los atributos especificados en 'by'.
// Recibe como parámetros:
//   - by: criterios de búsqueda (atributos y operador)
//   - alert: puntero donde se almacenará el DTO resultante
//   - inTransaction: bandera para indicar si se está en una transacción
//   - module: módulo que invoca la consulta (para trazabilidad o logs)
//   - connData, clientConfig, serverConfig: configuración y datos de conexión a la base de datos
//
// Devuelve la conexión actualizada y un error en caso de producirse.
func GetAlert(by common_controllers.By, alert *AlertDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción del path de la tabla: esquema.tabla
	var alertPath string = AlertDBScheme + "." + AlertDBName

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Preparación de la consulta SQL
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}
	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr +
		` FROM ` + alertPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, AlertDBName, by.AttrsName, []string{}, []string{}, by.Operator, AlertDBScheme, AlertFieldDefinitions, true)

	// Ejecución de la consulta para obtener una sola fila
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)
	var alertPG AlertPgDB = AlertPgDB{}
	persistenceCtrl.Scan(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertCreationDate, &alertPG.AlertType, &alertPG.AlertPriority, &alertPG.AlertCode)
	*alert = alertPG.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetAlerts obtiene múltiples alertas de la base de datos aplicando filtros definidos en 'by'.
// Además, realiza JOIN con la tabla relacionada (por ejemplo, RelAlertVictimCase) para traer datos adicionales.
// Parámetros:
//   - by: criterios de búsqueda
//   - inTransaction: bandera para indicar si se ejecuta dentro de una transacción
//   - module: módulo invocante
//   - connData, clientConfig, serverConfig: parámetros de conexión y configuración de la BD
//
// Retorna la conexión actualizada, un slice de AlertDTO y un error en caso de fallo.
func GetAlerts(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AlertDTO, error) {
	// Inicialización del controlador de persistencia y definición de paths de tablas
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var alertPath string = AlertDBScheme + "." + AlertDBName
	var relAlertPath string = RelAlertVictimCaseDBScheme + "." + RelAlertVictimCaseDBName
	var alerts []AlertDTO = []AlertDTO{}

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	// Preparación de la consulta SQL para la alerta y la relación con victim case
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}
	var relAlertFieldsSlice []string = []string{"RelAlertVictimCase_CreationDate", "RelAlertVictimCase_Data"}
	var relAlertFieldsAliasSlice []string = []string{}

	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)
	var relAlertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relAlertFieldsSlice, relAlertFieldsAliasSlice, RelAlertVictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, RelAlertVictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr + `, ` + relAlertFieldsStr +
		` FROM ` + alertPath +
		` RIGHT JOIN ` + relAlertPath + ` ON (` + relAlertPath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_Alert"].DBName + ` = ` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, AlertDBName, by.AttrsName, []string{}, []string{}, by.Operator, AlertDBScheme, AlertFieldDefinitions, true)

	// Ejecución de la consulta
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	for persistenceCtrl.Next() {
		var alertPG AlertPgDB = AlertPgDB{}
		persistenceCtrl.ScanRow(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertCreationDate, &alertPG.AlertType, &alertPG.AlertPriority, &alertPG.AlertCode,
			&alertPG.AlertDataCreationDate, &alertPG.AlertData)
		alerts = append(alerts, alertPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	return alerts, nil
}

// GetAlertsByVictimCaseICode obtiene las alertas asociadas a un caso de víctima específico,
// identificado por el código vCaseIcode.
// Parámetros:
//   - vCaseIcode: código identificador del caso de víctima
//   - inTransaction, connData, clientConfig, serverConfig: parámetros de conexión y configuración
//
// Retorna la conexión actualizada, un slice de AlertDTO y un error en caso de fallo.
func GetAlertsByVictimCaseICode(vCaseIcode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AlertDTO, error) {
	// Inicialización del controlador de persistencia y definición de paths de tablas
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var alertPath string = AlertDBScheme + "." + AlertDBName
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var relAlertVictimCasePath string = RelAlertVictimCaseDBScheme + "." + RelAlertVictimCaseDBName
	var alerts []AlertDTO = []AlertDTO{}

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	// Preparación de la consulta SQL que une las tablas de alertas, relaciones y victim case
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode", "AlertData"}
	var alertFieldsAliasSlice []string = []string{}

	var relAlertFieldsSlice []string = []string{"RelAlertVictimCase_CreationDate", "RelAlertVictimCase_Data"}
	var relAlertFieldsAliasSlice []string = []string{}

	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)
	var relAlertVictimCaseStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relAlertFieldsSlice, relAlertFieldsAliasSlice, RelAlertVictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, RelAlertVictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr + `, ` + relAlertVictimCaseStr +
		` FROM ` + alertPath +
		` RIGHT JOIN ` + relAlertVictimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_Alert"].DBName + ` = ` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + `)` +
		` LEFT JOIN ` + victimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseIcode"].DBName + ` = $1 `
	// Ejecución de la consulta
	persistenceCtrl.QueryRow(context.Background(), query, vCaseIcode)
	for persistenceCtrl.Next() {
		var alertPG AlertPgDB = AlertPgDB{}
		persistenceCtrl.ScanRow(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertCreationDate, &alertPG.AlertType, &alertPG.AlertPriority, &alertPG.AlertCode,
			&alertPG.AlertDataCreationDate, &alertPG.AlertData)
		alerts = append(alerts, alertPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	return alerts, nil
}

// GetAlertsByUserICode obtiene las alertas asociadas a un usuario específico,
// identificado por su código userIcode. Se realizan múltiples JOIN para relacionar
// alertas, casos de víctimas y propietarios de casos.
// Parámetros:
//   - userIcode: código identificador del usuario
//   - inTransaction, connData, clientConfig, serverConfig: parámetros de conexión y configuración
//
// Retorna la conexión actualizada, un slice de AlertDTO y un error en caso de fallo.
func GetAlertsByUserICode(userIcode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AlertDTO, error) {
	// Inicialización del controlador de persistencia y definición de paths de tablas
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var alertPath string = AlertDBScheme + "." + AlertDBName
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var caseOwnerPath string = CaseOwnerDBScheme + "." + CaseOwnerDBName
	var relAlertVictimCasePath string = RelAlertVictimCaseDBScheme + "." + RelAlertVictimCaseDBName
	var relCaseOwnerVictimCasePath string = RelCaseOwnerVictimCaseDBScheme + "." + RelCaseOwnerVictimCaseDBName
	var alerts []AlertDTO = []AlertDTO{}

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	// Preparación de la consulta SQL que une múltiples tablas para relacionar al usuario con sus alertas
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}

	var relAlertFieldsSlice []string = []string{"RelAlertVictimCase_CreationDate", "RelAlertVictimCase_Data"}
	var relAlertFieldsAliasSlice []string = []string{}

	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)
	var relAlertVictimCaseStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relAlertFieldsSlice, relAlertFieldsAliasSlice, RelAlertVictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, RelAlertVictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr + `, ` + relAlertVictimCaseStr +
		` FROM ` + alertPath +
		` RIGHT JOIN ` + relAlertVictimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_Alert"].DBName + ` = ` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + `)` +
		` LEFT JOIN ` + victimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		` LEFT JOIN ` + relCaseOwnerVictimCasePath + ` ON (` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `) ` +
		` LEFT JOIN ` + caseOwnerPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerId"].DBName + ` = ` + relCaseOwnerVictimCasePath + `.` + RelCaseOwnerVictimCaseFieldDefinitions["RelCaseOwnerVictimCase_CaseOwner"].DBName + `) ` +
		// Líneas comentadas para joins adicionales que pueden ser habilitados si se requieren:
		//` LEFT JOIN ` + userPath + ` ON (` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `)` +
		//` LEFT JOIN ` + relRolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserGeneralUser"].DBName + ` = ` + userPath + `.` + GeneralUserFieldDefinitions["GeneralUserId"].DBName + `)` +
		//` LEFT JOIN ` + rolePath + ` ON (` + relRolePath + `.` + RelRoleGeneralUserFieldDefinitions["RelRoleGeneralUserRole"].DBName + ` = ` + rolePath + `.` + RoleFieldDefinitions["RoleId"].DBName + `)` +
		` WHERE ` + caseOwnerPath + `.` + CaseOwnerFieldDefinitions["CaseOwnerGeneralUser"].DBName + ` = $1`

	// Ejecución de la consulta
	persistenceCtrl.Query(context.Background(), query, userIcode)
	for persistenceCtrl.Next() {
		var alertPG AlertPgDB = AlertPgDB{}
		persistenceCtrl.ScanRow(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertCreationDate, &alertPG.AlertType, &alertPG.AlertPriority, &alertPG.AlertCode,
			&alertPG.AlertDataCreationDate, &alertPG.AlertData)
		alerts = append(alerts, alertPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	return alerts, nil
}

// GetAlertsByTownCode obtiene las alertas correspondientes a un código de ciudad (townCode).
// Se realiza un JOIN con la tabla victimCase para relacionar la alerta con la localidad.
// Parámetros:
//   - townCode: código de la ciudad
//   - inTransaction, connData, clientConfig, serverConfig: parámetros de conexión y configuración
//
// Retorna la conexión actualizada, un slice de AlertDTO y un error en caso de fallo.
func GetAlertsByTownCode(townCode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AlertDTO, error) {
	// Inicialización del controlador de persistencia y definición de paths de tablas
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var alertPath string = AlertDBScheme + "." + AlertDBName
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var relAlertVictimCasePath string = RelAlertVictimCaseDBScheme + "." + RelAlertVictimCaseDBName
	var alerts []AlertDTO = []AlertDTO{}

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	// Preparación de la consulta SQL que une alertas y victim case
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode", "AlertData"}
	var alertFieldsAliasSlice []string = []string{}

	var relAlertFieldsSlice []string = []string{"RelAlertVictimCase_CreationDate", "RelAlertVictimCase_Data"}
	var relAlertFieldsAliasSlice []string = []string{}

	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)
	var relAlertVictimCaseStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relAlertFieldsSlice, relAlertFieldsAliasSlice, RelAlertVictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, RelAlertVictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr + `, ` + relAlertVictimCaseStr +
		` FROM ` + alertPath +
		` RIGHT JOIN ` + relAlertVictimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_Alert"].DBName + ` = ` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + `)` +
		` LEFT JOIN ` + victimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTowncode"].DBName + ` = $1 `
	// Ejecución de la consulta
	persistenceCtrl.QueryRow(context.Background(), query, townCode)
	for persistenceCtrl.Next() {
		var alertPG AlertPgDB = AlertPgDB{}
		persistenceCtrl.ScanRow(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertCreationDate, &alertPG.AlertType, &alertPG.AlertPriority, &alertPG.AlertCode,
			&alertPG.AlertDataCreationDate, &alertPG.AlertData)
		alerts = append(alerts, alertPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	return alerts, nil
}

// GetAlertsByTownCodeAndEntityBranchICode obtiene alertas filtradas por código de ciudad (townCode)
// y código de rama de entidad (entity branch i-code). Se realizan varios JOIN para relacionar alertas,
// victim case, momentos y ramas de entidad.
// Parámetros:
//   - townCode: código de la ciudad
//   - inTransaction, connData, clientConfig, serverConfig: parámetros de conexión y configuración
//
// Retorna la conexión actualizada, un slice de AlertDTO y un error en caso de fallo.
func GetAlertsByTownCodeAndEntityBranchICode(townCode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AlertDTO, error) {
	// Inicialización del controlador de persistencia y definición de paths de tablas
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var alertPath string = AlertDBScheme + "." + AlertDBName
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName
	var momentPath string = MomentDBScheme + "." + MomentDBName
	var entBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName
	var relAlertVictimCasePath string = RelAlertVictimCaseDBScheme + "." + RelAlertVictimCaseDBName
	var alerts []AlertDTO = []AlertDTO{}

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	// Preparación de la consulta SQL que une alertas, victim case, momentos y entidad rama
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertCreationDate", "AlertType", "AlertPriority", "AlertCode", "AlertData"}
	var alertFieldsAliasSlice []string = []string{}

	var relAlertFieldsSlice []string = []string{"RelAlertVictimCase_CreationDate", "RelAlertVictimCase_Data"}
	var relAlertFieldsAliasSlice []string = []string{}

	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)
	var relAlertVictimCaseStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relAlertFieldsSlice, relAlertFieldsAliasSlice, RelAlertVictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, RelAlertVictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr + `, ` + relAlertVictimCaseStr +
		` FROM ` + alertPath +
		` RIGHT JOIN ` + relAlertVictimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_Alert"].DBName + ` = ` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + `)` +
		` LEFT JOIN ` + victimCasePath + ` ON (` + relAlertVictimCasePath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		` LEFT JOIN ` + momentPath + ` ON (` + momentPath + `.` + MomentFieldDefinitions["MomentVictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)` +
		` LEFT JOIN ` + entBranchPath + ` ON (` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchId"].DBName + ` = ` + momentPath + `.` + MomentFieldDefinitions["MomentEntityBranch"].DBName + `)` +
		` WHERE ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseTowncode"].DBName + ` = $1 AND ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchICode"].DBName + ` = $2 `
	// Ejecución de la consulta
	persistenceCtrl.QueryRow(context.Background(), query, townCode)
	for persistenceCtrl.Next() {
		var alertPG AlertPgDB = AlertPgDB{}
		persistenceCtrl.ScanRow(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertCreationDate, &alertPG.AlertType, &alertPG.AlertPriority, &alertPG.AlertCode,
			&alertPG.AlertDataCreationDate, &alertPG.AlertData)
		alerts = append(alerts, alertPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return alerts, persistenceCtrl.Error
	}

	return alerts, nil
}

// GetAllAlerts obtiene todas las alertas registradas en la base de datos, realizando joins
// con las tablas de relaciones y victim case para traer información complementaria.
// Parámetros:
//   - inTransaction, connData, clientConfig, serverConfig: parámetros de conexión y configuración
//
// Retorna la conexión actualizada, un slice de AlertDTO con todas las alertas y un error en caso de fallo.
func GetAllAlerts(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]AlertDTO, error) {
	// Inicialización del controlador de persistencia y definición de paths de tablas
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var alertPath string = AlertDBScheme + "." + AlertDBName
	var relAlertPath string = RelAlertVictimCaseDBScheme + "." + RelAlertVictimCaseDBName
	var victimCasePath string = VictimCaseDBScheme + "." + VictimCaseDBName

	// Obtención de la conexión a la base de datos
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Preparación de la consulta SQL que une alertas, relaciones y victim case
	var alertFieldsSlice []string = []string{"AlertId", "AlertICode", "AlertType", "AlertPriority", "AlertCode"}
	var alertFieldsAliasSlice []string = []string{}
	var relAlertFieldsSlice []string = []string{"RelAlertVictimCase_CreationDate", "RelAlertVictimCase_Data"}
	var relAlertFieldsAliasSlice []string = []string{}

	var victimCaseFieldsSlice []string = []string{"VictimCaseICode"}
	var victimCaseFieldsAliasSlice []string = []string{}

	var victimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, victimCaseFieldsSlice, victimCaseFieldsAliasSlice, VictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, VictimCaseDBScheme, VictimCaseFieldDefinitions, true)
	var alertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, alertFieldsSlice, alertFieldsAliasSlice, AlertDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, AlertFieldDefinitions, true)
	var relAlertFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relAlertFieldsSlice, relAlertFieldsAliasSlice, RelAlertVictimCaseDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, AlertDBScheme, RelAlertVictimCaseFieldDefinitions, true)

	var query string = `SELECT ` + alertFieldsStr + `, ` + relAlertFieldsStr + `, ` + victimCaseFieldsStr +
		` FROM ` + alertPath +
		` RIGHT JOIN ` + relAlertPath + ` ON (` + relAlertPath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_Alert"].DBName + ` = ` + alertPath + `.` + AlertFieldDefinitions["AlertId"].DBName + `)` +
		` LEFT JOIN ` + victimCasePath + ` ON (` + relAlertPath + `.` + RelAlertVictimCaseFieldDefinitions["RelAlertVictimCase_VictimCase"].DBName + ` = ` + victimCasePath + `.` + VictimCaseFieldDefinitions["VictimCaseId"].DBName + `)`

	// Ejecución de la consulta
	persistenceCtrl.Query(context.Background(), query)
	var alerts []AlertDTO
	for persistenceCtrl.Next() {
		var alertPG AlertPgDB = AlertPgDB{}
		persistenceCtrl.ScanRow(&alertPG.AlertId, &alertPG.AlertICode, &alertPG.AlertType, &alertPG.AlertPriority,
			&alertPG.AlertCode, &alertPG.AlertCreationDate, &alertPG.AlertData, &alertPG.AlertVictimCaseICode)
		alerts = append(alerts, alertPG.ToDTO())
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return alerts, nil
}

// SetAlertDefaults establece valores predeterminados para una alerta según la acción a realizar.
// Actualmente, la función contempla las acciones de inserción (SQL_INSERT) y actualización (SQL_UPDATE),
// y puede ser extendida para definir valores por defecto específicos.
// Parámetros:
//   - alert: puntero al AlertDTO al cual se le aplicarán los valores por defecto
//   - action: acción que se va a realizar (por ejemplo, common_dao.SQL_INSERT o common_dao.SQL_UPDATE)
func SetAlertDefaults(alert *AlertDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Aquí se pueden definir valores por defecto antes de una inserción.
	case common_dao.SQL_UPDATE:
		// Aquí se pueden definir valores por defecto antes de una actualización.
	}
}

// ToDTO convierte un objeto AlertPgDB (con valores nulos manejados por sql.Null*) a su correspondiente AlertDTO.
func (obj *AlertPgDB) ToDTO() AlertDTO {
	var dto AlertDTO

	if obj.AlertId.Valid {
		dto.AlertId = uint64(obj.AlertId.Int64)
	}

	if obj.AlertICode.Valid {
		dto.AlertICode = obj.AlertICode.String
	}

	if obj.AlertCreationDate.Valid {
		dto.AlertCreationDate = obj.AlertCreationDate.Time
	}

	if obj.AlertType.Valid {
		dto.AlertType = obj.AlertType.String
	}

	if obj.AlertPriority.Valid {
		dto.AlertPriority = int(obj.AlertPriority.Int64)
	}

	if obj.AlertCode.Valid {
		dto.AlertCode = obj.AlertCode.String
	}

	if obj.AlertData.Valid {
		dto.AlertData = obj.AlertData.String
	}

	if obj.AlertVictimCaseICode.Valid {
		dto.AlertVictimCaseICode = obj.AlertVictimCaseICode.String
	}

	if obj.AlertDataCreationDate.Valid {
		dto.AlertDataCreationDate = obj.AlertDataCreationDate.Time
	}

	return dto
}
