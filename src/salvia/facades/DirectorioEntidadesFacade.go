// Package salvia_facades — DirectorioEntidadesFacade: renderiza pantalla del Directorio de Entidades.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// DirectorioEntidadesGET renderiza la pantalla del Directorio de Entidades.
func DirectorioEntidadesGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, _ := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_directorio_entidades", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	common_facades.RenderTemplate(c, "DirectorioEntidades", "salvia", "directorio_entidades/", salvia_config.HTML_Templates, "directorio_entidades", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":      "Directorio de Entidades",
			"currentUser":      s.Names + " " + s.LastNames,
			"currentUserICode": s.UserICode,
			"currentRole":      s.CurrentRole,
			"locale":           salvia_config.Locale,
			"lang":             s.Lang,
			"menu":             menu,
		}, utils.GetFullHtmlFuncMap())
}
