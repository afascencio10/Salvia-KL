━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Alerta inminente"
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — no hay botón: es la pregunta GS_18 del módulo Gestión Salvia
   Funciones: setBool('GS_18', …) · setVal('GS_19', …) · guardarModulo('GS') → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Decisión C7 (`registro-v3-data.md` § D10): la Alerta Inminente **no** es un botón del header; vive solo dentro de
> GS § E y, como GS exige la tríada completa, no puede saltarse la habilitación secuencial. El botón "⚠ Alerta
> inminente" del header y la bitácora/notificación a coordinación de la planeación **no se construyeron**.

INPUT: {
  GS_18:  radio "¿Activar alerta inminente? Sí | No"   → GS § E
  GS_19:  texto "Justificación de la alerta"             → <textarea>, visible solo si GS_18 = Sí
}

PASO 1 — Capturar la decisión
  setBool('GS_18', true|false) → answers['GS_18'] = 'true' | 'false'
  SI GS_18 = 'true': Vue muestra el campo de justificación (GS_19); no es obligatorio

PASO 2 — Persistir
  Al presionar "Guardar sección" (flow-E12, módulo GS):
    victim_case_form2_wants_imminent_alert = id de "Sí"/"No"     (bool_enum)
    victim_case_form2_imminent_alert_note  = GS_19                (texto)
  GS_18 sí es obligatoria para finalizar (flow-E15); GS_19 no.

PASO 3 — Efecto
  Ninguno más: no se crea entrada de bitácora, no se notifica a coordinación y ninguna otra pantalla lee estas dos columnas (verificado por búsqueda en el código).

Aviso relacionado: en TZ, si el nivel de riesgo calculado es alto o extremo se muestra "⚠️ Nivel de riesgo alto — considere activar la Alerta Inminente en Gestión Salvia." (flow-E11).

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Destino real de la alerta (notificación a coordinación): hoy solo queda el dato guardado | PASO 3 |
| G-V3-14 — si "Tentativa Feminicidio/Homicidio" debe habilitar TZ sin tríada completa: no implementado | — |
