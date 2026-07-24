# `registro-caso-v2` — Inventario de Eventos

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`registro-caso-v2-interface.md`](registro-caso-v2-interface.md) | Árbol de interfaz — estructura visual y condiciones de render |
| [`form-registro-caso-v2-data.md`](form-registro-caso-v2-data.md) | Estructura completa del formulario dinámico en BD (9 secciones, 128 preguntas) |

### Otros MDs

Ninguno todavía — pantalla nueva, sin changelog previo.

### Resumen de eventos

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | [Cuando carga la pantalla](Flujos/flow-E01-cuando-carga-pantalla.md) | Lifecycle | No (solo lee) |
| E-02 | [Cuando el formulario emite `answers-updated`](Flujos/flow-E02-cuando-se-actualizan-respuestas.md) | User Interaction | No (estado local reactivo) |
| E-03 | [Cuando se guarda la primera sección](Flujos/flow-E03-cuando-se-guarda-primera-seccion.md) | Backend/Scheduled | Sí (crea victim_case en Borrador) |
| E-04 | [Cuando se completa el formulario](Flujos/flow-E04-cuando-se-completa-formulario.md) | Backend/Scheduled | Sí (usuario, activación, calendario, ruta) |
| E-05 | [Cuando el formulario emite `form-completed`](Flujos/flow-E05-cuando-emite-form-completed.md) | User Interaction | No (GET de lectura) |
| E-06 | [Cuando presiona "Finalizar"](Flujos/flow-E06-cuando-presiona-finalizar.md) | User Interaction | No (navegación) |

> **Nota de alcance:** esta pantalla delega toda la interacción con el formulario a `dinamic-form`. Los 16 eventos internos de ese componente (navegación, inputs, validación, guardado por sección, repeaters, guardar borrador) están documentados en `dinamic-form-events.md` y no se repiten aquí. Este inventario cubre únicamente los eventos propios de la pantalla anfitriona — mismo criterio que `hacer-seguimiento-index.md`.

> **Nota de arquitectura (E-03/E-04):** a diferencia de Seguimiento (donde el caso ya existe antes de abrir el form), aquí el `victim_case` se crea progresivamente: **E-03** lo crea en estado `Borrador` (`bo`) al guardar la primera sección; **E-04** lo completa y activa (`ra`) al guardar la última. Ambos son disparados desde el mismo endpoint genérico `POST /api/v1/forms/saveSection`, con un dispatch SEGÚN `formId` — E-03 en el punto donde se crea un `form_submission` nuevo, E-04 en el punto donde `allAnswered == true` (mismo dispatch que ya usa `processFollowUpSubmission` hoy).

---

## Eventos identificados

---

### E-01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  Se ejecuta al montar el app Vue de la pantalla anfitriona.
              Primero resuelve el gate de autorización (fuera de dinamic-form).
              Tras autorizar, monta DinamicForm sin submissionId (formulario
              nuevo) o con uno existente si viene de un contacto previo.
              En paralelo carga el catálogo de enums (victim_case_form2_enums,
              agrupado por categoría → formState.enums.*) y departamentos
              (→ formState.departments), sin bloquear el render del form.
Requerido:    Sí
```

---

### E-02 — Cuando el formulario emite `answers-updated`

📄 [Ver flujo → flow-E02-cuando-se-actualizan-respuestas.md](Flujos/flow-E02-cuando-se-actualizan-respuestas.md)

```
Evento:       Cuando el formulario emite answers-updated
Tipo:         User Interaction
Descripción:  DinamicForm emite 'answers-updated' en cada guardado de sección
              (y en la carga inicial). La pantalla recalcula en formState:
              (1) wasPartner / partnerKnown / tamizajeAggressorCheck — a partir
              de relationshipWithPresumedAggressor y proximityPrincipalAggressor;
              (2) riskScore / riskLevel / hasRisk / riskBadgeText — misma
              fórmula de tamizaje que hoy, solo que recalculada en el host;
              (3) hasWorkplaceScope / hasSelectedViolenceTypes / violenceSubtypes
              — flags y opciones dinámicas para Sección 3;
              (4) las 3 cadenas geográficas (residencia, hechos, atención) —
              mismo patrón de cache por índice que Barreras en Seguimiento,
              generalizado a 3 instancias nombradas en vez de un repeater.
