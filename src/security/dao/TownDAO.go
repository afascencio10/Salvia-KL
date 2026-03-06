// Package security_daos contiene las implementaciones de acceso a datos (DAO)
// para las operaciones relacionadas con "Town" en el módulo de seguridad.
package security_daos

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	common_config "bitsflow/common/config"
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
)

var (
	// Nombres y esquemas para la entidad Town.
	TownEntityName string = "Town"
	TownJSONName   string = "town"
	TownDBName     string = "town"
	TownDBScheme   string = "security"

	// Definición de campos relacionados con las validaciones.
	// Se especifica para cada campo su nombre, nombre en base de datos, alias,
	// tipo de dato en el modelo, tamaño mínimo, tamaño máximo y si es requerido.
	TownFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"TownId":           {Name: "TownId", DBName: "town_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"TownICode":        {Name: "TownICode", DBName: "town_i_code", Alias: "", ModelType: "string", MinSize: 8, MaxSize: 36, Required: false},
		"TownCreationDate": {Name: "TownCreationDate", DBName: "town_creation_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"TownUpdateDate":   {Name: "TownUpdateDate", DBName: "town_update_date", Alias: "", ModelType: "datetime", MinSize: 0, MaxSize: 0, Required: false},
		"TownName":         {Name: "TownName", DBName: "town_name", Alias: "", ModelType: "string", MinSize: 6, MaxSize: 64, Required: true},
		"TownCode":         {Name: "TownCode", DBName: "town_code", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"TownType":         {Name: "TownType", DBName: "town_type", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 2, Required: true},
		"TownCity":         {Name: "TownCity", DBName: "city_id", Alias: "", ModelType: "uint", MinSize: 0, MaxSize: 0, Required: false},
		"TownLatitude":     {Name: "TownLatitude", DBName: "town_latitude", Alias: "", ModelType: "float", MinSize: -90, MaxSize: 90, Required: true},
		"TownLongitude":    {Name: "TownLatitude", DBName: "town_longitude", Alias: "", ModelType: "float", MinSize: -180, MaxSize: 180, Required: true},
	}
)

// TownDTO representa el Data Transfer Object para la entidad Town.
// Se utiliza para transportar la información entre capas sin exponer detalles internos.
type TownDTO struct {
	TownId           uint64    `json:"-"`
	TownICode        string    `json:"icode"`
	TownCreationDate time.Time `json:"creationDate"`
	TownUpdateDate   time.Time `json:"updateDate"`
	TownName         string    `json:"name"`
	TownCode         string    `json:"code"`
	TownType         string    `json:"type"`
	TownLatitude     string    `json:"latitude"`
	TownLongitude    string    `json:"longitude"`
	TownCity         uint64    `json:"-"`
}

// TownPgDB representa la estructura que mapea los campos de la entidad Town
// en la base de datos PostgreSQL, utilizando tipos que permiten valores nulos.
type TownPgDB struct {
	TownId           sql.NullInt64
	TownICode        sql.NullString
	TownCreationDate sql.NullTime
	TownUpdateDate   sql.NullTime
	TownName         sql.NullString
	TownCode         sql.NullString
	TownType         sql.NullString
	TownCity         sql.NullInt64
}

func (tw TownDTO) MarshalJSON() ([]byte, error) {
	type Alias TownDTO // Crea un alias para evitar recursión infinita

	// Formatea la hora como desees (por ejemplo, "2006-01-02 15:04:05")

	// Crea una estructura anónima con el formato deseado
	return json.Marshal(&struct {
		*Alias
		TownCreationDate string `json:"creationDate"`
		TownUpdateDate   string `json:"updateDate"`
	}{
		Alias:            (*Alias)(&tw),
		TownCreationDate: tw.TownCreationDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
		TownUpdateDate:   tw.TownUpdateDate.Format(common_config.DateTime.DATE_TIME_FORMAT),
	})
}

func (tw *TownDTO) UnmarshalJSON(data []byte) error {
	type Alias TownDTO

	// Estructura auxiliar con fechas como string
	aux := &struct {
		*Alias
		TownCreationDate string `json:"creationDate"`
		TownUpdateDate   string `json:"updateDate"`
	}{
		Alias: (*Alias)(tw),
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
	tw.TownCreationDate = parse(aux.TownCreationDate, common_config.DateTime.DATE_TIME_FORMAT)

	tw.TownUpdateDate = parse(aux.TownUpdateDate, common_config.DateTime.DATE_TIME_FORMAT)

	return nil
}

// GetTown obtiene un registro único de Town en base a criterios de búsqueda.
// Recibe un objeto 'By' que contiene los atributos de búsqueda, un puntero a TownDTO donde se almacenará
// el resultado, y parámetros de conexión y configuración para la base de datos.
// Retorna la conexión actualizada y un error en caso de fallo.
func GetTown(by common_controllers.By, town *TownDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construye la ruta completa de la tabla (esquema + nombre de tabla).
	var townPath string = TownDBScheme + "." + TownDBName

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Verifica si hubo error al establecer la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// Construcción de la lista de campos a seleccionar.
	var townFieldsSlice []string = []string{"TownId", "TownICode", "TownCreationDate", "TownUpdateDate", "TownName", "TownCode", "TownType", "TownCity"}
	var townFieldsAliasSlice []string = []string{}

	// Genera la parte de la sentencia SQL que especifica los campos.
	var townFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, townFieldsSlice, townFieldsAliasSlice, TownDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, TownDBScheme, TownFieldDefinitions, true)

	// Construye la sentencia SELECT completa, incluyendo la cláusula WHERE.
	var query string = `SELECT ` + townFieldsStr +
		` FROM ` + townPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, TownDBName, by.AttrsName, []string{}, []string{}, by.Operator, TownDBScheme, TownFieldDefinitions, true)

	// Imprime la consulta para fines de depuración.
	fmt.Printf(query, by.AttrsValue...)
	// Ejecuta la consulta en modo QueryRow.
	persistenceCtrl.QueryRow(context.Background(), query, by.AttrsValue...)

	// Escanea el resultado de la consulta en el objeto town.
	persistenceCtrl.Scan(&town.TownId, &town.TownICode, &town.TownCreationDate, &town.TownUpdateDate,
		&town.TownName, &town.TownCode, &town.TownType, &town.TownCity)

	// Verifica si ocurrió algún error durante la ejecución o el escaneo de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	return nil
}

// GetTowns obtiene múltiples registros de Town basándose en criterios de búsqueda.
// Recibe el objeto 'By' para definir la búsqueda, parámetros de transacción, módulo, y la configuración de la conexión.
// Retorna la conexión actualizada, una lista de TownDTO y un error en caso de ocurrir.
func GetTowns(by common_controllers.By, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]TownDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construye la ruta completa de la tabla.
	var townPath string = TownDBScheme + "." + TownDBName

	var town = TownDTO{}
	var towns []TownDTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Verifica si hubo error al establecer la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Construcción de la lista de campos a seleccionar.
	var townFieldsSlice []string = []string{"TownId", "TownICode", "TownName", "TownCode", "TownType"}
	var townFieldsAliasSlice []string = []string{}

	// Genera la parte de la sentencia SQL que especifica los campos.
	var townFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, townFieldsSlice, townFieldsAliasSlice, TownDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, TownDBScheme, TownFieldDefinitions, true)

	// Construye la sentencia SELECT completa con cláusula WHERE y ordena los resultados.
	var query string = `SELECT ` + townFieldsStr +
		` FROM ` + townPath +
		common_dao.GetSQL(common_dao.SQL_SELECT_WHERE_ONLY, by.AttrsName, by.AttrsAliasName, TownDBName, by.AttrsName, []string{}, []string{}, by.Operator, TownDBScheme, TownFieldDefinitions, true) +
		` ORDER BY ` + TownFieldDefinitions["TownName"].DBName + ` ASC`

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query, by.AttrsValue...)

	// Itera sobre los resultados y los almacena en el slice de TownDTO.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&town.TownId, &town.TownICode, &town.TownName, &town.TownCode, &town.TownType)
		towns = append(towns, town)
	}

	// Verifica si ocurrió algún error durante la ejecución o el procesamiento de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un 'by' para todos los DAOS.
	return towns, nil
}

