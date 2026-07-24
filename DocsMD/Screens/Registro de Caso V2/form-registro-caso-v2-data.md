# Formulario "Registro de Caso" — Estructura de Secciones y Preguntas

**Form ID:** `0a24ab30-3cfc-4861-b74d-65d21524bc00`
**Schema DB:** `salvia`
**Tablas involucradas:** `form`, `form_section`, `question`, `option`, `visibility_condition`, `render_modification`
**Script de creación:** `src/cmd/seed/seed_registro_caso.sql`

> Migración del formulario estático `set_victim_case.html` a `dinamic-form`. Excluye la grilla de "Asignación de Ruta" (sedes por momento/sector) — ver GAPS. Las 9 preguntas de discapacidad se modelan como `single` con una escala fija de 4 opciones en vez de star-rating.

---

## Enums compartidos (`salvia.victim_case_form2_enums`)

Todas las preguntas `dropdown`/`multiple` marcadas con `state_options_path = enums.<categoría>` **no tienen `option` propias en BD** — reutilizan el catálogo compartido `victim_case_form2_enums` (mismo usado hoy por `set_victim_case.html`, `set_victim_contact.html` y feminicidio), agrupado por columna `victim_case_form2_enums_category`. La pantalla anfitriona debe cargar este catálogo (o las categorías relevantes) y escribirlo en `formState.enums.<categoría>` como `[{label, value}]`.

> ✅ **Resuelto (revisión de código):**
> - `docType2` **no** vive en `victim_case_form2_enums` — es el mapa Go hardcodeado `common_config.DOCUMENT_TYPE_FORM2` (`common/config/Enums.go`, ~20 valores). Se materializó como 20 `option` rows estáticas directamente sobre la pregunta "Tipo de documento" (sin `state_options_path`).
> - `aggressor_occupation`: el código actual (`VictimCaseFacade.go:687,866`) en realidad puebla ese campo con la categoría **`victim_case_form2_relationship_with_presumed_aggressor_02`** — probablemente un nombre mal puesto en el legacy (el campo se llama "ocupación" pero muestra opciones de relación con el agresor). Se replicó el comportamiento **real** (no lo que sugiere el nombre) — `state_options_path = 'enums.victim_case_form2_relationship_with_presumed_aggressor_02'`. Confirmar con el equipo si esto es un bug a corregir antes de lanzar la V2.

---

## Secciones

| Order | ID | Nombre | Preguntas |
|---|---|---|---|
| 1 | `eb2782e7-3888-4f2f-a6f2-89ed791cbd14` | Datos de la Víctima | 16 |
| 2 | `02ee584e-16b6-4d69-a6a5-678cc3564e78` | Contacto de Apoyo | 4 |
| 3 | `ba4a4edc-6dc6-40f7-9f4f-3749a8397986` | Hechos | 18 |
| 4 | `71eb735c-d199-4d37-8b32-e2ad06195e47` | Agresor | 11 |
| 5 | `0f5bbfb1-474c-4260-b87e-251bc3c1eb32` | Tamizaje | 39 |
| 6 | `39785f1c-2de7-48e9-966e-9bf76b47b3bd` | Datos Personales | 34 |
| 7 | `b11b2408-6b06-4a9b-bb85-336f6297d5df` | Plan de Acción | 2 |
| 8 | `48a3f928-ccc5-4295-9304-b0285f942c3a` | Denuncia Fácil | 1 |
| 9 | `7d64c1ea-1da7-40b1-9f7e-ed401c8a019b` | Lugar de Atención | 3 |
| **Total** | | | **128** |

> **No existe sección de "Autorización de datos personales"** — es un gate a nivel de pantalla anfitriona (igual que `canEdit`/`loading`), se resuelve antes de montar `<dinamic-form>`.

---

## Sección 1 — Datos de la Víctima

| Order | Tipo | Req | Pregunta | `state_options_path` / Notas |
|---|---|---|---|---|
| 1 | `text` | ✅ | Nombres | |
| 2 | `text` | ✅ | Apellidos | |
| 3 | `text` | ❌ | Nombre identitario | |
| 4 | `number` | ✅ | Teléfono | |
| 5 | `dropdown` | ✅ | Tipo de documento | 20 `option` estáticas (mapa Go `DOCUMENT_TYPE_FORM2`, no `state_options_path`) |
| 6 | `text` | ✅ | Número de documento | |
| 7 | `dropdown` | ✅ | Zona de residencia | `enums.victim_case_form2_facts_zone` |
| 8 | `dropdown` | ✅ | Departamento de residencia | `departments` |
| 9 | `dropdown` | ✅ | Ciudad de residencia | `geo.residencia.cities` |
| 10 | `dropdown` | ✅ | Municipio de residencia | `geo.residencia.towns` |
| 11 | `text` | ✅ | Dirección de residencia | |
| 12 | `info` | ❌ | "Accesibilidad" (banner divisor) | |
| 13 | `boolean` | ✅ | ¿Tiene alguna discapacidad? | |
| 14 | `multiple` | ✅ | Ajustes razonables VBG | `enums.victim_case_form2_adjustments_gbv` · VC: Q13=true |
| 15 | `boolean` | ✅ | ¿Requiere intérprete de idioma? | |
| 16 | `text` | ❌ | ¿Cuál idioma? | VC: Q15=true |

