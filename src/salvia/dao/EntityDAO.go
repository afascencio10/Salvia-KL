// Package salvia_daos contiene las funciones DAO para gestionar operaciones
// relacionadas con la entidad "Entity". Provee métodos para obtener, listar y
// establecer valores por defecto en los registros de la entidad, facilitando
// la interacción con la base de datos.
package salvia_daos

import (
	// Controladores comunes para operaciones de persistencia.
	// Funciones comunes para la generación de SQL.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"    // Conexión y configuración de la base de datos.
	"bitsflow/common/utils" // Utilidades generales, como generación de UUID y definiciones de campos.
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

var (
	// EntityEntityName es el nombre interno de la entidad.
	EntityEntityName string = "Entity"
	// EntityJSONName es el nombre que se utiliza en las respuestas JSON.
	EntityJSONName string = "entity"
	// EntityDBName es el nombre de la tabla en la base de datos.
	EntityDBName string = "entity"
	// EntityDBScheme es el esquema de la base de datos en el que se encuentra la entidad.
	EntityDBScheme string = "salvia"

	// EntityFieldDefinitions define los atributos de la entidad y las validaciones asociadas.
	// Cada campo incluye su nombre interno, nombre en la BD, alias, tipo de modelo,
	// tamaños mínimos y máximos y si es obligatorio o no.
	EntityFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"EntityId":                   {Name: "EntityId", DBName: "entity_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"EntityICode":                {Name: "EntityICode", DBName: "entity_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"EntityCreationDate":         {Name: "EntityCreationDate", DBName: "entity_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"EntityUpdateDate":           {Name: "EntityUpdateDate", DBName: "entity_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"EntityName":                 {Name: "EntityName", DBName: "entity_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 254, Required: true},
		"EntityDescription":          {Name: "EntityName", DBName: "entity_name", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 254, Required: false},
		"EntityIsInteroperable":      {Name: "EntityIsInteroperable", DBName: "entity_is_interoperable", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"EntityInteroperabilityCode": {Name: "EntityInteroperabilityCode", DBName: "entity_interoperability_code", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 32, Required: true},
		"EntityResponseTime":         {Name: "EntityResponseTime", DBName: "entity_response_time", Alias: "", ModelType: "uint", Required: true},
		"EntitySector":               {Name: "EntitySector", DBName: "entity_sector", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
	}
)

// EntityDTO representa la estructura de datos de la entidad que se utiliza para
// transferir información entre las capas de la aplicación (por ejemplo, de la BD al API).
type EntityDTO struct {
	EntityId                   uint64    `json:"-"`
	EntityICode                string    `json:"icode"`
	EntityCreationDate         time.Time `json:"-"`
	EntityUpdateDate           time.Time `json:"-"`
	EntityName                 string    `json:"name"`
	EntityDescription          string    `json:"description"`
	EntityIsInteroperable      string    `json:"-"`
	EntityInteroperabilityCode string    `json:"code"`
	EntityResponseTime         uint64    `json:"-"`
	EntitySector               string    `json:"sector"`
	// EntityMoment es un campo adicional que se utiliza en formularios y no forma parte
	// directamente del modelo en la base de datos.
	EntityMoment string `json:"moment"`
}

// EntityPgDB representa la estructura de la entidad tal como se almacena en la base
// de datos, utilizando tipos nulos de SQL para manejar valores que pueden ser NULL.
type EntityPgDB struct {
	EntityId                   sql.NullInt64
	EntityICode                sql.NullString
	EntityCreationDate         sql.NullTime
	EntityUpdateDate           sql.NullTime
	EntityName                 sql.NullString
	EntityDescription          sql.NullString
	EntityIsInteroperable      sql.NullString
	EntityInteroperabilityCode sql.NullString
	EntityResponseTime         sql.NullInt64
	EntitySector               sql.NullString
}

func (e EntityDTO) MarshalJSON() ([]byte, error) {
	type Alias EntityDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		EntityCreationDate string `json:"creationDate"`
		EntityUpdateDate   string `json:"updateDate"`
	}{
		Alias:              (*Alias)(&e),
		EntityCreationDate: e.EntityCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		EntityUpdateDate:   e.EntityUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (e *EntityDTO) UnmarshalJSON(data []byte) error {
	type Alias EntityDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		EntityCreationDate string `json:"creationDate"`
		EntityUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(e),
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
	e.EntityCreationDate = parse(aux.EntityCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	e.EntityUpdateDate = parse(aux.EntityUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// GetEntity recupera un registro de la entidad desde la base de datos utilizando
// un filtro especificado en el parámetro "by".
//
// Parámetros:
//   - by: estructura que contiene los atributos y valores para filtrar la consulta.
//   - entity: puntero a EntityDTO donde se almacenará el resultado.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - El error encontrado.
func GetEntity(by common_controllers.By, entity *EntityDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construye el path completo de la tabla (esquema.tabla).
	var entityPath string = EntityDBScheme + "." + EntityDBName

	// Establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos que se seleccionarán en la consulta.
	var entityFieldsSlice []string = []string{
		"EntityId", "EntityICode", "EntityCreationDate", "EntityUpdateDate",
		"EntityName", "EntityDescription", "EntityIsInteroperable",
		"EntityInteroperabilityCode", "EntityResponseTime", "EntitySector",
	}
	var entityFieldsAliasSlice []string = []string{}

	// Genera la parte del SELECT con los campos definidos.
	var entityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityFieldsSlice, entityFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityDBScheme, EntityFieldDefinitions, true)

	// Construye la consulta SQL utilizando la función GetSQL para el WHERE.
	var query string = `SELECT ` + entityFieldsStr +
		` FROM ` + entityPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, EntityDBName, by.AttrsName, []string{}, []string{}, by.Operator, EntityDBScheme, EntityFieldDefinitions, true)

	// Ejecuta la consulta y mapea el resultado en un objeto temporal.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)
	var entityPG EntityPgDB = EntityPgDB{}
	persistenceCtrl.Scan(&entityPG.EntityId, &entityPG.EntityICode, &entityPG.EntityCreationDate, &entityPG.EntityUpdateDate, &entityPG.EntityName, &entityPG.EntityDescription, &entityPG.EntityIsInteroperable, &entityPG.EntityInteroperabilityCode,
		&entityPG.EntityResponseTime, &entityPG.EntitySector)

	// Convierte el resultado obtenido a DTO y lo asigna al parámetro de salida.
	var tmp EntityDTO = entityPG.ToDTO()
	*entity = tmp

	// Manejo de errores en la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetEntities recupera múltiples registros de la entidad según los criterios de filtrado
// especificados en el parámetro "by".
//
// Parámetros:
//   - by: criterios de filtrado para la consulta.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - un slice de EntityDTO con los registros encontrados y,
//     en caso de error, el error correspondiente.
func GetEntities(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entityPath string = EntityDBScheme + "." + EntityDBName

	var entities []EntityDTO

	// Establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar en la consulta.
	var entityFieldsSlice []string = []string{
		"EntityICode", "EntityCreationDate", "EntityUpdateDate", "EntityName",
		"EntityDescription", "EntityIsInteroperable", "EntityInteroperabilityCode",
		"EntityResponseTime", "EntitySector",
	}
	var entityFieldsAliasSlice []string = []string{}

	// Genera el string de campos para el SELECT.
	var entityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityFieldsSlice, entityFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityDBScheme, EntityFieldDefinitions, true)

	// Construye la consulta SQL incluyendo cláusulas WHERE y ORDER BY.
	var query string = `SELECT ` + entityFieldsStr +
		` FROM ` + entityPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, EntityDBName, by.AttrsName, []string{}, []string{}, by.Operator, EntityDBScheme, EntityFieldDefinitions, true) +
		` ORDER BY ` + entityPath + "." + EntityFieldDefinitions["EntityName"].DBName + ` ASC `

	// Ejecuta la consulta y obtiene el cursor para iterar sobre los resultados.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	for persistenceCtrl.Next() {
		var entityPg EntityPgDB = EntityPgDB{}
		// Mapea cada fila en la estructura EntityDTO.
		persistenceCtrl.ScanRow(&entityPg.EntityICode, &entityPg.EntityCreationDate, &entityPg.EntityUpdateDate, &entityPg.EntityName, &entityPg.EntityDescription,
			&entityPg.EntityIsInteroperable, &entityPg.EntityInteroperabilityCode, &entityPg.EntityResponseTime, &entityPg.EntitySector)
		entities = append(entities, entityPg.ToDTO())
	}

	// Manejo de errores durante la iteración.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: Aplicar esta metodología de traer 1 o más elementos con un "by" para todos los DAOS.
	return entities, nil
}

// GetAllEntites recupera todos los registros de la entidad sin aplicar filtros específicos.
//
// Parámetros:
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - un slice de EntityDTO con todos los registros y,
//     en caso de error, el error correspondiente.
func GetAllEntites(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityDTO, error) {
	// Declaración de variables.
	var entity EntityPgDB
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var entityFieldsSlice []string = []string{
		"EntityId", "EntityICode", "EntityCreationDate", "EntityUpdateDate",
		"EntityName", "EntityDescription", "EntityIsInteroperable",
		"EntityInteroperabilityCode", "EntityResponseTime", "EntitySector",
	}
	var entityFieldsAliasSlice []string = []string{}

	// Construye la consulta SQL para seleccionar todos los registros.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, entityFieldsSlice, entityFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, "", EntityDBScheme, EntityFieldDefinitions, true) +
		` ORDER BY ` + EntityDBScheme + "." + EntityDBName + `.` + EntityFieldDefinitions["EntityName"].DBName + ` ASC `

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var entities []EntityDTO
	for persistenceCtrl.Next() {
		entity = EntityPgDB{}
		// Mapea cada fila de la consulta en la estructura EntityPgDB.
		persistenceCtrl.ScanRow(&entity.EntityId, &entity.EntityICode, &entity.EntityCreationDate, &entity.EntityUpdateDate, &entity.EntityName, &entity.EntityDescription,
			&entity.EntityIsInteroperable, &entity.EntityInteroperabilityCode, &entity.EntityResponseTime, &entity.EntitySector)
		entities = append(entities, entity.ToDTO())
	}

	// Manejo de errores en la ejecución de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// GetEntitiesByMoment recupera registros de la entidad que están asociados a un
// "moment" (momento) específico. Esta función realiza un JOIN con la tabla
// RelEntityMoment para obtener los registros relacionados.
//
// Parámetros:
//   - moment: código o identificador del momento a filtrar.
//   - connData: datos de conexión actuales.
//   - clientConfig: configuración del cliente de base de datos.
//   - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
//   - un slice de EntityDTO con los registros filtrados y,
//     en caso de error, el error correspondiente.
//
// Nota: Se asume que las variables RelEntityMomentDBScheme, RelEntityMomentDBName y
// RelEntityMomentFieldDefinitions están definidas en otro lugar del código.
func GetEntitiesByMoment(moment string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entPath string = EntityDBScheme + "." + EntityDBName
	var relMomentPath string = RelEntityMomentDBScheme + "." + RelEntityMomentDBName

	// Establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar de la entidad y de la relación de "moment".
	var entFieldsSlice []string = []string{"EntityId", "EntityICode", "EntityName", "EntitySector"}
	var entFieldsAliasSlice []string = []string{}
	var relMomentFieldsSlice []string = []string{"RelEntityMomentCode"}
	var relMomentFieldsAliasSlice []string = []string{}

	// Genera el string de campos para ambas tablas.
	var entFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entFieldsSlice, entFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityDBScheme, EntityFieldDefinitions, true)
	var momentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relMomentFieldsSlice, relMomentFieldsAliasSlice, RelEntityMomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Construye la consulta SQL realizando un RIGHT JOIN entre la entidad y la tabla de momentos,
	// filtrando por el código del momento y ordenando los resultados.
	var query string = `SELECT ` + entFieldsStr + `, ` + momentFieldsStr +
		` FROM ` + entPath +
		` RIGHT JOIN ` + relMomentPath + ` ON (` + relMomentPath + `.` + RelEntityMomentFieldDefinitions["RelEntityMomentEntity"].DBName + ` = ` + entPath + `.` + EntityFieldDefinitions["EntityId"].DBName + `)` +
		` WHERE ` + relMomentPath + `.` + RelEntityMomentFieldDefinitions["RelEntityMomentCode"].DBName + ` = '` + moment + `' ORDER BY ` + entPath + `.` + EntityFieldDefinitions["EntityName"].DBName + ` ASC` +
		` ORDER BY ` + entPath + `.` + EntityFieldDefinitions["EntityName"].DBName + ` ASC `

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var entities []EntityDTO
	for persistenceCtrl.Next() {
		var entity = EntityDTO{}
		// Mapea cada fila de la consulta en la estructura EntityDTO.
		persistenceCtrl.ScanRow(&entity.EntityId, &entity.EntityICode, &entity.EntityName, &entity.EntitySector, &entity.EntityMoment)
		entities = append(entities, entity)
	}

	// Manejo de errores durante la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// SetEntityDefaults asigna valores por defecto a los campos de la entidad según
// la acción que se esté realizando (inserción o actualización).
//
// Parámetros:
//   - entity: puntero a la estructura EntityDTO que se modificará.
//   - action: acción que se va a realizar (por ejemplo, common_dao.SQL_INSERT o common_dao.SQL_UPDATE).
func SetEntityDefaults(entity *EntityDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción, se establecen las fechas de creación y actualización al momento actual,
		// y se genera un código único para la entidad.
		entity.EntityCreationDate = time.Now()
		entity.EntityUpdateDate = time.Now()
		entity.EntityICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para actualización, solo se actualiza la fecha de modificación.
		entity.EntityUpdateDate = time.Now()
	}
}

// ToDTO convierte una instancia de EntityPgDB (estructura de la BD) en un EntityDTO.
// Este método mapea los campos válidos de la base de datos al DTO correspondiente.
func (obj *EntityPgDB) ToDTO() EntityDTO {
	var dto EntityDTO

	if obj.EntityId.Valid {
		dto.EntityId = uint64(obj.EntityId.Int64)
	}

	if obj.EntityICode.Valid {
		dto.EntityICode = obj.EntityICode.String
	}

	if obj.EntityCreationDate.Valid {
		dto.EntityCreationDate = obj.EntityCreationDate.Time
	}

	if obj.EntityUpdateDate.Valid {
		dto.EntityUpdateDate = obj.EntityUpdateDate.Time
	}

	if obj.EntityName.Valid {
		dto.EntityName = obj.EntityName.String
	}

	if obj.EntityDescription.Valid {
		dto.EntityDescription = obj.EntityDescription.String
	}

	if obj.EntityIsInteroperable.Valid {
		dto.EntityIsInteroperable = obj.EntityIsInteroperable.String
	}

	if obj.EntityResponseTime.Valid {
		dto.EntityResponseTime = uint64(obj.EntityResponseTime.Int64)
	}

	if obj.EntitySector.Valid {
		dto.EntitySector = obj.EntitySector.String
	}

	return dto
}
