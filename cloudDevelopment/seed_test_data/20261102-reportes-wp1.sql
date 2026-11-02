-- Reportes de ejemplo para el sandbox cloud de Claude Code (WP-1: radicado, consulta
-- pública, historial del reporte y «Registro de caso»).
--
-- Sin esto el sandbox arranca con 0 reportes y cada sesión tendría que radicarlos y
-- gestionarlos a mano para ver el Listado de Reportes, Mis Reportes, el modal de gestión
-- con su historial, la consulta pública por radicado o la conversión a borrador de caso.
-- Trae 17 reportes con radicado «SEED-xxxx», uno por situación, y dos casos de Registro V3
-- enlazados a dos de ellos. Todo es sintético: nombres «… Semilla», teléfonos 300000000x /
-- 301000000x, descripciones de ejemplo. No sale de Supabase ni de una persona real.
--
--   Radicado    Tipo      Quién reporta    Estado en el Listado         Qué permite probar
--   SEED-SV22   Víctima   anónimo          Sin asignar                  consulta pública «recibido»
--   SEED-SR22   Tercero   anónimo          Sin asignar                  tercero con nombre propio
--   SEED-SF22   Familiar  anónimo          Sin asignar                  tipo sin destino en V3
--   SEED-SH22   Hombre    anónimo          Sin asignar                  sin persona afectada nombrada
--   SEED-SV23   Víctima   con sesión (op)  Sin asignar                  «Quién reporta»: usuario y rol
--   SEED-SR23   Tercero   con sesión (et)  Sin asignar                  ídem, otro rol
--   SEED-AV22   Víctima   anónimo          Asignado (op)                Mis Reportes de cloud.test.op
--   SEED-AF22   Familiar  anónimo          Asignado (op)                1 intento fallido hoy
--   SEED-CV22   Víctima   anónimo          Asignado (op)                contacto logrado → «Registro de caso»
--   SEED-CR22   Tercero   anónimo          Asignado (op)                ídem (tercero con reporte anónimo)
--   SEED-BV22   Víctima   anónimo          Borrador iniciado            caso V3 en borrador enlazado
--   SEED-TV22   Víctima   anónimo          Cerrado — válido             caso V3 terminado enlazado
--   SEED-XH22   Hombre    anónimo          Cerrado — inválido (ni)      cierre sin información mínima
--   SEED-XR22   Tercero   anónimo          Cerrado — inválido (nc)      esquema 3x3 completo (9 intentos)
--   SEED-XV22   Víctima   anónimo          Cerrado — inválido (dp)      cierre con observación
--   SEED-XF22   Familiar  anónimo          Cerrado — inválido (co)      cierre con observación
--   SEED-XV23   Víctima   anónimo          Cerrado — inválido           estado heredado 'i', sin motivo de cierre
--
-- Los casos de Registro V3 (SEED-BV22 en borrador, SEED-TV22 terminado) los crea la app al
-- convertir un reporte: usuario 'us' propio (cloud.seed.borrador / cloud.seed.terminado, clave
-- CloudDev!Test2026), el reporte enlazado en victim_case_victim_contact, CI_2, CI_4, HV_1 y CI_8
-- prellenados y el documento vacío. El formulario (victim_case_form2) repite lo que deja
-- CreateDraft (insertDummyForm2Row en victim_case_v3_repository.go): el id de «No» en los
-- enums, 0, '' y 1900-01-01. Si CreateDraft cambia, compáralo con un borrador creado desde la app.
--
-- Cuándo corre: fase 2 (apply-test-seeds.sh), después de que la app arranca, porque usa columnas
-- y tablas que crea ella (victim_contact_radicado, victim_contact_reporter_role,
-- case_timeline_event.victim_contact_id, victim_case_v3_answered_keys...). Necesita los 19
-- usuarios de seed_local_dev.sql (op, sv, et) y el catálogo real de enums.
--
-- Las fechas son relativas al momento de aplicarlo (hace N días) y en hora de Bogotá, como las
-- guarda la app. Idempotente: ids fijos «5eedNNNN-0000-4000-8000-…» y ON CONFLICT / NOT EXISTS;
-- volver a aplicarlo no duplica nada. Un archivo por requerimiento: si cambias algo, edita este.

-- ─── Ayudas de esta sesión (desaparecen al terminar) ─────────────────────
-- Id fijo y reconocible por tabla: kind 1 reporte, 2 form2, 3 form2_v2, 4 intento, 5 evento,
-- 6 caso, 8 perfil y 9 usuario del caso.
CREATE FUNCTION pg_temp.sid(kind int, n int) RETURNS text LANGUAGE sql IMMUTABLE AS
$$ SELECT '5eed' || lpad(kind::text, 4, '0') || '-0000-4000-8000-' || lpad(n::text, 12, '0') $$;

-- Id de una opción del catálogo por categoría y código (los ids dependen de cómo se cargó).
CREATE FUNCTION pg_temp.enum_id(cat text, code text) RETURNS bigint LANGUAGE sql STABLE AS
$$ SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums
    WHERE victim_case_form2_enums_category = cat AND victim_case_form2_enums_code = code $$;

