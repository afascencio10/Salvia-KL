package repository

import (
	"bitsflow/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type CaseReportDTO struct {
	VictimCaseId              int64     `gorm:"column:victim_case_id"`
	VictimCaseICode           string    `gorm:"column:victim_case_i_code"`
	VictimCaseStatus          string    `gorm:"column:victim_case_status"`
	VictimCaseVictimName      string    `gorm:"column:victim_case_victim_names"`
	VictimCaseVictimLastName  string    `gorm:"column:victim_case_victim_last_names"`
	VictimCaseVictimAge       *int      `gorm:"column:victim_case_victim_age"`
	VictimCaseVictimGender    string    `gorm:"column:victim_case_victim_gender"`
	VictimCaseVictimDocNumber string    `gorm:"column:victim_case_victim_doc_number"`
	VictimCaseTownCode        string    `gorm:"column:victim_case_victim_town_code"`
	VictimCaseTeam            string    `gorm:"column:victim_case_team"`
	VictimCaseCreatedAt       time.Time `gorm:"column:victim_case_creation_date"`
	VictimPhone               string    `gorm:"column:victim_phone"`
	HasForm1                  string    `gorm:"column:has_form1"`
	HasForm2                  string    `gorm:"column:has_form2"`
	TownName                  string    `gorm:"column:town_name"`
	DeptName                  string    `gorm:"column:dept_name"`

	// Campos enriquecidos de la víctima (Hoja Registro)
	VictimDocType          string `gorm:"column:victim_doc_type"`
	VictimAddress          string `gorm:"column:victim_address"`
	VictimEmail            string `gorm:"column:victim_email"`
	VictimIdentityName     string `gorm:"column:victim_identity_name"`
	VictimSexualOrientation string `gorm:"column:victim_sexual_orientation"`
	VictimOrigin           string `gorm:"column:victim_origin"`
	VictimOccupation       string `gorm:"column:victim_occupation"`
	VictimEthnicGroup      string `gorm:"column:victim_ethnic_group"`
	VictimIfAfro           string `gorm:"column:victim_if_afro"`
	VictimIfIndigenous     string `gorm:"column:victim_if_indigenous"`
	VictimIfPeasant        string `gorm:"column:victim_if_peasant"`
	VictimChildrenNumber   int    `gorm:"column:victim_children_number"`
	VictimMaritalStatus    string `gorm:"column:victim_marital_status"`
	VictimDisability       string `gorm:"column:victim_disability"`

	// Contacto de emergencia
	VictimContactNames   string `gorm:"column:victim_contact_names"`
	VictimContactPhone   string `gorm:"column:victim_contact_phone"`
	VictimContactKinship string `gorm:"column:victim_contact_kinship"`

	// Hechos
	FactsDescription string `gorm:"column:facts_description"`
	FactsDate        string `gorm:"column:facts_date"`
	FactsAddress     string `gorm:"column:facts_address"`
	FactsStartTime   string `gorm:"column:facts_start_time"`
	FactsEndTime     string `gorm:"column:facts_end_time"`
	FactsWeekday     string `gorm:"column:facts_weekday"`
	FactsOccurrence  string `gorm:"column:facts_occurrence"`

	// Violencias y Riesgos
	ViolenceScope string `gorm:"column:violence_scope"`
	ViolenceScene string `gorm:"column:violence_scene"`
	FemicideRisk  string `gorm:"column:femicide_risk"`
	RiskLevel     string `gorm:"column:risk_level"`

	// Agresor
	AggressorName         string `gorm:"column:aggressor_name"`
	AggressorDocType      string `gorm:"column:aggressor_doc_type"`
	AggressorDocNumber    string `gorm:"column:aggressor_doc_number"`
	AggressorAddress      string `gorm:"column:aggressor_address"`
	AggressorPhone        string `gorm:"column:aggressor_phone"`
	RelationshipAggressor string `gorm:"column:relationship_aggressor"`

	// Preguntas específicas valoración riesgo
	PhysicalViolenceIncreased      string `gorm:"column:physical_violence_increased"`
	SeparatedFromPartnerLastYear   string `gorm:"column:separated_last_year"`
	ThreatenedWithWeapon           string `gorm:"column:threatened_weapon"`
	ThreatenedToKillOrHarmChildren string `gorm:"column:threatened_children"`
	JealousAndViolent              string `gorm:"column:jealous_violent"`
	BelievesCapableOfKilling       string `gorm:"column:capable_killing"`

	// Nuevas preguntas de riesgo pareja / no pareja (Form 2)
	VictimHealthToBlackmail                   string `gorm:"column:victim_health_to_blackmail"`
	ThreatenedRevealSexualOrientation         string `gorm:"column:threatened_reveal_sexual_orientation"`
	StoppedSeekingHelp                        string `gorm:"column:stopped_seeking_help"`
	AggressorSexuallyHarassment2              string `gorm:"column:aggressor_sexually_harassment_2"`
	AggressorTakenAdvantagePhysicalVulnerabil string `gorm:"column:aggressor_taken_advantage_physical_vulnerabil"`
	ViolenceMotivatedByGender2                string `gorm:"column:violence_motivated_by_gender_2"`

	// Bloque 1: Autorización y Contacto
	AuthorizationAnswer string `gorm:"column:authorization_answer"`
	AdjustmentsGBV      string `gorm:"column:adjustments_gbv"`
	RequireInterpreter  string `gorm:"column:require_interpreter"`
	SupportContactEmail string `gorm:"column:support_contact_email"`

	// Bloque 2: Hechos Victimizantes (ampliación)
	FactsDeptName        string `gorm:"column:facts_dept_name"`
	FactsCityName        string `gorm:"column:facts_city_name"`
	FactsZone            string `gorm:"column:facts_zone"`
	FactsAddressForm2    string `gorm:"column:facts_address_form2"`
	ViolenceType         string `gorm:"column:violence_type"`
	ViolenceSubtype      string `gorm:"column:violence_subtype"`
	ViolenceScopeForm2   string `gorm:"column:violence_scope_form2"`
	WorkplaceSector      string `gorm:"column:workplace_sector"`
	RecurrenceAggression string `gorm:"column:recurrence_aggression"`

	// Bloque 3: Características del Agresor (ampliación)
	NumAggressors          string `gorm:"column:num_aggressors"`
	ProximityAggressor     string `gorm:"column:proximity_aggressor"`
	EconomicallyDependent  string `gorm:"column:economically_dependent"`
	AggressorGenderIdentity string `gorm:"column:aggressor_gender_identity"`

	// Bloque 4: Tamizaje (preguntas adicionales)
	AggressorPursuesSpies      string `gorm:"column:aggressor_pursues_spies"`
	AggressorHasAccessWeapons  string `gorm:"column:aggressor_has_access_weapons"`
	PartnerUnemployed          string `gorm:"column:partner_unemployed"`
	PartnerOtherDenunciations  string `gorm:"column:partner_other_denunciations"`
	AggressorPenalBackground   string `gorm:"column:aggressor_penal_background"`
	AggressorStrangulation     string `gorm:"column:aggressor_strangulation"`
	AggressorConsumesDrugs     string `gorm:"column:aggressor_consumes_drugs"`
	AggressorIsAlcoholic       string `gorm:"column:aggressor_is_alcoholic"`
	PartnerControls            string `gorm:"column:partner_controls"`
	PartnerThreatenedSuicide   string `gorm:"column:partner_threatened_suicide"`
	PartnerThreatenedDamage    string `gorm:"column:partner_threatened_damage"`
	ThoughtsOfSelfHarm         string `gorm:"column:thoughts_of_self_harm"`
	AggressorLimitsContact     string `gorm:"column:aggressor_limits_contact"`
	StillLivesWithAggressor    string `gorm:"column:still_lives_with_aggressor"`

	// Bloque 5: Datos Personales Víctima (ampliación)
	BirthDate                string `gorm:"column:birth_date_form2"`
	PhysicalDifficulties     string `gorm:"column:physical_difficulties"`
	Nationality              string `gorm:"column:nationality_form2"`
	SpecifiedNationality     string `gorm:"column:specified_nationality"`
	MigrationCondition       string `gorm:"column:migration_condition"`
	GenderIdentity           string `gorm:"column:gender_identity_form2"`
	AssignedSexAtBirth       string `gorm:"column:assigned_sex_at_birth"`
	SpecialProtectedPop      string `gorm:"column:special_protected_pop"`
	LastEducationLevel       string `gorm:"column:last_education_level"`
	IncomeGenerationMethod   string `gorm:"column:income_generation_method"`
	EmploymentRelationship   string `gorm:"column:employment_relationship"`
	ASPMode                  string `gorm:"column:asp_mode"`
	ApproxStartASP           string `gorm:"column:approx_start_asp"`
	ReasonASP                string `gorm:"column:reason_asp"`
	HousingTenancy           string `gorm:"column:housing_tenancy"`
	HousingStratum           string `gorm:"column:housing_stratum"`
	HasDependents            string `gorm:"column:has_dependents"`
	CurrentlyPregnant        string `gorm:"column:currently_pregnant"`
	ResidenceZone            string `gorm:"column:residence_zone"`

	// Bloque 6: Plan de Acción, Denuncia, Operación
	ActionPlan              string `gorm:"column:action_plan"`
	ManagementExplanation   string `gorm:"column:management_explanation"`
	AllowsEasyReport        string `gorm:"column:allows_easy_report"`
	AgentResponsible        string `gorm:"column:agent_responsible"`
	OwnerDescription        string `gorm:"column:owner_description"`
}

