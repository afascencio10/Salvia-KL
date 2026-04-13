package models

type FollowUpDetailResponse struct {
	FollowUp    FollowUpV2      `json:"followUp"`
	CaseInfo    VictimCaseLight `json:"caseInfo"`
	Barriers               []BarrierV2             `json:"barriers"`
	Permissions            Permissions             `json:"permissions"`
	EmergencyMeasures      []EmergencyMeasure      `json:"emergencyMeasures"`
	PsychosocialSupports   []PsychosocialSupport   `json:"psychosocialSupports"`
	EconomicStabilizations []EconomicStabilization `json:"economicStabilizations"`
}

type Permissions struct {
	CanEdit bool `json:"canEdit"`
}
