# Estructura de los 4 Formularios Psicosociales

Fuente: Sheet _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial`.

**Decisiones del líder (columna "Comentarios Felipe"):**
- `"Seleccionar número de dupla"` → **eliminada de todos los formularios** ("No va → Dupla asignada en sistema")
- `"¿La llamada fue efectiva?"` en Primer Contacto → **eliminada** ("No va → esto es antes de contestar")
- `"¿La atención es individual o en dupla?"` → **se mantiene** ("Aca contesta")

**Decisión del líder — ajuste 2 (Jul 2026):**
- Si "Continuar Primera Atención = Sí", **no se redirige** al formulario de Primera Atención. En su lugar, se muestra una nueva sección **"Primera Atención"** dentro del mismo formulario de Primer Contacto (originalmente S2; pasó a ser **S4** tras el ajuste 3 de Barreras, ver abajo).
- La pregunta "Fecha próxima atención" de S1 desaparece cuando "Continuar Primera Atención = Sí" (se mueve a la sección de Primera Atención).
- En el formulario de Cierre, se agrega la pregunta "Cerrar remisión" en la sección de Seguimiento (originalmente S2, pasó a ser **S4** tras el ajuste 3). La sección de Cierre (originalmente S3, ahora **S5**) solo se muestra si "Cerrar remisión = Sí".

**Decisión del líder — ajuste 3, registro de Barreras (Jul 2026):**
- Requerimiento surgido en reunión: el Formulario Psicosocial debe poder registrar barreras institucionales, igual que el Formulario de Seguimiento.
- Se agregan **2 secciones nuevas a los 4 formularios**, ubicadas **justo después de la sección de Contacto**: **"Seguimiento a Barreras"** e **"Identificación de Barreras"**. Ambas replican preguntas y funcionalidad de `DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md` y `form-seguimiento-barreras-repeater.md` (misma estructura ya en producción, no la versión simplificada de `seed_seguimiento.sql`).
- En la sección de Contacto de cada formulario se agrega, en la **última posición**, la pregunta **"Desde la atención anterior se han identificado barreras institucionales"** (boolean):
  - En **Primera Atención / Seguimiento / Cierre** → visible cuando "¿Es atención o solo contacto?" = `atencion`.
  - En **Primer Contacto** (que no tiene esa pregunta) → visible cuando "Continuar Primera Atención" = `true` (equivalente: solo aplica si la atención avanza en esta misma sesión).
- Respuesta de esa pregunta = Sí → habilita **"Identificación de Barreras"**. Respuesta = No/vacío → esa sección permanece oculta.
- **"Seguimiento a Barreras" NO depende de esta pregunta**: permanece habilitada sin importar la respuesta. Igual que en el Formulario de Seguimiento, su visibilidad real depende de si el caso tiene barreras activas (formState externo — decisión pendiente de definir el `trigger_state_path` exacto, igual que en `form-seguimiento-barreras-repeater.md`).
- La pantalla que aloja estos formularios debe usar el componente `dinamic-form`, igual que `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html`.

**Decisión del líder — ajuste 4, limpieza de preguntas y rebranding (Jul 2026):**
- Se **eliminan** de la Sección "Primera Atención" (Form 1 S4 y Form 2 S4) las preguntas **"¿Requiere algún ajuste razonable en el marco de la atención?"** y **"¿Requiere intérprete de idiomas y/o traducción?"** (antes Q2/Q3). La pregunta "Confirmación consentimiento persona de apoyo" (antes Q5, ahora Q3) dependía de la pregunta de intérprete eliminada; al quitarse esa dependencia, **pasa a ser siempre visible**. El resto de preguntas de la sección se renumeran (14 → 12).
- Se agrega la pregunta **"¿La atención es individual o en dupla?"** a la sección de Contacto de los formularios **Seguimiento** y **Cierre** (antes ausente en estos 2 formularios, ya existía en Primer Contacto y Primera Atención). Se ubica como **primera pregunta** de la sección, igual que en los otros 2 formularios. El resto de preguntas se recorre 1 posición (11 → 12; "¿Es atención o solo contacto?" pasa de Q8 a Q9; "Desde la atención anterior..." pasa de Q11 a Q12).
- El **Formulario de Seguimiento se renombra a "Atención Psicosocial"**, junto con sus secciones internas que contenían la palabra "Seguimiento": "Contacto Seguimiento" → "Contacto Atención Psicosocial", y la sección de contenido "Seguimiento" → "Atención Psicosocial". En el **Formulario de Cierre** (que mantiene su propio nombre) se aplica el mismo cambio de palabra a su sección "Seguimiento (en Cierre)" → "Atención Psicosocial (en Cierre)". La sección **"Seguimiento a Barreras" NO se renombra** en ningún formulario: mantiene su nombre fijo en los 4 formularios.
- Implementado en `src/cmd/seed/seed_psicosocial.sql`.

---

## Form 1: Primer Contacto

**Secciones:** 4 (S3 y S4 ocultas por defecto)

### Sección 1 — Primer contacto

> El orden de Q12, Q13 y Q14 es intencional: "Fecha próxima atención" y la nueva pregunta de barreras son las **últimas** de la sección, después de "Continuar Primera Atención", para que su visibilidad pueda evaluarse una vez que se conoce la respuesta al gatillo.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | ¿La atención es individual o en dupla? | single | ✅ | Individual / Dupla |
| 2 | ¿Se encuentra en un lugar seguro? | single | ✅ | Sí / No |
| 3 | ¿Se encuentra en riesgo inminente? | single | ✅ | Sí / No — *visible si Q2 = No* |
| 4 | Describa las acciones realizadas ante la situación de riesgo inminente | text | ✅ | *visible si Q3 = Sí* |
| 5 | ¿Hay voluntariedad para la atención? | single | ✅ | Sí / No |
| 6 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización — *visible si Q5 = Sí* |
| 7 | Compromisos | text | ✅ | *visible si Q5 = Sí* |
| 8 | Observaciones | text | ❌ | Siempre visible |
| 9 | Hay nuevos hechos de violencia | boolean | ✅ | — |
| 10 | Descripción de los hechos | text | ❌ | *visible si Q9 = true* |
| 11 | Fecha (de los hechos) | date | ❌ | *visible si Q9 = true* |
| 12 | Continuar Primera Atención | boolean | ✅ | **Gatillo de visibilidad de S4**: Si = true → S4 aparece en este mismo formulario |
| 13 | Fecha próxima atención | date | ✅ | *visible si Q5 = Sí **Y** Q12 = No* — se oculta cuando "Continuar = Sí" (aparece en S4) |
| 14 | Desde la atención anterior se han identificado barreras institucionales | boolean | ✅ | **ÚLTIMA** (Jul 2026) — *visible si Q12 = Sí* — **gatillo de visibilidad de S3 (Identificación de Barreras)** |

### Sección 2 — Seguimiento a Barreras *(nueva, Jul 2026)*

> **Visibilidad de sección:** No depende de ninguna pregunta de este formulario. Permanece habilitada sin importar la respuesta de Q14. Su visibilidad real (si el caso tiene barreras activas) queda pendiente de resolver vía formState externo, igual que en el Formulario de Seguimiento.
>
> Misma estructura y funcionalidad (repeater de 7 preguntas) que la Sección "Seguimiento a Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md`.

