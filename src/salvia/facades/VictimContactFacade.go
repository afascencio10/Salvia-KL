// Package salvia_facades contiene los manejadores (facades) para las operaciones de contacto de víctimas
// en la aplicación Salvia.
package salvia_facades

import (
	// Importaciones de configuraciones, controladores, utilidades y otros paquetes comunes y específicos de la aplicación.
	common_config "bitsflow/common/config"
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"
	"bytes"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// VictimContactPOST maneja la solicitud HTTP POST para crear un contacto de víctima.
// Realiza la verificación de permisos del usuario y, de ser correcto, envía el contenido del
// cuerpo de la solicitud al controlador correspondiente para procesar la creación del contacto.
func VictimContactPOST(c *gin.Context) {
	// Verificar permisos de usuario obteniendo la sesión y el rol actual.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	var checkCaptcha bool = true

	s, _ := utils.GetCommonSession(sessionID)

	// Si el usuario no tiene permisos para "set_victim_contact", se termina la ejecución.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_victim_contact", s.CurrentRole, c) {
		return
	}

	// Leer el cuerpo de la solicitud HTTP.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Enviar el contenido al controlador para crear el contacto de víctima.
	// Se utilizan las configuraciones de conexión a la base de datos (dbClientConfig, dbServerConfig).
	if v, found := c.Request.Header["User-Agent"]; found && v[0] == "flutter-client" {
		checkCaptcha = false
	}
	code, res := salvia_ctrl.SetVictimContact(buf.String(), checkCaptcha, &db.ConnData{}, dbClientConfig, dbServerConfig)
	// Enviar la respuesta HTTP con el código de estado, longitud del contenido y tipo MIME JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// VictimContactPUT maneja la solicitud HTTP PUT para actualizar (invalidar) un contacto de víctima.
// Realiza la verificación de permisos, extrae el identificador del contacto desde la URL y pasa
// el cuerpo de la solicitud al controlador para realizar la invalidación.
func VictimContactPUT(c *gin.Context) {
	var res string
	var code int = 400

	// Verificar permisos de usuario obteniendo la sesión y el rol actual.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Si el usuario no tiene permisos para "put_victim_contact", se termina la ejecución.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "put_victim_contact", s.CurrentRole, c) {
		return
	}

	// Leer el cuerpo de la solicitud HTTP.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Obtener el identificador del contacto de víctima desde los parámetros de la URL.
	id := c.Param("id")
	if id != "" {
		// Llamada al controlador para invalidar el contacto de víctima utilizando el código del identificador.
		code, res = salvia_ctrl.InvalidateVictimContactByICode(buf.String(), id, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Enviar la respuesta HTTP con el código de estado y el resultado generado.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// VictimContactPOST_Public maneja la solicitud HTTP POST pública para crear un contacto de víctima.
// A diferencia de la versión privada, esta función no realiza verificaciones de permisos.
func VictimContactPOST_Public(c *gin.Context) {
	var checkCaptcha bool = true
	// Leer el cuerpo de la solicitud HTTP.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Llamada al controlador para crear el contacto de víctima.
	if v, found := c.Request.Header["User-Agent"]; found && v[0] == "flutter-client" {
		checkCaptcha = false
	}
	code, res := salvia_ctrl.SetVictimContact(buf.String(), checkCaptcha, &db.ConnData{}, dbClientConfig, dbServerConfig)
	// Enviar la respuesta HTTP con el código de estado y el resultado generado.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// VictimContactGET maneja la solicitud HTTP GET para obtener información de contactos de víctimas.
// Realiza la verificación de permisos y, según la presencia de un identificador en la URL, obtiene
// un contacto específico o una lista de contactos, y además puede renderizar una plantilla HTML o
// devolver datos en formato JSON según el encabezado "Accept" de la solicitud.
func VictimContactGET(c *gin.Context) {
	// Verificar permisos de usuario obteniendo la sesión y el rol actual.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)

	// Variables para almacenar los resultados y el código de respuesta.
	var victimContactRes string
	var victimCaseRes string
	var code int = http.StatusOK
	var tplName string

	// Establecer encabezados para evitar el almacenamiento en caché de la respuesta.
	setHeaderNoCache(c)

	var count int

	// Obtener el identificador del contacto de víctima desde los parámetros de la URL.
	id := c.Param("id")
	filter := c.Param("f")
	page, _ := strconv.Atoi(c.Param("p"))

	if id != "" {
		// Verificar permiso específico para obtener un contacto de víctima.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_victim_contact", s.CurrentRole, c) {
			return
		}
		// Llamada al controlador para obtener el contacto de víctima específico.
		var vContact salvia_daos.VictimContactDTO
		code, victimContactRes, vContact = salvia_ctrl.GetVictimContactByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

		// Seleccionar la plantilla adecuada según el rol del usuario.
		switch s.CurrentRole {
		case "sv":
			if vContact.VictimContactForm1.VictimContactForm1ICode != "" {
				tplName = "get_victim_contact_sv_v1"
			} else {
				tplName = "get_victim_contact_sv_v2"
			}

		case "op":
			if vContact.VictimContactForm1.VictimContactForm1ICode != "" {
				tplName = "get_victim_contact_v1"
			} else {
				tplName = "get_victim_contact"
			}
		}

	} else {
		// Verificar permiso para obtener múltiples contactos de víctimas.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_victim_contacts", s.CurrentRole, c) {
			return
		}

		// Seleccionar la plantilla adecuada según el rol del usuario para múltiples contactos.
		switch s.CurrentRole {
		case "sv":
			tplName = "get_victim_contacts_sv"
		case "op":
			tplName = "get_victim_contacts"
		}

		// Llamada al controlador para obtener la lista de contactos sin caso asociado.
		switch filter {
		case "", "fcv":
			code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "v", &db.ConnData{}, dbClientConfig, dbServerConfig)
		case "fci":
			code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "i", &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	}

	// Revisar el encabezado "Accept" para determinar el tipo de respuesta (JSON o HTML).
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Devolver respuesta en formato JSON.
		c.DataFromReader(code, int64(len(victimContactRes)), gin.MIMEJSON, strings.NewReader(victimContactRes), nil)
	} else {
		// Preparar datos adicionales para renderizar la plantilla HTML.
		var menu map[string][]map[string]string
		var contactMenuTools []map[string]string = []map[string]string{}
		var caseMenuTools []map[string]string = []map[string]string{}

		if err == nil {
			menu = s.CurrentMenu
		}
		// Obtener herramientas de menú específicas para contactos de víctimas.
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_victim_contacts"]; found {
			contactMenuTools = v
		}

		// Obtener herramientas de menú específicas para casos de víctimas.
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_victim_cases"]; found {
			caseMenuTools = v
		}

		// Renderizar la plantilla HTML con los datos recopilados y configurados.
		common_routers.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_contact/", salvia_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":              salvia_config.Locale["sp"]["get_victim_contact_window_title"],
				"currentUser":              s.Names + " " + s.LastNames,
				"nav_rules":                salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_victim_case"]),
				"locale":                   salvia_config.Locale,
				"lang":                     s.Lang,
				"victimContacts":           victimContactRes,
				"victimCases":              victimCaseRes,
				"numPages":                 int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE))),
				"victimContactStatus":      salvia_config.VICTIM_CONTACT_STATUS[s.Lang],
				"menu":                     menu,
				"menuToolsContactTable":    contactMenuTools,
				"menuToolsCaseTable":       caseMenuTools,
				"gender":                   common_config.GENDER_IDENTITY,
				"docType":                  common_config.DOCUMENT_TYPE,
				"livingZone":               salvia_config.LIVING_ZONES,
				"genderIdentity":           common_config.GENDER_IDENTITY,
				"sexualOrientation":        salvia_config.SEXUAL_ORIENTATION,
				"origin":                   salvia_config.ORIGIN_PLACE,
				"occupation":               salvia_config.OCCUPATION,
				"ethnicGroup":              salvia_config.ETHNIC_GROUP,
				"salviaFormPath":           salvia_config.FormPaths[s.Lang]["VictimContactPUT"],
				"salviaVictimCaseFormPath": salvia_config.FormPaths[s.Lang]["VictimCaseGET"],
			}, utils.GetFullHtmlFuncMap(),
		)
	}
}

