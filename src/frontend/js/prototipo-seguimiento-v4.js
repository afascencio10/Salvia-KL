/**
 * Prototipo del formulario de seguimiento V.4.
 * Renderiza y captura respuestas en memoria. No llama al servidor.
 * Los estilos visuales los inyecta dinamic-form.js (clases df-*).
 */
(function injectProtoStyles() {
    if (document.getElementById('proto-v4-styles')) return;
    const style = document.createElement('style');
    style.id = 'proto-v4-styles';
    style.textContent = `
        .proto-banner {
            background: #fffbeb;
            border: 1px solid #fde68a;
            color: #92400e;
            border-radius: 12px;
            padding: 12px 16px;
            font-size: 14px;
            line-height: 1.45;
            margin-bottom: 16px;
        }
        .proto-context {
            background: #fff;
            border: 1px solid #e5e7eb;
            border-radius: 12px;
            padding: 16px 20px;
            margin-bottom: 16px;
            box-shadow: 0 1px 2px 0 rgb(0 0 0/.05);
        }
        .proto-context h3 {
            margin: 0 0 4px;
            font-size: 13px;
            font-weight: 700;
            color: #6d28d9;
            text-transform: uppercase;
            letter-spacing: .04em;
        }
        .proto-context p { margin: 0 0 12px; font-size: 13px; color: #6b7280; }
        .proto-grid { display: flex; flex-wrap: wrap; gap: 14px 28px; align-items: center; }
        .proto-grid label { display: flex; align-items: center; gap: 8px; font-size: 14px; color: #374151; margin: 0; }
        .proto-grid select, .proto-grid input[type="text"] {
            border: 1px solid #e5e7eb; border-radius: 8px; padding: 6px 10px; font-size: 14px;
        }
        .proto-heading {
            font-size: 13px; font-weight: 700; color: #6d28d9;
            text-transform: uppercase; letter-spacing: .04em; margin-top: 8px;
        }
        .proto-summary-sec { margin-bottom: 16px; }
        .proto-summary-sec h3 { font-size: 16px; margin: 0 0 8px; color: #111827; }
        .proto-row { padding: 8px 0; border-bottom: 1px solid #f3f4f6; }
        .proto-row-label { font-size: 13px; color: #6b7280; }
        .proto-row-value { font-size: 14px; color: #111827; font-weight: 500; }
        .proto-head { font-size: 13px; font-weight: 700; color: #6d28d9; margin-top: 10px; }
        .df-repeater .df-repeater { margin-top: 4px; }
        .proto-barrera-card {
            background: #f5f3ff;
            border: 1px solid #ddd6fe;
            border-radius: 12px;
            padding: 14px 16px 4px;
            margin: 8px 0 12px;
        }
        .proto-barrera-title {
            margin: 0 0 12px;
            font-size: 15px;
            font-weight: 700;
            color: #5b21b6;
        }
        .proto-panel {
            background: #fffbeb;
            border: 1px solid #fde68a;
            border-radius: 12px;
            padding: 14px 16px 4px;
            margin: 12px 0;
        }
        .proto-panel-title {
            margin: 0 0 12px;
            font-size: 13px;
            font-weight: 700;
            color: #92400e;
            text-transform: uppercase;
            letter-spacing: .04em;
        }
        @media (max-width: 800px) {
            #proto-app .df-wrapper { flex-direction: column; }
            #proto-app .df-sidebar { width: 100%; }
        }
    `;
    document.head.appendChild(style);
})();

var YN = [{ v: 'si', l: 'Sí' }, { v: 'no', l: 'No' }];
var SCALE = [1, 2, 3, 4, 5].map(function (n) { return { v: String(n), l: String(n) }; });

