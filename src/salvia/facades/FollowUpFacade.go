// Package salvia_facades contiene las funciones que actúan como fachada para la gestión de logs de casos
// en el sistema Salvia.
package salvia_facades

import (
	// Paquetes internos del proyecto

	"bitsflow/common/db"
	"bitsflow/common/utils"
	common_facades "bitsflow/common/facades"
	salvia_daos "bitsflow/salvia/dao"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"

	// Paquetes estándar
	"bytes"
	"net/http"
	"strings"

	// Paquetes de terceros
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// FollowUpPOST maneja las solicitudes POST para crear  un seguimiento de caso.
// La función realiza lo siguiente:
//   - Obtiene y valida la sesión del usuario.
//   - Verifica que el usuario tenga permisos para establecer un log de caso.
//   - Lee el cuerpo de la solicitud.
//   - Llama al controlador para procesar y almacenar el log.
//   - Envía la respuesta al cliente en formato JSON.
func FollowUpPOST(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_follow_up", s.CurrentRole, c) {
		return
	}

	// Se obtiene el parámetro 'id' de la URL
	victimCaseICode := c.Param("id")
	if victimCaseICode != "" {
		// Se crea un buffer para leer el cuerpo de la solicitud
		buf := new(bytes.Buffer)
		// Se lee todo el contenido del cuerpo de la solicitud
		buf.ReadFrom(c.Request.Body)
		// Se llama al controlador para establecer el log de caso.
		// Se pasan: el contenido del cuerpo, el id del caso, el origen ("salvia"),
		// la sesión del usuario y las configuraciones de base de datos (dbClientConfig, dbServerConfig).
		code, res = salvia_ctrl.SetFollowUp(buf.String(), victimCaseICode, &db.ConnData{}, *s, dbClientConfig, dbServerConfig)
	}

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// FollowUpPUT maneja las solicitudes PUT para actualizar un seguimiento de caso.
// La función realiza lo siguiente:
//   - Obtiene y valida la sesión del usuario.
//   - Verifica que el usuario tenga permisos para establecer un log de caso.
//   - Lee el cuerpo de la solicitud.
//   - Llama al controlador para procesar y almacenar el log.
//   - Envía la respuesta al cliente en formato JSON.
func FollowUpPUT(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_follow_up", s.CurrentRole, c) {
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
		code, res = salvia_ctrl.UpdateFollowUpByICode(buf.String(), id, &db.ConnData{}, *s, dbClientConfig, dbServerConfig)
	}

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// FollowUpGET maneja la solicitud GET para obtener el detalle de un seguimiento.
func FollowUpGET(c *gin.Context) {
	// Se obtiene la sesión del usuario y el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)

	// Establece cabeceras para evitar cache en la respuesta.
	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if err == nil {
		menu = s.CurrentMenu
	}

	id := c.Param("id")

	common_facades.RenderTemplate(c, salvia_daos.FollowUpEntityName, "salvia", "follow_up_detail/", salvia_config.HTML_Templates, "get_follow_up_detail", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Detalle de Seguimiento",
			"currentUser": s.Names + " " + s.LastNames,
			"nav_rules":   salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_follow_up"]),
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
			"followUpId":  id,
		}, utils.GetFullHtmlFuncMap())
}

// MyFollowUpsGET maneja la solicitud GET para la vista de "Mis Seguimientos".
func MyFollowUpsGET(c *gin.Context) {
	// Se obtiene la sesión del usuario y el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)

	// Establece cabeceras para evitar cache en la respuesta.
	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if err == nil {
		menu = s.CurrentMenu
	}

	// Renderiza el template "get_my_follow_ups" definido en salvia_config.HTML_Templates
	common_facades.RenderTemplate(c, salvia_daos.FollowUpEntityName, "salvia", "my_follow_ups/", salvia_config.HTML_Templates, "get_my_follow_ups", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Mis Seguimientos",
			"currentUser": s.Names + " " + s.LastNames,
			"agentName":   s.Names + " " + s.LastNames, // Placeholder para el nombre del agente
			"nav_rules":   salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_follow_up"]),
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
