// Package salvia_daos contiene los Data Access Objects (DAO) para la entidad RelCaseOwnerFeminicideRisk,
// encargados de interactuar con la base de datos.
package salvia_daos

import (
	// Importación de paquetes internos y externos necesarios para la funcionalidad del DAO.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	// Constantes y nombres relacionados con la entidad RelCaseOwnerFeminicideRisk.
	RelCaseOwnerFeminicideRisk         string = "RelCaseOwnerFeminicideRisk" // Nombre interno de la entidad.
	RelCaseOwnerFeminicideRiskJSONName string = "relCaseOwner"               // Nombre de la entidad en el JSON.
	RelCaseOwnerFeminicideRiskDBName   string = "rel_case_owner_feminicide"  // Nombre de la tabla en la base de datos.
	RelCaseOwnerFeminicideRiskDBScheme string = "salvia"                     // Esquema de la base de datos.

	// Atributos relacionados con las validaciones de campos.
	// Cada entrada del mapa define la relación entre el campo del JSON, su nombre en la base de datos,
	// el tipo de dato esperado en el modelo y otros parámetros de validación.
	RelCaseOwnerFeminicideRiskFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelCaseOwnerFeminicideRiskId":              {Name: "RelCaseOwnerFeminicideRiskId", DBName: "rel_case_owner_feminicide_risk_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelCaseOwnerFeminicideRisk_CaseOwner":      {Name: "RelCaseOwnerFeminicideRisk_CaseOwner", DBName: "case_owner_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerFeminicideRisk_FeminicideRisk": {Name: "RelCaseOwnerFeminicideRisk_FeminicideRisk", DBName: "feminicide_risk_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerFeminicideRisk_CreationDate":   {Name: "RelCaseOwnerFeminicideRisk_CreationDate", DBName: "rel_case_owner_feminicide_risk_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerFeminicideRisk_Status":         {Name: "RelCaseOwnerFeminicideRisk_Status", DBName: "rel_case_owner_feminicide_risk_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
	}
)

/*
RelCaseOwnerFeminicideRiskDTO representa la estructura de datos (DTO) para la entidad RelCaseOwnerFeminicideRisk.
Este DTO se utiliza para transferir datos entre las capas de la aplicación.
*/
type RelCaseOwnerFeminicideRiskDTO struct {
	RelCaseOwnerFeminicideRiskId              uint64    `json:"-"`
	RelCaseOwnerFeminicideRisk_CaseOwner      uint64    `json:"caseOwner"`
	RelCaseOwnerFeminicideRisk_FeminicideRisk uint64    `json:"feminicideRisk"`
	RelCaseOwnerFeminicideRisk_CreationDate   time.Time `json:"creationDate"`
	RelCaseOwnerFeminicideRisk_Status         string    `json:"status"`
}

/*
RelCaseOwnerFeminicideRiskPgDB representa la estructura que mapea la entidad en la base de datos.
Utiliza tipos sql.Null* para manejar valores nulos en la base de datos.
*/
type RelCaseOwnerFeminicideRiskPgDB struct {
	RelCaseOwnerFeminicideRiskId              sql.NullInt64
	RelCaseOwnerFeminicideRisk_CaseOwner      sql.NullInt64
	RelCaseOwnerFeminicideRisk_FeminicideRisk sql.NullInt64
	RelCaseOwnerFeminicideRisk_CreationDate   sql.NullTime
	RelCaseOwnerFeminicideRisk_Status         sql.NullString
}

