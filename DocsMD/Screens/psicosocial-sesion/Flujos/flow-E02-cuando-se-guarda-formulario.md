━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se guarda el formulario psicosocial
   Tipo: Frontend → Backend
   Estado: IMPLEMENTADO (Jul 2026) — incluye creación de barrier_v2/case_task desde
           "Identificación de Barreras" (sección 9) y resolución de "barreras activas" para
           "Seguimiento a Barreras" vía barrier_v2.team_contact_id (sección 9.5).
           ⚠️ REQUIERE ACCIÓN MANUAL: falta correr el ALTER TABLE de la columna
           team_contact_id contra Supabase (ver sección 9.5) antes de poder guardar barreras.
           Pendiente: actualizar/cerrar barreras desde "Seguimiento a Barreras" (sección 9.6).
           Bug corregido: "Continuar Primera Atención = Sí" no completaba la sesión por un dato
           de seed incorrecto (pregunta de Consentimiento como 'info' en vez de 'single' —
           sección 10).
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Equivalente psicosocial de `DocsMD/Screens/hacer-seguimiento/Flujos/flow-E04-cuando-se-procesa-submission.md`.
Reutiliza el mismo mecanismo genérico de `dinamic-form` + `form_service.go`, no una ruta nueva.

---

## 1. Cómo emite el evento `dinamic-form` (confirmado en código)

`dinamic-form.js` no tiene un evento "guardar formulario completo" explícito — cada clic en
"Siguiente" / "Guardar Seguimiento" llama a `saveSection()`, que hace `POST` al backend con las
respuestas de la sección actual. El backend (`FormService.SaveSection`) decide si el formulario
quedó completo:

```1719:1724:src/frontend/js/components/dinamic-form.js
const allAnswered = this.formStructure.sections
    ...
if (allAnswered) {
    console.log('[saveSection] formulario completado — emitiendo form-completed');
    this.$emit('form-completed');
}
```

El backend hace el mismo cálculo de forma independiente (fuente de verdad):

```1882:1899:src/salvia/service/form_service.go
allAnswered := true
...
if allAnswered {
    log.Printf("[saveSection] submissionId=%s → formulario COMPLETO, disparando OnEndFormSubmission", submissionID)
    actorID := input.ActorID
    if err := s.OnEndFormSubmission(ctx, input.FormID, submissionID, actorID); err != nil {
        log.Printf("[OnEndFormSubmission] error: %v\n", err)
    }
}
```

`registrar_sesion.html` ya escucha `@form-completed="onFormCompleted"` (implementado en el
evento E-01) — solo muestra el overlay de "sesión completada" y no ejecuta lógica de negocio.
**Toda la lógica de negocio vive en el backend**, disparada por `OnEndFormSubmission`.

### Punto de enganche real

```1904:1925:src/salvia/service/form_service.go
const (
    seguimientoFormID    = "2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff"
    barrierUpdateFormID  = "4d0aeb46-5af3-4c47-a0d5-c5bfc4d549ff"
    cierreCasoFormID     = "da8423ab-1a8c-47db-96b7-d10496df571a"
)

func (s *formService) OnEndFormSubmission(ctx context.Context, formID, submissionID, actorID string) error {
    switch formID {
    case seguimientoFormID:
        return s.processFollowUpSubmission(ctx, submissionID, actorID)
    case barrierUpdateFormID:
        return s.processBarrierUpdateSubmission(ctx, submissionID)
    case cierreCasoFormID:
        return s.processCaseClosureSubmission(ctx, submissionID, actorID)
    }
    return nil
}
```

**Plan:** agregar los 4 form IDs psicosociales (ya existentes en
`src/internal/constants/psicosocial_forms.go`) a este switch, delegando a una función nueva
`processPsicosocialSessionSubmission(ctx, formID, submissionID, actorID)`.

```go
case constants.FormIDPrimerContacto, constants.FormIDPrimeraAtencion,
     constants.FormIDAtencionPsicosocial, constants.FormIDCierre:
    return s.processPsicosocialSessionSubmission(ctx, formID, submissionID, actorID)
```

No se necesita ninguna ruta nueva ni cambio en el frontend — `SaveSection` ya es la ruta que
usa `registrar_sesion.html` a través de `dinamic-form` (la misma que usa `hacer_seguimiento.html`).

---

## 2. Gap de infraestructura: falta un repo para `team_contact`

`form_service.go` es 100% repo-based (`FormServiceDeps`); no tiene acceso directo a `*gorm.DB`.
Ya existe `PsychosocialSupportRepo repository.PsychosocialSupportRepository` inyectado (usado en
el PASO 4 de `processFollowUpSubmission` para crear derivaciones), pero **no existe un repo de
`team_contact`**. `psychosocial_detail_service.go` administra `team_contact` con `*gorm.DB` directo
porque es un servicio distinto, más nuevo, que no sigue (todavía) el patrón de repos.

**Plan:** crear `internal/repository/team_contact_repository.go` siguiendo el patrón exacto de
`psychosocial_support_repository.go`:

```go
type TeamContactRepository interface {
    Repository[models.TeamContact]
    FindByFormSubmissionID(ctx context.Context, submissionID string) (*models.TeamContact, error)
}
```

`Repository[T]` (genérico, ya existente en `base_repository.go`) da gratis `FindByID`, `Update`,
`UpdateFields`, `Create`. Solo hay que implementar el finder por `form_submission_id`.

