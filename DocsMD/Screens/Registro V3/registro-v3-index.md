# `Registro V3` — Índice

> ✅ **Implementado el 2026-10-14/15 y sincronizado con el código el 2026-10-31.** Ruta `GET /salvia/registro-v3`, roles `op` y `sv` (permiso `manage_registro_v3`), sin entrada en el menú lateral. Cada evento lleva su estado: ✅ implementado · 🟡 implementado con cambios respecto a la planeación · ⛔ no implementado. Resumen del cambio en [`changelogOct2026.md`](changelogOct2026.md).

## Navegación

| MD | Contenido |
|---|---|
| [registro-v3-interface.md](registro-v3-interface.md) | Árbol de interfaz real y tabla "planeado y no implementado" |
| [registro-v3-data.md](registro-v3-data.md) | Decisiones de arquitectura, mapeo a base de datos real y estado de los GAPS |

### Otros MDs

| MD | Contenido |
|---|---|
| [registro-v3-secciones.md](registro-v3-secciones.md) | Catálogo de preguntas por módulo (CI/DP/CA/HV/AG/TZ/GS) y divergencias con el catálogo del código |
| [registro-v3-parrafo-sintesis.md](registro-v3-parrafo-sintesis.md) | Reglas del párrafo síntesis (módulo GS y pantalla de Cierre) |
| [registro-v3-pendientes.md](registro-v3-pendientes.md) | Reconciliación planeación↔código del 2026-10-15: qué se corrigió y cómo se probó |
| [QA/registro-v3-qa-plan.md](QA/registro-v3-qa-plan.md) | Plan de QA con Playwright, hallazgos y bugs corregidos |
| [QA/](QA/) | Casos de prueba por tema: [cabezote](QA/E04-E06-cabezote-qa-cases.md) · [condicionales](QA/condicionales-formulario-qa-cases.md) · [entidades y barreras](QA/entidades-barreras-qa-cases.md) · [finalizar y alerta](QA/finalizar-y-alerta-qa-cases.md) · [gestión Salvia y remisión](QA/gestion-salvia-remision-qa-cases.md) · [tríada y guardado](QA/triada-y-guardado-qa-cases.md) · [tamizaje](QA/tamizaje-qa-cases.md) · [E2E](QA/registro-v3-e2e-qa-cases.md) · [implementación completa](QA/implementacion-completa-2026-10-15-qa-cases.md) |

### Changelogs

- [changelogOct2026.md](changelogOct2026.md)

---

## Resumen de eventos

| Evento | Tipo | Estado | Flujo |
|---|---|---|---|
| Cuando carga la pantalla de Registro | Lifecycle | 🟡 | 📄 [Ver flujo](Flujos/flow-E01-cuando-carga-la-pantalla.md) |
| Cuando responde la autorización de datos personales | User Interaction | ✅ | 📄 [Ver flujo](Flujos/flow-E02-cuando-responde-habeas-data.md) |
| Cuando selecciona el tipo de contacto | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E03-cuando-selecciona-tipo-contacto.md) |
| Cuando busca en el Cabezote | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E04-cuando-busca-en-cabezote.md) |
| Cuando presiona "Limpiar" en el Cabezote | User Interaction | ⛔ | 📄 [Ver flujo](Flujos/flow-E05-cuando-limpia-busqueda.md) |
| Cuando presiona "Nuevo Registro" | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E06-cuando-presiona-nuevo-registro.md) |
| Cuando presiona "Proceder al formulario" (hoy: "Abrir") | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E07-cuando-procede-al-formulario.md) |
| Cuando presiona "Nuevo Caso en Registro existente" | User Interaction | ⛔ | 📄 [Ver flujo](Flujos/flow-E08-cuando-presiona-nuevo-caso-en-registro.md) |
| Cuando presiona "Actualizar Caso" | User Interaction | ⛔ | 📄 [Ver flujo](Flujos/flow-E09-cuando-presiona-actualizar-caso.md) |
| Cuando cambia de módulo | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E10-cuando-cambia-de-modulo.md) |
| Cuando actualiza una respuesta del formulario | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E11-cuando-actualiza-una-respuesta.md) |
| Cuando guarda una sección/módulo | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E12-cuando-guarda-una-seccion.md) |
| Cuando presiona "Alerta inminente" | User Interaction | 🟡 (campos GS_18/GS_19, sin botón) | 📄 [Ver flujo](Flujos/flow-E13-cuando-presiona-alerta-inminente.md) |
| Cuando abandona la pantalla sin finalizar | User Interaction + Lifecycle | ✅ | 📄 [Ver flujo](Flujos/flow-E14-cuando-abandona-sin-finalizar.md) |
| Cuando presiona "Finalizar registro y enrutar a Seguimiento" | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E15-cuando-finaliza-y-enruta-a-seguimiento.md) |
| Cuando presiona un acceso de Rutas/Barreras/Notificaciones | User Interaction | ⛔ (fuera de alcance) | 📄 [Ver flujo](Flujos/flow-E16-cuando-presiona-rutas-barreras-notificaciones.md) |
| Cuando agrega o elimina una Entidad | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E17-cuando-agrega-o-elimina-una-entidad.md) |
| Cuando agrega o elimina una Barrera | User Interaction | 🟡 | 📄 [Ver flujo](Flujos/flow-E18-cuando-agrega-o-elimina-una-barrera.md) |

