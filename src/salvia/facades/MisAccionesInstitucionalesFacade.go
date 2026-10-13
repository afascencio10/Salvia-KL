// Package salvia_facades — MisAccionesInstitucionalesFacade: pantalla
// "Mis Acciones Institucionales" (solo rol et).
//
// RUTA: GET /salvia/mis-acciones-institucionales
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MisAccionesInstitucionalesGET renderiza la pantalla donde la entidad consulta y
// gestiona las acciones institucionales que se le solicitaron por cada caso.
//
// Mismo patrón que CasosEntidadGET: la sede (entityBranchId) se inyecta desde la
// sesión y el frontend nunca la envía, porque es lo que acota la pantalla a las
// acciones de esta entidad.
func MisAccionesInstitucionalesGET(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_mis_acciones_institucionales", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	common_facades.RenderTemplate(c, "MisAccionesInstitucionales", "salvia", "entity_case/", salvia_config.HTML_Templates, "mis_acciones_institucionales", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":    "Mis Acciones Institucionales",
			"currentUser":    s.Names + " " + s.LastNames,
			"currentRole":    s.CurrentRole,
			"entityBranchId": s.EntityBranchId,
			"locale":         salvia_config.Locale,
			"lang":           s.Lang,
			"menu":           menu,
		}, utils.GetFullHtmlFuncMap())
}
