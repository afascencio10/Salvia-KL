# `casos-component` — Inventario de Eventos

Componente: `CasosComponent`  
Archivo fuente: `src/frontend/components/casos-component.js` *(pendiente de crear)*  
Usado por: Pantallas que necesitan visualizar y filtrar listas de casos

> **Nota de alcance:** Este componente es reutilizable. Recibe props para configurar columnas visibles, filtros disponibles, filtro inicial, botones de acción por fila y modo de reasignación. La interacción hacia afuera ocurre a través de dos eventos emitidos: `action-clicked` (botón de acción en una fila) y `reasignar-casos` (reasignación masiva de casos seleccionados). La lógica de filtrado, búsqueda, ordenamiento y selección de casos es interna al componente.

---

## Eventos identificados

---

### E-01 — Cuando carga el componente

📄 [Ver flujo → flow-E01-cuando-carga-componente.md](./Flujos/flow-E01-cuando-carga-componente.md)

```
Evento:       Cuando carga el componente
Tipo:         Lifecycle
Descripción:  Se ejecuta al montar el componente. Lee el prop :defaultFilter
              para determinar el filtro inicial. Llama al backend para obtener
              la lista de casos que cumplen ese filtro. Inicializa el estado
              interno: filtro activo, texto de búsqueda vacío y ordenamiento
              por defecto (fecha de registro desc). Renderiza la tabla con
              las columnas definidas en :columns, ocultando las que indica
              :hiddenColumns, y pinta los botones de acción definidos en
              :buttons en cada fila. Si :reasignacion es true, inicializa
              selectedCases = [] y prepara la columna de checkbox.
Requerido:    Sí
```

---

### E-02 — Índice de filtros chip y dropdown

📄 [Ver índice → flow-E02-cuando-selecciona-filtro.md](./Flujos/flow-E02-cuando-selecciona-filtro.md)

E-02 se dividió en tres flujos independientes:

| Sub-evento | Filtro | Flujo |
|---|---|---|
| **E-09** | Casos nuevos (chip) | [flow-E09](./Flujos/flow-E09-cuando-filtra-casos-nuevos.md) |
| **E-10** | Nivel de riesgo (dropdown) | [flow-E10](./Flujos/flow-E10-cuando-filtra-nivel-riesgo.md) |
| **E-11** | Por equipo (dropdown) | [flow-E11](./Flujos/flow-E11-cuando-filtra-equipo.md) |
| **E-12** | Seguimientos ejecutados (dropdown) | [flow-E12](./Flujos/flow-E12-cuando-filtra-seguimientos-ejecutados.md) |
| **E-13** | Estado del caso (dropdown) | [flow-E13](./Flujos/flow-E13-cuando-filtra-estado-caso.md) |

---

### E-09 — Cuando filtra por casos nuevos

📄 [Ver flujo → flow-E09-cuando-filtra-casos-nuevos.md](./Flujos/flow-E09-cuando-filtra-casos-nuevos.md)

```
Evento:       Cuando filtra por casos nuevos
Tipo:         User Interaction
Descripción:  Toggle del chip "Casos nuevos". Filtra casos creados desde hoy
              hasta 5 días calendario antes (victim_case_creation_date,
              zona America/Bogota). Llama al backend con chip_filter=casos_nuevos.
Requerido:    Sí
```

---

### E-10 — Cuando filtra por nivel de riesgo

📄 [Ver flujo → flow-E10-cuando-filtra-nivel-riesgo.md](./Flujos/flow-E10-cuando-filtra-nivel-riesgo.md)

```
Evento:       Cuando filtra por nivel de riesgo
Tipo:         User Interaction
Descripción:  Cambio en el dropdown de riesgo. Filtra por
              victim_case_form2_risk_level (1-4: bajo/moderado/alto/extremo).
              Llama al backend con filter_key=riesgo y filter_value.
Requerido:    Sí
```

---

### E-11 — Cuando filtra por equipo

📄 [Ver flujo → flow-E11-cuando-filtra-equipo.md](./Flujos/flow-E11-cuando-filtra-equipo.md)

```
Evento:       Cuando filtra por equipo
Tipo:         User Interaction
Descripción:  Cambio en el dropdown "Por equipo". Filtra por
              victim_case.victim_case_team (equipo del caso, no del agente).
              Llama al backend con filter_key=equipo y filter_value.
Requerido:    Sí
```

---

### E-12 — Cuando filtra por número de seguimientos ejecutados

📄 [Ver flujo → flow-E12-cuando-filtra-seguimientos-ejecutados.md](./Flujos/flow-E12-cuando-filtra-seguimientos-ejecutados.md)

