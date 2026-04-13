# QA — Módulo Seguimientos del Área

Documentación técnica y plan de pruebas para la integración Frontend (Vue 3 CDN) + Backend (Go/Gin/GORM).

## URL de la página

```
https://localhost/salvia/seguimientos/area
```

---

## Modo Desarrollo (acceso sin login)

Para facilitar las pruebas durante el desarrollo, la página funciona sin requerir sesión activa. Cuando no hay usuario logueado, muestra TODOS los casos y seguimientos de todos los equipos.

### Cómo funciona

El archivo `salvia/facades/SeguimientosAreaFacade.go` tiene una variable al inicio:

```go
var seguimientosAreaDevMode = true   // ← true = sin login, trae todo
```

Cuando `seguimientosAreaDevMode = true`:
- No requiere sesión activa (no redirige al landing)
- No valida permisos por rol
- Inyecta `TeamId = ""` (vacío = sin filtro de equipo)
- Inyecta `IsAdmin = true` y `UserId = "dev-user"`
- Las APIs aceptan `team` vacío y retornan todos los registros

Adicionalmente, el `AuthMiddleware` en `common/facades/MainRouter.go` excluye `/salvia/seguimientos` de la validación de sesión.

### Cómo desactivarlo para producción

1. En `salvia/facades/SeguimientosAreaFacade.go`, cambiar:

```go
var seguimientosAreaDevMode = false  // ← false = requiere login y permisos
```

2. En `common/facades/MainRouter.go`, quitar `!strings.HasPrefix(location, "/salvia/seguimientos")` de la condición del `AuthMiddleware`:

```go
// ANTES (modo dev):
if ... && !strings.HasPrefix(location, "/salvia/seguimientos") {

// DESPUÉS (producción):
if ... {
```

3. En los handlers del controller (`followup_v2_controller.go`), restaurar la validación obligatoria de `team` en `ListByArea`, `AgentWorkload` y `FilterOptions`:

```go
// Descomentar estas líneas en cada handler:
if team == "" {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "team requerido"})
    return
}
```

Con estos 3 cambios, la página vuelve a exigir sesión + rol `sv` o `ad` + team asignado.

---

## 1. Qué se hizo

Se implementó la vista "Seguimientos del Área" que permite a los supervisores consultar, filtrar y reagendar los seguimientos de todos los agentes de su equipo.

La integración funciona así: Gin renderiza el HTML del template `seguimientos_area.html` inyectando variables del servidor (UserId, IsAdmin, TeamId) directamente en el `data()` de Vue. Una vez montada la app Vue, el frontend consume las APIs REST vía `fetch()` para obtener datos dinámicos sin recargar la página.

### Endpoints creados

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/salvia/seguimientos/area` | Renderiza el HTML (template Go + Vue) |
| `GET` | `/api/v1/seguimientos` | Lista paginada con filtros dinámicos |
| `GET` | `/api/v1/seguimientos/carga-agentes` | Carga pendiente agrupada por agente |
| `GET` | `/api/v1/seguimientos/filtros-opciones` | Catálogos (agentes del área + estados) |
| `PUT` | `/api/v1/seguimientos/:id/reagendar` | Reagenda un seguimiento |

### Archivos creados/modificados

| Archivo | Tipo |
|---|---|
| `frontend/html/salvia/follow_up_v2/seguimientos_area.html` | HTML root + Vue app |
| `frontend/js/components/seguimientos-filtros.js` | Componente Vue |
| `frontend/js/components/carga-agente-bar.js` | Componente Vue |
| `frontend/js/components/seguimientos-tabla.js` | Componente Vue |
| `salvia/facades/SeguimientosAreaFacade.go` | Facade Go (renderizado HTML) |
| `salvia/facades/MainRouter.go` | Registro de ruta `/seguimientos/area` |
| `salvia/controller/followup_v2_controller.go` | 4 handlers nuevos (API) |
| `salvia/service/followup_v2_service.go` | 4 métodos nuevos |
| `internal/repository/followup_repository.go` | 4 queries nuevas |

---

## 2. Cómo se hizo (Arquitectura Frontend)

### Flujo de datos unidireccional

```
seguimientos_area.html (Estado global en data())
    │
    ├── Props ↓ hacia componentes hijos
    │   ├── <carga-agente-bar :agentes="cargaAgentes">
    │   ├── <seguimientos-filtros :agentes="opcionesAgentes" :estados="opcionesEstados">
    │   └── <seguimientos-tabla :items="seguimientos" :total="total" :page="page">
    │
    └── Eventos ↑ desde componentes hijos
        ├── @filtrar → aplicarFiltros() → fetch GET /api/v1/seguimientos
        ├── @paginar → cambiarPagina() → fetch GET /api/v1/seguimientos?page=N
        ├── @reagendar → abrirReagendar() → modal → fetch PUT /api/v1/seguimientos/:id/reagendar
        ├── @seleccionar → seleccionarItem() → abre panel lateral de detalle
        └── @cambiar-fecha → cambiarFechaCarga() → fetch GET /api/v1/seguimientos/carga-agentes
