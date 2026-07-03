# `remisiones-psicosocial-component` — Interfaz del Componente

Tabla reutilizable para visualizar y filtrar remisiones de Atención Psicosocial. Consulta `salvia.psychosocial_support` enriquecida con datos del caso, dupla, profesional asignado y sesiones en `team_contact`. Emite eventos hacia el padre para navegación y reasignación masiva.

**Prerequisitos:**
- [M-01 — Migración schema dupla + psychosocial_support](./Flujos/flow-M01-migracion-schema-dupla-psychosocial-support.md)
- [M-02 — submitted_by_team, professional_id y team_contact](./Flujos/flow-M02-migracion-submitted-by-team-professional-id-team-contact.md)

Eventos: [remisiones-psicosocial-component-events.md](./remisiones-psicosocial-component-events.md)

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/remisiones-psicosocial-component.js` | Componente principal |
| `src/frontend/html/salvia/remisiones-psicosocial/remisiones_psicosocial_component.html` | Template Vue |
| `src/frontend/css/remisiones-psicosocial-component.css` | Estilos (prefijo `rps-`) |
| `src/internal/models/psychosocial_support.go` | Modelo — campos + constantes status |
| `src/internal/models/dupla.go` | Modelo tabla `salvia.dupla` |
| `src/internal/models/team_contact.go` | Modelo tabla `salvia.team_contact` *(M-02)* |

---

## Cambios de modelo

### `salvia.psychosocial_support` (M-01 + M-02)

| Campo | Tipo | Descripción |
|---|---|---|
| `submitted_by` | varchar(36), nullable | `general_user_i_code` del agente que remitió. **Solo lectura:** JOIN para **nombre** del remitente |
| `submitted_by_team` | varchar(50), nullable | Equipo del remitente **denormalizado**. Render y filtro E-13 sin JOIN |
| `dupla_id` | varchar(36), nullable | FK lógica a `salvia.dupla.id` |
| `professional_id` | varchar(36), nullable | Profesional asignado (psicóloga o trab. social). Reemplaza `agent_id` |
| `status` | varchar(30), default `'abierto'` | Estados en español snake_case |

**Constantes de status** (patrón `entity_letter.go`):

```go
const (
    PsychosocialSupportStatusAbierto      = "abierto"
    PsychosocialSupportStatusEnGestion    = "en_gestion"
    PsychosocialSupportStatusEnDevolucion = "en_devolucion"
    PsychosocialSupportStatusCerrado      = "cerrado"
)
```

| Valor BD | Label UI (card / badge) | Condición de negocio |
|---|---|---|
| `abierto` | Abiertos | Al crear la remisión |
| `en_gestion` | En gestión | Primer contacto logrado |
| `en_devolucion` | En devolución | No cumplió criterios |
| `cerrado` | Cerrados | Por cualquier motivo |

### `salvia.dupla`

| Campo | Tipo |
|---|---|
| `id` | uuid PK |
| `name` | varchar(36) |
| `psychologist_id` | varchar(36) |
| `social_worker_id` | varchar(36) |
| `created_at`, `updated_at`, `deleted_at` | timestamps |

### `salvia.team_contact` (M-02)

Registra contactos/sesiones del equipo de Atención Psicosocial vinculados a una remisión.

| Campo | Tipo | Uso en UI |
|---|---|---|
| `id` | uuid PK | — |
| `case_id` | varchar(36) | Vínculo al caso |
| `form_submission_id` | varchar(36) | — |
| `dupla_id` | varchar(36) | — |
| `psicosocial_id` | varchar(36) | FK lógica a `psychosocial_support.id` |
| `professional_id` | varchar(36) | Profesional que atendió la sesión |
| `team` | varchar(50) | — |
| `scheduled_date` | timestamptz | — |
| `is_completed` | boolean | Debe ser `true` para contar sesión |
| `is_psico_session` | boolean | Debe ser `true` para contar sesión |
| `status` | varchar(20) | — |
| `scheduled_time` | varchar(8) | — |
| `completed_at` | timestamp | — |
| `summary` | text | — |
| `created_at`, `updated_at`, `deleted_at` | timestamps | — |

**Regla de conteo de sesiones (barra de puntos y filtro E-10):**

```sql
COUNT(*) FROM salvia.team_contact
WHERE psicosocial_id = {remision.id}
  AND is_psico_session = true
  AND is_completed = true
  AND deleted_at IS NULL
