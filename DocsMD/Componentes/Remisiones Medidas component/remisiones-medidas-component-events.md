# `remisiones-medidas-component` — Inventario de Eventos

Componente: `remisiones-medidas-component`  
Archivo fuente: `src/frontend/js/components/remisiones-medidas-component.js`  
Usado por: [Remisiones Medidas](../../Screens/Remisiones%20Medidas/remisiones-medidas-index.md) y [Mis Remisiones Medidas](../../Screens/Mis%20Remisiones%20Medidas/mis-remisiones-medidas-index.md)

> **Nota de alcance:** El componente ya pinta las cinco vistas con un seed en memoria. Estos eventos describen el contrato al que se cablea: una remisión (`emergency_measure_referral`) con varias medidas (`emergency_measure.referral_id`). Los filtros se combinan con AND y se aplican después del recorte de `vista`. Las cards y el banner salen de ese recorte, no de los filtros. Reasignar y la revisión técnica persisten en la remisión. "Ver caso" y "Ver detalle" solo se emiten.
>
> Los modales viven en el mismo componente, así que sus eventos están en este inventario.
>
> **Prerequisito:** [M-01 — Cabecera de remisión y `referral_id`](./Flujos/flow-M01-migracion-emergency-measure-referral.md)

Interfaz detallada: [remisiones-medidas-component-interface.md](./remisiones-medidas-component-interface.md)

---

## Eventos de infraestructura

---

### M-01 — Cabecera de remisión y `referral_id`

📄 [Ver flujo → flow-M01-migracion-emergency-measure-referral.md](./Flujos/flow-M01-migracion-emergency-measure-referral.md)

```
Evento:       Migración — emergency_measure_referral y referral_id
Tipo:         Infrastructure / Backend
Descripción:  Crear salvia.emergency_measure_referral con estado, urgencia,
              agente, remitente y la revisión técnica. Agregar referral_id a
              salvia.emergency_measure. El status de la medida no es el de la
              bandeja. Definir constantes snake_case de estado, urgencia y
              resultado de revisión.
Requerido:    Sí
Estado:       Propuesto — no implementado
```

---

## Eventos identificados

---

### E-01 — Cuando carga el componente

📄 [Ver flujo → flow-E01-cuando-carga-componente.md](./Flujos/flow-E01-cuando-carga-componente.md)

```
Evento:       Cuando carga el componente
Tipo:         Lifecycle
Descripción:  Lee vista y agenteId. Pide el listado con el recorte de la vista,
              sin filtros de UI. Pinta 8 cards (supervisor), 6 cards (agente)
              o el banner de las otras tres vistas. Ordena riesgo inminente
              primero y, en empate, referred_at descendente.
Requerido:    Sí
Prerequisito: M-01
```

---

### E-02 — Cuando escribe nombre o documento

📄 [Ver flujo → flow-E02-cuando-escribe-busqueda.md](./Flujos/flow-E02-cuando-escribe-busqueda.md)

```
Evento:       Cuando el usuario escribe en el buscador
Tipo:         User Interaction
Descripción:  Debounce 400ms. q contra nombre completo de la víctima (ILIKE)
              o victim_case_victim_doc_number (ILIKE). Se combina con el resto
              de filtros. No mueve cards ni banner.
Requerido:    Sí
```

---

### E-03 — Cuando filtra por estado

📄 [Ver flujo → flow-E03-cuando-filtra-estado.md](./Flujos/flow-E03-cuando-filtra-estado.md)

```
Evento:       Cuando el usuario elige un estado
Tipo:         User Interaction
Descripción:  Solo supervisor (16 estados) y agente (5 estados de su bandeja).
              Envía el código snake_case. Equipo técnico, convenio y Cruz Roja
              no muestran el select.
Requerido:    Sí
```

---

### E-04 — Cuando filtra por agente asignado

📄 [Ver flujo → flow-E04-cuando-filtra-agente.md](./Flujos/flow-E04-cuando-filtra-agente.md)

