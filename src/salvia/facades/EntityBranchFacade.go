package salvia_facades

import (
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"
	security_daos "bitsflow/security/dao"
	"bytes"

	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// EntityBranchGET maneja las peticiones HTTP para obtener información de sucursales de entidades.
// Dependiendo de los parámetros de la URL y los permisos del usuario, se obtienen datos en formato JSON
// o se renderiza una plantilla HTML con los resultados.
// La función realiza las siguientes operaciones:
// 1. Obtiene la sesión del usuario y extrae el ID de sesión.
// 2. Recupera la sesión común asociada al usuario.
// 3. Lee los parámetros de la URL (townCode, entityICode y by) para determinar el tipo de consulta.
// 4. Según los parámetros, verifica los permisos y obtiene la información correspondiente de las sucursales.
// 5. Organiza los datos en estructuras anidadas (cuando es requerido) y genera la respuesta en JSON.
// 6. Si la cabecera "Accept" es "application/json", retorna la respuesta JSON; de lo contrario, renderiza una plantilla HTML.
func EntityBranchGET(c *gin.Context) {
	// Obtiene la sesión actual del contexto y extrae el ID de sesión almacenado bajo "userData"
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)

	// Recupera la sesión común usando el ID de sesión obtenido
	s, err := utils.GetCommonSession(sessionID)

	// Variable que contendrá el resultado final de la respuesta (en JSON o HTML)
	var entityBranchesRes string

	// Código HTTP para la respuesta, inicialmente OK (200)
	var code int = http.StatusOK

	// Nombre de la plantilla a utilizar para la renderización (se espera asignación en otra parte si es necesario)
	var tplName string

	// Mapa anidado para agrupar las sucursales por "moment", "sector" y "entity Icode"
	var entitiesByMomentAndSector map[string]map[string]map[string][]salvia_daos.EntityBranchDTO = make(map[string]map[string]map[string][]salvia_daos.EntityBranchDTO)

	// Slice que almacenará las sucursales de entidad obtenidas
	var branches []salvia_daos.EntityBranchDTO = []salvia_daos.EntityBranchDTO{}

	// Variable para almacenar la cadena de sucursales en formato string, cuando aplique
	var branchesStr string

	// Se obtienen todos los departamentos desde el controlador de seguridad.
	_, departments := security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Establece cabeceras HTTP para prevenir cacheo de la respuesta
	common_routers.SetHeaderNoCache(c)

	// Obtiene los parámetros de la URL:
	// - townCode: Código de la ciudad o municipio.
	// - entityICode: Código de identificación de la entidad.
	// - by: Parámetro para determinar el modo de consulta, que se traduce según el idioma de la sesión.
	townCode := c.Param("townCode")
	entityICode := c.Param("entityICode")
	by := c.Param("by")
	var deparment security_daos.DepartmentDTO
	var cities []security_daos.CityDTO

	if s.CurrentRole == "do" {
		_, _, deparment = security_ctrl.GetDepartmentByTownCode(s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		_, _, cities = security_ctrl.GetCitiesByDeparment(deparment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Caso 1: Si se proporciona townCode y el parámetro 'by' coincide con "router_get_city_with_moments"
	if townCode != "" && "router_get_city_with_moments" == salvia_config.TranslateLocale(by, s.Lang) {
		// Verifica si el usuario tiene permiso para obtener sucursales con momentos
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_entity_branches_with_moments", s.CurrentRole, c) {
			return
		}
		// Llama al controlador para obtener las sucursales por townCode con información de momentos.
		// Los parámetros dbClientConfig y dbServerConfig son configuraciones de la base de datos definidas globalmente.
		_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(townCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

		// Itera sobre las sucursales obtenidas para agruparlas en el mapa anidado por "moment", "sector" y "entity Icode"
		for _, e := range branches {
			// Inicializa el mapa para el "moment" si no existe
			if entitiesByMomentAndSector[e.EntityBranchMoment] == nil {
				entitiesByMomentAndSector[e.EntityBranchMoment] = map[string]map[string][]salvia_daos.EntityBranchDTO{}
			}
			// Inicializa el mapa para el "sector" si no existe
			if entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector] == nil {
				entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector] = map[string][]salvia_daos.EntityBranchDTO{}
				// Inicializa el slice para el "entity Icode" en el sector correspondiente
				entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector][e.EntityBranchEntityICode] = []salvia_daos.EntityBranchDTO{}
			}
			// Agrega la sucursal actual al slice correspondiente en el mapa anidado
			entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector][e.EntityBranchEntityICode] = append(entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector][e.EntityBranchEntityICode], e)
		}

		// Genera una respuesta JSON exitosa a partir del mapa de sucursales agrupadas
		entityBranchesRes = utils.CommMsgGetJSONSuccess(entitiesByMomentAndSector)
	} else if entityICode != "" && townCode != "" {
		// Caso 2: Si se proporciona un entityICode (y potencialmente un townCode), se verifica otro permiso
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_entity_branches_by_towncode_and_entity", s.CurrentRole, c) {
			return
		}
		// Llama al controlador para obtener sucursales basado en townCode y entityICode.
		// Se espera que la función retorne la cadena de sucursales en formato string.
		_, branchesStr, _ = salvia_ctrl.GetEntityBranchesByTownCodeAndEntityICode(entityICode, townCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		entityBranchesRes = branchesStr
	}

	// Verifica el header "Accept" para determinar el formato de respuesta solicitado por el cliente.
	// Si se solicita "application/json", se devuelve el JSON generado.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		c.DataFromReader(code, int64(len(entityBranchesRes)), gin.MIMEJSON, strings.NewReader(entityBranchesRes), nil)
	} else {
		// Si no se solicita JSON, se prepara la respuesta en formato HTML.
		tplName = "get_entity_branches"
		// Declara variable para almacenar el menú actual de la sesión
		var menu map[string][]map[string]string
		// Declara y inicializa un slice para herramientas del menú relacionadas con las sucursales
		var menuTools []map[string]string = []map[string]string{}

		// Si no hubo error al obtener la sesión común, asigna el menú actual
		if err == nil {
			menu = s.CurrentMenu
		}
		// Intenta obtener las herramientas del menú para la acción de "get_entity_branch_by_city"
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_entity_branches"]; found {
			menuTools = v
		}
		// Renderiza la plantilla HTML utilizando los datos obtenidos, tales como:
		// - entityBranches: La respuesta generada (en formato string) con la información de las sucursales.
		// - menu: El menú actual del usuario.
		// - menuToolsEntityBranchTable: Herramientas específicas para la tabla de sucursales.
		common_routers.RenderTemplate(
			c,
			salvia_daos.EntityBranchName, // Nombre de la entidad
			"salvia",                     // Nombre de la aplicación o módulo
			"entity_branch/",             // Ruta base de las plantillas para entity_branch
			salvia_config.HTML_Templates, // Colección de plantillas HTML
			tplName,                      // Nombre específico de la plantilla (si aplica)
			utils.GetFullHtmlTemplates(), // Función para obtener todas las plantillas HTML completas
			utils.DEFAULT_VIEW,           // Vista por defecto
			utils.DEFAULT_PANIC_TEMPLATE, // Plantilla para manejo de errores (panic)
			map[string]interface{}{ // Datos a enviar a la plantilla
				"windowTitle":                salvia_config.Locale["sp"]["get_victim_case_window_title"],
				"currentUser":                s.Names + " " + s.LastNames,
				"nav_rules":                  salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_victim_case"]),
				"locale":                     salvia_config.Locale,
				"lang":                       s.Lang,
				"entityBranches":             entityBranchesRes,
				"menu":                       menu,
				"menuToolsEntityBranchTable": menuTools,
				"departments":                departments,
				"department":                 deparment.DepartmentICode,
				"cities":                     cities,
				"sectors":                    salvia_config.SECTOR[s.Lang],
				"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
				"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
				"salviaEntityFormPath":       salvia_config.FormPaths[s.Lang]["EntityGET"],
				"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
				"salviaPOSTFormPath":         salvia_config.FormPaths[s.Lang]["EntityBranchPOST"],
			},
			utils.GetFullHtmlFuncMap(), // Funciones adicionales para utilizar en la plantilla
		)
	}
}

