# Gestión Salvia § C/D — Entidades y Barreras — Casos de prueba

Cubre E17/E18. Endpoints: `POST/DELETE /api/v1/victim-cases-v3/:icode/entidades`, `POST/DELETE /api/v1/victim-cases-v3/:icode/barreras`.

## Entidades (§ C)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E1 | GS_14 = "No" (no se ha reportado a ninguna entidad) | — | El formulario de agregar entidad sigue **visible** (hallazgo G3 — no está gateado por GS_14, a diferencia de barreras) — confirmar que efectivamente se puede agregar una entidad en este estado | Validar DOM |
| E2 | Agrega una entidad tipo "acudida", sector "Salud" | — | `POST /entidades` → 200 `{icode}`; en BD, `victim_case_entity_kind='acudida'`, `victim_case_entity_sector` = id de `victim_case_form2_sector_sa` | Request directo + leer BD |
| E3 | Agrega una entidad tipo "activación" sin elegir canal de activación | Kind="activacion" | El campo "Canal de activación" aparece en el formulario (condicional a `Kind==='activacion'`); si se envía sin canal, confirmar qué hace el backend (`ActivationChannel` vacío → `channelID` queda `nil`, no debería reventar) | Validar DOM + request directo |
| E4 | Agrega una entidad tipo "acudida" (canal no debería aparecer) | Kind="acudida" | El campo "Canal de activación" permanece oculto | Validar DOM |
| E5 | Agrega 2 entidades seguidas | — | Ambas aparecen en la tabla, cada una con su propio `icode`; el formulario se limpia (`entidadVacia()`) tras cada guardado | Validar DOM |
| E6 | Elimina una entidad de la lista | Al menos 1 entidad agregada | `DELETE /entidades/:entityIcode` → 200; la fila desaparece de la tabla; en BD, `victim_case_entity_deleted_at` queda con fecha (soft delete, no se borra la fila) | Request directo + leer BD |
| E7 | Agrega una entidad sin seleccionar sector | Sector vacío | El botón "+ Siguiente entidad" está deshabilitado (`:disabled="!nuevaEntidad.SectorCode"`) | Validar DOM |
| E8 | Reabre el Caso tras agregar entidades | Tras E2/E5 | `GET /:icode` devuelve las entidades en `entities`, se repueblan en la tabla | Request directo + validar DOM |

## Barreras (§ D)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| B1 | GS_17 = "No" | — | El formulario de agregar barrera permanece **oculto** (sí está gateado, a diferencia de entidades — contrastar con E1) | Validar DOM |
| B2 | GS_17 = "Sí" | Oculto | Aparece el formulario completo (sector, entidad, fecha, barreras específicas, transversales, descripción, gestión) | Validar DOM |
| B3 | Selecciona sector "Salud" | Sin barreras específicas visibles | Aparece el bloque "Barreras específicas del sector" con las 12 opciones reales de `victim_case_form2_barrier_health` | Validar DOM |
| B4 | Selecciona sector "Justicia" | — | 18 opciones reales de `victim_case_form2_barrier_justice` | Validar DOM |
| B5 | Selecciona sector "Protección" | — | 13 opciones de `victim_case_form2_barrier_protect` | Validar DOM |
| B6 | Selecciona sector "Otras instituciones" (`ot`) | — | **No** aparece bloque de "barreras específicas del sector" (`SECTOR_A_CATEGORIA_BARRERA_V3` no tiene entrada para `ot` — `barreraCategoriaSector()` devuelve `null`) — solo quedan visibles las 4 transversales | Validar DOM — caso que distingue sector con catálogo propio de sector sin catálogo propio |
| B7 | Cambia de sector "Salud" a "Justicia" tras marcar barreras específicas de Salud | Selección previa en `nuevaBarrera.Specific` | El bloque cambia a las opciones de Justicia; verificar si los icode de Salud marcados antes quedan "fantasma" en `nuevaBarrera.Specific` (el array no se limpia al cambiar de sector) | Validar DOM + inspeccionar `nuevaBarrera.Specific` |
| B8 | Marca barreras transversales (institucional, económica, territorial, diferencial) sin ninguna específica de sector | Sector "Otras instituciones" | Se pueden marcar igual, son independientes del sector elegido | Validar DOM |
| B9 | Intenta guardar sin descripción | `nuevaBarrera.Description` vacío | Botón deshabilitado (`:disabled="!nuevaBarrera.SectorCode || !nuevaBarrera.Description"`) | Validar DOM |
| B10 | Guarda una barrera completa | Sector + descripción + al menos 1 específica | `POST /barreras` → 200; en BD, `victim_case_barrier_specific` (JSONB) con los icode marcados, `victim_case_barrier_sector` resuelto | Request directo + leer BD |
| B11 | Marca gestión realizada (ej. "Alerta por barreras") | — | Se guarda en `victim_case_barrier_management` (JSONB) | Leer BD |
| B12 | Elimina una barrera | Al menos 1 agregada | Soft delete (`victim_case_barrier_deleted_at`), desaparece de la tabla | Request directo + leer BD |
| B13 | Reabre el Caso tras agregar barreras | — | Se repueblan en la tabla vía `GET /:icode` → `barriers` | Request directo + validar DOM |

## Nota

E1 vs. B1 es el caso que documenta la inconsistencia G3 del plan — ejecutar ambos y reportar el contraste explícitamente (no es una corrección, es documentar el comportamiento real para que el usuario decida si GS_14 debería gatear igual que GS_17).
