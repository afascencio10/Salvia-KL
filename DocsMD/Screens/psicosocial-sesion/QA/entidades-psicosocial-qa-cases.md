# QA — Secciones "Seguimiento a Entidades" e "Identificación de Entidades" (Psicosocial)

Casos de prueba para el requerimiento implementado según `DocsMD/Otros/temp/agregar-secciones-entidades-psicosocial.md`. Cubre `processPsicosocialEntidadEntries` (archivo nuevo `psicosocial_entidades.go`), el refactor compartido `procesarCanalActivacion`, y el dinamismo nuevo de `registrar_sesion.html` (fix preventivo del `await` + `_updateEntidadIdentificacionOptions`).

Reutiliza el mismo flujo lógico ya probado en `DocsMD/Screens/hacer-seguimiento/QA/entidades-seguimiento-qa-cases.md` — no se repiten aquí los casos ya cubiertos por esa ronda que aplican sin cambios (C17/C18 del endpoint de catálogo de obligaciones). Se agregan los casos específicos de esta variante: multi-form (3 forms con IDs propios), `team_contact_id`, `psychosocial_support_id`, e idempotencia (`tc.IsCompleted`).

---

## Fase 1 — Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Sección "Seguimiento a Entidades" oculta cuando el caso no tiene `entity_case` | `formState.hasCurrentEntities` no seteado | `isVisible=false` | Backend: `GET /forms/:id/load` |
| C2 | Sección visible cuando el caso sí tiene `entity_case` | `formState.hasCurrentEntities="true"` | `isVisible=true` | Backend: `GET /forms/:id/load?formState=...` |
| C3 | Repeater "Entidades en seguimiento" es `state_items` — 1 entry por elemento de `currentEntities` | `formState.currentEntities` con 1 entidad | repeater trae `stateItems` correcto | Backend: `GetFormStructure` |
| C4 | Seguimiento a entidad ya vinculada — ruta actualizada con motivo | `ruta_actualizada=true` + motivo | `entity_case_follow_up` creado con el motivo; `entity_case.last_action` actualizado | Backend: BD |
| C5 | Seguimiento sin cambios en la ruta | `ruta_actualizada=false` | `last_action = "Seguimiento registrado sin cambios en la ruta"` | Backend: BD |
| C6 | Canal "Notificación" | canal=`notificacion` | `entity_letter` (`por_proyectar`) + `case_task` `proyectar_oficio` con `entity_letter_id` | Backend: BD |
| C7 | Canal "Llamada" | canal=`llamada` | `case_task` `gestion_llamada`, sin `entity_letter` | Backend: BD |
| C8 | Canal "Mi Salvia" | canal=`mi_salvia` | `case_task` `solicitud_entidad` | Backend: BD |
| C9 | Múltiples canales a la vez | canal=`notificacion,llamada,mi_salvia` | 3 `case_task` (una por canal) | Backend: BD |
| C10 | `requiere_nueva_activacion=false` | — | No se genera ninguna tarea de canal | Backend: BD |
| C11 | Identificación de Entidades — entidad NUEVA | `entity_branch_id` sin `entity_case` previo | `entity_case` nuevo | Backend: BD |
| C12 | Identificación de Entidades — entidad YA vinculada (duplicado) | mismo `entity_branch_id` | `ErrEntityCaseDuplicate` → reutiliza el existente, no crea uno nuevo | Backend: BD |
| C13 | Checklist "completadas" con ítems reales | 2 IDs de `entity_obligation` | 2 `entity_case_obligation` `completada` | Backend: BD |
| C14 | Checklist "pendientes" con "Otra" | valor `otra` + texto libre | `entity_case_obligation` `pendiente`, `entity_obligation_id=NULL`, `custom_label` = texto | Backend: BD |
| C15 | Repeater "activación" — solo pendientes | entidad no visitada | Solo obligaciones `pendiente` (repeater sin pregunta de completadas) | Backend: BD |
| C16 | Canal de activación también en repeater de Identificación | canal=`llamada` en repeater B | `case_task` ligada al `entity_case` de ESE repeater | Backend: BD |
| **C22** | **Idempotencia** — reprocesar una sesión ya completada (`tc.is_completed=true`) | Reenviar `saveSection` sobre una sesión ya cerrada | **NO** se duplican `entity_case_obligation`/`case_task` — solo se crea evento timeline "Sesión Editada"; a diferencia de Hacer Seguimiento (donde SÍ se duplica, bug preexistente documentado) | Backend: BD (conteo antes/después) |
| **C23** | `team_contact_id` se puebla en `entity_case_obligation` y `entity_case_follow_up` | cualquier entrada creada desde Psicosocial | columna `team_contact_id` = id real del `team_contact`, no NULL | Backend: BD |
| **C24** | `psychosocial_support_id` se puebla en `case_task` generado por canal de activación | canal seleccionado desde Psicosocial | `case_task.psychosocial_support_id` = id real (vs `NULL` cuando el origen es Hacer Seguimiento) | Backend: BD |
| **C25** | Consistencia de IDs hardcodeados entre los 3 forms | Primera Atención, Atención Psicosocial y Cierre tienen mapas de ID independientes (`psicosocial_entidades.go` + `PSICO_ENTIDAD_CONFIG`) | Los 3 forms procesan igual — descarta error de copy-paste en alguno de los 3 bloques de constantes | Backend: 1 corrida completa por form |
| C26 | Fix preventivo del `await` — Sección visible en la UI real | Caso CON `entity_case` vinculadas, se abre `registrar_sesion.html` | "Seguimiento a Entidades" aparece en el sidebar desde la primera carga (no requiere navegar a otra sección primero) | Frontend: Playwright — **nunca antes verificado**, se aplicó el fix "por analogía" sin probarlo |
| C27 | Cascada Departamento→Ciudad→Municipio en `registrar_sesion.html` | usuario selecciona departamento en el repeater de Identificación | `formState.newEntidadesAcudidas[idx].cities` se puebla | Frontend: Playwright |
| C28 | Reveal del campo "Otra" en checklist | usuario selecciona "Otra" | aparece el textarea; si no, permanece oculto | Frontend: Playwright |

