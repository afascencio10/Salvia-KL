/**
 * mock_psicosocial.js
 * Simula las respuestas del endpoint GET /api/v1/psicosocial-support/:id/load
 * para cada caso de prueba de la pantalla de Registro de Sesión Psicosocial.
 *
 * Estructura del response real:
 * {
 *   formId:           string  (UUID del formulario seleccionado)
 *   submissionId:     string  (UUID del form_submission)
 *   victimInfo:       object
 *   psicosocialState: object
 *   formState:        object  (actualmente vacío — reservado para uso futuro)
 * }
 *
 * Uso: window.PSICO_MOCK.getCaso('caso1')
 */

(function (global) {

  // ─── IDs de formulario (pendientes de seed — valores placeholder) ──────────
  // Reemplazar con los UUIDs reales después de ejecutar seed_psicosocial.sql
  const FORM_IDS = {
    PRIMER_CONTACTO:   'FORM-PC-XXXX-SEED-PENDING',
    PRIMERA_ATENCION:  'FORM-PA-XXXX-SEED-PENDING',
    SEGUIMIENTO:       'FORM-SEG-XXXX-SEED-PENDING',
    CIERRE:            'FORM-CIE-XXXX-SEED-PENDING',
  };

  const SUBMISSION_ID = 'MOCK-SUBMISSION-00000001';

  // ─── Víctima de prueba (compartida por todos los casos) ───────────────────
  const VICTIM_INFO = {
    Names:             'Laura Valentina',
    LastNames:         'Moreno Ríos',
    Phone:             '310 456 7890',
    GenderIdentity:    'Mujer',
    TownName:          'Bogotá D.C.',
    RiskLevel:         3,
    CaseICode:         'CASO-TEST-001',
  };

  // ─── Estructura de secciones por formulario (para FormPreview) ────────────

  // Preguntas compartidas de la sección de Primera Atención
  // (usadas en PC-S2 y PA-S2 por igual)
  const PREGUNTAS_PRIMERA_ATENCION_CONTENIDO = [
    { label: 'Describa las acciones ante riesgo inminente',          type: 'text',     conditional: 'Si riesgo inminente (S1) = Sí' },
    { label: '¿Requiere ajuste razonable?',                          type: 'single',   options: ['Sí', 'No'] },
    { label: '¿Requiere intérprete de idiomas?',                     type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q2 = Sí' },
    { label: 'Consentimiento Informado para la Atención Psicosocial',type: 'single',   options: ['Sí', 'No'] },
    { label: 'Confirmación consentimiento persona de apoyo',         type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q3 = Sí' },
    { label: '¿Consentimiento contacto posterior para calidad?',     type: 'single',   options: ['Sí', 'No'] },
    { label: 'Ingresa por conducta suicida asociada a VBG o VpP',   type: 'single',   options: ['Sí', 'No'] },
    { label: 'Tipo de conducta suicida',                             type: 'single',   options: ['Ideación', 'Amenaza', 'Intento'], conditional: 'Si Q7 = Sí' },
    { label: 'Contenido de la atención',                             type: 'text' },
    { label: 'Plan de orientación',                                   type: 'multiple', options: ['Enrutamiento', 'Activación de ruta', 'Seguimiento', 'Medidas de emergencia', 'Plan de estabilización'] },
    { label: 'Plan de trabajo y recomendaciones',                    type: 'text' },
    { label: 'Compromisos',                                           type: 'text' },
    { label: 'Fecha próxima atención',                               type: 'date' },
    { label: 'Observaciones',                                         type: 'text' },
  ];

  const FORM_SECTIONS = {

    // ── Form 1: Primer Contacto ────────────────────────────────────────────
    // S2 (Primera Atención) es condicional: visible cuando Q13 (S1) = Sí
    PRIMER_CONTACTO: [
      {
        name: 'Primer contacto',
        order: 1,
        questions: [
          { label: '¿La atención es individual o en dupla?',               type: 'single',   options: ['Individual', 'Dupla'] },
          { label: '¿Se encuentra en un lugar seguro?',                    type: 'single',   options: ['Sí', 'No'] },
          { label: '¿Se encuentra en riesgo inminente?',                   type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q2 = No' },
          { label: 'Describa las acciones ante riesgo inminente',          type: 'text',     conditional: 'Si Q3 = Sí' },
          { label: '¿Hay voluntariedad para la atención?',                 type: 'single',   options: ['Sí', 'No'] },
          { label: 'Plan de orientación',                                   type: 'multiple', options: ['Enrutamiento', 'Activación de ruta', 'Seguimiento', 'Medidas de emergencia', 'Plan de estabilización'], conditional: 'Si Q5 = Sí' },
          { label: 'Compromisos',                                           type: 'text',     conditional: 'Si Q5 = Sí' },
          { label: 'Observaciones',                                         type: 'text' },
          { label: 'Hay nuevos hechos de violencia',                       type: 'boolean' },
          { label: 'Descripción de los hechos',                            type: 'text',     conditional: 'Si Q9 = Sí' },
          { label: 'Fecha de los hechos',                                  type: 'date',     conditional: 'Si Q9 = Sí' },
          { label: 'Continuar Primera Atención',                           type: 'boolean',  highlight: true, note: 'GATILLO S2: Si = Sí → Sección "Primera Atención" aparece en este mismo formulario. "Fecha próxima atención" se mueve a S2.' },
          { label: 'Fecha próxima atención',                               type: 'date',     conditional: 'Si Q5 = Sí Y Q12 = No', note: 'ÚLTIMA PREGUNTA — Se oculta cuando "Continuar = Sí" (aparece al final de S2). Visible solo cuando Q12 = No.' },
        ],
      },
      {
        name: 'Primera Atención',
        order: 2,
        conditionNote: 'Oculta por defecto. Visible cuando Q13 (S1 — Continuar Primera Atención) = Sí. El profesional completa todo en un único formulario.',
        questions: PREGUNTAS_PRIMERA_ATENCION_CONTENIDO,
      },
    ],

    // ── Form 2: Primera Atención ───────────────────────────────────────────
    // Solo para Escenario B: segunda llamada separada (Continuar = No en PC)
    PRIMERA_ATENCION: [
      {
        name: 'Contacto Primera Atención',
        order: 1,
        questions: [
          { label: '¿La atención es individual o en dupla?',               type: 'single',   options: ['Individual', 'Dupla'] },
          { label: '¿La llamada fue efectiva?',                            type: 'single',   options: ['Sí', 'No'] },
          { label: '¿Se encuentra en un lugar seguro?',                    type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q2 = Sí' },
          { label: '¿Se encuentra en riesgo inminente?',                   type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q3 = No' },
          { label: 'Hay nuevos hechos de violencia',                       type: 'boolean',  conditional: 'Si Q2 = Sí' },
          { label: 'Descripción de los hechos',                            type: 'text',     conditional: 'Si Q5 = Sí' },
          { label: 'Fecha de los hechos',                                  type: 'date',     conditional: 'Si Q5 = Sí' },
          { label: '¿Es atención o solo contacto?',                       type: 'multiple', options: ['Atención', 'Solo Contacto'], highlight: true, note: 'GATILLO S2: "Solo Contacto" oculta la sección siguiente' },
          { label: 'Observaciones del contacto',                           type: 'text' },
          { label: 'Fecha nueva',                                          type: 'date',     conditional: 'Si Q8 = Solo Contacto' },
        ],
      },
      {
        name: 'Primera Atención',
        order: 2,
        conditionNote: 'Visible solo cuando Q8 (sección anterior) = Atención',
        questions: PREGUNTAS_PRIMERA_ATENCION_CONTENIDO,
      },
    ],

    // ── Form 3: Seguimiento ────────────────────────────────────────────────
    SEGUIMIENTO: [
      {
        name: 'Contacto Seguimiento',
        order: 1,
        questions: [
          { label: '¿La llamada fue efectiva?',                            type: 'single',   options: ['Sí', 'No'] },
          { label: '¿Se encuentra en un lugar seguro?',                    type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q1 = Sí' },
          { label: '¿Se encuentra en riesgo inminente?',                   type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q2 = No' },
          { label: 'Describa las acciones ante riesgo inminente',          type: 'text',     conditional: 'Si Q3 = Sí' },
          { label: 'Hay nuevos hechos de violencia',                       type: 'boolean',  conditional: 'Si Q1 = Sí' },
          { label: 'Descripción de los hechos',                            type: 'text',     conditional: 'Si Q5 = Sí' },
          { label: 'Fecha de los hechos',                                  type: 'date',     conditional: 'Si Q5 = Sí' },
          { label: '¿Es atención o solo contacto?',                       type: 'multiple', options: ['Atención', 'Solo Contacto'], highlight: true, note: 'GATILLO S2: "Solo Contacto" oculta la sección siguiente' },
          { label: 'Observaciones del contacto',                           type: 'text' },
          { label: 'Fecha nueva',                                          type: 'date',     conditional: 'Si Q8 = Solo Contacto' },
        ],
      },
      {
        name: 'Seguimiento',
        order: 2,
        conditionNote: 'Visible solo cuando Q8 (sección anterior) = Atención',
        questions: [
          { label: 'Contenido de la atención',                             type: 'text' },
          { label: 'Plan de orientación',                                   type: 'multiple', options: ['Enrutamiento', 'Activación de ruta', 'Seguimiento', 'Medidas de emergencia', 'Plan de estabilización'] },
          { label: 'Compromisos',                                           type: 'text' },
          { label: 'Fecha próxima atención',                               type: 'date' },
          { label: 'Observaciones',                                         type: 'text' },
        ],
      },
    ],

    // ── Form 4: Cierre ─────────────────────────────────────────────────────
    // S3 es condicional: visible cuando S2-Q6 "Cerrar remisión" = Sí
    CIERRE: [
      {
        name: 'Contacto Cierre',
        order: 1,
        questions: [
          { label: '¿La llamada fue efectiva?',                            type: 'single',   options: ['Sí', 'No'] },
          { label: '¿Se encuentra en un lugar seguro?',                    type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q1 = Sí' },
          { label: '¿Se encuentra en riesgo inminente?',                   type: 'single',   options: ['Sí', 'No'],   conditional: 'Si Q2 = No' },
          { label: 'Describa las acciones ante riesgo inminente',          type: 'text',     conditional: 'Si Q3 = Sí' },
          { label: 'Hay nuevos hechos de violencia',                       type: 'boolean',  conditional: 'Si Q1 = Sí' },
          { label: 'Descripción de los hechos',                            type: 'text',     conditional: 'Si Q5 = Sí' },
          { label: 'Fecha de los hechos',                                  type: 'date',     conditional: 'Si Q5 = Sí' },
          { label: '¿Es atención o solo contacto?',                       type: 'multiple', options: ['Atención', 'Solo Contacto'], highlight: true, note: 'GATILLO S2: "Solo Contacto" oculta las secciones 2 y 3' },
          { label: 'Observaciones del contacto',                           type: 'text' },
          { label: 'Fecha nueva',                                          type: 'date',     conditional: 'Si Q8 = Solo Contacto' },
        ],
      },
      {
        name: 'Seguimiento (Cierre)',
        order: 2,
        conditionNote: 'Visible solo cuando Q8 (sección anterior) = Atención',
        questions: [
          { label: 'Contenido de la atención',                             type: 'text' },
          { label: 'Plan de orientación',                                   type: 'multiple', options: ['Enrutamiento', 'Activación de ruta', 'Seguimiento', 'Medidas de emergencia', 'Plan de estabilización'] },
          { label: 'Compromisos',                                           type: 'text' },
          { label: 'Fecha próxima atención',                               type: 'date' },
          { label: 'Observaciones',                                         type: 'text' },
          { label: 'Cerrar remisión',                                       type: 'boolean',  highlight: true, note: 'GATILLO S3: Si = Sí → Sección "Cierre" aparece. Si = No → el formulario actúa como seguimiento y status NO cambia a cerrado.' },
        ],
      },
      {
        name: 'Cierre',
        order: 3,
        conditionNote: 'Visible solo cuando: (1) Q8 sección 1 = Atención  Y  (2) Q6 sección 2 "Cerrar remisión" = Sí',
        questions: [
          { label: 'Motivo de cierre',                                      type: 'single',   options: ['Cumplimiento de objetivos', 'Cumplimiento esquema', 'No consentimiento', 'Imposibilidad del contacto (3x3)', 'Desistimiento del proceso'] },
          { label: 'Contenido de la atención',                             type: 'text' },
          { label: 'Plan de orientación',                                   type: 'multiple', options: ['Enrutamiento', 'Activación de ruta', 'Seguimiento', 'Medidas de emergencia', 'Plan de estabilización'] },
          { label: 'Temas trabajados durante la atención',                 type: 'text' },
          { label: 'Hay nuevos hechos de violencia',                       type: 'boolean' },
          { label: 'Descripción de los hechos',                            type: 'text',     conditional: 'Si Q5 = Sí' },
          { label: 'Fecha de los hechos',                                  type: 'date',     conditional: 'Si Q5 = Sí' },
        ],
      },
    ],
  };

  // ─── Casos de prueba ──────────────────────────────────────────────────────
  const CASOS = {

    /**
     * Caso 1: Primera llamada — solo Primer Contacto (Continuar = No)
     * El profesional no continúa con Primera Atención en esta sesión.
     * S2 permanece oculta; se agenda una próxima fecha.
     */
    caso1: {
      _meta: {
        title: 'Primer Contacto (sin continuar)',
        subtitle: 'Primera llamada. "Continuar Primera Atención = No". S2 oculta. Se agenda próxima llamada.',
        badge: { color: '#f59e0b', label: 'Caso 1' },
      },
      formId:       FORM_IDS.PRIMER_CONTACTO,
      submissionId: SUBMISSION_ID,
      formType:     'PRIMER_CONTACTO',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  false,
        yaHizoPrimeraAtencion: false,
        sessionCount:          0,
        status:                'abierto',
      },
      formState: {},
      canEdit: true,
      sections: FORM_SECTIONS.PRIMER_CONTACTO,
    },

    /**
     * Caso 2: Primera llamada — Primer Contacto + Primera Atención en la misma sesión
     * El profesional responde "Continuar Primera Atención = Sí".
     * S2 (Primera Atención) se despliega dentro del mismo formulario.
     * "Fecha próxima atención" desaparece de S1 y queda al final de S2.
     */
    caso2: {
      _meta: {
        title: 'Primer Contacto + Primera Atención (misma sesión)',
        subtitle: '"Continuar Primera Atención = Sí". S2 visible en el Form PC. Todo en un único formulario.',
        badge: { color: '#8b5cf6', label: 'Caso 2' },
      },
      formId:       FORM_IDS.PRIMER_CONTACTO,
      submissionId: SUBMISSION_ID,
      formType:     'PRIMER_CONTACTO',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  false,
        yaHizoPrimeraAtencion: false,
        sessionCount:          0,
        status:                'abierto',
      },
      formState: {},
      canEdit: true,
      sections: FORM_SECTIONS.PRIMER_CONTACTO,
      // Nota para el mock preview:
      // En este caso S2 debería renderizarse como visible desde el inicio
      // para simular que el profesional ya marcó "Continuar = Sí".
      _previewState: { continuarPrimeraAtencion: true },
    },

    /**
     * Caso 3: Segunda llamada — Primera Atención independiente (Escenario B)
     * ya_hizo_primer_contacto = true, ya_hizo_primera_atencion = false.
     * El profesional respondió "Continuar = No" en la sesión anterior.
     * Se carga el Form de Primera Atención con su sección de Contacto visible.
     */
    caso3: {
      _meta: {
        title: 'Primera Atención (llamada separada)',
        subtitle: 'Segunda llamada independiente. ya_hizo_primer_contacto = true. Form PA con sección de Contacto visible.',
        badge: { color: '#3b82f6', label: 'Caso 3' },
      },
      formId:       FORM_IDS.PRIMERA_ATENCION,
      submissionId: SUBMISSION_ID,
      formType:     'PRIMERA_ATENCION',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  true,
        yaHizoPrimeraAtencion: false,
        sessionCount:          0,
        status:                'en_gestion',
      },
      formState: {},
      canEdit: true,
      sections: FORM_SECTIONS.PRIMERA_ATENCION,
    },

    /**
     * Caso 4: Seguimiento regular (1 sesión completada)
     */
    caso4: {
      _meta: {
        title: 'Seguimiento (sesión 2)',
        subtitle: 'ya_hizo_primera_atencion = true, session_count = 1. Form Seguimiento.',
        badge: { color: '#10b981', label: 'Caso 4' },
      },
      formId:       FORM_IDS.SEGUIMIENTO,
      submissionId: SUBMISSION_ID,
      formType:     'SEGUIMIENTO',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  true,
        yaHizoPrimeraAtencion: true,
        sessionCount:          1,
        status:                'en_gestion',
      },
      formState: {},
      canEdit: true,
      sections: FORM_SECTIONS.SEGUIMIENTO,
    },

    /**
     * Caso 5: Cierre con "Cerrar remisión = No" (seguimiento dentro del Form Cierre)
     * session_count >= 3, pero el profesional decide no cerrar aún.
     * S3 permanece oculta.
     */
    caso5: {
      _meta: {
        title: 'Cierre — sin cerrar remisión (sesión 4)',
        subtitle: 'session_count = 3. "Cerrar remisión = No" → S3 oculta. Actúa como seguimiento. status no cambia.',
        badge: { color: '#f97316', label: 'Caso 5' },
      },
      formId:       FORM_IDS.CIERRE,
      submissionId: SUBMISSION_ID,
      formType:     'CIERRE',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  true,
        yaHizoPrimeraAtencion: true,
        sessionCount:          3,
        status:                'en_gestion',
      },
      formState: {},
      canEdit: true,
      sections: FORM_SECTIONS.CIERRE,
      _previewState: { cerrarRemision: false },
    },

    /**
     * Caso 6: Cierre definitivo — "Cerrar remisión = Sí"
     * S3 (Cierre) se muestra; al guardar status → cerrado.
     */
    caso6: {
      _meta: {
        title: 'Cierre definitivo (sesión 4)',
        subtitle: 'session_count = 3. "Cerrar remisión = Sí" → S3 visible. Al guardar: status = cerrado.',
        badge: { color: '#ef4444', label: 'Caso 6' },
      },
      formId:       FORM_IDS.CIERRE,
      submissionId: SUBMISSION_ID,
      formType:     'CIERRE',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  true,
        yaHizoPrimeraAtencion: true,
        sessionCount:          3,
        status:                'en_gestion',
      },
      formState: {},
      canEdit: true,
      sections: FORM_SECTIONS.CIERRE,
      _previewState: { cerrarRemision: true },
    },

    /**
     * Caso 7: Proceso cerrado (solo lectura)
     */
    caso7: {
      _meta: {
        title: 'Proceso cerrado (solo lectura)',
        subtitle: 'status = cerrado. canEdit = false. Se muestra el formulario en modo lectura.',
        badge: { color: '#6b7280', label: 'Caso 7' },
      },
      formId:       FORM_IDS.CIERRE,
      submissionId: SUBMISSION_ID,
      formType:     'CIERRE',
      victimInfo:   VICTIM_INFO,
      psicosocialState: {
        yaHizoPrimerContacto:  true,
        yaHizoPrimeraAtencion: true,
        sessionCount:          4,
        status:                'cerrado',
      },
      formState: {},
      canEdit: false,
      sections: FORM_SECTIONS.CIERRE,
    },
  };

  // ─── Labels de estado para UI ─────────────────────────────────────────────
  const STATUS_LABELS = {
    abierto:       { label: 'Abierto',       color: '#f59e0b', bg: '#fffbeb' },
    en_gestion:    { label: 'En gestión',    color: '#3b82f6', bg: '#eff6ff' },
    en_devolucion: { label: 'En devolución', color: '#f97316', bg: '#fff7ed' },
    cerrado:       { label: 'Cerrado',       color: '#6b7280', bg: '#f9fafb' },
  };

  const FORM_TYPE_LABELS = {
    PRIMER_CONTACTO:   'Primer Contacto',
    PRIMERA_ATENCION:  'Primera Atención',
    SEGUIMIENTO:       'Seguimiento',
    CIERRE:            'Cierre',
  };

  const QUESTION_TYPE_ICONS = {
    single:   '○',
    boolean:  '◉',
    multiple: '☑',
    text:     '≡',
    date:     '📅',
  };

  // ─── API pública ──────────────────────────────────────────────────────────
  global.PSICO_MOCK = {
    getCaso(id) { return CASOS[id] || null; },
    getAllCasos() { return Object.entries(CASOS).map(([id, c]) => ({ id, ...c })); },
    STATUS_LABELS,
    FORM_TYPE_LABELS,
    QUESTION_TYPE_ICONS,
    FORM_IDS,
  };

})(window);