// GetAllTowns recupera todos los registros de Town de la base de datos.
// No requiere criterios de búsqueda, por lo que obtiene todos los registros.
// Retorna la conexión actualizada, una lista de TownDTO y un error en caso de ocurrir.
func GetAllTowns(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]TownDTO, error) {
	// Declaración de variables locales.
	var town TownDTO
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Verifica si hubo error al establecer la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Construcción de la lista de campos a seleccionar.
	var townFieldsSlice []string = []string{"TownId", "TownICode", "TownName", "TownCode", "TownType"}
	var townFieldsAliasSlice []string = []string{}

	// Genera la sentencia SQL completa para seleccionar todos los registros.
	var query string = common_dao.GetSQL(common_dao.SQL_SELECT, townFieldsSlice, townFieldsAliasSlice, TownDBName, []string{}, []string{}, []string{}, "", TownDBScheme, TownFieldDefinitions, true)

	// Ejecuta la consulta.
	persistenceCtrl.Query(context.Background(), query)
	var towns []TownDTO
	// Itera sobre los resultados y los agrega al slice.
	for persistenceCtrl.Next() {
		town = TownDTO{}
		persistenceCtrl.ScanRow(&town.TownId, &town.TownICode, &town.TownName, &town.TownCode, &town.TownType)
		towns = append(towns, town)
	}

	// Verifica si ocurrió algún error durante la ejecución o el procesamiento de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	return towns, nil
}

