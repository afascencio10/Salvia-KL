## 2026-08-20 — Reabrir caso desde el detalle

Cuando un caso está cerrado (`victim_case_status = cd`), el header muestra el botón "Reabrir caso" a la derecha de "+ Nuevo seguimiento". Confirmar en el modal llama `POST /api/v1/casos/:id/reabrir`, pasa el caso a `ra` (Activo) y registra el evento de timeline "Reapertura de Caso". Mismos roles que Nuevo seguimiento; `ps`/`ts` no lo ven.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/case_closure_service.go` | Nuevo `ReabrirCaso` (cd → ra + timeline) |
| `src/internal/models/case_timeline_event.go` | Constante `TimelineTypeReabrirCaso` |
| `src/salvia/controller/case_detail_controller.go` | `POST /casos/:id/reabrir`; inyecta `CasoCierreService` |
| `src/main.go` | Pasa `casoCierreSvc` al controller de detalle |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Botón, modal de confirmación y `confirmarReabrirCaso` |
| `DocsMD/Screens/Detalle del Caso/detalle-caso-interface.md` | HeaderRight y modal documentados |
| `DocsMD/Screens/Detalle del Caso/detalle-caso-index.md` | Eventos E21 y E22 |
| `DocsMD/Screens/Detalle del Caso/Flujos/flow-E22-cuando-confirma-reapertura.md` | Flujo del POST y recarga |
| `DocsMD/Screens/Detalle del Caso/QA/E22-reabrir-caso-qa-cases.md` | Casos C1–C4 del botón y del POST |
