# Listado Reportes — Interfaz de la pantalla

Pantalla **nueva** (ruta `/salvia/reportes`, roles `op` y `sv`) para la funcionalidad "Gestión reportes" descrita en FigJam (canal `l25jrzyf`). Muestra **todos** los reportes (`victim_contact`), asignados o no, con filtros, y permite al supervisor (`sv`) reasignarlos a un agente integral. No modifica `/salvia/primer_contacto` (pantalla legada que cumple hoy un rol parecido) — ambas coexisten mientras se decide una migración futura (GAP-4 en [related-tables.md](./related-tables.md)).

Reutiliza `casos-component` y `reasignar-casos-modal` (ya en producción en `list_cases.html`), extendidos de forma aditiva — ver detalle de la extensión en `DocsMD/Componentes/Casos component/`. La pantalla en sí es una capa fina de composición, igual que ya existe para casos (`list_cases.html`).

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/reportes/list_reportes.html` | Pantalla — composición (nuevo) |
| `src/frontend/js/components/casos-component.js` | Tabla + filtros (existente, extendido — ver nota) |
| `src/frontend/html/salvia/case_component/case_component.html` | Template de `casos-component` (extendido: columnas de reporte, `sortOptions`, `itemLabelPlural`) |
| `src/frontend/js/components/reasignar-casos-modal.js` | Modal de reasignación masiva (existente, extendido) |
| `src/frontend/html/salvia/list-cases/reasignar_casos_modal.html` | Template del modal de reasignación (textos condicionales por `entityKind`) |
| `src/frontend/js/components/gestion-reporte-modal.js` | Popup "Ver Reporte" (nuevo — ver `DocsMD/Componentes/Gestion Reporte Modal/`) |
| `src/salvia/facades/VictimContactFacade.go` | Handlers `ListReportsGET` (esta pantalla) y `MyReportsGET` (Mis Reportes): validan el permiso y renderizan el HTML |
| `src/salvia/facades/MainRouter.go` | Registra `GET /salvia/reportes` (y `/salvia/mis-reportes`) |
| `src/salvia/config/Menu.go` | Permisos nuevos: `get_list_reports` (op, sv), `list_victim_contacts` (op, sv), `manage_victim_contacts` (op, sv), `reassign_victim_contacts` (sv) |
| `src/salvia/controller/victim_contacts_list_controller.go` · `service/victim_contacts_list_service.go` · `internal/repository/victim_contacts_list_repository.go` | `GET /api/v1/victim-contacts/list` |
| `src/salvia/controller/victim_contacts_reassign_controller.go` · `service/victim_contacts_reassign_service.go` · `internal/repository/victim_contacts_reassign_repository.go` | `POST /api/v1/victim-contacts/reasignar-bulk` |
| `src/internal/models/report_list_item.go` | `ReportListItem` (fila del listado) y `ReportDetail` (detalle del modal) |
| `src/frontend/html/sidebar/sidebar.html` | Entrada "Listado Reportes" del menú lateral (op, sv) |

Entry point backend: `GET /salvia/reportes` → `ListReportsGET` (en `VictimContactFacade.go`; no hay un facade aparte). Solo renderiza el HTML — la data llega por el endpoint JSON `GET /api/v1/victim-contacts/list`, igual que `list_cases.html` no trae los casos en el HTML inicial.

---

## Árbol de interfaz

```
list_reportes.html  (.salvia-main)
│
├── WelcomeSection
│   └── h2  "Reportes"
│
├── casos-component                                           ← EXTENDIDO (ver nota)
│   :list-endpoint="/api/v1/victim-contacts/list"              // prop nueva, default sigue siendo /api/v1/cases/list
│   :default-filter="{}"                                      // sin scope — ve todo (el componente lo normaliza a "sin filtro")
│   :columns="[
│       { key: 'victim_info',       label: 'Víctima / ID Reporte' },
│       { key: 'report_type',       label: 'Tipo' },            // NUEVO col.key
│       { key: 'contact_scheme',    label: 'Esquema 3x3' },      // NUEVO col.key
│       { key: 'report_status',     label: 'Estado' },           // NUEVO col.key
│       { key: 'assigned_agent',    label: 'Asignado a' },        // NUEVO col.key
│       { key: 'registration_date', label: 'Fecha Reporte' },    // ya existe en casos-component
│   ]"
│   :available-filters="[
│       { key: 'estado_reporte', type: 'dropdown',    label: 'Estado del reporte', options: [...] },  // NUEVO
│       { key: 'asignacion',     type: 'autocomplete', label: 'Asignación' },                          // NUEVO (incluye opción fija "Sin asignar")
│       { key: 'busqueda',       type: 'search',      label: 'Buscar por nombres, apellidos, documento o teléfono' }, // mecanismo heredado, traducción SQL nueva — ver E-07
│   ]"
│   :sort-options="[
│       { value: 'registration_date_desc', label: 'Fecha de Reporte DESC' },  // NUEVO prop — ver flow-E08
│       { value: 'registration_date_asc',  label: 'Fecha de Reporte ASC' },
│   ]"
│   :buttons="[{ id: 'ver', label: 'Ver Reporte' }]"
│   item-label-plural="reportes"                               // prop nuevo — "reporte(s) seleccionado(s)", "Reasignar Reportes", "No se encontraron reportes"
│   :reasignacion="canReasignarReportes"                       // currentRole === 'sv'; para op no existen ni las casillas ni el botón
│   @action-clicked="onReporteActionClicked"   // → E-04
│   @reasignar-casos="onReasignarReportes"     // → E-05
│
├── reasignar-casos-modal                                      ← EXTENDIDO (ver nota)
│   v-if="canReasignarReportes"            // solo se monta para sv
│   ref="reasignarModal"
│   entity-kind="reportes"                 // prop nueva — cambia candidatos, endpoint y textos
│   @closed="onReasignarModalClosed"       // no recarga: la selección de filas se conserva
│   @reassigned="onReasignarCompletado"    // → E-06 (recarga casos-component con reload())
│
└── gestion-reporte-modal                                       ← NUEVO componente
    ref="gestionModal"
    @updated="onReporteActualizado"        // recarga casos-component tras cerrar/intento