Requerido:    Sí
```

---

### E-03 — Cuando se guarda la primera sección

📄 [Ver flujo → flow-E03-cuando-se-guarda-primera-seccion.md](Flujos/flow-E03-cuando-se-guarda-primera-seccion.md)

```
Evento:       Cuando se guarda la primera sección
Tipo:         Backend/Scheduled
Descripción:  Hook nuevo dentro de POST /api/v1/forms/saveSection, en el
              punto donde se crea un form_submission nuevo (formSubmissionId
              venía vacío). Para el formId de Registro de Caso, crea de forma
              SÍNCRONA un victim_case + victim_case_form2 en estado
              'bo' (Borrador) con los campos ya disponibles de la Sección 1
              (nombres, tipo/número de documento — los mismos exigidos hoy
              como obligatorios por SetVictimCase), y vincula el
              form_submission al caso nuevo. Es idempotente — si el
              submission ya tiene un caso vinculado, no crea uno segundo.
Requerido:    Sí
```

---

### E-04 — Cuando se completa el formulario

📄 [Ver flujo → flow-E04-cuando-se-completa-formulario.md](Flujos/flow-E04-cuando-se-completa-formulario.md)

```
Evento:       Cuando se completa el formulario (todas las secciones respondidas)
Tipo:         Backend/Scheduled
Descripción:  Goroutine disparada por el mismo mecanismo que
              processFollowUpSubmission (OnEndFormSubmission, SEGÚN formID),
              cuando allAnswered == true tras el saveSection de la última
              sección. El victim_case YA EXISTE (creado en Borrador por E-03)
              — este evento lo completa: proyecta las 128 respuestas sobre
              victim_case_form2 (UPDATE completo), calcula riskScore/riskLevel,
              genera el usuario y contraseña de la víctima (bcrypt), activa
              el caso (status 'bo' → 'ra'), genera el calendario de
              seguimientos (FollowUpSvc.GenerateOrRecalculate), asigna
              equipo/agente y registra el evento en el timeline. Guarda el
              resultado (credenciales + caseId) para que el frontend lo
              recupere en E-05.
Requerido:    Sí
```

---

### E-05 — Cuando el formulario emite `form-completed`

📄 [Ver flujo → flow-E05-cuando-emite-form-completed.md](Flujos/flow-E05-cuando-emite-form-completed.md)

```
Evento:       Cuando el formulario emite form-completed
Tipo:         User Interaction
Descripción:  DinamicForm emite 'form-completed' usando el mismo criterio
              (isAnswered == true en todas las secciones visibles) que hoy
              dispara la redirección en hacer-seguimiento. La pantalla
              muestra el CompletedOverlay y hace un GET puntual al backend
              para obtener el resultado de E-04 (credenciales + caseId) —
              el endpoint espera brevemente al goroutine antes de responder,
              sin necesidad de que el cliente reintente en bucle.
Requerido:    Sí
```

---

### E-06 — Cuando presiona "Finalizar"

📄 [Ver flujo → flow-E06-cuando-presiona-finalizar.md](Flujos/flow-E06-cuando-presiona-finalizar.md)

```
Evento:       Cuando presiona "Finalizar"
Tipo:         User Interaction
Descripción:  Botón dentro del CompletedOverlay. Redirige a
              /salvia/casos/:caseId/detalle usando el caseId obtenido en E-05.
              Análogo a "Ver caso" en hacer-seguimiento.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia de la pantalla está cubierta? → E-06
- [x] ¿Los eventos emitidos por componentes hijos que disparan lógica en la pantalla están cubiertos? → E-02, E-05
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → E-03 (síncrono, dentro del request), E-04 (goroutine)
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

**Total: 6 eventos — 1 Lifecycle, 4 User Interaction, 2 Backend/Scheduled (de los cuales E-03 y E-04 escriben en el backend).**
**E-03 y E-04 son los eventos nuevos de esta migración (no existen en Seguimiento) — implementan la creación progresiva del caso (Borrador → Activo).**
**Todos los eventos tienen flujo documentado (módulo nuevo → modo Metódico, sin excepciones).**
