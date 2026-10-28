/**
 * Genera src/internal/seed/seguimiento_v4_form.sql a partir de las listas
 * del prototipo. No toca el formulario vigente.
 *
 *   node scripts/gen_seguimiento_v4_seed.js
 */
const fs = require('fs');
const vm = require('vm');
const crypto = require('crypto');

const FORM_ID = 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f';
const FORM_DESCRIPTION = 'Seguimiento V.4. Prototipo estable del 23 de octubre de 2026. Convive con el formulario vigente y todavía no ejecuta su lógica al completar.';
const proto = fs.readFileSync('src/frontend/js/prototipo-seguimiento-v4.js', 'utf8');
const start = proto.indexOf('var YN');
const end = proto.indexOf('function opts');
const ctx = vm.createContext({});
vm.runInContext(proto.slice(start, end), ctx);

function id(key) {
    const h = crypto.createHash('sha256').update('seguimiento-v4:' + key).digest();
    const b = Buffer.from(h.subarray(0, 16));
    b[6] = (b[6] & 0x0f) | 0x40;
    b[8] = (b[8] & 0x3f) | 0x80;
    const hex = b.toString('hex');
    return hex.slice(0, 8) + '-' + hex.slice(8, 12) + '-' + hex.slice(12, 16) + '-' + hex.slice(16, 20) + '-' + hex.slice(20);
}
function sql(s) {
    if (s == null) return 'NULL';
    return "'" + String(s).replace(/'/g, "''") + "'";
}

const sections = [];
const repeaters = [];
const questions = [];
const options = [];
const conditions = [];
const byKey = {};

function section(key, name, description) {
    const sec = { key, id: id('section:' + key), name, description, cursor: 1 };
    sections.push(sec);
    return sec;
}
function repeater(sec, key, name, itemName, extra) {
    extra = extra || {};
    const rep = {
        key,
        id: id('repeater:' + key),
        sectionId: sec.id,
        name,
        itemName,
        order: sec.cursor++,
        min: extra.min || 0,
        stateItems: extra.stateItems || null,
        parentKey: extra.parentKey || null,
        cursor: 1,
    };
    repeaters.push(rep);
    return rep;
}
function question(owner, spec) {
    const key = spec.key;
    if (byKey[key]) throw new Error('clave repetida ' + key);
    const inRepeater = owner.sectionId != null;
    const row = {
        key,
        id: id('q:' + key),
        formId: FORM_ID,
        sectionId: inRepeater ? owner.sectionId : owner.id,
        repeaterId: inRepeater ? owner.id : null,
        type: spec.type,
        description: spec.label,
        required: !!spec.required,
        statePath: spec.statePath || null,
        order: owner.cursor++,
    };
    if (row.description.length > 4000) throw new Error('pregunta larga ' + key);
    questions.push(row);
    byKey[key] = row;
    (spec.options || []).forEach(function (pair, i) {
        const label = pair[1];
        const value = pair[0];
        if (label.length > 255) throw new Error('opción >255 ' + key + ' ' + label.length + ' ' + label);
        if (value.length > 255) throw new Error('valor >255 ' + key);
        options.push({
            id: id('opt:' + key + ':' + value),
            questionId: row.id,
            label: label,
            value: value,
            order: i + 1,
        });
    });
    return row;
}
function when(targetKey, triggerKey, value, op) {
    conditions.push({
        id: id('vc:' + targetKey + ':' + triggerKey + ':' + (op || 'EQUALS') + ':' + value),
        targetType: byKey[targetKey] ? 'QUESTION' : 'REPEATER_GROUP',
        targetKey: targetKey,
        triggerKey: triggerKey,
        value: value,
        op: op || 'EQUALS',
        statePath: null,
    });
}
function whenState(targetKey, targetType, path, value) {
    conditions.push({
        id: id('vcstate:' + targetKey + ':' + path + ':' + value),
        targetType: targetType,
        targetKey: targetKey,
        triggerKey: 'risk.respondiente',
        value: value,
        op: 'EQUALS',
        statePath: path,
    });
}

const SECTOR_BARRERA = [
    ['salud', 'Salud'], ['justicia', 'Justicia'], ['proteccion', 'Protección'], ['transversales', 'Transversales'],
];
const GESTION_BARRERA = [
    ['activacion', 'Activación de ruta interinstitucional (Generar oficio)'],
    ['alerta', 'Alerta por barreras'],
];
const GESTION_SEG = [
    ['activacion', 'Activación de ruta'],
    ['alerta', 'Alerta por barreras'],
];
const VEREDA = [['cabecera', 'Cabecera'], ['rural', 'Vereda o zona rural'], ['no', 'No precisa']];
const RESPONDIENTE = [
    ['directo', 'Contacto directo con la víctima'],
    ['institucion', 'Contacto con institución'],
    ['indirecto', 'Contacto indirecto (familiar o persona de confianza)'],
];

function barrierQuestions(rep, prefix) {
    question(rep, { key: prefix + '.sector', type: 'single', label: 'Sector en el que se presentó la barrera', required: true, options: SECTOR_BARRERA });
    question(rep, { key: prefix + '.infoSalud', type: 'info', label: 'Salud: hospitales, IPS, EPS y puestos de salud.' });
    when(prefix + '.infoSalud', prefix + '.sector', 'salud');
    question(rep, { key: prefix + '.tiposSalud', type: 'multiple', label: 'Barreras en salud', options: ctx.SALUD_B });
    when(prefix + '.tiposSalud', prefix + '.sector', 'salud');
    question(rep, { key: prefix + '.instSalud', type: 'dropdown', label: '¿A qué institución acudió?', options: ctx.INST.salud.map(function (n) { return [slug(n), n]; }) });
    when(prefix + '.instSalud', prefix + '.sector', 'salud');
    question(rep, { key: prefix + '.otraTipoSalud', type: 'text', label: 'Otra barrera en sector Salud' });
    when(prefix + '.otraTipoSalud', prefix + '.sector', 'salud');

    question(rep, { key: prefix + '.infoJusticia', type: 'info', label: 'Justicia: fiscalía, policía, comisarías, jueces y medicina legal.' });
    when(prefix + '.infoJusticia', prefix + '.sector', 'justicia');
    question(rep, { key: prefix + '.tiposJusticia', type: 'multiple', label: 'Barreras en justicia', options: ctx.JUSTICIA_B });
    when(prefix + '.tiposJusticia', prefix + '.sector', 'justicia');
    question(rep, { key: prefix + '.instJusticia', type: 'dropdown', label: '¿A qué institución acudió?', options: ctx.INST.justicia.map(function (n) { return [slug(n), n]; }) });
    when(prefix + '.instJusticia', prefix + '.sector', 'justicia');
    question(rep, { key: prefix + '.otraTipoJusticia', type: 'text', label: 'Otra barrera en sector Justicia' });
    when(prefix + '.otraTipoJusticia', prefix + '.sector', 'justicia');

    question(rep, { key: prefix + '.infoProteccion', type: 'info', label: 'Protección: comisarías, ICBF, UNP, defensoría de familia e inspección.' });
    when(prefix + '.infoProteccion', prefix + '.sector', 'proteccion');
    question(rep, { key: prefix + '.tiposProteccion', type: 'multiple', label: 'Barreras en protección', options: ctx.PROTECCION_B });
    when(prefix + '.tiposProteccion', prefix + '.sector', 'proteccion');
    question(rep, { key: prefix + '.instProteccion', type: 'dropdown', label: '¿A qué institución acudió?', options: ctx.INST.proteccion.map(function (n) { return [slug(n), n]; }) });
    when(prefix + '.instProteccion', prefix + '.sector', 'proteccion');
    question(rep, { key: prefix + '.otraTipoProteccion', type: 'text', label: 'Otra barrera en el marco del sector Protección' });
    when(prefix + '.otraTipoProteccion', prefix + '.sector', 'proteccion');

    question(rep, { key: prefix + '.talento', type: 'multiple', label: 'Barreras institucionales y de talento humano', options: ctx.TALENTO });
    when(prefix + '.talento', prefix + '.sector', 'transversales');
    question(rep, { key: prefix + '.eco', type: 'multiple', label: 'Barreras económicas y socioeconómicas', options: ctx.ECO });
    when(prefix + '.eco', prefix + '.sector', 'transversales');
    question(rep, { key: prefix + '.terr', type: 'multiple', label: 'Barreras territoriales y geográficas', options: ctx.TERR_B });
    when(prefix + '.terr', prefix + '.sector', 'transversales');
    question(rep, { key: prefix + '.diferencial', type: 'multiple', label: 'Barreras por ausencia de enfoque diferencial', options: ctx.DIFERENCIAL });
    when(prefix + '.diferencial', prefix + '.sector', 'transversales');

    question(rep, { key: prefix + '.fecha', type: 'date', label: 'Fecha en la que se presentó la barrera (aproximada)', required: true });
    question(rep, { key: prefix + '.funcionario', type: 'text', label: 'Funcionario/a o dependencia donde se presentó la barrera', required: true });
    question(rep, { key: prefix + '.depto', type: 'dropdown', label: 'Departamento', required: true, statePath: 'statesColombia' });
    question(rep, { key: prefix + '.municipio', type: 'dropdown', label: 'Municipio', required: true, statePath: prefix + '.{_entryIndex}.municipios' });
    question(rep, { key: prefix + '.vereda', type: 'dropdown', label: 'Vereda o zona', required: true, options: VEREDA });
    question(rep, { key: prefix + '.descripcion', type: 'text', label: 'Descripción de la barrera (tiempo, modo y lugar)', required: true });
    question(rep, { key: prefix + '.gestion', type: 'multiple', label: 'Gestión de la barrera', required: true, options: GESTION_BARRERA });
    question(rep, { key: prefix + '.perspectiva', type: 'text', label: '¿Cuál es su perspectiva profesional frente a la barrera?', required: true });
    when(prefix + '.perspectiva', prefix + '.gestion', 'alerta', 'INCLUDES');
    question(rep, { key: prefix + '.infoAlerta', type: 'info', label: 'Alerta por barreras notificaría al equipo técnico. Esta versión todavía no envía la notificación.' });
    when(prefix + '.infoAlerta', prefix + '.gestion', 'alerta', 'INCLUDES');
    question(rep, { key: prefix + '.gestionTexto', type: 'text', label: 'Descripción de la gestión realizada con relación a las barreras', required: true });
}
function slug(n) {
    return n.normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_|_$/g, '');
}

