# QA — Tarea de validación de remisión psicosocial (E-04)

Casos de prueba para el paso nuevo agregado a `processFollowUpSubmission()` ([flow-E04](../Flujos/flow-E04-cuando-se-procesa-submission.md), PASO 4, caso `"atencion_psico"`): cuando se crea una remisión a Atención Psicosocial (`salvia.psychosocial_support`), también se crea una `salvia.case_task` (`category: "Psicosocial"`, `type: "validar_remision"`) sin asignar, para que el equipo psicosocial valide si la remisión es procedente.

> **Alcance:** solo se prueban los caminos de decisión que afectan directamente esta lógica nueva. No se re-prueban las demás ramas de E-04 (barreras, medidas de emergencia, cierre de caso, reasignación) — ya existían antes de este cambio y no fueron modificadas.

---

## Método de prueba

Backend puro, sin Playwright — contra el caso de prueba E2E ya usado por `qa-salvia` (`CASE_ID = 019eadd1-6df3-7728-b014-6db076cbbb1e`), reutilizando el mismo `follow_up_v2` / `form_submission` ya completo (`follow_up_id = fd053f17-5614-4206-8c9f-1a61012eacc8`, `form_submission_id = 317e7a2d-8827-40ad-94f9-660856e9d49d`):

1. Marcar el `follow_up_v2` de prueba como `PENDIENTE` en BD — `processFollowUpSubmission` trata `REALIZADO` como idempotente y no reprocesa.
2. Modificar directamente en BD (`salvia.answer`) las respuestas del `form_submission_id` para `qEquipos` (`e0d38cf5-…`), `qCriteriosPsico` (`71c42c4a-…`) y `qMedidasEmergencia` (`1a36260c-…`).
3. `POST http://localhost:9090/api/v1/forms/saveSection` guardando solo la Sección 6 (Cierre del caso, con `¿Realiza cierre? = false`) para disparar `allAnswered == true` sin cerrar el caso de prueba.
4. Verificar logs en terminal y los registros creados directamente en BD (`salvia.psychosocial_support`, `salvia.case_task`).

> Las requests no llevan cookie de sesión — `actorId` se pasa directo en el payload (`019a5234-9a65-7cb0-badc-670fe28ec4eb`, agente dueño del seguimiento de prueba). El bloque `"atencion_psico"` de PASO 4 no depende de `actorId` para decidir si crea o no la remisión/tarea, así que esto no afecta lo que se está probando.

> ⚠️ **Hallazgo de setup — `form-seguimiento-data.md` está desactualizado.** El primer intento con la Sección 6 "Cierre del caso" quedó `INCOMPLETO` porque el formulario real en BD ya tiene **8 secciones**, no 6: se agregaron `Seguimiento a Entidades` (order 5) y `Identificación de Entidades` (order 6), corriendo `Seguimiento de Caso`→7 y `Cierre del caso`→8 (mismo ID `3e03f685-…`, solo cambió el order). `Seguimiento a Entidades` quedó invisible por defecto (VC de sección sobre `formState.hasCurrentEntities`, que no se envía), pero `Identificación de Entidades` sí es visible y tiene 2 preguntas directas requeridas. Se respondieron ambas en `false` para dejar sus 2 repeaters internos (min_repetitions=0) sin necesidad de entradas:
> - `f1be8fba-8fc0-4db6-8c4f-e6c764da3249` — "¿Ha puesto en conocimiento de alguna entidad/institución los hechos de violencia reportados?" = `false`
> - `0b4e51c9-82b4-4e21-afe0-da1eaea43077` — "¿Está usted de acuerdo con que SALVIA realice la activación de la Ruta de Atención?" = `false`
>
> No se actualizó `form-seguimiento-data.md` con esta estructura nueva — está fuera del alcance de este cambio y no se tocó código de ese flujo. Queda como pendiente para quien mantenga esas 2 secciones.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Remisión psicosocial cumple criterios | `qEquipos` incluye `atencion_psico` (sin `medidas_emergencia`) · `qCriteriosPsico` incluye `criterio_obligatorio` + criterios que suman ≥ 3 puntos | Se crea 1 `psychosocial_support` (`status=ACTIVE`) y 1 `case_task` nueva (`category=Psicosocial`, `type=validar_remision`, `description="Validar que remision a psicosocial es valida"`, `assigned_user_id=''`, `status=ToDo`, `psychosocial_support_id` = id de la remisión creada, `case_id`/`follow_up_id` del seguimiento) | Query directo a BD + logs |
| C2 | Remisión psicosocial NO cumple criterios | `qEquipos` incluye `atencion_psico` · `qCriteriosPsico` NO incluye `criterio_obligatorio` (o suma < 3 puntos) | NO se crea `psychosocial_support` ni `case_task` de validación | Query directo a BD + logs |
| C3 | Exclusión por medidas de emergencia | `qEquipos` incluye `atencion_psico` **y** `medidas_emergencia` · criterios psicosociales sí cumplen | La derivación psicosocial se omite completamente (ni `psychosocial_support` ni `case_task`) por la regla de exclusión existente — `medidas_emergencia` sí se procesa normalmente | Query directo a BD + logs |