```

El estado vive exclusivamente en el `data()` del root. Los componentes hijos son "tontos": reciben datos por props y comunican acciones por eventos. Nunca hacen fetch directamente.

### Regla de delimiters

Todos los componentes y el root usan `delimiters: ['${', '}']` para evitar conflictos con los templates `{{ }}` de Go. Si un desarrollador olvida esto, Vue intentará interpretar las variables de Go y la página quedará en blanco.

### Orden de scripts (inamovible)

```html
<!-- 1. Cargar Vue.js y dependencias -->
{{ template "standard_scripts.html" .}}

<!-- 2. Crear la app Vue con createApp + data/methods -->
<script> var app = Vue.createApp({ ... }); </script>

<!-- 3. Registrar componentes (usan app.component()) -->
<script src="/static/js/components/carga-agente-bar.js"></script>
<script src="/static/js/components/seguimientos-filtros.js"></script>
<script src="/static/js/components/seguimientos-tabla.js"></script>

<!-- 4. Montar la app -->
<script> app.mount("#seguimientos-area-app"); </script>
```

Si se altera este orden (por ejemplo, cargar un componente antes del `createApp`), la variable `app` no existe y el `app.component()` falla silenciosamente.

---

## 3. Cómo validar que funciona (Plan de Pruebas QA)

### Prerequisitos

- Servidor corriendo: `go run main.go` desde `src/`
- Migración ejecutada: `migrations/add_followup_v2_calendar_fields.sql` y `migrations/add_user_team.sql`
- Al menos un usuario con rol `sv` y campo `general_user_team` = `RIESGO_BAJO` o `RIESGO_ALTO`
- Al menos 5 seguimientos en la tabla `salvia.follow_up_v2` con distintos estados y agentes

### Escenario 1: Filtrado por equipo del supervisor

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | Iniciar sesión como Supervisor General (rol `sv`, team `RIESGO_BAJO`) | Login exitoso |
| 2 | Navegar a `/salvia/seguimientos/area` | Se renderiza la página con sidebar morado |
| 3 | Verificar la tabla | Solo muestra seguimientos donde `team = 'RIESGO_BAJO'` |
| 4 | Abrir DevTools → Network → verificar la petición a `/api/v1/seguimientos` | El query param `team=RIESGO_BAJO` está presente |
| 5 | Repetir con Supervisor Integral (team `RIESGO_ALTO`) | Solo muestra seguimientos de `RIESGO_ALTO` |

### Escenario 2: Filtrado por agente desde la barra de carga

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | En la barra "CARGA POR AGENTE", observar los pills | Cada pill muestra el nombre del agente y su número de pendientes |
| 2 | Cambiar la fecha en el selector de la barra | Los pills se actualizan con la carga de esa fecha |
| 3 | En el filtro "Todos los agentes", seleccionar un agente específico | La tabla se filtra mostrando solo los seguimientos de ese agente |
| 4 | Verificar en Network que la petición incluye `agente_id=XXX` | El parámetro está presente y la respuesta solo contiene registros de ese agente |

### Escenario 3: Reagendar un seguimiento

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | En la tabla, ubicar un seguimiento con estado `PENDIENTE` | Debe mostrar el link "Ejecutar" en la columna Estado |
| 2 | Hacer clic en "Ejecutar" | Se abre el modal "Re agendar y priorizar" |
| 3 | Llenar: nueva fecha, motivo (obligatorios), hora y prioridad (opcionales) | El botón "Confirmar" se habilita |
| 4 | Hacer clic en "Confirmar" | El modal se cierra |
| 5 | Verificar en Network: `PUT /api/v1/seguimientos/:id/reagendar` | Respuesta 200 con `{"message": "seguimiento reagendado exitosamente"}` |
| 6 | Verificar la tabla | La fila se actualiza reactivamente (estado cambia a `REPROGRAMADO`) sin recargar la página |
| 7 | Verificar la barra de carga | El contador del agente se actualiza (un pendiente menos) |

### Escenario 4: Tabs de filtrado rápido

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | Hacer clic en tab "Vista rápida" | Muestra solo seguimientos `PENDIENTE` |
| 2 | Hacer clic en tab "Esta semana" | Filtra por rango lunes-domingo de la semana actual |
| 3 | Hacer clic en tab "Por fecha" | Aparecen inputs de fecha inicio/fin |
| 4 | Hacer clic en tab "Histórico" | Muestra todos los seguimientos sin filtro de estado |

### Escenario 5: Panel lateral de detalle

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | Hacer clic en cualquier fila de la tabla | Se abre un panel lateral derecho con el detalle del seguimiento |
| 2 | Verificar que muestra: Caso, Riesgo, Fecha, Hora, Estado, Seguimiento # | Todos los datos coinciden con la fila seleccionada |
| 3 | Hacer clic en la X del panel | El panel se cierra |

### Escenario 6: Acceso denegado

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | Iniciar sesión como Agente (rol `op`) | Login exitoso |
| 2 | Navegar manualmente a `/salvia/seguimientos/area` | Respuesta 401 (no tiene permiso `get_seguimientos_area`) |

### Escenario 7: Paginación

| Paso | Acción | Resultado esperado |
|---|---|---|
| 1 | Con más de 20 seguimientos, verificar que aparece la paginación | Muestra "Página 1 de N — X registros" |
| 2 | Hacer clic en "Siguiente" | La tabla carga la página 2, el contador se actualiza |
| 3 | Hacer clic en "Anterior" | Vuelve a la página 1 |

---

## 4. Comandos útiles para QA

```bash
# Levantar el servidor
cd c:\Desarrollo\SOG_SALVIA\src
go run main.go

