package models

// PsTsProfessionalRow fila de profesionales activos con rol ps o ts (RRM-03).
type PsTsProfessionalRow struct {
	ICode    string `gorm:"column:icode"`
	FullName string `gorm:"column:full_name"`
	Role     string `gorm:"column:role"`
}

// DuplaReassignRow dupla activa con nombres de integrantes (RRM-03).
type DuplaReassignRow struct {
	ID               string `gorm:"column:id"`
	Name             string `gorm:"column:name"`
	PsychologistName string `gorm:"column:psychologist_name"`
	SocialWorkerName string `gorm:"column:social_worker_name"`
}
