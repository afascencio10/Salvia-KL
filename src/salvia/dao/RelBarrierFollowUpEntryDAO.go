package salvia_daos

import (
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
	// RelBarrierFollowUpEntryEntityName es el nombre lógico de la entidad de relación entre Barrera y Entrada de Seguimiento.
	RelBarrierFollowUpEntryEntityName string = "RelBarrierFollowUpEntry"
	// RelBarrierFollowUpEntryJSONName es el nombre JSON para la entidad RelBarrierFollowUpEntry.
	RelBarrierFollowUpEntryJSONName string = "relBarrierFollowUpEntry"
	// RelBarrierFollowUpEntryDBName es el nombre de la tabla en la base de datos para RelBarrierFollowUpEntry.
	RelBarrierFollowUpEntryDBName string = "rel_barrier_follow_up_entry"
	// RelBarrierFollowUpEntryDBScheme es el esquema de la base de datos para la tabla RelBarrierFollowUpEntry.
	RelBarrierFollowUpEntryDBScheme string = "salvia"

	// RelBarrierFollowUpEntryFieldDefinitions define las propiedades de los campos de la entidad RelBarrierFollowUpEntry,
	// incluyendo su nombre, nombre en la DB, alias, tipo de modelo y si es requerido.
	RelBarrierFollowUpEntryFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelBarrierFollowUpEntryId":            {Name: "RelBarrierFollowUpEntryId", DBName: "rel_barrier_follow_up_entry_id", Alias: "", ModelType: "uint", Required: false},
		"RelBarrierFollowUpEntryCreationDate":  {Name: "RelBarrierFollowUpEntryCreationDate", DBName: "rel_barrier_follow_up_entry_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"RelBarrierFollowUpEntryBarrier":       {Name: "RelBarrierFollowUpEntryBarrier", DBName: "rel_barrier_follow_up_entry_barrier", Alias: "", ModelType: "uint", Required: true},
		"RelBarrierFollowUpEntryFollowUpEntry": {Name: "RelBarrierFollowUpEntryFollowUpEntry", DBName: "rel_barrier_follow_up_entry_follow_up_entry", Alias: "", ModelType: "uint", Required: true},
	}
)

// RelBarrierFollowUpEntryDTO representa el Data Transfer Object para la relación entre una Barrera y una Entrada de Seguimiento.
// Se utiliza para transferir datos entre las capas de la aplicación y la presentación.
type RelBarrierFollowUpEntryDTO struct {
	RelBarrierFollowUpEntryId            uint64 `json:"-"` // ID de la relación (no serializado a JSON).
	RelBarrierFollowUpEntryCreationDate  time.Time
	RelBarrierFollowUpEntryBarrier       BarrierDTO       `json:"barrier"`       // Objeto DTO de la barrera asociada.
	RelBarrierFollowUpEntryFollowUpEntry FollowUpEntryDTO `json:"followUpEntry"` // Objeto DTO de la entrada de seguimiento asociada.
}

// RelBarrierFollowUpEntryPgDB representa la estructura de la relación entre Barrera y Entrada de Seguimiento tal como se almacena en PostgreSQL.
// Utiliza tipos sql.Null para manejar valores nulos de la base de datos de forma segura.
type RelBarrierFollowUpEntryPgDB struct {
	RelBarrierFollowUpEntryId            sql.NullInt64 // ID de la relación.
	RelBarrierFollowUpEntryCreationDate  sql.NullTime
	RelBarrierFollowUpEntryBarrier       sql.NullInt64 // ID de la barrera asociada.
	RelBarrierFollowUpEntryFollowUpEntry sql.NullInt64 // ID de la entrada de seguimiento asociada.
}

