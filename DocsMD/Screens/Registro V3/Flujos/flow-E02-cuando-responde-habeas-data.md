━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando responde la autorización de datos personales
   Tipo: User Interaction
   Función: responderHabeasData(autoriza) → components/habeas-data-gate.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  autoriza: boolean   → click en BtnSí / BtnNo del HabeasDataGate
}

PASO 1 — Evaluar la respuesta

SI autoriza === true:
  → stage = 'tipoContacto'
  → Vue re-evalúa stage → se oculta HabeasDataGate, se muestra TipoContactoBlock
  → CONTINÚA FLUJO GENERAL

SI NO (autoriza === false):
  → stage = 'formularioCerrado'
  → No se crea ningún `victim_case`, no se guarda ninguna respuesta
  → Vue re-evalúa stage → se muestra MensajeCierre
  → TERMINAR ejecución — no hay nada más que hacer en esta pantalla

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                          | Paso afectado |
|-----------------------------------------------------------------|---------------|
| ¿Se debe registrar en algún log/analítica el número de personas que responden "No" a habeas data? (el Excel no lo pide explícitamente) | PASO 1 (rama No) |
