#!/bin/bash
# Aplica a la BD local del sandbox cloud (salvia_local) las migraciones FECHADAS de src/migrations/:
# AAAAMMDD_<tema>.sql (p. ej. 20261102_wp2_registro_v3_catalogos_y_columnas.sql).
#
# Por qué existe: un cambio de catálogo o de esquema que todavía no está en la BD real de pruebas
# (Supabase) se escribe como migración en src/migrations/ y allá se aplica A MANO (SQL editor).
# Sin este script, cada sandbox nuevo arrancaría SIN esos cambios —listas desplegables vacías,
# columnas que no existen (500)— y alguien tendría que aplicarlos a mano en cada sesión.
# setup-local-postgres.sh lo llama en CADA arranque (no solo la primera vez), antes de iniciar la
# app, así que un sandbox ya inicializado también recoge las migraciones que lleguen después.
#
# Qué se aplica: solo los archivos cuyo nombre empieza por 8 dígitos y guion bajo, en orden. Los demás
# archivos de src/migrations/ (add_*.sql, alter_*.sql...) son anteriores a esta convención y no se tocan.
#
# Reglas para una migración nueva (la misma que ya piden los archivos de seed_test_data/):
#   - IDEMPOTENTE: se vuelve a aplicar cada vez que el archivo cambia (IF NOT EXISTS, WHERE NOT EXISTS...).
#     Las filas de enum se buscan por NOMBRE (victim_case_form2_enums_name es UNIQUE), no por id.
#   - Solo puede asumir lo que existe ANTES de que la app arranque: el esquema de los dumps legacy +
#     schema_drift_patch.sql + los seeds de referencia. Si necesita una tabla que crea la app (GORM
#     AutoMigrate), falla en este arranque y se reintenta en el siguiente (o a mano con este script).
#   - `git add -f`: .gitignore ignora *.sql, y src/migrations/ no es una excepción.
#
# Llevar la cuenta: cada archivo aplicado queda registrado (nombre + sha256) en
# public.cloud_migration_applied, dentro de la misma BD: se aplican solo los archivos nuevos o
# modificados, y si se recrea la BD el registro se reinicia solo. Un archivo que abre su propia
# transacción (BEGIN; ... COMMIT;) se aplica tal cual; uno que no, se aplica en una sola (psql -1).
# Si un archivo falla se avisa con el error de Postgres y se sigue con el siguiente: la app arranca igual.
#
# Esto es SOLO del sandbox. En Supabase la migración se aplica a mano.
#
# Uso:
#   cloudDevelopment/apply-migrations.sh        # aplica las pendientes ahora
# Variables: DB_HOST/DB_PORT/DB_USER/DB_NAME (mismas que el setup).

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS_DIR="$REPO_ROOT/src/migrations"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-salvia_local}"

PSQL_BIN="$(command -v psql || ls -d /usr/lib/postgresql/*/bin/psql 2>/dev/null | sort -rV | head -n 1)"
if [ -z "$PSQL_BIN" ]; then
    echo "[migraciones] ERROR: no encontré psql." >&2
    exit 1
fi
# Sin los avisos NOTICE de "IF NOT EXISTS ... skipping" en cada corrida.
export PGOPTIONS="${PGOPTIONS:+$PGOPTIONS }-c client_min_messages=warning"
PSQL="$PSQL_BIN -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -X -q"

shopt -s nullglob
files=("$MIGRATIONS_DIR"/[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]_*.sql)
if [ ${#files[@]} -eq 0 ]; then
    echo "[migraciones] No hay migraciones fechadas (AAAAMMDD_<tema>.sql) en src/migrations/."
    exit 0
fi

$PSQL -c "CREATE TABLE IF NOT EXISTS public.cloud_migration_applied (
              filename   text PRIMARY KEY,
              sha256     text NOT NULL,
              applied_at timestamptz NOT NULL DEFAULT now());" >/dev/null || exit 1

applied=0; skipped=0; failed=0
for f in "${files[@]}"; do
    name="$(basename "$f")"
    case "$name" in
        *\'*)
            # El nombre se interpola en SQL para llevar el registro: sin comillas simples.
            echo "[migraciones] AVISO: $name tiene una comilla simple en el nombre; renómbralo (AAAAMMDD_<tema>.sql)." >&2
            failed=$((failed + 1))
            continue
            ;;
    esac
    sum="$(sha256sum "$f" | cut -d' ' -f1)"
    prev="$($PSQL -tA -c "SELECT sha256 FROM public.cloud_migration_applied WHERE filename = '$name'")"
    if [ "$prev" = "$sum" ]; then
        skipped=$((skipped + 1))
        continue
    fi
    # ¿El archivo abre su propia transacción? (BEGIN; / START TRANSACTION; al inicio de línea; el
    # BEGIN sin punto y coma de un bloque DO $$ ... $$ no cuenta.) Si no, se aplica en una sola.
    one_tx="-1"
    if grep -qiE '^[[:space:]]*(BEGIN|START[[:space:]]+TRANSACTION)[[:space:]]*;' "$f"; then
        one_tx=""
    fi
    if $PSQL $one_tx -v ON_ERROR_STOP=1 -f "$f" >/dev/null; then
        $PSQL -c "INSERT INTO public.cloud_migration_applied (filename, sha256) VALUES ('$name', '$sum')
                  ON CONFLICT (filename) DO UPDATE SET sha256 = EXCLUDED.sha256, applied_at = now();" >/dev/null
        applied=$((applied + 1))
        echo "[migraciones] Aplicada: $name"
    else
        failed=$((failed + 1))
        echo "[migraciones] AVISO: falló $name (no se registró; se reintenta en el próximo arranque o con cloudDevelopment/apply-migrations.sh). Revisa el error de Postgres de arriba." >&2
    fi
done

echo "[migraciones] Listo: $applied aplicada(s), $skipped ya estaban al día, $failed con error."
[ "$failed" -eq 0 ]
