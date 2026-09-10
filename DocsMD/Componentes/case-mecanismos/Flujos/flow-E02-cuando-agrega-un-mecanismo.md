━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega un mecanismo
   Tipo: User Interaction
   Función: guardar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: form { scope, departmentId, cityId, townCode, assignedUserId, salviaCommitment }

PASO 1 — Validar campos según tipo
  municipal: depto + ciudad + municipio + usuario
  departamental: depto + usuario
  nacional: usuario

PASO 2 — POST /api/v1/casos/{caseId}/mecanismos

Backend:
  - el usuario debe existir, estar activo y tener rol mc
  - debe tener articulating_mechanism del tipo pedido y coincidir con la ubicación
  - persiste case_articulating_mechanism con status = active (abierto)
  - crea CaseTimelineEvent categoría Mecanismos

SI ok:
  → cierra modal
  → recarga el listado (E01)
