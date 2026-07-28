# flow-E06 — Cuando presiona "Ver detalle"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver detalle"
   Tipo: User Interaction
   Función: goToRelationDetail(item)
   Estado: implementado (toast temporal — destino pendiente)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  item.entityBranchICode: entity_branch.entity_branch_i_code
  item.entityCaseId:      entity_case.id
}

PASO 1 — SI !item.entityBranchICode → toast error → TERMINAR

PASO 2 — Destino temporal
  → Mostrar toast: "Detalle de entidad pendiente (sede: {entityBranchICode})"
  → NO navegar (la ruta /salvia/entidad/:id aún no existe)

// Cuando exista la pantalla (case-entities E09):
// window.location.href = '/salvia/entidad/' + entityBranchICode

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                      | Paso afectado |
|----------------------------------------------------------|---------------|
| Crear pantalla /salvia/entidad/:id (case-entities)       | PASO 2        |
```
