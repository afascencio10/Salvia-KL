# Tamizaje — Casos de prueba

Endpoint: `POST /api/v1/victim-cases-v3/:icode/tamizaje`. Lógica en `salvia/service/victim_case_v3_service.go` (`SaveTamizaje`, `computeRiskLevelV3`).

## Grupo 1 — Gating de las 5 preguntas iniciales (D4)

Preguntas iniciales obligatorias (`inicial:true` en `TAMIZAJE_GENERALES_V3`): TZ_1 (aumento violencia física), TZ_2 (arma/objeto), TZ_3 (amenaza de muerte), TZ_4 (asfixia/estrangulamiento), TZ_7 (acceso a armas). TZ_5/TZ_6 son generales pero NO iniciales (solo se revelan).

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| TZ1 | Responde "No" a las 5 iniciales | Módulo TZ recién abierto | El resto del Tamizaje (TZ_5, TZ_6, TZ_8/VPP, rama pareja/no-pareja de 24 preguntas) permanece oculto; `tzRevelarResto=false` | Validar DOM |
| TZ2 | Responde "Sí" a **una sola** de las 5 iniciales (ej. solo TZ_3) | — | Se revela el resto inmediatamente (`tzToggle` hace `some()` sobre las 5, no exige todas) | Validar DOM |
| TZ3 | Responde "Sí" a TZ_3 y luego lo cambia a "No", sin haber marcado ninguna otra inicial en Sí | Tras TZ2 | El resto vuelve a ocultarse (`tzRevelarResto` se recalcula en cada `tzToggle`) | Validar DOM |
| TZ4 | Responde "Sí" a TZ_5 (no es inicial) sin tocar ninguna de las 5 iniciales | — | El resto permanece oculto — TZ_5/TZ_6 no disparan el gating aunque sean generales | Validar DOM |

## Grupo 2 — Rama pareja / no-pareja (detectada automáticamente desde AG_3)

`calcularWasPartner()` revisa si la etiqueta seleccionada en AG_3 contiene "pareja".

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| TZ5 | AG_3 = "Pareja íntima" | Módulo AG guardado antes de llegar a TZ | Tamizaje muestra las 24 preguntas de la rama pareja (`TAMIZAJE_PAREJA_V3`) | Validar DOM |
| TZ6 | AG_3 = "Expareja o antigua pareja" | — | También activa la rama pareja (contiene "pareja") | Validar DOM |
| TZ7 | AG_3 = "Hermano(a)" | — | Activa la rama no-pareja (`TAMIZAJE_NO_PAREJA_V3`) | Validar DOM |
| TZ8 | AG_3 no se ha respondido todavía | Llega a TZ sin haber guardado AG | `etiquetaSeleccionada('AG_3')` devuelve `''`, `wasPartner=false` por defecto → rama no-pareja | Validar DOM + confirmar que no revienta sin AG_3 |
| TZ9 | Cambia AG_3 de "Pareja íntima" a "Hermano(a)" después de haber respondido varias preguntas de la rama pareja | Rama pareja con respuestas | Tamizaje cambia a la rama no-pareja; las respuestas de la rama pareja (`TZP_1..24`) quedan en `answers` pero ya no se muestran ni se envían (el payload de `guardarTamizaje` construye según la rama **actual**) | Validar DOM + inspeccionar el payload enviado |

## Grupo 3 — Tipo de riesgo (D3: por identidad/orientación, NO por la pregunta VPP)

`calcularEsOsigd()`: `true` si DP_9 (orientación) no es "Heterosexual", o si DP_8 (identidad) contiene "trans" o "no binaria".

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| TZ10 | DP_9 = "Heterosexual", DP_8 = "Mujer" | — | `esOsigd=false` → `riskType='feminicidio'` | Request directo, revisar `riskType` en la respuesta |
| TZ11 | DP_9 = "Gay" | DP_8 cualquiera | `esOsigd=true` → `riskType='homicidio_abyecto'` | Request directo |
| TZ12 | DP_9 = "Heterosexual", DP_8 = "Hombre Trans" | — | `esOsigd=true` por la identidad, aunque la orientación sea heterosexual | Request directo — confirma que ambos campos se evalúan, no solo DP_9 |
| TZ13 | Responde "Sí" a TZ_8 (pregunta VPP) con DP_9="Heterosexual" y DP_8="Mujer" | — | `riskType` sigue siendo `'feminicidio'` — TZ_8 **no** decide el tipo, solo suma al `score` (D3) | Request directo — este es el caso que prueba explícitamente que VPP no es el criterio de tipo |

