━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda el formulario de nueva entidad
   Tipo: User Interaction
   Función: guardarEntidad()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  entidadId:  entity_branch_id elegido    → <select> Entidad del modal (poblado en flujo E06)
  objetivo:   texto libre                 → <textarea> del modal
  caseId:     icode del caso              → prop del componente
  userId:     icode del agente en sesión  → prop del componente
}

Validar formulario:
  • entidadId:  required
  • objetivo:   required, min_chars: 1

  SI validación con errores:
    → Mostrar error por campo (errores.entidad / errores.objetivo)
    → TERMINAR ejecución

  SI validación ok:
    → CONTINÚA FLUJO GENERAL

PASO 1 — Marcar estado de guardado
  guardando = true

PASO 2 — Crear la relación caso-entidad
  RelEntityBranchVictimCase.crear({
    entityBranchId: formNuevaEntidad.entidadId,   // del <select> Entidad
    objetivo:       formNuevaEntidad.objetivo,     // del <textarea>
    createdById:    userId                          // prop del componente
  })

  [MÉTODO] POST /api/v1/casos/{caseId}/entidades

  // payload:
  {
    "entityBranchId": formNuevaEntidad.entidadId,   // origen: selección en el modal (E06)
    "objetivo":        formNuevaEntidad.objetivo,    // origen: input del usuario
    "createdById":     userId                         // origen: prop userId
  }

  Backend ejecuta:
    1. Verificar que no exista ya una relación activa con el mismo (case_id, entity_branch_id)
       — índice único parcial ux_rel_entity_branch_victim_case_active
    2. INSERT en rel_entity_branch_victim_case

SI el backend responde error 409 (relación duplicada):
  → guardando = false
  → Mostrar error "Esta entidad ya está asociada a este caso"
  → TERMINAR ejecución (el modal permanece abierto para que el usuario corrija)

SI el backend responde error genérico (500 u otro):
  → guardando = false
  → Mostrar error genérico de guardado
  → TERMINAR ejecución

SI respuesta ok (201):
  → guardando = false
  → modalAgregarAbierto = false
  → CONTINÚA FLUJO GENERAL

PASO 3 — Recargar la lista de entidades del caso
→ Ver flujo: Cuando carga el componente (E01) — se re-ejecuta para reflejar la nueva entidad
  (incluye su oficiosCount = 0 y ultimaAccionFecha = null, ya que aún no tiene oficios)

PASO 4 — Mostrar confirmación de éxito (toast / mensaje)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                    | Paso afectado |
|-------------------------------------------------------------------------------------------|---------------|
| El endpoint POST /api/v1/casos/:caseId/entidades no existe — hay que crearlo junto con la validación de duplicados (409). | PASO 2 |
| Confirmar mensaje/UX exacto de éxito (toast, banner, etc.) — no especificado en el mockup original. | PASO 4 |
