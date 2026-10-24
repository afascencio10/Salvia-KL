━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Código: E-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Pantalla: Mis Remisiones — Medidas de Emergencia
Ruta: GET /salvia/mis-remisiones-medidas


INPUT: {
  session:     userData en la sesión Gin
  currentRole: rol de CommonSession
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — MisRemisionesMedidasGET
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Sesión

  Sin userData, o si GetCommonSession falla:
    → 307 a /static/landing.html
    → TERMINAR

PASO 2 — Permiso HTML

  CheckPermission(..., "get_mis_remisiones_medidas", s.CurrentRole)
  Hoy el permiso es true solo para sv.
  Si se niega, CheckPermission responde 401 y TERMINAR.

PASO 3 — Render

  Template mis_remisiones_medidas.html
  Extra: el mismo remisiones_medidas_component.html de la pantalla general
  windowTitle: "Mis Remisiones — Medidas de Emergencia"
  currentUser, currentRole, menu

  El handler no envía el icode del agente. El template lo deja fijo.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — mis_remisiones_medidas.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 4 — Montar el componente en vista agente

  <remisiones-medidas-component vista="agente" agente-id="am-1">

  Eso dispara E-01 del componente:

    GET /api/v1/emergency-measure-referrals?vista=agente&agent_id=am-1

  am-1 es temporal. El destino es no enviar agent_id: el backend usa el
  general_user_i_code de la sesión cuando exista el rol de Agente de Medidas.

  Hasta que ese endpoint exista, el componente sigue pintando el seed
  filtrado por agenteId am-1.

PASO 5 — Handlers de navegación

  onVerCaso y onVerDetalle no redirigen.
  E-14 y E-15 del componente emiten y la pantalla se queda en la bandeja.

  → FIN EJECUCIÓN ✓
