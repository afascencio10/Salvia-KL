## 2026-10-28 — Panel de tareas por fila

Cada fila tiene el link «Tareas». Abre un panel en la misma pantalla y monta `case-tasks` de esa remisión. No navega. La pantalla pasa `user-id` y `user-role` de la sesión. «Gestionar» abre el modal. La tarea de decidir vía aparece para `sm` y `sv`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas.html` | Props de sesión y scripts de `case-tasks` |
| `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas_component.html` | Link y panel |
| `src/frontend/js/components/remisiones-medidas-component.js` | `abrirTareas` / `cerrarTareas` |
| `src/salvia/facades/RemisionesMedidasFacade.go` | `currentUserId` |
| `DocsMD/Screens/Remisiones Medidas/remisiones-medidas-index.md` | E-16 |
| `DocsMD/Componentes/Remisiones Medidas component/remisiones-medidas-component-events.md` | E-16 |

## 2026-10-28 — Decidir vía abre el panel de tareas

El botón «Decidir vía de atención» deja de ir al detalle y abre el mismo panel de tareas. Al guardar la decisión, el listado se vuelve a pedir y la fila muestra Asignada o por escalar a Cruz Roja. «Escalar a Cruz Roja» sigue yendo al detalle.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/remisiones-medidas-component.js` | El botón es tipo `tareas` |
| `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas_component.html` | Recarga el listado al completar una tarea |
| `src/salvia/controller/case_task_controller.go` | El listado por remisión incluye la tarea del usuario de la sesión |

