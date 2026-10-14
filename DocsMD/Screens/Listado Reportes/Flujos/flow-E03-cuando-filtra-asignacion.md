━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por asignación
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: (A) escribe en el autocomplete "Asignación"
               (B) selecciona la opción fija "Sin asignar" (siempre primera
                   del dropdown de sugerencias, antes de cualquier resultado
                   de búsqueda de agentes)

INPUT: {
  text:          texto escrito           → autocompleteText['asignacion']
  opcionFija:    true si eligió "Sin asignar"
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js (extensión)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Rama "Sin asignar"

  SI usuario selecciona la opción fija "Sin asignar":
    autocompleteSelected['asignacion'] = { label: 'Sin asignar', value: '__unassigned__' }
    activeFilter = { key: 'asignacion', value: '__unassigned__' }
    currentPage = 1
    → GET {listEndpoint}?filter_key=asignacion&filter_value=__unassigned__&page=1...
    → TERMINAR ejecución (no vuelve a pedir la lista de agentes)


PASO 2 — Rama búsqueda de agente (mismo mecanismo de E-07/E-08, con otra fuente de datos)

  SI NO (usuario escribe texto o abre el desplegable):
    onAutocompleteInput('asignacion')  [debounce 400ms] → runAsignacionSearch('asignacion')
    SI aún no hay caché de agentes (_opAgentsCache):
      → GET /api/v1/agents/by-role?role=op     // agentes integrales; se pide UNA vez y se guarda en caché
    sugerencias = [ { value:'__unassigned__', label:'Sin asignar' } ]
                  + agentes cuyo nombre contiene el texto (filtrado en el cliente, sin distinguir mayúsculas)
    // No se usa /api/v1/agents/search: ese endpoint filtra por equipo (no por rol) y un reporte no
    // tiene equipo hasta que se le asigna un agente (GAP-3).

  SI usuario selecciona una sugerencia:
    activeFilter = { key: 'asignacion', value: agente.icode }
    → GET {listEndpoint}?filter_key=asignacion&filter_value={agente.icode}&page=1...


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/victim-contacts/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 3 — Traducir filter_value

  SI filter_key === 'asignacion' Y filter_value === '__unassigned__':
    WHERE victim_contact_agent_id IS NULL
  SI filter_key === 'asignacion' Y filter_value != '__unassigned__':
    WHERE victim_contact_agent_id = {filter_value}

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                 | Resolución |
|---------------------------------------------------------------------------------------|------------|
| ¿La opción "Sin asignar" se muestra siempre fija arriba, o solo cuando el texto está vacío? | Siempre es la primera opción del desplegable, también al filtrar por texto (caso `LR-E03-C4`) |
| `/api/v1/agents/search` hoy busca sin filtro de rol explícito — este filtro necesita acotar a `op` | No se usa: la lista sale de `GET /api/v1/agents/by-role?role=op` (solo agentes `op`; caso `LR-E03-C5`) |
