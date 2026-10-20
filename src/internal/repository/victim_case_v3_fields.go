package repository

import "fmt"

// victimCaseV3FormFields mapea cada pregunta de Registro V3 (por FieldKey
// "{MODULO}_{orden}", módulo = CI/DP/CA/HV/AG/TZ/GS, o TZP/TZNP para las
// ramas de Tamizaje) a su columna en victim_case_form2. Mismo patrón que
// victimCaseFormFields (V2).
//
// Fuente: DocsMD/Screens/Registro V3/registro-v3-secciones.md. No usa
// dinamic-form (decisión A1) — este mapeo lo consume directamente
// victim_case_v3_repository.go via FieldKeyV3(modulo, orden).

// FieldKeyV3 identifica una pregunta por su módulo (CI/DP/CA/HV/AG/TZ/GS,
// o TZP/TZNP para las ramas de Tamizaje) y orden dentro del módulo, igual
// que secciones.md.
func FieldKeyV3(modulo string, orden int) string {
	return fmt.Sprintf("%s_%d", modulo, orden)
}

var victimCaseV3FormFields = []VictimCaseFormField{ //nolint:gochecknoglobals
	// ── CI — Contacto e identificación ──────────────────────────────────────
	// Q2 nombres, Q5 tipo doc, Q6 núm doc se manejan aparte (bootstrap de
	// victim_case, igual que V2) — ver CreateDraftV3/UpdateDraftCoreFieldsV3.
	{FieldKeyV3("CI", 1), "victim_case_form2_victim_pronoun", KindEnum1},
	{FieldKeyV3("CI", 3), "victim_case_form2_victim_identity_name", KindText},
	{FieldKeyV3("CI", 4), "victim_case_form2_victim_phone", KindInt},
	{FieldKeyV3("CI", 7), "victim_case_form2_person_with_disability", KindBoolEnum},
	{FieldKeyV3("CI", 8), "", KindEnumN}, // Ajustes razonables VBG
	{FieldKeyV3("CI", 9), "victim_case_form2_require_language_interpreter", KindBoolEnum},

	// ── DP — Datos personales (incluye residencia) ──────────────────────────
	{FieldKeyV3("DP", 1), "victim_case_form2_birth_date", KindDate},
	{FieldKeyV3("DP", 2), "victim_case_form2_physical_mental_sensory_difficulties", KindBoolEnum},
	{FieldKeyV3("DP", 3), "", KindEnumN}, // actividades con dificultad (escala) -- ver nota abajo
	{FieldKeyV3("DP", 4), "", KindEnumN}, // Ley 1996
	{FieldKeyV3("DP", 5), "victim_case_form2_nationality", KindEnum1},
	{FieldKeyV3("DP", 50), "victim_case_form2_victim_email", KindText}, // Q5b correo víctima (obs. #12)
	{FieldKeyV3("DP", 6), "victim_case_form2_specified_nationality", KindEnum1},
	{FieldKeyV3("DP", 7), "victim_case_form2_migration_condition", KindEnum1},
	{FieldKeyV3("DP", 8), "victim_case_form2_gender_identity", KindEnum1},
	{FieldKeyV3("DP", 9), "victim_case_form2_sexual_orientation", KindEnum1},
	{FieldKeyV3("DP", 10), "victim_case_form2_assigned_sex_at_birth", KindEnum1},
	{FieldKeyV3("DP", 11), "", KindEnumN}, // Población especialmente protegida
	{FieldKeyV3("DP", 12), "victim_case_form2_ethnic_affiliation", KindEnum1},
	{FieldKeyV3("DP", 13), "victim_case_form2_indigenous_people", KindEnum1},
	{FieldKeyV3("DP", 14), "victim_case_form2_campesino_recognition", KindBoolEnum},
	{FieldKeyV3("DP", 15), "victim_case_form2_marital_status", KindEnum1},
	{FieldKeyV3("DP", 16), "victim_case_form2_last_education_level", KindEnum1},
	{FieldKeyV3("DP", 17), "victim_case_form2_occupation", KindEnum1},
	{FieldKeyV3("DP", 18), "victim_case_form2_income_generation_method", KindEnum1},
	{FieldKeyV3("DP", 19), "victim_case_form2_employment_relationship", KindEnum1},
	{FieldKeyV3("DP", 20), "", KindEnumN}, // Modalidad ASP
	{FieldKeyV3("DP", 21), "victim_case_form2_approx_start_asp", KindTimestamp},
	{FieldKeyV3("DP", 22), "", KindEnumN}, // Razón ASP
	{FieldKeyV3("DP", 23), "victim_case_form2_housing_tenancy_form", KindEnum1},
	{FieldKeyV3("DP", 24), "victim_case_form2_housing_stratum", KindEnum1},
	{FieldKeyV3("DP", 25), "", KindEnumN}, // Personas a cargo
	{FieldKeyV3("DP", 26), "victim_case_form2_currently_pregnant", KindBoolEnum},
	{FieldKeyV3("DP", 27), "", KindText},                                // Departamento residencia -- cascada UI, no se proyecta
	{FieldKeyV3("DP", 28), "victim_case_form2_residence_town", KindText}, // bootstrap también (ver CreateDraftV3)
	{FieldKeyV3("DP", 29), "victim_case_form2_residence_address", KindText},
	{FieldKeyV3("DP", 30), "victim_case_form2_residence_zone", KindEnum1},

	// ── CA — Contacto de apoyo ───────────────────────────────────────────────
	{FieldKeyV3("CA", 1), "victim_case_form2_support_contact_names", KindText},
	{FieldKeyV3("CA", 2), "victim_case_form2_support_contact_phone", KindInt},
	{FieldKeyV3("CA", 3), "victim_case_form2_support_contact_email", KindText},
	{FieldKeyV3("CA", 4), "victim_case_form2_support_contact_kinship", KindEnum1},

	// ── HV — Hechos de violencia ─────────────────────────────────────────────
	{FieldKeyV3("HV", 1), "victim_case_form2_facts_description", KindText},
	{FieldKeyV3("HV", 2), "victim_case_form2_facts_date", KindDate},
	{FieldKeyV3("HV", 3), "victim_case_form2_facts_start_time", KindTime},
	{FieldKeyV3("HV", 4), "", KindText}, // Departamento hechos -- cascada UI
	{FieldKeyV3("HV", 5), "victim_case_form2_facts_town_code", KindText},
	{FieldKeyV3("HV", 6), "victim_case_form2_facts_zone", KindEnum1},
	{FieldKeyV3("HV", 7), "victim_case_form2_facts_address", KindText},
	{FieldKeyV3("HV", 8), "victim_case_form2_scenario_violence", KindEnum1},
	{FieldKeyV3("HV", 9), "", KindEnumN},  // Tipo de violencia
	{FieldKeyV3("HV", 10), "", KindEnumN}, // Subtipo de violencia
	{FieldKeyV3("HV", 11), "", KindEnumN}, // Ámbito de la violencia
	{FieldKeyV3("HV", 12), "victim_case_form2_workplace_sector_occurrence", KindEnum1},
	{FieldKeyV3("HV", 13), "victim_case_form2_violence_motivated_by_gender", KindBoolEnum},
	{FieldKeyV3("HV", 14), "victim_case_form2_recurrence_aggression", KindEnum1},

	// ── AG — Características del agresor ─────────────────────────────────────
	{FieldKeyV3("AG", 1), "victim_case_form2_num_agressors", KindEnum1},
	{FieldKeyV3("AG", 2), "victim_case_form2_proximity_principal_aggressor", KindEnum1}, // "conoce y convive" -- enum reemplazado, D (categorías)
	{FieldKeyV3("AG", 3), "victim_case_form2_relationship_with_presumed_aggressor", KindEnum1},
	{FieldKeyV3("AG", 4), "victim_case_form2_aggressors_occupation", KindEnum1}, // corrige GAP G-12 de V2 -- categoría propia nueva
	{FieldKeyV3("AG", 5), "victim_case_form2_economically_dependent", KindBoolEnum},
	{FieldKeyV3("AG", 6), "victim_case_form2_aggressor_gender_identity", KindEnum1},
	{FieldKeyV3("AG", 7), "victim_case_form2_aggressor_names", KindText},
	{FieldKeyV3("AG", 8), "victim_case_form2_aggressor_doc_type", KindEnum1},
	{FieldKeyV3("AG", 9), "victim_case_form2_aggressor_doc_number", KindText},
	{FieldKeyV3("AG", 10), "victim_case_form2_aggressor_address", KindText},
	{FieldKeyV3("AG", 11), "victim_case_form2_aggressor_phone", KindInt},

	// ── TZ — Tamizaje V3 (columnas nuevas, no las 38 de V2) ──────────────────
	{FieldKeyV3("TZ", 1), "victim_case_form2_tzv3_g1", KindBoolEnum},
	{FieldKeyV3("TZ", 2), "victim_case_form2_tzv3_g2", KindBoolEnum},
	{FieldKeyV3("TZ", 3), "victim_case_form2_tzv3_g3", KindBoolEnum},
	{FieldKeyV3("TZ", 4), "victim_case_form2_tzv3_g4", KindBoolEnum},
	{FieldKeyV3("TZ", 5), "victim_case_form2_tzv3_g5", KindBoolEnum},
	{FieldKeyV3("TZ", 6), "victim_case_form2_tzv3_g6", KindBoolEnum},
	{FieldKeyV3("TZ", 7), "victim_case_form2_tzv3_g7", KindBoolEnum},
	{FieldKeyV3("TZ", 8), "victim_case_form2_tzv3_vpp", KindBoolEnum},

	// ── GS — Gestión Salvia ───────────────────────────────────────────────────
	{FieldKeyV3("GS", 1), "victim_case_form2_allows_easy_report", KindBoolEnum},
	{FieldKeyV3("GS", 2), "", KindEnumN}, // gestiones interinstitucionales
	{FieldKeyV3("GS", 3), "victim_case_form2_interinstitutional_actions_note", KindText},
	{FieldKeyV3("GS", 4), "victim_case_form2_saliva_management_explanation", KindText}, // reutiliza columna V2 (mismo propósito: bitácora inicial)
	{FieldKeyV3("GS", 5), "victim_case_form2_requires_referral", KindBoolEnum},
	{FieldKeyV3("GS", 6), "", KindEnumN},  // ¿Cuáles equipos?
	{FieldKeyV3("GS", 7), "", KindEnumN},  // ¿Cuáles medidas de emergencia?
	{FieldKeyV3("GS", 9), "", KindEnumN},  // Criterios Psicosocial
	{FieldKeyV3("GS", 10), "", KindEnumN}, // Criterio Hombres
	{FieldKeyV3("GS", 11), "", KindEnumN}, // Servicios Discapacidad
	{FieldKeyV3("GS", 12), "", KindEnumN}, // Criterios Estabilización
	{FieldKeyV3("GS", 13), "victim_case_form2_referral_justification", KindText},
	{FieldKeyV3("GS", 14), "victim_case_form2_reported_to_entity", KindBoolEnum},
	{FieldKeyV3("GS", 15), "victim_case_form2_agrees_route_activation", KindBoolEnum},
	{FieldKeyV3("GS", 16), "victim_case_form2_salvia_activates_routes", KindBoolEnum},
	{FieldKeyV3("GS", 17), "victim_case_form2_identifies_barriers", KindBoolEnum},
	{FieldKeyV3("GS", 18), "victim_case_form2_wants_imminent_alert", KindBoolEnum},
	{FieldKeyV3("GS", 19), "victim_case_form2_imminent_alert_note", KindText}, // justificación, solo si Q18=true
}

