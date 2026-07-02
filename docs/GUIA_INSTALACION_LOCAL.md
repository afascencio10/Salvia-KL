# Guía de Instalación Local - Sistema SALVIA
## Configuración Completa de Entorno de Desarrollo

---

## 📋 REQUISITOS PREVIOS

### Software Necesario

#### 1. Go (Golang)
```bash
# Verificar si está instalado
go version

# Si no está instalado, descargar desde:
# https://go.dev/dl/

# Versión requerida: Go 1.19 o superior
# Recomendado: Go 1.21+
```

#### 2. PostgreSQL
```bash
# Verificar si está instalado
psql --version

# Windows: Descargar desde
# https://www.postgresql.org/download/windows/

# Versión requerida: PostgreSQL 12+
# Recomendado: PostgreSQL 14+
```

#### 3. Git
```bash
# Verificar si está instalado
git --version

# Windows: Descargar desde
# https://git-scm.com/download/win
```

#### 4. Editor de Código (Opcional pero recomendado)
- Visual Studio Code con extensión Go
- GoLand (IDE específico para Go)

---

## 🗄️ PASO 1: CONFIGURAR BASE DE DATOS

### 1.1 Iniciar PostgreSQL

```bash
# Windows - Iniciar servicio PostgreSQL
# Buscar "Services" en el menú inicio
# Encontrar "postgresql-x64-XX" y hacer clic en "Start"

# O desde línea de comandos (como administrador)
net start postgresql-x64-14
```

### 1.2 Crear Base de Datos

```bash
# Abrir terminal y conectar a PostgreSQL
psql -U postgres

# Dentro de psql, ejecutar:
```

```sql
-- Crear usuario para la aplicación
CREATE USER salvia_admin WITH PASSWORD 'asd876.!@asdSDS5a36Z';

-- Crear base de datos
CREATE DATABASE salvia OWNER salvia_admin;

-- Otorgar privilegios
GRANT ALL PRIVILEGES ON DATABASE salvia TO salvia_admin;

-- Conectar a la base de datos
\c salvia

-- Crear esquema
CREATE SCHEMA IF NOT EXISTS salvia;
CREATE SCHEMA IF NOT EXISTS security;

-- Otorgar permisos en esquemas
GRANT ALL ON SCHEMA salvia TO salvia_admin;
GRANT ALL ON SCHEMA security TO salvia_admin;

-- Salir de psql
\q
```

### 1.3 Restaurar Backup (Si existe)

```bash
# Si tienes un archivo de backup (.sql o .dump)
psql -U salvia_admin -d salvia -f backup.sql

# O con pg_restore si es formato custom
pg_restore -U salvia_admin -d salvia backup.dump
```

### 1.4 Verificar Conexión

```bash
# Probar conexión
psql -U salvia_admin -d salvia -h localhost -p 5432

# Si conecta correctamente, la base de datos está lista
```

---

## 🔧 PASO 2: CONFIGURAR BACKEND (Go)

### 2.1 Clonar/Ubicar el Proyecto

```bash
# Si el proyecto ya está en tu máquina
cd C:\Salvia\MIN_Igualdad\Sistema_Informacio╠ün\Fuentes

# Si necesitas clonarlo desde Git
git clone <url-del-repositorio>
cd Fuentes
```

### 2.2 Verificar Estructura del Proyecto

```bash
# Listar archivos principales
dir

# Deberías ver:
# - main.go
# - go.mod
# - go.work
# - common/
# - salvia/
# - security/
# - frontend/
# - config/
```

### 2.3 Configurar Variables de Entorno

```bash
# Crear archivo .env en la raíz del proyecto
# Fuentes/.env
```

Contenido del archivo `.env`:
```bash
# Base de Datos
DB_HOSTNAME=localhost
DB_PORT=5432
DB_DATABASE=salvia
DB_USER=salvia_admin
DB_PASSWORD=asd876.!@asdSDS5a36Z

# Sesiones
SESSION_SECRET=cambiar-por-clave-aleatoria-segura

# Entorno
ENV=development

# Puerto (opcional)
PORT=443
```

### 2.4 Instalar Dependencias de Go

