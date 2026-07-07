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
                totalItems:   0,
                pendingCount: 0,

                filterDebounceTimer: null,

                showModal:      false,
                activeModal:    null,
                selectedOficio: null,

                // Locaciones cargadas al montar (E05 background)
                allDepartments: [],
                allCities:      [],

                modalForm: {
                    // Campos del modal proyectar v2
                    departmentId:         '',
                    cityId:               '',
                    townId:               '',
                    entityBranchId:       '',
                    entityName:           '',
                    officialDependency:   '',
                    subject:              '',
                    cities:               [],
                    towns:                [],
                    entityBranches:       [],
                    // Campos comunes
                    kofaxPath:            '',
                    prioridad:            'normal',
                    // Campos otros modales
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
            this.loadLocations();
        },

        computed: {
            numPages() {
                return Math.max(1, Math.ceil(this.totalItems / this.itemsPerPage));
            },

            paginatedOficios() {
                return this.oficios;
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

            /* ── Locaciones (carga en background al montar) ─────────────────────── */

            loadLocations() {
                var self = this;
                Promise.all([
                    new Promise(function (resolve) {
                        getData('/api/v1/locations/departments', function (status, response) {
                            if (status === 200) {
                                var items = Array.isArray(response) ? response : (response && response.items ? response.items : []);
                                self.allDepartments = items;
                            }
                            resolve();
                        }, true);
                    }),
                    new Promise(function (resolve) {
                        getData('/api/v1/locations/cities', function (status, response) {
                            if (status === 200) {
                                var items = Array.isArray(response) ? response : (response && response.items ? response.items : []);
                                self.allCities = items;
                            }
                            resolve();
                        }, true);
                    })
                ]);
            },

            /* ── Carga inicial ──────────────────────────────────────────────────── */

            buildOficiosUrl() {
                const params = new URLSearchParams();
                params.set('page', String(this.currentPage));
                params.set('limit', String(this.itemsPerPage));

                if (this.currentRole === 'op' || this.currentRole === 'ro') {
                    params.set('agentId', this.currentUserId);
                } else if (this.currentRole === 'an') {
                    params.set('notificationUserId', this.currentUserId);
                }

                if (this.currentTab === 'gestionar') {
                    params.set('manageableOnly', 'true');
                }
                if (this.filters.estado) {
                    params.set('state', this.filters.estado);
                }
                if (this.filters.identidad.trim()) {
                    params.set('identidad', this.filters.identidad.trim());
                }
                if (this.filters.entidad.trim()) {
                    params.set('entidad', this.filters.entidad.trim());
                }
                if (this.filters.radicado.trim()) {
                    params.set('numeroRadicado', this.filters.radicado.trim());
                }

                return '/api/v1/entity-letters?' + params.toString();
            },

            loadOficios() {
                const VALID_ROLES = ['op', 'an', 'ro'];
                if (!VALID_ROLES.includes(this.currentRole)) {
                    this.loadError = 'Rol no autorizado para acceder a esta pantalla.';
                    return;
                }

                this.isLoading = true;
                this.loadError = null;

                const url = this.buildOficiosUrl();

                getData(url, (function (status, response) {
                    this.isLoading = false;

                    if (status === 200) {
                        const items = Array.isArray(response)
                            ? response
                            : (response && response.items ? response.items : []);

                        this.oficios = items.map(el => this.mapApiToOficio(el));

                        if (response && typeof response.total === 'number') {
                            this.totalItems = response.total;
                        } else {
                            this.totalItems = items.length;
                        }
                        if (response && typeof response.pendingCount === 'number') {
                            this.pendingCount = response.pendingCount;
                        }

                        if (this.currentPage > 0 && this.currentPage >= this.numPages) {
                            this.currentPage = Math.max(0, this.numPages - 1);
                            this.loadOficios();
                        }
                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else {
                        this.loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.';
                    }
                }).bind(this), true);
            },

            canManageForRole(state) {
                if (this.currentRole === 'op' || this.currentRole === 'ro') {
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
                    barrierId:       el.barrierId      || '',
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
                    entidad:              el.entidad             || null,
                    nivel:                el.nivel               || null,
                    urlKofax:             el.urlKofax            || null,
                    asuntoRadicado:       el.asuntoRadicado      || null,
                    correoEntidad:        el.correoEntidad       || null,
                    numeroRadicado:       el.numeroRadicado      || null,
                    reasonCorrection:     el.reasonCorrection    || null,
                    officialDependency:   el.officialDependency  || null,
                    subject:              el.subject             || null,
                    townName:             el.townName            || null
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
                this.loadOficios();
            },

            goToPage(page) {
                if (page < 0 || page >= this.numPages) return;
                this.currentPage = page;
                window.scrollTo(0, 0);
                this.loadOficios();
            },

            /* ── Modal de gestión ───────────────────────────────────────────────── */

            openGestionarModal(oficio) {
                console.log('[NF] openGestionarModal id=' + oficio.id + ' status=' + oficio.status);

                this.selectedOficio = oficio;
                this.activeModal    = this.modalByStatus[oficio.status] || null;

                this.modalForm = {
                    // Proyectar v2 — siempre vacío (por_proyectar nunca regresa a ese estado)
                    departmentId:         '',
                    cityId:               '',
                    townId:               '',
                    entityBranchId:       '',
                    entityName:           '',
                    officialDependency:   '',
                    subject:              '',
                    cities:               [],
                    towns:                [],
                    entityBranches:       [],
                    // Campos comunes
                    kofaxPath:            oficio.urlKofax          || '',
                    prioridad:            oficio.letterPriority    || 'normal',
                    // Otros modales
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
                this.showModal          = false;
                this.activeModal        = null;
                this.selectedOficio     = null;
                this.saveError          = null;
                this.modalForm.reasonCorrection = '';
            },

            submitModal(action) {
                if (!this.selectedOficio) return;

                let payload = {
                    action: action,
                    userId: this.currentUserId
                };

                if (action === 'proyectar') {
                    if (!this.modalForm.departmentId) {
                        alert('El campo "Departamento" es requerido.');
                        return;
                    }
                    if (!this.modalForm.cityId) {
                        alert('El campo "Ciudad" es requerido.');
                        return;
                    }
                    if (!this.modalForm.townId) {
                        alert('El campo "Municipio" es requerido.');
                        return;
                    }
                    if (!this.modalForm.entityBranchId) {
                        alert('El campo "Entidad" es requerido.');
                        return;
                    }
                    var entityName = '';
                    if (this.modalForm.entityBranchId === 'otra') {
                        if (!this.modalForm.entityName) {
                            alert('El campo "Nombre de la entidad" es requerido.');
                            return;
                        }
                        entityName = this.modalForm.entityName;
                    } else {
                        var selected = this.modalForm.entityBranches.find(function (e) {
                            return e.id === this.modalForm.entityBranchId;
                        }.bind(this));
                        entityName = selected ? selected.name : '';
                    }
                    if (!this.modalForm.officialDependency) {
                        alert('El campo "Funcionario" es requerido.');
                        return;
                    }
                    if (!this.modalForm.subject) {
                        alert('El campo "Asunto" es requerido.');
                        return;
                    }
                    if (!this.modalForm.kofaxPath) {
                        alert('El campo "Ruta del oficio en el Kofax" es requerido.');
                        return;
                    }
                    payload.departmentId       = this.modalForm.departmentId;
                    payload.cityId             = this.modalForm.cityId;
                    payload.townId             = this.modalForm.townId;
                    payload.entityBranchId     = this.modalForm.entityBranchId !== 'otra' ? this.modalForm.entityBranchId : null;
                    payload.entityName         = entityName;
                    payload.officialDependency = this.modalForm.officialDependency;
                    payload.subject            = this.modalForm.subject;
                    payload.urlKofax           = this.modalForm.kofaxPath;
                    payload.priority           = this.modalForm.prioridad || 'normal';
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
                console.log('[NF] submitModal url=' + url + ' payload=' + JSON.stringify(payload));

                updateEntity(url, payload, (function (status, response) {
                    this.isSaving = false;

                    if (status === 200) {
                        setTimeout(() => {
                            const overlay = document.getElementById('successOverlay');
                            if (overlay) overlay.style.display = 'none';
                        }, 1200);

                        this.closeModal();
                        this.loadOficios();
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

            /* ── E08: Cuando selecciona departamento en modal proyectar ─────────── */

            onDepartmentChange() {
                var deptId = this.modalForm.departmentId;
                console.log('[NF] onDepartmentChange deptId=' + deptId + ' allCities=' + this.allCities.length);
                var filtered = deptId
                    ? this.allCities.filter(function (c) { return String(c.departmentId) === String(deptId); })
                    : [];
                console.log('[NF] cities filtered=' + filtered.length);
                this.modalForm.cities        = filtered;
                this.modalForm.cityId        = '';
                this.modalForm.townId        = '';
                this.modalForm.entityBranchId = '';
                this.modalForm.entityName    = '';
                this.modalForm.towns         = [];
                this.modalForm.entityBranches = [];
            },

            /* ── E09: Cuando selecciona ciudad en modal proyectar ───────────────── */

            onCityChange() {
                var cityId = this.modalForm.cityId;
                this.modalForm.townId         = '';
                this.modalForm.entityBranchId = '';
                this.modalForm.entityName     = '';
                this.modalForm.towns          = [];
                this.modalForm.entityBranches = [];

                if (!cityId) return;

                var self = this;
                getData('/api/v1/locations/towns?city_id=' + encodeURIComponent(cityId), function (status, response) {
                    if (status === 200) {
                        self.modalForm.towns = Array.isArray(response) ? response : (response && response.items ? response.items : []);
                    }
                }, true);
            },

            /* ── E10: Cuando selecciona municipio en modal proyectar ────────────── */

            onTownChange() {
                var townId = this.modalForm.townId;
                this.modalForm.entityBranchId = '';
                this.modalForm.entityName     = '';
                this.modalForm.entityBranches = [];

                if (!townId) return;

                var self = this;
                getData('/api/v1/entity-branches?town_code=' + encodeURIComponent(townId), function (status, response) {
                    if (status === 200) {
                        self.modalForm.entityBranches = Array.isArray(response) ? response : [];
                        console.log('[NF] entityBranches town_code=' + townId + ' count=' + self.modalForm.entityBranches.length + ' first=' + JSON.stringify(self.modalForm.entityBranches[0]));
                    }
                }, true);
            },

            logout() {
                logout('/seguridad/logout');
            }
        },

        watch: {
            filters: {
                deep: true,
                handler() {
                    this.currentPage = 0;
                    clearTimeout(this.filterDebounceTimer);
                    this.filterDebounceTimer = setTimeout(() => {
                        this.loadOficios();
                    }, 350);
                }
            }
        }
    });

    home.mount('#app');

    /* ── Configuración para el componente hijo OficiosList ──────────────────
     *
     * Angular-like: el "componente padre" (esta pantalla) le pasa parámetros
     * al "componente hijo" (oficios-list) antes de que éste cargue y se monte.
     *
     * El orden de ejecución garantiza que este objeto esté disponible cuando
     * oficios-list.js lee window.OficiosListConfig en su función mount().
     *
     * Para filtrar por caso:    { caseId:    '<uuid>' }
     * Para filtrar por barrera: { barrierId: '<uuid>' }
     * Para filtrar por usuario: { userId:    '<icode>' }
     * ─────────────────────────────────────────────────────────────────────── */
    /*window.OficiosListConfig = {
        caseId:    '019d020a-16cf-7c25-8744-5c1fd6f43a72',
        barrierId: '34e367b4-ea59-4936-b080-5e5491a60380',
    };*/

})();
