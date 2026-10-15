# Reporte del formulario actual de Hacer Seguimiento

Fotografía de lo que el formulario **hace hoy**, tomada de la base de pruebas (Supabase) y del código que la interpreta.

| Dato | Valor |
|---|---|
| Formulario | Formulario de Seguimiento |
| Form ID | `2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff` |
| Pantalla | `/salvia/hacer-seguimiento/:id` |
| Motor | `dinamic-form` (`src/frontend/js/components/dinamic-form.js`) |
| Procesamiento al completar | `processFollowUpSubmission` en `src/salvia/service/form_service.go` |
| Secciones vivas | 7 |
| Preguntas vivas | 87 (15 info + el resto respondibles) |
| Repeaters vivos | 5 |

`form-seguimiento-data.md` describe una versión anterior (sección “Rutas Diferenciales” y “Cierre del caso”). Esa sección sigue en la base con `deleted_at` y **no se renderiza**. Este reporte es el que corresponde al formulario que corre ahora.

---

## Cómo se guarda

Hay dos momentos distintos.

**1. Al pasar de sección (“Siguiente”).** `dinamic-form` valida lo visible y obligatorio de la sección y persiste respuestas en:

- `salvia.answer` (una fila por pregunta; los multiple se guardan como CSV)
- `salvia.repeater_entry` + `salvia.answer.repeater_entry_id` para cada ítem de un repeater

Eso no crea barreras, entidades, remisiones ni cierra el caso.

