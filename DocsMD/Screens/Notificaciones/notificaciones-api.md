# Notificaciones — API de listado paginado

Endpoint usado por la pantalla de Notificaciones para cargar oficios con paginación y filtros en base de datos.

## Endpoint principal

```
GET /api/v1/entity-letters
```

## Roles y alcance

| Rol | Tab | Parámetros clave |
|---|---|---|
| `op` / `ro` | Todos mis oficios | `agentId={userId}` |
| `op` / `ro` | Oficios por gestionar | `agentId={userId}&manageableOnly=true` |
| `an` | Todos | `listAll=true` |
| `an` | Mis Oficios | `mineOnly=true&notificationAgentId={userId}` |
| `an` | Oficios por gestionar | `listAll=true&manageableOnly=true` |

## Paginación

| Parámetro | Tipo | Default | Descripción |
|---|---|---|---|
| `page` | int | 0 | Página base-0 |
| `limit` | int | 20 | Registros por página (**requerido** para respuesta paginada) |

## Filtros comunes

| Parámetro | Descripción |
|---|---|
| `state` | Estado exacto del oficio |
| `identidad` | ILIKE sobre documento de la víctima |
| `entidad` | ILIKE sobre sector de barrera |
| `numeroRadicado` | ILIKE sobre número radicado |
| `manageableOnly` | Solo estados gestionables según rol |

## Filtros agente de notificaciones (rol `an`)

| Parámetro | Columna BD |
|---|---|
| `notificationUserIdReview` | `notification_user_id_review` |
| `notificationUserIdRadicado` | `notification_user_id_radicado` |
| `notificationUserIdResponse` | `notification_user_id_response` |

## Respuesta paginada

```json
{
  "items": [],
  "total": 42,
  "page": 0,
  "pageSize": 5,
  "pendingCount": 7
}
```

## Catálogo de agentes `an`

```
GET /api/v1/agents/by-role?role=an
```

```json
{
  "agents": [
    { "icode": "USR001", "names": "María", "lastNames": "García", "team": "..." }
  ]
}
```

## Campos nuevos en entity_letter

| JSON | Se asigna en acción |
|---|---|
| `notificationUserIdReview` | `revisar`, `por_corregir` (solo desde `para_revisar`) |
| `notificationUserIdRadicado` | `radicar` (modal aprobar) |
| `notificationUserIdResponse` | `registrar_respuesta` |

Los campos legacy `review_by` y `radicado_by` se mantienen con su lógica actual.
