# `case-entities` — Análisis de Base de Datos

> Este archivo responde dos preguntas concretas antes de implementar el componente:
> 1. ¿Existe una tabla que registre entidades ubicadas por departamento/ciudad/municipio?
> 2. ¿Existe una tabla que relacione un caso con una entidad?
>
> Y documenta qué hay que crear o modificar en la BD para soportar el componente.

---

## 1. Tablas existentes relevantes

### `salvia.entity` — la organización

Fuente: `src/salvia/dao/EntityDAO.go`

| Columna | Tipo | Notas |
|---|---|---|
| `entity_id` | uint (PK) | |
| `entity_i_code` | string | |
| `entity_name` | string | Requerido |
| `entity_description` | string | ⚠️ Bug existente: `EntityFieldDefinitions["EntityDescription"]` apunta por error a `DBName: "entity_name"` en vez de una columna propia — la descripción real de la entidad no tiene columna dedicada hoy. No es parte de este componente, pero queda anotado como hallazgo. |
| `entity_is_interoperable` | string(1) | Si la entidad tiene usuario/integración con Salvia |
| `entity_interoperability_code` | string | |
| `entity_response_time` | uint | |
| `entity_sector` | string(2) | Código de sector (mapea a Salud/Justicia/Protección, igual que `barrier_v2.Sector`) |

**No tiene ubicación** (ni department/city/town, ni dirección).

---

### `salvia.entity_branch` — la sede/sucursal (SÍ tiene ubicación)

Fuente: `src/salvia/dao/EntityBranchDAO.go`

| Columna | Tipo | Notas |
|---|---|---|
| `entity_branch_id` | uint (PK) | |
| `entity_branch_i_code` | string | |
| `entity_branch_name` | string | |
| `entity_branch_description` | string | |
| `entity_branch_address` | string | Dirección textual |
| `entity_branch_latitude` / `entity_branch_longitude` | float | Coordenadas |
| `entity_id` | uint (FK → `entity`) | A qué organización pertenece |
| `entity_branch_town_code` | string (FK → `security.town`) | **Aquí vive la ubicación** — resuelve a municipio → ciudad → departamento por el mismo mecanismo que ya usa `victim_case.town_code` |
| `entity_branch_source` | string(1) | `'u'` = creada por un usuario (se le antepone un marcador en el nombre al listarla) |

**Respuesta a la pregunta 1: sí existe.** `entity_branch` es la tabla de sedes ubicadas geográficamente (vía `town_code`, no columnas de depto/ciudad separadas — igual que el resto del sistema). `entity` es la organización "padre" sin ubicación propia. Un `entity` puede tener N `entity_branch`.

Ya existe además un endpoint para buscarlas por ubicación + sector, pensado originalmente para el flujo de oficios:

```
GET /api/v1/entity-branches?town_code={townCode}&sector={sectorDeBarrera}
```
(`src/salvia/controller/entity_branch_api_controller.go` — usado hoy en "el dropdown de entidades en el modal proyectar" de oficios)

> Nota aparte: existe también `salvia.directories` (`src/internal/models/directory.go`) — catálogo de fiscalías/comisarías/urgencias/líneas de emergencia para la app Flutter, ubicado solo por `city_id`. Es un catálogo distinto, de menor detalle, pensado para consulta pública desde la app móvil — no es la misma tabla que `entity`/`entity_branch` y no aplica a este componente.

---

### `salvia.entity_letter` (oficios) — SÍ referencia una entidad, pero solo en oficios ya proyectados

Fuente: `src/internal/models/entity_letter.go`

Hallazgo relevante que no estaba documentado en `case-oficios-interface.md`: el modelo real ya tiene bastantes más columnas que las que ese MD describía, incluyendo:

| Columna | Tipo | Notas |
|---|---|---|
| `case_id` | varchar(36) not null | |
| `barrier_id` | varchar(36) not null | |
| `entity_branch_id` | int64 nullable | **FK a `entity_branch`** |
| `department_id` / `city_id` / `town_id` | varchar nullable | Ubicación del oficio (puede no coincidir 1:1 con la de la sede si se corrige a mano) |
| `entidad` | varchar(255) nullable | Nombre de la entidad en texto libre (redundante con `entity_branch_id`, se usa como fallback/legacy) |