```bash
# Desde la carpeta Fuentes/
cd C:\Salvia\MIN_Igualdad\Sistema_Informacio╠ün\Fuentes

# Descargar todas las dependencias
go mod download

# Verificar que no hay errores
go mod verify

# Si hay problemas, limpiar caché y reintentar
go clean -modcache
go mod download
```

### 2.5 Compilar el Proyecto

```bash
# Compilar sin ejecutar (para verificar errores)
go build -o salvia.exe main.go

# Si compila sin errores, estás listo
```

---

## 🔐 PASO 3: CONFIGURAR CERTIFICADOS SSL

### 3.1 Verificar Certificados Existentes

```bash
# Verificar si existen los certificados
dir certs\

# Deberías ver:
# - salviaTest.crt
# - salviaTest.key
```

### 3.2 Generar Certificados de Desarrollo (Si no existen)

```bash
# Instalar OpenSSL si no lo tienes
# Windows: https://slproweb.com/products/Win32OpenSSL.html

# Crear carpeta certs si no existe
mkdir certs
cd certs

# Generar certificado autofirmado
openssl req -x509 -newkey rsa:4096 -keyout salviaTest.key -out salviaTest.crt -days 365 -nodes

# Responder las preguntas:
# Country Name: ES
# State: Madrid
# Locality: Madrid
# Organization: Ministerio de Igualdad
# Common Name: localhost
# Email: admin@localhost

cd ..
```

### 3.3 Confiar en el Certificado (Windows)

```bash
# Abrir el archivo .crt
# Doble clic en certs\salviaTest.crt

# Clic en "Install Certificate"
# Seleccionar "Local Machine"
# Seleccionar "Place all certificates in the following store"
# Buscar "Trusted Root Certification Authorities"
# Finalizar

# Esto evitará advertencias de seguridad en el navegador
```

---

## 🚀 PASO 4: EJECUTAR EL BACKEND

### 4.1 Ejecutar en Modo Desarrollo

```bash
# Desde la carpeta Fuentes/
cd C:\Salvia\MIN_Igualdad\Sistema_Informacio╠ün\Fuentes

# Ejecutar directamente con go run
go run main.go

# Deberías ver algo como:
# [GIN-debug] Listening and serving HTTPS on :443
```

### 4.2 Ejecutar Compilado

```bash
# Si ya compilaste el proyecto
.\salvia.exe

# O compilar y ejecutar en un solo paso
go build -o salvia.exe main.go && .\salvia.exe
```

### 4.3 Verificar que el Servidor Está Corriendo

```bash
# Abrir navegador y visitar:
https://localhost

# O con curl
curl -k https://localhost

# Deberías ver la página de inicio o ser redirigido a /static/landing.html
```

### 4.4 Solución de Problemas Comunes

#### Error: "bind: permission denied"
```bash
# Windows requiere permisos de administrador para puerto 443
# Ejecutar terminal como administrador

# O cambiar puerto en el código temporalmente
# common/facades/MainRouter.go
# Cambiar ":443" por ":8443"
```

#### Error: "database connection failed"
```bash
# Verificar que PostgreSQL está corriendo
net start postgresql-x64-14

# Verificar credenciales en config/db_config.json
# O en variables de entorno .env
```

#### Error: "certificate not found"
```bash
# Verificar que los certificados existen
dir certs\

# Si no existen, generarlos (ver Paso 3.2)
```

---

## 🎨 PASO 5: ENTENDER EL FRONTEND

### 5.1 Estructura del Frontend

```
frontend/
├── html/              # Páginas HTML estáticas
├── js/                # JavaScript (Vue.js)
│   ├── vue.js        # Framework Vue.js
│   ├── dao.js        # Acceso a datos (AJAX)
│   ├── formulario.js # Lógica de formularios
│   └── ...
├── css/              # Estilos CSS
├── templates/        # Plantillas Go HTML
│   ├── salvia/      # Templates del módulo Salvia
│   └── security/    # Templates del módulo Security
├── plugins/          # Librerías de terceros
├── images/           # Imágenes
└── landing.html      # Página de inicio
```

### 5.2 Cómo Funciona el Frontend

El frontend es **híbrido**:

1. **Servidor Go renderiza HTML** usando templates
2. **Vue.js maneja interactividad** en el cliente
3. **AJAX** para comunicación con backend

