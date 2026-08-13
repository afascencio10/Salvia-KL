-- =============================================================================
-- SEED: Formularios Psicosociales — 4 formularios independientes
--
-- Fuente: Sheet "Psicosocial Kreivo27.05.2026", hoja "Formularios Psicosocial"
-- Decisiones del líder (columna "Comentarios Felipe"):
--   - "Seleccionar número de dupla" ELIMINADA (Dupla asignada en sistema)
--   - "¿La llamada fue efectiva?" ELIMINADA de Primer Contacto
--
-- Ajustes Jul 2026:
--   - PC: Se agrega Sección 2 "Primera Atención" (oculta por defecto, gatillo = Q12 Continuar)
--   - PC S1: "Fecha próxima atención" es la ÚLTIMA pregunta (Q13), después de Q12 "Continuar
--     Primera Atención", para que su visibilidad pueda evaluarse al responder Q12.
--     Visible solo cuando Q5=Sí AND Q12=false (no continuar).
--   - CIE S2: Se agrega Q6 "Cerrar remisión" (boolean, gatillo de S3)
--   - CIE S3: Visibilidad ahora controlada por "Cerrar remisión = true" (no por "Es atención")
--
-- Ajustes Jul 2026 (registro de Barreras):
--   - Se agregan 2 secciones nuevas a los 4 formularios, ubicadas justo después de la
--     sección de Contacto: "Seguimiento a Barreras" e "Identificación de Barreras".
--     Ambas replican preguntas/funcionalidad del Formulario de Seguimiento (ver
--     DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md y
--     form-seguimiento-barreras-repeater.md).
--   - En la sección de Contacto de cada formulario se agrega, en la ÚLTIMA posición,
--     la pregunta "Desde la atención anterior se han identificado barreras
--     institucionales" (boolean). Gatillo de visibilidad de "Identificación de Barreras":
--       PC        → Q12 "Continuar Primera Atención" = true
--       PA/SEG/CIE→ Q8  "¿Es atención o solo contacto?" = atencion
--   - "Seguimiento a Barreras" NO depende de esta nueva pregunta: permanece habilitada
--     sin importar la respuesta. Al igual que en el Formulario de Seguimiento, su
--     visibilidad real depende de si el caso tiene barreras activas (formState externo,
--     ej. `currentBarriers`, poblado por el backend) — decisión pendiente de definir el
--     `trigger_state_path` exacto, igual que en form-seguimiento-barreras-repeater.md.
--     Por eso no se inserta un visibility_condition de sección para ella en este seed.
--
-- Ajustes Jul 2026 (ronda 2 — limpieza de preguntas y rebranding):
--   - PC S4 y PA S4 ("Primera Atención"): se ELIMINAN "¿Requiere ajuste razonable?" y
--     "¿Requiere intérprete de idiomas?" (antes Q2/Q3). "Confirmación consentimiento
--     persona de apoyo" (ahora Q3) ya NO depende de la pregunta de intérprete eliminada:
--     pasa a ser siempre visible. Las preguntas restantes se renumeran (14 → 12).
--   - SEG S1 y CIE S1 (Contacto): se agrega "¿La atención es individual o en dupla?"
--     como NUEVA primera pregunta (antes ausente en estos 2 formularios). El resto de
--     preguntas de la sección se recorre 1 posición (11 → 12; "Es atención" pasa de
--     Q8 a Q9, "Desde la atención anterior..." pasa de Q11 a Q12).
--   - Form 3 renombrado de "Seguimiento Psicosocial" a "Atención Psicosocial", junto con
--     sus secciones que contenían la palabra "Seguimiento" ("Contacto Seguimiento" →
--     "Contacto Atención Psicosocial", "Seguimiento" → "Atención Psicosocial"). En el
--     Form 4 (Cierre, que mantiene su nombre) se aplica el mismo cambio a su sección
--     "Seguimiento (Cierre)" → "Atención Psicosocial (Cierre)". "Seguimiento a Barreras"
--     NO se renombra en ningún formulario (nombre fijo de esa sección de barreras).
--
-- Ajustes Aug 2026 (diagrama):
--   - Consentimiento PC/PA: info (texto) + single (Sí/No); hide posteriores si No
--   - Agenda: ¿Agendar nueva sesión? + Fecha + Hora (UUIDs fijos en
--     internal/constants/psicosocial_agenda_questions.go)
--   - BD ya sembrada: aplicar también migrate_psicosocial_aug2026.sql
--
-- Formularios resultantes:
--   Form 1 — Primer Contacto     (4 secciones: S1=14 preguntas, S2=Seguimiento a Barreras,
--                                  S3=Identificación de Barreras, S4=12 preguntas condicional)
--   Form 2 — Primera Atención    (4 secciones: S1=11, S2=Seguimiento a Barreras,
--                                  S3=Identificación de Barreras, S4=12)
--   Form 3 — Atención Psicosocial (4 secciones: S1=12, S2=Seguimiento a Barreras,
--                                  S3=Identificación de Barreras, S4=5)
--   Form 4 — Cierre              (5 secciones: S1=12, S2=Seguimiento a Barreras,
--                                  S3=Identificación de Barreras, S4=6, S5=7)
--
-- Ejecutar:
--   PGPASSWORD='LegacySalvia2026@' psql \
--     -h aws-1-us-west-1.pooler.supabase.com \
--     -U "salvia_legacy.pwwelwfhauspznpatuqm" \
--     -d postgres -f seed_psicosocial.sql
-- =============================================================================

-- =============================================================================
-- FUNCIONES HELPER (temporales — schema pg_temp, se descartan solas al cerrar la
-- sesión de psql). Insertan el repeater_group + preguntas + opciones de cada
-- sección de barreras, dado el form_id y el form_section_id ya creados.
-- Estructura tomada de DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md
-- y form-seguimiento-barreras-repeater.md (misma estructura que el Formulario de
-- Seguimiento en producción).
-- =============================================================================

CREATE OR REPLACE FUNCTION pg_temp.seed_identificacion_barreras(p_form_id uuid, p_section_id uuid)
RETURNS varchar AS $func$
DECLARE
  v_rg   varchar;   -- repeater_group id
  v_q1   varchar;   -- Sector de la barrera
  v_q2   varchar;   -- Barreras identificadas en Salud
  v_q5   varchar;   -- Barreras identificadas en Justicia
  v_q8   varchar;   -- Barreras identificadas en Protección
  v_tmp  varchar;   -- reutilizada para preguntas que no se referencian después