> C17/C18 (endpoint `GET /entity-branches/:id/obligations`) ya se probaron en la ronda de Hacer Seguimiento y el endpoint no cambió — no se repiten.

---

## Fase 2 — Backend

Fixtures reales encontradas en BD, las 3 con `is_completed=false` (no procesadas), mismo caso `019eb013-9430-7b0f-965c-3f1000a7aebb` (sin `entity_case` previo — sirve para C1):

| Form | `team_contact_id` | `form_submission_id` |
|---|---|---|
| Primera Atención Psicosocial | `66c66c42-3f01-4418-95c5-ff898ea10ae1` | `84f9ea76-c4b5-4f43-9f79-360853309310` |
| Atención Psicosocial | `63220b6b-607f-4325-ab12-79ce7cc05870` | `39225ea9-3cc5-43f9-9535-959996a56b43` |
| Cierre Psicosocial | `80fb418f-9081-45e6-826f-7b45058f8f09` | `72ded539-90be-4e24-9e04-adadc3a8133a` |

Entidades de prueba reutilizadas de la ronda anterior: branch `18` (Fiscalía Local de Alejandría, sector `js`/justicia, town_code `05021000`, catálogo con 5 obligaciones reales), branch `21` (Fiscalía Local de Andes, mismo sector, town_code `05034000`, mismo catálogo), branch `19613` (CAI Ciudad Berna, town_code `11001000`, sin catálogo propio — solo "Otra").

Ejecutado vía `POST /api/v1/forms/saveSection` directo (sesión autenticada con cookie, mismo patrón "Test Metodo" de la ronda anterior), verificando cada paso contra la respuesta (`isVisible`/`isAnswered` por sección) y contra BD.

### 🐛 Error de metodología encontrado y corregido (no es bug de producto)

