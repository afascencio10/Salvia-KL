// Package salvia_facades provides HTTP handlers for operator assignment operations
// within the Salvia module of the BitsFlow application.
// This package includes endpoints for assigning operators either in bulk or individually,
// as well as rendering the corresponding HTML view for operator assignment.
package salvia_facades

import (
	"bitsflow/common/db"
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	security_ctrl "bitsflow/security/controllers"
	"net/http"

	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AssignOperatorsPOST handles POST requests for assigning operators.
// Depending on the "by" URL parameter, it either assigns operators in bulk ("all")
// or individually ("single").
// The function performs the following steps:
// 1. Retrieves the current session from the Gin context.
// 2. Extracts the session ID and obtains the common session data.
// 3. Checks if the current user has the necessary permission ("assign_operators").
// 4. Retrieves URL parameters ("by", "p1", and "p2") to determine the mode of assignment.
// 5. Calls the appropriate controller function to perform the assignment.
// 6. Sends back a JSON response with the resulting status and message.
//
// Note: The variables dbClientConfig and dbServerConfig (used for database configuration)
// are assumed to be defined elsewhere in the application.
func AssignOperatorsPOST(c *gin.Context) {
	// Retrieve the current session from the Gin context.
	session := sessions.Default(c)
	// Extract the session ID stored under the key "userData". It is expected to be a string.
	var sessionID string = session.Get("userData").(string)
	// Retrieve common session data using the session ID.
	s, _ := utils.GetCommonSession(sessionID)
	// Initialize response code with Internal Server Error as default.
	var code = http.StatusInternalServerError
	// Initialize an empty response string.
	var res string

	// Check if the current user has permission to assign operators.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "assign_operators", s.CurrentRole, c) {
		// If permission is not granted, exit the handler.
		return
	}

	// Retrieve URL parameters that determine the assignment mode and targets.
	by := c.Param("by")
	p1 := c.Param("p1")
	p2 := c.Param("p2")
	// If the "by" parameter is "all" and both p1 and p2 are provided,
	// perform a bulk operator assignment.
	if by == "all" && p1 != "" && p2 != "" {
		code, res = salvia_ctrl.AssignOperators(p1, p2, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
	} else if by == "single" && p1 != "" && p2 != "" {
		// If the "by" parameter is "single" and both p1 and p2 are provided,
		// perform an individual operator assignment.
		code, res = salvia_ctrl.AssignOperator(p1, p2, *s, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Write the JSON response with the appropriate HTTP status code.
	// The response is streamed from a string reader containing the result.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// AssignOperatorsPOST_GET handles GET requests to display the operator assignment interface.
// It performs the following operations:
// 1. Retrieves the current session from the Gin context and extracts session data.
// 2. Checks if the current user has the necessary permission ("assign_operators").
// 3. Fetches a list of general users with the operator role from the security controller.
// 4. Sets HTTP headers to disable caching for the response.
// 5. Prepares the menu data from the session (if available).
// 6. Renders the HTML template for the assign operators view, passing necessary data such as:
//   - Window title and current user information.
//   - Form action path for operator assignment.
//   - Locale and language settings.
//   - Navigation rules and menu data.
//   - List of available operator users.
//
// Note: The variables dbClientConfig and dbServerConfig (used for database configuration)
// are assumed to be defined elsewhere in the application.
func AssignOperatorsPOST_GET(c *gin.Context) {
	// Retrieve the current session from the Gin context.
	session := sessions.Default(c)
	// Extract the session ID stored under the key "userData". It is expected to be a string.
	var sessionID string = session.Get("userData").(string)
	// Retrieve common session data using the session ID.
	s, err := utils.GetCommonSession(sessionID)

	// Check if the current user has permission to assign operators.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "assign_operators", s.CurrentRole, c) {
		// If permission is not granted, exit the handler.
		return
	}

	// Retrieve a list of general users with the "op" (operator) role from the security controller.
	// The returned data is used to populate the list of available operators.
	_, res := security_ctrl.GetGeneralUsersByRole("op", &db.ConnData{}, dbClientConfig, dbServerConfig)

	// Set HTTP headers to disable caching of the response.
	common_facades.SetHeaderNoCache(c)

	// Initialize a variable to hold menu data.
	var menu map[string][]map[string]string

	// If no error occurred while retrieving session data, use the current menu from the session.
	if err == nil {
		menu = s.CurrentMenu
	}

	// Render the HTML template for the assign operators view.
	// The template is provided with various parameters required for dynamic content generation.
	common_facades.RenderTemplate(c,
		salvia_daos.VictimCaseEntityName, // Entity name for victim case.
		"salvia",                         // Module or section name.
		"victim_case/",                   // Template directory path.
		salvia_config.HTML_Templates,     // Map of HTML templates.
		"assign_operators",               // Specific template name.
		utils.GetFullHtmlTemplates(),     // Complete set of HTML templates.
		utils.DEFAULT_VIEW,               // Default view identifier.
		utils.DEFAULT_PANIC_TEMPLATE,     // Template used in case of panic.
		map[string]interface{}{ // Data passed to the template.
			"windowTitle":    salvia_config.Locale["sp"]["assign_operators_window_title"],
			"currentUser":    s.Names + " " + s.LastNames,
			"salviaFormPath": salvia_config.FormPaths[s.Lang]["AssignOperatorsPOST"],
			"locale":         salvia_config.Locale,
			"lang":           s.Lang,
			"menu":           menu,
			"users":          res,
			"nav_rules":      salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["assign_operators"]),
		},
		utils.GetFullHtmlFuncMap(), // Function map for template processing.
	)
}
