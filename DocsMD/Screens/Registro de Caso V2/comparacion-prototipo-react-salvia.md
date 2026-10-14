# Comparación: "Registro de Caso V2" (real, SOG_SALVIA-develop) vs. prototipo "Registro SALVIA" (React, repo Salvia)

> Compara el formulario **real** `dinamic-form` ya construido en este repo (128
> preguntas, 9 secciones — ver
> [registro-caso-v2-secciones.md](registro-caso-v2-secciones.md) y
> [form-registro-caso-v2-data.md](form-registro-caso-v2-data.md)) contra el
> prototipo de solo-frontend `Salvia/src/data/registroSalviaSchema.js` (71
> campos, 7 módulos CI/DP/CA/HV/AG/TZ/GS, desplegado en
> `https://salvia-mockups.vercel.app/registro-salvia`).
>
> Fuente de las columnas de BD: `internal/repository/victim_case_form_fields.go`
> (mapeo pregunta → columna) y `victim_case_form_repository.go` (cómo se
> persiste cada `Kind`). El prototipo no tiene backend — nada de lo que
> diligencia se guarda; las columnas que se le asignan abajo son **propuestas
> por analogía**, no columnas que existan hoy.

## Resumen ejecutivo

| | Cantidad |
|---|---:|
| Preguntas en V2 real | 128 (9 secciones) |
| Campos en el prototipo | 71 (7 módulos) |
| En ambas partes (con equivalencia conceptual) | ~68 |
| Solo en el prototipo (no existen en V2) | 3 |
| Solo en V2 (no están en el prototipo) | 57 |

Los 3 campos exclusivos del prototipo vienen todos del PDF de observaciones
(`correo_victima`, `requiere_remision`, `procesos_remision`) — el prototipo
no inventó nada fuera de ese encargo. Los 57 exclusivos de V2 se concentran
en 3 bloques grandes: la batería de discapacidad de 9 preguntas (Sección 6),
el bloque de Actividad Sexual Pagada — ASP (4 preguntas, Sección 6), y la
mayoría de Tamizaje (Sección 5, ~22 de 39) — ver detalle sección por sección.

**Nota de alcance:** el prototipo reorganiza el contenido en 7 módulos que no
corresponden 1:1 a las 9 secciones de V2 (p. ej. "¿Requiere intérprete?" vive
en la Sección 1 de V2 pero en el módulo **DP** del prototipo, no en **CI**).
Cuando digo "Ambas" abajo me refiero a que la *pregunta* existe en ambos
lados, sin importar en qué módulo/sección quedó.

---

## Leyenda — cómo se guarda cada `Kind` en BD (V2)

| Kind (Go) | Significado | Dónde cae |
|---|---|---|
| `text` / `int` / `date` / `time` / `timestamp` | Columna escalar plana | `victim_case_form2.<columna>` |
| `enum1` | Un solo valor de catálogo (dropdown) | FK en `victim_case_form2.<columna>` → `victim_case_form2_enums_id` |
| `enumN` | Varios valores (multiple) | **Sin columna propia** — tabla de relación `rel_victim_case_form2_enums_victim_case_form2` |
| `bool_enum` | Sí/No | **No es un booleano nativo** — FK → `victim_case_form2_enums` categoría `yes_no` |
| `scale` | Escala de discapacidad | Entero plano 0–3 en `victim_case_form2.<columna>` |
| *(sin entrada en el mapeo)* | Nombres, Apellidos, Tipo/Núm. documento | Columnas núcleo de `victim_case` (no de `form2`) — valores "bootstrap" que crean el caso en estado Borrador |
| *(sin entrada en el mapeo)* | Departamento/Ciudad (cascadas de UI), banners informativos | No se persisten — solo resuelven el `dropdown` de Municipio en pantalla |

---

## Sección 1 — Datos de la Víctima (16 preguntas) → módulos CI / DP del prototipo

