// Package repository — case_info_repository.go
// Repositorio para el componente <case-info>: carga toda la información de un caso.
package repository

import (
	"context"

	"gorm.io/gorm"
)

// CaseInfoRaw contiene los datos crudos de las 3 tablas para un caso.
// Los campos se resuelven con COALESCE para priorizar form2 sobre form1.
type CaseInfoRaw struct {
	// victim_case
	Nombres           string `gorm:"column:nombres"`
	Apellidos         string `gorm:"column:apellidos"`
	DocType           string `gorm:"column:doc_type"`
	DocNumber         string `gorm:"column:doc_number"`
	Status            string `gorm:"column:status"`
	TownCode          string `gorm:"column:town_code"`
	OwnerDescription  string `gorm:"column:owner_description"`
	AgentId           string `gorm:"column:agent_id"`
	AgentName         string `gorm:"column:agent_name"`
	CreationDate      string `gorm:"column:creation_date"`
	UpdateDate        string `gorm:"column:update_date"`

	// form1
	F1Age                    *int   `gorm:"column:f1_age"`
	F1Nationality            string `gorm:"column:f1_nationality"`
	F1NationalityOther       string `gorm:"column:f1_nationality_other"`
	F1Gender                 string `gorm:"column:f1_gender"`
	F1GenderIdentityOther    string `gorm:"column:f1_gender_identity_other"`
	F1SexualOrientationOther string `gorm:"column:f1_sexual_orientation_other"`
	F1Email                  string `gorm:"column:f1_email"`
	F1ForeignerStatus        string `gorm:"column:f1_foreigner_status"`
	F1Dependents             string `gorm:"column:f1_dependents"`
	F1ChildrenNumber         *int   `gorm:"column:f1_children_number"`
	F1ViolenceScene          string `gorm:"column:f1_violence_scene"`
	F1FemicideRisk           string `gorm:"column:f1_femicide_risk"`
	F1DeathThreats           string `gorm:"column:f1_death_threats"`
	F1AggressorHasWeapons    string `gorm:"column:f1_aggressor_has_weapons"`
	F1ViolenceBefore         string `gorm:"column:f1_violence_before"`
	F1PhysicalIncreased      string `gorm:"column:f1_physical_increased"`
	F1SeparatedLastYear      string `gorm:"column:f1_separated_last_year"`
	F1ThreatenedWeapon       string `gorm:"column:f1_threatened_weapon"`
	F1ThreatenedChildren     string `gorm:"column:f1_threatened_children"`
	F1JealousViolent         string `gorm:"column:f1_jealous_violent"`
	F1CapableKilling         string `gorm:"column:f1_capable_killing"`
	F1ImminentRisk           string `gorm:"column:f1_imminent_risk"`
	F1PreviouslyReported     string `gorm:"column:f1_previously_reported"`
	F1IfPreviouslyReported   string `gorm:"column:f1_if_previously_reported"`
	F1Aggressor              string `gorm:"column:f1_aggressor"`
	F1RelationshipAggressor  string `gorm:"column:f1_relationship_aggressor"`
	F1AggressorName          string `gorm:"column:f1_aggressor_name"`
	F1AggressorDocType       string `gorm:"column:f1_aggressor_doc_type"`
	F1AggressorDocNumber     string `gorm:"column:f1_aggressor_doc_number"`
	F1AggressorAddress       string `gorm:"column:f1_aggressor_address"`
	F1AggressorPhone         string `gorm:"column:f1_aggressor_phone"`
	F1IfAfro                 string `gorm:"column:f1_if_afro"`
	F1IfIndigenous           string `gorm:"column:f1_if_indigenous"`
	F1IfIndigenousTongue     string `gorm:"column:f1_if_indigenous_tongue"`
	F1IfPeasant              string `gorm:"column:f1_if_peasant"`
	F1IfArmedConflict        string `gorm:"column:f1_if_armed_conflict"`

	// form1 — hechos
	F1FactsOccurrence        string `gorm:"column:f1_facts_occurrence"`
	F1FactsStartTime         string `gorm:"column:f1_facts_start_time"`
	F1FactsEndTime           string `gorm:"column:f1_facts_end_time"`
	F1FactsWeekday           string `gorm:"column:f1_facts_weekday"`
	F1FactsDate              string `gorm:"column:f1_facts_date"`
	F1FactsDescription       string `gorm:"column:f1_facts_description"`
	F1ViolenceExperienced    string `gorm:"column:f1_violence_experienced"`
	F1ViolenceExperiencedOther string `gorm:"column:f1_violence_experienced_other"`
	F1ViolenceScope          string `gorm:"column:f1_violence_scope"`

	// form2 (prioridad)
	F2RiskLevel        *int   `gorm:"column:f2_risk_level"`
	F2Age              *int64 `gorm:"column:f2_age"`
	F2BirthDate        string `gorm:"column:f2_birth_date"`
	F2Phone            string `gorm:"column:f2_phone"`
	F2ContactPhone     string `gorm:"column:f2_contact_phone"`
	F2GenderIdentity   string `gorm:"column:f2_gender_identity"`
	F2SexualOrientation string `gorm:"column:f2_sexual_orientation"`
	F2ContactNames     string `gorm:"column:f2_contact_names"`
	F2ContactKinship   string `gorm:"column:f2_contact_kinship"`
	F2MaritalStatus    string `gorm:"column:f2_marital_status"`
	F2Disability       string `gorm:"column:f2_disability"`
	F2EthnicAffiliation string `gorm:"column:f2_ethnic_affiliation"`
	F2IndigenousPeople string `gorm:"column:f2_indigenous_people"`
	F2Campesino        string `gorm:"column:f2_campesino"`
	F2FactsDescription string `gorm:"column:f2_facts_description"`
	F2FactsDate        string `gorm:"column:f2_facts_date"`
	F2FactsStartTime   string `gorm:"column:f2_facts_start_time"`
	F2FactsAddress     string `gorm:"column:f2_facts_address"`
	F2ResidenceAddress string `gorm:"column:f2_residence_address"`
	F2IdentityName     string `gorm:"column:f2_identity_name"`
	F2RequireInterpreter string `gorm:"column:f2_require_interpreter"`
	F2NumAggressors    string `gorm:"column:f2_num_aggressors"`
	F2ProximityAggressor string `gorm:"column:f2_proximity_aggressor"`
	F2AggressorGender  string `gorm:"column:f2_aggressor_gender"`
	F2AggressorDocType string `gorm:"column:f2_aggressor_doc_type"`
	F2ScenarioViolence string `gorm:"column:f2_scenario_violence"`
	F2AggressorNames   string `gorm:"column:f2_aggressor_names"`
	F2AggressorDocNumber string `gorm:"column:f2_aggressor_doc_number"`
	F2AggressorAddress string `gorm:"column:f2_aggressor_address"`
	F2AggressorPhone   string `gorm:"column:f2_aggressor_phone"`
	F2RelationshipAggressor string `gorm:"column:f2_relationship_aggressor"`
	F2ManagementExplanation string `gorm:"column:f2_management_explanation"`

	// Ubicación resuelta
	CityName       string `gorm:"column:city_name"`
	DeptName       string `gorm:"column:dept_name"`
	TownName       string `gorm:"column:town_name"`
}

