// Package salvia_facades contiene las funciones que actúan como fachada para la gestión de logs de casos
// en el sistema Salvia.
package salvia_facades

import (
	// Paquetes internos del proyecto
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"

	// Paquetes estándar
	"bytes"
	"net/http"
	"strings"

	// Paquetes de terceros
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CaseLogPOST maneja las solicitudes POST para crear o actualizar un log de caso.
// La función realiza lo siguiente:
//   - Obtiene y valida la sesión del usuario.
//   - Verifica que el usuario tenga permisos para establecer un log de caso.
//   - Lee el cuerpo de la solicitud.
//   - Llama al controlador para procesar y almacenar el log.
//   - Envía la respuesta al cliente en formato JSON.
func CaseLogPOST(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_case_log", s.CurrentRole, c) {
		return
	}

	// Se obtiene el parámetro 'id' de la URL
	id := c.Param("id")
	if id != "" {
		// Se crea un buffer para leer el cuerpo de la solicitud
		buf := new(bytes.Buffer)
		// Se lee todo el contenido del cuerpo de la solicitud
		buf.ReadFrom(c.Request.Body)
		// Se llama al controlador para establecer el log de caso.
		// Se pasan: el contenido del cuerpo, el id del caso, el origen ("salvia"),
		// la sesión del usuario y las configuraciones de base de datos (dbClientConfig, dbServerConfig).
		code, res = salvia_ctrl.SetCaseLog(buf.String(), id, *s, dbClientConfig, dbServerConfig)
	}

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// CaseLogGET maneja las solicitudes GET para obtener la información de un log de caso.
// La función realiza lo siguiente:
//   - Obtiene y valida la sesión del usuario.
//   - Establece cabeceras para evitar el cacheo de la respuesta.
//   - Verifica los permisos del usuario para la acción solicitada.
//   - Dependiendo del parámetro 'by', decide si se obtiene un log individual o múltiples logs.
//   - Si la cabecera "Accept" indica JSON, se devuelve la respuesta en JSON.
//     En caso contrario, se renderiza una plantilla HTML con la información del log.
func CaseLogGET(c *gin.Context) {
	// Se obtiene la sesión actual del contexto
	session := sessions.Default(c)
	// Se extrae el identificador de sesión del usuario
	var sessionID string = session.Get("userData").(string)
	// Se recupera la información común de la sesión usando el ID, junto con cualquier error ocurrido
	s, err := utils.GetCommonSession(sessionID)

	// Variables para almacenar la respuesta del log, el código HTTP y el nombre de la plantilla a usar
	var caseLogRes string
	var code int = http.StatusUnauthorized
	var tplName string

	// Se configuran las cabeceras de la respuesta para deshabilitar el cacheo
	common_routers.SetHeaderNoCache(c)

	// Se extraen los parámetros 'by' e 'id' de la URL
	by := c.Param("by")
	id := c.Param("id")
	if id != "" {
		// Si el parámetro 'by' traducido al español coincide con "router_get_case_log_by_moment",
		// se procesará como una solicitud para obtener múltiples logs por momento.
		if "router_get_case_log_by_moment" == salvia_config.TranslateLocale(by, "sp") {
			// Verifica si el usuario tiene permiso para obtener múltiples logs de caso.
			if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_case_logs", s.CurrentRole, c) {
				return
			}

			// Se define el nombre de la plantilla a utilizar para múltiples logs
			tplName = "get_case_logs"
			// Se obtiene el log de casos por momento llamando al controlador, utilizando los parámetros necesarios
			code, caseLogRes = salvia_ctrl.GetCaseLogsByMomentICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		} else {
			// Para obtener un único log de caso, se verifica que el usuario tenga el permiso correspondiente.
			if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_case_log", s.CurrentRole, c) {
				return
			}

			// Se define el nombre de la plantilla para un único log de caso.
			tplName = "get_case_log"
			// Se obtiene el log de caso individual llamando al controlador.
			// Se ignora el tercer valor devuelto.
			code, caseLogRes, _ = salvia_ctrl.GetCaseLogByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}

	}

	// Se verifica si la solicitud espera una respuesta en formato JSON
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Se envía la respuesta en formato JSON
		c.DataFromReader(code, int64(len(caseLogRes)), gin.MIMEJSON, strings.NewReader(caseLogRes), nil)
	} else {
		// Variables para almacenar el menú de navegación y las herramientas específicas de logs de caso
		var menu map[string][]map[string]string
		var caseLogsMenuTools []map[string]string = []map[string]string{}

		// Si no hubo error al obtener la sesión, se asigna el menú actual de la sesión
		if err == nil {
			menu = s.CurrentMenu
		}

		// Se intenta obtener la configuración de herramientas de menú para logs de caso,
		// según el idioma y rol del usuario
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_case_logs"]; found {
			caseLogsMenuTools = v
		}

		// Se renderiza la plantilla HTML con los datos obtenidos y configurados,
		// pasando parámetros como título de la ventana, usuario actual, reglas de navegación, etc.
		common_routers.RenderTemplate(
			c,
			salvia_daos.VictimContactEntityName, // Nombre de la entidad
			"salvia",                            // Módulo o fuente
			"victim_case/",                      // Ruta de la plantilla
			salvia_config.HTML_Templates,        // Conjunto de plantillas HTML
			tplName,                             // Nombre de la plantilla a renderizar
			utils.GetFullHtmlTemplates(),        // Funciones HTML adicionales
			utils.DEFAULT_VIEW,                  // Vista por defecto
			utils.DEFAULT_PANIC_TEMPLATE,        // Plantilla en caso de pánico
			map[string]interface{}{
				"windowTitle":           salvia_config.Locale["sp"]["get_case_log_window_title"],
				"currentUser":           s.Names + " " + s.LastNames,
				"nav_rules":             salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_victim_case"]),
				"locale":                salvia_config.Locale,
				"menu":                  menu,
				"caseLogsMenuTools":     caseLogsMenuTools,
				"lang":                  s.Lang,
				"caseLogs":              caseLogRes,
				"salviaCaseLogFormPath": salvia_config.FormPaths[s.Lang]["CaseLogReportPOST"],
			},
			utils.GetFullHtmlFuncMap(), // Funciones adicionales para la plantilla
		)
	}
}
