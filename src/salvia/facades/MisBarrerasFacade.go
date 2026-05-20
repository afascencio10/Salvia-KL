package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// misBarrerasModalTemplates lista los partials de modales de Mis Barreras.
var misBarrerasModalTemplates = []string{
	"frontend/html/salvia/barriers/modal_gestionar.html",
}

// MisBarrerasGET renderiza la pantalla de Mis Barreras.
// Ruta: GET /salvia/mis-barreras
// Acceso: cualquier usuario autenticado (sin restricción de rol por ahora).
func MisBarrerasGET(c *gin.Context) {
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

	extraTemplates := append(utils.GetFullHtmlTemplates(), misBarrerasModalTemplates...)

	common_facades.RenderTemplate(
		c,
		"mis_barreras",
		"salvia",
		"barriers/",
		salvia_config.HTML_Templates,
		"mis_barreras",
		extraTemplates,
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Mis Barreras",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
