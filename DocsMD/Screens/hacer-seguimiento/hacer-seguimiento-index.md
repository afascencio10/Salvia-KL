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
| [`changelogJun2026.md`](changelogJun2026.md) | Cambios realizados en junio 2026 |

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
Descripción:  Se ejecuta al montar el app Vue de #seguimiento-app. Lanza en
              paralelo loadFollowUp() y loadLocations(). loadFollowUp()
              llama al backend para obtener el seguimiento, victimInfo
              (con riskLevel), caseStatus y activeBarriers (barreras
              activas del caso con status != MANAGED). En la primera carga
              el backend fija fu.active_barrier_ids con esos IDs; en cargas
              posteriores los usa directamente. Guarda caseRiskLevel y
              formState.currentBarriers. Calcula canEdit: (1) caseStatus
              'cd' → false permanente; (2) REALIZADO + completed_at → true
              solo si ≤ 5 días; (3) cualquier otro → true. loadLocations()
              fetchea departamentos (→ formState.statesColombia) y ciudades
              (→ allCities local) en background sin loader, para alimentar
              los dropdowns de ubicación de la Sección 4.
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
              (3) Procesa barreras (Sección 4 repeater) → crea 1 BarrierV2
              por entry con todos los campos (sector, specific_barriers,
              institutions, ubicación, barreras estructurales, descripción,
              gestión). Cambio respecto a versión anterior: antes se creaba
              1 registro por opción CSV; ahora 1 registro por entry completa.
              (3b) Procesa Seguimiento a Barreras (Sección 3 repeater) →
              por cada entry: relaciona con barrera por posición en
              fu.active_barrier_ids, crea evento timeline "Seguimiento a
              Barrera" con resumen (persiste, respuesta institucional,
              actuaciones), y si Q6 cierra = true → actualiza
              barrier_v2.status = MANAGED.
              (4) Procesa derivaciones a equipos según qEquipos (multi-select).
              (5) Procesa medidas de emergencia → 1 EmergencyMeasure por medida.
              (6) Marca follow-up como REALIZADO y crea evento timeline.
              (7) Si qCierraCaso == "true" → CasoCierreService.CerrarCaso.
              (8) reasignarCaso: evalúa nivel y confirmaciones → actualiza
              risk_level, ReasignarCalendario, evento "Reasignación de Caso".
Requerido:    Sí
```

---

### E-05 — Cuando el formulario emite `answers-updated`

📄 [Ver flujo → flow-E05-cuando-se-actualizan-respuestas.md](Flujos/flow-E05-cuando-se-actualizan-respuestas.md)

```
Evento:       Cuando el formulario emite answers-updated
Tipo:         User Interaction
Descripción:  DinamicForm emite 'answers-updated' cada vez que el usuario
              guarda una sección (y también en carga inicial). La pantalla
              ejecuta onAnswersUpdated() que delega a tres métodos privados:
              (1) _updatePsicosocialState: evalúa criterio_obligatorio +
              puntaje ≥ 3 → escribe formState.psysocialRemisionState.
              (2) _updateReasignacionState: lee factores protectores, de
              riesgo y extremos junto con caseRiskLevel → escribe
              shouldReassignCase, reassingText, canReassignHigh y
              canReassignLow en formState.
              (3) _updateBarreraLocationOptions: para cada entry del repeater
              de Sección 4 (Identificación de Barreras), filtra ciudades de
              allCities según el departamento seleccionado (con cache por
              índice) y fetchea municipios por ciudad solo cuando cambia
              (con cache por índice). Escribe formState.newBarriers[i].cities
              y formState.newBarriers[i].towns que dinamic-form usa vía
              stateOptionsPath en Q13 y Q14.
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