```
Usuario → Navegador
    ↓
    Solicita https://localhost/salvia/caso-victima
    ↓
Backend Go (Gin)
    ↓
    Renderiza template HTML con datos
    ↓
    Envía HTML al navegador
    ↓
Vue.js en el navegador
    ↓
    Añade interactividad (validaciones, AJAX)
    ↓
    Usuario interactúa con formulario
    ↓
    Vue.js envía datos via AJAX
    ↓
Backend Go procesa y responde JSON
```

### 5.3 No Requiere Compilación

```bash
# El frontend NO necesita npm install ni build
# Los archivos se sirven directamente desde /frontend

# El servidor Go sirve archivos estáticos:
# router.Use(static.Serve("/static", static.LocalFile("./frontend", true)))
```

### 5.4 Modificar Frontend

```bash
# Para modificar el frontend, simplemente edita los archivos:

# Editar JavaScript
notepad frontend\js\formulario.js

# Editar CSS
notepad frontend\css\styles.css

# Editar HTML
notepad frontend\html\casos.html

# Refrescar navegador (Ctrl+F5) para ver cambios
# No requiere reiniciar el servidor Go
```

---

## 📚 PASO 6: ENTENDER LA ESTRUCTURA DEL CÓDIGO

### 6.1 Flujo de una Petición HTTP

```
1. Cliente hace petición
   GET https://localhost/salvia/contacto-victima
   
2. Gin Router recibe petición
   common/facades/MainRouter.go
   
3. Middlewares procesan petición
   - limitMiddleware() → Control de concurrencia
   - ForceHTTPS() → Asegurar HTTPS
   - AuthMiddleware() → Verificar autenticación
   
4. Router dirige a Facade correcto
   salvia/facades/MainRouter.go
   → VictimContactGET()
   
5. Facade verifica permisos
   salvia/facades/VictimContactFacade.go
   → CheckPermission()
   
6. Facade llama Controller
   salvia/controllers/VictimContactController.go
   → GetVictimContactByICode()
   
7. Controller valida datos
   → ValidateJSONInput()
   
8. Controller llama DAO
   salvia/dao/VictimContactDAO.go
   → GetVictimContact()
   
9. DAO ejecuta query SQL
   → db.Query()
   
10. DAO devuelve DTO
    → VictimContactDTO
    
11. Controller procesa resultado
    → Formatea JSON
    
12. Facade renderiza respuesta
    → RenderTemplate() o JSON
    
13. Cliente recibe respuesta
    → HTML o JSON
```

### 6.2 Capas del Sistema

```
┌─────────────────────────────────────┐
│         CLIENTE (Navegador)         │
│     HTML + Vue.js + JavaScript      │
└─────────────────────────────────────┘
                 ↕ HTTPS
┌─────────────────────────────────────┐
│      FACADES (Capa Presentación)    │
│  - Maneja HTTP requests/responses   │
│  - Verifica sesiones y permisos     │
│  - Renderiza templates              │
│                                     │
│  Archivos: */facades/*Facade.go     │
└─────────────────────────────────────┘
                 ↕
┌─────────────────────────────────────┐
│    CONTROLLERS (Lógica Negocio)     │
│  - Valida datos de entrada          │
│  - Aplica reglas de negocio         │
│  - Coordina operaciones             │
│                                     │
│  Archivos: */controllers/*Ctrl.go   │
└─────────────────────────────────────┘
                 ↕
┌─────────────────────────────────────┐
│      DAO (Acceso a Datos)           │
│  - Ejecuta queries SQL              │
│  - Mapea DB ↔ DTOs                  │
│  - Gestiona transacciones           │
│                                     │
│  Archivos: */dao/*DAO.go            │
└─────────────────────────────────────┘
                 ↕
┌─────────────────────────────────────┐
│       PostgreSQL Database           │
│  - Almacena datos persistentes      │
└─────────────────────────────────────┘
```

### 6.3 Módulos Principales

#### Módulo COMMON (Compartido)
```
common/
├── config/        # Configuraciones globales
├── controllers/   # Controladores comunes
├── dao/           # DTOs y DAOs comunes
├── db/            # Gestión de conexiones DB
├── facades/       # Router principal y middlewares
└── utils/         # Utilidades (validación, sesiones, etc.)
```

