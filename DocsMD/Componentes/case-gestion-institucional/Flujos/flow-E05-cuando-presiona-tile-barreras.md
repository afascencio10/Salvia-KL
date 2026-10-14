━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona el tile "Barreras sin resolver"
   Tipo: User Interaction
   Función: abrirModalBarreras(estado = null)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  barriers: prop del componente   → recibida del padre, no fetcheada
  estado: null   // este evento siempre llama sin estado — ver flow-E06 para la variante filtrada
}

PASO 1 — Filtrar las barreras no resueltas

  filasFiltradas = barriers.filter(b => !barreraEsResuelta(b.status))
    // equivalente a: b.status === 'OPEN' || b.status === 'En Gestion'

PASO 2 — Construir el objeto modal

  modal = {
    tipo: 'barreras',
    titulo: 'Barreras sin resolver',
    filas: filasFiltradas,
  }

PASO 3 — Renderizar el modal

  → Se monta el ModalBox con la tabla: Sector | Estado | Descripción | Fecha
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
