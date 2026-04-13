-- Verificación del form de seguimiento en BD

-- 1. Secciones
SELECT 'SECTIONS' as tabla, id, name, "order"
FROM salvia.form_section
WHERE form_id = '2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff'
ORDER BY "order";

-- 2. Preguntas del repeater (sección 4)
SELECT 'REPEATER_QUESTIONS' as tabla, q.id, q.description, q.question_type, q."order", q.repeater_group_id
FROM salvia.question q
JOIN salvia.form_section s ON s.id::varchar = q.form_section_id
WHERE s.form_id = '2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff'
  AND s."order" = 4
  AND q.deleted_at IS NULL
ORDER BY q."order";

-- 3. TODAS las visibility_conditions del form
SELECT 'VISIBILITY_CONDITIONS' as tabla, vc.id, vc.target_type, vc.target_id, vc.trigger_question_id, vc.trigger_value, vc.operator
FROM salvia.visibility_condition vc
WHERE vc.target_id IN (
  SELECT id FROM salvia.question
  WHERE form_id = '2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff'
    AND deleted_at IS NULL
)
ORDER BY vc.trigger_value;

-- 4. Conteo de opciones por pregunta del repeater
SELECT 'OPTIONS_COUNT' as tabla, q.description, count(o.id) as total_options
FROM salvia.question q
LEFT JOIN salvia.option o ON o.question_id = q.id AND o.deleted_at IS NULL
WHERE q.form_id = '2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff'
  AND q.repeater_group_id IS NOT NULL
  AND q.deleted_at IS NULL
GROUP BY q.id, q.description
ORDER BY q."order";
