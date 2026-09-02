// Package salvia_facades — CasosEntidadDetalleFacade.go
// Sirve el HTML shell de la pantalla de detalle de un caso-entidad (rol et).
// Los datos los carga el propio HTML via fetch a /api/v1/entity-cases/:entityCaseId/detail.
// RUTA: GET /salvia/casos-entidad/:entityCaseId/detalle
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CasosEntidadDetalleGET verifica sesión y sirve el HTML shell del detalle caso-entidad.
func CasosEntidadDetalleGET(c *gin.Context) {
	sc := CheckAndGetSession(c, "get_casos_entidad")
	if sc == nil {
		return
	}
	common_facades.SetHeaderNoCache(c)

	entityCaseId := c.Param("entityCaseId")
	if entityCaseId == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/salvia/casos-entidad")
		return
	}

	vars := BaseTemplateVars(sc, sc.Session.Lang)
	vars["windowTitle"] = "Detalle Caso Entidad"
	vars["entityCaseId"] = entityCaseId
	vars["entityBranchId"] = sc.Session.EntityBranchId

	common_facades.RenderTemplate(
		c,
		"CasosEntidadDetalle",
		"salvia",
		"entity_case/",
		salvia_config.HTML_Templates,
		"caso_entidad_detalle",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		vars,
		utils.GetFullHtmlFuncMap(),
	)
}
