/**
 * CustomSelect — Componente Vue 3 accesible
 *
 * Patrón ARIA: Combobox con listbox (WAI-ARIA 1.2)
 * https://www.w3.org/WAI/ARIA/apg/patterns/combobox/
 *
 * Correcciones aplicadas:
 *  1. tabindex="-1" en botones × para que el foco Tab siempre llegue al combobox primero.
 *  2. onComboboxBlur usa setTimeout + nextTick para tolerar el ciclo del lector de pantalla.
 *  3. stopPropagation en todos los keydown relevantes para evitar interferencia del lector.
 *  4. scrollIntoView sin behavior:smooth para evitar eventos intermedios de foco.
 *  5. tabindex="-1" en los <li> de opciones para que no sean focusables por Tab/lector.
 *  6. Dos regiones live: polite para estado, assertive para navegación entre opciones.
 *  7. announceNavigation usa assertive para interrumpir al lector inmediatamente.
 *  8. announce y announceNavigation limpian con \u00A0 antes para forzar relecura.
 *  9. moveUp/moveDown anuncian posición y estado de selección de la opción.
 * 10. isOpening guard: evita que onComboboxBlur cierre el menú durante la apertura,
 *     cuando open() mueve el foco del combobox al searchEl (condición de carrera).
 */

const { ref, computed, nextTick, onBeforeUnmount } = Vue;

