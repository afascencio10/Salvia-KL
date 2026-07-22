# Flujo Psicosocial — Máquina de Estados y Selección de Formulario

Describe los escenarios posibles al registrar una sesión psicosocial, qué formulario se muestra en cada caso y cómo el backend actualiza el estado tras completarlo.

**Fuente:** Sheet _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial` + diagrama de casos.

---

## Arquitectura: 4 formularios independientes

Cada tipo de sesión tiene su propio formulario con secciones propias. El backend selecciona el form_id correcto leyendo el estado de `salvia.psychosocial_support`.

| Formulario | Form ID | Secciones |
|---|---|---|
| **Primer Contacto** | `FORM_ID_PRIMER_CONTACTO` | 1 — Primer contacto · 2 — Seguimiento a Barreras · 3 — Identificación de Barreras *(condicional)* · 4 — Primera atención *(condicional: visible si "Continuar PA = Sí")* |
| **Primera Atención** | `FORM_ID_PRIMERA_ATENCION` | 1 — Contacto · 2 — Seguimiento a Barreras · 3 — Identificación de Barreras *(condicional)* · 4 — Primera atención |
| **Atención Psicosocial** *(antes "Seguimiento")* | `FORM_ID_SEGUIMIENTO` | 1 — Contacto · 2 — Seguimiento a Barreras · 3 — Identificación de Barreras *(condicional)* · 4 — Atención Psicosocial |
| **Cierre** | `FORM_ID_CIERRE` | 1 — Contacto · 2 — Seguimiento a Barreras · 3 — Identificación de Barreras *(condicional)* · 4 — Atención Psicosocial (Cierre) · 5 — Cierre *(condicional: visible si "Cerrar remisión = Sí")* |

> Los Form IDs se generan con el seed y se almacenan como constantes en el backend Go. La constante `FORM_ID_SEGUIMIENTO` conserva su nombre en el código aunque el formulario ahora se muestre como "Atención Psicosocial".
>
> **Jul 2026:** Se agregaron las secciones "Seguimiento a Barreras" e "Identificación de Barreras" a los 4 formularios, replicando la funcionalidad del Formulario de Seguimiento (`DocsMD/Screens/hacer-seguimiento/`). Ver detalle en [Comportamiento de registro de Barreras](#comportamiento-de-registro-de-barreras-jul-2026) más abajo.
>
> **Jul 2026 (ajuste 4):** El formulario "Seguimiento" se renombró a **"Atención Psicosocial"** (junto con sus secciones internas que contenían la palabra "Seguimiento"). Ver [Decisión: Rebranding "Seguimiento" → "Atención Psicosocial"](#decisión-rebranding-seguimiento--atención-psicosocial-y-limpieza-de-preguntas-jul-2026).

---

## Variables de estado (`salvia.psychosocial_support`)

| Campo | Tipo | Valor inicial | Descripción |
|---|---|---|---|
| `ya_hizo_primer_contacto` | boolean | `false` | Se activa al completar el formulario de Primer Contacto |
| `ya_hizo_primera_atencion` | boolean | `false` | Se activa cuando la Primera Atención se completa (ya sea en el Form PC con S4 visible, o en el Form PA) |
| `session_count` | int | `0` | Sesiones completadas (Primera Atención + Seguimientos). Campo ya existente |
| `status` | varchar | `abierto` | `abierto` / `en_gestion` / `en_devolucion` / `cerrado`. Campo ya existente |

---

## Selección de formulario por escenario

### Escenario A — Primera llamada (Primer Contacto)
**Condición:** `ya_hizo_primer_contacto = false`

- Se carga **Form: Primer Contacto** (siempre, sea cual sea la respuesta a "Continuar Primera Atención")
- El formulario tiene 4 secciones: S1 (Primer contacto, siempre visible) + S2 (Seguimiento a Barreras, siempre habilitada) + S3 (Identificación de Barreras, condicional) + S4 (Primera Atención, oculta por defecto)
- **La llamada fue efectiva** no aparece: si el profesional está llenando el formulario, significa que la llamada ya fue respondida

**Pregunta gatillo al final de S1:** `Continuar Primera Atención` (boolean)

| Respuesta | Comportamiento en pantalla | Estado resultante |
|---|---|---|
| **Sí** | S4 (Primera Atención) se despliega. "Fecha próxima atención" desaparece de S1 y aparece al final de S4. El profesional completa ambas secciones en la misma sesión. La pregunta "Desde la atención anterior se han identificado barreras institucionales" también aparece al final de S1 | `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = true`, `session_count += 1`, `status = en_gestion` |
| **No** | S4 permanece oculta. "Fecha próxima atención" sigue visible en S1 para agendar la próxima llamada. La pregunta de barreras no se muestra (no aplica: no hubo atención en esta sesión) | `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = false`, `status = en_gestion` |

> **Nota:** Ya no existe un escenario de "redirección" al Form de Primera Atención desde Primer Contacto. Todo ocurre dentro del mismo formulario.

---

### Escenario B — Segunda llamada (Primera Atención, llamada separada)
**Condición:** `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = false`

- Se carga **Form: Primera Atención** (aplica cuando el profesional respondió "Continuar = No" en la sesión anterior)
- El profesional llena la sección de Contacto (incluyendo "¿Es atención o solo contacto?" y, si aplica, "Desde la atención anterior se han identificado barreras institucionales")
- S2 (Seguimiento a Barreras) siempre está habilitada; S3 (Identificación de Barreras) se muestra solo si la pregunta de barreras = Sí
- Si la sesión avanza, llena la sección "Primera Atención" (S4) con consentimiento, conducta suicida, etc. (Jul 2026: ya no incluye "ajuste razonable" ni "intérprete de idiomas", preguntas eliminadas)

**Al completar:**

| Condición | Estado resultante |
|---|---|
| Consentimiento = Sí | `ya_hizo_primera_atencion = true`, `session_count += 1`, `status = en_gestion` |
| Consentimiento = No | `ya_hizo_primera_atencion` permanece `false`, `status = en_devolucion`. Se registra sesión con `session_type = CIERRE_NO_CONSENTIMIENTO` |
| "Solo contacto" | `ya_hizo_primera_atencion` permanece `false`. Se registra `session_type = CONTACTO_SIN_ATENCION` |

---

### Escenario C — Seguimiento regular
**Condición:** `ya_hizo_primera_atencion = true`, `session_count >= 1`, `session_count < 3`

- Se carga **Form: Atención Psicosocial** *(antes "Seguimiento")*
- Sección 1 ("Contacto Atención Psicosocial") + Sección 2 ("Seguimiento a Barreras") + Sección 3 ("Identificación de Barreras", condicional) + Sección 4 ("Atención Psicosocial")

**Al completar:** `session_count += 1`

---

### Escenario D — Seguimiento con posibilidad de cierre (3 sesiones en adelante)
**Condición:** `ya_hizo_primera_atencion = true`, `session_count >= 3`

- Se carga **Form: Cierre**
- S1 (Contacto) + S2 (Seguimiento a Barreras) + S4 (Seguimiento) siempre visibles cuando hay atención
- S3 (Identificación de Barreras) condicional a la pregunta de barreras en S1
- S5 (Cierre) **solo visible** si el profesional responde "Cerrar remisión = Sí" en S4

**Al completar S4 con "Cerrar remisión = No":** `session_count += 1`, `status` permanece `en_gestion` (actúa como seguimiento regular)
**Al completar S5 (Cierre):** `session_count += 1`, `status = cerrado`

---

## Resumen de escenarios

| Escenario | `ya_hizo_pc` | `ya_hizo_pa` | `session_count` | Formulario cargado | Secciones visibles |
|---|---|---|---|---|---|
| A — Primera llamada, Continuar = No | false | false | 0 | Primer Contacto | S1 — Primer contacto + S2 — Seguimiento a Barreras (+ S3 si aplica) |
| A — Primera llamada, Continuar = Sí | false | false | 0 | Primer Contacto | S1 + S2 — Seguimiento a Barreras + S3 — Identif. Barreras (si aplica) + S4 — Primera atención |
| B — Segunda llamada separada | true | false | 0 | Primera Atención | S1 Contacto + S2 Seg. Barreras + S3 Identif. Barreras (si aplica) + S4 Primera atención |
| C — Seguimiento regular | true | true | 1–2 | Atención Psicosocial *(antes "Seguimiento")* | S1 Contacto + S2 Seg. Barreras + S3 Identif. Barreras (si aplica) + S4 Atención Psicosocial |
| D — Seguimiento sin cierre | true | true | ≥3 | Cierre | S1 Contacto + S2 Seg. Barreras + S3 Identif. Barreras (si aplica) + S4 Atención Psicosocial (Cerrar = No) |
| D — Cierre definitivo | true | true | ≥3 | Cierre | S1 Contacto + S2 Seg. Barreras + S3 Identif. Barreras (si aplica) + S4 Atención Psicosocial + S5 Cierre (Cerrar = Sí) |

> "S3 — Identificación de Barreras" solo aparece si la pregunta "Desde la atención anterior se han identificado barreras institucionales" (última de S1, o gatillada por "Continuar Primera Atención" en Primer Contacto) se responde "Sí". "S2 — Seguimiento a Barreras" siempre está habilitada.

---

## Secciones de "Contacto" compartidas

Las secciones de contacto de **Primera Atención**, **Atención Psicosocial** *(antes "Seguimiento")* y **Cierre** contienen el mismo conjunto base de preguntas. La diferencia es que en Primer Contacto las preguntas de contacto son propias de ese formulario (sin la pregunta "¿La llamada fue efectiva?" porque el form se llena durante la llamada activa).

### Preguntas de contacto — Primera Atención / Atención Psicosocial / Cierre
- ¿La atención es individual o en dupla? *(en PA ya existía; agregada a Atención Psicosocial y Cierre en Jul 2026 — antes ausente en esos 2 formularios)*
- ¿La llamada fue efectiva? *(solo en PA, Atención Psicosocial y CIE — no en PC)*
- ¿Se encuentra en un lugar seguro?
- ¿Se encuentra en riesgo inminente?
- Describa las acciones ante riesgo inminente *(condicional a riesgo=Sí)*
- Hay nuevos hechos de violencia
- Descripción de los hechos *(condicional a nuevos_hechos=true)*
- Fecha *(condicional a nuevos_hechos=true)*
- ¿Es atención o solo contacto? *(Atención / Solo Contacto)*
- Observaciones del contacto
- Fecha nueva *(para reagendar si es solo contacto)*
- **Desde la atención anterior se han identificado barreras institucionales** *(nueva, Jul 2026 — última pregunta, condicional a "¿Es atención o solo contacto? = Atención")*

> Jul 2026: en Atención Psicosocial y Cierre, al agregarse "¿Individual o en dupla?" como primera pregunta, el resto se recorrió 1 posición ("¿Es atención o solo contacto?" pasó de Q8 a Q9, "Desde la atención anterior..." pasó de Q11 a Q12). En Primera Atención el orden no cambió (ya tenía esta pregunta como Q1).

### Preguntas de contacto — Primer Contacto (diferencias)
- **Sin** "¿La llamada fue efectiva?" — el form se llena porque la llamada fue respondida
- **Con** "¿Hay voluntariedad para la atención?" en lugar de "¿Es atención o solo contacto?"
- **Con** "Continuar Primera Atención" (gatillo de S4 — Primera Atención)
- **Con** "Fecha próxima atención" condicional: visible solo cuando "Continuar = No"
- **Con** "Desde la atención anterior se han identificado barreras institucionales" al final (nueva, Jul 2026 — condicional a "Continuar Primera Atención = Sí", gatillo de S3)

---

## Comportamiento de "¿Es atención o solo contacto?"

Esta pregunta aparece en la sección de Contacto de los formularios **Primera Atención**, **Atención Psicosocial** *(antes "Seguimiento")* y **Cierre**. Actúa como un condicional que determina si la sesión avanza o si solo se registra el intento de contacto.

### Camino: Solo Contacto

| Qué se muestra | Qué se oculta | Impacto en estado |
|---|---|---|
| "Observaciones del contacto" + "Fecha nueva" en la sección de Contacto | Toda(s) la(s) sección(es) posterior(es) al Contacto | **Sin cambio en variables de estado**: `ya_hizo_primera_atencion` permanece `false` (Form PA), `session_count` no incrementa (Form Atención Psicosocial / CIE) |

El `team_contact` se crea con `session_type = CONTACTO_SIN_ATENCION`, `is_psico_session = false`, `is_completed = true`. La fecha nueva del campo se copia a `psychosocial_support.scheduled_at` para agendar la próxima llamada.

### Camino: Atención

La sección de contenido (Sección 4 en PA y Atención Psicosocial, Secciones 4 y 5 en CIE) se muestra con normalidad. Se aplican todas las reglas de estado descritas en los escenarios A–D. Adicionalmente, se habilita la pregunta "Desde la atención anterior se han identificado barreras institucionales" al final de la sección de Contacto (ver [Comportamiento de registro de Barreras](#comportamiento-de-registro-de-barreras-jul-2026)).

### Implementación en DinamicForm

La visibilidad de las secciones de contenido se controla con `visibility_condition` de tipo `SECTION`:

| Formulario | Sección controlada | Trigger | Valor |
|---|---|---|---|
| Primer Contacto | Sección 4 — Primera Atención | `Continuar Primera Atención` | `true` |
| Primera Atención | Sección 4 — Primera Atención | `¿Es atención o solo contacto?` (Q8) | `atencion` |
| Atención Psicosocial | Sección 4 — Atención Psicosocial | `¿Es atención o solo contacto?` (Q9) | `atencion` |
| Cierre | Sección 4 — Atención Psicosocial (Cierre) | `¿Es atención o solo contacto?` (Q9) | `atencion` |
| Cierre | Sección 5 — Cierre | `Cerrar remisión` | `true` |
| Primer Contacto | Sección 3 — Identificación de Barreras | `Desde la atención anterior...` | `true` |
| Primera Atención | Sección 3 — Identificación de Barreras | `Desde la atención anterior...` (Q11) | `true` |
| Atención Psicosocial / Cierre | Sección 3 — Identificación de Barreras | `Desde la atención anterior...` (Q12) | `true` |

La pregunta "Fecha nueva" en PA/SEG/CIE solo es visible cuando se selecciona `solo_contacto`.
La pregunta "Fecha próxima atención" en PC-S1 solo es visible cuando "Continuar Primera Atención = No".

> **Sección 2 — Seguimiento a Barreras** no tiene `visibility_condition` basada en pregunta: permanece siempre habilitada en los 4 formularios. Su visibilidad real (mostrar u ocultar según si el caso tiene barreras activas) depende de `formState` externo — decisión pendiente, igual que en `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md`.

---

## Comportamiento de "Cerrar remisión" (Form 4 — Cierre, Sección 4)

Esta pregunta decide si el profesional está cerrando la remisión o simplemente registrando un seguimiento más dentro del formulario de Cierre.

| Respuesta | Sección 5 | Impacto en estado |
|---|---|---|
| **Sí** | S5 (Cierre) se muestra | Al guardar: `session_count += 1`, `status = cerrado` |
| **No** | S5 permanece oculta | Al guardar: `session_count += 1`, `status` sin cambio (sigue `en_gestion`) |

---

## Comportamiento de registro de Barreras (Jul 2026)

Requerimiento surgido en reunión: el Formulario Psicosocial debe poder registrar barreras institucionales, igual que el Formulario de Seguimiento (`DocsMD/Screens/hacer-seguimiento/`).

### Secciones nuevas

Se agregan **2 secciones nuevas a los 4 formularios**, ubicadas **justo después de la sección de Contacto**:

| Sección | Posición | Contenido | Fuente de estructura |
|---|---|---|---|
| **Seguimiento a Barreras** | S2 (todos los formularios) | Repeater de 7 preguntas por entrada, para dar seguimiento a barreras ya identificadas en el caso | `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md` |
| **Identificación de Barreras** | S3 (todos los formularios) | Repeater de 22 preguntas por entrada (mín. 1 repetición), para registrar nuevas barreras institucionales | `DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md` |

Ambas replican exactamente preguntas, opciones y funcionalidad (incluyendo dropdowns en cascada de ubicación vía `stateOptionsPath`) de la versión ya en producción del Formulario de Seguimiento — **no** la versión simplificada que existía en `seed_seguimiento.sql`.

### Pregunta gatillo en la sección de Contacto

Se agrega, en la **última posición** de la sección de Contacto de cada formulario, la pregunta booleana **"Desde la atención anterior se han identificado barreras institucionales"**:

| Formulario | Visible cuando | Motivo |
|---|---|---|
| Primer Contacto | `Continuar Primera Atención = Sí` | Este formulario no tiene "¿Es atención o solo contacto?"; la pregunta de barreras solo aplica si la atención avanza en la misma sesión |
| Primera Atención | `¿Es atención o solo contacto? = Atención` | — |
| Atención Psicosocial *(antes "Seguimiento")* | `¿Es atención o solo contacto? = Atención` | — |
| Cierre | `¿Es atención o solo contacto? = Atención` | — |

### Reglas de visibilidad resultantes

| Respuesta a la pregunta gatillo | Sección "Identificación de Barreras" (S3) | Sección "Seguimiento a Barreras" (S2) |
|---|---|---|
| **Sí** | Se muestra | Se muestra (sin cambio) |
| **No** / sin responder | Permanece oculta | Se muestra (sin cambio) |

**"Seguimiento a Barreras" no depende de esta pregunta**: permanece habilitada sin importar la respuesta. Igual que en el Formulario de Seguimiento, su visibilidad real (si mostrarla o no según si el caso tiene barreras activas) depende de `formState` externo — decisión pendiente de definir el `trigger_state_path` exacto.

### Requisito de implementación en frontend

La pantalla que aloje estos 4 formularios debe usar el componente `dinamic-form`, tal como lo hace `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` (incluyendo el manejo de `formState.currentBarriers` / `formState.newBarriers` y los dropdowns dependientes de ubicación).

---

## Preguntas de "Nuevos Hechos" (presentes en todos los formularios)

Aparecen en la sección de contacto de cada formulario (y en la sección de cierre):

| Pregunta | Tipo | Req |
|---|---|---|
| Hay nuevos hechos de violencia | boolean | ✅ |
| Descripción de los hechos | text | ❌ (visible si anterior = true) |
| Fecha | date | ❌ (visible si anterior = true) |

---

## Lógica de actualización del estado — evento E-02 (implementado, Jul 2026)

Al completar el formulario, `dinamic-form` hace `POST` de la última sección a `SaveSection`; el
backend (no el frontend) calcula si el formulario quedó completo y — de ser así — llama de forma
**síncrona** (no en goroutine) a `OnEndFormSubmission`, que despacha a
`processPsicosocialSessionSubmission` (`src/salvia/service/form_service.go`). Detalle completo,
IDs de pregunta reales y decisiones del líder en
`DocsMD/Screens/psicosocial-sesion/Flujos/flow-E02-cuando-se-guarda-formulario.md`.

```
1. Leer team_contact por form_submission_id (nuevo TeamContactRepository)
   SI team_contact.is_completed == true → idempotente: solo registra "Sesión Editada" y termina
