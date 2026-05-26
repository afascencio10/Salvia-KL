# Formulario de Seguimiento — Estructura de Secciones y Preguntas

**Form ID:** `2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff`
**Schema DB:** `salvia`
**Tablas involucradas:** `form`, `form_section`, `question`, `option`, `visibility_condition`, `repeater_group`

---

## Secciones

| Order | ID | Nombre | Notas |
|---|---|---|---|
| 1 | `525203d6` | Valoración del Riesgo | Siempre visible |
| 2 | `0b7ff496` | Rutas Diferenciales | Siempre visible |
| 3 | `235f44f5` | Identificación de Barreras | Visible solo si Q7 de Valoración del Riesgo = `true` |
| 4 | `2a2347fd` | Seguimiento de Caso | Siempre visible |
| 5 | `3e03f685` | Cierre del caso | Siempre visible |

---

## Sección 1 — Valoración del Riesgo

**ID:** `525203d6-65c4-4ba6-ad30-70001b48a32a`

| Order | ID | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|---|
| 1 | `f8b69cd8` | `single` | ✅ | Respondiente del seguimiento | 3 opciones |
| 2 | `f8453544` | `text` | ❌ | Registre si se han presentado nuevos hechos de violencia desde el último seguimiento (Tiempo/Modo/Lugar) | |
| 3 | `a0fdcf67` | `multiple` | ❌ | Factores protectores presentes en el caso | 6 opciones. Alert si ≥3 seleccionados |
| 4 | `ec5bb242` | `multiple` | ❌ | Factores de riesgo presentes en el caso | 9 opciones. Alert si ≥4 seleccionados |
| 5 | `9cec4dcf` | `boolean` | ✅ | ¿Las acciones desplegadas han tenido efecto protector tangible en la percepción de seguridad de la mujer? | |
| 6 | `65f2d582` | `multiple` | ❌ | Factores de riesgo extremo. Si identifica uno o más, remita al equipo de riesgo. | 8 opciones. Alert si cualquiera seleccionado |
| 7 | `2eede1a4` | `boolean` | ❌ | ¿Desde la atención anterior se han identificado barreras institucionales que hayan contribuido a mantener o incrementar el riesgo? | Dispara visibilidad de sección Identificación de Barreras |
| 8 | `c2b02516` | `text` | ✅ | Describa de manera analítica cómo la integración de los factores protectores y riesgos presentes permite determinar la situación actual de riesgo | |

### Visibility Conditions — Valoración del Riesgo
- **Sección Identificación de Barreras** → visible cuando Q7 (`2eede1a4`) = `true` | `EQUALS`

### Metadata de alertas
| Pregunta | Condición | Mensaje |
|---|---|---|
| Q3 — Factores protectores | `min_3_selected` | Sugerencia: con 3 o más factores protectores y ningún riesgo extremo, considere remitir a seguimiento general |
| Q4 — Factores de riesgo | `min_4_selected` | Sugerencia: con 4 o más factores de riesgo, considere remitir al equipo de riesgo |
| Q6 — Riesgo extremo | `any_selected` | Al finalizar el seguimiento realice remisión al equipo de Riesgo |
| Q7 — Barreras | `hint_on_true` | Registre estas en el módulo de barreras |

---

## Sección 2 — Rutas Diferenciales

**ID:** `0b7ff496-5cfe-4c30-aa3f-7a68eea93bbb`

Patrón repetido: `boolean` de sector → `multiple` de instituciones (condicional).

| Order | ID | Tipo | Req | Pregunta | Condición visibilidad |
|---|---|---|---|---|---|
| 1 | `685c2c90` | `boolean` | ✅ | Sector salud | — |
| 2 | `579abd44` | `multiple` | ❌ | ¿A qué institución acudió? | Q1 = `true` \| `EQUALS` (7 opts) |
| 3 | `a337aa7f` | `boolean` | ✅ | Sector justicia | — |
| 4 | `f2505f75` | `multiple` | ❌ | ¿A qué institución acudió? | Q3 = `true` \| `EQUALS` (14 opts) |
| 5 | `d654b205` | `boolean` | ✅ | Sector protección | — |
| 6 | `d7690eb5` | `multiple` | ❌ | ¿A qué institución acudió? | Q5 = `true` \| `EQUALS` (10 opts) |
| 7 | `a2c0c480` | `boolean` | ✅ | Sector educación | — |
| 8 | `d193583a` | `multiple` | ❌ | ¿A qué institución acudió? | Q7 = `true` \| `EQUALS` (9 opts) |
| 9 | `a40d564e` | `boolean` | ✅ | Ministerio Público | — |
| 10 | `56ffd8a0` | `multiple` | ❌ | ¿A qué institución acudió? | Q9 = `true` \| `EQUALS` (4 opts) |
| 11 | `feb3bf05` | `boolean` | ✅ | ¿Acudió a algún sector/entidad territorial? | — |
| 12 | `a959c9e7` | `multiple` | ❌ | ¿A qué institución acudió? | Q11 = `true` \| `EQUALS` (8 opts) |
| 13 | `96f30975` | `boolean` | ✅ | ¿Se ha apoyado de alguna organización de la sociedad civil, comunitaria o internacional de Naciones Unidas? | — |
| 14 | `a547a193` | `multiple` | ❌ | ¿A qué institución acudió? | Q13 = `true` \| `EQUALS` (4 opts) |
| 15 | `7cbed66d` | `boolean` | ✅ | ¿Acudió a entidades clave con acciones complementarias y/o especializadas frente a la respuesta en casos de VBG/VPP? | — |
| 16 | `abd975e8` | `multiple` | ❌ | ¿A qué institución acudió? | Q15 = `true` \| `EQUALS` (9 opts) |
| 17 | `59bbfb48` | `text` | ❌ | ¿Acudió a otra institución interna (empleador o empleadora) o externa? | — |
| 18 | `19bfd6b7` | `multiple` | ❌ | ¿En el enrutamiento se empleó el enfoque diferencial? ¿Cuál? | — (11 opts) |
| 19 | `db799254` | `boolean` | ❌ | ¿Presenta barreras institucionales en cualquiera de los sectores de la ruta? | — |