// VictimContactPOST_GET maneja la solicitud HTTP GET para renderizar el formulario de creación
// de un contacto de víctima. Verifica los permisos del usuario, establece los encabezados adecuados
// y renderiza la plantilla HTML con los datos necesarios para el formulario.
func VictimContactPOST_GET(c *gin.Context) {
	// Verificar permisos de usuario obteniendo la sesión y el rol actual.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Si el usuario no tiene permisos para "set_victim_contact", se termina la ejecución.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_victim_contact", s.CurrentRole, c) {
		return
	}

	// Establecer encabezados para evitar que la respuesta se almacene en caché.
	common_routers.SetHeaderNoCache(c)

	// Renderizar la plantilla HTML para el formulario de contacto de víctima, pasando datos de configuración y navegación.
	common_routers.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_contact/", salvia_config.HTML_Templates, "set_victim_contact", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":               salvia_config.Locale["sp"]["set_victim_contact_window_title"],
			"currentUser":               s.Names + " " + s.LastNames,
			"gender":                    common_config.GENDER_IDENTITY,
			"docType":                   common_config.DOCUMENT_TYPE,
			"livingZone":                salvia_config.LIVING_ZONES,
			"genderIdentity":            common_config.GENDER_IDENTITY,
			"sexualOrientation":         salvia_config.SEXUAL_ORIENTATION,
			"origin":                    salvia_config.ORIGIN_PLACE,
			"occupation":                salvia_config.OCCUPATION,
			"language":                  common_config.LANGUAGE,
			"ethnicGroup":               salvia_config.ETHNIC_GROUP,
			"locale":                    salvia_config.Locale,
			"lang":                      s.Lang,
			"yes_no":                    salvia_daos.VictimCaseForm2Enums["yes_no"],
			"victimCaseForm2ReportType": salvia_daos.VictimCaseForm2Enums["victim_case_form2_report_type"],
			"formPath":                  salvia_config.FormPaths[s.Lang]["VictimContactPOST_GET"],
			"nav_rules":                 salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["set_victim_contact"]),
		}, utils.GetFullHtmlFuncMap())
}

