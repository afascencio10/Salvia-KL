/**
 * casos-entidad.js
 * Pantalla "Casos Entidad" (casos_entidad.html) — UI con datos de prueba.
 *
 * Requiere window.__CasosEntidadConfig:
 *   - windowTitle {string|any}
 *   - currentUser {string}
 *   - currentRole {string}
 *   - menu        {object}
 *
 * Nota: los datos vienen de MOCK_ENTITY_CASES (quemados). Se reemplazarán
 * por API cuando exista salvia.entity_case.
 */

(function () {
    'use strict';

    var PAGE_SIZE = 5;
    var FILTER_DEBOUNCE_MS = 350;

    /** Datos de prueba alineados al mockup de planeación. */
    var MOCK_ENTITY_CASES = [
        {
            entityCaseId: 'ec-001',
            caseId: 'case-001',
            caseCode: 'SAL-001',
            victimFullName: 'María González',
            document: '1023456780',
            city: 'San Salvador',
            riskLevel: 'Alto',
            caseStatus: 'Activo',
            relationDescription: 'Exigir la recepción formal de la denuncia por violencia intrafamiliar.',
            lastAction: {
                description: 'Oficio radicado exigiendo la recepción formal de la denuncia.',
                date: '2026-03-10'
            }
        },
        {
            entityCaseId: 'ec-002',
            caseId: 'case-002',
            caseCode: 'SAL-003',
            victimFullName: 'Carmen Flores',
            document: '1045678902',
            city: 'Sonsonate',
            riskLevel: 'Medio',
            caseStatus: 'En seguimiento',
            relationDescription: 'Solicitar copia de la denuncia interpuesta por violencia familiar.',
            lastAction: null
        },
        {
            entityCaseId: 'ec-003',
            caseId: 'case-003',
            caseCode: 'SAL-007',
            victimFullName: 'Ana Martínez',
            document: '1089012345',
            city: 'San Salvador',
            riskLevel: 'Extremo',
            caseStatus: 'Activo',
            relationDescription: 'Solicitar medidas de protección inmediata ante amenaza de feminicidio.',
            lastAction: {
                description: 'Oficio enviado solicitando medidas de protección.',
                date: '2026-03-18'
            }
        },
        {
            entityCaseId: 'ec-004',
            caseId: 'case-004',
            caseCode: 'SAL-012',
            victimFullName: 'Rosa Hernández',
            document: '1098765432',
            city: 'Santa Ana',
            riskLevel: 'Bajo',
            caseStatus: 'En seguimiento',
            relationDescription: 'Verificar el estado de la denuncia radicada en la sede local.',
            lastAction: {
                description: 'Respuesta de la entidad recibida y archivada.',
                date: '2026-02-28'
            }
        },
        {
            entityCaseId: 'ec-005',
            caseId: 'case-005',
            caseCode: 'SAL-015',
            victimFullName: 'Lucía Ramírez',
            document: '1011121314',
            city: 'La Libertad',
            riskLevel: 'Alto',
            caseStatus: 'Activo',
            relationDescription: 'Exigir informe de avance sobre la investigación policial.',
            lastAction: null
        },
        {
            entityCaseId: 'ec-006',
            caseId: 'case-006',
            caseCode: 'SAL-018',
            victimFullName: 'Patricia López',
            document: '1055667788',
            city: 'San Miguel',
            riskLevel: 'Medio',
            caseStatus: 'Activo',
            relationDescription: 'Solicitar traslado del expediente a la fiscalía competente.',
            lastAction: {
                description: 'Oficio de traslado proyectado para revisión.',
                date: '2026-03-05'
            }
        },
        {
            entityCaseId: 'ec-007',
            caseId: 'case-007',
            caseCode: 'SAL-021',
            victimFullName: 'Elena Vargas',
            document: '1022334455',
            city: 'Sonsonate',
            riskLevel: 'Alto',
            caseStatus: 'Cerrado',
            relationDescription: 'Confirmar cierre formal del trámite de denuncia.',
            lastAction: {
                description: 'Constancia de cierre recibida por la entidad.',
                date: '2026-01-20'
            }
        }
    ];

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

                sectorName: 'Justicia',
                allItems: [],
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
            filteredItems: function () {
                var doc = (this.filters.document || '').trim().toLowerCase();
                var city = (this.filters.city || '').trim().toLowerCase();

                return this.allItems.filter(function (item) {
                    var matchDoc = !doc || String(item.document || '').toLowerCase().indexOf(doc) !== -1;
                    var matchCity = !city || String(item.city || '').toLowerCase().indexOf(city) !== -1;
                    return matchDoc && matchCity;
                });
            },

            totalPages: function () {
                if (this.filteredItems.length === 0) {
                    return 0;
                }
                return Math.ceil(this.filteredItems.length / this.pageSize);
            },

            pagedItems: function () {
                var maxPage = Math.max(0, this.totalPages - 1);
                var page = Math.min(this.currentPage, maxPage);
                var start = page * this.pageSize;
                return this.filteredItems.slice(start, start + this.pageSize);
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

            // E01 — Cuando carga la pantalla (mock)
            reloadScreen: function () {
                if (this.currentRole !== 'et') {
                    this.loadError = 'Rol no autorizado para ver Casos Entidad.';
                    this.isLoading = false;
                    return;
                }

                this.isLoading = true;
                this.loadError = null;

                var self = this;
                // Simula latencia de red; luego se reemplaza por fetch al API.
                setTimeout(function () {
                    self.allItems = MOCK_ENTITY_CASES.slice();
                    self.currentPage = 0;
                    self.isLoading = false;
                }, 350);
            },

            // E02 — Cuando filtra por documento (mock, client-side)
            onFilterDocument: function () {
                var self = this;
                if (this._documentTimer) {
                    clearTimeout(this._documentTimer);
                }
                this._documentTimer = setTimeout(function () {
                    self.currentPage = 0;
                }, FILTER_DEBOUNCE_MS);
            },

            // E03 — Cuando filtra por ciudad (mock, client-side)
            onFilterCity: function () {
                var self = this;
                if (this._cityTimer) {
                    clearTimeout(this._cityTimer);
                }
                this._cityTimer = setTimeout(function () {
                    self.currentPage = 0;
                }, FILTER_DEBOUNCE_MS);
            },

            // E04 — Cuando cambia de página
            changePage: function (page) {
                if (page < 0 || page >= this.totalPages) {
                    return;
                }
                this.currentPage = page;
                window.scrollTo({ top: 0, behavior: 'smooth' });
            },

            // E05 — Cuando presiona "Ver caso" (destino pendiente)
            goToCase: function (item) {
                this.showToast('Ver caso pendiente de conectar (caseId: ' + item.caseId + ')', 'info');
            },

            // E06 — Cuando presiona "Ver detalle" (destino pendiente)
            goToRelationDetail: function (item) {
                this.showToast('Ver detalle pendiente de conectar (id: ' + item.entityCaseId + ')', 'info');
            }
        }
    });

    home.mount('#app');
})();
