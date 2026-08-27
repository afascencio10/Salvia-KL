# `reasignar-remisiones-modal` — Interfaz del Modal

Modal para reasignar una o varias remisiones de Atención Psicosocial seleccionadas desde `remisiones-psicosocial-component`. Lo consume la **pantalla padre** (p. ej. Historial de Remisiones, rol `sv`). Se abre cuando el padre recibe el evento `reasignar-remisiones` (E-15) y llama al método público `open(remisiones)`.

Permite elegir asignación **individual** (profesional `ps` o `ts`) o **en dupla**, y al confirmar actualiza `psychosocial_support`, los `team_contact` pendientes vinculados y las `case_task` pendientes (`ToDo`) de esas remisiones.

## Archivos relevantes (implementación futura)

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/reasignar-remisiones-modal.js` | Componente Vue |
| `src/frontend/css/reasignar-remisiones-modal.css` | Estilos (prefijo `rrm-`) |
| `src/frontend/html/salvia/remisiones-psicosocial/reasignar_remisiones_modal.html` | Partial HTML con `x-template` |

## Relación con otros componentes

```
historial_remisiones.html (padre — rol sv)
│
├── remisiones-psicosocial-component
│   :reasignacion="true"
│   └── emite reasignar-remisiones { remisiones }     → E-15
│
└── reasignar-remisiones-modal
    ├── ref="reasignarModal"
    ├── open(remisiones)  ← padre llama tras E-15
    └── eventos RRM-01 … RRM-05
```

---

## Árbol de interfaz

Referencia visual: mockups del modal (Image 1 — switch off; Image 2 — select profesional; Image 3 — select dupla).

```
reasignar-remisiones-modal
│
├── [v-if !visible]  → no renderiza nada
│
└── [v-if visible]
    ModalBackdrop  (.rrm-backdrop)
    │  @click.self → close()   // RRM-04
    │
    └── ModalDialog  (.rrm-dialog)  role="dialog" aria-modal="true"
        │
        ├── ModalHeader  (.rrm-header)
        │   ├── Title  (.rrm-title)  "Reasignar remisiones"
        │   └── Subtitle  (.rrm-subtitle)
        │       "{ remisiones.length } remisión(es) seleccionada(s)."
        │
        ├── ModalBody  (.rrm-body)
        │   │
        │   ├── DuplaSwitchRow  (.rrm-switch-row)
        │   │   ├── ToggleSwitch  v-model="asignarEnDupla"
        │   │   │   @change → onToggleDupla()   // RRM-02
        │   │   └── Label  "Asignar en dupla"
        │   │
        │   ├── [v-if !asignarEnDupla]  — modo profesional (Image 2)
        │   │   AssignmentSection
        │   │   ├── Label  "Profesional"
        │   │   ├── [v-if loadingOptions]  Spinner "Cargando..."
        │   │   ├── [v-else-if optionsError]  Mensaje error
        │   │   └── <select>  v-model="selectedProfessionalId"
        │   │       ├── <option value="">Seleccionar profesional...</option>
        │   │       ├── <optgroup label="Psicólogas">   ← rol ps
        │   │       │   └── <option :value="p.icode"> p.fullName
        │   │       └── <optgroup label="Trabajadoras Sociales">  ← rol ts
        │   │           └── <option :value="p.icode"> p.fullName
        │   │
        │   └── [v-if asignarEnDupla]  — modo dupla (Image 3)
        │       AssignmentSection
        │       ├── Label  "Dupla"
        │       ├── [v-if loadingOptions]  Spinner
        │       ├── [v-else-if optionsError]  Mensaje error
        │       └── <select>  v-model="selectedDuplaId"
        │           ├── <option value="">Seleccionar dupla...</option>
        │           └── <option :value="d.id"> d.label
        │               // "Dupla 1 — Alejandra Mora + Valentina Ospina"
        │
        └── ModalFooter  (.rrm-footer)
            ├── CancelBtn  (.rrm-btn--secondary)  "Cancelar"
            │   → close()   // RRM-04
            └── ReassignBtn  (.rrm-btn--primary)  "Reasignar"
                :disabled si !canConfirm || loadingOptions || saving
                → confirmReassign()   // RRM-05
