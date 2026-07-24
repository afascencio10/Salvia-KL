# Contingencia — Reasignación de casos a cualquier agente (cross-team)

> **Motivo:** contingencia operativa. El supervisor debe poder reasignar casos a **cualquier agente activo (rol `ro`)**, sin limitar el listado al equipo derivado del caso (riesgo / `caseTeam`).
>
> **Estrategia de implementación:** no eliminar la lógica actual. Comentarla y habilitar la rama nueva mediante un flag explícito en frontend. Patrón ya aplicado en `get_case_detail_sv.html` (`cargarOperadores` / `cargarAgentesRO`).

---

## Alcance

| Capa | ¿Cambia? | Detalle |
|---|---|---|
| **Modal `reasignar-casos-modal`** | **Sí** | Dejar de filtrar agentes por `resolvedTeam`; cargar todos los agentes `ro` activos |
| **Backend `GET /api/v1/equipo-operadores`** | **No** | Ya soporta omitir `?team=` → devuelve todos los agentes del rol |
| **Backend `POST /api/v1/cases/reasignar-bulk` (M-05)** | **No** | Ya permite agente de otro equipo; actualiza `follow_up_v2.team` con el equipo del agente destino |
| **`casos-component` E-14 (selección mismo equipo)** | **Sí** | Flag `REASSIGN_CROSS_TEAM_CONTINGENCY` — permite seleccionar casos de distintos equipos |

---

## Flag de contingencia (frontend)

Ubicación sugerida: `src/frontend/js/components/reasignar-casos-modal.js`

```javascript
// CONTINGENCIA: true = listar todos los agentes ro (sin filtro por equipo del caso).
// false = comportamiento original (M-01 resolveTeam + M-02 fetchAgentsByTeam).
var REASSIGN_CROSS_TEAM_CONTINGENCY = true;
```

Al revertir la contingencia: poner el flag en `false` y descomentar la lógica original.

---

## Cambios por evento

### M-01 — Apertura del modal

| Modo | Comportamiento |
|---|---|
| **Original** (comentado) | `resolveTeam(cases[0])` → si falla, error y no carga agentes → `fetchAgentsByTeam(resolvedTeam)` |
| **Contingencia** (activo) | No bloquear por equipo indeterminado. Opcional: calcular `resolvedTeam` solo para mostrar en banner informativo. Disparar `fetchAllAgents()` → M-02-C |

Flujo actualizado: [flow-M01-cuando-abre-modal-reasignacion.md](./Flujos/flow-M01-cuando-abre-modal-reasignacion.md)

### M-02 — Carga de agentes

| Modo | Endpoint | Filtro |
|---|---|---|
| **Original** (comentado) | `GET /api/v1/equipo-operadores?team={resolvedTeam}&role=ro` | `general_user_team IN teamVariants` |
| **Contingencia** (activo) | `GET /api/v1/equipo-operadores?role=ro` | Solo rol + status activo |

Nuevo flujo documentado: [flow-M02-C-cuando-carga-todos-los-agentes.md](./Flujos/flow-M02-C-cuando-carga-todos-los-agentes.md)

Referencia backend (`case_detail_controller.GetOperadoresByTeam`): si `team == ""`, la cláusula `AND gu.general_user_team IN ?` **no se aplica**.

Referencia previa en el proyecto:

```javascript
// get_case_detail_sv.html — cargarOperadores / cargarAgentesRO
// TEMPORALMENTE DESHABILITADO: filtro por team para permitir reasignación entre equipos
// var teamCaso = self.teamSegunRiesgo();
// if (teamCaso) url += '?team=' + encodeURIComponent(teamCaso);
fetch(url)
```

### M-03 — Selección de agente

Sin cambio de lógica. En contingencia conviene mostrar el equipo del agente en el `<select>` para que el supervisor distinga destinos:

```html
<option :value="agent.icode">
  ${ agent.fullName } — ${ agent.team || 'Sin equipo' }
</option>
```

### M-04 / M-05

Sin cambios funcionales. M-05 ya persiste reasignación cross-team:

- `victim_case.agent_id` → siempre el agente elegido
- `victim_case_team` → solo se rellena si estaba vacío (no se sobrescribe equipo del caso)
- `follow_up_v2` pendientes → `agent_id` y `team` del agente destino

