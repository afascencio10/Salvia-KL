━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega o elimina una Entidad
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — actualizado 2026-11-04 (WP-4: REG-17, REG-19, REG-20, REG-38*)
   Funciones: agregarEntidad() · eliminarEntidad(relId) · entidadFaltantes() · activacionAplica() · toggleCanal() · toggleAccion() · onInputBuscarEntidad() · seleccionarEntidadSugerida() · buscarEntidadesEntidad() · onCambiaEntidadSeleccionada() · marcarSoloLectura() → registro-v3-app.js · AddEntity() / RemoveEntity() → victim_case_v3_service.go · addEntityDraft() / removeEntityDraft() → victim_case_v3_draft.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **WP-4.** Mientras el Caso es **Borrador** (`bo`) las entidades se guardan en el **borrador JSON**
> (`victim_case.victim_case_v3_draft_data`), no en `entity_case`: al pasar el Caso a Iniciado se materializan (ver
> [flow-subflujo-materializar-borrador](flow-subflujo-materializar-borrador.md)). Después son **solo lectura**. El
> **tipo** lo fija el módulo: `acudida` en «Ruta de atención integral» (formulario de Seguimiento v4 «Identificación de
> entidades», D17) y `activacion` en «Plan de acción Salvia» (se despliega con GS_16 = Sí). Cada alta es una llamada al
> backend (no hay "guardar todo").

INPUT (agregar): {
  nuevaEntidad.Kind:           'acudida' (RA) | 'activacion' (GS)               → lo fija el módulo (watch de moduloActivo); ya no hay selector «Tipo»
  depto → ciudad → municipio → sector:  Departamento*, Ciudad*, Municipio*, Sector* en cascada obligatoria   → selects (solo filtro de búsqueda, no se envían)
  EntityBranchID:              sede elegida                                     → «Buscar entidad»* (texto con sugerencias; sustituye al select «Entidad»)
  FileNumber:                  N° de radicado (acudida: «Número de radicado (NUNC, medida de protección)…»*; activación: opcional)
  ReportDate:                  «Fecha en la que reportó a la entidad»* (solo acudida)
  OtherInfo:                   acudida «Otra información asociada a la Ruta»* · activación «Objetivo con la entidad»*   → objetivo
  Solicitud:                   «Solicitud a la entidad» (solo activación)
  CompletedActions / PendingActions:  chips de acciones institucionales (acudida: completadas y pendientes, excluyentes entre sí; activación: pendientes). Acudida: al menos una si la sede tiene acciones en el catálogo
  Agreement · DueDiligence:    acudida con acciones pendientes: preguntas propias (acuerdo* → debida diligencia*); activación: copias de GS_15 y GS_16 al agregar
  ActivationChannels[] · ManagementNotes:  canales* (Notificación SALVIA · Articulación… · Denuncia fácil · Mi Salvia · Ninguna, «Ninguna» excluyente) y «¿Qué gestiones se realizaron…?»*; solo si Agreement y DueDiligence = Sí
}

PASO 0 — Buscar la sede (REG-19)
  Departamento → [GET] /api/v1/locations/cities?department_id=<depto>        → entidadCiudades
  Ciudad       → [GET] /api/v1/locations/towns?city_id=<ciudad>              → entidadMunicipios
  Cada select está deshabilitado hasta elegir el anterior; cambiar un nivel limpia los inferiores, el texto de búsqueda y la entidad
  Municipio + Sector → UNA petición [GET] /api/v1/entity-branches?town_code=<municipio>&sector=<sector> → entidadesDisponibles
     (cada petición lleva un número: una respuesta atrasada, de un municipio o sector ya cambiado, se descarta)
  «Buscar entidad»: @input guarda el texto y, tras 250 ms sin teclear (debounce), filtra la lista YA cargada en el navegador (sin tildes ni mayúsculas, desde la primera letra); escribir nunca llama al servidor. Editar el texto después de elegir una sugerencia limpia la selección
  Elegir una sugerencia → [GET] /api/v1/entity-branches/<EntityBranchID>/obligations → chips (se excluye "otra"; sin acciones: «Esta entidad no tiene tareas configuradas en el catálogo.»)

PASO 1 — Agregar una entidad ("+ Siguiente entidad", habilitado con sede elegida y sin faltantes; si falta algo se lista «Falta: …»)
  [POST] /api/v1/victim-cases-v3/{icode}/entidades
  // payload: { Kind, EntityBranchID, FileNumber, Solicitud, ReportDate, ActivationChannels, Agreement, DueDiligence, ManagementNotes, CompletedActions, PendingActions, OtherInfo }
  // (ActivationChannel en texto —formato anterior no/ll/ms— se sigue aceptando)

  Backend (AddEntity) — estado leído de victim_case_status:
    SI el Caso NO está en Borrador: 409 { "error": "La sección de entidades y remisiones ya no es editable: el Caso dejó de ser borrador", "readOnly": true }
    SI la sede ya está en el borrador o activa en entity_case: 409 { "error": "Esta entidad ya está asociada a este caso" }
    SI NO: la entidad se agrega al JSON (id "d~<icode>~<uuid>", con los datos de la sede, `by` = agente) y `entityMeta` se completa al materializar
    Respuesta: { "icode": <id del borrador> }
    NOTA: no hay filas en entity_case/entity_case_obligation/case_task/entity_letter hasta materializar

  Cliente:
    SI error: errorEntidad = <error>; SI además es 409 readOnly: entitiesReadOnly = true (la pantalla pasa a solo lectura) y se recarga
    SI ok: nuevaEntidad se reinicia (mismo tipo); recargarEntidadesBarreras() → GET /api/v1/victim-cases-v3/{icode} → entidades, barreras, entityMeta, barrierMeta, entitiesReadOnly

PASO 2 — Eliminar una entidad ("Eliminar entidad", solo si !entitiesReadOnly)
  window.confirm("¿Eliminar esta entidad del registro? También se eliminarán sus barreras asociadas.")
  SI confirma: [DELETE] /api/v1/victim-cases-v3/entidades/{relId}
    Backend: id "d~…" → se quita del JSON la entidad Y sus barreras (solo en el borrador: aún no hay tareas ni filas); id de tabla (UUID) → 409 readOnly; id inexistente → 404
    recargarEntidadesBarreras()

PASO 3 — Qué se ve en la tarjeta (ambas listas): sector, nombre, ubicación, «Canales de activación» (de entityMeta o, en datos anteriores, de «Canal de activación: …» en solicitud), acciones institucionales con su estado (✓ Completada / Pendiente), tabla de barreras y «+ Agregar barrera a esta entidad» (flow-E18)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| **G-V3-22** — ✅ Resuelto (WP-4): en el borrador se eliminan juntas la entidad y sus barreras; después de materializar la sección es solo lectura | PASO 2 |
| «Otra ¿cuál?» (entidad que no aparece en el catálogo): fuera de alcance («entidades genéricas», pedido aparte del cliente) | PASO 0 |
| Sin máximo de entidades por Caso | PASO 1 |
| Fallas de `/locations`, `/entity-branches` y `/obligations` se tragan en silencio (listas vacías) | PASO 0 |
| La fecha de reporte y los datos de activación solo viven en `entityMeta` (Seguimiento tampoco los guarda en columna) | PASO 1 |
