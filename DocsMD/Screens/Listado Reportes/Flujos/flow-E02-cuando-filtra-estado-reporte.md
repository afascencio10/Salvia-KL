━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por estado del reporte
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en el <select> del filtro "Estado del reporte"
               dentro de casos-component (filter.type === 'dropdown')

INPUT: {
  value:  opción elegida   → 'asignado' | 'cerrado_valido' | 'cerrado_invalido' | 'sin_asignar' | '' (Todos / Todas)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js (extensión)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — setDropdownFilter('estado_reporte', value)

  SI value === activeDropdowns['estado_reporte'] (sin cambio):
    → TERMINAR ejecución
  SI value === "":
    quitar la clave 'estado_reporte' de activeDropdowns
  SI NO:
    activeDropdowns['estado_reporte'] = value

  searchText = ""     // heredado de setDropdownFilter: cambiar un dropdown limpia el buscador
  currentPage = 1
  cases = []

  // activeFilter NO se toca: el filtro de estado es COMBINABLE con el resto (asignación,
  // búsqueda y el scope de Mis Reportes, que va por la prop agentId). Una versión anterior
  // de este flujo reemplazaba activeFilter por { key:'estado_reporte', value }; eso habría
  // quitado el scope de Mis Reportes y mostrado reportes de otros agentes.


PASO 2 — fetchCases(): limpia la selección de filas y llama al backend

  GET {listEndpoint}?sort=registration_date&order={desc|asc}&page=1&page_size=20
                      [&filter_key={...}&filter_value={...}]     // scope o asignación, si hay
                      [&filter_estado_reporte={value}]           // solo si hay valor
                      [&search={searchText}]

  // 'estado_reporte' viaja por la misma lista dropdownParamKeys que riesgo, equipo,
  // seguimientos_ejecutados, estado_caso y barreras_activas (E-02/E-13 de casos-component).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/victim-contacts/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 3 — Traducir filter_estado_reporte a una condición SQL (se une con AND al resto)

  abierto = victim_contact_close_reason IS NULL AND victim_contact_status <> 'i'

  SI filter_estado_reporte === 'asignado':
    WHERE victim_contact_agent_id IS NOT NULL AND abierto
  SI filter_estado_reporte === 'sin_asignar':
    WHERE victim_contact_agent_id IS NULL AND abierto
  SI filter_estado_reporte === 'cerrado_valido':
    WHERE victim_contact_close_reason IN (enum categoría 'victim_contact_close_reason', código 'ap')
  SI filter_estado_reporte === 'cerrado_invalido':
    WHERE victim_contact_close_reason IN (enums de esa categoría con código <> 'ap': 'ni','nc','co','dp')
       OR (victim_contact_close_reason IS NULL AND victim_contact_status = 'i')   // reporte legado invalidado
  SI NO (vacío o ausente):
    sin condición adicional

  // Un reporte legado invalidado (status 'i') cuenta como cerrado: no es "asignado" ni
  // "sin asignar" (related-tables.md, nota de convivencia con el botón "Invalidar" legado).
  // 'cerrado_valido' solo existe por BD en esta fase: M-10 está diferido y /cerrar rechaza 'ap'.

  → Devuelve { reports, total, page, pageSize } (mismo shape de página que CasesListResponse,
    con el arreglo en "reports")

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                  | Resolución |
|-----------------------------------------------------------------------|------------|
| Etiquetas exactas del <select> (texto visible) | "Sin asignar", "Asignado", "Cerrado — válido", "Cerrado — inválido" (en ese orden, con raya larga) + la opción "Todos / Todas" que agrega el componente |
