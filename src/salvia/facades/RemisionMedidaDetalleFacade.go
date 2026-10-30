package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-gonic/gin"
)

// remisionMedidaBackFallbackPath es el listado de origen según el rol.
// am y ro entran desde su bandeja. El resto, desde Remisiones Medidas.
func remisionMedidaBackFallbackPath(role string) string {
	if role == "am" || role == "ro" {
		return "/salvia/mis-remisiones-medidas"
	}
	return "/salvia/remisiones-medidas"
}

// RemisionMedidaDetalleGET renderiza el detalle de una remisión de medidas.
// Ruta: GET /salvia/remision-medida/:id
// Acceso: get_remision_medida_detalle (sm, em, sc, cr, sv, am, ro).
func RemisionMedidaDetalleGET(c *gin.Context) {
	sc := CheckAndGetSession(c, "get_remision_medida_detalle")
	if sc == nil {
		return
	}

	common_facades.SetHeaderNoCache(c)

	vars := BaseTemplateVars(sc, sc.Session.Lang)
	vars["windowTitle"] = "Detalle de Remisión — Medida de Emergencia"
	vars["currentUserId"] = sc.Session.UserICode
	vars["remisionId"] = c.Param("id")
	vars["salviaBackFallbackPath"] = remisionMedidaBackFallbackPath(sc.Session.CurrentRole)

	common_facades.RenderTemplate(
		c,
		"remision_medida_detalle",
		"salvia",
		"remision-medida/",
		salvia_config.HTML_Templates,
		"remision_medida_detalle",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		vars,
		utils.GetFullHtmlFuncMap(),
	)
}
