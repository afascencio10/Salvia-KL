# flow-E08 — Cuando cierra modal de error al eliminar

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cierra modal de error al eliminar
   Tipo: User Interaction
   Función: closeDeleteBlockedModal()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  modal.visible: true
  modal.kind:    'deleteBlocked'
  blockedMessage: texto mostrado en el modal
}

Disparado por:
  • Botón "Entendido" / "Cerrar"
  • ✕ / backdrop / Escape (si se habilita)


PASO 1 — Cerrar y limpiar

  modal.visible = false
  modal.kind = null
  blockedMessage = null
  deleteTarget = null   // por si quedó residual del intento E07

  → Vue oculta el modal de error
  → La lista de duplas NO se modifica (no hubo soft-delete)

→ FIN EJECUCIÓN ✓
```
