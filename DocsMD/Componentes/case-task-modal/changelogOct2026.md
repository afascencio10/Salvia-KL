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

