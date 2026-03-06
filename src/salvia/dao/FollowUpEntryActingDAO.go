// Package salvia_daos contiene los Objetos de Acceso a Datos (DAO) para las entidades
// de la lógica de negocio de "salvia". Este archivo gestiona la persistencia de la entidad
// 'FollowUpEntryActing', que representa una acción o actuación específica asociada a una
// entrada de seguimiento (FollowUpEntry).
package salvia_daos

import (
	// Importaciones de paquetes del proyecto y de la librería estándar de Go.
	common_config "bitsflow/common/config"           // Proporciona acceso a configuraciones globales, como formatos de fecha.
	common_controllers "bitsflow/common/controllers" // Contiene el controlador de persistencia y otras utilidades.
	common_dao "bitsflow/common/dao"                 // Helpers y constantes para la capa de acceso a datos.
	"bitsflow/common/db"                             // Manejo de la conexión y configuración de la base de datos.
	"bitsflow/common/utils"                          // Funciones de utilidad como la generación de UUIDs.
	security_daos "bitsflow/security/dao"
	"context"      // Para el manejo de contextos en operaciones de base de datos.
	"database/sql" // Interfaz genérica para bases de datos SQL.
	"encoding/json"
	"errors"
	"fmt"  // Paquete para formateo de I/O, usado para imprimir errores.
	"time" // Para trabajar con fechas y horas.
)

// Este bloque var define los metadatos para la entidad FollowUpEntryActing.
// Centralizar esta información facilita el mantenimiento y asegura la consistencia
// al interactuar con la base de datos y al serializar datos.
var (
	// FollowUpEntryActingEntityName es el nombre de la entidad en la lógica de la aplicación.
	FollowUpEntryActingEntityName string = "FollowUpEntryActing"
	// FollowUpEntryActingJSONName es el nombre que se usará para la serialización JSON.
	FollowUpEntryActingJSONName string = "followUpEntryActing"
	// FollowUpEntryActingDBName es el nombre de la tabla en la base de datos.
	FollowUpEntryActingDBName string = "follow_up_entry_acting"
	// FollowUpEntryActingDBScheme es el esquema de la base de datos donde reside la tabla.
	FollowUpEntryActingDBScheme string = "salvia"

	// FollowUpEntryActingFieldDefinitions mapea los campos del struct a las columnas de la base de datos
	// y define sus propiedades. Es utilizado por helpers para construir consultas SQL dinámicas.
	FollowUpEntryActingFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"FollowUpEntryActingId":               {Name: "FollowUpEntryActingId", DBName: "follow_up_entry_acting_id", Alias: "", ModelType: "uint", Required: false},
		"FollowUpEntryActingICode":            {Name: "FollowUpEntryActingICode", DBName: "follow_up_entry_acting_i_code", Alias: "", ModelType: "string", MaxSize: 36, Required: false},
		"FollowUpEntryActingCreationDate":     {Name: "FollowUpEntryActingCreationDate", DBName: "follow_up_entry_acting_creation_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpEntryActingUpdateDate":       {Name: "FollowUpEntryActingUpdateDate", DBName: "follow_up_entry_acting_update_date", Alias: "", ModelType: "datetime", Required: false},
		"FollowUpEntryActingOwnerGeneralUser": {Name: "FollowUpEntryActingOwnerGeneralUser", DBName: "follow_up_entry_acting_owner_general_user", Alias: "", ModelType: "string", MaxSize: 36, Required: true},
		"FollowUpEntryActingDescription":      {Name: "FollowUpEntryActingDescription", DBName: "follow_up_entry_acting_description", Alias: "", ModelType: "string", MaxSize: 1000, Required: true},
		"FollowUpEntryActionStatus":           {Name: "FollowUpEntryActionStatus", DBName: "follow_up_entry_acting_status", Alias: "", ModelType: "string", MaxSize: 1, Required: true},
		// Este campo representa la clave foránea a la tabla 'follow_up_entry'.
		"FollowUpEntryActingRelBarrierFollowUpEntry": {Name: "FollowUpEntryActingRelBarrierFollowUpEntry", DBName: "follow_up_entry_acting_rel_barrier_follow_up_entry", Alias: "", ModelType: "uint", Required: false},
	}
)