// victimCaseV3TamizajeBranchFields -- ramas pareja/no-pareja del Tamizaje
// V3 (24 + 24, incluye ideación suicida en el ítem 24 de cada rama). Se
// activa solo una u otra según wasPartner -- ver UpsertFormV3FromAnswers.
var victimCaseV3TamizajePareja = []VictimCaseFormField{} //nolint:gochecknoglobals
var victimCaseV3TamizajeNoPareja = []VictimCaseFormField{} //nolint:gochecknoglobals

func init() {
	for i := 1; i <= 24; i++ {
		victimCaseV3TamizajePareja = append(victimCaseV3TamizajePareja, VictimCaseFormField{
			Key:    FieldKeyV3("TZP", i),
			Column: fmt.Sprintf("victim_case_form2_tzv3_p%d", i),
			Kind:   KindBoolEnum,
		})
		victimCaseV3TamizajeNoPareja = append(victimCaseV3TamizajeNoPareja, VictimCaseFormField{
			Key:    FieldKeyV3("TZNP", i),
			Column: fmt.Sprintf("victim_case_form2_tzv3_np%d", i),
			Kind:   KindBoolEnum,
		})
	}
}

// victimCaseV3EnumNCategories -- campos KindEnumN (sin columna propia, van a
// rel_victim_case_form2_enums_victim_case_form2) junto con las categorías de
// victim_case_form2_enums que agrupan sus opciones. Varias categorías por
// campo porque rel_* no distingue "a qué pregunta pertenece" cada relación
// -- solo se puede acotar por categoría (ver resolveEnumNUpdates). HV_10
// (subtipo de violencia) cubre las 7 categorías de subtipo a la vez porque
// puede tener seleccionados subtipos de varios tipos simultáneamente.
var victimCaseV3EnumNCategories = map[string][]string{ //nolint:gochecknoglobals
	FieldKeyV3("CI", 8):  {"victim_case_form2_adjustments_gbv"},
	FieldKeyV3("DP", 3):  {"victim_case_form2_activities_unable_to_perform"},
	FieldKeyV3("DP", 4):  {"victim_case_form2_law_1996"},
	FieldKeyV3("DP", 11): {"victim_case_form2_specially_protected_population"},
	FieldKeyV3("DP", 20): {"victim_case_form2_asp_mode"},
	FieldKeyV3("DP", 22): {"victim_case_form2_reason_asp"},
	FieldKeyV3("DP", 25): {"victim_case_form2_has_dependents"},
	FieldKeyV3("HV", 9):  {"victim_case_form2_type_violence_experienced"},
	FieldKeyV3("HV", 10): {
		"victim_case_form2_subtype_violence_experienced_fi", "victim_case_form2_subtype_violence_experienced_pl",
		"victim_case_form2_subtype_violence_experienced_po", "victim_case_form2_subtype_violence_experienced_ps",
		"victim_case_form2_subtype_violence_experienced_re", "victim_case_form2_subtype_violence_experienced_se",
		"victim_case_form2_subtype_violence_experienced_vi",
	},
	FieldKeyV3("HV", 11): {"victim_case_form2_scope_of_violence"},
	FieldKeyV3("GS", 2):  {"victim_case_form2_interinst_action"},
	FieldKeyV3("GS", 6):  {"victim_case_form2_referral_team"},
	FieldKeyV3("GS", 7):  {"victim_case_form2_emergency_measure"},
	FieldKeyV3("GS", 9):  {"victim_case_form2_psycho_criteria"},
	FieldKeyV3("GS", 10): {"victim_case_form2_men_criteria"},
	FieldKeyV3("GS", 11): {"victim_case_form2_disability_service"},
	FieldKeyV3("GS", 12): {"victim_case_form2_stabiliz_criteria"},
}

