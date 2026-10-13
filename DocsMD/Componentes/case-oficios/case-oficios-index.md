# case-oficios — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [case-oficios-interface.md](case-oficios-interface.md) |
| Uso | [case-oficios-usage.md](case-oficios-usage.md) |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga el componente | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-el-componente.md) |
| E02 | Cuando escribe en el buscador | User Interaction | — |
| E03 | Cuando selecciona un chip de tema | User Interaction | — |
| E04 | Cuando selecciona un estado en el filtro | User Interaction | — |
| E05 | Cuando presiona una card de oficio | User Interaction | — |
| E06 | Cuando cierra el modal de detalle | User Interaction | — |
| E07 | Cuando descarga el PDF del oficio | User Interaction | — |
| E08 | Cuando imprime el oficio | User Interaction | — |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-el-componente.md](Flujos/flow-E01-cuando-carga-el-componente.md)

**Evento:** Cuando carga el componente
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue (`mounted`), recibiendo `caseId` y/o `barrierId` como props. Si `barrierId` está presente hace `GET /api/v1/entity-letters?barrierId={barrierId}` (usado en Detalle de Barrera); si no, `GET /api/v1/entity-letters?caseId={caseId}` (usado en Detalle del Caso). Guarda el resultado. Si falla, muestra un mensaje de error con botón de reintentar; si no hay oficios, el estado vacío se muestra después de aplicar filtros (ver E02-E04).
**Requerido:** Sí

---

**Evento:** Cuando escribe en el buscador
**Tipo:** User Interaction
**Descripción:** El usuario escribe texto en el buscador. Filtra reactivamente (computed) los oficios ya cargados cuyo `entidad`, `url_kofax` o campo de asunto disponible contengan el texto (case-insensitive). No hace ninguna llamada al API — opera sobre los datos ya traídos en E01.
**Requerido:** Sí

---

**Evento:** Cuando selecciona un chip de tema
**Tipo:** User Interaction
**Descripción:** El usuario hace click en un chip de tema (ej. "Barreras"). Activa/desactiva ese filtro y recalcula la lista visible. Solo el chip "Barreras" es funcional hoy — los otros 3 dependen del GAP de schema documentado en `case-oficios-interface.md`.
**Requerido:** Sí

---

**Evento:** Cuando selecciona un estado en el filtro
**Tipo:** User Interaction
**Descripción:** El usuario elige un estado del `<select>` (o "Todos"). Filtra reactivamente los oficios visibles por `state === valor`. Client-side — no llama al API con `?state=`, porque combinarlo con `?caseId=` no es posible en el endpoint real (son mutuamente excluyentes en el backend).
**Requerido:** Sí

---

**Evento:** Cuando presiona una card de oficio
**Tipo:** User Interaction
**Descripción:** Setea `oficioSeleccionado = oficio` (el objeto ya está en memoria, sin fetch adicional) y el modal se muestra reactivamente con todo su detalle.
**Requerido:** Sí

---

**Evento:** Cuando cierra el modal de detalle
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Cerrar", el botón "×" o hace click en el overlay oscuro. Setea `oficioSeleccionado = null`. No emite ningún evento al padre.
**Requerido:** Sí

---

**Evento:** Cuando descarga el PDF del oficio
**Tipo:** User Interaction
**Descripción:** Abre `oficio.url_kofax` en una pestaña nueva (`window.open`), si existe. No hay lógica de servidor involucrada — es un link directo al archivo.
**Requerido:** No (solo si `url_kofax` existe)

---

**Evento:** Cuando imprime el oficio
**Tipo:** User Interaction
**Descripción:** Llama a `window.print()` sobre la vista actual del modal.
**Requerido:** No

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos? — sí, E01
- [x] ¿Todos los campos interactivos del formulario tienen evento? — buscador (E02), chips (E03), select estado (E04) cubiertos
- [x] ¿Las validaciones de formulario están cubiertas? — N/A, no hay formulario de entrada, solo filtros de lectura
- [x] ¿Los estados de error de carga y guardado están cubiertos? — error de carga sí (E01); no hay operación de guardado, es de solo lectura
- [x] ¿El callback al padre está cubierto? — N/A, el componente no emite eventos
