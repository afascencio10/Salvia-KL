# `reasignar-casos-modal` — Interfaz del Modal

Modal reutilizable para reasignar uno o varios casos seleccionados desde `casos-component`. Lo consume el **componente padre** de la pantalla (p. ej. Lista de casos). Se abre cuando el padre recibe el evento `reasignar-casos` (E-15) y llama al método público `open(cases)`.

Permite revisar los casos a reasignar, elegir un nuevo agente y confirmar.

> **Modo normal:** el agente destino pertenece al mismo equipo del caso (derivado de `caseTeam` o riesgo).
>
> **Modo contingencia (cross-team):** el supervisor puede elegir **cualquier agente `ro` activo**, sin filtro por equipo. Ver [contingencia-reasignacion-cross-team.md](./contingencia-reasignacion-cross-team.md).

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/reasignar-casos-modal.js` | Componente Vue — template y lógica |
| `src/frontend/css/reasignar-casos-modal.css` | Estilos del modal y estados visuales |
| [contingencia-reasignacion-cross-team.md](./contingencia-reasignacion-cross-team.md) | Plan de contingencia — reasignación a cualquier agente |
| `src/frontend/html/salvia/list-cases/reasignar_casos_modal.html` | Partial HTML con `x-template` del modal |

## Relación con otros componentes

```
list_cases.html (padre)
│
├── casos-component
│   └── emite reasignar-casos { cases }          → E-15
│
└── reasignar-casos-modal
    ├── ref="reasignarModal"
    ├── open(cases)  ← padre llama tras E-15
    └── (eventos M-01 … M-04 documentados abajo)
```

---

## Árbol de interfaz

```
reasignar-casos-modal
│
├── [v-if !visible]  → no renderiza nada
│
└── [v-if visible]
    ModalBackdrop  (.rcm-backdrop)
    │  @click.self → close()   // → M-04 (solo si no loading)
    │
    └── ModalDialog  (.rcm-dialog)  role="dialog" aria-modal="true"
        │
        ├── ModalHeader  (.rcm-header)
        │   ├── Title  (.rcm-title)  "Reasignar casos"
        │   └── CloseBtn  (.rcm-close)  "✕"
        │       → close()   // → M-04
        │
        ├── ModalBody  (.rcm-body)
        │   │
        │   ├── TeamBanner  (.rcm-team-banner)  [v-if !crossTeamContingency && resolvedTeam]
        │   │   └── "Equipo: { resolvedTeam }"
        │   │   // Modo normal: equipo calculado al abrir (M-01)
        │   │   // Contingencia: oculto o mensaje alternativo (ver § Modo contingencia)
        │   │
        │   ├── CasesSection  (.rcm-cases-section)
        │   │   ├── SectionTitle  (.rcm-section-title)  "Casos seleccionados ({ cases.length })"
        │   │   └── CasesList  (.rcm-cases-list)
        │   │       └── CaseRow × N  [v-for cases]  (.rcm-case-row)
        │   │           ├── VictimName  (.rcm-case-victim)
        │   │           │   case.names + " " + case.lastNames
        │   │           ├── RiskBadge  (.rcm-case-risk)
        │   │           │   riskLabel(case.riskStatus)   // Bajo | Moderado | Alto | Extremo | —
        │   │           └── AssignedPerson  (.rcm-case-owner)
        │   │               ownerFullName(case) || "Sin asignar"
        │   │               // ownerNames + ownerLastNames (rel_case_owner activo)
        │   │
        │   ├── AgentSection  (.rcm-agent-section)
        │   │   ├── Label  (.rcm-label)  "Nueva persona asignada"
        │   │   │
        │   │   ├── [v-if loadingAgents]
        │   │   │   LoadingAgents  (.rcm-agents-loading)
        │   │   │   └── Spinner + "Cargando agentes..."
        │   │   │
        │   │   ├── [v-else-if agentsError]
        │   │   │   AgentsError  (.rcm-agents-error)
        │   │   │   └── agentsError
        │   │   │
        │   │   └── [v-else]
        │   │       <select>  (.rcm-select)  v-model="selectedAgentIcode"
        │   │       ├── <option value="">  "Seleccione un agente"
        │   │       └── <option> × N  [v-for agents]
        │   │           :value="agent.icode"
        │   │           Modo normal:     agent.fullName
        │   │           Contingencia:   agent.fullName + ' — ' + agent.team
        │   │           → onAgentChange()   // → M-03
        │   │
        │   └── WarningBox  (.rcm-warning)
        │       ├── Icon  (.rcm-warning-icon)  ⚠
        │       └── Text  (.rcm-warning-text)
        │           "Los cambios son irreversibles. El caso y sus seguimientos
        │            pendientes quedarán asignados a la nueva persona seleccionada."
        │
        └── ModalFooter  (.rcm-footer)
            ├── CancelBtn  (.rcm-btn.rcm-btn--secondary)  "Cancelar"
            │   → close()   // → M-04
            └── ReassignBtn  (.rcm-btn.rcm-btn--primary)  "Reasignar"
                :disabled si !selectedAgentIcode || loadingAgents || saving
                → confirmReassign()   // evento de guardado — pendiente de planeación
