# Notificaciones — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [notificaciones-interface.md](./notificaciones-interface.md) |
| API paginada | [notificaciones-api.md](./notificaciones-api.md) |

---

## Roles con acceso

| Rol | Código | Alcance en pantalla |
|---|---|---|
| Agente de Seguimiento | `op` | Oficios asignados como `agent_id`; gestiona `por_proyectar` y `en_correccion` |
| Revisor Operativo | `ro` | Mismo alcance que `op` |
| Agente de Notificaciones | `an` | Oficios asignados como `notification_user_id`; gestiona revisión, aprobación, radicación y respuesta |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando cambia de tab | User Interaction | [📄 flow-E02-cuando-cambia-de-tab.md](./Flujos/flow-E02-cuando-cambia-de-tab.md) |
| E03 | Cuando aplica filtros | User Interaction | [📄 flow-E03-cuando-aplica-filtros.md](./Flujos/flow-E03-cuando-aplica-filtros.md) |
| E04 | Cuando cambia de página | User Interaction | [📄 flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md) |
| E05 | Cuando presiona "Gestionar" | User Interaction | [📄 flow-E05-cuando-presiona-gestionar.md](./Flujos/flow-E05-cuando-presiona-gestionar.md) |
| E06 | Cuando cancela el modal | User Interaction | — |
| E07 | Cuando confirma acción en modal | User Interaction | [📄 flow-E07-cuando-confirma-accion-modal.md](./Flujos/flow-E07-cuando-confirma-accion-modal.md) |
| E08 | Cuando selecciona departamento en modal proyectar | User Interaction | — |
| E09 | Cuando selecciona ciudad en modal proyectar | User Interaction | — |
| E10 | Cuando selecciona municipio en modal proyectar | User Interaction | — |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue. Verifica que el rol del usuario sea válido (`op`, `ro` o `an`), construye la URL del API con paginación server-side y carga la primera página de oficios desde BD.
**Requerido:** Sí

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** El usuario presiona el tab "Todos mis oficios" o "Oficios por gestionar". Cambia `currentTab`, resetea la página a 0 y recarga desde BD con `manageableOnly=true` si aplica.
**Requerido:** Sí

📄 [Ver flujo → flow-E02-cuando-cambia-de-tab.md](./Flujos/flow-E02-cuando-cambia-de-tab.md)

---

**Nombre del evento:** Cuando aplica filtros
**Tipo:** User Interaction
**Descripción:** El usuario escribe o selecciona en cualquiera de los 4 filtros (estado, identidad, entidad, radicado). Tras debounce de 350 ms, consulta BD con los filtros como query params y resetea la página a 0.
**Requerido:** Sí

📄 [Ver flujo → flow-E03-cuando-aplica-filtros.md](./Flujos/flow-E03-cuando-aplica-filtros.md)

---

**Nombre del evento:** Cuando cambia de página
**Tipo:** User Interaction
**Descripción:** El usuario presiona Anterior, Siguiente o un número de página. Cambia `currentPage`, hace scroll al inicio y consulta la nueva página en BD manteniendo tab y filtros activos.
**Requerido:** Sí

📄 [Ver flujo → flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md)

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
**Descripción:** El usuario presiona el botón de acción principal del modal. Valida los campos requeridos según el tipo de acción, construye el payload y llama a `PUT /api/v1/entity-letters/:id/action`. Si el API responde 200, recarga la página actual de oficios desde BD y cierra el modal.
**Requerido:** Sí

📄 [Ver flujo → flow-E07-cuando-confirma-accion-modal.md](./Flujos/flow-E07-cuando-confirma-accion-modal.md)

---

**Nombre del evento:** Cuando selecciona departamento en modal proyectar
**Tipo:** User Interaction
**Descripción:** Filtra `allCities` en memoria por `departmentId`. Resetea ciudad, municipio, entidad y sus listas de opciones. Sin fetch.
**Requerido:** Sí

---

**Nombre del evento:** Cuando selecciona ciudad en modal proyectar
**Tipo:** User Interaction
**Descripción:** Fetchea municipios del API (`GET /api/v1/locations/towns?city_id={cityId}`). Resetea municipio, entidad y sus listas.
**Requerido:** Sí

---

**Nombre del evento:** Cuando selecciona municipio en modal proyectar
**Tipo:** User Interaction
**Descripción:** Fetchea sedes del municipio (`GET /api/v1/entity-branches?town_code={townId}`). La respuesta incluye `id` (PK integer, FK a `salvia.entity_branch`), `icode` y `name`. El `<option>` usa `id` como valor. Resetea entidad. La opción "Otra entidad" (`value="otra"`) siempre está presente al final. Vue reactivamente muestra/oculta el input de nombre libre según la selección.
**Requerido:** Sí

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Todos los envíos de formulario están incluidos?
- [x] ¿Paginación y filtros se ejecutan en BD?
- [x] ¿Roles op, ro y an documentados?
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica en esta pantalla