```

---

## Método público de apertura

| Método | Parámetros | Descripción |
|---|---|---|
| `open(remisiones)` | `Array<PsychosocialListItem>` | Abre el modal con las remisiones de E-15. Dispara RRM-01. |

```javascript
onReasignarRemisiones: function(payload) {
    this.$refs.reasignarModal.open(payload.remisiones);
}
```

---

## Estado interno del componente

| Variable | Tipo | Inicial | Descripción |
|---|---|---|---|
| `visible` | `Boolean` | `false` | Modal abierto/cerrado |
| `remisiones` | `Array` | `[]` | Remisiones a reasignar (copia del payload E-15) |
| `asignarEnDupla` | `Boolean` | `false` | Switch "Asignar en dupla" |
| `selectedProfessionalId` | `String` | `''` | `general_user_i_code` del profesional elegido |
| `selectedDuplaId` | `String` | `''` | `dupla.id` elegida |
| `professionalGroups` | `Array` | `[]` | Grupos ps / ts para `<optgroup>` |
| `duplaOptions` | `Array` | `[]` | Duplas con label enriquecido |
| `loadingOptions` | `Boolean` | `false` | Carga de select en curso |
| `optionsError` | `String \| null` | `null` | Error al cargar opciones |
| `saving` | `Boolean` | `false` | POST reasignación en curso |
| `saveError` | `String \| null` | `null` | Error al guardar |

### Computed `canConfirm`

```javascript
canConfirm = asignarEnDupla
  ? selectedDuplaId !== ''
  : selectedProfessionalId !== ''
```

---

## Modos de asignación (mutuamente excluyentes)

| Switch | Campo en `psychosocial_support` | Campo en `team_contact` (pendientes) |
|---|---|---|
| OFF — Profesional | `professional_id = icode`, `dupla_id = NULL` | `professional_id = icode`, `dupla_id = NULL` |
| ON — Dupla | `dupla_id = id`, `professional_id = NULL` | `dupla_id = id`, `professional_id = NULL` |

> La remisión se asigna **o** a un profesional **o** a una dupla, nunca ambos (coherente con columna ESTADO Y ASIGNACIÓN del listado).

---

## Endpoints backend

### GET `/api/v1/psychosocial-support/profesionales-reasignacion`

Usuarios activos con roles `ps` o `ts`, agrupados para el `<select>` (Image 2).

**Criterios SQL:**
- `general_user_status = 'e'`
- Rol en `security.role.role_code IN ('ps', 'ts')`
- JOIN `general_user_profile` para nombre

**Response 200:**

```json
{
  "groups": [
    {
      "role": "ps",
      "label": "Psicólogas",
      "professionals": [
        { "icode": "...", "fullName": "Alejandra Mora" }
      ]
    },
    {
      "role": "ts",
      "label": "Trabajadoras Sociales",
      "professionals": [
        { "icode": "...", "fullName": "Valentina Ospina" }
      ]
    }
  ]
}
```

### GET `/api/v1/duplas/reasignacion`

Duplas activas con nombres de integrantes para el `<select>` (Image 3).

**Response 200:**

```json
{
  "duplas": [
    {
      "id": "uuid-dupla",
      "name": "Dupla 1",
      "label": "Dupla 1 — Alejandra Mora + Valentina Ospina",
      "psychologistName": "Alejandra Mora",
      "socialWorkerName": "Valentina Ospina"
    }
  ]
}
```

`label` = `{name} — {psychologistName} + {socialWorkerName}`

> Reutilizable ampliando el endpoint existente `GET /api/v1/duplas` si se prefiere un solo contrato con query `?enriched=true`.

### POST `/api/v1/psychosocial-support/reasignar-bulk`

Persiste la reasignación (RRM-05).

**Request:**

```json
{
  "remision_ids": ["uuid-1", "uuid-2"],
  "assign_mode": "professional",
  "professional_id": "icode-profesional",
  "dupla_id": null
}
```

```json
{
  "remision_ids": ["uuid-1"],
  "assign_mode": "dupla",
  "professional_id": null,
  "dupla_id": "uuid-dupla"
}
```

**Response 200:**

```json
{
  "ok": true,
  "remisiones_updated": 2,
  "team_contacts_updated": 5,
  "case_tasks_updated": 3
}
```

**Permiso sugerido:** `reassign_remisiones_psicosocial` → rol `sv` (misma pantalla que Historial de Remisiones).

---

## Reglas de negocio al confirmar (RRM-05)

Por cada `remision_id` en el lote:

### 1. Actualizar `salvia.psychosocial_support`

**Modo profesional:**
```sql
UPDATE salvia.psychosocial_support
SET professional_id = $professionalId,
    dupla_id        = NULL,
    updated_at      = NOW()
