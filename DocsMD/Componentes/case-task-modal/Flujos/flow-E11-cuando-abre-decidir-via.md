━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre decidir vía de riesgo inminente
   Tipo: User Interaction
   Función: open(taskId) + _initForm()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: "Gestionar" en case-tasks, cuando la tarea es decidir_via_riesgo_inminente.
La carga de la tarea es E02. Este flujo es lo que pinta el modal después.

INPUT: {
  tarea.type: "decidir_via_riesgo_inminente"
  tarea.assignedRole: "sm"
}


PASO 1 — Título y subtítulo

  Título:    "Decidir vía de atención — Riesgo inminente"
  Subtítulo: "Elige una vía y asigna un agente de medidas."


PASO 2 — Formulario local

  form.via = ""
  form.agentId = ""

  Vía A y Vía B son radios. Solo una puede quedar marcada.

  Vía A — Asignación inmediata a Agente de Medidas
    El agente prioriza contacto y entrevista de caracterización de inmediato.

  Vía B — Derivar a Supervisión del Convenio
    Escala a Cruz Roja para entrega inmediata (hotel/transporte), sin esperar caracterización.

  Asignar agente carga GET /api/v1/emergency-measure-referrals/agentes.
  El catálogo son los usuarios activos con rol am.


PASO 3 — Cerrar sin guardar

  Cancelar, la X o el overlay llaman cancelar().
  No escriben form_data ni cambian la tarea.


PASO 4 — Confirmar decisión

  El botón se habilita con una vía y un agente.
  El guardado está en E09.
  Vía A crea intentar_contacto_3x3 para el agente.
  Vía B crea escalar_cruz_roja para el rol sc.
  Esos dos modales se abren después, desde su propia tarea.
  Guardar y Confirmar de esos modales cierran y no persisten.
