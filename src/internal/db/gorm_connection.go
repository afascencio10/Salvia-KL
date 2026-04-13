// Package db provee la conexión GORM para el nuevo patrón de repositorios.
// Coexiste con la conexión pgx existente en common/db sin modificarla.
package db

import (
	commondb "bitsflow/common/db"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB construye una instancia *gorm.DB usando la misma configuración
// que el pool pgx existente (db_config.json vía common/db.DBClientConfig).
// El pool se limita a 5 conexiones para coexistir con el pool pgx de 80.
func NewGormDB(cfg commondb.DBClientConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.Hostname, cfg.Port, cfg.UserName, cfg.Password, cfg.DatabaseName,
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
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	return db, nil
}
