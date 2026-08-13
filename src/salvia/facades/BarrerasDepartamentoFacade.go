package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// BarrerasDepartamentoGET renderiza la pantalla de Barreras Departamento.
// Ruta: GET /salvia/barreras-departamento
// Roles permitidos: en (Enlace Territorial)
func BarrerasDepartamentoGET(c *gin.Context) {
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

	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_barreras_departamento", s.CurrentRole, c) {
		return
	}

	common_facades.RenderTemplate(
		c,
		"barreras_departamento",
		"salvia",
		"barriers/",
		salvia_config.HTML_Templates,
		"barreras_departamento",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Barreras Departamento",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