Agregar `TeamContactRepo repository.TeamContactRepository` a `FormServiceDeps` / `formService`,
y wiring en `main.go` (mismo lugar donde se construye `PsychosocialSupportRepo`).

---

## 3. `processPsicosocialSessionSubmission` — lógica propuesta

Espejo de `processFollowUpSubmission`, pero resolviendo el contexto desde `team_contact` en vez
de `follow_up_v2`.

```
PASO 1 — Cargar el team_contact por FormSubmissionID
  tc = teamContactRepo.FindByFormSubmissionID(submissionID)
  SI no existe → error (no debería pasar; el submission siempre nace de un team_contact
                 creado en LoadSession)

  SI tc.IsCompleted == true:
    → Ya procesado (idempotente). Registrar evento "Sesión Editada" en el timeline (igual
      que PASO de processFollowUpSubmission con fu.Status == REALIZADO) y TERMINAR.

PASO 2 — Cargar psychosocial_support por tc.PsicosocialID
  ps = psychosocialSupportRepo.FindByID(*tc.PsicosocialID)

PASO 3 — Construir answerMap desde las respuestas directas del submission
  answers = answerRepo.FindDirectBySubmissionID(submissionID)
  answerMap = { questionId: value }

PASO 4 — Determinar sessionType según el formulario (formID) + respuestas clave
  (tabla completa en la sección 4 de este documento — IDs reales ya capturados de Supabase)

PASO 5 — Actualizar team_contact
  tc.SessionType = sessionType
  tc.IsCompleted = true
  tc.CompletedAt = now()
  SI sessionType == CONTACTO_SIN_ATENCION:
    tc.IsPsicoSession = false   // no cuenta como sesión para el contador 1-6

PASO 6 — Actualizar psychosocial_support según sessionType (tabla PASO 5, sección 4)

PASO 7 — INSERT case_timeline_event (categoría Psicosocial)

PASO 8 — (Fuera de alcance de este plan, ver sección 5) Procesar repeaters de Barreras
```

---

## 4. Determinación de `sessionType` por formulario — IDs reales (capturados de Supabase, 2026-07-15)

> Los form IDs (`constants.FormIDPrimerContacto`, etc.) ya están en
> `src/internal/constants/psicosocial_forms.go`. Los IDs de pregunta abajo son nuevos —
> capturados directamente de la BD porque el seed usa `gen_random_uuid()` (no son literales
> fijos como en `seed_seguimiento.sql`).

### Form Primer Contacto (`FormIDPrimerContacto`)

| Constante propuesta | ID | Pregunta | Valores |
|---|---|---|---|
| `qPCContinuarPA` | `aa46351a-f6c0-4666-af3e-e0f1f7b6dc8c` | Continuar Primera Atención (S1) | boolean `"true"`/`"false"` |
| `qPCConsentimiento` | `c3296c87-d7e3-4ea1-8a0f-cbf77c561d0f` | Consentimiento Informado (S4, solo si Continuar=true) | `si` / `no` |

```
SI answerMap[qPCContinuarPA] == "true":
    SI answerMap[qPCConsentimiento] == "si":
        sessionType = PRIMER_CONTACTO_CON_ATENCION
    SINO:
        sessionType = PRIMER_CONTACTO_SIN_CONSENTIMIENTO  // Continuar=Sí pero sin consentimiento
SINO:
    sessionType = PRIMER_CONTACTO
```

> **Decisión del líder (resuelta, ver sección 6 punto A):** Continuar=Sí + Consentimiento=No
> avanza igual que un Primer Contacto regular (`ya_hizo_primer_contacto=true`,
> `status=en_gestion`, `ya_hizo_primera_atencion` permanece `false`). Se registra un
> `sessionType` distinto (`PRIMER_CONTACTO_SIN_CONSENTIMIENTO`) solo para trazabilidad/reporting;
> el efecto en `psychosocial_support` es idéntico al de `PRIMER_CONTACTO`.

### Form Primera Atención (`FormIDPrimeraAtencion`)

| Constante propuesta | ID | Pregunta | Valores |
|---|---|---|---|
| `qPAEsAtencion` | `3117fd01-6a31-4595-9888-dcdcc96d84e7` | ¿Es atención o solo contacto? (S1) | `atencion` / `solo_contacto` |
| `qPAConsentimiento` | `7256b91e-861b-48b3-9216-2ac13f0ae889` | Consentimiento Informado (S4) | `si` / `no` |
| `qPAFechaNueva` | `64d63b79-edee-464b-be56-1104efd31a46` | Fecha nueva (S1, solo si Solo Contacto) | date |

```
SI answerMap[qPAEsAtencion] == "solo_contacto":
    sessionType = CONTACTO_SIN_ATENCION
SINO SI answerMap[qPAConsentimiento] == "si":
    sessionType = PRIMERA_ATENCION
SINO:
    sessionType = CIERRE_NO_CONSENTIMIENTO
```

### Form Atención Psicosocial (`FormIDAtencionPsicosocial`)

| Constante propuesta | ID | Pregunta | Valores |
|---|---|---|---|
| `qSEGEsAtencion` | `16de7674-4349-4acc-9d0d-656d1d5e5f10` | ¿Es atención o solo contacto? (S1) | `atencion` / `solo_contacto` |
| `qSEGFechaNueva` | `72ce49d2-f853-4f4a-9f1a-f795f4d514c3` | Fecha nueva (S1, solo si Solo Contacto) | date |

