# Guardado por sección y Tríada (Borrador→Iniciado) — Casos de prueba

Cubre E12. Endpoint: `POST /api/v1/victim-cases-v3/:icode/seccion/:modulo`. Mecanismo verificado manualmente contra la BD real antes de escribir estos casos (ver hallazgo H3 en `registro-v3-qa-plan.md`) — los números de caso T1-T6 reproducen exactamente esa verificación para que quede como regresión.

Campos de la tríada (`triadaFieldsV3`, `salvia/service/victim_case_v3_service.go`): CI_1,2,3,4,5,6,7,9 · DP_1,2,5,8,9,11,12,14,17,27,28,29 · HV_1,2,4,6,8,9,11,13,14 · AG_1,2,3. 32 campos en total, CA no aporta ninguno.

## Casos — progresión de la tríada

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| T1 | Guarda CI completo (8 campos) recién creado el Borrador | `victim_case.status='bo'`, ningún guardado previo | `POST /seccion/CI` → `{ok:true, triadaCompleta:false}`; status sigue `bo` | Request directo + leer BD |
| T2 | Guarda DP completo tras T1 | Tras T1 | `triadaCompleta:false` (faltan HV/AG); status sigue `bo` | Request directo + leer BD |
| T3 | Guarda HV completo tras T2 | Tras T2 | `triadaCompleta:false` (falta AG); status sigue `bo` | Request directo + leer BD |
| T4 | Guarda AG completo tras T3 (último módulo de la tríada) | Tras T3 | `triadaCompleta:true`; status cambia a `in` | Request directo + leer BD |
| T5 | Guarda CI con **un solo campo faltante** (ej. omite CI_9) y luego completa los otros 31 | CI/DP/HV/AG completos salvo CI_9 | Mientras falte CI_9, `triadaCompleta` permanece `false` en cada guardado; al enviar CI_9, pasa a `true` | Request directo — probar el límite exacto, no solo "todo o nada" |
| T6 | Guarda un módulo de la tríada ya completa (ej. vuelve a guardar CI con cambios menores) | Caso ya en `in` (tras T4) | `UpdateStatusIfCurrent` solo transiciona si el status actual es `bo` — confirmar que el status **no** se toca (sigue en `in`, o en el estado que tenga) al re-guardar | Leer BD antes y después |
| T7 | Guarda CI con un campo vacío después de haberlo respondido (ej. borra CI_3) | CI_3 ya respondido antes | El campo vuelve a mandarse en el payload con `""` — confirmar si `victim_case_v3_answered_keys` lo sigue contando como "respondido" (la marca es por presencia de la key, no por valor no vacío) y si eso es el comportamiento esperado | Leer BD (`victim_case_v3_answered_keys`) + request directo |
| T8 | Guarda un módulo que no es de la tríada (GS o TZ) | Cualquier estado | `triadaCompleta` siempre `false` en la respuesta, sin reevaluar nada (la función retorna temprano para módulos fuera de CI/DP/CA/HV/AG) | Request directo |
| T9 | CA no aporta a la tríada aunque se guarde completo | CI/DP/HV/AG incompletos, solo CA completo | `triadaCompleta` sigue `false` — CA_1-4 no están en `triadaFieldsV3` | Request directo |

## Casos — persistencia de respuestas multi-valor (enumN)

Mecanismo: `rel_victim_case_form2_enums_victim_case_form2`, sin columna que distinga a qué pregunta pertenece cada relación — el borrado/inserción se acota por categoría (`saveEnumNFields`).

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| T10 | Guarda CI_8 (ajustes razonables) con 3 opciones marcadas | CI_8 sin respuestas previas | 3 filas nuevas en `rel_victim_case_form2_enums_victim_case_form2`, todas de categoría `victim_case_form2_adjustments_gbv` | Leer BD |
| T11 | Vuelve a guardar CI_8 con una selección distinta (2 opciones, 1 repetida de T10) | Tras T10 | Las 3 filas viejas de esa categoría se borran y se insertan las 2 nuevas — total final = 2, no 5 | Leer BD (contar filas antes/después) |
| T12 | Guarda CI_8 (una categoría) y HV_9 (otra categoría) en guardados separados | — | Guardar HV_9 **no borra** las relaciones de CI_8 — el borrado está acotado por categoría, no es un DELETE global de la fila | Leer BD — confirmar que ambas categorías coexisten |
| T13 | Guarda HV_10 con subtipos de 2 tipos distintos (ej. `_fi` y `_se`) en el mismo guardado | HV_9 ya tiene "Física"+"Sexual" marcados | El borrado/inserción de HV_10 cubre **las 7 categorías de subtipo a la vez** (`victimCaseV3EnumNCategories["HV_10"]`) — confirmar que no quedan relaciones huérfanas de una categoría de subtipo no usada en este guardado | Leer BD |

## Caso — reapertura de un Borrador (`GetAnswers` vs. respuestas reales)

Corregido durante la preparación de este plan (mismo hallazgo H3): `GetCase` ahora devuelve también `answeredKeys`, y el frontend (`abrirCaso`, `registro-v3-app.js`) filtra `answers` por esa lista antes de repoblar — los valores dummy ya no deberían aparecer preseleccionados. T14 verifica esa corrección como regresión.

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| T14 | Reabre un Borrador con CI parcialmente guardado (T1) y el resto de módulos nunca tocados | Tras T1 | `GET /:icode` → `answers` trae TODOS los campos (reales + dummy, ej. `AG_1` con el icode de "No", `DP_1` con fecha `1900-01-01`), pero `answeredKeys` solo trae las 8 keys de CI realmente guardadas en T1; en el DOM, al reabrir, solo los campos de CI deben aparecer preseleccionados — AG_1/DP_1/etc. deben verse vacíos/sin seleccionar | Request directo al JSON (comparar `answers` vs `answeredKeys`) + validar DOM tras reabrir |
| T15 | Reabre un Caso ya en `Iniciado` con los 4 módulos de la tríada completos | Tras T4 | Los 32 campos de la tríada aparecen preseleccionados con sus valores reales; los campos de GS/TZ (nunca guardados) siguen sin preselección | Validar DOM |