-- ─── Los 17 reportes ─────────────────────────────────────────────────────
-- tipo: v víctima · r tercero · fm familiar · ho hombre. anonimato (solo r): pr nombre propio · an anónimo.
-- reporta: login de quien reporta con sesión (NULL = anónimo). nom/ape/tel_reporta: quien reporta
-- (en v no hay; en r solo si es de nombre propio). llamada/hora: ¿acepta que lo llamen? y a qué hora
-- (la de la víctima en v y r; la de quien reporta en fm y ho). dias: hace cuántos días se radicó.
-- agente: login del agente asignado. cierre: ni | nc | co | dp. cierre_dias: hace cuántos días se cerró.
CREATE TEMP TABLE seed_reporte (
    n int PRIMARY KEY, radicado text, tipo text, anonimato text, reporta text,
    nombres text, apellidos text, tel_victima bigint,
    nom_reporta text, ape_reporta text, tel_reporta bigint,
    llamada text, hora time, parentesco text, escenario text,
    dias int, agente text, cierre text, cierre_obs text, cierre_dias int, heredado boolean DEFAULT false
);
INSERT INTO seed_reporte
    (n, radicado, tipo, anonimato, reporta, nombres, apellidos, tel_victima, nom_reporta, ape_reporta, tel_reporta, llamada, hora, parentesco, escenario, dias, agente, cierre, cierre_obs, cierre_dias, heredado)
VALUES
 (1,  'SEED-SV22', 'v',  NULL, NULL,                'Valeria',  'Semilla Uno',        3000000001, NULL,       NULL,              NULL,       'y', '10:00', NULL, 'víctima, sin asignar',                              1, NULL,            NULL, NULL, NULL, false),
 (2,  'SEED-SR22', 'r',  'pr', NULL,                'Teresa',   'Semilla Dos',        3000000002, 'Tomás',    'Semilla Reporta', NULL,       'y', '11:30', NULL, 'tercero con nombre propio, sin asignar',            1, NULL,            NULL, NULL, NULL, false),
 (3,  'SEED-SF22', 'fm', NULL, NULL,                'Fernanda', 'Semilla Tres',       NULL,       'Felipe',   'Semilla Familia', 3010000003, 'y', '14:30', 'fa', 'familiar, sin asignar',                             2, NULL,            NULL, NULL, NULL, false),
 (4,  'SEED-SH22', 'ho', NULL, NULL,                '',         '',                   NULL,       'Héctor',   'Semilla Cuatro',  3010000004, 'n', NULL,    NULL, 'hombre, sin asignar',                               2, NULL,            NULL, NULL, NULL, false),
 (5,  'SEED-SV23', 'v',  NULL, 'cloud.test.op',     'Valentina','Semilla Cinco',      3000000005, NULL,       NULL,              NULL,       'y', '09:00', NULL, 'víctima con sesión (op), sin asignar',              0, NULL,            NULL, NULL, NULL, false),
 (6,  'SEED-SR23', 'r',  'pr', 'cloud.test.et',     'Tatiana',  'Semilla Seis',       3000000006, 'Ernesto',  'Semilla Entidad', NULL,       'y', '16:15', NULL, 'tercero con sesión (et), sin asignar',              0, NULL,            NULL, NULL, NULL, false),
 (7,  'SEED-AV22', 'v',  NULL, NULL,                'Adriana',  'Semilla Siete',      3000000007, NULL,       NULL,              NULL,       'y', '16:00', NULL, 'víctima asignada a op, sin intentos',               3, 'cloud.test.op', NULL, NULL, NULL, false),
 (8,  'SEED-AF22', 'fm', NULL, NULL,                'Alicia',   'Semilla Ocho',       NULL,       'Amparo',   'Semilla Familia', 3010000008, 'y', '08:30', 'fs', 'familiar asignado a op, un intento fallido hoy',    3, 'cloud.test.op', NULL, NULL, NULL, false),
 (9,  'SEED-CV22', 'v',  NULL, NULL,                'Camila',   'Semilla Nueve',      3000000009, NULL,       NULL,              NULL,       'y', '15:00', NULL, 'víctima asignada a op, contacto logrado',           4, 'cloud.test.op', NULL, NULL, NULL, false),
 (10, 'SEED-CR22', 'r',  'an', NULL,                'Carolina', 'Semilla Diez',       3000000010, NULL,       NULL,              NULL,       'y', '13:00', NULL, 'tercero anónimo asignado a op, contacto logrado',   4, 'cloud.test.op', NULL, NULL, NULL, false),
 (11, 'SEED-BV22', 'v',  NULL, NULL,                'Beatriz',  'Semilla Once',       3000000011, NULL,       NULL,              NULL,       'y', '10:30', NULL, 'víctima con borrador de caso',                      5, 'cloud.test.op', NULL, NULL, NULL, false),
 (12, 'SEED-TV22', 'v',  NULL, NULL,                'Tamara',   'Semilla Doce',       3000000012, NULL,       NULL,              NULL,       'y', '09:15', NULL, 'víctima con caso terminado',                        6, 'cloud.test.op', NULL, NULL, NULL, false),
 (13, 'SEED-XH22', 'ho', NULL, NULL,                '',         '',                   NULL,       'Hugo',     'Semilla Trece',   3010000013, 'n', NULL,    NULL, 'hombre, cerrado por no información mínima',         7, NULL,            'ni', NULL, 6, false),
 (14, 'SEED-XR22', 'r',  'pr', NULL,                'Ximena',   'Semilla Catorce',    3000000014, 'Roberto',  'Semilla Reporta', NULL,       'y', '09:00', NULL, 'tercero cerrado por no contacto (3x3 completo)',    8, 'cloud.test.op', 'nc', NULL, 3, false),
 (15, 'SEED-XV22', 'v',  NULL, NULL,                'Vanessa',  'Semilla Quince',     3000000015, NULL,       NULL,              NULL,       'y', '14:00', NULL, 'víctima cerrada por decisión de la persona',        7, 'cloud.test.op', 'dp', 'La persona decidió no continuar con el proceso.', 4, false),
 (16, 'SEED-XF22', 'fm', NULL, NULL,                'Fabiana',  'Semilla Dieciseis',  NULL,       'Francisco','Semilla Familia', 3010000016, 'y', '17:00', 'ot', 'familiar cerrado por no competencia',               9, 'cloud.test.op', 'co', 'Se orienta a la entidad competente.', 5, false),
 (17, 'SEED-XV23', 'v',  NULL, NULL,                'Viviana',  'Semilla Diecisiete', 3000000017, NULL,       NULL,              NULL,       'n', NULL,    NULL, 'víctima con estado heredado inválido',             10, NULL,            NULL, NULL, NULL, true);

