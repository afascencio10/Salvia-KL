━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando busca por texto
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: escribe en el input del filtro `busqueda` (filter.type === 'search')
               dentro de casos-component

INPUT: {
  text:  texto escrito   → searchText
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js (sin cambios — ya documentado como E-03)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Debounce y llamada al backend

  onSearchInput()  [debounce 400ms]
  currentPage = 1
  GET {listEndpoint}?filter_key={activeFilter.key}&filter_value={activeFilter.value}
                      &search={searchText}&page=1&page_size=20&sort={sortOption}

  // Mecanismo idéntico al ya documentado para casos (E-03 de
  // casos-component-events.md) — no se repite aquí. Lo único nuevo es
  // cómo el backend traduce `search` cuando listEndpoint apunta a
  // /api/v1/victim-contacts/list.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/victim-contacts/list (nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — Traducir `search` a condiciones SQL sobre victim_contact

  // Adaptado de buildVictimContactsWithoutCaseWhere (VictimContactDAO.go:316-372),
  // la función que ya usa /salvia/primer_contacto hoy — mismo origen de
  // datos, pero aquí los 4 campos se combinan con OR en un solo parámetro
  // `search` (allá son 4 inputs AND-combinados: Names, LastNames, Phone,
  // Document).

  SI search != "":
    WHERE (
      victim_contact.victim_contact_names      ILIKE '%search%'
      OR victim_contact.victim_contact_last_names ILIKE '%search%'
      OR (victim_contact.victim_contact_names || ' ' || victim_contact.victim_contact_last_names) ILIKE '%search%'
                                                           // nombre completo: "Ana Pérez" coincide aunque estén en columnas distintas
      OR EXISTS (
           SELECT 1 FROM victim_contact_form2 vcf2
           WHERE vcf2.victim_contact_form2_victim_contact = victim_contact.victim_contact_id
             AND vcf2.victim_contact_form2_victim_col_phone::text ILIKE '%search%'
         )
      OR EXISTS (
           SELECT 1 FROM victim_contact_form1 vcf1
           WHERE vcf1.victim_contact_form1_victim_contact = victim_contact.victim_contact_id
             AND vcf1.victim_contact_form1_phone::text ILIKE '%search%'
         )
      OR EXISTS (
           SELECT 1 FROM victim_contact_form1 vcf1d
           WHERE vcf1d.victim_contact_form1_victim_contact = victim_contact.victim_contact_id
             AND vcf1d.victim_contact_form1_doc_number::text ILIKE '%search%'
         )
    )

  // Nota: un reporte de tercero (Form2) no tiene columna de documento — no
  // aparecerá al buscar por documento, igual que hoy en /salvia/primer_contacto.
  // Un reporte V2 (victim_contact_form2_v2) reutiliza el mismo
  // victim_contact_form2 subyacente, así que esta consulta ya lo cubre sin
  // cambios adicionales.


PASO 3 — Responder

  Igual al resto de /api/v1/victim-contacts/list — { reports, total, page, pageSize }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Resolución |
|------------------------------------------------------------------------------------|------------|
| Confirmar que `victim_contact_form2_v2` no requiere una condición EXISTS propia — related-tables.md GAP-1 | Confirmado: un reporte V2 también tiene fila en `victim_contact_form2` (el teléfono vive ahí), así que el EXISTS existente lo cubre; verificado con reportes creados por `/salvia/public-v2` (caso `LR-E07-C2`) |