---

## Sección 3 — Identificación de Barreras

**ID:** `235f44f5-106d-4b86-8b5d-e087d04fd0d9`
**Visibilidad:** Solo visible cuando Q7 de Valoración del Riesgo = `true`
**Tipo:** Contiene un `repeater_group`

### Repeater Group
| Campo | Valor |
|---|---|
| ID | `5fd3ecdc-2e5f-4b31-97ef-8a994580586a` |
| Nombre | Barreras identificadas |
| min_repetitions | 1 |
| max_repetitions | — (sin límite) |

### Preguntas (dentro del repeater)

| Order | ID | Tipo | Req | Pregunta |
|---|---|---|---|---|
| 1 | `f19378b6` | `dropdown` | ✅ | Sector de la barrera |
| 2 | `5fc1f2af` | `multiple` | ❌ | Barreras identificadas en Salud |
| 3 | `2bec977e` | `multiple` | ❌ | Barreras identificadas en Justicia |
| 4 | `66c9fc1e` | `multiple` | ❌ | Barreras identificadas en Protección |

---

## Sección 4 — Seguimiento de Caso

**ID:** `2a2347fd-19d2-408c-a5de-02e854418479`

| Order | ID | Tipo | Req | Pregunta | Condición visibilidad |
|---|---|---|---|---|---|
| 1 | `bb7a2307` | `text` | ✅ | Gestión realizada en el seguimiento | — |
| 2 | `fc64033e` | `boolean` | ✅ | ¿Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA? | — |
| 3 | `e0d38cf5` | `multiple` | ✅ | ¿Cuáles equipos? | Q2 = `true` \| `EQUALS` (5 opts) |
| 4 | `1a36260c` | `multiple` | ❌ | ¿Cuáles medidas de emergencia? | Q3 CONTAINS `medidas_emergencia` (6 opts) |
| 5 | `71c42c4a` | `multiple` | ❌ | Criterios de remisión — Atención Psicosocial | Q3 CONTAINS `atencion_psico` (7 opts) |
| 6 | `f7edf4fc` | `text` | ❌ | Criterios de remisión — Atención Hombres | Q3 CONTAINS `atencion_hombres` |
| 7 | `47b122b1` | `multiple` | ❌ | Criterios de remisión — Salvia Dignidad | Q3 CONTAINS `salvia_dignidad` (2 opts) |
| 8 | `28accaa6` | `multiple` | ❌ | Criterios de remisión — Estabilización | Q3 CONTAINS `estabilizacion` (3 opts) |
| 9 | `1bdc8b52` | `text` | ✅ | Describa los elementos que evidencia para realizar la remisión | Q2 = `true` \| `EQUALS` |

### Opciones de Q3 — ¿Cuáles equipos?
| Value | Label |
|---|---|
| `atencion_psico` | Atención Psicosocial |
| `medidas_emergencia` | Medidas de Emergencia |
| `estabilizacion` | Estabilización |
| `atencion_hombres` | Atención Hombres |
| `salvia_dignidad` | Salvia Dignidad |

---

## Sección 5 — Cierre del caso

**ID:** `3e03f685-0a41-4070-8ca7-1bc5f1e6b699`

| Order | ID | Tipo | Req | Pregunta | Condición visibilidad |
|---|---|---|---|---|---|
| 1 | `08950a38` | `boolean` | ✅ | ¿Realiza cierre del caso? | — |
| 2 | `95fb963e` | `single` | ❌ | Motivo del cierre | Q1 = `true` \| `EQUALS` (5 opts) |
| 3 | `e259ff16` | `text` | ❌ | Otro motivo ¿cuál? | Q2 CONTAINS `otro` |
| 4 | `50ab05f3` | `text` | ❌ | Describa la causa del cierre | Q1 = `true` \| `EQUALS` |
| 5 | `0e7b61c8` | `boolean` | ❌ | ¿Realizó acciones institucionales por el cierre? | Q1 = `true` \| `EQUALS` |

### Opciones de Q2 — Motivo del cierre
| Value | Label |
|---|---|
| `perdida_contacto` | Pérdida de contacto |
| `solicitud_expresa` | Solicitud expresa de la persona de finalizar el proceso |
| `cumplimiento_plan` | Cumplimiento del plan de atención |
| `no_corresponde` | No corresponde al ámbito, población o naturaleza de la atención |
| `otro` | Otro ¿cuál? |

---

## Resumen general

| Sección | Preguntas | Visibility Conditions | Repeater |
|---|---|---|---|
| Valoración del Riesgo | 8 | 1 (sobre sección) | No |
| Rutas Diferenciales | 19 | 8 | No |
| Identificación de Barreras | 4 | — | Sí (min 1) |
| Seguimiento de Caso | 9 | 6 | No |
| Cierre del caso | 5 | 4 | No |
| **Total** | **45** | **19** | |
