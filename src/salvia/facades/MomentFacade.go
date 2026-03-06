// Package salvia_facades contiene las funciones de fachada que permiten la interacción
// entre la capa de presentación y los controladores, gestionando las operaciones relacionadas
// con los "momentos" en la aplicación.
package salvia_facades

import (
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	"bytes"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MomentPUT maneja las solicitudes HTTP PUT para actualizar un "momento" específico.
// Realiza las siguientes operaciones:
// 1. Recupera la sesión del usuario y extrae el identificador de sesión.
// 2. Lee el cuerpo de la solicitud para obtener los datos a actualizar.
// 3. Obtiene los parámetros de la URL necesarios para identificar el "momento" (id, momentCode y entityBranchIcode).
// 4. Verifica que los parámetros no estén vacíos.
// 5. Comprueba que el usuario tenga permisos para realizar la actualización.
// 6. Invoca al controlador correspondiente para efectuar la actualización.
// 7. Envía la respuesta HTTP con el código de estado y el contenido resultante.
//
// Parámetros:
//   - c: Contexto de Gin que encapsula la solicitud HTTP, la respuesta y otros metadatos.
func MomentPUT(c *gin.Context) {
	// Obtiene la sesión actual a partir del contexto.
	session := sessions.Default(c)
	// Extrae el identificador de sesión almacenado en "userData" y lo convierte a string.
	var sessionID string = session.Get("userData").(string)
	// Recupera la sesión común asociada al usuario usando el identificador.
	s, _ := utils.GetCommonSession(sessionID)

	// Inicializa la respuesta con un mensaje vacío y el código HTTP por defecto (400 - Bad Request).
	var res string
	var code int = 400

	// Crea un buffer para leer el cuerpo de la solicitud HTTP.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	// Extrae los parámetros de la URL que identifican el "momento".
	id := c.Param("id")
	mCode := c.Param("momentCode")
	ebIcode := c.Param("entityBranchIcode")

	// Verifica que todos los parámetros requeridos estén presentes.
	if id != "" && mCode != "" && ebIcode != "" {
		// Verifica que el usuario tenga permisos para actualizar el "momento".
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_moment", s.CurrentRole, c) {
			// Si no tiene permisos, se finaliza la ejecución sin enviar respuesta adicional.
			return
		}
		// Llama al controlador para actualizar el "momento" pasando los datos leídos del cuerpo,
		// la sesión del usuario y los parámetros identificadores.
		// Los parámetros "security", "&db.ConnData{}", "dbClientConfig" y "dbServerConfig" se utilizan
		// para configurar la actualización y la conexión a la base de datos.
		code, res = salvia_ctrl.UpdateMomentByVictimCaseICodeAndMomentCodeAndEntityBranch(
			buf.String(), *s, id, mCode, ebIcode, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Envía la respuesta HTTP al cliente utilizando el código, la longitud de la respuesta y
	// el contenido en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
