# case-entities — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [case-entities-interface.md](case-entities-interface.md) |
| Uso | [case-entities-usage.md](case-entities-usage.md) |
| Análisis de BD (tablas existentes + propuesta) | [related-tables.md](related-tables.md) |

---

## Resumen de eventos

> A pedido del usuario, **todos** los eventos tienen flujo documentado, incluidos los triviales (fuera de lo que normalmente pide `dev-docs-system.md`).

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga el componente | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-el-componente.md) |
| E02 | Cuando escribe en el buscador | User Interaction | [📄 Ver flujo](Flujos/flow-E02-cuando-escribe-en-el-buscador.md) |
| E03 | Cuando selecciona un sector en el filtro | User Interaction | [📄 Ver flujo](Flujos/flow-E03-cuando-selecciona-sector-filtro.md) |
| E04 | Cuando selecciona el filtro "¿Con barreras activas?" | User Interaction | [📄 Ver flujo](Flujos/flow-E04-cuando-selecciona-filtro-barreras-activas.md) |
| E05 | Cuando presiona "+ Agregar entidad" | User Interaction | [📄 Ver flujo](Flujos/flow-E05-cuando-presiona-agregar-entidad.md) |
| E06 | Cuando cambia ubicación o sector en el modal | User Interaction | [📄 Ver flujo](Flujos/flow-E06-cuando-cambia-ubicacion-en-modal.md) |
| E07 | Cuando guarda el formulario de nueva entidad | User Interaction | [📄 Ver flujo](Flujos/flow-E07-cuando-guarda-nueva-entidad.md) |
| E08 | Cuando cierra el modal de agregar entidad | User Interaction | [📄 Ver flujo](Flujos/flow-E08-cuando-cierra-modal-agregar.md) |
| E09 | Cuando presiona "Ver Detalle" de una entidad | User Interaction | [📄 Ver flujo](Flujos/flow-E09-cuando-presiona-ver-detalle.md) |

> ~~E10 — Cuando quita una entidad del caso~~ — **eliminado del alcance**: por decisión del usuario, esta versión no incluye funcionalidad para quitar una entidad ya agregada (ver `case-entities-interface.md`).

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-el-componente.md](Flujos/flow-E01-cuando-carga-el-componente.md)

**Evento:** Cuando carga el componente
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue (`mounted`), recibiendo `caseId` como prop. Hace `GET /api/v1/casos/:caseId/entidades`. Mientras carga, muestra el spinner. Si hay error, muestra el mensaje de error. Si la respuesta está vacía, muestra el empty state.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E02-cuando-escribe-en-el-buscador.md](Flujos/flow-E02-cuando-escribe-en-el-buscador.md)

**Evento:** Cuando escribe en el buscador
**Tipo:** User Interaction
**Descripción:** Filtra reactivamente (computed) las entidades ya cargadas cuyo `entityBranchName` contenga el texto (case-insensitive). No hace ninguna llamada al API — opera sobre los datos ya traídos en E01.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E03-cuando-selecciona-sector-filtro.md](Flujos/flow-E03-cuando-selecciona-sector-filtro.md)

**Evento:** Cuando selecciona un sector en el filtro
**Tipo:** User Interaction
**Descripción:** El usuario elige un sector del `<select>` (o "Todos"). Filtra reactivamente las entidades visibles por `sector === valor`. Client-side, sobre los datos ya cargados.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E04-cuando-selecciona-filtro-barreras-activas.md](Flujos/flow-E04-cuando-selecciona-filtro-barreras-activas.md)

**Evento:** Cuando selecciona el filtro "¿Con barreras activas?"
**Tipo:** User Interaction
**Descripción:** El usuario elige "Sí" / "No" / "Todos". Filtra reactivamente las entidades visibles por `barrerasActivasCount > 0` (o `=== 0` para "No"). Client-side, sobre los datos ya cargados — el conteo por entidad ya viene resuelto desde el backend en E01, gracias a `barrier_v2.entity_branch_id` (ver `related-tables.md`).
**Requerido:** Sí

---

📄 [Ver flujo → flow-E05-cuando-presiona-agregar-entidad.md](Flujos/flow-E05-cuando-presiona-agregar-entidad.md)

