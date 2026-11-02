━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando responde "¿Tiene usuario?"
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — WP-1 T8 (2026-11-02) (ver QA/E12-tiene-usuario-qa-cases.md)
   Funciones: elegirUsuario(opcion) · goTo('auth') · window.location.assign(login)   (set_victim_contact_v2.html)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  opcion:        'si' | 'no'                  → botones «Sí, iniciar sesión» / «No, continuar sin usuario»
  reporterName:  nombre de quien tiene sesión → variable de plantilla (E-01); solo existe para los roles en, op y et
}

El paso 'usuario' es el paso inicial del wizard (E-01). Con una sesión de otro rol (sv, ad, ro…) el wizard arranca en
'auth', sin este paso ni mensaje (el reporte igual queda asociado a su usuario — ver E-10 PASO 1b).

SI reporterName ≠ "":                           // variante «Está reportando como …»
  → Solo existe el botón «Continuar» → goTo('auth')
  → CONTINÚA FLUJO GENERAL

SI NO:
  SEGÚN opcion:
    CASO 'no': → goTo('auth')                    // flujo anónimo, igual que antes
    CASO 'si': → window.location.assign('/seguridad/login?next=' + encodeURIComponent('/salvia/public-v2/nuevo'))
               → TERMINAR ejecución               // el login devuelve a esta página (Login E-02, lista blanca exacta)

PASO 2 — Al volver del login la página se carga de nuevo (E-01): ahora trae la sesión y muestra la variante con sesión.

No hay llamada a backend desde este evento. El envío del reporte (E-10) sigue exigiendo captcha aunque haya sesión.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Con una sesión de otro rol el reporte queda asociado sin que la persona lo vea (no hay mensaje): decisión de la revisión de WP-1 | PASO inicial |
