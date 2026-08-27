# QA — Límite de caracteres en campos de texto

Bug: "Sin límite de caracteres en formularios de tareas" — varios `<input>`/`<textarea>` de `case-task-modal.js` no tenían `maxlength`. Ver tabla completa en [case-task-modal-interface.md § Límite de caracteres por campo](../case-task-modal-interface.md).

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Campos de `gestion_llamada` truncan al escribir más del límite | Modal abierto con `tarea.type = gestion_llamada`, "Otra entidad" seleccionada, `generaOficio = true` | `Nombre de la entidad` corta a 150, `Funcionario contactado` a 50, `Descripción / Notas` a 500, `Asunto del oficio` a 80, `Ruta Kofax` a 255 | Playwright — `fill()` con límite+50 chars, `toHaveValue()` con el string truncado |
| C2 | Textareas de `comite_caso` truncan al escribir más del límite | Modal abierto con `tarea.type = comite_caso`, las 3 decisiones que revelan textarea marcadas | `Observaciones del oficio`, `...para el agente` y `...del mecanismo` cortan a 500 c/u | Playwright — igual técnica que C1 |

> `proyectar_oficio` no se probó aparte — comparte los mismos límites (`funcionario`=50, `asunto`=80, `rutaKofax`=255, `entidadNombre`=150) ya cubiertos por C1 en `gestion_llamada`; hubiera sido redundante.

---

## Fase 3 — Frontend (Playwright, headless)

Test: `qa-salvia/tests/case-task-modal/limite-caracteres.spec.ts`.

```
CI=1 npx playwright test tests/case-task-modal/limite-caracteres.spec.ts

✓ gestion_llamada — campos de texto respetan su límite de caracteres (28.8s)
✓ comite_caso — observaciones respetan el límite de 500 caracteres (27.5s)
2 passed
```

### Hallazgo durante el setup — el panel dev para abrir tareas ya no existe

El panel `input[placeholder="UUID de CaseTask"]` + botón "Abrir modal" en `get_case_detail_sv.html` quedó **comentado** ("pruebas finalizadas") — reemplazado por el tab real **"📝 Tareas"** (componente `case-tasks.js`, ya mergeado). Se actualizó `detalle-caso.page.ts` → `abrirTaskModal()` para navegar a ese tab, ubicar la tarea vía `GET /api/v1/case-tasks/:id` (para obtener su `description`, no expuesta como id en el DOM) y hacer click en su botón "Gestionar". Ese botón solo es visible si el usuario logueado puede completar la tarea (`puedeCompletar`: rol `sv`, o `assignedUserId === userId`) — se agregaron helpers `assignTask()` y `setTaskDescription()` en `helpers/db.ts` para asignar las tareas fixture a `TEST_USER` y darles una descripción única (evita colisión con ~24 tareas históricas del mismo caso E2E compartido que reutilizan la misma descripción de producción).

Este fix de navegación queda disponible para cualquier otro test de `case-task-modal`/`case-task-history` que dependiera del panel viejo.

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `limite-caracteres.spec.ts` › `gestion_llamada — campos de texto respetan su límite de caracteres` | C1 |
| `limite-caracteres.spec.ts` › `comite_caso — observaciones respetan el límite de 500 caracteres` | C2 |

**Resultado:** 2/2 casos cubiertos y verdes.
