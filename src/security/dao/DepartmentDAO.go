// Package security_daos proporciona las funciones de acceso a datos (DAO) para la entidad Department,
// permitiendo la interacción con la base de datos de forma estructurada y validada.
package security_daos

import (
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
	// DepartmentEntityName es el nombre de la entidad en el sistema.
	DepartmentEntityName string = "Department"
	// DepartmentJSONName es el nombre que se utiliza para la entidad en formato JSON.
	DepartmentJSONName string = "department"
	// DepartmentDBName es el nombre de la tabla en la base de datos.
	DepartmentDBName string = "department"
	// DepartmentDBScheme es el esquema de la base de datos donde se encuentra la tabla.
	DepartmentDBScheme string = "security"

	// DepartmentFieldDefinitions contiene la definición y validaciones de los campos de Department.
	// La clave del mapa es el nombre del campo y el valor es una estructura que define las propiedades del mismo.
	DepartmentFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"DepartmentId":           {Name: "DepartmentId", DBName: "department_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"DepartmentICode":        {Name: "DepartmentICode", DBName: "department_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"DepartmentCreationDate": {Name: "DepartmentCreationDate", DBName: "department_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"DepartmentUpdateDate":   {Name: "DepartmentUpdateDate", DBName: "department_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"DepartmentName":         {Name: "DepartmentName", DBName: "department_name", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 64, Required: true},
		"DepartmentCode":         {Name: "DepartmentCode", DBName: "department_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
	}
)

// DepartmentDTO representa el objeto de transferencia de datos para un departamento.
// Se utiliza para enviar y recibir datos en la capa de negocio o en la comunicación JSON.
type DepartmentDTO struct {
	DepartmentId           uint64    `json:"-"`
	DepartmentICode        string    `json:"icode"`
	DepartmentCreationDate time.Time `json:"creationDate"`
	DepartmentUpdateDate   time.Time `json:"updateDate"`
	DepartmentName         string    `json:"name"`
	DepartmentCode         string    `json:"code"`
}

// DepartmentPgDB representa la estructura que mapea los campos de la tabla 'department'
// en la base de datos, utilizando tipos que permiten manejar valores nulos.
type DepartmentPgDB struct {
	DepartmentId           sql.NullInt64
	DepartmentICode        sql.NullString
	DepartmentCreationDate sql.NullTime
	DepartmentUpdateDate   sql.NullTime
	DepartmentName         sql.NullString
	DepartmentCode         sql.NullString
}

func (dp DepartmentDTO) MarshalJSON() ([]byte, error) {
	type Alias DepartmentDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		DepartmentCreationDate string `json:"creationDate"`
		DepartmentUpdateDate   string `json:"updateDate"`
	}{
		Alias:                  (*Alias)(&dp),
		DepartmentCreationDate: dp.DepartmentCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		DepartmentUpdateDate:   dp.DepartmentUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (dp *DepartmentDTO) UnmarshalJSON(data []byte) error {
	type Alias DepartmentDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		DepartmentCreationDate string `json:"creationDate"`
		DepartmentUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(dp),
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
	dp.DepartmentCreationDate = parse(aux.DepartmentCreationDate, common_config.DateTime.DATE_TIME_FORMAT)
	dp.DepartmentUpdateDate = parse(aux.DepartmentUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// GetDepartment obtiene un único registro de un departamento de la base de datos,
// según los criterios de filtrado especificados en el parámetro 'by'.
// Recibe además un puntero a DepartmentDTO para almacenar el resultado y parámetros
// de configuración de la conexión.
// Retorna la conexión actualizada y un error en caso de ocurrir problemas durante la consulta.
func GetDepartment(by common_controllers.By, department *DepartmentDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Define la ruta completa (esquema y tabla) para la consulta.
	var departmentPath string = DepartmentDBScheme + "." + DepartmentDBName

	// Configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Define los campos que se seleccionarán en la consulta.
	var departmentFieldsSlice []string = []string{"DepartmentId", "DepartmentICode", "DepartmentCreationDate", "DepartmentUpdateDate", "DepartmentName", "DepartmentCode"}
	var departmentFieldsAliasSlice []string = []string{}

	// Construye dinámicamente la lista de campos SQL a seleccionar.
	var departmentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, departmentFieldsSlice, departmentFieldsAliasSlice, DepartmentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, DepartmentDBScheme, DepartmentFieldDefinitions, true)

	// Construye el query SQL completo, incluyendo la cláusula WHERE generada dinámicamente.
	var query string = `SELECT ` + departmentFieldsStr +
		` FROM ` + departmentPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, DepartmentDBName, by.AttrsName, []string{}, []string{}, by.Operator, DepartmentDBScheme, DepartmentFieldDefinitions, true)

	// Ejecuta el query y obtiene una fila de resultado.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Escanea los valores obtenidos en la estructura DepartmentDTO.
	persistenceCtrl.Scan(&department.DepartmentId, &department.DepartmentICode, &department.DepartmentCreationDate, &department.DepartmentUpdateDate,
		&department.DepartmentName, &department.DepartmentCode)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetDepartments obtiene uno o varios registros de departamentos de la base de datos,
// basándose en los criterios de filtrado especificados en el parámetro 'by'.
// Recibe parámetros de configuración para la conexión y retorna un slice de DepartmentDTO,
// junto con la conexión actualizada y un error en caso de producirse.
func GetDepartments(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]DepartmentDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Define la ruta completa (esquema y tabla) para la consulta.
	var departmentPath string = DepartmentDBScheme + "." + DepartmentDBName

	var department = DepartmentDTO{}
	var departments []DepartmentDTO

	// Configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var departmentFieldsSlice []string = []string{"DepartmentId", "DepartmentICode", "DepartmentCreationDate", "DepartmentUpdateDate", "DepartmentName", "DepartmentCode"}
	var departmentFieldsAliasSlice []string = []string{}

	// Construye la parte de selección de campos del query.
	var departmentFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, departmentFieldsSlice, departmentFieldsAliasSlice, DepartmentDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, DepartmentDBScheme, DepartmentFieldDefinitions, true)

	// Construye el query SQL completo con la cláusula WHERE.
	var query string = `SELECT ` + departmentFieldsStr +
		` FROM ` + departmentPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, DepartmentDBName, by.AttrsName, []string{}, []string{}, by.Operator, DepartmentDBScheme, DepartmentFieldDefinitions, true)

	// Ejecuta el query para obtener múltiples registros.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre cada fila obtenida y la escanea en la estructura DepartmentDTO.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&department.DepartmentId, &department.DepartmentICode, &department.DepartmentCreationDate, &department.DepartmentUpdateDate, &department.DepartmentName, &department.DepartmentCode)
		departments = append(departments, department)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return departments, nil
}

// GetAllDepartment obtiene todos los registros de departamentos de la base de datos sin aplicar filtros.
// Recibe parámetros de configuración de la conexión y retorna un slice de DepartmentDTO,
// junto con la conexión actualizada y un error en caso de producirse.
func GetAllDepartment(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]DepartmentDTO, error) {
	// Declaración de variables locales.
	var department DepartmentDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Configura la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Define los campos a seleccionar.
	var departmentFieldsSlice []string = []string{"DepartmentId", "DepartmentICode", "DepartmentCreationDate", "DepartmentUpdateDate", "DepartmentName", "DepartmentCode"}
	var departmentFieldsAliasSlice []string = []string{}

	// Construye el query SQL para seleccionar todos los registros de la tabla.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, departmentFieldsSlice, departmentFieldsAliasSlice, DepartmentDBName, []string{}, []string{}, []string{}, "", DepartmentDBScheme, DepartmentFieldDefinitions, true)

	// Ejecuta el query.
	persistenceCtrl.Query(context.Background(), query)
	var departments []DepartmentDTO
	// Itera sobre cada fila del resultado, escaneándola en la estructura DepartmentDTO.
	for persistenceCtrl.Next() {
		department = DepartmentDTO{}
		persistenceCtrl.ScanRow(&department.DepartmentId, &department.DepartmentICode, &department.DepartmentCreationDate, &department.DepartmentUpdateDate,
			&department.DepartmentName, &department.DepartmentCode)
		departments = append(departments, department)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return departments, nil
}

// SetDepartmentDefaults establece valores predeterminados en un objeto DepartmentDTO
// según la acción que se vaya a realizar (inserción o actualización).
func SetDepartmentDefaults(department *DepartmentDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// Para una inserción, se asigna la fecha actual como creación y actualización,
		// y se genera un nuevo UUID para DepartmentICode.
		department.DepartmentCreationDate = time.Now()
		department.DepartmentUpdateDate = time.Now()
		department.DepartmentICode = utils.GetUUID()

	case common_dao.SQL_UPDATE:
		// Para una actualización, solo se actualiza la fecha de modificación.
		department.DepartmentUpdateDate = time.Now()
	}
}

// PgDBToDTO convierte un objeto DepartmentPgDB (representación de la base de datos)
// a un objeto DepartmentDTO (utilizado en la capa de negocio).
// Retorna el objeto DepartmentDTO resultante.
func (obj *DepartmentPgDB) PgDBToDTO() DepartmentDTO {
	var dto DepartmentDTO

	if obj.DepartmentId.Valid {
		dto.DepartmentId = uint64(obj.DepartmentId.Int64)
	}

	if obj.DepartmentICode.Valid {
		dto.DepartmentICode = obj.DepartmentICode.String
	}

	if obj.DepartmentCreationDate.Valid {
		dto.DepartmentCreationDate = obj.DepartmentCreationDate.Time
	}

	if obj.DepartmentUpdateDate.Valid {
		dto.DepartmentUpdateDate = obj.DepartmentUpdateDate.Time
	}

	if obj.DepartmentName.Valid {
		dto.DepartmentName = obj.DepartmentName.String
	}

	if obj.DepartmentCode.Valid {
		dto.DepartmentCode = obj.DepartmentCode.String
	}

	return dto
}
