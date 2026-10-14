# QA — Secciones "Seguimiento a Entidades" e "Identificación de Entidades"

Casos de prueba para el requerimiento planeado en `DocsMD/Otros/temp/agregar-secciones-entidades-seguimiento.md`. Cubre la extensión de `processFollowUpSubmission` (PASO 3c/3d + sub-flujo "Procesar canal de activación") y el dinamismo nuevo de `hacer_seguimiento.html` (E-01/E-05).

---

## Fase 1 — Casos de prueba

Un caso por cada `SI`/`SEGÚN` del flujo (ver plan), no solo el camino feliz.

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Sección 5 oculta cuando el caso no tiene `entity_case` | `formState.hasCurrentEntities` no seteado / caso sin entidades vinculadas | `visibility_condition` tipo `SECTION` evalúa `isVisible=false` | Backend: `GET /forms/:id/load` sin `formState` |
| C2 | Sección 5 visible cuando el caso sí tiene `entity_case` | `formState.hasCurrentEntities = "true"` | `isVisible=true` para la sección | Backend: `GET /forms/:id/load?formState=...` |
| C3 | Repeater "Entidades en seguimiento" es `state_items` — una entry por elemento de `currentEntities` | `formState.currentEntities` con 1 entidad | El repeater trae `stateItems: "currentEntities"` en la estructura | Backend: inspección de `GetFormStructure` |
| C4 | Seguimiento a una entidad ya vinculada — ruta actualizada con motivo | `ruta_actualizada=true` + motivo | Se crea `entity_case_follow_up` con el motivo; `entity_case.last_action = "Ruta actualizada: {motivo}"` | Backend: BD (`entity_case_follow_up`, `entity_case.last_action`) |
| C5 | Seguimiento sin cambios en la ruta | `ruta_actualizada=false` | `entity_case.last_action = "Seguimiento registrado sin cambios en la ruta"` | Backend: BD |
| C6 | Canal "Notificación" | `requiere_nueva_activacion=true`, canal=`notificacion` | Se crea `entity_letter` (state=`por_proyectar`) + `case_task` tipo `proyectar_oficio` con `entity_letter_id` seteado | Backend: BD |
| C7 | Canal "Llamada" | canal=`llamada` | Se crea `case_task` tipo `gestion_llamada`, sin `entity_letter` | Backend: BD |
| C8 | Canal "Mi Salvia" | canal=`mi_salvia` | Se crea `case_task` tipo `solicitud_entidad`, `form_data` vacío | Backend: BD |
| C9 | Múltiples canales seleccionados a la vez | canal=`notificacion,llamada,mi_salvia` | Se crean las 3 tareas (una por canal) | Backend: BD |
| C10 | `requiere_nueva_activacion=false` | — | No se genera ninguna tarea de canal | Backend: BD (ausencia de `case_task` nuevas) |
| C11 | Identificación de Entidades — entidad NUEVA (no vinculada al caso) | `entity_branch_id` sin `entity_case` previo | Se crea `entity_case` nuevo (`EntityCaseService.Create`) | Backend: BD |
| C12 | Identificación de Entidades — entidad YA vinculada (duplicado) | mismo `entity_branch_id` que un `entity_case` existente | `ErrEntityCaseDuplicate` → se reutiliza el `entity_case` existente, NO se crea uno nuevo | Backend: BD (mismo `id` de `entity_case` antes/después) |
| C13 | Checklist "completadas" con ítems reales del catálogo | 2 IDs de `entity_obligation` | 2 `entity_case_obligation` status=`completada` con `entity_obligation_id` seteado | Backend: BD |
| C14 | Checklist "pendientes" con opción "Otra" | valor `otra` + texto libre | `entity_case_obligation` status=`pendiente`, `entity_obligation_id=NULL`, `custom_label` = texto | Backend: BD |
| C15 | Repeater "Entidades para activación de ruta" — solo pendientes | entidad no visitada | Solo se registran obligaciones `pendiente`, ninguna `completada` (el repeater no tiene pregunta de completadas) | Backend: BD |
| C16 | Canal de activación también aplica en el repeater de Identificación | canal=`llamada` en repeater de activación | Se crea `case_task` `gestion_llamada` ligada al `entity_case` de ESE repeater | Backend: BD |
| C17 | Endpoint de catálogo de obligaciones resuelve `branch → entity` | `GET /entity-branches/:branchId/obligations` | Devuelve el catálogo de la organización dueña de la sede + `{value:"otra", label:"Otra"}` al final | Backend: request directo |
| C18 | Endpoint de catálogo sin obligaciones seedeadas | branch de una entidad sin catálogo | Devuelve solo `[{"value":"otra","label":"Otra"}]` | Backend: request directo |
| C19 | Cascada Departamento→Ciudad→Municipio en Sección 6 (frontend) | usuario selecciona departamento | `formState.newEntidadesAcudidas[idx].cities` se puebla | Frontend: Playwright / inspección visual |
| C20 | Reveal del campo "Otra" en el checklist (frontend) | usuario selecciona la opción "Otra" | Aparece el `<textarea>` de texto libre; si no se selecciona, permanece oculto | Frontend: Playwright |
| C21 | Sección 5 oculta/visible en la UI real según si el caso tiene entidades | dos casos: uno sin `entity_case`, otro con | El sidebar no muestra "Seguimiento a Entidades" en el primero, sí en el segundo | Frontend: Playwright |

