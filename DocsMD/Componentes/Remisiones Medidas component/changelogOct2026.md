## 2026-10-30 — Orden, paginación, rol ro y datos de la víctima

El listado ordena por `created_at` descendente en todas las vistas. Cada página pide 10 filas al backend. La celda del caso muestra documento y celular, y ya no el id. El rol `ro` abre Mis remisiones y ve las remisiones que él envió, en cualquier estado.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/repository/emergency_measure_referral_repository.go` | Orden, conteo filtrado, `LIMIT`/`OFFSET` y teléfono |
| `src/salvia/service/emergency_measure_referral_service.go` | Página fija de 10 y `victimPhone` |
| `src/salvia/controller/emergency_measure_referral_controller.go` | Query `page` y vista `remitente` |
| `src/salvia/config/RemisionesMedidas.go` | Rol `ro` → vista `remitente` |
| `src/salvia/config/Menu.go` | Permisos de bandeja y detalle para `ro` |
| `src/frontend/js/components/remisiones-medidas-component.js` | Paginación y filtro de estados de `ro` |
| `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas_component.html` | Doc/Tel y barra de páginas |
| `src/frontend/html/sidebar/sidebar.html` | Menú Mis Remisiones Medidas para `ro` |

## 2026-10-28 — E-16, panel de tareas

El link «Tareas» de cada fila abre un panel con `<case-tasks>` filtrado por `referralId`. Props nuevas: `userId` y `userRole`. Lo usan la bandeja general y Mis Remisiones. La tarea `decidir_via_riesgo_inminente` solo aparece si el rol de la sesión es `sm` o `sv`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/remisiones-medidas-component.js` | Panel de tareas |
| `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas_component.html` | Link «Tareas» en las dos filas |
| `src/frontend/js/components/case-tasks.js` | Prop `referralId` y visibilidad por `assignedRole` |
| `DocsMD/Componentes/Remisiones Medidas component/remisiones-medidas-component-interface.md` | Props y árbol |
| `DocsMD/Componentes/Remisiones Medidas component/Flujos/flow-E16-cuando-presiona-tareas.md` | Flujo |

## 2026-10-28 — Decidir vía abre el panel

El botón «Decidir vía de atención» usa el mismo panel que «Tareas». Al completar la decisión, el listado se recarga. La tarea del agente aparece en ese panel para su usuario, y la de Cruz Roja para el rol `sc`.

