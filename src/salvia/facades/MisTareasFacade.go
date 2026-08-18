// Package salvia_facades — MisTareasFacade: bandeja consolidada de tareas del agente.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-gonic/gin"
)

// MisTareasGET renderiza la pantalla "Mis Tareas".
// Los datos los carga el propio HTML vía GET /api/v1/case-tasks/mine.
// Ruta: GET /salvia/mis-tareas
// Roles permitidos: op (Agente de Seguimiento)
func MisTareasGET(c *gin.Context) {
	sc := CheckAndGetSession(c, "get_mis_tareas")
	if sc == nil {
		return
	}
	common_facades.SetHeaderNoCache(c)

	vars := BaseTemplateVars(sc, sc.Session.Lang)
	vars["windowTitle"] = "Mis Tareas"

	common_facades.RenderTemplate(
		c,
		"MisTareas",
		"salvia",
		"mis-tareas/",
		salvia_config.HTML_Templates,
		"mis_tareas",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		vars,
		utils.GetFullHtmlFuncMap(),
	)
}
