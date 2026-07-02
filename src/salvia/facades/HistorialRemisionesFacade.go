// Package salvia_facades — HistorialRemisionesFacade: pantalla Historial de Remisiones (solo rol sv).
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// remisionesPsicosocialComponentTemplates — partial del componente Vue remisiones-psicosocial-component.
var remisionesPsicosocialComponentTemplates = []string{
	"frontend/html/salvia/remisiones-psicosocial/remisiones_psicosocial_component.html",
}

// HistorialRemisionesGET renderiza la pantalla "Historial de Remisiones" con datos mock (backend pendiente).
func HistorialRemisionesGET(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_historial_remisiones", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	extraTemplates := append(utils.GetFullHtmlTemplates(), remisionesPsicosocialComponentTemplates...)

	common_facades.RenderTemplate(c, "HistorialRemisiones", "salvia", "historial-remisiones/", salvia_config.HTML_Templates, "historial_remisiones", extraTemplates, utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Historial de Remisiones",
			"currentUser": s.Names + " " + s.LastNames,
			"currentRole": s.CurrentRole,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
