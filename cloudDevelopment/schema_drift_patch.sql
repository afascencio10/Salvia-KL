-- Parche de ESQUEMA para la BD local del sandbox cloud: cierra el drift entre los
-- dumps legacy (src/security/scripts + src/salvia/scripts db_postgres_creation.sql)
-- y la BD real de pruebas (Supabase), que se fue ajustando a mano con el tiempo.
-- Sin este parche la BD recién creada NO coincide con lo que el código espera y
-- varias pantallas responden 500 ("column ... does not exist", "relation ... does
-- not exist"): Lista de casos, Casos Entidad, Mecanismos, Mis barreras, remisiones
-- psicosociales y de medidas de emergencia.
--
-- Cómo se generó (2026-11-01): se comparó el catálogo de esquema de Supabase
-- (pg_catalog: columnas, constraints, índices — solo metadatos, ninguna fila de
-- datos), leído vía Composio, contra la BD local ya creada por
-- setup-local-postgres.sh y con la app arrancada. Para cerrar un drift nuevo, ver
-- "Esquema: cuando aparece 'column ... does not exist'" en
-- actualizar-seed-con-composio.md.
--
-- Idempotente (IF NOT EXISTS / guardas). Se aplica DESPUÉS de los dos dumps legacy
-- y ANTES de los seeds, así que solo toca tablas que existen en ese momento (las
-- que crea la app con GORM aparecen después y ya coinciden con Supabase).
--
-- NO se replica, a propósito:
--   * Claves foráneas. Las columnas nuevas de victim_case_form2 apuntan a
--     victim_case_form2_enums, que local está incompleta (790 de 912 filas): la FK
--     haría fallar cualquier dato de prueba que use una de las opciones faltantes.
--   * victim_case_barrier y victim_case_entity: existen en Supabase pero el código
--     Go no las usa.
--   * salvia.entity.entity_city_code (NOT NULL en Supabase): solo la usan los
--     controllers de migración, y exigirla rompería el seed de referencia.
--   * Índices únicos parciales de entity_case y follow_up_v2 (tablas de GORM, que
--     aún no existen cuando corre este parche).
--   * Las 8 tablas feminicide* que solo existen en los dumps legacy (Supabase ya
--     no las tiene) y las diferencias donde local es MÁS ancho que Supabase.

-- 1) security.email: en la BD real la PK se llama e_mail_id (el dump legacy dice email_id).
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'security' AND table_name = 'email' AND column_name = 'email_id')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'security' AND table_name = 'email' AND column_name = 'e_mail_id') THEN
        ALTER TABLE security.email RENAME COLUMN email_id TO e_mail_id;
    END IF;
END $$;

