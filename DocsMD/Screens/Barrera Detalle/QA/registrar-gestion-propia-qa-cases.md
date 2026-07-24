# QA — Cuando el Enlace registra una gestión propia

Casos de prueba derivados del flujo lógico de `req-registrar-gestion-propia-enlace.md` (`DocsMD/Otros/temp/`). Cubre el endpoint `POST /api/v1/case-tasks/gestion-propia` y su UI en `case-tasks.js` / `case-task-history.js` dentro de la pantalla Detalle de Barrera.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Enlace registra gestión propia en barrera de su departamento | Rol `en`, barrera `OPEN` en el mismo departamento asignado | `case_task` creada con `status=Done`, `type=gestion_propia`; evento en timeline; barrera pasa a `En Gestion` | Playwright + validación DOM + `GET /barriers-v2/:id/detail` |
| C2 | Botón oculto si la barrera es de otro departamento | Rol `en`, barrera de un departamento distinto al asignado | El botón "+ Registrar gestión propia" no se renderiza | Playwright — `toHaveCount(0)` |
| C3 | Botón oculto para roles distintos de `en` | Rol `sv` (u otro) sobre cualquier barrera | El botón no se renderiza | Playwright — `toHaveCount(0)` |
| C4 | Rechazo por rol en el backend (defensa en profundidad) | Sesión válida con rol ≠ `en`, POST directo a la API | `403` `"su rol no tiene permiso..."` | Playwright — `page.request.post` |
| C5 | Rechazo por departamento en el backend | Sesión `en`, POST directo con `barrierId` de otro departamento | `403` `"esta barrera no pertenece a tu departamento asignado"` | Playwright — `page.request.post` |
| C6 | Sin sesión | Contexto de navegador sin cookie de sesión | `401` `"no autenticado"` | Playwright — contexto limpio |
| C7 | Tipo de gestión inválido | Sesión `en`, `tipo` fuera de las 4 opciones válidas | `400` | Playwright — `page.request.post` |
| C8 | Validación de formulario — descripción vacía | Modal abierto, `descripcion` vacía o solo espacios | Botón "Registrar como completada" deshabilitado | Playwright — `toBeDisabled()` |
| C9 | Cancelar no crea la tarea | Modal abierto con datos, click en "Cancelar" | Modal se cierra, no se llama al endpoint | Playwright — `not.toBeVisible()` |
| C10 | Ver detalle de una gestión propia ya completada | Tarea `type=gestion_propia` ya `Done` | `case-task-history` muestra título "Gestión Propia", descripción y `formData.subtipo` | Playwright — `case-task-history.js` |

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `TC-GP-01` — flujo feliz completo (crear + UI + transición de barrera + timeline) | C1 |
| `TC-GP-02` — botón oculto, otro departamento | C2 |
| `TC-GP-03` — botón oculto, otro rol | C3 |
| `TC-GP-04` — POST directo, rol inválido | C4 |
| `TC-GP-05` — POST directo, departamento inválido | C5 |
| `TC-GP-06` — POST directo, sin sesión | C6 |
| `TC-GP-07` — POST directo, tipo inválido | C7 |
| `TC-GP-08` — validación de formulario (disabled + cancelar) | C8, C9 |
| `TC-GP-09` — ver detalle de tarea completada | C10 |

Archivo de tests: `qa-salvia/tests/case-tasks-gestion-propia/registrar-gestion-propia.spec.ts`
Page object: `qa-salvia/pages/barrera-detalle.page.ts`
Fixtures DB: `qa-salvia/helpers/db.ts` → `resetBarrierStatus`, `limpiarGestionesPropia`

## Backend

No se probaron endpoints adicionales por separado — `POST /api/v1/case-tasks/gestion-propia` es nuevo y se cubre íntegramente por los casos C1, C4-C7 vía Playwright `page.request` (equivalente a un test de integración de API, sin pasar por la UI para esos casos).

## Datos de prueba

| Variable (`.env.local`) | Valor | Rol |
|---|---|---|
| `GESTION_PROPIA_EN_LOGIN` / `_PASS` | `test.claude.en` | Enlace Territorial (`en`), `general_user_assigned_department='2'` |
| `GESTION_PROPIA_BARRIER_SAME_DEPT` | barrera real, `department_id='2'` | — |
| `GESTION_PROPIA_BARRIER_OTHER_DEPT` | barrera real, `department_id='29'` | — |
