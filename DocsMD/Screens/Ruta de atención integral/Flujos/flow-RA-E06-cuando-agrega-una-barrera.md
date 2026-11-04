━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega una barrera a una entidad acudida
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: abrirBarreraParaEntidad(e) · agregarBarrera() · armarPayloadBarrera() → registro-v3-app.js · AddBarrier() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Es el mismo mecanismo de [flow-E18](../../Registro%20V3/Flujos/flow-E18-cuando-agrega-o-elimina-una-barrera.md): formulario de barrera de Seguimiento v4 (sin competencias) dentro de la tarjeta de la entidad.

PASO 1 — «+ Agregar barrera a esta entidad» (GS_17 = Sí y Caso en Borrador): abre el formulario con `EntityBranchID` de la tarjeta
PASO 2 — Sector* → preguntas del sector · fecha* · funcionario* · departamento* → municipio* · vereda o zona* · descripción* · gestión* (activación de ruta / alerta) · [con alerta] perspectiva profesional* · descripción de la gestión*; el botón se habilita sin faltantes
PASO 3 — [POST] /api/v1/victim-cases-v3/{icode}/barreras → borrador JSON (409 `readOnly` si el Caso ya salió de Borrador)
PASO 4 — Al iniciarse el Caso: `barrier_v2` + oficio y tarea por «Activación de ruta» y/o tarea de comité por «Alerta por barreras»

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| REG-29 (dónde se gestionan las barreras) sigue pendiente del cliente: el botón está en las tarjetas de ambas listas (D8) | PASO 1 |
