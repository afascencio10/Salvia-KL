# Repeater de Barreras — Estructura de Preguntas

**Sección:** Identificación de Barreras (Sección 4 del Formulario de Seguimiento)
**Repeater Group ID:** `5fd3ecdc-2e5f-4b31-97ef-8a994580586a`
**Nombre del repeater:** Barreras identificadas
**item_name:** Barrera
**min_repetitions:** 1
**max_repetitions:** sin límite

> Cada entrada del repeater representa una barrera individual. Todas las preguntas (sectoriales y estructurales) se registran por barrera.

---

## Estructura de preguntas por entrada

### 1. Sector de la barrera
| Campo | Valor |
|---|---|
| ID actual | `f19378b6-55c5-4fdf-b765-7ebcc3978741` |
| Tipo | `dropdown` |
| Requerido | ✅ |
| Condición | — siempre visible |

| Label | Value |
|---|---|
| Salud | `salud` |
| Justicia | `justicia` |
| Protección | `proteccion` |
| Otras instituciones | `otras_instituciones` |
| Barrera Transversal | `barrera_transversal` |

> ℹ️ **Barrera Transversal** no tiene bloque de preguntas propio — al seleccionarlo solo se muestran las preguntas comunes (Q12 en adelante).

---

### Bloque Salud
*Visible cuando Q1 = `salud` (EQUALS)*

#### 2. Barreras identificadas en Salud
| Campo | Valor |
|---|---|
| ID actual | `5fc1f2af-cc30-41f0-aa31-4731e5cb674c` |
| Tipo | `multiple` |
| Requerido | ❌ |

| # | Label | Value |
|---|---|---|
| 1 | Dificultades de aseguramiento o afiliación en salud | `dificultades_aseguramiento` |
| 2 | Negaciones en los servicios de urgencias | `negacion_urgencias` |
| 3 | Fallas en la calidad del servicio de urgencias | `fallas_calidad_urgencias` |
| 4 | Demoras en la asignación de citas y continuidad de tratamientos | `demoras_citas` |
| 5 | Faltas de activación de protocolos para violencia sexual | `falta_protocolo_violencia_sexual` |
| 6 | Faltas de activación de protocolos para ataques con agentes químicos | `falta_protocolo_agentes_quimicos` |
| 7 | Demoras y dificultades en la valoración medicolegal | `demoras_valoracion_medicolegal` |
| 8 | Negaciones en el acceso a la Interrupción Voluntaria del Embarazo (IVE) | `negacion_ive` |
| 9 | Exigencias indebidas de autorización para procedimientos en personas con discapacidad cognitiva y/o psicosocial | `exigencias_autorizacion_discapacidad` |
| 10 | Atenciones en salud no centradas en cosmovisiones y prácticas culturales propias | `falta_enfoque_cultural_salud` |
| 11 | Negaciones o imposiciones de procedimientos a personas con diversidad sexual o identidad de género diversa | `negacion_procedimientos_osigd` |
| 12 | Otras barreras en salud | `otras_barreras_salud` |

#### 3. ¿A qué institución acudió? (Salud)
| Campo | Valor |
|---|---|
| ID actual | `e88fb2bf-7196-4bea-91e6-fe78ccdedac1` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | Q1 = `salud` (EQUALS) |

| # | Label | Value |
|---|---|---|
| 1 | Hospital Público | `hospital_publico` |
| 2 | Hospital Privado | `hospital_privado` |
| 3 | IPS | `ips` |
| 4 | EPS | `eps` |
| 5 | Puesto de Salud | `puesto_salud` |
| 6 | Consultorio médico | `consultorio_medico` |
| 7 | No sabe / No responde | `no_sabe` |

#### 4. Otra barrera en Salud ¿cuál?
| Campo | Valor |
|---|---|
| ID actual | `bebf6e6c-0200-4b53-886c-c01f591eaa62` |
| Tipo | `text` |
| Requerido | ❌ |
| Condición | Q2 CONTAINS `otras_barreras_salud` |

---

### Bloque Justicia
*Visible cuando Q1 = `justicia` (EQUALS)*

#### 5. Barreras identificadas en Justicia
| Campo | Valor |
|---|---|
| ID actual | `2bec977e-97c7-42d7-a00a-b536af8038eb` |
| Tipo | `multiple` |
| Requerido | ❌ |

