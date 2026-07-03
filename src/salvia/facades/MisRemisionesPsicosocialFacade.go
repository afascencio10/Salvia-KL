// Package salvia_facades — MisRemisionesPsicosocialFacade: pantalla Mis remisiones Psicosocial (roles ps, ts).
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MisRemisionesPsicosocialGET renderiza remisiones del profesional logueado (directas + vía dupla).
func MisRemisionesPsicosocialGET(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_mis_remisiones_psicosocial", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	extraTemplates := append(utils.GetFullHtmlTemplates(), remisionesPsicosocialComponentTemplates...)

	common_facades.RenderTemplate(c, "MisRemisionesPsicosocial", "salvia", "mis-remisiones-psicosocial/", salvia_config.HTML_Templates, "mis_remisiones_psicosocial", extraTemplates, utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Mis remisiones Psicosocial",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
			"menu":          menu,
		}, utils.GetFullHtmlFuncMap())
}
