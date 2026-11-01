#!/bin/bash
# Levanta una BD Postgres local dentro del sandbox cloud de Claude Code,
# donde no hay salida de red hacia la BD de pruebas compartida (Supabase).
# Ver cloudDevelopment/README.md.
#
# No-op fuera del sandbox cloud: solo actúa si CLAUDE_CODE_REMOTE=true
# (variable que exporta automáticamente el sandbox). En local, run_dev.sh
# sigue usando config/db_config.json exactamente como antes — este script
# no lo toca.
#
# Uso: se fuentea (source) desde run_dev.sh, no se ejecuta aparte — así los
# `export` de abajo quedan disponibles para el `go run` que sigue.

if [ "$CLAUDE_CODE_REMOTE" != "true" ]; then
    return 0 2>/dev/null || exit 0
fi

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PG_BASE="${SALVIA_LOCAL_PG_DIR:-/var/lib/salvia-local-pg}"
PGDATA="$PG_BASE/data"
PG_LOG="$PG_BASE/log.txt"
SCHEMA_MARKER="$PG_BASE/.schema_applied"

export DB_HOST="127.0.0.1"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="salvia_local"

# ─── 1. Binarios de Postgres (cache: solo se instalan la primera vez) ──────
if ! command -v pg_ctl >/dev/null 2>&1; then
    echo "[cloudDevelopment] Instalando postgresql (una sola vez, se cachea en la imagen del sandbox)..."
    apt-get update -qq && apt-get install -y -qq postgresql >/dev/null
fi
PG_BIN="$(dirname "$(command -v pg_ctl)")"

# ─── 2. Inicializar + arrancar el cluster local (idempotente) ─────────────
mkdir -p "$PG_BASE"
export LC_ALL=C LANG=C   # evita locale issues al arrancar postgres

if [ ! -f "$PGDATA/PG_VERSION" ]; then
    echo "[cloudDevelopment] Inicializando BD Postgres local en $PGDATA..."
    "$PG_BIN/initdb" -D "$PGDATA" -U postgres --auth=trust -E UTF8 --locale=C >/dev/null
fi

if ! "$PG_BIN/pg_ctl" -D "$PGDATA" status >/dev/null 2>&1; then
    echo "[cloudDevelopment] Arrancando Postgres local en el puerto $DB_PORT..."
    "$PG_BIN/pg_ctl" -D "$PGDATA" -o "-p $DB_PORT -k /tmp -h 127.0.0.1" -l "$PG_LOG" start
fi

# Espera a que acepte conexiones (hasta ~10s)
for i in $(seq 1 20); do
    if PGPASSWORD="$DB_PASSWORD" "$PG_BIN/psql" -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "SELECT 1" >/dev/null 2>&1; then
        break
    fi
    sleep 0.5
done

# ─── 3. Esquema + seed (solo la primera vez por sandbox) ──────────────────
if [ ! -f "$SCHEMA_MARKER" ]; then
    echo "[cloudDevelopment] Creando esquema (security + salvia) y datos de prueba..."

    PSQL="$PG_BIN/psql -h $DB_HOST -p $DB_PORT -U $DB_USER"

    $PSQL -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 \
        || $PSQL -d postgres -c "CREATE DATABASE $DB_NAME;"

    $PSQL -d "$DB_NAME" -c "CREATE SCHEMA IF NOT EXISTS security; CREATE SCHEMA IF NOT EXISTS salvia;
        CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"; CREATE EXTENSION IF NOT EXISTS pgcrypto; CREATE EXTENSION IF NOT EXISTS unaccent;"

    # ON_ERROR_STOP=0: estos dumps legacy traen un puñado de ALTER TABLE /
    # GRANT ya conocidos y no destructivos que fallan en una BD nueva — ver
    # cloudDevelopment/README.md ("Errores esperados al aplicar el esquema").
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=0 -f "$REPO_ROOT/src/security/scripts/db_postgres_creation.sql" >/dev/null 2>&1 || true
    # Se aplica DOS veces: hay ALTER TABLE (FKs) que referencian tablas
    # definidas más abajo en el mismo archivo salvia — en la 2da pasada ya
    # existen. Ver README.
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=0 -f "$REPO_ROOT/src/salvia/scripts/db_postgres_creation.sql" >/dev/null 2>&1 || true
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=0 -f "$REPO_ROOT/src/salvia/scripts/db_postgres_creation.sql" >/dev/null 2>&1 || true

    # Columna legacy que normalmente agrega la propia app al arrancar
    # (migrateLegacyTables en main.go) — se adelanta aquí para que el seed
    # pueda asignarle equipo al usuario de prueba rol 'ro' desde el inicio.
    $PSQL -d "$DB_NAME" -c "ALTER TABLE security.general_user ADD COLUMN IF NOT EXISTS general_user_team VARCHAR(50);" >/dev/null

    # Catálogo/referencia real (roles, geografía, entidades, motor de formularios)
    # — ver cloudDevelopment/actualizar-seed-con-composio.md. Antes de los
    # usuarios de prueba: estos dependen de que los roles ya existan.
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=1 -f "$REPO_ROOT/cloudDevelopment/seed_reference_data.sql"
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=1 -f "$REPO_ROOT/cloudDevelopment/seed_local_dev.sql"

    # Datos de Supabase traídos para probar requerimientos (casos, seguimientos,
    # remisiones...). Un archivo por requerimiento, permanentes entre sesiones —
    # ver cloudDevelopment/seed_test_data/README.md. Cada uno va en una sola
    # transacción (-1): si falla, se revierte completo y se avisa, sin tumbar el
    # arranque por el dataset de otro requerimiento. Van después de los usuarios
    # de prueba porque sus FK a usuarios apuntan a cloud.test.*.
    for seed_file in "$REPO_ROOT"/cloudDevelopment/seed_test_data/*.sql; do
        [ -e "$seed_file" ] || continue
        $PSQL -d "$DB_NAME" -1 -v ON_ERROR_STOP=1 -f "$seed_file" >/dev/null \
            || echo "[cloudDevelopment] AVISO: falló $(basename "$seed_file") (se revirtió completo); se continúa sin ese dataset."
    done

    touch "$SCHEMA_MARKER"
    echo "[cloudDevelopment] Listo. Usuarios de prueba: cloud.test.ad / cloud.test.te / cloud.test.ro (clave: CloudDev!Test2026)"
else
    echo "[cloudDevelopment] BD local ya inicializada (borrar $SCHEMA_MARKER para forzar recreación de esquema/seed)."
fi

set +e   # no propagar errexit al resto de run_dev.sh
