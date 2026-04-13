-- =============================================================================
-- FIX: Reemplaza la pregunta única "Barreras identificadas en este sector"
-- por 3 preguntas separadas (una por sector) con VisibilityCondition cada una.
-- Nota: form_section.id y form_section.form_id son tipo UUID en la BD.
--       El resto de IDs son varchar(36).
-- =============================================================================

DO $$
DECLARE
  v_form_id        varchar := '2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff';
  v_rg_id          varchar;
  v_s4_id          varchar;
  v_q_sector       varchar;
  v_q_old_barreras varchar;
  v_q_salud        varchar;
  v_q_justicia     varchar;
  v_q_proteccion   varchar;
BEGIN

  -- form_section usa tipo uuid → castear en comparación
  SELECT id::varchar INTO v_s4_id
  FROM salvia.form_section
  WHERE form_id = v_form_id::uuid AND "order" = 4;

  -- repeater_group usa varchar → comparación directa
  SELECT id INTO v_rg_id
  FROM salvia.repeater_group
  WHERE form_section_id = v_s4_id;

  -- question usa varchar → comparación directa
  SELECT id INTO v_q_sector
  FROM salvia.question
  WHERE form_id = v_form_id
    AND repeater_group_id = v_rg_id
    AND question_type = 'dropdown';

  SELECT id INTO v_q_old_barreras
  FROM salvia.question
  WHERE form_id = v_form_id
    AND repeater_group_id = v_rg_id
    AND question_type = 'multiple';

  RAISE NOTICE 's4=% rg=% q_sector=% q_old=%', v_s4_id, v_rg_id, v_q_sector, v_q_old_barreras;

  -- ── Soft delete pregunta antigua y sus opciones ───────────────────────────
  UPDATE salvia.option  SET deleted_at = now() WHERE question_id = v_q_old_barreras AND deleted_at IS NULL;
  UPDATE salvia.question SET deleted_at = now() WHERE id = v_q_old_barreras AND deleted_at IS NULL;

  -- ── Barreras en Salud ─────────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, "order")
  VALUES (gen_random_uuid()::varchar, v_form_id, v_s4_id, v_rg_id, 'multiple', 'Barreras identificadas en Salud', 2)
  RETURNING id INTO v_q_salud;

  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid()::varchar, v_q_salud, 'No Aplica',                                                    'salud_no_aplica',     1),
    (gen_random_uuid()::varchar, v_q_salud, 'Demoras en asignación de citas y continuidad de tratamientos', 'salud_demoras',       2),
    (gen_random_uuid()::varchar, v_q_salud, 'Dificultades de aseguramiento o afiliación en salud',          'salud_aseguramiento', 3),
    (gen_random_uuid()::varchar, v_q_salud, 'Falta de activación de protocolos para violencia sexual',      'salud_protocolos',    4),
    (gen_random_uuid()::varchar, v_q_salud, 'Negación en los servicios de urgencias',                       'salud_urgencias',     5),
    (gen_random_uuid()::varchar, v_q_salud, 'Negativa en el acceso a la IVE',                               'salud_ive',           6),
    (gen_random_uuid()::varchar, v_q_salud, 'Otras barreras en salud',                                      'salud_otras',         7);

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid()::varchar, 'QUESTION', v_q_salud, v_q_sector, 'salud', 'EQUALS');

  -- ── Barreras en Justicia ──────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, "order")
  VALUES (gen_random_uuid()::varchar, v_form_id, v_s4_id, v_rg_id, 'multiple', 'Barreras identificadas en Justicia', 3)
  RETURNING id INTO v_q_justicia;

  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid()::varchar, v_q_justicia, 'No Aplica',                                           'justicia_no_aplica', 1),
    (gen_random_uuid()::varchar, v_q_justicia, 'Ausencia de representación judicial',                 'justicia_repres',    2),
    (gen_random_uuid()::varchar, v_q_justicia, 'Cargas probatorias injustificadas o excesivas',       'justicia_probat',    3),
    (gen_random_uuid()::varchar, v_q_justicia, 'Falta de celeridad en la investigación',              'justicia_celeridad', 4),
    (gen_random_uuid()::varchar, v_q_justicia, 'Negativa para recibir la denuncia',                   'justicia_denuncia',  5),
    (gen_random_uuid()::varchar, v_q_justicia, 'Tipificación errónea del delito',                     'justicia_tipif',     6),
    (gen_random_uuid()::varchar, v_q_justicia, 'Otras barreras en justicia',                          'justicia_otras',     7);

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid()::varchar, 'QUESTION', v_q_justicia, v_q_sector, 'justicia', 'EQUALS');

  -- ── Barreras en Protección ────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, "order")
  VALUES (gen_random_uuid()::varchar, v_form_id, v_s4_id, v_rg_id, 'multiple', 'Barreras identificadas en Protección', 4)
  RETURNING id INTO v_q_proteccion;

  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid()::varchar, v_q_proteccion, 'No Aplica',                                                   'prot_no_aplica',  1),
    (gen_random_uuid()::varchar, v_q_proteccion, 'Demoras en la emisión de medidas de protección urgentes',     'prot_demoras',    2),
    (gen_random_uuid()::varchar, v_q_proteccion, 'Fallas en la valoración del riesgo',                          'prot_valoracion', 3),
    (gen_random_uuid()::varchar, v_q_proteccion, 'Incumplimiento de medidas de protección sin consecuencias',   'prot_incumplim',  4),
    (gen_random_uuid()::varchar, v_q_proteccion, 'Medidas de protección ineficaces',                            'prot_ineficaces', 5),
    (gen_random_uuid()::varchar, v_q_proteccion, 'Otras barreras en protección',                                'prot_otras',      6);

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid()::varchar, 'QUESTION', v_q_proteccion, v_q_sector, 'proteccion', 'EQUALS');

  RAISE NOTICE 'Fix completado. salud=% justicia=% proteccion=%', v_q_salud, v_q_justicia, v_q_proteccion;

END $$;