### Sección 3 — Identificación de Barreras *(nueva, Jul 2026)*

> **Visibilidad de sección:** Visible cuando Q14 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = Sí.
>
> Misma estructura y funcionalidad (repeater de 22 preguntas, min. 1 repetición) que la Sección "Identificación de Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md`.

### Sección 4 — Primera Atención

> **Visibilidad de sección:** Oculta por defecto. Se muestra cuando Q12 (S1 — Continuar Primera Atención) = true.
>
> Contiene las mismas preguntas que la Sección 4 del Form 2 (Primera Atención). El profesional no sale del formulario de Primer Contacto.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si riesgo inminente (S1-Q3) = Sí* |
| 2 | Consentimiento Informado para la Atención Psicosocial *(texto largo + pregunta)* | single | ✅ | Sí / No — texto completo del consentimiento como label |
| 3 | Confirmación consentimiento persona de apoyo | single | ✅ | Sí / No — *siempre visible (Jul 2026: antes dependía de "¿Intérprete?", pregunta eliminada)* |
| 4 | ¿La persona da su consentimiento para ser contactada posteriormente para evaluar calidad? | single | ✅ | Sí / No |
| 5 | Ingresa por conducta suicida asociada a VBG o VpP | single | ✅ | Sí / No |
| 6 | Tipo de conducta suicida | single | ❌ | Ideación / Amenaza / Intento — *visible si Q5 = Sí* |
| 7 | Contenido de la atención | text | ✅ | — |
| 8 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 9 | Plan de trabajo y recomendaciones | text | ✅ | — |
| 10 | Compromisos | text | ✅ | — |
| 11 | Fecha próxima atención | date | ✅ | *Movida desde S1 — aplica cuando Continuar Primera Atención = Sí* |
| 12 | Observaciones | text | ❌ | — |

