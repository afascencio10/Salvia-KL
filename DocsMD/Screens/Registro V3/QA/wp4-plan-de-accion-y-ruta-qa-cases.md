# QA — Registro V3 · Plan de acción, Ruta de atención integral, entidades y barreras (WP-4 · REG-14/15/16/17/18/19/20/38*, REG-04 c)

Casos de [`flow-E10`](../Flujos/flow-E10-cuando-cambia-de-modulo.md), [`flow-E11`](../Flujos/flow-E11-cuando-actualiza-una-respuesta.md), [`flow-E12`](../Flujos/flow-E12-cuando-guarda-una-seccion.md), [`flow-E15`](../Flujos/flow-E15-cuando-finaliza-y-enruta-a-seguimiento.md), [`flow-E17`](../Flujos/flow-E17-cuando-agrega-o-elimina-una-entidad.md), [`flow-E18`](../Flujos/flow-E18-cuando-agrega-o-elimina-una-barrera.md) y el [sub-flujo de materialización](../Flujos/flow-subflujo-materializar-borrador.md). Casos del módulo nuevo: [`../../Ruta de atención integral/QA/`](../../Ruta%20de%20atenci%C3%B3n%20integral/QA/ruta-atencion-integral-qa-cases.md).

> **Estado: ejecutado el 2026-11-04 sobre el código ya rebasado en `develop` — 318 de 318 ✅** (API + BD: 142 · pantalla con Playwright: 176).

**Método.** API primero (curl/`fetch` con la sesión del navegador + lectura de la BD local) y Playwright solo para lo que la API no cubre. Cada caso crea su propio Caso (documento aleatorio); la tríada se marca como respondida por SQL (`victim_case_v3_answered_keys`) o se completa por la pantalla. Scripts: `DocsMD/Otros/temp/tests-wp4/` — `bash run-all.sh` corre toda la batería contra `./run_dev.sh`.

## API + BD

| Script | Casos clave | Resultado |
|---|---|---|
| `t9-api-test.sh` (63) | Entidad acudida + activación con 4 canales; duplicado → 409; barreras en la sede; **el borrador no toca `entity_case`, `barrier_v2` ni `case_task`**; `GetCase` devuelve los mismos campos y `entityMeta`/`barrierMeta`; eliminar la entidad borra sus barreras solo en el JSON; borrar dos veces → 404; **materializar al completar la tríada**: 1 entidad con sus 3 obligaciones (1 completada, 2 pendientes, sin «otra»), 1 barrera `OPEN`, **6 tareas** y **2 oficios** con los vínculos correctos, JSON vaciado con la meta conservada; idempotencia; **409 de solo lectura** en las 4 operaciones; `entitiesReadOnly = true` | ✅ 63/63 |
| `t9b-api-test.sh` (57) | Remisiones: psicosocial (criterio obligatorio + 3 puntos, `referral_criteria` en el vocabulario de Seguimiento, tarea «Validar remisión»), hombres, estabilización, discapacidad ×2, medidas de emergencia (`alojamiento`, `transporte`); psicosocial omitida con medidas; criterios insuficientes; `GS_5` = No; acudida **sin pendientes** o **sin acuerdo** no crea tareas; **fallo a mitad** (barrera con sector de 80 caracteres): 500, el Caso sigue en `bo`, nada a medias, el JSON intacto y, corregido, pasa; sede ya existente (se reutiliza y se le agregan sus acciones); clientes anteriores (`no,ll,ms`, gestión 01–06, `Competences`); «Ninguna» no crea tareas; **reintento al leer** | ✅ 57/57 |
| `t5-api-test.sh` (13) | `GET /enums`: 15 opciones de `GS_20` (01…14, 99 al final); guardar 3 y recargar las mismas 3; «Ninguno» deja solo el 99; vaciar; no es obligatoria | ✅ 13/13 |
| `t1-finalize-test.sh` (9) | Finalizar sin `GS_2` → 200 y `te`; sin `GS_4`, `GS_5`, `GS_14`, `GS_16` o `GS_17` → 400 con `modulosIncompletos: ["GS"]`; sin `GS_20` → 200 | ✅ 9/9 |

