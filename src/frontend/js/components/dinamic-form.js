/**
 * dinamic-form
 * Componente reutilizable para renderizar un formulario dinámico por secciones.
 * Soporta: single, multiple, boolean, dropdown, text, date, datetime, repeater.
 * Estilos basados en el prototipo Salvia.
 * Props:
 *   - formId (String): ID del formulario a cargar
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
            display: flex;
            gap: 24px;
            align-items: flex-start;
        }

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
        .df-section-item.active  { background: #f5f3ff; color: #6d28d9; font-weight: 600; cursor: default; }
        .df-section-item.completed { color: #4b5563; }
        .df-section-item.completed:hover { background: #f9fafb; }
        .df-section-item.pending { color: #d1d5db; cursor: not-allowed; }
        .df-badge {
            width: 20px; height: 20px;
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            font-size: 11px; font-weight: 700;
            flex-shrink: 0;
        }
        .df-badge.active    { background: #7c3aed; color: #fff; }
        .df-badge.completed { background: #22c55e; color: #fff; }
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

        /* ── Main panel ── */
        .df-main {
            flex: 1; min-width: 0;
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
        .df-question > label,
        .df-repeater-item label {
            display: block;
            font-size: 14px; font-weight: 500; color: #374151;
            margin-bottom: 4px; margin-top: 0;
            min-height: unset;
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
        .df-hint { font-size: 12px; color: #9ca3af; font-style: italic; }

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

/* ─── Componente Vue ─────────────────────────────────────────────────────── */
app.component('dinamic-form', {
    delimiters: ['${', '}'],
    props: {
        formId: { type: String, required: true },
    },
    data() {
        return {
            currentIndex: 0,
            answers: {},
            sections: DF_SCHEMA,
        };
    },
    computed: {
        currentSection() { return this.sections[this.currentIndex]; },
        totalSections()  { return this.sections.length; },
        progressPercent(){ return Math.round(((this.currentIndex + 1) / this.totalSections) * 100); },
        isFirst() { return this.currentIndex === 0; },
        isLast()  { return this.currentIndex === this.totalSections - 1; },
    },
    methods: {
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

        goNext() { if (!this.isLast) this.currentIndex++; },
        goPrev() { if (!this.isFirst) this.currentIndex--; },
        goToSection(index) { if (index <= this.currentIndex) this.currentIndex = index; },
    },
    template: `
<div class="df-wrapper">

    <!-- Sidebar -->
    <aside class="df-sidebar">
        <div class="df-sidebar-header">
            <p class="df-sidebar-title">Secciones</p>
        </div>
        <div class="df-sidebar-nav">
            <button
                v-for="(section, index) in sections"
                :key="section.id"
                type="button"
                class="df-section-item"
                :class="getSectionState(index)"
                :disabled="index > currentIndex"
                @click="goToSection(index)"
            >
                <span class="df-badge" :class="getSectionState(index)">
                    <span v-if="getSectionState(index) === 'completed'">✓</span>
                    <span v-else>\${ index + 1 }</span>
                </span>
                <span>\${ section.title }</span>
            </button>
        </div>
        <div class="df-progress-wrap">
            <div class="df-progress-label">
                <span>Progreso</span>
                <span>\${ currentIndex + 1 } / \${ totalSections }</span>
            </div>
            <div class="df-progress-bar-bg">
                <div class="df-progress-bar-fill" :style="{ width: progressPercent + '%' }"></div>
            </div>
        </div>
    </aside>

    <!-- Main panel -->
    <div class="df-main">

        <div class="df-section-header">
            <div class="df-section-header-row">
                <div class="df-section-number">\${ currentSection.id }</div>
                <h2 class="df-section-title">\${ currentSection.title }</h2>
            </div>
            <p v-if="currentSection.description" class="df-section-desc">\${ currentSection.description }</p>
        </div>

        <div class="df-questions">
            <template v-for="question in currentSection.questions" :key="question.id">

                <!-- ── REPEATER ── -->
                <div v-if="question.type === 'repeater'" class="df-question">
                    <label>\${ question.label }</label>
                    <div class="df-repeater">
                        <div
                            v-for="(item, idx) in getRepeaterItems(question.id)"
                            :key="idx"
                            class="df-repeater-item"
                        >
                            <div class="df-repeater-item-header">
                                <span class="df-repeater-item-title">\${ question.label } #\${ idx + 1 }</span>
                                <button type="button" class="df-repeater-delete" @click="repeaterRemove(question.id, idx)">
                                    ✕ Eliminar
                                </button>
                            </div>

                            <template v-for="field in question.fields" :key="field.id">
                                <div
                                    v-if="!field.dynamicOptions || getDynamicOptions(field, item).length > 0 || !field.dynamicOptions"
                                    class="df-question"
                                    style="gap:0"
                                >
                                    <label style="margin-bottom:4px;margin-top:0;min-height:unset">
                                        \${ field.label }
                                        <span v-if="field.required" class="df-required">*</span>
                                    </label>

                                    <!-- dropdown inside repeater -->
                                    <select
                                        v-if="field.type === 'dropdown'"
                                        class="df-select"
                                        :value="item[field.id] || ''"
                                        @change="repeaterSet(question.id, idx, field.id, $event.target.value)"
                                    >
                                        <option value="">Selecciona una opción...</option>
                                        <option v-for="opt in field.options" :key="opt" :value="opt">\${ opt }</option>
                                    </select>

                                    <!-- multiple inside repeater -->
                                    <template v-else-if="field.type === 'multiple'">
                                        <p v-if="getDynamicOptions(field, item).length === 0" class="df-hint">
                                            Selecciona primero el sector para ver las opciones.
                                        </p>
                                        <div v-else class="df-btn-group">
                                            <button
                                                v-for="opt in getDynamicOptions(field, item)"
                                                :key="opt"
                                                type="button"
                                                class="df-opt-btn"
                                                :class="{ 'df-selected': Array.isArray(item[field.id]) && item[field.id].includes(opt) }"
                                                @click="repeaterToggleMultiple(question.id, idx, field.id, opt)"
                                            >\${ opt }</button>
                                        </div>
                                    </template>
                                </div>
                            </template>
                        </div>

                        <button type="button" class="df-repeater-add" @click="repeaterAdd(question)">
                            \${ question.addLabel || '+ Agregar' }
                        </button>
                    </div>
                </div>

                <!-- ── Campos normales (con visibilidad condicional) ── -->
                <div
                    v-else-if="isVisible(question, sectionAnswers())"
                    class="df-question"
                >
                    <label>
                        \${ question.label }
                        <span v-if="question.required" class="df-required">*</span>
                    </label>

                    <!-- single -->
                    <div v-if="question.type === 'single'" class="df-btn-group">
                        <button
                            v-for="opt in question.options" :key="opt"
                            type="button" class="df-opt-btn"
                            :class="{ 'df-selected': getAnswer(question.id) === opt }"
                            @click="setAnswer(question.id, opt)"
                        >\${ opt }</button>
                    </div>

                    <!-- multiple -->
                    <div v-else-if="question.type === 'multiple'" class="df-btn-group">
                        <button
                            v-for="opt in question.options" :key="opt"
                            type="button" class="df-opt-btn"
                            :class="{ 'df-selected': Array.isArray(getAnswer(question.id)) && getAnswer(question.id).includes(opt) }"
                            @click="toggleMultiple(question.id, opt)"
                        >
                            <span v-if="Array.isArray(getAnswer(question.id)) && getAnswer(question.id).includes(opt)">✓ </span>\${ opt }
                        </button>
                    </div>

                    <!-- boolean -->
                    <div v-else-if="question.type === 'boolean'" class="df-bool-group">
                        <button type="button" class="df-bool-btn"
                            :class="{ 'df-selected': getAnswer(question.id) === true }"
                            @click="setAnswer(question.id, true)">Sí</button>
                        <button type="button" class="df-bool-btn"
                            :class="{ 'df-selected': getAnswer(question.id) === false }"
                            @click="setAnswer(question.id, false)">No</button>
                    </div>

                    <!-- dropdown -->
                    <select
                        v-else-if="question.type === 'dropdown'"
                        class="df-select"
                        :value="getAnswer(question.id) || ''"
                        @change="setAnswer(question.id, $event.target.value)"
                    >
                        <option value="">Selecciona una opción...</option>
                        <option v-for="opt in question.options" :key="opt" :value="opt">\${ opt }</option>
                    </select>

                    <!-- text (simple o multiline) -->
                    <textarea
                        v-else-if="question.type === 'text' && question.multiline"
                        class="df-textarea"
                        :placeholder="question.placeholder || ''"
                        :value="getAnswer(question.id) || ''"
                        @input="setAnswer(question.id, $event.target.value)"
                    ></textarea>
                    <input
                        v-else-if="question.type === 'text'"
                        class="df-input" type="text"
                        :placeholder="question.placeholder || ''"
                        :value="getAnswer(question.id) || ''"
                        @input="setAnswer(question.id, $event.target.value)"
                    />

                    <!-- date -->
                    <input
                        v-else-if="question.type === 'date'"
                        class="df-input" type="date"
                        :value="getAnswer(question.id) || ''"
                        @input="setAnswer(question.id, $event.target.value)"
                    />

                    <!-- datetime -->
                    <input
                        v-else-if="question.type === 'datetime'"
                        class="df-input" type="datetime-local"
                        :value="getAnswer(question.id) || ''"
                        @input="setAnswer(question.id, $event.target.value)"
                    />
                </div>

            </template>
        </div>

        <div class="df-nav">
            <button type="button" class="df-btn-prev" :disabled="isFirst" @click="goPrev">← Anterior</button>
            <button type="button" class="df-btn-next" @click="goNext">
                <span v-if="isLast">Guardar seguimiento ✓</span>
                <span v-else>Siguiente →</span>
            </button>
        </div>

    </div>
</div>
    `,
});
