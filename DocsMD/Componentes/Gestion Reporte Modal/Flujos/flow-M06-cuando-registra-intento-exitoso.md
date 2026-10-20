━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando registra un intento exitoso
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón verde "Intento Exitoso"
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
  body: { result: 'exitoso' }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/:icode/intento
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — Validar precondiciones

  Idéntico a M-05 PASO 2 (reporte existe, no cerrado, canManage).


PASO 3 — Insertar el intento

  INSERT INTO victim_contact_attempt (...)
  VALUES (gen_random_uuid(), :icode, 'exitoso', <ahora en UTC>, sesión.userICode, <ahora en UTC>)

  // Un intento exitoso NO se cuenta para intentosHoy/diasCompletos del
  // esquema 3x3 (esos contadores son exclusivamente de intentos 'fallido') —
  // una vez hay contacto, el esquema 3x3 deja de ser relevante.
  //
  // Por qué se persiste (no es solo bitácora): si el agente cierra el modal
  // antes de completar la decisión post-contacto (M-08/M-09), reabrir el
  // reporte necesita saber que ya hubo contacto para aterrizar en
  // step='postContacto' — ver flow-M01 PASO 7. Sin esta fila, el modal no
  // tendría cómo distinguir "nunca contactó" de "contactó pero no terminó
  // de decidir", y el agente podría terminar cerrando por "no contacto" un
  // reporte donde sí hubo contacto.


PASO 4 — Responder

  { ok: true, intentosHoy, diasCompletos, quotaReached, agentName }   // mismo handler que M-05; el cliente solo usa agentName

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — transición automática
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 5 — Pasar a vista post-contacto

  saving = false
  intentos.push({ date: ahora (ISO), result: 'exitoso', agentName: response.agentName })

  esperar 400ms   // mismo retardo visual que el mockup (setTimeout 400ms),
                  // da tiempo a que el agente vea el intento registrado
                  // en el log antes de que cambie la vista
  step = 'postContacto'
  pendingClose = ''   // ninguna de las 2 opciones elegida aún (competencia / decisión — ver nota en gestion-reporte-modal-interface.md sobre M-10 diferido)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|------------------------------------------------------------------------------------|---------------|
| ¿Un intento exitoso registrado por error se puede "deshacer" (volver a esquema3x3)? El mockup no lo contempla — es unidireccional. | PASO 5 — se asume que no, consistente con el mockup |
