━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "No información mínima"
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón "No información mínima"
               (step 'validacionPrevia', solo si canManage)

Precondición: step === 'validacionPrevia' Y canManage === true

INPUT: {
  icode:  reporte abierto en el modal
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Llamar a cerrar

  saving = true
  POST /api/v1/victim-contacts/{icode}/cerrar
  body: { close_reason: 'ni' }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/:icode/cerrar (nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — Validar precondiciones

  SI no hay sesión o sesión.role NOT IN ('op', 'sv') (permiso manage_victim_contacts):
    → HTTP 401   // CheckPermission responde 401, no 403
  vc = buscar victim_contact por icode
  SI vc no existe:
    → HTTP 404
  SI vc.close_reason != null O vc.victim_contact_status = 'i' (reporte legado invalidado):
    → HTTP 409 "el reporte ya está cerrado"
    → TERMINAR ejecución
  // GAP-2 confirmado 2026-10-11: op y sv pueden gestionar cualquier
  // reporte — ya no se exige vc.agent_id === sesión.userICode.


PASO 3 — Resolver el enum de close_reason

  reasonRow = SELECT * FROM victim_case_form2_enums
              WHERE victim_case_form2_enums_category = 'victim_contact_close_reason'
                AND victim_case_form2_enums_code = 'ni'


PASO 4 — Actualizar victim_contact

  UPDATE victim_contact SET
    victim_contact_close_reason = reasonRow.id,
    victim_contact_closed_date  = NOW()
  WHERE victim_contact_i_code = :icode


PASO 5 — Responder

  { ok: true, closeReason: 'ni', category: 'Inválido' }

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — transición de vista
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 6 — Mostrar confirmación

  saving = false
  closeReasonLabel = 'Inválido por No información Mínima'
  categoria = 'Inválido'
  step = 'cerrado'

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Ninguno — GAP-2 quedó resuelto (op y sv, ver PASO 2).