// FollowUpEntryActingDTO (Data Transfer Object) representa una actuación o acción
// dentro de una entrada de seguimiento en la capa de aplicación.
type FollowUpEntryActingDTO struct {
	FollowUpEntryActingId               uint64    `json:"-"`
	FollowUpEntryActingICode            string    `json:"icode"`
	FollowUpEntryActingCreationDate     time.Time `json:"creationDate"`
	FollowUpEntryActingUpdateDate       time.Time `json:"updateDate"`
	FollowUpEntryActingOwnerGeneralUser string    `json:"ownerGeneralUser"`
	FollowUpEntryActingDescription      string    `json:"description"`
	FollowUpEntryActionStatus           string    `json:"actionStatus"`

	// Permite anidar la información de la entrada de seguimiento a la que pertenece esta actuación.
	FollowUpEntryActingRelBarrierFollowUpEntry RelBarrierFollowUpEntryDTO   `json:"relBarrierFollowUpEntry"`
	FollowUpEntryActingOwnerGeneralUserDTO     security_daos.GeneralUserDTO `json:"owner"`
}

// FollowUpEntryActingPgDB representa la estructura de la tabla 'follow_up_entry_acting'
// en la base de datos. Utiliza tipos `sql.Null*` para manejar correctamente columnas
// que pueden ser nulas.
type FollowUpEntryActingPgDB struct {
	FollowUpEntryActingId               sql.NullInt64
	FollowUpEntryActingICode            sql.NullString
	FollowUpEntryActingCreationDate     sql.NullTime
	FollowUpEntryActingUpdateDate       sql.NullTime
	FollowUpEntryActingOwnerGeneralUser sql.NullString
	FollowUpEntryActingDescription      sql.NullString
	FollowUpEntryActionStatus           sql.NullString
	// Clave foránea que referencia a 'follow_up_entry'.
	FollowUpEntryActingRelBarrierFollowUpEntry sql.NullInt64
}