func EntityBranchPOST(c *gin.Context) {
	// Se obtiene la sesión actual del contexto
	session := sessions.Default(c)
	// Se extrae el identificador de sesión del usuario
	var sessionID string = session.Get("userData").(string)
	// Se recupera la información común de la sesión usando el ID
	s, _ := utils.GetCommonSession(sessionID)
	// Se inicializan las variables para el código HTTP y la respuesta
	var code int = http.StatusBadRequest
	var res string

	// Se verifica si el usuario tiene permiso para establecer un log de caso.
	// La verificación se basa en la configuración de permisos por rol.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_entity_branch", s.CurrentRole, c) {
		return
	}

	// Se obtiene el parámetro 'id' de la URL

	// Se crea un buffer para leer el cuerpo de la solicitud
	buf := new(bytes.Buffer)
	// Se lee todo el contenido del cuerpo de la solicitud
	buf.ReadFrom(c.Request.Body)
	// Se llama al controlador para establecer el log de caso.
	// Se pasan: el contenido del cuerpo, el id del caso, el origen ("salvia"),
	// la sesión del usuario y las configuraciones de base de datos (dbClientConfig, dbServerConfig).
	code, res = salvia_ctrl.SetEntityBranchByUser(buf.String(), *s, dbClientConfig, dbServerConfig)

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

func EntityBranchPUT(c *gin.Context) {
	// Se obtiene la sesión actual del contexto
	session := sessions.Default(c)
	// Se extrae el identificador de sesión del usuario
	var sessionID string = session.Get("userData").(string)
	// Se recupera la información común de la sesión usando el ID
	s, _ := utils.GetCommonSession(sessionID)
	// Se inicializan las variables para el código HTTP y la respuesta
	var code int = http.StatusBadRequest
	var res string

	// Se verifica si el usuario tiene permiso para establecer un log de caso.
	// La verificación se basa en la configuración de permisos por rol.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_entity_branch", s.CurrentRole, c) {
		return
	}

	// Se obtiene el parámetro 'id' de la URL

	// Se crea un buffer para leer el cuerpo de la solicitud
	buf := new(bytes.Buffer)
	// Se lee todo el contenido del cuerpo de la solicitud
	buf.ReadFrom(c.Request.Body)
	// Se llama al controlador para establecer el log de caso.
	// Se pasan: el contenido del cuerpo, el id del caso, el origen ("salvia"),
	// la sesión del usuario y las configuraciones de base de datos (dbClientConfig, dbServerConfig).
	code, res = salvia_ctrl.UpdateEntityBranchByUser(buf.String(), *s, dbClientConfig, dbServerConfig)

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
