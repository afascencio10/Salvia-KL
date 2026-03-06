// El paquete salvia_daos contiene los Objetos de Acceso a Datos (DAO) para las entidades
// relacionadas con la lógica de negocio de "salvia". Estos DAO encapsulan la lógica de
// interacción con la base de datos para cada entidad.
package salvia_daos

import (
	// Importaciones de paquetes del proyecto y de la librería estándar de Go.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers" // Proporciona controladores de persistencia.
	common_dao "bitsflow/common/dao"                 // Contiene constantes y helpers para operaciones DAO.
	"bitsflow/common/db"                             // Maneja la configuración y conexión a la base de datos.
	"bitsflow/common/utils"                          // Utilidades generales, como la generación de UUIDs.
	"context"                                        // Para el manejo de contextos en operaciones de base de datos.
	"database/sql"                                   // Proporciona la interfaz genérica para bases de datos SQL.
	"errors"
	"fmt" // Paquete para formateo de I/O, usado aquí para imprimir errores.
	"time"
)

// Este bloque var define metadatos para la entidad Barrier.
// Estos valores se utilizan para construir consultas SQL de forma dinámica y consistente,
// y para mantener una única fuente de verdad sobre los nombres de la entidad, tabla y esquema.
var (
	// BarrierEntityName es el nombre lógico de la entidad en la aplicación.
	BarrierEntityName string = "Barrier"
	// BarrierJSONName es el nombre utilizado para la serialización/deserialización JSON.
	BarrierJSONName string = "barrier"
	// BarrierDBName es el nombre de la tabla en la base de datos.
	BarrierDBName string = "barrier"
	// BarrierDBScheme es el nombre del esquema de la base de datos donde se encuentra la tabla.
	BarrierDBScheme string = "salvia"

	// BarrierFieldDefinitions es un mapa que define las propiedades de cada campo de la entidad Barrier.
	// Se utiliza para la validación y la construcción dinámica de consultas SQL, mapeando los campos del
	// struct de Go a las columnas de la base de datos y especificando sus tipos y restricciones.
	BarrierFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"BarrierId":            {Name: "BarrierId", DBName: "barrier_id", Alias: "", ModelType: "uint", Required: false},
		"BarrierICode":         {Name: "BarrierICode", DBName: "barrier_i_code", Alias: "", ModelType: "string", MaxSize: 36, Required: false},
		"BarrierCreationDate":  {Name: "BarrierCreationDate", DBName: "barrier_creation_date", Alias: "", ModelType: "datetime", Required: false},
		"BarrierUpdateDate":    {Name: "BarrierUpdateDate", DBName: "barrier_update_date", Alias: "", ModelType: "datetime", Required: false},
		"BarrierName":          {Name: "BarrierName", DBName: "barrier_name", Alias: "", ModelType: "string", MaxSize: 36, Required: true},
		"BarrierDescription":   {Name: "BarrierDescription", DBName: "barrier_description", Alias: "", ModelType: "string", MaxSize: 500, Required: false},
		"BarrierSectorBarrier": {Name: "BarrierSectorBarrier", DBName: "sector_barrier_id", Alias: "", ModelType: "uint", Required: true},
	}
)

// BarrierDTO (Data Transfer Object) define la estructura de datos para la entidad Barrier.
// Se utiliza para transferir datos entre las capas de la aplicación (por ejemplo, en las
// respuestas de una API) y contiene los campos que son relevantes para el cliente.
type BarrierDTO struct {
	// BarrierId es el identificador único de la barrera. Se omite en JSON (`-`).
	BarrierId uint64 `json:"-"`
	// BarrierICode es un código identificador interno único (posiblemente un UUID).
	BarrierICode string `json:"icode"`
	//BarrierCreationDate Fecha de creación
	BarrierCreationDate time.Time `json:"creationDate"`
	//BarrierCreationDate Fecha de modificación
	BarrierUpdateDate time.Time `json:"updateDate"`
	// BarrierName es el nombre de la barrera.
	BarrierName string `json:"name"`
	// BarrierDescription es una descripción detallada de la barrera.
	BarrierDescription string `json:"description"`
	// BarrierSectorBarrier es el sector al que pertenece la barrera.
	BarrierSectorBarrier SectorBarrierDTO `json:"sectorBarrier"`

	//Campos que se usan en el formulario
	BarrierFollowUpEntriesActing []FollowUpEntryActingDTO `json:"acting"`
}

// BarrierPgDB representa la estructura de una fila de la tabla 'barrier' en la base de datos PostgreSQL.
// Utiliza los tipos `sql.Null*` para manejar correctamente los valores que pueden ser NULL en la base de datos,
// evitando errores en tiempo de ejecución al escanear los resultados de una consulta.
type BarrierPgDB struct {
	BarrierId            sql.NullInt64  // Mapea a la columna 'barrier_id' (BIGINT).
	BarrierICode         sql.NullString // Mapea a la columna 'barrier_i_code' (CHARACTER VARYING).
	BarrierCreationDate  sql.NullTime
	BarrierUpdateDate    sql.NullTime
	BarrierName          sql.NullString // Mapea a la columna 'barrier_name' (CHARACTER VARYING).
	BarrierDescription   sql.NullString // Mapea a la columna 'barrier_description' (CHARACTER VARYING).
	BarrierSectorBarrier sql.NullInt64  // Mapea a la columna 'barrier_sector' (CHARACTER VARYING).
}

