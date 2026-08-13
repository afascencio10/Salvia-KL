/**
 * dinamic-form
 * Componente reutilizable para renderizar un formulario dinámico por secciones.
 * Soporta: single, multiple, boolean, dropdown, text, date, time, datetime, number, repeater, info.
 * Estilos basados en el prototipo Salvia.
 * Props:
 *   - formId      (String, required): ID del formulario a cargar
 *   - submissionId (String, optional): ID del intento existente a continuar
 */

/* ─── Inyección de estilos ──────────────────────────────────────────────── */
(function injectDinamicFormStyles() {
    if (document.getElementById('dinamic-form-styles')) return;
    const style = document.createElement('style');
    style.id = 'dinamic-form-styles';
    style.textContent = `
        /* Tailwind colors used:
           violet-50 #f5f3ff  violet-200 #ddd6fe  violet-400 #a78bfa
           violet-500 #8b5cf6  violet-600 #7c3aed  violet-700 #6d28d9
           gray-50 #f9fafb  gray-100 #f3f4f6  gray-200 #e5e7eb
           gray-300 #d1d5db  gray-400 #9ca3af  gray-500 #6b7280
           gray-600 #4b5563  gray-700 #374151  gray-900 #111827
           green-500 #22c55e  red-400 #f87171
           shadow-sm: 0 1px 2px 0 rgb(0 0 0/.05) */

        .df-wrapper {
            position: relative;
            display: flex;
            gap: 24px;
            align-items: flex-start;
        }

        /* ── Guardar borrador (flotante) ── */
        .df-draft-save-btn {
            position: absolute;
            top: -44px;
            right: 0;
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            font-weight: 500;
            background: #fff;
            color: #7c3aed;
            border: 1px solid #ddd6fe;
            border-radius: 8px;
            padding: 7px 14px;
            cursor: pointer;
            box-shadow: 0 2px 6px rgba(0,0,0,0.08);
            transition: background 0.15s, border-color 0.15s, color 0.15s;
            z-index: 5;
        }
        .df-draft-save-btn:hover:not(:disabled) { background: #f5f3ff; border-color: #a78bfa; }
        .df-draft-save-btn:disabled { opacity: 0.5; cursor: not-allowed; }
        .df-draft-save-btn.df-draft-saved { color: #16a34a; border-color: #bbf7d0; background: #f0fdf4; }
        .df-draft-save-btn.df-draft-error { color: #dc2626; border-color: #fecaca; background: #fef2f2; }

        /* ── Sidebar ── */
        .df-sidebar {
            width: 208px;
            flex-shrink: 0;
            background: #fff;
            border: 1px solid #e5e7eb;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 1px 2px 0 rgb(0 0 0/.05);
        }
        .df-sidebar-header {
            padding: 12px 16px;
            border-bottom: 1px solid #f3f4f6;
        }
        .df-sidebar-title {
            font-size: 11px;
            font-weight: 600;
            color: #9ca3af;
            text-transform: uppercase;
            letter-spacing: 0.07em;
            margin: 0;
        }
        .df-sidebar-nav {
            padding: 8px;
            display: flex;
            flex-direction: column;
            gap: 2px;
        }
        .df-section-item {
            width: 100%;
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px 12px;
            border-radius: 8px;
            font-size: 14px;
            text-align: left;
            background: transparent;
            border: none;
            cursor: pointer;
            transition: background 0.15s, color 0.15s;
            line-height: 1.3;
        }
        .df-section-item.active    { background: #f5f3ff; color: #6d28d9; font-weight: 600; }
        .df-section-item.completed { color: #4b5563; }
        .df-section-item.completed:hover { background: #f9fafb; }
        .df-section-item.enabled   { color: #4b5563; }
        .df-section-item.enabled:hover { background: #f9fafb; }
        .df-section-item.pending   { color: #d1d5db; cursor: not-allowed; }
        .df-badge {
            width: 20px; height: 20px;
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            font-size: 11px; font-weight: 700;
            flex-shrink: 0;
        }
        .df-badge.active    { background: #7c3aed; color: #fff; }
        .df-badge.completed { background: #22c55e; color: #fff; }
        .df-badge.enabled   { background: #e5e7eb; color: #6b7280; }
        .df-badge.pending   { background: #f3f4f6; color: #d1d5db; }
        .df-progress-wrap {
            padding: 12px 16px;
            border-top: 1px solid #f3f4f6;
        }
        .df-progress-label {
            display: flex; justify-content: space-between;
            font-size: 11px; color: #9ca3af; margin-bottom: 6px;
        }
        .df-progress-bar-bg {
            height: 6px; background: #f3f4f6;
            border-radius: 999px; overflow: hidden;
        }
        .df-progress-bar-fill {
            height: 100%; background: #8b5cf6;
            border-radius: 999px; transition: width 0.3s ease;
        }

        /* ── Saving overlay ── */
        .df-saving-overlay {
            position: absolute; inset: 0;
            background: rgba(255,255,255,0.75);
            border-radius: 12px;
            display: flex; align-items: center; justify-content: center;
            gap: 10px;
            font-size: 14px; color: #6b7280;
            z-index: 10;
        }

        /* ── Main panel ── */
        .df-main {
            flex: 1; min-width: 0;
            position: relative;
            background: #fff;
            border: 1px solid #e5e7eb;
            border-radius: 12px;
            box-shadow: 0 1px 2px 0 rgb(0 0 0/.05);
            overflow: hidden;
        }
        .df-section-header {
            padding: 14px 24px;
            border-bottom: 1px solid #f3f4f6;
            background: #f9fafb;
        }
        .df-section-header-row {
            display: flex; align-items: center; gap: 8px; margin-bottom: 2px;
        }
        .df-section-number {
            width: 24px; height: 24px; border-radius: 50%;
            background: #7c3aed; color: #fff;
            font-size: 11px; font-weight: 700;
            display: flex; align-items: center; justify-content: center;
            flex-shrink: 0;
        }
        .df-section-title {
            font-size: 16px; font-weight: 700; color: #111827; margin: 0;
        }
        .df-section-desc {
            font-size: 14px; color: #6b7280; margin: 0; padding-left: 32px;
        }

        /* ── Questions ── */
        .df-questions {
            padding: 20px 24px;
            display: flex; flex-direction: column; gap: 16px;
        }
        .df-question { display: flex; flex-direction: column; }

        /* Divisor debajo de cada pregunta de nivel superior */
        .df-questions > .df-question {
            border-bottom: 1px solid #f0f0f0;
            padding-bottom: 14px;
            margin-bottom: 20px;
        }

        /* Divisor y separación entre preguntas dentro del repeater */
        .df-repeater-item .df-question {
            border-bottom: 1px solid #f0f0f0;
            padding-bottom: 14px;
            margin-bottom: 20px;
        }
        .df-repeater-item .df-question:last-child {
            border-bottom: none;
            padding-bottom: 0;
            margin-bottom: 0;
        }

        /* Estilos base del label (aplica a todos, incluyendo repeater interior) */
        .df-question > label,
        .df-repeater-item label {
            display: block;
            font-size: 14px; font-weight: 500; color: #374151;
            margin-bottom: 4px; margin-top: 0;
            min-height: unset;
        }

        /* Label de preguntas de nivel superior: fuente +25% + ícono */
        .df-questions > .df-question > label,
        .df-repeater-item .df-question > label {
            display: flex;
            align-items: flex-start;
            gap: 8px;
            font-size: 17.5px;
            margin-bottom: 8px;
            line-height: 1.4;
        }
        .df-questions > .df-question > label::before,
        .df-repeater-item .df-question > label::before {
            content: '';
            display: block;
            width: 7px; height: 7px;
            min-width: 7px;
            border-radius: 50%;
            background: #7c3aed;
            margin-top: 6px;
            flex-shrink: 0;
        }

        .df-required { color: #f87171; margin-left: 4px; }

        /* text / date / datetime */
        .df-input, .df-textarea, .df-select {
            width: 100%; font-size: 14px;
            border: 1px solid #e5e7eb; border-radius: 8px;
            padding: 7px 12px; background: #fff; color: #111827;
            outline: none; box-sizing: border-box;
            transition: box-shadow 0.15s, border-color 0.15s;
            font-family: inherit;
        }
        .df-textarea { min-height: 72px; resize: vertical; }
        .df-input:focus, .df-textarea:focus, .df-select:focus {
            border-color: #a78bfa;
            box-shadow: 0 0 0 2px #ddd6fe;
        }

        /* boolean */
        .df-bool-group { display: flex; gap: 12px; }
        .df-bool-btn {
            flex: 1; padding: 7px 0;
            border-radius: 8px; font-size: 14px; font-weight: 500;
            border: 1px solid #e5e7eb; background: #fff; color: #4b5563;
            cursor: pointer; transition: background 0.15s, border-color 0.15s, color 0.15s;
        }
        .df-bool-btn:hover:not(.df-selected) { border-color: #ddd6fe; background: #f5f3ff; }
        .df-bool-btn.df-selected { background: #7c3aed; border-color: #7c3aed; color: #fff; }

        /* single / multiple (chips) */
        .df-btn-group { display: flex; flex-wrap: wrap; gap: 8px; }
        .df-opt-btn {
            padding: 6px 12px; border-radius: 8px;
            font-size: 14px; font-weight: 500;
            border: 1px solid #e5e7eb; background: #fff; color: #4b5563;
            cursor: pointer; transition: background 0.15s, border-color 0.15s, color 0.15s;
        }
        .df-opt-btn:hover:not(.df-selected) { border-color: #ddd6fe; background: #f5f3ff; }
        .df-opt-btn.df-selected { background: #7c3aed; border-color: #7c3aed; color: #fff; }

        /* repeater */
        .df-repeater { display: flex; flex-direction: column; gap: 12px; }
        .df-repeater-item {
            border: 1px solid #e5e7eb; border-radius: 12px;
            padding: 16px; background: #f9fafb;
            display: flex; flex-direction: column; gap: 12px;
        }
        .df-repeater-item-header {
            display: flex; align-items: center; justify-content: space-between;
        }
        .df-repeater-item-title {
            font-size: 12px; font-weight: 600; color: #6b7280; margin: 0;
        }
        .df-repeater-delete {
            font-size: 12px; color: #f87171; background: none; border: none;
            cursor: pointer; padding: 0; transition: color 0.15s;
        }
        .df-repeater-delete:hover { color: #ef4444; }
        .df-repeater-add {
            width: 100%; padding: 12px;
            border: 2px dashed #ddd6fe; border-radius: 12px;
            background: transparent; color: #7c3aed;
            font-size: 14px; font-weight: 500;
            cursor: pointer; transition: border-color 0.15s, background 0.15s;
            box-sizing: border-box;
        }
        .df-repeater-add:hover { border-color: #a78bfa; background: #f5f3ff; }

        /* multiple hint */
        .df-hint  { font-size: 12px; color: #9ca3af; font-style: italic; }
        .df-error { font-size: 12px; color: #f87171; margin-top: 4px; }

        /* ── Navigation ── */
        .df-nav {
            padding: 12px 24px;
            border-top: 1px solid #f3f4f6; background: #f9fafb;
            display: flex; align-items: center; justify-content: space-between;
        }
        .df-btn-prev {
            font-size: 14px; color: #6b7280;
            border: 1px solid #e5e7eb; border-radius: 8px;
            padding: 8px 16px; background: transparent; cursor: pointer;
            transition: background 0.15s;
        }
        .df-btn-prev:hover:not(:disabled) { background: #fff; }
        .df-btn-prev:disabled { opacity: 0.3; cursor: not-allowed; }
        .df-btn-next {
            font-size: 14px; font-weight: 500;
            background: #7c3aed; color: #fff;
            border: 1px solid #7c3aed; border-radius: 8px;
            padding: 8px 20px; cursor: pointer;
            transition: background 0.15s, border-color 0.15s;
        }
        .df-btn-next:hover { background: #6d28d9; border-color: #6d28d9; }

        @keyframes df-spin { to { transform: rotate(360deg); } }

        /* ── Validation popup ── */
        .df-popup-backdrop {
            position: fixed; inset: 0;
            background: rgba(0,0,0,0.35);
            z-index: 100;
            display: flex; align-items: center; justify-content: center;
        }
        .df-popup {
            background: #fff;
            border-radius: 12px;
            padding: 28px 32px;
            max-width: 360px; width: 90%;
            box-shadow: 0 8px 32px rgba(0,0,0,0.18);
            display: flex; flex-direction: column; align-items: center; gap: 16px;
            text-align: center;
        }
        .df-popup-icon {
            width: 44px; height: 44px; border-radius: 50%;
            background: #fef2f2;
            display: flex; align-items: center; justify-content: center;
            font-size: 20px;
        }
        .df-popup-title {
            font-size: 15px; font-weight: 600; color: #111827; margin: 0;
        }
        .df-popup-btn {
            padding: 8px 24px; border-radius: 8px;
            background: #7c3aed; color: #fff;
            border: none; font-size: 14px; font-weight: 500;
            cursor: pointer; transition: background 0.15s;
        }
        .df-popup-btn:hover { background: #6d28d9; }

        /* ── Info banner ── */
        .df-info-banner {
            display: flex;
            align-items: flex-start;
            gap: 10px;
            background: #eff6ff;
            border: 1px solid #bfdbfe;
            border-radius: 8px;
            padding: 12px 16px;
            font-size: 14px;
            color: #1e40af;
            line-height: 1.5;
            white-space: pre-wrap;
        }
        .df-info-banner-icon {
            flex-shrink: 0;
            font-size: 16px;
            margin-top: 1px;
        }

        /* ── Read-only ── */
        .df-input:disabled, .df-textarea:disabled, .df-select:disabled {
            background: #f9fafb; color: #6b7280; cursor: default; opacity: 1;
        }
        .df-opt-btn:disabled, .df-bool-btn:disabled { cursor: default; opacity: 1; }
        .df-opt-btn:disabled:not(.df-selected),
        .df-bool-btn:disabled:not(.df-selected) {
            background: #f9fafb; border-color: #e5e7eb; color: #9ca3af;
        }
        .df-readonly-badge {
            display: inline-flex; align-items: center; gap: 4px;
            font-size: 11px; font-weight: 600; color: #9ca3af;
            background: #f3f4f6; border: 1px solid #e5e7eb;
            border-radius: 6px; padding: 2px 8px;
            text-transform: uppercase; letter-spacing: 0.05em;
        }
    `;
    document.head.appendChild(style);
})();

