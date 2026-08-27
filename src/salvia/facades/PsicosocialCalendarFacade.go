// Package salvia_facades — PsicosocialCalendarFacade: renderiza pantalla del calendario.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PsicosocialCalendarGET renderiza la pantalla del calendario psicosocial.
func PsicosocialCalendarGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, _ := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	fmt.Printf("[PsicosocialCalendarGET] userICode=%s currentRole=%q\n", s.UserICode, s.CurrentRole)
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_psychosocial_calendar", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	var menu map[string][]map[string]string
	if s != nil {
		menu = s.CurrentMenu
	}

	common_facades.RenderTemplate(c, "PsicosocialCalendar", "salvia", "psicosocial_calendar/", salvia_config.HTML_Templates, "psicosocial_calendar", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Calendario",
			"currentUser": s.Names + " " + s.LastNames,
			"currentUserICode": s.UserICode,
			"currentRole": s.CurrentRole,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"menu":        menu,
		}, utils.GetFullHtmlFuncMap())
}
