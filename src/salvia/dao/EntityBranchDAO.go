// Package salvia_daos contiene los Data Access Objects (DAO) para manejar las operaciones
// de persistencia relacionadas con la entidad "EntityBranch" en la base de datos Salvia.
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

// Variables globales que definen nombres de entidad, nombres JSON, nombres de tablas y esquema en BD,
// además de la definición de los campos utilizados para validaciones y generación dinámica de consultas.
var (
	// EntityBranchName es el nombre de la entidad en el modelo.
	EntityBranchName string = "EntityBranch"
	// EntityBranchJSONName es el nombre de la entidad en formato JSON.
	EntityBranchJSONName string = "entityBranch"
	// EntityBranchDBName es el nombre de la tabla en la base de datos.
	EntityBranchDBName string = "entity_branch"
	// EntityBranchDBScheme es el esquema de la base de datos en el que se encuentra la tabla.
	EntityBranchDBScheme string = "salvia"

	EntityBranchUserMark string = ""

	// EntityBranchFieldDefinitions define los atributos del modelo, sus equivalentes en la BD,
	// tipo de dato, tamaño mínimo/máximo y si son requeridos. Se utiliza para validaciones y generación de SQL.
	EntityBranchFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"EntityBranchId":           {Name: "EntityBranchId", DBName: "entity_branch_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"EntityBranchICode":        {Name: "EntityBranchICode", DBName: "entity_branch_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"EntityBranchCreationDate": {Name: "EntityBranchCreationDate", DBName: "entity_branch_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"EntityBranchUpdateDate":   {Name: "EntityBranchUpdateDate", DBName: "entity_branch_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"EntityBranchName":         {Name: "EntityBranchName", DBName: "entity_branch_name", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 254, Required: true},
		"EntityBranchDescription":  {Name: "EntityBranchDescription", DBName: "entity_branch_description", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 254, Required: true},
		"EntityBranchAddress":      {Name: "EntityBranchAddress", DBName: "entity_branch_address", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 254, Required: true},
		"EntityBranchLatitude":     {Name: "EntityBranchLatitude", DBName: "entity_branch_latitude", Alias: "", ModelType: "float", MinSize: -90, MaxSize: 90, Required: false},
		"EntityBranchLongitude":    {Name: "EntityBranchLongitude", DBName: "entity_branch_longitude", Alias: "", ModelType: "float", MinSize: -180, MaxSize: 180, Required: false},
		"EntityBranchEntity":       {Name: "EntityBranchEntity", DBName: "entity_id", Alias: "", ModelType: "uint", Required: false},
		"EntityBranchTownCode":     {Name: "EntityBranchTownCode", DBName: "entity_branch_town_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"EntityBranchSource":       {Name: "EntityBranchSource", DBName: "entity_branch_source", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
	}
)

// EntityBranchDTO representa el objeto de transferencia de datos (DTO) para la entidad EntityBranch.
// Incluye tanto campos que se almacenan en la base de datos como campos adicionales utilizados en formularios.
type EntityBranchDTO struct {
	EntityBranchId           uint64    `json:"-"`
	EntityBranchICode        string    `json:"icode"`
	EntityBranchCreationDate time.Time `json:"creationDate"`
	EntityBranchUpdateDate   time.Time `json:"updateDate"`
	EntityBranchName         string    `json:"name"`
	EntityBranchDescription  string    `json:"description"`
	EntityBranchAddress      string    `json:"address"`
	EntityBranchLatitude     float64   `json:"latitude"`
	EntityBranchLongitude    float64   `json:"longitude"`
	EntityBranchTownCode     string    `json:"townCode"`
	EntityBranchSource       string    `json:"branchSource"`
	EntityBranchEntity       EntityDTO `json:"-"`
	EntityBranchSector       string    `json:"sector"`
	EntityBranchMoment       string    `json:"moment"`
	EntityBranchEntityICode  string    `json:"entityICode"`
	EntityBranchEntityName   string    `json:"entityName"`
}

// EntityBranchPgDB representa la estructura de la tabla "entity_branch" en la base de datos PostgreSQL.
// Se utiliza sql.Null* para manejar valores nulos provenientes de la BD.
type EntityBranchPgDB struct {
	EntityBranchId           sql.NullInt64
	EntityBranchICode        sql.NullString
	EntityBranchCreationDate sql.NullTime
	EntityBranchUpdateDate   sql.NullTime
	EntityBranchName         sql.NullString
	EntityBranchDescription  sql.NullString
	EntityBranchAddress      sql.NullString
	EntityBranchLatitude     sql.NullFloat64
	EntityBranchLongitude    sql.NullFloat64
	EntityBranchTownCode     sql.NullString
	EntityBranchEntity       sql.NullInt64
	EntityBranchSource       sql.NullString
}

func (eb EntityBranchDTO) MarshalJSON() ([]byte, error) {
	type Alias EntityBranchDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		EntityBranchCreationDate string `json:"creationDate"`
		EntityBranchUpdateDate   string `json:"updateDate"`
	}{
		Alias:                    (*Alias)(&eb),
		EntityBranchCreationDate: eb.EntityBranchCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		EntityBranchUpdateDate:   eb.EntityBranchUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (eb *EntityBranchDTO) UnmarshalJSON(data []byte) error {
	type Alias EntityBranchDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		EntityBranchCreationDate string `json:"creationDate"`
		EntityBranchUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(eb),
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
	eb.EntityBranchCreationDate = parse(aux.EntityBranchCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	eb.EntityBranchUpdateDate = parse(aux.EntityBranchUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetEntityBranch inserta un nuevo registro de EntityBranch en la base de datos.
// Recibe un puntero a EntityBranchDTO, un indicador de transacción, el módulo y la configuración de conexión.
// Devuelve la conexión actualizada y un error en caso de fallar.
func SetEntityBranch(entityBranch *EntityBranchDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Define los campos que se insertarán en la BD.
	var entityBranchFieldsSlice []string = []string{"EntityBranchICode", "EntityBranchCreationDate", "EntityBranchUpdateDate", "EntityBranchName", "EntityBranchDescription", "EntityBranchAddress",
		"EntityBranchLatitude", "EntityBranchLongitude", "EntityBranchTownCode", "EntityBranchEntity", "EntityBranchSource"}
	var entityBranchFieldsAliasSlice []string = []string{}

	// Genera la consulta SQL de inserción de forma dinámica.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, entityBranchFieldsSlice, entityBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{"EntityBranchId"}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, false)

	// Ejecuta la consulta pasando los valores del DTO.
	persistenceCtrl.QueryRow(context.Background(), query,
		entityBranch.EntityBranchICode, entityBranch.EntityBranchCreationDate, entityBranch.EntityBranchUpdateDate,
		entityBranch.EntityBranchName, entityBranch.EntityBranchDescription, entityBranch.EntityBranchAddress, entityBranch.EntityBranchLatitude,
		entityBranch.EntityBranchLongitude, entityBranch.EntityBranchTownCode, entityBranch.EntityBranchEntity.EntityId, entityBranch.EntityBranchSource)

	// Escanea el resultado para obtener el ID generado.
	persistenceCtrl.Scan(&entityBranch.EntityBranchId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateEntityBranch actualiza un registro existente de EntityBranch en la base de datos.
// Recibe el DTO con los datos a actualizar y retorna la conexión actualizada junto con un error si ocurre.
func UpdateEntityBranch(entityBranch *EntityBranchDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos que se actualizarán.
	var entityBranchFieldsSlice []string = []string{"EntityBranchUpdateDate", "EntityBranchName", "EntityBranchDescription", "EntityBranchAddress", "EntityBranchLatitude", "EntityBranchLongitude", "EntityBranchTownCode", "EntityBranchEntity"}
	var entityBranchFieldsAliasSlice []string = []string{}

	// Genera la consulta SQL de actualización de forma dinámica.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, entityBranchFieldsSlice, entityBranchFieldsAliasSlice, EntityBranchDBName, []string{"EntityBranchId"}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, false)

	// Ejecuta la consulta pasando los valores correspondientes.
	persistenceCtrl.Exec(context.Background(), query, entityBranch.EntityBranchId,
		entityBranch.EntityBranchUpdateDate, entityBranch.EntityBranchName, entityBranch.EntityBranchDescription, entityBranch.EntityBranchAddress,
		entityBranch.EntityBranchLatitude, entityBranch.EntityBranchLongitude, entityBranch.EntityBranchTownCode, entityBranch.EntityBranchEntity.EntityId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}
	// Verifica que se haya afectado al menos una fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetEntityBranch recupera un único registro de EntityBranch basado en un conjunto de condiciones.
// El parámetro 'by' especifica los atributos y valores para filtrar la consulta.
// Devuelve la conexión actualizada y el error (si ocurre).
func GetEntityBranch(by common_controllers.By, entityBranch *EntityBranchDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entityBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var entityBranchFieldsSlice []string = []string{"EntityBranchId", "EntityBranchICode", "EntityBranchCreationDate", "EntityBranchUpdateDate", "EntityBranchName", "EntityBranchDescription", "EntityBranchAddress",
		"EntityBranchLatitude", "EntityBranchLongitude", "EntityBranchTownCode", "EntityBranchEntity", "EntityBranchSource"}
	var entityBranchFieldsAliasSlice []string = []string{}

	// Genera la parte de la consulta que indica los campos a seleccionar.
	var entityBranchFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityBranchFieldsSlice, entityBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)

	// Genera la consulta SQL completa con cláusula WHERE basada en el parámetro 'by'.
	var query string = `SELECT ` + entityBranchFieldsStr +
		` FROM ` + entityBranchPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, EntityBranchDBName, by.AttrsName, []string{}, []string{}, by.Operator, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)

	fmt.Printf(query, by.AttrsValue...)
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se crea una instancia temporal para almacenar los datos de la entidad asociada.
	var entity EntityDTO = EntityDTO{}

	// Escanea el resultado de la consulta en los campos del DTO.
	persistenceCtrl.Scan(&entityBranch.EntityBranchId, &entityBranch.EntityBranchICode, &entityBranch.EntityBranchCreationDate, &entityBranch.EntityBranchUpdateDate,
		&entityBranch.EntityBranchName, &entityBranch.EntityBranchDescription, &entityBranch.EntityBranchAddress, &entityBranch.EntityBranchLatitude,
		&entityBranch.EntityBranchLongitude, &entityBranch.EntityBranchTownCode, &entity.EntityId, &entityBranch.EntityBranchSource)

	entityBranch.EntityBranchEntity = entity
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetEntityBranches recupera múltiples registros de EntityBranch basándose en condiciones definidas en 'by'.
// Devuelve la conexión actualizada, un slice de EntityBranchDTO y un error en caso de ocurrir.
func GetEntityBranches(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityBranchDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entityBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName

	var entityBranch = EntityBranchDTO{}
	var entityBranchs []EntityBranchDTO

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var entityBranchFieldsSlice []string = []string{"EntityBranchId", "EntityBranchICode", "EntityBranchCreationDate", "EntityBranchUpdateDate", "EntityBranchName", "EntityBranchDescription",
		"EntityBranchAddress", "EntityBranchLatitude", "EntityBranchLongitude", "EntityBranchSource"}
	var entityBranchFieldsAliasSlice []string = []string{}

	// Genera la parte de la consulta SQL para la selección de campos.
	var entityBranchFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityBranchFieldsSlice, entityBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)

	// Genera la consulta SQL completa con cláusula WHERE.
	var query string = `SELECT ` + entityBranchFieldsStr +
		` FROM ` + entityBranchPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, EntityBranchDBName, by.AttrsName, []string{}, []string{}, by.Operator, EntityBranchDBScheme, EntityBranchFieldDefinitions, true) +
		`ORDER BY ` + entityBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchSource"].DBName + ` DESC, ` +
		entityBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchName"].DBName + ` ASC`

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre los resultados y construye el slice de DTOs.
	for persistenceCtrl.Next() {

		persistenceCtrl.ScanRow(&entityBranch.EntityBranchId, &entityBranch.EntityBranchICode, &entityBranch.EntityBranchCreationDate, &entityBranch.EntityBranchUpdateDate,
			&entityBranch.EntityBranchName, &entityBranch.EntityBranchDescription, &entityBranch.EntityBranchAddress, &entityBranch.EntityBranchLatitude, &entityBranch.EntityBranchLongitude,
			&entityBranch.EntityBranchSource)

		if entityBranch.EntityBranchSource == "u" {
			entityBranch.EntityBranchName = EntityBranchUserMark + entityBranch.EntityBranchName
		}
		entityBranchs = append(entityBranchs, entityBranch)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un 'by' para todos los DAOS.
	return entityBranchs, nil
}

// GetAllEntityBranch recupera todos los registros de EntityBranch, realizando un JOIN con la entidad asociada
// para incluir información adicional (por ejemplo, sector). Se ordenan los resultados por nombre.
func GetAllEntityBranch(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityBranchDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var entityBranch EntityBranchDTO
	var entBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName
	// 'entityPath' se utiliza para acceder a la tabla de entidades asociadas.
	var entityPath string = EntityDBScheme + "." + EntityDBName

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar para la tabla entity_branch.
	var entBranchFieldsSlice []string = []string{"EntityBranchId", "EntityBranchICode", "EntityBranchName", "EntityBranchAddress", "EntityBranchLatitude", "EntityBranchLongitude", "EntityBranchSector", "EntityBranchSource"}
	var entBranchFieldsAliasSlice []string = []string{}

	// Define los campos a seleccionar de la tabla de entidades.
	var entityFieldsSlice []string = []string{"EntitySector"}
	var entityFieldsAliasSlice []string = []string{}

	// Genera las partes de la consulta para cada tabla.
	var entBranchFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entBranchFieldsSlice, entBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)
	var entityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityFieldsSlice, entityFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityDBScheme, EntityFieldDefinitions, true)

	// Construye la consulta SQL con un RIGHT JOIN para relacionar ambas tablas y ordena los resultados.
	var query string = `SELECT ` + entBranchFieldsStr + `, ` + entityFieldsStr +
		` FROM ` + entBranchPath +
		` RIGHT JOIN ` + entityPath + ` ON (` + entityPath + `.` + EntityFieldDefinitions["EntityId"].DBName + ` = ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchEntity"].DBName + `)` +
		` WHERE TRUE ORDER BY ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchSource"].DBName + ` DESC, ` +
		entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchName"].DBName + ` ASC`

	persistenceCtrl.Query(context.Background(), query)
	var entityBranchs []EntityBranchDTO
	for persistenceCtrl.Next() {
		entityBranch = EntityBranchDTO{}
		persistenceCtrl.ScanRow(&entityBranch.EntityBranchId, &entityBranch.EntityBranchICode, &entityBranch.EntityBranchName,
			&entityBranch.EntityBranchAddress, &entityBranch.EntityBranchLatitude, &entityBranch.EntityBranchLongitude, &entityBranch.EntityBranchSector, &entityBranch.EntityBranchSource)

		if entityBranch.EntityBranchSource == "u" {
			entityBranch.EntityBranchName = EntityBranchUserMark + entityBranch.EntityBranchName
		}
		entityBranchs = append(entityBranchs, entityBranch)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entityBranchs, nil
}

// GetEntityBranchesBySector recupera un registro de EntityBranch filtrado por sector.
// Realiza un JOIN con la tabla de entidades para obtener información adicional.
// Parámetros:
//   - sector: Sector por el cual filtrar.
//   - entBranch: Puntero al DTO donde se almacenarán los datos recuperados.
func GetEntityBranchesBySector(sector string, entBranch *EntityBranchDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName
	var entityPath string = EntityDBScheme + "." + EntityDBName

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos a seleccionar de la tabla entity_branch.
	var entBranchFieldsSlice []string = []string{"EntityBranchId", "EntityBranchICode", "EntityBranchName", "EntityBranchAddress", "EntityBranchLatitude", "EntityBranchLongitude", "EntityBranchSector", "EntityBranchSource"}
	var entBranchFieldsAliasSlice []string = []string{}

	// Define los campos a seleccionar de la tabla de entidades.
	var entityFieldsSlice []string = []string{"EntitySector"}
	var entityFieldsAliasSlice []string = []string{}

	// Genera las partes de la consulta SQL.
	var entBranchFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entBranchFieldsSlice, entBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)
	var entityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityFieldsSlice, entityFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)

	// Construye la consulta SQL con un JOIN y filtra por el sector.
	var query string = `SELECT ` + entBranchFieldsStr + `, ` + entityFieldsStr +
		` FROM ` + entBranchPath +
		` RIGHT JOIN ` + entityPath + ` ON (` + entityPath + `.` + EntityFieldDefinitions["EntityId"].DBName + ` = ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchEntity"].DBName + `)` +
		` WHERE ` + entityPath + `.` + EntityFieldDefinitions["EntitySector"].DBName + ` = '` + sector + `' ORDER BY ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchSource"].DBName + ` DESC, ` +
		entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchName"].DBName + ` ASC`

	// Ejecuta la consulta.
	persistenceCtrl.QueryRow(context.Background(), query, sector)

	// Escanea el resultado en el DTO.
	persistenceCtrl.Scan(&entBranch.EntityBranchId, &entBranch.EntityBranchICode, &entBranch.EntityBranchName, &entBranch.EntityBranchAddress,
		&entBranch.EntityBranchLatitude, &entBranch.EntityBranchLongitude, &entBranch.EntityBranchSector, &entBranch.EntityBranchSource)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetEntityBranchesByTownCodeWithMoments recupera registros de EntityBranch filtrados por townCode,
// incluyendo datos de la entidad asociada y momentos relacionados.
// Realiza JOINs con las tablas Entity, RelEntityMoment y Town.
// Parámetros:
//   - townCode: Código del municipio para filtrar.
func GetEntityBranchesByTownCodeWithMoments(townCode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityBranchDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName
	var entityPath string = EntityDBScheme + "." + EntityDBName
	var relMomentPath string = RelEntityMomentDBScheme + "." + RelEntityMomentDBName
	var townPath string = security_daos.TownDBScheme + "." + security_daos.TownDBName

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar de entity_branch.
	var entBranchFieldsSlice []string = []string{"EntityBranchId", "EntityBranchICode", "EntityBranchName", "EntityBranchAddress", "EntityBranchTownCode", "EntityBranchSource"}
	var entBranchFieldsAliasSlice []string = []string{}

	// Define los campos a seleccionar de la tabla de entidades.
	var entityFieldsSlice []string = []string{"EntityICode", "EntityName", "EntitySector"}
	var entityFieldsAliasSlice []string = []string{}

	// Define los campos a seleccionar de la tabla de momentos.
	var relMomentFieldsSlice []string = []string{"RelEntityMomentCode"}
	var relMomentFieldsAliasSlice []string = []string{}

	// Genera las partes de la consulta SQL para cada tabla.
	var entBranchFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entBranchFieldsSlice, entBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)
	var entityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityFieldsSlice, entityFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityDBScheme, EntityFieldDefinitions, true)
	var momentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relMomentFieldsSlice, relMomentFieldsAliasSlice, RelEntityMomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Construye la consulta SQL con múltiples JOINs y filtra por el townCode.
	var query string = `SELECT ` + entBranchFieldsStr + `, ` + entityFieldsStr + `, ` + momentFieldsStr +
		` FROM ` + entBranchPath +
		` RIGHT JOIN ` + entityPath + ` ON (` + entityPath + `.` + EntityFieldDefinitions["EntityId"].DBName + ` = ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchEntity"].DBName + `)` +
		` RIGHT JOIN ` + relMomentPath + ` ON (` + relMomentPath + `.` + RelEntityMomentFieldDefinitions["RelEntityMomentEntity"].DBName + ` = ` + entityPath + `.` + EntityFieldDefinitions["EntityId"].DBName + `)` +
		` LEFT JOIN ` + townPath + ` ON (` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + ` = ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchTownCode"].DBName + `)` +
		` WHERE ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + ` = $1` +
		` ORDER BY ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchSource"].DBName + ` DESC, ` +
		entityPath + `.` + EntityFieldDefinitions["EntityName"].DBName + ` ASC, ` +
		entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchName"].DBName + ` ASC`

	// Ejecuta la consulta pasando el townCode como parámetro.
	persistenceCtrl.Query(context.Background(), query, townCode)
	var entities []EntityBranchDTO
	for persistenceCtrl.Next() {
		var branch = EntityBranchDTO{}
		persistenceCtrl.ScanRow(&branch.EntityBranchId, &branch.EntityBranchICode, &branch.EntityBranchName, &branch.EntityBranchAddress, &branch.EntityBranchTownCode, &branch.EntityBranchSource,
			&branch.EntityBranchEntityICode, &branch.EntityBranchEntityName, &branch.EntityBranchSector, &branch.EntityBranchMoment)

		if branch.EntityBranchSource == "u" {
			branch.EntityBranchName = EntityBranchUserMark + branch.EntityBranchName
		}
		entities = append(entities, branch)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// GetEntityBranchesByTownCodeAndEntityBranchIcodeWithMoments recupera registros de EntityBranch filtrados
// tanto por townCode como por entityBranchICode, incluyendo datos de la entidad y momentos relacionados.
// Realiza JOINs con las tablas Entity, RelEntityMoment y Town.
func GetEntityBranchesByTownCodeAndEntityBranchIcodeWithMoments(townCode string, entityBranchICode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]EntityBranchDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName
	var entityPath string = EntityDBScheme + "." + EntityDBName
	var relMomentPath string = RelEntityMomentDBScheme + "." + RelEntityMomentDBName
	var townPath string = security_daos.TownDBScheme + "." + security_daos.TownDBName

	// Se obtiene la conexión a la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar de entity_branch.
	var entBranchFieldsSlice []string = []string{"EntityBranchId", "EntityBranchICode", "EntityBranchName", "EntityBranchAddress", "EntityBranchTownCode", "EntityBranchSource"}
	var entBranchFieldsAliasSlice []string = []string{}

	// Define los campos a seleccionar de la tabla de entidades.
	var entityFieldsSlice []string = []string{"EntityICode", "EntityName", "EntitySector"}
	var entityFieldsAliasSlice []string = []string{}

	// Define los campos a seleccionar de la tabla de momentos.
	var relMomentFieldsSlice []string = []string{"RelEntityMomentCode"}
	var relMomentFieldsAliasSlice []string = []string{}

	// Genera las partes de la consulta SQL para cada tabla.
	var entBranchFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entBranchFieldsSlice, entBranchFieldsAliasSlice, EntityBranchDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityBranchDBScheme, EntityBranchFieldDefinitions, true)
	var entityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entityFieldsSlice, entityFieldsAliasSlice, EntityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, EntityDBScheme, EntityFieldDefinitions, true)
	var momentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relMomentFieldsSlice, relMomentFieldsAliasSlice, RelEntityMomentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelEntityMomentDBScheme, RelEntityMomentFieldDefinitions, true)

	// Construye la consulta SQL con múltiples JOINs y filtra por townCode y entityBranchICode.
	var query string = `SELECT ` + entBranchFieldsStr + `, ` + entityFieldsStr + `, ` + momentFieldsStr +
		` FROM ` + entBranchPath +
		` RIGHT JOIN ` + entityPath + ` ON (` + entityPath + `.` + EntityFieldDefinitions["EntityId"].DBName + ` = ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchEntity"].DBName + `)` +
		` RIGHT JOIN ` + relMomentPath + ` ON (` + relMomentPath + `.` + RelEntityMomentFieldDefinitions["RelEntityMomentEntity"].DBName + ` = ` + entityPath + `.` + EntityFieldDefinitions["EntityId"].DBName + `)` +
		` LEFT JOIN ` + townPath + ` ON (` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + ` = ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchTownCode"].DBName + `)` +
		` WHERE ` + townPath + `.` + security_daos.TownFieldDefinitions["TownCode"].DBName + ` = $1 AND ` +
		entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchICode"].DBName + ` = $2 ` +
		` ORDER BY ` + entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchSource"].DBName + ` DESC, ` +
		entityPath + `.` + EntityFieldDefinitions["EntityName"].DBName + ` ASC, ` +
		entBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchName"].DBName + ` ASC`

	// Ejecuta la consulta pasando los parámetros correspondientes.
	persistenceCtrl.Query(context.Background(), query, townCode, entityBranchICode)
	var entities []EntityBranchDTO
	for persistenceCtrl.Next() {
		var branch = EntityBranchDTO{}
		persistenceCtrl.ScanRow(&branch.EntityBranchId, &branch.EntityBranchICode, &branch.EntityBranchName, &branch.EntityBranchAddress, &branch.EntityBranchTownCode, &branch.EntityBranchSource,
			&branch.EntityBranchEntityICode, &branch.EntityBranchEntityName, &branch.EntityBranchSector, &branch.EntityBranchMoment)

		if branch.EntityBranchSource == "u" {
			branch.EntityBranchName = EntityBranchUserMark + branch.EntityBranchName
		}
		entities = append(entities, branch)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entities, nil
}

// SetEntityBranchDefaults establece valores por defecto en un objeto EntityBranchDTO según la acción a realizar.
// Por ejemplo, en una inserción se inicializa la fecha de creación, la fecha de actualización y se asigna un UUID.
func SetEntityBranchDefaults(entityBranch *EntityBranchDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		entityBranch.EntityBranchCreationDate = time.Now()
		entityBranch.EntityBranchUpdateDate = time.Now()
		entityBranch.EntityBranchICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		entityBranch.EntityBranchUpdateDate = time.Now()
	}
}

