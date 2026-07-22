# flow-E07 — Cuando confirma eliminar dupla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma eliminar dupla
   Tipo: User Interaction
   Función: confirmDeleteDupla()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  deleteTarget.id:    id de la dupla               → state (E05)
  deleteTarget.name:  nombre para mensajes         → state (E05)
}


PASO 1 — Guard clauses

SI deleteTarget es null:
  → TERMINAR

SI isDeleting === true:
  → TERMINAR


PASO 2 — Activar estado de borrado
  isDeleting = true
  deleteError = null


PASO 3 — Solicitar eliminado lógico

  DELETE /api/v1/duplas/{deleteTarget.id}


┌─────────────────────────────────────────────────────────────┐
│  SUB-FLUJO BACKEND: Validar uso y soft-delete               │
└─────────────────────────────────────────────────────────────┘

  B1. Buscar dupla por id donde deleted_at IS NULL
      SI no existe → 404

  B2. Verificar si la dupla está en uso en una sesión/remisión NO cerrada

      Condición de bloqueo (cualquiera de las dos):

      (A) Remisión directa con la dupla y status distinto de cerrado:
      ```sql
      SELECT 1
      FROM salvia.psychosocial_support ps
      WHERE ps.dupla_id = :duplaId
        AND ps.deleted_at IS NULL
        AND ps.status <> 'cerrado'   -- PsychosocialSupportStatusCerrado
      LIMIT 1
      ```

      (B) Sesión (team_contact) con la dupla, cuya remisión asociada
          no está cerrada:
      ```sql
      SELECT 1
      FROM salvia.team_contact tc
      JOIN salvia.psychosocial_support ps
        ON ps.id::text = tc.psicosocial_id
       AND ps.deleted_at IS NULL
      WHERE tc.dupla_id = :duplaId
        AND tc.deleted_at IS NULL
        AND ps.status <> 'cerrado'
      LIMIT 1
      ```

      SI (A) O (B) encuentra filas:
        → 409 {
            "error": "dupla_en_uso",
            "message": "Esta dupla no se puede eliminar porque se está utilizando en una sesión."
          }
        → FIN SUB-FLUJO (error) — NO se hace soft-delete

      SI NO hay uso activo:
        → CONTINÚA

      Nota: si solo existen referencias con ps.status = 'cerrado',
      la eliminación SÍ se permite. El dupla_id histórico puede quedar
      apuntando a la dupla soft-deleted (auditoría).

  B3. Soft delete
      ```sql
      UPDATE salvia.dupla
      SET deleted_at = NOW(), updated_at = NOW()
      WHERE id = :id AND deleted_at IS NULL
      ```

  B4. NO modificar security.general_user ni nullificar dupla_id
      en remisiones/sesiones históricas.

  B5. Respuesta 200 { "ok": true }

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


PASO 4 — Manejar respuesta

SI status === 200:
  → isDeleting = false
  → Cerrar modal de confirmación (limpieza como E06)
  → Toast: "Dupla eliminada. Los profesionales quedaron disponibles."
  → Recargar datos → Ver flujo: Cuando carga la pantalla (E01)
  → FIN ✓

SI status === 409 (dupla en uso):
  → isDeleting = false
  → Cerrar modal de confirmación (limpieza deleteTarget / kind delete)
  → Abrir modal de error:
      modal.kind = 'deleteBlocked'
      modal.visible = true
      blockedMessage = response.message
        || 'Esta dupla no se puede eliminar porque se está utilizando en una sesión.'
  → Vue muestra modal de error (ver interfaz)
  → Cierre del modal → Ver evento: Cuando cierra modal de error al eliminar (E08)
  → TERMINAR

SI status === 404:
  → isDeleting = false
  → Cerrar modal de confirmación
  → Toast / aviso: 'La dupla ya no existe o fue eliminada.'
  → Recargar datos → E01
  → TERMINAR

SI status === 401:
  → Redirigir a login
  → TERMINAR

SI otro error:
  → isDeleting = false
  → deleteError = 'No se pudo eliminar la dupla. Intenta de nuevo.'
  → Modal de confirmación permanece abierto con el mensaje
  → TERMINAR
```

---

## Efecto en UI tras soft-delete exitoso (recarga E01)

| Antes | Después |
|---|---|
| Miembros con badge “En dupla” | Badge “Disponible” |
| Fila en “Duplas activas (N)” | Desaparece; contador N-1 |
| Aparecen de nuevo en selects de E02 create | Sí |

---

## Modelos involucrados en la validación

| Tabla / modelo | Campo | Criterio |
|---|---|---|
| `salvia.psychosocial_support` | `dupla_id`, `status` | Bloquea si `status ∈ { abierto, en_gestion, en_devolucion }` |
| `salvia.team_contact` | `dupla_id`, `psicosocial_id` | Bloquea si el `psychosocial_support` joinado no está `cerrado` |

Estados de remisión (`psychosocial_support.go`):
`abierto` | `en_gestion` | `en_devolucion` | `cerrado`

---

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Paso afectado |
|---|---|
| ¿Cast `ps.id::text = tc.psicosocial_id` o unificar tipos en BD? | B2-(B) |
| Copy exacto del toast de éxito | "Dupla eliminada. Los profesionales quedaron disponibles." |
| Título exacto del modal de error (“No se puede eliminar” vs otro) | "No se puede eliminar" (interfaz) |
