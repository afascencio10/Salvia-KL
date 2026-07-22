# flow-E02 — Cuando abre modal crear o editar dupla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre modal crear o editar dupla
   Tipo: User Interaction
   Función: openDuplaModal(mode, dupla?)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  mode:   'create' | 'edit'                         → botón "+ Nueva dupla" o "Editar"
  dupla:  objeto de la fila (solo si mode === 'edit') → fila de "Duplas activas"
  psychologists:   lista cargada en E01             → state
  socialWorkers:   lista cargada en E01             → state
  duplas:          lista cargada en E01             → state
}

Precondiciones:
  isLoading === false
  isSaving === false
  modal no está ya en proceso de guardado


PASO 1 — Determinar modo del modal

SEGÚN mode:
  CASO 'create':
    → modal.mode = 'create'
    → modal.title = 'Nueva dupla'
    → modal.primaryLabel = 'Crear dupla'   // o "Guardar" según copy final
    → editingDuplaId = null
    → CONTINÚA

  CASO 'edit':
    → modal.mode = 'edit'
    → modal.title = 'Editar dupla'
    → modal.primaryLabel = 'Guardar cambios'
    → editingDuplaId = dupla.id
    → CONTINÚA

  DEFAULT:
    → TERMINAR (modo inválido)


PASO 2 — Inicializar formulario

SI mode === 'create':
  → form.name = ''                         // nombre libre; lo escribe el supervisor
  → form.psychologistId = ''
  → form.socialWorkerId = ''

SI mode === 'edit':
  → form.name = dupla.name
  → form.psychologistId = dupla.psychologistId
  → form.socialWorkerId = dupla.socialWorkerId

saveError = null
warningPs = null
warningTs = null


PASO 3 — Calcular profesionales disponibles para los selects

  Regla: un profesional está "ocupado" si pertenece a una dupla activa
  distinta de la que se está editando.

  occupiedIds = set vacío
  PARA CADA d EN duplas:
    SI editingDuplaId != null Y d.id === editingDuplaId:
      → omitir (los miembros actuales siguen seleccionables)
    SI NO:
      → occupiedIds.add(d.psychologistId)
      → occupiedIds.add(d.socialWorkerId)

  availablePsychologists = psychologists.filter(p => !occupiedIds.has(p.icode))
  availableSocialWorkers = socialWorkers.filter(p => !occupiedIds.has(p.icode))

  // En edit, los miembros actuales ya están incluidos porque su dupla se omitió


PASO 4 — Avisos de disponibilidad (UI)

SI availablePsychologists.length === 0:
  → warningPs = 'Todas las psicólogas están asignadas a una dupla.'
SI NO:
  → warningPs = null

SI availableSocialWorkers.length === 0:
  → warningTs = 'Todas las trabajadoras sociales están asignadas a una dupla.'
SI NO:
  → warningTs = null


PASO 5 — Abrir modal

  modal.kind = 'form'
  modal.visible = true
  → Vue muestra el modal con título, campos y avisos
  → Selects poblados con availablePsychologists / availableSocialWorkers

→ FIN EJECUCIÓN ✓
```

---

## canSave (computado, usado en E04 / UI)

```
canSave =
  form.name.trim().length > 0
  AND form.psychologistId !== ''
  AND form.socialWorkerId !== ''
  AND !isSaving
```

En modo `create`, si no hay opciones en algún select, `canSave` queda en false (no hay forma de completar el par).

---

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Paso afectado |
|---|---|
| Copy exacto del botón primario en create (“Crear dupla” vs “Guardar”) | PASO 1 |
| ¿Recargar profesionales/duplas desde API al abrir el modal (vs usar state de E01)? | PASO 3 |
