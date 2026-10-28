-- Agrega Seguimiento a medidas de emergencia sin borrar envíos ni el resto del formulario.
-- Si la sección ya está, no hace nada. Generado por scripts/gen_seguimiento_v4_seed.js. No editar a mano.
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM salvia.form
     WHERE id = 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f' AND deleted_at IS NULL
  ) THEN
    RETURN;
  END IF;

  IF EXISTS (
    SELECT 1 FROM salvia.form_section WHERE id = 'fc583234-5a64-45f1-aa9e-46e499fe6da4'
  ) THEN
    RETURN;
  END IF;

  UPDATE salvia.form_section
     SET "order" = "order" + 1
   WHERE form_id = 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f'
     AND "order" >= 2
     AND deleted_at IS NULL;

  INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES
    ('fc583234-5a64-45f1-aa9e-46e499fe6da4', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'Seguimiento a medidas de emergencia', 'Solo cuando el caso ya tiene medidas de emergencia, en cualquier estado', 2);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order") VALUES
    ('b8a1ab54-83eb-47a6-a148-09a3856a1512', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'text', '¿De qué manera la medida ha contribuido a la gestión de su seguridad y la reducción del riesgo desde el último seguimiento?', true, NULL, 1),
    ('ab2698c7-c8a2-4173-8653-3a46789a5ac3', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'text', '¿Cómo se ha sentido con las medidas de emergencia recibidas y qué cambios han generado en su situación actual?', true, NULL, 2),
    ('1e13234a-963f-40fc-a8b4-8dcbb91821ea', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'text', '¿Cómo perciben las personas beneficiarias las medidas otorgadas y de qué manera estas han impactado su bienestar o seguridad? (Si las personas beneficiarias son menores de edad u otras personas dependientes, la valoración puede realizarse a partir del reporte de la participante).', true, NULL, 3),
    ('4085efbd-b581-4ce2-a785-a046c37b5f23', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'text', 'Además de las medidas de emergencia, ¿qué alternativas o estrategias está construyendo para fortalecer su protección y bienestar a mediano plazo?', true, NULL, 4),
    ('b0071c7f-3489-42bc-9a3e-ff27886105e9', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'info', 'Registre los acuerdos, acciones o gestiones que realizará la persona o los beneficiarios', false, NULL, 5),
    ('1af7af19-a097-42cf-9178-9d2b094d320c', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'text', 'Compromisos', true, NULL, 6),
    ('a9495e3a-a25c-4d05-9516-268d8dfce5f1', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'info', 'Registre información relevante identificada durante el seguimiento que no haya sido consignada en otros campos y que aporte a la comprensión del caso.', false, NULL, 7),
    ('e771e46d-3d5a-4575-b385-d5891c1cf400', 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', NULL, 'text', 'Observaciones', true, NULL, 8);

  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator) VALUES
    ('1cb02c1e-ed01-4da8-b451-ae95f3a7e6ea', 'SECTION', 'fc583234-5a64-45f1-aa9e-46e499fe6da4', '86a2fa68-74e3-4e65-86bf-28dfdf2ab8c3', 'hasMedidas', 'true', 'EQUALS');
END $$;
