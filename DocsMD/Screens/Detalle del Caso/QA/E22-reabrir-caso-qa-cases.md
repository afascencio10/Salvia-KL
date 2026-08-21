# QA — Reabrir caso (E-21 / E-22)

Casos de prueba para el botón "Reabrir caso" en `/salvia/casos/:id/detalle` ([flow-E22](../Flujos/flow-E22-cuando-confirma-reapertura.md)).

| Caso | Descripción | Condición inicial | Resultado esperado | Método |
|---|---|---|---|---|
| C1 | Botón visible en caso cerrado | `victim_case_status = cd` y rol `sv` / `ro` / `op` asignado | El header muestra "Reabrir caso" a la derecha de "+ Nuevo seguimiento" | UI |
| C2 | Botón ausente si el caso no está cerrado | `victim_case_status = ra` (u otro distinto de `cd`) | "Reabrir caso" no está en el DOM | UI |
| C3 | Confirmar reapertura | Caso `cd`; confirmar en el modal | `POST /api/v1/casos/:id/reabrir` 200; BD: `victim_case_status = ra`; timeline type `Reapertura de Caso`; badge "Activo"; el botón desaparece | Request + BD + UI |
| C4 | Cancelar modal | Caso `cd`; abrir modal y Cancelar | Status sigue `cd`; no hay evento de timeline nuevo | UI + BD |
