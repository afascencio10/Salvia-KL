# Plan — Tipificación en Barrera Detalle

> ⚠️ **MD temporal de planificación.** Validar antes de implementar código.
> Última actualización: 2026-06-24 — refleja el estado post-pull (95b47f6).

---

## Contexto

El requerimiento original tenía dos partes:

| Parte | Descripción | Estado |
|---|---|---|
| 1 | Guardar barreras correctamente al procesar el seguimiento: un registro `barrier_v2` por entrada del repeater, tipologías como CSV en `specific_barriers`, y crear `case_task` / `entity_letter` por cada acción de gestión seleccionada | ✅ **Implementado** en commit `95b47f6` |
| 2 | Mostrar la tipificación de la barrera y la entidad en la pantalla **Barrera Detalle** | 🔴 **Pendiente** |

Este MD cubre lo que falta.

---

## Aclaración: cómo funciona el guardado de barreras (PASO 3 del E04)

El código en `form_service.go → processFollowUpSubmission()` crea **UN registro `barrier_v2` por cada entrada del repeater** de la Sección 4. Las tipologías seleccionadas en Q2/Q5/Q8 (según el sector) quedan guardadas como CSV en el campo `specific_barriers`. El campo `specific_institutions` también es CSV.

Este es el comportamiento correcto. La descripción del flow MD (PASO 3.1) no lo refleja bien — ve el **Cambio 1** debajo.

---

## Cambio 1 — Corrección del flow MD (E04, PASO 3.1)

El flow `flow-E04-cuando-se-procesa-submission.md` en PASO 3.1 todavía describe la creación de barreras así:

```
Por cada opción del CSV:
  → DB.barrier_v2.Create({ sector, description: opcion_individual, status: "OPEN" })
```

Eso es incorrecto. El código real crea **un solo registro por entry**, con el CSV completo en `SpecificBarriers`. El flujo debe actualizarse:

```
PASO 3 — Por cada entry del repeater rgBarreras:

  3.1 Leer respuestas → entryMap
      sector            = entryMap[qBarreraSector]
      specificBarriers  = entryMap[sectorBarrierQ[sector]]   // CSV: "negativa_recibir_denuncia,tipificacion_erronea"
      specificInstitutions = entryMap[sectorInstitutionQ[sector]]  // CSV

  3.2 Crear UN registro barrier_v2:
      DB.barrier_v2.Create({
        case_id:               fu.case_id,
        follow_up_id:          fu.id,
        created_by_id:         actorId,
        status:                "OPEN",
        sector:                sector,
        specific_barriers:     specificBarriers,       // CSV completo
        specific_institutions: specificInstitutions,   // CSV completo
        other_barrier_desc:    entryMap[sectorOtherBarrierQ[sector]],
        institution_name:      entryMap[qBarreraOtraInstitucion],
        department_id, city_id, town_id, ...           // resto de campos comunes
        description:           entryMap[qBarreraDescripcion],
        management_actions:    entryMap[qBarreraGestion],
      })
      → newBarrierID = barrier_v2.id
      → barrierCount++

  3.3 Crear tareas y oficios por gestión (ya implementado — ver req-tareas-por-gestion-barrera.md)
      ...
```

**Archivo a modificar:** `DocsMD/Screens/hacer-seguimiento/Flujos/flow-E04-cuando-se-procesa-submission.md`

---

## Cambio 2 — Barrera Detalle: tipificación, entidad y carga real

### Estado actual de la pantalla

- Renderiza datos mock hardcoded en `data()`. No hay llamada a API.
- Muestra: sector, org (entidad), descripción.
- Falta: **tipificación** (labels de `specific_barriers`) y **entidad real** (labels de `specific_institutions` o `institution_name`).

### Cambios en la interfaz

**Tab Información General → `BarrierDetailRow`:**

```
ANTES:
  FieldSector:    barrera.sector
  FieldEntidad:   barrera.org          ← mock
  FieldDescripcion: barrera.description

DESPUÉS:
  FieldSector:       barrera.sector
  FieldTipificacion: barrera.tipificacion   ← [NUEVO] lista de labels de specific_barriers
  FieldEntidad:      barrera.institution    ← renombrar org → institution, dato real
  FieldDescripcion:  barrera.description
```

**Tipificación:** los values de `specific_barriers` (ej. `"negativa_recibir_denuncia,tipificacion_erronea"`) se resuelven a labels legibles en el backend antes de llegar al frontend. El frontend los recibe ya como lista de strings y los muestra.

**Entidad:** si `specific_institutions` tiene values → resolverlos a labels. Si `institution_name` tiene valor (sector `otras_instituciones`) → usarlo directamente.

### Cambios en E01 — Cuando carga la pantalla

**Comportamiento actual:** `mounted()` solo hace `document.getElementById('app').style.display = 'block'`.

**Comportamiento deseado:** `mounted()` dispara `GET /api/v1/barriers-v2/:barrierICode` y puebla los datos reales.

**Nuevo flujo E01:**

```
PASO 1 — document.getElementById('app').style.display = 'block'

PASO 2 — GET /api/v1/barriers-v2/:barrierICode

  SI error:
    → Mostrar estado de error
    → TERMINAR

  SI ok:
    → poblar barrera:
        sector          = data.sector
        tipificacion    = data.typologyLabels    // lista de strings ya resueltos
        institution     = data.institution       // label resuelto de specific_institutions / institution_name
        description     = data.description
        status          = data.status
        identifiedAt    = data.createdAt
        active          = data.status === 'OPEN'
        pendingTasks    = data.pendingTasks      // desde case_task WHERE barrier_id AND status != Done
        completedTasks  = data.completedTasks    // desde case_task WHERE barrier_id AND status == Done
        victimName      = data.victimName        // ver GAP-01

PASO 3 — Vue renderiza con datos reales
```

