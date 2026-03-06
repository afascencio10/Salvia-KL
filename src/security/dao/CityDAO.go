// Package security_daos contiene las funciones y estructuras de acceso a datos para la entidad City.
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
	// CityEntityName es el nombre de la entidad de ciudad.
	CityEntityName string = "City"
	// CityJSONName es el nombre usado en JSON para representar una ciudad.
	CityJSONName string = "city"
	// CityDBName es el nombre de la tabla de ciudades en la base de datos.
	CityDBName string = "city"
	// CityDBScheme es el esquema de la base de datos donde se encuentra la tabla de ciudades.
	CityDBScheme string = "security"

	// CityFieldDefinitions contiene la definición de los campos de la entidad City, incluyendo mapeos de nombres,
	// tipos de datos, tamaños mínimos y máximos, y si el campo es obligatorio.
	// La clave del mapa es el nombre del campo en el modelo y el valor es una estructura FieldDefinition con
	// la configuración correspondiente.
	CityFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"CityId":           {Name: "CityId", DBName: "city_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"CityICode":        {Name: "CityICode", DBName: "city_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"CityCreationDate": {Name: "CityCreationDate", DBName: "city_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"CityUpdateDate":   {Name: "CityUpdateDate", DBName: "city_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"CityName":         {Name: "CityName", DBName: "city_name", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 64, Required: true},
		"CityCode":         {Name: "CityCode", DBName: "city_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"CityTownName":     {Name: "CityName", DBName: "city_town_name", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 128, Required: true},
		"CityType":         {Name: "CityType", DBName: "city_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"CityDepartment":   {Name: "CityDepartment", DBName: "department_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
	}
)

// CityDTO representa la estructura de datos de una ciudad que se utiliza para transferir datos
// entre capas de la aplicación.
type CityDTO struct {
	CityId           uint64    `json:"-"`
	CityICode        string    `json:"icode"`
	CityCreationDate time.Time `json:"creationDate"`
	CityUpdateDate   time.Time `json:"updatDate"`
	CityName         string    `json:"name"`
	CityCode         string    `json:"code"`
	CityDepartment   uint64    `json:"-"`
}

// CityPgDB representa la estructura de la entidad City para la interacción con la base de datos PostgreSQL.
// Utiliza tipos sql.NullXXX para manejar valores nulos en la base de datos.
type CityPgDB struct {
	CityId           sql.NullInt64
	CityICode        sql.NullString
	CityCreationDate sql.NullTime
	CityUpdateDate   sql.NullTime
	CityName         sql.NullString
	CityCode         sql.NullString
	CityType         sql.NullString
	CityDepartment   sql.NullInt64
}

func (ct CityDTO) MarshalJSON() ([]byte, error) {
	type Alias CityDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		CityCreationDate string `json:"creationDate"`
		CityUpdateDate   string `json:"updateDate"`
	}{
		Alias:            (*Alias)(&ct),
		CityCreationDate: ct.CityCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		CityUpdateDate:   ct.CityUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (ct *CityDTO) UnmarshalJSON(data []byte) error {
	type Alias CityDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		CityCreationDate string `json:"creationDate"`
		CityUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(ct),
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
	ct.CityCreationDate = parse(aux.CityCreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	ct.CityUpdateDate = parse(aux.CityUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// GetCity obtiene un registro de ciudad basado en los criterios especificados.
//
// Parámetros:
// - by: estructura que contiene los criterios de filtrado (nombres, alias, valores y operador).
// - city: puntero a una estructura CityDTO donde se almacenará el resultado.
// - inTransaction: indica si la consulta se realiza dentro de una transacción.
// - module: nombre del módulo que invoca la función (para fines de logging o manejo específico).
// - connData: datos de conexión a la base de datos.
// - clientConfig: configuración del cliente de base de datos.
// - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
// - connData: datos de conexión actualizados.
// - error: error en caso de que ocurra algún fallo en la consulta o escaneo.
func GetCity(by common_controllers.By, city *CityDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción de la ruta completa de la tabla (esquema.tabla).
	var cityPath string = CityDBScheme + "." + CityDBName

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Definición de los campos que se desean seleccionar en la consulta.
	var cityFieldsSlice []string = []string{"CityId", "CityICode", "CityName", "CityCode", "CityDepartment"}
	var cityFieldsAliasSlice []string = []string{}

	// Genera la parte de la sentencia SELECT con los campos indicados.
	var cityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, cityFieldsSlice, cityFieldsAliasSlice, CityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CityDBScheme, CityFieldDefinitions, true)

	// Construcción completa de la consulta SQL utilizando la configuración de los campos y filtros.
	var query string = `SELECT ` + cityFieldsStr +
		` FROM ` + cityPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, CityDBName, by.AttrsName, []string{}, []string{}, by.Operator, CityDBScheme, CityFieldDefinitions, true)

	// Ejecuta la consulta en modo QueryRow para obtener un único registro.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Escanea los resultados y los asigna a la estructura CityDTO.
	persistenceCtrl.Scan(&city.CityId, &city.CityICode, &city.CityName, &city.CityCode, &city.CityDepartment)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetCities obtiene múltiples registros de ciudades que cumplen con los criterios especificados.
//
// Parámetros:
// - by: estructura que contiene los criterios de filtrado.
// - inTransaction: indica si la consulta se realiza dentro de una transacción.
// - module: nombre del módulo invocante.
// - connData: datos de conexión a la base de datos.
// - clientConfig: configuración del cliente de base de datos.
// - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
// - connData: datos de conexión actualizados.
// - []CityDTO: slice de estructuras CityDTO con los registros obtenidos.
// - error: error en caso de que ocurra algún fallo en la consulta o en el procesamiento de resultados.
func GetCities(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CityDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción de la ruta completa de la tabla.
	var cityPath string = CityDBScheme + "." + CityDBName

	var city = CityDTO{}
	var cities []CityDTO

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar en la consulta.
	var cityFieldsSlice []string = []string{"CityId", "CityICode", "CityName", "CityCode"}
	var cityFieldsAliasSlice []string = []string{}

	// Genera la parte de la sentencia SELECT con los campos definidos.
	var cityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, cityFieldsSlice, cityFieldsAliasSlice, CityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CityDBScheme, CityFieldDefinitions, true)

	// Construcción de la consulta SQL completa, incluyendo filtros y ordenamiento.
	var query string = `SELECT ` + cityFieldsStr +
		` FROM ` + cityPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, CityDBName, by.AttrsName, []string{}, []string{}, by.Operator, CityDBScheme, CityFieldDefinitions, true) +
		` ORDER BY ` + CityFieldDefinitions["CityName"].DBName

	// Ejecuta la consulta para obtener múltiples registros.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre cada registro obtenido y lo agrega al slice de ciudades.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&city.CityId, &city.CityICode, &city.CityName, &city.CityCode)
		cities = append(cities, city)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return cities, nil
}

// GetAllCities obtiene todos los registros de ciudades de la base de datos.
//
// Parámetros:
// - inTransaction: indica si la consulta se realiza dentro de una transacción.
// - module: nombre del módulo invocante.
// - connData: datos de conexión a la base de datos.
// - clientConfig: configuración del cliente de base de datos.
// - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
// - connData: datos de conexión actualizados.
// - []CityDTO: slice de estructuras CityDTO con todos los registros de ciudades.
// - error: error en caso de que ocurra algún fallo durante la consulta o el procesamiento.
func GetAllCities(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CityDTO, error) {
	// Definición de variables y controlador de persistencia.
	var city CityDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar.
	var cityFieldsSlice []string = []string{"CityId", "CityICode", "CityName", "CityCode"}
	var cityFieldsAliasSlice []string = []string{}

	// Genera la consulta SQL completa para seleccionar todos los registros.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, cityFieldsSlice, cityFieldsAliasSlice, CityDBName, []string{}, []string{}, []string{}, "", CityDBScheme, CityFieldDefinitions, true)

	// Ejecución de la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var citys []CityDTO
	// Itera sobre cada registro obtenido y lo agrega al slice de ciudades.
	for persistenceCtrl.Next() {
		city = CityDTO{}
		persistenceCtrl.ScanRow(&city.CityId, &city.CityICode,
			&city.CityName, &city.CityCode)

		citys = append(citys, city)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return citys, nil
}

// SearchCitiesByDeparment busca ciudades filtrando por el código del departamento y una cadena de consulta.
//
// Parámetros:
// - departmentICode: código identificador del departamento.
// - queryStr: cadena de búsqueda que se aplicará a campos específicos de la ciudad.
// - inTransaction: indica si la consulta se realiza dentro de una transacción.
// - module: nombre del módulo invocante.
// - connData: datos de conexión a la base de datos.
// - clientConfig: configuración del cliente de base de datos.
// - serverConfig: configuración del servidor de base de datos.
//
// Retorna:
// - connData: datos de conexión actualizados.
// - []CityDTO: slice de estructuras CityDTO con los registros que cumplen los criterios de búsqueda.
// - error: error en caso de fallo durante la ejecución de la consulta.
//
// Nota: Esta función utiliza variables y definiciones relacionadas con "Department" (como DepartmentDBScheme,
// DepartmentDBName y DepartmentFieldDefinitions) que se asumen definidas en otro lugar del proyecto.
func SearchCitiesByDeparment(departmentICode string, queryStr string, inTransaction bool, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]CityDTO, error) {
	// Inicialización del controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construcción de la ruta completa de la tabla City y la tabla Department.
	var cityPath string = CityDBScheme + "." + CityDBName
	var departmentPath string = DepartmentDBScheme + "." + DepartmentDBName

	var city = CityDTO{}
	var cities []CityDTO

	// Se obtiene la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Definición de los campos a seleccionar en la consulta.
	var cityFieldsSlice []string = []string{"CityId", "CityICode", "CityName", "CityCode", "CityTownCode"}
	var cityFieldsAliasSlice []string = []string{}

	// Genera la parte de la sentencia SELECT con los campos definidos.
	var cityFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, cityFieldsSlice, cityFieldsAliasSlice, CityDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, CityDBScheme, CityFieldDefinitions, true)

	// Construcción de la consulta SQL completa, incluyendo filtros para relacionar las tablas City y Department,
	// y aplicando la cadena de búsqueda a campos específicos.
	var query string = `SELECT ` + cityFieldsStr +
		` FROM ` + cityPath + `, ` + departmentPath +
		` WHERE ` + cityPath + `.` + CityFieldDefinitions["CityDepartment"].DBName + ` = '` + departmentPath + `.` + DepartmentFieldDefinitions["DeparmentId"].DBName +
		` AND ` + departmentPath + `.` + DepartmentFieldDefinitions["DeparmentICode"].DBName + ` = '$1' ` +
		` AND ( ` +
		cityPath + `.` + CityFieldDefinitions["CityName"].DBName + ` LIKE '%$2%' OR ` +
		cityPath + `.` + CityFieldDefinitions["CityCode"].DBName + ` LIKE '%$2%' OR ` +
		` )` +
		`' ORDER BY ` + cityPath + `.` + CityFieldDefinitions["CityName"].DBName + ` ASC`

	// Ejecuta la consulta con los parámetros departmentICode y queryStr.
	persistenceCtrl.Query(context.Background(), query, departmentICode, queryStr)

	// Itera sobre cada registro obtenido y lo agrega al slice de ciudades.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&city.CityId, &city.CityICode, &city.CityCreationDate, &city.CityUpdateDate, &city.CityName, &city.CityCode)
		cities = append(cities, city)
	}

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un by para todos los DAOS
	return cities, nil
}

// SetCityDefaults asigna valores por defecto a los campos de una ciudad según la acción a realizar.
//
// Parámetros:
// - city: puntero a la estructura CityDTO a la que se asignarán los valores por defecto.
// - action: cadena que especifica la acción (por ejemplo, SQL_INSERT o SQL_UPDATE) para determinar qué campos se deben inicializar.
func SetCityDefaults(city *CityDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		// En inserción, se asignan la fecha de creación y actualización actuales, y se genera un nuevo UUID.
		city.CityCreationDate = time.Now()
		city.CityUpdateDate = time.Now()
		city.CityICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		// En actualización, solo se actualiza la fecha de modificación.
		city.CityUpdateDate = time.Now()
	}
}

