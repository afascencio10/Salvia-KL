━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por nivel de riesgo
   Tipo: User Interaction
   Código: E-12
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `nivel_riesgo`

> Filtra por el nivel de riesgo del **caso** asociado, no de la remisión.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Opciones:

| value | label |
|---|---|
| `bajo` | Bajo |
| `moderado` | Moderado |
| `alto` | Alto |
| `extremo` | Crítico |

PASO 1 — activeFilters['nivel_riesgo'] = value (o delete si "")

PASO 2 — fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Traducción UI → entero:

| filter_nivel_riesgo | victim_case_form2_risk_level |
|---|---|
| `bajo` | 1 |
| `moderado` | 2 |
| `alto` | 3 |
| `extremo` | 4 |

```sql
AND vf2.victim_case_form2_risk_level = {nivel_entero}
```

JOIN `victim_case_form2` ya presente en query base (E-01).