-- 2) Tablas que existen en Supabase y no en los dumps legacy (las usa el código Go).
CREATE TABLE IF NOT EXISTS salvia.follow_up_attempts (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    follow_up_id character varying(36) NOT NULL,
    reason text,
    was_answered boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT follow_up_attempts_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_follow_up_attempts_deleted_at ON salvia.follow_up_attempts USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_follow_up_attempts_follow_up_id ON salvia.follow_up_attempts USING btree (follow_up_id);

CREATE SEQUENCE IF NOT EXISTS salvia.victim_contact_attempt_victim_contact_attempt_id_seq;
CREATE TABLE IF NOT EXISTS salvia.victim_contact_attempt (
    victim_contact_attempt_id integer NOT NULL DEFAULT nextval('salvia.victim_contact_attempt_victim_contact_attempt_id_seq'::regclass),
    victim_contact_attempt_i_code character varying(36) NOT NULL,
    victim_contact_attempt_victim_contact character varying(36) NOT NULL,
    victim_contact_attempt_result character varying(10) NOT NULL,
    victim_contact_attempt_date timestamp without time zone NOT NULL,
    victim_contact_attempt_agent character varying(64) NOT NULL,
    victim_contact_attempt_creation_date timestamp without time zone NOT NULL,
    victim_contact_attempt_deleted_at timestamp without time zone,
    CONSTRAINT victim_contact_attempt_i_code_key UNIQUE (victim_contact_attempt_i_code),
    CONSTRAINT victim_contact_attempt_pkey PRIMARY KEY (victim_contact_attempt_id)
);
ALTER SEQUENCE salvia.victim_contact_attempt_victim_contact_attempt_id_seq OWNED BY salvia.victim_contact_attempt.victim_contact_attempt_id;

CREATE SEQUENCE IF NOT EXISTS salvia.victim_contact_form1_victim_contact_form1_id_seq;
CREATE TABLE IF NOT EXISTS salvia.victim_contact_form1 (
    victim_contact_form1_id integer NOT NULL DEFAULT nextval('salvia.victim_contact_form1_victim_contact_form1_id_seq'::regclass),
    victim_contact_form1_i_code character varying(36),
    victim_contact_form1_creation_date timestamp without time zone NOT NULL,
    victim_contact_form1_update_date timestamp without time zone NOT NULL,
    victim_contact_form1_general_user character varying(36),
    victim_contact_form1_nick character varying(32),
    victim_contact_form1_doc_type character varying(2),
    victim_contact_form1_doc_number character varying(32),
    victim_contact_form1_birth_date timestamp without time zone,
    victim_contact_form1_town_code character varying(8),
    victim_contact_form1_address character varying(120),
    victim_contact_form1_phone character varying(10),
    victim_contact_form1_gender_identity character varying(2),
    victim_contact_form1_sexual_orientation character varying(2),
    victim_contact_form1_origin character varying(1),
    victim_contact_form1_occupation character varying(2),
    victim_contact_form1_occupation_other character varying(32),
    victim_contact_form1_facts_description text,
    victim_contact_form1_victim_contact bigint NOT NULL,
    CONSTRAINT victim_contact_form1_pkey PRIMARY KEY (victim_contact_form1_id),
    CONSTRAINT victim_contact_form1_victim_contact_form1_i_code_key UNIQUE (victim_contact_form1_i_code)
);
ALTER SEQUENCE salvia.victim_contact_form1_victim_contact_form1_id_seq OWNED BY salvia.victim_contact_form1.victim_contact_form1_id;

CREATE SEQUENCE IF NOT EXISTS salvia.victim_contact_form2_v2_victim_contact_form2_v2_id_seq;
CREATE TABLE IF NOT EXISTS salvia.victim_contact_form2_v2 (
    victim_contact_form2_v2_id integer NOT NULL DEFAULT nextval('salvia.victim_contact_form2_v2_victim_contact_form2_v2_id_seq'::regclass),
    victim_contact_form2_v2_i_code character varying(36) NOT NULL,
    victim_contact_form2_v2_creation_date timestamp without time zone NOT NULL,
    victim_contact_form2_v2_update_date timestamp without time zone NOT NULL,
    victim_contact_form2_v2_victim_contact_form2 bigint NOT NULL,
    victim_contact_form2_v2_authorization_answer bigint NOT NULL,
    victim_contact_form2_v2_anonimato bigint,
    victim_contact_form2_v2_parentesco bigint,
    victim_contact_form2_v2_reporter_last_names character varying(64),
    victim_contact_form2_v2_reporter_email character varying(120),
    victim_contact_form2_v2_reporter_will_receive_call bigint,
    victim_contact_form2_v2_reporter_best_contact_time time without time zone,
    CONSTRAINT victim_contact_form2_v2_form2_key UNIQUE (victim_contact_form2_v2_victim_contact_form2),
    CONSTRAINT victim_contact_form2_v2_i_code_key UNIQUE (victim_contact_form2_v2_i_code),
    CONSTRAINT victim_contact_form2_v2_pkey PRIMARY KEY (victim_contact_form2_v2_id)
);
ALTER SEQUENCE salvia.victim_contact_form2_v2_victim_contact_form2_v2_id_seq OWNED BY salvia.victim_contact_form2_v2.victim_contact_form2_v2_id;

-- 3) Columnas que existen en Supabase y no en los dumps legacy.
ALTER TABLE salvia.follow_up_entry_acting
    ADD COLUMN IF NOT EXISTS follow_up_entry_id bigint;
ALTER TABLE salvia.rel_alert_victim_case
    ADD COLUMN IF NOT EXISTS rel_alert_victim_case_victim_moment_id bigint;
ALTER TABLE salvia.rel_case_owner_victim_case
    ADD COLUMN IF NOT EXISTS rel_case_owner_victim_case_status character varying(1) NOT NULL DEFAULT 'a'::character varying;
ALTER TABLE salvia.victim_case
    ADD COLUMN IF NOT EXISTS victim_case_victim_doc_type character varying(2),
    ADD COLUMN IF NOT EXISTS victim_case_victim_town_code character varying(8) NOT NULL,
    ADD COLUMN IF NOT EXISTS victim_case_form_submission_id character varying(36),
    ADD COLUMN IF NOT EXISTS victim_case_new_user_credentials jsonb,
    ADD COLUMN IF NOT EXISTS victim_case_v3_answered_keys text;
