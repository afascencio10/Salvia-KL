// Package salvia_daos proporciona las funciones y estructuras para gestionar la relación entre alertas y casos de víctimas
// en la base de datos. Incluye la definición de DTOs, estructuras para el manejo de datos SQL y funciones para insertar y
// convertir registros entre la representación de la base de datos y el DTO.
package salvia_daos

import (
	// Importación de paquetes comunes para configuración, controladores, acceso a datos, manejo de base de datos y utilidades.
	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

var (
	// Constantes y variables para la configuración de la tabla de relación entre alertas y casos de víctimas.
	RelAlertVictimCase         string = "RelAlertVictimCase"    // Nombre lógico de la entidad
	RelAlertVictimCaseJSONName string = "relAlert"              // Nombre del campo en JSON
	RelAlertVictimCaseDBName   string = "rel_alert_victim_case" // Nombre de la tabla en la base de datos
	RelAlertVictimCaseDBScheme string = "salvia"                // Esquema de la base de datos

	// Definiciones de los campos que se utilizan para las validaciones y el mapeo entre el JSON, el modelo y la base de datos.
	// La definición incluye el nombre, nombre en la base de datos, alias, tipo de dato del modelo, tamaño mínimo, tamaño máximo
	// y si es obligatorio o no.
	RelAlertVictimCaseFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"RelAlertVictimCase_Id":           {Name: "RelAlertVictimCase_Id", DBName: "rel_alert_victim_case_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelAlertVictimCase_Alert":        {Name: "RelAlertVictimCase_Alert", DBName: "rel_alert_victim_case_alert_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelAlertVictimCase_VictimCase":   {Name: "RelAlertVictimCase_VictimCase", DBName: "rel_alert_victim_case_victim_case_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: true},
		"RelAlertVictimCase_Moment":       {Name: "RelAlertVictimCase_Moment", DBName: "rel_alert_victim_case_victim_moment_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"RelAlertVictimCase_CreationDate": {Name: "RelAlertVictimCase_CreationDate", DBName: "rel_alert_victim_case_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: true},
		"RelAlertVictimCase_Data":         {Name: "RelAlertVictimCase_Data", DBName: "rel_alert_victim_case_data", Alias: "", ModelType: "string", MinSize: 3, MaxSize: 128, Required: false},
	}
)

// RelAlertVictimCaseDTO representa el objeto de transferencia de datos para la relación entre una alerta y un caso de víctima.
// Este DTO se utiliza para intercambiar información entre las diferentes capas de la aplicación.
type RelAlertVictimCaseDTO struct {
	RelAlertVictimCaseId            uint64    `json:"-"`            // Identificador único de la relación (no se expone en JSON)
	RelAlertVictimCase_Alert        uint64    `json:"alert"`        // Identificador de la alerta asociada
	RelAlertVictimCase_VictimCase   uint64    `json:"victimCase"`   // Identificador del caso de víctima asociado
	RelAlertVictimCase_Moment       uint64    `json:"moment"`       // Identificador del momento relacionado (si aplica)
	RelAlertVictimCase_CreationDate time.Time `json:"creationDate"` // Fecha de creación de la relación
	RelAlertVictimCase_Data         string    `json:"data"`         // Datos adicionales relacionados con la relación
}

// RelAlertVictimCasePgDB representa la estructura utilizada para mapear los datos de la base de datos,
// usando tipos nulos de SQL para manejar valores que pueden ser NULL.
type RelAlertVictimCasePgDB struct {
	RelAlertVictimCaseId            sql.NullInt64  // Mapea el campo identificador de la relación en la base de datos
	RelAlertVictimCase_Alert        sql.NullInt64  // Mapea el identificador de la alerta
	RelAlertVictimCase_VictimCase   sql.NullInt64  // Mapea el identificador del caso de víctima
	RelAlertVictimCase_Moment       sql.NullInt64  // Mapea el identificador del momento relacionado
	RelAlertVictimCase_CreationDate sql.NullTime   // Mapea la fecha de creación de la relación
	RelAlertVictimCase_Data         sql.NullString // Mapea los datos adicionales (puede ser NULL)
}

func (rcwvc RelAlertVictimCaseDTO) MarshalJSON() ([]byte, error) {
	type Alias RelAlertVictimCaseDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		RelAlertVictimCase_CreationDate string `json:"creationDate"`
	}{
		Alias:                           (*Alias)(&rcwvc),
		RelAlertVictimCase_CreationDate: rcwvc.RelAlertVictimCase_CreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (rcwvc *RelAlertVictimCaseDTO) UnmarshalJSON(data []byte) error {
	type Alias RelAlertVictimCaseDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		RelAlertVictimCase_CreationDate string `json:"creationDate"`
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
	rcwvc.RelAlertVictimCase_CreationDate = parse(aux.RelAlertVictimCase_CreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetRelAlertVictimCase inserta un nuevo registro de relación entre alerta y caso de víctima en la base de datos.
// Recibe el DTO con los datos a insertar, información sobre la transacción, configuración del cliente y servidor de la base de datos,
// y retorna la conexión utilizada junto con cualquier error que se haya producido.
//
// Parámetros:
//   - relAlertVictimCase: Puntero al DTO que contiene los datos de la relación.
//   - connData: Datos de conexión actuales; puede ser actualizado durante la operación.
//   - clientConfig: Configuración del cliente de la base de datos.
//   - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//   - error: Error encontrado durante la operación, o nil si la operación fue exitosa.
func SetRelAlertVictimCase(relAlertVictimCase *RelAlertVictimCaseDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia que facilita operaciones comunes sobre la base de datos.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configuración y obtención de la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		// Se registra el error en la conexión y se retorna.
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos que serán insertados en la tabla.
	var relAlertVictimCaseFieldsSlice []string = []string{
		"RelAlertVictimCase_Alert",
		"RelAlertVictimCase_VictimCase",
		"RelAlertVictimCase_Moment",
		"RelAlertVictimCase_CreationDate",
		"RelAlertVictimCase_Data",
	}
	// En este caso no se utilizan alias para los campos en la consulta.
	var relAlertVictimCaseFieldsAliasSlice []string = []string{}

	// Generación de la consulta SQL de inserción utilizando la función GetSQL.
	var query string = common_dao.GetSQL(
		common_dao.SQL_INSERT,
		relAlertVictimCaseFieldsSlice,
		relAlertVictimCaseFieldsAliasSlice,
		RelAlertVictimCaseDBName,
		[]string{},
		[]string{},
		[]string{"RelAlertVictimCase_Id"},
		common_dao.SQL_AND,
		RelAlertVictimCaseDBScheme,
		RelAlertVictimCaseFieldDefinitions,
		false,
	)

	// Ejecución de la consulta SQL, pasando los parámetros necesarios.
	persistenceCtrl.QueryRow(
		context.Background(),
		query,
		relAlertVictimCase.RelAlertVictimCase_Alert,
		relAlertVictimCase.RelAlertVictimCase_VictimCase,
		relAlertVictimCase.RelAlertVictimCase_Moment,
		relAlertVictimCase.RelAlertVictimCase_CreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		relAlertVictimCase.RelAlertVictimCase_Data,
	)

	// Escaneo del resultado para obtener el ID generado de la inserción.
	persistenceCtrl.Scan(&relAlertVictimCase.RelAlertVictimCaseId)
	if persistenceCtrl.Error != nil {
		// En caso de error, se registra la consulta y el error producido.
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Retorna la conexión utilizada y nil en error, indicando éxito en la operación.
	return nil
}

// SetRelAlertVictimCaseDefaults asigna valores por defecto al DTO de relación entre alerta y caso de víctima
// en función de la acción a realizar (por ejemplo, inserción o actualización).
//
// Parámetros:
//   - relEntity: Puntero al DTO al que se le asignarán los valores por defecto.
//   - action: Tipo de acción (e.g., common_dao.SQL_INSERT o common_dao.SQL_UPDATE).
func SetRelAlertVictimCaseDefaults(relEntity *RelAlertVictimCaseDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se establece la fecha de creación al momento actual.
		relEntity.RelAlertVictimCase_CreationDate = time.Now()
	case common_dao.SQL_UPDATE:
		// Para una actualización, actualmente no se asigna ningún valor por defecto.
	}
}

// PgDBToDTO convierte la estructura de datos proveniente de la base de datos (que utiliza tipos SQL nulos)
// a un DTO (Data Transfer Object) que se utiliza en la capa de negocio.
// Solo se convierten los campos que son válidos (no nulos) en la estructura de la base de datos.
//
// Retorna:
//   - RelAlertVictimCaseDTO: Objeto de transferencia de datos con la información convertida.
func (obj *RelAlertVictimCasePgDB) PgDBToDTO() RelAlertVictimCaseDTO {
	var dto RelAlertVictimCaseDTO

	// Conversión del campo RelAlertVictimCase_CreationDate si es válido.
	if obj.RelAlertVictimCase_CreationDate.Valid {
		dto.RelAlertVictimCase_CreationDate = obj.RelAlertVictimCase_CreationDate.Time
	}

	// Conversión del campo RelAlertVictimCase_Data si es válido.
	if obj.RelAlertVictimCase_CreationDate.Valid {
		dto.RelAlertVictimCase_Data = obj.RelAlertVictimCase_Data.String
	}

	// Conversión del campo RelAlertVictimCaseId si es válido.
	if obj.RelAlertVictimCaseId.Valid {
		dto.RelAlertVictimCaseId = uint64(obj.RelAlertVictimCaseId.Int64)
	}

	// Conversión del campo RelAlertVictimCase_Alert si es válido.
	if obj.RelAlertVictimCase_Alert.Valid {
		dto.RelAlertVictimCase_Alert = uint64(obj.RelAlertVictimCase_Alert.Int64)
	}

	// Conversión del campo RelAlertVictimCase_VictimCase si es válido.
	if obj.RelAlertVictimCase_VictimCase.Valid {
		dto.RelAlertVictimCase_VictimCase = uint64(obj.RelAlertVictimCase_VictimCase.Int64)
	}

	// Conversión del campo RelAlertVictimCase_Moment si es válido.
	if obj.RelAlertVictimCase_Moment.Valid {
		dto.RelAlertVictimCase_Moment = uint64(obj.RelAlertVictimCase_Moment.Int64)
	}

	return dto
}