### Cambios en el backend

**¿Ya existe el endpoint?** No. El controlador `barrier_v2_gin_controller.go` solo tiene `GET /api/v1/barriers-v2` (lista filtrada por `createdById`). No hay endpoint para obtener una barrera por ID.

**¿Ya existe el servicio?** Sí. `BarrierV2Service.GetByID(ctx, id)` existe en `src/salvia/service/barrier_v2_service.go` (línea 33), pero retorna el modelo plano sin víctima ni labels resueltos.

**Lo que hay que agregar:**

| Qué | Dónde | Detalle |
|---|---|---|
| Endpoint `GET /api/v1/barriers-v2/:id` | `src/salvia/controller/barrier_v2_gin_controller.go` | Nuevo handler `GetByID` |
| Método enriquecido en el servicio | `src/salvia/service/barrier_v2_service.go` | `GetByIDWithDetails(ctx, id)` retorna struct con labels + víctima + tareas |
| Método en el repositorio (si hace falta) | `src/internal/repository/barrier_v2_repository.go` | Query `barrier_v2` JOIN `victim_case` + `case_task` WHERE `barrier_id` |
| Resolución de labels | Dentro del servicio | Ver GAP-02 |

**Shape de la respuesta del endpoint:**

```json
{
  "id": "uuid",
  "sector": "justicia",
  "typologyLabels": ["Negativa para recibir la denuncia", "Tipificación errónea del delito"],
  "institution": "Fiscalía General de la Nación",
  "description": "...",
  "status": "OPEN",
  "createdAt": "2026-06-01",
  "managementActions": "activacion_ruta_interinstitucional",
  "victimName": "...",
  "pendingTasks": [{ "id": "...", "label": "Activación de ruta interinstitucional" }],
  "completedTasks": []
}
```

### Cambios en BD

No se requieren columnas nuevas. Los campos `specific_barriers` e `specific_institutions` ya son `text` en `barrier_v2`. La resolución de labels se hace en el backend.

---

## Resumen ejecutivo de cambios pendientes

### Backend

| Cambio | Archivo | Urgencia |
|---|---|---|
| Endpoint `GET /api/v1/barriers-v2/:id` | `src/salvia/controller/barrier_v2_gin_controller.go` | Requerido |
| Método `GetByIDWithDetails` en servicio | `src/salvia/service/barrier_v2_service.go` | Requerido |
| Query enriquecida (barrier + case_task) | `src/internal/repository/barrier_v2_repository.go` | Requerido |
| Resolver `specific_barriers` → labels | En el servicio | Requerido — ver GAP-02 |

### Frontend

| Cambio | Archivo |
|---|---|
| E01: llamada a `GET /api/v1/barriers-v2/:barrierICode` en `mounted()` | `src/frontend/html/salvia/barriers/barrera_detalle.html` |
| Agregar campo `tipificacion` al `data()` de Vue | `barrera_detalle.html` |
| Renombrar `org` → `institution` en `data()` | `barrera_detalle.html` |
| Renderizar `tipificacion` como lista en la UI (tab Información General) | `barrera_detalle.html` |

### MDs a actualizar

| MD | Cambio |
|---|---|
| `Screens/hacer-seguimiento/Flujos/flow-E04-cuando-se-procesa-submission.md` | Reescribir PASO 3.1 para reflejar que crea UN registro por entry con CSV en `specific_barriers` |
| `Screens/Barrera Detalle/barrera-detalle-interface.md` | Agregar `FieldTipificacion`, renombrar `FieldEntidad` de `org` → `institution`, quitar nota de mock |
| `Screens/Barrera Detalle/barrera-detalle-index.md` | Actualizar descripción de E01 |
| `Screens/Barrera Detalle/Flujos/flow-E01-cuando-carga-pantalla.md` | Reemplazar flujo mock con el flujo real |

---

## Gaps y decisiones pendientes

| # | Gap / Decisión | Impacto |
|---|---|---|
| GAP-01 | **¿Qué datos de la víctima se muestran?** La UI mock muestra `victimName`, `age`, `location`. ¿El endpoint retorna los tres? `age` y `location` requieren JOINs adicionales (form2 → birthDate, geography). Definir si se incluyen o se simplifica la vista. | Complejidad del query de repositorio |
| GAP-02 | **Resolución de labels para tipificación e institución.** `specific_barriers` almacena values como `negativa_recibir_denuncia`. Opciones: (a) mapa hardcoded en Go igual al que ya existe en `form_service.go` para `gestionLabels` — fácil de implementar, requiere mantenimiento manual; (b) JOIN con tabla `option` WHERE `question_id` IN (qBarreraSalud, qBarreraJusticia, qBarreraProteccion). Definir cuál. | Complejidad del nuevo servicio |
| GAP-03 | **Prioridad en Barrera Detalle.** La UI mock muestra una etiqueta de prioridad (Alto/Medio/Bajo). El modelo `barrier_v2` no tiene campo `priority`. ¿Se agrega? ¿O se elimina de la UI? | Posible cambio en BD + AutoMigrate |
| GAP-04 | **Identificador en la URL.** La fachada Go pasa `barrierICode`. `barrier_v2.id` es el UUID (PK). Confirmar: ¿la URL `/salvia/barreras/:id` usa el UUID directamente como barrierICode? De ser así no se necesita columna adicional. | Sin cambio en BD si se usa el UUID |