```

---

## Método público de apertura

| Método | Parámetros | Descripción |
|---|---|---|
| `open(cases)` | `Array<CaseListItem>` | Abre el modal con los casos recibidos desde E-15. Dispara M-01. |

El padre lo invoca así:

```javascript
onReasignarCasos: function(payload) {
    this.$refs.reasignarModal.open(payload.cases);
}
```

---

## Estado interno del componente

| Variable | Tipo | Inicial | Descripción |
|---|---|---|---|
| `visible` | `Boolean` | `false` | Controla si el modal está abierto |
| `cases` | `Array<Object>` | `[]` | Casos a reasignar (copia del payload E-15) |
| `resolvedTeam` | `String` | `''` | Equipo del caso (modo normal: filtra agentes; contingencia: solo informativo u oculto) |
| `agents` | `Array<AgentOption>` | `[]` | Agentes disponibles en el select |
| `crossTeamContingency` | `Boolean` | `REASSIGN_CROSS_TEAM_CONTINGENCY` | Flag de contingencia — ver doc dedicado |
| `selectedAgentIcode` | `String` | `''` | `icode` del agente elegido en el select |
| `selectedAgent` | `Object \| null` | `null` | Objeto completo del agente seleccionado |
| `loadingAgents` | `Boolean` | `false` | Carga de agentes en curso (M-02) |
| `agentsError` | `String \| null` | `null` | Error al cargar agentes |
| `saving` | `Boolean` | `false` | Reservado para el evento de guardado (pendiente) |

### Shape de `AgentOption` (respuesta backend)

| Campo | Tipo | Fuente |
|---|---|---|
| `icode` | `string` | `general_user.general_user_i_code` |
| `fullName` | `string` | `general_user_profile_names + ' ' + general_user_profile_last_names` |
| `team` | `string` | `general_user.general_user_team` |

---

## Modo contingencia (cross-team)

Activado con `REASSIGN_CROSS_TEAM_CONTINGENCY = true` en `reasignar-casos-modal.js`.

| Aspecto | Modo normal | Contingencia |
|---|---|---|
| Resolución de equipo (M-01 PASO 3) | Obligatoria; error si no se deriva | **Comentada** — no bloquea apertura |
| Carga de agentes | M-02 `fetchAgentsByTeam(resolvedTeam)` | M-02-C `fetchAllAgents()` sin `?team=` |
| TeamBanner | Visible con `resolvedTeam` | Oculto o mensaje de contingencia |
| Opciones del select | Solo `fullName` | `fullName — team` |
| Backend | `?team={resolvedTeam}` | Sin param `team` |
| M-05 (guardado) | Igual | Igual — ya soporta agente de otro equipo |

Plan completo: [contingencia-reasignacion-cross-team.md](./contingencia-reasignacion-cross-team.md)

---

## Regla de resolución de equipo (`resolvedTeam`) — modo normal

> En contingencia esta sección queda **comentada** en código; se conserva aquí para reversión.

Se calcula al abrir el modal (M-01) a partir del primer caso de `cases`. Todos los casos seleccionados comparten el mismo equipo (garantizado por E-14 en `casos-component`).

| Prioridad | Condición | `resolvedTeam` |
|---|---|---|
| 1 | `case.caseTeam` no está vacío | Valor de `caseTeam` (`victim_case.victim_case_team`) |
| 2 | `case.caseTeam` vacío y `riskStatus` es `bajo`, `moderado` o `medio` | `"Riesgo bajo"` |
| 3 | `case.caseTeam` vacío y `riskStatus` es `alto` o `extremo` | `"Riesgo alto"` |
| 4 | No se puede determinar equipo | Mostrar error en modal; no cargar agentes |

### Tabla de derivación por riesgo (fallback)

| `riskStatus` | Equipo derivado |
|---|---|
| `bajo` | Riesgo bajo |
| `moderado` | Riesgo bajo |
| `medio` | Riesgo bajo |
| `alto` | Riesgo alto |
| `extremo` | Riesgo alto |
| `null` / vacío / desconocido | No derivable — error |

---

## Datos mostrados por caso seleccionado

| Campo UI | Fuente en `CaseListItem` | Formato |
|---|---|---|
| Nombre víctima | `names` + `lastNames` | Texto libre |
| Nivel de riesgo | `riskStatus` | Etiqueta: Bajo, Moderado, Alto, Extremo, — |
| Persona asignada | `ownerNames` + `ownerLastNames` | Nombre completo o "Sin asignar" |

---

## Fuente de datos — Agentes

### Contingencia (M-02-C)

```
GET /api/v1/equipo-operadores?role=ro
```

Sin `team` → todos los agentes activos con rol `ro`.

### Modo normal — Agentes por equipo (M-02)

### Tablas involucradas

| Tabla | Esquema | Uso |
|---|---|---|
| `general_user` | `security` | Filtro por `general_user_team` y `general_user_status = 'e'` |
| `general_user_profile` | `security` | Nombre del agente (`general_user_profile_names`, `general_user_profile_last_names`) |
| `rel_role_general_user` | `security` | Relación usuario ↔ rol |
| `role` | `security` | Filtro `role_code = 'ro'` (agentes operativos) |

### Referencia en código existente

| Ubicación | Descripción |
|---|---|
| `src/security/dao/GeneralUserDAO.go` | DAO de `general_user`; campo `GeneralUserTeam` → `general_user_team`; perfil vía `GeneralUserProfile` |
| `src/internal/repository/agent_light_repository.go` | `FindAllByRoleAndTeam(roleCode, team)` — JOIN con `general_user_profile` |
| `src/salvia/controller/case_detail_controller.go` | `GET /api/v1/equipo-operadores?team={team}` — endpoint existente reutilizable |

### Endpoint propuesto para M-02

```
GET /api/v1/equipo-operadores?team={resolvedTeam}
```

Respuesta esperada:

```json
[
  { "icode": "abc-123", "fullName": "María López", "team": "Riesgo bajo" },
  { "icode": "def-456", "fullName": "Juan Pérez",  "team": "Riesgo bajo" }
]
```

> Si en implementación se prefiere un endpoint dedicado (`/api/v1/agents/by-team`), la consulta debe replicar el mismo criterio: `general_user` + `general_user_profile` + rol `ro` + `general_user_team = ?`.

---

## Eventos del modal (resumen)

| ID | Evento | Documento |
|---|---|---|
| M-01 | Cuando abre el modal | [flow-M01](./Flujos/flow-M01-cuando-abre-modal-reasignacion.md) |
| M-02 | Cuando carga agentes por equipo *(modo normal)* | [flow-M02](./Flujos/flow-M02-cuando-carga-agentes-por-equipo.md) |
| M-02-C | Cuando carga todos los agentes *(contingencia)* | [flow-M02-C](./Flujos/flow-M02-C-cuando-carga-todos-los-agentes.md) |
| M-03 | Cuando selecciona un agente | [flow-M03](./Flujos/flow-M03-cuando-selecciona-agente.md) |
| M-04 | Cuando cancela o cierra el modal | [flow-M04](./Flujos/flow-M04-cuando-cancela-modal.md) |
| M-05 | Cuando confirma reasignación (guardado) | [flow-M05](./Flujos/flow-M05-cuando-confirma-reasignacion.md) |

Inventario completo: [reasignar-casos-modal-events.md](./reasignar-casos-modal-events.md)

---

## Estados visuales del botón "Reasignar"

| Condición | Estado |
|---|---|
| Sin agente seleccionado | `:disabled` |
| Cargando agentes (`loadingAgents`) | `:disabled` |
| Error al cargar agentes | `:disabled` |
| Agente seleccionado y agentes cargados | Habilitado → dispara M-05 |

---

## Persistencia al confirmar (M-05)

### Endpoint

```
POST /api/v1/cases/reasignar-bulk
```

| Campo request | Tipo | Descripción |
|---|---|---|
| `case_icodes` | `string[]` | `i_code` de cada caso seleccionado |
| `agent_icode` | `string` | Agente elegido en el modal |

El supervisor (ActorID / ActorName del timeline) se obtiene en el **servidor** desde la sesión, no desde el body.

### Reglas de actualización por caso

| Entidad | Campo | Regla |
|---|---|---|
| `victim_case` | `agent_id` | Siempre → `agent_icode` seleccionado |
| `victim_case` | `victim_case_team` | Solo si vacío → `general_user_team` del agente |
| `follow_up_v2` | `agent_id` | Seguimientos con `status <> 'REALIZADO'` del caso |
| `case_timeline_event` | — | Un registro por caso (ver tabla abajo) |

### Evento de timeline por caso

| Campo | Valor |
|---|---|
| `Category` | `Seguimientos` (`TimelineCategorySeguimientos`) |
| `Type` | `Reasignación de Caso` (`TimelineTypeReasignacionCaso`) |
| `ActorID` | `icode` del supervisor logueado |
| `ActorName` | Nombre completo del supervisor logueado |
| `EventUserID` | Igual que `ActorID` |
| `Description` | Texto descriptivo con agente destino y supervisor |

### Evento emitido al padre tras éxito

| Evento | Payload |
|---|---|
| `reassigned` | `{ cases, agent, updated }` |

El padre debe recargar `casos-component` (método público `reload()` pendiente de implementar).
