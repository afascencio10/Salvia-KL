━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario filtra por estado
   Tipo: User Interaction
   Código: E-03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: select "Estado remisión"
Visible si vista es supervisor o agente. Equipo técnico, convenio y Cruz Roja
no lo renderizan: su recorte ya es la bandeja (E-01).


INPUT: {
  status: código snake_case, o vacío para "Estado remisión"
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Opciones

  supervisor: los 16 estados de M-01, label visible y código en el value
  agente: solo estos cinco
    riesgo_inminente_por_confirmar_seguridad
    asignada
    en_gestion_documental
    aprobadas
    en_seguimiento

PASO 2 — Volver a pedir E-01 con &status={código}

  La opción vacía quita el parámetro.
  Cards, banner y alcance no cambian.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  AND referral.status = {status}

  Si el código no está en el catálogo de M-01, la petición responde 400.
  En vista agente, un status fuera de los cinco de su bandeja también
  responde 400: no amplía el recorte.

  → FIN EJECUCIÓN ✓
