#!/bin/bash
# Fase 2 del seed del sandbox cloud: aplica cloudDevelopment/seed_test_data/*.sql
# a la BD local (salvia_local).
#
# Por qué es una fase aparte: las tablas del dominio nuevo (GORM: case_task,
# follow_up_v2, emergency_measure_referral, dupla, articulating_mechanism, ...)
# las crea la APP al arrancar (AutoMigrate en main.go). Antes del primer arranque
# no existen, así que setup-local-postgres.sh no puede sembrar datos en ellas.
# Este script corre cuando la app ya terminó de migrar.
#
# Llevar la cuenta: cada archivo aplicado queda registrado (nombre + sha256) en
# public.cloud_seed_applied, dentro de la misma BD. Se aplican solo los archivos
# nuevos o modificados; si se recrea la BD, el registro se reinicia solo.
# Cada archivo va en una sola transacción (psql -1): si falla se revierte
# completo, se avisa con el error de Postgres y se sigue con el siguiente.
#
# Uso:
#   cloudDevelopment/apply-test-seeds.sh          # aplica los pendientes ahora
#   cloudDevelopment/apply-test-seeds.sh --wait   # primero espera a que la app responda
#                                                 # (así lo lanza setup-local-postgres.sh)
# Variables: DB_HOST/DB_PORT/DB_USER/DB_NAME (mismas que el setup),
#            SALVIA_APP_URL (por defecto http://127.0.0.1:9090).

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SEED_DIR="$REPO_ROOT/cloudDevelopment/seed_test_data"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-salvia_local}"
APP_URL="${SALVIA_APP_URL:-http://127.0.0.1:9090}"

PSQL_BIN="$(command -v psql || ls -d /usr/lib/postgresql/*/bin/psql 2>/dev/null | sort -rV | head -n 1)"
if [ -z "$PSQL_BIN" ]; then
    echo "[seed_test_data] ERROR: no encontré psql." >&2
    exit 1
fi
PSQL="$PSQL_BIN -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -X -q"

if [ "${1:-}" = "--wait" ]; then
    # La app imprime "Servidor corriendo" y empieza a servir DESPUÉS de migrar
    # (AutoMigrate + seeds propios), así que el login respondiendo = tablas listas.
    # Margen amplio: la primera compilación de `go run` puede tardar varios minutos.
    ready=""
    for _ in $(seq 1 900); do
        if curl -fs -o /dev/null --max-time 2 "$APP_URL/seguridad/login"; then
            ready=1
            break
        fi
        sleep 1
    done
    if [ -z "$ready" ]; then
        echo "[seed_test_data] La app no respondió en $APP_URL tras 15 min; no se aplicaron seeds. Cuando esté arriba: cloudDevelopment/apply-test-seeds.sh" >&2
        exit 1
    fi
fi

shopt -s nullglob
files=("$SEED_DIR"/*.sql)
if [ ${#files[@]} -eq 0 ]; then
    echo "[seed_test_data] No hay archivos en cloudDevelopment/seed_test_data/."
    exit 0
fi

$PSQL -c "CREATE TABLE IF NOT EXISTS public.cloud_seed_applied (
              filename   text PRIMARY KEY,
              sha256     text NOT NULL,
              applied_at timestamptz NOT NULL DEFAULT now());" >/dev/null || exit 1

applied=0; skipped=0; failed=0
for f in "${files[@]}"; do
    name="$(basename "$f")"
    case "$name" in
        *\'*)
            # El nombre se interpola en SQL para llevar el registro: sin comillas simples.
            echo "[seed_test_data] AVISO: $name tiene una comilla simple en el nombre; renómbralo (AAAAMMDD-<slug>.sql)." >&2
            failed=$((failed + 1))
            continue
            ;;
    esac
    sum="$(sha256sum "$f" | cut -d' ' -f1)"
    prev="$($PSQL -tA -c "SELECT sha256 FROM public.cloud_seed_applied WHERE filename = '$name'")"
    if [ "$prev" = "$sum" ]; then
        skipped=$((skipped + 1))
        continue
    fi
    # Un solo -1: el archivo y su registro entran o se revierten juntos.
    if $PSQL -1 -v ON_ERROR_STOP=1 -f "$f" \
        -c "INSERT INTO public.cloud_seed_applied (filename, sha256) VALUES ('$name', '$sum')
            ON CONFLICT (filename) DO UPDATE SET sha256 = EXCLUDED.sha256, applied_at = now();" >/dev/null; then
        applied=$((applied + 1))
        echo "[seed_test_data] Aplicado: $name"
    else
        failed=$((failed + 1))
        echo "[seed_test_data] AVISO: falló $name (se revirtió completo). Corrige el archivo y vuelve a correr cloudDevelopment/apply-test-seeds.sh" >&2
    fi
done

echo "[seed_test_data] Listo: $applied aplicado(s), $skipped ya estaban al día, $failed con error."
[ "$failed" -eq 0 ]