func (rcwvc RelCaseOwnerFeminicideRiskDTO) MarshalJSON() ([]byte, error) {
	type Alias RelCaseOwnerFeminicideRiskDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		RelCaseOwnerFeminicideRisk_CreationDate string `json:"creationDate"`
	}{
		Alias:                                   (*Alias)(&rcwvc),
		RelCaseOwnerFeminicideRisk_CreationDate: rcwvc.RelCaseOwnerFeminicideRisk_CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (rcwvc *RelCaseOwnerFeminicideRiskDTO) UnmarshalJSON(data []byte) error {
	type Alias RelCaseOwnerFeminicideRiskDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		RelCaseOwnerFeminicideRisk_CreationDate string `json:"creationDate"`
	}{
		Alias: (*Alias)(rcwvc),
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
	rcwvc.RelCaseOwnerFeminicideRisk_CreationDate = parse(aux.RelCaseOwnerFeminicideRisk_CreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

/*
SetRelCaseOwnerFeminicideRisk inserta un nuevo registro de RelCaseOwnerFeminicideRisk en la base de datos.

Parámetros:
  - relCaseOwnerFeminicideRisk: Puntero al DTO que contiene los datos a insertar.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:
  - Un error en caso de ocurrir alguno.
*/
func SetRelCaseOwnerFeminicideRisk(relCaseOwnerFeminicideRisk *RelCaseOwnerFeminicideRiskDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a insertar en la base de datos.
	var relCaseOwnerFeminicideRiskFieldsSlice []string = []string{"RelCaseOwnerFeminicideRisk_CaseOwner", "RelCaseOwnerFeminicideRisk_FeminicideRisk", "RelCaseOwnerFeminicideRisk_CreationDate", "RelCaseOwnerFeminicideRisk_Status"}
	var relCaseOwnerFeminicideRiskFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, relCaseOwnerFeminicideRiskFieldsSlice, relCaseOwnerFeminicideRiskFieldsAliasSlice, RelCaseOwnerFeminicideRiskDBName, []string{}, []string{}, []string{"RelCaseOwnerFeminicideRiskId"}, common_dao.SQL_AND, RelCaseOwnerFeminicideRiskDBScheme, RelCaseOwnerFeminicideRiskFieldDefinitions, false)

	// Ejecución de la consulta SQL y asignación del ID generado al DTO.
	persistenceCtrl.QueryRow(context.Background(), query,
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_CaseOwner,
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_FeminicideRisk,
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_Status)
	persistenceCtrl.Scan(&relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRiskId)

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

/*
GetRelCaseOwnerFeminicideRisks recupera registros de RelCaseOwnerFeminicideRisk desde la base de datos
según los criterios de filtrado proporcionados.

Parámetros:
  - by: Estructura que contiene los nombres y valores de los atributos por los que se filtra.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:

  - Una lista de DTOs con los registros encontrados.
  - Un error en caso de ocurrir alguno.
*/
func GetRelCaseOwnerFeminicideRisks(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelCaseOwnerFeminicideRiskDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción de la ruta completa de la tabla en la base de datos.
	var relCaseOwnerFeminicideRiskPath string = RelCaseOwnerFeminicideRiskDBScheme + "." + RelCaseOwnerFeminicideRiskDBName
	var relCaseOwnerFeminicideRisks []RelCaseOwnerFeminicideRiskDTO = []RelCaseOwnerFeminicideRiskDTO{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return relCaseOwnerFeminicideRisks, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar en la consulta.
	var relCaseOwnerFeminicideRiskFieldsSlice []string = []string{"RelCaseOwnerFeminicideRiskId", "RelCaseOwnerFeminicideRisk_CaseOwner", "RelCaseOwnerFeminicideRisk_FeminicideRisk", "RelCaseOwnerFeminicideRisk_CreationDate", "RelCaseOwnerFeminicideRisk_Status"}
	var relCaseOwnerFeminicideRiskFieldsAliasSlice []string = []string{}

	// Construcción de la lista de campos para la consulta SQL.
	var relCaseOwnerFeminicideRiskFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY,
		relCaseOwnerFeminicideRiskFieldsSlice,
		relCaseOwnerFeminicideRiskFieldsAliasSlice,
		RelCaseOwnerFeminicideRiskDBName,
		[]string{}, []string{}, []string{},
		common_dao.SQL_AND,
		RelCaseOwnerFeminicideRiskDBScheme,
		RelCaseOwnerFeminicideRiskFieldDefinitions,
		true)

	// Construcción de la consulta SQL completa con cláusula WHERE basada en los criterios de filtrado.
	var query string = `SELECT ` + relCaseOwnerFeminicideRiskFieldsStr +
		` FROM ` + relCaseOwnerFeminicideRiskPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY,
			by.AttrsName,
			by.AttrsAliasName,
			RelCaseOwnerFeminicideRiskDBName,
			by.AttrsName,
			[]string{}, []string{},
			by.Operator,
			RelCaseOwnerFeminicideRiskDBScheme,
			RelCaseOwnerFeminicideRiskFieldDefinitions,
			true)

	// Ejecución de la consulta con los parámetros de filtrado.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre los resultados y los mapea al DTO correspondiente.
	for persistenceCtrl.Next() {
		var relCaseOwnerFeminicideRisk RelCaseOwnerFeminicideRiskDTO = RelCaseOwnerFeminicideRiskDTO{}

		persistenceCtrl.ScanRow(&relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRiskId,
			&relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_CaseOwner,
			&relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_FeminicideRisk,
			&relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_CreationDate,
			&relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_Status)

		relCaseOwnerFeminicideRisks = append(relCaseOwnerFeminicideRisks, relCaseOwnerFeminicideRisk)
	}

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return relCaseOwnerFeminicideRisks, persistenceCtrl.Error
	}

	return relCaseOwnerFeminicideRisks, nil
}

/*
UpdateRelCaseOwnerFeminicideRisk actualiza el campo "status" de un registro de RelCaseOwnerFeminicideRisk en la base de datos.

Parámetros:
  - relCaseOwnerFeminicideRisk: Puntero al DTO con los datos actualizados.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:

  - Un error en caso de ocurrir alguno.
*/
func UpdateRelCaseOwnerFeminicideRisk(relCaseOwnerFeminicideRisk *RelCaseOwnerFeminicideRiskDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a actualizar (en este caso, solo el campo "status").
	var relCaseOwnerFeminicideRiskFieldsSlice []string = []string{"RelCaseOwnerFeminicideRisk_Status"}
	var relCaseOwnerFeminicideRiskFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE,
		relCaseOwnerFeminicideRiskFieldsSlice,
		relCaseOwnerFeminicideRiskFieldsAliasSlice,
		RelCaseOwnerFeminicideRiskDBName,
		[]string{"RelCaseOwnerFeminicideRiskId"},
		[]string{}, []string{},
		common_dao.SQL_AND,
		RelCaseOwnerFeminicideRiskDBScheme,
		RelCaseOwnerFeminicideRiskFieldDefinitions,
		false)

	// Ejecución de la consulta SQL con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRiskId,
		relCaseOwnerFeminicideRisk.RelCaseOwnerFeminicideRisk_Status)

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica si se actualizó al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

