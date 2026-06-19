━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario filtra por nombre o teléfono
   Tipo: User Interaction
   Código: E-02
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en NameFilter o PhoneFilter dentro de FilterGroup
               Cada input dispara su propio handler — debounce 400ms compartido

INPUT: {
  searchName:    texto del input de nombre    → v-model de .rc-name-input
  searchPhone:   texto del input de teléfono  → v-model de .rc-phone-input
}

> **Nota:** Paginación server-side. Los parámetros `search_name` y `search_phone`
> son independientes y se combinan con AND cuando ambos tienen valor.

> **Campos buscados en BD:**
> - **Nombre (search_name):**
>   - `victim_contact.victim_contact_names` — ILIKE parcial
>   - `victim_contact.victim_contact_last_names` — ILIKE parcial
>   - Concatenación `names + ' ' + lastNames` — ILIKE parcial
> - **Teléfono (search_phone):**
>   - `victim_contact_form2.victim_contact_form2_victim_col_phone` — ILIKE parcial


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reportes-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Actualizar estado y resetear paginación (debounce 400ms)

  searchName  = searchName.trim()
  searchPhone = searchPhone.trim()
  currentPage = 1
  reports     = []
  loading     = true


PASO 2 — Consultar backend

  GET /api/v1/reports/list
    ?search_name={searchName}             // omitido si vacío
    &search_phone={searchPhone}           // omitido si vacío
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI ambos vacíos:
    → Lista completa de reportes sin caso (status = 'v')

  SI respuesta no ok:
    → loadError = data.error || 'Error al filtrar reportes'
    → loading = false

  SI respuesta ok:
    → reports, totalReports; loading = false


PASO 3 — Actualizar vista

  filteredReports = reports
  totalPages      = Math.ceil(totalReports / pageSize)
  → Tabla o EmptyState "No se encontraron reportes"

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/reports/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Base query: misma de E-01 (sin caso + status = 'v')

SI search_name != "":
    AND (
      vc.victim_contact_names ILIKE '%' || search_name || '%'
      OR vc.victim_contact_last_names ILIKE '%' || search_name || '%'
      OR (vc.victim_contact_names || ' ' || vc.victim_contact_last_names) ILIKE '%' || search_name || '%'
    )

SI search_phone != "":
    AND vcf2.victim_contact_form2_victim_col_phone::text ILIKE '%' || search_phone || '%'

El resto (ordenamiento, paginación, response) es idéntico a E-01.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  EJEMPLOS DE COMBINACIÓN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| searchName | searchPhone | Resultado |
|---|---|---|
| `"María"` | *(vacío)* | Reportes sin caso cuyo nombre o apellido contenga "María" |
| *(vacío)* | `"300"` | Reportes sin caso cuyo teléfono contenga "300" |
| `"García"` | `"3001234567"` | Reportes que cumplan **ambas** condiciones (AND) |
| *(vacío)* | *(vacío)* | Todos los reportes sin caso válidos |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Debounce | 400ms por input |
| Inputs separados | Sí — nombre y teléfono independientes |
| Combinación | AND cuando ambos tienen valor |
| Alcance base | Siempre sin caso + status = 'v' (E-01) |
