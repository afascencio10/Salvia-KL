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

-- =============================================================================
-- 5. DERIVACIÓN A EQUIPOS SALVIA — Sep 2026
--
-- Agrega la sección "Derivación a equipos SALVIA" en los 4 formularios psicosociales
-- como última sección, visible solo cuando la sesión fue atención efectiva.
-- UUIDs fijos de preguntas: serie d1000NXX (N=formulario, XX=orden).
-- Idempotente: ON CONFLICT (id) DO NOTHING en question, INSERT WHERE NOT EXISTS
-- en form_section y visibility_condition.
-- =============================================================================

DO $$
DECLARE
  v_form_pc  uuid := '439b57e6-07ea-4da4-9721-8ed28c6ca43f';
  v_form_pa  uuid := '501ab3d7-8382-4447-96a6-f463152693c6';
  v_form_seg uuid := '7a7b61b3-7bc3-4088-a74d-0974e54a3563';
  v_form_cie uuid := 'c31026f8-7ce2-49b9-8990-04b44ed4513c';

  -- Sección IDs (generados al vuelo — idempotencia via nombre+form_id)
  v_pc_deriv_s  uuid;
  v_pa_deriv_s  uuid;
  v_seg_deriv_s uuid;
  v_cie_deriv_s uuid;

  -- Question IDs fijos (serie d1000NXX)
  -- PC
  v_pc_q1  varchar := 'd1000101-0000-4000-8000-000000000001';
  v_pc_q2  varchar := 'd1000102-0000-4000-8000-000000000001';
  v_pc_q3  varchar := 'd1000103-0000-4000-8000-000000000001';
  v_pc_q4  varchar := 'd1000104-0000-4000-8000-000000000001';
  v_pc_q5  varchar := 'd1000105-0000-4000-8000-000000000001';
  v_pc_q6  varchar := 'd1000106-0000-4000-8000-000000000001';
  v_pc_q7  varchar := 'd1000107-0000-4000-8000-000000000001';
  -- PA
  v_pa_q1  varchar := 'd1000201-0000-4000-8000-000000000002';
  v_pa_q2  varchar := 'd1000202-0000-4000-8000-000000000002';
  v_pa_q3  varchar := 'd1000203-0000-4000-8000-000000000002';
  v_pa_q4  varchar := 'd1000204-0000-4000-8000-000000000002';
  v_pa_q5  varchar := 'd1000205-0000-4000-8000-000000000002';
  v_pa_q6  varchar := 'd1000206-0000-4000-8000-000000000002';
  v_pa_q7  varchar := 'd1000207-0000-4000-8000-000000000002';
  -- SEG
  v_seg_q1 varchar := 'd1000301-0000-4000-8000-000000000003';
  v_seg_q2 varchar := 'd1000302-0000-4000-8000-000000000003';
  v_seg_q3 varchar := 'd1000303-0000-4000-8000-000000000003';
  v_seg_q4 varchar := 'd1000304-0000-4000-8000-000000000003';
  v_seg_q5 varchar := 'd1000305-0000-4000-8000-000000000003';
  v_seg_q6 varchar := 'd1000306-0000-4000-8000-000000000003';
  v_seg_q7 varchar := 'd1000307-0000-4000-8000-000000000003';
  -- CIE
  v_cie_q1 varchar := 'd1000401-0000-4000-8000-000000000004';
  v_cie_q2 varchar := 'd1000402-0000-4000-8000-000000000004';
  v_cie_q3 varchar := 'd1000403-0000-4000-8000-000000000004';
  v_cie_q4 varchar := 'd1000404-0000-4000-8000-000000000004';
  v_cie_q5 varchar := 'd1000405-0000-4000-8000-000000000004';
  v_cie_q6 varchar := 'd1000406-0000-4000-8000-000000000004';
  v_cie_q7 varchar := 'd1000407-0000-4000-8000-000000000004';

  -- IDs de las preguntas gatillo (ya existentes en la BD)
  -- PC:  Q12 "Continuar Primera Atención"
  v_pc_trigger  varchar := 'aa46351a-f6c0-4666-af3e-e0f1f7b6dc8c';
  -- PA:  Q8  "¿Es atención o solo contacto?"
  v_pa_trigger  varchar := '3117fd01-6a31-4595-9888-dcdcc96d84e7';
  -- SEG: Q9  "¿Es atención o solo contacto?"
  v_seg_trigger varchar := '16de7674-4349-4acc-9d0d-656d1d5e5f10';
  -- CIE: Q9  "¿Es atención o solo contacto?"
  v_cie_trigger varchar := 'be185336-e025-4768-a1a2-3ab9e2281450';

