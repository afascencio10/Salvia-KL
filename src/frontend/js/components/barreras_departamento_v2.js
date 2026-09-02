/**
 * barreras_departamento_v2.js
 * Pantalla "Barreras Departamento" (barreras_departamento_v2.html).
 *
 * Requiere window.__BarrerasDepartamentoConfig:
 *   - windowTitle, currentUser, currentRole, currentUserId, assignedDepartmentId
 */

(function () {
    'use strict';

    var PAGE_SIZE = 10;
    var FILTER_DEBOUNCE_MS = 350;

    var cfg = window.__BarrerasDepartamentoConfig || {};

    var riskLevelLabels = { 1: 'Bajo', 2: 'Moderado', 3: 'Alto', 4: 'Extremo' };
    var estadoLabels = {
        OPEN: 'Abierta',
        'En Gestion': 'En gestión',
        Articulada: 'Articulada',
        MANAGED: 'Gestionada'
    };
    var sectorLabels = {
        salud: 'Salud',
        justicia: 'Justicia',
        proteccion: 'Protección',
        otras_instituciones: 'Otras instituciones',
        barrera_transversal: 'Barrera transversal'
    };

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data: function () {
            return {
                windowTitle: cfg.windowTitle || 'Barreras Departamento',
                currentUser: cfg.currentUser || '',
                currentRole: cfg.currentRole || '',
                currentUserId: cfg.currentUserId || '',
                assignedDepartmentId: cfg.assignedDepartmentId || '',

                isLoading: false,
                loadError: null,

                items: [],
                totalItems: 0,
                cities: [],

                filters: {
                    docNumber: '',
                    entidad: '',
                    cityId: '',
                    status: ''
                },
                currentPage: 0,
                pageSize: PAGE_SIZE,
                _filterTimer: null
            };
        },

        computed: {
            totalPages: function () {
                if (this.totalItems === 0) {
                    return 0;
                }
                return Math.ceil(this.totalItems / this.pageSize);
            },

            visiblePages: function () {
                var max = this.totalPages;
                var cur = this.currentPage;
                var delta = 2;
                var range = [];
                var result = [];
                var last;
                var i;
                for (i = 0; i < max; i++) {
                    if (i === 0 || i === max - 1 || (i >= cur - delta && i <= cur + delta)) {
                        range.push(i);
                    }
                }
                for (i = 0; i < range.length; i++) {
                    if (last !== undefined) {
                        if (range[i] - last === 2) {
                            result.push(last + 1);
                        } else if (range[i] - last > 2) {
                            result.push('...');
                        }
                    }
                    result.push(range[i]);
                    last = range[i];
                }
                return result;
            }
        },

        mounted: function () {
            document.title = typeof this.windowTitle === 'string'
                ? this.windowTitle
                : String(this.windowTitle || 'Barreras Departamento');
            var appEl = document.getElementById('app');
            if (appEl) {
                appEl.style.display = 'block';
            }
            this.reloadScreen();
        },

        methods: {
            slugify: function (str) {
                if (!str) {
                    return '';
                }
                return String(str)
                    .toLowerCase()
                    .normalize('NFD')
                    .replace(/[\u0300-\u036f]/g, '')
                    .replace(/\s+/g, '-');
            },

            formatFecha: function (value) {
                if (!value) {
                    return '—';
                }
                var d = new Date(value);
                if (isNaN(d.getTime())) {
                    return value;
                }
                var y = d.getFullYear();
                var m = String(d.getMonth() + 1).padStart(2, '0');
                var day = String(d.getDate()).padStart(2, '0');
                return y + '-' + m + '-' + day;
            },

            mapApiItem: function (b) {
                var victimFull = [b.victimName, b.victimLastName]
                    .filter(function (s) { return s && String(s).trim(); })
                    .join(' ') || '—';
                var sectorKey = b.sector || '';
                return {
                    id: b.id,
                    caseId: b.caseId || '',
                    victimName: victimFull,
                    docNumber: b.victimDocNumber || '—',
                    victimPhone: (b.victimPhone && String(b.victimPhone).trim()) || '—',
                    cityName: b.cityName || '—',
                    sectorKey: sectorKey,
                    sector: sectorLabels[sectorKey] || sectorKey || '—',
                    description: b.description || '—',
                    entityBranchName: b.entityBranchName || '—',
                    estadoLabel: estadoLabels[b.status] || b.status || '—',
                    priority: riskLevelLabels[b.riskLevel] || 'Desconocido',
                    fechaIdentificacion: this.formatFecha(b.barrierDate || b.createdAt)
                };
            },

            reloadScreen: function () {
                if (this.currentRole !== 'en') {
                    this.loadError = 'Rol no autorizado para acceder a esta pantalla.';
                    this.isLoading = false;
                    return;
                }
                this.filters.docNumber = '';
                this.filters.entidad = '';
                this.filters.cityId = '';
                this.filters.status = '';
                this.currentPage = 0;
                this.items = [];
                this.totalItems = 0;
                this.loadCities();
                this.loadBarreras();
            },

            loadCities: function () {
                var self = this;
                if (!this.assignedDepartmentId) {
                    this.cities = [];
                    return Promise.resolve();
                }
                return fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(this.assignedDepartmentId), {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (!res.ok) {
                            throw new Error('cities');
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        self.cities = Array.isArray(data) ? data : [];
                    })
                    .catch(function () {
                        self.cities = [];
                    });
            },

            loadBarreras: function () {
                var self = this;
                this.isLoading = true;
                this.loadError = null;

                var params = new URLSearchParams();
                params.set('page', String(this.currentPage));
                params.set('pageSize', String(this.pageSize));
                if ((this.filters.docNumber || '').trim()) {
                    params.set('docNumber', this.filters.docNumber.trim());
                }
                if ((this.filters.entidad || '').trim()) {
                    params.set('entidad', this.filters.entidad.trim());
                }
                if ((this.filters.cityId || '').trim()) {
                    params.set('cityId', this.filters.cityId.trim());
                }
                if ((this.filters.status || '').trim()) {
                    params.set('status', this.filters.status.trim());
                }

                return fetch('/api/v1/barriers-v2/department?' + params.toString(), {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (res.status === 403) {
                            return res.json().then(function (body) {
                                var msg = (body && body.error) || '';
                                if (msg.indexOf('departamento') !== -1) {
                                    throw new Error('Tu usuario no tiene un departamento asignado. Contacta a un administrador para que te lo asigne.');
                                }
                                throw new Error('No tienes permiso para ver esta lista.');
                            });
                        }
                        if (!res.ok) {
                            throw new Error('department-barriers');
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        if (data === null) {
                            return;
                        }
                        var raw = (data && data.items) ? data.items : [];
                        self.items = raw.map(function (b) { return self.mapApiItem(b); });
                        self.totalItems = (data && typeof data.total === 'number') ? data.total : 0;
                        self.isLoading = false;
                    })
                    .catch(function (err) {
                        self.loadError = (err && err.message && err.message !== 'department-barriers')
                            ? err.message
                            : 'No se pudo cargar la lista de barreras. Intenta de nuevo.';
                        self.isLoading = false;
                    });
            },

            onFilterText: function () {
                var self = this;
                if (this._filterTimer) {
                    clearTimeout(this._filterTimer);
                }
                this._filterTimer = setTimeout(function () {
                    self.currentPage = 0;
                    self.loadBarreras();
                }, FILTER_DEBOUNCE_MS);
            },

            onFilterSelect: function () {
                this.currentPage = 0;
                this.loadBarreras();
            },

            changePage: function (page) {
                if (page < 0 || page >= this.totalPages) {
                    return;
                }
                this.currentPage = page;
                window.scrollTo({ top: 0, behavior: 'smooth' });
                this.loadBarreras();
            },

            goToCase: function (item) {
                if (!item || !item.caseId) {
                    return;
                }
                window.location.href = '/salvia/casos/' + encodeURIComponent(item.caseId) + '/detalle';
            },

            goToBarrier: function (item) {
                if (!item || !item.id) {
                    return;
                }
                window.location.href = '/salvia/barreras/' + encodeURIComponent(item.id);
            }
        }
    });

    home.mount('#app');
})();
