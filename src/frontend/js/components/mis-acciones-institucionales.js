/**
 * mis-acciones-institucionales.js
 * Pantalla "Mis Acciones Institucionales" (mis_acciones_institucionales.html).
 *
 * Al cargar y en cada cambio de pestaña o filtro:
 *   GET /api/v1/mis-acciones-institucionales?status=&document=&cityId=&page=&pageSize=
 *
 * Esa respuesta trae todo lo que la pantalla necesita: las filas de la pestaña
 * activa, los contadores de las dos pestañas, las opciones del filtro de ciudad y
 * el catálogo de estados para el modal. La sede sale de la sesión en el backend, no
 * se envía desde acá.
 *
 * Requiere window.__MisAccionesInstitucionalesConfig:
 *   - windowTitle, currentUser, currentRole, entityBranchId, menu
 */

(function () {
    'use strict';

    var PAGE_SIZE = 10;
    var FILTER_DEBOUNCE_MS = 350;

    var cfg = window.__MisAccionesInstitucionalesConfig || {};

    var app = Vue.createApp({
        delimiters: ['${', '}'],

        data: function () {
            return {
                windowTitle: cfg.windowTitle || 'Mis Acciones Institucionales',
                currentUser: cfg.currentUser || '',
                currentRole: cfg.currentRole || '',
                entityBranchId: Number(cfg.entityBranchId) || 0,
                menu: cfg.menu || {},

                isLoading: false,
                loadError: null,

                entityName: '',
                status: 'pendiente',
                items: [],
                totalItems: 0,
                counts: { pendientes: 0, completadas: 0 },
                cityOptions: [],
                estadoOptions: [],

                checklists: {},

                // El modal se prellena con la respuesta que ya exista, así que
                // "Editar respuesta" no obliga a volver a escribir todo.
                modal: { visible: false, item: null, radicado: '', ampliacion: '', estado: '' },
                isSaving: false,
                saveError: null,

                filters: { document: '', cityId: '' },
                currentPage: 0,
                pageSize: PAGE_SIZE,

                toast: { visible: false, message: '', type: 'info' },
                _toastTimer: null,
                _documentTimer: null
            };
        },

        computed: {
            totalPages: function () {
                if (this.totalItems === 0) return 0;
                return Math.ceil(this.totalItems / this.pageSize);
            },

            hasSede: function () {
                return this.entityBranchId > 0;
            }
        },

        mounted: function () {
            document.title = typeof this.windowTitle === 'string'
                ? this.windowTitle
                : String(this.windowTitle || 'Mis Acciones Institucionales');
            var appEl = document.getElementById('app');
            if (appEl) appEl.style.display = 'block';
            this.reloadScreen();
        },

        methods: {
            riskClass: function (level) {
                return String(level || '')
                    .toLowerCase()
                    .normalize('NFD')
                    .replace(/[̀-ͯ]/g, '')
                    .replace(/\s+/g, '');
            },

            // El catálogo de estados llega del backend, así que la etiqueta se resuelve
            // contra él y no contra un mapa duplicado en el frontend.
            estadoLabel: function (value) {
                for (var i = 0; i < this.estadoOptions.length; i++) {
                    if (this.estadoOptions[i].value === value) return this.estadoOptions[i].label;
                }
                return value || '';
            },

            // El checklist viene indexado por caso y no repetido en cada fila, porque
            // varias acciones del mismo caso lo comparten. Se devuelve un objeto vacío
            // por defecto para que la plantilla no tenga que comprobar nulos.
            checklistDe: function (item) {
                var checklist = item && this.checklists[item.entityCaseId];
                if (!checklist) return { total: 0, completadas: 0, items: [] };
                return {
                    total: checklist.total || 0,
                    completadas: checklist.completadas || 0,
                    items: Array.isArray(checklist.items) ? checklist.items : []
                };
            },

            // Porcentaje para la barra de avance de la tarjeta. Se calcula acá y no
            // en la plantilla para no dejar aritmética suelta en el HTML.
            progresoDe: function (item) {
                var checklist = this.checklistDe(item);
                if (!checklist.total) return 0;
                return Math.round((checklist.completadas / checklist.total) * 100);
            },

            verDetalle: function (item) {
                if (!item || !item.entityCaseId) {
                    this.showToast('No se encontró el identificador de la relación.', 'info');
                    return;
                }
                // Misma pantalla de detalle que ya usa Casos Entidad.
                window.location.href = '/salvia/casos-entidad/' + encodeURIComponent(item.entityCaseId) + '/detalle';
            },

            showToast: function (message, type) {
                var self = this;
                if (this._toastTimer) clearTimeout(this._toastTimer);
                this.toast = { visible: true, message: message, type: type || 'info' };
                this._toastTimer = setTimeout(function () { self.toast.visible = false; }, 2800);
            },

            reloadScreen: function () {
                if (this.currentRole !== 'et') {
                    this.loadError = 'Rol no autorizado para ver Mis Acciones Institucionales.';
                    this.isLoading = false;
                    return;
                }
                // Sin sede no hay error que reportar: la plantilla muestra un mensaje
                // explicando qué hacer, que es más útil que un bloque de error.
                if (!this.hasSede) {
                    this.isLoading = false;
                    return;
                }
                this.status = 'pendiente';
                this.filters = { document: '', cityId: '' };
                this.currentPage = 0;
                this.items = [];
                this.totalItems = 0;
                this.checklists = {};
                this.cerrarModal();
                this.loadAcciones();
            },

            loadAcciones: function () {
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
                params.set('status', this.status);
                params.set('page', String(this.currentPage));
                params.set('pageSize', String(this.pageSize));
                if ((this.filters.document || '').trim()) {
                    params.set('document', this.filters.document.trim());
                }
                if (this.filters.cityId) {
                    params.set('cityId', this.filters.cityId);
                }

                return fetch('/api/v1/mis-acciones-institucionales?' + params.toString(), {
                    credentials: 'same-origin'
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (!res.ok) {
                            return res.json().catch(function () { return {}; }).then(function (body) {
                                throw new Error((body && body.error) || 'Error ' + res.status);
                            });
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        if (!data) return;
                        self.entityName = data.entityName || '';
                        self.items = Array.isArray(data.items) ? data.items : [];
                        self.totalItems = Number(data.total) || 0;
                        self.counts = data.counts || { pendientes: 0, completadas: 0 };
                        self.cityOptions = Array.isArray(data.cityOptions) ? data.cityOptions : [];
                        self.estadoOptions = Array.isArray(data.estadoOptions) ? data.estadoOptions : [];
                        self.checklists = data.checklists || {};
                        self.isLoading = false;

                        // Al completar la última acción de una página, esa página deja de
                        // existir y el usuario vería un listado vacío habiendo datos. La
                        // condición solo puede bajar de página, así que no cicla.
                        var ultimaPagina = Math.max(0, Math.ceil(self.totalItems / self.pageSize) - 1);
                        if (self.items.length === 0 && self.totalItems > 0 && self.currentPage > ultimaPagina) {
                            self.currentPage = ultimaPagina;
                            return self.loadAcciones();
                        }
                    })
                    .catch(function (err) {
                        self.loadError = err && err.message
                            ? err.message
                            : 'No se pudieron cargar las acciones institucionales.';
                        self.isLoading = false;
                    });
            },

            changeStatus: function (status) {
                if (this.status === status || this.isLoading) return;
                this.status = status;
                this.currentPage = 0;
                this.loadAcciones();
            },

            // Los filtros vuelven a la primera página: si el usuario estaba en la
            // página 3 y filtra, el resultado puede tener una sola página y quedaría
            // viendo un listado vacío.
            onFilterDocument: function () {
                var self = this;
                if (this._documentTimer) clearTimeout(this._documentTimer);
                this._documentTimer = setTimeout(function () {
                    self.currentPage = 0;
                    self.loadAcciones();
                }, FILTER_DEBOUNCE_MS);
            },

            onFilterCity: function () {
                this.currentPage = 0;
                this.loadAcciones();
            },

            changePage: function (page) {
                if (page < 0 || page >= this.totalPages || page === this.currentPage) return;
                this.currentPage = page;
                this.loadAcciones();
            },

            abrirModalCompletar: function (item) {
                if (!item) return;
                this.saveError = null;
                // Prellenado con lo ya guardado: el modal sirve igual para completar y
                // para corregir, porque el backend hace upsert sobre (caso, acción).
                this.modal = {
                    visible: true,
                    item: item,
                    radicado: item.radicado || '',
                    ampliacion: item.ampliacion || '',
                    estado: item.estado || ''
                };
            },

            cerrarModal: function () {
                if (this.isSaving) return;
                this.modal = { visible: false, item: null, radicado: '', ampliacion: '', estado: '' };
                this.saveError = null;
            },

            guardarRespuesta: function () {
                var self = this;
                if (this.isSaving || !this.modal.item) return;

                // El estado es lo único obligatorio, igual que en el backend: sin él la
                // respuesta no dice si la acción quedó atendida.
                if (!this.modal.estado) {
                    this.saveError = 'Selecciona el estado general de la respuesta.';
                    return;
                }

                this.isSaving = true;
                this.saveError = null;

                return fetch('/api/v1/mis-acciones-institucionales/respuesta', {
                    method: 'POST',
                    credentials: 'same-origin',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        entityCaseId: this.modal.item.entityCaseId,
                        entityObligationId: this.modal.item.entityObligationId,
                        radicado: this.modal.radicado,
                        ampliacion: this.modal.ampliacion,
                        estado: this.modal.estado
                    })
                })
                    .then(function (res) {
                        if (res.status === 401) {
                            window.location.href = '/static/landing.html';
                            return null;
                        }
                        if (!res.ok) {
                            return res.json().catch(function () { return {}; }).then(function (body) {
                                throw new Error((body && body.error) || 'Error ' + res.status);
                            });
                        }
                        return res.json();
                    })
                    .then(function (data) {
                        if (!data) return;
                        self.isSaving = false;
                        self.cerrarModal();
                        self.showToast('Respuesta registrada.', 'success');
                        // Se recarga en vez de parchear la fila en memoria: al completar,
                        // la acción cambia de pestaña y los dos contadores se mueven.
                        return self.loadAcciones();
                    })
                    .catch(function (err) {
                        self.isSaving = false;
                        self.saveError = err && err.message
                            ? err.message
                            : 'No se pudo registrar la respuesta.';
                    });
            }
        }
    });

    app.mount('#app');
})();
