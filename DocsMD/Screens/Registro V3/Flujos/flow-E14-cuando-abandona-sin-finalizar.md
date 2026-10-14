━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abandona la pantalla sin finalizar
   Tipo: User Interaction (primario) + Lifecycle (respaldo)
   Función: guardarYSalir() → set_victim_case_v3.html · también beforeUnmount()/beforeunload como respaldo
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Camino primario confirmado al revisar `Mockup_Reportes_Registro (1).html`:
> existe un botón explícito **"Finalizar llamada (guardar y salir)"** en la
> navegación inferior del Formulario — el agente lo presiona deliberadamente
> al colgar la llamada, no se depende solo de detectar el cierre de pestaña.
> El respaldo de `beforeunload` (ya documentado en la versión anterior de
> este flujo) se mantiene para el caso en que el agente cierre sin presionar
> el botón.

INPUT: {
  stage: etapa actual de la pantalla       → estado de la pantalla
  caso:  Caso en curso, si aplica           → estado de la pantalla
}

SI stage !== 'formulario':
  → No hay nada que guardar (Habeas data y Cabezote no persisten estado parcial)
  → TERMINAR ejecución

SI stage === 'formulario':
  → PASO 1 — Asegurar guardado del módulo actual
    SI el módulo actual tiene cambios sin guardar:
      → Ver flujo: Cuando guarda una sección/módulo (ejecuta el mismo saveSection, sin exigir que las obligatorias [1er contacto] estén completas — guardado parcial de cierre de llamada)
  → PASO 2 — El Caso permanece en su estado actual (Borrador o Iniciado)
    → Queda visible en el Cabezote para que cualquier agente lo retome (ver flow-E07)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                       | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| El guardado de cierre de llamada no exige obligatorias [1er contacto] completas — confirmar que `saveSection` soporta un modo "parcial" sin fallar la validación normal | PASO 1 |
| Detección real de cierre de pestaña/navegador (`beforeunload`) es poco confiable para disparar una llamada async — confirmar si basta con el guardado automático periódico en vez de depender de este evento | Todo el evento |
