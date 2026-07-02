# barrier-follow-up-timeline — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [barrier-follow-up-timeline-interface.md](barrier-follow-up-timeline-interface.md) |
| Uso | [barrier-follow-up-timeline-usage.md](barrier-follow-up-timeline-usage.md) |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga el componente | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-el-componente.md) |
| E02 | Cuando cambia el prop barrierId | Lifecycle | — |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-el-componente.md](Flujos/flow-E01-cuando-carga-el-componente.md)

**Evento:** Cuando carga el componente
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue (`mounted`). Llama `GET /api/v1/barriers/:barrierId/follow-ups` para cargar los seguimientos de la barrera. Ordena los resultados por `createdAt` ascendente (más antiguo primero). Mientras carga, muestra el spinner. Si hay error, muestra el mensaje de error. Si la respuesta está vacía, muestra el empty state.
**Requerido:** Sí

---

**Evento:** Cuando cambia el prop barrierId
**Tipo:** Lifecycle
**Descripción:** Watcher sobre el prop `barrierId`. Si el padre cambia el ID de la barrera, el componente vuelve a llamar al API para recargar los seguimientos correspondientes al nuevo ID. Resetea `cargando`, `error` y `seguimientos` antes de la nueva carga.
**Requerido:** Sí

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos?
- [x] ¿Se cubre el cambio dinámico del prop?
- [x] ¿Los estados de carga, error y vacío están cubiertos?
- [ ] ¿Hay llamadas a API? — Sí: `GET /api/v1/barriers/:barrierId/follow-ups` (pendiente de crear en el backend)
