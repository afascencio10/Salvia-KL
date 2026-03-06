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
	"bytes"
	"net/http"
	"strings"

	// Paquetes de terceros
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// FollowUpEntryPOST maneja las solicitudes POST para crear  un seguimiento de caso.
// La función realiza lo siguiente:
//   - Obtiene y valida la sesión del usuario.
//   - Verifica que el usuario tenga permisos para establecer un log de caso.
//   - Lee el cuerpo de la solicitud.
//   - Llama al controlador para procesar y almacenar el log.
//   - Envía la respuesta al cliente en formato JSON.
func FollowUpEntryActingPOST(c *gin.Context) {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_follow_up_entry_acting", s.CurrentRole, c) {
		return
	}

	barrierICode := c.Param("b")
	entryICode := c.Param("e")
	if barrierICode != "" && entryICode != "" {
		// Se crea un buffer para leer el cuerpo de la solicitud
		buf := new(bytes.Buffer)
		// Se lee todo el contenido del cuerpo de la solicitud
		buf.ReadFrom(c.Request.Body)
		// Se llama al controlador para establecer el log de caso.
		// Se pasan: el contenido del cuerpo, el id del caso, el origen ("salvia"),
		// la sesión del usuario y las configuraciones de base de datos (dbClientConfig, dbServerConfig).
		code, res = salvia_ctrl.SetFollowUpEntryActing(buf.String(), barrierICode, entryICode, &db.ConnData{}, *s, dbClientConfig, dbServerConfig)
	}

	// Se envía la respuesta al cliente en formato JSON utilizando el código y la respuesta obtenida.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
