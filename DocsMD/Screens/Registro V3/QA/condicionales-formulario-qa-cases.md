# Formulario — Condicionales (CI/DP/HV/AG) — Casos de prueba

Cubre E10/E11 (navegación y actualización de respuestas). Esta es la cobertura principal pedida: cada campo condicional del catálogo (`frontend/js/components/registro-v3-campos.js`) tiene al menos un caso positivo, uno negativo y, donde aplica, un caso de transición (activar→desactivar la condición). Mecanismos de condicional usados en el código (`registro-v3-app.js`, método `mostrarCampo`):

- `condicional` — compara contra un booleano literal (`CI_7 === 'true'`)
- `condicionalEnumCode` — compara el `code` real del enum seleccionado (no el texto)
- `condicionalEnumLabel` — compara si la etiqueta seleccionada **contiene** un texto
- `subtipoDe` — visible si el campo padre (multiselección) tiene ≥1 opción marcada
- pseudo-condiciones (`__gestacionAplica`, `__haySubtipo`, `__ambitoLaboral`, `__agresorConocido`) — lógica ad-hoc en `pseudoCondicion()`

## Grupo 1 — CI_7 (discapacidad) → CI_8 (ajustes razonables) y DP_4 (Ley 1996)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D1 | Marca CI_7 = Sí | CI_8 y DP_4 ocultos (CI_7 sin responder o en No) | CI_8 (checklist, categoría `victim_case_form2_adjustments_gbv`, 12 opciones reales) y DP_4 (categoría `victim_case_form2_law_1996`, 4 opciones) se muestran | Validar DOM |
| D2 | Marca CI_7 = No, habiendo marcado antes CI_8 | CI_7=Sí, CI_8 con 2 opciones marcadas | CI_8/DP_4 se ocultan. **Verificar si el valor queda en `answers` aunque oculto** (el código no limpia `answers['CI_8']` al ocultar — confirmar si eso es el comportamiento esperado o si debería limpiarse antes de guardar) | Validar DOM + inspeccionar `this.answers` vía consola |
| D3 | Guarda la sección CI con CI_7=No pero CI_8 con selección residual de D2 | Tras D2, presiona "Guardar sección" | El payload enviado a `POST /seccion/CI` — confirmar si incluye `CI_8` con la selección vieja (el filtro de guardado es por prefijo de módulo, no por visibilidad actual) | Leer el request real (Playwright network) + comparar contra BD (`rel_victim_case_form2_enums_victim_case_form2`) |
| D4 | Nunca marca CI_7 | Caso nuevo | CI_8/DP_4 permanecen ocultos todo el módulo | Validar DOM |

## Grupo 2 — DP_2 (dificultades) → DP_3 (actividades)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D5 | Marca DP_2 = Sí | DP_3 oculto | DP_3 aparece con las 9 opciones reales de `victim_case_form2_activities_unable_to_perform` | Validar DOM |
| D6 | Marca DP_2 = No | DP_3 visible con selección previa | DP_3 se oculta | Validar DOM |

## Grupo 3 — DP_5 (nacionalidad) → DP_6 (país específico) — por **código**, no por texto

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D7 | Selecciona DP_5 = "Extranjera" (code `ex`) | DP_6 oculto | DP_6 aparece, categoría `victim_case_form2_specified_nationality` (249 opciones reales — catálogo de países) | Validar DOM |
| D8 | Selecciona DP_5 = "Colombiana" (code `co`) | — | DP_6 permanece oculto | Validar DOM |
| D9 | Selecciona DP_5 = "Apátrida" (code `ap`) | — | DP_6 permanece oculto (solo se activa por `ex`, no por cualquier valor distinto de `co`) | Validar DOM — este es el caso que distingue "por código exacto" de "por texto" |

## Grupo 4 — DP_12 (grupo étnico) → DP_13 (pueblo indígena) — por **texto de etiqueta**

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D10 | Selecciona DP_12 = "Indígena" | DP_13 oculto | DP_13 aparece, categoría `victim_case_form2_indigenous_people` (122 opciones reales) | Validar DOM |
| D11 | Selecciona DP_12 = "Afrocolombiana(o)" | — | DP_13 permanece oculto | Validar DOM |
| D12 | Selecciona DP_12 = "Pueblo Rrom / Gitano" | — | DP_13 permanece oculto (el texto no contiene "Indígena" aunque sea otra etnia) | Validar DOM — confirma que el match es por substring exacto, no por "es una etnia distinta de afro" |