-- Instantes derivados (hora de Bogotá, naive, como guarda la app victim_contact_creation_date).
ALTER TABLE seed_reporte ADD COLUMN creado timestamp, ADD COLUMN cerrado_en timestamptz;
UPDATE seed_reporte
   SET creado = (now() AT TIME ZONE 'America/Bogota') - dias * interval '1 day' - n * interval '7 minutes',
       cerrado_en = CASE WHEN cierre IS NOT NULL THEN now() - cierre_dias * interval '1 day' END;

-- Texto de cada tipo y de cada motivo de cierre (los mismos que usa la app en el historial).
CREATE TEMP TABLE seed_tipo (tipo text PRIMARY KEY, etiqueta text);
INSERT INTO seed_tipo VALUES ('v', 'Víctima'), ('r', 'Tercero'), ('fm', 'Familiar'), ('ho', 'Hombre');
CREATE TEMP TABLE seed_cierre (code text PRIMARY KEY, etiqueta text);
INSERT INTO seed_cierre VALUES
    ('ni', 'Inválido por No información Mínima'),
    ('nc', 'Inválido por Cierre por no contacto'),
    ('co', 'Inválido por Cierre por no competencia'),
    ('dp', 'Inválido por Cierre por decisión de la persona');

-- ─── victim_contact ──────────────────────────────────────────────────────
-- Anónimo: victim_contact_general_user = '' y sin rol. Con sesión: icode del usuario y su rol.
-- Sin agente = NULL; asignado = icode del agente y equipo '' (op no tiene equipo).
INSERT INTO salvia.victim_contact (
    victim_contact_i_code, victim_contact_creation_date, victim_contact_update_date,
    victim_contact_general_user, victim_contact_names, victim_contact_last_names,
    victim_contact_status, victim_contact_latitude, victim_contact_longitude,
    victim_contact_agent_id, victim_contact_team,
    victim_contact_close_reason, victim_contact_close_observation, victim_contact_closed_date,
    victim_contact_radicado, victim_contact_reporter_role)
SELECT pg_temp.sid(1, r.n), r.creado, COALESCE(r.cerrado_en AT TIME ZONE 'America/Bogota', r.creado),
       COALESCE('user-icode-' || r.reporta, ''), r.nombres, r.apellidos,
       CASE WHEN r.heredado THEN 'i' ELSE 'v' END, 0, 0,
       'user-icode-' || r.agente, CASE WHEN r.agente IS NOT NULL THEN '' END,
       CASE WHEN r.cierre IS NOT NULL THEN pg_temp.enum_id('victim_contact_close_reason', r.cierre) END,
       r.cierre_obs, r.cerrado_en,
       r.radicado, substring(r.reporta FROM 'cloud\.test\.(.*)')
  FROM seed_reporte r
ON CONFLICT DO NOTHING;

-- ─── victim_contact_form2 ────────────────────────────────────────────────
-- Lo mismo que guarda el formulario público V2 según el tipo: el teléfono y la hora de la víctima
-- solo en v y r; quien reporta (nombre y teléfono) solo en fm y ho (en r, solo el nombre);
-- «cuidado» y «la víctima sabe» solo en r y fm / r.
INSERT INTO salvia.victim_contact_form2 (
    victim_case_form2_i_code, victim_case_form2_creation_date, victim_case_form2_update_date,
    victim_contact_form2_will_receive_call, victim_contact_form2_has_care_role, victim_contact_form2_victim_aware_of_report,
    victim_contact_form2_reporter_names, victim_contact_form2_reporter_phone, victim_contact_form2_victim_col_phone,
    victim_contact_form2_facts_description, victim_contact_form2_best_contact_time,
    victim_contact_form2_report_type, victim_contact_form2_victim_contact)