BEGIN

  -- ── 1. Crear las 4 secciones (idempotente por nombre + form_id) ─────────────

  SELECT id INTO v_pc_deriv_s FROM salvia.form_section
  WHERE form_id = v_form_pc AND name = 'Derivación a equipos SALVIA' AND deleted_at IS NULL
  LIMIT 1;
  IF v_pc_deriv_s IS NULL THEN
    INSERT INTO salvia.form_section (id, form_id, name, description, "order")
    VALUES (gen_random_uuid(), v_form_pc,
            'Derivación a equipos SALVIA',
            'Visible cuando la sesión fue atención efectiva (trigger: Q12 Continuar Primera Atención = true). '
            'Permite registrar necesidades de remisión a equipos SALVIA generadas durante la atención.', 5)
    RETURNING id INTO v_pc_deriv_s;
  END IF;

  SELECT id INTO v_pa_deriv_s FROM salvia.form_section
  WHERE form_id = v_form_pa AND name = 'Derivación a equipos SALVIA' AND deleted_at IS NULL
  LIMIT 1;
  IF v_pa_deriv_s IS NULL THEN
    INSERT INTO salvia.form_section (id, form_id, name, description, "order")
    VALUES (gen_random_uuid(), v_form_pa,
            'Derivación a equipos SALVIA',
            'Visible cuando la sesión fue atención efectiva (trigger: ¿Es atención o solo contacto? = atencion). '
            'Permite registrar necesidades de remisión a equipos SALVIA generadas durante la atención.', 7)
    RETURNING id INTO v_pa_deriv_s;
  END IF;

  SELECT id INTO v_seg_deriv_s FROM salvia.form_section
  WHERE form_id = v_form_seg AND name = 'Derivación a equipos SALVIA' AND deleted_at IS NULL
  LIMIT 1;
  IF v_seg_deriv_s IS NULL THEN
    INSERT INTO salvia.form_section (id, form_id, name, description, "order")
    VALUES (gen_random_uuid(), v_form_seg,
            'Derivación a equipos SALVIA',
            'Visible cuando la sesión fue atención efectiva (trigger: ¿Es atención o solo contacto? = atencion). '
            'Permite registrar necesidades de remisión a equipos SALVIA generadas durante la atención.', 7)
    RETURNING id INTO v_seg_deriv_s;
  END IF;

  SELECT id INTO v_cie_deriv_s FROM salvia.form_section
  WHERE form_id = v_form_cie AND name = 'Derivación a equipos SALVIA' AND deleted_at IS NULL
  LIMIT 1;
  IF v_cie_deriv_s IS NULL THEN
    INSERT INTO salvia.form_section (id, form_id, name, description, "order")
    VALUES (gen_random_uuid(), v_form_cie,
            'Derivación a equipos SALVIA',
            'Visible cuando la sesión fue atención efectiva (trigger: ¿Es atención o solo contacto? = atencion). '
            'Se ubica después de la sección de atención y antes de Cierre.', 7)
    RETURNING id INTO v_cie_deriv_s;
  END IF;

  -- Orden único: Derivación después de atención, Cierre al final.
  UPDATE salvia.form_section SET "order" = 6
   WHERE form_id = v_form_pc AND name = 'Cierre' AND deleted_at IS NULL;
  UPDATE salvia.form_section SET "order" = 7
   WHERE form_id IN (v_form_pa, v_form_seg, v_form_cie)
     AND name = 'Derivación a equipos SALVIA' AND deleted_at IS NULL;
  UPDATE salvia.form_section SET "order" = 8
   WHERE form_id IN (v_form_pa, v_form_seg, v_form_cie)
     AND name = 'Cierre' AND deleted_at IS NULL;

  -- ── 2. Insertar preguntas con ON CONFLICT (id) DO NOTHING ───────────────────
  -- Opciones de equipos (reutilizadas en los 4 formularios)

  -- ─── Form PC ─────────────────────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q1, v_form_pc::varchar, v_pc_deriv_s::varchar, 'boolean',
          '¿Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA?',
          TRUE, 1)
  ON CONFLICT (id) DO NOTHING;

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q2, v_form_pc::varchar, v_pc_deriv_s::varchar, 'multiple', '¿Cuáles equipos?', TRUE, 2)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pc_q2, t.label, t.value, t.ord
  FROM (VALUES
    ('Medidas de Emergencia','medidas_emergencia',1),
    ('Estabilización','estabilizacion',2),
    ('Atención Hombres','atencion_hombres',3),
    ('Salvia Dignidad','salvia_dignidad',4),
    ('Discapacidad','discapacidad',5)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pc_q2 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q3, v_form_pc::varchar, v_pc_deriv_s::varchar, 'multiple', '¿Cuáles medidas de emergencia?', FALSE, 3)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pc_q3, t.label, t.value, t.ord
  FROM (VALUES
    ('Medida de protección','medida_proteccion',1),
    ('Carabina / Botón de pánico','boton_panico',2),
    ('Casa de acogida','casa_acogida',3),
    ('Reubicación','reubicacion',4),
    ('Gestión seguridad física en el hogar','gestion_seguridad_fisica',5),
    ('Subsidio de vivienda temporal','subsidio_vivienda',6),
    ('Medida de alejamiento','medida_alejamiento',7),
    ('Custodia y cuidado de hijos e hijas','custodia_hijos',8)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pc_q3 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q4, v_form_pc::varchar, v_pc_deriv_s::varchar, 'multiple', 'Criterios de remisión — Atención Hombres', FALSE, 4)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pc_q4, 'Criterio de remisión', 'criterio_hombres', 1
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pc_q4 AND value = 'criterio_hombres');

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q5, v_form_pc::varchar, v_pc_deriv_s::varchar, 'multiple', 'Criterios de remisión — Salvia Dignidad', FALSE, 5)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pc_q5, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer con orientación sexual diversa (lesbiana, bisexual)','osigd_mujer',1),
    ('Persona no binaria','persona_no_binaria',2)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pc_q5 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q6, v_form_pc::varchar, v_pc_deriv_s::varchar, 'multiple', 'Criterios de remisión — Estabilización', FALSE, 6)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pc_q6, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer víctima de violencia con necesidades de estabilización económica','estabilizacion_economica',1),
    ('Mujer en situación de vulnerabilidad extrema o indigencia','vulnerabilidad_extrema',2),
    ('Mujer desplazada o migrante que requiere apoyo para su estabilización social','desplazada_migrante',3)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pc_q6 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pc_q7, v_form_pc::varchar, v_pc_deriv_s::varchar, 'text',
          'Describa los elementos que evidencia para realizar la remisión', TRUE, 7)
  ON CONFLICT (id) DO NOTHING;

  -- ─── Form PA ─────────────────────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q1, v_form_pa::varchar, v_pa_deriv_s::varchar, 'boolean',
          '¿Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA?',
          TRUE, 1)
  ON CONFLICT (id) DO NOTHING;

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q2, v_form_pa::varchar, v_pa_deriv_s::varchar, 'multiple', '¿Cuáles equipos?', TRUE, 2)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pa_q2, t.label, t.value, t.ord
  FROM (VALUES
    ('Medidas de Emergencia','medidas_emergencia',1),
    ('Estabilización','estabilizacion',2),
    ('Atención Hombres','atencion_hombres',3),
    ('Salvia Dignidad','salvia_dignidad',4),
    ('Discapacidad','discapacidad',5)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pa_q2 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q3, v_form_pa::varchar, v_pa_deriv_s::varchar, 'multiple', '¿Cuáles medidas de emergencia?', FALSE, 3)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pa_q3, t.label, t.value, t.ord
  FROM (VALUES
    ('Medida de protección','medida_proteccion',1),
    ('Carabina / Botón de pánico','boton_panico',2),
    ('Casa de acogida','casa_acogida',3),
    ('Reubicación','reubicacion',4),
    ('Gestión seguridad física en el hogar','gestion_seguridad_fisica',5),
    ('Subsidio de vivienda temporal','subsidio_vivienda',6),
    ('Medida de alejamiento','medida_alejamiento',7),
    ('Custodia y cuidado de hijos e hijas','custodia_hijos',8)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pa_q3 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q4, v_form_pa::varchar, v_pa_deriv_s::varchar, 'multiple', 'Criterios de remisión — Atención Hombres', FALSE, 4)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pa_q4, 'Criterio de remisión', 'criterio_hombres', 1
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pa_q4 AND value = 'criterio_hombres');

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q5, v_form_pa::varchar, v_pa_deriv_s::varchar, 'multiple', 'Criterios de remisión — Salvia Dignidad', FALSE, 5)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pa_q5, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer con orientación sexual diversa (lesbiana, bisexual)','osigd_mujer',1),
    ('Persona no binaria','persona_no_binaria',2)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pa_q5 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q6, v_form_pa::varchar, v_pa_deriv_s::varchar, 'multiple', 'Criterios de remisión — Estabilización', FALSE, 6)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_pa_q6, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer víctima de violencia con necesidades de estabilización económica','estabilizacion_economica',1),
    ('Mujer en situación de vulnerabilidad extrema o indigencia','vulnerabilidad_extrema',2),
    ('Mujer desplazada o migrante que requiere apoyo para su estabilización social','desplazada_migrante',3)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_pa_q6 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_pa_q7, v_form_pa::varchar, v_pa_deriv_s::varchar, 'text',
          'Describa los elementos que evidencia para realizar la remisión', TRUE, 7)
  ON CONFLICT (id) DO NOTHING;

  -- ─── Form SEG ─────────────────────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q1, v_form_seg::varchar, v_seg_deriv_s::varchar, 'boolean',
          '¿Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA?',
          TRUE, 1)
  ON CONFLICT (id) DO NOTHING;

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q2, v_form_seg::varchar, v_seg_deriv_s::varchar, 'multiple', '¿Cuáles equipos?', TRUE, 2)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_seg_q2, t.label, t.value, t.ord
  FROM (VALUES
    ('Medidas de Emergencia','medidas_emergencia',1),
    ('Estabilización','estabilizacion',2),
    ('Atención Hombres','atencion_hombres',3),
    ('Salvia Dignidad','salvia_dignidad',4),
    ('Discapacidad','discapacidad',5)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_seg_q2 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q3, v_form_seg::varchar, v_seg_deriv_s::varchar, 'multiple', '¿Cuáles medidas de emergencia?', FALSE, 3)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_seg_q3, t.label, t.value, t.ord
  FROM (VALUES
    ('Medida de protección','medida_proteccion',1),
    ('Carabina / Botón de pánico','boton_panico',2),
    ('Casa de acogida','casa_acogida',3),
    ('Reubicación','reubicacion',4),
    ('Gestión seguridad física en el hogar','gestion_seguridad_fisica',5),
    ('Subsidio de vivienda temporal','subsidio_vivienda',6),
    ('Medida de alejamiento','medida_alejamiento',7),
    ('Custodia y cuidado de hijos e hijas','custodia_hijos',8)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_seg_q3 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q4, v_form_seg::varchar, v_seg_deriv_s::varchar, 'multiple', 'Criterios de remisión — Atención Hombres', FALSE, 4)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_seg_q4, 'Criterio de remisión', 'criterio_hombres', 1
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_seg_q4 AND value = 'criterio_hombres');

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q5, v_form_seg::varchar, v_seg_deriv_s::varchar, 'multiple', 'Criterios de remisión — Salvia Dignidad', FALSE, 5)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_seg_q5, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer con orientación sexual diversa (lesbiana, bisexual)','osigd_mujer',1),
    ('Persona no binaria','persona_no_binaria',2)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_seg_q5 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q6, v_form_seg::varchar, v_seg_deriv_s::varchar, 'multiple', 'Criterios de remisión — Estabilización', FALSE, 6)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_seg_q6, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer víctima de violencia con necesidades de estabilización económica','estabilizacion_economica',1),
    ('Mujer en situación de vulnerabilidad extrema o indigencia','vulnerabilidad_extrema',2),
    ('Mujer desplazada o migrante que requiere apoyo para su estabilización social','desplazada_migrante',3)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_seg_q6 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_seg_q7, v_form_seg::varchar, v_seg_deriv_s::varchar, 'text',
          'Describa los elementos que evidencia para realizar la remisión', TRUE, 7)
  ON CONFLICT (id) DO NOTHING;

  -- ─── Form CIE ─────────────────────────────────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q1, v_form_cie::varchar, v_cie_deriv_s::varchar, 'boolean',
          '¿Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA?',
          TRUE, 1)
  ON CONFLICT (id) DO NOTHING;

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q2, v_form_cie::varchar, v_cie_deriv_s::varchar, 'multiple', '¿Cuáles equipos?', TRUE, 2)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_cie_q2, t.label, t.value, t.ord
  FROM (VALUES
    ('Medidas de Emergencia','medidas_emergencia',1),
    ('Estabilización','estabilizacion',2),
    ('Atención Hombres','atencion_hombres',3),
    ('Salvia Dignidad','salvia_dignidad',4),
    ('Discapacidad','discapacidad',5)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_cie_q2 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q3, v_form_cie::varchar, v_cie_deriv_s::varchar, 'multiple', '¿Cuáles medidas de emergencia?', FALSE, 3)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_cie_q3, t.label, t.value, t.ord
  FROM (VALUES
    ('Medida de protección','medida_proteccion',1),
    ('Carabina / Botón de pánico','boton_panico',2),
    ('Casa de acogida','casa_acogida',3),
    ('Reubicación','reubicacion',4),
    ('Gestión seguridad física en el hogar','gestion_seguridad_fisica',5),
    ('Subsidio de vivienda temporal','subsidio_vivienda',6),
    ('Medida de alejamiento','medida_alejamiento',7),
    ('Custodia y cuidado de hijos e hijas','custodia_hijos',8)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_cie_q3 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q4, v_form_cie::varchar, v_cie_deriv_s::varchar, 'multiple', 'Criterios de remisión — Atención Hombres', FALSE, 4)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_cie_q4, 'Criterio de remisión', 'criterio_hombres', 1
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_cie_q4 AND value = 'criterio_hombres');

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q5, v_form_cie::varchar, v_cie_deriv_s::varchar, 'multiple', 'Criterios de remisión — Salvia Dignidad', FALSE, 5)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_cie_q5, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer con orientación sexual diversa (lesbiana, bisexual)','osigd_mujer',1),
    ('Persona no binaria','persona_no_binaria',2)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_cie_q5 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q6, v_form_cie::varchar, v_cie_deriv_s::varchar, 'multiple', 'Criterios de remisión — Estabilización', FALSE, 6)
  ON CONFLICT (id) DO NOTHING;
  INSERT INTO salvia.option (id, question_id, label, value, "order")
  SELECT gen_random_uuid()::varchar, v_cie_q6, t.label, t.value, t.ord
  FROM (VALUES
    ('Mujer víctima de violencia con necesidades de estabilización económica','estabilizacion_economica',1),
    ('Mujer en situación de vulnerabilidad extrema o indigencia','vulnerabilidad_extrema',2),
    ('Mujer desplazada o migrante que requiere apoyo para su estabilización social','desplazada_migrante',3)
  ) AS t(label,value,ord)
  WHERE NOT EXISTS (SELECT 1 FROM salvia.option WHERE question_id = v_cie_q6 AND value = t.value);

  INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
  VALUES (v_cie_q7, v_form_cie::varchar, v_cie_deriv_s::varchar, 'text',
          'Describa los elementos que evidencia para realizar la remisión', TRUE, 7)
  ON CONFLICT (id) DO NOTHING;

  -- ── 3. Visibility conditions de preguntas (idempotente WHERE NOT EXISTS) ─────

  -- PC — Q2 visible si Q1 = true
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pc_q2, v_pc_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_q2 AND trigger_question_id = v_pc_q1 AND trigger_value = 'true');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pc_q3, v_pc_q2, 'medidas_emergencia', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_q3 AND trigger_question_id = v_pc_q2 AND trigger_value = 'medidas_emergencia');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pc_q4, v_pc_q2, 'atencion_hombres', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_q4 AND trigger_question_id = v_pc_q2 AND trigger_value = 'atencion_hombres');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pc_q5, v_pc_q2, 'salvia_dignidad', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_q5 AND trigger_question_id = v_pc_q2 AND trigger_value = 'salvia_dignidad');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pc_q6, v_pc_q2, 'estabilizacion', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_q6 AND trigger_question_id = v_pc_q2 AND trigger_value = 'estabilizacion');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pc_q7, v_pc_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_q7 AND trigger_question_id = v_pc_q1 AND trigger_value = 'true');

  -- PA — Q2..Q7
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pa_q2, v_pa_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_q2 AND trigger_question_id = v_pa_q1 AND trigger_value = 'true');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pa_q3, v_pa_q2, 'medidas_emergencia', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_q3 AND trigger_question_id = v_pa_q2 AND trigger_value = 'medidas_emergencia');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pa_q4, v_pa_q2, 'atencion_hombres', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_q4 AND trigger_question_id = v_pa_q2 AND trigger_value = 'atencion_hombres');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pa_q5, v_pa_q2, 'salvia_dignidad', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_q5 AND trigger_question_id = v_pa_q2 AND trigger_value = 'salvia_dignidad');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pa_q6, v_pa_q2, 'estabilizacion', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_q6 AND trigger_question_id = v_pa_q2 AND trigger_value = 'estabilizacion');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_pa_q7, v_pa_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_q7 AND trigger_question_id = v_pa_q1 AND trigger_value = 'true');

  -- SEG — Q2..Q7
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_seg_q2, v_seg_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_q2 AND trigger_question_id = v_seg_q1 AND trigger_value = 'true');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_seg_q3, v_seg_q2, 'medidas_emergencia', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_q3 AND trigger_question_id = v_seg_q2 AND trigger_value = 'medidas_emergencia');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_seg_q4, v_seg_q2, 'atencion_hombres', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_q4 AND trigger_question_id = v_seg_q2 AND trigger_value = 'atencion_hombres');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_seg_q5, v_seg_q2, 'salvia_dignidad', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_q5 AND trigger_question_id = v_seg_q2 AND trigger_value = 'salvia_dignidad');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_seg_q6, v_seg_q2, 'estabilizacion', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_q6 AND trigger_question_id = v_seg_q2 AND trigger_value = 'estabilizacion');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_seg_q7, v_seg_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_q7 AND trigger_question_id = v_seg_q1 AND trigger_value = 'true');

  -- CIE — Q2..Q7
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_cie_q2, v_cie_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_q2 AND trigger_question_id = v_cie_q1 AND trigger_value = 'true');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_cie_q3, v_cie_q2, 'medidas_emergencia', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_q3 AND trigger_question_id = v_cie_q2 AND trigger_value = 'medidas_emergencia');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_cie_q4, v_cie_q2, 'atencion_hombres', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_q4 AND trigger_question_id = v_cie_q2 AND trigger_value = 'atencion_hombres');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_cie_q5, v_cie_q2, 'salvia_dignidad', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_q5 AND trigger_question_id = v_cie_q2 AND trigger_value = 'salvia_dignidad');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_cie_q6, v_cie_q2, 'estabilizacion', 'CONTAINS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_q6 AND trigger_question_id = v_cie_q2 AND trigger_value = 'estabilizacion');

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'QUESTION', v_cie_q7, v_cie_q1, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_q7 AND trigger_question_id = v_cie_q1 AND trigger_value = 'true');

  -- ── 4. Visibility conditions de secciones ────────────────────────────────────

  -- PC  S5 "Derivación" visible cuando Q12 "Continuar Primera Atención" = true
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'SECTION', v_pc_deriv_s::varchar, v_pc_trigger, 'true', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pc_deriv_s::varchar AND trigger_question_id = v_pc_trigger AND trigger_value = 'true');

  -- PA  S5 "Derivación" visible cuando Q8 "¿Es atención o solo contacto?" = atencion
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'SECTION', v_pa_deriv_s::varchar, v_pa_trigger, 'atencion', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_pa_deriv_s::varchar AND trigger_question_id = v_pa_trigger AND trigger_value = 'atencion');

  -- SEG S5 "Derivación" visible cuando Q9 "¿Es atención o solo contacto?" = atencion
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'SECTION', v_seg_deriv_s::varchar, v_seg_trigger, 'atencion', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_seg_deriv_s::varchar AND trigger_question_id = v_seg_trigger AND trigger_value = 'atencion');

  -- CIE S6 "Derivación" visible cuando Q9 "¿Es atención o solo contacto?" = atencion
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  SELECT gen_random_uuid()::varchar, 'SECTION', v_cie_deriv_s::varchar, v_cie_trigger, 'atencion', 'EQUALS'
  WHERE NOT EXISTS (SELECT 1 FROM salvia.visibility_condition WHERE target_id = v_cie_deriv_s::varchar AND trigger_question_id = v_cie_trigger AND trigger_value = 'atencion');

  RAISE NOTICE '=== MIGRATE Sep 2026: Derivación a equipos SALVIA ===';
  RAISE NOTICE '  PC  sección: %  Q1: %', v_pc_deriv_s, v_pc_q1;
  RAISE NOTICE '  PA  sección: %  Q1: %', v_pa_deriv_s, v_pa_q1;
  RAISE NOTICE '  SEG sección: %  Q1: %', v_seg_deriv_s, v_seg_q1;
  RAISE NOTICE '  CIE sección: %  Q1: %', v_cie_deriv_s, v_cie_q1;

END $$;

COMMIT;
