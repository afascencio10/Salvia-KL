# `set_victim_case` — Interfaz de la Pantalla

Pantalla de registro de un caso nuevo de violencia basada en genero.  
Ruta: `/salvia/casos/nuevo` (tambien `/salvia/casos/:id/nuevo` cuando viene de un contacto previo)  
Template: `frontend/html/salvia/victim_case/set_victim_case.html`  
Facade: `VictimCasePOST_GET` (GET) / `VictimCasePOST` (POST)  
Controlador: `SetVictimCase`

---

## Modelo de datos principal

```
data.victimCase
  ├── names                          // string — nombres de la victima
  ├── lastNames                      // string — apellidos
  ├── docType                        // string — codigo tipo documento
  ├── docNumber                      // string — numero documento
  ├── townCode                       // string — municipio de atencion (select encadenado)
  ├── entityBranchesByMomentIdx      // object — sedes seleccionadas por momento/sector
  ├── newUser                        // object — {login, pass} devuelto tras crear el caso
  └── form2                          // object — formulario principal (ver abajo)
```

---

## Arbol de interfaz

```
container
│
├── [v-if !authorizationAnswer || authorizationAnswer.code != 'y']
│   └── SeccionAutorizacion
│       ├── Titulo  "Autorizacion de datos personales"
│       ├── TextoLegal  (Ley 1581 de 2012, derechos ARCO)
│       └── SelectSingle  authorizationAnswer  (yes_no)
│
├── [v-if authorizationAnswer.code == 'y']  ── Formulario principal ──
│
│   ├── AlertDebidaDiligencia
│   │   └── Texto informativo SALVIA / Ministerio
│   │
│   ├── ══ SECCION: Datos de la Victima ══════════════════════════════
│   │   ├── Row
│   │   │   ├── Input:text  names *                         // Nombres
│   │   │   └── Input:text  lastNames *                     // Apellidos
│   │   ├── Row
│   │   │   ├── Input:text  form2.identityName              // Nombre identitario
│   │   │   └── Input:number  form2.phone *                 // Telefono
│   │   ├── Row
│   │   │   ├── SelectEnum  docType * (docType2)            // Tipo documento
│   │   │   └── Input:text  docNumber *                     // Numero documento
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.residenceZone * (victim_case_form2_facts_zone)
│   │   │   └── SelectEncadenado  Departamento3             // Depto residencia
│   │   ├── Row
│   │   │   ├── [v-if cities3] SelectEncadenado  Ciudad3
│   │   │   └── [v-if towns3]  SelectEncadenado  form2.residenceTownCode
│   │   └── Row
│   │       └── Input:text  form2.residenceAddress *        // Direccion residencia
│   │
│   ├── ── Accesibilidad ───────────────────────────────────────
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.personWithDisability * (yes_no)
│   │   │   └── SelectEnum:multiple  form2.adjustmentsGBV * (victim_case_form2_adjustments_gbv)
│   │   └── Row
│   │       ├── SelectEnum  form2.requireLanguageInterpreter * (yes_no)
│   │       └── [v-if requireLanguageInterpreter.code=='y']
│   │           Input:text  form2.languageAssistance
│   │
│   ├── ══ SECCION: Contacto de Apoyo ═══════════════════════════
│   │   ├── Row
│   │   │   ├── Input:text  form2.supportContactNames
│   │   │   └── Input:number  form2.supportContactPhone
│   │   └── Row
│   │       ├── Input:email  form2.supportContactEmail
│   │       └── SelectEnum  form2.supportContactKinship (victim_case_form2_support_contact_kinship)
│   │
│   ├── ══ SECCION: Hechos ══════════════════════════════════════
│   │   ├── Row
│   │   │   └── Textarea  form2.factsDescription *          // Relato de hechos
│   │   ├── Row
│   │   │   ├── DateTime  form2.factsDate * (DD/MM/YYYY)    // Fecha hechos
│   │   │   └── DateTime  form2.factsStartTime * (hh:mm)    // Hora hechos
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.factsZone * (victim_case_form2_facts_zone)
│   │   │   └── SelectEncadenado  Departamento2             // Depto hechos
│   │   ├── Row
│   │   │   ├── [v-if cities2] SelectEncadenado  Ciudad2
│   │   │   └── [v-if towns2]  SelectEncadenado  form2.factsTownCode
│   │   ├── Row
│   │   │   ├── Input:text  form2.factsAddress *            // Direccion hechos
│   │   │   └── SelectEnum  form2.scenarioViolence * (victim_case_form2_scenario_violence)
│   │   ├── Row
│   │   │   ├── SelectEnum:multiple  form2.typeViolenceExperienced * (victim_case_form2_type_violence_experienced)
│   │   │   │   @change → violenceTypeChanged (recalcula subtipos)
│   │   │   └── [v-if typeViolenceExperienced.length > 0]
│   │   │       SelectEnum:multiple  form2.subtypeViolenceExperienced * (violenceSubtypes — dinamico)
│   │   ├── Row
│   │   │   ├── SelectEnum:multiple  form2.scopeOfViolence * (victim_case_form2_scope_of_violence)
│   │   │   │   @change → checkWorkplaceSectorOccurrence
│   │   │   └── [v-if workplaceSectorOccurrence]
│   │   │       SelectEnum  form2.workplaceSectorOccurrence * (victim_case_form2_workplace_sector_occurrence)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.violenceMotivatedByGender * (yes_no)
│   │   │   └── SelectEnum  form2.reportedPreviously * (yes_no)
│   │   ├── [v-if reportedPreviously.code == 'y'] Row
│   │   │   ├── SelectEnum:multiple  form2.whoReportTo * (victim_case_form2_who_report_to)
│   │   │   └── SelectEnum  form2.attentionWasAppropriate * (yes_no)
│   │   └── Row
│   │       └── SelectEnum  form2.recurrenceAggression * (victim_case_form2_recurrence_aggression)
│   │
│   ├── ══ SECCION: Agresor ═════════════════════════════════════
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.numAgressors * (victim_case_form2_num_agressors)
│   │   │   └── SelectEnum  form2.proximityPrincipalAggressor * (victim_case_form2_proximity_principal_aggressor)
│   │   │       @change → proximityWithAggressorChanged (determina wasPartner/partnerKnown)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.relationshipWithPresumedAggressor * (victim_case_form2_relationship_with_presumed_aggressor_01)
│   │   │   │   @change → relationWithAggressorChanged
│   │   │   └── SelectEnum  form2.aggressorOccupation * (aggressor_occupation)
│   │   ├── [v-if partnerKnown] Row
│   │   │   └── SelectEnum  form2.economicallyDependent * (yes_no)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.aggressorGenderIdentity * (victim_case_form2_aggressor_gender_identity)
│   │   │   └── Input:text  form2.aggressorNames
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.aggressorDocType (victim_case_form2_victim_doc_type)
│   │   │   └── Input:text  form2.aggressorDocNumber
│   │   └── Row
│   │       ├── Input:text  form2.aggressorAddress
│   │       └── Input:number  form2.aggressorPhone
│   │
│   ├── ══ SECCION: Tamizaje ════════════════════════════════════
│   │   │   Todas las preguntas son SelectEnum (yes_no), @change → updateTamizajeScore
│   │   │
│   │   ├── ── Preguntas comunes (6) — siempre visibles ────────
│   │   │   ├── form2.aggressorViolencePhysicalIncrease *
│   │   │   ├── form2.aggressorWeaponUsed *
│   │   │   ├── form2.aggressorThreatKill *
│   │   │   ├── form2.aggressorPursuesSpiesDestroys *
│   │   │   ├── form2.aggressorCapableOfKilling *
│   │   │   └── form2.aggressorHasAccessToWeapons *
│   │   │       @change → updateTamizajeAggressor (activa/desactiva bloque pareja o no-pareja)
│   │   │
│   │   ├── [v-if tamizajeAggressorCheck && wasPartner]
│   │   │   ── Preguntas pareja intima (18) ────────────────────
│   │   │   ├── form2.partnerUnemployed *
│   │   │   ├── form2.partnerOtherDenunciations *
│   │   │   ├── form2.aggressorHasPenalBackground *
│   │   │   ├── form2.stoppedSeekingHelp *
│   │   │   ├── form2.aggressorForcedSex *
│   │   │   ├── form2.aggressorAttemptedStrangulation *
│   │   │   ├── form2.aggressorConsumesDrugs *
│   │   │   ├── form2.aggressorIsAlcoholic *
│   │   │   ├── form2.partnerControls *
│   │   │   ├── form2.aggressorHadHitInVulnerability *
│   │   │   ├── form2.victimHealthToBlackmail *
│   │   │   ├── form2.partnerThreatenedSuicide *
│   │   │   ├── form2.partnerThreatenedDamageMembers *
│   │   │   ├── form2.thoughtsOfSelfHarm *
│   │   │   ├── form2.aggressorLimitsContactSupportNetworks *
│   │   │   ├── form2.threatenedRevealSexualOrientation *
│   │   │   ├── form2.stillLivesWithAggressor *
│   │   │   └── form2.aggressorViolentlyJealous *
│   │   │
│   │   ├── [v-if tamizajeAggressorCheck && !wasPartner]
│   │   │   ── Preguntas no-pareja (14) ────────────────────────
│   │   │   ├── form2.aggressorTakenAdvantagePhysicalVulnerability *
│   │   │   ├── form2.aggressorUnemployed *
│   │   │   ├── form2.aggressorHasPenalBackground2 *
│   │   │   ├── form2.aggressorSexuallyHarassment *
│   │   │   ├── form2.aggressorSexuallyHarassment2 *
│   │   │   ├── form2.violenceMotivatedByGender2 *
│   │   │   ├── form2.aggressorUseDrugs *
│   │   │   ├── form2.aggressorIsAlcoholic2 *
│   │   │   ├── form2.aggressorControls *
│   │   │   ├── form2.aggressorThreatenedDamageMembers *
│   │   │   ├── form2.thoughtsOfSelfHarm2 *
│   │   │   ├── form2.aggressorCommonSpaces *
│   │   │   ├── form2.aggressorHierarchy *
│   │   │   └── form2.aggressorUsedPositionAuthority *
│   │   │
│   │   └── [v-if riskLevel > 0] RiskBadge
│   │       ├── Label  "Puntaje de riesgo: {riskScore}"
│   │       └── Label  "Nivel de riesgo: {riskLevel}"
│   │       CSS dinamico: riesgo_bajo | riesgo_moderado | riesgo_alto | riesgo_extremo
│   │
│   ├── ══ SECCION: Datos Personales ════════════════════════════
│   │   ├── Row
│   │   │   ├── DateTime  form2.birthDate * (DD/MM/YYYY)
│   │   │   └── SelectEnum  form2.physicalMentalSensoryDifficulties * (yes_no)
│   │   │
│   │   ├── [v-if physicalMentalSensoryDifficulties.code == 'y']
│   │   │   ├── InfoEscala  (4 estrellas: sin dificultad → no puede hacerlo)
│   │   │   ├── StarRating × 9
│   │   │   │   ├── activitiesUnableToHear
│   │   │   │   ├── activitiesUnableToTalk
│   │   │   │   ├── activitiesUnableToSee
│   │   │   │   ├── activitiesUnableToMove
│   │   │   │   ├── activitiesUnableToTake
│   │   │   │   ├── activitiesUnableToUnderstand
│   │   │   │   ├── activitiesUnableToEat
│   │   │   │   ├── activitiesUnableToInteract
│   │   │   │   └── activitiesUnableToDoEveryday
│   │   │   └── SelectEnum:multiple  form2.law1996 * (victim_case_form2_law_1996)
│   │   │
│   │   ├── Row
│   │   │   └── SelectEnum  form2.nationality * (victim_case_form2_nationality)
│   │   ├── [v-if nationality.code == 'ex'] Row
│   │   │   ├── SelectEnum  form2.specifiedNationality * (victim_case_form2_specified_nationality)
│   │   │   └── SelectEnum  form2.migrationCondition * (victim_case_form2_migration_condition)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.genderIdentity * (victim_case_form2_gender_identity)
│   │   │   └── SelectEnum  form2.sexualOrientation * (victim_case_form2_sexual_orientation)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.assignedSexAtBirth * (victim_case_form2_assigned_sex_at_birth)
│   │   │   └── SelectEnum:multiple  form2.speciallyProtectedPopulation * (victim_case_form2_specially_protected_population)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.ethnicAffiliation * (victim_case_form2_ethnic_affiliation)
│   │   │   └── [v-if ethnicAffiliation.code == 'in']
│   │   │       SelectEnum  form2.indigenousPeople * (victim_case_form2_indigenous_people)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.campesinoRecognition * (yes_no)
│   │   │   └── SelectEnum  form2.maritalStatus * (victim_case_form2_marital_status)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.lastEducationLevel * (victim_case_form2_last_education_level)
│   │   │   └── SelectEnum  form2.occupation * (victim_case_form2_occupation)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.incomeGenerationMethod * (victim_case_form2_income_generation_method)
│   │   │   ├── [v-if incomeGenerationMethod.code == 'pr']
│   │   │   │   SelectEnum:multiple  form2.aspMode * (victim_case_form2_asp_mode)
│   │   │   └── [v-if incomeGenerationMethod.code == 'em']
│   │   │       SelectEnum  form2.employmentRelationship * (victim_case_form2_employment_relationship)
│   │   ├── [v-if incomeGenerationMethod.code == 'pr'] Row
│   │   │   ├── Input:text  form2.approxStartAsp *
│   │   │   └── SelectEnum:multiple  form2.reasonASP * (victim_case_form2_reason_asp)
│   │   ├── Row
│   │   │   ├── SelectEnum  form2.housingTenancyForm * (victim_case_form2_housing_tenancy_form)
│   │   │   └── SelectEnum  form2.housingStratum * (victim_case_form2_housing_stratum)
│   │   └── Row
│   │       ├── SelectEnum:multiple  form2.hasDependents * (victim_case_form2_has_dependents)
│   │       └── SelectEnum  form2.currentlyPregnant * (yes_no)
│   │
│   ├── ══ SECCION: Plan de Accion ══════════════════════════════
│   │   ├── Row
│   │   │   └── SelectEnum:multiple  form2.actionPlan * (victim_case_form2_action_plan)
│   │   └── Row
│   │       └── Textarea  form2.salivaManagementExplanation *
│   │
│   ├── ══ SECCION: Denuncia Facil ══════════════════════════════
│   │   └── Row
│   │       └── SelectEnum  form2.allowsEasyReport * (yes_no)
│   │
│   ├── ══ SECCION: Lugar de Atencion ═══════════════════════════
│   │   ├── SelectEncadenado  Departamento → Ciudad → Municipio (townCode)
│   │   │
│   │   └── [v-if entitiesByMomentAndSector]
│   │       ══ SECCION: Asignacion de Ruta ══════════════════════
│   │       └── MomentGrid
│   │           └── MomentSection × N  [v-for moments]
│   │               └── SectorGroup × N  [v-for sectors]
│   │                   └── EntityBranchSelect × N  [v-for branches]
│   │                       entityBranchesByMomentIdx[moment][sector][entity]
│   │
│   └── ActionButtons
│       ├── BtnVolver  (.btn-warning)  → submit('back')
│       └── BtnGuardar  (.btn-success)  → submit('save')
│
└── FinishedOverlay  (#finishedOverlay)  [display:none → visible tras exito]
    ├── Titulo  "Caso registrado exitosamente"
    ├── Credenciales
    │   ├── Usuario:  data.victimCase.newUser.login
    │   └── Clave:    data.victimCase.newUser.pass
    └── Botones
        ├── BtnSMS
        └── BtnFinalizar  → finished()
```

