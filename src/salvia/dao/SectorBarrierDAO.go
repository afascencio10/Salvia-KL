// El paquete salvia_daos contiene los Objetos de Acceso a Datos (DAO) para las entidades
// de la lógica de negocio de "salvia". Este archivo en particular maneja la persistencia
// de la entidad 'SectorBarrier' (Barrera de Sector).
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

// Este bloque var define metadatos constantes para la entidad SectorBarrier.
// Estos valores centralizan la configuración de la entidad, facilitando el mantenimiento y
// asegurando la consistencia en la construcción de consultas SQL y en la serialización de datos.
var (
	// SectorBarrierEntityName es el nombre lógico de la entidad dentro de la aplicación.
	SectorBarrierEntityName string = "SectorBarrier"
	// SectorBarrierJSONName es el nombre que se utilizará para la entidad en formatos JSON.
	SectorBarrierJSONName string = "sectorBarrier"
	// SectorBarrierDBName es el nombre exacto de la tabla en la base de datos.
	SectorBarrierDBName string = "sector_barrier"
	// SectorBarrierDBScheme es el esquema de la base de datos al que pertenece la tabla.
	SectorBarrierDBScheme string = "salvia"

	// SectorBarrierFieldDefinitions es un mapa que detalla las propiedades de cada campo del modelo SectorBarrier.
	// Asocia los nombres de los campos del struct de Go con los nombres de las columnas de la base de datos,
	// y define metadatos como el tipo de dato, tamaño y si es requerido. Es crucial para la construcción
	// dinámica de consultas SQL y para validaciones.
	SectorBarrierFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"SectorBarrierId":           {Name: "SectorBarrierId", DBName: "sector_barrier_id", Alias: "", ModelType: "uint", Required: false},
		"SectorBarrierICode":        {Name: "SectorBarrierICode", DBName: "sector_barrier_i_code", Alias: "", ModelType: "string", MinSize: 36, MaxSize: 36, Required: true},
		"SectorBarrierCreationDate": {Name: "SectorBarrierCreationDate", DBName: "sector_barrier_creation_date", Alias: "", ModelType: "datetime", Required: true},
		"SectorBarrierUpdateDate":   {Name: "SectorBarrierUpdateDate", DBName: "sector_barrier_update_date", Alias: "", ModelType: "datetime", Required: true},
		"SectorBarrierName":         {Name: "SectorBarrierName", DBName: "sector_barrier_name", Alias: "", ModelType: "string", MaxSize: 36, Required: true},
		"SectorBarrierDescription":  {Name: "SectorBarrierDescription", DBName: "sector_barrier_description", Alias: "", ModelType: "string", MaxSize: 500, Required: false},
		"SectorBarrierSector":       {Name: "SectorBarrierSector", DBName: "sector_barrier_sector", Alias: "", ModelType: "string", MaxSize: 2, Required: false},
	}
)

// SectorBarrierDTO (Data Transfer Object) es la estructura que representa una barrera de sector en la capa de aplicación.
// Se utiliza para intercambiar datos con el exterior (ej. API) y para la lógica de negocio.
type SectorBarrierDTO struct {
	SectorBarrierId           uint64    `json:"-"`
	SectorBarrierICode        string    `json:"icode"`
	SectorBarrierCreationDate time.Time `json:"creationDate"`
	SectorBarrierUpdateDate   time.Time `json:"updateDate"`
	SectorBarrierName         string    `json:"name"`
	SectorBarrierDescription  string    `json:"description"`
	SectorBarrierSector       string    `json:"sector"`
}

// SectorBarrierPgDB representa la estructura de la tabla 'sector_barrier' de la base de datos.
// Utiliza tipos `sql.Null*` para manejar de forma segura las columnas que pueden contener
// valores NULL, evitando errores al escanear los resultados de las consultas.
type SectorBarrierPgDB struct {
	SectorBarrierId           sql.NullInt64
	SectorBarrierICode        sql.NullString
	SectorBarrierCreationDate sql.NullTime
	SectorBarrierUpdateDate   sql.NullTime
	SectorBarrierName         sql.NullString
	SectorBarrierDescription  sql.NullString
	SectorBarrierSector       sql.NullString
}

