# `psicosocial-sesion` — Índice de Documentación

Pantalla para que el agente psicosocial registre una sesión realizada. Comparte arquitectura con `hacer-seguimiento`: tarjeta de la víctima + componente `DinamicForm`. Usa **4 formularios independientes** (uno por tipo de sesión).

**Ruta:** `/salvia/psicosocial/registrar/:psicosocial_id`  
**Fuente de preguntas:** Google Sheet — _Psicosocial Kreivo27.05.2026_, hoja `Formularios Psicosocial`

---

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`interfaz.md`](interfaz.md) | Árbol de interfaz, overlay de completado |
| [`flujo-psicosocial.md`](flujo-psicosocial.md) | Máquina de estados A–D, formularios, efectos al completar |
| [`related-tables.md`](related-tables.md) | Tablas BD + migraciones |
| [`form-psicosocial-data.md`](form-psicosocial-data.md) | Secciones, preguntas, visibilidad (incl. Aug 2026) |
| [`changelogAug2026.md`](changelogAug2026.md) | Barreras §9.6, consentimiento info+hide, hechos, agenda |
| [`Flujos/flow-E01-cuando-carga-pantalla.md`](Flujos/flow-E01-cuando-carga-pantalla.md) | Carga + selección de formulario |
| [`Flujos/flow-E02-cuando-se-guarda-formulario.md`](Flujos/flow-E02-cuando-se-guarda-formulario.md) | Guardado: session_type, agenda, barreras, hechos, timeline |

---

## Formularios según estado (`psychosocial_support`)

| Escenario | Condición | Formulario |
|---|---|---|
| A | `ya_hizo_primer_contacto = false` | Primer Contacto (S4 Primera Atención si Continuar=Sí, **mismo form**) |
| B | PC hecho, PA no | Primera Atención |
| C | PA hecho, `session_count < 3` | Atención Psicosocial |
| D | PA hecho, `session_count >= 3` | Cierre |

> No hay `skip_contact` ni redirección a otro form al Continuar PA: todo queda en Form PC.

| Form | Secciones |
|---|---|
| 1 Primer Contacto | Contacto · Seg. Barreras · Identif. Barreras · Primera Atención |
| 2 Primera Atención | Contacto · Seg. Barreras · Identif. Barreras · Primera Atención |
| 3 Atención Psicosocial | Contacto · Seg. Barreras · Identif. Barreras · Atención Psicosocial |
| 4 Cierre | Contacto · Seg. Barreras · Identif. Barreras · Atención (Cierre) · Cierre |

---

## Cambios Aug 2026 (diagrama)

- Consentimiento PC/PA: `info` (texto) + `single` (Sí/No); si No → oculta resto de la sección.
- Agenda: ¿Agendar? → Fecha + Hora; chequeo disponibilidad 2h; skip schedule si ocupado.
- Seguimiento a Barreras §9.6: `barrier_follow_up` + timeline + `MANAGED`.
- Hechos de violencia → timeline `Hechos del caso`.

Seed/migración: [`src/cmd/seed/seed_psicosocial.sql`](../../../src/cmd/seed/seed_psicosocial.sql), [`migrate_psicosocial_aug2026.sql`](../../../src/cmd/seed/migrate_psicosocial_aug2026.sql).

---

## Componentes relacionados

| Componente | Relación |
|---|---|
| `remisiones-psicosocial-component` | Acceso vía remisión |
| `DinamicForm` | Render del formulario (tipos incl. `time`, `info`) |
| `hacer-seguimiento` | Referencia de barreras / hechos / tareas |