var PROTECTORES = [
    ['medidas_vigentes', 'Medidas de protección vigentes, monitoreadas, efectivas y adecuadas a las condiciones diferenciales de la persona'],
    ['lugar_seguro', 'Permanencia en un lugar seguro o traslado reciente a un entorno protegido y respetuoso de las condiciones diferenciales de la persona'],
    ['separacion', 'Separación física y sostenida de la persona agresora'],
    ['redes', 'Redes de apoyo activas, seguras y respetuosas de las condiciones diferenciales de la persona'],
    ['psico', 'Acompañamiento psicosocial o psicoterapéutico'],
    ['autonomia', 'Cuenta con autonomía económica'],
    ['efecto_protector', 'Acciones institucionales con efecto protector tangible en la percepción de seguridad de la persona'],
    ['ninguna', 'Ninguna'],
];
var RIESGOS = [
    ['habitacional', 'Condiciones habitacionales inseguras o inestables (incluye mudanzas forzadas o frecuentes asociadas a la violencia)'],
    ['grupo_armado', 'Agresor pertenece a un grupo organizado al margen de la ley'],
    ['spa', 'Agresor con consumo problemático de sustancias psicoactivas ilegales y legales'],
    ['miedo_denunciar', 'Miedo a denunciar o separarse por represalias'],
    ['bajo_reconocimiento', 'Bajo reconocimiento del riesgo que implica ausencia de implementación de esquema de autoprotección'],
    ['embarazo', 'Embarazo con riesgo o gestación vulnerable'],
    ['discapacidad', 'Discapacidad física, cognitiva y/o psicosocial'],
    ['rural', 'Residencia rural aislada o instituciones distantes'],
    ['poder', 'Persona agresora con poder institucional o territorial'],
    ['fuerza_publica', 'Agresor pertenece a la fuerza pública'],
    ['dependencia_osigd', 'Dependencia respecto de personas o entornos que rechazan, sancionan o buscan ocultar la OSIGD'],
    ['aislamiento_osigd', 'Aislamiento o debilitamiento de redes por rechazo relacionado con la OSIGD'],
    ['territorio_osigd', 'Permanencia en un territorio o entorno que controla, sanciona o reprime las OSIGD'],
    ['ninguna', 'Ninguna'],
];
var EXTREMOS = [
    ['amenazas', 'Amenazas de muerte explícitas'],
    ['escalamiento', 'Escalamiento en frecuencia o gravedad de la violencia'],
    ['sexual', 'Contacto sexual no consentido (último mes)'],
    ['aislamiento', 'Aislamiento, retención o privación de la libertad'],
    ['armas', 'Acceso reciente del agresor a armas'],
    ['bienes', 'Destrucción o control de bienes'],
    ['vigilancia', 'Aumento de vigilancia, control o monitoreo'],
    ['tentativa', 'Tentativa de feminicidio u homicidio abyecto o fútil'],
    ['ninguna', 'Ninguna'],
];
var SECTORES = [
    ['salud', 'Salud'],
    ['justicia', 'Justicia'],
    ['proteccion', 'Protección'],
    ['educacion', 'Educación'],
    ['ministerio', 'Ministerio Público'],
    ['territorial', 'Entidad territorial'],
    ['osc', 'Organización de la sociedad civil, comunitaria o internacional de Naciones Unidas'],
    ['clave', 'Entidades clave con acciones complementarias y/o especializadas frente a VBG/VPP'],
    ['otra', 'Otra'],
];
var ENTIDADES = {
    salud: ['Hospital Público', 'Hospital Privado', 'IPS', 'EPS', 'Puesto de Salud', 'Consultorio médico'],
    justicia: ['Fiscalía General de la Nación', 'Comisaría de Familia', 'CAI — Policía', 'Instituto Nacional de Medicina Legal', 'Inspección de Policía', 'Juzgado'],
    proteccion: ['Comisaría de Familia', 'ICBF', 'Unidad Nacional de Protección', 'Defensoría de Familia', 'Inspección de Policía'],
    educacion: ['Institución educativa', 'Secretaría de Educación'],
    ministerio: ['Defensoría del Pueblo', 'Procuraduría', 'Personería'],
    territorial: ['Alcaldía', 'Gobernación', 'Secretaría de la Mujer'],
    osc: ['Organización de mujeres', 'Agencia de Naciones Unidas'],
    clave: ['Casa de la Mujer', 'Casa LGBTI'],
    otra: [],
};
var ACCIONES = {
    salud: ['Atención en urgencias', 'Activación de protocolo de violencia sexual', 'Valoración medicolegal', 'Continuidad del tratamiento'],
    justicia: ['Recepción de denuncia', 'Asignación de radicado', 'Medidas de protección', 'Notificación a la persona agresora'],
    proteccion: ['Medida de protección', 'Valoración del riesgo', 'Acompañamiento policial', 'Seguimiento a la medida'],
    educacion: ['Orientación a la comunidad educativa', 'Ruta de activación escolar', 'Acompañamiento a la víctima'],
    ministerio: ['Orientación en derechos', 'Acompañamiento al proceso', 'Seguimiento institucional'],
    territorial: ['Atención en secretaría de la mujer', 'Oferta de estabilización', 'Medidas de atención integral'],
    osc: ['Acompañamiento psicosocial', 'Orientación jurídica', 'Albergue o apoyo humanitario'],
    clave: ['Atención especializada en VBG', 'Atención a personas OSIGD', 'Acompañamiento jurídico'],
    otra: ['Orientación', 'Acompañamiento', 'Otra gestión'],
};
var TERR = [
    { v: 'bogota', l: 'Bogotá D.C.', munis: [
        { v: 'usaquen', l: 'Usaquén', ciudades: ['Usaquén', 'Zona rural'] },
        { v: 'kennedy', l: 'Kennedy', ciudades: ['Kennedy', 'Zona rural'] },
        { v: 'suba', l: 'Suba', ciudades: ['Suba', 'Zona rural'] },
    ]},
    { v: 'antioquia', l: 'Antioquia', munis: [
        { v: 'medellin', l: 'Medellín', ciudades: ['Medellín', 'Zona rural'] },
        { v: 'bello', l: 'Bello', ciudades: ['Bello', 'Zona rural'] },
    ]},
    { v: 'valle', l: 'Valle del Cauca', munis: [
        { v: 'cali', l: 'Cali', ciudades: ['Cali', 'Zona rural'] },
    ]},
];
var CANALES = [
    ['notificacion', 'Notificación SALVIA'],
    ['articulacion', 'Articulación institucional telefónica o por mensajería'],
    ['denuncia', 'Denuncia fácil'],
    ['misalvia', 'Mi Salvia'],
    ['ninguna', 'Ninguna'],
];
var SALUD_B = [
    ['aseguramiento', 'Dificultades de aseguramiento o afiliación en salud'],
    ['urgencias', 'Negaciones en los servicios de urgencias'],
    ['calidad', 'Fallas en la calidad del servicio de urgencias'],
    ['citas', 'Demoras en la asignación de citas y continuidad de tratamientos'],
    ['protocolo_sexual', 'Faltas de activación de protocolos para violencia sexual'],
    ['quimicos', 'Faltas de activación de protocolos para ataques con agentes químicos'],
    ['medicolegal', 'Demoras y dificultades en la valoración medicolegal'],
    ['ive', 'Negaciones en el acceso a la Interrupción Voluntaria del Embarazo (IVE)'],
    ['discapacidad', 'Exigencias indebidas de autorización para procedimientos en personas con discapacidad cognitiva y/o psicosocial'],
    ['cosmovision', 'Atenciones en salud no centradas en cosmovisiones y prácticas culturales propias'],
    ['osigd', 'Negaciones o imposiciones de procedimientos a personas con diversidad sexual o identidad de género diversa'],
    ['diferenciales', 'Negación, restricción o condicionamiento de servicios por falta de adecuación a las necesidades diferenciales de la persona'],
];
var JUSTICIA_B = [
    ['negativa', 'Negativa institucional para recibir la denuncia'],
    ['no_confrontacion', 'Denegación al derecho a la no confrontación'],
    ['notificacion', 'Deficiencias en la notificación y recolección de datos del agresor'],
    ['documentacion', 'Falta de entrega de información y documentación procesal'],
    ['cargas', 'Imposición de cargas probatorias injustificadas'],
    ['celeridad', 'Falta de celeridad en la investigación judicial'],
    ['tipificacion', 'Tipificación errónea del delito'],
    ['traslado', 'Traslado inadecuado o incompleto del caso / expediente judicial'],
    ['representacion', 'Ausencia o dificultad de acceso a representación judicial'],
    ['terminos', 'Riesgo de vencimiento de términos'],
    ['libertad', 'Medidas preventivas de libertad que ponen en riesgo la vida de la mujer y su núcleo familiar'],
    ['confidencial', 'Divulgación de información confidencial del proceso'],
    ['prejuicio', 'Omisión o minimización del posible móvil de prejuicio en la investigación'],
];
var PROTECCION_B = [
    ['demora_mp', 'Ausencia o demoras en la adopción de medidas de protección urgentes o definitivas en Comisaría de Familia'],
    ['fisicas', 'Ausencia o dilación en adopción de medidas de protección físicas'],
    ['integral', 'Demora o ausencia en la activación de medidas de atención integral'],
    ['incumplimiento', 'Incumplimiento de medidas de protección sin respuesta institucional'],
    ['ineficaces', 'Medidas de protección ineficaces o mal implementadas'],
    ['valoracion', 'Fallas en la valoración y actualización del riesgo'],
    ['visitas', 'Omisión de restricción de visitas ante riesgo de feminicidio'],
    ['policial', 'Omisión de apoyo policial en residencia o lugar de trabajo'],
];
var INST = {
    salud: ['Hospital Público', 'Hospital Privado', 'IPS', 'EPS', 'Puesto de Salud', 'Consultorio médico', 'No sabe / No responde'],
    justicia: ['CAI — Policía', 'CAIVAS', 'Casas de Justicia', 'CAVIV', 'Comisaría de Familia', 'Estación de Policía', 'Fiscalía General de la Nación', 'Inspección de Policía', 'Instituto Nacional de Medicina Legal', 'Jueces Civiles o Promiscuos Municipales', 'Jueces de Control de Garantías o Jueces Penales', 'Policía Judicial', 'Unidad de reacción inmediata', 'No sabe / No responde'],
    proteccion: ['Comisarías de familia', 'Fiscalía (funciones de protección)', 'ICBF', 'Unidad Nacional de Protección', 'Jueces Civiles o Promiscuos Municipales', 'Inspección de Policía', 'Defensoría de Familia', 'Ejército Nacional', 'Armada Colombiana', 'No sabe / No responde'],
};
var TALENTO = [
    ['enfoque', 'Falta de aplicación del enfoque o perspectiva de género / revictimización o violencia institucional'],
    ['protocolos', 'Falta de aplicación o desactualización de protocolos de atención en VBG'],
    ['capacidad', 'Capacidad institucional limitada'],
    ['desarticulacion', 'Desarticulación institucional y ausencia de gestión integral de casos'],
    ['orientacion', 'Insuficiencia en la orientación en derechos a las víctimas'],
];
var ECO = [
    ['autonomia', 'Falta de autonomía económica'],
    ['pobreza', 'Pobreza y desigualdad estructural'],
];
var TERR_B = [
    ['ausencia', 'Ausencia institucional, infraestructura y conectividad insuficiente o instalaciones distantes'],
];
var DIFERENCIAL = [
    ['osigd', 'OSIGD'],
    ['salud_mental', 'Salud mental'],
    ['etnico', 'Origen étnico (indígena, afrodescendiente, raizal, palenquera)'],
    ['ciclo', 'Ciclo de vida (niñez, adolescencia, adultez mayor)'],
    ['discapacidad', 'Discapacidad física o cognitiva'],
    ['migrantes', 'Personas migrantes y refugiadas'],
    ['libertad', 'Personas privadas de la libertad'],
    ['calle', 'Situación de calle'],
    ['trata', 'Víctimas de trata'],
    ['sexuales', 'Actividades sexuales pagas'],
    ['lideresas', 'Lideresas y defensoras de derechos humanos'],
    ['conflicto', 'Víctimas del conflicto armado'],
];
var EQUIPOS = [
    ['medidas', 'Medidas de Emergencia'],
    ['psico', 'Atención psicosocial'],
    ['osignh', 'Psicosocial Personas OSIGNH'],
    ['hombres', 'Atención Hombres'],
    ['discapacidad', 'Equipo Discapacidad'],
    ['estabilizacion', 'Estabilización'],
];
var MEDIDAS_TIPO = [
    ['alojamiento', 'Alojamiento'],
    ['transporte', 'Transporte'],
    ['alimentacion', 'Alimentación'],
    ['vestuario', 'Vestuario'],
    ['apoyo_psico', 'Apoyo psicosocial'],
    ['otras', 'Otras medidas de emergencia'],
];
var PSICO_CRIT = [
    ['poblacion', 'Mujeres en todas sus diversidades, mayores de 18 años, víctimas de VBG o VpP, con voluntariedad y en post-emergencia (no riesgo inminente, no desborde emocional)'],
    ['suicida', 'Conducta suicida derivada de VBG o VpP — 3 puntos'],
    ['temor', 'Temor frente a la activación de la ruta — 1 punto'],
    ['inter', 'Interseccionalidad (gestación, étnico, recursos limitados, etc.) — 2 puntos'],
    ['territorio', 'Condiciones territoriales y de orden público — 1 punto'],
    ['sin_acceso', 'Sin acceso a la atención psicosocial o psicológica — 1 punto'],
    ['naturalizacion', 'Naturalización de la VBG y roles de género — 1 punto'],
];
var PSICO_PTS = { poblacion: 0, suicida: 3, temor: 1, inter: 2, territorio: 1, sin_acceso: 1, naturalizacion: 1 };
var ELEMENTOS = [
    ['escucha', 'Escucha activa y contención emocional en crisis'],
    ['psicoeducacion', 'Psicoeducación del ciclo de la violencia, telaraña del abuso y/o impactos de la violencia'],
    ['necesidades', 'Apoyo en identificación de necesidades específicas y riesgo'],
    ['despatologizacion', 'Despatologización de experiencias y reacciones asociadas a la violencia'],
    ['derechos', 'Reconocimiento y comprensión de derechos de víctimas'],
    ['estado', 'Orientación sobre acciones disponibles y responsabilidad del Estado'],
    ['salud_mental', 'Sensibilización y orientación sobre servicios de salud mental'],
    ['plan_seguridad', 'Co-construcción del Plan de Seguridad'],
    ['senales', 'Identificación de señales tempranas de riesgo y estrategias de respuesta'],
    ['corresponsabilidad', 'Fortalecimiento de la corresponsabilidad de la persona en su seguridad'],
    ['acceso', 'Orientación para acceso a servicios, programas o apoyos institucionales'],
    ['redes', 'Identificación de redes seguras y afirmativas'],
    ['recursos', 'Identificación y fortalecimiento de recursos personales'],
    ['autonomia', 'Fortalecimiento de la autonomía económica y vinculación a proyecto de vida'],
    ['mandatos', 'Cuestionamiento de mandatos culturales que sostienen la violencia'],
    ['seguimiento_ruta', 'Seguimiento a la ruta, actualización de riesgo y revisión de barreras'],
    ['decisiones', 'Acompañamiento en la toma de decisiones informadas respetando la autonomía'],
    ['autocuidado', 'Promoción de autocuidado y estrategias de protección digital y física'],
    ['creencias', 'Promoción de reflexión crítica sobre creencias que justifican control y dependencia'],
    ['narrativas', 'Identificación de creencias que sostienen la violencia y proponer narrativas alternativas'],
];
var DEMO_ENTIDADES = [
    {
        id: 'comisaria',
        nombre: 'Comisaría de Familia — Kennedy',
        barreras: [
            { id: 'b1', nombre: 'Demora en medidas de protección' },
            { id: 'b2', nombre: 'Falta de seguimiento a la medida' },
        ],
        pendientes: ['Adopción de medida de protección', 'Seguimiento a la medida'],
    },
    {
        id: 'hospital',
        nombre: 'Hospital Público — Kennedy',
        barreras: [],
        pendientes: ['Atención en urgencias'],
    },
];
var DEMO_MEDIDAS = [
    { id: 'alojamiento', nombre: 'Alojamiento temporal' },
    { id: 'transporte', nombre: 'Transporte' },
];
var CIERRE_MOTIVOS = [
    ['no_desea', 'La persona no desea que Salvia gestione la barrera'],
    ['no_gestionable', 'Barrera no gestionable desde la competencia institucional'],
    ['mal_clasificada', 'Se clasificó incorrectamente la barrera'],
    ['superada', 'La barrera fue superada de manera efectiva'],
];
var CONCLUSION_MOTIVOS = [
    ['contacto', 'Pérdida de contacto'],
    ['solicitud', 'Solicitud expresa de la persona para finalizar el proceso'],
    ['plan', 'Cumplimiento del plan de acción Salvia'],
    ['corresponsabilidad', 'No corresponsabilidad por parte de la persona'],
    ['no_clasifica', 'La información disponible no permite clasificar el caso como VBG o VpP'],
    ['otro', 'Otro'],
];

