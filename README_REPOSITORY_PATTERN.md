# Generic Repository Pattern — SALVIA Backend

Guía de onboarding para la nueva capa de repositorio genérico con Go Generics y GORM.

---

## ¿Qué se creó?

| Archivo | Descripción |
|---|---|
| `src/internal/repository/base_repository.go` | Interfaz `Repository[T any]` e implementación genérica con GORM |
| `src/internal/repository/followup_repository.go` | Repositorio específico `FollowUpRepository` con método de negocio `FindByCaseID` |
| `src/internal/models/followup_v2.go` | Modelo GORM `FollowUpV2` para la tabla `follow_up_v2` |

> **Importante**: Esta capa es **aditiva**. No modifica ni reemplaza el patrón DAO existente en `common/dao`, `salvia/dao` y `security/dao`. Ambas capas coexisten.

---

## ¿Cómo se estructura el patrón?

```
src/internal/
├── models/
│   └── followup_v2.go        ← struct GORM (tabla + columnas)
└── repository/
    ├── base_repository.go    ← Repository[T any] genérico
    └── followup_repository.go ← FollowUpRepository específico
```

### Jerarquía de tipos

```
Repository[T any]          ← interfaz genérica (CRUD + paginación)
    └── repository[T any]  ← implementación con *gorm.DB

FollowUpRepository         ← interfaz específica (compone Repository[FollowUpV2])
    └── followUpRepository ← embebe repository[FollowUpV2] + agrega FindByCaseID
```

### Métodos disponibles en todo repositorio

| Método | Descripción |
|---|---|
| `Create(ctx, *T) error` | Inserta un registro; GORM escanea el ID generado |
| `Update(ctx, *T) error` | Actualiza todos los campos (requiere PK) |
| `Delete(ctx, id) error` | Soft-delete por ID |
| `FindByID(ctx, id) (*T, error)` | Busca por PK; retorna `gorm.ErrRecordNotFound` si no existe |
| `FindWithPagination(ctx, page, pageSize) (PageResult[T], error)` | Paginación offset-based, base-0 |

---

## Guía paso a paso: agregar una nueva tabla

### Paso 1 — Instalar dependencias (una sola vez)

```bash
# Desde src/ (donde está go.mod)
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
```

Actualizar la directiva de versión en `go.mod` (los generics requieren Go 1.18+):

```
go 1.21
```

### Paso 2 — Crear el modelo GORM

Crear `src/internal/models/mi_entidad.go`:

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type MiEntidad struct {
    ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
    Nombre    string         `gorm:"type:varchar(100);not null"`
    // ... más campos
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**Reglas de naming GORM**: el struct `MiEntidad` mapea a la tabla `mi_entidades` por convención. Para usar otro nombre:

```go
func (MiEntidad) TableName() string { return "salvia.mi_entidad" }
```

### Paso 3 — Crear el repositorio específico

Crear `src/internal/repository/mi_entidad_repository.go`:

```go
package repository

import (
    "bitsflow/internal/models"
    "context"
    "gorm.io/gorm"
)

// MiEntidadRepository extiende el CRUD genérico con métodos de negocio.
type MiEntidadRepository interface {
    Repository[models.MiEntidad]
    // Agrega aquí los métodos específicos de tu dominio:
    FindByNombre(ctx context.Context, nombre string) ([]models.MiEntidad, error)
}

type miEntidadRepository struct {
    repository[models.MiEntidad]
    db *gorm.DB
}

func NewMiEntidadRepository(db *gorm.DB) MiEntidadRepository {
    return &miEntidadRepository{
        repository: repository[models.MiEntidad]{db: db},
        db:         db,
    }
}

func (r *miEntidadRepository) FindByNombre(ctx context.Context, nombre string) ([]models.MiEntidad, error) {
    var items []models.MiEntidad
    result := r.db.WithContext(ctx).Where("nombre = ?", nombre).Find(&items)
    return items, result.Error
}
```

### Paso 4 — Inicializar GORM y el repositorio

En el punto de entrada de tu módulo (e.g., `main.go` o un archivo de configuración):

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "bitsflow/internal/repository"
)

