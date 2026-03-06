// El paquete salvia_daos contiene los Objetos de Acceso a Datos (DAO) para las entidades
// de la lógica de negocio de "salvia". Este archivo en particular maneja la persistencia
// de la entidad 'FollowUp' (Seguimiento).
package salvia_daos

import (
	// Importaciones de paquetes del proyecto y de la librería estándar de Go.
	common_config "bitsflow/common/config"           // Proporciona acceso a configuraciones globales, como formatos de fecha y hora.
	common_controllers "bitsflow/common/controllers" // Contiene controladores de persistencia y estructuras para consultas dinámicas.
	common_dao "bitsflow/common/dao"                 // Contiene constantes (como SQL_INSERT) y helpers para operaciones DAO.
	"bitsflow/common/db"                             // Maneja la configuración y la conexión con la base de datos.
	"bitsflow/common/utils"                          // Contiene funciones de utilidad, como la generación de UUIDs.
	"context"                                        // Permite manejar contextos, útil para controlar cancelaciones y timeouts en las consultas.
	"database/sql"                                   // Proporciona la interfaz genérica para trabajar con bases de datos SQL.
	"encoding/json"
	"errors" // Para la creación de errores personalizados.
	"fmt"    // Paquete para formateo de I/O, usado para imprimir errores de SQL.
	"time"   // Proporciona funcionalidad para medir y mostrar el tiempo.
)

