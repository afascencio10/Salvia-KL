// Package security_routers contiene los controladores HTTP para las operaciones relacionadas
// con el manejo de usuarios generales en el módulo de seguridad.
package security_routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	common_config "bitsflow/common/config"
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"

	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"
	security_daos "bitsflow/security/dao"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GeneralUserPOST maneja la petición HTTP POST para crear o actualizar un usuario general.
// Recupera la sesión del usuario, verifica los permisos necesarios y procesa el cuerpo
// de la petición para registrar o modificar la información del usuario a través del controlador correspondiente.
func GeneralUserPOST(c *gin.Context) {
	// Se obtiene la sesión actual y se extrae el ID de sesión almacenado en "userData".
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	// Se recupera la sesión común usando el identificador.
	s, _ := utils.GetCommonSession(sessionID)
	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	// Verificamos que el usuario tenga permiso para "set_general_user".
	if !utils.CheckPermission(security_config.PermissionsByRole, "set_general_user", s.CurrentRole, c) {
		return
	}

	// Se lee el cuerpo de la petición.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Se llama al controlador para procesar la creación o actualización del usuario general.
	code, res := security_ctrl.SetGeneralUser(buf.String(), *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
	// Se envía la respuesta en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// GeneralUserLOGIN_POST procesa la petición HTTP POST para el login, reseteo o recuperación de contraseña
// del usuario general. Dependiendo del parámetro "by" de la URL, se delega la acción correspondiente en el controlador.
func GeneralUserLOGIN_POST(c *gin.Context) {
	// Se lee el cuerpo de la petición.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	var navStr string
	var res string
	var code int
	var usr security_daos.GeneralUserDTO

	// Se obtiene el parámetro "by" para determinar el tipo de acción (login, forgot o reset).
	by := c.Param("by")

	if by == "forgot" {
		// Procesa la solicitud para recuperar la contraseña.
		code, navStr, _ = security_ctrl.SetResetPassword(buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)
	} else if by == "reset" {
		// Procesa la solicitud para actualizar la contraseña tras un reseteo.
		code, navStr = security_ctrl.UpdateGeneralUserByReset(buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)
	} else if by == "" {
		// Procesa la solicitud de login normal.
		var checkCaptcha bool = true
		if v, found := c.Request.Header["User-Agent"]; found && v[0] == "flutter-client" {
			checkCaptcha = false
		}

		//checkCaptcha = false

		code, res, usr = security_ctrl.LoginGeneralUser(buf.String(), checkCaptcha, &db.ConnData{}, dbClientConfig, dbServerConfig)
		navStr = res
		if code == http.StatusOK {

			// Si el usuario ya tiene una sesión activa, se cierra la sesión anterior.
			s, error := utils.GetCommonSessionByLogin(usr.GeneralUserLogin)
			if error == nil {
				utils.RemoveCommonSession(s.SessionID, true)
			}

			// Se toma el primer rol asignado al usuario; en el futuro, se permitirá elegir el rol.
			var currentRole string
			if len(usr.GeneralUserRoleCodes) > 0 {
				currentRole = usr.GeneralUserRoleCodes[0]
				// Se genera un nuevo ID de sesión.
				var sessionID = utils.GetUUID()

				// Se almacena el nuevo ID de sesión en la sesión actual.
				session := sessions.Default(c)
				session.Set("userData", sessionID)
				session.Save()

				// Se inicializa el menú que será mostrado al usuario, fusionando opciones según sus roles.
				var menu map[string][]map[string]string = map[string][]map[string]string{}

				for _, r := range usr.GeneralUserRoleCodes {
					if v, found := security_config.Menu[usr.GeneralUserLanguage][r]; found {
						utils.MergeMaps(menu, v)
					}
				}

				// Se determina la regla de navegación (navigation rules) según el rol actual.
				switch currentRole {
				case "ad":
					navMap := security_config.TranslateNavigationRule(usr.GeneralUserLanguage, security_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := security_config.Menu[usr.GeneralUserLanguage]["ad"]; found {
						utils.MergeMaps(menu, v)
					}
				case "op":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["op"]; found {
						utils.MergeMaps(menu, v)
					}
				case "ro":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["ro"]; found {
						utils.MergeMaps(menu, v)
					}
				case "do":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["do"]; found {
						utils.MergeMaps(menu, v)
					}
				case "no":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["no"]; found {
						utils.MergeMaps(menu, v)
					}
				case "et":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["et"]; found {
						utils.MergeMaps(menu, v)
					}
				case "us":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["us"]; found {
						utils.MergeMaps(menu, v)
					}
				case "sv":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["sv"]; found {
						utils.MergeMaps(menu, v)
					}

				case "fo":
					navMap := salvia_config.TranslateNavigationRule(usr.GeneralUserLanguage, salvia_config.NAVIGATION_RULES["login"])
					navMap["role"] = currentRole
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
					if v, found := salvia_config.Menu[usr.GeneralUserLanguage]["fo"]; found {
						utils.MergeMaps(menu, v)
					}

				case "an":
					// Agente de Notificaciones: redirigir directamente a la pantalla de notificaciones.
					navMap := map[string]string{
						"login": "/salvia/notificaciones",
						"role":  currentRole,
					}
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)

				case "en":
					// Agente de Mis Barreras: redirigir directamente a la pantalla de mis barreras.
					navMap := map[string]string{
						"login": "/salvia/mis-barreras",
						"role":  currentRole,
					}
					navByte, _ := json.Marshal(navMap)
					navStr = string(navByte)
				}

				// Se crea una sesión común con la información del usuario autenticado.
				var cs utils.CommonSession = utils.CommonSession{
					UserICode:        usr.GeneralUserICode,
					Roles:            usr.GeneralUserRoleCodes,
					CurrentRole:      currentRole,
					Lang:             usr.GeneralUserLanguage,
					CurrentMenu:      menu,
					Names:            usr.GeneralUserGeneralUserProfile.GeneralUserProfileNames,
					LastNames:        usr.GeneralUserGeneralUserProfile.GeneralUserProfileLastNames,
					TownCode:         usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownCode,
					TownICode:        usr.GeneralUserGeneralUserProfile.GeneralUserProfileTown.TownICode,
					EntityBrandICode: usr.GeneralUserGeneralUserProfile.GeneralUserProfileEntityBranchSelected,
					UserLogin:        usr.GeneralUserLogin,
					SessionID:        sessionID,
					Team:             usr.GeneralUserTeam,
				}
				// Se almacena la sesión común.
				utils.AddCommonSession(sessionID, &cs)
			} else {
				code = http.StatusForbidden
			}
		}
	}

	// Se envía la respuesta con el código HTTP correspondiente y el contenido JSON.
	c.DataFromReader(code, int64(len(navStr)), gin.MIMEJSON, strings.NewReader(navStr), nil)
}

// GeneralUserLOGOUT_POST maneja la petición HTTP POST para el logout del usuario general.
// Elimina la información de la sesión y remueve la sesión común asociada.
func GeneralUserLOGOUT_POST(c *gin.Context) {
	// Se obtiene la sesión actual y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	// Se elimina la información de "userData" y se guarda la sesión actualizada.
	session.Delete("userData")
	session.Save()

	// Se elimina la sesión común del sistema.
	utils.RemoveCommonSession(sessionID, true)

	// Se envía una respuesta vacía con código HTTP 200.
	c.DataFromReader(http.StatusOK, int64(len("")), gin.MIMEJSON, strings.NewReader(""), nil)
}

// GeneralUserLOGIN_GET maneja la petición HTTP GET para mostrar la página de login del usuario general.
// Prepara los datos necesarios, como el captcha y parámetros opcionales para el reseteo de contraseña, y renderiza la plantilla HTML.
func GeneralUserLOGIN_GET(c *gin.Context) {
	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	// Se obtiene el parámetro opcional para reseteo de contraseña.
	var resetPswd string = c.Param("id")
	// Se genera un nuevo ID para el captcha.
	captchaID := captcha.New()

	// Se renderiza la plantilla de login con los datos necesarios.
	common_routers.RenderTemplate(c, security_daos.GeneralUserEntityName, "security", "", security_config.HTML_Templates, "login", utils.GetLoginHtmlTemplates(), utils.LOGIN_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": security_config.Locale["sp"]["login_window_title"],
			"language":    common_config.LANGUAGE,
			"locale":      security_config.Locale,
			"lang":        "sp",
			"resetPswd":   resetPswd,
			"CaptchaID":   captchaID,
			"nav_rules":   "",
		}, utils.GetLoginHtmlFuncMap())
}

// GeneralUserLOGIN_CAPTCHA genera y envía la imagen CAPTCHA para la autenticación.
// Si el ID del captcha no es válido, se responde con un error HTTP 400.
func GeneralUserLOGIN_CAPTCHA(c *gin.Context) {
	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	// Se obtiene el ID del captcha desde los parámetros de la URL.
	captchaID := c.Param("id")
	media := c.Param("media")
	// Si el captcha no existe o no se puede recargar, se retorna un error.
	if captchaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CAPTCHA inválido"})
		return
	}
	if media == "img" {
		c.Header("Content-Type", "image/png")
		captcha.WriteImage(c.Writer, captchaID, 240, 80)
	} else {
		if media == "audio" {
			c.Header("Content-Type", "audio/wav")
			captcha.WriteAudio(c.Writer, captchaID, "en")
		}
	}
}

