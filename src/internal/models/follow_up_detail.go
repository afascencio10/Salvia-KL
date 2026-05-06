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
}

type Permissions struct {
	CanEdit bool `json:"canEdit"`
}
