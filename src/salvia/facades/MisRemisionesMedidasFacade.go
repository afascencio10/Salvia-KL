// Package salvia_facades — MisRemisionesMedidasFacade: pantalla Mis Remisiones — Medidas de Emergencia.
// Vista fija de agente (am-1). El mapeo currentRole → vista queda pendiente.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MisRemisionesMedidasGET renderiza la bandeja del agente con data quemada.
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
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