> Jul 2026: se eliminaron "¿Requiere algún ajuste razonable...?" y "¿Requiere intérprete de idiomas...?" (antes Q2/Q3). 14 → 12 preguntas.

---

## Form 2: Primera Atención

**Secciones:** 4

> Este formulario se usa **únicamente en Escenario B**: segunda llamada independiente donde `ya_hizo_primer_contacto = true` y `ya_hizo_primera_atencion = false`. La Primera Atención realizada dentro del Primer Contacto (Escenario A, Continuar = Sí) ya no usa este formulario.

### Sección 1 — Contacto Primera Atención

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | ¿La atención es individual o en dupla? | single | ✅ | Individual / Dupla |
| 2 | ¿La llamada fue efectiva? | single | ✅ | Sí / No |
| 3 | ¿Se encuentra en un lugar seguro? | single | ✅ | Sí / No — *visible si Q2 = Sí* |
| 4 | ¿Se encuentra en riesgo inminente? | single | ✅ | Sí / No — *visible si Q3 = No* |
| 5 | Hay nuevos hechos de violencia | boolean | ✅ | *visible si Q2 = Sí* |
| 6 | Descripción de los hechos | text | ❌ | *visible si Q5 = true* |
| 7 | Fecha (de los hechos) | date | ❌ | *visible si Q5 = true* |
| 8 | ¿Es atención o solo contacto? | multiple | ✅ | Atención / Solo Contacto — **gatillo de visibilidad de Sección 4** |
| 9 | Observaciones del contacto | text | ❌ | Visible siempre |
| 10 | Fecha nueva | date | ❌ | *visible si Q8 = Solo Contacto* — para reagendar cuando no hay atención |
| 11 | Desde la atención anterior se han identificado barreras institucionales | boolean | ✅ | **ÚLTIMA** (Jul 2026) — *visible si Q8 = Atención* — **gatillo de visibilidad de Sección 3 (Identificación de Barreras)** |

### Sección 2 — Seguimiento a Barreras *(nueva, Jul 2026)*

> **Visibilidad de sección:** No depende de ninguna pregunta de este formulario. Permanece habilitada sin importar la respuesta de Q11. Su visibilidad real (si el caso tiene barreras activas) queda pendiente de resolver vía formState externo, igual que en el Formulario de Seguimiento.
>
> Misma estructura y funcionalidad (repeater de 7 preguntas) que la Sección "Seguimiento a Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md`.

