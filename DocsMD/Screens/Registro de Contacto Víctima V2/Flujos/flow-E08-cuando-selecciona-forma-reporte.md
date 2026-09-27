━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona la forma del reporte (anonimato)
   Tipo: User Interaction
   Función: v-model="form.anonimato" (paso 'formulario', solo tipoReporte === 'Tercero')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  anonimato:  valor seleccionado ('Anonimo' | 'Propio')  → <select> del usuario
}

PASO 1 — Actualizar form.anonimato con el valor seleccionado

PASO 2 — Vue re-evalúa el computed mostrarNombreReportante
  SI form.anonimato === 'Propio':
    → mostrarNombreReportante = true
    → Se muestran los campos Nombres, Apellidos y Correo del reportante
    → CONTINÚA FLUJO GENERAL

  SI form.anonimato === 'Anonimo':
    → mostrarNombreReportante = false
    → Los campos de identidad del reportante se ocultan; si ya tenían valor,
      permanecen en form pero no se envían en E-10 (ver PASO de armado de
      payload en flow-E10)
    → CONTINÚA FLUJO GENERAL

No hay llamada a backend. Este es el mismo evento que la reactividad del
formulario — no se documenta el re-render como un evento aparte (criterio de
identifyTriggers.md sobre reactividad de framework).

form.anonimato se guarda en E-10 en `victim_contact_form2_v2.anonimato`, campo
nuevo y dedicado (decisión 5B) — no reutiliza `victim_case_form2_report_type_details`.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|--------------------------------------------------------------------|---------------|
| Si un reporte "Anónimo" con nombre ya digitado debe limpiar esos campos al cambiar a Anónimo, o solo dejar de enviarlos | PASO 2 |