// GeneralUserGET maneja la petición HTTP GET para obtener la información de usuarios generales.
// Dependiendo de si se provee un identificador, se obtiene la información de un usuario específico o de todos los usuarios.
// La respuesta puede ser en formato JSON o renderizar una plantilla HTML.
func GeneralUserGET(c *gin.Context) {
	// Se obtiene la sesión y se recupera la sesión común.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)

	var res string
	var code int = http.StatusOK
	var tplName string
	var cities []security_daos.CityDTO = []security_daos.CityDTO{}

	_, _, deparment := security_ctrl.GetDepartmentByTownCode(s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

	_, _, cities = security_ctrl.GetCitiesByDeparment(deparment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	var count int

	// Se obtiene el parámetro "id" para determinar si se solicita un usuario específico.
	id := c.Param("id")
	filter := c.Param("f")
	page, _ := strconv.Atoi(c.Param("p"))

	if id != "" {
		// Se verifica el permiso para obtener la información de un usuario en particular.
		if !utils.CheckPermission(security_config.PermissionsByRole, "get_general_user", s.CurrentRole, c) {
			return
		}
		switch s.CurrentRole {
		case "ad":
			tplName = "get_general_user"
			code, res, _ = security_ctrl.GetGeneralUserByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		case "do":
			tplName = "get_general_user_do"
			code, res, _ = security_ctrl.GetDepartmentGeneralUserByICodeAndTownCode(id, s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}

	} else {
		// Se verifica el permiso para obtener la lista de usuarios.
		if !utils.CheckPermission(security_config.PermissionsByRole, "get_general_users", s.CurrentRole, c) {
			return
		}
		switch s.CurrentRole {
		case "ad":
			tplName = "get_general_users"
			switch filter {
			case "", "all":
				code, res, count = security_ctrl.GetGeneralUserByAll(page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "a", "e":
				code, res, count = security_ctrl.GetGeneralUserByStatus("e", page, &db.ConnData{}, dbClientConfig, dbServerConfig)

			case "d":
				code, res, count = security_ctrl.GetGeneralUserByStatus("d", page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}

		case "do":
			tplName = "get_general_users_do"
			code, res = security_ctrl.GetDepartmentGeneralUsersByTownCodeAndRoleCode("et", s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	}

	// Si el cliente solicita una respuesta JSON, se envía directamente; de lo contrario, se renderiza la plantilla HTML.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
	} else {
		var menu map[string][]map[string]string
		var menuTools []map[string]string = []map[string]string{}

		if err == nil {
			menu = s.CurrentMenu
		}
		if v, found := security_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_general_users"]; found {
			menuTools = v
		}
		// Se renderiza la plantilla con los datos del usuario y las opciones de menú.
		common_routers.RenderTemplate(c, security_daos.GeneralUserEntityName, "security", "general_user/", security_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":                  security_config.Locale["sp"]["get_general_user_window_title"],
				"currentUser":                  s.Names + " " + s.LastNames,
				"users":                        res,
				"numPages":                     int(math.Ceil(float64(count) / float64(security_config.NUM_ITEMS_PER_PAGE))),
				"menu":                         menu,
				"menuToolsUserTable":           menuTools,
				"lang":                         s.Lang,
				"locale":                       security_config.Locale,
				"department":                   deparment.DepartmentICode,
				"cities":                       cities,
				"docType":                      common_config.DOCUMENT_TYPE,
				"status":                       security_config.USER_STATUS,
				"nav_rules":                    security_config.TranslateNavigationRule("sp", security_config.NAVIGATION_RULES["get_general_user"]),
				"salviaPlainFilesFormPath":     salvia_config.FormPaths[s.Lang]["PlainFilesGET"],
				"salviaVictimCaseFormPath":     salvia_config.FormPaths[s.Lang]["VictimCaseGET"],
				"salviaEntityBranchesFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
				"salviaGeneralUserGETFormPath": security_config.FormPaths[s.Lang]["GeneralUserGET"],
			}, utils.GetFullHtmlFuncMap())
	}
}

// GeneralUserPUT_GET maneja la petición HTTP GET para mostrar la página de actualización (o habilitación/deshabilitación)
// de un usuario general. Obtiene la información actual del usuario, roles, departamentos y entidades, y renderiza la plantilla adecuada.
func GeneralUserPUT_GET(c *gin.Context) {
	// Se obtiene la sesión actual y la sesión común.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	var res string
	var roles []security_daos.RoleDTO
	var departments string
	var entities string

	// Se obtienen los departamentos y entidades disponibles.
	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
	_, entities, _ = salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	// Se obtienen los parámetros "id" y "by" para determinar la acción a realizar.
	id := c.Param("id")
	by := c.Param("by")
	var templateName string = "update_general_user"
	if id != "" {
		_, _, roles = security_ctrl.GetRoleByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
		if security_config.TranslateLocale(by, s.Lang) == "router_update_general_user_disable" || security_config.TranslateLocale(by, s.Lang) == "router_update_general_user_enable" {
			// Verifica que el usuario tenga permiso para deshabilitar el usuario.
			if !utils.CheckPermission(security_config.PermissionsByRole, "update_general_user_disable", s.CurrentRole, c) {
				return
			}
			templateName = "disable_general_user"
			switch s.CurrentRole {
			case "ad":
				_, res, _ = security_ctrl.GetGeneralUserByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "do":
				_, res, _ = security_ctrl.GetDepartmentGeneralUserByICodeAndTownCode(id, s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
		} else {
			// Verifica que el usuario tenga permiso para actualizar la información del usuario.
			if !utils.CheckPermission(security_config.PermissionsByRole, "update_general_user", s.CurrentRole, c) {
				return
			}
			templateName = "update_general_user"
			switch s.CurrentRole {
			case "ad":
				_, res, _ = security_ctrl.GetGeneralUserByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "do":
				_, res, _ = security_ctrl.GetDepartmentGeneralUserByICodeAndTownCode(id, s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
		}
	}

	// Se renderiza la plantilla de actualización con la información y opciones disponibles.
	common_routers.RenderTemplate(c, security_daos.GeneralUserEntityName, "security", "general_user/", security_config.HTML_Templates, templateName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":                security_config.Locale["sp"]["update_general_user_window_title"],
			"currentUser":                s.Names + " " + s.LastNames,
			"users":                      res,
			"gender":                     common_config.GENDER_IDENTITY,
			"docType":                    common_config.DOCUMENT_TYPE,
			"language":                   common_config.LANGUAGE,
			"locale":                     security_config.Locale,
			"roles":                      roles,
			"entities":                   entities,
			"departments":                departments,
			"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
			"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
			"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
			"lang":                       s.Lang,
			"nav_rules":                  security_config.TranslateNavigationRule("sp", security_config.NAVIGATION_RULES["update_general_user"]),
		}, utils.GetFullHtmlFuncMap())
}

// GeneralUserDELETE_GET maneja la petición HTTP GET para mostrar la página de eliminación de un usuario general.
// Verifica los permisos y obtiene la información del usuario a eliminar para renderizar la plantilla correspondiente.
func GeneralUserDELETE_GET(c *gin.Context) {
	// Se obtiene la sesión y la sesión común.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Se verifica que el usuario tenga permiso para eliminar usuarios generales.
	if !utils.CheckPermission(security_config.PermissionsByRole, "remove_general_user", s.CurrentRole, c) {
		return
	}

	var res string

	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	// Se obtiene el identificador del usuario a eliminar.
	id := c.Param("id")
	if id != "" {
		_, res, _ = security_ctrl.GetGeneralUserByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Se renderiza la plantilla de eliminación con la información del usuario.
	common_routers.RenderTemplate(c, security_daos.GeneralUserEntityName, "security", "general_user/", security_config.HTML_Templates, "remove_general_user", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": security_config.Locale["sp"]["remove_general_user_window_title"],
			"users":       res,
			"gender":      common_config.GENDER_IDENTITY,
			"docType":     common_config.DOCUMENT_TYPE,
			"language":    common_config.LANGUAGE,
			"locale":      security_config.Locale,
			"lang":        s.Lang,
			"currentUser": s.Names + " " + s.LastNames,
			"nav_rules":   security_config.TranslateNavigationRule("sp", security_config.NAVIGATION_RULES["remove_general_user"]),
		}, utils.GetFullHtmlFuncMap())
}

// GeneralUserPOST_GET maneja la petición HTTP GET para mostrar la página de creación de un nuevo usuario general.
// Verifica los permisos, obtiene los roles, departamentos y entidades necesarios, y renderiza la plantilla para la creación.
func GeneralUserPOST_GET(c *gin.Context) {
	// Se obtiene la sesión y la sesión común.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)
	var tplName string

	// Se verifica que el usuario tenga permiso para crear usuarios generales.
	if !utils.CheckPermission(security_config.PermissionsByRole, "set_general_user", s.CurrentRole, c) {
		return
	}

	var roles []security_daos.RoleDTO
	var departments string
	var entities string
	var cities string

	// Se obtienen los departamentos y entidades disponibles.
	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
	_, entities, _ = salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se obtienen todos los roles disponibles.
	_, _, roles = security_ctrl.GetRoles(&db.ConnData{}, dbClientConfig, dbServerConfig)

	_, cities, _ = security_ctrl.GetCitiesByTownCode(s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

	switch s.CurrentRole {
	case "ad":
		tplName = "set_general_user"
	case "do":
		tplName = "set_general_user_do"
	}

	// Se configuran las cabeceras para evitar la caché.
	common_routers.SetHeaderNoCache(c)

	// Se renderiza la plantilla para la creación de un usuario general.
	common_routers.RenderTemplate(c, security_daos.GeneralUserEntityName, "security", "general_user/", security_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":                security_config.Locale["sp"]["set_general_user_window_title"],
			"gender":                     common_config.GENDER_IDENTITY,
			"docType":                    common_config.DOCUMENT_TYPE,
			"language":                   common_config.LANGUAGE,
			"locale":                     security_config.Locale,
			"roles":                      roles,
			"lang":                       s.Lang,
			"currentUser":                s.Names + " " + s.LastNames,
			"entities":                   entities,
			"cities":                     cities,
			"departments":                departments,
			"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
			"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
			"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
			"nav_rules":                  security_config.TranslateNavigationRule("sp", security_config.NAVIGATION_RULES["set_general_user"]),
		}, utils.GetFullHtmlFuncMap())
}

// GeneralUserPUT maneja la petición HTTP PUT para actualizar la información de un usuario general.
// Dependiendo del parámetro "by" en la URL, determina si se trata de una deshabilitación, habilitación o actualización normal.
func GeneralUserPUT(c *gin.Context) {
	// Se obtiene la sesión y la sesión común.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	var res string
	var code int
	// Se lee el cuerpo de la petición que contiene los datos a actualizar.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	// Se obtienen el identificador del usuario y el parámetro "by" para determinar la acción.
	id := c.Param("id")
	by := c.Param("by")
	if by != "" && "router_update_general_user_disable" == security_config.TranslateLocale(by, s.Lang) {
		// Verifica el permiso para deshabilitar el usuario.
		if !utils.CheckPermission(security_config.PermissionsByRole, "update_general_user_disable", s.CurrentRole, c) {
			return
		}
		code, res = security_ctrl.DisableGeneralUserByICode(buf.String(), id, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
	} else if by != "" && "router_update_general_user_enable" == security_config.TranslateLocale(by, s.Lang) {
		// Verifica el permiso para habilitar el usuario.
		if !utils.CheckPermission(security_config.PermissionsByRole, "update_general_user_enable", s.CurrentRole, c) {
			return
		}
		code, res = security_ctrl.EnableGeneralUserByICode(buf.String(), id, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
	} else {
		// Verifica el permiso para actualizar la información del usuario.
		if !utils.CheckPermission(security_config.PermissionsByRole, "update_general_user", s.CurrentRole, c) {
			return
		}
		code, res = security_ctrl.UpdateGeneralUserByICode(buf.String(), id, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Se envía la respuesta en formato JSON con el código HTTP correspondiente.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// GeneralUserDELETE maneja la petición HTTP DELETE para eliminar un usuario general.
// Verifica los permisos, ejecuta la eliminación a través del controlador y devuelve la respuesta correspondiente.
func GeneralUserDELETE(c *gin.Context) {
	// Se obtiene la sesión y la sesión común.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Se verifica que el usuario tenga permiso para eliminar usuarios generales.
	if !utils.CheckPermission(security_config.PermissionsByRole, "remove_general_user", s.CurrentRole, c) {
		return
	}

	var res string
	var code int
	// Se obtiene el identificador del usuario a eliminar.
	id := c.Param("id")
	fmt.Println(id)

	// Se llama al controlador para eliminar el usuario.
	code, res = security_ctrl.RemoveGeneralUserByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se envía la respuesta en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
