━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega o elimina una Entidad
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: agregarEntidad() · eliminarEntidad(relId) · onCambiaDeptoEntidad() · onCambiaCiudadEntidad() · buscarEntidadesEntidad() · onCambiaEntidadSeleccionada() → registro-v3-app.js · AddEntity() / RemoveEntity() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Las entidades se guardan en el modelo real de **Gestión institucional** (`salvia.entity_case` y
> `entity_case_obligation`), no en un repeater de `dinamic-form`. Cada alta es una llamada al backend (no hay "guardar
> todo" ni `registrar-lote`).

INPUT (agregar): {
  nuevaEntidad.Kind:           'acudida' | 'activacion'                        → <select> "Tipo"
  depto → ciudad → municipio:  cascada de ubicación                            → selects (solo filtro de búsqueda, no se envían)
  sector:                      salud | justicia | proteccion | otras_instituciones   → <select> "Sector" (solo filtro)
  EntityBranchID:              sede elegida (entity_branch)                    → <select> "Entidad"
  FileNumber:                  N° radicado, opcional                           → <input>
  ActivationChannel:           no | ll | ms (solo Kind 'activacion')           → <select> "Canal de activación"
  OtherInfo:                   "Objetivo con la entidad" (acudida) o "Información adicional" (activación)   → <input>
  CompletedActions:            ids de obligaciones completadas (solo 'acudida')   → checklist
  PendingActions:              ids de obligaciones pendientes                      → checklist
}

PASO 0 — Buscar la sede (cada select dispara su carga)
  Departamento → [GET] /api/v1/locations/cities?department_id=<depto>        → entidadCiudades
  Ciudad       → [GET] /api/v1/locations/towns?city_id=<ciudad>              → entidadMunicipios
  Municipio + Sector → [GET] /api/v1/entity-branches?town_code=<municipio>&sector=<sector>   → entidadesDisponibles
  Cambiar un nivel limpia los inferiores y la entidad elegida
  Entidad → [GET] /api/v1/entity-branches/<EntityBranchID>/obligations → obligacionesEntidad (se excluye "otra")

PASO 1 — Agregar una entidad ("+ Siguiente entidad", habilitado con una sede elegida)
  [POST] /api/v1/victim-cases-v3/{icode}/entidades
  // payload:
  {
    "Kind": …, "EntityBranchID": Number(EntityBranchID), "FileNumber": …,
    "ActivationChannel": Kind === 'activacion' ? ActivationChannel : "",
    "CompletedActions": Kind === 'acudida' ? CompletedActions : [],
    "PendingActions": PendingActions, "OtherInfo": …
  }

  Backend (AddEntity):
    SI la sede ya está activa en el Caso (EntityCaseRepository.ExistsActive): 409 { "error": "Esta entidad ya está asociada a este caso" }
    entity_case ← { case_id = icode, entity_branch_id, objetivo = OtherInfo, solicitud = "Canal de activación: <código>" (si hay canal), radicado = FileNumber, created_by = agente }
    PARA CADA id de CompletedActions (status 'completada') y PendingActions (status 'pendiente'): entity_case_obligation si no existía; "otra" y vacíos se ignoran
    Respuesta: { "icode": <id de entity_case> }
    NOTA: el tipo (acudida/activación) no se guarda como columna; solo deja huella el canal en `solicitud`. No se crean `case_task` ni `entity_letter`.

  Cliente:
    SI error: errorEntidad = <error> (p. ej. el 409) y no cambia nada más
    SI ok: entidadGeoPorBranch[sede] = {depto, ciudad, municipio, sector} (para las barreras); nuevaEntidad se reinicia; recargarEntidadesBarreras() → GET /api/v1/victim-cases-v3/{icode} → entidades y barreras

PASO 2 — Eliminar una entidad ("Eliminar entidad" en su tarjeta)
  window.confirm("¿Eliminar esta entidad del registro? También se eliminarán sus barreras asociadas.")
  SI cancela: TERMINAR ejecución
  SI confirma:
    [DELETE] /api/v1/victim-cases-v3/entidades/{relId}   → baja del entity_case (borrado del repositorio genérico)
    recargarEntidadesBarreras()

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| **G-V3-22** — el texto del `confirm` promete borrar las barreras asociadas, pero `RemoveEntity` solo elimina el `entity_case`: sus barreras siguen en `barrier_v2`, dejan de verse bajo una tarjeta (se agrupan por sede) y siguen contando en la tarjeta "Barreras" del Cierre. Por lectura de código; no verificado en navegador | PASO 2 |
| Sin automatización de `case_task`/`entity_letter` por canal de activación (sí la hay en Seguimiento): decisión de alcance del 2026-10-15 | PASO 1 |
| Sin máximo de entidades por Caso | PASO 1 |
| Fallas de `/locations`, `/entity-branches` y `/obligations` se tragan en silencio (listas vacías) | PASO 0 |
