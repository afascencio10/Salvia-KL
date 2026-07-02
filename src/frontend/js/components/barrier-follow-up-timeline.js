(function injectBarrierFollowUpStyles() {
    if (document.getElementById('bft-styles')) return;
    var style = document.createElement('style');
    style.id = 'bft-styles';
    style.textContent = `
        .bft-container { padding: 4px 0 }
        .bft-loading { text-align:center;padding:32px 0;color:#9ca3af;font-size:.85rem }
        .bft-error { background:#fef2f2;border:1px solid #fecaca;border-radius:8px;padding:12px 16px;color:#dc2626;font-size:.82rem }
        .bft-empty { text-align:center;padding:32px 0;color:#9ca3af;font-size:.84rem }
        .bft-list { display:flex;flex-direction:column;gap:12px }
        .bft-card { border:1px solid #e5e7eb;border-radius:10px;padding:14px 16px;background:#fff }
        .bft-card-header { display:flex;align-items:center;justify-content:space-between;margin-bottom:10px }
        .bft-fecha { font-size:.75rem;color:#6b7280;font-weight:600 }
        .bft-autor { font-size:.75rem;color:#374151;font-weight:500 }
        .bft-tag-row { margin-bottom:10px }
        .bft-tag { display:inline-block;padding:2px 10px;border-radius:20px;font-size:.7rem;font-weight:700;letter-spacing:.02em }
        .bft-tag--green { background:#dcfce7;color:#15803d;border:1px solid #bbf7d0 }
        .bft-tag--orange { background:#fff7ed;color:#c2410c;border:1px solid #fed7aa }
        .bft-tag--gray { background:#f3f4f6;color:#6b7280;border:1px solid #e5e7eb }
        .bft-field { margin-bottom:8px }
        .bft-field:last-child { margin-bottom:0 }
        .bft-label { font-size:.68rem;color:#9ca3af;font-weight:700;text-transform:uppercase;letter-spacing:.04em;display:block;margin-bottom:3px }
        .bft-value { font-size:.82rem;color:#374151 }
        .bft-text { font-size:.82rem;color:#374151;line-height:1.5 }
        .bft-chips { display:flex;flex-wrap:wrap;gap:5px;margin-top:2px }
        .bft-chip { background:#eff6ff;color:#1d4ed8;border:1px solid #bfdbfe;border-radius:20px;padding:2px 9px;font-size:.7rem;font-weight:500 }
    `;
    document.head.appendChild(style);
})();

app.component('barrier-follow-up-timeline', {
    delimiters: ['${', '}'],
    props: {
        barrierId: { type: String, required: true },
    },
    data: function() {
        return {
            cargando: false,
            error: null,
            seguimientos: [],
        };
    },
    watch: {
        barrierId: function(newId) {
            if (newId) this.cargar();
        },
    },
    mounted: function() {
        if (this.barrierId) this.cargar();
    },
    methods: {
        async cargar() {
            this.cargando = true;
            this.error = null;
            this.seguimientos = [];
            try {
                var res = await fetch('/api/v1/barriers-v2/' + this.barrierId + '/follow-ups');
                if (!res.ok) throw new Error('Error ' + res.status);
                this.seguimientos = await res.json();
            } catch (e) {
                this.error = 'No se pudieron cargar los seguimientos.';
            } finally {
                this.cargando = false;
            }
        },
        formatFecha: function(iso) {
            if (!iso) return '';
            var d = new Date(iso);
            if (d.getFullYear() < 2000) return '';
            var meses = ['ene','feb','mar','abr','may','jun','jul','ago','sep','oct','nov','dic'];
            return d.getDate() + ' ' + meses[d.getMonth()] + '. ' + d.getFullYear();
        },
        resolveRespuesta: function(val) {
            var map = {
                'respuesta_oficial': 'Respondió de manera oficial a Salvia',
                'contacto_victima': 'Se contactó con la víctima',
                'sin_respuesta': 'No se obtuvo respuesta',
            };
            return map[val] || val;
        },
        resolveGestionCSV: function(csv) {
            if (!csv) return [];
            var map = {
                'activacion_ruta_interinstitucional': 'Activación de ruta interinstitucional',
                'alerta_barreras': 'Alerta por barreras',
                'articulacion_institucional': 'Articulación institucional',
                'escalamiento_organismo_control': 'Escalamiento a organismo de control',
                'gestion_llamada': 'Gestión administrativa - Llamada',
                'orientacion_llamada': 'Orientación y enrutamiento - Llamada',
            };
            return csv.split(',').map(function(v) { var k = v.trim(); return map[k] || k; });
        },
        resolveCierre: function(val) {
            var map = {
                'resuelta': 'Resuelta de manera efectiva',
                'superada_parcialmente': 'Superada parcialmente',
                'instalada_entidad': 'Instalada en entidad competente',
                'no_gestionable': 'No gestionable desde la competencia institucional',
                'no_voluntad': 'Expresa no voluntad de accionar institucional',
                'clasificacion_incorrecta': 'Clasificada incorrectamente',
            };
            return map[val] || val;
        },
        tagClass: function(s) {
            if (s.closesBarrier) return 'bft-tag--green';
            if (s.persists) return 'bft-tag--orange';
            return 'bft-tag--gray';
        },
        tagLabel: function(s) {
            if (s.closesBarrier) return 'Cierra barrera';
            if (s.persists) return 'Persiste';
            return 'No persiste';
        },
    },
    template: `
        <div class="bft-container">
            <div v-if="cargando" class="bft-loading">Cargando seguimientos...</div>
            <div v-else-if="error" class="bft-error">\${ error }</div>
            <template v-else>
                <div v-if="seguimientos.length === 0" class="bft-empty">
                    No hay seguimientos registrados para esta barrera.
                </div>
                <div v-else class="bft-list">
                    <div v-for="s in seguimientos" :key="s.id" class="bft-card">
                        <div class="bft-card-header">
                            <span class="bft-fecha">\${ formatFecha(s.createdAt) }</span>
                            <span class="bft-autor">\${ s.actorName || s.createdById }</span>
                        </div>
                        <div class="bft-tag-row">
                            <span class="bft-tag" :class="tagClass(s)">\${ tagLabel(s) }</span>
                        </div>
                        <div v-if="s.institutionalResponse" class="bft-field">
                            <span class="bft-label">Respuesta institucional</span>
                            <span class="bft-value">\${ resolveRespuesta(s.institutionalResponse) }</span>
                        </div>
                        <div v-if="s.managementActions" class="bft-field">
                            <span class="bft-label">Gestión realizada</span>
                            <div class="bft-chips">
                                <span v-for="chip in resolveGestionCSV(s.managementActions)" :key="chip" class="bft-chip">\${ chip }</span>
                            </div>
                        </div>
                        <div v-if="s.actions" class="bft-field">
                            <span class="bft-label">Actuaciones</span>
                            <p class="bft-text">\${ s.actions }</p>
                        </div>
                        <div v-if="s.closesBarrier && s.closureReason" class="bft-field">
                            <span class="bft-label">Motivo del cierre</span>
                            <span class="bft-value">\${ resolveCierre(s.closureReason) }</span>
                        </div>
                    </div>
                </div>
            </template>
        </div>
    `,
});
