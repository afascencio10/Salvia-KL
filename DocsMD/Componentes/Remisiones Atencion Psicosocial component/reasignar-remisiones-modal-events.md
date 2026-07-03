# `reasignar-remisiones-modal` — Inventario de Eventos

Componente: `ReasignarRemisionesModal`  
Archivo fuente: `src/frontend/js/components/reasignar-remisiones-modal.js` *(pendiente de crear)*  
Usado por: Pantallas padre con `remisiones-psicosocial-component` y `:reasignacion="true"` (p. ej. Historial de Remisiones — rol `sv`)

> **Nota de alcance:** Modal independiente del listado. Se abre cuando el padre recibe `reasignar-remisiones` (E-15). RRM-05 persiste en `psychosocial_support` y `team_contact` (solo contactos con `is_completed = false`).

Interfaz detallada: [reasignar-remisiones-modal-interface.md](./reasignar-remisiones-modal-interface.md)

---

## Eventos identificados

### RRM-01 — Cuando abre el modal de reasignación

📄 [Ver flujo → flow-RRM01-cuando-abre-modal-reasignacion.md](./Flujos/flow-RRM01-cuando-abre-modal-reasignacion.md)

```
Evento:       Cuando abre el modal de reasignación
Tipo:         Lifecycle / User Interaction
Descripción:  El padre llama open(remisiones) tras E-15. Visible = true,
              guarda remisiones, resetea switch y selección, dispara RRM-03
              (carga profesionales — modo por defecto).
Requerido:    Sí
```

---

### RRM-02 — Cuando alterna el switch "Asignar en dupla"

📄 [Ver flujo → flow-RRM02-cuando-alterna-asignar-en-dupla.md](./Flujos/flow-RRM02-cuando-alterna-asignar-en-dupla.md)

```
Evento:       Cuando alterna "Asignar en dupla"
Tipo:         User Interaction
Descripción:  Cambia asignarEnDupla. Limpia selección del select anterior
              y dispara RRM-03 con el modo correspondiente (profesionales o duplas).
Requerido:    Sí
```

---

### RRM-03 — Cuando carga opciones del select

📄 [Ver flujo → flow-RRM03-cuando-carga-opciones-asignacion.md](./Flujos/flow-RRM03-cuando-carga-opciones-asignacion.md)

```
Evento:       Cuando carga opciones del select
Tipo:         Backend read
Descripción:  SI asignarEnDupla === false → GET profesionales-reasignacion (ps/ts agrupados).
              SI asignarEnDupla === true  → GET duplas/reasignacion (label enriquecido).
Requerido:    Sí
```

---

### RRM-04 — Cuando cancela o cierra el modal

📄 [Ver flujo → flow-RRM04-cuando-cancela-modal.md](./Flujos/flow-RRM04-cuando-cancela-modal.md)

```
Evento:       Cuando cancela o cierra el modal
Tipo:         User Interaction
Descripción:  Cancelar, ✕ o backdrop. Cierra sin guardar, limpia estado,
              emite 'closed'. No recarga tabla.
Requerido:    Sí
```

---

### RRM-05 — Cuando confirma la reasignación (guardado)

📄 [Ver flujo → flow-RRM05-cuando-confirma-reasignacion.md](./Flujos/flow-RRM05-cuando-confirma-reasignacion.md)

```
Evento:       Cuando confirma la reasignación
Tipo:         User Interaction / Backend write
Descripción:  POST reasignar-bulk. Actualiza professional_id o dupla_id en
              psychosocial_support y en team_contact pendientes (is_completed=false).
              Emite 'reassigned'; padre recarga listado.
Requerido:    Sí
```

---

## Integración E-15 → Modal

```
[E-15 ReassignBtn] → emit('reasignar-remisiones')
       │
       ▼
[Padre] reasignarModal.open(remisiones)  → RRM-01
```

---

## Checklist de completitud

- [x] Apertura desde E-15 → RRM-01
- [x] Switch dupla / profesional → RRM-02 + RRM-03
- [x] Select agrupado ps/ts (Image 2) → RRM-03
- [x] Select duplas enriquecidas (Image 3) → RRM-03
- [x] Cancelar sin guardar → RRM-04
- [x] Persistencia psychosocial_support + team_contact → RRM-05
- [x] Recarga tabla padre → RRM-05 + reload() del componente

---

## Resumen

| # | Evento | Tipo | Persiste |
|---|---|---|---|
| RRM-01 | Abre modal | Lifecycle | No |
| RRM-02 | Alterna switch dupla | User Interaction | No |
| RRM-03 | Carga opciones select | Backend read | No |
| RRM-04 | Cancela / cierra | User Interaction | No |
| RRM-05 | Confirma reasignación | Backend write | Sí |

**Total: 5 eventos — 1 Lifecycle, 2 User Interaction, 1 Backend read, 1 Backend write**
