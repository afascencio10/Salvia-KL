# Arquitectura del Sistema SALVIA
## Ministerio de Igualdad - Sistema de Información para Víctimas

---

## 📐 ÍNDICE

1. [Arquitectura General](#1-arquitectura-general)
2. [Arquitectura Backend](#2-arquitectura-backend)
3. [Arquitectura Frontend](#3-arquitectura-frontend)
4. [Comunicación entre Componentes](#4-comunicación-entre-componentes)
5. [Modelo de Datos](#5-modelo-de-datos)
6. [Flujos de Datos](#6-flujos-de-datos)
7. [Seguridad](#7-seguridad)

---

## 1. ARQUITECTURA GENERAL

### 1.1 Vista de Alto Nivel

```
┌─────────────────────────────────────────────────────────────────┐
│                        CAPA CLIENTE                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Navegador  │  │   Navegador  │  │   Navegador  │         │
│  │   Chrome     │  │   Firefox    │  │     Edge     │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│         │                  │                  │                 │
│         └──────────────────┴──────────────────┘                 │
│                            │                                     │
│                      HTTPS (TLS 1.2+)                           │
└────────────────────────────┼────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                    CAPA PRESENTACIÓN                            │
│                            │                                     │
│  ┌─────────────────────────▼──────────────────────────┐        │
│  │         Gin Web Framework (Go)                     │        │
│  │  ┌──────────────────────────────────────────────┐ │        │
│  │  │  Middlewares                                 │ │        │
│  │  │  • Rate Limiting (9999 concurrent)           │ │        │
│  │  │  • HTTPS Redirect                            │ │        │
│  │  │  • Session Management                        │ │        │
│  │  │  • Authentication                            │ │        │
│  │  │  • Static File Server                        │ │        │
│  │  └──────────────────────────────────────────────┘ │        │
│  │                                                     │        │
│  │  ┌──────────────────────────────────────────────┐ │        │
│  │  │  Router                                      │ │        │
│  │  │  • /salvia/*    → Salvia Module             │ │        │
│  │  │  • /seguridad/* → Security Module           │ │        │
│  │  │  • /public/*    → Public Endpoints          │ │        │
│  │  │  • /static/*    → Static Files              │ │        │
│  │  └──────────────────────────────────────────────┘ │        │
│  └─────────────────────────────────────────────────────┘        │
└────────────────────────────┼────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                    CAPA LÓGICA DE NEGOCIO                       │
│                            │                                     │
│  ┌─────────────────────────▼──────────────────────────┐        │
│  │              MÓDULOS FUNCIONALES                   │        │
│  │                                                     │        │
│  │  ┌──────────────┐  ┌──────────────┐  ┌─────────┐ │        │
│  │  │   COMMON     │  │    SALVIA    │  │SECURITY │ │        │
│  │  │              │  │              │  │         │ │        │
│  │  │ • Config     │  │ • Casos      │  │ • Auth  │ │        │
│  │  │ • Utils      │  │ • Contactos  │  │ • Users │ │        │
│  │  │ • DB Pool    │  │ • Alertas    │  │ • Roles │ │        │
│  │  │ • Validation │  │ • Entidades  │  │ • Login │ │        │
│  │  │ • Sessions   │  │ • Momentos   │  │         │ │        │
│  │  └──────────────┘  └──────────────┘  └─────────┘ │        │
│  │                                                     │        │
│  │  Cada módulo tiene 3 capas:                        │        │
│  │  • Facades     (HTTP Handlers)                     │        │
│  │  • Controllers (Business Logic)                    │        │
│  │  • DAOs        (Data Access)                       │        │
│  └─────────────────────────────────────────────────────┘        │
└────────────────────────────┼────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                    CAPA DE ACCESO A DATOS                       │
│                            │                                     │
│  ┌─────────────────────────▼──────────────────────────┐        │
│  │         Connection Pool Manager                    │        │
│  │         (pgx/v4 - PostgreSQL Driver)               │        │
│  │                                                     │        │
│  │  Pool Size: 80 conexiones                          │        │
│  │  Timeout: Configurable                             │        │
│  │  Retry Logic: Implementado                         │        │
│  └─────────────────────────────────────────────────────┘        │
└────────────────────────────┼────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    CAPA DE PERSISTENCIA                         │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              PostgreSQL 12+                              │  │
│  │                                                           │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │  │
│  │  │   Schema:    │  │   Schema:    │  │   Schema:    │  │  │
│  │  │   security   │  │   salvia     │  │   common     │  │  │
│  │  │              │  │              │  │              │  │  │
│  │  │ • Users      │  │ • Cases      │  │ • Logs       │  │  │
│  │  │ • Roles      │  │ • Contacts   │  │ • Config     │  │  │
│  │  │ • Locations  │  │ • Alerts     │  │ • Audit      │  │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘  │  │
│  └──────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### 1.2 Tecnologías Principales

| Capa | Tecnología | Versión | Propósito |
|------|-----------|---------|-----------|
| Backend | Go (Golang) | 1.19+ | Lenguaje principal |
| Web Framework | Gin | 1.9.1 | HTTP routing y middleware |
| Base de Datos | PostgreSQL | 12+ | Persistencia de datos |
| DB Driver | pgx/v4 | 4.18.3 | Driver PostgreSQL |
| Frontend | Vue.js | 3.x | Framework JavaScript |
| UI Framework | Bootstrap | 4.x | Componentes UI |
| Mapas | Leaflet | Latest | Visualización geográfica |
| Sesiones | gin-contrib/sessions | 0.0.5 | Gestión de sesiones |
| TLS/SSL | OpenSSL | 1.1+ | Certificados HTTPS |

---

## 2. ARQUITECTURA BACKEND

### 2.1 Patrón de Capas (3-Tier Architecture)

```
┌─────────────────────────────────────────────────────────────────┐
│                      CAPA 1: FACADES                            │
│                   (Presentation Layer)                          │
│                                                                  │
│  Responsabilidades:                                             │
│  • Recibir peticiones HTTP                                      │
│  • Validar sesiones y permisos                                  │
│  • Parsear parámetros de URL y body                             │
│  • Renderizar templates HTML o devolver JSON                    │
│  • Manejar errores HTTP (4xx, 5xx)                              │
│                                                                  │
│  Archivos: */facades/*Facade.go                                 │
│                                                                  │
│  Ejemplo:                                                        │
│  func VictimContactPOST(c *gin.Context) {                       │
│      session := sessions.Default(c)                             │
│      if !CheckPermission(...) { return }                        │
│      code, res := controller.SetVictimContact(...)              │
│      c.JSON(code, res)                                          │
│  }                                                               │
└─────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                   CAPA 2: CONTROLLERS                           │
│                   (Business Logic Layer)                        │
│                                                                  │
│  Responsabilidades:                                             │
│  • Implementar reglas de negocio                                │
│  • Validar datos de entrada (tipos, rangos, formatos)          │
│  • Coordinar operaciones entre múltiples DAOs                   │
│  • Manejar transacciones complejas                              │
│  • Transformar datos entre formatos                             │
│  • Aplicar lógica de autorización específica                    │
│                                                                  │
│  Archivos: */controllers/*Controller.go                         │
│                                                                  │
│  Ejemplo:                                                        │
│  func SetVictimContact(dataInput string, ...) (int, string) {   │
│      // Validar JSON                                            │
│      dto := ValidateJSONInput(dataInput)                        │
│      // Llamar DAO                                              │
│      err := dao.SetVictimContact(dto, ...)                      │
│      // Formatear respuesta                                     │
│      return 200, FormatJSON(dto)                                │
│  }                                                               │
└─────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      CAPA 3: DAOs                               │
│                   (Data Access Layer)                           │
│                                                                  │
│  Responsabilidades:                                             │
│  • Ejecutar queries SQL (SELECT, INSERT, UPDATE, DELETE)        │
│  • Mapear resultados DB ↔ DTOs (Data Transfer Objects)         │
│  • Gestionar conexiones del pool                                │
│  • Manejar transacciones de base de datos                       │
│  • Implementar prepared statements (prevenir SQL injection)     │
│                                                                  │
│  Archivos: */dao/*DAO.go                                        │
│                                                                  │
│  Ejemplo:                                                        │
│  func SetVictimContact(dto *VictimContactDTO, ...) error {      │
│      conn := GetConnection(...)                                 │
│      query := "INSERT INTO victim_contact (...) VALUES (...)"   │
│      _, err := conn.Exec(ctx, query, params...)                 │
│      ReleaseConnection(conn)                                    │
│      return err                                                  │
│  }                                                               │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Estructura de Módulos

```
Fuentes/
│
├── main.go                          # Punto de entrada
│   └── Inicializa routers y arranca servidor
│
├── common/                          # Módulo compartido
│   ├── config/                      # Configuraciones globales
│   │   ├── Enums.go                # Constantes y enumeraciones
│   │   ├── DateTime.go             # Formatos de fecha/hora
│   │   └── Locale.go               # Traducciones i18n
│   │
│   ├── controllers/                 # Controladores comunes
│   │   ├── PersistenceController.go # Gestión de conexiones
│   │   └── By.go                   # Criterios de búsqueda
│   │
│   ├── dao/                         # DTOs comunes
│   │   └── CommonDTO.go            # Estructuras compartidas
│   │
│   ├── db/                          # Gestión de base de datos
│   │   ├── PostgresConnection.go   # Pool de conexiones
│   │   ├── DBClientConfig.go       # Configuración cliente
│   │   └── DBServerConfig.go       # Configuración servidor
│   │
│   ├── facades/                     # Router principal
│   │   └── MainRouter.go           # Configuración Gin
│   │
│   └── utils/                       # Utilidades
│       ├── CommonSession.go        # Gestión de sesiones
│       ├── EntityValidator.go      # Validación de datos
│       ├── EntityOperator.go       # Operaciones genéricas
│       ├── CommMsg.go              # Mensajes de respuesta
│       ├── FileManager.go          # Gestión de archivos
│       └── HTMLHelper.go           # Helpers para templates
│
├── salvia/                          # Módulo de casos de víctimas
│   ├── config/                      # Configuración específica
│   │   ├── Locale.go               # Traducciones del módulo
│   │   ├── PermissionsByRole.go    # Permisos por rol
│   │   ├── MenuTools.go            # Menús de navegación
│   │   └── FormPaths.go            # Rutas de formularios
│   │
│   ├── controllers/                 # Lógica de negocio
│   │   ├── VictimContactController.go
│   │   ├── VictimCaseController.go
│   │   ├── AlertController.go
│   │   ├── MomentController.go
│   │   ├── EntityController.go
│   │   └── CaseLogController.go
│   │
│   ├── dao/                         # Acceso a datos
│   │   ├── VictimContactDAO.go     # CRUD contactos
│   │   ├── VictimCaseDAO.go        # CRUD casos
│   │   ├── AlertDAO.go             # CRUD alertas
│   │   ├── MomentDAO.go            # CRUD momentos
│   │   ├── EntityDAO.go            # CRUD entidades
│   │   ├── EntityBranchDAO.go      # CRUD sucursales
│   │   └── CaseLogDAO.go           # CRUD logs
│   │
│   ├── facades/                     # Handlers HTTP
│   │   ├── MainRouter.go           # Rutas del módulo
│   │   ├── VictimContactFacade.go
│   │   ├── VictimCaseFacade.go
│   │   ├── AlertFacade.go
│   │   ├── MomentFacade.go
│   │   └── EntityFacade.go
│   │
│   └── scripts/                     # Scripts SQL
│       └── db_postgres_creation.sql
│
└── security/                        # Módulo de seguridad
    ├── config/                      # Configuración de seguridad
    │   ├── Locale.go
    │   ├── PermissionsByRole.go
    │   └── FormPaths.go
    │
    ├── controllers/                 # Lógica de autenticación
    │   ├── GeneralUserController.go
    │   ├── CityController.go
    │   └── TownController.go
    │
    ├── dao/                         # Acceso a datos de usuarios
    │   ├── GeneralUserDAO.go
    │   ├── RoleDAO.go
    │   ├── CityDAO.go
    │   ├── TownDAO.go
    │   └── DepartmentDAO.go
    │
    ├── facades/                     # Handlers de autenticación
    │   ├── MainRouter.go
    │   ├── GeneralUserFacade.go
    │   ├── CityFacade.go
    │   └── TownFacade.go
    │
    └── scripts/                     # Scripts SQL
        └── db_postgres_creation.sql
```


### 2.3 Flujo de Procesamiento de Petición

```
┌──────────────────────────────────────────────────────────────────┐
│ 1. CLIENTE ENVÍA PETICIÓN                                        │
│    POST https://localhost/salvia/contacto-victima                │
│    Body: { "names": "María", "phone": "3001234567", ... }       │
└────────────────────────────┬─────────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────────┐
│ 2. GIN ROUTER RECIBE PETICIÓN                                    │
│    router.POST("/salvia/contacto-victima", VictimContactPOST)   │
└────────────────────────────┬─────────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────────┐
│ 3. MIDDLEWARES PROCESAN (en orden)                               │
│    ┌──────────────────────────────────────────────────────────┐ │
│    │ a) limitMiddleware()                                     │ │
│    │    • Verifica límite de 9999 peticiones concurrentes    │ │
│    │    • Si excede: bloquea hasta que haya espacio          │ │
│    └──────────────────────────────────────────────────────────┘ │
│    ┌──────────────────────────────────────────────────────────┐ │
│    │ b) ForceHTTPS()                                          │ │
│    │    • Verifica que la conexión sea HTTPS                 │ │
│    │    • Si es HTTP: redirige a HTTPS                       │ │
│    └─────