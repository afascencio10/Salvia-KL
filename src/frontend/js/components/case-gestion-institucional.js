(function injectCaseGestionInstitucionalStyles() {
    if (document.getElementById('cgi-styles')) return;
    var style = document.createElement('style');
    style.id = 'cgi-styles';
    style.textContent = `
        .cgi-card { background:linear-gradient(135deg, #faf8ff 0%, #ffffff 60%);border:1px solid #ede9fe;border-radius:14px;padding:18px 20px;margin-bottom:24px;padding-bottom:20px }
        .cgi-header { display:flex;align-items:center;gap:8px;margin-bottom:14px }
        .cgi-header h3 { margin:0;font-size:.86rem;font-weight:700;color:#111827 }
        .cgi-loading { text-align:center;padding:24px 0;color:#9ca3af;font-size:.85rem }
        .cgi-error { background:#fef2f2;border:1px solid #fecaca;border-radius:8px;padding:12px 16px;color:#dc2626;font-size:.82rem;display:flex;align-items:center;justify-content:space-between;gap:12px }
        .cgi-error button { background:#dc2626;color:#fff;border:none;border-radius:6px;padding:5px 12px;font-size:.75rem;font-weight:600;cursor:pointer;white-space:nowrap }

        .cgi-kpis { display:grid;grid-template-columns:repeat(2, 1fr);gap:12px;margin-bottom:20px }
        @media (min-width: 768px) { .cgi-kpis { grid-template-columns:repeat(4, 1fr) } }
        .cgi-tile { border-radius:10px;padding:12px 14px;text-align:left;cursor:pointer;border:1px solid transparent;transition:filter .15s,box-shadow .15s;width:100%;font:inherit }
        .cgi-tile:hover { filter:brightness(0.98);box-shadow:0 1px 4px rgba(0,0,0,.06) }
        .cgi-tile-label { font-size:.72rem;font-weight:600;display:flex;align-items:center;gap:6px;margin:0 }
        .cgi-tile-value { font-size:1.35rem;font-weight:700;margin:4px 0 0 }
        .cgi-tone-violet { background:#f3edfe;border-color:#ede9fe }
        .cgi-tone-violet .cgi-tile-label { color:#7c3aed }
        .cgi-tone-violet .cgi-tile-value { color:#5106A7 }
        .cgi-tone-green { background:#f0fdf4;border-color:#dcfce7 }
        .cgi-tone-green .cgi-tile-label { color:#16a34a }
        .cgi-tone-green .cgi-tile-value { color:#15803d }
        .cgi-tone-amber { background:#fffbeb;border-color:#fef3c7 }
        .cgi-tone-amber .cgi-tile-label { color:#d97706 }
        .cgi-tone-amber .cgi-tile-value { color:#92400e }
        .cgi-tone-red { background:#fef2f2;border-color:#fecaca }
        .cgi-tone-red .cgi-tile-label { color:#dc2626 }
        .cgi-tone-red .cgi-tile-value { color:#b91c1c }
        .cgi-tone-gray { background:#f9fafb;border-color:#f3f4f6 }
        .cgi-tone-gray .cgi-tile-label { color:#6b7280 }
        .cgi-tone-gray .cgi-tile-value { color:#374151;font-size:.95rem }

        .cgi-section-label { font-size:.68rem;font-weight:700;color:#9ca3af;text-transform:uppercase;letter-spacing:.05em;margin:0 0 12px }
        .cgi-sectores { margin-bottom:20px }
        .cgi-sectores-grid { display:grid;grid-template-columns:repeat(3, 1fr);gap:8px }
        .cgi-sector-item { background:none;border:none;cursor:pointer;display:flex;flex-direction:column;align-items:center;text-align:center;gap:2px;font:inherit;padding:4px }
        .cgi-sector-item:hover { opacity:.82 }
        .cgi-sector-nombre { font-size:.78rem;font-weight:600;color:#374151;margin-top:6px }
        .cgi-sector-meta { font-size:.68rem;color:#9ca3af }

        .cgi-estados { display:grid;grid-template-columns:1fr;gap:16px;padding-top:16px;border-top:1px solid #ede9fe }
        @media (min-width: 640px) { .cgi-estados { grid-template-columns:1fr 1fr } }
        .cgi-badges { display:flex;flex-wrap:wrap;gap:8px }
        .cgi-badge { font-size:.75rem;font-weight:600;padding:4px 10px;border-radius:999px;border:1px solid transparent;cursor:pointer;font:inherit }
        .cgi-badge:hover { filter:brightness(0.96) }
        .cgi-badge-barrera-open { background:#fef2f2;color:#b91c1c;border-color:#fecaca }
        .cgi-badge-barrera-gestion { background:#fff7ed;color:#c2410c;border-color:#fed7aa }
        .cgi-badge-barrera-resuelta { background:#f0fdf4;color:#15803d;border-color:#bbf7d0 }
        .cgi-badge-oficio-gestion { background:#fffbeb;color:#92400e;border-color:#fde68a }
        .cgi-badge-oficio-radicado { background:#f0fdfa;color:#0f766e;border-color:#99f6e4 }
        .cgi-badge-oficio-respondido { background:#f0fdf4;color:#15803d;border-color:#bbf7d0 }

        .cgi-modal-overlay { position:fixed;inset:0;background:rgba(0,0,0,.35);display:flex;align-items:center;justify-content:center;z-index:1000;padding:16px }
        .cgi-modal { background:#fff;border-radius:14px;width:100%;max-width:720px;max-height:85vh;overflow:hidden;display:flex;flex-direction:column;box-shadow:0 10px 40px rgba(0,0,0,.2) }
        .cgi-modal-header { display:flex;align-items:flex-start;justify-content:space-between;padding:18px 22px;border-bottom:1px solid #f3f4f6 }
        .cgi-modal-header h3 { margin:0;font-size:.95rem;font-weight:700;color:#111827 }
        .cgi-modal-header p { margin:2px 0 0;font-size:.74rem;color:#9ca3af }
        .cgi-modal-close { background:none;border:none;font-size:1.2rem;color:#9ca3af;cursor:pointer;line-height:1;padding:4px }
        .cgi-modal-body { overflow-y:auto;flex:1 }
        .cgi-modal-empty { text-align:center;padding:40px 0;color:#9ca3af;font-size:.85rem }
        .cgi-modal-table { width:100%;font-size:.82rem;border-collapse:collapse }
        .cgi-modal-table thead { position:sticky;top:0;background:#f9fafb }
        .cgi-modal-table th { text-align:left;padding:10px 18px;font-size:.68rem;font-weight:700;color:#6b7280;text-transform:uppercase;letter-spacing:.04em }
        .cgi-modal-table td { padding:10px 18px;border-top:1px solid #f3f4f6;color:#374151 }
        .cgi-chip-sector { font-size:.7rem;font-weight:700;padding:2px 8px;border-radius:6px;background:#ede9fe;color:#5106A7 }
        .cgi-chip-estado { font-size:.7rem;font-weight:600;padding:2px 8px;border-radius:999px }
        .cgi-modal-footer { display:flex;justify-content:flex-end;padding:14px 22px;border-top:1px solid #f3f4f6;background:#fafafa }
        .cgi-btn-cerrar { background:#fff;border:1px solid #e5e7eb;color:#374151;border-radius:8px;padding:7px 16px;font-size:.82rem;cursor:pointer }
    `;
    document.head.appendChild(style);
})();

