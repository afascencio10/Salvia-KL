━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando ordena por fecha de reporte
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambia la opción del `<select>` "Ordenar por" dentro de
               casos-component

Precondición: ninguna — reemplaza por completo al SortSelector actual

INPUT: {
  value:  opción elegida   → 'registration_date_desc' | 'registration_date_asc'
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js (requiere el prop nuevo sortOptions)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Renderizar el <select> desde el prop, no desde opciones fijas

  <select :value="sortOption">
    <option v-for="opt in (sortOptions || DEFAULT_SORT_OPTIONS)" :value="opt.value">
      {{ opt.label }}
    </option>
  </select>

  // DEFAULT_SORT_OPTIONS = las 4 actuales (registration_date_*, next_follow_up_*)
  // — así get_my_cases.html/list_cases.html/case_component.html, que no pasan
  // este prop, siguen viendo exactamente lo mismo que hoy.
  //
  // list_reportes.html pasa sortOptions con solo 2 entradas (ver
  // listado-reportes-interface.md) — "Próximo seguimiento" nunca aparece
  // para un reporte, porque no existe follow_up_v2 para victim_contact.


PASO 2 — onSortChange(event)  [mismo mecanismo ya documentado como E-04]

  sortOption = value
  currentPage = 1
  GET {listEndpoint}?...&sort=registration_date&order={desc|asc}&page=1...

  // El backend de /api/v1/victim-contacts/list solo soporta
  // ORDER BY victim_contact_creation_date — no tiene un equivalente de
  // "next_follow_up" porque esa opción nunca se le pide. Un sort distinto no
  // da error: se ignora y se ordena por fecha de reporte (order distinto de
  // 'asc' se toma como desc).

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Ninguno.
