package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"

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
// NOTA: Sin validación de sesión — acceso público para desarrollo/prototipo.
func NotificacionesGET(c *gin.Context) {
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
			"windowTitle": "Notificaciones Salvia",
			"currentUser": "",
			"locale":      salvia_config.Locale,
			"lang":        "sp",
			"currentRole": "",
		},
		utils.GetFullHtmlFuncMap(),
	)
}
