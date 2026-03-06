// Package security_routers contiene los controladores HTTP para las rutas relacionadas con la seguridad,
// específicamente para la obtención de información de ciudades y pueblos.
package security_routers

import (
	// Importación de paquetes comunes para la conexión a la base de datos y utilidades.
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"

	// Importación de paquetes de configuración, controladores y DAO específicos de seguridad.
	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"
	security_daos "bitsflow/security/dao"

	"html/template"
	"net/http"
	"strings"

	// Importación de paquetes externos para la gestión de sesiones y el framework Gin.
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// TownGET es un manejador HTTP que procesa solicitudes GET para obtener información de pueblos
// basándose en el código de una ciudad. Realiza la validación de sesión, verifica permisos de acceso,
// consulta la base de datos y, dependiendo del encabezado "Accept", responde en formato JSON o renderiza
// una plantilla HTML.
func TownGET(c *gin.Context) {
	// Recupera la sesión actual del contexto.
	session := sessions.Default(c)
	// Obtiene el identificador de sesión del usuario almacenado en la sesión.
	var sessionID string = session.Get("userData").(string)
	// Obtiene la información común de la sesión a partir del identificador.
	s, err := utils.GetCommonSession(sessionID)

	// Variables para almacenar la respuesta, el código HTTP, el nombre de la plantilla y datos de la ciudad.
	var res string
	var code int = http.StatusOK
	var tplName string
	var city security_daos.CityDTO

	// Establece los encabezados HTTP para evitar el almacenamiento en caché.
	common_routers.SetHeaderNoCache(c)

	// Obtiene parámetros de la URL.
	by := c.Param("by")
	id := c.Param("id")
	// Traduce y verifica que la ruta corresponda al caso esperado ("router_get_town_by_city_code").
	if security_config.TranslateLocale(by, s.Lang) == "router_get_town_by_city_code" {
		// Verifica que el rol actual del usuario tenga permiso para obtener información del pueblo.
		if !utils.CheckPermission(security_config.PermissionsByRole, "get_town_by_city_code", s.CurrentRole, c) {
			return
		}

		// Consulta la base de datos para obtener la información de la ciudad mediante el código.
		code, res, city = security_ctrl.GetCityByCode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		// Si la consulta fue exitosa, obtiene la lista de pueblos asociados a la ciudad.
		if code == http.StatusOK {
			code, res = security_ctrl.GetTownsByCity(city.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	} else {
		// Si el parámetro 'by' no es el esperado, no se realiza ninguna acción.
		return
	}

	// Determina el formato de respuesta basado en el encabezado "Accept" de la solicitud.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Responde en formato JSON directamente.
		c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
	} else {
		// Variables para almacenar el menú de navegación y herramientas adicionales para el usuario.
		var menu map[string][]map[string]string
		var menuTools []map[string]string = []map[string]string{}

		// Si no hubo error al obtener la sesión, asigna el menú actual de la sesión.
		if err == nil {
			menu = s.CurrentMenu
		}
		// Recupera las herramientas del menú específicas para el rol y lenguaje del usuario.
		if v, found := security_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_town_towns"]; found {
			menuTools = v
		}
		// Renderiza la plantilla HTML con la información de los pueblos, menú y herramientas.
		common_routers.RenderTemplate(c, security_daos.TownEntityName, "security", "town/", security_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"towns":              res,
				"menu":               menu,
				"menuToolsUserTable": menuTools}, template.FuncMap{})
	}
}

// TownGET_Public es un manejador HTTP que procesa solicitudes GET públicas para obtener información
// de pueblos basándose en el código de una ciudad. No requiere validación de sesión y solo devuelve
// respuestas en formato JSON; en caso de solicitar otro formato, devuelve un error de solicitud incorrecta.
func TownGET_Public(c *gin.Context) {

	// Variables para almacenar la respuesta y el código HTTP.
	var res string
	var code int = http.StatusOK

	// Variable para almacenar la información de la ciudad obtenida.
	var city security_daos.CityDTO

	// Establece los encabezados HTTP para evitar el almacenamiento en caché.
	common_routers.SetHeaderNoCache(c)

	// Obtiene parámetros de la URL.
	by := c.Param("by")
	id := c.Param("id")
	// Verifica que el parámetro 'by' traducido al español ("sp") corresponda al caso esperado.
	if security_config.TranslateLocale(by, "sp") == "router_get_town_by_city_code" {
		// Consulta la base de datos para obtener la información de la ciudad mediante el código.
		code, res, city = security_ctrl.GetCityByCode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		// Si la consulta fue exitosa, obtiene la lista de pueblos asociados a la ciudad.
		if code == http.StatusOK {
			code, res = security_ctrl.GetTownsByCity(city.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	} else {
		// Si el parámetro 'by' no es el esperado, no se realiza ninguna acción.
		return
	}

	// Determina el formato de respuesta basado en el encabezado "Accept" de la solicitud.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Responde en formato JSON directamente.
		c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
	} else {
		// Para formatos distintos a JSON, devuelve una respuesta vacía con código Bad Request.
		c.DataFromReader(http.StatusBadRequest, 0, gin.MIMEJSON, strings.NewReader(""), nil)
	}
}
