-- Pregunta de riesgo inminente en el seguimiento vigente.
-- Visible al elegir Medidas de Emergencia. No reescribe envíos.
DO $$
DECLARE
  sec_id varchar;
  frm_id varchar;
  next_order int;
BEGIN
  IF EXISTS (
    SELECT 1 FROM salvia.question WHERE id = 'b7e4c2a1-9f30-4c58-8d16-5a2e7b3c0d48'
  ) THEN
    RETURN;
  END IF;

  SELECT form_id, form_section_id, "order" + 1
    INTO frm_id, sec_id, next_order
    FROM salvia.question
   WHERE id = '1a36260c-33a4-4ebd-bffb-e387d7964b96'
     AND deleted_at IS NULL;

  IF sec_id IS NULL THEN
    RETURN;
  END IF;

  UPDATE salvia.question
     SET "order" = "order" + 1
   WHERE form_section_id = sec_id
     AND "order" >= next_order
     AND deleted_at IS NULL;

  INSERT INTO salvia.question (
    id, form_id, form_section_id, repeater_group_id, question_type,
    description, required, state_options_path, "order"
  ) VALUES (
    'b7e4c2a1-9f30-4c58-8d16-5a2e7b3c0d48',
    frm_id,
    sec_id,
    NULL,
    'boolean',
    '¿La remisión a medidas de emergencia es de riesgo inminente?',
    true,
    NULL,
    next_order
  );

  INSERT INTO salvia.visibility_condition (
    id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator
  ) VALUES (
    'e8f5d3b2-0a41-4d69-9e27-6b3f8c4d1e59',
    'QUESTION',
    'b7e4c2a1-9f30-4c58-8d16-5a2e7b3c0d48',
    'e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc',
    NULL,
    'medidas_emergencia',
    'INCLUDES'
  );
END $$;
