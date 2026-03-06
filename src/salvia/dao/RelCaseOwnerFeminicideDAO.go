// Package salvia_daos contiene los Data Access Objects (DAO) para la entidad RelCaseOwnerFeminicide,
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
	// Constantes y nombres relacionados con la entidad RelCaseOwnerFeminicide.
	RelCaseOwnerFeminicide         string = "RelCaseOwnerFeminicide"    // Nombre interno de la entidad.
	RelCaseOwnerFeminicideJSONName string = "relCaseOwner"              // Nombre de la entidad en el JSON.
	RelCaseOwnerFeminicideDBName   string = "rel_case_owner_feminicide" // Nombre de la tabla en la base de datos.
	RelCaseOwnerFeminicideDBScheme string = "salvia"                    // Esquema de la base de datos.

	// Atributos relacionados con las validaciones de campos.
	// Cada entrada del mapa define la relación entre el campo del JSON, su nombre en la base de datos,
	// el tipo de dato esperado en el modelo y otros parámetros de validación.
	RelCaseOwnerFeminicideFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelCaseOwnerFeminicideId":            {Name: "RelCaseOwnerFeminicideId", DBName: "rel_case_owner_feminicide_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelCaseOwnerFeminicide_CaseOwner":    {Name: "RelCaseOwnerFeminicide_CaseOwner", DBName: "case_owner_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerFeminicide_Feminicide":   {Name: "RelCaseOwnerFeminicide_Feminicide", DBName: "feminicide_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerFeminicide_CreationDate": {Name: "RelCaseOwnerFeminicide_CreationDate", DBName: "rel_case_owner_feminicide_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"RelCaseOwnerFeminicide_Status":       {Name: "RelCaseOwnerFeminicide_Status", DBName: "rel_case_owner_feminicide_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
	}
)

/*
RelCaseOwnerFeminicideDTO representa la estructura de datos (DTO) para la entidad RelCaseOwnerFeminicide.
Este DTO se utiliza para transferir datos entre las capas de la aplicación.
*/
type RelCaseOwnerFeminicideDTO struct {
	RelCaseOwnerFeminicideId            uint64    `json:"-"`
	RelCaseOwnerFeminicide_CaseOwner    uint64    `json:"caseOwner"`
	RelCaseOwnerFeminicide_Feminicide   uint64    `json:"feminicide"`
	RelCaseOwnerFeminicide_CreationDate time.Time `json:"creationDate"`
	RelCaseOwnerFeminicide_Status       string    `json:"status"`
}

/*
RelCaseOwnerFeminicidePgDB representa la estructura que mapea la entidad en la base de datos.
Utiliza tipos sql.Null* para manejar valores nulos en la base de datos.
*/
type RelCaseOwnerFeminicidePgDB struct {
	RelCaseOwnerFeminicideId            sql.NullInt64
	RelCaseOwnerFeminicide_CaseOwner    sql.NullInt64
	RelCaseOwnerFeminicide_Feminicide   sql.NullInt64
	RelCaseOwnerFeminicide_CreationDate sql.NullTime
	RelCaseOwnerFeminicide_Status       sql.NullString
}

