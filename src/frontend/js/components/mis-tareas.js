/**
 * mis-tareas.js
 * Lógica Vue para la pantalla "Mis Tareas" (mis_tareas.html).
 *
 * La paginación y el filtro por estado se resuelven en el servidor
 * (GET /api/v1/case-tasks/mine); el agente se toma de la sesión, no de la URL.
 */

(function () {
    'use strict';

    var cfg = window.__MisTareasConfig || {};

    var TIPO_LABELS = {
        orientacion_llamada: 'Orientación llamada',
        gestion_llamada:     'Gestión llamada',
        gestion_propia:      'Gestión propia',
        proyectar_oficio:    'Proyectar oficio',
        comite_caso:         'Comité de caso'
    };

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data: function () {
            return {
                currentUser: cfg.currentUser || '',
                currentRole: cfg.currentRole || '',

                isLoading: false,
                loadError: null,

                estadoOpciones: [
                    { value: '',     label: 'Todas' },
                    { value: 'ToDo', label: 'Pendientes' },
                    { value: 'Done', label: 'Completadas' }
                ],
                filtroEstado: 'ToDo',

                tareas:   [],
                total:    0,
                page:     1,
                pageSize: 20
            };
        },

        mounted: function () {
            document.getElementById('app').style.display = 'block';
            this.cargarTareas();
        },

        computed: {
            numPages: function () {
                return Math.max(1, Math.ceil(this.total / this.pageSize));
            },
            filtroEstadoTexto: function () {
                if (this.filtroEstado === 'ToDo') return 'pendientes';
                if (this.filtroEstado === 'Done') return 'completadas';
                return 'asignadas';
            }
        },

        methods: {

            cargarTareas: function () {
                var self = this;
                this.isLoading = true;
                this.loadError = null;

                var url = '/api/v1/case-tasks/mine'
                    + '?page=' + this.page
                    + '&pageSize=' + this.pageSize
                    + '&status=' + encodeURIComponent(this.filtroEstado);

                getData(url, function (status, response) {
                    self.isLoading = false;

                    if (status === 200 && response) {
                        self.tareas   = (response.items || []).map(self.mapApiToTarea);
                        self.total    = response.total || 0;
                        self.page     = response.page || 1;
                        self.pageSize = response.pageSize || self.pageSize;

                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else if (status === 403) {
                        self.loadError = 'Tu rol no tiene acceso a esta pantalla.';
                    } else {
                        self.loadError = 'No se pudo cargar la lista de tareas. Intenta de nuevo.';
                    }
                }, true);
            },

            /* Transforma un CaseTaskWithRelations al formato de la tabla */
            mapApiToTarea: function (t) {
                var victima = [t.victimName, t.victimLastName]
                    .filter(function (s) { return s && s.trim(); })
                    .join(' ') || '—';

                var completada = t.status === 'Done';

                return {
                    id:          t.id,
                    caseId:      t.caseId || '',
                    caseCode:    t.caseCode || t.caseId || '—',
                    victima:     victima,
                    tipoLabel:   this.tipoLabel(t.type, t.category),
                    descripcion: t.description || '—',
                    completada:  completada,
                    estadoLabel: completada ? 'Completada' : 'Pendiente',
                    fecha:       this.formatFecha(completada ? (t.completedAt || t.createdAt) : t.createdAt)
                };
            },

            tipoLabel: function (type, category) {
                if (!type) return category || '—';
                if (TIPO_LABELS[type]) return TIPO_LABELS[type];
                return type.replace(/_/g, ' ').replace(/^./, function (c) { return c.toUpperCase(); });
            },

            formatFecha: function (value) {
                if (!value) return '—';
                var d = new Date(value);
                if (isNaN(d.getTime())) return value;
                return d.toLocaleDateString('es-CO', { timeZone: 'America/Bogota', year: 'numeric', month: '2-digit', day: '2-digit' });
            },

            cambiarEstado: function (valor) {
                if (this.filtroEstado === valor) return;
                this.filtroEstado = valor;
                this.page = 1;
                this.cargarTareas();
            },

            irAPagina: function (n) {
                if (n < 1 || n > this.numPages || n === this.page) return;
                this.page = n;
                this.cargarTareas();
                window.scrollTo(0, 0);
            },

            logout: function () {
                logout('/seguridad/logout');
            }
        }
    });

    home.mount('#app');
})();