---

## Selects encadenados (cascading)

El formulario tiene **3 cadenas** de Departamento → Ciudad → Municipio:

| Cadena | Variable depto | Variable ciudad | Variable municipio | Uso |
|--------|----------------|-----------------|-------------------|-----|
| 1 | `departmentSelected` | `citySelected` | `victimCase.townCode` | Lugar de atencion |
| 2 | `department2Selected` | `city2Selected` | `form2.factsTownCode` | Lugar de hechos |
| 3 | `department3Selected` | `city3Selected` | `form2.residenceTownCode` | Residencia victima |

Cada cadena usa `@change` para cargar la lista siguiente via API:
- `changeDepartmentN($event)` → carga ciudades
- `changeCityN($event)` → carga municipios

---

## Campos condicionales

| Campo visible cuando... | Condicion |
|---|---|
| Subtipos violencia | `typeViolenceExperienced.length > 0` |
| A quien denuncio | `reportedPreviously.code == 'y'` |
| Dependencia economica | `partnerKnown == true` |
| Tamizaje pareja (18 preguntas) | `tamizajeAggressorCheck && wasPartner` |
| Tamizaje no-pareja (14 preguntas) | `tamizajeAggressorCheck && !wasPartner` |
| Sector laboral ocurrencia | `workplaceSectorOccurrence == true` |
| Escala de dificultades (9 star ratings) | `physicalMentalSensoryDifficulties.code == 'y'` |
| Ley 1996 | `physicalMentalSensoryDifficulties.code == 'y'` |
| Nacionalidad especifica + condicion migratoria | `nationality.code == 'ex'` |
| Pueblo indigena | `ethnicAffiliation.code == 'in'` |
| Modalidad ASP + razon + fecha inicio | `incomeGenerationMethod.code == 'pr'` |
| Relacion laboral | `incomeGenerationMethod.code == 'em'` |
| Interprete idioma | `requireLanguageInterpreter.code == 'y'` |

