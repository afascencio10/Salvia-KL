# Listado Reportes — Inventario de Eventos

Pantalla: `GET /salvia/reportes` (roles `op`, `sv`)
Template: `list_reportes.html` (nuevo)
API de listado: `GET /api/v1/victim-contacts/list` (nuevo, análogo a `/api/v1/cases/list`)

> Esta pantalla es una composición fina sobre `casos-component` + `reasignar-casos-modal` + `gestion-reporte-modal`. Solo se documentan aquí los eventos **propios de la pantalla** (qué prop pasa, cómo reacciona a lo que emiten sus hijos). La mecánica interna de filtros genéricos, búsqueda, paginación y ordenamiento está documentada una sola vez en `DocsMD/Componentes/Casos component/casos-component-events.md` y no se repite.

---

## Eventos identificados

### E-01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  El handler ListReportsGET verifica el permiso (op/sv; otro rol →
              401) y renderiza list_reportes.html. Monta casos-component con
              listEndpoint, columns, availableFilters, sortOptions,
              itemLabelPlural, buttons y reasignacion (solo para sv) propios de
              reportes (ver interface.md). casos-component dispara su propio
              E-01 contra el nuevo endpoint.
Requerido:    Sí
```

---

### E-02 — Cuando filtra por estado del reporte

📄 [Ver flujo → flow-E02-cuando-filtra-estado-reporte.md](./Flujos/flow-E02-cuando-filtra-estado-reporte.md)

```
Evento:       Cuando filtra por estado del reporte
Tipo:         User Interaction
Descripción:  Dropdown nuevo (Sin asignar / Asignado / Cerrado — válido /
              Cerrado — inválido). Extiende casos-component con el filtro
              estado_reporte (no existe para casos), combinable con el resto
              de filtros. Llama al backend con filter_estado_reporte.
Requerido:    Sí
```

---

### E-03 — Cuando filtra por asignación

📄 [Ver flujo → flow-E03-cuando-filtra-asignacion.md](./Flujos/flow-E03-cuando-filtra-asignacion.md)

```
Evento:       Cuando filtra por asignación
Tipo:         User Interaction
Descripción:  Autocomplete de agentes (reutiliza E-07/E-08 de casos-component)
              con una opción fija adicional "Sin asignar" que no dispara
              búsqueda de agentes y filtra victim_contact_agent_id IS NULL.
Requerido:    Sí
```

---

### E-04 — Cuando presiona "Ver Reporte" en una fila

📄 [Ver flujo → flow-E04-cuando-presiona-ver-reporte.md](./Flujos/flow-E04-cuando-presiona-ver-reporte.md)

```
Evento:       Cuando presiona "Ver Reporte" en una fila
Tipo:         User Interaction
Descripción:  casos-component emite action-clicked { buttonId:'ver', case }.
              La pantalla escucha, valida buttonId==='ver' y abre
              gestion-reporte-modal vía $refs.gestionModal.open(report).
Requerido:    Sí
```

---

### E-05 — Cuando presiona "Reasignar Reportes"

📄 [Ver flujo → flow-E05-cuando-presiona-reasignar-reportes.md](./Flujos/flow-E05-cuando-presiona-reasignar-reportes.md)

```
Evento:       Cuando presiona "Reasignar Reportes"
Tipo:         User Interaction
Descripción:  Solo visible para sv (botón nace de casos-component con
              reasignacion=true). casos-component emite reasignar-casos
              { cases: reportesSeleccionados }. La pantalla abre
              reasignar-casos-modal con entityKind='reportes'.
Requerido:    Sí
```

---

### E-06 — Cuando confirma la reasignación

📄 [Ver flujo → flow-E06-cuando-confirma-reasignacion.md](./Flujos/flow-E06-cuando-confirma-reasignacion.md)

```
Evento:       Cuando confirma la reasignación
Tipo:         User Interaction / Backend write
Descripción:  reasignar-casos-modal (extendido) emite reassigned tras llamar
              POST /api/v1/victim-contacts/reasignar-bulk. La pantalla recarga
              casos-component (reload()) para reflejar el nuevo agente/equipo.
Requerido:    Sí
```

---

### E-07 — Cuando busca por texto

📄 [Ver flujo → flow-E07-cuando-busca-por-texto.md](./Flujos/flow-E07-cuando-busca-por-texto.md)

```
Evento:       Cuando busca por texto
Tipo:         User Interaction
Descripción:  El mecanismo de UI (debounce 400ms, reset de página) es el ya
              documentado como E-03 de casos-component — no cambia. Lo que sí
              es nuevo es la traducción a SQL: victim_contact no tiene la
              forma de victim_case, así que el backend GET
              /api/v1/victim-contacts/list debe reutilizar la lógica ya
              probada de GetVictimContactsWithoutVictimCase (teléfono en
              Form1 o Form2, documento solo en Form1, nombres/apellidos
              directos) en vez de inventar una nueva.
Requerido:    Sí
```

---

### E-08 — Cuando ordena por fecha de reporte

📄 [Ver flujo → flow-E08-cuando-ordena-por-fecha.md](./Flujos/flow-E08-cuando-ordena-por-fecha.md)

```
Evento:       Cuando ordena por fecha de reporte
Tipo:         User Interaction
Descripción:  El mecanismo (reset de página, recarga con sort/order) es el ya
              documentado como E-04 de casos-component. Lo nuevo: el <select>
              hoy tiene 4 opciones escritas literal en el template
              (incluyendo "Próximo seguimiento", que no aplica a un reporte).
              Requiere el prop nuevo sortOptions para mostrar solo
              "Fecha de Reporte ASC/DESC" — cubre literal el pedido de
              FigJam ("Ordenar por fecha de reporte ascendente y
              descendente").
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿Carga inicial? → E-01
- [x] ¿Filtros nuevos específicos de reportes? → E-02, E-03
- [x] ¿Acción de ver detalle? → E-04
- [x] ¿Reasignación masiva (disparo + confirmación)? → E-05, E-06
- [x] ¿Búsqueda por texto (traducción específica a `victim_contact`)? → E-07
- [x] ¿Ordenamiento (opciones específicas de reporte)? → E-08
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No
- [x] ¿Hay scheduled tasks o webhooks? → No

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga la pantalla | Lifecycle | No (solo lee) |
| E-02 | Cuando filtra por estado del reporte | User Interaction | No (solo lee) |
| E-03 | Cuando filtra por asignación | User Interaction | No (solo lee) |
| E-04 | Cuando presiona "Ver Reporte" | User Interaction | No (navegación/apertura de modal) |
| E-05 | Cuando presiona "Reasignar Reportes" | User Interaction | No (apertura de modal) |
| E-06 | Cuando confirma la reasignación | User Interaction / Backend write | Sí |
| E-07 | Cuando busca por texto | User Interaction | No (solo lee) |
| E-08 | Cuando ordena por fecha de reporte | User Interaction | No (solo lee) |

**Total: 8 eventos propios de la pantalla — 1 Lifecycle, 6 User Interaction, 1 Backend write.**