| # | Pregunta | Dónde | Prototipo | Columna/tabla V2 |
|---|---|---|---|---|
| 1 | Nombres | Ambas | CI.nombres | `victim_case` (bootstrap, no form2) |
| 2 | Apellidos | Ambas | CI.apellidos | `victim_case` (bootstrap) |
| 3 | Nombre identitario | Ambas | CI.nombre_identitario | `victim_case_form2_identity_name` (text) |
| 4 | Teléfono | Ambas | CI.telefono | `victim_case_form2_victim_phone` (int) |
| 5 | Tipo de documento | Ambas | CI.tipo_documento | `victim_case` (bootstrap) + 20 opciones estáticas (mapa Go `DOCUMENT_TYPE_FORM2`) |
| 6 | Número de documento | Ambas | CI.numero_documento | `victim_case` (bootstrap) |
| 7 | Zona de residencia | Solo V2 | — | `victim_case_form2_residence_zone` (enum1) |
| 8 | Departamento de residencia | Ambas | CI.departamento_residencia | cascada UI, sin persistencia propia |
| 9 | Ciudad de residencia | Solo V2 | — | cascada UI, sin persistencia (el prototipo saltó directo depto→municipio) |
| 10 | Municipio de residencia | Ambas | CI.municipio_residencia | `victim_case` (bootstrap) **y** `victim_case_form2_residence_town` (text) |
| 11 | Dirección de residencia | Ambas | CI.direccion_residencia | `victim_case_form2_residence_address` (text) |
| 12 | Accesibilidad *(banner)* | Solo V2 | — | sin persistencia (divisor visual) |
| 13 | ¿Tiene alguna discapacidad? | Ambas* | DP.tiene_discapacidad | `victim_case_form2_person_with_disability` (bool_enum) |
| 14 | Ajustes razonables VBG | Solo V2 | — | sin columna propia → `rel_victim_case_form2_enums_victim_case_form2` |
| 15 | ¿Requiere intérprete de idioma? | Ambas | DP.requiere_interprete | `victim_case_form2_require_language_interpreter` (bool_enum) |
| 16 | ¿Cuál idioma? | Solo V2 | — | `victim_case_form2_language_assistance` (text) |

\* La 13 cambia de módulo (CI en V2 → DP en el prototipo) y de redacción
("tiene dificultades" vs. "se reconoce con discapacidad"), pero es la misma
pregunta-gate.

## Sección 2 — Contacto de Apoyo (4 preguntas) → módulo CA

| # | Pregunta | Dónde | Prototipo | Columna V2 |
|---|---|---|---|---|
| 1 | Nombres del contacto | Ambas | CA.contacto_nombres | `victim_case_form2_support_contact_names` (text) |
| 2 | Teléfono del contacto | Ambas | CA.contacto_telefono | `victim_case_form2_support_contact_phone` (int) |
| 3 | Correo del contacto | Solo V2 | — | `victim_case_form2_support_contact_email` (text) |
| 4 | Parentesco del contacto | Ambas | CA.contacto_parentesco | `victim_case_form2_support_contact_kinship` (enum1) |

