package db

type DBClientConfig struct {
	Hostname     string `json:"hostname"`
	Port         string `json:"port"`
	DatabaseName string `json:"database"`
	UserName     string `json:"user"`
	Password     string `json:"password"`
	SSLMode      string `json:"sslmode,omitempty"`

	// Credenciales exclusivas para el pool GORM.
	// Si están vacías, GORM usa UserName/Password igual que legacy.
	GormUser     string `json:"gorm_user,omitempty"`
	GormPassword string `json:"gorm_password,omitempty"`
}

// AsGormConfig retorna una copia de la config con las credenciales GORM.
func (c DBClientConfig) AsGormConfig() DBClientConfig {
	if c.GormUser != "" {
		c.UserName = c.GormUser
		c.Password = c.GormPassword
	}
	return c
}
