## 2026-11-02 — «Registro de caso»: borrador de Registro V3 desde un reporte (WP-1 T7, REP-03)

Desde el modal de gestión de un reporte (Víctima o Tercero con contacto exitoso) se crea un borrador de Registro V3 prellenado y se abre en el formulario. Corte vertical nuevo: no se modifica `CreateDraft` ni los archivos compartidos de V3 (solo +1 línea y +1 método en `registro-v3-app.js`).

| Archivo | Cambio |
|---|---|
| `src/salvia/controller/victim_case_v3_from_report_controller.go` | Nuevo — `POST /api/v1/victim-cases-v3/desde-reporte/:icode` |
| `src/salvia/service/victim_case_v3_from_report_service.go` | Nuevo — validación previa, credenciales y orquestación |
| `src/internal/repository/victim_case_v3_from_report_repository.go` | Nuevo — la transacción de conversión (reporte bloqueado `FOR UPDATE`, `CreateDraft` en un savepoint, vínculo, prellenado y hito) |
| `src/frontend/js/components/registro-v3-app.js` | +1 llamada en `mounted()` y el método `abrirDesdeReporte()` |
| `src/main.go` | cableado del controller nuevo |

Los borradores creados desde un reporte quedan **sin documento** (tipo y número vacíos) hasta que el agente complete CI_5 y CI_6: el documento temporal `cc/0000000000` de `CreateDraft` hacía que `fanOutPerfilV3` replicara el Perfil entre borradores sin documento (ver el GAP de `Flujos/flow-E19-cuando-carga-con-registro-de-origen.md`).

El reporte NO se cierra: queda «Borrador iniciado» (estado derivado) y se lee «Cerrado — válido» cuando el caso queda en `te` o `cd`; ningún código escribe `ap`.