Estos campos (`entity_branch_id`, `department_id`, `city_id`, `town_id`) se llenan en la acción **"proyectar"** (`PUT /api/v1/entity-letters/:id/action`, `action: "proyectar"`) — es decir, **un oficio recién creado (`por_proyectar`) todavía no tiene entidad asociada.** Solo los oficios que ya avanzaron a `para_revisar` o un estado posterior tienen `entity_branch_id` poblado.

---

### `salvia.barrier_v2` — hoy no referencia ninguna entidad → se agrega FK opcional

Fuente: `src/internal/models/barrier_v2.go`

Tiene `SpecificInstitutions` (CSV de códigos de catálogo hardcodeado) e `InstitutionName` (texto libre, solo cuando se elige "otras_instituciones") — ambos son **texto**, no FK a `entity`/`entity_branch`. También tiene su propia ubicación (`DepartmentID`/`CityID`/`TownID`), independiente de cualquier entidad.

**Decisión de producto:** se agrega `EntityBranchID *int64` (nullable) al struct `BarrierV2`, columna `entity_branch_id`. Esto resuelve el GAP de "barreras activas por entidad" (ver sección 4) y habilita el filtro "¿Con barreras activas?" del componente. Al ser nullable, no rompe barreras existentes ni obliga a asociar entidad en el formulario de registro de barrera — queda como campo opcional a llenar cuando se conozca la sede específica.

### `salvia.case_task` — tampoco referencia entidad

Fuente: `src/internal/models/case_task.go`. Se relaciona con `barrier_id`, `entity_letter_id`, `emergency_measure_id`, `psychosocial_support_id`, `economic_stabilization_id` — ninguno es `entity_branch_id`.

### `salvia.case_owner` + `salvia.rel_case_owner_victim_case` (legacy) — no es lo que se busca

Fuente: `src/salvia/dao/CaseOwnerDAO.go`, `RelCaseOwnerVictimCaseDAO.go`. `case_owner` sí tiene `entity_branch_id`, pero representa **el profesional responsable** (`case_owner_general_user`) vinculado a una sede — y `rel_case_owner_victim_case` vincula ese profesional a **un** caso (relación 1 a 1, marcada legacy en `tables-database.md`). No es una lista de entidades relacionadas con un caso; es "quién es el dueño del caso", un concepto distinto.

---

## 2. Respuesta a la pregunta 2

**No existe una tabla que relacione directamente un caso (`victim_case`) con una o varias entidades (`entity`/`entity_branch`).**

Las únicas señales indirectas hoy son:
- `entity_letter.entity_branch_id` — pero solo para oficios ya proyectados (deja fuera barreras sin oficio, y oficios en `por_proyectar`).
- `case_owner` / `rel_case_owner_victim_case` — relación 1:1 legacy de responsable, no de "entidades con las que trabaja el caso".

Esto confirma lo que ya sugería el mockup React (`SeccionEntidades.jsx`, botón "+ Agregar entidad"): hoy no hay forma automática de saber "con qué entidades está trabajando este caso" — hay que registrarlo explícitamente.

---

## 3. Propuesta de cambio en BD

### Nueva tabla: `salvia.entity_case`

Nombre confirmado por el usuario: `entity_case` (no se sigue la convención `rel_<entidad1>_<entidad2>` de las tablas legacy — se prefiere un nombre corto de dominio, igual que `case_task` o `case_owner`). Implementada con el patrón **moderno** (GORM + uuid + soft delete) que ya usan las tablas más recientes del dominio (`barrier_v2`, `entity_letter`, `case_task`) en vez del patrón DAO legacy (id numérico autoincremental) que usan `entity`/`entity_branch`/`case_owner`.

