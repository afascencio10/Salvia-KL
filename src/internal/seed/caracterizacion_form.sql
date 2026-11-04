-- Formulario dinámico "Entrevista de Caracterización" (Medidas de Emergencia).
-- Id fijo: f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d
-- FASE 2b-1: solo siembra la definición del formulario (secciones, preguntas,
-- opciones, repeaters y condicionales). NO precarga answers, NO enlaza el
-- submission a la tarea/caso y NO crea medidas al completar (eso es FASE 2b-2).
-- No está registrado en OnEndFormSubmission, así que completarlo solo guarda
-- respuestas (igual que Seguimiento V.4 en su día).
--
-- Una sola sentencia (bloque DO) para db.Exec. UUID fijos (no gen_random_uuid()).
-- Si el formulario ya tiene envíos, aborta y no borra nada.
-- Catálogos cerrados EPS y Régimen SGSSS: sembrados como option inline (ver doc
-- de mapeo §4, decisión G-2). Los catálogos que EXISTEN se resuelven en el front
-- vía state_options_path (enums.victim_case_form2_*, departments, geo.residencia.*)
-- y por eso esas preguntas no llevan option inline.
DO $$
DECLARE
  n_envios integer;
  -- i_code de las opciones "Otra/Otro" (uuid_generate_v4 por entorno, no fijo).
  ot_sexual_orientation text;
  ot_occupation text;
BEGIN

SELECT COUNT(*) INTO n_envios FROM salvia.form_submission
  WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d';
IF n_envios > 0 THEN
  RAISE EXCEPTION 'el formulario de entrevista de caracterización ya tiene envíos y no se reemplaza';
END IF;

-- El motor guarda y compara el i_code de la opción (no el código corto 'ot') para
-- los dropdown con state_options_path = enums.*. Ese i_code es un UUID generado por
-- entorno, así que se resuelve aquí contra el catálogo de la BD actual.
SELECT victim_case_form2_enums_i_code::text INTO ot_sexual_orientation
  FROM salvia.victim_case_form2_enums
  WHERE victim_case_form2_enums_category = 'victim_case_form2_sexual_orientation'
    AND victim_case_form2_enums_code = 'ot' LIMIT 1;
SELECT victim_case_form2_enums_i_code::text INTO ot_occupation
  FROM salvia.victim_case_form2_enums
  WHERE victim_case_form2_enums_category = 'victim_case_form2_occupation'
    AND victim_case_form2_enums_code = 'ot' LIMIT 1;
IF ot_sexual_orientation IS NULL OR ot_occupation IS NULL THEN
  RAISE EXCEPTION 'no se encontró la opción ot en victim_case_form2_enums (sexual_orientation / occupation)';
END IF;

