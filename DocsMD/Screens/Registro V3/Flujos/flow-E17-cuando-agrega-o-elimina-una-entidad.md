━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega o elimina una Entidad
   Tipo: User Interaction
   Función: agregarEntidad(tipo) / eliminarEntidad(tipo, index) → components/dinamic-form.js (repeater genérico)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Modela el repeater real "Identificación de Entidades" de Seguimiento
> (`repeater_group` `ee42ef36...` para Acudidas y `5b5b7713...` para
> Activación) — ver `registro-v3-data.md` § D9. Mecanismo genérico de
> `dinamic-form`, no exclusivo de Registro V3.

INPUT: {
  tipo:  'acudidas' | 'activacion'   → qué repeater (según cuál boolean disparador esté en true)
  index: posición del ítem a eliminar, solo aplica para eliminarEntidad
}

PASO 1 — Agregar una entidad

SI se presiona "+ Siguiente entidad":
  → Crea un nuevo ítem vacío en `respuestas.entidades{Acudidas|Activacion}`
  → Renderiza sus campos (Sector, Departamento, Ciudad, Municipio, Entidad, radicado, fecha, canal si es Activación, acciones completadas/pendientes)
  → El nuevo ítem no cuenta para completitud hasta que se respondan sus campos obligatorios
  → CONTINÚA FLUJO GENERAL

PASO 2 — Eliminar una entidad

SI se presiona "✕ Eliminar" sobre un ítem:
  → Mostrar confirmación "¿Eliminar esta entidad del registro?"
  SI confirma:
    → Quita el ítem en `index` del arreglo correspondiente
    → Recalcula la barra de progreso del módulo GS
    → CONTINÚA FLUJO GENERAL
  SI NO confirma:
    → TERMINAR ejecución, no se elimina nada

PASO 3 — Recalcular cascada Sector → Entidad dentro del ítem

Cuando cambia el campo "Sector de la entidad" de un ítem ya agregado:
  → Recalcula las opciones de "Entidad" para ESE ítem específico (catálogo `entidadesDisponibles` por sector, dinámico — ver GAP G-V3-21)
  → Si ya había una Entidad seleccionada y deja de ser válida para el nuevo sector, se limpia

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                     | Paso afectado |
|-----------------------------------------------------------------------------------------------|---------------|
| G-V3-21 — tabla real donde `dinamic-form` persiste las respuestas de un repeater (no es `victim_case_form2`) | PASO 1 |
| Fuente exacta del catálogo dinámico `entidadesDisponibles` por sector (¿tabla maestra propia, o el mismo catálogo que usa Seguimiento hoy?) | PASO 3 |
| Si hay un máximo de entidades por repeater (Seguimiento no define `max_repetitions` en estos dos repeaters — confirmar si V3 debe limitarlo) | PASO 1 |