## Sección 2 — Contacto de Apoyo

| Order | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|
| 1 | `text` | ❌ | Nombres del contacto de apoyo | |
| 2 | `number` | ❌ | Teléfono del contacto de apoyo | |
| 3 | `text` | ❌ | Correo del contacto de apoyo | |
| 4 | `dropdown` | ❌ | Parentesco del contacto de apoyo | `enums.victim_case_form2_support_contact_kinship` |

## Sección 3 — Hechos

| Order | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|
| 1 | `text` | ✅ | Relato de los hechos | |
| 2 | `date` | ✅ | Fecha de los hechos | |
| 3 | `datetime` | ✅ | Hora de inicio de los hechos | adaptado desde time-only |
| 4 | `dropdown` | ✅ | Zona de los hechos | `enums.victim_case_form2_facts_zone` |
| 5 | `dropdown` | ✅ | Departamento de los hechos | `departments` |
| 6 | `dropdown` | ✅ | Ciudad de los hechos | `geo.hechos.cities` |
| 7 | `dropdown` | ✅ | Municipio de los hechos | `geo.hechos.towns` |
| 8 | `text` | ✅ | Dirección de los hechos | |
| 9 | `dropdown` | ✅ | Escenario de la violencia | `enums.victim_case_form2_scenario_violence` |
| 10 | `multiple` | ✅ | Tipo de violencia experimentada | `enums.victim_case_form2_type_violence_experienced` |
| 11 | `multiple` | ✅ | Subtipo de violencia experimentada | `violenceSubtypes` (dinámico) · VC: `hasSelectedViolenceTypes`=true |
| 12 | `multiple` | ✅ | Ámbito de la violencia | `enums.victim_case_form2_scope_of_violence` |
| 13 | `dropdown` | ✅ | Sector laboral de ocurrencia | `enums.victim_case_form2_workplace_sector_occurrence` · VC: `hasWorkplaceScope`=true |
| 14 | `boolean` | ✅ | ¿Motivada por género? | |
| 15 | `boolean` | ✅ | ¿Denunció previamente? | |
| 16 | `multiple` | ✅ | ¿A quién denunció? | `enums.victim_case_form2_who_report_to` · VC: Q15=true |
| 17 | `boolean` | ✅ | ¿Atención apropiada? | VC: Q15=true |
| 18 | `dropdown` | ✅ | Recurrencia de la agresión | `enums.victim_case_form2_recurrence_aggression` |

## Sección 4 — Agresor

| Order | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|
| 1 | `dropdown` | ✅ | Número de agresores | `enums.victim_case_form2_num_agressors` |
| 2 | `dropdown` | ✅ | Proximidad con el agresor principal | `enums.victim_case_form2_proximity_principal_aggressor` |
| 3 | `dropdown` | ✅ | Relación con el presunto agresor | `enums.victim_case_form2_relationship_with_presumed_aggressor_01` |
| 4 | `dropdown` | ✅ | Ocupación del agresor | `enums.victim_case_form2_relationship_with_presumed_aggressor_02` (ver nota "Resuelto" arriba — posible bug de nombre en el legacy) |
| 5 | `boolean` | ✅ | ¿Depende económicamente? | VC: `partnerKnown`=true |
| 6 | `dropdown` | ✅ | Identidad de género del agresor | `enums.victim_case_form2_aggressor_gender_identity` |
| 7 | `text` | ❌ | Nombres del agresor | |
| 8 | `dropdown` | ❌ | Tipo de documento del agresor | `enums.victim_case_form2_victim_doc_type` |
| 9 | `text` | ❌ | Número de documento del agresor | |
| 10 | `text` | ❌ | Dirección del agresor | |
| 11 | `number` | ❌ | Teléfono del agresor | |

## Sección 5 — Tamizaje