const riesgo = section('riesgo', 'Valoración del riesgo', 'Seguimiento al riesgo de feminicidio, homicidio abyecto o fútil');
question(riesgo, { key: 'risk.respondiente', type: 'single', label: 'Indique quién brindó la información durante este seguimiento', required: true, options: RESPONDIENTE });
question(riesgo, { key: 'risk.info', type: 'info', label: 'Los factores son indicadores técnicos de riesgo de feminicidio, homicidio abyecto o fútil y otras violencias graves. No son preguntas directas a la víctima. Registre solo hechos desde el último contacto hasta hoy.' });
question(riesgo, { key: 'risk.hechos', type: 'boolean', label: '¿Se registraron nuevos hechos de violencia desde el último seguimiento?', required: true });
question(riesgo, { key: 'risk.hechosTexto', type: 'text', label: 'Registre los nuevos hechos: tiempo, modo, lugar y relación de la víctima con la persona agresora', required: true });
when('risk.hechosTexto', 'risk.hechos', 'true');
question(riesgo, { key: 'risk.hechosFecha', type: 'date', label: 'Fecha en la que se presentaron los hechos', required: true });
when('risk.hechosFecha', 'risk.hechos', 'true');
question(riesgo, { key: 'risk.protectores', type: 'multiple', label: 'Indique si estos factores protectores se presentan en el caso', required: true, options: ctx.PROTECTORES });
question(riesgo, { key: 'risk.riesgos', type: 'multiple', label: 'Indique si estos factores de riesgo se presentan en el caso', required: true, options: ctx.RIESGOS });
question(riesgo, { key: 'risk.extremos', type: 'multiple', label: 'Si identifica algún factor de riesgo extremo, remita al equipo de riesgo. Verifique si ya hay articulación interinstitucional y priorice la protección inmediata', required: true, options: ctx.EXTREMOS });
question(riesgo, { key: 'risk.descripcion', type: 'text', label: 'Describa la situación de riesgo actual y cómo los factores de riesgo y protectores influyen en que este se mantenga o modifique', required: true });

