━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega o elimina una Barrera
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — actualizado 2026-11-04 (WP-4: REG-18, REG-04 c, REG-38*)
   Funciones: abrirBarreraParaEntidad(e) · cerrarBarreraParaEntidad() · setBarrera() · toggleBarreraOpcion() · barreraPreguntasVisibles() · barreraFaltantes() · armarPayloadBarrera() · agregarBarrera() · eliminarBarrera(id) · etiquetaGestionBarrera(csv) → registro-v3-app.js · AddBarrier() / RemoveBarrier() → victim_case_v3_service.go · addBarrierDraft() / removeBarrierDraft() → victim_case_v3_draft.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **WP-4.** El formulario es **el de Seguimiento v4** (REG-18): se pinta recorriendo `BARRERA_V4_V3`, un catálogo
> generado del seed de Seguimiento (`scripts/gen_registro_v3_barrera_catalog.py`); ya **no** se capturan competencias.
> Las barreras **cuelgan de una entidad** (tarjeta de ambas listas, D8: mientras REG-29 no se defina) y, mientras el Caso
> es Borrador, se guardan en el **borrador JSON**; al iniciarse el Caso se materializan en `barrier_v2` con las tareas y
> oficios de su gestión ([sub-flujo](flow-subflujo-materializar-borrador.md)). Después son **solo lectura**.

INPUT (agregar): {
  e:  entidad de la tarjeta   → EntityBranchID = e.entityBranchId
  Sector*:                      Salud | Justicia | Protección | Transversales
  por sector:                   «Barreras en …» (multiple) · «¿A qué institución acudió?» · «Otra barrera en sector …»   → Specific[] · Institution · OtherDetail
  Transversales:                4 listas (institucionales y de talento humano · económicas y socioeconómicas · territoriales y geográficas · por ausencia de enfoque diferencial) → Institutional[] · Economic[] · Territorial[] · Differential[]
  Fecha aproximada* · Funcionario/a o dependencia* · Departamento* → Municipio* · Vereda o zona* (cabecera | rural | no) · Descripción (tiempo, modo y lugar)*   → Date · Official · Department · City · Zone · Description
  Gestión de la barrera*:       activacion («Activación de ruta interinstitucional (Generar oficio)») | alerta («Alerta por barreras»), múltiple   → Management[]
  Perspectiva profesional*:     solo si la gestión incluye alerta (+ aviso «Alerta por barreras notificaría al equipo técnico…»)   → ProfessionalPerspective
  Descripción de la gestión*:   → ManagementDescription
}

PASO 1 — Abrir el formulario ("+ Agregar barrera a esta entidad")
  Solo se ofrece si GS_17 = Sí, el Caso sigue en Borrador (!entitiesReadOnly) y no hay otro formulario abierto.
  nuevaBarrera = barreraVacia(); EntityBranchID = e.entityBranchId; barreraParaEntidad = e.entityBranchId
  Visibilidad por condiciones del catálogo (`cuando`): las preguntas de cada sector, las 4 transversales y la perspectiva aparecen según lo elegido; cambiar el sector limpia lo que dejó de aplicar; elegir un Departamento carga sus municipios ([GET] /api/v1/locations/cities?department_id=)

PASO 2 — Agregar ("+ Agregar barrera", deshabilitado mientras `barreraFaltantes()` liste obligatorias visibles sin responder)
  [POST] /api/v1/victim-cases-v3/{icode}/barreras
  // payload: { EntityBranchID, SectorCode, Specific, Institution, OtherDetail, Institutional, Economic, Territorial, Differential, Date, Official, Department, City, Zone, Description, Management, ProfessionalPerspective, ManagementDescription }
  // (Competences y los códigos de gestión 01–06 de clientes anteriores se siguen aceptando)

  Backend (AddBarrier) — estado leído de victim_case_status:
    SI el Caso NO está en Borrador: 409 { "error": "La sección de entidades y remisiones ya no es editable…", "readOnly": true }
    SI NO: la barrera se agrega al JSON (id "d~<icode>~<uuid>")
    Respuesta: { "icode": <id del borrador> }
    Al materializar → barrier_v2: status OPEN · sector · specific_barriers · specific_institutions = Institution · other_barrier_desc · structural_* · barrier_date · official_dependency · department_id · city_id · description · management_actions (CSV) · entity_branch_id; Zone, ProfessionalPerspective y ManagementDescription van a `barrierMeta[id]`
    Tareas y oficios: gestión `activacion` → oficio (por proyectar) + tarea Barreras · proyectar_oficio «Activación de ruta interinstitucional (Generar oficio)»; gestión `alerta` → tarea Barreras · comite_caso «Alerta por barreras» (sin responsable)

  Cliente: SI error → errorBarrera = <error> (y, si es 409 readOnly, entitiesReadOnly = true); SI ok → cerrarBarreraParaEntidad() y recargarEntidadesBarreras()

PASO 3 — Cancelar ("Cancelar"): cerrarBarreraParaEntidad() descarta el formulario sin llamar al backend

PASO 4 — Eliminar una barrera (botón "Eliminar" de la tabla, solo si !entitiesReadOnly)
  window.confirm("¿Eliminar esta barrera del registro?")
  SI confirma: [DELETE] /api/v1/victim-cases-v3/barreras/{id} → id "d~…": se quita del JSON; UUID de tabla: 409 readOnly; inexistente: 404 → recargarEntidadesBarreras()

PASO 5 — La columna «Gestión» de la tabla muestra etiquetas (REG-04 c): `etiquetaGestionBarrera(managementActions)` traduce los 2 valores de Seguimiento v4 y los 6 slugs históricos, deja tal cual un valor desconocido y muestra «—» si no hay

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Sin mínimo de barreras: GS_17 = Sí no obliga a registrar ninguna | PASO 1 |
| REG-29 (dónde se gestionan las barreras) y REG-31 (textos finales) siguen pendientes del cliente: el botón aparece en las tarjetas de ambas listas (D8) y los textos son los de Seguimiento v4 | PASO 1 |
| Los valores `activacion` y `alerta` no tienen etiqueta en `barrera_detalle.html`, `barrier-follow-up-timeline.js`, `psicosocial_barreras.go`, `followup_v2_service.go` ni `form_service.go` («Para integración») | PASO 5 |
| «Alerta por barreras» solo crea la tarea del comité: la notificación al equipo técnico todavía no se envía (igual que Seguimiento v4) | PASO 2 |
