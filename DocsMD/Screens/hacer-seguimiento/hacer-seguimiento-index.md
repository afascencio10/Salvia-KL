# `hacer-seguimiento` — Inventario de Eventos

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`interfaz.md`](interfaz.md) | Árbol de interfaz — estructura visual y condiciones de render |
| [`related-tables.md`](related-tables.md) | Tablas de BD usadas por esta pantalla |

### Changelogs

| Archivo | Descripción |
|---|---|
| [`changelogMay2025.md`](changelogMay2025.md) | Cambios realizados en mayo 2025 |
| [`changelogMay2026.md`](changelogMay2026.md) | Cambios realizados en mayo 2026 |

### Otros MDs

| Archivo | Descripción |
|---|---|
| [`form-seguimiento-data.md`](form-seguimiento-data.md) | Datos y estructura del formulario dinámico de seguimiento |

### Resumen de eventos

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | [Cuando carga la pantalla](Flujos/flow-E01-cuando-carga-pantalla.md) | Lifecycle | No (solo lee) |
| E-02 | Cuando el formulario emite `form-completed` | User Interaction | No (estado local) |
| E-03 | Cuando presiona "Ver caso" | User Interaction | No (navegación) |
| E-04 | [Cuando el backend procesa el seguimiento completado](Flujos/flow-E04-cuando-se-procesa-submission.md) | Backend/Scheduled | Sí (follow-up, timeline, caso) |
| E-05 | [Cuando el formulario emite `answers-updated`](Flujos/flow-E05-cuando-se-actualizan-respuestas.md) | User Interaction | No (estado local reactivo) |

---



Pantalla: `/salvia/hacer-seguimiento/:id`  
Archivo fuente: `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html`  
Entry point backend: `GET /salvia/hacer-seguimiento/:id` → `HacerSeguimientoFacade.go`

> **Nota de alcance:** Esta pantalla delega toda la interacción con el formulario al componente `dinamic-form`. Los 15 eventos internos de ese componente (navegación, inputs, validación, guardado por sección, repeaters) están documentados en `dinamic-form-events.md` y no se repiten aquí. Este inventario cubre únicamente los eventos propios de la pantalla anfitriona.

---

## Eventos identificados

---

### E-01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  Se ejecuta al montar el app Vue de #seguimiento-app. Llama al
              backend con el :id de la URL para obtener los datos del
              seguimiento (formId, submissionId, status, completed_at),
              los datos de la víctima (victimInfo) y el estado del caso
              (caseStatus). Calcula canEdit con la siguiente prioridad:
              (1) si caseStatus === 'cd' (caso cerrado) → canEdit = false
              permanente, sin importar ninguna otra condición;
              (2) si status === 'REALIZADO' y completed_at existe →
              canEdit = true solo si han pasado ≤ 5 días desde completed_at;
              (3) en cualquier otro caso → canEdit = true.
              Pasa formId, submissionId y canEdit como props a DinamicForm.
              Si falla la carga, setea loadError para mostrar el ErrorAlert.
Requerido:    Sí
```

---

### E-02 — Cuando el formulario emite `form-completed`

```
Evento:       Cuando el formulario emite form-completed
Tipo:         User Interaction
Descripción:  El componente DinamicForm emite 'form-completed' cuando el
              usuario guarda la última sección y todas las secciones visibles
              quedan respondidas. La pantalla captura este evento y setea
              formCompleted = true, lo que monta el CompletedOverlay sobre
              toda la UI con el mensaje de éxito y el botón "Ver caso".
Requerido:    Sí
```

---

### E-03 — Cuando presiona "Ver caso"

```
Evento:       Cuando presiona "Ver caso"
Tipo:         User Interaction
Descripción:  Se ejecuta desde el botón dentro de CompletedOverlay, visible
              únicamente cuando formCompleted = true. Llama a goToCase() que
              redirige al navegador a /salvia/casos/:caseId/detalle usando
              el caseId obtenido durante la carga inicial (E-01).
Requerido:    Sí
```

---

### E-04 — Cuando el backend procesa el seguimiento completado

📄 [Ver flujo → flow-E04-cuando-se-procesa-submission.md](Flujos/flow-E04-cuando-se-procesa-submission.md)

```
Evento:       Cuando el backend procesa el seguimiento completado
Tipo:         Backend/Scheduled
Descripción:  Goroutine disparada por dinamic-form cuando allAnswered == true.
              No bloquea la respuesta HTTP. Ejecuta
              processFollowUpSubmission(submissionId, actorId), que:
              (1) Si el follow-up ya es REALIZADO → crea evento "Seguimiento
              Editado" en el timeline y termina (idempotente).
              (2) Construye answerMap con todas las respuestas directas.
              (3) Procesa barreras del repeater → crea registros BarrierV2.
              (4) Procesa derivaciones a equipos según qEquipos (multi-select):
                - atencion_psico: valida criterio_obligatorio + ≥ 3 puntos;
                  se omite si medidas_emergencia también fue seleccionado
                  (regla de exclusión);
                - atencion_hombres: valida que criterio_hombres esté marcado
                  → crea MenTeamRemision;
                - discapacidad: crea un DiscapacidadRemision por cada
                  servicio seleccionado (apoyo_lsc, enfoque_discapacidad);
                - estabilizacion: valida ≥ 1 criterio seleccionado
                  → crea EconomicStabilization.
              (5) Procesa medidas de emergencia → crea un EmergencyMeasure
              por cada medida seleccionada en qMedidasEmergencia.
              (6) Marca el follow-up como REALIZADO y crea evento en timeline.
              (7) Si qCierraCaso == "true" → llama CasoCierreService.CerrarCaso.
Requerido:    Sí
```

---

### E-05 — Cuando el formulario emite `answers-updated`

📄 [Ver flujo → flow-E05-cuando-se-actualizan-respuestas.md](Flujos/flow-E05-cuando-se-actualizan-respuestas.md)

```
Evento:       Cuando el formulario emite answers-updated
Tipo:         User Interaction
Descripción:  DinamicForm emite 'answers-updated' cada vez que el usuario
              guarda una sección. La pantalla captura el payload (answers.directAnswers)
              y ejecuta onAnswersUpdated() para computar el estado de remisión
              al equipo psicosocial en tiempo real. Busca la respuesta a la
              pregunta de Criterios de remisión — Atención Psicosocial
              (qCriteriosPsico) y evalúa si cumple con la regla:
              debe incluir 'criterio_obligatorio' Y sumar ≥ 3 puntos según
              el puntaje de cada criterio. El resultado se escribe en
              formState.psysocialRemisionState ('Remisión: SI cumple' /
              'Remisión: NO cumple' / '') y se pasa como prop formState
              al componente dinamic-form, que lo usa para resolver el
              render_modification REPLACE del banner informativo de
              remisión psicosocial.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia de la pantalla está cubierta? → E-03
- [x] ¿Los eventos emitidos por componentes hijos que disparan lógica en la pantalla están cubiertos? → E-02, E-05
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → E-04
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

**Total: 5 eventos — 1 Lifecycle, 3 User Interaction, 1 Backend/Scheduled**  
**E-04 es el único evento que escribe en el backend desde esta pantalla.**  
**El guardado sección a sección ocurre dentro de `dinamic-form` (E-04 de ese componente) y es transparente para esta pantalla.**
