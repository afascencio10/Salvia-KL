// Package salvia_daos contiene los Data Access Objects (DAO) para la entidad Moment,
// incluyendo funciones para insertar, actualizar, eliminar y recuperar registros de Moment.
package salvia_daos

import (
	// Importación de configuraciones, controladores y utilidades comunes
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
	// Constantes que representan los nombres de la entidad y de la tabla en distintos contextos.
	MomentEntityName string = "Moment"
	MomentJSONName   string = "moment"
	MomentDBName     string = "moment"
	MomentDBScheme   string = "salvia"

	// Definición de campos y validaciones para la entidad Moment.
	// Estos campos indican la correspondencia entre los nombres del JSON, los de la base de datos
	// y el tipo de dato esperado en el modelo. El booleano 'Required' indica si el campo es obligatorio.
	MomentFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"MomentId":                   {Name: "MomentId", DBName: "moment_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"MomentICode":                {Name: "MomentICode", DBName: "moment_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: true},
		"MomentCreationDate":         {Name: "MomentCreationDate", DBName: "moment_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"MomentUpdateDate":           {Name: "MomentUpdateDate", DBName: "moment_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"MomentApprovalCancellation": {Name: "MomentApprovalCancellation", DBName: "moment_approval_cancellation", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"MomentCode":                 {Name: "MomentCode", DBName: "moment_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: false},
		"MomentStatus":               {Name: "MomentStatus", DBName: "moment_status", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: true},
		"MomentApprovalSource":       {Name: "MomentApprovalSource", DBName: "moment_approval_source", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 1, Required: false},
		"MomentVictimCase":           {Name: "MomentVictimCase", DBName: "moment_victim_case", Alias: "", ModelType: "uint", Required: false},
		"MomentEntityBranch":         {Name: "MomentEntityBranch", DBName: "moment_entity_branch", Alias: "", ModelType: "uint", Required: false},
		"MomentApprovalOwner":        {Name: "MomentApprovalOwner", DBName: "moment_approval_owner", Alias: "", ModelType: "uint", Required: false},
		"MomentApprovalDescription":  {Name: "MomentApprovalDescription", DBName: "moment_approval_description", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 50000, Required: false},
	}
)

// MomentDTO representa la entidad Moment en el modelo de datos (Data Transfer Object).
// Contiene tanto los campos directamente mapeados a la base de datos como algunos campos de formulario.
type MomentDTO struct {
	MomentId                   uint64          `json:"-"`
	MomentICode                string          `json:"icode"`
	MomentCreationDate         time.Time       `json:"creationDate"`
	MomentUpdateDate           time.Time       `json:"updateDate"`
	MomentApprovalCancellation string          `json:"canceled"`
	MomentCode                 string          `json:"code"`
	MomentStatus               string          `json:"status"`
	MomentApprovalSource       string          `json:"approvalSource"`
	MomentApprovalDescription  string          `json:"approvalDescription"`
	MomentVictimCase           VictimCaseDTO   `json:"victimCase"`
	MomentEntityBranch         EntityBranchDTO `json:"entityBranch"`
	MomentApprovalOwner        CaseOwnerDTO    `json:"approvalOwner"`
	// Campos adicionales para el formulario que no forman parte directa del modelo en la BD.
	MomentCaseLogs []CaseLogDTO `json:"logs"`
}

// MomentPgDB representa la estructura de la entidad Moment tal como se almacena en la base de datos,
// utilizando los tipos sql.Null* para manejar valores nulos.
type MomentPgDB struct {
	MomentId                   sql.NullInt64
	MomentICode                sql.NullString
	MomentCreationDate         sql.NullTime
	MomentUpdateDate           sql.NullTime
	MomentApprovalCancellation sql.NullString
	MomentCode                 sql.NullString
	MomentStatus               sql.NullString
	MomentApprovalSource       sql.NullString
	MomentApprovalDescription  sql.NullString
	MomentVictimCase           sql.NullInt64
	MomentEntityBranch         sql.NullInt64
	MomentApprovalOwner        sql.NullInt64
	// Campos adicionales del formulario que no se almacenan en la base de datos.
}

func (m MomentDTO) MarshalJSON() ([]byte, error) {
	type Alias MomentDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		MomentCreationDate string `json:"creationDate"`
		MomentUpdateDate   string `json:"updateDate"`
	}{
		Alias:              (*Alias)(&m),
		MomentCreationDate: m.MomentCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		MomentUpdateDate:   m.MomentUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (m *MomentDTO) UnmarshalJSON(data []byte) error {
	type Alias MomentDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		MomentCreationDate string `json:"creationDate"`
		MomentUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(m),
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
	m.MomentCreationDate = parse(aux.MomentCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	m.MomentUpdateDate = parse(aux.MomentUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// SetMoment inserta un nuevo registro de Moment en la base de datos.
//
// Parámetros:
// - moment: Puntero a la estructura MomentDTO que contiene los datos a insertar.
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - error en caso de fallo, o nil si la operación fue exitosa.
func SetMoment(moment *MomentDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a insertar y generación de la consulta SQL.
	var momentFieldsSlice []string = []string{
		"MomentICode", "MomentCreationDate", "MomentUpdateDate",
		"MomentApprovalCancellation", "MomentCode", "MomentStatus",
		"MomentApprovalSource", "MomentVictimCase", "MomentEntityBranch",
	}
	var momentFieldsAliasSlice []string = []string{}
	var query string = common_dao.GetSQL(
		common_dao.SQL_INSERT,
		momentFieldsSlice,
		momentFieldsAliasSlice,
		MomentDBName,
		[]string{},
		[]string{},
		[]string{"MomentId"},
		common_dao.SQL_AND,
		MomentDBScheme,
		MomentFieldDefinitions,
		false,
	)

	// Ejecución de la consulta para insertar el registro.
	persistenceCtrl.QueryRow(
		context.Background(),
		query,
		moment.MomentICode,
		moment.MomentCreationDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		moment.MomentUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		moment.MomentApprovalCancellation,
		moment.MomentCode,
		moment.MomentStatus,
		moment.MomentApprovalSource,
		moment.MomentVictimCase.VictimCaseId,
		moment.MomentEntityBranch.EntityBranchId,
	)
	// Se extrae el ID generado del registro insertado.
	persistenceCtrl.Scan(&moment.MomentId)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// UpdateMomentById actualiza un registro de Moment identificado por su ID.
//
// Parámetros:
// - moment: Puntero a la estructura MomentDTO que contiene los datos a actualizar.
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - error en caso de fallo o si no se afectó ninguna fila.
func UpdateMomentById(moment *MomentDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a actualizar y generación de la consulta SQL.
	var momentFieldsSlice []string = []string{
		"MomentUpdateDate", "MomentApprovalCancellation", "MomentCode",
		"MomentStatus", "MomentApprovalSource", "MomentApprovalDescription",
		"MomentVictimCase", "MomentEntityBranch", "MomentApprovalOwner",
	}
	var momentFieldsAliasSlice []string = []string{}
	var query string = common_dao.GetSQL(
		common_dao.SQL_UPDATE,
		momentFieldsSlice,
		momentFieldsAliasSlice,
		MomentDBName,
		[]string{"MomentId"},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		MomentDBScheme,
		MomentFieldDefinitions,
		false,
	)

	// Ejecución de la consulta para actualizar el registro.
	persistenceCtrl.Exec(
		context.Background(),
		query,
		moment.MomentId,
		moment.MomentUpdateDate.Format(common_config.DateTime.DB_DATE_TIME_FORMAT),
		moment.MomentApprovalCancellation,
		moment.MomentCode,
		moment.MomentStatus,
		moment.MomentApprovalSource,
		moment.MomentApprovalDescription,
		moment.MomentVictimCase.VictimCaseId,
		moment.MomentEntityBranch.EntityBranchId,
		moment.MomentApprovalOwner.CaseOwnerId,
	)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica si la actualización afectó alguna fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveMomentById elimina un registro de Moment identificado por su ID.
//
// Parámetros:
// - moment: Puntero a la estructura MomentDTO que contiene el ID del registro a eliminar.
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - error en caso de fallo o si no se afectó ninguna fila.
func RemoveMomentById(moment *MomentDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Generación de la consulta SQL para eliminar el registro.
	var momentFieldsSlice []string = []string{}
	var momentFieldsAliasSlice []string = []string{}
	var query string = common_dao.GetSQL(common_dao.SQL_DELETE, momentFieldsSlice, momentFieldsAliasSlice, MomentDBName, []string{"MomentId"}, []string{}, []string{}, common_dao.SQL_AND, MomentDBScheme, MomentFieldDefinitions, false)

	// Ejecución de la consulta para eliminar el registro.
	persistenceCtrl.Exec(context.Background(), query, moment.MomentId)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica si la eliminación afectó alguna fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// RemoveMomentByVictimCaseId elimina los registros de Moment asociados a un VictimCase específico.
//
// Parámetros:
// - victimCaseId: ID del VictimCase cuyo Moment se desea eliminar.
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - error en caso de fallo o si no se afectó ninguna fila.
func RemoveMomentByVictimCaseId(victimCaseId uint64, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Generación de la consulta SQL para eliminar el registro por VictimCase.
	var momentFieldsSlice []string = []string{}
	var momentFieldsAliasSlice []string = []string{}
	var query string = common_dao.GetSQL(
		common_dao.SQL_DELETE,
		momentFieldsSlice,
		momentFieldsAliasSlice,
		MomentDBName,
		[]string{"MomentVictimCase"},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		MomentDBScheme,
		MomentFieldDefinitions,
		false,
	)

	// Ejecución de la consulta.
	persistenceCtrl.Exec(context.Background(), query, victimCaseId)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Verifica si la eliminación afectó alguna fila.
	if persistenceCtrl.RowsAffected == 0 {
		return errors.New(common_config.Locale["sp"]["common_db_no_rows_affected"])
	}

	return nil
}

// GetMoment recupera un registro de Moment basado en un criterio de búsqueda especificado.
//
// Parámetros:
// - by: Estructura que contiene el atributo y valor para filtrar la búsqueda.
// - moment: Puntero a MomentDTO donde se asignarán los datos recuperados.
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - error en caso de fallo.
func GetMoment(by common_controllers.By, moment *MomentDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia y definición de la ruta completa de la tabla.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	var momentPath string = MomentDBScheme + "." + MomentDBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos a recuperar y generación de la cláusula SELECT.
	var momentFieldsSlice []string = []string{
		"MomentId", "MomentICode", "MomentCreationDate", "MomentUpdateDate",
		"MomentApprovalCancellation", "MomentCode", "MomentStatus",
		"MomentApprovalSource", "MomentApprovalDescription",
		"MomentVictimCase", "MomentEntityBranch", "MomentApprovalOwner",
	}
	var momentFieldsAliasSlice []string = []string{}
	var momentFieldsStr = common_dao.GetSQL(
		common_dao.SQL_SELECT_FIELDS_ONLY,
		momentFieldsSlice,
		momentFieldsAliasSlice,
		MomentDBName,
		[]string{},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		MomentDBScheme,
		MomentFieldDefinitions,
		true,
	)

	// Generación de la consulta SELECT completa, utilizando el filtro definido en 'by'.
	var query string = `SELECT ` + momentFieldsStr +
		` FROM ` + momentPath +
		common_dao.GetSQL(
			common_dao.SQL_SELECT_WHERE_ONLY,
			by.AttrsName,
			by.AttrsAliasName,
			MomentDBName,
			by.AttrsName,
			[]string{},
			[]string{},
			by.Operator,
			MomentDBScheme,
			MomentFieldDefinitions,
			true,
		)

	// Ejecución de la consulta.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Mapeo de los resultados a la estructura MomentPgDB y conversión a DTO.
	var momentPG = MomentPgDB{}
	persistenceCtrl.Scan(
		&momentPG.MomentId, &momentPG.MomentICode, &momentPG.MomentCreationDate, &momentPG.MomentUpdateDate,
		&momentPG.MomentApprovalCancellation, &momentPG.MomentCode, &momentPG.MomentStatus,
		&momentPG.MomentApprovalSource, &momentPG.MomentApprovalDescription,
		&momentPG.MomentVictimCase, &momentPG.MomentEntityBranch, &momentPG.MomentApprovalOwner,
	)
	var tmp MomentDTO = momentPG.ToDTO()
	*moment = tmp

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetMoments recupera una lista de registros de Moment basados en un criterio de búsqueda.
// Además, realiza un LEFT JOIN con la tabla de EntityBranch para obtener información relacionada.
//
// Parámetros:
// - by: Estructura que contiene el atributo y valor para filtrar la búsqueda.
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - Slice de MomentDTO con los registros obtenidos.
// - error en caso de fallo.
func GetMoments(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]MomentDTO, error) {
	// Inicialización de variables y del controlador de persistencia.
	var moment MomentPgDB
	var momentPath string = MomentDBScheme + "." + MomentDBName
	// Las siguientes variables hacen referencia a la tabla EntityBranch, definida externamente.
	var entityBranchPath string = EntityBranchDBScheme + "." + EntityBranchDBName
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a recuperar de la tabla Moment.
	var momentFieldsSlice []string = []string{
		"MomentId", "MomentICode", "MomentCreationDate", "MomentUpdateDate",
		"MomentApprovalCancellation", "MomentCode", "MomentStatus",
		"MomentApprovalSource", "MomentApprovalDescription",
		"MomentVictimCase", "MomentEntityBranch", "MomentApprovalOwner",
	}
	var momentFieldsAliasSlice []string = []string{}
	var momentFieldsStr = common_dao.GetSQL(
		common_dao.SQL_SELECT_FIELDS_ONLY,
		momentFieldsSlice,
		momentFieldsAliasSlice,
		MomentDBName,
		[]string{},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		MomentDBScheme,
		MomentFieldDefinitions,
		true,
	)

	// Definición de los campos a recuperar de la tabla EntityBranch.
	var entityBranchFieldsSlice []string = []string{
		"EntityBranchId", "EntityBranchICode", "EntityBranchName", "EntityBranchTownCode", "EntityBranchEntity",
	}
	var entityBranchFieldsAliasSlice []string = []string{}
	var entityBranchFieldsStr = common_dao.GetSQL(
		common_dao.SQL_SELECT_FIELDS_ONLY,
		entityBranchFieldsSlice,
		entityBranchFieldsAliasSlice,
		EntityBranchDBName,
		[]string{},
		[]string{},
		[]string{},
		common_dao.SQL_AND,
		EntityBranchDBScheme,
		EntityBranchFieldDefinitions,
		true,
	)

	// Construcción de la consulta con LEFT JOIN para relacionar Moment con EntityBranch.
	var query string = `SELECT ` + momentFieldsStr + ", " + entityBranchFieldsStr +
		` FROM ` + momentPath +
		` LEFT JOIN ` + entityBranchPath + ` ON (` +
		entityBranchPath + `.` + EntityBranchFieldDefinitions["EntityBranchId"].DBName + ` = ` +
		momentPath + `.` + MomentFieldDefinitions["MomentEntityBranch"].DBName + `)` +
		common_dao.GetSQL(
			common_dao.SQL_SELECT_WHERE_ONLY,
			by.AttrsName,
			by.AttrsAliasName,
			MomentDBName,
			by.AttrsName,
			[]string{},
			[]string{},
			by.Operator,
			MomentDBScheme,
			MomentFieldDefinitions,
			true,
		)

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)
	var moments []MomentDTO
	// Itera sobre cada registro recuperado.
	for persistenceCtrl.Next() {
		var branch EntityBranchPgDB = EntityBranchPgDB{}
		moment = MomentPgDB{}
		persistenceCtrl.ScanRow(
			&moment.MomentId, &moment.MomentICode, &moment.MomentCreationDate, &moment.MomentUpdateDate,
			&moment.MomentApprovalCancellation, &moment.MomentCode, &moment.MomentStatus,
			&moment.MomentApprovalSource, &moment.MomentApprovalDescription,
			&moment.MomentVictimCase, &moment.MomentEntityBranch, &moment.MomentApprovalOwner,
			&branch.EntityBranchId, &branch.EntityBranchICode, &branch.EntityBranchName,
			&branch.EntityBranchTownCode, &branch.EntityBranchEntity,
		)

		// Conversión de la entidad de base de datos a DTO y asignación del branch.
		var m MomentDTO = moment.ToDTO()
		m.MomentEntityBranch = branch.ToDTO()
		moments = append(moments, m)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return moments, nil
}

// GetAllMoment recupera todos los registros de Moment sin aplicar filtros.
//
// Parámetros:
// - connData: Datos de la conexión a la base de datos.
// - clientConfig: Configuración del cliente de la base de datos.
// - serverConfig: Configuración del servidor de la base de datos.
//
// Retorna:
//
// - Slice de MomentDTO con todos los registros recuperados.
// - error en caso de fallo.
func GetAllMoment(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]MomentDTO, error) {
	var moment MomentDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a recuperar y generación de la consulta SELECT.
	var momentFieldsSlice []string = []string{
		"MomentId", "MomentICode", "MomentStatus", "MomentCreationDate",
		"MomentUpdateDate", "MomentApprovalCancellation", "MomentCode",
		"MomentApprovalSource", "MomentVictimCase", "MomentEntityBranch", "MomentApprovalOwner",
	}
	var momentFieldsAliasSlice []string = []string{}
	var query string = common_dao.GetSQL(
		common_dao.SQL_SELECT,
		momentFieldsSlice,
		momentFieldsAliasSlice,
		MomentDBName,
		[]string{},
		[]string{},
		[]string{},
		"",
		MomentDBScheme,
		MomentFieldDefinitions,
		true,
	)

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var moments []MomentDTO
	// Iteración sobre cada registro recuperado.
	for persistenceCtrl.Next() {
		moment = MomentDTO{}
		persistenceCtrl.ScanRow(
			&moment.MomentId, &moment.MomentICode, &moment.MomentStatus,
			&moment.MomentCreationDate, &moment.MomentUpdateDate,
			&moment.MomentApprovalCancellation, &moment.MomentCode,
			&moment.MomentApprovalSource,
			&moment.MomentVictimCase.VictimCaseId, &moment.MomentEntityBranch.EntityBranchId,
			&moment.MomentApprovalOwner.CaseOwnerId,
		)
		moments = append(moments, moment)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return moments, nil
}

// SetMomentDefaults asigna valores por defecto a los campos de Moment según la acción a realizar.
// Por ejemplo, para una inserción se asigna la fecha actual y se genera un UUID.
func SetMomentDefaults(moment *MomentDTO, action string) {
	switch action {
	// Para inserción: se asigna la fecha de creación y se genera un UUID único.
	case common_dao.SQL_INSERT:
		moment.MomentCreationDate = time.Now()
		moment.MomentICode = utils.GetUUID()
	// Para actualización: se asigna la fecha de actualización.
	case common_dao.SQL_UPDATE:
		moment.MomentUpdateDate = time.Now()
	}
}

// ToDTO convierte una instancia de MomentPgDB (la representación en base de datos) a MomentDTO,
// la representación utilizada en el modelo de datos.
func (obj *MomentPgDB) ToDTO() MomentDTO {
	var dto MomentDTO

	if obj.MomentId.Valid {
		dto.MomentId = uint64(obj.MomentId.Int64)
	}
	if obj.MomentICode.Valid {
		dto.MomentICode = obj.MomentICode.String
	}
	if obj.MomentStatus.Valid {
		dto.MomentStatus = obj.MomentStatus.String
	}
	if obj.MomentCreationDate.Valid {
		dto.MomentCreationDate = obj.MomentCreationDate.Time
	}
	if obj.MomentUpdateDate.Valid {
		dto.MomentUpdateDate = obj.MomentUpdateDate.Time
	}
	if obj.MomentApprovalCancellation.Valid {
		dto.MomentApprovalCancellation = obj.MomentApprovalCancellation.String
	}
	if obj.MomentCode.Valid {
		dto.MomentCode = obj.MomentCode.String
	}
	if obj.MomentApprovalSource.Valid {
		dto.MomentApprovalSource = obj.MomentApprovalSource.String
	}
	if obj.MomentApprovalDescription.Valid {
		dto.MomentApprovalDescription = obj.MomentApprovalDescription.String
	}
	if obj.MomentVictimCase.Valid {
		dto.MomentVictimCase = VictimCaseDTO{VictimCaseId: uint64(obj.MomentVictimCase.Int64)}
	}
	if obj.MomentApprovalOwner.Valid {
		dto.MomentApprovalOwner = CaseOwnerDTO{CaseOwnerId: uint64(obj.MomentApprovalOwner.Int64)}
	}
	if obj.MomentEntityBranch.Valid {
		dto.MomentEntityBranch = EntityBranchDTO{EntityBranchId: uint64(obj.MomentEntityBranch.Int64)}
	}

	return dto
}
