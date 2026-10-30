/**
 * remisiones-medidas-component
 * Listado de remisiones de Medidas de Emergencia.
 * La prop `vista` decide cards, filtros, recorte y acciones.
 * `puedeReasignar` la encienden los roles sm y sv.
 * El recorte real lo aplica el API según el rol de la sesión.
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[remisiones-medidas-component] Cargue este script después de var app/home.');
        return;
    }

    var RESULTADOS_REVISION = [
        { value: 'aprobado_para_caracterizacion_me', label: 'Aprobado para caracterización ME' },
        { value: 'devuelto_para_ajuste', label: 'Devuelto para ajuste' },
        { value: 'devuelto_para_acompanamiento', label: 'Devuelto para acompañamiento' },
    ];

    var ESTADOS = [
        { value: 'riesgo_inminente_por_decidir_via', label: 'Riesgo inminente: por decidir vía' },
        { value: 'riesgo_inminente_por_escalar_cruz_roja', label: 'Riesgo inminente: por escalar a Cruz Roja' },
        { value: 'riesgo_inminente_por_entregar_medida', label: 'Riesgo inminente: por entregar medida' },
        { value: 'riesgo_inminente_por_confirmar_seguridad', label: 'Riesgo inminente: por confirmar seguridad' },
        { value: 'pendiente_de_revision', label: 'Pendiente de Revisión' },
        { value: 'por_asignar', label: 'Por asignar' },
        { value: 'asignada', label: 'Asignada' },
        { value: 'en_gestion_documental', label: 'En gestión documental' },
        { value: 'caracterizado', label: 'Caracterizado' },
        { value: 'por_revisar_trazabilidad', label: 'Por revisar trazabilidad' },
        { value: 'por_cotizar', label: 'Por cotizar' },
        { value: 'por_aprobar_cotizaciones', label: 'Por aprobar cotizaciones' },
        { value: 'aprobadas', label: 'Aprobadas' },
        { value: 'en_seguimiento', label: 'En Seguimiento' },
        { value: 'devuelta_para_justificacion', label: 'Devuelta para justificación' },
        { value: 'devuelta_a_seguimiento', label: 'Devuelta a seguimiento' },
        { value: 'descartada', label: 'Descartada' },
        { value: 'cerrada_por_no_contacto', label: 'Cerrada por no contacto' },
    ];

    var ESTADOS_AGENTE = [
        'riesgo_inminente_por_confirmar_seguridad',
        'asignada',
        'en_gestion_documental',
        'aprobadas',
        'en_seguimiento',
    ];

    var ESTADO_MOD = {
        'Riesgo inminente: por decidir vía': 'decidir',
        'Riesgo inminente: por escalar a Cruz Roja': 'escalar',
        'Riesgo inminente: por entregar medida': 'entregar',
        'Riesgo inminente: por confirmar seguridad': 'confirmar',
        'Pendiente de Revisión': 'revision',
        'Por asignar': 'por-asignar',
        'Asignada': 'asignada',
        'En gestión documental': 'gestion',
        'Caracterizado': 'caracterizado',
        'Por revisar trazabilidad': 'trazabilidad',
        'Por cotizar': 'cotizar',
        'Por aprobar cotizaciones': 'cotizaciones',
        'Aprobadas': 'aprobadas',
        'En Seguimiento': 'seguimiento',
        'Devuelta para justificación': 'justificacion',
        'Devuelta a seguimiento': 'devuelta',
    };

    var TONO_CARD = {
        total: 'total',
        riesgo: 'riesgo',
        revision: 'revision',
        porAsignar: 'asignar',
        asignadas: 'asignadas',
        gestion: 'gestion',
        caracterizado: 'caracterizado',
        aprobadas: 'aprobadas',
        seguimiento: 'seguimiento',
    };

    var ICONO_BANNER = {
        equipo: '⚖️',
        convenio: '📂',
        'cruz-roja': '🤝',
    };

    function estadosDe(values) {
        return ESTADOS.filter(function(estado) {
            return values.indexOf(estado.value) !== -1;
        });
    }

    function aFila(item) {
        return {
            id: item.id,
            caso: {
                id: item.caseICode || '',
                nombre: item.caseName || '',
                docNumero: item.docNumber || '',
                celular: item.victimPhone || '',
                nivelRiesgo: item.riskLabel || '',
                municipio: item.municipality || '',
            },
            remitidoPor: { nombre: item.submittedByName || '' },
            fechaRemision: item.referredAt || '',
            estado: item.statusLabel || '',
            medidasSolicitadas: item.measures || [],
            agenteAsignado: item.assignedAgent || null,
            urgencia: item.urgencyLabel || '',
            revisionTecnica: item.technicalReview ? {
                profesional: item.technicalReview.professionalName || '',
                concepto: item.technicalReview.concept || '',
                resultado: item.technicalReview.resultLabel || '',
                fecha: item.technicalReview.date || '',
            } : null,
        };
    }

    vueApp.component('remisiones-medidas-component', {
        delimiters: ['${', '}'],
        template: window.__remisionesMedidasComponentTpl || '#tpl-remisiones-medidas-component',

        props: {
            vista: { type: String, default: 'supervisor' },
            agenteId: { type: String, default: '' },
            puedeReasignar: { type: Boolean, default: false },
            userId: { type: String, default: '' },
            userRole: { type: String, default: '' },
        },

        emits: ['ver-caso', 'ver-detalle'],

        data: function() {
            return {
                agentes: [],
                profesionales: [],
                resultadosRevision: RESULTADOS_REVISION,
                remisiones: [],
                alcanceTotal: 0,
                currentPage: 1,
                pageSize: 10,
                totalFiltrado: 0,
                cards: [],
                banner: null,
                loading: true,
                loadError: '',
                searchText: '',
                searchTimer: null,
                reqSeq: 0,
                filtroEstado: '',
                filtroAgente: '',
                filtroUrgencia: '',
                selectedIds: [],
                reasignarAbierto: false,
                reasignarAgenteId: '',
                reasignarError: '',
                revisionAbierta: false,
                revisionItem: null,
                revisionProfesional: '',
                revisionConcepto: '',
                revisionResultado: '',
                revisionError: '',
                accionError: '',
                abriendoRevisionId: '',
                abriendoAsignacionId: '',
                panelTareas: null,
            };
        },

        computed: {
            esAgente: function() { return this.vista === 'agente'; },
            muestraFiltroEstado: function() {
                return this.vista === 'supervisor' || this.vista === 'agente';
            },
            muestraFiltroAgente: function() { return this.vista === 'supervisor'; },
            placeholderBusqueda: function() {
                return this.esAgente ? 'Nombre o número de identidad' : 'Nombre o doc. identidad';
            },
            labelLimpiar: function() {
                return this.esAgente ? 'Limpiar' : 'Limpiar filtros';
            },
            estadosFiltro: function() {
                if (this.userRole === 'ro') return ESTADOS;
                return this.esAgente ? estadosDe(ESTADOS_AGENTE) : ESTADOS;
            },
            totalPages: function() {
                var pages = Math.ceil((this.totalFiltrado || 0) / (this.pageSize || 10));
                return pages > 0 ? pages : 1;
            },
            claseGrilla: function() {
                return this.puedeReasignar ? 'rme-grid--check' : 'rme-grid--3';
            },
            visibles: function() {
                return this.remisiones;
            },
            todosSeleccionados: function() {
                var self = this;
                return this.visibles.length > 0 && this.visibles.every(function(r) {
                    return self.selectedIds.indexOf(r.id) !== -1;
                });
            },
            emptyMessage: function() {
                if (this.loadError) return this.loadError;
                if (this.loading) return '';
                if (this.esAgente && this.alcanceTotal === 0) {
                    if (this.userRole === 'ro') return 'No has remitido casos.';
                    return 'No tienes remisiones asignadas.';
                }
                return 'No hay remisiones que coincidan con los filtros.';
            },
            revisionLista: function() {
                return !!(this.revisionProfesional && this.revisionConcepto.trim() && this.revisionResultado);
            },
        },

        mounted: function() {
            this.cargar();
            if (this.muestraFiltroAgente || this.puedeReasignar) {
                this.cargarAgentes();
            }
        },

        methods: {
            inicial: function(nombre) {
                return nombre ? nombre.charAt(0) : '';
            },
            claseRiesgo: function(nivel) {
                if (nivel === 'Extremo' || nivel === 'Alto (feminicidio)' || nivel === 'Crítico') return 'rme-risk--critico';
                if (nivel === 'Alto') return 'rme-risk--alto';
                return 'rme-risk--medio';
            },
            claseUrgencia: function(urgencia) {
                return urgencia === 'Riesgo inminente' ? 'rme-urgencia--inminente' : 'rme-urgencia--normal';
            },
            claseEstado: function(estado) {
                return 'rme-estado--' + (ESTADO_MOD[estado] || 'default');
            },
            isSelected: function(id) {
                return this.selectedIds.indexOf(id) !== -1;
            },
            toggleOne: function(id) {
                var idx = this.selectedIds.indexOf(id);
                if (idx === -1) this.selectedIds.push(id);
                else this.selectedIds.splice(idx, 1);
            },
            toggleAll: function() {
                if (this.todosSeleccionados) {
                    this.selectedIds = [];
                    return;
                }
                this.selectedIds = this.visibles.map(function(r) { return r.id; });
            },
            onSearchInput: function() {
                var self = this;
                clearTimeout(this.searchTimer);
                this.reqSeq += 1;
                this.loading = true;
                this.loadError = '';
                this.searchTimer = setTimeout(function() { self.cargarDesdeInicio(); }, 400);
            },
            limpiarFiltros: function() {
                clearTimeout(this.searchTimer);
                this.searchText = '';
                this.filtroEstado = '';
                this.filtroAgente = '';
                this.filtroUrgencia = '';
                this.cargarDesdeInicio();
            },
            cargarDesdeInicio: function() {
                this.currentPage = 1;
                this.selectedIds = [];
                this.cargar();
            },
            changePage: function(newPage) {
                if (newPage < 1 || newPage > this.totalPages || newPage === this.currentPage) return;
                this.currentPage = newPage;
                this.selectedIds = [];
                this.cargar();
            },
            cargar: function() {
                var self = this;
                this.reqSeq += 1;
                var seq = this.reqSeq;
                var params = new URLSearchParams();
                if (this.searchText.trim()) params.set('q', this.searchText.trim());
                if (this.filtroEstado) params.set('status', this.filtroEstado);
                if (this.filtroUrgencia) params.set('urgency', this.filtroUrgencia);
                if (this.vista === 'supervisor' && this.filtroAgente) params.set('agent_id', this.filtroAgente);
                if (this.vista === 'agente' && this.agenteId) params.set('agent_id', this.agenteId);
                params.set('page', String(this.currentPage || 1));
                this.loading = true;
                this.loadError = '';
                fetch('/api/v1/emergency-measure-referrals?' + params.toString(), { credentials: 'same-origin' })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function(data) {
                        if (seq !== self.reqSeq) return;
                        self.alcanceTotal = data.alcance || 0;
                        self.totalFiltrado = data.total || 0;
                        self.pageSize = data.pageSize || 10;
                        self.currentPage = data.page || self.currentPage;
                        self.remisiones = (data.items || []).map(aFila);
                        self.cards = (data.cards || []).map(function(card) {
                            return {
                                key: card.key,
                                label: card.label,
                                valor: card.value,
                                tono: TONO_CARD[card.key] || 'total',
                            };
                        });
                        self.banner = data.banner ? {
                            tono: data.banner.tono,
                            texto: data.banner.texto,
                            icono: ICONO_BANNER[data.banner.tono] || '',
                        } : null;
                        self.loading = false;
                    })
                    .catch(function(err) {
                        if (seq !== self.reqSeq) return;
                        self.remisiones = [];
                        self.alcanceTotal = 0;
                        self.totalFiltrado = 0;
                        self.cards = [];
                        self.banner = null;
                        self.loading = false;
                        self.loadError = err.message || 'No se pudo cargar el listado.';
                    });
            },
            cargarAgentes: function() {
                var self = this;
                fetch('/api/v1/emergency-measure-referrals/agentes', { credentials: 'same-origin' })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function(data) { self.agentes = data.items || []; })
                    .catch(function() { self.agentes = []; });
            },
            cargarProfesionales: function() {
                var self = this;
                if (this.profesionales.length) return Promise.resolve();
                return fetch('/api/v1/emergency-measure-referrals/profesionales-revision', { credentials: 'same-origin' })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function(data) { self.profesionales = data.items || []; });
            },
            leerJSON: function(res) {
                return res.json().catch(function() { return {}; }).then(function(body) {
                    if (!res.ok) {
                        var mensaje = (body && body.error) ? body.error : 'No se pudo completar la operación.';
                        throw new Error(mensaje);
                    }
                    return body || {};
                });
            },
            accionDe: function(item) {
                var estado = item.estado;
                var vista = this.vista;
                if (vista === 'supervisor' && estado === 'Riesgo inminente: por decidir vía') {
                    return { tipo: 'tareas', label: '🔥 Decidir vía de atención', clase: 'rme-action--decidir rme-action--pulse' };
                }
                if (vista === 'supervision-convenio' && estado === 'Riesgo inminente: por escalar a Cruz Roja') {
                    return { tipo: 'detalle', label: '🔥 Escalar a Cruz Roja', clase: 'rme-action--escalar' };
                }
                if (vista === 'cruz-roja' && estado === 'Riesgo inminente: por entregar medida') {
                    return { tipo: 'detalle', label: '🔥 Entregar medida inmediata', clase: 'rme-action--entregar' };
                }
                if ((vista === 'supervisor' || vista === 'equipo-tecnico') && estado === 'Pendiente de Revisión') {
                    return { tipo: 'revisar', label: '⚖️ Revisar (Equipo Técnico)', clase: 'rme-action--revisar' };
                }
                if (this.puedeReasignar && estado === 'Por asignar') {
                    return { tipo: 'asignar', label: '👤 Asignar agente', clase: 'rme-action--asignar' };
                }
                if (vista === 'supervisor' && estado === 'Caracterizado') {
                    return { tipo: 'detalle', label: '📄 Revisar caracterización', clase: 'rme-action--caracterizacion' };
                }
                if (vista === 'supervision-convenio' && estado === 'Por revisar trazabilidad') {
                    return { tipo: 'detalle', label: '📂 Revisar trazabilidad', clase: 'rme-action--trazabilidad' };
                }
                if (vista === 'supervision-convenio' && estado === 'Por aprobar cotizaciones') {
                    return { tipo: 'detalle', label: '✅ Aprobar cotizaciones', clase: 'rme-action--cotizaciones' };
                }
                if (vista === 'cruz-roja' && estado === 'Por cotizar') {
                    return { tipo: 'detalle', label: '💰 Cargar cotizaciones', clase: 'rme-action--cotizar' };
                }
                if (vista === 'agente' && estado === 'Riesgo inminente: por confirmar seguridad') {
                    return { tipo: 'detalle', label: '🔥 Confirmar seguridad', clase: 'rme-action--confirmar' };
                }
                return null;
            },
            onAccion: function(item) {
                var accion = this.accionDe(item);
                if (!accion) return;
                if (accion.tipo === 'revisar' && this.vista === 'equipo-tecnico') this.abrirRevisionTarea(item);
                else if (accion.tipo === 'revisar') this.abrirRevision(item);
                else if (accion.tipo === 'asignar') this.abrirAsignacionTarea(item);
                else if (accion.tipo === 'tareas') this.abrirTareas(item);
                else this.onVerDetalle(item);
            },
            onVerCaso: function(item) {
                this.$emit('ver-caso', { caseICode: item.caso.id, remisionId: item.id });
            },
            onVerDetalle: function(item) {
                this.$emit('ver-detalle', { remisionId: item.id, caseICode: item.caso.id });
            },
            abrirTareas: function(item) {
                this.panelTareas = { caseICode: item.caso.id, remisionId: item.id, nombre: item.caso.nombre };
            },
            cerrarTareas: function() {
                this.panelTareas = null;
            },
            abrirReasignarMasivo: function() {
                this.reasignarAgenteId = '';
                this.reasignarError = '';
                this.reasignarAbierto = true;
                this.cargarAgentes();
            },
            abrirAsignacionTarea: function(item) {
                var self = this;
                if (!item || this.abriendoAsignacionId) return;
                this.accionError = '';
                this.abriendoAsignacionId = item.id;
                fetch('/api/v1/case-tasks?referralId=' + encodeURIComponent(item.id), { credentials: 'same-origin' })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function(tareas) {
                        var lista = Array.isArray(tareas) ? tareas : [];
                        var tarea = lista.find(function(t) {
                            return t.type === 'asignar_agente_medidas' && t.status === 'ToDo';
                        });
                        if (!tarea || !self.$refs.taskModal) {
                            self.accionError = 'Esta remisión no tiene una asignación pendiente.';
                            return;
                        }
                        self.$refs.taskModal.open(tarea.id);
                    })
                    .catch(function() {
                        self.accionError = 'No se pudo abrir la asignación.';
                    })
                    .finally(function() { self.abriendoAsignacionId = ''; });
            },
            cerrarReasignar: function() {
                this.reasignarAbierto = false;
                this.reasignarAgenteId = '';
                this.reasignarError = '';
            },
            confirmarReasignar: function() {
                var self = this;
                if (!this.reasignarAgenteId || !this.selectedIds.length) return;
                this.reasignarError = '';
                fetch('/api/v1/emergency-measure-referrals/reasignar', {
                    method: 'POST',
                    credentials: 'same-origin',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ ids: this.selectedIds.slice(), agentId: this.reasignarAgenteId }),
                })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function() {
                        self.selectedIds = [];
                        self.cerrarReasignar();
                        self.cargar();
                    })
                    .catch(function(err) {
                        self.reasignarError = err.message || 'No se pudo reasignar.';
                    });
            },
            abrirRevisionTarea: function(item) {
                var self = this;
                if (!item || this.abriendoRevisionId) return;
                this.accionError = '';
                this.abriendoRevisionId = item.id;
                fetch('/api/v1/case-tasks?referralId=' + encodeURIComponent(item.id), { credentials: 'same-origin' })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function(tareas) {
                        var lista = Array.isArray(tareas) ? tareas : [];
                        var tarea = lista.find(function(t) {
                            return t.type === 'revision_equipo_tecnico' && t.status === 'ToDo';
                        });
                        if (!tarea || !self.$refs.taskModal) {
                            self.accionError = 'Esta remisión no tiene una revisión pendiente.';
                            return;
                        }
                        self.$refs.taskModal.open(tarea.id);
                    })
                    .catch(function() {
                        self.accionError = 'No se pudo abrir la revisión.';
                    })
                    .finally(function() { self.abriendoRevisionId = ''; });
            },
            abrirRevision: function(item) {
                this.revisionItem = item;
                this.revisionProfesional = '';
                this.revisionConcepto = '';
                this.revisionResultado = '';
                this.revisionError = '';
                this.revisionAbierta = true;
                this.cargarProfesionales();
            },
            cerrarRevision: function() {
                this.revisionAbierta = false;
                this.revisionItem = null;
                this.revisionProfesional = '';
                this.revisionConcepto = '';
                this.revisionResultado = '';
                this.revisionError = '';
            },
            confirmarRevision: function() {
                var self = this;
                if (!this.revisionLista || !this.revisionItem) return;
                this.revisionError = '';
                fetch('/api/v1/emergency-measure-referrals/' + encodeURIComponent(this.revisionItem.id) + '/revision', {
                    method: 'POST',
                    credentials: 'same-origin',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        professionalId: this.revisionProfesional,
                        concepto: this.revisionConcepto.trim(),
                        resultado: this.revisionResultado,
                    }),
                })
                    .then(function(res) { return self.leerJSON(res); })
                    .then(function() {
                        self.cerrarRevision();
                        self.cargar();
                    })
                    .catch(function(err) {
                        self.revisionError = err.message || 'No se pudo guardar la revisión.';
                    });
            },
        },
    });
})();