const medidas = section('medidas', 'Seguimiento a medidas de emergencia', 'Solo cuando el caso ya tiene medidas de emergencia, en cualquier estado');
whenState('medidas', 'SECTION', 'hasMedidas', 'true');
question(medidas, { key: 'med.contribucion', type: 'text', label: '¿De qué manera la medida ha contribuido a la gestión de su seguridad y la reducción del riesgo desde el último seguimiento?', required: true });
question(medidas, { key: 'med.sentir', type: 'text', label: '¿Cómo se ha sentido con las medidas de emergencia recibidas y qué cambios han generado en su situación actual?', required: true });
question(medidas, { key: 'med.beneficiarias', type: 'text', label: '¿Cómo perciben las personas beneficiarias las medidas otorgadas y de qué manera estas han impactado su bienestar o seguridad? (Si las personas beneficiarias son menores de edad u otras personas dependientes, la valoración puede realizarse a partir del reporte de la participante).', required: true });
question(medidas, { key: 'med.alternativas', type: 'text', label: 'Además de las medidas de emergencia, ¿qué alternativas o estrategias está construyendo para fortalecer su protección y bienestar a mediano plazo?', required: true });
question(medidas, { key: 'med.infoCompromisos', type: 'info', label: 'Registre los acuerdos, acciones o gestiones que realizará la persona o los beneficiarios' });
question(medidas, { key: 'med.compromisos', type: 'text', label: 'Compromisos', required: true });
question(medidas, { key: 'med.infoObservaciones', type: 'info', label: 'Registre información relevante identificada durante el seguimiento que no haya sido consignada en otros campos y que aporte a la comprensión del caso.' });
question(medidas, { key: 'med.observaciones', type: 'text', label: 'Observaciones', required: true });

