━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona un anillo de sector (Salud, Justicia o Protección)
   Tipo: User Interaction
   Función: abrirModalObligaciones(sector)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  resumenObligaciones.items: array ya cargado en E01   → estado del componente
  sector: código de 2 letras ('he' | 'js' | 'pt')       → click del usuario sobre el anillo
}

PASO 1 — Filtrar las filas por sector

  filasFiltradas = resumenObligaciones.items.filter(item => item.sector === sector)

PASO 2 — Construir el objeto modal

  modal = {
    tipo: 'obligaciones',
    titulo: `Obligaciones completadas — ${sectorLabels[sector]}`,   // ej. "— Justicia"
    filas: filasFiltradas,
  }

PASO 3 — Renderizar el modal

  → Misma tabla que flow-E03 (Entidad | Sector | Obligación | Estado), ya acotada al sector
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
