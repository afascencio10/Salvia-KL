/**
 * case-info
 * Componente reutilizable: información completa de un caso.
 * Soporta modo "modal" (overlay) e "inline" (embebido en la página).
 *
 * Props:
 *   - caseId  (String, required): icode del caso
 *   - visible (Boolean, required): controla visibilidad
 *   - mode    (String, default 'modal'): 'modal' | 'inline'
 * Events:
 *   - @close: emitido al cerrar (solo en modo modal)
 */

(function injectCaseInfoStyles() {
    if (document.getElementById('case-info-styles')) return;
    var style = document.createElement('style');
    style.id = 'case-info-styles';
    style.textContent = `
        .ci-backdrop { position:fixed;inset:0;background:rgba(0,0,0,.5);display:flex;align-items:center;justify-content:center;z-index:10000;animation:ciFadeIn .2s ease }
        @keyframes ciFadeIn { from{opacity:0} to{opacity:1} }
        .ci-modal { background:#fff;border-radius:16px;width:95%;max-width:1000px;max-height:90vh;display:flex;overflow:hidden;box-shadow:0 25px 60px rgba(0,0,0,.2);animation:ciSlideUp .25s ease }
        @keyframes ciSlideUp { from{transform:translateY(20px);opacity:0} to{transform:translateY(0);opacity:1} }
        .ci-sidebar { width:220px;background:#f9fafb;border-right:1px solid #f3f4f6;padding:20px 0;overflow-y:auto;flex-shrink:0 }
        .ci-sidebar-title { font-size:.68rem;font-weight:700;color:#9ca3af;text-transform:uppercase;letter-spacing:.05em;padding:0 16px;margin:0 0 12px }
        .ci-nav-item { display:flex;align-items:center;gap:8px;padding:8px 16px;font-size:.78rem;color:#6b7280;cursor:pointer;border-left:3px solid transparent;transition:all .12s }
        .ci-nav-item:hover { background:#f3f4f6;color:#374151 }
        .ci-nav-item.active { background:#ede9fe;color:#5106A7;border-left-color:#5106A7;font-weight:600 }
        .ci-nav-item i { width:16px;text-align:center;font-size:.75rem }
        .ci-main { flex:1;overflow-y:auto;display:flex;flex-direction:column }
        .ci-header { position:sticky;top:0;z-index:1;background:#fff;display:flex;align-items:center;justify-content:space-between;padding:18px 24px 14px;border-bottom:1px solid #f3f4f6 }
        .ci-title { font-size:1.05rem;font-weight:700;color:#111827;margin:0;display:flex;align-items:center;gap:8px }
        .ci-close { background:none;border:none;cursor:pointer;color:#9ca3af;font-size:1.1rem;padding:4px 8px;border-radius:6px;line-height:1 }
        .ci-close:hover { background:#f3f4f6;color:#374151 }
        .ci-body { padding:20px 24px;flex:1 }
        .ci-section { margin-bottom:8px;border:1px solid #e5e7eb;border-radius:10px;overflow:hidden }
        .ci-section-header { display:flex;align-items:center;gap:10px;padding:12px 16px;cursor:pointer;background:#fff;transition:background .12s;user-select:none }
        .ci-section-header:hover { background:#f9fafb }
        .ci-section-header.open { background:#faf5ff;border-bottom:1px solid #ede9fe }
        .ci-section-icon { width:28px;height:28px;border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:.75rem;color:#fff;flex-shrink:0 }
        .ci-section-label { font-size:.82rem;font-weight:600;color:#374151;flex:1 }
        .ci-section-count { font-size:.68rem;color:#9ca3af;background:#f3f4f6;padding:2px 8px;border-radius:10px }
        .ci-section-chevron { font-size:.7rem;color:#9ca3af;transition:transform .2s }
        .ci-section-header.open .ci-section-chevron { transform:rotate(180deg) }
        .ci-section-body { padding:14px 16px;background:#fff }
        .ci-grid { display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:14px 20px }
        .ci-field { display:flex;flex-direction:column;gap:2px }
        .ci-label { font-size:.68rem;font-weight:600;color:#9ca3af;text-transform:uppercase;letter-spacing:.03em }
        .ci-value { font-size:.83rem;font-weight:500;color:#111827;word-break:break-word }
        .ci-loading { display:flex;align-items:center;justify-content:center;gap:10px;padding:60px 0;color:#9ca3af;font-size:.9rem }
        .ci-error { padding:20px;background:#fef2f2;border:1px solid #fecaca;border-radius:8px;color:#b91c1c;font-size:.85rem;text-align:center }
        @media(max-width:768px) {
            .ci-modal { flex-direction:column;max-width:100%;border-radius:12px 12px 0 0;max-height:95vh }
            .ci-sidebar { width:100%;overflow-x:auto;padding:12px;border-right:none;border-bottom:1px solid #f3f4f6;display:flex;flex-direction:row }
            .ci-nav-item { white-space:nowrap;border-left:none;border-bottom:2px solid transparent;padding:6px 12px }
            .ci-nav-item.active { border-bottom-color:#5106A7;border-left-color:transparent }
            .ci-grid { grid-template-columns:repeat(2,1fr) }
        }
    `;
    document.head.appendChild(style);
})();

