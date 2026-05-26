/**
 * notifications.js — Lógica Vue de la pantalla de Notificaciones Salvia (EntityLetter)
 *
 * Las variables de sesión (currentUser, currentRole, currentUserId) se inyectan
 * desde el template Go como window.NotifConfig antes de cargar este archivo.
 */

(function () {
    'use strict';

    const cfg = window.NotifConfig || {};

    var home = Vue.createApp({
        delimiters: ['${', '}'],
        data() {
            return {
                currentUser:   cfg.currentUser   || '',
                currentRole:   cfg.currentRole   || '',
                currentUserId: cfg.currentUserId || '',

                isLoading: false,
                loadError: null,

                currentTab: 'todos',

                filters: {
                    estado:    '',
                    identidad: '',
                    entidad:   '',
                    radicado:  ''
                },

                itemsPerPage: 5,
                currentPage:  0,

                showModal:      false,
                activeModal:    null,
                selectedOficio: null,

                modalForm: {
                    nivel:                '',
                    kofaxPath:            '',
                    prioridad:            'normal',
                    entidad:              '',
                    asunto:               '',
                    correoEntidad:        '',
                    numeroRadicado:       '',
                    fechaRespuesta:       '',
                    correoRemitente:      '',
                    asuntoRespuesta:      '',
                    respuestaRecibidaPor: '',
                    reasonCorrection:     ''
                },

                isSaving:  false,
                saveError: null,

                modalByStatus: {
                    por_proyectar:       'proyectar',
                    para_revisar:        'revisar',
                    en_correccion:       'corregir',
                    aprobacion_juridica: 'aprobar',
                    para_radicar:        'radicar',
                    radicado:            'registrar_respuesta'
                },

                statusLabels: {
                    por_proyectar:       'Por proyectar',
                    para_revisar:        'Para revisar',
                    en_correccion:       'En corrección (proyección)',
                    aprobacion_juridica: 'Aprobación jurídica',
                    para_radicar:        'Para radicar',
                    radicado:            'Radicado',
                    respondido:          'Respondido'
                },

                priorityLabels: {
                    normal: 'Normal',
                    alta:   'Alta'
                },

                oficios: []
            };
        },

        mounted() {
            document.getElementById('app').style.display = 'block';
            this.loadOficios();
        },

        computed: {
            pendingCount() {
                return this.oficios.filter(o => o.canManage).length;
            },

            filteredOficios() {
                return this.oficios.filter(o => {
                    const tabOk      = this.currentTab === 'todos' || o.canManage;
                    const estadoOk   = !this.filters.estado || o.status === this.filters.estado;
                    const idOk       = !this.filters.identidad || o.docNumber.includes(this.filters.identidad);
                    const entidadOk  = !this.filters.entidad ||
                        o.barrierOrg.toLowerCase().includes(this.filters.entidad.toLowerCase());
                    const radicadoOk = !this.filters.radicado ||
                        (o.numeroRadicado && o.numeroRadicado.toLowerCase().includes(this.filters.radicado.toLowerCase()));
                    return tabOk && estadoOk && idOk && entidadOk && radicadoOk;
                });
            },

            numPages() {
                return Math.max(1, Math.ceil(this.filteredOficios.length / this.itemsPerPage));
            },

            paginatedOficios() {
                const start = this.currentPage * this.itemsPerPage;
                return this.filteredOficios.slice(start, start + this.itemsPerPage);
            },

            visiblePages() {
                const max    = this.numPages;
                const cur    = this.currentPage;
                const delta  = 2;
                const range  = [];
                const result = [];
                let   last;

                for (let i = 0; i < max; i++) {
                    if (i === 0 || i === max - 1 || (i >= cur - delta && i <= cur + delta)) {
                        range.push(i);
                    }
                }
                for (let i of range) {
                    if (last !== undefined) {
                        if (i - last === 2) result.push(last + 1);
                        else if (i - last > 2) result.push('...');
                    }
                    result.push(i);
                    last = i;
                }
                return result;
            }
        },

        methods: {

            /* ── Carga inicial ──────────────────────────────────────────────────── */

            loadOficios() {
                const VALID_ROLES = ['op', 'an'];
                if (!VALID_ROLES.includes(this.currentRole)) {
                    this.loadError = 'Rol no autorizado para acceder a esta pantalla.';
                    return;
                }

                this.isLoading = true;
                this.loadError = null;

                let url = '/api/v1/entity-letters?limit=100&page=0';
                if (this.currentRole === 'op') {
                    url += '&agentId=' + encodeURIComponent(this.currentUserId);
                } else if (this.currentRole === 'an') {
                    url += '&notificationUserId=' + encodeURIComponent(this.currentUserId);
                }

                getData(url, (function (status, response) {
                    this.isLoading = false;

                    if (status === 200) {
                        const items = Array.isArray(response)
                            ? response
                            : (response && response.items ? response.items : []);

                        if (items.length > 0) {
                            this.oficios = items
                                .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
                                .map(el => this.mapApiToOficio(el));
                        }
                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else {
                        this.loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.';
                    }
                }).bind(this), true);
            },

            canManageForRole(state) {
                if (this.currentRole === 'op') {
                    return ['por_proyectar', 'en_correccion'].includes(state);
                }
                if (this.currentRole === 'an') {
                    return ['para_revisar', 'aprobacion_juridica', 'para_radicar', 'radicado'].includes(state);
                }
                return false;
            },

            mapApiToOficio(el) {
                const fecha = el.createdAt ? el.createdAt.split('T')[0] : null;

                const victimFull = [el.victimName, el.victimLastName]
                    .filter(s => s && s.trim())
                    .join(' ') || '—';

                const sectors = (el.barrierSector && el.barrierSector.trim())
                    ? [el.barrierSector]
                    : [];

                const rawPriority = (el.priority || 'normal').toLowerCase().trim();
                const letterPriority = rawPriority === 'alta' ? 'alta' : 'normal';

                return {
                    id:              el.id,
                    victimName:      victimFull,
                    caseCode:        el.caseCode       || el.caseId   || '—',
                    docNumber:       el.victimDocNumber || '—',
                    caseId:          el.caseId         || '',
                    priority:        'Alto',
                    letterPriority:  letterPriority,
                    sectors:         sectors,
                    barrierOrg:      el.barrierSector      || '—',
                    barrierDesc:     el.barrierDescription || '—',
                    status:          el.state,
                    canManage:       this.canManageForRole(el.state),
                    title:           '—',
                    radicado:        null,
                    fecha:           fecha,
                    correo:          '—',
                    entidad:         el.entidad          || null,
                    nivel:           el.nivel            || null,
                    urlKofax:        el.urlKofax         || null,
                    asuntoRadicado:  el.asuntoRadicado   || null,
                    correoEntidad:   el.correoEntidad    || null,
                    numeroRadicado:  el.numeroRadicado   || null,
                    reasonCorrection: el.reasonCorrection || null
                };
            },

            /* ── Utilidades ─────────────────────────────────────────────────────── */

            slugify(str) {
                if (!str) return '';
                return str.toLowerCase()
                    .normalize('NFD')
                    .replace(/[\u0300-\u036f]/g, '')
                    .replace(/\s+/g, '-');
            },

            priorityLabel(priority) {
                return this.priorityLabels[priority] || priority;
            },

            switchTab(tab) {
                this.currentTab  = tab;
                this.currentPage = 0;
            },

            goToPage(page) {
                if (page < 0 || page >= this.numPages) return;
                this.currentPage = page;
                window.scrollTo(0, 0);
            },

            /* ── Modal de gestión ───────────────────────────────────────────────── */

            openGestionarModal(oficio) {
                console.log('[Notificaciones] Gestionar oficio:', oficio);

                this.selectedOficio = oficio;
                this.activeModal    = this.modalByStatus[oficio.status] || null;

                this.modalForm = {
                    nivel:                oficio.nivel             || '',
                    kofaxPath:            oficio.urlKofax          || '',
                    prioridad:            oficio.letterPriority    || 'normal',
                    entidad:              oficio.entidad           || oficio.barrierOrg || '',
                    asunto:               oficio.asuntoRadicado    || oficio.title || '',
                    correoEntidad:        oficio.correoEntidad     || oficio.correo || '',
                    numeroRadicado:       oficio.numeroRadicado    || '',
                    fechaRespuesta:       '',
                    correoRemitente:      oficio.correo            || '',
                    asuntoRespuesta:      '',
                    respuestaRecibidaPor: '',
                    reasonCorrection:     ''
                };

                this.saveError = null;
                this.showModal = true;
            },

            closeModal() {
                this.showModal                  = false;
                this.activeModal                = null;
                this.selectedOficio             = null;
                this.saveError                  = null;
                this.modalForm.reasonCorrection = '';
            },

            submitModal(action) {
                if (!this.selectedOficio) return;

                let payload = {
                    action: action,
                    userId: this.currentUserId
                };

                if (action === 'proyectar') {
                    if (!this.modalForm.nivel) {
                        alert('El campo "Nivel" es requerido.');
                        return;
                    }
                    if (!this.modalForm.entidad) {
                        alert('El campo "Entidad" es requerido.');
                        return;
                    }
                    if (!this.modalForm.kofaxPath) {
                        alert('El campo "Ruta del oficio en el Kofax" es requerido.');
                        return;
                    }
                    payload.nivel    = this.modalForm.nivel;
                    payload.entidad  = this.modalForm.entidad;
                    payload.urlKofax = this.modalForm.kofaxPath;
                    payload.priority = this.modalForm.prioridad || 'normal';
                }

                if (action === 'por_corregir') {
                    if (!this.modalForm.reasonCorrection) {
                        alert('El campo "Razón de corrección" es requerido para marcar el oficio por corregir.');
                        return;
                    }
                    payload.reasonCorrection = this.modalForm.reasonCorrection;
                }

                if (action === 'registrar_respuesta') {
                    if (!this.modalForm.fechaRespuesta) {
                        alert('El campo "Fecha de respuesta" es requerido.');
                        return;
                    }
                    if (!this.modalForm.correoRemitente) {
                        alert('El campo "Correo del remitente" es requerido.');
                        return;
                    }
                    if (!this.modalForm.asuntoRespuesta) {
                        alert('El campo "Asunto de la respuesta" es requerido.');
                        return;
                    }
                    if (!this.modalForm.respuestaRecibidaPor) {
                        alert('El campo "Respuesta recibida por" es requerido.');
                        return;
                    }
                    payload.responseDate     = this.modalForm.fechaRespuesta;
                    payload.correoRemitente  = this.modalForm.correoRemitente;
                    payload.asuntoRespuesta  = this.modalForm.asuntoRespuesta;
                    payload.responseReviewBy = this.modalForm.respuestaRecibidaPor;
                }

                if (action === 'radicar') {
                    if (!this.modalForm.asunto) {
                        alert('El campo "Asunto" es requerido.');
                        return;
                    }
                    if (!this.modalForm.correoEntidad) {
                        alert('El campo "Correo entidad" es requerido.');
                        return;
                    }
                    if (!this.modalForm.numeroRadicado) {
                        alert('El campo "Número radicado" es requerido.');
                        return;
                    }
                    payload.asuntoRadicado = this.modalForm.asunto;
                    payload.correoEntidad  = this.modalForm.correoEntidad;
                    payload.numeroRadicado = this.modalForm.numeroRadicado;
                }

                this.isSaving  = true;
                this.saveError = null;

                const url = '/api/v1/entity-letters/' + this.selectedOficio.id + '/action';

                updateEntity(url, payload, (function (status, response) {
                    this.isSaving = false;

                    if (status === 200) {
                        setTimeout(() => {
                            const overlay = document.getElementById('successOverlay');
                            if (overlay) overlay.style.display = 'none';
                        }, 1200);

                        const idx = this.oficios.findIndex(o => o.id === this.selectedOficio.id);
                        if (idx !== -1) {
                            const updated = response;
                            this.oficios[idx] = {
                                ...this.oficios[idx],
                                status:    updated.state,
                                canManage: this.canManageForRole(updated.state)
                            };
                        }
                        this.closeModal();
                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else if (status === 422) {
                        this.saveError = (response && response.error) ? response.error : 'Transición de estado no permitida.';
                    } else {
                        this.saveError = (response && response.error)
                            ? response.error
                            : 'No se pudo guardar el oficio. Intenta de nuevo.';
                    }
                }).bind(this));
            },

            logout() {
                logout('/seguridad/logout');
            }
        },

        watch: {
            filters: {
                deep: true,
                handler() { this.currentPage = 0; }
            }
        }
    });

    home.mount('#app');
})();
