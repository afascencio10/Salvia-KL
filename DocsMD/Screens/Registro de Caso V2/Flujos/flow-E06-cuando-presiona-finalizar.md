━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Finalizar"
   Tipo: User Interaction
   Función: goToCase()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  caseId:  UUID del caso recién activado   → obtenido en E-05, PASO 2/3
}


PASO 1 — Redirigir al detalle del caso

  SI caseId disponible:
    window.location.assign(`/salvia/casos/${caseId}/detalle`)

  SI NO disponible (E-05 terminó en 202 sin confirmar el resultado):
    window.location.assign('/salvia/lista-casos')
    // Fallback razonable — el caso existe en BD aunque el frontend
    // no haya podido confirmar su ID a tiempo

  → FIN EJECUCIÓN ✓

Idéntico en espíritu a "Ver caso" (E-03) de hacer-seguimiento — sin gaps.