var CI_SECTIONS = [
    { id: 'encabezado',      label: 'Encabezado',           icon: 'fa-users',              color: '#6b7280' },
    { id: 'victima',         label: 'Datos de la víctima',  icon: 'fa-user',               color: '#5106A7' },
    { id: 'etnicos',         label: 'Étnicos y sociales',   icon: 'fa-globe',              color: '#10b981' },
    { id: 'contacto',        label: 'Contacto y familia',   icon: 'fa-phone',              color: '#f59e0b' },
    { id: 'ubicacion',       label: 'Ubicación',            icon: 'fa-map-marker-alt',     color: '#8b5cf6' },
    { id: 'hechos',          label: 'Hechos',               icon: 'fa-file-alt',           color: '#ef4444' },
    { id: 'agresor',         label: 'Datos del agresor',    icon: 'fa-user-slash',         color: '#dc2626' },
    { id: 'riesgo',          label: 'Identificación riesgo',icon: 'fa-exclamation-triangle',color: '#f97316' },
];


/* ─── Contenido de secciones (compartido entre modal e inline) ─── */
var CI_SECTIONS_TPL = `
    <div v-if="sec.id==='encabezado'" class="ci-grid">
        <div class="ci-field" v-if="val(info.encabezado.agenteAsignado)"><span class="ci-label">Agente asignado</span><span class="ci-value" style="color:#5106A7;font-weight:600">\${ info.encabezado.agenteAsignado }</span></div>
        <div class="ci-field" v-if="val(info.encabezado.funcionarios)"><span class="ci-label">Historial funcionarios</span><span class="ci-value">\${ info.encabezado.funcionarios }</span></div>
        <div class="ci-field" v-if="val(info.encabezado.fechaCreacion)"><span class="ci-label">Fecha creación</span><span class="ci-value">\${ info.encabezado.fechaCreacion }</span></div>
        <div class="ci-field" v-if="val(info.encabezado.fechaModificacion)"><span class="ci-label">Última modificación</span><span class="ci-value">\${ info.encabezado.fechaModificacion }</span></div>
        <div class="ci-field" v-if="val(info.encabezado.estado)"><span class="ci-label">Estado</span><span class="ci-value">\${ info.encabezado.estado }</span></div>
    </div>
    <div v-if="sec.id==='victima'" class="ci-grid">
        <div class="ci-field" v-if="val(info.victima.nombres)"><span class="ci-label">Nombres</span><span class="ci-value">\${ info.victima.nombres }</span></div>
        <div class="ci-field" v-if="val(info.victima.apellidos)"><span class="ci-label">Apellidos</span><span class="ci-value">\${ info.victima.apellidos }</span></div>
        <div class="ci-field" v-if="val(info.victima.edad)"><span class="ci-label">Edad</span><span class="ci-value">\${ info.victima.edad } años</span></div>
        <div class="ci-field" v-if="val(info.victima.fechaNacimiento)"><span class="ci-label">Fecha nacimiento</span><span class="ci-value">\${ info.victima.fechaNacimiento }</span></div>
        <div class="ci-field" v-if="val(info.victima.tipoDocumento)"><span class="ci-label">Tipo documento</span><span class="ci-value">\${ info.victima.tipoDocumento }</span></div>
        <div class="ci-field" v-if="val(info.victima.numeroDocumento)"><span class="ci-label">Número documento</span><span class="ci-value">\${ info.victima.numeroDocumento }</span></div>
        <div class="ci-field" v-if="val(info.victima.nacionalidad)"><span class="ci-label">Nacionalidad</span><span class="ci-value">\${ info.victima.nacionalidad }</span></div>
        <div class="ci-field" v-if="val(info.victima.municipio)"><span class="ci-label">Municipio</span><span class="ci-value">\${ info.victima.municipio }</span></div>
        <div class="ci-field" v-if="val(info.victima.correoElectronico)"><span class="ci-label">Correo</span><span class="ci-value">\${ info.victima.correoElectronico }</span></div>
        <div class="ci-field" v-if="val(info.datosPersonales.telefono)"><span class="ci-label">Teléfono</span><span class="ci-value">\${ info.datosPersonales.telefono }</span></div>
        <div class="ci-field" v-if="val(info.datosPersonales.genero)"><span class="ci-label">Género</span><span class="ci-value">\${ info.datosPersonales.genero }</span></div>
        <div class="ci-field" v-if="val(info.datosPersonales.identidadGenero)"><span class="ci-label">Identidad de género</span><span class="ci-value">\${ info.datosPersonales.identidadGenero }</span></div>
        <div class="ci-field" v-if="val(info.datosPersonales.orientacionSexual)"><span class="ci-label">Orientación sexual</span><span class="ci-value">\${ info.datosPersonales.orientacionSexual }</span></div>
        <div class="ci-field" v-if="val(info.datosPersonales.condicionMigratoria)"><span class="ci-label">Condición migratoria</span><span class="ci-value">\${ info.datosPersonales.condicionMigratoria }</span></div>
        <div class="ci-field" v-if="val(info.contacto.estadoCivil)"><span class="ci-label">Estado civil</span><span class="ci-value">\${ info.contacto.estadoCivil }</span></div>
        <div class="ci-field" v-if="val(info.victima.direccionResidencia)"><span class="ci-label">Dirección de residencia</span><span class="ci-value">\${ info.victima.direccionResidencia }</span></div>
        <div class="ci-field" v-if="val(info.victima.nombreIdentitario)"><span class="ci-label">Nombre identitario</span><span class="ci-value">\${ info.victima.nombreIdentitario }</span></div>
    </div>
    <div v-if="sec.id==='etnicos'" class="ci-grid">
        <div class="ci-field" v-if="val(info.etnicos.grupoEtnico)"><span class="ci-label">Grupo étnico</span><span class="ci-value">\${ info.etnicos.grupoEtnico }</span></div>
        <div class="ci-field" v-if="val(info.etnicos.afrodescendiente)"><span class="ci-label">Afrodescendiente</span><span class="ci-value">\${ siNo(info.etnicos.afrodescendiente) }</span></div>
        <div class="ci-field" v-if="val(info.etnicos.indigena)"><span class="ci-label">Indígena</span><span class="ci-value">\${ info.etnicos.indigena }</span></div>
        <div class="ci-field" v-if="val(info.etnicos.campesino)"><span class="ci-label">Campesino</span><span class="ci-value">\${ info.etnicos.campesino }</span></div>
        <div class="ci-field" v-if="val(info.etnicos.victimaConflicto)"><span class="ci-label">Víctima conflicto</span><span class="ci-value">\${ siNo(info.etnicos.victimaConflicto) }</span></div>
        <div class="ci-field" v-if="val(info.contacto.discapacidad)"><span class="ci-label">Discapacidad</span><span class="ci-value">\${ info.contacto.discapacidad }</span></div>
    </div>
    <div v-if="sec.id==='contacto'" class="ci-grid">
        <div class="ci-field"><span class="ci-label">Nombre del contacto</span><span class="ci-value">\${ info.contacto.nombreContacto || 'No registra' }</span></div>
        <div class="ci-field"><span class="ci-label">Teléfono del contacto</span><span class="ci-value">\${ info.contacto.telefonoContacto || 'No registra' }</span></div>
        <div class="ci-field"><span class="ci-label">Parentesco</span><span class="ci-value">\${ info.contacto.parentesco || 'No registra' }</span></div>
        <div class="ci-field"><span class="ci-label">Personas a cargo</span><span class="ci-value">\${ info.contacto.personasCargo || 'No registra' }</span></div>
        <div class="ci-field" v-if="val(info.contacto.numeroHijos)"><span class="ci-label">Número de hijos</span><span class="ci-value">\${ info.contacto.numeroHijos }</span></div>
    </div>
    <div v-if="sec.id==='ubicacion'" class="ci-grid">
        <div class="ci-field" v-if="val(info.ubicacion.departamento)"><span class="ci-label">Departamento</span><span class="ci-value">\${ info.ubicacion.departamento }</span></div>
        <div class="ci-field" v-if="val(info.ubicacion.ciudad)"><span class="ci-label">Ciudad</span><span class="ci-value">\${ info.ubicacion.ciudad }</span></div>
        <div class="ci-field" v-if="val(info.ubicacion.municipio)"><span class="ci-label">Municipio</span><span class="ci-value">\${ info.ubicacion.municipio }</span></div>
    </div>
    <div v-if="sec.id==='hechos'" class="ci-grid">
        <div class="ci-field" v-if="val(info.hechos.descripcion)" style="grid-column:1/-1">
            <span class="ci-label">Descripción</span>
            <span class="ci-value" :style="{ display:'-webkit-box', '-webkit-line-clamp': hechosExpandido ? 'unset' : '4', '-webkit-box-orient':'vertical', overflow: hechosExpandido ? 'visible' : 'hidden' }">\${ info.hechos.descripcion }</span>
            <button v-if="info.hechos.descripcion && info.hechos.descripcion.length > 200" @click="hechosExpandido = !hechosExpandido" style="background:none;border:none;color:#5106A7;font-size:.75rem;cursor:pointer;padding:4px 0;font-weight:600">\${ hechosExpandido ? '▲ Ver menos' : '▼ Ver más' }</button>
        </div>
        <div class="ci-field" v-if="val(info.hechos.fechaHechos)"><span class="ci-label">Fecha hechos</span><span class="ci-value">\${ info.hechos.fechaHechos }</span></div>
        <div class="ci-field" v-if="val(info.hechos.horario)"><span class="ci-label">Horario</span><span class="ci-value">\${ info.hechos.horario }</span></div>
        <div class="ci-field" v-if="val(info.hechos.escenarioViolencia)"><span class="ci-label">Escenario</span><span class="ci-value">\${ info.hechos.escenarioViolencia }</span></div>
        <div class="ci-field" v-if="val(info.hechos.direccionHechos)"><span class="ci-label">Dirección de los hechos</span><span class="ci-value">\${ info.hechos.direccionHechos }</span></div>
        <div class="ci-field" v-if="val(info.hechos.riesgoFeminicida)"><span class="ci-label">Riesgo feminicida</span><span class="ci-value">\${ siNo(info.hechos.riesgoFeminicida) }</span></div>
    </div>
    <div v-if="sec.id==='agresor'" class="ci-grid">
        <div class="ci-field" v-if="val(info.agresor.tipoAgresor)"><span class="ci-label">Tipo agresor</span><span class="ci-value">\${ info.agresor.tipoAgresor }</span></div>
        <div class="ci-field" v-if="val(info.agresor.relacion)"><span class="ci-label">Relación</span><span class="ci-value">\${ info.agresor.relacion }</span></div>
        <div class="ci-field" v-if="val(info.agresor.nombre)"><span class="ci-label">Nombre</span><span class="ci-value">\${ info.agresor.nombre }</span></div>
        <div class="ci-field" v-if="val(info.agresor.tipoDocumento)"><span class="ci-label">Tipo documento</span><span class="ci-value">\${ info.agresor.tipoDocumento }</span></div>
        <div class="ci-field" v-if="val(info.agresor.numeroDocumento)"><span class="ci-label">Documento</span><span class="ci-value">\${ info.agresor.numeroDocumento }</span></div>
        <div class="ci-field" v-if="val(info.agresor.direccion)"><span class="ci-label">Dirección</span><span class="ci-value">\${ info.agresor.direccion }</span></div>
        <div class="ci-field" v-if="val(info.agresor.telefono)"><span class="ci-label">Teléfono</span><span class="ci-value">\${ info.agresor.telefono }</span></div>
        <div class="ci-field" v-if="val(info.agresor.numAgresores)"><span class="ci-label">Número de agresores</span><span class="ci-value">\${ info.agresor.numAgresores }</span></div>
        <div class="ci-field" v-if="val(info.agresor.proximidad)"><span class="ci-label">Proximidad del agresor</span><span class="ci-value">\${ info.agresor.proximidad }</span></div>
        <div class="ci-field" v-if="val(info.agresor.generoAgresor)"><span class="ci-label">Género del agresor</span><span class="ci-value">\${ info.agresor.generoAgresor }</span></div>
    </div>
    <div v-if="sec.id==='riesgo'" class="ci-grid">
        <div class="ci-field" v-if="val(info.riesgo.nivelRiesgoTexto)"><span class="ci-label">Nivel de riesgo</span><span class="ci-value">\${ info.riesgo.nivelRiesgoTexto }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.amenazasMuerte)"><span class="ci-label">Amenazas de muerte</span><span class="ci-value">\${ siNo(info.riesgo.amenazasMuerte) }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.agresorTieneArmas)"><span class="ci-label">Agresor tiene armas</span><span class="ci-value">\${ siNo(info.riesgo.agresorTieneArmas) }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.violenciaPrevia)"><span class="ci-label">Violencia previa</span><span class="ci-value">\${ siNo(info.riesgo.violenciaPrevia) }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.celosoViolento)"><span class="ci-label">Celoso y violento</span><span class="ci-value">\${ siNo(info.riesgo.celosoViolento) }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.creeCapazMatar)"><span class="ci-label">Cree capaz de matar</span><span class="ci-value">\${ siNo(info.riesgo.creeCapazMatar) }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.riesgoInminente)"><span class="ci-label">Riesgo inminente</span><span class="ci-value">\${ siNo(info.riesgo.riesgoInminente) }</span></div>
        <div class="ci-field" v-if="val(info.riesgo.denunciaPrevia)"><span class="ci-label">Denuncia previa</span><span class="ci-value">\${ siNo(info.riesgo.denunciaPrevia) }</span></div>
    </div>
`;


