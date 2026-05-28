# `casos-component` — Inventario de Eventos

Componente: `CasosComponent`  
Archivo fuente: `src/frontend/components/casos-component.js` *(pendiente de crear)*  
Usado por: Pantallas que necesitan visualizar y filtrar listas de casos

> **Nota de alcance:** Este componente es reutilizable. Recibe props para configurar columnas visibles, filtros disponibles, filtro inicial y botones de acción por fila. Toda la interacción hacia afuera ocurre a través de un único evento emitido cuando el usuario presiona un botón de acción. La lógica de filtrado, búsqueda y ordenamiento es interna al componente.

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
              :buttons en cada fila.
Requerido:    Sí
```

---

### E-02 — Cuando el usuario selecciona un filtro

```
Evento:       Cuando el usuario selecciona un filtro
Tipo:         User Interaction
Descripción:  El usuario hace clic en una de las opciones de la barra de
              filtros: Casos nuevos, Por nivel de riesgo, Por equipo, o
              Persona que tiene caso asignado. Actualiza el filtro activo
              en el estado interno y realiza una nueva consulta al backend
              con los parámetros del filtro seleccionado. Resetea el texto
              de búsqueda. Mantiene el ordenamiento activo vigente. Actualiza
              la tabla con los casos resultantes.
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

```
Evento:       Cuando el usuario cambia el ordenamiento
Tipo:         User Interaction
Descripción:  El usuario selecciona uno de los criterios de ordenamiento
              disponibles: por fecha del primer seguimiento sin ejecutar,
              o por fecha de registro. Reordena la lista de casos actualmente
              visible del lado del cliente, sin disparar una nueva consulta
              al backend. Si el mismo criterio ya está activo, invierte el
              orden (ascendente / descendente).
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

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia del componente está cubierta? → E-02, E-03, E-04, E-05, E-06, E-07, E-08
- [x] ¿Los eventos emitidos hacia el padre están cubiertos? → E-05
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga el componente | Lifecycle | No (solo lee) |
| E-02 | Cuando el usuario selecciona un filtro chip o dropdown | User Interaction | No (solo lee) |
| E-03 | Cuando el usuario escribe en el buscador | User Interaction | No (solo lee) |
| E-04 | Cuando el usuario cambia el ordenamiento | User Interaction | No (ordenamiento local) |
| E-05 | Cuando el usuario presiona un botón de acción | User Interaction | No (emite hacia padre) |
| E-06 | Cuando el usuario cambia de página | User Interaction | No (solo lee) |
| E-07 | Cuando el usuario escribe en el autocomplete de persona asignada | User Interaction | No (solo lee — busca agentes) |
| E-08 | Cuando el usuario selecciona o limpia una opción del autocomplete | User Interaction | No (solo lee — filtra casos) |

**Total: 8 eventos — 1 Lifecycle, 7 User Interaction**  
**Ningún evento escribe en el backend desde este componente.**  
**La acción resultante del botón presionado (E-05) es responsabilidad del componente padre que consume `casos-component`.**