// SearchTownsByCity busca registros de Town asociados a una ciudad específica,
// utilizando el código interno de la ciudad (cityICode) y un string de consulta para filtrar por nombre o código.
// Nota: Las variables CityDBScheme, CityDBName y CityFieldDefinitions deben estar definidas en el contexto.
func SearchTownsByCity(cityICode string, queryStr string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) ([]TownDTO, error) {
	// Inicializa el controlador de persistencia.
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}
	// Construye las rutas completas de las tablas Town y City.
	var townPath string = TownDBScheme + "." + TownDBName
	var cityPath string = CityDBScheme + "." + CityDBName

	var town = TownDTO{}
	var towns []TownDTO

	// Se establece la conexión a la base de datos.
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)

	// Verifica si hubo error al establecer la conexión.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}

	// Construcción de la lista de campos a seleccionar.
	var townFieldsSlice []string = []string{"TownId", "TownICode", "TownName", "TownCode", "TownTownCode", "TownType"}
	var townFieldsAliasSlice []string = []string{}

	// Genera la parte de la sentencia SQL que especifica los campos.
	var townFieldsStr = common_dao.GetSQL(common_dao.SQL_SELECT_FIELDS_ONLY, townFieldsSlice, townFieldsAliasSlice, TownDBName, []string{}, []string{}, []string{}, common_dao.SQL_AND, TownDBScheme, TownFieldDefinitions, true)

	// Construye la consulta SQL uniendo las tablas Town y City, y aplicando filtros de búsqueda.
	var query string = `SELECT ` + townFieldsStr +
		` FROM ` + townPath + `, ` + cityPath +
		` WHERE ` + townPath + `.` + TownFieldDefinitions["TownCity"].DBName + ` = '` + cityPath + `.` + CityFieldDefinitions["CityId"].DBName +
		` AND ` + cityPath + `.` + CityFieldDefinitions["CityICode"].DBName + ` = '$1' ` +
		` AND ( ` +
		townPath + `.` + TownFieldDefinitions["TownName"].DBName + ` LIKE '%$2%' OR ` +
		townPath + `.` + TownFieldDefinitions["TownCode"].DBName + ` LIKE '%$2%' OR ` +
		` )` +
		`' ORDER BY ` + townPath + `.` + TownFieldDefinitions["TownName"].DBName + ` ASC`

	// Ejecuta la consulta pasando los parámetros cityICode y queryStr.
	persistenceCtrl.Query(context.Background(), query, cityICode, queryStr)

	// Itera sobre los resultados y los almacena en el slice de TownDTO.
	for persistenceCtrl.Next() {
		persistenceCtrl.ScanRow(&town.TownId, &town.TownICode, &town.TownName, &town.TownCode, &town.TownType)
		towns = append(towns, town)
	}

	// Verifica si ocurrió algún error durante la ejecución o el procesamiento de la consulta.
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query: ", query)
		fmt.Println("SQL Error: ", persistenceCtrl.Error)
		return nil, persistenceCtrl.Error
	}
	// TODO: aplicar esta metodología de traer 1 o más elementos con un 'by' para todos los DAOS.
	return towns, nil
}

// SetTownDefaults asigna valores por defecto a los campos de un TownDTO según la acción a realizar.
// La acción puede ser SQL_INSERT o SQL_UPDATE, y se actualizan los campos de fecha y código según corresponda.
func SetTownDefaults(town *TownDTO, action string) {
	switch action {
	case common_dao.SQL_INSERT:
		town.TownCreationDate = time.Now()
		town.TownUpdateDate = time.Now()
		town.TownICode = utils.GetUUID()
	case common_dao.SQL_UPDATE:
		town.TownUpdateDate = time.Now()
	}
}

// ToDTO convierte una instancia de TownPgDB (estructura de base de datos) en un TownDTO.
// Realiza la verificación de valores nulos y la conversión de tipos según corresponda.
func (obj *TownPgDB) ToDTO() TownDTO {
	var dto TownDTO

	if obj.TownId.Valid {
		dto.TownId = uint64(obj.TownId.Int64)
	}

	if obj.TownICode.Valid {
		dto.TownICode = obj.TownICode.String
	}

	if obj.TownCreationDate.Valid {
		dto.TownCreationDate = obj.TownCreationDate.Time
	}

	if obj.TownUpdateDate.Valid {
		dto.TownUpdateDate = obj.TownUpdateDate.Time
	}

	if obj.TownName.Valid {
		dto.TownName = obj.TownName.String
	}

	if obj.TownCode.Valid {
		dto.TownCode = obj.TownCode.String
	}

	if obj.TownType.Valid {
		dto.TownType = obj.TownType.String
	}

	if obj.TownCity.Valid {
		dto.TownCity = uint64(obj.TownCity.Int64)
	}
	return dto
}
