-- =============================================================================
-- MIGRATE: Psicosocial Aug 2026 — consentimiento info+hide, agenda, exclusividad
--          "¿Es atención o solo contacto?" (multiple→single)
--
-- Aplica sobre BD ya sembrada (no re-ejecuta seed_psicosocial.sql completo).
-- UUIDs de agenda = src/internal/constants/psicosocial_agenda_questions.go
--
-- Nota: salvia.visibility_condition NO tiene deleted_at. Tampoco se filtra
-- deleted_at en question/option aquí (la BD viva del form engine puede no
-- tener esa columna aunque el modelo Go la declare).
--
-- Ejecutar con usuario con permisos INSERT/UPDATE sobre salvia.question / option /
-- visibility_condition / repeater_group / render_modification.
-- =============================================================================

BEGIN;

-- =============================================================================
-- 1. CONSENTIMIENTO PC S4 / PA S4 — info banner + single + hide posteriores
-- =============================================================================

DO $$
DECLARE
  v_pc_consent varchar := 'c3296c87-d7e3-4ea1-8a0f-cbf77c561d0f';
  v_pa_consent varchar := '7256b91e-861b-48b3-9216-2ac13f0ae889';
  v_pc_section varchar;
  v_pa_section varchar;
  v_pc_form    varchar := '439b57e6-07ea-4da4-9721-8ed28c6ca43f';
  v_pa_form    varchar := '501ab3d7-8382-4447-96a6-f463152693c6';
  v_info_pc    varchar := 'b1000001-0000-4000-8000-000000000001';
  v_info_pa    varchar := 'b1000002-0000-4000-8000-000000000002';
  v_consent_txt text := E'CONSENTIMIENTO/DESISTIMIENTO INFORMADO PARA LA ATENCIÓN PSICOSOCIAL EN LA LÍNEA 155\n\n'
    'Por medio de este documento se establecen los términos de confidencialidad, los riesgos, las excepciones de esta y los derechos de la persona usuaria durante las conversaciones que se sostengan en el marco de la atención psicosocial telefónica de la Línea 155 SALVIA.\n\n'
    '1) Es un espacio individual gratuito que tiene como objetivo propiciar una reflexión sobre las violencias que enfrentan las mujeres y personas no binarias.\n'
    '2) Durante el proceso pueden surgir momentos o temas que despierten emociones intensas.\n'
    '3) No se hará con el objetivo de diagnosticarle alguna afectación en su salud emocional.\n'
    '4) Se entenderá que solo usted conoce completamente su propia situación.\n'
    '5) Durante las sesiones se realizarán preguntas que nos permitan identificar sus necesidades.\n'
    '6) Se realizará un máximo de cuatro sesiones de atención psicosocial telefónica de una hora c/u.\n'
    '7) La información que brinde se usará para hacer seguimiento a su proceso.\n'
    '8) De acuerdo con la Ley 1090 de 2006 y Ley 53 de 1977, la información podrá ser compartida en casos judiciales o de riesgo inminente.';
  r record;
BEGIN
  SELECT form_section_id INTO v_pc_section FROM salvia.question WHERE id = v_pc_consent;
  SELECT form_section_id INTO v_pa_section FROM salvia.question WHERE id = v_pa_consent;

  IF v_pc_section IS NOT NULL THEN
    UPDATE salvia.question SET "order" = "order" + 1
    WHERE form_section_id = v_pc_section AND "order" >= 2
      AND id <> v_info_pc;

    INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
    VALUES (v_info_pc, v_pc_form, v_pc_section, 'info', v_consent_txt, FALSE, 2)
    ON CONFLICT (id) DO UPDATE SET description = EXCLUDED.description, question_type = 'info', "order" = 2;

    UPDATE salvia.question
    SET description = '¿Acepta el consentimiento informado para la atención psicosocial?',
        question_type = 'single',
        "order" = 3
    WHERE id = v_pc_consent;

    FOR r IN
      SELECT id FROM salvia.question
      WHERE form_section_id = v_pc_section AND "order" > 3
        AND question_type <> 'info'
    LOOP
      INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
      SELECT gen_random_uuid(), 'QUESTION', r.id, v_pc_consent, 'si', 'EQUALS'
      WHERE NOT EXISTS (
        SELECT 1 FROM salvia.visibility_condition
        WHERE target_id = r.id AND trigger_question_id = v_pc_consent AND trigger_value = 'si'
      );
    END LOOP;
  END IF;

  IF v_pa_section IS NOT NULL THEN
    UPDATE salvia.question SET "order" = "order" + 1
    WHERE form_section_id = v_pa_section AND "order" >= 2
      AND id <> v_info_pa;

    INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
    VALUES (v_info_pa, v_pa_form, v_pa_section, 'info', v_consent_txt, FALSE, 2)
    ON CONFLICT (id) DO UPDATE SET description = EXCLUDED.description, question_type = 'info', "order" = 2;

    UPDATE salvia.question
    SET description = '¿Acepta el consentimiento informado para la atención psicosocial?',
        question_type = 'single',
        "order" = 3
    WHERE id = v_pa_consent;

    FOR r IN
      SELECT id FROM salvia.question
      WHERE form_section_id = v_pa_section AND "order" > 3
        AND question_type <> 'info'
    LOOP
      INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
      SELECT gen_random_uuid(), 'QUESTION', r.id, v_pa_consent, 'si', 'EQUALS'
      WHERE NOT EXISTS (
        SELECT 1 FROM salvia.visibility_condition
        WHERE target_id = r.id AND trigger_question_id = v_pa_consent AND trigger_value = 'si'
      );
    END LOOP;
  END IF;
END $$;

-- =============================================================================
-- 2. AGENDA — ¿Agendar nueva sesión? + Hora (UUIDs fijos) junto a fechas existentes
-- =============================================================================

DO $$
DECLARE
  slots text[][] := ARRAY[
    ARRAY['a1000001-0000-4000-8000-000000000001','a2000001-0000-4000-8000-000000000001','58ce2d34-24d2-4e73-bf95-26a2c608f8e6','439b57e6-07ea-4da4-9721-8ed28c6ca43f'],
    ARRAY['a1000002-0000-4000-8000-000000000002','a2000002-0000-4000-8000-000000000002','fd2fb664-de83-4069-a161-6348dfef48bf','439b57e6-07ea-4da4-9721-8ed28c6ca43f'],
    ARRAY['a1000003-0000-4000-8000-000000000003','a2000003-0000-4000-8000-000000000003','64d63b79-edee-464b-be56-1104efd31a46','501ab3d7-8382-4447-96a6-f463152693c6'],
    ARRAY['a1000004-0000-4000-8000-000000000004','a2000004-0000-4000-8000-000000000004','d68c7334-74bb-47b1-a2ee-f1d3da04627b','501ab3d7-8382-4447-96a6-f463152693c6'],
    ARRAY['a1000005-0000-4000-8000-000000000005','a2000005-0000-4000-8000-000000000005','72ce49d2-f853-4f4a-9f1a-f795f4d514c3','7a7b61b3-7bc3-4088-a74d-0974e54a3563'],
    ARRAY['a1000006-0000-4000-8000-000000000006','a2000006-0000-4000-8000-000000000006','7040a37d-f346-4bf7-9b80-6286b9da62c5','7a7b61b3-7bc3-4088-a74d-0974e54a3563'],
    ARRAY['a1000007-0000-4000-8000-000000000007','a2000007-0000-4000-8000-000000000007','1b4d09f0-e5ca-4b4e-9d74-428478a113c6','c31026f8-7ce2-49b9-8990-04b44ed4513c'],
    ARRAY['a1000008-0000-4000-8000-000000000008','a2000008-0000-4000-8000-000000000008','3ce6ff5a-f139-4417-9255-155207e9a970','c31026f8-7ce2-49b9-8990-04b44ed4513c']
  ];
  i int;
  v_agendar varchar;
  v_hora    varchar;
  v_fecha   varchar;
  v_form    varchar;
  v_section varchar;
  v_order   int;