#### Módulo SALVIA (Casos de Víctimas)
```
salvia/
├── config/        # Configuración específica
├── controllers/   # Lógica de negocio de casos
├── dao/           # Modelos de datos de casos
└── facades/       # Handlers HTTP de casos
```

#### Módulo SECURITY (Autenticación)
```
security/
├── config/        # Configuración de seguridad
├── controllers/   # Lógica de autenticación
├── dao/           # Modelos de usuarios y roles
└── facades/       # Handlers de login/logout
```


### 6.4 Archivos Clave para Entender

#### main.go (Punto de Entrada)
```go
package main

func main() {
    // 1. Inicializa router principal
    var router *gin.Engine = common_routers.InitRouter()
    
    // 2. Registra rutas del módulo Salvia
    salvia_facades.StartRouter(router)
    
    // 3. Registra rutas del módulo Security
    security_routers.StartRouter(router)
    
    // 4. Inicia servidor HTTPS
    common_routers.StartRouter()
}
```

#### common/facades/MainRouter.go (Configuración Global)
```go
func InitRouter() *gin.Engine {
    // Crea instancia de Gin
    router = gin.Default()
    
    // Configura sesiones
    store := cookie.NewStore([]byte("secret"))
    router.Use(sessions.Sessions("default_session", store))
    
    // Añade middlewares
    router.Use(limitMiddleware())      // Control concurrencia
    router.Use(ForceHTTPS())           // Forzar HTTPS
    router.Use(AuthMiddleware())       // Autenticación
    
    // Sirve archivos estáticos
    router.Use(static.Serve("/static", static.LocalFile("./frontend", true)))
    
    // Carga templates HTML
    router.LoadHTMLGlob("frontend/templates/**/*")
    
    return router
}
```

#### salvia/facades/MainRouter.go (Rutas del Módulo)
```go
func StartRouter(router *gin.Engine) {
    // Grupo de rutas protegidas
    secRouter := router.Group("/salvia")
    {
        // POST /salvia/contacto-victima
        secRouter.POST("/contacto-victima", VictimContactPOST)
        
        // GET /salvia/contacto-victima
        secRouter.GET("/contacto-victima", VictimContactGET)
        
        // PUT /salvia/contacto-victima/:id/invalidar
        secRouter.PUT("/contacto-victima/:id/invalidar", VictimContactPUT)
    }
    
    // Grupo de rutas públicas (sin autenticación)
    publicRouter := router.Group("/public")
    {
        // POST /public/contacto-victima/nuevo
        publicRouter.POST("/contacto-victima/nuevo", VictimContactPOST_Public)
    }
}
```

#### salvia/dao/VictimContactDAO.go (Modelo de Datos)
```go
// DTO (Data Transfer Object)
type VictimContactDTO struct {
    VictimContactId           uint64    `json:"id"`
    VictimContactICode        string    `json:"icode"`
    VictimContactCreationDate time.Time `json:"creationDate"`
    VictimContactNames        string    `json:"names"`
    VictimContactPhone        string    `json:"phone"`
    // ... más campos
}

// Función para obtener contacto de BD
func GetVictimContact(by By, victimContact *VictimContactDTO, ...) error {
    // Construye query SQL
    query := "SELECT * FROM salvia.victim_contact WHERE ..."
    
    // Ejecuta query
    rows, err := conn.Query(ctx, query, params...)
    
    // Mapea resultado a DTO
    // ...
    
    return err
}
```

---

## 🔍 PASO 7: COMANDOS BÁSICOS DE DESARROLLO

### 7.1 Comandos Go Útiles

```bash
# Ver todas las dependencias
go list -m all

# Actualizar dependencias
go get -u ./...

# Limpiar caché de módulos
go clean -modcache

# Formatear código
go fmt ./...

# Verificar código (linting básico)
go vet ./...

# Ver documentación de un paquete
go doc <paquete>

# Ejemplo: Ver doc de gin
go doc github.com/gin-gonic/gin

# Compilar para producción (optimizado)
go build -ldflags="-s -w" -o salvia.exe main.go
```

### 7.2 Comandos PostgreSQL Útiles

