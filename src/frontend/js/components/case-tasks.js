/**
 * case-tasks
 * Componente reutilizable: gestión de tareas de un caso.
 * Muestra tareas pendientes y completadas, permite completarlas con modales según tipo.
 *
 * Props:
 *   - caseId   (String, required): icode del caso
 *   - userId   (String, optional): icode del usuario logueado (para permisos)
 *   - userRole (String, optional): rol del usuario ('sv', 'ro', 'op', 'ps', 'ts', 'et')
 *   - allowPsicoManage (Boolean, optional): si true, ps/ts pueden gestionar
 *     tareas asignadas a ellos (p. ej. detalle de remisión psicosocial).
 *     Default false: en detalle de caso ps/ts son solo lectura.
 * Events:
 *   - @task-completed: emitido cuando se completa una tarea
 */

(function injectCaseTasksStyles() {
    if (document.getElementById('case-tasks-styles')) return;
    var style = document.createElement('style');
    style.id = 'case-tasks-styles';
    style.textContent = `
        .ct-container { padding: 16px 0 }
        .ct-summary { display:flex;gap:16px;margin-bottom:20px }
        .ct-card { flex:1;border:1px solid #e5e7eb;border-radius:10px;padding:16px;text-align:center }
        .ct-card-label { font-size:.72rem;color:#6b7280;text-transform:uppercase;font-weight:600;margin:0 0 4px }
        .ct-card-value { font-size:1.4rem;font-weight:700;margin:0 }
        .ct-card-pending .ct-card-value { color:#f97316 }
        .ct-card-done .ct-card-value { color:#22c55e }
        .ct-list { display:flex;flex-direction:column;gap:10px }
        .ct-list::before { content:none;display:none }
        .ct-task { border:none!important;border-radius:10px;padding:14px 16px;display:flex;align-items:center;justify-content:space-between;transition:background .12s;background:#f9fafb;box-shadow:none }
        .ct-task:hover { background:#f9fafb }
        .ct-task-info { flex:1 }
        .ct-task-category { font-size:.68rem;color:#9ca3af;text-transform:uppercase;font-weight:600 }
        .ct-task-desc { font-size:.84rem;color:#374151;margin:4px 0 0;font-weight:500 }
        .ct-task-meta { font-size:.72rem;color:#9ca3af;margin:4px 0 0 }
        .ct-badge { display:inline-block;padding:2px 8px;border-radius:10px;font-size:.68rem;font-weight:600 }
        .ct-badge-todo { background:#fff7ed;color:#f97316;border:1px solid #fed7aa }
        .ct-badge-done { background:#f0fdf4;color:#22c55e;border:1px solid #bbf7d0 }
        .ct-badge-devuelto { background:#fffbeb;color:#b45309;border:1px solid #fde68a }
        .ct-btn-complete { background:#5106A7;color:#fff;border:none;padding:6px 14px;border-radius:6px;font-size:.78rem;font-weight:600;cursor:pointer }
        .ct-btn-complete:hover { background:#3d0480 }
        .ct-btn-complete:disabled { opacity:.5;cursor:not-allowed }
        .ct-empty { text-align:center;padding:40px 0;color:#9ca3af;font-size:.85rem }
        .ct-loading { text-align:center;padding:40px 0;color:#9ca3af }
        .ct-intentos { border:1px solid #e5e7eb;border-radius:12px;padding:14px 16px;margin-bottom:16px;background:#fff }
        .ct-intentos-head { display:flex;align-items:center;gap:10px;flex-wrap:wrap }
        .ct-intentos-head strong { font-size:.9rem;color:#111827 }
        .ct-intentos-count { margin-left:auto;font-size:.84rem;font-weight:700;color:#6b7280 }
        .ct-intentos-count.is-limite { color:#dc2626 }
        .ct-intentos p { margin:8px 0 0;font-size:.82rem;color:#9ca3af }
        .ct-modal-backdrop { position:fixed;inset:0;background:rgba(0,0,0,.5);display:flex;align-items:center;justify-content:center;z-index:10000 }
        .ct-modal { background:#fff;border-radius:12px;width:95%;max-width:550px;padding:24px;box-shadow:0 20px 50px rgba(0,0,0,.2) }
        .ct-modal-title { font-size:1rem;font-weight:700;margin:0 0 16px;color:#111827 }
        .ct-modal-field { margin-bottom:14px }
        .ct-modal-label { font-size:.75rem;font-weight:600;color:#6b7280;margin-bottom:4px;display:block }
        .ct-modal-input { width:100%;border:1px solid #d1d5db;border-radius:8px;padding:10px 12px;font-size:.84rem;resize:vertical;min-height:80px }
        .ct-modal-footer { display:flex;justify-content:flex-end;gap:10px;margin-top:16px }
        .ct-modal-btn-cancel { background:#f3f4f6;color:#374151;border:none;padding:8px 16px;border-radius:6px;font-size:.82rem;cursor:pointer }
        .ct-modal-btn-save { background:#5106A7;color:#fff;border:none;padding:8px 16px;border-radius:6px;font-size:.82rem;font-weight:600;cursor:pointer }
        .ct-modal-btn-save:disabled { opacity:.5;cursor:not-allowed }
        .ct-tabs { display:flex;gap:0;border-bottom:2px solid #e5e7eb;margin-bottom:16px }
        .ct-tab { padding:8px 16px;font-size:.8rem;font-weight:600;color:#6b7280;cursor:pointer;border-bottom:2px solid transparent;margin-bottom:-2px }
        .ct-tab.active { color:#5106A7;border-bottom-color:#5106A7 }
    `;
    document.head.appendChild(style);
})();

