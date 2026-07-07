# Notificaciones — API de listado paginado

Endpoint usado por la pantalla de Notificaciones para cargar oficios con paginación y filtros en base de datos.

## Endpoint

```
GET /api/v1/entity-letters
```

## Roles y parámetros de usuario

| Rol | Parámetro requerido | Descripción |
|---|---|---|
| `op` | `agentId` | Agente de seguimiento — oficios donde es agent_id |
| `ro` | `agentId` | Revisor operativo — mismo alcance que `op` |
| `an` | `notificationUserId` | Agente de notificaciones — oficios donde es notification_user_id |

La validación de rol ocurre en frontend al montar la pantalla. El backend filtra por el ID enviado.

## Parámetros de paginación

| Parámetro | Tipo | Default | Descripción |
|---|---|---|---|
| `page` | int | 0 | Página base-0 |
| `limit` | int | 20 | Registros por página. **Requerido** para respuesta paginada |

> Sin `limit`, el endpoint devuelve el array completo (compatibilidad con `oficios-list.js`).

## Parámetros de filtro

| Parámetro | Tipo | Descripción |
|---|---|---|
| `state` | string | Estado exacto del oficio (`por_proyectar`, `para_revisar`, etc.) |
| `identidad` | string | ILIKE sobre `victim_case.victim_case_victim_doc_number` |
| `entidad` | string | ILIKE sobre `barrier_v2.sector` |
| `numeroRadicado` | string | ILIKE sobre `entity_letter.numero_radicado` |
| `manageableOnly` | bool | Si `true`, solo estados gestionables según rol del usuario |

### Estados gestionables (`manageableOnly=true`)

| Rol | Estados incluidos |
|---|---|
| `op`, `ro` | `por_proyectar`, `en_correccion` |
| `an` | `para_revisar`, `aprobacion_juridica`, `para_radicar`, `radicado` |

## Respuesta paginada

```json
{
  "items": [ /* EntityLetterWithRelations[] */ ],
  "total": 42,
  "page": 0,
  "pageSize": 5,
  "pendingCount": 7
}
```

| Campo | Descripción |
|---|---|
| `items` | Oficios de la página actual, enriquecidos con caso y barrera |
| `total` | Total de registros que coinciden con filtros activos |
| `page` | Página devuelta |
| `pageSize` | Tamaño de página usado |
| `pendingCount` | Oficios gestionables del usuario (sin aplicar filtros de UI) |

## Ejemplo — Agente de seguimiento, tab gestionar, página 2

```
GET /api/v1/entity-letters?page=1&limit=5&agentId=USR001&manageableOnly=true
```

## Ejemplo — Agente de notificaciones con filtros

```
GET /api/v1/entity-letters?page=0&limit=5&notificationUserId=USR002&state=para_revisar&identidad=52.123
```

## Implementación backend

| Archivo | Responsabilidad |
|---|---|
| `src/salvia/controller/entity_letter_controller.go` | Parseo de query params y respuesta HTTP |
| `src/salvia/service/entity_letter_service.go` | Estados gestionables por rol |
| `src/internal/repository/entity_letter_repository.go` | SQL con JOINs, filtros ILIKE y COUNT paginado |
