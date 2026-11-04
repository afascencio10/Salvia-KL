━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando elimina una barrera
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: eliminarBarrera(id) → registro-v3-app.js · RemoveBarrier() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { id: id de la barrera (`d~…` en el borrador) }   → «Eliminar» en la tabla de la tarjeta (solo si `!entitiesReadOnly`)

PASO 1 — window.confirm("¿Eliminar esta barrera del registro?"); SI cancela: TERMINAR
PASO 2 — [DELETE] /api/v1/victim-cases-v3/barreras/{id}
  Backend: id `d~…` → se quita del JSON; UUID de tabla → 409 `readOnly`; inexistente → 404
PASO 3 — `recargarEntidadesBarreras()`

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Después de materializar la barrera no se elimina desde Registro V3 (solo lectura); los cambios se hacen en el detalle del caso | PASO 2 |