## API que usa la pantalla

Todas bajo `/api/v1`, con sesión y permiso `manage_registro_v3` (sin sesión 401; sin permiso, 401 del middleware). Los errores de servicio salen como 500 genérico, salvo los indicados.

| Endpoint | Evento | Qué hace |
|---|---|---|
| `GET /victim-cases-v3/enums` | E-01 | Catálogo de opciones `{code, icode, label}` por categoría de `victim_case_form2_enums` (el frontend resuelve todos los dropdowns y checkboxes contra él) |
| `GET /victim-cases-v3/buscar?tipo_documento&documento&telefono` | E-04, E-11 | Casos que coinciden → `{registros: [{icode, names, lastNames, docType, docNumber, phone, status, registroCompleto}]}`. 400 si no hay documento ni teléfono. Si viene documento busca **solo por documento** (más el tipo si se envía) e ignora el teléfono; sin documento, busca por teléfono |
| `POST /victim-cases-v3/draft` `{tipo_documento, documento, telefono, municipio_residencia, canal}` | E-06 | Crea el Borrador (`general_user` + `victim_case` `bo` + `victim_case_form2` con valores dummy) → `{icode}`. El frontend solo manda tipo, documento y teléfono |
| `GET /victim-cases-v3/:icode` | E-07, E-17, E-18 | Estado completo: `summary`, `answers`, `answeredKeys`, `triadaCompleta`, `entities`, `barriers`, `tzResultado` |
| `POST /victim-cases-v3/:icode/seccion/:modulo` `{answers}` | E-12, E-14 | Guarda CI/DP/CA/HV/AG/GS → `{ok, triadaCompleta}`; pasa `bo`→`in` cuando la tríada queda completa |
| `POST /victim-cases-v3/:icode/tamizaje` `{answers, wasPartner, esOsigd, subtipoFeminicidio, subtipoHomicidioAbyecto}` | E-11 | Calcula tipo, nivel y puntaje → `{riskType, riskLevel, score, tieFlag}` |
| `POST /victim-cases-v3/:icode/finalizar` | E-15 | Valida obligatorias y pasa `in`→`te` → `{ok, denunciaFacilSolicitada, equiposRemisionCSV}`; 400 `{error, modulosIncompletos}` |
| `POST /victim-cases-v3/:icode/entidades` · `DELETE /victim-cases-v3/entidades/:entityIcode` | E-17 | Alta (409 si la sede ya está en el Caso) y baja lógica en `entity_case` |
| `POST /victim-cases-v3/:icode/barreras` · `DELETE /victim-cases-v3/barreras/:barrierIcode` | E-18 | Alta y baja en `barrier_v2` |
| `GET /locations/departments`, `/locations/cities?department_id`, `/locations/towns?city_id` · `GET /entity-branches?town_code&sector` · `GET /entity-branches/:id/obligations` · `GET /competences?sector` | E-17, E-18 | Catálogos compartidos con Gestión institucional y Seguimiento |

---

## Inventario de eventos

**Evento:** Cuando carga la pantalla de Registro
**Tipo:** Lifecycle · **Estado:** 🟡
**Descripción:** Renderiza la página en `etapa = 'habeas'`. `mounted()` carga el catálogo de enums y los departamentos y registra el respaldo `beforeunload` (E-14). **No** lee ningún `icode` de la URL: no hay entrada directa a un Caso.
**Requerido:** Sí

---

**Evento:** Cuando responde la autorización de datos personales
**Tipo:** User Interaction · **Estado:** ✅
**Descripción:** "Acepto y continuar" lleva a `etapa = 'tipoContacto'`. "No acepto" lleva a `formularioCerrado` y no crea ningún registro.
**Requerido:** Sí

---

**Evento:** Cuando selecciona el tipo de contacto
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** Cinco botones. "Víctima - Caso Nuevo" pasa al Cabezote; las otras cuatro (Seguimiento, Tercero, Familiares, Hombres) solo muestran un aviso de que se gestionan en otra pantalla.
**Requerido:** Sí

---

**Evento:** Cuando busca en el Cabezote
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** Busca por documento y/o teléfono (nunca por nombre) en `GET /victim-cases-v3/buscar`. Muestra una tabla única con el estado crudo del Caso y un badge Completo/Incompleto del Registro. No clasifica en reportes/registros/borradores ni consulta reportes pendientes.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Limpiar" en el Cabezote
**Tipo:** User Interaction · **Estado:** ⛔
**Descripción:** Planeado para resetear criterios y resultados. El botón no existe.
**Requerido:** No

