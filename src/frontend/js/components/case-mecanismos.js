(function injectCaseMecanismosStyles() {
    if (document.getElementById('cm-styles')) return;
    var style = document.createElement('style');
    style.id = 'cm-styles';
    style.textContent = `
        .cm-card { background:#fff;border:1px solid #e5e7eb;border-radius:12px;padding:16px 18px }
        .cm-header { display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:14px }
        .cm-title { display:flex;align-items:center;gap:8px;font-size:1rem;font-weight:700;color:#111827;margin:0 }
        .cm-title i { color:#5106A7 }
        .cm-btn-agregar { background:#5106A7;color:#fff;border:none;border-radius:8px;padding:8px 14px;font-size:.8rem;font-weight:600;cursor:pointer;white-space:nowrap }
        .cm-btn-agregar:hover { background:#3d0780 }
        .cm-loading, .cm-empty { text-align:center;padding:36px 12px;color:#9ca3af }
        .cm-empty-icon { font-size:2.4rem;color:#d1d5db;margin-bottom:8px }
        .cm-empty p { margin:0;font-size:.88rem;color:#6b7280 }
        .cm-error { background:#fef2f2;border:1px solid #fecaca;border-radius:8px;padding:12px 16px;color:#dc2626;font-size:.82rem;display:flex;align-items:center;justify-content:space-between;gap:12px }
        .cm-error button { background:#dc2626;color:#fff;border:none;border-radius:6px;padding:5px 12px;font-size:.75rem;font-weight:600;cursor:pointer }
        .cm-list { display:flex;flex-direction:column;gap:12px }
        .cm-item { border:1px solid #e5e7eb;border-radius:10px;padding:14px 16px;background:#fff }
        .cm-item-head { display:flex;align-items:center;gap:8px;flex-wrap:wrap;margin-bottom:6px }
        .cm-badge { font-size:.68rem;font-weight:700;padding:3px 10px;border-radius:999px;text-transform:capitalize;letter-spacing:.02em }
        .cm-badge-municipal { background:#e0f2fe;color:#0369a1 }
        .cm-badge-departamental { background:#ede9fe;color:#6d28d9 }
        .cm-badge-nacional { background:#dcfce7;color:#166534 }
        .cm-item-name { font-size:.92rem;font-weight:700;color:#111827 }
        .cm-ubicacion { font-size:.78rem;color:#6b7280;display:flex;align-items:center;gap:5px;margin-bottom:10px }
        .cm-ubicacion i { color:#dc2626 }
        .cm-compromiso-label { font-size:.72rem;font-weight:600;color:#6b7280;margin-bottom:2px }
        .cm-compromiso-text { font-size:.84rem;color:#374151;margin-bottom:10px }
        .cm-comments { margin-top:8px;padding-top:8px;border-top:1px solid #f3f4f6 }
        .cm-comments-toggle { background:none;border:none;padding:0;font-size:.78rem;font-weight:600;color:#374151;cursor:pointer }
        .cm-comments-empty { font-size:.8rem;color:#9ca3af;font-style:italic;margin:6px 0 8px }
        .cm-comment { font-size:.82rem;color:#374151;margin:8px 0;padding:8px 10px;background:#f9fafb;border-radius:8px }
        .cm-comment-meta { font-size:.7rem;color:#9ca3af;margin-top:4px }
        .cm-comment-form { display:flex;flex-direction:column;gap:6px;margin-top:8px }
        .cm-comment-form textarea { width:100%;font-size:.82rem;border:1px solid #d1d5db;border-radius:8px;padding:8px 10px;outline:none;box-sizing:border-box;font-family:inherit;resize:vertical;min-height:56px }
        .cm-comment-form textarea:focus { border-color:#5106A7 }
        .cm-btn-comentar { align-self:flex-end;background:#5106A7;color:#fff;border:none;border-radius:8px;padding:6px 12px;font-size:.75rem;font-weight:600;cursor:pointer }
        .cm-btn-comentar:disabled { opacity:.5;cursor:not-allowed }
        .cm-modal-overlay { position:fixed;inset:0;background:rgba(0,0,0,.4);display:flex;align-items:center;justify-content:center;z-index:1000;padding:16px }
        .cm-modal { background:#fff;border-radius:14px;width:100%;max-width:560px;max-height:90vh;overflow-y:auto;box-shadow:0 10px 40px rgba(0,0,0,.2) }
        .cm-modal-header { display:flex;align-items:flex-start;justify-content:space-between;padding:18px 20px 8px }
        .cm-modal-header h3 { margin:0;font-size:1.05rem;font-weight:700;color:#111827 }
        .cm-modal-sub { margin:4px 0 0;font-size:.8rem;color:#6b7280 }
        .cm-modal-close { background:none;border:none;font-size:1.2rem;color:#9ca3af;cursor:pointer;line-height:1 }
        .cm-modal-body { padding:8px 20px 16px;display:flex;flex-direction:column;gap:14px }
        .cm-field label { display:block;font-size:.78rem;font-weight:600;color:#374151;margin-bottom:4px }
        .cm-req { color:#dc2626 }
        .cm-field select, .cm-field textarea { width:100%;font-size:.85rem;border:1px solid #d1d5db;border-radius:8px;padding:8px 10px;outline:none;box-sizing:border-box;background:#fff }
        .cm-field select:focus, .cm-field textarea:focus { border-color:#5106A7 }
        .cm-field textarea { resize:vertical;min-height:80px;font-family:inherit }
        .cm-row { display:grid;grid-template-columns:1fr 1fr 1fr;gap:10px }
        @media (max-width: 640px) { .cm-row { grid-template-columns:1fr } }
        .cm-hint { font-size:.72rem;color:#9ca3af;margin-top:4px }
        .cm-field-error { color:#dc2626;font-size:.72rem;margin-top:3px }
        .cm-modal-footer { display:flex;justify-content:flex-end;gap:10px;padding:14px 20px;border-top:1px solid #f3f4f6 }
        .cm-btn-cancel { background:#fff;border:1px solid #d1d5db;color:#374151;border-radius:8px;padding:8px 16px;font-size:.82rem;font-weight:600;cursor:pointer }
        .cm-btn-guardar { background:#5106A7;color:#fff;border:none;border-radius:8px;padding:8px 16px;font-size:.82rem;font-weight:600;cursor:pointer }
        .cm-btn-guardar:disabled { opacity:.5;cursor:not-allowed }
    `;
    document.head.appendChild(style);
})();

