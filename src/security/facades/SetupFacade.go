// Package security_routers — handlers HTTP para los endpoints de setup/administración.
// Estos endpoints están protegidos por API Key (header X-Setup-Key) y no requieren sesión.
package security_routers

import (
	"bytes"
	"net/http"
	"os"
	"strings"

	"bitsflow/common/db"
	security_ctrl "bitsflow/security/controllers"

	"github.com/gin-gonic/gin"
)

// securityKeyDefault es la clave por defecto para los endpoints de administración.
// Puede sobreescribirse con la variable de entorno SETUP_API_KEY.
const securityKeyDefault = "s4lv1a_kr31vo"

// adminAPIKeyMiddleware valida el header X-Security-Key.
// Prioridad: variable de entorno SETUP_API_KEY → clave por defecto.
func adminAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := os.Getenv("SETUP_API_KEY")
		if apiKey == "" {
			apiKey = securityKeyDefault
		}

		clientKey := c.GetHeader("X-Security-Key")
		if clientKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Header X-Security-Key requerido",
			})
			return
		}

		if clientKey != apiKey {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "X-Security-Key inválida",
			})
			return
		}

		c.Next()
	}
}

// setupAPIKeyMiddleware es un alias de adminAPIKeyMiddleware para compatibilidad interna.
var setupAPIKeyMiddleware = adminAPIKeyMiddleware

// SetupCreateUserPOST maneja POST /api/v1/setup/usuario
// Crea un nuevo usuario en el sistema con los datos recibidos en el body JSON.
//
// Body esperado:
//
//	{
//	  "login":             "john.doe",          // obligatorio
//	  "pass":              "Password123!",       // obligatorio
//	  "team":              "Equipo A",           // opcional
//	  "roleCodes":         ["op"],               // obligatorio — array de códigos de rol
//	  "names":             "John",               // obligatorio
//	  "lastNames":         "Doe",                // obligatorio
//	  "email":             "john@example.com",   // obligatorio
//	  "phone":             "3001234567",         // opcional (default: 0000000000)
//	  "docType":           "CC",                 // obligatorio
//	  "docNumber":         "12345678",           // obligatorio
//	  "gender":            "M",                  // obligatorio
//	  "lang":              "sp",                 // opcional (default: sp)
//	  "entityBranchICode": "abc-123",            // opcional — iCode de la sede
//	  "townCode":          "11001000"            // opcional (default: 11001000 = Bogotá)
//	}
//
// Respuesta exitosa (201):  {"success":true,"icode":"...","login":"..."}
// Códigos de rol disponibles: ad, do, et, op, ro, no, sv, fo, an, en
func SetupCreateUserPOST(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	code, res := security_ctrl.SetupCreateUser(buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// MigrateExcelPOST maneja POST /api/v1/admin/migrate/excel
// Recibe un archivo .xlsx como multipart/form-data (campo "file"),
// parsea la hoja "Usuarios en BD" y ejecuta la migración masiva en el servidor.
//
// Uso desde Postman:
//   Body → form-data → Key: "file", Type: File, Value: <archivo .xlsx>
//
// Respuesta: { sheet, total_rows, valid_rows, skipped_rows,
//              created_count, updated_count, failed_count,
//              created[], updated[], failed[], parsed_but_skipped[] }
func MigrateExcelPOST(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'file' requerido (multipart/form-data)"})
		return
	}

	// Validar extensión
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".xlsx") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo debe ser .xlsx"})
		return
	}

	// Leer bytes del archivo
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo abrir el archivo: " + err.Error()})
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	code, res := security_ctrl.MigrateFromExcel(buf.Bytes(), dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// MigrateUsersPOST maneja POST /api/v1/admin/migrate/users
// Crea o actualiza masivamente usuarios a partir de un array JSON.
//
// Body esperado:
//
//	{
//	  "users": [
//	    {
//	      "login":     "ana.garcia",         // obligatorio
//	      "names":     "Ana",                // obligatorio al crear; ignorado al actualizar
//	      "lastNames": "García López",       // obligatorio al crear; ignorado al actualizar
//	      "email":     "ana@correo.com",     // obligatorio al crear; ignorado al actualizar
//	      "phone":     "3001234567",         // opcional
//	      "docType":   "CC",                 // obligatorio al crear
//	      "docNumber": "1234567890",         // obligatorio al crear
//	      "gender":    "ma",                 // opcional (default: ma)
//	      "team":      "Riesgo alto",        // Riesgo alto | Riesgo bajo | SIN_EQUIPO | Hombres
//	      "roleCodes": ["ro"],               // sv → supervisor; cualquier otro → ro (operador)
//	      "pass":      "Password123!",       // obligatorio al crear; ignorado al actualizar
//	      "lang":      "sp",                 // opcional
//	      "townCode":  "11001000"            // opcional
//	    }
//	  ]
//	}
//
// Respuesta: { total_processed, created_count, updated_count, skipped_count, created[], updated[], skipped[] }
func MigrateUsersPOST(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	code, res := security_ctrl.MigrateBatchUsers(buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// MigrateTestUsersPOST maneja POST /api/v1/admin/migrate/test-users
// Crea credenciales temporales de prueba con login "test.{login}" para verificación en producción.
//
// Body esperado:
//
//	{
//	  "testPassword": "Salvia@Test2026!",    // opcional — contraseña para todos los test users
//	  "users": [
//	    {
//	      "login":     "ana.garcia",         // obligatorio — se crea como "test.ana.garcia"
//	      "names":     "Ana",
//	      "lastNames": "García López",
//	      "email":     "ana@correo.com",     // opcional — se usa como "test.ana@correo.com"
//	      "docNumber": "1234567890",         // opcional — se guarda como "TEST_1234567890"
//	      "team":      "Riesgo alto",
//	      "roleCodes": ["ro"]
//	    }
//	  ]
//	}
func MigrateTestUsersPOST(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	code, res := security_ctrl.MigrateTestUsers(buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// SetupUpdateUserRoleTeamPATCH maneja PATCH /api/v1/setup/usuario/:login/rol-equipo
// Actualiza el rol y/o el equipo de un usuario existente, identificado por su login.
//
// Parámetro de ruta: :login — login del usuario a modificar
//
// Body esperado:
//
//	{
//	  "roleCodes": ["sv"],      // opcional — reemplaza TODOS los roles actuales
//	  "team":      "Equipo B"   // opcional — nuevo nombre de equipo
//	}
//
// Nota: si se envían roleCodes, se eliminan todos los roles actuales y se asignan los nuevos.
// Debe enviarse al menos uno de los dos campos.
//
// Respuesta exitosa (200): {"success":true,"login":"...","icode":"..."}
func SetupUpdateUserRoleTeamPATCH(c *gin.Context) {
	login := c.Param("login")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	code, res := security_ctrl.SetupUpdateUserRoleTeam(login, buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
