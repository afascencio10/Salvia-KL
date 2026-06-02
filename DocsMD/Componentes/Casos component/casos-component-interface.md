# `casos-component` — Interfaz del Componente

Tabla reutilizable para visualizar y filtrar casos de víctimas. Acepta props para configurar columnas visibles, filtros disponibles, filtro inicial y botones de acción por fila. Emite un evento hacia el padre cuando el usuario presiona un botón de acción.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/components/casos-component.js` | Componente principal — template y lógica *(pendiente de crear)* |
| `src/frontend/css/casos-component.css` | Estilos de layout y estados visuales *(pendiente de crear)* |

---

## Árbol de interfaz

```
casos-component  (.cc-wrapper)
│
├── [v-if loading]
│   LoadingState  (.cc-loading)
│   └── Spinner + "Cargando casos..."
│
├── [v-else-if loadError]
│   ErrorState  (.cc-error)
│   └── AlertText  (.alert.alert-danger)  loadError
│
└── [v-else]
    │
    ├── FilterBar  (.cc-filter-bar)
    │   │  // Extensible: agregar filtros nuevos al prop :availableFilters
    │   │  // sin cambiar el template. Tipos soportados:
    │   │  //   'chip' | 'dropdown' | 'search' | 'autocomplete'
    │   │  // Cada tipo de filtro siempre llama al backend al activarse.
    │   │
    │   ├── FilterGroup  (.cc-filter-group)
    │   │   └── FilterItem × N  [v-for availableFilters]
    │   │       │
    │   │       ├── [filter.type === 'chip']
    │   │       │   // Botón toggle — activa/desactiva. Llama al backend (→ E-02)
    │   │       │   // Ejemplo: "Casos nuevos"
    │   │       │   Chip  (.cc-chip)  filter.label
    │   │       │   :class="cc-chip--active" si activeFilter.key === filter.key
    │   │       │   → toggleChipFilter(filter.key)   // → E-02
    │   │       │       SI ya estaba activo → activeFilter = defaultFilter (reset) + llama backend
    │   │       │       SI no estaba activo → activeFilter = { key: filter.key } + llama backend
    │   │       │
    │   │       ├── [filter.type === 'dropdown']
    │   │       │   // Select con opciones fijas o cargadas al montar (E-01).
    │   │       │   // Cambiar selección llama al backend (→ E-02)
    │   │       │   // Ejemplos: "Nivel de riesgo", "Por equipo"
    │   │       │   DropdownFilter  (.cc-filter-dropdown)
    │   │       │   ├── Label  (.cc-filter-label)  filter.label
    │   │       │   └── <select>  (.cc-select)
    │   │       │       ├── <option value="">  "Todos / Todas"
    │   │       │       └── <option> × N  [v-for filter.options]  value=opt.value
    │   │       │           opt.label
    │   │       │       :class="cc-select--active" si activeFilter.key === filter.key
    │   │       │       → setDropdownFilter(filter.key, value)   // → E-02
    │   │       │           SI value === "" → activeFilter = defaultFilter (reset) + llama backend
    │   │       │           SI value !== "" → activeFilter = { key: filter.key, value } + llama backend
    │   │       │
    │   │       ├── [filter.type === 'autocomplete']
    │   │       │   // Campo de texto con sugerencias en dropdown flotante.
    │   │       │   // Tiene sus propios eventos: escribir busca agentes (→ E-07),
    │   │       │   // seleccionar filtra casos (→ E-08).
    │   │       │   // Ejemplo: "Persona asignada"
    │   │       │   AutocompleteFilter  (.cc-autocomplete-filter)
    │   │       │   ├── Label  (.cc-filter-label)  filter.label
    │   │       │   │
    │   │       │   ├── [v-if !autocompleteSelected[filter.key]]
    │   │       │   │   // Estado: escribiendo — muestra el input
    │   │       │   │   InputWrap  (.cc-autocomplete-input-wrap)
    │   │       │   │   ├── <input type="text">  (.cc-autocomplete-input)
    │   │       │   │   │   placeholder=filter.label
    │   │       │   │   │   v-model=autocompleteText[filter.key]
    │   │       │   │   │   → onAutocompleteInput(filter.key)  [debounce 400ms → E-07]
    │   │       │   │   └── [v-if autocompleteLoading[filter.key]]
    │   │       │   │       Spinner  (.cc-autocomplete-spinner)
    │   │       │   │
    │   │       │   ├── [v-else]
    │   │       │   │   // Estado: seleccionado — muestra el tag con la opción elegida
    │   │       │   │   SelectedTag  (.cc-autocomplete-tag)
    │   │       │   │   ├── TagLabel  (.cc-tag-label)
    │   │       │   │   │   autocompleteSelected[filter.key].label
    │   │       │   │   └── ClearBtn  (.cc-tag-clear)  "✕"
    │   │       │   │       → clearAutocomplete(filter.key)   // → E-08 (limpia y resetea filtro)
    │   │       │   │
    │   │       │   └── [v-if autocompleteSuggestions[filter.key].length > 0]
    │   │       │       // Dropdown flotante de sugerencias (visible mientras escribe)
    │   │       │       SuggestionsDropdown  (.cc-autocomplete-dropdown)
    │   │       │       │   :class="cc-dropdown--loading" si autocompleteLoading[filter.key]
    │   │       │       ├── SuggestionItem × N  [v-for autocompleteSuggestions[filter.key]]
    │   │       │       │   (.cc-autocomplete-item)
    │   │       │       │   opt.label   // "Nombres + Apellidos" del agente
    │   │       │       │   → selectAutocompleteSuggestion(filter.key, opt)   // → E-08
    │   │       │       └── [v-if autocompleteSuggestions[filter.key].length === 0
    │   │       │                && !autocompleteLoading[filter.key]
    │   │       │                && autocompleteText[filter.key] !== ""]
    │   │       │           NoResults  (.cc-autocomplete-empty)  "Sin resultados"
    │   │       │
    │   │       └── [filter.type === 'search']
    │   │           // Input de texto libre — busca casos por ID o teléfono.
    │   │           // Combina con el filtro activo. Llama al backend (→ E-03) con debounce 400ms.
    │   │           // Ejemplo: "Buscar por ID o teléfono"
    │   │           SearchFilter  (.cc-search-filter)
    │   │           ├── Label  (.cc-filter-label)  filter.label
    │   │           └── <input type="text">  (.cc-search-input)
    │   │               placeholder=filter.label
    │   │               v-model=searchText
    │   │               → onSearchInput()  [debounce 400ms → E-03]
    │   │
    │   └── SortSelector  (.cc-sort-selector)
    │       ├── Label  (.cc-filter-label)  "Ordenar por"
    │       └── <select>  (.cc-select)  :value="sortOption"
    │           ├── <option value="registration_date_desc">  "Fecha de registro DESC"
    │           ├── <option value="registration_date_asc">   "Fecha de registro ASC"
    │           ├── <option value="next_follow_up_desc">   "Próximo seguimiento DESC"
    │           └── <option value="next_follow_up_asc">    "Próximo seguimiento ASC"
    │           → onSortChange(event)   // → E-04 (recarga backend)
    │
    └── TableWrap  (.cc-table-wrap)
        │
        ├── [v-if filteredCases.length === 0]
        │   EmptyState  (.cc-empty)
        │   └── "No se encontraron casos"
        │
        └── [v-else]
            │
            ├── <table.cc-table>
            │   │
            │   ├── <thead>
            │   │   └── <tr>
            │   │       ├── <th> × N  [v-for visibleColumns]  (.cc-th)
            │   │       │   └── column.label
            │   │       └── [v-if buttons.length > 0]
            │   │           <th>  (.cc-th)  "Acciones"
            │   │
            │   └── <tbody>
            │       └── <tr> × N  [v-for filteredCases]  (.cc-row)
            │           │
            │           ├── [col.key === 'victim_info']
            │           │   <td>  (.cc-cell-victim)
            │           │   ├── Name  (.cc-victim-name)
            │           │   │   case.names + " " + case.lastNames
            │           │   └── ICode  (.cc-victim-icode)
            │           │       case.i_code
            │           │
            │           ├── [col.key === 'operator']
            │           │   <td>  (.cc-cell)
            │           │   ├── [case.ownerNames existe]
            │           │   │   case.ownerNames + " " + case.ownerLastNames
            │           │   └── [v-else]  "—"
            │           │   // ownerNames viene de rel_case_owner_victim_case (status='a')
            │           │   // → perfil del agente asignado
            │           │
            │           ├── [col.key === 'registration_date']
            │           │   <td>  (.cc-cell)
            │           │   └── case.creationDate  (formateada DD/MM/YYYY)
            │           │
            │           ├── [col.key === 'risk_level']
            │           │   <td>  (.cc-cell)
            │           │   └── RiskBadge  (.cc-risk-badge)
            │           │       State: alto | medio | bajo | sin_dato
            │           │       :class según case.riskStatus
            │           │
            │           ├── [col.key === 'team']
            │           │   <td>  (.cc-cell)
            │           │   └── case.ownerTeam || "—"
            │           │   // ownerTeam viene de general_user_team del agente asignado
            │           │
            │           ├── [col.key === 'assigned_person']
            │           │   <td>  (.cc-cell)
            │           │   └── case.ownerNames + " " + case.ownerLastNames || "—"
            │           │
            │           ├── [col.key === 'next_follow_up_date']
            │           │   <td>  (.cc-cell-follow-up)
            │           │   ├── [case.nextFollowUpDate existe]
            │           │   │   DateText  (.cc-follow-up-date)
            │           │   │   case.nextFollowUpDate  (formateada DD/MM/YYYY)
            │           │   └── [v-else]
            │           │       NoPending  (.cc-follow-up-none)  "—"
            │           │
            │           ├── [col.key === 'completed_follow_ups']
            │           │   <td>  (.cc-cell)
            │           │   └── case.completedFollowUpsCount  (entero; "0" si es 0)
            │           │   // COUNT follow_up_v2 WHERE status = 'REALIZADO'
            │           │
            │           └── [v-if buttons.length > 0]
            │               <td>  (.cc-cell-actions)
            │               └── ActionBtn × N  [v-for buttons]  (.cc-action-btn)
            │                   btn.label
            │                   → emitActionClicked(btn.id, case)
            │
            └── PaginationBar  (.cc-pagination)
                ├── PrevBtn  (.cc-page-btn)  "← Anterior"
                │   :disabled si currentPage === 1
                │   → changePage(currentPage - 1)
                ├── PageInfo  (.cc-page-info)
                │   "Página {currentPage} de {totalPages}"
                └── NextBtn  (.cc-page-btn)  "Siguiente →"
                    :disabled si currentPage === totalPages
                    → changePage(currentPage + 1)