### Sección 3 — Identificación de Barreras *(nueva, Jul 2026)*

> **Visibilidad de sección:** Visible cuando Q11 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = Sí.
>
> Misma estructura y funcionalidad (repeater de 22 preguntas, min. 1 repetición) que la Sección "Identificación de Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md`.

### Sección 4 — Primera Atención

> **Visibilidad de sección:** Solo se muestra cuando `¿Es atención o solo contacto? = Atención`. Si la respuesta es "Solo Contacto", esta sección queda oculta y el formulario termina en la Sección 1 (solo se muestran "Observaciones del contacto", "Fecha nueva" y la pregunta de barreras). `ya_hizo_primera_atencion` permanece en `false`.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si riesgo inminente (S1-Q4) = Sí* |
| 2 | Consentimiento Informado para la Atención Psicosocial *(texto largo + pregunta)* | single | ✅ | Sí / No — texto completo del consentimiento como label |
| 3 | Confirmación consentimiento persona de apoyo | single | ✅ | Sí / No — *siempre visible (Jul 2026: antes dependía de "¿Intérprete?", pregunta eliminada)* |
| 4 | ¿La persona da su consentimiento para ser contactada posteriormente para evaluar calidad? | single | ✅ | Sí / No |
| 5 | Ingresa por conducta suicida asociada a VBG o VpP | single | ✅ | Sí / No |
| 6 | Tipo de conducta suicida | single | ❌ | Ideación / Amenaza / Intento — *visible si Q5 = Sí* |
| 7 | Contenido de la atención | text | ✅ | — |
| 8 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 9 | Plan de trabajo y recomendaciones | text | ✅ | — |
| 10 | Compromisos | text | ✅ | — |
| 11 | Fecha próxima atención | date | ✅ | — |
| 12 | Observaciones | text | ❌ | — |

> Jul 2026: se eliminaron "¿Requiere algún ajuste razonable...?" y "¿Requiere intérprete de idiomas...?" (antes Q2/Q3). 14 → 12 preguntas.

---

## Form 3: Atención Psicosocial *(antes "Seguimiento" — renombrado Jul 2026)*

**Secciones:** 4

### Sección 1 — Contacto Atención Psicosocial *(antes "Contacto Seguimiento")*

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | ¿La atención es individual o en dupla? | single | ✅ | Individual / Dupla — **NUEVA (Jul 2026)**, antes ausente en este formulario |
| 2 | ¿La llamada fue efectiva? | single | ✅ | Sí / No |
| 3 | ¿Se encuentra en un lugar seguro? | single | ✅ | Sí / No — *visible si Q2 = Sí* |
| 4 | ¿Se encuentra en riesgo inminente? | single | ✅ | Sí / No — *visible si Q3 = No* |
| 5 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si Q4 = Sí* |
| 6 | Hay nuevos hechos de violencia | boolean | ✅ | *visible si Q2 = Sí* |
| 7 | Descripción de los hechos | text | ❌ | *visible si Q6 = true* |
| 8 | Fecha (de los hechos) | date | ❌ | *visible si Q6 = true* |
| 9 | ¿Es atención o solo contacto? | multiple | ✅ | Atención / Solo Contacto — **gatillo de visibilidad de Sección 4** |
| 10 | Observaciones del contacto | text | ❌ | Visible siempre |
| 11 | Fecha nueva | date | ❌ | *visible si Q9 = Solo Contacto* — para reagendar cuando no hay atención |
| 12 | Desde la atención anterior se han identificado barreras institucionales | boolean | ✅ | **ÚLTIMA** (Jul 2026) — *visible si Q9 = Atención* — **gatillo de visibilidad de Sección 3 (Identificación de Barreras)** |

> Jul 2026: se agregó Q1 "¿La atención es individual o en dupla?" (antes ausente); el resto de preguntas se recorrió 1 posición (11 → 12).

