━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Cerrar por no contacto"
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón del CompleteBanner, step 'esquema3x3'

Precondición: step === 'esquema3x3' Y diasCompletos >= 3 Y canManage === true

INPUT: {
  icode:  reporte abierto en el modal
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Llamar a cerrar

  saving = true
  POST /api/v1/victim-contacts/{icode}/cerrar
  body: { close_reason: 'nc' }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/:icode/cerrar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — Validar precondiciones

  SI no hay sesión o sesión.role NOT IN ('op', 'sv'): → HTTP 401   // CheckPermission responde 401, no 403
  vc = buscar victim_contact por icode
  SI vc no existe:
    → HTTP 404
  SI vc.close_reason != null O vc.victim_contact_status = 'i' (legado invalidado):
    → HTTP 409
  // GAP-2 confirmado 2026-10-11: op y sv pueden gestionar cualquier reporte.

  SI NO — validación adicional específica de este motivo:
  // Corrección (2026-10-12, respuesta a comentario de Felipe): un "día
  // completo" NO es un día con cualquier cantidad de intentos fallidos —
  // es un día con AL MENOS 3. Verificado contra la fuente literal del
  // mockup (diasCompletosDe en Mockup_Reportes_Registro.html):
  //   agrupar intentos fallidos por día
  //   → contar cuántos de esos días tienen >= 3 intentos fallidos
  // "3x3" es literal: 3 intentos fallidos, en cada uno de 3 días distintos
  // → mínimo 9 intentos fallidos en total para poder cerrar por esta vía,
  // no 3 como decía una versión anterior de esta planeación.
  // La columna es timestamp naive en UTC: se declara UTC y luego se convierte a Bogotá.
  fallidosPorDia = SELECT (victim_contact_attempt_date AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota')::date AS dia, COUNT(*) AS n
                   FROM victim_contact_attempt
                   WHERE victim_contact_attempt_victim_contact = :icode
                     AND victim_contact_attempt_result = 'fallido'
                     AND victim_contact_attempt_deleted_at IS NULL
                   GROUP BY dia
  diasCompletos = COUNT(dia EN fallidosPorDia WHERE n >= 3)

  SI diasCompletos < 3:
    → HTTP 409 "el esquema 3x3 no está completo todavía"
    → TERMINAR ejecución

  // Esta validación server-side es necesaria porque el botón solo se oculta
  // en el cliente — sin ella, una llamada directa al endpoint podría cerrar
  // un reporte sin haber completado el esquema.


PASO 3 — Resolver enum y actualizar

  reasonRow = enum categoría 'victim_contact_close_reason', code 'nc'

  UPDATE victim_contact SET
    victim_contact_close_reason = reasonRow.id,
    victim_contact_closed_date  = NOW()
  WHERE victim_contact_i_code = :icode


PASO 4 — Responder

  { ok: true, closeReason: 'nc', category: 'Inválido' }

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — transición de vista
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 5 — Mostrar confirmación

  saving = false
  closeReasonLabel = 'Inválido por Cierre por no contacto'
  categoria = 'Inválido'
  step = 'cerrado'

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Ninguno.
