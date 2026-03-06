// Package salvia_facades proporciona funciones fachada para manejar operaciones relacionadas con Salvia,
// incluyendo la carga de archivos en texto plano y la renderización de plantillas asociadas.
// Este paquete actúa como intermediario entre la capa HTTP y la lógica de controladores y DAOs.
package salvia_facades

import (
	"bitsflow/common/db"
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	security_config "bitsflow/security/config"
	"net/http"

	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoadPlainFilePOST maneja las solicitudes POST para cargar archivos en texto plano.
// Esta función realiza las siguientes operaciones:
// 1. Obtiene la sesión del usuario a través del contexto Gin.
// 2. Verifica que el usuario tenga los permisos necesarios para la operación.
// 3. Recupera el archivo cargado desde el formulario HTTP.
// 4. Dependiendo del parámetro "srv", invoca la función controladora correspondiente para procesar el archivo.
// 5. Retorna una respuesta HTTP en formato JSON con el código de estado y mensaje apropiados.
//
// Nota: Las variables globales dbClientConfig y dbServerConfig se asumen definidas en otro lugar del proyecto.
func LoadPlainFilePOST(c *gin.Context) {
	// Se obtiene la sesión actual del usuario.
	session := sessions.Default(c)
	// Se extrae el identificador de sesión del usuario almacenado en "userData".
	var sessionID string = session.Get("userData").(string)
	// Se recupera la sesión común utilizando el identificador de sesión.
	s, _ := utils.GetCommonSession(sessionID)
	// Se inicializan la variable 'code' con el estado HTTP 500 y 'res' para almacenar el mensaje de respuesta.
	var code = http.StatusInternalServerError
	var res string

	// Verifica que el rol actual del usuario tenga permiso para cargar archivos en texto plano.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "load_plain_files", s.CurrentRole, c) {
		return
	}
	// Se intenta obtener el archivo subido con la clave "file" del formulario.
	file, error := c.FormFile("file")

	if error == nil {
		// Se recuperan parámetros de la URL para determinar el servicio y, en algunos casos, la entidad.
		srv := c.Param("srv")
		entityICode := c.Param("id")
		// Dependiendo del parámetro 'srv', se llama a la función controladora correspondiente.
		if srv == "victim_service" {
			code, res = salvia_ctrl.UpdateMomentByPlainFile(file, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
		} else if srv == "entity_branch" && entityICode != "" {
			code, res = salvia_ctrl.SetUpdateEntityBranchByPlainFile(file, entityICode, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
		} else {
			// Si los parámetros son insuficientes o inválidos, se retorna un error de solicitud incorrecta.
			code = http.StatusBadRequest
			res = `{"loadPlainFiles":["Error interno: número de parámetros insuficiente"]}`
		}
	} else {
		// Si ocurre un error al obtener el archivo, se retorna un error indicando que el archivo no es válido.
		code = http.StatusBadRequest
		res = `{"loadPlainFiles":["Por favor adjunte un archivo válido"]}`
	}

	// Se escribe la respuesta HTTP con el código de estado y el contenido JSON generado.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// LoadPlainFilePOST_GET maneja las solicitudes GET para la carga de archivos en texto plano.
// Esta función realiza las siguientes operaciones:
// 1. Obtiene la sesión del usuario y verifica los permisos para la operación.
// 2. Recupera datos relacionados con las entidades que se requieren para la vista.
// 3. Configura encabezados HTTP para evitar el cacheo de la respuesta.
// 4. Renderiza la plantilla HTML correspondiente, pasando parámetros tales como título, menú, y datos de usuario.
//
// Nota: Se utiliza el paquete common_facades para renderizar la plantilla y establecer encabezados.
func LoadPlainFilePOST_GET(c *gin.Context) {
	// Se obtiene la sesión actual del usuario.
	session := sessions.Default(c)
	// Se extrae el identificador de sesión del usuario almacenado en "userData".
	var sessionID string = session.Get("userData").(string)
	// Se recupera la sesión común utilizando el identificador de sesión.
	s, err := utils.GetCommonSession(sessionID)
	// Se inicializa una lista vacía de entidades.
	var entities []salvia_daos.EntityDTO = []salvia_daos.EntityDTO{}

	// Verifica que el rol actual del usuario tenga permiso para cargar archivos en texto plano.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "load_plain_files", s.CurrentRole, c) {
		return
	}

	// Se establecen encabezados HTTP para prevenir el cacheo de la página.
	common_facades.SetHeaderNoCache(c)

	// Se declara una variable para almacenar el menú actual del usuario.
	var menu map[string][]map[string]string

	// Si la recuperación de la sesión fue exitosa, se asigna el menú actual.
	if err == nil {
		menu = s.CurrentMenu
	}

	// Se obtienen los datos de todas las entidades utilizando la función del controlador.
	_, _, entities = salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se renderiza la plantilla HTML para la carga de archivos en texto plano, pasando los parámetros necesarios.
	common_facades.RenderTemplate(
		c,
		salvia_daos.VictimCaseEntityName,
		"salvia",
		"plain_files/",
		salvia_config.HTML_Templates,
		"load_plain_files",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":    salvia_config.Locale["sp"]["load_plain_files_window_title"],
			"currentUser":    s.Names + " " + s.LastNames,
			"salviaFormPath": salvia_config.FormPaths[s.Lang]["PlainFilesPOST"],
			"locale":         salvia_config.Locale,
			"lang":           s.Lang,
			"menu":           menu,
			"entities":       entities,
			"nav_rules":      security_config.TranslateNavigationRule(s.Lang, security_config.NAVIGATION_RULES["load_plain_files"]),
		},
		utils.GetFullHtmlFuncMap())
}