39 preguntas `boolean` (salvo la última, `info`): 6 comunes (order 1-6, siempre visibles) + 18 de pareja íntima (order 7-24, VC: `wasPartner`=true) + 14 de no-pareja (order 25-38, VC: `wasPartner`=false) + 1 banner de riesgo (order 39, `info`, VC: `hasRisk`=true, `render_modification` REPLACE sobre `riskBadgeText`).

Texto exacto de cada pregunta, fórmula de puntaje y tabla de umbrales (pareja 0-4/5-8/9-15/16-24, no-pareja 0-2/3-5/6-8/9-20): ver [`registro-caso-interface.md`](../Registro%20de%20Caso/registro-caso-interface.md) sección "Lógica del tamizaje" — la fórmula **no cambia**, solo se recalcula ahora en `formState` vía `answers-updated` (ver [flow-E02](Flujos/flow-E02-cuando-se-actualizan-respuestas.md)).

## Sección 6 — Datos Personales

| Order | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|
| 1 | `date` | ✅ | Fecha de nacimiento | |
| 2 | `boolean` | ✅ | ¿Tiene dificultades físicas/mentales/sensoriales? | |
| 3-11 | `single` | ✅ | 9 preguntas de dificultad (oír, hablar, ver, moverse, tomar, entender, comer, interactuar, cotidianas) | 4 opciones fijas · VC: Q2=true |
| 12 | `multiple` | ✅ | Ley 1996 | `enums.victim_case_form2_law_1996` · VC: Q2=true |
| 13 | `dropdown` | ✅ | Nacionalidad | `enums.victim_case_form2_nationality` |
| 14 | `dropdown` | ❌ | Nacionalidad específica | `enums.victim_case_form2_specified_nationality` · VC: Q13='ex' |
| 15 | `dropdown` | ❌ | Condición migratoria | `enums.victim_case_form2_migration_condition` · VC: Q13='ex' |
| 16 | `dropdown` | ✅ | Identidad de género | `enums.victim_case_form2_gender_identity` |
| 17 | `dropdown` | ✅ | Orientación sexual | `enums.victim_case_form2_sexual_orientation` |
| 18 | `dropdown` | ✅ | Sexo asignado al nacer | `enums.victim_case_form2_assigned_sex_at_birth` |
| 19 | `multiple` | ✅ | Población especialmente protegida | `enums.victim_case_form2_specially_protected_population` |
| 20 | `dropdown` | ✅ | Afiliación étnica | `enums.victim_case_form2_ethnic_affiliation` |
| 21 | `dropdown` | ❌ | Pueblo indígena | `enums.victim_case_form2_indigenous_people` · VC: Q20='in' |
| 22 | `boolean` | ✅ | ¿Reconocimiento como campesina? | |
| 23 | `dropdown` | ✅ | Estado civil | `enums.victim_case_form2_marital_status` |
| 24 | `dropdown` | ✅ | Último nivel educativo | `enums.victim_case_form2_last_education_level` |
| 25 | `dropdown` | ✅ | Ocupación | `enums.victim_case_form2_occupation` |
| 26 | `dropdown` | ✅ | Forma de generación de ingresos | `enums.victim_case_form2_income_generation_method` |
| 27 | `multiple` | ❌ | Modalidad ASP | `enums.victim_case_form2_asp_mode` · VC: Q26='pr' |
| 28 | `dropdown` | ❌ | Relación laboral | `enums.victim_case_form2_employment_relationship` · VC: Q26='em' |
| 29 | `text` | ❌ | Inicio aproximado ASP | VC: Q26='pr' |
| 30 | `multiple` | ❌ | Razón ASP | `enums.victim_case_form2_reason_asp` · VC: Q26='pr' |
| 31 | `dropdown` | ✅ | Forma de tenencia de vivienda | `enums.victim_case_form2_housing_tenancy_form` |
| 32 | `dropdown` | ✅ | Estrato de vivienda | `enums.victim_case_form2_housing_stratum` |
| 33 | `multiple` | ✅ | Personas a cargo | `enums.victim_case_form2_has_dependents` |
| 34 | `boolean` | ✅ | ¿Actualmente en embarazo? | |

### Opciones fijas — escala de dificultad (Q3-Q11)

| Label | Value |
|---|---|
| Sin dificultad | `sin_dificultad` |
| Alguna dificultad | `alguna_dificultad` |
| Mucha dificultad | `mucha_dificultad` |
| No puede hacerlo | `no_puede` |

## Sección 7 — Plan de Acción

| Order | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|
| 1 | `multiple` | ✅ | Plan de acción | `enums.victim_case_form2_action_plan` |
| 2 | `text` | ✅ | Explicación de la gestión Salvia | |

## Sección 8 — Denuncia Fácil

