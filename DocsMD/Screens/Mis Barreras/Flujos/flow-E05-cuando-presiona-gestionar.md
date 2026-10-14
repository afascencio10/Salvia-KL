# flow-E05 — Cuando presiona Gestionar

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona Gestionar
   Tipo: User Interaction
   Función: openGestionar(barrera)
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: fila de la tabla (barrierId, caseId, taskId?, entityBranchId, entityName, description)

PASO 1 — this.$refs.taskModal.openArticulacion({
  taskId, barrierId, caseId, entityBranchId, entityName, barrierDescription
})

SI taskId:
  → case-task-modal.open(taskId) + formulario articulacion_territorial
SI NO:
  → tarea en memoria type=articulacion_territorial, description="Gestión propia"

Al @completed:
  toast "Gestión guardada."
  loadBarreras()  (refresh)
```
