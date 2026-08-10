# Casos Entidad — Tablas relacionadas

> Actualizado con la relación usuario `et` ↔ sede (`entity_branch`).

---

## Cadena de relación (alcance del listado)

```
security.general_user
    │  entity_branch_id   ← campo NUEVO (nullable; obligatorio en práctica para rol et)
    ▼
salvia.entity_branch          (sede = el usuario et)
    │  entity_id                 entity_branch_town_code
    ▼                            ▼
salvia.entity                 ciudad de la sede (implícita)
(nombre + sector en header)
    │
    ▼
salvia.entity_case
  WHERE entity_branch_id = :sedeDelUsuario
  AND deleted_at IS NULL
    │  case_id
    ▼
salvia.victim_case (+ form2 riesgo, town→city del caso para display)
```

**Regla de negocio:** un usuario con rol `et` = una sede (`entity_branch`). No cambia de entidad. El listado **no** agrega por organización: solo la sede del usuario.

---

## Cambio de schema — `security.general_user` (aplicado)

| Columna nueva | Tipo | Nulo | Descripción |
|---|---|---|---|
| `entity_branch_id` | bigint | Sí (NULL para roles ≠ et) | FK lógica → `salvia.entity_branch.entity_branch_id` |

Ver script: [`migration-general-user-entity-branch-id.sql`](./migration-general-user-entity-branch-id.sql)

### Impacto en código (implementado)

| Pieza | Cambio |
|---|---|
| `GeneralUserDAO` / DTO | Campo `GeneralUserEntityBranchId` → `entity_branch_id` |
| Login / `CommonSession` | `EntityBranchId` en sesión |
| Facade / API Casos Entidad | Sede desde sesión; listado `ec.entity_branch_id = ?` |
| `main.go` → `migrateLegacyTables` | AutoMigrate de `security.general_user` crea la columna en deploy |
| Alta/edición de usuarios `et` | **Pendiente:** UI setup para obligar sede |

> Nota: `EntityBrandICode` (perfil) sigue existiendo por legado; la fuente de verdad para Casos Entidad es `general_user.entity_branch_id`.

---

## `salvia.entity_case`

| Columna | Tipo | Descripción |
|---|---|---|
| `id` | varchar(36) PK | UUID del vínculo |
| `case_id` | varchar(36) | → `victim_case.victim_case_i_code` |
| `entity_branch_id` | int64 | → sede; **filtro del listado et** |
| `objetivo` | text | “Relación con entidad” |
| `last_action` | text | Texto libre (no timeline) |
| `created_by_id` | varchar(36) | Auditoría |
| `created_at` / `updated_at` / `deleted_at` | — | Soft-delete GORM |

---

## Tablas de apoyo

| Tabla | Uso |
|---|---|
| `security.general_user` | `entity_branch_id` del usuario et |
| `entity_branch` | Sede del usuario; `entity_id` + `town_code` |
| `entity` | Nombre y sector en el header |
| `victim_case` | Víctima, documento, estado, ciudad del caso (display) |
| `victim_case_form2` | Nivel de riesgo |

---

## Contrato API (ajustado)

### Listado — `GET /api/v1/entity-cases`

| Param | Origen | Notas |
|---|---|---|
| *(sede)* | **Sesión** `EntityBranchId` | Required; si falta → 403/400 |
| `document` | query (E02) | ILIKE parcial opcional |
| `page` | query (E04) | 0-based |
| `pageSize` | query | default **5** |

```
WHERE ec.deleted_at IS NULL
  AND ec.entity_branch_id = :sessionEntityBranchId
  AND [document ILIKE si aplica]
ORDER BY ec.updated_at DESC
LIMIT 5 OFFSET page*5
```

Respuesta sugerida: `{ items, total, page, pageSize, meta: { entityName, sectorName, entityBranchId } }`

### Endpoints que esta pantalla **deja de usar**

| Endpoint | Motivo |
|---|---|
| `GET /api/v1/entities` | Ya no hay selector de organización |
| `GET /api/v1/entities/:id/cities` | Ya no hay filtro de ciudad |

(Pueden quedar para otros módulos.)

---

## DTO / fila UI

| Campo UI | Origen |
|---|---|
| `entityCaseId` | `entity_case.id` |
| `caseId` / `caseCode` | `case_id` / `victim_case_i_code` |
| `entityBranchICode` | sede (E06) |
| `victimFullName` | nombres + apellidos del caso |
| `document` | doc del caso |
| `city` | ciudad del caso (solo display en fila) |
| `riskLevel` / `caseStatus` | form2 / status |
| `objetivo` / `lastAction` | `entity_case` |
| Header `entityName` / `sectorName` | `entity` vía sede del usuario |

---

## GAPS / pendientes

| Ítem | Impacto |
|---|---|
| Poblar `entity_branch_id` en usuarios et | E01 (datos) |
| UI alta/edición usuario: obligar sede para rol et | Operación |
| Pantalla `/salvia/entidad/:id` | E06 |
| Mecanismo para editar `last_action` | Contenido de columna |
