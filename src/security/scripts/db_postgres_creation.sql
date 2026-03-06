-- Schema de seguridad
SET search_path = security;

CREATE TABLE general_user (
    general_user_id SERIAL NOT NULL,
    general_user_i_code CHARACTER VARYING(36) NOT NULL,
    general_user_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    general_user_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    general_user_login CHARACTER VARYING(128) NOT NULL,
    general_user_password CHARACTER VARYING(512) NOT NULL,
    general_user_status CHARACTER VARYING(10) NOT NULL,
    general_user_language CHARACTER VARYING(2) NOT NULL,
    general_user_general_user_profile BIGINT ,
    CONSTRAINT general_user_pkey PRIMARY KEY (general_user_id),
    CONSTRAINT general_user_general_user_i_code_key UNIQUE (general_user_i_code),
    CONSTRAINT general_user_general_user_login_key UNIQUE (general_user_login)
);

CREATE TABLE reset_password (
    reset_password_id SERIAL NOT NULL,
    reset_password_i_code CHARACTER VARYING(36) NOT NULL,
    reset_password_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    reset_password_general_user BIGINT NOT NULL,
    CONSTRAINT reset_password_pkey PRIMARY KEY (reset_password_id),
    CONSTRAINT reset_password_reset_password_i_code_key UNIQUE (reset_password_i_code)
);


CREATE TABLE general_user_profile (
    general_user_profile_id SERIAL NOT NULL,
    general_user_profile_i_code CHARACTER VARYING(36) NOT NULL,
    general_user_profile_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    general_user_profile_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    general_user_profile_gender CHARACTER VARYING(2) ,
    general_user_profile_nick CHARACTER VARYING(120) ,
    general_user_profile_description CHARACTER VARYING(1024) ,
    general_user_profile_names CHARACTER VARYING(128) NOT NULL,
    general_user_profile_last_names CHARACTER VARYING(128) ,
    general_user_profile_doc_type CHARACTER VARYING(2) ,
    general_user_profile_doc_number CHARACTER VARYING(32) ,
    general_user_profile_town CHARACTER VARYING(8) ,
    CONSTRAINT general_user_profile_pkey PRIMARY KEY (general_user_profile_id),
    CONSTRAINT general_user_profile_general_user_profile_i_code_key UNIQUE (general_user_profile_i_code)
);

CREATE TABLE role (
    role_id SERIAL NOT NULL,
    role_i_code CHARACTER VARYING(36) NOT NULL,
    role_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    role_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    role_code CHARACTER VARYING(32) NOT NULL,
    role_name CHARACTER VARYING(128) NOT NULL,
    role_description CHARACTER VARYING(256) NOT NULL,
    CONSTRAINT role_pkey PRIMARY KEY (role_id),
    CONSTRAINT role_role_i_code_key UNIQUE (role_i_code),
    CONSTRAINT role_role_code_key UNIQUE (role_code)
);

CREATE TABLE rel_role_general_user (
    rel_role_general_user_id SERIAL NOT NULL,
    rel_role_general_user_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    role_id BIGINT NOT NULL,
    general_user_id BIGINT NOT NULL,
    CONSTRAINT rel_role_general_user_pkey PRIMARY KEY (rel_role_general_user_id)
);

CREATE TABLE email (
    email_id SERIAL NOT NULL,
    e_mail_i_code CHARACTER VARYING(36) NOT NULL,
    e_mail_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    e_mail_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    e_mail_data CHARACTER VARYING(128) NOT NULL,
    e_mail_general_user_profile BIGINT NOT NULL,
    CONSTRAINT email_pkey PRIMARY KEY (email_id),
    CONSTRAINT email_e_mail_i_code_key UNIQUE (e_mail_i_code)
);

CREATE TABLE phone_number (
    phone_number_id SERIAL NOT NULL,
    phone_number_i_code CHARACTER VARYING(36) NOT NULL,
    phone_number_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    phone_number_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    phone_number_data CHARACTER VARYING(16) NOT NULL,
    phone_number_general_user_profile BIGINT NOT NULL,
    CONSTRAINT phone_number_pkey PRIMARY KEY (phone_number_id),
    CONSTRAINT phone_number_phone_number_i_code_key UNIQUE (phone_number_i_code)
);

