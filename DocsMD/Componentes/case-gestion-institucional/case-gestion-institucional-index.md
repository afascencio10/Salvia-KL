# case-gestion-institucional — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [case-gestion-institucional-interface.md](case-gestion-institucional-interface.md) |
| Uso | [case-gestion-institucional-usage.md](case-gestion-institucional-usage.md) |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga el componente | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-el-componente.md) |
| E02 | Cuando presiona el tile "Entidades agregadas" | User Interaction | [📄 Ver flujo](Flujos/flow-E02-cuando-presiona-tile-entidades.md) |
| E03 | Cuando presiona el tile "Obligaciones completadas" | User Interaction | [📄 Ver flujo](Flujos/flow-E03-cuando-presiona-tile-obligaciones.md) |
| E04 | Cuando presiona un anillo de sector | User Interaction | [📄 Ver flujo](Flujos/flow-E04-cuando-presiona-anillo-sector.md) |
| E05 | Cuando presiona el tile "Barreras sin resolver" | User Interaction | [📄 Ver flujo](Flujos/flow-E05-cuando-presiona-tile-barreras.md) |
| E06 | Cuando presiona un badge de "Barreras por estado" | User Interaction | [📄 Ver flujo](Flujos/flow-E06-cuando-presiona-badge-barrera-estado.md) |
| E07 | Cuando presiona el tile "Oficios en gestión" | User Interaction | [📄 Ver flujo](Flujos/flow-E07-cuando-presiona-tile-oficios.md) |
| E08 | Cuando presiona un badge de "Oficios por estado" | User Interaction | [📄 Ver flujo](Flujos/flow-E08-cuando-presiona-badge-oficio-estado.md) |
| E09 | Cuando cierra el modal de detalle | User Interaction | [📄 Ver flujo](Flujos/flow-E09-cuando-cierra-modal.md) |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-el-componente.md](Flujos/flow-E01-cuando-carga-el-componente.md)

**Evento:** Cuando carga el componente
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente (`mounted`), recibiendo `caseId` y `barriers` como props. Dispara en paralelo `GET /api/v1/casos/:caseId/entidades`, `GET /api/v1/casos/:caseId/obligaciones-resumen` (nuevo, agregado en una sola consulta) y `GET /api/v1/entity-letters?caseId=`. `barriers` no se fetchea — ya llega resuelto desde el padre. Mientras carga muestra el spinner; si algo falla, el mensaje de error con reintentar.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E02-cuando-presiona-tile-entidades.md](Flujos/flow-E02-cuando-presiona-tile-entidades.md)

**Evento:** Cuando presiona el tile "Entidades agregadas"
**Tipo:** User Interaction
**Descripción:** Abre el modal de detalle con `tipo: 'entidades'` y todas las entidades del caso como filas (sector, nombre, ubicación, oficios, barreras activas — mismos campos que ya resuelve `GET /api/v1/casos/:caseId/entidades`).
**Requerido:** Sí

---

📄 [Ver flujo → flow-E03-cuando-presiona-tile-obligaciones.md](Flujos/flow-E03-cuando-presiona-tile-obligaciones.md)

**Evento:** Cuando presiona el tile "Obligaciones completadas"
**Tipo:** User Interaction
**Descripción:** Abre el modal con `tipo: 'obligaciones'` y todas las filas de `resumenObligaciones.items` (ya traídas en E01), cada una con su estado (completada/pendiente).
**Requerido:** Sí

---

📄 [Ver flujo → flow-E04-cuando-presiona-anillo-sector.md](Flujos/flow-E04-cuando-presiona-anillo-sector.md)

**Evento:** Cuando presiona un anillo de sector (Salud, Justicia o Protección)
**Tipo:** User Interaction
**Descripción:** Abre el modal de obligaciones, igual que el evento anterior, pero filtrado a `item.sector === sector` únicamente.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E05-cuando-presiona-tile-barreras.md](Flujos/flow-E05-cuando-presiona-tile-barreras.md)

**Evento:** Cuando presiona el tile "Barreras sin resolver"
**Tipo:** User Interaction
**Descripción:** Abre el modal con `tipo: 'barreras'` y las barreras de la prop `barriers` cuyo estado no sea "Resuelta" (es decir, `OPEN` o `En Gestion`).
**Requerido:** Sí

---

📄 [Ver flujo → flow-E06-cuando-presiona-badge-barrera-estado.md](Flujos/flow-E06-cuando-presiona-badge-barrera-estado.md)

**Evento:** Cuando presiona un badge de "Barreras por estado" (Abierta / En gestión / Resuelta)
**Tipo:** User Interaction
**Descripción:** Abre el modal de barreras filtrado exactamente a ese estado (client-side, sobre la prop `barriers` ya cargada).
**Requerido:** Sí

---

📄 [Ver flujo → flow-E07-cuando-presiona-tile-oficios.md](Flujos/flow-E07-cuando-presiona-tile-oficios.md)

**Evento:** Cuando presiona el tile "Oficios en gestión"
**Tipo:** User Interaction
**Descripción:** Abre el modal con `tipo: 'oficios'` y los oficios del caso cuyo `state` no sea `radicado` ni `respondido`.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E08-cuando-presiona-badge-oficio-estado.md](Flujos/flow-E08-cuando-presiona-badge-oficio-estado.md)

**Evento:** Cuando presiona un badge de "Oficios por estado" (En gestión / Radicado / Respondido)
**Tipo:** User Interaction
**Descripción:** Abre el modal de oficios filtrado a ese bucket exacto (client-side, sobre los oficios ya cargados en E01).
**Requerido:** Sí

---

📄 [Ver flujo → flow-E09-cuando-cierra-modal.md](Flujos/flow-E09-cuando-cierra-modal.md)

**Evento:** Cuando cierra el modal de detalle
**Tipo:** User Interaction
**Descripción:** El usuario presiona "×", "Cerrar", hace click en el overlay, o presiona Escape. Limpia el estado `modal = null`.
**Requerido:** Sí

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos? → E01
- [x] ¿Todo elemento clickeable del árbol de interfaz tiene evento? → E02–E08
- [x] ¿El cierre del modal está cubierto? → E09
- [x] ¿Hay estados de error/vacío cubiertos? → E01 (error), cada modal (vacío, ver interface.md)
- [x] Flujos lógicos de los 9 eventos documentados