func (sb SectorBarrierDTO) MarshalJSON() ([]byte, error) {
	type Alias SectorBarrierDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		SectorBarrierCreationDate string `json:"creationDate"`
		SectorBarrierUpdateDate   string `json:"updateDate"`
	}{
		Alias:                     (*Alias)(&sb),
		SectorBarrierCreationDate: sb.SectorBarrierCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		SectorBarrierUpdateDate:   sb.SectorBarrierUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (sb *SectorBarrierDTO) UnmarshalJSON(data []byte) error {
	type Alias SectorBarrierDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		SectorBarrierCreationDate string `json:"creationDate"`
		SectorBarrierUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(sb),
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
	sb.SectorBarrierCreationDate = parse(aux.SectorBarrierCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	sb.SectorBarrierUpdateDate = parse(aux.SectorBarrierUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetSectorBarrier inserta una nueva barrera de sector en la base de datos.
// Recibe el objeto barrera de sector, información de transacción, módulo, y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func SetSectorBarrier(sectorBarrier *SectorBarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var sbFieldsSlice []string = []string{
		"SectorBarrierICode",
		"SectorBarrierCreationDate",
		"SectorBarrierUpdateDate",
		"SectorBarrierName",
		"SectorBarrierDescription",
		"SectorBarrierSector",
	}
	var sbFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{"SectorBarrierId"}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, false)

	persistenceCtrl.QueryRow(context.Background(), query,
		sectorBarrier.SectorBarrierICode,
		sectorBarrier.SectorBarrierCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		sectorBarrier.SectorBarrierUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		sectorBarrier.SectorBarrierName,
		sectorBarrier.SectorBarrierDescription,
		sectorBarrier.SectorBarrierSector,
	)
	persistenceCtrl.Scan(&sectorBarrier.SectorBarrierId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateSectorBarrierByICode actualiza los datos de una barrera de sector en la BD, utilizando su ICode como referencia.
// Recibe el objeto barrera de sector con los nuevos datos, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func UpdateSectorBarrierByICode(sectorBarrier *SectorBarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var sbFieldsSlice []string = []string{
		"SectorBarrierUpdateDate",
		"SectorBarrierName",
		"SectorBarrierDescription",
		"SectorBarrierSector",
	}
	var sbFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{"SectorBarrierICode"}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query,
		sectorBarrier.SectorBarrierICode,
		sectorBarrier.SectorBarrierUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		sectorBarrier.SectorBarrierName,
		sectorBarrier.SectorBarrierDescription,
		sectorBarrier.SectorBarrierSector,
	)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveSectorBarrierByICode elimina una barrera de sector de la base de datos utilizando su ICode.
// Recibe el objeto barrera de sector, información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func RemoveSectorBarrierByICode(sectorBarrier *SectorBarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var sbFieldsSlice []string = []string{}
	var sbFieldsAliasSlice []string = []string{}

	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{"SectorBarrierICode"}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, false)

	persistenceCtrl.Exec(context.Background(), query, sectorBarrier.SectorBarrierICode)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetSectorBarrier consulta una barrera de sector basado en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto barrera de sector a completar,
// información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada y un error en caso de producirse.
func GetSectorBarrier(by common_controllers.By, sectorBarrier *SectorBarrierDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	var sbFieldsSlice []string = []string{
		"SectorBarrierId",
		"SectorBarrierICode",
		"SectorBarrierCreationDate",
		"SectorBarrierUpdateDate",
		"SectorBarrierName",
		"SectorBarrierDescription",
		"SectorBarrierSector",
	}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	var query string = `SELECT ` + sbFieldsStr +
		` FROM ` + sectorBarrierPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, SectorBarrierDBName, by.AttrsName, []string{}, []string{}, by.Operator, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var sectorBarrierPg = SectorBarrierPgDB{}

	persistenceCtrl.Scan(
		&sectorBarrierPg.SectorBarrierId,
		&sectorBarrierPg.SectorBarrierICode,
		&sectorBarrierPg.SectorBarrierCreationDate,
		&sectorBarrierPg.SectorBarrierUpdateDate,
		&sectorBarrierPg.SectorBarrierName,
		&sectorBarrierPg.SectorBarrierDescription,
		&sectorBarrierPg.SectorBarrierSector,
	)

	*sectorBarrier = sectorBarrierPg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetSectorBarriers consulta barreras de sector basado en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto barrera de sector a completar,
// información de transacción, módulo y datos de conexión.
// Retorna las barreras de sector o un error en caso de producirse.
func GetSectorBarriers(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]SectorBarrierDTO, error) {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var sbFieldsSlice []string = []string{
		"SectorBarrierId",
		"SectorBarrierICode",
		"SectorBarrierCreationDate",
		"SectorBarrierUpdateDate",
		"SectorBarrierName",
		"SectorBarrierDescription",
		"SectorBarrierSector",
	}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	var query string = `SELECT ` + sbFieldsStr +
		` FROM ` + sectorBarrierPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, SectorBarrierDBName, by.AttrsName, []string{}, []string{}, by.Operator, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query)
	var sectorBarriers []SectorBarrierDTO
	for persistenceCtrl.Next() {
		var sectorBarrierPg = SectorBarrierPgDB{}
		persistenceCtrl.ScanRow(
			&sectorBarrierPg.SectorBarrierId,
			&sectorBarrierPg.SectorBarrierICode,
			&sectorBarrierPg.SectorBarrierCreationDate,
			&sectorBarrierPg.SectorBarrierUpdateDate,
			&sectorBarrierPg.SectorBarrierName,
			&sectorBarrierPg.SectorBarrierDescription,
			&sectorBarrierPg.SectorBarrierSector,
		)

		var sectorBarrier SectorBarrierDTO = sectorBarrierPg.ToDTO()
		sectorBarriers = append(sectorBarriers, sectorBarrier)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return sectorBarriers, nil
}