```
Evento:       Cuando el usuario elige agente asignado
Tipo:         User Interaction
Descripción:  Solo supervisor. agent_id con el icode, o el literal sin-asignar
              (assigned_agent_id IS NULL). El catálogo sale de GET .../agentes.
Requerido:    Sí
```

---

### E-05 — Cuando filtra por urgencia

📄 [Ver flujo → flow-E05-cuando-filtra-urgencia.md](./Flujos/flow-E05-cuando-filtra-urgencia.md)

```
Evento:       Cuando el usuario elige urgencia
Tipo:         User Interaction
Descripción:  En las cinco vistas. urgency = riesgo_inminente o normal.
Requerido:    Sí
```

---

### E-06 — Cuando limpia los filtros

📄 [Ver flujo → flow-E06-cuando-limpia-filtros.md](./Flujos/flow-E06-cuando-limpia-filtros.md)

```
Evento:       Cuando el usuario presiona "Limpiar" o "Limpiar filtros"
Tipo:         User Interaction
Descripción:  Vacía búsqueda, estado, agente y urgencia. Conserva la vista y,
              en agente, agenteId. Vuelve a pedir el listado del recorte.
Requerido:    Sí
```

---

### E-07 — Cuando selecciona una fila o todas

📄 [Ver flujo → flow-E07-cuando-selecciona-remision.md](./Flujos/flow-E07-cuando-selecciona-remision.md)

```
Evento:       Cuando el usuario marca una fila o el checkbox del encabezado
Tipo:         User Interaction
Descripción:  Solo supervisor. La selección es local. "Todas" marca las filas
              visibles; si ya estaban todas marcadas, vacía la selección entera.
              Muestra la barra "↺ Reasignar N seleccionada(s)".
Requerido:    Sí
```

---

### E-08 — Cuando abre reasignar

📄 [Ver flujo → flow-E08-cuando-abre-reasignar.md](./Flujos/flow-E08-cuando-abre-reasignar.md)

```
Evento:       Cuando abre el modal de reasignación
Tipo:         User Interaction
Descripción:  Desde la barra (usa la selección actual) o desde "Asignar agente"
              (deja seleccionada solo esa fila). Carga GET .../agentes.
              No persiste.
Requerido:    Sí
```

---

### E-09 — Cuando cancela reasignar

📄 [Ver flujo → flow-E09-cuando-cancela-reasignar.md](./Flujos/flow-E09-cuando-cancela-reasignar.md)

```
Evento:       Cuando cancela el modal de reasignación
Tipo:         User Interaction
Descripción:  Cierra el modal y limpia el select. No toca assigned_agent_id
              ni el estado. La selección de filas se queda como estaba al abrir.
Requerido:    Sí
```

---

### E-10 — Cuando confirma reasignar

📄 [Ver flujo → flow-E10-cuando-confirma-reasignar.md](./Flujos/flow-E10-cuando-confirma-reasignar.md)

```
Evento:       Cuando confirma la reasignación
Tipo:         User Interaction / Backend write
Descripción:  POST .../reasignar con los ids seleccionados y el agentId.
              Cada remisión queda con assigned_agent_id y status asignada.
              Cierra el modal, vacía la selección y recarga E-01.
Requerido:    Sí
```

---

### E-11 — Cuando abre la revisión técnica

📄 [Ver flujo → flow-E11-cuando-abre-revision.md](./Flujos/flow-E11-cuando-abre-revision.md)

```
Evento:       Cuando abre el modal de revisión
Tipo:         User Interaction
Descripción:  Botón "Revisar (Equipo Técnico)". Solo supervisor y equipo
              técnico, estado pendiente_de_revision. Carga el catálogo de
              profesionales. No persiste.
Requerido:    Sí
```

---

### E-12 — Cuando cancela la revisión

📄 [Ver flujo → flow-E12-cuando-cancela-revision.md](./Flujos/flow-E12-cuando-cancela-revision.md)

```
Evento:       Cuando cancela el modal de revisión
Tipo:         User Interaction
Descripción:  Cierra y limpia profesional, concepto y resultado. No cambia
              la remisión.
Requerido:    Sí
```

---

### E-13 — Cuando confirma la revisión