```
SI answerMap[qSEGEsAtencion] == "solo_contacto":
    sessionType = CONTACTO_SIN_ATENCION
SINO:
    sessionType = ATENCION_PSICOSOCIAL   // (enum sin cambios; "SEGUIMIENTO" en flujo-psicosocial.md
                                          //  es el nombre pre-rebranding — ver sección 6, punto B)
```

### Form Cierre (`FormIDCierre`)

| Constante propuesta | ID | Pregunta | Valores |
|---|---|---|---|
| `qCIEEsAtencion` | `be185336-e025-4768-a1a2-3ab9e2281450` | ¿Es atención o solo contacto? (S1) | `atencion` / `solo_contacto` |
| `qCIEFechaNueva` | `1b4d09f0-e5ca-4b4e-9d74-428478a113c6` | Fecha nueva (S1, solo si Solo Contacto) | date |
| `qCIECerrarRemision` | `2edd39af-45da-4d24-ac67-fdc8b9b0b6ac` | Cerrar remisión (S4) | boolean `"true"`/`"false"` |
| `qCIEMotivoCierre` | `3dd8fd25-20d9-476b-9bfd-495a239c91e6` | Motivo de cierre (S5, solo si Cerrar=Sí) | `cumplimiento_objetivos` / `cumplimiento_esquema` / `no_consentimiento` / `imposibilidad_contacto` / `desistimiento` |

```
SI answerMap[qCIEEsAtencion] == "solo_contacto":
    sessionType = CONTACTO_SIN_ATENCION
SINO SI answerMap[qCIECerrarRemision] == "true":
    sessionType = CIERRE
SINO:
    sessionType = ATENCION_PSICOSOCIAL   // actúa como un seguimiento más, sin cerrar
```

---

## 5. Actualización de `psychosocial_support` según `sessionType`

| `sessionType` | `ya_hizo_primer_contacto` | `ya_hizo_primera_atencion` | `session_count` | `status` | `scheduled_at` |
|---|---|---|---|---|---|
| `PRIMER_CONTACTO` | `true` | sin cambio (`false`) | sin cambio | `en_gestion` | — |
| `PRIMER_CONTACTO_CON_ATENCION` *(nuevo)* | `true` | `true` | `+1` | `en_gestion` | — |
| `PRIMER_CONTACTO_SIN_CONSENTIMIENTO` *(nuevo)* | `true` | sin cambio (`false`) | sin cambio | `en_gestion` | — |
| `PRIMERA_ATENCION` | sin cambio | `true` | `+1` | `en_gestion` | — |
| `CONTACTO_SIN_ATENCION` | sin cambio | sin cambio | sin cambio | sin cambio | `= fecha nueva` del formulario |
| `ATENCION_PSICOSOCIAL` | sin cambio | sin cambio | `+1` | sin cambio | — |
| `CIERRE` | sin cambio | sin cambio | `+1` | `cerrado` | — |
| `CIERRE_NO_CONSENTIMIENTO` *(nuevo)* | sin cambio | sin cambio (`false`) | sin cambio | `en_devolucion` | — |

