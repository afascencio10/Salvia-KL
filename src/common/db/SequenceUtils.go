package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// ── Helpers de error ──────────────────────────────────────────────────────────

// IsDuplicateKeyError retorna true si err es un error PostgreSQL de clave
// duplicada (código 23505 — unique_violation).
func IsDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// ── Conexión auxiliar ─────────────────────────────────────────────────────────

func freshConn(clientConfig *DBClientConfig) (*pgx.Conn, error) {
	sslmode := clientConfig.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		clientConfig.UserName, clientConfig.Password,
		clientConfig.Hostname, clientConfig.Port,
		clientConfig.DatabaseName, sslmode)
	return pgx.Connect(context.Background(), dsn)
}

// ── Sincronización individual ─────────────────────────────────────────────────

// SyncSequence corrige una secuencia de PostgreSQL que quedó desincronizada
// respecto a los datos reales de la tabla. Abre una conexión directa (fuera
// del pool) para poder ejecutar el setval incluso si la conexión del pool está
// en una transacción abortada.
//
// Parámetros:
//   - seqFQN:   nombre completamente calificado de la secuencia, ej: "salvia.victim_case_form2_victim_case_form2_id_seq"
//   - tableFQN: nombre completamente calificado de la tabla,    ej: "salvia.victim_case_form2"
//   - idCol:    nombre de la columna PK serial,                 ej: "victim_case_form2_id"
func SyncSequence(clientConfig *DBClientConfig, seqFQN, tableFQN, idCol string) error {
	conn, err := freshConn(clientConfig)
	if err != nil {
		return fmt.Errorf("SyncSequence: no se pudo conectar: %w", err)
	}
	defer conn.Close(context.Background())

	query := fmt.Sprintf(
		"SELECT setval('%s', COALESCE((SELECT MAX(%s) FROM %s), 1))",
		seqFQN, idCol, tableFQN,
	)
	if _, err = conn.Exec(context.Background(), query); err != nil {
		return fmt.Errorf("SyncSequence: error al resetear secuencia '%s': %w", seqFQN, err)
	}

	fmt.Printf("[SyncSequence] Secuencia '%s' resincronizada.\n", seqFQN)
	return nil
}

// ── Sincronización masiva ─────────────────────────────────────────────────────

// SyncAllSequences detecta y corrige todas las secuencias desincronizadas del
// schema indicado. Usa pg_sequences + pg_depend para descubrir automáticamente
// qué tabla y columna corresponde a cada secuencia, sin necesidad de listarlas
// manualmente.
//
// Llamar una vez al arrancar el servidor garantiza que ninguna secuencia quede
// detrás del MAX(id) real, evitando errores de duplicate key al insertar.
func SyncAllSequences(clientConfig *DBClientConfig, schema string) error {
	conn, err := freshConn(clientConfig)
	if err != nil {
		return fmt.Errorf("SyncAllSequences: no se pudo conectar: %w", err)
	}
	defer conn.Close(context.Background())

	// Descubre todas las secuencias del schema con su tabla y columna asociadas.
	discoverSQL := `
		SELECT
			s.schemaname || '.' || s.sequencename   AS seq_fqn,
			s.schemaname || '.' || t.relname         AS table_fqn,
			a.attname                                AS col_name
		FROM pg_sequences s
		JOIN pg_class sc
			ON sc.relname        = s.sequencename
			AND sc.relnamespace  = (SELECT oid FROM pg_namespace WHERE nspname = s.schemaname)
		JOIN pg_depend d
			ON d.objid    = sc.oid
			AND d.deptype = 'a'
		JOIN pg_class t
			ON t.oid = d.refobjid
		JOIN pg_attribute a
			ON a.attrelid = t.oid
			AND a.attnum  = d.refobjsubid
		WHERE s.schemaname = $1
		ORDER BY seq_fqn`

	rows, err := conn.Query(context.Background(), discoverSQL, schema)
	if err != nil {
		return fmt.Errorf("SyncAllSequences: error consultando secuencias: %w", err)
	}

	type seqInfo struct{ seqFQN, tableFQN, colName string }
	var sequences []seqInfo
	for rows.Next() {
		var si seqInfo
		if err := rows.Scan(&si.seqFQN, &si.tableFQN, &si.colName); err != nil {
			rows.Close()
			return fmt.Errorf("SyncAllSequences: error leyendo fila: %w", err)
		}
		sequences = append(sequences, si)
	}
	rows.Close() // libera la conexión antes de ejecutar los setval

	if len(sequences) == 0 {
		fmt.Printf("[SyncAllSequences] No se encontraron secuencias en el schema '%s'.\n", schema)
		return nil
	}

	fixed := 0
	var lastErr error
	for _, si := range sequences {
		q := fmt.Sprintf(
			"SELECT setval('%s', COALESCE((SELECT MAX(%s) FROM %s), 1))",
			si.seqFQN, si.colName, si.tableFQN,
		)
		if _, err := conn.Exec(context.Background(), q); err != nil {
			fmt.Printf("[SyncAllSequences] Error en '%s': %v\n", si.seqFQN, err)
			lastErr = err
			continue
		}
		fixed++
	}

	fmt.Printf("[SyncAllSequences] %d/%d secuencias del schema '%s' sincronizadas.\n",
		fixed, len(sequences), schema)
	return lastErr
}

// ── Retry automático ──────────────────────────────────────────────────────────

// WithSequenceRetry ejecuta fn y, si falla con duplicate key (secuencia
// desincronizada), sincroniza todas las secuencias del schema y reintenta fn
// una sola vez. Si el segundo intento también falla, retorna ese error.
//
// Uso típico en un controlador:
//
//	err = db.WithSequenceRetry(&clientConfig, "salvia", func() error {
//	    return salvia_daos.SetVictimCaseForm2(form, connData, &clientConfig, serverConfig)
//	})
func WithSequenceRetry(clientConfig *DBClientConfig, schema string, fn func() error) error {
	err := fn()
	if err == nil || !IsDuplicateKeyError(err) {
		return err
	}

	fmt.Printf("[WithSequenceRetry] Duplicate key en schema '%s' — resincronizando secuencias...\n", schema)
	if syncErr := SyncAllSequences(clientConfig, schema); syncErr != nil {
		// Loguea el error de sync pero retorna el error original
		fmt.Printf("[WithSequenceRetry] Error al resincronizar: %v\n", syncErr)
		return err
	}

	fmt.Println("[WithSequenceRetry] Reintentando operación...")
	return fn()
}
