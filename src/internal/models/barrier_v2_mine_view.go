package models

// BarrierV2MineItem es el DTO de lectura para "Mis Barreras".
// Lista barrier_v2 asignadas al Enlace Territorial (user_territory_link_id),
// enriquecidas con victim_case, entity_branch y una tarea ToDo opcional
// de tipo articulacion_territorial asignada al mismo usuario.
type BarrierV2MineItem struct {
	ID          string `gorm:"column:id"           json:"id"`
	CaseID      string `gorm:"column:case_id"      json:"caseId"`
	Sector      string `gorm:"column:sector"       json:"sector"`
	Description string `gorm:"column:description"  json:"description"`
	Status      string `gorm:"column:status"       json:"status"`

	EntityBranchID   *int64 `gorm:"column:entity_branch_id"   json:"entityBranchId,omitempty"`
	EntityBranchName string `gorm:"column:entity_branch_name" json:"entityBranchName"`

	VictimName      string `gorm:"column:victim_name"       json:"victimName"`
	VictimLastName  string `gorm:"column:victim_last_name"  json:"victimLastName"`
	VictimDocNumber string `gorm:"column:victim_doc_number" json:"victimDocNumber"`
	VictimPhone     string `gorm:"column:victim_phone"      json:"victimPhone"`
	CaseCode        string `gorm:"column:case_code"         json:"caseCode"`

	CityID   string `gorm:"column:city_id"   json:"cityId"`
	CityName string `gorm:"column:city_name" json:"cityName"`

	// Tarea ToDo de articulación (vacío si no existe).
	TaskID          string `gorm:"column:task_id"          json:"taskId,omitempty"`
	TaskType        string `gorm:"column:task_type"        json:"taskType,omitempty"`
	TaskDescription string `gorm:"column:task_description" json:"taskDescription,omitempty"`
}