WHERE id = $remisionId
  AND deleted_at IS NULL
```

**Modo dupla:**
```sql
UPDATE salvia.psychosocial_support
SET dupla_id        = $duplaId,
    professional_id = NULL,
    updated_at      = NOW()
WHERE id = $remisionId
  AND deleted_at IS NULL
```

### 2. Actualizar `salvia.team_contact` (solo pendientes)

Solo registros **no completados** de esa remisión:

```sql
UPDATE salvia.team_contact
SET professional_id = $professionalIdOrNull,
    dupla_id        = $duplaIdOrNull,
    updated_at      = NOW()
WHERE psicosocial_id = $remisionId::text
  AND is_completed = false
  AND deleted_at IS NULL
```

> No modificar `team_contact` con `is_completed = true` (sesiones ya realizadas conservan asignación histórica).

### 3. Reasignar `salvia.case_task` (solo pendientes)

Solo tareas **no realizadas** (`status = 'ToDo'`) vinculadas a esa remisión vía `psychosocial_support_id`.

`$taskAssigneeId`:
- Modo profesional → `professional_id` del request
- Modo dupla → `dupla.psychologist_id` (la TS de la dupla no recibe las tareas)

```sql
UPDATE salvia.case_task
SET assigned_user_id = $taskAssigneeId,
    updated_at       = NOW()
WHERE BTRIM(psychosocial_support_id::text) = BTRIM($remisionId)
  AND status = 'ToDo'
  AND deleted_at IS NULL
```

> No modificar `case_task` con `status = 'Done'` (tareas ya realizadas conservan asignación histórica).

### Transacción

- Una transacción DB para todo el lote.
- Si un `remision_id` no existe → omitir (`skipped++`) o abortar según decisión de implementación (recomendado: abortar y rollback).

---

## Eventos del modal

| ID | Evento | Documento |
|---|---|---|
| RRM-01 | Cuando abre el modal | [flow-RRM01](./Flujos/flow-RRM01-cuando-abre-modal-reasignacion.md) |
| RRM-02 | Cuando alterna "Asignar en dupla" | [flow-RRM02](./Flujos/flow-RRM02-cuando-alterna-asignar-en-dupla.md) |
| RRM-03 | Cuando carga opciones del select | [flow-RRM03](./Flujos/flow-RRM03-cuando-carga-opciones-asignacion.md) |
| RRM-04 | Cuando cancela o cierra | [flow-RRM04](./Flujos/flow-RRM04-cuando-cancela-modal.md) |
| RRM-05 | Cuando confirma reasignación | [flow-RRM05](./Flujos/flow-RRM05-cuando-confirma-reasignacion.md) |

Inventario: [reasignar-remisiones-modal-events.md](./reasignar-remisiones-modal-events.md)

---

## Eventos emitidos hacia el padre

| Evento | Cuándo | Payload |
|---|---|---|
| `closed` | Cierre sin éxito (RRM-04) | `{}` |
| `reassigned` | Éxito (RRM-05) | `{ remisiones, assignMode, updated }` |

El padre debe llamar `this.$refs.remisionesComponent.reload()` para refrescar tabla y limpiar selección.

---

## Precondiciones heredadas de E-14

- Ninguna remisión seleccionada tiene `status = 'cerrado'`.
- Selección solo de la página actual (no cross-page).
- Pantalla con `:reasignacion="true"`.

El modal confía en el payload de E-15 (sin revalidar estados mezclados; sí excluye cerradas en backend).
