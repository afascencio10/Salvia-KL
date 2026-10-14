/**
 * barriers.js
 * Lógica Vue para la pantalla "Mis Barreras" (barrieris.html).
 *
 * Expone `home` / `app` para registrar case-task-modal antes del mount.
 * Requiere window.__MisBarrerasConfig: currentUser, currentRole, currentUserId.
 */

(function () {
    'use strict';

    const cfg = window.__MisBarrerasConfig || {};

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data() {
            return {
                currentUser:   cfg.currentUser   || '',
                currentRole:   cfg.currentRole   || '',
                currentUserId: cfg.currentUserId || '',
                assignedDepartmentId: cfg.assignedDepartmentId || '',

                isLoading: false,
                loadError: null,

                filters: {
                    identidad: '',
                    entidad:   '',
                    cityId:    ''
                },

                cities: [],

                itemsPerPage: 5,
                currentPage:  0,

                barreras: [],
                successToast: null
            };
        },

        mounted() {
            document.getElementById('app').style.display = 'block';
            this.loadCities();
            this.loadBarreras();
        },

        computed: {
            filteredBarreras() {
                return this.barreras.filter(b => {
                    const idOk = !this.filters.identidad ||
                        (b.docNumber && b.docNumber.includes(this.filters.identidad));
                    const entidadOk = !this.filters.entidad ||
                        (b.sector && b.sector.toLowerCase().includes(this.filters.entidad.toLowerCase())) ||
                        (b.entityBranchName && b.entityBranchName.toLowerCase().includes(this.filters.entidad.toLowerCase()));
                    const cityOk = !this.filters.cityId ||
                        String(b.cityId) === String(this.filters.cityId);
                    return idOk && entidadOk && cityOk;
                });
            },

            numPages() {
                return Math.max(1, Math.ceil(this.filteredBarreras.length / this.itemsPerPage));
            },

            paginatedBarreras() {
                const start = this.currentPage * this.itemsPerPage;
                return this.filteredBarreras.slice(start, start + this.itemsPerPage);
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
            loadCities() {
                if (!this.assignedDepartmentId) {
                    this.cities = [];
                    return;
                }
                fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(this.assignedDepartmentId), {
                    credentials: 'same-origin'
                })
                    .then(res => res.ok ? res.json() : [])
                    .then(data => {
                        this.cities = Array.isArray(data) ? data : [];
                    })
                    .catch(() => {
                        this.cities = [];
                    });
            },

            loadBarreras() {
                /* El GET de la pantalla ya exige get_mis_barreras. El API
                   /barriers-v2/mine vuelve a validar rol en. No bloquear aquí:
                   un currentRole mal inyectado dejaba el listado vacío. */
                this.isLoading = true;
                this.loadError = null;

                getData('/api/v1/barriers-v2/mine', (function (status, response) {
                    this.isLoading = false;

                    if (status === 200) {
                        const items = Array.isArray(response)
                            ? response
                            : (response && response.items ? response.items : []);
                        this.barreras = items.map(row => this.mapApiToBarrera(row));
                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else if (status === 403) {
                        this.loadError = 'No tienes permiso para ver estas barreras.';
                    } else {
                        this.loadError = 'No se pudo cargar la lista de barreras. Intenta de nuevo.';
                    }
                }).bind(this), true);
            },

            mapApiToBarrera(row) {
                const victimFull = [row.victimName, row.victimLastName]
                    .filter(s => s && String(s).trim())
                    .join(' ') || '—';

                return {
                    id:               row.id,
                    barrierId:        row.id,
                    caseId:           row.caseId || '',
                    victimName:       victimFull,
                    docNumber:        row.victimDocNumber || '—',
                    victimPhone:      (row.victimPhone && String(row.victimPhone).trim()) || '—',
                    cityId:           row.cityId || '',
                    cityName:         row.cityName || '—',
                    sector:           row.sector || '—',
                    description:      row.description || '—',
                    entityBranchId:   row.entityBranchId || null,
                    entityBranchName: row.entityBranchName || '',
                    taskId:           row.taskId || '',
                    taskTypeLabel:    this.taskTypeLabel(row.taskType),
                    taskDescription:  row.taskDescription || ''
                };
            },

            taskTypeLabel(type) {
                const labels = {
                    articulacion_territorial: 'Articulación territorial',
                    compromiso_con_enlace:    'Compromiso con Enlace',
                    gestion_propia:           'Gestión propia'
                };
                if (!type) return '';
                if (labels[type]) return labels[type];
                return type.replace(/_/g, ' ').replace(/^./, c => c.toUpperCase());
            },

            slugify(str) {
                if (!str) return '';
                return str.toLowerCase()
                    .normalize('NFD')
                    .replace(/[\u0300-\u036f]/g, '')
                    .replace(/\s+/g, '-');
            },

            goToPage(page) {
                if (page < 0 || page >= this.numPages) return;
                this.currentPage = page;
                window.scrollTo(0, 0);
            },

            openGestionar(barrera) {
                if (!this.$refs.taskModal) return;
                this.$refs.taskModal.openArticulacion({
                    taskId:             barrera.taskId || null,
                    barrierId:          barrera.barrierId,
                    caseId:             barrera.caseId,
                    entityBranchId:     barrera.entityBranchId,
                    entityName:         barrera.entityBranchName,
                    barrierDescription: barrera.description
                });
            },

            onTaskCompleted() {
                this.showSuccessToast('Gestión guardada.');
                this.loadBarreras();
            },

            showSuccessToast(message) {
                this.successToast = message;
                setTimeout(() => { this.successToast = null; }, 3000);
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

    window.home = home;
    window.app = home;
})();
