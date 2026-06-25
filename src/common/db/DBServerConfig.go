package db

// DefaultPoolSize define el tamaño del pool pgx compartido entre módulos.
// Ajustar aquí afecta salvia y security simultáneamente.
// Desarrollo local (Supabase free): 2 — Producción: 80
// Subido a 10 para evitar serialización de queries del dashboard KPIs.
const DefaultPoolSize uint = 10

type DBServerConfig struct {
	PoolSize uint `json:"size"`
}