// VictimContactPOST_GET_Public maneja la solicitud HTTP GET pública para renderizar el formulario
// de creación de un contacto de víctima. A diferencia de la versión privada, esta función no verifica
// permisos y además incluye elementos públicos como captcha y datos de departamentos.
func VictimContactPOST_GET_Public(c *gin.Context) {
	// Declaración de variable para almacenar la información de departamentos.
	var departments string

	// Establecer encabezados para evitar que la respuesta se almacene en caché.
	common_routers.SetHeaderNoCache(c)

	// Obtener la lista de departamentos utilizando el controlador de seguridad.
	// Se ignora el primer valor de retorno y se asigna el resultado a 'departments'.
	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Generar un nuevo ID de captcha para la validación humana en el formulario.
	captchaID := captcha.New()

	// Revisar el encabezado "Accept" para determinar el tipo de respuesta (JSON o HTML).
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {

		type VictimContactForm2 struct {
			YesNo                     []salvia_daos.VictimCaseForm2EnumsDTO `json:"yesNo"`
			VictimCaseForm2ReportType []salvia_daos.VictimCaseForm2EnumsDTO `json:"reportType"`
			CaptchaID                 string                                `json:"captchaID"`
		}

		type VictimContactFields struct {
			Fields VictimContactForm2 `json:"fields"`
		}

		var formFields VictimContactForm2 = VictimContactForm2{YesNo: salvia_daos.VictimCaseForm2Enums["yes_no"], VictimCaseForm2ReportType: salvia_daos.VictimCaseForm2Enums["victim_case_form2_report_type"], CaptchaID: captchaID}
		var fields VictimContactFields = VictimContactFields{Fields: formFields}
		fieldsStr, _ := utils.DTOtoString(fields)
		// Devolver respuesta en formato JSON.
		c.DataFromReader(http.StatusOK, int64(len(fieldsStr)), gin.MIMEJSON, strings.NewReader(fieldsStr), nil)
	} else {
		// Renderizar la plantilla HTML para el formulario público de contacto de víctima, pasando
		// tanto datos de configuración como rutas de formularios para seguridad.
		common_routers.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_contact/", salvia_config.HTML_Templates, "set_victim_contact", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":               salvia_config.Locale["sp"]["set_victim_contact_window_title"],
				"gender":                    common_config.GENDER_IDENTITY,
				"docType":                   common_config.DOCUMENT_TYPE,
				"livingZone":                salvia_config.LIVING_ZONES,
				"genderIdentity":            common_config.GENDER_IDENTITY,
				"sexualOrientation":         salvia_config.SEXUAL_ORIENTATION,
				"origin":                    salvia_config.ORIGIN_PLACE,
				"occupation":                salvia_config.OCCUPATION,
				"language":                  common_config.LANGUAGE,
				"locale":                    salvia_config.Locale,
				"ethnicGroup":               salvia_config.ETHNIC_GROUP,
				"lang":                      "sp",
				"yes_no":                    salvia_daos.VictimCaseForm2Enums["yes_no"],
				"victimCaseForm2ReportType": salvia_daos.VictimCaseForm2Enums["victim_case_form2_report_type"],
				"CaptchaID":                 captchaID,
				"formPath":                  salvia_config.FormPaths["sp"]["VictimContactPOST_GET_Public"],
				"securityCityFormPath":      security_config.FormPaths["sp"]["CityGET_Public"],
				"securityTownFormPath":      security_config.FormPaths["sp"]["TownGET_Public"],
				"nav_rules":                 salvia_config.TranslateNavigationRule("sp", salvia_config.NAVIGATION_RULES["set_victim_contact"]),
				"departments":               departments,
			},

			utils.GetFullHtmlFuncMap())
	}

}

// setHeaderNoCache es una función auxiliar que establece los encabezados HTTP necesarios
// para evitar el almacenamiento en caché de la respuesta. Esto asegura que el cliente
// siempre reciba la información más reciente.
func setHeaderNoCache(c *gin.Context) {
	// Agregar encabezados para evitar el almacenamiento en caché.
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}
