# Flujo Psicosocial — Máquina de Estados y Selección de Formulario

Describe los escenarios posibles al registrar una sesión psicosocial, qué formulario se muestra en cada caso y cómo el backend actualiza el estado tras completarlo.

**Fuente:** Sheet _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial` + diagrama de casos.

---

## Arquitectura: 4 formularios independientes

Cada tipo de sesión tiene su propio formulario con secciones propias. El backend selecciona el form_id correcto leyendo el estado de `salvia.psychosocial_support`.

| Formulario | Form ID | Secciones |
|---|---|---|
| **Primer Contacto** | `FORM_ID_PRIMER_CONTACTO` | 1 — Primer contacto · 2 — Primera atención *(condicional: visible si "Continuar PA = Sí")* |
| **Primera Atención** | `FORM_ID_PRIMERA_ATENCION` | 1 — Contacto · 2 — Primera atención |
| **Seguimiento** | `FORM_ID_SEGUIMIENTO` | 1 — Contacto · 2 — Seguimiento |
| **Cierre** | `FORM_ID_CIERRE` | 1 — Contacto · 2 — Seguimiento · 3 — Cierre *(condicional: visible si "Cerrar remisión = Sí")* |

> Los Form IDs se generan con el seed y se almacenan como constantes en el backend Go.

---

## Variables de estado (`salvia.psychosocial_support`)

| Campo | Tipo | Valor inicial | Descripción |
|---|---|---|---|
| `ya_hizo_primer_contacto` | boolean | `false` | Se activa al completar el formulario de Primer Contacto |
| `ya_hizo_primera_atencion` | boolean | `false` | Se activa cuando la Primera Atención se completa (ya sea en el Form PC con S2 visible, o en el Form PA) |
| `session_count` | int | `0` | Sesiones completadas (Primera Atención + Seguimientos). Campo ya existente |
| `status` | varchar | `abierto` | `abierto` / `en_gestion` / `en_devolucion` / `cerrado`. Campo ya existente |

---

## Selección de formulario por escenario

### Escenario A — Primera llamada (Primer Contacto)
**Condición:** `ya_hizo_primer_contacto = false`

- Se carga **Form: Primer Contacto** (siempre, sea cual sea la respuesta a "Continuar Primera Atención")
- El formulario tiene 2 secciones: S1 (Primer contacto, siempre visible) + S2 (Primera Atención, oculta por defecto)
- **La llamada fue efectiva** no aparece: si el profesional está llenando el formulario, significa que la llamada ya fue respondida

**Pregunta gatillo al final de S1:** `Continuar Primera Atención` (boolean)

| Respuesta | Comportamiento en pantalla | Estado resultante |
|---|---|---|
| **Sí** | S2 (Primera Atención) se despliega. "Fecha próxima atención" desaparece de S1 y aparece al final de S2. El profesional completa ambas secciones en la misma sesión | `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = true`, `session_count += 1`, `status = en_gestion` |
| **No** | S2 permanece oculta. "Fecha próxima atención" sigue visible en S1 para agendar la próxima llamada | `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = false`, `status = en_gestion` |

> **Nota:** Ya no existe un escenario de "redirección" al Form de Primera Atención desde Primer Contacto. Todo ocurre dentro del mismo formulario.

---

### Escenario B — Segunda llamada (Primera Atención, llamada separada)
**Condición:** `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = false`

- Se carga **Form: Primera Atención** (aplica cuando el profesional respondió "Continuar = No" en la sesión anterior)
- El profesional llena la sección de Contacto (incluyendo "¿Es atención o solo contacto?")
- Si la sesión avanza, llena la sección "Primera Atención" con consentimiento, ajuste razonable, etc.

**Al completar:**

| Condición | Estado resultante |
|---|---|
| Consentimiento = Sí | `ya_hizo_primera_atencion = true`, `session_count += 1`, `status = en_gestion` |
| Consentimiento = No | `ya_hizo_primera_atencion` permanece `false`, `status = en_devolucion`. Se registra sesión con `session_type = CIERRE_NO_CONSENTIMIENTO` |
| "Solo contacto" | `ya_hizo_primera_atencion` permanece `false`. Se registra `session_type = CONTACTO_SIN_ATENCION` |

---

### Escenario C — Seguimiento regular
**Condición:** `ya_hizo_primera_atencion = true`, `session_count >= 1`, `session_count < 3`

- Se carga **Form: Seguimiento**
- Sección 1 ("Contacto Seguimiento") + Sección 2 ("Seguimiento")

**Al completar:** `session_count += 1`

---

### Escenario D — Seguimiento con posibilidad de cierre (3 sesiones en adelante)
**Condición:** `ya_hizo_primera_atencion = true`, `session_count >= 3`

- Se carga **Form: Cierre**
- S1 (Contacto) + S2 (Seguimiento) siempre visibles cuando hay atención
- S3 (Cierre) **solo visible** si el profesional responde "Cerrar remisión = Sí" en S2

**Al completar S2 con "Cerrar remisión = No":** `session_count += 1`, `status` permanece `en_gestion` (actúa como seguimiento regular)
**Al completar S3 (Cierre):** `session_count += 1`, `status = cerrado`

---

## Resumen de escenarios

| Escenario | `ya_hizo_pc` | `ya_hizo_pa` | `session_count` | Formulario cargado | Secciones visibles |
|---|---|---|---|---|---|
| A — Primera llamada, Continuar = No | false | false | 0 | Primer Contacto | Solo S1 — Primer contacto |
| A — Primera llamada, Continuar = Sí | false | false | 0 | Primer Contacto | S1 + S2 — Primera atención |
| B — Segunda llamada separada | true | false | 0 | Primera Atención | S1 Contacto + S2 Primera atención |
| C — Seguimiento regular | true | true | 1–2 | Seguimiento | S1 Contacto + S2 Seguimiento |
| D — Seguimiento sin cierre | true | true | ≥3 | Cierre | S1 Contacto + S2 Seguimiento (Cerrar = No) |
| D — Cierre definitivo | true | true | ≥3 | Cierre | S1 Contacto + S2 Seguimiento + S3 Cierre (Cerrar = Sí) |

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
- **Con** "Continuar Primera Atención" al final (gatillo de S2)
- **Con** "Fecha próxima atención" condicional: visible solo cuando "Continuar = No"

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
| Primer Contacto | Sección 2 — Primera Atención | `Continuar Primera Atención` | `true` |
| Primera Atención | Sección 2 — Primera Atención | `¿Es atención o solo contacto?` | `atencion` |
| Seguimiento | Sección 2 — Seguimiento | `¿Es atención o solo contacto?` | `atencion` |
| Cierre | Sección 2 — Seguimiento (Cierre) | `¿Es atención o solo contacto?` | `atencion` |
| Cierre | Sección 3 — Cierre | `Cerrar remisión` | `true` |

La pregunta "Fecha nueva" en PA/SEG/CIE solo es visible cuando se selecciona `solo_contacto`.
La pregunta "Fecha próxima atención" en PC-S1 solo es visible cuando "Continuar Primera Atención = No".

---

## Comportamiento de "Cerrar remisión" (Form 4 — Cierre, Sección 2)

Esta pregunta decide si el profesional está cerrando la remisión o simplemente registrando un seguimiento más dentro del formulario de Cierre.

| Respuesta | Sección 3 | Impacto en estado |
|---|---|---|
| **Sí** | S3 (Cierre) se muestra | Al guardar: `session_count += 1`, `status = cerrado` |
| **No** | S3 permanece oculta | Al guardar: `session_count += 1`, `status` sin cambio (sigue `en_gestion`) |

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
3. Determinar session_type según el formulario y respuestas clave:

     Form Primer Contacto:
       - "Continuar Primera Atención" = true  → session_type = PRIMER_CONTACTO_CON_ATENCION
       - "Continuar Primera Atención" = false → session_type = PRIMER_CONTACTO

     Form Primera Atención:
       - "¿Es atención o solo contacto?" = atencion      → session_type = PRIMERA_ATENCION
       - "¿Es atención o solo contacto?" = solo_contacto → session_type = CONTACTO_SIN_ATENCION

     Form Seguimiento:
       - "¿Es atención o solo contacto?" = atencion      → session_type = SEGUIMIENTO
       - "¿Es atención o solo contacto?" = solo_contacto → session_type = CONTACTO_SIN_ATENCION

     Form Cierre:
       - "¿Es atención o solo contacto?" = solo_contacto → session_type = CONTACTO_SIN_ATENCION
       - "¿Es atención o solo contacto?" = atencion Y "Cerrar remisión" = true  → session_type = CIERRE
       - "¿Es atención o solo contacto?" = atencion Y "Cerrar remisión" = false → session_type = SEGUIMIENTO

4. Actualizar team_contact.session_type, team_contact.is_completed = true, completed_at = NOW()
   Si solo_contacto: is_psico_session = false

5. Actualizar psychosocial_support según session_type:
     - PRIMER_CONTACTO:              ya_hizo_primer_contacto = true, status = en_gestion
     - PRIMER_CONTACTO_CON_ATENCION: ya_hizo_primer_contacto = true,
                                     ya_hizo_primera_atencion = true,
                                     session_count += 1,
                                     status = en_gestion
     - PRIMERA_ATENCION:             ya_hizo_primera_atencion = true, session_count += 1
     - SEGUIMIENTO:                  session_count += 1
     - CIERRE:                       session_count += 1, status = cerrado
     - CONTACTO_SIN_ATENCION:        scheduled_at = fecha_nueva del formulario (sin cambio en contadores)

6. INSERT salvia.case_timeline_event
```

---

## Decisión: "Seleccionar número de dupla" eliminada

Comentario del lider en el Sheet: _"No va → Dupla asignada en sistema"_.  
La dupla se gestiona desde `psychosocial_support.dupla_id` (asignada por el supervisor). No es una pregunta del formulario.

## Decisión: "¿La llamada fue efectiva?" en Primer Contacto eliminada

Comentario del lider: _"No va → esto es antes de contestar"_.  
El formulario de Primer Contacto se llena porque la llamada YA fue efectiva. Si no fue efectiva, el profesional no abre el formulario (registra el intento de contacto por otro medio, TBD).

## Decisión: Primera Atención dentro del Primer Contacto (Jul 2026)

El líder indicó que si el profesional responde "Continuar Primera Atención = Sí", toda la atención debe quedar registrada en el mismo formulario de Primer Contacto, no en un formulario aparte. Esto simplifica el flujo: el profesional no navega a otra pantalla y el registro queda atómico en un único `form_submission`. Se elimina la dependencia de `formState.skipContact`.
