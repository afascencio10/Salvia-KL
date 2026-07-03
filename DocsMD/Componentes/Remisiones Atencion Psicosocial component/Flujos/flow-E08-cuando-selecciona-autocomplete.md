━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona o limpia el autocomplete de profesional asignada
   Tipo: User Interaction
   Código: E-08
   Decisión relacionada: [DEC-E08-01 — filtro incluye dupla](./flow-decision-E08-filtro-profesional-incluye-dupla.md)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  (A) clic en SuggestionItem del dropdown
  (B) clic en "✕" del SelectedTag


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONTEXTO — Asignación directa vs dupla
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Un profesional puede estar vinculado a remisiones de dos maneras:

| Vía | Dónde se guarda |
|---|---|
| Directa | `psychosocial_support.professional_id` |
| Por dupla | `psychosocial_support.dupla_id` → `dupla.psychologist_id` o `dupla.social_worker_id` |

Al filtrar por profesional, el listado debe incluir **ambas** vías (ver DEC-E08-01).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — caso (A) seleccionar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Guardar selección

  autocompleteSelected['profesional_asignada'] = {
    value: opt.value,
    label: opt.label,
    roleLabel: opt.roleLabel
  }
  autocompleteText['profesional_asignada'] = ""
  autocompleteSuggestions = []

PASO 2 — Activar filtro y recargar

  activeFilters['professional_id'] = opt.value   // general_user_i_code
  currentPage                      = 1
  selectedRemisiones               = []
  loading                          = true

PASO 3 — fetchRemisiones() con todos los filtros activos

  GET /api/v1/psychosocial-support/list
    &filter_professional_id={opt.value}
    (+ demás filter_*)

  → El backend resuelve duplas del profesional; el frontend no envía dupla_id.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — caso (B) limpiar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Limpiar estado autocomplete

  autocompleteSelected['profesional_asignada'] = null
  autocompleteText['profesional_asignada']     = ""

PASO 2 — Quitar filtro y recargar

  delete activeFilters['professional_id']
  currentPage        = 1
  selectedRemisiones = []

PASO 3 — fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — cláusula WHERE (DEC-E08-01)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND (
  BTRIM(ps.professional_id::text) = BTRIM({filter_professional_id})
  OR ps.dupla_id IN (
    SELECT d.id
    FROM salvia.dupla d
    WHERE d.deleted_at IS NULL
      AND (
        BTRIM(d.psychologist_id::text) = BTRIM({filter_professional_id})
        OR BTRIM(d.social_worker_id::text) = BTRIM({filter_professional_id})
      )
  )
)
```

**Semántica:** “Todas las remisiones vinculadas a este profesional”, ya sea por asignación directa o como miembro de una dupla.

**Combinación con E-11:** si además hay `filter_dupla_id`, ambos filtros se aplican con AND (intersección).

> **Cambio post-reunión:** reemplaza `agent_id` / `filter_agent_id` de la versión anterior.
>
> **Cambio DEC-E08-01:** reemplaza el filtro simple `ps.professional_id = ?` que omitía remisiones solo asignadas vía dupla.