const CustomSelect = {
    directives: {
        clickOutside: {
            mounted(el, binding) {
                el.__clickOutside__ = (e) => {
                    if (!el.contains(e.target)) binding.value(e);
                };
                document.addEventListener('mousedown', el.__clickOutside__);
            },
            unmounted(el) {
                document.removeEventListener('mousedown', el.__clickOutside__);
            }
        }
    },

    template: `
        <!-- Región polite: apertura, cierre, resultados de búsqueda -->
        <span
            role="status"
            aria-live="polite"
            aria-atomic="true"
            class="sr-only"
        >{{ liveMessage }}</span>

        <!-- Región assertive: navegación activa entre opciones -->
        <span
            role="alert"
            aria-live="assertive"
            aria-atomic="true"
            class="sr-only"
        >{{ liveNavigationMessage }}</span>

        <div
            class="custom-select"
            :class="{ active: isOpen }"
            v-click-outside="close"
            ref="rootEl"
        >
            <!--
                Área visual del control: muestra placeholder, valor single o tags multiple.
                aria-hidden="true" porque toda la semántica va en el botón combobox.
                @click="toggle" aquí para que el área visible siempre abra el dropdown,
                incluso cuando hay tags y el trigger tiene pointer-events reducidos.
            -->
            <div class="custom-select-face" aria-hidden="true" @click="toggle">

                <!-- SINGLE: placeholder o valor seleccionado -->
                <template v-if="mode === 'single'">
                    <span v-if="!selectedOption || !selectedOption[valueKey]" class="custom-select-placeholder">
                        {{ placeholder }}
                    </span>
                    <span v-else class="custom-select-single-value">
                        {{ getOptionLabel(selectedOption) }}
                    </span>
                </template>

                <!-- MULTIPLE: placeholder o tags con × -->
                <template v-else>
                    <span v-if="selectedOptions.length === 0" class="custom-select-placeholder">
                        {{ placeholder }}
                    </span>
                    <span
                        v-for="opt in selectedOptions"
                        :key="getOptionValue(opt)"
                        class="custom-select-tag"
                    >
                        <span class="custom-select-tag-label">{{ getOptionLabel(opt) }}</span>
                        <!--
                            tabindex="-1": el botón × es alcanzable programáticamente y con
                            click/lector, pero NO entra en el flujo de Tab. Así el foco de Tab
                            siempre llega al combobox trigger primero y el lector lee el label.
                        -->
                        <button
                            type="button"
                            class="custom-select-tag-remove"
                            :aria-label="'Eliminar ' + getOptionLabel(opt)"
                            @click.stop="removeOption(opt)"
                            @blur="onComboboxBlur"
                            tabindex="-1"
                        >×</button>
                    </span>
                </template>

                <svg
                    class="custom-select-arrow"
                    :class="{ open: isOpen }"
                    width="16" height="16" viewBox="0 0 16 16"
                    aria-hidden="true" focusable="false"
                >
                    <path d="M4 6L8 10L12 6" stroke="currentColor" stroke-width="2" fill="none"/>
                </svg>
            </div>

            <!--
                COMBOBOX: botón overlay que recibe el foco real y los eventos de teclado.
                En modo multiple con tags, el CSS puede reducir su área visual, pero
                el click principal se gestiona también desde custom-select-face arriba.
            -->
            <button
                type="button"
                class="custom-select-trigger"
                :class="{ 'has-tags': mode === 'multiple' && selectedOptions.length > 0 }"
                role="combobox"
                :aria-expanded="isOpen"
                aria-haspopup="listbox"
                :aria-controls="listboxId"
                :aria-activedescendant="activeDescendant || undefined"
                :aria-label="comboboxAriaLabel"
                :aria-multiselectable="mode === 'multiple' ? 'true' : undefined"
                @click="toggle"
                @keydown="onComboboxKeydown"
                @blur="onComboboxBlur"
                ref="comboboxEl"
            ></button>

            <transition name="dropdown">
                <div
                    v-if="isOpen"
                    class="custom-select-dropdown"
                    ref="dropdownEl"
                >
                    <!--
                        BUSCADOR: searchbox independiente fuera del listbox.
                        Anuncia el número de resultados al escribir.
                    -->
                    <div
                        v-if="options.length > 2"
                        class="custom-select-search"
                        role="search"
                    >
                        <label :for="searchId" class="sr-only">{{ searchPlaceholder }}</label>
                        <input
                            type="text"
                            :id="searchId"
                            v-model="searchQuery"
                            :placeholder="searchPlaceholder"
                            autocomplete="off"
                            :aria-controls="listboxId"
                            :aria-activedescendant="activeDescendant || undefined"
                            @keydown="onSearchKeydown"
                            @blur="onComboboxBlur"
                            ref="searchEl"
                        />
                    </div>

                    <!--
                        LISTBOX: el foco lógico se comunica vía aria-activedescendant
                        desde el combobox/search. El foco real del DOM nunca entra aquí.
                        Los <li> tienen tabindex="-1" para que el lector no los recorra
                        con su cursor virtual de forma independiente.
                    -->
                    <ul
                        :id="listboxId"
                        role="listbox"
                        :aria-multiselectable="mode === 'multiple'"
                        :aria-label="placeholder"
                        class="custom-select-options"
                        ref="listboxEl"
                    >
                        <li
                            v-if="filteredOptions.length === 0"
                            role="option"
                            aria-disabled="true"
                            aria-selected="false"
                            class="custom-select-no-results"
                        >
                            Sin resultados
                        </li>
                        <li
                            v-for="(option, index) in filteredOptions"
                            :key="getOptionValue(option)"
                            :id="optionId(index)"
                            role="option"
                            :aria-selected="isSelected(option)"
                            class="custom-select-option"
                            :class="{
                                selected: isSelected(option),
                                focused: index === focusedIndex
                            }"
                            tabindex="-1"
                            @mousedown.prevent="selectOption(option)"
                        >
                            <span
                                :class="mode === 'single'
                                    ? 'custom-select-radio'
                                    : 'custom-select-checkbox'"
                                aria-hidden="true"
                            ></span>
                            <span class="custom-select-option-label">
                                {{ getOptionLabel(option) }}
                            </span>
                        </li>
                    </ul>

                    <div class="custom-select-actions">
                        <button
                            type="button"
                            class="custom-select-btn custom-select-btn-clear"
                            :aria-label="mode === 'single'
                                ? 'Limpiar selección'
                                : 'Limpiar todas las selecciones'"
                            @mousedown.prevent="clearAll"
                            @blur="onComboboxBlur"
                        >
                            Limpiar
                        </button>
                    </div>
                </div>
            </transition>
        </div>
    `,

    props: {
        modelValue:        { default: null },
        options:           { type: Array, required: true },
        mode:              { type: String, default: 'multiple' },
        placeholder:       { type: String, default: 'Selecciona opciones...' },
        searchPlaceholder: { type: String, default: 'Buscar...' },
        valueKey:          { type: String, default: 'value' },
        labelKey:          { type: String, default: 'label' },
        returnObject:      { type: Boolean, default: false },
        ariaLabel:         { type: String, default: null }
    },

    emits: ['update:modelValue'],

    setup(props, { emit }) {
        /* ── Refs DOM ─────────────────────────────────────────────── */
        const rootEl     = ref(null);
        const comboboxEl = ref(null);
        const dropdownEl = ref(null);
        const searchEl   = ref(null);
        const listboxEl  = ref(null);

        /* ── Estado ───────────────────────────────────────────────── */
        const isOpen       = ref(false);
        const isOpening    = ref(false); // Guard: evita que onComboboxBlur cierre durante apertura
        const searchQuery  = ref('');
        const focusedIndex = ref(-1);

        /* ── Regiones live ────────────────────────────────────────── */
        // polite  → apertura, cierre, resultados de búsqueda
        // assertive → navegación activa entre opciones (interrumpe al lector)
        const liveMessage           = ref('');
        const liveNavigationMessage = ref('');

        /* ── IDs únicos por instancia ─────────────────────────────── */
        const uid       = Math.random().toString(36).slice(2, 9);
        const listboxId = `cs-listbox-${uid}`;
        const searchId  = `cs-search-${uid}`;
        const optionId  = (i) => `cs-opt-${uid}-${i}`;

        /* ── Helpers ──────────────────────────────────────────────── */
        const getOptionValue = (opt) => opt[props.valueKey];
        const getOptionLabel = (opt) => opt[props.labelKey];

        /* ── Opciones filtradas ───────────────────────────────────── */
        const filteredOptions = computed(() => {
            const q = searchQuery.value.trim().toLowerCase();
            if (!q) return props.options;
            return props.options.filter(o =>
                getOptionLabel(o).toLowerCase().includes(q)
            );
        });

        /* ── Selección ────────────────────────────────────────────── */
        const selectedOption = computed(() => {
            if (props.mode !== 'single' || !props.modelValue) return null;
            if (props.returnObject) return props.modelValue;
            return props.options.find(o => getOptionValue(o) === props.modelValue) ?? null;
        });

        const selectedOptions = computed(() => {
            if (props.mode !== 'multiple' || !Array.isArray(props.modelValue)) return [];
            if (props.returnObject) return props.modelValue;
            return props.options.filter(o => props.modelValue.includes(getOptionValue(o)));
        });

        const isSelected = (opt) => {
            if (props.mode === 'single') {
                if (!props.modelValue) return false;
                return props.returnObject
                    ? getOptionValue(opt) === getOptionValue(props.modelValue)
                    : getOptionValue(opt) === props.modelValue;
            }
            if (!Array.isArray(props.modelValue)) return false;
            return props.returnObject
                ? props.modelValue.some(v => getOptionValue(v) === getOptionValue(opt))
                : props.modelValue.includes(getOptionValue(opt));
        };

        /* ── Etiqueta accesible del combobox ─────────────────────── */
        const comboboxAriaLabel = computed(() => {
            if (props.ariaLabel) return props.ariaLabel;
            if (props.mode === 'single') {
                return selectedOption.value
                    ? `${props.placeholder}: ${getOptionLabel(selectedOption.value)}`
                    : props.placeholder;
            }
            if (selectedOptions.value.length === 0) return props.placeholder;
            return `${props.placeholder}: ${selectedOptions.value.map(getOptionLabel).join(', ')}`;
        });

        /* ── aria-activedescendant ────────────────────────────────── */
        const activeDescendant = computed(() =>
            isOpen.value && focusedIndex.value >= 0
                ? optionId(focusedIndex.value)
                : null
        );

        /* ── Anuncios ─────────────────────────────────────────────── */
        // Limpia con \u00A0 primero para forzar al lector a detectar el cambio
        // aunque el mensaje nuevo sea igual al anterior.
        const announce = (msg) => {
            liveMessage.value = '\u00A0';
            nextTick(() => { liveMessage.value = msg; });
        };

        const announceNavigation = (msg) => {
            liveNavigationMessage.value = '\u00A0';
            nextTick(() => { liveNavigationMessage.value = msg; });
        };

        /* ── Scroll para mantener la opción enfocada visible ─────── */
        // Sin behavior:'smooth' para evitar eventos intermedios de foco
        // que algunos lectores de pantalla interpretan como cambios de foco.
        const scrollFocusedIntoView = () => {
            nextTick(() => {
                const el = listboxEl.value?.querySelector('.custom-select-option.focused');
                el?.scrollIntoView({ block: 'nearest' });
            });
        };

        /* ── Apertura / cierre ────────────────────────────────────── */
        const open = async () => {
            if (isOpen.value) return;
            isOpen.value   = true;
            isOpening.value = true; // Activa guard antes de mover el foco
            searchQuery.value = '';

            // Pre-enfocar la opción actualmente seleccionada en modo single
            if (props.mode === 'single' && props.modelValue) {
                const val = props.returnObject ? getOptionValue(props.modelValue) : props.modelValue;
                focusedIndex.value = filteredOptions.value.findIndex(o => getOptionValue(o) === val);
            } else {
                focusedIndex.value = -1;
            }

            await nextTick();
            if (props.options.length > 2 && searchEl.value) {
                searchEl.value.focus();
            }

            // Desactiva el guard después de que el blur del combobox
            // (disparado por el .focus() de searchEl) ya fue procesado.
            // 200ms > 150ms del setTimeout de onComboboxBlur.
            setTimeout(() => { isOpening.value = false; }, 200);

            announce(
                `Lista desplegada. ${filteredOptions.value.length} opciones. ` +
                `Use las flechas para navegar, Enter para seleccionar, Escape para cerrar.`
            );
        };

        const close = () => {
            if (!isOpen.value) return;
            isOpen.value = false;
            searchQuery.value = '';
            focusedIndex.value = -1;
            announce('Lista cerrada.');
        };

        const toggle = () => isOpen.value ? close() : open();

        /* ── Navegación por teclado (foco lógico) ─────────────────── */
        const moveUp = (opts) => {
            focusedIndex.value = focusedIndex.value <= 0
                ? opts.length - 1
                : focusedIndex.value - 1;
            const opt = opts[focusedIndex.value];
            const pos = focusedIndex.value + 1;
            const sel = isSelected(opt) ? ', seleccionado' : '';
            // assertive para interrumpir cualquier lectura en curso del lector
            announceNavigation(`${getOptionLabel(opt)}${sel}, ${pos} de ${opts.length}`);
            scrollFocusedIntoView();
        };

        const moveDown = (opts) => {
            focusedIndex.value = (focusedIndex.value + 1) % opts.length;
            const opt = opts[focusedIndex.value];
            const pos = focusedIndex.value + 1;
            const sel = isSelected(opt) ? ', seleccionado' : '';
            // assertive para interrumpir cualquier lectura en curso del lector
            announceNavigation(`${getOptionLabel(opt)}${sel}, ${pos} de ${opts.length}`);
            scrollFocusedIntoView();
        };

        /* ── Keydown del combobox (botón principal) ───────────────── */
        const onComboboxKeydown = (e) => {
            const opts = filteredOptions.value;

            switch (e.key) {
                case 'Enter':
                case ' ':
                    e.preventDefault();
                    e.stopPropagation();
                    if (!isOpen.value) { open(); return; }
                    if (focusedIndex.value >= 0) selectOption(opts[focusedIndex.value]);
                    break;
                case 'ArrowDown':
                    e.preventDefault();
                    e.stopPropagation();
                    if (!isOpen.value) { open(); return; }
                    if (opts.length) moveDown(opts);
                    break;
                case 'ArrowUp':
                    e.preventDefault();
                    e.stopPropagation();
                    if (!isOpen.value) { open(); return; }
                    if (opts.length) moveUp(opts);
                    break;
                case 'Escape':
                    e.preventDefault();
                    e.stopPropagation();
                    close();
                    comboboxEl.value?.focus();
                    break;
                case 'Tab':
                    close();
                    break;
            }
        };

        /* ── Keydown del buscador ─────────────────────────────────── */
        const onSearchKeydown = (e) => {
            const opts = filteredOptions.value;

            switch (e.key) {
                case 'ArrowDown':
                    e.preventDefault();
                    e.stopPropagation();
                    if (opts.length) moveDown(opts);
                    break;
                case 'ArrowUp':
                    e.preventDefault();
                    e.stopPropagation();
                    if (opts.length) moveUp(opts);
                    break;
                case 'Enter':
                    e.preventDefault();
                    e.stopPropagation();
                    if (focusedIndex.value >= 0) selectOption(opts[focusedIndex.value]);
                    break;
                case 'Escape':
                    e.preventDefault();
                    e.stopPropagation();
                    close();
                    comboboxEl.value?.focus();
                    break;
                case 'Tab':
                    close();
                    break;
            }
        };

        /* ── Blur: robusto para lectores de pantalla ──────────────── */
        // Los lectores de pantalla pueden mover el foco del sistema brevemente
        // durante su ciclo de lectura, haciendo que relatedTarget llegue como null
        // aunque el foco real siga dentro del widget. El setTimeout de 150ms
        // espera a que el lector termine su ciclo antes de verificar el foco real.
        // nextTick dentro garantiza que Vue también haya procesado cambios pendientes.
        // isOpening guard: si open() acaba de mover el foco a searchEl, ignoramos
        // el blur del combobox para no cerrar el menú recién abierto.
        const onComboboxBlur = (e) => {
            setTimeout(() => {
                nextTick(() => {
                    if (isOpening.value) return; // Guard activo: apertura en curso
                    const active = document.activeElement;
                    if (rootEl.value && rootEl.value.contains(active)) return;
                    close();
                });
            }, 150);
        };

        /* ── Seleccionar / deseleccionar ─────────────────────────── */
        const selectOption = (opt) => {
            const label = getOptionLabel(opt);
            const val   = props.returnObject ? opt : getOptionValue(opt);

            if (props.mode === 'single') {
                emit('update:modelValue', val);
                announce(`${label} seleccionado.`);
                close();
                comboboxEl.value?.focus();
            } else {
                const current = Array.isArray(props.modelValue) ? props.modelValue : [];
                let nowSelected;

                if (props.returnObject) {
                    const idx = current.findIndex(i => getOptionValue(i) === getOptionValue(opt));
                    nowSelected = idx < 0;
                    emit('update:modelValue', nowSelected
                        ? [...current, opt]
                        : current.filter((_, i) => i !== idx)
                    );
                } else {
                    const v = getOptionValue(opt);
                    nowSelected = !current.includes(v);
                    emit('update:modelValue', nowSelected
                        ? [...current, v]
                        : current.filter(x => x !== v)
                    );
                }

                announceNavigation(`${label} ${nowSelected ? 'seleccionado' : 'deseleccionado'}.`);
            }
        };

        const removeOption = (opt) => {
            if (props.mode !== 'multiple') return;
            const current = Array.isArray(props.modelValue) ? props.modelValue : [];
            const label   = getOptionLabel(opt);

            if (props.returnObject) {
                emit('update:modelValue', current.filter(i => getOptionValue(i) !== getOptionValue(opt)));
            } else {
                emit('update:modelValue', current.filter(v => v !== getOptionValue(opt)));
            }
            announce(`${label} eliminado.`);
        };

        const clearAll = () => {
            emit('update:modelValue', props.mode === 'single' ? null : []);
            announce('Selección limpiada.');
            close();
            comboboxEl.value?.focus();
        };

        /* ── Anunciar resultados de búsqueda al escribir ─────────── */
        let searchAnnounceTimer = null;
        const watchSearch = computed(() => searchQuery.value);

        Vue.watchEffect(() => {
            if (!watchSearch.value && watchSearch.value !== '') return;
            const count = filteredOptions.value.length;
            clearTimeout(searchAnnounceTimer);
            searchAnnounceTimer = setTimeout(() => {
                if (isOpen.value) announce(`${count} resultado${count !== 1 ? 's' : ''}.`);
            }, 300);
        });

        onBeforeUnmount(() => {
            clearTimeout(searchAnnounceTimer);
            isOpening.value = false; // Limpieza defensiva
        });

        return {
            // Refs DOM
            rootEl, comboboxEl, dropdownEl, searchEl, listboxEl,
            // Estado
            isOpen, isOpening, searchQuery, focusedIndex,
            // Regiones live
            liveMessage, liveNavigationMessage,
            // IDs
            listboxId, searchId, optionId,
            // Computed
            filteredOptions, selectedOption, selectedOptions,
            isSelected, comboboxAriaLabel, activeDescendant,
            // Helpers
            getOptionValue, getOptionLabel,
            // Acciones
            toggle, close, selectOption, removeOption, clearAll,
            // Eventos
            onComboboxKeydown, onSearchKeydown, onComboboxBlur
        };
    }
};