2. Leer psychosocial_support por team_contact.psicosocial_id
3. Leer respuestas directas del form_submission → answerMap
4. Determinar session_type según el formulario y respuestas clave (resolvePsicosocialSessionType):

     Form Primer Contacto:
       - "Continuar Primera Atención" = true  Y consentimiento S4 = si → PRIMER_CONTACTO_CON_ATENCION
       - "Continuar Primera Atención" = true  Y consentimiento S4 = no → PRIMER_CONTACTO_SIN_CONSENTIMIENTO
       - "Continuar Primera Atención" = false                          → PRIMER_CONTACTO

     Form Primera Atención:
       - "¿Es atención o solo contacto?" = solo_contacto → CONTACTO_SIN_ATENCION
       - "¿Es atención o solo contacto?" = atencion Y consentimiento S4 = si → PRIMERA_ATENCION
       - "¿Es atención o solo contacto?" = atencion Y consentimiento S4 = no → CIERRE_NO_CONSENTIMIENTO

     Form Atención Psicosocial:
       - "¿Es atención o solo contacto?" = solo_contacto → CONTACTO_SIN_ATENCION
       - "¿Es atención o solo contacto?" = atencion      → ATENCION_PSICOSOCIAL

     Form Cierre:
       - "¿Es atención o solo contacto?" = solo_contacto        → CONTACTO_SIN_ATENCION
       - atencion Y "Cerrar remisión" = true                    → CIERRE
       - atencion Y "Cerrar remisión" = false                   → ATENCION_PSICOSOCIAL

