package models

// AgentLight representa información básica de un agente del sistema
type AgentLight struct {
	ICode     string `json:"icode" gorm:"column:general_user_i_code;primaryKey"`
	Login     string `json:"login" gorm:"column:general_user_login"`
	Names     string `json:"names" gorm:"column:general_user_profile_names"`
	LastNames string `json:"last_names" gorm:"column:general_user_profile_last_names"`
	Team      string `json:"team" gorm:"column:general_user_team"`
}
