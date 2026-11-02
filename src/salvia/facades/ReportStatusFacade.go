// Package salvia_facades — T4/REP-02: Consulta pública de reporte.
//
// Handler de la página pública (sin sesión) donde la persona escribe el
// radicado de su reporte V2 y ve su estado. La consulta en sí es
// POST /api/v1/public/consulta-reporte (salvia/controller/report_status_controller.go).
package salvia_facades

import (
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-gonic/gin"
)

// ReportStatusGET_Public (CR-E01) renderiza la pantalla «Consulte el estado de
// su reporte». No requiere sesión: la URI contiene /public y AuthMiddleware la
// deja pasar.
func ReportStatusGET_Public(c *gin.Context) {
	common_routers.SetHeaderNoCache(c)

	common_routers.RenderTemplate(c, "ReportStatus", "salvia", "victim_contact/", salvia_config.HTML_Templates, "get_report_status", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   salvia_config.Locale["sp"]["report_status_window_title"],
			"lang":          "sp",
			"locale":        salvia_config.Locale,
			"apiPath":       salvia_config.FormPaths["sp"]["ReportStatusPOST_Public"],
			"newReportPath": salvia_config.FormPaths["sp"]["VictimContactPOST_GET_Public_V2"],
		},
		utils.GetFullHtmlFuncMap())
}