# Insertar datos de prueba (ajustar UUIDs y team según tu BD)
psql -h localhost -U root -d salvia_pruebas_q9xt -c "
INSERT INTO salvia.follow_up_v2 (case_id, agent_id, team, risk_status, scheduled_date, status, sequence_number)
VALUES
    ('caso-001', 'agente-lopez', 'RIESGO_BAJO', 'ALTO', '2026-04-15', 'PENDIENTE', 1),
    ('caso-001', 'agente-lopez', 'RIESGO_BAJO', 'ALTO', '2026-04-17', 'PENDIENTE', 2),
    ('caso-002', 'agente-ramirez', 'RIESGO_BAJO', 'MODERADO', '2026-04-16', 'PENDIENTE', 1),
    ('caso-003', 'agente-diaz', 'RIESGO_ALTO', 'EXTREMO', '2026-04-11', 'REALIZADO', 1),
    ('caso-003', 'agente-diaz', 'RIESGO_ALTO', 'EXTREMO', '2026-04-12', 'PENDIENTE', 2);
"

# Asignar team a un supervisor de prueba
psql -h localhost -U root -d salvia_pruebas_q9xt -c "
UPDATE security.general_user SET general_user_team = 'RIESGO_BAJO' WHERE general_user_login = 'supervisor_test';
"

# Probar API directamente con cURL
curl -k "https://localhost/api/v1/seguimientos?team=RIESGO_BAJO&page=0&limit=20"
curl -k "https://localhost/api/v1/seguimientos/carga-agentes?team=RIESGO_BAJO&fecha=2026-04-15"
curl -k "https://localhost/api/v1/seguimientos/filtros-opciones?team=RIESGO_BAJO"
```
