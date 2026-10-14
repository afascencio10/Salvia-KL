# Mis Reportes — Interfaz de la pantalla

Pantalla **nueva** (ruta `/salvia/mis-reportes`, rol `op` únicamente; el diseño original decía `/salvia/reportes/mis-reportes`). Es **la misma composición** que [`Listado Reportes`](../Listado%20Reportes/listado-reportes-interface.md), pre-filtrada al agente que tiene la sesión abierta y sin el modo de reasignación — igual relación que ya existe hoy entre `list_cases.html` y `get_my_cases.html`. No tiene lógica propia distinta; por eso este documento es intencionalmente corto.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/reportes/get_my_reportes.html` | Pantalla — composición (nuevo) |
| `src/frontend/js/components/casos-component.js` | Tabla + filtros (existente, extendido — ver `Listado Reportes`) |
| `src/frontend/js/components/gestion-reporte-modal.js` | Popup "Ver Reporte" (nuevo, compartido con `Listado Reportes`) |
| `src/salvia/facades/VictimContactFacade.go` | Handler `MyReportsGET`: valida el permiso `get_my_reports` y renderiza el HTML |
| `src/salvia/facades/MainRouter.go` | Registra `GET /salvia/mis-reportes` |
| `src/salvia/config/Menu.go` | Permiso `get_my_reports` (solo `op`) |
| `src/frontend/html/sidebar/sidebar.html` | Entrada "Mis Reportes" del menú lateral (solo `op`) |

Entry point backend: `GET /salvia/mis-reportes` → `MyReportsGET`, un handler propio en `VictimContactFacade.go` (no una rama del de `Listado Reportes`). Inyecta `currentUserId` (= `s.UserICode`, de la sesión) en el template; la data llega por el mismo endpoint JSON `GET /api/v1/victim-contacts/list`.

---

## Árbol de interfaz

```
get_my_reportes.html  (.salvia-main)
│
├── WelcomeSection
│   └── h2  "Mis Reportes"
│
├── casos-component                                           ← mismas props que Listado Reportes, salvo:
│   :list-endpoint="/api/v1/victim-contacts/list"
│   :agent-id="currentUserId"                                  // fuerza persona_asignada = usuario de la sesión (no editable); lo inyecta el backend desde s.UserICode
│   :default-filter="{}"
│   :columns="[
│       { key: 'victim_info',       label: 'Víctima / ID Reporte' },
│       { key: 'report_type',       label: 'Tipo' },
│       { key: 'contact_scheme',    label: 'Esquema 3x3' },
│       { key: 'report_status',     label: 'Estado' },
│       // sin 'assigned_agent' — siempre es el usuario actual, redundante
│       { key: 'registration_date', label: 'Fecha Reporte' },
│   ]"
│   :available-filters="[
│       { key: 'estado_reporte', type: 'dropdown', label: 'Estado del reporte', options: [...] },
│       { key: 'busqueda',       type: 'search',   label: 'Buscar por nombres, apellidos, documento o teléfono' },
│       // sin 'asignacion' — no aplica, la pantalla ya está acotada a "yo"
│   ]"
│   :sort-options="[
│       { value: 'registration_date_desc', label: 'Fecha de Reporte DESC' },
│       { value: 'registration_date_asc',  label: 'Fecha de Reporte ASC' },
│   ]"
│   :buttons="[{ id: 'ver', label: 'Ver Reporte' }]"
│   item-label-plural="reportes"
│   :reasignacion="false"
│   @action-clicked="onReporteActionClicked"   // → E-02, idéntico a Listado Reportes E-04
│
└── gestion-reporte-modal
    ref="gestionModal"
    @updated="onReporteActualizado"     // recarga casos-component
```

---

## Diferencias vs. `Listado Reportes`

| Aspecto | Listado Reportes | Mis Reportes |
|---|---|---|
| Rol | `op`, `sv` | `op` únicamente (`sv` y los demás reciben 401) |
| Scope | Sin scope (`:default-filter="{}"`) | `:agent-id="currentUserId"` — el componente fuerza `filter_key=persona_asignada&filter_value={agentId}`, oculta el filtro de equipo y bloquea el autocomplete; no se puede quitar |
| Filtro "Asignación" | Sí | No (redundante — todo es "mío") |
| Columna "Asignado a" | Sí | No (redundante) |
| `reasignacion` | `true` (solo sv) | `false` siempre |
| Botón "Reasignar Reportes" | Visible para sv | No existe |

Todo lo demás (estado del reporte, búsqueda — misma traducción SQL que [Listado Reportes → flow-E07](../Listado%20Reportes/Flujos/flow-E07-cuando-busca-por-texto.md) —, ordenamiento — mismo `sortOptions` que [flow-E08](../Listado%20Reportes/Flujos/flow-E08-cuando-ordena-por-fecha.md) —, el popup "Ver Reporte" con validación/3x3/cierre) es **idéntico** — mismo componente, mismos endpoints.

---

## GAP — resuelto

| # | Pregunta | Resolución |
|---|---|---|
| GAP-5 | ¿El scope `persona_asignada` es fijo y el usuario no lo puede quitar (a diferencia de casos, donde es un filtro más)? | Sí. Se resolvió con la prop `agentId` de `casos-component` (la misma que usa `get_my_cases.html`) y no con `defaultFilter`: el scope sobrevive a cualquier filtro, búsqueda u orden (casos `MR-E01-C5` y `MR-E01-C6`) y sale de la sesión, no de la URL (`MR-E01-C4`) |