/* ─── Schema de secciones ───────────────────────────────────────────────── */
var DF_BARRERAS_POR_SECTOR = {
    'Salud': [
        'No Aplica',
        'Demoras en asignación de citas y continuidad de tratamientos',
        'Dificultades de aseguramiento o afiliación en salud',
        'Falta de activación de protocolos para violencia sexual',
        'Negación en los servicios de urgencias',
        'Negativa en el acceso a la Interrupción Voluntaria del Embarazo (IVE)',
        'Otras barreras en salud',
    ],
    'Justicia': [
        'No Aplica',
        'Ausencia de representación judicial',
        'Cargas probatorias injustificadas o excesivas',
        'Falta de celeridad en la investigación',
        'Negativa para recibir la denuncia',
        'Tipificación errónea del delito',
        'Otras barreras en justicia',
    ],
    'Protección': [
        'No Aplica',
        'Demoras en la emisión de medidas de protección urgentes',
        'Fallas en la valoración del riesgo',
        'Incumplimiento de medidas de protección sin consecuencias',
        'Medidas de protección ineficaces',
        'Otras barreras en protección',
    ],
};

var DF_PROFESIONALES = [
    'Ana María Mojica Quiroz',
    'Andres Eduardo Barbosa Deaquiz',
    'Andres Felipe Suarez Cabra',
    'Cristian Alexander Rodriguez Alarcon',
    'Daniel Esteban Acosta Rodríguez',
    'Dayan Vargas Sánchez',
    'José Felipe Calixto',
    'Luz Virginia Gomez Morales',
    'Tatiana Geraldine Montalván Caicedo',
];

var DF_SCHEMA = [
    {
        id: 1,
        title: 'Inicio de Seguimiento',
        description: 'Información básica sobre cómo se realizó el seguimiento',
        questions: [
            {
                id: 'nivel_riesgo',
                type: 'single',
                label: '¿Cuál fue el nivel de riesgo identificado en el registro?',
                options: ['Extremo', 'Alto', 'Moderado', 'Bajo', 'Sin Riesgo'],
                required: true,
            },
            {
                id: 'equipo_atencion',
                type: 'single',
                label: 'Equipo que realiza la atención',
                options: ['Agente Integral', 'Riesgo de Feminicidio', 'Seguimiento General', 'Notificación Salvia', 'Psicosocial', 'Masculinidades'],
                required: true,
            },
            {
                id: 'nombre_profesional',
                type: 'dropdown',
                label: 'Nombre del profesional que realiza la atención',
                options: DF_PROFESIONALES,
                required: true,
            },
            {
                id: 'efectividad_llamada',
                type: 'single',
                label: 'Efectividad de la llamada',
                options: [
                    'Llamada Efectiva - Se logra comunicación',
                    'Llamada NO efectiva - NO hay comunicación',
                ],
                required: true,
            },
            {
                id: 'numero_seguimiento',
                type: 'single',
                label: '¿Qué número de seguimiento está registrando?',
                options: ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10'],
                required: true,
            },
            {
                id: 'fecha_seguimiento',
                type: 'date',
                label: 'Fecha del seguimiento',
                required: true,
            },
            {
                id: 'hora_inicio',
                type: 'datetime',
                label: 'Fecha y hora de inicio del seguimiento',
                required: true,
            },
            {
                id: 'codigo_caso',
                type: 'text',
                label: 'Código interno del caso',
                placeholder: 'Ej. SAL-2024-0042',
                required: false,
            },
        ],
    },
    {
        id: 2,
        title: 'Valoración del Riesgo',
        description: 'Evaluación de nuevos hechos y respondiente del seguimiento',
        questions: [
            {
                id: 'respondiente',
                type: 'single',
                label: 'Respondiente del seguimiento',
                options: [
                    'Contacto con la víctima',
                    'Contacto indirecto - Familiar o persona conocida',
                    'Contacto Con Institución',
                ],
                required: true,
            },
            {
                id: 'nuevos_hechos',
                type: 'boolean',
                label: '¿Se han presentado nuevos hechos de violencia desde el último seguimiento?',
                required: true,
            },
            {
                id: 'descripcion_hechos',
                type: 'text',
                multiline: true,
                label: 'Descripción de nuevos hechos de violencia (Tiempo/Modo/Lugar)',
                required: true,
                dependsOn: { id: 'nuevos_hechos', value: true },
            },
        ],
    },
    {
        id: 3,
        title: 'Seguimiento de Caso',
        description: 'Gestión realizada y variaciones identificadas en el caso',
        questions: [
            {
                id: 'info_psicosocial',
                type: 'text',
                multiline: true,
                label: 'Información psicosocial relevante para el seguimiento',
                required: true,
            },
            {
                id: 'gestion_realizada',
                type: 'text',
                multiline: true,
                label: 'Gestión realizada en el seguimiento',
                placeholder: 'Acciones realizadas sobre la ruta, nivel psicosocial, medidas de emergencia...',
                required: true,
            },
            {
                id: 'variacion_riesgo',
                type: 'boolean',
                label: '¿Hubo variación en el nivel de riesgo desde el último seguimiento?',
                required: true,
            },
            {
                id: 'descripcion_variacion',
                type: 'text',
                multiline: true,
                label: 'Describa la variación del riesgo',
                placeholder: 'Registrar si el riesgo aumentó o disminuyó, señales de alerta, cambios en el contexto...',
                required: true,
                dependsOn: { id: 'variacion_riesgo', value: true },
            },
            {
                id: 'necesidades_derivacion',
                type: 'boolean',
                label: 'Se generaron necesidades inmediatas que requieran derivación a los equipos Salvia',
                required: true,
            },
            {
                id: 'equipos_derivacion',
                type: 'multiple',
                label: '¿Cuáles equipos?',
                options: ['Medidas de Emergencia', 'Atención Psicosocial', 'Estabilización'],
                required: true,
                dependsOn: { id: 'necesidades_derivacion', value: true },
            },
            {
                id: 'medidas_emergencia',
                type: 'multiple',
                label: '¿Cuáles medidas de emergencia?',
                options: ['Alojamiento', 'Transporte', 'Alimentación', 'Vestuario', 'Apoyo psicosocial', 'Otras ME'],
                required: true,
                dependsOn: { id: 'equipos_derivacion', operator: 'includes', value: 'Medidas de Emergencia' },
            },
        ],
    },
    {
        id: 4,
        title: 'Identificación de Barreras',
        description: 'Registro de barreras institucionales identificadas durante el seguimiento',
        questions: [
            {
                id: 'barreras_list',
                type: 'repeater',
                label: 'Barreras',
                addLabel: '+ Agregar barrera',
                fields: [
                    {
                        id: 'sector',
                        type: 'dropdown',
                        label: 'Sector de la barrera',
                        options: ['Salud', 'Justicia', 'Protección'],
                        required: true,
                    },
                    {
                        id: 'barreras_sector',
                        type: 'multiple',
                        label: 'Barreras identificadas en este sector',
                        required: true,
                        dynamicOptions: {
                            dependsOn: 'sector',
                            map: DF_BARRERAS_POR_SECTOR,
                        },
                    },
                ],
            },
        ],
    },
    {
        id: 5,
        title: 'Empalme de Seguimiento',
        description: 'Registro de empalme y riesgo inminente',
        questions: [
            {
                id: 'realiza_empalme',
                type: 'boolean',
                label: 'Se realiza empalme del seguimiento',
                required: true,
            },
            {
                id: 'riesgo_inminente',
                type: 'single',
                label: '¿Este seguimiento es de riesgo inminente?',
                options: [
                    'No',
                    'Sí. Requiere seguimiento en 4 horas',
                    'Sí. Requiere seguimiento en 8 horas',
                ],
                required: true,
                dependsOn: { id: 'realiza_empalme', value: true },
            },
            {
                id: 'fecha_hora_accion',
                type: 'datetime',
                label: 'Fecha y hora de la acción a realizar',
                required: true,
                dependsOn: { id: 'realiza_empalme', value: true },
            },
            {
                id: 'equipo_empalme',
                type: 'single',
                label: 'Equipo que realiza la atención',
                options: ['Riesgo de Feminicidio', 'Seguimiento General', 'Agente Integral'],
                required: true,
                dependsOn: { id: 'realiza_empalme', value: true },
            },
            {
                id: 'profesional_empalme',
                type: 'dropdown',
                label: 'Nombre del profesional para empalme',
                options: DF_PROFESIONALES,
                required: true,
                dependsOn: { id: 'realiza_empalme', value: true },
            },
        ],
    },
];

