// Package salvia_facades — MisRemisionesMedidasFacade: pantalla Mis Remisiones — Medidas de Emergencia.
// Solo el rol am. El id del agente sale de la sesión.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MisRemisionesMedidasGET renderiza la bandeja del agente de medidas.
func MisRemisionesMedidasGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_mis_remisiones_medidas", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	extraTemplates := append(utils.GetFullHtmlTemplates(), remisionesMedidasComponentTemplates...)

	common_facades.RenderTemplate(c, "MisRemisionesMedidas", "salvia", "mis-remisiones-medidas/", salvia_config.HTML_Templates, "mis_remisiones_medidas", extraTemplates, utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Mis Remisiones — Medidas de Emergencia",
			"currentUser": s.Names + " " + s.LastNames,
			"currentRole": s.CurrentRole,
			"agentICode":  s.UserICode,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