📄 [Ver flujo → flow-E13-cuando-confirma-revision.md](./Flujos/flow-E13-cuando-confirma-revision.md)

```
Evento:       Cuando guarda la revisión técnica
Tipo:         User Interaction / Backend write
Descripción:  POST .../{id}/revision. Guarda profesional, concepto, resultado
              y la fecha del día. El resultado pasa el estado a por_asignar,
              devuelta_para_justificacion o devuelta_a_seguimiento.
              Cierra el modal y recarga E-01.
Requerido:    Sí
```

---

### E-14 — Cuando presiona "Ver caso"

📄 [Ver flujo → flow-E14-cuando-presiona-ver-caso.md](./Flujos/flow-E14-cuando-presiona-ver-caso.md)

```
Evento:       Cuando el usuario presiona "Ver caso"
Tipo:         User Interaction
Descripción:  Emite ver-caso con caseICode y remisionId. El componente no
              navega. Los padres actuales no redirigen: la ruta de destino
              de estos casos todavía no está definida.
Requerido:    Sí
```

---

### E-15 — Cuando presiona "Ver detalle" o el botón de acción

📄 [Ver flujo → flow-E15-cuando-presiona-ver-detalle.md](./Flujos/flow-E15-cuando-presiona-ver-detalle.md)

```
Evento:       Cuando el usuario abre el detalle de la remisión
Tipo:         User Interaction
Descripción:  Un solo evento para el link "Ver detalle →" y para los botones
              de fila de tipo detalle (decidir vía, escalar, entregar medida,
              caracterización, trazabilidad, cotizaciones, cargar cotizaciones,
              confirmar seguridad). Emite ver-detalle con remisionId y
              caseICode. El negocio de cada botón queda en la pantalla de
              detalle, que todavía no existe.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿Migración de schema documentada? → M-01
- [x] ¿Ciclo de vida inicial cubierto? → E-01
- [x] ¿Filtros cubiertos y combinables (AND)? → E-02, E-03, E-04, E-05, E-06
- [x] ¿Selección masiva? → E-07, solo supervisor
- [x] ¿Reasignar? → E-08, E-09, E-10
- [x] ¿Revisión técnica? → E-11, E-12, E-13
- [x] ¿Navegación emitida al padre? → E-14, E-15
- [x] ¿Cards y banner sobre el recorte, no sobre los filtros? → E-01
- [x] ¿Flujos detallados escritos? → M-01, E-01 … E-15
- [ ] ¿M-01 y los endpoints implementados? → Pendiente
- [ ] ¿Mapeo currentRole → vista? → Pendiente
- [ ] ¿Pantalla de detalle de la remisión? → Pendiente (E-15 solo emite)

---

## Resumen

| # | Evento | Tipo | Persiste |
|---|---|---|---|
| M-01 | Cabecera de remisión y `referral_id` | Infrastructure | Sí (DDL, propuesto) |
| E-01 | Cuando carga el componente | Lifecycle | No (solo lee) |
| E-02 | Cuando escribe nombre o documento | User Interaction | No (solo lee) |
| E-03 | Cuando filtra por estado | User Interaction | No (solo lee) |
| E-04 | Cuando filtra por agente asignado | User Interaction | No (solo lee) |
| E-05 | Cuando filtra por urgencia | User Interaction | No (solo lee) |
| E-06 | Cuando limpia los filtros | User Interaction | No (solo lee) |
| E-07 | Cuando selecciona una fila o todas | User Interaction | No |
| E-08 | Cuando abre reasignar | User Interaction | No (solo lee catálogo) |
| E-09 | Cuando cancela reasignar | User Interaction | No |
| E-10 | Cuando confirma reasignar | User Interaction / Backend write | Sí |
| E-11 | Cuando abre la revisión técnica | User Interaction | No (solo lee catálogo) |
| E-12 | Cuando cancela la revisión | User Interaction | No |
| E-13 | Cuando confirma la revisión | User Interaction / Backend write | Sí |
| E-14 | Cuando presiona "Ver caso" | User Interaction | No (emite) |
| E-15 | Cuando presiona "Ver detalle" o el botón de acción | User Interaction | No (emite) |
