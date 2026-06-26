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
    `;
    document.head.appendChild(style);
})();

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
            saveError:     null,

            // Ubicaciones
            departamentos: [],
            allCiudades:   [],
            ciudades:      [],
            municipios:    [],
            entidades:     [],

            // Formulario (shape varía según tarea.type)
            form: {},

            // Datos del entity_letter vinculado (para Corregir oficio)
            oficioVinculado: null,
            cargandoOficio:  false,
        };
    },

    computed: {
        titulo() {
            if (!this.tarea) return 'Completar tarea';
            return {
                gestion_llamada:  'Gestión de Llamada',
                proyectar_oficio: 'Proyectar Oficio',
                comite_caso:      'Decisiones del Comité',
                'Corregir oficio': 'Corregir Oficio',
            }[this.tarea.type] || 'Completar tarea';
        },

        iconoTipo() {
            if (!this.tarea) return 'fa-tasks';
            return {
                gestion_llamada:  'fa-phone-alt',
                proyectar_oficio: 'fa-file-signature',
                comite_caso:      'fa-users',
                'Corregir oficio': 'fa-pencil-alt',
            }[this.tarea.type] || 'fa-tasks';
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
                    if (!f.departamentoId || !f.ciudadId || !f.municipioId || !f.entidadId) return false;
                    if (f.entidadId === 'otra' && !f.entidadNombre) return false;
                    if (!f.funcionario || !f.asunto || !f.rutaKofax) return false;
                    return true;
                case 'comite_caso':
                    if (!f.decisiones || f.decisiones.length === 0) return false;
                    if (f.decisiones.includes('mecanismo_articulador') && !f.nivelMecanismo) return false;
                    return true;
                case 'Corregir oficio':
                    // Sin campos obligatorios — solo confirmación
                    return !this.cargandoOficio;
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
            this.entidades     = [];
            this.form          = {};

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

        cancelar() {
            this.visible         = false;
            this.tarea           = null;
            this.form            = {};
            this.saveError       = null;
            this.errorTarea      = null;
            this.ciudades        = [];
            this.municipios      = [];
            this.entidades       = [];
            this.oficioVinculado = null;
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

        /* ── Init formulario por tipo ────────────────────────────────────── */

        _initForm() {
            switch (this.tarea.type) {
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
                        entidadId:      null,
                        entidadNombre:  '',
                        funcionario:    '',
                        asunto:         '',
                        rutaKofax:      '',
                    };
                    break;
                case 'comite_caso':
                    this.form = {
                        decisiones:                   [],
                        observacionesOficio:          '',
                        observacionesRecomendaciones: '',
                        nivelMecanismo:               null,
                        observacionesMecanismo:       '',
                    };
                    break;
                case 'Corregir oficio':
                    this.form = {};
                    console.log('[CTM] E01 Corregir oficio — form inicializado | tarea:', this.tarea.id);
                    this._cargarOficioVinculado();
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
            this.entidades  = [];
            console.log('[CTM] E02 dept seleccionado:', deptId, '| ciudades disponibles:', this.ciudades.length);
        },

        async onSelectCiudad() {
            const ciudadId  = this.form.ciudadId;
            this.form.municipioId = null;
            this.form.entidadId   = null;
            this.form.entidadNombre = '';
            this.municipios = [];
            this.entidades  = [];
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
            const townCode  = this.form.municipioId;
            this.form.entidadId   = null;
            this.form.entidadNombre = '';
            this.entidades  = [];
            if (!townCode) return;
            console.log('[CTM] E04 municipio seleccionado:', townCode, '— cargando entidades...');
            try {
                const res = await fetch('/api/v1/entity-branches?town_code=' + encodeURIComponent(townCode));
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
        },

        onChangeGeneraOficio() {
            console.log('[CTM] E06 generaOficio seleccionado →', this.form.generaOficio);
        },

        /* ── Comité decisions ────────────────────────────────────────────── */

        onChangeDecision(valor) {
            const idx = this.form.decisiones.indexOf(valor);
            if (idx > -1) {
                this.form.decisiones.splice(idx, 1);
                if (valor === 'oficio')                this.form.observacionesOficio          = '';
                if (valor === 'recomendaciones_agente') this.form.observacionesRecomendaciones = '';
                if (valor === 'mecanismo_articulador') {
                    this.form.nivelMecanismo        = null;
                    this.form.observacionesMecanismo = '';
                }
            } else {
                this.form.decisiones.push(valor);
            }
            console.log('[CTM] E07 decisión', idx > -1 ? 'removida' : 'agregada', '→', valor, '| decisiones:', JSON.stringify(this.form.decisiones));
        },

        tieneDecision(valor) {
            return this.form.decisiones && this.form.decisiones.includes(valor);
        },

        onChangeNivelMecanismo() {
            console.log('[CTM] E08 nivelMecanismo seleccionado →', this.form.nivelMecanismo);
        },

        /* ── Build formData ──────────────────────────────────────────────── */

        _buildFormData() {
            const f = this.form;
            const entidadNombreLabel = (f.entidadId && f.entidadId !== 'otra')
                ? (this.entidades.find(e => String(e.id) === String(f.entidadId)) || {}).name || ''
                : f.entidadNombre;

            switch (this.tarea.type) {
                case 'gestion_llamada':
                    return {
                        departamentoId: f.departamentoId,
                        ciudadId:       f.ciudadId,
                        municipioId:    f.municipioId,
                        entidadId:      f.entidadId !== 'otra' ? Number(f.entidadId) : null,
                        entidadNombre:  entidadNombreLabel,
                        funcionario:    f.funcionario,
                        descripcion:    f.descripcion || null,
                        generaOficio:   !!f.generaOficio,
                        asunto:         f.generaOficio ? f.asunto   : null,
                        rutaKofax:      f.generaOficio ? f.rutaKofax : null,
                    };
                case 'proyectar_oficio':
                    return {
                        departamentoId: f.departamentoId,
                        ciudadId:       f.ciudadId,
                        municipioId:    f.municipioId,
                        entidadId:      f.entidadId !== 'otra' ? Number(f.entidadId) : null,
                        entidadNombre:  entidadNombreLabel,
                        funcionario:    f.funcionario,
                        asunto:         f.asunto,
                        rutaKofax:      f.rutaKofax,
                    };
                case 'comite_caso':
                    return {
                        decisiones:                   f.decisiones,
                        observacionesOficio:          this.tieneDecision('oficio')                ? f.observacionesOficio          || null : null,
                        observacionesRecomendaciones: this.tieneDecision('recomendaciones_agente') ? f.observacionesRecomendaciones || null : null,
                        nivelMecanismo:               this.tieneDecision('mecanismo_articulador')  ? f.nivelMecanismo                       : null,
                        observacionesMecanismo:       this.tieneDecision('mecanismo_articulador')  ? f.observacionesMecanismo       || null : null,
                    };
                case 'Corregir oficio':
                    return {};
            }
        },

        /* ── Submit ──────────────────────────────────────────────────────── */

        async confirmar() {
            if (!this.formularioValido || this.guardando) return;
            this.guardando = true;
            this.saveError = null;

            const formData = this._buildFormData();
            const payload  = { userId: this.currentUserId, formData };
            console.log('[CTM] E09 confirmar — tarea:', this.tarea.id, '| type:', this.tarea.type);
            console.log('[CTM] E09 payload formData:', JSON.stringify(formData));

            try {
                const res = await fetch('/api/v1/case-tasks/' + this.tarea.id + '/complete', {
                    method:  'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body:    JSON.stringify(payload),
                });

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
                this.cancelar();

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
<div class="ctm-modal">

    <!-- Header -->
    <div class="ctm-header">
        <div class="ctm-header-left">
            <div class="ctm-header-icon">
                <i :class="'fas ' + iconoTipo"></i>
            </div>
            <div>
                <p class="ctm-header-title">\${ titulo }</p>
                <p v-if="tarea" class="ctm-header-sub">ID \${ tarea.id.substring(0,8) }…</p>
            </div>
        </div>
        <button class="ctm-close-btn" @click="cancelar" :disabled="guardando">✕</button>
    </div>

    <!-- Loader overlay mientras guarda -->
    <div v-if="guardando" style="position:absolute;inset:0;background:rgba(255,255,255,.65);border-radius:16px;display:flex;align-items:center;justify-content:center;z-index:10;flex-direction:column;gap:10px">
        <svg style="width:32px;height:32px;animation:ctm-spin .7s linear infinite" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
            <path d="M12 2a10 10 0 0 1 10 10" stroke="#1d4ed8" stroke-width="3" stroke-linecap="round"/>
        </svg>
        <span style="font-size:13px;font-weight:600;color:#1d4ed8">Guardando...</span>
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

            <div class="ctm-field">
                <label class="ctm-label">Sede / Entidad <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.entidadId" :disabled="!form.municipioId" @change="onSelectEntidad">
                    <option :value="null" disabled>Selecciona entidad</option>
                    <option v-for="e in entidades" :key="e.id" :value="e.id">\${ e.name }</option>
                    <option value="otra">Otra entidad…</option>
                </select>
            </div>

            <div v-if="form.entidadId === 'otra'" class="ctm-field">
                <label class="ctm-label">Nombre de la entidad <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.entidadNombre" placeholder="Escribe el nombre de la entidad"/>
            </div>

            <p class="ctm-section-title">Contacto</p>

            <div class="ctm-field">
                <label class="ctm-label">Funcionario contactado <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.funcionario" placeholder="Nombre del funcionario"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Descripción / Notas</label>
                <textarea class="ctm-textarea" v-model="form.descripcion" rows="3" placeholder="Resumen de la gestión..."></textarea>
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
                    <input class="ctm-input" v-model="form.asunto" placeholder="Asunto…"/>
                </div>
                <div class="ctm-field">
                    <label class="ctm-label">Ruta Kofax / URL del archivo <span class="req">*</span></label>
                    <input class="ctm-input" v-model="form.rutaKofax" placeholder="\\\\servidor\\ruta\\archivo.pdf"/>
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
                <label class="ctm-label">Sede / Entidad <span class="req">*</span></label>
                <select class="ctm-select" v-model="form.entidadId" :disabled="!form.municipioId">
                    <option :value="null" disabled>Selecciona entidad</option>
                    <option v-for="e in entidades" :key="e.id" :value="e.id">\${ e.name }</option>
                    <option value="otra">Otra entidad…</option>
                </select>
            </div>

            <div v-if="form.entidadId === 'otra'" class="ctm-field">
                <label class="ctm-label">Nombre de la entidad <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.entidadNombre" placeholder="Escribe el nombre de la entidad"/>
            </div>

            <p class="ctm-section-title">Datos del oficio</p>

            <div class="ctm-field">
                <label class="ctm-label">Dependencia del funcionario <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.funcionario" placeholder="Ej. Departamento de Protección Social"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Asunto <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.asunto" placeholder="Asunto del oficio…"/>
            </div>

            <div class="ctm-field">
                <label class="ctm-label">Ruta Kofax / URL del archivo <span class="req">*</span></label>
                <input class="ctm-input" v-model="form.rutaKofax" placeholder="\\\\servidor\\ruta\\archivo.pdf"/>
            </div>

        </template>

        <!-- Formulario comite_caso -->
        <template v-else-if="tarea && tarea.type === 'comite_caso'">

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

            <template v-if="tieneDecision('oficio')">
                <p class="ctm-section-title" style="margin-top:16px">Oficio</p>
                <div class="ctm-field">
                    <label class="ctm-label">Observaciones del oficio</label>
                    <textarea class="ctm-textarea" v-model="form.observacionesOficio" rows="2" placeholder="Observaciones sobre el oficio a generar…"></textarea>
                </div>
            </template>

            <template v-if="tieneDecision('recomendaciones_agente')">
                <p class="ctm-section-title" style="margin-top:16px">Recomendaciones</p>
                <div class="ctm-field">
                    <label class="ctm-label">Observaciones para el agente</label>
                    <textarea class="ctm-textarea" v-model="form.observacionesRecomendaciones" rows="2" placeholder="Recomendaciones al agente del caso…"></textarea>
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
                    <label class="ctm-label">Observaciones del mecanismo</label>
                    <textarea class="ctm-textarea" v-model="form.observacionesMecanismo" rows="2" placeholder="Observaciones sobre el mecanismo…"></textarea>
                </div>
            </template>

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

        <!-- Error al guardar -->
        <div v-if="saveError" class="ctm-error-banner">
            <i class="fas fa-exclamation-circle"></i> \${ saveError }
        </div>

    </div>

    <!-- Footer -->
    <div class="ctm-footer">
        <button class="ctm-btn ctm-btn-cancel" @click="cancelar" :disabled="guardando">
            Cancelar
        </button>
        <button class="ctm-btn ctm-btn-confirm" @click="confirmar"
            :disabled="!formularioValido || guardando || cargandoTarea">
            <template v-if="guardando">
                <i class="fas fa-circle-notch fa-spin"></i> Guardando…
            </template>
            <template v-else>
                <template v-if="tarea && tarea.type === 'Corregir oficio'">
                    Marcar como corregido
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
