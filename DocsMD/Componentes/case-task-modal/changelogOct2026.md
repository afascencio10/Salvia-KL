# Changelog — case-task-modal — octubre 2026

## 2026-10-28 — Tipo `decidir_via_riesgo_inminente`

El modal abre la tarea de decidir vía. Muestra Vía A y Vía B como checkboxes y el texto «Puedes elegir una opción o solicitar la gestión simultánea de ambas.». Cancelar, la X y el overlay cierran sin escribir. «Confirmar decisión» se pinta deshabilitado y no llama a `PUT /case-tasks/:id/complete`.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | Formulario del tipo. `formularioValido` queda en false |
| `DocsMD/Componentes/case-task-modal/case-task-modal-interface.md` | Título, formulario y botón |
| `DocsMD/Componentes/case-task-modal/case-task-modal-index.md` | E11 |
| `DocsMD/Componentes/case-task-modal/Flujos/flow-E11-cuando-abre-decidir-via.md` | Apertura |
| `DocsMD/Componentes/case-task-modal/Flujos/flow-E09-cuando-confirma-formulario.md` | Este tipo no persiste |

## 2026-10-28 — Confirmar la decisión de vía

Las vías son radios: solo una. El select «Asignar agente» carga los Agentes de Medidas activos. Confirmar decisión exige vía y agente, y llama a `PUT /case-tasks/:id/complete`.

En una transacción se asigna el agente, la remisión pasa a `asignada` (Vía A) o a `riesgo_inminente_por_escalar_cruz_roja` (Vía B), se crea la tarea siguiente y la de decidir vía queda Done.

Vía A crea `intentar_contacto_3x3` para el agente. Vía B crea `escalar_cruz_roja` para el rol `sc`. Guardar y Confirmar de esos dos modales cierran y no escriben.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | Radios, agente y los dos modales de solo cierre |
| `src/salvia/service/case_task_service.go` | `completeDecidirVia` |
| `src/internal/models/case_task.go` | Tipos y descripciones de las tareas siguientes |

## 2026-10-28 — Timeline al decidir la vía

Al confirmar la decisión, el caso recibe «Vía de atención decidida». La descripción es «Vía A — asignación al agente» o «Vía B — escalar a Cruz Roja». El evento se escribe después de la transacción: si falla, la decisión queda.

La tarea `escalar_cruz_roja` nace con la descripción «Escalar a Cruz Roja — Riesgo inminente». El título de la fila y el del modal siguen en «Escalar a Cruz Roja».

## 2026-10-28 — Confirmar el escalamiento a Cruz Roja

Al confirmar, la remisión pasa a `riesgo_inminente_por_entregar_medida` y se crea la tarea `entregar_medida_inmediata` para el rol `cr`, con la descripción «Entregar medida inmediata (sin caracterización previa)». El caso recibe «Escalamiento a Cruz Roja confirmado»: «Supervisión del Convenio confirmó el escalamiento a Cruz Roja». Si el timeline falla, la confirmación queda.

Esa tarea abre el modal «Entregar medida inmediata» con las medidas de la remisión y comentarios. Registrar entrega cierra y no guarda.


## 2026-10-28 — Confirmar el escalamiento a Cruz Roja

Al confirmar, la remisión pasa a `riesgo_inminente_por_entregar_medida` y se crea la tarea `entregar_medida_inmediata` para el rol `cr`, con la descripción «Entregar medida inmediata (sin caracterización previa)». El caso recibe «Escalamiento a Cruz Roja confirmado»: «Supervisión del Convenio confirmó el escalamiento a Cruz Roja». Si el timeline falla, el escalamiento queda.

El modal de esa tarea muestra las medidas remitidas y un campo de comentarios. «Registrar entrega» cierra y no guarda.

## 2026-10-28 — Registrar la entrega inmediata

Cruz Roja elige las medidas que entrega. El comentario puede ir vacío. «Registrar entrega» llama a `PUT /case-tasks/:id/complete`.

En una transacción, las medidas elegidas quedan en `ENTREGADA`, la remisión pasa a `riesgo_inminente_por_confirmar_seguridad` y se crea `confirmar_lugar_seguro` para el agente ya asignado. El caso recibe «Medida entregada por Cruz Roja», con los nombres de las medidas y el comentario si lo hay. Si el timeline falla, la entrega queda.

La tarea nueva abre el modal con «Confirma entrega de medidas». Confirmar y Cancelar cierran y no guardan.

## 2026-10-29 — Confirmar la entrega de la medida

Al confirmar, la remisión pasa a `asignada` y se crea `intentar_contacto_3x3` para el mismo agente. El caso recibe «Entrega de medida confirmada»: «El agente de medidas confirmó la entrega de la medida.» Si el timeline falla, la confirmación queda.

El 3×3 abre su modal. Guardar cierra y no persiste.

## 2026-10-29 — Revisión del Equipo Técnico

Al guardar un seguimiento viejo o V.4 con medidas que no son de riesgo inminente, nace `revision_equipo_tecnico` para el rol `em`. El caso recibe «Medidas enviadas a revisión» con los nombres de las medidas. Si la tarea no se crea, no hay timeline y la remisión queda. Las derivaciones psicosociales no crean esta tarea.

El modal es «Revisión Equipo Técnico», sin el select de profesional. Guardar revisión exige concepto y un resultado. Aprobar deja la remisión por asignar. Devolver se describe más abajo.

Al aprobar, la remisión pasa a `por_asignar` y se crea `asignar_agente_medidas` para el rol `sm`. El caso recibe «Enviada a asignación de agente»: «Se envió para asignación de agente.» Si el timeline falla, la aprobación queda.

Esa tarea abre «Reasignar remisiones» con el select de agentes. Reasignar asigna el agente, deja la remisión en `asignada` y crea `intentar_contacto_3x3`. El caso recibe «Agente de medidas asignado»: «Se asignó el agente de medidas.» Si el timeline falla, la asignación queda.

En la bandeja, «Asignar agente» abre esa tarea. La barra «Reasignar» sigue en su modal: si la remisión estaba en `por_asignar`, nace el mismo 3×3 y se cierra la tarea; en otro estado solo cambia el agente y pasa a `asignada`.

## 2026-10-29 — Devolución y reenvío

Al devolver para ajuste o para acompañamiento, la remisión pasa a `devuelta_para_justificacion` o `devuelta_a_seguimiento` y nace `ajustar_remision_devuelta` para quien remitió. El caso recibe «Devolución del Equipo Técnico», con el concepto.

El modal muestra el hilo de `GET /case-tasks/:id/chat` y el último comentario del equipo. Reenviar vuelve la remisión a `pendiente_de_revision` y crea otra revisión para `em`. Descartar la deja en `descartada`, pasa a `CLOSED` las medidas `ACTIVE` de esa remisión y no crea revisión. Si el timeline falla, la operación queda.

La revisión que ya tiene vueltas muestra el mismo hilo encima del concepto.

## 2026-10-29 — Contacto 3×3

Guardar el intento exige Sí o No. Si contestó, nace `entrevista_caracterizacion` para el mismo agente y el caso recibe «Contacto con la víctima». Ese modal abre y Guardar no persiste.

Si no contestó, el intento queda cerrado. En el primero y en el segundo nace otro `intentar_contacto_3x3`. A la tercera, la remisión pasa a `cerrada_por_no_contacto`, las medidas `ACTIVE` pasan a `CLOSED` y el caso recibe «Remisión cerrada por no contacto». Si el timeline falla, el intento queda.

Arriba de las tareas de la remisión se ven los intentos fallidos, de 0 a 3.