func (rcwvc RelBarrierFollowUpEntryDTO) MarshalJSON() ([]byte, error) {
	type Alias RelBarrierFollowUpEntryDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		RelBarrierFollowUpEntryCreationDate string `json:"creationDate"`
	}{
		Alias:                               (*Alias)(&rcwvc),
		RelBarrierFollowUpEntryCreationDate: rcwvc.RelBarrierFollowUpEntryCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (rcwvc *RelBarrierFollowUpEntryDTO) UnmarshalJSON(data []byte) error {
	type Alias RelBarrierFollowUpEntryDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		RelBarrierFollowUpEntryCreationDate string `json:"creationDate"`
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
	rcwvc.RelBarrierFollowUpEntryCreationDate = parse(aux.RelBarrierFollowUpEntryCreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetRelBarrierFollowUpEntry inserta una nueva relación entre una barrera y una entrada de seguimiento en la base de datos.
// Parámetros:
//   - rel: puntero al DTO que contiene los datos de la relación a insertar.
//   - connData: datos de la conexión actual a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de producirse algún fallo.
func SetRelBarrierFollowUpEntry(rel *RelBarrierFollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelBarrierFollowUpEntryBarrier", "RelBarrierFollowUpEntryFollowUpEntry"}
	var relFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, relFieldsSlice, relFieldsAliasSlice, RelBarrierFollowUpEntryDBName, []string{}, []string{}, []string{"RelBarrierFollowUpEntryId"}, common_dao.SQL_AND, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query, rel.RelBarrierFollowUpEntryBarrier.BarrierId, rel.RelBarrierFollowUpEntryFollowUpEntry.FollowUpEntryId)

	persistenceCtrl.Scan(&rel.RelBarrierFollowUpEntryId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// RemoveRelBarrierFollowUpEntry elimina una relación específica entre una barrera y una entrada de seguimiento en la base de datos.
// Parámetros:
//   - rel: puntero al DTO que contiene el ID de la relación a eliminar.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func RemoveRelBarrierFollowUpEntry(rel *RelBarrierFollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{}
	var relFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, relFieldsSlice, relFieldsAliasSlice, RelBarrierFollowUpEntryDBName, []string{}, []string{"RelBarrierFollowUpEntryId"}, []string{}, common_dao.SQL_AND, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, rel.RelBarrierFollowUpEntryId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	if persistenceCtrl.RowsAffected == 0 {
		return errors.New("common_db_no_rows_affected")
	}

	return nil
}

// RemoveRelBarrierFollowUpEntries elimina una o más relaciones entre barreras y entradas de seguimiento basándose en criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (nombres de atributos, alias, operadores y valores).
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func RemoveRelBarrierFollowUpEntries(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, by.AttrsName, by.AttrsAliasName, RelBarrierFollowUpEntryDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, by.AttrsValue...)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelBarrierFollowUpEntry obtiene una relación específica entre una barrera y una entrada de seguimiento según los criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda (atributos, alias, operador y valores).
//   - rel: puntero al DTO donde se almacenarán los datos obtenidos.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna un error en caso de fallo.
func GetRelBarrierFollowUpEntry(by common_controllers.By, rel *RelBarrierFollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var relPath string = RelBarrierFollowUpEntryDBScheme + "." + RelBarrierFollowUpEntryDBName
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var followUpEntryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelBarrierFollowUpEntryId", "RelBarrierFollowUpEntryBarrier", "RelBarrierFollowUpEntryFollowUpEntry"}
	var relFieldsAliasSlice []string = []string{}

	var barrierFieldsSlice []string = []string{"BarrierICode", "BarrierName", "BarrierDescription", "BarrierSectorBarrier"}
	var barrierFieldsAliasSlice []string = []string{}

	var followUpEntryFieldsSlice []string = []string{"FollowUpEntryICode", "FollowUpEntryCreationDate", "FollowUpEntryUpdateDate", "FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone", "FollowUpEntrySector", "FollowUpEntryPersonVisitedEntity", "FollowUpEntryReceivedAttention",
		"FollowUpEntryComments", "FollowUpEntryCaseDocumentsPrepared", "FollowUpEntryStatus", "FollowUpEntryFollowUp"}
	var followUpEntryFieldsAliasSlice []string = []string{}

	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelBarrierFollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, true)
	var barrierFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, true)
	var followUpEntryFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, followUpEntryFieldsSlice, followUpEntryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	var query string = `SELECT ` + relFieldsStr + `, ` + barrierFieldsStr + `, ` + followUpEntryFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + barrierPath + ` ON (` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryBarrier"].DBName + `)` +
		` LEFT JOIN ` + followUpEntryPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryFollowUpEntry"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelBarrierFollowUpEntryDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var barrier BarrierPgDB = BarrierPgDB{}
	var followUpEntry FollowUpEntryPgDB = FollowUpEntryPgDB{}

	persistenceCtrl.Scan(&rel.RelBarrierFollowUpEntryId, &barrier.BarrierId, &followUpEntry.FollowUpEntryId,
		&barrier.BarrierICode, &barrier.BarrierName, &barrier.BarrierDescription, &barrier.BarrierSectorBarrier, &followUpEntry.FollowUpEntryICode, &followUpEntry.FollowUpEntryCreationDate, &followUpEntry.FollowUpEntryUpdateDate,
		&followUpEntry.FollowUpEntryOwnerGeneralUser, &followUpEntry.FollowUpEntryWasDone, &followUpEntry.FollowUpEntrySector, &followUpEntry.FollowUpEntryPersonVisitedEntity, &followUpEntry.FollowUpEntryReceivedAttention,
		&followUpEntry.FollowUpEntryComments, &followUpEntry.FollowUpEntryCaseDocumentsPrepared, &followUpEntry.FollowUpEntryStatus, &followUpEntry.FollowUpEntryFollowUp)

	rel.RelBarrierFollowUpEntryBarrier = barrier.ToDTO()
	rel.RelBarrierFollowUpEntryFollowUpEntry = followUpEntry.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetRelBarrierFollowUpEntries obtiene una lista de relaciones entre barreras y entradas de seguimiento basándose en criterios de búsqueda.
// Parámetros:
//   - by: estructura que define los criterios de búsqueda.
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna una lista de DTOs con los datos de las relaciones y un error en caso de fallo.
func GetRelBarrierFollowUpEntries(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelBarrierFollowUpEntryDTO, error) {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var relPath string = RelBarrierFollowUpEntryDBScheme + "." + RelBarrierFollowUpEntryDBName
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var followUpEntryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	var rel RelBarrierFollowUpEntryDTO
	var rels []RelBarrierFollowUpEntryDTO

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelBarrierFollowUpEntryId", "RelBarrierFollowUpEntryBarrier", "RelBarrierFollowUpEntryFollowUpEntry"}
	var relFieldsAliasSlice []string = []string{}

	var barrierFieldsSlice []string = []string{"BarrierICode", "BarrierName", "BarrierDescription"}
	var barrierFieldsAliasSlice []string = []string{}

	var followUpEntryFieldsSlice []string = []string{"FollowUpEntryICode", "FollowUpEntryCreationDate", "FollowUpEntryUpdateDate", "FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone", "FollowUpEntrySector", "FollowUpEntryPersonVisitedEntity", "FollowUpEntryReceivedAttention",
		"FollowUpEntryComments", "FollowUpEntryCaseDocumentsPrepared", "FollowUpEntryStatus", "FollowUpEntryFollowUp"}
	var followUpEntryFieldsAliasSlice []string = []string{}

	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelBarrierFollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, true)
	var barrierFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, true)
	var followUpEntryFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, followUpEntryFieldsSlice, followUpEntryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	var sbFieldsSlice []string = []string{"SectorBarrierId", "SectorBarrierICode", "SectorBarrierName", "SectorBarrierSector"}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	var query string = `SELECT ` + relFieldsStr + `, ` + barrierFieldsStr + `, ` + followUpEntryFieldsStr + ", " + sbFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + barrierPath + ` ON (` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryBarrier"].DBName + `)` +
		` LEFT JOIN ` + followUpEntryPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryFollowUpEntry"].DBName + `)` +
		` LEFT JOIN ` + sectorBarrierPath + ` ON (` + sectorBarrierPath + `.` + SectorBarrierFieldDefinitions["SectorBarrierId"].DBName + ` = ` + barrierPath + `.` + BarrierFieldDefinitions["BarrierSectorBarrier"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, RelBarrierFollowUpEntryDBName, by.AttrsName, []string{}, []string{}, by.Operator, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	for persistenceCtrl.Next() {
		var barrierPg BarrierPgDB = BarrierPgDB{}
		var followUpEntryPg FollowUpEntryPgDB = FollowUpEntryPgDB{}
		var sectorBarrierPg = SectorBarrierPgDB{}

		persistenceCtrl.ScanRow(&rel.RelBarrierFollowUpEntryId, &barrierPg.BarrierId, &followUpEntryPg.FollowUpEntryId,
			&barrierPg.BarrierICode, &barrierPg.BarrierName, &barrierPg.BarrierDescription, &followUpEntryPg.FollowUpEntryICode, &followUpEntryPg.FollowUpEntryCreationDate,
			&followUpEntryPg.FollowUpEntryUpdateDate, &followUpEntryPg.FollowUpEntryOwnerGeneralUser, &followUpEntryPg.FollowUpEntryWasDone, &followUpEntryPg.FollowUpEntrySector,
			&followUpEntryPg.FollowUpEntryPersonVisitedEntity, &followUpEntryPg.FollowUpEntryReceivedAttention, &followUpEntryPg.FollowUpEntryComments,
			&followUpEntryPg.FollowUpEntryCaseDocumentsPrepared, &followUpEntryPg.FollowUpEntryStatus, &followUpEntryPg.FollowUpEntryFollowUp,
			&sectorBarrierPg.SectorBarrierId, &sectorBarrierPg.SectorBarrierICode, &sectorBarrierPg.SectorBarrierName, &sectorBarrierPg.SectorBarrierSector)

		rel.RelBarrierFollowUpEntryBarrier = barrierPg.ToDTO()
		rel.RelBarrierFollowUpEntryFollowUpEntry = followUpEntryPg.ToDTO()
		rel.RelBarrierFollowUpEntryBarrier.BarrierSectorBarrier = sectorBarrierPg.ToDTO()

		rels = append(rels, rel)
	}
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// GetAllRelBarrierFollowUpEntry obtiene todas las relaciones entre barreras y entradas de seguimiento sin aplicar filtros.
// Parámetros:
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna una lista de DTOs con todas las relaciones y un error en caso de fallo.
func GetAllRelBarrierFollowUpEntry(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelBarrierFollowUpEntryDTO, error) {
	var rel RelBarrierFollowUpEntryDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelBarrierFollowUpEntryId", "RelBarrierFollowUpEntryBarrier", "RelBarrierFollowUpEntryFollowUpEntry"}
	var relFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, relFieldsSlice, relFieldsAliasSlice, RelBarrierFollowUpEntryDBName, []string{}, []string{}, []string{}, "", RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query)
	var rels []RelBarrierFollowUpEntryDTO
	for persistenceCtrl.Next() {
		rel = RelBarrierFollowUpEntryDTO{}
		var barrier BarrierPgDB = BarrierPgDB{}
		var followUpEntry FollowUpEntryPgDB = FollowUpEntryPgDB{}

		persistenceCtrl.ScanRow(&rel.RelBarrierFollowUpEntryId, &barrier.BarrierId, &followUpEntry.FollowUpEntryId)

		rel.RelBarrierFollowUpEntryBarrier = barrier.ToDTO()
		rel.RelBarrierFollowUpEntryFollowUpEntry = followUpEntry.ToDTO()

		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// GetAllRelBarrierFollowUpEntryWithBarrierAndFollowUpEntry obtiene todas las relaciones entre barreras y entradas de seguimiento
// junto con los detalles asociados de la barrera y la entrada de seguimiento.
// Parámetros:
//   - connData: datos de la conexión a la base de datos.
//   - clientConfig: configuración del cliente de la base de datos.
//   - serverConfig: configuración del servidor de la base de datos.
//
// Retorna una lista de DTOs con las relaciones y sus datos asociados, y un error en caso de fallo.
func GetAllRelBarrierFollowUpEntryWithBarrierAndFollowUpEntry(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]RelBarrierFollowUpEntryDTO, error) {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	var relPath string = RelBarrierFollowUpEntryDBScheme + "." + RelBarrierFollowUpEntryDBName
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var followUpEntryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var relFieldsSlice []string = []string{"RelBarrierFollowUpEntryId", "RelBarrierFollowUpEntryBarrier", "RelBarrierFollowUpEntryFollowUpEntry"}
	var relFieldsAliasSlice []string = []string{}

	var barrierFieldsSlice []string = []string{"BarrierICode", "BarrierName", "BarrierDescription", "BarrierSectorBarrier"}
	var barrierFieldsAliasSlice []string = []string{}

	var followUpEntryFieldsSlice []string = []string{"FollowUpEntryICode", "FollowUpEntryCreationDate", "FollowUpEntryUpdateDate", "FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone", "FollowUpEntrySector", "FollowUpEntryPersonVisitedEntity", "FollowUpEntryReceivedAttention",
		"FollowUpEntryComments", "FollowUpEntryCaseDocumentsPrepared", "FollowUpEntryStatus", "FollowUpEntryFollowUp"}
	var followUpEntryFieldsAliasSlice []string = []string{}

	var relFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, relFieldsSlice, relFieldsAliasSlice, RelBarrierFollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, RelBarrierFollowUpEntryDBScheme, RelBarrierFollowUpEntryFieldDefinitions, true)
	var barrierFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, true)
	var followUpEntryFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, followUpEntryFieldsSlice, followUpEntryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	var query string = `SELECT ` + relFieldsStr + ", " + barrierFieldsStr + ", " + followUpEntryFieldsStr +
		` FROM ` + relPath +
		` LEFT JOIN ` + barrierPath + ` ON (` + barrierPath + `.` + BarrierFieldDefinitions["BarrierId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryBarrier"].DBName + `)` +
		` LEFT JOIN ` + followUpEntryPath + ` ON (` + followUpEntryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryId"].DBName + ` = ` + relPath + `.` + RelBarrierFollowUpEntryFieldDefinitions["RelBarrierFollowUpEntryFollowUpEntry"].DBName + `)`

	persistenceCtrl.Query(context.Background(), query)
	var rels []RelBarrierFollowUpEntryDTO
	for persistenceCtrl.Next() {
		var barrier BarrierPgDB = BarrierPgDB{}
		var followUpEntry FollowUpEntryPgDB = FollowUpEntryPgDB{}
		var rel = RelBarrierFollowUpEntryDTO{}

		persistenceCtrl.ScanRow(&rel.RelBarrierFollowUpEntryId, &barrier.BarrierId, &followUpEntry.FollowUpEntryId,
			&barrier.BarrierICode, &barrier.BarrierName, &barrier.BarrierDescription, &barrier.BarrierSectorBarrier, &followUpEntry.FollowUpEntryICode, &followUpEntry.FollowUpEntryCreationDate,
			&followUpEntry.FollowUpEntryUpdateDate, &followUpEntry.FollowUpEntryOwnerGeneralUser, &followUpEntry.FollowUpEntryWasDone,
			&followUpEntry.FollowUpEntrySector, &followUpEntry.FollowUpEntryPersonVisitedEntity, &followUpEntry.FollowUpEntryReceivedAttention,
			&followUpEntry.FollowUpEntryComments, &followUpEntry.FollowUpEntryCaseDocumentsPrepared, &followUpEntry.FollowUpEntryStatus, &followUpEntry.FollowUpEntryFollowUp)

		rel.RelBarrierFollowUpEntryBarrier = barrier.ToDTO()
		rel.RelBarrierFollowUpEntryFollowUpEntry = followUpEntry.ToDTO()

		rels = append(rels, rel)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return rels, nil
}

// SetRelBarrierFollowUpEntryDefaults asigna valores por defecto a ciertos campos de la relación en función de la acción que se realizará.
// Parámetros:
//   - rel: puntero al DTO de la relación a modificar.
//   - action: acción que se va a realizar (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetRelBarrierFollowUpEntryDefaults(rel *RelBarrierFollowUpEntryDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		rel.RelBarrierFollowUpEntryCreationDate = time.Now()
		// No default values for insert
	case common_dao.SQL_UPDATE:
		// No default values for update
	}
}

