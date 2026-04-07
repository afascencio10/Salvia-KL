# Análisis del Código - Sistema SALVIA
## Ministerio de Igualdad - Sistema de Información

---

## 1. ESTRUCTURA DEL CÓDIGO

### 1.1 Arquitectura General
El proyecto sigue una **arquitectura modular en capas** basada en Go (Golang) con el framework **Gin** para el servidor web.

```
Fuentes/
├── main.go                    # Punto de entrada de la aplicación
├── common/                    # Módulo común compartido
│   ├── config/               # Configuraciones globales
│   ├── controllers/          # Controladores comunes
│   ├── dao/                  # Data Access Objects comunes
│   ├── db/                   # Gestión de conexiones a base de datos
│   ├── facades/              # Capa de presentación/routing común
│   ├── routers/              # Enrutadores HTTP
│   ├── scripts/              # Scripts auxiliares
│   └── utils/                # Utilidades compartidas
├── salvia/                    # Módulo de negocio principal (casos de víctimas)
│   ├── config/               # Configuración específica de Salvia
│   ├── controllers/          # Lógica de negocio
│   ├── dao/                  # Modelos y acceso a datos
│   ├── facades/              # Handlers HTTP
│   └── routers/              # Rutas específicas
├── security/                  # Módulo de seguridad y autenticación
│   ├── config/
│   ├── controllers/
│   ├── dao/
│   ├── facades/
│   └── routers/
├── frontend/                  # Recursos del cliente
│   ├── html/                 # Páginas HTML
│   ├── js/                   # JavaScript (Vue.js)
│   ├── css/                  # Estilos
│   ├── templates/            # Plantillas Go HTML
│   └── plugins/              # Librerías de terceros
└── config/                    # Configuración de base de datos
    └── db_config.json
```

### 1.2 Patrón de Capas

El sistema implementa una **arquitectura de 3 capas**:

1. **Capa de Presentación (Facades)**
   - Maneja las peticiones HTTP
   - Valida sesiones y permisos
   - Renderiza templates HTML o devuelve JSON
   - Ejemplo: `VictimContactFacade.go`

2. **Capa de Lógica de Negocio (Controllers)**
   - Procesa la lógica de negocio
   - Valida datos de entrada
   - Coordina operaciones entre DAOs
   - Ejemplo: `VictimContactController.go`

3. **Capa de Acceso a Datos (DAO)**
   - Interactúa directamente con PostgreSQL
   - Define estructuras de datos (DTOs)
   - Ejecuta queries SQL
   - Ejemplo: `VictimContactDAO.go`

---

## 2. DISEÑO Y PATRONES

### 2.1 Patrones de Diseño Identificados

#### **Facade Pattern**
Los facades actúan como punto de entrada unificado para cada módulo:
```go
// MainRouter.go inicializa todos los facades
func main() {
    var router *gin.Engine = common_routers.InitRouter()
    salvia_facades.StartRouter(router)
    security_routers.StartRouter(router)
    common_routers.StartRouter()
}
```

#### **DTO (Data Transfer Object)**
Cada entidad tiene su DTO con validaciones:
```go
type VictimContactDTO struct {
    VictimContactId           uint      `json:"victimContactId"`
    VictimContactICode        string    `json:"victimContactICode"`
    VictimContactCreationDate time.Time `json:"victimContactCreationDate"`
    // ... más campos
}
```

#### **Repository Pattern (DAO)**
Separación clara entre lógica de negocio y acceso a datos:
- `SetVictimContact()` - Crear/actualizar
- `GetVictimContact()` - Obtener por criterio
- `GetVictimContacts()` - Listar
- `InvalidateVictimContactByICode()` - Soft delete

#### **Middleware Pattern**
Uso extensivo de middlewares en Gin:
- `AuthMiddleware()` - Autenticación
- `limitMiddleware()` - Control de concurrencia (9999 peticiones)
- `ForceHTTPS()` - Redirección HTTP → HTTPS

