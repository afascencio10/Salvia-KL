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
			COALESCE(NULLIF(vmg.victim_case_form2_enums_name, ''), '') AS violence_motivated_by_gender_2

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
			a.value AS answer_value,
			vc.victim_case_victim_doc_number AS victim_doc_number,
			fu.completed_at AS follow_up_date
		FROM salvia.answer a
		JOIN salvia.question q ON q.id = a.question_id
		LEFT JOIN salvia.form_section s ON s.id = q.form_section_id::uuid
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
