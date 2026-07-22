// Package db provee la conexión GORM para el nuevo patrón de repositorios.
// Coexiste con la conexión pgx existente en common/db sin modificarla.
package db

import (
	commondb "bitsflow/common/db"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB construye una instancia *gorm.DB usando la misma configuración
// que el pool pgx existente (db_config.json vía common/db.DBClientConfig).
// El pool se limita a 12 conexiones para coexistir con el pool pgx de 80 (antes
// era 2 — subido para permitir concurrencia real en /admin/migrate/follow-up,
// que hasta ahora serializaba todo el paralelismo del script de migración detrás
// de solo 2 conexiones. Revisar si vuelve a bajarse tras terminar la migración
// de Bases Salvia, si el presupuesto total contra Supabase lo amerita).
// Incluye keepalive periódico para mantener la conexión activa con Supabase pooler.
func NewGormDB(cfg commondb.DBClientConfig) (*gorm.DB, error) {
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Hostname, cfg.Port, cfg.UserName, cfg.Password, cfg.DatabaseName, sslmode,
	)
	// PreferSimpleProtocol evita prepared statements en el driver pgx.
	// Sin esto, el pooler de Supabase puede devolver intermitentemente:
	// "prepared statement stmtcache_... does not exist" (SQLSTATE 26000).
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(12)
	sqlDB.SetMaxIdleConns(12)
	// No rotar conexiones — evita errores de DNS intermitentes con Supabase pooler
	// Si la conexión se cae, GORM la recrea automáticamente en el siguiente uso

	// Iniciar keepalive en background
	startKeepalive(sqlDB, 30*time.Second)

	return db, nil
}

// ─── Retry Helper ────────────────────────────────────────────────────────────

const (
	DefaultMaxRetries = 3
	DefaultRetryDelay = 600 * time.Millisecond
)

// WithRetry ejecuta una función que usa GORM y reintenta automáticamente
// si el error es de conexión/DNS. Útil para queries que pueden fallar
// por problemas intermitentes con el pooler de Supabase.
//
// Uso:
//
//	var results []MyStruct
//	err := db.WithRetry(func() error {
//	    return gormDB.Raw("SELECT ...").Scan(&results).Error
//	})
func WithRetry(fn func() error) error {
	return WithRetryN(fn, DefaultMaxRetries, DefaultRetryDelay)
}

// WithRetryN ejecuta fn con reintentos configurables.
func WithRetryN(fn func() error, maxRetries int, delay time.Duration) error {
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(delay * time.Duration(attempt))
			log.Printf("[GORM-RETRY] Reintentando (intento %d/%d)...", attempt+1, maxRetries)
		}
		err = fn()
		if err == nil || !IsConnectionError(err) {
			return err
		}
		log.Printf("[GORM-RETRY] Error de conexión detectado: %v", err)
	}
	return err
}

// IsConnectionError determina si un error es de conexión/DNS (retriable).
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	connectionErrors := []string{
		"hostname resolving error",
		"no such host",
		"connection refused",
		"connection reset",
		"broken pipe",
		"failed to connect",
		"i/o timeout",
		"network is unreachable",
		"connection timed out",
		"dial tcp",
		"eof",
		"prepared statement",
		"stmtcache_",
		"sqlstate 26000",
	}
	for _, ce := range connectionErrors {
		if strings.Contains(msg, ce) {
			return true
		}
	}
	return false
}

// ─── Keepalive ───────────────────────────────────────────────────────────────

// startKeepalive lanza una goroutine que hace Ping a la BD periódicamente
// para mantener la conexión activa y evitar que el pooler la cierre por inactividad.
func startKeepalive(sqlDB *sql.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := sqlDB.PingContext(ctx); err != nil {
				log.Printf("[GORM-KEEPALIVE] Ping falló: %v", err)
			}
			cancel()
		}
	}()
}
