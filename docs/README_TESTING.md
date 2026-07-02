# Guía de Testing — SALVIA Fase 2

## A) Ejecutar los tests

```bash
# Desde src/ (donde está go.mod)

# Todos los tests del módulo salvia con output detallado
go test ./salvia/... -v

# Solo el servicio
go test ./salvia/service/... -v

# Solo el controlador
go test ./salvia/controller/... -v
```

## B) Ver el coverage

```bash
# Generar el reporte de cobertura
go test ./salvia/... -coverprofile=coverage.out

# Abrir el reporte visual en el navegador
go tool cover -html=coverage.out

# Ver resumen por función en la terminal
go tool cover -func=coverage.out
```

## C) Por qué usamos Mocks y no PostgreSQL

En los tests de **servicio** y **controlador** no necesitamos una base de datos real. Aquí está el razonamiento:

| Criterio | Mock | PostgreSQL real |
|---|---|---|
| Velocidad | Milisegundos | Segundos (conexión + query) |
| Aislamiento | Total: cada test controla exactamente qué retorna el repo/servicio | Depende del estado de la BD |
| Reproducibilidad | Siempre igual, sin datos residuales | Requiere fixtures y limpieza |
| CI/CD | Sin dependencias externas | Requiere Docker o instancia dedicada |
| Qué se prueba | Lógica de negocio y manejo de errores HTTP | Queries SQL y mapeo GORM |

**Regla práctica del proyecto:**

- `service_test.go` → mockea el **repositorio**: verifica que el servicio traduce errores correctamente (`gorm.ErrRecordNotFound` → `ErrFollowUpNotFound`) y delega bien al repo.
- `controller_test.go` → mockea el **servicio**: verifica que el controlador responde con el código HTTP correcto según el resultado del servicio.
- Tests de **integración** (con PostgreSQL real vía `testcontainers-go`) se reservan para el paquete `repository`, donde sí importa que el SQL generado por GORM sea correcto.

Esta separación sigue el principio de **probar una sola capa a la vez**, haciendo los tests rápidos, deterministas y fáciles de mantener.
