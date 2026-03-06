// Package salvia_daos proporciona las funcionalidades para interactuar con la
// entidad RelEntityMoment en la base de datos, incluyendo la conversión entre
// los modelos de la base de datos y los DTOs, y la ejecución de consultas SQL.
package salvia_daos

import (
	"context"
	"database/sql"
	"fmt"

	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
)

var (
	// RelEntityMomentRelEntityMomentCode es el código identificador de la entidad RelEntityMoment.
	RelEntityMomentRelEntityMomentCode string = "RelEntityMoment"
	// RelEntityMomentJSONName es el nombre que se usará en el JSON para representar un RelEntityMoment.
	RelEntityMomentJSONName string = "relMoment"
	// RelEntityMomentDBName es el nombre de la tabla en la base de datos para la entidad RelEntityMoment.
	RelEntityMomentDBName string = "rel_entity_moment"
	// RelEntityMomentDBScheme es el esquema de la base de datos donde se encuentra la tabla.
	RelEntityMomentDBScheme string = "salvia"

	// RelEntityMomentFieldDefinitions define las propiedades y validaciones de los campos de RelEntityMoment.
	// Cada campo se asocia a una definición que especifica su nombre, el nombre en la base de datos,
	// tipo de dato en el modelo, tamaño mínimo y máximo, y si es requerido o no.
	RelEntityMomentFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelEntityMomentId":     {Name: "RelEntityMomentId", DBName: "rel_entity_moment_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelEntityMomentCode":   {Name: "RelEntityMomentCode", DBName: "rel_entity_moment_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"RelEntityMomentEntity": {Name: "RelEntityMomentEntity", DBName: "rel_entity_moment_entity", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
	}
)

// RelEntityMomentDTO representa el Data Transfer Object para la entidad RelEntityMoment.
// Este DTO se utiliza para transferir datos entre las capas de la aplicación sin exponer
// detalles de implementación de la base de datos.
type RelEntityMomentDTO struct {
	RelEntityMomentId     uint64    `json:"-"`
	RelEntityMomentCode   string    `json:"code"`
	RelEntityMomentEntity EntityDTO `json:"-"`
}

// RelEntityMomentPgDB representa la estructura utilizada para mapear los resultados
// de las consultas SQL desde una base de datos PostgreSQL a un objeto Go.
type RelEntityMomentPgDB struct {
	RelEntityMomentId     sql.NullInt64
	RelEntityMomentCode   sql.NullString
	RelEntityMomentEntity sql.NullInt64
}

// GetRelEntityMoment obtiene un único registro de RelEntityMoment basado en los criterios
// especificados en 'by'. Se encarga de establecer la conexión, ejecutar la consulta y
// mapear los resultados al DTO proporcionado.
//
// Parámetros:
//   - by: estructura que define los atributos y operadores para la cláusula WHERE.
//   - relEntityMoment: puntero al DTO donde se almacenará el resultado.
//   - connData: datos de conexión que se pueden actualizar en la operación.
//   - clientConfig: configuración específica del cliente para la base de datos.
//   - serverConfig: configuración específica del servidor para la base de datos.
//
// Retorna:
//   - Un error en caso de que falle la operación.
func GetRelEntityMoment(by common_controllers.By, relEntityMoment *RelEntityMomentDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia para gestionar la conexión y consultas.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construye el path completo de la tabla en la base de datos.
	var relEntityPath string = RelEntityMomentDBScheme + "." + RelEntityMomentDBName

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var relEntityFieldsSlice []string = []string{"RelEntityMomentId", "RelEntityMomentCode", "RelEntityMomentEntity"}
	var relEntityFieldsAliasSlice []string = []string{}

	// Genera la parte de la consulta SQL correspondiente a la selección de campos.
	var relEntityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relEntityFieldsSlice, relEntityFieldsAliasSlice, RelEntityMomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Genera la consulta SQL completa, incluyendo la cláusula WHERE basada en 'by'.
	var query string = `SELECT ` + relEntityFieldsStr +
		` FROM ` + relEntityPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelEntityMomentDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Imprime la consulta para propósitos de depuración.
	fmt.Printf(query, by.AttrsValue...)

	// Se crea una variable para almacenar temporalmente los datos de la entidad relacionada.
	var entity EntityPgDB = EntityPgDB{}
	// Ejecuta la consulta y obtiene una única fila.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Escanea los resultados de la consulta en las variables correspondientes del DTO.
	persistenceCtrl.Scan(&relEntityMoment.RelEntityMomentId, &relEntityMoment.RelEntityMomentCode, &entity.EntityId)

	// Convierte la entidad obtenida de la base de datos a su representación DTO.
	relEntityMoment.RelEntityMomentEntity = entity.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelEntityMoments obtiene una lista de registros de RelEntityMoment que cumplen
// con los criterios especificados en 'by'. Realiza la conexión a la base de datos,
// ejecuta la consulta SQL y mapea cada resultado a un DTO.
//
// Parámetros:
//   - by: criterios para la cláusula WHERE de la consulta.
//   - connData: datos de conexión a la base de datos.
//   - clientConfig: configuración del cliente para la base de datos.
//   - serverConfig: configuración del servidor para la base de datos.
//
// Retorna:
//   - Un slice de RelEntityMomentDTO con los registros encontrados.
//   - Un error en caso de que falle la operación.
func GetRelEntityMoments(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelEntityMomentDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var relEntityPath string = RelEntityMomentDBScheme + "." + RelEntityMomentDBName

	var relEntity = RelEntityMomentDTO{}
	var relEntitys []RelEntityMomentDTO

	// Establece la conexión con la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var relEntityFieldsSlice []string = []string{"RelEntityMomentId", "RelEntityMomentCode", "RelEntityMomentEntity"}
	var relEntityFieldsAliasSlice []string = []string{}

	// Construye la parte de la consulta SQL para seleccionar campos.
	var relEntityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relEntityFieldsSlice, relEntityFieldsAliasSlice, RelEntityMomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Construye la consulta SQL completa con la cláusula WHERE.
	var query string = `SELECT ` + relEntityFieldsStr +
		` FROM ` + relEntityPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelEntityMomentDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Imprime la consulta generada para depuración.
	fmt.Printf(query, by.AttrsValue...)
	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre los resultados y mapea cada fila a un DTO.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&relEntity.RelEntityMomentId, &relEntity.RelEntityMomentCode, &relEntity.RelEntityMomentEntity)
		relEntitys = append(relEntitys, relEntity)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un 'by' para todos los DAOS.
	return relEntitys, nil
}