// PgDBToDTO convierte una estructura CityPgDB (modelo de base de datos) a una estructura CityDTO (modelo de transferencia de datos).
//
// Retorna:
// - CityDTO: estructura con los datos convertidos, validando cada campo nulo antes de asignarlo.
func (obj *CityPgDB) PgDBToDTO() CityDTO {
	var dto CityDTO

	// Conversión y asignación del ID de la ciudad.
	if obj.CityId.Valid {
		dto.CityId = uint64(obj.CityId.Int64)
	}

	// Conversión y asignación del código único de la ciudad.
	if obj.CityICode.Valid {
		dto.CityICode = obj.CityICode.String
	}

	// Conversión y asignación de la fecha de creación.
	if obj.CityCreationDate.Valid {
		dto.CityCreationDate = obj.CityCreationDate.Time
	}

	// Conversión y asignación de la fecha de actualización.
	if obj.CityUpdateDate.Valid {
		dto.CityUpdateDate = obj.CityUpdateDate.Time
	}

	// Conversión y asignación del nombre de la ciudad.
	if obj.CityName.Valid {
		dto.CityName = obj.CityName.String
	}

	// Conversión y asignación del código de la ciudad.
	if obj.CityCode.Valid {
		dto.CityCode = obj.CityCode.String
	}

	// Conversión y asignación del identificador del departamento asociado a la ciudad.
	if obj.CityDepartment.Valid {
		dto.CityDepartment = uint64(obj.CityDepartment.Int64)
	}
	return dto
}