func initGORM() (*gorm.DB, error) {
    dsn := "host=localhost user=salvia password=secret dbname=salvia_db port=5432 sslmode=require"
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// En tu función de setup:
gormDB, err := initGORM()
if err != nil {
    log.Fatal(err)
}

// Ajustar el pool para alinearse con el pool pgx existente (80 conexiones)
sqlDB, _ := gormDB.DB()
sqlDB.SetMaxOpenConns(80)
sqlDB.SetMaxIdleConns(10)

miRepo := repository.NewMiEntidadRepository(gormDB)
```

### Paso 5 — Usar el repositorio en un handler Gin

```go
func GetMiEntidadHandler(repo repository.MiEntidadRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        entity, err := repo.FindByID(c.Request.Context(), id)
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, gin.H{"error": "no encontrado"})
            return
        }
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, entity)
    }
}
```

### Paso 6 — Migración de la tabla (opcional, auto-migrate)

```go
// Solo en desarrollo; en producción usar migraciones SQL explícitas.
gormDB.AutoMigrate(&models.MiEntidad{})
```

---

## Diferencias con el patrón DAO existente

| Aspecto | DAO existente (pgx/v4) | Repositorio genérico (GORM) |
|---|---|---|
| Boilerplate por entidad | Alto (SQL manual, scan manual) | Mínimo (solo modelo + interfaz específica) |
| Control del SQL | Total | Parcial (GORM genera el SQL) |
| Soft-delete | Manual | Automático con `gorm.DeletedAt` |
| Paginación | Manual (`GetOffsetQuery`) | Incluida en `FindWithPagination` |
| Generics | No | Sí (`[T any]`) |
| Casos de uso ideales | Queries complejas, joins, lógica existente | CRUD estándar, nuevas entidades |

---

## Preguntas frecuentes

**¿Puedo usar ambas capas en el mismo handler?**
Sí. Son completamente independientes. Puedes llamar a un DAO existente y a un repositorio GORM en el mismo flujo.

**¿Cómo manejo transacciones con el repositorio genérico?**
Pasa una instancia de `*gorm.DB` con transacción activa:
```go
tx := gormDB.Begin()
repo := repository.NewFollowUpRepository(tx)
// ... operaciones
tx.Commit()
```

**¿El soft-delete afecta las queries automáticamente?**
Sí. GORM agrega `WHERE deleted_at IS NULL` en todas las queries por defecto cuando el modelo tiene `gorm.DeletedAt`. Para incluir registros eliminados usa `db.Unscoped()`.

---

## Transacciones Complejas: Patrón Unit of Work

Para operaciones que involucran múltiples tablas que deben persistirse de forma atómica, el repositorio específico puede encapsular la transacción internamente. El método `SubmitKoboForm` en `SubmissionRepository` es el ejemplo canónico.

### ¿Qué hace `SubmitKoboForm`?

Persiste tres entidades en una sola transacción, garantizando que o todo se guarda o nada se guarda:

```
CaseTimeline  →  FormSubmission  →  []Answer
   (ancla           (vinculado          (vinculadas
   forense)         al timeline)        al submission)
```

### Flujo interno

```go
tx := r.db.WithContext(ctx).Begin()

// 1. Evento de auditoría (append-only, sin DeletedAt)
tx.Create(timelineEvent)

// 2. Submission anclado al timeline
submission.TimelineID = timelineEvent.ID
tx.Create(submission)

// 3. Respuestas en batch insert
for i := range answers { answers[i].FormSubmissionID = submission.ID }
tx.Create(&answers)

tx.Commit()
```

Si cualquier paso falla, se ejecuta `tx.Rollback()` antes de retornar el error, evitando locks y datos parciales en la BD.

### Cómo replicar este patrón para otras operaciones multi-tabla

1. Recibe `*gorm.DB` por inyección en el constructor del repositorio.
2. Abre la transacción con `db.WithContext(ctx).Begin()`.
3. Agrega `defer func() { if recover() != nil { tx.Rollback() } }()` para cubrir panics.
4. Ejecuta cada operación con `tx.Create/Save/Delete`; ante cualquier error llama `tx.Rollback()` y retorna.
5. Cierra con `tx.Commit().Error`.

> El `CaseTimeline` es **append-only** (sin `DeletedAt`). Nunca se elimina ni actualiza: es el registro forense inmutable de lo que ocurrió.