### Sección 2 — Seguimiento a Barreras *(nueva, Jul 2026 — nombre fijo, no se renombra)*

> **Visibilidad de sección:** No depende de ninguna pregunta de este formulario. Permanece habilitada sin importar la respuesta de Q12. Su visibilidad real (si el caso tiene barreras activas) queda pendiente de resolver vía formState externo, igual que en el Formulario de Seguimiento (hacer-seguimiento).
>
> Misma estructura y funcionalidad (repeater de 7 preguntas) que la Sección "Seguimiento a Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md`.

### Sección 3 — Identificación de Barreras *(nueva, Jul 2026)*

> **Visibilidad de sección:** Visible cuando Q12 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = Sí.
>
> Misma estructura y funcionalidad (repeater de 22 preguntas, min. 1 repetición) que la Sección "Identificación de Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md`.

### Sección 4 — Atención Psicosocial *(antes "Seguimiento")*

> **Visibilidad de sección:** Solo se muestra cuando `¿Es atención o solo contacto? = Atención`. Si la respuesta es "Solo Contacto", esta sección queda oculta. `session_count` no incrementa.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Contenido de la atención | text | ✅ | — |
| 2 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 3 | Compromisos | text | ✅ | — |
| 4 | Fecha próxima atención | date | ✅ | — |
| 5 | Observaciones | text | ❌ | — |

---

## Form 4: Cierre

**Secciones:** 5

> Tiene las mismas secciones de Seguimiento (Contacto + Barreras + Seguimiento) más la sección de Cierre. Aplica cuando `session_count >= 3`.

### Sección 1 — Contacto Cierre

*(Idéntica a Sección 1 del Form Atención Psicosocial, incluyendo la pregunta de barreras)*

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | ¿La atención es individual o en dupla? | single | ✅ | Individual / Dupla — **NUEVA (Jul 2026)**, antes ausente en este formulario |
| 2 | ¿La llamada fue efectiva? | single | ✅ | Sí / No |
| 3 | ¿Se encuentra en un lugar seguro? | single | ✅ | Sí / No — *visible si Q2 = Sí* |
| 4 | ¿Se encuentra en riesgo inminente? | single | ✅ | Sí / No — *visible si Q3 = No* |
| 5 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si Q4 = Sí* |
| 6 | Hay nuevos hechos de violencia | boolean | ✅ | *visible si Q2 = Sí* |
| 7 | Descripción de los hechos | text | ❌ | *visible si Q6 = true* |
| 8 | Fecha (de los hechos) | date | ❌ | *visible si Q6 = true* |
| 9 | ¿Es atención o solo contacto? | multiple | ✅ | Atención / Solo Contacto — **gatillo de visibilidad de Sección 4** |
| 10 | Observaciones del contacto | text | ❌ | Visible siempre |
| 11 | Fecha nueva | date | ❌ | *visible si Q9 = Solo Contacto* — para reagendar cuando no hay atención |
| 12 | Desde la atención anterior se han identificado barreras institucionales | boolean | ✅ | **ÚLTIMA** (Jul 2026) — *visible si Q9 = Atención* — **gatillo de visibilidad de Sección 3 (Identificación de Barreras)** |

> Jul 2026: se agregó Q1 "¿La atención es individual o en dupla?" (antes ausente); el resto de preguntas se recorrió 1 posición (11 → 12).

### Sección 2 — Seguimiento a Barreras *(nueva, Jul 2026 — nombre fijo, no se renombra)*

> **Visibilidad de sección:** No depende de ninguna pregunta de este formulario. Permanece habilitada sin importar la respuesta de Q12. Su visibilidad real (si el caso tiene barreras activas) queda pendiente de resolver vía formState externo, igual que en el Formulario de Seguimiento (hacer-seguimiento).
>
> Misma estructura y funcionalidad (repeater de 7 preguntas) que la Sección "Seguimiento a Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md`.

