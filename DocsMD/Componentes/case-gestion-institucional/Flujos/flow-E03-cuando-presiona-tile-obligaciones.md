━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona el tile "Obligaciones completadas"
   Tipo: User Interaction
   Función: abrirModalObligaciones(sector = null)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  resumenObligaciones.items: array ya cargado en E01   → estado del componente
  sector: null   // este evento siempre llama sin sector — ver flow-E04 para la variante filtrada
}

PASO 1 — Construir el objeto modal

  modal = {
    tipo: 'obligaciones',
    titulo: 'Obligaciones completadas',
    filas: resumenObligaciones.items,   // todas, sin filtrar por sector
  }

PASO 2 — Renderizar el modal

  → Se monta el ModalBox con la tabla: Entidad | Sector | Obligación | Estado
    (Estado: "✅ Completada" si item.completada, "⭕ Pendiente" si no)
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
