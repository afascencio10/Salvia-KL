# flow-E06 — Cuando cancela confirmación de eliminar

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cancela confirmación de eliminar
   Tipo: User Interaction
   Función: closeDeleteConfirm()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  modal.visible: true
  modal.kind:    'delete'
  deleteTarget:  dupla pendiente de borrar
}

Disparado por:
  • Botón "Cancelar"
  • ✕ / backdrop / Escape (si se habilita)


PASO 1 — Verificar que se puede cerrar

SI isDeleting === true:
  → No cerrar
  → TERMINAR

SI isDeleting === false:
  → CONTINÚA FLUJO GENERAL


PASO 2 — Cerrar y limpiar

  modal.visible = false
  modal.kind = null
  deleteTarget = null
  deleteError = null

  → Vue oculta el modal
  → No hay cambios en BD ni en listas

→ FIN EJECUCIÓN ✓
```