```
Evento:       Cuando filtra por número de seguimientos ejecutados
Tipo:         User Interaction
Descripción:  Cambio en el dropdown numérico "Seguimientos ejecutados".
              Filtra casos cuyo conteo de follow_up_v2 con status REALIZADO
              coincide exactamente con el valor seleccionado (misma métrica
              que la columna completed_follow_ups en E-01).
              Llama al backend con dropdown_filter_key=seguimientos_ejecutados
              y dropdown_filter_value.
Requerido:    Sí
```

---

### E-13 — Cuando filtra por estado del caso

📄 [Ver flujo → flow-E13-cuando-filtra-estado-caso.md](./Flujos/flow-E13-cuando-filtra-estado-caso.md)

```
Evento:       Cuando filtra por estado del caso
Tipo:         User Interaction
Descripción:  Cambio en el dropdown "Estado del caso". Filtra por
              victim_case.victim_case_status (códigos ra/is/cd/ex/r/fc).
              Opciones quemadas en el padre; etiquetas según labelEstado
              de get_case_detail_sv. Combinable con otros filtros.
              Llama al backend con filter_estado_caso={código}.
Requerido:    Sí
```

---

### E-03 — Cuando el usuario escribe en el buscador

📄 [Ver flujo → flow-E03-cuando-escribe-buscador.md](./Flujos/flow-E03-cuando-escribe-buscador.md)

```
Evento:       Cuando el usuario escribe en el buscador
Tipo:         User Interaction
Descripción:  El usuario escribe en el campo de búsqueda (filter.type='search')
              para buscar casos por número de ID (victim_case_i_code) o por
              teléfono de la víctima. Con un debounce de 400ms, resetea
              currentPage a 1 y llama al backend combinando el filtro activo
              vigente con el texto de búsqueda como parámetro adicional.
              Si el campo queda vacío, hace la misma llamada sin el parámetro
              de búsqueda para restaurar la lista completa del filtro activo.
              No altera el filtro activo ni el ordenamiento.
Requerido:    Sí
```

---

### E-04 — Cuando el usuario cambia el ordenamiento

📄 [Ver flujo → flow-E04-cuando-cambia-ordenamiento.md](./Flujos/flow-E04-cuando-cambia-ordenamiento.md)

```
Evento:       Cuando el usuario cambia el ordenamiento
Tipo:         User Interaction
Descripción:  El usuario elige una de cuatro opciones explícitas en el
              select: fecha de registro o próximo seguimiento, cada una en
              ASC o DESC. Por defecto (E-01): Fecha de registro DESC.
              Resetea a página 1 y recarga casos desde el backend con
              sort y order; mantiene filtros activos.
Requerido:    Sí
```

---

### E-07 — Cuando el usuario escribe en el autocomplete de persona asignada

📄 [Ver flujo → flow-E07-cuando-escribe-autocomplete.md](./Flujos/flow-E07-cuando-escribe-autocomplete.md)

```
Evento:       Cuando el usuario escribe en el autocomplete de persona asignada
Tipo:         User Interaction
Descripción:  El usuario escribe en el input del filtro de tipo 'autocomplete'
              (persona_asignada). Con debounce de 400ms, llama al backend para
              buscar agentes cuyo nombre o apellido coincida con el texto
              ingresado. Muestra los resultados como sugerencias en el
              dropdown flotante debajo del input. Si el campo queda vacío,
              limpia las sugerencias sin llamar al backend.
Requerido:    Sí
```

---

### E-08 — Cuando el usuario selecciona o limpia una opción del autocomplete

📄 [Ver flujo → flow-E08-cuando-selecciona-autocomplete.md](./Flujos/flow-E08-cuando-selecciona-autocomplete.md)

```
Evento:       Cuando el usuario selecciona o limpia una opción del autocomplete
Tipo:         User Interaction
Descripción:  Se dispara en dos situaciones:
              (A) El usuario hace clic en una sugerencia del dropdown del
              autocomplete de persona_asignada. Establece la opción
              seleccionada como filtro activo, cierra el dropdown y llama al
              backend para recargar los casos filtrados por ese agente.
              (B) El usuario presiona "✕" en el SelectedTag. Limpia la
              selección, resetea el filtro al defaultFilter y llama al
              backend para recargar los casos sin ese filtro.
Requerido:    Sí
```

---

### E-06 — Cuando el usuario cambia de página

📄 [Ver flujo → flow-E06-cuando-cambia-pagina.md](./Flujos/flow-E06-cuando-cambia-pagina.md)

```
Evento:       Cuando el usuario cambia de página
Tipo:         User Interaction
Descripción:  El usuario presiona "← Anterior" o "Siguiente →" en la
              PaginationBar. Actualiza currentPage en el estado interno
              y realiza una nueva consulta al backend con el mismo filtro
              activo, ordenamiento y el nuevo número de página. Actualiza
              la tabla con los casos de la nueva página. Scroll al inicio
              de la tabla al cargar los nuevos resultados.
Requerido:    Sí
```

