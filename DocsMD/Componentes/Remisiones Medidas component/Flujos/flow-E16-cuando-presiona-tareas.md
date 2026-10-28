━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario abre las tareas de la remisión
   Tipo: User Interaction
   Código: E-16
   Función: abrirTareas(item) / cerrarTareas()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  - el link "Tareas" de la fila, en la bandeja general y en Mis Remisiones
  - el botón "Decidir vía de atención" del supervisor
No navega. No emite al padre.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Abrir el panel

  panelTareas = la fila
  El panel muestra "Tareas · {nombre}" y el componente:

    <case-tasks
      :case-id="panelTareas.caseICode"
      :referral-id="panelTareas.remisionId"
      :user-id="userId"
      :user-role="userRole"
    >

  userId y userRole llegan del padre (la sesión).


PASO 2 — Pedir las tareas de esa remisión

  → GET /api/v1/case-tasks?referralId={remisionId}

  Lo pueden pedir sm, em, sc, cr, am y sv.
  Si la sesión no es sv, la respuesta trae las tareas cuyo assigned_role
  coincide con el role_code y las cuyo assigned_user_id es el de la sesión.
  sv recibe todas las de la remisión.
  No abre la bandeja completa del caso.


PASO 3 — Qué se ve

  La tarea decidir_via_riesgo_inminente tiene assigned_role = sm.
  La ve y la gestiona sm. También sv.

  La tarea intentar_contacto_3x3 no tiene rol: la ve el agente cuyo
  assigned_user_id es el de la sesión, y también sv.

  La tarea escalar_cruz_roja tiene assigned_role = sc.
  La ve sc y sv.

  "Asignado a" muestra el nombre del rol o el del usuario.
  "Gestionar" llama a case-task-modal.open(taskId).


PASO 4 — Al completar una tarea

  case-tasks emite task-completed.
  El listado vuelve a pedir las remisiones.
  La fila muestra el estado nuevo. El panel sigue abierto.


PASO 5 — Cerrar

  "Cerrar" limpia panelTareas.
  No escribe la tarea ni cambia la remisión.