**2. Al completar todas las secciones visibles.** El backend lanza `processFollowUpSubmission`. Si el seguimiento **ya está REALIZADO**, solo crea un evento de timeline “Seguimiento Editado” y termina: no vuelve a crear barreras, entidades, remisiones ni reasigna. Si todavía no está realizado, ejecuta los efectos de la sección [Qué crea al completar](#qué-crea-al-completar).

Las preguntas `info` no se responden y no bloquean el avance.

---

## Cómo se decide qué se ve

Una pregunta, repeater o sección se muestra si **todas** sus condiciones aplicables se cumplen (`EQUALS` o `CONTAINS`). Sin condiciones, se muestra siempre.

Hay dos orígenes de condición:

- **Respuesta de una pregunta anterior** (`trigger_question_id`).
- **Estado de la pantalla** (`trigger_state_path` sobre `formState`), que la pantalla escribe en `hacer_seguimiento.html` mientras el agente responde.

Los dropdowns de departamento, ciudad, municipio, entidad y obligaciones **no tienen opciones en la tabla `option`**. Las arma la pantalla y el formulario las lee por `state_options_path`.

Textos que cambian en pantalla (no son preguntas):

| Pregunta info | Qué muestra |
|---|---|
| “Reasignación:” | Sustituye ese texto por `formState.reassingText` |
| “Estado de la Remisión:” | Sustituye “Remisión:” por `formState.psysocialRemisionState` (`Remisión: SI cumple` / `NO cumple`) |
| “Seguimiento a Barrera” | Se reemplaza por el nombre de la barrera activa de esa posición |
| “Entidad en seguimiento” | Se reemplaza por el nombre de la entidad del caso de esa posición |

Los `metadata.alert` y `metadata.hint` de algunas preguntas **están en la base y no se pintan**. El componente no lee `metadata`.

---

## Sección 1 — Valoración del Riesgo

Siempre visible. 15 preguntas, sin repeater.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 0 | info | | Aviso de que los factores no son preguntas directas a la víctima y que solo se registran hechos posteriores al registro inicial | Siempre |
| 1 | single | sí | Respondiente del seguimiento | Siempre |
| 2 | boolean | sí | ¿Se registraron nuevos hechos de violencia? | Siempre |
| 3 | text | sí | Nuevos hechos: tiempo, modo, lugar y relación con la persona agresora | Q2 = Sí |
| 4 | date | sí | Fecha en que ocurrieron los hechos | Q2 = Sí |
| 5 | multiple | no | Factores protectores | Siempre |
| 6 | multiple | no | Factores de riesgo | Siempre |
| 7 | boolean | sí | ¿Las acciones desplegadas han tenido efecto protector tangible? | Siempre |
| 8 | multiple | no | Factores de riesgo extremo | Siempre |
| 9 | info | | Reasignación (texto dinámico) | `formState.shouldReassignCase` = `true` |
| 10 | boolean | sí | ¿Confirmar reasignación del caso a riesgo alto? | `formState.canReassignHigh` = `true` |
| 11 | boolean | sí | ¿Confirmar reasignación del caso a riesgo bajo? | `formState.canReassignLow` = `true` |
| 12 | info | | Invita a describir otro factor en la pregunta siguiente | Siempre |
| 13 | text | sí | Situación de riesgo actual y cómo los factores la aumentan, disminuyen o mantienen | Siempre |
| 14 | boolean | no | ¿Se identificaron barreras institucionales que mantienen o incrementan el riesgo? | Siempre. Si es Sí, abre la sección 4 |

**Respondiente:** `victima` Contacto con la víctima · `indirecto` Contacto indirecto (familiar o conocida) · `institucion` Contacto con institución.

**Factores protectores:** `medidas_proteccion`, `lugar_seguro`, `separacion_fisica`, `redes_apoyo`, `acompanamiento_psicosocial`, `autonomia_economica`, `ninguno`.

**Factores de riesgo:** `condiciones_habitacionales`, `agresor_grupo_organizado`, `agresor_sustancias`, `miedo_denunciar`, `bajo_reconocimiento_riesgo`, `embarazo_riesgo`, `discapacidad`, `residencia_rural`, `agresor_poder_institucional`, `ninguno`.

**Riesgo extremo:** `amenazas_muerte`, `tentativa_feminicidio`, `contacto_sexual`, `aislamiento`, `acceso_armas`, `destruccion_bienes`, `vigilancia_control`, `escalamiento_violencia`, `ninguno`.

La opción `ninguno` no cuenta para la reasignación.

**Qué calcula la pantalla al guardar la sección** (aún no escribe el caso):

- Nivel actual 1 o 2, y al menos un extremo (distinto de `ninguno`): muestra el aviso de reasignación automática a riesgo alto. No pide confirmación.
- Nivel actual 1 o 2, sin extremos, y 4 o más riesgos: muestra el aviso y la pregunta de confirmar riesgo alto.
- Nivel actual 3 o más, sin extremos, y 3 o más protectores: muestra el aviso y la pregunta de confirmar riesgo bajo.
- En cualquier otro caso esas tres preguntas quedan ocultas.

Q14 = Sí es la única llave de la sección 4. El hint “Registre estas en el módulo de barreras” está en metadata y no se muestra.

**Qué crea al completar el formulario**

| Respuesta | Registro |
|---|---|
| Q2 = Sí | Evento de timeline categoría General, tipo “Hechos del caso”, con la descripción (Q3) y la fecha (Q4; si no parsea, usa ahora) |
| Extremo en caso de nivel 1–2 | Sube `victim_case_form2.risk_level` a 4, regenera seguimientos pendientes del equipo Riesgo alto, actualiza equipo y agente del caso, evento “Reasignación de Caso”. No usa la confirmación |
| Confirmación de riesgo alto = Sí y 4 o más riesgos, en caso de nivel 1–2 | Igual, con nivel 3 |
| Confirmación de riesgo bajo = Sí, sin extremos y 3 o más protectores, en caso de nivel 3+ | Nivel 2, reasigna el agente de los pendientes al equipo Riesgo bajo (no regenera el calendario), actualiza equipo y agente, evento “Reasignación de Caso” |
| Resto de esta sección | Solo queda en `answer` |

---

## Sección 2 — Identificación de Entidades

Siempre visible. Tres preguntas directas y dos repeaters. El agente agrega filas con “+ Agregar Entidad”. Mínimo 0.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 1 | boolean | sí | ¿Ha puesto en conocimiento de alguna entidad los hechos? | Siempre |
| 3 | single | sí | ¿Está de acuerdo con que SALVIA active la Ruta de Atención? (`si` / `no`) | Siempre |
| 4 | boolean | sí | ¿Salvia debe activar rutas en el marco de la debida diligencia? | Siempre |

La pregunta de acuerdo **no condiciona** a la de debida diligencia: las tres se ven siempre. El acuerdo tampoco crea un registro propio.

### Repeater “Entidades a las que ya acudió”

Visible solo si “¿Ha puesto en conocimiento?” = Sí.

| # | Tipo | Obl. | Pregunta | Opciones |
|---|---|---|---|---|
| 1 | dropdown | sí | Sector | `salud`, `justicia`, `proteccion`, `otras_instituciones` |
| 2 | dropdown | sí | Departamento | `formState.statesColombia` |
| 3 | dropdown | sí | Ciudad | ciudades del departamento elegido |
| 4 | dropdown | sí | Municipio | municipios de la ciudad (`/api/v1/locations/towns`) |
| 5 | dropdown | sí | Entidad | sedes de `/api/v1/entity-branches?town_code=&sector=`. El valor guardado es el id de la sede |
| 6 | text | no | N° de radicado, NUNC, MP, etc. | |
| 7 | date | sí | Fecha de reporte | |
| 8 | multiple | no | Acciones institucionales completadas | catálogo `/api/v1/entity-branches/{id}/obligations` (incluye “Otra”) |
| 9 | text | sí | Describe el otro ítem completado | solo si Q8 contiene `otra` |
| 10 | multiple | no | Acciones institucionales pendientes | el mismo catálogo |
| 11 | text | sí | Describe el otro ítem pendiente | solo si Q10 contiene `otra` |
| 12 | text | sí | Otra información asociada a la Ruta | |

### Repeater “Entidades para activación de ruta”

Visible solo si “¿Salvia debe activar rutas…?” = Sí.

| # | Tipo | Obl. | Pregunta | Opciones |
|---|---|---|---|---|
| 1 | dropdown | sí | Sector | igual que el repeater anterior |
| 2–5 | dropdown | sí | Departamento, Ciudad, Municipio, Entidad | la misma cascada, en `formState.newEntidadesActivacion` |
| 6 | multiple | no | Acciones institucionales pendientes | catálogo de obligaciones de la sede |
| 7 | text | sí | Describe el otro ítem pendiente | solo si Q6 contiene `otra` |
| 8 | multiple | sí | Canal de activación | `notificacion` Notificación (tarea) · `llamada` Llamada (tarea) · `mi_salvia` Mi Salvia (tarea) |

**Qué crea al completar**, una fila por entidad elegida. Si la sede ya está en el caso, reutiliza ese `entity_case` y no duplica.

| Repeater | Registros |
|---|---|
| Ya acudió | `entity_case` (caso, sede, radicado, objetivo = “Otra información asociada a la Ruta”). Por cada obligación completada: `entity_case_obligation` en estado completada. Por cada pendiente: la misma tabla en estado pendiente. Si el valor es `otra`, guarda el texto libre en `custom_label` en lugar del id de catálogo. La fecha de reporte solo queda en `answer` |
| Activación | Igual, sin radicado ni objetivo, y solo obligaciones pendientes. Además, por cada canal: ver [Canal de activación](#canal-de-activación) |

---

## Sección 3 — Seguimiento a Entidades

La sección solo aparece si el caso **ya tiene** entidades (`formState.hasCurrentEntities` = `true`). Esas entidades se cargan al abrir la pantalla.

Repeater “Entidades en seguimiento”. No se agregan ni se borran filas: hay una por entidad existente, en el mismo orden. El banner info muestra el nombre de la entidad.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 1 | info | | Nombre de la entidad | Siempre, dentro de la fila |
| 2 | multiple | no | Acciones institucionales que ya se completaron | Opciones = obligaciones **pendientes** de esa entidad. El valor es el id de la fila `entity_case_obligation`, no el del catálogo |
| 3 | boolean | no | ¿La Ruta se actualizó? | Siempre |
| 4 | text | sí | Por qué se actualizó la Ruta | Q3 = Sí |
| 5 | boolean | no | ¿Se requiere nueva activación? | Siempre |
| 6 | multiple | sí | Canal de activación | Q5 = Sí. Mismas tres opciones: `notificacion`, `llamada`, `mi_salvia` |

**Qué crea o actualiza al completar**, emparejando fila N con la entidad N del caso:

| Respuesta | Efecto |
|---|---|
| Cualquier fila | `entity_case_follow_up` (ruta actualizada, motivo, si requiere nueva activación, canal) y actualiza `entity_case.last_action` |
| Obligaciones marcadas como completadas | Esas filas de `entity_case_obligation` pasan a completada. No crea filas nuevas |
| ¿Se requiere nueva activación? = Sí | Una tarea (y oficio, si el canal es notificación) por cada canal. Ver abajo |

---

## Canal de activación

Lo usan el repeater de activación de la sección 2 y la nueva activación de la sección 3. Un canal por valor seleccionado. La tarea queda asignada al agente que llena el seguimiento, categoría `Entidades`, estado `ToDo`.

| Valor | Qué crea |
|---|---|
| `notificacion` | `entity_letter` en estado `por_proyectar` (ligado a la sede) y `case_task` tipo `proyectar_oficio`, descripción “Proyectar oficio de activación de ruta” |
| `llamada` | `case_task` tipo `gestion_llamada`, descripción “Llamar a la entidad para gestionar la activación” |
| `mi_salvia` | `case_task` tipo `solicitud_entidad`, descripción “Hacer solicitudes a la entidad sobre obligaciones pendientes” |

---

## Sección 4 — Identificación de Barreras

Visible solo si en la sección 1 la pregunta de barreras institucionales = Sí.

Repeater “Barreras identificadas”. Mínimo 1. Botón “+ Agregar Barrera”. Los multiple de **este** repeater se pintan como dropdown con chips, no como botones.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 1 | dropdown | sí | Sector de la barrera | Siempre en la fila |
| 2 | multiple | no | Barreras identificadas en Salud (12) | Sector = `salud` |
| 4 | text | no | Otra barrera en Salud ¿cuál? | El multiple de salud contiene `otras_barreras_salud` |
| 5 | multiple | no | Barreras identificadas en Justicia (18) | Sector = `justicia` |
| 7 | text | no | Otra barrera en Justicia ¿cuál? | Contiene `otras_barreras_justicia` |
| 8 | multiple | no | Barreras identificadas en Protección (13) | Sector = `proteccion` |
| 10 | text | no | Otra barrera en Protección ¿cuál? | Contiene `otras_barreras_proteccion` |
| 12 | dropdown | no | Entidad a la que se asignará la barrera | Siempre. Opciones = entidades que el caso ya tenía al abrir la pantalla (`currentEntitiesOptions`). No incluye las que se acaban de escribir en la sección 2 de esta misma sesión |
| 16 | multiple | no | Barreras institucionales y de talento humano (5) | Siempre |
| 17 | multiple | no | Barreras económicas y socioeconómicas (2) | Siempre |
| 18 | multiple | no | Barreras territoriales y geográficas (1) | Siempre |
| 19 | multiple | no | Barreras por ausencia de enfoque diferencial (12) | Siempre |
| 20 | date | sí | Fecha aproximada | Siempre |
| 21 | text | sí | Funcionario/a o dependencia | Siempre |
| 22 | text | sí | Descripción de la barrera | Siempre |
| 23 | multiple | sí | Gestión de la barrera | Siempre |

**Sector:** `salud`, `justicia`, `proteccion`, `otras_instituciones`, `barrera_transversal`.

En `otras_instituciones` y `barrera_transversal` no aparece el bloque de barreras específicas. La pregunta de “nombre de institución” y las de departamento, ciudad y municipio de la barrera están borradas: ya no se muestran. El código de la pantalla todavía intenta armar ciudades, municipios y competencias a partir de esas preguntas borradas; no hay pregunta viva que las consuma.

**Gestión de la barrera**

| Valor | Label | Qué crea |
|---|---|---|
| `orientacion_llamada` | Orientación y enrutamiento - Llamada | Nada además de la barrera |
| `gestion_llamada` | Gestión administrativa - Llamada | `case_task` tipo `gestion_llamada`, asignada al agente |
| `alerta_barreras` | Alerta por barreras | `case_task` tipo `comite_caso`, **sin** agente asignado |
| `activacion_ruta_interinstitucional` | Activación de ruta interinstitucional | `entity_letter` `por_proyectar` ligado a la barrera + `case_task` tipo `proyectar_oficio` |
| `articulacion_institucional` | Articulación institucional | Igual: oficio + tarea `proyectar_oficio` |
| `escalamiento_organismo_control` | Escalamiento a organismo de control | Igual: oficio + tarea `proyectar_oficio` |

**Qué crea la fila, siempre:** un `barrier_v2` en estado `OPEN`, con sector, barreras específicas del sector, texto de “otra barrera”, las cuatro barreras estructurales, fecha, funcionario, descripción, gestión y, si eligió entidad, `entity_branch_id`. Categoría de las tareas: `Barreras`.

**Salud (12):** `dificultades_aseguramiento`, `negacion_urgencias`, `fallas_calidad_urgencias`, `demoras_citas`, `falta_protocolo_violencia_sexual`, `falta_protocolo_agentes_quimicos`, `demoras_valoracion_medicolegal`, `negacion_ive`, `exigencias_autorizacion_discapacidad`, `falta_enfoque_cultural_salud`, `negacion_procedimientos_osigd`, `otras_barreras_salud`.

**Justicia (18):** `negativa_recibir_denuncia`, `dilacion_registro_denuncia`, `denegacion_no_confrontacion`, `falta_info_documental`, `formalismo_excesivo`, `cargas_probatorias_excesivas`, `falta_celeridad_investigacion`, `incumplimiento_ordenes_judiciales`, `tipificacion_erronea`, `traslado_inadecuado_expediente`, `dificultad_reasignacion_municipio`, `falta_representacion_judicial`, `vacios_normativos_vbg`, `riesgo_vencimiento_terminos`, `medidas_libertad_riesgo`, `falta_proteccion_info_confidencial`, `deficiencias_notificacion_agresor`, `otras_barreras_justicia`.

**Protección (13):** `demoras_medidas_comisaria`, `dilacion_medidas_fisicas`, `demora_atencion_integral_comisaria`, `incumplimiento_medidas_proteccion`, `denegacion_no_confrontacion`, `filtracion_info_confidencial`, `medidas_ineficaces`, `falta_seguimiento_medidas`, `fallas_valoracion_riesgo`, `medidas_insuficientes`, `omision_restriccion_visitas`, `omision_apoyo_policial`, `otras_barreras_proteccion`.

**Institucionales y de talento humano:** `falta_enfoque_genero`, `falta_protocolos_vbg`, `capacidad_institucional_limitada`, `desarticulacion_institucional`, `insuficiencia_orientacion_derechos`.

**Económicas:** `falta_autonomia_economica`, `pobreza_desigualdad`.

**Territoriales:** `ausencia_institucional_territorial`.

**Enfoque diferencial:** `barrera_osigd`, `barrera_salud_mental`, `barrera_etnica`, `barrera_ciclo_vida`, `barrera_discapacidad`, `barrera_migrantes`, `barrera_privadas_libertad`, `barrera_situacion_calle`, `barrera_trata`, `barrera_actividades_sexuales_pagas`, `barrera_lideresas`, `barrera_conflicto_armado`.

---

## Sección 5 — Seguimiento a Barreras

La sección **siempre está en el menú**. El repeater “Seguimiento a Barreras Activas” no deja agregar filas: trae una por barrera activa del caso (`formState.currentBarriers`, las que no están en `MANAGED`). Si el caso no tiene, la sección se ve vacía y se puede continuar. El botón de la base dice “Siguiente Barrera”, pero al ser un repeater de estado ese botón no se pinta.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 1 | info | | Nombre de la barrera | Siempre en la fila |
| 2 | boolean | sí | ¿Persiste la barrera? | Siempre |
| 3 | single | sí | ¿Hubo respuesta institucional? | Siempre |
| 4 | multiple | sí | Gestión de la barrera | Siempre. Mismas 6 opciones de la sección 4 |
| 5 | text | sí | Actuaciones realizadas | Siempre |
| 6 | boolean | sí | ¿Se realiza cierre de la barrera? | Siempre |
| 7 | single | sí | Motivo del cierre | Q6 = Sí |

**Respuesta institucional:** `respuesta_oficial` La entidad respondió de manera oficial a Salvia · `contacto_victima` La entidad se contactó con la víctima · `sin_respuesta` No se obtuvo respuesta.

**Motivo del cierre:** `no_voluntad`, `no_gestionable`, `clasificacion_incorrecta`, `instalada_entidad`, `resuelta`, `superada_parcialmente`.

La fila N se ata a la barrera activa N (`follow_up_v2.active_barrier_ids`).

**Qué crea al completar**

| Respuesta | Efecto |
|---|---|
| Cada fila con barrera emparejada | `barrier_follow_up` (persiste, respuesta institucional, gestión, actuaciones, si cierra, motivo) |
| Cada fila | Evento de timeline categoría Barreras, tipo “Seguimiento a Barrera”. El texto junta si persiste, la respuesta institucional y las actuaciones |
| Q6 = Sí | Esa `barrier_v2` pasa a `MANAGED` |

La gestión de esta sección **no** crea tareas ni oficios. Eso solo ocurre con las barreras nuevas de la sección 4.

---

## Sección 6 — Seguimiento de Caso

Siempre visible.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 1 | text | sí | Acciones relevantes del seguimiento (riesgo, ruta, orientación) | Siempre |
| 2 | multiple | no | Gestiones interinstitucionales a partir de los hallazgos | Siempre |
| 3 | text | no | Acciones previstas, entidades y propósito | Siempre (no depende de Q2) |
| 4 | boolean | sí | ¿Hay necesidades que requieran remisión a equipos SALVIA? | Siempre |
| 5 | multiple | sí | ¿Cuáles equipos? | Q4 = Sí |
| 6 | info | | Aviso de que Medidas de Emergencia y Atención Psicosocial no se remiten juntas; solo se remite a Medidas | Equipos contiene `medidas_emergencia` **y** `atencion_psico` |
| 7 | multiple | no | ¿Cuáles medidas de emergencia? | Equipos contiene `medidas_emergencia` |
| 8 | info | | Regla de puntaje y estado SI/NO cumple | Equipos contiene `atencion_psico` |
| 9 | multiple | no | Criterios de remisión — Atención Psicosocial | Equipos contiene `atencion_psico` |
| 10 | multiple | sí | Criterio de Atención Hombres | Equipos contiene `atencion_hombres` |
| 11 | multiple | sí | Servicios del Equipo de Discapacidad | Equipos contiene `discapacidad` |
| 12 | multiple | sí | Criterios de Estabilización (al menos 1) | Equipos contiene `estabilizacion` |
| 13 | text | sí | Elementos que evidencian la remisión | Q4 = Sí |
| 14 | boolean | sí | ¿Se realiza reasignación con otro equipo o profesional? | Siempre |
| 15 | single | sí | Equipo al que reasigna el caso | Q14 = Sí |

**Gestiones interinstitucionales (Q2):** `orientacion_enrutamiento`, `gestion_administrativa`, `activacion_ruta`, `articulacion_institucional`, `escalamiento_control`. Q2 y Q3 solo quedan en `answer`. No crean tareas.

**Equipos (Q5):** `medidas_emergencia`, `atencion_psico`, `estabilizacion`, `atencion_hombres`, `salvia_dignidad`, `discapacidad`.

**Medidas:** `alojamiento`, `transporte`, `alimentacion`, `vestuario`, `apoyo_psico`, `otras_me`.

**Criterios psicosociales y puntaje** (el banner y el backend usan la misma regla):

| Valor | Puntos |
|---|---|
| `criterio_obligatorio` | obligatorio, 0 puntos |
| `conducta_suicida` | 3 |
| `interseccionalidad` | 2 |
| `sin_ruta`, `condiciones_territoriales`, `sin_acceso_psico`, `naturalizacion_vbg` | 1 cada uno |

La remisión psicosocial se crea solo si está el criterio obligatorio **y** la suma es al menos 3.

**Hombres:** una sola opción, `criterio_hombres` (“Hombres diversos victimarios o posibles victimarios”). Sin esa opción no se crea la remisión. El hint de población está en metadata y no se muestra.

**Discapacidad:** `apoyo_lsc`, `enfoque_discapacidad`. Un registro por servicio.

**Estabilización:** `sin_ingresos`, `vivienda_temporal`, `personas_cargo`. Basta uno.

**Equipo de reasignación manual (Q15):** `riesgo`, `general`, `hombres`, `salvia_dignidad`. Esta pareja de preguntas **no reasigna el caso**. La reasignación real es la de la sección 1 (factores y confirmaciones).

**Qué crea al completar**, si Q4 = Sí y el equipo está marcado:

| Equipo | Condición extra | Registro |
|---|---|---|
| `atencion_psico` | No puede ir junto con `medidas_emergencia`. Obligatorio + mínimo 3 puntos | `psychosocial_support` tipo `derivacion`, estado `abierto`, con remitente, equipo del remitente, criterios y el texto de evidencia. Más `case_task` categoría Psicosocial, tipo `validar_remision`, sin agente |
| `atencion_hombres` | Debe estar `criterio_hombres` | `men_team_remision` tipo `derivacion`, estado `ACTIVE` |
| `discapacidad` | Al menos un servicio | Un `discapacidad_remision` por servicio, estado `ACTIVE` |
| `estabilizacion` | Al menos un criterio | `economic_stabilization` tipo `derivacion`, estado `ACTIVE` |
| `salvia_dignidad` | — | Nada. La opción se guarda y el backend no tiene caso para ella |
| `medidas_emergencia` | — | Un `emergency_measure` por medida, estado `ACTIVE` |

Si se marcan medidas y psicosocial a la vez, se crean las medidas y se omite la remisión psicosocial.

---

## Sección 7 — Conclusión del ciclo de atención

Siempre visible.

| # | Tipo | Obl. | Pregunta | Cuándo se ve |
|---|---|---|---|---|
| 1 | boolean | sí | ¿Realiza conclusión del ciclo de atención? | Siempre |
| 2 | single | no | Motivo | Q1 = Sí |
| 3 | text | no | Otro motivo ¿cuál? | Motivo = `otro` |
| 4 | text | no | Razones que justifican la conclusión | Q1 = Sí |
| 5 | boolean | no | ¿Realizó acciones institucionales a partir de la conclusión? | Q1 = Sí |

**Motivos:** `perdida_contacto`, `solicitud_expresa`, `cumplimiento_plan`, `no_corresponde` (la información no permite clasificar el caso como VBG o VPP), `otro`.

La pregunta de detallar esas gestiones institucionales está borrada. Q5 se guarda y **no dispara** tareas ni oficios.

**Qué crea si Q1 = Sí:** el caso pasa a estado `cd`. Evento de timeline “Cierre de Caso” con motivo (si es `otro`, usa el texto libre) y la descripción. Si el caso ya estaba cerrado, no hace nada. Un fallo de este paso no deshace el seguimiento ya marcado REALIZADO.

---

## Qué crea al completar

Siempre, la primera vez que el formulario queda completo:

1. `follow_up_v2.status` = `REALIZADO` y `completed_at`.
2. `follow_up_attempt` con `was_answered = true` y motivo “Sí contestó”.
3. Evento de timeline “Seguimiento Ejecutado”. El texto cuenta cuántas barreras nuevas y cuántas remisiones a equipos SALVIA se crearon (las medidas cuentan dentro de ese número; las entidades no).

Además, según las respuestas:

| Si el agente… | Se crea o se actualiza |
|---|---|
| Reporta nuevos hechos | Timeline “Hechos del caso” |
| Califica para reasignación por factores y, cuando aplica, confirma | Nivel de riesgo, calendario o agente, equipo del caso, timeline “Reasignación de Caso” |
| Registra entidades a las que ya acudió o para activar | `entity_case` y `entity_case_obligation` |
| Elige canal de activación | `case_task` de entidades y, si es notificación, `entity_letter` |
| Da seguimiento a entidades ya del caso | `entity_case_follow_up`, `last_action`, obligaciones marcadas como completadas y, si pide nueva activación, tareas |
| Identifica barreras | `barrier_v2` y, según la gestión, `case_task` y `entity_letter` |
| Da seguimiento a barreras activas | `barrier_follow_up`, timeline “Seguimiento a Barrera” y, si cierra, `barrier_v2.status = MANAGED` |
| Remite a un equipo que cumple su regla | La tabla de ese equipo (psicosocial, hombres, discapacidad, estabilización) y, en psicosocial, la tarea `validar_remision` |
| Elige medidas de emergencia | `emergency_measure` |
| Concluye el ciclo | Caso en `cd` y timeline “Cierre de Caso” |

**Respuestas que solo quedan en `answer`:** respondiente, efecto protector, texto analítico de riesgo, acuerdo con activar la ruta, fecha de reporte de la entidad, gestiones interinstitucionales de la sección 6 y su descripción, reasignación manual a otro equipo (Q14 y Q15), acciones institucionales de la conclusión (Q5), y todos los textos info.

**Si se vuelve a guardar un seguimiento ya REALIZADO:** solo el evento “Seguimiento Editado”. Las respuestas sí se actualizan en `answer`, pero no se reprocesan barreras, entidades, remisiones, cierre ni reasignación.