CREATE TABLE town (
    town_id SERIAL NOT NULL,
    town_i_code CHARACTER VARYING(36) NOT NULL,
    town_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    town_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    town_code CHARACTER VARYING(5) NOT NULL,
    town_name CHARACTER VARYING(128) NOT NULL,
    town_type CHARACTER VARYING(2) NOT NULL,
    town_latitude DOUBLE PRECISION,
    town_longitude DOUBLE PRECISION,
    city_id BIGINT ,
    CONSTRAINT town_pkey PRIMARY KEY (town_id),
    CONSTRAINT town_town_i_code_key UNIQUE (town_i_code),
    CONSTRAINT town_town_code_key UNIQUE (town_code)
);


CREATE TABLE city (
    city_id SERIAL NOT NULL,
    city_i_code CHARACTER VARYING(36) NOT NULL,
    city_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    city_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    city_code CHARACTER VARYING(5) NOT NULL,
    city_name CHARACTER VARYING(64) NOT NULL,
    department_id BIGINT ,
    CONSTRAINT city_pkey PRIMARY KEY (city_id),
    CONSTRAINT city_city_i_code_key UNIQUE (city_i_code),
    CONSTRAINT city_city_code_key UNIQUE (city_code)
);

CREATE TABLE department (
    department_id SERIAL NOT NULL,
    department_i_code CHARACTER VARYING(36) NOT NULL,
    department_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    department_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    department_code CHARACTER VARYING(2) NOT NULL,
    department_name CHARACTER VARYING(64) NOT NULL,
    country_id BIGINT ,
    CONSTRAINT department_pkey PRIMARY KEY (department_id),
    CONSTRAINT department_department_i_code_key UNIQUE (department_i_code),
    CONSTRAINT department_department_code_key UNIQUE (department_code),
    CONSTRAINT department_department_name_key UNIQUE (department_name)
);

CREATE TABLE country (
    country_id SERIAL NOT NULL,
    country_i_code CHARACTER VARYING(36) NOT NULL,
    country_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    country_update_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    country_code CHARACTER VARYING(2) NOT NULL,
    country_name CHARACTER VARYING(32) NOT NULL,
    CONSTRAINT country_pkey PRIMARY KEY (country_id),
    CONSTRAINT country_country_i_code_key UNIQUE (country_i_code),
    CONSTRAINT country_country_code_key UNIQUE (country_code),
    CONSTRAINT country_country_name_key UNIQUE (country_name)
);

CREATE INDEX general_user_login_idx ON security.general_user(general_user_login);
CREATE INDEX general_user_i_code_idx ON security.general_user(general_user_i_code);

ALTER TABLE general_user 
    ADD CONSTRAINT general_user_general_user_profile
        FOREIGN KEY ( general_user_general_user_profile) 
            REFERENCES general_user_profile(general_user_profile_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE reset_password ADD CONSTRAINT reset_password_general_user
    FOREIGN KEY ( reset_password_general_user) REFERENCES general_user(general_user_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE general_user_profile
    ADD CONSTRAINT general_user_profile_city
        FOREIGN KEY ( general_user_profile_city) 
            REFERENCES city(city_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE rel_role_general_user 
    ADD CONSTRAINT rel_general_user_id
        FOREIGN KEY ( general_user_id)
            REFERENCES general_user(general_user_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;



ALTER TABLE rel_role_general_user 
    ADD CONSTRAINT rel_role_id
        FOREIGN KEY ( role_id)
            REFERENCES role(role_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

            
ALTER TABLE email 
    ADD CONSTRAINT general_user_profile
        FOREIGN KEY ( general_user_profile) 
            REFERENCES general_user_profile(general_user_profile_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE phone_number
    ADD CONSTRAINT general_user_profile
        FOREIGN KEY ( general_user_profile) 
            REFERENCES general_user_profile(general_user_profile_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE city 
    ADD CONSTRAINT city_department
        FOREIGN KEY ( department_id)
            REFERENCES department(department_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;

ALTER TABLE town 
    ADD CONSTRAINT town_city
        FOREIGN KEY ( city_id) REFERENCES city(city_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION ;