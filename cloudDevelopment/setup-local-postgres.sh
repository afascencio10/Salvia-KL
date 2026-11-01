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
# En Debian/Ubuntu los binarios del servidor (pg_ctl, initdb) viven en
# /usr/lib/postgresql/<versión>/bin y NO están en el PATH. Se busca primero en
# el PATH y luego ahí (versión más alta primero).
find_pg_bin() {
    local d
    if command -v pg_ctl >/dev/null 2>&1; then
        dirname "$(command -v pg_ctl)"
        return 0
    fi
    for d in $(ls -d /usr/lib/postgresql/*/bin 2>/dev/null | sort -rV); do
        if [ -x "$d/pg_ctl" ]; then
            echo "$d"
            return 0
        fi
    done
    return 1
}

if ! PG_BIN="$(find_pg_bin)"; then
    echo "[cloudDevelopment] Instalando postgresql (una sola vez, se cachea en la imagen del sandbox)..."
    apt-get update -qq && apt-get install -y -qq postgresql >/dev/null
    PG_BIN="$(find_pg_bin)"
fi

# ─── 2. Inicializar + arrancar el cluster local (idempotente) ─────────────
# Postgres se niega a correr como root ("initdb: cannot be run as root"), y el
# sandbox cloud corre como root. En ese caso initdb/pg_ctl se ejecutan como el
# usuario del SO `postgres` (lo crea el paquete postgresql; si faltara se crea
# aquí). psql (cliente) sí puede correr como root: entra por TCP con auth trust.
PG_RUN=""   # prefijo para los comandos de servidor; vacío si no somos root
if [ "$(id -u)" = "0" ]; then
    if ! id postgres >/dev/null 2>&1; then
        useradd --system --home-dir "$PG_BASE" --shell /bin/bash postgres
    fi
    PG_RUN="runuser -u postgres --"
fi

mkdir -p "$PG_BASE"
if [ -n "$PG_RUN" ]; then
    chown postgres "$PG_BASE"
    # El usuario postgres tiene que poder entrar a $PG_BASE (p. ej. un
    # SALVIA_LOCAL_PG_DIR bajo un directorio 700 de root no sirve).
    if ! $PG_RUN test -w "$PG_BASE"; then
        echo "[cloudDevelopment] ERROR: el usuario 'postgres' no puede escribir en $PG_BASE; usa otro SALVIA_LOCAL_PG_DIR (p. ej. /var/lib/salvia-local-pg)." >&2
        return 1 2>/dev/null || exit 1
    fi
fi
export LC_ALL=C LANG=C   # evita locale issues al arrancar postgres

if [ ! -f "$PGDATA/PG_VERSION" ]; then
    echo "[cloudDevelopment] Inicializando BD Postgres local en $PGDATA..."
    $PG_RUN "$PG_BIN/initdb" -D "$PGDATA" -U postgres --auth=trust -E UTF8 --locale=C >/dev/null
fi

if ! $PG_RUN "$PG_BIN/pg_ctl" -D "$PGDATA" status >/dev/null 2>&1; then
    echo "[cloudDevelopment] Arrancando Postgres local en el puerto $DB_PORT..."
    $PG_RUN "$PG_BIN/pg_ctl" -D "$PGDATA" -o "-p $DB_PORT -k /tmp -h 127.0.0.1" -l "$PG_LOG" start >/dev/null
fi

# Espera a que acepte conexiones (hasta ~10s); si no llega, falla con el log a la vista
PG_READY=""
for i in $(seq 1 20); do
    if PGPASSWORD="$DB_PASSWORD" "$PG_BIN/psql" -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "SELECT 1" >/dev/null 2>&1; then
        PG_READY=1
        break
    fi
    sleep 0.5
done
if [ -z "$PG_READY" ]; then
    echo "[cloudDevelopment] ERROR: Postgres no aceptó conexiones en $DB_HOST:$DB_PORT. Últimas líneas de $PG_LOG:" >&2
    tail -n 15 "$PG_LOG" >&2 || true
    return 1 2>/dev/null || exit 1
fi

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

    # Drift entre esos dumps y la BD real de pruebas (Supabase): columnas y tablas
    # que allá se agregaron a mano. Sin esto varias pantallas responden 500
    # ("column ... does not exist"). A diferencia de los dumps, este parche es
    # idempotente y está verificado, así que va estricto: si falla, el setup se
    # detiene a la vista. Ver cloudDevelopment/schema_drift_patch.sql.
    PGOPTIONS="-c client_min_messages=warning" \
        $PSQL -d "$DB_NAME" -q -v ON_ERROR_STOP=1 -f "$REPO_ROOT/cloudDevelopment/schema_drift_patch.sql"

    # Columnas que normalmente agrega la propia app al arrancar
    # (migrateLegacyTables en main.go) — se adelantan aquí, con el mismo tipo
    # que declara la app, para que el seed de usuarios pueda fijar los
    # atributos por rol desde el primer arranque (equipo del rol 'ro',
    # departamento del 'en', sede del 'et', mecanismo del 'mc').
    $PSQL -d "$DB_NAME" -c "ALTER TABLE security.general_user
        ADD COLUMN IF NOT EXISTS general_user_team VARCHAR(50),
        ADD COLUMN IF NOT EXISTS general_user_assigned_department VARCHAR(20),
        ADD COLUMN IF NOT EXISTS entity_branch_id BIGINT,
        ADD COLUMN IF NOT EXISTS articulating_mechanism_id VARCHAR(36);" >/dev/null

    # Catálogo/referencia real (roles, geografía, entidades, motor de formularios)
    # — ver cloudDevelopment/actualizar-seed-con-composio.md. Antes de los
    # usuarios de prueba: estos dependen de que los roles ya existan.
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=1 -f "$REPO_ROOT/cloudDevelopment/seed_reference_data.sql"
    $PSQL -d "$DB_NAME" -v ON_ERROR_STOP=1 -f "$REPO_ROOT/cloudDevelopment/seed_local_dev.sql"

    touch "$SCHEMA_MARKER"
    echo "[cloudDevelopment] Listo. Un usuario de prueba por rol: cloud.test.<rol> (ad, op, sv, ro, do, no, et, en, an, us, ps, ts, te, mc, sm, em, sc, cr, am). Clave: CloudDev!Test2026"
else
    echo "[cloudDevelopment] BD local ya inicializada (borrar $SCHEMA_MARKER para forzar recreación de esquema/seed)."
fi

# ─── 4. Preparar el arranque de la app + fase 2 del seed ──────────────────
# Si quedó una instancia anterior de la app escuchando (el binario hijo de
# `go run` sobrevive si se mata solo `go run`), se detiene: ocuparía el puerto
# y seguiría sirviendo código viejo. Ver cloudDevelopment/stop-app.sh.
bash "$REPO_ROOT/cloudDevelopment/stop-app.sh" || true

# Fase 2: los datasets de cloudDevelopment/seed_test_data/ dependen de tablas que
# crea la APP al arrancar (GORM AutoMigrate), así que se aplican en segundo plano
# apenas la app responde. Solo los archivos nuevos o modificados; la salida sale
# con el prefijo [seed_test_data] en esta misma consola. Ver apply-test-seeds.sh.
bash "$REPO_ROOT/cloudDevelopment/apply-test-seeds.sh" --wait &

set +e   # no propagar errexit al resto de run_dev.sh
