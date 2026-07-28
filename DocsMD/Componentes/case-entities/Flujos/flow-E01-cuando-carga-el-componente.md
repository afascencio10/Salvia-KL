━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Función: mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  caseId:  icode del caso (victim_case_i_code)   → prop inyectada por el padre (case_detail)
}

PASO 1 — Inicializar estado
  cargando = true, error = null, entidades = []

PASO 2 — Consultar entidades del caso
  GET /api/v1/casos/{caseId}/entidades

  Backend ejecuta:
    1. Query entity_case WHERE case_id = caseId AND deleted_at IS NULL
    2. JOIN entity_branch (por entity_branch_id) + entity (por entity_branch.entity_id) → nombre, sector, dirección, lat/long
    3. JOIN town → city → department (por entity_branch.entity_branch_town_code) → nombres de ubicación
    4. Para cada relación: COUNT(entity_letter) WHERE case_id = caseId AND entity_branch_id = X → oficiosCount
    5. Para cada relación: COUNT(barrier_v2) WHERE case_id = caseId AND entity_branch_id = X AND status IN ('OPEN','En Gestion') → barrerasActivasCount
    6. Para cada relación: entity_case.last_action → lastAction (se lee tal cual, no se calcula)
    7. ORDER BY entity_case.created_at ASC

SI respuesta con error HTTP:
  → error = mensaje de error
  → cargando = false
  → Pantalla muestra ErrorMsg + BtnReintentar
  → TERMINAR ejecución

SI respuesta ok (200):
  → entidades = resultado (puede ser array vacío)
  → cargando = false
  → CONTINÚA FLUJO GENERAL

PASO 3 — Vue renderiza
  SI entidades.length === 0:
    → Muestra EmptyState "No hay entidades registradas para este caso"
  SI NO:
    → Muestra Grid de EntidadCard

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                        | Paso afectado |
|-----------------------------------------------------------------------------|---------------|
| El endpoint GET /api/v1/casos/:caseId/entidades no existe — hay que crearlo, junto con el join descrito. | PASO 2 |
| barrerasActivasCount será 0 para todas las entidades hasta que se implemente el cambio (fuera de alcance) que permite elegir entity_branch_id al registrar una barrera. | PASO 2.5 |
| Falta definir el mecanismo para poblar/editar entity_case.last_action (¿formulario propio? ¿se actualiza en otro evento?). | PASO 2.6 |
