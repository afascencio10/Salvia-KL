/**
 * casos-component
 * Tabla reutilizable para visualizar y filtrar casos de víctimas.
 *
 * Props:
 *   :agentId          string (opcional)                  — icode del agente; fuerza filtro persona_asignada
 *   :defaultFilter    { key, value? }                   — filtro inicial
 *   :availableFilters Array<FilterDef>                  — filtros a mostrar en la barra
 *   :columns          Array<{ key, label }>              — columnas visibles
 *   :hiddenColumns    Array<string>                      — keys de columnas a ocultar
 *   :buttons          Array<{ id, label }>               — botones de acción por fila
 *   :reasignacion     boolean                            — modo reasignación masiva
 *
 * Emite:
 *   action-clicked    { buttonId, case }                 — botón de fila presionado
 *   reasignar-casos   { cases }                          — E-15 (pendiente en padre)
 *
 * Eventos internos implementados por archivo:
 *   E-01  mounted         — carga inicial de casos (incluye openBarriers desde barrier_v2)
 *   E-09  toggleChipFilter — filtro chip casos nuevos (hoy −5 días, America/Bogota)
 *   E-10  setDropdownFilter — filtro dropdown nivel de riesgo (combinable)
 *   E-11  setDropdownFilter — filtro dropdown por equipo (no combina con agentId)
 *   E-12  setDropdownFilter — filtro dropdown seguimientos ejecutados 0–10 (combinable)
 *   E-13  setDropdownFilter — filtro dropdown estado del caso (combinable)
 *   E-16  setDropdownFilter — filtro dropdown barreras activas (combinable)
 *   E-03  onSearchInput   — búsqueda con debounce
 *   E-04  onSortChange    — cambio de ordenamiento
 *   E-05  emitActionClicked
 *   E-06  changePage
 *   E-07  onAutocompleteInput
 *   E-08  selectAutocompleteSuggestion / clearAutocomplete
 *   E-14  toggleCaseSelection / toggleSelectAllCurrentPage
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[casos-component] No se encontró app/home. Cargue este script después de var app = home.');
        return;
    }

    vueApp.component('casos-component', {
    delimiters: ['${', '}'],

    props: {
        agentId:          { type: String,  default: '' },
        defaultFilter:    { type: Object,  required: true },
        availableFilters: { type: Array,   default: function() { return []; } },
        columns:          { type: Array,   required: true },
        hiddenColumns:    { type: Array,   default: function() { return []; } },
        buttons:          { type: Array,   default: function() { return []; } },
        reasignacion:     { type: Boolean, default: false },
    },

    emits: ['action-clicked', 'reasignar-casos'],

    data: function() {
        return {
            // --- Estado principal ---
            activeFilter:    null,
            activeChipKey:   null,
            activeDropdowns:   {}, // { riesgo: 'alto', seguimientos_ejecutados: '2', ... } — combinables
            searchText:    '',
            sortBy:        'registration_date',
            sortOrder:     'desc',
            currentPage:   1,
            pageSize:      20,
            totalCases:    0,
            cases:         [],
            loading:       true,
            loadError:     null,
            selectedCases: [],
            tableAlert: {
                visible: false,
                message: '',
            },

            // --- Estado autocomplete (indexado por filter.key) ---
            autocompleteText:        {},
            autocompleteSelected:    {},
            autocompleteSuggestions: {},
            autocompleteLoading:     {},
            autocompleteError:       {},
            autocompleteOpen:        {},

            // --- Timers internos para debounce ---
            _autocompleteTimers: {},
            _searchTimer:        null,
            _tableAlertTimer:    null,
            _onDocumentClick:    null,
        };
    },

    computed: {
        visibleColumns: function() {
            var hidden = this.hiddenColumns;
            return this.columns.filter(function(col) {
                return !hidden.includes(col.key);
            });
        },

        totalPages: function() {
            return Math.ceil(this.totalCases / this.pageSize) || 1;
        },

        filteredCases: function() {
            return this.cases;
        },

        // Filtros visibles en la barra (equipo oculto en Mis casos — E-11)
        filtersForBar: function() {
            if (!this.agentId) {
                return this.availableFilters;
            }
            return this.availableFilters.filter(function(f) {
                return f.key !== 'equipo';
            });
        },

        sortOption: function() {
            return this.sortBy + '_' + this.sortOrder;
        },

        totalCasesFormatted: function() {
            return this.totalCases.toLocaleString('es-CO');
        },

        selectedTeam: function() {
            if (!this.selectedCases.length) {
                return null;
            }
            return this.getCaseTeam(this.selectedCases[0]);
        },

        currentPageSelectableCases: function() {
            if (!this.reasignacion) {
                return [];
            }
            var team = this.selectedTeam;
            if (!team) {
                return this.filteredCases.slice();
            }
            var self = this;
            return this.filteredCases.filter(function(c) {
                return self.getCaseTeam(c) === team;
            });
        },

        isAllCurrentPageSelected: function() {
            var selectable = this.currentPageSelectableCases;
            if (!selectable.length) {
                return false;
            }
            var self = this;
            return selectable.every(function(c) {
                return self.isCaseSelected(c);
            });
        },

        isSomeCurrentPageSelected: function() {
            var self = this;
            return this.filteredCases.some(function(c) {
                return self.isCaseSelected(c);
            }) && !this.isAllCurrentPageSelected;
        },
    },

    watch: {
        reasignacion: function(val) {
            if (!val) {
                this.clearSelectedCases();
            }
        },
    },

    mounted: function() {
        if (this.agentId) {
            this.activeFilter = { key: 'persona_asignada', value: this.agentId };
        } else {
            this.activeFilter = this.normalizeDefaultActiveFilter();
        }

        // Inicializar estado de cada filtro tipo autocomplete
        var self = this;
        this.availableFilters.forEach(function(f) {
            if (f.type === 'autocomplete') {
                self.autocompleteText[f.key]        = '';
                self.autocompleteSelected[f.key]    = null;
                self.autocompleteSuggestions[f.key] = [];
                self.autocompleteLoading[f.key]     = false;
                self.autocompleteError[f.key]       = false;
                self.autocompleteOpen[f.key]        = false;
            }
        });

        this._onDocumentClick = this.handleDocumentClick.bind(this);
        document.addEventListener('click', this._onDocumentClick);

        var self = this;
        this._onWindowResize = function() {
            self.updateTableScrollWidth();
        };
        window.addEventListener('resize', this._onWindowResize);

        this.fetchCases();
    },

    beforeUnmount: function() {
        if (this._onDocumentClick) {
            document.removeEventListener('click', this._onDocumentClick);
        }
        if (this._onWindowResize) {
            window.removeEventListener('resize', this._onWindowResize);
        }
        if (this._tableResizeObserver) {
            this._tableResizeObserver.disconnect();
            this._tableResizeObserver = null;
        }
        if (this._tableAlertTimer) {
            clearTimeout(this._tableAlertTimer);
            this._tableAlertTimer = null;
        }
    },

    methods: {

        // ----------------------------------------------------------------
        // TOOLTIP CSS (chip Casos nuevos — E-09; evita Bootstrap/Popper)
        // ----------------------------------------------------------------

        chipTooltipText: function(filter) {
            if (!filter || filter.key !== 'casos_nuevos') {
                return '';
            }
            return 'Muestra los casos más recientes: creados hoy y en los 5 días calendario anteriores (hora Colombia).';
        },

        // ----------------------------------------------------------------
        // SCROLL HORIZONTAL SINCRONIZADO (barra superior + inferior)
        // ----------------------------------------------------------------

        scheduleTableScrollSync: function() {
            var self = this;
            this.$nextTick(function() {
                self.updateTableScrollWidth();
                self.observeTableScroll();
            });
        },

        updateTableScrollWidth: function() {
            var main = this.$refs.tableScrollMain;
            var top = this.$refs.tableScrollTop;
            var inner = this.$refs.tableScrollTopInner;
            if (!main || !top || !inner) {
                return;
            }

            var table = main.querySelector('.cc-table');
            if (!table) {
                top.classList.add('cc-table-scroll--hidden');
                return;
            }

            var scrollWidth = table.scrollWidth;
            inner.style.width = scrollWidth + 'px';
            inner.style.height = '1px';

            var needsScroll = scrollWidth > main.clientWidth + 1;
            var rowCount = this.cases ? this.cases.length : 0;
            var showTopScroll = needsScroll && rowCount >= 5;

            if (showTopScroll) {
                top.classList.remove('cc-table-scroll--hidden');
                top.scrollLeft = main.scrollLeft;
            } else {
                top.classList.add('cc-table-scroll--hidden');
                top.scrollLeft = 0;
                if (!needsScroll) {
                    main.scrollLeft = 0;
                }
            }
        },

        onTableScrollMain: function() {
            if (this._scrollSyncing) {
                return;
            }
            var top = this.$refs.tableScrollTop;
            var main = this.$refs.tableScrollMain;
            if (!top || !main || top.classList.contains('cc-table-scroll--hidden')) {
                return;
            }
            this._scrollSyncing = true;
            top.scrollLeft = main.scrollLeft;
            this._scrollSyncing = false;
        },

        onTableScrollTop: function() {
            if (this._scrollSyncing) {
                return;
            }
            var top = this.$refs.tableScrollTop;
            var main = this.$refs.tableScrollMain;
            if (!top || !main) {
                return;
            }
            this._scrollSyncing = true;
            main.scrollLeft = top.scrollLeft;
            this._scrollSyncing = false;
        },

        observeTableScroll: function() {
            var self = this;
            var main = this.$refs.tableScrollMain;
            if (!main) {
                return;
            }

            var table = main.querySelector('.cc-table');
            if (!table) {
                return;
            }

            if (this._tableResizeObserver) {
                this._tableResizeObserver.disconnect();
            }

            if (typeof ResizeObserver === 'undefined') {
                return;
            }

            this._tableResizeObserver = new ResizeObserver(function() {
                self.updateTableScrollWidth();
            });
            this._tableResizeObserver.observe(table);
        },

        // ----------------------------------------------------------------
        // HELPERS DE FILTRO
        // ----------------------------------------------------------------

        normalizeDefaultActiveFilter: function() {
            var df = Object.assign({}, this.defaultFilter);
            if (!df.key || df.key === 'casos_nuevos' || df.key === 'riesgo' || df.key === 'equipo' ||
                df.key === 'seguimientos_ejecutados' || df.key === 'estado_caso' || df.key === 'barreras_activas') {
                return { key: '' };
            }
            if (df.key === 'persona_asignada' && !df.value) {
                return { key: '' };
            }
            return df;
        },

        closeAutocompleteDropdown: function(key) {
            this.autocompleteOpen[key] = false;
            this.autocompleteSuggestions[key] = [];
            this.autocompleteLoading[key] = false;
        },

        handleDocumentClick: function(event) {
            if (!this.$el) {
                return;
            }

            var wrapper = event.target.closest('[data-autocomplete-key]');
            var clickedKey = wrapper ? wrapper.getAttribute('data-autocomplete-key') : null;
            var self = this;

            this.availableFilters.forEach(function(f) {
                if (f.type !== 'autocomplete' || !self.autocompleteOpen[f.key]) {
                    return;
                }
                if (clickedKey === f.key) {
                    return;
                }
                self.closeAutocompleteDropdown(f.key);
            });
        },

        // ----------------------------------------------------------------
        // HELPERS DE ORDENAMIENTO
        // ----------------------------------------------------------------

        parseSortOption: function(sortOption) {
            var valid = {
                registration_date_desc: { sortBy: 'registration_date', sortOrder: 'desc' },
                registration_date_asc:  { sortBy: 'registration_date', sortOrder: 'asc' },
                next_follow_up_desc:    { sortBy: 'next_follow_up',    sortOrder: 'desc' },
                next_follow_up_asc:     { sortBy: 'next_follow_up',    sortOrder: 'asc' },
            };
            return valid[sortOption] || null;
        },

        applySortChange: function(sortOption) {
            var parsed = this.parseSortOption(sortOption);
            if (!parsed) {
                return;
            }
            if (parsed.sortBy === this.sortBy && parsed.sortOrder === this.sortOrder) {
                return;
            }

            this.sortBy = parsed.sortBy;
            this.sortOrder = parsed.sortOrder;
            this.currentPage = 1;
            this.cases = [];
            this.fetchCases({ errorMessage: 'Error al ordenar los casos' });
        },

        // ----------------------------------------------------------------
        // E-01 — Carga inicial y recarga de casos desde el backend
        // ----------------------------------------------------------------

        clearSelectedCases: function() {
            this.selectedCases = [];
        },

        getCaseTeam: function(caseObj) {
            return (caseObj && caseObj.caseTeam) ? String(caseObj.caseTeam).trim() : '';
        },

        getCaseSelectionKey: function(caseObj) {
            if (!caseObj) {
                return '';
            }
            return caseObj.i_code || String(caseObj.id || '');
        },

        isCaseSelected: function(caseObj) {
            var key = this.getCaseSelectionKey(caseObj);
            return this.selectedCases.some(function(c) {
                return this.getCaseSelectionKey(c) === key;
            }.bind(this));
        },

        hideTableAlert: function() {
            if (this._tableAlertTimer) {
                clearTimeout(this._tableAlertTimer);
                this._tableAlertTimer = null;
            }
            this.tableAlert.visible = false;
            this.tableAlert.message = '';
        },

        showTableAlert: function(message) {
            var self = this;
            if (self._tableAlertTimer) {
                clearTimeout(self._tableAlertTimer);
            }
            self.tableAlert.message = message;
            self.tableAlert.visible = true;
            self._tableAlertTimer = setTimeout(function() {
                self.hideTableAlert();
            }, 4500);
        },

        showTeamMismatchToast: function() {
            this.showTableAlert('Solo puedes seleccionar casos de un mismo equipo.');
        },

        canSelectCase: function(caseObj) {
            var team = this.getCaseTeam(caseObj);
            if (!team) {
                return false;
            }
            if (!this.selectedCases.length) {
                return true;
            }
            return team === this.selectedTeam;
        },

        toggleCaseSelection: function(caseObj, event) {
            if (!this.reasignacion) {
                return;
            }

            var checked = event.target.checked;
            var key = this.getCaseSelectionKey(caseObj);

            if (checked) {
                if (!this.getCaseTeam(caseObj)) {
                    event.target.checked = false;
                    return;
                }
                if (!this.canSelectCase(caseObj)) {
                    event.target.checked = false;
                    this.showTeamMismatchToast();
                    return;
                }
                if (!this.isCaseSelected(caseObj)) {
                    this.selectedCases = this.selectedCases.concat([caseObj]);
                }
                return;
            }

            this.selectedCases = this.selectedCases.filter(function(c) {
                return this.getCaseSelectionKey(c) !== key;
            }.bind(this));
        },

        toggleSelectAllCurrentPage: function(event) {
            if (!this.reasignacion) {
                return;
            }

            var checked = event.target.checked;
            var pageCases = this.filteredCases;

            if (!checked) {
                var pageKeys = {};
                pageCases.forEach(function(c) {
                    pageKeys[this.getCaseSelectionKey(c)] = true;
                }.bind(this));
                this.selectedCases = this.selectedCases.filter(function(c) {
                    return !pageKeys[this.getCaseSelectionKey(c)];
                }.bind(this));
                return;
            }

            if (!pageCases.length) {
                event.target.checked = false;
                return;
            }

            var teamsOnPage = {};
            pageCases.forEach(function(c) {
                var team = this.getCaseTeam(c);
                if (team) {
                    teamsOnPage[team] = true;
                }
            }.bind(this));
            var uniqueTeams = Object.keys(teamsOnPage);

            if (uniqueTeams.length > 1 && !this.selectedCases.length) {
                event.target.checked = false;
                this.showTeamMismatchToast();
                return;
            }

            var targetTeam = this.selectedTeam || uniqueTeams[0] || '';
            if (!targetTeam) {
                event.target.checked = false;
                return;
            }

            if (this.selectedTeam && this.selectedTeam !== targetTeam) {
                event.target.checked = false;
                this.showTeamMismatchToast();
                return;
            }

            var self = this;
            var toAdd = pageCases.filter(function(c) {
                return self.getCaseTeam(c) === targetTeam && !self.isCaseSelected(c);
            });

            if (!toAdd.length && !this.isAllCurrentPageSelected) {
                event.target.checked = false;
                return;
            }

            this.selectedCases = this.selectedCases.concat(toAdd);
        },

        reload: function() {
            this.fetchCases({ errorMessage: 'Error al recargar los casos' });
        },

        fetchCases: function(options) {
            options = options || {};
            var self = this;
            self.clearSelectedCases();
            self.loading = true;
            self.loadError = null;

            var params = new URLSearchParams({
                sort:      self.sortBy,
                order:     self.sortOrder,
                page:      String(self.currentPage),
                page_size: String(self.pageSize),
            });

            if (self.agentId) {
                params.set('filter_key', 'persona_asignada');
                params.set('filter_value', self.agentId);
            } else if (self.activeFilter && self.activeFilter.key &&
                       self.activeFilter.key !== 'casos_nuevos' &&
                       self.activeFilter.key !== 'riesgo' &&
                       self.activeFilter.key !== 'equipo' &&
                       self.activeFilter.key !== 'seguimientos_ejecutados' &&
                       self.activeFilter.key !== 'estado_caso' &&
                       self.activeFilter.key !== 'barreras_activas') {
                params.set('filter_key', self.activeFilter.key);
                if (self.activeFilter.value) {
                    params.set('filter_value', self.activeFilter.value);
                }
            }

            if (self.activeChipKey === 'casos_nuevos') {
                params.set('chip_filter', 'casos_nuevos');
            }

            var dropdownParamKeys = ['riesgo', 'equipo', 'seguimientos_ejecutados', 'estado_caso', 'barreras_activas'];
            dropdownParamKeys.forEach(function(dk) {
                var val = self.activeDropdowns[dk];
                if (!val) {
                    return;
                }
                if (self.agentId && dk === 'equipo') {
                    return;
                }
                params.set('filter_' + dk, val);
            });

            var search = (self.searchText || '').trim();
            if (search) {
                params.set('search', search);
            }

            fetch('/api/v1/cases/list?' + params.toString())
                .then(function(res) {
                    return res.json().then(function(data) {
                        return { ok: res.ok, data: data };
                    });
                })
                .then(function(result) {
                    if (!result.ok) {
                        self.loadError = (result.data && result.data.error)
                            || options.errorMessage
                            || 'Error al cargar los casos';
                        self.cases = [];
                        self.totalCases = 0;
                    } else {
                        self.cases = result.data.cases || [];
                        self.totalCases = result.data.total || 0;
                        if (options.scrollToTop && self.$el) {
                            self.$el.scrollIntoView({ behavior: 'smooth', block: 'start' });
                        }
                    }
                    self.loading = false;
                    self.scheduleTableScrollSync();
                })
                .catch(function() {
                    self.loadError = options.errorMessage || 'Error al cargar los casos';
                    self.cases = [];
                    self.totalCases = 0;
                    self.loading = false;
                    self.scheduleTableScrollSync();
                });
        },

        // ----------------------------------------------------------------
        // HELPERS DE PRESENTACIÓN
        // ----------------------------------------------------------------

        formatDate: function(dateStr) {
            if (!dateStr) return '—';
            var d = new Date(dateStr);
            if (isNaN(d.getTime())) return '—';
            var dd   = String(d.getDate()).padStart(2, '0');
            var mm   = String(d.getMonth() + 1).padStart(2, '0');
            var yyyy = d.getFullYear();
            return dd + '/' + mm + '/' + yyyy;
        },

        riskBadgeClass: function(riskStatus) {
            var map = {
                'extremo':     'cc-risk-badge extremo',
                'alto':        'cc-risk-badge alto',
                'moderado':    'cc-risk-badge moderado',
                'medio':       'cc-risk-badge moderado',
                'bajo':        'cc-risk-badge bajo',
                'desconocido': 'cc-risk-badge desconocido',
            };
            return map[riskStatus] || 'cc-risk-badge desconocido';
        },

        riskBadgeLabel: function(riskStatus) {
            var map = {
                'extremo':     'Extremo',
                'alto':        'Alto',
                'moderado':    'Moderado',
                'medio':       'Moderado',
                'bajo':        'Bajo',
                'desconocido': 'Desconocido',
            };
            return map[riskStatus] || 'Desconocido';
        },

        dropdownSelectValue: function(filterKey) {
            return this.activeDropdowns[filterKey] || '';
        },

        ownerFullName: function(caseObj) {
            if (!caseObj.ownerNames) return '—';
            return (caseObj.ownerNames + ' ' + (caseObj.ownerLastNames || '')).trim();
        },

        completedFollowUpsDisplay: function(caseObj) {
            var n = caseObj.completedFollowUpsCount;
            if (n === undefined || n === null) return '0';
            return String(n);
        },

        caseStatusLabel: function(code) {
            var map = {
                ra: 'Activo',
                is: 'Con novedad',
                cd: 'Cerrado',
                ex: 'Vencido',
                r:  'Por aprobar',
                fc: 'Recontacto',
            };
            return map[code] || 'Desconocido';
        },

        caseStatusBadgeClass: function(code) {
            var known = { ra: true, is: true, cd: true, ex: true, r: true, fc: true };
            var slug = known[code] ? code : 'desconocido';
            return 'cc-case-status-badge cc-case-status-' + slug;
        },

        barrierSectorLabel: function(code) {
            var map = {
                salud: 'Salud',
                justicia: 'Justicia',
                proteccion: 'Protección',
                otras_instituciones: 'Otras instituciones',
                barrera_transversal: 'Barrera Transversal',
            };
            if (!code) {
                return '';
            }
            return map[String(code).toLowerCase()] || code;
        },

        formatOpenBarriers: function(caseObj) {
            var barriers = (caseObj && caseObj.openBarriers) ? caseObj.openBarriers : [];
            if (!barriers.length) {
                return '';
            }
            var sectors = [];
            var self = this;
            barriers.forEach(function(b) {
                var label = self.barrierSectorLabel(b.sector);
                if (label && sectors.indexOf(label) === -1) {
                    sectors.push(label);
                }
            });
            if (!sectors.length) {
                return 'Abierta';
            }
            return 'Abierta → Sector: ' + sectors.join(', ');
        },

        autocompleteMinLength: function() {
            return 3;
        },

        autocompleteShowDropdown: function(key) {
            if (!this.autocompleteOpen[key]) {
                return false;
            }
            if (this.autocompleteLoading[key]) {
                return true;
            }
            if (this.autocompleteSuggestions[key] && this.autocompleteSuggestions[key].length > 0) {
                return true;
            }
            var text = (this.autocompleteText[key] || '').trim();
            return text.length >= this.autocompleteMinLength();
        },

        autocompleteShowNoResults: function(key) {
            var text = (this.autocompleteText[key] || '').trim();
            return text.length >= this.autocompleteMinLength() &&
                !this.autocompleteLoading[key] &&
                this.autocompleteSuggestions[key] &&
                this.autocompleteSuggestions[key].length === 0;
        },

        // ----------------------------------------------------------------
        // E-09 — Filtro chip "Casos nuevos" (toggle on/off)
        // ----------------------------------------------------------------

        toggleChipFilter: function(key) {
            if (this.activeChipKey === key) {
                this.activeChipKey = null;
            } else {
                this.activeChipKey = key;
            }

            this.searchText = '';
            this.currentPage = 1;
            this.cases = [];
            this.fetchCases({ errorMessage: 'Error al aplicar el filtro' });
        },

        // ----------------------------------------------------------------
        // E-10 — Filtro dropdown (riesgo, equipo, …)
        // ----------------------------------------------------------------

        setDropdownFilter: function(key, value) {
            if (this.agentId && key === 'equipo') {
                return;
            }

            value = value || '';

            if ((this.activeDropdowns[key] || '') === value) {
                return;
            }

            var next = Object.assign({}, this.activeDropdowns);
            if (value === '') {
                delete next[key];
            } else {
                next[key] = value;
            }
            this.activeDropdowns = next;

            this.searchText = '';
            this.currentPage = 1;
            this.cases = [];
            this.fetchCases({ errorMessage: 'Error al aplicar el filtro' });
        },

        // ----------------------------------------------------------------
        // E-03 — Búsqueda con debounce
        // ----------------------------------------------------------------

        onSearchInput: function() {
            var self = this;
            if (self._searchTimer) {
                clearTimeout(self._searchTimer);
            }
            self._searchTimer = setTimeout(function() {
                self.searchText = (self.searchText || '').trim();
                self.currentPage = 1;
                self.cases = [];
                self.fetchCases({ errorMessage: 'Error al buscar casos' });
            }, 400);
        },

        // ----------------------------------------------------------------
        // E-04 — Cambio de ordenamiento
        // ----------------------------------------------------------------

        onSortChange: function(event) {
            this.applySortChange(event.target.value);
        },

        // ----------------------------------------------------------------
        // E-05 — Botón de acción en fila
        // ----------------------------------------------------------------

        emitActionClicked: function(btnId, caseObj) {
            this.$emit('action-clicked', { buttonId: btnId, case: caseObj });
        },

        // ----------------------------------------------------------------
        // E-06 — Cambio de página
        // ----------------------------------------------------------------

        changePage: function(newPage) {
            if (newPage < 1 || newPage > this.totalPages) {
                return;
            }
            if (newPage === this.currentPage) {
                return;
            }

            this.currentPage = newPage;
            this.cases = [];
            this.fetchCases({
                errorMessage: 'Error al cargar la página',
                scrollToTop:  true,
            });
        },

        // ----------------------------------------------------------------
        // E-07 — Autocomplete: input con debounce (mín. 3 caracteres)
        // ----------------------------------------------------------------

        onAutocompleteInput: function(key) {
            var self = this;
            if (self._autocompleteTimers[key]) {
                clearTimeout(self._autocompleteTimers[key]);
            }
            self._autocompleteTimers[key] = setTimeout(function() {
                self.runAutocompleteSearch(key);
            }, 400);
        },

        runAutocompleteSearch: function(key) {
            var self = this;
            var queryText = (self.autocompleteText[key] || '').trim();
            var minLen = self.autocompleteMinLength();

            if (queryText === '') {
                self.autocompleteSuggestions[key] = [];
                self.autocompleteLoading[key] = false;
                self.autocompleteError[key] = false;
                self.autocompleteOpen[key] = false;
                return;
            }

            if (queryText.length < minLen) {
                self.autocompleteSuggestions[key] = [];
                self.autocompleteLoading[key] = false;
                self.autocompleteError[key] = false;
                self.autocompleteOpen[key] = false;
                return;
            }

            self.autocompleteOpen[key] = true;
            self.autocompleteLoading[key] = true;
            self.autocompleteSuggestions[key] = [];
            self.autocompleteError[key] = false;

            var params = new URLSearchParams({
                q:     queryText,
                limit: '10',
            });

            fetch('/api/v1/agents/search?' + params.toString())
                .then(function(res) {
                    return res.json().then(function(data) {
                        return { ok: res.ok, data: data };
                    });
                })
                .then(function(result) {
                    if (!result.ok) {
                        self.autocompleteError[key] = true;
                        self.autocompleteSuggestions[key] = [];
                    } else {
                        var agents = result.data.agents || [];
                        self.autocompleteSuggestions[key] = agents.map(function(a) {
                            return {
                                value: a.icode,
                                label: ((a.names || '') + ' ' + (a.lastNames || '')).trim(),
                            };
                        });
                    }
                    self.autocompleteLoading[key] = false;
                })
                .catch(function() {
                    self.autocompleteError[key] = true;
                    self.autocompleteSuggestions[key] = [];
                    self.autocompleteLoading[key] = false;
                });
        },

        // ----------------------------------------------------------------
        // E-08 — Autocomplete: seleccionar o limpiar persona asignada
        // ----------------------------------------------------------------

        selectAutocompleteSuggestion: function(key, opt) {
            if (!opt || !opt.value || this.agentId) {
                return;
            }

            this.autocompleteSelected[key] = {
                value: opt.value,
                label: opt.label,
            };
            this.autocompleteText[key] = '';
            this.autocompleteError[key] = false;
            this.closeAutocompleteDropdown(key);

            this.activeFilter = { key: key, value: opt.value };
            this.searchText = '';
            this.currentPage = 1;
            this.cases = [];
            this.fetchCases({ errorMessage: 'Error al filtrar por persona asignada' });
        },

        emitReasignarCasos: function() {
            if (!this.reasignacion || !this.selectedCases.length) {
                return;
            }
            this.$emit('reasignar-casos', { cases: this.selectedCases.slice() });
        },

        clearAutocomplete: function(key) {
            if (this.agentId) {
                return;
            }

            this.autocompleteSelected[key] = null;
            this.autocompleteText[key] = '';
            this.autocompleteError[key] = false;
            this.closeAutocompleteDropdown(key);

            this.activeFilter = this.normalizeDefaultActiveFilter();
            this.searchText = '';
            this.currentPage = 1;
            this.cases = [];
            this.fetchCases({ errorMessage: 'Error al cargar los casos' });
        },
    },

    // Template capturado en la pantalla host (window.__casosComponentTpl) antes de mount,
    // porque #app reemplaza el DOM y el x-template dentro de body deja de existir.
    template: (typeof window !== 'undefined' && window.__casosComponentTpl)
        ? window.__casosComponentTpl
        : '#tpl-casos-component',
    });
})();