BEGIN
  INSERT INTO salvia.repeater_group (id, form_section_id, name, item_name, "order", min_repetitions)
  VALUES (gen_random_uuid(), p_section_id, 'Barreras identificadas', 'Barrera', 1, 1)
  RETURNING id INTO v_rg;

  -- Q1: Sector de la barrera
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'dropdown', 'Sector de la barrera', TRUE, 1)
  RETURNING id INTO v_q1;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_q1, 'Salud',                'salud',                1),
    (gen_random_uuid(), v_q1, 'Justicia',             'justicia',             2),
    (gen_random_uuid(), v_q1, 'Protección',           'proteccion',           3),
    (gen_random_uuid(), v_q1, 'Otras instituciones',  'otras_instituciones',  4),
    (gen_random_uuid(), v_q1, 'Barrera Transversal',  'barrera_transversal',  5);

  -- ── Bloque Salud [visible: Q1 = salud] ──────────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras identificadas en Salud', FALSE, 2)
  RETURNING id INTO v_q2;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_q2, 'Dificultades de aseguramiento o afiliación en salud',                                                          'dificultades_aseguramiento',           1),
    (gen_random_uuid(), v_q2, 'Negaciones en los servicios de urgencias',                                                                      'negacion_urgencias',                   2),
    (gen_random_uuid(), v_q2, 'Fallas en la calidad del servicio de urgencias',                                                                'fallas_calidad_urgencias',             3),
    (gen_random_uuid(), v_q2, 'Demoras en la asignación de citas y continuidad de tratamientos',                                               'demoras_citas',                        4),
    (gen_random_uuid(), v_q2, 'Faltas de activación de protocolos para violencia sexual',                                                     'falta_protocolo_violencia_sexual',     5),
    (gen_random_uuid(), v_q2, 'Faltas de activación de protocolos para ataques con agentes químicos',                                         'falta_protocolo_agentes_quimicos',     6),
    (gen_random_uuid(), v_q2, 'Demoras y dificultades en la valoración medicolegal',                                                           'demoras_valoracion_medicolegal',       7),
    (gen_random_uuid(), v_q2, 'Negaciones en el acceso a la Interrupción Voluntaria del Embarazo (IVE)',                                      'negacion_ive',                         8),
    (gen_random_uuid(), v_q2, 'Exigencias indebidas de autorización para procedimientos en personas con discapacidad cognitiva y/o psicosocial','exigencias_autorizacion_discapacidad', 9),
    (gen_random_uuid(), v_q2, 'Atenciones en salud no centradas en cosmovisiones y prácticas culturales propias',                              'falta_enfoque_cultural_salud',        10),
    (gen_random_uuid(), v_q2, 'Negaciones o imposiciones de procedimientos a personas con diversidad sexual o identidad de género diversa',    'negacion_procedimientos_osigd',       11),
    (gen_random_uuid(), v_q2, 'Otras barreras en salud',                                                                                       'otras_barreras_salud',                12);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_q2, v_q1, 'salud', 'EQUALS');

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', '¿A qué institución acudió?', FALSE, 3)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Hospital Público',      'hospital_publico',   1),
    (gen_random_uuid(), v_tmp, 'Hospital Privado',      'hospital_privado',   2),
    (gen_random_uuid(), v_tmp, 'IPS',                   'ips',                3),
    (gen_random_uuid(), v_tmp, 'EPS',                   'eps',                4),
    (gen_random_uuid(), v_tmp, 'Puesto de Salud',       'puesto_salud',       5),
    (gen_random_uuid(), v_tmp, 'Consultorio médico',    'consultorio_medico', 6),
    (gen_random_uuid(), v_tmp, 'No sabe / No responde', 'no_sabe',            7);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q1, 'salud', 'EQUALS');

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text', 'Otra barrera en Salud ¿cuál?', FALSE, 4)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q2, 'otras_barreras_salud', 'CONTAINS');

  -- ── Bloque Justicia [visible: Q1 = justicia] ────────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras identificadas en Justicia', FALSE, 5)
  RETURNING id INTO v_q5;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_q5, 'Negativa institucional para recibir la denuncia',                                                             'negativa_recibir_denuncia',          1),
    (gen_random_uuid(), v_q5, 'Dilación en la recepción y registro de la denuncia',                                                          'dilacion_registro_denuncia',         2),
    (gen_random_uuid(), v_q5, 'Denegación al derecho a la no confrontación',                                                                 'denegacion_no_confrontacion',        3),
    (gen_random_uuid(), v_q5, 'Falta de entrega de información y documentación procesal',                                                    'falta_info_documental',              4),
    (gen_random_uuid(), v_q5, 'Formalismo excesivo en los trámites judiciales',                                                              'formalismo_excesivo',                5),
    (gen_random_uuid(), v_q5, 'Imposición de cargas probatorias injustificadas o excesivas',                                                 'cargas_probatorias_excesivas',       6),
    (gen_random_uuid(), v_q5, 'Falta de celeridad en la investigación judicial',                                                             'falta_celeridad_investigacion',      7),
    (gen_random_uuid(), v_q5, 'Incumplimiento de órdenes judiciales de protección o sanción',                                                'incumplimiento_ordenes_judiciales',  8),
    (gen_random_uuid(), v_q5, 'Tipificación errónea del delito',                                                                             'tipificacion_erronea',               9),
    (gen_random_uuid(), v_q5, 'Traslado inadecuado o incompleto del expediente judicial',                                                    'traslado_inadecuado_expediente',    10),
    (gen_random_uuid(), v_q5, 'Dificultades en la reasignación de procesos de Comisaría o Defensoría de Familia tras cambio de municipio',   'dificultad_reasignacion_municipio', 11),
    (gen_random_uuid(), v_q5, 'Ausencia o dificultad de acceso a representación judicial',                                                   'falta_representacion_judicial',     12),
    (gen_random_uuid(), v_q5, 'Vacíos normativos y definiciones restrictivas en VBG',                                                        'vacios_normativos_vbg',             13),
    (gen_random_uuid(), v_q5, 'Riesgo de vencimiento de términos',                                                                           'riesgo_vencimiento_terminos',       14),
    (gen_random_uuid(), v_q5, 'Medidas preventivas de libertad dirigidas al perpetrador que ponen en riesgo a la mujer y su núcleo familiar','medidas_libertad_riesgo',           15),
    (gen_random_uuid(), v_q5, 'Falta de protección en la divulgación de información confidencial del proceso',                              'falta_proteccion_info_confidencial',16),
    (gen_random_uuid(), v_q5, 'Deficiencias en la notificación y recolección de datos del agresor',                                          'deficiencias_notificacion_agresor', 17),
    (gen_random_uuid(), v_q5, 'Otras barreras en justicia',                                                                                  'otras_barreras_justicia',           18);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_q5, v_q1, 'justicia', 'EQUALS');

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', '¿A qué institución acudió?', FALSE, 6)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'CAI - Policía (Comando de Atención Inmediata)',                       'cai_policia',              1),
    (gen_random_uuid(), v_tmp, 'CAIVAS (Centro de Atención Integral a Víctimas de Violencia Sexual)', 'caivas',                   2),
    (gen_random_uuid(), v_tmp, 'Casas de Justicia',                                                   'casas_justicia',           3),
    (gen_random_uuid(), v_tmp, 'CAVIV (Centro de Atención Integral contra la Violencia Intrafamiliar)','caviv',                   4),
    (gen_random_uuid(), v_tmp, 'Comisaría de Familia',                                                'comisaria_familia',        5),
    (gen_random_uuid(), v_tmp, 'Estación de Policía',                                                 'estacion_policia',         6),
    (gen_random_uuid(), v_tmp, 'Fiscalía General de la Nación',                                       'fiscalia',                 7),
    (gen_random_uuid(), v_tmp, 'Inspección de Policía',                                                'inspeccion_policia',       8),
    (gen_random_uuid(), v_tmp, 'Instituto Nacional de Medicina Legal',                                'medicina_legal',           9),
    (gen_random_uuid(), v_tmp, 'Jueces Civiles o Promiscuos Municipales',                             'jueces_civiles',          10),
    (gen_random_uuid(), v_tmp, 'Jueces de Control de Garantías o Jueces Penales',                     'jueces_penales',          11),
    (gen_random_uuid(), v_tmp, 'Policía Judicial',                                                    'policia_judicial',        12),
    (gen_random_uuid(), v_tmp, 'Unidad de Reacción Inmediata',                                        'unidad_reaccion_inmediata',13),
    (gen_random_uuid(), v_tmp, 'No sabe / No responde',                                               'no_sabe',                 14);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q1, 'justicia', 'EQUALS');

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text', 'Otra barrera en Justicia ¿cuál?', FALSE, 7)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q5, 'otras_barreras_justicia', 'CONTAINS');

  -- ── Bloque Protección [visible: Q1 = proteccion] ────────────────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras identificadas en Protección', FALSE, 8)
  RETURNING id INTO v_q8;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_q8, 'Ausencia o demoras en la adopción de medidas de protección urgentes o definitivas en Comisaría de Familia', 'demoras_medidas_comisaria',          1),
    (gen_random_uuid(), v_q8, 'Ausencia o dilación en adopción de medidas de protección físicas',                                          'dilacion_medidas_fisicas',           2),
    (gen_random_uuid(), v_q8, 'Demora o ausencia en la activación de medidas de atención integral en Comisaría de Familia',               'demora_atencion_integral_comisaria', 3),
    (gen_random_uuid(), v_q8, 'Incumplimiento de medidas de protección sin respuesta institucional',                                       'incumplimiento_medidas_proteccion',  4),
    (gen_random_uuid(), v_q8, 'Denegación del derecho a la no confrontación',                                                              'denegacion_no_confrontacion',        5),
    (gen_random_uuid(), v_q8, 'Divulgación o filtración de información confidencial',                                                      'filtracion_info_confidencial',       6),
    (gen_random_uuid(), v_q8, 'Medidas de protección ineficaces o mal implementadas',                                                      'medidas_ineficaces',                 7),
    (gen_random_uuid(), v_q8, 'Falta de seguimiento institucional a las medidas de protección',                                            'falta_seguimiento_medidas',          8),
    (gen_random_uuid(), v_q8, 'Fallas en la valoración y actualización del riesgo',                                                        'fallas_valoracion_riesgo',           9),
    (gen_random_uuid(), v_q8, 'Medidas de protección insuficientes acordes a la situación de riesgo',                                      'medidas_insuficientes',             10),
    (gen_random_uuid(), v_q8, 'Omisión de restricción de visitas ante riesgo de feminicidio',                                              'omision_restriccion_visitas',       11),
    (gen_random_uuid(), v_q8, 'Omisión de apoyo policial en residencia o lugar de trabajo',                                                'omision_apoyo_policial',            12),
    (gen_random_uuid(), v_q8, 'Otras barreras en protección',                                                                              'otras_barreras_proteccion',         13);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_q8, v_q1, 'proteccion', 'EQUALS');

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', '¿A qué institución acudió?', FALSE, 9)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Comisarías de Familia',                                    'comisaria_familia',         1),
    (gen_random_uuid(), v_tmp, 'Fiscalía General de la Nación (funciones de protección)',  'fiscalia_proteccion',       2),
    (gen_random_uuid(), v_tmp, 'ICBF (Instituto Colombiano de Bienestar Familiar)',        'icbf',                      3),
    (gen_random_uuid(), v_tmp, 'Unidad Nacional de Protección',                            'unidad_nacional_proteccion',4),
    (gen_random_uuid(), v_tmp, 'Jueces Civiles o Promiscuos Municipales',                  'jueces_civiles',            5),
    (gen_random_uuid(), v_tmp, 'Inspección de Policía',                                     'inspeccion_policia',       6),
    (gen_random_uuid(), v_tmp, 'Defensoría de Familia',                                     'defensoria_familia',       7),
    (gen_random_uuid(), v_tmp, 'Ejército Nacional',                                         'ejercito_nacional',        8),
    (gen_random_uuid(), v_tmp, 'Armada Colombiana',                                         'armada_colombiana',        9),
    (gen_random_uuid(), v_tmp, 'No sabe / No responde',                                     'no_sabe',                 10);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q1, 'proteccion', 'EQUALS');

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text', 'Otra barrera en Protección ¿cuál?', FALSE, 10)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q8, 'otras_barreras_proteccion', 'CONTAINS');

  -- ── Otras instituciones [visible: Q1 = otras_instituciones] ─────────────────
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text', 'Nombre de la institución donde se presentó la barrera', FALSE, 11)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q1, 'otras_instituciones', 'EQUALS');

  -- ── Comunes — siempre visibles ───────────────────────────────────────────────
  -- Ubicación en cascada: sin opciones en la tabla `option`; el frontend las carga
  -- vía stateOptionsPath (mismo patrón que hacer_seguimiento.html: formState.statesColombia /
  -- formState.newBarriers[i].cities / .towns, con soporte de {_entryIndex} por entry del repeater).
  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'dropdown', 'Departamento donde se presentó la barrera', TRUE, 'statesColombia', 12)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'dropdown', 'Ciudad donde se presentó la barrera', TRUE, 'newBarriers.{_entryIndex}.cities', 13)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'dropdown', 'Municipio donde se presentó la barrera', TRUE, 'newBarriers.{_entryIndex}.towns', 14)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras institucionales y de talento humano', FALSE, 15)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Falta de aplicación del enfoque de género / Revictimización o violencia institucional',    'falta_enfoque_genero',              1),
    (gen_random_uuid(), v_tmp, 'Falta de aplicación o desactualización de protocolos de atención en VBG',                 'falta_protocolos_vbg',              2),
    (gen_random_uuid(), v_tmp, 'Capacidad institucional limitada',                                                        'capacidad_institucional_limitada',  3),
    (gen_random_uuid(), v_tmp, 'Desarticulación institucional y ausencia de gestión integral de casos',                   'desarticulacion_institucional',     4),
    (gen_random_uuid(), v_tmp, 'Insuficiencia en la orientación en derechos a las víctimas',                              'insuficiencia_orientacion_derechos',5);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras económicas y socioeconómicas', FALSE, 16)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Falta de autonomía económica',        'falta_autonomia_economica', 1),
    (gen_random_uuid(), v_tmp, 'Pobreza y desigualdad estructural',   'pobreza_desigualdad',       2);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras territoriales y geográficas', FALSE, 17)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Ausencia institucional, infraestructura y conectividad insuficiente o instalaciones distantes', 'ausencia_institucional_territorial', 1);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Barreras por ausencia de enfoque diferencial', FALSE, 18)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Barreras para personas con orientaciones sexuales e identidades de género diversas (OSIGD)', 'barrera_osigd',                      1),
    (gen_random_uuid(), v_tmp, 'Barreras para personas con diagnósticos o condiciones de salud mental',                     'barrera_salud_mental',               2),
    (gen_random_uuid(), v_tmp, 'Barreras para población con origen étnico (indígena, afrodescendiente, raizal, palenquera)','barrera_etnica',                     3),
    (gen_random_uuid(), v_tmp, 'Barreras asociadas al ciclo de vida (niñez, adolescencia, adultez mayor)',                  'barrera_ciclo_vida',                 4),
    (gen_random_uuid(), v_tmp, 'Barreras para personas con discapacidad física o cognitiva',                               'barrera_discapacidad',               5),
    (gen_random_uuid(), v_tmp, 'Barreras para personas migrantes y refugiadas',                                            'barrera_migrantes',                  6),
    (gen_random_uuid(), v_tmp, 'Barreras para personas privadas de la libertad',                                           'barrera_privadas_libertad',          7),
    (gen_random_uuid(), v_tmp, 'Barreras para personas en situación de calle',                                             'barrera_situacion_calle',            8),
    (gen_random_uuid(), v_tmp, 'Barreras para víctimas de trata de personas',                                              'barrera_trata',                      9),
    (gen_random_uuid(), v_tmp, 'Barreras para personas en actividades sexuales pagas',                                    'barrera_actividades_sexuales_pagas',10),
    (gen_random_uuid(), v_tmp, 'Barreras para lideresas y defensoras de derechos humanos',                                 'barrera_lideresas',                 11),
    (gen_random_uuid(), v_tmp, 'Barreras para víctimas del conflicto armado',                                              'barrera_conflicto_armado',          12);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'date', 'Fecha en la que se presentó la barrera (aproximada)', TRUE, 19)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text', 'Funcionario/a o dependencia donde se presentó la barrera', TRUE, 20)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text', 'Descripción de la barrera', TRUE, 21)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Gestión de la barrera', TRUE, 22)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Orientación y enrutamiento - Llamada',      'orientacion_llamada',                1),
    (gen_random_uuid(), v_tmp, 'Gestión administrativa - Llamada',         'gestion_llamada',                    2),
    (gen_random_uuid(), v_tmp, 'Activación de ruta interinstitucional',    'activacion_ruta_interinstitucional', 3),
    (gen_random_uuid(), v_tmp, 'Articulación institucional',               'articulacion_institucional',        4),
    (gen_random_uuid(), v_tmp, 'Escalamiento a organismo de control',      'escalamiento_organismo_control',    5),
    (gen_random_uuid(), v_tmp, 'Alerta por barreras',                      'alerta_barreras',                    6);

  RETURN v_rg;
