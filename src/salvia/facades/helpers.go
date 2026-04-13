// Package salvia_facades — helpers.go
// Funciones compartidas por todas las fachadas del módulo salvia.
package salvia_facades

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionContext agrupa sesión y menú para no repetir las mismas líneas en cada fachada.
type SessionContext struct {
	Session utils.CommonSession
	Menu    map[string][]map[string]string
}

// CheckAndGetSession extrae la sesión y verifica el permiso.
// Retorna nil si la sesión no existe o el rol no tiene el permiso — ya escribió la respuesta HTTP.
func CheckAndGetSession(c *gin.Context, permission string) *SessionContext {
	session := sessions.Default(c)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return nil
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, permission, s.CurrentRole, c) {
		return nil
	}
	return &SessionContext{Session: *s, Menu: s.CurrentMenu}
}

// BaseTemplateVars construye el mapa base de variables que toda pantalla necesita.
func BaseTemplateVars(sc *SessionContext, lang string) map[string]interface{} {
	return map[string]interface{}{
		"currentUser": sc.Session.Names + " " + sc.Session.LastNames,
		"menu":        sc.Menu,
		"lang":        lang,
		"locale":      salvia_config.Locale,
	}
}
