# Arquitectura de Software - Sistema SALVIA
## Ministerio de Igualdad - Sistema de Atención a Víctimas de Violencia de Género

---

## 📐 ÍNDICE

1. [Visión General](#visión-general)
2. [Arquitectura de Alto Nivel](#arquitectura-de-alto-nivel)
3. [Arquitectura Backend](#arquitectura-backend)
4. [Arquitectura Frontend](#arquitectura-frontend)
5. [Comunicación entre Componentes](#comunicación-entre-componentes)
6. [Modelo de Datos](#modelo-de-datos)
7. [Flujos de Datos](#flujos-de-datos)
8. [Patrones de Diseño](#patrones-de-diseño)

---

## 1. VISIÓN GENERAL

### 1.1 Descripción del Sistema

SALVIA es un sistema web para la gestión integral de casos de víctimas de violencia de género, que permite:

- Registro de contactos iniciales de víctimas
- Gestión completa de casos
- Asignación de operadores y entidades de atención
- Seguimiento de momentos del proceso de atención
- Generación de alertas y reportes
- Control de acceso basado en roles

### 1.2 Tecnologías Principales

```
┌─────────────────────────────────────────────┐
│           STACK TECNOLÓGICO                 │
├─────────────────────────────────────────────┤
│ Backend:    Go 1.19 + Gin Framework         │
│ Frontend:   HTML + Vue.js 3 + JavaScript    │
│ Base Datos: PostgreSQL 12+                  │
│ Servidor:   HTTPS (TLS 1.2+)                │
│ Sesiones:   Cookie-based (in-memory)        │
└─────────────────────────────────────────────┘
```

### 1.3 Características Arquitectónicas

- **Arquitectura Monolítica Modular**: Un solo ejecutable con módulos separados
- **Patrón MVC Adaptado**: Facades (View) → Controllers → DAOs (Model)
- **Renderizado Híbrido**: Server-side (Go templates) + Client-side (Vue.js)
- **Autenticación Basada en Sesiones**: Cookie-based con middleware
- **Pool de Conexiones**: 80 conexiones concurrentes a PostgreSQL

---

## 2. ARQUITECTURA DE ALTO NIVEL

### 2.1 Diagrama de Contexto

```
┌──────────────────────────────────────────────────────────────┐
│                    SISTEMA SALVIA                            │
│                                                              │
│  ┌────────────┐         ┌──────────────┐                   │
│  │  Víctimas  │────────▶│   Frontend   │                   │
│  │  (Público) │         │  (Vue.js +   │                   │
│  └────────────┘         │   HTML)      │                   │
│                         └──────┬───────┘                   │
│                                │                            │
│  ┌────────────┐                │                            │
│  │ Operadores │────────────────┤                            │
│  │   (Auth)   │                │                            │
│  └────────────┘                ▼                            │
│                         ┌──────────────┐                   │
│  ┌────────────┐         │   Backend    │                   │
│  │Supervisores│────────▶│  (Go + Gin)  │                   │
│  │   (Auth)   │         └──────┬───────┘                   │
│  └────────────┘                │                            │
│                                 │                            │
│  ┌────────────┐                ▼                            │
│  │   Admin    │         ┌──────────────┐                   │
│  │   (Auth)   │────────▶│  PostgreSQL  │                   │
│  └────────────┘         │   Database   │                   │
│                         └──────────────┘                   │
└──────────────────────────────────────────────────────────────┘
```

### 2.2 Arquitectura de Capas

```
┌─────────────────────────────────────────────────────────────┐
│                    CAPA DE PRESENTACIÓN                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Browser    │  │   Vue.js     │  │  Templates   │     │
│  │   (HTTPS)    │  │  Components  │  │   Go HTML    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            ▼ HTTPS
┌─────────────────────────────────────────────────────────────┐
│                   CAPA DE APLICACIÓN                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              GIN ROUTER + MIDDLEWARES                 │  │
│  │  • Authentication  • Rate Limiting  • HTTPS Force    │  │
│  └──────────────────────────────────────────────────────┘  │
│                            ▼                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   FACADES    │  │   FACADES    │  │   FACADES    │     │
│  │   (Common)   │  │   (Salvia)   │  │  (Security)  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   CAPA DE LÓGICA DE NEGOCIO                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ CONTROLLERS  │  │ CONTROLLERS  │  │ CONTROLLERS  │     │
│  │   (Common)   │  │   (Salvia)   │  │  (Security)  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   CAPA DE ACCESO A DATOS                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │     DAOs     │  │     DAOs     │  │     DAOs     │     │
│  │   (Common)   │  │   (Salvia)   │  │  (Security)  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                            ▼                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         CONNECTION POOL (80 conexiones)              │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   CAPA DE PERSISTENCIA                      │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              PostgreSQL Database                      │  │
│  │  • Schema: security  • Schema: salvia                │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```


---

## 3. ARQUITECTURA BACKEND

### 3.1 Estructura Modular

```
Backend (Go)
│
├── main.go                    # Punto de entrada
│
├── common/                    # Módulo compartido
│   ├── config/               # Configuraciones globales
│   ├── controllers/          # Lógica común
│   ├── dao/                  # DTOs comunes
│   ├── db/                   # Gestión de conexiones
│   ├── facades/              # Router principal
│   └── utils/                # Utilidades
│
├── salvia/                    # Módulo de casos
│   ├── config/               # Config específica
│   ├── controllers/          # Lógica de casos
│   │   ├── VictimContactController.go
│   │   ├── VictimCaseController.go
│   │   ├── AlertController.go
│   │   ├── MomentController.go
│   │   └── EntityController.go
│   ├── dao/                  # Modelos de datos
│   │   ├── VictimContactDAO.go
│   │   ├── VictimCaseDAO.go
│   │   ├── AlertDAO.go
│   │   ├── MomentDAO.go
│   │   └── EntityDAO.go
│   └── facades/              # HTTP Handlers
│       ├── MainRouter.go
│       ├── VictimContactFacade.go
│       ├── VictimCaseFacade.go
│       ├── AlertFacade.go
│       └── MomentFacade.go
│
└── security/                  # Módulo de seguridad
    ├── config/               # Config de seguridad
    ├── controllers/          # Lógica de auth
    │   ├── GeneralUserController.go
    │   ├── CityController.go
    │   └── TownController.go
    ├── dao/                  # Modelos de usuarios
    │   ├── GeneralUserDAO.go
    │   ├── RoleDAO.go
    │   ├── CityDAO.go
    │   └── TownDAO.go
    └── facades/              # HTTP Handlers
        ├── MainRouter.go
        ├── GeneralUserFacade.go
        ├── CityFacade.go
        └── TownFacade.go
```

### 3.2 Componentes Backend Detallados

#### 3.2.1 Main Router (common/facades/MainRouter.go)

```go
Responsabilidades:
├── Inicializar Gin Engine
├── Configurar middlewares globales
│   ├── limitMiddleware()      → Control de concurrencia (9999 req)
│   ├── ForceHTTPS()           → Redirección HTTP → HTTPS
│   ├── AuthMiddleware()       → Verificación de sesión
│   └── sessions.Sessions()    → Gestión de sesiones
├── Servir archivos estáticos (/static → /frontend)
├── Cargar templates HTML
└── Iniciar servidor HTTPS (puerto 443)
```

#### 3.2.2 Facades (Capa de Presentación)

```go
Responsabilidades:
├── Recibir peticiones HTTP
├── Extraer parámetros (URL, Query, Body)
├── Verificar sesión activa
├── Verificar permisos del usuario
├── Llamar a Controllers
├── Renderizar respuesta
│   ├── HTML (templates Go)
│   └── JSON (API responses)
└── Manejar errores HTTP
```

**Ejemplo: VictimContactFacade.go**
```
VictimContactPOST()        → Crear contacto
VictimContactGET()         → Obtener contacto(s)
VictimContactPUT()         → Actualizar/invalidar
VictimContactPOST_GET()    → Renderizar formulario
VictimContactPOST_Public() → Endpoint público
```

#### 3.2.3 Controllers (Capa de Lógica de Negocio)

```go
Responsabilidades:
├── Validar datos de entrada
│   ├── Tipos de datos
│   ├── Rangos de valores
│   ├── Campos requeridos
│   └── Formato (emails, teléfonos, etc.)
├── Aplicar reglas de negocio
│   ├── Cálculo de edad
│   ├── Evaluación de riesgo
│   ├── Asignación de operadores
│   └── Generación de alertas
├── Coordinar operaciones entre DAOs
├── Gestionar transacciones
└── Formatear respuestas
```

**Ejemplo: VictimContactController.go**
```
SetVictimContact()                    → Crear/actualizar
GetVictimContactByICode()             → Obtener por ID
GetVictimContactsWithoutVictimCase()  → Listar sin caso
InvalidateVictimContactByICode()      → Soft delete
```

#### 3.2.4 DAOs (Capa de Acceso a Datos)

```go
Responsabilidades:
├── Definir DTOs (Data Transfer Objects)
├── Mapear DTOs ↔ Base de Datos
├── Ejecutar queries SQL
│   ├── SELECT
│   ├── INSERT
│   ├── UPDATE
│   └── DELETE
├── Gestionar conexiones del pool
├── Manejar transacciones
└── Convertir tipos de datos (Go ↔ PostgreSQL)
```

**Ejemplo: VictimContactDAO.go**
```
VictimContactDTO              → Estructura de datos
SetVictimContact()            → INSERT/UPDATE
GetVictimContact()            → SELECT por criterio
GetVictimContacts()           → SELECT múltiples
InvalidateVictimContact()     → UPDATE status
```

### 3.3 Gestión de Conexiones a Base de Datos

```
┌─────────────────────────────────────────────┐
│         CONNECTION POOL                     │
│                                             │
│  ┌─────┐ ┌─────┐ ┌─────┐       ┌─────┐    │
│  │Conn1│ │Conn2│ │Conn3│  ...  │Conn80│   │
│  └─────┘ └─────┘ └─────┘       └─────┘    │
│     ▲       ▲       ▲              ▲       │
│     │       │       │              │       │
│     └───────┴───────┴──────────────┘       │
│                  │                          │
└──────────────────┼──────────────────────────┘
                   │
         ┌─────────▼─────────┐
         │  GetConnection()  │
         │  - Obtiene conn   │
         │  - Crea si no hay │
         └─────────┬─────────┘
                   │
         ┌─────────▼─────────┐
         │   DAO ejecuta     │
         │   query SQL       │
         └─────────┬─────────┘
                   │
         ┌─────────▼─────────┐
         │ ReleaseConnection()│
         │ - Devuelve al pool│
         └───────────────────┘
```

**Código:**
```go
// common/db/PostgresConnection.go
func GetConnection(connData *ConnData, 
                   clientConfig *DBClientConfig, 
                   serverConfig *DBServerConfig) (*ConnData, error) {
    // Obtiene conexión del pool o crea nueva
    // Pool size: 80 conexiones
}

func ReleaseConnection(connData *ConnData) error {
    // Devuelve conexión al pool
}
```

### 3.4 Middlewares

```
┌──────────────────────────────────────────────────────┐
│              PIPELINE DE MIDDLEWARES                 │
└──────────────────────────────────────────────────────┘
                        │
                        ▼
        ┌───────────────────────────┐
        │   limitMiddleware()       │
        │   Control concurrencia    │
        │   Max: 9999 peticiones    │
        └───────────┬───────────────┘
                    ▼
        ┌───────────────────────────┐
        │   ForceHTTPS()            │
        │   Redirige HTTP → HTTPS   │
        └───────────┬───────────────┘
                    ▼
        ┌───────────────────────────┐
        │   sessions.Sessions()     │
        │   Gestiona cookies        │
        └───────────┬───────────────┘
                    ▼
        ┌───────────────────────────┐
        │   AuthMiddleware()        │
        │   Verifica autenticación  │
        │   Redirige a login si no  │
        └───────────┬───────────────┘
                    ▼
        ┌───────────────────────────┐
        │   Handler (Facade)        │
        │   Procesa petición        │
        └───────────────────────────┘
```

---

## 4. ARQUITECTURA FRONTEND

### 4.1 Estructura Frontend

```
frontend/
│
├── html/                      # Páginas HTML estáticas
│   ├── casos.html
│   ├── contactos.html
│   └── reportes.html
│
├── js/                        # JavaScript
│   ├── vue.js                # Framework Vue.js 3
│   ├── dao.js                # Capa de acceso a datos (AJAX)
│   ├── formulario.js         # Lógica de formularios
│   ├── helper.js             # Funciones auxiliares
│   ├── html_helper.js        # Helpers para HTML
│   ├── login.js              # Lógica de login
│   ├── registro.js           # Lógica de registro
│   ├── logs.js               # Gestión de logs
│   ├── momento.js            # Gestión de momentos
│   └── view.js               # Gestión de vistas
│
├── css/                       # Estilos CSS
│   ├── bootstrap.min.css
│   ├── styles.css
│   └── custom.css
│
├── templates/                 # Templates Go HTML
│   ├── salvia/
│   │   ├── victim_contact/
│   │   │   ├── get_victim_contact.html
│   │   │   └── set_victim_contact.html
│   │   └── victim_case/
│   │       ├── get_victim_case.html
│   │       └── set_victim_case.html
│   └── security/
│       └── user/
│           ├── login.html
│           └── profile.html
│
├── plugins/                   # Librerías de terceros
│   ├── leaflet/              # Mapas
│   ├── fontawesome/          # Iconos
│   └── bootstrap/            # UI Framework
│
├── images/                    # Imágenes
├── videos/                    # Videos
└── landing.html              # Página de inicio
```

### 4.2 Arquitectura Híbrida Frontend

```
┌─────────────────────────────────────────────────────────┐
│              RENDERIZADO HÍBRIDO                        │
└─────────────────────────────────────────────────────────┘

SERVER-SIDE (Go Templates)              CLIENT-SIDE (Vue.js)
┌──────────────────────┐                ┌──────────────────┐
│  Backend Go          │                │  Navegador       │
│                      │                │                  │
│  1. Recibe petición  │                │  4. Recibe HTML  │
│  2. Obtiene datos DB │                │  5. Vue.js init  │
│  3. Renderiza HTML   │───────────────▶│  6. Interactivo  │
│     con template     │     HTTPS      │                  │
│                      │                │  7. Usuario      │
│                      │                │     interactúa   │
│                      │                │                  │
│                      │◀───────────────│  8. AJAX request │
│  9. Procesa request  │     HTTPS      │                  │
│ 10. Devuelve JSON    │───────────────▶│ 11. Actualiza UI │
└──────────────────────┘                └──────────────────┘
```

### 4.3 Componentes Frontend

#### 4.3.1 DAO.js (Capa de Acceso a Datos)

```javascript
Responsabilidades:
├── Realizar peticiones AJAX
│   ├── GET
│   ├── POST
│   ├── PUT
│   └── DELETE
├── Manejar respuestas
│   ├── Success (200-299)
│   ├── Errors (400-599)
│   └── Network errors
├── Gestionar cookies/sesiones
└── Formatear datos
```

**Ejemplo:**
```javascript
class DAO {
    async get(url) {
        const response = await fetch(url, {
            method: 'GET',
            credentials: 'same-origin'
        });
        return response.json();
    }
    
    async post(url, data) {
        const response = await fetch(url, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data),
            credentials: 'same-origin'
        });
        return response.json();
    }
}
```

#### 4.3.2 Vue.js Components

```javascript
Responsabilidades:
├── Gestionar estado de la UI
├── Validación de formularios
├── Binding de datos (v-model)
├── Eventos de usuario
├── Renderizado reactivo
└── Comunicación con DAO
```

**Ejemplo:**
```javascript
new Vue({
    el: '#app',
    data: {
        victimContact: {
            names: '',
            phone: '',
            // ... más campos
        }
    },
    methods: {
        async submitForm() {
            const dao = new DAO();
            const result = await dao.post(
                '/salvia/contacto-victima',
                this.victimContact
            );
            // Manejar resultado
        }
    }
});
```


---

## 5. COMUNICACIÓN ENTRE COMPONENTES

### 5.1 Flujo Completo de una Petición

```
┌─────────────────────────────────────────────────────────────────┐
│                    FLUJO DE PETICIÓN HTTP                       │
└─────────────────────────────────────────────────────────────────┘

1. USUARIO INTERACTÚA
   │
   ├─▶ Navegador
   │   └─▶ Vue.js captura evento
   │       └─▶ DAO.js prepara petición
   │
2. PETICIÓN HTTPS
   │
   ├─▶ POST https://localhost/salvia/contacto-victima
   │   Headers: {
   │     Content-Type: application/json,
   │     Cookie: session_id=...
   │   }
   │   Body: {
   │     names: "María",
   │     phone: "3001234567",
   │     ...
   │   }
   │
3. SERVIDOR RECIBE
   │
   ├─▶ Gin Router
   │   └─▶ Middlewares
   │       ├─▶ limitMiddleware()     ✓ Permite (< 9999)
   │       ├─▶ ForceHTTPS()          ✓ Ya es HTTPS
   │       ├─▶ sessions.Sessions()   ✓ Carga sesión
   │       └─▶ AuthMiddleware()      ✓ Usuario autenticado
   │
4. ROUTING
   │
   ├─▶ MainRouter identifica ruta
   │   └─▶ /salvia/contacto-victima → VictimContactPOST()
   │
5. FACADE PROCESA
   │
   ├─▶ VictimContactFacade.VictimContactPOST()
   │   ├─▶ Obtiene sesión
   │   ├─▶ Verifica permiso "set_victim_contact"
   │   ├─▶ Lee body de la petición
   │   └─▶ Llama a Controller
   │
6. CONTROLLER VALIDA
   │
   ├─▶ VictimContactController.SetVictimContact()
   │   ├─▶ Parsea JSON a DTO
   │   ├─▶ Valida campos
   │   │   ├─▶ names: 3-32 caracteres ✓
   │   │   ├─▶ phone: 6-10 dígitos ✓
   │   │   └─▶ ... más validaciones
   │   ├─▶ Aplica reglas de negocio
   │   └─▶ Llama a DAO
   │
7. DAO PERSISTE
   │
   ├─▶ VictimContactDAO.SetVictimContact()
   │   ├─▶ Obtiene conexión del pool
   │   ├─▶ Construye query SQL
   │   │   INSERT INTO salvia.victim_contact (...)
   │   │   VALUES (...)
   │   ├─▶ Ejecuta query
   │   ├─▶ Obtiene ID generado
   │   ├─▶ Libera conexión
   │   └─▶ Retorna DTO con ID
   │
8. RESPUESTA ASCIENDE
   │
   ├─▶ DAO → Controller
   │   └─▶ Controller formatea respuesta JSON
   │       └─▶ Controller → Facade
   │           └─▶ Facade envía respuesta HTTP
   │
9. CLIENTE RECIBE
   │
   ├─▶ HTTP 200 OK
   │   Body: {
   │     "success": true,
   │     "data": {
   │       "icode": "VCT-2024-001",
   │       "names": "María",
   │       ...
   │     }
   │   }
   │
10. FRONTEND ACTUALIZA
    │
    └─▶ DAO.js recibe respuesta
        └─▶ Vue.js actualiza UI
            └─▶ Usuario ve confirmación
```

### 5.2 Comunicación Backend ↔ Base de Datos

```
┌────────────────────────────────────────────────────────┐
│         COMUNICACIÓN CON BASE DE DATOS                 │
└────────────────────────────────────────────────────────┘

DAO                    Connection Pool         PostgreSQL
│                            │                      │
│ 1. GetConnection()         │                      │
├───────────────────────────▶│                      │
│                            │ 2. Obtiene/Crea conn │
│                            ├─────────────────────▶│
│                            │                      │
│ 3. Retorna ConnData        │                      │
│◀───────────────────────────┤                      │
│                            │                      │
│ 4. Ejecuta Query           │                      │
├────────────────────────────┼─────────────────────▶│
│                            │                      │
│                            │ 5. Resultado         │
│◀────────────────────────────┼──────────────────────┤
│                            │                      │
│ 6. ReleaseConnection()     │                      │
├───────────────────────────▶│                      │
│                            │ 7. Devuelve al pool  │
│                            │                      │
```

### 5.3 Comunicación Frontend ↔ Backend

```
┌────────────────────────────────────────────────────────┐
│         COMUNICACIÓN FRONTEND ↔ BACKEND                │
└────────────────────────────────────────────────────────┘

Vue.js Component       DAO.js           Backend Facade
│                        │                     │
│ 1. Usuario submit      │                     │
├───────────────────────▶│                     │
│                        │ 2. AJAX POST        │
│                        ├────────────────────▶│
│                        │                     │
│                        │                     │ 3. Procesa
│                        │                     │    Valida
│                        │                     │    Persiste
│                        │                     │
│                        │ 4. JSON Response    │
│                        │◀────────────────────┤
│                        │                     │
│ 5. Callback            │                     │
│◀───────────────────────┤                     │
│                        │                     │
│ 6. Actualiza UI        │                     │
│                        │                     │
```

### 5.4 Autenticación y Sesiones

```
┌────────────────────────────────────────────────────────┐
│              FLUJO DE AUTENTICACIÓN                    │
└────────────────────────────────────────────────────────┘

1. LOGIN
   │
   Usuario → POST /seguridad/login
   │         Body: {login: "user", password: "pass"}
   │
   ├─▶ GeneralUserFacade.GeneralUserLOGIN_POST()
   │   ├─▶ Valida CAPTCHA
   │   ├─▶ Busca usuario en BD
   │   ├─▶ Verifica contraseña (hash)
   │   ├─▶ Crea sesión
   │   │   └─▶ utils.AddCommonSession(sessionID, userData)
   │   ├─▶ Guarda sessionID en cookie
   │   │   └─▶ session.Set("userData", sessionID)
   │   └─▶ Redirige a dashboard
   │
2. PETICIONES SUBSECUENTES
   │
   Usuario → GET /salvia/contacto-victima
   │         Cookie: session_id=abc123...
   │
   ├─▶ AuthMiddleware()
   │   ├─▶ Lee cookie
   │   ├─▶ Obtiene sessionID
   │   ├─▶ Busca sesión en memoria
   │   │   └─▶ utils.GetCommonSession(sessionID)
   │   ├─▶ Si existe → Continúa
   │   └─▶ Si no existe → Redirige a login
   │
3. VERIFICACIÓN DE PERMISOS
   │
   ├─▶ Facade verifica permiso específico
   │   └─▶ utils.CheckPermission(
   │         permissions,
   │         "set_victim_contact",
   │         userRole,
   │         context
   │       )
   │
4. LOGOUT
   │
   Usuario → POST /seguridad/logout
   │
   ├─▶ GeneralUserFacade.GeneralUserLOGOUT_POST()
   │   ├─▶ Elimina sesión de memoria
   │   │   └─▶ utils.DeleteCommonSession(sessionID)
   │   ├─▶ Invalida cookie
   │   └─▶ Redirige a login
```

---

## 6. MODELO DE DATOS

### 6.1 Diagrama Entidad-Relación General

```
┌─────────────────────────────────────────────────────────────┐
│                    ESQUEMA: security                        │
└─────────────────────────────────────────────────────────────┘

    ┌──────────────────┐
    │   country        │
    │ ──────────────── │
    │ PK country_id    │
    │    country_code  │
    │    country_name  │
    └────────┬─────────┘
             │ 1
             │
             │ N
    ┌────────▼─────────┐
    │   department     │
    │ ──────────────── │
    │ PK department_id │
    │    department_code│
    │    department_name│
    │ FK country_id    │
    └────────┬─────────┘
             │ 1
             │
             │ N
    ┌────────▼─────────┐
    │      city        │
    │ ──────────────── │
    │ PK city_id       │
    │    city_code     │
    │    city_name     │
    │ FK department_id │
    └────────┬─────────┘
             │ 1
             │
             │ N
    ┌────────▼─────────┐
    │      town        │
    │ ──────────────── │
    │ PK town_id       │
    │    town_code     │
    │    town_name     │
    │    latitude      │
    │    longitude     │
    │ FK city_id       │
    └──────────────────┘


┌──────────────────┐         ┌──────────────────┐
│ general_user     │    1    │      role        │
│ ──────────────── │◀────────│ ──────────────── │
│ PK user_id       │    N    │ PK role_id       │
│    user_i_code   │         │    role_code     │
│    login         │         │    role_name     │
│    password      │         │    description   │
│    status        │         └──────────────────┘
│    language      │                 ▲
│ FK profile_id    │                 │
└────────┬─────────┘                 │
         │ 1                         │
         │                           │
         │ 1                         │
┌────────▼─────────┐                 │
│ general_user_    │                 │
│    profile       │                 │
│ ──────────────── │                 │
│ PK profile_id    │                 │
│    names         │                 │
│    last_names    │                 │
│    doc_type      │                 │
│    doc_number    │                 │
│    gender        │                 │
│ FK town_code     │                 │
└──────────────────┘                 │
                                     │
┌────────────────────────────────────┘
│
│  rel_role_general_user
│  (Tabla de relación N:M)
│
└─────────────────────────────────────┐
                                      │
┌─────────────────────────────────────▼──┐
│           ESQUEMA: salvia              │
└────────────────────────────────────────┘
```