const seg = section('segEntidades', 'Seguimiento a entidades', 'Solo si el caso ya tiene entidades');
whenState('segEntidades', 'SECTION', 'hasEntidades', 'true');
question(seg, { key: 'seg.info', type: 'info', label: 'El listado lo trae el sistema. Este apartado integra las acciones nuevas de la entidad, las barreras nuevas y el seguimiento a las barreras que ya tenía.' });
question(seg, { key: 'seg.entidadesSel', type: 'multiple', label: 'Selección de la entidad a actualizar', required: true, statePath: 'currentEntities' });
const segRep = repeater(seg, 'segEntidades', 'Entidades del caso', 'Entidad', { stateItems: 'currentEntities' });
question(segRep, { key: 'seg.completadas', type: 'multiple', label: 'Acciones institucionales que ya se completaron', statePath: 'currentEntities.{_entryIndex}.obligaciones' });
question(segRep, { key: 'seg.infoBarreras', type: 'info', label: 'Seguimiento a barreras de esta entidad. Las barreras ya asociadas las trae el sistema.' });
whenState('seg.infoBarreras', 'QUESTION', 'currentEntities.{_entryIndex}.hasBarriers', 'true');
const segBarrerasActivas = repeater(seg, 'segBarrerasActivas', 'Seguimiento a barreras de esta entidad', 'Barrera', {
    parentKey: 'segEntidades',
    stateItems: 'currentEntities.{_entryIndex}.barreras',
});
whenState('segBarrerasActivas', 'REPEATER_GROUP', 'currentEntities.{_entryIndex}.hasBarriers', 'true');
question(segBarrerasActivas, { key: 'seg.barreraActiva.persisten', type: 'boolean', label: '¿Las barreras persisten?', required: true });
question(segBarrerasActivas, { key: 'seg.barreraActiva.cierre', type: 'boolean', label: '¿Se realiza cierre de la barrera?', required: true });
question(segBarrerasActivas, { key: 'seg.barreraActiva.motivoCierre', type: 'single', label: 'Motivo del cierre', required: true, options: ctx.CIERRE_MOTIVOS });
when('seg.barreraActiva.motivoCierre', 'seg.barreraActiva.cierre', 'true');
question(segBarrerasActivas, { key: 'seg.barreraActiva.gestion', type: 'single', label: 'Gestión de la barrera por parte de Salvia', options: GESTION_SEG });
question(segBarrerasActivas, { key: 'seg.barreraActiva.infoAlerta', type: 'info', label: 'Alerta por barreras: al guardar de verdad, esto notificaría al equipo técnico. Esta versión no envía la notificación.' });
when('seg.barreraActiva.infoAlerta', 'seg.barreraActiva.gestion', 'alerta');
question(segRep, { key: 'seg.infoPend', type: 'info', label: 'Pendientes que trae el sistema para esta entidad.' });
question(segRep, { key: 'seg.pendientes', type: 'multiple', label: 'Acciones institucionales pendientes', statePath: 'currentEntities.{_entryIndex}.obligaciones' });
question(segRep, { key: 'seg.otraPendiente', type: 'text', label: 'Otra acción pendiente', required: true });
when('seg.otraPendiente', 'seg.pendientes', 'otra', 'INCLUDES');
question(segRep, { key: 'seg.actualizacion', type: 'text', label: 'Describa cualquier actualización de la ruta desde el último seguimiento', required: true });
question(segRep, { key: 'seg.infoNueva', type: 'info', label: 'Registro de barreras y activación de ruta' });
question(segRep, { key: 'seg.nuevaBarrera', type: 'boolean', label: '¿Se identifica alguna nueva barrera en este seguimiento con relación a esta entidad?', required: true });
const segBarreras = repeater(seg, 'segBarreras', 'Barreras nuevas de esta entidad', 'Barrera', { parentKey: 'segEntidades' });
when('segBarreras', 'seg.nuevaBarrera', 'true');
barrierQuestions(segBarreras, 'seg.barrera');
question(segRep, { key: 'seg.acuerdo', type: 'boolean', label: '¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes?', required: true });
question(segRep, { key: 'seg.debida', type: 'boolean', label: '¿Salvia debe activar rutas en el marco de la debida diligencia?', required: true });
when('seg.debida', 'seg.acuerdo', 'true');
question(segRep, { key: 'seg.canales', type: 'multiple', label: 'Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención', required: true, options: ctx.CANALES });
when('seg.canales', 'seg.acuerdo', 'true');
when('seg.canales', 'seg.debida', 'true');
question(segRep, { key: 'seg.gestiones', type: 'text', label: 'Describa las gestiones realizadas para la activación de rutas y gestión de barreras en esta entidad' });
when('seg.gestiones', 'seg.acuerdo', 'true');
when('seg.gestiones', 'seg.debida', 'true');