## Grupo 5 — HV_9 (tipo de violencia, múltiple) → HV_10 (subtipo, dinámico) — el condicional más complejo

`HV_10` no tiene una sola categoría fija: para cada tipo marcado en HV_9, se une el catálogo de `victim_case_form2_subtype_violence_experienced_{código}` correspondiente (`fi`/`pl`/`po`/`ps`/`re`/`se`/`vi`).

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D13 | Marca un solo tipo ("Física") | HV_10 oculto | HV_10 aparece con exactamente las 7 opciones de `_fi` (Ataque con agente químico, Golpizas, Quemadura, Sofocamiento, Agresión con objetos, Tortura, Secuestro) | Validar DOM, contar opciones |
| D14 | Marca dos tipos ("Física" + "Sexual") | Partiendo de D13 | HV_10 ahora muestra la **unión** de `_fi` (7) + `_se` (5) = 12 opciones, sin duplicados | Validar DOM, contar opciones |
| D15 | Desmarca "Física" dejando solo "Sexual" | Partiendo de D14, con alguna opción de `_fi` ya marcada en HV_10 | HV_10 muestra solo las 5 de `_se`. Verificar si las marcas que eran de `_fi` quedan "fantasma" en `answers['HV_10']` (ya no se muestran como checkbox pero el icode sigue en el CSV guardado) | Validar DOM + inspeccionar `this.answers['HV_10']` |
| D16 | Desmarca todos los tipos | Partiendo de cualquier estado con tipos marcados | HV_10 se oculta completamente (`subtipoDe` exige ≥1 marcado en HV_9) | Validar DOM |
| D17 | Intenta registrar "Tentativa Feminicidio" como subtipo de violencia Física | — | **No existe esa opción en el catálogo real** (ver hallazgo G1 del plan) — confirmar que efectivamente no aparece ninguna opción con ese texto en ninguna de las 7 categorías de subtipo | Validar DOM (ausencia) — este caso documenta el gap G1, no se espera que pase, es para no asumir que existe |

## Grupo 6 — HV_11 (ámbito, múltiple, máx. 2) → HV_12 (sector laboral) — pseudo-condición por texto

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D18 | Marca "Ámbito laboral" en HV_11 | HV_12 oculto | HV_12 aparece (categoría `victim_case_form2_workplace_sector_occurrence`, 4 opciones) | Validar DOM |
| D19 | Marca un 3er ámbito tras ya tener 2 marcados | HV_11 con 2 opciones marcadas (tope `max:2`) | El 3er click no agrega nada — `toggleMarcado` respeta el `max` | Validar DOM (sigue con 2 marcadas) |
| D20 | Desmarca "Ámbito laboral" dejando otro ámbito marcado | Partiendo de D18 | HV_12 se oculta | Validar DOM |
| D21 | Marca "Ámbito educativo" + "Ámbito de la salud" (sin laboral) | — | HV_12 permanece oculto en todo momento | Validar DOM |

## Grupo 7 — AG_2 (proximidad del agresor) → AG_5 (depende económicamente) — pseudo-condición con ambigüedad conocida (G4)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D22 | Selecciona AG_2 = "Persona conocida con la que convive" | AG_5 oculto | AG_5 aparece | Validar DOM |
| D23 | Selecciona AG_2 = "Persona conocida con la que no convive" | — | AG_5 aparece (el texto no contiene "desconocid") | Validar DOM |
| D24 | Selecciona AG_2 = "Persona desconocida" | — | AG_5 se oculta (contiene "desconocid") | Validar DOM |
| D25 | Selecciona AG_2 = "No se puede determinar / no informa" | — | **AG_5 queda visible** (el texto no contiene "desconocid", aunque conceptualmente tampoco se conoce la relación) — confirmar con el equipo de producto si este es el comportamiento deseado (G4) | Validar DOM — reportar el resultado real, no asumir que es incorrecto |

## Grupo 8 — DP_26 (gestación) — condición simplificada a "siempre visible"

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| D26 | Abre el módulo DP con cualquier combinación de DP_8/DP_9/DP_10 | — | DP_26 siempre visible (`__gestacionAplica` devuelve `true` sin condición real — simplificación documentada en el código) | Validar DOM |

## Resumen de cobertura esperada

Cada condicional del catálogo queda cubierto por al menos un caso positivo y uno negativo, salvo donde se indica explícitamente que es un gap conocido (D17, G1) o una ambigüedad a validar con producto (D25, G4) en vez de un defecto claro.