ALTER TABLE salvia.victim_case_form2
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g1 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g2 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g3 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g4 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g5 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g6 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_g7 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_vpp integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p1 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np1 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p2 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np2 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p3 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np3 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p4 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np4 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p5 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np5 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p6 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np6 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p7 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np7 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p8 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np8 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p9 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np9 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p10 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np10 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p11 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np11 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p12 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np12 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p13 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np13 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p14 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np14 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p15 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np15 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p16 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np16 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p17 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np17 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p18 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np18 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p19 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np19 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p20 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np20 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p21 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np21 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p22 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np22 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p23 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np23 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_p24 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_np24 integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_risk_type integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_risk_level integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_risk_score integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_tzv3_risk_tie_flag integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_requires_referral integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_referral_justification character varying(2000),
    ADD COLUMN IF NOT EXISTS victim_case_form2_reported_to_entity integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_salvia_activates_routes integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_agrees_route_activation integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_identifies_barriers integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_wants_imminent_alert integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_imminent_alert_note character varying(1000),
    ADD COLUMN IF NOT EXISTS victim_case_form2_interinstitutional_actions_note character varying(2000),
    ADD COLUMN IF NOT EXISTS victim_case_form2_victim_email character varying(128),
    ADD COLUMN IF NOT EXISTS victim_case_form2_victim_pronoun integer,
    ADD COLUMN IF NOT EXISTS victim_case_form2_victim_identity_name character varying(256);
ALTER TABLE salvia.victim_contact
    ADD COLUMN IF NOT EXISTS victim_contact_latitude double precision,
    ADD COLUMN IF NOT EXISTS victim_contact_longitude double precision;

-- 4) Columnas que local dejó más angostas que Supabase (p. ej. un i_code de 36 no cabe en 32).
ALTER TABLE salvia.case_owner ALTER COLUMN case_owner_general_user TYPE character varying(36);  -- local: character varying(32)
ALTER TABLE salvia.victim_case ALTER COLUMN victim_case_owner_description TYPE character varying(1024);  -- local: character varying(512)
ALTER TABLE salvia.victim_case_form2 ALTER COLUMN victim_case_form2_facts_address TYPE character varying(150);  -- local: character varying(32)
ALTER TABLE salvia.victim_case_form2 ALTER COLUMN victim_case_form2_aggressor_address TYPE character varying(150);  -- local: character varying(32)
ALTER TABLE salvia.victim_case_form2 ALTER COLUMN victim_case_form2_residence_address TYPE character varying(150);  -- local: character varying(32)

-- 4b) Nulabilidad y defaults: columnas que local dejó NOT NULL (o sin default) y que en
--     Supabase aceptan NULL (o traen default). Sin esto, crear un caso nuevo, que aún no
--     está aprobado, falla con violación de NOT NULL.
ALTER TABLE salvia.follow_up_entry_acting ALTER COLUMN follow_up_entry_acting_rel_barrier_follow_up_entry DROP NOT NULL;
ALTER TABLE salvia.moment ALTER COLUMN moment_approval_owner DROP NOT NULL;
ALTER TABLE salvia.moment ALTER COLUMN moment_victim_case DROP NOT NULL;
ALTER TABLE salvia.victim_case ALTER COLUMN victim_case_approved_by DROP NOT NULL;
ALTER TABLE salvia.case_owner ALTER COLUMN case_owner_num_cases SET DEFAULT 0;
ALTER TABLE salvia.entity_branch ALTER COLUMN entity_branch_source SET DEFAULT 's'::character varying;
ALTER TABLE salvia.rel_case_owner_victim_case ALTER COLUMN rel_case_owner_victim_case_creation_date SET DEFAULT now();
ALTER TABLE salvia.victim_case ALTER COLUMN victim_case_owner_description SET DEFAULT '...'::character varying;

-- 5) Restricción e índice que existen en Supabase y no en los dumps legacy.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'moment_moment_i_code_key') THEN
        ALTER TABLE salvia.moment ADD CONSTRAINT moment_moment_i_code_key UNIQUE (moment_i_code);
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_salvia_victim_case_form_submission_id ON salvia.victim_case USING btree (victim_case_form_submission_id);