```

---

## Props del componente

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `defaultFilter` | `Object` | Sí | — | Filtro inicial. Shape: `{ key: string, value?: string }` |
| `availableFilters` | `Array<FilterDef>` | No | `[]` | Filtros a mostrar en la FilterBar. Ver tabla de FilterDef abajo. |
| `columns` | `Array<Object>` | Sí | — | Columnas a mostrar. Cada objeto: `{ key: string, label: string }` |
| `hiddenColumns` | `Array<string>` | No | `[]` | Keys de columnas a ocultar aunque estén en `columns` |
| `buttons` | `Array<Object>` | No | `[]` | Botones de acción por fila. Cada objeto: `{ id: string, label: string }` |

### Shape de `FilterDef`

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `key` | `string` | Sí | Identificador único del filtro. El backend lo recibe como `filter_key`. |
| `type` | `'chip' \| 'dropdown' \| 'search' \| 'autocomplete'` | Sí | Tipo de control a renderizar. |
| `label` | `string` | Sí | Texto visible en la UI (label y placeholder). |
| `options` | `Array<{ value, label }>` | Solo si `type === 'dropdown'` | Opciones del select. Puede pasarse vacío si las opciones se cargan del backend al montar (E-01). |

### Comportamiento de cada tipo al interactuar

| `type` | Acción del usuario | Llama al backend | Evento |
|---|---|---|---|
| `chip` | Click toggle activa/desactiva | Sí — recarga casos | E-02 |
| `dropdown` | Cambia opción del select | Sí — recarga casos | E-02 |
| `autocomplete` | Escribe en el input | Sí — busca agentes | E-07 |
| `autocomplete` | Selecciona/limpia sugerencia | Sí — recarga casos | E-08 |
| `search` | Escribe en el input | Sí — busca casos (debounce 400ms) | E-03 |

### Filtros actualmente definidos

| `key` | `type` | `label` | `options` |
|---|---|---|---|
| `casos_nuevos` | `chip` | `"Casos nuevos"` | — |
| `riesgo` | `dropdown` | `"Nivel de riesgo"` | Fijas: `{ value:'alto', label:'Alto' }`, `medio`, `bajo` |
| `equipo` | `dropdown` | `"Por equipo"` | Dinámicas — cargadas en E-01 |
| `seguimientos_ejecutados` | `dropdown` | `"Seguimientos ejecutados"` | Fijas: `{ value:'0', label:'0' }` … `{ value:'10', label:'10' }` (→ E-12) |
| `persona_asignada` | `autocomplete` | `"Persona asignada"` | — (se buscan on-demand al escribir → E-07) |
| `busqueda` | `search` | `"Buscar por ID o teléfono"` | — |

---

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `action-clicked` | Usuario presiona un botón de acción en una fila | `{ buttonId: string, case: Object }` |

---

## Estados visuales del RiskBadge

| `riskStatus` | Clase CSS | Texto visible |
|---|---|---|
| `'alto'` | `.cc-risk-badge.alto` | Alto |
| `'medio'` | `.cc-risk-badge.medio` | Medio |
| `'bajo'` | `.cc-risk-badge.bajo` | Bajo |
| `null` / vacío | `.cc-risk-badge.sin-dato` | — |

---

## Keys de columna disponibles

| Key | Fuente de dato | Descripción |
|---|---|---|
| `victim_info` | `VictimCaseLight.names` + `lastNames` + `i_code` | Nombre completo e ID del caso |
| `operator` | `rel_case_owner_victim_case` (status=`'a'`) → perfil del agente | Profesional actualmente asignado al caso |
| `registration_date` | `VictimCaseLight.creationDate` | Fecha de creación del caso |
| `risk_level` | `FollowUpV2.risk_status` | Nivel de riesgo del último seguimiento activo |
| `team` | `AgentLight.team` (del agente asignado vía `rel_case_owner`) | Equipo del profesional asignado al caso |
| `assigned_person` | `rel_case_owner_victim_case` (status=`'a'`) → perfil del agente | Persona con el caso asignado |
| `next_follow_up_date` | `FollowUpV2.scheduled_date` (status=`PENDIENTE`, más próximo) | Fecha del próximo seguimiento sin ejecutar |
| `completed_follow_ups` | `COUNT(follow_up_v2)` donde `status = 'REALIZADO'` y `case_id = victim_case_i_code` | Número de seguimientos ya realizados del caso |