#### **Session Management**
Sistema de sesiones personalizado con almacenamiento en memoria:
```go
store := cookie.NewStore([]byte("secret"))
store.Options(sessions.Options{
    MaxAge:   14400,  // 4 horas
    HttpOnly: false,
    Secure:   true,
})
```

### 2.2 Tecnologías y Frameworks

**Backend:**
- **Go 1.19**
- **Gin** - Framework web
- **pgx/v4** - Driver PostgreSQL
- **gin-contrib/sessions** - Gestión de sesiones
- **dchest/captcha** - Validación CAPTCHA

**Frontend:**
- **Vue.js 3** - Framework JavaScript
- **Bootstrap** - UI Framework
- **Leaflet** - Mapas interactivos
- **FontAwesome** - Iconos

**Base de Datos:**
- **PostgreSQL** con pool de conexiones (80 conexiones)

**Seguridad:**
- **TLS/HTTPS** obligatorio
- Certificados en `certs/salviaTest.crt`

---

## 3. DEUDA TÉCNICA

### 3.1 Crítica (Alta Prioridad)

#### **1. Credenciales Hardcodeadas**
```json
// config/db_config.json
{
    "password":"asd876.!@asdSDS5a36Z"  // ❌ Contraseña en texto plano
}
```
**Riesgo:** Exposición de credenciales en repositorio  
**Solución:** Usar variables de entorno o gestores de secretos (HashiCorp Vault, AWS Secrets Manager)

#### **2. Secret Key Débil**
```go
store := cookie.NewStore([]byte("secret"))  // ❌ Clave trivial
```
**Riesgo:** Sesiones vulnerables a ataques  
**Solución:** Generar clave aleatoria de 32+ bytes y almacenarla de forma segura

#### **3. Gestión de Sesiones en Memoria**
```go
// CommonSession.go - Sesiones almacenadas en memoria
var sessions map[string]CommonSession = map[string]CommonSession{}
```
**Riesgo:** Pérdida de sesiones al reiniciar, no escalable  
**Solución:** Usar Redis o base de datos para persistencia

#### **4. Sin Manejo de Errores Robusto**
```go
err := router.RunTLS(":443", "certs/salviaTest.crt", "certs/salviaTest.key")
println(err)  // ❌ Solo imprime el error
```
**Riesgo:** Fallos silenciosos en producción  
**Solución:** Logging estructurado (logrus, zap) y monitoreo

#### **5. Configuración HTTPS Insegura**
```go
store.Options(sessions.Options{
    HttpOnly: false,  // ❌ Vulnerable a XSS
    Secure:   true,
    SameSite: http.SameSiteDefaultMode,  // ❌ Debería ser Strict/Lax
})
```
**Riesgo:** Cookies accesibles desde JavaScript  
**Solución:** `HttpOnly: true`, `SameSite: http.SameSiteStrictMode`

### 3.2 Alta (Media Prioridad)

#### **6. Código Comentado**
```go
//go initWebSocket()
//prepareSound()
//gin.SetMode(gin.ReleaseMode)
```
**Problema:** Código muerto que genera confusión  
**Solución:** Eliminar o documentar por qué está comentado

#### **7. Falta de Tests**
No se encontraron archivos de pruebas (`*_test.go`)  
**Riesgo:** Regresiones no detectadas  
**Solución:** Implementar tests unitarios e integración

#### **8. Validaciones Duplicadas**
Validaciones en múltiples capas sin centralización:
- Frontend (JavaScript)
- Facades (Go)
- Controllers (Go)
- DAO (Go)

**Problema:** Mantenimiento complejo  
**Solución:** Centralizar validaciones en una capa

#### **9. Dependencias Desactualizadas**
```go
go 1.19  // Versión antigua (actual: 1.23)
```
**Riesgo:** Vulnerabilidades de seguridad  
**Solución:** Actualizar a Go 1.21+ y dependencias

