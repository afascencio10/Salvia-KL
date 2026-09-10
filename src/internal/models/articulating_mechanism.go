package models

import (
	"time"

	"gorm.io/gorm"
)

// ArticulatingMechanism es el registro 1:1 del usuario con rol mc:
// tipo (municipal / departamental / nacional) y la ubicación que le corresponde.
type ArticulatingMechanism struct {
	ID   string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	Type string `gorm:"type:varchar(20);not null;index;column:type" json:"type"`

	DepartmentID   *string `gorm:"type:varchar(36);column:department_id" json:"departmentId,omitempty"`
	CityID         *string `gorm:"type:varchar(36);column:city_id" json:"cityId,omitempty"`
	TownCode       *string `gorm:"type:varchar(16);column:town_code" json:"townCode,omitempty"`
	DepartmentName *string `gorm:"type:varchar(128);column:department_name" json:"departmentName,omitempty"`
	CityName       *string `gorm:"type:varchar(128);column:city_name" json:"cityName,omitempty"`
	TownName       *string `gorm:"type:varchar(128);column:town_name" json:"townName,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ArticulatingMechanism) TableName() string {
	return "salvia.articulating_mechanism"
}
