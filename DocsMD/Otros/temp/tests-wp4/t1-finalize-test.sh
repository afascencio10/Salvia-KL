#!/usr/bin/env bash
# Prueba API de T1 (REG-14): GS_2 ya no es obligatoria para finalizar; las demás obligatorias de GS siguen exigiéndose.
set -u
source "$(dirname "$0")/wp4-lib.sh"

claves_todas() { # todas las llaves posibles de los módulos, menos las que se pasan como argumentos
  local excluir=" $* " k out=()
  for pref_max in CI:20 DP:60 CA:20 HV:30 AG:20 TZ:10 GS:20; do
    local pref=${pref_max%%:*} max=${pref_max##*:}
    for i in $(seq 1 "$max"); do k="${pref}_${i}"; [[ "$excluir" == *" $k "* ]] || out+=("$k"); done
  done
  (IFS=,; echo "${out[*]}")
}

echo "== 1. Todo respondido salvo GS_2 → finaliza"
IC=$(new_case $((RANDOM + 910000000))); echo "   caso $IC"
SQL "update salvia.victim_case set victim_case_status='in', victim_case_v3_answered_keys='$(claves_todas GS_2)' where victim_case_i_code='$IC'" >/dev/null
R=$(api POST "$IC/finalizar")
check "Finalizar sin GS_2 → 200" "$(code_of "$R")" "200"
check "el Caso quedó Terminado" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "te"

echo "== 2. Falta GS_4 (sigue obligatoria) → 400 con el módulo GS"
IC=$(new_case $((RANDOM + 920000000))); echo "   caso $IC"
SQL "update salvia.victim_case set victim_case_status='in', victim_case_v3_answered_keys='$(claves_todas GS_2 GS_4)' where victim_case_i_code='$IC'" >/dev/null
R=$(api POST "$IC/finalizar")
check "Finalizar sin GS_4 → 400" "$(code_of "$R")" "400"
check "informa el módulo GS incompleto" "$(body_of "$R" | jq -r '.modulosIncompletos|join(",")')" "GS"

echo "== 3. GS_14 y GS_17 (que viven en la Ruta de atención integral) siguen siendo obligatorias"
for k in GS_14 GS_17 GS_16 GS_5; do
  IC=$(new_case $((RANDOM + 930000000)))
  SQL "update salvia.victim_case set victim_case_status='in', victim_case_v3_answered_keys='$(claves_todas GS_2 $k)' where victim_case_i_code='$IC'" >/dev/null
  R=$(api POST "$IC/finalizar")
  check "sin $k → 400 (módulo GS)" "$(code_of "$R")/$(body_of "$R" | jq -r '.modulosIncompletos|join(",")')" "400/GS"
done

echo "== 4. GS_20 (elementos psicosociales) no es obligatoria"
IC=$(new_case $((RANDOM + 940000000)))
SQL "update salvia.victim_case set victim_case_status='in', victim_case_v3_answered_keys='$(claves_todas GS_2 GS_20)' where victim_case_i_code='$IC'" >/dev/null
R=$(api POST "$IC/finalizar")
check "sin GS_20 → 200" "$(code_of "$R")" "200"

echo; echo "RESULTADO: $PASS ok / $FAIL fallos"; exit $FAIL