#### **10. Pool de Conexiones Fijo**
```go
var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: 80}
```
**Problema:** No configurable dinámicamente  
**Solución:** Hacer configurable por entorno

### 3.3 Media (Baja Prioridad)

#### **11. Nombres de Variables Inconsistentes**
```go
var translatedModule string = salvia_config.Locale["sp"][module]
var translatedEntity string
var translatedNew string = salvia_config.Locale["sp"]["new"]
```
**Problema:** Mezcla de español e inglés  
**Solución:** Estandarizar nomenclatura

#### **12. Funciones Muy Largas**
Facades con funciones de 200+ líneas  
**Problema:** Difícil mantenimiento  
**Solución:** Refactorizar en funciones más pequeñas

#### **13. Sin Documentación de API**
No hay especificación OpenAPI/Swagger  
**Problema:** Dificulta integración  
**Solución:** Generar documentación automática

#### **14. Logging Insuficiente**
```go
println(err)  // ❌ No estructurado
```
**Problema:** Dificulta debugging en producción  
**Solución:** Implementar logging estructurado

#### **15. Sin Rate Limiting por Usuario**
```go
sem = make(chan struct{}, 9999)  // Límite global
```
**Problema:** Un usuario puede consumir todos los recursos  
**Solución:** Rate limiting por IP/usuario

---

## 4. FUNCIONAMIENTO

### 4.1 Flujo de Inicio de Aplicación

```
1. main.go
   ↓
2. common_routers.InitRouter()
   - Configura Gin
   - Establece middlewares (auth, HTTPS, límite concurrencia)
   - Carga sesiones
   - Configura archivos estáticos
   ↓
3. salvia_facades.StartRouter(router)
   - Registra rutas de módulo Salvia
   ↓
4. security_routers.StartRouter(router)
   - Registra rutas de seguridad
   ↓
5. common_routers.StartRouter()
   - Inicia servidor HTTPS en puerto 443
   - Redirige HTTP (80) → HTTPS (443)
```

### 4.2 Flujo de Petición HTTP

**Ejemplo: Crear Contacto de Víctima**

```
1. Cliente → POST /salvia/contacto-victima
   ↓
2. Middleware: limitMiddleware()
   - Verifica límite de concurrencia
   ↓
3. Middleware: ForceHTTPS()
   - Asegura conexión HTTPS
   ↓
4. Middleware: AuthMiddleware()
   - Valida sesión activa
   - Redirige a login si no autenticado
   ↓
5. Facade: VictimContactPOST()
   - Obtiene sesión del usuario
   - Verifica permiso "set_victim_contact"
   - Lee body de la petición
   ↓
6. Controller: SetVictimContact()
   - Valida JSON de entrada
   - Valida campos según FieldDefinitions
   - Maneja errores de validación
   ↓
7. DAO: SetVictimContact()
   - Obtiene conexión del pool
   - Ejecuta INSERT/UPDATE en PostgreSQL
   - Libera conexión
   ↓
8. Respuesta JSON al cliente
   - Éxito: {"success": true, "data": {...}}
   - Error: {"success": false, "errors": {...}}
```

### 4.3 Gestión de Base de Datos

**Pool de Conexiones:**
```go
// PostgresConnection.go
func GetConnection(connData *ConnData, clientConfig *DBClientConfig, 
                   serverConfig *DBServerConfig) (*ConnData, error) {
    // Obtiene conexión del pool o crea nueva
    // Pool size: 80 conexiones
}

func ReleaseConnection(connData *ConnData) error {
    // Devuelve conexión al pool
}
```

**Transacciones:**
```go
// Parámetro inTransaction en DAOs
connData, err := dao.SetVictimContact(
    victimContact, 
    true,  // inTransaction = true
    module, 
    connData, 
    clientConfig, 
    serverConfig
)
```

### 4.4 Sistema de Permisos