END;
$func$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION pg_temp.seed_seguimiento_barreras(p_form_id uuid, p_section_id uuid)
RETURNS varchar AS $func$
DECLARE
  v_rg   varchar;   -- repeater_group id
  v_q6   varchar;   -- ¿Se realiza cierre de la barrera?
  v_tmp  varchar;
BEGIN
  INSERT INTO salvia.repeater_group (id, form_section_id, name, item_name, "order", min_repetitions)
  VALUES (gen_random_uuid(), p_section_id, 'Seguimiento a Barreras Activas', 'Barrera', 1, 0)
  RETURNING id INTO v_rg;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'info', 'Seguimiento a Barrera', FALSE, 1)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'boolean', '¿Persiste la barrera?', TRUE, 2)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'single', '¿Hubo respuesta institucional?', TRUE, 3)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Sí. La entidad respondió de manera oficial a Salvia', 'respuesta_oficial', 1),
    (gen_random_uuid(), v_tmp, 'Sí. La entidad se contactó con la víctima',          'contacto_victima',  2),
    (gen_random_uuid(), v_tmp, 'No se obtuvo respuesta',                              'sin_respuesta',     3);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'multiple', 'Gestión de la barrera', TRUE, 4)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Orientación y enrutamiento - Llamada',      'orientacion_llamada',                1),
    (gen_random_uuid(), v_tmp, 'Gestión administrativa - Llamada',         'gestion_llamada',                    2),
    (gen_random_uuid(), v_tmp, 'Activación de ruta interinstitucional',    'activacion_ruta_interinstitucional', 3),
    (gen_random_uuid(), v_tmp, 'Articulación institucional',               'articulacion_institucional',        4),
    (gen_random_uuid(), v_tmp, 'Escalamiento a organismo de control',      'escalamiento_organismo_control',    5),
    (gen_random_uuid(), v_tmp, 'Alerta por barreras',                      'alerta_barreras',                    6);

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'text',
          'Actuaciones realizadas y descripción de la gestión realizada con relación a las barreras', TRUE, 5)
  RETURNING id INTO v_tmp;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'boolean', '¿Se realiza cierre de la barrera?', TRUE, 6)
  RETURNING id INTO v_q6;

  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, "order")
  VALUES (gen_random_uuid(), p_form_id, p_section_id, v_rg, 'single', 'Motivo del cierre', TRUE, 7)
  RETURNING id INTO v_tmp;
  INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
    (gen_random_uuid(), v_tmp, 'Expresa no voluntad de accionar institucional',                                                     'no_voluntad',           1),
    (gen_random_uuid(), v_tmp, 'Barrera no gestionable desde la competencia institucional',                                         'no_gestionable',        2),
    (gen_random_uuid(), v_tmp, 'Se clasificó incorrectamente la barrera',                                                           'clasificacion_incorrecta',3),
    (gen_random_uuid(), v_tmp, 'La barrera no fue resuelta directamente, pero quedó formalmente instalada en la entidad competente','instalada_entidad',      4),
    (gen_random_uuid(), v_tmp, 'La barrera fue resuelta de manera efectiva',                                                        'resuelta',                5),
    (gen_random_uuid(), v_tmp, 'La barrera fue superada parcialmente y no requiere más gestión inmediata',                          'superada_parcialmente',   6);
  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
  VALUES (gen_random_uuid(), 'QUESTION', v_tmp, v_q6, 'true', 'EQUALS');

  RETURN v_rg;
END;
$func$ LANGUAGE plpgsql;

DO $$
DECLARE
  -- ── FORMS ─────────────────────────────────────────────────────────────────
  v_form_pc   uuid;   -- Form 1: Primer Contacto
  v_form_pa   uuid;   -- Form 2: Primera Atención
  v_form_seg  uuid;   -- Form 3: Seguimiento
  v_form_cie  uuid;   -- Form 4: Cierre

  -- ── SECCIONES — Form 1: Primer Contacto ───────────────────────────────────
  v_pc_s1     uuid;   -- Primer contacto (siempre visible)
  v_pc_s2     uuid;   -- Primera Atención (oculta; visible cuando Q12 Continuar = true)

  -- ── SECCIONES — Form 2: Primera Atención ──────────────────────────────────
  v_pa_s1     uuid;   -- Contacto Primera Atención
  v_pa_s2     uuid;   -- Primera Atención

  -- ── SECCIONES — Form 3: Seguimiento ───────────────────────────────────────
  v_seg_s1    uuid;   -- Contacto Seguimiento
  v_seg_s2    uuid;   -- Seguimiento

  -- ── SECCIONES — Form 4: Cierre ────────────────────────────────────────────
  v_cie_s1    uuid;   -- Contacto Cierre
  v_cie_s2    uuid;   -- Seguimiento (en Cierre)
  v_cie_s3    uuid;   -- Cierre

  -- ── SECCIONES DE BARRERAS (nuevas Jul 2026 — en los 4 formularios) ────────
  -- Ubicadas justo después de la sección de Contacto de cada formulario.
  v_pc_bar_seg_s   uuid;   -- PC:  Seguimiento a Barreras   (order 2)
  v_pc_bar_id_s    uuid;   -- PC:  Identificación de Barreras (order 3)
  v_pa_bar_seg_s   uuid;   -- PA:  Seguimiento a Barreras   (order 2)
  v_pa_bar_id_s    uuid;   -- PA:  Identificación de Barreras (order 3)
  v_seg_bar_seg_s  uuid;   -- SEG: Seguimiento a Barreras   (order 2)
  v_seg_bar_id_s   uuid;   -- SEG: Identificación de Barreras (order 3)
  v_cie_bar_seg_s  uuid;   -- CIE: Seguimiento a Barreras   (order 2)
  v_cie_bar_id_s   uuid;   -- CIE: Identificación de Barreras (order 3)
  v_rg_tmp         varchar; -- reutilizada solo para capturar el retorno de las funciones helper

  -- ── PREGUNTAS: Form 1 Sección 1 — Primer Contacto ─────────────────────────
  -- Orden ajustado Jul 2026: "Fecha próxima atención" movida al final (Q13)
  v_pc_q1     uuid;   -- ¿Individual o en dupla?
  v_pc_q2     uuid;   -- ¿Lugar seguro?
  v_pc_q3     uuid;   -- ¿Riesgo inminente?
  v_pc_q4     uuid;   -- Acciones ante riesgo inminente
  v_pc_q5     uuid;   -- ¿Voluntariedad?
  v_pc_q6     uuid;   -- Plan de orientación
  v_pc_q7     uuid;   -- Compromisos
  v_pc_q8     uuid;   -- Observaciones
  v_pc_q9     uuid;   -- Hay nuevos hechos de violencia
  v_pc_q10    uuid;   -- Descripción de los hechos
  v_pc_q11    uuid;   -- Fecha de los hechos
  v_pc_q12    uuid;   -- Continuar Primera Atención  ← gatillo de S4
  v_pc_q13    uuid;   -- Fecha próxima atención  ← visible si Q5=si Y Q12=false
  v_pc_q14    uuid;   -- Desde la atención anterior... ← ÚLTIMA; visible si Q12=true; gatillo de S3 (Identificación de Barreras)

  -- ── PREGUNTAS: Form 1 Sección 4 — Primera Atención (condicional en PC) ────
  -- Jul 2026: se eliminaron "¿Ajuste razonable?" y "¿Intérprete de idiomas?" (antes Q2/Q3).
  -- "Consentimiento persona de apoyo" (ahora Q3) ya no depende de esa pregunta: es siempre visible.
  v_pc_a_q1   uuid;   -- Acciones ante riesgo inminente
  v_pc_a_q2   uuid;   -- Consentimiento Informado
  v_pc_a_q3   uuid;   -- Consentimiento persona de apoyo (siempre visible)
  v_pc_a_q4   uuid;   -- Consentimiento contacto posterior para calidad
  v_pc_a_q5   uuid;   -- Ingresa por conducta suicida
  v_pc_a_q6   uuid;   -- Tipo de conducta suicida
  v_pc_a_q7   uuid;   -- Contenido de la atención
  v_pc_a_q8   uuid;   -- Plan de orientación
  v_pc_a_q9   uuid;   -- Plan de trabajo y recomendaciones
  v_pc_a_q10  uuid;   -- Compromisos
  v_pc_a_q11  uuid;   -- Fecha próxima atención
  v_pc_a_q12  uuid;   -- Observaciones

  -- ── PREGUNTAS: Form 2 Sección 1 — Contacto Primera Atención ──────────────
  v_pa_c_q1   uuid;   -- ¿Individual o en dupla?
  v_pa_c_q2   uuid;   -- ¿La llamada fue efectiva?
  v_pa_c_q3   uuid;   -- ¿Lugar seguro?
  v_pa_c_q4   uuid;   -- ¿Riesgo inminente?
  v_pa_c_q5   uuid;   -- Hay nuevos hechos de violencia
  v_pa_c_q6   uuid;   -- Descripción de los hechos
  v_pa_c_q7   uuid;   -- Fecha de los hechos
  v_pa_c_q8   uuid;   -- ¿Es atención o solo contacto?
  v_pa_c_q9   uuid;   -- Observaciones del contacto
  v_pa_c_q10  uuid;   -- Fecha nueva
  v_pa_c_q11  uuid;   -- Desde la atención anterior... ← ÚLTIMA; visible si Q8=atencion; gatillo de S3 (Identificación de Barreras)

  -- ── PREGUNTAS: Form 2 Sección 4 — Primera Atención ───────────────────────
  -- Jul 2026: se eliminaron "¿Ajuste razonable?" y "¿Intérprete de idiomas?" (antes Q2/Q3).
  -- "Consentimiento persona de apoyo" (ahora Q3) ya no depende de esa pregunta: es siempre visible.
  v_pa_a_q1   uuid;   -- Acciones ante riesgo inminente
  v_pa_a_q2   uuid;   -- Consentimiento Informado
  v_pa_a_q3   uuid;   -- Consentimiento persona de apoyo (siempre visible)
  v_pa_a_q4   uuid;   -- Consentimiento contacto posterior para calidad
  v_pa_a_q5   uuid;   -- Ingresa por conducta suicida
  v_pa_a_q6   uuid;   -- Tipo de conducta suicida
  v_pa_a_q7   uuid;   -- Contenido de la atención
  v_pa_a_q8   uuid;   -- Plan de orientación
  v_pa_a_q9   uuid;   -- Plan de trabajo y recomendaciones
  v_pa_a_q10  uuid;   -- Compromisos
  v_pa_a_q11  uuid;   -- Fecha próxima atención
  v_pa_a_q12  uuid;   -- Observaciones

  -- ── PREGUNTAS: Form 3 Sección 1 — Contacto Atención Psicosocial ──────────
  -- Jul 2026: se agrega Q1 "¿Individual o en dupla?" (antes ausente en este formulario);
  -- el resto de preguntas se recorre 1 posición.
  v_seg_c_q1  uuid;   -- ¿Individual o en dupla?
  v_seg_c_q2  uuid;   -- ¿La llamada fue efectiva?
  v_seg_c_q3  uuid;   -- ¿Lugar seguro?
  v_seg_c_q4  uuid;   -- ¿Riesgo inminente?
  v_seg_c_q5  uuid;   -- Acciones ante riesgo inminente
  v_seg_c_q6  uuid;   -- Hay nuevos hechos de violencia
  v_seg_c_q7  uuid;   -- Descripción de los hechos
  v_seg_c_q8  uuid;   -- Fecha de los hechos
  v_seg_c_q9  uuid;   -- ¿Es atención o solo contacto?
  v_seg_c_q10 uuid;   -- Observaciones del contacto
  v_seg_c_q11 uuid;   -- Fecha nueva
  v_seg_c_q12 uuid;   -- Desde la atención anterior... ← ÚLTIMA; visible si Q9=atencion; gatillo de S3 (Identificación de Barreras)

  -- ── PREGUNTAS: Form 3 Sección 2 — Seguimiento ────────────────────────────
  v_seg_s_q1  uuid;   -- Contenido de la atención
  v_seg_s_q2  uuid;   -- Plan de orientación
  v_seg_s_q3  uuid;   -- Compromisos
  v_seg_s_q4  uuid;   -- Fecha próxima atención
  v_seg_s_q5  uuid;   -- Observaciones

  -- ── PREGUNTAS: Form 4 Sección 1 — Contacto Cierre ────────────────────────
  -- Jul 2026: se agrega Q1 "¿Individual o en dupla?" (antes ausente en este formulario);
  -- el resto de preguntas se recorre 1 posición.
  v_cie_c_q1  uuid;   -- ¿Individual o en dupla?
  v_cie_c_q2  uuid;   -- ¿La llamada fue efectiva?
  v_cie_c_q3  uuid;   -- ¿Lugar seguro?
  v_cie_c_q4  uuid;   -- ¿Riesgo inminente?
  v_cie_c_q5  uuid;   -- Acciones ante riesgo inminente
  v_cie_c_q6  uuid;   -- Hay nuevos hechos de violencia
  v_cie_c_q7  uuid;   -- Descripción de los hechos
  v_cie_c_q8  uuid;   -- Fecha de los hechos
  v_cie_c_q9  uuid;   -- ¿Es atención o solo contacto?
  v_cie_c_q10 uuid;   -- Observaciones del contacto
  v_cie_c_q11 uuid;   -- Fecha nueva
  v_cie_c_q12 uuid;   -- Desde la atención anterior... ← ÚLTIMA; visible si Q9=atencion; gatillo de S3 (Identificación de Barreras)

  -- ── PREGUNTAS: Form 4 Sección 2 — Seguimiento (en Cierre) ────────────────
  v_cie_sv_q1 uuid;   -- Contenido de la atención
  v_cie_sv_q2 uuid;   -- Plan de orientación
  v_cie_sv_q3 uuid;   -- Compromisos
  v_cie_sv_q4 uuid;   -- Fecha próxima atención
  v_cie_sv_q5 uuid;   -- Observaciones
  v_cie_sv_q6 uuid;   -- Cerrar remisión  ← gatillo de S5

  -- ── PREGUNTAS: Form 4 Sección 3 — Cierre ─────────────────────────────────
  v_cie_ci_q1 uuid;   -- Motivo de cierre
  v_cie_ci_q2 uuid;   -- Contenido de la atención
  v_cie_ci_q3 uuid;   -- Plan de orientación
  v_cie_ci_q4 uuid;   -- Temas trabajados
  v_cie_ci_q5 uuid;   -- Hay nuevos hechos de violencia
  v_cie_ci_q6 uuid;   -- Descripción de los hechos
  v_cie_ci_q7 uuid;   -- Fecha de los hechos