type FollowUpReportDTO struct {
	models.FollowUpV2
	VictimDocNumber string `gorm:"column:victim_doc_number"`
}

type AnswerDTO struct {
	CaseICode        string     `gorm:"column:case_icode"`
	FollowUpID       string     `gorm:"column:follow_up_id"`
	FormSubmissionID string     `gorm:"column:form_submission_id"`
	QuestionDesc     string     `gorm:"column:question_desc"`
	AnswerValue      string     `gorm:"column:answer_value"`
	VictimDocNumber  string     `gorm:"column:victim_doc_number"`
	FollowUpDate     *time.Time `gorm:"column:follow_up_date"`
}

type TimelineReportDTO struct {
	models.CaseTimelineEvent
	VictimDocNumber string `gorm:"column:victim_doc_number"`
}

// ContactReportDTO representa una fila del reporte consolidado de contactos (reportes).
// Un contacto tiene exactamente un Form1 (reporte propio) o un Form2 (reporte de tercero);
// las columnas del formulario que no aplica llegan vacías desde el SQL.
type ContactReportDTO struct {
	VictimContactId   int64     `gorm:"column:victim_contact_id"`
	VictimContactICode string   `gorm:"column:victim_contact_i_code"`
	CreationDate      time.Time `gorm:"column:victim_contact_creation_date"`
	UpdateDate        time.Time `gorm:"column:victim_contact_update_date"`
	StatusDescription string    `gorm:"column:victim_contact_status_description"`
	Names             string    `gorm:"column:victim_contact_names"`
	LastNames         string    `gorm:"column:victim_contact_last_names"`
	Latitude          float64   `gorm:"column:victim_contact_latitude"`
	Longitude         float64   `gorm:"column:victim_contact_longitude"`
	FormType          string    `gorm:"column:form_type"`

	Form1Nick              string `gorm:"column:f1_nick"`
	Form1DocType           string `gorm:"column:f1_doc_type"`
	Form1DocNumber         string `gorm:"column:f1_doc_number"`
	Form1BirthDate         string `gorm:"column:f1_birth_date"`
	Form1TownCode          string `gorm:"column:f1_town_code"`
	Form1TownName          string `gorm:"column:f1_town_name"`
	Form1Address           string `gorm:"column:f1_address"`
	Form1Phone             string `gorm:"column:f1_phone"`
	Form1GenderIdentity    string `gorm:"column:f1_gender_identity"`
	Form1SexualOrientation string `gorm:"column:f1_sexual_orientation"`
	Form1Origin            string `gorm:"column:f1_origin"`
	Form1Occupation        string `gorm:"column:f1_occupation"`
	Form1OccupationOther   string `gorm:"column:f1_occupation_other"`
	Form1FactsDescription  string `gorm:"column:f1_facts_description"`

	Form2ReporterNames       string `gorm:"column:f2_reporter_names"`
	Form2ReporterPhone       string `gorm:"column:f2_reporter_phone"`
	Form2VictimColPhone      string `gorm:"column:f2_victim_col_phone"`
	Form2FactsDescription    string `gorm:"column:f2_facts_description"`
	Form2BestContactTime     string `gorm:"column:f2_best_contact_time"`
	Form2WillReceiveCall     string `gorm:"column:f2_will_receive_call"`
	Form2HasCareRole         string `gorm:"column:f2_has_care_role"`
	Form2VictimAware         string `gorm:"column:f2_victim_aware"`
	Form2ReportType     string `gorm:"column:f2_report_type"`
	Form2AdjustmentsGBV string `gorm:"column:f2_adjustments_gbv"`
}

