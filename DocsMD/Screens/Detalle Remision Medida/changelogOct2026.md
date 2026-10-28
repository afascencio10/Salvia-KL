## 2026-10-28 — El tab Tareas filtra por la remisión

El tab Tareas deja de pedir el caso entero. Pasa `referral-id` y `case-tasks` llama a `GET /api/v1/case-tasks?referralId=`. La tarea de decidir vía la ven y la gestionan `sm` y `sv`. El resto de roles no la ve. «Gestionar» abre el modal de decidir vía.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/remision-medida/remision_medida_detalle.html` | `<case-tasks>` con `referral-id` |
| `src/salvia/controller/case_task_controller.go` | Listado por `referralId` y por `assigned_role` |
| `DocsMD/Screens/Detalle Remision Medida/detalle-remision-medida-interface.md` | Tab Tareas |
| `DocsMD/Screens/Detalle Remision Medida/detalle-remision-medida-index.md` | Filtro por remisión |

## 2026-10-28 — El detalle refresca el estado al guardar la vía

Al completar una tarea del tab, la pantalla vuelve a pedir la remisión sin desmontar el tab. El agente asignado ve su tarea `intentar_contacto_3x3` porque el listado por remisión también entrega las tareas de su usuario.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/remision-medida/remision_medida_detalle.html` | `refrescar` al completar una tarea |
| `src/salvia/controller/case_task_controller.go` | Filtro por `assigned_user_id` además del rol |

