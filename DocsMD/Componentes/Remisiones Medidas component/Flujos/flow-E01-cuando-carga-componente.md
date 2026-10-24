━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Código: E-01
   Prerequisito: M-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: montar `remisiones-medidas-component`

Hoy el `mounted` no llama al API: clona el seed y calcula alcance, cards y filas en el cliente. Este flujo es el contrato al que se cablea.


INPUT: {
  vista:      recorte y presentación     → prop vista (default supervisor)
  agenteId:   icode del agente           → prop agenteId (default am-1)
              solo lo usa vista = agente
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — remisiones-medidas-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Estado inicial de los filtros de UI

  searchText     = ""
  filtroEstado   = ""
  filtroAgente   = ""
  filtroUrgencia = ""
  selectedIds    = []

  La vista y agenteId no son filtros de UI. No se vacían en E-06.


PASO 2 — Pedir el listado del recorte

  GET /api/v1/emergency-measure-referrals
    ?vista={vista}

  SI vista === "agente":
    → agregar &agent_id={agenteId}

  agent_id es temporal. Cuando exista el rol de Agente de Medidas, Mis
  Remisiones deja de enviar el id y el backend lo toma de la sesión
  (general_user_i_code). Un agent_id distinto al de la sesión se rechaza.

  vista también es temporal. Cuando cada vista tenga rol, el backend la
  deduce de la sesión. El query vista deja de ser la fuente del recorte.

  Sin page ni page_size. La pantalla no pagina.

  Los filtros de UI (E-02 … E-05) se agregan a esta misma URL. En la carga
  inicial van vacíos.


PASO 3 — Recorte que aplica el backend (antes de los filtros de UI)

  deleted_at IS NULL en la remisión.

  vista = supervisor
    → todas las remisiones

  vista = equipo-tecnico
    → status = pendiente_de_revision

  vista = supervision-convenio
    → status IN (
         por_revisar_trazabilidad,
         por_aprobar_cotizaciones,
         riesgo_inminente_por_escalar_cruz_roja
       )

  vista = cruz-roja
    → status IN (
         aprobadas,
         por_cotizar,
         riesgo_inminente_por_entregar_medida
       )

  vista = agente
    → assigned_agent_id = agent_id
      AND status IN (
           riesgo_inminente_por_confirmar_seguridad,
           asignada,
           en_gestion_documental,
           aprobadas,
           en_seguimiento
         )

  Ese conjunto es el alcance. Las cards, el banner y el denominador del pie
  salen de ahí. Búsqueda, estado, agente y urgencia recortan solo `items`.


PASO 4 — Respuesta

  {
    "alcance": 15,
    "cards": [ { "key": "total", "label": "Total", "value": 15 } ],
    "banner": null,
    "items": [ ReferralListItem ]
  }

  supervisor → cards con las 8 métricas de la interfaz, banner null
  agente    → cards con las 6 métricas, banner null
  equipo-tecnico, supervision-convenio, cruz-roja
            → cards null, banner { "tono", "texto" } con el copy de la interfaz

  ReferralListItem:

  {
    "id": "uuid",
    "caseICode": "SAL-001",
    "caseName": "Ana Torres Velásquez",
    "docNumber": "1098234567",
    "riskLabel": "Alto",
    "municipality": "Cali",
    "measures": ["Alojamiento temporal", "Kit de alimentación"],
    "urgency": "riesgo_inminente",
    "urgencyLabel": "Riesgo inminente",
    "submittedByName": "María González",
    "referredAt": "2026-06-05",
    "status": "caracterizado",
    "statusLabel": "Caracterizado",
    "assignedAgent": { "id": "icode", "nombre": "Pedro Ramírez" },
    "technicalReview": {
      "professionalName": "Camila Andrea Rojas Pinzón",
      "concept": "...",
      "result": "aprobado_para_caracterizacion_me",
      "resultLabel": "Aprobado para caracterización ME",
      "date": "2026-06-11"
    }
  }

  assignedAgent y technicalReview van null cuando no hay dato.
  measures sale de emergency_measure.type WHERE referral_id = id
  ORDER BY created_at.
  riskLabel usa el mapa de case_info_service (form 2: 1 Bajo, 2 Moderado,
  3 Alto, 4 Extremo).
  municipality es town_name del municipio de la víctima.
  submittedByName es el nombre del general_user de submitted_by.
  El listado no pinta reason ni rejection_reason.

  Orden de items:

    ORDER BY (urgency = 'riesgo_inminente') DESC, referred_at DESC


PASO 5 — Pintar

  SI cards tiene elementos → fila de cards (8 o 6). No son botones.
  SI banner no es null     → banner. No hay cards.
  Tabla con items.
  Pie: "{items.length} de {alcance} remisiones".

  SI vista = agente AND alcance = 0:
    → "No tienes remisiones asignadas."
  SI items está vacío y alcance > 0:
    → "No hay remisiones que coincidan con los filtros."

  SI la petición falla:
    → la tabla queda vacía y se muestra el mensaje de error de la respuesta
    → no se inventan filas del seed


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — lectura
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Handler nuevo. No reutiliza el List de EmergencyMeasure: ese lista medidas
  sueltas, sin cabecera ni los estados de esta bandeja.

  Joins de cada item, alineados al listado de Atención Psicosocial:

    salvia.victim_case            por victim_case_i_code = case_id
    salvia.victim_case_form2      por el caso, para risk_level
    directorio de municipio       town_name
    security.general_user         submitted_by, assigned_agent_id,
                                  review_professional_id
    security.general_user_profile nombres y apellidos
    salvia.emergency_measure      type por referral_id

  Las cards son COUNT sobre el alcance, con los mismos criterios de la
  interfaz. No usan los filtros de UI.

  → FIN EJECUCIÓN ✓
