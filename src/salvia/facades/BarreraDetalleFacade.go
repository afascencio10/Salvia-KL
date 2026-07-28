package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// BarreraDetalleGET renderiza la pantalla de detalle interno de una barrera.
// Ruta: GET /salvia/barreras/:id
// Acceso: cualquier usuario autenticado (sin restricción de rol por ahora).
func BarreraDetalleGET(c *gin.Context) {
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

	barrierICode := c.Param("id")

	common_facades.RenderTemplate(
		c,
		"barrera_detalle",
		"salvia",
		"barriers/",
		salvia_config.HTML_Templates,
		"barrera_detalle",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Detalle de Barrera",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"currentUserAssignedDepartment": s.AssignedDepartmentID,
			"barrierICode":  barrierICode,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
