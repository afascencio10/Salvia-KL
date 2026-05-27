/**
 * oficios-list.js — Componente reutilizable de lectura de Oficios (EntityLetter)
 *
 * Props opcionales (via data-attributes en el div de montaje o vars Go template):
 *   data-case-id      → filtra por caseId    → GET /api/v1/entity-letters?caseId=...
 *   data-barrier-id   → filtra por barrierId → GET /api/v1/entity-letters?barrierId=...
 *   data-user-id      → filtra por userId    → GET ?agentId=... + GET ?notificationUserId=... (merge)
 *
 * Si no se pasa ningún prop, carga la lista paginada general.
 *
 * Respuesta con relaciones (para userId): EntityLetterWithRelations
 *   Campos extra: barrierSector, barrierDescription, victimName, victimLastName,
 *                 victimDocNumber, caseCode
 *
 * Respuesta sin relaciones (para caseId/barrierId/ninguno): EntityLetter
 *   Los campos case y barrier del modal solo muestran el ID con enlace.
 */

(function () {
    'use strict';

    /* ────────────────────────────────────────────────────────────────
       Etiquetas y constantes
    ──────────────────────────────────────────────────────────────── */

    var STATUS_LABELS = {
        radicado:            'Radicado',
        por_proyectar:       'Por proyectar',
        para_revisar:        'Para revisar',
        en_correccion:       'En corrección',
        aprobacion_juridica: 'Aprobación jurídica',
        para_radicar:        'Para radicar',
        respondido:          'Respondido',
    };

    var NIVEL_LABELS = {
        nacional:      'Nacional',
        departamental: 'Departamental',
        municipal:     'Municipal',
    };

    var PRIORITY_LABELS = {
        normal: 'Normal',
        alta:   'Alta',
    };

    /* ────────────────────────────────────────────────────────────────
       Definición del componente Vue
    ──────────────────────────────────────────────────────────────── */

    var OficiosListComponent = {
        delimiters: ['${', '}'],
        template:   '#tpl-oficios-list',

        props: {
            caseId:    { default: null },
            barrierId: { default: null },
            userId:    { default: null },
        },

        data: function () {
            return {
                oficios:        [],
                selectedOficio: null,
                isLoading:      false,
                loadError:      null,
            };
        },

        computed: {
            /* La API ya filtra; el computed solo ordena por fecha desc */
            sortedOficios: function () {
                return this.oficios.slice().sort(function (a, b) {
                    return new Date(b.created_at) - new Date(a.created_at);
                });
            },
        },

        mounted: function () {
            this.loadOficios();
        },

        methods: {

            /* ── Carga de datos ─────────────────────────────────────────── */

            loadOficios: function () {
                var self = this;
                self.isLoading = true;
                self.loadError = null;
                self.oficios   = [];

                if (self.caseId) {
                    /* Filtra por caso → devuelve []EntityLetter (sin relaciones) */
                    var url = '/api/v1/entity-letters?caseId=' + encodeURIComponent(self.caseId);
                    self._fetchSingle(url, false);

                } else if (self.barrierId) {
                    /* Filtra por barrera → devuelve []EntityLetter (sin relaciones) */
                    var url = '/api/v1/entity-letters?barrierId=' + encodeURIComponent(self.barrierId);
                    self._fetchSingle(url, false);

                } else if (self.userId) {
                    /* Filtra por usuario como agente O como notificador.
                       Dos llamadas paralelas que se fusionan y deduplicán por id.
                       Devuelven []EntityLetterWithRelations (con datos de caso y barrera). */
                    var url1 = '/api/v1/entity-letters?agentId=' + encodeURIComponent(self.userId);
                    var url2 = '/api/v1/entity-letters?notificationUserId=' + encodeURIComponent(self.userId);
                    self._fetchParallel(url1, url2, true);

                } else {
                    /* Sin filtros → lista paginada general (sin relaciones) */
                    self._fetchSingle('/api/v1/entity-letters?limit=100&page=0', false);
                }
            },

            /* GET único → normaliza y asigna self.oficios */
            _fetchSingle: function (url, withRelations) {
                var self = this;
                getData(url, function (status, response) {
                    self.isLoading = false;
                    if (status === 200) {
                        var items = self._toArray(response);
                        self.oficios = items.map(function (el) {
                            return self.mapApiToOficio(el, withRelations);
                        });
                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else {
                        self.loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.';
                    }
                }, true);
            },

            /* Dos GETs paralelos → fusiona y deduplica por id */
            _fetchParallel: function (url1, url2, withRelations) {
                var self        = this;
                var combined    = [];
                var completed   = 0;
                var hasError    = false;

                function onDone(status, response) {
                    completed++;
                    if (status === 200) {
                        combined = combined.concat(self._toArray(response));
                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                        return;
                    } else {
                        hasError = true;
                    }

                    if (completed < 2) return;

                    self.isLoading = false;

                    if (hasError && combined.length === 0) {
                        self.loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.';
                        return;
                    }

                    /* Deduplicar por id */
                    var seen = {};
                    var deduped = combined.filter(function (el) {
                        if (seen[el.id]) return false;
                        seen[el.id] = true;
                        return true;
                    });

                    self.oficios = deduped.map(function (el) {
                        return self.mapApiToOficio(el, withRelations);
                    });
                }

                getData(url1, onDone, true);
                getData(url2, onDone, true);
            },

            /* Normaliza la respuesta de la API a un array */
            _toArray: function (response) {
                if (Array.isArray(response)) return response;
                if (response && Array.isArray(response.items)) return response.items;
                return [];
            },

            /* ── Mapeo de respuesta API → estructura interna ────────────── */

            /**
             * Convierte un registro de la API (camelCase) al objeto interno del componente.
             *
             * @param {object}  el             - Registro de la API (EntityLetter o EntityLetterWithRelations)
             * @param {boolean} withRelations  - true si el registro incluye datos de caso y barrera
             */
            mapApiToOficio: function (el, withRelations) {
                /* ── Relaciones (solo disponibles en queries por userId) ── */
                var caseData    = null;
                var barrierData = null;

                if (withRelations) {
                    var victimName = [el.victimName, el.victimLastName]
                        .filter(function (s) { return s && s.trim(); })
                        .join(' ');

                    if (victimName || el.caseCode) {
                        caseData = {
                            victimName: victimName || '—',
                            caseCode:   el.caseCode         || el.caseId || '—',
                            docNumber:  el.victimDocNumber  || '—',
                            priority:   'Alto',
                        };
                    }

                    if (el.barrierSector) {
                        barrierData = {
                            sectors: [el.barrierSector],
                            org:     el.barrierSector,
                            desc:    el.barrierDescription || '',
                        };
                    }
                }

                return {
                    /* Identificadores */
                    id:                   el.id,
                    case_id:              el.caseId              || null,
                    barrier_id:           el.barrierId           || null,
                    /* Estado */
                    state:                el.state,
                    priority:             el.priority            || 'normal',
                    /* Auditoría */
                    agent_id:             el.agentId             || null,
                    notification_user_id: el.notificationUserId  || null,
                    review_by:            el.reviewBy            || null,
                    radicado_by:          el.radicadoBy          || null,
                    register_by:          el.registerBy          || null,
                    /* Entidad */
                    entidad:              el.entidad             || null,
                    nivel:                el.nivel               || null,
                    url_kofax:            el.urlKofax            || null,
                    correo_entidad:       el.correoEntidad       || null,
                    /* Radicación */
                    asunto_radicado:      el.asuntoRadicado      || null,
                    numero_radicado:      el.numeroRadicado      || null,
                    correo_remitente:     el.correoRemitente     || null,
                    asunto_respuesta:     el.asuntoRespuesta     || null,
                    response_date:        el.responseDate        || null,
                    response_review_by:   el.responseReviewBy    || null,
                    /* Corrección */
                    reason_correction:    el.reasonCorrection    || null,
                    /* Fechas */
                    created_at:           el.createdAt,
                    updated_at:           el.updatedAt,
                    /* Relaciones enriquecidas */
                    case:                 caseData,
                    barrier:              barrierData,
                };
            },

            /* ── Utilidades de presentación ─────────────────────────────── */

            statusLabel: function (state) {
                return STATUS_LABELS[state] || state;
            },

            nivelLabel: function (nivel) {
                return NIVEL_LABELS[nivel] || nivel || '—';
            },

            priorityLabel: function (priority) {
                return PRIORITY_LABELS[priority] || priority || '—';
            },

            stateClass: function (state) {
                return 'ol-badge--' + (state || 'default');
            },

            nivelClass: function (nivel) {
                return 'ol-badge--nivel-' + (nivel || 'default');
            },

            priorityClass: function (priority) {
                return 'ol-badge--priority-' + (priority || 'normal');
            },

            truncateUrl: function (url, maxLen) {
                maxLen = maxLen || 40;
                if (!url) return '—';
                return url.length > maxLen ? url.slice(0, maxLen) + '...' : url;
            },

            formatDate: function (isoStr) {
                if (!isoStr) return '—';
                try {
                    var d = new Date(isoStr);
                    return d.toLocaleDateString('es-CO', {
                        year: 'numeric', month: '2-digit', day: '2-digit',
                    }) + ' ' + d.toLocaleTimeString('es-CO', {
                        hour: '2-digit', minute: '2-digit',
                    });
                } catch (e) {
                    return isoStr;
                }
            },

            /* ── Modal ──────────────────────────────────────────────────── */

            openDetail: function (oficio) {
                this.selectedOficio = oficio;
            },

            closeDetail: function () {
                this.selectedOficio = null;
            },

            downloadPDF: function (oficio) {
                if (oficio.url_kofax) {
                    window.open(oficio.url_kofax, '_blank');
                }
            },

            printOficio: function () {
                window.print();
            },

            /* Recarga los datos (útil como botón de reintento) */
            retry: function () {
                this.loadOficios();
            },
        },
    };

    /* ────────────────────────────────────────────────────────────────
       Auto-montaje
       Prioridad de parámetros (de mayor a menor):
         1. data-* attributes en el div  (Go template — uso server-side)
         2. window.OficiosListConfig     (JS del padre — uso Angular-like)
    ──────────────────────────────────────────────────────────────── */

    function mount() {
        var el = document.getElementById('oficios-list-root');
        if (!el) return;

        /* Configuración publicada por la pantalla padre (ej: notifications.js) */
        var cfg = window.OficiosListConfig || {};

        var props = {
            caseId:    el.dataset.caseId    || cfg.caseId    || null,
            barrierId: el.dataset.barrierId || cfg.barrierId || null,
            userId:    el.dataset.userId    || cfg.userId    || null,
        };

        Vue.createApp(OficiosListComponent, props).mount(el);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', mount);
    } else {
        mount();
    }

    window.OficiosList = { mount: mount };

})();
