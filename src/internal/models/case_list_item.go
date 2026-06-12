package models

import "time"

// OpenBarrierItem barrera OPEN asociada a un caso en el listado (E-01).
type OpenBarrierItem struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Sector string `json:"sector"`
}

// CaseListItem representa un caso en la respuesta de GET /api/v1/cases/list.
type CaseListItem struct {
	ID               int64      `json:"id"`
	ICode            string     `json:"i_code"`
	Names            string     `json:"names"`
	LastNames        string     `json:"lastNames"`
	DocNumber        string     `json:"docNumber"`
	VictimPhone      string     `json:"victimPhone"`
	CreationDate     time.Time  `json:"creationDate"`
	Status           string     `json:"status"`
	OwnerNames       string     `json:"ownerNames"`
	OwnerLastNames   string     `json:"ownerLastNames"`
	OwnerTeam        string     `json:"ownerTeam"`
	CaseTeam         string     `json:"caseTeam"` // victim_case.victim_case_team (E-11)
	RiskStatus       *string    `json:"riskStatus"`
	NextFollowUpDate        *time.Time `json:"nextFollowUpDate"`
	CompletedFollowUpsCount int               `json:"completedFollowUpsCount"`
	OpenBarriers            []OpenBarrierItem `json:"openBarriers"`
}