// CaseInfoRepository define el acceso a datos para la información completa del caso.
type CaseInfoRepository interface {
	GetFullInfoByICode(ctx context.Context, caseICode string) (*CaseInfoRaw, error)
	GetPlanAtencionByICode(ctx context.Context, caseICode string) ([]string, error)
}

type caseInfoRepository struct {
	db *gorm.DB
}

func NewCaseInfoRepository(db *gorm.DB) CaseInfoRepository {
	return &caseInfoRepository{db: db}
}

func (r *caseInfoRepository) GetFullInfoByICode(ctx context.Context, caseICode string) (*CaseInfoRaw, error) {
	var raw CaseInfoRaw
	sql := `
		SELECT
			-- victim_case
			COALESCE(vc.victim_case_victim_names, '')          AS nombres,
			COALESCE(vc.victim_case_victim_last_names, '')     AS apellidos,
			COALESCE(vc.victim_case_victim_doc_type, '')       AS doc_type,
			COALESCE(vc.victim_case_victim_doc_number, '')     AS doc_number,
			COALESCE(vc.victim_case_status, '')                AS status,
			COALESCE(vc.victim_case_victim_town_code, '')      AS town_code,
			COALESCE(vc.victim_case_owner_description, '')     AS owner_description,
			COALESCE(vc.agent_id, '')                          AS agent_id,
			COALESCE((SELECT gup2.general_user_profile_names || ' ' || gup2.general_user_profile_last_names FROM security.general_user gu2 JOIN security.general_user_profile gup2 ON gup2.general_user_profile_id = gu2.general_user_general_user_profile WHERE gu2.general_user_i_code = vc.agent_id LIMIT 1), '') AS agent_name,
			TO_CHAR(vc.victim_case_creation_date, 'YYYY-MM-DD HH24:MI:SS') AS creation_date,
			TO_CHAR(vc.victim_case_update_date, 'YYYY-MM-DD HH24:MI:SS')   AS update_date,

			-- form1
			f1.victim_case_form1_age                                          AS f1_age,
			COALESCE(f1.victim_case_form1_victim_nationality, '')             AS f1_nationality,
			COALESCE(f1.victim_case_form1_victim_nationality_other, '')       AS f1_nationality_other,
			COALESCE(f1.victim_case_form1_victim_gender, '')                  AS f1_gender,
			COALESCE(f1.victim_case_form1_victim_gender_identity_other, '')   AS f1_gender_identity_other,
			COALESCE(f1.victim_case_form1_victim_sexual_orientation_other,'') AS f1_sexual_orientation_other,
			COALESCE(f1.victim_case_form1_victim_e_mail, '')                  AS f1_email,
			COALESCE(f1.victim_case_form1_victim_foreigner_immigration_status,'') AS f1_foreigner_status,
			COALESCE(f1.victim_case_form1_victim_dependents, '')              AS f1_dependents,
			f1.victim_case_form1_victim_children_number                       AS f1_children_number,
			COALESCE(f1.victim_case_form1_victim_violence_scene, '')          AS f1_violence_scene,
			COALESCE(f1.victim_case_form1_victim_femicide_risk, '')           AS f1_femicide_risk,
			COALESCE(f1.victim_case_form1_victim_death_threats, '')           AS f1_death_threats,
			COALESCE(f1.victim_case_form1_victim_aggressor_has_weapons, '')   AS f1_aggressor_has_weapons,
			COALESCE(f1.victim_case_form1_experienced_physical_or_sexual_violence_befor,'') AS f1_violence_before,
			COALESCE(f1.victim_case_form1_physical_violence_increased, '')    AS f1_physical_increased,
			COALESCE(f1.victim_case_form1_separated_from_partner_last_year,'') AS f1_separated_last_year,
			COALESCE(f1.victim_case_form1_threatened_with_weapon, '')         AS f1_threatened_weapon,
			COALESCE(f1.victim_case_form1_threatened_to_kill_or_harm_children,'') AS f1_threatened_children,
			COALESCE(f1.victim_case_form1_jealous_and_violent, '')            AS f1_jealous_violent,
			COALESCE(f1.victim_case_form1_believes_capable_of_killing, '')    AS f1_capable_killing,
			COALESCE(f1.victim_case_form1_victim_imminent_risk, '')           AS f1_imminent_risk,
			COALESCE(f1.victim_case_form1_victim_previously_reported_situation,'') AS f1_previously_reported,
			COALESCE(f1.victim_case_form1_victim_if_previously_reported, '')  AS f1_if_previously_reported,
			COALESCE(f1.victim_case_form1_victim_aggressor, '')               AS f1_aggressor,
			COALESCE(f1.victim_case_form1_victim_relationship_with_aggressor,'') AS f1_relationship_aggressor,
			COALESCE(f1.victim_case_form1_victim_aggressor_name, '')          AS f1_aggressor_name,
			COALESCE(f1.victim_case_form1_victim_aggressor_doc_type, '')      AS f1_aggressor_doc_type,
			COALESCE(f1.victim_case_form1_victim_aggressor_doc_number, '')    AS f1_aggressor_doc_number,
			COALESCE(f1.victim_case_form1_victim_aggressor_address, '')       AS f1_aggressor_address,
			COALESCE(f1.victim_case_form1_victim_aggressor_phone, '')         AS f1_aggressor_phone,
			COALESCE(f1.victim_case_form1_victim_if_afro, '')                 AS f1_if_afro,
			COALESCE(f1.victim_case_form1_victim_if_indigenous, '')           AS f1_if_indigenous,
			COALESCE(f1.victim_case_form1_victim_if_indigenous_tongue, '')    AS f1_if_indigenous_tongue,
			COALESCE(f1.victim_case_form1_victim_if_peasant, '')              AS f1_if_peasant,
			COALESCE(f1.victim_case_form1_victim_if_armed_conflict, '')       AS f1_if_armed_conflict,

			-- form1 hechos
			COALESCE(f1.victim_case_form1_facts_occurrence, '')               AS f1_facts_occurrence,
			COALESCE(f1.victim_case_form1_facts_start_time::text, '')         AS f1_facts_start_time,
			COALESCE(f1.victim_case_form1_facts_end_time::text, '')           AS f1_facts_end_time,
			COALESCE(f1.victim_case_form1_facts_weekday::text, '')            AS f1_facts_weekday,
			COALESCE(TO_CHAR(f1.victim_case_form1_facts_date, 'DD/MM/YYYY'), '') AS f1_facts_date,
			COALESCE(f1.victim_case_form1_facts_description, '')              AS f1_facts_description,
			COALESCE(f1.victim_case_form1_victim_violence_experienced, '')    AS f1_violence_experienced,
			COALESCE(f1.victim_case_form1_victim_violence_experienced_other, '') AS f1_violence_experienced_other,
			COALESCE(f1.victim_case_form1_victim_violence_scope, '')          AS f1_violence_scope,

			-- form2 (prioridad)
			f2.victim_case_form2_risk_level                                   AS f2_risk_level,
			EXTRACT(YEAR FROM AGE(NOW(), f2.victim_case_form2_birth_date))::int AS f2_age,
			COALESCE(TO_CHAR(f2.victim_case_form2_birth_date, 'YYYY-MM-DD'), '') AS f2_birth_date,
			COALESCE(f2.victim_case_form2_victim_phone::text, '')             AS f2_phone,
			COALESCE(f2.victim_case_form2_support_contact_phone::text, '')    AS f2_contact_phone,
			COALESCE(gi.victim_case_form2_enums_name, '')                     AS f2_gender_identity,
			COALESCE(so.victim_case_form2_enums_name, '')                     AS f2_sexual_orientation,
			COALESCE(f2.victim_case_form2_support_contact_names, '')          AS f2_contact_names,
			COALESCE(kin.victim_case_form2_enums_name, '')                    AS f2_contact_kinship,
			COALESCE(ms.victim_case_form2_enums_name, '')                     AS f2_marital_status,
			COALESCE(dis.victim_case_form2_enums_name, '')                    AS f2_disability,
			COALESCE(eth.victim_case_form2_enums_name, '')                    AS f2_ethnic_affiliation,
			COALESCE(ind.victim_case_form2_enums_name, '')                    AS f2_indigenous_people,
			COALESCE(cam.victim_case_form2_enums_name, '')                    AS f2_campesino,
			COALESCE(f2.victim_case_form2_facts_description, '')              AS f2_facts_description,
			COALESCE(TO_CHAR(f2.victim_case_form2_facts_date, 'YYYY-MM-DD'), '') AS f2_facts_date,
			COALESCE(f2.victim_case_form2_facts_start_time::text, '')               AS f2_facts_start_time,
			COALESCE(f2.victim_case_form2_facts_address, '')                      AS f2_facts_address,
			COALESCE(f2.victim_case_form2_residence_address, '')                   AS f2_residence_address,
			COALESCE(f2.victim_case_form2_identity_name, '')                       AS f2_identity_name,
			COALESCE(ri.victim_case_form2_enums_name, '')                          AS f2_require_interpreter,
			COALESCE(nag.victim_case_form2_enums_name, '')                         AS f2_num_aggressors,
			COALESCE(prox.victim_case_form2_enums_name, '')                        AS f2_proximity_aggressor,
			COALESCE(agi.victim_case_form2_enums_name, '')                         AS f2_aggressor_gender,
			COALESCE(adt.victim_case_form2_enums_name, '')                         AS f2_aggressor_doc_type,
			COALESCE(sv.victim_case_form2_enums_name, '')                     AS f2_scenario_violence,
			COALESCE(f2.victim_case_form2_aggressor_names, '')                AS f2_aggressor_names,
			COALESCE(f2.victim_case_form2_aggressor_doc_number, '')           AS f2_aggressor_doc_number,
			COALESCE(f2.victim_case_form2_aggressor_address, '')              AS f2_aggressor_address,
			COALESCE(f2.victim_case_form2_aggressor_phone::text, '')          AS f2_aggressor_phone,
			COALESCE(rel.victim_case_form2_enums_name, '')                    AS f2_relationship_aggressor,
			COALESCE(f2.victim_case_form2_saliva_management_explanation, '')   AS f2_management_explanation,

			-- Ubicación
			COALESCE(c.city_name, '')                                         AS city_name,
			COALESCE(d.department_name, '')                                   AS dept_name,
			COALESCE(t.town_name, '')                                         AS town_name

		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form1 f1
			ON f1.victim_case_form1_victim_case = vc.victim_case_id
		LEFT JOIN salvia.victim_case_form2 f2
			ON f2.victim_case_form2_victim_case = vc.victim_case_id
		LEFT JOIN salvia.victim_case_form2_enums gi
			ON gi.victim_case_form2_enums_id = f2.victim_case_form2_gender_identity
		LEFT JOIN salvia.victim_case_form2_enums so
			ON so.victim_case_form2_enums_id = f2.victim_case_form2_sexual_orientation
		LEFT JOIN salvia.victim_case_form2_enums kin
			ON kin.victim_case_form2_enums_id = f2.victim_case_form2_support_contact_kinship
		LEFT JOIN salvia.victim_case_form2_enums ms
			ON ms.victim_case_form2_enums_id = f2.victim_case_form2_marital_status
		LEFT JOIN salvia.victim_case_form2_enums dis
			ON dis.victim_case_form2_enums_id = f2.victim_case_form2_person_with_disability
		LEFT JOIN salvia.victim_case_form2_enums eth
			ON eth.victim_case_form2_enums_id = f2.victim_case_form2_ethnic_affiliation
		LEFT JOIN salvia.victim_case_form2_enums ind
			ON ind.victim_case_form2_enums_id = f2.victim_case_form2_indigenous_people
		LEFT JOIN salvia.victim_case_form2_enums cam
			ON cam.victim_case_form2_enums_id = f2.victim_case_form2_campesino_recognition
		LEFT JOIN salvia.victim_case_form2_enums sv
			ON sv.victim_case_form2_enums_id = f2.victim_case_form2_scenario_violence
		LEFT JOIN salvia.victim_case_form2_enums rel
			ON rel.victim_case_form2_enums_id = f2.victim_case_form2_relationship_with_presumed_aggressor
		LEFT JOIN salvia.victim_case_form2_enums ri
			ON ri.victim_case_form2_enums_id = f2.victim_case_form2_require_language_interpreter
		LEFT JOIN salvia.victim_case_form2_enums nag
			ON nag.victim_case_form2_enums_id = f2.victim_case_form2_num_agressors
		LEFT JOIN salvia.victim_case_form2_enums prox
			ON prox.victim_case_form2_enums_id = f2.victim_case_form2_proximity_principal_aggressor
		LEFT JOIN salvia.victim_case_form2_enums agi
			ON agi.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_gender_identity
		LEFT JOIN salvia.victim_case_form2_enums adt
			ON adt.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_doc_type
		LEFT JOIN security.town t
			ON t.town_code = vc.victim_case_victim_town_code
		LEFT JOIN security.city c
			ON c.city_id = t.city_id
		LEFT JOIN security.department d
			ON d.department_id = c.department_id
		WHERE vc.victim_case_i_code = ?
		LIMIT 1
	`
	err := r.db.WithContext(ctx).Raw(sql, caseICode).Scan(&raw).Error
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func (r *caseInfoRepository) GetPlanAtencionByICode(ctx context.Context, caseICode string) ([]string, error) {
	type enumRow struct {
		Name string `gorm:"column:enum_name"`
	}
	var enums []enumRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT e.victim_case_form2_enums_name AS enum_name
		FROM salvia.rel_victim_case_form2_enums_victim_case_form2 rel
		JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = rel.victim_case_form2_enums_id
		WHERE rel.victim_case_form2_id = (
			SELECT victim_case_form2_id FROM salvia.victim_case_form2
			WHERE victim_case_form2_victim_case = (
				SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ? LIMIT 1
			) LIMIT 1
		)
		AND e.victim_case_form2_enums_category LIKE '%action_plan%'
	`, caseICode).Scan(&enums).Error
	if err != nil {
		return nil, err
	}
	var result []string
	for _, e := range enums {
		result = append(result, e.Name)
	}
	return result, nil
}