Nota curiosa: el prototipo sí tiene un correo, pero es el de la **víctima**
(`DP.correo_victima`, exclusivo del prototipo — Observación #12 del PDF), no
el del contacto de apoyo. Son dos campos distintos, cada uno ausente en el
otro lado.

## Sección 3 — Hechos (18 preguntas) → módulo HV

| # | Pregunta | Dónde | Prototipo | Columna V2 |
|---|---|---|---|---|
| 1 | Relato de los hechos | Ambas | HV.relato_hechos | `victim_case_form2_facts_description` (text) |
| 2 | Fecha de los hechos | Ambas | HV.fecha_hechos | `victim_case_form2_facts_date` (date) |
| 3 | Hora de inicio | Ambas† | HV.hora_hechos | `victim_case_form2_facts_start_time` (time) |
| 4 | Zona de los hechos | Ambas | HV.zona_hechos | `victim_case_form2_facts_zone` (enum1) |
| 5 | Departamento de los hechos | Ambas | HV.departamento_hechos | cascada UI |
| 6 | Ciudad de los hechos | Solo V2 | — | cascada UI, sin persistencia |
| 7 | Municipio de los hechos | Ambas | HV.municipio_hechos | `victim_case_form2_facts_town_code` (text) |
| 8 | Dirección de los hechos | Solo V2 | — | `victim_case_form2_facts_address` (text) |
| 9 | Escenario de la violencia | Ambas | HV.escenario_violencia | `victim_case_form2_scenario_violence` (enum1) |
| 10 | Tipo de violencia experimentada | Ambas | HV.tipo_violencia | sin columna propia → tabla de relación (enumN) |
| 11 | Subtipo de violencia experimentada | Ambas | HV.subtipo_violencia | sin columna propia → tabla de relación (enumN) |
| 12 | Ámbito de la violencia | Ambas | HV.ambito_violencia | sin columna propia → tabla de relación (enumN) |
| 13 | Sector laboral de ocurrencia | Solo V2 | — | `victim_case_form2_workplace_sector_occurrence` (enum1) |
| 14 | ¿Motivada por género? | Ambas | HV.motivada_genero | `victim_case_form2_violence_motivated_by_gender` (bool_enum) |
| 15 | ¿Denunció previamente? | Ambas | HV.denuncio_previamente | `victim_case_form2_reported_previously` (bool_enum) |
| 16 | ¿A quién denunció? | Solo V2 | — | sin columna propia → tabla de relación (enumN) |
| 17 | ¿Atención apropiada? | Solo V2 | — | `victim_case_form2_attention_was_appropriate` (bool_enum) |
| 18 | Recurrencia de la agresión | Ambas | HV.recurrencia_agresion | `victim_case_form2_recurrence_aggression` (enum1) |

† En V2 es un `datetime`; en el prototipo es un texto libre ("Ej: 19:30") —
mismo concepto, widget más simple.

## Sección 4 — Agresor (11 preguntas) → módulo AG

| # | Pregunta | Dónde | Prototipo | Columna V2 |
|---|---|---|---|---|
| 1 | Número de agresores | Ambas | AG.numero_agresores | `victim_case_form2_num_agressors` (enum1) |
| 2 | Proximidad con el agresor principal | Ambas | AG.proximidad_agresor | `victim_case_form2_proximity_principal_aggressor` (enum1) |
| 3 | Relación con el presunto agresor | Ambas | AG.relacion_agresor | `victim_case_form2_relationship_with_presumed_aggressor` (enum1) |
| 4 | Ocupación del agresor | Ambas‡ | AG.ocupacion_agresor | `victim_case_form2_aggressors_occupation` (enum1) |
| 5 | ¿Depende económicamente? | Solo V2 | — | `victim_case_form2_economically_dependent` (bool_enum) |
| 6 | Identidad de género del agresor | Ambas | AG.identidad_genero_agresor | `victim_case_form2_aggressor_gender_identity` (enum1) |
| 7 | Nombres del agresor | Ambas | AG.nombres_agresor | `victim_case_form2_aggressor_names` (text) |
| 8 | Tipo de documento del agresor | Solo V2 | — | `victim_case_form2_aggressor_doc_type` (enum1) |
| 9 | Número de documento del agresor | Solo V2 | — | `victim_case_form2_aggressor_doc_number` (text) |
| 10 | Dirección del agresor | Solo V2 | — | `victim_case_form2_aggressor_address` (text) |
| 11 | Teléfono del agresor | Ambas | AG.telefono_agresor | `victim_case_form2_aggressor_phone` (int) |

‡ Ojo con esta: `form-registro-caso-v2-data.md` documenta que en el código
real (`VictimCaseFacade.go:687,866`) esta pregunta **no** usa sus propias
opciones — puebla el campo con el catálogo
`victim_case_form2_relationship_with_presumed_aggressor_02` (ver GAP G-12,
posible bug de nombre heredado del legacy). El prototipo usa
`OCUPACION_AGRESOR_OPCIONES` propia y limpia — no hereda esa rareza.

## Sección 5 — Tamizaje (39 preguntas) → módulo TZ

**Esta es la divergencia más grande, y es deliberada.** El propio código del
prototipo lo dice (`registroSalviaSchema.js:200-208`): es una versión
*simplificada a propósito* ("banco de preguntas más corto para este mockup de
apoyo"), no una réplica. V2 tiene 38 preguntas booleanas reales + 1 banner;
el prototipo tiene 15.

| Bloque | V2 | Prototipo | Con equivalente en el otro lado |
|---|:---:|:---:|---|
| Comunes (siempre visibles) | 6 | 5 | ~5 — ver detalle abajo |
| Pareja íntima | 18 | 6 | 6 limpios |
| No pareja | 14 | 4 | 3 limpios + 1 aproximado |
| Banner de riesgo | 1 (info) | — | el prototipo no tiene banner, calcula el nivel inline |

Equivalencias encontradas (las demás preguntas de cada rama son exclusivas de V2):

- **Comunes:** `tz_incremento_violencia`≈`aggressor_violence_physical_increase`,
  `tz_uso_arma`≈`aggressor_weapon_used` (y también cubre
  `aggressor_has_access_to_weapons`, que V2 trae como pregunta separada),
  `tz_amenaza_muerte`≈`aggressor_threat_kill`,
  `tz_capaz_matar`≈`aggressor_capable_of_killing`. El prototipo no tiene
  equivalente de `aggressor_pursues_spies_destroys` (persigue/espía/destruye).
  `tz_ideacion_autolesion` es un caso especial: en V2 la pregunta de
  autolesión está **duplicada dentro de cada rama**
  (`thoughts_of_self_harm` #20 en pareja, `thoughts_of_self_harm_2` #35 en
  no-pareja); el prototipo la promovió a "común" única (Observación #15 fue
  la que reorganizó cuándo se abren las ramas, este fue un ajuste de diseño
  asociado).
- **Pareja:** `tz_p_convive`≈`still_lives_with_aggressor`,
  `tz_p_controla`≈`partner_controls`, `tz_p_forzo_sexo`≈`aggressor_forced_sex`,
  `tz_p_intento_ahorcar`≈`aggressor_attempted_strangulation`,
  `tz_p_consumo_alcohol`≈`aggressor_is_alcoholic`,
  `tz_p_consumo_sustancias`≈`aggressor_consumes_drugs` (estos dos últimos
  separados por la Observación #4 del PDF, igual que en V2 que ya los trae
  separados). Exclusivas de V2: desempleo, otras denuncias, antecedentes
  penales, dejó de buscar ayuda, golpeó en vulnerabilidad, usa la salud para
  chantajear, amenazó suicidio, amenazó dañar a miembros de la familia,
  limita redes de apoyo, amenazó revelar orientación sexual, celoso violento
  (11 preguntas).
- **No pareja:** `tz_np_jerarquia`≈`aggressor_hierarchy`,
  `tz_np_espacios_comunes`≈`aggressor_common_spaces`,
  `tz_np_antecedentes`≈`aggressor_has_penal_background_2`.
  `tz_np_acoso_reiterado` es solo **aproximado** a
  `aggressor_sexually_harassment` (V2 trae esa pregunta duplicada, 28 y 29).
  Exclusivas de V2: aprovechó vulnerabilidad física, desempleo, motivada por
  género, usa drogas, es alcohólico, controla, amenazó dañar miembros,
  autolesión (duplicada, ya cubierta arriba), usó la posición de autoridad
  (9 preguntas).

Todas las 38 preguntas reales de V2 se guardan igual: `bool_enum` → FK a
`victim_case_form2_enums` categoría `yes_no`, columna propia en
`victim_case_form2` (ver listado completo en
[form-registro-caso-v2-data.md](form-registro-caso-v2-data.md) o
directamente en `victim_case_form_fields.go:72-111`). El prototipo no
persiste nada; solo suma puntos en memoria para `calcularNivelRiesgoRegistro`.

## Sección 6 — Datos Personales (34 preguntas) → módulo DP (resto)

| # | Pregunta | Dónde | Prototipo | Columna V2 |
|---|---|---|---|---|
| 1 | Fecha de nacimiento | Ambas | DP.fecha_nacimiento | `victim_case_form2_birth_date` (date) |
| 2 | ¿Dificultades físicas/mentales/sensoriales? | Ambas§ | DP.tiene_discapacidad | `victim_case_form2_physical_mental_sensory_difficulties` (bool_enum) |
| 3–11 | 9 preguntas de escala (oír, hablar, ver, moverse, tomar, entender, comer, interactuar, cotidianas) | **Solo V2** | — | `victim_case_form2_activities_unable_to_{hear,talk,see,move,take,understand,eat,interact,do_everyday}` (scale, entero 0-3 cada una) |
| 12 | Ley 1996 | Ambas | DP.ley_1996 | sin columna propia → tabla de relación (enumN) |
| 13 | Nacionalidad | Ambas | DP.nacionalidad | `victim_case_form2_nationality` (enum1) |
| 14 | Nacionalidad específica | Solo V2 | — | `victim_case_form2_specified_nationality` (enum1) |
| 15 | Condición migratoria | Solo V2 | — | `victim_case_form2_migration_condition` (enum1) |
| 16 | Identidad de género | Ambas | DP.identidad_genero | `victim_case_form2_gender_identity` (enum1) |
| 17 | Orientación sexual | Ambas | DP.orientacion_sexual | `victim_case_form2_sexual_orientation` (enum1) |
| 18 | Sexo asignado al nacer | Ambas | DP.sexo_asignado_nacer | `victim_case_form2_assigned_sex_at_birth` (enum1) |
| 19 | Población especialmente protegida | Ambas | DP.poblacion_protegida | sin columna propia → tabla de relación (enumN) |
| 20 | Afiliación étnica | Ambas | DP.grupo_etnico | `victim_case_form2_ethnic_affiliation` (enum1) |
| 21 | Pueblo indígena | Solo V2 | — | `victim_case_form2_indigenous_people` (enum1) |
| 22 | ¿Reconocimiento como campesina? | Ambas | DP.reconocimiento_campesina | `victim_case_form2_campesino_recognition` (bool_enum) |
| 23 | Estado civil | Ambas | DP.estado_civil | `victim_case_form2_marital_status` (enum1) |
| 24 | Último nivel educativo | Ambas | DP.ultimo_nivel_educativo | `victim_case_form2_last_education_level` (enum1) |
| 25 | Ocupación | Ambas | DP.ocupacion | `victim_case_form2_occupation` (enum1) |
| 26 | Forma de generación de ingresos | Ambas | DP.forma_generacion_ingresos | `victim_case_form2_income_generation_method` (enum1) |
| 27 | Modalidad ASP | Solo V2 | — | sin columna propia → tabla de relación (enumN) |
| 28 | Relación laboral | Solo V2 | — | `victim_case_form2_employment_relationship` (enum1) |
| 29 | Inicio aproximado ASP | Solo V2 | — | `victim_case_form2_approx_start_asp` (timestamp) |
| 30 | Razón ASP | Solo V2 | — | sin columna propia → tabla de relación (enumN) |
| 31 | Forma de tenencia de vivienda | Ambas | DP.tenencia_vivienda | `victim_case_form2_housing_tenancy_form` (enum1) |
| 32 | Estrato de vivienda | Ambas | DP.estrato_vivienda | `victim_case_form2_housing_stratum` (enum1) |
| 33 | Personas a cargo | Ambas | DP.personas_a_cargo | sin columna propia → tabla de relación (enumN) |
| 34 | ¿Actualmente en embarazo? | Ambas¶ | DP.actualmente_embarazada | `victim_case_form2_currently_pregnant` (bool_enum) |

§ Mismo caso que la S1Q13: es el mismo gate que `tiene_discapacidad`, V2 la
repite aquí por construcción del formulario original (`set_victim_case.html`
ya la tenía dos veces).

¶ En V2 es incondicional; el prototipo la condiciona a
`sexo_asignado_nacer = Femenino` (Observación #8) — mismo campo, regla de
visibilidad más fina del lado del prototipo.

**El hueco más grande de toda la comparación está aquí:** las 9 preguntas de
dificultad (3–11) son el instrumento corto del Washington Group que alimenta
la certificación de discapacidad — el prototipo nunca las implementó, solo
se quedó con el booleano-gate. Y el bloque ASP completo (27, 29, 30 + la 28
de relación laboral) tampoco existe — son datos sensibles sobre actividad
sexual pagada que el prototipo omite enteramente, no solo simplifica.

## Sección 7 — Plan de Acción (2 preguntas) → módulo GS

| # | Pregunta | Dónde | Prototipo | Columna V2 |
|---|---|---|---|---|
| 1 | Plan de acción | Ambas | GS.plan_accion | sin columna propia → tabla de relación (enumN) |
| 2 | Explicación de la gestión Salvia | Ambas | GS.explicacion_gestion | `victim_case_form2_saliva_management_explanation` (text) |

`GS.requiere_remision` y `GS.procesos_remision` **no tienen equivalente en
V2** — son exclusivos del prototipo. Mi lectura: V2 probablemente empaqueta
esas opciones de remisión dentro del único multiselect "Plan de acción"
(`enums.victim_case_form2_action_plan`); el prototipo lo descompuso en un
booleano-gate + multiselect aparte. No alcancé a confirmar si las opciones
reales de `victim_case_form2_action_plan` incluyen específicamente
"Medidas de emergencia / Psicosocial / Estabilización económica / Atención
Especializada" — valdría la pena revisar el seed
(`src/cmd/seed/seed_registro_caso.sql`) si se quiere decidir si esto es una
mejora de UX a trasladar al V2 real o una divergencia a corregir en el
prototipo.

## Sección 8 — Denuncia Fácil (1 pregunta) → sin equivalente en el prototipo

| # | Pregunta | Dónde | Columna V2 |
|---|---|---|---|
| 1 | ¿Permite el uso del canal de Denuncia Fácil? | **Solo V2** | `victim_case_form2_allows_easy_report` (bool_enum) |

Sección completa ausente del prototipo — no hay ningún control relacionado
con "Denuncia Fácil" en `registroSalviaSchema.js`.

## Sección 9 — Lugar de Atención (3 preguntas) → sin equivalente en el prototipo

| # | Pregunta | Dónde | Columna V2 |
|---|---|---|---|
| 1 | Departamento de atención | Solo V2 | cascada UI |
| 2 | Ciudad de atención | Solo V2 | cascada UI |
| 3 | Municipio de atención | Solo V2 | `victim_case.victim_case_victim_town_code` (no `form2`; se actualiza al activar el caso) |

También ausente por completo. El prototipo nunca pregunta dónde se atiende
el caso (se asume implícitamente la sede del agente).

---

## Campos exclusivos del prototipo (no existen en V2 real)

| Campo | Módulo | Origen | Columna propuesta (no existe hoy) |
|---|---|---|---|
| `correo_victima` | DP | Observación #12 del PDF | `victim_case_form2_victim_email` (text) |
| `requiere_remision` | GS | Decisión de diseño del prototipo | `victim_case_form2_requires_referral` (bool_enum) |
| `procesos_remision` | GS | Decisión de diseño del prototipo | sin columna propia → tabla de relación (enumN), análogo a `plan_accion` |

---

## Notas para decidir qué hacer con esto

- **Si el prototipo se usa para validar con el equipo antes de tocar el V2
  real:** los 3 campos exclusivos del prototipo (correo víctima, remisión)
  son candidatos razonables para proponer como *nuevas* preguntas en el V2 —
  nacieron de observaciones reales sobre el mockup, no de la nada.
- **Si la pregunta es "qué le falta al prototipo para ser fiel a V2":** los
  bloques grandes a priorizar, en orden de impacto, son: (1) la batería de
  discapacidad de 9 preguntas, (2) Tamizaje completo (quedó explícitamente
  simplificado, no es un olvido), (3) Denuncia Fácil y Lugar de Atención
  (secciones enteras ausentes), (4) el bloque ASP.
- **GAP ya documentado y que sigue abierto en el V2 real, no en el
  prototipo:** G-12 (`aggressor_occupation` usa por error las opciones de
  "relación con el agresor") — el prototipo no hereda ese bug porque tiene su
  propio catálogo de ocupación del agresor.
