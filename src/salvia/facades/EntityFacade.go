package salvia_facades

import (
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// EntityGET maneja las peticiones HTTP para obtener información de sucursales de entidades.
// Dependiendo de los parámetros de la URL y los permisos del usuario, se obtienen datos en formato JSON
// o se renderiza una plantilla HTML con los resultados.
// La función realiza las siguientes operaciones:
// 1. Obtiene la sesión del usuario y extrae el ID de sesión.
// 2. Recupera la sesión común asociada al usuario.
// 3. Lee los parámetros de la URL (townCode, entityICode y by) para determinar el tipo de consulta.
// 4. Según los parámetros, verifica los permisos y obtiene la información correspondiente de las sucursales.
// 5. Organiza los datos en estructuras anidadas (cuando es requerido) y genera la respuesta en JSON.
// 6. Si la cabecera "Accept" es "application/json", retorna la respuesta JSON; de lo contrario, renderiza una plantilla HTML.
func EntityGET(c *gin.Context) {
	// Obtiene la sesión actual del contexto y extrae el ID de sesión almacenado bajo "userData"
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)

	// Recupera la sesión común usando el ID de sesión obtenido
	s, err := utils.GetCommonSession(sessionID)

	// Variable que contendrá el resultado final de la respuesta (en JSON o HTML)
	var entities string

	// Código HTTP para la respuesta, inicialmente OK (200)
	var code int = http.StatusOK

	// Nombre de la plantilla a utilizar para la renderización (se espera asignación en otra parte si es necesario)
	var tplName string

	// Establece cabeceras HTTP para prevenir cacheo de la respuesta
	common_routers.SetHeaderNoCache(c)

	// Obtiene los parámetros de la URL:
	// - townCode: Código de la ciudad o municipio.
	// - entityICode: Código de identificación de la entidad.
	// - by: Parámetro para determinar el modo de consulta, que se traduce según el idioma de la sesión.
	sectorCode := c.Param("sectorCode")

	// Caso 1: Si se proporciona townCode y el parámetro 'by' coincide con "router_get_city_with_moments"
	if sectorCode != "" {
		// Verifica si el usuario tiene permiso para obtener sucursales con momentos
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_entity", s.CurrentRole, c) {
			return
		}
		_, entities, _ = salvia_ctrl.GetEntitiesBySector(sectorCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

	}
	// Verifica el header "Accept" para determinar el formato de respuesta solicitado por el cliente.
	// Si se solicita "application/json", se devuelve el JSON generado.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		c.DataFromReader(code, int64(len(entities)), gin.MIMEJSON, strings.NewReader(entities), nil)
	} else {
		// Si no se solicita JSON, se prepara la respuesta en formato HTML.
		tplName = "get_entity"
		// Declara variable para almacenar el menú actual de la sesión
		var menu map[string][]map[string]string
		// Declara y inicializa un slice para herramientas del menú relacionadas con las sucursales
		var menuTools []map[string]string = []map[string]string{}

		// Si no hubo error al obtener la sesión común, asigna el menú actual
		if err == nil {
			menu = s.CurrentMenu
		}
		// Intenta obtener las herramientas del menú para la acción de "get_entity_branch_by_city"
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_entity_branch_by_city"]; found {
			menuTools = v
		}
		// Renderiza la plantilla HTML utilizando los datos obtenidos, tales como:
		// - entityBranches: La respuesta generada (en formato string) con la información de las sucursales.
		// - menu: El menú actual del usuario.
		// - menuToolsEntityTable: Herramientas específicas para la tabla de sucursales.
		common_routers.RenderTemplate(
			c,
			salvia_daos.EntityEntityName, // Nombre de la entidad,
			"salvia",                     // Nombre de la aplicación o módulo
			"entity_branch/",             // Ruta base de las plantillas para entity_branch
			salvia_config.HTML_Templates, // Colección de plantillas HTML
			tplName,                      // Nombre específico de la plantilla (si aplica)
			utils.GetFullHtmlTemplates(), // Función para obtener todas las plantillas HTML completas
			utils.DEFAULT_VIEW,           // Vista por defecto
			utils.DEFAULT_PANIC_TEMPLATE, // Plantilla para manejo de errores (panic)
			map[string]interface{}{ // Datos a enviar a la plantilla
				"windowTitle":          salvia_config.Locale["sp"]["get_victim_case_window_title"],
				"currentUser":          s.Names + " " + s.LastNames,
				"nav_rules":            salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_victim_case"]),
				"locale":               salvia_config.Locale,
				"lang":                 s.Lang,
				"entities":             entities,
				"menu":                 menu,
				"menuToolsEntityTable": menuTools,
			},
			utils.GetFullHtmlFuncMap(), // Funciones adicionales para utilizar en la plantilla
		)
	}
}