// Códigos reales de entity.entity_sector — mismos que ya usa case-entities.js.
var CGI_SECTOR_LABELS = { he: 'Salud', js: 'Justicia', pt: 'Protección' };
var CGI_SECTORES_FIJOS = ['he', 'js', 'pt'];

var CGI_ESTADOS_BARRERA = ['Abierta', 'En gestión', 'Resuelta'];
var CGI_BUCKETS_OFICIO = ['En gestión', 'Radicado', 'Respondido'];

app.component('case-gestion-institucional', {
    delimiters: ['${', '}'],
    props: {
        caseId: { type: String, required: true },
        barriers: { type: Array, required: true },
    },
    data: function() {
        return {
            cargando: false,
            error: null,

            entidades: [],
            resumenObligaciones: { global: { completadas: 0, total: 0 }, porSector: [], items: [] },
            oficios: [],

            sectorLabels: CGI_SECTOR_LABELS,
            modal: null,
        };
    },
    computed: {
        pctGlobal: function() {
            var g = this.resumenObligaciones.global;
            return g.total > 0 ? Math.round((g.completadas / g.total) * 100) : null;
        },
        toneGlobal: function() {
            return this.toneParaPorcentaje(this.pctGlobal);
        },
        porSectorCompleto: function() {
            var self = this;
            var porSector = this.resumenObligaciones.porSector || [];
            return CGI_SECTORES_FIJOS.map(function(sector) {
                var fila = porSector.find(function(s) { return s.sector === sector; });
                var total = fila ? fila.total : 0;
                var completadas = fila ? fila.completadas : 0;
                return {
                    sector: sector,
                    cantidadEntidades: fila ? fila.cantidadEntidades : 0,
                    completadas: completadas,
                    total: total,
                    pct: total > 0 ? Math.round((completadas / total) * 100) : null,
                };
            });
        },
        barrerasPorEstado: function() {
            var self = this;
            return {
                'Abierta': this.barriers.filter(function(b) { return b.status === 'OPEN'; }).length,
                'En gestión': this.barriers.filter(function(b) { return b.status === 'En Gestion'; }).length,
                'Resuelta': this.barriers.filter(function(b) { return self.barreraEsResuelta(b.status); }).length,
            };
        },
        barrerasSinResolver: function() {
            return this.barrerasPorEstado['Abierta'] + this.barrerasPorEstado['En gestión'];
        },
        oficiosPorBucket: function() {
            return {
                'En gestión': this.oficios.filter(function(o) { return o.state !== 'radicado' && o.state !== 'respondido'; }).length,
                'Radicado': this.oficios.filter(function(o) { return o.state === 'radicado'; }).length,
                'Respondido': this.oficios.filter(function(o) { return o.state === 'respondido'; }).length,
            };
        },
        estadosBarrera: function() { return CGI_ESTADOS_BARRERA; },
        bucketsOficio: function() { return CGI_BUCKETS_OFICIO; },
    },
    mounted: function() {
        this.cargar();
        this._onKeydown = this.onKeydown.bind(this);
        window.addEventListener('keydown', this._onKeydown);
    },
    unmounted: function() {
        window.removeEventListener('keydown', this._onKeydown);
    },
    methods: {
        onKeydown: function(e) {
            if (e.key === 'Escape') this.cerrarModal();
        },

        barreraEsResuelta: function(status) {
            return status === 'Articulada' || status === 'MANAGED' || status === 'CLOSED' ||
                status === 'FINISH' || status === 'CLOSE' || status === 'RESOLVED';
        },

        toneParaPorcentaje: function(pct) {
            if (pct == null) return 'gray';
            if (pct >= 80) return 'green';
            if (pct >= 40) return 'amber';
            return 'red';
        },

        badgeToneBarrera: function(estado) {
            if (estado === 'Abierta') return 'cgi-badge-barrera-open';
            if (estado === 'En gestión') return 'cgi-badge-barrera-gestion';
            return 'cgi-badge-barrera-resuelta';
        },
        badgeToneOficio: function(bucket) {
            if (bucket === 'En gestión') return 'cgi-badge-oficio-gestion';
            if (bucket === 'Radicado') return 'cgi-badge-oficio-radicado';
            return 'cgi-badge-oficio-respondido';
        },

        // Anillo SVG — mismo enfoque que el mockup: relleno = severidad, track gris claro.
        // radius=30 fijo en el <svg> del template (viewBox 72x72) — mismo valor usado aquí.
        ringCircumference: function() {
            return 2 * Math.PI * 30;
        },
        ringOffset: function(pct) {
            var c = this.ringCircumference();
            return c - (c * (pct || 0)) / 100;
        },

        formatFecha: function(iso) {
            if (!iso) return '—';
            var d = new Date(iso);
            if (isNaN(d.getTime())) return iso;
            return d.toLocaleDateString('es-CO', { year: 'numeric', month: '2-digit', day: '2-digit' });
        },

        // Las 3 fuentes (entidades, obligaciones-resumen, oficios) se cargan de forma
        // independiente -- un fallo transitorio en una (500, timeout) no debe tumbar
        // las otras dos ni bloquear el panel completo. Antes un solo `!resultados[i].ok`
        // en cualquiera de las 3 marcaba error=true y ocultaba TODO el resumen, aunque
        // las otras 2 hubieran cargado bien y el caso simplemente no tuviera datos.
        // `this.error` queda solo como red de seguridad para un fallo catastrófico
        // fuera de las 3 fuentes (ver catch de más abajo); en operación normal nunca
        // se setea, así que el panel siempre se renderiza con lo que sí haya cargado.
        async cargarFuente(url, parse, onError) {
            try {
                var res = await fetch(url);
                if (!res.ok) throw new Error('http_' + res.status);
                return parse(await res.json());
            } catch (e) {
                console.warn('[case-gestion-institucional] error cargando ' + url + ':', e);
                return onError();
            }
        },
        async cargar() {
            this.cargando = true;
            this.error = null;
            try {
                var self = this;
                var resultados = await Promise.all([
                    this.cargarFuente(
                        '/api/v1/casos/' + this.caseId + '/entidades',
                        function(data) { return Array.isArray(data) ? data : []; },
                        function() { return []; }
                    ),
                    this.cargarFuente(
                        '/api/v1/casos/' + this.caseId + '/obligaciones-resumen',
                        function(data) { return data || { global: { completadas: 0, total: 0 }, porSector: [], items: [] }; },
                        function() { return { global: { completadas: 0, total: 0 }, porSector: [], items: [] }; }
                    ),
                    this.cargarFuente(
                        '/api/v1/entity-letters?caseId=' + encodeURIComponent(this.caseId),
                        function(data) { return Array.isArray(data) ? data : (data.items || []); },
                        function() { return []; }
                    ),
                ]);

                this.entidades = resultados[0];
                this.resumenObligaciones = resultados[1];
                this.oficios = resultados[2];
            } catch (e) {
                console.warn('[case-gestion-institucional] error inesperado cargando resumen:', e);
                this.error = 'No se pudo cargar el resumen de gestión institucional.';
            } finally {
                this.cargando = false;
            }
        },

        abrirModalEntidades() {
            this.modal = { tipo: 'entidades', titulo: 'Entidades agregadas', filas: this.entidades };
        },
        abrirModalObligaciones(sector) {
            var items = this.resumenObligaciones.items || [];
            var filas = sector ? items.filter(function(i) { return i.sector === sector; }) : items;
            var titulo = sector
                ? 'Obligaciones completadas — ' + (this.sectorLabels[sector] || sector)
                : 'Obligaciones completadas';
            this.modal = { tipo: 'obligaciones', titulo: titulo, filas: filas };
        },
        abrirModalBarreras(estado) {
            var self = this;
            var filas;
            if (!estado) {
                filas = this.barriers.filter(function(b) { return b.status === 'OPEN' || b.status === 'En Gestion'; });
            } else if (estado === 'Abierta') {
                filas = this.barriers.filter(function(b) { return b.status === 'OPEN'; });
            } else if (estado === 'En gestión') {
                filas = this.barriers.filter(function(b) { return b.status === 'En Gestion'; });
            } else {
                filas = this.barriers.filter(function(b) { return self.barreraEsResuelta(b.status); });
            }
            this.modal = { tipo: 'barreras', titulo: estado ? ('Barreras — ' + estado) : 'Barreras sin resolver', filas: filas };
        },
        abrirModalOficios(bucket) {
            var filas;
            if (bucket === 'Radicado') {
                filas = this.oficios.filter(function(o) { return o.state === 'radicado'; });
            } else if (bucket === 'Respondido') {
                filas = this.oficios.filter(function(o) { return o.state === 'respondido'; });
            } else {
                filas = this.oficios.filter(function(o) { return o.state !== 'radicado' && o.state !== 'respondido'; });
            }
            this.modal = { tipo: 'oficios', titulo: 'Oficios — ' + bucket, filas: filas };
        },
        cerrarModal() {
            this.modal = null;
        },
    },
    template: `
        <div class="cgi-card">
            <div class="cgi-header">
                <span>📊</span>
                <h3>Resumen de Gestión Institucional</h3>
            </div>

            <div v-if="cargando" class="cgi-loading">Cargando resumen...</div>

            <div v-else-if="error" class="cgi-error">
                <span>\${ error }</span>
                <button @click="cargar">Reintentar</button>
            </div>

            <div v-else>
                <div class="cgi-kpis">
                    <button type="button" class="cgi-tile cgi-tone-violet" @click="abrirModalEntidades">
                        <p class="cgi-tile-label"><span>🏛️</span>Entidades agregadas</p>
                        <p class="cgi-tile-value">\${ entidades.length }</p>
                    </button>

                    <button type="button" class="cgi-tile" :class="'cgi-tone-' + toneGlobal" @click="abrirModalObligaciones(null)">
                        <p class="cgi-tile-label"><span>📋</span>Obligaciones completadas</p>
                        <p class="cgi-tile-value">
                            <template v-if="pctGlobal != null">\${ resumenObligaciones.global.completadas }/\${ resumenObligaciones.global.total } · \${ pctGlobal }%</template>
                            <template v-else>Sin obligaciones</template>
                        </p>
                    </button>

                    <button type="button" class="cgi-tile" :class="barrerasSinResolver > 0 ? 'cgi-tone-red' : 'cgi-tone-green'" @click="abrirModalBarreras(null)">
                        <p class="cgi-tile-label"><span>⚠️</span>Barreras sin resolver</p>
                        <p class="cgi-tile-value">\${ barrerasSinResolver }</p>
                    </button>

                    <button type="button" class="cgi-tile" :class="oficiosPorBucket['En gestión'] > 0 ? 'cgi-tone-amber' : 'cgi-tone-green'" @click="abrirModalOficios('En gestión')">
                        <p class="cgi-tile-label"><span>📄</span>Oficios en gestión</p>
                        <p class="cgi-tile-value">\${ oficiosPorBucket['En gestión'] }</p>
                    </button>
                </div>

                <div class="cgi-sectores">
                    <p class="cgi-section-label">% Obligaciones completadas por sector</p>
                    <div class="cgi-sectores-grid">
                        <button type="button" class="cgi-sector-item" v-for="s in porSectorCompleto" :key="s.sector" @click="abrirModalObligaciones(s.sector)">
                            <svg width="72" height="72" viewBox="0 0 72 72" style="transform:rotate(-90deg)">
                                <circle cx="36" cy="36" r="30" fill="none" stroke-width="7" style="stroke:#f3f4f6" />
                                <circle v-if="s.pct != null" cx="36" cy="36" r="30" fill="none" stroke-width="7" stroke-linecap="round"
                                    :style="{ stroke: s.pct >= 80 ? '#22c55e' : (s.pct >= 40 ? '#f59e0b' : '#ef4444'), strokeDasharray: ringCircumference(), strokeDashoffset: ringOffset(s.pct) }" />
                            </svg>
                            <span class="cgi-sector-nombre">\${ sectorLabels[s.sector] }</span>
                            <span class="cgi-sector-meta">\${ s.cantidadEntidades } entidad\${ s.cantidadEntidades === 1 ? '' : 'es' }</span>
                            <span class="cgi-sector-meta">\${ s.total > 0 ? (s.completadas + '/' + s.total + ' obligaciones') : 'Sin obligaciones' }</span>
                        </button>
                    </div>
                </div>

                <div class="cgi-estados">
                    <div>
                        <p class="cgi-section-label">Barreras por estado</p>
                        <div class="cgi-badges">
                            <button type="button" class="cgi-badge" :class="badgeToneBarrera(estado)" v-for="estado in estadosBarrera" :key="estado" @click="abrirModalBarreras(estado)">
                                \${ estado }: \${ barrerasPorEstado[estado] }
                            </button>
                        </div>
                    </div>
                    <div>
                        <p class="cgi-section-label">Oficios por estado</p>
                        <div class="cgi-badges">
                            <button type="button" class="cgi-badge" :class="badgeToneOficio(bucket)" v-for="bucket in bucketsOficio" :key="bucket" @click="abrirModalOficios(bucket)">
                                \${ bucket }: \${ oficiosPorBucket[bucket] }
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="modal" class="cgi-modal-overlay" @click.self="cerrarModal">
                <div class="cgi-modal">
                    <div class="cgi-modal-header">
                        <div>
                            <h3>\${ modal.titulo }</h3>
                            <p>\${ modal.filas.length } resultado\${ modal.filas.length === 1 ? '' : 's' }</p>
                        </div>
                        <button class="cgi-modal-close" @click="cerrarModal">×</button>
                    </div>
                    <div class="cgi-modal-body">
                        <div v-if="modal.filas.length === 0" class="cgi-modal-empty">No hay registros para este indicador</div>
                        <table v-else class="cgi-modal-table">
                            <thead>
                                <tr v-if="modal.tipo === 'entidades'">
                                    <th>Sector</th><th>Entidad</th><th>Ubicación</th><th>Oficios</th><th>Barreras activas</th>
                                </tr>
                                <tr v-else-if="modal.tipo === 'obligaciones'">
                                    <th>Entidad</th><th>Sector</th><th>Obligación</th><th>Estado</th>
                                </tr>
                                <tr v-else-if="modal.tipo === 'barreras'">
                                    <th>Sector</th><th>Estado</th><th>Descripción</th><th>Fecha</th>
                                </tr>
                                <tr v-else>
                                    <th>Entidad</th><th>Estado</th><th>Fecha</th>
                                </tr>
                            </thead>
                            <tbody>
                                <template v-if="modal.tipo === 'entidades'">
                                    <tr v-for="e in modal.filas" :key="e.relId">
                                        <td><span class="cgi-chip-sector">\${ sectorLabels[e.sector] || e.sector }</span></td>
                                        <td>\${ e.entityBranchName }</td>
                                        <td>\${ e.townName }, \${ e.cityName }</td>
                                        <td>\${ e.oficiosCount }</td>
                                        <td>\${ e.barrerasActivasCount }</td>
                                    </tr>
                                </template>
                                <template v-else-if="modal.tipo === 'obligaciones'">
                                    <tr v-for="(o, i) in modal.filas" :key="i">
                                        <td>\${ o.entidad }</td>
                                        <td><span class="cgi-chip-sector">\${ sectorLabels[o.sector] || o.sector }</span></td>
                                        <td>\${ o.obligacion }</td>
                                        <td>
                                            <span class="cgi-chip-estado" :style="o.completada ? 'background:#dcfce7;color:#15803d' : 'background:#f3f4f6;color:#6b7280'">
                                                \${ o.completada ? '✅ Completada' : '⭕ Pendiente' }
                                            </span>
                                        </td>
                                    </tr>
                                </template>
                                <template v-else-if="modal.tipo === 'barreras'">
                                    <tr v-for="b in modal.filas" :key="b.id">
                                        <td><span class="cgi-chip-sector">\${ b.sector }</span></td>
                                        <td>\${ b.status }</td>
                                        <td>\${ b.description || b.specificBarriers || '—' }</td>
                                        <td>\${ formatFecha(b.createdAt) }</td>
                                    </tr>
                                </template>
                                <template v-else>
                                    <tr v-for="o in modal.filas" :key="o.id">
                                        <td>\${ o.entidad || '—' }</td>
                                        <td>\${ o.state }</td>
                                        <td>\${ formatFecha(o.createdAt) }</td>
                                    </tr>
                                </template>
                            </tbody>
                        </table>
                    </div>
                    <div class="cgi-modal-footer">
                        <button class="cgi-btn-cerrar" @click="cerrarModal">Cerrar</button>
                    </div>
                </div>
            </div>
        </div>
    `,
});
