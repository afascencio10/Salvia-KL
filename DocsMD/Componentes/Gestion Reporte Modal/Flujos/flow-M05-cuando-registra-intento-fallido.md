━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando registra un intento fallido
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón rojo "Intento Fallido"
               (step 'esquema3x3', solo si canManage)

Precondición: step === 'esquema3x3' Y canManage === true

INPUT: {
  icode:  reporte abierto en el modal
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Llamar al backend

  saving = true
  POST /api/v1/victim-contacts/{icode}/intento
  body: { result: 'fallido' }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/:icode/intento (nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — Validar precondiciones

  SI no hay sesión o sesión.role NOT IN ('op', 'sv') (permiso manage_victim_contacts):
    → HTTP 401   // CheckPermission responde 401, no 403
  SI body.result falta, está vacío o no es 'fallido' / 'exitoso':
    → HTTP 400
  vc = buscar victim_contact por icode
  SI vc no existe:
    → HTTP 404
  SI vc.close_reason != null O vc.victim_contact_status = 'i' (reporte legado invalidado):
    → HTTP 409 "el reporte ya está cerrado"
    → TERMINAR ejecución
  // GAP-2 confirmado 2026-10-11: op y sv pueden gestionar cualquier reporte.


PASO 3 — Insertar el intento

  INSERT INTO victim_contact_attempt (
    victim_contact_attempt_i_code, victim_contact_attempt_victim_contact,
    victim_contact_attempt_result, victim_contact_attempt_date,
    victim_contact_attempt_agent, victim_contact_attempt_creation_date
  ) VALUES (
    gen_random_uuid(), :icode, 'fallido', <ahora en UTC>, sesión.userICode, <ahora en UTC>
  )
  // <ahora en UTC> = time.Now().UTC() desde Go: la columna es timestamp sin zona y el sistema la lee
  // como UTC (con time.Now() a secas la hora quedaba en la zona del servidor).

  // No hay límite duro que RECHACE un 4° intento fallido el mismo día — el
  // mockup solo avisa la cuota (ver PASO 5), nunca bloquea el registro.


PASO 4 — Recalcular contadores (zona America/Bogota, igual que casos_nuevos)

  intentosHoy = COUNT intentos fallidos de HOY para este icode
  // día Bogotá = (fecha AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota')::date

  // "Día completo" = día con >= 3 intentos fallidos, no cualquier día con
  // intentos (corrección 2026-10-12 — ver flow-M07 PASO 2 para el detalle
  // y la fuente exacta en el mockup, diasCompletosDe()).
  fallidosPorDia = agrupar intentos fallidos de este icode POR fecha::date, contar cada grupo
  diasCompletos  = COUNT(días EN fallidosPorDia CON count >= 3)


PASO 5 — Responder

  { ok: true, intentosHoy, diasCompletos, quotaReached: intentosHoy >= 3,
    agentName: nombre de quien registra (de la sesión) }   // el modal lo usa en el log sin recargar el detalle

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — actualizar vista
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 6 — Refrescar contadores y log

  saving = false
  intentos.push({ date: ahora (ISO), result: 'fallido', agentName: response.agentName })
  intentosHoy = response.intentosHoy
  diasCompletos = response.diasCompletos

  SI response.intentosHoy >= 3 Y response.diasCompletos < 3:
    mostrar QuotaBanner (naranja) — "puede seguir intentando hoy"
  SI response.diasCompletos >= 3:
    mostrar CompleteBanner con botón "Cerrar por no contacto" (→ M-07)

  // En ningún caso se cierra el reporte automáticamente aquí — solo se
  // HABILITA la opción. El cierre real requiere la acción explícita M-07.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Resolución |
|------------------------------------------------------------------------------------|------------|
| ¿Debe existir algún límite duro (p. ej. no más de N intentos fallidos en total) además del aviso de cuota diaria? | No: sin límite duro, solo el aviso informativo (confirmado: un 4.º y 5.º fallido del mismo día se aceptan, caso `M05-C3`) |
