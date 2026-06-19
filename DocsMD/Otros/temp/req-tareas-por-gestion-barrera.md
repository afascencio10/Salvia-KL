# Requerimiento — Crear tareas por gestión de barrera al completar seguimiento

> ⚠️ **MD temporal de planificación.** Pendiente de validación antes de implementar.

---

## ¿Qué pide el requerimiento?

Cuando se completa el formulario de seguimiento y el backend procesa la Sección 4 (Identificación de Barreras), por cada barrera identificada se deben crear tareas en `case_task` — **una tarea por cada opción seleccionada** en la pregunta de gestión (Q23 `gestión de la barrera`).

**Ejemplo:**
- Barrera #1 tiene seleccionadas: `orientacion_llamada`, `articulacion_institucional`
  → Se crean 2 tareas para esa barrera
- Barrera #2 tiene seleccionada: `escalamiento_organismo_control`
  → Se crea 1 tarea para esa barrera

---

## Cambios en la interfaz

**Ninguno.** El formulario ya recolecta la pregunta de gestión (Q23, tipo `multiple`) dentro del repeater de Sección 4. No se requiere agregar ni modificar ningún campo visible para el usuario.

---

## Eventos afectados

Solo se modifica **E-04**.

| Evento | Tipo | Cambio |
|---|---|---|
| E-04 — Cuando el backend procesa el seguimiento completado | Backend/Scheduled | Se extiende el PASO 3 para crear tareas por cada opción de gestión seleccionada en cada barrera |

---

## Cambio en el flujo de E-04

El flujo actual de E-04 en el PASO 3 crea un registro `barrier_v2` por cada entry del repeater de Sección 4. La modificación agrega un sub-paso **inmediatamente después** de crear cada `barrier_v2`:

### Flujo modificado — PASO 3 (extracto)

```
PASO 3 — Procesar barreras (Sección 4 repeater)

  ↺ PARA CADA entry del repeater de Sección 4:

    3.1 — Leer respuestas de la entry:
            sector           → Q1  (dropdown)
            gestion_values   → Q23 (multiple, CSV)

    3.2 — Crear registro barrier_v2:
            DB.Create(barrier_v2 {
              case_id:     caseId,
              follow_up_id: followUpId,
              sector:      sector,
              description: [...campos según MD existente...],
              status:      "OPEN"
            })
            → newBarrierID = barrier_v2.id  ← guardar para el paso siguiente

    3.3 — [NUEVO] Crear tareas por gestión de la barrera:

            SI gestion_values está vacío o no existe:
              → No crear tareas
              → CONTINÚA al siguiente entry

            SI NO:
              → Parsear CSV → lista de values seleccionados

              ↺ PARA CADA gestion_value EN lista:

                DB.Create(case_task {
                  category:        "Barrera",
                  type:            gestion_value,
                  description:     label del gestion_value  ← ver tabla de labels
                  assigned_user_id: actorId,
                  status:          "ToDo",
                  case_id:         caseId,
                  follow_up_id:    followUpId,
                  barrier_id:      newBarrierID
                })

              → FIN PARA CADA gestion_value
              → CONTINÚA al siguiente entry

  → FIN PARA CADA entry
```

---

## Tabla de opciones de gestión (Q23)

La pregunta `gestión de la barrera` es de tipo `multiple`. Cada value seleccionado genera una tarea. El `description` de la tarea se toma del label correspondiente:

| `type` (value guardado) | `description` (label legible) |
|---|---|
| `orientacion_llamada` | Orientación y enrutamiento - Llamada |
| `gestion_llamada` | Gestión administrativa - Llamada |
| `activacion_ruta_interinstitucional` | Activación de ruta interinstitucional |
| `articulacion_institucional` | Articulación institucional |
| `escalamiento_organismo_control` | Escalamiento a organismo de control |
| `alerta_barreras` | Alerta por barreras |

---

## Campos del `case_task` a poblar

| Campo | Valor |
|---|---|
| `category` | `"Barrera"` |
| `type` | value de la opción de gestión (ej. `"orientacion_llamada"`) |
| `description` | label legible de la opción (ej. `"Orientación y enrutamiento - Llamada"`) |
| `assigned_user_id` | `actorId` — el agente que completó el seguimiento |
| `status` | `"ToDo"` |
| `case_id` | `caseId` del seguimiento |
| `follow_up_id` | ID del `follow_up_v2` procesado |
| `barrier_id` | ID del `barrier_v2` recién creado en el paso 3.2 |

---

## Archivo a modificar

| Archivo | Cambio |
|---|---|
| `src/salvia/service/follow_up_submission_service.go` (o el archivo donde vive `processFollowUpSubmission`) | Agregar lógica de creación de tareas dentro del loop de barreras del PASO 3 |

> ⚠️ **GAP:** Confirmar el nombre exacto del archivo donde está implementado el PASO 3 de `processFollowUpSubmission` antes de implementar.

---

## Preguntas para validar antes de implementar

1. ¿El `assigned_user_id` de la tarea debe ser siempre el agente que completó el seguimiento (`actorId`)? ¿O debe asignarse a otro usuario según alguna regla?
2. ¿Se deben crear tareas aunque la barrera tenga `gestion_values` vacío, o solo cuando hay al menos una opción seleccionada?
3. ¿La `category` de la tarea debe ser `"Barrera"` u otro valor? ¿Hay un catálogo de categorías definido?
