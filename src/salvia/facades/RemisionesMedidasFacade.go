// Package salvia_facades — RemisionesMedidasFacade: pantalla Remisiones Medidas de Emergencia.
// Vista fija de supervisor. El mapeo currentRole → vista queda pendiente.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// remisionesMedidasComponentTemplates — partial Vue del listado compartido.
var remisionesMedidasComponentTemplates = []string{
	"frontend/html/salvia/remisiones-medidas/remisiones_medidas_component.html",
}

// RemisionesMedidasGET renderiza el listado general con data quemada (vista supervisor).
func RemisionesMedidasGET(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_remisiones_medidas", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	extraTemplates := append(utils.GetFullHtmlTemplates(), remisionesMedidasComponentTemplates...)

	common_facades.RenderTemplate(c, "RemisionesMedidas", "salvia", "remisiones-medidas/", salvia_config.HTML_Templates, "remisiones_medidas", extraTemplates, utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Remisiones Medidas de Emergencia",
			"currentUser": s.Names + " " + s.LastNames,
			"currentRole": s.CurrentRole,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
