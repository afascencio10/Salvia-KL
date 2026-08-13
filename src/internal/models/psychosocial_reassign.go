package models

// PsTsProfessionalRow fila de profesionales activos con rol ps o ts (RRM-03).
type PsTsProfessionalRow struct {
	ICode    string `gorm:"column:icode"`
	FullName string `gorm:"column:full_name"`
	Role     string `gorm:"column:role"`
}

// DuplaReassignRow dupla activa con ids y nombres de integrantes (RRM-03 / RRM-05).
type DuplaReassignRow struct {
	ID               string `gorm:"column:id"`
	Name             string `gorm:"column:name"`
	PsychologistID   string `gorm:"column:psychologist_id"`
	PsychologistName string `gorm:"column:psychologist_name"`
	SocialWorkerID   string `gorm:"column:social_worker_id"`
	SocialWorkerName string `gorm:"column:social_worker_name"`
}