| # | Label | Value |
|---|---|---|
| 1 | Negativa institucional para recibir la denuncia | `negativa_recibir_denuncia` |
| 2 | Dilación en la recepción y registro de la denuncia | `dilacion_registro_denuncia` |
| 3 | Denegación al derecho a la no confrontación | `denegacion_no_confrontacion` |
| 4 | Falta de entrega de información y documentación procesal | `falta_info_documental` |
| 5 | Formalismo excesivo en los trámites judiciales | `formalismo_excesivo` |
| 6 | Imposición de cargas probatorias injustificadas o excesivas | `cargas_probatorias_excesivas` |
| 7 | Falta de celeridad en la investigación judicial | `falta_celeridad_investigacion` |
| 8 | Incumplimiento de órdenes judiciales de protección o sanción | `incumplimiento_ordenes_judiciales` |
| 9 | Tipificación errónea del delito | `tipificacion_erronea` |
| 10 | Traslado inadecuado o incompleto del expediente judicial | `traslado_inadecuado_expediente` |
| 11 | Dificultades en la reasignación de procesos de Comisaría o Defensoría de Familia tras cambio de municipio | `dificultad_reasignacion_municipio` |
| 12 | Ausencia o dificultad de acceso a representación judicial | `falta_representacion_judicial` |
| 13 | Vacíos normativos y definiciones restrictivas en VBG | `vacios_normativos_vbg` |
| 14 | Riesgo de vencimiento de términos | `riesgo_vencimiento_terminos` |
| 15 | Medidas preventivas de libertad dirigidas al perpetrador que ponen en riesgo a la mujer y su núcleo familiar | `medidas_libertad_riesgo` |
| 16 | Falta de protección en la divulgación de información confidencial del proceso | `falta_proteccion_info_confidencial` |
| 17 | Deficiencias en la notificación y recolección de datos del agresor | `deficiencias_notificacion_agresor` |
| 18 | Otras barreras en justicia | `otras_barreras_justicia` |

#### 6. ¿A qué institución acudió? (Justicia)
| Campo | Valor |
|---|---|
| ID actual | `4b4997fa-67a0-4479-852d-df9e7bfb2b3e` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | Q1 = `justicia` (EQUALS) |

| # | Label | Value |
|---|---|---|
| 1 | CAI - Policía (Comando de Atención Inmediata) | `cai_policia` |
| 2 | CAIVAS (Centro de Atención Integral a Víctimas de Violencia Sexual) | `caivas` |
| 3 | Casas de Justicia | `casas_justicia` |
| 4 | CAVIV (Centro de Atención Integral contra la Violencia Intrafamiliar) | `caviv` |
| 5 | Comisaría de Familia | `comisaria_familia` |
| 6 | Estación de Policía | `estacion_policia` |
| 7 | Fiscalía General de la Nación | `fiscalia` |
| 8 | Inspección de Policía | `inspeccion_policia` |
| 9 | Instituto Nacional de Medicina Legal | `medicina_legal` |
| 10 | Jueces Civiles o Promiscuos Municipales | `jueces_civiles` |
| 11 | Jueces de Control de Garantías o Jueces Penales | `jueces_penales` |
| 12 | Policía Judicial | `policia_judicial` |
| 13 | Unidad de Reacción Inmediata | `unidad_reaccion_inmediata` |
| 14 | No sabe / No responde | `no_sabe` |

#### 7. Otra barrera en Justicia ¿cuál?
| Campo | Valor |
|---|---|
| ID actual | `96f0c507-64c1-446b-b07d-83d5ad1d7172` |
| Tipo | `text` |
| Requerido | ❌ |
| Condición | Q5 CONTAINS `otras_barreras_justicia` |

---

### Bloque Protección
*Visible cuando Q1 = `proteccion` (EQUALS)*

#### 8. Barreras identificadas en Protección
| Campo | Valor |
|---|---|
| ID actual | `66c9fc1e-9b5e-4ad4-999f-7aeb483d84dc` |
| Tipo | `multiple` |
| Requerido | ❌ |

| # | Label | Value |
|---|---|---|
| 1 | Ausencia o demoras en la adopción de medidas de protección urgentes o definitivas en Comisaría de Familia | `demoras_medidas_comisaria` |
| 2 | Ausencia o dilación en adopción de medidas de protección físicas | `dilacion_medidas_fisicas` |
| 3 | Demora o ausencia en la activación de medidas de atención integral en Comisaría de Familia | `demora_atencion_integral_comisaria` |
| 4 | Incumplimiento de medidas de protección sin respuesta institucional | `incumplimiento_medidas_proteccion` |
| 5 | Denegación del derecho a la no confrontación | `denegacion_no_confrontacion` |
| 6 | Divulgación o filtración de información confidencial | `filtracion_info_confidencial` |
| 7 | Medidas de protección ineficaces o mal implementadas | `medidas_ineficaces` |
| 8 | Falta de seguimiento institucional a las medidas de protección | `falta_seguimiento_medidas` |
| 9 | Fallas en la valoración y actualización del riesgo | `fallas_valoracion_riesgo` |
| 10 | Medidas de protección insuficientes acordes a la situación de riesgo | `medidas_insuficientes` |
| 11 | Omisión de restricción de visitas ante riesgo de feminicidio | `omision_restriccion_visitas` |
| 12 | Omisión de apoyo policial en residencia o lugar de trabajo | `omision_apoyo_policial` |
| 13 | Otras barreras en protección | `otras_barreras_proteccion` |

