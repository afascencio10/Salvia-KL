# flow-E05 — Cuando abre confirmación de eliminar

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre confirmación de eliminar
   Tipo: User Interaction
   Función: openDeleteConfirm(dupla)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  dupla:  fila seleccionada de "Duplas activas"   → botón "Eliminar"
}

Precondiciones:
  isLoading === false
  isDeleting === false
  modal de formulario no está guardando


PASO 1 — Guardar dupla objetivo

  deleteTarget = {
    id:   dupla.id,
    name: dupla.name
  }
  deleteError = null


PASO 2 — Abrir modal de confirmación

  modal.kind = 'delete'
  modal.visible = true

  → Vue muestra:
      Título: "Eliminar dupla"
      Texto:  "¿Estás seguro de que deseas eliminar {deleteTarget.name}?
               Esta acción libera a sus miembros para ser asignados a otra dupla."
      Botones: Cancelar | Eliminar (destructivo)

→ FIN EJECUCIÓN ✓

// No valida uso en remisiones/sesiones aquí.
// Esa validación ocurre en E07 al confirmar (backend).
```
