# Mis Barreras — Tablas relacionadas

## Cadena de alcance

```
security.general_user (sesión UserICode)
    ▼
salvia.barrier_v2
  WHERE user_territory_link_id = :UserICode
    AND status IN ('OPEN', 'En Gestion')
    AND deleted_at IS NULL
    │
    ├─ victim_case ON case_id = victim_case_i_code
    ├─ entity_branch ON entity_branch_id
    ├─ security.city ON city_id::text = b.city_id
    └─ LATERAL case_task ToDo type=articulacion_territorial
         assigned_user_id = :UserICode
```

El cliente no envía el icode del enlace. `requireEnlaceAccess` exige rol `en` y departamento asignado.

La columna `user_territory_link_id` se agrega por AutoMigrate; **quién la escribe** al crear/activar la barrera está fuera de este requerimiento.

---

## `salvia.barrier_v2` (columnas usadas)

| Columna | Uso |
|---|---|
| `id` | PK; enlace a detalle |
| `case_id` | → `victim_case.victim_case_i_code` |
| `user_territory_link_id` | Filtro de alcance (icode del Enlace) |
| `city_id` | JOIN a `security.city`; filtro de ciudad |
| `entity_branch_id` | Entidad del modal / filtro |
| `sector` | Badge |
| `description` | Texto de la barrera |
| `status` | Solo `OPEN` y `En Gestion` |
| `deleted_at` | Soft-delete |

---

## `salvia.case_task`

| Tipo | Estado | Uso |
|---|---|---|
| `articulacion_territorial` | ToDo (listado) / Done (al guardar) | Columna Tarea; se completa o se crea ya Done |
| `compromiso_con_enlace` | ToDo | Se crea si el compromiso requiere seguimiento; `form_data.fechaSeguimiento` |

---

## Contrato API

`GET /api/v1/barriers-v2/mine` — array de `BarrierV2MineItem` (incluye `victimPhone`, `cityId`, `cityName`; `caseId` solo para el link).

`POST /api/v1/case-tasks/articulacion-territorial`

```json
{
  "barrierId": "uuid",
  "taskId": "uuid opcional",
  "formData": {
    "fechaContacto": "YYYY-MM-DD",
    "compromiso": "texto",
    "requiereSeguimiento": true,
    "fechaSeguimiento": "YYYY-MM-DD",
    "sugerenciasRemision": "texto o null",
    "entityBranchId": 1
  }
}
```

Auth: sesión + rol `en`. La barrera debe tener `user_territory_link_id` = usuario de sesión.
