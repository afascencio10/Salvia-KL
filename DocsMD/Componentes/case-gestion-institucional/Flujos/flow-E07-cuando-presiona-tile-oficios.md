━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona el tile "Oficios en gestión"
   Tipo: User Interaction
   Función: abrirModalOficios('En gestión')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  oficios: array ya cargado en E01   → estado del componente
  bucket: 'En gestión'   // fijo — este tile siempre abre el mismo bucket
}

PASO 1 — Filtrar los oficios en gestión

  filasFiltradas = oficios.filter(o => o.state !== 'radicado' && o.state !== 'respondido')

PASO 2 — Construir el objeto modal

  modal = {
    tipo: 'oficios',
    titulo: 'Oficios — En gestión',
    filas: filasFiltradas,
  }

PASO 3 — Renderizar el modal

  → Se monta el ModalBox con la tabla: Entidad | Estado | Fecha
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
