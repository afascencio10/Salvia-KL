#!/usr/bin/env bash
# Prueba API+BD de T9 (REG-38*): borrador JSON, cascada, materialización, solo lectura.
set -u
source "$(dirname "$0")/wp4-lib.sh"

DOC=$((RANDOM + 910000000))
IC=$(new_case "$DOC")
echo "== Caso $IC (doc $DOC)"
check "el Caso nace en Borrador" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "bo"

O1=4cd330bd-bb75-44a2-9bc1-d2a50da9d63b   # Recibir la denuncia formalmente (sede 101)
O2=0e3b3d24-5b47-4013-99c9-fc1edb72bb6a   # FIR
O3=b4a0c171-fe98-4623-a342-b91e3d0a0803   # radicado NUNC

echo "== 1. Agregar entidad acudida (sede 101) con acciones completadas/pendientes y 4 canales"
R=$(api POST "$IC/entidades" "{\"kind\":\"acudida\",\"entityBranchId\":101,\"fileNumber\":\"RAD-77\",\"completedActions\":[\"$O1\"],\"pendingActions\":[\"$O2\",\"$O3\",\"otra\"],\"otherInfo\":\"Objetivo de prueba\",\"solicitud\":\"Solicitud de prueba\",\"agreement\":true,\"dueDiligence\":true,\"activationChannels\":[\"notificacion\",\"misalvia\",\"articulacion\",\"denuncia\"],\"reportDate\":\"2026-10-30\",\"managementNotes\":\"Se llamó a la entidad\"}")
check "POST entidades → 200" "$(code_of "$R")" "200"
E1=$(body_of "$R" | jq -r .icode)
check "el id del borrador lleva prefijo d~<icode>~" "$([[ $E1 == d~$IC~* ]] && echo si || echo no)" "si"

echo "== 2. Agregar entidad de activación (sede 102) con un solo canal"
R=$(api POST "$IC/entidades" '{"kind":"activacion","entityBranchId":102,"agreement":true,"dueDiligence":true,"activationChannels":["notificacion"]}')
check "POST entidades (activación) → 200" "$(code_of "$R")" "200"
E2=$(body_of "$R" | jq -r .icode)

echo "== 3. Duplicado de sede → 409"
R=$(api POST "$IC/entidades" '{"kind":"activacion","entityBranchId":102}')
check "sede repetida → 409" "$(code_of "$R")" "409"

echo "== 4. Barreras: una en la sede 101 (activación + alerta) y otra en la 102 (alerta)"
B1BODY='{"sectorCode":"justicia","entityBranchId":101,"specific":["x1"],"institution":"Fiscalía","institutional":["a"],"date":"2026-10-20","official":"Fiscal Pérez","department":"11","city":"11001","zone":"cabecera","description":"No recibieron la denuncia","management":["activacion","alerta"],"professionalPerspective":"Hay riesgo","managementDescription":"Se activó la ruta"}'
R=$(api POST "$IC/barreras" "$B1BODY"); check "POST barreras (101) → 200" "$(code_of "$R")" "200"; B1=$(body_of "$R" | jq -r .icode)
R=$(api POST "$IC/barreras" '{"sectorCode":"salud","entityBranchId":102,"description":"Otra barrera","management":["alerta"],"date":"2026-10-21"}')
check "POST barreras (102) → 200" "$(code_of "$R")" "200"; B2=$(body_of "$R" | jq -r .icode)

echo "== 5. El borrador NO toca las tablas"
check "entity_case del caso = 0" "$(SQL "select count(*) from salvia.entity_case where case_id='$IC'")" "0"
check "barrier_v2 del caso = 0" "$(SQL "select count(*) from salvia.barrier_v2 where case_id='$IC'")" "0"
check "case_task del caso = 0" "$(SQL "select count(*) from salvia.case_task where case_id='$IC'")" "0"
check "JSON: 2 entidades y 2 barreras" "$(SQL "select jsonb_array_length(victim_case_v3_draft_data->'entities') || '/' || jsonb_array_length(victim_case_v3_draft_data->'barriers') from salvia.victim_case where victim_case_i_code='$IC'")" "2/2"

