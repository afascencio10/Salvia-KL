━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario filtra por urgencia
   Tipo: User Interaction
   Código: E-05
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: select "Urgencia"
Visible en las cinco vistas.


INPUT: {
  urgency: "riesgo_inminente", "normal", o vacío
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Opciones

  vacío              label "Urgencia"
  riesgo_inminente   label "Riesgo inminente"
  normal             label "Normal"

PASO 2 — Volver a pedir E-01 con &urgency={valor}

  La opción vacía quita el parámetro.
  Cards, banner y alcance no cambian.
  El orden de E-01 se mantiene dentro del resultado: si se pidió
  riesgo_inminente, todas las filas quedan en ese grupo.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  AND referral.urgency = {urgency}

  Otro valor responde 400.

  → FIN EJECUCIÓN ✓
