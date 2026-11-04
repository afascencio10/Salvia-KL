━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 SUB-FLUJO: Materializar el borrador de entidades y barreras (no es un evento de usuario)
   Tipo: Backend — corre dentro de «Cuando guarda una sección/módulo» (E12), cuando el Caso pasa de Borrador a Iniciado
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-38*)
   Funciones: draftFlushRepoV3.UpdateStatusIfCurrent() · flushDraft() · materializarEntidadV3() → victim_case_v3_draft.go · crearTareasMaterializacionV3() · procesarCanalesV3() · procesarGestionBarreraV3() · crearRemisionesV3() → victim_case_v3_tasks.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Disparador: `SaveSection` completa la tríada y llama `repo.UpdateStatusIfCurrent(icode, "bo", "in")`. El servicio
> envuelve el repositorio (opción `WithVictimCaseV3Draft`) con un decorador que **vacía el borrador antes de cambiar el
> estado**, sin editar `SaveSection`. Si el vaciado falla, el estado no cambia y el guardado responde 500: se reintenta
> en el siguiente guardado.

INPUT: victim_case.victim_case_v3_draft_data (entities, barriers, entityMeta, barrierMeta, materializedAt) + respuestas GS_5–GS_12 de victim_case_form2

PASO 1 — Transacción de filas (una sola; si algo falla se revierte todo y el Caso sigue en `bo`)
  SELECT … FOR UPDATE de la fila de victim_case; SI el estado ya no es 'bo': no hace nada (otro guardado ganó la carrera)
  PARA CADA entidad del JSON (en orden):
    SI la sede ya está en entity_case (activa) → se reutiliza; SI NO → entity_case { case_id, entity_branch_id, objetivo, solicitud, radicado = fileNumber, created_by = by, created_at = createdAt }
    PARA CADA id de completed → entity_case_obligation 'completada'; de pending → 'pendiente' (sin «otra»; sin duplicar una acción ya existente)
    entityMeta[entityBranchId] = { kind, agreement, dueDiligence, activationChannels, reportDate, managementNotes, pending, completed }
  PARA CADA barrera del JSON → barrier_v2 { status OPEN, sector, specific_barriers, specific_institutions, other_barrier_desc, structural_*, barrier_date, official_dependency, department_id, city_id, description, management_actions, entity_branch_id }; barrierMeta[id] = { zone, professionalPerspective, managementDescription }
  entities = []; barriers = []; materializedAt = ahora (solo la primera vez); COMMIT

PASO 2 — Tareas y oficios (después del COMMIT, mejor esfuerzo: un fallo se registra en el log y no revierte las filas)
  Canales de cada entidad, solo si agreement Y dueDiligence (y, en una acudida, si tiene acciones pendientes):
    notificacion → entity_letter (por_proyectar, agente = by, sede) + case_task Entidades · proyectar_oficio «Proyectar oficio de activación de ruta» (con entity_case_id y entity_letter_id)
    misalvia     → case_task Entidades · solicitud_entidad «Hacer solicitudes a la entidad sobre obligaciones pendientes»
    articulacion → case_task Entidades · gestion_llamada «Articulación institucional telefónica o por mensajería»
    denuncia     → case_task Entidades · gestion_llamada «Denuncia fácil»
    ninguna      → nada
  Gestión de cada barrera:
    activacion → entity_letter (barrier_id, por_proyectar) + case_task Barreras · proyectar_oficio «Activación de ruta interinstitucional (Generar oficio)»
    alerta     → case_task Barreras · comite_caso «Alerta por barreras» (sin responsable: pool del comité)
  Cada tarea: Status ToDo, asignada a quien registró (by), sin follow_up_id

PASO 3 — Remisiones a otros equipos (una sola vez por Caso: solo si materializedAt era nulo; se omiten si GS_5 ≠ Sí)
  Lee GS_6 (equipos), GS_7 (medidas) y GS_9–GS_12 (criterios) de victim_case_form2 y aplica las reglas de Seguimiento (`processEquiposYMedidas`):
    ap Atención psicosocial → psychosocial_support (derivacion/abierto, criterios, remitente y su equipo) + case_task Psicosocial · validar_remision; requiere criterio obligatorio y 3 puntos; se omite si también se pidió «me»
    ah Atención hombres → men_team_remision (derivacion/ACTIVE) si GS_10 tiene el criterio obligatorio
    di Discapacidad → un discapacidad_remision por servicio (apoyo_lsc, enfoque_discapacidad)
    es Estabilización → economic_stabilization (derivacion/ACTIVE) con al menos un criterio
    me Medidas de emergencia → emergency_measure_referral (pendiente_de_revision, normal) + una emergency_measure por tipo (alojamiento, transporte, alimentacion, vestuario, apoyo_psico, otras_me)
    sd Salvia Dignidad → sin registro (igual que Seguimiento)

PASO 4 — Después: el Caso queda en Iniciado; la sección de entidades, barreras y remisiones es de solo lectura (AddEntity/RemoveEntity/AddBarrier/RemoveBarrier → 409 readOnly); GetCase devuelve entitiesReadOnly = true y las filas de las tablas con su entityMeta/barrierMeta

REINTENTO AL LEER: si un Caso ya fuera de borrador conserva JSON pendiente (el vaciado anterior falló a mitad), GetCase / ListEntities / ListBarriers lo materializan (pasos 1 y 2; **sin remisiones**). Es idempotente: una segunda lectura no duplica nada

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Las respuestas de remisión (GS_5–GS_12) deben quedar guardadas ANTES de completar la tríada: después no se crean remisiones (D20; el backend no bloquea las respuestas, solo la pantalla) | PASO 3 |
| Las tareas, oficios y remisiones replican los de Seguimiento en archivos nuevos: conviene extraerlos a un helper compartido | PASO 2–3 |
| Un fallo después del COMMIT (tareas, oficios, remisiones) solo se registra en el log; no hay reintento automático | PASO 2–3 |
