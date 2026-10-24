━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Código: E-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Pantalla: Remisiones Medidas de Emergencia
Ruta: GET /salvia/remisiones-medidas


INPUT: {
  session:     userData en la sesión Gin
  currentRole: rol de CommonSession
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — RemisionesMedidasGET
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Sesión

  Sin userData, o si GetCommonSession falla:
    → 307 a /static/landing.html
    → TERMINAR

PASO 2 — Permiso HTML

  CheckPermission(..., "get_remisiones_medidas", s.CurrentRole)
  El permiso es true para sm, em, sc, cr y sv.
  Si se niega, CheckPermission responde 401 y TERMINAR.

PASO 3 — Render

  Template remisiones_medidas.html
  Extra: remisiones_medidas_component.html
  windowTitle: "Remisiones Medidas de Emergencia"
  currentUser, currentRole, menu


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — remisiones_medidas.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 4 — Montar el componente según el rol

  sm y sv → vista="supervisor" y :puede-reasignar="true"
  em  → vista="equipo-tecnico"
  sc  → vista="supervision-convenio"
  cr  → vista="cruz-roja"

  Eso dispara E-01 del componente:
    GET /api/v1/emergency-measure-referrals?vista={vista}

  Hasta que ese endpoint exista, el componente sigue pintando el seed.
  El padre no envía agenteId.

  Al iniciar sesión, sm, em, sc y cr aterrizan en esta ruta.
  sv conserva su entrada habitual y llega aquí por el menú.

PASO 5 — Handlers de navegación

  onVerCaso y onVerDetalle no redirigen.
  E-14 y E-15 del componente emiten y la pantalla se queda en el listado.

  → FIN EJECUCIÓN ✓
