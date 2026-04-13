-- =============================================================================
-- SEED: 4 FormSubmissions de prueba
-- Sub 1: Sin respuestas
-- Sub 2: Respondido S1
-- Sub 3: Respondido S1 y S2
-- Sub 4: Respondido S1, S2, S3 y S4 (con repeater de barreras)
-- =============================================================================

DO $$
DECLARE
  v_form_id varchar := '2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff';
  v_rg_id   varchar := '5fd3ecdc-2e5f-4b31-97ef-8a994580586a';

  -- Submissions
  v_sub1 varchar;
  v_sub2 varchar;
  v_sub3 varchar;
  v_sub4 varchar;

  -- RepeaterEntry para sub4
  v_entry1 varchar;

  -- Question IDs — Sección 1
  q_nivel_riesgo       varchar := 'a1937f54-9282-48bb-bf6b-2a8bd9933c85';
  q_equipo_atencion    varchar := '3cc676e1-dbb6-437e-ad06-784d0e89600e';
  q_nombre_prof        varchar := '6c4429c6-b828-456e-adc1-5556557ca117';
  q_efectividad        varchar := '813b67ff-18e6-494d-899d-f7551b453673';
  q_num_seguimiento    varchar := '8f5aed6d-56e7-4a51-8db8-12cfe248e4ed';
  q_fecha_seg          varchar := '491dc915-6fab-432e-b3e6-0c03e83e3afd';
  q_hora_inicio        varchar := 'f567a718-5515-414d-a9b7-91d7ce575f89';
  q_codigo_caso        varchar := '81d7b17f-3438-4d07-b72f-b45f858d56c4';

  -- Question IDs — Sección 2
  q_respondiente       varchar := 'f8b69cd8-c684-4ec2-b1fb-3853c9cfffed';
  q_nuevos_hechos      varchar := 'da90b1c9-45a1-4e2a-83dd-a2ae77cf5227';
  q_desc_hechos        varchar := 'f8453544-2a8c-461d-9986-302ee4719492';

  -- Question IDs — Sección 3
  q_info_psico         varchar := '6891fbec-da1d-46e6-9c3b-3bb824d19631';
  q_gestion            varchar := 'bb7a2307-e461-47b4-b34e-401ca807efe7';
  q_variacion          varchar := 'be3e7cde-426c-4fde-8cd9-a7453566ead7';
  q_desc_variacion     varchar := '4433a431-dbef-42d7-bc36-a17374588902';
  q_necesidades        varchar := 'fc64033e-1056-4771-9134-d054ebca8d6f';
  q_equipos_der        varchar := 'e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc';
  q_medidas_em         varchar := '1a36260c-33a4-4ebd-bffb-e387d7964b96';

  -- Question IDs — Sección 4 (repeater)
  q_sector             varchar := 'f19378b6-55c5-4fdf-b765-7ebcc3978741';
  q_barreras_salud     varchar := '5fc1f2af-cc30-41f0-aa31-4731e5cb674c';

BEGIN

-- =============================================================================
-- SUBMISSION 1 — Sin respuestas
-- =============================================================================
INSERT INTO salvia.form_submission (id, form_id)
VALUES (gen_random_uuid()::varchar, v_form_id)
RETURNING id INTO v_sub1;

RAISE NOTICE 'Sub1 (sin respuestas): %', v_sub1;

-- =============================================================================
-- SUBMISSION 2 — Respondido S1
-- =============================================================================
INSERT INTO salvia.form_submission (id, form_id)
VALUES (gen_random_uuid()::varchar, v_form_id)
RETURNING id INTO v_sub2;

INSERT INTO salvia.answer (id, form_submission_id, question_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub2, q_nivel_riesgo,    'alto'),
  (gen_random_uuid()::varchar, v_sub2, q_equipo_atencion, 'agente_integral'),
  (gen_random_uuid()::varchar, v_sub2, q_nombre_prof,     'daniel_acosta'),
  (gen_random_uuid()::varchar, v_sub2, q_efectividad,     'efectiva'),
  (gen_random_uuid()::varchar, v_sub2, q_num_seguimiento, '1'),
  (gen_random_uuid()::varchar, v_sub2, q_fecha_seg,       '2026-04-10'),
  (gen_random_uuid()::varchar, v_sub2, q_hora_inicio,     '2026-04-10T09:00'),
  (gen_random_uuid()::varchar, v_sub2, q_codigo_caso,     'SAL-2026-0087');

RAISE NOTICE 'Sub2 (S1): %', v_sub2;

-- =============================================================================
-- SUBMISSION 3 — Respondido S1 y S2
-- =============================================================================
INSERT INTO salvia.form_submission (id, form_id)
VALUES (gen_random_uuid()::varchar, v_form_id)
RETURNING id INTO v_sub3;

