package repository

import "fmt"

// FieldKey identifica una pregunta del form "Registro de Caso" por su posición
// (orden de sección, orden de pregunta dentro de la sección) tal como quedaron
// creadas por src/cmd/seed/seed_registro_caso.sql. Se usa posición en vez de
// descripción porque varias preguntas de Tamizaje (pareja / no-pareja) comparten
// texto idéntico pero mapean a columnas distintas (ej. "¿La violencia fue
// motivada por género?" aparece en Sección 3 Y en Tamizaje no-pareja).
func FieldKey(sectionOrder, questionOrder int) string {
	return fmt.Sprintf("S%dQ%d", sectionOrder, questionOrder)
}

// victimCaseFormFields mapea cada pregunta (por FieldKey) a su columna en
// victim_case_form2 y al tipo de conversión necesaria. Las preguntas que no
// aparecen aquí (Nombres, Apellidos, Tipo/Número de documento, las 3 cadenas
// de Departamento/Ciudad que solo alimentan cascadas de UI, el banner de
// Accesibilidad, el banner de riesgo, y Municipio de atención) se manejan por
// fuera de este mapeo genérico -- ver CreateDraft y ActivateFromAnswers.
var victimCaseFormFields = []VictimCaseFormField{ //nolint:gochecknoglobals
	// ── Sección 1 — Datos de la Víctima ──────────────────────────────────────
	{FieldKey(1, 3), "victim_case_form2_identity_name", KindText},
	{FieldKey(1, 4), "victim_case_form2_victim_phone", KindInt},
	{FieldKey(1, 7), "victim_case_form2_residence_zone", KindEnum1},
	{FieldKey(1, 10), "victim_case_form2_residence_town", KindText},
	{FieldKey(1, 11), "victim_case_form2_residence_address", KindText},
	{FieldKey(1, 13), "victim_case_form2_person_with_disability", KindBoolEnum},
	{FieldKey(1, 14), "", KindEnumN}, // Ajustes razonables VBG
	{FieldKey(1, 15), "victim_case_form2_require_language_interpreter", KindBoolEnum},
	{FieldKey(1, 16), "victim_case_form2_language_assistance", KindText},

	// ── Sección 2 — Contacto de Apoyo ────────────────────────────────────────
	{FieldKey(2, 1), "victim_case_form2_support_contact_names", KindText},
	{FieldKey(2, 2), "victim_case_form2_support_contact_phone", KindInt},
	{FieldKey(2, 3), "victim_case_form2_support_contact_email", KindText},
	{FieldKey(2, 4), "victim_case_form2_support_contact_kinship", KindEnum1},

	// ── Sección 3 — Hechos ────────────────────────────────────────────────────
	{FieldKey(3, 1), "victim_case_form2_facts_description", KindText},
	{FieldKey(3, 2), "victim_case_form2_facts_date", KindDate},
	{FieldKey(3, 3), "victim_case_form2_facts_start_time", KindTime},
	{FieldKey(3, 4), "victim_case_form2_facts_zone", KindEnum1},
	{FieldKey(3, 7), "victim_case_form2_facts_town_code", KindText},
	{FieldKey(3, 8), "victim_case_form2_facts_address", KindText},
	{FieldKey(3, 9), "victim_case_form2_scenario_violence", KindEnum1},
	{FieldKey(3, 10), "", KindEnumN}, // Tipo de violencia experimentada
	{FieldKey(3, 11), "", KindEnumN}, // Subtipo de violencia experimentada
	{FieldKey(3, 12), "", KindEnumN}, // Ámbito de la violencia
	{FieldKey(3, 13), "victim_case_form2_workplace_sector_occurrence", KindEnum1},
	{FieldKey(3, 14), "victim_case_form2_violence_motivated_by_gender", KindBoolEnum},
	{FieldKey(3, 15), "victim_case_form2_reported_previously", KindBoolEnum},
	{FieldKey(3, 16), "", KindEnumN}, // A quién denunció
	{FieldKey(3, 17), "victim_case_form2_attention_was_appropriate", KindBoolEnum},
	{FieldKey(3, 18), "victim_case_form2_recurrence_aggression", KindEnum1},

	// ── Sección 4 — Agresor ───────────────────────────────────────────────────
	{FieldKey(4, 1), "victim_case_form2_num_agressors", KindEnum1},
	{FieldKey(4, 2), "victim_case_form2_proximity_principal_aggressor", KindEnum1},
	{FieldKey(4, 3), "victim_case_form2_relationship_with_presumed_aggressor", KindEnum1},
	{FieldKey(4, 4), "victim_case_form2_aggressors_occupation", KindEnum1}, // plural -- confirmado en BD
	{FieldKey(4, 5), "victim_case_form2_economically_dependent", KindBoolEnum},
	{FieldKey(4, 6), "victim_case_form2_aggressor_gender_identity", KindEnum1},
	{FieldKey(4, 7), "victim_case_form2_aggressor_names", KindText},
	{FieldKey(4, 8), "victim_case_form2_aggressor_doc_type", KindEnum1},
	{FieldKey(4, 9), "victim_case_form2_aggressor_doc_number", KindText},
	{FieldKey(4, 10), "victim_case_form2_aggressor_address", KindText},
	{FieldKey(4, 11), "victim_case_form2_aggressor_phone", KindInt},

	// ── Sección 5 — Tamizaje ──────────────────────────────────────────────────
	// Comunes (6)
	{FieldKey(5, 1), "victim_case_form2_aggressor_violence_physical_increase", KindBoolEnum},
	{FieldKey(5, 2), "victim_case_form2_aggressor_weapon_used", KindBoolEnum},
	{FieldKey(5, 3), "victim_case_form2_aggressor_threat_kill", KindBoolEnum},
	{FieldKey(5, 4), "victim_case_form2_aggressor_pursues_spies_destroys", KindBoolEnum},
	{FieldKey(5, 5), "victim_case_form2_aggressor_capable_of_killing", KindBoolEnum},
	{FieldKey(5, 6), "victim_case_form2_aggressor_has_access_to_weapons", KindBoolEnum},
	// Pareja íntima (18) -- order 7..24
	{FieldKey(5, 7), "victim_case_form2_partner_unemployed", KindBoolEnum},
	{FieldKey(5, 8), "victim_case_form2_partner_other_denunciations", KindBoolEnum},
	{FieldKey(5, 9), "victim_case_form2_aggressor_has_penal_background", KindBoolEnum},
	{FieldKey(5, 10), "victim_case_form2_stopped_seeking_help", KindBoolEnum},
	{FieldKey(5, 11), "victim_case_form2_aggressor_forced_sex", KindBoolEnum},
	{FieldKey(5, 12), "victim_case_form2_aggressor_attempted_strangulation", KindBoolEnum},
	{FieldKey(5, 13), "victim_case_form2_aggressor_consumes_drugs", KindBoolEnum},
	{FieldKey(5, 14), "victim_case_form2_aggressor_is_alcoholic", KindBoolEnum},
	{FieldKey(5, 15), "victim_case_form2_partner_controls", KindBoolEnum},
	{FieldKey(5, 16), "victim_case_form2_aggressor_had_hit_in_vulnerability", KindBoolEnum},
	{FieldKey(5, 17), "victim_case_form2_victim_health_to_blackmail", KindBoolEnum},
	{FieldKey(5, 18), "victim_case_form2_partner_threatened_suicide", KindBoolEnum},
	{FieldKey(5, 19), "victim_case_form2_partner_threatened_damage_members", KindBoolEnum},
	{FieldKey(5, 20), "victim_case_form2_thoughts_of_self_harm", KindBoolEnum},
	{FieldKey(5, 21), "victim_case_form2_aggressor_limits_contact_support_networks", KindBoolEnum},
	{FieldKey(5, 22), "victim_case_form2_threatened_reveal_sexual_orientation", KindBoolEnum},
	{FieldKey(5, 23), "victim_case_form2_still_lives_with_aggressor", KindBoolEnum},
	{FieldKey(5, 24), "victim_case_form2_aggressor_violently_jealous", KindBoolEnum},
	// No pareja (14) -- order 25..38
	{FieldKey(5, 25), "victim_case_form2_aggressor_taken_advantage_physical_vulnerabil", KindBoolEnum}, // truncado a 63 chars en BD
	{FieldKey(5, 26), "victim_case_form2_aggressor_unemployed", KindBoolEnum},
	{FieldKey(5, 27), "victim_case_form2_aggressor_has_penal_background_2", KindBoolEnum},
	{FieldKey(5, 28), "victim_case_form2_aggressor_sexually_harassment", KindBoolEnum},
	{FieldKey(5, 29), "victim_case_form2_aggressor_sexually_harassment_2", KindBoolEnum},
	{FieldKey(5, 30), "victim_case_form2_violence_motivated_by_gender_2", KindBoolEnum},
	{FieldKey(5, 31), "victim_case_form2_aggressor_use_drugs", KindBoolEnum},
	{FieldKey(5, 32), "victim_case_form2_aggressor_is_alcoholic_2", KindBoolEnum},
	{FieldKey(5, 33), "victim_case_form2_aggressor_controls", KindBoolEnum},
	{FieldKey(5, 34), "victim_case_form2_aggressor_threatened_damage_members", KindBoolEnum},
	{FieldKey(5, 35), "victim_case_form2_thoughts_of_self_harm_2", KindBoolEnum},
	{FieldKey(5, 36), "victim_case_form2_aggressor_common_spaces", KindBoolEnum},
	{FieldKey(5, 37), "victim_case_form2_aggressor_hierarchy", KindBoolEnum},
	{FieldKey(5, 38), "victim_case_form2_aggressor_used_position_authority", KindBoolEnum},
	// order 39 = banner de riesgo (info, sin respuesta) -- omitido

	// ── Sección 6 — Datos Personales ─────────────────────────────────────────
	{FieldKey(6, 1), "victim_case_form2_birth_date", KindDate},
	{FieldKey(6, 2), "victim_case_form2_physical_mental_sensory_difficulties", KindBoolEnum},
	{FieldKey(6, 3), "victim_case_form2_activities_unable_to_hear", KindScale},
	{FieldKey(6, 4), "victim_case_form2_activities_unable_to_talk", KindScale},
	{FieldKey(6, 5), "victim_case_form2_activities_unable_to_see", KindScale},
	{FieldKey(6, 6), "victim_case_form2_activities_unable_to_move", KindScale},
	{FieldKey(6, 7), "victim_case_form2_activities_unable_to_take", KindScale},
	{FieldKey(6, 8), "victim_case_form2_activities_unable_to_understand", KindScale},
	{FieldKey(6, 9), "victim_case_form2_activities_unable_to_eat", KindScale},
	{FieldKey(6, 10), "victim_case_form2_activities_unable_to_interact", KindScale},
	{FieldKey(6, 11), "victim_case_form2_activities_unable_to_do_everyday", KindScale},
	{FieldKey(6, 12), "", KindEnumN}, // Ley 1996
	{FieldKey(6, 13), "victim_case_form2_nationality", KindEnum1},
	{FieldKey(6, 14), "victim_case_form2_specified_nationality", KindEnum1},
	{FieldKey(6, 15), "victim_case_form2_migration_condition", KindEnum1},
	{FieldKey(6, 16), "victim_case_form2_gender_identity", KindEnum1},
	{FieldKey(6, 17), "victim_case_form2_sexual_orientation", KindEnum1},
	{FieldKey(6, 18), "victim_case_form2_assigned_sex_at_birth", KindEnum1},
	{FieldKey(6, 19), "", KindEnumN}, // Población especialmente protegida
	{FieldKey(6, 20), "victim_case_form2_ethnic_affiliation", KindEnum1},
	{FieldKey(6, 21), "victim_case_form2_indigenous_people", KindEnum1},
	{FieldKey(6, 22), "victim_case_form2_campesino_recognition", KindBoolEnum},
	{FieldKey(6, 23), "victim_case_form2_marital_status", KindEnum1},
	{FieldKey(6, 24), "victim_case_form2_last_education_level", KindEnum1},
	{FieldKey(6, 25), "victim_case_form2_occupation", KindEnum1},
	{FieldKey(6, 26), "victim_case_form2_income_generation_method", KindEnum1},
	{FieldKey(6, 27), "", KindEnumN}, // Modalidad ASP
	{FieldKey(6, 28), "victim_case_form2_employment_relationship", KindEnum1},
	{FieldKey(6, 29), "victim_case_form2_approx_start_asp", KindTimestamp},
	{FieldKey(6, 30), "", KindEnumN}, // Razón ASP
	{FieldKey(6, 31), "victim_case_form2_housing_tenancy_form", KindEnum1},
	{FieldKey(6, 32), "victim_case_form2_housing_stratum", KindEnum1},
	{FieldKey(6, 33), "", KindEnumN}, // Personas a cargo
	{FieldKey(6, 34), "victim_case_form2_currently_pregnant", KindBoolEnum},

	// ── Sección 7 — Plan de Acción ────────────────────────────────────────────
	{FieldKey(7, 1), "", KindEnumN}, // Plan de acción
	{FieldKey(7, 2), "victim_case_form2_saliva_management_explanation", KindText},

	// ── Sección 8 — Denuncia Fácil ────────────────────────────────────────────
	{FieldKey(8, 1), "victim_case_form2_allows_easy_report", KindBoolEnum},

	// Sección 9 (Lugar de Atención) no proyecta sobre victim_case_form2 --
	// Municipio de atención (S9Q3) se maneja aparte en ActivateFromAnswers
	// para actualizar victim_case_victim_town_code.
}