func (fuea FollowUpEntryActingDTO) MarshalJSON() ([]byte, error) {
	type Alias FollowUpEntryActingDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		FollowUpEntryActingCreationDate string `json:"creationDate"`
		FollowUpEntryActingUpdateDate   string `json:"updateDate"`
	}{
		Alias:                           (*Alias)(&fuea),
		FollowUpEntryActingCreationDate: fuea.FollowUpEntryActingCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		FollowUpEntryActingUpdateDate:   fuea.FollowUpEntryActingUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (fuea *FollowUpEntryActingDTO) UnmarshalJSON(data []byte) error {
	type Alias FollowUpEntryActingDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		FollowUpEntryActingCreationDate string `json:"creationDate"`
		FollowUpEntryActingUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(fuea),
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
	fuea.FollowUpEntryActingCreationDate = parse(aux.FollowUpEntryActingCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	fuea.FollowUpEntryActingUpdateDate = parse(aux.FollowUpEntryActingUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetFollowUpEntryActing inserta un nuevo registro de actuación de seguimiento en la base de datos.
// Recibe el objeto de actuación, información de conexión y configuración de la base de datos.
// Retorna un error si la operación falla.
func SetFollowUpEntryActing(acting *FollowUpEntryActingDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se define la lista de campos que se insertarán en la BD.
	var actingFieldsSlice []string = []string{"FollowUpEntryActingICode", "FollowUpEntryActingCreationDate", "FollowUpEntryActingUpdateDate", "FollowUpEntryActingOwnerGeneralUser", "FollowUpEntryActingDescription", "FollowUpEntryActionStatus", "FollowUpEntryActingRelBarrierFollowUpEntry"}
	var actingFieldsAliasSlice []string = []string{}

	// Se genera la query de inserción utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_INSERT, actingFieldsSlice, actingFieldsAliasSlice, FollowUpEntryActingDBName, []string{}, []string{}, []string{"FollowUpEntryActingId"}, common_dao.SQL_AND, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.QueryRow(context.Background(), query,
		acting.FollowUpEntryActingICode,
		acting.FollowUpEntryActingCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		acting.FollowUpEntryActingUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		acting.FollowUpEntryActingOwnerGeneralUser,
		acting.FollowUpEntryActingDescription,
		acting.FollowUpEntryActionStatus,
		acting.FollowUpEntryActingRelBarrierFollowUpEntry.RelBarrierFollowUpEntryId)
	persistenceCtrl.Scan(&acting.FollowUpEntryActingId)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateFollowUpEntryActingByICode actualiza los datos de una actuación de seguimiento en la BD, utilizando su ICode como referencia.
// Recibe el objeto actuación con los nuevos datos, información de conexión y configuración de la base de datos.
// Retorna un error si la operación falla.
func UpdateFollowUpEntryActingByICode(acting *FollowUpEntryActingDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos que se actualizarán.
	var actingFieldsSlice []string = []string{"FollowUpEntryActingUpdateDate", "FollowUpEntryActingDescription", "FollowUpEntryActionStatus"}
	var actingFieldsAliasSlice []string = []string{}

	// Se genera la query de actualización utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_UPDATE, actingFieldsSlice, actingFieldsAliasSlice, FollowUpEntryActingDBName, []string{"FollowUpEntryActingICode"}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, false)

	// Se ejecuta la query con los parámetros correspondientes.
	persistenceCtrl.Exec(context.Background(), query, acting.FollowUpEntryActingICode,
		acting.FollowUpEntryActingUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		acting.FollowUpEntryActingDescription,
		acting.FollowUpEntryActionStatus)

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

// RemoveFollowUpEntryActingByICode elimina una actuación de seguimiento de la base de datos utilizando su ICode.
// Recibe el objeto actuación, información de conexión y configuración de la base de datos.
// Retorna un error si la operación falla.
func RemoveFollowUpEntryActingByICode(acting *FollowUpEntryActingDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos y condiciones para la eliminación.
	var actingFieldsSlice []string = []string{}
	var actingFieldsAliasSlice []string = []string{}

	// Se genera la query de eliminación utilizando la función GetSQL.
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, actingFieldsSlice, actingFieldsAliasSlice, FollowUpEntryActingDBName, []string{}, []string{"FollowUpEntryActingICode"}, []string{}, common_dao.SQL_AND, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, false)

	// Se ejecuta la query con el parámetro correspondiente.
	persistenceCtrl.Exec(context.Background(), query, acting.FollowUpEntryActingICode)

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

// GetFollowUpEntryActing consulta una actuación de seguimiento basada en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto actuación a completar,
// información de conexión y configuración de la base de datos.
// Retorna un error si la operación falla.
func GetFollowUpEntryActing(by common_controllers.By, acting *FollowUpEntryActingDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla FollowUpEntryActing en la BD.
	var actingPath string = FollowUpEntryActingDBScheme + "." + FollowUpEntryActingDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var actingFieldsSlice []string = []string{"FollowUpEntryActingId", "FollowUpEntryActingICode", "FollowUpEntryActingCreationDate", "FollowUpEntryActingUpdateDate", "FollowUpEntryActingOwnerGeneralUser", "FollowUpEntryActingDescription", "FollowUpEntryActionStatus", "FollowUpEntryActingRelBarrierFollowUpEntry"}
	var actingFieldsAliasSlice []string = []string{}
	var actingFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, actingFieldsSlice, actingFieldsAliasSlice, FollowUpEntryActingDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + actingFieldsStr +
		` FROM ` + actingPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FollowUpEntryActingDBName, by.AttrsName, []string{}, []string{}, by.Operator, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, true)

	// Se ejecuta la query con los parámetros de filtrado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	var actingPg = FollowUpEntryActingPgDB{}

	// Se escanean los resultados de la query.
	persistenceCtrl.Scan(
		&actingPg.FollowUpEntryActingId,
		&actingPg.FollowUpEntryActingICode,
		&actingPg.FollowUpEntryActingCreationDate,
		&actingPg.FollowUpEntryActingUpdateDate,
		&actingPg.FollowUpEntryActingOwnerGeneralUser,
		&actingPg.FollowUpEntryActingDescription,
		&actingPg.FollowUpEntryActionStatus,
		&actingPg.FollowUpEntryActingRelBarrierFollowUpEntry,
	)

	// Se convierte el objeto de base de datos a DTO.
	*acting = actingPg.ToDTO()

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetFollowUpEntryActing consulta las actuaciones de seguimiento basada en atributos de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de filtrado, el objeto actuación a completar,
// información de conexión y configuración de la base de datos.
// Retorna las actuaciones o un error si la operación falla.
func GetFollowUpEntriesActing(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FollowUpEntryActingDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path de la tabla FollowUpEntryActing en la BD.
	var actingPath string = FollowUpEntryActingDBScheme + "." + FollowUpEntryActingDBName
	var userPath string = security_daos.GeneralUserDBScheme + "." + security_daos.GeneralUserDBName
	var profilePath string = security_daos.GeneralUserProfileDBScheme + "." + security_daos.GeneralUserProfileDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar.
	var actingFieldsSlice []string = []string{"FollowUpEntryActingId", "FollowUpEntryActingICode", "FollowUpEntryActingCreationDate", "FollowUpEntryActingUpdateDate", "FollowUpEntryActingOwnerGeneralUser", "FollowUpEntryActingDescription", "FollowUpEntryActionStatus", "FollowUpEntryActingRelBarrierFollowUpEntry"}
	var actingFieldsAliasSlice []string = []string{}
	var actingFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, actingFieldsSlice, actingFieldsAliasSlice, FollowUpEntryActingDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, true)

	var usrFieldsSlice []string = []string{"GeneralUserId", "GeneralUserICode", "GeneralUserCreationDate", "GeneralUserUpdateDate", "GeneralUserLogin", "GeneralUserStatus", "GeneralUserLanguage"}
	var usrFieldsAliasSlice []string = []string{}

	var profileFieldsSlice []string = []string{"GeneralUserProfileId", "GeneralUserProfileICode", "GeneralUserProfileCreationDate", "GeneralUserProfileUpdateDate", "GeneralUserProfileGender", "GeneralUserProfileNick", "GeneralUserProfileDescription", "GeneralUserProfileNames", "GeneralUserProfileLastNames", "GeneralUserProfileDocType", "GeneralUserProfileDocNumber"}
	var profileFieldsAliasSlice []string = []string{}

	// Se generan los strings de campos.
	var usrFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, usrFieldsSlice, usrFieldsAliasSlice, security_daos.GeneralUserDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.GeneralUserDBScheme, security_daos.GeneralUserFieldDefinitions, true)
	var profileFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, profileFieldsSlice, profileFieldsAliasSlice, security_daos.GeneralUserProfileDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, security_daos.GeneralUserProfileDBScheme, security_daos.GeneralUserProfileFieldDefinitions, true)

	// Se genera la query de selección utilizando la función GetSQL y añadiendo condiciones dinámicas.
	var query string = `SELECT ` + actingFieldsStr + ", " + usrFieldsStr + ", " + profileFieldsStr +
		` FROM ` + actingPath +
		` LEFT JOIN ` + userPath + ` ON (` + actingPath + `.` + FollowUpEntryActingFieldDefinitions["FollowUpEntryActingOwnerGeneralUser"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserICode"].DBName + `)` +
		` LEFT JOIN ` + profilePath + ` ON (` + profilePath + `.` + security_daos.GeneralUserProfileFieldDefinitions["GeneralUserProfileId"].DBName + ` = ` + userPath + `.` + security_daos.GeneralUserFieldDefinitions["GeneralUserGeneralUserProfile"].DBName + `)` +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, FollowUpEntryActingDBName, by.AttrsName, []string{}, []string{}, by.Operator, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, true)

	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	var actings []FollowUpEntryActingDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var actingPg = FollowUpEntryActingPgDB{}
		var profilePg security_daos.GeneralUserProfilePgDB = security_daos.GeneralUserProfilePgDB{}
		var userPg = security_daos.GeneralUserPgDB{}

		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(
			&actingPg.FollowUpEntryActingId,
			&actingPg.FollowUpEntryActingICode,
			&actingPg.FollowUpEntryActingCreationDate,
			&actingPg.FollowUpEntryActingUpdateDate,
			&actingPg.FollowUpEntryActingOwnerGeneralUser,
			&actingPg.FollowUpEntryActingDescription,
			&actingPg.FollowUpEntryActionStatus,
			&actingPg.FollowUpEntryActingRelBarrierFollowUpEntry,
			&userPg.GeneralUserId, &userPg.GeneralUserICode, &userPg.GeneralUserCreationDate, &userPg.GeneralUserUpdateDate, &userPg.GeneralUserLogin,
			&userPg.GeneralUserStatus, &userPg.GeneralUserLanguage,
			&profilePg.GeneralUserProfileId, &profilePg.GeneralUserProfileICode, &profilePg.GeneralUserProfileCreationDate, &profilePg.GeneralUserProfileUpdateDate, &profilePg.GeneralUserProfileGender, &profilePg.GeneralUserProfileNick, &profilePg.GeneralUserProfileDescription, &profilePg.GeneralUserProfileNames, &profilePg.GeneralUserProfileLastNames, &profilePg.GeneralUserProfileDocType, &profilePg.GeneralUserProfileDocNumber,
		)
		var acting FollowUpEntryActingDTO = actingPg.ToDTO()
		var user security_daos.GeneralUserDTO = userPg.ToDTO()

		user.GeneralUserGeneralUserProfile = profilePg.ToDTO()
		acting.FollowUpEntryActingOwnerGeneralUserDTO = user
		actings = append(actings, acting)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return actings, nil
}

// GetAllFollowUpEntryActing consulta todos los registros de actuaciones de seguimiento existentes en la base de datos.
// Recibe información de conexión y configuración de la base de datos.
// Retorna un slice de FollowUpEntryActingDTO y un error en caso de producirse.
func GetAllFollowUpEntriesActing(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]FollowUpEntryActingDTO, error) {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var actingPath string = FollowUpEntryActingDBScheme + "." + FollowUpEntryActingDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Se definen los campos a consultar para las actuaciones.
	var actingFieldsSlice []string = []string{"FollowUpEntryActingId", "FollowUpEntryActingICode", "FollowUpEntryActingCreationDate", "FollowUpEntryActingUpdateDate", "FollowUpEntryActingOwnerGeneralUser", "FollowUpEntryActingDescription", "FollowUpEntryActionStatus", "FollowUpEntryActingRelBarrierFollowUpEntry"}
	var actingFieldsAliasSlice []string = []string{}
	var actingFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, actingFieldsSlice, actingFieldsAliasSlice, FollowUpEntryActingDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, FollowUpEntryActingDBScheme, FollowUpEntryActingFieldDefinitions, true)

	// Se genera la query para obtener todas las actuaciones.
	var query string = `SELECT ` + actingFieldsStr +
		` FROM ` + actingPath +
		` WHERE TRUE`

	// Se ejecuta la query.
	persistenceCtrl.Query(context.Background(), query)
	var actings []FollowUpEntryActingDTO
	// Se itera sobre cada registro obtenido.
	for persistenceCtrl.Next() {
		var actingPg = FollowUpEntryActingPgDB{}
		var entryPg = FollowUpEntryPgDB{}
		// Se escanean los datos de cada fila.
		persistenceCtrl.ScanRow(
			&actingPg.FollowUpEntryActingId,
			&actingPg.FollowUpEntryActingICode,
			&actingPg.FollowUpEntryActingCreationDate,
			&actingPg.FollowUpEntryActingUpdateDate,
			&actingPg.FollowUpEntryActingOwnerGeneralUser,
			&actingPg.FollowUpEntryActingDescription,
			&actingPg.FollowUpEntryActionStatus,
			&actingPg.FollowUpEntryActingRelBarrierFollowUpEntry,
			&entryPg.FollowUpEntryId,
		)
		var acting FollowUpEntryActingDTO = actingPg.ToDTO()
		actings = append(actings, acting)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return actings, nil
}

// SetFollowUpEntryActingDefaults asigna valores por defecto a los campos de una actuación de seguimiento,
// dependiendo de la acción que se esté realizando (insertar o actualizar).
func SetFollowUpEntryActingDefaults(acting *FollowUpEntryActingDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para inserción se asignan fechas actuales y se genera un UUID para el ICode.
		acting.FollowUpEntryActionStatus = "p"
		acting.FollowUpEntryActingCreationDate = time.Now()
		acting.FollowUpEntryActingUpdateDate = time.Now()
		acting.FollowUpEntryActingICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// Para actualización solo se actualiza la fecha de modificación.
		acting.FollowUpEntryActingUpdateDate = time.Now()
	}
}