SELECT pg_temp.sid(2, r.n), r.creado, r.creado,
       CASE WHEN r.tipo IN ('v', 'r')  THEN pg_temp.enum_id('yes_no', r.llamada) END,
       CASE WHEN r.tipo IN ('r', 'fm') THEN pg_temp.enum_id('yes_no', 'y') END,
       CASE WHEN r.tipo = 'r'          THEN pg_temp.enum_id('yes_no', 'y') END,
       COALESCE(r.nom_reporta, ''), COALESCE(r.tel_reporta, 0), COALESCE(r.tel_victima, 0),
       '[Semilla] ' || r.escenario || '. Descripción de ejemplo de los hechos; no contiene datos reales.',
       CASE WHEN r.tipo IN ('v', 'r') AND r.llamada = 'y' THEN r.hora ELSE time '00:00:00' END,
       pg_temp.enum_id('victim_case_form2_report_type_v2', r.tipo),
       vc.victim_contact_id
  FROM seed_reporte r
  JOIN salvia.victim_contact vc ON vc.victim_contact_i_code = pg_temp.sid(1, r.n)
ON CONFLICT DO NOTHING;

-- ─── victim_contact_form2_v2 ─────────────────────────────────────────────
-- Autorización de datos: siempre «Sí» (sin ella no se radica). Anonimato solo en r, parentesco solo
-- en fm; llamada y hora de quien reporta solo en fm y ho.
INSERT INTO salvia.victim_contact_form2_v2 (
    victim_contact_form2_v2_i_code, victim_contact_form2_v2_creation_date, victim_contact_form2_v2_update_date,
    victim_contact_form2_v2_victim_contact_form2, victim_contact_form2_v2_authorization_answer,
    victim_contact_form2_v2_anonimato, victim_contact_form2_v2_parentesco,
    victim_contact_form2_v2_reporter_last_names,
    victim_contact_form2_v2_reporter_will_receive_call, victim_contact_form2_v2_reporter_best_contact_time)
SELECT pg_temp.sid(3, r.n), r.creado, r.creado,
       f.victim_contact_form2_id, pg_temp.enum_id('yes_no', 'y'),
       CASE WHEN r.tipo = 'r'  THEN pg_temp.enum_id('victim_case_form2_anonimato', r.anonimato) END,
       CASE WHEN r.tipo = 'fm' THEN pg_temp.enum_id('victim_case_form2_parentesco', r.parentesco) END,
       r.ape_reporta,
       CASE WHEN r.tipo IN ('fm', 'ho') THEN pg_temp.enum_id('yes_no', r.llamada) END,
       CASE WHEN r.tipo IN ('fm', 'ho') AND r.llamada = 'y' THEN r.hora END
  FROM seed_reporte r
  JOIN salvia.victim_contact_form2 f ON f.victim_case_form2_i_code = pg_temp.sid(2, r.n)
ON CONFLICT DO NOTHING;

-- Ajuste razonable de cada reporte: «Sin ajuste» (sc), la opción por defecto del formulario.
INSERT INTO salvia.rel_victim_case_form2_enums_victim_contact_form2 (victim_case_form2_enums_id, victim_contact_form2_id)
SELECT pg_temp.enum_id('victim_case_form2_adjustments_gbv', 'sc'), f.victim_contact_form2_id
  FROM seed_reporte r
  JOIN salvia.victim_contact_form2 f ON f.victim_case_form2_i_code = pg_temp.sid(2, r.n)
 WHERE NOT EXISTS (
        SELECT 1 FROM salvia.rel_victim_case_form2_enums_victim_contact_form2 x
         WHERE x.victim_contact_form2_id = f.victim_contact_form2_id
           AND x.victim_case_form2_enums_id = pg_temp.enum_id('victim_case_form2_adjustments_gbv', 'sc'));

-- ─── Intentos de contacto ────────────────────────────────────────────────
-- resultado: fallido | exitoso. Con dias+hora: ese día (Bogotá) a esa hora; con hace_min: hace N
-- minutos (los de «hoy»). El de SEED-XR22 es el esquema 3x3 completo: 3 fallidos al día durante 3 días.
CREATE TEMP TABLE seed_intento (n int, resultado text, dias int, hora time, hace_min int);
INSERT INTO seed_intento (n, resultado, dias, hora, hace_min) VALUES
 (8,  'fallido', NULL, NULL,    90),
 (9,  'fallido', 3,    '09:10', NULL),
 (9,  'exitoso', 2,    '10:30', NULL),
 (10, 'exitoso', 2,    '11:05', NULL),
 (11, 'exitoso', 3,    '09:40', NULL),
 (12, 'exitoso', 4,    '09:20', NULL),
 (14, 'fallido', 6,    '09:00', NULL), (14, 'fallido', 6, '09:30', NULL), (14, 'fallido', 6, '10:00', NULL),
 (14, 'fallido', 5,    '09:00', NULL), (14, 'fallido', 5, '09:30', NULL), (14, 'fallido', 5, '10:00', NULL),
 (14, 'fallido', 4,    '09:00', NULL), (14, 'fallido', 4, '09:30', NULL), (14, 'fallido', 4, '10:00', NULL),
 (15, 'exitoso', 5,    '14:15', NULL),
 (16, 'exitoso', 6,    '15:45', NULL);