// PgDBToDTO convierte una instancia de RelBarrierFollowUpEntryPgDB (modelo de datos de PostgreSQL)
// a su correspondiente objeto de transferencia de datos (DTO) RelBarrierFollowUpEntryDTO.
// Retorna el DTO con los datos convertidos.
func (obj *RelBarrierFollowUpEntryPgDB) PgDBToDTO() RelBarrierFollowUpEntryDTO {
	var dto RelBarrierFollowUpEntryDTO

	if obj.RelBarrierFollowUpEntryId.Valid {
		dto.RelBarrierFollowUpEntryId = uint64(obj.RelBarrierFollowUpEntryId.Int64)
	}

	if obj.RelBarrierFollowUpEntryCreationDate.Valid {
		dto.RelBarrierFollowUpEntryCreationDate = obj.RelBarrierFollowUpEntryCreationDate.Time
	}

	if obj.RelBarrierFollowUpEntryBarrier.Valid {
		dto.RelBarrierFollowUpEntryBarrier = BarrierDTO{BarrierId: uint64(obj.RelBarrierFollowUpEntryBarrier.Int64)}
	}

	if obj.RelBarrierFollowUpEntryFollowUpEntry.Valid {
		dto.RelBarrierFollowUpEntryFollowUpEntry = FollowUpEntryDTO{FollowUpEntryId: uint64(obj.RelBarrierFollowUpEntryFollowUpEntry.Int64)}
	}

	return dto
}