BEGIN

-- =============================================================================
-- 1. FORMS (4 formularios)
-- =============================================================================

INSERT INTO salvia.form (id, name, description, status) VALUES (
  gen_random_uuid(),
  'Primer Contacto Psicosocial',
  'Primera llamada. ya_hizo_primer_contacto = false. S1 siempre visible; S2 "Primera Atención" visible si "Continuar = Sí".',
  'active'
) RETURNING id INTO v_form_pc;

INSERT INTO salvia.form (id, name, description, status) VALUES (
  gen_random_uuid(),
  'Primera Atención Psicosocial',
  'Escenario B: segunda llamada separada. ya_hizo_primer_contacto = true, ya_hizo_primera_atencion = false.',
  'active'
) RETURNING id INTO v_form_pa;

-- Jul 2026: renombrado de "Seguimiento Psicosocial" a "Atención Psicosocial" (y sus
-- secciones internas que contenían la palabra "Seguimiento", excepto "Seguimiento a
-- Barreras", que mantiene su nombre en los 4 formularios).
INSERT INTO salvia.form (id, name, description, status) VALUES (
  gen_random_uuid(),
  'Atención Psicosocial',
  'Sesiones de seguimiento (antes "Seguimiento Psicosocial"). ya_hizo_primera_atencion = true, session_count entre 1 y 2.',
  'active'
) RETURNING id INTO v_form_seg;

INSERT INTO salvia.form (id, name, description, status) VALUES (
  gen_random_uuid(),
  'Cierre Psicosocial',
  'Sesión con posibilidad de cierre. ya_hizo_primera_atencion = true, session_count >= 3. S3 visible si "Cerrar remisión = Sí".',
  'active'
) RETURNING id INTO v_form_cie;

RAISE NOTICE '=== FORM IDs ===';
RAISE NOTICE 'PC:  %', v_form_pc;
RAISE NOTICE 'PA:  %', v_form_pa;
RAISE NOTICE 'SEG: %', v_form_seg;
RAISE NOTICE 'CIE: %', v_form_cie;

-- =============================================================================
-- 2. SECCIONES
--
-- Cada formulario agrega, justo después de la sección de Contacto, las 2
-- secciones nuevas de Barreras. El resto de secciones existentes se recorre
-- 2 posiciones ("order" +2).
-- =============================================================================

-- Form 1 — Primer Contacto (4 secciones)
INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pc, 'Primer contacto',
        'Sección siempre visible. Incluye Q12 "Continuar Primera Atención" como gatillo de S4 y Q14 como gatillo de S3.', 1)
RETURNING id INTO v_pc_s1;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pc, 'Seguimiento a Barreras',
        'Repeater de seguimiento a barreras activas del caso. No depende de ninguna pregunta de este formulario: su visibilidad real depende de si el caso tiene barreras activas (formState externo, pendiente de definir igual que en el Formulario de Seguimiento).', 2)
RETURNING id INTO v_pc_bar_seg_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pc, 'Identificación de Barreras',
        'Repeater de registro de nuevas barreras institucionales. Visible cuando Q14 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = true.', 3)
RETURNING id INTO v_pc_bar_id_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pc, 'Primera Atención',
        'Oculta por defecto. Visible cuando Q12 "Continuar Primera Atención" = true. Mismas preguntas que PA S4.', 4)
RETURNING id INTO v_pc_s2;

-- Form 2 — Primera Atención (4 secciones)
INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pa, 'Contacto Primera Atención',
        'Preguntas de contacto para llamada separada (Escenario B). Incluye Q11 como gatillo de S3.', 1)
RETURNING id INTO v_pa_s1;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pa, 'Seguimiento a Barreras',
        'Repeater de seguimiento a barreras activas del caso. No depende de ninguna pregunta de este formulario: su visibilidad real depende de si el caso tiene barreras activas (formState externo, pendiente de definir igual que en el Formulario de Seguimiento).', 2)
RETURNING id INTO v_pa_bar_seg_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pa, 'Identificación de Barreras',
        'Repeater de registro de nuevas barreras institucionales. Visible cuando Q11 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = true.', 3)
RETURNING id INTO v_pa_bar_id_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_pa, 'Primera Atención',
        'Consentimiento, conducta suicida, contenido de la atención (Jul 2026: se eliminaron ajuste razonable e intérprete).', 4)
RETURNING id INTO v_pa_s2;

-- Form 3 — Atención Psicosocial (antes "Seguimiento", 4 secciones)
INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_seg, 'Contacto Atención Psicosocial',
        'Preguntas de contacto previas a la sesión. Incluye Q1 "¿Individual o en dupla?" (nueva, Jul 2026) y Q12 como gatillo de S3.', 1)
RETURNING id INTO v_seg_s1;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_seg, 'Seguimiento a Barreras',
        'Repeater de seguimiento a barreras activas del caso. No depende de ninguna pregunta de este formulario: su visibilidad real depende de si el caso tiene barreras activas (formState externo, pendiente de definir igual que en el Formulario de Seguimiento).', 2)
RETURNING id INTO v_seg_bar_seg_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_seg, 'Identificación de Barreras',
        'Repeater de registro de nuevas barreras institucionales. Visible cuando Q11 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = true.', 3)
RETURNING id INTO v_seg_bar_id_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_seg, 'Atención Psicosocial', 'Contenido de la sesión de atención psicosocial (antes "Seguimiento").', 4)
RETURNING id INTO v_seg_s2;

-- Form 4 — Cierre (5 secciones)
INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_cie, 'Contacto Cierre',
        'Preguntas de contacto previas al cierre. Incluye Q1 "¿Individual o en dupla?" (nueva, Jul 2026) y Q12 como gatillo de S3.', 1)
RETURNING id INTO v_cie_s1;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_cie, 'Seguimiento a Barreras',
        'Repeater de seguimiento a barreras activas del caso. No depende de ninguna pregunta de este formulario: su visibilidad real depende de si el caso tiene barreras activas (formState externo, pendiente de definir igual que en el Formulario de Seguimiento).', 2)
RETURNING id INTO v_cie_bar_seg_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_cie, 'Identificación de Barreras',
        'Repeater de registro de nuevas barreras institucionales. Visible cuando Q11 de S1 ("Desde la atención anterior se han identificado barreras institucionales") = true.', 3)
RETURNING id INTO v_cie_bar_id_s;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_cie, 'Atención Psicosocial (Cierre)',
        'Contenido de atención psicosocial (antes "Seguimiento (Cierre)"). Q6 "Cerrar remisión" controla la visibilidad de S5.', 4)
RETURNING id INTO v_cie_s2;

INSERT INTO salvia.form_section (id, form_id, name, description, "order")
VALUES (gen_random_uuid(), v_form_cie, 'Cierre',
        'Visible solo cuando CIE-SV-Q6 "Cerrar remisión" = true.', 5)
RETURNING id INTO v_cie_s3;

RAISE NOTICE '=== SECCIÓN IDs ===';
RAISE NOTICE 'PC  S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4 Primera Atención: %', v_pc_s1, v_pc_bar_seg_s, v_pc_bar_id_s, v_pc_s2;
RAISE NOTICE 'PA  S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4 Primera Atención: %', v_pa_s1, v_pa_bar_seg_s, v_pa_bar_id_s, v_pa_s2;
RAISE NOTICE 'SEG S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4 Atención Psicosocial: %', v_seg_s1, v_seg_bar_seg_s, v_seg_bar_id_s, v_seg_s2;
RAISE NOTICE 'CIE S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4 Atención Psicosocial: % | S5 Cierre: %', v_cie_s1, v_cie_bar_seg_s, v_cie_bar_id_s, v_cie_s2, v_cie_s3;