```bash
# Conectar a la base de datos
psql -U salvia_admin -d salvia

# Dentro de psql:

# Listar todas las tablas
\dt salvia.*

# Describir estructura de una tabla
\d salvia.victim_contact

# Ver datos de una tabla
SELECT * FROM salvia.victim_contact LIMIT 10;

# Contar registros
SELECT COUNT(*) FROM salvia.victim_contact;

# Ver usuarios conectados
SELECT * FROM pg_stat_activity;

# Salir
\q
```

### 7.3 Comandos Git Útiles

```bash
# Ver estado del repositorio
git status

# Ver cambios realizados
git diff

# Crear rama para desarrollo
git checkout -b feature/mi-nueva-funcionalidad

# Guardar cambios
git add .
git commit -m "Descripción de cambios"

# Ver historial
git log --oneline

# Deshacer cambios no guardados
git checkout -- <archivo>

# Volver a commit anterior
git reset --hard HEAD~1
```

---

## 🧪 PASO 8: PROBAR LA APLICACIÓN

### 8.1 Acceder a la Aplicación

```bash
# 1. Asegurarse que el servidor está corriendo
go run main.go

# 2. Abrir navegador
https://localhost

# 3. Deberías ver la página de inicio (landing.html)
```

### 8.2 Probar Login

```bash
# URL de login
https://localhost/seguridad/login

# Si no tienes usuario, necesitas crear uno en la BD:
```

```sql
-- Conectar a PostgreSQL
psql -U salvia_admin -d salvia

-- Crear usuario de prueba (ajustar según tu esquema)
INSERT INTO security.general_user (
    general_user_i_code,
    general_user_names,
    general_user_last_names,
    general_user_email,
    general_user_password,  -- Debe estar hasheado
    general_user_role,
    general_user_status
) VALUES (
    'USR001',
    'Admin',
    'Test',
    'admin@test.com',
    'hash_de_contraseña',  -- Usar bcrypt
    'ad',
    'A'
);
```

### 8.3 Probar Endpoints con Curl

```bash
# Probar endpoint público
curl -k https://localhost/public/contacto-victima/nuevo

# Probar endpoint protegido (sin autenticación)
curl -k https://localhost/salvia/contacto-victima
# Debería redirigir a login

# Probar con sesión (después de login)
curl -k -b cookies.txt https://localhost/salvia/contacto-victima
```

### 8.4 Probar con Postman

```
1. Descargar Postman: https://www.postman.com/downloads/

2. Desactivar verificación SSL:
   Settings → General → SSL certificate verification: OFF

3. Crear colección "SALVIA API"

4. Añadir requests:
   - GET https://localhost/salvia/contacto-victima
   - POST https://localhost/salvia/contacto-victima
   - PUT https://localhost/salvia/contacto-victima/:id/invalidar

5. Configurar cookies/sesiones en Postman
```

---

## 🐛 PASO 9: DEBUGGING

### 9.1 Debugging con VS Code

**Instalar extensión Go:**
1. Abrir VS Code
2. Ir a Extensions (Ctrl+Shift+X)
3. Buscar "Go"
4. Instalar extensión oficial de Go

**Configurar launch.json:**
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/main.go",
            "env": {
                "DB_HOSTNAME": "localhost",
                "DB_PORT": "5432",
                "DB_DATABASE": "salvia",
                "DB_USER": "salvia_admin",
                "DB_PASSWORD": "asd876.!@asdSDS5a36Z"
            }
        }
    ]
}
```

**Usar breakpoints:**
1. Hacer clic en el margen izquierdo del código (aparece punto rojo)
2. Presionar F5 para iniciar debugging
3. El código se detendrá en los breakpoints
4. Inspeccionar variables en el panel izquierdo

### 9.2 Logs de Debugging

```go
// Añadir logs temporales en el código
import "log"

func VictimContactPOST(c *gin.Context) {
    log.Println("=== DEBUG: Entrando a VictimContactPOST ===")
    
    session := sessions.Default(c)
    sessionID := session.Get("userData")
    log.Printf("DEBUG: Session ID: %v\n", sessionID)
    
    // ... resto del código
}
```

### 9.3 Ver Logs del Servidor

```bash
# Los logs aparecen en la terminal donde ejecutaste go run
# Buscar líneas como:
# [GIN] 2024/02/23 - 10:30:45 | 200 | 45.2ms | 127.0.0.1 | GET "/salvia/contacto-victima"

# Para guardar logs en archivo:
go run main.go > logs.txt 2>&1

