# `case-entities` — Guía de uso

Lista las entidades institucionales relacionadas con un caso, con filtros y un modal para agregar nuevas relaciones. "Ver Detalle" navega directamente (no emite evento); no tiene funcionalidad de quitar una entidad en esta versión.

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `caseId` | String | Sí | — | `victim_case_i_code` del caso cuyas entidades se listan |
| `userId` | String | Sí | — | icode del agente en sesión — se envía como `createdById` al agregar una entidad |
| `userRole` | String | Sí | — | Rol del usuario en sesión — controla si se muestra el botón "+ Agregar entidad" (`sv`/`op`/`ro`). El backend valida el permiso real de forma independiente; este prop solo gobierna la UI. |

---

## Eventos emitidos

Ninguno. "Ver Detalle" navega directamente a `/salvia/entidad/:id` (ver GAP en `case-entities-index.md` — esa ruta es nueva, no existe hoy). El resto de acciones (agregar, filtrar) son enteramente internas — no requieren intervención del padre.

> No hay funcionalidad de quitar una entidad en esta versión (decisión de producto — ver `case-entities-interface.md`).

---

## Qué hace el componente por sí solo

Al montarse, hace `GET /api/v1/casos/:caseId/entidades` y muestra las entidades como cards con sus contadores de oficios y barreras activas, y su última acción (`lastAction`), ya resueltos por el backend. Maneja internamente los estados de carga, error y vacío, los filtros (sector, buscador, barreras activas) y el modal de agregar entidad (con selector de ubicación en cascada + sector, ambos obligatorios, + búsqueda de sedes existentes). No requiere que el padre le pase la lista de barreras u oficios — a diferencia del mockup React original, todo el cálculo de contadores ocurre en el backend.

---

## Integración

```html
{{ template "components/case_entities.html" . }}
```

En el template del padre (`get_case_detail_sv.html`, tab "Gestión institucional" — **antes** de la lista de barreras, igual que `SeccionEntidades` en el mockup):

```html
<case-entities
  :case-id="caseICode"
  :user-id="userICode"
  :user-role="userRole">
</case-entities>
```

---

## Endpoints que consume

| Acción | Método | URL | Permiso requerido |
|---|---|---|---|
| Listar entidades del caso | GET | `/api/v1/casos/:caseId/entidades` | Sesión válida + rol con acceso a `get_case_detail_sv` (todos menos `ad`) |
| Agregar entidad al caso | POST | `/api/v1/casos/:caseId/entidades` | Lo anterior **y además** rol `sv`, `op` o `ro` (403 si no) |

### Endpoint reutilizado (ya existe)

| Acción | Método | URL |
|---|---|---|
| Buscar sedes disponibles por ubicación/sector (para el modal) | GET | `/api/v1/entity-branches?town_code={townCode}&sector={sector}` |

### Shape esperado — `GET /api/v1/casos/:caseId/entidades`

```json
[
  {
    "relId": "uuid",
    "entityBranchId": 12,
    "entityBranchICode": "uuid",
    "entityBranchName": "Hospital Nacional San Rafael",
    "sector": "he",
    "address": "25 Av. Norte, Barrio San Miguelito",
    "departmentName": "San Salvador",
    "cityName": "San Salvador",
    "townName": "San Salvador Centro",
    "objetivo": "Atención médica y valoración de lesiones.",
    "oficiosCount": 2,
    "barrerasActivasCount": 0,
    "lastAction": null,
    "createdById": "icode del agente",
    "createdAt": "2026-06-01T10:00:00Z"
  }
]
```

> **Notas:**
> - `sector` viene como código de 2 caracteres (`entity.entity_sector`) — el frontend resuelve la etiqueta igual que ya hace `case_detail` con el sector de barreras.
> - `barrerasActivasCount` cuenta `barrier_v2` por `case_id` + `entity_branch_id` (status `OPEN` o `En Gestion`). Será `0` para todas las entidades hasta que el formulario de registro de barrera permita elegir `entity_branch_id` (cambio fuera de alcance de este componente — ver `related-tables.md`).
> - `lastAction` es el campo `entity_case.last_action` tal cual — texto libre editado directamente, no una fecha calculada. `null` mientras no se le dé un mecanismo para poblarlo (GAP pendiente).
