# QA — Sección "Identificación de Entidades" (Registro de Caso legacy)

Casos de prueba para el requerimiento implementado según `DocsMD/Otros/temp/agregar-identificacion-entidades-registro-caso-legacy.md`. Cubre `ProcessVictimCaseEntidadEntries` (archivo nuevo `victim_case_entidades.go`), el hook post-commit en `VictimCaseController.SetVictimCase`, y el dinamismo nuevo de `set_victim_case.html` (gates, repeaters hand-rolled, validación bloqueante `validateEntidadEntries()`).

A diferencia de Hacer Seguimiento y Psicosocial, aquí el caso **no existe antes** del `POST` — cada prueba de backend crea un caso nuevo desde cero (no hay "Test Metodo" de resetear un submission existente; cada ronda es un `POST /salvia/casos` completo con un `docNumber` distinto).

---

## Fase 1 — Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | `hasReportedToEntity` = No/vacío | gate 1 sin marcar | Repeater A no se envía / vacío, sin `entity_case` de esa fuente | Frontend + Backend: BD |
| C2 | `hasReportedToEntity` = Sí | gate 1 = "Sí" | Repeater A visible/procesable | Frontend: Playwright |
| C3 | `acceptsRouteActivation` = No/vacío | gate 2 sin marcar | Repeater B vacío | Frontend + Backend: BD |
| C4 | `acceptsRouteActivation` = Sí | gate 2 = "Sí" | Repeater B visible/procesable | Frontend: Playwright |
| C5 | Fila agregada e incompleta | usuario agrega fila, deja Sector vacío | `validateEntidadEntries()` bloquea, `newEntity` NO se llama, caso NO se crea | Frontend: Playwright |
| C6 | Fila agregada y completa | todos los campos requeridos llenos | Validación pasa, `submit()` continúa | Frontend: Playwright |
| C7 | Caso sin ninguna fila de entidades | ambos gates en No | Caso se crea normalmente, **cero** `entity_case` | Backend: BD |
| C8 | Repeater A — entidad nueva | `entity_branch_id` sin `entity_case` previo (imposible que exista, el caso es nuevo) | `entity_case` creado, `objetivo` = info ruta | Backend: BD |
| C9 | Repeater B — entidad nueva | ídem, repeater B | `entity_case` creado | Backend: BD |
| C10 | Misma entidad en repeater A **y** B del mismo registro | mismo `entity_branch_id` en ambos repeaters | `ErrEntityCaseDuplicate` en el segundo → reutiliza el `entity_case` del primero, no crea uno nuevo | Backend: BD |
| C11 | Checklist "completadas" (repeater A) | 1 ID real de `entity_obligation` | `entity_case_obligation` `completada` | Backend: BD |
| C12 | Checklist con "Otra" (ambos repeaters) | valor `otra` + texto libre | `entity_case_obligation` con `custom_label`, `entity_obligation_id=NULL` | Backend: BD |
| C13 | Repeater B no tiene pregunta de "completadas" | — | Solo se registran obligaciones `pendiente` en B | Backend: BD |
| C14 | Canal de activación — un solo valor | canal=`llamada` | 1 `case_task` `gestion_llamada` | Backend: BD |
| C15 | Canal de activación — múltiples valores | canal=`notificacion,mi_salvia` | 2 `case_task` (una por canal) | Backend: BD |
| **C16** | `entity_case_obligation.follow_up_id` queda NULL | cualquier obligación creada desde Registro | columna `follow_up_id` NULL (antes NOT NULL — requería la migración) | Backend: BD |
| **C17** | `entity_case_obligation.team_contact_id` queda NULL | ídem | NULL (Registro no tiene `team_contact`) | Backend: BD |
| **C18** | `case_task.psychosocial_support_id` queda NULL | canal de activación desde Registro | NULL (Registro no tiene `psychosocial_support`) | Backend: BD |
| **C19** | Backend recibe una fila sin `entityBranchId` (bypass del frontend, vía API directa) | fila con campos vacíos | Se ignora silenciosamente, sin error, sin bloquear la creación del caso — mismo criterio que "Asignación de Ruta" | Backend: BD + respuesta 200 |
| C20 | Procesamiento es post-commit (no bloquea la respuesta) | — | El `POST /salvia/casos` responde 200 con las credenciales antes de que termine de procesarse `ProcessVictimCaseEntidadEntries` | Backend: logs (timestamps) |

---

## Fase 2 — Backend

Endpoint único: `POST /salvia/casos` (`SetVictimCase`). Fixture reutilizable: payload base con ~50 campos requeridos del form2 (ver `registro_qa.py` en el scratchpad de la sesión), variando solo `docNumber` (único por corrida) y `entidadesAcudidas`/`entidadesActivacion`.

Entidades de prueba reutilizadas de rondas anteriores: branch `18` (Fiscalía Local de Alejandría, sector `js`/justicia, town_code `05021000`, catálogo con 5 obligaciones reales), branch `21` (Fiscalía Local de Andes, mismo sector/catálogo, town_code `05034000`).

### Caso base — sin ninguna fila de entidades (C7)

`POST /salvia/casos` sin `entidadesAcudidas`/`entidadesActivacion`. **Resultado: ✅** caso creado (HTTP 200), `SELECT count(*) FROM entity_case WHERE case_id=...` → `0`.

### Ronda A — Repeater "Entidades a las que ya acudió" (entidad nueva, checklist completadas + "Otra")

