━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el autocomplete de persona asignada
   Tipo: User Interaction
   Código: E-07
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
    → autocompleteError['persona_asignada']         = false
    → FIN EJECUCIÓN ✓  (campo vacío — no buscar)

  SI queryText.length < 3:
    → No llamar al backend todavía (esperar más caracteres)
    → FIN EJECUCIÓN ✓

  SI queryText.length >= 3:
    → CONTINÚA PASO 2


PASO 2 — Mostrar indicador de carga y limpiar sugerencias anteriores

  autocompleteLoading['persona_asignada']     = true
  autocompleteSuggestions['persona_asignada'] = []
  autocompleteError['persona_asignada']       = false


PASO 3 — Buscar agentes en el backend

  GET /api/v1/agents/search?q={queryText}&limit=10

  // Busca agentes cuyo nombre o apellido contenga el texto (ILIKE)
  // Mínimo 3 caracteres; debounce 400ms

  SI respuesta no ok (status != 2xx):
    → autocompleteLoading['persona_asignada'] = false
    → autocompleteSuggestions['persona_asignada'] = []
    → autocompleteError['persona_asignada'] = true  (borde rojo en input)
    → TERMINAR ejecución

  SI respuesta ok:
    → agentes = data.agents   // Array<{ icode, names, lastNames, team }>
    → CONTINÚA PASO 4


PASO 4 — Mapear resultados a sugerencias y mostrar dropdown

  autocompleteSuggestions['persona_asignada'] = agentes.map(a => ({
    value: a.icode,                             // general_user_i_code
    label: a.names + " " + a.lastNames          // perfil visible en dropdown
  }))

  autocompleteLoading['persona_asignada'] = false

  SI autocompleteSuggestions['persona_asignada'].length === 0:
    → Mostrar NoResults "Sin resultados" en el dropdown (solo si queryText.length >= 3)

  SI autocompleteSuggestions['persona_asignada'].length > 0:
    → Mostrar SuggestionsDropdown con los resultados

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/agents/search
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  q:      texto de búsqueda    → query param (mín. 3 caracteres)
  limit:  máximo de resultados → query param (default: 10, máx. 50)
}

PASO 5 — Buscar agentes por nombre o apellido

  SELECT
    gu.general_user_i_code            AS icode,
    gup.general_user_profile_names    AS names,
    gup.general_user_profile_last_names AS last_names,
    gu.general_user_team              AS team
  FROM security.general_user gu
  JOIN security.general_user_profile gup
       ON gup.general_user_profile_id = gu.general_user_general_user_profile
  WHERE gu.general_user_status = 'e'
    AND gu.general_user_team IN ('Riesgo bajo', 'Riesgo alto', 'Riesgo Alto', 'Hombres')
    AND (
        gup.general_user_profile_names ILIKE '%' || q || '%'
     OR gup.general_user_profile_last_names ILIKE '%' || q || '%'
     OR (gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names) ILIKE '%' || q || '%'
    )
  ORDER BY gup.general_user_profile_last_names ASC, gup.general_user_profile_names ASC
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
✅  Decisiones aplicadas (E-07 implementado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Tabla de nombres | `security.general_user_profile` (JOIN con `general_user`) |
| ID del agente | `security.general_user.general_user_i_code` |
| Estado activo | `general_user_status = 'e'` |
| Equipos permitidos | `Riesgo bajo`, `Riesgo alto`, `Riesgo Alto`, `Hombres` |
| Endpoint | `GET /api/v1/agents/search` |
| Mínimo de caracteres | 3 |
| Debounce | 400ms |
