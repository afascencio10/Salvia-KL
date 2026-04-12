-- =============================================================================
-- SEED: Formulario de Seguimiento
-- Estructura completa: secciones, preguntas, opciones, repeater de barreras
-- y condiciones de visibilidad.
-- Ejecutar contra la BD de Salvia:
--   psql "$DATABASE_URL" -f seed_seguimiento.sql
-- =============================================================================

DO $$
DECLARE
  -- ── IDs principales ────────────────────────────────────────────────────────
  v_form_id   uuid;

  -- ── Secciones ──────────────────────────────────────────────────────────────
  v_s1 uuid; -- Inicio de Seguimiento
  v_s2 uuid; -- Valoración del Riesgo
  v_s3 uuid; -- Seguimiento de Caso
  v_s4 uuid; -- Identificación de Barreras
  v_s5 uuid; -- Empalme de Seguimiento

  -- ── Preguntas Sección 1 ───────────────────────────────────────────────────
  v_q_nivel_riesgo        uuid;
  v_q_equipo_atencion     uuid;
  v_q_nombre_profesional  uuid;
  v_q_efectividad         uuid;
  v_q_numero_seguimiento  uuid;
  v_q_fecha_seguimiento   uuid;
  v_q_hora_inicio         uuid;
  v_q_codigo_caso         uuid;

  -- ── Preguntas Sección 2 ───────────────────────────────────────────────────
  v_q_respondiente  uuid;
  v_q_nuevos_hechos uuid;
  v_q_desc_hechos   uuid; -- visible si nuevos_hechos = 'true'

  -- ── Preguntas Sección 3 ───────────────────────────────────────────────────
  v_q_info_psico       uuid;
  v_q_gestion          uuid;
  v_q_variacion        uuid;
  v_q_desc_variacion   uuid; -- visible si variacion = 'true'
  v_q_necesidades      uuid;
  v_q_equipos_der      uuid; -- visible si necesidades = 'true'
  v_q_medidas_em       uuid; -- visible si equipos_der includes 'Medidas de Emergencia'

  -- ── Sección 4: Repeater de Barreras ───────────────────────────────────────
  v_rg_barreras        uuid; -- RepeaterGroup
  v_q_sector           uuid; -- pregunta dentro del repeater
  v_q_barreras_sector  uuid; -- pregunta dentro del repeater

  -- ── Preguntas Sección 5 ───────────────────────────────────────────────────
  v_q_empalme        uuid;
  v_q_riesgo_inm     uuid; -- visible si empalme = 'true'
  v_q_fecha_accion   uuid; -- visible si empalme = 'true'
  v_q_equipo_empalme uuid; -- visible si empalme = 'true'
  v_q_prof_empalme   uuid; -- visible si empalme = 'true'

BEGIN

-- =============================================================================
-- FORM
-- =============================================================================
INSERT INTO salvia.form (id, name, description, status)
VALUES (gen_random_uuid(), 'Formulario de Seguimiento',
        'Registro de seguimientos de casos activos en la plataforma Salvia', 'active')
RETURNING id INTO v_form_id;

RAISE NOTICE 'Form creado: %', v_form_id;

-- =============================================================================
-- SECCIONES
-- =============================================================================
INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
  (gen_random_uuid(), v_form_id, 'Inicio de Seguimiento',
   'Información básica sobre cómo se realizó el seguimiento', 1)
RETURNING id INTO v_s1;

INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
  (gen_random_uuid(), v_form_id, 'Valoración del Riesgo',
   'Evaluación de nuevos hechos y respondiente del seguimiento', 2)
RETURNING id INTO v_s2;

INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
  (gen_random_uuid(), v_form_id, 'Seguimiento de Caso',
   'Gestión realizada y variaciones identificadas en el caso', 3)
RETURNING id INTO v_s3;

INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
  (gen_random_uuid(), v_form_id, 'Identificación de Barreras',
   'Registro de barreras institucionales identificadas durante el seguimiento', 4)
RETURNING id INTO v_s4;

INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
  (gen_random_uuid(), v_form_id, 'Empalme de Seguimiento',
   'Registro de empalme y riesgo inminente', 5)
RETURNING id INTO v_s5;

-- =============================================================================
-- SECCIÓN 1 — Inicio de Seguimiento
-- =============================================================================

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'single',
        '¿Cuál fue el nivel de riesgo identificado en el registro?', 1)
RETURNING id INTO v_q_nivel_riesgo;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_nivel_riesgo, 'Extremo',    'extremo',    1),
  (gen_random_uuid(), v_q_nivel_riesgo, 'Alto',       'alto',       2),
  (gen_random_uuid(), v_q_nivel_riesgo, 'Moderado',   'moderado',   3),
  (gen_random_uuid(), v_q_nivel_riesgo, 'Bajo',       'bajo',       4),
  (gen_random_uuid(), v_q_nivel_riesgo, 'Sin Riesgo', 'sin_riesgo', 5);

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'single',
        'Equipo que realiza la atención', 2)
RETURNING id INTO v_q_equipo_atencion;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_equipo_atencion, 'Agente Integral',          'agente_integral',   1),
  (gen_random_uuid(), v_q_equipo_atencion, 'Riesgo de Feminicidio',    'riesgo_feminicidio',2),
  (gen_random_uuid(), v_q_equipo_atencion, 'Seguimiento General',      'seguimiento_gral',  3),
  (gen_random_uuid(), v_q_equipo_atencion, 'Notificación Salvia',      'notif_salvia',      4),
  (gen_random_uuid(), v_q_equipo_atencion, 'Psicosocial',              'psicosocial',       5),
  (gen_random_uuid(), v_q_equipo_atencion, 'Masculinidades',           'masculinidades',    6);

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'dropdown',
        'Nombre del profesional que realiza la atención', 3)
RETURNING id INTO v_q_nombre_profesional;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_nombre_profesional, 'Ana María Mojica Quiroz',               'ana_mojica',       1),
  (gen_random_uuid(), v_q_nombre_profesional, 'Andres Eduardo Barbosa Deaquiz',        'andres_barbosa',   2),
  (gen_random_uuid(), v_q_nombre_profesional, 'Andres Felipe Suarez Cabra',            'andres_suarez',    3),
  (gen_random_uuid(), v_q_nombre_profesional, 'Cristian Alexander Rodriguez Alarcon',  'cristian_rod',     4),
  (gen_random_uuid(), v_q_nombre_profesional, 'Daniel Esteban Acosta Rodríguez',       'daniel_acosta',    5),
  (gen_random_uuid(), v_q_nombre_profesional, 'Dayan Vargas Sánchez',                  'dayan_vargas',     6),
  (gen_random_uuid(), v_q_nombre_profesional, 'José Felipe Calixto',                   'jose_calixto',     7),
  (gen_random_uuid(), v_q_nombre_profesional, 'Luz Virginia Gomez Morales',            'luz_gomez',        8),
  (gen_random_uuid(), v_q_nombre_profesional, 'Tatiana Geraldine Montalván Caicedo',   'tatiana_mont',     9);

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'single',
        'Efectividad de la llamada', 4)
RETURNING id INTO v_q_efectividad;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_efectividad, 'Llamada Efectiva - Se logra comunicación',         'efectiva',     1),
  (gen_random_uuid(), v_q_efectividad, 'Llamada NO efectiva - NO hay comunicación',         'no_efectiva',  2);

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'single',
        '¿Qué número de seguimiento está registrando?', 5)
RETURNING id INTO v_q_numero_seguimiento;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_numero_seguimiento, '1',  '1',  1),
  (gen_random_uuid(), v_q_numero_seguimiento, '2',  '2',  2),
  (gen_random_uuid(), v_q_numero_seguimiento, '3',  '3',  3),
  (gen_random_uuid(), v_q_numero_seguimiento, '4',  '4',  4),
  (gen_random_uuid(), v_q_numero_seguimiento, '5',  '5',  5),
  (gen_random_uuid(), v_q_numero_seguimiento, '6',  '6',  6),
  (gen_random_uuid(), v_q_numero_seguimiento, '7',  '7',  7),
  (gen_random_uuid(), v_q_numero_seguimiento, '8',  '8',  8),
  (gen_random_uuid(), v_q_numero_seguimiento, '9',  '9',  9),
  (gen_random_uuid(), v_q_numero_seguimiento, '10', '10', 10);

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'date',
        'Fecha del seguimiento', 6)
