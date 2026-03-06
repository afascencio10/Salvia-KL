// Package security_routers contiene los controladores de rutas relacionadas con la seguridad,
// en particular para la gestión de ciudades basadas en departamentos.
package security_routers

import (
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"

	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"
	security_daos "bitsflow/security/dao"

	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CityGET maneja las solicitudes HTTP para obtener las ciudades de un departamento.
// Esta función verifica los permisos del usuario autenticado, obtiene la información del
// departamento y posteriormente las ciudades correspondientes. El resultado se devuelve
// en formato JSON o renderizado mediante plantilla HTML según la cabecera "Accept" de la solicitud.
func CityGET(c *gin.Context) {
	// Se obtiene la sesión actual y se extrae el ID de usuario almacenado en "userData".
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	// Se obtiene la sesión común a partir del ID, que incluye configuraciones de idioma, menú, etc.
	s, err := utils.GetCommonSession(sessionID)

	// Variables para almacenar la respuesta, código HTTP, nombre de plantilla y datos del departamento.
	var res string
	var code int = http.StatusOK
	var tplName string
	var department security_daos.DepartmentDTO

	// Se configura la cabecera de la respuesta para evitar el caché.
	common_routers.SetHeaderNoCache(c)

	// Se extraen los parámetros "by" e "id" de la URL.
	by := c.Param("by")
	id := c.Param("id")
	// Se compara el parámetro "by" con el valor esperado traducido según el idioma del usuario.
	if "router_get_city_by_department" == security_config.TranslateLocale(by, s.Lang) {
		// Se verifica que el usuario tenga permiso para la acción "get_city_by_department".
		if !utils.CheckPermission(security_config.PermissionsByRole, "get_city_by_department", s.CurrentRole, c) {
			return
		}

		// Se obtiene el departamento a partir del código proporcionado.
		code, res, department = security_ctrl.GetDepartmentByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		// Si la obtención del departamento fue exitosa, se obtienen las ciudades asociadas.
		if code == http.StatusOK {
			code, res, _ = security_ctrl.GetCitiesByDeparment(department.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	} else {
		// Si el parámetro "by" no coincide con el esperado, se retorna sin hacer nada.
		return
	}

	// Se verifica el encabezado "Accept" de la solicitud para determinar el formato de respuesta.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Si se acepta JSON, se envía la respuesta como JSON.
		c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
	} else {
		// Preparación de variables para renderizar la plantilla HTML.
		var menu map[string][]map[string]string
		// Inicialización de herramientas del menú específicas del usuario.
		var menuTools []map[string]string = []map[string]string{}

		// Si la sesión se obtuvo correctamente, se utiliza el menú actual del usuario.
		if err == nil {
			menu = s.CurrentMenu
		}
		// Se obtienen las herramientas del menú correspondientes a la acción "menu_tool_get_town_cities"
		// según el idioma y el rol del usuario.
		if v, found := security_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_town_cities"]; found {
			menuTools = v
		}
		// Se renderiza la plantilla HTML con los datos obtenidos (ciudades, menú y herramientas).
		common_routers.RenderTemplate(c, security_daos.CityEntityName, "security", "city/", security_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE, map[string]interface{}{"cities": res, "menu": menu, "menuToolsUserTable": menuTools}, template.FuncMap{})
	}
}

// CityGET_Public maneja las solicitudes públicas (sin autenticación) para obtener las ciudades de un departamento.
// La función es similar a CityGET pero no requiere verificación de permisos y utiliza el idioma español ("sp")
// de forma predeterminada. Solo se devuelve la respuesta en formato JSON; si no se solicita JSON, se responde
// con un código de error.
func CityGET_Public(c *gin.Context) {

	// Variables para almacenar la respuesta y el código HTTP.
	var res string
	var code int = http.StatusOK

	// Variable para almacenar los datos del departamento.
	var department security_daos.DepartmentDTO

	// Se configura la cabecera de la respuesta para evitar el caché.
	common_routers.SetHeaderNoCache(c)

	// Se extraen los parámetros "by" e "id" de la URL.
	by := c.Param("by")
	id := c.Param("id")
	// Se compara el parámetro "by" con el valor esperado traducido al idioma español ("sp").
	if "router_get_city_by_department" == security_config.TranslateLocale(by, "sp") {
		// Se obtiene el departamento a partir del código proporcionado.
		code, res, department = security_ctrl.GetDepartmentByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		// Si la obtención del departamento fue exitosa, se obtienen las ciudades asociadas.
		if code == http.StatusOK {
			code, res, _ = security_ctrl.GetCitiesByDeparment(department.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	} else {
		// Si el parámetro "by" no coincide con el esperado, se retorna sin hacer nada.
		return
	}

	// Se verifica el encabezado "Accept" de la solicitud para determinar el formato de respuesta.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Si se acepta JSON, se envía la respuesta en formato JSON.
		c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
	} else {
		// Si no se solicita JSON, se responde con un código de error (Bad Request).
		c.DataFromReader(http.StatusBadRequest, 0, gin.MIMEJSON, strings.NewReader(""), nil)
	}
}
