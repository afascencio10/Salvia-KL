# `Ruta de atención integral` — Índice

> ✅ **Implementado el 2026-11-04 (WP-4 · REG-17)** y probado en navegador (ver [`QA/`](QA/ruta-atencion-integral-qa-cases.md)). Es el **8.º módulo** de [Registro V3](../Registro%20V3/registro-v3-index.md), entre «Tamizaje de riesgo» y «Plan de acción Salvia» (pestaña `RA`, sin candado). Sirve para **registrar a qué rutas o entidades acudió la persona y si encontró barreras en cada una**; es informativo: las activaciones de ruta se registran en el Plan de acción. Replica la sección «Identificación de entidades» de Seguimiento v4 con controles de Vue, **sin el componente de formulario dinámico**.

## Navegación

| MD | Contenido |
|---|---|
| [ruta-atencion-integral-interface.md](ruta-atencion-integral-interface.md) | Árbol de interfaz del módulo |
| [related-tables.md](related-tables.md) | Tablas que usa el módulo y cómo se llenan |
| [QA/ruta-atencion-integral-qa-cases.md](QA/ruta-atencion-integral-qa-cases.md) | Casos de prueba y resultados |

### Otros MDs

- [`../Registro V3/registro-v3-index.md`](../Registro%20V3/registro-v3-index.md): pantalla padre (módulos CI a GS y el Cierre).
- [`../Registro V3/Flujos/flow-subflujo-materializar-borrador.md`](../Registro%20V3/Flujos/flow-subflujo-materializar-borrador.md): sub-flujo que vacía el borrador JSON a las tablas al iniciar el Caso.
- [`../Registro V3/Flujos/flow-E17-cuando-agrega-o-elimina-una-entidad.md`](../Registro%20V3/Flujos/flow-E17-cuando-agrega-o-elimina-una-entidad.md) y [`flow-E18`](../Registro%20V3/Flujos/flow-E18-cuando-agrega-o-elimina-una-barrera.md): el formulario de entidad y de barrera se comparten con el Plan de acción.

### Changelogs

- [changelogNov2026.md](changelogNov2026.md)

---

## Resumen de eventos

Los eventos de este módulo se numeran aparte (RA-E01 a RA-E08), para no chocar con los E01–E23 de Registro V3.

| Evento | Tipo | Estado | Flujo |
|---|---|---|---|
| Cuando carga el módulo «Ruta de atención integral» | Lifecycle | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E01-cuando-carga-el-modulo.md) |
| Cuando responde «¿Ha puesto en conocimiento de alguna entidad…?» (`GS_14`) | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E02-cuando-responde-ya-se-reporto.md) |
| Cuando agrega una entidad acudida | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E03-cuando-agrega-una-entidad-acudida.md) |
| Cuando elimina una entidad acudida | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E04-cuando-elimina-una-entidad-acudida.md) |
| Cuando responde «¿Se identificaron barreras?» (`GS_17`) | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E05-cuando-responde-se-identificaron-barreras.md) |
| Cuando agrega una barrera a una entidad acudida | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E06-cuando-agrega-una-barrera.md) |
| Cuando elimina una barrera | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E07-cuando-elimina-una-barrera.md) |
| Cuando guarda la sección «Ruta de atención integral» | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-RA-E08-cuando-guarda-la-seccion.md) |

## API que usa el módulo

Todas bajo `/api/v1`, con sesión y permiso `manage_registro_v3`. **No hay endpoints propios:** se usan los de Registro V3.

