// Package salvia_facades — PsicosocialSesionFacade: pantalla de registro de sesión psicosocial.
package salvia_facades

import (
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RegistrarSesionPsicosocialGET renderiza la pantalla para registrar una sesión psicosocial.
//
// Ruta: GET /salvia/psicosocial/registrar/:id
//
// El parámetro :id es el UUID de salvia.psychosocial_support.
// Cuando :id = "mock", el frontend detecta el query param ?mock=casoN y usa datos simulados
// en lugar de llamar al backend real (útil para validar el flujo sin ejecutar el seed).
func RegistrarSesionPsicosocialGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_registrar_sesion_psicosocial", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	psicosocialID := c.Param("id")

	common_facades.RenderTemplate(
		c,
		"registrar_sesion_psicosocial",
		"salvia",
		"psicosocial/",
		salvia_config.HTML_Templates,
		"registrar_sesion_psicosocial",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   "Registrar Sesión Psicosocial",
			"currentUser":   s.Names + " " + s.LastNames,
			"currentRole":   s.CurrentRole,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
			"psicosocialId": psicosocialID,
			"userICode":     s.UserICode,
		},
		utils.GetFullHtmlFuncMap(),
	)
}

// TestCasosPsicosocialGET renderiza la pantalla de selección de casos de prueba.
// Permite navegar a cada escenario del flujo psicosocial con datos simulados.
//
// Ruta: GET /salvia/psicosocial/test
func TestCasosPsicosocialGET(c *gin.Context) {
	session := sessions.Default(c)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_test_psicosocial", s.CurrentRole, c) {
		return
	}

	common_facades.SetHeaderNoCache(c)

	common_facades.RenderTemplate(
		c,
		"test_casos_psicosocial",
		"salvia",
		"psicosocial/",
		salvia_config.HTML_Templates,
		"test_casos_psicosocial",
		utils.GetFullHtmlTemplates(),
		utils.DEFAULT_VIEW,
		utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": "Test — Casos Psicosocial",
			"currentUser": s.Names + " " + s.LastNames,
			"currentRole": s.CurrentRole,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
		},
		utils.GetFullHtmlFuncMap(),
	)
}