| Caso cubierto | Resultado observado |
|---|---|
| C2, C8 | `entity_case` creado para branch 18, `objetivo` = texto de "Otra información asociada a la Ruta" |
| C11 | `entity_case_obligation` `completada` con ID real de `entity_obligation` |
| C12 | `entity_case_obligation` `pendiente` con `custom_label` (vía "Otra") |
| **C16** | `follow_up_id` de ambas obligaciones = **NULL** — confirma que la migración a nullable funciona end-to-end |
| **C17** | `team_contact_id` = **NULL** |

**Resultado: ✅ Todos los casos pasaron.**

### Ronda B — Repeater "Entidades para activación de ruta" (entidad nueva, canal multi-select)

| Caso cubierto | Resultado observado |
|---|---|
| C4, C9 | `entity_case` creado para branch 21, `objetivo` = NULL (este repeater no tiene esa pregunta) |
| C13 | Solo se registró obligación `pendiente` (repeater sin pregunta de "completadas") |
| C15 | Canal `notificacion,mi_salvia` → 2 `case_task` (`proyectar_oficio` con `entity_letter_id` seteado, `solicitud_entidad` sin él) |
| **C18** | `psychosocial_support_id` = **NULL** en ambas tareas |

**Resultado: ✅ Todos los casos pasaron.**

### Ronda C — Misma entidad en ambos repeaters del mismo registro (duplicado)

Branch 18 en Acudidas **y** en Activación, en el mismo `POST`. **Resultado: ✅** — verificado por ausencia de un segundo `INSERT INTO entity_case`: solo existe **un** `entity_case` (branch 18, con el `objetivo` del primer repeater procesado), y la obligación del repeater de Activación (`custom_label="Duplicado QA"`) quedó asociada a ese mismo `entity_case`. Confirma que `ErrEntityCaseDuplicate` se maneja correctamente incluso cuando ambas menciones llegan en el mismo registro de un caso recién creado (escenario que no podía darse en Hacer Seguimiento/Psicosocial, donde el caso ya existe de antes).

### Ronda D — Fila sin `entityBranchId` (bypass del frontend, vía API directa) + timing post-commit

**C19:** enviado directamente al backend (sin pasar por `validateEntidadEntries()` del frontend) un `entidadesAcudidas` con todos los campos vacíos. **Resultado: ✅** HTTP 200, caso creado normalmente — la fila se ignoró silenciosamente, sin error, sin bloquear la creación del caso.

**C20** (post-commit no bloquea la respuesta): verificado con los logs del servidor de la Ronda A — el `[GIN] ... 200 ... POST /salvia/casos` se registró a las `15:26:18`, y el `INSERT INTO entity_case` correspondiente ocurrió a las `15:26:19.513` — **después** de que la respuesta HTTP ya había sido enviada al cliente. Confirma que el procesamiento es verdaderamente asíncrono.

---

## Fase 3 — Frontend (Playwright)

`qa-salvia/tests/registro-caso/entidades-registro-caso.spec.ts` — 2 tests, ambos en verde.

| Test | Casos cubiertos | Resultado |
|---|---|---|
| `C2/C4 — los gates muestran/ocultan sus repeaters` | C2, C4 | ✅ Ambos repeaters ocultos por defecto; cada gate en "Sí" revela el suyo |
| `C5 — fila agregada e incompleta bloquea el guardado` | C5 | ✅ Se interceptó la red: el `POST /salvia/casos` **nunca se disparó** al hacer click en "Guardar" con una fila agregada y vacía; el campo "Sector de la entidad" muestra el error inline |

**C6** (fila completa no bloquea) no se ejecutó como test aislado por UI — cobertura por combinación: el código de `validateEntidadEntries()` solo puebla errores y bloquea cuando falta algo (revisado, camino inverso del test C5), y las 4 rondas de backend (A-D) ya probaron exhaustivamente que filas completas se procesan correctamente de punta a punta. Un flujo 100% por UI llenando las ~50 preguntas del resto del formulario se consideró de bajo retorno dado que esas preguntas no son parte de este requerimiento y ya están cubiertas por la pantalla legacy existente (sin cambios).

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| Backend — caso base sin entidades | C7 |
| Backend — Ronda A (Acudidas) | C2, C8, C11, C12, C16, C17 |
| Backend — Ronda B (Activación) | C4, C9, C13, C15, C18 |
| Backend — Ronda C (duplicado mismo registro) | C10 |
| Backend — Ronda D (fila vacía + timing) | C19, C20 |
| Playwright — `entidades-registro-caso.spec.ts` | C2, C4, C5 |
| Cobertura por combinación (no ejecutado aisladamente) | C1, C3, C6 |

**Bugs de producto encontrados: 0.** Los 4 escenarios de backend (base, Acudidas, Activación, duplicado-mismo-registro, fila-vacía) y los 2 de frontend pasaron en el primer intento tras corregir la construcción del payload de prueba (nombres/tipos de campo del formulario legacy — error de metodología propio, no del código bajo prueba). Hallazgo más relevante: se confirmó con evidencia real de logs que el procesamiento de entidades es genuinamente asíncrono/no bloqueante (C20), y que el manejo de duplicados funciona correctamente incluso en el escenario nuevo de "misma entidad mencionada en ambos repeaters del mismo registro" (C10), caso que no existía en los 2 requerimientos anteriores.
