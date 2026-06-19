# Requerimiento — Crear tareas por gestión de barrera al completar seguimiento

> ⚠️ **MD temporal de planificación.** Pendiente de validación antes de implementar.

---

## ¿Qué pide el requerimiento?

Cuando se completa el formulario de seguimiento y el backend procesa la Sección 4 (Identificación de Barreras), por cada barrera identificada se deben crear registros según la gestión seleccionada (Q23):

- **Opciones que generan solo una `case_task`:**
  - `orientacion_llamada` — Orientación y enrutamiento - Llamada
  - `gestion_llamada` — Gestión administrativa - Llamada
  - `alerta_barreras` — Alerta por barreras

- **Opciones que generan un `entity_letter` + una `case_task` vinculada al oficio:**
  - `activacion_ruta_interinstitucional` — Activación de ruta interinstitucional
  - `articulacion_institucional` — Articulación institucional
  - `escalamiento_organismo_control` — Escalamiento a organismo de control

---

## Cambios en la interfaz

**Ninguno.** El formulario ya recolecta la pregunta de gestión (Q23, tipo `multiple`) dentro del repeater de Sección 4. No se requiere agregar ni modificar ningún campo visible para el usuario.

---

## Eventos afectados

Solo se modifica **E-04**.

| Evento | Tipo | Cambio |
|---|---|---|
| E-04 — Cuando el backend procesa el seguimiento completado | Backend/Scheduled | Se extiende el PASO 3 para crear tareas y oficios por cada opción de gestión seleccionada en cada barrera |

---

## Cambio en el flujo de E-04

### Flujo modificado — PASO 3 (extracto)

```
PASO 3 — Procesar barreras (Sección 4 repeater)

  ↺ PARA CADA entry del repeater de Sección 4:

    3.1 — Leer respuestas de la entry:
            sector         → Q1  (dropdown)
            gestion_values → Q23 (multiple, CSV)
            [resto de campos según lógica existente]

    3.2 — Crear registro barrier_v2:
            DB.Create(barrier_v2 { ... })
            → newBarrierID = barrier_v2.id  ← guardar para el paso siguiente

    3.3 — [NUEVO] Crear tareas y oficios por gestión de la barrera:

            SI gestion_values está vacío:
              → No crear nada
              → CONTINÚA al siguiente entry

            SI NO:
              → Parsear CSV → lista de gestion_values seleccionados

              ↺ PARA CADA gestion_value EN lista:

                SEGÚN gestion_value:

                  CASO "orientacion_llamada" | "gestion_llamada" | "alerta_barreras":
                    → Crear case_task (sin oficio)
                    DB.Create(case_task {
                      category:         "Barrera",
                      type:             gestion_value,
                      description:      label[gestion_value],
                      assigned_user_id: actorId,
                      status:           "ToDo",
                      case_id:          caseId,
                      follow_up_id:     followUpId,
                      barrier_id:       newBarrierID,
                      entity_letter_id: nil
                    })
                    → CONTINÚA

                  CASO "activacion_ruta_interinstitucional" | "articulacion_institucional" | "escalamiento_organismo_control":
                    → Crear entity_letter primero
                    DB.Create(entity_letter {
                      barrier_id: newBarrierID,
                      case_id:    caseId,
                      state:      "por_proyectar",
                      agent_id:   actorId
                    })
                    → newLetterID = entity_letter.id

                    → Crear case_task vinculada al oficio
                    DB.Create(case_task {
                      category:         "Barrera",
                      type:             gestion_value,
                      description:      label[gestion_value],
                      assigned_user_id: actorId,
                      status:           "ToDo",
                      case_id:          caseId,
                      follow_up_id:     followUpId,
                      barrier_id:       newBarrierID,
                      entity_letter_id: newLetterID
                    })
                    → CONTINÚA

              → FIN PARA CADA gestion_value
              → CONTINÚA al siguiente entry

  → FIN PARA CADA entry
```

---

## Tabla de opciones de gestión (Q23)

| `gestion_value` | `description` (label) | Genera oficio |
|---|---|---|
| `orientacion_llamada` | Orientación y enrutamiento - Llamada | No |
| `gestion_llamada` | Gestión administrativa - Llamada | No |
| `alerta_barreras` | Alerta por barreras | No |
| `activacion_ruta_interinstitucional` | Activación de ruta interinstitucional | **Sí** |
| `articulacion_institucional` | Articulación institucional | **Sí** |
| `escalamiento_organismo_control` | Escalamiento a organismo de control | **Sí** |

---

## Campos a poblar por registro

### `case_task` (todos los casos)

| Campo | Valor |
|---|---|
| `category` | `"Barrera"` |
| `type` | value de la opción (ej. `"orientacion_llamada"`) |
| `description` | label legible de la opción |
| `assigned_user_id` | `actorId` — el agente que completó el seguimiento |
| `status` | `"ToDo"` |
| `case_id` | `caseId` del seguimiento |
| `follow_up_id` | ID del `follow_up_v2` procesado |
| `barrier_id` | ID del `barrier_v2` recién creado |
| `entity_letter_id` | ID del `entity_letter` recién creado (solo si genera oficio, `nil` si no) |

### `entity_letter` (solo opciones que requieren oficio)

| Campo | Valor |
|---|---|
| `barrier_id` | ID del `barrier_v2` recién creado |
| `case_id` | `caseId` del seguimiento |
| `state` | `"por_proyectar"` |
| `agent_id` | `actorId` — el agente que completó el seguimiento |

---

## Archivos a modificar

| Archivo | Cambio |
|---|---|
| Archivo donde vive `processFollowUpSubmission` | Agregar lógica del paso 3.3 dentro del loop de barreras |

> ⚠️ **GAP:** Confirmar el nombre exacto del archivo donde está implementado `processFollowUpSubmission` antes de implementar.

---

## Preguntas para validar antes de implementar

1. ¿El `agent_id` del `entity_letter` y el `assigned_user_id` de la `case_task` deben ser siempre `actorId` (el agente que hizo el seguimiento)? ¿O hay algún caso donde se asigna a otro usuario?
2. ¿Si `gestion_values` está vacío no se crea nada, o igual se crea la barrera sin tareas?
