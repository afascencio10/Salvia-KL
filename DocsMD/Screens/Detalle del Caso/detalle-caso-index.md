# Detalle del Caso — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [detalle-caso-interface.md](detalle-caso-interface.md) |

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando cambia de tab | User Interaction | — |
| E03 | Cuando presiona "📋 Info caso" | User Interaction | — |
| E04 | Cuando presiona "Reasignar" (caso) | User Interaction | — |
| E05 | Cuando presiona "👤 Asignarme este caso" | User Interaction | — |
| E06 | Cuando presiona "+ Nuevo seguimiento" | User Interaction | — |
| E07 | Cuando presiona "Ver →" en banner de próximo seguimiento | User Interaction | — |
| E08 | Cuando presiona "Ver tareas" en banner de tareas pendientes | User Interaction | — |
| E09 | Cuando guarda nuevo seguimiento | User Interaction | — |
| E10 | Cuando confirma reasignación de caso | User Interaction | — |
| E11 | Cuando expande/colapsa detalle de una barrera | User Interaction | — |
| E12 | Cuando presiona "▶ Iniciar" seguimiento | User Interaction | — |
| E13 | Cuando presiona "✏️ Editar" seguimiento | User Interaction | — |
| E14 | Cuando confirma edición de seguimiento | User Interaction | — |
| E15 | Cuando presiona "📋 Reasignar" seguimiento | User Interaction | — |
| E16 | Cuando confirma reasignación de seguimiento | User Interaction | — |
| E17 | Cuando expande/colapsa historial de intentos | User Interaction | — |
| E18 | Cuando presiona "🔄 Reintentar" | User Interaction | — |
| E19 | Cuando el modal de contacto completa una gestión | User Interaction | — |
| E20 | Cuando el modal de completar tarea emite "completed" | User Interaction | — |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

**Evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** Dispara dos llamadas paralelas: carga el detalle completo del caso (datos, seguimientos, barreras, derivaciones, timeline) y carga la lista de todos los agentes del equipo. Al completar la carga del detalle, dispara además la carga de tareas del caso (`cargarTareasCaso`) y, en paralelo, operadores y agentes RO filtrados por team.
**Requerido:** Sí

---

**Evento:** Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** Actualiza `tabActiva` con el id del tab presionado. Vue re-renderiza el contenido del panel mostrando la sección correspondiente (info, derivaciones, barreras, seguimientos, tareas, timeline).
**Requerido:** Sí

---

**Evento:** Cuando presiona "📋 Info caso"
**Tipo:** User Interaction
**Descripción:** Pone `mostrarInfoCaso = true`, lo que hace visible el componente `CaseInfo` en modo modal con el detalle completo del caso.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Reasignar" (caso)
**Tipo:** User Interaction
**Descripción:** Limpia `reasignarOperador` y pone `modalReasignar = true`. Visible únicamente para rol `sv`. El dropdown muestra operadores filtrados por el team del caso.
**Requerido:** Sí

---

**Evento:** Cuando presiona "👤 Asignarme este caso"
**Tipo:** User Interaction
**Descripción:** Visible solo para rol `op` cuando el caso no tiene agente asignado o el agente no es el usuario actual. POST `/api/v1/casos/:id/reasignar` con `operador_icode = userICode` (sin modal intermedio — es una acción directa). El backend actualiza `agent_id`/`victim_case_team`, reasigna al usuario los seguimientos `PENDIENTE` y las `CaseTask` con estado `ToDo`, y actualiza `victim_case_owner_description`. En paralelo registra evento `REASIGNACION` en el timeline. Recarga datos al finalizar.
**Requerido:** Sí

---

**Evento:** Cuando presiona "+ Nuevo seguimiento"
**Tipo:** User Interaction
**Descripción:** Abre el modal de nuevo seguimiento. Precalcula un agente por defecto: si el rol es `ro`, se preselecciona a sí mismo; si no, y ya existen seguimientos, preselecciona el agente del último seguimiento registrado. Resetea el formulario, los flags de validación y el error del modal. Deshabilitado si `followUpsV2.length >= 8`.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Ver →" en banner de próximo seguimiento
**Tipo:** User Interaction
**Descripción:** Asigna `tabActiva = 'seguimientos'` para llevar al usuario directamente a la pestaña de seguimientos.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Ver tareas" en banner de tareas pendientes
**Tipo:** User Interaction
**Descripción:** Asigna `tabActiva = 'tareas'`. Existen dos banners con el mismo trigger y la misma acción: uno global (visible sobre cualquier tab, arriba de los tabs) y otro repetido dentro del tab Gestión institucional — ambos se basan en `caseTasks` con `status === 'ToDo'`.
**Requerido:** Sí

---

**Evento:** Cuando guarda nuevo seguimiento
**Tipo:** User Interaction
**Descripción:** Valida fecha (requerida, no pasada) y agente (requerido solo para `sv`; para otros roles el campo va fijo/deshabilitado). Si válido y no se alcanzó el máximo de 8 seguimientos, POST `/api/v1/casos/:id/seguimiento` con agente, fecha, hora y notas. El backend crea el `follow_up_v2` y registra un evento en el timeline. Al completar, recarga los datos de la pantalla y muestra toast de éxito.
**Requerido:** Sí

---

