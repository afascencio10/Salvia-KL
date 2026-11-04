#!/usr/bin/env bash
# Corre toda la batería de WP-4 contra la app local (arriba en :9090) y resume los resultados.
# Uso: bash DocsMD/Otros/temp/tests-wp4/run-all.sh        (requiere Playwright global: NODE_PATH=$(npm root -g))
DIR="$(cd "$(dirname "$0")" && pwd)"
export NODE_PATH="$(npm root -g 2>/dev/null)" COOKIE_FILE=/tmp/wp4-cookie-op.txt
cd /tmp && node "$DIR/login-cookie.js" op "$COOKIE_FILE" >/dev/null || { echo "no se pudo iniciar sesión"; exit 1; }
TOTAL_OK=0; TOTAL_FAIL=0
run() { # run <nombre> <comando...>
  local nombre=$1; shift
  local salida; salida=$("$@" 2>&1); local linea; linea=$(grep -E "RESULTADO" <<<"$salida" | tail -1)
  local ok fail; ok=$(sed -E 's/.*RESULTADO: ([0-9]+) ok.*/\1/' <<<"$linea"); fail=$(sed -E 's/.*\/ ([0-9]+) fallos.*/\1/' <<<"$linea")
  if [[ "$ok" =~ ^[0-9]+$ ]]; then TOTAL_OK=$((TOTAL_OK+ok)); TOTAL_FAIL=$((TOTAL_FAIL+fail)); printf '%-28s %s\n' "$nombre" "$linea"
  else TOTAL_FAIL=$((TOTAL_FAIL+1)); printf '%-28s SIN RESULTADO\n' "$nombre"; echo "$salida" | tail -5; fi
  grep -E "❌" <<<"$salida" | head -5
}
run "API T9 borrador"        bash "$DIR/t9-api-test.sh"
run "API T9b remisiones"     bash "$DIR/t9b-api-test.sh"
run "API T5 GS_20"           bash "$DIR/t5-api-test.sh"
run "API T1 finalizar"        bash "$DIR/t1-finalize-test.sh"
for t in t2 t3 t4 t5 t6 t7 t8 solo-lectura cierre; do run "UI $t" node "$DIR/ui-$t.js"; done
echo "----------------------------------------"
echo "TOTAL: $TOTAL_OK ok / $TOTAL_FAIL fallos"
exit $TOTAL_FAIL