/* ─── getObjectProperty ──────────────────────────────────────────────────────
 * Recorre un objeto anidado siguiendo un path con puntos ("a.b.c").
 * Retorna undefined si algún segmento no existe.
 */
function getObjectProperty(object, path) {
    if (path === undefined || path === null) {
        return object;
    }
    const parts = path.split('.');
    return parts.reduce((obj, key) => {
        if (!obj) return undefined;
        return obj[key];
    }, object);
}

/* ─── applyRenderModifications ──────────────────────────────────────────────
 * Aplica las render_modification de un campo de texto.
 * Orden: SETs primero (reemplazan todo el texto), REPLACEs después (subcadena).
 * Solo aplica si el valor resuelto del formState no es '' ni undefined/null.
 *
 * entryIndex (opcional): índice 0-based de la entry del repeater.
 * Si el statePath contiene '{_entryIndex}', se reemplaza por el índice antes
 * de resolver el path. Ej: "barriers.{_entryIndex}.name" con entryIndex=1
 * → "barriers.1.name".
 * Si entryIndex no se pasa o el path no contiene el placeholder, se resuelve
 * el path tal cual (comportamiento estándar).
 */
function applyRenderModifications(text, modifications, fieldName, formState, entryIndex) {
    const applicable = (modifications || []).filter(m => m.targetField === fieldName);
    if (!applicable.length) return text;

    let result = text || '';

    const resolvePath = (statePath) => {
        const path = (entryIndex !== undefined && entryIndex !== null)
            ? statePath.replace('{_entryIndex}', String(entryIndex))
            : statePath;
        const resolved = getObjectProperty(formState || {}, path);
        console.log('[applyRenderModifications:resolvePath]', { statePath, entryIndex, resolvedPath: path, resolved });
        return resolved;
    };

    // 1. SET — reemplaza todo el texto
    for (const mod of applicable.filter(m => m.modificationType === 'SET')) {
        const raw = resolvePath(mod.statePath);
        if (raw === undefined || raw === null || String(raw) === '') continue;
        result = String(raw);
    }

    // 2. REPLACE — reemplaza subcadena específica (todas las ocurrencias)
    for (const mod of applicable.filter(m => m.modificationType === 'REPLACE')) {
        const raw = resolvePath(mod.statePath);
        if (raw === undefined || raw === null || String(raw) === '') continue;
        if (!mod.searchString) continue;
        result = result.split(mod.searchString).join(String(raw));
    }

    return result;
}

/* ─── resolveQuestionOptions ─────────────────────────────────────────────────
 * Devuelve el array de opciones a renderizar para una pregunta.
 *
 * Si question.stateOptionsPath está definido, es la fuente de verdad:
 *   - Resuelve el path en formState (con soporte de {_entryIndex}).
 *   - Valida que el resultado sea un array de objetos { label, value }.
 *   - Si no resuelve de forma válida → devuelve [] (sin opciones).
 * Si no está definido → devuelve question.options (comportamiento estándar).
 *
 * entryIndex (opcional): índice 0-based de la entry del repeater.
 */
function resolveQuestionOptions(question, formState, entryIndex) {
    if (!question.stateOptionsPath) return question.options || [];

    const path = (entryIndex !== undefined && entryIndex !== null)
        ? question.stateOptionsPath.replace('{_entryIndex}', String(entryIndex))
        : question.stateOptionsPath;

    const raw = getObjectProperty(formState || {}, path);

    if (!Array.isArray(raw) || raw.length === 0) return [];

    // Validar que cada elemento tenga label y value
    const valid = raw.every(o => o && typeof o === 'object' && 'label' in o && 'value' in o);
    return valid ? raw : [];
}

/* ─── checkVisibility ────────────────────────────────────────────────────────
 * Réplica exacta en JS de checkVisibility (Go — form_service.go).
 * formState: objeto externo opcional para condiciones de tipo triggerStatePath.
 */

function findQuestion(fs, questionId) {
    for (const sec of fs.sections) {
        for (const q of (sec.questions || [])) {
            if (q.id === questionId) {
                return { description: q.description, sectionOrder: sec.order, repeaterGroupId: null, orderInSection: q.order, orderInRepeater: 0 };
            }
        }
        for (const r of (sec.repeaters || [])) {
            for (const q of (r.questions || [])) {
                if (q.id === questionId) {
                    return { description: q.description, sectionOrder: sec.order, repeaterGroupId: r.id, orderInSection: r.order, orderInRepeater: q.order };
                }
            }
        }
    }
    return null;
}

function isBeforeInFlow(trigger, item) {
    return trigger.sectionOrder < item.sectionOrder ||
        (trigger.sectionOrder === item.sectionOrder && trigger.orderInSection < item.orderInSection);
}

function checkVisibility(fs, submission, entryAnswers, item, formState) {
    if (!item.conditions || item.conditions.length === 0) return { name: item.name, visible: true };

    const directAnswers = (submission && submission.directAnswers) || [];
    const applicable = [];

    for (const cond of item.conditions) {
        // ── Condición por formState (triggerStatePath) ──────────────────────
        if (cond.triggerStatePath) {
            applicable.push(cond);
            continue;
        }
        // ── Condición por pregunta (triggerQuestionId) ──────────────────────
        const trigger = findQuestion(fs, cond.triggerQuestionId);
        if (!trigger) continue;
        if (trigger.repeaterGroupId !== null) {
            if (!item.repeaterGroupId || item.repeaterGroupId !== trigger.repeaterGroupId) continue;
            if (!entryAnswers) continue;
            if (trigger.orderInRepeater >= item.orderInRepeater) continue;
        } else {
            if (!isBeforeInFlow(trigger, item)) continue;
        }
        applicable.push(cond);
    }

    if (applicable.length === 0) return { name: item.name, visible: true };

    const failedConditions = [];
    for (const cond of applicable) {
        const trigVal = cond.triggerValue ?? '';
        let condFailed = false;

        if (cond.triggerStatePath) {
            // Resolver valor desde formState (igual que resolveStatePath en Go)
            const raw = getObjectProperty(formState || {}, cond.triggerStatePath);
            const stateVal = (raw === undefined || raw === null) ? '' : String(raw);
            switch (cond.operator.toLowerCase()) {
                case 'equals':             condFailed = stateVal !== trigVal; break;
                case 'includes':
                case 'contains':           condFailed = !stateVal.includes(trigVal); break;
                default:                   condFailed = false;
            }
            if (condFailed) failedConditions.push({ ...cond });
        } else {
            const trigger = findQuestion(fs, cond.triggerQuestionId);
            const pool    = trigger.repeaterGroupId !== null ? (entryAnswers || []) : directAnswers;
            const answer  = pool.find(a => a.questionId === cond.triggerQuestionId) || null;
            condFailed = !answer;
            if (answer) {
                switch (cond.operator.toLowerCase()) {
                    case 'equals':             condFailed = answer.value !== trigVal; break;
                    case 'includes':
                    case 'contains':           condFailed = !answer.value.includes(trigVal); break;
                    default:                   condFailed = false;
                }
            }
            if (condFailed) failedConditions.push({ ...cond, triggerQuestionName: trigger.description });
        }
    }

    return failedConditions.length > 0
        ? { name: item.name, visible: false, failedConditions }
        : { name: item.name, visible: true };
}

function sectionItem(sec) {
    return { id: sec.id, name: sec.name, conditions: sec.conditions || [], sectionOrder: sec.order, orderInSection: 0, repeaterGroupId: null, orderInRepeater: 0 };
}
function repeaterItem(r, sectionOrder) {
    return { id: r.id, name: r.name, conditions: r.conditions || [], sectionOrder, orderInSection: r.order, repeaterGroupId: null, orderInRepeater: 0 };
}
function directQuestionItem(q, sectionOrder) {
    return { id: q.id, name: q.description, conditions: q.conditions || [], sectionOrder, orderInSection: q.order, repeaterGroupId: null, orderInRepeater: 0 };
}
function repeaterQuestionItem(q, sectionOrder, repeaterOrder, repeaterGroupId) {
    return { id: q.id, name: q.description, conditions: q.conditions || [], sectionOrder, orderInSection: repeaterOrder, repeaterGroupId, orderInRepeater: q.order };
}

/* ─── buildSectionRenderData ─────────────────────────────────────────────── */
/**
 * Construye la estructura de datos para renderizar una sección:
 * formItems es un array ordenado que mezcla preguntas directas y repeaterGroups,
 * cada uno con su(s) respuesta(s) asociada(s) del submission actual.
 *
 * @param {Object} section    - SectionStructure (de formStructure.sections)
 * @param {Object|null} submission - SubmissionStructure (de formSubmission)
 * @returns {{ section, formItems }}
 */
