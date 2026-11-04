━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando elimina una entidad acudida
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: eliminarEntidad(relId) → registro-v3-app.js · RemoveEntity() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { relId: id de la entidad (`d~…` en el borrador) }   → «Eliminar entidad» (solo si `!entitiesReadOnly`)

PASO 1 — window.confirm("¿Eliminar esta entidad del registro? También se eliminarán sus barreras asociadas."); SI cancela: TERMINAR
PASO 2 — [DELETE] /api/v1/victim-cases-v3/entidades/{relId}
  Backend: id `d~…` → se quitan del JSON la entidad y **todas sus barreras** (aún no hay filas ni tareas que cancelar); UUID de tabla (ya materializada) → 409 `readOnly`; id inexistente → 404
PASO 3 — `recargarEntidadesBarreras()`; SI fue 409 `readOnly` la pantalla pasa a solo lectura

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| G-V3-22 resuelto: la cascada de borrado ocurre solo en el borrador | PASO 2 |
