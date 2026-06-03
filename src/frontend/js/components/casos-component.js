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
 *
 * Emite:
 *   action-clicked    { buttonId, case }                 — botón de fila presionado
 *
 * Eventos internos implementados por archivo:
 *   E-01  mounted         — carga inicial de casos
 *   E-09  toggleChipFilter — filtro chip casos nuevos (hoy −5 días, America/Bogota)
 *   E-10  setDropdownFilter — filtro dropdown nivel de riesgo (combinable)
 *   E-11  setDropdownFilter — filtro dropdown por equipo (no combina con agentId)
 *   E-12  setDropdownFilter — filtro dropdown seguimientos ejecutados 0–10 (combinable)
 *   E-03  onSearchInput   — búsqueda con debounce
 *   E-04  onSortChange    — cambio de ordenamiento
 *   E-05  emitActionClicked
 *   E-06  changePage
 *   E-07  onAutocompleteInput
 *   E-08  selectAutocompleteSuggestion / clearAutocomplete
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
    },

    emits: ['action-clicked'],

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

        this.fetchCases();
    },

    beforeUnmount: function() {
        if (this._onDocumentClick) {
            document.removeEventListener('click', this._onDocumentClick);
        }
    },

    methods: {

        // ----------------------------------------------------------------
        // HELPERS DE FILTRO
        // ----------------------------------------------------------------

        normalizeDefaultActiveFilter: function() {
            var df = Object.assign({}, this.defaultFilter);
            if (!df.key || df.key === 'casos_nuevos' || df.key === 'riesgo' || df.key === 'equipo' ||
                df.key === 'seguimientos_ejecutados') {
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

        fetchCases: function(options) {
            options = options || {};
            var self = this;
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
                       self.activeFilter.key !== 'seguimientos_ejecutados') {
                params.set('filter_key', self.activeFilter.key);
                if (self.activeFilter.value) {
                    params.set('filter_value', self.activeFilter.value);
                }
            }

            if (self.activeChipKey === 'casos_nuevos') {
                params.set('chip_filter', 'casos_nuevos');
            }

            var dropdownParamKeys = ['riesgo', 'equipo', 'seguimientos_ejecutados'];
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
                })
                .catch(function() {
                    self.loadError = options.errorMessage || 'Error al cargar los casos';
                    self.cases = [];
                    self.totalCases = 0;
                    self.loading = false;
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
