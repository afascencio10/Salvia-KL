━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona el tile "Entidades agregadas"
   Tipo: User Interaction
   Función: abrirModal('entidades')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  entidades: array ya cargado en E01   → estado del componente
}

PASO 1 — Construir el objeto modal

  modal = {
    tipo: 'entidades',
    titulo: 'Entidades agregadas',
    filas: entidades,   // sin transformar — ya trae sector, entityBranchName, ubicación, oficiosCount, barrerasActivasCount
  }

PASO 2 — Renderizar el modal

  → Vue re-evalúa `v-if modal` → true
  → Se monta el ModalBox con la tabla: Sector | Entidad | Ubicación | Oficios | Barreras activas
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
