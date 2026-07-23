(function injectCaseEntitiesStyles() {
    if (document.getElementById('ce-styles')) return;
    var style = document.createElement('style');
    style.id = 'ce-styles';
    style.textContent = `
        .ce-container { padding: 4px 0 }
        .ce-loading { text-align:center;padding:32px 0;color:#9ca3af;font-size:.85rem }
        .ce-error { background:#fef2f2;border:1px solid #fecaca;border-radius:8px;padding:12px 16px;color:#dc2626;font-size:.82rem;display:flex;align-items:center;justify-content:space-between;gap:12px }
        .ce-error button { background:#dc2626;color:#fff;border:none;border-radius:6px;padding:5px 12px;font-size:.75rem;font-weight:600;cursor:pointer;white-space:nowrap }
        .ce-empty { text-align:center;padding:32px 0;color:#9ca3af;font-size:.84rem }
        .ce-filtros { display:flex;align-items:center;gap:10px;flex-wrap:wrap;margin-bottom:16px }
        .ce-filtros select, .ce-buscador { font-size:.82rem;border:1px solid #d1d5db;border-radius:8px;padding:7px 10px;background:#fff;outline:none }
        .ce-buscador { flex:1;min-width:180px }
        .ce-btn-agregar { margin-left:auto;background:#5106A7;color:#fff;border:none;border-radius:8px;padding:8px 14px;font-size:.8rem;font-weight:600;cursor:pointer;white-space:nowrap }
        .ce-grid { display:grid;grid-template-columns:repeat(auto-fill, minmax(280px, 1fr));gap:14px }
        .ce-card { background:#fff;border:1px solid #e5e7eb;border-radius:12px;padding:14px 16px;display:flex;flex-direction:column;gap:6px }
        .ce-card-header { display:flex;align-items:center;gap:8px;flex-wrap:wrap }
        .ce-badge-sector { font-size:.68rem;font-weight:700;padding:2px 8px;border-radius:6px;background:#ede9fe;color:#5106A7;text-transform:uppercase;letter-spacing:.03em }
        .ce-nombre { font-size:.9rem;font-weight:700;color:#111827 }
        .ce-ubicacion { font-size:.78rem;color:#6b7280 }
        .ce-direccion { font-size:.78rem;color:#6b7280 }
        .ce-objetivo { font-size:.82rem;color:#374151;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden }
        .ce-stats { display:flex;flex-wrap:wrap;gap:10px;font-size:.74rem;color:#6b7280;margin-top:4px }
        .ce-actions { margin-top:8px;display:flex;justify-content:flex-end }
        .ce-btn-detalle { background:none;border:none;color:#5106A7;font-size:.78rem;font-weight:600;cursor:pointer;padding:0 }
        .ce-modal-overlay { position:fixed;inset:0;background:rgba(0,0,0,.4);display:flex;align-items:center;justify-content:center;z-index:1000;padding:16px }
        .ce-modal { background:#fff;border-radius:14px;width:100%;max-width:480px;max-height:90vh;overflow-y:auto;box-shadow:0 10px 40px rgba(0,0,0,.2) }
        .ce-modal-header { display:flex;align-items:center;justify-content:space-between;padding:16px 20px;border-bottom:1px solid #f3f4f6 }
        .ce-modal-header h3 { margin:0;font-size:1rem;font-weight:700;color:#111827 }
        .ce-modal-close { background:none;border:none;font-size:1.2rem;color:#9ca3af;cursor:pointer;line-height:1 }
        .ce-modal-body { padding:16px 20px;display:flex;flex-direction:column;gap:14px }
        .ce-field label { display:block;font-size:.78rem;font-weight:600;color:#374151;margin-bottom:4px }
        .ce-field select, .ce-field textarea { width:100%;font-size:.85rem;border:1px solid #d1d5db;border-radius:8px;padding:8px 10px;outline:none;box-sizing:border-box }
        .ce-field textarea { resize:vertical;min-height:70px;font-family:inherit }
        .ce-field-error { color:#dc2626;font-size:.72rem;margin-top:3px }
        .ce-modal-footer { display:flex;justify-content:flex-end;gap:10px;padding:14px 20px;border-top:1px solid #f3f4f6 }
        .ce-btn-cancel { background:#fff;border:1px solid #d1d5db;color:#374151;border-radius:8px;padding:8px 16px;font-size:.82rem;font-weight:600;cursor:pointer }
        .ce-btn-guardar { background:#5106A7;color:#fff;border:none;border-radius:8px;padding:8px 16px;font-size:.82rem;font-weight:600;cursor:pointer }
        .ce-btn-guardar:disabled { opacity:.5;cursor:not-allowed }
        .ce-hint { font-size:.75rem;color:#9ca3af }
    `;
    document.head.appendChild(style);
})();

