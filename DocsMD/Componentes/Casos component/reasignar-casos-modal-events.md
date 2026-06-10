# `reasignar-casos-modal` — Inventario de Eventos

Componente: `ReasignarCasosModal`  
Archivo fuente: `src/frontend/js/components/reasignar-casos-modal.js` *(pendiente de crear)*  
Usado por: Pantallas padre que consumen `casos-component` con `:reasignacion="true"` (p. ej. Lista de casos)

> **Nota de alcance:** Este modal es independiente de `casos-component`. Se abre desde el padre cuando recibe `reasignar-casos` (E-15). M-05 persiste la reasignación en `victim_case`, `follow_up_v2` y `case_timeline_event`.

---

## Eventos identificados

---

### M-01 — Cuando abre el modal de reasignación

📄 [Ver flujo → flow-M01-cuando-abre-modal-reasignacion.md](./Flujos/flow-M01-cuando-abre-modal-reasignacion.md)

```
Evento:       Cuando abre el modal de reasignación
Tipo:         Lifecycle / User Interaction
Descripción:  El padre llama open(cases) tras recibir reasignar-casos (E-15).
              El modal se hace visible, guarda los casos, resuelve el equipo
              (caseTeam o derivado por riesgo) y dispara la carga de agentes (M-02).
Requerido:    Sí
```

---

### M-02 — Cuando carga agentes por equipo

📄 [Ver flujo → flow-M02-cuando-carga-agentes-por-equipo.md](./Flujos/flow-M02-cuando-carga-agentes-por-equipo.md)

```
Evento:       Cuando carga agentes por equipo
Tipo:         Lifecycle / Backend read
Descripción:  Tras resolver resolvedTeam (M-01), consulta al backend los agentes
              activos (rol ro) cuyo general_user_team coincide con el equipo.
              Obtiene el nombre desde general_user_profile. Puebla el <select>
              de "Nueva persona asignada".
Requerido:    Sí
```

---

### M-03 — Cuando selecciona un agente

📄 [Ver flujo → flow-M03-cuando-selecciona-agente.md](./Flujos/flow-M03-cuando-selecciona-agente.md)

```
Evento:       Cuando selecciona un agente
Tipo:         User Interaction
Descripción:  El usuario elige un agente en el <select>. Actualiza
              selectedAgentIcode y selectedAgent en el estado interno.
              Habilita el botón "Reasignar" (la acción de guardado es M-05).
Requerido:    Sí
```

---

### M-04 — Cuando cancela o cierra el modal

📄 [Ver flujo → flow-M04-cuando-cancela-modal.md](./Flujos/flow-M04-cuando-cancela-modal.md)

```
Evento:       Cuando cancela o cierra el modal
Tipo:         User Interaction
Descripción:  El usuario presiona "Cancelar", la "✕" del encabezado o el backdrop.
              Cierra el modal sin guardar, limpia el estado interno y emite
              'closed' hacia el padre. No recarga la tabla.
Requerido:    Sí
```

---

### M-05 — Cuando confirma la reasignación (guardado)

📄 [Ver flujo → flow-M05-cuando-confirma-reasignacion.md](./Flujos/flow-M05-cuando-confirma-reasignacion.md)

```
Evento:       Cuando confirma la reasignación
Tipo:         User Interaction / Backend write
Descripción:  El usuario presiona "Reasignar". El modal llama a
              POST /api/v1/cases/reasignar-bulk con los i_code de los casos
              y el agent_icode seleccionado. El backend (en transacción, por
              cada caso): actualiza agent_id en victim_case; asigna
              victim_case_team solo si estaba vacío; actualiza agent_id en
              follow_up_v2 con status <> REALIZADO; crea un CaseTimelineEvent
              con Category "Seguimientos", Type "Reasignación de Caso" y
              ActorID/ActorName del supervisor logueado. Al éxito emite
              'reassigned' y el padre recarga la tabla.
Requerido:    Sí
```

---

## Eventos emitidos hacia el padre

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `closed` | Modal cerrado sin reasignación exitosa (M-04) | `{}` |
| `reassigned` | Reasignación exitosa (M-05) | `{ cases, agent, updated }` |

---

## Checklist de completitud

- [x] ¿La apertura del modal desde E-15 está cubierta? → M-01
- [x] ¿La carga de agentes por equipo está cubierta? → M-02
- [x] ¿La selección de agente en el select está cubierta? → M-03
- [x] ¿La cancelación / cierre sin guardar está cubierta? → M-04
- [x] ¿El guardado de la reasignación está cubierto? → M-05
- [x] ¿La recarga de la tabla tras éxito está cubierta? → Responsabilidad padre + M-05

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| M-01 | Cuando abre el modal | Lifecycle | No |
| M-02 | Cuando carga agentes por equipo | Backend read | No (solo lee) |
| M-03 | Cuando selecciona un agente | User Interaction | No (estado interno) |
| M-04 | Cuando cancela o cierra el modal | User Interaction | No |
| M-05 | Cuando confirma la reasignación | User Interaction / Backend write | Sí |

**Total: 5 eventos — 1 Lifecycle, 1 Backend read, 2 User Interaction, 1 Backend write**
