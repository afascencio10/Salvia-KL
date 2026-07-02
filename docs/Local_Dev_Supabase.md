# Configuración de DB — Local vs Producción

Guía para configurar la conexión a base de datos según el entorno.

---

## 🖥️ Desarrollo local (Supabase free plan)

### Arrancar el servidor

```bash
cd src
PORT=9090 go run main.go
```

### 1. `src/config/db_config.json`

```json
{
    "hostname": "aws-1-us-west-1.pooler.supabase.com",
    "port": "5432",
    "database": "postgres",
    "user": "postgres.pwwelwfhauspznpatuqm",
    "password": "Salvia2026@",
    "sslmode": "require"
}
```

### 2. `src/salvia/facades/MainRouter.go`

```go
var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: 3}
```

### 3. `src/security/facades/MainRouter.go`

```go
var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: 3}
```

### 4. `src/internal/db/gorm_connection.go`

```go
sqlDB.SetMaxOpenConns(2)
sqlDB.SetMaxIdleConns(1)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
sqlDB.SetConnMaxIdleTime(4 * time.Minute)
```

---

## 🚀 Producción

### 1. `src/config/db_config.json`

```json
{
    "hostname": "<HOST_PRODUCCION>",
    "port": "5432",
    "database": "<DB_PRODUCCION>",
    "user": "<USUARIO_PRODUCCION>",
    "password": "<PASSWORD_PRODUCCION>",
    "sslmode": "require"
}
```

### 2. `src/salvia/facades/MainRouter.go`

```go
var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: 80}
```

### 3. `src/security/facades/MainRouter.go`

```go
var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: 80}
```

### 4. `src/internal/db/gorm_connection.go`

```go
sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(5)
```

---

## Resumen de conexiones

| Componente              | Local (Supabase free) | Producción |
|-------------------------|-----------------------|------------|
| Legacy pgx — salvia     | 3                     | 80         |
| Legacy pgx — security   | 3                     | 80         |
| GORM                    | 2                     | 10         |
| **Total**               | **8**                 | **170**    |

---

## ¿Por qué se limitan las conexiones en local?

El **Session Pooler de Supabase en plan free** tiene un límite de ~10 clientes simultáneos.
Con los valores de producción (170 conexiones) el servidor falla al arrancar con:

```
FATAL: MaxClientsInSessionMode: max clients reached
```

En producción se recomienda usar una instancia PostgreSQL sin restricción de clientes
(instancia dedicada, Supabase Pro, o Render).
