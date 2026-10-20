━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver Reporte" en una fila
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón de acción "Ver Reporte" de una fila
               (columna Acciones, generada por casos-component desde :buttons)

INPUT: {
  buttonId:  'ver'
  report:    objeto reporte completo de la fila (ReportListItem)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js (sin cambios — ya documentado como E-05)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Emitir hacia el padre

  this.$emit('action-clicked', { buttonId: 'ver', case: report })

  // Mismo mecanismo genérico ya usado por cualquier botón de :buttons.
  // casos-component no sabe qué significa 'ver' — es la pantalla quien decide.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — list_reportes.html (padre)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — onActionClicked({ buttonId, case })

  SI buttonId !== 'ver':
    → TERMINAR ejecución (no hay otros botones en esta pantalla)

  SI NO:
    this.$refs.gestionModal.open(case)
    // Abre gestion-reporte-modal — toda la lógica de validación previa,
    // esquema 3x3 y cierre vive ahí (M-01 en adelante). Ver
    // DocsMD/Componentes/Gestion Reporte Modal/

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Ninguno — planeación completa del modal en
`DocsMD/Componentes/Gestion Reporte Modal/gestion-reporte-modal-events.md`.
GAP-2 (¿puede `sv` gestionar, no solo ver?) quedó resuelto: sí, `sv` abre el
modal con los mismos botones de acción que `op`, sin modo solo-lectura.
