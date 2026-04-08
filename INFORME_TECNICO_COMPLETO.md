# INFORME TÉCNICO COMPLETO - SISTEMA SALVIA
## Ministerio de Igualdad - Sistema de Atención a Víctimas de Violencia de Género

**Fecha de Análisis:** Diciembre 2024  
**Versión Go:** 1.19  
**Framework Backend:** Gin  
**Framework Frontend:** Vue.js 3 + HTML/JavaScript  
**Base de Datos:** PostgreSQL 12+

---

## TABLA DE CONTENIDOS

1. [Arquitectura de Solución](#1-arquitectura-de-solución)
2. [Análisis de Archivos Más Grandes y Complejos](#2-análisis-de-archivos-más-grandes-y-complejos)
3. [Deuda Técnica](#3-deuda-técnica)
4. [Vulnerabilidades de Seguridad](#4-vulnerabilidades-de-seguridad)
5. [Resumen Ejecutivo](#5-resumen-ejecutivo)

---

## 1. ARQUITECTURA DE SOLUCIÓN

### 1.1 Patrón Arquitectónico Principal

**Patrón Identificado:** Arquitectura Monolítica Modular con patrón MVC adaptado (Facades-Controllers-DAOs)

**Características:**
- Monolito modular dividido en 3 módulos principales: `common`, `salvia`, `security`
- Separación de responsabilidades por capas
- Renderizado híbrido: Server-side (Go templates) + Client-side (Vue.js)
- Autenticación basada en sesiones con cookies
- Pool de conexiones a PostgreSQL (80 conexiones concurrentes)

### 1.2 Diagrama de Capas (ASCII)

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CAPA DE PRESENTACIÓN                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │   Browser    │  │   Vue.js 3   │  │  Go HTML Templates       │  │
│  │   (HTTPS)    │  │  Components  │  │  (Server-side render)    │  │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ▼ HTTPS
┌─────────────────────────────────────────────────────────────────────┐
│                   CAPA DE MIDDLEWARES                               │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  limitMiddleware() → ForceHTTPS() → sessions.Sessions()     │   │
│  │  → AuthMiddleware() → Handler                                │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   CAPA DE FACADES (HTTP Handlers)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │   common/    │  │   salvia/    │  │   security/              │  │
│  │   facades    │  │   facades    │  │   facades                │  │
│  │              │  │              │  │                          │  │
│  │ MainRouter   │  │ VictimCase   │  │ GeneralUser              │  │
│  │              │  │ VictimContact│  │ City/Town                │  │
│  │              │  │ Alert        │  │                          │  │
│  │              │  │ Moment       │  │                          │  │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   CAPA DE CONTROLLERS (Lógica de Negocio)           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │   common/    │  │   salvia/    │  │   security/              │  │
│  │ controllers  │  │ controllers  │  │ controllers              │  │
│  │              │  │              │  │                          │  │
│  │ General      │  │ VictimCase   │  │ GeneralUser              │  │
│  │ Persistence  │  │ VictimContact│  │ City/Department/Town     │  │
│  │              │  │ Alert        │  │ Role                     │  │
│  │              │  │ Entity       │  │ ResetPassword            │  │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   CAPA DE DAOs (Acceso a Datos)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │   common/    │  │   salvia/    │  │   security/              │  │
│  │     dao      │  │     dao      │  │     dao                  │  │
│  │              │  │              │  │                          │  │
│  │ GenericDAO   │  │ VictimCase   │  │ GeneralUser              │  │
│  │              │  │ VictimContact│  │ GeneralUserProfile       │  │
│  │              │  │ Alert        │  │ City/Department/Town     │  │
│  │              │  │ Entity       │  │ Role/RelRoleGeneralUser  │  │
│  │              │  │ Moment       │  │ EMail/PhoneNumber        │  │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   CAPA DE CONEXIÓN                                  │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │         CONNECTION POOL (80 conexiones)                      │   │
│  │  GetConnection() → Ejecuta Query → ReleaseConnection()      │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   CAPA DE PERSISTENCIA                              │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              PostgreSQL Database                             │   │
│  │  • Schema: security  • Schema: salvia                        │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### 1.3 Flujo Completo de una Request HTTP

#### Backend Go:

```
1. CLIENTE ENVÍA REQUEST
   │
   ├─▶ POST https://localhost/salvia/contacto-victima
   │   Headers: { Content-Type: application/json, Cookie: session_id=... }
   │   Body: { names: "María", phone: "3001234567", ... }
   │
2. GIN ROUTER RECIBE
   │
   ├─▶ Middlewares Pipeline:
   │   │
   │   ├─▶ limitMiddleware()
   │   │   └─▶ Control de concurrencia (máx 9999 peticiones)
   │   │       ✓ Permite si < 9999
   │   │
   │   ├─▶ ForceHTTPS()
   │   │   └─▶ Verifica protocolo HTTPS
   │   │       ✓ Ya es HTTPS, continúa
   │   │
   │   ├─▶ sessions.Sessions()
   │   │   └─▶ Carga sesión desde cookie
   │   │       ✓ Sesión válida encontrada
   │   │
   │   └─▶ AuthMiddleware()
   │       └─▶ Verifica autenticación
   │           ├─▶ Lee cookie session_id
   │           ├─▶ Busca en utils.GetCommonSession(sessionID)
   │           ├─▶ Si existe → Continúa
   │           └─▶ Si no existe → Redirige a /static/landing.html
   │
3. ROUTING
   │
   ├─▶ MainRouter identifica ruta
   │   └─▶ /salvia/contacto-victima → VictimContactPOST()
   │
4. FACADE PROCESA (VictimContactFacade.go)
   │
   ├─▶ VictimContactPOST(c *gin.Context)
   │   ├─▶ Obtiene sesión: session.Get("userData")
   │   ├─▶ Recupera CommonSession: utils.GetCommonSession(sessionID)
   │   ├─▶ Verifica permiso: utils.CheckPermission(
   │   │       salvia_config.PermissionsByRole,
   │   │       "set_victim_contact",
   │   │       s.CurrentRole,
   │   │       c
   │   │   )
   │   │   └─▶ Si no tiene permiso → HTTP 403 Forbidden
   │   │
   │   ├─▶ Lee body: buf.ReadFrom(c.Request.Body)
   │   └─▶ Llama a Controller
   │
5. CONTROLLER VALIDA (VictimContactController.go)
   │
   ├─▶ SetVictimContact(dataInput, module, connData, dbClientConfig, dbServerConfig)
   │   ├─▶ Parsea JSON a DTO
   │   ├─▶ Valida campos:
   │   │   ├─▶ names: 3-32 caracteres ✓
   │   │   ├─▶ phone: 6-10 dígitos ✓
   │   │   ├─▶ email: formato válido ✓
   │   │   └─▶ ... más validaciones
   │   │
   │   ├─▶ Aplica reglas de negocio:
   │   │   ├─▶ Calcula edad desde birthDate
   │   │   ├─▶ Genera ICode único
   │   │   └─▶ Establece status inicial
   │   │
   │   └─▶ Llama a DAO
   │
6. DAO PERSISTE (VictimContactDAO.go)
   │
   ├─▶ SetVictimContact(dto, connData, dbClientConfig, dbServerConfig)
   │   ├─▶ Obtiene conexión: db.GetConnection()
   │   ├─▶ Construye query SQL:
   │   │   INSERT INTO salvia.victim_contact (
   │   │       victim_contact_names,
   │   │       victim_contact_phone,
   │   │       ...
   │   │   ) VALUES ($1, $2, ...)
   │   │
   │   ├─▶ Ejecuta query: connData.Conn.QueryRow(ctx, query, args...)
   │   ├─▶ Obtiene ID generado: Scan(&dto.VictimContactId)
   │   ├─▶ Libera conexión: db.ReleaseConnection(connData)
   │   └─▶ Retorna DTO con ID
   │
7. RESPUESTA ASCIENDE
   │
   ├─▶ DAO → Controller
   │   └─▶ Controller formatea respuesta JSON:
   │       { "success": true, "data": { "icode": "VCT-2024-001", ... } }
   │       └─▶ Controller → Facade
   │           └─▶ Facade envía respuesta HTTP:
   │               c.DataFromReader(code, len(res), gin.MIMEJSON, ...)
   │
8. CLIENTE RECIBE
   │
   └─▶ HTTP 200 OK
       Body: { "success": true, "data": { ... } }
```

#### Frontend Vue.js:

```
1. USUARIO INTERACTÚA
   │
   ├─▶ Vue.js Component captura evento (v-on:click="submitForm")
   │
2. COMPONENT PREPARA REQUEST
   │
   ├─▶ methods: {
   │       submitForm() {
   │           newEntity(
   │               '/salvia/contacto-victima',
   │               this.victimContact,
   │               this.handleResponse
   │           )
   │       }
   │   }
   │
3. DAO.JS ENVÍA REQUEST (dao.js)
   │
   ├─▶ newEntity(url, obj, callback)
   │   └─▶ setData(url, obj, callback, "POST")
   │       ├─▶ Muestra loading overlay
   │       ├─▶ Crea XMLHttpRequest
   │       ├─▶ xhr.open("POST", url)
   │       ├─▶ xhr.setRequestHeader("Content-Type", "application/json")
   │       ├─▶ xhr.send(JSON.stringify(obj))
   │       └─▶ Espera respuesta
   │
4. RECIBE RESPUESTA
   │
   ├─▶ xhr.onload = function() {
   │       if (xhr.status === 200) {
   │           // Muestra success overlay
   │           response = JSON.parse(xhr.responseText)
   │           callback(xhr.status, response)
   │       } else if (xhr.status === 400) {
   │           // Muestra fail overlay
   │           callback(xhr.status, response)
   │       }
   │   }
   │
5. COMPONENT ACTUALIZA UI
   │
   └─▶ handleResponse(status, response) {
           if (status === 200) {
               // Actualiza datos reactivos
               this.victimContacts.push(response.data)
               // Muestra mensaje de éxito
           } else {
               // Muestra errores de validación
           }
       }
```

### 1.4 Módulos Principales y Responsabilidades

#### Módulo `common`
**Responsabilidades:**
- Configuración global del servidor (MainRouter.go)
- Middlewares de autenticación y seguridad
- Gestión de conexiones a base de datos (PostgresConnection.go)
- Utilidades compartidas (CommonSession, FileManager, EMailManager)
- Validaciones genéricas (EntityValidator)
- Configuración de enums y locales

**Archivos clave:**
- `common/facades/MainRouter.go` (172 líneas)
- `common/db/PostgresConnection.go`
- `common/utils/CommonSession.go`
- `common/utils/EntityValidator.go`

#### Módulo `salvia`
**Responsabilidades:**
- Gestión de casos de víctimas (VictimCase)
- Gestión de contactos de víctimas (VictimContact)
- Gestión de alertas (Alert)
- Gestión de momentos del proceso (Moment)
- Gestión de entidades de atención (Entity, EntityBranch)
- Asignación de operadores
- Generación de reportes

**Archivos clave:**
- `salvia/facades/VictimCaseFacade.go` (695 líneas) ⚠️
- `salvia/controllers/VictimCaseController.go` (1779 líneas) ⚠️⚠️⚠️
- `salvia/dao/VictimCaseDAO.go`
- `salvia/facades/VictimContactFacade.go` (274 líneas)

#### Módulo `security`
**Responsabilidades:**
- Autenticación y autorización de usuarios
- Gestión de usuarios (GeneralUser)
- Gestión de roles y permisos
- Gestión de ubicaciones geográficas (Country, Department, City, Town)
- Reseteo de contraseñas
- Gestión de perfiles de usuario

**Archivos clave:**
- `security/facades/GeneralUserFacade.go` (530 líneas) ⚠️
- `security/controllers/GeneralUserController.go` (1582 líneas) ⚠️⚠️⚠️
- `security/dao/GeneralUserDAO.go`
- `security/config/Menu.go`


### 1.5 Dependencias entre Paquetes

#### Dependencias Go (go.mod):

```
bitsflow (main)
├── common
│   ├── config
│   ├── controllers
│   ├── dao
│   ├── db
│   ├── facades
│   └── utils
├── salvia
│   ├── config → common/config
│   ├── controllers → common/db, common/utils, salvia/dao, security/dao
│   ├── dao → common/config, common/dao, common/db, security/dao
│   └── facades → common/facades, common/utils, salvia/config, salvia/controllers
└── security
    ├── config → security/dao
    ├── controllers → common/db, common/utils, security/dao
    ├── dao → common/config, common/dao, common/db
    └── facades → common/facades, common/utils, security/config, security/controllers
```

**Dependencias Externas Principales:**
- `github.com/gin-gonic/gin v1.9.1` - Framework HTTP
- `github.com/jackc/pgx/v4 v4.18.3` - Driver PostgreSQL
- `github.com/gin-contrib/sessions v0.0.5` - Gestión de sesiones
- `github.com/dchest/captcha v1.0.0` - Generación de captchas
- `github.com/google/uuid v1.6.0` - Generación de UUIDs

#### Dependencias Frontend:

**JavaScript:**
- Vue.js 3 (vue.js)
- Bootstrap 4 (bootstrap.min.js)
- Leaflet (mapas)
- DataTables (tablas interactivas)
- jQuery (plugins)

**CSS:**
- Bootstrap 4
- FontAwesome
- Leaflet CSS
- Estilos personalizados

### 1.6 Decisiones Arquitectónicas Relevantes

#### ✅ Decisiones Acertadas:

1. **Separación modular por dominio** (common, salvia, security)
2. **Pool de conexiones** (80 conexiones) para optimizar acceso a BD
3. **Middleware pipeline** para separación de concerns
4. **Uso de prepared statements** en la mayoría de queries SQL
5. **Validación en múltiples capas** (Facade, Controller, DAO)

#### ⚠️ Decisiones Cuestionables:

1. **Sesiones en memoria** (no persistentes, se pierden al reiniciar)
   - **Impacto:** Usuarios deben re-autenticarse tras cada reinicio del servidor
   - **Alternativa:** Redis o base de datos para sesiones

2. **Renderizado híbrido** (Server-side + Client-side)
   - **Impacto:** Complejidad adicional, duplicación de lógica
   - **Alternativa:** SPA completo con Vue.js o SSR puro

3. **Límite de concurrencia hardcodeado** (9999 peticiones)
   - **Impacto:** No configurable, puede ser insuficiente o excesivo
   - **Alternativa:** Configuración externa

4. **Ausencia de caché** para datos estáticos (departamentos, ciudades)
   - **Impacto:** Queries repetitivas a BD
   - **Alternativa:** Caché en memoria o Redis

5. **Manejo de errores inconsistente**
   - **Impacto:** Dificulta debugging y monitoreo
   - **Alternativa:** Logging estructurado (logrus, zap)

6. **Sin versionado de API**
   - **Impacto:** Dificulta evolución sin breaking changes
   - **Alternativa:** /api/v1/, /api/v2/

---

## 2. ANÁLISIS DE ARCHIVOS MÁS GRANDES Y COMPLEJOS

### 2.1 Ranking por Tamaño y Complejidad

| # | Archivo | Líneas | Funciones | Complejidad | Prioridad |
|---|---------|--------|-----------|-------------|-----------|
| 1 | `salvia/controllers/VictimCaseController.go` | 1779 | 26 | **CRÍTICA** | 🔴 ALTA |
| 2 | `security/controllers/GeneralUserController.go` | 1582 | 19 | **CRÍTICA** | 🔴 ALTA |
| 3 | `salvia/facades/VictimCaseFacade.go` | 695 | 7 | **ALTA** | 🟠 MEDIA |
| 4 | `security/facades/GeneralUserFacade.go` | 530 | 11 | **ALTA** | 🟠 MEDIA |
| 5 | `salvia/facades/VictimContactFacade.go` | 274 | 7 | **MEDIA** | 🟡 BAJA |

### 2.2 Análisis Detallado

---

#### 🔴 ARCHIVO #1: `salvia/controllers/VictimCaseController.go`

**Métricas:**
- **Líneas:** 1779
- **Funciones:** 26
- **Complejidad Ciclomática Promedio:** 15-20 (ALTA)
- **Responsabilidades:** 8+ (VIOLACIÓN SRP)

**Structs/Interfaces Principales:**
- Ninguno (solo funciones)

**Funciones Principales:**

| Función | Líneas | Complejidad | Responsabilidades |
|---------|--------|-------------|-------------------|
| `SetVictimCase()` | ~405 | 25+ | Crear caso, validar, crear usuario, asignar operador, enviar email |
| `UpdateVictimCase()` | ~397 | 20+ | Actualizar caso, validar, gestionar momentos, actualizar usuario |
| `GetVictimCaseByICode()` | ~140 | 8 | Obtener caso por ID |
| `UpdateGeneralUserByICode()` | ~393 | 18+ | Actualizar usuario completo |
| `AssignOperators()` | ~80 | 10 | Reasignar operadores en masa |

**Responsabilidades Actuales:**
1. ✅ Validación de datos de entrada
2. ✅ Aplicación de reglas de negocio
3. ❌ Creación de usuarios (debería estar en security)
4. ❌ Envío de emails (debería estar en common/utils)
5. ❌ Generación de contraseñas (debería estar en security)
6. ❌ Gestión de transacciones complejas
7. ❌ Formateo de respuestas JSON
8. ❌ Cálculo de edad y validaciones de fecha

**Responsabilidades que NO debería tener (Violaciones SRP):**
- ❌ Creación de usuarios del sistema (líneas 100-250)
- ❌ Envío de emails de notificación (líneas 300-350)
- ❌ Generación de nombres de usuario y contraseñas (líneas 1722-1743)
- ❌ Gestión directa de transacciones de BD
- ❌ Formateo de respuestas HTTP

**Funciones con Mayor Complejidad Ciclomática:**

1. **`SetVictimCase()`** (Líneas 26-431)
   - **Complejidad:** ~25 ramas
   - **Problemas:**
     - 405 líneas (debería ser <50)
     - 15+ niveles de indentación
     - Maneja creación de usuario, caso, momentos, emails
     - Transacción manual sin rollback automático
   - **Código problemático:**
   ```go
   // Líneas 100-250: Creación de usuario embebida
   if vCase.VictimCaseNewUser.GeneralUserLogin != "" {
       // 150 líneas de lógica de creación de usuario
       // Debería estar en security_ctrl.CreateUser()
   }
   ```

2. **`UpdateVictimCase()`** (Líneas 431-828)
   - **Complejidad:** ~20 ramas
   - **Problemas:**
     - 397 líneas
     - Lógica compleja de diff de momentos
     - Múltiples queries sin transacción explícita
   - **Código problemático:**
   ```go
   // Líneas 600-700: Lógica de diff de momentos
   momentsToDelete, momentsToInsert := cleanMoments(...)
   // Debería estar en un servicio separado
   ```

3. **`LoginGeneralUser()`** (security/controllers/GeneralUserController.go, Líneas 28-176)
   - **Complejidad:** ~12 ramas
   - **Problemas:**
     - Validación de captcha embebida
     - Verificación de contraseña sin rate limiting
     - Múltiples queries sin transacción

**Dependencias Externas:**
- `common/db` - Conexiones a BD
- `common/utils` - Utilidades (sesiones, validaciones, emails)
- `salvia/dao` - Acceso a datos de salvia
- `security/dao` - Acceso a datos de security (ACOPLAMIENTO)
- `security/controllers` - Controladores de security (ACOPLAMIENTO FUERTE)

**Propuesta de Refactoring:**

```
VictimCaseController.go (1779 líneas)
│
├─▶ VictimCaseService.go (300 líneas)
│   ├─▶ CreateVictimCase()
│   ├─▶ UpdateVictimCase()
│   ├─▶ GetVictimCase()
│   └─▶ ApproveVictimCase()
│
├─▶ VictimCaseValidator.go (200 líneas)
│   ├─▶ ValidateVictimCaseData()
│   ├─▶ ValidateMoments()
│   └─▶ ValidateEntityBranches()
│
├─▶ VictimCaseUserService.go (250 líneas)
│   ├─▶ CreateUserForVictim()
│   ├─▶ GenerateCredentials()
│   └─▶ SendCredentialsEmail()
│
├─▶ VictimCaseMomentService.go (200 líneas)
│   ├─▶ SyncMoments()
│   ├─▶ AddMoment()
│   ├─▶ RemoveMoment()
│   └─▶ UpdateMoment()
│
├─▶ VictimCaseOperatorService.go (150 líneas)
│   ├─▶ AssignOperator()
│   ├─▶ ReassignOperators()
│   └─▶ GetOperatorCases()
│
└─▶ VictimCaseReportService.go (200 líneas)
    ├─▶ GenerateReport()
    ├─▶ FilterCases()
    └─▶ ExportReport()
```

**Esfuerzo Estimado:** 40-60 horas

---

#### 🔴 ARCHIVO #2: `security/controllers/GeneralUserController.go`

**Métricas:**
- **Líneas:** 1582
- **Funciones:** 19
- **Complejidad Ciclomática Promedio:** 12-18 (ALTA)
- **Responsabilidades:** 7+ (VIOLACIÓN SRP)

**Funciones Principales:**

| Función | Líneas | Complejidad | Responsabilidades |
|---------|--------|-------------|-------------------|
| `SetGeneralUser()` | ~389 | 20+ | Crear usuario, validar, gestionar emails/teléfonos/roles |
| `UpdateGeneralUserByICode()` | ~393 | 18+ | Actualizar usuario completo |
| `LoginGeneralUser()` | ~148 | 12 | Login, validar captcha, verificar contraseña |
| `UpdateGeneralUserByReset()` | ~144 | 10 | Resetear contraseña |

**Responsabilidades Actuales:**
1. ✅ Validación de datos de usuario
2. ✅ Autenticación (login)
3. ❌ Gestión de emails (debería ser EmailService)
4. ❌ Gestión de teléfonos (debería ser PhoneService)
5. ❌ Gestión de roles (debería ser RoleService)
6. ❌ Envío de emails de notificación
7. ❌ Validación de captcha (debería ser middleware)

**Funciones con Mayor Complejidad:**

1. **`SetGeneralUser()`** (Líneas 320-709)
   - **Complejidad:** ~20 ramas
   - **Problemas:**
     - 389 líneas
     - Gestiona emails, teléfonos, roles en el mismo método
     - Múltiples queries sin transacción explícita
   - **Código problemático:**
   ```go
   // Líneas 400-500: Gestión de emails embebida
   for _, mail := range usr.GeneralUserEMails {
       // Validación y persistencia de emails
       // Debería estar en EmailService
   }
   ```

2. **`UpdateGeneralUserByICode()`** (Líneas 889-1282)
   - **Complejidad:** ~18 ramas
   - **Problemas:**
     - 393 líneas
     - Lógica de diff de emails/teléfonos/roles
     - Sin transacción explícita

**Propuesta de Refactoring:**

```
GeneralUserController.go (1582 líneas)
│
├─▶ UserService.go (250 líneas)
│   ├─▶ CreateUser()
│   ├─▶ UpdateUser()
│   ├─▶ DeleteUser()
│   └─▶ GetUser()
│
├─▶ AuthService.go (200 líneas)
│   ├─▶ Login()
│   ├─▶ Logout()
│   ├─▶ ValidateCredentials()
│   └─▶ ValidateCaptcha()
│
├─▶ UserEmailService.go (150 líneas)
│   ├─▶ AddEmail()
│   ├─▶ RemoveEmail()
│   ├─▶ UpdateEmail()
│   └─▶ SyncEmails()
│
├─▶ UserPhoneService.go (150 líneas)
│   ├─▶ AddPhone()
│   ├─▶ RemovePhone()
│   ├─▶ UpdatePhone()
│   └─▶ SyncPhones()
│
├─▶ UserRoleService.go (150 líneas)
│   ├─▶ AssignRole()
│   ├─▶ RevokeRole()
│   ├─▶ UpdateRoles()
│   └─▶ SyncRoles()
│
└─▶ PasswordResetService.go (200 líneas)
    ├─▶ RequestReset()
    ├─▶ ValidateResetToken()
    ├─▶ ResetPassword()
    └─▶ SendResetEmail()
```

**Esfuerzo Estimado:** 35-50 horas

---

#### 🟠 ARCHIVO #3: `salvia/facades/VictimCaseFacade.go`

**Métricas:**
- **Líneas:** 695
- **Funciones:** 7
- **Complejidad Ciclomática Promedio:** 10-15 (MEDIA-ALTA)

**Funciones Principales:**

| Función | Líneas | Complejidad | Responsabilidades |
|---------|--------|-------------|-------------------|
| `VictimCaseGET()` | ~220 | 15+ | Obtener casos, filtrar por rol, renderizar |
| `VictimCasePUT_GET()` | ~156 | 12 | Preparar datos para actualización |
| `VictimCasePOST_GET()` | ~132 | 10 | Preparar datos para creación |

**Problemas Identificados:**
1. Funciones demasiado largas (>100 líneas)
2. Lógica de negocio en la capa de presentación
3. Múltiples responsabilidades por función
4. Código duplicado entre funciones GET

**Propuesta de Refactoring:**

```
VictimCaseFacade.go (695 líneas)
│
├─▶ VictimCaseHandler.go (200 líneas)
│   ├─▶ HandleCreate()
│   ├─▶ HandleUpdate()
│   ├─▶ HandleGet()
│   └─▶ HandleDelete()
│
├─▶ VictimCaseRenderer.go (200 líneas)
│   ├─▶ RenderCaseList()
│   ├─▶ RenderCaseDetail()
│   ├─▶ RenderCaseForm()
│   └─▶ RenderCaseReport()
│
└─▶ VictimCasePermissionChecker.go (100 líneas)
    ├─▶ CheckCreatePermission()
    ├─▶ CheckUpdatePermission()
    ├─▶ CheckViewPermission()
    └─▶ CheckDeletePermission()
```

**Esfuerzo Estimado:** 20-30 horas


---

## 3. DEUDA TÉCNICA

### 3.1 Patrones Problemáticos

#### 3.1.1 Código Duplicado

**HALLAZGO #1: Lógica de validación de sesión duplicada**
- **Archivos:** Todos los facades (10+ archivos)
- **Líneas:** ~5-10 líneas por archivo
- **Impacto:** MEDIO
- **Esfuerzo:** 4 horas

**Código Actual:**
```go
// Repetido en cada facade
session := sessions.Default(c)
var sessionID string = session.Get("userData").(string)
s, _ := utils.GetCommonSession(sessionID)
```

**Propuesta de Fix:**
```go
// common/utils/SessionHelper.go
func GetCurrentSession(c *gin.Context) (*CommonSession, error) {
    session := sessions.Default(c)
    sessionID, ok := session.Get("userData").(string)
    if !ok {
        return nil, errors.New("invalid session")
    }
    return GetCommonSession(sessionID)
}

// Uso en facades
s, err := utils.GetCurrentSession(c)
if err != nil {
    c.JSON(401, gin.H{"error": "Unauthorized"})
    return
}
```

---

**HALLAZGO #2: Lógica de renderizado de templates duplicada**
- **Archivos:** Todos los facades con renderizado HTML
- **Líneas:** ~20-30 líneas por archivo
- **Impacto:** MEDIO
- **Esfuerzo:** 6 horas

**Código Actual:**
```go
// Repetido en múltiples facades
common_facades.RenderTemplate(c, entityName, "salvia", "folder/", 
    config.HTML_Templates, "template_name", 
    utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, 
    utils.DEFAULT_PANIC_TEMPLATE,
    map[string]interface{}{
        "windowTitle": "...",
        "currentUser": s.Names + " " + s.LastNames,
        "menu": menu,
        // ... 20+ campos más
    }, utils.GetFullHtmlFuncMap())
```

**Propuesta de Fix:**
```go
// common/utils/TemplateRenderer.go
type TemplateContext struct {
    Session      *CommonSession
    WindowTitle  string
    Data         interface{}
    ExtraFields  map[string]interface{}
}

func RenderWithContext(c *gin.Context, ctx *TemplateContext, 
                       templateName string) {
    fields := map[string]interface{}{
        "windowTitle": ctx.WindowTitle,
        "currentUser": ctx.Session.Names + " " + ctx.Session.LastNames,
        "menu":        ctx.Session.CurrentMenu,
        "lang":        ctx.Session.Lang,
        "locale":      getLocale(ctx.Session.Lang),
    }
    // Merge con campos extra
    for k, v := range ctx.ExtraFields {
        fields[k] = v
    }
    // Renderizar
}
```

---

**HALLAZGO #3: Lógica de diff de slices duplicada**
- **Archivos:** 
  - `security/controllers/GeneralUserController.go` (líneas 1495-1572)
  - `salvia/controllers/VictimCaseController.go` (líneas 1755-1779)
- **Impacto:** MEDIO
- **Esfuerzo:** 3 horas

**Código Duplicado:**
```go
// Repetido para emails, phones, roles, moments
func cleanEmails(mailsToDelete []EMailDTO, mails []EMailDTO) 
    ([]EMailDTO, []EMailDTO) {
    // Lógica de diff
}
func cleanPhones(phonesToDelete []PhoneNumberDTO, phones []PhoneNumberDTO) 
    ([]PhoneNumberDTO, []PhoneNumberDTO) {
    // Misma lógica
}
func cleanRoles(rolesToDelete []RoleDTO, roles []RoleDTO) 
    ([]RoleDTO, []RoleDTO) {
    // Misma lógica
}
```

**Propuesta de Fix:**
```go
// common/utils/SliceHelper.go
type Identifiable interface {
    GetID() uint64
}

func DiffSlices[T Identifiable](old, new []T) (toDelete, toInsert []T) {
    oldMap := make(map[uint64]T)
    for _, item := range old {
        oldMap[item.GetID()] = item
    }
    
    for _, item := range new {
        if _, exists := oldMap[item.GetID()]; !exists {
            toInsert = append(toInsert, item)
        } else {
            delete(oldMap, item.GetID())
        }
    }
    
    for _, item := range oldMap {
        toDelete = append(toDelete, item)
    }
    return
}
```

---

#### 3.1.2 Funciones Largas (>50 líneas)

| Archivo | Función | Líneas | Impacto | Esfuerzo |
|---------|---------|--------|---------|----------|
| `salvia/controllers/VictimCaseController.go` | `SetVictimCase()` | 405 | CRÍTICO | 16h |
| `salvia/controllers/VictimCaseController.go` | `UpdateVictimCase()` | 397 | CRÍTICO | 14h |
| `security/controllers/GeneralUserController.go` | `UpdateGeneralUserByICode()` | 393 | CRÍTICO | 14h |
| `security/controllers/GeneralUserController.go` | `SetGeneralUser()` | 389 | CRÍTICO | 14h |
| `salvia/facades/VictimCaseGET()` | `VictimCaseGET()` | 220 | ALTO | 8h |
| `salvia/facades/VictimCaseFacade.go` | `VictimCasePUT_GET()` | 156 | ALTO | 6h |
| `security/facades/GeneralUserFacade.go` | `GeneralUserLOGIN_POST()` | 148 | ALTO | 6h |

**Total de funciones >50 líneas:** 25+  
**Esfuerzo total estimado:** 120-150 horas

---

#### 3.1.3 Magic Numbers y Strings Hardcodeados

**HALLAZGO #4: Límite de concurrencia hardcodeado**
- **Archivo:** `common/facades/MainRouter.go`
- **Línea:** 23
- **Impacto:** BAJO
- **Esfuerzo:** 1 hora

**Código Actual:**
```go
var sem = make(chan struct{}, 9999) // Limitar a 9999 peticiones concurrentes
```

**Fix:**
```go
// common/config/ServerConfig.go
const (
    MAX_CONCURRENT_REQUESTS = 9999
    MAX_MULTIPART_MEMORY    = 50 << 20 // 50 MB
    SESSION_MAX_AGE         = 14400    // 4 horas
)

// MainRouter.go
var sem = make(chan struct{}, common_config.MAX_CONCURRENT_REQUESTS)
```

---

**HALLAZGO #5: Tamaños de validación hardcodeados**
- **Archivos:** Múltiples DAOs
- **Impacto:** MEDIO
- **Esfuerzo:** 8 horas

**Código Actual:**
```go
// salvia/dao/VictimContactDAO.go
VictimContactFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
    "VictimContactNames": {
        Name: "VictimContactNames", 
        MinSize: 3,  // Magic number
        MaxSize: 32, // Magic number
        Required: true
    },
}
```

**Fix:**
```go
// common/config/ValidationRules.go
const (
    MIN_NAME_LENGTH     = 3
    MAX_NAME_LENGTH     = 32
    MIN_PHONE_LENGTH    = 6
    MAX_PHONE_LENGTH    = 10
    MIN_EMAIL_LENGTH    = 5
    MAX_EMAIL_LENGTH    = 128
    MIN_PASSWORD_LENGTH = 8
    MAX_PASSWORD_LENGTH = 32
)
```

---

**HALLAZGO #6: Strings de status hardcodeados**
- **Archivos:** Múltiples controllers
- **Impacto:** MEDIO
- **Esfuerzo:** 4 horas

**Código Actual:**
```go
// Disperso en múltiples archivos
vCase.VictimCaseStatus = "pe" // pending
vCase.VictimCaseStatus = "ap" // approved
usr.GeneralUserStatus = "ac"  // active
usr.GeneralUserStatus = "in"  // inactive
```

**Fix:**
```go
// common/config/StatusCodes.go
const (
    STATUS_PENDING  = "pe"
    STATUS_APPROVED = "ap"
    STATUS_ACTIVE   = "ac"
    STATUS_INACTIVE = "in"
)
```

---

#### 3.1.4 Comentarios TODO, FIXME, HACK

**HALLAZGO #7: Comentarios TODO sin resolver**
- **Archivo:** `salvia/facades/VictimCaseFacade.go`
- **Líneas:** 450-455
- **Impacto:** BAJO
- **Esfuerzo:** 2 horas

**Código:**
```go
/*
for _, m := range victimCase.VictimCaseMoments {
    if momentsWithSector[m.MomentCode] == nil {
        momentsWithSector[m.MomentCode] = make(map[uint64]salvia_daos.MomentDTO)
    }
    momentsWithSector[m.MomentCode][m.MomentEntityBranch.EntityBranchId] = m
}*/
// TODO: Descomentar cuando se implemente la lógica de momentos
```

---

**HALLAZGO #8: Comentario sobre selección de rol**
- **Archivo:** `security/facades/GeneralUserFacade.go`
- **Línea:** 91
- **Impacto:** MEDIO
- **Esfuerzo:** 8 horas

**Código:**
```go
// Se toma el primer rol asignado al usuario; en el futuro, se permitirá elegir el rol.
var currentRole string
if len(usr.GeneralUserRoleCodes) > 0 {
    currentRole = usr.GeneralUserRoleCodes[0]
}
```

**Propuesta:** Implementar selector de rol en el login

---

#### 3.1.5 Componentes Vue.js con >300 líneas

**HALLAZGO #9: Archivo home.vue**
- **Archivo:** `frontend/home.vue`
- **Líneas:** Desconocido (no analizado en detalle)
- **Impacto:** MEDIO
- **Esfuerzo:** 6 horas

**Propuesta:** Dividir en componentes más pequeños

---

### 3.2 Antipatrones Detectados

#### 3.2.1 God Structs/Packages

**HALLAZGO #10: VictimCaseDTO con 80+ campos**
- **Archivo:** `salvia/dao/VictimCaseDAO.go`
- **Líneas:** 1-200
- **Impacto:** ALTO
- **Esfuerzo:** 20 horas

**Problema:**
```go
type VictimCaseDTO struct {
    // 80+ campos
    VictimCaseId uint64
    VictimCaseICode string
    VictimCaseCreationDate time.Time
    // ... 77 campos más
}
```

**Propuesta de Fix:**
```go
// Dividir en múltiples DTOs
type VictimCaseDTO struct {
    ID           uint64
    ICode        string
    CreationDate time.Time
    Status       string
    Victim       VictimInfoDTO
    Facts        FactsDTO
    Aggressor    AggressorDTO
    Moments      []MomentDTO
}

type VictimInfoDTO struct {
    Names           string
    LastNames       string
    DocType         string
    DocNumber       string
    BirthDate       time.Time
    Contact         ContactInfoDTO
    Demographics    DemographicsDTO
}

type FactsDTO struct {
    Occurrence      string
    Date            time.Time
    StartTime       time.Time
    EndTime         time.Time
    Description     string
    Location        LocationDTO
}
```


---

**HALLAZGO #11: VictimCaseController como God Object**
- **Archivo:** `salvia/controllers/VictimCaseController.go`
- **Líneas:** 1779
- **Impacto:** CRÍTICO
- **Esfuerzo:** 40 horas

**Problema:** Un solo archivo con 26 funciones y múltiples responsabilidades

**Propuesta:** Ver sección 2.2 (Refactoring propuesto)

---

#### 3.2.2 Acoplamiento Fuerte

**HALLAZGO #12: salvia/controllers depende de security/controllers**
- **Archivos:** 
  - `salvia/controllers/VictimCaseController.go` (importa security_ctrl)
- **Línea:** 15
- **Impacto:** ALTO
- **Esfuerzo:** 12 horas

**Código Problemático:**
```go
import (
    security_ctrl "bitsflow/security/controllers"
)

// Línea 150
code, res := security_ctrl.SetGeneralUser(...)
```

**Problema:** Módulo salvia no debería conocer implementación de security

**Propuesta de Fix:**
```go
// common/interfaces/UserService.go
type UserService interface {
    CreateUser(userData UserCreateDTO) (UserDTO, error)
    UpdateUser(userID string, userData UserUpdateDTO) error
    GetUser(userID string) (UserDTO, error)
}

// security/services/UserServiceImpl.go
type UserServiceImpl struct {
    // implementación
}

// salvia/controllers/VictimCaseController.go
type VictimCaseController struct {
    userService common_interfaces.UserService
}
```

---

**HALLAZGO #13: Uso de recover() genérico**
- **Archivo:** No encontrado en análisis inicial
- **Impacto:** MEDIO
- **Esfuerzo:** N/A

**Nota:** No se encontraron usos de recover() en el código analizado, lo cual es positivo.

---

#### 3.2.3 Logging Inadecuado

**HALLAZGO #14: fmt.Println() en código de producción**
- **Archivo:** `security/facades/GeneralUserFacade.go`
- **Línea:** 521
- **Impacto:** MEDIO
- **Esfuerzo:** 8 horas

**Código Problemático:**
```go
id := c.Param("id")
fmt.Println(id) // Logging a stdout
```

**Propuesta de Fix:**
```go
// Implementar logging estructurado
import "github.com/sirupsen/logrus"

var log = logrus.New()

// En el código
log.WithFields(logrus.Fields{
    "user_id": id,
    "action":  "delete_user",
}).Info("Deleting user")
```

---

**HALLAZGO #15: Ausencia de logging de errores**
- **Archivos:** Múltiples
- **Impacto:** ALTO
- **Esfuerzo:** 16 horas

**Código Problemático:**
```go
s, _ := utils.GetCommonSession(sessionID) // Error ignorado
```

**Propuesta de Fix:**
```go
s, err := utils.GetCommonSession(sessionID)
if err != nil {
    log.WithError(err).Error("Failed to get session")
    c.JSON(401, gin.H{"error": "Unauthorized"})
    return
}
```

---

#### 3.2.4 Manejo de Errores Ignorado

**HALLAZGO #16: Errores ignorados con _ = err**
- **Archivos:** Múltiples facades y controllers
- **Impacto:** ALTO
- **Esfuerzo:** 20 horas

**Ejemplos:**
```go
// common/facades/MainRouter.go, línea 82
err := router.RunTLS(":443", "certs/salviaTest.crt", "certs/salviaTest.key")
println(err) // Solo imprime, no maneja

// Múltiples facades
s, _ := utils.GetCommonSession(sessionID) // Error ignorado
```

**Propuesta de Fix:**
```go
// MainRouter.go
if err := router.RunTLS(":443", "certs/salviaTest.crt", "certs/salviaTest.key"); err != nil {
    log.WithError(err).Fatal("Failed to start HTTPS server")
}

// Facades
s, err := utils.GetCommonSession(sessionID)
if err != nil {
    log.WithError(err).Warn("Invalid session")
    c.Redirect(302, "/static/landing.html")
    return
}
```

---

#### 3.2.5 Frontend: Mutación Directa de Props

**HALLAZGO #17: Uso de XMLHttpRequest en lugar de Fetch API**
- **Archivo:** `frontend/js/dao.js`
- **Líneas:** 1-300
- **Impacto:** BAJO
- **Esfuerzo:** 4 horas

**Código Actual:**
```javascript
var getData = function (url, callback, hideDialogs) {
  const xhr = new XMLHttpRequest()
  xhr.open("GET", url);
  xhr.send();
  // ...
}
```

**Propuesta de Fix:**
```javascript
// Usar Fetch API moderna
async function getData(url, hideDialogs = false) {
  try {
    if (!hideDialogs) showLoading();
    
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      credentials: 'same-origin'
    });
    
    if (!hideDialogs) hideLoading();
    
    if (response.redirected && response.url.includes('login')) {
      alert('Sesión expirada');
      window.location.href = '/';
      return;
    }
    
    if (response.ok) {
      if (!hideDialogs) showSuccess();
      return await response.json();
    } else {
      if (!hideDialogs) showError();
      throw new Error(`HTTP ${response.status}`);
    }
  } catch (error) {
    console.error('Network error:', error);
    throw error;
  }
}
```

---

### 3.3 Cobertura y Testing

#### 3.3.1 Ausencia de Tests

**HALLAZGO #18: Sin archivos *_test.go**
- **Impacto:** CRÍTICO
- **Esfuerzo:** 200+ horas

**Paquetes sin tests:**
- ❌ `common/controllers` (0 tests)
- ❌ `common/utils` (0 tests)
- ❌ `common/db` (0 tests)
- ❌ `salvia/controllers` (0 tests)
- ❌ `salvia/facades` (0 tests)
- ❌ `salvia/dao` (0 tests)
- ❌ `security/controllers` (0 tests)
- ❌ `security/facades` (0 tests)
- ❌ `security/dao` (0 tests)

**Propuesta:**
```
Prioridad 1 (40h):
├─▶ common/utils/EntityValidator_test.go
├─▶ common/utils/CommonSession_test.go
├─▶ security/controllers/GeneralUserController_test.go
└─▶ salvia/controllers/VictimCaseController_test.go

Prioridad 2 (60h):
├─▶ common/db/PostgresConnection_test.go
├─▶ security/dao/GeneralUserDAO_test.go
├─▶ salvia/dao/VictimCaseDAO_test.go
└─▶ salvia/dao/VictimContactDAO_test.go

Prioridad 3 (100h):
└─▶ Tests de integración para todos los facades
```

---

**HALLAZGO #19: Frontend sin tests**
- **Impacto:** ALTO
- **Esfuerzo:** 80 horas

**Archivos sin tests:**
- ❌ `frontend/js/dao.js`
- ❌ `frontend/js/formulario.js`
- ❌ `frontend/js/login.js`
- ❌ `frontend/home.vue`

**Propuesta:**
```javascript
// frontend/tests/dao.test.js (usando Jest)
describe('getData', () => {
  it('should fetch data successfully', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ data: 'test' })
      })
    );
    
    const result = await getData('/api/test');
    expect(result).toEqual({ data: 'test' });
  });
  
  it('should handle session expiration', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        redirected: true,
        url: 'https://example.com/login'
      })
    );
    
    await getData('/api/test');
    expect(window.location.href).toBe('/');
  });
});
```

---

### 3.4 Dependencias

#### 3.4.1 Dependencias Go Desactualizadas

**HALLAZGO #20: Go 1.19 (EOL)**
- **Archivo:** `go.mod`
- **Línea:** 3
- **Impacto:** ALTO
- **Esfuerzo:** 8 horas

**Código Actual:**
```go
go 1.19
```

**Problema:** Go 1.19 llegó a End of Life en agosto 2023

**Fix:**
```go
go 1.22 // o 1.23 (última versión estable)
```

**Pasos:**
1. Actualizar go.mod
2. Ejecutar `go mod tidy`
3. Probar compilación
4. Ejecutar tests (cuando existan)
5. Probar en staging

---

**HALLAZGO #21: Dependencias sin mantenimiento activo**
- **Archivo:** `go.mod`
- **Impacto:** MEDIO
- **Esfuerzo:** 4 horas

**Dependencias a revisar:**
```go
github.com/gin-gonic/contrib v0.0.0-20201101042839-6a891bf89f19
// Última actualización: 2020 (4 años sin updates)
```

**Propuesta:** Migrar a alternativas mantenidas o fork interno

---

#### 3.4.2 Dependencias Frontend

**HALLAZGO #22: jQuery innecesario**
- **Archivos:** `frontend/plugins/jquery/`
- **Impacto:** BAJO
- **Esfuerzo:** 12 horas

**Problema:** jQuery es innecesario con Vue.js 3

**Propuesta:** Eliminar jQuery y migrar plugins a Vue.js

---

**HALLAZGO #23: Múltiples librerías de UI**
- **Archivos:** `frontend/plugins/`
- **Impacto:** MEDIO
- **Esfuerzo:** 40 horas

**Problema:** Bootstrap + múltiples plugins aumentan el bundle size

**Propuesta:** Consolidar en un solo framework UI (Vuetify, Element Plus, etc.)

---

### 3.5 Resumen de Deuda Técnica

| Categoría | Hallazgos | Impacto | Esfuerzo Total |
|-----------|-----------|---------|----------------|
| Código Duplicado | 3 | MEDIO | 13h |
| Funciones Largas | 25+ | CRÍTICO | 120-150h |
| Magic Numbers | 3 | MEDIO | 13h |
| God Objects | 2 | CRÍTICO | 60h |
| Acoplamiento Fuerte | 1 | ALTO | 12h |
| Logging Inadecuado | 2 | ALTO | 24h |
| Errores Ignorados | 1 | ALTO | 20h |
| Sin Tests | 2 | CRÍTICO | 280h |
| Dependencias | 3 | ALTO | 24h |
| **TOTAL** | **42+** | - | **566-596h** |


---

## 4. VULNERABILIDADES DE SEGURIDAD

### 4.1 Inyecciones

#### 4.1.1 SQL Injection

**✅ BUENA PRÁCTICA: Uso de Prepared Statements**

El código utiliza correctamente prepared statements con placeholders en la mayoría de queries:

```go
// Ejemplo correcto en VictimCaseDAO.go
query := `INSERT INTO salvia.victim_case (
    victim_case_names,
    victim_case_phone
) VALUES ($1, $2)`

err := connData.Conn.QueryRow(ctx, query, 
    dto.VictimCaseNames,
    dto.VictimCasePhone
).Scan(&dto.VictimCaseId)
```

**HALLAZGO #24: Posible SQL Injection en queries dinámicas**
- **Archivo:** No encontrado en análisis inicial
- **Severidad:** N/A
- **Estado:** ✅ No se encontraron vulnerabilidades de SQL Injection

---

#### 4.1.2 Command Injection

**HALLAZGO #25: Sin uso de exec.Command**
- **Severidad:** N/A
- **Estado:** ✅ No se encontró uso de exec.Command con input del usuario

---

#### 4.1.3 Template Injection

**HALLAZGO #26: Uso seguro de html/template**
- **Archivo:** `common/facades/MainRouter.go`
- **Severidad:** ✅ SEGURO

**Código:**
```go
import "html/template"

t := template.New(configTemplate[configTemplateName]).Funcs(funcMap)
t, err := t.ParseFiles(...)
err = t.Execute(&b, templateFields)
```

**Análisis:** Go's html/template escapa automáticamente el contenido, previniendo XSS.

---

### 4.2 Autenticación y Autorización

#### 4.2.1 Endpoints sin Autenticación

**HALLAZGO #27: Endpoint público sin rate limiting**
- **Archivo:** `salvia/facades/VictimContactFacade.go`
- **Línea:** 82
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 4 horas

**Código Vulnerable:**
```go
func VictimContactPOST_Public(c *gin.Context) {
    // Sin autenticación, sin rate limiting
    buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    code, res := salvia_ctrl.SetVictimContact(...)
    c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, ...)
}
```

**Riesgo:** 
- Abuso del endpoint (spam)
- Creación masiva de contactos falsos
- DoS por agotamiento de recursos

**Prueba de Concepto:**
```bash
# Script para crear 1000 contactos falsos
for i in {1..1000}; do
  curl -X POST https://salvia.minigualdadyequidad.gov.co/public/salvia/contacto-victima \
    -H "Content-Type: application/json" \
    -d '{
      "names": "Fake'$i'",
      "phone": "300000'$i'",
      "email": "fake'$i'@test.com"
    }'
done
```

**Fix Recomendado:**
```go
// Implementar rate limiting con redis o memoria
import "github.com/ulule/limiter/v3"
import "github.com/ulule/limiter/v3/drivers/store/memory"

var rateLimiter = limiter.New(memory.NewStore(), limiter.Rate{
    Period: 1 * time.Hour,
    Limit:  5, // 5 contactos por hora por IP
})

func VictimContactPOST_Public(c *gin.Context) {
    // Rate limiting por IP
    context, err := rateLimiter.Get(c, c.ClientIP())
    if err != nil {
        c.JSON(500, gin.H{"error": "Internal error"})
        return
    }
    
    if context.Reached {
        c.JSON(429, gin.H{"error": "Too many requests"})
        return
    }
    
    // Continuar con la lógica normal
    buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    code, res := salvia_ctrl.SetVictimContact(...)
    c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, ...)
}
```

---

#### 4.2.2 Credenciales Hardcodeadas

**HALLAZGO #28: Secret key hardcodeado en sesiones**
- **Archivo:** `common/facades/MainRouter.go`
- **Línea:** 39
- **Severidad:** 🔴 CRÍTICO
- **Esfuerzo:** 2 horas

**Código Vulnerable:**
```go
store := cookie.NewStore([]byte("secret"))
```

**Riesgo:**
- Cualquiera con acceso al código puede firmar cookies válidas
- Posible session hijacking
- Imposible rotar la clave sin cambiar código

**Prueba de Concepto:**
```go
// Atacante puede crear sesiones válidas
import "github.com/gorilla/sessions"

store := sessions.NewCookieStore([]byte("secret"))
session := sessions.NewSession(store, "default_session")
session.Values["userData"] = "admin-session-id"
// Generar cookie firmada válida
```

**Fix Recomendado:**
```go
// config/secrets.go
import "os"

func GetSessionSecret() []byte {
    secret := os.Getenv("SESSION_SECRET")
    if secret == "" {
        log.Fatal("SESSION_SECRET environment variable not set")
    }
    if len(secret) < 32 {
        log.Fatal("SESSION_SECRET must be at least 32 characters")
    }
    return []byte(secret)
}

// MainRouter.go
store := cookie.NewStore(config.GetSessionSecret())
```

**Configuración:**
```bash
# .env
SESSION_SECRET=tu-clave-secreta-muy-larga-y-aleatoria-de-al-menos-32-caracteres
```

---

**HALLAZGO #29: Sin validación de expiración de sesión**
- **Archivo:** `common/utils/CommonSession.go`
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 6 horas

**Código Actual:**
```go
// Sesiones en memoria sin timestamp de expiración
var commonSessions map[string]CommonSession = map[string]CommonSession{}

func AddCommonSession(sessionID string, cs CommonSession) {
    commonSessions[sessionID] = cs
}
```

**Riesgo:**
- Sesiones nunca expiran en memoria
- Posible memory leak
- Sesiones robadas válidas indefinidamente

**Fix Recomendado:**
```go
type CommonSession struct {
    UserICode        string
    Roles            []string
    CurrentRole      string
    // ... otros campos
    CreatedAt        time.Time
    LastAccessedAt   time.Time
    ExpiresAt        time.Time
}

func AddCommonSession(sessionID string, cs CommonSession) {
    cs.CreatedAt = time.Now()
    cs.LastAccessedAt = time.Now()
    cs.ExpiresAt = time.Now().Add(4 * time.Hour) // 4 horas
    commonSessions[sessionID] = cs
}

func GetCommonSession(sessionID string) (CommonSession, error) {
    cs, exists := commonSessions[sessionID]
    if !exists {
        return CommonSession{}, errors.New("session not found")
    }
    
    // Verificar expiración
    if time.Now().After(cs.ExpiresAt) {
        delete(commonSessions, sessionID)
        return CommonSession{}, errors.New("session expired")
    }
    
    // Actualizar último acceso
    cs.LastAccessedAt = time.Now()
    commonSessions[sessionID] = cs
    
    return cs, nil
}

// Goroutine para limpiar sesiones expiradas
func CleanExpiredSessions() {
    ticker := time.NewTicker(10 * time.Minute)
    for range ticker.C {
        now := time.Now()
        for id, session := range commonSessions {
            if now.After(session.ExpiresAt) {
                delete(commonSessions, id)
                log.WithField("session_id", id).Info("Expired session cleaned")
            }
        }
    }
}
```

---

#### 4.2.3 Lógica de Permisos Bypasseable

**HALLAZGO #30: Verificación de permisos inconsistente**
- **Archivos:** Múltiples facades
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 8 horas

**Código Vulnerable:**
```go
// Algunos facades verifican permisos
if !utils.CheckPermission(salvia_config.PermissionsByRole, 
    "set_victim_case", s.CurrentRole, c) {
    return
}

// Otros no verifican
func SomeHandler(c *gin.Context) {
    // Sin verificación de permisos
    buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    code, res := controller.DoSomething(...)
}
```

**Riesgo:** Escalación de privilegios si se olvida verificar permisos

**Fix Recomendado:**
```go
// Middleware de permisos
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        s, err := utils.GetCurrentSession(c)
        if err != nil {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        
        if !utils.CheckPermission(
            salvia_config.PermissionsByRole,
            permission,
            s.CurrentRole,
            c,
        ) {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Uso en routers
router.POST("/salvia/caso-victima", 
    RequirePermission("set_victim_case"),
    VictimCasePOST)
```

---

#### 4.2.4 Manejo Inseguro de JWT

**HALLAZGO #31: Sin uso de JWT**
- **Severidad:** N/A
- **Estado:** ✅ El sistema usa sesiones basadas en cookies, no JWT

---

### 4.3 Exposición de Datos

#### 4.3.1 Structs con Campos Sensibles Expuestos

**HALLAZGO #32: Password expuesto en JSON**
- **Archivo:** `security/dao/GeneralUserDAO.go`
- **Línea:** ~50
- **Severidad:** 🔴 CRÍTICO
- **Esfuerzo:** 2 horas

**Código Vulnerable:**
```go
type GeneralUserDTO struct {
    GeneralUserId       uint64 `json:"id"`
    GeneralUserLogin    string `json:"login"`
    GeneralUserPassword string `json:"password"` // ⚠️ EXPUESTO
    GeneralUserStatus   string `json:"status"`
    // ...
}
```

**Riesgo:** Password hash expuesto en respuestas JSON

**Fix Recomendado:**
```go
type GeneralUserDTO struct {
    GeneralUserId       uint64 `json:"id"`
    GeneralUserLogin    string `json:"login"`
    GeneralUserPassword string `json:"-"` // Omitir en JSON
    GeneralUserStatus   string `json:"status"`
    // ...
}
```

---

**HALLAZGO #33: Datos sensibles en logs**
- **Archivo:** `security/facades/GeneralUserFacade.go`
- **Línea:** 521
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 4 horas

**Código Vulnerable:**
```go
fmt.Println(id) // Puede imprimir IDs de usuario
```

**Fix:** Ver HALLAZGO #14 (implementar logging estructurado sin datos sensibles)

---


#### 4.3.2 Respuestas de Error Reveladoras

**HALLAZGO #34: Stack traces en respuestas de error**
- **Archivo:** `common/facades/MainRouter.go`
- **Líneas:** 120-140
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 4 horas

**Código Vulnerable:**
```go
t, err := t.ParseFiles(...)
if err != nil {
    utils.SetError(errorMap, entity, "default", 
        common_config.Enums.GLOBAL_ERROR, "", common_config.Locale)
    res = utils.CommMsgGetJSONErrors(errorMap)
    code = 400
}
```

**Riesgo:** Errores genéricos pueden revelar información interna

**Fix Recomendado:**
```go
// Logging interno detallado
log.WithError(err).WithFields(logrus.Fields{
    "entity":   entity,
    "template": configTemplate[configTemplateName],
}).Error("Template parsing failed")

// Respuesta genérica al cliente
if err != nil {
    c.JSON(500, gin.H{
        "error": "Internal server error",
        "code":  "TEMPLATE_ERROR"
    })
    return
}
```

---

#### 4.3.3 Variables de Entorno sin Valor por Defecto Seguro

**HALLAZGO #35: Sin uso de variables de entorno**
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 8 horas

**Problema:** Configuración hardcodeada en código

**Archivos Afectados:**
- `common/facades/MainRouter.go` (puerto 443 hardcodeado)
- `config/db_config.json` (credenciales en archivo)

**Fix Recomendado:**
```go
// config/env.go
import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort      int
    DBHost          string
    DBPort          int
    DBName          string
    DBUser          string
    DBPassword      string
    SessionSecret   string
    TLSCertPath     string
    TLSKeyPath      string
}

func LoadConfig() (*Config, error) {
    port, _ := strconv.Atoi(getEnv("SERVER_PORT", "443"))
    dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
    
    return &Config{
        ServerPort:    port,
        DBHost:        getEnv("DB_HOST", "localhost"),
        DBPort:        dbPort,
        DBName:        getEnv("DB_NAME", "salvia"),
        DBUser:        getEnv("DB_USER", ""),
        DBPassword:    getEnv("DB_PASSWORD", ""),
        SessionSecret: getEnv("SESSION_SECRET", ""),
        TLSCertPath:   getEnv("TLS_CERT_PATH", "certs/salviaTest.crt"),
        TLSKeyPath:    getEnv("TLS_KEY_PATH", "certs/salviaTest.key"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

---

### 4.4 Validación de Entrada

#### 4.4.1 Validación de Tipo, Longitud y Formato

**✅ BUENA PRÁCTICA: Validación exhaustiva en DAOs**

```go
// VictimContactDAO.go
VictimContactFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
    "VictimContactNames": {
        Name: "VictimContactNames",
        ModelType: "string",
        MinSize: 3,
        MaxSize: 32,
        Required: true
    },
    "VictimContactPhone": {
        Name: "VictimContactPhone",
        ModelType: "string",
        MinSize: 6,
        MaxSize: 10,
        Required: true
    },
}
```

**HALLAZGO #36: Validación de email insuficiente**
- **Archivo:** `common/utils/EntityValidator.go` (asumido)
- **Severidad:** 🟡 BAJO
- **Esfuerzo:** 2 horas

**Propuesta:**
```go
import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
    if len(email) < 5 || len(email) > 128 {
        return false
    }
    return emailRegex.MatchString(email)
}
```

---

#### 4.4.2 Archivos Subidos sin Validación

**HALLAZGO #37: Validación de archivos subidos**
- **Archivo:** `salvia/facades/LoadFilesFacade.go`
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 6 horas

**Código a Revisar:**
```go
// MainRouter.go
router.MaxMultipartMemory = 50 << 20 // 50 MB
```

**Propuesta de Fix:**
```go
func ValidateUploadedFile(file *multipart.FileHeader) error {
    // Validar tamaño
    if file.Size > 10*1024*1024 { // 10 MB
        return errors.New("file too large")
    }
    
    // Validar extensión
    allowedExtensions := []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"}
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if !contains(allowedExtensions, ext) {
        return errors.New("invalid file type")
    }
    
    // Validar MIME type
    f, err := file.Open()
    if err != nil {
        return err
    }
    defer f.Close()
    
    buffer := make([]byte, 512)
    _, err = f.Read(buffer)
    if err != nil {
        return err
    }
    
    mimeType := http.DetectContentType(buffer)
    allowedMimes := []string{
        "application/pdf",
        "image/jpeg",
        "image/png",
        "application/msword",
    }
    
    if !contains(allowedMimes, mimeType) {
        return errors.New("invalid MIME type")
    }
    
    return nil
}
```

---

#### 4.4.3 URLs Externas sin Validación (SSRF)

**HALLAZGO #38: Sin llamadas a URLs externas**
- **Severidad:** N/A
- **Estado:** ✅ No se encontraron llamadas HTTP a URLs externas con input del usuario

---

#### 4.4.4 Rate Limiting Ausente

**HALLAZGO #39: Sin rate limiting en login**
- **Archivo:** `security/facades/GeneralUserFacade.go`
- **Línea:** 55
- **Severidad:** 🔴 CRÍTICO
- **Esfuerzo:** 6 horas

**Código Vulnerable:**
```go
func GeneralUserLOGIN_POST(c *gin.Context) {
    // Sin rate limiting
    buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    code, res, usr := security_ctrl.LoginGeneralUser(...)
}
```

**Riesgo:** Ataques de fuerza bruta contra credenciales

**Prueba de Concepto:**
```bash
# Ataque de fuerza bruta
for password in $(cat passwords.txt); do
  curl -X POST https://salvia.minigualdadyequidad.gov.co/seguridad/login \
    -H "Content-Type: application/json" \
    -d "{\"login\":\"admin\",\"password\":\"$password\",\"captcha\":\"\"}"
done
```

**Fix Recomendado:**
```go
import "github.com/ulule/limiter/v3"

var loginLimiter = limiter.New(memory.NewStore(), limiter.Rate{
    Period: 15 * time.Minute,
    Limit:  5, // 5 intentos cada 15 minutos
})

func GeneralUserLOGIN_POST(c *gin.Context) {
    // Rate limiting por IP + login
    identifier := c.ClientIP() + ":" + getLoginFromRequest(c)
    context, err := loginLimiter.Get(c, identifier)
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Internal error"})
        return
    }
    
    if context.Reached {
        log.WithFields(logrus.Fields{
            "ip":    c.ClientIP(),
            "login": getLoginFromRequest(c),
        }).Warn("Login rate limit exceeded")
        
        c.JSON(429, gin.H{
            "error": "Too many login attempts. Try again in 15 minutes."
        })
        return
    }
    
    // Continuar con login normal
    buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    code, res, usr := security_ctrl.LoginGeneralUser(...)
    
    // Si login falla, no resetear el contador
    if code != 200 {
        log.WithFields(logrus.Fields{
            "ip":    c.ClientIP(),
            "login": getLoginFromRequest(c),
        }).Warn("Failed login attempt")
    }
    
    c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, ...)
}
```

---

### 4.5 Configuración Insegura

#### 4.5.1 Modo Debug Activo

**HALLAZGO #40: Gin en modo debug**
- **Archivo:** `main.go`
- **Línea:** 14 (comentado)
- **Severidad:** 🟡 BAJO
- **Esfuerzo:** 1 hora

**Código:**
```go
//gin.SetMode(gin.ReleaseMode)
```

**Fix:**
```go
// Configurar según entorno
if os.Getenv("ENV") == "production" {
    gin.SetMode(gin.ReleaseMode)
} else {
    gin.SetMode(gin.DebugMode)
}
```

---

#### 4.5.2 CORS Configurado Inseguramente

**HALLAZGO #41: Sin configuración de CORS**
- **Severidad:** 🟡 BAJO
- **Esfuerzo:** 2 horas

**Estado:** No se encontró configuración de CORS

**Propuesta:**
```go
import "github.com/gin-contrib/cors"

func InitRouter() *gin.Engine {
    router := gin.Default()
    
    // Configurar CORS
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://salvia.minigualdadyequidad.gov.co"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    return router
}
```

---

#### 4.5.3 Headers de Seguridad HTTP Ausentes

**HALLAZGO #42: Sin headers de seguridad**
- **Archivo:** `common/facades/MainRouter.go`
- **Severidad:** 🟠 MEDIO
- **Esfuerzo:** 3 horas

**Headers Faltantes:**
- Content-Security-Policy
- X-Frame-Options
- X-Content-Type-Options
- Strict-Transport-Security
- X-XSS-Protection

**Fix Recomendado:**
```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevenir clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // Prevenir MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // XSS Protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // HSTS (solo HTTPS)
        c.Header("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains")
        
        // Content Security Policy
        c.Header("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "font-src 'self' data:; "+
            "connect-src 'self'")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        c.Next()
    }
}

