# Estructura de los 4 Formularios Psicosociales

Fuente: Sheet _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial`.

**Decisiones del líder (columna "Comentarios Felipe"):**
- `"Seleccionar número de dupla"` → **eliminada de todos los formularios** ("No va → Dupla asignada en sistema")
- `"¿La llamada fue efectiva?"` en Primer Contacto → **eliminada** ("No va → esto es antes de contestar")
- `"¿La atención es individual o en dupla?"` → **se mantiene** ("Aca contesta")

**Decisión del líder — ajuste 2 (Jul 2026):**
- Si "Continuar Primera Atención = Sí", **no se redirige** al formulario de Primera Atención. En su lugar, se muestra una nueva sección **S2 — Primera Atención** dentro del mismo formulario de Primer Contacto.
- La pregunta "Fecha próxima atención" de S1 desaparece cuando "Continuar Primera Atención = Sí" (se mueve a S2).
- En el formulario de Cierre, se agrega la pregunta "Cerrar remisión" en S2 (Seguimiento). La sección S3 (Cierre) solo se muestra si "Cerrar remisión = Sí".

---

## Form 1: Primer Contacto

**Secciones:** 2 (S2 oculta por defecto)

### Sección 1 — Primer contacto

> El orden de Q12 y Q13 es intencional: "Fecha próxima atención" es la **última** pregunta de la sección, después de "Continuar Primera Atención", para que su visibilidad pueda evaluarse una vez que se conoce la respuesta al gatillo.

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
| 12 | Continuar Primera Atención | boolean | ✅ | **Gatillo de visibilidad de S2**: Si = true → S2 aparece en este mismo formulario |
| 13 | Fecha próxima atención | date | ✅ | **ÚLTIMA** — *visible si Q5 = Sí **Y** Q12 = No* — se oculta cuando "Continuar = Sí" (aparece en S2) |

### Sección 2 — Primera Atención

> **Visibilidad de sección:** Oculta por defecto. Se muestra cuando Q13 (S1 — Continuar Primera Atención) = true.
>
> Contiene las mismas preguntas que la Sección 2 del Form 2 (Primera Atención). El profesional no sale del formulario de Primer Contacto.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si riesgo inminente (S1-Q3) = Sí* |
| 2 | ¿Requiere algún ajuste razonable en el marco de la atención? | single | ✅ | Sí / No |
| 3 | ¿Requiere intérprete de idiomas y/o traducción? | single | ✅ | Sí / No — *visible si Q2 = Sí* |
| 4 | Consentimiento Informado para la Atención Psicosocial *(texto largo + pregunta)* | single | ✅ | Sí / No — texto completo del consentimiento como label |
| 5 | Confirmación consentimiento persona de apoyo | single | ✅ | Sí / No — *visible si Q3 = Sí* |
| 6 | ¿La persona da su consentimiento para ser contactada posteriormente para evaluar calidad? | single | ✅ | Sí / No |
| 7 | Ingresa por conducta suicida asociada a VBG o VpP | single | ✅ | Sí / No |
| 8 | Tipo de conducta suicida | single | ❌ | Ideación / Amenaza / Intento — *visible si Q7 = Sí* |
| 9 | Contenido de la atención | text | ✅ | — |
| 10 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 11 | Plan de trabajo y recomendaciones | text | ✅ | — |
| 12 | Compromisos | text | ✅ | — |
| 13 | Fecha próxima atención | date | ✅ | *Movida desde S1 — aplica cuando Continuar Primera Atención = Sí* |
| 14 | Observaciones | text | ❌ | — |

---

## Form 2: Primera Atención

**Secciones:** 2

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
| 8 | ¿Es atención o solo contacto? | multiple | ✅ | Atención / Solo Contacto — **gatillo de visibilidad de Sección 2** |
| 9 | Observaciones del contacto | text | ❌ | Visible siempre |
| 10 | Fecha nueva | date | ❌ | *visible si Q8 = Solo Contacto* — para reagendar cuando no hay atención |

### Sección 2 — Primera Atención

> **Visibilidad de sección:** Solo se muestra cuando `¿Es atención o solo contacto? = Atención`. Si la respuesta es "Solo Contacto", esta sección queda oculta y el formulario termina en la Sección 1 (solo se muestran "Observaciones del contacto" y "Fecha nueva"). `ya_hizo_primera_atencion` permanece en `false`.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si riesgo inminente (S1-Q4) = Sí* |
| 2 | ¿Requiere algún ajuste razonable en el marco de la atención? | single | ✅ | Sí / No |
| 3 | ¿Requiere intérprete de idiomas y/o traducción? | single | ✅ | Sí / No — *visible si Q2 = Sí* |
| 4 | Consentimiento Informado para la Atención Psicosocial *(texto largo + pregunta)* | single | ✅ | Sí / No — texto completo del consentimiento como label |
| 5 | Confirmación consentimiento persona de apoyo | single | ✅ | Sí / No — *visible si Q3 = Sí* |
| 6 | ¿La persona da su consentimiento para ser contactada posteriormente para evaluar calidad? | single | ✅ | Sí / No |
| 7 | Ingresa por conducta suicida asociada a VBG o VpP | single | ✅ | Sí / No |
| 8 | Tipo de conducta suicida | single | ❌ | Ideación / Amenaza / Intento — *visible si Q7 = Sí* |
| 9 | Contenido de la atención | text | ✅ | — |
| 10 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 11 | Plan de trabajo y recomendaciones | text | ✅ | — |
| 12 | Compromisos | text | ✅ | — |
| 13 | Fecha próxima atención | date | ✅ | — |
| 14 | Observaciones | text | ❌ | — |

---

## Form 3: Seguimiento

**Secciones:** 2

### Sección 1 — Contacto Seguimiento

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | ¿La llamada fue efectiva? | single | ✅ | Sí / No |
| 2 | ¿Se encuentra en un lugar seguro? | single | ✅ | Sí / No — *visible si Q1 = Sí* |
| 3 | ¿Se encuentra en riesgo inminente? | single | ✅ | Sí / No — *visible si Q2 = No* |
| 4 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si Q3 = Sí* |
| 5 | Hay nuevos hechos de violencia | boolean | ✅ | *visible si Q1 = Sí* |
| 6 | Descripción de los hechos | text | ❌ | *visible si Q5 = true* |
| 7 | Fecha (de los hechos) | date | ❌ | *visible si Q5 = true* |
| 8 | ¿Es atención o solo contacto? | multiple | ✅ | Atención / Solo Contacto — **gatillo de visibilidad de Sección 2** |
| 9 | Observaciones del contacto | text | ❌ | Visible siempre |
| 10 | Fecha nueva | date | ❌ | *visible si Q8 = Solo Contacto* — para reagendar cuando no hay atención |

### Sección 2 — Seguimiento

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

**Secciones:** 3

> Tiene las mismas secciones de Seguimiento (Contacto + Seguimiento) más la sección de Cierre. Aplica cuando `session_count >= 3`.

### Sección 1 — Contacto Cierre

*(Idéntica a Sección 1 del Form Seguimiento)*

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | ¿La llamada fue efectiva? | single | ✅ | Sí / No |
| 2 | ¿Se encuentra en un lugar seguro? | single | ✅ | Sí / No — *visible si Q1 = Sí* |
| 3 | ¿Se encuentra en riesgo inminente? | single | ✅ | Sí / No — *visible si Q2 = No* |
| 4 | Describa las acciones ante la situación de riesgo inminente | text | ✅ | *visible si Q3 = Sí* |
| 5 | Hay nuevos hechos de violencia | boolean | ✅ | *visible si Q1 = Sí* |
| 6 | Descripción de los hechos | text | ❌ | *visible si Q5 = true* |
| 7 | Fecha (de los hechos) | date | ❌ | *visible si Q5 = true* |
| 8 | ¿Es atención o solo contacto? | multiple | ✅ | Atención / Solo Contacto — **gatillo de visibilidad de Sección 2** |
| 9 | Observaciones del contacto | text | ❌ | Visible siempre |
| 10 | Fecha nueva | date | ❌ | *visible si Q8 = Solo Contacto* — para reagendar cuando no hay atención |

### Sección 2 — Seguimiento (en Cierre)

> **Visibilidad de sección:** Solo se muestra cuando `¿Es atención o solo contacto? = Atención`. Si la respuesta es "Solo Contacto", esta sección y la Sección 3 quedan ocultas. `session_count` no incrementa.

| # | Pregunta | Tipo | Req | Opciones / Notas |
|---|---|---|---|---|
| 1 | Contenido de la atención | text | ✅ | — |
| 2 | Plan de orientación | multiple | ❌ | Enrutamiento / Activación de ruta / Seguimiento / Medidas de emergencia / Plan de estabilización |
| 3 | Compromisos | text | ✅ | — |
| 4 | Fecha próxima atención | date | ✅ | — |
| 5 | Observaciones | text | ❌ | — |
| 6 | Cerrar remisión | boolean | ✅ | **Gatillo de visibilidad de S3**: Si = true → Sección de Cierre aparece. Si = false → S3 permanece oculta y el formulario actúa como un seguimiento más |

### Sección 3 — Cierre

> **Visibilidad de sección:** Solo se muestra cuando **ambas** condiciones se cumplen:
> 1. `¿Es atención o solo contacto? = Atención` (S1 — misma condición que S2)
> 2. `Cerrar remisión = Sí` (S2-Q6)
>
> Si "Cerrar remisión = No", el formulario termina en S2 y el registro actúa como un seguimiento regular más. `status` no cambia a `cerrado`.

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

### Plan de orientación (presente en Form 1 S1 y S2, Form 2 S2, Form 3 S2, Form 4 S2 y S3)

| Valor | Label |
|---|---|
| `enrutamiento` | Enrutamiento |
| `activacion_ruta` | Activación de ruta |
| `seguimiento` | Seguimiento |
| `medidas_emergencia` | Medidas de emergencia |
| `plan_estabilizacion` | Plan de estabilización |

### ¿La atención es individual o en dupla? (Form 1 S1 y Form 2 S1)

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

### Tipo de conducta suicida (Form 1 S2 y Form 2 S2)

| Valor | Label |
|---|---|
| `ideacion` | Ideación |
| `amenaza` | Amenaza |
| `intento` | Intento |

### Motivo de cierre (Form 4 S3)

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
| Form 1 — Primer Contacto | S2 — Primera Atención | S1-Q12 Continuar Primera Atención | `= true` |
| Form 2 — Primera Atención | S2 — Primera Atención | S1-Q8 ¿Es atención o solo contacto? | `= atencion` |
| Form 3 — Seguimiento | S2 — Seguimiento | S1-Q8 ¿Es atención o solo contacto? | `= atencion` |
| Form 4 — Cierre | S2 — Seguimiento (Cierre) | S1-Q8 ¿Es atención o solo contacto? | `= atencion` |
| Form 4 — Cierre | S3 — Cierre | S2-Q6 Cerrar remisión | `= true` (y S2 visible) |

---

## Resumen de conteo

| Formulario | Secciones | Preguntas | Opciones |
|---|---|---|---|
| Form 1: Primer Contacto | 2 (S2 condicional) | 13 (S1) + 14 (S2) = 27 | ~20 |
| Form 2: Primera Atención | 2 | 10 (S1) + 14 (S2) = 24 | ~25 |
| Form 3: Seguimiento | 2 | 10 (S1) + 5 (S2) = 15 | ~15 |
| Form 4: Cierre | 3 | 10 (S1) + 6 (S2) + 7 (S3) = 23 | ~20 |
| **Total** | **9** | **~89** | **~80** |
