━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el formulario emite "form-completed"
   Tipo: User Interaction
   Función: onFormCompleted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  (ninguno — 'form-completed' no lleva payload, igual que en Seguimiento)
  submissionId:  UUID del FormSubmission   → estado interno de la pantalla
}

Mecanismo resuelto (decisión del usuario): usar el mismo criterio que ya
usa `dinamic-form` internamente para emitir este evento — todas las
secciones visibles con `isAnswered == true` — y, en ese preciso momento,
hacer un GET puntual al backend para obtener el resultado de
`processVictimCaseSubmission` (E-04). No hay polling con reintentos en el
cliente: el endpoint espera internamente (server-side) a que el goroutine
termine, con un timeout corto.


PASO 1 — Marcar el formulario como completado

  formCompleted = true
  → Monta CompletedOverlay con estado "cargando credenciales..."


PASO 2 — Solicitar el resultado al backend

  GET /api/v1/victim-case-forms/{submissionId}/result

  ┌──────────────────────────────────────────────────────────────┐
  │  BACKEND — GET .../result                                     │
  └──────────────────────────────────────────────────────────────┘

    El endpoint espera brevemente a que processVictimCaseSubmission (E-04,
    goroutine ya en curso) termine de escribir en form_submission_result —
    evita el round-trip de reintentos en el cliente ante la condición de
    carrera entre "el saveSection ya respondió" y "el goroutine aún no
    terminó" (bcrypt + varios inserts ~200-500ms típico).

    PARA cada intento (hasta 5, con 300ms de espera entre intentos, ~1.5s máx):
      DB.form_submission_result.FindBySubmissionID({ submissionId })
      SI existe:
        → Responder 200 { caseId, newUser: { login, pass } }
        → TERMINAR
      SI NO existe:
        → Esperar 300ms → siguiente intento

    SI se agotan los 5 intentos sin resultado:
      → Responder 202 { status: "processing" }
      // El caso puede seguir procesándose; el frontend debe manejar esto
      // como fallback, no como error


PASO 3 — Mostrar credenciales en el overlay

  SI respuesta 200:
    → Mostrar Usuario: newUser.login  /  Clave: newUser.pass
    → caseId = respuesta.caseId
    → Habilitar BtnFinalizar (ver E-06)

  SI respuesta 202 (aún procesando):
    → Mostrar: "El caso se está registrando, esto puede tardar unos segundos..."
    → Reintentar el GET una sola vez más tras 2s
    → SI sigue en 202:
        → Mostrar: "El caso se registró correctamente, pero hubo un problema
          mostrando las credenciales — consúltalas en el detalle del caso"
        → Habilitar BtnFinalizar igual, con redirección a /salvia/lista-casos
          en vez de al detalle (no hay caseId confirmado — ver E-06)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                          | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| Tabla `form_submission_result` — mismo gap de flow-E04 PASO 11 (necesita expiración/limpieza) | PASO 2 |
| Afinar timeouts exactos (300ms×5, +2s de reintento extra) contra tiempos reales medidos en producción | PASO 2-3 |
