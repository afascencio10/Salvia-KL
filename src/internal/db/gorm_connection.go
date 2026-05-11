// Package db provee la conexión GORM para el nuevo patrón de repositorios.
// Coexiste con la conexión pgx existente en common/db sin modificarla.
package db

import (
	commondb "bitsflow/common/db"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB construye una instancia *gorm.DB usando la misma configuración
// que el pool pgx existente (db_config.json vía common/db.DBClientConfig).
// El pool se limita a 5 conexiones para coexistir con el pool pgx de 80.
func NewGormDB(cfg commondb.DBClientConfig) (*gorm.DB, error) {
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Hostname, cfg.Port, cfg.UserName, cfg.Password, cfg.DatabaseName, sslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(7)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // rota conexiones cada 5 min
	sqlDB.SetConnMaxIdleTime(4 * time.Minute) // cierra conexiones idle tras 4 min
	return db, nil
}