### Sección 3 — Identificación de Barreras *(nueva, Jul 2026)*

> **Visibilidad de sección:** Visible cuando Q12 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = Sí.
>
> Misma estructura y funcionalidad (repeater de 22 preguntas, min. 1 repetición) que la Sección "Identificación de Barreras" del Formulario de Seguimiento — ver `DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md`.

### Sección 4 — Atención Psicosocial (en Cierre) *(antes "Seguimiento (en Cierre)")*

> **Visibilidad de sección:** Solo se muestra cuando `¿Es atención o solo contacto? = Atención`. Si la respuesta es "Solo Contacto", esta sección y la Sección 5 quedan ocultas. `session_count` no incrementa.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Contenido de la atención | text | ✅ | — |
| 2 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 3 | Compromisos | text | ✅ | — |
| 4 | Fecha próxima atención | date | ✅ | — |
| 5 | Observaciones | text | ❌ | — |
| 6 | Cerrar remisión | boolean | ✅ | **Gatillo de visibilidad de S5**: Si = true → Sección de Cierre aparece. Si = false → S5 permanece oculta y el formulario actúa como un seguimiento más |

### Sección 5 — Cierre

> **Visibilidad de sección:** Solo se muestra cuando **ambas** condiciones se cumplen:
> 1. `¿Es atención o solo contacto? = Atención` (S1 — misma condición que S4)
> 2. `Cerrar remisión = Sí` (S4-Q6)
>
> Si "Cerrar remisión = No", el formulario termina en S4 y el registro actúa como un seguimiento regular más. `status` no cambia a `cerrado`.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Motivo de cierre | single | ✅ | Cumplimiento de objetivos / Cumplimiento esquema / No consentimiento / Imposibilidad del contacto (3x3) / Desistimiento del proceso |
| 2 | Contenido de la atención | text | ✅ | — |
| 3 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 4 | Temas trabajados durante la atención | text | ✅ | — |
| 5 | Hay nuevos hechos de violencia | boolean | ✅ | — |
| 6 | Descripción de los hechos | text | ❌ | *visible si Q5 = true* |
| 7 | Fecha (de los hechos) | date | ❌ | *visible si Q5 = true* |

---

## Resumen de opciones compartidas

### Plan de orientación (presente en Form 1 S1 y S4, Form 2 S4, Form 3 S4, Form 4 S4 y S5)

| Valor | Label |
|---|---|
| `enrutamiento` | Enrutamiento |
| `activacion_ruta` | Activación de ruta |
| `seguimiento` | Seguimiento |
| `medidas_emergencia` | Medidas de emergencia |
| `plan_estabilizacion` | Plan de estabilización |

### ¿La atención es individual o en dupla? (Form 1 S1, Form 2 S1, Form 3 S1 y Form 4 S1)

> Jul 2026: se agregó a Form 3 y Form 4 (antes solo en Form 1 y Form 2).

| Valor | Label |
|---|---|
| `individual` | Individual |
| `dupla` | Dupla |

### ¿La llamada fue efectiva? / Lugar seguro / Riesgo inminente / Consentimientos

| Valor | Label |
|---|---|
| `si` | Sí |
| `no` | No |

### ¿Es atención o solo contacto? (Form 2 S1, Form 3 S1, Form 4 S1)

| Valor | Label |
|---|---|
| `atencion` | Atención |
| `solo_contacto` | Solo Contacto |

### Tipo de conducta suicida (Form 1 S4 y Form 2 S4)

| Valor | Label |
|---|---|
| `ideacion` | Ideación |
| `amenaza` | Amenaza |
| `intento` | Intento |

### Motivo de cierre (Form 4 S5)

| Valor | Label |
|---|---|
| `cumplimiento_objetivos` | Cumplimiento de objetivos |
| `cumplimiento_esquema` | Cumplimiento esquema |
| `no_consentimiento` | No consentimiento |
| `imposibilidad_contacto` | Imposibilidad del contacto (3x3) |
| `desistimiento` | Desistimiento del proceso |

