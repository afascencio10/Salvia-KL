━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma la reasignación (guardado)
   Tipo: User Interaction / Backend write
   ID: RRM-05
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en "Reasignar" (.rrm-btn--primary)

Precondiciones:
  canConfirm === true
  remisiones.length > 0
  loadingOptions === false
  optionsError === null
  saving === false


INPUT (frontend): {
  remisionIds:     remisiones.map(r => r.id)
  assignMode:      asignarEnDupla ? 'dupla' : 'professional'
  professionalId:  selectedProfessionalId  // solo si assignMode === 'professional'
  duplaId:         selectedDuplaId         // solo si assignMode === 'dupla'
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Validar canConfirm

  SI !canConfirm → TERMINAR


PASO 2 — Confirmación opcional (Swal)

  "¿Reasignar {N} remisión(es) a {nombre seleccionado}?"
  SI cancela → TERMINAR


PASO 3 — Guardado

  saving = true
  saveError = null

  POST /api/v1/psychosocial-support/reasignar-bulk
  Body: {
    remision_ids:    remisionIds,
    assign_mode:     assignMode,
    professional_id: assignMode === 'professional' ? professionalId : null,
    dupla_id:        assignMode === 'dupla' ? duplaId : null
  }


PASO 4 — Respuesta

  SI error:
    saving = false
    saveError = mensaje
    → TERMINAR

  SI ok:
    emit('reassigned', {
      remisiones:  remisiones,
      assignMode:  assignMode,
      updated:     result.remisiones_updated
    })
    → close()  // RRM-04
    → FIN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```javascript
onReasignarCompletado: function(payload) {
  this.$refs.remisionesComponent.reload();
  // Toast opcional: "{payload.updated} remisión(es) reasignada(s)"
}
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST reasignar-bulk
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivos sugeridos:
  src/salvia/controller/psychosocial_reassign_controller.go
  src/salvia/service/psychosocial_reassign_service.go
  src/internal/repository/psychosocial_reassign_repository.go

Response 200:
```json
{
  "ok": true,
  "remisiones_updated": 2,
  "team_contacts_updated": 5
}
```


PASO B1 — Validación

  → Sesión válida (supervisor `sv`)
  → remision_ids no vacío
  → assign_mode ∈ { 'professional', 'dupla' }
  → SI professional: professional_id presente; usuario activo rol ps o ts
  → SI dupla: dupla_id presente; dupla activa con integrantes activos


PASO B2 — Transacción

  Iniciar transacción
  remisionesUpdated = 0
  teamContactsUpdated = 0

  Para cada remisionId en remision_ids:

    ── B2.1 — Verificar remisión ──
    Tabla: salvia.psychosocial_support
    SI no existe → Rollback + 404/400

    ── B2.2 — Actualizar psychosocial_support ──

    Modo professional:
      UPDATE salvia.psychosocial_support
      SET professional_id = $professionalId,
          dupla_id        = NULL,
          updated_at      = NOW()
      WHERE id = $remisionId AND deleted_at IS NULL

    Modo dupla:
      UPDATE salvia.psychosocial_support
      SET dupla_id        = $duplaId,
          professional_id = NULL,
          updated_at      = NOW()
      WHERE id = $remisionId AND deleted_at IS NULL

    remisionesUpdated++

    ── B2.3 — Actualizar team_contact pendientes ──
    Tabla: salvia.team_contact
    Modelo: models.TeamContact

    Criterio: psicosocial_id = remisionId AND is_completed = false

    Modo professional:
      UPDATE salvia.team_contact
      SET professional_id = $professionalId,
          dupla_id        = NULL,
          updated_at      = NOW()
      WHERE psicosocial_id = $remisionId::text
        AND is_completed = false
        AND deleted_at IS NULL

    Modo dupla:
      UPDATE salvia.team_contact
      SET dupla_id        = $duplaId,
          professional_id = NULL,
          updated_at      = NOW()
      WHERE psicosocial_id = $remisionId::text
        AND is_completed = false
        AND deleted_at IS NULL

    teamContactsUpdated += filas afectadas

  Commit transacción
  SI falla → Rollback completo


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  TABLAS Y MODELOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Tabla | Modelo | Operación |
|---|---|---|
| salvia.psychosocial_support | PsychosocialSupport | UPDATE professional_id / dupla_id (exclusivos) |
| salvia.team_contact | TeamContact | UPDATE solo WHERE is_completed = false |
| salvia.dupla | Dupla | READ (validación modo dupla) |
| security.general_user | — | READ (validación profesional) |

Nota: `team_contact.psicosocial_id` es VARCHAR; comparar con `psychosocial_support.id::text`.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  QUÉ NO SE MODIFICA
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Entidad | Motivo |
|---|---|
| team_contact con is_completed = true | Conservar asignación histórica de sesiones realizadas |
| psychosocial_support.status | La reasignación no cambia el estado del flujo |
| victim_case / follow_up_v2 | Fuera de alcance de este modal |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  DIAGRAMA
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```
[Reasignar] → confirmReassign() [RRM-05]
       │
       ▼
POST /api/v1/psychosocial-support/reasignar-bulk
       │
       ├── Por cada remisionId ──────────────────────────────┐
       │    1. UPDATE psychosocial_support (prof OR dupla)   │
       │    2. UPDATE team_contact (is_completed = false)    │
       └─────────────────────────────────────────────────────┘
       │
       ▼
emit('reassigned') → padre reload() remisiones-psicosocial-component
```