function opts(pairs) {
    return pairs.map(function (p) { return { v: p[0], l: p[1] }; });
}
function optsNames(names) {
    return names.map(function (n) {
        var v = n.normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_|_$/g, '');
        return { v: v, l: n };
    });
}
function uid() {
    return Math.random().toString(36).slice(2, 10);
}
function getPath(obj, path) {
    return path.split('.').reduce(function (cur, key) {
        return cur == null ? undefined : cur[key];
    }, obj);
}
function setPath(obj, path, value) {
    var parts = path.split('.');
    var cur = obj;
    for (var i = 0; i < parts.length - 1; i++) {
        if (cur[parts[i]] == null) cur[parts[i]] = /^\d+$/.test(parts[i + 1]) ? [] : {};
        cur = cur[parts[i]];
    }
    cur[parts[parts.length - 1]] = value;
}
function emptyEntidad() {
    return {
        uid: uid(), sector: '', depto: '', municipio: '', ciudad: '', entidad: '', otraEntidad: '',
        radicado: '', fecha: '', completadas: [], pendientes: [], otraRuta: '', canales: [],
        acuerdo: '', debidaDiligencia: '', gestionesRuta: '', nuevaBarrera: '', barrerasNuevas: [],
    };
}
function emptyBarreraEntidad() {
    return {
        uid: uid(), sector: '', tipos: [], otraTipo: '', institucion: '', otraInst: '',
        talento: [], eco: [], terr: [], diferencial: [],
        fecha: '', funcionario: '', depto: '', municipio: '', vereda: '',
        descripcion: '', gestion: [], gestionTexto: '', perspectiva: '',
    };
}
function emptySeg() {
    return {
        completadas: [], pendientes: [], otraPendiente: '', actualizacion: '',
        nuevaBarrera: '', barrerasNuevas: [], acuerdo: '', debidaDiligencia: '',
        canales: [], gestionesRuta: '',
        barreras: {},
    };
}
function humanField(field) {
    var v = field.value;
    if (field.type === 'multi') {
        if (!v || !v.length) return '';
        return v.map(function (x) {
            var o = (field.options || []).find(function (op) { return op.v === x; });
            return o ? o.l : x;
        }).join(' · ');
    }
    if (v == null || String(v).trim() === '') return '';
    var found = (field.options || []).find(function (op) { return op.v === v; });
    if (found) return found.l;
    if (v === 'si') return 'Sí';
    if (v === 'no') return 'No';
    return String(v);
}

var FieldView = {
    props: { field: { type: Object, required: true } },
    emits: ['set'],
    methods: {
        isOn: function (v) {
            if (this.field.type === 'multi') return (this.field.value || []).indexOf(v) !== -1;
            return this.field.value === v;
        },
        blocked: function (v) {
            return (this.field.disabledValues || []).indexOf(v) !== -1;
        },
        choose: function (v) {
            if (this.blocked(v)) return;
            if (this.field.type === 'multi') {
                var cur = (this.field.value || []).slice();
                var ex = this.field.exclusive;
                var next;
                if (cur.indexOf(v) !== -1) next = cur.filter(function (x) { return x !== v; });
                else if (ex && v === ex) next = [ex];
                else next = cur.concat(v).filter(function (x) { return x !== ex; });
                this.$emit('set', this.field.path, next);
                return;
            }
            this.$emit('set', this.field.path, this.field.value === v ? '' : v);
        },
        onInput: function (e) {
            this.$emit('set', this.field.path, e.target.value);
        },
    },
    template: `
        <div v-if="field.type==='single' || field.type==='bool'" :class="field.type==='bool' ? 'df-bool-group' : 'df-btn-group'">
            <button type="button" v-for="o in field.options" :key="o.v"
                :class="[field.type==='bool' ? 'df-bool-btn' : 'df-opt-btn', { 'df-selected': isOn(o.v) }]"
                @click="choose(o.v)">{{ o.l }}</button>
        </div>
        <div v-else-if="field.type==='multi'" class="df-btn-group">
            <button type="button" v-for="o in field.options" :key="o.v"
                class="df-opt-btn" :class="{ 'df-selected': isOn(o.v) }"
                :disabled="blocked(o.v)" @click="choose(o.v)">{{ o.l }}</button>
        </div>
        <textarea v-else-if="field.type==='text'" class="df-textarea" :value="field.value || ''" @input="onInput"></textarea>
        <input v-else-if="field.type==='input' || field.type==='date'" class="df-input"
            :type="field.type==='date' ? 'date' : 'text'" :value="field.value || ''" @input="onInput">
        <select v-else-if="field.type==='select'" class="df-select" :value="field.value || ''" @change="onInput">
            <option value="">Seleccione</option>
            <option v-for="o in field.options" :key="o.v" :value="o.v">{{ o.l }}</option>
        </select>
    `,
};

