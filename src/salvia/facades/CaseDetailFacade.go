// Package salvia_facades ÔÇö CaseDetailFacade.go
// Sirve el HTML shell de la pantalla de detalle de caso (rol sv).
// Los datos los carga el propio HTML via fetch a /api/v1/casos/:id/detalle.
// RUTA: GET /salvia/casos/:id/detalle
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CaseDetailGET verifica sesi├│n y sirve el HTML shell.
// El HTML hace fetch a /api/v1/casos/:id/detalle para obtener los datos.
func CaseDetailGET(c *gin.Context) {
	sc := CheckAndGetSession(c, "get_case_detail_sv")
	if sc == nil {
		return
	}
	common_facades.SetHeaderNoCache(c)

	caseICode := c.Param("id")
	if caseICode == "" {
		c.Redirect(http.StatusTemporaryRedirect, salvia_config.FormPaths[sc.Session.Lang]["VictimCaseGET"])
		return
	}

	vars := BaseTemplateVars(sc, sc.Session.Lang)
	vars["windowTitle"]                 = "Detalle del Caso"
	vars["caseICode"]                   = caseICode
	vars["userRole"]                    = sc.Session.CurrentRole
	vars["userTeam"]                    = sc.Session.Team
	vars["userICode"]                   = sc.Session.UserICode
	vars["salviaVictimCaseGETFormPath"] = salvia_config.FormPaths[sc.Session.Lang]["VictimCaseGET"]

	common_facades.RenderTemplate(
		c,
		salvia_daos.VictimCaseEntityName,
		"salvia",
		"case_detail/",
		salvia_config.HTML_Templates,
		"get_case_detail_sv",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		vars,
		utils.GetFullHtmlFuncMap(),
	)
}