```

---

## Constantes de UI (frontend)

```js
const SESSION_DOTS  = 6;                  // puntos fijos en la barra
const SESSION_LABEL = 'Sesiones (4 - 6)'; // quemado en UI
```

Barra de progreso: **6 puntos** fijos; pintar los primeros `sessionCount` (origen: subquery `team_contact`). Label siempre **"Sesiones (4 - 6)"** — no dinámico.

---

## Props del componente

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `mostrarCards` | `Boolean` | No | `false` | Si es `true`, renderiza la fila de cards de resumen encima de los filtros |
| `defaultFilter` | `Object` | No | `{}` | Filtro inicial fijo. Acota remisiones a un profesional o dupla |
| `reasignacion` | `Boolean` | No | `false` | Activa checkbox y botón "Reasignar" |
| `pageSize` | `Number` | No | `20` | Registros por página |

### Shape de `defaultFilter`

Filtro de alcance que se aplica al montar y **no se elimina** con "Limpiar filtros" (E-05). El padre envía **uno** de los dos criterios según la pantalla:

| Campo | Tipo | Uso |
|---|---|---|
| `professional_id` | `string` | Remisiones vinculadas al icode: asignación directa **o** vía dupla donde es `psychologist_id` / `social_worker_id`. Ver [DEC-E08-01](./Flujos/flow-decision-E08-filtro-profesional-incluye-dupla.md) |
| `dupla_id` | `string` | Remisiones donde `psychosocial_support.dupla_id` = id de **una dupla concreta** |

Ejemplos de pantallas padre:

```html
<!-- Listado general (supervisor) -->
<remisiones-psicosocial-component :mostrar-cards="true" />

<!-- Mis remisiones — todo lo del profesional (directo + duplas donde participa) -->
<remisiones-psicosocial-component
  :default-filter="{ professional_id: sessionAgentIcode }"
/>

<!-- Mis remisiones — acotado a una dupla concreta (opcional) -->
<remisiones-psicosocial-component
  :default-filter="{ dupla_id: sessionDuplaId }"
/>
```

Al montar con `defaultFilter`:
- Se copia a `activeFilters` y se envía en cada consulta al backend.
- Si incluye `professional_id`, el autocomplete "Profesional asignada" puede mostrarse pre-seleccionado (opcional en UI).
- Si incluye `dupla_id`, el dropdown "Dupla asignada" queda pre-seleccionado.

---

## Cards de resumen (`mostrarCards === true`)

Ubicación: **encima** de la FilterBar. **Cinco cards** en fila horizontal: **Total** + los **4 estados** de `PsychosocialSupportStatus`.

```
SummaryCardsRow  (.rps-summary-cards)
├── SummaryCard  (.rps-card.rps-card--total)
│   ├── Label  "Total remisiones"
│   └── Value  stats.total
├── SummaryCard  (.rps-card.rps-card--abierto)
│   ├── Label  "Abiertos"
│   └── Value  stats.abierto
├── SummaryCard  (.rps-card.rps-card--en-gestion)
│   ├── Label  "En gestión"
│   └── Value  stats.enGestion
├── SummaryCard  (.rps-card.rps-card--en-devolucion)
│   ├── Label  "En devolución"
│   └── Value  stats.enDevolucion
└── SummaryCard  (.rps-card.rps-card--cerrado)
    ├── Label  "Cerrados"
    └── Value  stats.cerrado