// SetBarrier inserta una nueva barrera en la base de datos.
// Recibe el objeto barrera, información de transacción, módulo, y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func SetBarrier(barrier *BarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se define la lista de campos que se insertarán en la BD.
	var barrierFieldsSlice []string = []string{"BarrierICode", "BarrierCreationDate", "BarrierUpdateDate", "BarrierName", "BarrierDescription", "BarrierSectorBarrier"}
	var barrierFieldsAliasSlice []string = []string{}

	// Se genera la query de inserción utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{"BarrierId"}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query, barrier.BarrierICode, barrier.BarrierCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), barrier.BarrierUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT), barrier.BarrierName, barrier.BarrierDescription, barrier.BarrierSectorBarrier)
	persistenceCtrl.Scan(&barrier.BarrierId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateBarrierByICode actualiza los datos de una barrera en la BD, utilizando su ICode como referencia.
// Recibe el objeto barrera con los nuevos datos, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func UpdateBarrierByICode(barrier *BarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se actualizarán.
	var barrierFieldsSlice []string = []string{"BarrierName", "BarrierDescription"}
	var barrierFieldsAliasSlice []string = []string{}

	// Se genera la query de actualización utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{"BarrierICode"}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		barrier.BarrierName, barrier.BarrierDescription, barrier.BarrierICode)

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

// RemoveBarrierByICode elimina una barrera de la base de datos utilizando su ICode.
// Recibe el objeto barrera, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func RemoveBarrierByICode(barrier *BarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos y condiciones para la eliminación.
	var barrierFieldsSlice []string = []string{}
	var barrierFieldsAliasSlice []string = []string{}

	// Se genera la query de eliminación utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{"BarrierICode"}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, false)

	// Se ejecuta la query con el parámetro correspondiente.
	persistenceCtrl.Exec(context.Background(), query, barrier.BarrierICode)

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

// GetBarrier consulta la barrera basada en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto barrera a completar,
// información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func GetBarrier(by common_controllers.By, barrier *BarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla Barrier en la BD.
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var barrierFieldsSlice []string = []string{"BarrierId", "BarrierICode", "BarrierCreationDate", "BarrierUpdateDate", "BarrierName", "BarrierDescription", "BarrierSectorBarrier"}
	var barrierFieldsAliasSlice []string = []string{}
	var barrierFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, true)

	var sbFieldsSlice []string = []string{"SectorBarrierId", "SectorBarrierICode", "SectorBarrierName", "SectorBarrierSector"}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + barrierFieldsStr + ", " + sbFieldsStr +
		` FROM ` + barrierPath +
		` LEFT JOIN ` + sectorBarrierPath + ` ON (` + sectorBarrierPath + `.` + SectorBarrierFieldDefinitions["SectorBarrierId"].DBName + ` = ` + barrierPath + `.` + BarrierFieldDefinitions["BarrierSectorBarrier"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, BarrierDBName, by.AttrsName, []string{}, []string{}, by.Operator, BarrierDBScheme, BarrierFieldDefinitions, true)

	// Se ejecuta la query con los parámetros de filtrado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var barrierPg = BarrierPgDB{}
	var sectorBarrierPg = SectorBarrierPgDB{}

	// Se escanean los resultados de la query.
	persistenceCtrl.Scan(&barrierPg.BarrierId, &barrierPg.BarrierICode, &barrierPg.BarrierCreationDate, &barrierPg.BarrierUpdateDate, &barrierPg.BarrierName, &barrierPg.BarrierDescription, &barrierPg.BarrierSectorBarrier,
		&sectorBarrierPg.SectorBarrierId, &sectorBarrierPg.SectorBarrierICode, &sectorBarrierPg.SectorBarrierName, &sectorBarrierPg.SectorBarrierSector)

	// Se convierte el objeto de base de datos a DTO.
	*barrier = barrierPg.ToDTO()
	barrier.BarrierSectorBarrier = sectorBarrierPg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetBarrier consulta barreras basado en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto barrera a completar,
