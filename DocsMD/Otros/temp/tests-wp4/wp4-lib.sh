#!/usr/bin/env bash
# Helpers de prueba de WP-4 (API + BD). Uso: source wp4-lib.sh
export PGPASSWORD=postgres
BASE=http://127.0.0.1:9090/api/v1/victim-cases-v3
COOKIE_FILE=${COOKIE_FILE:-/tmp/wp4-cookie-op.txt}
SQL() { psql -h 127.0.0.1 -U postgres -d salvia_local -At -F ' | ' -c "$1"; }
# api METHOD PATH [BODY] -> imprime "<status>\t<body>"
api() {
  local m=$1 p=$2 b=${3:-}
  if [ -n "$b" ]; then
    curl -s -w '\n%{http_code}' -X "$m" "$BASE/$p" -H "Cookie: $(cat $COOKIE_FILE)" -H 'Content-Type: application/json' -d "$b"
  else
    curl -s -w '\n%{http_code}' -X "$m" "$BASE/$p" -H "Cookie: $(cat $COOKIE_FILE)"
  fi
}
# body_of / code_of sobre la salida de api()
body_of() { sed '$d' <<<"$1"; }
code_of() { tail -n1 <<<"$1"; }
PASS=0; FAIL=0
check() { # check "descripción" actual esperado
  if [ "$2" == "$3" ]; then PASS=$((PASS+1)); echo "  ✅ $1"; else FAIL=$((FAIL+1)); echo "  ❌ $1  (esperado: [$3] obtenido: [$2])"; fi
}
new_case() { # new_case DOC -> icode
  local r; r=$(api POST draft "{\"tipo_documento\":\"cc\",\"documento\":\"$1\",\"telefono\":\"3001112233\"}")
  body_of "$r" | jq -r .icode
}
# Marca toda la tríada como respondida (atajo de prueba; el guardado real exige ~35 campos).
mark_triada() {
  SQL "update salvia.victim_case set victim_case_v3_answered_keys='CI_1,CI_2,CI_3,CI_4,CI_5,CI_6,CI_7,CI_9,DP_1,DP_2,DP_5,DP_8,DP_9,DP_11,DP_12,DP_14,DP_17,DP_27,DP_28,DP_29,HV_1,HV_2,HV_4,HV_6,HV_8,HV_9,HV_11,HV_13,HV_14,AG_1,AG_2,AG_3' where victim_case_i_code='$1'" >/dev/null
}
ic() { SQL "select victim_case_form2_enums_i_code from salvia.victim_case_form2_enums where victim_case_form2_enums_category='$1' and victim_case_form2_enums_code='$2'"; }
