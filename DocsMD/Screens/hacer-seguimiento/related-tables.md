# Tablas relacionadas — `hacer-seguimiento`

Tablas de BD que lee o escribe la pantalla `/salvia/hacer-seguimiento/:id` y su goroutine `processFollowUpSubmission`.

---

## `salvia.follow_up_v2`

Seguimiento programado. La pantalla lo carga al iniciar (E-01) y `processFollowUpSubmission` lo marca como REALIZADO.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso VBG asociado |
| `form_submission_id` | uuid | Sí | Submission del formulario dinámico vinculado |
| `status` | varchar | No | `PENDIENTE` / `REALIZADO` / `NO_REALIZADO` |
| `scheduled_date` | timestamptz | Sí | Fecha programada del seguimiento |
| `completed_at` | timestamptz | Sí | Fecha en que fue marcado REALIZADO |
| `sequence` | int | Sí | Número de seguimiento en el caso |

---

## `salvia.answer`

Respuestas individuales del formulario. `processFollowUpSubmission` las lee para construir `answerMap`.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `form_submission_id` | uuid | No | Submission al que pertenece |
| `question_id` | uuid | No | Pregunta respondida |
| `repeater_entry_id` | uuid | Sí | Null si es respuesta directa |
| `value` | text | Sí | Valor como string (CSV para multi-select) |

---

## `salvia.barrier_v2`

Barreras institucionales identificadas en el seguimiento. Creadas en PASO 3 de `processFollowUpSubmission`.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | No | Seguimiento que la generó |
| `sector` | varchar | No | `salud` / `justicia` / `proteccion` |
| `description` | text | No | Opción seleccionada por el profesional |
| `status` | varchar | No | `OPEN` / `CLOSED` |

---

## `salvia.case_task`

Tareas generadas por `processFollowUpSubmission`: gestión de barreras (PASO 3.2, `category: "Barreras"`) y validación de remisión psicosocial (PASO 4, `category: "Psicosocial"`, `type: "validar_remision"`, creada sin asignar cuando se deriva a `atencion_psico`).

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | Sí | Seguimiento que la generó |
| `category` | varchar | No | `Barreras` / `Psicosocial` |
| `type` | varchar | No | `gestion_llamada` / `proyectar_oficio` / `comite_caso` / `validar_remision` |
| `description` | text | No | Texto legible de la tarea |
| `assigned_user_id` | varchar | No | ICode del agente asignado — vacío (`''`) si queda sin asignar |
| `status` | varchar | No | `ToDo` / `Done` |
| `barrier_id` | uuid | Sí | Barrera asociada (solo tareas de gestión de barreras) |
| `entity_letter_id` | uuid | Sí | Oficio asociado (solo si la gestión de la barrera generó uno) |
| `psychosocial_support_id` | uuid | Sí | Remisión psicosocial asociada (solo tarea `validar_remision`) |

---

## `salvia.psychosocial_support`

Remisión al equipo de Atención Psicosocial. Creada solo si se cumplen criterios y no hay exclusión con medidas de emergencia.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | No | Seguimiento que la generó |
| `type` | varchar | No | `derivacion` |
| `status` | varchar | No | `ACTIVE` |

---

## `salvia.men_team_remision`

Remisión al equipo de Atención Hombres. Creada solo si `criterio_hombres` está seleccionado.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | No | Seguimiento que la generó |
| `type` | varchar | No | `derivacion` |
| `status` | varchar | No | `ACTIVE` |

---

## `salvia.discapacidad_remision`

Remisión al equipo de Discapacidad. Se crea un registro por cada servicio seleccionado.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | No | Seguimiento que la generó |
| `service` | varchar | No | `apoyo_lsc` / `enfoque_discapacidad` |
| `status` | varchar | No | `ACTIVE` |
| `notes` | text | Sí | Notas opcionales |

---

## `salvia.economic_stabilization`

Remisión al equipo de Estabilización Económica. Creada solo si se seleccionó al menos 1 criterio.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | No | Seguimiento que la generó |
| `type` | varchar | No | `derivacion` |
| `status` | varchar | No | `ACTIVE` |

---

## `salvia.emergency_measure`

Medida de emergencia adoptada para la víctima. Se crea un registro por cada medida seleccionada.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | No | Seguimiento que la generó |
| `type` | varchar | No | `alojamiento` / `transporte` / `alimentacion` / `vestuario` / `apoyo_psico` / `otras_me` |
| `status` | varchar | No | `ACTIVE` |

---

## `salvia.case_timeline_event`

Eventos del timeline del caso. `processFollowUpSubmission` crea 1 evento al completar el seguimiento y 1 adicional si el caso es editado después de REALIZADO.

| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| `id` | uuid | No | PK |
| `case_id` | varchar | No | ICode del caso |
| `follow_up_id` | uuid | Sí | Seguimiento asociado |
| `category` | varchar | No | `Seguimientos` |
| `type` | varchar | No | `Seguimiento Ejecutado` / `Seguimiento Editado` |
| `icon` | varchar | Sí | Icono visual del evento |
| `color` | varchar | Sí | Color hex del evento |
| `description` | text | Sí | Texto descriptivo del evento |
| `event_user_id` | varchar | Sí | ICode del agente que ejecutó la acción |
| `actor_name` | varchar | Sí | Nombre resuelto del agente |