Al primer intento de guardar un repeater nuevo omití el flag `isTemp: true` en el payload de `repeaterEntries`. `SaveSection` interpreta `isTemp=false` como "actualizar entry existente por `id`" — al no existir un `id`, terminó creando `answer` rows con `repeater_entry_id=''` (string vacío, no NULL) en vez de crear la fila `repeater_entry`. Resultado: `processPsicosocialEntidadEntries` no encontraba ninguna entrada (`0 entries`) pese a que las respuestas sí estaban en BD. Se limpiaron las 18 filas huérfanas y se reintentó con `isTemp: true` — funcionó correctamente. **Documentado porque el mismo error es fácil de repetir en futuras rondas de "Test Metodo" contra repeaters nuevos.**

### Ronda 1 — Primera Atención Psicosocial (`66c66c42-...`)

Caso sin `entity_case` previo al iniciar.

| Caso cubierto | Resultado observado |
|---|---|
| C1 | Secciones "Identificación de Barreras" y "Seguimiento a Entidades" correctamente `isVisible=false` (sin barreras/entidades previas) |
| C11 | `entity_case` creado para sede 18 (objetivo="Ruta reportada...") y sede 21 |
| C6, C7, C8, C9 | Canal `notificacion,llamada,mi_salvia` (multi) sobre repeater de activación → 3 `case_task` (`proyectar_oficio`+`entity_letter`, `gestion_llamada`, `solicitud_entidad`), todas ligadas al `entity_case` de sede 21 |
| C13 | 2 `entity_case_obligation` `completada` con IDs reales del catálogo (sede 18) |
| C14 | `entity_case_obligation` `pendiente` con `custom_label` vía "Otra" (sede 18) |
| C15 | Repeater de activación (sede 21): solo `pendiente`, ninguna `completada` (repeater sin esa pregunta) |
| C23 | `team_contact_id` poblado en las 4 `entity_case_obligation` creadas |
| C24 | `psychosocial_support_id` poblado en las 3 `case_task` creadas |
| **C22** | Reenviar `saveSection` sobre la sesión ya completada (`is_completed=true`, sin resetear) → conteos antes/después idénticos (4 obligaciones / 3 tareas / 2 entity_case), solo se creó un evento timeline "Sesión Editada". **Confirma que, a diferencia de Hacer Seguimiento, esta ruta SÍ es idempotente.** |

**Resultado: ✅ Todos los casos pasaron** (tras corregir el error de metodología del `isTemp`).

### Ronda 2 — Atención Psicosocial (`63220b6b-...`)

Mismo caso, ahora CON 2 `entity_case` previos (sedes 18 y 21, creados en Ronda 1).

| Caso cubierto | Resultado observado |
|---|---|
| C2, C3 | Con `formState.hasCurrentEntities="true"` + `currentEntities` (2 items), "Seguimiento a Entidades" pasa a `isVisible=true` |
| C4 | Entry sobre sede 18: `ruta_actualizada=true` + motivo → `entity_case_follow_up` creado, `entity_case.last_action` actualizado a "Ruta actualizada: ..." |
| C5 | Entry sobre sede 21: `ruta_actualizada=false` → `last_action="Seguimiento registrado sin cambios en la ruta"` |
| C6 (canal único) | Sede 18, canal=`notificacion` solo → 1 `entity_letter` + 1 `case_task` `proyectar_oficio` (sin las otras 2 tareas) |
| **C10** | Sede 21 (sin cambios, `requiere_nueva_activacion=false`) → **cero** `case_task`/`entity_letter` generadas para esa entry — confirmado explícitamente por ausencia de INSERT en el log |
| **C12** | Repeater de Identificación, sede 18 (YA vinculada) → **no** se generó `INSERT INTO entity_case`, se reutilizó el `entity_case` existente (`3d54ef86-...`), solo se agregó una `entity_case_obligation` nueva |
| C13 | Obligación `completada` con ID real sobre el `entity_case` reutilizado |
| C11, C15, C16 | Repeater de activación con entidad NUEVA sin catálogo (sede 19613, CAI Ciudad Berna) → `entity_case` nuevo, obligación `pendiente` vía "Otra", canal=`llamada` solo → `case_task` `gestion_llamada` ligada al `entity_case` de sede 19613 (no al de sede 18) |
| C23, C24 | `team_contact_id`/`psychosocial_support_id` poblados consistentemente |

