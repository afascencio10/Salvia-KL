package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-gonic/gin"
)

// remisionBackFallbackPath es el listado de origen según el rol:
// supervisor (sv) entra desde Historial de Remisiones; ps/ts desde Mis remisiones.
func remisionBackFallbackPath(role string) string {
	if role == "sv" {
		return "/salvia/historial-remisiones"
	}
	return "/salvia/mis-remisiones-psicosocial"
}

// RemisionPsicosocialDetalleGET renderiza la pantalla de detalle de una remisión psicosocial.
// Ruta: GET /salvia/remision-psicosocial/:id
// Acceso: roles con permiso get_remision_psicosocial_detalle (ps, ts, sv).
func RemisionPsicosocialDetalleGET(c *gin.Context) {
	sc := CheckAndGetSession(c, "get_remision_psicosocial_detalle")
	if sc == nil {
		return
	}

	common_facades.SetHeaderNoCache(c)

	remisionId := c.Param("id")
	vars := BaseTemplateVars(sc, sc.Session.Lang)
	vars["windowTitle"] = "Detalle de Remisión Psicosocial"
	vars["currentUserId"] = sc.Session.UserICode
	vars["userTeam"] = sc.Session.Team
	vars["remisionId"] = remisionId
	vars["salviaBackFallbackPath"] = remisionBackFallbackPath(sc.Session.CurrentRole)

	common_facades.RenderTemplate(
		c,
		"remision_psicosocial_detalle",
		"salvia",
		"remision-psicosocial/",
		salvia_config.HTML_Templates,
		"remision_psicosocial_detalle",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		vars,
		utils.GetFullHtmlFuncMap(),
	)
}