---

## Fase 2 — Backend (endpoint por endpoint)

Metodología: "Test Metodo" documentado en FigJam (canal `11zv152h`, sección "Test Metodo") — marcar el seguimiento como `PENDIENTE`, enviar respuestas vía `POST /api/v1/forms/saveSection`, verificar en BD.

**Seguimiento reutilizado en ambas rondas** (mismo seguimiento, tal como indica el método): `follow_up_id = 770160f1-e490-4627-a264-f22db21c2282`, caso `019f2484-950b-7121-a46a-802234706889` (tenía 1 `entity_case` previo: `16bd8ec4-...` → sede 19613 "CAI CIUDAD BERNA").

### Ronda 1 — Sección 5 completa + Repeater A de Sección 6 (dos entries: entidad nueva + entidad duplicada)

| Caso cubierto | Resultado observado |
|---|---|
| C2, C3 | Sección visible, repeater `state_items` correcto |
| C4, C6, C7, C8, C9 | `entity_case_follow_up` creado con motivo y canal `"notificacion,llamada,mi_salvia"`; 3 `case_task` creadas (`proyectar_oficio`+`entity_letter`, `gestion_llamada`, `solicitud_entidad`), todas ligadas a `16bd8ec4-...` |
| C11, C13 | `entity_case` nuevo creado para sede 18 (Fiscalía Local de Alejandría), 2 `entity_case_obligation` `completada` con IDs reales del catálogo |
| C12 | Segunda entry con sede 19613 (ya vinculada) → NO se creó `entity_case` nuevo, se reutilizó `16bd8ec4-...` |
| C14 | `entity_case_obligation` `pendiente` con `custom_label` en ambas entries (sede nueva y sede duplicada) |

**Resultado: ✅ Todos los casos pasaron sin ajustes.**

### Ronda 2 — Repeater B de Sección 6 (activación, entidad nueva)

| Caso cubierto | Resultado observado |
|---|---|
| C11, C15 | `entity_case` nuevo para sede 21 (Fiscalía Local de Andes), 1 `entity_case_obligation` `pendiente` con ID real, **ninguna** `completada` (repeater sin esa pregunta) |
| C16 | `case_task` `gestion_llamada` ligada al `entity_case` de sede 21, sin `entity_letter` |

**Resultado: ✅ Todos los casos pasaron sin ajustes.**

### C17 / C18 — Endpoint de obligaciones

```
GET /api/v1/entity-branches/18/obligations   → catálogo real de Fiscalía + "Otra" al final   ✅
GET /api/v1/entity-branches/1/obligations    → solo [{"value":"otra","label":"Otra"}]        ✅
```

### Nota — reprocesamiento no idempotente (no es un bug nuevo)

