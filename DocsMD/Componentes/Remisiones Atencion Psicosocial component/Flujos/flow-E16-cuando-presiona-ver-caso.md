━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona "Ver caso"
   Tipo: User Interaction
   Código: E-16
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en LinkVerCaso "Ver caso →" dentro del bloque CASO


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Emitir evento

  $emit('ver-caso', {
    caseICode: remision.caseICode,
    remision:  remision
  })

PASO 2 — El componente no navega por sí mismo

  // El padre escucha @ver-caso y redirige, p. ej.:
  // window.location = `/salvia/case_detail/${payload.caseICode}`

  → FIN EJECUCIÓN ✓
