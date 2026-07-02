package models

// AgentPsicosocialLight agente del equipo psicosocial para autocomplete E-07.
type AgentPsicosocialLight struct {
	ICode     string `json:"icode" gorm:"column:icode"`
	Names     string `json:"names" gorm:"column:names"`
	LastNames string `json:"lastNames" gorm:"column:last_names"`
	Team      string `json:"team" gorm:"column:team"`
	Specialty string `json:"specialty" gorm:"column:specialty"`
}
