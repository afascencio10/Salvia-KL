━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver" en un comentario
   Tipo: User Interaction
   Función: abrirDetalle(gestion, scope) · cerrarDetalle()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

La sección de comentarios ya no usa `case_articulating_mechanism_comment`.
Lista las `case_task` `gestion_mecanismo` Done del escalamiento (`gestiones` en el GET).

INPUT: {
  gestion.result
  gestion.formData
  gestion.assignedUserName
  gestion.completedAt
  m.scope
}

PASO 1 — El usuario pulsa "Ver" junto al result
  → detalleGestion = gestion
  → detalleScope = scope
  → detalleAbierto = true

PASO 2 — Overlay de solo lectura
  Título: "Gestión del mecanismo"
  Subtítulo: Municipal | Departamental | Nacional
  Pares label/valor de formData (campos vacíos no se muestran)
  Labels alineados al formulario Excel / buildGestionMecanismoResult

PASO 3 — Cerrar (botón Cerrar, × o overlay)
  → detalleAbierto = false

El endpoint POST /api/v1/casos/{caseId}/mecanismos/{mid}/comentarios
sigue existiendo; la card no lo llama.