const ident = section('idEntidades', 'Identificación de entidades', 'Entidades a las que la persona ya acudió, o a las que otras activaron ruta');
question(ident, { key: 'ident.conocimiento', type: 'boolean', label: '¿Ha puesto en conocimiento de alguna entidad o institución los hechos de violencia reportados?', required: true });
const identRep = repeater(ident, 'identEntidades', 'Entidades a las que ya acudió', 'Entidad');
when('identEntidades', 'ident.conocimiento', 'true');
question(identRep, { key: 'ident.sector', type: 'dropdown', label: 'Sector de la entidad', required: true, options: ctx.SECTORES });
question(identRep, { key: 'ident.depto', type: 'dropdown', label: 'Departamento', required: true, statePath: 'statesColombia' });
question(identRep, { key: 'ident.municipio', type: 'dropdown', label: 'Municipio', required: true, statePath: 'identEntidades.{_entryIndex}.municipios' });
question(identRep, { key: 'ident.ciudad', type: 'dropdown', label: 'Ciudad', required: true, statePath: 'identEntidades.{_entryIndex}.ciudades' });
question(identRep, { key: 'ident.entidad', type: 'dropdown', label: 'Entidad', required: true, statePath: 'identEntidades.{_entryIndex}.entidades' });
question(identRep, { key: 'ident.otraEntidad', type: 'text', label: 'Otra ¿cuál?', required: true });
when('ident.otraEntidad', 'ident.entidad', 'otra');
question(identRep, { key: 'ident.radicado', type: 'text', label: 'Número de radicado (NUNC, medida de protección). Escriba tipo y número', required: true });
question(identRep, { key: 'ident.fecha', type: 'date', label: 'Fecha en la que reportó a la entidad', required: true });
question(identRep, { key: 'ident.completadas', type: 'multiple', label: 'Acciones institucionales completadas', required: true, statePath: 'identEntidades.{_entryIndex}.obligaciones' });
question(identRep, { key: 'ident.pendientes', type: 'multiple', label: 'Acciones institucionales pendientes', required: true, statePath: 'identEntidades.{_entryIndex}.obligaciones' });
question(identRep, { key: 'ident.otraRuta', type: 'text', label: 'Otra información asociada a la Ruta', required: true });
question(identRep, { key: 'ident.nuevaBarrera', type: 'boolean', label: '¿Se identifica alguna nueva barrera con relación a esta entidad?', required: true });
const identBarreras = repeater(ident, 'identBarreras', 'Barreras nuevas de esta entidad', 'Barrera', { parentKey: 'identEntidades' });
when('identBarreras', 'ident.nuevaBarrera', 'true');
barrierQuestions(identBarreras, 'ident.barrera');
question(identRep, { key: 'ident.acuerdo', type: 'boolean', label: '¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes?', required: true });
question(identRep, { key: 'ident.debida', type: 'boolean', label: '¿Salvia debe activar rutas en el marco de la debida diligencia?', required: true });
when('ident.debida', 'ident.acuerdo', 'true');
question(identRep, { key: 'ident.canales', type: 'multiple', label: 'Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención', required: true, options: ctx.CANALES });
when('ident.canales', 'ident.acuerdo', 'true');
when('ident.canales', 'ident.debida', 'true');
question(identRep, { key: 'ident.gestiones', type: 'text', label: '¿Qué gestiones se realizaron desde Salvia para la activación de ruta?', required: true });
when('ident.gestiones', 'ident.acuerdo', 'true');
when('ident.gestiones', 'ident.debida', 'true');