-- =============================================================================
-- 2b. CONTENIDO DE LAS SECCIONES DE BARRERAS (repeaters)
--
-- Se pobla aquí, inmediatamente después de crear las secciones, usando las
-- funciones helper definidas al inicio del script.
-- =============================================================================

v_rg_tmp := pg_temp.seed_seguimiento_barreras(v_form_pc, v_pc_bar_seg_s);
RAISE NOTICE 'PC  — repeater Seguimiento a Barreras: %', v_rg_tmp;
v_rg_tmp := pg_temp.seed_identificacion_barreras(v_form_pc, v_pc_bar_id_s);
RAISE NOTICE 'PC  — repeater Identificación de Barreras: %', v_rg_tmp;

v_rg_tmp := pg_temp.seed_seguimiento_barreras(v_form_pa, v_pa_bar_seg_s);
RAISE NOTICE 'PA  — repeater Seguimiento a Barreras: %', v_rg_tmp;
v_rg_tmp := pg_temp.seed_identificacion_barreras(v_form_pa, v_pa_bar_id_s);
RAISE NOTICE 'PA  — repeater Identificación de Barreras: %', v_rg_tmp;

v_rg_tmp := pg_temp.seed_seguimiento_barreras(v_form_seg, v_seg_bar_seg_s);
RAISE NOTICE 'SEG — repeater Seguimiento a Barreras: %', v_rg_tmp;
v_rg_tmp := pg_temp.seed_identificacion_barreras(v_form_seg, v_seg_bar_id_s);
RAISE NOTICE 'SEG — repeater Identificación de Barreras: %', v_rg_tmp;

v_rg_tmp := pg_temp.seed_seguimiento_barreras(v_form_cie, v_cie_bar_seg_s);
RAISE NOTICE 'CIE — repeater Seguimiento a Barreras: %', v_rg_tmp;
v_rg_tmp := pg_temp.seed_identificacion_barreras(v_form_cie, v_cie_bar_id_s);
RAISE NOTICE 'CIE — repeater Identificación de Barreras: %', v_rg_tmp;

-- =============================================================================
-- 3. FORM 1 — PRIMER CONTACTO — Sección 1
-- Orden ajustado Jul 2026:
--   Q1-Q8   : preguntas de contacto/voluntariedad/observaciones
--   Q9-Q11  : hechos de violencia
--   Q12     : Continuar Primera Atención  (gatillo de S4)
--   Q13     : Fecha próxima atención      (visible si Q5=si AND Q12=false)
--   Q14     : Desde la atención anterior... (ÚLTIMA — visible si Q12=true; gatillo de S3)
-- =============================================================================

-- Q1: ¿Individual o en dupla?
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'single', '¿La atención es individual o en dupla?', TRUE, 1)
RETURNING id INTO v_pc_q1;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_q1, 'Individual', 'individual', 1),
  (gen_random_uuid(), v_pc_q1, 'Dupla',      'dupla',      2);

-- Q2: ¿Lugar seguro?
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'single', '¿Se encuentra en un lugar seguro?', TRUE, 2)
RETURNING id INTO v_pc_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_q2, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_q2, 'No', 'no', 2);

-- Q3: ¿Riesgo inminente?  [visible: Q2 = no]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'single', '¿Se encuentra en riesgo inminente?', TRUE, 3)
RETURNING id INTO v_pc_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_q3, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_q3, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q3::varchar, v_pc_q2::varchar, 'no', 'EQUALS');

-- Q4: Acciones ante riesgo inminente  [visible: Q3 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'text',
        'Describa las acciones realizadas ante la situación de riesgo inminente', TRUE, 4)
RETURNING id INTO v_pc_q4;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q4::varchar, v_pc_q3::varchar, 'si', 'EQUALS');

-- Q5: ¿Voluntariedad?
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'single', '¿Hay voluntariedad para la atención?', TRUE, 5)
RETURNING id INTO v_pc_q5;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_q5, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_q5, 'No', 'no', 2);

-- Q6: Plan de orientación  [visible: Q5 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'multiple', 'Plan de orientación', FALSE, 6)
RETURNING id INTO v_pc_q6;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_q6, 'Enrutamiento',          'enrutamiento',       1),
  (gen_random_uuid(), v_pc_q6, 'Activación de ruta',    'activacion_ruta',    2),
  (gen_random_uuid(), v_pc_q6, 'Seguimiento',           'seguimiento',        3),
  (gen_random_uuid(), v_pc_q6, 'Medidas de emergencia', 'medidas_emergencia', 4),
  (gen_random_uuid(), v_pc_q6, 'Plan de estabilización','plan_estabilizacion',5);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q6::varchar, v_pc_q5::varchar, 'si', 'EQUALS');

-- Q7: Compromisos  [visible: Q5 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'text', 'Compromisos', TRUE, 7)
RETURNING id INTO v_pc_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q7::varchar, v_pc_q5::varchar, 'si', 'EQUALS');

-- Q8: Observaciones (siempre visible)
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'text', 'Observaciones', FALSE, 8)
RETURNING id INTO v_pc_q8;

-- Q9: Hay nuevos hechos de violencia
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'boolean', 'Hay nuevos hechos de violencia', TRUE, 9)
RETURNING id INTO v_pc_q9;

-- Q10: Descripción de los hechos  [visible: Q9 = true]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'text', 'Descripción de los hechos', FALSE, 10)
RETURNING id INTO v_pc_q10;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q10::varchar, v_pc_q9::varchar, 'true', 'EQUALS');

-- Q11: Fecha de los hechos  [visible: Q9 = true]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'date', 'Fecha de los hechos', FALSE, 11)
RETURNING id INTO v_pc_q11;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q11::varchar, v_pc_q9::varchar, 'true', 'EQUALS');

-- Q12: Continuar Primera Atención  ← GATILLO de Sección 2
-- Si = true → S2 (Primera Atención) se muestra en este mismo formulario.
-- Si = false → S2 permanece oculta; Q13 (Fecha próxima) aparece para agendar la próxima llamada.
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'boolean', 'Continuar Primera Atención', TRUE, 12)
RETURNING id INTO v_pc_q12;

-- Q13–Q15: ¿Agendar? + Fecha + Hora (Aug 2026)
-- Visible cuando: Q5 = si AND Q12 = false. Fecha/Hora solo si Agendar = si.
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000001-0000-4000-8000-000000000001', v_form_pc, v_pc_s1, 'single', '¿Agendar nueva sesión?', TRUE, 13);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000001-0000-4000-8000-000000000001', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000001-0000-4000-8000-000000000001', 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000001-0000-4000-8000-000000000001', v_pc_q5::varchar, 'si', 'EQUALS');
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000001-0000-4000-8000-000000000001', v_pc_q12::varchar, 'false', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('58ce2d34-24d2-4e73-bf95-26a2c608f8e6', v_form_pc, v_pc_s1, 'date', 'Fecha próxima atención', TRUE, 14)
RETURNING id INTO v_pc_q13;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', '58ce2d34-24d2-4e73-bf95-26a2c608f8e6', 'a1000001-0000-4000-8000-000000000001', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000001-0000-4000-8000-000000000001', v_form_pc, v_pc_s1, 'time', 'Hora próxima atención', TRUE, 15);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000001-0000-4000-8000-000000000001', 'a1000001-0000-4000-8000-000000000001', 'si', 'EQUALS');

-- Q16: Desde la atención anterior se han identificado barreras institucionales
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s1, 'boolean',
        'Desde la atención anterior se han identificado barreras institucionales', TRUE, 16)
RETURNING id INTO v_pc_q14;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_q14::varchar, v_pc_q12::varchar, 'true', 'EQUALS');

RAISE NOTICE 'Form PC S1 — preguntas insertadas. Gatillo Continuar (Q12): %, Fecha próxima (Q14): %, Barreras (Q16): %',
             v_pc_q12, v_pc_q13, v_pc_q14;

-- =============================================================================
-- 4. FORM 1 — PRIMER CONTACTO — Sección 4 (Primera Atención — condicional)
--
-- Mismas preguntas que Form 2 S4. La visibilidad de la sección entera
-- está controlada por Q12 (Continuar Primera Atención) = true.
-- Se inserta al final (ver sección 8 — VISIBILIDAD DE SECCIONES).
-- =============================================================================

-- PC-A-Q1: Acciones ante riesgo inminente  [visible: PC-S1-Q3 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'text',
        'Describa las acciones realizadas ante la situación de riesgo inminente', TRUE, 1)
RETURNING id INTO v_pc_a_q1;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q1::varchar, v_pc_q3::varchar, 'si', 'EQUALS');

-- PC-A-Q2: Consentimiento Informado — banner INFO (texto largo)
-- PC-A-Q2b: ¿Acepta el consentimiento? — single Sí/No (alimenta resolvePsicosocialSessionType)
-- Aug 2026: se separa info + single; Q3+ visibles solo si consentimiento = si.
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'info',
        E'CONSENTIMIENTO/DESISTIMIENTO INFORMADO PARA LA ATENCIÓN PSICOSOCIAL EN LA LÍNEA 155\n\n'
        'Por medio de este documento se establecen los términos de confidencialidad, los riesgos, las excepciones de esta y los derechos de la persona usuaria durante las conversaciones que se sostengan en el marco de la atención psicosocial telefónica de la Línea 155 SALVIA.\n\n'
        '1) Es un espacio individual gratuito que tiene como objetivo propiciar una reflexión sobre las violencias que enfrentan las mujeres y personas no binarias.\n'
        '2) Durante el proceso pueden surgir momentos o temas que despierten emociones intensas.\n'
        '3) No se hará con el objetivo de diagnosticarle alguna afectación en su salud emocional.\n'
        '4) Se entenderá que solo usted conoce completamente su propia situación.\n'
        '5) Durante las sesiones se realizarán preguntas que nos permitan identificar sus necesidades.\n'
        '6) Se realizará un máximo de cuatro sesiones de atención psicosocial telefónica de una hora c/u.\n'
        '7) La información que brinde se usará para hacer seguimiento a su proceso.\n'
        '8) De acuerdo con la Ley 1090 de 2006 y Ley 53 de 1977, la información podrá ser compartida en casos judiciales o de riesgo inminente.',
        FALSE, 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('c3296c87-d7e3-4ea1-8a0f-cbf77c561d0f'::uuid, v_form_pc, v_pc_s2, 'single',
        '¿Acepta el consentimiento informado para la atención psicosocial?',
        TRUE, 3)
RETURNING id INTO v_pc_a_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_a_q2, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_a_q2, 'No', 'no', 2);

-- PC-A-Q3: Consentimiento persona de apoyo (antes Q5) — visible si consentimiento = si
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'single',
        'Confirmación consentimiento persona de apoyo: ¿Acepta usted, como persona de apoyo, los mismos términos de confidencialidad y reserva descritos en este documento?',
        TRUE, 4)
RETURNING id INTO v_pc_a_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_a_q3, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_a_q3, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q3::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A-Q4: Consentimiento contacto posterior para calidad
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'single',
        '¿La persona da su consentimiento para ser contactada de manera posterior con fin de evaluar la calidad del servicio?',
        TRUE, 5)