BEGIN
  FOR i IN 1..array_length(slots, 1) LOOP
    v_agendar := slots[i][1];
    v_hora    := slots[i][2];
    v_fecha   := slots[i][3];
    v_form    := slots[i][4];

    SELECT form_section_id, "order" INTO v_section, v_order
    FROM salvia.question WHERE id = v_fecha;

    IF v_section IS NULL THEN
      RAISE NOTICE 'Fecha % no encontrada — skip slot %', v_fecha, i;
      CONTINUE;
    END IF;

    -- Idempotencia: si Agendar ya existe, no re-desplazar orders
    IF EXISTS (SELECT 1 FROM salvia.question WHERE id = v_agendar) THEN
      RAISE NOTICE 'Agenda slot % ya aplicado (agendar=%) — skip shift', i, v_agendar;
      CONTINUE;
    END IF;

    UPDATE salvia.question SET "order" = "order" + 2
    WHERE form_section_id = v_section AND "order" > v_order
      AND id NOT IN (v_agendar, v_hora, v_fecha);

    UPDATE salvia.question SET "order" = v_order + 1 WHERE id = v_fecha;

    INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
    VALUES (v_agendar, v_form, v_section, 'single', '¿Agendar nueva sesión?', TRUE, v_order)
    ON CONFLICT (id) DO UPDATE SET description = EXCLUDED.description, "order" = EXCLUDED."order";

    INSERT INTO salvia.option (id, question_id, label, value, "order")
    SELECT gen_random_uuid(), v_agendar, 'Sí', 'si', 1
    WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_agendar AND value = 'si');
    INSERT INTO salvia.option (id, question_id, label, value, "order")
    SELECT gen_random_uuid(), v_agendar, 'No', 'no', 2
    WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_agendar AND value = 'no');

    INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
    VALUES (v_hora, v_form, v_section, 'time', 'Hora próxima atención', TRUE, v_order + 2)
    ON CONFLICT (id) DO UPDATE SET description = EXCLUDED.description, question_type = 'time', "order" = EXCLUDED."order";

    INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
    SELECT gen_random_uuid(), 'QUESTION', v_agendar, vc.trigger_question_id, vc.trigger_value, vc.operator
    FROM salvia.visibility_condition vc
    WHERE vc.target_id = v_fecha
      AND NOT EXISTS (
        SELECT 1 FROM salvia.visibility_condition x
        WHERE x.target_id = v_agendar AND x.trigger_question_id = vc.trigger_question_id
          AND x.trigger_value = vc.trigger_value
      );

    INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
    SELECT gen_random_uuid(), 'QUESTION', v_fecha, v_agendar, 'si', 'EQUALS'
    WHERE NOT EXISTS (
      SELECT 1 FROM salvia.visibility_condition
      WHERE target_id = v_fecha AND trigger_question_id = v_agendar AND trigger_value = 'si'
    );

    INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
    SELECT gen_random_uuid(), 'QUESTION', v_hora, v_agendar, 'si', 'EQUALS'
    WHERE NOT EXISTS (
      SELECT 1 FROM salvia.visibility_condition
      WHERE target_id = v_hora AND trigger_question_id = v_agendar AND trigger_value = 'si'
    );

    RAISE NOTICE 'Agenda slot % OK — section=% fecha=% agendar=% hora=%', i, v_section, v_fecha, v_agendar, v_hora;
  END LOOP;
END $$;

-- =============================================================================
-- 3. state_items / render_modification Seguimiento a Barreras (idempotente)
-- =============================================================================

UPDATE salvia.repeater_group SET state_items = 'currentBarriers'
WHERE id IN (
  'cb2fb8b8-0431-4088-b626-06e1069432fa',
  '10122569-e6c2-4449-a7ff-add378815eb8',
  '94718173-b809-41e4-86d5-378de3434dca',
  '3cb23f26-af5a-4052-8987-28294f68e5dc'
)
AND (state_items IS NULL OR state_items <> 'currentBarriers');

INSERT INTO salvia.render_modification (target_type, target_id, target_field, modification_type, state_path)
SELECT 'question', q.id, 'description', 'SET', 'currentBarriers.{_entryIndex}.barrierName'
FROM salvia.question q
WHERE q.repeater_group_id IN (
  'cb2fb8b8-0431-4088-b626-06e1069432fa',
  '10122569-e6c2-4449-a7ff-add378815eb8',
  '94718173-b809-41e4-86d5-378de3434dca',
  '3cb23f26-af5a-4052-8987-28294f68e5dc'
)
AND q.question_type = 'info'
AND q.description ILIKE 'Seguimiento a Barrera%'
AND NOT EXISTS (
  SELECT 1 FROM salvia.render_modification rm
  WHERE rm.target_id = q.id AND rm.state_path = 'currentBarriers.{_entryIndex}.barrierName'
);

-- =============================================================================
-- 4. ¿Es atención o solo contacto? — exclusive choice (multiple → single)
--    PA / Atención Psicosocial / Cierre. Idempotente.
-- =============================================================================

UPDATE salvia.question
SET question_type = 'single'
WHERE description = '¿Es atención o solo contacto?'
  AND question_type = 'multiple';

-- Respuestas ya guardadas con ambas opciones (CSV de multiple) → un solo valor.
-- Preferimos solo_contacto cuando ambas están presentes (gatillo de "contacto sin atención").
UPDATE salvia.answer
SET value = 'solo_contacto'
WHERE question_id IN (
  SELECT id FROM salvia.question WHERE description = '¿Es atención o solo contacto?'
)
AND value LIKE '%atencion%'
AND value LIKE '%solo_contacto%';

COMMIT;
