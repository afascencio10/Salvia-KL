/**
 * case-timeline
 * Componente reutilizable para renderizar el timeline de eventos de un caso.
 * Props:
 *   - caseId    (String, required): icode del caso
 *   - barrierId (String, optional): filtra eventos por barrier_id
 */

/* ─── Inyección de estilos ──────────────────────────────────────────────── */
(function injectCaseTimelineStyles() {
    if (document.getElementById('case-timeline-styles')) return;
    const style = document.createElement('style');
    style.id = 'case-timeline-styles';
    style.textContent = `
        .ct-wrapper {
            display: flex;
            flex-direction: column;
            gap: 0;
        }

        /* ── Filtros ── */
        .ct-filters {
            display: flex;
            flex-wrap: wrap;
            gap: 6px;
            margin-bottom: 16px;
        }
        .ct-filter-btn {
            padding: 5px 12px;
            border-radius: 999px;
            font-size: 12px;
            font-weight: 500;
            border: 1px solid #e5e7eb;
            background: #fff;
            color: #6b7280;
            cursor: pointer;
            transition: background 0.15s, border-color 0.15s, color 0.15s;
            white-space: nowrap;
        }
        .ct-filter-btn:hover:not(.active) {
            background: #f9fafb;
            border-color: #d1d5db;
        }
        .ct-filter-btn.active {
            background: #eff6ff;
            border-color: #93c5fd;
            color: #1d4ed8;
        }
        .ct-filter-count {
            display: inline-block;
            background: #e5e7eb;
            color: #6b7280;
            border-radius: 999px;
            font-size: 10px;
            font-weight: 600;
            padding: 0 6px;
            margin-left: 4px;
            min-width: 18px;
            text-align: center;
        }
        .ct-filter-btn.active .ct-filter-count {
            background: #bfdbfe;
            color: #1d4ed8;
        }

        /* ── Timeline vertical ── */
        .ct-list {
            display: flex;
            flex-direction: column;
            position: relative;
        }
        .ct-list::before {
            content: '';
            position: absolute;
            left: 19px;
            top: 0;
            bottom: 0;
            width: 2px;
            background: #e5e7eb;
            z-index: 0;
        }

        /* ── Item ── */
        .ct-item {
            display: flex;
            gap: 14px;
            align-items: flex-start;
            position: relative;
            padding-bottom: 20px;
        }
        .ct-item:last-child {
            padding-bottom: 0;
        }

        /* ── Icono ── */
        .ct-icon-wrap {
            flex-shrink: 0;
            width: 38px;
            height: 38px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 15px;
            color: #fff;
            z-index: 1;
            position: relative;
            box-shadow: 0 0 0 3px #fff;
        }

        /* ── Card ── */
        .ct-card {
            flex: 1;
            min-width: 0;
            background: #fff;
            border: 1px solid #e5e7eb;
            border-radius: 10px;
            padding: 10px 14px;
            box-shadow: 0 1px 2px 0 rgb(0 0 0 / .04);
        }
        .ct-card-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            margin-bottom: 4px;
            flex-wrap: wrap;
        }
        .ct-type-badge {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            font-size: 11px;
            font-weight: 700;
            padding: 3px 10px;
            border-radius: 999px;
            white-space: nowrap;
            letter-spacing: 0.01em;
        }
        .ct-date {
            font-size: 11px;
            color: #9ca3af;
            white-space: nowrap;
        }
        .ct-desc {
            font-size: 13px;
            color: #374151;
            line-height: 1.5;
            word-break: break-word;
        }
        .ct-actor {
            margin-top: 6px;
            display: flex;
            align-items: center;
            gap: 5px;
            font-size: 11px;
            color: #6b7280;
        }

        /* ── Estados ── */
        .ct-empty {
            text-align: center;
            padding: 32px 16px;
            color: #9ca3af;
            font-size: 14px;
        }
        .ct-empty-icon {
            font-size: 28px;
            margin-bottom: 8px;
        }
        .ct-loading {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 10px;
            padding: 32px 0;
            color: #9ca3af;
            font-size: 14px;
        }
        @keyframes ct-spin { to { transform: rotate(360deg); } }
        .ct-spinner {
            width: 18px; height: 18px;
            animation: ct-spin 0.8s linear infinite;
            flex-shrink: 0;
        }
        .ct-error {
            padding: 14px 16px;
            background: #fef2f2;
            border: 1px solid #fecaca;
            border-radius: 8px;
            color: #b91c1c;
            font-size: 13px;
        }
    `;
    document.head.appendChild(style);
})();

/* ─── Categorías del timeline ───────────────────────────────────────────── */
var CT_CATEGORIAS = [
    { value: '',               label: 'Todos' },
    { value: 'General',        label: 'General' },
    { value: 'Barreras',       label: 'Barreras' },
    { value: 'Seguimientos',   label: 'Seguimientos' },
    { value: 'Oficios',        label: 'Oficios' },
    { value: 'Medidas',        label: 'Medidas' },
    { value: 'Psicosocial',    label: 'Psicosocial' },
    { value: 'Estabilización', label: 'Estabilización' },
];

/* ─── Helper: color de fondo por defecto cuando el evento no trae color ── */
function ctDefaultColor(category) {
    var map = {
        'Seguimientos':   '#1d4ed8',
        'Barreras':       '#b91c1c',
        'Medidas':        '#c2410c',
        'Psicosocial':    '#5106A7',
        'Estabilización': '#15803d',
        'Oficios':        '#0f766e',
        'General':        '#4b5563',
    };
    return map[category] || '#4b5563';
}