/* ─── Bloque de secciones reutilizable ─── */
var CI_SECTIONS_BLOCK = `
    <div v-for="sec in sections" :key="sec.id" class="ci-section" :ref="'section-' + sec.id">
        <div :class="['ci-section-header', isOpen(sec.id) ? 'open' : '']" @click="toggle(sec.id)">
            <div class="ci-section-icon" :style="{ background: sec.color }"><i :class="'fa ' + sec.icon"></i></div>
            <span class="ci-section-label">\${ sec.label }</span>
            <span class="ci-section-count">\${ fieldCount(sec.id) } campos</span>
            <i class="fa fa-chevron-down ci-section-chevron"></i>
        </div>
        <div v-if="isOpen(sec.id)" class="ci-section-body">` + CI_SECTIONS_TPL + `</div>
    </div>
`;

app.component('case-info', {
    delimiters: ['${', '}'],
    props: {
        caseId:  { type: String, required: true },
        visible: { type: Boolean, required: true },
        mode:    { type: String, default: 'modal' },
    },
    emits: ['close'],
    data: function() {
        return {
            loading: false,
            error: null,
            info: null,
            sections: CI_SECTIONS,
            openSections: { encabezado: true, victima: true, etnicos: true, contacto: true, ubicacion: true, hechos: true, agresor: true, riesgo: true },
            activeNav: 'encabezado',
            hechosExpandido: false,
        };
    },
    watch: {
        visible: function(val) {
            if (val && !this.info) this.cargar();
        }
    },
    mounted: function() {
        // En modo inline, cargar inmediatamente si es visible
        if (this.mode === 'inline' && this.visible) this.cargar();
    },
    methods: {
        async cargar() {
            this.loading = true;
            this.error = null;
            try {
                var res = await fetch('/api/v1/casos/' + this.caseId + '/info-completa');
                if (!res.ok) throw new Error('Error ' + res.status);
                this.info = await res.json();
            } catch (e) {
                this.error = e.message;
            } finally {
                this.loading = false;
            }
        },
        cerrar() { this.$emit('close'); },
        toggle(id) { this.openSections[id] = !this.openSections[id]; this.activeNav = id; },
        navTo(id) {
            var newState = {};
            for (var i = 0; i < this.sections.length; i++) newState[this.sections[i].id] = false;
            newState[id] = true;
            this.openSections = newState;
            this.activeNav = id;
            var el = this.$refs['section-' + id];
            if (el && el[0]) el[0].scrollIntoView({ behavior: 'smooth', block: 'start' });
        },
        isOpen(id) { return !!this.openSections[id]; },
        val(v) { return (v === null || v === undefined || v === '' || v === '0') ? null : v; },
        siNo(v) {
            if (v === '1' || v === 'true' || v === 'Si' || v === 'si') return 'Sí';
            if (v === '0' || v === 'false' || v === 'No' || v === 'no') return 'No';
            return v || null;
        },
        fieldCount(sectionId) {
            if (!this.info || !this.info[sectionId]) return 0;
            var obj = this.info[sectionId]; var count = 0;
            for (var k in obj) { if (this.val(obj[k])) count++; }
            return count;
        },
    },
    template: `
<div v-if="visible">
    <!-- MODO MODAL -->
    <div v-if="mode === 'modal'" class="ci-backdrop" @click.self="cerrar">
        <div class="ci-modal">
            <div class="ci-sidebar">
                <p class="ci-sidebar-title">Secciones</p>
                <div v-for="sec in sections" :key="sec.id" :class="['ci-nav-item', activeNav === sec.id ? 'active' : '']" @click="navTo(sec.id)">
                    <i :class="'fa ' + sec.icon"></i> \${ sec.label }
                </div>
            </div>
            <div class="ci-main">
                <div class="ci-header">
                    <p class="ci-title"><i class="fa fa-folder-open" style="color:#5106A7"></i> Información del caso</p>
                    <button class="ci-close" @click="cerrar">✕</button>
                </div>
                <div class="ci-body">
                    <div v-if="loading" class="ci-loading"><i class="fa fa-spinner fa-spin"></i> Cargando...</div>
                    <div v-else-if="error" class="ci-error">\${ error }</div>
                    <template v-else-if="info">` + CI_SECTIONS_BLOCK + `</template>
                </div>
            </div>
        </div>
    </div>
    <!-- MODO INLINE -->
    <div v-else>
        <div v-if="loading" class="ci-loading"><i class="fa fa-spinner fa-spin"></i> Cargando...</div>
        <div v-else-if="error" class="ci-error">\${ error }</div>
        <template v-else-if="info">` + CI_SECTIONS_BLOCK + `</template>
    </div>
</div>
    `
});