/*
SetRelCaseOwnerFeminicideRiskDefaults establece valores predeterminados en el DTO de RelCaseOwnerFeminicideRisk
según la acción a realizar (inserción o actualización).

Parámetros:
  - relEntity: Puntero al DTO que se modificará.
  - action: Acción que se va a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
*/
func SetRelCaseOwnerFeminicideRiskDefaults(relEntity *RelCaseOwnerFeminicideRiskDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción, se establece la fecha de creación actual y el estado activo ("a").
		relEntity.RelCaseOwnerFeminicideRisk_CreationDate = time.Now()
		relEntity.RelCaseOwnerFeminicideRisk_Status = "a"
	case common_dao.SQL_UPDATE:
		// Para actualización, se puede cambiar el estado a inactivo ("i").
		relEntity.RelCaseOwnerFeminicideRisk_Status = "i"
	}
}

/*
ToDTO convierte una estructura de la base de datos (RelCaseOwnerFeminicideRiskPgDB) en su correspondiente DTO (RelCaseOwnerFeminicideRiskDTO).
Maneja la validación de campos nulos y realiza la conversión de tipos según corresponda.
*/
func (obj *RelCaseOwnerFeminicideRiskPgDB) ToDTO() RelCaseOwnerFeminicideRiskDTO {
	var dto RelCaseOwnerFeminicideRiskDTO

	if obj.RelCaseOwnerFeminicideRiskId.Valid {
		dto.RelCaseOwnerFeminicideRiskId = uint64(obj.RelCaseOwnerFeminicideRiskId.Int64)
	}

	if obj.RelCaseOwnerFeminicideRisk_CaseOwner.Valid {
		dto.RelCaseOwnerFeminicideRisk_CaseOwner = uint64(obj.RelCaseOwnerFeminicideRisk_CaseOwner.Int64)
	}

	if obj.RelCaseOwnerFeminicideRisk_FeminicideRisk.Valid {
		dto.RelCaseOwnerFeminicideRisk_FeminicideRisk = uint64(obj.RelCaseOwnerFeminicideRisk_FeminicideRisk.Int64)
	}

	if obj.RelCaseOwnerFeminicideRisk_CreationDate.Valid {
		dto.RelCaseOwnerFeminicideRisk_CreationDate = obj.RelCaseOwnerFeminicideRisk_CreationDate.Time
	}

	if obj.RelCaseOwnerFeminicideRisk_Status.Valid {
		dto.RelCaseOwnerFeminicideRisk_Status = obj.RelCaseOwnerFeminicideRisk_Status.String
	}

	return dto
}