function buildSectionRenderData(section, submission, fs, formState) {
    const sectionOrder = section.order;

    // ── 1. Índices de respuestas ────────────────────────────────────────────
    const answerByQuestionId  = {};
    const entriesByRepeaterId = {};

    if (submission) {
        (submission.directAnswers || []).forEach(a => {
            answerByQuestionId[a.questionId] = a;
        });
        (submission.repeaterEntries || []).forEach(entry => {
            if (!entriesByRepeaterId[entry.repeaterGroupId])
                entriesByRepeaterId[entry.repeaterGroupId] = [];
            entriesByRepeaterId[entry.repeaterGroupId].push(entry);
        });
    }

    // ── 2. Preguntas directas ───────────────────────────────────────────────
    const questionItems = (section.questions || []).map(q => {
        const visItem  = directQuestionItem(q, sectionOrder);
        const vis      = fs ? checkVisibility(fs, submission, null, visItem, formState) : { visible: true };
        return {
            type:      'question',
            order:     q.order,
            question:  q,
            answer:    answerByQuestionId[q.id] || null,
            isVisible: vis.visible,
        };
    });

    // ── 3. Repeater groups con sus entries y respuestas ────────────────────
    const repeaterItems = (section.repeaters || []).map(r => {
        const visItem  = repeaterItem(r, sectionOrder);
        const vis      = fs ? checkVisibility(fs, submission, null, visItem, formState) : { visible: true };

        let entries;

        if (r.stateItems) {
            // ── Repeater controlado por estado ──────────────────────────────
            // El array del formState es la fuente de verdad del conteo de entries.
            // Las entries guardadas en el submission se usan para precargar respuestas (por iteración).
            const stateArray = getObjectProperty(formState || {}, r.stateItems);

            if (!Array.isArray(stateArray) || stateArray.length === 0) {
                entries = [];
            } else {
                // Índice de entries guardadas por iteración para precargar respuestas
                const savedByIteration = {};
                (entriesByRepeaterId[r.id] || []).forEach(e => {
                    savedByIteration[e.iteration] = e;
                });

                entries = stateArray.map((_, idx) => {
                    const iteration  = idx + 1;
                    const savedEntry = savedByIteration[iteration];
                    // ID estable para temp entries — evita resetear localAnswers al re-renderizar
                    const entry = savedEntry || {
                        id:              `temp-state-${r.id}-${iteration}`,
                        repeaterGroupId: r.id,
                        iteration,
                        isTemp:          true,
                        answers:         [],
                    };
                    return {
                        entry,
                        questions: (r.questions || []).map(q => {
                            const qVisItem = repeaterQuestionItem(q, sectionOrder, r.order, r.id);
                            const qVis     = fs ? checkVisibility(fs, submission, entry.answers, qVisItem, formState) : { visible: true };
                            return {
                                question:  q,
                                answer:    (entry.answers || []).find(a => a.questionId === q.id) || null,
                                isVisible: qVis.visible,
                            };
                        }),
                    };
                });
            }
        } else {
            // ── Repeater estándar ────────────────────────────────────────────
            const rawEntries = entriesByRepeaterId[r.id] || [];
            entries = rawEntries.map(entry => ({
                entry,
                questions: (r.questions || []).map(q => {
                    const qVisItem = repeaterQuestionItem(q, sectionOrder, r.order, r.id);
                    const qVis     = fs ? checkVisibility(fs, submission, entry.answers, qVisItem, formState) : { visible: true };
                    return {
                        question:  q,
                        answer:    (entry.answers || []).find(a => a.questionId === q.id) || null,
                        isVisible: qVis.visible,
                    };
                }),
            }));
        }

        return {
            type:      'repeater',
            order:     r.order,
            repeater:  r,
            entries,
            isVisible: vis.visible,
        };
    });

    // ── 4. Mezclar y ordenar por order ─────────────────────────────────────
    const formItems = [...questionItems, ...repeaterItems]
        .sort((a, b) => a.order - b.order);

    return { section, formItems };
}

/* ─── Reglas de negocio del formulario "vivo" de seguimiento ────────────────
 * Pregunta "¿Cuáles equipos?" (remisión a equipos SALVIA): Medidas de Emergencia
 * y Atención Psicosocial no se pueden remitir al mismo tiempo, se debe elegir una.
 */
const Q_EQUIPOS_SALVIA        = 'e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc';
const OPT_MEDIDAS_EMERGENCIA  = 'medidas_emergencia';
const OPT_ATENCION_PSICO      = 'atencion_psico';

