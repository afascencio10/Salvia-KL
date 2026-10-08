-- =============================================================================
-- REPARACIÓN EN ARRANQUE: tarea de validación faltante en remisiones psicosociales
--
-- POR QUÉ EXISTE
-- La pantalla de detalle habilita "recepcionar / devolver" solo cuando la
-- remisión tiene una case_task 'validar_remision' en ToDo (ver validationPending
-- en psychosocial_detail_service.go). Las remisiones creadas por el endpoint
-- manual nacían sin ella, así que no había forma de validarlas. Eso ya se corrigió
-- en el código, pero las remisiones creadas antes siguen rotas.
--
-- Esto no es solo un parche retroactivo: se deja corriendo en cada arranque como
-- red de seguridad. En la vía del seguimiento la tarea se crea "best effort" (si
-- falla, solo queda un WARN en el log), y nada impide que mañana aparezca una
-- tercera vía de creación que vuelva a olvidarla. Esta consulta repara cualquiera
-- de esos casos sin que nadie tenga que darse cuenta.
--
-- Deliberadamente NO lleva lista de cédulas ni de ids: repara la condición, no
-- unos casos puntuales. Mantener documentos de víctimas dentro del repositorio no
-- es aceptable, y además no hace falta.
--
-- ES UNA SOLA SENTENCIA: db.Exec lo exige, y RowsAffected dice cuántas se crearon.
--
-- QUÉ DEJA FUERA
--   - Remisiones cerradas o en devolución: no hay nada que validar.
--   - Las que ya tienen la tarea (ToDo o Done): no duplica ni vuelve a pedir una
--     validación ya hecha.
--   - Las que tienen un evento "Remisión Validada": se validaron alguna vez,
--     aunque la tarea no esté.
--   - Las que ya tienen contactos registrados: están siendo atendidas, y crearles
--     la tarea ahora les bloquearía el registro de contactos hasta validar
--     (puedeGestionarContactos exige que no haya validación pendiente).
--     Decisión del equipo: esas se revisan a mano, no se tocan desde acá.
-- =============================================================================

INSERT INTO salvia.case_task
    (id, category, type, description, assigned_user_id, status,
     case_id, follow_up_id, psychosocial_support_id, created_at, updated_at)
SELECT gen_random_uuid(),
       'Psicosocial',
       'validar_remision',
       'Validar que remision a psicosocial es valida',
       '',          -- sin asignar, igual que las que crea el flujo de seguimiento
       'ToDo',
       ps.case_id,
       ps.follow_up_id,
       ps.id,
       NOW(),
       NOW()
FROM salvia.psychosocial_support ps
WHERE ps.deleted_at IS NULL
  AND ps.status NOT IN ('cerrado', 'en_devolucion')
  AND NOT EXISTS (
      SELECT 1 FROM salvia.case_task t
      WHERE t.psychosocial_support_id = ps.id
        AND t.type = 'validar_remision'
        AND t.deleted_at IS NULL
  )
  AND NOT EXISTS (
      SELECT 1 FROM salvia.case_timeline_event e
      WHERE e.psychosocial_support_id = ps.id
        AND e.type = 'Remisión Validada'
  )
  AND NOT EXISTS (
      SELECT 1 FROM salvia.team_contact tc
      WHERE tc.psicosocial_id = ps.id::text
        AND tc.deleted_at IS NULL
  );