// Este bloque var define metadatos constantes para la entidad FollowUp.
// Estos valores centralizan la configuración de la entidad, facilitando el mantenimiento y
// asegurando la consistencia en la construcción de consultas SQL y en la serialización de datos.
var (
	// FollowUpEntityName es el nombre lógico de la entidad dentro de la aplicación.
	FollowUpEntityName string = "FollowUp"
	// FollowUpJSONName es el nombre que se utilizará para la entidad en formatos JSON.
	FollowUpJSONName string = "followUp"
	// FollowUpDBName es el nombre exacto de la tabla en la base de datos.
	FollowUpDBName string = "follow_up"
	// FollowUpDBScheme es el esquema de la base de datos al que pertenece la tabla.
	FollowUpDBScheme string = "salvia"

	// FollowUpFieldDefinitions es un mapa que detalla las propiedades de cada campo del modelo FollowUp.
	// Asocia los nombres de los campos del struct de Go con los nombres de las columnas de la base de datos,
	// y define metadatos como el tipo de dato, tamaño y si es requerido. Es crucial para la construcción
	// dinámica de consultas SQL y para validaciones.
	FollowUpFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FollowUpId":           {Name: "FollowUpId", DBName: "follow_up_id", Alias: "", ModelType: "uint", Required: false},
		"FollowUpICode":        {Name: "FollowUpICode", DBName: "follow_up_i_code", Alias: "", ModelType: "string", MinSize: 36, MaxSize: 36, Required: false},
		"FollowUpCreationDate": {Name: "FollowUpCreationDate", DBName: "follow_up_creation_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpUpdateDate":   {Name: "FollowUpUpdateDate", DBName: "follow_up_update_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpPhysicalViolenceWitnessedByFamily":  {Name: "FollowUpPhysicalViolenceWitnessedByFamily", DBName: "follow_up_physical_violence_witnessed_by_family", Alias: "", ModelType: "string", MaxSize: 1, Required: true},
		"FollowUpViolenceEscalation":                 {Name: "FollowUpViolenceEscalation", DBName: "follow_up_violence_escalation", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpAssaultWithWeapon":                  {Name: "FollowUpAssaultWithWeapon", DBName: "follow_up_assault_with_weapon", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpRecentControllingOrJealousBehavior": {Name: "FollowUpRecentControllingOrJealousBehavior", DBName: "follow_up_recent_controlling_or_jealous_behavior", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpViolenceHistoryWithExPartner":       {Name: "FollowUpViolenceHistoryWithExPartner", DBName: "follow_up_violence_history_with_ex_partner", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpViolenceHistoryWithOthers":          {Name: "FollowUpViolenceHistoryWithOthers", DBName: "follow_up_violence_history_with_others", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpSubstanceAbuse":                     {Name: "FollowUpSubstanceAbuse", DBName: "follow_up_substance_abuse", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpViolenceJustification":              {Name: "FollowUpViolenceJustification", DBName: "follow_up_violence_justification", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpVictimVulnerability":                {Name: "FollowUpVictimVulnerability", DBName: "follow_up_victim_vulnerability", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpRiskLevel":                          {Name: "FollowUpRiskLevel", DBName: "follow_up_risk_level", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"FollowUpOwnerGeneralUser":                   {Name: "FollowUpOwnerGeneralUser", DBName: "follow_up_owner_general_user", Alias: "", ModelType: "string", MinSize: 36, MaxSize: 36, Required: true},
		"FollowUpStatus":                             {Name: "FollowUpStatus", DBName: "follow_up_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
	}
)

// FollowUpDTO (Data Transfer Object) es la estructura que representa un seguimiento en la capa de aplicación.
// Se utiliza para intercambiar datos con el exterior (ej. API) y para la lógica de negocio.
// Contiene campos que parecen corresponder a un cuestionario de evaluación de riesgos.
type FollowUpDTO struct {
	FollowUpId                                 uint64    `json:"-"`
	FollowUpICode                              string    `json:"icode"`
	FollowUpCreationDate                       time.Time `json:"creationDate"`
	FollowUpUpdateDate                         time.Time `json:"updateDate"`
	FollowUpPhysicalViolenceWitnessedByFamily  string    `json:"physicalViolenceWitnessedByFamily"`
	FollowUpViolenceEscalation                 string    `json:"violenceEscalation"`
	FollowUpAssaultWithWeapon                  string    `json:"assaultWithWeapon"`
	FollowUpRecentControllingOrJealousBehavior string    `json:"recentControllingOrJealousBehavior"`
	FollowUpViolenceHistoryWithExPartner       string    `json:"violenceHistoryWithExPartner"`
	FollowUpViolenceHistoryWithOthers          string    `json:"violenceHistoryWithOthers"`
	FollowUpSubstanceAbuse                     string    `json:"substanceAbuse"`
	FollowUpViolenceJustification              string    `json:"violenceJustification"`
	FollowUpVictimVulnerability                string    `json:"victimVulnerability"`
	FollowUpRiskLevel                          string    `json:"riskLevel"`
	FollowUpOwnerGeneralUser                   string    `json:"ownerGeneralUser"`
	FollowUpStatus                             string    `json:"status"`
	// FollowUpEntries representa la relación uno-a-muchos con las entradas de seguimiento.
	// Este campo se poblará con las entradas asociadas a este seguimiento principal.
	FollowUpEntries []FollowUpEntryDTO `json:"entries"`
}

// FollowUpPgDB representa la estructura de la tabla 'follow_up' de la base de datos.
// Utiliza tipos `sql.Null*` para manejar de forma segura las columnas que pueden contener
// valores NULL, evitando errores al escanear los resultados de las consultas.
type FollowUpPgDB struct {
	FollowUpId                                 sql.NullInt64
	FollowUpICode                              sql.NullString
	FollowUpCreationDate                       sql.NullTime
	FollowUpUpdateDate                         sql.NullTime
	FollowUpPhysicalViolenceWitnessedByFamily  sql.NullString
	FollowUpViolenceEscalation                 sql.NullString
	FollowUpAssaultWithWeapon                  sql.NullString
	FollowUpRecentControllingOrJealousBehavior sql.NullString
	FollowUpViolenceHistoryWithExPartner       sql.NullString
	FollowUpViolenceHistoryWithOthers          sql.NullString
	FollowUpSubstanceAbuse                     sql.NullString
	FollowUpViolenceJustification              sql.NullString
	FollowUpVictimVulnerability                sql.NullString
	FollowUpRiskLevel                          sql.NullString
	FollowUpOwnerGeneralUser                   sql.NullString
	FollowUpStatus                             sql.NullString
}

func (fu FollowUpDTO) MarshalJSON() ([]byte, error) {
	type Alias FollowUpDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		FollowUpCreationDate string `json:"creationDate"`
		FollowUpUpdateDate   string `json:"updateDate"`
	}{
		Alias:                (*Alias)(&fu),
		FollowUpCreationDate: fu.FollowUpCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FollowUpUpdateDate:   fu.FollowUpUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (fu *FollowUpDTO) UnmarshalJSON(data []byte) error {
	type Alias FollowUpDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		FollowUpCreationDate string `json:"creationDate"`
		FollowUpUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(fu),
	}

	// Unmarshal estructural (NO falla por fechas)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Helper de parseo tolerante
	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{} // fecha inválida → zero value
		}
		return t
	}

	// Parseo seguro
	fu.FollowUpCreationDate = parse(aux.FollowUpCreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	fu.FollowUpUpdateDate = parse(aux.FollowUpUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetFollowUp inserta un nuevo seguimiento en la base de datos.
// Recibe el objeto seguimiento, información de transacción, módulo, y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func SetFollowUp(followUp *FollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se define la lista de campos que se insertarán en la BD.
	var fuFieldsSlice []string = []string{
		"FollowUpICode",
		"FollowUpCreationDate",
		"FollowUpUpdateDate",
		"FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation",
		"FollowUpAssaultWithWeapon",
		"FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner",
		"FollowUpViolenceHistoryWithOthers",
		"FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification",
		"FollowUpVictimVulnerability",
		"FollowUpRiskLevel",
		"FollowUpOwnerGeneralUser",
		"FollowUpStatus",
	}
	var fuFieldsAliasSlice []string = []string{}

	// Se genera la query de inserción utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, fuFieldsSlice, fuFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{"FollowUpId"}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query,
		followUp.FollowUpICode,
		followUp.FollowUpCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUp.FollowUpUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUp.FollowUpPhysicalViolenceWitnessedByFamily,
		followUp.FollowUpViolenceEscalation,
		followUp.FollowUpAssaultWithWeapon,
		followUp.FollowUpRecentControllingOrJealousBehavior,
		followUp.FollowUpViolenceHistoryWithExPartner,
		followUp.FollowUpViolenceHistoryWithOthers,
		followUp.FollowUpSubstanceAbuse,
		followUp.FollowUpViolenceJustification,
		followUp.FollowUpVictimVulnerability,
		followUp.FollowUpRiskLevel,
		followUp.FollowUpOwnerGeneralUser,
		followUp.FollowUpStatus,
	)
	persistenceCtrl.Scan(&followUp.FollowUpId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateFollowUpByICode actualiza los datos de un seguimiento en la BD, utilizando su ICode como referencia.
// Recibe el objeto seguimiento con los nuevos datos, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func UpdateFollowUpByICode(followUp *FollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se actualizarán.
	var fuFieldsSlice []string = []string{
		"FollowUpUpdateDate",
		"FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation",
		"FollowUpAssaultWithWeapon",
		"FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner",
		"FollowUpViolenceHistoryWithOthers",
		"FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification",
		"FollowUpVictimVulnerability",
		"FollowUpRiskLevel",
		"FollowUpOwnerGeneralUser",
		"FollowUpStatus",
	}
	var fuFieldsAliasSlice []string = []string{}

	// Se genera la query de actualización utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, fuFieldsSlice, fuFieldsAliasSlice, FollowUpDBName, []string{"FollowUpICode"}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		followUp.FollowUpICode,
		followUp.FollowUpUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUp.FollowUpPhysicalViolenceWitnessedByFamily,
		followUp.FollowUpViolenceEscalation,
		followUp.FollowUpAssaultWithWeapon,
		followUp.FollowUpRecentControllingOrJealousBehavior,
		followUp.FollowUpViolenceHistoryWithExPartner,
		followUp.FollowUpViolenceHistoryWithOthers,
		followUp.FollowUpSubstanceAbuse,
		followUp.FollowUpViolenceJustification,
		followUp.FollowUpVictimVulnerability,
		followUp.FollowUpRiskLevel,
		followUp.FollowUpOwnerGeneralUser,
		followUp.FollowUpStatus,
	)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se valida que se hayan afectado filas en la actualización.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveFollowUpByICode elimina un seguimiento de la base de datos utilizando su ICode.
// Recibe el objeto seguimiento, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func RemoveFollowUpByICode(followUp *FollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos y condiciones para la eliminación.
	var fuFieldsSlice []string = []string{}
	var fuFieldsAliasSlice []string = []string{}

	// Se genera la query de eliminación utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, fuFieldsSlice, fuFieldsAliasSlice, FollowUpDBName, []string{}, []string{"FollowUpICode"}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, false)

	// Se ejecuta la query con el parámetro correspondiente.
	persistenceCtrl.Exec(context.Background(), query, followUp.FollowUpICode)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se valida que se hayan afectado filas en la eliminación.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetFollowUp consulta un seguimiento basado en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto seguimiento a completar,
// información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func GetFollowUp(by common_controllers.By, followUp *FollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla FollowUp en la BD.
	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var fuFieldsSlice []string = []string{
		"FollowUpId",
		"FollowUpICode",
		"FollowUpCreationDate",
		"FollowUpUpdateDate",
		"FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation",
		"FollowUpAssaultWithWeapon",
		"FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner",
		"FollowUpViolenceHistoryWithOthers",
		"FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification",
		"FollowUpVictimVulnerability",
		"FollowUpRiskLevel",
		"FollowUpOwnerGeneralUser",
		"FollowUpStatus",
	}
	var fuFieldsAliasSlice []string = []string{}
	var fuFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fuFieldsSlice, fuFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + fuFieldsStr +
		` FROM ` + followUpPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FollowUpDBName, by.AttrsName, []string{}, []string{}, by.Operator, FollowUpDBScheme, FollowUpFieldDefinitions, true)

	// Se ejecuta la query con los parámetros de filtrado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var followUpPg = FollowUpPgDB{}

	// Se escanean los resultados de la query.
	persistenceCtrl.Scan(
		&followUpPg.FollowUpId,
		&followUpPg.FollowUpICode,
		&followUpPg.FollowUpCreationDate,
		&followUpPg.FollowUpUpdateDate,
		&followUpPg.FollowUpPhysicalViolenceWitnessedByFamily,
		&followUpPg.FollowUpViolenceEscalation,
		&followUpPg.FollowUpAssaultWithWeapon,
		&followUpPg.FollowUpRecentControllingOrJealousBehavior,
		&followUpPg.FollowUpViolenceHistoryWithExPartner,
		&followUpPg.FollowUpViolenceHistoryWithOthers,
		&followUpPg.FollowUpSubstanceAbuse,
		&followUpPg.FollowUpViolenceJustification,
		&followUpPg.FollowUpVictimVulnerability,
		&followUpPg.FollowUpRiskLevel,
		&followUpPg.FollowUpOwnerGeneralUser,
		&followUpPg.FollowUpStatus,
	)

	// Se convierte el objeto de base de datos a DTO.
	*followUp = followUpPg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetFollowUp consulta seguimientos basado en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto seguimiento a completar,
// información de transacción, módulo y datos de conexión.
// Retorna los seguimientos o un error en caso de producirse.
func GetFollowUps(by common_controllers.By, followUp *FollowUpDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FollowUpDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla FollowUp en la BD.
	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var fuFieldsSlice []string = []string{
		"FollowUpId",
		"FollowUpICode",
		"FollowUpCreationDate",
		"FollowUpUpdateDate",
		"FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation",
		"FollowUpAssaultWithWeapon",
		"FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner",
		"FollowUpViolenceHistoryWithOthers",
		"FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification",
		"FollowUpVictimVulnerability",
		"FollowUpRiskLevel",
		"FollowUpOwnerGeneralUser",
		"FollowUpStatus",
	}
	var fuFieldsAliasSlice []string = []string{}
	var fuFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fuFieldsSlice, fuFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, true)

	// Se genera la query para obtener todos los seguimientos.
	var query string = `SELECT ` + fuFieldsStr +
		` FROM ` + followUpPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FollowUpDBName, by.AttrsName, []string{}, []string{}, by.Operator, FollowUpDBScheme, FollowUpFieldDefinitions, true)

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	var followUps []FollowUpDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var followUpPg = FollowUpPgDB{}
		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(
			&followUpPg.FollowUpId,
			&followUpPg.FollowUpICode,
			&followUpPg.FollowUpCreationDate,
			&followUpPg.FollowUpUpdateDate,
			&followUpPg.FollowUpPhysicalViolenceWitnessedByFamily,
			&followUpPg.FollowUpViolenceEscalation,
			&followUpPg.FollowUpAssaultWithWeapon,
			&followUpPg.FollowUpRecentControllingOrJealousBehavior,
			&followUpPg.FollowUpViolenceHistoryWithExPartner,
			&followUpPg.FollowUpViolenceHistoryWithOthers,
			&followUpPg.FollowUpSubstanceAbuse,
			&followUpPg.FollowUpViolenceJustification,
			&followUpPg.FollowUpVictimVulnerability,
			&followUpPg.FollowUpRiskLevel,
			&followUpPg.FollowUpOwnerGeneralUser,
			&followUpPg.FollowUpStatus,
		)

		var followUp FollowUpDTO = followUpPg.ToDTO()
		followUps = append(followUps, followUp)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return followUps, nil
}

// GetAllFollowUps consulta todos los seguimientos existentes en la base de datos.
// Recibe información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada, un slice de FollowUpDTO y un error en caso de producirse.
func GetAllFollowUps(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FollowUpDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var followUpPath string = FollowUpDBScheme + "." + FollowUpDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar para los seguimientos.
	var fuFieldsSlice []string = []string{
		"FollowUpId",
		"FollowUpICode",
		"FollowUpCreationDate",
		"FollowUpUpdateDate",
		"FollowUpPhysicalViolenceWitnessedByFamily",
		"FollowUpViolenceEscalation",
		"FollowUpAssaultWithWeapon",
		"FollowUpRecentControllingOrJealousBehavior",
		"FollowUpViolenceHistoryWithExPartner",
		"FollowUpViolenceHistoryWithOthers",
		"FollowUpSubstanceAbuse",
		"FollowUpViolenceJustification",
		"FollowUpVictimVulnerability",
		"FollowUpRiskLevel",
		"FollowUpOwnerGeneralUser",
		"FollowUpStatus",
	}
	var fuFieldsAliasSlice []string = []string{}
	var fuFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, fuFieldsSlice, fuFieldsAliasSlice, FollowUpDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpDBScheme, FollowUpFieldDefinitions, true)

	// Se genera la query para obtener todos los seguimientos.
	var query string = `SELECT ` + fuFieldsStr +
		` FROM ` + followUpPath +
		` WHERE TRUE`

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query)
	var followUps []FollowUpDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var followUpPg = FollowUpPgDB{}
		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(
			&followUpPg.FollowUpId,
			&followUpPg.FollowUpICode,
			&followUpPg.FollowUpCreationDate,
			&followUpPg.FollowUpUpdateDate,
			&followUpPg.FollowUpPhysicalViolenceWitnessedByFamily,
			&followUpPg.FollowUpViolenceEscalation,
			&followUpPg.FollowUpAssaultWithWeapon,
			&followUpPg.FollowUpRecentControllingOrJealousBehavior,
			&followUpPg.FollowUpViolenceHistoryWithExPartner,
			&followUpPg.FollowUpViolenceHistoryWithOthers,
			&followUpPg.FollowUpSubstanceAbuse,
			&followUpPg.FollowUpViolenceJustification,
			&followUpPg.FollowUpVictimVulnerability,
			&followUpPg.FollowUpRiskLevel,
			&followUpPg.FollowUpOwnerGeneralUser,
			&followUpPg.FollowUpStatus,
		)

		var followUp FollowUpDTO = followUpPg.ToDTO()
		followUps = append(followUps, followUp)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return followUps, nil
}

// SetFollowUpDefaults asigna valores por defecto a los campos de un seguimiento,
// dependiendo de la acción que se esté realizando (insertar o actualizar).
func SetFollowUpDefaults(followUp *FollowUpDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción se asigna estado 'p', fechas actuales y se genera un UUID para el ICode.
		followUp.FollowUpStatus = "p"
		followUp.FollowUpCreationDate = time.Now()
		followUp.FollowUpUpdateDate = time.Now()
		followUp.FollowUpICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para actualización solo se actualiza la fecha de modificación.
		followUp.FollowUpUpdateDate = time.Now()
	}
}

// ToDTO convierte un objeto FollowUpPgDB obtenido de la base de datos en un objeto FollowUpDTO.
// Se encarga de verificar la validez de los campos nulos y asignarlos correctamente.
func (obj *FollowUpPgDB) ToDTO() FollowUpDTO {
	var dto FollowUpDTO

	if obj.FollowUpId.Valid {
		dto.FollowUpId = uint64(obj.FollowUpId.Int64)
	}

	if obj.FollowUpICode.Valid {
		dto.FollowUpICode = obj.FollowUpICode.String
	}

	if obj.FollowUpCreationDate.Valid {
		dto.FollowUpCreationDate = obj.FollowUpCreationDate.Time
	}

	if obj.FollowUpUpdateDate.Valid {
		dto.FollowUpUpdateDate = obj.FollowUpUpdateDate.Time
	}

	if obj.FollowUpPhysicalViolenceWitnessedByFamily.Valid {
		dto.FollowUpPhysicalViolenceWitnessedByFamily = obj.FollowUpPhysicalViolenceWitnessedByFamily.String
	}

	if obj.FollowUpViolenceEscalation.Valid {
		dto.FollowUpViolenceEscalation = obj.FollowUpViolenceEscalation.String
	}

	if obj.FollowUpAssaultWithWeapon.Valid {
		dto.FollowUpAssaultWithWeapon = obj.FollowUpAssaultWithWeapon.String
	}

	if obj.FollowUpRecentControllingOrJealousBehavior.Valid {
		dto.FollowUpRecentControllingOrJealousBehavior = obj.FollowUpRecentControllingOrJealousBehavior.String
	}

	if obj.FollowUpViolenceHistoryWithExPartner.Valid {
		dto.FollowUpViolenceHistoryWithExPartner = obj.FollowUpViolenceHistoryWithExPartner.String
	}

	if obj.FollowUpViolenceHistoryWithOthers.Valid {
		dto.FollowUpViolenceHistoryWithOthers = obj.FollowUpViolenceHistoryWithOthers.String
	}

	if obj.FollowUpSubstanceAbuse.Valid {
		dto.FollowUpSubstanceAbuse = obj.FollowUpSubstanceAbuse.String
	}

	if obj.FollowUpViolenceJustification.Valid {
		dto.FollowUpViolenceJustification = obj.FollowUpViolenceJustification.String
	}

	if obj.FollowUpVictimVulnerability.Valid {
		dto.FollowUpVictimVulnerability = obj.FollowUpVictimVulnerability.String
	}

	if obj.FollowUpRiskLevel.Valid {
		dto.FollowUpRiskLevel = obj.FollowUpRiskLevel.String
	}

	if obj.FollowUpOwnerGeneralUser.Valid {
		dto.FollowUpOwnerGeneralUser = obj.FollowUpOwnerGeneralUser.String
	}

	if obj.FollowUpStatus.Valid {
		dto.FollowUpStatus = obj.FollowUpStatus.String
	}

	return dto
}
