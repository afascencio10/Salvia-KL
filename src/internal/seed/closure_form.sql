-- =============================================================================
-- SEED: Formulario dinámico "Cierre de proceso psicosocial"
--
-- POR QUÉ EXISTE ESTE ARCHIVO
-- El flujo 3x3 abre este formulario por un UUID fijo en el código
-- (ClosureFormID en src/salvia/service/psychosocial_3x3_service.go). Si la fila
-- no está en la base, el modal de cierre abre y el motor dinamic-form responde
-- "Error 404 al cargar el formulario": el agente no puede cerrar el caso por
-- imposibilidad de contacto.
--
-- La documentación canónica (.knowledge/2-data-dictionary/closure-form-model.md)
-- decía que el seed existía, pero el archivo nunca se versionó (*.sql está en
-- .gitignore), así que el formulario solo vivía en las bases donde alguien lo
-- había insertado a mano. Para que eso no vuelva a pasar, este archivo va
-- embebido en el binario y la aplicación lo aplica sola al arrancar
-- (src/internal/seed/seed.go). Los UUIDs salen de ese documento, que sigue
-- siendo la fuente canónica.
--
-- TODO VA EN UN SOLO BLOQUE DO: así es una única sentencia, que es lo que admite
-- db.Exec, y a la vez queda atómico sin BEGIN/COMMIT explícitos. Mismo patrón que
-- src/cmd/seed/seed_psicosocial.sql.
--
-- ES IDEMPOTENTE: se puede aplicar sobre una base que ya lo tenga. Los
-- ON CONFLICT no reescriben lo existente y los UPDATE solo reviven filas
-- borradas lógicamente, que producen el mismo 404 que si no existieran.
--
-- Aplicarlo a mano (además del arranque):
--   go run cmd/seed/main.go internal/seed/closure_form.sql
--   psql "$DATABASE_URL" -f src/internal/seed/closure_form.sql
--
-- Verificar (debe dar 1, 1, 8, 12, 5):
--   SELECT
--     (SELECT COUNT(*) FROM salvia.form         WHERE id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND deleted_at IS NULL),
--     (SELECT COUNT(*) FROM salvia.form_section WHERE form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND deleted_at IS NULL),
--     (SELECT COUNT(*) FROM salvia.question     WHERE form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND deleted_at IS NULL),
--     (SELECT COUNT(*) FROM salvia.option o JOIN salvia.question q ON q.id = o.question_id
--        WHERE q.form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND o.deleted_at IS NULL),
--     (SELECT COUNT(*) FROM salvia.visibility_condition vc JOIN salvia.question q ON q.id = vc.target_id
--        WHERE q.form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b');
-- =============================================================================

DO $$
BEGIN

-- ── A. Formulario ────────────────────────────────────────────────────────────
INSERT INTO salvia.form (id, name, description, status) VALUES (
  'fcc7dc8d-835b-4559-9283-2ea8b8e6092b',
  'Cierre de proceso psicosocial',
  'Formulario para registrar el cierre de un proceso de Atención Psicosocial',
  'active'
) ON CONFLICT (id) DO NOTHING;

-- ── B. Sección ───────────────────────────────────────────────────────────────
INSERT INTO salvia.form_section (id, form_id, name, "order") VALUES (
  '0821f9a1-3bf2-4349-8cb8-5b0c337716fe',
  'fcc7dc8d-835b-4559-9283-2ea8b8e6092b',
  'Cierre',
  1
) ON CONFLICT (id) DO NOTHING;

-- ── C. Preguntas ─────────────────────────────────────────────────────────────
-- Las preguntas 7 y 8 quedan required = false a propósito: son obligatorias solo
-- cuando son visibles, y de esa validación condicional se encarga el frontend.
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order") VALUES
  ('c253f68c-866e-4093-ba80-a41bbb96b0ae', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'single',   '¿La llamada fue efectiva?',            true,  1),
  ('1c6d97ac-dc83-4e73-b747-14cd93ad422d', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'single',   'Motivo de cierre',                     true,  2),
  ('036f6945-d55a-496d-b343-b9c57eca954f', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'text',     'Contenido de la atención',             true,  3),
  ('96146a85-c194-470a-a2c8-c93413c88fc2', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'multiple', 'Plan de orientación',                  false, 4),
  ('9628c09c-5e5d-4d2a-82ff-d944491044ac', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'text',     'Temas trabajados durante la atención', true,  5),
  ('555d2810-c0b1-4422-95cc-02479fc061f4', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'boolean',  'Hay nuevos hechos de violencia',       true,  6),
  ('39c3f808-b100-42d6-83fd-78f7aec00f96', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'text',     'Descripción de los hechos',            false, 7),
  ('5e0713eb-90ac-4183-bcdf-f72d1caeb489', 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b', '0821f9a1-3bf2-4349-8cb8-5b0c337716fe', 'date',     'Fecha',                                false, 8)
ON CONFLICT (id) DO NOTHING;

-- ── D. Opciones ──────────────────────────────────────────────────────────────
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  -- ¿La llamada fue efectiva?
  ('4dff9f5d-36cf-4db7-a767-52b0d3add5f2', 'c253f68c-866e-4093-ba80-a41bbb96b0ae', 'Sí', 'si', 1),
  ('8f1bdbbd-1dcc-4b4f-a66f-52897df0e4d5', 'c253f68c-866e-4093-ba80-a41bbb96b0ae', 'No', 'no', 2),
  -- Motivo de cierre — 'imposibilidad_contacto_3x3' es el que preselecciona el
  -- flujo 3x3; sin esa opción el cierre reportado no se puede completar.
  ('53b8a4eb-d506-4215-870a-577e670e4907', '1c6d97ac-dc83-4e73-b747-14cd93ad422d', 'Cumplimiento de objetivos',        'cumplimiento_objetivos',     1),
  ('2d17669f-6f61-4c02-9aea-e9267e567bf0', '1c6d97ac-dc83-4e73-b747-14cd93ad422d', 'Cumplimiento esquema',             'cumplimiento_esquema',       2),
  ('9c25c5f0-13f2-499b-9c8f-5d39bdaf61ef', '1c6d97ac-dc83-4e73-b747-14cd93ad422d', 'No consentimiento',                'no_consentimiento',          3),
  ('d2c6065a-9da8-443e-a1d6-2cd5eee47a87', '1c6d97ac-dc83-4e73-b747-14cd93ad422d', 'Imposibilidad del contacto (3x3)', 'imposibilidad_contacto_3x3', 4),
  ('d054dd41-b927-4f85-b7bb-2949dc1576c9', '1c6d97ac-dc83-4e73-b747-14cd93ad422d', 'Desistimiento del proceso',        'desistimiento_proceso',      5),
  -- Plan de orientación
  ('d7dac006-bc8e-4d5b-83b3-e914073eb671', '96146a85-c194-470a-a2c8-c93413c88fc2', 'Enrutamiento',           'enrutamiento',        1),
  ('09fa69e1-9f65-4402-bf7c-4bad61aa450b', '96146a85-c194-470a-a2c8-c93413c88fc2', 'Activación de ruta',     'activacion_ruta',     2),
  ('9e1d08f3-68d3-4130-bbb1-c8440bd92b85', '96146a85-c194-470a-a2c8-c93413c88fc2', 'Seguimiento',            'seguimiento',         3),
  ('2cb7d162-5b36-4516-a28b-4abf8a4adde9', '96146a85-c194-470a-a2c8-c93413c88fc2', 'Medidas de emergencia',  'medidas_emergencia',  4),
  ('01d19bc2-9349-466a-aedf-bfd5fcc5fa4b', '96146a85-c194-470a-a2c8-c93413c88fc2', 'Plan de estabilización', 'plan_estabilizacion', 5)
ON CONFLICT (id) DO NOTHING;

-- ── E. Condiciones de visibilidad ────────────────────────────────────────────
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator) VALUES
  ('fae70d75-6d71-492f-bb36-54234f085bf7', 'QUESTION', '036f6945-d55a-496d-b343-b9c57eca954f', 'c253f68c-866e-4093-ba80-a41bbb96b0ae', 'si',   'EQUALS'),
  ('66aa281a-9a24-41d6-9c09-6cb73f44e485', 'QUESTION', '96146a85-c194-470a-a2c8-c93413c88fc2', 'c253f68c-866e-4093-ba80-a41bbb96b0ae', 'si',   'EQUALS'),
  ('61bd463f-c2a1-4da7-8bd0-029602e08e86', 'QUESTION', '9628c09c-5e5d-4d2a-82ff-d944491044ac', 'c253f68c-866e-4093-ba80-a41bbb96b0ae', 'si',   'EQUALS'),
  ('b7e5bd86-97fc-4bc4-b618-4d4a3bdf71bc', 'QUESTION', '39c3f808-b100-42d6-83fd-78f7aec00f96', '555d2810-c0b1-4422-95cc-02479fc061f4', 'true', 'EQUALS'),
  ('8ece8f82-75ec-4634-a44a-e51aa0bc0858', 'QUESTION', '5e0713eb-90ac-4183-bcdf-f72d1caeb489', '555d2810-c0b1-4422-95cc-02479fc061f4', 'true', 'EQUALS')
ON CONFLICT (id) DO NOTHING;

-- ── F. Reparación: revivir lo borrado lógicamente ────────────────────────────
-- Una fila con deleted_at produce el mismo 404 que una fila ausente, y el
-- ON CONFLICT de arriba no la toca.
UPDATE salvia.form         SET deleted_at = NULL, status = 'active'
  WHERE id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND deleted_at IS NOT NULL;
UPDATE salvia.form_section SET deleted_at = NULL
  WHERE form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND deleted_at IS NOT NULL;
UPDATE salvia.question     SET deleted_at = NULL
  WHERE form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b' AND deleted_at IS NOT NULL;
UPDATE salvia.option       SET deleted_at = NULL
  WHERE question_id IN (SELECT id FROM salvia.question WHERE form_id = 'fcc7dc8d-835b-4559-9283-2ea8b8e6092b')
    AND deleted_at IS NOT NULL;

END $$;