RETURNING id INTO v_pc_a_q4;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_a_q4, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_a_q4, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q4::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A-Q5: Conducta suicida
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'single',
        'Ingresa por conducta suicida asociada a VBG o VpP', TRUE, 6)
RETURNING id INTO v_pc_a_q5;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_a_q5, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pc_a_q5, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q5::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A-Q6: Tipo de conducta suicida  [visible: Q5 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'single', 'Tipo de conducta suicida', FALSE, 7)
RETURNING id INTO v_pc_a_q6;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_a_q6, 'Ideación', 'ideacion', 1),
  (gen_random_uuid(), v_pc_a_q6, 'Amenaza',  'amenaza',  2),
  (gen_random_uuid(), v_pc_a_q6, 'Intento',  'intento',  3);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q6::varchar, v_pc_a_q5::varchar, 'si', 'EQUALS');

-- PC-A-Q7: Contenido de la atención
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'text', 'Contenido de la atención', TRUE, 8)
RETURNING id INTO v_pc_a_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q7::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A-Q8: Plan de orientación
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'multiple', 'Plan de orientación', FALSE, 9)
RETURNING id INTO v_pc_a_q8;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pc_a_q8, 'Enrutamiento',          'enrutamiento',       1),
  (gen_random_uuid(), v_pc_a_q8, 'Activación de ruta',    'activacion_ruta',    2),
  (gen_random_uuid(), v_pc_a_q8, 'Seguimiento',           'seguimiento',        3),
  (gen_random_uuid(), v_pc_a_q8, 'Medidas de emergencia', 'medidas_emergencia', 4),
  (gen_random_uuid(), v_pc_a_q8, 'Plan de estabilización','plan_estabilizacion',5);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q8::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A-Q9: Plan de trabajo y recomendaciones
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'text', 'Plan de trabajo y recomendaciones', TRUE, 10)
RETURNING id INTO v_pc_a_q9;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q9::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A-Q10: Compromisos
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'text', 'Compromisos', TRUE, 11)
RETURNING id INTO v_pc_a_q10;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q10::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

-- PC-A: ¿Agendar nueva sesión? + Fecha + Hora (UUIDs fijos Aug 2026)
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000002-0000-4000-8000-000000000002'::uuid, v_form_pc, v_pc_s2, 'single', '¿Agendar nueva sesión?', TRUE, 12);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000002-0000-4000-8000-000000000002'::uuid, 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000002-0000-4000-8000-000000000002'::uuid, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000002-0000-4000-8000-000000000002', v_pc_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('fd2fb664-de83-4069-a161-6348dfef48bf'::uuid, v_form_pc, v_pc_s2, 'date', 'Fecha próxima atención', TRUE, 13)
RETURNING id INTO v_pc_a_q11;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'fd2fb664-de83-4069-a161-6348dfef48bf', 'a1000002-0000-4000-8000-000000000002', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000002-0000-4000-8000-000000000002'::uuid, v_form_pc, v_pc_s2, 'time', 'Hora próxima atención', TRUE, 14);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000002-0000-4000-8000-000000000002', 'a1000002-0000-4000-8000-000000000002', 'si', 'EQUALS');

-- PC-A observaciones
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pc, v_pc_s2, 'text', 'Observaciones', FALSE, 15)
RETURNING id INTO v_pc_a_q12;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pc_a_q12::varchar, v_pc_a_q2::varchar, 'si', 'EQUALS');

RAISE NOTICE 'Form PC S4 — preguntas insertadas (consentimiento info+single, agenda condicional Aug 2026).';

-- =============================================================================
-- 5. FORM 2 — PRIMERA ATENCIÓN
-- =============================================================================

-- ─── Sección 1: Contacto Primera Atención ────────────────────────────────────

-- PA-C-Q1: ¿Individual o en dupla?
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'single', '¿La atención es individual o en dupla?', TRUE, 1)
RETURNING id INTO v_pa_c_q1;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_c_q1, 'Individual', 'individual', 1),
  (gen_random_uuid(), v_pa_c_q1, 'Dupla',      'dupla',      2);

-- PA-C-Q2: ¿La llamada fue efectiva?
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'single', '¿La llamada fue efectiva?', TRUE, 2)
RETURNING id INTO v_pa_c_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_c_q2, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_c_q2, 'No', 'no', 2);

-- PA-C-Q3: ¿Lugar seguro?  [visible: Q2 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'single', '¿Se encuentra en un lugar seguro?', TRUE, 3)
RETURNING id INTO v_pa_c_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_c_q3, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_c_q3, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_c_q3::varchar, v_pa_c_q2::varchar, 'si', 'EQUALS');

-- PA-C-Q4: ¿Riesgo inminente?  [visible: Q3 = no]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'single', '¿Se encuentra en riesgo inminente?', TRUE, 4)
RETURNING id INTO v_pa_c_q4;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_c_q4, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_c_q4, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_c_q4::varchar, v_pa_c_q3::varchar, 'no', 'EQUALS');

-- PA-C-Q5: Hay nuevos hechos de violencia  [visible: Q2 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'boolean', 'Hay nuevos hechos de violencia', TRUE, 5)
RETURNING id INTO v_pa_c_q5;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_c_q5::varchar, v_pa_c_q2::varchar, 'si', 'EQUALS');

-- PA-C-Q6: Descripción de los hechos  [visible: Q5 = true]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'text', 'Descripción de los hechos', FALSE, 6)
RETURNING id INTO v_pa_c_q6;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_c_q6::varchar, v_pa_c_q5::varchar, 'true', 'EQUALS');

-- PA-C-Q7: Fecha de los hechos  [visible: Q5 = true]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'date', 'Fecha de los hechos', FALSE, 7)
RETURNING id INTO v_pa_c_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_c_q7::varchar, v_pa_c_q5::varchar, 'true', 'EQUALS');

-- PA-C-Q8: ¿Es atención o solo contacto?
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'single', '¿Es atención o solo contacto?', TRUE, 8)
RETURNING id INTO v_pa_c_q8;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_c_q8, 'Atención',      'atencion',      1),
  (gen_random_uuid(), v_pa_c_q8, 'Solo Contacto', 'solo_contacto', 2);

-- PA-C-Q9: Observaciones del contacto (siempre visible)
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'text', 'Observaciones del contacto', FALSE, 9)
RETURNING id INTO v_pa_c_q9;

-- PA-C: ¿Agendar? + Fecha nueva + Hora [visible: Q8 = solo_contacto]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000003-0000-4000-8000-000000000003', v_form_pa, v_pa_s1, 'single', '¿Agendar nueva sesión?', TRUE, 10);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000003-0000-4000-8000-000000000003', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000003-0000-4000-8000-000000000003', 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000003-0000-4000-8000-000000000003', v_pa_c_q8::varchar, 'solo_contacto', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('64d63b79-edee-464b-be56-1104efd31a46', v_form_pa, v_pa_s1, 'date', 'Fecha nueva', FALSE, 11)
RETURNING id INTO v_pa_c_q10;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', '64d63b79-edee-464b-be56-1104efd31a46', 'a1000003-0000-4000-8000-000000000003', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000003-0000-4000-8000-000000000003', v_form_pa, v_pa_s1, 'time', 'Hora próxima atención', TRUE, 12);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000003-0000-4000-8000-000000000003', 'a1000003-0000-4000-8000-000000000003', 'si', 'EQUALS');

-- PA-C-Q13: Desde la atención anterior...
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s1, 'boolean',
        'Desde la atención anterior se han identificado barreras institucionales', TRUE, 13)
RETURNING id INTO v_pa_c_q11;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_c_q11::varchar, v_pa_c_q8::varchar, 'atencion', 'EQUALS');

-- ─── Sección 4: Primera Atención ─────────────────────────────────────────────

-- PA-A-Q1: Acciones ante riesgo inminente  [visible: PA-C-Q4 = si]
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'text',
        'Describa las acciones realizadas ante la situación de riesgo inminente', TRUE, 1)
RETURNING id INTO v_pa_a_q1;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q1::varchar, v_pa_c_q4::varchar, 'si', 'EQUALS');

-- PA-A-Q2: Consentimiento — banner INFO + single (Aug 2026)
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'info',
        E'CONSENTIMIENTO/DESISTIMIENTO INFORMADO PARA LA ATENCIÓN PSICOSOCIAL EN LA LÍNEA 155\n\n'
        'Por medio de este documento se establecen los términos de confidencialidad, los riesgos, las excepciones de esta y los derechos de la persona usuaria durante las conversaciones que se sostengan en el marco de la atención psicosocial telefónica de la Línea 155 SALVIA.\n\n'
        '1) Es un espacio individual gratuito que tiene como objetivo propiciar una reflexión sobre las violencias que enfrentan las mujeres y personas no binarias.\n'
        '2) Durante el proceso pueden surgir momentos o temas que despierten emociones intensas.\n'
        '3) No se hará con el objetivo de diagnosticarle alguna afectación en su salud emocional.\n'
        '4) Se entenderá que solo usted conoce completamente su propia situación.\n'
        '5) Durante las sesiones se realizarán preguntas que nos permitan identificar sus necesidades.\n'
        '6) Se realizará un máximo de cuatro sesiones de atención psicosocial telefónica de una hora c/u.\n'
        '7) La información que brinde se usará para hacer seguimiento a su proceso.\n'
        '8) De acuerdo con la Ley 1090 de 2006 y Ley 53 de 1977, la información podrá ser compartida en casos judiciales o de riesgo inminente.',
        FALSE, 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('7256b91e-861b-48b3-9216-2ac13f0ae889', v_form_pa, v_pa_s2, 'single',
        '¿Acepta el consentimiento informado para la atención psicosocial?', TRUE, 3)
RETURNING id INTO v_pa_a_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_a_q2, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_a_q2, 'No', 'no', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'single',
        'Confirmación consentimiento persona de apoyo: ¿Acepta usted, como persona de apoyo, los mismos términos de confidencialidad y reserva descritos en este documento?',
        TRUE, 4)
RETURNING id INTO v_pa_a_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_a_q3, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_a_q3, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q3::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'single',
        '¿La persona da su consentimiento para ser contactada de manera posterior con fin de evaluar la calidad del servicio?',
        TRUE, 5)
RETURNING id INTO v_pa_a_q4;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_a_q4, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_a_q4, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q4::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'single',
        'Ingresa por conducta suicida asociada a VBG o VpP', TRUE, 6)
RETURNING id INTO v_pa_a_q5;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_a_q5, 'Sí', 'si', 1),
  (gen_random_uuid(), v_pa_a_q5, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q5::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'single', 'Tipo de conducta suicida', FALSE, 7)
RETURNING id INTO v_pa_a_q6;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_a_q6, 'Ideación', 'ideacion', 1),
  (gen_random_uuid(), v_pa_a_q6, 'Amenaza',  'amenaza',  2),
  (gen_random_uuid(), v_pa_a_q6, 'Intento',  'intento',  3);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q6::varchar, v_pa_a_q5::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'text', 'Contenido de la atención', TRUE, 8)
RETURNING id INTO v_pa_a_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q7::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'multiple', 'Plan de orientación', FALSE, 9)
RETURNING id INTO v_pa_a_q8;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_pa_a_q8, 'Enrutamiento',          'enrutamiento',       1),
  (gen_random_uuid(), v_pa_a_q8, 'Activación de ruta',    'activacion_ruta',    2),
  (gen_random_uuid(), v_pa_a_q8, 'Seguimiento',           'seguimiento',        3),
  (gen_random_uuid(), v_pa_a_q8, 'Medidas de emergencia', 'medidas_emergencia', 4),
  (gen_random_uuid(), v_pa_a_q8, 'Plan de estabilización','plan_estabilizacion',5);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q8::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'text', 'Plan de trabajo y recomendaciones', TRUE, 10)
