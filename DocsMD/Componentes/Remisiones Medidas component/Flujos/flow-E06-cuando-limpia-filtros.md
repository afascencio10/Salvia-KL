━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario limpia los filtros
   Tipo: User Interaction
   Código: E-06
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: botón "Limpiar filtros" (vistas de gestión) o "Limpiar" (agente)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Vaciar solo los filtros de UI

  searchText     = ""
  filtroEstado   = ""
  filtroAgente   = ""
  filtroUrgencia = ""

  No cambia vista.
  No cambia agenteId.
  No vacía selectedIds. Una fila que siga visible permanece marcada.

PASO 2 — Volver a pedir E-01

  La URL queda con vista y, si es agente, con el agent_id del recorte.
  Sin q, status, agent_id de filtro ni urgency.

  El alcance no se modifica: limpiar devuelve las filas del recorte, que
  es lo que el pie ya usaba como denominador.

  → FIN EJECUCIÓN ✓