ALTER TABLE seed_intento ADD COLUMN en timestamptz, ADD COLUMN k int, ADD COLUMN n_dia int;
UPDATE seed_intento
   SET en = CASE WHEN dias IS NOT NULL
                 THEN (date_trunc('day', now() AT TIME ZONE 'America/Bogota') - dias * interval '1 day' + hora) AT TIME ZONE 'America/Bogota'
                 ELSE now() - hace_min * interval '1 minute' END;
-- k: orden del intento dentro del reporte; n_dia: cuántos fallidos van ese día (el «N hoy» del historial).
UPDATE seed_intento i
   SET k = o.k, n_dia = o.n_dia
  FROM (SELECT ctid AS id,
               row_number() OVER (PARTITION BY n ORDER BY en) AS k,
               row_number() OVER (PARTITION BY n, (en AT TIME ZONE 'America/Bogota')::date, resultado ORDER BY en) AS n_dia
          FROM seed_intento) o
 WHERE i.ctid = o.id;

-- victim_contact_attempt_date va en UTC (timestamp sin zona), como lo escribe la app.
INSERT INTO salvia.victim_contact_attempt (
    victim_contact_attempt_i_code, victim_contact_attempt_victim_contact, victim_contact_attempt_result,
    victim_contact_attempt_date, victim_contact_attempt_agent, victim_contact_attempt_creation_date)
SELECT pg_temp.sid(4, i.n * 100 + i.k), pg_temp.sid(1, i.n), i.resultado,
       i.en AT TIME ZONE 'UTC', 'user-icode-' || r.agente, i.en AT TIME ZONE 'UTC'
  FROM seed_intento i
  JOIN seed_reporte r ON r.n = i.n
ON CONFLICT DO NOTHING;

-- ─── Casos de Registro V3 enlazados (SEED-BV22 y SEED-TV22) ──────────────
-- estado: bo borrador · te terminado. El caso se crea el día del contacto logrado, un rato después.
CREATE TEMP TABLE seed_caso (n int PRIMARY KEY, estado text, login text, dias int, hora time);
INSERT INTO seed_caso VALUES
 (11, 'bo', 'cloud.seed.borrador',  3, '10:00'),
 (12, 'te', 'cloud.seed.terminado', 4, '09:45');
ALTER TABLE seed_caso ADD COLUMN creado timestamp;
UPDATE seed_caso
   SET creado = date_trunc('day', now() AT TIME ZONE 'America/Bogota') - dias * interval '1 day' + hora;

-- Usuario 'us' del caso (el que la app genera con credenciales aleatorias), con la clave del sandbox.
INSERT INTO security.general_user_profile (
    general_user_profile_i_code, general_user_profile_creation_date, general_user_profile_update_date,
    general_user_profile_names, general_user_profile_last_names,
    general_user_profile_doc_type, general_user_profile_doc_number, general_user_profile_town)
SELECT pg_temp.sid(8, c.n), c.creado, c.creado, 'Pendiente', 'Pendiente', 'cc', '0000000000', '00000000'
  FROM seed_caso c
ON CONFLICT DO NOTHING;

INSERT INTO security.general_user (
    general_user_i_code, general_user_creation_date, general_user_update_date,
    general_user_login, general_user_password, general_user_status, general_user_language,
    general_user_general_user_profile)
SELECT pg_temp.sid(9, c.n), c.creado, c.creado, c.login,
       '$2y$10$tPgKXMz9c/TE7LY/iHG4luH5l3XQGmA9.ZA2CvJHmhNXqhYQ5Lz5G',   -- CloudDev!Test2026 (igual que seed_local_dev.sql)
       'e', 'sp', gup.general_user_profile_id
  FROM seed_caso c
  JOIN security.general_user_profile gup ON gup.general_user_profile_i_code = pg_temp.sid(8, c.n)
ON CONFLICT DO NOTHING;

INSERT INTO security.rel_role_general_user (rel_role_general_user_creation_date, role_id, general_user_id)
SELECT c.creado, ro.role_id, gu.general_user_id
  FROM seed_caso c
  JOIN security.general_user gu ON gu.general_user_i_code = pg_temp.sid(9, c.n)
  JOIN security.role ro ON ro.role_code = 'us'
 WHERE NOT EXISTS (SELECT 1 FROM security.rel_role_general_user x
                    WHERE x.general_user_id = gu.general_user_id AND x.role_id = ro.role_id);

