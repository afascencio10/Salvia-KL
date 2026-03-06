// Package salvia_daos contiene los Data Access Objects (DAO) para la entidad RelCaseOwnerVictimCase,
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
	// Constantes y nombres relacionados con la entidad RelCaseOwnerVictimCase.
	RelCaseOwnerVictimCase         string = "RelCaseOwnerVictimCase"     // Nombre interno de la entidad.
	RelCaseOwnerVictimCaseJSONName string = "relCaseOwner"               // Nombre de la entidad en el JSON.
	RelCaseOwnerVictimCaseDBName   string = "rel_case_owner_victim_case" // Nombre de la tabla en la base de datos.
	RelCaseOwnerVictimCaseDBScheme string = "salvia"                     // Esquema de la base de datos.

	// Atributos relacionados con las validaciones de campos.
	// Cada entrada del mapa define la relación entre el campo del JSON, su nombre en la base de datos,
	// el tipo de dato esperado en el modelo y otros parámetros de validación.
	RelCaseOwnerVictimCaseFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelCaseOwnerVictimCaseId":            {Name: "RelCaseOwnerVictimCaseId", DBName: "rel_case_owner_victim_case_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelCaseOwnerVictimCase_CaseOwner":    {Name: "RelCaseOwnerVictimCase_CaseOwner", DBName: "case_owner_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerVictimCase_VictimCase":   {Name: "RelCaseOwnerVictimCase_VictimCase", DBName: "victim_case_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerVictimCase_CreationDate": {Name: "RelCaseOwnerVictimCase_CreationDate", DBName: "rel_case_owner_victim_case_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerVictimCase_Status":       {Name: "RelCaseOwnerVictimCase_Status", DBName: "rel_case_owner_victim_case_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
	}
)

/*
RelCaseOwnerVictimCaseDTO representa la estructura de datos (DTO) para la entidad RelCaseOwnerVictimCase.
Este DTO se utiliza para transferir datos entre las capas de la aplicación.
*/
type RelCaseOwnerVictimCaseDTO struct {
	RelCaseOwnerVictimCaseId            uint64    `json:"-"`
	RelCaseOwnerVictimCase_CaseOwner    uint64    `json:"caseOwner"`
	RelCaseOwnerVictimCase_VictimCase   uint64    `json:"victimCase"`
	RelCaseOwnerVictimCase_CreationDate time.Time `json:"creationDate"`
	RelCaseOwnerVictimCase_Status       string    `json:"status"`
}

/*
RelCaseOwnerVictimCasePgDB representa la estructura que mapea la entidad en la base de datos.
Utiliza tipos sql.Null* para manejar valores nulos en la base de datos.
*/
type RelCaseOwnerVictimCasePgDB struct {
	RelCaseOwnerVictimCaseId            sql.NullInt64
	RelCaseOwnerVictimCase_CaseOwner    sql.NullInt64
	RelCaseOwnerVictimCase_VictimCase   sql.NullInt64
	RelCaseOwnerVictimCase_CreationDate sql.NullTime
	RelCaseOwnerVictimCase_Status       sql.NullString
}

