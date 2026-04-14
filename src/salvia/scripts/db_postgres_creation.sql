-- Schema de seguridad
SET search_path = salvia;

CREATE TABLE victim_contact (
   victim_contact_id SERIAL NOT NULL,
   victim_contact_i_code CHARACTER VARYING(36) NOT NULL,
   victim_contact_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_contact_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_contact_general_user CHARACTER VARYING(36) NOT NULL,
   victim_contact_names CHARACTER VARYING(32) ,
   victim_contact_last_names CHARACTER VARYING(32) ,
   victim_contact_nick CHARACTER VARYING(32) ,
   victim_contact_doc_type CHARACTER VARYING(2) ,
   victim_contact_doc_number CHARACTER VARYING(32) ,
   victim_contact_birth_date TIMESTAMP WITHOUT TIME ZONE ,
   victim_contact_living_zone CHARACTER VARYING(2) ,
   victim_contact_address CHARACTER VARYING(120) ,
   victim_contact_living_latitude DOUBLE PRECISION ,
   victim_contact_living_longitude DOUBLE PRECISION ,
   victim_contact_phone CHARACTER VARYING(10) ,
   victim_contact_gender_identity CHARACTER VARYING(2) ,
   victim_contact_sexual_orientation CHARACTER VARYING(2) ,
   victim_contact_origin CHARACTER VARYING(1) ,
   victim_contact_occupation CHARACTER VARYING(2) ,
   victim_contact_occupation_other CHARACTER VARYING(32) ,
   victim_contact_status CHARACTER VARYING(1) NOT NULL,
   victim_contact_status_description CHARACTER VARYING(256) ,
   victim_contact_facts_description TEXT,
   CONSTRAINT victim_contact_pkey PRIMARY KEY (victim_contact_id),
   CONSTRAINT victim_contact_victim_contact_i_code_key UNIQUE (victim_contact_i_code)
);

CREATE TABLE victim_case (
   victim_case_id SERIAL NOT NULL,
   victim_case_i_code CHARACTER VARYING(36) NOT NULL,
   victim_case_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_case_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_case_status CHARACTER VARYING(2) NOT NULL,
   victim_case_general_user CHARACTER VARYING(36) NOT NULL,
   victim_case_victim_names CHARACTER VARYING(32) ,
   victim_case_victim_last_names CHARACTER VARYING(32) ,
   victim_case_form2_victim_doc_type CHARACTER VARYING(2) ,
   victim_case_victim_doc_number CHARACTER VARYING(32) ,
   victim_case_approved_by BIGINT NOT NULL,
   victim_case_victim_contact BIGINT ,
   victim_case_follow_up BIGINT ,
   
   victim_case_owner_description CHARACTER VARYING(512) ,

   
   CONSTRAINT victim_case_pkey PRIMARY KEY (victim_case_id),
   CONSTRAINT victim_case_victim_case_i_code_key UNIQUE (victim_case_i_code)
);

CREATE TABLE salvia.victim_case_form1
(
   victim_case_form1_id SERIAL NOT NULL,
   victim_case_form1_i_code character varying(36) NOT NULL,
   victim_case_form1_creation_date timestamp without time zone NOT NULL,
   victim_case_form1_update_date timestamp without time zone NOT NULL,
   victim_case_form1_victim_nick character varying(32) ,
   victim_case_form1_victim_birth_date timestamp without time zone NOT NULL,
   victim_case_form1_victim_address character varying(120) NOT NULL,
   victim_case_form1_victim_living_latitude double precision,
   victim_case_form1_victim_living_longitude double precision,
   victim_case_form1_victim_phone character varying(10) NOT NULL,
   victim_case_form1_victim_gender_identity character varying(2) NOT NULL,
   victim_case_form1_victim_sexual_orientation character varying(2) NOT NULL,
   victim_case_form1_victim_origin character varying(1) NOT NULL,
   victim_case_form1_victim_occupation character varying(2) NOT NULL,
   victim_case_form1_victim_occupation_other character varying(32) ,
   victim_case_form1_victim_ethnic_group character varying(2) NOT NULL,
   victim_case_form1_victim_ethnic_group_other character varying(32) ,
   victim_case_form1_victim_contact_names character varying(32) NOT NULL,
   victim_case_form1_victim_contact_phone character varying(10) NOT NULL,
   victim_case_form1_victim_contact_kinship character varying(32) NOT NULL,
   victim_case_form1_victim_children_number BIGINT NOT NULL,
   victim_case_form1_victim_marital_status character varying(2) NOT NULL,
   victim_case_form1_victim_marital_status_other character varying(32) NOT NULL,
   victim_case_form1_victim_children_age character varying(32) ,
   victim_case_form1_victim_disability character varying(2) NOT NULL,
   victim_case_form1_victim_special_support character varying(128) ,
   victim_case_form1_facts_occurrence character varying(2) NOT NULL,
   victim_case_form1_facts_weekday bigint,
   victim_case_form1_facts_date timestamp without time zone NOT NULL,
   victim_case_form1_facts_description text NOT NULL,
   victim_case_form1_victim_violence_experienced character varying(2) NOT NULL,
   victim_case_form1_victim_violence_experienced_other character varying(32),
   victim_case_form1_victim_violence_scope character varying(2) NOT NULL,
   victim_case_form1_victim_femicide_risk character varying(1) NOT NULL,
   victim_case_form1_victim_aggressor character varying(1) NOT NULL,
   victim_case_form1_victim_relationship_with_aggressor character varying(2),
   victim_case_form1_victim_aggressor_name character varying(64),
   victim_case_form1_victim_aggressor_doc_type character varying(2),
   victim_case_form1_victim_aggressor_doc_number character varying(32),
   victim_case_form1_victim_aggressor_address character varying(128),
   victim_case_form1_victim_aggressor_phone character varying(10),
   victim_case_form1_facts_start_time time without time zone,
   victim_case_form1_facts_end_time time without time zone,
   victim_case_form1_victim_e_mail character varying(128),
   victim_case_form1_age bigint,
   victim_case_form1_victim_nationality character varying(1),
   victim_case_form1_victim_nationality_other character varying(32),
   victim_case_form1_victim_foreigner_immigration_status character varying(1),
   victim_case_form1_victim_gender character varying(1),
   victim_case_form1_victim_gender_identity_other character varying(32),
   victim_case_form1_victim_sexual_orientation_other character varying(32),
   victim_case_form1_victim_dependents character varying(1),
   victim_case_form1_victim_death_threats character varying(1),
   victim_case_form1_victim_aggressor_has_weapons character varying(1),
   victim_case_form1_experienced_physical_or_sexual_violence_befor character varying(1),
   victim_case_form1_victim_previously_reported_situation character varying(1),
   victim_case_form1_victim_if_previously_reported character varying(2),
   victim_case_form1_victim_imminent_risk character varying(1),
   victim_case_form1_violence_town_code character varying(8),
   victim_case_form1_victim_if_afro character varying(1),
   victim_case_form1_victim_if_indigenous character varying(2),
   victim_case_form1_victim_if_indigenous_tongue character varying(2),
   victim_case_form1_victim_if_peasant character varying(1),
   victim_case_form1_victim_if_armed_conflict character varying(1),
   victim_case_form1_victim_violence_scene character varying(2),
   victim_case_form1_physical_violence_increased character varying(1),
   victim_case_form1_separated_from_partner_last_year character varying(1),
   victim_case_form1_threatened_with_weapon character varying(1),
   victim_case_form1_threatened_to_kill_or_harm_children character varying(1),
   victim_case_form1_jealous_and_violent character varying(1),
   victim_case_form1_believes_capable_of_killing character varying(1),
   
   victim_case_form1_victim_case BIGINT NOT NULL,
   CONSTRAINT victim_case_form1_pkey PRIMARY KEY (victim_case_form1_id),
   CONSTRAINT victim_case_form1_victim_case_form1_i_code_key UNIQUE (victim_case_form1_i_code)
)


CREATE TABLE case_owner (
   case_owner_id SERIAL NOT NULL,
   case_owner_i_code CHARACTER VARYING(36) NOT NULL,
   case_owner_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   case_owner_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   case_owner_general_user CHARACTER VARYING(32) NOT NULL,
   case_owner_num_cases INTEGER NOT NULL,
   attention_line_id BIGINT,
   entity_branch_id BIGINT,
   CONSTRAINT case_owner_pkey PRIMARY KEY (case_owner_id),
   CONSTRAINT case_owner_case_owner_i_code_key UNIQUE (case_owner_i_code)
);

CREATE TABLE entity_branch (
   entity_branch_id SERIAL NOT NULL,
   entity_branch_i_code CHARACTER VARYING(36) NOT NULL,
   entity_branch_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   entity_branch_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   entity_branch_name CHARACTER VARYING(254) NOT NULL,
   entity_branch_description CHARACTER VARYING(512) ,
   entity_branch_address CHARACTER VARYING(254) NOT NULL,
   entity_branch_latitude DOUBLE PRECISION ,
   entity_branch_longitude DOUBLE PRECISION ,
   entity_branch_town_code CHARACTER VARYING(8) NOT NULL,
   entity_branch_source CHARACTER VARYING(2) NOT NULL,
   entity_branch_owner_general_user CHARACTER VARYING(36),
   entity_id BIGINT ,
   CONSTRAINT entity_branch_pkey PRIMARY KEY (entity_branch_id),
   CONSTRAINT entity_branch_entity_branch_i_code_key UNIQUE (entity_branch_i_code)

);

CREATE TABLE entity (
   entity_id SERIAL NOT NULL,
   entity_i_code CHARACTER VARYING(36) NOT NULL,
   entity_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   entity_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   entity_name CHARACTER VARYING(254) NOT NULL,
   entity_description CHARACTER VARYING(512) NOT NULL,
   entity_is_interoperable CHARACTER VARYING(1) NOT NULL,
   entity_interoperability_code CHARACTER VARYING(32) ,
   entity_response_time BIGINT NOT NULL,
   entity_start_edu_content TEXT NOT NULL,
   entity_fail_edu_content TEXT NOT NULL,
   
   entity_sector CHARACTER VARYING(2) NOT NULL,
   CONSTRAINT entity_pkey PRIMARY KEY (entity_id),
   CONSTRAINT entity_entity_i_code_key UNIQUE (entity_i_code)
);

CREATE TABLE rel_entity_moment (
   rel_entity_moment_id SERIAL NOT NULL,
   rel_entity_moment_code CHARACTER VARYING(2) NOT NULL,
   rel_entity_moment_entity BIGINT NOT NULL,
   CONSTRAINT rel_entity_moment_pkey PRIMARY KEY (rel_entity_moment_id)
);

CREATE TABLE rel_case_owner_victim_case (
   rel_case_owner_victim_case_id SERIAL NOT NULL,
   case_owner_id BIGINT NOT NULL,
   victim_case_id BIGINT NOT NULL,
   rel_case_owner_victim_case_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   CONSTRAINT rel_case_owner_victim_case_pkey PRIMARY KEY (rel_case_owner_victim_case_id)
);



CREATE TABLE attention_line (
   attention_line_id SERIAL NOT NULL,
   attention_line_i_code CHARACTER VARYING(36) NOT NULL,
   attention_line_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   attention_line_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   attention_line_name CHARACTER VARYING(254) NOT NULL,
   attention_line_description CHARACTER VARYING(512) NOT NULL,
   CONSTRAINT attention_line_pkey PRIMARY KEY (attention_line_id),
   CONSTRAINT attention_line_attention_line_i_code_key UNIQUE (attention_line_i_code)
);

CREATE TABLE moment (
   moment_id SERIAL NOT NULL,
   moment_i_code CHARACTER VARYING(36) NOT NULL,
   moment_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   moment_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   moment_approval_cancellation CHARACTER VARYING(1) NOT NULL,
   moment_approval_source CHARACTER VARYING(1) ,
   moment_victim_case BIGINT NOT NULL,
   moment_approval_owner BIGINT NOT NULL,
   moment_entity_branch BIGINT NOT NULL,
   moment_code CHARACTER VARYING(2) NOT NULL,
   moment_approval_description TEXT,
   moment_status CHARACTER VARYING(1) NOT NULL,
   CONSTRAINT moment_pkey PRIMARY KEY (moment_id)
);

CREATE TABLE case_log (
   case_log_id SERIAL NOT NULL,
   case_log_i_code CHARACTER VARYING(36) NOT NULL,
   case_log_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   case_log_description TEXT NOT NULL,
   case_log_general_user CHARACTER VARYING(36) NOT NULL,
   moment_id BIGINT ,
   CONSTRAINT case_log_pkey PRIMARY KEY (case_log_id),
   CONSTRAINT case_log_case_log_i_code_key UNIQUE (case_log_i_code)
);

CREATE TABLE alert (
   alert_id SERIAL NOT NULL,
   alert_i_code CHARACTER VARYING(36) NOT NULL,
   alert_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   alert_type CHARACTER VARYING(1) NOT NULL,
   alert_priority INTEGER NOT NULL,
   alert_code CHARACTER VARYING(64) NOT NULL,
   CONSTRAINT alert_pkey PRIMARY KEY (alert_id),
   CONSTRAINT alert_alert_i_code_key UNIQUE (alert_i_code)

);

CREATE TABLE rel_alert_victim_case (
   rel_alert_victim_case_id SERIAL NOT NULL,
   rel_alert_victim_case_alert_id BIGINT NOT NULL,
   rel_alert_victim_case_victim_case_id BIGINT NOT NULL,
   rel_alert_victim_case_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   rel_alert_victim_case_data TEXT,
   CONSTRAINT rel_alert_victim_case_pkey PRIMARY KEY (rel_alert_victim_case_id)
);

