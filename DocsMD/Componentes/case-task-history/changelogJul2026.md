## 2026-07-02 — Implementación de case-task-history + enriquecimiento de GetByID

Se implementó el componente `case-task-history` (modal de solo lectura para ver el detalle de una `case_task` completada), se montó en la pantalla `Detalle del Caso` junto a `case-task-modal`, y se actualizó `GetByID` para que devuelva `assignedUserName` — resolviendo el GAP documentado en el diseño inicial. También se escribieron 4 tests E2E de solo lectura (uno por tipo de tarea) que verifican el detalle sin volver a completar las tareas.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-history.js` | Nuevo — componente completo (estilos, template, `open`/`cerrar`) |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Monta `<case-task-history ref="taskHistory">` + panel dev de prueba |
| `src/salvia/controller/case_task_controller.go` | `GetByID` enriquecido con `assignedUserName`; refactor a `taskToJSON`/`lookupUserName` compartidos con el listado |
| `src/config/db_config.json` | Usuarios de BD actualizados a `salvia_legacy2`/`salvia_gorm2` (pool separado, mismas contraseñas) |
| `qa-salvia/.env.local` | `DB_USER` actualizado a `salvia_gorm2` |
| `qa-salvia/pages/detalle-caso.page.ts` | Agregados `abrirTaskHistory`, `historyTitulo`, `historyBody`, `cerrarTaskHistory` |
| `qa-salvia/tests/case-task-modal/ver-detalle-tarea-completada.spec.ts` | Nuevo — 4 tests de solo lectura (uno por tipo) |
