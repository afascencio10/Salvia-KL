// Package salvia_facades — KpiFacade: endpoints del embudo analítico.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// KpiFunnelGET — endpoint JSON con todas las métricas del embudo.
// Query params: ?from=YYYY-MM-DD&to=YYYY-MM-DD
func KpiFunnelGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, _ := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_kpis", s.CurrentRole, c) {
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	code, res := salvia_ctrl.GetFunnelKPIs(from, to, &db.ConnData{}, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// KpiDashboardGET — renderiza la pantalla del dashboard KPIs.
func KpiDashboardGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, _ := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_kpis", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	common_facades.RenderTemplate(c, "Kpi", "salvia", "kpis_dashboard/", salvia_config.HTML_Templates, "kpis_dashboard", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Dashboard KPIs",
			"currentUser": s.Names + " " + s.LastNames,
			"currentRole": s.CurrentRole,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