**Basado en Roles:**
```go
// Roles identificados:
// - "sv" (Supervisor)
// - "op" (Operador)
// - "ad" (Admin - inferido)

// Verificación de permisos
if !utils.CheckPermission(
    salvia_config.PermissionsByRole, 
    "set_victim_contact", 
    s.CurrentRole, 
    c
) {
    return  // 403 Forbidden
}
```

### 4.5 Renderizado de Templates

**Sistema Híbrido:**
- **Backend:** Go templates (`html/template`)
- **Frontend:** Vue.js 3 para interactividad

```go
common_routers.RenderTemplate(
    c, 
    entityName, 
    module, 
    htmlFolder, 
    configTemplate, 
    templateName, 
    extraTemplates, 
    viewTemplate, 
    errorTemplate,
    templateFields,  // Datos para el template
    funcMap          // Funciones personalizadas
)
```

### 4.6 Módulos Principales

#### **Módulo SALVIA (Casos de Víctimas)**
Entidades principales:
- `VictimContact` - Contactos iniciales
- `VictimCase` - Casos completos de víctimas
- `Alert` - Alertas del sistema
- `CaseLog` - Registro de actividades
- `Moment` - Estados del caso
- `EntityBranch` - Entidades de atención

#### **Módulo SECURITY (Autenticación)**
Entidades principales:
- `GeneralUser` - Usuarios del sistema
- `City` - Ciudades
- `Town` - Municipios
- `Role` - Roles y permisos

---

## 5. RECOMENDACIONES

### 5.1 Seguridad (Urgente)
1. ✅ Migrar credenciales a variables de entorno
2. ✅ Implementar rotación de secretos
3. ✅ Configurar `HttpOnly: true` en cookies
4. ✅ Actualizar a Go 1.21+
5. ✅ Implementar rate limiting por usuario
6. ✅ Añadir CORS configurado correctamente
7. ✅ Implementar CSP (Content Security Policy)

### 5.2 Arquitectura (Corto Plazo)
1. ✅ Migrar sesiones a Redis
2. ✅ Implementar logging estructurado (zap/logrus)
3. ✅ Añadir health checks (`/health`, `/ready`)
4. ✅ Implementar graceful shutdown
5. ✅ Dockerizar la aplicación

### 5.3 Calidad de Código (Medio Plazo)
1. ✅ Implementar tests unitarios (>70% cobertura)
2. ✅ Añadir tests de integración
3. ✅ Configurar CI/CD (GitHub Actions/GitLab CI)
4. ✅ Implementar linting (golangci-lint)
5. ✅ Generar documentación OpenAPI

### 5.4 Escalabilidad (Largo Plazo)
1. ✅ Separar frontend en SPA independiente
2. ✅ Implementar caché (Redis)
3. ✅ Considerar microservicios para módulos
4. ✅ Implementar message queue (RabbitMQ/Kafka)
5. ✅ Añadir observabilidad (Prometheus/Grafana)

---

## 6. CONCLUSIONES

### Fortalezas
- ✅ Arquitectura modular bien definida
- ✅ Separación clara de responsabilidades (Facade/Controller/DAO)
- ✅ Uso de HTTPS obligatorio
- ✅ Sistema de permisos basado en roles
- ✅ Pool de conexiones a base de datos

### Debilidades Críticas
- ❌ Credenciales en texto plano
- ❌ Gestión de sesiones en memoria
- ❌ Falta de tests
- ❌ Configuración de seguridad mejorable
- ❌ Sin logging estructurado

### Evaluación General
**Puntuación: 6.5/10**

El sistema tiene una base arquitectónica sólida pero requiere mejoras urgentes en seguridad y observabilidad antes de considerarse production-ready para un entorno crítico como el Ministerio de Igualdad.

---

**Fecha de Análisis:** 22 de febrero de 2026  
**Analista:** Kiro AI Assistant  
**Versión del Documento:** 1.0
