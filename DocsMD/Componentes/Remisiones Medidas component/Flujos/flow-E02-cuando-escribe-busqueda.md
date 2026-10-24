━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el buscador
   Tipo: User Interaction
   Código: E-02
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input de búsqueda (v-model searchText)


INPUT: {
  q: texto escrito
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Debounce 400ms y volver a pedir E-01

  GET /api/v1/emergency-measure-referrals
    ?vista={vista}
    &q={texto}
    + filtroEstado, filtroAgente, filtroUrgencia si están activos
    + agent_id si vista = agente

  q vacío no se envía.

  Las cards, el banner y alcance no cambian: el backend los calcula antes
  de aplicar q. Solo cambia items y, con eso, el numerador del pie.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  AND (
    nombre completo de la víctima ILIKE %q%
    OR victim_case_victim_doc_number ILIKE %q%
  )

  El nombre completo es victim_case_victim_names || ' ' || victim_case_victim_last_names.

  Se combina con AND con el recorte de la vista y con status, agent_id y urgency.

  → FIN EJECUCIÓN ✓
