package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RegistroCasoV2GET sirve la pantalla "Registro de Caso" basada en dinamic-form.
// Ver DocsMD/Screens/Registro de Caso V2/registro-caso-v2-interface.md
func RegistroCasoV2GET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		c.Abort()
		return
	}

	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_victim_case", s.CurrentRole, c) {
		return
	}

	common_facades.RenderTemplate(
		c,
		"set_victim_case_v2",
		"salvia",
		"victim_case/",
		salvia_config.HTML_Templates,
		"set_victim_case_v2",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Registro de Caso",
			"currentUser": s.Names + " " + s.LastNames,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"userICode":   s.UserICode,
			"userRole":    s.CurrentRole,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