RETURNING id INTO v_q_fecha_seguimiento;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'datetime',
        'Fecha y hora de inicio del seguimiento', 7)
RETURNING id INTO v_q_hora_inicio;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s1, 'text',
        'Código interno del caso', 8)
RETURNING id INTO v_q_codigo_caso;

-- =============================================================================
-- SECCIÓN 2 — Valoración del Riesgo
-- =============================================================================

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s2, 'single',
        'Respondiente del seguimiento', 1)
RETURNING id INTO v_q_respondiente;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_respondiente, 'Contacto con la víctima',                        'victima',     1),
  (gen_random_uuid(), v_q_respondiente, 'Contacto indirecto - Familiar o persona conocida','indirecto',   2),
  (gen_random_uuid(), v_q_respondiente, 'Contacto Con Institución',                        'institucion', 3);

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s2, 'boolean',
        '¿Se han presentado nuevos hechos de violencia desde el último seguimiento?', 2)
RETURNING id INTO v_q_nuevos_hechos;

-- Visible solo si nuevos_hechos = 'true'
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s2, 'text',
        'Descripción de nuevos hechos de violencia (Tiempo/Modo/Lugar)', 3)
RETURNING id INTO v_q_desc_hechos;

INSERT INTO salvia.visibility_condition
  (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES
  (gen_random_uuid(), 'QUESTION', v_q_desc_hechos::varchar,
   v_q_nuevos_hechos::varchar, 'true', 'EQUALS');

-- =============================================================================
-- SECCIÓN 3 — Seguimiento de Caso
-- =============================================================================

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'text',
        'Información psicosocial relevante para el seguimiento', 1)
RETURNING id INTO v_q_info_psico;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'text',
        'Gestión realizada en el seguimiento', 2)
RETURNING id INTO v_q_gestion;

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'boolean',
        '¿Hubo variación en el nivel de riesgo desde el último seguimiento?', 3)
RETURNING id INTO v_q_variacion;

-- Visible solo si variacion = 'true'
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'text',
        'Describa la variación del riesgo', 4)
RETURNING id INTO v_q_desc_variacion;

INSERT INTO salvia.visibility_condition
  (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES
  (gen_random_uuid(), 'QUESTION', v_q_desc_variacion::varchar,
   v_q_variacion::varchar, 'true', 'EQUALS');

-- ──────────────────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'boolean',
        'Se generaron necesidades inmediatas que requieran derivación a los equipos Salvia', 5)
RETURNING id INTO v_q_necesidades;

-- Visible solo si necesidades = 'true'
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'multiple',
        '¿Cuáles equipos?', 6)
RETURNING id INTO v_q_equipos_der;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_equipos_der, 'Medidas de Emergencia',  'medidas_emergencia', 1),
  (gen_random_uuid(), v_q_equipos_der, 'Atención Psicosocial',   'atencion_psico',     2),
  (gen_random_uuid(), v_q_equipos_der, 'Estabilización',         'estabilizacion',     3);

INSERT INTO salvia.visibility_condition
  (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES
  (gen_random_uuid(), 'QUESTION', v_q_equipos_der::varchar,
   v_q_necesidades::varchar, 'true', 'EQUALS');

-- Visible solo si equipos_der includes 'Medidas de Emergencia'
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s3, 'multiple',
        '¿Cuáles medidas de emergencia?', 7)
