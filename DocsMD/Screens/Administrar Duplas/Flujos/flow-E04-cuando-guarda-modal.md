# flow-E04 — Cuando guarda modal

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda modal
   Tipo: User Interaction
   Función: saveDuplaModal()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  modal.mode:        'create' | 'edit'           → state (E02)
  editingDuplaId:    id | null                   → state (E02)
  form.name:         string (libre)              → input usuario
  form.psychologistId: general_user_i_code       → select
  form.socialWorkerId: general_user_i_code       → select
}


PASO 1 — Guard clause

SI isSaving === true:
  → TERMINAR

SI modal.kind !== 'form' O modal.visible !== true:
  → TERMINAR


PASO 2 — Validar formulario (frontend)

Validar:
  • name:            required, trim, max 36 chars  // nombre libre; no se autogenera
  • psychologistId:  required
  • socialWorkerId:  required
  • psychologistId !== socialWorkerId   // defensa; roles distintos en UI

SI validación con errores:
  → Mostrar error por campo / toast
  → TERMINAR ejecución

SI validación ok:
  → CONTINÚA FLUJO GENERAL


PASO 3 — Activar estado de guardado
  isSaving = true
  saveError = null


PASO 4 — Persistir según modo

SEGÚN modal.mode:

  CASO 'create':
    POST /api/v1/duplas
    // payload:
    {
      "name":            form.name.trim(),          // input libre del supervisor
      "psychologistId":  form.psychologistId,       // select
      "socialWorkerId":  form.socialWorkerId        // select
    }
    → CONTINÚA

  CASO 'edit':
    PUT /api/v1/duplas/{editingDuplaId}
    // payload:
    {
      "name":            form.name.trim(),
      "psychologistId":  form.psychologistId,
      "socialWorkerId":  form.socialWorkerId
    }
    → CONTINÚA


┌─────────────────────────────────────────────────────────────┐
│  SUB-FLUJO BACKEND: Crear o actualizar dupla                │
└─────────────────────────────────────────────────────────────┘

  B1. Validar payload (mismos campos required + max length)

  B2. Verificar unicidad de nombre entre duplas activas
      ```sql
      SELECT id FROM salvia.dupla
      WHERE deleted_at IS NULL
        AND name = :nameTrimmed
        AND (:editingId IS NULL OR id <> :editingId)
      ```
      SI existe:
        → 409 { "error": "nombre de la dupla en uso" }
        → FIN SUB-FLUJO (error)

  B3. Verificar que psychologistId es usuario activo con role 'ps'
      Verificar que socialWorkerId es usuario activo con role 'ts'
      SI alguno inválido:
        → 400 { "error": "Profesional inválido o inactivo" }
        → FIN SUB-FLUJO (error)

  B4. Verificar unicidad de psicóloga en duplas activas
      (la trabajadora social SÍ puede repetirse en varias duplas)
      ```sql
      SELECT id FROM salvia.dupla
      WHERE deleted_at IS NULL
        AND BTRIM(psychologist_id::text) = BTRIM(:psychologistId)
        AND (:editingId IS NULL OR id <> :editingId)
      ```
      SI existe conflicto:
        → 409 { "error": "La psicóloga ya pertenece a otra dupla" }
        → FIN SUB-FLUJO (error)

  B5. Persistir
      CREATE:
        INSERT salvia.dupla { name, psychologist_id, social_worker_id }
      UPDATE:
        UPDATE salvia.dupla
        SET name, psychologist_id, social_worker_id, updated_at = now()
        WHERE id = :editingId AND deleted_at IS NULL
        SI 0 filas → 404

  B6. Respuesta 200/201 con la dupla enriquecida (mismos campos que GET lista)

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


PASO 5 — Manejar respuesta HTTP

SI status === 200 o 201:
  → isSaving = false
  → Cerrar modal (misma limpieza que E03)
  → Toast éxito: create → "Dupla creada" | edit → "Cambios guardados"
  → Recargar datos de pantalla → Ver flujo: Cuando carga la pantalla (E01)
  → FIN ✓

SI status === 409 Y mensaje contiene "nombre de la dupla en uso":
  → isSaving = false
  → saveError = 'nombre de la dupla en uso'
  → Resaltar campo Nombre (para que el supervisor lo cambie)
  → Modal permanece abierto
  → TERMINAR

SI status === 409 (conflicto de miembros):
  → isSaving = false
  → saveError = mensaje del API
  → Modal permanece abierto
  → TERMINAR

SI status === 400:
  → isSaving = false
  → saveError = mensaje de validación
  → Modal permanece abierto
  → TERMINAR

SI status === 401:
  → Redirigir a login
  → TERMINAR

SI otro error:
  → isSaving = false
  → saveError = 'No se pudo guardar la dupla. Intenta de nuevo.'
  → Modal permanece abierto
  → TERMINAR
```

---

## Notas de implementación

- Extender `DuplaRepository`: `Create`, `Update`, `SoftDelete`, `ListActiveEnriched`, `FindActiveConflictByMembers`, `FindActiveByName`.
- El endpoint `GET /api/v1/duplas/reasignacion` (RRM-03) sigue siendo de solo lectura; esta pantalla es la dueña del CRUD.
- Comparación de nombre: exacta sobre el valor trimmeado, solo entre filas con `deleted_at IS NULL`. Una dupla soft-deleted puede reutilizar el mismo nombre.

---

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Paso afectado |
|---|---|
| Paths definitivos de API | `POST /api/v1/duplas`, `PUT /api/v1/duplas/:id` (implementado) |
| ¿Comparación de nombre case-insensitive? (hoy: exacta tras trim) | B2 |
| Copy de toasts | create → "La dupla se creó exitosamente" / edit → "La dupla se editó exitosamente" |
