# case-task-history — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [case-task-history-interface.md](case-task-history-interface.md) |
| Uso | [case-task-history-usage.md](case-task-history-usage.md) |

**Changelogs:** [changelogJul2026.md](changelogJul2026.md)

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga el componente | Lifecycle | — |
| E02 | Cuando se abre con una tarea | User Interaction | [📄 Ver flujo](Flujos/flow-E02-cuando-se-abre-con-una-tarea.md) |
| E03 | Cuando cierra el modal | User Interaction | — |

---

## Inventario de eventos

**Evento:** Cuando carga el componente
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue (`mounted`). Solo inicializa el estado (`visible = false`, `tarea = null`) — a diferencia de `case-task-modal`, este componente no tiene dropdowns ni catálogos de ubicación que precargar, así que no hace ninguna llamada al API al montarse.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E02-cuando-se-abre-con-una-tarea.md](Flujos/flow-E02-cuando-se-abre-con-una-tarea.md)

**Evento:** Cuando se abre con una tarea
**Tipo:** User Interaction
**Descripción:** El componente padre (`case-tasks.js`, listado de tareas del caso) llama al método público `open(taskId)`. El componente muestra el modal y hace `GET /api/v1/case-tasks/:taskId` para cargar la tarea, ya enriquecida con `assignedUserName`. Una vez cargada, renderiza su detalle según `tarea.type`.
**Requerido:** Sí

---

**Evento:** Cuando cierra el modal
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Cerrar", el botón "✕" del header o hace click en el overlay oscuro. Setea `visible = false` y limpia `tarea` y `error`. No emite ningún evento al padre.
**Requerido:** Sí

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos? — sí, vía `open(taskId)` (E02); el montaje (E01) no carga nada porque no hay catálogos
- [x] ¿Todos los campos interactivos del formulario tienen evento? — N/A, componente de solo lectura sin campos interactivos
- [x] ¿Las validaciones de formulario están cubiertas? — N/A, no hay formulario
- [x] ¿Los estados de error de carga y guardado están cubiertos? — error de carga sí (E02); no hay operación de guardado
- [x] ¿El callback al padre está cubierto? — N/A, el componente no emite eventos
