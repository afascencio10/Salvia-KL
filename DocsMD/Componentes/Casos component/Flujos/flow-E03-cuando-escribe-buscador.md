━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el buscador
   Tipo: User Interaction
   Código: E-03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en el SearchFilter (type='search') dentro de FilterGroup
               Se activa en cada cambio del campo — con debounce de 400ms

INPUT: {
  searchText:   texto escrito por el usuario   → v-model del <input type="text">
}

> **Nota:** La búsqueda llama al backend porque el componente usa paginación server-side.
> El parámetro `search` se combina en paralelo con chip, dropdowns, persona asignada y agentId.

> **Campos buscados en BD:**
> - **Documento:** `victim_case.victim_case_victim_doc_number` — ILIKE parcial
> - **Teléfono:** misma fuente que la columna Víctima (E-01):
>   1. `victim_case_form2.victim_case_form2_victim_phone` (teléfono copiado al crear el caso)
>   2. Fallback: `victim_contact_form1.victim_contact_form1_phone` vía FK `victim_case_victim_contact`


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Actualizar estado y resetear paginación (debounce 400ms)

  searchText  = searchText.trim()
  currentPage = 1
  cases       = []


PASO 2 — Consultar backend (con o sin texto)

  GET /api/v1/cases/list
    ?filter_key=persona_asignada          // si agentId o autocomplete activo
    &filter_value={...}
    &chip_filter=casos_nuevos             // si chip activo (E-09)
    &dropdown_filter_key=riesgo|equipo    // si dropdown activo (E-10/E-11)
    &dropdown_filter_value={...}
    &search={searchText}                  // omitido o vacío → sin filtro de búsqueda
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI searchText === "":
    → Restaura la lista del filtro activo sin cláusula search

  SI respuesta no ok:
    → loadError = data.error || 'Error al buscar casos'

  SI respuesta ok:
    → cases, totalCases; loading = false


PASO 3 — Actualizar vista

  filteredCases = cases
  → Tabla o EmptyState "No se encontraron casos"

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (parámetro search)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SI search != "":

    AND (
      vc.victim_case_victim_doc_number ILIKE '%' || search || '%'
      OR vf2.victim_case_form2_victim_phone::text ILIKE '%' || search || '%'
      OR EXISTS (
        SELECT 1
        FROM   salvia.victim_contact_form1 vcf1
        WHERE  vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
        AND    vcf1.victim_contact_form1_phone::text ILIKE '%' || search || '%'
      )
    )

  // vf2 ya está en el JOIN base de E-01
  // victim_case_victim_contact es FK al contacto, no el número de teléfono

SI search == "" o no viene:
  → Sin cláusula adicional

El resto (ordenamiento, paginación, response) es idéntico a E-01.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas (E-03 implementado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Debounce | 400ms |
| Combinación con otros filtros | Sí — parámetro `search` aditivo |
| Teléfono sin contacto vinculado | Se busca en `victim_case_form2_victim_phone` aunque falte form1 |