-- El caso: nombres del reporte, documento vacío, agente y equipo vacíos, y las preguntas que la
-- conversión deja respondidas (CI_2 nombres, CI_4 teléfono, HV_1 hechos, CI_8 ajuste razonable).
INSERT INTO salvia.victim_case (
    victim_case_i_code, victim_case_creation_date, victim_case_update_date, victim_case_status,
    victim_case_general_user, victim_case_victim_names, victim_case_victim_last_names,
    victim_case_victim_doc_type, victim_case_victim_doc_number, victim_case_victim_town_code,
    victim_case_victim_contact, victim_case_new_user_credentials, victim_case_v3_answered_keys,
    victim_case_team, agent_id)
SELECT pg_temp.sid(6, c.n), c.creado, c.creado, c.estado,
       pg_temp.sid(9, c.n), r.nombres, r.apellidos,
       '', '', '00000000',
       vc.victim_contact_id, jsonb_build_object('login', c.login, 'pass', 'CloudDev!Test2026'), 'CI_2,CI_4,HV_1,CI_8',
       '', ''
  FROM seed_caso c
  JOIN seed_reporte r ON r.n = c.n
  JOIN salvia.victim_contact vc ON vc.victim_contact_i_code = pg_temp.sid(1, c.n)
ON CONFLICT DO NOTHING;

-- Formulario del caso (victim_case_form2): una fila con las columnas NOT NULL resueltas como lo
-- hace CreateDraft, más los datos del reporte (hechos y teléfono de la víctima). Se arma desde una
-- plantilla en JSON para no repetir 120 columnas; lo que la plantilla no cubra y sea NOT NULL sin
-- valor por defecto se rellena según su tipo, igual que la «red de seguridad» de CreateDraft.
DO $$
DECLARE
    v_no  bigint;
    tpl   jsonb;
    c     record;
    m     record;
    f2    jsonb;
    cols  text;
