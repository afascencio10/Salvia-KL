/**
 * case-professionals.js
 * Sección "Profesionales de atención" de la pestaña Derivaciones del detalle
 * del caso. Muestra una tabla con Usuario, Rol y Notas: una fila por persona
 * que ha intervenido en el caso (desde cualquier módulo), con sus roles
 * concatenados y la consolidación cronológica de sus intervenciones.
 *
 * Consume GET /api/v1/casos/:id/profesionales-atencion.
 */
(function injectCaseProfessionalsStyles() {
    if (document.getElementById('cp-styles')) return;
    var style = document.createElement('style');
    style.id = 'cp-styles';
    style.textContent = `
        .cp-table-wrap { border:1px solid #e5e7eb;border-radius:10px;overflow:hidden;background:#fff }
        .cp-table { width:100%;border-collapse:collapse;font-size:.82rem }
        .cp-table thead tr { background:#f8fafc;color:#334155;text-transform:uppercase;font-size:.68rem;letter-spacing:.04em }
        .cp-table th { text-align:left;padding:9px 14px;font-weight:700 }
        .cp-table td { padding:10px 14px;color:#374151;border-top:1px solid #f1f5f9;vertical-align:top }
        .cp-usuario { font-weight:600;color:#111827;white-space:nowrap }
        .cp-rol { color:#5106A7;font-weight:600;white-space:nowrap }
        .cp-notas { color:#4b5563;line-height:1.45 }
        .cp-loading, .cp-empty { padding:16px;font-size:.82rem;color:#9ca3af;font-style:italic }
        .cp-error { background:#fef2f2;border:1px solid #fecaca;border-radius:8px;padding:12px 16px;color:#dc2626;font-size:.82rem }
    `;
    document.head.appendChild(style);
})();

app.component('case-professionals', {
    delimiters: ['${', '}'],
    props: {
        caseId: { type: String, required: true },
    },
    data: function() {
        return {
            cargando: false,
            error: null,
            profesionales: [],
        };
    },
    mounted: function() {
        this.cargar();
    },
    methods: {
        cargar: function() {
            var self = this;
            self.cargando = true;
            self.error = null;
            fetch('/api/v1/casos/' + self.caseId + '/profesionales-atencion')
                .then(function(res) {
                    if (!res.ok) throw new Error('Error ' + res.status);
                    return res.json();
                })
                .then(function(data) {
                    self.profesionales = Array.isArray(data) ? data : [];
                    self.cargando = false;
                })
                .catch(function() {
                    self.error = 'No se pudieron cargar los profesionales de atención de este caso.';
                    self.cargando = false;
                });
        },
    },
    template: `
        <div class="cd-deriv-section" style="background:#faf5ff;border-color:#d8b4fe">
            <div class="cd-deriv-header" style="border-bottom-color:#d8b4fe">
                👥 Profesionales de atención
                <span class="cd-count" style="color:#5106A7">\${ profesionales.length } persona\${ profesionales.length !== 1 ? 's' : '' }</span>
            </div>

            <div v-if="cargando" class="cp-loading">Cargando profesionales de atención...</div>
            <div v-else-if="error" class="cp-error">\${ error }</div>
            <div v-else-if="profesionales.length === 0" class="cd-deriv-empty">Sin intervenciones registradas</div>

            <div v-else style="padding:12px 16px">
                <div class="cp-table-wrap">
                    <table class="cp-table">
                        <thead>
                            <tr>
                                <th>Usuario</th>
                                <th>Rol</th>
                                <th>Notas</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, idx) in profesionales" :key="idx">
                                <td class="cp-usuario">\${ p.usuario || '—' }</td>
                                <td class="cp-rol">\${ p.rol || '—' }</td>
                                <td class="cp-notas">\${ p.notas || '—' }</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    `,
});
