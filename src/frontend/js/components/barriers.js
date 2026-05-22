/**
 * barriers.js
 * Lógica Vue para la pantalla "Mis Barreras" (barrieris.html).
 *
 * Requiere que el HTML haya expuesto previamente la configuración del servidor
 * en el objeto global window.__MisBarrerasConfig con las propiedades:
 *   - currentUser   {string}
 *   - currentRole   {string}
 *   - currentUserId {string}
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

                isLoading: false,
                loadError: null,

                currentTab: 'por_articular',

                filters: {
                    identidad: '',
                    entidad:   ''
                },

                itemsPerPage: 5,
                currentPage:  0,

                /* ── Estado del modal ──────────────────────────────────────── */
                showModal:       false,
                selectedBarrera: null,
                isSaving:        false,
                saveError:       null,
                modalForm: {
                    gestion: ''
                },

                /* ── Datos cargados desde la API ───────────────────────────── */
                barreras: [],

                /* ── Feedback de éxito ─────────────────────────────────────── */
                successToast: null
            };
        },

        mounted() {
            document.getElementById('app').style.display = 'block';
            this.loadBarreras();
        },

        /* ── Computed ──────────────────────────────────────────────────────── */
        computed: {

            porArticularCount() {
                return this.barreras.filter(b => !b.articulada).length;
            },

            articuladasCount() {
                return this.barreras.filter(b => b.articulada).length;
            },

            filteredBarreras() {
                return this.barreras.filter(b => {
                    const tabOk    = this.currentTab === 'por_articular' ? !b.articulada : b.articulada;
                    const idOk     = !this.filters.identidad || b.docNumber.includes(this.filters.identidad);
                    const sectorOk = !this.filters.entidad ||
                        b.sector.toLowerCase().includes(this.filters.entidad.toLowerCase());
                    return tabOk && idOk && sectorOk;
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

        /* ── Methods ───────────────────────────────────────────────────────── */
        methods: {

            /* ── Evento: Cuando carga la pantalla ────────────────────────── */

            loadBarreras() {
                console.log('loadBarreras');
                /* Solo el enlace territorial (en) puede acceder a esta pantalla */
                if (this.currentRole !== 'en') {
                    this.loadError = 'Rol no autorizado para acceder a esta pantalla.';
                    return;
                }

                this.isLoading = true;
                this.loadError = null;

                const url = '/api/v1/case-tasks?assignedUserId=' + encodeURIComponent(this.currentUserId);

                getData(url, (function (status, response) {
                    this.isLoading = false;

                    if (status === 200) {
                        const items = Array.isArray(response)
                            ? response
                            : (response && response.items ? response.items : []);

                        this.barreras = items.map(ct => this.mapApiToBarrera(ct));

                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else {
                        this.loadError = 'No se pudo cargar la lista de barreras. Intenta de nuevo.';
                    }
                }).bind(this), true);
            },

            /* Transforma un CaseTaskWithRelations de la API al formato de la tabla */
            mapApiToBarrera(ct) {
                const victimFull = [ct.victimName, ct.victimLastName]
                    .filter(s => s && s.trim())
                    .join(' ') || '—';

                return {
                    id:          ct.id,
                    barrierId:   ct.barrierId || '',
                    victimName:  victimFull,
                    caseCode:    ct.caseCode || ct.caseId || '—',
                    docNumber:   ct.victimDocNumber || '—',
                    caseId:      ct.caseId || '',
                    priority:    'Alto',
                    sector:      ct.barrierSector || '—',
                    description: ct.barrierDescription || ct.description || '—',
                    gestion:     ct.result || null,
                    articulada:  ct.status === 'Done',
                    canManage:   ct.status === 'ToDo'
                };
            },

            /* ── Utilidades ─────────────────────────────────────────────── */

            slugify(str) {
                if (!str) return '';
                return str.toLowerCase()
                    .normalize('NFD')
                    .replace(/[\u0300-\u036f]/g, '')
                    .replace(/\s+/g, '-');
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

            /* ── Modal Gestionar ─────────────────────────────────────────── */

            openGestionar(barrera) {
                this.selectedBarrera   = barrera;
                this.modalForm.gestion = '';
                this.saveError         = null;
                this.showModal         = true;
            },

            closeModal() {
                this.showModal         = false;
                this.selectedBarrera   = null;
                this.saveError         = null;
                this.modalForm.gestion = '';
            },

            submitGestionar() {
                if (!this.modalForm.gestion || !this.modalForm.gestion.trim()) {
                    this.saveError = 'La descripción de la gestión es requerida.';
                    return;
                }

                this.isSaving  = true;
                this.saveError = null;

                const url  = '/api/v1/case-tasks/' + this.selectedBarrera.id + '/complete';
                const body = JSON.stringify({ result: this.modalForm.gestion });

                fetch(url, {
                    method:  'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body:    body
                })
                .then(res => res.json().then(data => ({ status: res.status, data })))
                .then(({ status, data }) => {
                    this.isSaving = false;

                    if (status === 200) {
                        /* Actualizar la fila localmente:
                           · CaseTask → articulada / sin botón Gestionar
                           · Gestión realizada visible en la columna */
                        const idx = this.barreras.findIndex(b => b.id === this.selectedBarrera.id);
                        if (idx !== -1) {
                            this.barreras[idx].articulada = true;
                            this.barreras[idx].canManage  = false;
                            this.barreras[idx].gestion    = this.modalForm.gestion;
                        }
                        this.closeModal();
                        this.showSuccessToast('Barrera articulada exitosamente.');
                    } else {
                        this.saveError = (data && data.error)
                            ? data.error
                            : 'No se pudo guardar la gestión. Intenta de nuevo.';
                    }
                })
                .catch(() => {
                    this.isSaving  = false;
                    this.saveError = 'Error de conexión. Intenta de nuevo.';
                });
            },

            /* Muestra un toast de éxito temporal (3 s) */
            showSuccessToast(message) {
                this.successToast = message;
                setTimeout(() => { this.successToast = null; }, 3000);
            },

            logout() {
                logout('/seguridad/logout');
            }
        },

        /* ── Watchers ──────────────────────────────────────────────────────── */
        watch: {
            filters: {
                deep: true,
                handler() { this.currentPage = 0; }
            }
        }
    });

    home.mount('#app');
})();
