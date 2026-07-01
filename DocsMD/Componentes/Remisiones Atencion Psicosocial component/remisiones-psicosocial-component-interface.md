# `remisiones-psicosocial-component` — Interfaz del Componente

Tabla reutilizable para visualizar y filtrar remisiones de Atención Psicosocial. Consulta `salvia.psychosocial_support` enriquecida con datos del caso, dupla asignada y agente remitente. Emite eventos hacia el padre para navegación y reasignación masiva.

**Prerequisito:** [M-01 — Migración schema](./Flujos/flow-M01-migracion-schema-dupla-psychosocial-support.md)

Eventos: [remisiones-psicosocial-component-events.md](./remisiones-psicosocial-component-events.md)

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/components/remisiones-psicosocial-component.js` | Componente principal *(pendiente)* |
| `src/frontend/css/remisiones-psicosocial-component.css` | Estilos *(pendiente)* |
| `src/internal/models/psychosocial_support.go` | Modelo — campos nuevos + constantes status |
| `src/internal/models/dupla.go` | Modelo tabla `salvia.dupla` *(pendiente)* |

---

## Cambios de modelo (M-01)

### `salvia.psychosocial_support`

| Campo | Tipo | Descripción |
|---|---|---|
| `submitted_by` | varchar(36), nullable | `general_user_i_code` del agente que remitió. **Solo lectura:** el componente lee el id guardado y resuelve nombre + equipo vía JOIN a `general_user` / `general_user_profile`. No implementa la lógica que lo persiste al crear la remisión |
| `dupla_id` | varchar(36), nullable | FK lógica a `salvia.dupla.id` |
| `agent_id` | varchar(36), nullable | Profesional asignado (psicóloga o trab. social) |
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

| Valor BD | Label UI | Condición de negocio |
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

---

## Constantes de UI (frontend)

```js
const MAX_SESSIONS = 6;   // quemado hasta definir origen real
```

Barra de progreso: `MAX_SESSIONS` puntos; llenar `sessionCount`. Label: `"Sesiones (0-{MAX_SESSIONS})"`.

---

## Props del componente

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `mostrarCards` | `Boolean` | No | `false` | Si es `true`, renderiza la fila de cards de resumen encima de los filtros |
| `defaultFilter` | `Object` | No | `{}` | Filtro inicial fijo. Permite cargar la pantalla acotada a remisiones de un profesional o dupla |
| `reasignacion` | `Boolean` | No | `false` | Activa checkbox y botón "Reasignar" |
| `pageSize` | `Number` | No | `20` | Registros por página |

### Shape de `defaultFilter`

Filtro de alcance que se aplica al montar y **no se elimina** con "Limpiar filtros" (E-05). El padre envía **uno** de los dos criterios según la pantalla:

| Campo | Tipo | Uso |
|---|---|---|
| `agent_id` | `string` | Remisiones donde `psychosocial_support.agent_id` = icode del psicólogo/trab. social |
| `dupla_id` | `string` | Remisiones donde `psychosocial_support.dupla_id` = id de la dupla |

Ejemplos de pantallas padre:

```html
<!-- Listado general de psicosocial (supervisor) -->
<remisiones-psicosocial-component :mostrar-cards="true" />

<!-- Mis remisiones — psicólogo asignado directamente -->
<remisiones-psicosocial-component
  :default-filter="{ agent_id: sessionAgentIcode }"
/>

<!-- Mis remisiones — vista por dupla -->
<remisiones-psicosocial-component
  :default-filter="{ dupla_id: sessionDuplaId }"
/>
```

Al montar con `defaultFilter`:
- Se copia a `activeFilters` y se envía en cada consulta al backend.
- Si incluye `agent_id`, el autocomplete "Profesional asignada" puede mostrarse pre-seleccionado (opcional en UI).
- Si incluye `dupla_id`, el dropdown "Dupla asignada" queda pre-seleccionado.

---

## Cards de resumen (`mostrarCards === true`)

Ubicación: **encima** de la FilterBar. Cuatro cards en fila horizontal (mockup).

```
SummaryCardsRow  (.rps-summary-cards)
├── SummaryCard  (.rps-card.rps-card--total)
│   ├── Label  "Total remisiones"
│   └── Value  stats.total
├── SummaryCard  (.rps-card.rps-card--pendiente)
│   ├── Label  "Pendiente asignación"
│   └── Value  stats.pendienteAsignacion
├── SummaryCard  (.rps-card.rps-card--proceso)
│   ├── Label  "En proceso"
│   └── Value  stats.enProceso
└── SummaryCard  (.rps-card.rps-card--cerradas)
    ├── Label  "Cerradas"
    └── Value  stats.cerradas