type ReportRepository interface {
	FindCasesInDateRange(ctx context.Context, start, end time.Time) ([]CaseReportDTO, error)
	FindFollowUpsByCaseICodes(ctx context.Context, caseICodes []string) ([]FollowUpReportDTO, error)
	FindAnswersByFormSubmissions(ctx context.Context, submissionIDs []string) ([]AnswerDTO, error)
	FindTimelineEventsByCaseICodes(ctx context.Context, caseICodes []string) ([]TimelineReportDTO, error)
	FindContactsInDateRange(ctx context.Context, start, end time.Time) ([]ContactReportDTO, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) FindCasesInDateRange(ctx context.Context, start, end time.Time) ([]CaseReportDTO, error) {
	var cases []CaseReportDTO
	query := `
		SELECT 
			vc.victim_case_id,
			vc.victim_case_i_code,
			vc.victim_case_victim_names,
			vc.victim_case_victim_last_names,
			vc.victim_case_victim_doc_number,
			COALESCE(vf1.victim_case_form1_age, EXTRACT(YEAR FROM AGE(NOW(), vf2.victim_case_form2_birth_date))::int) AS victim_case_victim_age,
			COALESCE(NULLIF(vf1.victim_case_form1_victim_gender, ''), gi.victim_case_form2_enums_name, '') AS victim_case_victim_gender,
			vc.victim_case_victim_town_code,
			vc.victim_case_team,
			vc.victim_case_creation_date,
			vc.victim_case_status,
			(
				COALESCE(
					NULLIF(vf2.victim_case_form2_victim_phone::text, ''),
					NULLIF(vf1.victim_case_form1_victim_phone::text, ''),
					NULLIF(vcf1.victim_contact_form1_phone::text, ''),
					''
				)
			) AS victim_phone,
			CASE WHEN vf1.victim_case_form1_victim_case IS NOT NULL THEN 'Sí' ELSE 'No' END AS has_form1,
			CASE WHEN vf2.victim_case_form2_victim_case IS NOT NULL THEN 'Sí' ELSE 'No' END AS has_form2,
			COALESCE(c.city_name, '') AS town_name,
			COALESCE(d.department_name, '') AS dept_name,

			-- Datos enriquecidos de la víctima
			COALESCE(NULLIF(vc.victim_case_victim_doc_type, ''), '') AS victim_doc_type,
			COALESCE(NULLIF(vf2.victim_case_form2_residence_address, ''), NULLIF(vcf1.victim_contact_form1_address, ''), '') AS victim_address,
			COALESCE(NULLIF(vf1.victim_case_form1_victim_e_mail, ''), '') AS victim_email,
			COALESCE(NULLIF(vf2.victim_case_form2_identity_name, ''), NULLIF(vcf1.victim_contact_form1_nick, ''), '') AS victim_identity_name,
			COALESCE(NULLIF(so.victim_case_form2_enums_name, ''), NULLIF(vcf1.victim_contact_form1_sexual_orientation, ''), '') AS victim_sexual_orientation,
			COALESCE(NULLIF(vcf1.victim_contact_form1_origin, ''), '') AS victim_origin,
			COALESCE(NULLIF(oc2.victim_case_form2_enums_name, ''), NULLIF(vf1.victim_case_form1_victim_occupation_other, ''), NULLIF(vf1.victim_case_form1_victim_occupation, ''), '') AS victim_occupation,
			COALESCE(NULLIF(eth.victim_case_form2_enums_name, ''), '') AS victim_ethnic_group,
			COALESCE(vf1.victim_case_form1_victim_if_afro, '') AS victim_if_afro,
			COALESCE(vf1.victim_case_form1_victim_if_indigenous, '') AS victim_if_indigenous,
			COALESCE(vf1.victim_case_form1_victim_if_peasant, '') AS victim_if_peasant,
			COALESCE(vf1.victim_case_form1_victim_children_number, 0) AS victim_children_number,
			COALESCE(NULLIF(ms.victim_case_form2_enums_name, ''), '') AS victim_marital_status,
			COALESCE(NULLIF(dis.victim_case_form2_enums_name, ''), '') AS victim_disability,

			-- Contacto de emergencia
			COALESCE(NULLIF(vf2.victim_case_form2_support_contact_names, ''), NULLIF(vf1.victim_case_form1_victim_contact_names, ''), '') AS victim_contact_names,
			COALESCE(NULLIF(vf2.victim_case_form2_support_contact_phone::text, ''), NULLIF(vf1.victim_case_form1_victim_contact_phone::text, ''), '') AS victim_contact_phone,
			COALESCE(NULLIF(kin.victim_case_form2_enums_name, ''), NULLIF(vf1.victim_case_form1_victim_contact_kinship, ''), '') AS victim_contact_kinship,

			-- Hechos
			COALESCE(NULLIF(vf2.victim_case_form2_facts_description, ''), NULLIF(vf1.victim_case_form1_facts_description, ''), '') AS facts_description,
			COALESCE(NULLIF(TO_CHAR(vf2.victim_case_form2_facts_date, 'YYYY-MM-DD'), ''), NULLIF(TO_CHAR(vf1.victim_case_form1_facts_date, 'YYYY-MM-DD'), ''), '') AS facts_date,
			COALESCE(NULLIF(vf2.victim_case_form2_facts_address, ''), '') AS facts_address,
			COALESCE(NULLIF(vf2.victim_case_form2_facts_start_time::text, ''), NULLIF(vf1.victim_case_form1_facts_start_time::text, ''), '') AS facts_start_time,
			COALESCE(NULLIF(vf1.victim_case_form1_facts_end_time::text, ''), '') AS facts_end_time,
			COALESCE(NULLIF(vf1.victim_case_form1_facts_weekday::text, ''), '') AS facts_weekday,
			COALESCE(NULLIF(vf1.victim_case_form1_facts_occurrence, ''), '') AS facts_occurrence,

			-- Violencias y Riesgos
			COALESCE(NULLIF(vf1.victim_case_form1_victim_violence_scope, ''), '') AS violence_scope,
			COALESCE(NULLIF(sv.victim_case_form2_enums_name, ''), NULLIF(vf1.victim_case_form1_victim_violence_scene, ''), '') AS violence_scene,
			COALESCE(NULLIF(vf1.victim_case_form1_victim_femicide_risk, ''), '') AS femicide_risk,
			CASE vf2.victim_case_form2_risk_level
				WHEN 1 THEN 'Bajo'
				WHEN 2 THEN 'Moderado'
				WHEN 3 THEN 'Alto'
				WHEN 4 THEN 'Extremo'
				ELSE 'Desconocido'
			END AS risk_level,

			-- Agresor
			COALESCE(NULLIF(vf2.victim_case_form2_aggressor_names, ''), NULLIF(vf1.victim_case_form1_victim_aggressor_name, ''), '') AS aggressor_name,
			COALESCE(NULLIF(vf1.victim_case_form1_victim_aggressor_doc_type, ''), '') AS aggressor_doc_type,
			COALESCE(NULLIF(vf2.victim_case_form2_aggressor_doc_number, ''), NULLIF(vf1.victim_case_form1_victim_aggressor_doc_number, ''), '') AS aggressor_doc_number,
			COALESCE(NULLIF(vf2.victim_case_form2_aggressor_address, ''), NULLIF(vf1.victim_case_form1_victim_aggressor_address, ''), '') AS aggressor_address,
			COALESCE(NULLIF(vf2.victim_case_form2_aggressor_phone::text, ''), NULLIF(vf1.victim_case_form1_victim_aggressor_phone, ''), '') AS aggressor_phone,
			COALESCE(NULLIF(rel.victim_case_form2_enums_name, ''), NULLIF(vf1.victim_case_form1_victim_relationship_with_aggressor, ''), '') AS relationship_aggressor,

			-- Preguntas específicas valoración riesgo
			COALESCE(NULLIF(vf1.victim_case_form1_physical_violence_increased, ''), '') AS physical_violence_increased,
			COALESCE(NULLIF(vf1.victim_case_form1_separated_from_partner_last_year, ''), '') AS separated_last_year,
			COALESCE(NULLIF(vf1.victim_case_form1_threatened_with_weapon, ''), '') AS threatened_weapon,
			COALESCE(NULLIF(vf1.victim_case_form1_threatened_to_kill_or_harm_children, ''), '') AS threatened_children,
			COALESCE(NULLIF(vf1.victim_case_form1_jealous_and_violent, ''), '') AS jealous_violent,
			COALESCE(NULLIF(vf1.victim_case_form1_believes_capable_of_killing, '') , '') AS capable_killing,

			-- Nuevas preguntas de riesgo pareja / no pareja (Form 2)
			COALESCE(NULLIF(htb.victim_case_form2_enums_name, ''), '') AS victim_health_to_blackmail,
			COALESCE(NULLIF(tro.victim_case_form2_enums_name, ''), '') AS threatened_reveal_sexual_orientation,
			COALESCE(NULLIF(ssh.victim_case_form2_enums_name, ''), '') AS stopped_seeking_help,
			COALESCE(NULLIF(ash.victim_case_form2_enums_name, ''), '') AS aggressor_sexually_harassment_2,
			COALESCE(NULLIF(avp.victim_case_form2_enums_name, ''), '') AS aggressor_taken_advantage_physical_vulnerabil,
			COALESCE(NULLIF(vmg.victim_case_form2_enums_name, ''), '') AS violence_motivated_by_gender_2,

			-- Bloque 1: Autorización y Contacto
			'' AS authorization_answer,
			COALESCE((
				SELECT string_agg(adj_e.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 adj_rel
				JOIN salvia.victim_case_form2_enums adj_e ON adj_e.victim_case_form2_enums_id = adj_rel.victim_case_form2_enums_id
				WHERE adj_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND adj_e.victim_case_form2_enums_category LIKE '%adjustments_gbv%'
			), '') AS adjustments_gbv,
			COALESCE(NULLIF(interp.victim_case_form2_enums_name, ''), '') AS require_interpreter,
			COALESCE(NULLIF(vf2.victim_case_form2_support_contact_email, ''), '') AS support_contact_email,

			-- Bloque 2: Hechos Victimizantes (ampliación)
			COALESCE(facts_dept.department_name, '') AS facts_dept_name,
			COALESCE(facts_city.city_name, '') AS facts_city_name,
			COALESCE(NULLIF(fzone.victim_case_form2_enums_name, ''), '') AS facts_zone,
			COALESCE(NULLIF(vf2.victim_case_form2_facts_address, ''), '') AS facts_address_form2,
			COALESCE((
				SELECT string_agg(vte.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 vt_rel
				JOIN salvia.victim_case_form2_enums vte ON vte.victim_case_form2_enums_id = vt_rel.victim_case_form2_enums_id
				WHERE vt_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND (vte.victim_case_form2_enums_category LIKE '%violence_experienced%' OR vte.victim_case_form2_enums_category LIKE '%type_of_experienced_violence%')
			), '') AS violence_type,
			COALESCE((
				SELECT string_agg(ste.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 st_rel
				JOIN salvia.victim_case_form2_enums ste ON ste.victim_case_form2_enums_id = st_rel.victim_case_form2_enums_id
				WHERE st_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND ste.victim_case_form2_enums_category LIKE '%subtype_violence%'
			), '') AS violence_subtype,
			COALESCE((
				SELECT string_agg(soe.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 so_rel
				JOIN salvia.victim_case_form2_enums soe ON soe.victim_case_form2_enums_id = so_rel.victim_case_form2_enums_id
				WHERE so_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND soe.victim_case_form2_enums_category LIKE '%scope_of_violence%'
			), '') AS violence_scope_form2,
			COALESCE(NULLIF(wps.victim_case_form2_enums_name, ''), '') AS workplace_sector,
			COALESCE(NULLIF(recur.victim_case_form2_enums_name, ''), '') AS recurrence_aggression,

			-- Bloque 3: Características del Agresor (ampliación)
			COALESCE(NULLIF(nag.victim_case_form2_enums_name, ''), '') AS num_aggressors,
			COALESCE(NULLIF(prox.victim_case_form2_enums_name, ''), '') AS proximity_aggressor,
			COALESCE(NULLIF(ecodep.victim_case_form2_enums_name, ''), '') AS economically_dependent,
			COALESCE(NULLIF(aggi.victim_case_form2_enums_name, ''), '') AS aggressor_gender_identity,

			-- Bloque 4: Tamizaje (preguntas adicionales)
			COALESCE(NULLIF(b4_pursues.victim_case_form2_enums_name, ''), '') AS aggressor_pursues_spies,
			COALESCE(NULLIF(b4_weapons.victim_case_form2_enums_name, ''), '') AS aggressor_has_access_weapons,
			COALESCE(NULLIF(b4_unemp.victim_case_form2_enums_name, ''), '') AS partner_unemployed,
			COALESCE(NULLIF(b4_denunc.victim_case_form2_enums_name, ''), '') AS partner_other_denunciations,
			COALESCE(NULLIF(b4_penal.victim_case_form2_enums_name, ''), '') AS aggressor_penal_background,
			COALESCE(NULLIF(b4_strang.victim_case_form2_enums_name, ''), '') AS aggressor_strangulation,
			COALESCE(NULLIF(b4_drugs.victim_case_form2_enums_name, ''), '') AS aggressor_consumes_drugs,
			COALESCE(NULLIF(b4_alcohol.victim_case_form2_enums_name, ''), '') AS aggressor_is_alcoholic,
			COALESCE(NULLIF(b4_controls.victim_case_form2_enums_name, ''), '') AS partner_controls,
			COALESCE(NULLIF(b4_suicide.victim_case_form2_enums_name, ''), '') AS partner_threatened_suicide,
			COALESCE(NULLIF(b4_damage.victim_case_form2_enums_name, ''), '') AS partner_threatened_damage,
			COALESCE(NULLIF(b4_selfharm.victim_case_form2_enums_name, ''), '') AS thoughts_of_self_harm,
			COALESCE(NULLIF(b4_limits.victim_case_form2_enums_name, ''), '') AS aggressor_limits_contact,
			COALESCE(NULLIF(b4_lives.victim_case_form2_enums_name, ''), '') AS still_lives_with_aggressor,

			-- Bloque 5: Datos Personales Víctima (ampliación)
			COALESCE(TO_CHAR(vf2.victim_case_form2_birth_date, 'YYYY-MM-DD'), '') AS birth_date_form2,
			COALESCE(NULLIF(b5_diff.victim_case_form2_enums_name, ''), '') AS physical_difficulties,
			COALESCE(NULLIF(b5_nat.victim_case_form2_enums_name, ''), '') AS nationality_form2,
			COALESCE(NULLIF(b5_snat.victim_case_form2_enums_name, ''), '') AS specified_nationality,
			COALESCE(NULLIF(b5_mig.victim_case_form2_enums_name, ''), '') AS migration_condition,
			COALESCE(NULLIF(gi.victim_case_form2_enums_name, ''), '') AS gender_identity_form2,
			COALESCE(NULLIF(b5_sex.victim_case_form2_enums_name, ''), '') AS assigned_sex_at_birth,
			COALESCE((
				SELECT string_agg(spp_e.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 spp_rel
				JOIN salvia.victim_case_form2_enums spp_e ON spp_e.victim_case_form2_enums_id = spp_rel.victim_case_form2_enums_id
				WHERE spp_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND spp_e.victim_case_form2_enums_category LIKE '%specially_protected%'
			), '') AS special_protected_pop,
			COALESCE(NULLIF(b5_edu.victim_case_form2_enums_name, ''), '') AS last_education_level,
			COALESCE(NULLIF(b5_income.victim_case_form2_enums_name, ''), '') AS income_generation_method,
			COALESCE(NULLIF(b5_employ.victim_case_form2_enums_name, ''), '') AS employment_relationship,
			COALESCE((
				SELECT string_agg(asp_e.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 asp_rel
				JOIN salvia.victim_case_form2_enums asp_e ON asp_e.victim_case_form2_enums_id = asp_rel.victim_case_form2_enums_id
				WHERE asp_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND asp_e.victim_case_form2_enums_category LIKE '%asp_mode%'
			), '') AS asp_mode,
			COALESCE(TO_CHAR(vf2.victim_case_form2_approx_start_asp, 'YYYY-MM-DD'), '') AS approx_start_asp,
			COALESCE((
				SELECT string_agg(rasp_e.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 rasp_rel
				JOIN salvia.victim_case_form2_enums rasp_e ON rasp_e.victim_case_form2_enums_id = rasp_rel.victim_case_form2_enums_id
				WHERE rasp_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND rasp_e.victim_case_form2_enums_category LIKE '%reason_asp%'
			), '') AS reason_asp,
			COALESCE(NULLIF(b5_housing.victim_case_form2_enums_name, ''), '') AS housing_tenancy,
			COALESCE(NULLIF(b5_stratum.victim_case_form2_enums_name, ''), '') AS housing_stratum,
			COALESCE((
				SELECT string_agg(dep_e.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 dep_rel
				JOIN salvia.victim_case_form2_enums dep_e ON dep_e.victim_case_form2_enums_id = dep_rel.victim_case_form2_enums_id
				WHERE dep_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND dep_e.victim_case_form2_enums_category LIKE '%has_dependents%'
			), '') AS has_dependents,
			COALESCE(NULLIF(b5_preg.victim_case_form2_enums_name, ''), '') AS currently_pregnant,
			'' AS residence_zone,

			-- Bloque 6: Plan de Acción, Denuncia, Operación
			COALESCE((
				SELECT string_agg(ap_e.victim_case_form2_enums_name, ', ')
				FROM salvia.rel_victim_case_form2_enums_victim_case_form2 ap_rel
				JOIN salvia.victim_case_form2_enums ap_e ON ap_e.victim_case_form2_enums_id = ap_rel.victim_case_form2_enums_id
				WHERE ap_rel.victim_case_form2_id = vf2.victim_case_form2_id
				AND ap_e.victim_case_form2_enums_category LIKE '%action_plan%'
			), '') AS action_plan,
			COALESCE(NULLIF(vf2.victim_case_form2_saliva_management_explanation, ''), '') AS management_explanation,
			COALESCE(NULLIF(b6_easy.victim_case_form2_enums_name, ''), '') AS allows_easy_report,
			COALESCE((SELECT gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names FROM security.general_user gu JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile WHERE gu.general_user_i_code = vc.agent_id LIMIT 1), '') AS agent_responsible,
			COALESCE(vc.victim_case_owner_description, '') AS owner_description

		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form1 vf1 ON vf1.victim_case_form1_victim_case = vc.victim_case_id
		LEFT JOIN salvia.victim_case_form2 vf2 ON vf2.victim_case_form2_victim_case = vc.victim_case_id
		LEFT JOIN salvia.victim_contact vco ON vco.victim_contact_id = vc.victim_case_victim_contact
		LEFT JOIN (
			SELECT DISTINCT ON (victim_contact_form1_victim_contact) *
			FROM salvia.victim_contact_form1
			ORDER BY victim_contact_form1_victim_contact, victim_contact_form1_id DESC
		) vcf1 ON vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
		LEFT JOIN salvia.victim_case_form2_enums gi ON gi.victim_case_form2_enums_id = vf2.victim_case_form2_gender_identity
		LEFT JOIN salvia.victim_case_form2_enums rel ON rel.victim_case_form2_enums_id = vf2.victim_case_form2_relationship_with_presumed_aggressor
		LEFT JOIN salvia.victim_case_form2_enums so ON so.victim_case_form2_enums_id = vf2.victim_case_form2_sexual_orientation
		LEFT JOIN salvia.victim_case_form2_enums ms ON ms.victim_case_form2_enums_id = vf2.victim_case_form2_marital_status
		LEFT JOIN salvia.victim_case_form2_enums dis ON dis.victim_case_form2_enums_id = vf2.victim_case_form2_person_with_disability
		LEFT JOIN salvia.victim_case_form2_enums eth ON eth.victim_case_form2_enums_id = vf2.victim_case_form2_ethnic_affiliation
		LEFT JOIN salvia.victim_case_form2_enums kin ON kin.victim_case_form2_enums_id = vf2.victim_case_form2_support_contact_kinship
		LEFT JOIN salvia.victim_case_form2_enums sv ON sv.victim_case_form2_enums_id = vf2.victim_case_form2_scenario_violence
		LEFT JOIN salvia.victim_case_form2_enums oc2 ON oc2.victim_case_form2_enums_id = vf2.victim_case_form2_occupation
		LEFT JOIN salvia.victim_case_form2_enums htb ON htb.victim_case_form2_enums_id = vf2.victim_case_form2_victim_health_to_blackmail
		LEFT JOIN salvia.victim_case_form2_enums tro ON tro.victim_case_form2_enums_id = vf2.victim_case_form2_threatened_reveal_sexual_orientation
		LEFT JOIN salvia.victim_case_form2_enums ssh ON ssh.victim_case_form2_enums_id = vf2.victim_case_form2_stopped_seeking_help
		LEFT JOIN salvia.victim_case_form2_enums ash ON ash.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_sexually_harassment_2
		LEFT JOIN salvia.victim_case_form2_enums avp ON avp.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_taken_advantage_physical_vulnerabil
		LEFT JOIN salvia.victim_case_form2_enums vmg ON vmg.victim_case_form2_enums_id = vf2.victim_case_form2_violence_motivated_by_gender_2
		LEFT JOIN salvia.victim_case_form2_enums interp ON interp.victim_case_form2_enums_id = vf2.victim_case_form2_require_language_interpreter
		LEFT JOIN security.town facts_town ON facts_town.town_code = vf2.victim_case_form2_facts_town_code
		LEFT JOIN security.city facts_city ON facts_city.city_id = facts_town.city_id
		LEFT JOIN security.department facts_dept ON facts_dept.department_id = facts_city.department_id
		LEFT JOIN salvia.victim_case_form2_enums fzone ON fzone.victim_case_form2_enums_id = vf2.victim_case_form2_facts_zone
		LEFT JOIN salvia.victim_case_form2_enums wps ON wps.victim_case_form2_enums_id = vf2.victim_case_form2_workplace_sector_occurrence
		LEFT JOIN salvia.victim_case_form2_enums recur ON recur.victim_case_form2_enums_id = vf2.victim_case_form2_recurrence_aggression
		LEFT JOIN salvia.victim_case_form2_enums nag ON nag.victim_case_form2_enums_id = vf2.victim_case_form2_num_agressors
		LEFT JOIN salvia.victim_case_form2_enums prox ON prox.victim_case_form2_enums_id = vf2.victim_case_form2_proximity_principal_aggressor
		LEFT JOIN salvia.victim_case_form2_enums ecodep ON ecodep.victim_case_form2_enums_id = vf2.victim_case_form2_economically_dependent
		LEFT JOIN salvia.victim_case_form2_enums aggi ON aggi.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_gender_identity
		LEFT JOIN salvia.victim_case_form2_enums b4_pursues ON b4_pursues.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_pursues_spies_destroys
		LEFT JOIN salvia.victim_case_form2_enums b4_weapons ON b4_weapons.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_has_access_to_weapons
		LEFT JOIN salvia.victim_case_form2_enums b4_unemp ON b4_unemp.victim_case_form2_enums_id = vf2.victim_case_form2_partner_unemployed
		LEFT JOIN salvia.victim_case_form2_enums b4_denunc ON b4_denunc.victim_case_form2_enums_id = vf2.victim_case_form2_partner_other_denunciations
		LEFT JOIN salvia.victim_case_form2_enums b4_penal ON b4_penal.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_has_penal_background
		LEFT JOIN salvia.victim_case_form2_enums b4_strang ON b4_strang.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_attempted_strangulation
		LEFT JOIN salvia.victim_case_form2_enums b4_drugs ON b4_drugs.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_consumes_drugs
		LEFT JOIN salvia.victim_case_form2_enums b4_alcohol ON b4_alcohol.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_is_alcoholic
		LEFT JOIN salvia.victim_case_form2_enums b4_controls ON b4_controls.victim_case_form2_enums_id = vf2.victim_case_form2_partner_controls
		LEFT JOIN salvia.victim_case_form2_enums b4_suicide ON b4_suicide.victim_case_form2_enums_id = vf2.victim_case_form2_partner_threatened_suicide
		LEFT JOIN salvia.victim_case_form2_enums b4_damage ON b4_damage.victim_case_form2_enums_id = vf2.victim_case_form2_partner_threatened_damage_members
		LEFT JOIN salvia.victim_case_form2_enums b4_selfharm ON b4_selfharm.victim_case_form2_enums_id = vf2.victim_case_form2_thoughts_of_self_harm
		LEFT JOIN salvia.victim_case_form2_enums b4_limits ON b4_limits.victim_case_form2_enums_id = vf2.victim_case_form2_aggressor_limits_contact_support_networks
		LEFT JOIN salvia.victim_case_form2_enums b4_lives ON b4_lives.victim_case_form2_enums_id = vf2.victim_case_form2_still_lives_with_aggressor
		LEFT JOIN salvia.victim_case_form2_enums b5_diff ON b5_diff.victim_case_form2_enums_id = vf2.victim_case_form2_physical_mental_sensory_difficulties
		LEFT JOIN salvia.victim_case_form2_enums b5_nat ON b5_nat.victim_case_form2_enums_id = vf2.victim_case_form2_nationality
		LEFT JOIN salvia.victim_case_form2_enums b5_snat ON b5_snat.victim_case_form2_enums_id = vf2.victim_case_form2_specified_nationality
		LEFT JOIN salvia.victim_case_form2_enums b5_mig ON b5_mig.victim_case_form2_enums_id = vf2.victim_case_form2_migration_condition
		LEFT JOIN salvia.victim_case_form2_enums b5_sex ON b5_sex.victim_case_form2_enums_id = vf2.victim_case_form2_assigned_sex_at_birth
		LEFT JOIN salvia.victim_case_form2_enums b5_edu ON b5_edu.victim_case_form2_enums_id = vf2.victim_case_form2_last_education_level
		LEFT JOIN salvia.victim_case_form2_enums b5_income ON b5_income.victim_case_form2_enums_id = vf2.victim_case_form2_income_generation_method
		LEFT JOIN salvia.victim_case_form2_enums b5_employ ON b5_employ.victim_case_form2_enums_id = vf2.victim_case_form2_employment_relationship
		LEFT JOIN salvia.victim_case_form2_enums b5_housing ON b5_housing.victim_case_form2_enums_id = vf2.victim_case_form2_housing_tenancy_form
		LEFT JOIN salvia.victim_case_form2_enums b5_stratum ON b5_stratum.victim_case_form2_enums_id = vf2.victim_case_form2_housing_stratum
		LEFT JOIN salvia.victim_case_form2_enums b5_preg ON b5_preg.victim_case_form2_enums_id = vf2.victim_case_form2_currently_pregnant
		LEFT JOIN salvia.victim_case_form2_enums b6_easy ON b6_easy.victim_case_form2_enums_id = vf2.victim_case_form2_allows_easy_report
		LEFT JOIN security.town t ON t.town_code = vc.victim_case_victim_town_code
		LEFT JOIN security.city c ON c.city_id = t.city_id
		LEFT JOIN security.department d ON d.department_id = c.department_id
		WHERE vc.victim_case_creation_date BETWEEN ? AND ?
		ORDER BY vc.victim_case_creation_date ASC
	`
	err := r.db.WithContext(ctx).Raw(query, start, end).Scan(&cases).Error
	return cases, err
}

func (r *reportRepository) FindFollowUpsByCaseICodes(ctx context.Context, caseICodes []string) ([]FollowUpReportDTO, error) {
	var followUps []FollowUpReportDTO
	if len(caseICodes) == 0 {
		return followUps, nil
	}
	query := `
		SELECT 
			fu.*,
			vc.victim_case_victim_doc_number AS victim_doc_number
		FROM salvia.follow_up_v2 fu
		JOIN salvia.victim_case vc ON vc.victim_case_i_code = fu.case_id
		WHERE fu.case_id IN ? AND fu.deleted_at IS NULL
		ORDER BY fu.scheduled_date ASC, fu.sequence_number ASC
	`
	err := r.db.WithContext(ctx).Raw(query, caseICodes).Scan(&followUps).Error
	return followUps, err
}

func (r *reportRepository) FindAnswersByFormSubmissions(ctx context.Context, submissionIDs []string) ([]AnswerDTO, error) {
	var answers []AnswerDTO
	if len(submissionIDs) == 0 {
		return answers, nil
	}
	query := `
		SELECT 
			fu.case_id AS case_icode,
			fu.id AS follow_up_id,
			a.form_submission_id,
			q.description AS question_desc,
			COALESCE(
				NULLIF(opt.label, ''),
				NULLIF(town_ans.town_name, ''),
				NULLIF(city_ans.city_name, ''),
				NULLIF(dept_ans.department_name, ''),
				a.value
			) AS answer_value,
			vc.victim_case_victim_doc_number AS victim_doc_number,
			fu.completed_at AS follow_up_date
		FROM salvia.answer a
		JOIN salvia.question q ON q.id = a.question_id
		LEFT JOIN salvia.form_section s ON s.id = q.form_section_id::uuid
		LEFT JOIN salvia.option opt ON opt.question_id = a.question_id AND opt.value = a.value AND opt.deleted_at IS NULL
		LEFT JOIN security.town town_ans ON (town_ans.town_code = a.value OR town_ans.town_id::text = a.value) AND LENGTH(a.value) >= 1 AND a.value ~ '^\d+$'
		LEFT JOIN security.city city_ans ON city_ans.city_id::text = a.value AND LENGTH(a.value) <= 5 AND a.value ~ '^\d+$' AND town_ans.town_id IS NULL
		LEFT JOIN security.department dept_ans ON dept_ans.department_id::text = a.value AND LENGTH(a.value) <= 3 AND a.value ~ '^\d+$' AND town_ans.town_id IS NULL AND city_ans.city_id IS NULL
		JOIN salvia.follow_up_v2 fu ON fu.form_submission_id = a.form_submission_id
		JOIN salvia.victim_case vc ON vc.victim_case_i_code = fu.case_id
		WHERE a.form_submission_id IN ?
		  AND a.deleted_at IS NULL
		  AND fu.deleted_at IS NULL
		ORDER BY fu.completed_at ASC, fu.case_id ASC, fu.sequence_number ASC, s.order ASC, q.order ASC
	`
	err := r.db.WithContext(ctx).Raw(query, submissionIDs).Scan(&answers).Error
	return answers, err
}

func (r *reportRepository) FindTimelineEventsByCaseICodes(ctx context.Context, caseICodes []string) ([]TimelineReportDTO, error) {
	var events []TimelineReportDTO
	if len(caseICodes) == 0 {
		return events, nil
	}
	query := `
		SELECT 
			ev.*,
			vc.victim_case_victim_doc_number AS victim_doc_number
		FROM salvia.case_timeline_event ev
		JOIN salvia.victim_case vc ON vc.victim_case_i_code = ev.case_id
		WHERE ev.case_id IN ? AND ev.deleted_at IS NULL
		ORDER BY ev.date ASC, ev.created_at ASC
	`
	err := r.db.WithContext(ctx).Raw(query, caseICodes).Scan(&events).Error
	return events, err
}

func (r *reportRepository) FindContactsInDateRange(ctx context.Context, start, end time.Time) ([]ContactReportDTO, error) {
	var contacts []ContactReportDTO
	query := `
		SELECT
			vc.victim_contact_id,
			vc.victim_contact_i_code,
			vc.victim_contact_creation_date,
			vc.victim_contact_update_date,
			COALESCE(vc.victim_contact_status_description, '') AS victim_contact_status_description,
			COALESCE(vc.victim_contact_names, '') AS victim_contact_names,
			COALESCE(vc.victim_contact_last_names, '') AS victim_contact_last_names,
			COALESCE(vc.victim_contact_latitude, 0) AS victim_contact_latitude,
			COALESCE(vc.victim_contact_longitude, 0) AS victim_contact_longitude,
			CASE
				WHEN f1.victim_contact_form1_id IS NOT NULL THEN 'f1'
				WHEN f2.victim_contact_form2_id IS NOT NULL THEN 'f2'
				ELSE ''
			END AS form_type,

			COALESCE(f1.victim_contact_form1_nick, '') AS f1_nick,
			COALESCE(f1.victim_contact_form1_doc_type, '') AS f1_doc_type,
			COALESCE(f1.victim_contact_form1_doc_number::text, '') AS f1_doc_number,
			COALESCE(TO_CHAR(f1.victim_contact_form1_birth_date, 'DD/MM/YYYY'), '') AS f1_birth_date,
			COALESCE(f1.victim_contact_form1_town_code, '') AS f1_town_code,
			COALESCE(t.town_name, '') AS f1_town_name,
			COALESCE(f1.victim_contact_form1_address, '') AS f1_address,
			COALESCE(f1.victim_contact_form1_phone::text, '') AS f1_phone,
			COALESCE(f1.victim_contact_form1_gender_identity, '') AS f1_gender_identity,
			COALESCE(f1.victim_contact_form1_sexual_orientation, '') AS f1_sexual_orientation,
			COALESCE(f1.victim_contact_form1_origin, '') AS f1_origin,
			COALESCE(f1.victim_contact_form1_occupation, '') AS f1_occupation,
			COALESCE(f1.victim_contact_form1_occupation_other, '') AS f1_occupation_other,
			COALESCE(f1.victim_contact_form1_facts_description, '') AS f1_facts_description,

			COALESCE(f2.victim_contact_form2_reporter_names, '') AS f2_reporter_names,
			COALESCE(f2.victim_contact_form2_reporter_phone::text, '') AS f2_reporter_phone,
			COALESCE(f2.victim_contact_form2_victim_col_phone::text, '') AS f2_victim_col_phone,
			COALESCE(f2.victim_contact_form2_facts_description, '') AS f2_facts_description,
			COALESCE(TO_CHAR(f2.victim_contact_form2_best_contact_time, 'HH24:MI'), '') AS f2_best_contact_time,
			COALESCE(wrc.victim_case_form2_enums_name, '') AS f2_will_receive_call,
			COALESCE(hcr.victim_case_form2_enums_name, '') AS f2_has_care_role,
			COALESCE(vaw.victim_case_form2_enums_name, '') AS f2_victim_aware,
			COALESCE(rt.victim_case_form2_enums_name, '') AS f2_report_type,
			COALESCE(rel_enums.adjustments_gbv, '') AS f2_adjustments_gbv

		FROM salvia.victim_contact vc
		LEFT JOIN (
			SELECT DISTINCT ON (victim_contact_form1_victim_contact) *
			FROM salvia.victim_contact_form1
			ORDER BY victim_contact_form1_victim_contact, victim_contact_form1_id DESC
		) f1 ON f1.victim_contact_form1_victim_contact = vc.victim_contact_id
		LEFT JOIN (
			SELECT DISTINCT ON (victim_contact_form2_victim_contact) *
			FROM salvia.victim_contact_form2
			ORDER BY victim_contact_form2_victim_contact, victim_contact_form2_id DESC
		) f2 ON f2.victim_contact_form2_victim_contact = vc.victim_contact_id
		LEFT JOIN security.town t ON t.town_code = f1.victim_contact_form1_town_code
		LEFT JOIN salvia.victim_case_form2_enums wrc ON wrc.victim_case_form2_enums_id = f2.victim_contact_form2_will_receive_call
		LEFT JOIN salvia.victim_case_form2_enums hcr ON hcr.victim_case_form2_enums_id = f2.victim_contact_form2_has_care_role
		LEFT JOIN salvia.victim_case_form2_enums vaw ON vaw.victim_case_form2_enums_id = f2.victim_contact_form2_victim_aware_of_report
		LEFT JOIN salvia.victim_case_form2_enums rt ON rt.victim_case_form2_enums_id = f2.victim_contact_form2_report_type
		LEFT JOIN (
			SELECT rel.victim_contact_form2_id AS form2_id,
			       STRING_AGG(e.victim_case_form2_enums_name, ', ' ORDER BY e.victim_case_form2_enums_id) AS adjustments_gbv
			FROM salvia.rel_victim_case_form2_enums_victim_contact_form2 rel
			JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = rel.victim_case_form2_enums_id
			WHERE e.victim_case_form2_enums_category = 'victim_case_form2_adjustments_gbv'
			GROUP BY rel.victim_contact_form2_id
		) rel_enums ON rel_enums.form2_id = f2.victim_contact_form2_id
		WHERE vc.victim_contact_creation_date BETWEEN ? AND ?
		ORDER BY vc.victim_contact_creation_date ASC
	`
	err := r.db.WithContext(ctx).Raw(query, start, end).Scan(&contacts).Error
	return contacts, err
}