## Grupo 4 — Empate / override por subtipo (B6) — actualmente inalcanzable (G1)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| TZ14 | Intenta marcar "Tentativa Feminicidio" y "Tentativa Homicidio abyecto o fútil" como subtipos de HV para forzar `tieFlag=true` | — | No se puede — esas opciones no existen en el catálogo real de subtipo (ver G1). El payload que llega a `/tamizaje` nunca tendrá `subtipoFeminicidio=true` ni `subtipoHomicidioAbyecto=true` desde la UI real. Documentar como no ejecutable hasta que se resuelva G1 | No aplica (caso de gap, no de ejecución) |
| TZ15 | Llama directamente al endpoint `POST /tamizaje` con `subtipoFeminicidio:true, subtipoHomicidioAbyecto:true` en el body (bypasseando la UI) | — | El backend sí soporta el empate: `riskType='feminicidio'`, `tieFlag=true` (lógica de `SaveTamizaje` en el service, confirmada por lectura de código) | Request directo al endpoint, sin pasar por la UI — prueba el backend aislado del gap de catálogo del frontend |

## Grupo 5 — Score y umbrales de nivel de riesgo

Umbrales (`computeRiskLevelV3`): feminicidio+pareja → ≤4 bajo, ≤8 moderado, ≤15 alto, resto extremo. feminicidio+no-pareja → ≤2 bajo, ≤5 moderado, ≤8 alto, resto extremo. homicidio_abyecto (ambas) → ≤5 bajo, ≤11 moderado, ≤16 alto, resto extremo. El score cuenta cuántas respuestas de TODO el payload (7 generales + VPP + 24 de la rama) son `"true"`.

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| TZ16 | 0 respuestas en "Sí" | feminicidio + pareja | `score=0`, `riskLevel='bajo'` | Request directo |
| TZ17 | Exactamente 4 "Sí" | feminicidio + pareja | `riskLevel='bajo'` (límite exacto del umbral) | Request directo — caso de frontera |
| TZ18 | Exactamente 5 "Sí" | feminicidio + pareja | `riskLevel='moderado'` (primer valor que cruza el umbral de "bajo") | Request directo — caso de frontera, par con TZ17 |
| TZ19 | Exactamente 2 "Sí" | feminicidio + **no**-pareja | `riskLevel='bajo'` | Request directo |
| TZ20 | Exactamente 3 "Sí" | feminicidio + no-pareja | `riskLevel='moderado'` (el umbral de no-pareja es más bajo que el de pareja — confirma que la rama cambia el umbral, no solo el catálogo de preguntas) | Request directo — caso de frontera |
| TZ21 | 20 "Sí" de 32 posibles | homicidio_abyecto (cualquier rama) | `riskLevel='alto'` (≤16 sería alto, 20>16 → `extremo`; recalcular con el valor exacto al ejecutar, documentar aquí el score real usado) | Request directo |
| TZ22 | Todas las respuestas en "Sí" (32/32) | cualquier tipo/rama | `riskLevel='extremo'` en los tres escenarios de tipo/rama | Request directo, repetir para feminicidio+pareja, feminicidio+no-pareja y homicidio_abyecto |

## Nota para Playwright

Los casos de frontera (TZ17/TZ18/TZ19/TZ20) dependen de marcar un número exacto de preguntas en "Sí" — más fácil de ejecutar con `playwright-cli` llamando directo al endpoint con un payload armado a mano (ver Fase 2/3 de `qa-process.md`) que clickeando 4-5 radios específicos en la UI; una vez confirmado el backend, sí vale la pena un paso por UI para 1-2 casos representativos (ej. TZ17 completo por clicks) para confirmar que el payload que arma el frontend coincide con lo esperado.
