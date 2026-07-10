# Flujo Psicosocial — Máquina de Estados y Selección de Formulario

Describe los escenarios posibles al registrar una sesión psicosocial, qué formulario se muestra en cada caso y cómo el backend actualiza el estado tras completarlo.

**Fuente:** Sheet _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial` + diagrama de casos.

---

## Arquitectura: 4 formularios independientes

Cada tipo de sesión tiene su propio formulario con secciones propias. El backend selecciona el form_id correcto leyendo el estado de `salvia.psychosocial_support`.

| Formulario | Form ID | Secciones |
|---|---|---|
| **Primer Contacto** | `FORM_ID_PRIMER_CONTACTO` | 1 — Primer contacto |
| **Primera Atención** | `FORM_ID_PRIMERA_ATENCION` | 1 — Contacto · 2 — Primera atención |
| **Seguimiento** | `FORM_ID_SEGUIMIENTO` | 1 — Contacto · 2 — Seguimiento |
| **Cierre** | `FORM_ID_CIERRE` | 1 — Contacto · 2 — Seguimiento · 3 — Cierre |

> Los Form IDs se generan con el seed y se almacenan como constantes en el backend Go.

---

## Variables de estado (`salvia.psychosocial_support`)

| Campo | Tipo | Valor inicial | Descripción |
|---|---|---|---|
| `ya_hizo_primer_contacto` | boolean | `false` | Se activa al completar el formulario de Primer Contacto |
| `ya_hizo_primera_atencion` | boolean | `false` | Se activa cuando la Primera Atención se completa con consentimiento = Sí |
| `session_count` | int | `0` | Sesiones completadas (Primera Atención + Seguimientos). Campo ya existente |
| `status` | varchar | `abierto` | `abierto` / `en_gestion` / `en_devolucion` / `cerrado`. Campo ya existente |

---

## Selección de formulario por escenario

### Escenario A — Primera llamada (Primer Contacto)
**Condición:** `ya_hizo_primer_contacto = false`

- Se carga **Form: Primer Contacto**
- El formulario tiene 1 sola sección con preguntas de contacto + preguntas de hechos + pregunta clave "Continuar Primera Atención"
- **La llamada fue efectiva** no aparece: si el profesional está llenando el formulario, significa que la llamada ya fue respondida

**Pregunta gatillo al final:** `Continuar Primera Atención` (boolean)

| Respuesta | Acción siguiente |
|---|---|
| **Sí** | Al guardar, el backend actualiza `ya_hizo_primer_contacto = true`. La pantalla redirige al formulario de **Primera Atención** omitiendo la sección "Contacto" (carga la sección "Primera atención" directamente via `formState.skip_contact = true`) |
| **No** | Al guardar, el backend actualiza `ya_hizo_primer_contacto = true`. Se agenda una nueva llamada. El próximo escenario será B |

---

### Escenario B — Segunda llamada (Primera Atención, llamada separada)
**Condición:** `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = false`

- Se carga **Form: Primera Atención** con la sección "Contacto Primera Atención" visible
- El profesional llena la sección de Contacto (incluyendo "¿Es atención o solo contacto?")
- Si la sesión avanza, llena la sección "Primera Atención" con consentimiento, ajuste razonable, etc.

**Al completar:**

| Condición | Estado resultante |
|---|---|
| Consentimiento = Sí | `ya_hizo_primera_atencion = true`, `session_count += 1`, `status = en_gestion` |
| Consentimiento = No | `ya_hizo_primera_atencion` permanece `false`, `status = en_devolucion` o sin cambio. Se registra sesión con `session_type = CIERRE_NO_CONSENTIMIENTO` |
| Persona de apoyo no consiente | Requiere reagendamiento. Estado sin cambio |

---

### Escenario C — Seguimiento regular
**Condición:** `ya_hizo_primera_atencion = true`, `session_count >= 1`, `session_count < 3`

- Se carga **Form: Seguimiento**
- Sección 1 ("Contacto Seguimiento") + Sección 2 ("Seguimientos")

**Al completar:** `session_count += 1`

---

### Escenario D — Seguimiento con posibilidad de cierre (3 sesiones en adelante)
**Condición:** `ya_hizo_primera_atencion = true`, `session_count >= 3`

- Se carga **Form: Cierre**
- Sección 1 ("Contacto") + Sección 2 ("Seguimientos") + Sección 3 ("Cierre")
- El profesional puede completar solo las secciones 1 y 2 (seguimiento normal) o también la 3 (cierre)

**Al completar Sección 3 (Cierre):** `session_count += 1`, `status = cerrado`  
**Al completar solo Secciones 1-2:** `session_count += 1`, `status = en_gestion`

---

## Resumen de escenarios

| Escenario | `ya_hizo_pc` | `ya_hizo_pa` | `session_count` | Formulario cargado | Secciones visibles |
|---|---|---|---|---|---|
| A — Primera llamada | false | false | 0 | Primer Contacto | Primer contacto |
| A (variante) → Continuar=Sí | false | false | 0 | Primera Atención | Solo "Primera atención" (sin contacto) |
| B — Segunda llamada | true | false | 0 | Primera Atención | Contacto + Primera atención |
| C — Seguimiento regular | true | true | 1–2 | Seguimiento | Contacto + Seguimiento |
| D — Seguimiento con cierre | true | true | ≥3 | Cierre | Contacto + Seguimiento + Cierre |

---

## Secciones de "Contacto" compartidas

Las secciones de contacto de **Primera Atención**, **Seguimiento** y **Cierre** contienen el mismo conjunto base de preguntas. La diferencia es que en Primer Contacto las preguntas de contacto son propias de ese formulario (sin la pregunta "¿La llamada fue efectiva?" porque el form se llena durante la llamada activa).

### Preguntas de contacto — Primera Atención / Seguimiento / Cierre
- ¿La llamada fue efectiva? *(solo en PA, SEG y CIE — no en PC)*
- ¿Se encuentra en un lugar seguro?
- ¿Se encuentra en riesgo inminente?
- Describa las acciones ante riesgo inminente *(condicional a riesgo=Sí)*
- Hay nuevos hechos de violencia
- Descripción de los hechos *(condicional a nuevos_hechos=true)*
- Fecha *(condicional a nuevos_hechos=true)*
- ¿Es atención o solo contacto? *(Atención / Solo Contacto)*
- Observaciones del contacto
- Fecha nueva *(para reagendar si es solo contacto)*

### Preguntas de contacto — Primer Contacto (diferencias)
- **Sin** "¿La llamada fue efectiva?" — el form se llena porque la llamada fue respondida
- **Con** "¿Hay voluntariedad para la atención?" en lugar de "¿Es atención o solo contacto?"
- **Con** "Continuar Primera Atención" al final

---

## Comportamiento de "¿Es atención o solo contacto?"

Esta pregunta aparece en la sección de Contacto de los formularios **Primera Atención**, **Seguimiento** y **Cierre**. Actúa como un condicional que determina si la sesión avanza o si solo se registra el intento de contacto.

### Camino: Solo Contacto

| Qué se muestra | Qué se oculta | Impacto en estado |
|---|---|---|
| "Observaciones del contacto" + "Fecha nueva" en la sección de Contacto | Toda(s) la(s) sección(es) posterior(es) al Contacto | **Sin cambio en variables de estado**: `ya_hizo_primera_atencion` permanece `false` (Form PA), `session_count` no incrementa (Form SEG / CIE) |

El `team_contact` se crea con `session_type = CONTACTO_SIN_ATENCION`, `is_psico_session = false`, `is_completed = true`. La fecha nueva del campo se copia a `psychosocial_support.scheduled_at` para agendar la próxima llamada.

### Camino: Atención

La sección de contenido (Sección 2 en PA y SEG, Secciones 2 y 3 en CIE) se muestra con normalidad. Se aplican todas las reglas de estado descritas en los escenarios A–D.

### Implementación en DinamicForm

La visibilidad de las secciones de contenido se controla con `visibility_condition` de tipo `SECTION`:

| Formulario | Sección controlada | Trigger | Valor |
|---|---|---|---|
| Primera Atención | Sección 2 — Primera Atención | `¿Es atención o solo contacto?` | `atencion` |
| Seguimiento | Sección 2 — Seguimiento | `¿Es atención o solo contacto?` | `atencion` |
| Cierre | Sección 2 — Seguimiento (Cierre) | `¿Es atención o solo contacto?` | `atencion` |
| Cierre | Sección 3 — Cierre | `¿Es atención o solo contacto?` | `atencion` |

La pregunta "Fecha nueva" solo es visible cuando se selecciona `solo_contacto`, ya que es el campo para reagendar la sesión.

---

## Preguntas de "Nuevos Hechos" (presentes en todos los formularios)

Aparecen en la sección de contacto de cada formulario (y en la sección de cierre):

| Pregunta | Tipo | Req |
|---|---|---|
| Hay nuevos hechos de violencia | boolean | ✅ |
| Descripción de los hechos | text | ❌ (visible si anterior = true) |
| Fecha | date | ❌ (visible si anterior = true) |

---

## Lógica de actualización del estado (goroutine `processPsicoSession`)

Al completar el formulario (`form-completed`), el backend ejecuta en background:

```
1. Leer psychosocial_support por psicosocial_id del team_contact
2. Leer respuestas del form_submission
3. Determinar session_type según el formulario y respuesta de "¿Es atención o solo contacto?":
     - Form Primer Contacto            → session_type = PRIMER_CONTACTO
     - Form Primera Atención + atencion  → session_type = PRIMERA_ATENCION
     - Form Primera Atención + solo_contacto → session_type = CONTACTO_SIN_ATENCION
     - Form Seguimiento + atencion     → session_type = SEGUIMIENTO
     - Form Seguimiento + solo_contacto → session_type = CONTACTO_SIN_ATENCION
     - Form Cierre (S3 completada)     → session_type = CIERRE
     - Form Cierre (S3 no completada)  → session_type = SEGUIMIENTO
     - Form Cierre + solo_contacto     → session_type = CONTACTO_SIN_ATENCION
