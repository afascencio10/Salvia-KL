package models

// DepartmentLight es una proyección de solo lectura de security.department.
type DepartmentLight struct {
	DepartmentId   uint64 `gorm:"column:department_id;primaryKey"`
	DepartmentName string `gorm:"column:department_name"`
}

func (DepartmentLight) TableName() string { return "security.department" }

// CityLight es una proyección de solo lectura de security.city.
type CityLight struct {
	CityId       uint64 `gorm:"column:city_id;primaryKey"`
	CityICode    string `gorm:"column:city_i_code"`
	CityName     string `gorm:"column:city_name"`
	DepartmentId uint64 `gorm:"column:department_id"`
}

func (CityLight) TableName() string { return "security.city" }

// CityICodeLight proyecta city_i_code, city_name y department_id para resolver ciudades en directorios.
type CityICodeLight struct {
	CityICode    string `gorm:"column:city_i_code"`
	CityName     string `gorm:"column:city_name"`
	DepartmentId uint64 `gorm:"column:department_id"`
}

func (CityICodeLight) TableName() string { return "security.city" }

// LocationOption es el formato { label, value } que consume stateOptionsPath en dinamic-form.
type LocationOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// CityLocationOption extiende LocationOption con el departamento para filtrado en el frontend.
type CityLocationOption struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	ICode        string `json:"iCode"`
	DepartmentId uint64 `json:"departmentId"`
}