const plan = section('plan', 'Seguimiento al plan de acción', 'Remisiones, acompañamiento psicosocial y descripción');
question(plan, { key: 'plan.remision', type: 'boolean', label: 'Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA', required: true });
question(plan, { key: 'plan.equipos', type: 'multiple', label: 'Equipos', required: true, options: ctx.EQUIPOS });
when('plan.equipos', 'plan.remision', 'true');
question(plan, { key: 'plan.infoExclusion', type: 'info', label: 'Remitió a medidas de emergencia. Atención psicosocial y psicosocial OSIGNH quedan deshabilitadas.' });
when('plan.infoExclusion', 'plan.remision', 'true');
when('plan.infoExclusion', 'plan.equipos', 'medidas', 'INCLUDES');
question(plan, { key: 'plan.medidasTipos', type: 'multiple', label: 'Medidas de emergencia a las que remite', required: true, options: ctx.MEDIDAS_TIPO });
when('plan.medidasTipos', 'plan.remision', 'true');
when('plan.medidasTipos', 'plan.equipos', 'medidas', 'INCLUDES');
question(plan, { key: 'plan.psico', type: 'multiple', label: 'Criterios de remisión a psicosocial', required: true, options: ctx.PSICO_CRIT });
when('plan.psico', 'plan.remision', 'true');
when('plan.psico', 'plan.equipos', 'psico', 'INCLUDES');
question(plan, { key: 'plan.osignh', type: 'multiple', label: 'Criterios de psicosocial para personas OSIGNH', required: true, options: [
    ['obligatorio', 'Ser OSIGNH, mayor de 18 años, víctima de VpP, en post-emergencia, con voluntariedad y disponibilidad. Obligatorio'],
    ['eps', 'No aplica si ya recibe atención psicosocial por EPS, salvo barreras por orientación sexual o identidad de género'],
] });
when('plan.osignh', 'plan.remision', 'true');
when('plan.osignh', 'plan.equipos', 'osignh', 'INCLUDES');
question(plan, { key: 'plan.infoEps', type: 'info', label: 'Si ya recibe atención por EPS y no hay barreras por orientación sexual o identidad de género, la remisión a OSIGNH no aplica.' });
when('plan.infoEps', 'plan.remision', 'true');
when('plan.infoEps', 'plan.osignh', 'eps', 'INCLUDES');
question(plan, { key: 'plan.infoHombres', type: 'info', label: 'Criterio de Atención Hombres: hombres diversos victimarios o posibles victimarios.' });
when('plan.infoHombres', 'plan.remision', 'true');
when('plan.infoHombres', 'plan.equipos', 'hombres', 'INCLUDES');
question(plan, { key: 'plan.discapacidad', type: 'multiple', label: 'Equipo Discapacidad', required: true, options: [
    ['lsc', 'Apoyo comunicativo en LSC'],
    ['enfoque', 'Apoyo con enfoque de discapacidad (ley de capacidad legal)'],
] });
when('plan.discapacidad', 'plan.remision', 'true');
when('plan.discapacidad', 'plan.equipos', 'discapacidad', 'INCLUDES');
question(plan, { key: 'plan.estabilizacion', type: 'multiple', label: 'Estabilización', required: true, options: [
    ['ingresos', 'No tiene ninguna forma de ingresos o tiene ingresos de forma inestable'],
    ['vivienda', 'La persona reside actualmente en vivienda temporal o familiar'],
    ['cargo', 'La persona tiene personas a cargo que limitan su tiempo para trabajar o estudiar'],
    ['depende', 'Depende económicamente de la persona agresora'],
] });
when('plan.estabilizacion', 'plan.remision', 'true');
when('plan.estabilizacion', 'plan.equipos', 'estabilizacion', 'INCLUDES');
question(plan, { key: 'plan.elementos', type: 'multiple', label: 'Seleccione los elementos psicosociales aplicados en este seguimiento', required: true, options: ctx.ELEMENTOS });
question(plan, { key: 'plan.acciones', type: 'text', label: 'Describa las acciones realizadas en el seguimiento', required: true });
question(plan, { key: 'plan.ajustes', type: 'text', label: 'Describa los ajustes que realiza al plan de acción de Salvia', required: true });
question(plan, { key: 'plan.equipoReasigna', type: 'dropdown', label: 'Equipo al que reasigna el caso', required: true, options: [
    ['riesgo', 'Seguimiento de Riesgo'],
    ['general', 'Seguimiento General'],
    ['feminicidio', 'Abordaje Integral del Feminicidio'],
    ['ninguno', 'Ninguno'],
] });

const conclusion = section('conclusion', 'Conclusión del ciclo de atención', 'Criterio, motivo y acción institucional');
question(conclusion, { key: 'con.concluye', type: 'boolean', label: '¿Se concluye el ciclo de atención?', required: true });
question(conclusion, { key: 'con.info', type: 'info', label: 'Registre de forma clara y sustentada las razones.' });
when('con.info', 'con.concluye', 'true');
question(conclusion, { key: 'con.motivo', type: 'single', label: 'Motivo de la conclusión del ciclo de atención', required: true, options: ctx.CONCLUSION_MOTIVOS });
when('con.motivo', 'con.concluye', 'true');
question(conclusion, { key: 'con.otro', type: 'text', label: 'Otro ¿cuál?', required: true });
when('con.otro', 'con.motivo', 'otro');
question(conclusion, { key: 'con.infoDesc', type: 'info', label: 'Registrar de forma clara la descripción del criterio y la orientación brindada sobre continuidad de derechos y rutas disponibles: intentos de contacto (fechas, medios, número de intentos y tiempo transcurrido sin respuesta); solicitud de la persona (qué solicita, razones expresadas); avances y acciones completadas con evidencia y los criterios que permiten concluir que el plan de atención se cumplió; no corresponsabilidad (explique la ausencia de gestiones o de compromiso con el acompañamiento); no clasificación como VBG/VpP (explicar por qué no se clasifica).' });
when('con.infoDesc', 'con.concluye', 'true');
question(conclusion, { key: 'con.descripcion', type: 'text', label: 'Descripción del criterio y de la orientación brindada sobre continuidad de derechos y rutas disponibles', required: true });
when('con.descripcion', 'con.concluye', 'true');
question(conclusion, { key: 'con.accion', type: 'boolean', label: '¿Se realizó alguna acción institucional en la conclusión del caso?', required: true });
when('con.accion', 'con.concluye', 'true');
question(conclusion, { key: 'con.entidad', type: 'dropdown', label: 'Seleccione la entidad con la que realizará acción institucional', required: true, statePath: 'entidadesConocidas' });
when('con.entidad', 'con.accion', 'true');
question(conclusion, { key: 'con.otraEntidad', type: 'text', label: 'Otra ¿cuál?', required: true });
when('con.otraEntidad', 'con.entidad', 'otra');
question(conclusion, { key: 'con.canales', type: 'multiple', label: 'Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención', required: true, options: ctx.CANALES });
when('con.canales', 'con.accion', 'true');

