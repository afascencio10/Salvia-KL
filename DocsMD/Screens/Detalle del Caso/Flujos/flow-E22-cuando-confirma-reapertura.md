━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
EVENTO: Cuando confirma reapertura
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: click en "Confirmar" del modal Reabrir caso (E-21)
Archivo:       src/frontend/html/salvia/case_detail/get_case_detail_sv.html → confirmarReabrirCaso()
Backend:       POST /api/v1/casos/:id/reabrir → CasoCierreService.ReabrirCaso()

INPUT: {
  caseICode:  i_code del caso           → CASE_ICODE (plantilla Go)
  actorId:    i_code del usuario        → userICode (plantilla Go)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — confirmarReabrirCaso
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Marcar loading

  reabriendo = true


PASO 2 — Llamar al backend

  POST /api/v1/casos/{caseICode}/reabrir
  body: { actor_id: actorId }

  SI respuesta no ok:
    → reabriendo = false
    → showToast('error', data.error || 'Error al reabrir el caso')
    → TERMINAR (el modal permanece abierto)

  SI error de red:
    → reabriendo = false
    → showToast('error', 'Error de conexión al reabrir el caso')
    → TERMINAR


PASO 3 — Cerrar modal y recargar

  modalReabrir = false
  reabriendo = false
  showToast('success', 'Caso reabierto correctamente')
  cargarDatos()   // badge pasa a Activo; el botón Reabrir se oculta (status != cd)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — ReabrirCaso
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO B1 — Validar caseICode

  SI caseICode vacío:
    → Retornar error


PASO B2 — Cargar el caso

  DB.victim_case.FindByICode({ iCode: caseICode })

  SI no existe:
    → Retornar error

  SI status != "cd":
    → Log: "no está cerrado — se omite"
    → Retornar nil (idempotente)


PASO B3 — Actualizar status a Activo

  DB.victim_case.UpdateStatus({ iCode: caseICode, status: "ra" })

  SI falla el update:
    → Retornar error (aborta)


PASO B4 — Registrar evento en el timeline

  SI caseTimelineRepo es nil:
    → Log WARN
    → Retornar nil (el status ya cambió)

  DB.case_timeline_events.Create({
    case_id:       caseICode,
    category:      "General",
    type:          "Reapertura de Caso",
    icon:          "arrows-rotate",
    color:         "#0f766e",
    description:   "Caso reabierto",
    event_user_id: actorId,
    date:          now(),
  })

  SI falla el insert:
    → Log WARN (no aborta — el caso ya está reabierto)

  → FIN EJECUCIÓN


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| No se restaura el status previo al cierre (siempre pasa a `ra`)  | B3            |
| No se pide motivo de reapertura                                  | Modal / B4    |
