package models

// FollowUpVictimInfo contiene los datos de la víctima para la pantalla de detalle.
// Replica repository.VictimCaseInfo para evitar dependencia circular.
type FollowUpVictimInfo struct {
	Names             string `json:"Names"`
	LastNames         string `json:"LastNames"`
	TownName          string `json:"TownName"`
	Phone             string `json:"Phone"`
	GenderIdentity    string `json:"GenderIdentity"`
	SexualOrientation string `json:"SexualOrientation"`
	ContactPhone      string `json:"ContactPhone"`
	Age               *int64 `json:"Age"`
}

type FollowUpDetailResponse struct {
	FollowUp               FollowUpV2              `json:"followUp"`
	CaseInfo               VictimCaseLight         `json:"caseInfo"`
	VictimInfo             *FollowUpVictimInfo     `json:"victimInfo,omitempty"`
	Barriers               []BarrierV2             `json:"barriers"`
	Permissions            Permissions             `json:"permissions"`
	EmergencyMeasures      []EmergencyMeasure      `json:"emergencyMeasures"`
	PsychosocialSupports   []PsychosocialSupport   `json:"psychosocialSupports"`
	EconomicStabilizations []EconomicStabilization `json:"economicStabilizations"`
	FormAnswers            *FollowUpFormAnswers    `json:"formAnswers,omitempty"`
}

type Permissions struct {
	CanEdit bool `json:"canEdit"`
}

// FollowUpFormAnswers reúne, en texto legible, un subconjunto fijo de respuestas
// del formulario dinámico de seguimiento (leídas por form_submission_id) para
// mostrarlas en el tab Resumen de Detalle de Seguimiento.
type FollowUpFormAnswers struct {
	RiskAnalysis      string                 `json:"riskAnalysis"`      // Valoración del Riesgo — análisis de factores
	CaseManagement    string                 `json:"caseManagement"`    // Seguimiento de Caso — gestión realizada
	ReferralEvidence  string                 `json:"referralEvidence"`  // Seguimiento de Caso — elementos que evidencian la remisión
	BarrierFollowUps  []BarrierAnswerSummary `json:"barrierFollowUps"`  // Seguimiento a Barreras (repeater, barreras ya activas)
	BarrierIdentified []BarrierAnswerSummary `json:"barrierIdentified"` // Identificación de Barreras (repeater, barreras nuevas)
}

// BarrierAnswerSummary es una entrada de uno de los 2 repeaters de barreras.
// Actuaciones queda vacío para entries de Identificación de Barreras (esa
// pregunta solo existe en Seguimiento a Barreras).
type BarrierAnswerSummary struct {
	Actuaciones string `json:"actuaciones,omitempty"`
	Gestion     string `json:"gestion,omitempty"`
}
