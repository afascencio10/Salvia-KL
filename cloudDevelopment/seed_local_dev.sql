-- Usuarios de prueba sintéticos para la BD local del sandbox cloud de Claude
-- Code. La geografía y el catálogo de roles YA NO están acá — se cargan
-- completos y reales desde cloudDevelopment/seed_reference_data.sql (ver ese
-- archivo y actualizar-seed-con-composio.md). Este archivo solo contiene lo
-- que SÍ debe ser sintético: cuentas de prueba, nunca usuarios reales.
--
-- Idempotente: se puede correr varias veces sin duplicar filas.
-- Aplicar DESPUÉS de seed_reference_data.sql (depende de que los roles
-- 'ad'/'te'/'ro' ya existan).

-- ─── Usuarios de prueba ──────────────────────────────────────────────────
-- Password para los 3: CloudDev!Test2026  (bcrypt cost 10, mismo esquema que
-- la app usa — ver security/controllers/GeneralUserController.go bcrypt.*)
-- Generado con: htpasswd -bnBC 10 x 'CloudDev!Test2026' | cut -d: -f2
-- Solo para esta BD local del sandbox — no es una cuenta real.

INSERT INTO security.general_user_profile (general_user_profile_i_code, general_user_profile_creation_date, general_user_profile_update_date, general_user_profile_names, general_user_profile_last_names, general_user_profile_town)
VALUES
    ('profile-icode-cloud.test.ad', NOW(), NOW(), 'Cloud', 'Admin', '11001000'),
    ('profile-icode-cloud.test.te', NOW(), NOW(), 'Cloud', 'Tecnico', '11001000'),
    ('profile-icode-cloud.test.ro', NOW(), NOW(), 'Cloud', 'Operador', '11001000')
ON CONFLICT (general_user_profile_i_code) DO NOTHING;

INSERT INTO security.general_user (general_user_i_code, general_user_creation_date, general_user_update_date, general_user_login, general_user_password, general_user_status, general_user_language, general_user_general_user_profile)
SELECT 'user-icode-' || v.login, NOW(), NOW(), v.login, '$2y$10$tPgKXMz9c/TE7LY/iHG4luH5l3XQGmA9.ZA2CvJHmhNXqhYQ5Lz5G', 'e', 'sp', gup.general_user_profile_id
FROM (VALUES ('cloud.test.ad'), ('cloud.test.te'), ('cloud.test.ro')) AS v(login)
JOIN security.general_user_profile gup ON gup.general_user_profile_i_code = 'profile-icode-' || v.login
ON CONFLICT (general_user_login) DO NOTHING;

-- El rol 'ro' exige general_user_team en ('Riesgo alto'|'Riesgo bajo') — ver
-- GeneralUserController.go. Se setea vía migrateLegacyTables (columna
-- general_user_team) una vez la app la haya creado en el primer arranque;
-- si esta columna ya existe (arranques posteriores), el UPDATE aplica directo.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'security' AND table_name = 'general_user' AND column_name = 'general_user_team'
    ) THEN
        UPDATE security.general_user SET general_user_team = 'Riesgo alto'
        WHERE general_user_login = 'cloud.test.ro';
    END IF;
END $$;

INSERT INTO security.rel_role_general_user (rel_role_general_user_creation_date, role_id, general_user_id)
SELECT NOW(), r.role_id, gu.general_user_id
FROM (VALUES ('cloud.test.ad', 'ad'), ('cloud.test.te', 'te'), ('cloud.test.ro', 'ro')) AS v(login, role_code)
JOIN security.role r ON r.role_code = v.role_code
JOIN security.general_user gu ON gu.general_user_login = v.login
WHERE NOT EXISTS (
    SELECT 1 FROM security.rel_role_general_user rr
    WHERE rr.general_user_id = gu.general_user_id AND rr.role_id = r.role_id
);