RETURNING id INTO v_q_medidas_em;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_medidas_em, 'Alojamiento',        'alojamiento',    1),
  (gen_random_uuid(), v_q_medidas_em, 'Transporte',         'transporte',     2),
  (gen_random_uuid(), v_q_medidas_em, 'Alimentación',       'alimentacion',   3),
  (gen_random_uuid(), v_q_medidas_em, 'Vestuario',          'vestuario',      4),
  (gen_random_uuid(), v_q_medidas_em, 'Apoyo psicosocial',  'apoyo_psico',    5),
  (gen_random_uuid(), v_q_medidas_em, 'Otras ME',           'otras_me',       6);

INSERT INTO salvia.visibility_condition
  (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES
  (gen_random_uuid(), 'QUESTION', v_q_medidas_em::varchar,
   v_q_equipos_der::varchar, 'medidas_emergencia', 'CONTAINS');

-- =============================================================================
-- SECCIÓN 4 — Identificación de Barreras (Repeater)
-- =============================================================================

-- RepeaterGroup
INSERT INTO salvia.repeater_group
  (id, form_section_id, name, item_name, "order", min_repetitions)
VALUES
  (gen_random_uuid(), v_s4, 'Barreras identificadas', 'Barrera', 1, 0)
RETURNING id INTO v_rg_barreras;

-- Pregunta: sector (dentro del repeater)
INSERT INTO salvia.question
  (id, form_id, form_section_id, repeater_group_id, question_type, description, "order")
VALUES
  (gen_random_uuid(), v_form_id, v_s4, v_rg_barreras::varchar, 'dropdown',
   'Sector de la barrera', 1)
RETURNING id INTO v_q_sector;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_sector, 'Salud',      'salud',      1),
  (gen_random_uuid(), v_q_sector, 'Justicia',   'justicia',   2),
  (gen_random_uuid(), v_q_sector, 'Protección', 'proteccion', 3);

-- Pregunta: barreras_sector (dentro del repeater, multiple)
INSERT INTO salvia.question
  (id, form_id, form_section_id, repeater_group_id, question_type, description, "order")
VALUES
  (gen_random_uuid(), v_form_id, v_s4, v_rg_barreras::varchar, 'multiple',
   'Barreras identificadas en este sector', 2)
RETURNING id INTO v_q_barreras_sector;

-- Opciones Salud
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] No Aplica',                                                     'salud_no_aplica',    1),
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] Demoras en asignación de citas y continuidad de tratamientos',  'salud_demoras',      2),
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] Dificultades de aseguramiento o afiliación en salud',           'salud_aseguramiento',3),
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] Falta de activación de protocolos para violencia sexual',       'salud_protocolos',   4),
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] Negación en los servicios de urgencias',                        'salud_urgencias',    5),
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] Negativa en el acceso a la IVE',                                'salud_ive',          6),
  (gen_random_uuid(), v_q_barreras_sector, '[Salud] Otras barreras en salud',                                       'salud_otras',        7);

-- Opciones Justicia
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] No Aplica',                                                  'justicia_no_aplica', 8),
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] Ausencia de representación judicial',                        'justicia_repres',    9),
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] Cargas probatorias injustificadas o excesivas',              'justicia_probat',    10),
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] Falta de celeridad en la investigación',                     'justicia_celeridad', 11),
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] Negativa para recibir la denuncia',                          'justicia_denuncia',  12),
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] Tipificación errónea del delito',                            'justicia_tipif',     13),
  (gen_random_uuid(), v_q_barreras_sector, '[Justicia] Otras barreras en justicia',                                 'justicia_otras',     14);

-- Opciones Protección
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_barreras_sector, '[Protección] No Aplica',                                                  'prot_no_aplica', 15),
  (gen_random_uuid(), v_q_barreras_sector, '[Protección] Demoras en la emisión de medidas de protección urgentes',    'prot_demoras',   16),
  (gen_random_uuid(), v_q_barreras_sector, '[Protección] Fallas en la valoración del riesgo',                         'prot_valoracion',17),
  (gen_random_uuid(), v_q_barreras_sector, '[Protección] Incumplimiento de medidas de protección sin consecuencias',  'prot_incumplim', 18),
  (gen_random_uuid(), v_q_barreras_sector, '[Protección] Medidas de protección ineficaces',                           'prot_ineficaces',19),
  (gen_random_uuid(), v_q_barreras_sector, '[Protección] Otras barreras en protección',                               'prot_otras',     20);

