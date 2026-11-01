━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Proceder al formulario"
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — hoy es el botón "Abrir" de cada fila del Cabezote
   Funciones: abrirCaso(icode) → registro-v3-app.js · GetCase() → victim_case_v3_controller.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> La planeación distinguía "Proceder al formulario" (solo Borradores, con trazabilidad de retome) de otros accesos.
> En el código hay **un solo acceso**: el botón "Abrir" de cualquier fila (Borrador, Iniciado o Terminado). No se
> registra quién retoma ni cuándo (`/retome` no existe). `abrirCaso` también lo usa "+ Nuevo registro" (flow-E06).

INPUT: {
  icode:  del Caso de la fila   → resultadosBusqueda[i].icode (o el icode que devuelve el draft)
}

PASO 1 — Cargar el Caso completo
  [GET] /api/v1/victim-cases-v3/{icode}
  Respuesta (CaseStateV3): {
    summary:        { icode, names, lastNames, docType, docNumber, status },
    answers:        todas las respuestas guardadas, por FieldKey (CI_1, DP_5, …), con valores dummy incluidos,
    answeredKeys:   las keys que el agente respondió alguna vez,
    triadaCompleta: ¿todas las obligatorias de la tríada ya respondidas?,
    entities:       entidades del Caso (entity_case con sede, sector y ubicación),
    barriers:       barreras del Caso (barrier_v2),
    tzResultado:    { riskType, riskLevel, score } o null si el Tamizaje nunca se guardó
  }

PASO 2 — Poblar el estado
  icode = icode; resumenCaso = data.summary
  answers = SOLO las keys que están en answeredKeys   // las demás son valores dummy que el agente nunca respondió
  entidades = data.entities; barreras = data.barriers; tzResultado = data.tzResultado; triadaCompleta = data.triadaCompleta

PASO 3 — Abrir el Formulario
  etapa = 'formulario'; moduloActivo = 'CI'    // siempre entra por CI, no por "el último módulo con progreso"
  Vue muestra el encabezado, las pestañas (TZ/GS bloqueadas si !triadaCompleta) y el módulo CI

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Planeado y no implementado: trazabilidad del retome (`POST …/retome`) y entrada en el módulo con menor progreso | PASO 1 · PASO 3 |
| Si la petición falla, la pantalla no muestra ningún error y se queda en el Cabezote | PASO 1 |
