━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona «Ir al borrador del caso»
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — WP-1 T7 (2026-11-02)
   Función: irAlBorrador()    (paso 'borradorIniciado' del modal)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { draftCaseIcode ← detalle del reporte (M-01) }

PASO 1 — SI draftCaseIcode = '' → TERMINAR ejecución      // no debería ocurrir: el botón solo existe en el paso 'borradorIniciado'
PASO 2 — window.location.assign('/salvia/registro-v3?caso=' + encodeURIComponent(draftCaseIcode))
         → Registro V3 E19: abre el borrador en el formulario

No hay llamada a backend desde el modal.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| El botón lo ven op y sv (quienes abren el modal); si el borrador lo creó otro agente igual se puede retomar | PASO 1 |