// victimCaseV3EnumCategories -- catálogo completo de categorías que expone
// GET /victim-cases-v3/enums (une las de arriba con las de los campos
// KindEnum1/KindBoolEnum que el frontend necesita resolver a icode, y las
// de Entidades/Barreras que no viven en victim_case_form2). No es necesario
// listar aquí 'yes_no' (los booleanos se resuelven aparte en el frontend).
var victimCaseV3EnumCategories = []string{ //nolint:gochecknoglobals
	"victim_case_form2_pronoun",
	"victim_case_form2_nationality", "victim_case_form2_specified_nationality", "victim_case_form2_migration_condition",
	"victim_case_form2_gender_identity", "victim_case_form2_sexual_orientation", "victim_case_form2_assigned_sex_at_birth",
	"victim_case_form2_ethnic_affiliation", "victim_case_form2_indigenous_people", "victim_case_form2_marital_status",
	"victim_case_form2_last_education_level", "victim_case_form2_occupation", "victim_case_form2_income_generation_method",
	"victim_case_form2_employment_relationship", "victim_case_form2_housing_tenancy_form", "victim_case_form2_housing_stratum",
	"victim_case_form2_residence_zone", "victim_case_form2_facts_zone", "victim_case_form2_scenario_violence",
	"victim_case_form2_workplace_sector_occurrence", "victim_case_form2_recurrence_aggression", "victim_case_form2_num_agressors",
	"victim_case_form2_proximity_principal_aggressor", "victim_case_form2_relationship_with_presumed_aggressor_01",
	"victim_case_form2_relationship_with_presumed_aggressor_02", "victim_case_form2_aggressors_occupation",
	"victim_case_form2_aggressor_gender_identity", "victim_case_form2_victim_doc_type", "victim_case_form2_support_contact_kinship",
	"victim_case_form2_adjustments_gbv", "victim_case_form2_activities_unable_to_perform", "victim_case_form2_law_1996",
	"victim_case_form2_specially_protected_population", "victim_case_form2_asp_mode", "victim_case_form2_reason_asp",
	"victim_case_form2_has_dependents", "victim_case_form2_type_violence_experienced",
	"victim_case_form2_subtype_violence_experienced_fi", "victim_case_form2_subtype_violence_experienced_pl",
	"victim_case_form2_subtype_violence_experienced_po", "victim_case_form2_subtype_violence_experienced_ps",
	"victim_case_form2_subtype_violence_experienced_re", "victim_case_form2_subtype_violence_experienced_se",
	"victim_case_form2_subtype_violence_experienced_vi", "victim_case_form2_scope_of_violence",
	"victim_case_form2_tzv3_risk_type", "victim_case_form2_tzv3_risk_level",
	"victim_case_form2_referral_team", "victim_case_form2_emergency_measure", "victim_case_form2_psycho_criteria",
	"victim_case_form2_men_criteria", "victim_case_form2_disability_service", "victim_case_form2_stabiliz_criteria",
	"victim_case_form2_interinst_action", "victim_case_form2_sector", "victim_case_form2_activation_channel",
	"victim_case_form2_barrier_health", "victim_case_form2_barrier_justice", "victim_case_form2_barrier_protect",
	"victim_case_form2_barrier_institutional", "victim_case_form2_barrier_economic", "victim_case_form2_barrier_territorial",
	"victim_case_form2_barrier_differential", "victim_case_form2_barrier_management", "victim_case_form2_barrier_response",
	"victim_case_form2_barrier_close_reason",
}
