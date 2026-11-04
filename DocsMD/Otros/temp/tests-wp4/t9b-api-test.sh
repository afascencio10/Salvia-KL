#!/usr/bin/env bash
# Prueba API+BD de T9 (REG-38*), parte 2: remisiones a otros equipos y casos límite.
set -u
source "$(dirname "$0")/wp4-lib.sh"

complete_triada() { # complete_triada IC  -> guarda AG_7 con la tríada ya marcada
  mark_triada "$1"
  api POST "$1/seccion/AG" '{"answers":{"AG_7":"Agresor de prueba"}}'
}
cnt() { SQL "select count(*) from salvia.$1 where case_id='$2'"; }

TEAM=victim_case_form2_referral_team
echo "== A. Remisiones: ap + ah + es + di + sd (sin medidas) con criterios válidos"
DOC=$((RANDOM + 920000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
EQ="$(ic $TEAM ap),$(ic $TEAM ah),$(ic $TEAM es),$(ic $TEAM di),$(ic $TEAM sd)"
PSY="$(ic victim_case_form2_psycho_criteria ob),$(ic victim_case_form2_psycho_criteria cs)"
R=$(api POST "$IC/seccion/GS" "{\"answers\":{\"GS_5\":\"true\",\"GS_6\":\"$EQ\",\"GS_9\":\"$PSY\",\"GS_10\":\"$(ic victim_case_form2_men_criteria ob)\",\"GS_11\":\"$(ic victim_case_form2_disability_service ls),$(ic victim_case_form2_disability_service ed)\",\"GS_12\":\"$(ic victim_case_form2_stabiliz_criteria si)\"}}")
check "SaveSection GS (remisión) → 200" "$(code_of "$R")" "200"
check "GS_5/GS_6 leídos de vuelta" "$(body_of "$(api GET "$IC")" | jq -r '.answers.GS_5')" "true"
check "antes de la tríada no hay remisiones" "$(cnt psychosocial_support "$IC")/$(cnt men_team_remision "$IC")/$(cnt economic_stabilization "$IC")/$(cnt discapacidad_remision "$IC")" "0/0/0/0"
R=$(complete_triada "$IC"); check "tríada completa → 200" "$(code_of "$R")" "200"
check "estado = Iniciado" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "in"
check "psychosocial_support = 1 (derivacion/abierto)" "$(SQL "select count(*)||'/'||min(type)||'/'||min(status) from salvia.psychosocial_support where case_id='$IC'")" "1/derivacion/abierto"
check "referral_criteria en el vocabulario de Seguimiento" "$(SQL "select referral_criteria from salvia.psychosocial_support where case_id='$IC'")" "criterio_obligatorio,conducta_suicida"
check "submitted_by = agente del caso" "$(SQL "select submitted_by from salvia.psychosocial_support where case_id='$IC'")" "user-icode-cloud.test.op"
check "tarea «Validar remisión a psicosocial» (Psicosocial/validar_remision)" "$(SQL "select count(*) from salvia.case_task where case_id='$IC' and category='Psicosocial' and type='validar_remision' and psychosocial_support_id is not null")" "1"
check "men_team_remision = 1 (derivacion/ACTIVE)" "$(SQL "select count(*)||'/'||min(type)||'/'||min(status) from salvia.men_team_remision where case_id='$IC'")" "1/derivacion/ACTIVE"
check "economic_stabilization = 1 (derivacion/ACTIVE)" "$(SQL "select count(*)||'/'||min(type)||'/'||min(status) from salvia.economic_stabilization where case_id='$IC'")" "1/derivacion/ACTIVE"
check "discapacidad_remision = 2 (apoyo_lsc, enfoque_discapacidad)" "$(SQL "select string_agg(service, ',' order by service) from salvia.discapacidad_remision where case_id='$IC'")" "apoyo_lsc,enfoque_discapacidad"
check "sin medidas de emergencia" "$(cnt emergency_measure_referral "$IC")" "0"
check "materializedAt marcado" "$(SQL "select (victim_case_v3_draft_data->>'materializedAt') is not null from salvia.victim_case where victim_case_i_code='$IC'")" "t"
R=$(complete_triada "$IC")
check "guardar de nuevo no duplica remisiones" "$(cnt psychosocial_support "$IC")/$(cnt men_team_remision "$IC")/$(cnt economic_stabilization "$IC")/$(cnt discapacidad_remision "$IC")" "1/1/1/2"

echo "== B. Remisiones: me + ap (la psicosocial se omite por incompatibilidad) con medidas al + tr"
DOC=$((RANDOM + 930000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
EQ="$(ic $TEAM me),$(ic $TEAM ap)"
ME="$(ic victim_case_form2_emergency_measure al),$(ic victim_case_form2_emergency_measure tr)"
PSY="$(ic victim_case_form2_psycho_criteria ob),$(ic victim_case_form2_psycho_criteria cs)"
api POST "$IC/seccion/GS" "{\"answers\":{\"GS_5\":\"true\",\"GS_6\":\"$EQ\",\"GS_7\":\"$ME\",\"GS_9\":\"$PSY\"}}" >/dev/null
complete_triada "$IC" >/dev/null
check "emergency_measure_referral = 1 (pendiente_de_revision/normal)" "$(SQL "select count(*)||'/'||min(status)||'/'||min(urgency) from salvia.emergency_measure_referral where case_id='$IC'")" "1/pendiente_de_revision/normal"
check "medidas: alojamiento y transporte (ACTIVE)" "$(SQL "select string_agg(type||':'||status, ',' order by type) from salvia.emergency_measure where case_id='$IC'")" "alojamiento:ACTIVE,transporte:ACTIVE"
check "submitted_by del referral" "$(SQL "select submitted_by from salvia.emergency_measure_referral where case_id='$IC'")" "user-icode-cloud.test.op"
check "psicosocial omitida" "$(cnt psychosocial_support "$IC")" "0"

echo "== C. Criterios insuficientes: no se crea la remisión"
DOC=$((RANDOM + 940000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
EQ="$(ic $TEAM ap),$(ic $TEAM ah),$(ic $TEAM es)"
api POST "$IC/seccion/GS" "{\"answers\":{\"GS_5\":\"true\",\"GS_6\":\"$EQ\",\"GS_9\":\"$(ic victim_case_form2_psycho_criteria ob)\"}}" >/dev/null
complete_triada "$IC" >/dev/null
check "psicosocial (ob con 0 puntos) → no se crea" "$(cnt psychosocial_support "$IC")" "0"
check "hombres sin criterio → no se crea" "$(cnt men_team_remision "$IC")" "0"
check "estabilización sin criterios → no se crea" "$(cnt economic_stabilization "$IC")" "0"
check "aun así el Caso pasó a Iniciado" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "in"

echo "== D. GS_5 = No: no hay remisiones"
DOC=$((RANDOM + 950000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
api POST "$IC/seccion/GS" "{\"answers\":{\"GS_5\":\"false\",\"GS_6\":\"$(ic $TEAM ah)\",\"GS_10\":\"$(ic victim_case_form2_men_criteria ob)\"}}" >/dev/null
complete_triada "$IC" >/dev/null
check "sin remisión si GS_5 = No" "$(cnt men_team_remision "$IC")" "0"

echo "== E. Canales: la acudida SIN pendientes y la que no tiene acuerdo no crean tareas"
DOC=$((RANDOM + 960000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
R=$(api POST "$IC/entidades" '{"kind":"acudida","entityBranchId":101,"agreement":true,"dueDiligence":true,"activationChannels":["notificacion","misalvia"]}')
check "acudida sin pendientes se agrega" "$(code_of "$R")" "200"
R=$(api POST "$IC/entidades" '{"kind":"acudida","entityBranchId":102,"pendingActions":["4cd330bd-bb75-44a2-9bc1-d2a50da9d63b"],"agreement":true,"dueDiligence":false,"activationChannels":["notificacion"]}')
check "acudida con pendientes pero sin debida diligencia se agrega" "$(code_of "$R")" "200"
R=$(api POST "$IC/entidades" '{"kind":"acudida","entityBranchId":9217,"pendingActions":["x"],"agreement":true,"dueDiligence":true,"activationChannels":["articulacion"]}')
check "acudida con pendientes + acuerdo + debida diligencia se agrega" "$(code_of "$R")" "200"
complete_triada "$IC" >/dev/null
check "entity_case = 3" "$(cnt entity_case "$IC")" "3"
check "solo la tercera genera tarea (articulación)" "$(SQL "select string_agg(description, ',') from salvia.case_task where case_id='$IC'")" "Articulación institucional telefónica o por mensajería"
check "sin oficios" "$(cnt entity_letter "$IC")" "0"

echo "== F. Fallo a mitad: el Caso sigue en Borrador, el JSON intacto, y se recupera"
DOC=$((RANDOM + 970000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
api POST "$IC/entidades" '{"kind":"activacion","entityBranchId":101,"agreement":true,"dueDiligence":true,"activationChannels":["misalvia"]}' >/dev/null
LONG=$(printf 'x%.0s' $(seq 1 80))
R=$(api POST "$IC/barreras" "{\"sectorCode\":\"$LONG\",\"entityBranchId\":101,\"description\":\"rompe el INSERT\",\"management\":[\"alerta\"]}")
BAD=$(body_of "$R" | jq -r .icode)
check "el borrador acepta la barrera defectuosa" "$(code_of "$R")" "200"
R=$(complete_triada "$IC")
check "la tríada falla (500) porque la materialización falla" "$(code_of "$R")" "500"
check "el Caso sigue en Borrador" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "bo"
check "no quedó nada a medias: entity_case = 0" "$(cnt entity_case "$IC")" "0"
check "no quedó nada a medias: barrier_v2 = 0" "$(cnt barrier_v2 "$IC")" "0"
check "no quedó nada a medias: case_task = 0" "$(cnt case_task "$IC")" "0"
check "el JSON conserva entidad y barrera" "$(SQL "select jsonb_array_length(victim_case_v3_draft_data->'entities') || '/' || jsonb_array_length(victim_case_v3_draft_data->'barriers') from salvia.victim_case where victim_case_i_code='$IC'")" "1/1"
api DELETE "barreras/$BAD" >/dev/null
R=$(complete_triada "$IC")
check "corregido el borrador, la tríada pasa (200)" "$(code_of "$R")" "200"
check "ahora sí: Iniciado, 1 entity_case y 1 tarea" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")/$(cnt entity_case "$IC")/$(cnt case_task "$IC")" "in/1/1"

echo "== G. Sede que ya existe en las tablas: duplicado al agregar y reutilización al materializar"
DOC=$((RANDOM + 980000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
api POST "$IC/entidades" '{"kind":"acudida","entityBranchId":102,"pendingActions":["4cd330bd-bb75-44a2-9bc1-d2a50da9d63b"]}' >/dev/null
SQL "insert into salvia.entity_case (id, case_id, entity_branch_id, created_by_id, created_at, updated_at) values ('11111111-2222-3333-4444-$(printf '%012d' $RANDOM)', '$IC', 102, 'test', now(), now())" >/dev/null
SQL "insert into salvia.entity_case (id, case_id, entity_branch_id, created_by_id, created_at, updated_at) values ('11111111-2222-3333-4444-$(printf '%012d' $((RANDOM+50000)))', '$IC', 101, 'test', now(), now())" >/dev/null
R=$(api POST "$IC/entidades" '{"kind":"acudida","entityBranchId":101}')
check "sede 101 ya en las tablas → 409" "$(code_of "$R")" "409"
complete_triada "$IC" >/dev/null
check "la sede 102 (JSON + tabla) no se duplica" "$(SQL "select count(*) from salvia.entity_case where case_id='$IC' and entity_branch_id=102")" "1"
check "se le agregó la obligación pendiente" "$(SQL "select count(*) from salvia.entity_case_obligation eco join salvia.entity_case ec on ec.id=eco.entity_case_id where ec.case_id='$IC' and ec.entity_branch_id=102 and eco.status='pendiente'")" "1"

echo "== H. Clientes antiguos (compatibilidad)"
DOC=$((RANDOM + 990000000)); IC=$(new_case "$DOC"); echo "   caso $IC"
check "GET de un Caso sin JSON → 200 y listas vacías" "$(code_of "$(api GET "$IC")")/$(body_of "$(api GET "$IC")" | jq -r '(.entities|length|tostring)+(.barriers|length|tostring)')" "200/00"
api POST "$IC/entidades" '{"kind":"activacion","entityBranchId":102,"activationChannel":"no,ll,ms","agreement":true,"dueDiligence":true}' >/dev/null
G=$(body_of "$(api GET "$IC")")
check "canal viejo (no,ll,ms) → notificacion, articulacion, misalvia" "$(jq -r '.entityMeta["102"].activationChannels|join(",")' <<<"$G")" "notificacion,articulacion,misalvia"
api POST "$IC/barreras" '{"sectorCode":"salud","entityBranchId":102,"description":"vieja","management":["03","06"],"competences":["c1"],"date":"2026-10-21"}' >/dev/null
G=$(body_of "$(api GET "$IC")")
check "gestión 03/06 → slugs de Psicosocial/Seguimiento" "$(jq -r '.barriers[0].managementActions' <<<"$G")" "activacion_ruta_interinstitucional,alerta_barreras"
check "competencias se conservan" "$(jq -r '.barriers[0].competences' <<<"$G")" "c1"
complete_triada "$IC" >/dev/null
check "gestión legada no crea tareas de barrera" "$(SQL "select count(*) from salvia.case_task where case_id='$IC' and category='Barreras'")" "0"
check "los 3 canales viejos sí crean sus 3 tareas de entidad" "$(SQL "select count(*) from salvia.case_task where case_id='$IC' and category='Entidades'")" "3"

echo
echo "RESULTADO: $PASS ok / $FAIL fallos"
exit $FAIL