RETURNING id INTO v_pa_a_q9;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q9::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'text', 'Compromisos', TRUE, 11)
RETURNING id INTO v_pa_a_q10;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q10::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000004-0000-4000-8000-000000000004', v_form_pa, v_pa_s2, 'single', '¿Agendar nueva sesión?', TRUE, 12);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000004-0000-4000-8000-000000000004', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000004-0000-4000-8000-000000000004', 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000004-0000-4000-8000-000000000004', v_pa_a_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('d68c7334-74bb-47b1-a2ee-f1d3da04627b', v_form_pa, v_pa_s2, 'date', 'Fecha próxima atención', TRUE, 13)
RETURNING id INTO v_pa_a_q11;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'd68c7334-74bb-47b1-a2ee-f1d3da04627b', 'a1000004-0000-4000-8000-000000000004', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000004-0000-4000-8000-000000000004', v_form_pa, v_pa_s2, 'time', 'Hora próxima atención', TRUE, 14);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000004-0000-4000-8000-000000000004', 'a1000004-0000-4000-8000-000000000004', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_pa, v_pa_s2, 'text', 'Observaciones', FALSE, 15)
RETURNING id INTO v_pa_a_q12;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_pa_a_q12::varchar, v_pa_a_q2::varchar, 'si', 'EQUALS');

RAISE NOTICE 'Form PA S4 — consentimiento info+single + agenda Aug 2026.';

-- =============================================================================
-- 6. FORM 3 — ATENCIÓN PSICOSOCIAL (antes "SEGUIMIENTO")
-- =============================================================================

-- ─── Sección 1: Contacto Atención Psicosocial (antes "Contacto Seguimiento") ──

-- SEG-C-Q1: ¿Individual o en dupla?  ← NUEVA (Jul 2026)
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'single', '¿La atención es individual o en dupla?', TRUE, 1)
RETURNING id INTO v_seg_c_q1;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_seg_c_q1, 'Individual', 'individual', 1),
  (gen_random_uuid(), v_seg_c_q1, 'Dupla',      'dupla',      2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'single', '¿La llamada fue efectiva?', TRUE, 2)
RETURNING id INTO v_seg_c_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_seg_c_q2, 'Sí', 'si', 1),
  (gen_random_uuid(), v_seg_c_q2, 'No', 'no', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'single', '¿Se encuentra en un lugar seguro?', TRUE, 3)
RETURNING id INTO v_seg_c_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_seg_c_q3, 'Sí', 'si', 1),
  (gen_random_uuid(), v_seg_c_q3, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q3::varchar, v_seg_c_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'single', '¿Se encuentra en riesgo inminente?', TRUE, 4)
RETURNING id INTO v_seg_c_q4;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_seg_c_q4, 'Sí', 'si', 1),
  (gen_random_uuid(), v_seg_c_q4, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q4::varchar, v_seg_c_q3::varchar, 'no', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'text',
        'Describa las acciones realizadas ante la situación de riesgo inminente', TRUE, 5)
RETURNING id INTO v_seg_c_q5;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q5::varchar, v_seg_c_q4::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'boolean', 'Hay nuevos hechos de violencia', TRUE, 6)
RETURNING id INTO v_seg_c_q6;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q6::varchar, v_seg_c_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'text', 'Descripción de los hechos', FALSE, 7)
RETURNING id INTO v_seg_c_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q7::varchar, v_seg_c_q6::varchar, 'true', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'date', 'Fecha de los hechos', FALSE, 8)
RETURNING id INTO v_seg_c_q8;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q8::varchar, v_seg_c_q6::varchar, 'true', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'single', '¿Es atención o solo contacto?', TRUE, 9)
RETURNING id INTO v_seg_c_q9;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_seg_c_q9, 'Atención',      'atencion',      1),
  (gen_random_uuid(), v_seg_c_q9, 'Solo Contacto', 'solo_contacto', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'text', 'Observaciones del contacto', FALSE, 10)
RETURNING id INTO v_seg_c_q10;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000005-0000-4000-8000-000000000005', v_form_seg, v_seg_s1, 'single', '¿Agendar nueva sesión?', TRUE, 11);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000005-0000-4000-8000-000000000005', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000005-0000-4000-8000-000000000005', 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000005-0000-4000-8000-000000000005', v_seg_c_q9::varchar, 'solo_contacto', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('72ce49d2-f853-4f4a-9f1a-f795f4d514c3', v_form_seg, v_seg_s1, 'date', 'Fecha nueva', FALSE, 12)
RETURNING id INTO v_seg_c_q11;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', '72ce49d2-f853-4f4a-9f1a-f795f4d514c3', 'a1000005-0000-4000-8000-000000000005', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000005-0000-4000-8000-000000000005', v_form_seg, v_seg_s1, 'time', 'Hora próxima atención', TRUE, 13);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000005-0000-4000-8000-000000000005', 'a1000005-0000-4000-8000-000000000005', 'si', 'EQUALS');

-- SEG-C-Q14: Desde la atención anterior...
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s1, 'boolean',
        'Desde la atención anterior se han identificado barreras institucionales', TRUE, 14)
RETURNING id INTO v_seg_c_q12;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_seg_c_q12::varchar, v_seg_c_q9::varchar, 'atencion', 'EQUALS');

-- ─── Sección 4: Atención Psicosocial (antes "Seguimiento") ───────────────────

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s2, 'text', 'Contenido de la atención', TRUE, 1)
RETURNING id INTO v_seg_s_q1;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s2, 'multiple', 'Plan de orientación', FALSE, 2)
RETURNING id INTO v_seg_s_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_seg_s_q2, 'Enrutamiento',          'enrutamiento',       1),
  (gen_random_uuid(), v_seg_s_q2, 'Activación de ruta',    'activacion_ruta',    2),
  (gen_random_uuid(), v_seg_s_q2, 'Seguimiento',           'seguimiento',        3),
  (gen_random_uuid(), v_seg_s_q2, 'Medidas de emergencia', 'medidas_emergencia', 4),
  (gen_random_uuid(), v_seg_s_q2, 'Plan de estabilización','plan_estabilizacion',5);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s2, 'text', 'Compromisos', TRUE, 3)
RETURNING id INTO v_seg_s_q3;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000006-0000-4000-8000-000000000006', v_form_seg, v_seg_s2, 'single', '¿Agendar nueva sesión?', TRUE, 4);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000006-0000-4000-8000-000000000006', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000006-0000-4000-8000-000000000006', 'No', 'no', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('7040a37d-f346-4bf7-9b80-6286b9da62c5', v_form_seg, v_seg_s2, 'date', 'Fecha próxima atención', TRUE, 5)
RETURNING id INTO v_seg_s_q4;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', '7040a37d-f346-4bf7-9b80-6286b9da62c5', 'a1000006-0000-4000-8000-000000000006', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000006-0000-4000-8000-000000000006', v_form_seg, v_seg_s2, 'time', 'Hora próxima atención', TRUE, 6);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000006-0000-4000-8000-000000000006', 'a1000006-0000-4000-8000-000000000006', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_seg, v_seg_s2, 'text', 'Observaciones', FALSE, 7)
RETURNING id INTO v_seg_s_q5;

RAISE NOTICE 'Form SEG — preguntas insertadas.';

-- =============================================================================
-- 7. FORM 4 — CIERRE
-- =============================================================================

-- ─── Sección 1: Contacto Cierre ──────────────────────────────────────────────

-- CIE-C-Q1: ¿Individual o en dupla?  ← NUEVA (Jul 2026)
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'single', '¿La atención es individual o en dupla?', TRUE, 1)
RETURNING id INTO v_cie_c_q1;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_c_q1, 'Individual', 'individual', 1),
  (gen_random_uuid(), v_cie_c_q1, 'Dupla',      'dupla',      2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'single', '¿La llamada fue efectiva?', TRUE, 2)
RETURNING id INTO v_cie_c_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_c_q2, 'Sí', 'si', 1),
  (gen_random_uuid(), v_cie_c_q2, 'No', 'no', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'single', '¿Se encuentra en un lugar seguro?', TRUE, 3)
RETURNING id INTO v_cie_c_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_c_q3, 'Sí', 'si', 1),
  (gen_random_uuid(), v_cie_c_q3, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q3::varchar, v_cie_c_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'single', '¿Se encuentra en riesgo inminente?', TRUE, 4)
RETURNING id INTO v_cie_c_q4;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_c_q4, 'Sí', 'si', 1),
  (gen_random_uuid(), v_cie_c_q4, 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q4::varchar, v_cie_c_q3::varchar, 'no', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'text',
        'Describa las acciones realizadas ante la situación de riesgo inminente', TRUE, 5)
RETURNING id INTO v_cie_c_q5;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q5::varchar, v_cie_c_q4::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'boolean', 'Hay nuevos hechos de violencia', TRUE, 6)
RETURNING id INTO v_cie_c_q6;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q6::varchar, v_cie_c_q2::varchar, 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'text', 'Descripción de los hechos', FALSE, 7)
RETURNING id INTO v_cie_c_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q7::varchar, v_cie_c_q6::varchar, 'true', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'date', 'Fecha de los hechos', FALSE, 8)
RETURNING id INTO v_cie_c_q8;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q8::varchar, v_cie_c_q6::varchar, 'true', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'single', '¿Es atención o solo contacto?', TRUE, 9)
RETURNING id INTO v_cie_c_q9;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_c_q9, 'Atención',      'atencion',      1),
  (gen_random_uuid(), v_cie_c_q9, 'Solo Contacto', 'solo_contacto', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'text', 'Observaciones del contacto', FALSE, 10)
RETURNING id INTO v_cie_c_q10;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000007-0000-4000-8000-000000000007', v_form_cie, v_cie_s1, 'single', '¿Agendar nueva sesión?', TRUE, 11);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000007-0000-4000-8000-000000000007', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000007-0000-4000-8000-000000000007', 'No', 'no', 2);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a1000007-0000-4000-8000-000000000007', v_cie_c_q9::varchar, 'solo_contacto', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('1b4d09f0-e5ca-4b4e-9d74-428478a113c6', v_form_cie, v_cie_s1, 'date', 'Fecha nueva', FALSE, 12)
RETURNING id INTO v_cie_c_q11;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', '1b4d09f0-e5ca-4b4e-9d74-428478a113c6', 'a1000007-0000-4000-8000-000000000007', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000007-0000-4000-8000-000000000007', v_form_cie, v_cie_s1, 'time', 'Hora próxima atención', TRUE, 13);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000007-0000-4000-8000-000000000007', 'a1000007-0000-4000-8000-000000000007', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s1, 'boolean',
        'Desde la atención anterior se han identificado barreras institucionales', TRUE, 14)
RETURNING id INTO v_cie_c_q12;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_c_q12::varchar, v_cie_c_q9::varchar, 'atencion', 'EQUALS');