echo "== 6. GetCase devuelve el borrador con los mismos campos que las tablas"
R=$(api GET "$IC"); G=$(body_of "$R")
check "GET → 200" "$(code_of "$R")" "200"
check "entities=2, barriers=2" "$(jq -r '(.entities|length|tostring)+"/"+(.barriers|length|tostring)' <<<"$G")" "2/2"
check "relId de la entidad 101 = id del borrador" "$(jq -r '.entities[]|select(.entityBranchId==101)|.relId' <<<"$G")" "$E1"
check "la entidad trae nombre y sector de la sede" "$(jq -r '.entities[]|select(.entityBranchId==101)|.sector' <<<"$G")" "js"
check "radicado y objetivo viajan" "$(jq -r '.entities[]|select(.entityBranchId==101)|.radicado+"|"+.objetivo' <<<"$G")" "RAD-77|Objetivo de prueba"
check "la entidad 101 cuenta 1 barrera activa" "$(jq -r '.entities[]|select(.entityBranchId==101)|.barrerasActivasCount' <<<"$G")" "1"
check "entityMeta[101].activationChannels = 4 canales" "$(jq -r '.entityMeta["101"].activationChannels|length' <<<"$G")" "4"
check "entityMeta[101].kind = acudida" "$(jq -r '.entityMeta["101"].kind' <<<"$G")" "acudida"
check "barrierMeta[b1].professionalPerspective" "$(jq -r --arg id "$B1" '.barrierMeta[$id].professionalPerspective' <<<"$G")" "Hay riesgo"
check "la barrera d~ trae management como CSV" "$(jq -r --arg id "$B1" '.barriers[]|select(.id==$id)|.managementActions' <<<"$G")" "activacion,alerta"
check "entitiesReadOnly = false en borrador" "$(jq -r '.entitiesReadOnly' <<<"$G")" "false"

echo "== 7. Eliminar la entidad 102 borra también sus barreras (solo en el JSON)"
R=$(api DELETE "entidades/$E2"); check "DELETE entidad d~ → 200" "$(code_of "$R")" "200"
G=$(body_of "$(api GET "$IC")")
check "queda 1 entidad y 1 barrera" "$(jq -r '(.entities|length|tostring)+"/"+(.barriers|length|tostring)' <<<"$G")" "1/1"
check "la barrera que queda es la de la sede 101" "$(jq -r '.barriers[0].entityBranchId|tostring' <<<"$G")" "101"
R=$(api DELETE "entidades/$E2"); check "borrar de nuevo → 404 (no existe)" "$(code_of "$R")" "404"

echo "== 8. Quitar y volver a poner una barrera"
R=$(api DELETE "barreras/$B1"); check "DELETE barrera d~ → 200" "$(code_of "$R")" "200"
check "sin barreras" "$(jq -r '.barriers|length' <<<"$(body_of "$(api GET "$IC")")")" "0"
R=$(api POST "$IC/barreras" "$B1BODY"); B1=$(body_of "$R" | jq -r .icode)
check "barrera nueva en el borrador" "$(code_of "$R")" "200"

