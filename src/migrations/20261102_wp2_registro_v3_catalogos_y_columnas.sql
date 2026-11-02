-- WP-2 · Registro V3 — catálogos y columnas nuevas.
--   REG-07  Opciones «Rural» / «Urbano» de la zona de residencia (DP_30)
--   REG-11  Departamentos de los hechos y de residencia (HV_4, DP_27) y municipios con texto libre
--   (el bloque de REG-08 se agrega más abajo, uno por tarea)
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

-- ── REG-11 · Hechos de violencia y residencia: columnas que no alcanzaban ─────────────────────
-- Departamento de los hechos (HV_4) y de residencia (DP_27): la pantalla los pedía pero no tenían
-- columna, así que nunca se guardaban. NULL = sin respuesta.
ALTER TABLE salvia.victim_case_form2
    ADD COLUMN IF NOT EXISTS victim_case_form2_facts_department     varchar(100),
    ADD COLUMN IF NOT EXISTS victim_case_form2_residence_department varchar(100);

-- Municipio de los hechos (HV_5) y de residencia (DP_28): la pantalla lo recoge como texto libre
-- (p. ej. «Barranquilla») pero las columnas «_town_code» medían 8 caracteres y el guardado fallaba con
-- un 500. Ampliar un varchar no reescribe la tabla.
ALTER TABLE salvia.victim_case_form2
    ALTER COLUMN victim_case_form2_facts_town_code TYPE varchar(100),
    ALTER COLUMN victim_case_form2_residence_town  TYPE varchar(100);
ALTER TABLE salvia.victim_case
    ALTER COLUMN victim_case_victim_town_code TYPE varchar(100);

COMMIT;