-- =============================================================================
-- SECCIÓN 5 — Empalme de Seguimiento
-- =============================================================================

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s5, 'boolean',
        'Se realiza empalme del seguimiento', 1)
RETURNING id INTO v_q_empalme;

-- Las 4 preguntas siguientes son visibles solo si empalme = 'true'
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s5, 'single',
        '¿Este seguimiento es de riesgo inminente?', 2)
RETURNING id INTO v_q_riesgo_inm;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_riesgo_inm, 'No',                                         'no',          1),
  (gen_random_uuid(), v_q_riesgo_inm, 'Sí. Requiere seguimiento en 4 horas',        'si_4h',       2),
  (gen_random_uuid(), v_q_riesgo_inm, 'Sí. Requiere seguimiento en 8 horas',        'si_8h',       3);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s5, 'datetime',
        'Fecha y hora de la acción a realizar', 3)
RETURNING id INTO v_q_fecha_accion;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s5, 'single',
        'Equipo que realiza la atención', 4)
RETURNING id INTO v_q_equipo_empalme;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_equipo_empalme, 'Riesgo de Feminicidio',  'riesgo_feminicidio', 1),
  (gen_random_uuid(), v_q_equipo_empalme, 'Seguimiento General',    'seguimiento_gral',   2),
  (gen_random_uuid(), v_q_equipo_empalme, 'Agente Integral',        'agente_integral',    3);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, "order")
VALUES (gen_random_uuid(), v_form_id, v_s5, 'dropdown',
        'Nombre del profesional para empalme', 5)
RETURNING id INTO v_q_prof_empalme;

INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_q_prof_empalme, 'Ana María Mojica Quiroz',               'ana_mojica',    1),
  (gen_random_uuid(), v_q_prof_empalme, 'Andres Eduardo Barbosa Deaquiz',        'andres_barbosa',2),
  (gen_random_uuid(), v_q_prof_empalme, 'Andres Felipe Suarez Cabra',            'andres_suarez', 3),
  (gen_random_uuid(), v_q_prof_empalme, 'Cristian Alexander Rodriguez Alarcon',  'cristian_rod',  4),
  (gen_random_uuid(), v_q_prof_empalme, 'Daniel Esteban Acosta Rodríguez',       'daniel_acosta', 5),
  (gen_random_uuid(), v_q_prof_empalme, 'Dayan Vargas Sánchez',                  'dayan_vargas',  6),
  (gen_random_uuid(), v_q_prof_empalme, 'José Felipe Calixto',                   'jose_calixto',  7),
  (gen_random_uuid(), v_q_prof_empalme, 'Luz Virginia Gomez Morales',            'luz_gomez',     8),
  (gen_random_uuid(), v_q_prof_empalme, 'Tatiana Geraldine Montalván Caicedo',   'tatiana_mont',  9);

-- Condiciones de visibilidad para Sección 5 (trigger: empalme = 'true')
INSERT INTO salvia.visibility_condition
  (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES
  (gen_random_uuid(), 'QUESTION', v_q_riesgo_inm::varchar,    v_q_empalme::varchar, 'true', 'EQUALS'),
  (gen_random_uuid(), 'QUESTION', v_q_fecha_accion::varchar,  v_q_empalme::varchar, 'true', 'EQUALS'),
  (gen_random_uuid(), 'QUESTION', v_q_equipo_empalme::varchar,v_q_empalme::varchar, 'true', 'EQUALS'),
  (gen_random_uuid(), 'QUESTION', v_q_prof_empalme::varchar,  v_q_empalme::varchar, 'true', 'EQUALS');

-- =============================================================================
RAISE NOTICE 'Seed completado. Form ID: %', v_form_id;
RAISE NOTICE 'Sección 1: % | Sección 2: % | Sección 3: % | Sección 4: % | Sección 5: %',
  v_s1, v_s2, v_s3, v_s4, v_s5;

END $$;