// En InitRouter()
router.Use(SecurityHeadersMiddleware())
```

---

#### 4.5.4 Timeouts Ausentes en Clientes HTTP

**HALLAZGO #43: Sin timeouts en XMLHttpRequest**
- **Archivo:** `frontend/js/dao.js`
- **Severidad:** 🟡 BAJO
- **Esfuerzo:** 2 horas

**Código Vulnerable:**
```javascript
const xhr = new XMLHttpRequest()
xhr.open("GET", url);
xhr.send();
// Sin timeout configurado
```

**Fix:**
```javascript
const xhr = new XMLHttpRequest()
xhr.timeout = 30000; // 30 segundos
xhr.ontimeout = function() {
    console.error("Request timed out");
    callback(408, null); // Request Timeout
};
xhr.open("GET", url);
xhr.send();
```

---

#### 4.5.5 Variables de Entorno Sensibles Expuestas

**HALLAZGO #44: Sin uso de variables de entorno en frontend**
- **Severidad:** N/A
- **Estado:** ✅ No se encontraron variables de entorno expuestas con prefijo VITE_ o VUE_APP_

---

### 4.6 Resumen de Vulnerabilidades

| ID | Vulnerabilidad | Archivo | Severidad | OWASP | Esfuerzo |
|----|----------------|---------|-----------|-------|----------|
| #27 | Endpoint público sin rate limiting | VictimContactFacade.go:82 | 🟠 MEDIO | A05:2021 | 4h |
| #28 | Secret key hardcodeado | MainRouter.go:39 | 🔴 CRÍTICO | A02:2021 | 2h |
| #29 | Sin validación de expiración de sesión | CommonSession.go | 🟠 MEDIO | A07:2021 | 6h |
| #30 | Verificación de permisos inconsistente | Múltiples facades | 🟠 MEDIO | A01:2021 | 8h |
| #32 | Password expuesto en JSON | GeneralUserDAO.go:50 | 🔴 CRÍTICO | A01:2021 | 2h |
| #33 | Datos sensibles en logs | GeneralUserFacade.go:521 | 🟠 MEDIO | A09:2021 | 4h |
| #34 | Stack traces en respuestas | MainRouter.go:120-140 | 🟠 MEDIO | A05:2021 | 4h |
| #35 | Sin variables de entorno | Múltiples | 🟠 MEDIO | A05:2021 | 8h |
| #36 | Validación de email insuficiente | EntityValidator.go | 🟡 BAJO | A03:2021 | 2h |
| #37 | Validación de archivos subidos | LoadFilesFacade.go | 🟠 MEDIO | A03:2021 | 6h |
| #39 | Sin rate limiting en login | GeneralUserFacade.go:55 | 🔴 CRÍTICO | A07:2021 | 6h |
| #40 | Gin en modo debug | main.go:14 | 🟡 BAJO | A05:2021 | 1h |
| #42 | Sin headers de seguridad | MainRouter.go | 🟠 MEDIO | A05:2021 | 3h |
| #43 | Sin timeouts en HTTP | dao.js | 🟡 BAJO | A05:2021 | 2h |

**Total de Vulnerabilidades:** 14  
**Esfuerzo Total:** 58 horas


---

## 5. RESUMEN EJECUTIVO

### 5.1 Hallazgos por Categoría y Severidad

| Categoría | 🔴 Crítico | 🟠 Alto | 🟠 Medio | 🟡 Bajo | Total |
|-----------|-----------|---------|----------|---------|-------|
| **Deuda Técnica** | 4 | 3 | 6 | 0 | 13 |
| **Vulnerabilidades de Seguridad** | 3 | 0 | 8 | 3 | 14 |
| **Arquitectura** | 0 | 2 | 3 | 0 | 5 |
| **Testing** | 2 | 0 | 0 | 0 | 2 |
| **TOTAL** | **9** | **5** | **17** | **3** | **34** |

### 5.2 Top 10 Issues Más Críticos

#### 🔴 #1: Secret Key Hardcodeado en Sesiones (HALLAZGO #28)
- **Severidad:** CRÍTICA
- **Archivo:** `common/facades/MainRouter.go:39`
- **Impacto:** Cualquiera puede firmar cookies válidas y secuestrar sesiones
- **Esfuerzo:** 2 horas
- **Prioridad:** INMEDIATA

**Acción Requerida:**
```bash
# 1. Generar secret aleatorio
openssl rand -base64 32

