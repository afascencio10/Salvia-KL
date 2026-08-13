/**
 * barreras_departamento.js
 * Lógica Vue para la pantalla "Barreras Departamento" (barreras_departamento.html).
 *
 * Requiere que el HTML haya expuesto previamente la configuración del servidor
 * en el objeto global window.__BarrerasDepartamentoConfig con las propiedades:
 *   - currentUser   {string}
 *   - currentRole   {string}
 *   - currentUserId {string}
 */

(function () {
    'use strict';

    const cfg = window.__BarrerasDepartamentoConfig || {};

    const riskLevelLabels = { 1: 'Bajo', 2: 'Moderado', 3: 'Alto', 4: 'Extremo' };
    const estadoLabels = { OPEN: 'Articulando', 'En Gestion': 'Articulando', Articulada: 'Resuelto' };

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data() {
            return {
                currentUser:   cfg.currentUser   || '',
                currentRole:   cfg.currentRole   || '',
                currentUserId: cfg.currentUserId || '',

                isLoading: false,
                loadError: null,

                filters: {
                    identidad: '',
                    entidad:   '',
                    ciudad:    ''
                },

                itemsPerPage: 10,
                currentPage:  0,

                /* ── Datos cargados desde la API ───────────────────────────── */
                barreras: []
            };
        },

        mounted() {
            document.getElementById('app').style.display = 'block';
            this.loadBarreras();
        },

        /* ── Computed ──────────────────────────────────────────────────────── */
        computed: {

            filteredBarreras() {
                return this.barreras.filter(b => {
                    const idOk     = !this.filters.identidad || b.docNumber.includes(this.filters.identidad);
                    const entidadOk = !this.filters.entidad ||
                        b.entityBranchName.toLowerCase().includes(this.filters.entidad.toLowerCase());
                    const ciudadOk = !this.filters.ciudad ||
                        b.cityName.toLowerCase().includes(this.filters.ciudad.toLowerCase());
                    return idOk && entidadOk && ciudadOk;
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

            loadBarreras() {
                /* Solo el enlace territorial (en) puede acceder a esta pantalla */
                if (this.currentRole !== 'en') {
                    this.loadError = 'Rol no autorizado para acceder a esta pantalla.';
                    return;
                }

                this.isLoading = true;
                this.loadError = null;

                getData('/api/v1/barriers-v2/department', (function (status, response) {
                    this.isLoading = false;

                    if (status === 200) {
                        const items = Array.isArray(response) ? response : [];
                        this.barreras = items.map(b => this.mapApiToBarrera(b));

                    } else if (status === 401) {
                        location.assign('/static/landing.html');
                    } else if (status === 403) {
                        this.loadError = 'Tu usuario no tiene un departamento asignado. Contacta a un administrador para que te lo asigne.';
                    } else {
                        this.loadError = 'No se pudo cargar la lista de barreras. Intenta de nuevo.';
                    }
                }).bind(this), true);
            },

            /* Transforma un BarrierV2DepartmentItem de la API al formato de la tabla */
            mapApiToBarrera(b) {
                const victimFull = [b.victimName, b.victimLastName]
                    .filter(s => s && s.trim())
                    .join(' ') || '—';

                return {
                    id:                b.id,
                    caseId:            b.caseId || '',
                    victimName:        victimFull,
                    caseCode:          b.caseCode || b.caseId || '—',
                    docNumber:         b.victimDocNumber || '—',
                    sector:            b.sector || '—',
                    description:       b.description || '—',
                    entityBranchName:  b.entityBranchName || '—',
                    cityName:          b.cityName || '—',
                    estadoLabel:       estadoLabels[b.status] || b.status || '—',
                    priority:          riskLevelLabels[b.riskLevel] || 'Desconocido',
                    fechaIdentificacion: this.formatFecha(b.barrierDate || b.createdAt)
                };
            },

            formatFecha(value) {
                if (!value) return '—';
                const d = new Date(value);
                if (isNaN(d.getTime())) return value;
                return d.toLocaleDateString('es-CO', { year: 'numeric', month: '2-digit', day: '2-digit' });
            },

            /* ── Utilidades ─────────────────────────────────────────────── */

            slugify(str) {
                if (!str) return '';
                return str.toLowerCase()
                    .normalize('NFD')
                    .replace(/[̀-ͯ]/g, '')
                    .replace(/\s+/g, '-');
            },

            goToPage(page) {
                if (page < 0 || page >= this.numPages) return;
                this.currentPage = page;
                window.scrollTo(0, 0);
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
