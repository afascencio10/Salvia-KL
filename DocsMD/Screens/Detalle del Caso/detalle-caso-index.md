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
| E04 | Cuando presiona "Ver →" en banner de próximo seguimiento | User Interaction | — |
| E05 | Cuando presiona "Reasignar" (caso) | User Interaction | — |
| E06 | Cuando guarda nuevo seguimiento | User Interaction | — |
| E07 | Cuando confirma reasignación de caso | User Interaction | — |
| E08 | Cuando presiona "▶ Iniciar" seguimiento | User Interaction | — |
| E09 | Cuando presiona "✏️ Editar" seguimiento | User Interaction | — |
| E10 | Cuando confirma edición de seguimiento | User Interaction | — |
| E11 | Cuando confirma reasignación de seguimiento | User Interaction | — |
| E12 | Cuando expande/colapsa historial de intentos | User Interaction | — |
| E13 | Cuando presiona "🔄 Reintentar" | User Interaction | — |
| E14 | Cuando el modal de contacto completa una gestión | User Interaction | — |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

**Evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** Dispara dos llamadas paralelas: carga el detalle completo del caso (datos, seguimientos, barreras, derivaciones, timeline) y carga la lista de todos los agentes del equipo para resolver nombres en la UI.
**Requerido:** Sí

---

**Evento:** Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** Actualiza `tabActiva` con el id del tab presionado. Vue re-renderiza el contenido del panel mostrando la sección correspondiente (info, derivaciones, barreras, seguimientos, timeline).
**Requerido:** Sí

---

**Evento:** Cuando presiona "📋 Info caso"
**Tipo:** User Interaction
**Descripción:** Pone `mostrarInfoCaso = true`, lo que hace visible el componente `CaseInfo` en modo modal con el detalle completo del caso.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Ver →" en banner de próximo seguimiento
**Tipo:** User Interaction
**Descripción:** Asigna `tabActiva = 'seguimientos'` para llevar al usuario directamente a la pestaña de seguimientos.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Reasignar" (caso)
**Tipo:** User Interaction
**Descripción:** Limpia `reasignarOperador` y pone `modalReasignar = true`. Visible únicamente para rol `sv`. El dropdown muestra operadores filtrados por el team del caso.
**Requerido:** Sí

---

**Evento:** Cuando guarda nuevo seguimiento
**Tipo:** User Interaction
**Descripción:** Valida fecha (requerida, no pasada) y agente (requerido para `sv`). Si válido, POST `/api/v1/casos/:id/seguimiento` con agente, fecha, hora y notas. El backend crea el `follow_up_v2` y registra un evento en el timeline. Al completar, recarga los datos de la pantalla.
**Requerido:** Sí

---

**Evento:** Cuando confirma reasignación de caso
**Tipo:** User Interaction
**Descripción:** POST `/api/v1/casos/:id/reasignar` con el icode del nuevo operador. El backend actualiza `agent_id` y `victim_case_team` en `victim_case`, reasigna los seguimientos PENDIENTE al nuevo operador y actualiza `victim_case_owner_description`. En paralelo registra evento REASIGNACION en el timeline. Recarga datos.
**Requerido:** Sí

---

**Evento:** Cuando presiona "▶ Iniciar" seguimiento
**Tipo:** User Interaction
**Descripción:** Valida que la fecha programada del seguimiento ya haya llegado. Si no, muestra toast de error. Si sí, abre `FollowUpContactModal` con los datos del seguimiento (id, caseId, nombre caso, intentos) para que el agente registre el contacto.
**Requerido:** Sí

---

**Evento:** Cuando presiona "✏️ Editar" seguimiento
**Tipo:** User Interaction
**Descripción:** Pre-llena `editarSeg` con la fecha y hora actuales del seguimiento y pone `modalEditarSeg = true`. Solo visible para `sv` o `ro` dueño del caso.
**Requerido:** Sí

---

**Evento:** Cuando confirma edición de seguimiento
**Tipo:** User Interaction
**Descripción:** Valida que la nueva fecha no sea pasada. PUT `/api/v1/seguimiento/:segId/editar` con la nueva fecha/hora. El backend actualiza `scheduled_date` y `scheduled_time` en `follow_up_v2` y registra evento de edición en el timeline. Recarga datos.
**Requerido:** Sí

---

**Evento:** Cuando confirma reasignación de seguimiento
**Tipo:** User Interaction
**Descripción:** PUT `/api/v1/seguimiento/:segId/reasignar` con el nuevo `agent_id`. En caso de éxito, POST al timeline con el evento REASIGNACION_SEGUIMIENTO. Recarga datos. La lista de agentes disponibles se filtra por team del caso para `op/ro`, y usa `agentesRO` para `sv`.
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

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos?
- [x] ¿Todos los botones de la UI tienen evento?
- [x] ¿Las validaciones de formularios están cubiertas?
- [x] ¿Los estados de error están cubiertos?
- [x] ¿Los callbacks de componentes hijos están cubiertos?