/* ─── Componente Vue ─────────────────────────────────────────────────────── */
app.component('dinamic-form', {
    delimiters: ['${', '}'],
    props: {
        formId:       { type: String,  required: true  },
        submissionId: { type: String,  required: false, default: null },
        canEdit:      { type: Boolean, required: false, default: true },
        formState:    { type: Object,  required: false, default: () => ({}) },
    },
    watch: {
        formState: {
            deep: true,
            handler(newState) {
                if (!this.currentSectionRender) return;
                this._reevaluateVisibility(this._buildTempSubmission(), newState);
                this._clearStaleOptionAnswers(newState);
            },
        },
    },
    data() {
        return {
            // ── Estado de carga ──────────────────────────────────────────
            loading: true,
            saving:  false,
            error:   null,

            // ── Estado principal (poblado por loadForm) ──────────────────
            formStructure:        null,   // { id, name, sections: [...] }
            formSubmission:       null,   // { id, directAnswers, repeaterEntries } | null
            currentSection:       null,   // SectionStructure de la sección activa
            currentSectionRender: null,   // resultado de buildSectionRenderData — para renderizar
            localAnswers:         {},     // respuestas editables en memoria { questionId: value }
            questionErrors:       {},     // errores de validación { answerKey: string | null }

            // ── Estado legacy (a migrar) ─────────────────────────────────
            currentIndex: 0,
            answers: {},
            sections: DF_SCHEMA,

            // ── Validación ───────────────────────────────────────────────
            showValidationError: false,
            repeaterErrors: {},

            // ── Guardar borrador ─────────────────────────────────────────
            draftSaved: false,
            draftError: false,
        };
    },
    computed: {
        // ── Secciones visibles del form real ─────────────────────────────
        visibleSections() {
            if (!this.formStructure) return [];
            return this.formStructure.sections.filter(s => s.isVisible);
        },
        firstUnansweredSection() {
            return this.visibleSections.find(s => !s.isAnswered) || null;
        },
        answeredCount() {
            return this.visibleSections.filter(s => s.isAnswered).length;
        },
        progressPercent() {
            if (!this.visibleSections.length) return 0;
            return Math.round((this.answeredCount / this.visibleSections.length) * 100);
        },

        isFirstSection() {
            if (!this.currentSection || !this.visibleSections.length) return true;
            return this.visibleSections[0].id === this.currentSection.id;
        },
        isLastSection() {
            if (!this.currentSection || !this.visibleSections.length) return false;
            return this.visibleSections[this.visibleSections.length - 1].id === this.currentSection.id;
        },

        // Hay al menos una respuesta cargada en la sección actual (controla si se
        // muestra el botón "Guardar borrador" — no tiene sentido si no hay nada que guardar).
        hasDraftAnswers() {
            return Object.values(this.localAnswers).some(v => v !== '' && v !== null && v !== undefined);
        },

        // ── Legacy (a migrar) ─────────────────────────────────────────────
        legacyCurrentSection() { return this.sections[this.currentIndex]; },
        totalSections()        { return this.sections.length; },
        isFirst() { return this.currentIndex === 0; },
        isLast()  { return this.currentIndex === this.totalSections - 1; },
    },
    async mounted() {
        await this.loadForm();
    },
    methods: {
        /* ── Render modifications ── */
        applyRenderModifications(text, modifications, fieldName, formState, entryIndex) {
            const applicable = (modifications || []).filter(m => m.targetField === fieldName);
            if (applicable.length && text && text.includes && text.includes('Barrera') || (applicable.length && entryIndex !== undefined)) {
                console.log('[applyRenderModifications]', {
                    fieldName,
                    text,
                    entryIndex,
                    statePaths: applicable.map(m => m.statePath),
                    formStateKeys: Object.keys(formState || {}),
                    currentBarriers: (formState || {}).currentBarriers,
                });
            }
            return applyRenderModifications(text, modifications, fieldName, formState, entryIndex);
        },
        resolveQuestionOptions(question, formState, entryIndex) {
            return resolveQuestionOptions(question, formState, entryIndex);
        },

        /* ── Sidebar ── */
        isSectionActive(section) {
            return this.currentSection && section.id === this.currentSection.id;
        },
        isSectionEnabled(section) {
            // Habilitada = respondida  O  es la primera sin responder
            return section.isAnswered ||
                   (this.firstUnansweredSection && section.id === this.firstUnansweredSection.id);
        },
        sidebarSectionClass(section) {
            if (this.isSectionActive(section))  return 'active';
            if (section.isAnswered)             return 'completed';
            if (this.isSectionEnabled(section)) return 'enabled';   // primera no respondida, no activa
            return 'pending';
        },
        sidebarBadgeClass(section) {
            if (section.isAnswered)             return 'completed';
            if (this.isSectionActive(section))  return 'active';
            if (this.isSectionEnabled(section)) return 'enabled';
            return 'pending';
        },
        goToSectionById(section) {
            if (!this.isSectionEnabled(section)) return;
            this.currentSection       = section;
            this.currentSectionRender = buildSectionRenderData(
                section,
                this.formSubmission,
                this.formStructure,
                this.formState,
            );
            this.initLocalAnswers();
        },

        /* ── Carga inicial del formulario ── */
        async loadForm() {
            this.loading = true;
            this.error   = null;
            console.log('[dinamic-form] mounted — formId:', this.formId, '| submissionId:', this.submissionId);
            try {
                const params = new URLSearchParams();
                if (this.submissionId) params.set('submissionId', this.submissionId);
                const fsJson = JSON.stringify(this.formState);
                if (fsJson && fsJson !== '{}') params.set('formState', fsJson);
                const url = `/api/v1/forms/${this.formId}/load` + (params.toString() ? '?' + params.toString() : '');

                console.log('[dinamic-form] GET', url);
                const res = await fetch(url);
                console.log('[dinamic-form] response status:', res.status);

                if (!res.ok) throw new Error(`Error ${res.status} al cargar el formulario`);

                const data = await res.json();
                console.log('[dinamic-form] data recibida:', data);

                this.formStructure        = data.formStructure;
                this.formSubmission       = data.formSubmission  ?? null;
                this.currentSection       = data.currentSection;
                this.currentSectionRender = buildSectionRenderData(
                    this.currentSection,
                    this.formSubmission,
                    this.formStructure,
                    this.formState,
                );
                this.initLocalAnswers();
                this._emitAnswersUpdated();

                console.log('[dinamic-form] estado actualizado:');
                console.log('  formStructure       :', this.formStructure);
                console.log('  formSubmission      :', this.formSubmission);
                console.log('  currentSection      :', this.currentSection);
                console.log('  currentSectionRender:', this.currentSectionRender);
            } catch (e) {
                console.error('[dinamic-form] error en loadForm:', e);
                this.error = e.message;
            } finally {
                this.loading = false;
                console.log('[dinamic-form] loading finalizado');
            }
        },

        getSectionState(index) {
            if (index < this.currentIndex) return 'completed';
            if (index === this.currentIndex) return 'active';
            return 'pending';
        },

        /* ── Respuestas planas (por sección) ── */
        getAnswer(questionId) {
            const key = this.currentSection.id + '_' + questionId;
            return this.answers[key] !== undefined ? this.answers[key] : null;
        },
        setAnswer(questionId, value) {
            const key = this.currentSection.id + '_' + questionId;
            this.answers = { ...this.answers, [key]: value };
        },

        /* ── Visibilidad condicional ── */
        isVisible(question, contextAnswers) {
            if (!question.dependsOn) return true;
            const dep = question.dependsOn;
            const ctx = contextAnswers || this.sectionAnswers();
            const val = ctx[dep.id] !== undefined ? ctx[dep.id] : null;
            if (dep.operator === 'includes') {
                return Array.isArray(val) && val.includes(dep.value);
            }
            return val === dep.value;
        },
        sectionAnswers() {
            const prefix = this.currentSection.id + '_';
            const result = {};
            Object.keys(this.answers).forEach(k => {
                if (k.startsWith(prefix)) result[k.slice(prefix.length)] = this.answers[k];
            });
            return result;
        },

        /* ── Opciones dinámicas de campos multiple ── */
        getDynamicOptions(question, contextAnswers) {
            if (!question.dynamicOptions) return question.options || [];
            const depVal = contextAnswers ? contextAnswers[question.dynamicOptions.dependsOn] : null;
            return question.dynamicOptions.map[depVal] || [];
        },

        /* ── Toggle para campos multiple ── */
        toggleMultiple(questionId, option) {
            const current = this.getAnswer(questionId) || [];
            const arr = Array.isArray(current) ? current : [];
            const next = arr.includes(option)
                ? arr.filter(v => v !== option)
                : [...arr, option];
            this.setAnswer(questionId, next);
        },

        /* ── Repeater ── */
        getRepeaterItems(questionId) {
            const val = this.getAnswer(questionId);
            return Array.isArray(val) ? val : [];
        },
        repeaterAdd(question) {
            const items = this.getRepeaterItems(question.id);
            const blank = {};
            question.fields.forEach(f => { blank[f.id] = f.type === 'multiple' ? [] : ''; });
            this.setAnswer(question.id, [...items, blank]);
        },
        repeaterRemove(questionId, idx) {
            const items = this.getRepeaterItems(questionId);
            this.setAnswer(questionId, items.filter((_, i) => i !== idx));
        },
        repeaterSet(questionId, idx, fieldId, value) {
            const items = [...this.getRepeaterItems(questionId)];
            items[idx] = { ...items[idx], [fieldId]: value };
            this.setAnswer(questionId, items);
        },
        repeaterToggleMultiple(questionId, idx, fieldId, option) {
            const items = [...this.getRepeaterItems(questionId)];
            const current = Array.isArray(items[idx][fieldId]) ? items[idx][fieldId] : [];
            items[idx] = {
                ...items[idx],
                [fieldId]: current.includes(option)
                    ? current.filter(v => v !== option)
                    : [...current, option],
            };
            this.setAnswer(questionId, items);
        },

        /* ── Respuestas locales (editables antes de guardar) ── */
        answerKey(questionId, entryId) {
            return entryId ? `${entryId}__${questionId}` : questionId;
        },
        initLocalAnswers() {
            const map = {};
            if (!this.currentSectionRender) return;
            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question') {
                    if (item.question.questionTypeId === 'info') continue;
                    const savedVal = item.answer?.value ?? '';
                    map[item.question.id] = this._sanitizeOptionAnswer(item.question, savedVal, this.formState, undefined);
                } else {
                    for (const e of item.entries) {
                        const entryIndex = e.entry.iteration - 1;
                        for (const q of e.questions) {
                            if (q.question.questionTypeId === 'info') continue;
                            const savedVal = q.answer?.value ?? '';
                            map[this.answerKey(q.question.id, e.entry.id)] = this._sanitizeOptionAnswer(q.question, savedVal, this.formState, entryIndex);
                        }
                    }
                }
            }
            this.localAnswers = map;
        },

        // Limpia un valor guardado que ya no existe en las opciones actuales de la pregunta.
        // Para single/dropdown: devuelve '' si el valor no está en opciones.
        // Para multiple: filtra del CSV solo los valores que siguen existiendo.
        // Si la pregunta no tiene stateOptionsPath devuelve el valor sin cambios.
        _sanitizeOptionAnswer(question, value, formState, entryIndex) {
            if (!question.stateOptionsPath || !value) return value;
            const opts  = resolveQuestionOptions(question, formState, entryIndex);
            const valid = new Set(opts.map(o => o.value));
            if (question.questionTypeId === 'multiple') {
                const kept = value.split(',').map(v => v.trim()).filter(v => v && valid.has(v));
                return kept.join(',');
            }
            return valid.has(value) ? value : '';
        },
        getLocalAnswer(questionId, entryId = null) {
            return this.localAnswers[this.answerKey(questionId, entryId)] ?? '';
        },
        setLocalAnswer(questionId, value, entryId = null) {
            this.localAnswers = { ...this.localAnswers, [this.answerKey(questionId, entryId)]: value };
        },
        isMultipleSelected(questionId, optionValue, entryId = null) {
            const val = this.getLocalAnswer(questionId, entryId);
            return val ? val.split(',').map(v => v.trim()).includes(optionValue) : false;
        },
        toggleMultipleAnswer(questionId, optionValue, entryId = null) {
            const current = this.getLocalAnswer(questionId, entryId);
            const arr = current ? current.split(',').map(v => v.trim()).filter(Boolean) : [];
            const idx = arr.indexOf(optionValue);
            if (idx === -1) arr.push(optionValue); else arr.splice(idx, 1);
            this.setLocalAnswer(questionId, arr.join(','), entryId);
        },

        /* ── Responder pregunta: valida + recalcula visibilidad ── */
        onAnswer(questionId, value, entryId = null) {
            // 1. Guardar respuesta
            this.setLocalAnswer(questionId, value, entryId);

            // 2. Validar (v1: solo requerido)
            const key      = this.answerKey(questionId, entryId);
            const question = this._findQuestionInRender(questionId, entryId);
            const errors   = { ...this.questionErrors };
            if (question && question.required && value === '') {
                errors[key] = 'Este campo es requerido';
            } else {
                delete errors[key];
            }
            this.questionErrors = errors;

            // 3. Re-evaluar visibilidad con las respuestas actuales
            this._reevaluateVisibility(this._buildTempSubmission());

            // 4. Limpiar respuestas de items ocultos y re-evaluar en cascada
            if (this._clearHiddenAnswers()) {
                this._reevaluateVisibility(this._buildTempSubmission());
            }

            // 5. Emitir respuestas actuales al padre
            console.log('[onAnswer] emitiendo answers-updated | questionId:', questionId, '| value:', value, '| entryId:', entryId);
            this._emitAnswersUpdated();
        },

        onToggleMultipleAnswer(questionId, optionValue, entryId = null) {
            this.toggleMultipleAnswer(questionId, optionValue, entryId);
            this.onAnswer(questionId, this.getLocalAnswer(questionId, entryId), entryId);
        },

        onBlur(questionId, entryId = null) {
            console.log('[onBlur] questionId:', questionId, '| entryId:', entryId);
            const question = this._findQuestionInRender(questionId, entryId);
            console.log('[onBlur] question encontrada:', question);
            if (!question || !question.required) {
                console.log('[onBlur] saliendo — no required o no encontrada');
                return;
            }
            const value = this.getLocalAnswer(questionId, entryId);
            console.log('[onBlur] value:', JSON.stringify(value));
            if (value === '') {
                const key = this.answerKey(questionId, entryId);
                console.log('[onBlur] seteando error en key:', key);
                this.questionErrors = { ...this.questionErrors, [key]: 'Este campo es requerido' };
            }
        },

        // Construye un tempSubmission combinando base (sin sección actual) + respuestas actuales
        _buildTempSubmission() {
            const base  = this.submissionWithoutCurrentSection() ?? { directAnswers: [], repeaterEntries: [] };
            const fresh = this.collectSectionAnswers();
            return {
                ...base,
                directAnswers:   [...(base.directAnswers   || []), ...fresh.directAnswers],
                repeaterEntries: [...(base.repeaterEntries || []), ...fresh.repeaterEntries],
            };
        },

        // Emite 'answers-updated' con la sección actual y las respuestas completas del submission
        _emitAnswersUpdated() {
            if (!this.currentSection) return;
            const payload = {
                currentSectionId: this.currentSection.id,
                answers:          this._buildTempSubmission(),
            };
            console.log('[_emitAnswersUpdated] emitiendo answers-updated | payload:', payload);
            this.$emit('answers-updated', payload);
        },

        // Re-evalúa isVisible de todos los items de la sección y de secciones posteriores
        _reevaluateVisibility(tempSubmission, formState) {
            const fs  = this.formStructure;
            const ord = this.currentSection.order;
            const fst = formState !== undefined ? formState : this.formState;

            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question') {
                    item.isVisible = checkVisibility(fs, tempSubmission, null,
                        directQuestionItem(item.question, ord), fst).visible;
                } else if (item.type === 'repeater') {
                    item.isVisible = checkVisibility(fs, tempSubmission, null,
                        repeaterItem(item.repeater, ord), fst).visible;

                    // ── Re-sincronizar entries si el repeater está controlado por estado ──
                    if (item.repeater.stateItems) {
                        const stateArray = getObjectProperty(fst || {}, item.repeater.stateItems);
                        const newCount     = Array.isArray(stateArray) ? stateArray.length : 0;
                        const currentCount = item.entries.length;

                        if (newCount !== currentCount) {
                            // Índice de entries actuales por iteración para preservar respuestas
                            const byIteration = {};
                            item.entries.forEach(ed => { byIteration[ed.entry.iteration] = ed; });

                            const newLocalAnswers = { ...this.localAnswers };

                            // Eliminar claves de entries que ya no existen
                            for (let i = newCount; i < currentCount; i++) {
                                const ed = byIteration[i + 1];
                                if (ed) {
                                    for (const qData of ed.questions) {
                                        delete newLocalAnswers[this.answerKey(qData.question.id, ed.entry.id)];
                                    }
                                }
                            }

                            // Construir nuevo array de entries
                            const newEntries = [];
                            for (let i = 0; i < newCount; i++) {
                                const iteration = i + 1;
                                const existing  = byIteration[iteration];
                                if (existing) {
                                    newEntries.push(existing);
                                } else {
                                    const newEntry = {
                                        id:              `temp-state-${item.repeater.id}-${iteration}`,
                                        repeaterGroupId: item.repeater.id,
                                        iteration,
                                        isTemp:          true,
                                        answers:         [],
                                    };
                                    const questions = (item.repeater.questions || []).map(q => ({
                                        question:  q,
                                        answer:    null,
                                        isVisible: true,
                                    }));
                                    newEntries.push({ entry: newEntry, questions });
                                    for (const q of (item.repeater.questions || [])) {
                                        if (q.questionTypeId !== 'info') {
                                            newLocalAnswers[this.answerKey(q.id, newEntry.id)] = '';
                                        }
                                    }
                                }
                            }

                            item.entries       = newEntries;
                            this.localAnswers  = newLocalAnswers;
                        }
                    }

                    for (const entryData of item.entries) {
                        const entryAnswers = (tempSubmission.repeaterEntries.find(e => e.id === entryData.entry.id) || {}).answers || [];
                        for (const qData of entryData.questions) {
                            qData.isVisible = checkVisibility(fs, tempSubmission, entryAnswers,
                                repeaterQuestionItem(qData.question, ord, item.repeater.order, item.repeater.id), fst).visible;
                        }
                    }
                }
            }

            for (const sec of fs.sections) {
                if (sec.order > ord) {
                    sec.isVisible = checkVisibility(fs, tempSubmission, null, sectionItem(sec), fst).visible;
                }
            }

        },

        // Limpia respuestas huérfanas de preguntas con stateOptionsPath.
        // Solo debe llamarse cuando cambia formState (E-16), nunca desde onAnswer.
        _clearStaleOptionAnswers(formState) {
            const fst        = formState !== undefined ? formState : this.formState;
            const newAnswers = { ...this.localAnswers };
            let anyStale     = false;
            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question') {
                    if (!item.question.stateOptionsPath) continue;
                    const key       = this.answerKey(item.question.id);
                    const sanitized = this._sanitizeOptionAnswer(item.question, newAnswers[key] || '', fst, undefined);
                    if (sanitized !== (newAnswers[key] || '')) { newAnswers[key] = sanitized; anyStale = true; }
                } else if (item.type === 'repeater') {
                    for (const entryData of item.entries) {
                        const entryIndex = entryData.entry.iteration - 1;
                        for (const qData of entryData.questions) {
                            if (!qData.question.stateOptionsPath) continue;
                            const key       = this.answerKey(qData.question.id, entryData.entry.id);
                            const sanitized = this._sanitizeOptionAnswer(qData.question, newAnswers[key] || '', fst, entryIndex);
                            if (sanitized !== (newAnswers[key] || '')) { newAnswers[key] = sanitized; anyStale = true; }
                        }
                    }
                }
            }
            if (anyStale) this.localAnswers = newAnswers;
        },

        // Vacía localAnswers de los items ocultos. Retorna true si limpió algo (para cascada).
        _clearHiddenAnswers() {
            let anyCleared = false;
            const newAnswers = { ...this.localAnswers };

            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question' && !item.isVisible) {
                    const k = this.answerKey(item.question.id);
                    if (newAnswers[k]) { newAnswers[k] = ''; anyCleared = true; }
                } else if (item.type === 'repeater') {
                    for (const entryData of item.entries) {
                        for (const qData of entryData.questions) {
                            if (!qData.isVisible) {
                                const k = this.answerKey(qData.question.id, entryData.entry.id);
                                if (newAnswers[k]) { newAnswers[k] = ''; anyCleared = true; }
                            }
                        }
                    }
                }
            }

            if (anyCleared) this.localAnswers = newAnswers;
            return anyCleared;
        },

        // Helper: encuentra la question en currentSectionRender por id (y entryId para repeaters)
        _findQuestionInRender(questionId, entryId = null) {
            if (!this.currentSectionRender) return null;
            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question' && item.question.id === questionId && !entryId) {
                    return item.question;
                }
                if (item.type === 'repeater') {
                    for (const entryData of item.entries) {
                        if (entryData.entry.id !== entryId) continue;
                        const qData = entryData.questions.find(q => q.question.id === questionId);
                        if (qData) return qData.question;
                    }
                }
            }
            return null;
        },

        /* ── Colectar respuestas de la sección actual ── */
        collectSectionAnswers() {
            const directAnswers   = [];
            const repeaterEntries = [];

            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question') {
                    // info banners no generan respuesta
                    if (item.question.questionTypeId === 'info') continue;
                    const value = this.getLocalAnswer(item.question.id);
                    if (value !== '') directAnswers.push({ questionId: item.question.id, value });

                } else if (item.type === 'repeater') {
                    for (const entryData of item.entries) {
                        const answers = [];
                        for (const qData of entryData.questions) {
                            const value = this.getLocalAnswer(qData.question.id, entryData.entry.id);
                            if (value !== '') answers.push({ questionId: qData.question.id, value });
                        }
                        repeaterEntries.push({
                            id:              entryData.entry.id,
                            repeaterGroupId: item.repeater.id,
                            iteration:       entryData.entry.iteration,
                            isTemp:          entryData.entry.isTemp || false,
                            answers,
                        });
                    }
                }
            }

            return { directAnswers, repeaterEntries };
        },

        /* ── Submission sin respuestas de la sección actual ── */
        submissionWithoutCurrentSection() {
            if (!this.formSubmission) return null;

            const questionIds     = new Set();
            const repeaterGroupIds = new Set();

            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question')  questionIds.add(item.question.id);
                if (item.type === 'repeater')  repeaterGroupIds.add(item.repeater.id);
            }

            return {
                ...this.formSubmission,
                directAnswers:   (this.formSubmission.directAnswers   || []).filter(a => !questionIds.has(a.questionId)),
                repeaterEntries: (this.formSubmission.repeaterEntries || []).filter(e => !repeaterGroupIds.has(e.repeaterGroupId)),
            };
        },

        /* ── Repeater: agregar / eliminar entries ── */
        addRepeaterEntry(item) {
            const tempId    = 'temp-' + Date.now() + '-' + Math.random().toString(36).slice(2, 8);
            const iteration = item.entries.length + 1;
            const newEntry  = {
                id:              tempId,
                repeaterGroupId: item.repeater.id,
                iteration,
                isTemp:          true,
            };

            const sectionOrder = this.currentSection.order;
            const questions = (item.repeater.questions || []).map(q => {
                const qVisItem = repeaterQuestionItem(q, sectionOrder, item.repeater.order, item.repeater.id);
                const qVis     = checkVisibility(this.formStructure, this.formSubmission, [], qVisItem, this.formState);
                return { question: q, answer: null, isVisible: qVis.visible };
            });

            item.entries.push({ entry: newEntry, questions });

            const newAnswers = { ...this.localAnswers };
            for (const q of (item.repeater.questions || [])) {
                // info banners no tienen respuesta — no crear clave en localAnswers
                if (q.questionTypeId === 'info') continue;
                newAnswers[this.answerKey(q.id, tempId)] = '';
            }
            this.localAnswers = newAnswers;

            const min = item.repeater.minRepetitions || 0;
            if (min > 0 && item.entries.length >= min) {
                const { [item.repeater.id]: _, ...rest } = this.repeaterErrors;
                this.repeaterErrors = rest;
            }

            this._emitAnswersUpdated();
        },

        removeRepeaterEntry(item, entryData) {
            const idx = item.entries.indexOf(entryData);
            if (idx === -1) return;

            item.entries.splice(idx, 1);
            item.entries.forEach((e, i) => { e.entry.iteration = i + 1; });

            const prefix  = entryData.entry.id + '__';
            const cleaned = {};
            for (const [k, v] of Object.entries(this.localAnswers)) {
                if (!k.startsWith(prefix)) cleaned[k] = v;
            }
            this.localAnswers = cleaned;

            const min = item.repeater.minRepetitions || 0;
            if (min > 0 && item.entries.length < min) {
                const itemName = item.repeater.itemName || item.repeater.name;
                this.repeaterErrors = {
                    ...this.repeaterErrors,
                    [item.repeater.id]: `Se deben agregar mínimo ${min} ${itemName}`,
                };
            }

            this._emitAnswersUpdated();
        },

        /* ── Validar sección completa antes de guardar ── */
        validateCurrentSection() {
            const errors         = {};
            const repeaterErrors = {};
            let valid = true;

            for (const item of this.currentSectionRender.formItems) {
                if (item.type === 'question' && item.isVisible && item.question.required && item.question.questionTypeId !== 'info') {
                    if (this.getLocalAnswer(item.question.id) === '') {
                        errors[this.answerKey(item.question.id)] = 'Este campo es requerido';
                        valid = false;
                    }
                }

                if (item.type === 'question' && item.isVisible && item.question.id === Q_EQUIPOS_SALVIA) {
                    const selected = this.getLocalAnswer(item.question.id).split(',').map(v => v.trim());
                    if (selected.includes(OPT_MEDIDAS_EMERGENCIA) && selected.includes(OPT_ATENCION_PSICO)) {
                        errors[this.answerKey(item.question.id)] = 'No se pueden seleccionar Medidas de Emergencia y Apoyo Psicosocial al mismo tiempo';
                        valid = false;
                    }
                }

                if (item.type === 'repeater' && item.isVisible) {
                    const min = item.repeater.minRepetitions || 0;
                    if (min > 0 && item.entries.length < min) {
                        const baseName = item.repeater.itemName || item.repeater.name;
                        const itemName = min > 1 ? baseName + 's' : baseName;
                        repeaterErrors[item.repeater.id] = `Se deben agregar mínimo ${min} ${itemName}`;
                        valid = false;
                    }
                    for (const entryData of item.entries) {
                        for (const qData of entryData.questions) {
                            if (qData.isVisible && qData.question.required && qData.question.questionTypeId !== 'info') {
                                if (this.getLocalAnswer(qData.question.id, entryData.entry.id) === '') {
                                    errors[this.answerKey(qData.question.id, entryData.entry.id)] = 'Este campo es requerido';
                                    valid = false;
                                }
                            }
                        }
                    }
                }
            }

            this.questionErrors  = { ...this.questionErrors,  ...errors };
            this.repeaterErrors  = { ...this.repeaterErrors,  ...repeaterErrors };
            return valid;
        },

        /* ── Guardar sección y avanzar ── */
        async saveSection() {
            if (!this.validateCurrentSection()) {
                this.showValidationError = true;
                return;
            }

            this.saving = true;
            const fresh = this.collectSectionAnswers();
            const body  = {
                formId:           this.formId,
                formSectionId:    this.currentSection.id,
                formSubmissionId: this.formSubmission?.id ?? '',
                directAnswers:    fresh.directAnswers,
                repeaterEntries:  fresh.repeaterEntries,
                formState:        this.formState,
            };

            console.log('[saveSection] enviando sección:', this.currentSection.name);
            console.log('[saveSection] body:', body);

            try {
                const res = await fetch('/api/v1/forms/saveSection', {
                    method:  'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body:    JSON.stringify(body),
                });

                console.log('[saveSection] response status:', res.status);
                if (!res.ok) throw new Error('Hubo un error al guardar la información, por favor verifique su conexión a internet y vuelva a intentarlo');

                const data = await res.json();
                console.log('[saveSection] data recibida:', data);

                // Actualizar estado con la respuesta del servidor
                this.formStructure  = data.formStructure;
                this.formSubmission = data.formSubmission ?? null;
                console.log('[saveSection] formSubmission actualizado:', this.formSubmission?.id);
                this._emitAnswersUpdated();

                // Detectar formulario completado (todas las secciones visibles respondidas)
                const allAnswered = this.formStructure.sections
                    .filter(s => s.isVisible)
                    .every(s => s.isAnswered);
                if (allAnswered) {
                    console.log('[saveSection] formulario completado — emitiendo form-completed');
                    this.$emit('form-completed');
                    return;
                }

                // Navegación: primera opción → siguiente sección visible
                // Fallback → currentSection del servidor (si la siguiente no está enabled)
                const savedId     = this.currentSection.id;
                const newVisible  = this.formStructure.sections.filter(s => s.isVisible);
                const savedIdx    = newVisible.findIndex(s => s.id === savedId);
                const nextSection = savedIdx >= 0 && savedIdx < newVisible.length - 1
                    ? newVisible[savedIdx + 1]
                    : null;

                console.log('[saveSection] nextSection:', nextSection?.name ?? 'ninguna');
                console.log('[saveSection] firstUnanswered:', this.firstUnansweredSection?.name ?? 'ninguna');

                if (nextSection && (nextSection.isAnswered || this.firstUnansweredSection?.id === nextSection.id)) {
                    console.log('[saveSection] navegando a siguiente sección visible:', nextSection.name);
                    this.goToSectionById(nextSection);
                } else if (data.currentSection) {
                    console.log('[saveSection] fallback → currentSection del server:', data.currentSection.name);
                    this.goToSectionById(data.currentSection);
                }

            } catch (e) {
                console.error('[saveSection] error:', e);
                this.error = e.message;
            } finally {
                this.saving = false;
            }
        },

        /* ── Guardar borrador: omite validación, no navega, no dispara form-completed ── */
        async saveDraft() {
            this.saving = true;
            const fresh = this.collectSectionAnswers();
            const body  = {
                formId:           this.formId,
                formSectionId:    this.currentSection.id,
                formSubmissionId: this.formSubmission?.id ?? '',
                directAnswers:    fresh.directAnswers,
                repeaterEntries:  fresh.repeaterEntries,
                formState:        this.formState,
            };

            console.log('[saveDraft] enviando borrador:', this.currentSection.name);
            console.log('[saveDraft] body:', body);

            try {
                const res = await fetch('/api/v1/forms/saveSection', {
                    method:  'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body:    JSON.stringify(body),
                });

                console.log('[saveDraft] response status:', res.status);
                if (!res.ok) throw new Error('Hubo un error al guardar el borrador');

                const data = await res.json();
                console.log('[saveDraft] data recibida:', data);

                // Actualizar estado con la respuesta del servidor — sin navegar ni
                // evaluar form-completed: un borrador no significa que el usuario
                // haya terminado la sección.
                this.formStructure  = data.formStructure;
                this.formSubmission = data.formSubmission ?? null;
                this._emitAnswersUpdated();

                this.draftSaved = true;
                setTimeout(() => { this.draftSaved = false; }, 2000);
            } catch (e) {
                console.error('[saveDraft] error:', e);
                this.draftError = true;
                setTimeout(() => { this.draftError = false; }, 2000);
            } finally {
                this.saving = false;
            }
        },

        /* ── Navegación entre secciones ── */
        goNext() {
            const idx = this.visibleSections.findIndex(s => s.id === this.currentSection.id);
            if (idx < this.visibleSections.length - 1) this.goToSectionById(this.visibleSections[idx + 1]);
        },
        goPrev() {
            const idx = this.visibleSections.findIndex(s => s.id === this.currentSection.id);
            if (idx > 0) this.goToSectionById(this.visibleSections[idx - 1]);
        },
        goToSection(index) { if (index <= this.currentIndex) this.currentIndex = index; },
    },
    template: `
<div class="df-wrapper">

    <!-- Guardar borrador -->
    <button
        v-if="canEdit && hasDraftAnswers"
        type="button"
        class="df-draft-save-btn"
        :class="{ 'df-draft-saved': draftSaved, 'df-draft-error': draftError }"
        :disabled="saving"
        @click="saveDraft"
    >
        <span v-if="draftSaved">✓ Guardado</span>
        <span v-else-if="draftError">⚠ Error al guardar</span>
        <span v-else>💾 Guardar borrador</span>
    </button>

    <!-- Loader -->
    <div v-if="loading" style="display:flex;align-items:center;justify-content:center;width:100%;padding:48px 0;gap:12px;color:#6b7280;font-size:14px;">
        <svg style="width:20px;height:20px;animation:df-spin 0.8s linear infinite;flex-shrink:0" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
            <path d="M12 2a10 10 0 0 1 10 10" stroke="#7c3aed" stroke-width="3" stroke-linecap="round"/>
        </svg>
        Cargando formulario...
    </div>

    <!-- Error -->
    <div v-else-if="error" style="width:100%;padding:24px;background:#fef2f2;border:1px solid #fecaca;border-radius:8px;color:#b91c1c;font-size:14px;">
        \${ error }
    </div>

    <!-- Formulario -->
    <template v-else>

    <!-- Sidebar -->
    <aside class="df-sidebar">
        <div class="df-sidebar-header">
            <p class="df-sidebar-title">Secciones</p>
        </div>
        <div class="df-sidebar-nav">
            <button
                v-for="section in visibleSections"
                :key="section.id"
                type="button"
                class="df-section-item"
                :class="sidebarSectionClass(section)"
                :disabled="!isSectionEnabled(section)"
                @click="goToSectionById(section)"
            >
                <span class="df-badge" :class="sidebarBadgeClass(section)">
                    <span v-if="section.isAnswered">✓</span>
                    <span v-else>\${ section.order }</span>
                </span>
                <span>\${ applyRenderModifications(section.name, section.modifications, 'name', formState) }</span>
            </button>
        </div>
        <div class="df-progress-wrap">
            <div class="df-progress-label">
                <span>Progreso</span>
                <span>\${ answeredCount } / \${ visibleSections.length }</span>
            </div>
            <div class="df-progress-bar-bg">
                <div class="df-progress-bar-fill" :style="{ width: progressPercent + '%' }"></div>
            </div>
        </div>
    </aside>

    <!-- Main panel -->
    <div class="df-main">

        <!-- Saving overlay -->
        <div v-if="saving" class="df-saving-overlay">
            <svg style="width:18px;height:18px;animation:df-spin 0.8s linear infinite;flex-shrink:0" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
                <path d="M12 2a10 10 0 0 1 10 10" stroke="#7c3aed" stroke-width="3" stroke-linecap="round"/>
            </svg>
            Guardando...
        </div>

        <!-- Encabezado de sección -->
        <div class="df-section-header">
            <div class="df-section-header-row">
                <div class="df-section-number">\${ currentSection.order }</div>
                <h2 class="df-section-title">\${ applyRenderModifications(currentSection.name, currentSection.modifications, 'name', formState) }</h2>
                <span v-if="!canEdit" class="df-readonly-badge">🔒 Solo lectura</span>
            </div>
            <p v-if="currentSection.description" class="df-section-desc">\${ applyRenderModifications(currentSection.description, currentSection.modifications, 'description', formState) }</p>
        </div>

        <!-- Items de la sección -->
        <div class="df-questions" v-if="currentSectionRender">
            <template v-for="item in currentSectionRender.formItems" :key="item.type === 'question' ? item.question.id : item.repeater.id">

                <!-- ── Info banner ── -->
                <div v-if="item.type === 'question' && item.isVisible && item.question.questionTypeId === 'info'" class="df-info-banner">
                    <span class="df-info-banner-icon">ℹ️</span>
                    <span>\${ applyRenderModifications(item.question.description, item.question.modifications, 'description', formState) }</span>
                </div>

                <!-- ── Pregunta directa ── -->
                <div v-else-if="item.type === 'question' && item.isVisible" class="df-question">
                    <label>
                        \${ applyRenderModifications(item.question.description, item.question.modifications, 'description', formState) }
                        <span v-if="item.question.required" class="df-required">*</span>
                    </label>

                    <!-- single -->
                    <div v-if="item.question.questionTypeId === 'single'" class="df-btn-group">
                        <button v-for="opt in resolveQuestionOptions(item.question, formState)" :key="opt.id"
                            type="button" class="df-opt-btn"
                            :class="{ 'df-selected': getLocalAnswer(item.question.id) === opt.value }"
                            :disabled="!canEdit"
                            @click="onAnswer(item.question.id, opt.value)"
                        >\${ opt.label }</button>
                    </div>

                    <!-- dropdown -->
                    <select v-else-if="item.question.questionTypeId === 'dropdown'"
                        class="df-select"
                        :disabled="!canEdit"
                        :value="getLocalAnswer(item.question.id)"
                        @change="onAnswer(item.question.id, $event.target.value)"
                        @blur="onBlur(item.question.id)"
                    >
                        <option value="">Selecciona una opción...</option>
                        <option v-for="opt in resolveQuestionOptions(item.question, formState)" :key="opt.id" :value="opt.value">\${ opt.label }</option>
                    </select>

                    <!-- boolean -->
                    <div v-else-if="item.question.questionTypeId === 'boolean'" class="df-bool-group">
                        <button type="button" class="df-bool-btn"
                            :class="{ 'df-selected': getLocalAnswer(item.question.id) === 'true' }"
                            :disabled="!canEdit"
                            @click="onAnswer(item.question.id, 'true')">Sí</button>
                        <button type="button" class="df-bool-btn"
                            :class="{ 'df-selected': getLocalAnswer(item.question.id) === 'false' }"
                            :disabled="!canEdit"
                            @click="onAnswer(item.question.id, 'false')">No</button>
                    </div>

                    <!-- multiple -->
                    <div v-else-if="item.question.questionTypeId === 'multiple'" class="df-btn-group">
                        <button v-for="opt in resolveQuestionOptions(item.question, formState)" :key="opt.id"
                            type="button" class="df-opt-btn"
                            :class="{ 'df-selected': isMultipleSelected(item.question.id, opt.value) }"
                            :disabled="!canEdit"
                            @click="onToggleMultipleAnswer(item.question.id, opt.value)"
                        >
                            <span v-if="isMultipleSelected(item.question.id, opt.value)">✓ </span>\${ opt.label }
                        </button>
                    </div>

                    <!-- text -->
                    <textarea v-else-if="item.question.questionTypeId === 'text'"
                        class="df-textarea"
                        :disabled="!canEdit"
                        :value="getLocalAnswer(item.question.id)"
                        @input="onAnswer(item.question.id, $event.target.value)"
                        @blur="onBlur(item.question.id)"
                    ></textarea>

                    <!-- number -->
                    <input v-else-if="item.question.questionTypeId === 'number'"
                        class="df-input" type="number"
                        :disabled="!canEdit"
                        :value="getLocalAnswer(item.question.id)"
                        @input="onAnswer(item.question.id, $event.target.value)"
                    />

                    <!-- date -->
                    <input v-else-if="item.question.questionTypeId === 'date'"
                        class="df-input" type="date"
                        :disabled="!canEdit"
                        :value="getLocalAnswer(item.question.id)"
                        @input="onAnswer(item.question.id, $event.target.value)"
                        @blur="onBlur(item.question.id)"
                    />

                    <!-- time -->
                    <input v-else-if="item.question.questionTypeId === 'time'"
                        class="df-input" type="time"
                        :disabled="!canEdit"
                        :value="getLocalAnswer(item.question.id)"
                        @input="onAnswer(item.question.id, $event.target.value)"
                        @blur="onBlur(item.question.id)"
                    />

                    <!-- datetime -->
                    <input v-else-if="item.question.questionTypeId === 'datetime'"
                        class="df-input" type="datetime-local"
                        :disabled="!canEdit"
                        :value="getLocalAnswer(item.question.id)"
                        @input="onAnswer(item.question.id, $event.target.value)"
                        @blur="onBlur(item.question.id)"
                    />
                    <span v-if="questionErrors[item.question.id]" class="df-error">\${ questionErrors[item.question.id] }</span>
                </div>

                <!-- ── Repeater group ── -->
                <div v-else-if="item.type === 'repeater' && item.isVisible" class="df-question">
                    <label>
                        \${ applyRenderModifications(item.repeater.name, item.repeater.modifications, 'name', formState) }
                        <span v-if="item.repeater.minRepetitions > 0" class="df-required">*</span>
                    </label>
                    <div class="df-repeater">

                        <!-- entries -->
                        <div v-for="entryData in item.entries" :key="entryData.entry.id" class="df-repeater-item">
                            <div class="df-repeater-item-header">
                                <span class="df-repeater-item-title">\${ applyRenderModifications(item.repeater.itemName || item.repeater.name, item.repeater.modifications, 'item_name', formState) } #\${ entryData.entry.iteration }</span>
                                <button v-if="canEdit && !item.repeater.stateItems" type="button" class="df-repeater-delete" @click="removeRepeaterEntry(item, entryData)">✕ Eliminar</button>
                            </div>

                            <template v-for="qData in entryData.questions" :key="qData.question.id">
                                <!-- info banner dentro de repeater -->
                                <div v-if="qData.isVisible && qData.question.questionTypeId === 'info'" class="df-info-banner">
                                    <span class="df-info-banner-icon">ℹ️</span>
                                    <span>\${ applyRenderModifications(qData.question.description, qData.question.modifications, 'description', formState, entryData.entry.iteration - 1) }</span>
                                </div>

                                <div v-else-if="qData.isVisible" class="df-question">
                                    <label>
                                        \${ applyRenderModifications(qData.question.description, qData.question.modifications, 'description', formState, entryData.entry.iteration - 1) }
                                        <span v-if="qData.question.required" class="df-required">*</span>
                                    </label>

                                    <!-- single -->
                                    <div v-if="qData.question.questionTypeId === 'single'" class="df-btn-group">
                                        <button v-for="opt in resolveQuestionOptions(qData.question, formState, entryData.entry.iteration - 1)" :key="opt.id"
                                            type="button" class="df-opt-btn"
                                            :class="{ 'df-selected': getLocalAnswer(qData.question.id, entryData.entry.id) === opt.value }"
                                            :disabled="!canEdit"
                                            @click="onAnswer(qData.question.id, opt.value, entryData.entry.id)"
                                        >\${ opt.label }</button>
                                    </div>

                                    <!-- dropdown -->
                                    <select v-else-if="qData.question.questionTypeId === 'dropdown'"
                                        class="df-select"
                                        :disabled="!canEdit"
                                        :value="getLocalAnswer(qData.question.id, entryData.entry.id)"
                                        @change="onAnswer(qData.question.id, $event.target.value, entryData.entry.id)"
                                        @blur="onBlur(qData.question.id, entryData.entry.id)"
                                    >
                                        <option value="">Selecciona una opción...</option>
                                        <option v-for="opt in resolveQuestionOptions(qData.question, formState, entryData.entry.iteration - 1)" :key="opt.id" :value="opt.value">\${ opt.label }</option>
                                    </select>

                                    <!-- boolean -->
                                    <div v-else-if="qData.question.questionTypeId === 'boolean'" class="df-bool-group">
                                        <button type="button" class="df-bool-btn"
                                            :class="{ 'df-selected': getLocalAnswer(qData.question.id, entryData.entry.id) === 'true' }"
                                            :disabled="!canEdit"
                                            @click="onAnswer(qData.question.id, 'true', entryData.entry.id)">Sí</button>
                                        <button type="button" class="df-bool-btn"
                                            :class="{ 'df-selected': getLocalAnswer(qData.question.id, entryData.entry.id) === 'false' }"
                                            :disabled="!canEdit"
                                            @click="onAnswer(qData.question.id, 'false', entryData.entry.id)">No</button>
                                    </div>

                                    <!-- multiple -->
                                    <div v-else-if="qData.question.questionTypeId === 'multiple'" class="df-btn-group">
                                        <button v-for="opt in resolveQuestionOptions(qData.question, formState, entryData.entry.iteration - 1)" :key="opt.id"
                                            type="button" class="df-opt-btn"
                                            :class="{ 'df-selected': isMultipleSelected(qData.question.id, opt.value, entryData.entry.id) }"
                                            :disabled="!canEdit"
                                            @click="onToggleMultipleAnswer(qData.question.id, opt.value, entryData.entry.id)"
                                        >
                                            <span v-if="isMultipleSelected(qData.question.id, opt.value, entryData.entry.id)">✓ </span>\${ opt.label }
                                        </button>
                                    </div>

                                    <!-- text -->
                                    <textarea v-else-if="qData.question.questionTypeId === 'text'"
                                        class="df-textarea"
                                        :disabled="!canEdit"
                                        :value="getLocalAnswer(qData.question.id, entryData.entry.id)"
                                        @input="onAnswer(qData.question.id, $event.target.value, entryData.entry.id)"
                                        @blur="onBlur(qData.question.id, entryData.entry.id)"
                                    ></textarea>

                                    <!-- number -->
                                    <input v-else-if="qData.question.questionTypeId === 'number'"
                                        class="df-input" type="number"
                                        :disabled="!canEdit"
                                        :value="getLocalAnswer(qData.question.id, entryData.entry.id)"
                                        @input="onAnswer(qData.question.id, $event.target.value, entryData.entry.id)"
                                    />

                                    <!-- date -->
                                    <input v-else-if="qData.question.questionTypeId === 'date'"
                                        class="df-input" type="date"
                                        :disabled="!canEdit"
                                        :value="getLocalAnswer(qData.question.id, entryData.entry.id)"
                                        @input="onAnswer(qData.question.id, $event.target.value, entryData.entry.id)"
                                        @blur="onBlur(qData.question.id, entryData.entry.id)"
                                    />

                                    <!-- time -->
                                    <input v-else-if="qData.question.questionTypeId === 'time'"
                                        class="df-input" type="time"
                                        :disabled="!canEdit"
                                        :value="getLocalAnswer(qData.question.id, entryData.entry.id)"
                                        @input="onAnswer(qData.question.id, $event.target.value, entryData.entry.id)"
                                        @blur="onBlur(qData.question.id, entryData.entry.id)"
                                    />

                                    <!-- datetime -->
                                    <input v-else-if="qData.question.questionTypeId === 'datetime'"
                                        class="df-input" type="datetime-local"
                                        :disabled="!canEdit"
                                        :value="getLocalAnswer(qData.question.id, entryData.entry.id)"
                                        @input="onAnswer(qData.question.id, $event.target.value, entryData.entry.id)"
                                        @blur="onBlur(qData.question.id, entryData.entry.id)"
                                    />
                                    <span v-if="questionErrors[answerKey(qData.question.id, entryData.entry.id)]" class="df-error">\${ questionErrors[answerKey(qData.question.id, entryData.entry.id)] }</span>
                                </div>
                            </template>
                        </div>

                        <button v-if="canEdit && !item.repeater.stateItems" type="button" class="df-repeater-add" @click="addRepeaterEntry(item)">+ Agregar \${ applyRenderModifications(item.repeater.itemName || item.repeater.name, item.repeater.modifications, 'item_name', formState) }</button>
                    </div>
                    <span v-if="repeaterErrors[item.repeater.id]" class="df-error">\${ repeaterErrors[item.repeater.id] }</span>
                </div>

            </template>
        </div>

        <!-- Navegación -->
        <div class="df-nav">
            <button type="button" class="df-btn-prev" :disabled="isFirstSection" @click="goPrev">← Anterior</button>
            <button v-if="!isLastSection || canEdit" type="button" class="df-btn-next" @click="canEdit ? saveSection() : goNext()">
                <span v-if="isLastSection && canEdit">Guardar seguimiento ✓</span>
                <span v-else>Siguiente →</span>
            </button>
        </div>

    </div>

    </template><!-- v-else -->

    <!-- Popup validación -->
    <div v-if="showValidationError" class="df-popup-backdrop" @click.self="showValidationError = false">
        <div class="df-popup">
            <div class="df-popup-icon">⚠️</div>
            <p class="df-popup-title">Algunas respuestas no son válidas, revisa el formulario</p>
            <button type="button" class="df-popup-btn" @click="showValidationError = false">Entendido</button>
        </div>
    </div>

</div>
    `,
});
