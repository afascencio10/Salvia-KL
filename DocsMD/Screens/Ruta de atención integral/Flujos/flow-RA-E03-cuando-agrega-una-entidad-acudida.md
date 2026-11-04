━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega una entidad acudida
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: agregarEntidad() · entidadFaltantes() · activacionAplica() · toggleCanal() · toggleAccion() → registro-v3-app.js · AddEntity() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Es el mismo mecanismo de [flow-E17](../../Registro%20V3/Flujos/flow-E17-cuando-agrega-o-elimina-una-entidad.md) con `Kind = 'acudida'`; aquí solo se resume lo propio del módulo.

INPUT: cascada Departamento* → Ciudad* → Municipio* → Sector* y «Buscar entidad»* (REG-19) · acciones completadas y pendientes (chips) · «Número de radicado…»* · «Fecha en la que reportó a la entidad»* · «Otra información asociada a la Ruta»* · [con acciones pendientes] acuerdo* → debida diligencia* → canales* y «¿Qué gestiones se realizaron…?»* (D17)

PASO 1 — «+ Siguiente entidad» se habilita con la sede elegida y sin faltantes (`entidadFaltantes()`; si falta algo se lista «Falta: …»)
PASO 2 — [POST] /api/v1/victim-cases-v3/{icode}/entidades  { Kind: 'acudida', EntityBranchID, FileNumber, ReportDate, OtherInfo, CompletedActions, PendingActions, Agreement, DueDiligence, ActivationChannels, ManagementNotes, Solicitud: '' }
  Agreement/DueDiligence son `true` solo si hay acciones pendientes y se respondió «Sí» (y «Sí»); los canales solo viajan si ambas son «Sí»
  Backend: Caso en Borrador → se agrega al JSON (`d~<icode>~<uuid>`, 409 si la sede ya está); Caso ya Iniciado → 409 `readOnly`
PASO 3 — Cliente: ok → el formulario se reinicia (tipo `acudida`) y `recargarEntidadesBarreras()`; 409 `readOnly` → `entitiesReadOnly = true`
PASO 4 — Al iniciarse el Caso: `entity_case` (+ `entity_case_obligation`) y, con acuerdo y debida diligencia en «Sí» y acciones pendientes, las tareas y oficios de cada canal (igual que una activación; ver el [sub-flujo](../../Registro%20V3/Flujos/flow-subflujo-materializar-borrador.md)). Sin acciones pendientes no se crea nada

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| «Otra ¿cuál?» (entidad que no está en el catálogo): fuera de alcance («entidades genéricas») | PASO 1 |
| La fecha de reporte solo vive en `entityMeta` | PASO 2 |
