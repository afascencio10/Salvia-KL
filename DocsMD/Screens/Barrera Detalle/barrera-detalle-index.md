# Barrera Detalle — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [barrera-detalle-interface.md](barrera-detalle-interface.md) |

**Changelogs:** [changelogJul2026.md](changelogJul2026.md)

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando cambia de tab | User Interaction | — |
| E03 | Cuando el Enlace registra una gestión propia | User Interaction | [📄 Ver flujo](../../Otros/temp/req-registrar-gestion-propia-enlace.md) |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

**Evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** `mounted()` muestra el div `#app` y hace `GET /api/v1/barriers-v2/:id/detail` para cargar la barrera, datos de la víctima y ubicación. `barrierICode` viene del path param inyectado por `BarreraDetalleFacade.go`.
**Requerido:** Sí

---

**Evento:** Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** Actualiza `activeTab` ('info' | 'tareas' | 'timeline'). El tab "Tareas" monta `<case-tasks>` (pendientes/completadas + gestión propia); el tab "Timeline" monta `<case-timeline>`.
**Requerido:** Sí

---

📄 [Ver flujo → req-registrar-gestion-propia-enlace.md](../../Otros/temp/req-registrar-gestion-propia-enlace.md)

**Evento:** Cuando el Enlace registra una gestión propia
**Tipo:** User Interaction
**Descripción:** Aplica solo al rol `en` (Enlace Territorial) sobre una barrera de su mismo departamento asignado. Desde el tab "Tareas", completa tipo + descripción en el modal y confirma. El frontend llama `POST /api/v1/case-tasks/gestion-propia`, que crea la `case_task` ya en `Done` (sin pasar por `ToDo`), registra un evento de timeline, y transiciona la barrera de `OPEN` a `"En Gestion"` si aplicaba. Implementado dentro del componente reutilizable `case-tasks.js` (sin doc propio en `Componentes/`).
**Requerido:** Sí

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos?
- [x] ¿Todos los botones de la UI tienen evento?
- [x] ¿Hay llamadas a API? — Sí, pantalla completamente integrada al backend.
