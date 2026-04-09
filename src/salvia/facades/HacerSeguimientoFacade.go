package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// HacerSeguimientoGET renderiza la página para ejecutar un seguimiento.
// Ruta: GET /salvia/hacer-seguimiento/:id
func HacerSeguimientoGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

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
		},
		utils.GetFullHtmlFuncMap(),
	)
}