# Ver logs en tiempo real:
tail -f logs.txt  # Linux/Mac
Get-Content logs.txt -Wait  # Windows PowerShell
```

---

## 📝 PASO 10: MODIFICAR EL CÓDIGO

### 10.1 Añadir Nueva Ruta

**1. Definir ruta en MainRouter.go:**
```go
// salvia/facades/MainRouter.go
func StartRouter(router *gin.Engine) {
    secRouter := router.Group("/salvia")
    {
        // ... rutas existentes
        
        // ✅ NUEVA RUTA
        secRouter.GET("/mi-nueva-ruta", MiNuevaRutaGET)
    }
}
```

**2. Crear handler en Facade:**
```go
// salvia/facades/MiNuevaRutaFacade.go
package salvia_facades

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func MiNuevaRutaGET(c *gin.Context) {
    // Verificar autenticación
    session := sessions.Default(c)
    sessionID := session.Get("userData").(string)
    
    // Verificar permisos
    s, _ := utils.GetCommonSession(sessionID)
    if !utils.CheckPermission(salvia_config.PermissionsByRole, "mi_permiso", s.CurrentRole, c) {
        return
    }
    
    // Llamar controller
    code, res := salvia_ctrl.MiNuevaFuncion()
    
    // Devolver respuesta
    c.JSON(code, gin.H{"data": res})
}
```

**3. Crear controller:**
```go
// salvia/controllers/MiNuevoController.go
package salvia_ctrl

func MiNuevaFuncion() (int, string) {
    // Lógica de negocio
    
    // Llamar DAO si es necesario
    data, err := salvia_daos.GetMisDatos()
    
    if err != nil {
        return 500, "Error"
    }
    
    return 200, data
}
```

**4. Reiniciar servidor:**
```bash
# Detener servidor (Ctrl+C)
# Volver a ejecutar
go run main.go
```

### 10.2 Modificar Template HTML

```bash
# Editar template
notepad frontend\templates\salvia\victim_contact\get_victim_contact.html

# Añadir nuevo campo
<div class="form-group">
    <label>Nuevo Campo:</label>
    <input type="text" v-model="nuevoCampo" class="form-control">
</div>

# Refrescar navegador (Ctrl+F5)
# No requiere reiniciar servidor
```

### 10.3 Añadir Nueva Tabla en BD

```sql
-- Conectar a PostgreSQL
psql -U salvia_admin -d salvia

-- Crear nueva tabla
CREATE TABLE salvia.mi_nueva_tabla (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    fecha_creacion TIMESTAMP DEFAULT NOW()
);

-- Otorgar permisos
GRANT ALL ON salvia.mi_nueva_tabla TO salvia_admin;
GRANT ALL ON SEQUENCE salvia.mi_nueva_tabla_id_seq TO salvia_admin;
```

---

## 🔄 PASO 11: WORKFLOW DE DESARROLLO

### 11.1 Ciclo de Desarrollo Típico

```
1. Crear rama de feature
   git checkout -b feature/nueva-funcionalidad

2. Modificar código
   - Editar archivos .go
   - Editar templates HTML
   - Editar JavaScript

3. Probar localmente
   go run main.go
   Abrir https://localhost

4. Verificar errores
   go vet ./...
   go fmt ./...

5. Commit cambios
   git add .
   git commit -m "Añadir nueva funcionalidad"

6. Push a repositorio
   git push origin feature/nueva-funcionalidad

7. Crear Pull Request
   (En GitHub/GitLab)

8. Code Review
   (Otro desarrollador revisa)

9. Merge a main
   (Después de aprobación)

10. Deploy a producción
    (Proceso automatizado o manual)
```

### 11.2 Buenas Prácticas

```bash
# ✅ HACER:
- Probar cambios localmente antes de commit
- Escribir mensajes de commit descriptivos
- Mantener commits pequeños y enfocados
- Actualizar documentación cuando cambias funcionalidad
- Usar ramas para cada feature
- Hacer pull de main frecuentemente

