# QA — Validar remisión psicosocial + tarea "Justificar Remisión"

Pantallas: `/salvia/remision-psicosocial/:id` (`remision_psicosocial_detalle.html`, modal "Validar remisión") y Detalle del Caso → tab Tareas (`case-task-modal.js`, tipo `justificar_remision`). Cubre 3 cambios relacionados, en orden causal:

1. **Fix de integración:** `processFollowUpSubmission` (`form_service.go`) creaba la `case_task` de validación con `type: "validar_remision_psicosocial"`, pero el modal "Validar remisión" (`ValidateRemision`, `psychosocial_detail_service.go`) busca `type = 'validar_remision'` para completarla. Se alineó el `type`.
2. Cuando la remisión se devuelve (`isValid=false`), `ValidateRemision` crea una `case_task` `justificar_remision` asignada a `victim_case.agent_id` (el agente dueño del caso — no el agente del seguimiento puntual, campo distinto). Si se valida como correcta, no se crea nada extra.
3. **Nuevo:** `case-task-modal.js` ahora reconoce `justificar_remision` — el agente ve el motivo de la devolución (solo lectura) y responde en un textarea requerido. Al completarla (`sideEffectsJustificarRemision`, `case_task_service.go`): la remisión vuelve a `abierto`, y se crea una nueva `case_task` `validar_remision` — asignada al `psychologist_id` de la `dupla` de la remisión, o al `professional_id` individual si no hay dupla, o **sin asignar** si la remisión no tiene ninguno de los dos.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Validar como correcta | Remisión con `case_task validar_remision` en `ToDo`, banner "Revisar" visible | Modal "✓ Es válida" → confirmar: banner desaparece, `validationPending=false`, `validar_remision → Done`, **no** se crea `justificar_remision` | Playwright — DOM (banner) + API (`/detail`, `/case-tasks`) |
| C2 | Devolver remisión con motivo | Mismo fixture, reseteado | Modal "✕ No cumple criterios" + motivo → confirmar (botón disabled hasta escribir motivo): banner desaparece, `status = en_devolucion`, `validar_remision → Done`, **se crea** `justificar_remision` (`ToDo`, `assignedUserId = victim_case.agent_id`, `description` incluye el motivo) | Playwright — DOM + API |
| C3 | Completar Justificar Remisión — remisión **con dupla** | `case_task justificar_remision` en `ToDo`, remisión con `dupla_id` set | Modal "Justificar Remisión" muestra el motivo (`tarea.description`) + textarea "Tu respuesta" (botón disabled hasta escribir) → confirmar: `justificar_remision → Done` con `formData.respuesta`, `status → abierto`, se crea **nueva** `validar_remision` (`ToDo`) asignada al `psychologist_id` de la dupla | Playwright — DOM (`case-task-modal.js`) + API |
| C4 | Completar Justificar Remisión — remisión **sin dupla ni profesional** | `case_task justificar_remision` en `ToDo`, remisión sin `dupla_id` ni `professional_id` | Igual que C3, pero la nueva `validar_remision` queda **sin asignar** (`assignedUserId = ''`) | Backend directo (`PUT /case-tasks/:id/complete` + query BD) — no repetido por Playwright, mismo mecanismo que C3, solo cambia el resultado de una consulta SQL (`dupla_id`/`professional_id` nulos) |

> No se probó el caso "reason vacío" en C2 (validación 400 del controller) por separado — es de Juan Ruiz, no del cambio bajo prueba, y el botón "Devolver remisión" ya queda `disabled` en la UI hasta que se escribe el motivo (cubierto implícitamente en C2). Tampoco se probó la rama `professional_id` individual (sin dupla) por separado — es el mismo `assignedTo := professionalID` sin el paso extra de resolver la dupla, ya cubierto en principio por C3/C4.

---

## Fixtures reutilizables

`qa-salvia/helpers/db.ts`:

| Fixture (`psychosocial_support.id`) | dupla_id | Reset helper |
|---|---|---|
| `0199e2b1-…0001` | — (sin asignar) | `resetValidarRemisionFixture(id)` — status → `abierto`, recrea `validar_remision` ToDo |
| `0199e2b1-…0002` | `c1ea6d8b-…` ("QA Dupla Claude", psicóloga `019fba3f-…5e94-…`) | `resetJustificarRemisionFixture(id, assignedUserId, description)` — status → `en_devolucion`, recrea `justificar_remision` ToDo asignada a `assignedUserId` con `description` única |

Ambas sobre el caso E2E compartido `019eadd1-…` / `follow_up_id = fd053f17-…`.

---

## Fase 2 — Backend directo (C4)

```
PUT /api/v1/case-tasks/:id/complete   { userId, formData: { respuesta } }
→ 200, justificar_remision → Done

GET /api/v1/case-tasks?caseId=...&psychosocialId=0199e2b1-…0001
→ validar_remision | ToDo | assignedUserId=''
→ justificar_remision | Done | assignedUserId=''

GET /api/v1/psychosocial-support/.../detail → status: abierto
```

## Fase 3 — Frontend (Playwright, headless)

Test: `qa-salvia/tests/psicosocial-sesion/validar-remision-justificar.spec.ts`.

```
CI=1 npx playwright test tests/psicosocial-sesion/validar-remision-justificar.spec.ts

✓ validar como correcta — completa la tarea, NO crea Justificar Remisión (9.6s)
✓ devolver remisión con motivo — completa la tarea y SÍ crea Justificar Remisión asignada al agente del caso (14.5s)
✓ completar Justificar Remisión — reabre la remisión y reasigna Validar Remisión al psicólogo de la dupla (33.8s)
3 passed
```

### Hallazgos durante las pruebas

- **`CASE_AGENT_ID` mal asumido en el primer intento de C2:** se asumió `follow_up_v2.agent_id` (agente del seguimiento puntual, `019a5234-…`), pero el código correctamente usa `victim_case.agent_id` (agente actual del caso, `019edb41-…` en este fixture — puede diferir tras una reasignación). Error de la aserción del test, no del código.
- **C3 requirió timeouts más generosos:** el caso E2E compartido acumuló >70 `case_task` históricas durante esta sesión de QA; el tab Tareas (`case-tasks.js`) las renderiza todas de una vez, más lento que la pantalla liviana de detalle de remisión que usan C1/C2. Se subió `test.setTimeout` a 90s y los timeouts de click/overlay en `detalle-caso.page.ts` — confirmado con BD directa que la primera corrida "fallida" en realidad sí había completado todo correctamente server-side; era únicamente el `expect` de Playwright quedándose sin tiempo.

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `validar-remision-justificar.spec.ts` › `validar como correcta` | C1 |
| `validar-remision-justificar.spec.ts` › `devolver remisión con motivo` | C2 |
| `validar-remision-justificar.spec.ts` › `completar Justificar Remisión ... dupla` | C3 |
| Backend directo (`PUT /complete` + query BD, ver Fase 2) | C4 |

**Resultado:** 4/4 casos cubiertos y verdes.

---

## Nota de alcance

No se documentó el árbol de interfaz completo (`interfaz.md`) ni el inventario de eventos (`index.md`) de la pantalla `remision-psicosocial-detalle` — no existían antes de este cambio y su documentación completa es un trabajo aparte, no solicitado en este QA. El tipo `justificar_remision` de `case-task-modal.js` sí está documentado en `DocsMD/Componentes/case-task-modal/`.
