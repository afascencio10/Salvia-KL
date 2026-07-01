━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 EVENTO: Migración de schema — tabla dupla y ajuste psychosocial_support
   Tipo: Infrastructure / Backend
   Código: M-01
   Prerequisito de: E-01, E-11 y todos los flujos que lean los campos nuevos
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **Alcance:** Documenta los cambios de modelo y BD que deben existir **antes**
> de implementar el componente. El campo `submitted_by` se define en schema;
> este componente solo lo **lee** para mostrar nombre y equipo del remitente.

INPUT: {
  gormDB:   conexión GORM al arrancar la aplicación   → src/main.go
  models:   structs Go a migrar                      → src/internal/models/
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 1 — Crear modelo Dupla
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/internal/models/dupla.go`

```go
// Dupla representa un par psicóloga + trabajador social del equipo de Atención Psicosocial.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.dupla (
//	    id                 VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    name               VARCHAR(36) NOT NULL,
//	    psychologist_id    VARCHAR(36) NOT NULL,
//	    social_worker_id   VARCHAR(36) NOT NULL,
//	    created_at         TIMESTAMPTZ,
//	    updated_at         TIMESTAMPTZ,
//	    deleted_at         TIMESTAMPTZ
//	);
type Dupla struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"type:varchar(36);not null"                            json:"name"`
	PsychologistID  string         `gorm:"type:varchar(36);not null;column:psychologist_id"     json:"psychologistId"`
	SocialWorkerID  string         `gorm:"type:varchar(36);not null;column:social_worker_id"    json:"socialWorkerId"`
	CreatedAt       time.Time      `                                                             json:"createdAt"`
	UpdatedAt       time.Time      `                                                             json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                                                 json:"deletedAt,omitempty"`
}

func (Dupla) TableName() string { return "salvia.dupla" }
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 2 — Actualizar modelo PsychosocialSupport
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/internal/models/psychosocial_support.go`

Agregar campos y constantes de estado (mismo patrón que `entity_letter.go`):

```go
// Flujo de estados:
//
//	abierto       → al crear la remisión
//	en_gestion    → se logró primer contacto
//	en_devolucion → no cumplió criterios
//	cerrado       → por cualquier motivo
//
// SQL columnas nuevas en salvia.psychosocial_support:
//   submitted_by  VARCHAR(36)  — general_user_i_code del agente remitente (lectura en UI)
//   dupla_id      VARCHAR(36)  — FK salvia.dupla.id
//   agent_id      VARCHAR(36)  — general_user_i_code del profesional asignado
```

Campos nuevos en el struct:

| Campo Go | Columna BD | Tipo | Nullable |
|---|---|---|---|
| `SubmittedBy` | `submitted_by` | varchar(36) | Sí |
| `DuplaID` | `dupla_id` | varchar(36) | Sí |
| `AgentID` | `agent_id` | varchar(36) | Sí |

Constantes de estado:

```go
const (
	PsychosocialSupportStatusAbierto      = "abierto"
	PsychosocialSupportStatusEnGestion    = "en_gestion"
	PsychosocialSupportStatusEnDevolucion = "en_devolucion"
	PsychosocialSupportStatusCerrado      = "cerrado"
)
```

Cambios adicionales en el struct existente:

| Campo | Cambio |
|---|---|
| `Status` | `default:'abierto'` (reemplaza `'ACTIVE'`) |
| `Status` gorm tag | `type:varchar(30)` (ampliar de varchar(20)) |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 3 — Registrar AutoMigrate en main.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/main.go`

Agregar al slice de AutoMigrate:

```go
&models.Dupla{},
&models.PsychosocialSupport{},  // ya existe — AutoMigrate agrega columnas nuevas
```

Orden: `Dupla` **antes** de `PsychosocialSupport` si se agrega FK explícita en el futuro.
Por ahora la FK es lógica (varchar), no constraint en BD.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 4 — Migración de datos legacy (status ACTIVE)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Registros existentes pueden tener `status = 'ACTIVE'` (valor actual al crear remisión).

Ejecutar una sola vez (SQL o script de migración):

```sql
UPDATE salvia.psychosocial_support
SET    status = 'abierto'
WHERE  status = 'ACTIVE';
```

No modificar `submitted_by` en registros históricos desde esta migración.
El componente de listado lo leerá cuando esté poblado.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 5 — Actualizar Create en psychosocial_support_service
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/salvia/service/psychosocial_support_service.go`

```go
if ps.Status == "" {
    ps.Status = models.PsychosocialSupportStatusAbierto  // "abierto"
}
```

Reemplazar el default anterior `"ACTIVE"`.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 6 — Actualizar form_service.go (creación desde seguimiento)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/salvia/service/form_service.go`

Al crear la remisión en `case "atencion_psico"`:

```go
Status: models.PsychosocialSupportStatusAbierto,  // "abierto" en lugar de "ACTIVE"
// submitted_by, dupla_id, agent_id → NO se setean aquí (fuera de alcance)
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 7 — Verificación post-migración
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Checklist:

- [ ] Tabla `salvia.dupla` existe con columnas esperadas
- [ ] Columnas `submitted_by`, `dupla_id`, `agent_id` existen en `psychosocial_support`
- [ ] Default de `status` es `'abierto'`
- [ ] Registros legacy `ACTIVE` migrados a `abierto`
- [ ] Constantes Go exportadas y usadas en servicio de creación

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ARCHIVOS AFECTADOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Archivo | Acción |
|---|---|
| `src/internal/models/dupla.go` | Crear |
| `src/internal/models/psychosocial_support.go` | Modificar |
| `src/main.go` | Agregar `&models.Dupla{}` a AutoMigrate |
| `src/salvia/service/psychosocial_support_service.go` | Default status |
| `src/salvia/service/form_service.go` | Status al crear remisión |
| Script SQL one-shot | Migrar `ACTIVE` → `abierto` |
