━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el autocomplete de profesional asignada
   Tipo: User Interaction
   Código: E-07
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en AutocompleteFilter `profesional_asignada` — debounce 400ms

> **Alcance:** Solo busca agentes para sugerencias. No recarga la tabla.
> El filtrado ocurre al seleccionar (→ E-08).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — queryText = trim(autocompleteText['profesional_asignada'])

  SI queryText === "":
    → autocompleteSuggestions = []
    → FIN

  SI queryText.length < 3:
    → FIN (esperar más caracteres)

PASO 2 — Indicador de carga

  autocompleteLoading = true
  autocompleteSuggestions = []

PASO 3 — Buscar agentes

  GET /api/v1/agents/search-psicosocial?q={queryText}&limit=10

  SI error → autocompleteLoading = false; autocompleteError = true

  SI ok:
    autocompleteSuggestions = data.agents.map(a => ({
      value: a.icode,
      label: a.names + " " + a.lastNames,
      team:  a.team
    }))
    autocompleteLoading = false

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/agents/search-psicosocial
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
SELECT
    gu.general_user_i_code              AS icode,
    gup.general_user_profile_names      AS names,
    gup.general_user_profile_last_names AS last_names,
    gu.general_user_team                AS team
FROM security.general_user gu
JOIN security.general_user_profile gup
     ON gup.general_user_profile_id = gu.general_user_general_user_profile
WHERE gu.general_user_status = 'e'
  AND gu.general_user_team IN ('psicologia', 'trab. social')
  AND (
        gup.general_user_profile_names ILIKE '%' || q || '%'
     OR gup.general_user_profile_last_names ILIKE '%' || q || '%'
     OR (gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names) ILIKE '%' || q || '%'
      )
ORDER BY gup.general_user_profile_last_names, gup.general_user_profile_names
LIMIT {limit}
```

**Equipos permitidos (decisión aplicada):**

| `general_user_team` | Rol UI |
|---|---|
| `psicologia` | Psicóloga |
| `trab. social` | Trab. Social |

Mínimo 3 caracteres · debounce 400ms · máx. 10 resultados.