**Evento:** Cuando confirma reasignación de caso
**Tipo:** User Interaction
**Descripción:** POST `/api/v1/casos/:id/reasignar` con el icode del nuevo operador. El backend actualiza `agent_id` y `victim_case_team` en `victim_case`, reasigna al nuevo operador los seguimientos `PENDIENTE` (`follow_up_v2`) y las tareas pendientes (`case_task` con `status = 'ToDo'`), y actualiza `victim_case_owner_description`. En paralelo registra evento `REASIGNACION` en el timeline. Recarga datos.
**Requerido:** Sí

---

**Evento:** Cuando expande/colapsa detalle de una barrera
**Tipo:** User Interaction
**Descripción:** Toggle de `barrExpandido`: si el id de la barrera ya estaba seleccionado lo limpia, si no lo asigna. Vue muestra u oculta la descripción de la barrera y el componente `CaseTasks` filtrado por esa barrera.
**Requerido:** Sí

---

**Evento:** Cuando presiona "▶ Iniciar" seguimiento
**Tipo:** User Interaction
**Descripción:** Valida que la fecha programada del seguimiento ya haya llegado. Si no, muestra toast de error. Si sí, abre `FollowUpContactModal` con los datos del seguimiento (id, caseId, nombre caso, intentos) para que el agente registre el contacto.
**Requerido:** Sí

---

**Evento:** Cuando presiona "✏️ Editar" seguimiento
**Tipo:** User Interaction
**Descripción:** Pre-llena `editarSeg` con la fecha y hora actuales del seguimiento y pone `modalEditarSeg = true`. Visible para `sv`, o para `ro`/`op` dueño del caso (`caso.agentId === userICode`).
**Requerido:** Sí

---

**Evento:** Cuando confirma edición de seguimiento
**Tipo:** User Interaction
**Descripción:** Valida que la nueva fecha no sea pasada. PUT `/api/v1/seguimiento/:segId/editar` con la nueva fecha/hora. El backend actualiza `scheduled_date` y `scheduled_time` en `follow_up_v2` y registra evento de edición en el timeline. Recarga datos.
**Requerido:** Sí

---

**Evento:** Cuando presiona "📋 Reasignar" seguimiento
**Tipo:** User Interaction
**Descripción:** Guarda el seguimiento activo en `reasignarSegActual` y abre `modalReasignarSeg`. Si el rol actual es `op`, preselecciona `reasignarSegAgente` a sí mismo (único agente disponible para ese rol). Visible para `sv`, o para `ro`/`op` agente propio del seguimiento (`seg.agent_id === userICode`).
**Requerido:** Sí

---

**Evento:** Cuando confirma reasignación de seguimiento
**Tipo:** User Interaction
**Descripción:** PUT `/api/v1/seguimiento/:segId/reasignar` con el nuevo `agent_id`. En caso de éxito, POST al timeline con el evento `REASIGNACION_SEGUIMIENTO`. Recarga datos. La lista de agentes disponibles (`agentesParaReasignarSeg`) depende del rol: `sv` ve `agentesRO` (ya filtrado por team del caso); `op` solo se ve a sí mismo (un único elemento); `ro` ve `todosAgentes` filtrado por el team del caso.
**Requerido:** Sí

---

**Evento:** Cuando expande/colapsa historial de intentos
**Tipo:** User Interaction
**Descripción:** Toggle de `segExpandido`: si el id ya estaba seleccionado lo limpia, si no lo asigna. Vue muestra u oculta el bloque de intentos del seguimiento correspondiente. Solo aplica a seguimientos que tengan `follow_up_attempts`.
**Requerido:** Sí

---

**Evento:** Cuando presiona "🔄 Reintentar"
**Tipo:** User Interaction
**Descripción:** Llama `reintentar()` que invoca `cargarDatos()` de nuevo. Solo visible cuando la pantalla está en estado de error no-404.
**Requerido:** Sí

---

**Evento:** Cuando el modal de contacto completa una gestión
**Tipo:** User Interaction
**Descripción:** Callback `onFollowUpContactCompleted(data)` emitido por `FollowUpContactModal`. Si `data.type === 'rescheduled'` o `'closed'` recarga los datos de la pantalla para reflejar el nuevo estado del seguimiento.
**Requerido:** Sí

---

**Evento:** Cuando el modal de completar tarea emite "completed"
**Tipo:** User Interaction
**Descripción:** Callback `onCaseTaskCompleted(tarea)` cableado a `<case-task-modal ref="taskModal">` montado a nivel de pantalla. Recarga las tareas del caso (`cargarTareasCaso`). **Actualmente inactivo en la práctica**: el único control que abría este modal (`$refs.taskModal.open(testTaskId)`) vivía en un bloque de UI de desarrollo dentro del tab Derivaciones que hoy está comentado. `case-tasks.js` maneja sus propias tareas con una instancia interna independiente de `case-task-modal`.
**Requerido:** No — código vivo pero sin disparador activo; candidato a limpieza o a reconectar si se retoma la UI de prueba.

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos?
- [x] ¿Todos los botones de la UI tienen evento? — incluye "Asignarme este caso", "+ Nuevo seguimiento" y "Reasignar seguimiento" (abrir), antes no documentados
- [x] ¿Las validaciones de formularios están cubiertas?
- [x] ¿Los estados de error están cubiertos?
- [x] ¿Los callbacks de componentes hijos están cubiertos? — incluye el callback de `case-task-modal`, documentado como inactivo