```

| Card | Métrica | Criterio SQL |
|---|---|---|
| Total remisiones | `stats.total` | COUNT con el mismo alcance de filtros activos (`defaultFilter` + filtros UI) |
| Pendiente asignación | `stats.pendienteAsignacion` | `dupla_id IS NULL AND agent_id IS NULL` |
| En proceso | `stats.enProceso` | `status = 'en_gestion'` |
| Cerradas | `stats.cerradas` | `status = 'cerrado'` |

Estilos sugeridos (mockup):

| Card | Fondo | Texto |
|---|---|---|
| Total remisiones | blanco / neutro | gris oscuro |
| Pendiente asignación | amarillo claro | marrón-amarillo |
| En proceso | azul claro | azul oscuro |
| Cerradas | verde claro | verde oscuro |

Las cards se recargan junto con el listado (mismos filtros activos). Si `mostrarCards === false`, no se llama al endpoint de stats.

---

## Árbol de interfaz

```
remisiones-psicosocial-component  (.rps-wrapper)
│
├── [v-if mostrarCards]
│   SummaryCardsRow  (.rps-summary-cards)     → datos de E-01 / fetchStats
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

### Bloque REMISIÓN — `submitted_by`

| Elemento UI | Fuente |
|---|---|
| Remitido por | `submitted_by` → JOIN `general_user_profile` (nombres + apellidos) |
| Equipo remitente | `submitted_by` → JOIN `general_user.general_user_team` |

> El componente **solo lee** el id almacenado en `submitted_by`. No implementa cómo se persiste al crear la remisión.

### Bloque REMISIÓN — placeholders

| Elemento | Valor UI actual |
|---|---|
| ReferralTypeBadge | **"Por consultar"** |
| ExtraTags | **"Por consultar"** |

> Fuente de datos pendiente de definir. Documentar aquí cuando se conozca el origen.

---

## Filtros

| `key` | `type` | Label | Criterio |
|---|---|---|---|
| `numero_identidad` | search | Número de identidad | `victim_case_victim_doc_number` ILIKE |
| `telefono` | search | Teléfono | Teléfono víctima COALESCE |
| `estado_remision` | dropdown | Estado remisión | `psychosocial_support.status` |
| `sesiones_completadas` | dropdown | Sesiones completadas | `session_count` exacto 0-6 |
| `dupla_asignada` | dropdown | Dupla asignada | `dupla_id` — opciones de `dupla.name` |
| `profesional_asignada` | autocomplete | Profesional asignada | `agent_id` — equipos `psicologia`, `trab. social` |
| `nivel_riesgo` | dropdown | Nivel de riesgo | `victim_case_form2_risk_level` del caso |
| `equipo_remitente` | dropdown | Equipo remitente | `general_user_team` del usuario en `submitted_by` (solo lectura) |

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
  "pendienteAsignacion": 2,
  "enProceso": 3,
  "cerradas": 3
}
```

Acepta los mismos query params de filtro que el listado (incluidos `filter_agent_id`, `filter_dupla_id` de `defaultFilter`).

Query params del listado (combinables AND):

| Param | Filtro |
|---|---|
| `filter_numero_identidad` | E-03 |
| `filter_telefono` | E-04 |
| `filter_estado_remision` | E-09 |
| `filter_sesiones_completadas` | E-10 |
| `filter_dupla_id` | E-11 |
| `filter_agent_id` | E-08 |
| `filter_nivel_riesgo` | E-12 |
| `filter_equipo_remitente` | E-13 |
| `page`, `page_size` | E-06 |
| `sort`, `order` | default `created_at` / `desc` |