BEGIN
    SELECT victim_case_form2_enums_id INTO STRICT v_no
      FROM salvia.victim_case_form2_enums
     WHERE victim_case_form2_enums_category = 'yes_no' AND victim_case_form2_enums_code = 'n';

    tpl := (SELECT jsonb_object_agg('victim_case_form2_' || k, v_no)   -- el id de «No» en todos los enums
              FROM unnest(ARRAY[
            'aggressor_attempted_strangulation', 'aggressor_capable_of_killing',
            'aggressor_common_spaces', 'aggressor_consumes_drugs', 'aggressor_controls',
            'aggressor_doc_type', 'aggressor_forced_sex', 'aggressor_gender_identity',
            'aggressor_had_hit_in_vulnerability', 'aggressor_has_access_to_weapons',
            'aggressor_has_penal_background', 'aggressor_has_penal_background_2',
            'aggressor_hierarchy', 'aggressor_is_alcoholic', 'aggressor_is_alcoholic_2',
            'aggressor_limits_contact_support_networks', 'aggressor_pursues_spies_destroys',
            'aggressor_sexually_harassment', 'aggressor_sexually_harassment_2',
            'aggressor_taken_advantage_physical_vulnerabil', 'aggressor_threat_kill',
            'aggressor_threatened_damage_members', 'aggressor_unemployed',
            'aggressor_use_drugs', 'aggressor_used_position_authority',
            'aggressor_violence_physical_increase', 'aggressor_violently_jealous',
            'aggressor_weapon_used', 'aggressors_occupation', 'agrees_route_activation',
            'allows_easy_report', 'assigned_sex_at_birth', 'attention_was_appropriate',
            'campesino_recognition', 'currently_pregnant', 'economically_dependent',
            'employment_relationship', 'ethnic_affiliation', 'facts_zone',
            'gender_identity', 'housing_stratum', 'housing_tenancy_form',
            'identifies_barriers', 'income_generation_method', 'indigenous_people',
            'last_education_level', 'marital_status', 'migration_condition', 'nationality',
            'num_agressors', 'occupation', 'partner_controls',
            'partner_other_denunciations', 'partner_threatened_damage_members',
            'partner_threatened_suicide', 'partner_unemployed', 'person_with_disability',
            'physical_mental_sensory_difficulties', 'proximity_principal_aggressor',
            'recurrence_aggression', 'relationship_with_presumed_aggressor',
            'reported_previously', 'reported_to_entity', 'require_language_interpreter',
            'requires_referral', 'residence_zone', 'salvia_activates_routes',
            'scenario_violence', 'sexual_orientation', 'specified_nationality',
            'still_lives_with_aggressor', 'stopped_seeking_help', 'support_contact_kinship',
            'thoughts_of_self_harm', 'thoughts_of_self_harm_2',
            'threatened_reveal_sexual_orientation', 'tzv3_g1', 'tzv3_g2', 'tzv3_g3',
            'tzv3_g4', 'tzv3_g5', 'tzv3_g6', 'tzv3_g7', 'tzv3_vpp',
            'victim_health_to_blackmail', 'victim_pronoun', 'violence_motivated_by_gender',
            'violence_motivated_by_gender_2', 'wants_imminent_alert',
            'workplace_sector_occurrence'
              ]) AS k)
        || (SELECT jsonb_object_agg('victim_case_form2_' || k, 0)      -- números en 0
              FROM unnest(ARRAY[
            'activities_unable_to_do_everyday', 'activities_unable_to_eat',
            'activities_unable_to_hear', 'activities_unable_to_interact',
            'activities_unable_to_move', 'activities_unable_to_see',
            'activities_unable_to_take', 'activities_unable_to_talk',
            'activities_unable_to_understand', 'aggressor_phone', 'support_contact_phone'
              ]) AS k)
        || (SELECT jsonb_object_agg('victim_case_form2_' || k, ''::text)   -- textos vacíos
              FROM unnest(ARRAY[
            'aggressor_address', 'aggressor_doc_number', 'aggressor_names', 'facts_address',
            'facts_town_code', 'identity_name', 'imminent_alert_note',
            'interinstitutional_actions_note', 'language_assistance',
            'referral_justification', 'residence_address', 'residence_town',
            'saliva_management_explanation', 'support_contact_email',
            'support_contact_names', 'victim_email', 'victim_identity_name'
              ]) AS k)
        || jsonb_build_object(                                          -- fechas y hora marcador
            'victim_case_form2_birth_date', '1900-01-01',
            'victim_case_form2_facts_date', '1900-01-01',
            'victim_case_form2_facts_start_time', '00:00:00');

    FOR c IN
        SELECT vcase.victim_case_id, s.n, s.creado,
               f.victim_contact_form2_facts_description AS hechos,
               f.victim_contact_form2_victim_col_phone  AS telefono
          FROM seed_caso s
          JOIN salvia.victim_case vcase ON vcase.victim_case_i_code = pg_temp.sid(6, s.n)
          JOIN salvia.victim_contact_form2 f ON f.victim_case_form2_i_code = pg_temp.sid(2, s.n)
         WHERE NOT EXISTS (SELECT 1 FROM salvia.victim_case_form2 x
                            WHERE x.victim_case_form2_victim_case = vcase.victim_case_id)
    LOOP
        f2 := tpl || jsonb_build_object(
            'victim_case_form2_i_code', pg_temp.sid(7, c.n),
            'victim_case_form2_creation_date', c.creado,
            'victim_case_form2_update_date', c.creado,
            'victim_case_form2_approx_start_asp', c.creado,
            'victim_case_form2_victim_case', c.victim_case_id,
            'victim_case_form2_facts_description', c.hechos,        -- HV_1
            'victim_case_form2_victim_phone', c.telefono);          -- CI_4

        -- Red de seguridad: columnas NOT NULL sin valor por defecto que la plantilla no conoce.
        FOR m IN
            SELECT column_name, data_type FROM information_schema.columns
             WHERE table_schema = 'salvia' AND table_name = 'victim_case_form2'
               AND is_nullable = 'NO' AND column_default IS NULL AND NOT (f2 ? column_name)
        LOOP
            f2 := f2 || jsonb_build_object(m.column_name,
                CASE WHEN m.data_type LIKE 'timestamp%' THEN to_jsonb(c.creado)
                     WHEN m.data_type = 'date' THEN to_jsonb('1900-01-01'::text)
                     WHEN m.data_type LIKE 'time%' THEN to_jsonb('00:00:00'::text)
                     WHEN m.data_type IN ('bigint', 'integer', 'smallint', 'numeric') THEN to_jsonb(0)
                     WHEN m.data_type = 'boolean' THEN to_jsonb(false)
                     ELSE to_jsonb(''::text) END);
        END LOOP;

        SELECT string_agg(quote_ident(k), ', ') INTO cols FROM jsonb_object_keys(f2) AS k;
        EXECUTE format('INSERT INTO salvia.victim_case_form2 (%1$s) SELECT %1$s FROM jsonb_populate_record(NULL::salvia.victim_case_form2, %2$L::jsonb)', cols, f2);
    END LOOP;
END $$;

-- CI_8 (ajuste razonable): el caso hereda las opciones del reporte, como en la conversión.
INSERT INTO salvia.rel_victim_case_form2_enums_victim_case_form2 (victim_case_form2_enums_id, victim_case_form2_id)
SELECT rel.victim_case_form2_enums_id, f2.victim_case_form2_id
  FROM seed_caso s
  JOIN salvia.victim_case_form2 f2 ON f2.victim_case_form2_i_code = pg_temp.sid(7, s.n)
  JOIN salvia.victim_contact_form2 vf ON vf.victim_case_form2_i_code = pg_temp.sid(2, s.n)
  JOIN salvia.rel_victim_case_form2_enums_victim_contact_form2 rel ON rel.victim_contact_form2_id = vf.victim_contact_form2_id
 WHERE NOT EXISTS (SELECT 1 FROM salvia.rel_victim_case_form2_enums_victim_case_form2 x
                    WHERE x.victim_case_form2_id = f2.victim_case_form2_id
                      AND x.victim_case_form2_enums_id = rel.victim_case_form2_enums_id);

