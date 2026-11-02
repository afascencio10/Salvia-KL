-- WP-2 · Registro V3 — catálogos y columnas nuevas.
--   REG-07  Opciones «Rural» / «Urbano» de la zona de residencia (DP_30)
--   (los bloques de REG-08 y REG-11 se agregan más abajo, uno por tarea)
--
-- Idempotente: se puede aplicar más de una vez. Las filas del catálogo se buscan por NOMBRE
-- (victim_case_form2_enums_name es UNIQUE), no por (categoría, código), para no chocar con
-- filas que la BD real ya tenga con otro código.
-- Los textos que ve la persona están en src/salvia/config/Locale.go (clave = name).
--
-- Aplicar a mano en la BD (Supabase SQL editor o psql). Esta carpeta no la ejecuta la app.

BEGIN;

-- La secuencia del id del catálogo puede haber quedado por detrás del máximo (el catálogo se
-- cargó con ids explícitos). Se adelanta si hace falta; nunca se retrocede.
DO $$
DECLARE
    seq text := pg_get_serial_sequence('salvia.victim_case_form2_enums', 'victim_case_form2_enums_id');
    mx  integer;
BEGIN
    IF seq IS NOT NULL THEN
        SELECT MAX(victim_case_form2_enums_id) INTO mx FROM salvia.victim_case_form2_enums;
        IF mx IS NOT NULL THEN
            EXECUTE format('SELECT setval(%L, GREATEST(%s, (SELECT last_value FROM %s)))', seq, mx, seq);
        END IF;
    END IF;
END $$;

-- ── REG-07 · Zona de residencia (DP_30) ──────────────────────────────────────────────────
-- La categoría victim_case_form2_residence_zone existía como columna y campo pero sin opciones,
-- así que la lista desplegable salía vacía.
INSERT INTO salvia.victim_case_form2_enums
    (victim_case_form2_enums_i_code, victim_case_form2_enums_name,
     victim_case_form2_enums_code,   victim_case_form2_enums_category)
SELECT gen_random_uuid()::text, v.name, v.code, 'victim_case_form2_residence_zone'
FROM (VALUES
        ('victim_case_form2_residence_zone_ru', 'ru'),   -- Rural
        ('victim_case_form2_residence_zone_ur', 'ur')    -- Urbano
     ) AS v(name, code)
WHERE NOT EXISTS (
    SELECT 1 FROM salvia.victim_case_form2_enums e
    WHERE e.victim_case_form2_enums_name = v.name
);

COMMIT;
