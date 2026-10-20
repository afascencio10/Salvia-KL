package salvia_facades

import (
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RegistroV3GET renderiza la pantalla de Registro V3 (DocsMD/Screens/
// Registro V3). Sin dinamic-form (decisión A1) -- página estática que
// monta registro-v3-app.js, toda la lógica de guardado vive en
// /api/v1/victim-cases-v3/* (ver victim_case_v3_controller.go).
func RegistroV3GET(c *gin.Context) {
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	if !utils.CheckPermission(salvia_config.PermissionsByRole, "manage_registro_v3", s.CurrentRole, c) {
		return
	}

	setHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	extraTemplates := utils.GetFullHtmlTemplates()
	common_routers.RenderTemplate(c, "VictimCaseV3", "salvia", "registro_v3/", salvia_config.HTML_Templates, "registro_v3", extraTemplates, utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Registro V3",
			"currentUser": s.Names + " " + s.LastNames,
			"currentRole": s.CurrentRole,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