---

## Cambios de interfaz (`reasignar_casos_modal.html`)

| Elemento | Original | Contingencia |
|---|---|---|
| `.rcm-team-banner` | `"Equipo: { resolvedTeam }"` | Ocultar (`v-if="!crossTeamContingency && resolvedTeam"`) **o** texto informativo: *"Modo contingencia: puede elegir agente de cualquier equipo"* |
| Mensaje lista vacía | *"No hay agentes disponibles para este equipo."* | *"No hay agentes disponibles."* |
| `<select>` opciones | `agent.fullName` | `agent.fullName — agent.team` |
| Errores de carga | *"...del equipo"* | *"Error al cargar los agentes"* |

Detalle en [reasignar-casos-modal-interface.md](./reasignar-casos-modal-interface.md#modo-contingencia-cross-team).

---

## Selección de casos (E-14)

**Alcance mínimo (solo modal):** E-14 puede permanecer igual. El supervisor selecciona casos de un mismo `caseTeam` (como hoy) pero elige agente destino de **cualquier** equipo.

**Alcance ampliado (opcional):** si la contingencia también requiere reasignar en lote casos de equipos distintos, habría que comentar en `casos-component.js` la regla de mismo equipo (E-14 PASO 2) con el mismo patrón de flag. Documentar en [flow-E14-cuando-selecciona-caso-reasignacion.md](./Flujos/flow-E14-cuando-selecciona-caso-reasignacion.md).

---

## Checklist de implementación

### Frontend — `reasignar-casos-modal.js`

- [ ] Agregar `REASSIGN_CROSS_TEAM_CONTINGENCY = true`
- [ ] En `open()`: rama `if (REASSIGN_CROSS_TEAM_CONTINGENCY)` → `fetchAllAgents()`; `else` → lógica original comentada o en bloque `/* ORIGINAL */`
- [ ] Implementar `fetchAllAgents()` (M-02-C): `GET /api/v1/equipo-operadores?role=ro` sin `team`
- [ ] Comentar (no borrar) `resolveTeam` + llamada `fetchAgentsByTeam` en rama original
- [ ] Ajustar mensajes de error genéricos en rama contingencia
- [ ] Exponer `crossTeamContingency` en `data` o `computed` para el template

### Frontend — `reasignar_casos_modal.html`

- [ ] Condicionar `.rcm-team-banner`
- [ ] Mostrar equipo en opciones del select
- [ ] Actualizar texto de lista vacía

### Backend

- [ ] **Ningún cambio requerido** si se usa `GET /api/v1/equipo-operadores` sin `team`

### Opcional — `casos-component.js` (E-14)

- [x] Flag `REASSIGN_CROSS_TEAM_CONTINGENCY` — omite regla de mismo equipo en `toggleCaseSelection` y `toggleSelectAllCurrentPage`

### Reversión post-contingencia

- [ ] `REASSIGN_CROSS_TEAM_CONTINGENCY = false`
- [ ] Descomentar lógica M-01 / M-02 original
- [ ] Restaurar textos UI del modal

---

## Archivos de planeación actualizados

| Documento | Cambio |
|---|---|
| Este archivo | Plan maestro de contingencia |
| `reasignar-casos-modal-interface.md` | Sección modo contingencia + UI |
| `reasignar-casos-modal-events.md` | Evento M-02-C + nota contingencia |
| `Flujos/flow-M01-...` | Rama contingencia en PASO 3–4 |
| `Flujos/flow-M02-...` | Marcado como lógica original (comentar al activar contingencia) |
| `Flujos/flow-M02-C-...` | **Nuevo** — carga todos los agentes |
| `Flujos/flow-E14-...` | Nota alcance opcional |

---

## Riesgos / consideraciones operativas

1. **Seguimientos pendientes** pasan al `team` del agente destino aunque el `victim_case_team` del caso no cambie si ya tenía valor — comportamiento ya definido en M-05.
2. **Listado grande** de agentes en el `<select>`; aceptable para contingencia. Mejora futura: autocomplete (fuera de alcance).
3. **Coherencia con detalle de caso:** `get_case_detail_sv.html` ya usa el mismo patrón sin filtro por team; conviene mantener flags alineados durante la contingencia.
