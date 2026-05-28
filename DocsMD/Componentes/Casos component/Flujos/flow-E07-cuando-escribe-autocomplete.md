━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el autocomplete de persona asignada
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en el <input type="text"> del AutocompleteFilter
               (filter.key === 'persona_asignada', filter.type === 'autocomplete')
               Con debounce de 400ms — se ejecuta cuando el usuario deja de escribir

INPUT: {
  filterKey:    'persona_asignada'             → constante, identificador del filtro
  queryText:    texto escrito por el usuario   → v-model autocompleteText['persona_asignada']
}

> **Alcance:** Este evento solo busca agentes para mostrar en el dropdown de sugerencias.
> No modifica el filtro activo ni recarga los casos. La acción de filtrar
> ocurre únicamente cuando el usuario selecciona una sugerencia (→ Ver E-08).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si hay texto mínimo para buscar

  queryText = autocompleteText['persona_asignada'].trim()

  SI queryText === "":
    → autocompleteSuggestions['persona_asignada'] = []
    → autocompleteLoading['persona_asignada']     = false
    → FIN EJECUCIÓN ✓  (campo vacío — no buscar)

  SI queryText.length < 2:
    → No llamar al backend todavía (esperar más caracteres)
    → FIN EJECUCIÓN ✓

  SI queryText.length >= 2:
    → CONTINÚA PASO 2


PASO 2 — Mostrar indicador de carga y limpiar sugerencias anteriores

  autocompleteLoading['persona_asignada']     = true
  autocompleteSuggestions['persona_asignada'] = []


PASO 3 — Buscar agentes en el backend

  GET /api/v1/agents/search?q={queryText}&limit=10

  // Busca agentes cuyo nombre o apellido contenga el texto (ILIKE)
  // El límite de 10 sugerencias es suficiente para el dropdown

  SI respuesta no ok (status != 2xx):
    → autocompleteLoading['persona_asignada'] = false
    → autocompleteSuggestions['persona_asignada'] = []
    → Mostrar error silencioso en el input (borde rojo, sin bloquear)
    → TERMINAR ejecución

  SI respuesta ok:
    → agentes = data.agents   // Array<{ icode, names, lastNames, team }>
    → CONTINÚA PASO 4


PASO 4 — Mapear resultados a sugerencias y mostrar dropdown

  autocompleteSuggestions['persona_asignada'] = agentes.map(a => ({
    value: a.icode,                             // el icode es el valor del filtro
    label: a.names + " " + a.lastNames          // texto visible en el dropdown
  }))

  autocompleteLoading['persona_asignada'] = false

  SI autocompleteSuggestions['persona_asignada'].length === 0:
    → Mostrar NoResults "Sin resultados" en el dropdown

  SI autocompleteSuggestions['persona_asignada'].length > 0:
    → Mostrar SuggestionsDropdown con los resultados

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/agents/search
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  q:      texto de búsqueda    → query param
  limit:  máximo de resultados → query param (default: 10)
}

PASO 5 — Buscar agentes por nombre o apellido

  SELECT
    general_user_i_code   AS icode,
    general_user_profile_names      AS names,
    general_user_profile_last_names AS last_names,
    general_user_team     AS team
  FROM [tabla_agentes]                              // ⚠️ GAP: tabla exacta del perfil del agente
  WHERE (
      general_user_profile_names      ILIKE '%' || q || '%'
    OR general_user_profile_last_names ILIKE '%' || q || '%'
  )
  AND [campo_estado] = 'activo'                     // ⚠️ GAP: campo que indica si el agente está activo
  ORDER BY general_user_profile_last_names ASC
  LIMIT limit

PASO 6 — Retornar resultados

  200 {
    agents: [
      {
        icode:     general_user_i_code,
        names:     general_user_profile_names,
        lastNames: general_user_profile_last_names,
        team:      general_user_team
      },
      ...
    ]
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| Tabla exacta del perfil del agente (¿security.general_user_profile?)             | PASO 5        |
| Campo que indica si el agente está activo/inactivo en esa tabla                  | PASO 5        |
| ¿Se filtra por rol? ¿Solo agentes con rol de profesional de caso?                | PASO 5        |
| Endpoint exacto para buscar agentes (puede diferir de /api/v1/agents/search)     | PASO 3        |
| Mínimo de caracteres para disparar la búsqueda: ¿2 es correcto o se prefiere 3? | PASO 1        |
| ¿El debounce de 400ms es correcto para el autocomplete?                          | PASO 1        |
