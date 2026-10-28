# Detalle de Remisión — Medida de Emergencia — Interfaz

Lectura de una remisión. El color es el esmeralda de la bandeja (`#059669`), prefijo `rmd-`.

## Archivos

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/remision-medida/remision_medida_detalle.html` | Template Vue |
| `src/frontend/css/remision-medida-detalle.css` | Estilos `rmd-` |
| `src/salvia/facades/RemisionMedidaDetalleFacade.go` | Sesión, permiso y variables del template |
| `src/salvia/controller/emergency_measure_referral_controller.go` | `GET /emergency-measure-referrals/:id` |
| `src/frontend/js/components/case-tasks.js` | Tab Tareas, con el `caseICode` |
| `src/frontend/js/components/case-timeline.js` | Tab Timeline, con el `caseICode` |

## Árbol

```
RemisionMedidaDetalle  (.rmd-page)
├── Loading
├── Error                          mensaje + volver
└── Contenido
    ├── Breadcrumb
    │   ├── "← Medidas"            salviaGoBack al listado del rol
    │   └── Nombre de la víctima   texto, no link
    ├── Header
    │   ├── Inicial, nombre
    │   ├── Badge riesgo           riskLabel del caso
    │   ├── Badge estado           statusLabel
    │   ├── Badge urgencia         "Urgencia " + urgencyLabel
    │   ├── Documento · municipio
    │   ├── Agente asignado        o "Sin asignar"
    │   └── Chips de medidas
    ├── Tabs                       Info Medida | Tareas | Timeline
    ├── Info Medida
    │   ├── Datos de la remisión
    │   │   ├── Medidas, urgencia, remitido por, fecha, motivo
    │   │   ├── Motivo de devolución          solo si hay rejectionReason
    │   │   └── Agente asignado
    │   ├── Situación de la víctima
    │   │   ├── Agresor activo                aviso rojo si urgency es riesgo_inminente
    │   │   ├── Hijos menores                 form 1 más reciente
    │   │   └── Requiere intérprete           form 2 más reciente
    │   └── Asignación                        solo sm y sv → E-02
    ├── Tareas                     <case-tasks :case-id :referral-id :user-id :user-role>
    └── Timeline                   <case-timeline :case-id>
```

No hay descripción de la situación ni observaciones. No hay select de rol.

## Qué lee cada campo

| Campo en pantalla | Origen |
|---|---|
| Nombre, documento, municipio, riesgo | Mismos joins del listado |
| Estado, urgencia, medidas, agente | `emergency_measure_referral` y sus medidas |
| Remitido por | Nombre del `submitted_by`. Debajo, `submitted_by_team` si viene |
| Fecha | `referred_at` |
| Motivo | `reason`. Vacío se pinta "—" |
| Motivo de devolución | `rejection_reason`. El bloque no se pinta si está vacío |
| Agresor activo | No es una columna. "Sí — riesgo inminente" solo si `urgency` es `riesgo_inminente`; si no, "No reportado" |
| Hijos | `victim_case_form1_victim_children_number` y `victim_case_form1_victim_children_age` del form 1 con mayor id. Sin número, "No" |
| Intérprete | Código del enum `victim_case_form2_require_language_interpreter` del form 2 con mayor id. `y` es "Sí". Vacío u otro código es "No" |

Tareas pide `GET /api/v1/case-tasks?referralId=` de esta remisión. La ven y las gestionan las sesiones cuyo `role_code` coincide con `assigned_role`, y también quien tenga la tarea en `assigned_user_id`. `sv` ve todas las de la remisión. Al completar una tarea, la pantalla vuelve a pedir la remisión para refrescar el estado. Timeline sigue pidiendo el caso entero.
