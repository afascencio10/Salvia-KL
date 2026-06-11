# `casos-component` — Interfaz del Componente

Tabla reutilizable para visualizar y filtrar casos de víctimas. Acepta props para configurar columnas visibles, filtros disponibles, filtro inicial, botones de acción por fila y modo de reasignación masiva. Emite eventos hacia el padre cuando el usuario presiona un botón de acción en una fila o cuando solicita reasignar los casos seleccionados.

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
        ├── [v-if reasignacion && selectedCases.length > 0]
        │   TableToolbar  (.cc-table-toolbar)
        │   └── ReassignBtn  (.cc-reassign-btn)  "Reasignar Casos"
        │       → emitReasignarCasos()   // → E-15
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
            │   │       ├── [v-if reasignacion]
            │   │       │   <th>  (.cc-th.cc-th-select)
            │   │       │   └── SelectAllCheckbox  (.cc-select-all-checkbox)  "Todos"
            │   │       │       → toggleSelectAllCurrentPage()   // → E-14
            │   │       │   // Columna de checkbox — solo visible si :reasignacion === true
            │   │       ├── <th> × N  [v-for visibleColumns]  (.cc-th)
            │   │       │   └── column.label
            │   │       └── [v-if buttons.length > 0]
            │   │           <th>  (.cc-th)  "Acciones"
            │   │
            │   └── <tbody>
            │       └── <tr> × N  [v-for filteredCases]  (.cc-row)
            │           │
            │           ├── [v-if reasignacion]
            │           │   <td>  (.cc-cell-select)
            │           │   └── <input type="checkbox">  (.cc-case-checkbox)
            │           │       :checked si case está en selectedCases
            │           │       → toggleCaseSelection(case)   // → E-14
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
            │           ├── [col.key === 'case_status']
            │           │   <td>  (.cc-cell)
            │           │   └── caseStatusLabel(case.status)
            │           │   // victim_case_status: ra→Activo, is→Con novedad, cd→Cerrado, …
            │           │
            │           ├── [col.key === 'barriers']
            │           │   <td>  (.cc-cell-barriers)
            │           │   ├── [case.openBarriers.length > 0]
            │           │   │   BarrierText  (.cc-barriers-text)
            │           │   │   formatOpenBarriers(case.openBarriers)
            │           │   │   // Ejemplo: "Abierta → Sector: Salud, Protección"
            │           │   │   // Solo barreras con status OPEN (E-01). Sectores únicos,
            │           │   │   // etiquetas legibles vía barrierSectorLabel().
            │           │   └── [v-else]
            │           │       // Sin barreras OPEN → celda vacía (sin "—")
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
| `reasignacion` | `Boolean` | No | `false` | Activa el modo de reasignación masiva. Si es `true`, muestra una columna de checkbox en la tabla y permite seleccionar uno o varios casos para reasignarlos. |

### Modo reasignación (`reasignacion === true`)

| Elemento | Comportamiento |
|---|---|
| Columna de checkbox | Primera columna de la tabla. Checkbox "Todos" en el encabezado y un checkbox por fila. |
| Alcance de selección | Solo casos de la página actual (`filteredCases`). Al cambiar página o recargar, se limpia la selección. |
| Regla de equipo | Solo se pueden seleccionar casos con el mismo `caseTeam`. Si se intenta mezclar equipos, se muestra alerta flotante sobre la tabla (`.cc-table-alert`) y se rechaza la selección (→ E-14). |
| Estado interno `selectedCases` | Array de objetos caso seleccionados. Inicializado vacío en E-01. |
| Botón "Reasignar Casos" | Visible solo si `selectedCases.length > 0`. Ubicado en la parte superior derecha, encima de la tabla (`.cc-table-toolbar`). |
| Emisión al padre | Al presionar el botón emite `reasignar-casos` con los casos seleccionados (→ E-15). El padre abre `reasignar-casos-modal` (ver [reasignar-casos-modal-interface.md](./reasignar-casos-modal-interface.md)). |
| Cambio de página / recarga | Al cambiar de página (E-06) o recargar casos por filtro u ordenamiento, se limpia `selectedCases` y se oculta el botón. |

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
| `estado_caso` | `dropdown` | `"Estado del caso"` | Fijas: `ra` Activo, `is` Con novedad, `cd` Cerrado, `ex` Vencido, `r` Por aprobar, `fc` Recontacto (→ E-13) |
| `barreras_activas` | `dropdown` | `"Barreras activas"` | Fija: `{ value:'OPEN', label:'Abierta' }` (→ E-16) |
| `persona_asignada` | `autocomplete` | `"Persona asignada"` | — (se buscan on-demand al escribir → E-07) |
| `busqueda` | `search` | `"Buscar por ID o teléfono"` | — |

---

## Formato de la columna `barriers`

| Condición | Texto en celda |
|---|---|
| Sin barreras con `status = 'OPEN'` | *(vacío — no mostrar "—")* |
| Una o más barreras OPEN | `Abierta → Sector: {sectores}` |

Sectores agregados sin repetir, traducidos desde `barrier_v2.sector`:

| Código BD | Etiqueta UI |
|---|---|
| `salud` | Salud |
| `justicia` | Justicia |
| `proteccion` | Protección |
| `otras_instituciones` | Otras instituciones |
| `barrera_transversal` | Barrera Transversal |

Ejemplo: caso con dos barreras OPEN (salud + proteccion) → `Abierta → Sector: Salud, Protección`

---

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `action-clicked` | Usuario presiona un botón de acción en una fila | `{ buttonId: string, case: Object }` |
| `reasignar-casos` | Usuario presiona "Reasignar Casos" con al menos un caso seleccionado (`reasignacion === true`) | `{ cases: Array<Object> }` — objetos completos de los casos en `selectedCases` |

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
| `case_status` | `victim_case.victim_case_status` (código) → etiqueta vía `caseStatusLabel` | Estado del caso (Activo, Cerrado, etc.) |
| `barriers` | `salvia.barrier_v2` (`status = 'OPEN'`, `case_id = victim_case_i_code`) → sectores únicos vía `formatOpenBarriers` | Barreras abiertas del caso y sector de cada una. Celda vacía si no hay barreras OPEN. |