/* ─── Helper: tinte rgba a partir de hex ────────────────────────────────── */
function ctTint(hex, alpha) {
    if (alpha == null) alpha = 0.12;
    if (!hex) return 'rgba(75, 85, 99, ' + alpha + ')';
    var h = hex.replace('#', '');
    if (h.length === 3) h = h.split('').map(function(c){ return c + c; }).join('');
    var r = parseInt(h.substr(0, 2), 16);
    var g = parseInt(h.substr(2, 2), 16);
    var b = parseInt(h.substr(4, 2), 16);
    return 'rgba(' + r + ',' + g + ',' + b + ',' + alpha + ')';
}

/* ─── Helper: formatear fecha ─────────────────────────────────────────── */
function ctFormatDate(raw) {
    if (!raw) return '';
    var d = new Date(raw);
    if (isNaN(d.getTime())) return raw;
    return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })
        + ' ' + d.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
}

/* ─── Componente Vue ─────────────────────────────────────────────────────── */
app.component('case-timeline', {
    delimiters: ['${', '}'],
    props: {
        caseId:      { type: String, required: true },
        barrierId:   { type: String, required: false, default: '' },
        showFilters: { type: Boolean, required: false, default: true },
    },
    data() {
        return {
            loading:   true,
            error:     null,
            events:    [],
            filtro:    '',
            categorias: CT_CATEGORIAS,
        };
    },
    computed: {
        eventosFiltrados() {
            if (!this.filtro) return this.events;
            return this.events.filter(function(ev) { return ev.category === this.filtro; }, this);
        },
        countPorCategoria() {
            var counts = {};
            var self = this;
            CT_CATEGORIAS.forEach(function(cat) {
                if (cat.value === '') {
                    counts[''] = self.events.length;
                } else {
                    counts[cat.value] = self.events.filter(function(ev) { return ev.category === cat.value; }).length;
                }
            });
            return counts;
        },
    },
    async mounted() {
        await this.cargar();
    },
    methods: {
        async cargar() {
            this.loading = true;
            this.error   = null;
            try {
                var url = '/api/v1/casos/' + this.caseId + '/timeline-events';
                if (this.barrierId) url += '?barrierId=' + encodeURIComponent(this.barrierId);
                var res = await fetch(url);
                if (!res.ok) throw new Error('Error ' + res.status + ' al cargar el timeline');
                var data = await res.json();
                this.events = Array.isArray(data) ? data : [];
            } catch (e) {
                this.error = e.message;
            } finally {
                this.loading = false;
            }
        },
        iconColor(ev) {
            return ev.color || ctDefaultColor(ev.category);
        },
        badgeStyle(ev) {
            var c = this.iconColor(ev);
            return {
                color: c,
                background: ctTint(c, 0.12),
                border: '1px solid ' + ctTint(c, 0.32),
            };
        },
        actorLabel(ev) {
            return ev.actor_full_name || ev.actor_name || '';
        },
        iconClass(icon) {
            if (!icon) return '';
            return icon.startsWith('fa ') ? icon : 'fa fa-' + icon;
        },
        formatFecha(raw) {
            return ctFormatDate(raw);
        },
    },
    template: `
<div class="ct-wrapper">

    <!-- Cargando -->
    <div v-if="loading" class="ct-loading">
        <svg class="ct-spinner" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
            <path d="M12 2a10 10 0 0 1 10 10" stroke="#3b82f6" stroke-width="3" stroke-linecap="round"/>
        </svg>
        Cargando timeline...
    </div>

    <!-- Error -->
    <div v-else-if="error" class="ct-error">\${ error }</div>

    <!-- Contenido -->
    <template v-else>

        <!-- Filtros por categoría -->
        <div v-if="showFilters" class="ct-filters">
            <button
                v-for="cat in categorias"
                :key="cat.value"
                type="button"
                :class="['ct-filter-btn', filtro === cat.value ? 'active' : '']"
                @click="filtro = cat.value"
            >
                \${ cat.label }
                <span class="ct-filter-count">\${ countPorCategoria[cat.value] || 0 }</span>
            </button>
        </div>

        <!-- Vacío -->
        <div v-if="eventosFiltrados.length === 0" class="ct-empty">
            <div class="ct-empty-icon">📋</div>
            <p>No hay eventos en este timeline</p>
        </div>

        <!-- Lista de eventos -->
        <div v-else class="ct-list">
            <div v-for="ev in eventosFiltrados" :key="ev.id" class="ct-item">

                <!-- Icono -->
                <div class="ct-icon-wrap" :style="{ background: iconColor(ev) }">
                    <i v-if="ev.icon" :class="iconClass(ev.icon)" style="font-size:14px"></i>
                    <span v-else style="font-size:14px">•</span>
                </div>

                <!-- Card -->
                <div class="ct-card">
                    <div class="ct-card-header">
                        <span class="ct-type-badge" :style="badgeStyle(ev)">
                            \${ ev.type || ev.category || '—' }
                        </span>
                        <span class="ct-date">\${ formatFecha(ev.date || ev.created_at) }</span>
                    </div>
                    <div v-if="ev.description" class="ct-desc">\${ ev.description }</div>
                    <div v-if="actorLabel(ev)" class="ct-actor">
                        <i class="fa fa-user" style="font-size:10px"></i>
                        \${ actorLabel(ev) }
                    </div>
                </div>

            </div>
        </div>

    </template>

</div>
    `,
});
