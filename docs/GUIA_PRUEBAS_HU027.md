# Guía de Pruebas — HU-027: Calendario de Seguimientos

Tutorial paso a paso para validar la funcionalidad del calendario de seguimientos en tu máquina local.

---

## 1. Instalaciones requeridas

Todos los comandos se ejecutan desde la carpeta `src/` del proyecto:

```bash
cd c:\Desarrollo\SOG_SALVIA\src
```

### 1.1 Instalar la CLI de Swagger (swag)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verifica que se instaló correctamente:

```bash
swag --version
```

Si el comando no se encuentra, agrega `%GOPATH%\bin` a tu variable de entorno `PATH`:
- `GOPATH` por defecto en Windows es `C:\Users\TU_USUARIO\go`
- Agrega `C:\Users\TU_USUARIO\go\bin` al PATH del sistema

### 1.2 Descargar dependencias del proyecto

```bash
go mod tidy
```

Esto descarga automáticamente todas las dependencias nuevas (swaggo, testify, GORM, etc.) que se agregaron al `go.mod`.

---

## 2. Generación y levantamiento de Swagger

### 2.1 Ejecutar la migración SQL

Antes de levantar el servidor, ejecuta la migración en tu base de datos PostgreSQL:

```sql
-- Conectarse a la BD configurada en src/config/db_config.json
-- y ejecutar el contenido de:
-- src/migrations/add_followup_v2_calendar_fields.sql

ALTER TABLE salvia.follow_up_v2
    ADD COLUMN IF NOT EXISTS team                VARCHAR(50),
    ADD COLUMN IF NOT EXISTS risk_status         VARCHAR(20),
    ADD COLUMN IF NOT EXISTS scheduled_date      DATE,
    ADD COLUMN IF NOT EXISTS is_completed        BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS sequence_number     INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS form_submission_id  VARCHAR(36);
```

Puedes ejecutarlo con `psql`, DBeaver, pgAdmin o cualquier cliente SQL.

### 2.2 Generar la documentación Swagger

Desde `src/`:

```bash
swag init -g main.go --parseDependency --parseInternal --output docs
```

Esto crea/actualiza la carpeta `src/docs/` con:
- `docs.go` — importado automáticamente por el código
- `swagger.json` — especificación OpenAPI
- `swagger.yaml` — misma especificación en YAML

### 2.3 Levantar el servidor

```bash
go run main.go
```

El servidor arranca en modo TLS local en el puerto 443.

### 2.4 Abrir Swagger UI

Abre tu navegador y ve a:

```
https://localhost/swagger/index.html
```

El navegador mostrará una advertencia de certificado autofirmado. Haz clic en "Avanzado" → "Continuar a localhost (no seguro)".

Deberías ver la interfaz de Swagger UI con todos los endpoints agrupados por tags. Busca el grupo "Seguimientos".

---

## 3. Guía de prueba manual (HU-027)

### 3.1 Generar el calendario de seguimientos

1. En Swagger UI, busca el tag "Seguimientos"
2. Haz clic en `POST /api/v1/cases/{victim_case_id}/follow-ups/generate`
3. Haz clic en "Try it out"
4. En el campo `victim_case_id` escribe un ID de caso existente en tu BD (o inventa uno para pruebas, ej: `caso-prueba-001`)
5. En el body, pega este JSON:

```json
{
  "risk_level": 3,
  "agent_id": "agente-prueba-001",
  "team": "Equipo Psicosocial"
}

```

Valores válidos para `risk_level`:
| Valor | Nivel    | Días S1→S4          |
|-------|----------|---------------------|
| 1     | Bajo     | +5, +15, +30, +60   |
| 2     | Moderado | +2, +15, +30, +45   |
| 3     | Alto     | +1, +3, +15, +30    |
| 4     | Extremo  | +1, +2, +3, +15     |

6. Haz clic en "Execute"
7. Respuesta esperada: `201 Created` con un array de 4 seguimientos, cada uno con:
   - `sequence_number` de 1 a 4
   - `status` = "PENDIENTE"
   - `scheduled_date` calculada desde hoy según la matriz
   - `risk_status` = "ALTO" (para risk_level=3)

### 3.2 Consultar el calendario generado

1. Haz clic en `GET /api/v1/cases/{victim_case_id}/follow-ups`
2. Haz clic en "Try it out"
3. Escribe el mismo `victim_case_id` que usaste en el paso anterior (ej: `caso-prueba-001`)
4. Haz clic en "Execute"
5. Respuesta esperada: `200 OK` con los 4 seguimientos ordenados por `scheduled_date` ASC

### 3.3 Probar la recalculación (cambio de riesgo)

1. Vuelve al `POST .../generate`
2. Usa el mismo `victim_case_id` pero cambia el `risk_level`:

```json
{
  "risk_level": 4,
  "agent_id": "agente-prueba-001",
  "team": "Equipo Psicosocial"
}
```