```

| Card | Métrica | Criterio SQL |
|---|---|---|
| Total remisiones | `stats.total` | `COUNT(*)` con el mismo alcance de filtros activos (`defaultFilter` + filtros UI) |
| Abiertos | `stats.abierto` | `status = 'abierto'` |
| En gestión | `stats.enGestion` | `status = 'en_gestion'` |
| En devolución | `stats.enDevolucion` | `status = 'en_devolucion'` |
| Cerrados | `stats.cerrado` | `status = 'cerrado'` |

> **Relación esperada:** `stats.total` = suma de los cuatro contadores por status **cuando no hay filtro de estado activo** y cada remisión tiene exactamente un status. Con `filter_estado_remision` activo, `total` refleja el subconjunto filtrado y solo la card del status seleccionado tendrá valor > 0 (salvo empates en datos).

Estilos sugeridos:

| Card | Clase CSS | Fondo | Texto |
|---|---|---|---|
| Total remisiones | `.rps-card--total` | blanco / neutro | gris oscuro |
| Abiertos | `.rps-card--abierto` | amarillo claro | marrón-amarillo |
| En gestión | `.rps-card--en-gestion` | azul claro | azul oscuro |
| En devolución | `.rps-card--en-devolucion` | naranja claro | naranja oscuro |
| Cerrados | `.rps-card--cerrado` | verde claro | verde oscuro |

Las cards se recargan junto con el listado (mismos filtros activos). Si `mostrarCards === false`, no se llama al endpoint de stats.

---

## Árbol de interfaz

```
remisiones-psicosocial-component  (.rps-wrapper)
│
├── [v-if mostrarCards]
│   SummaryCardsRow  (.rps-summary-cards)     → 5 cards: total + 4 status (E-01)
│
├── FilterBar  (.rps-filter-bar)
│   ├── FilterRow1
│   │   ├── SearchFilter  numero_identidad     → E-03
│   │   ├── SearchFilter  telefono             → E-04
│   │   ├── DropdownFilter estado_remision     → E-09
│   │   └── DropdownFilter sesiones_completadas → E-10
│   ├── FilterRow2
│   │   ├── DropdownFilter dupla_asignada      → E-11
│   │   ├── AutocompleteFilter profesional_asignada → E-07 / E-08
│   │   ├── DropdownFilter nivel_riesgo        → E-12
│   │   └── DropdownFilter equipo_remitente    → E-13
│   └── ClearFiltersBtn  "Limpiar filtros"     → E-05
│
└── TableWrap
    ├── [v-if reasignacion && selectedRemisiones.length > 0]
    │   ReassignBtn  "Reasignar"               → E-15
    ├── <table>  (CASO | REMISIÓN | ESTADO Y ASIGNACIÓN)
    └── PaginationBar                            → E-06