```go
// src/internal/models/entity_case.go
package models

import (
    "time"
    "gorm.io/gorm"
)

// EntityCase relaciona una sede de entidad (entity_branch) con un caso.
// Representa "esta víctima/caso está gestionando algo con esta entidad" — se crea
// manualmente por un agente, no se infiere automáticamente de barreras u oficios.
type EntityCase struct {
    ID             string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
    CaseID         string         `gorm:"type:varchar(36);not null;index;column:case_id" json:"caseId"`
    EntityBranchID int64          `gorm:"not null;index;column:entity_branch_id" json:"entityBranchId"`
    Objetivo       *string        `gorm:"type:text;column:objetivo" json:"objetivo,omitempty"`
    LastAction     *string        `gorm:"type:text;column:last_action" json:"lastAction,omitempty"`
    CreatedByID    string         `gorm:"type:varchar(36);column:created_by_id" json:"createdById"`
    CreatedAt      time.Time      `json:"createdAt"`
    UpdatedAt      time.Time      `json:"updatedAt"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityCase) TableName() string {
    return "salvia.entity_case"
}
```

| Campo | Origen |
|---|---|
| `case_id` | `victim_case.victim_case_i_code` (mismo string usado en todo el resto del sistema, ej. `entity_letter.case_id`) |
| `entity_branch_id` | Elegido por el agente en el modal "+ Agregar entidad", vía `GET /api/v1/entity-branches?town_code=&sector=` (endpoint ya existente) |
| `objetivo` | Texto libre — "¿qué debe gestionar la víctima en esta entidad?" (mismo campo que ya proponía el mockup) |
| `last_action` | **Campo string editable, no calculado.** Guarda la última acción con esta entidad tal como la registre un agente. Pendiente definir en un siguiente ciclo *cómo y cuándo* se le da valor (¿un campo de texto libre en un formulario? ¿se actualiza automáticamente en algún paso del flujo de oficios?) — ver GAP en sección 4. |
| `created_by_id` | Agente que registró la relación (auditoría) |

### Pasos de implementación

1. Crear el archivo de modelo de arriba en `src/internal/models/`.
2. Agregar `EntityBranchID *int64` (nullable, columna `entity_branch_id`) al struct `BarrierV2` en `src/internal/models/barrier_v2.go`.
3. Agregar `&models.EntityCase{}` y (si no estaba) `&models.BarrierV2{}` sigue en la lista de `AutoMigrate` en `src/main.go` (línea ~85) y en `src/cmd/migrate/main.go`. GORM AutoMigrate crea la tabla y agrega la columna nueva solo — `salvia_gorm` ya es owner del schema (ver `DocsMD/General/tables-database.md`).
4. **Único paso manual necesario:** crear el índice único parcial (GORM no lo genera desde tags de struct):
   ```sql
   CREATE UNIQUE INDEX ux_entity_case_active
     ON salvia.entity_case (case_id, entity_branch_id)
     WHERE deleted_at IS NULL;
   ```
   Evita duplicar la misma entidad activa en el mismo caso. Ejecutar con el usuario `salvia_gorm` en el puerto `6543` (session mode, requerido para DDL — ver `tables-database.md`).

No se requiere modificar `entity` ni `entity_branch` — se reutilizan tal cual.

---

## 4. Qué se puede derivar automáticamente vs. qué queda como GAP

| Dato que pide el mockup | ¿Se puede obtener hoy? | Cómo |
|---|---|---|
| Nombre, sector, dirección, ubicación de la entidad | ✅ Sí | Join `entity_branch` + `entity` + `town`/`city`/`department` (mismo patrón que `victim_case`) |
| Oficios con esta entidad | ✅ Sí | `COUNT(entity_letter) WHERE case_id = ? AND entity_branch_id = ?` — confiable, ya que `entity_branch_id` existe en `entity_letter` |
| "Última acción" con la entidad | ✅ Sí (decisión de producto) | Ya **no** se calcula. Es el campo `entity_case.last_action` (string), editado directamente — ver sección 3. Pendiente definir en un siguiente ciclo el mecanismo exacto para poblarlo. |
| Barreras activas de esta entidad | ✅ Resuelto | Se agrega `barrier_v2.entity_branch_id` (nullable) — ver sección 3. Con esto: `COUNT(barrier_v2) WHERE case_id = ? AND entity_branch_id = ? AND status IN ('OPEN', 'En Gestion')`. Habilita el filtro "¿Con barreras activas?" del componente (deja de estar bloqueado). |

### Barreras activas por entidad — GAP resuelto

Se adopta la opción de agregar `entity_branch_id` nullable a `barrier_v2` (antes presentada como alternativa 2). Queda pendiente, fuera del alcance de este componente, el cambio en el formulario de registro de barrera (`hacer_seguimiento.html` / `follow_up_v2`) para que el agente pueda elegir la sede al registrar la barrera — sin ese cambio, el campo existirá en la BD pero se poblará en 0 registros hasta que se implemente esa parte. El componente `case-entities` puede construirse ya asumiendo que el campo existe (cuenta 0 mientras no se registre ninguna barrera con `entity_branch_id`).
