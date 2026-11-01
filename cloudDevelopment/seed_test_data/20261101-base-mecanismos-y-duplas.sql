-- Base de configuración de los usuarios de prueba por rol (cloud.test.<rol>):
--   * catálogo de mecanismos articuladores (real, 6 filas) y el mecanismo del
--     usuario cloud.test.mc (Casos Mecanismo se acota por mecanismo);
--   * una dupla de prueba psicóloga + trabajadora social (cloud.test.ps /
--     cloud.test.ts) para Administrar Duplas y la asignación del flujo 3x3.
--
-- Origen: Supabase (proyecto pwwelwfhauspznpatuqm), leído vía Composio el
-- 2026-11-01: salvia.articulating_mechanism completa (6 filas, catálogo). La
-- dupla es sintética (los usuarios de prueba no existen en Supabase).
--
-- Vive en seed_test_data/ y no en seed_reference_data.sql / seed_local_dev.sql
-- porque salvia.articulating_mechanism y salvia.dupla las crea la APP al arrancar
-- (GORM); este archivo se aplica después, ver apply-test-seeds.sh.
-- Idempotente. Requiere los usuarios de seed_local_dev.sql.

-- Mecanismos articuladores: 2 nacionales, 2 departamentales (Risaralda, Bogotá), 2 municipales.
INSERT INTO salvia.articulating_mechanism (id, type, department_id, city_id, town_code, department_name, city_name, town_name, created_at, updated_at, deleted_at) VALUES
('547a725f-819a-4b2f-96e3-8e4578b4474b','nacional',NULL,NULL,NULL,NULL,NULL,NULL,'2026-09-10 16:40:47.122529+00','2026-09-10 16:40:47.122529+00',NULL),
('6762c0a3-4bc9-4963-86e7-c21619445b6c','municipal','20','8680','66170008','RISARALDA','DOSQUEBRADAS','AGUAZUL','2026-09-10 16:45:09.231888+00','2026-09-10 16:45:09.231888+00',NULL),
('7df585ba-1008-461e-8929-3c1b8ee1c88b','nacional',NULL,NULL,NULL,NULL,NULL,NULL,'2026-10-19 05:51:27.538742+00','2026-10-19 05:51:27.538742+00',NULL),
('cb1ce4b4-a075-4141-913b-15efbea98a44','departamental','20',NULL,NULL,'RISARALDA',NULL,NULL,'2026-09-10 16:43:16.559156+00','2026-09-10 16:43:16.559156+00',NULL),
('ef1a7d03-9a6a-4081-bb99-4b9bb8f37c8a','departamental','3',NULL,NULL,'BOGOTÁ, D.C.',NULL,NULL,'2026-10-19 05:54:36.647049+00','2026-10-19 05:54:36.647049+00',NULL),
('ef1cf837-77cd-4d24-a64b-14b97d19b00f','municipal','3','7993','11001000','BOGOTÁ, D.C.','BOGOTÁ, D.C.','BOGOTÁ, DISTRITO CAPITAL','2026-10-18 04:31:12.897311+00','2026-10-18 04:31:12.897311+00',NULL)
ON CONFLICT (id) DO NOTHING;

-- El usuario de prueba del rol 'mc' pertenece al mecanismo municipal de Bogotá.
UPDATE security.general_user
SET articulating_mechanism_id = 'ef1cf837-77cd-4d24-a64b-14b97d19b00f'
WHERE general_user_login = 'cloud.test.mc';

-- Dupla de prueba: guarda el i_code del usuario psicólogo y trabajador social.
INSERT INTO salvia.dupla (name, psychologist_id, social_worker_id, created_at, updated_at)
SELECT 'Dupla cloud.test', 'user-icode-cloud.test.ps', 'user-icode-cloud.test.ts', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM salvia.dupla WHERE name = 'Dupla cloud.test' AND deleted_at IS NULL);