app.component('case-mecanismos', {
    delimiters: ['${', '}'],
    props: {
        caseId: { type: String, required: true },
        userId: { type: String, required: true },
        userRole: { type: String, required: true },
    },
    data: function() {
        return {
            cargando: false,
            error: null,
            mecanismos: [],
            commentsOpen: {},

            modalAbierto: false,
            guardando: false,
            errores: {},
            departamentos: [],
            ciudades: [],
            municipios: [],
            usuarios: [],
            cargandoUsuarios: false,
            form: {
                scope: 'municipal',
                departmentId: '',
                cityId: '',
                townCode: '',
                assignedUserId: '',
                salviaCommitment: '',
            },
        };
    },
    computed: {
        puedeAgregar: function() {
            return this.userRole === 'sv' || this.userRole === 'op' || this.userRole === 'ro';
        },
        formValido: function() {
            if (!this.form.scope || !this.form.assignedUserId) return false;
            if (this.form.scope === 'municipal') {
                return !!(this.form.departmentId && this.form.cityId && this.form.townCode);
            }
            if (this.form.scope === 'departamental') {
                return !!this.form.departmentId;
            }
            return true;
        },
        mechanismSelectLabel: function() {
            if (this.form.scope === 'municipal') return 'Mecanismo Municipal';
            if (this.form.scope === 'departamental') return 'Mecanismo Departamental';
            if (this.form.scope === 'nacional') return 'Mecanismo Nacional';
            return 'Mecanismo';
        },
        usersReady: function() {
            if (this.form.scope === 'nacional') return true;
            if (this.form.scope === 'departamental') return !!this.form.departmentId;
            if (this.form.scope === 'municipal') {
                return !!(this.form.departmentId && this.form.cityId && this.form.townCode);
            }
            return false;
        },
    },
    mounted: function() {
        this.cargar();
    },
    methods: {
        formVacio: function() {
            return {
                scope: 'municipal',
                departmentId: '',
                cityId: '',
                townCode: '',
                assignedUserId: '',
                salviaCommitment: '',
            };
        },
        cargar: function() {
            var self = this;
            self.cargando = true;
            self.error = null;
            fetch('/api/v1/casos/' + self.caseId + '/mecanismos')
                .then(function(res) {
                    if (!res.ok) throw new Error('Error ' + res.status);
                    return res.json();
                })
                .then(function(data) {
                    self.mecanismos = Array.isArray(data) ? data : [];
                    self.cargando = false;
                })
                .catch(function() {
                    self.error = 'No se pudieron cargar los mecanismos de este caso.';
                    self.cargando = false;
                });
        },
        scopeLabel: function(scope) {
            var map = { municipal: 'Municipal', departamental: 'Departamental', nacional: 'Nacional' };
            return map[scope] || scope;
        },
        badgeClass: function(scope) {
            return 'cm-badge cm-badge-' + (scope || '');
        },
        locationText: function(m) {
            var parts = [];
            if (m.townName) parts.push(m.townName);
            if (m.cityName) parts.push(m.cityName);
            if (m.departmentName) parts.push(m.departmentName);
            if (parts.length >= 3) return parts[0] + ', ' + parts[1] + ' (' + parts[2] + ')';
            return parts.join(', ');
        },
        formatFecha: function(iso) {
            if (!iso) return '';
            var d = new Date(iso);
            if (isNaN(d.getTime())) return iso;
            var pad = function(n) { return n < 10 ? '0' + n : '' + n; };
            return pad(d.getDate()) + '/' + pad(d.getMonth() + 1) + '/' + d.getFullYear() +
                ' ' + pad(d.getHours()) + ':' + pad(d.getMinutes());
        },
        toggleComments: function(id) {
            this.commentsOpen[id] = this.commentsOpen[id] === false;
        },
        commentsOf: function(m) {
            return Array.isArray(m.comments) ? m.comments : [];
        },

        abrirModal: function() {
            this.form = this.formVacio();
            this.ciudades = [];
            this.municipios = [];
            this.usuarios = [];
            this.errores = {};
            this.guardando = false;
            this.modalAbierto = true;
            if (this.departamentos.length === 0) this.cargarDepartamentos();
            this.cargarUsuarios();
        },
        cerrarModal: function() {
            this.modalAbierto = false;
        },
        cargarDepartamentos: function() {
            var self = this;
            fetch('/api/v1/locations/departments')
                .then(function(res) { return res.json(); })
                .then(function(data) { self.departamentos = Array.isArray(data) ? data : []; })
                .catch(function() { self.departamentos = []; });
        },
        onCambiaAlcance: function() {
            this.form.departmentId = '';
            this.form.cityId = '';
            this.form.townCode = '';
            this.form.assignedUserId = '';
            this.ciudades = [];
            this.municipios = [];
            this.usuarios = [];
            this.errores = {};
            this.cargarUsuarios();
        },
        onCambiaDepartamento: function() {
            this.form.cityId = '';
            this.form.townCode = '';
            this.form.assignedUserId = '';
            this.ciudades = [];
            this.municipios = [];
            this.usuarios = [];
            if (!this.form.departmentId) return;
            if (this.form.scope === 'municipal') {
                var self = this;
                fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(this.form.departmentId))
                    .then(function(res) { return res.json(); })
                    .then(function(data) { self.ciudades = Array.isArray(data) ? data : []; })
                    .catch(function() { self.ciudades = []; });
            }
            this.cargarUsuarios();
        },
        onCambiaCiudad: function() {
            this.form.townCode = '';
            this.form.assignedUserId = '';
            this.municipios = [];
            this.usuarios = [];
            if (!this.form.cityId) {
                this.cargarUsuarios();
                return;
            }
            var self = this;
            fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(this.form.cityId))
                .then(function(res) { return res.json(); })
                .then(function(data) { self.municipios = Array.isArray(data) ? data : []; })
                .catch(function() { self.municipios = []; });
            this.cargarUsuarios();
        },
        onCambiaMunicipio: function() {
            this.form.assignedUserId = '';
            this.cargarUsuarios();
        },
        cargarUsuarios: function() {
            if (!this.usersReady) {
                this.usuarios = [];
                return;
            }
            var self = this;
            var url = '/api/v1/mecanismos/usuarios?scope=' + encodeURIComponent(this.form.scope);
            if (this.form.departmentId) url += '&department_id=' + encodeURIComponent(this.form.departmentId);
            if (this.form.cityId) url += '&city_id=' + encodeURIComponent(this.form.cityId);
            if (this.form.townCode) url += '&town_id=' + encodeURIComponent(this.form.townCode);
            self.cargandoUsuarios = true;
            fetch(url)
                .then(function(res) {
                    if (!res.ok) throw new Error('Error ' + res.status);
                    return res.json();
                })
                .then(function(data) {
                    self.usuarios = Array.isArray(data) ? data : [];
                    self.cargandoUsuarios = false;
                })
                .catch(function() {
                    self.usuarios = [];
                    self.cargandoUsuarios = false;
                });
        },
        guardar: function() {
            var self = this;
            var errores = {};
            if (!this.form.scope) errores.scope = 'Selecciona el tipo de mecanismo';
            if (this.form.scope === 'municipal') {
                if (!this.form.departmentId) errores.departmentId = 'Selecciona el departamento';
                if (!this.form.cityId) errores.cityId = 'Selecciona la ciudad';
                if (!this.form.townCode) errores.townCode = 'Selecciona el municipio';
            }
            if (this.form.scope === 'departamental' && !this.form.departmentId) {
                errores.departmentId = 'Selecciona el departamento';
            }
            if (!this.form.assignedUserId) errores.assignedUserId = 'Selecciona un mecanismo';
            this.errores = errores;
            if (Object.keys(errores).length > 0) return;

            self.guardando = true;
            fetch('/api/v1/casos/' + self.caseId + '/mecanismos', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    scope: self.form.scope,
                    assignedUserId: self.form.assignedUserId,
                    departmentId: self.form.departmentId || '',
                    cityId: self.form.cityId || '',
                    townCode: self.form.townCode || '',
                    salviaCommitment: (self.form.salviaCommitment || '').trim(),
                }),
            })
                .then(function(res) {
                    if (!res.ok) {
                        return res.json().then(function(d) {
                            throw new Error(d.error || 'Error ' + res.status);
                        });
                    }
                    return res.json();
                })
                .then(function() {
                    self.guardando = false;
                    self.modalAbierto = false;
                    self.cargar();
                })
                .catch(function(err) {
                    self.guardando = false;
                    self.errores.assignedUserId = err.message || 'Error al guardar. Intente nuevamente.';
                });
        },
    },
    template: `
        <div class="cm-card">
            <div class="cm-header">
                <h3 class="cm-title"><i class="fas fa-balance-scale"></i> Mecanismos</h3>
                <button v-if="puedeAgregar" class="cm-btn-agregar" @click="abrirModal">+ Agregar mecanismo</button>
            </div>

            <div v-if="cargando" class="cm-loading">Cargando mecanismos...</div>
            <div v-else-if="error" class="cm-error">
                <span>\${ error }</span>
                <button @click="cargar">Reintentar</button>
            </div>
            <div v-else-if="mecanismos.length === 0" class="cm-empty">
                <div class="cm-empty-icon"><i class="fas fa-balance-scale"></i></div>
                <p>Este caso no ha sido escalado a ningún mecanismo articulador</p>
            </div>
            <div v-else class="cm-list">
                <div v-for="m in mecanismos" :key="m.id" class="cm-item">
                    <div class="cm-item-head">
                        <span :class="badgeClass(m.scope)">\${ scopeLabel(m.scope) }</span>
                        <span class="cm-item-name">\${ m.assignedUserName }</span>
                    </div>
                    <div class="cm-ubicacion" v-if="locationText(m)">
                        <i class="fas fa-map-marker-alt"></i>
                        \${ locationText(m) }
                    </div>
                    <template v-if="m.salviaCommitment">
                        <div class="cm-compromiso-label">Compromiso de Salvia</div>
                        <div class="cm-compromiso-text">\${ m.salviaCommitment }</div>
                    </template>
                    <div class="cm-comments">
                        <button class="cm-comments-toggle" @click="toggleComments(m.id)">
                            Comentarios del mecanismo (\${ commentsOf(m).length })
                        </button>
                        <div v-if="commentsOpen[m.id] !== false">
                            <div v-if="commentsOf(m).length === 0" class="cm-comments-empty">Sin comentarios aún</div>
                            <div v-for="c in commentsOf(m)" :key="c.id" class="cm-comment">
                                <div>\${ c.body }</div>
                                <div class="cm-comment-meta">\${ c.authorName } · \${ formatFecha(c.createdAt) }</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="modalAbierto" class="cm-modal-overlay" @click.self="cerrarModal">
                <div class="cm-modal">
                    <div class="cm-modal-header">
                        <div>
                            <h3>Agregar mecanismo articulador</h3>
                            <p class="cm-modal-sub">Instancia a la que se escaló o se va a escalar este caso</p>
                        </div>
                        <button class="cm-modal-close" @click="cerrarModal">×</button>
                    </div>
                    <div class="cm-modal-body">
                        <div class="cm-field">
                            <label>Tipo mecanismo <span class="cm-req">*</span></label>
                            <select v-model="form.scope" @change="onCambiaAlcance">
                                <option value="municipal">Municipal</option>
                                <option value="departamental">Departamental</option>
                                <option value="nacional">Nacional</option>
                            </select>
                        </div>

                        <div v-if="form.scope === 'municipal'" class="cm-row">
                            <div class="cm-field">
                                <label>Departamento <span class="cm-req">*</span></label>
                                <select v-model="form.departmentId" @change="onCambiaDepartamento">
                                    <option value="">Selecciona...</option>
                                    <option v-for="d in departamentos" :key="d.value" :value="d.value">\${ d.label }</option>
                                </select>
                                <div v-if="errores.departmentId" class="cm-field-error">\${ errores.departmentId }</div>
                            </div>
                            <div class="cm-field">
                                <label>Ciudad <span class="cm-req">*</span></label>
                                <select v-model="form.cityId" :disabled="!form.departmentId" @change="onCambiaCiudad">
                                    <option value="">Selecciona...</option>
                                    <option v-for="c in ciudades" :key="c.value" :value="c.value">\${ c.label }</option>
                                </select>
                                <div v-if="errores.cityId" class="cm-field-error">\${ errores.cityId }</div>
                            </div>
                            <div class="cm-field">
                                <label>Municipio <span class="cm-req">*</span></label>
                                <select v-model="form.townCode" :disabled="!form.cityId" @change="onCambiaMunicipio">
                                    <option value="">Selecciona...</option>
                                    <option v-for="t in municipios" :key="t.value" :value="t.value">\${ t.label }</option>
                                </select>
                                <div v-if="errores.townCode" class="cm-field-error">\${ errores.townCode }</div>
                            </div>
                        </div>

                        <template v-if="form.scope === 'departamental'">
                            <div class="cm-field">
                                <label>Departamento <span class="cm-req">*</span></label>
                                <select v-model="form.departmentId" @change="onCambiaDepartamento">
                                    <option value="">Selecciona...</option>
                                    <option v-for="d in departamentos" :key="'d-'+d.value" :value="d.value">\${ d.label }</option>
                                </select>
                                <div v-if="errores.departmentId" class="cm-field-error">\${ errores.departmentId }</div>
                            </div>
                        </template>

                        <div v-if="form.scope" class="cm-field">
                            <label>\${ mechanismSelectLabel } <span class="cm-req">*</span></label>
                            <select v-model="form.assignedUserId" :disabled="!usersReady || cargandoUsuarios">
                                <option value="">\${ cargandoUsuarios ? 'Cargando...' : 'Selecciona...' }</option>
                                <option v-for="u in usuarios" :key="u.icode" :value="u.icode">\${ u.fullName }</option>
                            </select>
                            <div v-if="usersReady && !cargandoUsuarios && usuarios.length === 0" class="cm-hint">
                                No hay usuarios mecanismo en esta ubicación
                            </div>
                            <div v-if="errores.assignedUserId" class="cm-field-error">\${ errores.assignedUserId }</div>
                        </div>

                        <div class="cm-field">
                            <label>Compromiso de Salvia (opcional)</label>
                            <textarea v-model="form.salviaCommitment" maxlength="2000"
                                placeholder="Si ya se adquirió un compromiso, regístralo aquí..."></textarea>
                        </div>
                    </div>
                    <div class="cm-modal-footer">
                        <button class="cm-btn-cancel" @click="cerrarModal">Cancelar</button>
                        <button class="cm-btn-guardar" :disabled="!formValido || guardando" @click="guardar">
                            \${ guardando ? 'Guardando...' : 'Agregar mecanismo' }
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `,
});
