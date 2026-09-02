# Barreras Departamento — Tablas relacionadas

## Cadena de alcance

```
security.general_user
    │  general_user_assigned_department  → CommonSession.AssignedDepartmentID
    ▼
salvia.barrier_v2
  WHERE department_id = :AssignedDepartmentID
    AND deleted_at IS NULL
    [AND victim_case.victim_case_victim_doc_number ILIKE :docNumber]
    [AND entity_branch.entity_branch_name ILIKE :entidad]
    [AND b.city_id = :cityId]
    [AND status = :status]
```

El cliente no envía el departamento. Si `AssignedDepartmentID` está vacío, el API responde 403.

---

## `salvia.barrier_v2` (columnas usadas)

| Columna | Uso |
|---|---|
| `id` | PK; enlace a detalle |
| `case_id` | → `victim_case.victim_case_i_code` |
| `department_id` | Filtro de alcance (varchar; igual al id de sesión) |
| `city_id` | JOIN a `security.city` con `city_id::text = b.city_id` |
| `entity_branch_id` | Nombre de entidad en la columna BARRERA |
| `sector` | Badge de tipología |
| `description` | Texto de la barrera |
| `status` | `OPEN`, `En Gestion`, `Articulada`, `MANAGED` |
| `barrier_date` | Fecha mostrada |
| `created_at` | Orden del listado |
| `deleted_at` | Soft-delete; se excluyen |

---

## Tablas de apoyo

| Tabla | Uso |
|---|---|
| `security.general_user` | `general_user_assigned_department` del Enlace |
| `salvia.victim_case` | Nombre, apellidos, documento |
| `salvia.victim_case_form2` | Riesgo y `victim_case_form2_victim_phone` |
| `salvia.victim_contact_form1` | Teléfono de respaldo si form2 no tiene |
| `salvia.entity_branch` | `entity_branch_name` |
| `security.city` | Nombre en fila; catálogo del select vía `/api/v1/locations/cities` |

---

## Contrato API

`GET /api/v1/barriers-v2/department`

Query: `docNumber`, `entidad`, `cityId`, `status`, `page` (0-based), `pageSize` (default 10).

Auth: sesión + rol `en` + departamento asignado.

Respuesta:

```json
{
  "items": [],
  "total": 0,
  "page": 0,
  "pageSize": 10
}
```

Cada item sigue `BarrierV2DepartmentItem` (incluye `victimPhone`; el `caseCode` no se muestra en la columna CASO).
