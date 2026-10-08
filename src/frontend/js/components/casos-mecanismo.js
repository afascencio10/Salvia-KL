/**
 * casos-mecanismo.js
 * Pantalla "Casos Mecanismo" (casos_mecanismo.html).
 *
 * Expone `home` / `app` para registrar case-task-modal antes del mount.
 */

(function () {
    'use strict';

    var PAGE_SIZE = 5;
    var FILTER_DEBOUNCE_MS = 350;
    var cfg = window.__CasosMecanismoConfig || {};

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data: function () {
            return {
                windowTitle: cfg.windowTitle || 'Casos Mecanismo',
                currentUser: cfg.currentUser || '',
                currentRole: cfg.currentRole || '',
                currentUserId: cfg.currentUserId || '',
                menu: cfg.menu || {},

                isLoading: false,
                loadError: null,
                items: [],
                totalItems: 0,

                filters: {
                    document: '',
                    city: ''
                },
                currentPage: 0,
                pageSize: PAGE_SIZE,

                toast: { visible: false, message: '', type: 'info' },
                _toastTimer: null,
                _documentTimer: null,
                _cityTimer: null
            };
        },

        computed: {
            totalPages: function () {
                if (this.totalItems === 0) {
                    return 0;
                }
                return Math.ceil(this.totalItems / this.pageSize);
            }
        },

        mounted: function () {
            document.title = typeof this.windowTitle === 'string'
                ? this.windowTitle
                : String(this.windowTitle || 'Casos Mecanismo');
            var appEl = document.getElementById('app');
            if (appEl) {
                appEl.style.display = 'block';
            }
            this.reloadScreen();
        },

        methods: {
            riskClass: function (level) {
                return String(level || '')
                    .toLowerCase()
                    .normalize('NFD')
                    .replace(/[\u0300-\u036f]/g, '')
                    .replace(/\s+/g, '');
            },

            statusClass: function (status) {
                return String(status || '')
                    .toLowerCase()
                    .normalize('NFD')
                    .replace(/[\u0300-\u036f]/g, '')
                    .replace(/\s+/g, '');
            },

            scopeLabel: function (scope) {
                var map = {
                    municipal: 'Municipal',
                    departamental: 'Departamental',
                    nacional: 'Nacional'
                };
                return map[scope] || scope || '';
            },

            formatDate: function (value) {
                if (!value) {
                    return '';
                }
                var d = new Date(value);
                if (isNaN(d.getTime())) {
                    return String(value).slice(0, 10);
                }
                // Zona Colombia (SALVIA solo se usa en Colombia): evita desfase de día
                // por la zona del navegador.
                var parts = new Intl.DateTimeFormat('en-CA', {
                    timeZone: 'America/Bogota', year: 'numeric', month: '2-digit', day: '2-digit'
                }).formatToParts(d).reduce(function(acc, p) { acc[p.type] = p.value; return acc; }, {});
                return parts.year + '-' + parts.month + '-' + parts.day;
            },

            showToast: function (message, type) {
                var self = this;
                if (this._toastTimer) {
                    clearTimeout(this._toastTimer);
                }
                this.toast = {
                    visible: true,
                    message: message,
                    type: type || 'info'
                };
                this._toastTimer = setTimeout(function () {
                    self.toast.visible = false;
                }, 2800);
            },

            reloadScreen: function () {
                if (this.currentRole !== 'mc') {
                    this.loadError = 'Rol no autorizado para ver Casos Mecanismo.';
                    this.isLoading = false;
                    return;
                }
                this.currentPage = 0;
                this.loadCases();
            },

            loadCases: function () {
                var self = this;
                this.isLoading = true;
                this.loadError = null;

                var params = new URLSearchParams();
                params.set('page', String(this.currentPage));
                params.set('pageSize', String(this.pageSize));
                if ((this.filters.document || '').trim()) {
                    params.set('document', this.filters.document.trim());
                }
                if ((this.filters.city || '').trim()) {
                    params.set('city', this.filters.city.trim());
                }

                return fetch('/api/v1/mecanismos/mine?' + params.toString(), {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (res.status === 403) {
                            throw new Error('No tienes permiso para ver estos casos.');
                        }
                        if (!res.ok) {
                            throw new Error('mecanismos-mine');
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        if (data === null) {
                            return;
                        }
                        self.items = (data && data.items) ? data.items : [];
                        self.totalItems = (data && typeof data.total === 'number') ? data.total : 0;
                        self.isLoading = false;
                    })
                    .catch(function (err) {
                        self.loadError = (err && err.message && err.message !== 'mecanismos-mine')
                            ? err.message
                            : 'No se pudo cargar los casos del mecanismo.';
                        self.isLoading = false;
                    });
            },

            onFilterDocument: function () {
                var self = this;
                if (this._documentTimer) {
                    clearTimeout(this._documentTimer);
                }
                this._documentTimer = setTimeout(function () {
                    self.currentPage = 0;
                    self.loadCases();
                }, FILTER_DEBOUNCE_MS);
            },

            onFilterCity: function () {
                var self = this;
                if (this._cityTimer) {
                    clearTimeout(this._cityTimer);
                }
                this._cityTimer = setTimeout(function () {
                    self.currentPage = 0;
                    self.loadCases();
                }, FILTER_DEBOUNCE_MS);
            },

            changePage: function (page) {
                if (page < 0 || page >= this.totalPages) {
                    return;
                }
                this.currentPage = page;
                window.scrollTo({ top: 0, behavior: 'smooth' });
                this.loadCases();
            },

            goToCase: function (item) {
                var code = item && (item.caseCode || item.caseId);
                if (!code) {
                    this.showToast('No se encontró el identificador del caso.', 'info');
                    return;
                }
                window.location.href = '/salvia/casos/' + encodeURIComponent(code) + '/detalle';
            },

            openGestionar: function (item) {
                if (!item || !this.$refs.taskModal || !this.$refs.taskModal.openMecanismo) {
                    return;
                }
                this.$refs.taskModal.openMecanismo({
                    mechanismId: item.id,
                    caseId: item.caseId,
                    scope: item.scope
                });
            },

            onTaskCompleted: function () {
                this.showToast('Gestión guardada.', 'success');
                this.loadCases();
            },

            logout: function () {
                logout('/seguridad/logout');
            }
        }
    });

    window.home = home;
    window.app = home;
})();
