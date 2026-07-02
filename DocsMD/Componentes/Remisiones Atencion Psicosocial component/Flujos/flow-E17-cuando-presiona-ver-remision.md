━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona "Ver remisión"
   Tipo: User Interaction
   Código: E-17
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en LinkVerRemision "Ver remisión →" dentro del bloque REMISIÓN


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Emitir evento

  $emit('ver-remision', {
    remisionId: remision.id,
    followUpId: remision.followUpId,
    remision:   remision
  })

PASO 2 — El componente no navega por sí mismo

  // El padre define la ruta destino (detalle remisión o seguimiento).
  // Ruta tentativa: `/salvia/follow_up_detail/${followUpId}`

  → FIN EJECUCIÓN ✓