echo "== 9. Completar la tríada → materializa"
mark_triada "$IC"
check "aún Borrador antes del guardado final" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "bo"
R=$(api POST "$IC/seccion/AG" '{"answers":{"AG_7":"Agresor de prueba"}}')
check "SaveSection AG → 200" "$(code_of "$R")" "200"
check "triadaCompleta = true" "$(body_of "$R" | jq -r .triadaCompleta)" "true"
check "estado = Iniciado" "$(SQL "select victim_case_status from salvia.victim_case where victim_case_i_code='$IC'")" "in"
check "entity_case = 1" "$(SQL "select count(*) from salvia.entity_case where case_id='$IC'")" "1"
check "entity_case.radicado/objetivo/solicitud" "$(SQL "select radicado||'|'||objetivo||'|'||solicitud from salvia.entity_case where case_id='$IC'")" "RAD-77|Objetivo de prueba|Solicitud de prueba"
check "obligaciones: 1 completada + 2 pendientes (sin «otra»)" "$(SQL "select string_agg(status||':'||n, ',' order by status) from (select eco.status, count(*) n from salvia.entity_case_obligation eco join salvia.entity_case ec on ec.id=eco.entity_case_id where ec.case_id='$IC' group by 1) t")" "completada:1,pendiente:2"
check "barrier_v2 = 1 (OPEN)" "$(SQL "select count(*)||'/'||min(status) from salvia.barrier_v2 where case_id='$IC'")" "1/OPEN"
check "barrier_v2.management_actions" "$(SQL "select management_actions from salvia.barrier_v2 where case_id='$IC'")" "activacion,alerta"
check "barrier_v2.specific_institutions = institución" "$(SQL "select specific_institutions from salvia.barrier_v2 where case_id='$IC'")" "Fiscalía"
check "barrier_v2.entity_branch_id = 101" "$(SQL "select entity_branch_id from salvia.barrier_v2 where case_id='$IC'")" "101"
echo "   -- tareas creadas:"
SQL "select category||' / '||type||' / '||description||' / asignada='||assigned_user_id from salvia.case_task where case_id='$IC' order by category, type, description"
check "case_task = 6 (4 canales + oficio de barrera + alerta)" "$(SQL "select count(*) from salvia.case_task where case_id='$IC'")" "6"
check "oficios (entity_letter) = 2 (canal notificación + barrera)" "$(SQL "select count(*) from salvia.entity_letter where case_id='$IC'")" "2"
check "tarea del canal Notificación enlaza al oficio" "$(SQL "select count(*) from salvia.case_task where case_id='$IC' and type='proyectar_oficio' and category='Entidades' and entity_letter_id is not null and entity_case_id is not null")" "1"
check "tarea de la barrera enlaza oficio y barrera" "$(SQL "select count(*) from salvia.case_task where case_id='$IC' and type='proyectar_oficio' and category='Barreras' and entity_letter_id is not null and barrier_id is not null")" "1"
check "alerta por barreras sin responsable (pool de comité)" "$(SQL "select count(*) from salvia.case_task where case_id='$IC' and type='comite_caso' and assigned_user_id=''")" "1"
check "JSON vaciado: entities=0 barriers=0" "$(SQL "select jsonb_array_length(victim_case_v3_draft_data->'entities') || '/' || jsonb_array_length(victim_case_v3_draft_data->'barriers') from salvia.victim_case where victim_case_i_code='$IC'")" "0/0"
check "JSON: materializedAt marcado" "$(SQL "select (victim_case_v3_draft_data->>'materializedAt') is not null from salvia.victim_case where victim_case_i_code='$IC'")" "t"
check "JSON: entityMeta conserva los canales" "$(SQL "select victim_case_v3_draft_data->'entityMeta'->'101'->>'agreement' from salvia.victim_case where victim_case_i_code='$IC'")" "true"
check "JSON: barrierMeta conserva la perspectiva" "$(SQL "select count(*) from salvia.victim_case, jsonb_each(victim_case_v3_draft_data->'barrierMeta') e where victim_case_i_code='$IC' and e.value->>'professionalPerspective'='Hay riesgo'")" "1"

echo "== 10. Idempotencia: guardar otra vez no duplica nada"
R=$(api POST "$IC/seccion/AG" '{"answers":{"AG_7":"Agresor de prueba"}}')
check "SaveSection AG de nuevo → 200" "$(code_of "$R")" "200"
check "entity_case sigue en 1" "$(SQL "select count(*) from salvia.entity_case where case_id='$IC'")" "1"
check "barrier_v2 sigue en 1" "$(SQL "select count(*) from salvia.barrier_v2 where case_id='$IC'")" "1"
check "case_task sigue en 6" "$(SQL "select count(*) from salvia.case_task where case_id='$IC'")" "6"

echo "== 11. Solo lectura tras materializar (409)"
R=$(api POST "$IC/entidades" '{"kind":"acudida","entityBranchId":9217}'); check "AddEntity → 409" "$(code_of "$R")" "409"
check "mensaje de solo lectura" "$(body_of "$R" | jq -r .readOnly)" "true"
R=$(api POST "$IC/barreras" "$B1BODY"); check "AddBarrier → 409" "$(code_of "$R")" "409"
G=$(body_of "$(api GET "$IC")")
EREL=$(jq -r '.entities[0].relId' <<<"$G"); BID=$(jq -r '.barriers[0].id' <<<"$G")
R=$(api DELETE "entidades/$EREL"); check "RemoveEntity (UUID de tabla) → 409" "$(code_of "$R")" "409"
R=$(api DELETE "barreras/$BID"); check "RemoveBarrier (UUID de tabla) → 409" "$(code_of "$R")" "409"
check "las filas siguen ahí" "$(SQL "select count(*) from salvia.entity_case where case_id='$IC' and deleted_at is null")/$(SQL "select count(*) from salvia.barrier_v2 where case_id='$IC' and deleted_at is null")" "1/1"
check "GetCase.entitiesReadOnly = true" "$(jq -r '.entitiesReadOnly' <<<"$G")" "true"
check "GetCase sigue devolviendo la entidad (de la tabla)" "$(jq -r '.entities|length' <<<"$G")" "1"
check "GetCase.entityMeta[101] sigue disponible" "$(jq -r '.entityMeta["101"].dueDiligence' <<<"$G")" "true"
check "GetCase.barrierMeta[barrera] disponible" "$(jq -r --arg id "$BID" '.barrierMeta[$id].zone' <<<"$G")" "cabecera"

echo
echo "RESULTADO: $PASS ok / $FAIL fallos   (caso $IC)"
exit $FAIL
