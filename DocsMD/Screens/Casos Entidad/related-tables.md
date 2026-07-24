# Casos Entidad — Tablas relacionadas

> Schema confirmado vía `src/internal/models/entity_case.go` y análisis en [case-entities/related-tables.md](../../Componentes/case-entities/related-tables.md).

---

## Cadena de relación (alcance del listado)

```
usuario et
    │  (relación usuario↔entidad: AÚN NO DEFINIDA)
    │  E01 carga catálogo completo de salvia.entity
    │  E07: el usuario elige qué organización consultar
    ▼
salvia.entity                 (organización elegida)
    │  entity_id
    ▼
salvia.entity_branch          (todas las sedes de esa organización)
    │
    ▼
salvia.entity_case            WHERE entity_branch.entity_id = entidadElegida
    │  case_id
    ▼
victim_case / victim_contact  (datos de la fila)
```

Una fila de `entity_case` apunta a una **sede** (`entity_branch_id`). El listado agrupa por **organización** (`entity_id`): muestra todos los `entity_case` cuyas sedes pertenecen a la entidad elegida en el selector.

---

## `salvia.entity_case`

Fuente: `src/internal/models/entity_case.go`

| Columna | Tipo | Nulo | Descripción |
|---|---|---|---|
| `id` | varchar(36) PK | No | UUID del vínculo sede–caso |
| `case_id` | varchar(36) | No | FK lógica → `victim_case.victim_case_i_code` |
| `entity_branch_id` | int64 | No | FK → `entity_branch.entity_branch_id` |
| `objetivo` | text | Sí | Texto de “Relación con entidad” (qué debe gestionar la víctima) |
| `last_action` | text | Sí | Última acción — **texto libre, no se calcula** del timeline |
| `created_by_id` | varchar(36) | Sí | Agente que creó la relación |
| `created_at` / `updated_at` | timestamptz | — | Auditoría |
| `deleted_at` | timestamptz | Sí | Soft-delete (GORM) |

Índice único parcial (activo): `(case_id, entity_branch_id) WHERE deleted_at IS NULL`.

---

## Tablas de apoyo

| Tabla | Uso en la pantalla |
|---|---|
| `entity_branch` | Sede del vínculo; aporta `entity_id`, nombre, sector vía join a `entity`, ciudad vía `town_code` |
| `entity` | Organización del filtro E07 / header (nombre + `entity_sector`) |
| `victim_case` | Via `case_id`: código visible, estado, nivel de riesgo |
| `victim_contact` | Via caso: nombre de la víctima, documento |
| `security.town` (+ city) | Ciudad mostrada en la fila (ubicación del caso o de la sede — **GAP** cuál priorizar) |

> `case_timeline_event` **no** alimenta “Última acción” en esta pantalla. Se usa `entity_case.last_action`.

---

## DTO / fila del listado (contrato UI)

| Campo UI | Origen |
|---|---|
| `entityCaseId` / `relId` | `entity_case.id` |
| `caseId` | `entity_case.case_id` |
| `entityBranchId` | `entity_case.entity_branch_id` |
| `entityBranchICode` | `entity_branch.entity_branch_i_code` (para “Ver detalle”) |
| `victimFullName` | `victim_contact` |
| `caseCode` | `victim_case` (código visible) |
| `document` | `victim_contact` |
| `city` | ciudad del caso o de la sede — **GAP** |
| `riskLevel` | nivel de riesgo del caso |
| `caseStatus` | estado del caso |
| `relationDescription` | `entity_case.objetivo` |
| `lastAction` | `entity_case.last_action` (string o null → “Sin acciones registradas”) |
| `sectorName` | `entity.entity_sector` (label del header, ej. “Justicia”) |

Shape de referencia ya existente (por caso, no por listado et): `EntityCaseWithRelations` en el mismo modelo — útil como base, pero el listado de esta pantalla necesita además datos del caso/víctima.

---

## Filtros y paginación (contrato API propuesto)

| Query / contexto | Evento | Notas |
|---|---|---|
| `entityId` | E07 (y luego E02–E04) | Obligatoria para listar. E01 no preselecciona (relación usuario↔entidad pendiente) |
| `document` | E02 | Filtro sobre documento de la víctima — **GAP** match parcial vs exacto |
| `city` | E03 | Filtro por ciudad — **GAP** fuente (caso vs sede) |
| `page` | E04 | 0-based |
| `pageSize` | E01 / E04 | Valor fijo — **GAP** |

Query base sugerida:

```
entity_case
  JOIN entity_branch ON entity_branch_id
  JOIN entity ON entity_branch.entity_id
  JOIN victim_case ON case_id
  … victim_contact / town …
WHERE entity_case.deleted_at IS NULL
  AND entity_branch.entity_id = :entidadActiva
  AND [filtros documento / ciudad]
ORDER BY entity_case.updated_at DESC  -- GAP: criterio de orden
```

---

## GAPS abiertos

| Decisión pendiente | Impacto |
|---|---|
| Relación usuario et ↔ entidad (hoy: selector con catálogo completo) | E01 / E07 |
| Ciudad del filtro/fila: ¿del caso o de la sede? | E03 + DTO |
| Endpoint catálogo de entities | E01 |
| Endpoint de listado paginado para rol `et` (`GET /api/v1/entity-cases` propuesto) | E02–E04, E07 |
| Destino “Ver caso” para rol `et` | E05 |
| “Ver detalle” → `/salvia/entidad/:entityBranchICode` (ruta planificada, aún no existe) | E06 |
| Mecanismo para poblar/editar `last_action` | Texto de última acción |
