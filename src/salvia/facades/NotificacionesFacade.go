package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// notifModalTemplates lista los 6 partials de modales de gestión de oficios.
// Cada archivo contiene un bloque {{ define "notifications/modal_xxx.html" }}
// que es referenciado desde notificaciones.html con {{ template "..." . }}.
var notifModalTemplates = []string{
	"frontend/html/salvia/notifications/modal_corregir.html",
	"frontend/html/salvia/notifications/modal_proyectar.html",
	"frontend/html/salvia/notifications/modal_revisar.html",
	"frontend/html/salvia/notifications/modal_radicar.html",
	"frontend/html/salvia/notifications/modal_registrar_respuesta.html",
	"frontend/html/salvia/notifications/modal_aprobar.html",
}

// NotificacionesGET renderiza la pantalla de Notificaciones Salvia.
// Ruta: GET /salvia/notificaciones
// Roles con acceso funcional: op (Agente Seguimiento), ro (Revisor Operativo), an (Agente de Notificaciones)
func NotificacionesGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionIDVal := session.Get("userData")
	if sessionIDVal == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}

	sessionID := sessionIDVal.(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}

	// Validar que el usuario tenga permiso para acceder a notificaciones.
	// Roles con acceso funcional en frontend: op, ro y an.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_notificaciones", s.CurrentRole, c) {
		return
	}

	extraTemplates := append(utils.GetFullHtmlTemplates(), notifModalTemplates...)

	common_facades.RenderTemplate(
		c,
		"notificaciones",
		"salvia",
		"notifications/",
		salvia_config.HTML_Templates,
		"notificaciones",
		extraTemplates,
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Notificaciones Salvia",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"currentUserId": s.UserICode,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
