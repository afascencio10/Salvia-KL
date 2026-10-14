━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver Reporte" en una fila
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón de acción "Ver Reporte" de una fila

INPUT: {
  buttonId:  'ver'
  report:    objeto reporte completo de la fila (ReportListItem)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — get_my_reportes.html (padre)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — onActionClicked({ buttonId, case })

  SI buttonId !== 'ver':
    → TERMINAR ejecución

  SI NO:
    this.$refs.gestionModal.open(case)

  // Idéntico a E-04 de Listado Reportes — op y sv pueden gestionar
  // cualquier reporte (GAP-2), así que el comportamiento es el mismo aquí
  // y en Listado Reportes, sin distinción por pertenencia.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Ninguno adicional — hereda los de `Listado Reportes` E-04 y los del propio
`gestion-reporte-modal` (ver su índice de eventos).
