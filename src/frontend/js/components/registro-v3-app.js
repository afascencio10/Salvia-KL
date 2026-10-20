// Registro V3 — app principal. Sin dinamic-form (decisión A1): Habeas Data
// -> Cabezote (búsqueda por documento y/o teléfono, D7) -> Formulario (7
// módulos CI/DP/CA/HV/AG/TZ/GS) -> Cierre. Usa registro-v3-campos.js como
// catálogo de preguntas y GET /victim-cases-v3/enums como catálogo de
// opciones (icode real de victim_case_form2_enums).
var RegistroV3App = Vue.createApp({
  delimiters: ['${', '}'],
  data: function () {
    return {
      etapa: 'habeas', // habeas | formularioCerrado | tipoContacto | cabezote | formulario | cierre
      habeasAceptado: false,
      tipoContactoNoDisponible: false,

      buscar: { tipoDocumento: '', documento: '', telefono: '' },
      resultadosBusqueda: [],
      buscando: false,
      errorBusqueda: '',

      icode: null,
      resumenCaso: null,
      answers: {},
      entidades: [],
      barreras: [],
      enumCatalog: {},
      triadaCompleta: false,

      moduloActivo: 'CI',
      modulos: [
        { id: 'CI', label: 'Contacto e identificación', campos: CAMPOS_CI_V3 },
        { id: 'DP', label: 'Datos personales', campos: CAMPOS_DP_V3 },
        { id: 'CA', label: 'Contacto de apoyo', campos: CAMPOS_CA_V3 },
        { id: 'HV', label: 'Hechos de violencia', campos: CAMPOS_HV_V3 },
        { id: 'AG', label: 'Características del agresor', campos: CAMPOS_AG_V3 },
        { id: 'TZ', label: 'Tamizaje de riesgo', campos: [] },
        { id: 'GS', label: 'Gestión Salvia', campos: CAMPOS_GS_A_V3 },
      ],
      guardando: false,
      mensajeGuardado: '',
      mensajeModuloBloqueado: '',
      mensajeCierre: '',
      avisoDuplicidad: '',
      errorFinalizar: '',

      // GS § C -- buscador real Departamento/Ciudad/Municipio + Sector ->
      // entity_branch, igual que case-entities.js / Seguimiento (ver
      // registro-v3-pendientes.md § B).
      geoDepartamentos: [],
      entidadCiudades: [], entidadMunicipios: [], entidadesDisponibles: [],
      entidadBuscando: false, obligacionesEntidad: [], obligacionesCargando: false,
      errorEntidad: '',
      // entidadGeoPorBranch cachea depto/ciudad/municipio/sector usados al
      // buscar cada entidad agregada (por entityBranchId) -- GS § D (Barreras)
      // los reusa para no repetir la búsqueda: la barrera es un sub-repeater
      // de la entidad, no un buscador independiente (instrucción 2026-10-15).
      entidadGeoPorBranch: {},
      // GS § D -- barreraParaEntidad es el entityBranchId de la entidad cuyo
      // formulario de "+ Agregar barrera" está abierto (null = ninguno).
      barreraParaEntidad: null,
      competenciasBarrera: [], competenciasCargando: false,
      errorBarrera: '',

      // Catálogos globales (registro-v3-campos.js) que el TEMPLATE referencia
      // directo (no a través de un método) -- BUG encontrado en QA
      // (2026-10-15): Vue, al compilar el template en el navegador, reescribe
      // cualquier identificador que no reconoce como dato/método/prop del
      // componente a _ctx.<nombre> en vez de dejarlo caer al scope global
      // (a diferencia de lo esperado de un `with` normal) -- una referencia
      // directa a una var global en el HTML (ej. "${TAMIZAJE_VPP_V3.label}")
      // revienta el render entero con "Cannot read properties of undefined".
      // Exponerlas aquí como datos del componente resuelve _ctx.X
      // correctamente. Cualquier catálogo nuevo que se use DIRECTO en el
      // HTML (no envuelto en un método) debe agregarse aquí también.
      TAMIZAJE_GENERALES_V3: TAMIZAJE_GENERALES_V3,
      TAMIZAJE_VPP_V3: TAMIZAJE_VPP_V3,
      SECTORES_ENTIDAD_BARRERA_V3: SECTORES_ENTIDAD_BARRERA_V3,
      CANAL_ACTIVACION_V3: CANAL_ACTIVACION_V3,

      // Tamizaje
      tzRevelarResto: false,
      tzResultado: null,
      tzGuardando: false,

      // GS § B remisión
      gsRequiereRemision: false,
      gsJustificacionRemision: '',

      // GS § C entidades
      nuevaEntidad: this.entidadVacia(),
      // GS § D barreras
      nuevaBarrera: this.barreraVacia(),
    };
  },

  computed: {
    moduloActualObj: function () {
      var self = this;
      return this.modulos.filter(function (m) { return m.id === self.moduloActivo; })[0];
    },
    equiposRemision: function () { return EQUIPOS_REMISION_V3; },
    equiposSeleccionados: function () {
      return (this.answers['GS_6'] || '').split(',').filter(function (v) { return v; });
    },
  },

  created: function () {
    this.camposPorId = {};
    var self = this;
    this.modulos.forEach(function (m) {
      m.campos.forEach(function (c) { self.camposPorId[c.id] = c; });
    });
    CAMPOS_GS_B_V3.concat(CAMPOS_GS_CRITERIOS_V3).forEach(function (c) { self.camposPorId[c.id] = c; });
  },

  mounted: function () {
    this.cargarEnumCatalog();
    this.cargarGeoDepartamentos();
    // D11 (flow-E14): respaldo para cuando el agente cierra la pestaña sin
    // presionar "Finalizar llamada". navigator.sendBeacon (no fetch) porque
    // es la API pensada para esto -- garantiza el intento de envío aunque
    // la página se descargue antes de que una promesa normal resuelva (el
    // propio GAP de flow-E14 señalaba que un fetch común es poco confiable
    // acá).
    var self = this;
    window.addEventListener('beforeunload', function () {
      if (self.etapa !== 'formulario') return;
      if (['CI', 'DP', 'CA', 'HV', 'AG', 'GS'].indexOf(self.moduloActivo) === -1) return;
      var payload = {};
      Object.keys(self.answers).forEach(function (k) {
        if (k.indexOf(self.moduloActivo + '_') === 0) payload[k] = self.answers[k];
      });
      var blob = new Blob([JSON.stringify({ answers: payload })], { type: 'application/json' });
      navigator.sendBeacon('/api/v1/victim-cases-v3/' + self.icode + '/seccion/' + self.moduloActivo, blob);
    });
  },

  methods: {
    // depto/ciudad/municipio/sector son estado LOCAL del buscador -- nunca se
    // mandan al backend, solo resuelven EntityBranchID (ver
    // buscarEntidadesEntidad). El resto sí se manda, en PascalCase idéntico
    // al campo Go correspondiente (repository.VictimCaseEntityInput).
    entidadVacia: function () {
      return {
        Kind: 'acudida', depto: '', ciudad: '', municipio: '', sector: '',
        EntityBranchID: '', FileNumber: '', ActivationChannel: '',
        CompletedActions: [], PendingActions: [], OtherInfo: '',
      };
    },
    // Acá Department/City/Town SÍ se mandan -- son la ubicación donde se
    // identificó la barrera (BarrierV2.DepartmentID/CityID/TownID), no solo
    // un filtro. SectorCode también se manda tal cual (BarrierV2.Sector).
    // Department/City/Town/SectorCode/EntityBranchID se heredan de la
    // entidad padre (abrirBarreraParaEntidad) -- ya no los edita el agente
    // acá directamente.
    barreraVacia: function () {
      return {
        Department: '', City: '', Town: '', SectorCode: '', EntityBranchID: '',
        Specific: [], OtherDetail: '',
        Institutional: [], Economic: [], Territorial: [], Differential: [],
        Date: '', Official: '', Description: '', Management: [], Competences: [],
      };
    },

    // ── Habeas data / Cabezote ──────────────────────────────────────────
    // D1/D2 (registro-v3-index.md § Inventario de eventos, "Cuando responde
    // la autorización..."/"Cuando selecciona el tipo de contacto"): antes
    // "Acepto y continuar" saltaba directo a Cabezote, sin la rama "No" ni
    // la etapa de Tipo de contacto.
    aceptarHabeasData: function () {
      this.habeasAceptado = true;
      this.etapa = 'tipoContacto';
    },
    rechazarHabeasData: function () {
      this.etapa = 'formularioCerrado';
    },
    seleccionarTipoContacto: function (tipo) {
      if (tipo === 'victima_caso_nuevo') {
        this.etapa = 'cabezote';
        return;
      }
      // C9: Víctima-Seguimiento, Tercero-Reporte, Familiares (feminicidio) y
      // Hombres abren formularios/pantallas distintas a esta. Esta pantalla
      // no conoce sus URLs reales (no están resueltas en ninguna fuente de
      // planeación -- G-V3-11 para Familiares/Hombres) -- se avisa en vez
      // de inventar un enlace.
      this.tipoContactoNoDisponible = true;
    },

    ejecutarBusqueda: function () {
      if (!this.buscar.documento && !this.buscar.telefono) {
        this.errorBusqueda = 'Digite al menos documento o teléfono.';
        return;
      }
      this.errorBusqueda = '';
      this.buscando = true;
      var self = this;
      var qs = new URLSearchParams();
      if (this.buscar.tipoDocumento) qs.set('tipo_documento', this.buscar.tipoDocumento);
      if (this.buscar.documento) qs.set('documento', this.buscar.documento);
      if (this.buscar.telefono) qs.set('telefono', this.buscar.telefono);
      fetch('/api/v1/victim-cases-v3/buscar?' + qs.toString())
        .then(function (r) { return r.json(); })
        .then(function (data) { self.resultadosBusqueda = data.registros || []; })
        .catch(function () { self.errorBusqueda = 'Error al buscar.'; })
        .finally(function () { self.buscando = false; });
    },

    crearNuevoRegistro: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/draft', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tipo_documento: this.buscar.tipoDocumento,
          documento: this.buscar.documento,
          telefono: this.buscar.telefono,
        }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          if (data.icode) { self.abrirCaso(data.icode); }
        });
    },

    abrirCaso: function (icode) {
      var self = this;
      fetch('/api/v1/victim-cases-v3/' + icode)
        .then(function (r) { return r.json(); })
        .then(function (data) {
          self.icode = icode;
          self.resumenCaso = data.summary || null;
          // Solo se repueblan las keys en answeredKeys -- el resto de
          // 'answers' son valores dummy (D8) que nunca respondió el
          // agente y no deben mostrarse como si lo hubiera hecho (ver
          // hallazgo H3, registro-v3-qa-plan.md).
          var todas = data.answers || {};
          var respondidas = data.answeredKeys || [];
          var filtradas = {};
          respondidas.forEach(function (k) {
            if (todas[k] !== undefined) filtradas[k] = todas[k];
          });
          self.answers = filtradas;
          self.entidades = data.entities || [];
          self.barreras = data.barriers || [];
          // BUG encontrado en QA (2026-10-15): tzResultado solo se llenaba
          // con la respuesta directa de guardarTamizaje() -- en memoria,
          // nunca releído del backend. Cualquier reapertura del Caso (o
          // Finalizar en otra sesión distinta a la que guardó el Tamizaje)
          // mostraba "Nivel de riesgo: —" en el Cierre aunque el resultado
          // real siguiera guardado en BD (GetCase ya lo manda ahora).
          self.tzResultado = data.tzResultado || null;
          // A2: antes solo se marcaba tras un guardado en vivo -- un Caso ya
          // completo se reabría con el sidebar bloqueado otra vez. GetCase
          // ahora manda este dato calculado con el mismo criterio que
          // SaveSection (ver victim_case_v3_service.go, isTriadaCompleteV3).
          self.triadaCompleta = !!data.triadaCompleta;
          self.etapa = 'formulario';
          self.moduloActivo = 'CI';
        });
    },

    cargarEnumCatalog: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/enums')
        .then(function (r) { return r.json(); })
        .then(function (data) { self.enumCatalog = data || {}; });
    },

    // ── Catálogo de opciones / etiquetas ─────────────────────────────────
    campoPorId: function (id) { return this.camposPorId[id]; },

    opcionesPara: function (campo) {
      if (!campo) return [];
      if (campo.opciones) {
        return campo.opciones.map(function (o) { return { code: o, icode: o, label: o }; });
      }
      if (campo.subtipoDe) {
        var tiposSeleccionados = (this.answers[campo.subtipoDe] || '').split(',').filter(function (v) { return v; });
        var tipoOpts = this.opcionesPara(this.camposPorId[campo.subtipoDe]);
        var categorias = [];
        var self = this;
        tiposSeleccionados.forEach(function (icode) {
          var opt = tipoOpts.filter(function (o) { return o.icode === icode; })[0];
          var suf = opt && TIPO_A_SUFIJO_SUBTIPO_V3[opt.code];
          if (suf) categorias.push('victim_case_form2_subtype_violence_experienced_' + suf);
        });
        var out = [];
        categorias.forEach(function (cat) { out = out.concat(self.enumCatalog[cat] || []); });
        return out;
      }
      if (campo.enumCategories) {
        var merged = [];
        var seen = {};
        var self2 = this;
        campo.enumCategories.forEach(function (cat) {
          (self2.enumCatalog[cat] || []).forEach(function (o) {
            if (!seen[o.icode]) { seen[o.icode] = true; merged.push(o); }
          });
        });
        return merged;
      }
      if (campo.enumCategory) {
        return this.enumCatalog[campo.enumCategory] || [];
      }
      return [];
    },

    etiqueta: function (campo, opt) {
      if (!campo || !opt) return '';
      var overrides = campo.enumCategory && ENUM_LABELS_OVERRIDE_V3[campo.enumCategory];
      if (overrides && overrides[opt.code]) return overrides[opt.code];
      if (campo.subtipoDe || campo.enumCategories) {
        // busca override en cualquiera de las categorías candidatas
        var cats = campo.enumCategories || Object.keys(ENUM_LABELS_OVERRIDE_V3);
        for (var i = 0; i < cats.length; i++) {
          var ov = ENUM_LABELS_OVERRIDE_V3[cats[i]];
          if (ov && ov[opt.code]) return ov[opt.code];
        }
      }
      return opt.label;
    },

    etiquetaSeleccionada: function (fieldId) {
      var campo = this.camposPorId[fieldId];
      var icode = this.answers[fieldId];
      if (!campo || !icode) return '';
      var opts = this.opcionesPara(campo);
      for (var i = 0; i < opts.length; i++) {
        if (opts[i].icode === icode) return this.etiqueta(campo, opts[i]);
      }
      return '';
    },

    multipleIncluyeTexto: function (fieldId, texto) {
      var campo = this.camposPorId[fieldId];
      var raw = this.answers[fieldId] || '';
      if (!raw || !campo) return false;
      var icodes = raw.split(',');
      var opts = this.opcionesPara(campo);
      texto = texto.toLowerCase();
      for (var i = 0; i < icodes.length; i++) {
        for (var j = 0; j < opts.length; j++) {
          if (opts[j].icode === icodes[i] && this.etiqueta(campo, opts[j]).toLowerCase().indexOf(texto) >= 0) return true;
        }
      }
      return false;
    },

    // ── Valores / visibilidad ────────────────────────────────────────────
    getVal: function (id) { return this.answers[id] || ''; },
    // D6 (flow-E11, auto-sugerencias): si el tipo de documento (CI_5) es
    // inequívocamente colombiano o extranjero, predilligencia DP_5
    // (Nacionalidad) -- nunca pisa una respuesta que el agente ya haya
    // dado a mano. cd/oe/si quedan ambiguos a propósito (Carné Diplomático,
    // Otro, Persona sin identificación no determinan nacionalidad por sí
    // solos).
    sugerirNacionalidadPorDocumento: function () {
      if (this.answers['DP_5']) return;
      var codigosColombianos = ['cc', 'ti', 'rc', 'nc'];
      var codigosExtranjeros = ['ce', 'ps', 'pp', 'pf', 'pt', 'dn', 'vs', 'vr', 'sc', 'mv', 'no', 'pn'];
      var code = this.codigoSeleccionado('CI_5');
      var opts = this.opcionesPara(this.camposPorId['DP_5']);
      var destino = null;
      if (codigosColombianos.indexOf(code) >= 0) destino = 'co';
      else if (codigosExtranjeros.indexOf(code) >= 0) destino = 'ex';
      if (!destino) return;
      var opt = opts.filter(function (o) { return o.code === destino; })[0];
      if (opt) this.answers['DP_5'] = opt.icode;
    },
    setVal: function (id, v) {
      this.answers[id] = v;
      if (id === 'CI_5') this.sugerirNacionalidadPorDocumento();
      if (id === 'CI_5' || id === 'CI_6') this.verificarDuplicidadDocumento();
    },
    // D7 (flow-E11, "verificación de duplicidad si el campo es el número de
    // documento"): reusa el mismo /buscar que usa el Cabezote (D7 de
    // arquitectura: documento, nunca nombre) -- no hace falta un endpoint
    // nuevo. Debounced 500ms para no buscar en cada tecla.
    verificarDuplicidadDocumento: function () {
      var self = this;
      this.avisoDuplicidad = '';
      clearTimeout(this._duplicidadTimer);
      var code = this.codigoSeleccionado('CI_5');
      var numero = (this.answers['CI_6'] || '').trim();
      if (!code || !numero) return;
      this._duplicidadTimer = setTimeout(function () {
        fetch('/api/v1/victim-cases-v3/buscar?tipo_documento=' + encodeURIComponent(code) + '&documento=' + encodeURIComponent(numero))
          .then(function (r) { return r.json(); })
          .then(function (data) {
            var otros = (data.registros || []).filter(function (r) { return r.icode !== self.icode; });
            if (otros.length) {
              self.avisoDuplicidad = 'Ya existe' + (otros.length > 1 ? 'n ' + otros.length + ' Registros' : ' un Registro') + ' con este documento: ' + otros.map(function (r) { return r.names + ' ' + r.lastNames; }).join(', ') + '.';
            }
          })
          .catch(function () {});
      }, 500);
    },
    // D6: tratamiento según CI_1 (pronombre) -- "determina el tratamiento
    // mostrado en cabezote (Sra./Sr./Sre.)". Solo de visualización, no es
    // una respuesta que se guarde en otro campo.
    tratamiento: function () {
      var code = this.codigoSeleccionado('CI_1');
      return { la: 'Sra.', el: 'Sr.', le: 'Sre.' }[code] || '';
    },
    getBool: function (id) { return this.answers[id] === 'true'; },
    setBool: function (id, v) { this.answers[id] = v ? 'true' : 'false'; },

    estaMarcado: function (id, icode) {
      return (this.answers[id] || '').split(',').indexOf(icode) >= 0;
    },
    toggleMarcado: function (id, icode, max) {
      var actuales = (this.answers[id] || '').split(',').filter(function (v) { return v; });
      var idx = actuales.indexOf(icode);
      if (idx >= 0) {
        actuales.splice(idx, 1);
      } else {
        if (max && actuales.length >= max) return;
        actuales.push(icode);
      }
      this.answers[id] = actuales.join(',');
    },

    // Resuelve el `code` (no el icode) de la respuesta actual de un campo
    // dropdown de un solo valor -- D12 lo necesita para "Q17 o Q18 = Persona
    // en ASP" sin repetir la lógica de condicionalEnumCode por cada campo.
    codigoSeleccionado: function (fieldId) {
      var campo = this.camposPorId[fieldId];
      var icode = this.answers[fieldId];
      if (!campo || !icode) return '';
      var opt = this.opcionesPara(campo).filter(function (o) { return o.icode === icode; })[0];
      return opt ? opt.code : '';
    },
    pseudoCondicion: function (nombre) {
      switch (nombre) {
        case '__gestacionAplica': return true;
        case '__haySubtipo': return !!(this.answers['HV_9'] && this.answers['HV_9'].length);
        case '__ambitoLaboral': return this.multipleIncluyeTexto('HV_11', 'laboral');
        case '__agresorConocido':
          var et = this.etiquetaSeleccionada('AG_2').toLowerCase();
          return !!et && et.indexOf('desconocid') === -1;
        // D12: DP_20/21/22 (modalidad/fecha/motivo ASP) -- secciones.md
        // "condicional (Q17 o Q18 = Persona en ASP)". DP_17 occupation
        // code 'as', DP_18 income_generation_method code 'pr'.
        case '__personaEnASP':
          return this.codigoSeleccionado('DP_17') === 'as' || this.codigoSeleccionado('DP_18') === 'pr';
        // D13: CI_6 (número de documento) no aplica si CI_5="Persona sin
        // identificación" (code 'si').
        case '__tieneNumeroDocumento':
          return this.codigoSeleccionado('CI_5') !== 'si';
        default: return true;
      }
    },
    // D13: máscara de validación de CI_6 según el tipo de documento (CI_5)
    // -- secciones.md trae una tabla completa que antes no se aplicaba.
    ci6Mascara: function () {
      return MASCARA_DOC_NUMERO_V3[this.codigoSeleccionado('CI_5')];
    },
    ci6Invalido: function () {
      var m = this.ci6Mascara();
      var v = (this.answers['CI_6'] || '').trim();
      return !!m && !!v && !m.regex.test(v);
    },

    mostrarCampo: function (campo) {
      if (campo.condicional) {
        if (campo.condicional.campo.indexOf('__') === 0) {
          return this.pseudoCondicion(campo.condicional.campo);
        }
        return this.answers[campo.condicional.campo] === String(campo.condicional.valor);
      }
      if (campo.condicionalEnumCode) {
        var c = campo.condicionalEnumCode;
        var opts = this.opcionesPara(this.camposPorId[c.campo]);
        var sel = this.answers[c.campo];
        var opt = opts.filter(function (o) { return o.icode === sel; })[0];
        return !!opt && opt.code === c.code;
      }
      // D14: igual que condicionalEnumCode pero contra una lista de codes
      // válidos en vez de uno solo (ej. DP_19 -- "Q18=contrato formal O independiente").
      if (campo.condicionalEnumCodeIn) {
        var ci = campo.condicionalEnumCodeIn;
        var optsIn = this.opcionesPara(this.camposPorId[ci.campo]);
        var selIn = this.answers[ci.campo];
        var optIn = optsIn.filter(function (o) { return o.icode === selIn; })[0];
        return !!optIn && ci.codes.indexOf(optIn.code) >= 0;
      }
      if (campo.condicionalEnumLabel) {
        var et = this.etiquetaSeleccionada(campo.condicionalEnumLabel.campo);
        return et.indexOf(campo.condicionalEnumLabel.contains) >= 0;
      }
      if (campo.subtipoDe) {
        return !!(this.answers[campo.subtipoDe] && this.answers[campo.subtipoDe].length);
      }
      return true;
    },

    camposVisibles: function (moduloId) {
      var self = this;
      var mod = this.modulos.filter(function (m) { return m.id === moduloId; })[0];
      if (!mod) return [];
      return mod.campos.filter(function (c) { return self.mostrarCampo(c); });
    },

    // A3 (flow-E12): "Validar formulario: Todas las preguntas [1er contacto]
    // visibles del módulo: required". Antes no existía ningún chequeo --
    // se podía guardar un módulo vacío. Solo mira las VISIBLES (camposVisibles
    // ya resuelve los condicionales) con req==='C1'; "false" cuenta como
    // respondido (es una respuesta real de un boolean), solo '' y undefined
    // cuentan como sin responder.
    obligatoriasFaltantesModulo: function (moduloId) {
      var self = this;
      return this.camposVisibles(moduloId).filter(function (c) {
        return c.req === 'C1' && !self.answers[c.id];
      });
    },
    // A4 (obs. #13 cliente): "solo letras" en nombres -- rechaza una entrada
    // compuesta ÚNICAMENTE por números. Antes `soloLetras: true` estaba
    // declarado en el catálogo (CI_2/CI_3/CA_1/AG_7) pero nada lo leía.
    esSoloNumeros: function (valor) {
      return /^[0-9]+$/.test((valor || '').trim());
    },
    campoSoloLetrasInvalido: function (campo) {
      return !!campo.soloLetras && this.esSoloNumeros(this.answers[campo.id]);
    },
    camposSoloLetrasInvalidosModulo: function (moduloId) {
      var self = this;
      return this.camposVisibles(moduloId).filter(function (c) { return self.campoSoloLetrasInvalido(c); });
    },
    // ── Guardar sección ──────────────────────────────────────────────────
    // enviarSeccion -- el POST real, sin ninguna validación. Lo usan tanto
    // guardarModulo (que valida antes de llamarlo) como guardarYSalir (D11,
    // flow-E14: "sin exigir que las obligatorias [1er contacto] estén
    // completas -- guardado parcial de cierre de llamada").
    enviarSeccion: function (moduloId) {
      var self = this;
      var payload = {};
      Object.keys(this.answers).forEach(function (k) {
        if (k.indexOf(moduloId + '_') === 0) payload[k] = self.answers[k];
      });
      return fetch('/api/v1/victim-cases-v3/' + this.icode + '/seccion/' + moduloId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ answers: payload }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          self.triadaCompleta = !!data.triadaCompleta;
          return data;
        });
    },
    guardarModulo: function (moduloId) {
      var self = this;
      var faltantes = this.obligatoriasFaltantesModulo(moduloId);
      if (faltantes.length) {
        this.mensajeGuardado = 'Faltan obligatorias de 1er contacto: ' + faltantes.map(function (c) { return c.label; }).join('; ');
        return;
      }
      var invalidos = this.camposSoloLetrasInvalidosModulo(moduloId);
      if (invalidos.length) {
        this.mensajeGuardado = 'Estos campos no pueden ser solo números: ' + invalidos.map(function (c) { return c.label; }).join('; ');
        return;
      }
      if (moduloId === 'CI' && this.ci6Invalido()) {
        this.mensajeGuardado = 'Número de documento inválido para el tipo seleccionado: ' + this.ci6Mascara().hint;
        return;
      }
      this.guardando = true;
      this.mensajeGuardado = '';
      this.enviarSeccion(moduloId)
        .then(function () { self.mensajeGuardado = 'Sección guardada.'; })
        .catch(function () { self.mensajeGuardado = 'Error al guardar la sección.'; })
        .finally(function () { self.guardando = false; });
    },
    // D11 (flow-E14, botón "Finalizar llamada (guardar y salir)"): guarda el
    // módulo actual SIN exigir las obligatorias -- es un cierre de llamada
    // a medias, no un guardado normal -- y vuelve al Cabezote. Antes no
    // existía ningún botón ni respaldo para este caso; el único guardado
    // disponible era el de guardarModulo, que SÍ exige obligatorias y
    // bloqueaba salir de una llamada incompleta.
    guardarYSalir: function () {
      var self = this;
      if (['CI', 'DP', 'CA', 'HV', 'AG', 'GS'].indexOf(this.moduloActivo) === -1) {
        this.volverACabezote();
        return;
      }
      this.guardando = true;
      this.enviarSeccion(this.moduloActivo)
        .catch(function () { /* best-effort -- igual se sale de la llamada aunque falle el guardado */ })
        .finally(function () {
          self.guardando = false;
          self.volverACabezote();
        });
    },
    volverACabezote: function () {
      this.etapa = 'cabezote';
      this.icode = null;
      this.answers = {};
      this.entidades = []; this.barreras = [];
    },

    // A1 (flow-E10): TZ/GS deben quedar bloqueados hasta que la tríada
    // (CI/DP/CA/HV/AG) esté completa -- antes no había ningún guard acá.
    moduloBloqueado: function (id) {
      return (id === 'TZ' || id === 'GS') && !this.triadaCompleta;
    },
    irAModulo: function (id) {
      if (this.moduloBloqueado(id)) {
        this.mensajeModuloBloqueado = 'Completa primero los módulos de Contacto, Datos personales, Contacto de apoyo, Hechos y Agresor para habilitar ' + (id === 'TZ' ? 'Tamizaje' : 'Gestión Salvia') + '.';
        return;
      }
      this.mensajeModuloBloqueado = '';
      this.moduloActivo = id;
      this.mensajeGuardado = '';
    },

    // ── Tamizaje ──────────────────────────────────────────────────────────
    tzCampoVisible: function (campo) {
      return campo.inicial || this.tzRevelarResto;
    },
    tzToggle: function (id, val) {
      this.answers[id] = val ? 'true' : 'false';
      var self = this;
      var algunaSi = TAMIZAJE_GENERALES_V3.some(function (c) { return c.inicial && self.answers[c.id] === 'true'; });
      this.tzRevelarResto = algunaSi;
    },
    calcularEsOsigd: function () {
      var orientacion = this.etiquetaSeleccionada('DP_9').toLowerCase();
      var identidad = this.etiquetaSeleccionada('DP_8').toLowerCase();
      return (orientacion && orientacion.indexOf('heterosexual') === -1) ||
        identidad.indexOf('trans') >= 0 || identidad.indexOf('no binaria') >= 0;
    },
    calcularWasPartner: function () {
      var et = this.etiquetaSeleccionada('AG_3').toLowerCase();
      return et.indexOf('pareja') >= 0;
    },
    guardarTamizaje: function () {
      var self = this;
      var wasPartner = this.calcularWasPartner();
      var esOsigd = this.calcularEsOsigd();
      var subtipoFem = this.multipleIncluyeTexto('HV_10', 'tentativa feminicidio');
      var subtipoHom = this.multipleIncluyeTexto('HV_10', 'tentativa homicidio');

      var tzAnswers = {};
      TAMIZAJE_GENERALES_V3.forEach(function (c) { tzAnswers[c.id] = self.answers[c.id] === 'true' ? 'true' : 'false'; });
      tzAnswers['TZ_8'] = this.answers['TZ_8'] === 'true' ? 'true' : 'false';
      var lista = wasPartner ? TAMIZAJE_PAREJA_V3 : TAMIZAJE_NO_PAREJA_V3;
      var prefijo = wasPartner ? 'TZP_' : 'TZNP_';
      lista.forEach(function (_, idx) {
        var id = prefijo + (idx + 1);
        tzAnswers[id] = self.answers[id] === 'true' ? 'true' : 'false';
      });

      this.tzGuardando = true;
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/tamizaje', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          answers: tzAnswers, wasPartner: wasPartner, esOsigd: esOsigd,
          subtipoFeminicidio: subtipoFem, subtipoHomicidioAbyecto: subtipoHom,
        }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) { self.tzResultado = data; })
        .finally(function () { self.tzGuardando = false; });
    },
    // D10 (flow-E11, "alertas de riesgo/ideación suicida"): antes
    // tzResultado solo se mostraba como informativo, sin resaltar los casos
    // que de verdad requieren acción inmediata.
    riesgoEsAlto: function () {
      return !!this.tzResultado && (this.tzResultado.riskLevel === 'alto' || this.tzResultado.riskLevel === 'extremo');
    },
    hayIdeacionSuicida: function () {
      return this.answers['TZP_24'] === 'true' || this.answers['TZNP_24'] === 'true';
    },
    tzWasPartner: function () { return this.calcularWasPartner(); },
    tzListaPareja: function () { return TAMIZAJE_PAREJA_V3; },
    tzListaNoPareja: function () { return TAMIZAJE_NO_PAREJA_V3; },

    // ── GS § B remisión ──────────────────────────────────────────────────
    equipoMeta: function (code) {
      return EQUIPOS_REMISION_V3.filter(function (e) { return e.value === code; })[0];
    },
    icodeForCode: function (enumCategory, code) {
      var opt = (this.enumCatalog[enumCategory] || []).filter(function (o) { return o.code === code; })[0];
      return opt ? opt.icode : null;
    },
    criterioPuntaje: function (categoria, codigosSeleccionados) {
      var meta = CRITERIOS_PUNTAJE_V3[categoria];
      if (!meta) return { puntos: 0, minimo: 0, cumple: codigosSeleccionados.length > 0 };
      var puntos = 0;
      codigosSeleccionados.forEach(function (c) { puntos += meta.puntos[c] || 0; });
      return { puntos: puntos, minimo: meta.minimo, cumple: puntos >= meta.minimo };
    },
    cumpleCriterioEquipo: function (equipo) {
      if (!equipo.criterioFieldId) return true;
      var campo = this.camposPorId[equipo.criterioFieldId];
      var opts = this.opcionesPara(campo);
      var seleccionados = (this.answers[equipo.criterioFieldId] || '').split(',').filter(function (v) { return v; });
      var codigos = seleccionados.map(function (icode) {
        var o = opts.filter(function (o) { return o.icode === icode; })[0];
        return o ? o.code : null;
      }).filter(function (c) { return c; });
      if (equipo.value === 'ap' || equipo.value === 'ah') {
        if (codigos.indexOf('ob') === -1) return false;
      }
      if (equipo.criterioCategory && CRITERIOS_PUNTAJE_V3[equipo.criterioCategory]) {
        var puntuables = codigos.filter(function (c) { return c !== 'ob'; });
        return this.criterioPuntaje(equipo.criterioCategory, puntuables).cumple;
      }
      return codigos.length > 0;
    },

    // ── GS § C/D -- buscador geográfico compartido ──────────────────────────
    // Reconstruido sobre entity_branch real (mismos 3 endpoints que usa
    // case-entities.js) -- ver registro-v3-pendientes.md § B.
    cargarGeoDepartamentos: function () {
      if (this.geoDepartamentos.length) return;
      var self = this;
      fetch('/api/v1/locations/departments')
        .then(function (r) { return r.json(); })
        .then(function (data) { self.geoDepartamentos = Array.isArray(data) ? data : []; })
        .catch(function () { self.geoDepartamentos = []; });
    },
    recargarEntidadesBarreras: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/' + this.icode)
        .then(function (r) { return r.json(); })
        .then(function (data) {
          self.entidades = data.entities || [];
          self.barreras = data.barriers || [];
        });
    },
    ubicacionEntidad: function (e) {
      return [e.departmentName, e.cityName, e.townName].filter(function (v) { return v; }).join(' · ');
    },
    // entidades[].sector viene de entity.entity_sector (código interno de 2
    // letras: he/js/pt/os -- ver entity_case_repository.go), mientras que
    // barreras[].sector es el vocabulario real por palabra que manda este
    // mismo formulario (salud/justicia/proteccion/otras_instituciones -- ver
    // VictimCaseBarrierInput.SectorCode). Dos catálogos distintos para el
    // mismo concepto -- este mapa cubre ambos para la tabla "ya agregadas".
    sectorLabelEntidadBarrera: function (code) {
      var labels = {
        he: 'Salud', js: 'Justicia', pt: 'Protección', os: 'Otras instituciones',
        salud: 'Salud', justicia: 'Justicia', proteccion: 'Protección', otras_instituciones: 'Otras instituciones',
      };
      return labels[code] || code || '—';
    },

    // ── GS § C Entidades ──────────────────────────────────────────────────
    onCambiaDeptoEntidad: function () {
      this.nuevaEntidad.ciudad = ''; this.nuevaEntidad.municipio = ''; this.nuevaEntidad.EntityBranchID = '';
      this.entidadCiudades = []; this.entidadMunicipios = []; this.entidadesDisponibles = []; this.obligacionesEntidad = [];
      if (!this.nuevaEntidad.depto) return;
      var self = this;
      fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(this.nuevaEntidad.depto))
        .then(function (r) { return r.json(); })
        .then(function (data) { self.entidadCiudades = Array.isArray(data) ? data : []; })
        .catch(function () { self.entidadCiudades = []; });
    },
    onCambiaCiudadEntidad: function () {
      this.nuevaEntidad.municipio = ''; this.nuevaEntidad.EntityBranchID = '';
      this.entidadMunicipios = []; this.entidadesDisponibles = []; this.obligacionesEntidad = [];
      if (!this.nuevaEntidad.ciudad) return;
      var self = this;
      fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(this.nuevaEntidad.ciudad))
        .then(function (r) { return r.json(); })
        .then(function (data) { self.entidadMunicipios = Array.isArray(data) ? data : []; })
        .catch(function () { self.entidadMunicipios = []; });
    },
    onCambiaMunicipioEntidad: function () { this.buscarEntidadesEntidad(); },
    onCambiaSectorEntidad: function () { this.buscarEntidadesEntidad(); },
    buscarEntidadesEntidad: function () {
      this.nuevaEntidad.EntityBranchID = ''; this.entidadesDisponibles = [];
      this.obligacionesEntidad = []; this.nuevaEntidad.CompletedActions = []; this.nuevaEntidad.PendingActions = [];
      if (!this.nuevaEntidad.municipio || !this.nuevaEntidad.sector) return;
      var self = this;
      self.entidadBuscando = true;
      fetch('/api/v1/entity-branches?town_code=' + encodeURIComponent(this.nuevaEntidad.municipio) + '&sector=' + encodeURIComponent(this.nuevaEntidad.sector))
        .then(function (r) { return r.json(); })
        .then(function (data) { self.entidadesDisponibles = Array.isArray(data) ? data : []; self.entidadBuscando = false; })
        .catch(function () { self.entidadesDisponibles = []; self.entidadBuscando = false; });
    },
    onCambiaEntidadSeleccionada: function () {
      this.obligacionesEntidad = []; this.nuevaEntidad.CompletedActions = []; this.nuevaEntidad.PendingActions = [];
      if (!this.nuevaEntidad.EntityBranchID) return;
      var self = this;
      self.obligacionesCargando = true;
      fetch('/api/v1/entity-branches/' + encodeURIComponent(this.nuevaEntidad.EntityBranchID) + '/obligations')
        .then(function (r) { return r.ok ? r.json() : []; })
        .then(function (data) {
          // "Otra" no tiene captura de texto libre en V3 todavía -- se excluye
          // del checklist, igual que case-entities.js.
          self.obligacionesEntidad = (Array.isArray(data) ? data : []).filter(function (o) { return o.value !== 'otra'; });
          self.obligacionesCargando = false;
        })
        .catch(function () { self.obligacionesEntidad = []; self.obligacionesCargando = false; });
    },
    agregarEntidad: function () {
      var self = this;
      if (!this.nuevaEntidad.EntityBranchID) return;
      var branchId = Number(this.nuevaEntidad.EntityBranchID);
      var body = {
        Kind: this.nuevaEntidad.Kind,
        EntityBranchID: branchId,
        FileNumber: this.nuevaEntidad.FileNumber,
        ActivationChannel: this.nuevaEntidad.Kind === 'activacion' ? this.nuevaEntidad.ActivationChannel : '',
        CompletedActions: this.nuevaEntidad.Kind === 'acudida' ? this.nuevaEntidad.CompletedActions : [],
        PendingActions: this.nuevaEntidad.PendingActions,
        OtherInfo: this.nuevaEntidad.OtherInfo,
      };
      // Depto/ciudad/municipio/sector usados para encontrarla -- GS § D
      // (Barreras) los reusa en vez de repetir la búsqueda (instrucción
      // 2026-10-15, ver abrirBarreraParaEntidad).
      var geo = { depto: this.nuevaEntidad.depto, ciudad: this.nuevaEntidad.ciudad, municipio: this.nuevaEntidad.municipio, sector: this.nuevaEntidad.sector };
      self.errorEntidad = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/entidades', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
      })
        .then(function (r) { return r.json().then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) {
          if (!res.ok) { self.errorEntidad = res.data.error || 'No se pudo agregar la entidad'; return; }
          self.entidadGeoPorBranch[branchId] = geo;
          self.nuevaEntidad = self.entidadVacia();
          self.entidadCiudades = []; self.entidadMunicipios = []; self.entidadesDisponibles = []; self.obligacionesEntidad = [];
          self.recargarEntidadesBarreras();
        });
    },
    eliminarEntidad: function (relId) {
      if (!window.confirm('¿Eliminar esta entidad del registro? También se eliminarán sus barreras asociadas.')) return;
      var self = this;
      fetch('/api/v1/victim-cases-v3/entidades/' + relId, { method: 'DELETE' })
        .then(function () { self.recargarEntidadesBarreras(); });
    },

    // ── GS § D Barreras -- sub-repeater de Entidades (instrucción 2026-10-15) ──
    // Cada barrera cuelga de una entidad ya agregada: NO repite la búsqueda
    // geográfica, hereda EntityBranchID/sector/ubicación de esa entidad --
    // igual que "Gestión institucional" asocia barrier_v2 a un entity_case
    // (misma pareja case_id+entity_branch_id, ver entity_case_repository.go
    // entityCaseWithRelationsSQL § barrerasActivasCount). Verificado contra
    // el esquema real de dinamic-form (tabla repeater_group): el form de
    // Seguimiento en sí mismo tiene "Barreras identificadas" y "Entidades a
    // las que ya acudió" como dos repeaters HERMANOS e independientes (cada
    // uno con su propio form_section_id, sin relación padre/hijo) -- la
    // asociación ahí es solo la coincidencia implícita de columnas, nunca un
    // repeater anidado de verdad. El anidamiento real (0, 1 o más barreras
    // por entidad, sin repetir el buscador) es una decisión propia de V3,
    // no una réplica de cómo Seguimiento arma su UI.
    barrerasDeEntidad: function (e) {
      return this.barreras.filter(function (b) { return b.entityBranchId === e.entityBranchId; });
    },
    barreraCategoriaSector: function () {
      return SECTOR_A_CATEGORIA_BARRERA_V3[this.nuevaBarrera.SectorCode] || null;
    },
    opcionesBarreraSector: function () {
      var cat = this.barreraCategoriaSector();
      if (!cat) return [];
      var self = this;
      return (this.enumCatalog[cat] || []).map(function (o) {
        return { code: o.code, icode: o.icode, label: (ENUM_LABELS_OVERRIDE_V3[cat] && ENUM_LABELS_OVERRIDE_V3[cat][o.code]) || o.label };
      });
    },
    opcionesTransversal: function (enumCategory) {
      var self = this;
      return (this.enumCatalog[enumCategory] || []).map(function (o) {
        return { code: o.code, icode: o.icode, label: (ENUM_LABELS_OVERRIDE_V3[enumCategory] && ENUM_LABELS_OVERRIDE_V3[enumCategory][o.code]) || o.label };
      });
    },
    toggleEnArray: function (arr, value) {
      var idx = arr.indexOf(value);
      if (idx >= 0) arr.splice(idx, 1); else arr.push(value);
    },
    categoriasTransversales: function () { return CATEGORIAS_BARRERA_TRANSVERSAL_V3; },
    opcionesGestionBarrera: function () { return this.opcionesTransversal(GESTION_BARRERA_ENUM_CATEGORY_V3); },
    // Traduce el sector interno de entity (he/js/pt/os, lo que trae `e.sector`
    // de EntityCaseWithRelations) al vocabulario real que espera BarrierV2 --
    // inverso del sectorMap de entity_branch_api_controller.go.
    sectorInternoAPalabra: function (code) {
      var map = { he: 'salud', js: 'justicia', pt: 'proteccion', os: 'otras_instituciones' };
      return map[code] || '';
    },
    abrirBarreraParaEntidad: function (e) {
      this.nuevaBarrera = this.barreraVacia();
      this.nuevaBarrera.EntityBranchID = e.entityBranchId;
      var geo = this.entidadGeoPorBranch[e.entityBranchId];
      if (geo) {
        this.nuevaBarrera.Department = geo.depto; this.nuevaBarrera.City = geo.ciudad; this.nuevaBarrera.Town = geo.municipio;
        this.nuevaBarrera.SectorCode = geo.sector;
      } else {
        // Entidad agregada en una sesión anterior (recargada del backend) --
        // no se cacheó la búsqueda que la encontró. Department/City/Town
        // quedan vacíos (son solo "dónde se identificó", redundante con la
        // ubicación propia de la sede -- ver BarrierV2.EntityBranchID); el
        // sector se traduce del código interno igual, para que el checklist
        // de "barreras específicas" y competencias cargue bien.
        this.nuevaBarrera.SectorCode = this.sectorInternoAPalabra(e.sector);
      }
      this.barreraParaEntidad = e.entityBranchId;
      this.errorBarrera = '';
      this.buscarCompetencias();
    },
    cerrarBarreraParaEntidad: function () {
      this.barreraParaEntidad = null;
      this.nuevaBarrera = this.barreraVacia();
    },
    buscarCompetencias: function () {
      this.competenciasBarrera = [];
      if (!this.nuevaBarrera.SectorCode) return;
      var self = this;
      self.competenciasCargando = true;
      fetch('/api/v1/competences?sector=' + encodeURIComponent(this.nuevaBarrera.SectorCode))
        .then(function (r) { return r.ok ? r.json() : []; })
        .then(function (data) { self.competenciasBarrera = Array.isArray(data) ? data : []; self.competenciasCargando = false; })
        .catch(function () { self.competenciasBarrera = []; self.competenciasCargando = false; });
    },
    agregarBarrera: function () {
      var self = this;
      if (!this.nuevaBarrera.EntityBranchID) return;
      var body = Object.assign({}, this.nuevaBarrera, { EntityBranchID: Number(this.nuevaBarrera.EntityBranchID) });
      self.errorBarrera = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/barreras', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
      })
        .then(function (r) { return r.json().then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) {
          if (!res.ok) { self.errorBarrera = res.data.error || 'No se pudo agregar la barrera'; return; }
          self.cerrarBarreraParaEntidad();
          self.recargarEntidadesBarreras();
        });
    },
    eliminarBarrera: function (id) {
      if (!window.confirm('¿Eliminar esta barrera del registro?')) return;
      var self = this;
      fetch('/api/v1/victim-cases-v3/barreras/' + id, { method: 'DELETE' })
        .then(function () { self.recargarEntidadesBarreras(); });
    },

    // ── Cierre ────────────────────────────────────────────────────────────
    // ── Párrafo síntesis (D3-D5, 2026-10-15) ────────────────────────────────
    // Reescrito para seguir registro-v3-parrafo-sintesis.md literal (obs.
    // #18 del cliente, la fuente de mayor prioridad de todo Registro V3).
    // Antes solo tenía una versión mínima de § A + un dump crudo de HV_1 --
    // § B (datos personales) y § D (agresor) no existían, y el párrafo
    // mezclaba resultado de Tamizaje/remisión (que NO son parte de este
    // spec) en su lugar. Esa info de Tamizaje/remisión ya se muestra aparte
    // (alertas de TZ, D10; mensajeCierre de Finalizar, C2/C3) -- no hace
    // falta repetirla acá.
    //
    // etiquetaExcluye: true si la etiqueta resuelta de fieldId contiene
    // alguno de los fragmentos de exclusión (la plantilla del cliente
    // describe las excepciones por texto, no por código de enum).
    etiquetaExcluye: function (fieldId, exclusiones) {
      var et = this.etiquetaSeleccionada(fieldId).toLowerCase();
      if (!et) return true;
      return exclusiones.some(function (frag) { return et.indexOf(frag.toLowerCase()) >= 0; });
    },
    // Encontrado en vivo con Playwright (2026-10-15): GetAnswers devuelve las
    // fechas en el formato crudo de Go ("2026-10-01 00:00:00 +0000 UTC") una
    // vez que un Caso se reabre desde el backend -- antes de esta ronda
    // nadie mostraba una fecha cruda en prosa (solo se usaban como `value`
    // de un <input type="date">, que el navegador tolera igual). sintesisC
    // es la primera vez que se concatena una fecha directo en una frase, así
    // que es la primera vez que se nota. Se arregla acá, en el borde de
    // visualización, no en GetAnswers (que sigue alimentando el <input
    // type="date"> correctamente tal cual).
    limpiarFecha: function (raw) {
      var m = /^\d{4}-\d{2}-\d{2}/.exec(raw || '');
      return m ? m[0] : (raw || '');
    },
    calcularEdad: function (fechaISO) {
      if (!fechaISO) return null;
      var nacimiento = new Date(fechaISO);
      if (isNaN(nacimiento.getTime())) return null;
      var hoy = new Date();
      var edad = hoy.getFullYear() - nacimiento.getFullYear();
      var noHaCumplido = (hoy.getMonth() < nacimiento.getMonth()) ||
        (hoy.getMonth() === nacimiento.getMonth() && hoy.getDate() < nacimiento.getDate());
      if (noHaCumplido) edad--;
      return edad >= 0 ? edad : null;
    },
    // § A — Contacto e identificación. Se omite si no hay nombres/apellidos.
    sintesisA: function () {
      if (!this.answers['CI_2']) return '';
      var tratamientoTexto = { 'Sra.': 'La Señora', 'Sr.': 'El Señor', 'Sre.': 'Le Señore' }[this.tratamiento()] || '';
      var identitario = this.answers['CI_3'] ? this.answers['CI_3'] + ' ' : '';
      var frase = (tratamientoTexto ? tratamientoTexto + ' ' : '') + identitario + this.answers['CI_2'] +
        ', identificada/e/o con ' + (this.etiquetaSeleccionada('CI_5') || 'documento') + ' Número ' + (this.answers['CI_6'] || 's/n') + '.';
      if (this.answers['CI_7'] === 'true') frase += ' Reconoce discapacidad.';
      if (this.answers['CI_9'] === 'true') frase += ' Requiere intérprete de idiomas.';
      return frase;
    },
    // § B — Datos personales (4 bloques independientes).
    sintesisB: function () {
      var self = this;
      var bloques = [];
      var nombre = this.answers['CI_2'] || '';

      var edad = this.calcularEdad(this.answers['DP_1']);
      if (edad !== null) {
        var nacCode = this.codigoSeleccionado('DP_5');
        var nacionalidad = nacCode === 'co' ? 'Colombiana(o)' : nacCode === 'ap' ? 'Apátrida' :
          nacCode === 'ex' ? 'de ' + (this.etiquetaSeleccionada('DP_6') || 'nacionalidad extranjera') : '';
        var b1 = nombre + ' tiene ' + edad + ' años' + (nacionalidad ? ', nacional ' + nacionalidad : '') + '.';
        if (this.answers['DP_27']) b1 += ' Reside en ' + this.answers['DP_27'] + (this.answers['DP_28'] ? ' en ' + this.answers['DP_28'] : '') + (this.etiquetaSeleccionada('DP_30') ? ' en ' + this.etiquetaSeleccionada('DP_30') : '') + '.';
        bloques.push(b1);
      }

      if (!this.etiquetaExcluye('DP_8', ['otro'])) {
        var b2 = 'Se reconoce ' + this.etiquetaSeleccionada('DP_8');
        if (!this.etiquetaExcluye('DP_10', [])) b2 += ', fue asignado al nacer ' + this.etiquetaSeleccionada('DP_10');
        if (!this.etiquetaExcluye('DP_9', ['otro', 'no estoy segura', 'no informa'])) b2 += ' con orientación sexual ' + this.etiquetaSeleccionada('DP_9');
        bloques.push(b2 + '.');
      }

      var b3partes = [];
      if (!this.etiquetaExcluye('DP_12', ['no se auto reconoce', 'no informa'])) b3partes.push('Se identifica ' + this.etiquetaSeleccionada('DP_12'));
      if (this.answers['DP_14'] === 'true') b3partes.push((b3partes.length ? ' y se reconoce' : 'Se reconoce') + ' como persona campesina');
      var poblaciones = (this.answers['DP_11'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['DP_11']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['DP_11'], o) : ''; })
        .filter(function (l) { return l && l.toLowerCase().indexOf('ninguno') === -1; });
      if (b3partes.length) {
        var b3 = b3partes.join('') + '.';
        if (poblaciones.length) b3 += ' Pertenece a población de especial protección constitucional ' + poblaciones.join(', ') + '.';
        bloques.push(b3);
      } else if (poblaciones.length) {
        bloques.push('Pertenece a población de especial protección constitucional ' + poblaciones.join(', ') + '.');
      }

      var b4partes = [];
      if (this.answers['DP_15']) b4partes.push('Persona ' + this.etiquetaSeleccionada('DP_15'));
      if (!this.etiquetaExcluye('DP_16', ['no informa'])) {
        var esc = this.etiquetaSeleccionada('DP_16');
        b4partes.push((b4partes.length ? ' ' : 'Persona ') + (esc.toLowerCase() === 'sin educación formal' ? esc : 'con ' + esc));
      }
      if (b4partes.length) {
        var b4 = b4partes.join('') + '.';
        if (this.answers['DP_18'] && this.answers['DP_17']) b4 += ' ' + this.etiquetaSeleccionada('DP_18') + ' en su ocupación ' + this.etiquetaSeleccionada('DP_17') + '.';
        if (!this.etiquetaExcluye('DP_23', ['no informa'])) b4 += ' Su vivienda es ' + this.etiquetaSeleccionada('DP_23');
        if (!this.etiquetaExcluye('DP_24', ['no informa', 'no sabe'])) b4 += ' en estrato ' + this.etiquetaSeleccionada('DP_24') + '.';
        else if (b4.slice(-1) !== '.') b4 += '.';
        var personasCargo = (this.answers['DP_25'] || '').split(',').filter(function (v) { return v; })
          .map(function (ic) { var o = self.opcionesPara(self.camposPorId['DP_25']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['DP_25'], o) : ''; }).filter(function (l) { return l; });
        if (personasCargo.length) b4 += ' Tiene personas o animales a su cargo: ' + personasCargo.join(', ') + '.';
        bloques.push(b4);
      }

      return bloques.join(' ');
    },
    // § C — Hechos de violencia. Se omite si no hay tipo de violencia ni relato.
    sintesisC: function () {
      if (!this.answers['HV_1'] && !this.answers['HV_9']) return '';
      var self = this;
      var tipos = (this.answers['HV_9'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['HV_9']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['HV_9'], o) : ''; }).filter(function (l) { return l; });
      var subtipos = (this.answers['HV_10'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['HV_10']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['HV_10'], o) : ''; }).filter(function (l) { return l; });
      var ambitos = (this.answers['HV_11'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['HV_11']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['HV_11'], o) : ''; }).filter(function (l) { return l; });

      var frase = 'Hecho de violencia' + (tipos.length ? ' de tipo(s) ' + tipos.join(', ') : '') +
        (subtipos.length ? ', ' + subtipos.join(', ') : '') +
        (ambitos.length ? ', ocurrido en ámbito(s) ' + ambitos.join(', ') : '') +
        (this.etiquetaSeleccionada('HV_8') ? ' en escenario ' + this.etiquetaSeleccionada('HV_8') : '') + '.';
      if (this.answers['HV_2']) {
        frase += ' Sucedido el ' + this.limpiarFecha(this.answers['HV_2']) + (this.answers['HV_3'] ? ' ' + this.answers['HV_3'] : '') +
          (this.answers['HV_4'] ? ' en ' + this.answers['HV_4'] : '') + (this.answers['HV_5'] ? ' en ' + this.answers['HV_5'] : '') +
          (this.etiquetaSeleccionada('HV_6') ? ' en ' + this.etiquetaSeleccionada('HV_6') : '') + '.';
      }
      if (this.answers['HV_13'] === 'true') frase += ' Motivado por razones de género.';
      if (this.etiquetaSeleccionada('HV_14')) frase += ' ' + this.etiquetaSeleccionada('HV_14') + '.';
      if (this.answers['HV_1']) frase += ' ' + this.answers['HV_1'];
      return frase;
    },
    // § D — Características del agresor. Se omite si no hay relación ni número de agresores.
    sintesisD: function () {
      if (!this.answers['AG_3'] && !this.answers['AG_1']) return '';
      var frase = '';
      if (!this.etiquetaExcluye('AG_1', ['no se puede determinar', 'no informa'])) {
        frase = 'Responsabiliza a ' + this.etiquetaSeleccionada('AG_1') + ' persona(s) de la agresión.';
      }
      if (this.etiquetaSeleccionada('AG_2')) frase += (frase ? ' Principalmente a' : 'A') + ' ' + this.etiquetaSeleccionada('AG_2').toLowerCase() + (this.answers['AG_7'] ? ', de nombre ' + this.answers['AG_7'] : '');
      if (this.etiquetaSeleccionada('AG_6')) frase += (frase.slice(-1) === '.' ? ' ' : ', ') + this.etiquetaSeleccionada('AG_6');
      if (!this.etiquetaExcluye('AG_3', ['sin relación', 'otro'])) frase += ', con quien tiene una relación de ' + this.etiquetaSeleccionada('AG_3');
      if (!this.etiquetaExcluye('AG_4', ['otro', 'no conoce'])) frase += ', cuya ocupación es ' + this.etiquetaSeleccionada('AG_4');
      return frase.trim() + (frase.slice(-1) !== '.' ? '.' : '');
    },
    parrafoSintesis: function () {
      return [this.sintesisA(), this.sintesisB(), this.sintesisC(), this.sintesisD()].filter(function (s) { return s; }).join(' ');
    },
    // C2/C3 (flow-E15): "generar Denuncia fácil si corresponde" / "confirmar
    // remisión si aplica" -- no hay integración real con la FGN ni un
    // workflow de remisión separado de la respuesta GS_6 para reutilizar
    // (ver FinalizeResultV3 en victim_case_v3_service.go), así que se
    // implementa como una confirmación explícita al agente en vez de un
    // envío silencioso.
    confirmacionFinalizar: function (data) {
      var partes = [];
      if (data.denunciaFacilSolicitada) partes.push('Se solicitó la Denuncia fácil ante la FGN.');
      if (data.equiposRemisionCSV) {
        var opts = this.opcionesPara(this.camposPorId['GS_6']);
        var icodes = data.equiposRemisionCSV.split(',').filter(function (v) { return v; });
        var nombres = icodes.map(function (ic) {
          var o = opts.filter(function (o) { return o.icode === ic; })[0];
          return o ? o.label : ic;
        });
        if (nombres.length) partes.push('Se confirma remisión a: ' + nombres.join(', ') + '.');
      }
      return partes.join(' ');
    },
    finalizarRegistro: function () {
      var self = this;
      self.errorFinalizar = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/finalizar', { method: 'POST' })
        .then(function (r) { return r.json().then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) {
          if (res.ok) {
            self.mensajeCierre = self.confirmacionFinalizar(res.data);
            self.etapa = 'cierre';
            return;
          }
          var faltan = res.data.modulosIncompletos || [];
          var nombres = faltan.map(function (id) {
            var m = self.modulos.filter(function (x) { return x.id === id; })[0];
            return m ? m.label : id;
          });
          self.errorFinalizar = (res.data.error || 'No se pudo finalizar') + (nombres.length ? ' — falta diligenciar: ' + nombres.join(', ') : '');
        });
    },
  },
});
