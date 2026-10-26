━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Código: E-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Pantalla: Detalle de Remisión — Medida de Emergencia
Ruta: GET /salvia/remision-medida/:id


INPUT: {
  session:     userData en la sesión Gin
  currentRole: rol de CommonSession
  id:          uuid de emergency_measure_referral, el remisionId de E-15
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — RemisionMedidaDetalleGET
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Sesión

  Sin userData, o si GetCommonSession falla:
    → 307 a /static/landing.html
    → TERMINAR

PASO 2 — Permiso HTML

  CheckPermission(..., "get_remision_medida_detalle", s.CurrentRole)
  El permiso es true para sm, em, sc, cr, sv y am.
  Si se niega, CheckPermission responde 401 y TERMINAR.

PASO 3 — Render

  Template remision_medida_detalle.html
  windowTitle: "Detalle de Remisión — Medida de Emergencia"
  remisionId, currentUser, currentRole, currentUserId, menu

  salviaBackFallbackPath:
    am → /salvia/mis-remisiones-medidas
    el resto → /salvia/remisiones-medidas

  El handler no lee la remisión. El template la pide en el paso 5.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — remision_medida_detalle.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 4 — Estados

  Al montar: cargando = true, tab = info.
  Si remisionId viene vacío: error, sin pedir el API.

PASO 5 — Leer

  GET /api/v1/emergency-measure-referrals/{remisionId}
  credentials: same-origin

  El API exige el mismo permiso y la sesión.
  VistaRemisionesMedidas(currentRole) elige el recorte.
  Si el rol no participa, 403.
  La consulta usa el scopeClause del listado y filtra por id.

  sm y sv: cualquier remisión no borrada.
  am: assigned_agent_id = icode de la sesión, y status en los cinco
       estados de su bandeja.
  em: status = pendiente_de_revision.
  sc: status en por_revisar_trazabilidad, por_aprobar_cotizaciones
       o riesgo_inminente_por_escalar_cruz_roja.
  cr: status en aprobadas, por_cotizar
       o riesgo_inminente_por_entregar_medida.

  Sin fila (no existe, está borrada o queda fuera del recorte): 404.
  El frontend pinta el error y el link de volver. No distingue esos tres casos.

  200 devuelve la cabecera del listado más:
    reason, rejectionReason, submittedByTeam,
    childrenNumber, childrenAges, requiresInterpreter.

  requiresInterpreter es true solo si el enum del form 2 tiene código y.
  childrenNumber es null cuando el form 1 no trae número.

PASO 6 — Pintar

  Cabecera, chips y tab Info Medida con esos campos.
  "Agresor activo" se deriva de urgency, no de una columna.
  La descripción y las observaciones del mockup no se pintan.

  Tareas, al activar el tab o al montar el componente:
    <case-tasks :case-id="caseICode" :user-id="currentUserId" :user-role="currentRole">

  Timeline:
    <case-timeline :case-id="caseICode">

  Los dos piden el caso completo. No reciben el id de la remisión.

  El bloque Asignación se pinta solo si currentRole es sm o sv.
  Abrirlo y confirmarlo es E-02.

  → FIN EJECUCIÓN ✓
