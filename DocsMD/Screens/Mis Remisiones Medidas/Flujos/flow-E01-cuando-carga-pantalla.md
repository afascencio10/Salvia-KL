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
  El permiso es true para am y ro.
  Si se niega, CheckPermission responde 401 y TERMINAR.

PASO 3 — Render

  Template mis_remisiones_medidas.html
  Extra: el mismo remisiones_medidas_component.html de la pantalla general
  windowTitle: "Mis Remisiones — Medidas de Emergencia"
  currentUser, currentRole, menu

  El handler pasa agentICode = s.UserICode. El template lo pone en agente-id.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — mis_remisiones_medidas.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 4 — Montar el componente en vista agente

  <remisiones-medidas-component vista="agente" agente-id="{UserICode}">

  Eso dispara E-01 del componente:

    GET /api/v1/emergency-measure-referrals?vista=agente&agent_id={UserICode}

  El backend de esa vista usa el icode de la sesión y rechaza otro agent_id.

  Hasta que ese endpoint exista, el componente sigue pintando el seed
  filtrado por el icode real. Los ids del seed (am-1, am-2, am-3) no
  coinciden, así que la bandeja abre vacía.

  Al iniciar sesión, am aterriza en esta ruta.

PASO 5 — Handlers de navegación

  onVerDetalle redirige a /salvia/remision-medida/{remisionId}.
  Es E-15 del componente. La pantalla de detalle es de lectura.

  onVerCaso no redirige. E-14 emite y la bandeja se queda.

  → FIN EJECUCIÓN ✓
