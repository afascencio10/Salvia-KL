# `psicosocial-sesion` — Índice de Documentación

Pantalla para que el agente psicosocial registre una sesión realizada. Comparte arquitectura con `hacer-seguimiento`: tarjeta de la víctima + componente `DinamicForm`. A diferencia de `hacer-seguimiento`, usa **4 formularios independientes** (uno por tipo de sesión) en lugar de un único formulario con secciones condicionales.

**Ruta planificada:** `/salvia/psicosocial/registrar/:psicosocial_id`  
**Fuente de preguntas:** Google Sheet — _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial`

---

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`interfaz.md`](interfaz.md) | Árbol de interfaz: componentes, estados visuales, variantes del overlay de completado |
| [`flujo-psicosocial.md`](flujo-psicosocial.md) | Máquina de estados: 4 escenarios, qué formulario se carga y cómo cambia el estado del caso |
| [`related-tables.md`](related-tables.md) | Tablas de BD involucradas: tablas existentes reutilizadas + campos nuevos a migrar |
| [`form-psicosocial-data.md`](form-psicosocial-data.md) | Estructura de los 4 formularios: secciones, preguntas, opciones y condiciones de visibilidad |
| [`Flujos/flow-E01-cuando-carga-pantalla.md`](Flujos/flow-E01-cuando-carga-pantalla.md) | Evento E-01: carga inicial — selección de formulario, resolución de team_contact, formState |

---

## Decisión arquitectural: 4 formularios independientes

Recomendado por el líder técnico de `hacer-seguimiento`. Cada tipo de sesión tiene su propio formulario, lo que facilita el mantenimiento, la evolución independiente y la claridad de cada estado del proceso.

| # | Formulario | Escenario de uso | Secciones |
|---|---|---|---|
| 1 | **Primer Contacto** | `ya_hizo_primer_contacto = false` | 1 — Primer contacto |
| 2 | **Primera Atención** | `ya_hizo_primer_contacto = true`, `ya_hizo_primera_atencion = false` | 1 — Contacto · 2 — Primera atención |
| 3 | **Seguimiento** | `ya_hizo_primera_atencion = true`, `session_count` entre 1 y 2 | 1 — Contacto · 2 — Seguimiento |
| 4 | **Cierre** | `ya_hizo_primera_atencion = true`, `session_count >= 3` | 1 — Contacto · 2 — Seguimiento · 3 — Cierre |

**Variante del Escenario A:** Al completar el Formulario de Primer Contacto con "Continuar Primera Atención = Sí", el sistema carga el Formulario de Primera Atención omitiendo la sección de Contacto (`formState.skip_contact = true`).

---

## Tablas de BD involucradas

### Tablas existentes (reutilizadas con nuevos campos)

| Tabla | Operación | Campos nuevos |
|---|---|---|
| `salvia.psychosocial_support` | SELECT al cargar, UPDATE al completar | `ya_hizo_primer_contacto`, `ya_hizo_primera_atencion` |
| `salvia.team_contact` | INSERT por cada sesión registrada | `session_type` |

### Tablas del sistema DinamicForm (nuevos registros vía seed)

| Tabla | Tipo | Operación |
|---|---|---|
| `salvia.form` | Existente | INSERT 4 formularios nuevos |
| `salvia.form_section` | Existente | INSERT 8 secciones |
| `salvia.question` | Existente | INSERT ~74 preguntas |
| `salvia.option` | Existente | INSERT ~71 opciones |
| `salvia.visibility_condition` | Existente | INSERT ~15 condiciones de visibilidad |
| `salvia.form_submission` | Existente | INSERT por cada sesión iniciada |
| `salvia.answer` | Existente | INSERT por cada respuesta |
| `salvia.case_timeline_event` | Existente | INSERT al completar |

---

## Estado en `psychosocial_support` que maneja el flujo

| Campo | Tipo | Determina |
|---|---|---|
| `ya_hizo_primer_contacto` | boolean | Si se usa el Form 1 o Form 2+ |
| `ya_hizo_primera_atencion` | boolean | Si se usa el Form 3 o Form 4 |
| `session_count` | int | Diferencia entre Form 3 (Seguimiento) y Form 4 (Cierre) |
| `status` | varchar | Estado visible en el componente de remisiones |

---

## Componentes relacionados

| Componente | Relación |
|---|---|
| `remisiones-psicosocial-component` | Lista las remisiones y permite acceder a esta pantalla vía evento `ver-remision` |
| `DinamicForm` | Componente que renderiza el formulario seleccionado según el escenario |
| `hacer-seguimiento` | Pantalla de referencia — misma arquitectura de tarjeta + DinamicForm |
