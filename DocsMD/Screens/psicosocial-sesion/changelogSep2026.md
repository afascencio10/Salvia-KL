# Changelog Septiembre 2026 — `psicosocial-sesion`

---

## 2026-09-17 — OSIGHN sin tipo ni número de documento

Se retiran **Tipo de documento** y **Número de documento** de la sección OSIGHN en los 4 formularios. Siguen el banner de nombre identitario vs jurídico, nombres, pronombres, identidad de género, sexo asignado y orientación sexual (8 preguntas, order 1–8). Soft-delete (`deleted_at`); no se reinsertan al re-ejecutar el migrate original.

| Área | Archivos |
|---|---|
| Docs | `form-psicosocial-data.md`, este changelog |
| Seed | [`src/cmd/seed/migrate_psicosocial_osighn_drop_doc_sep2026.sql`](../../../src/cmd/seed/migrate_psicosocial_osighn_drop_doc_sep2026.sql); ajuste de [`migrate_psicosocial_osighn_sep2026.sql`](../../../src/cmd/seed/migrate_psicosocial_osighn_sep2026.sql) |
| Constantes | [`src/internal/constants/psicosocial_osighn_questions.go`](../../../src/internal/constants/psicosocial_osighn_questions.go) |

---

## 2026-09-16 — Sección OSIGHN en los 4 formularios

Nueva sección **OSIGHN** (orientación sexual e identidad de género no hegemónicas) como **segunda sección** de cada formulario psicosocial. Preguntas tomadas de la hoja `FORMULARIO OSIGHN` del Excel de Drive (14-09-2026).

### Gatillo en Contacto

Tras «Continuar Primera Atención» (PC) o «¿Es atención o solo contacto?» (PA / Atención / Cierre) se muestra:

**¿Tiene una orientación sexual e identidad de género no hegemónicas (OSIGHN)?** (`single` Sí/No)

- PC: visible si Continuar Primera Atención = `true`
- PA / SEG / CIE: visible si ¿Es atención o solo contacto? = `atencion`

Si = **Sí** → se muestra la sección OSIGHN (order 2). Si = **No** o no hay atención, la sección permanece oculta. `dinamic-form` no exige respuestas de secciones ocultas.

### Sección OSIGHN (idéntica en los 4 forms)

Banner `info` sobre nombre identitario vs jurídico, nombre identitario, nombre en documento, pronombres, identidad de género (+ «¿Cuál?» si Otro), sexo asignado al nacer, orientación sexual. Tipo y número de documento se retiraron el 2026-09-17.

Las respuestas viven en `form_submission` / `answer`. No se copian a ficha de víctima.

### Orden de secciones (catálogo vivo + OSIGHN)

PC: Contacto · **OSIGHN** · Seg. Barreras · Identif. Barreras · Primera Atención · Derivación · Cierre  
PA / Atención / Cierre: Contacto · **OSIGHN** · Barreras · Entidades · contenido · Derivación · Cierre

### Archivos

| Área | Archivos |
|---|---|
| Docs | `form-psicosocial-data.md`, `flujo-psicosocial.md`, índice, este changelog |
| Seed | [`src/cmd/seed/migrate_psicosocial_osighn_sep2026.sql`](../../../src/cmd/seed/migrate_psicosocial_osighn_sep2026.sql) |
| Constantes | [`src/internal/constants/psicosocial_osighn_questions.go`](../../../src/internal/constants/psicosocial_osighn_questions.go) |
