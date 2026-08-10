/**
 * casos-entidad.js
 * Pantalla "Casos Entidad" (casos_entidad.html).
 *
 * E01 — Cuando carga la pantalla:
 *   GET /api/v1/entity-cases → listado de la sede del usuario (sesión) + meta entidad padre.
 *
 * Requiere window.__CasosEntidadConfig:
 *   - windowTitle, currentUser, currentRole, entityBranchId, menu
 */

(function () {
    'use strict';

    var PAGE_SIZE = 5;
    var FILTER_DEBOUNCE_MS = 350;

    var cfg = window.__CasosEntidadConfig || {};

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data: function () {
            return {
                windowTitle: cfg.windowTitle || 'Casos Entidad',
                currentUser: cfg.currentUser || '',
                currentRole: cfg.currentRole || '',
                entityBranchId: Number(cfg.entityBranchId) || 0,
                menu: cfg.menu || {},

                isLoading: false,
                loadError: null,

                entityName: '',
                sectorName: '',
                items: [],
                totalItems: 0,

                filters: {
                    document: ''
                },
                currentPage: 0,
                pageSize: PAGE_SIZE,

                toast: { visible: false, message: '', type: 'info' },
                _toastTimer: null,
                _documentTimer: null
            };
        },

        computed: {
            totalPages: function () {
                if (this.totalItems === 0) {
                    return 0;
                }
                return Math.ceil(this.totalItems / this.pageSize);
            },

            hasSede: function () {
                return this.entityBranchId > 0;
            }
        },

        mounted: function () {
            document.title = typeof this.windowTitle === 'string'
                ? this.windowTitle
                : String(this.windowTitle || 'Casos Entidad');
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

            // E01 — Cuando carga la pantalla
            reloadScreen: function () {
                if (this.currentRole !== 'et') {
                    this.loadError = 'Rol no autorizado para ver Casos Entidad.';
                    this.isLoading = false;
                    return;
                }

                if (!this.hasSede) {
                    this.loadError = 'Tu usuario no tiene una sede (entity_branch_id) asignada.';
                    this.isLoading = false;
                    return;
                }

                this.filters.document = '';
                this.currentPage = 0;
                this.items = [];
                this.totalItems = 0;
                this.loadEntityCases();
            },

            // SUB-FLUJO: Cargar listado paginado
            loadEntityCases: function () {
                var self = this;
                if (!this.hasSede) {
                    this.items = [];
                    this.totalItems = 0;
                    this.isLoading = false;
                    return Promise.resolve();
                }

                this.isLoading = true;
                this.loadError = null;

                var params = new URLSearchParams();
                params.set('page', String(this.currentPage));
                params.set('pageSize', String(this.pageSize));
                if ((this.filters.document || '').trim()) {
                    params.set('document', this.filters.document.trim());
                }

                return fetch('/api/v1/entity-cases?' + params.toString(), {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (res.status === 400) {
                            return res.json().then(function (body) {
                                throw new Error((body && body.error) || 'bad_request');
                            });
                        }
                        if (!res.ok) {
                            throw new Error('entity-cases');
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        if (data === null) {
                            return;
                        }
                        self.items = (data && data.items) ? data.items : [];
                        self.totalItems = (data && typeof data.total === 'number') ? data.total : 0;
                        if (data && data.entityName) {
                            self.entityName = data.entityName;
                        }
                        if (data && data.sectorName) {
                            self.sectorName = data.sectorName;
                        } else if (data && data.sector) {
                            self.sectorName = data.sector;
                        }
                        self.isLoading = false;
                    })
                    .catch(function (err) {
                        self.loadError = (err && err.message && err.message !== 'entity-cases')
                            ? err.message
                            : 'No se pudo cargar los casos de la entidad.';
                        self.isLoading = false;
                    });
            },

            // E02 — Cuando filtra por documento (parcial + debounce)
            onFilterDocument: function () {
                var self = this;
                if (this._documentTimer) {
                    clearTimeout(this._documentTimer);
                }
                this._documentTimer = setTimeout(function () {
                    self.currentPage = 0;
                    self.loadEntityCases();
                }, FILTER_DEBOUNCE_MS);
            },

            // E04 — Cuando cambia de página
            changePage: function (page) {
                if (page < 0 || page >= this.totalPages) {
                    return;
                }
                this.currentPage = page;
                window.scrollTo({ top: 0, behavior: 'smooth' });
                this.loadEntityCases();
            },

            // E05 — Cuando presiona "Ver caso"
            goToCase: function (item) {
                if (!item || !item.caseId) {
                    this.showToast('No se encontró el identificador del caso.', 'info');
                    return;
                }
                window.location.href = '/salvia/casos/' + encodeURIComponent(item.caseId) + '/detalle';
            },

            // E06 — Cuando presiona "Ver detalle"
            // Destino /salvia/entidad/:id aún no existe (GAP compartido con case-entities).
            goToRelationDetail: function (item) {
                if (!item || !item.entityBranchICode) {
                    this.showToast('No se encontró la sede de la entidad.', 'info');
                    return;
                }
                this.showToast(
                    'Detalle de entidad pendiente (sede: ' + item.entityBranchICode + ')',
                    'info'
                );
            }
        }
    });

    home.mount('#app');
})();
