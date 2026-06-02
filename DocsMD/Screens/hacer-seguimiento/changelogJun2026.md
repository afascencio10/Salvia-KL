# Changelog Junio 2026 — `hacer-seguimiento`

---

## 2026-06-01 — Correcciones y mejoras al flujo de reasignación

**Bug fix: `riskLevel` siempre llegaba como 0 al frontend.** La query `LoadVictimInfoByCaseID` en `followup_repository.go` no incluía `victim_case_form2_risk_level` en el SELECT. Se agregó `COALESCE(f2.victim_case_form2_risk_level, 0) AS "RiskLevel"`.

**Reasignación actualiza `victim_case`.** Después de `ReasignarCalendario`, se actualiza `victim_case_team` y `agent_id` en la tabla `victim_case` con el equipo destino y el agente recién asignado (leído del primer PENDIENTE).

**Filtro "ninguno" en factores de riesgo.** Las preguntas de factores protectores, de riesgo y de riesgo extremo tienen una opción `ninguno` que ahora se excluye antes de correr la lógica de reasignación, tanto en frontend (`_updateReasignacionState`) como en backend (`reasignarCaso`).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/repository/followup_repository.go` | Agrega `risk_level` al SELECT de `LoadVictimInfoByCaseID` |
| `src/internal/repository/victim_case_light_repository.go` | Nuevo método `UpdateTeamAndAgent` |
| `src/salvia/service/form_service.go` | `reasignarCaso`: filtro "ninguno", llamada a `UpdateTeamAndAgent` post-reasignación |
| `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` | `_updateReasignacionState`: filtro "ninguno" antes de contar factores |

---

## 2026-06-01 — Registrar Barreras y Seguimiento a Barreras

Implementación completa del flujo de identificación y seguimiento de barreras institucionales desde el formulario de seguimiento.

### Registrar Barreras (Sección 4)

Los dropdowns de Departamento (Q12), Ciudad (Q13) y Municipio (Q14) del repeater de Sección 4 ahora cargan sus opciones dinámicamente desde `formState` en lugar de la tabla `option`. Al cargar la pantalla se fetchean departamentos y ciudades en background; al guardar cada sección se actualizan las opciones por barrera en cascada (ciudad depende de departamento seleccionado, municipio depende de ciudad). El backend ahora guarda un registro `barrier_v2` completo por cada entry del repeater (22 campos), en lugar de un registro por opción CSV seleccionada.

### Seguimiento a Barreras (Sección 3)

La Sección 3 (Seguimiento a Barreras) ahora está controlada por `stateItems = 'currentBarriers'`: se renderiza una entry por barrera activa del caso. Al cargar el seguimiento, el backend consulta las barreras activas (status != MANAGED), guarda los IDs en `fu.active_barrier_ids` (fijados para este seguimiento) y los retorna como `activeBarriers: [{ id, barrierName }]`. El frontend popula `formState.currentBarriers` con este array. Al completar el formulario, el backend procesa cada entry del repeater de Sección 3, crea un evento de timeline por barrera y actualiza el status a MANAGED si se marcó cierre.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` | `loadLocations()`, `formState.currentBarriers/statesColombia/newBarriers`, `_updateBarreraLocationOptions()` |
| `src/internal/models/follow_up_v2.go` | Campo `ActiveBarrierIDs *string` |
| `src/internal/models/barrier_v2.go` | 14 campos nuevos (sector, ubicación, estructurales, detalles) |
| `src/internal/models/location_light.go` | Nuevo — `DepartmentLight`, `CityLight`, `LocationOption`, `CityLocationOption` |
| `src/internal/repository/location_repository.go` | Nuevo — `GetDepartments`, `GetCities`, `GetTownsByCityID` |
| `src/internal/repository/barrier_v2_repository.go` | `FindActiveByCaseID`, `FindByIDs`, `UpdateStatus` |
| `src/internal/repository/followup_repository.go` | `UpdateActiveBarrierIDs` |
| `src/salvia/controller/location_controller.go` | Nuevo — 3 endpoints GET /api/v1/locations/* |
| `src/salvia/service/followup_v2_service.go` | `loadActiveBarriers`, `buildBarrierName`, `LoadFollowUpResult` con `ActiveBarriers` |
| `src/salvia/service/form_service.go` | PASO 3 refactorizado (1 BarrierV2 por entry), PASO 3b (Sección 3), `buildBarrierFollowUpSummary`, logs en saveSection |
| `src/frontend/js/components/dinamic-form.js` | `resolveQuestionOptions` expuesto en `methods` |
| `src/main.go` | `locationRepo`, `locationCtrl`, `BarrierV2` en AutoMigrate |