```

---

## Nota de extensión — `casos-component`

Cambios aditivos, sin tocar el comportamiento actual de ningún consumidor existente (`get_my_cases.html`, `list_cases.html`, `case_component.html`). Detalle completo en [`casos-component-interface.md` → «Extensión: Gestión de reportes»](../../Componentes/Casos%20component/casos-component-interface.md).

| Cambio | Tipo | Detalle |
|---|---|---|
| Prop `listEndpoint` | Nuevo, opcional | Default `/api/v1/cases/list` (comportamiento actual intacto). Esta pantalla pasa `/api/v1/victim-contacts/list`, que responde `{ reports, total, page, pageSize }`: el componente acepta `cases` o `reports`. |
| Prop `sortOptions` | Nuevo, opcional | El `<select>` de orden itera este prop (2 opciones) en vez de las 4 fijas. Ver [flow-E08](./Flujos/flow-E08-cuando-ordena-por-fecha.md). |
| Prop `itemLabelPlural` | Nuevo, opcional | Default `'casos'`. Esta pantalla pasa `'reportes'`: texto del contador ("reporte(s) seleccionado(s)", con el singular derivado en `itemLabelSingular`), del botón ("Reasignar Reportes") y del vacío ("No se encontraron reportes"). |
| `col.key === 'victim_info'` | Variante para reportes | Si la fila trae `reportStatus`: nombre (o "Sin nombre") + `ID: {i_code}`; sin documento ni teléfono. |
| `col.key === 'report_type'` | Nueva rama en el `<template>` de celda | Muestra `report.reportType` — la etiqueta (Víctima/Tercero/Familiar/Hombre) la arma el servicio; `—` si el reporte no tiene tipo (los registrados por un agente, solo-`form1`, no lo tienen). |
| `col.key === 'contact_scheme'` | Nueva rama | `"Sin intentos"` o `"{totalFallidos}/9 intentos · {diasCompletos}/3 días"` — es el texto real de la columna de la bandeja en el mockup (`totalFallidosDe`/`diasCompletosDe`), no el contador "hoy" que solo aplica dentro del modal. |
| `col.key === 'report_status'` | Nueva rama | Badge: Asignado (azul) / Cerrado — válido (verde) / Cerrado — inválido (rojo) / Sin asignar (gris). |
| `col.key === 'assigned_agent'` | Nueva rama | `report.assignedAgentName \|\| "Sin asignar"` |
| Filtro `estado_reporte` (dropdown) | Entra por `activeDropdowns` | Misma lista `dropdownParamKeys` que `riesgo`/`equipo`/`estado_caso`: `filter_estado_reporte={asignado\|cerrado_valido\|cerrado_invalido\|sin_asignar}`. Combinable; no toca `activeFilter`. Ver [flow-E02](./Flujos/flow-E02-cuando-filtra-estado-reporte.md). |
| Filtro `asignacion` (autocomplete) | `runAsignacionSearch` (rama propia) | Agentes de `GET /api/v1/agents/by-role?role=op` (una vez, en caché, filtrados en el cliente) + opción fija `{value:'__unassigned__', label:'Sin asignar'}`: `filter_key=asignacion&filter_value={icode\|__unassigned__}`. Ver [flow-E03](./Flujos/flow-E03-cuando-filtra-asignacion.md). |

### Traducción del filtro `busqueda` contra `victim_contact` (corrección de planeación, 2026-10-11)

La planeación original decía "se hereda sin cambios de casos-component" — impreciso: el *mecanismo* (debounce 400ms, reset de página) sí se hereda, pero la *traducción a SQL* es nueva porque `victim_contact` no tiene la forma de `victim_case`. Existe código real ya probado en producción para esto — `GetVictimContactsWithoutVictimCase` (usada hoy por `/salvia/primer_contacto`) — que el backend nuevo debe reutilizar en vez de reinventar:

| Campo buscado | Fuente real (confirmada en código) |
|---|---|
| Nombres / Apellidos | `victim_contact.victim_contact_names` / `victim_contact_last_names` — columnas directas |
| Teléfono | `EXISTS` contra `victim_contact_form2.victim_col_phone` **O** `victim_contact_form1.phone` — un reporte de tercero (Form2) y uno propio (Form1) guardan el teléfono en tablas distintas; el código real ya corrige un bug donde solo se buscaba en Form2 |
| Documento | Solo `victim_contact_form1.doc_number` — Form2 (reporte de tercero) no tiene columna de documento, así que un reporte sin Form1 no aparece al buscar por documento |

Fuente: `VictimContactDAO.go:22-27` (`VictimContactListFilters`) y `VictimContactDAO.go:316-372` (`buildVictimContactsWithoutCaseWhere`). El endpoint nuevo `GET /api/v1/victim-contacts/list` debe adaptar esta misma lógica (un único parámetro `search` en vez de 4 campos separados, para combinar con el resto de filtros de `casos-component` vía OR en vez de AND) — ver [flow-E07](./Flujos/flow-E07-cuando-busca-por-texto.md).

---

## Nota de extensión — `reasignar-casos-modal`

| Cambio | Detalle |
|---|---|
| Prop `entityKind` | Nuevo, opcional. Default `'casos'` (comportamiento actual intacto). `'reportes'` pide los candidatos con `role=op` (`GET /api/v1/equipo-operadores?role=op`), guarda en `POST /api/v1/victim-contacts/reasignar-bulk` (body `{ report_icodes, agent_icode }`) y cambia los textos ("Reasignar reportes", "Reportes seleccionados (n)", advertencia y confirmación), sin `Doc:` ni badge de riesgo. |
| Resolución de equipo (M-01) | **Sin rama propia para reportes.** Que no se exija equipo (GAP-3: cualquier agente `op` activo es candidato) lo da hoy el modo contingencia (`REASSIGN_CROSS_TEAM_CONTINGENCY = true`). Con el flag en `false` el modal intentaría resolver el equipo desde `caseTeam`/`riskStatus` y fallaría para un reporte (ver la limitación en la interfaz del modal). |

Detalle completo en [`reasignar-casos-modal-interface.md` → «Extensión: Gestión de reportes»](../../Componentes/Casos%20component/reasignar-casos-modal-interface.md).

---

## Permisos (nuevo, `Menu.go`)

| Rol | Acceso a `/salvia/reportes` | Puede reasignar | Puede gestionar (intento/cierre) |
|---|---|---|---|
| `sv` | Sí | Sí | Sí — cualquier reporte (GAP-2, confirmado 2026-10-11) |
| `op` | Sí | No | Sí — cualquier reporte, no solo los propios (GAP-2) |
| `ro`, `an`, otros | No | No | No |

Claves de permiso en `Menu.go`: `get_list_reports` (pantalla), `list_victim_contacts` (listado JSON), `manage_victim_contacts` (detalle, intento y cierre) y `reassign_victim_contacts` (reasignación, solo `sv`). Un rol sin permiso recibe **401** —`CheckPermission` responde 401, no 403— tanto en la pantalla como en los endpoints.

## Entrada en el menú lateral

`sidebar.html` muestra el ítem **Listado Reportes** (`/salvia/reportes`, ícono `fa-file-alt`) a los roles `op` y `sv`, justo después de "Lista Casos". Ojo: el ítem "Reportes" que ya existía apunta a `/salvia/casos` (es el listado de casos), no a esta pantalla. Además, la pantalla legada `/salvia/primer_contacto` trae el enlace "Ver nuevo Listado de Reportes" (ver el índice).

---

## Relación con `Mis Reportes`

`Mis Reportes` es la **misma** composición con `:agent-id="currentUserId"` (el componente fuerza `filter_key=persona_asignada&filter_value={agentId}`, oculta el filtro de equipo y bloquea el autocomplete de asignación) y `reasignacion=false` — no es una pantalla con lógica propia distinta. Ver [`Mis Reportes`](../Mis%20Reportes/mis-reportes-interface.md).