// ToDTO convierte una instancia de EntityBranchPgDB en un EntityBranchDTO,
// transformando los valores nulos en los tipos correspondientes.
func (obj *EntityBranchPgDB) ToDTO() EntityBranchDTO {
	var dto EntityBranchDTO

	if obj.EntityBranchId.Valid {
		dto.EntityBranchId = uint64(obj.EntityBranchId.Int64)
	}

	if obj.EntityBranchICode.Valid {
		dto.EntityBranchICode = obj.EntityBranchICode.String
	}

	if obj.EntityBranchCreationDate.Valid {
		dto.EntityBranchCreationDate = obj.EntityBranchCreationDate.Time
	}

	if obj.EntityBranchUpdateDate.Valid {
		dto.EntityBranchUpdateDate = obj.EntityBranchUpdateDate.Time
	}

	if obj.EntityBranchName.Valid {
		dto.EntityBranchName = obj.EntityBranchName.String
	}

	if obj.EntityBranchDescription.Valid {
		dto.EntityBranchDescription = obj.EntityBranchDescription.String
	}

	if obj.EntityBranchAddress.Valid {
		dto.EntityBranchAddress = obj.EntityBranchAddress.String
	}

	if obj.EntityBranchLatitude.Valid {
		dto.EntityBranchLatitude = obj.EntityBranchLatitude.Float64
	}

	if obj.EntityBranchLongitude.Valid {
		dto.EntityBranchLongitude = obj.EntityBranchLongitude.Float64
	}

	if obj.EntityBranchTownCode.Valid {
		dto.EntityBranchTownCode = obj.EntityBranchTownCode.String
	}

	if obj.EntityBranchSource.Valid {
		dto.EntityBranchSource = obj.EntityBranchSource.String
	}

	if obj.EntityBranchEntity.Valid {
		dto.EntityBranchEntity = EntityDTO{EntityId: uint64(obj.EntityBranchEntity.Int64)}
	}

	return dto
}
