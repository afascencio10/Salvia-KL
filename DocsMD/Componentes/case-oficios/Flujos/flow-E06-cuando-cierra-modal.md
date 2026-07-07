━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cierra el modal de detalle
   Tipo: User Interaction
   Función: cerrarDetalle()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: botón "×" del header, botón "Cerrar" del footer, o click en el
overlay oscuro fuera de la caja del modal (@click.self)

INPUT: {
  (ninguno)
}


PASO 1 — Ocultar el modal

  oficioSeleccionado = null


PASO 2 — Vue re-evalúa reactivamente

  → El modal (.co-modal-overlay) desaparece (v-if="oficioSeleccionado" → false)
  → El grid de cards queda visible tal como estaba antes de abrir el modal
    (los filtros no se resetean)

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| Ninguno — evento trivial, sin llamadas a sistemas externos.             | —             |
