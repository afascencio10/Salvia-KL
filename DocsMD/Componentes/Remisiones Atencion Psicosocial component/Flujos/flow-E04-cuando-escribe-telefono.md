━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el filtro de teléfono
   Tipo: User Interaction
   Código: E-04
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en SearchFilter `telefono` — debounce 400ms

INPUT: {
  queryText:  texto ingresado  → v-model searchTelefono
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Debounce 400ms; queryText = trim(searchTelefono)

PASO 2 — Resetear paginación y selección

  currentPage        = 1
  selectedRemisiones = []
  loading            = true

PASO 3 — Actualizar activeFilters

  SI queryText === ""  → delete activeFilters['telefono']
  SI NO                → activeFilters['telefono'] = queryText

PASO 4 — fetchRemisiones()

  GET /api/v1/psychosocial-support/list
    &filter_telefono={queryText}
    (+ demás filtros activos)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — cláusula WHERE adicional
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Usar la misma expresión COALESCE de teléfono que E-01:

```sql
AND (
  COALESCE(
    NULLIF(vf2.victim_case_form2_victim_phone::text, ''),
    (SELECT vcf1.victim_contact_form1_phone::text
     FROM salvia.victim_contact_form1 vcf1
     WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
     LIMIT 1),
    ''
  ) ILIKE '%' || {filter_telefono} || '%'
)
```