app.component('case-tasks', {
    delimiters: ['${', '}'],
    props: {
        caseId:   { type: String, required: true },
        userId:   { type: String, default: '' },
        userRole: { type: String, default: '' },
        barrierId: { type: String, default: '' },
        psychosocialId: { type: String, default: '' },
        puedeGestionPropia: { type: Boolean, default: false },
        allowPsicoManage: { type: Boolean, default: false },
        entityBranchId: { default: null },
        entityName: { type: String, default: '' },
        barrierDescription: { type: String, default: '' },
        referralId: { type: String, default: '' },
    },
    emits: ['task-completed'],
    data: function() {
        return {
            loading: false,
            tasks: [],
            activeTab: 'pending',
            reasignarTask: null,
            reasignarAgente: '',
            agentesEquipo: [],
            reasignando: false,
            // Set de ids de tareas con una gestión directa en vuelo. Evita el
            // doble-click: el botón "Gestionar" se deshabilita mientras su id
            // esté aquí. Es un objeto plano usado como set (clave = task.id).
            gestionandoDirecto: {},
        };
    },
    computed: {
        scopedTasks: function() {
            var self = this;
            var tasks = this.tasks;
            // Rol Entidad ('et'): filtro original del equipo — solo sus tareas asignadas.
            if (this.userRole === 'et') {
                tasks = tasks.filter(function(t) {
                    return String(t.assignedUserId || '') === String(self.userId || '');
                });
            } else if (this.userRole === 'te') {
                // Equipo Técnico ('te'): ve las tareas asignadas a sí mismo y, además, las
                // tareas "comite_caso" del pool (sin dueño), disponibles para que cualquier
                // integrante del equipo técnico las tome y gestione.
                tasks = tasks.filter(function(t) {
                    if (String(t.assignedUserId || '') === String(self.userId || '')) return true;
                    if (t.type === 'comite_caso' && !t.assignedUserId) return true;
                    return false;
                });
            }
            return tasks.filter(function(t) {
                if (!t.assignedRole) return true;
                return self.userRole === t.assignedRole || self.userRole === 'sv';
            });
        },
        pendingTasks: function() {
            return this.scopedTasks.filter(function(t) { return t.status === 'ToDo'; });
        },
        completedTasks: function() {
            return this.scopedTasks.filter(function(t) { return t.status === 'Done'; });
        },
        visibleTasks: function() {
            return this.activeTab === 'pending' ? this.pendingTasks : this.completedTasks;
        },
        intentosFallidos: function() {
            if (!this.referralId) return 0;
            var n = 0;
            (this.tasks || []).forEach(function(t) {
                if (!t || t.type !== 'intentar_contacto_3x3' || t.status !== 'Done') return;
                var fd = t.formData;
                if (typeof fd === 'string') {
                    try { fd = JSON.parse(fd); } catch (e) { fd = null; }
                }
                if (fd && fd.contesto === 'no') n += 1;
            });
            return n;
        },
        intentosEtiqueta: function() {
            var n = this.intentosFallidos;
            if (n >= 10) return '10/10';
            return n + '/9';
        },
    },
    mounted: function() {
        this.cargar();
    },
    methods: {
        async cargar() {
            this.loading = true;
            try {
            var url = this.referralId
                ? '/api/v1/case-tasks?referralId=' + encodeURIComponent(this.referralId)
                : '/api/v1/case-tasks?caseId=' + this.caseId;
                if (this.barrierId) url += '&barrierId=' + this.barrierId;
                if (this.psychosocialId) url += '&psychosocialId=' + this.psychosocialId;
                var res = await fetch(url);
                if (!res.ok) throw new Error('Error ' + res.status);
                this.tasks = await res.json();
                if (!Array.isArray(this.tasks)) this.tasks = [];
            } catch (e) {
                console.error('Error cargando tareas:', e);
                this.tasks = [];
            } finally {
                this.loading = false;
            }
        },
        abrirModal: function(task) {
            if (this.$refs.taskModal) {
                this.$refs.taskModal.open(task.id);
            }
        },
        // esGestionDirecta: los 7 tipos nuevos de la cadena de medidas que, por
        // ahora, se completan con un click directo al endpoint existente SIN abrir
        // case-task-modal (no hay formularios reales todavía). Al completarse, el
        // backend mueve el estado del referral y crea la tarea siguiente.
        esGestionDirecta: function(task) {
            if (!task) return false;
            // revisar_caracterizacion NO está aquí: tiene su propio modal real
            // (case-task-modal), por eso su botón dice "Revisar" y abre el modal.
            // revisar_caracterizacion y revisar_trazabilidad NO están aquí: tienen
            // su propio modal real (case-task-modal), con botón "Revisar".
            // cargar_cotizaciones tampoco: tiene su modal real "Cargar cotizaciones".
            // aprobar_cotizaciones tampoco: tiene su modal real "Revisar y aprobar
            // cotizaciones" (botón "Revisar").
            // gestionar_territorio y marcar_otorgadas tampoco: tienen su modal real
            // (case-task-modal), con botón "Registrar".
            // entrevista_caracterizacion tampoco: ahora tiene su modal real de 3
            // pasos (botón "Gestionar"); ya NO es gestión directa. Hoy no queda
            // ningún tipo en gestión directa.
            var tiposDirectos = {};
            return !!tiposDirectos[task.type];
        },
        // labelBotonGestion: texto del botón de acción de una tarea pendiente.
        // La revisión de caracterización/trazabilidad dice "Revisar"; la carga de
        // cotizaciones dice "Cargar cotizaciones"; la gestión en territorio y el
        // marcado de otorgadas dicen "Registrar"; el resto conserva "Gestionar".
        // Todos abren el case-task-modal (menos los de gestión directa).
        labelBotonGestion: function(task) {
            if (task && (task.type === 'revisar_caracterizacion' || task.type === 'revisar_trazabilidad' || task.type === 'aprobar_cotizaciones')) return 'Revisar';
            if (task && task.type === 'cargar_cotizaciones') return 'Cargar cotizaciones';
            if (task && (task.type === 'gestionar_territorio' || task.type === 'marcar_otorgadas')) return 'Registrar';
            return 'Gestionar';
        },
        // parseFormData devuelve el form_data de una tarea como objeto plano
        // (el backend puede enviarlo como objeto o como string JSON).
        parseFormData: function(task) {
            var fd = task && task.formData;
            if (!fd) return {};
            if (typeof fd === 'string') {
                try { return JSON.parse(fd) || {}; } catch (e) { return {}; }
            }
            return fd;
        },
        // esRevisarDevuelta: la revisión de caracterización, la de trazabilidad o la
        // aprobación de cotizaciones se cerró como "devuelto" (el supervisor la
        // devolvió). Se usa para pintar el chip ámbar "Devuelto" y las observaciones.
        esRevisarDevuelta: function(task) {
            if (!task || (task.type !== 'revisar_caracterizacion' && task.type !== 'revisar_trazabilidad' && task.type !== 'aprobar_cotizaciones')) return false;
            return this.parseFormData(task).resultado === 'devuelto';
        },
        // observacionesDevolucion: texto de las observaciones guardadas al
        // devolver la revisión.
        observacionesDevolucion: function(task) {
            return String(this.parseFormData(task).observaciones || '').trim();
        },
        // ── Cierres de la entrevista de caracterización (rol am) ──
        // La entrevista puede cerrarse sin avanzar la cadena en dos ramas, ambas
        // dejan la remisión descartada: "Concluir gestión" (resultado
        // concluido_sin_me) y "No otorga consentimiento" (sin_consentimiento). Se
        // pintan con su propio chip + la observación debajo (mismo patrón visual que
        // el chip "Devuelto", pero con texto propio).
        esCierreCaracterizacion: function(task) {
            if (!task || task.type !== 'entrevista_caracterizacion') return false;
            var r = this.parseFormData(task).resultado;
            return r === 'concluido_sin_me' || r === 'sin_consentimiento';
        },
        // chipCierreTexto: etiqueta del chip según la rama de cierre.
        chipCierreTexto: function(task) {
            var r = this.parseFormData(task).resultado;
            if (r === 'concluido_sin_me') return 'Concluida sin ME';
            if (r === 'sin_consentimiento') return 'Sin consentimiento';
            return '';
        },
        // chipCierreDetalle: línea descriptiva bajo el título (motivo cuando aplica
        // + observación). concluido_sin_me trae motivo; sin_consentimiento solo obs.
        chipCierreDetalle: function(task) {
            var fd = this.parseFormData(task);
            var obs = String(fd.observacion || '').trim();
            if (fd.resultado === 'concluido_sin_me') {
                var motivos = {
                    desistimiento: 'Desistimiento de la ciudadana',
                    sin_riesgo: 'No hay riesgo de feminicidio y tiene rutas de atención activas',
                };
                var motivo = motivos[fd.motivo] || fd.motivo || '';
                return [motivo, obs].filter(Boolean).join(' — ');
            }
            return obs;
        },
        // gestionarDirecto completa la tarea de la cadena de medidas directamente,
        // sin modal. Reutiliza el mismo patrón de fetch/manejo de errores del
        // componente (401→landing, 422→error de transición, resto→mensaje del
        // backend). Tras el éxito, refresca la lista para que aparezca la tarea
        // siguiente. "resultado: aprobado" es el camino feliz hacia adelante
        // (solo importa para aprobar_cotizaciones; en los demás es irrelevante).
        async gestionarDirecto(task) {
            if (!task || this.gestionandoDirecto[task.id]) return;
            this.gestionandoDirecto[task.id] = true;
            try {
                var res = await fetch('/api/v1/case-tasks/' + task.id + '/complete', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        userId: this.userId,
                        formData: { resultado: 'aprobado' },
                    }),
                });
                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }
                var data = null;
                try { data = await res.json(); } catch (e) { data = null; }
                if (!res.ok) {
                    var msg = (data && data.error) ? data.error : ('Error ' + res.status);
                    alert(msg);
                    return;
                }
                // Éxito: refrescar el panel (aparece la tarea siguiente) y avisar
                // al padre para que recargue el detalle de la remisión.
                this.cargar();
                this.$emit('task-completed');
            } catch (e) {
                console.error('Error gestionando tarea de medidas:', e);
                alert('No se pudo completar la tarea. Intenta de nuevo.');
            } finally {
                delete this.gestionandoDirecto[task.id];
            }
        },
        onTaskCompleted: function() {
            this.cargar();
            this.$emit('task-completed');
        },
        abrirModalGestionPropia: function() {
            if (!this.$refs.taskModal) return;
            this.$refs.taskModal.openArticulacion({
                barrierId: this.barrierId,
                caseId: this.caseId,
                entityBranchId: this.entityBranchId,
                entityName: this.entityName || '',
                barrierDescription: this.barrierDescription || '',
            });
        },
        // esComiteDelPool: tarea "comite_caso" sin dueño que cualquier integrante del
        // Equipo Técnico ('te') puede tomar y gestionar.
        esComiteDelPool: function(task) {
            return task && task.type === 'comite_caso' && !task.assignedUserId;
        },
        // esSoloAsignado: tipos de tarea que SOLO puede gestionar su usuario asignado,
        // ni siquiera el supervisor. articulacion_territorial es del enlace territorial
        // (rol 'en'): solo él la gestiona. validar_remision es del profesional psicosocial.
        esSoloAsignado: function(task) {
            return task && (task.type === 'articulacion_territorial' || task.type === 'validar_remision');
        },
        puedeCompletar: function(task) {
            if ((this.userRole === 'ps' || this.userRole === 'ts') && !this.allowPsicoManage) return false;
            // Tareas exclusivas del asignado (enlace territorial / profesional): ni el sv
            // ni el te pueden gestionarlas; solo el usuario al que están asignadas.
            if (this.esSoloAsignado(task)) {
                return String(task.assignedUserId || '') === String(this.userId || '');
            }
            if (task.assignedRole) {
                return this.userRole === task.assignedRole || this.userRole === 'sv';
            }
            if (this.userRole === 'sv') return true;
            // El Equipo Técnico puede gestionar las comite_caso del pool.
            if (this.userRole === 'te' && this.esComiteDelPool(task)) return true;
            return String(task.assignedUserId || '') === String(this.userId || '');
        },
        esAsignadoOSv: function(task) {
            // Tareas exclusivas del asignado (articulacion_territorial del enlace y
            // validar_remision del profesional): solo las gestiona su asignado, NO el sv.
            if (this.esSoloAsignado(task)) {
                return String(task.assignedUserId || '') === String(this.userId || '');
            }
            if (task.assignedRole) {
                return this.userRole === task.assignedRole || this.userRole === 'sv';
            }
            // Para las demás tareas: sv siempre puede, o ser el asignado
            if (this.userRole === 'sv') return true;
            // El Equipo Técnico puede gestionar las comite_caso del pool aunque no le estén asignadas.
            if (this.userRole === 'te' && this.esComiteDelPool(task)) return true;
            return String(task.assignedUserId || '') === String(this.userId || '');
        },
        async abrirReasignar(task) {
            this.reasignarTask = task;
            this.reasignarAgente = '';
            this.reasignando = false;
            // Cargar agentes del equipo
            try {
                var res = await fetch('/api/v1/equipo-operadores');
                if (res.ok) {
                    this.agentesEquipo = await res.json();
                    if (!Array.isArray(this.agentesEquipo)) this.agentesEquipo = [];
                }
            } catch(e) { this.agentesEquipo = []; }
        },
        async confirmarReasignar() {
            if (!this.reasignarAgente || !this.reasignarTask) return;
            this.reasignando = true;
            try {
                var res = await fetch('/api/v1/case-tasks/' + this.reasignarTask.id + '/reassign', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ assignedUserId: this.reasignarAgente })
                });
                if (!res.ok) throw new Error('Error');
                this.reasignarTask = null;
                this.cargar();
            } catch(e) {
                console.error('Error reasignando:', e);
            } finally {
                this.reasignando = false;
            }
        },
        formatFecha: function(fecha) {
            if (!fecha) return '';
            // Fecha pura (sin 'T') se toma tal cual; ISO con hora se convierte a hora
            // Colombia antes de recortar (SALVIA se usa solo en Colombia).
            if (fecha.indexOf('T') === -1) return fecha.substring(0, 10);
            var d = new Date(fecha);
            if (isNaN(d.getTime())) return fecha.substring(0, 10);
            var parts = new Intl.DateTimeFormat('en-CA', {
                timeZone: 'America/Bogota', year: 'numeric', month: '2-digit', day: '2-digit'
            }).formatToParts(d).reduce(function(acc, p) { acc[p.type] = p.value; return acc; }, {});
            return parts.year + '-' + parts.month + '-' + parts.day;
        },
        labelTipo: function(task) {
            var type = task.type || '';
            var labels = {
                'escribir_oficio_cierre': 'Escribir oficio cierre de caso',
                'escribir_oficio_barrera': 'Escribir oficio barrera',
                'corregir_oficio_barrera': 'Corregir oficio barrera',
                'llamar_entidad_barrera': 'Llamar a entidad barrera',
                'gestion_enlace_territorial': 'Gestión barrera enlace territorial',
                'gestion_propia': 'Gestión propia',
                'articulacion_territorial': 'Articulación territorial',
                'compromiso_con_enlace': 'Compromiso con Enlace',
                'solicitud_entidad': 'Hacer solicitudes a la entidad',
                'respuesta_entidad': 'Responder solicitud de la entidad',
                'decidir_via_riesgo_inminente': 'Decidir vía de atención',
                'intentar_contacto_3x3': 'Intentar contacto (3×3)',
                'entrevista_caracterizacion': 'Contactar víctima y realizar entrevista de caracterización',
                'revisar_caracterizacion': 'Revisar caracterización',
                'revisar_trazabilidad': 'Revisar trazabilidad',
                'cargar_cotizaciones': 'Cargar cotizaciones',
                'aprobar_cotizaciones': 'Aprobar cotizaciones',
                'gestionar_territorio': 'Gestionar territorio',
                'marcar_otorgadas': 'Marcar medidas como otorgadas',
                'escalar_cruz_roja': 'Escalar a Cruz Roja',
                'entregar_medida_inmediata': 'Entregar medida inmediata (sin caracterización previa)',
                'confirmar_lugar_seguro': 'Confirmar que la persona está en lugar seguro',
                'revision_equipo_tecnico': 'Revisión Equipo Técnico',
                'asignar_agente_medidas': 'Asignar agente de medidas',
                'ajustar_remision_devuelta': 'Ajustar remisión devuelta por el Equipo Técnico',
            };
            return labels[type] || type || task.category || 'Tarea';
        },
        nombreAsignado: function(task) {
            if (task.assignedRole) {
                var nombres = {
                    sm: 'Supervisor de Medidas',
                    sv: 'Supervisor',
                    em: 'Equipo Técnico de Medidas',
                    sc: 'Supervisión del Convenio',
                    cr: 'Cruz Roja',
                    am: 'Agente de Medidas',
                    en: 'Enlace Territorial',
                };
                return nombres[task.assignedRole] || task.assignedRole;
            }
            return task.assignedUserName || 'Sin asignar';
        },
        tituloTarea: function(task) {
            if (task.type === 'solicitud_entidad') {
                return 'Hacer solicitudes a ' + (task.entityName || 'la entidad');
            }
            // comite_caso (Decisiones del Comité / mesa técnica) usa un título fijo:
            // el texto específico de cada una (justificación al abrirla, o el
            // contexto de la barrera cuando nace de "Alerta por barreras") se ve
            // recortado en el bloque "Solicitud" más abajo, no como título.
            if (task.type === 'comite_caso') {
                return 'Análisis de la mesa técnica';
            }
            return task.description || this.labelTipo(task);
        },
        // primerParrafo recorta un texto largo a su primer párrafo (o ~160
        // caracteres si esa línea es más larga), para que no rompa el layout de
        // la tarjeta. El texto completo se ve al entrar a "Gestionar".
        primerParrafo: function(texto) {
            if (!texto) return '';
            var limite = 160;
            var primeraLinea = texto.split('\n')[0].trim();
            if (primeraLinea.length > limite) {
                return primeraLinea.substring(0, limite).trim() + '…';
            }
            if (primeraLinea.length < texto.trim().length) {
                return primeraLinea + '…';
            }
            return primeraLinea;
        },
        // entidadTarea resuelve el nombre de la entidad relacionada a la tarea.
        // El backend expone la entidad por dos rutas: 'entityName' (cuando la tarea
        // apunta a una entidad concreta del caso vía entity_case) y 'barrierEntityName'
        // (la sede/entidad ligada a la barrera). Se prioriza la primera.
        entidadTarea: function(task) {
            return (task && (task.entityName || task.barrierEntityName)) || '';
        },
        // solicitudTarea es el texto de la solicitud que se le hace a la entidad,
        // o la justificación de la tarea comite_caso cuando no hay entidad. Se
        // recorta al primer párrafo para no romper el layout de la tarjeta.
        solicitudTarea: function(task) {
            return this.primerParrafo(task && task.description);
        },
        // esTareaEnlace indica si la tarjeta debe mostrar el bloque "Entidad /
        // Solicitud": cuando hay una entidad resuelta, la tarea nace de una
        // barrera, o es una tarea comite_caso (mesa técnica), donde "Solicitud"
        // muestra la justificación aunque no haya entidad.
        esTareaEnlace: function(task) {
            if (!task) return false;
            return !!this.entidadTarea(task) || !!task.barrierId || task.type === 'comite_caso';
        },
    },
    template: `
<div class="ct-container">
    <div v-if="loading" class="ct-loading"><i class="fa fa-spinner fa-spin"></i> Cargando tareas...</div>
    <template v-else>
        <div v-if="referralId" class="ct-intentos">
            <div class="ct-intentos-head">
                <strong>Intentos de contacto fallidos (3×3)</strong>
                <span class="ct-intentos-count" :class="{ 'is-limite': intentosFallidos >= 9 }">\${ intentosEtiqueta }</span>
            </div>
            <p v-if="intentosFallidos === 0">Sin intentos registrados.</p>
        </div>
        <!-- Tabs: Pendientes / Completadas + acción Registrar gestión propia -->
        <div v-if="tasks.length > 0 || puedeGestionPropia" style="display:flex;align-items:center;justify-content:space-between;gap:8px;margin-bottom:14px;border-bottom:1px solid #e5e7eb;padding-bottom:0">
            <div v-if="tasks.length > 0" style="display:flex;gap:4px">
                <button @click="activeTab = 'pending'" style="padding:8px 16px;font-size:.78rem;font-weight:600;border:none;cursor:pointer;border-bottom:2px solid transparent;background:none" :style="activeTab === 'pending' ? 'color:#5106A7;border-bottom-color:#5106A7' : 'color:#6b7280'">
                    Pendientes (\${ pendingTasks.length })
                </button>
                <button @click="activeTab = 'completed'" style="padding:8px 16px;font-size:.78rem;font-weight:600;border:none;cursor:pointer;border-bottom:2px solid transparent;background:none" :style="activeTab === 'completed' ? 'color:#5106A7;border-bottom-color:#5106A7' : 'color:#6b7280'">
                    Completadas (\${ completedTasks.length })
                </button>
            </div>
            <div v-else></div>
            <button v-if="puedeGestionPropia" @click="abrirModalGestionPropia" class="ct-btn-complete" style="flex-shrink:0;margin-bottom:8px">
                + Registrar gestión propia
            </button>
        </div>

        <!-- Lista de tareas pendientes -->
        <template v-if="activeTab === 'pending'">
            <div v-if="pendingTasks.length === 0" class="ct-empty">
                No hay tareas pendientes
            </div>
            <div v-else class="ct-list">
                <div v-for="task in pendingTasks" :key="task.id" class="ct-task" style="flex-direction:row;align-items:center;gap:12px;padding:14px 20px">
                    <div v-if="task.type === 'solicitud_entidad'" style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#dbeafe">
                        <i class="fa fa-folder" style="color:#2563eb;font-size:.82rem"></i>
                    </div>
                    <div v-else-if="task.type === 'respuesta_entidad'" style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#ecfdf5">
                        <i class="fa fa-reply" style="color:#059669;font-size:.82rem"></i>
                    </div>
                    <div v-else style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0" :style="task.type && task.type.indexOf('oficio') !== -1 ? 'background:#f3f4f6' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'background:#fce7f3' : task.type && task.type.indexOf('comite') !== -1 ? 'background:#fef9c3' : 'background:#ede9fe'">
                        <i :class="task.type && task.type.indexOf('oficio') !== -1 ? 'fa fa-file-alt' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'fa fa-phone' : task.type && task.type.indexOf('comite') !== -1 ? 'fa fa-clipboard-list' : 'fa fa-tasks'" :style="task.type && task.type.indexOf('oficio') !== -1 ? 'color:#1f2937' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'color:#db2777' : task.type && task.type.indexOf('comite') !== -1 ? 'color:#92400e' : 'color:#7c3aed'" style="font-size:.82rem"></i>
                    </div>
                    <div style="flex:1;min-width:0">
                        <div style="display:flex;align-items:center;gap:8px;margin-bottom:2px">
                            <span style="font-size:.7rem;font-weight:800;color:#1f2937;text-transform:uppercase;letter-spacing:.04em">\${ task.category }</span>
                            <span class="ct-badge ct-badge-todo">Pendiente</span>
                        </div>
                        <div style="font-size:.84rem;font-weight:700;color:#1f2937">\${ tituloTarea(task) }</div>
                        <!-- Entidad y solicitud: guían al enlace territorial sobre cómo gestionar -->
                        <div v-if="esTareaEnlace(task)" style="margin-top:4px;display:flex;flex-direction:column;gap:2px">
                            <div v-if="entidadTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-building" style="color:#2563eb;margin-right:5px"></i><strong>Entidad:</strong> \${ entidadTarea(task) }</div>
                            <div v-if="solicitudTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-file-lines" style="color:#7c3aed;margin-right:5px"></i><strong>Solicitud:</strong> \${ solicitudTarea(task) }</div>
                        </div>
                        <div style="font-size:.73rem;color:#9ca3af;margin-top:2px">Asignado a: \${ nombreAsignado(task) }</div>
                    </div>
                    <div style="display:flex;gap:6px;flex-shrink:0">
                        <!-- Tipos de la cadena de medidas: click directo (sin modal) al endpoint complete. -->
                        <button v-if="esGestionDirecta(task) && puedeCompletar(task) && esAsignadoOSv(task)" class="ct-btn-complete" :disabled="!!gestionandoDirecto[task.id]" @click="gestionarDirecto(task)">\${ gestionandoDirecto[task.id] ? 'Gestionando...' : 'Gestionar' }</button>
                        <!-- Resto de tipos: abren el case-task-modal como siempre. La
                             revisión de caracterización muestra "Revisar" en el botón. -->
                        <button v-else-if="puedeCompletar(task) && esAsignadoOSv(task)" class="ct-btn-complete" @click="abrirModal(task)">\${ labelBotonGestion(task) }</button>
                        <button v-else-if="puedeCompletar(task) && !esAsignadoOSv(task)" class="ct-btn-complete" disabled style="opacity:.5;cursor:not-allowed" title="Solo el profesional asignado puede gestionar esta tarea">\${ labelBotonGestion(task) }</button>
                        <button v-if="false && userRole === 'sv'" class="ct-btn-complete" style="background:#3b82f6" @click="abrirReasignar(task)">Reasignar</button>
                    </div>
                </div>
            </div>
        </template>

        <!-- Lista de tareas completadas -->
        <template v-if="activeTab === 'completed'">
            <div v-if="completedTasks.length === 0" class="ct-empty">
                No hay tareas completadas aún
            </div>
            <div v-else class="ct-list">
                <div v-for="task in completedTasks" :key="task.id" class="ct-task" style="flex-direction:row;align-items:center;gap:12px;padding:14px 20px;background:#f9fafb">
                    <div v-if="esRevisarDevuelta(task) || esCierreCaracterizacion(task)" style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#fef3c7">
                        <i class="fa fa-undo" style="color:#b45309;font-size:.82rem"></i>
                    </div>
                    <div v-else style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#dcfce7">
                        <i class="fa fa-check" style="color:#16a34a;font-size:.82rem"></i>
                    </div>
                    <div style="flex:1;min-width:0">
                        <div style="display:flex;align-items:center;gap:8px;margin-bottom:2px">
                            <span style="font-size:.7rem;font-weight:800;color:#1f2937;text-transform:uppercase;letter-spacing:.04em">\${ task.category }</span>
                            <span v-if="esRevisarDevuelta(task)" class="ct-badge ct-badge-devuelto">Devuelto</span>
                            <span v-else-if="esCierreCaracterizacion(task)" class="ct-badge ct-badge-devuelto">\${ chipCierreTexto(task) }</span>
                            <span v-else class="ct-badge ct-badge-done">Completada</span>
                        </div>
                        <div style="font-size:.84rem;font-weight:700;color:#1f2937">\${ tituloTarea(task) }</div>
                        <!-- Observaciones de la devolución (revisar caracterización devuelta) -->
                        <div v-if="esRevisarDevuelta(task) && observacionesDevolucion(task)" style="margin-top:4px;font-size:.74rem;color:#92400e"><i class="fa fa-comment-dots" style="margin-right:5px"></i>\${ observacionesDevolucion(task) }</div>
                        <!-- Motivo + observación del cierre de la entrevista (concluida sin ME / sin consentimiento) -->
                        <div v-if="esCierreCaracterizacion(task) && chipCierreDetalle(task)" style="margin-top:4px;font-size:.74rem;color:#92400e"><i class="fa fa-comment-dots" style="margin-right:5px"></i>\${ chipCierreDetalle(task) }</div>
                        <!-- Entidad y solicitud: guían al enlace territorial sobre cómo gestionar -->
                        <div v-if="esTareaEnlace(task)" style="margin-top:4px;display:flex;flex-direction:column;gap:2px">
                            <div v-if="entidadTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-building" style="color:#2563eb;margin-right:5px"></i><strong>Entidad:</strong> \${ entidadTarea(task) }</div>
                            <div v-if="solicitudTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-file-lines" style="color:#7c3aed;margin-right:5px"></i><strong>Solicitud:</strong> \${ solicitudTarea(task) }</div>
                        </div>
                        <div style="font-size:.73rem;color:#9ca3af;margin-top:2px">Completada por: \${ task.assignedUserName || '—' }</div>
                    </div>
                    <div style="display:flex;gap:6px;flex-shrink:0">
                        <button class="ct-btn-complete" style="background:#16a34a" @click="$refs.taskHistory.open(task.id)">Ver detalle</button>
                    </div>
                </div>
            </div>
        </template>
    </template>

    <!-- Modal de tareas se maneja desde el padre via @open-task-modal -->
    <!-- Modal ver detalle tarea completada (autosuficiente) -->
    <case-task-history ref="taskHistory"></case-task-history>
    <!-- Modal gestionar tarea pendiente (autosuficiente) -->
    <case-task-modal ref="taskModal" :current-user-id="userId" @completed="onTaskCompleted"></case-task-modal>
    <!-- Modal reasignar tarea -->
    <div v-if="reasignarTask" class="ct-modal-backdrop" @click.self="reasignarTask = null">
        <div class="ct-modal">
            <p class="ct-modal-title">Reasignar tarea</p>
            <div class="ct-modal-field">
                <label class="ct-modal-label">Seleccionar agente</label>
                <select v-model="reasignarAgente" class="ct-modal-input" style="min-height:auto;resize:none;padding:8px 12px">
                    <option value="" disabled>Seleccione un agente...</option>
                    <option v-for="ag in agentesEquipo" :key="ag.icode" :value="ag.icode">\${ ag.fullName || ag.icode }</option>
                </select>
            </div>
            <div class="ct-modal-footer">
                <button class="ct-modal-btn-cancel" @click="reasignarTask = null">Cancelar</button>
                <button class="ct-modal-btn-save" :disabled="!reasignarAgente || reasignando" @click="confirmarReasignar">
                    \${ reasignando ? 'Reasignando...' : 'Reasignar' }
                </button>
            </div>
        </div>
    </div>
</div>
    `
});
