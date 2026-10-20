━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Alerta inminente"
   Tipo: User Interaction
   Función: activarAlertaInminente() → set_victim_case_v3.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Este botón proviene del prototipo React construido en una ronda anterior
> de esta misma iniciativa. Al releer `Observaciones y ajustes al MockUp de
> registro 2.0.pdf` se confirma que **sí es un elemento real del mockup
> revisado por el cliente** (no fue invención del prototipo): la observación
> #16 del cliente pregunta textualmente *"¿Dónde quedó ubicado el botón de
> alerta inminente?"* — es decir, el botón existe en el mockup, pero su
> ubicación exacta quedó sin resolver incluso para el cliente mismo. La
> observación #2 agrega que la habilitación secuencial de TZ/GS "considera
> completitud de víctima+hecho+agresor, **pendiente incorporación de
> situación Alerta inminente**" — o sea, también está pendiente si una
> Alerta inminente debe poder saltarse ese requisito de tríada completa.
> Ninguna de las dos cosas se resuelve en este plan — se documentan como
> GAPS explícitos (G-V3-13 y G-V3-09).

INPUT: {
  icode:        identificador del Caso actual                → estado de la pantalla
  agentId:      id del Agente Integral                        → sesión
  origen:       'header' | 'GS'                                → desde dónde se presionó el botón
  justificacion: texto, solo si origen === 'GS'                → <textarea> de SubsecciónE (secciones.md § GS)
}

> Dos puntos de entrada al mismo flujo (decisión de esta ronda, ver
> `registro-v3-data.md` § D9): el botón del header (acceso rápido, sin
> justificación) y la acción dentro de GS (`registro-v3-secciones.md` § GS
> pregunta 18), que pide una justificación breve por quedar documentada
> como parte formal de la gestión del caso. Ambos ejecutan el mismo PASO 2.

PASO 1 — Confirmar la activación
  Mostrar modal de confirmación: "¿Confirma que desea activar una alerta inminente para este caso?"
  SI origen === 'GS': el modal incluye el campo de justificación (obligatorio antes de confirmar)

SI el agente confirma:
  → PASO 2 — Registrar y notificar
    [POST] /api/v1/victim-cases/{icode}/alerta-inminente
    // payload:
    { "agent_id": agentId, "fecha": now(), "origen": origen, "justificacion": justificacion || null }
    → Se agrega una entrada a la bitácora del Caso: "{agentId} activó una alerta inminente el {fecha}" (+ justificación si vino de GS)
    → SI origen === 'GS': victim_case_form2.wants_imminent_alert = true (ver registro-v3-data.md, columnas nuevas)
    NotificationService.alertarCoordinacion({ icode, motivo: 'alerta_inminente_manual' })
  → Mostrar confirmación visual al agente (banner o badge persistente en el header)
  → CONTINÚA FLUJO GENERAL

SI NO confirma:
  → Cierra el modal sin hacer nada
  → TERMINAR ejecución

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                              | Paso afectado |
|---------------------------------------------------------------------------------------|---------------|
| G-V3-13 — Dónde ubicar el botón (pregunta literal del cliente, sin responder) | Todo el evento |
| G-V3-09 — Si la Alerta inminente debe saltarse el requisito de tríada completa para habilitar TZ/GS (el cliente lo marca como "pendiente incorporación") | PASO 1 |
| Destino real de la notificación (¿misma bandeja que "Alerta por riesgo" de Tamizaje, o un canal separado de emergencia?) | PASO 2 |
