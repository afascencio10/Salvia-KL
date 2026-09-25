## 2026-09-24 — Restringir "Reabrir caso" a rol supervisor

El endpoint `POST /api/v1/casos/:id/reabrir` no validaba rol en el backend — solo el botón se ocultaba en frontend según reglas de UI, por lo que cualquier rol autenticado podía llamarlo directo. Ahora exige sesión válida y rol `sv` (403 en otro caso) vía `requireCaseDetailAccess`, el mismo helper que ya protege `case-entities`. De paso, `actor_id` del evento de timeline se toma de la sesión (`s.UserICode`) en vez de confiar en el body enviado por el cliente. En frontend, el botón "Reabrir caso" ahora depende solo de `puedeReasignar` (antes también lo veían `ro` y el `op` dueño del caso).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/controller/case_detail_controller.go` | `ReabrirCaso` usa `requireCaseDetailAccess` + chequeo `CurrentRole == "sv"`; `actor_id` viene de sesión |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Botón "Reabrir caso": `v-if` reducido a `puedeReasignar && status === 'cd'` |
| `DocsMD/Screens/Detalle del Caso/detalle-caso-interface.md` | Condición del botón y tabla "Controles por rol" actualizadas |
| `DocsMD/Screens/Detalle del Caso/detalle-caso-index.md` | Descripciones de E21 y E22 actualizadas |
| `DocsMD/Screens/Detalle del Caso/Flujos/flow-E22-cuando-confirma-reapertura.md` | Nuevo PASO B0 (validación de sesión y rol) |
