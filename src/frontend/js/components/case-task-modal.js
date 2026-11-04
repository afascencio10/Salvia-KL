/**
 * case-task-modal
 * Modal reutilizable para completar una CaseTask.
 * El padre lo abre vía: this.$refs.taskModal.open(taskId)
 * Emite @completed(tareaActualizada) al cerrar con éxito.
 *
 * Props:
 *   currentUserId (String, required) — icode del usuario en sesión
 */

/* ─── Estilos ────────────────────────────────────────────────────────────── */
(function injectCaseTaskModalStyles() {
    if (document.getElementById('case-task-modal-styles')) return;
    const style = document.createElement('style');
    style.id = 'case-task-modal-styles';
    style.textContent = `
        .ctm-overlay {
            position: fixed;
            inset: 0;
            background: rgba(0,0,0,.45);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 1050;
            padding: 16px;
        }

        .ctm-modal {
            position: relative;
            background: #fff;
            border-radius: 16px;
            width: 100%;
            max-width: 540px;
            max-height: 90vh;
            display: flex;
            flex-direction: column;
            box-shadow: 0 20px 40px rgba(0,0,0,.18);
            overflow: hidden;
        }

        .ctm-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 20px 24px 16px;
            border-bottom: 1px solid #e5e7eb;
            flex-shrink: 0;
        }
        .ctm-header-left {
            display: flex;
            align-items: center;
            gap: 12px;
        }
        .ctm-header-icon {
            width: 40px;
            height: 40px;
            border-radius: 10px;
            background: #eff6ff;
            color: #1d4ed8;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 18px;
            flex-shrink: 0;
        }
        .ctm-header-title {
            font-size: 16px;
            font-weight: 700;
            color: #111827;
            margin: 0;
        }
        .ctm-header-sub {
            font-size: 12px;
            color: #9ca3af;
            margin: 2px 0 0;
            max-width: 380px;
            line-height: 1.4;
        }
        .ctm-close-btn {
            background: none;
            border: none;
            font-size: 20px;
            color: #9ca3af;
            cursor: pointer;
            padding: 4px;
            line-height: 1;
            transition: color .15s;
        }
        .ctm-close-btn:hover { color: #4b5563; }

        .ctm-body {
            flex: 1;
            overflow-y: auto;
            padding: 20px 24px;
        }

        .ctm-field {
            margin-bottom: 14px;
        }
        .ctm-label {
            display: block;
            font-size: 12px;
            font-weight: 600;
            color: #374151;
            margin-bottom: 5px;
        }
        .ctm-label .req { color: #dc2626; }
        .ctm-input, .ctm-select, .ctm-textarea {
            width: 100%;
            padding: 9px 12px;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            font-size: 13px;
            color: #111827;
            background: #fff;
            box-sizing: border-box;
            font-family: inherit;
            transition: border-color .15s, box-shadow .15s;
            outline: none;
        }
        .ctm-input:focus, .ctm-select:focus, .ctm-textarea:focus {
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59,130,246,.12);
        }
        .ctm-textarea { resize: vertical; }
        .ctm-select:disabled, .ctm-input:disabled {
            background: #f9fafb;
            color: #9ca3af;
            cursor: not-allowed;
        }

        .ctm-row-2 {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 10px;
        }

        /* Toggle switch */
        .ctm-toggle-wrap {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px 12px;
            background: #f9fafb;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            cursor: pointer;
            user-select: none;
        }
        .ctm-toggle {
            position: relative;
            width: 36px;
            height: 20px;
            flex-shrink: 0;
        }
        .ctm-toggle input { opacity: 0; width: 0; height: 0; position: absolute; }
        .ctm-toggle-track {
            position: absolute;
            inset: 0;
            background: #d1d5db;
            border-radius: 999px;
            transition: background .2s;
        }
        .ctm-toggle input:checked + .ctm-toggle-track { background: #3b82f6; }
        .ctm-toggle-thumb {
            position: absolute;
            top: 2px;
            left: 2px;
            width: 16px;
            height: 16px;
            background: #fff;
            border-radius: 50%;
            transition: transform .2s;
            box-shadow: 0 1px 3px rgba(0,0,0,.2);
        }
        .ctm-toggle input:checked ~ .ctm-toggle-thumb { transform: translateX(16px); }
        .ctm-toggle-label { font-size: 13px; font-weight: 500; color: #374151; }

        /* Checkboxes de decisiones */
        .ctm-checkbox-group {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        .ctm-checkbox-item {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px 12px;
            border: 1.5px solid #e5e7eb;
            border-radius: 8px;
            cursor: pointer;
            transition: border-color .15s, background .15s;
        }
        .ctm-checkbox-item.checked {
            border-color: #3b82f6;
            background: #eff6ff;
        }
        .ctm-checkbox-item input[type="checkbox"] { display: none; }
        .ctm-checkbox-box {
            width: 18px;
            height: 18px;
            border: 2px solid #d1d5db;
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-shrink: 0;
            transition: border-color .15s, background .15s;
        }
        .ctm-checkbox-item.checked .ctm-checkbox-box {
            border-color: #3b82f6;
            background: #3b82f6;
        }
        .ctm-checkbox-check {
            color: #fff;
            font-size: 11px;
            display: none;
        }
        .ctm-checkbox-item.checked .ctm-checkbox-check { display: block; }
        .ctm-checkbox-text { font-size: 13px; color: #374151; }
        .ctm-checkbox-item.checked .ctm-checkbox-text { color: #1d4ed8; font-weight: 500; }

        .ctm-ob-list { display: flex; flex-direction: column; gap: 8px; }
        .ctm-ob-item {
            border: 1.5px solid #e5e7eb;
            border-radius: 10px;
            padding: 12px 14px;
            cursor: pointer;
            background: #fff;
            transition: background .15s, border-color .15s, color .15s;
        }
        .ctm-ob-item.selected {
            background: #5106A7;
            border-color: #5106A7;
            color: #fff;
        }
        .ctm-ob-item-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 10px;
        }
        .ctm-ob-item-label { font-size: 13px; font-weight: 500; line-height: 1.4; }
        .ctm-ob-item.selected .ctm-ob-item-label { color: #fff; }
        .ctm-ob-check { font-size: 13px; flex-shrink: 0; opacity: 0; }
        .ctm-ob-item.selected .ctm-ob-check { opacity: 1; color: #fff; }
        .ctm-ob-desc { margin-top: 10px; }
        .ctm-ob-desc .ctm-textarea { background: #fff; color: #111827; }

        .ctm-context-box {
            padding: 12px 14px;
            background: #f9fafb;
            border: 1.5px solid #e5e7eb;
            border-radius: 10px;
            margin-bottom: 16px;
        }
        .ctm-context-box .ctm-field:last-child { margin-bottom: 0; }

        /* Radio Sí/No para generaOficio */
        .ctm-radio-group {
            display: flex;
            gap: 10px;
        }
        .ctm-radio-item {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            padding: 10px 14px;
            border: 1.5px solid #e5e7eb;
            border-radius: 8px;
            cursor: pointer;
            font-size: 13px;
            font-weight: 500;
            color: #374151;
            transition: border-color .15s, background .15s;
        }
        .ctm-radio-item:has(input:checked) {
            border-color: #3b82f6;
            background: #eff6ff;
            color: #1d4ed8;
        }
        .ctm-radio-item input[type="radio"] { accent-color: #3b82f6; }

        /* Error banner */
        .ctm-error-banner {
            padding: 10px 14px;
            background: #fef2f2;
            border: 1px solid #fecaca;
            border-radius: 8px;
            color: #b91c1c;
            font-size: 13px;
            margin-top: 12px;
        }

        /* Loading */
        .ctm-loading {
            display: flex;
            flex-direction: column;
            align-items: center;
            padding: 40px 0;
            gap: 12px;
            color: #9ca3af;
            font-size: 13px;
        }
        @keyframes ctm-spin { to { transform: rotate(360deg); } }
        .ctm-spinner {
            width: 28px; height: 28px;
            animation: ctm-spin .7s linear infinite;
        }

        /* Separador de sección */
        .ctm-section-title {
            font-size: 11px;
            font-weight: 700;
            color: #9ca3af;
            text-transform: uppercase;
            letter-spacing: .05em;
            margin: 16px 0 10px;
        }
        .ctm-section-title:first-child { margin-top: 0; }

        /* Footer */
        .ctm-footer {
            display: flex;
            gap: 10px;
            padding: 16px 24px;
            border-top: 1px solid #e5e7eb;
            flex-shrink: 0;
        }
        .ctm-btn {
            flex: 1;
            padding: 10px 16px;
            border-radius: 8px;
            font-size: 13px;
            font-weight: 600;
            cursor: pointer;
            transition: background .15s, opacity .15s;
            border: none;
            outline: none;
        }
        .ctm-btn:disabled { opacity: .5; cursor: not-allowed; }
        .ctm-btn-cancel {
            background: #fff;
            border: 1px solid #e5e7eb;
            color: #374151;
        }
        .ctm-btn-cancel:hover:not(:disabled) { background: #f9fafb; }
        .ctm-btn-confirm {
            background: #1d4ed8;
            color: #fff;
        }
        .ctm-btn-confirm:hover:not(:disabled) { background: #1e40af; }
        .ctm-chip {
            display: inline-flex;
            align-items: center;
            border: 1px solid #e5e7eb;
            border-radius: 999px;
            padding: 6px 14px;
            font-size: .84rem;
            color: #374151;
            background: #fff;
            cursor: pointer;
        }
        .ctm-chip.is-on {
            background: #db2777;
            border-color: #db2777;
            color: #fff;
        }
        .ctm-resultado {
            display: block;
            width: 100%;
            text-align: left;
            background: #fff;
            border: 1px solid #e5e7eb;
            border-radius: 10px;
            padding: 10px 14px;
            margin-bottom: 8px;
            color: #111827;
            cursor: pointer;
        }
        .ctm-resultado.is-on {
            border-color: #1d4ed8;
            background: #1d4ed8;
            color: #fff;
            font-weight: 600;
            box-shadow: 0 0 0 3px rgba(29, 78, 216, .28);
        }
        .ctm-chat { display: flex; flex-direction: column; gap: 8px; margin-bottom: 14px; }
        .ctm-chat-msg { border-radius: 10px; padding: 8px 12px; font-size: .84rem; color: #111827; }
        .ctm-chat-msg strong { display: block; margin-bottom: 4px; }
        .ctm-chat-msg p { margin: 0; }
        .ctm-chat-equipo { background: #f3f4f6; }
        .ctm-chat-agente { background: #ede9fe; }
        .ctm-obs { background: #fef2f2; border: 1px solid #fecaca; border-radius: 10px; padding: 12px; margin-bottom: 14px; color: #b91c1c; }
        .ctm-obs strong { display: block; margin-bottom: 4px; }
        .ctm-obs p { margin: 0; }
        .ctm-check-descartar { display: flex; gap: 8px; align-items: center; margin: 4px 0 8px; color: #b91c1c; font-size: .9rem; }
        .ctm-modal.ctm-articulacion .ctm-btn-confirm {
            background: #5106A7;
        }
        .ctm-modal.ctm-articulacion .ctm-btn-confirm:hover:not(:disabled) { background: #3d0780; }
        .ctm-context-box {
            background: #f3f4f6;
            border-radius: 10px;
            padding: 12px 14px;
            font-size: 13px;
            color: #374151;
            line-height: 1.5;
            margin-bottom: 16px;
        }
        /* Panel de resumen (revisar caracterización) */
        .ctm-resumen-row { display: flex; gap: 10px; padding: 6px 0; font-size: 13px; }
        .ctm-resumen-row + .ctm-resumen-row { border-top: 1px solid #e5e7eb; }
        .ctm-resumen-row .ctm-resumen-label { flex: 0 0 150px; color: #6b7280; font-weight: 600; }
        .ctm-resumen-row .ctm-resumen-valor { flex: 1; color: #111827; }
        /* Toggle de resultado aprobado/devuelto (verde/rojo) */
        .ctm-resultado.is-aprobado { border-color: #16a34a; background: #16a34a; color: #fff; font-weight: 600; box-shadow: 0 0 0 3px rgba(22,163,74,.25); }
        .ctm-resultado.is-devuelto { border-color: #dc2626; background: #dc2626; color: #fff; font-weight: 600; box-shadow: 0 0 0 3px rgba(220,38,38,.25); }
        .ctm-btn-confirm.is-aprobar { background: #16a34a; }
        .ctm-btn-confirm.is-aprobar:hover:not(:disabled) { background: #15803d; }
        .ctm-btn-confirm.is-devolver { background: #dc2626; }
        .ctm-btn-confirm.is-devolver:hover:not(:disabled) { background: #b91c1c; }
        .ctm-toggle-row {
            display: flex;
            gap: 8px;
        }
        .ctm-toggle-btn {
            flex: 1;
            padding: 8px 12px;
            border-radius: 8px;
            border: 1px solid #e5e7eb;
            background: #fff;
            color: #374151;
            font-size: 13px;
            font-weight: 600;
            cursor: pointer;
        }
        .ctm-toggle-btn.active {
            background: #5106A7;
            border-color: #5106A7;
            color: #fff;
        }
        .ctm-hint {
            font-size: 12px;
            color: #6b7280;
            margin: 6px 0 0;
            line-height: 1.4;
        }

        /* Entrevista de caracterización (rol am) ───────────────────────────── */
        /* El modal se agranda en el paso 3 (entrevista completa). */
        .ctm-modal.ctm-caracterizacion-grande { max-width: 680px; }
        /* Línea verde de contacto efectivo confirmado. */
        .ctm-car-contacto {
            display: flex; align-items: center; gap: 8px;
            color: #15803d; font-size: 13px; font-weight: 600;
            margin: 8px 0 0;
        }
        .ctm-car-caso {
            font-size: 13px; color: #374151; margin: 0;
        }
        /* Botones toggle verde/rojo (continuar / concluir, otorga / no otorga). */
        .ctm-car-toggle {
            display: block; width: 100%; text-align: left;
            background: #fff; border: 1.5px solid #e5e7eb; border-radius: 10px;
            padding: 12px 16px; margin-bottom: 10px; color: #111827;
            font-size: 14px; font-weight: 600; cursor: pointer;
            transition: border-color .15s, background .15s, color .15s, box-shadow .15s;
        }
        .ctm-car-toggle:hover:not(.is-on) { border-color: #c7d2fe; }
        .ctm-car-toggle.is-on-green {
            border-color: #16a34a; background: #16a34a; color: #fff;
            box-shadow: 0 0 0 3px rgba(22,163,74,.25);
        }
        .ctm-car-toggle.is-on-red {
            border-color: #dc2626; background: #dc2626; color: #fff;
            box-shadow: 0 0 0 3px rgba(220,38,38,.25);
        }
        /* Bloque azul PLACEHOLDER del formulario dinámico (FASE 2). */
        .ctm-car-placeholder {
            background: #eff6ff; border: 1.5px dashed #93c5fd; border-radius: 10px;
            padding: 16px 18px; margin: 4px 0 16px; color: #1e40af;
        }
        .ctm-car-placeholder .ph-title { font-weight: 800; margin: 0 0 6px; font-size: 14px; }
        .ctm-car-placeholder .ph-body { margin: 0; font-size: 12.5px; line-height: 1.5; color: #1d4ed8; }
        .ctm-car-chip {
            display: inline-flex; align-items: center;
            border: 1.5px solid #e5e7eb; border-radius: 999px;
            padding: 7px 14px; font-size: .84rem; color: #374151;
            background: #fff; cursor: pointer; user-select: none;
            transition: border-color .15s, background .15s, color .15s;
        }
        .ctm-car-chip.is-on {
            background: #16a34a; border-color: #16a34a; color: #fff; font-weight: 600;
        }
        .ctm-car-footer {
            display: flex; gap: 10px; margin-top: 18px;
            padding-top: 16px; border-top: 1px solid #e5e7eb;
        }
    `;
    document.head.appendChild(style);
})();

// Sectores para el filtro de entidades de "Proyectar oficio" — GET /api/v1/entity-branches
// espera estas palabras (no códigos de 2 letras). Mismo catálogo que case-entities.js /
// notifications.js.
const CTM_SECTORES_MODAL = [
    { value: 'salud', label: 'Salud' },
    { value: 'justicia', label: 'Justicia' },
    { value: 'proteccion', label: 'Protección' },
    { value: 'otras_instituciones', label: 'Otras instituciones' },
];

// UUID fijo del formulario dinámico "Entrevista de Caracterización" (FASE 2b-1).
// Debe coincidir con el seed (internal/seed/caracterizacion_form.sql) y con la
// constante Go constants.FormIDCaracterizacion / seed.CaracterizacionFormID. Se
// embebe aquí igual que hacer_seguimiento.html embebe el UUID del form de
// seguimiento: el formId del nuevo formulario llega al componente <dinamic-form>
// desde esta constante, sin endpoint nuevo.
const CTM_CARACTERIZACION_FORM_ID = 'f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d';

