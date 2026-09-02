# `case-task-modal` — Guía de uso

Componente modal para completar tareas de tipo `gestion_llamada`, `proyectar_oficio`, `comite_caso`, `Corregir oficio`, `justificar_remision`, `solicitud_entidad`, `respuesta_entidad` y `articulacion_territorial`. El padre lo controla mediante una referencia Vue (`ref`) y los métodos públicos `open(taskId)` o `openArticulacion(ctx)`.

---

## Props

El componente no recibe props. Se controla exclusivamente mediante el método `open()` y el evento `@completed`.

---

## Método público

| Método | Parámetros | Descripción |
|---|---|---|
| `open(taskId)` | `taskId: string` — UUID de la `case_task` | Carga la tarea del API, inicializa el formulario según su tipo y muestra el modal |
| `openArticulacion(ctx)` | `{ taskId?, barrierId, caseId, entityBranchId, entityName, barrierDescription }` | Mis Barreras: mismo formulario de Articulación territorial con o sin tarea previa |

---

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `@completed` | Cuando la tarea se completa exitosamente (status 200) | Objeto `case_task` actualizado con `status: 'Done'` y `form_data` |

---

## Qué hace el componente por sí solo

Al montarse, carga la lista de departamentos y ciudades desde el API para pre-poblar los dropdowns de los formularios `gestion_llamada` y `proyectar_oficio`. No muestra el modal hasta que el padre llame a `open()`.

Si la tarea abierta es de tipo `Corregir oficio`, además hace un `GET /api/v1/entity-letters/:letterId` (usando `tarea.entityLetterId`) para mostrar en modo lectura la razón de corrección y la ruta Kofax del oficio devuelto.

Si la tarea abierta es de tipo `solicitud_entidad`, además hace un `GET /api/v1/entity-cases/:entityCaseId/obligations?status=pendiente` (usando `tarea.entityCaseId`) para pintar el checklist de obligaciones pendientes y el nombre de la organización.

---

## Integración en el HTML de la pantalla

```html
{{ template "tasks/case_task_modal.html" . }}
```

---

## Tipos de tarea soportados

| `case_task.type` | Título del modal | Formulario que renderiza |
|---|---|---|
| `gestion_llamada` | "Gestión de Llamada" | Ubicación + Funcionario + Descripción + Switch ¿Genera oficio? (condicional: Asunto + Ruta Kofax) |
| `proyectar_oficio` | "Proyectar Oficio" | Ubicación + Funcionario + Asunto + Ruta Kofax |
| `comite_caso` | "Decisiones del Comité" | Checkbox group de decisiones + campos condicionales por decisión |
| `Corregir oficio` | "Corregir Oficio" | Sin formulario editable — vista de solo lectura con la razón de corrección y la ruta Kofax del `entity_letter` vinculado (`tarea.entityLetterId`). Solo pide confirmación. |
| `justificar_remision` | "Justificar Remisión" | Caja de solo lectura con `tarea.description` (motivo de devolución del psicólogo/a) + textarea "Tu respuesta" (required, maxlength 500) |
| `solicitud_entidad` | "Solicitudes a {entityName}" | Checklist de obligaciones pendientes (`entity_case_obligation`); cada ítem seleccionado exige una descripción (maxlength 500). Botón "Enviar solicitudes". |
| `respuesta_entidad` | "Respuesta institucional" | Caja de solo lectura (obligación + solicitud del agente) + radicado opcional + ampliación opcional + select de estado institucional required. Botón "Enviar respuesta". |
| `articulacion_territorial` | "Gestionar barrera" | Fecha de contacto, compromiso, seguimiento condicional, sugerencias opcionales. Botón "Guardar". POST `/api/v1/case-tasks/articulacion-territorial`. |

---

## Ejemplo de uso

```javascript
// En el template de la pantalla padre:
// <case-task-modal ref="taskModal" @completed="onTaskCompleted"></case-task-modal>

methods: {
  abrirTarea(taskId) {
    this.$refs.taskModal.open(taskId);
  },
  onTaskCompleted(task) {
    // Reflejar el nuevo estado en la UI o recargar datos
    this.cargarDatos();
  }
}
```

---

## Endpoint que consume

| Acción | Método | URL |
|---|---|---|
| Cargar tarea | GET | `/api/v1/case-tasks/:taskId` |
| Completar tarea | PUT | `/api/v1/case-tasks/:taskId/complete` |
| Guardar articulación territorial | POST | `/api/v1/case-tasks/articulacion-territorial` |
| Cargar oficio vinculado (solo `Corregir oficio`) | GET | `/api/v1/entity-letters/:letterId` |
| Cargar obligaciones pendientes (solo `solicitud_entidad`) | GET | `/api/v1/entity-cases/:entityCaseId/obligations?status=pendiente` |
| Cargar departamentos | GET | `/api/v1/locations/departments` (al montar) |
| Cargar ciudades | GET | `/api/v1/locations/cities` (al montar) |
| Cargar municipios | GET | `/api/v1/locations/towns?city_id={id}` (al seleccionar ciudad) |
| Cargar sedes de entidades | GET | `/api/v1/entity-branches?town_code={code}` (al seleccionar municipio) |

