━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Función: mounted() → cargar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { caseId, userId, userRole }

PASO 1 — GET /api/v1/casos/{caseId}/mecanismos

SI error:
  → error = "No se pudieron cargar los mecanismos de este caso."
  → TERMINAR

SI ok:
  → mecanismos = resultado (array)
  → cada fila trae gestiones[] (case_task gestion_mecanismo Done: result, formData, completedAt, assignedUserName)
  → si vacío: empty state
  → si hay filas: lista con badge de alcance, usuario asignado, ubicación, compromiso y comentarios (= gestiones)