| Order | Tipo | Req | Pregunta |
|---|---|---|---|
| 1 | `boolean` | ✅ | ¿Permite el uso del canal de Denuncia Fácil? |

## Sección 9 — Lugar de Atención

| Order | Tipo | Req | Pregunta | Notas |
|---|---|---|---|---|
| 1 | `dropdown` | ✅ | Departamento de atención | `departments` |
| 2 | `dropdown` | ✅ | Ciudad de atención | `geo.atencion.cities` |
| 3 | `dropdown` | ✅ | Municipio de atención | `geo.atencion.towns` |

> **La grilla de "Asignación de Ruta" (sedes por momento/sector) NO está incluida** — decisión explícita del usuario. Se mantiene fuera de `dinamic-form`; la asignación de sedes se sigue haciendo por fuera (p. ej. `Asignar Operadores`) hasta que se decida cómo incorporarla.

---

## Resumen general

| Sección | Preguntas | Visibility Conditions | Options propias |
|---|---|---|---|
| Datos de la Víctima | 16 | 2 | 0 |
| Contacto de Apoyo | 4 | 0 | 0 |
| Hechos | 18 | 4 | 0 |
| Agresor | 11 | 1 | 0 |
| Tamizaje | 39 | 33 | 0 |
| Datos Personales | 34 | 17 | 36 (9 preguntas × 4) |
| Plan de Acción | 2 | 0 | 0 |
| Denuncia Fácil | 1 | 0 | 0 |
| Lugar de Atención | 3 | 0 | 0 |
| **Total** | **128** | **57** | **36** |

Más 1 `render_modification` (banner de riesgo, Sección 5).

---

## Estrategia de creación progresiva (Borrador → Activo)

Decisión del usuario: el `victim_case` se crea desde el **primer** `saveSection` (estado `Borrador`), y se activa solo al completar el formulario. Ver detalle completo en [`flow-E03-cuando-se-guarda-primera-seccion.md`](Flujos/flow-E03-cuando-se-guarda-primera-seccion.md) y [`flow-E04-cuando-se-completa-formulario.md`](Flujos/flow-E04-cuando-se-completa-formulario.md).

**Cambio de estado necesario:** `salvia.config.VICTIM_CASE_STATUS["sp"]` (`salvia/config/Enums.go:305`) no tiene un código para "Borrador" hoy — los códigos existentes son `fc, r, ra, c, ex, is, cd, iv`. Se necesita agregar un código nuevo, propuesto **`bo`** ("Borrador"), y al completar el formulario transicionar a **`ra`** ("enrutado aprobado" — el mismo estado que `SetVictimCase` ya asigna hoy al crear un caso como rol `op`, sin necesidad de inventar un código "Activo" nuevo).

⚠️ **Un `victim_case` en estado `bo` no debe aparecer en listados/KPIs/asignaciones existentes** hasta que se active — todas las consultas que hoy no filtran explícitamente por status (`cases_list_repository.go`, `KpiController.go`, `followup_repository.go`, etc.) deben revisarse para excluir `bo` por defecto, igual que ya excluyen `cd` en algunos casos.

---

## GAPS

| # | Descripción | Impacto |
|---|---|---|
| G-02 | El motor de `visibility_condition` solo soporta `EQUALS`/`CONTAINS`, no "mayor que" ni "no vacío" — por eso `hasWorkplaceScope`, `hasSelectedViolenceTypes`, `hasRisk`, `wasPartner`, `partnerKnown` se modelan como flags booleanos precalculados por el host en vez de condiciones directas sobre las respuestas | Medio — ya resuelto vía formState, pero exige que el host los calcule correctamente (ver flow-E02) |
| G-03 | Códigos exactos usados en las condiciones EQUALS (`'ex'`, `'in'`, `'pr'`, `'em'`, `'al'`) se tomaron de `registro-caso-index.md`/`registro-caso-interface.md` tal cual — deben confirmarse contra los `code` reales en `victim_case_form2_enums` antes de implementar | Medio |
| G-04 | La Asignación de Ruta queda completamente fuera de este formulario — falta decidir en qué momento/pantalla se asignará la sede cuando se retome esa fase | Bajo (fuera de alcance actual) |
| G-11 | Agregar el código `bo` a `VICTIM_CASE_STATUS` y auditar todas las consultas de listados/KPIs que no filtran explícitamente por status para excluirlo | Alto — bloquea que el caso en Borrador no se filtre correctamente |
| G-12 | Confirmar con el equipo si `aggressor_occupation` → `relationship_with_presumed_aggressor_02` es un bug del legacy a corregir, o comportamiento intencional a preservar tal cual | Medio |