// Sectores para el filtro de la lista — codigos reales de entity.entity_sector.
var CE_SECTOR_LABELS = { he: 'Salud', js: 'Justicia', pt: 'Protección', os: 'Otras instituciones' };

// Sectores para el buscador del modal — GET /api/v1/entity-branches espera estas
// palabras (no los códigos de 2 letras), ver entity_branch_api_controller.go.
var CE_SECTORES_MODAL = [
    { value: 'salud', label: 'Salud' },
    { value: 'justicia', label: 'Justicia' },
    { value: 'proteccion', label: 'Protección' },
    { value: 'otras_instituciones', label: 'Otras instituciones' },
];

app.component('case-entities', {
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
            entidades: [],

            busqueda: '',
            filtroSector: '',
            filtroBarrerasActivas: '',

            sectorLabels: CE_SECTOR_LABELS,
            sectoresModal: CE_SECTORES_MODAL,

            modalAgregarAbierto: false,
            departamentos: [],
            ciudades: [],
            municipios: [],
            entidadesDisponibles: [],
            cargandoEntidades: false,
            guardando: false,
            errores: {},
            form: { departamento: '', ciudad: '', municipio: '', sector: '', entidadId: '', objetivo: '' },
        };
    },
    computed: {
        entidadesFiltradas: function() {
            var self = this;
            var q = this.busqueda.trim().toLowerCase();
            return this.entidades.filter(function(e) {
                var porSector = !self.filtroSector || e.sector === self.filtroSector;
                var porTexto = !q || (e.entityBranchName || '').toLowerCase().includes(q);
                var porBarreras =
                    !self.filtroBarrerasActivas ? true :
                    self.filtroBarrerasActivas === 'si' ? e.barrerasActivasCount > 0 :
                    e.barrerasActivasCount === 0;
                return porSector && porTexto && porBarreras;
            });
        },
        formValido: function() {
            return !!this.form.entidadId && !!this.form.objetivo.trim();
        },
        // Cualquier rol con acceso a la pantalla de Detalle del Caso puede ver/filtrar
        // entidades (ya lo garantiza el backend al montar este componente) — pero solo
        // sv/op/ro pueden agregar una nueva. Ver case-entities-interface.md.
        puedeAgregar: function() {
            return this.userRole === 'sv' || this.userRole === 'op' || this.userRole === 'ro';
        },
    },
    mounted: function() {
        this.cargar();
    },
    methods: {
        cargar: function() {
            var self = this;
            self.cargando = true;
            self.error = null;
            fetch('/api/v1/casos/' + self.caseId + '/entidades')
                .then(function(res) {
                    if (!res.ok) throw new Error('Error ' + res.status);
                    return res.json();
                })
                .then(function(data) {
                    self.entidades = Array.isArray(data) ? data : [];
                    self.cargando = false;
                })
                .catch(function() {
                    self.error = 'No se pudieron cargar las entidades relacionadas con este caso.';
                    self.cargando = false;
                });
        },

        sectorLabel: function(code) {
            return this.sectorLabels[code] || code || '—';
        },

        // ── Modal: agregar entidad ──────────────────────────────────────────────
        abrirModalAgregar: function() {
            this.form = { departamento: '', ciudad: '', municipio: '', sector: '', entidadId: '', objetivo: '' };
            this.ciudades = [];
            this.municipios = [];
            this.entidadesDisponibles = [];
            this.cargandoEntidades = false;
            this.errores = {};
            this.guardando = false;
            this.modalAgregarAbierto = true;
            if (this.departamentos.length === 0) this.cargarDepartamentos();
        },
        cerrarModalAgregar: function() {
            this.modalAgregarAbierto = false;
        },
        cargarDepartamentos: function() {
            var self = this;
            fetch('/api/v1/locations/departments')
                .then(function(res) { return res.json(); })
                .then(function(data) { self.departamentos = Array.isArray(data) ? data : []; })
                .catch(function() { self.departamentos = []; });
        },
        onCambiaDepartamento: function() {
            this.form.ciudad = '';
            this.form.municipio = '';
            this.ciudades = [];
            this.municipios = [];
            this.entidadesDisponibles = [];
            if (!this.form.departamento) return;
            var self = this;
            fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(this.form.departamento))
                .then(function(res) { return res.json(); })
                .then(function(data) { self.ciudades = Array.isArray(data) ? data : []; })
                .catch(function() { self.ciudades = []; });
        },
        onCambiaCiudad: function() {
            this.form.municipio = '';
            this.municipios = [];
            this.entidadesDisponibles = [];
            if (!this.form.ciudad) return;
            var self = this;
            fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(this.form.ciudad))
                .then(function(res) { return res.json(); })
                .then(function(data) { self.municipios = Array.isArray(data) ? data : []; })
                .catch(function() { self.municipios = []; });
        },
        onCambiaMunicipio: function() {
            this.entidadesDisponibles = [];
            this.form.entidadId = '';
            this.buscarEntidadesDisponibles();
        },
        onCambiaSectorModal: function() {
            this.entidadesDisponibles = [];
            this.form.entidadId = '';
            this.buscarEntidadesDisponibles();
        },
        buscarEntidadesDisponibles: function() {
            if (!this.form.municipio || !this.form.sector) return;
            var self = this;
            self.cargandoEntidades = true;
            fetch('/api/v1/entity-branches?town_code=' + encodeURIComponent(this.form.municipio) + '&sector=' + encodeURIComponent(this.form.sector))
                .then(function(res) { return res.json(); })
                .then(function(data) {
                    self.entidadesDisponibles = Array.isArray(data) ? data : [];
                    self.cargandoEntidades = false;
                })
                .catch(function() {
                    self.entidadesDisponibles = [];
                    self.cargandoEntidades = false;
                });
        },
        guardarEntidad: function() {
            var self = this;
            var errores = {};
            if (!this.form.entidadId) errores.entidad = 'Selecciona una entidad';
            if (!this.form.objetivo.trim()) errores.objetivo = 'Describe el objetivo con la entidad';
            this.errores = errores;
            if (Object.keys(errores).length > 0) return;

            self.guardando = true;
            fetch('/api/v1/casos/' + self.caseId + '/entidades', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    entityBranchId: Number(self.form.entidadId),
                    objetivo: self.form.objetivo.trim(),
                    createdById: self.userId,
                }),
            })
                .then(function(res) {
                    if (res.status === 409) {
                        return res.json().then(function(d) { throw new Error(d.error || 'Esta entidad ya está asociada a este caso'); });
                    }
                    if (!res.ok) throw new Error('Error ' + res.status);
                    return res.json();
                })
                .then(function() {
                    self.guardando = false;
                    self.modalAgregarAbierto = false;
                    self.cargar();
                })
                .catch(function(err) {
                    self.guardando = false;
                    self.errores.entidad = err.message || 'Error al guardar. Intente nuevamente.';
                });
        },

        onVerDetalle: function(entidad) {
            window.location.href = '/salvia/entidad/' + entidad.entityBranchICode;
        },
    },
    template: `
        <div class="ce-container">
            <div v-if="cargando" class="ce-loading">Cargando entidades...</div>
            <div v-else-if="error" class="ce-error">
                <span>\${ error }</span>
                <button @click="cargar">Reintentar</button>
            </div>
            <template v-else>
                <div class="ce-filtros">
                    <select v-model="filtroSector">
                        <option value="">Todos los sectores</option>
                        <option v-for="(label, code) in sectorLabels" :key="code" :value="code">\${ label }</option>
                    </select>
                    <input type="text" class="ce-buscador" v-model="busqueda" placeholder="Buscar por nombre de entidad..." />
                    <select v-model="filtroBarrerasActivas">
                        <option value="">¿Con barreras activas?</option>
                        <option value="si">Sí</option>
                        <option value="no">No</option>
                    </select>
                    <button v-if="puedeAgregar" class="ce-btn-agregar" @click="abrirModalAgregar">+ Agregar entidad</button>
                </div>

                <div v-if="entidadesFiltradas.length === 0" class="ce-empty">
                    \${ entidades.length === 0 ? 'No hay entidades registradas para este caso' : 'No hay entidades que coincidan con los filtros' }
                </div>
                <div v-else class="ce-grid">
                    <div v-for="entidad in entidadesFiltradas" :key="entidad.relId" class="ce-card">
                        <div class="ce-card-header">
                            <span class="ce-badge-sector">\${ sectorLabel(entidad.sector) }</span>
                            <span class="ce-nombre">\${ entidad.entityBranchName }</span>
                        </div>
                        <div class="ce-ubicacion" v-if="entidad.townName || entidad.cityName || entidad.departmentName">
                            📍 \${ [entidad.departmentName, entidad.cityName, entidad.townName].filter(Boolean).join(' · ') }
                        </div>
                        <div class="ce-direccion" v-if="entidad.address">\${ entidad.address }</div>
                        <div class="ce-objetivo" v-if="entidad.objetivo">\${ entidad.objetivo }</div>
                        <div class="ce-stats">
                            <span>Oficios: \${ entidad.oficiosCount }</span>
                            <span>Barreras activas: \${ entidad.barrerasActivasCount }</span>
                            <span v-if="entidad.lastAction">Última acción: \${ entidad.lastAction }</span>
                            <span v-else>Sin acciones registradas</span>
                        </div>
                        <div class="ce-actions">
                            <button class="ce-btn-detalle" @click="onVerDetalle(entidad)">Ver Detalle →</button>
                        </div>
                    </div>
                </div>
            </template>

            <div v-if="modalAgregarAbierto" class="ce-modal-overlay" @click.self="cerrarModalAgregar">
                <div class="ce-modal">
                    <div class="ce-modal-header">
                        <h3>Agregar entidad al caso</h3>
                        <button class="ce-modal-close" @click="cerrarModalAgregar">×</button>
                    </div>
                    <div class="ce-modal-body">
                        <div class="ce-field">
                            <label>Departamento *</label>
                            <select v-model="form.departamento" @change="onCambiaDepartamento">
                                <option value="">Seleccione el departamento</option>
                                <option v-for="d in departamentos" :key="d.value" :value="d.value">\${ d.label }</option>
                            </select>
                        </div>
                        <div class="ce-field">
                            <label>Ciudad *</label>
                            <select v-model="form.ciudad" @change="onCambiaCiudad" :disabled="!form.departamento">
                                <option value="">Seleccione la ciudad</option>
                                <option v-for="c in ciudades" :key="c.value" :value="c.value">\${ c.label }</option>
                            </select>
                        </div>
                        <div class="ce-field">
                            <label>Municipio *</label>
                            <select v-model="form.municipio" @change="onCambiaMunicipio" :disabled="!form.ciudad">
                                <option value="">Seleccione el municipio</option>
                                <option v-for="m in municipios" :key="m.value" :value="m.value">\${ m.label }</option>
                            </select>
                        </div>
                        <div class="ce-field">
                            <label>Sector *</label>
                            <select v-model="form.sector" @change="onCambiaSectorModal" :disabled="!form.municipio">
                                <option value="">Seleccione el sector</option>
                                <option v-for="s in sectoresModal" :key="s.value" :value="s.value">\${ s.label }</option>
                            </select>
                        </div>

                        <div v-if="cargandoEntidades" class="ce-hint">Buscando sedes disponibles...</div>
                        <div v-else-if="form.sector && entidadesDisponibles.length === 0" class="ce-hint">
                            No hay sedes registradas en este municipio y sector. Puedes crear una desde Sedes.
                        </div>

                        <div class="ce-field" v-if="entidadesDisponibles.length > 0">
                            <label>Entidad *</label>
                            <select v-model="form.entidadId">
                                <option value="">Seleccione la entidad</option>
                                <option v-for="opt in entidadesDisponibles" :key="opt.id" :value="opt.id">\${ opt.name }</option>
                            </select>
                            <div v-if="errores.entidad" class="ce-field-error">\${ errores.entidad }</div>
                        </div>

                        <div class="ce-field">
                            <label>Objetivo con la entidad *</label>
                            <textarea v-model="form.objetivo" placeholder="¿Qué debe gestionar la víctima en esta entidad?"></textarea>
                            <div v-if="errores.objetivo" class="ce-field-error">\${ errores.objetivo }</div>
                        </div>
                    </div>
                    <div class="ce-modal-footer">
                        <button class="ce-btn-cancel" @click="cerrarModalAgregar">Cancelar</button>
                        <button class="ce-btn-guardar" :disabled="!formValido || guardando" @click="guardarEntidad">
                            \${ guardando ? 'Guardando...' : 'Agregar entidad' }
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `,
});