func (rcwvc RelCaseOwnerFeminicideDTO) MarshalJSON() ([]byte, error) {
	type Alias RelCaseOwnerFeminicideDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		RelCaseOwnerFeminicide_CreationDate string `json:"creationDate"`
	}{
		Alias:                               (*Alias)(&rcwvc),
		RelCaseOwnerFeminicide_CreationDate: rcwvc.RelCaseOwnerFeminicide_CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (rcwvc *RelCaseOwnerFeminicideDTO) UnmarshalJSON(data []byte) error {
	type Alias RelCaseOwnerFeminicideDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		RelCaseOwnerFeminicide_CreationDate string `json:"creationDate"`
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
	rcwvc.RelCaseOwnerFeminicide_CreationDate = parse(aux.RelCaseOwnerFeminicide_CreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

/*
SetRelCaseOwnerFeminicide inserta un nuevo registro de RelCaseOwnerFeminicide en la base de datos.

Parámetros:
  - relCaseOwnerFeminicide: Puntero al DTO que contiene los datos a insertar.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:
  - Un error en caso de ocurrir alguno.
*/
func SetRelCaseOwnerFeminicide(relCaseOwnerFeminicide *RelCaseOwnerFeminicideDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a insertar en la base de datos.
	var relCaseOwnerFeminicideFieldsSlice []string = []string{"RelCaseOwnerFeminicide_CaseOwner", "RelCaseOwnerFeminicide_Feminicide", "RelCaseOwnerFeminicide_CreationDate", "RelCaseOwnerFeminicide_Status"}
	var relCaseOwnerFeminicideFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de inserción.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, relCaseOwnerFeminicideFieldsSlice, relCaseOwnerFeminicideFieldsAliasSlice, RelCaseOwnerFeminicideDBName, []string{}, []string{}, []string{"RelCaseOwnerFeminicideId"}, common_dao.SQL_AND, RelCaseOwnerFeminicideDBScheme, RelCaseOwnerFeminicideFieldDefinitions, false)

	// Ejecución de la consulta SQL y asignación del ID generado al DTO.
	persistenceCtrl.QueryRow(context.Background(), query,
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_CaseOwner,
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Feminicide,
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Status)
	persistenceCtrl.Scan(&relCaseOwnerFeminicide.RelCaseOwnerFeminicideId)

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

/*
GetRelCaseOwnerFeminicides recupera registros de RelCaseOwnerFeminicide desde la base de datos
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
func GetRelCaseOwnerFeminicides(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelCaseOwnerFeminicideDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción de la ruta completa de la tabla en la base de datos.
	var relCaseOwnerFeminicidePath string = RelCaseOwnerFeminicideDBScheme + "." + RelCaseOwnerFeminicideDBName
	var relCaseOwnerFeminicides []RelCaseOwnerFeminicideDTO = []RelCaseOwnerFeminicideDTO{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return relCaseOwnerFeminicides, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar en la consulta.
	var relCaseOwnerFeminicideFieldsSlice []string = []string{"RelCaseOwnerFeminicideId", "RelCaseOwnerFeminicide_CaseOwner", "RelCaseOwnerFeminicide_Feminicide", "RelCaseOwnerFeminicide_CreationDate", "RelCaseOwnerFeminicide_Status"}
	var relCaseOwnerFeminicideFieldsAliasSlice []string = []string{}

	// Construcción de la lista de campos para la consulta SQL.
	var relCaseOwnerFeminicideFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY,
		relCaseOwnerFeminicideFieldsSlice,
		relCaseOwnerFeminicideFieldsAliasSlice,
		RelCaseOwnerFeminicideDBName,
		[]string{}, []string{}, []string{},
		common_dao.SQL_AND,
		RelCaseOwnerFeminicideDBScheme,
		RelCaseOwnerFeminicideFieldDefinitions,
		true)

	// Construcción de la consulta SQL completa con cláusula WHERE basada en los criterios de filtrado.
	var query string = `SELECT ` + relCaseOwnerFeminicideFieldsStr +
		` FROM ` + relCaseOwnerFeminicidePath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY,
			by.AttrsName,
			by.AttrsAliasName,
			RelCaseOwnerFeminicideDBName,
			by.AttrsName,
			[]string{}, []string{},
			by.Operator,
			RelCaseOwnerFeminicideDBScheme,
			RelCaseOwnerFeminicideFieldDefinitions,
			true)

	// Ejecución de la consulta con los parámetros de filtrado.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre los resultados y los mapea al DTO correspondiente.
	for persistenceCtrl.Next() {
		var relCaseOwnerFeminicide RelCaseOwnerFeminicideDTO = RelCaseOwnerFeminicideDTO{}

		persistenceCtrl.ScanRow(&relCaseOwnerFeminicide.RelCaseOwnerFeminicideId,
			&relCaseOwnerFeminicide.RelCaseOwnerFeminicide_CaseOwner,
			&relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Feminicide,
			&relCaseOwnerFeminicide.RelCaseOwnerFeminicide_CreationDate,
			&relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Status)

		relCaseOwnerFeminicides = append(relCaseOwnerFeminicides, relCaseOwnerFeminicide)
	}

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return relCaseOwnerFeminicides, persistenceCtrl.Error
	}

	return relCaseOwnerFeminicides, nil
}

/*
UpdateRelCaseOwnerFeminicide actualiza el campo "status" de un registro de RelCaseOwnerFeminicide en la base de datos.

Parámetros:
  - relCaseOwnerFeminicide: Puntero al DTO con los datos actualizados.
  - connData: Datos de conexión actuales.
  - clientConfig: Configuración del cliente de base de datos.
  - serverConfig: Configuración del servidor de base de datos.

Retorna:

  - Un error en caso de ocurrir alguno.
*/
func UpdateRelCaseOwnerFeminicide(relCaseOwnerFeminicide *RelCaseOwnerFeminicideDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a actualizar (en este caso, solo el campo "status").
	var relCaseOwnerFeminicideFieldsSlice []string = []string{"RelCaseOwnerFeminicide_Status"}
	var relCaseOwnerFeminicideFieldsAliasSlice []string = []string{}

	// Construcción de la consulta SQL de actualización.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE,
		relCaseOwnerFeminicideFieldsSlice,
		relCaseOwnerFeminicideFieldsAliasSlice,
		RelCaseOwnerFeminicideDBName,
		[]string{"RelCaseOwnerFeminicideId"},
		[]string{}, []string{},
		common_dao.SQL_AND,
		RelCaseOwnerFeminicideDBScheme,
		RelCaseOwnerFeminicideFieldDefinitions,
		false)

	// Ejecución de la consulta SQL con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		relCaseOwnerFeminicide.RelCaseOwnerFeminicideId,
		relCaseOwnerFeminicide.RelCaseOwnerFeminicide_Status)

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
SetRelCaseOwnerFeminicideDefaults establece valores predeterminados en el DTO de RelCaseOwnerFeminicide
según la acción a realizar (inserción o actualización).

Parámetros:
  - relEntity: Puntero al DTO que se modificará.
  - action: Acción que se va a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
*/
func SetRelCaseOwnerFeminicideDefaults(relEntity *RelCaseOwnerFeminicideDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción, se establece la fecha de creación actual y el estado activo ("a").
		relEntity.RelCaseOwnerFeminicide_CreationDate = time.Now()
		relEntity.RelCaseOwnerFeminicide_Status = "a"
	case common_dao.SQL_UPDATE:
		// Para actualización, se puede cambiar el estado a inactivo ("i").
		relEntity.RelCaseOwnerFeminicide_Status = "i"
	}
}

/*
ToDTO convierte una estructura de la base de datos (RelCaseOwnerFeminicidePgDB) en su correspondiente DTO (RelCaseOwnerFeminicideDTO).
Maneja la validación de campos nulos y realiza la conversión de tipos según corresponda.
*/
func (obj *RelCaseOwnerFeminicidePgDB) ToDTO() RelCaseOwnerFeminicideDTO {
	var dto RelCaseOwnerFeminicideDTO

	if obj.RelCaseOwnerFeminicideId.Valid {
		dto.RelCaseOwnerFeminicideId = uint64(obj.RelCaseOwnerFeminicideId.Int64)
	}

	if obj.RelCaseOwnerFeminicide_CaseOwner.Valid {
		dto.RelCaseOwnerFeminicide_CaseOwner = uint64(obj.RelCaseOwnerFeminicide_CaseOwner.Int64)
	}

	if obj.RelCaseOwnerFeminicide_Feminicide.Valid {
		dto.RelCaseOwnerFeminicide_Feminicide = uint64(obj.RelCaseOwnerFeminicide_Feminicide.Int64)
	}

	if obj.RelCaseOwnerFeminicide_CreationDate.Valid {
		dto.RelCaseOwnerFeminicide_CreationDate = obj.RelCaseOwnerFeminicide_CreationDate.Time
	}

	if obj.RelCaseOwnerFeminicide_Status.Valid {
		dto.RelCaseOwnerFeminicide_Status = obj.RelCaseOwnerFeminicide_Status.String
	}

	return dto
}
