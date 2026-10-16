package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PrototipoSeguimientoV4GET muestra el prototipo del formulario de seguimiento V.4.
// Cualquier usuario con sesión puede entrar: no consulta permisos por rol.
// La pantalla no escribe en base de datos.
func PrototipoSeguimientoV4GET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		c.Abort()
		return
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		c.Abort()
		return
	}

	common_facades.RenderTemplate(
		c,
		"prototipo_seguimiento_v4",
		"salvia",
		"follow_up_v2/",
		salvia_config.HTML_Templates,
		"prototipo_seguimiento_v4",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Prototipo seguimiento V.4",
			"currentUser": s.Names + " " + s.LastNames,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"userICode":   s.UserICode,
			"currentRole": s.CurrentRole,
			"userRole":    s.CurrentRole,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
