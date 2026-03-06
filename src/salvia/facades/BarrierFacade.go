// Package salvia_facades contiene las funciones que actúan como fachada para la gestión de logs de casos
// en el sistema Salvia.
package salvia_facades

import (
	// Paquetes internos del proyecto

	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"

	// Paquetes estándar

	"net/http"
	"strings"

	// Paquetes de terceros
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// BarrierGET maneja las solicitudes GET para obtener barreras
// La función realiza lo siguiente:
//   - Obtiene y valida la sesión del usuario.
//   - Verifica que el usuario tenga permisos para establecer un log de caso.
//   - Lee el cuerpo de la solicitud.
//   - Llama al controlador para obtener las barreras
//   - Envía la respuesta al cliente en formato JSON.
func BarrierGET(c *gin.Context) {
	// Se obtiene la sesión actual del contexto
	session := sessions.Default(c)
	// Se extrae el identificador de sesión del usuario
	var sessionID string = session.Get("userData").(string)
	// Se recupera la información común de la sesión usando el ID
	s, _ := utils.GetCommonSession(sessionID)
	// Se inicializan las variables para el código HTTP y la respuesta
	var code int = http.StatusBadRequest
	var res string

	// Se verifica si el usuario tiene permiso para establecer un log de caso.
	// La verificación se basa en la configuración de permisos por rol.
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_barrier", s.CurrentRole, c) {
		return
	}

	// Se obtiene el parámetro 'id' de la URL
	icode := c.Param("id")
	by := c.Param("by")

	if icode != "" {
		if by == "" {
			code, res = salvia_ctrl.GetBarrierByICode(icode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		} else if salvia_config.TranslateLocale(by, s.Lang) == "barrier_by_sector" {
			code, res = salvia_ctrl.GetBarriersBySectorBarrier(icode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}
	}

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
