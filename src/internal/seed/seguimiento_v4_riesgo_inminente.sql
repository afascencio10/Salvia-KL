-- Pregunta de riesgo inminente en el seguimiento V.4, junto a las medidas.
-- UUID de q:plan.riesgoInminente, el mismo que seguimientoV4ID en form_service_v4.go.
-- Visible si hay remisión a equipos y el equipo incluye medidas.
DO $$
DECLARE
  sec_id varchar := 'c9d66717-fb8c-49e8-847d-4c5e1f880a2a';
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM salvia.form
     WHERE id = 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f' AND deleted_at IS NULL
  ) THEN
    RETURN;
  END IF;

  IF EXISTS (
    SELECT 1 FROM salvia.question WHERE id = '10164e84-0e54-4ebc-80f2-8143c90820e1'
  ) THEN
    RETURN;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM salvia.form_section WHERE id = sec_id AND deleted_at IS NULL
  ) THEN
    RETURN;
  END IF;

  UPDATE salvia.question
     SET "order" = "order" + 1
   WHERE form_section_id = sec_id
     AND "order" >= 5
     AND deleted_at IS NULL;

  INSERT INTO salvia.question (
    id, form_id, form_section_id, repeater_group_id, question_type,
    description, required, state_options_path, "order"
  ) VALUES (
    '10164e84-0e54-4ebc-80f2-8143c90820e1',
    'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f',
    sec_id,
    NULL,
    'boolean',
    '¿La remisión a medidas de emergencia es de riesgo inminente?',
    true,
    NULL,
    5
  );

  INSERT INTO salvia.visibility_condition (
    id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator
  ) VALUES
    ('9924a3a8-8e66-4ec6-bfca-328cb24ce7bb', 'QUESTION', '10164e84-0e54-4ebc-80f2-8143c90820e1', 'f5287d9b-1a3d-4257-b4a6-f2b222cf4c23', NULL, 'true', 'EQUALS'),
    ('bac0709a-4c70-40f9-b91f-cae1ffe20a5e', 'QUESTION', '10164e84-0e54-4ebc-80f2-8143c90820e1', 'f589cd41-b3a9-4070-879e-a6e590bcbe23', NULL, 'medidas', 'INCLUDES');
END $$;