ALTER TABLE victim_case_form1 ADD CONSTRAINT victim_case_form1_victim_case
   FOREIGN KEY ( victim_case_form1_victim_case) REFERENCES victim_case(victim_case_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE victim_case ADD CONSTRAINT victim_case_follow_up
   FOREIGN KEY ( victim_case_follow_up) REFERENCES follow_up(follow_up_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_case ADD CONSTRAINT victim_case_victim_contact
   FOREIGN KEY ( victim_case_victim_contact) REFERENCES victim_contact(victim_contact_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE case_log ADD CONSTRAINT case_log_moment
   FOREIGN KEY ( moment_id) REFERENCES moment(moment_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_alert_victim_case ADD CONSTRAINT rel_victim_case_id
   FOREIGN KEY ( rel_alert_victim_case_victim_case_id) REFERENCES victim_case(victim_case_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_alert_victim_case ADD CONSTRAINT rel_alert_id
   FOREIGN KEY ( rel_alert_victim_case_alert_id) REFERENCES alert(alert_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE moment ADD CONSTRAINT moment_victim_case
   FOREIGN KEY ( moment_victim_case) REFERENCES case(case_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
ALTER TABLE moment ADD CONSTRAINT moment_approval_owner
   FOREIGN KEY ( moment_approval_owner) REFERENCES case_owner(case_owner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;ALTER TABLE moment ADD CONSTRAINT moment_entity_branch
   FOREIGN KEY ( moment_entity_branch) REFERENCES entity_branch(entity_branch_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE case_owner ADD CONSTRAINT case_owner_victim_case
   FOREIGN KEY ( victim_case_id) REFERENCES victim_case(victim_case_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE case_owner ADD CONSTRAINT case_owner_attention_line
   FOREIGN KEY ( attention_line_id) REFERENCES attention_line(attention_line_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE case_owner ADD CONSTRAINT case_owner_entity_branch
   FOREIGN KEY ( entity_branch_id) REFERENCES entity_branch(entity_branch_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE entity_branch ADD CONSTRAINT entity_branch_entity
   FOREIGN KEY ( entity_id) REFERENCES entity(entity_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE victim_case ADD CONSTRAINT victim_case_approved_by
   FOREIGN KEY ( victim_case_approved_by) REFERENCES case_owner(case_owner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_entity_moment ADD CONSTRAINT rel_entity_moment_entity
   FOREIGN KEY ( rel_entity_moment_entity) REFERENCES entity(entity_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_case_owner_victim_case ADD CONSTRAINT rel_victim_case_id
   FOREIGN KEY ( victim_case_id) REFERENCES victim_case(victim_case_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_case_owner_victim_case ADD CONSTRAINT rel_case_owner_id
   FOREIGN KEY ( case_owner_id) REFERENCES case_owner(case_owner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;





-- Scripts posteriores para la creación de seguimiento:

SET search_path = salvia;


CREATE TABLE rel_alert_follow_up (
   rel_alert_follow_up_id SERIAL NOT NULL,
   rel_alert_follow_up_data CHARACTER VARYING(128) NOT NULL,
   rel_alert_follow_up_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   rel_alert_follow_up_follow_up BIGINT NOT NULL,
   rel_alert_follow_up_alert BIGINT NOT NULL,
   CONSTRAINT rel_alert_follow_up_pkey PRIMARY KEY (rel_alert_follow_up_id)
);

CREATE TABLE follow_up_entry (
   follow_up_entry_id SERIAL NOT NULL,
   follow_up_entry_i_code CHARACTER VARYING(36) NOT NULL,
   follow_up_entry_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_entry_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_entry_completion_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_entry_was_done CHARACTER VARYING(1) ,
   follow_up_entry_sector CHARACTER VARYING(2) ,
   follow_up_entry_person_visited_entity CHARACTER VARYING(1) ,
   follow_up_entry_received_attention CHARACTER VARYING(1) ,
   follow_up_entry_comments CHARACTER VARYING(5000) ,
   follow_up_entry_case_documents_prepared CHARACTER VARYING(1) ,
   follow_up_entry_owner_general_user CHARACTER VARYING(36) NOT NULL,
   follow_up_entry_status CHARACTER VARYING(1) NOT NULL,
   follow_up_id BIGINT ,
   CONSTRAINT follow_up_entry_pkey PRIMARY KEY (follow_up_entry_id),
   CONSTRAINT follow_up_entry_follow_up_entry_i_code_key UNIQUE (follow_up_entry_i_code)
);

CREATE TABLE barrier (
   barrier_id SERIAL NOT NULL,
   barrier_i_code CHARACTER VARYING(36) NOT NULL,
   barrier_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   barrier_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   barrier_name CHARACTER VARYING(128) NOT NULL,
   barrier_description CHARACTER VARYING(500) ,
   sector_barrier_id BIGINT,
   CONSTRAINT barrier_pkey PRIMARY KEY (barrier_id),
   CONSTRAINT barrier_barrier_i_code_key UNIQUE (barrier_i_code)
);

CREATE TABLE sector_barrier (
   sector_barrier_id SERIAL NOT NULL,
   sector_barrier_i_code CHARACTER VARYING(36) NOT NULL,
   sector_barrier_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   sector_barrier_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   sector_barrier_name CHARACTER VARYING(36) NOT NULL,
   sector_barrier_description CHARACTER VARYING(500) ,
   sector_barrier_sector CHARACTER VARYING(2) ,
   CONSTRAINT sector_barrier_pkey PRIMARY KEY (sector_barrier_id),
   CONSTRAINT sector_barrier_sector_barrier_i_code_key UNIQUE (sector_barrier_i_code)
);


CREATE TABLE rel_barrier_follow_up_entry (
   rel_barrier_follow_up_entry_id SERIAL NOT NULL,
   rel_barrier_follow_up_entry_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   rel_barrier_follow_up_entry_barrier BIGINT NOT NULL,
   rel_barrier_follow_up_entry_follow_up_entry BIGINT NOT NULL,
   CONSTRAINT rel_barrier_follow_up_entry_pkey PRIMARY KEY (rel_barrier_follow_up_entry_id)
);

CREATE TABLE follow_up_entry_acting (
   follow_up_entry_acting_id SERIAL NOT NULL,
   follow_up_entry_acting_i_code CHARACTER VARYING(36) NOT NULL,
   follow_up_entry_acting_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_entry_acting_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_entry_acting_owner_general_user CHARACTER VARYING(36) NOT NULL,
   follow_up_entry_acting_description CHARACTER VARYING(1000) NOT NULL,
   follow_up_entry_acting_status CHARACTER VARYING(1) NOT NULL,
   follow_up_entry_acting_rel_barrier_follow_up_entry BIGINT NOT NULL,
   CONSTRAINT follow_up_entry_acting_pkey PRIMARY KEY (follow_up_entry_acting_id),
   CONSTRAINT follow_up_entry_acting_follow_up_entry_acting_i_code_key UNIQUE (follow_up_entry_acting_i_code)
);

CREATE TABLE follow_up (
   follow_up_id SERIAL NOT NULL,
   follow_up_i_code CHARACTER VARYING(36) NOT NULL,
   follow_up_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   follow_up_physical_violence_witnessed_by_family CHARACTER VARYING(1) NOT NULL,
   follow_up_violence_escalation CHARACTER VARYING(1) NOT NULL,
   follow_up_assault_with_weapon CHARACTER VARYING(1) NOT NULL,
   follow_up_recent_controlling_or_jealous_behavior CHARACTER VARYING(1) NOT NULL,
   follow_up_violence_history_with_ex_partner CHARACTER VARYING(1) NOT NULL,
   follow_up_violence_history_with_others CHARACTER VARYING(1) NOT NULL,
   follow_up_substance_abuse CHARACTER VARYING(1) NOT NULL,
   follow_up_violence_justification CHARACTER VARYING(1) NOT NULL,
   follow_up_victim_vulnerability CHARACTER VARYING(1) NOT NULL,
   follow_up_risk_level CHARACTER VARYING(1) NOT NULL,
   follow_up_owner_general_user CHARACTER VARYING(36) NOT NULL,
   follow_up_status CHARACTER VARYING(1) NOT NULL,
   CONSTRAINT follow_up_pkey PRIMARY KEY (follow_up_id),
   CONSTRAINT follow_up_follow_up_i_code_key UNIQUE (follow_up_i_code)
);


ALTER TABLE follow_up_entry_acting ADD CONSTRAINT follow_up_entry_acting_rel_barrier_follow_up_entry
   FOREIGN KEY ( follow_up_entry_acting_rel_barrier_follow_up_entry) REFERENCES rel_barrier_follow_up_entry(rel_barrier_follow_up_entry_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE rel_barrier_follow_up_entry ADD CONSTRAINT rel_barrier_follow_up_entry_barrier
   FOREIGN KEY ( rel_barrier_follow_up_entry_barrier) REFERENCES barrier(barrier_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_barrier_follow_up_entry ADD CONSTRAINT rel_barrier_follow_up_entry_follow_up_entry
   FOREIGN KEY ( rel_barrier_follow_up_entry_follow_up_entry) REFERENCES follow_up_entry(follow_up_entry_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE follow_up_entry ADD CONSTRAINT follow_up_entry_follow_up
   FOREIGN KEY ( follow_up_id) REFERENCES follow_up(follow_up_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_alert_follow_up ADD CONSTRAINT rel_alert_follow_up_follow_up
   FOREIGN KEY ( rel_alert_follow_up_follow_up) REFERENCES follow_up_entry(follow_up_entry_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;ALTER TABLE rel_alert_follow_up ADD CONSTRAINT rel_alert_follow_up_alert
   FOREIGN KEY ( rel_alert_follow_up_alert) REFERENCES alert(alert_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE barrier ADD CONSTRAINT barrier_sector_barrier
   FOREIGN KEY ( sector_barrier_id) REFERENCES sector_barrier(sector_barrier_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;



-- Datos


INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (1, '7274e166-955f-11ee-9005-3ca06777a813', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Entidad', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (2, '7274e166-955f-11ee-9005-3ca06777a812', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Justicia INMLCF', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (3, '7274e166-955f-11ee-9005-3ca06777a811', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Justicia juzgado', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (4, '7274e166-955f-11ee-9005-3ca06777a810', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Justicia Fiscalía', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (5, '7274e166-955f-11ee-9005-3ca06777a809', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Justicia General', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (6, '7274e166-955f-11ee-9005-3ca06777a808', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'UNP', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (7, '7274e166-955f-11ee-9005-3ca06777a807', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Protección Inspección', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (8, '7274e166-955f-11ee-9005-3ca06777a806', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Protección Juzgado Municipal', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (9, '7274e166-955f-11ee-9005-3ca06777a805', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Protección Policía', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (10, '7274e166-955f-11ee-9005-3ca06777a804', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Protección Juzgado Garantías', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (11, '7274e166-955f-11ee-9005-3ca06777a803', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Protección Fiscalía', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (12, '7274e166-955f-11ee-9005-3ca06777a802', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Protección Comisaría', NULL, NULL);
INSERT INTO salvia.sector_barrier (sector_barrier_id, sector_barrier_i_code, sector_barrier_creation_date, sector_barrier_update_date, sector_barrier_name, sector_barrier_description, sector_barrier_sector) VALUES (13, '7274e166-955f-11ee-9005-3ca06777a801', '2025-06-30 11:20:14.341421', '2025-06-30 11:20:14.341421', 'Salud', NULL, NULL);
SELECT pg_catalog.setval('salvia.sector_barrier_sector_barrier_id_seq', 13, true);


INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (1, '7274e166-955f-11ee-9005-3ca06777a507', 'Activación de la ruta de atención integral', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (2, '7274e166-955f-11ee-9005-3ca06777a506', 'Entrega de medicamentos', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (3, '7274e166-955f-11ee-9005-3ca06777a505', 'Continuidad en tratamiento', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (4, '7274e166-955f-11ee-9005-3ca06777a504', 'Programación de citas médicas', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (5, '7274e166-955f-11ee-9005-3ca06777a503', 'Atención en urgencias', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (6, '7274e166-955f-11ee-9005-3ca06777a502', 'Aseguramiento en salud, cobro de cuotas moderadoras', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (7, '7274e166-955f-11ee-9005-3ca06777a501', 'Actitud revictimizante de funcionarios(as)', NULL, 13, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (8, '7274e166-955f-11ee-9005-3ca06777a408', 'Actitud revictimizante de funcionarios(as)', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (9, '7274e166-955f-11ee-9005-3ca06777a407', 'Activación de la ruta de atención integral', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (10, '7274e166-955f-11ee-9005-3ca06777a406', 'No se realiza valoración del riesgo', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (11, '7274e166-955f-11ee-9005-3ca06777a405', 'Seguimiento a la medida de protección / atención', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (12, '7274e166-955f-11ee-9005-3ca06777a404', 'Otorgamiento de medida de atención', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (13, '7274e166-955f-11ee-9005-3ca06777a403', 'Otorgamiento de medida de protección definitiva', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (14, '7274e166-955f-11ee-9005-3ca06777a402', 'Otorgamiento de medida de protección provisional', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (15, '7274e166-955f-11ee-9005-3ca06777a401', 'Demora injustificada en la atención', NULL, 12, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (16, '7274e166-955f-11ee-9005-3ca06777a304', 'Actitud revictimizante de funcionarios(as)', NULL, 11, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (17, '7274e166-955f-11ee-9005-3ca06777a303', 'Activación de la ruta de atención integral', NULL, 11, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (18, '7274e166-955f-11ee-9005-3ca06777a302', 'Solicitud de medida de protección', NULL, 11, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (19, '7274e166-955f-11ee-9005-3ca06777a301', 'Demora injustificada en la atención', NULL, 11, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (20, '7274e166-955f-11ee-9005-3ca06777a202', 'Actitud revictimizante de funcionarios(as)', NULL, 10, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (21, '7274e166-955f-11ee-9005-3ca06777a201', 'Otorgamiento medida de protección provisional', NULL, 10, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (22, '7274e166-955f-11ee-9005-3ca06777a105', 'Actitud revictimizante de funcionarios(as)', NULL, 9, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (23, '7274e166-955f-11ee-9005-3ca06777a104', 'Activación de la ruta de atención integral', NULL, 9, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (24, '7274e166-955f-11ee-9005-3ca06777a103', 'Acudir como primer respondiente frente a actos urgentes, incluyendo capturas en flagrancia', NULL, 9, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (25, '7274e166-955f-11ee-9005-3ca06777a102', 'Incumplimiento a órdenes impartidas por autoridad administrativa', NULL, 9, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (26, '7274e166-955f-11ee-9005-3ca06777a101', 'Demora injustificada en la atención', NULL, 9, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (27, '7274e166-955f-11ee-9005-3ca06777a007', 'Actitud revictimizante de funcionarios(as)', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (28, '7274e166-955f-11ee-9005-3ca06777a006', 'Activación de la ruta de atención integral', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (29, '7274e166-955f-11ee-9005-3ca06777a005', 'Seguimiento a la medida de protección / atención', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (30, '7274e166-955f-11ee-9005-3ca06777a004', 'Otorgamiento medida de atención', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (31, '7274e166-955f-11ee-9005-3ca06777a003', 'Otorgamiento medida de protección definitiva', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (32, '7274e166-955f-11ee-9005-3ca06777a002', 'Otorgamiento medida de protección provisional', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (33, '7274e166-955f-11ee-9005-3ca06777a001', 'Demora injustificada en la atención', NULL, 8, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (34, '7274e166-955f-11ee-9005-3ca06777b007', 'Actitud revictimizante de funcionarios(as)', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (35, '7274e166-955f-11ee-9005-3ca06777b006', 'Activación de la ruta de atención integral', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (36, '7274e166-955f-11ee-9005-3ca06777b005', 'Seguimiento a la medida de protección / atención', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (37, '7274e166-955f-11ee-9005-3ca06777b004', 'Otorgamiento medida de atención ', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (38, '7274e166-955f-11ee-9005-3ca06777b003', 'Otorgamiento medida de protección definitiva ', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (39, '7274e166-955f-11ee-9005-3ca06777b002', 'Otorgamiento medida de protección provisional ', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (40, '7274e166-955f-11ee-9005-3ca06777b001', 'Demora injustificada en la atención', NULL, 7, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (41, '7274e166-955f-11ee-9005-3ca06777b104', 'Actitud revictimizante de funcionarios(as)', NULL, 6, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (42, '7274e166-955f-11ee-9005-3ca06777b103', 'Activación de la ruta de atención integral', NULL, 6, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (43, '7274e166-955f-11ee-9005-3ca06777b102', 'Valoración del riesgo asociado a la labor', NULL, 6, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (44, '7274e166-955f-11ee-9005-3ca06777b101', 'Demora injustificada en la atención', NULL, 6, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (45, '7274e166-955f-11ee-9005-3ca06777b206', 'Actitud revictimizante de funcionarios(as)', NULL, 5, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (46, '7274e166-955f-11ee-9005-3ca06777b205', 'Activación de la ruta de atención integral', NULL, 5, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (47, '7274e166-955f-11ee-9005-3ca06777b204', 'Valoración del riesgo', NULL, 5, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (48, '7274e166-955f-11ee-9005-3ca06777b203', 'Remisión a Instituto Nacional de Medicina Legal y Ciencias Forenses para valoración médico- legal', NULL, 5, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (49, '7274e166-955f-11ee-9005-3ca06777b202', 'Recepción de denuncia', NULL, 5, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (50, '7274e166-955f-11ee-9005-3ca06777b201', 'Demora injustificada en la atención', NULL, 5, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (51, '7274e166-955f-11ee-9005-3ca06777b305', 'Actitud revictimizante de funcionarios(as)', NULL, 4, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (52, '7274e166-955f-11ee-9005-3ca06777b304', 'Activación de la ruta de atención integral', NULL, 4, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (53, '7274e166-955f-11ee-9005-3ca06777b303', 'Impulso procesal (actos de investigación, ampliación de denuncia)', NULL, 4, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (54, '7274e166-955f-11ee-9005-3ca06777b302', 'Solicitudes injustificadas de aplazamiento', NULL, 4, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (55, '7274e166-955f-11ee-9005-3ca06777b301', 'Información acerca del estado del proceso o citaciones al mismo', NULL, 4, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (56, '7274e166-955f-11ee-9005-3ca06777b405', 'Actitud revictimizante de funcionarios(as)', NULL, 3, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (57, '7274e166-955f-11ee-9005-3ca06777b404', 'Activación de la ruta de atención integral', NULL, 3, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (58, '7274e166-955f-11ee-9005-3ca06777b403', 'Impulso procesal', NULL, 3, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (59, '7274e166-955f-11ee-9005-3ca06777b402', 'Aplazamientos injustificados', NULL, 3, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (60, '7274e166-955f-11ee-9005-3ca06777b401', 'Información acerca del estado del proceso o citaciones al mismo', NULL, 3, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (61, '7274e166-955f-11ee-9005-3ca06777b502', 'Actitud revictimizante de funcionarios(as)', NULL, 2, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (62, '7274e166-955f-11ee-9005-3ca06777b501', 'Realización del examen médico- legal', NULL, 2, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (63, '7274e166-955f-11ee-9005-3ca06777b604', 'Actitud revictimizante de funcionarios(as)', NULL, 1, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (64, '7274e166-955f-11ee-9005-3ca06777b603', 'Activación de la ruta de atención integral', NULL, 1, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (65, '7274e166-955f-11ee-9005-3ca06777b602', 'Prestación de las medidas de atención', NULL, 1, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
INSERT INTO salvia.barrier (barrier_id, barrier_i_code, barrier_name, barrier_description, sector_barrier_id, barrier_creation_date, barrier_update_date) VALUES (66, '7274e166-955f-11ee-9005-3ca06777b601', 'Demora injustificada en la atención', NULL, 1, '2025-07-02 16:33:35.540744', '2025-07-02 16:33:35.540744');
SELECT pg_catalog.setval('salvia.barrier_barrier_id_seq', 66, true);




-- Scripts para la creación del formulario 2

CREATE TABLE victim_case_form2 (
   victim_case_form2_id SERIAL NOT NULL,
   victim_case_form2_i_code CHARACTER VARYING(36) NOT NULL,
   victim_case_form2_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_case_form2_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_case_form2_identity_name CHARACTER VARYING(32),
   victim_case_form2_victim_phone BIGINT ,
   victim_case_form2_facts_description TEXT NOT NULL,
   victim_case_form2_facts_date DATE NOT NULL,
   victim_case_form2_facts_start_time TIME WITHOUT TIME ZONE NOT NULL,
   victim_case_form2_facts_town_code CHARACTER VARYING(8) NOT NULL,
   victim_case_form2_facts_zone BIGINT NOT NULL,
   victim_case_form2_facts_address CHARACTER VARYING(32),
   victim_case_form2_scenario_violence BIGINT NOT NULL,
   victim_case_form2_reported_previously BIGINT NOT NULL,
   victim_case_form2_recurrence_aggression BIGINT NOT NULL,
   victim_case_form2_num_agressors BIGINT NOT NULL,
   victim_case_form2_proximity_principal_aggressor BIGINT NOT NULL,
   victim_case_form2_relationship_with_presumed_aggressor BIGINT,
   victim_case_form2_economically_dependent BIGINT,
   victim_case_form2_aggressor_gender_identity BIGINT NOT NULL,
   victim_case_form2_aggressor_names CHARACTER VARYING(64) ,
   victim_case_form2_aggressor_doc_type BIGINT ,
   victim_case_form2_aggressor_doc_number CHARACTER VARYING(32) ,
   victim_case_form2_aggressor_address CHARACTER VARYING(32) ,
   victim_case_form2_aggressor_phone BIGINT ,
   victim_case_form2_aggressor_violence_physical_increase BIGINT NOT NULL,
   victim_case_form2_aggressor_weapon_used BIGINT NOT NULL,
   victim_case_form2_aggressor_threat_kill BIGINT NOT NULL,
   victim_case_form2_aggressor_pursues_spies_destroys BIGINT NOT NULL,
   victim_case_form2_aggressor_capable_of_killing BIGINT NOT NULL,
   victim_case_form2_aggressor_has_access_to_weapons BIGINT NOT NULL,
   victim_case_form2_partner_unemployed BIGINT,
   victim_case_form2_partner_other_denunciations BIGINT,
   victim_case_form2_aggressor_has_penal_background BIGINT,
   victim_case_form2_aggressor_forced_sex BIGINT,
   victim_case_form2_aggressor_attempted_strangulation BIGINT,
   victim_case_form2_aggressor_consumes_drugs BIGINT,
   victim_case_form2_aggressor_is_alcoholic BIGINT,
   victim_case_form2_partner_controls BIGINT,
   victim_case_form2_aggressor_had_hit_in_vulnerability BIGINT,
   victim_case_form2_partner_threatened_suicide BIGINT,
   victim_case_form2_partner_threatened_damage_members BIGINT,
   victim_case_form2_thoughts_of_self_harm BIGINT,
   victim_case_form2_aggressor_limits_contact_support_networks BIGINT,
   victim_case_form2_still_lives_with_aggressor BIGINT,
   victim_case_form2_aggressor_violently_jealous BIGINT,
   victim_case_form2_aggressor_unemployed BIGINT,
   victim_case_form2_aggressor_has_penal_background_2 BIGINT,
   victim_case_form2_aggressor_sexually_harassment BIGINT,
   victim_case_form2_aggressor_use_drugs BIGINT,
   victim_case_form2_aggressor_is_alcoholic_2 BIGINT,
   victim_case_form2_aggressor_controls BIGINT,
   victim_case_form2_aggressor_threatened_damage_members BIGINT,
   victim_case_form2_thoughts_of_self_harm_2 BIGINT,
   victim_case_form2_aggressor_common_spaces BIGINT,
   victim_case_form2_aggressor_hierarchy BIGINT,
   victim_case_form2_birth_date DATE NOT NULL,
   victim_case_form2_physical_mental_sensory_difficulties BIGINT NOT NULL,
   victim_case_form2_nationality BIGINT NOT NULL,
   victim_case_form2_specified_nationality BIGINT,
   victim_case_form2_migration_condition BIGINT,
   victim_case_form2_gender_identity BIGINT NOT NULL,
   victim_case_form2_sexual_orientation BIGINT NOT NULL,
   victim_case_form2_assigned_sex_at_birth BIGINT NOT NULL,
   victim_case_form2_ethnic_affiliation BIGINT NOT NULL,
   victim_case_form2_indigenous_people BIGINT,
   victim_case_form2_campesino_recognition BIGINT NOT NULL,
   victim_case_form2_marital_status BIGINT NOT NULL,
   victim_case_form2_last_education_level BIGINT NOT NULL,
   victim_case_form2_occupation BIGINT NOT NULL,
   victim_case_form2_income_generation_method BIGINT NOT NULL,
   victim_case_form2_employment_relationship BIGINT,
   victim_case_form2_approx_start_asp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_case_form2_housing_tenancy_form BIGINT NOT NULL,
   victim_case_form2_housing_stratum BIGINT NOT NULL,
   victim_case_form2_currently_pregnant BIGINT NOT NULL,
   victim_case_form2_residence_town CHARACTER VARYING(8) NOT NULL,
   victim_case_form2_residence_address CHARACTER VARYING(32) NOT NULL,
   victim_case_form2_residence_zone BIGINT NOT NULL,
   victim_case_form2_support_contact_names CHARACTER VARYING(64),
   victim_case_form2_support_contact_phone BIGINT,
   victim_case_form2_support_contact_email CHARACTER VARYING(128),
   victim_case_form2_support_contact_kinship BIGINT,
   victim_case_form2_saliva_management_explanation TEXT NOT NULL,
   victim_case_form2_activities_unable_to_hear INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_talk INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_see INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_move INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_take INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_understand INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_eat INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_interact INTEGER  NOT NULL,
   victim_case_form2_activities_unable_to_do_everyday INTEGER  NOT NULL,
   victim_case_form2_person_with_disability BIGINT NOT NULL,
   victim_case_form2_require_language_interpreter BIGINT NOT NULL,
   victim_case_form2_workplace_sector_occurrence BIGINT,
   victim_case_form2_violence_motivated_by_gender BIGINT NOT NULL,
   victim_case_form2_attention_was_appropriate BIGINT,
   victim_case_form2_aggressors_occupation BIGINT NOT NULL,
   victim_case_form2_language_assistance CHARACTER VARYING(32) ,
   victim_case_form2_stopped_seeking_help BIGINT ,
   victim_case_form2_victim_health_to_blackmail BIGINT ,
   victim_case_form2_threatened_reveal_sexual_orientation BIGINT ,
   victim_case_form2_aggressor_taken_advantage_physical_vulnerabil BIGINT ,
   victim_case_form2_aggressor_sexually_harassment_2 BIGINT ,
   victim_case_form2_aggressor_used_position_authority BIGINT ,
   victim_case_form2_allows_easy_report BIGINT ,
   victim_case_form2_violence_motivated_by_gender_2 BIGINT ,
   victim_case_form2_risk_score INTEGER ,
   victim_case_form2_risk_level INTEGER ,
   victim_case_form2_victim_case BIGINT NOT NULL,
   CONSTRAINT victim_case_form2_pkey PRIMARY KEY (victim_case_form2_id),
   CONSTRAINT victim_case_form2_victim_contact_form2_i_code_key UNIQUE (victim_case_form2_i_code)
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE victim_case_form2 TO salvia_admin;


CREATE TABLE victim_case_form2_enums (
   victim_case_form2_enums_id SERIAL NOT NULL,
   victim_case_form2_enums_i_code CHARACTER VARYING(36) NOT NULL,
   victim_case_form2_enums_name CHARACTER VARYING(96) NOT NULL,
   victim_case_form2_enums_code CHARACTER VARYING(2) NOT NULL,
   victim_case_form2_enums_category CHARACTER VARYING(64) NOT NULL,
   CONSTRAINT victim_case_form2_enums_pkey PRIMARY KEY (victim_case_form2_enums_id),
   CONSTRAINT victim_case_form2_enums_victim_case_form2_enums_i_code_key UNIQUE (victim_case_form2_enums_i_code),
   CONSTRAINT victim_case_form2_enums_victim_case_form2_enums_name_key UNIQUE (victim_case_form2_enums_name)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE victim_case_form2_enums TO salvia_admin;



CREATE TABLE rel_victim_case_form2_enums_victim_case_form2 (
   rel_victim_case_form2_enums_victim_case_form2_id SERIAL NOT NULL,
   victim_case_form2_enums_id BIGINT NOT NULL,
   victim_case_form2_id BIGINT NOT NULL,
   CONSTRAINT rel_victim_case_form2_enums_victim_case_form2_pkey PRIMARY KEY (rel_victim_case_form2_enums_victim_case_form2_id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE rel_victim_case_form2_enums_victim_case_form2 TO salvia_admin;


CREATE TABLE rel_victim_case_form2_enums_victim_contact_form2 (
   rel_victim_case_form2_enums_victim_contact_form2_id SERIAL NOT NULL,
   victim_case_form2_enums_id BIGINT NOT NULL,
   victim_contact_form2_id BIGINT NOT NULL,
   CONSTRAINT rel_victim_case_form2_enums_victim_contact_form2_pkey PRIMARY KEY (rel_victim_case_form2_enums_victim_contact_form2_id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE rel_victim_case_form2_enums_victim_contact_form2 TO salvia_admin;



CREATE TABLE victim_contact_form2 (
   victim_contact_form2_id SERIAL NOT NULL,
   victim_case_form2_i_code CHARACTER VARYING(36) NOT NULL,
   victim_case_form2_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_case_form2_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
   victim_contact_form2_will_receive_call BIGINT,
   victim_contact_form2_has_care_role BIGINT,
   victim_contact_form2_victim_aware_of_report BIGINT,
   victim_contact_form2_reporter_names CHARACTER VARYING(64) ,
   victim_contact_form2_reporter_phone BIGINT ,
   victim_contact_form2_victim_col_phone BIGINT NOT NULL,
   victim_contact_form2_facts_description TEXT NOT NULL,
   victim_contact_form2_best_contact_time TIME WITHOUT TIME ZONE NOT NULL,
   victim_contact_form2_report_type BIGINT NOT NULL,
   victim_contact_form2_report_type_details BIGINT,
   victim_contact_form2_victim_contact BIGINT NOT NULL,
   CONSTRAINT victim_contact_form2_pkey PRIMARY KEY (victim_contact_form2_id),
   CONSTRAINT victim_contact_form2_victim_case_form2_i_code_key UNIQUE (victim_case_form2_i_code)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE victim_contact_form2 TO salvia_admin;


ALTER TABLE victim_contact_form2 ADD CONSTRAINT victim_contact_form2_will_receive_call
FOREIGN KEY ( victim_contact_form2_will_receive_call) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_contact_form2 ADD CONSTRAINT victim_contact_form2_has_care_role
FOREIGN KEY ( victim_contact_form2_has_care_role) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_contact_form2 ADD CONSTRAINT victim_contact_form2_victim_aware_of_report
FOREIGN KEY ( victim_contact_form2_victim_aware_of_report) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_contact_form2 ADD CONSTRAINT victim_contact_form2_victim_contact
FOREIGN KEY ( victim_contact_form2_victim_contact) REFERENCES victim_contact(victim_contact_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_facts_zone
   FOREIGN KEY ( victim_case_form2_facts_zone) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_scenario_violence
   FOREIGN KEY ( victim_case_form2_scenario_violence) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_reported_previously
   FOREIGN KEY ( victim_case_form2_reported_previously) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_recurrence_aggression
   FOREIGN KEY ( victim_case_form2_recurrence_aggression) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_num_agressors
   FOREIGN KEY ( victim_case_form2_num_agressors) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_proximity_principal_aggressor
   FOREIGN KEY ( victim_case_form2_proximity_principal_aggressor) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_relationship_with_presumed_aggressor
   FOREIGN KEY ( victim_case_form2_relationship_with_presumed_aggressor) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_economically_dependent
   FOREIGN KEY ( victim_case_form2_economically_dependent) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_gender_identity
   FOREIGN KEY ( victim_case_form2_aggressor_gender_identity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_doc_type
   FOREIGN KEY ( victim_case_form2_aggressor_doc_type) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_violence_physical_increase
   FOREIGN KEY ( victim_case_form2_aggressor_violence_physical_increase) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_weapon_used
   FOREIGN KEY ( victim_case_form2_aggressor_weapon_used) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_threat_kill
   FOREIGN KEY ( victim_case_form2_aggressor_threat_kill) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_pursues_spies_destroys
   FOREIGN KEY ( victim_case_form2_aggressor_pursues_spies_destroys) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_capable_of_killing
   FOREIGN KEY ( victim_case_form2_aggressor_capable_of_killing) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_has_access_to_weapons
   FOREIGN KEY ( victim_case_form2_aggressor_has_access_to_weapons) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_partner_unemployed
   FOREIGN KEY ( victim_case_form2_partner_unemployed) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_partner_other_denunciations
   FOREIGN KEY ( victim_case_form2_partner_other_denunciations) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_has_penal_background
   FOREIGN KEY ( victim_case_form2_aggressor_has_penal_background) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_forced_sex
   FOREIGN KEY ( victim_case_form2_aggressor_forced_sex) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_attempted_strangulation
   FOREIGN KEY ( victim_case_form2_aggressor_attempted_strangulation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_consumes_drugs
   FOREIGN KEY ( victim_case_form2_aggressor_consumes_drugs) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_is_alcoholic
   FOREIGN KEY ( victim_case_form2_aggressor_is_alcoholic) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_partner_controls
   FOREIGN KEY ( victim_case_form2_partner_controls) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_had_hit_in_vulnerability
   FOREIGN KEY ( victim_case_form2_aggressor_had_hit_in_vulnerability) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_partner_threatened_suicide
   FOREIGN KEY ( victim_case_form2_partner_threatened_suicide) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_partner_threatened_damage_members
   FOREIGN KEY ( victim_case_form2_partner_threatened_damage_members) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_thoughts_of_self_harm
   FOREIGN KEY ( victim_case_form2_thoughts_of_self_harm) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_limits_contact_support_networks
   FOREIGN KEY ( victim_case_form2_aggressor_limits_contact_support_networks) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_still_lives_with_aggressor
   FOREIGN KEY ( victim_case_form2_still_lives_with_aggressor) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_violently_jealous
   FOREIGN KEY ( victim_case_form2_aggressor_violently_jealous) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_unemployed
   FOREIGN KEY ( victim_case_form2_aggressor_unemployed) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_has_penal_background_2
   FOREIGN KEY ( victim_case_form2_aggressor_has_penal_background_2) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_sexually_harassment
   FOREIGN KEY ( victim_case_form2_aggressor_sexually_harassment) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_use_drugs
   FOREIGN KEY ( victim_case_form2_aggressor_use_drugs) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_is_alcoholic_2
   FOREIGN KEY ( victim_case_form2_aggressor_is_alcoholic_2) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_controls
   FOREIGN KEY ( victim_case_form2_aggressor_controls) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_threatened_damage_members
   FOREIGN KEY ( victim_case_form2_aggressor_threatened_damage_members) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_thoughts_of_self_harm_2
   FOREIGN KEY ( victim_case_form2_thoughts_of_self_harm_2) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_common_spaces
   FOREIGN KEY ( victim_case_form2_aggressor_common_spaces) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_hierarchy
   FOREIGN KEY ( victim_case_form2_aggressor_hierarchy) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_physical_mental_sensory_difficulties
   FOREIGN KEY ( victim_case_form2_physical_mental_sensory_difficulties) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_nationality
   FOREIGN KEY ( victim_case_form2_nationality) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_specified_nationality
   FOREIGN KEY ( victim_case_form2_specified_nationality) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_migration_condition
   FOREIGN KEY ( victim_case_form2_migration_condition) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_gender_identity
   FOREIGN KEY ( victim_case_form2_gender_identity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_sexual_orientation
   FOREIGN KEY ( victim_case_form2_sexual_orientation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_assigned_sex_at_birth
   FOREIGN KEY ( victim_case_form2_assigned_sex_at_birth) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_ethnic_affiliation
   FOREIGN KEY ( victim_case_form2_ethnic_affiliation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_indigenous_people
   FOREIGN KEY ( victim_case_form2_indigenous_people) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_campesino_recognition
   FOREIGN KEY ( victim_case_form2_campesino_recognition) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_marital_status
   FOREIGN KEY ( victim_case_form2_marital_status) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_last_education_level
   FOREIGN KEY ( victim_case_form2_last_education_level) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_occupation
   FOREIGN KEY ( victim_case_form2_occupation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_income_generation_method
   FOREIGN KEY ( victim_case_form2_income_generation_method) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_housing_tenancy_form
   FOREIGN KEY ( victim_case_form2_housing_tenancy_form) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_housing_stratum
   FOREIGN KEY ( victim_case_form2_housing_stratum) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_currently_pregnant
   FOREIGN KEY ( victim_case_form2_currently_pregnant) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_residence_zone
   FOREIGN KEY ( victim_case_form2_residence_zone) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_support_contact_kinship
   FOREIGN KEY ( victim_case_form2_support_contact_kinship) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_victim_case
   FOREIGN KEY ( victim_case_form2_victim_case) REFERENCES victim_case(victim_case_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_employment_relationship
    FOREIGN KEY ( victim_case_form2_employment_relationship) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_stopped_seeking_help
   FOREIGN KEY ( victim_case_form2_stopped_seeking_help) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_victim_health_to_blackmail
   FOREIGN KEY ( victim_case_form2_victim_health_to_blackmail) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_threatened_reveal_sexual_orientation
    FOREIGN KEY ( victim_case_form2_threatened_reveal_sexual_orientation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_taken_advantage_physical_vulnerabil
   FOREIGN KEY ( victim_case_form2_aggressor_taken_advantage_physical_vulnerabil) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_sexually_harassment_2
   FOREIGN KEY ( victim_case_form2_aggressor_sexually_harassment_2) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_aggressor_used_position_authority
   FOREIGN KEY ( victim_case_form2_aggressor_used_position_authority) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE victim_case_form2 ADD CONSTRAINT victim_case_form2_violence_motivated_by_gender_2
   FOREIGN KEY ( victim_case_form2_violence_motivated_by_gender_2) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;



ALTER TABLE rel_victim_case_form2_enums_victim_case_form2 ADD CONSTRAINT rel_victim_case_form2_id
   FOREIGN KEY ( victim_case_form2_id) REFERENCES victim_case_form2(victim_case_form2_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE rel_victim_case_form2_enums_victim_case_form2 ADD CONSTRAINT rel_victim_case_form2_enums_id
   FOREIGN KEY ( victim_case_form2_enums_id) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE victim_contact_form2 ADD CONSTRAINT victim_contact_form2_report_type_details
       FOREIGN KEY ( victim_contact_form2_report_type_details) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE rel_victim_case_form2_enums_victim_contact_form2 ADD CONSTRAINT rel_victim_case_form2_id
   FOREIGN KEY ( victim_contact_form2_id) REFERENCES victim_contact_form2(victim_contact_form2_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE rel_victim_case_form2_enums_victim_contact_form2 ADD CONSTRAINT rel_victim_case_form2_enums_id
   FOREIGN KEY ( victim_case_form2_enums_id) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;





-- Inicialización de victim_case_form2_enums

INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'yes_no_y', 'y', 'yes_no');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'yes_no_n', 'n', 'yes_no');

INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_victim', 'v', 'victim_case_form2_report_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_reporter', 'r', 'victim_case_form2_report_type');

INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_details_ae', 'ae', 'victim_case_form2_report_type_details');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_details_et', 'et', 'victim_case_form2_report_type_details');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_details_ex', 'ex', 'victim_case_form2_report_type_details');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_details_ot', 'ot', 'victim_case_form2_report_type_details');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_report_type_details_an', 'an', 'victim_case_form2_report_type_details');

--victim_case_form2_victim_doc_type
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_cd', 'cd', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_cc', 'cc', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_ce', 'ce', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_no', 'no', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_nc', 'nc', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_mv', 'mv', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_dn', 'dn', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_oe', 'oe', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_pn', 'pn', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_ps', 'ps', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_pp', 'pp', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_pf', 'pf', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_pt', 'pt', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_si', 'si', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_rc', 'rc', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_sc', 'sc', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_ti', 'ti', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_vs', 'vs', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_vr', 'vr', 'victim_case_form2_victim_doc_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_victim_doc_type_np', 'np', 'victim_case_form2_victim_doc_type');
--victim_case_form2_facts_zone
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_facts_zone_cm', 'cm', 'victim_case_form2_facts_zone');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_facts_zone_pr', 'pr', 'victim_case_form2_facts_zone');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)	VALUES (public.uuid_generate_v4(), 'victim_case_form2_facts_zone_cp', 'cp', 'victim_case_form2_facts_zone');
-- victim_case_form2_scenario_violence
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_am',  'am', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ar',  'ar', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ca',  'ca', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_cr',  'cr', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ce',  'ce', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_co',  'co', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ed',  'ed', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ac',  'ac', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_te',  'te', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_cm',  'cm', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_in',  'in', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ex',  'ex', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ad',  'ad', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_fi',  'fi', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_es',  'es', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_gu',  'gu', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_cu',  'cu', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ci',  'ci', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_al',  'al', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_mi',  'mi', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ho',  'ho', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_of',  'of', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_pa',  'pa', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_si',  'si', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ta',  'ta', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_tr',  'tr', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_sp',  'sp', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_vp',  'vp', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_vi',  'vi', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_za',  'za', 'victim_case_form2_scenario_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scenario_violence_ot',  'ot', 'victim_case_form2_scenario_violence');
-- victim_case_form2_recurrence_aggression
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'victim_case_form2_recurrence_aggression_pr','pr','victim_case_form2_recurrence_aggression');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'victim_case_form2_recurrence_aggression_se','se','victim_case_form2_recurrence_aggression');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'victim_case_form2_recurrence_aggression_re','re','victim_case_form2_recurrence_aggression');
--victim_case_form2_num_agressors
INSERT INTO salvia.victim_case_form2_enums(   victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_num_agressors_01', '01', 'victim_case_form2_num_agressors');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_num_agressors_02', '02', 'victim_case_form2_num_agressors');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_num_agressors_03', '03', 'victim_case_form2_num_agressors');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_num_agressors_md', 'md', 'victim_case_form2_num_agressors');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_num_agressors_nd', 'nd', 'victim_case_form2_num_agressors');
--victim_case_form2_proximity_principal_aggressor
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_proximity_principal_aggressor_pc', 'pc', 'victim_case_form2_proximity_principal_aggressor');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_proximity_principal_aggressor_pn', 'pn', 'victim_case_form2_proximity_principal_aggressor');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_proximity_principal_aggressor_pd', 'pd', 'victim_case_form2_proximity_principal_aggressor');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_proximity_principal_aggressor_nd', 'nd', 'victim_case_form2_proximity_principal_aggressor');
--victim_case_form2_relationship_with_presumed_aggressor_01
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_am',  'am', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_en',  'en', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_he',  'he', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ma',  'ma', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_of',  'of', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_pa',  'pa', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_pi',  'pi', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ex',  'ex', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_rp',  'rp', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_su',  'su', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_co',  'co', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_pr',  'pr', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ce',  'ce', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ps',  'ps', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ci',  'ci', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_af',  'af', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ve',  've', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_cl',  'cl', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_sr',  'sr', 'victim_case_form2_relationship_with_presumed_aggressor_01');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_01_ot',  'ot', 'victim_case_form2_relationship_with_presumed_aggressor_01');
-- victim_case_form2_relationship_with_presumed_aggressor_02
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_dc', 'dc','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_de', 'de','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_mg', 'mg','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_ms', 'ms','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_pc', 'pc','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_sp', 'sp','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_lr', 'lr','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_rc', 'rc','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_pn', 'pn','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_en', 'en','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_an', 'an','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_fa', 'fa','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_di', 'di','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_is', 'is','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_ls', 'ls','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_as', 'as','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_ds', 'ds','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_ot', 'ot','victim_case_form2_relationship_with_presumed_aggressor_02');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_relationship_with_presumed_aggressor_02_ni', 'ni','victim_case_form2_relationship_with_presumed_aggressor_02');
-- victim_case_form2_aggressor_gender_identity
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_aggressor_gender_identity_ho', 'ho', 'victim_case_form2_aggressor_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_aggressor_gender_identity_ht', 'ht', 'victim_case_form2_aggressor_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_aggressor_gender_identity_mu', 'mu', 'victim_case_form2_aggressor_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_aggressor_gender_identity_mt', 'mt', 'victim_case_form2_aggressor_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_aggressor_gender_identity_ot', 'ot', 'victim_case_form2_aggressor_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_aggressor_gender_identity_ni', 'ni', 'victim_case_form2_aggressor_gender_identity');

--victim_case_form2_nationality
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_nationality_co', 'co', 'victim_case_form2_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_nationality_ex', 'ex', 'victim_case_form2_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_nationality_ap', 'ap', 'victim_case_form2_nationality');
-- victim_case_form2_specified_nationality
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_af','af', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ax','ax', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_al','al', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_dz','dz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_as','as', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ad','ad', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ao','ao', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ai','ai', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_aq','aq', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ag','ag', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ar','ar', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_am','am', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_aw','aw', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_au','au', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_at','at', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_az',  'az', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bs', 'bs', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bh', 'bh', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bd', 'bd', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bb', 'bb', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_by', 'by', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_be', 'be', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bz', 'bz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bj', 'bj', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bm', 'bm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bt', 'bt', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bo', 'bo', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bq', 'bq', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ba', 'ba', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bw', 'bw', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_br', 'br', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_io', 'io', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bn', 'bn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bg', 'bg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bf', 'bf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bi', 'bi', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kh', 'kh', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cm', 'cm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ca', 'ca', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cv', 'cv', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ky', 'ky', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cf', 'cf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_td', 'td', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cl', 'cl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cn', 'cn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cx', 'cx', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cc', 'cc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_co', 'co', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_km', 'km', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cg', 'cg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cd', 'cd', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ck', 'ck', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cr', 'cr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ci', 'ci', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_hr', 'hr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cu', 'cu', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cw', 'cw', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cy', 'cy', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_cz', 'cz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_dk', 'dk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_dj', 'dj', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_dm', 'dm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_do', 'do', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ec', 'ec', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_eg', 'eg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sv', 'sv', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gq', 'gq', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_er', 'er', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ee', 'ee', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sz', 'sz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_et', 'et', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_fk', 'fk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_fo', 'fo', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_fj', 'fj', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_fi', 'fi', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_fr', 'fr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gf', 'gf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pf', 'pf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tf', 'tf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ga', 'ga', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gm', 'gm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ge', 'ge', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_de', 'de', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gh', 'gh', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gi', 'gi', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gr', 'gr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gl', 'gl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gd', 'gd', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gp', 'gp', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gu', 'gu', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gt',  'gt',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gg',  'gg',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gn',  'gn',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gw',  'gw',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gy',  'gy',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ht',  'ht',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_hm',  'hm',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_va',  'va',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_hn',  'hn',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_hk',  'hk',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_hu',  'hu',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_is',  'is',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_in',  'in',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_id',  'id',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_xz',  'xz',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ir',  'ir',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_iq',  'iq',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ie',  'ie',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_im',  'im',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_il',  'il',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_it',  'it',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_jm',  'jm',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_jp',  'jp',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_je',  'je',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_jo',  'jo',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kz',  'kz',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ke',  'ke',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ki',  'ki',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kp',  'kp',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kr',  'kr',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kw',  'kw',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kg',  'kg',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_la',  'la',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lv',  'lv',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lb',  'lb',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ls',  'ls',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lr',  'lr',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ly',  'ly',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_li',  'li',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lt',  'lt',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lu',  'lu',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mo',  'mo',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mg',  'mg',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mw',  'mw',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_my',  'my',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mv',  'mv',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ml',  'ml',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mt',  'mt',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mh',  'mh',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mq',  'mq',  'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mr', 'mr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mu', 'mu', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_yt', 'yt', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mx', 'mx', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_fm', 'fm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_md', 'md', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mc', 'mc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mn', 'mn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_me', 'me', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ms', 'ms', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ma', 'ma', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mz', 'mz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mm', 'mm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_na', 'na', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_nr', 'nr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_np', 'np', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_nl', 'nl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_nc', 'nc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_nz', 'nz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ni', 'ni', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ne', 'ne', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ng', 'ng', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_nu', 'nu', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_nf', 'nf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mk', 'mk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mp', 'mp', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_no', 'no', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_om', 'om', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pk', 'pk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pw', 'pw', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ps', 'ps', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pa', 'pa', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pg', 'pg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_py', 'py', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pe', 'pe', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ph', 'ph', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pn', 'pn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pl', 'pl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pt', 'pt', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pr', 'pr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_qa', 'qa', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_re', 're', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ro', 'ro', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ru', 'ru', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_rw', 'rw', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_bl', 'bl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sh', 'sh', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_kn', 'kn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lc', 'lc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_mf', 'mf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_pm', 'pm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_vc', 'vc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ws', 'ws', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sm', 'sm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_st', 'st', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sa', 'sa', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sn', 'sn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_rs', 'rs', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sc', 'sc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sl', 'sl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sg', 'sg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sx', 'sx', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sk', 'sk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_si', 'si', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sb', 'sb', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_so', 'so', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_za', 'za', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gs', 'gs', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ss', 'ss', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_es', 'es', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_lk', 'lk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sd', 'sd', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sr', 'sr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sj', 'sj', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_se', 'se', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ch', 'ch', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_sy', 'sy', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tw', 'tw', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tj', 'tj', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tz', 'tz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_th', 'th', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tl', 'tl', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tg', 'tg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tk', 'tk', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_to', 'to', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tt', 'tt', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tn', 'tn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tr', 'tr', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tm', 'tm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tc', 'tc', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_tv', 'tv', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ug', 'ug', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ua', 'ua', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ae', 'ae', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_gb', 'gb', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_us', 'us', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_um', 'um', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_uy', 'uy', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_uz', 'uz', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_vu', 'vu', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ve', 've', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_vn', 'vn', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_vg', 'vg', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_vi', 'vi', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_wf', 'wf', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_eh', 'eh', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_ye', 'ye', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_zm', 'zm', 'victim_case_form2_specified_nationality');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specified_nationality_zw', 'zw', 'victim_case_form2_specified_nationality');
-- victim_case_form2_migration_condition
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_vi', 'vi', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_ex', 'ex', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_pm', 'pm', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_pe', 'pe', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_sr', 'sr', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_st', 'st', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_ic', 'ic', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_dv', 'dv', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_pa', 'pa', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_rr', 'rr', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_ni', 'ni', 'victim_case_form2_migration_condition');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_migration_condition_cr', 'cr', 'victim_case_form2_migration_condition');
--victim_case_form2_gender_identity
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_gender_identity_mu',  'mu',  'victim_case_form2_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_gender_identity_mt',  'mt',  'victim_case_form2_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_gender_identity_ho',  'ho',  'victim_case_form2_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_gender_identity_ht',  'ht',  'victim_case_form2_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_gender_identity_nb',  'nb',  'victim_case_form2_gender_identity');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_gender_identity_nf',  'nf',  'victim_case_form2_gender_identity');
-- victim_case_form2_sexual_orientation
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_he', 'he', 'victim_case_form2_sexual_orientation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_ga', 'ga', 'victim_case_form2_sexual_orientation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_le', 'le', 'victim_case_form2_sexual_orientation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_bi', 'bi', 'victim_case_form2_sexual_orientation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_ot', 'ot', 'victim_case_form2_sexual_orientation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_ns', 'ns', 'victim_case_form2_sexual_orientation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_sexual_orientation_ni', 'ni', 'victim_case_form2_sexual_orientation');
-- victim_case_form2_assigned_sex_at_birth
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_assigned_sex_at_birth_mu', 'mu', 'victim_case_form2_assigned_sex_at_birth');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_assigned_sex_at_birth_ho', 'ho', 'victim_case_form2_assigned_sex_at_birth');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_assigned_sex_at_birth_in', 'in', 'victim_case_form2_assigned_sex_at_birth');
-- victim_case_form2_ethnic_affiliation
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_af', 'af', 'victim_case_form2_ethnic_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_ne', 'ne', 'victim_case_form2_ethnic_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_pa', 'pa', 'victim_case_form2_ethnic_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_ar', 'ar', 'victim_case_form2_ethnic_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_in', 'in', 'victim_case_form2_ethnic_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_pr', 'pr', 'victim_case_form2_ethnic_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_ethnic_affiliation_ns', 'ns', 'victim_case_form2_ethnic_affiliation');
--victim_case_form2_indigenous_people
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ac',  'ac',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_am',  'am',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_wi',  'wi',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ya',  'ya',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_yr',  'yr',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_an',  'an',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ar',  'ar',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ww',  'ww',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ba',  'ba',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_br',  'br',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_bi',  'bi',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_be',  'be',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_bo',  'bo',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ka',  'ka',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ca',  'ca',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_kr',  'kr',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ce',  'ce',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ch',  'ch',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_co',  'co',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ko',  'ko',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ke',  'ke',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_pi',  'pi',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_aw',  'aw',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ku',  'ku',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_cu',  'cu',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_gu',  'gu',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_cr',  'cr',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_bn',  'bn',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ga',  'ga',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_de',  'de',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ta',  'ta',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_em',  'em',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_eb',  'eb',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ec',  'ec',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ep',  'ep',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ed',  'ed',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_nu',  'nu',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mi',  'mi',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ab',  'ab',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_qi',  'qi',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_gn',  'gn',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category)VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_go',  'go',  'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ji', 'ji', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_cm', 'cm', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_in', 'in', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_km', 'km', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_kf', 'kf', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_kg', 'kg', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_le', 'le', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ma', 'ma', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_hi', 'hi', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mk', 'mk', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ju', 'ju', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_kk', 'kk', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_hu', 'hu', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_jh', 'jh', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_jd', 'jd', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_nk', 'nk', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mb', 'mb', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mt', 'mt', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_je', 'je', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mr', 'mr', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mu', 'mu', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_no', 'no', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ok', 'ok', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_na', 'na', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_po', 'po', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_pa', 'pa', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_pr', 'pr', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_pt', 'pt', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ps', 'ps', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_pu', 'pu', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_as', 'as', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_qu', 'qu', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_sa', 'sa', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_si', 'si', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mp', 'mp', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_so', 'so', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_tu', 'tu', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ti', 'ti', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_tn', 'tn', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_tg', 'tg', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_tr', 'tr', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_tt', 'tt', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_to', 'to', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_tk', 'tk', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ts', 'ts', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_uk', 'uk', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_uw', 'uw', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ty', 'ty', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_wo', 'wo', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_wy', 'wy', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ui', 'ui', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_mm', 'mm', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_yi', 'yi', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_yg', 'yg', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_yk', 'yk', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_yn', 'yn', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_un', 'un', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_kp', 'kp', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_yt', 'yt', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ze', 'ze', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ua', 'ua', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_kn', 'kn', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ot', 'ot', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_qc', 'qc', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_k1', 'k1', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ay', 'ay', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_it', 'it', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_qb', 'qb', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_cl', 'cl', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_pn', 'pn', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ie', 'ie', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ip', 'ip', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_iv', 'iv', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_im', 'im', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ib', 'ib', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ye', 'ye', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ii', 'ii', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_nb', 'nb', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_gm', 'gm', 'victim_case_form2_indigenous_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_indigenous_people_ni', 'ni', 'victim_case_form2_indigenous_people');
--victim_case_form2_marital_status
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_so', 'so', 'victim_case_form2_marital_status');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_ul', 'ul', 'victim_case_form2_marital_status');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_ca', 'ca', 'victim_case_form2_marital_status');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_se', 'se', 'victim_case_form2_marital_status');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_di', 'di', 'victim_case_form2_marital_status');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_vi', 'vi', 'victim_case_form2_marital_status');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_marital_status_ni', 'ni', 'victim_case_form2_marital_status');
--victim_case_form2_last_education_level
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_sf', 'sf', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_pf', 'pf', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_sb', 'sb', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_bm', 'bm', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_tt', 'tt', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_pr', 'pr', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_es', 'es', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_do', 'do', 'victim_case_form2_last_education_level');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_last_education_level_ni', 'ni', 'victim_case_form2_last_education_level');
--victim_case_form2_occupation
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_sg', 'sg', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ca', 'ca', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_bc', 'bc', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ac', 'ac', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_cu', 'cu', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_td', 'td', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_cv', 'cv', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_cc', 'cc', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_se', 'se', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_tm', 'tm', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_co', 'co', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_mp', 'mp', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ao', 'ao', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_sc', 'sc', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ed', 'ed', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ar', 'ar', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_tp', 'tp', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_es', 'es', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_as', 'as', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_dh', 'dh', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_fp', 'fp', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_pe', 'pe', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_vt', 'vt', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_em', 'em', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ds', 'ds', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_af', 'af', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ni', 'ni', 'victim_case_form2_occupation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_occupation_ot', 'ot', 'victim_case_form2_occupation');
--victim_case_form2_income_generation_method
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_em', 'em', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_es', 'es', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_de', 'de', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_pe', 'pe', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_pr', 'pr', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_ta', 'ta', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_si', 'si', 'victim_case_form2_income_generation_method');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_income_generation_method_no', 'no', 'victim_case_form2_income_generation_method');
--victim_case_form2_income_generation_method
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_employment_relationship_pl', 'pl', 'victim_case_form2_employment_relationship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_employment_relationship_pr', 'pr', 'victim_case_form2_employment_relationship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_employment_relationship_co', 'co', 'victim_case_form2_employment_relationship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_employment_relationship_tr', 'tr', 'victim_case_form2_employment_relationship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_employment_relationship_ca', 'ca', 'victim_case_form2_employment_relationship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_employment_relationship_ps', 'ps', 'victim_case_form2_employment_relationship');

--victim_case_form2_housing_tenancy_form
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_pr', 'pr', 'victim_case_form2_housing_tenancy_form');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_ar', 'ar', 'victim_case_form2_housing_tenancy_form');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_fa', 'fa', 'victim_case_form2_housing_tenancy_form');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_eo', 'eo', 'victim_case_form2_housing_tenancy_form');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_al', 'al', 'victim_case_form2_housing_tenancy_form');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_sv', 'sv', 'victim_case_form2_housing_tenancy_form');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_tenancy_form_ni', 'ni', 'victim_case_form2_housing_tenancy_form');
--victim_case_form2_housing_stratum
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_se', 'se', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_01', '01', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_02', '02', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_03', '03', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_04', '04', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_05', '05', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_06', '06', 'victim_case_form2_housing_stratum');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_housing_stratum_na', 'na', 'victim_case_form2_housing_stratum');
--victim_case_form2_support_contact_kinship
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_ma', 'ma', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_he', 'he', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_hi', 'hi', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_pa', 'pa', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_ot', 'ot', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_am', 'am', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_je', 'je', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_cu', 'cu', 'victim_case_form2_support_contact_kinship');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_support_contact_kinship_no', 'no', 'victim_case_form2_support_contact_kinship');
--Mútiples -------------------------------
--victim_case_form2_type_violence_experienced
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_fi', 'fi', 'victim_case_form2_type_violence_experienced');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_ps', 'ps', 'victim_case_form2_type_violence_experienced');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_se', 'se', 'victim_case_form2_type_violence_experienced');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_po', 'po', 'victim_case_form2_type_violence_experienced');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_pl', 'pl', 'victim_case_form2_type_violence_experienced');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_re', 're', 'victim_case_form2_type_violence_experienced');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_type_violence_experienced_vi', 'vi', 'victim_case_form2_type_violence_experienced');
--victim_case_form2_subtype_violence_experienced_fi
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_ac', 'ac', 'victim_case_form2_subtype_violence_experienced_fi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_gp', 'gp', 'victim_case_form2_subtype_violence_experienced_fi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_qu', 'qu', 'victim_case_form2_subtype_violence_experienced_fi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_so', 'so', 'victim_case_form2_subtype_violence_experienced_fi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_ao', 'ao', 'victim_case_form2_subtype_violence_experienced_fi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_to', 'to', 'victim_case_form2_subtype_violence_experienced_fi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_fi_se', 'se', 'victim_case_form2_subtype_violence_experienced_fi');
--victim_case_form2_subtype_violence_experienced_ps
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_ps_sv', 'sv', 'victim_case_form2_subtype_violence_experienced_ps');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_ps_ia', 'ia', 'victim_case_form2_subtype_violence_experienced_ps');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_ps_av', 'av', 'victim_case_form2_subtype_violence_experienced_ps');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_ps_ce', 'ce', 'victim_case_form2_subtype_violence_experienced_ps');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_ps_ai', 'ai', 'victim_case_form2_subtype_violence_experienced_ps');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_ps_al', 'al', 'victim_case_form2_subtype_violence_experienced_ps');
--victim_case_form2_subtype_violence_experienced_se
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_se_ac', 'ac', 'victim_case_form2_subtype_violence_experienced_se');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_se_ab', 'ab', 'victim_case_form2_subtype_violence_experienced_se');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_se_tr', 'tr', 'victim_case_form2_subtype_violence_experienced_se');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_se_mu', 'mu', 'victim_case_form2_subtype_violence_experienced_se');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_se_ao', 'ao', 'victim_case_form2_subtype_violence_experienced_se');
--victim_case_form2_subtype_violence_experienced_po
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_po_ia', 'ia', 'victim_case_form2_subtype_violence_experienced_po');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_po_vm', 'vm', 'victim_case_form2_subtype_violence_experienced_po');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_po_nr', 'nr', 'victim_case_form2_subtype_violence_experienced_po');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_po_ri', 'ri', 'victim_case_form2_subtype_violence_experienced_po');
--victim_case_form2_subtype_violence_experienced_pl
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_pl_aa', 'aa', 'victim_case_form2_subtype_violence_experienced_pl');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_pl_cp', 'cp', 'victim_case_form2_subtype_violence_experienced_pl');
--victim_case_form2_subtype_violence_experienced_re
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_re_ab', 'ab', 'victim_case_form2_subtype_violence_experienced_re');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_re_es', 'es', 'victim_case_form2_subtype_violence_experienced_re');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_re_de', 'de', 'victim_case_form2_subtype_violence_experienced_re');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_re_ma', 'ma', 'victim_case_form2_subtype_violence_experienced_re');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_re_vi', 'vi', 'victim_case_form2_subtype_violence_experienced_re');
--victim_case_form2_subtype_violence_experienced_vi
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_vi_ah', 'ah', 'victim_case_form2_subtype_violence_experienced_vi');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_subtype_violence_experienced_vi_pr', 'pr', 'victim_case_form2_subtype_violence_experienced_vi');
--victim_case_form2_scope_of_violence
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ae', 'ae', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ac', 'ac', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_as', 'as', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ai', 'ai', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ic', 'ic', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_al', 'al', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ap', 'ap', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ri', 'ri', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_rn', 'rn', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ad', 'ad', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_co', 'co', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_de', 'de', 'victim_case_form2_scope_of_violence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_scope_of_violence_ce', 'ce', 'victim_case_form2_scope_of_violence');
--workplace_sector_occurrence
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_workplace_sector_occurrence_pu', 'pu', 'victim_case_form2_workplace_sector_occurrence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_workplace_sector_occurrence_pr', 'pr', 'victim_case_form2_workplace_sector_occurrence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_workplace_sector_occurrence_ei', 'ei', 'victim_case_form2_workplace_sector_occurrence');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_workplace_sector_occurrence_mt', 'mt', 'victim_case_form2_workplace_sector_occurrence');
--victim_case_form2_who_report_to
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_cf', 'cf', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_ip', 'ip', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_ci', 'ci', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_ep', 'ep', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_fg', 'fg', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_dp', 'dp', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_pg', 'pg', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_pm', 'pm', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_hc', 'hc', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_ic', 'ic', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_sm', 'sm', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_cj', 'cj', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_cv', 'cv', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_l1', 'l1', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_l2', 'l2', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_lp', 'lp', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_lo', 'lo', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_li', 'li', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_lc', 'lc', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_jc', 'jc', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_ie', 'ie', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_lt', 'lt', 'victim_case_form2_who_report_to');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_who_report_to_ot', 'ot', 'victim_case_form2_who_report_to');
--victim_case_form2_activities_unable_to_perform
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_oi', 'oi', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_ha', 'ha', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_ve', 've', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_mo', 'mo', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_ag', 'ag', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_en', 'en', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_co', 'co', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_re', 're', 'victim_case_form2_activities_unable_to_perform');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_activities_unable_to_perform_ac', 'ac', 'victim_case_form2_activities_unable_to_perform');
--victim_case_form2_adjustments_gbv
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_sc', 'sc', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_gi', 'gi', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_id', 'id', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_ap', 'ap', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_in', 'in', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_ad', 'ad', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_dt', 'dt', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_pr', 'pr', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_tr', 'tr', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_ni', 'ni', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_nr', 'nr', 'victim_case_form2_adjustments_gbv');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_adjustments_gbv_ns', 'ns', 'victim_case_form2_adjustments_gbv');
--victim_case_form2_law_1996
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_law_1996_va', 'va', 'victim_case_form2_law_1996');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_law_1996_ac', 'ac', 'victim_case_form2_law_1996');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_law_1996_di', 'di', 'victim_case_form2_law_1996');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_law_1996_na', 'na', 'victim_case_form2_law_1996');
--victim_case_form2_specially_protected_population
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_mb', 'mb', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_aj', 'aj', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_am', 'am', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_ls', 'ls', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_nd', 'nd', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_mc', 'mc', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_md', 'md', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_np', 'np', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_pp', 'pp', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_ph', 'ph', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_pl', 'pl', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_pd', 'pd', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_vc', 'vc', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_vf', 'vf', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_vd', 'vd', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_cf', 'cf', 'victim_case_form2_specially_protected_population');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_specially_protected_population_nn', 'nn', 'victim_case_form2_specially_protected_population');
--victim_case_form2_asp_mode
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_asp_mode_sp', 'sp', 'victim_case_form2_asp_mode');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_asp_mode_sv', 'sv', 'victim_case_form2_asp_mode');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_asp_mode_cp', 'cp', 'victim_case_form2_asp_mode');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_asp_mode_ar', 'ar', 'victim_case_form2_asp_mode');
--victim_case_form2_reason_asp
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_reason_asp_dp', 'dp', 'victim_case_form2_reason_asp');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_reason_asp_me', 'me', 'victim_case_form2_reason_asp');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_reason_asp_pi', 'pi', 'victim_case_form2_reason_asp');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_reason_asp_ma', 'ma', 'victim_case_form2_reason_asp');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_reason_asp_ni', 'ni', 'victim_case_form2_reason_asp');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_reason_asp_ot', 'ot', 'victim_case_form2_reason_asp');
--victim_case_form2_has_dependents
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_nn', 'nn', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_si', 'si', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_hi', 'hi', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_pa', 'pa', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_pr', 'pr', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_of', 'of', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_op', 'op', 'victim_case_form2_has_dependents');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_has_dependents_an', 'an', 'victim_case_form2_has_dependents');
--victim_case_form2_action_plan
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_ei', 'ei', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_ar', 'ar', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_ap', 'ap', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_pe', 'pe', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_me', 'me', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_ea', 'ea', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_se', 'se', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_rm', 'rm', 'victim_case_form2_action_plan');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'victim_case_form2_action_plan_ga', 'ga', 'victim_case_form2_action_plan');

-- MÓDULO FEMINICIDIO
-- feminicide_risk_form1_victim_sgsss_affiliation
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'feminicide_risk_form1_victim_sgsss_affiliation_su','su','feminicide_risk_form1_victim_sgsss_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'feminicide_risk_form1_victim_sgsss_affiliation_co','co','feminicide_risk_form1_victim_sgsss_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'feminicide_risk_form1_victim_sgsss_affiliation_re','re','feminicide_risk_form1_victim_sgsss_affiliation');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code,victim_case_form2_enums_name,victim_case_form2_enums_code,victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(),'feminicide_risk_form1_victim_sgsss_affiliation_na','na','feminicide_risk_form1_victim_sgsss_affiliation');
-- feminicide_risk_form1_victim_disability_type
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_df', 'df', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_ds', 'ds', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_da', 'da', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_dv', 'dv', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_so', 'so', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_di', 'di', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_dm', 'dm', 'feminicide_risk_form1_victim_disability_type');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_form1_victim_disability_type_mu', 'mu', 'feminicide_risk_form1_victim_disability_type');
-- feminicide_risk_financially_dependent_people
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_financially_dependent_people_gs', 'gs', 'feminicide_risk_financially_dependent_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_financially_dependent_people_lc', 'lc', 'feminicide_risk_financially_dependent_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_financially_dependent_people_m5', 'm5', 'feminicide_risk_financially_dependent_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_financially_dependent_people_pd', 'pd', 'feminicide_risk_financially_dependent_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_financially_dependent_people_pm', 'pm', 'feminicide_risk_financially_dependent_people');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_financially_dependent_people_ot', 'ot', 'feminicide_risk_financially_dependent_people');
-- feminicide_risk_victim_common_transport_mode
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_common_transport_mode_tr', 'tr', 'feminicide_risk_victim_common_transport_mode');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_common_transport_mode_fl', 'fl', 'feminicide_risk_victim_common_transport_mode');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_common_transport_mode_ae', 'ae', 'feminicide_risk_victim_common_transport_mode');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_common_transport_mode_mx', 'mx', 'feminicide_risk_victim_common_transport_mode');
-- feminicide_risk_places_visit_regularly
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_places_visit_regularly_ec', 'ec', 'feminicide_risk_places_visit_regularly');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_places_visit_regularly_cs', 'cs', 'feminicide_risk_places_visit_regularly');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_places_visit_regularly_fo', 'fo', 'feminicide_risk_places_visit_regularly');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_places_visit_regularly_cf', 'cf', 'feminicide_risk_places_visit_regularly');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_places_visit_regularly_lt', 'lt', 'feminicide_risk_places_visit_regularly');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_places_visit_regularly_ot', 'ot', 'feminicide_risk_places_visit_regularly');
-- feminicide_risk_victim_food_access_frequency
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_food_access_frequency_fm', 'fm', 'feminicide_risk_victim_food_access_frequency');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_food_access_frequency_cs', 'cs', 'feminicide_risk_victim_food_access_frequency');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_food_access_frequency_av', 'av', 'feminicide_risk_victim_food_access_frequency');
INSERT INTO salvia.victim_case_form2_enums(victim_case_form2_enums_i_code, victim_case_form2_enums_name, victim_case_form2_enums_code, victim_case_form2_enums_category) VALUES (public.uuid_generate_v4(), 'feminicide_risk_victim_food_access_frequency_nc', 'nc', 'feminicide_risk_victim_food_access_frequency');


--**************************************************************
-- Formularios de feminicidio

CREATE TABLE feminicide (
    feminicide_id SERIAL NOT NULL,
    feminicide_i_code CHARACTER VARYING(36) NOT NULL,
    feminicide_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_status CHARACTER VARYING(1) NOT NULL,
    feminicide_general_user CHARACTER VARYING(36) NOT NULL,
    feminicide_names CHARACTER VARYING(32) NOT NULL,
    feminicide_last_names CHARACTER VARYING(32) NOT NULL,
    feminicide_victim_doc_type CHARACTER VARYING(2) NOT NULL,
    feminicide_victim_doc_number CHARACTER VARYING(32) NOT NULL,
    CONSTRAINT feminicide_pkey PRIMARY KEY (feminicide_id),
    CONSTRAINT feminicide_feminicide_i_code_key UNIQUE (feminicide_i_code)
);

CREATE TABLE feminicide_risk (
    feminicide_risk_id SERIAL NOT NULL,
    feminicide_risk_i_code CHARACTER VARYING(36) NOT NULL,
    feminicide_risk_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_risk_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_risk_status CHARACTER VARYING(1) NOT NULL,
    feminicide_risk_general_user CHARACTER VARYING(36) NOT NULL,
    feminicide_risk_names CHARACTER VARYING(32) NOT NULL,
    feminicide_risk_last_names CHARACTER VARYING(32) NOT NULL,
    feminicide_risk_victim_doc_type CHARACTER VARYING(2) NOT NULL,
    feminicide_risk_victim_doc_number CHARACTER VARYING(32) NOT NULL,
    CONSTRAINT feminicide_risk_pkey PRIMARY KEY (feminicide_risk_id),
    CONSTRAINT feminicide_risk_feminicide_risk_i_code_key UNIQUE (feminicide_risk_i_code)
);


CREATE TABLE rel_case_owner_feminicide (
    rel_case_owner_feminicide_id SERIAL NOT NULL,
    rel_case_owner_feminicide_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    rel_case_owner_feminicide_status CHARACTER VARYING(1) NOT NULL,
    rel_case_owner_feminicide_case_owner BIGINT NOT NULL,
    rel_case_owner_feminicide_feminicide BIGINT NOT NULL,
    CONSTRAINT rel_case_owner_feminicide_pkey PRIMARY KEY (rel_case_owner_feminicide_id),
    CONSTRAINT rel_case_owner_feminicide_rel_case_owner_feminicide_i_code_key UNIQUE (rel_case_owner_feminicide_i_code)
);

CREATE TABLE rel_case_owner_feminicide_risk (
    rel_case_owner_feminicide_risk_id SERIAL NOT NULL,
    rel_case_owner_feminicide_risk_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    rel_case_owner_feminicide_risk_status CHARACTER VARYING(1) NOT NULL,
    rel_case_owner_feminicide_risk_case_owner BIGINT NOT NULL,
    rel_case_owner_feminicide_risk_feminicide_risk BIGINT NOT NULL,
    CONSTRAINT rel_case_owner_feminicide_risk_pkey PRIMARY KEY (rel_case_owner_feminicide_risk_id),
    CONSTRAINT rel_case_owner_feminicide_risk_rel_case_owner_feminicide_risk_i_code_key UNIQUE (rel_case_owner_feminicide_risk_i_code)
);

CREATE TABLE rel_victim_case_form2_enums_feminicide_form1 (
   rel_victim_case_form2_enums_feminicide_form1_id SERIAL NOT NULL,
   victim_case_form2_enums_id BIGINT NOT NULL,
   feminicide_form1_id BIGINT NOT NULL,
   CONSTRAINT rel_victim_case_form2_enums_feminicide_form1_pkey PRIMARY KEY (rel_victim_case_form2_enums_feminicide_form1_id)
);

CREATE TABLE rel_victim_case_form2_enums_feminicide_risk_form1 (
   rel_victim_case_form2_enums_feminicide_risk_form1_id SERIAL NOT NULL,
   victim_case_form2_enums_id BIGINT NOT NULL,
   feminicide_risk_form1_id BIGINT NOT NULL,
   CONSTRAINT rel_victim_case_form2_enums_feminicide_risk_form1_pkey PRIMARY KEY (rel_victim_case_form2_enums_feminicide_risk_form1_id)
);


CREATE TABLE feminicide_form1 (
    feminicide_form1_id SERIAL NOT NULL,
    feminicide_form1_i_code CHARACTER VARYING(36) NOT NULL,
    feminicide_form1_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_victim_identity_name CHARACTER VARYING(32) ,
    feminicide_form1_birth_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_death_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_victim_address TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_victim_zone BIGINT NOT NULL,
    feminicide_form1_victim_living_town TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_victim_marital_status BIGINT NOT NULL,
    feminicide_form1_victim_sex BIGINT NOT NULL,
    feminicide_form1_victim_gender_identity BIGINT NOT NULL,
    feminicide_form1_victim_sexual_orientation BIGINT NOT NULL,
    feminicide_form1_victim_ethnicity BIGINT NOT NULL,
    feminicide_form1_victim_indigenous_people BIGINT,
    feminicide_form1_victim_special_population BIGINT NOT NULL,
    feminicide_form1_victim_disability BIGINT NOT NULL,
    feminicide_form1_victim_disability_type BIGINT NOT NULL,
    feminicide_form1_presumed_aggressor_names CHARACTER VARYING(64) NOT NULL,
    feminicide_form1_presumed_aggressor_relation BIGINT NOT NULL,
    feminicide_form1_presumed_aggressor_known_vgb BIGINT NOT NULL,
    feminicide_form1_informant_names CHARACTER VARYING(64) NOT NULL,
    feminicide_form1_informant_identity_name CHARACTER VARYING(32) ,
    feminicide_form1_informant_doc_type CHARACTER VARYING(32) ,
    feminicide_form1_informant_doc_number CHARACTER VARYING(32) ,
    feminicide_form1_informant_birth_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_s_g_s_s_s_affiliation BIGINT NOT NULL,
    feminicide_form1_eps_name CHARACTER VARYING(64),
    feminicide_form1_informant_address TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_form1_informant_zone BIGINT NOT NULL,
    feminicide_form1_informant_living_town CHARACTER VARYING(8) NOT NULL,
    feminicide_form1_informant_phone CHARACTER VARYING(10) NOT NULL,
    feminicide_form1_emergency_contact_names CHARACTER VARYING(64) NOT NULL,
    feminicide_form1_emergency_contact_number CHARACTER VARYING(10) NOT NULL,
    feminicide_form1_informant_sex BIGINT NOT NULL,
    feminicide_form1_informant_gender_identity BIGINT NOT NULL,
    feminicide_form1_informant_sexual_orientation BIGINT NOT NULL,
    feminicide_form1_informant_ethnicity BIGINT NOT NULL,
    feminicide_form1_informant_indigenous_people BIGINT,
    feminicide_form1_informant_migratory_status BIGINT NOT NULL,
    feminicide_form1_informant_migration_situation BIGINT NOT NULL,
    feminicide_form1_informant_highest_education_level BIGINT NOT NULL,
    feminicide_form1_informant_special_population BIGINT NOT NULL,
    feminicide_form1_informant_current_employment BIGINT NOT NULL,
    feminicide_form1_informant_employment_access BIGINT NOT NULL,
    feminicide_form1_informant_employment_impact_description TEXT ,
    feminicide_form1_informant_primary_occupation BIGINT NOT NULL,
    feminicide_form1_informant_disability BIGINT NOT NULL,
    feminicide_form1_informant_disability_type BIGINT NOT NULL,
    feminicide_form1_situation_after_feminicide TEXT NOT NULL,
    feminicide_form1_additional_information TEXT ,
    feminicide_form1_any_assistance_received BIGINT NOT NULL,
    feminicide_form1_household_expense_responsibility BIGINT NOT NULL,
    feminicide_form1_post_death_economic_assumption BIGINT NOT NULL,
    feminicide_form1_economic_assumption_explanation TEXT ,
    feminicide_form1_any_dependent_people BIGINT NOT NULL,
    feminicide_form1_any_public_or_private_entity BIGINT NOT NULL,
    feminicide_form1_any_public_or_private_entity_explanation TEXT ,
    feminicide_form1_public_transport_access BIGINT NOT NULL,
    feminicide_form1_preferred_transportation_mode BIGINT NOT NULL,
    feminicide_form1_preferred_transportation_mode_explanation TEXT ,
    feminicide_form1_transportation_cost_estimate TEXT ,
    feminicide_form1_transport_difficulty BIGINT NOT NULL,
    feminicide_form1_transport_difficulty_explanation TEXT ,
    feminicide_form1_economic_resources_for_transport BIGINT NOT NULL,
    feminicide_form1_debt_or_help_due_to_transport BIGINT NOT NULL,
    feminicide_form1_debt_impact_explanation TEXT ,
    feminicide_form1_transport_subsidy_received BIGINT NOT NULL,
    feminicide_form1_transport_subsidy_explanation TEXT ,
    feminicide_form1_safety_transportation_concern BIGINT NOT NULL,
    feminicide_form1_safety_transportation_explanation TEXT ,
    feminicide_form1_food_access_frequency BIGINT NOT NULL,
    feminicide_form1_food_access_explanation TEXT ,
    feminicide_form1_aggressor_food_restriction BIGINT NOT NULL,
    feminicide_form1_aggressor_food_restriction_explanation TEXT ,
    feminicide_form1_food_income_support BIGINT NOT NULL,
    feminicide_form1_food_income_support_explanation TEXT ,
    feminicide_form1_juridical_assistance_received BIGINT NOT NULL,
    feminicide_form1_juridical_assistance_explanation TEXT ,
    feminicide_form1_victim_representation BIGINT NOT NULL,
    feminicide_form1_victim_representation_explanation TEXT ,
    feminicide_form1_psychosocial_support_received BIGINT NOT NULL,
    feminicide_form1_psychosocial_support_explanation TEXT ,
    feminicide_form1_emergency_emotional_crisis BIGINT NOT NULL,
    feminicide_form1_emergency_emotional_crisis_explanation TEXT ,
    feminicide_form1_aggressor_same_residence BIGINT NOT NULL,
    feminicide_form1_aggressor_same_residence_explanation TEXT ,
    feminicide_form1_aggressor_location_known BIGINT NOT NULL,
    feminicide_form1_aggressor_location_known_explanation TEXT ,
    feminicide_form1_any_type_of_assistance_received BIGINT NOT NULL,
    feminicide_form1_any_type_of_assistance_received_explanation TEXT ,
    feminicide_form1_compensation_fund_insurance_coverage BIGINT NOT NULL,
    feminicide_form1_compensation_fund_insurance_cov_explanation TEXT ,
    feminicide_form1_funeral_subsidy_received BIGINT NOT NULL,
    feminicide_form1_funeral_subsidy_explanation TEXT ,
    feminicide_form1_funeral_funds_available BIGINT NOT NULL,
    feminicide_form1_funeral_funds_available_explanation TEXT ,
    feminicide_form1_funeral_cost_value BIGINT NOT NULL,
    feminicide_form1_rename_reputation_impact BIGINT NOT NULL,
    feminicide_form1_rename_reputation_explanation TEXT ,
    feminicide_form1_additional_needs_description TEXT ,
    feminicide_form1_action_plan BIGINT NOT NULL,
    feminicide_form1_summary TEXT ,
    feminicide_form1_feminicide BIGINT NOT NULL,

    feminicide_form1_any_assistance_received_city_hall BIGINT ,
    feminicide_form1_any_assistance_received_womens_office BIGINT ,
    feminicide_form1_any_assistance_received_other_entity BIGINT ,
    feminicide_form1_any_assistance_received_other BIGINT ,
    feminicide_form1_assistance_received CHARACTER VARYING(120) ,
    feminicide_form1_assistance_received_city_hall CHARACTER VARYING(120) ,
    feminicide_form1_assistance_received_womens_office CHARACTER VARYING(120) ,
    feminicide_form1_assistance_received_other_entity CHARACTER VARYING(120) ,
    feminicide_form1_assistance_received_other CHARACTER VARYING(120) ,
    
    feminicide_form1_family_father INTEGER ,
    feminicide_form1_family_mother INTEGER ,
    feminicide_form1_family_stepfather INTEGER ,
    feminicide_form1_family_stepmother INTEGER ,
    feminicide_form1_family_partner INTEGER ,
    feminicide_form1_family_sibling_1 INTEGER ,
    feminicide_form1_family_sibling_2 INTEGER ,
    feminicide_form1_family_sibling_3 INTEGER ,
    feminicide_form1_family_sibling_4 INTEGER ,
    feminicide_form1_family_sibling_5 INTEGER ,
    feminicide_form1_family_son_daughter_1 INTEGER ,
    feminicide_form1_family_son_daughter_2 INTEGER ,
    feminicide_form1_family_son_daughter_3 INTEGER ,
    feminicide_form1_family_son_daughter_4 INTEGER ,
    feminicide_form1_family_son_daughter_5 INTEGER ,
    feminicide_form1_family_grandmother INTEGER ,
    feminicide_form1_family_grandfather INTEGER ,
    feminicide_form1_family_other_member CHARACTER VARYING(120) ,

    CONSTRAINT feminicide_form1_pkey PRIMARY KEY (feminicide_form1_id),
    CONSTRAINT feminicide_form1_feminicide_form1_i_code_key UNIQUE (feminicide_form1_i_code)

);


ALTER TABLE rel_victim_case_form2_enums_feminicide_form1 ADD CONSTRAINT rel_victim_case_form2_enums_feminicide_form1_id
   FOREIGN KEY ( feminicide_form1_id) REFERENCES feminicide_form1(feminicide_form1_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_victim_case_form2_enums_feminicide_form1 ADD CONSTRAINT rel_victim_case_form2_enums_feminicide_enums_id
   FOREIGN KEY ( victim_case_form2_enums_id) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE rel_victim_case_form2_enums_feminicide_risk_form1 ADD CONSTRAINT rel_victim_case_form2_enums_feminicide_risk_form1_id
   FOREIGN KEY ( feminicide_risk_form1_id) REFERENCES feminicide_risk_form1(feminicide_risk_form1_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_victim_case_form2_enums_feminicide_risk_form1 ADD CONSTRAINT rel_victim_case_form2_enums_feminicide_risk_enums_id
   FOREIGN KEY ( victim_case_form2_enums_id) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;


ALTER TABLE rel_case_owner_feminicide ADD CONSTRAINT rel_case_owner_feminicide_case_owner
    FOREIGN KEY ( rel_case_owner_feminicide_case_owner) REFERENCES case_owner(case_owner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE rel_case_owner_feminicide ADD CONSTRAINT rel_case_owner_feminicide_feminicide
    FOREIGN KEY ( rel_case_owner_feminicide_feminicide) REFERENCES feminicide(feminicide_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_case_owner_feminicide_risk ADD CONSTRAINT rel_case_owner_feminicide_risk_case_owner
    FOREIGN KEY ( rel_case_owner_feminicide_risk_case_owner) REFERENCES case_owner(case_owner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE rel_case_owner_feminicide_risk ADD CONSTRAINT rel_case_owner_feminicide_risk_feminicide_risk
    FOREIGN KEY ( rel_case_owner_feminicide_risk_feminicide_risk) REFERENCES feminicide_risk(feminicide_risk_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;



ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_zone
    FOREIGN KEY ( feminicide_form1_victim_zone) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_marital_status
    FOREIGN KEY ( feminicide_form1_victim_marital_status) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_sex
    FOREIGN KEY ( feminicide_form1_victim_sex) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_gender_identity
    FOREIGN KEY ( feminicide_form1_victim_gender_identity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_sexual_orientation
    FOREIGN KEY ( feminicide_form1_victim_sexual_orientation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_ethnicity
    FOREIGN KEY ( feminicide_form1_victim_ethnicity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_indigenous_people
    FOREIGN KEY ( feminicide_form1_victim_indigenous_people) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_special_population
    FOREIGN KEY ( feminicide_form1_victim_special_population) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_disability
    FOREIGN KEY ( feminicide_form1_victim_disability) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_disability_type
    FOREIGN KEY ( feminicide_form1_victim_disability_type) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_presumed_aggressor_relation
    FOREIGN KEY ( feminicide_form1_presumed_aggressor_relation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_presumed_aggressor_known_vgb
    FOREIGN KEY ( feminicide_form1_presumed_aggressor_known_vgb) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_s_g_s_s_s_affiliation
    FOREIGN KEY ( feminicide_form1_s_g_s_s_s_affiliation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_zone
    FOREIGN KEY ( feminicide_form1_informant_zone) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_sex
    FOREIGN KEY ( feminicide_form1_informant_sex) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_gender_identity
    FOREIGN KEY ( feminicide_form1_informant_gender_identity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_sexual_orientation
    FOREIGN KEY ( feminicide_form1_informant_sexual_orientation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_ethnicity
    FOREIGN KEY ( feminicide_form1_informant_ethnicity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_indigenous_people
    FOREIGN KEY ( feminicide_form1_informant_indigenous_people) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_migratory_status
    FOREIGN KEY ( feminicide_form1_informant_migratory_status) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_migration_situation
    FOREIGN KEY ( feminicide_form1_informant_migration_situation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_highest_education_level
    FOREIGN KEY ( feminicide_form1_informant_highest_education_level) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_current_employment
    FOREIGN KEY ( feminicide_form1_informant_current_employment) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_employment_access
    FOREIGN KEY ( feminicide_form1_informant_employment_access) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_primary_occupation
    FOREIGN KEY ( feminicide_form1_informant_primary_occupation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_disability
    FOREIGN KEY ( feminicide_form1_informant_disability) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_informant_disability_type
    FOREIGN KEY ( feminicide_form1_informant_disability_type) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_household_expense_responsibility
    FOREIGN KEY ( feminicide_form1_household_expense_responsibility) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_post_death_economic_assumption
    FOREIGN KEY ( feminicide_form1_post_death_economic_assumption) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_dependent_people
    FOREIGN KEY ( feminicide_form1_any_dependent_people) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_public_or_private_entity
    FOREIGN KEY ( feminicide_form1_any_public_or_private_entity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_public_transport_access
    FOREIGN KEY ( feminicide_form1_public_transport_access) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_preferred_transportation_mode
    FOREIGN KEY ( feminicide_form1_preferred_transportation_mode) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_transport_difficulty
    FOREIGN KEY ( feminicide_form1_transport_difficulty) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_economic_resources_for_transport
    FOREIGN KEY ( feminicide_form1_economic_resources_for_transport) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_debt_or_help_due_to_transport
    FOREIGN KEY ( feminicide_form1_debt_or_help_due_to_transport) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_transport_subsidy_received
    FOREIGN KEY ( feminicide_form1_transport_subsidy_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_safety_transportation_concern
    FOREIGN KEY ( feminicide_form1_safety_transportation_concern) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_food_access_frequency
    FOREIGN KEY ( feminicide_form1_food_access_frequency) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_aggressor_food_restriction
    FOREIGN KEY ( feminicide_form1_aggressor_food_restriction) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_food_income_support
    FOREIGN KEY ( feminicide_form1_food_income_support) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_juridical_assistance_received
    FOREIGN KEY ( feminicide_form1_juridical_assistance_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_victim_representation
    FOREIGN KEY ( feminicide_form1_victim_representation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_psychosocial_support_received
    FOREIGN KEY ( feminicide_form1_psychosocial_support_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_emergency_emotional_crisis
    FOREIGN KEY ( feminicide_form1_emergency_emotional_crisis) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_aggressor_same_residence
    FOREIGN KEY ( feminicide_form1_aggressor_same_residence) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_aggressor_location_known
    FOREIGN KEY ( feminicide_form1_aggressor_location_known) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_type_of_assistance_received
    FOREIGN KEY ( feminicide_form1_any_type_of_assistance_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_compensation_fund_insurance_coverage
    FOREIGN KEY ( feminicide_form1_compensation_fund_insurance_coverage) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
   
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_funeral_subsidy_received
    FOREIGN KEY ( feminicide_form1_funeral_subsidy_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_funeral_funds_available
    FOREIGN KEY ( feminicide_form1_funeral_funds_available) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_funeral_cost_value
    FOREIGN KEY ( feminicide_form1_funeral_cost_value) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_rename_reputation_impact
    FOREIGN KEY ( feminicide_form1_rename_reputation_impact) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_action_plan
    FOREIGN KEY ( feminicide_form1_action_plan) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_feminicide
    FOREIGN KEY ( feminicide_form1_feminicide) REFERENCES feminicide(feminicide_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_assistance_received_city_hall
    FOREIGN KEY ( feminicide_form1_any_assistance_received_city_hall) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_assistance_received_womens_office
    FOREIGN KEY ( feminicide_form1_any_assistance_received_womens_office) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_assistance_received_other_entity
    FOREIGN KEY ( feminicide_form1_any_assistance_received_other_entity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_form1 ADD CONSTRAINT feminicide_form1_any_assistance_received_other
    FOREIGN KEY ( feminicide_form1_any_assistance_received_other) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;




CREATE TABLE feminicide_risk_form1 (
    feminicide_risk_form1_id SERIAL NOT NULL,
    feminicide_risk_form1_i_code CHARACTER VARYING(36) NOT NULL,
    feminicide_risk_form1_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_risk_form1_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_risk_form1_victim_identity_name CHARACTER VARYING(32) ,
    feminicide_risk_form1_birth_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_risk_form1_victim_address CHARACTER VARYING(32) NOT NULL,
    feminicide_risk_form1_victim_living_zone BIGINT NOT NULL,
    feminicide_risk_form1_victim_living_town CHARACTER VARYING(8) NOT NULL,
    feminicide_risk_form1_victim_s_g_s_s_s_affiliation BIGINT NOT NULL,
    feminicide_risk_form1_eps_name CHARACTER VARYING(64),
    feminicide_risk_form1_contact_phone CHARACTER VARYING(10) NOT NULL,
    feminicide_risk_form1_contact_emergency_contact_names CHARACTER VARYING(64) NOT NULL,
    feminicide_risk_form1_contact_emergency_contact_number CHARACTER VARYING(10) NOT NULL,
    feminicide_risk_form1_victim_marital_status BIGINT NOT NULL,
    feminicide_risk_form1_victim_sex BIGINT NOT NULL,
    feminicide_risk_form1_victim_gender_identity BIGINT NOT NULL,
    feminicide_risk_form1_victim_sexual_orientation BIGINT NOT NULL,
    feminicide_risk_form1_victim_ethnic_affiliation BIGINT NOT NULL,
    feminicide_risk_form1_victim_indigenous_people BIGINT,
    feminicide_risk_form1_victim_is_migrant BIGINT NOT NULL,
    feminicide_risk_form1_victim_migration_status BIGINT,
    feminicide_risk_form1_victim_max_education_level BIGINT NOT NULL,
    feminicide_risk_form1_victim_is_special_population BIGINT NOT NULL,
    feminicide_risk_form1_victim_currently_has_job BIGINT NOT NULL,
    feminicide_risk_form1_victim_job_explanation TEXT ,
    feminicide_risk_form1_victim_abandoned_job_due_to_risk BIGINT NOT NULL,
    feminicide_risk_form1_victim_main_occupation BIGINT NOT NULL,
    feminicide_risk_form1_victim_has_disability BIGINT NOT NULL,
    feminicide_risk_form1_victim_disability_type BIGINT,
    feminicide_risk_form1_victim_risk_description TEXT NOT NULL,
    feminicide_risk_form1_victim_additional_info TEXT ,
    feminicide_risk_form1_victim_is_economic_provider BIGINT NOT NULL,
    feminicide_risk_form1_victim_economic_provider_explanation TEXT ,
    feminicide_risk_form1_victim_has_familiar_support BIGINT NOT NULL,
    feminicide_risk_form1_victim_familiar_support_type BIGINT NOT NULL,
    feminicide_risk_form1_victim_public_transport_access BIGINT NOT NULL,
    feminicide_risk_form1_victim_common_transport_mode BIGINT NOT NULL,
    feminicide_risk_form1_victim_transport_mode_explanation TEXT ,
    feminicide_risk_form1_victim_estimated_travel_cost INTEGER NOT NULL,
    feminicide_risk_form1_victim_difficulties_with_transport BIGINT NOT NULL,
    feminicide_risk_form1_victim_difficulties_explanation TEXT ,
    feminicide_risk_form1_victim_economic_resources_for_transport BIGINT NOT NULL,
    feminicide_risk_form1_victim_economic_resources_explanation TEXT ,
    feminicide_risk_form1_victim_has_debt_or_help_due_to_transport BIGINT NOT NULL,
    feminicide_risk_form1_victim_debt_impact_details TEXT ,
    feminicide_risk_form1_victim_received_transport_subsidy BIGINT NOT NULL,
    feminicide_risk_form1_victim_transport_subsidy_explanation TEXT ,
    feminicide_risk_form1_victim_safety_avoided_transport BIGINT NOT NULL,
    feminicide_risk_form1_victim_safety_avoided_explanation TEXT ,
    feminicide_risk_form1_victim_food_access_frequency BIGINT NOT NULL,
    feminicide_risk_form1_victim_food_access_explanation TEXT ,
    feminicide_risk_form1_victim_agressor_food_restriction BIGINT NOT NULL,
    feminicide_risk_form1_victim_agressor_food_restrict_explanat TEXT ,
    feminicide_risk_form1_victim_family_fixed_income BIGINT NOT NULL,
    feminicide_risk_form1_victim_family_fixed_income_explanation TEXT ,
    feminicide_risk_form1_victim_is_only_provider_for_food BIGINT NOT NULL,
    feminicide_risk_form1_victim_is_only_provider_explanation TEXT ,
    feminicide_risk_form1_victim_juridical_assistance_received BIGINT NOT NULL,
    feminicide_risk_form1_victim_juridical_assistance_explain TEXT ,
    feminicide_risk_form1_victim_wants_juridical_assistance BIGINT NOT NULL,
    feminicide_risk_form1_victim_representations_of_victims BIGINT NOT NULL,
    feminicide_risk_form1_victim_representations_explanation TEXT ,
    feminicide_risk_form1_victim_psychosocial_support_received BIGINT NOT NULL,
    feminicide_risk_form1_victim_psychosocial_supp_received_explain TEXT ,
    feminicide_risk_form1_victim_urgent_emotional_crisis BIGINT NOT NULL,
    feminicide_risk_form1_victim_urgent_crisis_explanation TEXT ,
    feminicide_risk_form1_aggressor_same_residence BIGINT NOT NULL,
    feminicide_risk_form1_aggressor_same_residence_explanation TEXT ,
    feminicide_risk_form1_aggressor_knows_victim_location BIGINT NOT NULL,
    feminicide_risk_form1_aggressor_knows_victim_loc_explanation TEXT ,
    feminicide_risk_form1_victim_housing_help_received BIGINT NOT NULL,
    feminicide_risk_form1_victim_housing_help_explain TEXT ,
    feminicide_risk_form1_victim_abandon_clothing BIGINT NOT NULL,
    feminicide_risk_form1_victim_abandon_clothing_explain TEXT ,
    feminicide_risk_form1_victim_clothing_help_received BIGINT NOT NULL,
    feminicide_risk_form1_victim_clothing_help_explain TEXT ,
    feminicide_risk_form1_victim_other_needs TEXT ,
    feminicide_risk_form1_interview_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feminicide_risk_form1_summary TEXT ,
    feminicide_risk_form1_informant_indigenous_people BIGINT ,

    feminicide_risk_form1_any_assistance_received_city_hall BIGINT ,
    feminicide_risk_form1_any_assistance_received_womens_office BIGINT ,
    feminicide_risk_form1_any_assistance_received_other_entity BIGINT ,
    feminicide_risk_form1_any_assistance_received_other BIGINT ,
    feminicide_risk_form1_assistance_received CHARACTER VARYING(120) ,
    feminicide_risk_form1_assistance_received_city_hall CHARACTER VARYING(120) ,
    feminicide_risk_form1_assistance_received_womens_office CHARACTER VARYING(120) ,
    feminicide_risk_form1_assistance_received_other_entity CHARACTER VARYING(120) ,
        feminicide_risk_form1_assistance_received_other CHARACTER VARYING(120) ,

    feminicide_risk_form1_family_father INTEGER ,
    feminicide_risk_form1_family_mother INTEGER ,
    feminicide_risk_form1_family_stepfather INTEGER ,
    feminicide_risk_form1_family_stepmother INTEGER ,
    feminicide_risk_form1_family_partner INTEGER ,
    feminicide_risk_form1_family_sibling_1 INTEGER ,
    feminicide_risk_form1_family_sibling_2 INTEGER ,
    feminicide_risk_form1_family_sibling_3 INTEGER ,
    feminicide_risk_form1_family_sibling_4 INTEGER ,
    feminicide_risk_form1_family_sibling_5 INTEGER ,
    feminicide_risk_form1_family_son_daughter_1 INTEGER ,
    feminicide_risk_form1_family_son_daughter_2 INTEGER ,
    feminicide_risk_form1_family_son_daughter_3 INTEGER ,
    feminicide_risk_form1_family_son_daughter_4 INTEGER ,
    feminicide_risk_form1_family_son_daughter_5 INTEGER ,
    feminicide_risk_form1_family_grandmother INTEGER ,
    feminicide_risk_form1_family_grandfather INTEGER ,
    feminicide_risk_form1_family_other_member CHARACTER VARYING(120) ,


    feminicide_form1_feminicide_risk BIGINT NOT NULL,
    CONSTRAINT feminicide_risk_form1_pkey PRIMARY KEY (feminicide_risk_form1_id),
    CONSTRAINT feminicide_risk_form1_feminicide_risk_form1_i_code_key UNIQUE (feminicide_risk_form1_i_code)
);


ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_zone
   FOREIGN KEY ( feminicide_risk_form1_victim_zone) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_s_g_s_s_s_affiliation
   FOREIGN KEY ( feminicide_risk_form1_victim_s_g_s_s_s_affiliation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_marital_status
   FOREIGN KEY ( feminicide_risk_form1_victim_marital_status) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_sex
   FOREIGN KEY ( feminicide_risk_form1_victim_sex) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_gender_identity
   FOREIGN KEY ( feminicide_risk_form1_victim_gender_identity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_sexual_orientation
   FOREIGN KEY ( feminicide_risk_form1_victim_sexual_orientation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_ethnic_affiliation
   FOREIGN KEY ( feminicide_risk_form1_victim_ethnic_affiliation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_is_migrant
   FOREIGN KEY ( feminicide_risk_form1_victim_is_migrant) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_migration_status
   FOREIGN KEY ( feminicide_risk_form1_victim_migration_status) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_max_education_level
   FOREIGN KEY ( feminicide_risk_form1_victim_max_education_level) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_is_special_population
   FOREIGN KEY ( feminicide_risk_form1_victim_is_special_population) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_currently_has_job
   FOREIGN KEY ( feminicide_risk_form1_victim_currently_has_job) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_abandoned_job_due_to_risk
   FOREIGN KEY ( feminicide_risk_form1_victim_abandoned_job_due_to_risk) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_main_occupation
   FOREIGN KEY ( feminicide_risk_form1_victim_main_occupation) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_has_disability
   FOREIGN KEY ( feminicide_risk_form1_victim_has_disability) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_disability_type
   FOREIGN KEY ( feminicide_risk_form1_victim_disability_type) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_received_help_or_subsidy
   FOREIGN KEY ( feminicide_risk_form1_victim_received_help_or_subsidy) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_is_economic_provider
   FOREIGN KEY ( feminicide_risk_form1_victim_is_economic_provider) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_has_familiar_support
   FOREIGN KEY ( feminicide_risk_form1_victim_has_familiar_support) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_familiar_support_type
   FOREIGN KEY ( feminicide_risk_form1_victim_familiar_support_type) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_public_transport_access
   FOREIGN KEY ( feminicide_risk_form1_victim_public_transport_access) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_common_transport_mode
   FOREIGN KEY ( feminicide_risk_form1_victim_common_transport_mode) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_difficulties_with_transport
   FOREIGN KEY ( feminicide_risk_form1_victim_difficulties_with_transport) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_economic_resources_for_transport
   FOREIGN KEY ( feminicide_risk_form1_victim_economic_resources_for_transport) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_has_debt_or_help_due_to_transport
   FOREIGN KEY ( feminicide_risk_form1_victim_has_debt_or_help_due_to_transport) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_received_transport_subsidy
   FOREIGN KEY ( feminicide_risk_form1_victim_received_transport_subsidy) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_safety_avoided_transport
   FOREIGN KEY ( feminicide_risk_form1_victim_safety_avoided_transport) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_food_access_frequency
   FOREIGN KEY ( feminicide_risk_form1_victim_food_access_frequency) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_agressor_food_restriction
   FOREIGN KEY ( feminicide_risk_form1_victim_agressor_food_restriction) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_family_fixed_income
   FOREIGN KEY ( feminicide_risk_form1_victim_family_fixed_income) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_is_only_provider_for_food
   FOREIGN KEY ( feminicide_risk_form1_victim_is_only_provider_for_food) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_juridical_assistance_received
   FOREIGN KEY ( feminicide_risk_form1_victim_juridical_assistance_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_wants_juridical_assistance
   FOREIGN KEY ( feminicide_risk_form1_victim_wants_juridical_assistance) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_representations_of_victims
   FOREIGN KEY ( feminicide_risk_form1_victim_representations_of_victims) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_psychosocial_support_received
   FOREIGN KEY ( feminicide_risk_form1_victim_psychosocial_support_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_urgent_emotional_crisis
   FOREIGN KEY ( feminicide_risk_form1_victim_urgent_emotional_crisis) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_aggressor_same_residence
   FOREIGN KEY ( feminicide_risk_form1_aggressor_same_residence) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_aggressor_knows_victim_location
   FOREIGN KEY ( feminicide_risk_form1_aggressor_knows_victim_location) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_housing_help_received
   FOREIGN KEY ( feminicide_risk_form1_victim_housing_help_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_abandon_clothing
   FOREIGN KEY ( feminicide_risk_form1_victim_abandon_clothing) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;
    
ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_victim_clothing_help_received
    FOREIGN KEY ( feminicide_risk_form1_victim_clothing_help_received) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_form1_feminicide_risk
    FOREIGN KEY ( feminicide_form1_feminicide_risk) REFERENCES feminicide_risk(feminicide_risk_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_any_assistance_received_city_hall
   FOREIGN KEY ( feminicide_risk_form1_any_assistance_received_city_hall) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_any_assistance_received_womens_office
   FOREIGN KEY ( feminicide_risk_form1_any_assistance_received_womens_office) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE feminicide_risk_form1 ADD CONSTRAINT feminicide_risk_form1_any_assistance_received_other_entity
   FOREIGN KEY ( feminicide_risk_form1_any_assistance_received_other_entity) REFERENCES victim_case_form2_enums(victim_case_form2_enums_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;




GRANT SELECT, UPDATE, USAGE ON ALL SEQUENCES IN SCHEMA salvia TO salvia_admin;
GRANT SELECT, UPDATE, USAGE ON ALL SEQUENCES IN SCHEMA security TO salvia_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA salvia TO salvia_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA security TO salvia_admin;