> **Corrección vs. versión anterior de este MD:** el prefijo real es `/api/v1/case-tasks/`, no `/api/v1/tasks/` — verificado en `src/frontend/js/components/case-task-modal.js` (métodos `open()` y `confirmar()`) y en los specs de `qa-salvia/tests/case-task-modal/`.
>
> ✅ **Resuelto — era una regresión, no una inconsistencia de diseño.** El endpoint `PUT /api/v1/case-tasks/:id/complete` + `CompleteWithFormData` existía desde el commit original (`34784f5`, "feat: Modal para manejar cada tipo de tarea"). Un commit posterior de otro desarrollador (`cfd1739`, un día después) lo eliminó accidentalmente en un merge — quitó la ruta `PUT`, el método del controller, el método del service y la dependencia `EntityLetterSvc` que necesitaba. Se restauró la lógica completa (controller + service + wiring en `main.go`) en esta sesión, idéntica a la del commit original. Compila limpio (`go build ./salvia/... ./internal/... .`).

### Payload de `PUT /api/v1/case-tasks/:taskId/complete`

```json
{
  "userId":   "icode del usuario en sesión",
  "formData": { }
}
```

#### `formData` para `gestion_llamada`
```json
{
  "departamentoId":     "string",
  "departamentoNombre": "string — label legible, resuelto de la lista de departamentos",
  "ciudadId":           "string",
  "ciudadNombre":       "string — label legible",
  "municipioId":        "string",
  "municipioNombre":    "string — label legible",
  "entidadId":          "integer | null (null si 'otra')",
  "entidadNombre":      "string",
  "funcionario":        "string",
  "descripcion":        "string | null",
  "generaOficio":       true,
  "asunto":             "string | null",
  "rutaKofax":          "string | null"
}
```

#### `formData` para `proyectar_oficio`
```json
{
  "departamentoId":     "string",
  "departamentoNombre": "string — label legible",
  "ciudadId":           "string",
  "ciudadNombre":       "string — label legible",
  "municipioId":        "string",
  "municipioNombre":    "string — label legible",
  "entidadId":          "integer | null (null si 'otra')",
  "entidadNombre":      "string",
  "funcionario":        "string",
  "asunto":             "string",
  "rutaKofax":          "string"
}
```

> **`*Nombre` — por qué existen:** además del ID (necesario porque el backend lo usa para poblar `entity_letter.department_id/city_id/town_id`), `_buildFormData()` ahora resuelve y guarda el label legible de cada `<select>` en el momento de completar la tarea. Esto permite que una vista de solo lectura (ej. historial de tareas completadas) muestre la ubicación sin tener que re-consultar los catálogos de ubicación ni depender de que un ID siga siendo válido más adelante.

#### `formData` para `comite_caso`
```json
{
  "decisiones":                   ["activar_enlace", "oficio", "recomendaciones_agente", "mecanismo_articulador"],
  "decisionesTexto":              ["Activar Enlace", "Generar Oficio", "Recomendaciones al Agente", "Mecanismo Articulador"],
  "observacionesOficio":          "string | null",
  "observacionesRecomendaciones": "string | null",
  "nivelMecanismo":               "municipal | departamental | nacional | null",
  "nivelMecanismoTexto":          "Municipal | Departamental | Nacional | null",
  "observacionesMecanismo":       "string | null"
}
```

> `decisionesTexto` y `nivelMecanismoTexto` son labels legibles añadidos para lectura (ej. historial de tareas completadas) — `decisiones` y `nivelMecanismo` se conservan tal cual porque el backend los usa para su lógica (`decisiones.includes('oficio')`, etc.).
>
> Nota de datos reales (registros anteriores a este cambio): `nivelMecanismo` se guardaba en minúscula (`"departamental"`). Con este cambio se mantiene así — el valor legible ahora vive aparte en `nivelMecanismoTexto`.

#### `formData` para `Corregir oficio`
```json
{}
```

Objeto vacío — no hay campos que enviar. El componente solo confirma la corrección; el backend usa `case_task.entity_letter_id` (no el payload) para saber qué `entity_letter` transicionar de `en_correccion` a `para_revisar`.

#### `formData` para `justificar_remision`
```json
{
  "respuesta": "string — justificación del agente, requerida, máx. 500 caracteres"
}
```

El motivo de devolución que el agente está respondiendo **no** viaja en el payload — ya está en `tarea.description` (lo dejó `sideEffectsValidateRemision`/`ValidateRemision` al crear esta `case_task`: `"Justificar remisión devuelta — Motivo: {motivo}"`). El backend, al completar esta tarea, reabre la `psychosocial_support` vinculada y crea una nueva `case_task` `validar_remision` — ver tabla "Efectos de lado por tipo de tarea" en `case-task-modal-interface.md`.

#### `formData` para `solicitud_entidad`
```json
{
  "entityName": "string — nombre de la organización",
  "solicitudes": [
    {
      "obligationId": "uuid de entity_case_obligation",
      "label": "string — texto del ítem",
      "descripcion": "string — requerida, máx. 500 caracteres"
    }
  ]
}
```

El backend crea una `case_task` `respuesta_entidad` por cada elemento de `solicitudes`, asignada al único usuario de `security.general_user` con la `entity_branch_id` de la sede. Si no hay exactamente un usuario, responde 400 y no completa la tarea original.

#### `formData` para `respuesta_entidad`
```json
{
  "solicitudTaskId": "uuid de la solicitud_entidad origen",
  "obligationId": "uuid de entity_case_obligation",
  "label": "string — obligación (heredado)",
  "descripcionSolicitud": "string — texto que escribió el agente (heredado)",
  "radicado": "string | null — opcional, máx. 80",
  "ampliacion": "string | null — opcional, máx. 500",
  "estadoInstitucional": "iniciado | en_proceso | concluido",
  "estadoInstitucionalTexto": "Iniciado | En proceso | Concluido"
}
```

El backend solo valida el estado, guarda el JSON y marca `Done`. No cambia `entity_case_obligation.status`.