3. Respuesta esperada: `201 Created` con 4 nuevos seguimientos (los anteriores fueron reprogramados con soft-delete)
4. Verifica con el `GET` que ahora las fechas corresponden al nivel Extremo (+1, +2, +3, +15 días)

### 3.4 Probar sin cambio de riesgo

1. Vuelve al `POST .../generate` con el mismo `risk_level` que ya tiene el caso (4)
2. Respuesta esperada: `200 OK` o `201 Created` retornando los seguimientos existentes sin modificar (no se crean duplicados)

### 3.5 Probar validación de body inválido

1. Envía un body con `risk_level` fuera de rango:

```json
{
  "risk_level": 5,
  "agent_id": "agente-prueba-001"
}
```

2. Respuesta esperada: `400 Bad Request` con mensaje de error de validación

### 3.6 Probar caso sin seguimientos

1. Usa un `victim_case_id` que no tenga seguimientos (ej: `caso-inexistente-xyz`)
2. Llama al `GET .../follow-ups`
3. Respuesta esperada: `404 Not Found` con `{"error": "no se encontraron seguimientos para el caso"}`

### 3.7 Prueba con cURL (alternativa a Swagger UI)

```bash
# Generar calendario (risk_level=3, Alto)
curl -k -X POST https://localhost/api/v1/cases/caso-prueba-001/follow-ups/generate \
  -H "Content-Type: application/json" \
  -d '{"risk_level": 3, "agent_id": "agente-001", "team": "Equipo A"}'

# Consultar calendario
curl -k https://localhost/api/v1/cases/caso-prueba-001/follow-ups

# Buscar seguimiento por ID
curl -k "https://localhost/api/v1/cases/caso-prueba-001/follow-ups/by-id?id=UUID_DEL_SEGUIMIENTO"

# Listar paginado
curl -k "https://localhost/api/v1/cases/caso-prueba-001/follow-ups/list?page=0&limit=10"
```

El flag `-k` ignora la verificación del certificado autofirmado.

---

## 4. Pruebas unitarias automatizadas

### 4.1 Correr solo los tests de la HU-027

```bash
# Desde src/

# Tests del servicio (lógica de negocio)
go test ./salvia/service/ -run "TestGenerateOrRecalculate|TestGetByCaseID" -v

# Tests del controlador (endpoints HTTP)
go test ./salvia/controller/ -run "TestGetCalendar|TestGenerateCalendar|TestGetByIDHandler|TestListHandler" -v

# Tests del repositorio (mock de queries)
go test ./internal/repository/ -run "TestFindByCaseID|TestSoftDeletePending|TestBulkCreate" -v
```

### 4.2 Correr TODOS los tests del proyecto

```bash
go test ./... -v
```

### 4.3 Ver reporte de cobertura en HTML

```bash
# Generar el archivo de cobertura
go test ./salvia/... ./internal/repository/... -coverprofile=coverage_hu027.out

# Abrir el reporte visual en el navegador
go tool cover -html=coverage_hu027.out

# Ver resumen por función en la terminal
go tool cover -func=coverage_hu027.out
```

### 4.4 Output esperado (todos los tests pasan)

```
=== RUN   TestGetByIDHandler
=== RUN   TestGetByIDHandler/200_-_registro_encontrado
=== RUN   TestGetByIDHandler/400_-_falta_parámetro_id
=== RUN   TestGetByIDHandler/404_-_registro_no_encontrado
=== RUN   TestGetByIDHandler/500_-_error_de_infraestructura
--- PASS: TestGetByIDHandler
=== RUN   TestListHandler
--- PASS: TestListHandler
=== RUN   TestGetCalendar_200_WithResults
--- PASS: TestGetCalendar_200_WithResults
=== RUN   TestGetCalendar_404_Empty
--- PASS: TestGetCalendar_404_Empty
=== RUN   TestGenerateCalendar_201_Created
--- PASS: TestGenerateCalendar_201_Created
=== RUN   TestGenerateCalendar_400_InvalidRiskLevel
--- PASS: TestGenerateCalendar_400_InvalidRiskLevel
PASS
```

---

## Troubleshooting

| Problema | Solución |
|---|---|
| `swag: command not found` | Agrega `%GOPATH%\bin` al PATH del sistema y reinicia la terminal |
| `cannot find package "bitsflow/docs"` | Ejecuta `swag init -g main.go --parseDependency --parseInternal --output docs` |
| Swagger redirige al landing de SALVIA | Verifica que el `AuthMiddleware` en `MainRouter.go` excluye `/swagger` y `/api/v1` |
| `connection refused` al levantar | Verifica que `src/config/db_config.json` apunta a una BD PostgreSQL accesible |
| Tests fallan con `undefined: ginQueryInt` | Asegúrate de que `form_controller.go` está en el mismo package `controller` |