-- ─── Historial de cada reporte (categoría «Reportes» del timeline) ───────
-- Un hito por hecho, con el mismo texto, ícono y color que escribe la app: recibido (actor: quien
-- reporta con sesión), asignado (por cloud.test.sv), intentos fallidos y contactos logrados,
-- borrador de caso iniciado y cierre como inválido. El cierre válido no tiene hito propio: lo
-- deriva el estado del caso. Id = 5eed0005-…: n * 1000 + el orden del hito.
INSERT INTO salvia.case_timeline_event (
    id, event_type, description, category, type, icon, color, date, created_at,
    event_user_id, victim_contact_id, case_id,
    task_id, actor_id, actor_name, barrier_id, follow_up_id, entity_letter_id,
    emergency_measure_id, psychosocial_support_id, articulating_mechanism_id, economic_stabilization_id)
SELECT pg_temp.sid(5, e.n * 1000 + e.ord), e.event_type, e.descr, 'Reportes', e.tipo, e.icon, e.color, e.at, e.at,
       e.actor, pg_temp.sid(1, e.n), e.case_icode,
       '', '', '', '', '', '', '', '', '', ''
  FROM (
    -- Reporte recibido
    SELECT r.n, 1 AS ord, 'REPORTE_RECIBIDO' AS event_type, 'Reporte recibido' AS tipo, 'inbox' AS icon, '#1d4ed8' AS color,
           'Reporte de ' || t.etiqueta || ' recibido. Radicado ' || r.radicado || '.' AS descr,
           COALESCE('user-icode-' || r.reporta, '') AS actor,
           r.creado AT TIME ZONE 'America/Bogota' AS at, NULL::text AS case_icode
      FROM seed_reporte r JOIN seed_tipo t ON t.tipo = r.tipo
    UNION ALL
    -- Reporte asignado: una hora después de radicarse
    SELECT r.n, 2, 'REPORTE_ASIGNADO', 'Reporte asignado', 'user-check', '#0f766e',
           'Asignado a ' || (SELECT trim(gup.general_user_profile_names || ' ' || COALESCE(gup.general_user_profile_last_names, ''))
                               FROM security.general_user gu
                               JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
                              WHERE gu.general_user_login = r.agente)
                     || ' por ' || (SELECT trim(gup.general_user_profile_names || ' ' || COALESCE(gup.general_user_profile_last_names, ''))
                                      FROM security.general_user gu
                                      JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
                                     WHERE gu.general_user_login = 'cloud.test.sv') || '.',
           'user-icode-cloud.test.sv',
           (r.creado AT TIME ZONE 'America/Bogota') + interval '1 hour', NULL
      FROM seed_reporte r WHERE r.agente IS NOT NULL
    UNION ALL
    -- Intento fallido / contacto logrado (orden 10 + k)
    SELECT i.n, 10 + i.k,
           CASE i.resultado WHEN 'fallido' THEN 'INTENTO_FALLIDO' ELSE 'CONTACTO_LOGRADO' END,
           CASE i.resultado WHEN 'fallido' THEN 'Intento de contacto fallido' ELSE 'Contacto logrado' END,
           CASE i.resultado WHEN 'fallido' THEN 'phone-slash' ELSE 'handshake' END,
           CASE i.resultado WHEN 'fallido' THEN '#c2410c' ELSE '#15803d' END,
           CASE i.resultado WHEN 'fallido' THEN 'Intento fallido registrado (' || i.n_dia || ' hoy).'
                ELSE 'Contacto logrado el ' || to_char(i.en AT TIME ZONE 'America/Bogota', 'DD/MM/YYYY HH24:MI') || '.' END,
           'user-icode-' || r.agente, i.en, NULL
      FROM seed_intento i JOIN seed_reporte r ON r.n = i.n
    UNION ALL
    -- Borrador de caso iniciado (cuelga del reporte y del caso)
    SELECT s.n, 91, 'BORRADOR_CASO_INICIADO', 'Borrador de caso iniciado', 'folder-open', '#5106A7',
           'Se inició el borrador del caso ' || pg_temp.sid(6, s.n) || ' a partir del reporte.',
           'user-icode-cloud.test.op', s.creado AT TIME ZONE 'America/Bogota', pg_temp.sid(6, s.n)
      FROM seed_caso s
    UNION ALL
    -- Reporte cerrado como inválido: lo cierra el agente asignado (o el supervisor si no hay)
    SELECT r.n, 90, 'REPORTE_CERRADO_INVALIDO', 'Reporte cerrado como inválido', 'circle-xmark', '#b91c1c',
           'Motivo: ' || cl.etiqueta || '.' || COALESCE(' ' || r.cierre_obs, ''),
           COALESCE('user-icode-' || r.agente, 'user-icode-cloud.test.sv'), r.cerrado_en, NULL
      FROM seed_reporte r JOIN seed_cierre cl ON cl.code = r.cierre
  ) e
ON CONFLICT (id) DO NOTHING;

DROP TABLE seed_intento, seed_caso, seed_cierre, seed_tipo, seed_reporte;