---

### E-05 — Cuando el usuario presiona un botón de acción en una fila

📄 [Ver flujo → flow-E05-cuando-presiona-boton-accion.md](./Flujos/flow-E05-cuando-presiona-boton-accion.md)

```
Evento:       Cuando el usuario presiona un botón de acción
Tipo:         User Interaction
Descripción:  El usuario presiona uno de los botones de la columna de
              acciones dentro de una fila de la tabla. Los botones
              disponibles son los definidos en el prop :buttons (array de
              objetos con id y etiqueta). El componente emite el evento
              'action-clicked' hacia el componente padre con el identificador
              del botón presionado y el objeto completo del caso de esa fila.
              El padre es quien decide qué acción ejecutar (navegar, abrir
              modal, etc.).
Requerido:    Sí
```

---

### E-14 — Cuando el usuario selecciona o deselecciona un caso para reasignación

📄 [Ver flujo → flow-E14-cuando-selecciona-caso-reasignacion.md](./Flujos/flow-E14-cuando-selecciona-caso-reasignacion.md)

```
Evento:       Cuando el usuario selecciona o deselecciona un caso para reasignación
Tipo:         User Interaction
Descripción:  Solo aplica si el prop :reasignacion es true. El usuario marca o
              desmarca el checkbox de una fila o el checkbox "Todos" del
              encabezado. Solo puede seleccionar casos de la página actual y
              todos deben compartir el mismo caseTeam; si intenta mezclar
              equipos, muestra alerta flotante sobre la tabla y rechaza la
              selección. Actualiza
              selectedCases y muestra u oculta el botón "Reasignar Casos".
              No emite eventos hacia el padre.
Requerido:    Condicional (solo si :reasignacion === true)
```

---

### E-15 — Cuando el usuario presiona "Reasignar Casos"

📄 [Ver flujo → flow-E15-cuando-presiona-reasignar-casos.md](./Flujos/flow-E15-cuando-presiona-reasignar-casos.md)

Modal asociado: [reasignar-casos-modal-interface.md](./reasignar-casos-modal-interface.md) · [reasignar-casos-modal-events.md](./reasignar-casos-modal-events.md)

```
Evento:       Cuando el usuario presiona "Reasignar Casos"
Tipo:         User Interaction
Descripción:  Solo aplica si el prop :reasignacion es true y hay al menos un
              caso en selectedCases. El usuario presiona el botón "Reasignar
              Casos" ubicado encima de la tabla. El componente emite el evento
              'reasignar-casos' hacia el padre con el array de casos
              seleccionados. El padre abre reasignar-casos-modal (M-01).
              El componente no realiza la reasignación ni llama al backend.
Requerido:    Condicional (solo si :reasignacion === true)
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia del componente está cubierta? → E-09, E-10, E-11, E-12, E-13, E-03, E-04, E-05, E-06, E-07, E-08, E-14, E-15
- [x] ¿Los eventos emitidos hacia el padre están cubiertos? → E-05, E-15
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga el componente | Lifecycle | No (solo lee) |
| E-02 | Índice de filtros chip/dropdown | — | — |
| E-09 | Cuando filtra por casos nuevos | User Interaction | No (solo lee) |
| E-10 | Cuando filtra por nivel de riesgo | User Interaction | No (solo lee) |
| E-11 | Cuando filtra por equipo | User Interaction | No (solo lee) |
| E-12 | Cuando filtra por seguimientos ejecutados | User Interaction | No (solo lee) |
| E-13 | Cuando filtra por estado del caso | User Interaction | No (solo lee) |
| E-03 | Cuando el usuario escribe en el buscador | User Interaction | No (solo lee) |
| E-04 | Cuando el usuario cambia el ordenamiento | User Interaction | No (solo lee, ORDER BY en servidor) |
| E-05 | Cuando el usuario presiona un botón de acción | User Interaction | No (emite hacia padre) |
| E-06 | Cuando el usuario cambia de página | User Interaction | No (solo lee) |
| E-07 | Cuando el usuario escribe en el autocomplete de persona asignada | User Interaction | No (solo lee — busca agentes) |
| E-08 | Cuando el usuario selecciona o limpia una opción del autocomplete | User Interaction | No (solo lee — filtra casos) |
| E-14 | Cuando el usuario selecciona o deselecciona un caso para reasignación | User Interaction | No (estado interno) |
| E-15 | Cuando el usuario presiona "Reasignar Casos" | User Interaction | No (emite hacia padre) |

**Total: 15 eventos — 1 Lifecycle, 13 User Interaction, 1 Índice**  
**Ningún evento escribe en el backend desde este componente.**  
**Las acciones resultantes de E-05 y E-15 son responsabilidad del componente padre que consume `casos-component`.**