# 2. Configurar variable de entorno
export SESSION_SECRET="tu-secret-generado"

# 3. Modificar código (ver HALLAZGO #28)
```

---

#### 🔴 #2: Password Hash Expuesto en JSON (HALLAZGO #32)
- **Severidad:** CRÍTICA
- **Archivo:** `security/dao/GeneralUserDAO.go:50`
- **Impacto:** Hashes de contraseñas expuestos en respuestas API
- **Esfuerzo:** 2 horas
- **Prioridad:** INMEDIATA

**Acción Requerida:**
```go
// Agregar json:"-" al campo password
GeneralUserPassword string `json:"-"`
```

---

#### 🔴 #3: Sin Rate Limiting en Login (HALLAZGO #39)
- **Severidad:** CRÍTICA
- **Archivo:** `security/facades/GeneralUserFacade.go:55`
- **Impacto:** Vulnerable a ataques de fuerza bruta
- **Esfuerzo:** 6 horas
- **Prioridad:** ALTA

**Acción Requerida:** Implementar rate limiting (ver HALLAZGO #39)

---

#### 🔴 #4: VictimCaseController.go - 1779 Líneas (HALLAZGO #1)
- **Severidad:** CRÍTICA (Deuda Técnica)
- **Archivo:** `salvia/controllers/VictimCaseController.go`
- **Impacto:** Mantenibilidad extremadamente baja, alto riesgo de bugs
- **Esfuerzo:** 40-60 horas
- **Prioridad:** ALTA

**Acción Requerida:** Refactorizar en 6 archivos (ver sección 2.2)

---

#### 🔴 #5: GeneralUserController.go - 1582 Líneas (HALLAZGO #2)
- **Severidad:** CRÍTICA (Deuda Técnica)
- **Archivo:** `security/controllers/GeneralUserController.go`
- **Impacto:** Mantenibilidad extremadamente baja
- **Esfuerzo:** 35-50 horas
- **Prioridad:** ALTA

**Acción Requerida:** Refactorizar en 6 archivos (ver sección 2.2)

---

#### 🔴 #6: Sin Tests Unitarios (HALLAZGO #18)
- **Severidad:** CRÍTICA
- **Impacto:** Sin garantías de calidad, regresiones frecuentes
- **Esfuerzo:** 200+ horas
- **Prioridad:** ALTA

**Acción Requerida:** Implementar tests (ver sección 3.3.1)

---

#### 🔴 #7: Go 1.19 EOL (HALLAZGO #20)
- **Severidad:** CRÍTICA
- **Archivo:** `go.mod:3`
- **Impacto:** Sin parches de seguridad
- **Esfuerzo:** 8 horas
- **Prioridad:** ALTA

**Acción Requerida:**
```bash
# Actualizar a Go 1.22 o 1.23
go mod edit -go=1.22
go mod tidy
go build
```

---

#### 🟠 #8: Endpoint Público sin Rate Limiting (HALLAZGO #27)
- **Severidad:** MEDIA
- **Archivo:** `salvia/facades/VictimContactFacade.go:82`
- **Impacto:** Vulnerable a spam y DoS
- **Esfuerzo:** 4 horas
- **Prioridad:** MEDIA

---

#### 🟠 #9: Sin Validación de Expiración de Sesión (HALLAZGO #29)
- **Severidad:** MEDIA
- **Archivo:** `common/utils/CommonSession.go`
- **Impacto:** Sesiones válidas indefinidamente, memory leak
- **Esfuerzo:** 6 horas
- **Prioridad:** MEDIA

---

#### 🟠 #10: Sin Headers de Seguridad HTTP (HALLAZGO #42)
- **Severidad:** MEDIA
- **Archivo:** `common/facades/MainRouter.go`
- **Impacto:** Vulnerable a clickjacking, XSS, MIME sniffing
- **Esfuerzo:** 3 horas
- **Prioridad:** MEDIA

---

### 5.3 Estimación Total de Deuda Técnica

| Categoría | Hallazgos | Esfuerzo (horas) | Costo Estimado* |
|-----------|-----------|------------------|-----------------|
| **Refactoring Crítico** | 2 | 75-110 | $15,000 - $22,000 |
| **Implementación de Tests** | 2 | 280 | $56,000 |
| **Vulnerabilidades Críticas** | 3 | 10 | $2,000 |
| **Vulnerabilidades Medias** | 8 | 48 | $9,600 |
| **Deuda Técnica Media** | 13 | 200 | $40,000 |
| **Mejoras Menores** | 6 | 30 | $6,000 |
| **TOTAL** | **34** | **643-678** | **$128,600 - $135,600** |

*Costo estimado a $200 USD/hora (tarifa promedio senior developer)

---

### 5.4 Roadmap de Correcciones

#### 🚨 FASE 1: Quick Wins Críticos (2-3 semanas)

**Objetivo:** Resolver vulnerabilidades críticas de seguridad

| # | Tarea | Esfuerzo | Prioridad |
|---|-------|----------|-----------|
| 1 | Mover secret key a variable de entorno (#28) | 2h | 🔴 INMEDIATA |
| 2 | Ocultar password en JSON (#32) | 2h | 🔴 INMEDIATA |
| 3 | Implementar rate limiting en login (#39) | 6h | 🔴 ALTA |
| 4 | Implementar rate limiting en endpoint público (#27) | 4h | 🟠 ALTA |
| 5 | Agregar headers de seguridad HTTP (#42) | 3h | 🟠 ALTA |
| 6 | Implementar validación de expiración de sesión (#29) | 6h | 🟠 ALTA |
| 7 | Actualizar Go a 1.22+ (#20) | 8h | 🔴 ALTA |
| 8 | Configurar modo release en producción (#40) | 1h | 🟡 MEDIA |

**Total Fase 1:** 32 horas (~1 sprint)

---

#### 🔧 FASE 2: Mejoras de Arquitectura (2-3 meses)

**Objetivo:** Reducir complejidad y mejorar mantenibilidad

| # | Tarea | Esfuerzo | Prioridad |
|---|-------|----------|-----------|
| 1 | Refactorizar VictimCaseController.go | 40-60h | 🔴 ALTA |
| 2 | Refactorizar GeneralUserController.go | 35-50h | 🔴 ALTA |
| 3 | Refactorizar VictimCaseFacade.go | 20-30h | 🟠 MEDIA |
| 4 | Implementar logging estructurado | 16h | 🟠 ALTA |
| 5 | Extraer código duplicado a utilidades | 13h | 🟠 MEDIA |
| 6 | Implementar variables de entorno | 8h | 🟠 MEDIA |
| 7 | Desacoplar salvia de security | 12h | 🟠 MEDIA |

**Total Fase 2:** 144-189 horas (~4-5 sprints)

---

#### 🧪 FASE 3: Testing y Calidad (3-4 meses)

**Objetivo:** Garantizar calidad y prevenir regresiones

| # | Tarea | Esfuerzo | Prioridad |
|---|-------|----------|-----------|
| 1 | Tests unitarios - Prioridad 1 (utils, auth) | 40h | 🔴 ALTA |
| 2 | Tests unitarios - Prioridad 2 (DAOs) | 60h | 🟠 MEDIA |
| 3 | Tests de integración - Facades | 100h | 🟠 MEDIA |
| 4 | Tests frontend (Jest) | 80h | 🟠 MEDIA |
| 5 | Configurar CI/CD con tests | 20h | 🟠 MEDIA |

**Total Fase 3:** 300 horas (~8 sprints)

---

#### 🚀 FASE 4: Optimizaciones (2-3 meses)

**Objetivo:** Mejorar rendimiento y experiencia de usuario

| # | Tarea | Esfuerzo | Prioridad |
|---|-------|----------|-----------|
| 1 | Implementar caché para datos estáticos | 16h | 🟡 MEDIA |
| 2 | Migrar sesiones a Redis | 24h | 🟡 MEDIA |
| 3 | Optimizar queries SQL | 40h | 🟡 MEDIA |
| 4 | Migrar frontend a Fetch API | 4h | 🟡 BAJA |
| 5 | Consolidar librerías UI frontend | 40h | 🟡 BAJA |
| 6 | Implementar versionado de API | 16h | 🟡 MEDIA |

**Total Fase 4:** 140 horas (~4 sprints)

---

### 5.5 Métricas de Calidad Actuales vs Objetivo

| Métrica | Actual | Objetivo | Gap |
|---------|--------|----------|-----|
| **Cobertura de Tests** | 0% | 80% | -80% |
| **Líneas por Función (Promedio)** | 150 | <50 | -100 |
| **Complejidad Ciclomática (Promedio)** | 15-20 | <10 | -5 a -10 |
| **Archivos >300 Líneas** | 5 | 0 | -5 |
| **Vulnerabilidades Críticas** | 3 | 0 | -3 |
| **Vulnerabilidades Totales** | 14 | <5 | -9 |
| **Deuda Técnica (horas)** | 643-678 | <100 | -543 a -578 |
| **Versión Go** | 1.19 (EOL) | 1.22+ | Actualizar |

---

### 5.6 Recomendaciones Estratégicas

#### 1. Implementar Desarrollo Guiado por Tests (TDD)
- **Beneficio:** Prevenir regresiones, mejorar diseño
- **Esfuerzo:** Capacitación + 20% tiempo adicional inicial
- **ROI:** Reducción 60% de bugs en producción

#### 2. Adoptar Arquitectura Hexagonal
- **Beneficio:** Desacoplamiento, testabilidad
- **Esfuerzo:** 3-4 meses de refactoring
- **ROI:** Reducción 50% tiempo de mantenimiento

#### 3. Implementar CI/CD con Gates de Calidad
- **Beneficio:** Automatización, calidad consistente
- **Esfuerzo:** 2-3 semanas setup inicial
- **ROI:** Reducción 80% tiempo de deployment

#### 4. Migrar a Microservicios (Largo Plazo)
- **Beneficio:** Escalabilidad, independencia de equipos
- **Esfuerzo:** 6-12 meses
- **ROI:** Escalabilidad horizontal, deployment independiente

#### 5. Implementar Observabilidad (Logging, Metrics, Tracing)
- **Beneficio:** Debugging rápido, monitoreo proactivo
- **Esfuerzo:** 1-2 meses
- **ROI:** Reducción 70% tiempo de resolución de incidentes

---

### 5.7 Conclusiones

#### Fortalezas del Sistema:
✅ Uso correcto de prepared statements (previene SQL Injection)  
✅ Validación exhaustiva en múltiples capas  
✅ Separación modular por dominio  
✅ Uso de HTTPS con TLS  
✅ Middleware pipeline bien estructurado  

#### Debilidades Críticas:
❌ Sin tests unitarios ni de integración (0% cobertura)  
❌ Archivos extremadamente largos (>1500 líneas)  
❌ Vulnerabilidades de seguridad críticas (secret hardcodeado, password expuesto)  
❌ Sin rate limiting en endpoints críticos  
❌ Go 1.19 sin soporte (EOL)  
❌ Deuda técnica estimada en 643-678 horas  

#### Riesgo General:
🔴 **ALTO** - El sistema requiere intervención inmediata en seguridad y refactoring urgente para garantizar mantenibilidad a largo plazo.

---

### 5.8 Próximos Pasos Inmediatos

**Semana 1-2:**
1. ✅ Mover secret key a variable de entorno
2. ✅ Ocultar password en respuestas JSON
3. ✅ Implementar rate limiting en login
4. ✅ Agregar headers de seguridad HTTP

**Semana 3-4:**
5. ✅ Actualizar Go a 1.22+
6. ✅ Implementar logging estructurado
7. ✅ Configurar entorno de staging para pruebas
8. ✅ Documentar proceso de deployment

**Mes 2:**
9. ✅ Comenzar refactoring de VictimCaseController
10. ✅ Implementar primeros tests unitarios
11. ✅ Configurar CI/CD básico
12. ✅ Auditoría de seguridad externa

---

## ANEXOS

### A. Comandos Útiles

```bash
# Análisis de código
go vet ./...
golangci-lint run

# Tests
go test ./... -v -cover

# Actualizar dependencias
go get -u ./...
go mod tidy

# Build
go build -o salvia main.go

# Generar certificados de prueba
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

### B. Herramientas Recomendadas

**Go:**
- golangci-lint (linting)
- go-critic (análisis estático)
- gocyclo (complejidad ciclomática)
- gosec (seguridad)
- testify (testing)

**Frontend:**
- ESLint (linting JavaScript)
- Jest (testing)
- Cypress (E2E testing)
- Lighthouse (performance)

**DevOps:**
- GitHub Actions / GitLab CI
- Docker
- Kubernetes (futuro)
- Prometheus + Grafana (monitoreo)

---

**FIN DEL INFORME**

---

*Generado el: Diciembre 2024*  
*Analista: Sistema Automatizado de Análisis de Código*  
*Versión: 1.0*