const repByKey = {};
repeaters.forEach(function (r) { repByKey[r.key] = r; });

function targetId(cond) {
    if (cond.targetType === 'SECTION') {
        const sec = sections.find(function (s) { return s.key === cond.targetKey; });
        if (!sec) throw new Error('sección ' + cond.targetKey);
        return sec.id;
    }
    if (cond.targetType === 'REPEATER_GROUP') return repByKey[cond.targetKey].id;
    return byKey[cond.targetKey].id;
}

const lines = [];
lines.push('-- Formulario de seguimiento V.4. Id fijo: ' + FORM_ID);
lines.push('-- No modifica el formulario vigente 2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff.');
lines.push('-- Una sola sentencia (bloque DO) para db.Exec.');
lines.push('-- Si el formulario ya tiene envíos, aborta y no borra nada.');
lines.push('-- Generado por scripts/gen_seguimiento_v4_seed.js. No editar a mano.');
lines.push('DO $$');
lines.push('DECLARE');
lines.push('  n_envios integer;');
lines.push('BEGIN');
lines.push('');
lines.push('SELECT COUNT(*) INTO n_envios FROM salvia.form_submission');
lines.push('  WHERE form_id = ' + sql(FORM_ID) + ';');
lines.push('IF n_envios > 0 THEN');
lines.push("  RAISE EXCEPTION 'el formulario de seguimiento V.4 ya tiene envíos y no se reemplaza';");
lines.push('END IF;');
lines.push('');
lines.push('DELETE FROM salvia.visibility_condition');
lines.push('  WHERE trigger_question_id IN (SELECT id FROM salvia.question WHERE form_id = ' + sql(FORM_ID) + ')');
lines.push('     OR target_id IN (SELECT id FROM salvia.question WHERE form_id = ' + sql(FORM_ID) + ')');
lines.push('     OR target_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = ' + sql(FORM_ID) + ')');
lines.push('     OR target_id IN (');
lines.push('          SELECT id FROM salvia.repeater_group');
lines.push('           WHERE form_section_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = ' + sql(FORM_ID) + ')');
lines.push('        );');
lines.push('DELETE FROM salvia.answer');
lines.push('  WHERE question_id IN (SELECT id FROM salvia.question WHERE form_id = ' + sql(FORM_ID) + ');');
lines.push('DELETE FROM salvia.repeater_entry');
lines.push('  WHERE repeater_group_id IN (');
lines.push('        SELECT id FROM salvia.repeater_group');
lines.push('         WHERE form_section_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = ' + sql(FORM_ID) + ')');
lines.push('      );');
lines.push('DELETE FROM salvia.option');
lines.push('  WHERE question_id IN (SELECT id FROM salvia.question WHERE form_id = ' + sql(FORM_ID) + ');');
lines.push('DELETE FROM salvia.question WHERE form_id = ' + sql(FORM_ID) + ';');
lines.push('DELETE FROM salvia.repeater_group');
lines.push('  WHERE form_section_id IN (SELECT id::text FROM salvia.form_section WHERE form_id = ' + sql(FORM_ID) + ');');
lines.push('DELETE FROM salvia.form_section WHERE form_id = ' + sql(FORM_ID) + ';');
lines.push('');
lines.push('INSERT INTO salvia.form (id, name, description, status) VALUES (');
lines.push('  ' + sql(FORM_ID) + ',');
lines.push("  'Seguimiento V.4',");
lines.push('  ' + sql(FORM_DESCRIPTION) + ',');
lines.push("  'active'");
lines.push(') ON CONFLICT (id) DO NOTHING;');
lines.push('');
lines.push('INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES');
lines.push(sections.map(function (s, i) {
    return '  (' + [sql(s.id), sql(FORM_ID), sql(s.name), sql(s.description), i + 1].join(', ') + ')';
}).join(',\n') + '\nON CONFLICT (id) DO NOTHING;');
lines.push('');
lines.push('INSERT INTO salvia.repeater_group (id, form_section_id, name, item_name, "order", min_repetitions, state_items, parent_repeater_group_id) VALUES');
lines.push(repeaters.map(function (r) {
    const parent = r.parentKey ? sql(repByKey[r.parentKey].id) : 'NULL';
    return '  (' + [sql(r.id), sql(r.sectionId), sql(r.name), sql(r.itemName), r.order, r.min, sql(r.stateItems), parent].join(', ') + ')';
}).join(',\n') + '\nON CONFLICT (id) DO NOTHING;');
lines.push('');
lines.push('INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order") VALUES');
lines.push(questions.map(function (q) {
    return '  (' + [sql(q.id), sql(q.formId), sql(q.sectionId), sql(q.repeaterId), sql(q.type), sql(q.description), q.required ? 'true' : 'false', sql(q.statePath), q.order].join(', ') + ')';
}).join(',\n') + '\nON CONFLICT (id) DO NOTHING;');
lines.push('');
lines.push('INSERT INTO salvia.option (id, question_id, label, value, "order") VALUES');
lines.push(options.map(function (o) {
    return '  (' + [sql(o.id), sql(o.questionId), sql(o.label), sql(o.value), o.order].join(', ') + ')';
}).join(',\n') + '\nON CONFLICT (id) DO NOTHING;');
lines.push('');
lines.push('INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator) VALUES');
lines.push(conditions.map(function (c) {
    const trigger = byKey[c.triggerKey];
    if (!trigger) throw new Error('trigger ' + c.triggerKey);
    return '  (' + [sql(c.id), sql(c.targetType), sql(targetId(c)), sql(trigger.id), sql(c.statePath), sql(c.value), sql(c.op)].join(', ') + ')';
}).join(',\n') + '\nON CONFLICT (id) DO NOTHING;');
lines.push('');
lines.push('UPDATE salvia.form SET deleted_at = NULL, status = ' + sql('active') + ',');
lines.push('  name = ' + sql('Seguimiento V.4') + ',');
lines.push('  description = ' + sql(FORM_DESCRIPTION));
lines.push('  WHERE id = ' + sql(FORM_ID) + ';');
lines.push('');
lines.push('END $$;');
lines.push('');

