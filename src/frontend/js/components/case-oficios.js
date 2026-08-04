/**
 * case-oficios
 * Lista los oficios (EntityLetter) de un caso o de una barrera en formato
 * de cards, con buscador, chips de tema y filtro por estado. Click en una
 * card abre un modal de solo lectura con el detalle completo del oficio.
 *
 * Props:
 *   caseId    (String, opcional) — icode del caso. Requerido si no se pasa barrierId.
 *   barrierId (String, opcional) — id de la barrera. Si se pasa, tiene prioridad
 *                                   sobre caseId y los chips de tema se ocultan.
 *
 * Nota: componente nuevo y separado de oficios-list.js (que sigue usándose
 * tal cual en la pantalla de Notificaciones, en formato tabla).
 */

/* ─── Estilos ────────────────────────────────────────────────────────────── */
(function injectCaseOficiosStyles() {
    if (document.getElementById('case-oficios-styles')) return;
    const style = document.createElement('style');
    style.id = 'case-oficios-styles';
    style.textContent = `
        .co-container { padding: 4px 0; }

        /* Filtros */
        .co-filtros {
            display: flex;
            flex-direction: column;
            gap: 10px;
            margin-bottom: 16px;
        }
        .co-buscador {
            width: 100%;
            padding: 9px 12px;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            font-size: 13px;
            outline: none;
            box-sizing: border-box;
        }
        .co-buscador:focus { border-color: #5106A7; box-shadow: 0 0 0 3px rgba(81,6,167,.12); }

        .co-filtros-row {
            display: flex;
            align-items: center;
            gap: 10px;
            flex-wrap: wrap;
        }
        .co-chips-tema { display: flex; gap: 8px; flex-wrap: wrap; }
        .co-chip {
            padding: 5px 12px;
            border-radius: 999px;
            border: 1.5px solid #e5e7eb;
            background: #fff;
            color: #6b7280;
            font-size: 12px;
            font-weight: 600;
            cursor: pointer;
            transition: border-color .15s, background .15s, color .15s;
        }
        .co-chip:hover:not(:disabled) { border-color: #c4b5fd; }
        .co-chip.co-chip-activo {
            background: #ede9fe;
            border-color: #5106A7;
            color: #5106A7;
        }
        .co-chip:disabled { opacity: .45; cursor: not-allowed; }

        .co-select-estado {
            padding: 7px 10px;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            font-size: 12px;
            color: #374151;
            background: #fff;
            outline: none;
        }

        /* Grid de cards */
        .co-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
            gap: 14px;
        }
        .co-card {
            border: 1px solid #e5e7eb;
            border-radius: 10px;
            padding: 14px 16px;
            cursor: pointer;
            transition: border-color .15s, box-shadow .15s;
            background: #fff;
        }
        .co-card:hover { border-color: #c4b5fd; box-shadow: 0 2px 8px rgba(81,6,167,.08); }
        .co-card-header {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 8px;
            margin-bottom: 8px;
        }
        .co-card-entidad { font-size: 13px; font-weight: 700; color: #111827; display: flex; align-items: center; gap: 6px; }
        .co-card-icon { color: #5106A7; font-size: 12px; }
        .co-card-badges { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 8px; }
        .co-card-asunto { font-size: 12px; color: #4b5563; margin-bottom: 8px; line-height: 1.4; }
        .co-card-doc { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #5106A7; margin-bottom: 6px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .co-card-fecha { font-size: 11px; color: #9ca3af; }

        .co-badge {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 999px;
            font-size: 10.5px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: .02em;
        }
        .co-badge--por_proyectar   { background: #fef3c7; color: #92400e; }
        .co-badge--para_revisar    { background: #dbeafe; color: #1d4ed8; }
        .co-badge--en_correccion   { background: #fee2e2; color: #b91c1c; }
        .co-badge--aprobacion_juridica { background: #ede9fe; color: #5106A7; }
        .co-badge--para_radicar    { background: #cffafe; color: #0e7490; }
        .co-badge--radicado        { background: #dcfce7; color: #15803d; }
        .co-badge--respondido      { background: #f0fdf4; color: #166534; }
        .co-badge--default         { background: #f3f4f6; color: #374151; }
        .co-badge--nivel { background: #f3f4f6; color: #4b5563; }
        .co-badge--priority-alta   { background: #fee2e2; color: #b91c1c; }
        .co-badge--priority-normal { background: #f3f4f6; color: #6b7280; }

        /* Estados */
        .co-empty, .co-loading, .co-error {
            text-align: center;
            padding: 40px 0;
            color: #9ca3af;
            font-size: 13px;
        }
        .co-error { color: #b91c1c; }
        .co-btn-retry {
            margin-top: 10px;
            padding: 7px 16px;
            background: #5106A7;
            color: #fff;
            border: none;
            border-radius: 8px;
            font-size: 12px;
            font-weight: 600;
            cursor: pointer;
        }
        @keyframes co-spin { to { transform: rotate(360deg); } }
        .co-spinner {
            width: 26px; height: 26px;
            margin: 0 auto 10px;
            border: 3px solid #e5e7eb;
            border-top-color: #5106A7;
            border-radius: 50%;
            animation: co-spin .7s linear infinite;
        }

        /* Modal */
        .co-modal-overlay {
            position: fixed;
            inset: 0;
            background: rgba(0,0,0,.45);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 1050;
            padding: 16px;
        }
        .co-modal {
            background: #fff;
            border-radius: 16px;
            width: 100%;
            max-width: 560px;
            max-height: 90vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            box-shadow: 0 20px 40px rgba(0,0,0,.18);
        }
        .co-modal-header {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            padding: 18px 22px 14px;
            border-bottom: 1px solid #e5e7eb;
        }
        .co-modal-title { font-size: 15px; font-weight: 700; color: #111827; display: block; }
        .co-modal-subtitle { font-size: 12px; color: #9ca3af; display: block; margin-top: 2px; }
        .co-modal-close { background: none; border: none; font-size: 20px; color: #9ca3af; cursor: pointer; line-height: 1; }
        .co-modal-close:hover { color: #4b5563; }
        .co-modal-body { flex: 1; overflow-y: auto; padding: 18px 22px; }

        .co-msection { margin-bottom: 18px; }
        .co-msection:last-child { margin-bottom: 0; }
        .co-msection-title {
            font-size: 11px;
            font-weight: 700;
            color: #9ca3af;
            text-transform: uppercase;
            letter-spacing: .04em;
            margin-bottom: 10px;
        }
        .co-msection--warning { background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 12px 14px; }
        .co-msection--warning .co-msection-title { color: #b91c1c; }
        .co-msection--relation { background: #f9fafb; border-radius: 8px; padding: 12px 14px; }

        .co-mfields { display: grid; grid-template-columns: 1fr 1fr; gap: 10px 16px; }
        .co-mfield--full { grid-column: 1 / -1; }
        .co-mfield-label { display: block; font-size: 10.5px; font-weight: 600; color: #9ca3af; text-transform: uppercase; margin-bottom: 2px; }
        .co-mfield-value { font-size: 13px; color: #111827; word-break: break-word; }
        .co-mfield-mono { font-family: monospace; font-size: 11.5px; }
        .co-link { color: #5106A7; text-decoration: none; }
        .co-link:hover { text-decoration: underline; }

        .co-modal-actions {
            display: flex;
            gap: 10px;
            margin-top: 20px;
            padding-top: 16px;
            border-top: 1px solid #e5e7eb;
        }
        .co-btn-action {
            flex: 1;
            padding: 9px 14px;
            border-radius: 8px;
            font-size: 12.5px;
            font-weight: 600;
            cursor: pointer;
            border: 1px solid #e5e7eb;
            background: #fff;
            color: #374151;
        }
        .co-btn-action:hover { background: #f9fafb; }
        .co-btn-action--close { background: #5106A7; color: #fff; border-color: #5106A7; }
        .co-btn-action--close:hover { background: #3d0480; }
    `;
    document.head.appendChild(style);
})();

