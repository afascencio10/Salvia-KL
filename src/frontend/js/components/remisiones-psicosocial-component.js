/**
 * remisiones-psicosocial-component
 * Listado de remisiones de Atención Psicosocial (E-01 — carga desde API).
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[remisiones-psicosocial-component] Cargue este script después de var app/home.');
        return;
    }

    var SESSION_DOTS = 6;
    var SESSION_LABEL = 'Sesiones (4 - 6)';

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
                autocompleteLoading: false,
                autocompleteError: false,
                duplaOptions: [],
                equipoRemitenteOptions: [],
                stats: {
                    total: 0,
                    abierto: 0,
                    enGestion: 0,
                    enDevolucion: 0,
                    cerrado: 0,
                },
                currentPage: 1,
                totalRemisiones: 0,
                remisiones: [],
                loading: true,
                loadError: null,
                selectedRemisiones: [],
                tableAlert: { visible: false, message: '' },
                _searchTimerNumeroIdentidad: null,
                _searchTimerTelefono: null,
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
                    { value: 'extremo', label: 'Extremo' },
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

            filteredRemisiones: function() {
                return this.remisiones;
            },

            totalPages: function() {
                return Math.max(1, Math.ceil(this.totalRemisiones / this.pageSize));
            },

            paginatedRemisiones: function() {
                return this.remisiones;
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
            this.loadAll();
        },

        methods: {
            initFromProps: function() {
                this.scopeFilter = Object.assign({}, this.defaultFilter || {});
                this.activeFilters = Object.assign({}, this.scopeFilter);
                if (this.scopeFilter.professional_id) {
                    this.activeFilters.professional_id = this.scopeFilter.professional_id;
                }
                if (this.scopeFilter.dupla_id) {
                    this.activeFilters.dupla_id = this.scopeFilter.dupla_id;
                }
            },

            buildQueryParams: function() {
                var params = new URLSearchParams();
                params.set('page', String(this.currentPage));
                params.set('page_size', String(this.pageSize));
                params.set('sort', 'created_at');
                params.set('order', 'desc');

                var f = this.activeFilters;
                if (f.professional_id) {
                    params.set('filter_professional_id', f.professional_id);
                }
                if (f.dupla_id) {
                    params.set('filter_dupla_id', f.dupla_id);
                }
                if (f.estado_remision) {
                    params.set('filter_estado_remision', f.estado_remision);
                }
                if (f.sesiones_completadas !== undefined && f.sesiones_completadas !== '') {
                    params.set('filter_sesiones_completadas', String(f.sesiones_completadas));
                }
                if (f.equipo_remitente) {
                    params.set('filter_equipo_remitente', f.equipo_remitente);
                }
                if (f.nivel_riesgo) {
                    params.set('filter_nivel_riesgo', f.nivel_riesgo);
                }
                if (f.numero_identidad) {
                    params.set('filter_numero_identidad', f.numero_identidad);
                }
                if (f.telefono) {
                    params.set('filter_telefono', f.telefono);
                }
                return params;
            },

            loadAll: function() {
                var self = this;
                self.loading = true;
                self.loadError = null;

                var tasks = [
                    self.fetchDuplas(),
                    self.fetchEquiposRemitentes(),
                    self.fetchRemisiones(false),
                ];
                if (self.mostrarCards) {
                    tasks.push(self.fetchStats());
                }

                Promise.all(tasks)
                    .catch(function(err) {
                        self.loadError = (err && err.message) || 'Error al cargar las remisiones';
                    })
                    .finally(function() {
                        self.loading = false;
                    });
            },

            reload: function() {
                this.currentPage = 1;
                this.selectedRemisiones = [];
                this.loadAll();
            },

            refreshData: function() {
                var self = this;
                self.loading = true;
                self.loadError = null;
                self.remisiones = [];
                var tasks = [self.fetchRemisiones(false)];
                if (self.mostrarCards) {
                    tasks.push(self.fetchStats());
                }
                Promise.all(tasks)
                    .catch(function(err) {
                        self.loadError = (err && err.message) || 'Error al cargar las remisiones';
                    })
                    .finally(function() {
                        self.loading = false;
                    });
            },

            fetchRemisiones: function(scrollToTop) {
                var self = this;
                var params = self.buildQueryParams();
                return fetch('/api/v1/psychosocial-support/list?' + params.toString())
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            throw new Error((result.data && result.data.error) || 'Error al cargar las remisiones');
                        }
                        self.remisiones = result.data.remisiones || [];
                        self.totalRemisiones = result.data.total || 0;
                        if (scrollToTop && self.$el) {
                            var el = self.$el.closest('.rps-wrapper');
                            if (el) {
                                el.scrollIntoView({ behavior: 'smooth', block: 'start' });
                            }
                        }
                    });
            },

            fetchStats: function() {
                var self = this;
                var params = self.buildQueryParams();
                params.delete('page');
                params.delete('page_size');
                return fetch('/api/v1/psychosocial-support/stats?' + params.toString())
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            throw new Error((result.data && result.data.error) || 'Error al cargar estadísticas');
                        }
                        self.stats = {
                            total: result.data.total || 0,
                            abierto: result.data.abierto || 0,
                            enGestion: result.data.enGestion || 0,
                            enDevolucion: result.data.enDevolucion || 0,
                            cerrado: result.data.cerrado || 0,
                        };
                    });
            },

            fetchDuplas: function() {
                var self = this;
                return fetch('/api/v1/duplas')
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            throw new Error((result.data && result.data.error) || 'Error al cargar duplas');
                        }
                        self.duplaOptions = result.data.duplas || [];
                    });
            },

            fetchEquiposRemitentes: function() {
                var self = this;
                return fetch('/api/v1/psychosocial-support/equipos-remitentes')
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            throw new Error((result.data && result.data.error) || 'Error al cargar equipos');
                        }
                        self.equipoRemitenteOptions = result.data.teams || [];
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

            onSearchInput: function(fieldKey) {
                var self = this;
                var timerProp = fieldKey === 'telefono'
                    ? '_searchTimerTelefono'
                    : '_searchTimerNumeroIdentidad';
                if (self[timerProp]) {
                    clearTimeout(self[timerProp]);
                }
                self[timerProp] = setTimeout(function() {
                    self.applySearchFilter(fieldKey);
                }, 400);
            },

            applySearchFilter: function(fieldKey) {
                var queryText = '';
                if (fieldKey === 'numero_identidad') {
                    queryText = (this.searchNumeroIdentidad || '').trim();
                } else if (fieldKey === 'telefono') {
                    queryText = (this.searchTelefono || '').trim();
                }
                if (queryText === '') {
                    delete this.activeFilters[fieldKey];
                } else {
                    this.activeFilters[fieldKey] = queryText;
                }
                this.onFilterChange();
            },

            onFilterChange: function() {
                this.currentPage = 1;
                this.selectedRemisiones = [];
                this.refreshData();
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
                this.autocompleteLoading = false;
                this.autocompleteError = false;
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
                var q = (this.autocompleteText.profesional_asignada || '').trim();
                if (q === '') {
                    this.autocompleteSuggestions = [];
                    this.autocompleteLoading = false;
                    this.autocompleteError = false;
                    return;
                }
                if (q.length < 3) {
                    this.autocompleteSuggestions = [];
                    this.autocompleteLoading = false;
                    this.autocompleteError = false;
                    return;
                }

                var self = this;
                self.autocompleteLoading = true;
                self.autocompleteError = false;
                self.autocompleteSuggestions = [];

                fetch('/api/v1/agents/search-psicosocial?q=' + encodeURIComponent(q) + '&limit=10')
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            throw new Error((result.data && result.data.error) || 'Error al buscar profesionales');
                        }
                        var agents = result.data.agents || [];
                        self.autocompleteSuggestions = agents.map(function(a) {
                            var label = ((a.names || '') + ' ' + (a.lastNames || '')).trim();
                            return {
                                value: a.icode,
                                label: label,
                                team: a.team || '',
                                specialty: a.specialty || '',
                                roleLabel: a.roleLabel || self.professionalRoleLabel(a.specialty || a.team),
                            };
                        });
                    })
                    .catch(function() {
                        self.autocompleteError = true;
                        self.autocompleteSuggestions = [];
                    })
                    .finally(function() {
                        self.autocompleteLoading = false;
                    });
            },

            autocompleteDisplayLabel: function(opt) {
                if (!opt) {
                    return '';
                }
                var role = opt.roleLabel || this.professionalRoleLabel(opt.specialty || opt.team);
                return role ? opt.label + ' - ' + role : opt.label;
            },

            autocompleteShowDropdown: function() {
                if (this.autocompleteSelected.profesional_asignada) {
                    return false;
                }
                if (this.autocompleteLoading) {
                    return true;
                }
                if (this.autocompleteSuggestions.length > 0) {
                    return true;
                }
                var text = (this.autocompleteText.profesional_asignada || '').trim();
                return text.length >= 3 && !this.autocompleteError;
            },

            autocompleteShowNoResults: function() {
                var text = (this.autocompleteText.profesional_asignada || '').trim();
                return text.length >= 3 &&
                    !this.autocompleteLoading &&
                    !this.autocompleteError &&
                    this.autocompleteSuggestions.length === 0;
            },

            selectAutocomplete: function(opt) {
                this.autocompleteSelected.profesional_asignada = {
                    value: opt.value,
                    label: opt.label,
                    team: opt.team || '',
                    specialty: opt.specialty || '',
                    roleLabel: opt.roleLabel || this.professionalRoleLabel(opt.specialty || opt.team),
                };
                this.autocompleteText.profesional_asignada = '';
                this.autocompleteSuggestions = [];
                this.autocompleteError = false;
                this.activeFilters.professional_id = opt.value;
                this.onFilterChange();
            },

            clearAutocomplete: function() {
                this.autocompleteSelected.profesional_asignada = null;
                this.autocompleteText.profesional_asignada = '';
                this.autocompleteSuggestions = [];
                this.autocompleteError = false;
                delete this.activeFilters.professional_id;
                if (this.scopeFilter.professional_id) {
                    this.activeFilters.professional_id = this.scopeFilter.professional_id;
                }
                this.onFilterChange();
            },

            changePage: function(page) {
                if (page < 1 || page > this.totalPages || page === this.currentPage) {
                    return;
                }
                this.currentPage = page;
                this.selectedRemisiones = [];
                this.remisiones = [];
                var self = this;
                self.loading = true;
                self.loadError = null;
                var tasks = [self.fetchRemisiones(true)];
                if (self.mostrarCards) {
                    tasks.push(self.fetchStats());
                }
                Promise.all(tasks)
                    .catch(function(err) {
                        self.loadError = (err && err.message) || 'Error al cargar las remisiones';
                    })
                    .finally(function() {
                        self.loading = false;
                    });
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

            riskStatusFromLevel: function(level) {
                var map = { 1: 'bajo', 2: 'moderado', 3: 'alto', 4: 'extremo' };
                return map[level] || 'desconocido';
            },

            riskBadgeLabel: function(riskStatus) {
                var map = {
                    extremo: 'Extremo',
                    alto: 'Alto',
                    moderado: 'Moderado',
                    medio: 'Moderado',
                    bajo: 'Bajo',
                    desconocido: 'Desconocido',
                };
                return map[riskStatus] || 'Desconocido';
            },

            riskLabel: function(level) {
                return this.riskBadgeLabel(this.riskStatusFromLevel(level));
            },

            riskBadgeClass: function(level) {
                var slug = this.riskStatusFromLevel(level);
                var map = {
                    extremo: 'extremo',
                    alto: 'alto',
                    moderado: 'moderado',
                    medio: 'moderado',
                    bajo: 'bajo',
                    desconocido: 'desconocido',
                };
                return map[slug] || 'desconocido';
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
                if (team === 'psicosocial') {
                    return 'Profesional';
                }
                return team || 'Profesional';
            },

            formatDate: function(iso) {
                if (!iso) {
                    return '—';
                }
                return String(iso).substring(0, 10);
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
