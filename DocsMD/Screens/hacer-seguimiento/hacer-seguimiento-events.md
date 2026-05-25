# `hacer-seguimiento` — Inventario de Eventos

Pantalla: `/salvia/hacer-seguimiento/:id`  
Archivo fuente: `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html`  
Entry point backend: `GET /salvia/hacer-seguimiento/:id` → `HacerSeguimientoFacade.go`

> **Nota de alcance:** Esta pantalla delega toda la interacción con el formulario al componente `dinamic-form`. Los 15 eventos internos de ese componente (navegación, inputs, validación, guardado por sección, repeaters) están documentados en `dinamic-form-events.md` y no se repiten aquí. Este inventario cubre únicamente los eventos propios de la pantalla anfitriona.

---

## Eventos identificados

---

### E-01 — Cuando carga la pantalla

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

```
Evento:       Cuando el backend procesa el seguimiento completado
Tipo:         Backend/Scheduled
Descripción:  Goroutine disparada por el form component (E-04 de dinamic-form)
              en el momento en que todas las secciones visibles quedan
              respondidas. No bloquea la respuesta HTTP. Ejecuta
              processFollowUpSubmission(submissionId, actorId), que:
              marca el follow-up como REALIZADO, registra el evento en el
              timeline, genera los intentos del calendario, procesa barreras,
              medidas de emergencia, apoyo psicosocial y estabilización
              económica. Adicionalmente, si el profesional marcó cierre de
              caso en la Sección 5, llama a CasoCierreService.CerrarCaso
              para actualizar victim_case_status a "cd" y registrar el
              evento de cierre en el timeline.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia de la pantalla está cubierta? → E-03
- [x] ¿Los eventos emitidos por componentes hijos que disparan lógica en la pantalla están cubiertos? → E-02
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → E-04
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga la pantalla | Lifecycle | No (solo lee) |
| E-02 | Cuando el formulario emite `form-completed` | User Interaction | No (estado local) |
| E-03 | Cuando presiona "Ver caso" | User Interaction | No (navegación) |
| E-04 | Cuando el backend procesa el seguimiento completado | Backend/Scheduled | Sí (follow-up, timeline, caso) |

**Total: 4 eventos — 1 Lifecycle, 2 User Interaction, 1 Backend/Scheduled**  
**E-04 es el único evento que escribe en el backend desde esta pantalla.**  
**El guardado sección a sección ocurre dentro de `dinamic-form` (E-04 de ese componente) y es transparente para esta pantalla.**