/* ─── Constantes de presentación ─────────────────────────────────────────── */
const CO_STATUS_LABELS = {
    por_proyectar:       'Por proyectar',
    para_revisar:        'Para revisar',
    en_correccion:       'En corrección',
    aprobacion_juridica: 'Aprobación jurídica',
    para_radicar:        'Para radicar',
    radicado:            'Radicado',
    respondido:          'Respondido',
};
const CO_NIVEL_LABELS = {
    nacional:      'Nacional',
    departamental: 'Departamental',
    municipal:     'Municipal',
};
const CO_PRIORITY_LABELS = {
    normal: 'Normal',
    alta:   'Alta',
};

/* ─── Componente Vue ──────────────────────────────────────────────────────── */
app.component('case-oficios', {
    delimiters: ['${', '}'],

    props: {
        caseId:    { type: String, default: null },
        barrierId: { type: String, default: null },
    },

    data() {
        return {
            cargando:  false,
            error:     null,
            oficios:   [],

            buscador:            '',
            temaActivo:          null,
            estadoActivo:        null,
            oficioSeleccionado:  null,
        };
    },

    computed: {
        oficiosFiltrados() {
            const q = this.buscador.trim().toLowerCase();
            return this.oficios.filter(o => {
                if (q) {
                    const asunto = o.subject || o.asuntoRadicado || o.asuntoRespuesta || '';
                    const matches = (o.entidad || '').toLowerCase().includes(q)
                        || (o.urlKofax || '').toLowerCase().includes(q)
                        || asunto.toLowerCase().includes(q);
                    if (!matches) return false;
                }
                if (this.temaActivo === 'barrera' && !o.barrierId) return false;
                if (this.estadoActivo && o.state !== this.estadoActivo) return false;
                return true;
            });
        },
    },

    mounted() {
        this.cargar();
    },

    methods: {
        /* ── Carga de datos ──────────────────────────────────────────────── */

        async cargar() {
            this.cargando = true;
            this.error    = null;
            this.oficios  = [];
            try {
                const query = this.barrierId
                    ? 'barrierId=' + encodeURIComponent(this.barrierId)
                    : 'caseId=' + encodeURIComponent(this.caseId);
                const res = await fetch('/api/v1/entity-letters?' + query);
                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }
                if (!res.ok) throw new Error('Error ' + res.status);
                const data = await res.json();
                this.oficios = Array.isArray(data) ? data : (data.items || []);
            } catch (e) {
                this.error = 'No se pudieron cargar los oficios de ' + (this.barrierId ? 'esta barrera.' : 'este caso.');
            } finally {
                this.cargando = false;
            }
        },

        /* ── Filtros ──────────────────────────────────────────────────────── */

        toggleTema(tema) {
            if (tema !== 'barrera') return; // otros temas deshabilitados — GAP de schema
            this.temaActivo = this.temaActivo === tema ? null : tema;
        },

        /* ── Modal ────────────────────────────────────────────────────────── */

        abrirDetalle(oficio) {
            this.oficioSeleccionado = oficio;
        },

        cerrarDetalle() {
            this.oficioSeleccionado = null;
        },

        /* ── Presentación ─────────────────────────────────────────────────── */

        asuntoResuelto(oficio) {
            return oficio.subject || oficio.asuntoRadicado || oficio.asuntoRespuesta || '';
        },

        statusLabel(state)     { return CO_STATUS_LABELS[state] || state; },
        nivelLabel(nivel)      { return CO_NIVEL_LABELS[nivel] || nivel || '—'; },
        priorityLabel(pri)     { return CO_PRIORITY_LABELS[pri] || pri || '—'; },
        stateClass(state)      { return 'co-badge--' + (state || 'default'); },
        priorityClass(pri)     { return 'co-badge--priority-' + (pri || 'normal'); },

        truncateUrl(url, maxLen) {
            maxLen = maxLen || 40;
            if (!url) return '—';
            return url.length > maxLen ? url.slice(0, maxLen) + '...' : url;
        },

        formatDate(isoStr) {
            if (!isoStr) return '—';
            try {
                const d = new Date(isoStr);
                return d.toLocaleDateString('es-CO', { day: '2-digit', month: '2-digit', year: 'numeric' })
                    + ' ' + d.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
            } catch (e) {
                return isoStr;
            }
        },
    },

    template: `
<div class="co-container">

    <!-- Cargando -->
    <div v-if="cargando" class="co-loading">
        <div class="co-spinner"></div>
        Cargando oficios...
    </div>

    <!-- Error -->
    <div v-else-if="error" class="co-error">
        <i class="fas fa-exclamation-triangle"></i> \${ error }
        <div><button class="co-btn-retry" @click="cargar">Reintentar</button></div>
    </div>

    <template v-else>

        <!-- Filtros -->
        <div class="co-filtros">
            <input
                class="co-buscador"
                v-model="buscador"
                placeholder="Buscar por entidad, ruta o asunto..."
            />
            <div class="co-filtros-row">
                <div v-if="!barrierId" class="co-chips-tema">
                    <button
                        class="co-chip"
                        :class="{ 'co-chip-activo': temaActivo === 'barrera' }"
                        @click="toggleTema('barrera')">
                        Barreras
                    </button>
                    <button class="co-chip" disabled title="Requiere agregar emergency_measure_id a entity_letter">
                        Medidas de Emergencia
                    </button>
                    <button class="co-chip" disabled title="Requiere agregar psychosocial_support_id a entity_letter">
                        Apoyo Psicosocial
                    </button>
                    <button class="co-chip" disabled title="Requiere agregar economic_stabilization_id a entity_letter">
                        Estabilización Económica
                    </button>
                </div>
                <select class="co-select-estado" v-model="estadoActivo">
                    <option :value="null">Todos los estados</option>
                    <option value="por_proyectar">Por proyectar</option>
                    <option value="para_revisar">Para revisar</option>
                    <option value="en_correccion">En corrección</option>
                    <option value="aprobacion_juridica">Aprobación jurídica</option>
                    <option value="para_radicar">Para radicar</option>
                    <option value="radicado">Radicado</option>
                    <option value="respondido">Respondido</option>
                </select>
            </div>
        </div>

        <!-- Empty state -->
        <div v-if="oficiosFiltrados.length === 0" class="co-empty">
            <i class="fas fa-inbox"></i><br>
            No hay oficios que coincidan con los filtros.
        </div>

        <!-- Grid de cards -->
        <div v-else class="co-grid">
            <div
                v-for="oficio in oficiosFiltrados"
                :key="oficio.id"
                class="co-card"
                @click="abrirDetalle(oficio)">

                <div class="co-card-header">
                    <span class="co-card-entidad"><i class="fas fa-file-alt co-card-icon"></i> \${ oficio.entidad || '—' }</span>
                    <span class="co-badge" :class="stateClass(oficio.state)">\${ statusLabel(oficio.state) }</span>
                </div>

                <div class="co-card-badges">
                    <span v-if="oficio.nivel" class="co-badge co-badge--nivel">\${ nivelLabel(oficio.nivel) }</span>
                    <span class="co-badge" :class="priorityClass(oficio.priority)">\${ priorityLabel(oficio.priority) }</span>
                </div>

                <div v-if="asuntoResuelto(oficio)" class="co-card-asunto">\${ asuntoResuelto(oficio) }</div>

                <div v-if="oficio.urlKofax" class="co-card-doc">
                    <i class="fas fa-file-pdf"></i> \${ truncateUrl(oficio.urlKofax) }
                </div>

                <div class="co-card-fecha">\${ formatDate(oficio.createdAt) }</div>
            </div>
        </div>

    </template>

    <!-- ════════════════════════════════════════
         Modal de detalle
         ════════════════════════════════════════ -->
    <div v-if="oficioSeleccionado" class="co-modal-overlay" @click.self="cerrarDetalle">
        <div class="co-modal">

            <div class="co-modal-header">
                <div>
                    <span class="co-modal-title">Detalle del Oficio</span>
                    <span v-if="oficioSeleccionado.urlKofax" class="co-modal-subtitle">\${ truncateUrl(oficioSeleccionado.urlKofax, 55) }</span>
                </div>
                <button class="co-modal-close" @click="cerrarDetalle">×</button>
            </div>

            <div class="co-modal-body">

                <!-- Oficio -->
                <div class="co-msection">
                    <div class="co-msection-title"><i class="fas fa-file-alt"></i> Oficio</div>
                    <div class="co-mfields">
                        <div class="co-mfield">
                            <span class="co-mfield-label">Estado</span>
                            <span class="co-mfield-value">
                                <span class="co-badge" :class="stateClass(oficioSeleccionado.state)">\${ statusLabel(oficioSeleccionado.state) }</span>
                            </span>
                        </div>
                        <div class="co-mfield">
                            <span class="co-mfield-label">Prioridad</span>
                            <span class="co-mfield-value">
                                <span class="co-badge" :class="priorityClass(oficioSeleccionado.priority)">\${ priorityLabel(oficioSeleccionado.priority) }</span>
                            </span>
                        </div>
                        <div class="co-mfield">
                            <span class="co-mfield-label">Entidad</span>
                            <span class="co-mfield-value">
                                \${ oficioSeleccionado.entidad || '—' }
                                <span v-if="oficioSeleccionado.nivel" class="co-badge co-badge--nivel">\${ nivelLabel(oficioSeleccionado.nivel) }</span>
                            </span>
                        </div>
                        <div v-if="oficioSeleccionado.correoEntidad" class="co-mfield">
                            <span class="co-mfield-label">Correo entidad</span>
                            <span class="co-mfield-value"><a :href="'mailto:' + oficioSeleccionado.correoEntidad" class="co-link">\${ oficioSeleccionado.correoEntidad }</a></span>
                        </div>
                        <div v-if="oficioSeleccionado.urlKofax" class="co-mfield co-mfield--full">
                            <span class="co-mfield-label">Documento</span>
                            <span class="co-mfield-value">
                                <i class="fas fa-file-pdf" style="color:#5106A7;margin-right:4px"></i>
                                <a :href="oficioSeleccionado.urlKofax" target="_blank" rel="noopener" class="co-link">\${ oficioSeleccionado.urlKofax }</a>
                            </span>
                        </div>
                        <div class="co-mfield">
                            <span class="co-mfield-label">Creado</span>
                            <span class="co-mfield-value">\${ formatDate(oficioSeleccionado.createdAt) }</span>
                        </div>
                        <div class="co-mfield">
                            <span class="co-mfield-label">Actualizado</span>
                            <span class="co-mfield-value">\${ formatDate(oficioSeleccionado.updatedAt) }</span>
                        </div>
                    </div>
                </div>

                <!-- Radicación y respuesta -->
                <div v-if="oficioSeleccionado.numeroRadicado || oficioSeleccionado.asuntoRadicado || oficioSeleccionado.responseDate" class="co-msection">
                    <div class="co-msection-title"><i class="fas fa-stamp"></i> Radicación y Respuesta</div>
                    <div class="co-mfields">
                        <div v-if="oficioSeleccionado.numeroRadicado" class="co-mfield">
                            <span class="co-mfield-label">Número radicado</span>
                            <span class="co-mfield-value">\${ oficioSeleccionado.numeroRadicado }</span>
                        </div>
                        <div v-if="oficioSeleccionado.asuntoRadicado" class="co-mfield co-mfield--full">
                            <span class="co-mfield-label">Asunto radicado</span>
                            <span class="co-mfield-value">\${ oficioSeleccionado.asuntoRadicado }</span>
                        </div>
                        <div v-if="oficioSeleccionado.correoRemitente" class="co-mfield">
                            <span class="co-mfield-label">Correo remitente</span>
                            <span class="co-mfield-value"><a :href="'mailto:' + oficioSeleccionado.correoRemitente" class="co-link">\${ oficioSeleccionado.correoRemitente }</a></span>
                        </div>
                        <div v-if="oficioSeleccionado.asuntoRespuesta" class="co-mfield co-mfield--full">
                            <span class="co-mfield-label">Asunto respuesta</span>
                            <span class="co-mfield-value">\${ oficioSeleccionado.asuntoRespuesta }</span>
                        </div>
                        <div v-if="oficioSeleccionado.responseDate" class="co-mfield">
                            <span class="co-mfield-label">Fecha respuesta</span>
                            <span class="co-mfield-value">\${ formatDate(oficioSeleccionado.responseDate) }</span>
                        </div>
                    </div>
                </div>

                <!-- Corrección -->
                <div v-if="oficioSeleccionado.reasonCorrection" class="co-msection co-msection--warning">
                    <div class="co-msection-title"><i class="fas fa-exclamation-triangle"></i> Corrección requerida</div>
                    <div class="co-mfield-value">\${ oficioSeleccionado.reasonCorrection }</div>
                </div>

                <!-- Barrera relacionada -->
                <div v-if="oficioSeleccionado.barrierId" class="co-msection co-msection--relation">
                    <div class="co-msection-title"><i class="fas fa-shield-alt"></i> Barrera relacionada</div>
                    <a :href="'/salvia/barreras/' + oficioSeleccionado.barrierId" class="co-link">Ver barrera →</a>
                </div>
                <div class="co-modal-actions">
                    <button class="co-btn-action co-btn-action--close" @click="cerrarDetalle">Cerrar</button>
                </div>

            </div>
        </div>
    </div>

</div>
    `,
});
