-- REG-16 (WP-4 T5) · Catálogo «Elementos psicosociales integrados en el acompañamiento» (GS_20 de Registro V3).
--
-- 14 opciones del cliente (códigos 01 a 14) más «Ninguno» (código 99, excluyente; queda al final porque el
-- listado se ordena por código). Categoría nueva: victim_case_form2_psychosocial_elements. El `name` técnico lleva
-- el prefijo «Elemento psicosocial:» y se recorta a 96 caracteres (límite de la columna); la etiqueta que ve la
-- agente sale de registro-v3-campos.js (ENUM_LABELS_OVERRIDE_V3), con el texto exacto del cliente.
-- Correspondencia con las 20 opciones de Seguimiento v4 (seed seguimiento_v4_form.sql): ver la planeación WP-4 (T5).
--
-- Idempotente: no inserta una opción si la categoría ya tiene ese código. En el sandbox cloud se aplica a mano:
--   psql -h 127.0.0.1 -U postgres -d salvia_local -f src/migrations/20261102_wp4_elementos_psicosociales.sql
INSERT INTO salvia.victim_case_form2_enums
    (victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)
SELECT gen_random_uuid()::text,
       LEFT('Elemento psicosocial: ' || v.label, 96),
       v.code,
       'victim_case_form2_psychosocial_elements'
FROM (VALUES
    ('01', 'Contención emocional en crisis'),
    ('02', 'Comprensión del ciclo de la violencia'),
    ('03', 'Identificación y fortalecimiento de recursos personales'),
    ('04', 'Despatologización de experiencias y reacciones asociadas a la violencia'),
    ('05', 'Co-construcción del Plan de Seguridad'),
    ('06', 'Identificación de señales tempranas de riesgo y estrategias de respuesta'),
    ('07', 'Orientación para acceso a servicios, programas o apoyos institucionales que fortalecen autonomía'),
    ('08', 'Fortalecimiento de la autonomía económica y vinculación a proyecto de vida'),
    ('09', 'Cuestionamiento de mandatos culturales que sostienen la violencia'),
    ('10', 'Identificación de redes seguras y afirmativas (familiares, comunitarias, institucionales)'),
    ('11', 'Reconocimiento y comprensión de los derechos de las personas víctimas de VBG y VpP'),
    ('12', 'Orientación sobre acciones disponibles para la garantía de derechos y la responsabilidad del Estado'),
    ('13', 'Fortalecimiento de la corresponsabilidad de la persona en la gestión de su seguridad'),
    ('14', 'Sensibilización y orientación sobre servicios de salud mental para acompañar los impactos de la violencia'),
    ('99', 'Ninguno')
) AS v(code, label)
WHERE NOT EXISTS (
    SELECT 1 FROM salvia.victim_case_form2_enums e
    WHERE e.victim_case_form2_enums_category = 'victim_case_form2_psychosocial_elements'
      AND e.victim_case_form2_enums_code = v.code
);
