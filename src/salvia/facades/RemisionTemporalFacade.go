// Package salvia_facades — RemisionTemporalFacade: pantalla temporal que hospeda
// el card del flujo 3x3 de Atención Psicosocial mientras no exista la pantalla
// interna definitiva de Remisión.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RemisionTemporalGET renderiza la pantalla temporal /salvia/remision-temporal/:id
// donde :id = psychosocial_support_id. Monta <psychosocial-contact-card>.
func RemisionTemporalGET(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_psychosocial_contact_attempts", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	psicosocialID := c.Param("id")

	common_facades.RenderTemplate(
		c,
		"RemisionTemporal",
		"salvia",
		"remision-temporal/",
		salvia_config.HTML_Templates,
		"remision_temporal",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Remisión Psicosocial",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"psicosocialId": psicosocialID,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