const seen = {};
[sections, repeaters, questions, options, conditions].forEach(function (list) {
    list.forEach(function (row) {
        if (seen[row.id]) throw new Error('id repetido ' + row.id + ' ' + (row.key || row.label || ''));
        seen[row.id] = true;
    });
});

const out = 'src/internal/seed/seguimiento_v4_form.sql';
fs.writeFileSync(out, lines.join('\n'));

const medidasSec = sections.find(function (s) { return s.key === 'medidas'; });
const medidasQuestions = questions.filter(function (q) { return q.sectionId === medidasSec.id; });
const medidasConditions = conditions.filter(function (c) { return c.targetKey === 'medidas'; });
const add = [];
add.push('-- Agrega Seguimiento a medidas de emergencia sin borrar envíos ni el resto del formulario.');
add.push('-- Si la sección ya está, no hace nada. Generado por scripts/gen_seguimiento_v4_seed.js. No editar a mano.');
add.push('DO $$');
add.push('BEGIN');
add.push('  IF NOT EXISTS (');
add.push('    SELECT 1 FROM salvia.form');
add.push('     WHERE id = ' + sql(FORM_ID) + ' AND deleted_at IS NULL');
add.push('  ) THEN');
add.push('    RETURN;');
add.push('  END IF;');
add.push('');
add.push('  IF EXISTS (');
add.push('    SELECT 1 FROM salvia.form_section WHERE id = ' + sql(medidasSec.id));
add.push('  ) THEN');
add.push('    RETURN;');
add.push('  END IF;');
add.push('');
add.push('  UPDATE salvia.form_section');
add.push('     SET "order" = "order" + 1');
add.push('   WHERE form_id = ' + sql(FORM_ID));
add.push('     AND "order" >= 2');
add.push('     AND deleted_at IS NULL;');
add.push('');
add.push('  INSERT INTO salvia.form_section (id, form_id, name, description, "order") VALUES');
add.push('    (' + [sql(medidasSec.id), sql(FORM_ID), sql(medidasSec.name), sql(medidasSec.description), '2'].join(', ') + ');');
add.push('');
add.push('  INSERT INTO salvia.question (id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order") VALUES');
add.push(medidasQuestions.map(function (q) {
    return '    (' + [sql(q.id), sql(q.formId), sql(q.sectionId), sql(q.repeaterId), sql(q.type), sql(q.description), q.required ? 'true' : 'false', sql(q.statePath), q.order].join(', ') + ')';
}).join(',\n') + ';');
add.push('');
add.push('  INSERT INTO salvia.visibility_condition (id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator) VALUES');
add.push(medidasConditions.map(function (c) {
    const trigger = byKey[c.triggerKey];
    if (!trigger) throw new Error('trigger ' + c.triggerKey);
    return '    (' + [sql(c.id), sql(c.targetType), sql(targetId(c)), sql(trigger.id), sql(c.statePath), sql(c.value), sql(c.op)].join(', ') + ')';
}).join(',\n') + ';');
add.push('END $$;');
add.push('');
const addOut = 'src/internal/seed/seguimiento_v4_medidas.sql';
fs.writeFileSync(addOut, add.join('\n'));

console.log(FORM_ID);
console.log('sections', sections.length);
console.log('repeaters', repeaters.length);
console.log('questions', questions.length);
console.log('options', options.length);
console.log('conditions', conditions.length);
console.log('bytes', fs.statSync(out).size);
console.log('medidas', fs.statSync(addOut).size);