# ❌ NO HACER:
- Commitear credenciales o secretos
- Hacer commits gigantes con muchos cambios
- Modificar directamente la rama main
- Dejar código comentado sin explicación
- Ignorar errores de go vet
```

---

## 📚 PASO 12: RECURSOS Y DOCUMENTACIÓN

### 12.1 Documentación Oficial

**Go:**
- https://go.dev/doc/
- https://go.dev/tour/
- https://pkg.go.dev/

**Gin Framework:**
- https://gin-gonic.com/docs/
- https://github.com/gin-gonic/gin

**PostgreSQL:**
- https://www.postgresql.org/docs/
- https://www.postgresqltutorial.com/

**Vue.js:**
- https://vuejs.org/guide/
- https://v3.vuejs.org/

### 12.2 Tutoriales Recomendados

**Go Básico:**
- Tour of Go: https://go.dev/tour/
- Go by Example: https://gobyexample.com/

**Gin Framework:**
- https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/

**PostgreSQL con Go:**
- https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/

### 12.3 Herramientas Útiles

**Gestión de BD:**
- pgAdmin: https://www.pgadmin.org/
- DBeaver: https://dbeaver.io/

**Testing de APIs:**
- Postman: https://www.postman.com/
- Insomnia: https://insomnia.rest/

**Debugging:**
- Delve (Go debugger): https://github.com/go-delve/delve

---

## ❓ PASO 13: SOLUCIÓN DE PROBLEMAS COMUNES

### Problema 1: "go: command not found"
```bash
# Solución: Añadir Go al PATH

# Windows:
# 1. Buscar "Environment Variables" en el menú inicio
# 2. Editar "Path" en System Variables
# 3. Añadir: C:\Go\bin
# 4. Reiniciar terminal

# Verificar:
go version
```

### Problema 2: "psql: command not found"
```bash
# Solución: Añadir PostgreSQL al PATH

# Windows:
# Añadir a PATH: C:\Program Files\PostgreSQL\14\bin

# Verificar:
psql --version
```

### Problema 3: "connection refused" al conectar a BD
```bash
# Verificar que PostgreSQL está corriendo
net start postgresql-x64-14

# Verificar puerto
netstat -an | findstr 5432

# Verificar credenciales en config/db_config.json
```

### Problema 4: "certificate signed by unknown authority"
```bash
# Solución: Confiar en el certificado autofirmado
# Ver Paso 3.3

# O usar curl con -k (inseguro, solo desarrollo)
curl -k https://localhost
```

### Problema 5: "port 443 already in use"
```bash
# Ver qué proceso usa el puerto
netstat -ano | findstr :443

# Matar proceso (reemplazar PID)
taskkill /PID <numero> /F

# O cambiar puerto en el código temporalmente
```

### Problema 6: Cambios en código no se reflejan
```bash
# Detener servidor (Ctrl+C)
# Limpiar build cache
go clean -cache

# Recompilar y ejecutar
go run main.go

# Para frontend, hacer hard refresh
# Ctrl+Shift+R o Ctrl+F5
```

---

## 🎯 RESUMEN RÁPIDO

### Comandos Esenciales

```bash
# 1. Iniciar PostgreSQL
net start postgresql-x64-14

# 2. Ir a carpeta del proyecto
cd C:\Salvia\MIN_Igualdad\Sistema_Informacio╠ün\Fuentes

# 3. Descargar dependencias (primera vez)
go mod download

# 4. Ejecutar aplicación
go run main.go

# 5. Abrir navegador
https://localhost

# 6. Detener aplicación
Ctrl+C
```

### Estructura Rápida

```
Fuentes/
├── main.go              # ← INICIO AQUÍ
├── common/facades/      # ← Router principal
├── salvia/facades/      # ← Rutas de casos
├── security/facades/    # ← Rutas de login
├── frontend/            # ← HTML/CSS/JS
└── config/              # ← Configuración BD
```

### Flujo de Petición

```
Cliente → Gin Router → Middleware → Facade → Controller → DAO → PostgreSQL
```

---

## 📞 AYUDA ADICIONAL

Si tienes problemas:

1. **Revisar logs** en la terminal donde ejecutaste `go run main.go`
2. **Verificar BD** con `psql -U salvia_admin -d salvia`
3. **Verificar certificados** en carpeta `certs/`
4. **Limpiar caché** con `go clean -cache -modcache`
5. **Consultar documentación** de Go, Gin, PostgreSQL

---

**Documento creado:** 22 de febrero de 2026  
**Versión:** 1.0  
**Autor:** Andres Roldan