---

## Logica del tamizaje (riskScore / riskLevel)

El tamizaje calcula un puntaje sumando respuestas "si" y determina el nivel de riesgo.

**Pareja intima** (`wasPartner == true`): 6 comunes + 18 especificas = 24 preguntas max

| Nivel | Puntaje |
|-------|---------|
| 1 - Bajo | 0–4 |
| 2 - Moderado | 5–8 |
| 3 - Alto | 9–15 |
| 4 - Extremo | 16–24 |

**No pareja** (`wasPartner == false`): 6 comunes + 14 especificas = 20 preguntas max

| Nivel | Puntaje |
|-------|---------|
| 1 - Bajo | 0–2 |
| 2 - Moderado | 3–5 |
| 3 - Alto | 6–8 |
| 4 - Extremo | 9–20 |

---

## Tipos de input usados

| Tipo | Componente | Cantidad aprox. |
|---|---|---|
| `Input:text` | `labeled_input.html` (type=text) | ~12 |
| `Input:number` | `labeled_input.html` (type=number) | ~3 |
| `Input:email` | `labeled_input.html` (type=email) | 1 |
| `Textarea` | `labeled_textarea.html` | 2 |
| `DateTime` | `date_time.html` | 3 |
| `SelectEnum` (single) | `labeled_select_entity.html` | ~35 |
| `SelectEnum` (multiple) | `labeled_select_entity.html` | ~12 |
| `SelectEncadenado` | `<select>` nativo con @change | 9 (3 cadenas x 3) |
| `StarRating` | `labeled_star_rating.html` | 9 |

---

## Flujo POST

1. Usuario hace click en "Guardar" → `submit('save')`
2. Vue serializa `data.victimCase` como JSON
3. POST a `/salvia/casos` (o `/salvia/casos/:id` si viene de contacto)
4. Backend: `SetVictimCase` valida, crea usuario, inserta caso, genera calendario de seguimientos
5. Si exito → muestra `FinishedOverlay` con login/password
6. Si error → muestra errores campo por campo en `errors.form2.*`