-- ─── Sección 4: Atención Psicosocial (en Cierre) — antes "Seguimiento (en Cierre)" ──

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s2, 'text', 'Contenido de la atención', TRUE, 1)
RETURNING id INTO v_cie_sv_q1;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s2, 'multiple', 'Plan de orientación', FALSE, 2)
RETURNING id INTO v_cie_sv_q2;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_sv_q2, 'Enrutamiento',          'enrutamiento',       1),
  (gen_random_uuid(), v_cie_sv_q2, 'Activación de ruta',    'activacion_ruta',    2),
  (gen_random_uuid(), v_cie_sv_q2, 'Seguimiento',           'seguimiento',        3),
  (gen_random_uuid(), v_cie_sv_q2, 'Medidas de emergencia', 'medidas_emergencia', 4),
  (gen_random_uuid(), v_cie_sv_q2, 'Plan de estabilización','plan_estabilizacion',5);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s2, 'text', 'Compromisos', TRUE, 3)
RETURNING id INTO v_cie_sv_q3;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a1000008-0000-4000-8000-000000000008', v_form_cie, v_cie_s2, 'single', '¿Agendar nueva sesión?', TRUE, 4);
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), 'a1000008-0000-4000-8000-000000000008', 'Sí', 'si', 1),
  (gen_random_uuid(), 'a1000008-0000-4000-8000-000000000008', 'No', 'no', 2);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('3ce6ff5a-f139-4417-9255-155207e9a970', v_form_cie, v_cie_s2, 'date', 'Fecha próxima atención', TRUE, 5)
RETURNING id INTO v_cie_sv_q4;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', '3ce6ff5a-f139-4417-9255-155207e9a970', 'a1000008-0000-4000-8000-000000000008', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES ('a2000008-0000-4000-8000-000000000008', v_form_cie, v_cie_s2, 'time', 'Hora próxima atención', TRUE, 6);
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', 'a2000008-0000-4000-8000-000000000008', 'a1000008-0000-4000-8000-000000000008', 'si', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s2, 'text', 'Observaciones', FALSE, 7)
RETURNING id INTO v_cie_sv_q5;

-- CIE-SV-Q8: Cerrar remisión  ← GATILLO de Sección 5
INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s2, 'boolean', 'Cerrar remisión', TRUE, 8)
RETURNING id INTO v_cie_sv_q6;

RAISE NOTICE 'Form CIE S4 — Q8 "Cerrar remisión" (gatillo S5): %', v_cie_sv_q6;

-- ─── Sección 5: Cierre ────────────────────────────────────────────────────────

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'single', 'Motivo de cierre', TRUE, 1)
RETURNING id INTO v_cie_ci_q1;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_ci_q1, 'Cumplimiento de objetivos',       'cumplimiento_objetivos', 1),
  (gen_random_uuid(), v_cie_ci_q1, 'Cumplimiento esquema',            'cumplimiento_esquema',   2),
  (gen_random_uuid(), v_cie_ci_q1, 'No consentimiento',               'no_consentimiento',      3),
  (gen_random_uuid(), v_cie_ci_q1, 'Imposibilidad del contacto (3x3)','imposibilidad_contacto', 4),
  (gen_random_uuid(), v_cie_ci_q1, 'Desistimiento del proceso',       'desistimiento',          5);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'text', 'Contenido de la atención', TRUE, 2)
RETURNING id INTO v_cie_ci_q2;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'multiple', 'Plan de orientación', FALSE, 3)
RETURNING id INTO v_cie_ci_q3;
INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES
  (gen_random_uuid(), v_cie_ci_q3, 'Enrutamiento',          'enrutamiento',       1),
  (gen_random_uuid(), v_cie_ci_q3, 'Activación de ruta',    'activacion_ruta',    2),
  (gen_random_uuid(), v_cie_ci_q3, 'Seguimiento',           'seguimiento',        3),
  (gen_random_uuid(), v_cie_ci_q3, 'Medidas de emergencia', 'medidas_emergencia', 4),
  (gen_random_uuid(), v_cie_ci_q3, 'Plan de estabilización','plan_estabilizacion',5);

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'text', 'Temas trabajados durante la atención', TRUE, 4)
RETURNING id INTO v_cie_ci_q4;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'boolean', 'Hay nuevos hechos de violencia', TRUE, 5)
RETURNING id INTO v_cie_ci_q5;

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'text', 'Descripción de los hechos', FALSE, 6)
RETURNING id INTO v_cie_ci_q6;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_ci_q6::varchar, v_cie_ci_q5::varchar, 'true', 'EQUALS');

INSERT INTO salvia.question (id, form_id, form_section_id, question_type, description, required, "order")
VALUES (gen_random_uuid(), v_form_cie, v_cie_s3, 'date', 'Fecha de los hechos', FALSE, 7)
RETURNING id INTO v_cie_ci_q7;
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'QUESTION', v_cie_ci_q7::varchar, v_cie_ci_q5::varchar, 'true', 'EQUALS');

RAISE NOTICE 'Form CIE — preguntas insertadas.';

-- =============================================================================
-- 8. VISIBILIDAD DE SECCIONES
--
-- Insertar al final para que todas las preguntas gatillo ya existan.
--
-- PC  S3 (Identificación de Barreras) : trigger = PC-S1-Q14  "Desde la atención anterior..." = true
-- PC  S4 (Primera Atención en PC)     : trigger = PC-S1-Q12  "Continuar Primera Atención" = true
-- PA  S3 (Identificación de Barreras) : trigger = PA-C-Q11   "Desde la atención anterior..." = true
-- PA  S4 (Primera Atención)           : trigger = PA-C-Q8    "¿Es atención o solo contacto?" = atencion
-- SEG S3 (Identificación de Barreras) : trigger = SEG-C-Q12  "Desde la atención anterior..." = true
-- SEG S4 (Atención Psicosocial)       : trigger = SEG-C-Q9   "¿Es atención o solo contacto?" = atencion
-- CIE S3 (Identificación de Barreras) : trigger = CIE-C-Q12  "Desde la atención anterior..." = true
-- CIE S4 (Atención Psicosocial-Cierre): trigger = CIE-C-Q9   "¿Es atención o solo contacto?" = atencion
-- CIE S5 (Cierre)                     : trigger = CIE-SV-Q6  "Cerrar remisión" = true
--
-- Nota: "Seguimiento a Barreras" (S2 en los 4 formularios) NO recibe visibility_condition
-- aquí — permanece habilitada sin importar la respuesta de ninguna pregunta del formulario.
-- Su visibilidad real (según barreras activas del caso) queda pendiente de resolver vía
-- formState externo, igual que en el Formulario de Seguimiento (ver
-- form-seguimiento-barreras-repeater.md, "Pendientes / Decisiones abiertas").
-- =============================================================================

-- Form 1: PC — S3 "Identificación de Barreras" visible cuando Q14 = true
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_pc_bar_id_s::varchar, v_pc_q14::varchar, 'true', 'EQUALS');

-- Form 1: PC — S4 "Primera Atención" visible cuando Q12 (Continuar) = true
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_pc_s2::varchar, v_pc_q12::varchar, 'true', 'EQUALS');

-- Form 2: PA — S3 "Identificación de Barreras" visible cuando PA-C-Q11 = true
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_pa_bar_id_s::varchar, v_pa_c_q11::varchar, 'true', 'EQUALS');

-- Form 2: PA — S4 "Primera Atención" visible cuando PA-C-Q8 = atencion
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_pa_s2::varchar, v_pa_c_q8::varchar, 'atencion', 'EQUALS');

-- Form 3: SEG — S3 "Identificación de Barreras" visible cuando SEG-C-Q12 = true
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_seg_bar_id_s::varchar, v_seg_c_q12::varchar, 'true', 'EQUALS');

-- Form 3: SEG — S4 "Atención Psicosocial" visible cuando SEG-C-Q9 = atencion
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_seg_s2::varchar, v_seg_c_q9::varchar, 'atencion', 'EQUALS');

-- Form 4: CIE — S3 "Identificación de Barreras" visible cuando CIE-C-Q12 = true
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_cie_bar_id_s::varchar, v_cie_c_q12::varchar, 'true', 'EQUALS');

-- Form 4: CIE — S4 "Atención Psicosocial (Cierre)" visible cuando CIE-C-Q9 = atencion
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_cie_s2::varchar, v_cie_c_q9::varchar, 'atencion', 'EQUALS');

-- Form 4: CIE — S5 "Cierre" visible cuando CIE-SV-Q6 "Cerrar remisión" = true
-- (anteriormente controlada por CIE-C-Q9, ahora por la nueva pregunta Q6 de S4)
INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_value, operator)
VALUES (gen_random_uuid(), 'SECTION', v_cie_s3::varchar, v_cie_sv_q6::varchar, 'true', 'EQUALS');

RAISE NOTICE 'Visibility conditions de secciones insertadas.';
RAISE NOTICE '  PC  S3 (Barreras-Id) trigger: Q14 Barreras anteriores → %', v_pc_q14;
RAISE NOTICE '  PC  S4 trigger: Q12 Continuar=true  → %', v_pc_q12;
RAISE NOTICE '  PA  S3 (Barreras-Id) trigger: Q11 Barreras anteriores → %', v_pa_c_q11;
RAISE NOTICE '  PA  S4 trigger: Q8  Es atención     → %', v_pa_c_q8;
RAISE NOTICE '  SEG S3 (Barreras-Id) trigger: Q12 Barreras anteriores → %', v_seg_c_q12;
RAISE NOTICE '  SEG S4 trigger: Q9  Es atención     → %', v_seg_c_q9;
RAISE NOTICE '  CIE S3 (Barreras-Id) trigger: Q12 Barreras anteriores → %', v_cie_c_q12;
RAISE NOTICE '  CIE S4 trigger: Q9  Es atención     → %', v_cie_c_q9;
RAISE NOTICE '  CIE S5 trigger: Q6  Cerrar remisión → %', v_cie_sv_q6;

-- =============================================================================
-- RESUMEN FINAL
-- =============================================================================
RAISE NOTICE '';
RAISE NOTICE '=== SEED PSICOSOCIAL COMPLETADO ===';
RAISE NOTICE 'Formularios: PC=% | PA=% | SEG=% | CIE=%', v_form_pc, v_form_pa, v_form_seg, v_form_cie;
RAISE NOTICE 'Secciones:';
RAISE NOTICE '  PC  S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4 (condicional): %', v_pc_s1, v_pc_bar_seg_s, v_pc_bar_id_s, v_pc_s2;
RAISE NOTICE '  PA  S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4: %', v_pa_s1, v_pa_bar_seg_s, v_pa_bar_id_s, v_pa_s2;
RAISE NOTICE '  SEG S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4: %', v_seg_s1, v_seg_bar_seg_s, v_seg_bar_id_s, v_seg_s2;
RAISE NOTICE '  CIE S1: % | S2 Barreras-Seg: % | S3 Barreras-Id: % | S4: % | S5: %', v_cie_s1, v_cie_bar_seg_s, v_cie_bar_id_s, v_cie_s2, v_cie_s3;
RAISE NOTICE 'Preguntas gatillo:';
RAISE NOTICE '  PC  Q12 Continuar Primera Atención:            %', v_pc_q12;
RAISE NOTICE '  PC  Q13 Fecha próxima (S1):                    %', v_pc_q13;
RAISE NOTICE '  PC  Q14 Barreras anteriores (última S1):       %', v_pc_q14;
RAISE NOTICE '  PA  Q11 Barreras anteriores (última S1):       %', v_pa_c_q11;
RAISE NOTICE '  SEG Q12 Barreras anteriores (última S1):       %', v_seg_c_q12;
RAISE NOTICE '  CIE Q12 Barreras anteriores (última S1):       %', v_cie_c_q12;
RAISE NOTICE '  CIE SV-Q6 Cerrar remisión:                     %', v_cie_sv_q6;

END $$;
