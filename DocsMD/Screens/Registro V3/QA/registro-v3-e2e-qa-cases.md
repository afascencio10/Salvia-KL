# Registro V3 — Casos end-to-end

Recorridos completos que ejercitan varios de los casos atómicos de los otros archivos a la vez. Útiles para formalizar como `.spec.ts` de regresión una vez confirmado el flujo paso a paso con `playwright-cli` (criterio de `qa-process.md`).

## E2E-1 — Camino feliz completo: Borrador → Iniciado → Terminado, con remisión, entidad y barrera

| Paso | Acción | Casos atómicos relacionados |
|---|---|---|
| 1 | Carga la pantalla, acepta Habeas Data | C1/C2 |
| 2 | Busca por documento nuevo (sin coincidencias) | C7 |
| 3 | Crea "Nuevo Registro" con documento + teléfono | C9 |
| 4 | Llena CI completo (incluye marcar discapacidad=Sí para activar CI_8/DP_4) | D1, T1 |
| 5 | Llena DP completo (incluye nacionalidad="Extranjera" para activar DP_6, y una etnia="Indígena" para activar DP_13) | D7, D10, T2 |
| 6 | Llena HV completo (marca 2 tipos de violencia para activar HV_10 con unión de catálogos, y "Ámbito laboral" para activar HV_12) | D14, D18, T3 |
| 7 | Llena AG completo (relación="Pareja íntima", para que Tamizaje tome la rama pareja) | T4 (la tríada se completa aquí → status debe pasar a `Iniciado`) |
| 8 | Verifica en el DOM que aparece el aviso de tríada completa y que TZ/GS están desbloqueados | — |
| 9 | Responde Tamizaje: marca 1 de las 5 iniciales en Sí, completa el resto de la rama pareja con un número conocido de "Sí" (ej. exactamente 8, para caer en nivel "moderado" si es feminicidio+pareja) | TZ2, TZ18 (ajustar al umbral real) |
| 10 | Guarda Tamizaje, verifica `riskType`/`riskLevel`/`score` en el DOM | — |
| 11 | En GS § A, responde denuncia fácil y gestiones interinstitucionales | — |
| 12 | En GS § B, activa remisión, marca "Atención Psicosocial" con el obligatorio + 2 criterios que sumen ≥3 pts | R4, R8 |
| 13 | En GS § C, agrega una entidad tipo "acudida" | E2 |
| 14 | En GS § D, activa identifica barreras, selecciona sector "Salud", marca 2 barreras específicas + 1 transversal, agrega | B2, B3, B10 |
| 15 | En GS § E, activa alerta inminente con justificación | A3, A5 |
| 16 | Guarda la sección GS | — |
| 17 | Presiona "Finalizar registro y enrutar a Seguimiento" | F2 |
| 18 | Verifica el párrafo síntesis en la pantalla de Cierre | — |
| 19 | Verifica en BD: `victim_case.status='te'`, todas las tablas relacionadas (`victim_case_form2`, `victim_case_entity`, `victim_case_barrier`, relaciones multi-valor) con los datos esperados | — |

## E2E-2 — Abandona el Borrador a mitad de camino y lo retoma

| Paso | Acción | Casos atómicos relacionados |
|---|---|---|
| 1 | Crea un Borrador nuevo, llena solo CI completo, guarda | T1 |
| 2 | Cierra la pestaña / navega a otra pantalla sin finalizar | — |
| 3 | Vuelve al Cabezote, busca por el mismo documento | C8 (debe aparecer el Borrador existente, no permitir crear uno duplicado) |
| 4 | Abre el Caso encontrado | C11 |
| 5 | Verifica que solo CI aparece preseleccionado — DP/HV/AG/GS/TZ deben verse vacíos, sin valores dummy visibles | T14 |
| 6 | Completa DP, HV, AG | T2, T3, T4 |
| 7 | Verifica que el status pasó a `Iniciado` justo al completar AG, no antes | T4 |

## E2E-3 — Rama no-pareja con esOsigd y nivel de riesgo extremo

| Paso | Acción | Casos atómicos relacionados |
|---|---|---|
| 1 | Crea Borrador, completa la tríada con AG_3="Vecino(a)" (no-pareja) y DP_9="Bisexual" (activa esOsigd) | TZ7, TZ11 |
| 2 | En Tamizaje, responde "Sí" a todas las preguntas disponibles (generales + VPP + las 24 de la rama no-pareja) | TZ22 |
| 3 | Guarda y verifica `riskType='homicidio_abyecto'`, `riskLevel='extremo'` | TZ11, TZ22 |

## Notas generales de ejecución

- Limpiar siempre los datos creados (`victim_case`, `victim_case_form2`, `general_user`, `general_user_profile`, entidades/barreras asociadas) entre corridas de E2E para no contaminar las búsquedas de duplicados de otros casos.
- E2E-1 es el candidato más claro para formalizar como `.spec.ts` de regresión (se va a correr repetidamente); E2E-2 y E2E-3 pueden quedar como exploración con `playwright-cli` hasta que se decida si ameritan convertirse en script.