-- Limpieza previa (mismo patrón que el seed de Seguimiento V.4).
DELETE FROM salvia.visibility_condition
  WHERE trigger_question_id IN (SELECT id FROM salvia.question WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d')
     OR target_id IN (SELECT id FROM salvia.question WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d')
     OR target_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d')
     OR target_id IN (
          SELECT id FROM salvia.repeater_group
           WHERE form_section_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d')
        );
DELETE FROM salvia.answer
  WHERE question_id IN (SELECT id FROM salvia.question WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d');
DELETE FROM salvia.repeater_entry
  WHERE repeater_group_id IN (
        SELECT id FROM salvia.repeater_group
         WHERE form_section_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d')
      );
DELETE FROM salvia.option
  WHERE question_id IN (SELECT id FROM salvia.question WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d');
DELETE FROM salvia.question WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d';
DELETE FROM salvia.repeater_group
  WHERE form_section_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d');
DELETE FROM salvia.form_section WHERE form_id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d';

INSERT INTO salvia.form (id, name, description, status) VALUES (
  'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d',
  'Entrevista de Caracterización',
  'Entrevista de Caracterización de Medidas de Emergencia. FASE 2b-1 v3 del 2026: siembra la definición del formulario (8 secciones + resultado). No precarga answers, no enlaza el submission a la tarea ni crea medidas al completar.',
  'active'
) ON CONFLICT (id) DO NOTHING;

-- ─── Secciones ──────────────────────────────────────────────────────────────
INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
  ('a0000001-0000-4000-8000-000000000001', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'Preámbulo', 'Fecha de la entrevista de caracterización', 1),
  ('a0000001-0000-4000-8000-000000000002', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'I. Datos generales de la usuaria en riesgo de feminicidio', 'Datos de identificación, ubicación, contacto y caracterización sociodemográfica', 2),
  ('a0000001-0000-4000-8000-000000000003', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'II. Información familiar', 'Conformación familiar, miembros y proveeduría del hogar', 3),
  ('a0000001-0000-4000-8000-000000000004', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'III. Apoyo para alojamiento', 'Convivencia o cercanía del presunto agresor e interés en traslado', 4),
  ('a0000001-0000-4000-8000-000000000005', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'IV. Movilización y transporte', 'Necesidad y tipo de apoyo para transporte', 5),
  ('a0000001-0000-4000-8000-000000000006', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'V. Alimentación', 'Afectación del acceso a la alimentación', 6),
  ('a0000001-0000-4000-8000-000000000007', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'VI. Atención psicojurídica', 'Orientación jurídica, representación y acompañamiento psicosocial', 7),
  ('a0000001-0000-4000-8000-000000000008', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'VII. Vestuario', 'Abandono de vestuario propio o de personas a cargo', 8),
  ('a0000001-0000-4000-8000-000000000009', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'VIII. Otro tipo de medidas de emergencia', 'Otras necesidades relacionadas con la VBG/VpP', 9),
  ('a0000001-0000-4000-8000-000000000010', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'Resultado de la caracterización', 'Hallazgos, plan de acción, beneficiarios, medidas solicitadas, documentos y estado', 10)
ON CONFLICT (id) DO NOTHING;

-- ─── Grupos repetibles ────────────────────────────────────────────────────────
-- RG-1 Miembros de la familia (Sección II). RG-2 Beneficiarios y "Medida
-- solicitada" (Sección Resultado). min_repetitions = 0, sin máximo (decisión usuario).
INSERT INTO salvia.repeater_group (id, form_section_id, name, item_name, "order", min_repetitions, state_items, parent_repeater_group_id) VALUES
  ('b0000001-0000-4000-8000-000000000001', 'a0000001-0000-4000-8000-000000000003', 'Miembros de la familia', 'Miembro', 3, 0, NULL, NULL),
  ('b0000001-0000-4000-8000-000000000002', 'a0000001-0000-4000-8000-000000000010', 'Beneficiarios de las medidas de emergencia', 'Beneficiario', 5, 0, NULL, NULL),
  ('b0000001-0000-4000-8000-000000000003', 'a0000001-0000-4000-8000-000000000010', 'Medida solicitada', 'Medida', 6, 0, NULL, NULL)
ON CONFLICT (id) DO NOTHING;

-- ─── Preguntas ────────────────────────────────────────────────────────────────
INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order") VALUES
  -- Preámbulo
  ('c0000001-0000-4000-8000-000000000001', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000001', NULL, 'date', 'Fecha de la entrevista de caracterización', true, NULL, 1),
  -- Sección I
  ('c0000001-0000-4000-8000-000000000101', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', '¿Cuáles son sus nombres?', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000102', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', '¿Cuáles son sus apellidos?', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000103', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Nombre identitario', false, NULL, 3),
  ('c0000001-0000-4000-8000-000000000104', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', '¿Qué tipo de documento de identificación tiene?', true, NULL, 4),
  ('c0000001-0000-4000-8000-000000000105', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'date', 'Fecha de nacimiento', true, NULL, 5),
  ('c0000001-0000-4000-8000-000000000106', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'number', 'Edad', true, NULL, 6),
  ('c0000001-0000-4000-8000-000000000107', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Régimen de afiliación al SGSSS', true, NULL, 7),
  ('c0000001-0000-4000-8000-000000000108', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Nombre de la EPS', true, NULL, 8),
  ('c0000001-0000-4000-8000-000000000109', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Dirección o indicaciones de su vivienda', true, NULL, 9),
  ('c0000001-0000-4000-8000-000000000110', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Zona en la que reside', true, 'enums.victim_case_form2_facts_zone', 10),
  ('c0000001-0000-4000-8000-000000000111', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Departamento', true, 'departments', 11),
  ('c0000001-0000-4000-8000-000000000112', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Ciudad o Municipio', true, 'geo.residencia.cities', 12),
  ('c0000001-0000-4000-8000-000000000113', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Corregimiento o Vereda', true, NULL, 13),
  ('c0000001-0000-4000-8000-000000000114', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'number', '¿Número de teléfono de contacto en Colombia?', true, NULL, 14),
  ('c0000001-0000-4000-8000-000000000115', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Nombre de persona de contacto de apoyo/emergencia', true, NULL, 15),
  ('c0000001-0000-4000-8000-000000000116', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'number', 'Número de teléfono del contacto de apoyo', true, NULL, 16),
  ('c0000001-0000-4000-8000-000000000117', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Estado civil', true, 'enums.victim_case_form2_marital_status', 17),
  ('c0000001-0000-4000-8000-000000000118', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Sexo', true, 'enums.victim_case_form2_assigned_sex_at_birth', 18),
  ('c0000001-0000-4000-8000-000000000119', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', '¿Cuál es su identidad de género?', true, 'enums.victim_case_form2_gender_identity', 19),
  ('c0000001-0000-4000-8000-000000000120', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Orientación sexual', true, 'enums.victim_case_form2_sexual_orientation', 21),
  ('c0000001-0000-4000-8000-00000000012a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Orientación sexual — Otra ¿cuál?', false, NULL, 22),
  ('c0000001-0000-4000-8000-000000000121', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Pertenencia étnica', true, 'enums.victim_case_form2_ethnic_affiliation', 23),
  ('c0000001-0000-4000-8000-000000000122', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'boolean', '¿Se reconoce como perteneciente a un pueblo indígena?', true, NULL, 24),
  ('c0000001-0000-4000-8000-00000000012b', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', '¿Cuál pueblo indígena?', false, 'enums.victim_case_form2_indigenous_people', 25),
  ('c0000001-0000-4000-8000-000000000123', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'boolean', '¿Hace parte de la población migrante?', true, NULL, 27),
  ('c0000001-0000-4000-8000-000000000124', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Situación migratoria', true, 'enums.victim_case_form2_migration_condition', 28),
  ('c0000001-0000-4000-8000-000000000125', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Nivel educativo más alto alcanzado', true, 'enums.victim_case_form2_last_education_level', 29),
  ('c0000001-0000-4000-8000-000000000126', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'multiple', '¿Pertenece a alguna población especial?', true, 'enums.victim_case_form2_specially_protected_population', 30),
  ('c0000001-0000-4000-8000-000000000127', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'boolean', '¿Se reconoce como persona con discapacidad?', true, NULL, 31),
  ('c0000001-0000-4000-8000-000000000128', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'multiple', '¿Requiere ajuste razonable para la atención del caso VBG/VpP?', true, 'enums.victim_case_form2_adjustments_gbv', 32),
  ('c0000001-0000-4000-8000-000000000129', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', 'Ocupación / actividad principal', true, 'enums.victim_case_form2_occupation', 33),
  ('c0000001-0000-4000-8000-00000000012d', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Ocupación — Otra ¿cuál?', false, NULL, 34),
  ('c0000001-0000-4000-8000-000000000130', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'dropdown', '¿Forma de generación de ingresos actual o últimos 3 meses?', true, 'enums.victim_case_form2_income_generation_method', 35),
  ('c0000001-0000-4000-8000-00000000013a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Describa las circunstancias de generación de ingresos', true, NULL, 36),
  ('c0000001-0000-4000-8000-000000000131', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Hechos victimizantes (se extraen de SALVIA sin modificar)', true, NULL, 37),
  ('c0000001-0000-4000-8000-000000000132', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Información complementaria sobre la situación de riesgo y medidas de emergencia', true, NULL, 38),
  ('c0000001-0000-4000-8000-000000000133', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'boolean', '¿Ha recibido ayuda o subsidio de la Alcaldía por la situación de VBG/VpP?', false, NULL, 39),
  ('c0000001-0000-4000-8000-000000000134', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Describa la ayuda recibida (Alcaldía)', false, NULL, 40),
  ('c0000001-0000-4000-8000-000000000135', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'boolean', '¿Ha recibido ayuda o subsidio de la Secretaría u oficina de la mujer?', false, NULL, 41),
  ('c0000001-0000-4000-8000-000000000136', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Describa la ayuda recibida (Secretaría u oficina de la mujer)', false, NULL, 42),
  ('c0000001-0000-4000-8000-000000000137', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'boolean', '¿Ha recibido ayuda o subsidio de otras entidades estatales?', false, NULL, 43),
  ('c0000001-0000-4000-8000-000000000138', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000002', NULL, 'text', 'Describa la ayuda recibida (otras entidades estatales)', false, NULL, 44),
  -- Sección II — directas de sección
  ('c0000001-0000-4000-8000-000000000201', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', NULL, 'single', 'Conformación familiar', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000203', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', NULL, 'boolean', '¿Tiene personas o animales a cargo?', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000204', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', NULL, 'dropdown', '¿Quiénes están a su cargo?', true, 'enums.victim_case_form2_has_dependents', 4),
  ('c0000001-0000-4000-8000-000000000205', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', NULL, 'single', 'Para cubrir gastos del hogar, ¿usted es la proveedora o depende de otra persona?', true, NULL, 5),
  ('c0000001-0000-4000-8000-00000000020a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', NULL, 'text', 'Describa las razones de su rol de proveedora o de dependencia', true, NULL, 6),
  -- Sección II — preguntas del repeater RG-1 (Miembros de la familia)
  ('c0000001-0000-4000-8000-000000000211', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', 'b0000001-0000-4000-8000-000000000001', 'single', 'Parentesco', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000212', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', 'b0000001-0000-4000-8000-000000000001', 'text', 'Otro integrante ¿cuál?', false, NULL, 2),
  ('c0000001-0000-4000-8000-000000000213', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000003', 'b0000001-0000-4000-8000-000000000001', 'number', 'Edad', true, NULL, 3),
  -- Sección III
  ('c0000001-0000-4000-8000-000000000301', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000004', NULL, 'single', '¿El presunto agresor vive en la misma vivienda o alrededores?', true, NULL, 1),
  ('c0000001-0000-4000-8000-00000000030a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000004', NULL, 'text', 'Describa la convivencia o cercanía con el presunto agresor', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000302', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000004', NULL, 'single', '¿Estaría interesada en cambiar de vivienda, municipio o ciudad?', true, NULL, 3),
  ('c0000001-0000-4000-8000-00000000030b', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000004', NULL, 'text', 'Describa el lugar de traslado o las razones para no trasladarse', true, NULL, 4),
  ('c0000001-0000-4000-8000-000000000303', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000004', NULL, 'single', '¿El presunto agresor conoce la ubicación y ha habido acoso u hostigamiento?', true, NULL, 5),
  ('c0000001-0000-4000-8000-00000000030c', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000004', NULL, 'text', 'Describa el conocimiento de la ubicación o el hostigamiento', true, NULL, 6),
  -- Sección IV
  ('c0000001-0000-4000-8000-000000000401', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000005', NULL, 'boolean', '¿Requiere transporte para movilizarse a otro lugar y no tiene recursos?', true, NULL, 1),
  ('c0000001-0000-4000-8000-00000000040a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000005', NULL, 'text', 'Describa los recursos con los que cuenta para transporte', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000402', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000005', NULL, 'multiple', '¿Requiere apoyo para transporte para…?', true, NULL, 3),
  ('c0000001-0000-4000-8000-000000000403', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000005', NULL, 'single', '¿Qué tipo de transporte requiere?', false, NULL, 4),
  ('c0000001-0000-4000-8000-00000000040b', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000005', NULL, 'text', 'Amplíe el tipo de transporte', false, NULL, 5),
  -- Sección V
  ('c0000001-0000-4000-8000-000000000501', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000006', NULL, 'single', '¿Se ha afectado el acceso a la alimentación para usted y sus dependientes?', true, NULL, 1),
  ('c0000001-0000-4000-8000-00000000050a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000006', NULL, 'text', 'Justifique la respuesta de alimentación', true, NULL, 2),
  -- Sección VI
  ('c0000001-0000-4000-8000-000000000601', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000007', NULL, 'boolean', '¿Está interesada en orientación o asesoría jurídica por entidad pública?', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000602', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000007', NULL, 'boolean', '¿Cuenta con representación de víctimas?', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000603', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000007', NULL, 'boolean', '¿Ha recibido acompañamiento psicosocial por entidad pública por los hechos VBG/VpP?', true, NULL, 3),
  ('c0000001-0000-4000-8000-00000000060a', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000007', NULL, 'text', 'Describa el proceso, las barreras y el estado emocional', true, NULL, 4),
  -- Sección VII
  ('c0000001-0000-4000-8000-000000000701', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000008', NULL, 'boolean', '¿Ha tenido que abandonar su vestuario o el de las personas a cargo?', true, NULL, 1),
  ('c0000001-0000-4000-8000-00000000070b', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000008', NULL, 'text', 'Información adicional de vestuario', false, NULL, 2),
  -- Sección VIII
  ('c0000001-0000-4000-8000-000000000801', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000009', NULL, 'text', '¿Qué otro tipo de necesidades se han presentado en relación con la VBG/VpP?', true, NULL, 1),
  -- Sección Resultado — directas
  ('c0000001-0000-4000-8000-000000000901', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', NULL, 'text', 'Resumen de hallazgos de la entrevista', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000902', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', NULL, 'text', 'Plan de acción', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000903', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', NULL, 'info', 'Nombre del profesional del Equipo de Abordaje Integral del Feminicidio: se toma de la asignación del profesional que gestiona el caso (solo lectura, no se captura en el formulario).', false, NULL, 3),
  ('c0000001-0000-4000-8000-000000000904', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', NULL, 'info', 'Registre a continuación los beneficiarios de las medidas de emergencia (no se repite la beneficiaria principal).', false, NULL, 4),
  ('c0000001-0000-4000-8000-000000000905', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', NULL, 'text', 'Documentos para aprobación por Supervisión ME (URL de la carpeta: Registro SALVIA, formato de remisión a Cruz Roja y declaración juramentada)', true, NULL, 7),
  ('c0000001-0000-4000-8000-000000000906', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', NULL, 'single', 'Estado del caso: Caracterizado', true, NULL, 8),
  -- Sección Resultado — preguntas del repeater RG-2 (Beneficiarios)
  ('c0000001-0000-4000-8000-000000000911', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'text', 'Nombres y apellidos', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000912', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'number', 'Edad', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000913', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'dropdown', 'Tipo de documento', true, NULL, 3),
  ('c0000001-0000-4000-8000-000000000914', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'text', 'Número de documento', true, NULL, 4),
  ('c0000001-0000-4000-8000-000000000915', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'boolean', '¿Presenta discapacidad o diagnóstico en salud física o mental?', false, NULL, 5),
  ('c0000001-0000-4000-8000-000000000916', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'text', '¿Cuál es la discapacidad o diagnóstico?', false, NULL, 6),
  ('c0000001-0000-4000-8000-000000000917', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000002', 'text', 'Parentesco con la titular del caso', true, NULL, 7),
  -- Sección Resultado — preguntas del repeater RG-3 (Medida solicitada)
  ('c0000001-0000-4000-8000-000000000921', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000003', 'dropdown', 'Medida', true, NULL, 1),
  ('c0000001-0000-4000-8000-000000000922', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000003', 'boolean', '¿Se solicita?', true, NULL, 2),
  ('c0000001-0000-4000-8000-000000000923', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000003', 'text', 'Tiempo (días o meses)', false, NULL, 3),
  ('c0000001-0000-4000-8000-000000000924', 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d', 'a0000001-0000-4000-8000-000000000010', 'b0000001-0000-4000-8000-000000000003', 'text', 'Observaciones', false, NULL, 4)
ON CONFLICT (id) DO NOTHING;

-- ─── Opciones (option inline) ───────────────────────────────────────────────
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  -- I.4 Tipo de documento (mapa DOCUMENT_TYPE_FORM2, igual que Registro de Caso)
  ('d0000001-0000-4000-8000-000000000401', 'c0000001-0000-4000-8000-000000000104', 'Carné Diplomático', 'cd', 1),
  ('d0000001-0000-4000-8000-000000000402', 'c0000001-0000-4000-8000-000000000104', 'Cédula de Ciudadanía', 'cc', 2),
  ('d0000001-0000-4000-8000-000000000403', 'c0000001-0000-4000-8000-000000000104', 'Cédula de Extranjería', 'ce', 3),
  ('d0000001-0000-4000-8000-000000000404', 'c0000001-0000-4000-8000-000000000104', 'Certificado de Nacido Vivo del país de origen', 'no', 4),
  ('d0000001-0000-4000-8000-000000000405', 'c0000001-0000-4000-8000-000000000104', 'Certificado de Nacido Vivo en Colombia', 'nc', 5),
  ('d0000001-0000-4000-8000-000000000406', 'c0000001-0000-4000-8000-000000000104', 'Certificado de Registro Administrativo de Migrantes Venezolanos - RAMV', 'mv', 6),
  ('d0000001-0000-4000-8000-000000000407', 'c0000001-0000-4000-8000-000000000104', 'Documento Extranjero o Documento Nacional de Identidad (del país de origen -DNI )', 'dn', 7),
  ('d0000001-0000-4000-8000-000000000408', 'c0000001-0000-4000-8000-000000000104', 'Otro no especificado', 'oe', 8),
  ('d0000001-0000-4000-8000-000000000409', 'c0000001-0000-4000-8000-000000000104', 'Partida de Nacimiento de su país de origen', 'pn', 9),
  ('d0000001-0000-4000-8000-00000000040a', 'c0000001-0000-4000-8000-000000000104', 'Pasaporte', 'ps', 10),
  ('d0000001-0000-4000-8000-00000000040b', 'c0000001-0000-4000-8000-000000000104', 'Permiso Especial de Permanencia (PEP)', 'pp', 11),
  ('d0000001-0000-4000-8000-00000000040c', 'c0000001-0000-4000-8000-000000000104', 'Permiso Especial de Permanencia para el Fomento de la Formalización (PEP FF)', 'pf', 12),
  ('d0000001-0000-4000-8000-00000000040d', 'c0000001-0000-4000-8000-000000000104', 'Permiso por Protección Temporal', 'pt', 13),
  ('d0000001-0000-4000-8000-00000000040e', 'c0000001-0000-4000-8000-000000000104', 'Persona sin identificación', 'si', 14),
  ('d0000001-0000-4000-8000-00000000040f', 'c0000001-0000-4000-8000-000000000104', 'Registro Civil de Nacimiento', 'rc', 15),
  ('d0000001-0000-4000-8000-000000000410', 'c0000001-0000-4000-8000-000000000104', 'Salvoconducto', 'sc', 16),
  ('d0000001-0000-4000-8000-000000000411', 'c0000001-0000-4000-8000-000000000104', 'Tarjeta Identidad', 'ti', 17),
  ('d0000001-0000-4000-8000-000000000412', 'c0000001-0000-4000-8000-000000000104', 'Visa', 'vs', 18),
  ('d0000001-0000-4000-8000-000000000413', 'c0000001-0000-4000-8000-000000000104', 'Visa de Refugiado', 'vr', 19),
  ('d0000001-0000-4000-8000-000000000414', 'c0000001-0000-4000-8000-000000000104', 'NIT de otro país', 'np', 20),
  -- I.7 Régimen SGSSS (SEMBRAR, option inline — decisión G-2)
  ('d0000001-0000-4000-8000-000000000701', 'c0000001-0000-4000-8000-000000000107', 'Contributivo', 'contributivo', 1),
  ('d0000001-0000-4000-8000-000000000702', 'c0000001-0000-4000-8000-000000000107', 'Subsidiado', 'subsidiado', 2),
  ('d0000001-0000-4000-8000-000000000703', 'c0000001-0000-4000-8000-000000000107', 'No asegurado', 'no_asegurado', 3),
  ('d0000001-0000-4000-8000-000000000704', 'c0000001-0000-4000-8000-000000000107', 'Régimen de excepción o especial', 'excepcion_especial', 4),
  -- I.8 EPS (SEMBRAR, option inline — decisión G-2)
  ('d0000001-0000-4000-8000-000000000801', 'c0000001-0000-4000-8000-000000000108', 'COOSALUD EPS-S', 'coosalud', 1),
  ('d0000001-0000-4000-8000-000000000802', 'c0000001-0000-4000-8000-000000000108', 'NUEVA EPS', 'nueva_eps', 2),
  ('d0000001-0000-4000-8000-000000000803', 'c0000001-0000-4000-8000-000000000108', 'MUTUAL SER', 'mutual_ser', 3),
  ('d0000001-0000-4000-8000-000000000804', 'c0000001-0000-4000-8000-000000000108', 'SALUD TOTAL EPS S.A.', 'salud_total', 4),
  ('d0000001-0000-4000-8000-000000000805', 'c0000001-0000-4000-8000-000000000108', 'EPS SANITAS', 'sanitas', 5),
  ('d0000001-0000-4000-8000-000000000806', 'c0000001-0000-4000-8000-000000000108', 'EPS SURA', 'sura', 6),
  ('d0000001-0000-4000-8000-000000000807', 'c0000001-0000-4000-8000-000000000108', 'FAMISANAR', 'famisanar', 7),
  ('d0000001-0000-4000-8000-000000000808', 'c0000001-0000-4000-8000-000000000108', 'SERVICIO OCCIDENTAL DE SALUD EPS SOS', 'sos', 8),
  ('d0000001-0000-4000-8000-000000000809', 'c0000001-0000-4000-8000-000000000108', 'SALUD MÍA', 'salud_mia', 9),
  ('d0000001-0000-4000-8000-00000000080a', 'c0000001-0000-4000-8000-000000000108', 'COMFENALCO VALLE', 'comfenalco_valle', 10),
  ('d0000001-0000-4000-8000-00000000080b', 'c0000001-0000-4000-8000-000000000108', 'COMPENSAR EPS', 'compensar', 11),
  ('d0000001-0000-4000-8000-00000000080c', 'c0000001-0000-4000-8000-000000000108', 'EPM — EMPRESAS PÚBLICAS DE MEDELLÍN', 'epm', 12),
  ('d0000001-0000-4000-8000-00000000080d', 'c0000001-0000-4000-8000-000000000108', 'FONDO DE PASIVO SOCIAL DE FERROCARRILES NACIONALES DE COLOMBIA', 'ferrocarriles', 13),
  ('d0000001-0000-4000-8000-00000000080e', 'c0000001-0000-4000-8000-000000000108', 'CAJACOPI ATLÁNTICO', 'cajacopi', 14),
  ('d0000001-0000-4000-8000-00000000080f', 'c0000001-0000-4000-8000-000000000108', 'CAPRESOCA', 'capresoca', 15),
  ('d0000001-0000-4000-8000-000000000810', 'c0000001-0000-4000-8000-000000000108', 'COMFACHOCÓ', 'comfachoco', 16),
  ('d0000001-0000-4000-8000-000000000811', 'c0000001-0000-4000-8000-000000000108', 'COMFAORIENTE', 'comfaoriente', 17),
  ('d0000001-0000-4000-8000-000000000812', 'c0000001-0000-4000-8000-000000000108', 'EPS FAMILIAR DE COLOMBIA', 'eps_familiar', 18),
  ('d0000001-0000-4000-8000-000000000813', 'c0000001-0000-4000-8000-000000000108', 'ASMET SALUD', 'asmet_salud', 19),
  ('d0000001-0000-4000-8000-000000000814', 'c0000001-0000-4000-8000-000000000108', 'EMSSANAR E.S.S.', 'emssanar', 20),
  ('d0000001-0000-4000-8000-000000000815', 'c0000001-0000-4000-8000-000000000108', 'CAPITAL SALUD EPS-S', 'capital_salud', 21),
  ('d0000001-0000-4000-8000-000000000816', 'c0000001-0000-4000-8000-000000000108', 'SAVIA SALUD EPS', 'savia_salud', 22),
  ('d0000001-0000-4000-8000-000000000817', 'c0000001-0000-4000-8000-000000000108', 'DUSAKAWI EPSI', 'dusakawi', 23),
  ('d0000001-0000-4000-8000-000000000818', 'c0000001-0000-4000-8000-000000000108', 'ASOCIACIÓN INDÍGENA DEL CAUCA AIC-EPSI', 'aic', 24),
  ('d0000001-0000-4000-8000-000000000819', 'c0000001-0000-4000-8000-000000000108', 'ANAS WAYUU EPSI', 'anas_wayuu', 25),
  ('d0000001-0000-4000-8000-00000000081a', 'c0000001-0000-4000-8000-000000000108', 'MALLAMAS EPSI', 'mallamas', 26),
  ('d0000001-0000-4000-8000-00000000081b', 'c0000001-0000-4000-8000-000000000108', 'PIJAOS SALUD EPSI', 'pijaos', 27),
  ('d0000001-0000-4000-8000-00000000081c', 'c0000001-0000-4000-8000-000000000108', 'FUERZAS MILITARES Y POLICÍA NACIONAL', 'ffmm_policia', 28),
  ('d0000001-0000-4000-8000-00000000081d', 'c0000001-0000-4000-8000-000000000108', 'FONDO NACIONAL DE PRESTACIONES SOCIALES DEL MAGISTERIO - FOMAG', 'fomag', 29),
  ('d0000001-0000-4000-8000-00000000081e', 'c0000001-0000-4000-8000-000000000108', 'Otra', 'otra', 30),
  -- I.20b / I.29b condición de "Otra": las preguntas gatillo (I.20 orientación
  -- sexual, I.29 ocupación) usan enums cuyo catálogo SÍ tiene opción "Otra"
  -- (código 'ot'); la condición compara contra su i_code (ver visibility_condition).
  -- I.19 (identidad de género) e I.22b (pueblo indígena) NO tienen opción "Otra"
  -- en su enum: sus campos de texto libre se eliminaron (decisión usuario, Opción A).
  -- II.1 Conformación familiar (SEMBRAR)
  ('d0000001-0000-4000-8000-000000002001', 'c0000001-0000-4000-8000-000000000201', 'Nuclear', 'nuclear', 1),
  ('d0000001-0000-4000-8000-000000002002', 'c0000001-0000-4000-8000-000000000201', 'Mononuclear', 'mononuclear', 2),
  ('d0000001-0000-4000-8000-000000002003', 'c0000001-0000-4000-8000-000000000201', 'Extensa o ampliada', 'extensa', 3),
  ('d0000001-0000-4000-8000-000000002004', 'c0000001-0000-4000-8000-000000000201', 'Reconstruida', 'reconstruida', 4),
  ('d0000001-0000-4000-8000-000000002005', 'c0000001-0000-4000-8000-000000000201', 'De crianza', 'crianza', 5),
  ('d0000001-0000-4000-8000-000000002006', 'c0000001-0000-4000-8000-000000000201', 'Poliamorosa', 'poliamorosa', 6),
  ('d0000001-0000-4000-8000-000000002007', 'c0000001-0000-4000-8000-000000000201', 'Unipersonal', 'unipersonal', 7),
  ('d0000001-0000-4000-8000-000000002008', 'c0000001-0000-4000-8000-000000000201', 'Homoparental', 'homoparental', 8),
  ('d0000001-0000-4000-8000-000000002009', 'c0000001-0000-4000-8000-000000000201', 'Elegida', 'elegida', 9),
  -- II.5 Proveedora / Dependiente (SEMBRAR)
  ('d0000001-0000-4000-8000-000000002051', 'c0000001-0000-4000-8000-000000000205', 'Proveedora', 'proveedora', 1),
  ('d0000001-0000-4000-8000-000000002052', 'c0000001-0000-4000-8000-000000000205', 'Dependiente', 'dependiente', 2),
  -- RG-1 Parentesco (SEMBRAR)
  ('d0000001-0000-4000-8000-000000002111', 'c0000001-0000-4000-8000-000000000211', 'Padre', 'padre', 1),
  ('d0000001-0000-4000-8000-000000002112', 'c0000001-0000-4000-8000-000000000211', 'Madre', 'madre', 2),
  ('d0000001-0000-4000-8000-000000002113', 'c0000001-0000-4000-8000-000000000211', 'Padrastro', 'padrastro', 3),
  ('d0000001-0000-4000-8000-000000002114', 'c0000001-0000-4000-8000-000000000211', 'Madrastra', 'madrastra', 4),
  ('d0000001-0000-4000-8000-000000002115', 'c0000001-0000-4000-8000-000000000211', 'Pareja', 'pareja', 5),
  ('d0000001-0000-4000-8000-000000002116', 'c0000001-0000-4000-8000-000000000211', 'Hermana(o)', 'hermana_o', 6),
  ('d0000001-0000-4000-8000-000000002117', 'c0000001-0000-4000-8000-000000000211', 'Hija(o)', 'hija_o', 7),
  ('d0000001-0000-4000-8000-000000002118', 'c0000001-0000-4000-8000-000000000211', 'Abuela(o)', 'abuela_o', 8),
  ('d0000001-0000-4000-8000-000000002119', 'c0000001-0000-4000-8000-000000000211', 'Otro integrante', 'otro', 9),
  -- III.1 / III.2 / III.3 Sí/No (SEMBRAR)
  ('d0000001-0000-4000-8000-000000003011', 'c0000001-0000-4000-8000-000000000301', 'Sí', 'si', 1),
  ('d0000001-0000-4000-8000-000000003012', 'c0000001-0000-4000-8000-000000000301', 'No', 'no', 2),
  ('d0000001-0000-4000-8000-000000003021', 'c0000001-0000-4000-8000-000000000302', 'Sí', 'si', 1),
  ('d0000001-0000-4000-8000-000000003022', 'c0000001-0000-4000-8000-000000000302', 'No', 'no', 2),
  ('d0000001-0000-4000-8000-000000003031', 'c0000001-0000-4000-8000-000000000303', 'Sí', 'si', 1),
  ('d0000001-0000-4000-8000-000000003032', 'c0000001-0000-4000-8000-000000000303', 'No', 'no', 2),
  -- IV.2 Apoyo para transporte (SEMBRAR, multiple)
  ('d0000001-0000-4000-8000-000000004021', 'c0000001-0000-4000-8000-000000000402', 'Movilización al alojamiento', 'alojamiento', 1),
  ('d0000001-0000-4000-8000-000000004022', 'c0000001-0000-4000-8000-000000000402', 'Diligencias judiciales', 'judiciales', 2),
  ('d0000001-0000-4000-8000-000000004023', 'c0000001-0000-4000-8000-000000000402', 'Actividades de estabilización', 'estabilizacion', 3),
  ('d0000001-0000-4000-8000-000000004024', 'c0000001-0000-4000-8000-000000000402', 'No aplica', 'no_aplica', 4),
  -- IV.3 Tipo de transporte (SEMBRAR, single)
  ('d0000001-0000-4000-8000-000000004031', 'c0000001-0000-4000-8000-000000000403', 'Terrestre', 'terrestre', 1),
  ('d0000001-0000-4000-8000-000000004032', 'c0000001-0000-4000-8000-000000000403', 'Aéreo', 'aereo', 2),
  ('d0000001-0000-4000-8000-000000004033', 'c0000001-0000-4000-8000-000000000403', 'Fluvial', 'fluvial', 3),
  ('d0000001-0000-4000-8000-000000004034', 'c0000001-0000-4000-8000-000000000403', 'Multimodal', 'multimodal', 4),
  ('d0000001-0000-4000-8000-000000004035', 'c0000001-0000-4000-8000-000000000403', 'No aplica', 'no_aplica', 5),
  -- V.1 Frecuencia de afectación de la alimentación (SEMBRAR, single)
  ('d0000001-0000-4000-8000-000000005011', 'c0000001-0000-4000-8000-000000000501', 'Siempre', 'siempre', 1),
  ('d0000001-0000-4000-8000-000000005012', 'c0000001-0000-4000-8000-000000000501', 'Casi siempre', 'casi_siempre', 2),
  ('d0000001-0000-4000-8000-000000005013', 'c0000001-0000-4000-8000-000000000501', 'Algunas veces', 'algunas_veces', 3),
  ('d0000001-0000-4000-8000-000000005014', 'c0000001-0000-4000-8000-000000000501', 'Nunca', 'nunca', 4),
  -- R.7 Estado del caso (SEMBRAR, single — informativo, sin efecto de dominio, G-6)
  ('d0000001-0000-4000-8000-000000009061', 'c0000001-0000-4000-8000-000000000906', 'Caracterizado', 'caracterizado', 1),
  ('d0000001-0000-4000-8000-000000009062', 'c0000001-0000-4000-8000-000000000906', 'No caracterizado', 'no_caracterizado', 2),
  -- RG-2 Beneficiario — Tipo de documento (mismo mapa DOCUMENT_TYPE_FORM2 que I.4)
  ('d0000001-0000-4000-8000-000000009131', 'c0000001-0000-4000-8000-000000000913', 'Carné Diplomático', 'cd', 1),
  ('d0000001-0000-4000-8000-000000009132', 'c0000001-0000-4000-8000-000000000913', 'Cédula de Ciudadanía', 'cc', 2),
  ('d0000001-0000-4000-8000-000000009133', 'c0000001-0000-4000-8000-000000000913', 'Cédula de Extranjería', 'ce', 3),
  ('d0000001-0000-4000-8000-000000009134', 'c0000001-0000-4000-8000-000000000913', 'Certificado de Nacido Vivo del país de origen', 'no', 4),
  ('d0000001-0000-4000-8000-000000009135', 'c0000001-0000-4000-8000-000000000913', 'Certificado de Nacido Vivo en Colombia', 'nc', 5),
  ('d0000001-0000-4000-8000-000000009136', 'c0000001-0000-4000-8000-000000000913', 'Certificado de Registro Administrativo de Migrantes Venezolanos - RAMV', 'mv', 6),
  ('d0000001-0000-4000-8000-000000009137', 'c0000001-0000-4000-8000-000000000913', 'Documento Extranjero o Documento Nacional de Identidad (del país de origen -DNI )', 'dn', 7),
  ('d0000001-0000-4000-8000-000000009138', 'c0000001-0000-4000-8000-000000000913', 'Otro no especificado', 'oe', 8),
  ('d0000001-0000-4000-8000-000000009139', 'c0000001-0000-4000-8000-000000000913', 'Partida de Nacimiento de su país de origen', 'pn', 9),
  ('d0000001-0000-4000-8000-00000000913a', 'c0000001-0000-4000-8000-000000000913', 'Pasaporte', 'ps', 10),
  ('d0000001-0000-4000-8000-00000000913b', 'c0000001-0000-4000-8000-000000000913', 'Permiso Especial de Permanencia (PEP)', 'pp', 11),
  ('d0000001-0000-4000-8000-00000000913c', 'c0000001-0000-4000-8000-000000000913', 'Permiso Especial de Permanencia para el Fomento de la Formalización (PEP FF)', 'pf', 12),
  ('d0000001-0000-4000-8000-00000000913d', 'c0000001-0000-4000-8000-000000000913', 'Permiso por Protección Temporal', 'pt', 13),
  ('d0000001-0000-4000-8000-00000000913e', 'c0000001-0000-4000-8000-000000000913', 'Persona sin identificación', 'si', 14),
  ('d0000001-0000-4000-8000-00000000913f', 'c0000001-0000-4000-8000-000000000913', 'Registro Civil de Nacimiento', 'rc', 15),
  ('d0000001-0000-4000-8000-000000009140', 'c0000001-0000-4000-8000-000000000913', 'Salvoconducto', 'sc', 16),
  ('d0000001-0000-4000-8000-000000009141', 'c0000001-0000-4000-8000-000000000913', 'Tarjeta Identidad', 'ti', 17),
  ('d0000001-0000-4000-8000-000000009142', 'c0000001-0000-4000-8000-000000000913', 'Visa', 'vs', 18),
  ('d0000001-0000-4000-8000-000000009143', 'c0000001-0000-4000-8000-000000000913', 'Visa de Refugiado', 'vr', 19),
  ('d0000001-0000-4000-8000-000000009144', 'c0000001-0000-4000-8000-000000000913', 'NIT de otro país', 'np', 20),
  -- RG-3 Medida solicitada — 9 medidas del Excel (etiquetas legibles; el mapeo a
  -- tipos de emergency_measure y la creación de medidas se hace en FASE 2b-2).
  ('d0000001-0000-4000-8000-000000009211', 'c0000001-0000-4000-8000-000000000921', 'Alojamiento transitorio (hotel) hasta 15 días', 'alojamiento_transitorio', 1),
  ('d0000001-0000-4000-8000-000000009212', 'c0000001-0000-4000-8000-000000000921', 'Alojamiento temporal (arriendo) hasta 3 meses', 'alojamiento_temporal', 2),
  ('d0000001-0000-4000-8000-000000009213', 'c0000001-0000-4000-8000-000000000921', 'Transporte', 'transporte', 3),
  ('d0000001-0000-4000-8000-000000009214', 'c0000001-0000-4000-8000-000000000921', 'Alimentación', 'alimentacion', 4),
  ('d0000001-0000-4000-8000-000000009215', 'c0000001-0000-4000-8000-000000000921', 'Vestuario', 'vestuario', 5),
  ('d0000001-0000-4000-8000-000000009216', 'c0000001-0000-4000-8000-000000000921', 'Aseo', 'aseo', 6),
  ('d0000001-0000-4000-8000-000000009217', 'c0000001-0000-4000-8000-000000000921', 'Atención psicosocial', 'apoyo_psico', 7),
  ('d0000001-0000-4000-8000-000000009218', 'c0000001-0000-4000-8000-000000000921', 'Apoyo exequial (familiares de víctimas de feminicidio)', 'apoyo_exequial', 8),
  ('d0000001-0000-4000-8000-000000009219', 'c0000001-0000-4000-8000-000000000921', 'Acción simbólica (familiares de víctimas de feminicidio)', 'accion_simbolica', 9)
ON CONFLICT (id) DO NOTHING;

-- ─── Condicionales (visibility_condition) ─────────────────────────────────────
-- Columnas: id, target_type, target_id, trigger_question_id, trigger_state_path,
-- trigger_value, operator. Las preguntas "siempre visibles" (I.30b, II.5b, III.*b,
-- IV.1b/IV.3b, V.1b, VI.3b) NO llevan condición: se muestran siempre.
-- NOTA: solo dos dropdowns con catálogo por enum tienen opción "Otra" con código
-- 'ot' confirmado en Locale.go: I.20 orientación sexual (sexual_orientation_ot='Otro')
-- e I.29 ocupación (occupation_ot='Otra'). OJO: el valor que guarda el motor es el
-- i_code de la opción, no 'ot'; por eso I.20b e I.29b usan el i_code resuelto arriba.
-- I.19 identidad de género (gender_identity) e I.22b pueblo indígena (indigenous_people)
-- NO tienen opción "Otra/Otro" en su catálogo, así que sus ramas de texto libre
-- (I.19b / I.22c) se eliminaron por decisión del usuario (Opción A). Si el negocio
-- requiere la captura "Otra" luego, se agregará la opción al enum en Locale.go
-- (fuente única), no aquí.
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator) VALUES
  -- I.20b "Otra orientación sexual" ← I.20 = i_code de la opción 'ot'
  ('e0000001-0000-4000-8000-000000000002', 'QUESTION', 'c0000001-0000-4000-8000-00000000012a', 'c0000001-0000-4000-8000-000000000120', NULL, ot_sexual_orientation, 'EQUALS'),
  -- I.22b "¿Cuál pueblo indígena?" ← I.22 = true
  ('e0000001-0000-4000-8000-000000000003', 'QUESTION', 'c0000001-0000-4000-8000-00000000012b', 'c0000001-0000-4000-8000-000000000122', NULL, 'true', 'EQUALS'),
  -- I.24 "Situación migratoria" ← I.23 = true
  ('e0000001-0000-4000-8000-000000000005', 'QUESTION', 'c0000001-0000-4000-8000-000000000124', 'c0000001-0000-4000-8000-000000000123', NULL, 'true', 'EQUALS'),
  -- I.29b "Otra ocupación" ← I.29 = i_code de la opción 'ot'
  ('e0000001-0000-4000-8000-000000000006', 'QUESTION', 'c0000001-0000-4000-8000-00000000012d', 'c0000001-0000-4000-8000-000000000129', NULL, ot_occupation, 'EQUALS'),
  -- I.33a-txt ← I.33a = true
  ('e0000001-0000-4000-8000-000000000007', 'QUESTION', 'c0000001-0000-4000-8000-000000000134', 'c0000001-0000-4000-8000-000000000133', NULL, 'true', 'EQUALS'),
  -- I.33b-txt ← I.33b = true
  ('e0000001-0000-4000-8000-000000000008', 'QUESTION', 'c0000001-0000-4000-8000-000000000136', 'c0000001-0000-4000-8000-000000000135', NULL, 'true', 'EQUALS'),
  -- I.33c-txt ← I.33c = true
  ('e0000001-0000-4000-8000-000000000009', 'QUESTION', 'c0000001-0000-4000-8000-000000000138', 'c0000001-0000-4000-8000-000000000137', NULL, 'true', 'EQUALS'),
  -- II.4 "¿Quiénes están a su cargo?" ← II.3 = true
  ('e0000001-0000-4000-8000-00000000000a', 'QUESTION', 'c0000001-0000-4000-8000-000000000204', 'c0000001-0000-4000-8000-000000000203', NULL, 'true', 'EQUALS'),
  -- RG-1 "Otro integrante ¿cuál?" ← Parentesco = 'otro' (dentro del mismo grupo)
  ('e0000001-0000-4000-8000-00000000000b', 'QUESTION', 'c0000001-0000-4000-8000-000000000212', 'c0000001-0000-4000-8000-000000000211', NULL, 'otro', 'EQUALS'),
  -- VII.1b "Información adicional de vestuario" ← VII.1 = true
  ('e0000001-0000-4000-8000-00000000000c', 'QUESTION', 'c0000001-0000-4000-8000-00000000070b', 'c0000001-0000-4000-8000-000000000701', NULL, 'true', 'EQUALS'),
  -- RG-2 "¿Cuál es la discapacidad o diagnóstico?" ← beneficiario discapacidad = true (dentro del grupo)
  ('e0000001-0000-4000-8000-00000000000d', 'QUESTION', 'c0000001-0000-4000-8000-000000000916', 'c0000001-0000-4000-8000-000000000915', NULL, 'true', 'EQUALS')
ON CONFLICT (id) DO NOTHING;

UPDATE salvia.form SET deleted_at = NULL, status = 'active',
  name = 'Entrevista de Caracterización',
  description = 'Entrevista de Caracterización de Medidas de Emergencia. FASE 2b-1 v3 del 2026: siembra la definición del formulario (8 secciones + resultado). No precarga answers, no enlaza el submission a la tarea ni crea medidas al completar.'
  WHERE id = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d';

END $$;
