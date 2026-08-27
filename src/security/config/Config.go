package security_config

import (
	"net"
	"os"
	"strings"
)

// APP_DOMAIN es la URL base por defecto de la aplicación, usada para construir
// los enlaces que se envían por correo cuando no es posible determinar el
// dominio desde la petición HTTP. Puede sobreescribirse con la variable de
// entorno APP_DOMAIN; si no está definida se usa el dominio de producción.
var APP_DOMAIN string = getAppDomain()

// allowedAppHosts es la lista de hosts desde los que la aplicación puede ser
// servida. Solo estos hosts se aceptan al construir enlaces dinámicos a partir
// de la cabecera Host de la petición; cualquier otro valor se descarta para
// prevenir ataques de "password reset poisoning" (inyección de un dominio
// atacante en el enlace del correo).
var allowedAppHosts = map[string]bool{
	"salvia-colombia.co":         true, // producción
	"pruebas.salvia-colombia.co": true, // desarrollo
	"localhost":                  true, // entorno local
	"127.0.0.1":                  true, // entorno local
}

func getAppDomain() string {
	if v := os.Getenv("APP_DOMAIN"); v != "" {
		return v
	}
	return "https://salvia-colombia.co"
}

// ResolveAppDomain construye la URL base de la aplicación a partir del host de
// la petición HTTP (cabecera Host). Si el host está en la lista de permitidos,
// el enlace apunta al mismo ambiente desde el que el usuario hizo la solicitud
// (producción, desarrollo o local). Si no, se usa APP_DOMAIN como respaldo.
func ResolveAppDomain(requestHost string) string {
	host := requestHost
	// El Host puede incluir puerto (p.ej. "localhost:8080"); se separa para
	// validar solo el nombre, pero el enlace conserva el puerto original.
	if h, _, err := net.SplitHostPort(requestHost); err == nil {
		host = h
	}
	if allowedAppHosts[strings.ToLower(host)] {
		return "https://" + requestHost
	}
	return APP_DOMAIN
}

const EMAIL_SERVER_HOST_PATH string = "mail.soggroup.com"
const EMAIL_SERVER_HOST_PORT string = "465"
const EMAIL_SERVER_HOST_USERNAME string = "salvia_bot@soggroup.com"
const EMAIL_SERVER_HOST_PASS string = "^voo,*hPgLw8"

const EMAIL_SERVER_FROM_USER string = "salvia_bot@soggroup.com"
const EMAIL_SERVER_FROM_NAME string = "Salvia Bot"
