━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega un mecanismo
   Tipo: User Interaction
   Función: guardar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: form { scope, departmentId, cityId, townCode, assignedUserId, salviaCommitment }

PASO 1 — Validar campos según alcance
  municipal: depto + ciudad + municipio + usuario
  departamental: depto + usuario (ciudad opcional)
  nacional: usuario

PASO 2 — POST /api/v1/casos/{caseId}/mecanismos

Backend:
  - el usuario debe existir, estar activo y tener rol mc
  - debe coincidir con el filtro geo del alcance
  - persiste case_articulating_mechanism
  - crea CaseTimelineEvent categoría Mecanismos

SI ok:
  → cierra modal
  → recarga el listado (E01)