func (rcwvc RelCaseOwnerVictimCaseDTO) MarshalJSON() ([]byte, error) {
	type Alias RelCaseOwnerVictimCaseDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		RelCaseOwnerVictimCase_CreationDate string `json:"creationDate"`
	}{
		Alias:                               (*Alias)(&rcwvc),
		RelCaseOwnerVictimCase_CreationDate: rcwvc.RelCaseOwnerVictimCase_CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (rcwvc *RelCaseOwnerVictimCaseDTO) UnmarshalJSON(data []byte) error {
	type Alias RelCaseOwnerVictimCaseDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		RelCaseOwnerVictimCase_CreationDate string `json:"creationDate"`
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
	rcwvc.RelCaseOwnerVictimCase_CreationDate = parse(aux.RelCaseOwnerVictimCase_CreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

/*
SetRelCaseOwnerVictimCase inserta un nuevo registro de RelCaseOwnerVictimCase en la base de datos.

Parámetros:
  - relCaseOwnerVictimCase: Puntero al DTO que contiene los datos a insertar.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:
  - Un error en caso de ocurrir alguno.
*/
func SetRelCaseOwnerVictimCase(relCaseOwnerVictimCase *RelCaseOwnerVictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a insertar en la base de datos.
	var relCaseOwnerVictimCaseFieldsSlice []string = []string{"RelCaseOwnerVictimCase_CaseOwner", "RelCaseOwnerVictimCase_VictimCase", "RelCaseOwnerVictimCase_CreationDate", "RelCaseOwnerVictimCase_Status"}
	var relCaseOwnerVictimCaseFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, relCaseOwnerVictimCaseFieldsSlice, relCaseOwnerVictimCaseFieldsAliasSlice, RelCaseOwnerVictimCaseDBName, []string{}, []string{}, []string{"RelCaseOwnerVictimCaseId"}, common_dao.SQL_AND, RelCaseOwnerVictimCaseDBScheme, RelCaseOwnerVictimCaseFieldDefinitions, false)

	// Ejecución de la consulta SQL y asignación del ID generado al DTO.
	persistenceCtrl.QueryRow(context.Background(), query,
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner,
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase,
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_Status)
	persistenceCtrl.Scan(&relCaseOwnerVictimCase.RelCaseOwnerVictimCaseId)

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

/*
GetRelCaseOwnerVictimCases recupera registros de RelCaseOwnerVictimCase desde la base de datos
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
func GetRelCaseOwnerVictimCases(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelCaseOwnerVictimCaseDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción de la ruta completa de la tabla en la base de datos.
	var relCaseOwnerVictimCasePath string = RelCaseOwnerVictimCaseDBScheme + "." + RelCaseOwnerVictimCaseDBName
	var relCaseOwnerVictimCases []RelCaseOwnerVictimCaseDTO = []RelCaseOwnerVictimCaseDTO{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return relCaseOwnerVictimCases, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar en la consulta.
	var relCaseOwnerVictimCaseFieldsSlice []string = []string{"RelCaseOwnerVictimCaseId", "RelCaseOwnerVictimCase_CaseOwner", "RelCaseOwnerVictimCase_VictimCase", "RelCaseOwnerVictimCase_CreationDate", "RelCaseOwnerVictimCase_Status"}
	var relCaseOwnerVictimCaseFieldsAliasSlice []string = []string{}

	// Construcción de la lista de campos para la consulta SQL.
	var relCaseOwnerVictimCaseFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY,
		relCaseOwnerVictimCaseFieldsSlice,
		relCaseOwnerVictimCaseFieldsAliasSlice,
		RelCaseOwnerVictimCaseDBName,
		[]string{}, []string{}, []string{},
		common_dao.SQL_AND,
		RelCaseOwnerVictimCaseDBScheme,
		RelCaseOwnerVictimCaseFieldDefinitions,
		true)

	// Construcción de la consulta SQL completa con cláusula WHERE basada en los criterios de filtrado.
	var query string = `SELECT ` + relCaseOwnerVictimCaseFieldsStr +
		` FROM ` + relCaseOwnerVictimCasePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY,
			by.AttrsName,
			by.AttrsAliasName,
			RelCaseOwnerVictimCaseDBName,
			by.AttrsName,
			[]string{}, []string{},
			by.Operator,
			RelCaseOwnerVictimCaseDBScheme,
			RelCaseOwnerVictimCaseFieldDefinitions,
			true)

	// Ejecución de la consulta con los parámetros de filtrado.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre los resultados y los mapea al DTO correspondiente.
	for persistenceCtrl.Next() {
		var relCaseOwnerVictimCase RelCaseOwnerVictimCaseDTO = RelCaseOwnerVictimCaseDTO{}

		persistenceCtrl.ScanRow(&relCaseOwnerVictimCase.RelCaseOwnerVictimCaseId,
			&relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CaseOwner,
			&relCaseOwnerVictimCase.RelCaseOwnerVictimCase_VictimCase,
			&relCaseOwnerVictimCase.RelCaseOwnerVictimCase_CreationDate,
			&relCaseOwnerVictimCase.RelCaseOwnerVictimCase_Status)

		relCaseOwnerVictimCases = append(relCaseOwnerVictimCases, relCaseOwnerVictimCase)
	}

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return relCaseOwnerVictimCases, persistenceCtrl.Error
	}

	return relCaseOwnerVictimCases, nil
}

/*
UpdateRelCaseOwnerVictimCase actualiza el campo "status" de un registro de RelCaseOwnerVictimCase en la base de datos.

Parámetros:
  - relCaseOwnerVictimCase: Puntero al DTO con los datos actualizados.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:

  - Un error en caso de ocurrir alguno.
*/
func UpdateRelCaseOwnerVictimCase(relCaseOwnerVictimCase *RelCaseOwnerVictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a actualizar (en este caso, solo el campo "status").
	var relCaseOwnerVictimCaseFieldsSlice []string = []string{"RelCaseOwnerVictimCase_Status"}
	var relCaseOwnerVictimCaseFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE,
		relCaseOwnerVictimCaseFieldsSlice,
		relCaseOwnerVictimCaseFieldsAliasSlice,
		RelCaseOwnerVictimCaseDBName,
		[]string{"RelCaseOwnerVictimCaseId"},
		[]string{}, []string{},
		common_dao.SQL_AND,
		RelCaseOwnerVictimCaseDBScheme,
		RelCaseOwnerVictimCaseFieldDefinitions,
		false)

	// Ejecución de la consulta SQL con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		relCaseOwnerVictimCase.RelCaseOwnerVictimCaseId,
		relCaseOwnerVictimCase.RelCaseOwnerVictimCase_Status)

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
SetRelCaseOwnerVictimCaseDefaults establece valores predeterminados en el DTO de RelCaseOwnerVictimCase
según la acción a realizar (inserción o actualización).

Parámetros:
  - relEntity: Puntero al DTO que se modificará.
  - action: Acción que se va a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
*/
func SetRelCaseOwnerVictimCaseDefaults(relEntity *RelCaseOwnerVictimCaseDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción, se establece la fecha de creación actual y el estado activo ("a").
		relEntity.RelCaseOwnerVictimCase_CreationDate = time.Now()
		relEntity.RelCaseOwnerVictimCase_Status = "a"
	case common_dao.SQL_UPDATE:
		// Para actualización, se puede cambiar el estado a inactivo ("i").
		relEntity.RelCaseOwnerVictimCase_Status = "i"
	}
}

/*
ToDTO convierte una estructura de la base de datos (RelCaseOwnerVictimCasePgDB) en su correspondiente DTO (RelCaseOwnerVictimCaseDTO).
Maneja la validación de campos nulos y realiza la conversión de tipos según corresponda.
*/
func (obj *RelCaseOwnerVictimCasePgDB) ToDTO() RelCaseOwnerVictimCaseDTO {
	var dto RelCaseOwnerVictimCaseDTO

	if obj.RelCaseOwnerVictimCaseId.Valid {
		dto.RelCaseOwnerVictimCaseId = uint64(obj.RelCaseOwnerVictimCaseId.Int64)
	}

	if obj.RelCaseOwnerVictimCase_CaseOwner.Valid {
		dto.RelCaseOwnerVictimCase_CaseOwner = uint64(obj.RelCaseOwnerVictimCase_CaseOwner.Int64)
	}

	if obj.RelCaseOwnerVictimCase_VictimCase.Valid {
		dto.RelCaseOwnerVictimCase_VictimCase = uint64(obj.RelCaseOwnerVictimCase_VictimCase.Int64)
	}

	if obj.RelCaseOwnerVictimCase_CreationDate.Valid {
		dto.RelCaseOwnerVictimCase_CreationDate = obj.RelCaseOwnerVictimCase_CreationDate.Time
	}

	if obj.RelCaseOwnerVictimCase_Status.Valid {
		dto.RelCaseOwnerVictimCase_Status = obj.RelCaseOwnerVictimCase_Status.String
	}

	return dto
}