> Coincide con el borrador ya existente en `flujo-psicosocial.md` (sección "Lógica de
> actualización del estado"), ajustado con los 2 casos nuevos que surgen de las preguntas reales
> (`PRIMER_CONTACTO_SIN_CONSENTIMIENTO`) y renombrado `SEGUIMIENTO` → `ATENCION_PSICOSOCIAL` para
> ser consistente con el rebranding Jul 2026 y con las constantes ya definidas en
> `models.TeamContact` (`SessionTypeAtencionPsicosocial`).

---

## 6. Decisiones del líder (resueltas — Jul 2026)

**A. Primer Contacto con "Continuar = Sí" pero Consentimiento = No en S4.**
**Resuelto:** avanza igual que un Primer Contacto regular — `ya_hizo_primer_contacto = true`,
`status = en_gestion`, `ya_hizo_primera_atencion` permanece `false`. No queda en
`en_devolucion` (a diferencia del caso equivalente en Primera Atención). Implementado en
`resolvePsicosocialSessionType` / tabla de la sección 5.

**B. Constantes de `session_type` nuevas.**
**Resuelto:** se agregaron directo en `internal/models/team_contact.go`, junto a las 4
existentes: `SessionTypePrimerContactoConAtencion`, `SessionTypePrimerContactoSinConsentimiento`,
`SessionTypeContactoSinAtencion`, `SessionTypeCierreNoConsentimiento`.

**C. Alcance de esta iteración — procesamiento de Barreras.**
**Actualizado (segunda iteración, Jul 2026):** ya se implementó la creación de `barrier_v2` +
`case_task`/`entity_letter` a partir de "Identificación de Barreras" — ver sección 9. Sigue
fuera de alcance únicamente "Seguimiento a Barreras" (actualizar/cerrar barreras existentes,
equivalente al PASO 3b de `processFollowUpSubmission`), porque requiere primero decidir cómo se
determinan las "barreras activas" del caso sin un campo equivalente a
`follow_up_v2.active_barrier_ids` en `psychosocial_support`.

**D. `CONTACTO_SIN_ATENCION` y "Fecha nueva" → `psychosocial_support.scheduled_at`.**
**Resuelto** — `ScheduledAt *time.Time` ya existe en `models.PsychosocialSupport` (campo
preexistente, no uno de los agregados en esta iteración). No hay gap de columna; solo falta
escribir el `UPDATE` en el PASO 6.

**E. Redirect tras completar.** `registrar_sesion.html` (`onFormCompleted` → `goToRemision`) ya
redirige a `/salvia/remision-psicosocial/{psicosocialId}` — confirmado, no requiere cambios.

---

## 7. Cambios implementados

| # | Archivo | Cambio |
|---|---|---|
| 1 | `internal/repository/team_contact_repository.go` | **Nuevo.** `TeamContactRepository` (genérico + `FindByFormSubmissionID`) |
| 2 | `internal/models/team_contact.go` | Constantes `SessionType*` nuevas agregadas |
| 3 | `salvia/service/form_service.go` | `FormServiceDeps`/`formService`: `TeamContactRepo` inyectado. Los 4 `FormID*` agregados al switch de `OnEndFormSubmission`. `processPsicosocialSessionSubmission` + `resolvePsicosocialSessionType` + `psicosocialSessionTimelineType` implementadas |
| 4 | `main.go` | `teamContactRepo := repository.NewTeamContactRepository(gormDB)` + wiring en `FormServiceDeps` |

No requirió cambios en `registrar_sesion.html` ni en `dinamic-form.js` — el evento ya estaba
conectado correctamente en el frontend (E-01). Compilado sin errores
(`go build ./internal/... ./salvia/... .`).

### Pendiente

- Procesamiento de "Seguimiento a Barreras" (actualizar/cerrar barreras existentes) — ver
  sección 9. "Identificación de Barreras" (crear `barrier_v2`/tareas/oficios) ya está
  implementado (sección 9).
- Probar el flujo end-to-end en los 4 escenarios (A–D de `flujo-psicosocial.md`) una vez el
  usuario pueda agendar sesiones y completar formularios reales.

---

## 8. Agendamiento automático desde "Fecha próxima atención" (Jul 2026, implementado)

La pregunta **"Fecha próxima atención"** se repite en una sección distinta de cada uno de los 4
formularios (y en 2 secciones mutuamente excluyentes dentro de Primer Contacto, según la
respuesta de "Continuar Primera Atención"). Cuando se responde con una fecha, se debe generar un
**nuevo `team_contact`** con la sesión agendada — ocurre en el mismo evento de guardado (E-02),
no en un evento aparte.

### IDs reales (capturados de Supabase, 2026-07-15)

| Formulario | Sección | ID | Condición de visibilidad |
|---|---|---|---|
| Primer Contacto | S1 — Primer contacto (Q13) — "Fecha nueva" | `58ce2d34-24d2-4e73-bf95-26a2c608f8e6` | Visible si Continuar Primera Atención = No |
| Primer Contacto | S4 — Primera Atención (Q11) — "Fecha próxima atención" | `fd2fb664-de83-4069-a161-6348dfef48bf` | Visible si Continuar Primera Atención = Sí |
| Primera Atención | S1 — Contacto (Q?) — "Fecha nueva" | `64d63b79-edee-464b-be56-1104efd31a46` | Visible si "¿Es atención o solo contacto?" = Solo Contacto |
| Primera Atención | S4 — Primera Atención (Q11) — "Fecha próxima atención" | `d68c7334-74bb-47b1-a2ee-f1d3da04627b` | Sección visible si "Es atención" = Atención |
| Atención Psicosocial | S1 — Contacto Atención Psicosocial (Q?) — "Fecha nueva" | `72ce49d2-f853-4f4a-9f1a-f795f4d514c3` | Visible si "¿Es atención o solo contacto?" = Solo Contacto |
| Atención Psicosocial | S4 — Atención Psicosocial (Q4) — "Fecha próxima atención" | `7040a37d-f346-4bf7-9b80-6286b9da62c5` | Sección visible si "Es atención" = Atención |
| Cierre | S1 — Contacto (Q?) — "Fecha nueva" | `1b4d09f0-e5ca-4b4e-9d74-428478a113c6` | Visible si "¿Es atención o solo contacto?" = Solo Contacto |
| Cierre | S4 — Atención Psicosocial (Cierre) (Q4) — "Fecha próxima atención" | `3ce6ff5a-f139-4417-9255-155207e9a970` | Sección visible si "Es atención" = Atención |

`extractFechaProximaAtencion(formID, answerMap)` revisa los 2 IDs candidatos de cada formulario
(mutuamente excluyentes por visibilidad — ambos dependen del mismo trigger "¿Es atención o solo
contacto?"/"Continuar Primera Atención") y retorna el primer valor no vacío.

**🐛 Bug corregido (Jul 2026):** la primera implementación solo incluía el candidato de S4 para
Primera Atención / Atención Psicosocial / Cierre — el de S1 ("Fecha nueva" en la sección
Contacto, visible cuando se responde "Solo Contacto") no estaba en la lista. Efecto: el
formulario se guardaba bien (la respuesta sí quedaba en `salvia.answer`), pero
`extractFechaProximaAtencion` retornaba `""` y nunca se llamaba a
`scheduleNextPsicosocialContact` — no se creaba el `team_contact` agendado. Reportado por el
usuario tras probar el escenario "Solo Contacto" en Atención Psicosocial (y confirmado que
aplicaba igual a Primera Atención y Cierre). Corregido agregando los 3 IDs de S1 faltantes
(`qPAFechaNuevaS1`, `qSEGFechaNuevaS1`, `qCIEFechaNuevaS1`) a sus respectivos candidatos.

### Lógica (`scheduleNextPsicosocialContact`)

```
SI sessionType == CIERRE:
    → NO agendar (la remisión se está cerrando en esta misma sesión)
SINO SI extractFechaProximaAtencion(formID, answerMap) != "":
    → Parsear fecha (formato "2006-01-02")
    → Crear NUEVO team_contact:
         case_id          = ps.CaseID
         psicosocial_id   = ps.ID
         is_psico_session = true
         is_completed     = false
         status           = "agendada"
         scheduled_date   = fecha parseada
         professional_id / dupla_id = el que tenga ps (nunca ambos — mismo criterio que
                                       applyContactAssignment en psychosocial_detail_service.go)
         form_id          = NULL  ← se resuelve en su propio evento E-01, no ahora
         session_type     = NULL  ← idem
```

**Por qué `form_id`/`session_type` quedan vacíos:** a diferencia del `team_contact` que se
acaba de completar (cuyo `form_id` se fijó en su propio E-01 y ya no cambia), el contacto nuevo
todavía no tiene sesión — su formulario correcto depende del estado de `psychosocial_support`
**en el momento en que se abra esa sesión futura** (que puede ser distinto del estado actual: p.
ej. `session_count` ya habrá incrementado en este mismo guardado). Se resuelve de forma perezosa,
igual que hace `LoadSession` con un `team_contact` heredado sin `form_id`.

**Idempotencia:** si el submission se vuelve a guardar (edición posterior), el chequeo de
`tc.IsCompleted == true` al inicio de `processPsicosocialSessionSubmission` corta la ejecución
antes de llegar a esta sección — no se duplica el contacto agendado.

**Supuesto a validar con el líder:** por ahora solo se bloquea el agendamiento cuando
`sessionType == CIERRE` (cierre explícito). Los casos `CIERRE_NO_CONSENTIMIENTO` (Primera
Atención sin consentimiento) y `PRIMER_CONTACTO_SIN_CONSENTIMIENTO` sí agendan la próxima sesión
si el formulario trae una fecha — se decidió así porque en ambos casos el proceso continúa
(`status` no queda en `cerrado`), pero es un supuesto, no una decisión explícita del líder.

---

## 9. Procesamiento de Barreras — "Identificación de Barreras" (Jul 2026, implementado)

Segunda iteración anunciada en el punto C de la sección 6. Reemplaza el estado "fuera de
alcance" — ya se crea `barrier_v2` (+ `case_task`/`entity_letter` según gestión) por cada
entrada agregada al repeater **"Identificación de Barreras"**, en cualquiera de los 4
formularios psicosociales. Replica 1:1 el PASO 3 de `processFollowUpSubmission`
(`hacer_seguimiento`), solo que relacionado con la remisión psicosocial en vez del follow-up.

**Actualizado (Jul 2026, implementado — ver sección 9.5):** la sección **"Seguimiento a
Barreras"** (mostrar las barreras activas de la remisión para que el agente les dé seguimiento)
ya resuelve de dónde salen esas "barreras activas": en vez de replicar el campo
`follow_up_v2.active_barrier_ids` (CSV de IDs guardado en el propio follow_up), se decidió con
los líderes agregar `barrier_v2.team_contact_id` y resolver las barreras activas de una remisión
buscando todos los `team_contact` de ese `psychosocial_support`. La actualización/cierre de esas
barreras (equivalente al PASO 3b de `processFollowUpSubmission`) sigue pendiente — ver 9.6.

### 9.1 Frontend — cascada de ubicación (Departamento → Ciudad → Municipio)

Las preguntas "Departamento/Ciudad/Municipio donde se presentó la barrera" ya tenían
`state_options_path` seteado desde el seed original (`statesColombia` /
`newBarriers.{_entryIndex}.cities` / `newBarriers.{_entryIndex}.towns` — igual patrón que
`hacer_seguimiento.html`), pero `registrar_sesion.html` no tenía la lógica que puebla esos
paths en `formState`. Se agregó, replicando exactamente `loadLocations()` /
`_updateBarreraLocationOptions()` de `hacer_seguimiento.html`:

- `mounted()` ahora llama `loadLocations()` (fetch de `/api/v1/locations/departments` y
  `/api/v1/locations/cities`), en paralelo con `loadPsicosocial()`.
- `loadPsicosocial()` ahora *mergea* `formState` en vez de reemplazarlo, para no perder lo que
  `loadLocations()` ya haya escrito (las dos llamadas se disparan juntas, sin orden garantizado).
- El componente escucha `@answers-updated` del `dinamic-form` (antes no estaba conectado) →
  `onAnswersUpdated(payload)` → `_updateBarreraLocationOptions(payload)`.
- Como la sección "Identificación de Barreras" está clonada en cada uno de los 4 formularios
  (repeater group ID y preguntas Departamento/Ciudad distintos por formulario),
  `_updateBarreraLocationOptions` resuelve la config correcta con un mapa cliente
  `PSICO_BARRIER_LOCATION_CONFIG[this.formId] → { groupId, deptQ, cityQ }`.

### 9.2 Backend — `processPsicosocialBarrierEntries` (nuevo, `psicosocial_barreras.go`)

Se agregó un archivo nuevo `salvia/service/psicosocial_barreras.go` con:

- `psicosocialBarrierQuestions`: struct con el repeater group ID + los IDs de las 22 preguntas
  de "Identificación de Barreras".
- `psicosocialBarrierQuestionsByForm`: mapa `formID → psicosocialBarrierQuestions` con los 4
  juegos de IDs reales (capturados de Supabase, ver tabla abajo).
- `processPsicosocialBarrierEntries(ctx, formID, submissionID, actorID, ps)`: por cada entry del
  repeater, crea un `BarrierV2` y, por cada opción marcada en "Gestión de la barrera", crea el
  `case_task` (y `entity_letter` si la gestión lo requiere) correspondiente — misma lógica que
  `processFollowUpSubmission` (labels de gestión, `gestionSinTarea`, `gestionConOficio`,
  `taskType`), con dos diferencias:
  - `BarrierV2.FollowUpID` / `CaseTask.FollowUpID` usan `ps.FollowUpID` (el follow-up de origen
    de la remisión psicosocial) en vez de un follow-up "actual" — mantiene la barrera
    relacionada con el historial del caso.
  - `CaseTask.PsychosocialSupportID` se fija en `ps.ID` (campo ya existente en el modelo,
    pensado justo para este caso) para poder filtrar tareas por remisión psicosocial.
- Se llama desde `processPsicosocialSessionSubmission`, justo después de agendar la próxima
  sesión (sección 8) y antes de crear el evento de timeline. Si falla, solo se loguea como
  advertencia — no aborta el resto del guardado (igual criterio que otras sub-tareas
  best-effort de la función). El conteo de barreras creadas se agrega a la descripción del
  evento de timeline si es mayor a 0.

### 9.3 IDs reales (capturados de Supabase, 2026-07-15) — repeater "Identificación de Barreras"

| Formulario | Repeater Group ID |
|---|---|
| Primer Contacto Psicosocial | `4e9a169b-76d7-4f4d-ab9d-c807cff457d5` |
| Primera Atención Psicosocial | `b7fe655d-7e4e-4f67-b997-a1bf7512210e` |
| Atención Psicosocial | `0fbc7da0-8349-4bcd-acd2-da3c5240b4cd` |
| Cierre Psicosocial | `cb9b49e5-91f0-430e-9c0f-dc0235ef396f` |

Las 22 preguntas de cada repeater (Sector, Salud/Justicia/Protección + instituciones + "otra
barrera", Nombre institución, Departamento/Ciudad/Municipio, 4 barreras estructurales, Fecha,
Funcionario, Descripción, Gestión) están documentadas con sus IDs completos en
`psicosocialBarrierQuestionsByForm` (`salvia/service/psicosocial_barreras.go`) — no se
duplican aquí para evitar que la tabla y el código queden desincronizados.

**Valores de opciones verificados contra la referencia (`hacer_seguimiento`):** "Sector de la
barrera" y "Gestión de la barrera" usan los mismos `value` (`salud`/`justicia`/`proteccion` +
`otras_instituciones`/`barrera_transversal` extra sin sub-preguntas propias;
`orientacion_llamada`/`gestion_llamada`/`activacion_ruta_interinstitucional`/
`articulacion_institucional`/`escalamiento_organismo_control`/`alerta_barreras`) — la lógica de
`gestionLabels`/`gestionConOficio`/`gestionSinTarea` reutilizada sin cambios sigue siendo válida.

### 9.4 Cambios implementados (resumen)

| # | Archivo | Cambio |
|---|---|---|
| 1 | `salvia/service/psicosocial_barreras.go` | **Nuevo.** Config por formulario + `processPsicosocialBarrierEntries` |
| 2 | `salvia/service/form_service.go` | Llamada a `processPsicosocialBarrierEntries` dentro de `processPsicosocialSessionSubmission`; descripción del timeline incluye conteo de barreras |
| 3 | `frontend/html/salvia/psicosocial/registrar_sesion.html` | `@answers-updated` conectado; `loadLocations()` + `onAnswersUpdated`/`_updateBarreraLocationOptions` (config `PSICO_BARRIER_LOCATION_CONFIG` por formId); `formState` inicial con `statesColombia`/`newBarriers`; merge de `formState` en `loadPsicosocial()` |

No requirió migraciones — `state_options_path` de Departamento/Ciudad/Municipio ya estaba seteado
correctamente desde el seed original en los 4 formularios (verificado contra Supabase antes de
implementar). Compilado sin errores (`go build ./internal/... ./salvia/... .`).

---

## 9.5 "Seguimiento a Barreras" — resolver las barreras activas vía `team_contact_id` (Jul 2026, implementado)

**Decisión de los líderes:** en vez de replicar `follow_up_v2.active_barrier_ids` (un CSV de IDs
guardado en el propio contenedor padre — follow_up en `hacer_seguimiento`), se agrega una
relación directa `barrier_v2.team_contact_id` → `team_contact.id`. Cada barrera creada por
`processPsicosocialBarrierEntries` (sección 9.2) queda marcada con el `team_contact` de la
sesión donde se identificó. Al cargar la pantalla, se buscan **todos** los `team_contact` de ese
`psychosocial_support` y luego las `barrier_v2` (no `MANAGED`) cuyo `team_contact_id` esté en ese
conjunto — son las "barreras activas" de la remisión, mostradas en el repeater "Seguimiento a
Barreras Activas" vía `stateItems: currentBarriers`.

**Por qué este enfoque y no el de `hacer_seguimiento`:** `active_barrier_ids` funciona bien
cuando hay un solo contenedor "vivo" por caso (el follow_up actual). En la remisión psicosocial
el equivalente sería `psychosocial_support`, pero agregarle un campo mutable de ese tipo implica
recalcularlo/reescribirlo en cada sesión (más estado a sincronizar). Con `team_contact_id` la
relación queda fija en el momento de creación de la barrera y no requiere ningún campo adicional
en `psychosocial_support` — se resuelve siempre "hacia atrás" desde `team_contact`.

### Cambios

**1. Modelo — `internal/models/barrier_v2.go`**

```go
TeamContactID *string `gorm:"type:varchar(36);column:team_contact_id;index" json:"teamContactId,omitempty"`
```

Nullable — las barreras creadas por `hacer_seguimiento` (o remisiones anteriores a este cambio)
quedan con `team_contact_id = NULL`; solo importa para las creadas desde los formularios
psicosociales.

**2. `salvia/service/psicosocial_barreras.go` — `processPsicosocialBarrierEntries`**

Firma ahora recibe `teamContactID string` (el `tc.ID` que se está completando en
`processPsicosocialSessionSubmission`) y lo fija en cada `BarrierV2` creado.

**3. Repositorios nuevos**

- `TeamContactRepository.FindByPsicosocialID(ctx, psicosocialID) ([]models.TeamContact, error)`
  — todos los `team_contact` (completados o no) de una remisión.
- `BarrierV2Repository.FindActiveByTeamContactIDs(ctx, teamContactIDs) ([]models.BarrierV2, error)`
  — filtra `team_contact_id IN (?) AND status != 'MANAGED'`.

Ninguno de los dos se usa todavía desde fuera de `psychosocial_detail_service.go` (que resuelve
esto con una query directa vía `s.db`, igual estilo que el resto del archivo — ver punto 4), pero
quedan expuestos en la interfaz del repositorio para quien necesite esta consulta desde otro
service.

**4. `salvia/service/psychosocial_detail_service.go` — `LoadSession` (evento E-01)**

Se agregó `loadActivePsicosocialBarriers(ctx, psicosocialID) []ActiveBarrierInfo` (reutiliza el
struct `ActiveBarrierInfo`/`buildBarrierName` ya definidos en `followup_v2_service.go`, mismo
paquete `service`) y se puebla `FormState["currentBarriers"]` en la respuesta de `LoadSession`.
Antes `FormState` se devolvía siempre vacío (`map[string]interface{}{}`).

**5. Frontend — sin cambios en `registrar_sesion.html`**

No hizo falta tocar el frontend: `formState.currentBarriers` llega directo del backend en la
respuesta de `/psychosocial-support/:id/load` y ya se mergea correctamente gracias al cambio de
la sección 9.1 (`this.formState = { ...this.formState, ...(data.formState || {}) }`).
`dinamic-form.js` lee `currentBarriers` vía `stateItems` del repeater (ver punto 6) sin lógica
adicional del componente `registrar_sesion.html`.

**6. Configuración de formulario (datos, no código) — `state_items` + `render_modification`**

Los 4 repeaters "Seguimiento a Barreras Activas" no tenían `state_items` seteado (a diferencia
del de `hacer_seguimiento`, que sí lo tenía desde antes) ni el `render_modification` que muestra
el nombre de la barrera en la pregunta info "Seguimiento a Barrera". Se ejecutó contra Supabase:

```sql
UPDATE salvia.repeater_group SET state_items = 'currentBarriers' WHERE id IN (
  'cb2fb8b8-0431-4088-b626-06e1069432fa', -- Atención Psicosocial
  '10122569-e6c2-4449-a7ff-add378815eb8', -- Cierre
  '94718173-b809-41e4-86d5-378de3434dca', -- Primer Contacto
  '3cb23f26-af5a-4052-8987-28294f68e5dc'  -- Primera Atención
);

INSERT INTO salvia.render_modification (target_type, target_id, target_field, modification_type, state_path)
VALUES
  ('question', '752865b1-1dbd-4111-8ea1-1425c3997bba', 'description', 'SET', 'currentBarriers.{_entryIndex}.barrierName'), -- Atención Psicosocial
  ('question', 'c0d96633-901f-49f2-bc56-337f05656f41', 'description', 'SET', 'currentBarriers.{_entryIndex}.barrierName'), -- Cierre
  ('question', '8fa34c29-652c-49e2-ad44-d84aa3b74add', 'description', 'SET', 'currentBarriers.{_entryIndex}.barrierName'), -- Primer Contacto
  ('question', '1e65b5f2-41b2-46b6-937b-e5a37fc4ba8a', 'description', 'SET', 'currentBarriers.{_entryIndex}.barrierName'); -- Primera Atención
```

Ambos ya quedaron aplicados en la base (ejecutados contra Supabase con el mismo mecanismo que se
usó para migrar el seed).

### ⚠️ Pendiente manual — columna `barrier_v2.team_contact_id`

Igual que pasó con `team_contact.form_id`/`session_type` (ver sección 2), el usuario de base de
datos de la app **no tiene permisos de `ALTER TABLE`** sobre `barrier_v2` (`ERROR: must be owner
of table barrier_v2`, SQLSTATE 42501). Se agregó el `ADD COLUMN IF NOT EXISTS` defensivo en
`main.go` (junto a los de `victim_case`) y `models.BarrierV2` ya está en la lista de
`AutoMigrate`, pero ambos van a fallar en silencio (log `[WARN]`) contra Supabase con el usuario
actual. **Falta ejecutar manualmente, con un usuario con permisos de owner:**

```sql
ALTER TABLE salvia.barrier_v2 ADD COLUMN IF NOT EXISTS team_contact_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_barrier_v2_team_contact_id ON salvia.barrier_v2(team_contact_id);
```

Hasta que se ejecute, `processPsicosocialBarrierEntries` va a fallar al crear cualquier barrera
(columna inexistente) — mismo síntoma que el error `column "form_id" of relation "team_contact"
does not exist` que ya se vio antes con `team_contact`.

## 9.6 Pendiente — actualizar/cerrar barreras desde "Seguimiento a Barreras"

Mostrar las barreras activas (9.5) resuelve la mitad del flujo. Falta el equivalente al PASO 3b
de `processFollowUpSubmission`: cuando el agente responde el repeater "Seguimiento a Barreras"
(¿persiste?, gestión, ¿se cierra?, motivo de cierre), hay que:

- Relacionar cada entry del repeater con su `BarrierV2.ID` — como ahora `currentBarriers` ya
  trae `{id, barrierName}` en el mismo orden en que se renderizan las entries (`stateItems`), la
  relación por posición es directa (mismo mecanismo que `activeIDs[idx]` en
  `processFollowUpSubmission`, pero ahora la fuente es `currentBarriers[idx].id` en vez de
  `fu.ActiveBarrierIDs`).
- Si `¿Se realiza cierre de la barrera? = true` → `BarrierV2Repository.UpdateStatus(ctx, id, MANAGED)`.
- Crear un `BarrierFollowUp` (¿persiste?, respuesta institucional, gestión, actuaciones, cierre,
  motivo) — el modelo ya existe (`internal/models/barrier_follow_up.go`), solo falta usarlo desde
  el flujo psicosocial (usaría `FollowUpID: ps.FollowUpID`, igual criterio que 9.2).
- Registrar un evento de timeline por cada seguimiento a barrera (igual que
  `buildBarrierFollowUpSummary` en `processFollowUpSubmission`).

No implementado todavía — queda para la siguiente iteración una vez validado el flujo de 9.5 en
producción.

## 10. Bug corregido — "Continuar Primera Atención = Sí" no completaba la sesión (Jul 2026)

**Síntoma reportado:** en el formulario "Primer Contacto", al marcar "Continuar Primera Atención" =
Sí (se habilita la Sección 4 "Primera Atención" dentro del mismo formulario) y responder esa
sección, al presionar "Guardar" las respuestas quedaban persistidas pero:
- `team_contact` no se marcaba como completo... en realidad sí se marcaba, pero con el
  `sessionType` incorrecto.
- `psychosocial_support.ya_hizo_primera_atencion` seguía en `false` y `session_count` no subía.
- El nuevo `team_contact` agendado con la "Fecha próxima atención" de la Sección 4 no se creaba
  en los casos donde el flujo se interrumpía antes de llegar a esa parte del código.

**Causa raíz:** la pregunta "PC-A-Q2" (Consentimiento Informado, primera pregunta de la Sección 4
"Primera Atención" dentro de Primer Contacto, id `c3296c87-d7e3-4ea1-8a0f-cbf77c561d0f`) se había
seedeado con `question_type = 'info'` (banner informativo, sin control de respuesta) en vez de
`'single'` (Sí/No) — a diferencia de su equivalente exacto en el formulario "Primera Atención"
(`7256b91e-861b-48b3-9216-2ac13f0ae889`), que sí quedó correctamente como `'single'`.

Efecto: al ser tipo `info`, ni el frontend (`dinamic-form.js` → `validateCurrentSection`) ni el
backend (`validateAnswer` → `case "info": return validateAnswerOK`) exigían ni permitían responder
Sí/No a esa pregunta — el agente solo veía el texto del consentimiento, sin ningún control
interactivo. Como consecuencia, `answerMap[qPCConsentimiento]` llegaba siempre vacío a
`resolvePsicosocialSessionType`, que interpretaba esto como "no hay consentimiento" y devolvía
`SessionTypePrimerContactoSinConsentimiento` en lugar de `SessionTypePrimerContactoConAtencion`:

```go
if answerMap[qPCConsentimiento] == "si" {
    return models.SessionTypePrimerContactoConAtencion, ""
}
return models.SessionTypePrimerContactoSinConsentimiento, ""
```

Con `SessionTypePrimerContactoSinConsentimiento`, `processPsicosocialSessionSubmission` solo
marca `ya_hizo_primer_contacto = true` y no toca `ya_hizo_primera_atencion` ni `session_count`
(ver switch en la sección 5) — exactamente el síntoma reportado.

**Corrección aplicada:**
1. `src/cmd/seed/seed_psicosocial.sql` (bloque `PC-A-Q2`): `question_type` cambiado de `'info'` a
   `'single'` (ya tenía las opciones Sí/No creadas correctamente).
2. Corrección aplicada también directamente en la BD viva (Supabase) sobre la pregunta existente
   `c3296c87-d7e3-4ea1-8a0f-cbf77c561d0f`, ya que el seed no se re-ejecuta automáticamente.

No se requirió ningún cambio de código Go — la lógica de `resolvePsicosocialSessionType` y
`processPsicosocialSessionSubmission` ya era correcta; el problema era exclusivamente de datos
(seed). El agendamiento automático de "Fecha próxima atención" (sección 8) ya contemplaba los IDs
de S1 y S4 de Primer Contacto desde el principio, así que una vez corregido el `sessionType`,
la creación del nuevo `team_contact` agendado funciona sin cambios adicionales.