#### 9. ¿A qué institución acudió? (Protección)
| Campo | Valor |
|---|---|
| ID actual | `78474c82-61b9-4a0c-beca-439274813c03` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | Q1 = `proteccion` (EQUALS) |

| # | Label | Value |
|---|---|---|
| 1 | Comisarías de Familia | `comisaria_familia` |
| 2 | Fiscalía General de la Nación (funciones de protección) | `fiscalia_proteccion` |
| 3 | ICBF (Instituto Colombiano de Bienestar Familiar) | `icbf` |
| 4 | Unidad Nacional de Protección | `unidad_nacional_proteccion` |
| 5 | Jueces Civiles o Promiscuos Municipales | `jueces_civiles` |
| 6 | Inspección de Policía | `inspeccion_policia` |
| 7 | Defensoría de Familia | `defensoria_familia` |
| 8 | Ejército Nacional | `ejercito_nacional` |
| 9 | Armada Colombiana | `armada_colombiana` |
| 10 | No sabe / No responde | `no_sabe` |

#### 10. Otra barrera en Protección ¿cuál?
| Campo | Valor |
|---|---|
| ID actual | `68f6bf06-a6a8-443f-9a1e-001148d4eb45` |
| Tipo | `text` |
| Requerido | ❌ |
| Condición | Q8 CONTAINS `otras_barreras_proteccion` |

---

### Bloque Otras instituciones
*Visible cuando Q1 = `otras_instituciones` (EQUALS)*

#### 11. Nombre de la institución donde se presentó la barrera
| Campo | Valor |
|---|---|
| ID actual | `a1573bc3-28f6-485e-8d8b-b3567db42ae3` |
| Tipo | `text` |
| Requerido | ❌ |
| Condición | Q1 = `otras_instituciones` \| EQUALS |

---

### Preguntas comunes
*Siempre visibles — aplican a toda entrada del repeater independientemente del sector*

#### 12. Departamento donde se presentó la barrera
| Campo | Valor |
|---|---|
| ID actual | `31c7f8ba-880e-4c9a-89f0-1a6a1e43b7b9` |
| Tipo | `dropdown` |
| Requerido | ✅ |
| Condición | — |
| Fuente | Tabla `security.department` |

> ⚠️ **Sin opciones en BD:** las opciones de esta pregunta NO se registran en la tabla `option`. El frontend las carga directamente desde `security.department`.

#### 13. Ciudad donde se presentó la barrera
| Campo | Valor |
|---|---|
| ID actual | `c6f2d54a-3c61-4ae7-b2bf-50a884031fa4` |
| Tipo | `dropdown` |
| Requerido | ✅ |
| Condición | — |
| Fuente | Tabla `security.city` filtrada por Q12 |

> ⚠️ **Sin opciones en BD:** las opciones de esta pregunta NO se registran en la tabla `option`. El frontend las carga desde `security.city` filtrando por el departamento seleccionado en Q12.

#### 14. Municipio donde se presentó la barrera
| Campo | Valor |
|---|---|
| ID actual | `8225d03f-8de9-4ff0-9345-67b5e71bf02d` |
| Tipo | `dropdown` |
| Requerido | ✅ |
| Condición | — |
| Fuente | Tabla `security.town` filtrada por Q13 |

> ⚠️ **Sin opciones en BD:** las opciones de esta pregunta NO se registran en la tabla `option`. El frontend las carga desde `security.town` filtrando por la ciudad seleccionada en Q13.

#### 15. Barreras institucionales y de talento humano
| Campo | Valor |
|---|---|
| ID actual | `64754fb5-04a8-43e4-ac86-d60f0a84010b` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Falta de aplicación del enfoque de género / Revictimización o violencia institucional | `falta_enfoque_genero` |
| 2 | Falta de aplicación o desactualización de protocolos de atención en VBG | `falta_protocolos_vbg` |
| 3 | Capacidad institucional limitada | `capacidad_institucional_limitada` |
| 4 | Desarticulación institucional y ausencia de gestión integral de casos | `desarticulacion_institucional` |
| 5 | Insuficiencia en la orientación en derechos a las víctimas | `insuficiencia_orientacion_derechos` |

#### 16. Barreras económicas y socioeconómicas
| Campo | Valor |
|---|---|
| ID actual | `94dfc417-b59f-4afe-b63b-c6839a8dfecc` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Falta de autonomía económica | `falta_autonomia_economica` |
| 2 | Pobreza y desigualdad estructural | `pobreza_desigualdad` |