-- S1
INSERT INTO salvia.answer (id, form_submission_id, question_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub3, q_nivel_riesgo,    'extremo'),
  (gen_random_uuid()::varchar, v_sub3, q_equipo_atencion, 'riesgo_feminicidio'),
  (gen_random_uuid()::varchar, v_sub3, q_nombre_prof,     'tatiana_mont'),
  (gen_random_uuid()::varchar, v_sub3, q_efectividad,     'efectiva'),
  (gen_random_uuid()::varchar, v_sub3, q_num_seguimiento, '3'),
  (gen_random_uuid()::varchar, v_sub3, q_fecha_seg,       '2026-04-09'),
  (gen_random_uuid()::varchar, v_sub3, q_hora_inicio,     '2026-04-09T14:30'),
  (gen_random_uuid()::varchar, v_sub3, q_codigo_caso,     'SAL-2026-0091');

-- S2 (nuevos_hechos=true → incluye descripcion)
INSERT INTO salvia.answer (id, form_submission_id, question_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub3, q_respondiente,  'victima'),
  (gen_random_uuid()::varchar, v_sub3, q_nuevos_hechos, 'true'),
  (gen_random_uuid()::varchar, v_sub3, q_desc_hechos,   'El día 08/04 la víctima reportó amenazas verbales por parte del agresor en horas de la noche en el domicilio.');

RAISE NOTICE 'Sub3 (S1+S2): %', v_sub3;

-- =============================================================================
-- SUBMISSION 4 — Respondido S1, S2, S3 y S4 (con repeater)
-- =============================================================================
INSERT INTO salvia.form_submission (id, form_id)
VALUES (gen_random_uuid()::varchar, v_form_id)
RETURNING id INTO v_sub4;

-- S1
INSERT INTO salvia.answer (id, form_submission_id, question_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub4, q_nivel_riesgo,    'moderado'),
  (gen_random_uuid()::varchar, v_sub4, q_equipo_atencion, 'psicosocial'),
  (gen_random_uuid()::varchar, v_sub4, q_nombre_prof,     'luz_gomez'),
  (gen_random_uuid()::varchar, v_sub4, q_efectividad,     'no_efectiva'),
  (gen_random_uuid()::varchar, v_sub4, q_num_seguimiento, '2'),
  (gen_random_uuid()::varchar, v_sub4, q_fecha_seg,       '2026-04-08'),
  (gen_random_uuid()::varchar, v_sub4, q_hora_inicio,     '2026-04-08T11:15'),
  (gen_random_uuid()::varchar, v_sub4, q_codigo_caso,     'SAL-2026-0075');

-- S2 (sin nuevos hechos)
INSERT INTO salvia.answer (id, form_submission_id, question_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub4, q_respondiente,  'indirecto'),
  (gen_random_uuid()::varchar, v_sub4, q_nuevos_hechos, 'false');

-- S3 (variacion=true, necesidades=true con medidas de emergencia)
INSERT INTO salvia.answer (id, form_submission_id, question_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub4, q_info_psico,     'La víctima presenta signos de ansiedad elevada. Refiere dificultades para dormir y miedo constante.'),
  (gen_random_uuid()::varchar, v_sub4, q_gestion,        'Se realizó seguimiento telefónico. Se activó ruta de atención psicosocial y se coordinó con equipo de medidas de emergencia.'),
  (gen_random_uuid()::varchar, v_sub4, q_variacion,      'true'),
  (gen_random_uuid()::varchar, v_sub4, q_desc_variacion, 'El nivel de riesgo aumentó de Bajo a Moderado. La víctima reporta que el agresor ha retomado contacto.'),
  (gen_random_uuid()::varchar, v_sub4, q_necesidades,    'true'),
  (gen_random_uuid()::varchar, v_sub4, q_equipos_der,    'medidas_emergencia,atencion_psico'),
  (gen_random_uuid()::varchar, v_sub4, q_medidas_em,     'alojamiento,transporte');

-- S4 — Repeater: 1 entrada con sector Salud
INSERT INTO salvia.repeater_entry (id, form_submission_id, repeater_group_id, iteration)
VALUES (gen_random_uuid()::varchar, v_sub4, v_rg_id, 1)
RETURNING id INTO v_entry1;

INSERT INTO salvia.answer (id, form_submission_id, question_id, repeater_entry_id, value) VALUES
  (gen_random_uuid()::varchar, v_sub4, q_sector,        v_entry1, 'salud'),
  (gen_random_uuid()::varchar, v_sub4, q_barreras_salud, v_entry1, 'salud_urgencias,salud_protocolos');

RAISE NOTICE 'Sub4 (S1+S2+S3+S4): % | entry: %', v_sub4, v_entry1;

END $$;
