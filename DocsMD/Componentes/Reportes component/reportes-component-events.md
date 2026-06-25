# `reportes-component` — Inventario de Eventos

Componente: `ReportesComponent`  
Archivo fuente: `src/frontend/components/reportes-component.js` *(pendiente de crear)*  
Usado por: Pantallas que necesitan visualizar y filtrar listas de reportes de VBG

> **Nota de alcance:** Este componente es reutilizable. Recibe props para configurar columnas visibles, columnas ocultas y botones de acción opcionales por fila (tantos como defina el padre). Solo lista reportes **sin caso asociado** y con `victim_contact_status = 'v'`. La interacción hacia afuera ocurre exclusivamente a través del evento `action-clicked` — el componente no ejecuta acciones por sí mismo. Los filtros son dos inputs independientes (nombre y teléfono). La lógica de ordenamiento y paginación es interna.

---

## Eventos identificados

---

### E-01 — Cuando carga el componente

📄 [Ver flujo → flow-E01-cuando-carga-componente.md](./Flujos/flow-E01-cuando-carga-componente.md)

```
Evento:       Cuando carga el componente
Tipo:         Lifecycle
Descripción:  Se ejecuta al montar el componente. Llama al backend para obtener
              la lista paginada de reportes sin caso asociado (status = 'v')
              ordenados por fecha de registro DESC. Inicializa el estado interno:
              filtros de nombre y teléfono vacíos, ordenamiento por defecto
              (fecha de registro desc) y página 1. Renderiza la tabla con las
              columnas definidas en :columns, ocultando las que indica
              :hiddenColumns. Si :buttons tiene elementos, pinta la columna
              de acciones con un botón por cada entrada del array.
Requerido:    Sí
```

---

### E-02 — Cuando el usuario filtra por nombre o teléfono

📄 [Ver flujo → flow-E02-cuando-escribe-buscador.md](./Flujos/flow-E02-cuando-escribe-buscador.md)

```
Evento:       Cuando el usuario filtra por nombre o teléfono
Tipo:         User Interaction
Descripción:  El usuario escribe en uno o ambos inputs de filtro:
              (A) Nombre de la víctima — busca en nombres y apellidos.
              (B) Teléfono de la víctima — busca en victim_col_phone.
              Con debounce de 400ms, resetea currentPage a 1 y llama al
              backend enviando search_name y search_phone como parámetros
              independientes. Si ambos tienen valor, se combinan con AND.
              Si un campo queda vacío, no aplica filtro sobre ese campo.
              No altera el ordenamiento.
Requerido:    Sí
```

---

### E-03 — Cuando el usuario cambia el ordenamiento

📄 [Ver flujo → flow-E03-cuando-cambia-ordenamiento.md](./Flujos/flow-E03-cuando-cambia-ordenamiento.md)

```
Evento:       Cuando el usuario cambia el ordenamiento
Tipo:         User Interaction
Descripción:  El usuario elige entre fecha de registro ASC o DESC en el
              SortSelector. Resetea a página 1 y recarga reportes desde el
              backend con sort y order; mantiene los filtros de nombre y
              teléfono activos.
Requerido:    Sí
```

---

### E-04 — Cuando el usuario cambia de página

📄 [Ver flujo → flow-E04-cuando-cambia-pagina.md](./Flujos/flow-E04-cuando-cambia-pagina.md)

```
Evento:       Cuando el usuario cambia de página
Tipo:         User Interaction
Descripción:  El usuario presiona "← Anterior" o "Siguiente →" en la
              PaginationBar. Actualiza currentPage y realiza una nueva consulta
              al backend con los mismos filtros de nombre/teléfono, ordenamiento
              y el nuevo número de página. Scroll al inicio de la tabla al cargar
              los nuevos resultados.
Requerido:    Sí
```

---

### E-05 — Cuando el usuario presiona un botón de acción en una fila

📄 [Ver flujo → flow-E05-cuando-presiona-boton-accion.md](./Flujos/flow-E05-cuando-presiona-boton-accion.md)

```
Evento:       Cuando el usuario presiona un botón de acción
Tipo:         User Interaction
Descripción:  Solo aplica si el prop :buttons tiene al menos un elemento.
              El usuario presiona uno de los botones de la columna de acciones
              dentro de una fila. Los botones disponibles son los definidos en
              el prop :buttons — el padre puede pasar tantos como requiera.
              El componente emite el evento 'action-clicked' hacia el padre con
              el identificador del botón presionado y el objeto completo del
              reporte. El componente no interpreta buttonId ni ejecuta ninguna
              acción; toda la lógica es responsabilidad del padre.
Requerido:    Condicional (solo si :buttons.length > 0)
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia del componente está cubierta? → E-02, E-03, E-04, E-05
- [x] ¿Los eventos emitidos hacia el padre están cubiertos? → E-05
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga el componente | Lifecycle | No (solo lee) |
| E-02 | Cuando el usuario filtra por nombre o teléfono | User Interaction | No (solo lee) |
| E-03 | Cuando el usuario cambia el ordenamiento | User Interaction | No (solo lee, ORDER BY en servidor) |
| E-04 | Cuando el usuario cambia de página | User Interaction | No (solo lee) |
| E-05 | Cuando el usuario presiona un botón de acción | User Interaction | No (emite hacia padre) |

**Total: 5 eventos — 1 Lifecycle, 4 User Interaction**  
**Ningún evento escribe en el backend desde este componente.**  
**Las acciones resultantes de E-05 son responsabilidad del componente padre que consume `reportes-component`.**