var SeguimientoV4App = {
    components: { FieldView: FieldView },
    data: function () {
        return {
            ctx: { hasEntidades: true, hasMedidas: true, momento: '1', extemporaneo: false },
            a: {
                risk: {
                    respondiente: '', hechos: '', hechosTexto: '', hechosFecha: '',
                    protectores: [], riesgos: [], extremos: [], descripcion: '',
                },
                med: { '1': {}, '2': {}, '3': {} },
                ext: { profesional: '', justificacion: '', fecha: '' },
                entidadesSel: [],
                seg: {},
                ident: { conocimiento: '', entidades: [] },
                plan: {
                    remision: '', equipos: [], medidasTipos: [], psico: [], osignh: [],
                    discapacidad: [], estabilizacion: [], elementos: [], acciones: '', ajustes: '',
                },
                con: { concluye: '', motivo: '', otro: '', descripcion: '', accion: '', entidad: '', otraEntidad: '', canales: [] },
            },
            currentSectionId: 'riesgo',
            passed: {},
            errors: {},
            popup: false,
            popupText: '',
            done: false,
        };
    },
    computed: {
        sections: function () { return this.buildSections(); },
        visibleIds: function () { return this.sections.map(function (s) { return s.id; }); },
        current: function () {
            var id = this.currentSectionId;
            return this.sections.find(function (s) { return s.id === id; }) || this.sections[0];
        },
        currentIndex: function () {
            var id = this.current && this.current.id;
            var i = this.visibleIds.indexOf(id);
            return i < 0 ? 0 : i;
        },
        progressPct: function () {
            if (this.done) return 100;
            var n = this.visibleIds.length || 1;
            return Math.round((this.currentIndex / n) * 100);
        },
        summary: function () { return this.done ? this.buildSummary() : []; },
    },
    watch: {
        visibleIds: function (ids) {
            if (ids.indexOf(this.currentSectionId) === -1) {
                this.currentSectionId = ids[0] || 'riesgo';
            }
        },
        ctx: {
            deep: true,
            handler: function () { this.ensureStructures(); },
        },
    },
    methods: {
        val: function (path) { return getPath(this.a, path); },
        field: function (path, type, label, opt) {
            opt = opt || {};
            var value = getPath(this.a, path);
            if (value == null) value = type === 'multi' ? [] : '';
            return {
                kind: 'field',
                path: path,
                type: type,
                label: label,
                required: !!opt.required,
                hint: opt.hint || '',
                value: value,
                options: opt.options || (type === 'bool' ? YN : []),
                exclusive: opt.exclusive || '',
                disabledValues: opt.disabledValues || [],
                error: this.errors[path] || '',
            };
        },
        info: function (key, text) { return { kind: 'info', key: key, text: text }; },
        heading: function (key, text) { return { kind: 'heading', key: key, text: text }; },
        ensureStructures: function () {
            if (this.a.ident.conocimiento === 'si' && this.a.ident.entidades.length === 0) {
                this.a.ident.entidades.push(emptyEntidad());
            }
        },
        onSet: function (path, value) {
            if (path === 'plan.equipos' && value.indexOf('medidas') !== -1) {
                value = value.filter(function (v) { return v !== 'psico' && v !== 'osignh'; });
            }
            setPath(this.a, path, value);
            delete this.errors[path];
            var parts = path.split('.');
            var last = parts[parts.length - 1];
            var parent = parts.slice(0, -1).join('.');
            if (last === 'depto') {
                setPath(this.a, parent + '.municipio', '');
                setPath(this.a, parent + '.ciudad', '');
            }
            if (last === 'municipio') setPath(this.a, parent + '.ciudad', '');
            if (last === 'depto' || last === 'municipio') {
                if (getPath(this.a, parent + '.vereda') != null) setPath(this.a, parent + '.vereda', '');
            }
            if (last === 'sector') {
                ['entidad', 'otraEntidad', 'tipos', 'institucion', 'otraTipo', 'otraInst', 'talento', 'eco', 'terr', 'diferencial', 'completadas', 'pendientes'].forEach(function (k) {
                    var cur = getPath(this.a, parent + '.' + k);
                    if (Array.isArray(cur)) setPath(this.a, parent + '.' + k, []);
                    else if (cur != null) setPath(this.a, parent + '.' + k, '');
                }, this);
            }
            if (last === 'entidad' && value !== 'otra') {
                var otra = getPath(this.a, parent + '.otraEntidad');
                if (otra != null) setPath(this.a, parent + '.otraEntidad', '');
            }
            if (path === 'entidadesSel') {
                (value || []).forEach(function (id) {
                    if (!this.a.seg[id]) this.a.seg[id] = emptySeg();
                }, this);
            }
            if (last === 'nuevaBarrera' && value === 'si') {
                var nuevas = getPath(this.a, parent + '.barrerasNuevas');
                if (!nuevas || !nuevas.length) setPath(this.a, parent + '.barrerasNuevas', [emptyBarreraEntidad()]);
            }
            this.ensureStructures();
        },
        territoryFields: function (base) {
            var depto = this.val(base + '.depto');
            var muni = this.val(base + '.municipio');
            var d = TERR.find(function (t) { return t.v === depto; });
            var m = d && d.munis.find(function (x) { return x.v === muni; });
            return [
                this.field(base + '.depto', 'select', 'Departamento', { required: true, options: TERR.map(function (t) { return { v: t.v, l: t.l }; }) }),
                this.field(base + '.municipio', 'select', 'Municipio', { required: true, options: d ? d.munis.map(function (x) { return { v: x.v, l: x.l }; }) : [] }),
                this.field(base + '.ciudad', 'select', 'Ciudad', { required: true, options: m ? m.ciudades.map(function (c) { return { v: c, l: c }; }) : [] }),
            ];
        },
        yn: function (path, label, textMode, textLabel) {
            var out = [this.field(path, 'bool', label, { required: true })];
            var v = this.val(path);
            var show = textMode === 'always' || (textMode === 'si' && v === 'si') || (textMode === 'no' && v === 'no');
            if (show) out.push(this.field(path + 'Texto', 'text', textLabel || 'Describa', { required: true }));
            return out;
        },
        scalePair: function (path, label) {
            return [
                this.field(path, 'single', label, { required: true, options: SCALE, hint: '1 a 5' }),
                this.field(path + 'Texto', 'text', 'Justifique la puntuación', { required: true }),
            ];
        },
        buildSections: function () {
            var sections = [];
            sections.push(this.sectionRiesgo());
            if (this.ctx.hasMedidas) sections.push(this.sectionMedidas());
            if (this.ctx.hasEntidades) sections.push(this.sectionSegEntidades());
            sections.push(this.sectionIdent());
            sections.push(this.sectionPlan());
            sections.push(this.sectionConclusion());
            return sections;
        },
        sectionRiesgo: function () {
            var b = [];
            b.push(this.info('r0', 'Los factores son indicadores técnicos de riesgo de feminicidio, homicidio abyecto o fútil y otras violencias graves. No son preguntas directas a la víctima. Registre solo hechos desde el último contacto hasta hoy.'));
            b.push(this.field('risk.respondiente', 'single', 'Indique quién brindó la información durante este seguimiento', {
                required: true,
                options: opts([
                    ['directo', 'Contacto directo con la víctima'],
                    ['institucion', 'Contacto con institución'],
                    ['indirecto', 'Contacto indirecto (familiar o persona de confianza)'],
                ]),
            }));
            b.push(this.field('risk.hechos', 'bool', '¿Se registraron nuevos hechos de violencia desde el último seguimiento?', { required: true }));
            if (this.a.risk.hechos === 'si') {
                b.push(this.field('risk.hechosTexto', 'text', 'Registre los nuevos hechos: tiempo, modo, lugar y relación de la víctima con la persona agresora', { required: true }));
                b.push(this.field('risk.hechosFecha', 'date', 'Fecha en la que se presentaron los hechos', { required: true }));
            }
            b.push(this.field('risk.protectores', 'multi', 'Indique si estos factores protectores se presentan en el caso', {
                required: true, exclusive: 'ninguna', options: opts(PROTECTORES),
                hint: 'La primera opción del Excel pide adaptar la redacción para personas OSIG. Con 3 o más protectores y ningún riesgo extremo, el sistema sugiere seguimiento general.',
            }));
            b.push(this.field('risk.riesgos', 'multi', 'Indique si estos factores de riesgo se presentan en el caso', {
                required: true, exclusive: 'ninguna', options: opts(RIESGOS),
                hint: 'Con 4 o más factores de riesgo, el sistema sugiere remisión al equipo de riesgo.',
            }));
            b.push(this.field('risk.extremos', 'multi', 'Si identifica algún factor de riesgo extremo, remita al equipo de riesgo. Verifique si ya hay articulación interinstitucional y priorice la protección inmediata', {
                required: true, exclusive: 'ninguna', options: opts(EXTREMOS),
                hint: 'Si no identifica ninguno, marque Ninguna. La hoja no trae esa opción; sin ella la pregunta obligatoria no se podría cerrar.',
            }));
            var nProt = (this.a.risk.protectores || []).filter(function (x) { return x !== 'ninguna'; }).length;
            var nRisk = (this.a.risk.riesgos || []).filter(function (x) { return x !== 'ninguna'; }).length;
            var nExt = (this.a.risk.extremos || []).filter(function (x) { return x !== 'ninguna'; }).length;
            if (nExt > 0) b.push(this.info('r-ext', 'Sugerencia del sistema: hay riesgo extremo. Remita al equipo de riesgo, verifique la articulación interinstitucional y priorice la protección inmediata.'));
            if (nProt >= 3 && nExt === 0) b.push(this.info('r-prot', 'Sugerencia del sistema: hay 3 o más factores protectores y ningún riesgo extremo. Considere la remisión a seguimiento general.'));
            if (nRisk >= 4) b.push(this.info('r-risk', 'Sugerencia del sistema: hay 4 o más factores de riesgo. Considere la remisión al equipo de riesgo.'));
            b.push(this.field('risk.descripcion', 'text', 'Describa la situación de riesgo actual y cómo los factores de riesgo y protectores influyen en que este se mantenga o modifique', { required: true }));
            return { id: 'riesgo', title: 'Valoración del riesgo', desc: 'Seguimiento al riesgo de feminicidio, homicidio abyecto o fútil', blocks: b };
        },
        sectionMedidas: function () {
            var m = this.ctx.momento;
            var p = 'med.' + m;
            var b = [];
            var objetivos = {
                '1': 'Primer seguimiento, a la semana. Verifique la seguridad inmediata, la entrega de las medidas y las necesidades urgentes no cubiertas.',
                '2': 'Segundo seguimiento, a los 15 días. Verifique la activación de las rutas institucionales y fortalezca las redes de apoyo.',
                '3': 'Tercer seguimiento, al mes. Evalúe el avance hacia la estabilización y el funcionamiento de las rutas.',
            };
            b.push(this.info('m0', objetivos[m] + ' Las preguntas van en lenguaje sencillo, sin pedir que la persona repita los hechos. Las medidas de este caso (las trae el sistema) son: ' + DEMO_MEDIDAS.map(function (x) { return x.nombre; }).join(', ') + '.'));
            if (this.ctx.extemporaneo) {
                b.push(this.info('m-ext', 'Este seguimiento es extemporáneo. Registre profesional, justificación y fecha. El cuestionario es el del momento elegido arriba.'));
                b.push(this.field('ext.profesional', 'input', 'Profesional que realiza este seguimiento', { required: true }));
                b.push(this.field('ext.justificacion', 'text', 'Justificación del seguimiento extemporáneo', { required: true }));
                b.push(this.field('ext.fecha', 'date', 'Fecha del seguimiento extemporáneo', { required: true }));
            }
            if (m === '1') b = b.concat(this.momento1(p));
            if (m === '2') b = b.concat(this.momento2(p));
            if (m === '3') b = b.concat(this.momento3(p));
            return { id: 'medidas', title: 'Seguimiento a medidas de emergencia', desc: 'Solo porque el caso tiene medidas activas', blocks: b };
        },
        momento1: function (p) {
            var b = [];
            b = b.concat(this.scalePair(p + '.riesgo', 'En una escala de 1 a 5 (1 muy bajo, 5 muy alto), ¿cómo califica hoy el riesgo para usted y su familia?'));
            b = b.concat(this.scalePair(p + '.seguridad', 'En una escala de 1 a 5 (1 completamente segura, 5 peligro permanente), ¿cómo se siente actualmente?'));
            b.push(this.field(p + '.necesidades', 'single', 'Actualmente considera que sus necesidades básicas están', {
                required: true,
                options: opts([['total', 'Totalmente cubiertas'], ['parcial', 'Parcialmente cubiertas'], ['no', 'No cubiertas']]),
            }));
            b.push(this.field(p + '.necesidadesTexto', 'text', 'Describa el estado de las necesidades básicas', { required: true }));
            b = b.concat(this.yn(p + '.segura', '¿Se siente segura en el lugar donde está viviendo actualmente?', 'always'));
            b = b.concat(this.yn(p + '.contacto', '¿Ha tenido algún contacto con el agresor desde que iniciaron las medidas?', 'always'));
            b = b.concat(this.yn(p + '.amenazas', '¿Ha recibido amenazas, llamadas, mensajes o intentos de ubicarla?', 'always'));
            b = b.concat(this.yn(p + '.situacion', '¿Ha ocurrido alguna situación nueva que aumente su preocupación por su seguridad?', 'always'));
            b = b.concat(this.yn(p + '.hijas', '¿Sus hijas, hijos o personas a cargo se encuentran seguras?', 'always'));
            b.push(this.heading('m1-est', 'Estado de cada medida otorgada'));
            var self = this;
            DEMO_MEDIDAS.forEach(function (med) {
                var bp = p + '.' + med.id;
                b.push(self.heading('h-' + med.id, med.nombre));
                b.push(self.field(bp + '.recibio', 'bool', '¿La recibió oportunamente?', { required: true }));
                if (self.val(bp + '.recibio') === 'no') {
                    b.push(self.field(bp + '.recibioTexto', 'text', '¿Por qué no se entregó a tiempo?', { required: true }));
                }
                if (self.val(bp + '.recibio') === 'si') {
                    b.push(self.field(bp + '.utilizo', 'bool', '¿Ha podido utilizarla sin inconvenientes?', { required: true }));
                    if (self.val(bp + '.utilizo') === 'no') {
                        b.push(self.field(bp + '.utilizoTexto', 'text', 'Describa el inconveniente', { required: true }));
                    }
                    b.push(self.field(bp + '.dificultad', 'bool', '¿Ha presentado alguna dificultad para hacer efectiva la medida?', { required: true }));
                }
            });
            b = b.concat(this.yn(p + '.urgente', '¿Existe alguna necesidad urgente que aún no haya podido resolver?', 'si', 'Describa la necesidad urgente'));
            b.push(this.heading('m1-red', 'Redes de apoyo'));
            b = b.concat(this.yn(p + '.cuenta', '¿Cuenta con alguna persona con la que puede contar en este momento?', 'si', '¿Quién?'));
            b = b.concat(this.yn(p + '.org', 'Además de familiares o amistades, ¿hay alguna persona, organización, comunidad, iglesia, lideresa, colectivo o institución que le esté brindando apoyo?', 'si', '¿Cuál?'));
            b = b.concat(this.yn(p + '.primero', 'Si ocurriera una emergencia hoy, ¿sabe a quién acudir primero?', 'si', '¿A quién?'));
            b = b.concat(this.yn(p + '.traslado', '¿Cuenta con un lugar seguro al que pueda trasladarse si fuera necesario?', 'si', '¿Cuál?'));
            return b;
        },
        momento2: function (p) {
            var b = [];
            b = b.concat(this.scalePair(p + '.riesgo', 'En una escala de 1 a 5 (1 muy bajo, 5 muy alto), ¿cómo califica hoy el riesgo para usted y su familia?'));
            b = b.concat(this.scalePair(p + '.seguridad', 'En una escala de 1 a 5 (1 completamente segura, 5 peligro permanente), ¿cómo se siente actualmente?'));
            b = b.concat(this.yn(p + '.continuan', '¿Las medidas continúan respondiendo a sus necesidades?', 'always', 'Explique si la medida cubrió o no la necesidad'));
            b = b.concat(this.yn(p + '.ajuste', '¿Ha requerido algún ajuste?', 'si', 'Describa los ajustes'));
            b.push(this.field(p + '.denuncia', 'bool', '¿Ha presentado denuncia?', { required: true }));
            if (this.val(p + '.denuncia') === 'no') {
                b.push(this.field(p + '.denunciaTexto', 'text', 'Motivos por los que no ha denunciado', { required: true }));
            }
            if (this.val(p + '.denuncia') === 'si') {
                var rutas = [
                    ['comisaria', '¿Ha acudido a la Comisaría de Familia?'],
                    ['fiscalia', '¿Ha recibido atención por parte de la Fiscalía?'],
                    ['policia', '¿La Policía ha intervenido en su caso?'],
                    ['defensoria', '¿Ha recibido acompañamiento de la Defensoría del Pueblo?'],
                    ['salud', '¿Ha recibido atención en salud física o mental?'],
                    ['alcaldia', '¿Ha recibido apoyo de la Alcaldía, Gobernación u otra entidad?'],
                    ['apoyos', '¿Está recibiendo apoyos adicionales de otras instituciones?'],
                ];
                rutas.forEach(function (r) {
                    b = b.concat(this.yn(p + '.' + r[0], r[1], 'si', 'Describa el tipo de atención y lo otorgado'));
                }, this);
                b = b.concat(this.yn(p + '.radicado', '¿Le asignaron número de radicado o noticia criminal?', 'si', 'Escriba el número'));
                b = b.concat(this.yn(p + '.proteccion', '¿Le otorgaron medidas de protección?', 'si', 'Indique el número de medidas de emergencia'));
            }
            b.push(this.heading('m2-red', 'Redes y recursos'));
            b.push(this.field(p + '.redesPersonales', 'text', 'Redes personales: ¿quiénes la están acompañando actualmente?', { required: true }));
            b = b.concat(this.yn(p + '.redesCom', 'Redes comunitarias: ¿hace parte de algún grupo de mujeres, organización, iglesia u otro espacio?', 'si', '¿Cuál?'));
            b = b.concat(this.yn(p + '.redesInst', 'Redes institucionales: ¿conoce qué entidades puede contactar si el riesgo aumenta?', 'si', '¿Cuáles?'));
            b.push(this.field(p + '.recursos', 'text', 'Recursos personales: ¿qué fortalezas le han ayudado hasta ahora?', { required: true }));
            b.push(this.field(p + '.cuentaCon', 'bool', 'Actualmente, ¿cuenta con recursos de sostén (lugar, cuidado, salud, comunicación, transporte, documentos, apoyo emocional)?', { required: true }));
            if (this.val(p + '.cuentaCon') === 'si') {
                b = b.concat(this.yn(p + '.lugar', '¿Un lugar seguro?', 'si', '¿Cuál?'));
                b = b.concat(this.yn(p + '.cuidador', '¿Alguien que cuide a sus hijos si fuera necesario?', 'si', '¿Quién?'));
                b = b.concat(this.yn(p + '.accesoSalud', '¿Acceso a salud?', 'si', 'Describa'));
                b = b.concat(this.yn(p + '.telefono', '¿Teléfono o internet?', 'si', 'Describa'));
                b = b.concat(this.yn(p + '.transporte', '¿Transporte?', 'si', '¿Qué tipo de transporte?'));
                b = b.concat(this.yn(p + '.documentos', '¿Documentos personales?', 'no', '¿Qué documentos no tiene?'));
                b = b.concat(this.yn(p + '.apoyo', '¿Apoyo emocional?', 'no', '¿Quién puede brindarlo?'));
            }
            return b;
        },
        momento3: function (p) {
            var b = [];
            b = b.concat(this.scalePair(p + '.riesgo', 'En una escala de 1 a 5 (1 muy bajo, 5 muy alto), ¿cómo califica hoy el riesgo para usted y su familia?'));
            b = b.concat(this.scalePair(p + '.seguridad', 'En una escala de 1 a 5 (1 completamente segura, 5 peligro permanente), ¿cómo se siente actualmente?'));
            b = b.concat(this.yn(p + '.continuan', '¿Las medidas continúan respondiendo a sus necesidades?', 'always', 'Explique si la medida cubrió o no la necesidad'));
            b = b.concat(this.yn(p + '.ajuste', '¿Ha requerido algún ajuste?', 'si', 'Describa los ajustes'));
            b.push(this.field(p + '.avance', 'text', '¿Cómo han avanzado sus procesos?', { required: true }));
            b = b.concat(this.yn(p + '.audiencias', '¿Ha tenido audiencias?', 'si', 'Describa el proceso jurídico'));
            b = b.concat(this.yn(p + '.revictimizacion', '¿Ha sentido revictimización?', 'si', 'Describa las situaciones'));
            b = b.concat(this.yn(p + '.representacion', '¿Cuenta actualmente con representación jurídica?', 'si', '¿Cuál?'));
            b = b.concat(this.yn(p + '.ingresos', '¿Actualmente tiene ingresos?', 'si', '¿De qué tipo?'));
            b = b.concat(this.yn(p + '.empleo', '¿Ha podido conseguir empleo?', 'si', 'Describa'));
            b = b.concat(this.yn(p + '.productiva', '¿Está desarrollando alguna actividad productiva?', 'si', 'Describa'));
            b.push(this.field(p + '.necesidades', 'text', '¿Cuáles siguen siendo sus principales necesidades?', { required: true }));
            b = b.concat(this.yn(p + '.redes', '¿Han cambiado sus redes de apoyo?', 'si', 'Describa el cambio'));
            b.push(this.field(p + '.recursosHoy', 'text', '¿Qué recursos siente que hoy tiene para afrontar una nueva situación de riesgo?', { required: true }));
            b.push(this.field(p + '.fortalecer', 'text', '¿Qué recursos considera que aún necesita fortalecer?', { required: true }));
            return b;
        },
        sectionSegEntidades: function () {
            var b = [];
            b.push(this.info('e0', 'El listado lo trae el sistema. Este apartado integra las acciones nuevas de la entidad, las barreras nuevas y el seguimiento a las barreras que ya tenía.'));
            var lista = DEMO_ENTIDADES.map(function (e) {
                var extra = e.barreras.length ? ' — barreras: ' + e.barreras.map(function (x) { return x.nombre; }).join('; ') : ' — sin barreras previas';
                return e.nombre + extra;
            }).join('\n');
            b.push(this.info('e-list', 'Entidades del caso:\n' + lista));
            b.push(this.field('entidadesSel', 'multi', 'Selección de la entidad a actualizar', {
                required: true,
                options: DEMO_ENTIDADES.map(function (e) { return { v: e.id, l: e.nombre }; }),
            }));
            var items = [];
            (this.a.entidadesSel || []).forEach(function (id) {
                var ent = DEMO_ENTIDADES.find(function (e) { return e.id === id; });
                if (!ent || !this.a.seg[id]) return;
                items.push({
                    key: id,
                    title: ent.nombre,
                    canDelete: false,
                    fields: this.fieldsSegEntidad(ent),
                });
            }, this);
            if (items.length) {
                b.push({ kind: 'repeater', key: 'seg-cards', canAdd: false, items: items });
            }
            return { id: 'segEntidades', title: 'Seguimiento a entidades', desc: 'Solo porque el caso ya tiene entidades', blocks: b };
        },
        fieldsSegEntidad: function (ent) {
            var base = 'seg.' + ent.id;
            var row = this.a.seg[ent.id] || emptySeg();
            var f = [];
            f.push(this.field(base + '.completadas', 'multi', 'Acciones institucionales que ya se completaron', { options: optsNames(ACCIONES.proteccion.concat(ACCIONES.salud)) }));
            if (ent.barreras.length) {
                f.push(this.heading(ent.id + '-segb', 'Seguimiento a barreras de esta entidad'));
                ent.barreras.forEach(function (bar) {
                    f.push(this.tarjetaSeguimientoBarrera(base, ent.id, bar, (row.barreras || {})[bar.id] || {}));
                }, this);
            }
            f.push(this.info(ent.id + '-pend', 'Pendientes que trae el sistema: ' + ent.pendientes.join(', ') + '.'));
            f.push(this.field(base + '.pendientes', 'multi', 'Acciones institucionales pendientes', {
                options: optsNames(ACCIONES.proteccion.concat(ACCIONES.salud)).concat([{ v: 'otra', l: 'Otro' }]),
            }));
            if ((row.pendientes || []).indexOf('otra') !== -1) {
                f.push(this.field(base + '.otraPendiente', 'text', 'Otra acción pendiente', { required: true }));
            }
            f.push(this.field(base + '.actualizacion', 'text', 'Describa cualquier actualización de la ruta desde el último seguimiento', {
                required: true, hint: 'Si no hay novedad, puede escribir No o no aplica.',
            }));
            f.push(this.heading(ent.id + '-act', 'Registro de barreras y activación de ruta'));
            f.push(this.field(base + '.nuevaBarrera', 'bool', '¿Se identifica alguna nueva barrera en este seguimiento con relación a esta entidad?', {
                required: true,
                hint: 'Si marca Sí, el registro de la barrera se abre aquí, asociado a esta entidad. En el sector puede elegir Transversales.',
            }));
            if (row.nuevaBarrera === 'si') {
                f.push(this.barrerasNuevasRepeater(base + '.barrerasNuevas', row.barrerasNuevas, 'seg:' + ent.id));
            }
            var activacion = [
                this.field(base + '.acuerdo', 'bool', '¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes?', { required: true }),
            ];
            if (row.acuerdo === 'si') {
                activacion.push(this.field(base + '.debidaDiligencia', 'bool', '¿Salvia debe activar rutas en el marco de la debida diligencia?', { required: true }));
                activacion.push(this.field(base + '.canales', 'multi', 'Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención', {
                    required: true, exclusive: 'ninguna', options: opts(CANALES),
                }));
                activacion.push(this.field(base + '.gestionesRuta', 'text', 'Describa las gestiones realizadas para la activación de rutas y gestión de barreras en esta entidad'));
            }
            f.push({
                kind: 'panel',
                key: ent.id + '-actruta',
                title: 'Activación de ruta',
                fields: activacion,
            });
            return f;
        },
        sectionIdent: function () {
            var b = [];
            b.push(this.field('ident.conocimiento', 'bool', '¿Ha puesto en conocimiento de alguna entidad o institución los hechos de violencia reportados?', { required: true }));
            if (this.a.ident.conocimiento === 'si') {
                var items = this.a.ident.entidades.map(function (ent, i) {
                    return {
                        key: ent.uid,
                        title: 'Entidad ' + (i + 1) + (ent.sector ? ' — ' + (SECTORES.find(function (s) { return s[0] === ent.sector; }) || ['', ent.sector])[1] : ''),
                        canDelete: this.a.ident.entidades.length > 1,
                        remove: 'ident',
                        index: i,
                        fields: this.fieldsIdent(i),
                    };
                }, this);
                b.push({
                    kind: 'repeater', key: 'ident-rep', canAdd: true, addId: 'ident',
                    addLabel: 'Agregar entidad a la que ya acudió', items: items,
                });
            }
            return { id: 'idEntidades', title: 'Identificación de entidades', desc: 'Entidades a las que la persona ya acudió, o a las que otras activaron ruta', blocks: b };
        },
        fieldsIdent: function (i) {
            var base = 'ident.entidades.' + i;
            var ent = this.a.ident.entidades[i] || {};
            var f = [];
            f.push(this.field(base + '.sector', 'select', 'Sector de la entidad', { required: true, options: opts(SECTORES) }));
            f = f.concat(this.territoryFields(base));
            var entOpts = optsNames(ENTIDADES[ent.sector] || []).concat([{ v: 'otra', l: 'Otra' }]);
            f.push(this.field(base + '.entidad', 'select', 'Entidad', { required: true, options: entOpts }));
            if (ent.entidad === 'otra') f.push(this.field(base + '.otraEntidad', 'input', 'Otra ¿cuál?', { required: true }));
            f.push(this.field(base + '.radicado', 'input', 'Número de radicado (NUNC, medida de protección). Escriba tipo y número', { required: true, hint: 'Ejemplo: NUNC 2026120956789' }));
            f.push(this.field(base + '.fecha', 'date', 'Fecha en la que reportó a la entidad', { required: true }));
            var acc = optsNames(ACCIONES[ent.sector] || ACCIONES.otra);
            f.push(this.field(base + '.completadas', 'multi', 'Acciones institucionales completadas', { required: true, options: acc }));
            f.push(this.field(base + '.pendientes', 'multi', 'Acciones institucionales pendientes', { required: true, options: acc }));
            f.push(this.field(base + '.otraRuta', 'text', 'Otra información asociada a la Ruta', { required: true }));
            f.push(this.field(base + '.nuevaBarrera', 'bool', '¿Se identifica alguna nueva barrera con relación a esta entidad?', {
                required: true,
                hint: 'Si marca Sí, el registro de la barrera se abre aquí, asociado a esta entidad.',
            }));
            if (ent.nuevaBarrera === 'si') {
                f.push(this.barrerasNuevasRepeater(base + '.barrerasNuevas', ent.barrerasNuevas, 'identb:' + i));
            }
            f.push(this.field(base + '.debidaDiligencia', 'bool', '¿Salvia debe activar rutas en el marco de la debida diligencia?', { required: true }));
            if (ent.debidaDiligencia === 'si') {
                f.push(this.field(base + '.acuerdo', 'bool', '¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes?', { required: true }));
                f.push(this.field(base + '.canales', 'multi', 'Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención', {
                    required: true, exclusive: 'ninguna', options: opts(CANALES),
                }));
                f.push(this.field(base + '.gestionesRuta', 'text', '¿Qué gestiones se realizaron desde Salvia para la activación de ruta?', { required: true }));
            }
            return f;
        },
        tarjetaSeguimientoBarrera: function (base, entId, bar, ans) {
            var bb = base + '.barreras.' + bar.id;
            var fields = [];
            fields.push(this.field(bb + '.persisten', 'bool', '¿Las barreras persisten?', { required: true }));
            fields.push(this.field(bb + '.cierre', 'bool', '¿Se realiza cierre de la barrera?', { required: true }));
            if (ans.cierre === 'si') {
                fields.push(this.field(bb + '.motivoCierre', 'single', 'Motivo del cierre', { required: true, options: opts(CIERRE_MOTIVOS) }));
            }
            fields.push(this.field(bb + '.gestion', 'single', 'Gestión de la barrera por parte de Salvia', {
                options: opts([
                    ['activacion', 'Activación de ruta'],
                    ['alerta', 'Alerta por barreras'],
                ]),
                hint: 'Activación de ruta queda registrada como gestión y no abre otra sección. Alerta por barreras notifica al equipo técnico.',
            }));
            if (ans.gestion === 'alerta') {
                fields.push(this.info(entId + '-' + bar.id + '-al', 'Alerta por barreras: al guardar de verdad, esto notificaría al equipo técnico. Este prototipo no envía la notificación.'));
            }
            return { kind: 'card', key: bb, title: bar.nombre, fields: fields };
        },
        barrerasNuevasRepeater: function (listPath, list, addId) {
            var items = (list || []).map(function (bar, i) {
                return {
                    key: bar.uid,
                    title: 'Barrera de esta entidad ' + (i + 1),
                    canDelete: (list || []).length > 1,
                    remove: addId,
                    index: i,
                    fields: this.fieldsBarreraEntidad(listPath + '.' + i, bar),
                };
            }, this);
            return {
                kind: 'repeater',
                key: listPath,
                canAdd: true,
                addId: addId,
                addLabel: 'Agregar otra barrera de esta entidad',
                items: items,
            };
        },
        fieldsBarreraEntidad: function (base, bar) {
            bar = bar || {};
            var f = [];
            f.push(this.field(base + '.sector', 'single', 'Sector en el que se presentó la barrera', {
                required: true,
                options: opts([
                    ['salud', 'Salud'],
                    ['justicia', 'Justicia'],
                    ['proteccion', 'Protección'],
                    ['transversales', 'Transversales'],
                ]),
            }));
            if (bar.sector === 'transversales') {
                f.push(this.heading(base + '-est', 'Barreras transversales'));
                f.push(this.field(base + '.talento', 'multi', 'Barreras institucionales y de talento humano', { options: opts(TALENTO) }));
                f.push(this.field(base + '.eco', 'multi', 'Barreras económicas y socioeconómicas', { options: opts(ECO) }));
                f.push(this.field(base + '.terr', 'multi', 'Barreras territoriales y geográficas', { options: opts(TERR_B) }));
                f.push(this.field(base + '.diferencial', 'multi', 'Barreras por ausencia de enfoque diferencial', { options: opts(DIFERENCIAL) }));
            }
            if (bar.sector === 'salud') {
                f.push(this.info(base + '-hs', 'Salud: hospitales, IPS, EPS y puestos de salud.'));
                f.push(this.field(base + '.tipos', 'multi', 'Barreras en salud', { options: opts(SALUD_B) }));
                f.push(this.field(base + '.institucion', 'select', '¿A qué institución acudió?', { options: optsNames(INST.salud) }));
                f.push(this.field(base + '.otraTipo', 'text', 'Otra barrera en sector Salud'));
            }
            if (bar.sector === 'justicia') {
                f.push(this.info(base + '-hj', 'Justicia: fiscalía, policía, comisarías, jueces y medicina legal.'));
                f.push(this.field(base + '.tipos', 'multi', 'Barreras en justicia', { options: opts(JUSTICIA_B) }));
                f.push(this.field(base + '.institucion', 'select', '¿A qué institución acudió?', { options: optsNames(INST.justicia) }));
                f.push(this.field(base + '.otraTipo', 'text', 'Otra barrera en sector Justicia'));
            }
            if (bar.sector === 'proteccion') {
                f.push(this.info(base + '-hp', 'Protección: comisarías, ICBF, UNP, defensoría de familia e inspección.'));
                f.push(this.field(base + '.tipos', 'multi', 'Barreras en protección', { options: opts(PROTECCION_B) }));
                f.push(this.field(base + '.institucion', 'select', '¿A qué institución acudió?', { options: optsNames(INST.proteccion) }));
                f.push(this.field(base + '.otraTipo', 'text', 'Otra barrera en el marco del sector Protección'));
            }
            f = f.concat(this.fieldsBarreraComunes(base, bar));
            return f;
        },
        fieldsBarreraComunes: function (base, bar) {
            var f = [];
            f.push(this.field(base + '.fecha', 'date', 'Fecha en la que se presentó la barrera (aproximada)', { required: true }));
            f.push(this.field(base + '.funcionario', 'input', 'Funcionario/a o dependencia donde se presentó la barrera', { required: true }));
            f = f.concat(this.territoryFields(base).slice(0, 2));
            var depto = this.val(base + '.depto');
            var muni = this.val(base + '.municipio');
            if (depto && muni) {
                f.push(this.field(base + '.vereda', 'select', 'Vereda o zona', {
                    required: true,
                    options: opts([['cabecera', 'Cabecera'], ['rural', 'Vereda o zona rural'], ['no', 'No precisa']]),
                }));
            }
            f.push(this.field(base + '.descripcion', 'text', 'Descripción de la barrera (tiempo, modo y lugar)', { required: true }));
            f.push(this.field(base + '.gestion', 'multi', 'Gestión de la barrera', {
                required: true,
                options: opts([
                    ['activacion', 'Activación de ruta interinstitucional (Generar oficio)'],
                    ['alerta', 'Alerta por barreras'],
                ]),
                hint: 'Activación de ruta queda registrada como gestión y no abre otra sección. Alerta por barreras informa al equipo técnico.',
            }));
            if ((bar.gestion || []).indexOf('alerta') !== -1) {
                f.push(this.info(base + '-al', 'Alerta por barreras notificaría al equipo técnico. Este prototipo no envía la notificación.'));
                f.push(this.field(base + '.perspectiva', 'text', '¿Cuál es su perspectiva profesional frente a la barrera?', { required: true }));
            }
            f.push(this.field(base + '.gestionTexto', 'text', 'Descripción de la gestión realizada con relación a las barreras', { required: true }));
            return f;
        },
        sectionPlan: function () {
            var b = [];
            b.push(this.field('plan.remision', 'bool', 'Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA', { required: true }));
            if (this.a.plan.remision === 'si') {
                var disabled = (this.a.plan.equipos || []).indexOf('medidas') !== -1 ? ['psico', 'osignh'] : [];
                b.push(this.field('plan.equipos', 'multi', 'Equipos', {
                    required: true, options: opts(EQUIPOS), disabledValues: disabled,
                    hint: 'Si remite a medidas de emergencia no puede remitir a atención psicosocial ni a psicosocial OSIGNH.',
                }));
                if (disabled.length) b.push(this.info('p-ex', 'Remitió a medidas de emergencia. Atención psicosocial y psicosocial OSIGNH quedan deshabilitadas.'));
                if ((this.a.plan.equipos || []).indexOf('medidas') !== -1) {
                    b.push(this.field('plan.medidasTipos', 'multi', 'Medidas de emergencia a las que remite', { required: true, options: opts(MEDIDAS_TIPO) }));
                }
                if ((this.a.plan.equipos || []).indexOf('psico') !== -1) {
                    b.push(this.field('plan.psico', 'multi', 'Criterios de remisión a psicosocial', { required: true, options: opts(PSICO_CRIT) }));
                    var pts = (this.a.plan.psico || []).reduce(function (s, id) { return s + (PSICO_PTS[id] || 0); }, 0);
                    b.push(this.info('p-pts', 'Puntaje marcado en criterios con puntaje: ' + pts + '. La hoja V.4 no fija un mínimo para bloquear la remisión; el prototipo solo lo muestra.'));
                }
                if ((this.a.plan.equipos || []).indexOf('osignh') !== -1) {
                    b.push(this.field('plan.osignh', 'multi', 'Criterios de psicosocial para personas OSIGNH', {
                        required: true,
                        options: opts([
                            ['obligatorio', 'Ser OSIGNH, mayor de 18 años, víctima de VpP, en post-emergencia, con voluntariedad y disponibilidad. Obligatorio'],
                            ['eps', 'No aplica si ya recibe atención psicosocial por EPS, salvo barreras por orientación sexual o identidad de género'],
                        ]),
                    }));
                    if ((this.a.plan.osignh || []).indexOf('eps') !== -1) {
                        b.push(this.info('p-eps', 'Si ya recibe atención por EPS y no hay barreras por orientación sexual o identidad de género, la remisión a OSIGNH no aplica.'));
                    }
                }
                if ((this.a.plan.equipos || []).indexOf('hombres') !== -1) {
                    b.push(this.info('p-hom', 'Criterio de Atención Hombres: hombres diversos victimarios o posibles victimarios.'));
                }
                if ((this.a.plan.equipos || []).indexOf('discapacidad') !== -1) {
                    b.push(this.field('plan.discapacidad', 'multi', 'Equipo Discapacidad', {
                        required: true,
                        options: opts([
                            ['lsc', 'Apoyo comunicativo en LSC'],
                            ['enfoque', 'Apoyo con enfoque de discapacidad (ley de capacidad legal)'],
                        ]),
                    }));
                }
                if ((this.a.plan.equipos || []).indexOf('estabilizacion') !== -1) {
                    b.push(this.field('plan.estabilizacion', 'multi', 'Estabilización', {
                        required: true,
                        options: opts([
                            ['ingresos', 'No tiene ninguna forma de ingresos o tiene ingresos de forma inestable'],
                            ['vivienda', 'La persona reside actualmente en vivienda temporal o familiar'],
                            ['cargo', 'La persona tiene personas a cargo que limitan su tiempo para trabajar o estudiar'],
                            ['depende', 'Depende económicamente de la persona agresora'],
                        ]),
                    }));
                }
            }
            b.push(this.field('plan.elementos', 'multi', 'Seleccione los elementos psicosociales aplicados en este seguimiento', { required: true, options: opts(ELEMENTOS) }));
            b.push(this.field('plan.acciones', 'text', 'Describa las acciones realizadas en el seguimiento', { required: true }));
            b.push(this.field('plan.ajustes', 'text', 'Describa los ajustes que realiza al plan de acción de Salvia', {
                required: true,
                hint: 'En la versión definitiva este texto se acumula al plan de acción, de forma similar a los nuevos hechos.',
            }));
            return { id: 'plan', title: 'Seguimiento al plan de acción', desc: 'Remisiones, acompañamiento psicosocial y descripción', blocks: b };
        },
        sectionConclusion: function () {
            var b = [];
            b.push(this.field('con.concluye', 'bool', '¿Se concluye el ciclo de atención?', { required: true }));
            if (this.a.con.concluye === 'si') {
                b.push(this.info('c0', 'Registre de forma clara y sustentada las razones.'));
                b.push(this.field('con.motivo', 'single', 'Motivo de la conclusión del ciclo de atención', { required: true, options: opts(CONCLUSION_MOTIVOS) }));
                if (this.a.con.motivo === 'otro') b.push(this.field('con.otro', 'input', 'Otro ¿cuál?', { required: true }));
                b.push(this.info('c-desc', 'Registrar de forma clara la descripción del criterio y la orientación brindada sobre continuidad de derechos y rutas disponibles:\n- Intentos de contacto (fechas, medios, número de intentos y tiempo transcurrido sin respuesta)\n- Solicitud de la persona (qué solicita, razones expresadas)\n- Avances y acciones completadas con evidencia y los criterios que permiten concluir que el plan de atención se cumplió\n- No corresponsabilidad (explique la ausencia de gestiones o de compromiso con el acompañamiento)\n- No clasificación como VBG/VpP (explicar por qué no se clasifica)'));
                b.push(this.field('con.descripcion', 'text', 'Descripción del criterio y de la orientación brindada sobre continuidad de derechos y rutas disponibles', { required: true }));
                b.push(this.field('con.accion', 'bool', '¿Se realizó alguna acción institucional en la conclusión del caso?', { required: true }));
                if (this.a.con.accion === 'si') {
                    var conocidas = [];
                    var vista = {};
                    this.a.ident.entidades.forEach(function (e) {
                        var item = null;
                        if (e.entidad === 'otra' && e.otraEntidad) item = { v: 'txt:' + e.otraEntidad, l: e.otraEntidad };
                        else if (e.entidad) {
                            var cat = optsNames(ENTIDADES[e.sector] || []);
                            var found = cat.find(function (o) { return o.v === e.entidad; });
                            if (found) item = found;
                        }
                        if (item && !vista[item.v]) { vista[item.v] = true; conocidas.push(item); }
                    });
                    DEMO_ENTIDADES.forEach(function (e) {
                        if ((this.a.entidadesSel || []).indexOf(e.id) !== -1 && !vista[e.id]) {
                            vista[e.id] = true;
                            conocidas.push({ v: e.id, l: e.nombre });
                        }
                    }, this);
                    conocidas.push({ v: 'otra', l: 'Otra' });
                    b.push(this.field('con.entidad', 'select', 'Seleccione la entidad con la que realizará acción institucional', { required: true, options: conocidas }));
                    if (this.a.con.entidad === 'otra') b.push(this.field('con.otraEntidad', 'input', 'Otra ¿cuál?', { required: true }));
                    b.push(this.field('con.canales', 'multi', 'Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención', {
                        required: true, exclusive: 'ninguna', options: opts(CANALES),
                    }));
                }
            }
            return { id: 'conclusion', title: 'Conclusión del ciclo de atención', desc: 'Criterio, motivo y acción institucional', blocks: b };
        },
        isEmpty: function (field) {
            var v = field.value;
            if (field.type === 'multi') return !v || !v.length;
            return v == null || String(v).trim() === '';
        },
        missingOf: function (section) {
            var missing = [];
            var self = this;
            function walk(blocks) {
                (blocks || []).forEach(function (block) {
                    if (block.kind === 'field' && block.required && self.isEmpty(block)) missing.push(block);
                    if (block.kind === 'card' || block.kind === 'panel') walk(block.fields);
                    if (block.kind === 'repeater') {
                        (block.items || []).forEach(function (item) { walk(item.fields); });
                    }
                });
            }
            if (section) walk(section.blocks);
            if (section && section.id === 'plan' && this.a.plan.remision === 'si' && (this.a.plan.equipos || []).indexOf('osignh') !== -1) {
                if ((this.a.plan.osignh || []).indexOf('obligatorio') === -1 && !missing.some(function (m) { return m.path === 'plan.osignh'; })) {
                    missing.push({ path: 'plan.osignh' });
                }
            }
            return missing;
        },
        scrollToError: function () {
            var self = this;
            this.$nextTick(function () {
                var el = self.$el.querySelector('.df-error');
                if (el && el.scrollIntoView) el.scrollIntoView({ behavior: 'smooth', block: 'center' });
            });
        },
        next: function () {
            var missing = this.missingOf(this.current);
            if (missing.length) {
                var errors = {};
                missing.forEach(function (m) { errors[m.path] = 'Obligatorio'; });
                this.errors = errors;
                this.popupText = 'Revise los campos marcados en esta sección.';
                this.popup = true;
                this.scrollToError();
                return;
            }
            this.passed[this.currentSectionId] = true;
            this.errors = {};
            var ids = this.visibleIds;
            var i = ids.indexOf(this.currentSectionId);
            if (i < ids.length - 1) {
                this.currentSectionId = ids[i + 1];
                window.scrollTo(0, 0);
                return;
            }
            for (var s = 0; s < this.sections.length; s++) {
                var miss = this.missingOf(this.sections[s]);
                if (miss.length) {
                    var errs = {};
                    miss.forEach(function (m) { errs[m.path] = 'Obligatorio'; });
                    this.errors = errs;
                    this.popupText = 'Falta completar "' + this.sections[s].title + '".';
                    this.popup = true;
                    this.currentSectionId = this.sections[s].id;
                    this.scrollToError();
                    return;
                }
            }
            this.done = true;
            window.scrollTo(0, 0);
        },
        prev: function () {
            this.errors = {};
            var ids = this.visibleIds;
            var i = ids.indexOf(this.currentSectionId);
            if (i > 0) this.currentSectionId = ids[i - 1];
        },
        goTo: function (id) {
            if (this.visibleIds.indexOf(id) === -1) return;
            this.errors = {};
            this.currentSectionId = id;
            window.scrollTo(0, 0);
        },
        addItem: function (id) {
            if (id === 'ident') this.a.ident.entidades.push(emptyEntidad());
            if (String(id).indexOf('seg:') === 0) {
                this.a.seg[id.slice(4)].barrerasNuevas.push(emptyBarreraEntidad());
            }
            if (String(id).indexOf('identb:') === 0) {
                this.a.ident.entidades[Number(id.slice(7))].barrerasNuevas.push(emptyBarreraEntidad());
            }
        },
        removeItem: function (id, index) {
            if (id === 'ident' && this.a.ident.entidades.length > 1) this.a.ident.entidades.splice(index, 1);
            if (String(id).indexOf('seg:') === 0) {
                var segList = this.a.seg[id.slice(4)].barrerasNuevas;
                if (segList.length > 1) segList.splice(index, 1);
            }
            if (String(id).indexOf('identb:') === 0) {
                var identList = this.a.ident.entidades[Number(id.slice(7))].barrerasNuevas;
                if (identList.length > 1) identList.splice(index, 1);
            }
        },
        sectionState: function (id) {
            if (id === this.currentSectionId) return 'active';
            if (this.passed[id]) return 'completed';
            return 'enabled';
        },
        buildSummary: function () {
            var out = [{
                id: 'ctx',
                title: 'Contexto simulado del caso',
                rows: [
                    { label: 'Entidades previas', value: this.ctx.hasEntidades ? 'Sí' : 'No' },
                    { label: 'Medidas de emergencia', value: this.ctx.hasMedidas ? 'Sí' : 'No' },
                    { label: 'Momento de medidas', value: this.ctx.momento === '1' ? 'Primera semana' : this.ctx.momento === '2' ? '15 días' : 'Primer mes' },
                    { label: 'Seguimiento extemporáneo', value: this.ctx.extemporaneo ? 'Sí' : 'No' },
                ],
            }];
            var self = this;
            function walk(blocks, rows) {
                (blocks || []).forEach(function (block) {
                    if (block.kind === 'field') {
                        var text = humanField(block);
                        if (text) rows.push({ label: block.label, value: text });
                    }
                    if (block.kind === 'card' || block.kind === 'panel') {
                        rows.push({ label: block.title, value: '', head: true });
                        walk(block.fields, rows);
                    }
                    if (block.kind === 'repeater') {
                        (block.items || []).forEach(function (item) {
                            rows.push({ label: item.title, value: '', head: true });
                            walk(item.fields, rows);
                        });
                    }
                });
            }
            this.sections.forEach(function (sec) {
                var rows = [];
                walk(sec.blocks, rows);
                out.push({ id: sec.id, title: sec.title, rows: rows });
            });
            return out;
        },
    },
    template: `
<div>
    <div class="welcome-section mb-3">
        <h2><b>Prototipo — Formulario de seguimiento V.4</b></h2>
    </div>
    <div class="proto-banner">
        Esta pantalla no guarda nada. Sirve para revisar el formulario y sus condiciones de visibilidad.
        Al finalizar, las respuestas quedan solo en esta página.
    </div>
    <div class="proto-context">
        <h3>Contexto del caso</h3>
        <p>Simula lo que el sistema ya sabría del caso. No es una respuesta del formulario.</p>
        <div class="proto-grid">
            <label><input type="checkbox" v-model="ctx.hasEntidades"> Hay entidades registradas antes</label>
            <label><input type="checkbox" v-model="ctx.hasMedidas"> Hay medidas de emergencia activas</label>
            <label v-if="ctx.hasMedidas">
                Momento
                <select v-model="ctx.momento">
                    <option value="1">1 — primera semana</option>
                    <option value="2">2 — 15 días</option>
                    <option value="3">3 — primer mes</option>
                </select>
            </label>
            <label v-if="ctx.hasMedidas"><input type="checkbox" v-model="ctx.extemporaneo"> Seguimiento extemporáneo</label>
        </div>
    </div>

    <div v-if="!done" class="df-wrapper">
        <aside class="df-sidebar">
            <div class="df-sidebar-header"><p class="df-sidebar-title">Secciones</p></div>
            <nav class="df-sidebar-nav">
                <button type="button" v-for="(sec, i) in sections" :key="sec.id"
                    class="df-section-item" :class="sectionState(sec.id)"
                    @click="goTo(sec.id)">
                    <span class="df-badge" :class="sectionState(sec.id)">{{ i + 1 }}</span>
                    <span>{{ sec.title }}</span>
                </button>
            </nav>
            <div class="df-progress-wrap">
                <div class="df-progress-label"><span>Avance</span><span>{{ progressPct }}%</span></div>
                <div class="df-progress-bar-bg"><div class="df-progress-bar-fill" :style="{ width: progressPct + '%' }"></div></div>
            </div>
        </aside>
        <section class="df-main" v-if="current">
            <header class="df-section-header">
                <div class="df-section-header-row">
                    <span class="df-section-number">{{ currentIndex + 1 }}</span>
                    <h3 class="df-section-title">{{ current.title }}</h3>
                </div>
                <p class="df-section-desc">{{ current.desc }}</p>
            </header>
            <div class="df-questions">
                <template v-for="block in current.blocks" :key="block.key || block.path">
                    <div v-if="block.kind==='info'" class="df-info-banner">
                        <span class="df-info-banner-icon">ℹ️</span>
                        <span>{{ block.text }}</span>
                    </div>
                    <div v-else-if="block.kind==='heading'" class="proto-heading">{{ block.text }}</div>
                    <div v-else-if="block.kind==='field'" class="df-question">
                        <label>{{ block.label }}<span v-if="block.required" class="df-required">*</span></label>
                        <div v-if="block.hint" class="df-hint">{{ block.hint }}</div>
                        <field-view :field="block" @set="onSet"></field-view>
                        <div v-if="block.error" class="df-error">{{ block.error }}</div>
                    </div>
                    <div v-else-if="block.kind==='repeater'" class="df-repeater">
                        <div v-for="item in block.items" :key="item.key" class="df-repeater-item">
                            <div class="df-repeater-item-header">
                                <p class="df-repeater-item-title">{{ item.title }}</p>
                                <button v-if="item.canDelete" type="button" class="df-repeater-delete" @click="removeItem(item.remove, item.index)">Eliminar</button>
                            </div>
                            <template v-for="sub in item.fields" :key="sub.key || sub.path">
                                <div v-if="sub.kind==='info'" class="df-info-banner">
                                    <span class="df-info-banner-icon">ℹ️</span>
                                    <span>{{ sub.text }}</span>
                                </div>
                                <div v-else-if="sub.kind==='heading'" class="proto-heading">{{ sub.text }}</div>
                                <div v-else-if="sub.kind==='field'" class="df-question">
                                    <label>{{ sub.label }}<span v-if="sub.required" class="df-required">*</span></label>
                                    <div v-if="sub.hint" class="df-hint">{{ sub.hint }}</div>
                                    <field-view :field="sub" @set="onSet"></field-view>
                                    <div v-if="sub.error" class="df-error">{{ sub.error }}</div>
                                </div>
                                <div v-else-if="sub.kind==='card' || sub.kind==='panel'" :class="sub.kind==='card' ? 'proto-barrera-card' : 'proto-panel'">
                                    <p :class="sub.kind==='card' ? 'proto-barrera-title' : 'proto-panel-title'">{{ sub.title }}</p>
                                    <template v-for="leaf in sub.fields" :key="leaf.key || leaf.path">
                                        <div v-if="leaf.kind==='info'" class="df-info-banner">
                                            <span class="df-info-banner-icon">ℹ️</span>
                                            <span>{{ leaf.text }}</span>
                                        </div>
                                        <div v-else-if="leaf.kind==='heading'" class="proto-heading">{{ leaf.text }}</div>
                                        <div v-else-if="leaf.kind==='field'" class="df-question">
                                            <label>{{ leaf.label }}<span v-if="leaf.required" class="df-required">*</span></label>
                                            <div v-if="leaf.hint" class="df-hint">{{ leaf.hint }}</div>
                                            <field-view :field="leaf" @set="onSet"></field-view>
                                            <div v-if="leaf.error" class="df-error">{{ leaf.error }}</div>
                                        </div>
                                    </template>
                                </div>
                                <div v-else-if="sub.kind==='repeater'" class="df-repeater">
                                    <div v-for="nested in sub.items" :key="nested.key" class="df-repeater-item">
                                        <div class="df-repeater-item-header">
                                            <p class="df-repeater-item-title">{{ nested.title }}</p>
                                            <button v-if="nested.canDelete" type="button" class="df-repeater-delete" @click="removeItem(nested.remove, nested.index)">Eliminar</button>
                                        </div>
                                        <template v-for="leaf in nested.fields" :key="leaf.key || leaf.path">
                                            <div v-if="leaf.kind==='info'" class="df-info-banner">
                                                <span class="df-info-banner-icon">ℹ️</span>
                                                <span>{{ leaf.text }}</span>
                                            </div>
                                            <div v-else-if="leaf.kind==='heading'" class="proto-heading">{{ leaf.text }}</div>
                                            <div v-else-if="leaf.kind==='field'" class="df-question">
                                                <label>{{ leaf.label }}<span v-if="leaf.required" class="df-required">*</span></label>
                                                <div v-if="leaf.hint" class="df-hint">{{ leaf.hint }}</div>
                                                <field-view :field="leaf" @set="onSet"></field-view>
                                                <div v-if="leaf.error" class="df-error">{{ leaf.error }}</div>
                                            </div>
                                        </template>
                                    </div>
                                    <button v-if="sub.canAdd" type="button" class="df-repeater-add" @click="addItem(sub.addId)">{{ sub.addLabel }}</button>
                                </div>
                            </template>
                        </div>
                        <button v-if="block.canAdd" type="button" class="df-repeater-add" @click="addItem(block.addId)">{{ block.addLabel }}</button>
                    </div>
                </template>
            </div>
            <div class="df-nav">
                <button type="button" class="df-btn-prev" :disabled="currentIndex===0" @click="prev">Anterior</button>
                <button type="button" class="df-btn-next" @click="next">{{ currentIndex === sections.length - 1 ? 'Finalizar prototipo' : 'Siguiente' }}</button>
            </div>
        </section>
    </div>

    <div v-else>
        <div class="df-info-banner" style="margin-bottom:16px;">
            <span class="df-info-banner-icon">ℹ️</span>
            <span>Respuestas capturadas en esta pantalla. No se envió nada al servidor.</span>
        </div>
        <div v-for="sec in summary" :key="sec.id" class="df-repeater-item proto-summary-sec">
            <h3>{{ sec.title }}</h3>
            <p v-if="!sec.rows.length" class="proto-row-label">Sin respuestas en esta sección.</p>
            <div v-for="(row, idx) in sec.rows" :key="sec.id + '-' + idx" :class="row.head ? 'proto-head' : 'proto-row'">
                <div :class="row.head ? '' : 'proto-row-label'">{{ row.label }}</div>
                <div v-if="!row.head" class="proto-row-value">{{ row.value }}</div>
            </div>
        </div>
        <div class="df-nav" style="border:1px solid #e5e7eb;border-radius:12px;margin-top:12px;">
            <button type="button" class="df-btn-prev" @click="done=false">Volver al formulario</button>
            <span></span>
        </div>
    </div>

    <div v-if="popup" class="df-popup-backdrop" @click.self="popup=false">
        <div class="df-popup">
            <div class="df-popup-icon">!</div>
            <p class="df-popup-title">{{ popupText }}</p>
            <button type="button" class="df-popup-btn" @click="popup=false">Entendido</button>
        </div>
    </div>
</div>
    `,
};