## Pantalla (Playwright)

| Script | Casos clave | Resultado |
|---|---|---|
| `ui-t2.js` (25) | REG-18: solo se ve «Sector…» y las obligatorias; Justicia muestra sus 3 preguntas (13 opciones) y no las de otros sectores; Transversales muestra las 4 listas y al volver a Justicia se limpió lo anterior; Municipio deshabilitado sin Departamento; la perspectiva y el aviso aparecen solo con «Alerta»; el botón se habilita con todo lo obligatorio; lo que llega al backend (`barrierMeta`: zona, perspectiva, descripción de la gestión) | ✅ |
| `ui-t3.js` (20) | REG-20: datos anteriores («no,ll» y «Canal de activación: ll») se ven con etiquetas de Seguimiento; sin `GS_15` y `GS_16` en «Sí» no se piden canales; 5 opciones; «Ninguna» excluyente; agregar con canales y gestiones; `entityMeta` con `agreement`/`dueDiligence` copiados; `solicitud` ya no lleva la nota | ✅ |
| `ui-t4.js` (14) | REG-04 c: 2 gestiones nuevas, 6 históricas (01–06), valor desconocido tal cual y «—»; ninguna celda deja un slug | ✅ |
| `ui-t5.js` (10) | REG-16: texto exacto de la pregunta y de las 15 opciones; «Ninguno» excluyente en ambos sentidos; guardar la sección y reabrir el Caso conserva la selección | ✅ |
| `ui-t6.js` (23) | REG-19: asteriscos y cascada deshabilitada; «Buscar entidad» no aparece sin municipio y sector; **una sola petición de sedes**; el filtro no se aplica antes de 250 ms; «bolivar» encuentra «…BOLÍVAR»; escribir no envía ninguna petición; sin coincidencias se avisa; elegir carga las acciones (1 petición); chips marcados; editar el texto limpia la selección; **una respuesta atrasada no pisa a la nueva** | ✅ |
| `ui-t7.js` (24) | REG-15: pestaña «Plan de acción Salvia»; orden vertical b, c, d, e, h, i y § E; sin títulos A–D, sin `GS_14`/`GS_17`/`GS_2`; `GS_15` siempre visible; `GS_16` despliega el formulario sin selector de tipo; Objetivo obligatorio; el JSON guarda solicitud, radicado, objetivo y las 2 acciones; la tarjeta las lista como «Pendiente: …» | ✅ |
| `ui-t8.js` (34) | REG-17: pestaña entre Tamizaje y Plan de acción sin candado con la tríada incompleta; acudida solo en la Ruta y activación solo en el Plan; preguntas y obligatoriedad de Seguimiento; **con acciones pendientes** aparecen acuerdo → debida diligencia → canales y gestiones (D17); borrador con `kind = acudida`; `GS_17` = Sí muestra el aviso y el botón de barrera en cada tarjeta; «Guardar sección» guarda `GS_14` y `GS_17`; con `GS_17` = Sí y sin entidades se muestra el formulario | ✅ |
| `ui-solo-lectura.js` (21) | D20: tras materializar no hay formulario, «Eliminar» ni «+ Agregar barrera»; la remisión queda deshabilitada (`GS_5` solo si ya estaba respondida); `GS_15`/`GS_16`/`GS_14`/`GS_17` siguen respondibles; las tarjetas muestran canales y acciones del snapshot; **una pantalla atrasada** (otro actor completó la tríada) pasa a solo lectura por el 409; al completar la tríada desde la pantalla se recarga y queda en solo lectura | ✅ |
| `ui-cierre.js` (5) | El Cierre muestra «Canales: Notificación SALVIA, Denuncia fácil» en la fila de la entidad y la descripción de la barrera (tras el reintento al leer) | ✅ |

## Observaciones

- Durante el rebase sobre `develop` (WP-2 y WP-3) solo hubo que adaptar las pruebas de pantalla al menú vertical y al bloqueo exclusivo del Tamizaje; el código de WP-4 pasó sin cambios.
- Seguimiento v4 muestra «Ninguna» como quinta opción de canal (la planeación decía 4): se muestra igual y no crea tarea.