**Evento:** Cuando presiona "+ Agregar entidad"
**Tipo:** User Interaction
**Descripción:** Abre el modal de agregar entidad, reseteando el formulario (departamento/ciudad/municipio/sector/entidad/objetivo) y los errores de validación. El botón solo se renderiza (`v-if="puedeAgregar"`) para roles `sv`/`op`/`ro`; el backend además rechaza con 403 la petición de creación para cualquier otro rol, aunque se llame directamente a la API.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E06-cuando-cambia-ubicacion-en-modal.md](Flujos/flow-E06-cuando-cambia-ubicacion-en-modal.md)

**Evento:** Cuando cambia ubicación o sector en el modal
**Tipo:** User Interaction
**Descripción:** Al elegir/cambiar departamento, ciudad, municipio o sector dentro del modal, resetea el nivel de ubicación dependiente inferior. Departamento, ciudad, municipio y sector son todos obligatorios; solo cuando los cuatro tienen valor se llama `GET /api/v1/entity-branches?town_code=&sector=` para cargar las sedes disponibles.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E07-cuando-guarda-nueva-entidad.md](Flujos/flow-E07-cuando-guarda-nueva-entidad.md)

**Evento:** Cuando guarda el formulario de nueva entidad
**Tipo:** User Interaction
**Descripción:** Valida que haya una entidad y un objetivo. Si válido, `POST /api/v1/casos/:caseId/entidades`. Al completar, cierra el modal, recarga la lista (E01) y muestra confirmación.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E08-cuando-cierra-modal-agregar.md](Flujos/flow-E08-cuando-cierra-modal-agregar.md)

**Evento:** Cuando cierra el modal de agregar entidad
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Cancelar", la "×", o hace click en el overlay. Cierra el modal sin guardar cambios.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E09-cuando-presiona-ver-detalle.md](Flujos/flow-E09-cuando-presiona-ver-detalle.md)

**Evento:** Cuando presiona "Ver Detalle" de una entidad
**Tipo:** User Interaction
**Descripción:** Navega a `/salvia/entidad/:id` (ruta nueva, resuelta por decisión de producto — ver GAP más abajo).
**Requerido:** Sí

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos? — sí, E01
- [x] ¿Todos los campos interactivos tienen evento? — buscador (E02), sector (E03), barreras activas (E04) cubiertos
- [x] ¿Las validaciones de formulario están cubiertas? — E07 (entidad + objetivo requeridos)
- [x] ¿Los estados de error de carga están cubiertos? — E01
- [x] ¿El callback al padre está cubierto? — N/A: "Ver Detalle" ahora navega directamente (E09), ya no emite evento al padre
- [x] ¿Hay llamadas a API pendientes de crear? — No, ya existen: `GET/POST /api/v1/casos/:caseId/entidades` (`entity_case_controller.go`). `GET /api/v1/entity-branches` ya existía antes de este componente.
- [x] ¿Los permisos por rol están cubiertos? — sí: lectura para cualquier rol con acceso a `get_case_detail_sv`, escritura restringida a `sv`/`op`/`ro` (401/403 verificados con Playwright para `ps` y `ad`, ver `tests/api/case-entities.api.spec.ts` y `tests/e2e/case-entities.spec.ts`).

## GAPS de esta pantalla (no técnicos de BD, de producto/navegación)

| Decisión pendiente | Afecta |
|---|---|
| "Ver Detalle" navega a `/salvia/entidad/:id` — **ruta nueva, no existe hoy en el código real** (solo existe `/salvia/sedes`, la pantalla de listado de sedes). Falta definir/crear la pantalla de detalle de entidad; queda fuera del alcance de este componente pero es una dependencia directa de E09. | E09 |
| Mecanismo para poblar/editar `entity_case.last_action` — se definió que es un campo string editable, pero falta definir en qué paso del flujo se le da valor (¿lo escribe el agente manualmente? ¿se actualiza en algún otro evento?). | E01, futuro evento de edición |
| El campo `barrier_v2.entity_branch_id` nuevo solo se puede llenar si se modifica el formulario de registro de barrera (`hacer_seguimiento.html`) — cambio fuera de alcance de este componente. Mientras no se implemente, el conteo de "barreras activas" será 0 para todas las entidades. | E04 |
