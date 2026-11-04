# QA — Ruta de atención integral (WP-4 · REG-17)

Casos de los eventos RA-E01 a RA-E08 ([índice](../ruta-atencion-integral-index.md)). Los de entidad y barrera comparten mecanismo con Registro V3: ver también [`../../Registro V3/QA/wp4-plan-de-accion-y-ruta-qa-cases.md`](../../Registro%20V3/QA/wp4-plan-de-accion-y-ruta-qa-cases.md).

> **Estado: ejecutado el 2026-11-04 — `ui-t8.js` 34/34 ✅ y `ui-solo-lectura.js` 21/21 ✅** (más los casos de API y BD de la batería, 318 en total; todo sobre el código rebasado en `develop`). Scripts en `DocsMD/Otros/temp/tests-wp4/`.

| Caso | Evento | Descripción | Resultado esperado | Estado |
|---|---|---|---|:---:|
| RA-1 | RA-E01 | Abrir un Caso en Borrador con la tríada incompleta y mirar las pestañas | «Ruta de atención integral» está entre Tamizaje y Plan de acción y **no tiene candado** (solo el Tamizaje); se abre | ✅ |
| RA-2 | RA-E01 | Sin responder `GS_14` | No hay tarjetas ni formulario de entidad | ✅ |
| RA-3 | RA-E02 | `GS_14` = Sí | Se ve solo la entidad **acudida** (la de activación no); aparece el formulario | ✅ |
| RA-4 | RA-E03 | Elegir sede sin radicado, fecha, acciones ni otra información | «+ Siguiente entidad» deshabilitado y «Falta: …» | ✅ |
| RA-5 | RA-E03 | Marcar una acción completada, sin pendientes | Se habilita; **no** aparece el bloque de acuerdo ni de canales | ✅ |
| RA-6 | RA-E03 | Marcar además acciones pendientes | Aparece «¿La persona está de acuerdo…?» y el botón se deshabilita | ✅ |
| RA-7 | RA-E03 | Acuerdo = Sí → debida diligencia → Sí | Aparecen los canales (5 opciones) y «¿Qué gestiones se realizaron…?»; sin ellos sigue deshabilitado | ✅ |
| RA-8 | RA-E03 | Agregar con canales Notificación + Mi Salvia | En el borrador: `kind = acudida`, radicado, fecha, objetivo, `agreement`/`dueDiligence` en `true`, 2 canales, 2 pendientes y 1 completada; **ninguna fila** en `entity_case` ni `case_task` | ✅ |
| RA-9 | RA-E03 (API) | Acudida con pendientes + acuerdo + debida diligencia, al iniciar el Caso | Se crean las tareas y oficios de cada canal; **sin pendientes** o sin acuerdo no se crea nada | ✅ |
| RA-10 | RA-E05 | `GS_17` = Sí | Aviso «Registra cada barrera con el botón…» y un «+ Agregar barrera a esta entidad» por tarjeta; el botón abre el formulario de Seguimiento v4 | ✅ |
| RA-11 | RA-E05 | `GS_17` = Sí, `GS_14` = No y sin entidades | Se muestra el formulario de entidad y «Primero agrega la entidad.» | ✅ |
| RA-12 | RA-E06/07 (API) | Agregar y eliminar barreras y entidades en el borrador | Eliminar la entidad quita sus barreras; borrar de nuevo → 404 | ✅ |
| RA-13 | RA-E08 | «Guardar sección» | «Sección guardada.»; `GS_14` y `GS_17` quedan como respondidas y en `victim_case_form2` | ✅ |
| RA-14 | RA-E08 | Completar la tríada desde la pantalla | El Caso pasa a Iniciado, el módulo se recarga y queda en **solo lectura** | ✅ |
| RA-15 | Solo lectura | Caso ya Iniciado | Sin formulario, sin «Eliminar» y sin «+ Agregar barrera»; `GS_14`/`GS_17` siguen respondibles; las altas por API responden 409 `readOnly` | ✅ |
| RA-16 | Solo lectura | La pantalla seguía abierta cuando otro actor completó la tríada y se intenta agregar | El 409 pone la pantalla en solo lectura y no se crea la entidad | ✅ |
