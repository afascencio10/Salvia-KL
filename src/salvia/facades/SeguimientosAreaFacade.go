package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ══════════════════════════════════════════════════════════════════════════════
// MODO DESARROLLO: Cambiar a false para exigir sesión y permisos en producción.
// Cuando es true, la página se renderiza sin login y muestra TODOS los casos.
// ══════════════════════════════════════════════════════════════════════════════
var seguimientosAreaDevMode = true

// SeguimientosAreaGET renderiza la página de Seguimientos del Área para supervisores.
// Ruta: GET /salvia/seguimientos/area
func SeguimientosAreaGET(c *gin.Context) {

	// ── Modo desarrollo: sin sesión, sin permisos, trae todo ─────────────────
	if seguimientosAreaDevMode {
		common_facades.RenderTemplate(
			c,
			"seguimientos_area",
			"salvia",
			"follow_up_v2/",
			salvia_config.HTML_Templates,
			"seguimientos_area",
			utils.GetFullHtmlTemplates(),
			utils.DEFAULT_VIEW,
			utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle": "Seguimientos del Área (DEV)",
				"currentUser": "Desarrollador Local",
				"locale":      salvia_config.Locale,
				"lang":        "sp",
				"UserId":      "dev-user",
				"IsAdmin":     true,
				"TeamId":      "", // vacío = trae todos los equipos
			},
			utils.GetFullHtmlFuncMap(),
		)
		return
	}

	// ── Modo producción: requiere sesión y permiso ───────────────────────────
	session := sessions.Default(c)
	userData := session.Get("userData")
	if userData == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		c.Abort()
		return
	}

	sessionID := userData.(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		c.Abort()
		return
	}

	// Validar permiso: solo sv y ad
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_seguimientos_area", s.CurrentRole, c) {
		return
	}

	common_facades.RenderTemplate(
		c,
		"seguimientos_area",
		"salvia",
		"follow_up_v2/",
		salvia_config.HTML_Templates,
		"seguimientos_area",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Seguimientos del Área",
			"currentUser": s.Names + " " + s.LastNames,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"UserId":      s.UserICode,
			"IsAdmin":     s.CurrentRole == "ad",
			"TeamId":      s.Team,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
