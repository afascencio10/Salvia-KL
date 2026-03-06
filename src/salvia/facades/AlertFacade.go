// Package salvia_facades provides facade functions for handling alert-related HTTP requests
// in the Salvia module. It manages permission verification, session management, data retrieval,
// and response rendering (both JSON and HTML) according to the client's request.
package salvia_facades

import (
	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AlertGET handles HTTP GET requests for retrieving alerts.
// It verifies user permissions, retrieves session data, processes URL parameters,
// and returns alerts either as a JSON response or by rendering an HTML template.
func AlertGET(c *gin.Context) {
	// Retrieve the current session and extract the sessionID from "userData".
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	// Get common session data (which includes user role, language, etc.) based on sessionID.
	s, err := utils.GetCommonSession(sessionID)

	// Define variables to hold the alerts response, HTTP status code, and template name.
	var alertsRes string
	var code int = http.StatusOK
	var tplName string

	// Set headers to disable caching for the current response.
	common_routers.SetHeaderNoCache(c)

	// Retrieve URL parameters "id" and "by" from the request.
	id := c.Param("id")
	by := c.Param("by")

	// If no specific alert id or 'by' parameter is provided, process a request for all alerts.
	if id == "" && by == "" {
		// Check if the current user has permission to retrieve all alerts.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_alerts_by_all", s.CurrentRole, c) {
			return
		}

		// Set the template name to be used for rendering.
		tplName = "get_alerts"
		// Retrieve alerts based on the user's role.
		// For role "op" (operator) and "et", retrieve alerts by case owner.
		switch s.CurrentRole {
		case "op":
			code, alertsRes = salvia_ctrl.GetAlertsByCaseOwner(s.UserICode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		case "sv":
			// For role "sv" (supervisor), retrieve all alerts.
			code, alertsRes = salvia_ctrl.GetAlertByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
		case "et":
			code, alertsRes = salvia_ctrl.GetAlertsByCaseOwner(s.UserICode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}

		// If an alert id is provided and the translated "by" parameter indicates a victim case query:
	} else if id != "" && salvia_config.TranslateLocale(by, s.Lang) == "router_get_alert_by_victim_case" {
		// Check if the user has permission to get alerts by victim case.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_alerts_by_victim_case", s.CurrentRole, c) {
			return
		}

		tplName = "get_alerts"
		// Retrieve alerts specific to a victim case identified by the provided id.
		code, alertsRes, _ = salvia_ctrl.GetAlertsByVictimCaseICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

		// If the translated "by" parameter indicates a query by town code:
	} else if salvia_config.TranslateLocale(by, s.Lang) == "router_get_alert_by_town_code" {
		// Check if the user has permission to get alerts by town code.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_alerts_by_town_code", s.CurrentRole, c) {
			return
		}
		tplName = "get_alerts"
		// Retrieve alerts based on the user's town code.
		code, alertsRes = salvia_ctrl.GetAlertsByTownCode(s.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

	}

	// Check if the request expects a JSON response based on the "Accept" header.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Write the alerts response as JSON using DataFromReader.
		c.DataFromReader(code, int64(len(alertsRes)), gin.MIMEJSON, strings.NewReader(alertsRes), nil)
	} else {
		// Prepare variables for HTML template rendering.
		var menu map[string][]map[string]string
		var menuToolsAlertsTable []map[string]string = []map[string]string{}

		// If the common session was retrieved without error, use the current menu from the session.
		if err == nil {
			menu = s.CurrentMenu
		}

		// Retrieve additional menu tools for alerts if available in the configuration.
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_alerts"]; found {
			menuToolsAlertsTable = v
		}

		// Render the HTML template with the retrieved alerts and additional context.
		// The context includes window title, current user info, navigation rules, locale settings,
		// alerts data, menu configuration, alert type labels, and form paths.
		common_routers.RenderTemplate(
			c,
			salvia_daos.VictimContactEntityName,
			"salvia",
			"alert/",
			salvia_config.HTML_Templates,
			tplName,
			utils.GetFullHtmlTemplates(),
			utils.DEFAULT_VIEW,
			utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":          salvia_config.Locale["sp"]["get_alert_window_title"],
				"currentUser":          s.Names + " " + s.LastNames,
				"nav_rules":            salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_alerts"]),
				"locale":               salvia_config.Locale,
				"lang":                 s.Lang,
				"alerts":               alertsRes,
				"menu":                 menu,
				"alertType":            salvia_config.ALERT_TYPE_LOCALE[s.Lang],
				"salviaAlertFormPath":  salvia_config.FormPaths[s.Lang]["AlertGET"],
				"salviaFormPath":       salvia_config.FormPaths[s.Lang]["VictimCaseGET"],
				"menuToolsAlertsTable": menuToolsAlertsTable,
			},
			utils.GetFullHtmlFuncMap())
	}
}
