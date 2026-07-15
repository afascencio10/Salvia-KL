---
okf_version: "1.0"
type: Environment_Config
title: "Matriz de Entornos"
description: "Configuración y diferencias entre los entornos develop, testing/QA y producción de SALVIA, incluyendo base de datos, despliegue e inyección de secretos."
owner: "@sre-squad"
status: active
tags: [ops, environments, config, infrastructure]
last_updated: "2026-07-02"
---

# Matriz de Entornos

## Inventario de Entornos

| Entorno | Propósito | Base de Datos | Despliegue |
| :--- | :--- | :--- | :--- |
| `develop` | Desarrollo local y pruebas | **Supabase** PostgreSQL (`aws-1-us-west-1.pooler.supabase.com`) | Local (`go run`/`go build`) |
| `testing` (QA) | Pre-producción, QA | PostgreSQL definido en `secrets.DB_CONFIG_QA` | GitHub Actions → runner `[self-hosted, windows, qa]` (push a `testing`) |
| `production` | Usuarios reales | PostgreSQL **self-hosted (localhost)** definido en `secrets.DB_CONFIG_PRD` | GitHub Actions → runner `[self-hosted, windows, prd]` (push a `main`) |

## Configuración de Base de Datos

SALVIA lee la conexión desde `src/config/db_config.json`. En `develop` este archivo apunta a Supabase (versionado en el repo). En `testing` y `production`, el workflow de despliegue **sobrescribe** ese archivo con el contenido del secret correspondiente antes de compilar.

| Aspecto | develop | testing (QA) | production |
| :--- | :--- | :--- | :--- |
| Origen de `db_config.json` | Versionado en el repo | `secrets.DB_CONFIG_QA` | `secrets.DB_CONFIG_PRD` |
| Host | Supabase pooler (AWS us-west-1) | Definido en el secret | `localhost` (self-hosted) |
| `sslmode` | `require` | Según secret | Según secret |
| Usuarios | `salvia_legacy` (pgx) / `salvia_gorm` (GORM) | Según secret | Según secret |

> [!WARNING]
> `src/config/db_config.json` con credenciales de Supabase está **versionado en el repositorio**. Considerar moverlo también a un secret / archivo ignorado por git y rotar esas credenciales (ver [MEJORAS_SEGURIDAD_URGENTES](/docs/MEJORAS_SEGURIDAD_URGENTES.md)).

## Despliegue por Entorno

| Paso | testing (QA) | production |
| :--- | :--- | :--- |
| Trigger | `push` a rama `testing` | `push` a rama `main` |
| Runner | `[self-hosted, windows, qa]` | `[self-hosted, windows, prd]` |
| Escribir config | `db_config.json` ← `secrets.DB_CONFIG_QA` | `db_config.json` ← `secrets.DB_CONFIG_PRD` |
| Build | `cd src && go build -o ../salvia.exe main.go` | igual |
| Deploy | Detener servicio `SalviaApp`, backup del `.exe` anterior, copiar nuevo binario a `C:\salvia`, reiniciar servicio | igual |

```text
   develop (local)        testing / QA            production
  ┌───────────────┐      ┌───────────────┐      ┌───────────────┐
  │ Supabase PG   │      │ push: testing  │      │ push: main     │
  │ go run/build  │─────▶│ runner win/qa  │─────▶│ runner win/prd │
  │               │      │ DB_CONFIG_QA   │      │ DB_CONFIG_PRD  │
  └───────────────┘      └───────────────┘      └───────────────┘
                          servicio SalviaApp      servicio SalviaApp
```

## Servicio y Serving

| Aspecto | Detalle |
| :--- | :--- |
| Proceso | Servicio Windows `SalviaApp` ejecutando `C:\salvia\salvia.exe` |
| TLS | HTTPS con `certs/salvia.crt` / `certs/salvia.key`; si se define la env `PORT`, corre en HTTP por ese puerto |
| Frontend | Embebido en el binario (`go:embed`) — un redeploy actualiza UI y API a la vez |
| Backup | El deploy renombra el binario anterior a `salvia_old.exe` antes de copiar el nuevo |

> [!NOTE]
> No existe (aún) un entorno `staging` intermedio ni réplicas/auto-scaling: cada entorno es un único nodo. La estrategia de rollback documentada en [deployment-strategy.md](/.knowledge/4-operations/deployment-strategy.md) es genérica y **pendiente de adaptar** al swap de servicio Windows real.
