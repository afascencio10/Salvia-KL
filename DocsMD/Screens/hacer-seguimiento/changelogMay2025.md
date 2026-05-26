# Changelog Mayo 2025 — `hacer-seguimiento`

---

## 2025-05-25 — Cierre automático de caso al completar el formulario de seguimiento

Cuando el profesional marca cierre de caso en la Sección 5 del formulario de seguimiento y guarda esa sección, el backend detecta la respuesta en `processFollowUpSubmission` y llama al `CasoCierreService.CerrarCaso`, que actualiza `victim_case_status` a `"cd"` y registra un evento de timeline. La lógica de cierre es modular y puede invocarse desde cualquier parte del sistema.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/case_closure_service.go` | Nuevo servicio modular `CasoCierreService` con `CerrarCaso` |
| `src/salvia/service/form_service.go` | Paso 8 en `processFollowUpSubmission`: detecta cierre y llama `CasoCierreService` |
| `src/internal/repository/victim_case_light_repository.go` | Agrega método `UpdateStatus` a la interfaz e implementación |
| `src/main.go` | Instancia `casoCierreSvc` e inyecta en `FormServiceDeps` |

---

## 2025-05-25 — Bloqueo de edición del formulario cuando el caso está cerrado

Cuando un caso tiene `victim_case_status = "cd"`, el formulario de seguimiento se abre en modo solo lectura permanente, anulando la ventana de 5 días de edición. El backend retorna `caseStatus` en la respuesta de `loadFollowUp` y el frontend evalúa esa condición con prioridad máxima al calcular `canEdit`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/followup_v2_service.go` | `LoadFollowUpResult` incluye `CaseStatus`; `LoadFollowUp` consulta el status del caso |
| `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` | `loadFollowUp()` evalúa `caseStatus === 'cd'` como primera condición de `canEdit` |