// GetAllEntityMoments obtiene todos los registros de RelEntityMoment sin aplicar filtros en la consulta.
// Se establece la conexión, se ejecuta la consulta y se mapean los resultados a un slice de DTOs.
//
// Parámetros:
//   - connData: datos de conexión a la base de datos.
//   - clientConfig: configuración del cliente para la base de datos.
//   - serverConfig: configuración del servidor para la base de datos.
//
// Retorna:

// - Un slice de RelEntityMomentDTO con todos los registros.
// - Un error en caso de que ocurra algún fallo.
func GetAllEntityMoments(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelEntityMomentDTO, error) {
	var relEntity RelEntityMomentDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var relEntityFieldsSlice []string = []string{"RelEntityMomentId", "RelEntityMomentCode", "RelEntityMomentEntity"}
	var relEntityFieldsAliasSlice []string = []string{}

	// Genera la consulta SQL completa sin cláusula WHERE.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, relEntityFieldsSlice, relEntityFieldsAliasSlice, RelEntityMomentDBName, []string{}, []string{}, []string{}, "", RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var relEntitys []RelEntityMomentDTO
	// Itera sobre cada fila obtenida y mapea a DTO.
	for persistenceCtrl.Next() {
		relEntity = RelEntityMomentDTO{}
		persistenceCtrl.ScanRow(&relEntity.RelEntityMomentId, &relEntity.RelEntityMomentCode, &relEntity.RelEntityMomentEntity)
		relEntitys = append(relEntitys, relEntity)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return relEntitys, nil
}

// SetRelEntityMomentDefaults establece valores por defecto en un RelEntityMomentDTO
// dependiendo de la acción a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
// Actualmente, la función está preparada para extender la configuración de valores predeterminados.
func SetRelEntityMomentDefaults(relEntity *RelEntityMomentDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Implementar valores por defecto para la inserción.
	case common_dao.SQL_UPDATE:
		// Implementar valores por defecto para la actualización.
	}
}

// PgDBToDTO convierte un objeto RelEntityMomentPgDB (representación de la base de datos PostgreSQL)
// a su correspondiente Data Transfer Object (RelEntityMomentDTO).
//
// Retorna:
//   - Un RelEntityMomentDTO con los valores mapeados desde la base de datos.
func (obj *RelEntityMomentPgDB) PgDBToDTO() RelEntityMomentDTO {
	var dto RelEntityMomentDTO

	if obj.RelEntityMomentId.Valid {
		dto.RelEntityMomentId = uint64(obj.RelEntityMomentId.Int64)
	}

	if obj.RelEntityMomentCode.Valid {
		dto.RelEntityMomentCode = obj.RelEntityMomentCode.String
	}

	if obj.RelEntityMomentEntity.Valid {
		dto.RelEntityMomentEntity = EntityDTO{EntityId: uint64(obj.RelEntityMomentEntity.Int64)}
	}

	return dto
}
