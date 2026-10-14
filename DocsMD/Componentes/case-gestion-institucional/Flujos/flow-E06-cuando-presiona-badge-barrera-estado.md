━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona un badge de "Barreras por estado" (Abierta / En gestión / Resuelta)
   Tipo: User Interaction
   Función: abrirModalBarreras(estado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  barriers: prop del componente          → recibida del padre
  estado: 'Abierta' | 'En gestión' | 'Resuelta'   → click del usuario sobre el badge
}

PASO 1 — Filtrar las barreras según el estado seleccionado

  SEGÚN estado:
    CASO 'Abierta':
      → filasFiltradas = barriers.filter(b => b.status === 'OPEN')
      → CONTINÚA
    CASO 'En gestión':
      → filasFiltradas = barriers.filter(b => b.status === 'En Gestion')
      → CONTINÚA
    CASO 'Resuelta':
      → filasFiltradas = barriers.filter(b => barreraEsResuelta(b.status))
      → CONTINÚA

PASO 2 — Construir el objeto modal

  modal = {
    tipo: 'barreras',
    titulo: `Barreras — ${estado}`,
    filas: filasFiltradas,
  }

PASO 3 — Renderizar el modal

  → Misma tabla que flow-E05 (Sector | Estado | Descripción | Fecha), acotada al estado
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
