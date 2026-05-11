package db

// DefaultPoolSize define el tamaño del pool pgx compartido entre módulos.
// Ajustar aquí afecta salvia y security simultáneamente.
// Desarrollo local (Supabase free): 2 — Producción: 80
const DefaultPoolSize uint = 2

type DBServerConfig struct {
	PoolSize uint `json:"size"`
}
