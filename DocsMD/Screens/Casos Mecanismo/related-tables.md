# Casos Mecanismo — Tablas relacionadas

> Recorte del listado: usuario `mc` ↔ `assigned_user_id` + `status = 'active'`.

---

## Cadena de relación (alcance del listado)

```
security.general_user  (rol mc, sesión)
    │  general_user_i_code
    ▼
salvia.case_articulating_mechanism
  WHERE assigned_user_id = :sessionUserICode
    AND status = 'active'
    AND deleted_at IS NULL
    │  case_id
    ▼
salvia.victim_case (+ form2 riesgo, town→city)
    │
    └── salvia.case_task
          WHERE case_articulating_mechanism_id = cam.id
            AND status = 'Done'
          ORDER BY completed_at DESC
          LIMIT 1
          → result / completed_at
```

**Regla de negocio:** el listado **nunca** acepta un `userId` del cliente. El backend usa siempre el icode de sesión.

---

## `salvia.case_articulating_mechanism`

| Columna | Tipo | Descripción |
|---|---|---|
| `id` | varchar(36) PK | UUID del escalamiento |
| `case_id` | varchar(36) | → `victim_case.victim_case_i_code` |
| `scope` | varchar(20) | `municipal` \| `departamental` \| `nacional` |
| `status` | varchar(20) | En esta pantalla solo `active` |
| `assigned_user_id` | varchar(36) | Usuario `mc` asignado (**filtro del listado**) |
| `salvia_commitment` | text | Columna “Compromiso de Salvia” |
| `created_at` / `updated_at` / `deleted_at` | — | Soft-delete GORM |

El registro **no se cierra** al gestionar: cada gestión es una tarea nueva.

---

## `salvia.case_task` (gestión)

| Columna | Uso en esta pantalla |
|---|---|
| `type` | `gestion_mecanismo` |
| `category` | `Mecanismos` |
| `status` | Siempre `Done` al crear |
| `form_data` | JSON del formulario según `scope` |
| `result` | Resumen concatenado (`label: valor`. …) → columna “Último comentario” |
| `case_articulating_mechanism_id` | FK opcional al escalamiento (índice) |
| `assigned_user_id` | Usuario `mc` que guardó |
| `completed_at` | `now` al insertar |

---

## Tablas de apoyo

| Tabla | Uso |
|---|---|
| `victim_case` | Nombre, documento, `i_code`, status |
| `victim_case_form2` | Nivel de riesgo y teléfono de la víctima |
| `security.town` / `security.city` | Ciudad de la víctima (filtro + display) |

---

## Contrato API

### Listado — `GET /api/v1/mecanismos/mine`

| Param | Origen | Notas |
|---|---|---|
| *(usuario)* | **Sesión** `UserICode` | Required; rol `mc` |
| `document` | query (E02) | ILIKE parcial opcional |
| `city` | query (E03) | ILIKE parcial opcional sobre nombre de ciudad |
| `page` | query (E04) | 0-based |
| `pageSize` | query | default **5** |

```
WHERE cam.deleted_at IS NULL
  AND cam.status = 'active'
  AND cam.assigned_user_id = :sessionUserICode
  AND [document ILIKE si aplica]
  AND [city ILIKE si aplica]
ORDER BY cam.created_at DESC
LIMIT 5 OFFSET page*5
```

Respuesta: `{ items, total, page, pageSize }`

### Crear gestión — `POST /api/v1/case-tasks/gestion-mecanismo`

Body: `{ "mechanismId": "...", "formData": { ... } }`

El mecanismo debe existir, estar `active` y `assigned_user_id` = sesión. El formulario debe coincidir con `scope`. Inserta `case_task` Done; **nunca** actualiza una tarea previa. Tras el insert escribe un `case_timeline_event` (categoría `Mecanismos`, type `Gestión del mecanismo`). Si el timeline falla, se loguea WARN y la tarea permanece Done.

---

## DTO / fila UI

| Campo UI | Origen |
|---|---|
| `id` | `case_articulating_mechanism.id` |
| `caseId` / `caseCode` | `case_id` / `victim_case_i_code` (JSON; el i_code no se pinta en la fila) |
| `victimFullName` | nombres + apellidos del caso |
| `document` | documento de la víctima (fila: `Doc:`) |
| `phone` | `victim_case_form2_victim_phone` (fila: `Tel:`) |
| `city` | ciudad de la víctima |
| `riskLevel` / `caseStatusLabel` | form2 / status mapeado |
| `scope` | `municipal` \| `departamental` \| `nacional` |
| `salviaCommitment` | `salvia_commitment` |
| `lastComment` / `lastCommentAt` | `case_task.result` / `completed_at` |
