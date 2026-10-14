━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Función: mounted() → cargar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  caseId:   victim_case_i_code del caso   → prop
  barriers: array de barreras del caso    → prop (ya cargado por el padre, no se fetchea aquí)
}

PASO 1 — Activar estado de carga

  cargando = true
  error = null

PASO 2 — Disparar las 3 llamadas en paralelo

  Promise.all([
    fetch(`/api/v1/casos/${caseId}/entidades`),
    fetch(`/api/v1/casos/${caseId}/obligaciones-resumen`),
    fetch(`/api/v1/entity-letters?caseId=${caseId}`),
  ])

  SI cualquiera de las 3 responde con error (status != 200) o falla la red:
    → error = "No se pudo cargar el resumen de gestión institucional"
    → cargando = false
    → TERMINAR ejecución (se muestra ErrorMsg + BtnReintentar → vuelve a PASO 1)

  SI NO:
    → CONTINÚA FLUJO GENERAL

PASO 3 — Guardar los datos en el estado del componente

  entidades = respuesta de /entidades                     // array
  resumenObligaciones = respuesta de /obligaciones-resumen // { global, porSector, items }
  oficios = respuesta de /entity-letters                   // array

PASO 4 — Completar el desglose por sector con los sectores sin entidades

  PARA CADA sector EN ['he', 'js', 'pt']:
    SI sector NO está en resumenObligaciones.porSector:
      → agregar { sector, cantidadEntidades: 0, completadas: 0, total: 0 }

  → Angular/Vue re-evalúa los bindings reactivos
  → Se renderizan los 4 KPIs, los 3 anillos de sector, y los badges de barreras/oficios por estado
    (ver SUB-FLUJO: Cálculo de agregados derivados)

┌─────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Cálculo de agregados derivados (computed)   │
└─────────────────────────────────────────────────────────┘

  S1. pctGlobal = resumenObligaciones.global.total > 0
        ? round(completadas / total * 100)
        : null
  S2. toneObligaciones = toneParaPorcentaje(pctGlobal)   // <40 rojo, 40-79 ámbar, ≥80 verde, null gris
  S3. barrerasPorEstado = {
        Abierta:     barriers.filter(b => b.status === 'OPEN').length,
        'En gestión': barriers.filter(b => b.status === 'En Gestion').length,
        Resuelta:    barriers.filter(b => barreraEsResuelta(b.status)).length,
      }
  S4. barrerasSinResolver = barrerasPorEstado.Abierta + barrerasPorEstado['En gestión']
  S5. oficiosPorBucket = {
        'En gestión': oficios.filter(o => o.state !== 'radicado' && o.state !== 'respondido').length,
        Radicado:    oficios.filter(o => o.state === 'radicado').length,
        Respondido:  oficios.filter(o => o.state === 'respondido').length,
      }
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

PASO 5 — Finalizar carga

  cargando = false

→ FIN EJECUCIÓN ✓

✅ Confirmado: el campo de estado en la respuesta de `/entity-letters` es `.state` (verificado contra `case-oficios.js`, líneas 290/438/481) — usado tal cual en S5.