---

## Cobertura

Los 3 casos se ejecutaron de punta a punta contra el backend local (`go run .`, puerto 9090) el 2026-08-13, reutilizando y re-disparando el mismo `follow_up_v2`/`form_submission` de prueba en cada corrida (reset a `PENDIENTE` + edición directa de `salvia.answer` antes de cada `saveSection`).

| Caso | Resultado | Evidencia |
|---|---|---|
| C1 | ✅ Pasó | Log: `[processFollowUp] ✅ derivacion creada -> atencion_psico (id=3ef93b04-6bae-454a-96d2-fd509eec9f81)`. BD: `psychosocial_support.status='ACTIVE'` creado; `case_task` creada con `category='Psicosocial'`, `type='validar_remision'`, `description='Validar que remision a psicosocial es valida'`, `assigned_user_id=''`, `status='ToDo'`, `psychosocial_support_id='3ef93b04-…'` — todos los campos coinciden exactamente con lo pedido |
| C2 | ✅ Pasó | Log: `[processFollowUp] derivacion psicosocial NO cumple criterios (obligatorio=true, puntos=1) — omitida`. BD: conteo de `psychosocial_support` y de `case_task type=validar_remision` para el caso se mantuvo en 1 (sin cambio respecto al baseline post-C1) |
| C3 | ✅ Pasó | Log: `[processFollowUp] atencion_psico omitida — incompatible con medidas_emergencia seleccionada al mismo tiempo` seguido de `✅ medida emergencia creada tipo=alojamiento`. BD: `psychosocial_support` y `case_task` de validación sin cambio (siguen en 1); `emergency_measure` pasó de 0 a 1 — confirma que la exclusión es selectiva y no rompe el resto del PASO 4 |

**Resultado:** 3/3 casos pasaron. La lógica nueva crea la tarea exactamente cuando (y solo cuando) se crea la remisión psicosocial, con todos los campos pedidos, y no interfiere con las demás ramas de derivación (`medidas_emergencia` sigue funcionando con la exclusión activa).

| Flujo de prueba | Casos cubiertos |
|---|---|
| Backend manual (reset PENDIENTE + edición directa de `answer` + `saveSection`) | C1, C2, C3 |

---

## Nota de alcance

No se agregaron pruebas de frontend (Playwright). El cambio es puramente backend y no altera ningún render ni interacción de UI en `hacer-seguimiento`. `case-task-modal` y `case-task-history` todavía no reconocen el `type` `validar_remision` — completar/ver esta tarea desde la UI queda pendiente de otro requerimiento (ver nota en changelog).
