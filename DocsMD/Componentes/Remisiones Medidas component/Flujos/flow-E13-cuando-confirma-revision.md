━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda la revisión técnica
   Tipo: User Interaction / Backend write
   Código: E-13
   Prerequisito: M-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: botón "Guardar revisión", habilitado con profesional, concepto y resultado


INPUT: {
  id:              revisionItem.id
  professionalId:  revisionProfesional     → general_user_i_code
  concepto:        revisionConcepto.trim
  resultado:       revisionResultado       → uno de los tres códigos de M-01
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Guardar

  POST /api/v1/emergency-measure-referrals/{id}/revision
  Content-Type: application/json

  {
    "professionalId": "general_user_i_code",
    "concepto": "texto",
    "resultado": "aprobado_para_caracterizacion_me"
  }

PASO 2 — Si responde 200

  cerrar el modal (E-12)
  repetir E-01 con los filtros de UI que sigan activos

  La fila deja de mostrar "Revisar (Equipo Técnico)".
  Pasa a mostrar el badge del estado nuevo y el bloque
  "Concepto técnico ( {fecha} ): {resultado}" con "Por: {profesional}".
  El texto del concepto queda guardado y no se pinta en la fila.

  En supervisor, la card Pendiente de Revisión baja y sube la card del
  estado destino (Por asignar, si el resultado fue la aprobación).
  En equipo técnico la fila sale de la bandeja, porque el recorte solo
  incluye pendiente_de_revision.


PASO 3 — Si responde 4xx o 5xx

  El modal sigue abierto con los datos escritos.
  Se muestra el mensaje de error. El listado no se recarga.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  La remisión tiene que existir y estar en pendiente_de_revision.
  Si el estado es otro, responde 409 y no escribe.

  Mapa de resultado → status:

    aprobado_para_caracterizacion_me → por_asignar
    devuelto_para_ajuste             → devuelta_para_justificacion
    devuelto_para_acompanamiento     → devuelta_a_seguimiento

  UPDATE salvia.emergency_measure_referral
     SET review_professional_id = {professionalId},
         review_concept         = {concepto},
         review_result          = {resultado},
         review_date            = fecha del día en el servidor (DATE),
         status                 = {estado mapeado},
         updated_at             = now()
   WHERE id = {id}
     AND status = 'pendiente_de_revision'
     AND deleted_at IS NULL

  Rechaza con 400 si falta professionalId, si el concepto queda vacío
  al recortar espacios, si el resultado no es uno de los tres códigos,
  o si el profesional no está en GET .../profesionales-revision.
  Rechaza con 401 si el rol no tiene review_remisiones_medidas
  (sm, em y sv).
  Rechaza con 404 si el id no existe.

  No cambia assigned_agent_id ni las medidas hijas.
  Reemplaza la revisión anterior si la hubiera. No hay historial.

  Respuesta 200: el ReferralListItem actualizado (mismo shape de E-01).

  → FIN EJECUCIÓN ✓
