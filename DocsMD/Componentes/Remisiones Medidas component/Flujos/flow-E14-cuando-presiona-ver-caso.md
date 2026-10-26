━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona "Ver caso"
   Tipo: User Interaction
   Código: E-14
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en "Ver caso →" dentro de la columna Caso


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Emitir

  $emit('ver-caso', {
    caseICode:  item.caseICode,
    remisionId: item.id
  })

PASO 2 — El componente no navega

  El padre escucha @ver-caso.
  Remisiones Medidas y Mis Remisiones dejan el handler vacío: la ruta de
  destino de estos casos todavía no está definida para esta bandeja.

  → FIN EJECUCIÓN ✓
