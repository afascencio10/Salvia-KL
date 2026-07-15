package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RemisionPsicosocialDetalleGET renderiza la pantalla de detalle de una remisión psicosocial.
// Ruta: GET /salvia/remision-psicosocial/:id
// Acceso: cualquier usuario autenticado (sin restricción de rol).
func RemisionPsicosocialDetalleGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionIDVal := session.Get("userData")
	if sessionIDVal == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}

	sessionID := sessionIDVal.(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}

	remisionId := c.Param("id")

	common_facades.RenderTemplate(
		c,
		"remision_psicosocial_detalle",
		"salvia",
		"remision-psicosocial/",
		salvia_config.HTML_Templates,
		"remision_psicosocial_detalle",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Detalle de Remisión Psicosocial",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"userTeam":      s.Team,
			"remisionId":    remisionId,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
