## 2026-10-28 — Pregunta de riesgo inminente en el seguimiento

El seguimiento viejo y el V.4 preguntan si la remisión a medidas de emergencia es de riesgo inminente. La pregunta solo se ve cuando eligió Medidas de Emergencia. El seed es aditivo: no reescribe formularios que ya tienen envíos.

Si la respuesta es Sí, la cabecera nace con urgencia `riesgo_inminente` y estado `riesgo_inminente_por_decidir_via`. Cada medida elegida en esa remisión se guarda con `riesgo_inminente = true`. Se crea una `case_task` `decidir_via_riesgo_inminente` asignada al rol `sm`. Si es No, la remisión sigue en urgencia normal y pendiente de revisión, las medidas quedan en false y no hay tarea. Un fallo al crear la tarea no revierte la remisión.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/seed/seguimiento_riesgo_inminente.sql` | Pregunta booleana del formulario viejo |
| `src/internal/seed/seguimiento_v4_riesgo_inminente.sql` | Pregunta booleana del V.4 |
| `src/internal/models/emergency_measure.go` | Columna `riesgo_inminente` |
| `src/salvia/service/form_service.go` | Cabecera, medidas y tarea al procesar el seguimiento viejo |
| `src/salvia/service/form_service_v4.go` | Misma regla en el plan del V.4 |
| `DocsMD/Screens/hacer-seguimiento/form-seguimiento-data.md` | Pregunta y efecto al responder Sí |
| `DocsMD/Screens/hacer-seguimiento/formulario-seguimiento-v4.md` | Pregunta del plan y efecto al completar |
| `DocsMD/Screens/hacer-seguimiento/Flujos/flow-E04-cuando-se-procesa-submission.md` | PASO 5 de medidas |
