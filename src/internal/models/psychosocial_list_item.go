package models

import "time"

// PsychosocialListItem fila del listado GET /api/v1/psychosocial-support/list (E-01).
type PsychosocialListItem struct {
	ID               string    `json:"id"`
	CaseICode        string    `json:"caseICode"`
	FollowUpID       string    `json:"followUpId"`
	Status           string    `json:"status"`
	SessionCount     int       `json:"sessionCount"`
	CreatedAt        time.Time `json:"createdAt"`
	SubmittedByName  string    `json:"submittedByName"`
	SubmittedByTeam  string    `json:"submittedByTeam"`
	DuplaID          *string   `json:"duplaId"`
	DuplaName        string    `json:"duplaName"`
	PsychologistName string    `json:"psychologistName"`
	SocialWorkerName string    `json:"socialWorkerName"`
	ProfessionalID     *string   `json:"professionalId"`
	ProfessionalName   string    `json:"professionalName"`
	ProfessionalTeam   string    `json:"professionalTeam,omitempty"`
	ProfessionalRole   string    `json:"professionalRole"`
	VictimNames      string    `json:"victimNames"`
	VictimLastNames  string    `json:"victimLastNames"`
	DocNumber        string    `json:"docNumber"`
	VictimPhone      string    `json:"victimPhone"`
	Municipality     string    `json:"municipality"`
	RiskLevel        int       `json:"riskLevel"`
}

// DuplaOption ítem del catálogo GET /api/v1/duplas.
type DuplaOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PsychosocialListStats métricas GET /api/v1/psychosocial-support/stats.
type PsychosocialListStats struct {
	Total        int64 `json:"total"`
	Abierto      int64 `json:"abierto"`
	EnGestion    int64 `json:"enGestion"`
	EnDevolucion int64 `json:"enDevolucion"`
	Cerrado      int64 `json:"cerrado"`
}
