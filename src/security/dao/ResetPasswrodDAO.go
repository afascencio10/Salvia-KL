// Package security_daos contiene las funciones de acceso a datos (DAO) relacionadas con la entidad ResetPassword.
// Este paquete se encarga de realizar operaciones CRUD sobre la tabla de reset de contraseñas en la base de datos.
package security_daos

import (
	// Configuración general de la aplicación.
	// Controladores comunes, incluyendo el controlador de persistencia.
	// Funciones y constantes para construir queries SQL.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"    // Conexión y configuraciones de base de datos.
	"bitsflow/common/utils" // Utilidades generales, como manejo de UUID y definiciones de campos.
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	// Nombres y esquemas utilizados para la entidad ResetPassword.
	ResetPasswordEntityName string = "ResetPassword"
	ResetPasswordJSONName   string = "resetPassword"
	ResetPasswordDBName     string = "reset_password"
	ResetPasswordDBScheme   string = "security"

	// Definiciones de los campos para la validación de datos provenientes del JSON.
	// Cada campo se define mediante un objeto FieldDefinition que indica el nombre, nombre en BD, tipo de dato, tamaños mínimos y máximos, y si es requerido.
	ResetPasswordFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"ResetPasswordId":           {Name: "ResetPasswordId", DBName: "reset_password_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"ResetPasswordICode":        {Name: "ResetPasswordICode", DBName: "reset_password_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"ResetPasswordCreationDate": {Name: "ResetPasswordCreationDate", DBName: "reset_password_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"ResetPasswordGeneralUser":  {Name: "ResetPasswordGeneralUser", DBName: "reset_password_general_user", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
	}
)

// ResetPasswordDTO representa la estructura de datos utilizada en la capa de negocio para la entidad ResetPassword.
// GeneralUserDTO se asume definido en otro paquete y representa al usuario asociado.
type ResetPasswordDTO struct {
	ResetPasswordId           uint64         `json:"-"`
	ResetPasswordICode        string         `json:"-"`
	ResetPasswordCreationDate time.Time      `json:"-"`
	ResetPasswordGeneralUser  GeneralUserDTO `json:"-"`
}

// ResetPasswordPgDB representa la estructura que mapea los datos de la base de datos.
// Se utilizan tipos sql.Null* para manejar valores nulos provenientes de la BD.
type ResetPasswordPgDB struct {
	ResetPasswordId           sql.NullInt64
	ResetPasswordICode        sql.NullString
	ResetPasswordCreationDate sql.NullTime
	ResetPasswordGeneralUser  sql.NullInt64
}

// SetResetPassword inserta un nuevo registro de reset de contraseña en la base de datos.
// Parámetros:
//   - resetPassword: Puntero al DTO con los datos a insertar.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo o contexto que invoca la función.
//   - connData: Datos de conexión actuales.
//   - clientConfig: Configuración del cliente de la BD.
//   - serverConfig: Configuración del servidor de la BD.
//
// Devuelve la conexión actualizada y un error en caso de producirse.
func SetResetPassword(resetPassword *ResetPasswordDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia para manejar la conexión y ejecución de queries.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión con la BD utilizando la configuración proporcionada.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Definición de los campos a insertar.
	var resetPasswordFieldsSlice []string = []string{"ResetPasswordICode", "ResetPasswordCreationDate", "ResetPasswordGeneralUser"}
	var resetPasswordFieldsAliasSlice []string = []string{}

	// Se construye el query SQL para insertar un nuevo registro.
	var query string = common_dao.GetSQL(
		common_dao.SQL_INSERT,
		resetPasswordFieldsSlice,
		resetPasswordFieldsAliasSlice,
		ResetPasswordDBName,
		[]string{},
		[]string{},
		[]string{"ResetPasswordId"},
		common_dao.SQL_AND,
		ResetPasswordDBScheme,
		ResetPasswordFieldDefinitions,
		false,
	)

	// Se ejecuta el query y se pasan los valores a insertar.
	persistenceCtrl.QueryRow(context.Background(), query,
		resetPassword.ResetPasswordICode,
		resetPassword.ResetPasswordCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		resetPassword.ResetPasswordGeneralUser.GeneralUserId,
	)

	// Se escanea el valor retornado (por ejemplo, el ID generado) y se asigna al DTO.
	persistenceCtrl.Scan(&resetPassword.ResetPasswordId)

	// En caso de error, se imprime el query y el error y se retorna.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetResetPassword recupera un registro de reset de contraseña basado en condiciones específicas.
// Parámetros:
//   - by: Estructura que define los atributos y valores para la cláusula WHERE.
//   - resetPassword: Puntero al DTO donde se almacenará el registro recuperado.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo o contexto que invoca la función.
//   - connData: Datos de conexión actuales.
//   - clientConfig: Configuración del cliente de la BD.
//   - serverConfig: Configuración del servidor de la BD.
//
// Devuelve la conexión actualizada y un error en caso de producirse.
func GetResetPassword(by common_controllers.By, resetPassword *ResetPasswordDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path completo de la tabla (esquema.nombre_tabla).
	var resetPasswordPath string = ResetPasswordDBScheme + "." + ResetPasswordDBName

	// Se establece la conexión con la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var resetPasswordFieldsSlice []string = []string{"ResetPasswordId", "ResetPasswordICode", "ResetPasswordCreationDate", "ResetPasswordGeneralUser"}
	var resetPasswordFieldsAliasSlice []string = []string{}

	// Se construye la parte SELECT de la consulta SQL.
	var resetPasswordFieldsStr = common_dao.GetSQL(
		common_dao.SQL_SELECT_FIELDS_ONLY,
		resetPasswordFieldsSlice,
		resetPasswordFieldsAliasSlice,
		ResetPasswordDBName,
		[]string{},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		ResetPasswordDBScheme,
		ResetPasswordFieldDefinitions,
		true,
	)

	// Se construye la consulta completa, incluyendo la cláusula WHERE basada en los atributos proporcionados.
	var query string = `SELECT ` + resetPasswordFieldsStr +
		` FROM ` + resetPasswordPath +
		common_dao.GetSQL(
			common_dao.SQL_SELECT_WHERE_ONLY,
			by.AttrsName,
			by.AttrsAliasName,
			ResetPasswordDBName,
			by.AttrsName,
			[]string{},
			[]string{},
			by.Operator,
			ResetPasswordDBScheme,
			ResetPasswordFieldDefinitions,
			true,
		)

	// Se ejecuta el query pasando los valores para la cláusula WHERE.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Se escanean los resultados y se asignan a los campos del DTO.
	persistenceCtrl.Scan(
		&resetPassword.ResetPasswordId,
		&resetPassword.ResetPasswordICode,
		&resetPassword.ResetPasswordCreationDate,
		&resetPassword.ResetPasswordGeneralUser.GeneralUserId,
	)

	// Se comprueba si hubo algún error en la ejecución del query.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// RemoveResetPasswordByICode elimina un registro de reset de contraseña en la BD utilizando el ICode único.
// Parámetros:
//   - resetPasswordIcode: Valor del ICode que identifica el registro a eliminar.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo o contexto que invoca la función.
//   - connData: Datos de conexión actuales.
//   - clientConfig: Configuración del cliente de la BD.
//   - serverConfig: Configuración del servidor de la BD.
//
// Devuelve la conexión actualizada y un error en caso de producirse.
func RemoveResetPasswordByICode(resetPasswordIcode string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión con la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// No se requieren campos adicionales para el DELETE.
	var usrFieldsSlice []string = []string{}
	var usrFieldsAliasSlice []string = []string{}

	// Se construye el query SQL para eliminar el registro basándose en el ResetPasswordICode.
	var query string = common_dao.GetSQL(
		common_dao.SQL_DELETE,
		usrFieldsSlice,
		usrFieldsAliasSlice,
		ResetPasswordDBName,
		[]string{"ResetPasswordICode"},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		ResetPasswordDBScheme,
		ResetPasswordFieldDefinitions,
		false,
	)

	// Se ejecuta la consulta pasando el ICode.
	persistenceCtrl.Exec(context.Background(), query, resetPasswordIcode)

	// Se comprueba si hubo error durante la ejecución.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Si no se afectó ninguna fila, se retorna un error.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveResetPasswordByTime elimina registros de reset de contraseña que sean anteriores a un tiempo determinado.
// El parámetro timeout indica el tiempo en minutos.
// Parámetros:
//   - timeout: Tiempo en minutos para determinar la antigüedad del registro.
//   - inTransaction: Indica si la operación se ejecuta dentro de una transacción.
//   - module: Nombre del módulo o contexto que invoca la función.
//   - connData: Datos de conexión actuales.
//   - clientConfig: Configuración del cliente de la BD.
//   - serverConfig: Configuración del servidor de la BD.
//
// Devuelve la conexión actualizada y un error en caso de producirse.
func RemoveResetPasswordByTime(timeout uint, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Se instancia el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Se construye el path completo de la tabla.
	var resetPasswordPath string = ResetPasswordDBScheme + "." + ResetPasswordDBName

	// Se establece la conexión con la BD.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Se construye el query SQL para eliminar los registros cuyo campo ResetPasswordCreationDate
	// sea anterior a la fecha actual menos el intervalo de tiempo especificado (timeout en minutos).
	var query string = `DELETE FROM ` + resetPasswordPath +
		` WHERE ` + resetPasswordPath + `.` + ResetPasswordFieldDefinitions["ResetPasswordCreationDate"].DBName +
		` < NOW() - INTERVAL ` + strconv.FormatUint(uint64(timeout), 10) + `' minutes'`

	// Se ejecuta el query sin parámetros adicionales.
	persistenceCtrl.Exec(context.Background(), query)

	// Se comprueba si hubo error durante la ejecución.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// SetResetPasswordDefaults asigna valores predeterminados a los campos de ResetPasswordDTO según la acción a realizar.
// Parámetros:
//   - resetPassword: Puntero al DTO que se va a modificar.
//   - action: Acción que se está realizando (por ejemplo, SQL_INSERT o SQL_UPDATE).
func SetResetPasswordDefaults(resetPassword *ResetPasswordDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se asigna la fecha actual y se genera un nuevo UUID para el ICode.
		resetPassword.ResetPasswordCreationDate = time.Now()
		resetPassword.ResetPasswordICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// En caso de actualización, se pueden agregar otras asignaciones de valores predeterminados si es necesario.
	}
}

// PgDBToDTO convierte un objeto ResetPasswordPgDB (registro de BD) a un ResetPasswordDTO.
// Retorna el DTO con los valores mapeados desde la estructura de la base de datos.
func (obj *ResetPasswordPgDB) PgDBToDTO() ResetPasswordDTO {
	var dto ResetPasswordDTO

	// Verifica si el campo ResetPasswordId es válido y lo asigna al DTO.
	if obj.ResetPasswordId.Valid {
		dto.ResetPasswordId = uint64(obj.ResetPasswordId.Int64)
	}

	// Verifica si el campo ResetPasswordICode es válido y lo asigna.
	if obj.ResetPasswordICode.Valid {
		dto.ResetPasswordICode = obj.ResetPasswordICode.String
	}

	// Verifica si el campo ResetPasswordCreationDate es válido y lo asigna.
	if obj.ResetPasswordCreationDate.Valid {
		dto.ResetPasswordCreationDate = obj.ResetPasswordCreationDate.Time
	}

	// Verifica si el campo ResetPasswordGeneralUser es válido y lo asigna al DTO anidado.
	if obj.ResetPasswordGeneralUser.Valid {
		dto.ResetPasswordGeneralUser = GeneralUserDTO{GeneralUserId: uint64(obj.ResetPasswordGeneralUser.Int64)}
	}

	return dto
}