// ToDTO convierte un objeto FollowUpEntryActingPgDB obtenido de la base de datos en un objeto FollowUpEntryActingDTO.
// Se encarga de verificar la validez de los campos nulos y asignarlos correctamente.
func (obj *FollowUpEntryActingPgDB) ToDTO() FollowUpEntryActingDTO {
	var dto FollowUpEntryActingDTO

	if obj.FollowUpEntryActingId.Valid {
		dto.FollowUpEntryActingId = uint64(obj.FollowUpEntryActingId.Int64)
	}

	if obj.FollowUpEntryActingICode.Valid {
		dto.FollowUpEntryActingICode = obj.FollowUpEntryActingICode.String
	}

	if obj.FollowUpEntryActingCreationDate.Valid {
		dto.FollowUpEntryActingCreationDate = obj.FollowUpEntryActingCreationDate.Time
	}

	if obj.FollowUpEntryActingUpdateDate.Valid {
		dto.FollowUpEntryActingUpdateDate = obj.FollowUpEntryActingUpdateDate.Time
	}

	if obj.FollowUpEntryActingOwnerGeneralUser.Valid {
		dto.FollowUpEntryActingOwnerGeneralUser = obj.FollowUpEntryActingOwnerGeneralUser.String
	}

	if obj.FollowUpEntryActingDescription.Valid {
		dto.FollowUpEntryActingDescription = obj.FollowUpEntryActingDescription.String
	}

	if obj.FollowUpEntryActionStatus.Valid {
		dto.FollowUpEntryActionStatus = obj.FollowUpEntryActionStatus.String
	}

	if obj.FollowUpEntryActingRelBarrierFollowUpEntry.Valid {
		dto.FollowUpEntryActingRelBarrierFollowUpEntry = RelBarrierFollowUpEntryDTO{RelBarrierFollowUpEntryId: uint64(obj.FollowUpEntryActingRelBarrierFollowUpEntry.Int64)}
	}

	return dto
}
