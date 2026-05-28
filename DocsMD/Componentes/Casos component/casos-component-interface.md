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
    │   │
    │   ├── FilterChips  (.cc-filter-chips)
    │   │   │
    │   │   ├── [v-if 'casos_nuevos' IN availableFilters]
    │   │   │   CasosNuevosChip  (.cc-chip)  "Casos nuevos"
    │   │   │   :active si activeFilter.key === 'casos_nuevos'
    │   │   │   → selectFilter('casos_nuevos')
    │   │   │
    │   │   ├── [v-if 'riesgo' IN availableFilters]
    │   │   │   RiesgoDropdown  (.cc-dropdown)
    │   │   │   ├── Label  "Nivel de riesgo"
    │   │   │   └── <select>  (.cc-select)
    │   │   │       ├── <option value="">  "Todos"
    │   │   │       ├── <option value="alto">  "Alto"
    │   │   │       ├── <option value="medio">  "Medio"
    │   │   │       └── <option value="bajo">  "Bajo"
    │   │   │       → selectFilter('riesgo', value)
    │   │   │
    │   │   ├── [v-if 'equipo' IN availableFilters]
    │   │   │   EquipoDropdown  (.cc-dropdown)
    │   │   │   ├── Label  "Equipo"
    │   │   │   └── <select>  (.cc-select)
    │   │   │       ├── <option value="">  "Todos"
    │   │   │       └── <option> × N  [v-for availableTeams]
    │   │   │       → selectFilter('equipo', value)
    │   │   │
    │   │   └── [v-if 'persona_asignada' IN availableFilters]
    │   │       PersonaDropdown  (.cc-dropdown)
    │   │       ├── Label  "Persona asignada"
    │   │       └── <select>  (.cc-select)
    │   │           ├── <option value="">  "Todas"
    │   │           └── <option> × N  [v-for availableAgents]  value=agent.icode
    │   │           → selectFilter('persona_asignada', value)
    │   │
    │   └── SearchAndSort  (.cc-search-sort)
    │       │
    │       ├── SearchInput  (.cc-search)
    │       │   └── <input type="text">  placeholder="Buscar por ID o teléfono"
    │       │       → onSearch(event)
    │       │
    │       └── SortSelector  (.cc-sort)
    │           └── <select>  (.cc-select)
    │               ├── <option value="registration_date">  "Fecha de registro"
    │               └── <option value="next_follow_up">  "Próximo seguimiento"
    │               → onSortChange(event)
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
| `availableFilters` | `Array<string>` | No | `[]` | Keys de filtros a mostrar: `'casos_nuevos'`, `'riesgo'`, `'equipo'`, `'persona_asignada'` |
| `columns` | `Array<Object>` | Sí | — | Columnas a mostrar. Cada objeto: `{ key: string, label: string }` |
| `hiddenColumns` | `Array<string>` | No | `[]` | Keys de columnas a ocultar aunque estén en `columns` |
| `buttons` | `Array<Object>` | No | `[]` | Botones de acción por fila. Cada objeto: `{ id: string, label: string }` |

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