/* ─── Componente Vue ──────────────────────────────────────────────────────── */
app.component('case-task-modal', {
    delimiters: ['${', '}'],

    props: {
        currentUserId: { type: String, required: true },
    },

    data() {
        return {
            visible:       false,
            tarea:         null,
            cargandoTarea: false,
            errorTarea:    null,
            guardando:     false,
            completado:    false,
            saveError:     null,

            // Ubicaciones
            departamentos: [],
            allCiudades:   [],
            ciudades:      [],
            municipios:    [],
            entidades:     [],
            entidadBusqueda: '',
            entidadDropdownOpen: false,
            sectoresModal: CTM_SECTORES_MODAL,

            // Formulario (shape varía según tarea.type)
            form: {},

            // Entrevista de Caracterización (FASE 2b-1): formId fijo + formState con
            // los catálogos que el motor resuelve vía state_options_path. El
            // componente <dinamic-form> guarda sus propias secciones vía
            // POST /forms/saveSection de forma independiente (su propio submission);
            // en esta fase NO se enlaza ese submission a la tarea/caso (eso es 2b-2).
            caracterizacionFormId: CTM_CARACTERIZACION_FORM_ID,
            caracterizacionFormState: {
                // Resuelto por GET /victim-case-forms/enums (categorías
                // victim_case_form2_*) y GET /locations/*.
                enums: {},
                departments: [],
                geo: { residencia: { cities: [], towns: [] } },
            },
            caracterizacionFormStateCargado: false,
            // FASE 2b-2: submission precargado y enlazado a la tarea. El <dinamic-form>
            // solo se monta con submissionId y con el formState ya cargado.
            caracterizacionFormStateListo: false,
            caracterizacionSubmissionId: '',
            caracterizacionInitEstado: '',   // '' | 'cargando' | 'listo' | 'error'
            caracterizacionInitError: '',
            // true cuando el <dinamic-form> emitió form-completed (todas las secciones
            // visibles respondidas) o, si el submission se reutilizó, cuando /load lo confirma.
            caracterizacionFormCompleto: false,

            // Datos del entity_letter vinculado (para Corregir oficio)
            oficioVinculado: null,
            cargandoOficio:  false,
            cargandoObligaciones: false,
            errorObligaciones: null,
            entityName: '',
            articulacionCtx: {
                barrierId: '',
                caseId: '',
                entityBranchId: null,
                entityName: '',
                barrierDescription: '',
            },
            mecanismoCtx: {
                mechanismId: '',
                caseId: '',
                scope: '',
            },

            // Activar Enlace (comite_caso)
            enlacesDepartamento: [],   // enlaces 'en' del departamento del caso
            cargandoEnlaces: false,
            oeCiudades: [],            // cascada "otra entidad"
            oeMunicipios: [],
            oeEntidades: [],
            cargandoOeEntidades: false,
            agentesMedidas: [],
            cargandoAgentes: false,
            errorAgentes: '',
            intentosPrevios: 0,
            cargandoIntentos: false,
            medidasRemitidas: [],
            cargandoMedidas: false,
            errorMedidas: '',
            casoRevisionNombre: '',
            casoRevisionCodigo: '',
            cargandoCasoRevision: false,
            hiloChat: [],
            cargandoChat: false,
        };
    },

    computed: {
        entidadesFiltradas() {
            if (!this.entidadBusqueda) return this.entidades;
            var q = this.entidadBusqueda.toLowerCase();
            return this.entidades.filter(function(e) { return e.name.toLowerCase().indexOf(q) !== -1; });
        },
        titulo() {
            if (!this.tarea) return 'Completar tarea';
            if (this.tarea.type === 'solicitud_entidad') {
                return 'Solicitudes a ' + (this.entityName || this.tarea.entityName || 'la entidad');
            }
            if (this.tarea.type === 'articulacion_territorial') {
                return 'Gestionar barrera';
            }
            if (this.tarea.type === 'gestion_mecanismo') {
                var scope = (this.mecanismoCtx && this.mecanismoCtx.scope) || this.tarea.scope || '';
                if (scope === 'municipal') return 'Gestión del Comité Municipal';
                if (scope === 'departamental') return 'Gestión del Comité Departamental';
                if (scope === 'nacional') return 'Gestión de la Instancia Nacional';
                return 'Gestión del mecanismo';
            }
            return {
                gestion_llamada:  'Gestión de Llamada',
                proyectar_oficio: 'Proyectar Oficio',
                comite_caso:      'Decisiones del Comité',
                'Corregir oficio': 'Corregir Oficio',
                justificar_remision: 'Justificar Remisión',
                respuesta_entidad: 'Respuesta institucional',
                decidir_via_riesgo_inminente: 'Decidir vía de atención — Riesgo inminente',
                intentar_contacto_3x3: 'Registrar intento de contacto',
                entrevista_caracterizacion: 'Gestionar caso — Medidas de Emergencia',
                revisar_caracterizacion: 'Revisar caracterización',
                revisar_trazabilidad: 'Revisar trazabilidad del caso',
                cargar_cotizaciones: 'Cargar cotizaciones',
                aprobar_cotizaciones: 'Revisar y aprobar cotizaciones',
                gestionar_territorio: 'Gestionar medida en territorio',
                marcar_otorgadas: 'Marcar medidas como otorgadas',
                escalar_cruz_roja: 'Escalar a Cruz Roja',
                entregar_medida_inmediata: 'Entregar medida inmediata',
                confirmar_lugar_seguro: 'Confirmar que la persona está en lugar seguro',
                revision_equipo_tecnico: 'Revisión Equipo Técnico',
                asignar_agente_medidas: 'Reasignar remisiones',
                ajustar_remision_devuelta: 'Remisión devuelta por el Equipo Técnico',
            }[this.tarea.type] || 'Completar tarea';
        },

        intentoContacto() {
            return (this.intentosPrevios || 0) + 1;
        },

        articulacionSubtitulo() {
            const name = (this.articulacionCtx && this.articulacionCtx.entityName) || '';
            const desc = (this.articulacionCtx && this.articulacionCtx.barrierDescription) || '';
            return [name, desc].filter(Boolean).join(' — ');
        },

        iconoTipo() {
            if (!this.tarea) return 'fa-tasks';
            return {
                gestion_llamada:  'fa-phone-alt',
                proyectar_oficio: 'fa-file-signature',
                comite_caso:      'fa-users',
                'Corregir oficio': 'fa-pencil-alt',
                justificar_remision: 'fa-reply',
                solicitud_entidad: 'fa-folder',
                respuesta_entidad: 'fa-reply',
                articulacion_territorial: 'fa-link',
                gestion_mecanismo: 'fa-balance-scale',
                decidir_via_riesgo_inminente: 'fa-fire',
                intentar_contacto_3x3: 'fa-phone',
                entrevista_caracterizacion: 'fa-user-check',
                revisar_caracterizacion: 'fa-clipboard-check',
                revisar_trazabilidad: 'fa-folder-open',
                cargar_cotizaciones: 'fa-file-invoice-dollar',
                aprobar_cotizaciones: 'fa-file-invoice-dollar',
                gestionar_territorio: 'fa-map-marked-alt',
                marcar_otorgadas: 'fa-clipboard-check',
                escalar_cruz_roja: 'fa-ambulance',
                entregar_medida_inmediata: 'fa-box',
                confirmar_lugar_seguro: 'fa-shield-alt',
                revision_equipo_tecnico: 'fa-clipboard-check',
                asignar_agente_medidas: 'fa-user-plus',
                ajustar_remision_devuelta: 'fa-undo',
            }[this.tarea.type] || 'fa-tasks';
        },

        ultimoMensajeEquipo() {
            var items = this.hiloChat || [];
            for (var i = items.length - 1; i >= 0; i--) {
                if (items[i].lado === 'equipo') return items[i];
            }
            return null;
        },
        hiloAnterior() {
            var items = this.hiloChat || [];
            var ultimo = this.ultimoMensajeEquipo;
            if (!ultimo) return items;
            var out = [];
            var omitido = false;
            for (var i = items.length - 1; i >= 0; i--) {
                if (!omitido && items[i] === ultimo) {
                    omitido = true;
                    continue;
                }
                out.unshift(items[i]);
            }
            return out;
        },
        textoDevolucion() {
            var resultado = this.ultimoMensajeEquipo && this.ultimoMensajeEquipo.resultado;
            var motivo = resultado === 'devuelto_para_acompanamiento' ? 'para acompañamiento' : 'para ajuste';
            return 'El Equipo Técnico revisó esta remisión a Medidas de Emergencia y la devolvió ' + motivo + '. Corrige o complementa la información señalada y reenvíala para una nueva revisión, o descarta la remisión si ya no aplica.';
        },
        tareaFormData() {
            var fd = this.tarea && this.tarea.formData;
            if (!fd) return {};
            if (typeof fd === 'string') {
                try { return JSON.parse(fd) || {}; } catch (e) { return {}; }
            }
            return fd;
        },

        // medidasRequeridasTexto: medidas reales de la remisión (labels que expone
        // el backend en tarea.referralMeasures), unidas por coma. Vacío si no hay.
        medidasRequeridasTexto() {
            var m = (this.tarea && this.tarea.referralMeasures) || [];
            if (!Array.isArray(m) || m.length === 0) return '';
            return m.join(', ');
        },

        // resumenRevisarCaracterizacion: filas del panel de resumen (solo lectura)
        // del modal de revisar caracterización. Los datos reales los expone el backend
        // en tarea.caracterizacionResumen (form_data de la última entrevista Done:
        // hallazgos, planAccion, enlaceCruzRoja, enlaceDeclaracion). Si falta algo
        // se muestra "—". "Medidas requeridas" sale de las medidas reales de la remisión.
        resumenRevisarCaracterizacion() {
            var r = (this.tarea && this.tarea.caracterizacionResumen) || {};
            var v = function(x) { var t = String(x || '').trim(); return t || '—'; };
            var medidas = this.medidasRequeridasTexto;
            return [
                { label: 'Hallazgos', valor: v(r.hallazgos) },
                { label: 'Medidas requeridas', valor: medidas || 'Sin medidas registradas en la remisión.' },
                { label: 'Plan de acción', valor: v(r.planAccion) },
                { label: 'Remisión Cruz Roja', valor: v(r.enlaceCruzRoja) },
                { label: 'Declaración jurada', valor: v(r.enlaceDeclaracion) },
            ];
        },

        // resumenRevisarTrazabilidad: filas del panel de resumen (solo lectura) del
        // modal de revisar trazabilidad (rol sc). "Medidas requeridas" es dato real
        // (tarea.referralMeasures); "Plan de acción" es placeholder por ahora.
        resumenRevisarTrazabilidad() {
            var medidas = this.medidasRequeridasTexto;
            return [
                { label: 'Medidas requeridas', valor: medidas || 'Sin medidas registradas en la remisión.' },
                { label: 'Plan de acción', valor: 'Pendiente de la caracterización del agente.' },
            ];
        },

        // cotizacionMedidas: lista de labels de medidas para el modal de cargar
        // cotizaciones. Son las medidas reales de la remisión (tarea.referralMeasures);
        // si la remisión no tiene ninguna, se usan dos por defecto para poder probar
        // (fallback solo en el front; NO se crean en BD).
        cotizacionMedidas() {
            var m = (this.tarea && this.tarea.referralMeasures) || [];
            if (Array.isArray(m) && m.length > 0) return m.slice();
            return ['Alojamiento temporal', 'Kit de alimentación'];
        },

        // cotizacionesAprobarResumen: filas { label, valor } del panel de resumen
        // (solo lectura) del modal "Revisar y aprobar cotizaciones" (rol sc). Los
        // montos los cargó Cruz Roja en la tarea cargar_cotizaciones y el backend los
        // expone en tarea.cotizacionesCargadas ({ "<medida>": valor }). Fallback: si
        // no hay cotizaciones previas, lista las medidas de la remisión
        // (tarea.referralMeasures) con valor 0, para no dejar el panel vacío.
        cotizacionesAprobarResumen() {
            var cot = (this.tarea && this.tarea.cotizacionesCargadas) || null;
            var filas = [];
            if (cot && typeof cot === 'object' && Object.keys(cot).length > 0) {
                Object.keys(cot).forEach(function(label) {
                    var num = Number(cot[label]);
                    filas.push({ label: label, valor: isFinite(num) ? num : 0 });
                });
                return filas;
            }
            // Fallback: medidas de la remisión con valor 0.
            (this.cotizacionMedidas || []).forEach(function(label) {
                filas.push({ label: label, valor: 0 });
            });
            return filas;
        },

        // medidasAprobadasChips: labels de las medidas aprobadas de la remisión, para
        // los chips de solo lectura del modal "Marcar medidas como otorgadas" (rol cr).
        // Reutiliza tarea.referralMeasures (misma fuente que "Medidas requeridas").
        medidasAprobadasChips() {
            var m = (this.tarea && this.tarea.referralMeasures) || [];
            if (!Array.isArray(m)) return [];
            return m.slice();
        },

        // ── Entrevista de caracterización (rol am) ──
        // caracterizacionNombre / Codigo / FechaContacto: datos de cabecera del
        // PASO 1, expuestos por GET /case-tasks/:id (ver lookupCaracterizacionHeader
        // en el controller). Degradan a placeholder razonable si el backend no los
        // trae, para no romper el modal.
        caracterizacionNombre() {
            return (this.tarea && this.tarea.victimName) || 'Víctima sin nombre registrado';
        },
        caracterizacionCodigo() {
            return (this.tarea && this.tarea.caseCode) || (this.tarea && this.tarea.caseId) || '';
        },
        caracterizacionFechaContacto() {
            return (this.tarea && this.tarea.contactoEfectivoFecha) || '';
        },
        // caracterizacionMedidasCatalogo: catálogo de labels de tipos de medida para
        // los chips seleccionables de "Medidas requeridas". Fuente: los labels del
        // enum de tipos de medida (EmergencyMeasureTypeLabel del backend; aquí se
        // replican como lista fija de los tipos conocidos). Se UNE con las medidas ya
        // precargadas de la remisión (tarea.referralMeasures) por si alguna no está
        // en el catálogo base, así todas quedan visibles y marcables.
        caracterizacionMedidasCatalogo() {
            var base = [
                'Alojamiento',
                'Transporte',
                'Alimentación',
                'Vestuario',
                'Aseo',
                'Apoyo psicosocial',
                'Apoyo exequial',
                'Acción simbólica',
                'Otras ME',
                'Otras medidas de emergencia',
            ];
            var pre = (this.tarea && Array.isArray(this.tarea.referralMeasures)) ? this.tarea.referralMeasures : [];
            pre.forEach(function(m) { if (base.indexOf(m) === -1) base.push(m); });
            return base;
        },

        // totalCotizacionesAprobar: suma de los valores del resumen de cotizaciones.
        totalCotizacionesAprobar() {
            return (this.cotizacionesAprobarResumen || []).reduce(function(acc, f) {
                var num = Number(f.valor);
                return acc + (isFinite(num) ? num : 0);
            }, 0);
        },

        formularioValido() {
            if (!this.tarea || this.cargandoTarea) return false;
            const f = this.form;
            switch (this.tarea.type) {
                case 'gestion_llamada':
                    if (!f.departamentoId || !f.ciudadId || !f.municipioId || !f.entidadId) return false;
                    if (f.entidadId === 'otra' && !f.entidadNombre) return false;
                    if (!f.funcionario) return false;
                    if (f.generaOficio === null || f.generaOficio === undefined) return false;
                    if (f.generaOficio === true && (!f.asunto || !f.rutaKofax)) return false;
                    return true;
                case 'proyectar_oficio':
                    // Nombre de la entidad, Funcionario, Asunto, Ruta Kofax y Prioridad son
                    // opcionales — solo la ubicación y la entidad elegida son obligatorias.
                    if (!f.departamentoId || !f.ciudadId || !f.municipioId || !f.entidadId) return false;
                    return true;
                case 'comite_caso':
                    if (!f.analisisMesaTecnica || !f.analisisMesaTecnica.trim()) return false;
                    if (!f.decisiones || f.decisiones.length === 0) return false;
                    if (f.decisiones.includes('activar_enlace')) {
                        if (!f.solicitudEnlace || !f.solicitudEnlace.trim()) return false;
                        if (!f.entidadDestino) return false;
                        if (f.entidadDestino === 'otra' && !f.oeEntidadId) return false;
                        // Si hay varios enlaces, debe elegir uno
                        if (this.enlacesDepartamento.length > 1 && !f.enlaceSeleccionado) return false;
                    }
                    if (f.decisiones.includes('mecanismo_articulador')) {
                        if (!f.nivelMecanismo) return false;
                        if (!f.quienLlenaMecanismo) return false;
                    }
                    if (f.decisiones.includes('oficio')) {
                        // Debe definirse quién gestiona el oficio (equipo técnico o agente)
                        if (!f.gestionaOficio) return false;
                    }
                    return true;
                case 'Corregir oficio':
                    // Sin campos obligatorios — solo confirmación
                    return !this.cargandoOficio;
                case 'justificar_remision':
                    return !!(f.respuesta && f.respuesta.trim());
                case 'solicitud_entidad':
                    if (this.cargandoObligaciones || !f.items || !f.seleccion) return false;
                    var seleccionados = f.items.filter(function(item) {
                        return f.seleccion[item.id] && f.seleccion[item.id].selected;
                    });
                    if (seleccionados.length < 1) return false;
                    return seleccionados.every(function(item) {
                        var d = f.seleccion[item.id] && f.seleccion[item.id].descripcion;
                        return d && String(d).trim();
                    });
                case 'respuesta_entidad':
                    return !!(f.estadoInstitucional);
                case 'articulacion_territorial':
                    if (!f.fechaContacto || !(f.compromiso && String(f.compromiso).trim())) return false;
                    if (f.requiereSeguimiento === null || f.requiereSeguimiento === undefined) return false;
                    if (f.requiereSeguimiento && !f.fechaSeguimiento) return false;
                    return true;
                case 'gestion_mecanismo': {
                    var scope = (this.mecanismoCtx && this.mecanismoCtx.scope) || '';
                    if (!f.reporteGestion || !String(f.reporteGestion).trim()) return false;
                    if (scope === 'municipal') {
                        return !!(f.abordadoEnSesion);
                    }
                    if (scope === 'departamental') {
                        return !!(f.fechaProxima && String(f.detalleProxima || '').trim());
                    }
                    if (scope === 'nacional') {
                        return !!(f.asistenciasFecha && String(f.asistenciasTexto || '').trim());
                    }
                    return false;
                }
                case 'decidir_via_riesgo_inminente':
                    return (f.via === 'A' || f.via === 'B') && !!f.agentId;
                case 'intentar_contacto_3x3':
                    if (f.contesto !== 'si' && f.contesto !== 'no') return false;
                    if (this.cargandoIntentos) return false;
                    if (f.contesto === 'no' && this.intentoContacto === 9) {
                        return f.accion === 'seguir' || f.accion === 'cerrar';
                    }
                    return true;
                case 'revisar_caracterizacion':
                    // Válido si se aprueba, o si se devuelve con observaciones no vacías.
                    if (f.resultado === 'aprobado') return true;
                    if (f.resultado === 'devuelto') return !!(f.observaciones && String(f.observaciones).trim());
                    return false;
                case 'revisar_trazabilidad':
                    // Válido si se preaprueba, o si se devuelve con observaciones no vacías.
                    if (f.resultado === 'preaprobado') return true;
                    if (f.resultado === 'devuelto') return !!(f.observaciones && String(f.observaciones).trim());
                    return false;
                case 'cargar_cotizaciones': {
                    // Válido solo cuando TODAS las medidas mostradas tienen un valor
                    // numérico >= 0. Se acepta 0; NO se acepta vacío/null ni texto.
                    var medidas = this.cotizacionMedidas;
                    if (!Array.isArray(medidas) || medidas.length === 0) return false;
                    var valores = (f && f.cotizaciones) || {};
                    for (var i = 0; i < medidas.length; i++) {
                        var raw = valores[medidas[i]];
                        if (raw === undefined || raw === null || String(raw).trim() === '') return false;
                        var num = Number(raw);
                        if (!isFinite(num) || num < 0) return false;
                    }
                    return true;
                }
                case 'aprobar_cotizaciones':
                    // Válido si se aprueba, o si no cumple con observaciones no vacías.
                    if (f.resultado === 'aprobado') return true;
                    if (f.resultado === 'no_cumple') return !!(f.observaciones && String(f.observaciones).trim());
                    return false;
                case 'gestionar_territorio':
                    // Ambos obligatorios: fecha de gestión y acciones (trim) no vacías.
                    return !!(f.fechaGestion && String(f.fechaGestion).trim()) &&
                           !!(f.acciones && String(f.acciones).trim());
                case 'marcar_otorgadas':
                    // Ambos obligatorios: fecha de otorgamiento y comentarios (trim) no vacíos.
                    return !!(f.fechaOtorgamiento && String(f.fechaOtorgamiento).trim()) &&
                           !!(f.comentarios && String(f.comentarios).trim());
                case 'entrevista_caracterizacion':
                    // Validez según la rama activa. El botón de submit solo existe en
                    // los estados terminales de cada rama (ver plantilla/footer).
                    if (f.continua === 'concluir') {
                        // Rama A — concluir gestión: motivo + observación obligatorios.
                        return !!f.motivoConclusion && !!(f.observacionConcluir && String(f.observacionConcluir).trim());
                    }
                    if (f.continua === 'consentimiento' && f.consentimiento === 'no_otorga') {
                        // Rama B — no otorga: observación obligatoria.
                        return !!(f.observacionSinConsentimiento && String(f.observacionSinConsentimiento).trim());
                    }
                    if (f.continua === 'consentimiento' && f.consentimiento === 'otorga') {
                        // Rama caracterizado: obligatorios = hallazgos + plan de acción
                        // (medidas y enlaces son opcionales) Y el formulario dinámico
                        // completo (form-completed). Decisión de diseño FASE 2b-2.
                        return !!(f.hallazgos && String(f.hallazgos).trim()) &&
                               !!(f.planAccion && String(f.planAccion).trim()) &&
                               this.caracterizacionFormCompleto === true;
                    }
                    // En pasos intermedios no hay submit habilitado.
                    return false;
                case 'escalar_cruz_roja':
                case 'confirmar_lugar_seguro':
                    return true;
                case 'revision_equipo_tecnico':
                    return !!(f.concepto && String(f.concepto).trim()) && (
                        f.resultado === 'aprobado_para_caracterizacion_me' ||
                        f.resultado === 'devuelto_para_ajuste' ||
                        f.resultado === 'devuelto_para_acompanamiento'
                    );
                case 'ajustar_remision_devuelta':
                    return !!f.descartar || !!(f.ajuste && String(f.ajuste).trim());
                case 'asignar_agente_medidas':
                    return !!f.agentId && !this.cargandoAgentes;
                case 'entregar_medida_inmediata':
                    return !this.cargandoMedidas && Array.isArray(f.medidas) && f.medidas.length > 0;
                default:
                    return false;
            }
        },
    },

    async mounted() {
        try {
            const [deptRes, citiesRes] = await Promise.all([
                fetch('/api/v1/locations/departments'),
                fetch('/api/v1/locations/cities'),
            ]);
            if (deptRes.ok)   this.departamentos = await deptRes.json();
            if (citiesRes.ok) this.allCiudades   = await citiesRes.json();
        } catch (e) {
            console.warn('[case-task-modal] Error cargando ubicaciones al montar:', e);
        }
    },

    methods: {
        /* ── Public API ──────────────────────────────────────────────────── */

        async open(taskId) {
            this.visible       = true;
            this.tarea         = null;
            this.saveError     = null;
            this.errorTarea    = null;
            this.cargandoTarea = true;
            this.guardando     = false;
            this.ciudades      = [];
            this.municipios    = [];
            this.entidades     = []; this.entidadBusqueda = ""; this.entidadDropdownOpen = false;
            this.oficioVinculado = null;
            this.cargandoObligaciones = false;
            this.errorObligaciones = null;
            this.entityName = '';
            this.intentosPrevios = 0;
            this.cargandoIntentos = false;

            try {
                const res = await fetch('/api/v1/case-tasks/' + taskId);
                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }
                if (!res.ok) throw new Error('Error ' + res.status + ' al cargar la tarea.');
                this.tarea = await res.json();
                this._initForm();
            } catch (e) {
                this.errorTarea = 'No se pudo cargar la tarea. ' + e.message;
            } finally {
                this.cargandoTarea = false;
            }
        },

        async openArticulacion(ctx) {
            ctx = ctx || {};
            this.articulacionCtx = {
                barrierId: ctx.barrierId || '',
                caseId: ctx.caseId || '',
                entityBranchId: ctx.entityBranchId || null,
                entityName: ctx.entityName || '',
                barrierDescription: ctx.barrierDescription || '',
            };
            this.saveError = null;
            this.errorTarea = null;
            this.guardando = false;
            this.completado = false;
            if (ctx.taskId) {
                await this.open(ctx.taskId);
                return;
            }
            this.visible = true;
            this.cargandoTarea = false;
            this.tarea = {
                id: null,
                type: 'articulacion_territorial',
                description: 'Gestión propia',
                barrierId: ctx.barrierId || '',
                caseId: ctx.caseId || '',
            };
            this._initForm();
        },

        openMecanismo(ctx) {
            ctx = ctx || {};
            this.mecanismoCtx = {
                mechanismId: ctx.mechanismId || '',
                caseId: ctx.caseId || '',
                scope: ctx.scope || '',
            };
            this.saveError = null;
            this.errorTarea = null;
            this.guardando = false;
            this.completado = false;
            this.visible = true;
            this.cargandoTarea = false;
            this.tarea = {
                id: null,
                type: 'gestion_mecanismo',
                scope: this.mecanismoCtx.scope,
                description: 'Gestión mecanismo',
                caseId: ctx.caseId || '',
            };
            this._initForm();
        },

        cancelar() {
            this.visible         = false;
            this.tarea           = null;
            this.form            = {};
            this.saveError       = null;
            this.errorTarea      = null;
            this.guardando       = false;
            this.completado      = false;
            this.ciudades        = [];
            this.municipios      = [];
            this.entidades       = [];
            this.oficioVinculado = null;
            this.cargandoObligaciones = false;
            this.errorObligaciones = null;
            this.entityName = '';
            this.intentosPrevios = 0;
            this.cargandoIntentos = false;
            this.articulacionCtx = {
                barrierId: '',
                caseId: '',
                entityBranchId: null,
                entityName: '',
                barrierDescription: '',
            };
            this.mecanismoCtx = {
                mechanismId: '',
                caseId: '',
                scope: '',
            };
        },

        async _cargarOficioVinculado() {
            const letterId = this.tarea && this.tarea.entityLetterId;
            if (!letterId) {
                console.warn('[CTM] Corregir oficio — sin entityLetterId en la tarea');
                return;
            }
            this.cargandoOficio = true;
            console.log('[CTM] E01b cargando entity_letter:', letterId);
            try {
                const res = await fetch('/api/v1/entity-letters/' + letterId);
                if (res.ok) {
                    this.oficioVinculado = await res.json();
                    console.log('[CTM] E01b oficio cargado: state=' + this.oficioVinculado.state
                        + ' | reasonCorrection=' + this.oficioVinculado.reasonCorrection
                        + ' | urlKofax=' + this.oficioVinculado.urlKofax);
                } else {
                    console.warn('[CTM] E01b error cargando entity_letter:', res.status);
                }
            } catch (e) {
                console.warn('[CTM] E01b error de red cargando entity_letter:', e);
            } finally {
                this.cargandoOficio = false;
            }
        },

        async _cargarObligacionesPendientes() {
            const guardado = this.tareaFormData;
            if (guardado && Array.isArray(guardado.checklist)) {
                const items = guardado.checklist.map(function(item) {
                    return {
                        id: item.obligationId,
                        label: item.label || item.obligationId || '',
                    };
                });
                const seleccion = {};
                items.forEach(function(item) {
                    seleccion[item.id] = { selected: false, descripcion: '' };
                });
                this.entityName = (this.tarea && this.tarea.entityName) || '';
                this.form = { items: items, seleccion: seleccion };
                return;
            }
            const entityCaseId = this.tarea && this.tarea.entityCaseId;
            if (!entityCaseId) {
                this.errorObligaciones = 'Esta tarea no está vinculada a una entidad del caso.';
                return;
            }
            this.cargandoObligaciones = true;
            this.errorObligaciones = null;
            try {
                const res = await fetch('/api/v1/entity-cases/' + entityCaseId + '/obligations?status=pendiente');
                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }
                if (!res.ok) throw new Error('Error ' + res.status);
                const data = await res.json();
                this.entityName = data.entityName || this.tarea.entityName || '';
                const items = Array.isArray(data.items) ? data.items : [];
                const seleccion = {};
                items.forEach(function(item) {
                    seleccion[item.id] = { selected: false, descripcion: '' };
                });
                this.form = { items: items, seleccion: seleccion };
            } catch (e) {
                this.errorObligaciones = 'No se pudieron cargar las obligaciones pendientes.';
            } finally {
                this.cargandoObligaciones = false;
            }
        },

        toggleObligacion(itemId) {
            if (!this.form.seleccion[itemId]) {
                this.form.seleccion[itemId] = { selected: false, descripcion: '' };
            }
            this.form.seleccion[itemId].selected = !this.form.seleccion[itemId].selected;
            if (!this.form.seleccion[itemId].selected) {
                this.form.seleccion[itemId].descripcion = '';
            }
        },

        /* ── Entrevista de caracterización: formState + eventos del form ──── */

        // Arma el formState del <dinamic-form> con los catálogos que el motor
        // resuelve por state_options_path, replicando lo que hace
        // hacer_seguimiento.html (enums vía /victim-case-forms/enums, ubicación
        // vía /locations/*). Se cargan todas las ciudades en geo.residencia.cities
        // (sin cascada por departamento en esta fase de solo render); towns queda
        // vacío hasta 2b-2. Idempotente: no re-consulta si ya cargó.
        async cargarCaracterizacionFormState() {
            if (this.caracterizacionFormStateCargado) return;
            this.caracterizacionFormStateCargado = true;
            try {
                const [enumsRes, deptRes, citiesRes] = await Promise.all([
                    fetch('/api/v1/victim-case-forms/enums'),
                    fetch('/api/v1/locations/departments'),
                    fetch('/api/v1/locations/cities'),
                ]);
                if (enumsRes.ok)  this.caracterizacionFormState.enums = await enumsRes.json();
                if (deptRes.ok)   this.caracterizacionFormState.departments = await deptRes.json();
                if (citiesRes.ok) this.caracterizacionFormState.geo.residencia.cities = await citiesRes.json();
            } catch (e) {
                console.warn('[case-task-modal] Error cargando catálogos de caracterización:', e);
                // Se permite re-intento en una próxima apertura del modal.
                this.caracterizacionFormStateCargado = false;
            } finally {
                // El <dinamic-form> lee el formState al montarse: solo se monta cuando
                // los catálogos ya se intentaron cargar (con o sin éxito).
                this.caracterizacionFormStateListo = true;
            }
        },

        // FASE 2b-2: pide al backend el submission precargado (editable) y enlazado
        // a la tarea. Idempotente: el backend reutiliza el submission si ya existe.
        // Mientras no haya submissionId el <dinamic-form> no se monta.
        async iniciarCaracterizacionForm() {
            if (!this.tarea || !this.tarea.id) return;
            if (this.caracterizacionInitEstado === 'cargando') return;
            this.caracterizacionInitEstado = 'cargando';
            this.caracterizacionInitError = '';
            try {
                const res = await fetch('/api/v1/case-tasks/' + this.tarea.id + '/init-caracterizacion', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                });
                let body = {};
                try { body = await res.json(); } catch (e) { body = {}; }
                if (!res.ok || !body.submissionId) {
                    throw new Error(body.error || ('Error ' + res.status));
                }
                this.caracterizacionSubmissionId = body.submissionId;
                this.caracterizacionInitEstado = 'listo';
                if (body.reused) {
                    this.verificarCaracterizacionCompleta(body.submissionId);
                }
            } catch (e) {
                console.warn('[CTM] init-caracterizacion falló:', e);
                this.caracterizacionInitError = (e && e.message) || 'No se pudo preparar el formulario.';
                this.caracterizacionInitEstado = 'error';
            }
        },

        // Si el submission se reutilizó, el evento form-completed no se re-emite:
        // se consulta /forms/{id}/load y se marca completo si todas las secciones
        // visibles tienen isAnswered. Best-effort: si falla, el agente puede guardar
        // la última sección para que el formulario emita form-completed.
        async verificarCaracterizacionCompleta(submissionId) {
            try {
                const res = await fetch('/api/v1/forms/' + this.caracterizacionFormId + '/load?submissionId=' + encodeURIComponent(submissionId));
                if (!res.ok) return;
                const data = await res.json();
                const secciones = (data && data.formStructure && data.formStructure.sections) || [];
                const visibles = secciones.filter(function(s) { return s.isVisible; });
                if (visibles.length > 0 && visibles.every(function(s) { return s.isAnswered; })) {
                    this.caracterizacionFormCompleto = true;
                }
            } catch (e) {
                console.warn('[CTM] no se pudo verificar el estado del formulario reutilizado:', e);
            }
        },

        // El <dinamic-form> emite form-completed solo tras guardar una sección con
        // todas las secciones visibles respondidas. Habilita "Guardar entrevista"
        // (junto con hallazgos y plan de acción).
        onCaracterizacionFormCompleted() {
            this.caracterizacionFormCompleto = true;
        },
        onCaracterizacionAnswersUpdated() {
            // Sin efecto: el formulario guarda sus secciones por su cuenta.
        },

        _contarIntentosFallidos(tasks) {
            var n = 0;
            (tasks || []).forEach(function(t) {
                if (!t || t.type !== 'intentar_contacto_3x3' || t.status !== 'Done') return;
                var fd = t.formData;
                if (typeof fd === 'string') {
                    try { fd = JSON.parse(fd); } catch (e) { fd = null; }
                }
                if (fd && fd.contesto === 'no') n += 1;
            });
            return n;
        },

        async _cargarIntentosContacto() {
            this.intentosPrevios = 0;
            var id = this.tarea && this.tarea.emergencyMeasureReferralId;
            if (!id) return;
            this.cargandoIntentos = true;
            try {
                const res = await fetch('/api/v1/case-tasks?referralId=' + encodeURIComponent(id));
                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }
                if (!res.ok) return;
                const data = await res.json();
                this.intentosPrevios = this._contarIntentosFallidos(Array.isArray(data) ? data : []);
            } catch (e) {
                this.intentosPrevios = 0;
            } finally {
                this.cargandoIntentos = false;
            }
        },

        /* ── Init formulario por tipo ────────────────────────────────────── */

        _initForm() {
            switch (this.tarea.type) {
                case 'decidir_via_riesgo_inminente':
                    this.form = { via: '', agentId: '' };
                    this.cargarAgentesMedidas();
                    break;
                case 'intentar_contacto_3x3':
                    this.form = { contesto: '', nota: '', accion: '' };
                    this._cargarIntentosContacto();
                    break;
                case 'entrevista_caracterizacion': {
                    // Flujo de 3 pasos con estado interno en this.form:
                    //   paso:          1 (¿cómo continúa?) | 2 (consentimiento) | 3 (entrevista)
                    //   continua:      '' | 'consentimiento' | 'concluir'
                    //   consentimiento:'' | 'otorga' | 'no_otorga'
                    // Campos de las ramas:
                    //   Rama caracterizado (paso 3 → otorga): hallazgos, medidasRequeridas[],
                    //     planAccion, enlaceCruzRoja, enlaceDeclaracion.
                    //   Rama concluir (paso 1 → concluir): motivoConclusion, observacionConcluir.
                    //   Rama sin consentimiento (paso 2 → no_otorga): observacionSinConsentimiento.
                    // medidasRequeridas se precarga con las medidas de la remisión
                    // (tarea.referralMeasures) y queda editable (toggle sobre el catálogo).
                    var pre = (this.tarea && Array.isArray(this.tarea.referralMeasures))
                        ? this.tarea.referralMeasures.slice() : [];
                    this.form = {
                        paso: 1,
                        continua: '',
                        consentimiento: '',
                        // Rama caracterizado
                        hallazgos: '',
                        medidasRequeridas: pre,
                        planAccion: '',
                        enlaceCruzRoja: '',
                        enlaceDeclaracion: '',
                        // Rama concluir gestión
                        motivoConclusion: '',
                        observacionConcluir: '',
                        // Rama sin consentimiento
                        observacionSinConsentimiento: '',
                    };
                    // Carga los catálogos del formulario dinámico (enums + ubicación)
                    // una sola vez; el <dinamic-form> del paso 3 los consume vía
                    // state_options_path. No bloquea el render del modal.
                    // Cada apertura reinicia el estado del submission (el backend lo
                    // reutiliza si la tarea ya lo tiene enlazado).
                    this.caracterizacionSubmissionId = '';
                    this.caracterizacionInitEstado = '';
                    this.caracterizacionInitError = '';
                    this.caracterizacionFormCompleto = false;
                    this.cargarCaracterizacionFormState();
                    break;
                }
                case 'revisar_caracterizacion':
                    // resultado: '' | 'aprobado' | 'devuelto'; observaciones solo
                    // aplica cuando se devuelve. El resumen (medidas reales +
                    // placeholders) se arma en computeds a partir de this.tarea.
                    this.form = { resultado: '', observaciones: '' };
                    break;
                case 'revisar_trazabilidad':
                    // resultado: '' | 'preaprobado' | 'devuelto'; observaciones solo
                    // aplica cuando se devuelve. El resumen (medidas reales +
                    // placeholder del plan) se arma en computeds a partir de this.tarea.
                    this.form = { resultado: '', observaciones: '' };
                    break;
                case 'cargar_cotizaciones': {
                    // form.cotizaciones: objeto { '<label de medida>': '<valor>' } con
                    // un campo por cada medida de la remisión (o las 2 por defecto si
                    // no hay). Cada valor inicia vacío para forzar el diligenciamiento.
                    var cotizaciones = {};
                    (this.cotizacionMedidas || []).forEach(function(label) {
                        cotizaciones[label] = '';
                    });
                    this.form = { cotizaciones: cotizaciones };
                    break;
                }
                case 'aprobar_cotizaciones':
                    // resultado: '' | 'aprobado' | 'no_cumple'; observaciones solo
                    // aplica cuando no cumple (devolución). El panel de resumen con los
                    // montos cargados por Cruz Roja y el Total se arman en computeds a
                    // partir de this.tarea.cotizacionesCargadas.
                    this.form = { resultado: '', observaciones: '' };
                    break;
                case 'gestionar_territorio':
                    // fechaGestion: string ISO yyyy-mm-dd (input type=date);
                    // acciones: texto libre. Ambos obligatorios (se valida en el front).
                    this.form = { fechaGestion: '', acciones: '' };
                    break;
                case 'marcar_otorgadas':
                    // fechaOtorgamiento: string ISO yyyy-mm-dd (input type=date);
                    // comentarios: texto libre. Ambos obligatorios (se valida en el
                    // front). Los chips de medidas aprobadas se arman en el computed
                    // medidasAprobadasChips a partir de this.tarea.referralMeasures.
                    this.form = { fechaOtorgamiento: '', comentarios: '' };
                    break;
                case 'escalar_cruz_roja':
                    this.form = {};
                    break;
                case 'entregar_medida_inmediata':
                    this.form = { comentario: '', medidas: [] };
                    this.cargarMedidasRemitidas();
                    break;
                case 'confirmar_lugar_seguro':
                    this.form = {};
                    break;
                case 'revision_equipo_tecnico':
                    this.form = { concepto: '', resultado: '' };
                    this.cargarCasoRemision();
                    this.cargarChat();
                    break;
                case 'ajustar_remision_devuelta':
                    this.form = { ajuste: '', descartar: false };
                    this.cargarChat();
                    break;
                case 'asignar_agente_medidas':
                    this.form = { agentId: '' };
                    this.cargarAgentesMedidas();
                    break;
                case 'gestion_llamada':
                    this.form = {
                        departamentoId: null,
                        ciudadId:       null,
                        municipioId:    null,
                        entidadId:      null,
                        entidadNombre:  '',
                        funcionario:    '',
                        descripcion:    '',
                        generaOficio:   null,
                        asunto:         '',
                        rutaKofax:      '',
                    };
                    console.log('[CTM] E01 gestion_llamada — form inicializado | tarea:', this.tarea.id);
                    break;
                case 'proyectar_oficio':
                    this.form = {
                        departamentoId: null,
                        ciudadId:       null,
                        municipioId:    null,
                        // Precarga el sector de la barrera vinculada cuando coincide con el
                        // catálogo del buscador — el usuario puede cambiarlo si necesita otro.
                        sector:         this.sectoresModal.some(s => s.value === this.tarea.barrierSector)
                                          ? this.tarea.barrierSector : '',
                        entidadId:      null,
                        entidadNombre:  '',
                        funcionario:    '',
                        asunto:         '',
                        rutaKofax:      '',
                        prioridad:      'normal',
                    };
                    break;
                case 'comite_caso':
                    this.form = {
                        // Análisis general de la sesión (obligatorio, aplica a toda la mesa técnica)
                        analisisMesaTecnica:          '',
                        decisiones:                   [],
                        // Activar Enlace
                        solicitudEnlace:              '',
                        entidadDestino:               '',   // 'misma' | 'otra'
                        enlaceSeleccionado:           '',   // icode del enlace, si hay >1
                        // Otra entidad (solo si entidadDestino === 'otra')
                        oeDepartamentoId:             '',
                        oeCiudadId:                   '',
                        oeMunicipioId:                '',
                        oeSector:                     '',
                        oeEntidadId:                  '',
                        oeCanalActivacion:            '',
                        // Mecanismo
                        nivelMecanismo:               null,
                        quienLlenaMecanismo:          '',   // 'equipo_tecnico' | 'agente'
                        // Oficio: quién gestiona el oficio una vez generado
                        gestionaOficio:               '',   // 'equipo_tecnico' | 'agente'
                    };
                    // Cargar enlaces del departamento del caso (para decidir dropdown)
                    this._cargarEnlacesDepartamento();
                    break;
                case 'Corregir oficio':
                    this.form = {};
                    console.log('[CTM] E01 Corregir oficio — form inicializado | tarea:', this.tarea.id);
                    this._cargarOficioVinculado();
                    break;
                case 'justificar_remision':
                    this.form = { respuesta: '' };
                    break;
                case 'solicitud_entidad':
                    this.form = { items: [], seleccion: {} };
                    this._cargarObligacionesPendientes();
                    break;
                case 'respuesta_entidad':
                    this.form = { radicado: '', ampliacion: '', estadoInstitucional: '' };
                    break;
                case 'articulacion_territorial': {
                    const d = new Date();
                    const mm = String(d.getMonth() + 1).padStart(2, '0');
                    const dd = String(d.getDate()).padStart(2, '0');
                    this.form = {
                        fechaContacto: d.getFullYear() + '-' + mm + '-' + dd,
                        compromiso: '',
                        requiereSeguimiento: null,
                        fechaSeguimiento: '',
                        sugerenciasRemision: '',
                    };
                    break;
                }
                case 'gestion_mecanismo':
                    this.form = {
                        abordadoEnSesion: '',
                        fechaSesion: '',
                        compromisos: '',
                        fechaProxima: '',
                        detalleProxima: '',
                        reporteGestion: '',
                        remisionResponsables: '',
                        asistenciasFecha: '',
                        asistenciasTexto: '',
                    };
                    break;
            }
        },

        /* ── Location cascade ────────────────────────────────────────────── */

        onSelectDepartamento() {
            const deptId    = this.form.departamentoId;
            this.ciudades   = deptId ? this.allCiudades.filter(c => String(c.departmentId) === String(deptId)) : [];
            this.form.ciudadId    = null;
            this.form.municipioId = null;
            this.form.entidadId   = null;
            this.form.entidadNombre = '';
            this.municipios = [];
            this.entidades  = []; this.entidadBusqueda = ""; this.entidadDropdownOpen = false;
            console.log('[CTM] E02 dept seleccionado:', deptId, '| ciudades disponibles:', this.ciudades.length);
        },

        async onSelectCiudad() {
            const ciudadId  = this.form.ciudadId;
            this.form.municipioId = null;
            this.form.entidadId   = null;
            this.form.entidadNombre = '';
            this.municipios = [];
            this.entidades  = []; this.entidadBusqueda = ""; this.entidadDropdownOpen = false;
            if (!ciudadId) return;
            console.log('[CTM] E03 ciudad seleccionada:', ciudadId, '— cargando municipios...');
            try {
                const res = await fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(ciudadId));
                if (res.ok) {
                    const data = await res.json();
                    this.municipios = Array.isArray(data) ? data : (data.items || []);
                    console.log('[CTM] E03 municipios cargados:', this.municipios.length);
                }
            } catch (e) {
                console.warn('[case-task-modal] Error cargando municipios:', e);
            }
        },

        async onSelectMunicipio() {
            this.form.entidadId   = null;
            this.form.entidadNombre = '';
            this.entidades  = []; this.entidadBusqueda = ""; this.entidadDropdownOpen = false;
            await this._buscarEntidadesProyectar();
        },

        // Cuando selecciona sector (solo formulario proyectar_oficio) — el sector filtra
        // las entidades del municipio elegido, igual que en case-entities.js.
        async onSelectSector() {
            this.form.entidadId   = null;
            this.form.entidadNombre = '';
            this.entidades  = []; this.entidadBusqueda = ""; this.entidadDropdownOpen = false;
            await this._buscarEntidadesProyectar();
        },

        async _buscarEntidadesProyectar() {
            const townCode = this.form.municipioId;
            if (!townCode) return;
            // proyectar_oficio tiene su propio campo Sector (form.sector, siempre presente
            // como propiedad aunque esté vacío) y lo exige antes de listar entidades —
            // igual que el buscador de case-entities.js. gestion_llamada no tiene campo
            // propio: conserva el comportamiento previo, cargando igual con el filtro
            // automático e invisible de tarea.barrierSector si existe, sin bloquear.
            const tieneCampoSector = Object.prototype.hasOwnProperty.call(this.form, 'sector');
            const sector = tieneCampoSector ? this.form.sector : ((this.tarea && this.tarea.barrierSector) || '');
            if (tieneCampoSector && !sector) return;
            const url = '/api/v1/entity-branches?town_code=' + encodeURIComponent(townCode) + (sector ? '&sector=' + encodeURIComponent(sector) : '');
            console.log('[CTM] E04 municipio=' + townCode + ' sector=' + sector + ' — cargando entidades...');
            try {
                const res = await fetch(url);
                if (res.ok) {
                    const data = await res.json();
                    this.entidades = Array.isArray(data) ? data : (data.items || []);
                    console.log('[CTM] E04 entidades cargadas:', this.entidades.length);
                }
            } catch (e) {
                console.warn('[case-task-modal] Error cargando entidades:', e);
            }
        },

        onSelectEntidad() {
            const id   = this.form.entidadId;
            const name = id && id !== 'otra'
                ? (this.entidades.find(e => String(e.id) === String(id)) || {}).name || '?'
                : 'otra';
            console.log('[CTM] E05 entidad seleccionada: id=' + id + ' · "' + name + '"');
            this.entidadDropdownOpen = false;
        },

        selectEntidad(entidad) {
            this.form.entidadId = entidad.id;
            this.entidadBusqueda = entidad.name;
            this.entidadDropdownOpen = false;
            console.log('[CTM] E05 entidad seleccionada: id=' + entidad.id + ' · "' + entidad.name + '"');
        },

        selectOtraEntidad() {
            this.form.entidadId = 'otra';
            this.entidadBusqueda = '';
            this.entidadDropdownOpen = false;
        },

        onChangeGeneraOficio() {
            console.log('[CTM] E06 generaOficio seleccionado →', this.form.generaOficio);
        },

        /* ── Comité decisions ────────────────────────────────────────────── */

        onChangeDecision(valor) {
            const idx = this.form.decisiones.indexOf(valor);
            if (idx > -1) {
                this.form.decisiones.splice(idx, 1);
                if (valor === 'mecanismo_articulador') {
                    this.form.nivelMecanismo     = null;
                    this.form.quienLlenaMecanismo = '';
                }
                if (valor === 'oficio') {
                    this.form.gestionaOficio = '';
                }
            } else {
                this.form.decisiones.push(valor);
            }
            console.log('[CTM] E07 decisión', idx > -1 ? 'removida' : 'agregada', '→', valor, '| decisiones:', JSON.stringify(this.form.decisiones));
        },

        tieneDecision(valor) {
            return this.form.decisiones && this.form.decisiones.includes(valor);
        },

        // Selección de quién gestiona el oficio. Aunque el requerimiento pide dos
        // checkboxes, las opciones son excluyentes (el oficio lo gestiona el equipo
        // técnico O el agente asignado, nunca ambos): al marcar una se limpia la otra,
        // y volver a marcar la ya seleccionada la desmarca.
        onChangeGestionaOficio(valor) {
            this.form.gestionaOficio = (this.form.gestionaOficio === valor) ? '' : valor;
            console.log('[CTM] E09 gestionaOficio →', this.form.gestionaOficio);
        },

        /* ── Activar Enlace ──────────────────────────────────────────────── */

        // Carga los enlaces territoriales (rol 'en') del departamento del caso.
        // Si hay más de uno, el modal muestra un dropdown para elegir a cuál se dirige.
        async _cargarEnlacesDepartamento() {
            this.enlacesDepartamento = [];
            const caseId = this.tarea && this.tarea.caseId;
            if (!caseId) return;
            this.cargandoEnlaces = true;
            try {
                const res = await fetch('/api/v1/casos/' + encodeURIComponent(caseId) + '/enlaces');
                if (res.ok) {
                    const data = await res.json();
                    this.enlacesDepartamento = Array.isArray(data) ? data : [];
                    // Si hay exactamente uno, preseleccionarlo
                    if (this.enlacesDepartamento.length === 1) {
                        this.form.enlaceSeleccionado = this.enlacesDepartamento[0].icode;
                    }
                }
            } catch (e) {
                console.warn('[CTM] error cargando enlaces del departamento:', e);
            } finally {
                this.cargandoEnlaces = false;
            }
        },

        onCambiaEntidadDestino(valor) {
            this.form.entidadDestino = valor;
            // Limpiar cascada de "otra entidad" al cambiar
            if (valor !== 'otra') {
                this.form.oeDepartamentoId = '';
                this.form.oeCiudadId = '';
                this.form.oeMunicipioId = '';
                this.form.oeSector = '';
                this.form.oeEntidadId = '';
                this.form.oeCanalActivacion = '';
                this.oeCiudades = []; this.oeMunicipios = []; this.oeEntidades = [];
            }
        },

        onOeDepartamento() {
            this.oeCiudades = this.form.oeDepartamentoId
                ? this.allCiudades.filter(c => String(c.departmentId) === String(this.form.oeDepartamentoId)) : [];
            this.form.oeCiudadId = ''; this.form.oeMunicipioId = ''; this.form.oeEntidadId = '';
            this.oeMunicipios = []; this.oeEntidades = [];
        },

        async onOeCiudad() {
            this.form.oeMunicipioId = ''; this.form.oeEntidadId = '';
            this.oeMunicipios = []; this.oeEntidades = [];
            if (!this.form.oeCiudadId) return;
            try {
                const res = await fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(this.form.oeCiudadId));
                if (res.ok) { const d = await res.json(); this.oeMunicipios = Array.isArray(d) ? d : (d.items || []); }
            } catch (e) { console.warn('[CTM] error municipios otra entidad:', e); }
        },

        onOeMunicipio() { this.form.oeEntidadId = ''; this.oeEntidades = []; this._buscarOeEntidades(); },
        onOeSector()    { this.form.oeEntidadId = ''; this.oeEntidades = []; this._buscarOeEntidades(); },

        async _buscarOeEntidades() {
            if (!this.form.oeMunicipioId || !this.form.oeSector) return;
            this.cargandoOeEntidades = true;
            try {
                const res = await fetch('/api/v1/entity-branches?town_code=' + encodeURIComponent(this.form.oeMunicipioId) + '&sector=' + encodeURIComponent(this.form.oeSector));
                this.oeEntidades = res.ok ? (await res.json() || []) : [];
            } catch (e) {
                this.oeEntidades = [];
            } finally {
                this.cargandoOeEntidades = false;
            }
        },

        onChangeNivelMecanismo() {
            console.log('[CTM] E08 nivelMecanismo seleccionado →', this.form.nivelMecanismo);
        },

        /* ── Build formData ──────────────────────────────────────────────── */

        cargarAgentesMedidas() {
            var self = this;
            this.cargandoAgentes = true;
            this.errorAgentes = '';
            fetch('/api/v1/emergency-measure-referrals/agentes', { credentials: 'same-origin' })
                .then(function(res) { return res.json().then(function(body) { return { ok: res.ok, body: body }; }); })
                .then(function(res) {
                    self.agentesMedidas = (res.ok && res.body && res.body.items) ? res.body.items : [];
                    if (!res.ok) self.errorAgentes = 'No se pudieron cargar los agentes.';
                })
                .catch(function() {
                    self.agentesMedidas = [];
                    self.errorAgentes = 'No se pudieron cargar los agentes.';
                })
                .finally(function() { self.cargandoAgentes = false; });
        },

        etiquetaMedida(code) {
            var labels = {
                alojamiento: 'Alojamiento',
                transporte: 'Transporte',
                alimentacion: 'Alimentación',
                vestuario: 'Vestuario',
                aseo: 'Aseo',
                apoyo_psico: 'Apoyo psicosocial',
                exequial: 'Apoyo exequial',
                accion_simbolica: 'Acción simbólica',
                otras_me: 'Otras ME',
                otras: 'Otras medidas de emergencia',
            };
            return labels[code] || code;
        },

        cargarMedidasRemitidas() {
            var self = this;
            var id = this.tarea && this.tarea.emergencyMeasureReferralId;
            this.medidasRemitidas = [];
            this.errorMedidas = '';
            if (!id) {
                this.errorMedidas = 'No se encontró la remisión de la tarea.';
                return;
            }
            this.cargandoMedidas = true;
            fetch('/api/v1/emergency-measure-referrals/' + encodeURIComponent(id), { credentials: 'same-origin' })
                .then(function(res) { return res.json().then(function(body) { return { ok: res.ok, body: body }; }); })
                .then(function(res) {
                    var codes = (res.ok && res.body && res.body.measures) ? res.body.measures : [];
                    self.medidasRemitidas = codes.map(function(code) {
                        return { code: code, label: self.etiquetaMedida(code) };
                    });
                    if (!res.ok) self.errorMedidas = 'No se pudieron cargar las medidas.';
                })
                .catch(function() { self.errorMedidas = 'No se pudieron cargar las medidas.'; })
                .finally(function() { self.cargandoMedidas = false; });
        },

        cargarCasoRemision() {
            var self = this;
            var id = this.tarea && this.tarea.emergencyMeasureReferralId;
            this.casoRevisionNombre = '';
            this.casoRevisionCodigo = '';
            if (!id) return;
            this.cargandoCasoRevision = true;
            fetch('/api/v1/emergency-measure-referrals/' + encodeURIComponent(id), { credentials: 'same-origin' })
                .then(function(res) { return res.json().then(function(body) { return { ok: res.ok, body: body }; }); })
                .then(function(res) {
                    if (!res.ok || !res.body) return;
                    self.casoRevisionNombre = res.body.caseName || '';
                    self.casoRevisionCodigo = res.body.caseICode || '';
                })
                .catch(function() {})
                .finally(function() { self.cargandoCasoRevision = false; });
        },

        cargarChat() {
            var self = this;
            var id = this.tarea && this.tarea.id;
            this.hiloChat = [];
            if (!id) return;
            this.cargandoChat = true;
            fetch('/api/v1/case-tasks/' + encodeURIComponent(id) + '/chat', { credentials: 'same-origin' })
                .then(function(res) { return res.json().then(function(body) { return { ok: res.ok, body: body }; }); })
                .then(function(res) {
                    self.hiloChat = (res.ok && res.body && Array.isArray(res.body.items)) ? res.body.items : [];
                })
                .catch(function() { self.hiloChat = []; })
                .finally(function() { self.cargandoChat = false; });
        },

        medidaSeleccionada(code) {
            return Array.isArray(this.form.medidas) && this.form.medidas.indexOf(code) !== -1;
        },

        toggleMedida(code) {
            if (!Array.isArray(this.form.medidas)) this.form.medidas = [];
            if (this.medidaSeleccionada(code)) {
                this.form.medidas = this.form.medidas.filter(function(item) { return item !== code; });
                return;
            }
            this.form.medidas.push(code);
        },

        // ── Chips de "Medidas requeridas" de la entrevista de caracterización ──
        // Trabajan sobre this.form.medidasRequeridas (labels). Multi-selección
        // toggle: las precargadas de la remisión vienen marcadas y son editables.
        medidaCaracterizacionSeleccionada(label) {
            return Array.isArray(this.form.medidasRequeridas) && this.form.medidasRequeridas.indexOf(label) !== -1;
        },
        toggleMedidaCaracterizacion(label) {
            if (!Array.isArray(this.form.medidasRequeridas)) this.form.medidasRequeridas = [];
            if (this.medidaCaracterizacionSeleccionada(label)) {
                this.form.medidasRequeridas = this.form.medidasRequeridas.filter(function(item) { return item !== label; });
                return;
            }
            this.form.medidasRequeridas.push(label);
        },
        // Navegación de pasos del modal de entrevista (PASO 1 → 2 → 3).
        caracterizacionElegirContinua(valor) {
            this.form.continua = valor;
            // Resetea sub-estado al cambiar de rama en el paso 1.
            if (valor === 'concluir') {
                this.form.consentimiento = '';
            }
        },
        caracterizacionElegirConsentimiento(valor) {
            this.form.consentimiento = valor;
            // Al otorgar consentimiento (paso 3) se prepara el submission precargado.
            if (valor === 'otorga' && !this.caracterizacionSubmissionId) {
                this.iniciarCaracterizacionForm();
            }
        },

        // formatCOP: formatea un número como moneda COP con separador de miles y
        // símbolo $ (ej. 3232 → "$3.232"). Punto como separador de miles (es-CO).
        // Si el valor no es numérico devuelve "$0".
        formatCOP(valor) {
            var num = Number(valor);
            if (!isFinite(num)) num = 0;
            var entero = Math.round(num);
            var signo = entero < 0 ? '-' : '';
            var str = String(Math.abs(entero));
            var partes = '';
            while (str.length > 3) {
                partes = '.' + str.slice(-3) + partes;
                str = str.slice(0, -3);
            }
            return signo + '$' + str + partes;
        },

        _buildFormData() {
            const f = this.form;
            const entidadNombreLabel = (f.entidadId && f.entidadId !== 'otra')
                ? (this.entidades.find(e => String(e.id) === String(f.entidadId)) || {}).name || ''
                : f.entidadNombre;
            // Resuelve el label legible de un <select> a partir de su lista de opciones {value,label}
            const lookupLabel = (lista, valor) =>
                (lista.find(o => String(o.value) === String(valor)) || {}).label || '';

            switch (this.tarea.type) {
                case 'decidir_via_riesgo_inminente':
                    return { via: f.via, agentId: f.agentId };
                case 'entrevista_caracterizacion':
                    // Payload según la rama activa (ver FLUJO del modal). El backend
                    // bifurca por 'resultado'.
                    if (f.continua === 'concluir') {
                        return {
                            resultado: 'concluido_sin_me',
                            motivo: f.motivoConclusion,
                            observacion: String(f.observacionConcluir || '').trim(),
                        };
                    }
                    if (f.continua === 'consentimiento' && f.consentimiento === 'no_otorga') {
                        return {
                            resultado: 'sin_consentimiento',
                            observacion: String(f.observacionSinConsentimiento || '').trim(),
                        };
                    }
                    // Rama caracterizado (otorga consentimiento). Llaves de form_data
                    // consistentes para que FASE 2 y el modal de revisar las consuman.
                    return {
                        resultado: 'caracterizado',
                        consentimiento: 'otorgado',
                        submissionId: this.caracterizacionSubmissionId || '',
                        hallazgos: String(f.hallazgos || '').trim(),
                        medidasRequeridas: (f.medidasRequeridas || []).slice(),
                        planAccion: String(f.planAccion || '').trim(),
                        enlaceCruzRoja: String(f.enlaceCruzRoja || '').trim(),
                        enlaceDeclaracion: String(f.enlaceDeclaracion || '').trim(),
                    };
                case 'revisar_caracterizacion':
                    if (f.resultado === 'devuelto') {
                        return { resultado: 'devuelto', observaciones: String(f.observaciones || '').trim() };
                    }
                    return { resultado: 'aprobado' };
                case 'revisar_trazabilidad':
                    if (f.resultado === 'devuelto') {
                        return { resultado: 'devuelto', observaciones: String(f.observaciones || '').trim() };
                    }
                    // Camino feliz: 'preaprobado' (el backend lo acepta como
                    // sinónimo de 'aprobado' para disparar la transición a cotizar).
                    return { resultado: 'preaprobado' };
                case 'cargar_cotizaciones': {
                    // Arma { cotizaciones: { '<medida>': <valor numérico>, ... } }
                    // con un número por cada medida mostrada. El backend guarda este
                    // JSON en el form_data de la tarea para que el futuro modal de
                    // aprobar (rol sc) pueda leerlo.
                    var valores = (f && f.cotizaciones) || {};
                    var cotizaciones = {};
                    (this.cotizacionMedidas || []).forEach(function(label) {
                        var num = Number(valores[label]);
                        cotizaciones[label] = isFinite(num) ? num : 0;
                    });
                    return { cotizaciones: cotizaciones };
                }
                case 'aprobar_cotizaciones':
                    // No cumple → devolución automática: el backend exige observaciones
                    // y trata cualquier resultado distinto de 'aprobado' como devolución.
                    // Se envía resultado:'devuelto' para que la tarea cerrada luzca el
                    // chip ámbar "Devuelto" (el front lo detecta por resultado==='devuelto').
                    if (f.resultado === 'no_cumple') {
                        return { resultado: 'devuelto', observaciones: String(f.observaciones || '').trim() };
                    }
                    // Camino feliz: avanza a aprobadas y crea gestionar_territorio (rol en).
                    return { resultado: 'aprobado' };
                case 'gestionar_territorio':
                    // Payload persistido en el form_data de la tarea al cerrarla Done.
                    return {
                        fechaGestion: String(f.fechaGestion || '').trim(),
                        acciones: String(f.acciones || '').trim(),
                    };
                case 'marcar_otorgadas':
                    // Payload persistido en el form_data de la tarea al cerrarla Done.
                    return {
                        fechaOtorgamiento: String(f.fechaOtorgamiento || '').trim(),
                        comentarios: String(f.comentarios || '').trim(),
                    };
                case 'escalar_cruz_roja':
                case 'confirmar_lugar_seguro':
                    return { confirmado: true };
                case 'revision_equipo_tecnico':
                    return {
                        concepto: String(f.concepto || '').trim(),
                        resultado: f.resultado,
                    };
                case 'ajustar_remision_devuelta':
                    return {
                        ajuste: f.descartar ? '' : String(f.ajuste || '').trim(),
                        descartar: !!f.descartar,
                    };
                case 'asignar_agente_medidas':
                    return { agentId: f.agentId };
                case 'intentar_contacto_3x3': {
                    var contacto = {
                        contesto: f.contesto,
                        nota: String(f.nota || '').trim(),
                    };
                    if (f.contesto === 'no' && this.intentoContacto === 9 && f.accion === 'cerrar') {
                        contacto.cerrar = true;
                    }
                    if (f.contesto === 'no' && this.intentoContacto >= 10) {
                        contacto.cerrar = true;
                    }
                    return contacto;
                }
                case 'entregar_medida_inmediata':
                    return {
                        medidas: (f.medidas || []).slice(),
                        comentario: String(f.comentario || '').trim(),
                    };
                case 'gestion_llamada':
                    return {
                        departamentoId:     f.departamentoId,
                        departamentoNombre: lookupLabel(this.departamentos, f.departamentoId),
                        ciudadId:           f.ciudadId,
                        ciudadNombre:       lookupLabel(this.ciudades, f.ciudadId),
                        municipioId:        f.municipioId,
                        municipioNombre:    lookupLabel(this.municipios, f.municipioId),
                        entidadId:          f.entidadId !== 'otra' ? Number(f.entidadId) : null,
                        entidadNombre:      entidadNombreLabel,
                        funcionario:        f.funcionario,
                        descripcion:        f.descripcion || null,
                        generaOficio:       !!f.generaOficio,
                        asunto:             f.generaOficio ? f.asunto   : null,
                        rutaKofax:          f.generaOficio ? f.rutaKofax : null,
                    };
                case 'proyectar_oficio':
                    return {
                        departamentoId:     f.departamentoId,
                        departamentoNombre: lookupLabel(this.departamentos, f.departamentoId),
                        ciudadId:           f.ciudadId,
                        ciudadNombre:       lookupLabel(this.ciudades, f.ciudadId),
                        municipioId:        f.municipioId,
                        municipioNombre:    lookupLabel(this.municipios, f.municipioId),
                        entidadId:          f.entidadId !== 'otra' ? Number(f.entidadId) : null,
                        entidadNombre:      entidadNombreLabel,
                        funcionario:        f.funcionario || null,
                        asunto:             f.asunto || null,
                        rutaKofax:          f.rutaKofax || null,
                        prioridad:          f.prioridad || 'normal',
                    };
                case 'comite_caso': {
                    const DECISION_LABELS = {
                        activar_enlace:         'Activar Enlace',
                        oficio:                 'Generar Oficio',
                        recomendaciones_agente: 'Recomendaciones al Agente',
                        mecanismo_articulador:  'Mecanismo Articulador',
                    };
                    const NIVEL_LABELS = {
                        municipal:     'Municipal',
                        departamental: 'Departamental',
                        nacional:      'Nacional',
                    };
                    const QUIEN_LABELS = { equipo_tecnico: 'Equipo Técnico', agente: 'Agente de Seguimiento' };
                    const GESTIONA_OFICIO_LABELS = { equipo_tecnico: 'Equipo Técnico', agente: 'Agente asignado al caso' };
                    const activó = this.tieneDecision('activar_enlace');
                    const esOtra = activó && f.entidadDestino === 'otra';
                    const enlaceObj = this.enlacesDepartamento.find(e => e.icode === f.enlaceSeleccionado);
                    return {
                        // Análisis general de la sesión (aplica a toda la mesa técnica)
                        analisisMesaTecnica:          (f.analisisMesaTecnica || '').trim() || null,
                        decisiones:                   f.decisiones,
                        decisionesTexto:              f.decisiones.map(d => DECISION_LABELS[d] || d),
                        // Activar Enlace
                        solicitudEnlace:              activó ? (f.solicitudEnlace || '').trim() || null : null,
                        entidadDestino:               activó ? f.entidadDestino : null,        // 'misma' | 'otra'
                        enlaceId:                     activó ? (f.enlaceSeleccionado || null) : null,
                        enlaceNombre:                 activó && enlaceObj ? enlaceObj.fullName : null,
                        // Otra entidad (solo si aplica)
                        otraEntidadId:                esOtra ? Number(f.oeEntidadId) : null,
                        otraEntidadNombre:            esOtra ? ((this.oeEntidades.find(e => String(e.id) === String(f.oeEntidadId)) || {}).name || null) : null,
                        otraEntidadSector:            esOtra ? f.oeSector : null,
                        otraEntidadMunicipioId:       esOtra ? f.oeMunicipioId : null,
                        otraEntidadCanal:             esOtra ? (f.oeCanalActivacion || null) : null,
                        // Mecanismo
                        nivelMecanismo:               this.tieneDecision('mecanismo_articulador')  ? f.nivelMecanismo                       : null,
                        nivelMecanismoTexto:          this.tieneDecision('mecanismo_articulador')  ? (NIVEL_LABELS[f.nivelMecanismo] || null) : null,
                        quienLlenaMecanismo:          this.tieneDecision('mecanismo_articulador')  ? (f.quienLlenaMecanismo || null) : null,
                        quienLlenaMecanismoTexto:     this.tieneDecision('mecanismo_articulador')  ? (QUIEN_LABELS[f.quienLlenaMecanismo] || null) : null,
                        // Oficio: quién gestiona el oficio generado
                        gestionaOficio:               this.tieneDecision('oficio')  ? (f.gestionaOficio || null) : null,
                        gestionaOficioTexto:          this.tieneDecision('oficio')  ? (GESTIONA_OFICIO_LABELS[f.gestionaOficio] || null) : null,
                    };
                }
                case 'Corregir oficio':
                    return {};
                case 'justificar_remision':
                    return { respuesta: f.respuesta.trim() };
                case 'solicitud_entidad': {
                    const solicitudes = (f.items || []).filter(function(item) {
                        return f.seleccion[item.id] && f.seleccion[item.id].selected;
                    }).map(function(item) {
                        return {
                            obligationId: item.id,
                            label: item.label,
                            descripcion: String(f.seleccion[item.id].descripcion || '').trim(),
                        };
                    });
                    return {
                        entityName: this.entityName || (this.tarea && this.tarea.entityName) || '',
                        solicitudes: solicitudes,
                    };
                }
                case 'respuesta_entidad': {
                    const prev = this.tareaFormData;
                    const ESTADOS = { iniciado: 'Iniciado', en_proceso: 'En proceso', concluido: 'Concluido' };
                    return {
                        solicitudTaskId: prev.solicitudTaskId || null,
                        obligationId: prev.obligationId || null,
                        label: prev.label || '',
                        descripcionSolicitud: prev.descripcionSolicitud || '',
                        radicado: (f.radicado || '').trim() || null,
                        ampliacion: (f.ampliacion || '').trim() || null,
                        estadoInstitucional: f.estadoInstitucional,
                        estadoInstitucionalTexto: ESTADOS[f.estadoInstitucional] || f.estadoInstitucional,
                    };
                }
                case 'articulacion_territorial':
                    return {
                        fechaContacto: f.fechaContacto,
                        compromiso: String(f.compromiso || '').trim(),
                        requiereSeguimiento: !!f.requiereSeguimiento,
                        fechaSeguimiento: f.requiereSeguimiento ? (f.fechaSeguimiento || '') : '',
                        sugerenciasRemision: String(f.sugerenciasRemision || '').trim() || null,
                        entityBranchId: (this.articulacionCtx && this.articulacionCtx.entityBranchId) || null,
                    };
                case 'gestion_mecanismo': {
                    const SESION_LABELS = {
                        si_ordinaria: 'Sí, sesión ordinaria',
                        si_extraordinaria: 'Sí, sesión extraordinaria',
                        no: 'No',
                    };
                    const scope = (this.mecanismoCtx && this.mecanismoCtx.scope) || '';
                    const esSi = f.abordadoEnSesion === 'si_ordinaria' || f.abordadoEnSesion === 'si_extraordinaria';
                    const esNo = f.abordadoEnSesion === 'no';
                    return {
                        abordadoEnSesion: scope === 'municipal' ? (f.abordadoEnSesion || '') : '',
                        abordadoEnSesionTexto: scope === 'municipal' ? (SESION_LABELS[f.abordadoEnSesion] || '') : '',
                        fechaSesion: scope === 'municipal' && esSi ? (f.fechaSesion || '') : '',
                        compromisos: scope === 'municipal' && esSi ? String(f.compromisos || '').trim() : '',
                        fechaProxima: (scope === 'departamental' || (scope === 'municipal' && esNo))
                            ? (f.fechaProxima || '') : '',
                        detalleProxima: (scope === 'departamental' || (scope === 'municipal' && esNo))
                            ? String(f.detalleProxima || '').trim() : '',
                        reporteGestion: String(f.reporteGestion || '').trim(),
                        remisionResponsables: (scope === 'nacional')
                            ? ''
                            : String(f.remisionResponsables || '').trim(),
                        asistenciasFecha: scope === 'nacional' ? (f.asistenciasFecha || '') : '',
                        asistenciasTexto: scope === 'nacional' ? String(f.asistenciasTexto || '').trim() : '',
                    };
                }
            }
        },

        /* ── Submit ──────────────────────────────────────────────────────── */

        async confirmar() {
            if (!this.formularioValido || this.guardando) return;
            this.guardando = true;
            this.saveError = null;

            const formData = this._buildFormData();
            console.log('[CTM] E09 confirmar — tarea:', this.tarea.id, '| type:', this.tarea.type);
            console.log('[CTM] E09 payload formData:', JSON.stringify(formData));

            try {
                let res;
                if (this.tarea.type === 'articulacion_territorial') {
                    const barrierId = (this.articulacionCtx && this.articulacionCtx.barrierId) || this.tarea.barrierId;
                    const payload = { barrierId: barrierId, formData: formData };
                    if (this.tarea.id) payload.taskId = this.tarea.id;
                    res = await fetch('/api/v1/case-tasks/articulacion-territorial', {
                        method:  'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body:    JSON.stringify(payload),
                    });
                } else if (this.tarea.type === 'gestion_mecanismo') {
                    const payload = {
                        mechanismId: (this.mecanismoCtx && this.mecanismoCtx.mechanismId) || '',
                        formData: formData,
                    };
                    res = await fetch('/api/v1/case-tasks/gestion-mecanismo', {
                        method:  'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body:    JSON.stringify(payload),
                    });
                } else {
                    const payload = { userId: this.currentUserId, formData };
                    res = await fetch('/api/v1/case-tasks/' + this.tarea.id + '/complete', {
                        method:  'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body:    JSON.stringify(payload),
                    });
                }

                console.log('[CTM] E09 respuesta HTTP:', res.status);

                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }

                if (res.status === 422) {
                    const data = await res.json().catch(() => ({}));
                    console.warn('[CTM] E09 error 422:', data);
                    this.saveError = data.error || 'Transición no permitida.';
                    return;
                }

                if (!res.ok) {
                    const data = await res.json().catch(() => ({}));
                    console.warn('[CTM] E09 error', res.status, ':', data);
                    this.saveError = data.error || 'No se pudo completar la tarea. Intenta de nuevo.';
                    return;
                }

                const tareaActualizada = await res.json();
                console.log('[CTM] E09 tarea completada:', JSON.stringify(tareaActualizada));
                this.$emit('completed', tareaActualizada);
                this.completado = true;
                setTimeout(() => this.cancelar(), 1500);

            } catch (e) {
                console.error('[CTM] E09 error de red:', e);
                this.saveError = 'Error de red. Intenta de nuevo.';
            } finally {
                this.guardando = false;
            }
        },
    },

    template: `
<div>
<div v-if="visible" class="ctm-overlay" @click.self="cancelar">
<div class="ctm-modal" :class="{ 'ctm-articulacion': tarea && tarea.type === 'articulacion_territorial', 'ctm-caracterizacion-grande': tarea && tarea.type === 'entrevista_caracterizacion' && form.continua === 'consentimiento' && form.consentimiento === 'otorga' }">

    <!-- Header -->
    <div class="ctm-header">
        <div class="ctm-header-left">
            <div class="ctm-header-icon">
                <i :class="'fas ' + iconoTipo"></i>
            </div>
            <div>
                <p class="ctm-header-title">\${ titulo }</p>
                <p v-if="tarea && tarea.type === 'decidir_via_riesgo_inminente'" class="ctm-header-sub">Elige una vía y asigna un agente de medidas.</p>
                <p v-else-if="tarea && tarea.type === 'intentar_contacto_3x3'" class="ctm-header-sub">Protocolo 3×3 — registra el resultado del intento.</p>
                <p v-else-if="tarea && tarea.type === 'entrevista_caracterizacion'" class="ctm-header-sub">Entrevista de caracterización del Agente de Medidas.</p>
                <p v-else-if="tarea && tarea.type === 'articulacion_territorial'" class="ctm-header-sub">\${ articulacionSubtitulo }</p>
                <p v-else-if="tarea && tarea.type === 'entregar_medida_inmediata'" class="ctm-header-sub">Riesgo inminente: entrega sin esperar caracterización completa ni declaración jurada.</p>
                <p v-else-if="tarea && tarea.type === 'revision_equipo_tecnico'" class="ctm-header-sub">\${ cargandoCasoRevision ? 'Cargando caso…' : ('Caso ' + (casoRevisionNombre || 'sin nombre') + (casoRevisionCodigo ? ' · ' + casoRevisionCodigo : '')) }</p>
                <p v-else-if="tarea && tarea.type === 'revisar_caracterizacion'" class="ctm-header-sub">Verifica que la información y los documentos cargados estén completos, consistentes y con la calidad requerida.</p>
                <p v-else-if="tarea && tarea.type === 'revisar_trazabilidad'" class="ctm-header-sub">Verifica que la información de beneficiarios y la gestión del caso sea suficiente para continuar con el otorgamiento.</p>
                <p v-else-if="tarea && tarea.type === 'cargar_cotizaciones'" class="ctm-header-sub">Registra el valor cotizado para cada medida, de acuerdo con la Resolución 860 de 2024.</p>
                <p v-else-if="tarea && tarea.type === 'aprobar_cotizaciones'" class="ctm-header-sub">Verifica que los valores cotizados por Cruz Roja se ajusten a la Resolución 860 de 2024.</p>
                <p v-else-if="tarea && tarea.type === 'gestionar_territorio'" class="ctm-header-sub">Registra las acciones realizadas en campo para ejecutar la medida.</p>
                <p v-else-if="tarea && tarea.type === 'marcar_otorgadas'" class="ctm-header-sub">Confirma que las medidas fueron entregadas a la víctima.</p>
                <p v-else-if="tarea && tarea.type === 'asignar_agente_medidas'" class="ctm-header-sub">1 remisión seleccionada. Elige el agente destino.</p>
                <p v-else-if="tarea && tarea.id && tarea.type !== 'solicitud_entidad' && tarea.type !== 'respuesta_entidad' && tarea.type !== 'decidir_via_riesgo_inminente' && tarea.type !== 'intentar_contacto_3x3' && tarea.type !== 'escalar_cruz_roja' && tarea.type !== 'entregar_medida_inmediata' && tarea.type !== 'confirmar_lugar_seguro' && tarea.type !== 'revision_equipo_tecnico' && tarea.type !== 'asignar_agente_medidas' && tarea.type !== 'ajustar_remision_devuelta' && tarea.type !== 'entrevista_caracterizacion' && tarea.type !== 'revisar_caracterizacion' && tarea.type !== 'revisar_trazabilidad' && tarea.type !== 'cargar_cotizaciones' && tarea.type !== 'aprobar_cotizaciones' && tarea.type !== 'gestionar_territorio' && tarea.type !== 'marcar_otorgadas'" class="ctm-header-sub">ID \${ tarea.id.substring(0,8) }…</p>
                <p v-else-if="tarea && tarea.type === 'solicitud_entidad'" class="ctm-header-sub">Selecciona los ítems pendientes del checklist entidad sobre los que deseas hacer una solicitud y describe cada uno.</p>
            </div>
        </div>
        <button class="ctm-close-btn" @click="cancelar" :disabled="guardando">✕</button>
    </div>

    <!-- Loader / success overlay -->
    <div v-if="guardando || completado" style="position:absolute;inset:0;background:rgba(255,255,255,.8);border-radius:16px;display:flex;align-items:center;justify-content:center;z-index:10;flex-direction:column;gap:12px">
        <template v-if="completado">
            <div style="width:52px;height:52px;border-radius:50%;background:#dcfce7;display:flex;align-items:center;justify-content:center">
                <svg viewBox="0 0 24 24" fill="none" style="width:28px;height:28px">
                    <path d="M5 13l4 4L19 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
            </div>
            <span style="font-size:14px;font-weight:600;color:#15803d">¡Tarea completada!</span>
        </template>
        <template v-else>
            <svg style="width:32px;height:32px;animation:ctm-spin .7s linear infinite" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
                <path d="M12 2a10 10 0 0 1 10 10" stroke="#1d4ed8" stroke-width="3" stroke-linecap="round"/>
            </svg>
            <span style="font-size:13px;font-weight:600;color:#1d4ed8">Guardando...</span>
        </template>
    </div>

    <!-- Body -->
    <div class="ctm-body">

        <!-- Cargando tarea -->
        <div v-if="cargandoTarea" class="ctm-loading">
            <svg class="ctm-spinner" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
                <path d="M12 2a10 10 0 0 1 10 10" stroke="#1d4ed8" stroke-width="3" stroke-linecap="round"/>
            </svg>
            Cargando tarea...
        </div>

        <!-- Error al cargar -->
        <div v-else-if="errorTarea" class="ctm-error-banner">
            <i class="fas fa-exclamation-circle"></i> \${ errorTarea }
        </div>

        <!-- Formulario gestion_llamada -->
        <template v-else-if="tarea && tarea.type === 'gestion_llamada'">

            <p class="ctm-section-title">Ubicación de la entidad</p>

            <div class="ctm-field">
                <label class="ctm-label">Departamento <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.departamentoId" @change="onSelectDepartamento">
                    <option :value="null" disabled>Selecciona departamento</option>
                    <option v-for="d in departamentos" :key="d.value" :value="d.value">
                        \${ d.label }
                    </option>
                </select>
            </div>

            <div class="ctm-row-2">
                <div class="ctm-field">
                    <label class="ctm-label">Ciudad <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.ciudadId" @change="onSelectCiudad" :disabled="!form.departamentoId">
                        <option :value="null" disabled>Selecciona ciudad</option>
                        <option v-for="c in ciudades" :key="c.value" :value="c.value">
                            \${ c.label }
                        </option>
                    </select>
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Municipio <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.municipioId" @change="onSelectMunicipio" :disabled="!form.ciudadId">
                        <option :value="null" disabled>Selecciona municipio</option>
                        <option v-for="m in municipios" :key="m.value" :value="m.value">
                            \${ m.label }
                        </option>
                    </select>
                </div>
            </div>

            <div class="ctm-field" style="position:relative">
                <label class="ctm-label">Sede / Entidad <span class="req">*</span></label>
                <input class="ctm-input" v-model="entidadBusqueda" :disabled="!form.municipioId" placeholder="Buscar entidad..." @focus="entidadDropdownOpen = true" @input="entidadDropdownOpen = true; form.entidadId = null" autocomplete="off"/>
                <div v-if="entidadDropdownOpen && form.municipioId" style="position:absolute;z-index:100;top:100%;left:0;right:0;max-height:200px;overflow-y:auto;background:#fff;border:1px solid #d1d5db;border-radius:8px;box-shadow:0 4px 12px rgba(0,0,0,.1);margin-top:4px">
                    <div v-for="e in entidadesFiltradas" :key="e.id" @mousedown.prevent="selectEntidad(e)" style="padding:8px 12px;font-size:.82rem;cursor:pointer;border-bottom:1px solid #f3f4f6" :style="form.entidadId === e.id ? 'background:#ede9fe;color:#5106A7;font-weight:600' : ''" @mouseover="$event.target.style.background='#f9fafb'" @mouseout="$event.target.style.background=form.entidadId === e.id ? '#ede9fe' : ''">
                        \${ e.name }
                    </div>
                    <div @mousedown.prevent="selectOtraEntidad()" style="padding:8px 12px;font-size:.82rem;cursor:pointer;color:#6b7280;font-style:italic;border-top:1px solid #e5e7eb">
                        Otra entidad…
                    </div>
                    <div v-if="entidadesFiltradas.length === 0" style="padding:12px;font-size:.8rem;color:#9ca3af;text-align:center">
                        No se encontraron resultados
                    </div>
                </div>
            </div>

            <div v-if="form.entidadId === 'otra'" class="ctm-field">
                <label class="ctm-label">Nombre de la entidad <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.entidadNombre" placeholder="Escribe el nombre de la entidad" maxlength="5000"/>
            </div>

            <p class="ctm-section-title">Contacto</p>

            <div class="ctm-field">
                <label class="ctm-label">Funcionario contactado <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.funcionario" placeholder="Nombre del funcionario" maxlength="5000"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Descripción / Notas</label>
                <textarea class="ctm-textarea" v-model="form.descripcion" rows="3" placeholder="Resumen de la gestión..." maxlength="5000"></textarea>
            </div>

            <p class="ctm-section-title">Oficio</p>

            <div class="ctm-field">
                <label class="ctm-label">¿Genera oficio? <span class="req">*</span></label>
                <div class="ctm-radio-group">
                    <label class="ctm-radio-item">
                        <input type="radio" v-model="form.generaOficio" :value="true" @change="onChangeGeneraOficio"/>
                        Sí
                    </label>
                    <label class="ctm-radio-item">
                        <input type="radio" v-model="form.generaOficio" :value="false" @change="onChangeGeneraOficio"/>
                        No
                    </label>
                </div>
            </div>

            <template v-if="form.generaOficio === true">
                <div class="ctm-field">
                    <label class="ctm-label">Asunto del oficio <span class="req">*</span></label>
                    <input class="ctm-input" v-model="form.asunto" placeholder="Asunto…" maxlength="5000"/>
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Ruta Kofax / URL del archivo <span class="req">*</span></label>
                    <input class="ctm-input" v-model="form.rutaKofax" placeholder="\\\\servidor\\ruta\\archivo.pdf" maxlength="5000"/>
                </div>
            </template>

        </template>

        <!-- Formulario proyectar_oficio -->
        <template v-else-if="tarea && tarea.type === 'proyectar_oficio'">

            <p class="ctm-section-title">Ubicación de la entidad</p>

            <div class="ctm-field">
                <label class="ctm-label">Departamento <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.departamentoId" @change="onSelectDepartamento">
                    <option :value="null" disabled>Selecciona departamento</option>
                    <option v-for="d in departamentos" :key="d.value" :value="d.value">
                        \${ d.label }
                    </option>
                </select>
            </div>

            <div class="ctm-row-2">
                <div class="ctm-field">
                    <label class="ctm-label">Ciudad <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.ciudadId" @change="onSelectCiudad" :disabled="!form.departamentoId">
                        <option :value="null" disabled>Selecciona ciudad</option>
                        <option v-for="c in ciudades" :key="c.value" :value="c.value">
                            \${ c.label }
                        </option>
                    </select>
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Municipio <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.municipioId" @change="onSelectMunicipio" :disabled="!form.ciudadId">
                        <option :value="null" disabled>Selecciona municipio</option>
                        <option v-for="m in municipios" :key="m.value" :value="m.value">
                            \${ m.label }
                        </option>
                    </select>
                </div>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Sector <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.sector" @change="onSelectSector" :disabled="!form.municipioId">
                    <option value="">Selecciona sector</option>
                    <option v-for="s in sectoresModal" :key="s.value" :value="s.value">
                        \${ s.label }
                    </option>
                </select>
            </div>

            <div class="ctm-field" style="position:relative">
                <label class="ctm-label">Sede / Entidad <span class="req">*</span></label>
                <input class="ctm-input" v-model="entidadBusqueda" :disabled="!form.municipioId || !form.sector" placeholder="Buscar entidad..." @focus="entidadDropdownOpen = true" @input="entidadDropdownOpen = true; form.entidadId = null" autocomplete="off"/>
                <div v-if="entidadDropdownOpen && form.municipioId" style="position:absolute;z-index:100;top:100%;left:0;right:0;max-height:200px;overflow-y:auto;background:#fff;border:1px solid #d1d5db;border-radius:8px;box-shadow:0 4px 12px rgba(0,0,0,.1);margin-top:4px">
                    <div v-for="e in entidadesFiltradas" :key="e.id" @mousedown.prevent="selectEntidad(e)" style="padding:8px 12px;font-size:.82rem;cursor:pointer;border-bottom:1px solid #f3f4f6" :style="form.entidadId === e.id ? 'background:#ede9fe;color:#5106A7;font-weight:600' : ''" @mouseover="$event.target.style.background='#f9fafb'" @mouseout="$event.target.style.background=form.entidadId === e.id ? '#ede9fe' : ''">
                        \${ e.name }
                    </div>
                    <div @mousedown.prevent="selectOtraEntidad()" style="padding:8px 12px;font-size:.82rem;cursor:pointer;color:#6b7280;font-style:italic;border-top:1px solid #e5e7eb">
                        Otra entidad…
                    </div>
                    <div v-if="entidadesFiltradas.length === 0" style="padding:12px;font-size:.8rem;color:#9ca3af;text-align:center">
                        No se encontraron resultados
                    </div>
                </div>
            </div>

            <div v-if="form.entidadId === 'otra'" class="ctm-field">
                <label class="ctm-label">Nombre de la entidad</label>
                <input class="ctm-input" v-model="form.entidadNombre" placeholder="Escribe el nombre de la entidad" maxlength="5000"/>
            </div>

            <p class="ctm-section-title">Datos del oficio</p>

            <div class="ctm-field">
                <label class="ctm-label">Dependencia del funcionario</label>
                <input class="ctm-input" v-model="form.funcionario" placeholder="Ej. Departamento de Protección Social" maxlength="5000"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Asunto</label>
                <input class="ctm-input" v-model="form.asunto" placeholder="Asunto del oficio…" maxlength="5000"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Ruta Kofax / URL del archivo</label>
                <input class="ctm-input" v-model="form.rutaKofax" placeholder="\\\\servidor\\ruta\\archivo.pdf" maxlength="5000"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Prioridad del oficio</label>
                <select class="ctm-select" v-model="form.prioridad">
                    <option value="normal">Normal</option>
                    <option value="alta">Alta</option>
                </select>
            </div>

        </template>

        <!-- Formulario comite_caso -->
        <template v-else-if="tarea && tarea.type === 'comite_caso'">

            <!-- Motivo por el que se remitió el caso a mesa técnica (justificación al abrirla,
                 o contexto de la barrera cuando nace de "Alerta por barreras") -->
            <div v-if="tarea.description" style="margin-bottom:16px">
                <p class="ctm-section-title" style="color:#5106A7">
                    <i class="fas fa-info-circle"></i> Motivo de la mesa técnica
                </p>
                <div style="padding:12px 14px;background:#f5f3ff;border:1.5px solid #e9d5ff;border-radius:8px;color:#374151;font-size:13px;line-height:1.5;white-space:pre-wrap">
                    \${ tarea.description }
                </div>
            </div>

            <!-- Análisis de la mesa técnica (obligatorio, aplica a toda la sesión) -->
            <div class="ctm-field">
                <label class="ctm-label">Análisis de la mesa técnica <span class="req">*</span></label>
                <textarea class="ctm-textarea" v-model="form.analisisMesaTecnica" rows="4" placeholder="Resumen de lo que se determinó en la mesa técnica y por qué se toman estas decisiones…" maxlength="5000"></textarea>
            </div>

            <p class="ctm-section-title">Decisiones tomadas <span style="color:#dc2626">*</span></p>
            <p style="font-size:12px;color:#6b7280;margin:0 0 12px">Selecciona al menos una</p>

            <div class="ctm-checkbox-group">

                <label :class="['ctm-checkbox-item', tieneDecision('activar_enlace') ? 'checked' : '']"
                    @click="onChangeDecision('activar_enlace')">
                    <div class="ctm-checkbox-box">
                        <i class="fas fa-check ctm-checkbox-check"></i>
                    </div>
                    <span class="ctm-checkbox-text">Activar Enlace</span>
                </label>

                <label :class="['ctm-checkbox-item', tieneDecision('oficio') ? 'checked' : '']"
                    @click="onChangeDecision('oficio')">
                    <div class="ctm-checkbox-box">
                        <i class="fas fa-check ctm-checkbox-check"></i>
                    </div>
                    <span class="ctm-checkbox-text">Generar Oficio</span>
                </label>

                <label :class="['ctm-checkbox-item', tieneDecision('recomendaciones_agente') ? 'checked' : '']"
                    @click="onChangeDecision('recomendaciones_agente')">
                    <div class="ctm-checkbox-box">
                        <i class="fas fa-check ctm-checkbox-check"></i>
                    </div>
                    <span class="ctm-checkbox-text">Recomendaciones al Agente</span>
                </label>

                <label :class="['ctm-checkbox-item', tieneDecision('mecanismo_articulador') ? 'checked' : '']"
                    @click="onChangeDecision('mecanismo_articulador')">
                    <div class="ctm-checkbox-box">
                        <i class="fas fa-check ctm-checkbox-check"></i>
                    </div>
                    <span class="ctm-checkbox-text">Mecanismo Articulador</span>
                </label>

            </div>

            <!-- Campos condicionales -->

            <template v-if="tieneDecision('activar_enlace')">
                <div style="margin-top:16px;background:#f5f3ff;border:1px solid #e9d5ff;border-radius:10px;padding:14px 16px">
                    <p class="ctm-section-title" style="color:#5106A7;margin:0 0 10px">Activar Enlace Territorial</p>

                    <div class="ctm-field">
                        <label class="ctm-label">Solicitud para el enlace <span class="req">*</span></label>
                        <textarea class="ctm-textarea" v-model="form.solicitudEnlace" rows="3" placeholder="Explica la tarea que debe realizar el enlace territorial…" maxlength="5000"></textarea>
                    </div>

                    <!-- Dropdown de enlace solo si hay más de uno -->
                    <div class="ctm-field" v-if="enlacesDepartamento.length > 1">
                        <label class="ctm-label">¿A qué enlace se dirige? <span class="req">*</span></label>
                        <select class="ctm-select" v-model="form.enlaceSeleccionado">
                            <option value="">Selecciona un enlace…</option>
                            <option v-for="en in enlacesDepartamento" :key="en.icode" :value="en.icode">\${ en.fullName }</option>
                        </select>
                    </div>
                    <div class="ctm-field" v-else-if="enlacesDepartamento.length === 1">
                        <label class="ctm-label">Enlace</label>
                        <div style="font-size:13px;color:#374151;padding:6px 0">\${ enlacesDepartamento[0].fullName }</div>
                    </div>

                    <label class="ctm-label" style="margin-top:6px">¿A qué entidad debe ir la enlace? <span class="req">*</span></label>
                    <div style="display:flex;gap:10px;margin-top:6px">
                        <button type="button"
                            @click="onCambiaEntidadDestino('misma')"
                            :style="form.entidadDestino === 'misma'
                                ? 'flex:1;border:2px solid #5106A7;background:#ede9fe;color:#5106A7;border-radius:8px;padding:10px;font-size:.8rem;font-weight:700;cursor:pointer'
                                : 'flex:1;border:2px solid #e5e7eb;background:#fff;color:#374151;border-radius:8px;padding:10px;font-size:.8rem;font-weight:600;cursor:pointer'">
                            Misma entidad de la barrera<span v-if="tarea.barrierEntityName"> (\${ tarea.barrierEntityName })</span>
                        </button>
                        <button type="button"
                            @click="onCambiaEntidadDestino('otra')"
                            :style="form.entidadDestino === 'otra'
                                ? 'flex:1;border:2px solid #5106A7;background:#ede9fe;color:#5106A7;border-radius:8px;padding:10px;font-size:.8rem;font-weight:700;cursor:pointer'
                                : 'flex:1;border:2px solid #e5e7eb;background:#fff;color:#374151;border-radius:8px;padding:10px;font-size:.8rem;font-weight:600;cursor:pointer'">
                            Otra entidad
                        </button>
                    </div>

                    <!-- Formulario de otra entidad -->
                    <template v-if="form.entidadDestino === 'otra'">
                        <div class="ctm-field" style="margin-top:12px">
                            <label class="ctm-label">Sector <span class="req">*</span></label>
                            <select class="ctm-select" v-model="form.oeSector" @change="onOeSector">
                                <option value="">Selecciona un sector…</option>
                                <option v-for="s in sectoresModal" :key="s.value" :value="s.value">\${ s.label }</option>
                            </select>
                        </div>
                        <div class="ctm-field">
                            <label class="ctm-label">Departamento <span class="req">*</span></label>
                            <select class="ctm-select" v-model="form.oeDepartamentoId" @change="onOeDepartamento">
                                <option value="">Selecciona…</option>
                                <option v-for="d in departamentos" :key="d.value" :value="d.value">\${ d.label }</option>
                            </select>
                        </div>
                        <div class="ctm-field">
                            <label class="ctm-label">Ciudad <span class="req">*</span></label>
                            <select class="ctm-select" v-model="form.oeCiudadId" @change="onOeCiudad" :disabled="!form.oeDepartamentoId">
                                <option value="">Selecciona…</option>
                                <option v-for="c in oeCiudades" :key="c.value" :value="c.value">\${ c.label }</option>
                            </select>
                        </div>
                        <div class="ctm-field">
                            <label class="ctm-label">Municipio <span class="req">*</span></label>
                            <select class="ctm-select" v-model="form.oeMunicipioId" @change="onOeMunicipio" :disabled="!form.oeCiudadId">
                                <option value="">Selecciona…</option>
                                <option v-for="m in oeMunicipios" :key="m.value" :value="m.value">\${ m.label }</option>
                            </select>
                        </div>
                        <div class="ctm-field" v-if="cargandoOeEntidades">
                            <span style="font-size:12px;color:#6b7280">Buscando entidades…</span>
                        </div>
                        <div class="ctm-field" v-else-if="oeEntidades.length > 0">
                            <label class="ctm-label">Entidad <span class="req">*</span></label>
                            <select class="ctm-select" v-model="form.oeEntidadId">
                                <option value="">Selecciona la entidad…</option>
                                <option v-for="opt in oeEntidades" :key="opt.id" :value="opt.id">\${ opt.name }</option>
                            </select>
                        </div>
                        <div class="ctm-field">
                            <label class="ctm-label">Canal de activación</label>
                            <select class="ctm-select" v-model="form.oeCanalActivacion">
                                <option value="">Selecciona…</option>
                                <option value="notificacion">Notificación / Oficio</option>
                                <option value="llamada">Llamada</option>
                                <option value="correo">Correo electrónico</option>
                            </select>
                        </div>
                    </template>
                </div>
            </template>

            <template v-if="tieneDecision('mecanismo_articulador')">
                <p class="ctm-section-title" style="margin-top:16px">Mecanismo Articulador</p>
                <div class="ctm-field">
                    <label class="ctm-label">Nivel del mecanismo <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.nivelMecanismo" @change="onChangeNivelMecanismo">
                        <option :value="null" disabled>Selecciona nivel</option>
                        <option value="municipal">Municipal</option>
                        <option value="departamental">Departamental</option>
                        <option value="nacional">Nacional</option>
                    </select>
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">¿Quién llena el formulario del mecanismo? <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.quienLlenaMecanismo">
                        <option value="">Selecciona…</option>
                        <option value="equipo_tecnico">Equipo Técnico</option>
                        <option value="agente">Agente de Seguimiento del caso</option>
                    </select>
                </div>
            </template>

            <!-- Gestión del oficio: se despliega al elegir la decisión "Generar Oficio".
                 Define quién queda a cargo de gestionar el oficio generado. Va al final
                 del formulario, como pide el requerimiento. -->
            <template v-if="tieneDecision('oficio')">
                <p class="ctm-section-title" style="margin-top:16px">Gestión del oficio</p>
                <p style="font-size:.78rem;color:#6b7280;margin:0 0 8px">¿Quién gestiona el oficio una vez generado?</p>
                <div class="ctm-checkbox-group">
                    <label :class="['ctm-checkbox-item', form.gestionaOficio === 'equipo_tecnico' ? 'checked' : '']"
                        @click="onChangeGestionaOficio('equipo_tecnico')">
                        <div class="ctm-checkbox-box">
                            <i class="fas fa-check ctm-checkbox-check"></i>
                        </div>
                        <span class="ctm-checkbox-text">Lo gestiona el Equipo Técnico</span>
                    </label>

                    <label :class="['ctm-checkbox-item', form.gestionaOficio === 'agente' ? 'checked' : '']"
                        @click="onChangeGestionaOficio('agente')">
                        <div class="ctm-checkbox-box">
                            <i class="fas fa-check ctm-checkbox-check"></i>
                        </div>
                        <span class="ctm-checkbox-text">Lo gestiona el Agente asignado al caso</span>
                    </label>
                </div>
            </template>

        </template>

        <!-- Formulario justificar_remision -->
        <template v-else-if="tarea && tarea.type === 'justificar_remision'">

            <div style="margin-bottom:16px">
                <p class="ctm-section-title" style="color:#b91c1c">
                    <i class="fas fa-exclamation-circle"></i> Motivo de la devolución
                </p>
                <div style="padding:12px 14px;background:#fef2f2;border:1.5px solid #fecaca;border-radius:8px;color:#b91c1c;font-size:13px;line-height:1.5">
                    \${ tarea.description }
                </div>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Tu respuesta <span class="req">*</span></label>
                <textarea class="ctm-textarea" v-model="form.respuesta" rows="4" placeholder="Explica o justifica la remisión para que vuelva a revisión…" maxlength="5000"></textarea>
            </div>

        </template>

        <!-- Formulario Corregir oficio -->
        <template v-else-if="tarea && tarea.type === 'Corregir oficio'">

            <!-- Cargando datos del oficio -->
            <div v-if="cargandoOficio" class="ctm-loading">
                <svg class="ctm-spinner" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
                    <path d="M12 2a10 10 0 0 1 10 10" stroke="#1d4ed8" stroke-width="3" stroke-linecap="round"/>
                </svg>
                Cargando datos del oficio...
            </div>

            <template v-else-if="oficioVinculado">

                <!-- Razón de corrección -->
                <div v-if="oficioVinculado.reasonCorrection" style="margin-bottom:16px">
                    <p class="ctm-section-title" style="color:#b91c1c">
                        <i class="fas fa-exclamation-circle"></i> Razón de corrección
                    </p>
                    <div style="padding:12px 14px;background:#fef2f2;border:1.5px solid #fecaca;border-radius:8px;color:#b91c1c;font-size:13px;line-height:1.5">
                        \${ oficioVinculado.reasonCorrection }
                    </div>
                </div>

                <!-- Ruta Kofax -->
                <div v-if="oficioVinculado.urlKofax" class="ctm-field">
                    <label class="ctm-label">Oficio en Kofax</label>
                    <div style="display:flex;align-items:center;gap:10px;padding:10px 12px;background:#f9fafb;border:1.5px solid #e5e7eb;border-radius:8px;font-size:13px">
                        <i class="fas fa-file-alt" style="color:#6b7280;flex-shrink:0"></i>
                        <span style="word-break:break-all;color:#374151">\${ oficioVinculado.urlKofax }</span>
                    </div>
                </div>

                <p style="font-size:13px;color:#6b7280;margin:16px 0 0;line-height:1.5">
                    Revisa el oficio en el Kofax, aplica las correcciones indicadas y presiona
                    <strong>"Marcar como corregido"</strong> para enviarlo nuevamente a revisión.
                </p>

            </template>

            <div v-else style="color:#6b7280;font-size:13px;margin:8px 0">
                No se pudo cargar la información del oficio.
            </div>

        </template>

        <!-- Formulario solicitud_entidad -->
        <template v-else-if="tarea && tarea.type === 'solicitud_entidad'">

            <div v-if="cargandoObligaciones" class="ctm-loading">
                <svg class="ctm-spinner" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
                    <path d="M12 2a10 10 0 0 1 10 10" stroke="#1d4ed8" stroke-width="3" stroke-linecap="round"/>
                </svg>
                Cargando obligaciones...
            </div>

            <div v-else-if="errorObligaciones" class="ctm-error-banner">
                <i class="fas fa-exclamation-circle"></i> \${ errorObligaciones }
            </div>

            <template v-else>
                <p class="ctm-section-title">Checklist entidad pendientes</p>
                <div v-if="!form.items || form.items.length === 0" style="color:#6b7280;font-size:13px;margin:8px 0">
                    No hay obligaciones pendientes
                </div>
                <div v-else class="ctm-ob-list">
                    <div v-for="item in form.items" :key="item.id">
                        <div class="ctm-ob-item" :class="{ selected: form.seleccion[item.id] && form.seleccion[item.id].selected }" @click="toggleObligacion(item.id)">
                            <div class="ctm-ob-item-row">
                                <span class="ctm-ob-item-label">\${ item.label }</span>
                                <i class="fas fa-check ctm-ob-check"></i>
                            </div>
                        </div>
                        <div v-if="form.seleccion[item.id] && form.seleccion[item.id].selected" class="ctm-ob-desc" @click.stop>
                            <textarea class="ctm-textarea" v-model="form.seleccion[item.id].descripcion" rows="3" maxlength="5000" placeholder="Describe la solicitud para este ítem del checklist entidad..."></textarea>
                        </div>
                    </div>
                </div>
            </template>

        </template>

        <!-- Formulario respuesta_entidad -->
        <template v-else-if="tarea && tarea.type === 'respuesta_entidad'">
            <div class="ctm-context-box">
                <div class="ctm-field">
                    <label class="ctm-label">Obligación</label>
                    <div style="font-size:13px;color:#111827;line-height:1.45">\${ tareaFormData.label || '—' }</div>
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Solicitud</label>
                    <div style="font-size:13px;color:#374151;line-height:1.45;white-space:pre-wrap">\${ tareaFormData.descripcionSolicitud || '—' }</div>
                </div>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">En línea con las acciones adelantadas, ¿cuenta con un Nº de radicado, NUNC, MP, etc.?</label>
                <input class="ctm-input" v-model="form.radicado" placeholder="Ej: NUNC 123456789012" maxlength="5000"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">¿Desea ampliar información sobre la gestión y atención del caso?</label>
                <textarea class="ctm-textarea" v-model="form.ampliacion" rows="4" maxlength="5000" placeholder="Describe cómo se atendió esta acción institucional..."></textarea>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Estado general de la Respuesta Institucional <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.estadoInstitucional">
                    <option value="">Selecciona una opción...</option>
                    <option value="iniciado">Iniciado</option>
                    <option value="en_proceso">En proceso</option>
                    <option value="concluido">Concluido</option>
                </select>
            </div>
        </template>

        <!-- Formulario articulacion_territorial -->
        <template v-else-if="tarea && tarea.type === 'articulacion_territorial'">
            <div v-if="articulacionSubtitulo" class="ctm-context-box">
                \${ articulacionSubtitulo }
            </div>
            <!-- Entidad contactada: solo lectura, precargada con la entidad de la barrera -->
            <div class="ctm-field">
                <label class="ctm-label">Entidad contactada</label>
                <input class="ctm-input" type="text" readonly disabled
                    :value="(articulacionCtx && articulacionCtx.entityName) || 'Entidad de la barrera'"
                    style="background:#f3f4f6;color:#6b7280;cursor:not-allowed">
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Fecha de contacto <span class="req">*</span></label>
                <input class="ctm-input" type="date" v-model="form.fechaContacto">
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Compromiso adquirido <span class="req">*</span></label>
                <textarea class="ctm-textarea" rows="4" v-model="form.compromiso"
                    placeholder="Describe el compromiso adquirido con la entidad..."></textarea>
            </div>
            <div class="ctm-field">
                <label class="ctm-label">¿El compromiso requiere seguimiento? <span class="req">*</span></label>
                <div class="ctm-toggle-row">
                    <button type="button" class="ctm-toggle-btn" :class="{ active: form.requiereSeguimiento === true }"
                        @click="form.requiereSeguimiento = true">Sí</button>
                    <button type="button" class="ctm-toggle-btn" :class="{ active: form.requiereSeguimiento === false }"
                        @click="form.requiereSeguimiento = false">No</button>
                </div>
            </div>
            <div v-if="form.requiereSeguimiento" class="ctm-field">
                <label class="ctm-label">Fecha de seguimiento al compromiso <span class="req">*</span></label>
                <input class="ctm-input" type="date" v-model="form.fechaSeguimiento">
                <p class="ctm-hint">Se enviará un recordatorio de seguimiento en esta fecha y se generará una tarea con este vencimiento.</p>
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Sugerencias de remisión interna o externa</label>
                <textarea class="ctm-textarea" rows="3" v-model="form.sugerenciasRemision"
                    placeholder="Opcional"></textarea>
            </div>
        </template>

        <!-- Formulario gestion_mecanismo -->
        <template v-else-if="tarea && tarea.type === 'gestion_mecanismo'">
            <template v-if="mecanismoCtx && mecanismoCtx.scope === 'municipal'">
                <div class="ctm-field">
                    <label class="ctm-label">¿El caso ha sido abordado en sesión del Comité Municipal? <span class="req">*</span></label>
                    <select class="ctm-select" v-model="form.abordadoEnSesion">
                        <option value="">Selecciona una opción...</option>
                        <option value="si_ordinaria">Sí, sesión ordinaria</option>
                        <option value="si_extraordinaria">Sí, sesión extraordinaria</option>
                        <option value="no">No</option>
                    </select>
                    <p class="ctm-hint">Responda sí únicamente si la sesión ya sucedió.</p>
                </div>
                <template v-if="form.abordadoEnSesion === 'si_ordinaria' || form.abordadoEnSesion === 'si_extraordinaria'">
                    <div class="ctm-field">
                        <label class="ctm-label">Fecha de la sesión</label>
                        <input class="ctm-input" type="date" v-model="form.fechaSesion">
                    </div>
                    <div class="ctm-field">
                        <label class="ctm-label">Compromisos de la sesión</label>
                        <textarea class="ctm-textarea" rows="3" maxlength="5000" v-model="form.compromisos"
                            placeholder="Si el caso se abordó en varias sesiones, consigne los compromisos de todas."></textarea>
                    </div>
                </template>
                <template v-else-if="form.abordadoEnSesion === 'no'">
                    <div class="ctm-field">
                        <label class="ctm-label">Fecha de la próxima sesión</label>
                        <input class="ctm-input" type="date" v-model="form.fechaProxima">
                    </div>
                    <div class="ctm-field">
                        <label class="ctm-label">Detalle de la próxima sesión</label>
                        <textarea class="ctm-textarea" rows="3" maxlength="5000" v-model="form.detalleProxima"
                            placeholder="Indique cuándo se citará, según la urgencia del caso."></textarea>
                    </div>
                </template>
            </template>

            <template v-if="mecanismoCtx && mecanismoCtx.scope === 'departamental'">
                <div class="ctm-field">
                    <label class="ctm-label">Fecha de la próxima sesión o mesa técnica <span class="req">*</span></label>
                    <input class="ctm-input" type="date" v-model="form.fechaProxima">
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Detalle de la próxima sesión o mesa técnica <span class="req">*</span></label>
                    <textarea class="ctm-textarea" rows="3" maxlength="5000" v-model="form.detalleProxima"
                        placeholder="Indique cuándo se citará la próxima sesión o mesa técnica."></textarea>
                </div>
            </template>

            <div class="ctm-field">
                <label class="ctm-label">Reporte de gestión <span class="req">*</span></label>
                <textarea class="ctm-textarea" rows="4" maxlength="5000" v-model="form.reporteGestion"
                    placeholder="Acciones adelantadas y, si aplica, enlace a la documentación."></textarea>
                <p class="ctm-hint">Describa las acciones adelantadas. Si hay documentación de soporte, incluya el enlace.</p>
            </div>

            <div v-if="mecanismoCtx && mecanismoCtx.scope !== 'nacional'" class="ctm-field">
                <label class="ctm-label">Remisión a responsables</label>
                <textarea class="ctm-textarea" rows="3" maxlength="5000" v-model="form.remisionResponsables"
                    placeholder="Opcional"></textarea>
                <p class="ctm-hint">Según el Decreto 1710. Este campo se guarda; no dispara un escalamiento automático.</p>
            </div>

            <template v-if="mecanismoCtx && mecanismoCtx.scope === 'nacional'">
                <div class="ctm-field">
                    <label class="ctm-label">Fecha de asistencias técnicas <span class="req">*</span></label>
                    <input class="ctm-input" type="date" v-model="form.asistenciasFecha">
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Asistencias técnicas <span class="req">*</span></label>
                    <textarea class="ctm-textarea" rows="3" maxlength="5000" v-model="form.asistenciasTexto"
                        placeholder="Si no aplica, escriba «No aplica» y explique."></textarea>
                    <p class="ctm-hint">Si no aplica, escriba «No aplica» y la explicación.</p>
                </div>
            </template>
        </template>

        <template v-else-if="tarea && tarea.type === 'decidir_via_riesgo_inminente'">
            <label style="display:flex;gap:12px;align-items:flex-start;border:1px solid #e5e7eb;border-radius:12px;padding:14px 16px;margin-bottom:12px;cursor:pointer">
                <input type="radio" name="decidir-via" value="A" v-model="form.via" style="margin-top:4px">
                <span>
                    <strong style="display:block;color:#111827">Vía A — Asignación inmediata a Agente de Medidas</strong>
                    <span style="display:block;margin-top:4px;color:#4b5563;font-size:.84rem">El agente prioriza contacto y entrevista de caracterización de inmediato.</span>
                </span>
            </label>
            <label style="display:flex;gap:12px;align-items:flex-start;border:1px solid #e5e7eb;border-radius:12px;padding:14px 16px;margin-bottom:16px;cursor:pointer">
                <input type="radio" name="decidir-via" value="B" v-model="form.via" style="margin-top:4px">
                <span>
                    <strong style="display:block;color:#111827">Vía B — Derivar a Supervisión del Convenio</strong>
                    <span style="display:block;margin-top:4px;color:#4b5563;font-size:.84rem">Escala a Cruz Roja para entrega inmediata (hotel/transporte), sin esperar caracterización.</span>
                </span>
            </label>
            <div class="ctm-field">
                <label class="ctm-label">Asignar agente <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.agentId" :disabled="cargandoAgentes">
                    <option value="">\${ cargandoAgentes ? 'Cargando agentes…' : 'Selecciona un agente' }</option>
                    <option v-for="agente in agentesMedidas" :key="agente.id" :value="agente.id">\${ agente.nombre }</option>
                </select>
                <p v-if="errorAgentes" class="ctm-hint">\${ errorAgentes }</p>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'entrevista_caracterizacion'">

            <!-- Cabecera: caso + código + contacto efectivo (degrada si falta dato) -->
            <p class="ctm-car-caso">Caso \${ caracterizacionNombre }<span v-if="caracterizacionCodigo"> · \${ caracterizacionCodigo }</span></p>
            <p class="ctm-car-contacto">
                <i class="fas fa-check-circle"></i>
                <span v-if="caracterizacionFechaContacto">Contacto efectivo confirmado el \${ caracterizacionFechaContacto }</span>
                <span v-else>Contacto efectivo confirmado</span>
            </p>

            <!-- ── PASO 1 — ¿Cómo continúa el caso? ── -->
            <div style="margin-top:18px">
                <p class="ctm-label" style="margin-bottom:10px">¿Cómo continúa el caso?</p>
                <button type="button" class="ctm-car-toggle" :class="{ 'is-on-green': form.continua === 'consentimiento' }" @click="caracterizacionElegirContinua('consentimiento')">Continuar con consentimiento informado</button>
                <button type="button" class="ctm-car-toggle" :class="{ 'is-on-red': form.continua === 'concluir' }" @click="caracterizacionElegirContinua('concluir')">✕ Concluir gestión</button>
            </div>

            <!-- ── RAMA A — Concluir gestión ── -->
            <template v-if="form.continua === 'concluir'">
                <div style="margin-top:16px">
                    <p class="ctm-label" style="margin-bottom:10px">Motivo de conclusión <span class="req">*</span></p>
                    <button type="button" class="ctm-resultado" :class="{ 'is-on': form.motivoConclusion === 'desistimiento' }" @click="form.motivoConclusion = 'desistimiento'">Desistimiento de la ciudadana</button>
                    <button type="button" class="ctm-resultado" :class="{ 'is-on': form.motivoConclusion === 'sin_riesgo' }" @click="form.motivoConclusion = 'sin_riesgo'">No hay riesgo de feminicidio y tiene rutas de atención activas</button>
                    <div class="ctm-field" style="margin-top:12px">
                        <label class="ctm-label">Observación <span class="req">*</span></label>
                        <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.observacionConcluir" placeholder="Describe lo ocurrido en el intento de contacto o por qué el caso no cumple criterios de ME..."></textarea>
                    </div>
                </div>
            </template>

            <!-- ── PASO 2 — Consentimiento (si eligió continuar) ── -->
            <template v-else-if="form.continua === 'consentimiento'">
                <div style="margin-top:16px">
                    <p class="ctm-label" style="margin-bottom:10px">¿La víctima otorga el consentimiento informado para la entrevista de caracterización?</p>
                    <button type="button" class="ctm-car-toggle" :class="{ 'is-on-green': form.consentimiento === 'otorga' }" @click="caracterizacionElegirConsentimiento('otorga')">✓ Otorga consentimiento</button>
                    <button type="button" class="ctm-car-toggle" :class="{ 'is-on-red': form.consentimiento === 'no_otorga' }" @click="caracterizacionElegirConsentimiento('no_otorga')">✕ No otorga consentimiento</button>
                </div>

                <!-- ── RAMA B — No otorga consentimiento ── -->
                <template v-if="form.consentimiento === 'no_otorga'">
                    <div class="ctm-field" style="margin-top:12px">
                        <label class="ctm-label">Observación <span class="req">*</span></label>
                        <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.observacionSinConsentimiento" placeholder="Describe el contexto por el que la víctima no otorga el consentimiento informado..."></textarea>
                    </div>
                </template>

                <!-- ── PASO 3 — Otorga consentimiento: entrevista completa ── -->
                <template v-else-if="form.consentimiento === 'otorga'">
                    <p class="ctm-section-title" style="margin-top:18px">Entrevista de caracterización (consentimiento informado otorgado)</p>

                    <!-- Formulario dinámico de caracterización (FASE 2b-2). El backend
                         crea (o reutiliza) un submission PRECARGADO y EDITABLE desde el
                         Registro de Caso y lo enlaza a la tarea; el componente solo se
                         monta con ese submissionId y con los catálogos ya cargados.
                         Convive con los campos de resultado de abajo. -->
                    <div class="ctm-car-placeholder" style="background:#eff6ff;border-style:solid">
                        <p class="ph-title"><i class="fas fa-clipboard-list"></i> Entrevista de caracterización</p>
                        <p class="ph-body">Diligencia las secciones de la entrevista. Los datos que ya están en SALVIA vienen precargados y puedes editarlos. Guarda cada sección; "Guardar entrevista" se habilita al completar el formulario y los campos de resultado.</p>
                    </div>
                    <div v-if="caracterizacionInitEstado === 'error'" class="ctm-car-placeholder" style="background:#fef2f2;border-style:solid">
                        <p class="ph-title"><i class="fas fa-exclamation-triangle"></i> No se pudo preparar el formulario</p>
                        <p class="ph-body">\${ caracterizacionInitError }</p>
                        <button type="button" class="ctm-btn ctm-btn-cancel" @click="iniciarCaracterizacionForm">Reintentar</button>
                    </div>
                    <div v-else-if="!caracterizacionSubmissionId || !caracterizacionFormStateListo" class="ctm-hint" style="padding:12px 0">
                        <i class="fas fa-circle-notch fa-spin"></i> Preparando el formulario con los datos del caso…
                    </div>
                    <dinamic-form v-else
                        :key="caracterizacionSubmissionId"
                        :form-id="caracterizacionFormId"
                        :submission-id="caracterizacionSubmissionId"
                        :form-state="caracterizacionFormState"
                        @form-completed="onCaracterizacionFormCompleted"
                        @answers-updated="onCaracterizacionAnswersUpdated"></dinamic-form>
                    <p v-if="!caracterizacionFormCompleto" class="ctm-hint">Completa y guarda todas las secciones del formulario para habilitar "Guardar entrevista".</p>

                    <!-- Campos de RESULTADO (sí se implementan en fase 1) -->
                    <div class="ctm-field">
                        <label class="ctm-label">Hallazgos de la entrevista <span class="req">*</span></label>
                        <textarea class="ctm-textarea" rows="3" maxlength="4000" v-model="form.hallazgos" placeholder="Información adicional identificada durante la entrevista sobre la situación de riesgo de feminicidio..."></textarea>
                    </div>

                    <div class="ctm-field">
                        <label class="ctm-label">Medidas requeridas</label>
                        <div style="display:flex;flex-wrap:wrap;gap:8px">
                            <span class="ctm-car-chip" v-for="medida in caracterizacionMedidasCatalogo" :key="medida" :class="{ 'is-on': medidaCaracterizacionSeleccionada(medida) }" @click="toggleMedidaCaracterizacion(medida)">\${ medida }</span>
                        </div>
                        <p class="ctm-hint">Las medidas de la remisión vienen preseleccionadas y son editables. Marca o desmarca según la caracterización.</p>
                    </div>

                    <div class="ctm-field">
                        <label class="ctm-label">Plan de acción <span class="req">*</span></label>
                        <textarea class="ctm-textarea" rows="3" maxlength="4000" v-model="form.planAccion" placeholder="Medidas, atenciones, remisiones y acciones requeridas por los equipos de SALVIA de acuerdo con las necesidades identificadas..."></textarea>
                    </div>

                    <div class="ctm-field">
                        <label class="ctm-label">Enlace — remisión a Cruz Roja</label>
                        <input type="text" class="ctm-input" maxlength="1000" v-model="form.enlaceCruzRoja" placeholder="https://..." />
                    </div>

                    <div class="ctm-field">
                        <label class="ctm-label">Enlace — declaración jurada</label>
                        <input type="text" class="ctm-input" maxlength="1000" v-model="form.enlaceDeclaracion" placeholder="https://..." />
                    </div>
                </template>
            </template>

            <!-- Footer propio de la entrevista (el footer genérico se oculta). Los
                 botones de submit solo aparecen en los estados terminales de cada rama. -->
            <div class="ctm-car-footer">
                <button class="ctm-btn ctm-btn-cancel" @click="cancelar" :disabled="guardando">Cancelar</button>
                <button v-if="form.continua === 'concluir'" class="ctm-btn ctm-btn-confirm is-devolver" @click="confirmar" :disabled="!formularioValido || guardando">
                    <template v-if="guardando"><i class="fas fa-circle-notch fa-spin"></i> Guardando…</template>
                    <template v-else>Registrar conclusión</template>
                </button>
                <button v-else-if="form.continua === 'consentimiento' && form.consentimiento === 'no_otorga'" class="ctm-btn ctm-btn-confirm" @click="confirmar" :disabled="!formularioValido || guardando">
                    <template v-if="guardando"><i class="fas fa-circle-notch fa-spin"></i> Guardando…</template>
                    <template v-else>Registrar</template>
                </button>
                <button v-else-if="form.continua === 'consentimiento' && form.consentimiento === 'otorga'" class="ctm-btn ctm-btn-confirm is-aprobar" @click="confirmar" :disabled="!formularioValido || guardando">
                    <template v-if="guardando"><i class="fas fa-circle-notch fa-spin"></i> Guardando…</template>
                    <template v-else>Guardar entrevista</template>
                </button>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'intentar_contacto_3x3'">
            <p class="ctm-label" style="margin-bottom:10px">¿Contestó?</p>
            <div style="display:flex;gap:10px;margin-bottom:16px">
                <button type="button" class="ctm-btn ctm-btn-confirm" :style="form.contesto === 'si' ? '' : 'background:#fff;color:#166534;border:1px solid #86efac'" @click="form.contesto = 'si'; form.accion = ''">Sí</button>
                <button type="button" class="ctm-btn ctm-btn-cancel" :style="form.contesto === 'no' ? 'background:#b91c1c;border-color:#b91c1c;color:#fff;font-weight:600;box-shadow:0 0 0 3px rgba(185,28,28,.28)' : 'background:#fff;color:#b91c1c;border:1px solid #fecaca'" @click="form.contesto = 'no'">No</button>
            </div>
            <p v-if="cargandoIntentos" class="ctm-hint">Cargando intentos…</p>
            <div v-else-if="form.contesto === 'no' && intentoContacto === 9" style="margin-bottom:16px">
                <p class="ctm-label" style="margin-bottom:10px">El ciclo 3×3 está agotado</p>
                <div style="display:flex;gap:10px;flex-wrap:wrap">
                    <button type="button" class="ctm-btn ctm-btn-confirm" :style="form.accion === 'seguir' ? '' : 'background:#fff;color:#1d4ed8;border:1px solid #bfdbfe'" @click="form.accion = 'seguir'">Seguir intentando</button>
                    <button type="button" class="ctm-btn ctm-btn-cancel" :style="form.accion === 'cerrar' ? 'background:#b91c1c;border-color:#b91c1c;color:#fff;font-weight:600;box-shadow:0 0 0 3px rgba(185,28,28,.28)' : 'background:#fff;color:#b91c1c;border:1px solid #fecaca'" @click="form.accion = 'cerrar'">Cerrar remisión por no contacto</button>
                </div>
            </div>
            <p v-else-if="form.contesto === 'no' && intentoContacto >= 10" class="ctm-hint" style="color:#b91c1c">Este es el intento 10. Al guardar se cierra la remisión por no contacto.</p>
            <div class="ctm-field">
                <label class="ctm-label">Nota</label>
                <input class="ctm-input" type="text" maxlength="500" v-model="form.nota" placeholder="Ej: No contesta, fuera de área...">
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'escalar_cruz_roja'">
            <p style="margin:0;color:#111827;font-size:.95rem">Confirmar Derivar el caso a Cruz Roja para entrega inmediata de la medida</p>
        </template>

        <template v-else-if="tarea && tarea.type === 'confirmar_lugar_seguro'">
            <p style="margin:0;color:#111827;font-size:.95rem">Confirma entrega de medidas</p>
        </template>

        <template v-else-if="tarea && tarea.type === 'ajustar_remision_devuelta'">
            <p style="margin:0 0 14px;color:#374151;font-size:.9rem">\${ textoDevolucion }</p>
            <p v-if="cargandoChat" class="ctm-hint">Cargando historial…</p>
            <div v-else-if="hiloAnterior.length" class="ctm-chat">
                <div v-for="(msg, i) in hiloAnterior" :key="i" class="ctm-chat-msg" :class="msg.lado === 'agente' ? 'ctm-chat-agente' : 'ctm-chat-equipo'">
                    <strong>\${ msg.lado === 'agente' ? 'Agente de seguimiento' : 'Equipo Técnico' }</strong>
                    <p>\${ msg.texto }</p>
                </div>
            </div>
            <div v-if="ultimoMensajeEquipo" class="ctm-obs">
                <strong>Observaciones del Equipo Técnico</strong>
                <p>\${ ultimoMensajeEquipo.texto }</p>
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Ajuste o información complementada</label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.ajuste" :disabled="form.descartar" placeholder="Describe la información adicional o el ajuste realizado para que el Equipo Técnico pueda revisar nuevamente el caso..."></textarea>
            </div>
            <label class="ctm-check-descartar">
                <input type="checkbox" v-model="form.descartar">
                Descartar remisión definitivamente
            </label>
        </template>

        <template v-else-if="tarea && tarea.type === 'revision_equipo_tecnico'">
            <p v-if="cargandoChat" class="ctm-hint">Cargando historial…</p>
            <div v-else-if="hiloChat.length" class="ctm-chat">
                <div v-for="(msg, i) in hiloChat" :key="i" class="ctm-chat-msg" :class="msg.lado === 'agente' ? 'ctm-chat-agente' : 'ctm-chat-equipo'">
                    <strong>\${ msg.lado === 'agente' ? 'Agente de seguimiento' : 'Equipo Técnico' }</strong>
                    <p>\${ msg.texto }</p>
                </div>
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Concepto técnico</label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.concepto" placeholder="Análisis cualificado del caso: criterios identificados, condiciones de riesgo y recomendaciones de la revisión."></textarea>
            </div>
            <p class="ctm-label" style="margin-bottom:10px">Resultado de la revisión</p>
            <button type="button" class="ctm-resultado" :class="{ 'is-on': form.resultado === 'aprobado_para_caracterizacion_me' }" @click="form.resultado = 'aprobado_para_caracterizacion_me'">Aprobado para caracterización ME</button>
            <button type="button" class="ctm-resultado" :class="{ 'is-on': form.resultado === 'devuelto_para_ajuste' }" @click="form.resultado = 'devuelto_para_ajuste'">Devuelto para ajuste</button>
            <button type="button" class="ctm-resultado" :class="{ 'is-on': form.resultado === 'devuelto_para_acompanamiento' }" @click="form.resultado = 'devuelto_para_acompanamiento'">Devuelto para acompañamiento</button>
        </template>

        <template v-else-if="tarea && tarea.type === 'revisar_caracterizacion'">
            <!-- Resumen de la caracterización (solo lectura). "Medidas requeridas"
                 es dato real de la remisión; el resto son placeholders por ahora. -->
            <div class="ctm-context-box">
                <div v-for="fila in resumenRevisarCaracterizacion" :key="fila.label" class="ctm-resumen-row">
                    <span class="ctm-resumen-label">\${ fila.label }:</span>
                    <span class="ctm-resumen-valor">\${ fila.valor }</span>
                </div>
            </div>
            <p class="ctm-label" style="margin-bottom:10px">Resultado de la revisión</p>
            <button type="button" class="ctm-resultado" :class="{ 'is-aprobado': form.resultado === 'aprobado' }" @click="form.resultado = 'aprobado'">✓ Aprobado por ME</button>
            <button type="button" class="ctm-resultado" :class="{ 'is-devuelto': form.resultado === 'devuelto' }" @click="form.resultado = 'devuelto'">✕ Devuelto para ajuste</button>
            <div v-if="form.resultado === 'devuelto'" class="ctm-field" style="margin-top:12px">
                <label class="ctm-label">Observaciones para el Agente de Medidas</label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.observaciones" placeholder="Describe qué información o documento falta ajustar..."></textarea>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'revisar_trazabilidad'">
            <!-- Resumen del caso (solo lectura). "Medidas requeridas" es dato real
                 de la remisión; "Plan de acción" es placeholder por ahora. -->
            <div class="ctm-context-box">
                <div v-for="fila in resumenRevisarTrazabilidad" :key="fila.label" class="ctm-resumen-row">
                    <span class="ctm-resumen-label">\${ fila.label }:</span>
                    <span class="ctm-resumen-valor">\${ fila.valor }</span>
                </div>
            </div>
            <p class="ctm-label" style="margin-bottom:10px">Resultado de la revisión</p>
            <button type="button" class="ctm-resultado" :class="{ 'is-aprobado': form.resultado === 'preaprobado' }" @click="form.resultado = 'preaprobado'">✓ Preaprobado - enviar a Cruz Roja</button>
            <button type="button" class="ctm-resultado" :class="{ 'is-devuelto': form.resultado === 'devuelto' }" @click="form.resultado = 'devuelto'">✕ Devolver a Medidas de Emergencia</button>
            <div v-if="form.resultado === 'devuelto'" class="ctm-field" style="margin-top:12px">
                <label class="ctm-label">Observaciones para Medidas de Emergencia</label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.observaciones" placeholder="Ej: Información de beneficiarios insuficiente, falta nombre/edad/documento..."></textarea>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'cargar_cotizaciones'">
            <!-- Una fila por cada medida de la remisión (tarea.referralMeasures) o,
                 si no hay, las dos por defecto. Cada input pide el valor cotizado;
                 el botón solo se habilita cuando todas tienen un número >= 0. -->
            <p class="ctm-label" style="margin-bottom:10px">Valor cotizado por medida</p>
            <div class="ctm-field" v-for="medida in cotizacionMedidas" :key="medida" style="margin-bottom:12px">
                <label class="ctm-label">\${ medida }</label>
                <input type="number" min="0" step="any" class="ctm-input" v-model="form.cotizaciones[medida]" placeholder="Valor cotizado (COP)" />
            </div>
            <p class="ctm-hint">Cada medida debe tener un valor (se acepta 0). No se admiten campos vacíos.</p>
        </template>

        <template v-else-if="tarea && tarea.type === 'aprobar_cotizaciones'">
            <!-- Panel de resumen (solo lectura): cada medida con el valor cotizado por
                 Cruz Roja (tarea.cotizacionesCargadas), formateado en COP, y una fila
                 final con el Total. Fallback: medidas de la remisión con valor 0. -->
            <p class="ctm-label" style="margin-bottom:10px">Valores cotizados por Cruz Roja</p>
            <div class="ctm-context-box">
                <div v-for="fila in cotizacionesAprobarResumen" :key="fila.label" class="ctm-resumen-row">
                    <span class="ctm-resumen-label">\${ fila.label }:</span>
                    <span class="ctm-resumen-valor">\${ formatCOP(fila.valor) }</span>
                </div>
                <div class="ctm-resumen-row" style="border-top:1px solid #e5e7eb;margin-top:6px;padding-top:6px;font-weight:800">
                    <span class="ctm-resumen-label">Total:</span>
                    <span class="ctm-resumen-valor">\${ formatCOP(totalCotizacionesAprobar) }</span>
                </div>
            </div>
            <p class="ctm-label" style="margin-bottom:10px">Resultado de la revisión</p>
            <button type="button" class="ctm-resultado" :class="{ 'is-aprobado': form.resultado === 'aprobado' }" @click="form.resultado = 'aprobado'">✓ Aprobar cotizaciones</button>
            <button type="button" class="ctm-resultado" :class="{ 'is-devuelto': form.resultado === 'no_cumple' }" @click="form.resultado = 'no_cumple'">✕ No cumple</button>
            <div v-if="form.resultado === 'no_cumple'" class="ctm-field" style="margin-top:12px">
                <label class="ctm-label">Observaciones para Cruz Roja</label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.observaciones" placeholder="Explica qué valores no se ajustan a la Resolución 860 de 2024..."></textarea>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'gestionar_territorio'">
            <!-- Gestión en territorio (rol en). Ambos campos obligatorios: fecha de
                 gestión (input date nativo, valor ISO yyyy-mm-dd) y acciones realizadas. -->
            <div class="ctm-field">
                <label class="ctm-label">Fecha de gestión <span class="req">*</span></label>
                <input type="date" class="ctm-input" v-model="form.fechaGestion" />
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Acciones realizadas <span class="req">*</span></label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.acciones" placeholder="Describe las acciones: coordinación con albergue, traslado, entrega de kits, etc."></textarea>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'marcar_otorgadas'">
            <!-- Marcar medidas como otorgadas (rol cr). Chips de solo lectura con las
                 medidas aprobadas de la remisión (tarea.referralMeasures) + fecha de
                 otorgamiento y comentarios, ambos obligatorios. -->
            <p class="ctm-label" style="margin-bottom:10px">Medidas aprobadas</p>
            <div style="display:flex;flex-wrap:wrap;gap:8px;margin-bottom:16px">
                <span class="ctm-chip" v-for="medida in medidasAprobadasChips" :key="medida">\${ medida }</span>
                <span v-if="medidasAprobadasChips.length === 0" class="ctm-hint">Sin medidas registradas en la remisión.</span>
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Fecha de otorgamiento <span class="req">*</span></label>
                <input type="date" class="ctm-input" v-model="form.fechaOtorgamiento" />
            </div>
            <div class="ctm-field">
                <label class="ctm-label">Comentarios <span class="req">*</span></label>
                <textarea class="ctm-textarea" rows="4" maxlength="2000" v-model="form.comentarios" placeholder="Ej: Medidas entregadas en punto de atención norte. Beneficiaria recibió kit completo..."></textarea>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'asignar_agente_medidas'">
            <div class="ctm-field">
                <label class="ctm-label">Agente de Medidas</label>
                <select class="ctm-select" v-model="form.agentId" :disabled="cargandoAgentes">
                    <option value="">\${ cargandoAgentes ? 'Cargando agentes…' : 'Seleccionar agente…' }</option>
                    <option v-for="agente in agentesMedidas" :key="agente.id" :value="agente.id">\${ agente.nombre }</option>
                </select>
                <p v-if="errorAgentes" class="ctm-hint">\${ errorAgentes }</p>
            </div>
        </template>

        <template v-else-if="tarea && tarea.type === 'entregar_medida_inmediata'">
            <p class="ctm-label" style="margin-bottom:10px">Medidas entregadas</p>
            <p v-if="cargandoMedidas" class="ctm-hint">Cargando medidas…</p>
            <div v-else style="display:flex;flex-wrap:wrap;gap:8px;margin-bottom:16px">
                <button type="button" class="ctm-chip" :class="{ 'is-on': medidaSeleccionada(medida.code) }" v-for="medida in medidasRemitidas" :key="medida.code" @click="toggleMedida(medida.code)" :disabled="guardando">\${ medida.label }</button>
            </div>
            <p v-if="errorMedidas" class="ctm-hint">\${ errorMedidas }</p>
            <p v-else-if="!cargandoMedidas && medidasRemitidas.length === 0" class="ctm-hint">Esta remisión no tiene medidas.</p>
            <div class="ctm-field">
                <label class="ctm-label">Comentarios</label>
                <textarea class="ctm-textarea" rows="3" maxlength="500" v-model="form.comentario" placeholder="Ej: Hotel asignado en zona segura, traslado realizado el mismo día..."></textarea>
            </div>
        </template>

        <!-- Error al guardar -->
        <div v-if="saveError" class="ctm-error-banner">
            <i class="fas fa-exclamation-circle"></i> \${ saveError }
        </div>

    </div>

    <!-- Footer (genérico). La entrevista de caracterización usa su propio footer
         dentro del cuerpo porque sus botones de submit dependen del paso/rama. -->
    <div class="ctm-footer" v-if="!(tarea && tarea.type === 'entrevista_caracterizacion')">
        <button class="ctm-btn ctm-btn-cancel" @click="cancelar" :disabled="guardando">
            Cancelar
        </button>
        <button class="ctm-btn ctm-btn-confirm" @click="confirmar"
            :class="{ 'is-aprobar': tarea && ((tarea.type === 'revisar_caracterizacion' && form.resultado === 'aprobado') || (tarea.type === 'revisar_trazabilidad' && form.resultado === 'preaprobado') || (tarea.type === 'aprobar_cotizaciones' && form.resultado === 'aprobado') || tarea.type === 'cargar_cotizaciones' || tarea.type === 'gestionar_territorio' || tarea.type === 'marcar_otorgadas'), 'is-devolver': tarea && ((tarea.type === 'revisar_caracterizacion' || tarea.type === 'revisar_trazabilidad') && form.resultado === 'devuelto') || (tarea.type === 'aprobar_cotizaciones' && form.resultado === 'no_cumple') }"
            :disabled="!formularioValido || guardando || cargandoTarea">
            <template v-if="guardando">
                <i class="fas fa-circle-notch fa-spin"></i> Guardando…
            </template>
            <template v-else>
                <template v-if="tarea && tarea.type === 'Corregir oficio'">
                    Marcar como corregido
                </template>
                <template v-else-if="tarea && tarea.type === 'solicitud_entidad'">
                    Enviar solicitudes
                </template>
                <template v-else-if="tarea && tarea.type === 'respuesta_entidad'">
                    Enviar respuesta
                </template>
                <template v-else-if="tarea && tarea.type === 'articulacion_territorial'">
                    Guardar
                </template>
                <template v-else-if="tarea && tarea.type === 'gestion_mecanismo'">
                    Guardar
                </template>
                <template v-else-if="tarea && tarea.type === 'decidir_via_riesgo_inminente'">
                    Confirmar decisión
                </template>
                <template v-else-if="tarea && tarea.type === 'intentar_contacto_3x3'">
                    \${ form.contesto === 'no' && (intentoContacto >= 10 || (intentoContacto === 9 && form.accion === 'cerrar')) ? 'Cerrar remisión por no contacto' : 'Guardar' }
                </template>
                <template v-else-if="tarea && tarea.type === 'entrevista_caracterizacion'">
                    Guardar
                </template>
                <template v-else-if="tarea && tarea.type === 'escalar_cruz_roja'">
                    Confirmar
                </template>
                <template v-else-if="tarea && tarea.type === 'entregar_medida_inmediata'">
                    Registrar entrega
                </template>
                <template v-else-if="tarea && tarea.type === 'confirmar_lugar_seguro'">
                    Confirmar
                </template>
                <template v-else-if="tarea && tarea.type === 'revision_equipo_tecnico'">
                    Guardar revisión
                </template>
                <template v-else-if="tarea && tarea.type === 'revisar_caracterizacion'">
                    \${ form.resultado === 'devuelto' ? 'Devolver para ajuste' : 'Aprobar y enviar a Supervisión del Convenio' }
                </template>
                <template v-else-if="tarea && tarea.type === 'revisar_trazabilidad'">
                    \${ form.resultado === 'devuelto' ? 'Devolver a ME' : 'Preaprobar y enviar a Cruz Roja' }
                </template>
                <template v-else-if="tarea && tarea.type === 'cargar_cotizaciones'">
                    Enviar cotizaciones
                </template>
                <template v-else-if="tarea && tarea.type === 'aprobar_cotizaciones'">
                    \${ form.resultado === 'no_cumple' ? 'Devolver a Cruz Roja' : 'Aprobar cotizaciones' }
                </template>
                <template v-else-if="tarea && tarea.type === 'gestionar_territorio'">
                    Guardar gestión
                </template>
                <template v-else-if="tarea && tarea.type === 'marcar_otorgadas'">
                    Marcar como otorgada
                </template>
                <template v-else-if="tarea && tarea.type === 'ajustar_remision_devuelta'">
                    \${ form.descartar ? 'Descartar remisión' : 'Reenviar al Equipo Técnico' }
                </template>
                <template v-else-if="tarea && tarea.type === 'asignar_agente_medidas'">
                    Reasignar
                </template>
                <template v-else>
                    Completar tarea
                </template>
            </template>
        </button>
    </div>

</div>
</div>
</div>
    `,
});
