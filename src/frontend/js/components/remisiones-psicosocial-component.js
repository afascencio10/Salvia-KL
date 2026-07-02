/**
 * remisiones-psicosocial-component
 * Listado de remisiones de Atención Psicosocial (datos mock — backend pendiente).
 *
 * Props: mostrarCards, defaultFilter, reasignacion, pageSize
 * Emite: ver-caso, ver-remision, reasignar-remisiones
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[remisiones-psicosocial-component] Cargue este script después de var app/home.');
        return;
    }

    var SESSION_DOTS = 6;
    var SESSION_LABEL = 'Sesiones (4 - 6)';

    var MOCK_DUPLAS = [
        { id: 'dupla-1', name: 'Dupla 1' },
        { id: 'dupla-2', name: 'Dupla 2' },
    ];

    var MOCK_PROFESSIONALS = [
        { value: 'AG-PSY-01', label: 'Alejandra Mora', team: 'psicologia' },
        { value: 'AG-TS-01', label: 'Valentina Ospina', team: 'trab. social' },
        { value: 'AG-PSY-02', label: 'Carolina Ruiz', team: 'psicologia' },
    ];

    var MOCK_REMISIONES = [
        {
            id: 'rem-001', caseICode: 'SAL-001', followUpId: 'fu-001',
            status: 'abierto', sessionCount: 0, type: 'derivacion',
            createdAt: '2026-06-10T09:00:00Z',
            submittedBy: 'AG-RO-01', submittedByName: 'Agente López', submittedByTeam: 'Riesgo alto',
            duplaId: null, duplaName: null, psychologistName: null, socialWorkerName: null,
            professionalId: null, professionalName: null, professionalTeam: null,
            victimNames: 'María', victimLastNames: 'González', docNumber: '1023456789',
            victimPhone: '7777-1234', municipality: 'San Salvador', riskLevel: 3,
        },
        {
            id: 'rem-002', caseICode: 'SAL-002', followUpId: 'fu-002',
            status: 'abierto', sessionCount: 0, type: 'derivacion',
            createdAt: '2026-06-11T11:30:00Z',
            submittedBy: 'AG-RO-02', submittedByName: 'Agente Martínez', submittedByTeam: 'Riesgo bajo',
            duplaId: null, duplaName: null, psychologistName: null, socialWorkerName: null,
            professionalId: null, professionalName: null, professionalTeam: null,
            victimNames: 'Ana', victimLastNames: 'Pérez', docNumber: '1098765432',
            victimPhone: '7888-5678', municipality: 'San Miguel', riskLevel: 2,
        },
        {
            id: 'rem-003', caseICode: 'SAL-003', followUpId: 'fu-003',
            status: 'en_gestion', sessionCount: 2, type: 'derivacion',
            createdAt: '2026-06-05T14:00:00Z',
            submittedBy: 'AG-RO-01', submittedByName: 'Agente López', submittedByTeam: 'Riesgo alto',
            duplaId: 'dupla-1', duplaName: 'Dupla 1',
            psychologistName: 'Alejandra Mora', socialWorkerName: 'Valentina Ospina',
            professionalId: null, professionalName: null, professionalTeam: null,
            victimNames: 'Lucía', victimLastNames: 'Torres', docNumber: '1033445566',
            victimPhone: '7999-1111', municipality: 'Santa Ana', riskLevel: 4,
        },
        {
            id: 'rem-004', caseICode: 'SAL-004', followUpId: 'fu-004',
            status: 'en_gestion', sessionCount: 3, type: 'derivacion',
            createdAt: '2026-06-06T08:15:00Z',
            submittedBy: 'AG-RO-03', submittedByName: 'Agente Ramírez', submittedByTeam: 'Hombres',
            duplaId: 'dupla-1', duplaName: 'Dupla 1',
            psychologistName: 'Alejandra Mora', socialWorkerName: 'Valentina Ospina',
            professionalId: null, professionalName: null, professionalTeam: null,
            victimNames: 'Patricia', victimLastNames: 'Vega', docNumber: '1044556677',
            victimPhone: '7666-2222', municipality: 'San Salvador', riskLevel: 3,
        },
        {
            id: 'rem-005', caseICode: 'SAL-005', followUpId: 'fu-005',
            status: 'en_gestion', sessionCount: 1, type: 'derivacion',
            createdAt: '2026-06-07T16:45:00Z',
            submittedBy: 'AG-RO-02', submittedByName: 'Agente Martínez', submittedByTeam: 'Riesgo bajo',
            duplaId: null, duplaName: null, psychologistName: null, socialWorkerName: null,
            professionalId: 'AG-PSY-02', professionalName: 'Carolina Ruiz', professionalTeam: 'psicologia',
            victimNames: 'Sofía', victimLastNames: 'Castro', docNumber: '1055667788',
            victimPhone: '7555-3333', municipality: 'Usulután', riskLevel: 1,
        },
        {
            id: 'rem-006', caseICode: 'SAL-006', followUpId: 'fu-006',
            status: 'en_devolucion', sessionCount: 1, type: 'derivacion',
            createdAt: '2026-06-08T10:00:00Z',
            submittedBy: 'AG-RO-01', submittedByName: 'Agente López', submittedByTeam: 'Riesgo alto',
            duplaId: null, duplaName: null, psychologistName: null, socialWorkerName: null,
            professionalId: 'AG-PSY-01', professionalName: 'Alejandra Mora', professionalTeam: 'psicologia',
            victimNames: 'Rosa', victimLastNames: 'Jiménez', docNumber: '1066112233',
            victimPhone: '7111-9999', municipality: 'San Salvador', riskLevel: 2,
        },
        {
            id: 'rem-007', caseICode: 'SAL-007', followUpId: 'fu-007',
            status: 'cerrado', sessionCount: 4, type: 'derivacion',
            createdAt: '2026-05-22T13:20:00Z',
            submittedBy: 'AG-RO-03', submittedByName: 'Agente Ramírez', submittedByTeam: 'Hombres',
            duplaId: 'dupla-2', duplaName: 'Dupla 2',
            psychologistName: 'Carolina Ruiz', socialWorkerName: 'Valentina Ospina',
            professionalId: null, professionalName: null, professionalTeam: null,
            victimNames: 'Daniela', victimLastNames: 'Herrera', docNumber: '1077889900',
            victimPhone: '7333-5555', municipality: 'La Libertad', riskLevel: 3,
        },
        {
            id: 'rem-008', caseICode: 'SAL-008', followUpId: 'fu-008',
            status: 'cerrado', sessionCount: 5, type: 'derivacion',
            createdAt: '2026-05-25T09:30:00Z',
            submittedBy: 'AG-RO-02', submittedByName: 'Agente Martínez', submittedByTeam: 'Riesgo bajo',
            duplaId: null, duplaName: null, psychologistName: null, socialWorkerName: null,
            professionalId: 'AG-TS-01', professionalName: 'Valentina Ospina', professionalTeam: 'trab. social',
            victimNames: 'Gabriela', victimLastNames: 'Silva', docNumber: '1088990011',
            victimPhone: '7222-6666', municipality: 'Sonsonate', riskLevel: 4,
        },
    ];

    vueApp.component('remisiones-psicosocial-component', {
        delimiters: ['${', '}'],
        template: window.__remisionesPsicosocialComponentTpl || '#tpl-remisiones-psicosocial-component',

        props: {
            mostrarCards:  { type: Boolean, default: false },
            defaultFilter: { type: Object,  default: function() { return {}; } },
            reasignacion:  { type: Boolean, default: false },
            pageSize:      { type: Number,  default: 20 },
        },

        emits: ['ver-caso', 'ver-remision', 'reasignar-remisiones'],

        data: function() {
            return {
                sessionDots: SESSION_DOTS,
                sessionLabel: SESSION_LABEL,
                scopeFilter: {},
                activeFilters: {},
                searchNumeroIdentidad: '',
                searchTelefono: '',
                autocompleteText: { profesional_asignada: '' },
                autocompleteSelected: { profesional_asignada: null },
                autocompleteSuggestions: [],
                duplaOptions: MOCK_DUPLAS.slice(),
                equipoRemitenteOptions: ['Riesgo alto', 'Riesgo bajo', 'Hombres'],
                currentPage: 1,
                selectedRemisiones: [],
                tableAlert: { visible: false, message: '' },
                _searchTimer: null,
                _autocompleteTimer: null,
                _tableAlertTimer: null,
                estadoRemisionOptions: [
                    { value: 'abierto', label: 'Abiertos' },
                    { value: 'en_gestion', label: 'En gestión' },
                    { value: 'en_devolucion', label: 'En devolución' },
                    { value: 'cerrado', label: 'Cerrados' },
                ],
                nivelRiesgoOptions: [
                    { value: 'bajo', label: 'Bajo' },
                    { value: 'moderado', label: 'Moderado' },
                    { value: 'alto', label: 'Alto' },
                    { value: 'extremo', label: 'Crítico' },
                ],
            };
        },

        computed: {
            sesionesOptions: function() {
                var opts = [];
                for (var i = 0; i <= SESSION_DOTS; i++) {
                    opts.push(i);
                }
                return opts;
            },

            scopedRemisiones: function() {
                return this.applyFilters(MOCK_REMISIONES);
            },

            filteredRemisiones: function() {
                return this.scopedRemisiones;
            },

            stats: function() {
                var list = this.filteredRemisiones;
                var counts = { abierto: 0, en_gestion: 0, en_devolucion: 0, cerrado: 0 };
                list.forEach(function(r) {
                    if (counts[r.status] !== undefined) {
                        counts[r.status]++;
                    }
                });
                return {
                    total: list.length,
                    abierto: counts.abierto,
                    enGestion: counts.en_gestion,
                    enDevolucion: counts.en_devolucion,
                    cerrado: counts.cerrado,
                };
            },

            totalPages: function() {
                return Math.max(1, Math.ceil(this.filteredRemisiones.length / this.pageSize));
            },

            paginatedRemisiones: function() {
                var start = (this.currentPage - 1) * this.pageSize;
                return this.filteredRemisiones.slice(start, start + this.pageSize);
            },

            isAllPageSelected: function() {
                var page = this.paginatedRemisiones;
                if (!page.length) {
                    return false;
                }
                var self = this;
                return page.every(function(r) { return self.isSelected(r); });
            },

            isPagePartiallySelected: function() {
                var page = this.paginatedRemisiones;
                if (!page.length) {
                    return false;
                }
                var self = this;
                var any = page.some(function(r) { return self.isSelected(r); });
                return any && !this.isAllPageSelected;
            },
        },

        mounted: function() {
            this.initFromProps();
        },

        methods: {
            initFromProps: function() {
                this.scopeFilter = Object.assign({}, this.defaultFilter || {});
                this.activeFilters = Object.assign({}, this.scopeFilter);
                if (this.scopeFilter.professional_id) {
                    var prof = MOCK_PROFESSIONALS.find(function(p) {
                        return p.value === this.scopeFilter.professional_id;
                    }.bind(this));
                    if (prof) {
                        this.autocompleteSelected.profesional_asignada = {
                            value: prof.value,
                            label: prof.label,
                        };
                    }
                }
                this.reload();
            },

            reload: function() {
                this.currentPage = 1;
                this.selectedRemisiones = [];
            },

            applyFilters: function(source) {
                var self = this;
                var f = this.activeFilters;
                return source.filter(function(r) {
                    if (f.professional_id && r.professionalId !== f.professional_id) {
                        return false;
                    }
                    if (f.dupla_id && r.duplaId !== f.dupla_id) {
                        return false;
                    }
                    if (f.estado_remision && r.status !== f.estado_remision) {
                        return false;
                    }
                    if (f.sesiones_completadas !== undefined && f.sesiones_completadas !== '' &&
                        String(r.sessionCount) !== String(f.sesiones_completadas)) {
                        return false;
                    }
                    if (f.equipo_remitente && r.submittedByTeam !== f.equipo_remitente) {
                        return false;
                    }
                    if (f.nivel_riesgo) {
                        var riskMap = { bajo: 1, moderado: 2, alto: 3, extremo: 4 };
                        if (r.riskLevel !== riskMap[f.nivel_riesgo]) {
                            return false;
                        }
                    }
                    if (self.searchNumeroIdentidad.trim()) {
                        var doc = (r.docNumber || '').toLowerCase();
                        if (doc.indexOf(self.searchNumeroIdentidad.trim().toLowerCase()) === -1) {
                            return false;
                        }
                    }
                    if (self.searchTelefono.trim()) {
                        var tel = (r.victimPhone || '').toLowerCase();
                        if (tel.indexOf(self.searchTelefono.trim().toLowerCase()) === -1) {
                            return false;
                        }
                    }
                    return true;
                });
            },

            setDropdownFilter: function(key, value) {
                if (value === '') {
                    delete this.activeFilters[key];
                } else {
                    this.activeFilters[key] = value;
                }
                this.onFilterChange();
            },

            onSearchInput: function() {
                var self = this;
                if (self._searchTimer) {
                    clearTimeout(self._searchTimer);
                }
                self._searchTimer = setTimeout(function() {
                    self.onFilterChange();
                }, 400);
            },

            onFilterChange: function() {
                this.currentPage = 1;
                this.selectedRemisiones = [];
            },

            clearAllFilters: function() {
                this.activeFilters = Object.assign({}, this.scopeFilter);
                this.searchNumeroIdentidad = '';
                this.searchTelefono = '';
                this.autocompleteText.profesional_asignada = '';
                if (!this.scopeFilter.professional_id) {
                    this.autocompleteSelected.profesional_asignada = null;
                }
                this.autocompleteSuggestions = [];
                this.onFilterChange();
            },

            onAutocompleteInput: function() {
                var self = this;
                if (self._autocompleteTimer) {
                    clearTimeout(self._autocompleteTimer);
                }
                self._autocompleteTimer = setTimeout(function() {
                    self.runAutocompleteSearch();
                }, 400);
            },

            runAutocompleteSearch: function() {
                var q = (this.autocompleteText.profesional_asignada || '').trim().toLowerCase();
                if (q.length < 3) {
                    this.autocompleteSuggestions = [];
                    return;
                }
                this.autocompleteSuggestions = MOCK_PROFESSIONALS.filter(function(p) {
                    return p.label.toLowerCase().indexOf(q) !== -1;
                });
            },

            selectAutocomplete: function(opt) {
                this.autocompleteSelected.profesional_asignada = opt;
                this.autocompleteText.profesional_asignada = '';
                this.autocompleteSuggestions = [];
                this.activeFilters.professional_id = opt.value;
                this.onFilterChange();
            },

            clearAutocomplete: function() {
                this.autocompleteSelected.profesional_asignada = null;
                delete this.activeFilters.professional_id;
                if (this.scopeFilter.professional_id) {
                    this.activeFilters.professional_id = this.scopeFilter.professional_id;
                    var prof = MOCK_PROFESSIONALS.find(function(p) {
                        return p.value === this.scopeFilter.professional_id;
                    }.bind(this));
                    if (prof) {
                        this.autocompleteSelected.profesional_asignada = {
                            value: prof.value,
                            label: prof.label,
                        };
                    }
                }
                this.onFilterChange();
            },

            changePage: function(page) {
                if (page < 1 || page > this.totalPages) {
                    return;
                }
                this.currentPage = page;
                this.selectedRemisiones = [];
                var el = this.$el && this.$el.closest('.rps-wrapper');
                if (el) {
                    el.scrollIntoView({ behavior: 'smooth', block: 'start' });
                }
            },

            isSelected: function(r) {
                return this.selectedRemisiones.some(function(s) { return s.id === r.id; });
            },

            toggleRemisionSelection: function(r, checked) {
                if (!this.reasignacion) {
                    return;
                }
                if (checked) {
                    if (this.selectedRemisiones.length > 0 &&
                        this.selectedRemisiones[0].status !== r.status) {
                        this.showTableAlert('Solo puedes seleccionar remisiones con el mismo estado');
                        return;
                    }
                    if (!this.isSelected(r)) {
                        this.selectedRemisiones.push(r);
                    }
                } else {
                    this.selectedRemisiones = this.selectedRemisiones.filter(function(s) {
                        return s.id !== r.id;
                    });
                }
            },

            toggleSelectAllPage: function(ev) {
                var checked = ev.target.checked;
                var page = this.paginatedRemisiones;
                if (!checked) {
                    var pageIds = page.map(function(r) { return r.id; });
                    this.selectedRemisiones = this.selectedRemisiones.filter(function(s) {
                        return pageIds.indexOf(s.id) === -1;
                    });
                    return;
                }
                var targetStatus = this.selectedRemisiones.length
                    ? this.selectedRemisiones[0].status
                    : (page[0] ? page[0].status : null);
                var mixed = page.some(function(r) { return r.status !== targetStatus; });
                if (mixed && !this.selectedRemisiones.length) {
                    this.showTableAlert('Solo puedes seleccionar remisiones con el mismo estado');
                    ev.target.checked = false;
                    return;
                }
                var self = this;
                page.forEach(function(r) {
                    if (r.status === targetStatus && !self.isSelected(r)) {
                        self.selectedRemisiones.push(r);
                    }
                });
            },

            showTableAlert: function(msg) {
                var self = this;
                self.tableAlert = { visible: true, message: msg };
                if (self._tableAlertTimer) {
                    clearTimeout(self._tableAlertTimer);
                }
                self._tableAlertTimer = setTimeout(function() {
                    self.hideTableAlert();
                }, 4500);
            },

            hideTableAlert: function() {
                this.tableAlert.visible = false;
            },

            emitReasignar: function() {
                this.$emit('reasignar-remisiones', { remisiones: this.selectedRemisiones.slice() });
            },

            emitVerCaso: function(r) {
                this.$emit('ver-caso', { caseICode: r.caseICode, remision: r });
            },

            emitVerRemision: function(r) {
                this.$emit('ver-remision', {
                    remisionId: r.id,
                    followUpId: r.followUpId,
                    remision: r,
                });
            },

            riskLabel: function(level) {
                var map = { 1: 'Bajo', 2: 'Medio', 3: 'Alto', 4: 'Crítico' };
                return map[level] || '—';
            },

            riskBadgeClass: function(level) {
                var map = { 1: 'bajo', 2: 'moderado', 3: 'alto', 4: 'extremo' };
                return map[level] || 'bajo';
            },

            statusLabel: function(status) {
                var map = {
                    abierto: 'Abiertos',
                    en_gestion: 'En gestión',
                    en_devolucion: 'En devolución',
                    cerrado: 'Cerrados',
                };
                return map[status] || status;
            },

            statusBadgeClass: function(status) {
                return status || '';
            },

            professionalRoleLabel: function(team) {
                if (team === 'psicologia') {
                    return 'Psicóloga';
                }
                if (team === 'trab. social') {
                    return 'Trab. Social';
                }
                return team || '—';
            },

            formatDate: function(iso) {
                if (!iso) {
                    return '—';
                }
                return iso.substring(0, 10);
            },

            initials: function(name) {
                if (!name) {
                    return '?';
                }
                var parts = name.trim().split(/\s+/);
                return parts.map(function(p) { return p.charAt(0); }).join('').substring(0, 2).toUpperCase();
            },
        },
    });
})();
