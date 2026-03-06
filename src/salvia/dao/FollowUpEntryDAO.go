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
	// FollowUpEntryEntityName es el nombre lógico de la entidad FollowUpEntry.
	FollowUpEntryEntityName string = "FollowUpEntry"
	// FollowUpEntryJSONName es el nombre JSON para la entidad FollowUpEntry.
	FollowUpEntryJSONName string = "followUpEntry"
	// FollowUpEntryDBName es el nombre de la tabla en la base de datos para FollowUpEntry.
	FollowUpEntryDBName string = "follow_up_entry"
	// FollowUpEntryDBScheme es el esquema de la base de datos para la tabla FollowUpEntry.
	FollowUpEntryDBScheme string = "salvia"

	// FollowUpEntryFieldDefinitions define las propiedades de los campos de la entidad FollowUpEntry,
	// incluyendo su nombre, nombre en la DB, alias, tipo de modelo, tamaño máximo y si es requerido.
	FollowUpEntryFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FollowUpEntryId":                    {Name: "FollowUpEntryId", DBName: "follow_up_entry_id", Alias: "", ModelType: "uint", Required: false},
		"FollowUpEntryICode":                 {Name: "FollowUpEntryICode", DBName: "follow_up_entry_i_code", Alias: "", ModelType: "string", MaxSize: 36, Required: false},
		"FollowUpEntryCreationDate":          {Name: "FollowUpEntryCreationDate", DBName: "follow_up_entry_creation_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpEntryUpdateDate":            {Name: "FollowUpEntryUpdateDate", DBName: "follow_up_entry_update_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpEntryCompletionDate":        {Name: "FollowUpEntryCompletionDate", DBName: "follow_up_entry_completion_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpEntryOwnerGeneralUser":      {Name: "FollowUpEntryOwnerGeneralUser", DBName: "follow_up_entry_owner_general_user", Alias: "", ModelType: "string", MaxSize: 36, Required: false},
		"FollowUpEntryWasDone":               {Name: "FollowUpEntryWasDone", DBName: "follow_up_entry_was_done", Alias: "", ModelType: "string", MaxSize: 1, Required: false},
		"FollowUpEntrySector":                {Name: "FollowUpEntrySector", DBName: "follow_up_entry_sector", Alias: "", ModelType: "string", MaxSize: 2, Required: false},
		"FollowUpEntryPersonVisitedEntity":   {Name: "FollowUpEntryPersonVisitedEntity", DBName: "follow_up_entry_person_visited_entity", Alias: "", ModelType: "string", MaxSize: 1, Required: false},
		"FollowUpEntryReceivedAttention":     {Name: "FollowUpEntryReceivedAttention", DBName: "follow_up_entry_received_attention", Alias: "", ModelType: "string", MaxSize: 1, Required: false},
		"FollowUpEntryComments":              {Name: "FollowUpEntryComments", DBName: "follow_up_entry_comments", Alias: "", ModelType: "string", MaxSize: 5000, Required: false},
		"FollowUpEntryCaseDocumentsPrepared": {Name: "FollowUpEntryCaseDocumentsPrepared", DBName: "follow_up_entry_case_documents_prepared", Alias: "", ModelType: "string", MaxSize: 1, Required: false},
		"FollowUpEntryStatus":                {Name: "FollowUpEntryStatus", DBName: "follow_up_entry_status", Alias: "", ModelType: "string", MaxSize: 1, Required: true},
		"FollowUpEntryFollowUp":              {Name: "FollowUpEntryFollowUp", DBName: "follow_up_id", Alias: "", ModelType: "uint", Required: true},
	}
)

// FollowUpEntryDTO representa el Data Transfer Object para una entrada de seguimiento.
// Se utiliza para transferir datos entre las capas de la aplicación y la presentación.
type FollowUpEntryDTO struct {
	FollowUpEntryId                    uint64      `json:"-"`                     // ID de la entrada de seguimiento (no serializado a JSON).
	FollowUpEntryICode                 string      `json:"icode"`                 // Código interno de la entrada de seguimiento.
	FollowUpEntryCreationDate          time.Time   `json:"creationDate"`          // Fecha de creación de la entrada de seguimiento.
	FollowUpEntryUpdateDate            time.Time   `json:"updateDate"`            // Fecha de última actualización de la entrada de seguimiento.
	FollowUpEntryCompletionDate        time.Time   `json:"completionDate"`        // Fecha de vencimiento del seguimiento
	FollowUpEntryOwnerGeneralUser      string      `json:"ownerGeneralUser"`      // Usuario general propietario de la entrada de seguimiento.
	FollowUpEntryWasDone               string      `json:"wasDone"`               // Indica si se realizó el seguimiento ('S' o 'N').
	FollowUpEntrySector                string      `json:"sector"`                // Sector al que pertenece el seguimiento.
	FollowUpEntryPersonVisitedEntity   string      `json:"personVisitedEntity"`   // Entidad de la persona visitada.
	FollowUpEntryReceivedAttention     string      `json:"receivedAttention"`     // Indica si se recibió atención ('S' o 'N').
	FollowUpEntryComments              string      `json:"comments"`              // Comentarios adicionales sobre el seguimiento.
	FollowUpEntryCaseDocumentsPrepared string      `json:"caseDocumentsPrepared"` // Indica si se prepararon documentos del caso ('S' o 'N').
	FollowUpEntryStatus                string      `json:"status"`                // Estado de la entrada de seguimiento.
	FollowUpEntryFollowUp              FollowUpDTO `json:"followUp"`              // DTO del seguimiento relacionado.
	//Atributos útiles en la captura de datos de la UI
	FollowUpEntryIdentifiedBarriers []BarrierDTO `json:"identifiedBarriers"`
}

// FollowUpEntryPgDB representa la estructura de una entrada de seguimiento tal como se almacena en PostgreSQL.
// Utiliza tipos sql.Null para manejar valores nulos de la base de datos de forma segura.
type FollowUpEntryPgDB struct {
	FollowUpEntryId                    sql.NullInt64  // ID de la entrada de seguimiento.
	FollowUpEntryICode                 sql.NullString // Código interno de la entrada de seguimiento.
	FollowUpEntryCreationDate          sql.NullTime   // Fecha de creación de la entrada de seguimiento.
	FollowUpEntryUpdateDate            sql.NullTime   // Fecha de última actualización de la entrada de seguimiento.
	FollowUpEntryCompletionDate        sql.NullTime   // Fecha de vencimiento del seguimiento
	FollowUpEntryOwnerGeneralUser      sql.NullString // Usuario general propietario de la entrada de seguimiento.
	FollowUpEntryWasDone               sql.NullString // Indica si se realizó el seguimiento.
	FollowUpEntrySector                sql.NullString // Sector al que pertenece el seguimiento.
	FollowUpEntryPersonVisitedEntity   sql.NullString // Entidad de la persona visitada.
	FollowUpEntryReceivedAttention     sql.NullString // Indica si se recibió atención.
	FollowUpEntryComments              sql.NullString // Comentarios adicionales sobre el seguimiento.
	FollowUpEntryCaseDocumentsPrepared sql.NullString // Indica si se prepararon documentos del caso.
	FollowUpEntryStatus                sql.NullString // Estado de la entrada de seguimiento.
	FollowUpEntryFollowUp              sql.NullInt64  // ID del seguimiento relacionado.
}

func (fue FollowUpEntryDTO) MarshalJSON() ([]byte, error) {
	type Alias FollowUpEntryDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		FollowUpEntryCompletionDate string `json:"completionDate"`
		FollowUpEntryCreationDate   string `json:"creationDate"`
		FollowUpEntryUpdateDate     string `json:"updateDate"`
	}{
		Alias:                       (*Alias)(&fue),
		FollowUpEntryCompletionDate: fue.FollowUpEntryCompletionDate.Format(common_config.DateTime.DATE_FORMAT),
		FollowUpEntryCreationDate:   fue.FollowUpEntryCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FollowUpEntryUpdateDate:     fue.FollowUpEntryUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (fue *FollowUpEntryDTO) UnmarshalJSON(data []byte) error {
	// Alias para evitar la recursión infinita al unmarshalizar.
	type Alias FollowUpEntryDTO

	// Estructura auxiliar que recibe las fechas como strings.
	aux := &struct {
		*Alias
		FollowUpEntryCompletionDate string `json:"completionDate"`
		FollowUpEntryCreationDate   string `json:"creationDate"`
		FollowUpEntryUpdateDate     string `json:"updateDate"`
	}{
		Alias: (*Alias)(fue),
	}

	// Unmarshalizamos la estructura auxiliar (sin problemas con las fechas).
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Función helper para parsear fechas de forma tolerante.
	parse := func(value, layout string) time.Time {
		if value == "" {
			return time.Time{}
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			// Fecha inválida → valor cero.
			return time.Time{}
		}
		return t
	}

	// Asignamos los valores parseados a la estructura original.
	fue.FollowUpEntryCompletionDate = parse(aux.FollowUpEntryCompletionDate, common_config.DateTime.DATE_FORMAT)
	fue.FollowUpEntryCreationDate = parse(aux.FollowUpEntryCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	fue.FollowUpEntryUpdateDate = parse(aux.FollowUpEntryUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetFollowUpEntry inserta una nueva entrada de seguimiento en la base de datos.
// Recibe el objeto de entrada de seguimiento, información de transacción, módulo, y datos de conexión.
// Retorna un error en caso de producirse.
func SetFollowUpEntry(followUpEntry *FollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se define la lista de campos que se insertarán en la BD.
	var entryFieldsSlice []string = []string{
		"FollowUpEntryICode",
		"FollowUpEntryCreationDate",
		"FollowUpEntryUpdateDate",
		"FollowUpEntryCompletionDate",
		"FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone",
		"FollowUpEntrySector",
		"FollowUpEntryPersonVisitedEntity",
		"FollowUpEntryReceivedAttention",
		"FollowUpEntryComments",
		"FollowUpEntryCaseDocumentsPrepared",
		"FollowUpEntryStatus",
		"FollowUpEntryFollowUp",
	}
	var entryFieldsAliasSlice []string = []string{}

	// Se genera la query de inserción utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, entryFieldsSlice, entryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{"FollowUpEntryId"}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query,
		followUpEntry.FollowUpEntryICode,
		followUpEntry.FollowUpEntryCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUpEntry.FollowUpEntryUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUpEntry.FollowUpEntryCompletionDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUpEntry.FollowUpEntryOwnerGeneralUser,
		followUpEntry.FollowUpEntryWasDone,
		followUpEntry.FollowUpEntrySector,
		followUpEntry.FollowUpEntryPersonVisitedEntity,
		followUpEntry.FollowUpEntryReceivedAttention,
		followUpEntry.FollowUpEntryComments,
		followUpEntry.FollowUpEntryCaseDocumentsPrepared,
		followUpEntry.FollowUpEntryStatus,
		followUpEntry.FollowUpEntryFollowUp.FollowUpId,
	)
	persistenceCtrl.Scan(&followUpEntry.FollowUpEntryId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateFollowUpEntryByICode actualiza los datos de una entrada de seguimiento en la BD, utilizando su ICode como referencia.
// Recibe el objeto de entrada de seguimiento con los nuevos datos, información de transacción, módulo y datos de conexión.
// Retorna un error en caso de producirse.
func UpdateFollowUpEntryByICode(followUpEntry *FollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se actualizarán.
	var entryFieldsSlice []string = []string{
		"FollowUpEntryUpdateDate",
		"FollowUpEntryWasDone",
		"FollowUpEntrySector",
		"FollowUpEntryPersonVisitedEntity",
		"FollowUpEntryReceivedAttention",
		"FollowUpEntryComments",
		"FollowUpEntryCaseDocumentsPrepared",
		"FollowUpEntryStatus",
		"FollowUpEntryOwnerGeneralUser",
	}
	var entryFieldsAliasSlice []string = []string{}

	// Se genera la query de actualización utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, entryFieldsSlice, entryFieldsAliasSlice, FollowUpEntryDBName, []string{"FollowUpEntryICode"}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query,
		followUpEntry.FollowUpEntryICode,
		followUpEntry.FollowUpEntryUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		followUpEntry.FollowUpEntryWasDone,
		followUpEntry.FollowUpEntrySector,
		followUpEntry.FollowUpEntryPersonVisitedEntity,
		followUpEntry.FollowUpEntryReceivedAttention,
		followUpEntry.FollowUpEntryComments,
		followUpEntry.FollowUpEntryCaseDocumentsPrepared,
		followUpEntry.FollowUpEntryStatus,
		followUpEntry.FollowUpEntryOwnerGeneralUser,
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

// RemoveFollowUpEntryByICode elimina una entrada de seguimiento de la base de datos utilizando su ICode.
// Recibe el objeto de entrada de seguimiento, información de transacción, módulo y datos de conexión.
// Retorna un error en caso de producirse.
func RemoveFollowUpEntryByICode(followUpEntry *FollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos y condiciones para la eliminación.
	var entryFieldsSlice []string = []string{}
	var entryFieldsAliasSlice []string = []string{}

	// Se genera la query de eliminación utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, entryFieldsSlice, entryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{"FollowUpEntryICode"}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, false)

	// Se ejecuta la query con el parámetro correspondiente.
	persistenceCtrl.Exec(context.Background(), query, followUpEntry.FollowUpEntryICode)

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

// GetFollowUpEntry consulta una entrada de seguimiento basada en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto de entrada de seguimiento a completar,
// información de transacción, módulo y datos de conexión.
// Retorna un error en caso de producirse.
func GetFollowUpEntry(by common_controllers.By, followUpEntry *FollowUpEntryDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla FollowUpEntry en la BD.
	var entryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var entryFieldsSlice []string = []string{
		"FollowUpEntryId",
		"FollowUpEntryICode",
		"FollowUpEntryCreationDate",
		"FollowUpEntryUpdateDate",
		"FollowUpEntryCompletionDate",
		"FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone",
		"FollowUpEntrySector",
		"FollowUpEntryPersonVisitedEntity",
		"FollowUpEntryReceivedAttention",
		"FollowUpEntryComments",
		"FollowUpEntryCaseDocumentsPrepared",
		"FollowUpEntryStatus",
		"FollowUpEntryFollowUp",
	}
	var entryFieldsAliasSlice []string = []string{}

	var entryFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entryFieldsSlice, entryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + entryFieldsStr +
		` FROM ` + entryPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FollowUpEntryDBName, by.AttrsName, []string{}, []string{}, by.Operator, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	// Se ejecuta la query con los parámetros de filtrado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se preparan variables para capturar los datos.
	var followUpEntryPg = FollowUpEntryPgDB{}

	// Se escanean los resultados de la query.
	persistenceCtrl.Scan(
		&followUpEntryPg.FollowUpEntryId,
		&followUpEntryPg.FollowUpEntryICode,
		&followUpEntryPg.FollowUpEntryCreationDate,
		&followUpEntryPg.FollowUpEntryUpdateDate,
		&followUpEntryPg.FollowUpEntryCompletionDate,
		&followUpEntryPg.FollowUpEntryOwnerGeneralUser,
		&followUpEntryPg.FollowUpEntryWasDone,
		&followUpEntryPg.FollowUpEntrySector,
		&followUpEntryPg.FollowUpEntryPersonVisitedEntity,
		&followUpEntryPg.FollowUpEntryReceivedAttention,
		&followUpEntryPg.FollowUpEntryComments,
		&followUpEntryPg.FollowUpEntryCaseDocumentsPrepared,
		&followUpEntryPg.FollowUpEntryStatus,
		&followUpEntryPg.FollowUpEntryFollowUp,
	)

	// Se convierte el objeto de base de datos a DTO.
	*followUpEntry = followUpEntryPg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetFollowUpEntry consulta una entrada de seguimiento basada en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto de entrada de seguimiento a completar,
// información de transacción, módulo y datos de conexión.
// Retorna un error en caso de producirse.
func GetFollowUpEntries(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FollowUpEntryDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla FollowUpEntry en la BD.
	var entryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var entryFieldsSlice []string = []string{
		"FollowUpEntryId",
		"FollowUpEntryICode",
		"FollowUpEntryCreationDate",
		"FollowUpEntryUpdateDate",
		"FollowUpEntryCompletionDate",
		"FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone",
		"FollowUpEntrySector",
		"FollowUpEntryPersonVisitedEntity",
		"FollowUpEntryReceivedAttention",
		"FollowUpEntryComments",
		"FollowUpEntryCaseDocumentsPrepared",
		"FollowUpEntryStatus",
		"FollowUpEntryFollowUp",
	}
	var entryFieldsAliasSlice []string = []string{}

	var entryFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entryFieldsSlice, entryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + entryFieldsStr +
		` FROM ` + entryPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FollowUpEntryDBName, by.AttrsName, []string{}, []string{}, by.Operator, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true) +
		` ORDER BY ` + entryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryCompletionDate"].DBName + ` ASC `

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	var entries []FollowUpEntryDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var entryPg = FollowUpEntryPgDB{}

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(
			&entryPg.FollowUpEntryId,
			&entryPg.FollowUpEntryICode,
			&entryPg.FollowUpEntryCreationDate,
			&entryPg.FollowUpEntryUpdateDate,
			&entryPg.FollowUpEntryCompletionDate,
			&entryPg.FollowUpEntryOwnerGeneralUser,
			&entryPg.FollowUpEntryWasDone,
			&entryPg.FollowUpEntrySector,
			&entryPg.FollowUpEntryPersonVisitedEntity,
			&entryPg.FollowUpEntryReceivedAttention,
			&entryPg.FollowUpEntryComments,
			&entryPg.FollowUpEntryCaseDocumentsPrepared,
			&entryPg.FollowUpEntryStatus,
			&entryPg.FollowUpEntryFollowUp,
		)

		var entry FollowUpEntryDTO = entryPg.ToDTO()
		entries = append(entries, entry)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entries, nil
}

// GetAllFollowUpEntries consulta todas las entradas de seguimiento existentes en la base de datos.
// Recibe información de transacción, módulo y datos de conexión.
// Retorna un slice de FollowUpEntryDTO y un error en caso de producirse.
func GetAllFollowUpEntries(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FollowUpEntryDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var entryPath string = FollowUpEntryDBScheme + "." + FollowUpEntryDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var entryFieldsSlice []string = []string{
		"FollowUpEntryId",
		"FollowUpEntryICode",
		"FollowUpEntryCreationDate",
		"FollowUpEntryUpdateDate",
		"FollowUpEntryCompletionDate",
		"FollowUpEntryOwnerGeneralUser",
		"FollowUpEntryWasDone",
		"FollowUpEntrySector",
		"FollowUpEntryPersonVisitedEntity",
		"FollowUpEntryReceivedAttention",
		"FollowUpEntryComments",
		"FollowUpEntryCaseDocumentsPrepared",
		"FollowUpEntryStatus",
		"FollowUpEntryFollowUp",
	}
	var entryFieldsAliasSlice []string = []string{}

	var entryFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, entryFieldsSlice, entryFieldsAliasSlice, FollowUpEntryDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryDBScheme, FollowUpEntryFieldDefinitions, true)

	// Se genera la query para obtener todos las entradas de seguimiento.
	var query string = `SELECT ` + entryFieldsStr +
		` FROM ` + entryPath +
		` WHERE TRUE` +
		` ORDER BY ` + entryPath + `.` + FollowUpEntryFieldDefinitions["FollowUpEntryCompletionDate"].DBName + ` ASC `

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query)
	var entries []FollowUpEntryDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var entryPg = FollowUpEntryPgDB{}

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(
			&entryPg.FollowUpEntryId,
			&entryPg.FollowUpEntryICode,
			&entryPg.FollowUpEntryCreationDate,
			&entryPg.FollowUpEntryUpdateDate,
			&entryPg.FollowUpEntryCompletionDate,
			&entryPg.FollowUpEntryOwnerGeneralUser,
			&entryPg.FollowUpEntryWasDone,
			&entryPg.FollowUpEntrySector,
			&entryPg.FollowUpEntryPersonVisitedEntity,
			&entryPg.FollowUpEntryReceivedAttention,
			&entryPg.FollowUpEntryComments,
			&entryPg.FollowUpEntryCaseDocumentsPrepared,
			&entryPg.FollowUpEntryStatus,
			&entryPg.FollowUpEntryFollowUp,
		)

		var entry FollowUpEntryDTO = entryPg.ToDTO()
		entries = append(entries, entry)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return entries, nil
}

// SetFollowUpEntryDefaults asigna valores por defecto a los campos de una entrada de seguimiento,
// dependiendo de la acción que se esté realizando (insertar o actualizar).
func SetFollowUpEntryDefaults(followUpEntry *FollowUpEntryDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción se asigna estado 'p', fechas actuales y se genera un UUID para el ICode.
		followUpEntry.FollowUpEntryStatus = "u"
		followUpEntry.FollowUpEntryWasDone = "n"
		followUpEntry.FollowUpEntryCreationDate = time.Now()
		followUpEntry.FollowUpEntryUpdateDate = time.Now()
		followUpEntry.FollowUpEntryICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para actualización solo se actualiza la fecha de modificación.
		followUpEntry.FollowUpEntryUpdateDate = time.Now()
	}
}

// ToDTO convierte un objeto FollowUpEntryPgDB obtenido de la base de datos en un objeto FollowUpEntryDTO.
// Se encarga de verificar la validez de los campos nulos y asignarlos correctamente.
func (obj *FollowUpEntryPgDB) ToDTO() FollowUpEntryDTO {
	var dto FollowUpEntryDTO

	if obj.FollowUpEntryId.Valid {
		dto.FollowUpEntryId = uint64(obj.FollowUpEntryId.Int64)
	}

	if obj.FollowUpEntryICode.Valid {
		dto.FollowUpEntryICode = obj.FollowUpEntryICode.String
	}

	if obj.FollowUpEntryCreationDate.Valid {
		dto.FollowUpEntryCreationDate = obj.FollowUpEntryCreationDate.Time
	}

	if obj.FollowUpEntryUpdateDate.Valid {
		dto.FollowUpEntryUpdateDate = obj.FollowUpEntryUpdateDate.Time
	}

	if obj.FollowUpEntryCompletionDate.Valid {
		dto.FollowUpEntryCompletionDate = obj.FollowUpEntryCompletionDate.Time
	}

	if obj.FollowUpEntryOwnerGeneralUser.Valid {
		dto.FollowUpEntryOwnerGeneralUser = obj.FollowUpEntryOwnerGeneralUser.String
	}

	if obj.FollowUpEntryWasDone.Valid {
		dto.FollowUpEntryWasDone = obj.FollowUpEntryWasDone.String
	}

	if obj.FollowUpEntrySector.Valid {
		dto.FollowUpEntrySector = obj.FollowUpEntrySector.String
	}

	if obj.FollowUpEntryPersonVisitedEntity.Valid {
		dto.FollowUpEntryPersonVisitedEntity = obj.FollowUpEntryPersonVisitedEntity.String
	}

	if obj.FollowUpEntryReceivedAttention.Valid {
		dto.FollowUpEntryReceivedAttention = obj.FollowUpEntryReceivedAttention.String
	}

	if obj.FollowUpEntryComments.Valid {
		dto.FollowUpEntryComments = obj.FollowUpEntryComments.String
	}

	if obj.FollowUpEntryCaseDocumentsPrepared.Valid {
		dto.FollowUpEntryCaseDocumentsPrepared = obj.FollowUpEntryCaseDocumentsPrepared.String
	}

	if obj.FollowUpEntryStatus.Valid {
		dto.FollowUpEntryStatus = obj.FollowUpEntryStatus.String
	}

	if obj.FollowUpEntryFollowUp.Valid {
		dto.FollowUpEntryFollowUp = FollowUpDTO{FollowUpId: uint64(obj.FollowUpEntryFollowUp.Int64)}
	}

	return dto
}
