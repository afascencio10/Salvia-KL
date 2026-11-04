-- Usuarios de prueba sintéticos para la BD local del sandbox cloud de Claude
-- Code. La geografía y el catálogo de roles YA NO están acá — se cargan
-- completos y reales desde cloudDevelopment/seed_reference_data.sql (ver ese
-- archivo y actualizar-seed-con-composio.md). Este archivo solo contiene lo
-- que SÍ debe ser sintético: cuentas de prueba, nunca usuarios reales.
--
-- UN USUARIO POR CADA UNO DE LOS 19 ROLES de security.role: login
-- cloud.test.<role_code> (p. ej. cloud.test.ad, cloud.test.sv, cloud.test.ps).
-- Misma clave para todos: CloudDev!Test2026
--
-- Atributos que un rol necesita para que su pantalla funcione (el resto de los
-- roles no exige nada más que el rol):
--   ro → general_user_team = 'Riesgo alto'  (obligatorio, GeneralUserController.go)
--   en → general_user_assigned_department = '3'  (department_id de BOGOTÁ, D.C.;
--        acota el listado del Enlace territorial a ese departamento)
--   et → entity_branch_id = sede de prueba (se crea abajo; Casos Entidad se acota por sede)
--
-- Lo que depende de tablas que crea la app al arrancar (GORM) NO puede ir en este
-- archivo, que corre ANTES del primer arranque: el mecanismo articulador del
-- usuario 'mc' y la dupla de 'ps' + 'ts' están en
-- seed_test_data/20261101-base-mecanismos-y-duplas.sql (fase 2, ver
-- apply-test-seeds.sh).
--
-- Idempotente: se puede correr varias veces sin duplicar filas. Los atributos de
-- arriba se vuelven a fijar en cada corrida (son el estado esperado del usuario
-- de prueba).
--
-- Aplicar DESPUÉS de seed_reference_data.sql (necesita los roles y las
-- entidades). Requiere las columnas general_user_team,
-- general_user_assigned_department y entity_branch_id de security.general_user:
-- en un arranque limpio las pre-crea setup-local-postgres.sh (las crearía la app
-- después, y este seed corre antes).

-- ─── Usuarios de prueba (login ↔ rol) ────────────────────────────────────
-- Password: CloudDev!Test2026  (bcrypt cost 10, mismo esquema que la app usa —
-- ver security/controllers/GeneralUserController.go bcrypt.*)
-- Generado con: htpasswd -bnBC 10 x 'CloudDev!Test2026' | cut -d: -f2
-- Solo para esta BD local del sandbox — no es una cuenta real.
--
-- general_user_i_code = 'user-icode-<login>' (no UUID, a propósito): las
-- duplas y otras tablas guardan el i_code del usuario como texto, y la BD real
-- ya tiene usuarios con i_code que no son UUID.

CREATE TEMP TABLE cloud_test_users (login VARCHAR(64), role_code VARCHAR(32));
INSERT INTO cloud_test_users (login, role_code) VALUES
    ('cloud.test.ad', 'ad'),
    ('cloud.test.op', 'op'),
    ('cloud.test.sv', 'sv'),
    ('cloud.test.ro', 'ro'),
    ('cloud.test.do', 'do'),
    ('cloud.test.no', 'no'),
    ('cloud.test.et', 'et'),
    ('cloud.test.en', 'en'),
    ('cloud.test.an', 'an'),
    ('cloud.test.us', 'us'),
    ('cloud.test.ps', 'ps'),
    ('cloud.test.ts', 'ts'),
    ('cloud.test.te', 'te'),
    ('cloud.test.mc', 'mc'),
    ('cloud.test.sm', 'sm'),
    ('cloud.test.em', 'em'),
    ('cloud.test.sc', 'sc'),
    ('cloud.test.cr', 'cr'),
    ('cloud.test.am', 'am');

-- Sede de prueba para el rol 'et' (entidad: la primera del catálogo real).
INSERT INTO salvia.entity_branch (entity_branch_i_code, entity_branch_creation_date, entity_branch_update_date, entity_branch_name, entity_branch_description, entity_branch_address, entity_branch_town_code, entity_branch_source, entity_id)
SELECT 'branch-icode-cloud.test.et', NOW(), NOW(), 'Sede de prueba (cloud.test.et)', 'Sede sintética del usuario de prueba cloud.test.et', 'Calle 1 # 1-1', '11001000', 'u', (SELECT MIN(entity_id) FROM salvia.entity)
WHERE NOT EXISTS (SELECT 1 FROM salvia.entity_branch WHERE entity_branch_i_code = 'branch-icode-cloud.test.et');

INSERT INTO security.general_user_profile (general_user_profile_i_code, general_user_profile_creation_date, general_user_profile_update_date, general_user_profile_names, general_user_profile_last_names, general_user_profile_town)
SELECT 'profile-icode-' || u.login, NOW(), NOW(), 'Cloud', r.role_name, '11001000'
FROM cloud_test_users u
JOIN security.role r ON r.role_code = u.role_code
ON CONFLICT (general_user_profile_i_code) DO NOTHING;

INSERT INTO security.general_user (general_user_i_code, general_user_creation_date, general_user_update_date, general_user_login, general_user_password, general_user_status, general_user_language, general_user_general_user_profile)
SELECT 'user-icode-' || u.login, NOW(), NOW(), u.login, '$2y$10$tPgKXMz9c/TE7LY/iHG4luH5l3XQGmA9.ZA2CvJHmhNXqhYQ5Lz5G', 'e', 'sp', gup.general_user_profile_id
FROM cloud_test_users u
JOIN security.general_user_profile gup ON gup.general_user_profile_i_code = 'profile-icode-' || u.login
ON CONFLICT (general_user_login) DO NOTHING;

INSERT INTO security.rel_role_general_user (rel_role_general_user_creation_date, role_id, general_user_id)
SELECT NOW(), r.role_id, gu.general_user_id
FROM cloud_test_users u
JOIN security.role r ON r.role_code = u.role_code
JOIN security.general_user gu ON gu.general_user_login = u.login
WHERE NOT EXISTS (
    SELECT 1 FROM security.rel_role_general_user rr
    WHERE rr.general_user_id = gu.general_user_id AND rr.role_id = r.role_id
);

-- ─── Atributos por rol ───────────────────────────────────────────────────
UPDATE security.general_user SET general_user_team = 'Riesgo alto'
WHERE general_user_login = 'cloud.test.ro';

UPDATE security.general_user SET general_user_assigned_department = '3'
WHERE general_user_login = 'cloud.test.en';

UPDATE security.general_user
SET entity_branch_id = (SELECT entity_branch_id FROM salvia.entity_branch WHERE entity_branch_i_code = 'branch-icode-cloud.test.et')
WHERE general_user_login = 'cloud.test.et';

DROP TABLE cloud_test_users;