// GetAllSectorBarriers consulta todas las barreras de sector existentes en la base de datos.
// Recibe información de transacción, módulo y datos de conexión.
// Retorna la conexión actualizada, un slice de SectorBarrierDTO y un error en caso de producirse.
func GetAllSectorBarriers(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]SectorBarrierDTO, error) {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var sectorBarrierPath string = SectorBarrierDBScheme + "." + SectorBarrierDBName

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	var sbFieldsSlice []string = []string{
		"SectorBarrierId",
		"SectorBarrierICode",
		"SectorBarrierCreationDate",
		"SectorBarrierUpdateDate",
		"SectorBarrierName",
		"SectorBarrierDescription",
		"SectorBarrierSector",
	}
	var sbFieldsAliasSlice []string = []string{}
	var sbFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, sbFieldsSlice, sbFieldsAliasSlice, SectorBarrierDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, SectorBarrierDBScheme, SectorBarrierFieldDefinitions, true)

	var query string = `SELECT ` + sbFieldsStr +
		` FROM ` + sectorBarrierPath +
		` WHERE TRUE`

	persistenceCtrl.Query(context.Background(), query)
	var sectorBarriers []SectorBarrierDTO
	for persistenceCtrl.Next() {
		var sectorBarrierPg = SectorBarrierPgDB{}
		persistenceCtrl.ScanRow(
			&sectorBarrierPg.SectorBarrierId,
			&sectorBarrierPg.SectorBarrierICode,
			&sectorBarrierPg.SectorBarrierCreationDate,
			&sectorBarrierPg.SectorBarrierUpdateDate,
			&sectorBarrierPg.SectorBarrierName,
			&sectorBarrierPg.SectorBarrierDescription,
			&sectorBarrierPg.SectorBarrierSector,
		)

		var sectorBarrier SectorBarrierDTO = sectorBarrierPg.ToDTO()
		sectorBarriers = append(sectorBarriers, sectorBarrier)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return sectorBarriers, nil
}

// SetSectorBarrierDefaults asigna valores por defecto a los campos de una barrera de sector,
// dependiendo de la acción que se esté realizando (insertar o actualizar).
func SetSectorBarrierDefaults(sectorBarrier *SectorBarrierDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		sectorBarrier.SectorBarrierCreationDate = time.Now()
		sectorBarrier.SectorBarrierUpdateDate = time.Now()
		sectorBarrier.SectorBarrierICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		sectorBarrier.SectorBarrierUpdateDate = time.Now()
	}
}

// ToDTO convierte un objeto SectorBarrierPgDB obtenido de la base de datos en un objeto SectorBarrierDTO.
// Se encarga de verificar la validez de los campos nulos y asignarlos correctamente.
func (obj *SectorBarrierPgDB) ToDTO() SectorBarrierDTO {
	var dto SectorBarrierDTO

	if obj.SectorBarrierId.Valid {
		dto.SectorBarrierId = uint64(obj.SectorBarrierId.Int64)
	}

	if obj.SectorBarrierICode.Valid {
		dto.SectorBarrierICode = obj.SectorBarrierICode.String
	}

	if obj.SectorBarrierCreationDate.Valid {
		dto.SectorBarrierCreationDate = obj.SectorBarrierCreationDate.Time
	}

	if obj.SectorBarrierUpdateDate.Valid {
		dto.SectorBarrierUpdateDate = obj.SectorBarrierUpdateDate.Time
	}

	if obj.SectorBarrierName.Valid {
		dto.SectorBarrierName = obj.SectorBarrierName.String
	}

	if obj.SectorBarrierDescription.Valid {
		dto.SectorBarrierDescription = obj.SectorBarrierDescription.String
	}

	if obj.SectorBarrierSector.Valid {
		dto.SectorBarrierSector = obj.SectorBarrierSector.String
	}

	return dto
}