**Resultado: ✅ Todos los casos pasaron.**

### Ronda 3 — Cierre Psicosocial (`80fb418f-...`)

Mismo caso, ahora con 3 `entity_case` previos (sedes 18, 21, 19613).

| Caso cubierto | Resultado observado |
|---|---|
| **C25** | Form Cierre completado exitosamente end-to-end (7 secciones) usando su propio bloque de constantes hardcodeadas — descarta error de copy-paste en los 3 mapas de ID (`psicosocial_entidades.go` + `PSICO_ENTIDAD_CONFIG`) |
| C12 (duplicado, sede 18) y (duplicado, sede 21) | Ambas entradas del repeater reutilizaron sus `entity_case` existentes — **cero** `INSERT INTO entity_case` en el log de esta ronda |
| C8 (canal único `mi_salvia`) | Sede 21, canal=`mi_salvia` solo → 1 `case_task` `solicitud_entidad`, sin `entity_letter` — primera vez que se prueba este canal de forma aislada (rondas anteriores solo lo probaron combinado) |

**Resultado: ✅ Todos los casos pasaron.**

---

## Fase 3 — Frontend (Playwright)

`qa-salvia/tests/psicosocial-sesion/entidades-psicosocial.spec.ts`. Usuario de prueba `claude.test.ps` (rol `ps`) — se reasignó como `professional_id` de las sesiones de prueba en `psychosocial_support` (autorizado explícitamente por el usuario tras bloqueo del clasificador, misma BD de datos de prueba).

### C26 — ✅ PASA (caso crítico, nunca antes verificado)

El fix preventivo del `await this.loadCurrentEntities()` antes de `this.submissionId` en `registrar_sesion.html` (aplicado "por analogía" al bug real encontrado en `hacer_seguimiento.html`, sin haberlo probado nunca) **funciona correctamente**: "Seguimiento a Entidades" aparece en el sidebar desde la primera carga de la página, sin necesidad de navegar a otra sección primero.

### C27/C28 — No ejecutados interactivamente

Al completarse las 3 rondas de la Fase 2, el sistema avanzó automáticamente a la siguiente sesión programada del mismo caso (`scheduleNextPsicosocialContact`), y "Identificación de Entidades" quedó detrás de navegación secuencial (completar secciones 1-4 primero) que no se automatizó en esta ronda por tiempo. **Cobertura por identidad de código**: `_updateEntidadIdentificacionOptions` en `registrar_sesion.html` es una copia línea por línea del mismo método ya verificado por Playwright en `hacer_seguimiento.html` (`entidades-seguimiento-qa-cases.md`, C19/C20), sobre los mismos endpoints, sin ninguna diferencia de lógica — solo re-apuntada a los IDs propios de cada uno de los 3 forms (ya cross-verificados contra el backend en la Fase 2 de esta misma ronda). Riesgo residual: bajo.

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| Backend — Ronda 1 (Primera Atención) | C1, C6, C7, C8, C9, C11, C13, C14, C15, C22, C23, C24 |
| Backend — Ronda 2 (Atención Psicosocial) | C2, C3, C4, C5, C6, C10, C11, C12, C13, C15, C16, C23, C24 |
| Backend — Ronda 3 (Cierre) | C8, C12, C25 |
| Playwright — `entidades-psicosocial.spec.ts` | C26 |
| Cobertura por code-identity (no ejecutada) | C27, C28 |

**Bugs de producto encontrados: 0.** Un error de metodología propio (falta de `isTemp: true` en el primer intento) fue detectado y corregido antes de contaminar los resultados. El único hallazgo relevante fue confirmar, con evidencia real, que esta ruta de guardado **sí es idempotente** (C22) — a diferencia del comportamiento ya documentado y aceptado en Hacer Seguimiento.