```

---

## Columnas de la tabla

### Bloque CASO

| Elemento UI | Fuente |
|---|---|
| Nombre víctima | `victim_case` / formulario |
| Badge riesgo | `victim_case_form2_risk_level` |
| Código caso + documento | `caseICode`, `docNumber` |
| Teléfono | COALESCE teléfonos víctima |
| Municipio | Caso / formulario |
| "Ver caso →" | Emite `ver-caso` (E-16) |

### Bloque REMISIÓN

| Elemento UI | Fuente |
|---|---|
| Remitido por | `submitted_by` → JOIN `general_user_profile` |
| Equipo remitente | `submitted_by_team` (columna directa) |
| Fecha remisión | `created_at` |
| "Ver remisión →" | Emite `ver-remision` (E-17) |

> **Sin badge tipo ni tags extra** en esta columna (eliminados post-reunión).

### Bloque ESTADO Y ASIGNACIÓN

| Elemento UI | Fuente |
|---|---|
| Badge status | `psychosocial_support.status` → label UI |
| Label sesiones | Texto fijo **"Sesiones (4 - 6)"** |
| Barra de puntos | 6 puntos; pintar `sessionCount` desde `team_contact` |
| Resumen sesiones | `"{sessionCount} realizadas"` (opcional) |
| Tag dupla | `dupla.name` si `dupla_id` presente |
| Profesionales | Dupla (psicóloga + trab. social) o profesional directo vía `professional_id` |
| Sin asignar | Si `dupla_id IS NULL AND professional_id IS NULL` |

---

## Filtros

| `key` | `type` | Label | Criterio |
|---|---|---|---|
| `numero_identidad` | search | Número de identidad | `victim_case_victim_doc_number` ILIKE |
| `telefono` | search | Teléfono | Teléfono víctima COALESCE |
| `estado_remision` | dropdown | Estado remisión | `psychosocial_support.status` |
| `sesiones_completadas` | dropdown | Sesiones completadas | COUNT `team_contact` (0-6) exacto |
| `dupla_asignada` | dropdown | Dupla asignada | `dupla_id` — opciones de `dupla.name` |
| `profesional_asignada` | autocomplete | Profesional asignada | `filter_professional_id` — incluye asignación directa **y** remisiones de duplas donde el usuario es `psychologist_id` o `social_worker_id` ([DEC-E08-01](./Flujos/flow-decision-E08-filtro-profesional-incluye-dupla.md)) |
| `nivel_riesgo` | dropdown | Nivel de riesgo | `victim_case_form2_risk_level` del caso |
| `equipo_remitente` | dropdown | Equipo remitente | `submitted_by_team` (columna directa) |

### Autocomplete — equipos permitidos

| `general_user_team` | Rol mostrado |
|---|---|
| `psicologia` | Psicóloga |
| `trab. social` | Trab. Social |

Endpoint: `GET /api/v1/agents/search-psicosocial?q={text}&limit=10`

---

## Modo reasignación

| Elemento | Comportamiento |
|---|---|
| Checkbox | Primera columna; solo página actual |
| Regla | Mismo `status` en todas las seleccionadas |
| Botón | **"Reasignar"** — visible si hay selección |
| Emisión | `$emit('reasignar-remisiones', { remisiones })` → E-15 |
| Modal | **Fuera de alcance** — el padre lo implementará después |

---

## Eventos emitidos

| Evento | Cuándo | Payload |
|---|---|---|
| `ver-caso` | "Ver caso →" | `{ caseICode, remision }` |
| `ver-remision` | "Ver remisión →" | `{ remisionId, followUpId, remision }` |
| `reasignar-remisiones` | "Reasignar" | `{ remisiones: Array }` |

---

## Endpoints backend

| Método | Ruta | Evento |
|---|---|---|
| GET | `/api/v1/psychosocial-support/list` | E-01, E-03…E-13, E-06 |
| GET | `/api/v1/psychosocial-support/stats` | E-01 (si `mostrarCards`) |
| GET | `/api/v1/duplas` | E-01, E-11 |
| GET | `/api/v1/psychosocial-support/equipos-remitentes` | E-01, E-13 |
| GET | `/api/v1/agents/search-psicosocial` | E-07 |

### Response `GET /api/v1/psychosocial-support/stats`

```json
{
  "total": 8,
  "abierto": 2,
  "enGestion": 3,
  "enDevolucion": 0,
  "cerrado": 3
}
```

Acepta los mismos query params de filtro que el listado (incluidos `filter_professional_id`, `filter_dupla_id` de `defaultFilter`).

Query params del listado (combinables AND):

| Param | Filtro |
|---|---|
| `filter_numero_identidad` | E-03 |
| `filter_telefono` | E-04 |
| `filter_estado_remision` | E-09 |
| `filter_sesiones_completadas` | E-10 |
| `filter_dupla_id` | E-11 |
| `filter_professional_id` | E-08 — asignación directa **OR** duplas del profesional ([DEC-E08-01](./Flujos/flow-decision-E08-filtro-profesional-incluye-dupla.md)) |
| `filter_nivel_riesgo` | E-12 |
| `filter_equipo_remitente` | E-13 |
| `page`, `page_size` | E-06 |
| `sort`, `order` | default `created_at` / `desc` |

### Campo `sessionCount` en ítem del listado

Subquery por fila (misma regla que barra de puntos):

```sql
(
  SELECT COUNT(*)::int
  FROM salvia.team_contact tc
  WHERE tc.psicosocial_id = ps.id
    AND tc.is_psico_session = true
    AND tc.is_completed = true
    AND tc.deleted_at IS NULL
) AS session_count
```
