/**
 * remisiones-medidas-component
 * Listado de remisiones de Medidas de Emergencia con data quemada.
 * La prop `vista` decide cards, filtros, recorte y acciones.
 * `puedeReasignar` la encienden los roles sm y sv.
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[remisiones-medidas-component] Cargue este script después de var app/home.');
        return;
    }

    var AGENTES = [
        { id: 'am-1', nombre: 'Pedro Ramírez' },
        { id: 'am-2', nombre: 'Lucía Mendoza' },
        { id: 'am-3', nombre: 'Jorge Castillo' },
    ];

    var PROFESIONALES = [
        'Camila Andrea Rojas Pinzón',
        'Esteban David Linares Peña',
        'María Fernanda Oquendo Ibarra',
    ];

    var RESULTADOS_REVISION = [
        'Aprobado para caracterización ME',
        'Devuelto para ajuste',
        'Devuelto para acompañamiento',
    ];

    var RESULTADO_A_ESTADO = {
        'Aprobado para caracterización ME': 'Por asignar',
        'Devuelto para ajuste': 'Devuelta para justificación',
        'Devuelto para acompañamiento': 'Devuelta a seguimiento',
    };

    var ESTADOS = [
        'Riesgo inminente: por decidir vía',
        'Riesgo inminente: por escalar a Cruz Roja',
        'Riesgo inminente: por entregar medida',
        'Riesgo inminente: por confirmar seguridad',
        'Pendiente de Revisión',
        'Por asignar',
        'Asignada',
        'En gestión documental',
        'Caracterizado',
        'Por revisar trazabilidad',
        'Por cotizar',
        'Por aprobar cotizaciones',
        'Aprobadas',
        'En Seguimiento',
        'Devuelta para justificación',
        'Devuelta a seguimiento',
    ];

    var ESTADO_MOD = {
        'Riesgo inminente: por decidir vía': 'decidir',
        'Riesgo inminente: por escalar a Cruz Roja': 'escalar',
        'Riesgo inminente: por entregar medida': 'entregar',
        'Riesgo inminente: por confirmar seguridad': 'confirmar',
        'Pendiente de Revisión': 'revision',
        'Por asignar': 'por-asignar',
        'Asignada': 'asignada',
        'En gestión documental': 'gestion',
        'Caracterizado': 'caracterizado',
        'Por revisar trazabilidad': 'trazabilidad',
        'Por cotizar': 'cotizar',
        'Por aprobar cotizaciones': 'cotizaciones',
        'Aprobadas': 'aprobadas',
        'En Seguimiento': 'seguimiento',
        'Devuelta para justificación': 'justificacion',
        'Devuelta a seguimiento': 'devuelta',
    };

    var ESTADOS_RIESGO_INMINENTE = [
        'Riesgo inminente: por decidir vía',
        'Riesgo inminente: por escalar a Cruz Roja',
        'Riesgo inminente: por entregar medida',
        'Riesgo inminente: por confirmar seguridad',
    ];

    var ESTADOS_AGENTE = [
        'Riesgo inminente: por confirmar seguridad',
        'Asignada',
        'En gestión documental',
        'Aprobadas',
        'En Seguimiento',
    ];

    var ESTADOS_CONVENIO = [
        'Por revisar trazabilidad',
        'Por aprobar cotizaciones',
        'Riesgo inminente: por escalar a Cruz Roja',
    ];

    var ESTADOS_CRUZ_ROJA = [
        'Aprobadas',
        'Por cotizar',
        'Riesgo inminente: por entregar medida',
    ];

    function agente(id) {
        for (var i = 0; i < AGENTES.length; i++) {
            if (AGENTES[i].id === id) {
                return { id: AGENTES[i].id, nombre: AGENTES[i].nombre };
            }
        }
        return null;
    }

    function remision(data) {
        return data;
    }

    var SEED = [
        remision({
            id: 'MED-001',
            caso: { id: 'SAL-001', nombre: 'Ana Torres Velásquez', docNumero: '1098234567', nivelRiesgo: 'Crítico', municipio: 'Cali' },
            remitidoPor: { nombre: 'María González', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-05',
            estado: 'Caracterizado',
            medidasSolicitadas: ['Alojamiento temporal', 'Kit de alimentación'],
            agenteAsignado: agente('am-1'),
            motivoRemision: 'Víctima reportó amenaza directa del agresor en su domicilio. Requiere salida urgente y albergue provisional.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-002',
            caso: { id: 'SAL-002', nombre: 'Carmen Lucía Ríos', docNumero: '1023456789', nivelRiesgo: 'Alto', municipio: 'Palmira' },
            remitidoPor: { nombre: 'Juan Pérez', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-03',
            estado: 'Aprobadas',
            medidasSolicitadas: ['Transporte de emergencia', 'Kit de vestuario'],
            agenteAsignado: agente('am-2'),
            motivoRemision: 'Familia en zona rural sin acceso a servicios. Necesita transporte para traslado seguro.',
            urgencia: 'Normal',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-003',
            caso: { id: 'SAL-003', nombre: 'Yolanda Patiño Cruz', docNumero: '1156789012', nivelRiesgo: 'Crítico', municipio: 'Buenaventura' },
            remitidoPor: { nombre: 'Sandra Mora', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-01',
            estado: 'En Seguimiento',
            medidasSolicitadas: ['Alojamiento temporal', 'Kit de alimentación', 'Atención en salud'],
            agenteAsignado: agente('am-3'),
            motivoRemision: 'Caso de violencia extrema con lesiones físicas. Necesita atención médica y alojamiento de emergencia.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-004',
            caso: { id: 'SAL-004', nombre: 'Isabel Moreno Vega', docNumero: '1298765432', nivelRiesgo: 'Alto', municipio: 'Cali' },
            remitidoPor: { nombre: 'Luis Hernández', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-10',
            estado: 'Por asignar',
            medidasSolicitadas: ['Kit de alimentación', 'Kit de vestuario'],
            agenteAsignado: null,
            motivoRemision: 'Víctima desplazada sin recursos económicos. Requiere ayuda humanitaria de emergencia.',
            urgencia: 'Normal',
            revisionTecnica: {
                profesional: 'Camila Andrea Rojas Pinzón',
                concepto: 'El caso cumple con los criterios de riesgo de feminicidio y contacto reciente con el agresor. La información registrada es suficiente para continuar con la caracterización.',
                resultado: 'Aprobado para caracterización ME',
                fecha: '2026-06-11',
            },
        }),
        remision({
            id: 'MED-005',
            caso: { id: 'SAL-005', nombre: 'Rosa Martínez León', docNumero: '1367890123', nivelRiesgo: 'Alto', municipio: 'Cali' },
            remitidoPor: { nombre: 'María González', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-08',
            estado: 'Devuelta para justificación',
            medidasSolicitadas: ['Medida de protección urgente'],
            agenteAsignado: agente('am-1'),
            motivoRemision: 'Solicita medida de protección policial. Situación de riesgo documentada.',
            urgencia: 'Riesgo inminente',
            motivoRechazo: 'La información registrada no permite evidenciar con claridad el riesgo de feminicidio ni el contacto reciente con el agresor. Se requiere ajustar y complementar el registro antes de continuar con el proceso.',
            revisionTecnica: {
                profesional: 'Esteban David Linares Peña',
                concepto: 'No se evidencia con claridad el riesgo de feminicidio ni el contacto reciente con el agresor. La medida solicitada (protección policial) no corresponde a las medidas de emergencia del programa.',
                resultado: 'Devuelto para ajuste',
                fecha: '2026-06-09',
            },
        }),
        remision({
            id: 'MED-006',
            caso: { id: 'SAL-006', nombre: 'Claudia Suárez Reyes', docNumero: '1445678901', nivelRiesgo: 'Crítico', municipio: 'Buenaventura' },
            remitidoPor: { nombre: 'Jorge Ramírez', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-12',
            estado: 'Asignada',
            medidasSolicitadas: ['Alojamiento temporal', 'Transporte de emergencia'],
            agenteAsignado: agente('am-3'),
            motivoRemision: 'Víctima con orden de alejamiento activa que el agresor incumplió. Necesita reubicación inmediata.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-007',
            caso: { id: 'SAL-007', nombre: 'Paola Andrea Gutiérrez', docNumero: '1501122334', nivelRiesgo: 'Crítico', municipio: 'Cali' },
            remitidoPor: { nombre: 'Laura Jiménez', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-14',
            estado: 'Riesgo inminente: por decidir vía',
            medidasSolicitadas: ['Alojamiento temporal', 'Transporte de emergencia'],
            agenteAsignado: null,
            motivoRemision: 'Agente de seguimiento identifica contacto reciente con el agresor y riesgo de feminicidio. Remite para gestión de medidas de emergencia.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-008',
            caso: { id: 'SAL-008', nombre: 'Natalia Ximena Cárdenas', docNumero: '1502233445', nivelRiesgo: 'Alto', municipio: 'Palmira' },
            remitidoPor: { nombre: 'Equipo de Seguimiento', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-13',
            estado: 'Pendiente de Revisión',
            medidasSolicitadas: ['Kit de alimentación'],
            agenteAsignado: null,
            motivoRemision: 'La víctima manifestó que tiene medidas de protección activas, pero el equipo de seguimiento considera que no han sido efectivas para reducir el riesgo.',
            urgencia: 'Normal',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-009',
            caso: { id: 'SAL-009', nombre: 'Diana Carolina Ospina', docNumero: '1503344556', nivelRiesgo: 'Medio', municipio: 'Buenaventura' },
            remitidoPor: { nombre: 'Equipo de Seguimiento', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-09',
            estado: 'Devuelta a seguimiento',
            medidasSolicitadas: ['Kit de alimentación', 'Kit de vestuario'],
            agenteAsignado: null,
            motivoRemision: 'Equipo de Seguimiento remite el caso considerando que requiere medidas de emergencia inmediatas.',
            urgencia: 'Normal',
            revisionTecnica: {
                profesional: 'María Fernanda Oquendo Ibarra',
                concepto: 'El caso no cuenta, en este momento, con criterios suficientes para continuar hacia Medidas de Emergencia. Se requiere continuar el acompañamiento desde Seguimiento y reevaluar si se configuran los criterios más adelante.',
                resultado: 'Devuelto para acompañamiento',
                fecha: '2026-06-10',
            },
        }),
        remision({
            id: 'MED-010',
            caso: { id: 'SAL-010', nombre: 'Katherine Lozano Mesa', docNumero: '1504455667', nivelRiesgo: 'Alto', municipio: 'Cali' },
            remitidoPor: { nombre: 'María González', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-07',
            estado: 'Por revisar trazabilidad',
            medidasSolicitadas: ['Alojamiento temporal', 'Kit de alimentación'],
            agenteAsignado: agente('am-2'),
            motivoRemision: 'Víctima con antecedentes de agresión física reciente. Requiere alojamiento y alimentación de emergencia.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: {
                profesional: 'Esteban David Linares Peña',
                concepto: 'El caso cumple con los criterios de riesgo de feminicidio y contacto reciente con el agresor.',
                resultado: 'Aprobado para caracterización ME',
                fecha: '2026-06-06',
            },
        }),
        remision({
            id: 'MED-011',
            caso: { id: 'SAL-011', nombre: 'Angélica María Soto', docNumero: '1505566778', nivelRiesgo: 'Crítico', municipio: 'Palmira' },
            remitidoPor: { nombre: 'Juan Pérez', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-04',
            estado: 'Por cotizar',
            medidasSolicitadas: ['Transporte de emergencia', 'Atención en salud'],
            agenteAsignado: agente('am-3'),
            motivoRemision: 'Víctima requiere traslado urgente a otro municipio y atención médica por lesiones recientes.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: {
                profesional: 'Camila Andrea Rojas Pinzón',
                concepto: 'El caso cumple con los criterios de riesgo de feminicidio y contacto reciente con el agresor.',
                resultado: 'Aprobado para caracterización ME',
                fecha: '2026-06-03',
            },
        }),
        remision({
            id: 'MED-012',
            caso: { id: 'SAL-012', nombre: 'Luz Dary Bermúdez', docNumero: '1506677889', nivelRiesgo: 'Alto', municipio: 'Buenaventura' },
            remitidoPor: { nombre: 'Sandra Mora', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-02',
            estado: 'Por aprobar cotizaciones',
            medidasSolicitadas: ['Alojamiento temporal', 'Kit de vestuario'],
            agenteAsignado: agente('am-1'),
            motivoRemision: 'Víctima desplazada con sus hijos menores de edad. Requiere alojamiento y vestuario de emergencia.',
            urgencia: 'Normal',
            revisionTecnica: {
                profesional: 'María Fernanda Oquendo Ibarra',
                concepto: 'El caso cumple con los criterios de riesgo de feminicidio y contacto reciente con el agresor.',
                resultado: 'Aprobado para caracterización ME',
                fecha: '2026-06-01',
            },
        }),
        remision({
            id: 'MED-013',
            caso: { id: 'SAL-013', nombre: 'Mariana Elisa Ríos Pardo', docNumero: '1507788990', nivelRiesgo: 'Crítico', municipio: 'Cali' },
            remitidoPor: { nombre: 'Jorge Ramírez', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-15',
            estado: 'Riesgo inminente: por escalar a Cruz Roja',
            medidasSolicitadas: ['Alojamiento temporal', 'Transporte de emergencia'],
            agenteAsignado: null,
            motivoRemision: 'Agresor ingresó a la vivienda de la víctima con arma de fuego. Supervisor ME decidió derivar a Supervisión del Convenio para entrega inmediata.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-014',
            caso: { id: 'SAL-014', nombre: 'Tatiana Rocío Nieto', docNumero: '1508899001', nivelRiesgo: 'Crítico', municipio: 'Buenaventura' },
            remitidoPor: { nombre: 'Laura Jiménez', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-13',
            estado: 'Riesgo inminente: por entregar medida',
            medidasSolicitadas: ['Alojamiento temporal'],
            agenteAsignado: null,
            motivoRemision: 'Supervisión del Convenio escaló el caso a Cruz Roja. Pendiente entrega de hotel sin esperar caracterización completa.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
        remision({
            id: 'MED-015',
            caso: { id: 'SAL-015', nombre: 'Carolina Andrea Vélez', docNumero: '1509900112', nivelRiesgo: 'Crítico', municipio: 'Palmira' },
            remitidoPor: { nombre: 'María González', equipo: 'Agente de Seguimiento' },
            fechaRemision: '2026-06-11',
            estado: 'Riesgo inminente: por confirmar seguridad',
            medidasSolicitadas: ['Alojamiento temporal', 'Kit de alimentación'],
            agenteAsignado: agente('am-2'),
            motivoRemision: 'Cruz Roja ya entregó alojamiento de emergencia. Falta confirmar que la persona está en lugar seguro para retomar el curso normal del caso.',
            urgencia: 'Riesgo inminente',
            revisionTecnica: null,
        }),
    ];

    function clonar(valor) {
        return JSON.parse(JSON.stringify(valor));
    }

    function hoyISO() {
        var d = new Date();
        var mes = String(d.getMonth() + 1).padStart(2, '0');
        var dia = String(d.getDate()).padStart(2, '0');
        return d.getFullYear() + '-' + mes + '-' + dia;
    }

    function contar(lista, pred) {
        var n = 0;
        for (var i = 0; i < lista.length; i++) {
            if (pred(lista[i])) n++;
        }
        return n;
    }

    vueApp.component('remisiones-medidas-component', {
        delimiters: ['${', '}'],
        template: window.__remisionesMedidasComponentTpl || '#tpl-remisiones-medidas-component',

        props: {
            vista: { type: String, default: 'supervisor' },
            agenteId: { type: String, default: '' },
            puedeReasignar: { type: Boolean, default: false },
        },

        emits: ['ver-caso', 'ver-detalle'],

        data: function() {
            return {
                agentes: AGENTES,
                profesionales: PROFESIONALES,
                resultadosRevision: RESULTADOS_REVISION,
                remisiones: clonar(SEED),
                searchText: '',
                filtroEstado: '',
                filtroAgente: '',
                filtroUrgencia: '',
                selectedIds: [],
                reasignarAbierto: false,
                reasignarAgenteId: '',
                revisionAbierta: false,
                revisionItem: null,
                revisionProfesional: '',
                revisionConcepto: '',
                revisionResultado: '',
            };
        },

        computed: {
            esSupervisor: function() { return this.vista === 'supervisor'; },
            esAgente: function() { return this.vista === 'agente'; },
            muestraFiltroEstado: function() {
                return this.vista === 'supervisor' || this.vista === 'agente';
            },
            muestraFiltroAgente: function() { return this.vista === 'supervisor'; },
            placeholderBusqueda: function() {
                return this.esAgente ? 'Nombre o número de identidad' : 'Nombre o doc. identidad';
            },
            labelLimpiar: function() {
                return this.esAgente ? 'Limpiar' : 'Limpiar filtros';
            },
            estadosFiltro: function() {
                return this.esAgente ? ESTADOS_AGENTE : ESTADOS;
            },
            claseGrilla: function() {
                return this.puedeReasignar ? 'rme-grid--check' : 'rme-grid--3';
            },
            alcance: function() {
                var self = this;
                if (this.vista === 'equipo-tecnico') {
                    return this.remisiones.filter(function(r) { return r.estado === 'Pendiente de Revisión'; });
                }
                if (this.vista === 'supervision-convenio') {
                    return this.remisiones.filter(function(r) { return ESTADOS_CONVENIO.indexOf(r.estado) !== -1; });
                }
                if (this.vista === 'cruz-roja') {
                    return this.remisiones.filter(function(r) { return ESTADOS_CRUZ_ROJA.indexOf(r.estado) !== -1; });
                }
                if (this.vista === 'agente') {
                    return this.remisiones.filter(function(r) {
                        return r.agenteAsignado && r.agenteAsignado.id === self.agenteId && ESTADOS_AGENTE.indexOf(r.estado) !== -1;
                    });
                }
                return this.remisiones.slice();
            },
            visibles: function() {
                var texto = this.searchText;
                var estado = this.filtroEstado;
                var agenteFiltro = this.filtroAgente;
                var urgencia = this.filtroUrgencia;
                var lista = this.alcance.filter(function(r) {
                    if (texto) {
                        var q = texto.toLowerCase();
                        var porNombre = r.caso.nombre.toLowerCase().indexOf(q) !== -1;
                        var porDoc = r.caso.docNumero.indexOf(texto) !== -1;
                        if (!porNombre && !porDoc) return false;
                    }
                    if (estado && r.estado !== estado) return false;
                    if (urgencia && r.urgencia !== urgencia) return false;
                    if (agenteFiltro) {
                        if (agenteFiltro === 'sin-asignar') {
                            if (r.agenteAsignado) return false;
                        } else if (!r.agenteAsignado || r.agenteAsignado.id !== agenteFiltro) {
                            return false;
                        }
                    }
                    return true;
                });
                lista.sort(function(a, b) {
                    var ua = a.urgencia === 'Riesgo inminente' ? 0 : 1;
                    var ub = b.urgencia === 'Riesgo inminente' ? 0 : 1;
                    return ua - ub;
                });
                return lista;
            },
            todosSeleccionados: function() {
                var self = this;
                return this.visibles.length > 0 && this.visibles.every(function(r) {
                    return self.selectedIds.indexOf(r.id) !== -1;
                });
            },
            cards: function() {
                var lista = this.alcance;
                if (this.esSupervisor) {
                    return [
                        { key: 'total', label: 'Total', tono: 'total', valor: lista.length },
                        { key: 'riesgo', label: '🔥 Riesgo inminente', tono: 'riesgo', valor: contar(lista, function(r) { return ESTADOS_RIESGO_INMINENTE.indexOf(r.estado) !== -1; }) },
                        { key: 'revision', label: 'Pendiente de Revisión', tono: 'revision', valor: contar(lista, function(r) { return r.estado === 'Pendiente de Revisión'; }) },
                        { key: 'porAsignar', label: 'Por asignar', tono: 'asignar', valor: contar(lista, function(r) { return r.estado === 'Por asignar'; }) },
                        { key: 'asignadas', label: 'Asignadas', tono: 'asignadas', valor: contar(lista, function(r) { return r.estado === 'Asignada'; }) },
                        { key: 'gestion', label: 'En gestión documental', tono: 'gestion', valor: contar(lista, function(r) { return r.estado === 'En gestión documental'; }) },
                        { key: 'caracterizado', label: 'Caracterizado', tono: 'caracterizado', valor: contar(lista, function(r) { return r.estado === 'Caracterizado'; }) },
                        { key: 'aprobadas', label: 'Aprobadas', tono: 'aprobadas', valor: contar(lista, function(r) { return r.estado === 'Aprobadas'; }) },
                    ];
                }
                if (this.esAgente) {
                    return [
                        { key: 'total', label: 'Mis remisiones', tono: 'total', valor: lista.length },
                        { key: 'riesgo', label: '🔥 Riesgo inminente', tono: 'riesgo', valor: contar(lista, function(r) { return r.estado === 'Riesgo inminente: por confirmar seguridad'; }) },
                        { key: 'asignadas', label: 'Asignadas', tono: 'asignadas', valor: contar(lista, function(r) { return r.estado === 'Asignada'; }) },
                        { key: 'gestion', label: 'En gestión documental', tono: 'gestion', valor: contar(lista, function(r) { return r.estado === 'En gestión documental'; }) },
                        { key: 'aprobadas', label: 'Aprobadas', tono: 'aprobadas', valor: contar(lista, function(r) { return r.estado === 'Aprobadas'; }) },
                        { key: 'seguimiento', label: 'En Seguimiento', tono: 'seguimiento', valor: contar(lista, function(r) { return r.estado === 'En Seguimiento'; }) },
                    ];
                }
                return [];
            },
            banner: function() {
                if (this.vista === 'cruz-roja') {
                    return {
                        tono: 'cruz-roja',
                        icono: '🤝',
                        texto: 'Mostrando casos "Por cotizar", "Riesgo inminente: por entregar medida" (prioritario) y "Aprobadas" pendientes de otorgamiento.',
                    };
                }
                if (this.vista === 'equipo-tecnico') {
                    return {
                        tono: 'equipo',
                        icono: '⚖️',
                        texto: 'Bandeja "Pendiente de Revisión": casos remitidos a Medidas de Emergencia a la espera del concepto cualificado del Equipo Técnico.',
                    };
                }
                if (this.vista === 'supervision-convenio') {
                    return {
                        tono: 'convenio',
                        icono: '📂',
                        texto: 'Mostrando casos "Por revisar trazabilidad", "Por aprobar cotizaciones" y "Riesgo inminente: por escalar a Cruz Roja" (prioritario) a cargo de Supervisión del Convenio.',
                    };
                }
                return null;
            },
            emptyMessage: function() {
                if (this.esAgente && this.alcance.length === 0) {
                    return 'No tienes remisiones asignadas.';
                }
                return 'No hay remisiones que coincidan con los filtros.';
            },
            revisionLista: function() {
                return !!(this.revisionProfesional && this.revisionConcepto.trim() && this.revisionResultado);
            },
        },

        methods: {
            inicial: function(nombre) {
                return nombre ? nombre.charAt(0) : '';
            },
            claseRiesgo: function(nivel) {
                if (nivel === 'Crítico') return 'rme-risk--critico';
                if (nivel === 'Alto') return 'rme-risk--alto';
                return 'rme-risk--medio';
            },
            claseUrgencia: function(urgencia) {
                return urgencia === 'Riesgo inminente' ? 'rme-urgencia--inminente' : 'rme-urgencia--normal';
            },
            claseEstado: function(estado) {
                return 'rme-estado--' + (ESTADO_MOD[estado] || 'default');
            },
            isSelected: function(id) {
                return this.selectedIds.indexOf(id) !== -1;
            },
            toggleOne: function(id) {
                var idx = this.selectedIds.indexOf(id);
                if (idx === -1) this.selectedIds.push(id);
                else this.selectedIds.splice(idx, 1);
            },
            toggleAll: function() {
                if (this.todosSeleccionados) {
                    this.selectedIds = [];
                    return;
                }
                this.selectedIds = this.visibles.map(function(r) { return r.id; });
            },
            limpiarFiltros: function() {
                this.searchText = '';
                this.filtroEstado = '';
                this.filtroAgente = '';
                this.filtroUrgencia = '';
            },
            accionDe: function(item) {
                var estado = item.estado;
                var vista = this.vista;
                if (vista === 'supervisor' && estado === 'Riesgo inminente: por decidir vía') {
                    return { tipo: 'detalle', label: '🔥 Decidir vía de atención', clase: 'rme-action--decidir rme-action--pulse' };
                }
                if (vista === 'supervision-convenio' && estado === 'Riesgo inminente: por escalar a Cruz Roja') {
                    return { tipo: 'detalle', label: '🔥 Escalar a Cruz Roja', clase: 'rme-action--escalar' };
                }
                if (vista === 'cruz-roja' && estado === 'Riesgo inminente: por entregar medida') {
                    return { tipo: 'detalle', label: '🔥 Entregar medida inmediata', clase: 'rme-action--entregar' };
                }
                if ((vista === 'supervisor' || vista === 'equipo-tecnico') && estado === 'Pendiente de Revisión') {
                    return { tipo: 'revisar', label: '⚖️ Revisar (Equipo Técnico)', clase: 'rme-action--revisar' };
                }
                if (this.puedeReasignar && estado === 'Por asignar') {
                    return { tipo: 'asignar', label: '👤 Asignar agente', clase: 'rme-action--asignar' };
                }
                if (vista === 'supervisor' && estado === 'Caracterizado') {
                    return { tipo: 'detalle', label: '📄 Revisar caracterización', clase: 'rme-action--caracterizacion' };
                }
                if (vista === 'supervision-convenio' && estado === 'Por revisar trazabilidad') {
                    return { tipo: 'detalle', label: '📂 Revisar trazabilidad', clase: 'rme-action--trazabilidad' };
                }
                if (vista === 'supervision-convenio' && estado === 'Por aprobar cotizaciones') {
                    return { tipo: 'detalle', label: '✅ Aprobar cotizaciones', clase: 'rme-action--cotizaciones' };
                }
                if (vista === 'cruz-roja' && estado === 'Por cotizar') {
                    return { tipo: 'detalle', label: '💰 Cargar cotizaciones', clase: 'rme-action--cotizar' };
                }
                if (vista === 'agente' && estado === 'Riesgo inminente: por confirmar seguridad') {
                    return { tipo: 'detalle', label: '🔥 Confirmar seguridad', clase: 'rme-action--confirmar' };
                }
                return null;
            },
            onAccion: function(item) {
                var accion = this.accionDe(item);
                if (!accion) return;
                if (accion.tipo === 'revisar') this.abrirRevision(item);
                else if (accion.tipo === 'asignar') this.abrirAsignarUno(item.id);
                else this.onVerDetalle(item);
            },
            onVerCaso: function(item) {
                this.$emit('ver-caso', { caseICode: item.caso.id, remisionId: item.id });
            },
            onVerDetalle: function(item) {
                this.$emit('ver-detalle', { remisionId: item.id, caseICode: item.caso.id });
            },
            abrirReasignarMasivo: function() {
                this.reasignarAgenteId = '';
                this.reasignarAbierto = true;
            },
            abrirAsignarUno: function(id) {
                this.selectedIds = [id];
                this.reasignarAgenteId = '';
                this.reasignarAbierto = true;
            },
            cerrarReasignar: function() {
                this.reasignarAbierto = false;
                this.reasignarAgenteId = '';
            },
            confirmarReasignar: function() {
                var destino = null;
                for (var i = 0; i < AGENTES.length; i++) {
                    if (AGENTES[i].id === this.reasignarAgenteId) destino = AGENTES[i];
                }
                if (!destino) return;
                var ids = this.selectedIds.slice();
                this.remisiones.forEach(function(r) {
                    if (ids.indexOf(r.id) === -1) return;
                    r.agenteAsignado = { id: destino.id, nombre: destino.nombre };
                    r.estado = 'Asignada';
                });
                this.selectedIds = [];
                this.cerrarReasignar();
            },
            abrirRevision: function(item) {
                this.revisionItem = item;
                this.revisionProfesional = '';
                this.revisionConcepto = '';
                this.revisionResultado = '';
                this.revisionAbierta = true;
            },
            cerrarRevision: function() {
                this.revisionAbierta = false;
                this.revisionItem = null;
                this.revisionProfesional = '';
                this.revisionConcepto = '';
                this.revisionResultado = '';
            },
            confirmarRevision: function() {
                if (!this.revisionLista || !this.revisionItem) return;
                var id = this.revisionItem.id;
                var profesional = this.revisionProfesional;
                var concepto = this.revisionConcepto.trim();
                var resultado = this.revisionResultado;
                var estado = RESULTADO_A_ESTADO[resultado];
                this.remisiones.forEach(function(r) {
                    if (r.id !== id) return;
                    r.estado = estado;
                    r.revisionTecnica = {
                        profesional: profesional,
                        concepto: concepto,
                        resultado: resultado,
                        fecha: hoyISO(),
                    };
                });
                this.cerrarRevision();
            },
        },
    });
})();