| Endpoint | Evento | Qué hace |
|---|---|---|
| `GET /victim-cases-v3/:icode` | RA-E01 | Estado completo del Caso: respuestas, entidades, barreras, `entityMeta`, `barrierMeta` y `entitiesReadOnly` (lo carga `abrirCaso`) |
| `POST /victim-cases-v3/:icode/seccion/GS` | RA-E08 | Guarda las respuestas `GS_*`, entre ellas `GS_14` y `GS_17` |
| `POST /victim-cases-v3/:icode/entidades` | RA-E03 | Alta de entidad `acudida`: con el Caso en Borrador la guarda en el JSON (409 si la sede ya está); con el Caso ya Iniciado, 409 `readOnly` |
| `DELETE /victim-cases-v3/entidades/:entityIcode` | RA-E04 | Quita la entidad y sus barreras del borrador (404 si no existe; 409 `readOnly` si ya está materializada) |
| `POST /victim-cases-v3/:icode/barreras` | RA-E06 | Alta de barrera en el borrador (409 `readOnly` si el Caso ya salió de Borrador) |
| `DELETE /victim-cases-v3/barreras/:barrierIcode` | RA-E07 | Quita la barrera del borrador (409 `readOnly` si ya está materializada) |
| `GET /locations/…` · `GET /entity-branches?town_code&sector` · `GET /entity-branches/:id/obligations` | RA-E03 | Catálogos compartidos con Gestión institucional y Seguimiento |

---

## Inventario de eventos

**Evento:** Cuando carga el módulo «Ruta de atención integral»
**Tipo:** Lifecycle · **Estado:** ✅
**Descripción:** Al seleccionar la pestaña se muestra el bloque, el formulario de entidad toma el tipo `acudida` (`watch` de `moduloActivo`), se pintan las entidades acudidas ya guardadas con sus barreras y se decide si el formulario de entidad y el aviso de barreras están visibles.
**Requerido:** Sí

---

**Evento:** Cuando responde «¿Ha puesto en conocimiento de alguna entidad…?»
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** Guarda `GS_14` en el estado local y muestra u oculta las tarjetas y el formulario de entidad acudida, sin borrar nada: las entidades ya registradas se vuelven a ver al responder «Sí».
**Requerido:** Sí

---

**Evento:** Cuando agrega una entidad acudida
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** Busca la sede con la cascada obligatoria y «Buscar entidad» (como Gestión institucional), completa radicado, fecha de reporte, acciones completadas y pendientes y otra información (como Seguimiento v4) y, si tiene acciones pendientes, acuerdo, debida diligencia, canales y gestiones (D17); la guarda de tipo acudida en el borrador JSON y refresca la lista.
**Requerido:** Sí

---

**Evento:** Cuando elimina una entidad acudida
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** Pide confirmación, quita la entidad y todas sus barreras del borrador y refresca la lista. Después de materializar no hay botón (solo lectura).
**Requerido:** No

---

**Evento:** Cuando responde «¿Se identificaron barreras?»
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** Guarda `GS_17` en el estado local; con «Sí» habilita «+ Agregar barrera a esta entidad» en cada tarjeta y muestra el aviso de dónde registrarlas; si no hay entidades acudidas, también muestra el formulario de entidad aunque `GS_14` no sea «Sí» (y, si no hay ninguna entidad, agrega «Primero agrega la entidad.»).
**Requerido:** Sí

---

**Evento:** Cuando agrega una barrera a una entidad acudida
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** Abre el formulario de barrera de Seguimiento v4 (sin competencias), lo guarda en el borrador JSON y refresca la lista. Mientras REG-29 no se defina, el mismo botón aparece también en las tarjetas del Plan de acción (D8).
**Requerido:** No

---

**Evento:** Cuando elimina una barrera
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** Pide confirmación, quita la barrera del borrador y refresca la lista.
**Requerido:** No

---

**Evento:** Cuando guarda la sección «Ruta de atención integral»
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** «Guardar sección» propio: valida `GS_14` y `GS_17` y guarda el módulo GS (`enviarSeccion('GS')`); las entidades y barreras ya se guardaron al agregarlas. No avanza solo al siguiente módulo.
**Requerido:** Sí

---

## Checklist de validación de completitud

- [x] ¿Puede manejarse cada acción del usuario? Sí (RA-E01 a RA-E08).
- [x] ¿Están incluidos todos los envíos de formulario? Sí — RA-E03 y RA-E06 al agregar, RA-E08 por sección.
- [ ] ¿Hay actualizaciones en tiempo real? No aplica.
- [x] *(Backend)* ¿Hay tareas programadas? No; las tareas y oficios se crean al materializar el borrador (sub-flujo de Registro V3).
- [ ] *(Backend)* ¿Hay webhooks? No aplica.