#### 17. Barreras territoriales y geográficas
| Campo | Valor |
|---|---|
| ID actual | `cf162686-0327-4b28-9639-e4d947272483` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Ausencia institucional, infraestructura y conectividad insuficiente o instalaciones distantes | `ausencia_institucional_territorial` |

#### 18. Barreras por ausencia de enfoque diferencial
| Campo | Valor |
|---|---|
| ID actual | `81275638-2837-4ecc-8686-675dfeab4f32` |
| Tipo | `multiple` |
| Requerido | ❌ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Barreras para personas con orientaciones sexuales e identidades de género diversas (OSIGD) | `barrera_osigd` |
| 2 | Barreras para personas con diagnósticos o condiciones de salud mental | `barrera_salud_mental` |
| 3 | Barreras para población con origen étnico (indígena, afrodescendiente, raizal, palenquera) | `barrera_etnica` |
| 4 | Barreras asociadas al ciclo de vida (niñez, adolescencia, adultez mayor) | `barrera_ciclo_vida` |
| 5 | Barreras para personas con discapacidad física o cognitiva | `barrera_discapacidad` |
| 6 | Barreras para personas migrantes y refugiadas | `barrera_migrantes` |
| 7 | Barreras para personas privadas de la libertad | `barrera_privadas_libertad` |
| 8 | Barreras para personas en situación de calle | `barrera_situacion_calle` |
| 9 | Barreras para víctimas de trata de personas | `barrera_trata` |
| 10 | Barreras para personas en actividades sexuales pagas | `barrera_actividades_sexuales_pagas` |
| 11 | Barreras para lideresas y defensoras de derechos humanos | `barrera_lideresas` |
| 12 | Barreras para víctimas del conflicto armado | `barrera_conflicto_armado` |

#### 19. Fecha en la que se presentó la barrera (aproximada)
| Campo | Valor |
|---|---|
| ID actual | `4592d85f-8c11-4785-ad32-06aa810cc491` |
| Tipo | `date` |
| Requerido | ✅ |
| Condición | — |

#### 20. Funcionario/a o dependencia donde se presentó la barrera
| Campo | Valor |
|---|---|
| ID actual | `33c3961e-4252-4b02-b491-615b11d6bf56` |
| Tipo | `text` |
| Requerido | ✅ |
| Condición | — |

#### 21. Descripción de la barrera
| Campo | Valor |
|---|---|
| ID actual | `d2be610f-44eb-4526-b4ec-86eaaaba08c8` |
| Tipo | `text` |
| Requerido | ✅ |
| Condición | — |

#### 22. Gestión de la barrera
| Campo | Valor |
|---|---|
| ID actual | `572ad72a-8174-4ff3-9c56-5c8c65ac63ac` |
| Tipo | `multiple` |
| Requerido | ✅ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Orientación y enrutamiento - Llamada | `orientacion_llamada` |
| 2 | Gestión administrativa - Llamada | `gestion_llamada` |
| 3 | Activación de ruta interinstitucional | `activacion_ruta_interinstitucional` |
| 4 | Articulación institucional | `articulacion_institucional` |
| 5 | Escalamiento a organismo de control | `escalamiento_organismo_control` |
| 6 | Alerta por barreras | `alerta_barreras` |

---

## Resumen de VCs del repeater

| Target (Q#) | Trigger (Q#) | Trigger value | Operator |
|---|---|---|---|
| Q2 Barreras Salud | Q1 Sector | `salud` | EQUALS |
| Q3 Institución Salud | Q1 Sector | `salud` | EQUALS |
| Q4 Otra barrera Salud | Q2 Barreras Salud | `otras_barreras_salud` | CONTAINS |
| Q5 Barreras Justicia | Q1 Sector | `justicia` | EQUALS |
| Q6 Institución Justicia | Q1 Sector | `justicia` | EQUALS |
| Q7 Otra barrera Justicia | Q5 Barreras Justicia | `otras_barreras_justicia` | CONTAINS |
| Q8 Barreras Protección | Q1 Sector | `proteccion` | EQUALS |
| Q9 Institución Protección | Q1 Sector | `proteccion` | EQUALS |
| Q10 Otra barrera Protección | Q8 Barreras Protección | `otras_barreras_proteccion` | CONTAINS |
| Q11 Nombre institución | Q1 Sector | `otras_instituciones` | EQUALS |

---

## Pendientes / Decisiones abiertas

- [ ] Confirmar si el `dropdown` de sector (Q1) agrega más sectores en el futuro (ej. Educación)
- [ ] Definir si las preguntas estructurales (Q15–Q18) se muestran juntas bajo un título o subtítulo tipo `info`
