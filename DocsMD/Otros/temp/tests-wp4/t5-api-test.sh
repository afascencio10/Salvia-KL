#!/usr/bin/env bash
# Prueba API+BD de T5 (REG-16): GS_20 «elementos psicosociales».
set -u
source "$(dirname "$0")/wp4-lib.sh"
CAT=victim_case_form2_psychosocial_elements

echo "== 1. Catálogo"
R=$(api GET enums); E=$(body_of "$R")
check "GET enums → 200" "$(code_of "$R")" "200"
check "la categoría trae 15 opciones" "$(jq -r --arg c $CAT '.[$c]|length' <<<"$E")" "15"
check "ordenadas por código, «Ninguno» (99) al final" "$(jq -r --arg c $CAT '[.[$c][].code]|join(",")' <<<"$E")" "01,02,03,04,05,06,07,08,09,10,11,12,13,14,99"
check "las 15 filas existen en BD" "$(SQL "select count(*) from salvia.victim_case_form2_enums where victim_case_form2_enums_category='$CAT'")" "15"

echo "== 2. Guardar y recargar GS_20"
IC=$(new_case $((RANDOM + 930000000))); echo "   caso $IC"
I01=$(ic $CAT 01); I05=$(ic $CAT 05); I14=$(ic $CAT 14); I99=$(ic $CAT 99)
R=$(api POST "$IC/seccion/GS" "{\"answers\":{\"GS_20\":\"$I01,$I05,$I14\"}}")
check "guardar 3 opciones → 200" "$(code_of "$R")" "200"
G=$(body_of "$(api GET "$IC")")
check "GS_20 vuelve con las mismas 3 (en cualquier orden)" "$(jq -r '.answers.GS_20|split(",")|sort|join(",")' <<<"$G")" "$(echo "$I01,$I05,$I14" | tr ',' '\n' | sort | paste -sd, -)"
check "GS_20 figura como respondida" "$(jq -r '.answeredKeys|index("GS_20")!=null' <<<"$G")" "true"
check "3 filas en la relación de enums del Caso" "$(SQL "select count(*) from salvia.rel_victim_case_form2_enums_victim_case_form2 r join salvia.victim_case_form2_enums e on e.victim_case_form2_enums_id=r.victim_case_form2_enums_id join salvia.victim_case_form2 f on f.victim_case_form2_id=r.victim_case_form2_id join salvia.victim_case vc on vc.victim_case_id=f.victim_case_form2_victim_case where vc.victim_case_i_code='$IC' and e.victim_case_form2_enums_category='$CAT'")" "3"

echo "== 3. «Ninguno» deja solo el 99"
R=$(api POST "$IC/seccion/GS" "{\"answers\":{\"GS_20\":\"$I99\"}}")
check "guardar «Ninguno» → 200" "$(code_of "$R")" "200"
G=$(body_of "$(api GET "$IC")")
check "queda solo el 99" "$(jq -r '.answers.GS_20' <<<"$G")" "$I99"

echo "== 4. Vaciar la selección"
R=$(api POST "$IC/seccion/GS" '{"answers":{"GS_20":""}}')
check "guardar vacío → 200" "$(code_of "$R")" "200"
G=$(body_of "$(api GET "$IC")")
check "GS_20 ya no trae valores" "$(jq -r '.answers.GS_20 // ""' <<<"$G")" ""

echo "== 5. No es obligatoria para finalizar"
SQL "update salvia.victim_case set victim_case_status='in' where victim_case_i_code='$IC'" >/dev/null
R=$(api POST "$IC/finalizar"); B=$(body_of "$R")
check "Finalizar no menciona GS_20 entre los faltantes" "$(jq -r '(.modulosIncompletos // []) | tostring | contains("GS_20")' <<<"$B")" "false"

echo; echo "RESULTADO: $PASS ok / $FAIL fallos"; exit $FAIL
