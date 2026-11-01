━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega o elimina una Barrera
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: abrirBarreraParaEntidad(e) · cerrarBarreraParaEntidad() · agregarBarrera() · eliminarBarrera(id) · buscarCompetencias() → registro-v3-app.js · AddBarrier() / RemoveBarrier() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Las barreras se guardan en el modelo real (`salvia.barrier_v2`, el mismo de Seguimiento y la sesión psicosocial) y
> **cuelgan de una entidad ya agregada** (sub-repeater, decisión del 2026-10-15): no repiten la búsqueda geográfica;
> heredan sede, sector y ubicación. No hay una barrera sin entidad.

INPUT (agregar): {
  e:  entidad de la tarjeta   → entidades[i] (entityBranchId, sector interno he/js/pt/os)
  Specific[]:      barreras específicas del sector (códigos 01…)      → checklist (solo salud/justicia/protección)
  OtherDetail:     "¿Cuál otra barrera?"                              → <input>, solo sector otras_instituciones
  Institutional[] · Economic[] · Territorial[] · Differential[]       → 4 checklists transversales
  Competences[]:   competencias de la entidad                         → checklist (GET /api/v1/competences?sector=)
  Date · Official · Description:  fecha, funcionario/dependencia y descripción   → date, input, textarea (los tres obligatorios)
  Management[]:    "Gestión realizada" (códigos 01…06)                → checklist
}

PASO 1 — Abrir el formulario de una barrera ("+ Agregar barrera a esta entidad")
  Solo se ofrece si GS_17 = Sí y no hay otro formulario abierto.
  nuevaBarrera = barreraVacia(); EntityBranchID = e.entityBranchId
  SI la entidad se agregó en esta sesión (hay geo en entidadGeoPorBranch): Department/City/Town y SectorCode salen de esa búsqueda
  SI NO (entidad cargada del backend): SectorCode = traducción de e.sector (he→salud, js→justicia, pt→proteccion, os→otras_instituciones); Department/City/Town quedan vacíos
  barreraParaEntidad = e.entityBranchId
  buscarCompetencias(): [GET] /api/v1/competences?sector=<SectorCode> → competenciasBarrera

  Según el sector el formulario muestra: salud / justicia / protección → "Barreras específicas del sector" (12 / 18 / 13 opciones); otras_instituciones → "¿Cuál otra barrera?". Las 4 categorías transversales y la gestión salen siempre.

PASO 2 — Agregar ("+ Agregar barrera", habilitado solo con Descripción, Fecha y Funcionario)
  [POST] /api/v1/victim-cases-v3/{icode}/barreras
  // payload: nuevaBarrera completa con EntityBranchID numérico
  { "Department", "City", "Town", "SectorCode", "EntityBranchID", "Specific", "OtherDetail", "Institutional", "Economic", "Territorial", "Differential", "Date", "Official", "Description", "Management", "Competences" }

  Backend (AddBarrier) → barrier_v2:
    case_id = icode · status = OPEN (BarrierV2StatusOpen) · sector = SectorCode · entity_branch_id (si ≠ 0)
    listas → texto separado por comas (specific_barriers, structural_*, competences)
    management_actions: códigos 01…06 traducidos a slugs (orientacion_llamada, gestion_llamada, activacion_ruta_interinstitucional, articulacion_institucional, escalamiento_organismo_control, alerta_barreras)
    Respuesta: { "icode": <id de barrier_v2> }   (error de servicio: 500 "error al registrar la barrera")
    NOTA: no se crean `case_task` ni `entity_letter` por la gestión elegida (decisión de alcance 2026-10-15)

  Cliente: SI error → errorBarrera = <error>; SI ok → cerrarBarreraParaEntidad() y recargarEntidadesBarreras()

PASO 3 — Cancelar ("Cancelar"): cerrarBarreraParaEntidad() descarta el formulario sin llamar al backend

PASO 4 — Eliminar una barrera (botón "Eliminar" de la tabla de la entidad)
  window.confirm("¿Eliminar esta barrera del registro?")
  SI confirma: [DELETE] /api/v1/victim-cases-v3/barreras/{id} → recargarEntidadesBarreras()

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Sin mínimo de barreras: GS_17 = Sí no obliga a registrar ninguna | PASO 1 |
| Las barreras de una entidad eliminada quedan huérfanas (G-V3-22, ver flow-E17) | PASO 4 |
| La tabla de barreras de la tarjeta solo muestra descripción y la gestión en crudo (slugs) | PASO 2 |
