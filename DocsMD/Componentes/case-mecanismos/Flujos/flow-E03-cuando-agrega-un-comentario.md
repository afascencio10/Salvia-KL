━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega un comentario
   Tipo: User Interaction
   Función: agregarComentario(m)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { mechanismId, body }

PASO 1 — POST /api/v1/casos/{caseId}/mecanismos/{mid}/comentarios
  Body: { "body": "..." }

Backend guarda author_id, author_name (sesión) y created_at.
No escribe evento de Timeline.

SI ok:
  → append del comentario en la card
  → limpia el textarea