// información de transacción, módulo y datos de conexión.
// Retorna las barreras o error.
func GetBarriers(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]BarrierDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla Barrier en la BD.
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var barrierFieldsSlice []string = []string{"BarrierId", "BarrierICode", "BarrierCreationDate", "BarrierUpdateDate", "BarrierName", "BarrierDescription", "BarrierSectorBarrier"}
	var barrierFieldsAliasSlice []string = []string{}
	var barrierFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, true)

	var sbFieldsSlice []string = []string{"SectorBarrierId", "SectorBarrierICode", "SectorBarrierName", "SectorBarrierSector"}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + barrierFieldsStr + ", " + sbFieldsStr +
		` FROM ` + barrierPath +
		` LEFT JOIN ` + sectorBarrierPath + ` ON (` + sectorBarrierPath + `.` + SectorBarrierFieldDefinitions["SectorBarrierId"].DBName + ` = ` + barrierPath + `.` + BarrierFieldDefinitions["BarrierSectorBarrier"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, BarrierDBName, by.AttrsName, []string{}, []string{}, by.Operator, BarrierDBScheme, BarrierFieldDefinitions, true)

	// Se ejecuta la query con los parámetros de filtrado.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	var barriers []BarrierDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var barrierPg = BarrierPgDB{}
		var sectorBarrierPg = SectorBarrierPgDB{}
		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&barrierPg.BarrierId, &barrierPg.BarrierICode, &barrierPg.BarrierCreationDate, &barrierPg.BarrierUpdateDate, &barrierPg.BarrierName, &barrierPg.BarrierDescription, &barrierPg.BarrierSectorBarrier,
			&sectorBarrierPg.SectorBarrierId, &sectorBarrierPg.SectorBarrierICode, &sectorBarrierPg.SectorBarrierName, &sectorBarrierPg.SectorBarrierSector)

		var barrier BarrierDTO = barrierPg.ToDTO()
		barrier.BarrierSectorBarrier = sectorBarrierPg.ToDTO()
		barriers = append(barriers, barrier)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return barriers, nil
}

// GetAllBarrier consulta todas las barreras existentes en la base de datos.
// Recibe información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada, un slice de BarrierDTO y un error en caso de producirse.
func GetAllBarriers(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]BarrierDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var barrierPath string = BarrierDBScheme + "." + BarrierDBName
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar para las barreras.
	var barrierFieldsSlice []string = []string{"BarrierId", "BarrierICode", "BarrierCreationDate", "BarrierUpdateDate", "BarrierName", "BarrierDescription", "BarrierSectorBarrier"}
	var barrierFieldsAliasSlice []string = []string{}
	var barrierFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, barrierFieldsSlice, barrierFieldsAliasSlice, BarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, BarrierDBScheme, BarrierFieldDefinitions, true)

	var sbFieldsSlice []string = []string{"SectorBarrierId", "SectorBarrierICode", "SectorBarrierName", "SectorBarrierSector"}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	// Se genera la query para obtener todas las barreras.
	var query string = `SELECT ` + barrierFieldsStr + ", " + sbFieldsStr +
		` FROM ` + barrierPath +
		` LEFT JOIN ` + sectorBarrierPath + ` ON (` + sectorBarrierPath + `.` + SectorBarrierFieldDefinitions["SectorBarrierId"].DBName + ` = ` + barrierPath + `.` + BarrierFieldDefinitions["BarrierSectorBarrier"].DBName + `)` +
		` WHERE TRUE`

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query)
	var barriers []BarrierDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var barrierPg = BarrierPgDB{}
		var sectorBarrierPg = SectorBarrierPgDB{}
		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(&barrierPg.BarrierId, &barrierPg.BarrierICode, &barrierPg.BarrierCreationDate, &barrierPg.BarrierUpdateDate, &barrierPg.BarrierName, &barrierPg.BarrierDescription, &barrierPg.BarrierSectorBarrier)

		var barrier BarrierDTO = barrierPg.ToDTO()
		barrier.BarrierSectorBarrier = sectorBarrierPg.ToDTO()

		barriers = append(barriers, barrier)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return barriers, nil
}

// SetBarrierDefaults asigna valores por defecto a los campos de una barrera,
// dependiendo de la acción que se esté realizando (insertar o actualizar).
func SetBarrierDefaults(barrier *BarrierDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		barrier.BarrierCreationDate = time.Now()
		barrier.BarrierUpdateDate = time.Now()
		barrier.BarrierICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		barrier.BarrierUpdateDate = time.Now()
	}
}

// ToDTO convierte un objeto BarrierPgDB obtenido de la base de datos en un objeto BarrierDTO.
// Se encarga de verificar la validez de los campos nulos y asignarlos correctamente.
func (obj *BarrierPgDB) ToDTO() BarrierDTO {
	var dto BarrierDTO

	if obj.BarrierId.Valid {
		dto.BarrierId = uint64(obj.BarrierId.Int64)
	}

	if obj.BarrierICode.Valid {
		dto.BarrierICode = obj.BarrierICode.String
	}

	if obj.BarrierName.Valid {
		dto.BarrierName = obj.BarrierName.String
	}

	if obj.BarrierDescription.Valid {
		dto.BarrierDescription = obj.BarrierDescription.String
	}

	if obj.BarrierSectorBarrier.Valid {
		dto.BarrierSectorBarrier = SectorBarrierDTO{SectorBarrierId: uint64(obj.BarrierSectorBarrier.Int64)}
	}

	return dto
}