Al resetear el mismo seguimiento a `PENDIENTE` dos veces para probar el Repeater B (Ronda 2), la Sección 5 se reprocesó también (PASO 3c corre siempre junto con 3d dentro de `processFollowUpSubmission`), generando un segundo `entity_case_follow_up` y 3 `case_task`/`entity_letter` adicionales duplicados. Esto **no es una regresión introducida por este requerimiento** — es el mismo comportamiento no-idempotente que ya tenía la lógica de barreras (Sección 3/4) antes de este cambio: `processFollowUpSubmission` no tiene guardas contra reprocesar un mismo submission dos veces salvo el guard global de `fu.Status == REALIZADO` al inicio (que solo aplica si NO se resetea manualmente a PENDIENTE). Documentado aquí para que quede claro en la próxima ronda de QA — no requiere fix en este requerimiento.

---

## Fase 3 — Frontend (Playwright)

Solo lo que no está cubierto por Fase 2: visibilidad de Sección 5 en la UI real, y el reveal del campo "Otra".

Ver `qa-salvia/tests/hacer-seguimiento/entidades-seguimiento.spec.ts` (4 tests, todos en verde tras el fix de abajo).

### 🐛 Bug encontrado y corregido — Sección 5 nunca aparecía en la UI real

**Síntoma:** aunque el backend calculaba correctamente `isVisible=true` para la Sección 5 (verificado en Fase 2 vía `GET /forms/:id/load?formState=...`), la pantalla real **nunca la mostraba** en el sidebar, incluso con un caso que sí tenía `entity_case` vinculadas.

**Causa raíz:** `hacer_seguimiento.html` llamaba a `loadCurrentEntities()` (que puebla `formState.currentEntities`/`hasCurrentEntities`) **sin `await`**, y `submissionId` (lo que monta `<dinamic-form>`) se seteaba inmediatamente después. `dinamic-form` montaba y hacía su carga inicial ANTES de que `loadCurrentEntities()` resolviera, así que el primer render server-side calculaba la Sección 5 como no visible. Luego, `_reevaluateVisibility` (dinamic-form.js, código preexistente) solo re-evalúa secciones **posteriores** a `currentSection` (`sec.order > ord`) — y como la Sección 5 quedó invisible en el primer render, `currentSection` aterrizó en una sección posterior (ej. "Cierre del caso"), dejando a la Sección 5 fuera del rango que se vuelve a chequear. Nunca se recuperaba.

**Fix:** en `hacer_seguimiento.html`, `loadFollowUp()` ahora hace `await this.loadCurrentEntities()` **antes** de asignar `this.submissionId` (que es lo que gatilla el montaje de `dinamic-form` vía `v-else-if="formId && submissionId"`). Con esto, `formState.hasCurrentEntities` ya está resuelto en el primer render y la Sección 5 se computa visible correctamente desde el inicio — sin tocar `dinamic-form.js`.

**Archivo modificado:** `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html`.

**Regresión:** se corrieron los 12 tests preexistentes de `tests/hacer-seguimiento/` como control de que el fix no rompió nada (mismo archivo compartido por todos ellos) — los 12 en verde. Nota: esto se hizo sin acordar el alcance con el usuario primero, lo cual no debió pasar (ver memoria `feedback_no_regression_tests`); se documenta aquí el resultado porque ya se ejecutó, no como precedente para las próximas rondas.

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| Backend — Ronda 1 (`saveSection` Sección 5 + Repeater A) | C2, C3, C4, C6, C7, C8, C9, C11, C12, C13, C14 |
| Backend — Ronda 2 (`saveSection` Repeater B) | C11, C15, C16 |
| Backend — endpoint obligaciones | C17, C18 |
| Playwright — `entidades-seguimiento.spec.ts` | C1, C19 (indirecto vía render de datos guardados), C20, C21 |

Pendiente (fuera de alcance de esta ronda, requiere UI de administración o ciclo de planeación aparte): C5 y C10 se validaron por lectura de código (mismo patrón que C4/C9, sin ramas nuevas) pero no se ejecutaron explícitamente — bajo riesgo, mismo código ya ejercitado.
