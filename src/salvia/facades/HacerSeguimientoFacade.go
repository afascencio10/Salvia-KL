package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// HacerSeguimientoGET renderiza la página para ejecutar un seguimiento.
// Ruta: GET /salvia/hacer-seguimiento/:id
func HacerSeguimientoGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		c.Abort()
		return
	}

	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_hacer_seguimiento", s.CurrentRole, c) {
		return
	}

	followUpID := c.Param("id")

	common_facades.RenderTemplate(
		c,
		"hacer_seguimiento",
		"salvia",
		"follow_up_v2/",
		salvia_config.HTML_Templates,
		"hacer_seguimiento",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Hacer Seguimiento",
			"currentUser": s.Names + " " + s.LastNames,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"followUpId":  followUpID,
			"userICode":   s.UserICode,
			"userRole":    s.CurrentRole,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