4. Actualizar team_contact.session_type, team_contact.is_completed = true, completed_at = NOW()
   Si solo_contacto: is_psico_session = false
5. Actualizar psychosocial_support según session_type:
     - PRIMER_CONTACTO:       ya_hizo_primer_contacto = true, status = en_gestion
     - PRIMERA_ATENCION:      ya_hizo_primera_atencion = true, session_count += 1
     - SEGUIMIENTO:           session_count += 1
     - CIERRE:                session_count += 1, status = cerrado
     - CONTACTO_SIN_ATENCION: scheduled_at = fecha_nueva del formulario (sin cambio en contadores)
6. INSERT salvia.case_timeline_event
```

---

## Decisión: "Seleccionar número de dupla" eliminada

Comentario del lider en el Sheet: _"No va → Dupla asignada en sistema"_.  
La dupla se gestiona desde `psychosocial_support.dupla_id` (asignada por el supervisor). No es una pregunta del formulario.

## Decisión: "¿La llamada fue efectiva?" en Primer Contacto eliminada

Comentario del lider: _"No va → esto es antes de contestar"_.  
El formulario de Primer Contacto se llena porque la llamada YA fue efectiva. Si no fue efectiva, el profesional no abre el formulario (registra el intento de contacto por otro medio, TBD).
