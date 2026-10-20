━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda una sección/módulo
   Tipo: User Interaction
   Función: saveSection() → dinamic-form.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  moduloActual: módulo visible en este momento   → estado de la pantalla
  respuestas:   valores acumulados del módulo     → estado local (ver flow-E11)
  icode:        identificador del Caso            → estado de la pantalla
}

Validar formulario:
  • Todas las preguntas [1er contacto] visibles del módulo: required
  • Las [máx. 2º contacto] y [Deseable]: opcional en este guardado
  SI hay obligatorias [1er contacto] sin responder:
    → Mostrar error por campo
    → Mostrar toast "Completa las preguntas obligatorias de este módulo"
    → TERMINAR ejecución

  SI validación ok:
    → CONTINÚA FLUJO GENERAL

PASO 1 — Persistir el módulo

[POST] /api/v1/forms/saveSection
// payload:
{
  "icode":   icode,
  "section": moduloActual,
  "answers": respuestas   // solo las del módulo actual
}

PASO 2 — Transición de estado del Caso: Borrador → Iniciado

> **Corregido para quedar consistente con el umbral único de `triadaCompleta`**
> (ver sub-flujo "Recalcular habilitación secuencial" en `flow-E11` y
> `registro-v3-data.md`). Ya no se usa "primer guardado exitoso" — ese
> criterio no exigía ningún campo mínimo y permitía generar un Caso formal
> con una sola respuesta trivial. Ahora el mismo umbral que desbloquea TZ/GS
> es el que genera el Caso.

SI caso.estado === 'Borrador' Y triadaCompleta === true (recalculado tras persistir este guardado — todas las obligatorias [1er contacto] de CI+DP+CA+HV+AG ya respondidas):
  → caso.estado = 'Iniciado'
  → [PUT] /api/v1/victim-cases/{icode}/estado  { "estado": "Iniciado" }
  → Se genera/confirma el ID único de Caso en este momento (ver `registro-v3-data.md` § D8 — mecanismo de valores "dummy" para Borradores heredado de V2-real)
  → CONTINÚA FLUJO GENERAL

SI NO (caso.estado !== 'Borrador', o triadaCompleta sigue en false):
  → caso.estado se mantiene
  → CONTINÚA FLUJO GENERAL

PASO 3 — Confirmar al agente
  Mostrar toast "Sección guardada ✓"
  Vue re-evalúa BadgeEstadoCaso y la BarraProgreso del módulo (ya en 100% de obligatorias)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                      | Paso afectado |
|-------------------------------------------------------------------------------|---------------|
| Si `saveSection` falla (red, backend) ¿se reintenta automáticamente o se le informa al agente para que reintente manualmente? | PASO 1 |
