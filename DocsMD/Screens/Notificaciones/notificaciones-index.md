# Notificaciones — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [notificaciones-interface.md](./notificaciones-interface.md) |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando cambia de tab | User Interaction | — |
| E03 | Cuando aplica filtros | User Interaction | — |
| E04 | Cuando cambia de página | User Interaction | — |
| E05 | Cuando presiona "Gestionar" | User Interaction | [📄 flow-E05-cuando-presiona-gestionar.md](./Flujos/flow-E05-cuando-presiona-gestionar.md) |
| E06 | Cuando cancela el modal | User Interaction | — |
| E07 | Cuando confirma acción en modal | User Interaction | [📄 flow-E07-cuando-confirma-accion-modal.md](./Flujos/flow-E07-cuando-confirma-accion-modal.md) |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue. Verifica que el rol del usuario sea válido (`op` o `an`), construye la URL del API con el filtro de usuario correspondiente y carga la lista de oficios.
**Requerido:** Sí

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** El usuario presiona el tab "Todos mis oficios" o "Oficios por gestionar". Cambia `currentTab` y resetea la página a 0. Vue re-evalúa `filteredOficios` reactivamente mostrando solo los oficios del tab activo.
**Requerido:** Sí

---

**Nombre del evento:** Cuando aplica filtros
**Tipo:** User Interaction
**Descripción:** El usuario escribe o selecciona en cualquiera de los 4 filtros (estado, identidad, entidad, radicado). Vue filtra `filteredOficios` reactivamente y resetea la página a 0.
**Requerido:** Sí

---

**Nombre del evento:** Cuando cambia de página
**Tipo:** User Interaction
**Descripción:** El usuario presiona Anterior, Siguiente o un número de página en el paginador. Cambia `currentPage` y hace scroll al inicio de la página.
**Requerido:** Sí

---

**Nombre del evento:** Cuando presiona "Gestionar"
**Tipo:** User Interaction
**Descripción:** El usuario presiona el botón "Gestionar" en una fila de la tabla. Determina qué modal abrir según el estado actual del oficio, pre-carga el formulario del modal con los datos existentes y muestra el modal.
**Requerido:** Sí

📄 [Ver flujo → flow-E05-cuando-presiona-gestionar.md](./Flujos/flow-E05-cuando-presiona-gestionar.md)

---

**Nombre del evento:** Cuando cancela el modal
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Cancelar", la X del encabezado del modal o el overlay oscuro. Cierra el modal y limpia `selectedOficio`, `activeModal`, `saveError` y `modalForm.reasonCorrection`.
**Requerido:** Sí

---

**Nombre del evento:** Cuando confirma acción en modal
**Tipo:** User Interaction
**Descripción:** El usuario presiona el botón de acción principal del modal. Valida los campos requeridos según el tipo de acción, construye el payload y llama a `PATCH /api/v1/entity-letters/:id/action`. Si el API responde 200, actualiza el estado del oficio en la lista sin recargar la pantalla y cierra el modal.
**Requerido:** Sí

📄 [Ver flujo → flow-E07-cuando-confirma-accion-modal.md](./Flujos/flow-E07-cuando-confirma-accion-modal.md)

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Todos los envíos de formulario están incluidos?
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica en esta pantalla
