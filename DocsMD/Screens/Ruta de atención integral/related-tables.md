# `Ruta de atención integral` — Tablas relacionadas

El módulo **no tiene tablas propias**. Lo que registra se guarda en el borrador JSON del Caso y, al iniciarse el Caso, en las tablas de Gestión institucional y Seguimiento (ver [`registro-v3-data.md`](../Registro%20V3/registro-v3-data.md#borrador-json-de-entidades-y-barreras-wp-4--reg-38) y [el sub-flujo](../Registro%20V3/Flujos/flow-subflujo-materializar-borrador.md)).

| Tabla | Uso |
|---|---|
| `salvia.victim_case.victim_case_v3_draft_data` (jsonb, nueva en WP-4) | Mientras el Caso es Borrador: `entities`, `barriers`; siempre: `entityMeta` (por sede: tipo, acuerdo, debida diligencia, canales, fecha de reporte, gestiones, ids de acciones) y `barrierMeta` (por barrera: zona, perspectiva, descripción de la gestión) y `materializedAt` |
| `salvia.victim_case_form2` (+ `rel_victim_case_form2_enums_victim_case_form2`) | `GS_14` (`reported_to_entity`) y `GS_17` (`identifies_barriers`), guardadas con la sección GS |
| `salvia.entity_case` | Una fila por entidad acudida al materializar (`objetivo` = «otra información asociada a la Ruta», `radicado`, `created_by_id`); compartida con Gestión institucional y Seguimiento |
| `salvia.entity_case_obligation` | Una fila por acción completada (`completada`) o pendiente (`pendiente`) |
| `salvia.barrier_v2` | Una fila por barrera al materializar (`OPEN`, `entity_branch_id`, `management_actions`) |
| `salvia.entity_letter` · `salvia.case_task` | Oficios y tareas de los canales de activación (acudida con acciones pendientes, acuerdo y debida diligencia en «Sí») y de la gestión de cada barrera |
| `salvia.entity_branch` · `salvia.entity` · `salvia.entity_obligation` · `security.department/city/town` | Catálogos de la cascada, de la sede y de sus acciones institucionales (solo lectura) |