---

**Evento:** Cuando presiona "Nuevo Registro"
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** "+ Nuevo registro" (siempre habilitado) crea un Borrador `bo` con el documento y teléfono buscados y abre el Formulario en CI. No exige confirmar que se trata de otra persona.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Proceder al formulario"
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** Hoy es el botón "Abrir" de cada fila del Cabezote (`abrirCaso`): carga el Caso con sus respuestas y abre el Formulario en CI. No registra trazabilidad de quién lo retoma.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Nuevo Caso en Registro existente"
**Tipo:** User Interaction · **Estado:** ⛔
**Descripción:** Diferido junto con la separación Registro/Caso (decisión D1). No existe el botón ni el endpoint.
**Requerido:** No

---

**Evento:** Cuando presiona "Actualizar Caso"
**Tipo:** User Interaction · **Estado:** ⛔
**Descripción:** No existe el botón ni el endpoint de "nuevo episodio".
**Requerido:** No

---

**Evento:** Cuando cambia de módulo
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** Las pestañas CI…GS cambian `moduloActivo`. TZ y GS están bloqueadas (🔒) mientras `!triadaCompleta` y muestran un aviso en vez de navegar. No hay barras de progreso.
**Requerido:** Sí

---

**Evento:** Cuando actualiza una respuesta del formulario
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** Guarda el valor en `answers` y, en la misma ejecución, re-evalúa los condicionales del catálogo, sugiere la nacionalidad según el tipo de documento (nunca pisa una respuesta manual), valida "solo letras" y la máscara de CI_6, avisa de duplicidad de documento (debounce 500 ms), recalcula el tratamiento (Sra./Sr./Sre.) y el criterio de remisión de GS, y en TZ revela el resto de preguntas si alguna inicial es "Sí". El Tamizaje se calcula al presionar "Calcular y guardar Tamizaje", no en cada cambio.
**Requerido:** Sí

---

**Evento:** Cuando guarda una sección/módulo
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** "Guardar sección" valida las obligatorias [1er contacto] visibles, "solo letras" y la máscara de CI_6; envía solo las respuestas del módulo a `POST /victim-cases-v3/:icode/seccion/:modulo`. Si la tríada queda completa, el Caso pasa de Borrador a Iniciado.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Alerta inminente"
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** No hay botón ni bitácora: la alerta es la pregunta GS_18 ("¿Activar alerta inminente?") con su justificación GS_19, guardada con la sección GS. Se pierde la notificación a coordinación; solo queda el dato.
**Requerido:** No

---

**Evento:** Cuando abandona la pantalla sin finalizar
**Tipo:** User Interaction + Lifecycle · **Estado:** ✅
**Descripción:** "Finalizar llamada (guardar y salir)" guarda el módulo activo sin exigir obligatorias y vuelve al Cabezote. Como respaldo, `beforeunload` envía el módulo activo con `navigator.sendBeacon`.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Finalizar registro y enrutar a Seguimiento"
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** El backend valida las obligatorias incondicionales de los 7 módulos y pasa el Caso a Terminado (`te`). No genera Denuncia fácil ni remisiones: devuelve si se solicitó la denuncia (GS_1) y los equipos de GS_6 para mostrarlos como confirmación en el Cierre. No hay un traspaso aparte: Seguimiento lee las mismas tablas.
**Requerido:** Sí

---

**Evento:** Cuando presiona un acceso de Rutas/Barreras/Notificaciones
**Tipo:** User Interaction · **Estado:** ⛔
**Descripción:** Fuera de alcance de esta ronda; no existe el grupo "Registro de gestión".
**Requerido:** No

---

**Evento:** Cuando agrega o elimina una Entidad
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** En GS § C se elige sede por Departamento → Ciudad → Municipio + Sector y se agrega a `entity_case` con sus obligaciones (completadas/pendientes). Eliminar pide confirmación (el texto promete borrar también sus barreras, pero el backend solo da de baja la entidad: G-V3-22).
**Requerido:** No

---

**Evento:** Cuando agrega o elimina una Barrera
**Tipo:** User Interaction · **Estado:** 🟡
**Descripción:** Las barreras cuelgan de una entidad ya agregada (sub-repeater) y se guardan en `barrier_v2` heredando sede, sector y ubicación. Solo se ofrece "+ Agregar barrera" si GS_17 = Sí.
**Requerido:** No

---

## Checklist de validación de completitud

- [x] ¿Puede manejarse cada acción del usuario? Sí, salvo los eventos marcados ⛔.
- [x] ¿Está cubierta la carga inicial de datos? Sí — E01.
- [x] ¿Están incluidos todos los envíos de formulario? Sí — E12 (por sección), E11 (Tamizaje) y E15 (final).
- [ ] ¿Hay actualizaciones en tiempo real (sockets, notificaciones)? No aplica.
- [ ] *(Backend)* ¿Hay tareas programadas? No aplica.
- [ ] *(Backend)* ¿Hay webhooks? No aplica.