---

## Resumen de visibilidad de secciones

| Formulario | Sección | Trigger | Condición |
|---|---|---|---|
| Form 1 — Primer Contacto | S2 — Seguimiento a Barreras | *(ninguno — pendiente formState externo)* | siempre habilitada |
| Form 1 — Primer Contacto | S3 — Identificación de Barreras | S1-Q14 Desde la atención anterior... | `= true` |
| Form 1 — Primer Contacto | S4 — Primera Atención | S1-Q12 Continuar Primera Atención | `= true` |
| Form 2 — Primera Atención | S2 — Seguimiento a Barreras | *(ninguno — pendiente formState externo)* | siempre habilitada |
| Form 2 — Primera Atención | S3 — Identificación de Barreras | S1-Q11 Desde la atención anterior... | `= true` |
| Form 2 — Primera Atención | S4 — Primera Atención | S1-Q8 ¿Es atención o solo contacto? | `= atencion` |
| Form 3 — Atención Psicosocial | S2 — Seguimiento a Barreras | *(ninguno — pendiente formState externo)* | siempre habilitada |
| Form 3 — Atención Psicosocial | S3 — Identificación de Barreras | S1-Q12 Desde la atención anterior... | `= true` |
| Form 3 — Atención Psicosocial | S4 — Atención Psicosocial | S1-Q9 ¿Es atención o solo contacto? | `= atencion` |
| Form 4 — Cierre | S2 — Seguimiento a Barreras | *(ninguno — pendiente formState externo)* | siempre habilitada |
| Form 4 — Cierre | S3 — Identificación de Barreras | S1-Q12 Desde la atención anterior... | `= true` |
| Form 4 — Cierre | S4 — Atención Psicosocial (Cierre) | S1-Q9 ¿Es atención o solo contacto? | `= atencion` |
| Form 4 — Cierre | S5 — Cierre | S4-Q6 Cerrar remisión | `= true` (y S4 visible) |

> Las secciones "Seguimiento a Barreras" e "Identificación de Barreras" son idénticas (misma estructura, preguntas, opciones y repeaters) en los 4 formularios — ver `DocsMD/Screens/hacer-seguimiento/form-seguimiento-barreras-repeater.md` y `form-barreras-repeater.md`.

---

## Resumen de conteo

> Nota (Jul 2026): "Seguimiento a Barreras" agrega 1 repeater group con 7 preguntas por entrada; "Identificación de Barreras" agrega 1 repeater group con 22 preguntas por entrada (mín. 1 repetición). Estas 2 secciones se repiten idénticas en los 4 formularios.

| Formulario | Secciones | Preguntas (fuera de repeaters) | Opciones |
|---|---|---|---|
| Form 1: Primer Contacto | 4 (S3 y S4 condicionales) | 14 (S1) + 12 (S4) = 26 | ~18 |
| Form 2: Primera Atención | 4 (S3 condicional) | 11 (S1) + 12 (S4) = 23 | ~23 |
| Form 3: Atención Psicosocial | 4 (S3 condicional) | 12 (S1) + 5 (S4) = 17 | ~17 |
| Form 4: Cierre | 5 (S3 condicional) | 12 (S1) + 6 (S4) + 7 (S5) = 25 | ~22 |
| **Total (fuera de barreras)** | **17** | **~91** | **~80** |
| **+ Barreras (x4 formularios)** | 8 secciones (2 x 4) | +29 preguntas de repeater x 4 = +116 | ~15 por repeater |

> Jul 2026 (ajuste 4): Form 1/2 S4 pasó de 14 a 12 preguntas (-2, ajuste razonable/intérprete eliminadas). Form 3/4 S1 pasó de 11 a 12 preguntas (+1, "¿individual o en dupla?" agregada).