5. Actualizar team_contact: session_type, is_completed = true, completed_at = NOW()
   Si CONTACTO_SIN_ATENCION: is_psico_session = false

6. Actualizar psychosocial_support según session_type:
     - PRIMER_CONTACTO / PRIMER_CONTACTO_SIN_CONSENTIMIENTO:
                                     ya_hizo_primer_contacto = true, status = en_gestion
     - PRIMER_CONTACTO_CON_ATENCION: ya_hizo_primer_contacto = true,
                                     ya_hizo_primera_atencion = true,
                                     session_count += 1, status = en_gestion
     - PRIMERA_ATENCION:             ya_hizo_primera_atencion = true, session_count += 1,
                                      status = en_gestion
     - ATENCION_PSICOSOCIAL:         session_count += 1
     - CIERRE:                       session_count += 1, status = cerrado
     - CIERRE_NO_CONSENTIMIENTO:     status = en_devolucion
     - CONTACTO_SIN_ATENCION:        scheduled_at = fecha_nueva del formulario (sin cambio en contadores)

7. INSERT salvia.case_timeline_event (categoría Psicosocial)
```

> **Pendiente (segunda iteración):** el procesamiento de "Identificación de Barreras" /
> "Seguimiento a Barreras" (creación de `barrier_v2`, tareas, oficios) queda fuera de esta
> primera versión — sus respuestas quedan guardadas en `salvia.answer`/`repeater_entry` pero no
> generan entidades derivadas todavía.

---

## Decisión: "Seleccionar número de dupla" eliminada

Comentario del lider en el Sheet: _"No va → Dupla asignada en sistema"_.  
La dupla se gestiona desde `psychosocial_support.dupla_id` (asignada por el supervisor). No es una pregunta del formulario.

## Decisión: "¿La llamada fue efectiva?" en Primer Contacto eliminada

Comentario del lider: _"No va → esto es antes de contestar"_.  
El formulario de Primer Contacto se llena porque la llamada YA fue efectiva. Si no fue efectiva, el profesional no abre el formulario (registra el intento de contacto por otro medio, TBD).

## Decisión: Primera Atención dentro del Primer Contacto (Jul 2026)

El líder indicó que si el profesional responde "Continuar Primera Atención = Sí", toda la atención debe quedar registrada en el mismo formulario de Primer Contacto, no en un formulario aparte. Esto simplifica el flujo: el profesional no navega a otra pantalla y el registro queda atómico en un único `form_submission`. Se elimina la dependencia de `formState.skipContact`.

## Decisión: Registro de Barreras en los 4 formularios psicosociales (Jul 2026)

Requerimiento surgido en reunión: el Formulario Psicosocial debe poder registrar barreras institucionales, igual que el Formulario de Seguimiento.

- Se agregan las secciones **"Seguimiento a Barreras"** (S2) e **"Identificación de Barreras"** (S3) a los 4 formularios, justo después de la sección de Contacto, replicando la estructura de `DocsMD/Screens/hacer-seguimiento/` (no la versión simplificada de `seed_seguimiento.sql`).
- Se agrega la pregunta gatillo "Desde la atención anterior se han identificado barreras institucionales" al final de la sección de Contacto de cada formulario. En Primer Contacto (que no tiene "¿Es atención o solo contacto?"), el gatillo equivalente es "Continuar Primera Atención = Sí".
- "Identificación de Barreras" es condicional a esa pregunta = Sí. "Seguimiento a Barreras" permanece siempre habilitada (su visibilidad real depende de `formState` externo, decisión pendiente).
- Todas las secciones posteriores a Contacto en cada formulario se renumeraron (+2, o +3 en Cierre) para dar espacio a las 2 secciones nuevas. Ver tabla actualizada en la sección [Arquitectura](#arquitectura-4-formularios-independientes) y el detalle completo en [Comportamiento de registro de Barreras](#comportamiento-de-registro-de-barreras-jul-2026).
- La pantalla debe usar el componente `dinamic-form`, igual que `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html`.
- Implementado en `src/cmd/seed/seed_psicosocial.sql` y documentado en `form-psicosocial-data.md`.

## Decisión: Rebranding "Seguimiento" → "Atención Psicosocial" y limpieza de preguntas (Jul 2026)

Ajuste 4 sobre el registro de Barreras, con 3 cambios adicionales:

1. **Eliminación de preguntas en "Primera Atención"** (Form 1 S4 y Form 2 S4): se quitan "¿Requiere algún ajuste razonable...?" y "¿Requiere intérprete de idiomas...?" (antes Q2/Q3). "Confirmación consentimiento persona de apoyo" (antes Q5, ahora Q3) dependía de la pregunta de intérprete; al eliminarse esa dependencia, pasa a ser **siempre visible**. 14 → 12 preguntas en ambas secciones.
2. **Nueva pregunta "¿La atención es individual o en dupla?"** en la sección de Contacto de **Atención Psicosocial** y **Cierre** (antes ausente en esos 2 formularios; Primer Contacto y Primera Atención ya la tenían). Se agrega como primera pregunta; el resto se recorre 1 posición (11 → 12 preguntas).
3. **Renombrado del Formulario de Seguimiento a "Atención Psicosocial"**, junto con sus secciones internas que contenían la palabra "Seguimiento" ("Contacto Seguimiento" → "Contacto Atención Psicosocial", "Seguimiento" → "Atención Psicosocial"). En el Formulario de Cierre (que conserva su nombre) se aplica el mismo cambio a su sección "Seguimiento (en Cierre)" → "Atención Psicosocial (en Cierre)". La sección **"Seguimiento a Barreras" NO se renombra** en ningún formulario — es un nombre fijo, igual en los 4 formularios. Las constantes de código (`FORM_ID_SEGUIMIENTO`, `session_type = SEGUIMIENTO`) se mantienen sin cambios; solo cambia el nombre mostrado al usuario.

Implementado en `src/cmd/seed/seed_psicosocial.sql` y documentado en `form-psicosocial-data.md`.
