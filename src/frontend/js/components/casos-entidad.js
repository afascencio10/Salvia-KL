/**
 * casos-entidad.js
 * Pantalla "Casos Entidad" (casos_entidad.html).
 *
 * E01 — Cuando carga la pantalla:
 *   GET /api/v1/entities → toma la primera entidad → carga ciudades + listado.
 *
 * Requiere window.__CasosEntidadConfig:
 *   - windowTitle, currentUser, currentRole, menu
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
                menu: cfg.menu || {},

                isLoading: false,
                loadError: null,

                entities: [],
                cities: [],
                sectorName: '',
                items: [],
                totalItems: 0,

                filters: {
                    entityId: '',
                    document: '',
                    city: ''
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

            applySelectedEntityMeta: function () {
                var id = String(this.filters.entityId || '');
                var found = null;
                for (var i = 0; i < this.entities.length; i++) {
                    if (String(this.entities[i].id) === id) {
                        found = this.entities[i];
                        break;
                    }
                }
                this.sectorName = found ? (found.sectorName || found.sector || '') : '';
            },

            // E01 — Cuando carga la pantalla
            reloadScreen: function () {
                if (this.currentRole !== 'et') {
                    this.loadError = 'Rol no autorizado para ver Casos Entidad.';
                    this.isLoading = false;
                    return;
                }

                this.isLoading = true;
                this.loadError = null;
                this.filters.document = '';
                this.filters.city = '';
                this.currentPage = 0;
                this.items = [];
                this.totalItems = 0;
                this.cities = [];

                var self = this;
                fetch('/api/v1/entities', { credentials: 'same-origin' })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (!res.ok) {
                            throw new Error('entities');
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        if (data === null) {
                            return;
                        }
                        self.entities = Array.isArray(data) ? data : [];

                        if (self.entities.length === 0) {
                            self.filters.entityId = '';
                            self.sectorName = '';
                            self.isLoading = false;
                            return;
                        }

                        // Temporal: primera entidad del catálogo (relación usuario↔entidad pendiente).
                        self.filters.entityId = String(self.entities[0].id);
                        self.applySelectedEntityMeta();
                        return self.loadCitiesAndCases();
                    })
                    .catch(function () {
                        self.loadError = 'No se pudo cargar el catálogo de entidades.';
                        self.isLoading = false;
                    });
            },

            loadCitiesAndCases: function () {
                var self = this;
                var entityId = this.filters.entityId;
                if (!entityId) {
                    this.items = [];
                    this.totalItems = 0;
                    this.cities = [];
                    this.isLoading = false;
                    return Promise.resolve();
                }

                this.isLoading = true;
                this.loadError = null;

                return fetch('/api/v1/entities/' + encodeURIComponent(entityId) + '/cities', {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (!res.ok) {
                            return [];
                        }
                        return res.json();
                    })
                    .then(function (cities) {
                        if (cities === null) {
                            return null;
                        }
                        self.cities = Array.isArray(cities) ? cities : [];
                        return self.loadEntityCases();
                    })
                    .catch(function () {
                        self.cities = [];
                        return self.loadEntityCases();
                    });
            },

            // SUB-FLUJO: Cargar listado paginado
            loadEntityCases: function () {
                var self = this;
                var entityId = this.filters.entityId;
                if (!entityId) {
                    this.items = [];
                    this.totalItems = 0;
                    this.isLoading = false;
                    return Promise.resolve();
                }

                this.isLoading = true;
                this.loadError = null;

                var params = new URLSearchParams();
                params.set('entityId', entityId);
                params.set('page', String(this.currentPage));
                params.set('pageSize', String(this.pageSize));
                if ((this.filters.document || '').trim()) {
                    params.set('document', this.filters.document.trim());
                }
                if ((this.filters.city || '').trim()) {
                    params.set('city', this.filters.city.trim());
                }

                return fetch('/api/v1/entity-cases?' + params.toString(), {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
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
                        if (data && data.sectorName) {
                            self.sectorName = data.sectorName;
                        } else {
                            self.applySelectedEntityMeta();
                        }
                        self.isLoading = false;
                    })
                    .catch(function () {
                        self.loadError = 'No se pudo cargar los casos de la entidad.';
                        self.isLoading = false;
                    });
            },

            // E07 — Cuando cambia entidad
            onChangeEntity: function () {
                // Conserva document/city; solo resetea página (flow-E07).
                this.currentPage = 0;
                this.applySelectedEntityMeta();

                if (!this.filters.entityId) {
                    this.sectorName = '';
                    this.items = [];
                    this.totalItems = 0;
                    this.cities = [];
                    return;
                }

                // Ciudad del select puede no existir en la nueva entidad → limpiar solo city.
                this.filters.city = '';
                this.loadCitiesAndCases();
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

            // E03 — Cuando filtra por ciudad (select de ciudades de la entidad)
            onFilterCity: function () {
                this.currentPage = 0;
                this.loadEntityCases();
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
