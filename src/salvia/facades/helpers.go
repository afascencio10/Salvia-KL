// Package salvia_facades ÔÇö helpers.go
// Funciones compartidas por todas las fachadas del m├│dulo salvia.
package salvia_facades

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionContext agrupa sesi├│n y men├║ para no repetir las mismas l├¡neas en cada fachada.
type SessionContext struct {
	Session utils.CommonSession
	Menu    map[string][]map[string]string
}

// CheckAndGetSession extrae la sesi├│n y verifica el permiso.
// Retorna nil si la sesi├│n no existe o el rol no tiene el permiso ÔÇö ya escribi├│ la respuesta HTTP.
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

// BackFallbackPath devuelve el destino del botón "Volver" cuando no hay historial
// utilizable (pestaña nueva, URL pegada directo, recarga sobre el detalle).
//
// Debe ser el listado desde el que cada rol llega naturalmente al detalle:
//   - El respaldo por defecto es /salvia/lista-casos ("Lista Casos"), no
//     /salvia/casos, que es la pantalla de reportes — llegar ahí es justo el
//     síntoma reportado ("atrás lleva al reporte en lugar de a la lista").
//     Además /salvia/casos exige get_victim_cases, permiso más restringido que
//     el get_list_cases del listado.
//   - Los roles que no trabajan sobre casos vuelven a su propio listado; para
//     ellos /salvia/casos (Reportes) responde 401. ps y ts sí tienen
//     get_list_cases, así que el fallback por defecto aplica también a ellos.
func BackFallbackPath(role string) string {
	switch role {
	case "an":
		return "/salvia/notificaciones"
	case "en":
		return "/salvia/mis-barreras"
	case "ps", "ts":
		return "/salvia/mis-remisiones-psicosocial"
	case "mc":
		return "/salvia/casos-mecanismo"
	}
	return "/salvia/lista-casos"
}

// BaseTemplateVars construye el mapa base de variables que toda pantalla necesita.
func BaseTemplateVars(sc *SessionContext, lang string) map[string]interface{} {
	return map[string]interface{}{
		"currentUser": sc.Session.Names + " " + sc.Session.LastNames,
		"currentRole": sc.Session.CurrentRole,
		"menu":        sc.Menu,
		"lang":        lang,
		"locale":      salvia_config.Locale,
	}
}
